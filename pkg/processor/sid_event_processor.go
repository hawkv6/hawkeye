package processor

import (
	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type SidEventProcessor struct {
	log   *logrus.Entry
	graph graph.Graph
	cache cache.Cache
}

func NewSidEventProcessor(graph graph.Graph, cache cache.Cache) *SidEventProcessor {
	return &SidEventProcessor{
		log:   logging.DefaultLogger.WithField("subsystem", Subsystem),
		graph: graph,
		cache: cache,
	}
}

func (processor *SidEventProcessor) addSidtoCache(sid domain.Sid) {
	processor.log.Debugf("Add SRv6 SID %s to cache", sid.GetSid())
	processor.cache.StoreSid(sid)
}

func (processor *SidEventProcessor) deleteSidFromCache(key string) {
	processor.log.Debugf("Delete SRv6 SID %s from cache", key)
	sid := processor.cache.GetSidByKey(key)
	if sid == nil {
		processor.log.Debugf("SID with key %s does not exist in cache", key)
	}
	processor.cache.RemoveSid(sid)
}

func (processor *SidEventProcessor) ProcessSids(sids []domain.Sid) {
	for _, sid := range sids {
		processor.addSidtoCache(sid)
	}
}

func (processor *SidEventProcessor) HandleEvent(event domain.NetworkEvent) bool {
	switch eventType := event.(type) {
	case *domain.AddSidEvent:
		processor.log.Debugln("Received AddSidEvent: ", eventType.GetSid())
		processor.addSidtoCache(eventType.Sid)
	case *domain.DeleteSidEvent:
		processor.log.Debugln("Received DeleteSidEvent: ", eventType.GetKey())
		processor.deleteSidFromCache(eventType.GetKey())
	default:
		processor.log.Warnf("No handler found for event: %v", event)
	}
	return false
}
