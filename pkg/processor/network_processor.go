package processor

import (
	"time"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type NetworkProcessor struct {
	log                 *logrus.Entry
	graph               graph.Graph
	cache               cache.Cache
	eventChan           chan domain.NetworkEvent
	quitChan            chan struct{}
	updateChan          chan struct{}
	needsSubgraphUpdate bool
	eventDispatcher     *EventDispatcher
	nodeProcessor       NodeProcessor
	linkProcessor       LinkProcessor
	prefixProcessor     PrefixProcessor
	sidProcessor        SidProcessor
}

type EventOptions struct {
	NodeEventProcessor   NodeProcessor
	LinkEventProcessor   LinkProcessor
	PrefixEventProcessor PrefixProcessor
	SidEventProcessor    SidProcessor
	EventDispatcher      *EventDispatcher
}

func NewNetworkProcessor(graph graph.Graph, cache cache.Cache, eventChan chan domain.NetworkEvent, updateChan chan struct{}, eventOptions EventOptions) *NetworkProcessor {

	return &NetworkProcessor{
		log:                 logging.DefaultLogger.WithField("subsystem", Subsystem),
		graph:               graph,
		cache:               cache,
		eventChan:           eventChan,
		updateChan:          updateChan,
		quitChan:            make(chan struct{}),
		needsSubgraphUpdate: false,
		nodeProcessor:       eventOptions.NodeEventProcessor,
		linkProcessor:       eventOptions.LinkEventProcessor,
		prefixProcessor:     eventOptions.PrefixEventProcessor,
		sidProcessor:        eventOptions.SidEventProcessor,
		eventDispatcher:     eventOptions.EventDispatcher,
	}
}

func (processor *NetworkProcessor) ProcessNodes(nodes []domain.Node) {
	processor.nodeProcessor.ProcessNodes(nodes)
}

func (processor *NetworkProcessor) ProcessLinks(links []domain.Link) error {
	return processor.linkProcessor.ProcessLinks(links)
}

func (processor *NetworkProcessor) ProcessPrefixes(prefixes []domain.Prefix) {
	processor.prefixProcessor.ProcessPrefixes(prefixes)
}

func (processor *NetworkProcessor) ProcessSids(sids []domain.Sid) {
	processor.sidProcessor.ProcessSids(sids)
}

func (processor *NetworkProcessor) Start() {
	holdTime := helper.NetworkProcessorHoldTime
	processor.log.Infof("Starting processing network updates with hold time %s", holdTime.String())

	timer := time.NewTimer(holdTime)
	defer timer.Stop()
	mutexesLocked := false

	for {
		select {
		case event := <-processor.eventChan:
			if !mutexesLocked {
				processor.log.Debugln("Locking cache and graph mutexes")
				processor.cache.Lock()
				processor.graph.Lock()
				mutexesLocked = true
			}
			if processor.eventDispatcher.Dispatch(event) {
				processor.needsSubgraphUpdate = true
			}
			timer.Reset(holdTime)
		case <-timer.C:
			if mutexesLocked {
				processor.log.Debugln("Unlocking cache and graph mutexes")
				processor.cache.Unlock()
				processor.graph.Unlock()
				mutexesLocked = false
			}
			if processor.needsSubgraphUpdate {
				processor.graph.UpdateSubGraphs()
				processor.needsSubgraphUpdate = false
			}
			processor.updateChan <- struct{}{}
		case <-processor.quitChan:
			if mutexesLocked {
				processor.cache.Unlock()
				processor.graph.Unlock()
			}
			return
		}
	}
}

func (processor *NetworkProcessor) Stop() {
	processor.log.Infoln("Stopping network processor")
	close(processor.quitChan)
}
