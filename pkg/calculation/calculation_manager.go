package calculation

import (
	"fmt"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type CalculationManager struct {
	log                    *logrus.Entry
	cache                  cache.Cache
	graph                  graph.Graph
	calculationSetup       CalculationSetup
	calculationTransformer CalculationTransformer
	calculationUpdater     CalculationUpdater
	calculation            Calculation
	algorithm              uint32
}

func NewCalculationManager(cache cache.Cache, graph graph.Graph, calculationSetup CalculationSetup, calcultionTransformer CalculationTransformer, calculationUpdater CalculationUpdater) *CalculationManager {
	return &CalculationManager{
		log:                    logging.DefaultLogger.WithField("subsystem", subsystem),
		cache:                  cache,
		graph:                  graph,
		calculationSetup:       calculationSetup,
		calculationTransformer: calcultionTransformer,
		calculationUpdater:     calculationUpdater,
	}
}

func (manager *CalculationManager) lockElements() {
	manager.log.Debugln("Locking cache and graph mutexes")
	manager.graph.Lock()
	manager.cache.Lock()
}

func (manager *CalculationManager) unlockElements() {
	manager.log.Debugln("Unlocking cache and graph mutexes")
	manager.graph.Unlock()
	manager.cache.Unlock()
}

func (manager *CalculationManager) getGraphAndAlgorithm(graph graph.Graph, firstIntent domain.Intent) (graph.Graph, uint32) {
	if firstIntent.GetIntentType() == domain.IntentTypeFlexAlgo {
		algorithm := uint32(firstIntent.GetValues()[0].GetNumberValue())
		return graph.GetSubGraph(algorithm), algorithm
	} else {
		return graph, 0
	}
}

func (manager *CalculationManager) setUpCalculation(pathRequest domain.PathRequest) error {
	calculationOptions, err := manager.calculationSetup.PerformSetup(pathRequest)
	if err != nil {
		return fmt.Errorf("error performing setup: %w", err)
	}

	intents := pathRequest.GetIntents()
	calculationOptions.graph, manager.algorithm = manager.getGraphAndAlgorithm(manager.graph, manager.getFirstNonSfcIntent(intents))

	firstIntent := intents[0]
	switch firstIntent.GetIntentType() {
	case domain.IntentTypeSFC:
		return manager.setupServiceFunctionChainCalculation(firstIntent, calculationOptions)
	default:
		manager.calculation = NewShortestPathCalculation(calculationOptions)
	}
	return nil
}

func (manager *CalculationManager) getFirstNonSfcIntent(intents []domain.Intent) domain.Intent {
	if len(intents) > 1 && intents[0].GetIntentType() == domain.IntentTypeSFC {
		return intents[1]
	}
	return intents[0]
}

func (manager *CalculationManager) setupServiceFunctionChainCalculation(intent domain.Intent, calculationOptions *CalculationOptions) error {
	sfcCalculationOptions, err := manager.calculationSetup.PerformServiceFunctionChainSetup(intent, manager.algorithm)
	if err != nil {
		return fmt.Errorf("error setting up service function chain: %w", err)
	}
	manager.calculation = NewServiceFunctionChainCalculation(calculationOptions, sfcCalculationOptions)
	return nil
}

func (manager *CalculationManager) CalculateBestPath(pathRequest domain.PathRequest) (domain.PathResult, error) {
	manager.lockElements()
	defer manager.unlockElements()

	err := manager.setUpCalculation(pathRequest)
	if err != nil {
		return nil, err
	}

	path, err := manager.calculation.Execute()
	if err != nil {
		return nil, err
	}
	return manager.calculationTransformer.TransformResult(path, pathRequest, manager.algorithm), nil
}

func (manager *CalculationManager) getCalculationUpdateOptions(streamSession domain.StreamSession) *CalculationUpdateOptions {
	currentPathResult := streamSession.GetPathResult()
	currentAppliedSidList := currentPathResult.GetIpv6SidAddresses()
	manager.log.Debugln("SID list of current path is", currentAppliedSidList)
	pathRequest := streamSession.GetPathRequest()
	intents := pathRequest.GetIntents()
	weightKeys, calculationMode := manager.calculationSetup.GetWeightKeysandCalculationMode(intents)
	return &CalculationUpdateOptions{
		currentPathResult:     currentPathResult,
		currentAppliedSidList: currentAppliedSidList,
		weightKeys:            weightKeys,
		calculationMode:       calculationMode,
		pathRequest:           pathRequest,
	}
}
func (manager *CalculationManager) CalculatePathUpdate(streamSession domain.StreamSession) (domain.PathResult, error) {
	calculationUpdateOptions := manager.getCalculationUpdateOptions(streamSession)
	manager.log.Debugln("Recalculate path with new network state")
	newPathResult, err := manager.CalculateBestPath(calculationUpdateOptions.pathRequest)
	if err != nil {
		return nil, err
	}
	calculationUpdateOptions.newPathResult = newPathResult
	calculationUpdateOptions.streamSession = streamSession
	return manager.calculationUpdater.UpdateCalculation(calculationUpdateOptions)
}
