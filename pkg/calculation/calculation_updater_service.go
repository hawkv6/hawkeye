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

func (service *CalculationUpdaterService) calculateTotalCost(pathResult domain.PathResult, weightTypes []helper.WeightKey) error {
	var newTotalCost float64
	if len(weightTypes) == 1 && weightTypes[0] == helper.PacketLossKey {
		newTotalCost = 1.0
	} else {
		newTotalCost = 0.0
	}
	for _, edge := range pathResult.GetEdges() {
		graphEdge := service.graph.GetEdge(edge.GetId())
		if graphEdge == nil {
			return fmt.Errorf("Path not valid, edge not found in graph: %s", edge.GetId())
		}
		if len(weightTypes) == 1 {
			if weightTypes[0] == helper.PacketLossKey {
				newTotalCost *= 1 - graphEdge.GetWeight(weightTypes[0])
			} else {
				newTotalCost += graphEdge.GetWeight(weightTypes[0])
			}
		} else if len(weightTypes) == 2 {
			newTotalCost += edge.GetWeight(weightTypes[0])*float64(helper.TwoFactorWeights[0]) + edge.GetWeight(weightTypes[1])*float64(helper.TwoFactorWeights[1])
		} else {
			newTotalCost += edge.GetWeight(weightTypes[0])*float64(helper.ThreeFactorWeights[0]) + edge.GetWeight(weightTypes[1])*float64(helper.ThreeFactorWeights[1]) + edge.GetWeight(weightTypes[2])*float64(helper.ThreeFactorWeights[2])
		}
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

func (service *CalculationUpdaterService) calculateMinimumValue(pathResult domain.PathResult, weightType helper.WeightKey) error {
	minValue := math.Inf(1)
	var bottleneckEdge graph.Edge
	edges := pathResult.GetEdges()
	for _, edge := range edges {
		graphEdge := service.graph.GetEdge(edge.GetId())
		if graphEdge == nil {
			return fmt.Errorf("Path not valid, edge not found in graph: %s", edge.GetId())
		}
		cost := graphEdge.GetWeight(weightType)
		if cost < minValue {
			bottleneckEdge = graphEdge
			minValue = cost
		}
	}
	if pathResult.GetBottleneckEdge() != bottleneckEdge {
		service.log.Debugf("Bottleneck edge of current applied path changed from %v to %v", pathResult.GetBottleneckEdge(), bottleneckEdge)
		service.log.Debugf("Bottleneck value of current applied path changed from %f to %f", pathResult.GetBottleneckValue(), minValue)
		pathResult.SetBottleneckEdge(bottleneckEdge)
		pathResult.SetBottleneckValue(minValue)
	} else if pathResult.GetBottleneckValue() != minValue {
		service.log.Debugf("Bottleneck value of current applied path changed from %f to %f", pathResult.GetBottleneckValue(), minValue)
		pathResult.SetBottleneckValue(minValue)
	}
	return nil
}

func (service *CalculationUpdaterService) updateCurrentResult(weightKeys []helper.WeightKey, calculationMode CalculationMode, currentPathResult domain.PathResult) error {
	var err error
	if calculationMode == CalculationModeSum {
		err = service.calculateTotalCost(currentPathResult, weightKeys)
	} else {
		err = service.calculateMinimumValue(currentPathResult, weightKeys[0])
	}
	if err != nil {
		return fmt.Errorf("Current Path is not valid anymore, new path will be applied: %s", err)
	}
	return nil
}

func (service *CalculationUpdaterService) handleSumCalculationMode(currentPathResult, newPathResult domain.PathResult, streamSession domain.StreamSession) (*domain.PathResult, error) {
	newPathTotalCost := newPathResult.GetTotalCost()
	currentPathTotalCost := currentPathResult.GetTotalCost()
	if newPathTotalCost < currentPathTotalCost*(1-helper.FlappingThreshold) {
		service.log.Debugf("New path will be applied, cost of new path is by more than 10 percent smaller, current: %f to new: %f", currentPathTotalCost, newPathTotalCost)
		streamSession.SetPathResult(newPathResult)
		return &newPathResult, nil
	} else {
		service.log.Debugf("No path changes, cost of new path is not smaller by more than 10 percent, current: %f to new: %f", currentPathTotalCost, newPathTotalCost)
	}
	return nil, nil
}

func (service *CalculationUpdaterService) handleNonSumCalculationMode(currentPathResult, newPathResult domain.PathResult, streamSession domain.StreamSession) (*domain.PathResult, error) {
	newPathMinimumValue := newPathResult.GetBottleneckValue()
	newPathMinimumEdge := newPathResult.GetBottleneckEdge()
	currentPathMinimumValue := currentPathResult.GetBottleneckValue()
	currentPathMinimumEdge := currentPathResult.GetBottleneckEdge()
	if newPathMinimumValue < currentPathMinimumValue*(1-helper.FlappingThreshold) {
		service.log.Debugf("Bottleneck in old path was %v with value %g: ", currentPathMinimumEdge, currentPathMinimumValue)
		service.log.Debugf("Bottleneck in new path is %v with value %g: ", newPathMinimumEdge, newPathMinimumValue)
		service.log.Debugf("New Path will be applied, bottleneck of new path is by more than 10 percent smaller, current: %f to new: %f", currentPathMinimumValue, newPathMinimumValue)
		streamSession.SetPathResult(newPathResult)
		return &newPathResult, nil
	} else {
		service.log.Debugf("No path changes, cost of new path is not smaller by more than 10 percent, current: %f to new: %f", currentPathMinimumValue, newPathMinimumValue)
	}
	return nil, nil
}

func (service *CalculationUpdaterService) AreServicesStillValid(serviceSidList []string) bool {
	service.log.Debugln("Check if current services are still valid")
	for _, sid := range serviceSidList {
		if !service.cache.DoesServiceSidExist(sid) {
			service.log.Debugf("Service SID %s not found in cache, new path will be applied", sid)
			return false
		}
	}
	return true
}

func (service *CalculationUpdaterService) handlePathChange(weightKeys []helper.WeightKey, calculationMode CalculationMode, currentPathResult, newPathResult domain.PathResult, streamSession domain.StreamSession) (*domain.PathResult, error) {
	service.log.Debugln("Better Path found, check for applicability")
	service.log.Debugln("Validate current path and its cost")

	firstIntent := streamSession.GetPathRequest().GetIntents()[0]
	if firstIntent.GetIntentType() == domain.IntentTypeSFC {
		if !service.AreServicesStillValid(currentPathResult.GetServiceSidList()) {
			streamSession.SetPathResult(newPathResult)
			return &newPathResult, nil
		}
	}

	if err := service.updateCurrentResult(weightKeys, calculationMode, currentPathResult); err != nil {
		service.log.Errorln("Current Path is not valid anymore, new path will be applied: ", err)
		streamSession.SetPathResult(newPathResult)
		return &newPathResult, nil
	}

	service.log.Debugln("Current path is still valid, check if new path is better")
	if calculationMode == CalculationModeSum {
		return service.handleSumCalculationMode(currentPathResult, newPathResult, streamSession)
	} else {
		return service.handleNonSumCalculationMode(currentPathResult, newPathResult, streamSession)
	}
}

func (service *CalculationUpdaterService) UpdateCalculation(options *CalculationUpdateOptions) (*domain.PathResult, error) {
	if !reflect.DeepEqual(options.newPathResult.GetIpv6SidAddresses(), options.currentAppliedSidList) {
		return service.handlePathChange(options.weightKeys, options.calculationMode, options.currentPathResult, options.newPathResult, options.streamSession)
	} else {
		service.log.Debugln("No changes in path detected, update current path with new path cost")
		if err := service.updateCurrentResult(options.weightKeys, options.calculationMode, options.currentPathResult); err != nil {
			return nil, err
		}
	}
	service.log.Debugln("No path changes, current path is still valid")
	return nil, nil
}
