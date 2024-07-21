package calculation

import (
	"fmt"
	"math"
	"reflect"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type CalculationUpdaterService struct {
	log   *logrus.Entry
	cache cache.Cache
	graph graph.Graph
}

func NewCalculationUpdaterService(cache cache.Cache, graph graph.Graph) *CalculationUpdaterService {
	return &CalculationUpdaterService{
		log:   logging.DefaultLogger.WithField("subsystem", subsystem),
		cache: cache,
		graph: graph,
	}
}

func (*CalculationUpdaterService) getInitialTotalCost(weightTypes []helper.WeightKey) float64 {
	if len(weightTypes) == 1 && weightTypes[0] == helper.PacketLossKey {
		return 1.0
	}
	return 0.0
}

func (service *CalculationUpdaterService) getUpdatedEdge(edge graph.Edge) (graph.Edge, float64, error) {
	updatedEdge := service.graph.GetEdge(edge.GetId())
	if updatedEdge == nil {
		return nil, 0, fmt.Errorf("Path not valid, edge not found in updated graph: %s", edge.GetId())
	}
	return updatedEdge, 0, nil
}

func (service *CalculationUpdaterService) getUpdatedTotalCost(edges []graph.Edge, weightTypes []helper.WeightKey, newTotalCost float64) (float64, error) {
	for _, edge := range edges {
		updatedEdge, returnValue, err := service.getUpdatedEdge(edge)
		if err != nil {
			return returnValue, err
		}
		if len(weightTypes) == 1 {
			if weightTypes[0] == helper.PacketLossKey {
				newTotalCost *= 1 - updatedEdge.GetWeight(weightTypes[0])
			} else {
				newTotalCost += updatedEdge.GetWeight(weightTypes[0])
			}
		} else if len(weightTypes) == 2 {
			newTotalCost += edge.GetWeight(weightTypes[0])*float64(helper.TwoFactorWeights[0]) + edge.GetWeight(weightTypes[1])*float64(helper.TwoFactorWeights[1])
		} else {
			newTotalCost += edge.GetWeight(weightTypes[0])*float64(helper.ThreeFactorWeights[0]) + edge.GetWeight(weightTypes[1])*float64(helper.ThreeFactorWeights[1]) + edge.GetWeight(weightTypes[2])*float64(helper.ThreeFactorWeights[2])
		}
	}
	return newTotalCost, nil
}

func (service *CalculationUpdaterService) updateTotalCost(pathResult domain.PathResult, weightTypes []helper.WeightKey) error {
	newTotalCost := service.getInitialTotalCost(weightTypes)

	newTotalCost, err := service.getUpdatedTotalCost(pathResult.GetEdges(), weightTypes, newTotalCost)
	if err != nil {
		return err
	}
	if len(weightTypes) == 1 && weightTypes[0] == helper.PacketLossKey {
		newTotalCost = (1 - newTotalCost) * 100
	}
	currentTotalCost := pathResult.GetTotalCost()
	if currentTotalCost != newTotalCost {
		service.log.Debugf("Total cost of current applied path changed from %f to %f", currentTotalCost, newTotalCost)
		pathResult.SetTotalCost(newTotalCost)
	}
	return nil
}

func (service *CalculationUpdaterService) getUpdatedBottleneckValues(edges []graph.Edge, weightType helper.WeightKey, minValue float64, bottleneckEdge graph.Edge) (float64, graph.Edge, error) {
	for _, edge := range edges {
		updatedEdge, returnValue, err := service.getUpdatedEdge(edge)
		if err != nil {
			return returnValue, nil, err
		}
		cost := updatedEdge.GetWeight(weightType)
		if cost < minValue {
			bottleneckEdge = updatedEdge
			minValue = cost
		}
	}
	return minValue, bottleneckEdge, nil
}

func (service *CalculationUpdaterService) updateBottleneckValues(pathResult domain.PathResult, bottleneckEdge graph.Edge, minValue float64) {
	if pathResult.GetBottleneckEdge() != bottleneckEdge {
		service.log.Debugf("Bottleneck edge of current applied path changed from %v to %v", pathResult.GetBottleneckEdge(), bottleneckEdge)
		service.log.Debugf("Bottleneck value of current applied path changed from %f to %f", pathResult.GetBottleneckValue(), minValue)
		pathResult.SetBottleneckEdge(bottleneckEdge)
		pathResult.SetBottleneckValue(minValue)
	} else if pathResult.GetBottleneckValue() != minValue {
		service.log.Debugf("Bottleneck value of current applied path changed from %f to %f", pathResult.GetBottleneckValue(), minValue)
		pathResult.SetBottleneckValue(minValue)
	}
}

func (service *CalculationUpdaterService) updateMinimumValue(pathResult domain.PathResult, weightType helper.WeightKey) error {
	minValue := math.Inf(1)
	var bottleneckEdge graph.Edge
	minValue, bottleneckEdge, err := service.getUpdatedBottleneckValues(pathResult.GetEdges(), weightType, minValue, bottleneckEdge)
	if err != nil {
		return err
	}
	service.updateBottleneckValues(pathResult, bottleneckEdge, minValue)
	return nil
}

func (service *CalculationUpdaterService) updateCurrentResult(weightKeys []helper.WeightKey, calculationMode CalculationMode, currentPathResult domain.PathResult) error {
	var err error
	if calculationMode == CalculationModeSum {
		err = service.updateTotalCost(currentPathResult, weightKeys)
	} else {
		err = service.updateMinimumValue(currentPathResult, weightKeys[0])
	}
	if err != nil {
		return fmt.Errorf("Current Path is not valid anymore, new path will be applied: %s", err)
	}
	return nil
}

func (service *CalculationUpdaterService) updatePathIfCostImproved(currentPathResult, newPathResult domain.PathResult, streamSession domain.StreamSession) domain.PathResult {
	newPathTotalCost := newPathResult.GetTotalCost()
	currentPathTotalCost := currentPathResult.GetTotalCost()
	if newPathTotalCost < currentPathTotalCost*(1-helper.FlappingThreshold) {
		service.log.Debugf("New path will be applied, cost of new path is by more than 10 percent smaller, current: %f to new: %f", currentPathTotalCost, newPathTotalCost)
		streamSession.SetPathResult(newPathResult)
		return newPathResult
	}
	service.log.Debugf("No path changes, cost of new path is not smaller by more than 10 percent, current: %f to new: %f", currentPathTotalCost, newPathTotalCost)
	return nil
}

func (service *CalculationUpdaterService) updatePathIfMinimumImproved(currentPathResult, newPathResult domain.PathResult, streamSession domain.StreamSession) domain.PathResult {
	newPathMinimumValue := newPathResult.GetBottleneckValue()
	currentPathMinimumValue := currentPathResult.GetBottleneckValue()
	if newPathMinimumValue < currentPathMinimumValue*(1-helper.FlappingThreshold) {
		newPathMinimumEdge := newPathResult.GetBottleneckEdge()
		currentPathMinimumEdge := currentPathResult.GetBottleneckEdge()
		service.log.Debugf("Bottleneck in old path was %v with value %g: ", currentPathMinimumEdge, currentPathMinimumValue)
		service.log.Debugf("Bottleneck in new path is %v with value %g: ", newPathMinimumEdge, newPathMinimumValue)
		service.log.Debugf("New Path will be applied, bottleneck of new path is by more than 10 percent smaller, current: %f to new: %f", currentPathMinimumValue, newPathMinimumValue)
		streamSession.SetPathResult(newPathResult)
		return newPathResult
	}
	service.log.Debugf("No path changes, cost of new path is not smaller by more than 10 percent, current: %f to new: %f", currentPathMinimumValue, newPathMinimumValue)
	return nil
}

func (service *CalculationUpdaterService) currentServicesStillValid(serviceSidList []string) bool {
	service.log.Debugln("Check if current services are still valid")
	for _, sid := range serviceSidList {
		if !service.cache.DoesServiceSidExist(sid) {
			service.log.Debugf("Service SID %s not found in cache, new path will be applied", sid)
			return false
		}
	}
	return true
}

func (service *CalculationUpdaterService) currentServicesNotValidAnymore(firstIntent domain.Intent, currentPathResult domain.PathResult) bool {
	if firstIntent.GetIntentType() == domain.IntentTypeSFC {
		if !service.currentServicesStillValid(currentPathResult.GetServiceSidList()) {
			service.log.Debugln("Current services are not valid anymore, new path will be applied")
			return true
		}
	}
	service.log.Debugln("Current services are still valid, check if new path is better")
	return false
}

func (service *CalculationUpdaterService) currentPathNotValidAnymore(weightKeys []helper.WeightKey, calculationMode CalculationMode, currentPathResult domain.PathResult) bool {
	if err := service.updateCurrentResult(weightKeys, calculationMode, currentPathResult); err != nil {
		service.log.Errorln("Current Path is not valid anymore, new path will be applied: ", err)
		return true
	}
	return false
}

func (service *CalculationUpdaterService) handlePathChange(weightKeys []helper.WeightKey, calculationMode CalculationMode, currentPathResult, newPathResult domain.PathResult, streamSession domain.StreamSession) domain.PathResult {
	service.log.Debugln("Better Path found, check for applicability")
	service.log.Debugln("Validate current path and its cost")

	if service.currentServicesNotValidAnymore(streamSession.GetPathRequest().GetIntents()[0], currentPathResult) || service.currentPathNotValidAnymore(weightKeys, calculationMode, currentPathResult) {
		streamSession.SetPathResult(newPathResult)
		return newPathResult
	}

	if calculationMode == CalculationModeSum {
		return service.updatePathIfCostImproved(currentPathResult, newPathResult, streamSession)
	} else {
		return service.updatePathIfMinimumImproved(currentPathResult, newPathResult, streamSession)
	}
}

func (service *CalculationUpdaterService) UpdateCalculation(options *CalculationUpdateOptions) (domain.PathResult, error) {
	if !reflect.DeepEqual(options.newPathResult.GetIpv6SidAddresses(), options.currentAppliedSidList) {
		return service.handlePathChange(options.weightKeys, options.calculationMode, options.currentPathResult, options.newPathResult, options.streamSession), nil
	} else {
		service.log.Debugln("No changes in path detected, update current path with new path cost")
		if err := service.updateCurrentResult(options.weightKeys, options.calculationMode, options.currentPathResult); err != nil {
			return nil, err
		}
	}
	service.log.Debugln("No path changes, current path is still valid")
	return nil, nil
}
