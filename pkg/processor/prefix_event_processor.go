package processor

import (
	"fmt"

	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/sirupsen/logrus"
)

type PrefixEventProcessor struct {
	log          *logrus.Entry
	graph        graph.Graph
	cache        cache.Cache
	prefixCounts map[string]int
}

func NewPrefixEventProcessor(graph graph.Graph, cache cache.Cache) *PrefixEventProcessor {
	return &PrefixEventProcessor{
		log:          logrus.WithField("subsystem", Subsystem),
		graph:        graph,
		cache:        cache,
		prefixCounts: make(map[string]int),
	}
}

func (processor *PrefixEventProcessor) clearDuplicateAnnouncedPrefix(prefix domain.Prefix, networkAddress string, subnetLength uint8) {
	clientNetwork := processor.cache.GetClientNetworkByKey(prefix.GetKey())
	if clientNetwork == nil {
		processor.log.Debugf("Delete network %s/%d from cache since it's announced several times and thus not a client network: ", networkAddress, subnetLength)
		processor.prefixCounts[networkAddress]++
		processor.cache.RemoveClientNetwork(prefix)
	}
}

func (processor *PrefixEventProcessor) addNetworkToCache(prefix domain.Prefix, networkAddress string, subnetLength uint8) {
	processor.log.Debugf("Add network %s/%d to cache ", networkAddress, subnetLength)
	processor.prefixCounts[networkAddress] = 1
	processor.cache.StoreClientNetwork(prefix)
}

func (processor *PrefixEventProcessor) processPrefix(prefix domain.Prefix) {
	networkAddress := prefix.GetPrefix()
	subnetLength := prefix.GetPrefixLength()
	_, ok := processor.prefixCounts[networkAddress]
	if !ok {
		processor.addNetworkToCache(prefix, networkAddress, subnetLength)
	} else {
		processor.clearDuplicateAnnouncedPrefix(prefix, networkAddress, subnetLength)
	}
}

func (processor *PrefixEventProcessor) deleteClientNetwork(key string) error {
	prefix := processor.cache.GetClientNetworkByKey(key)
	if prefix == nil {
		return fmt.Errorf("Network with key %s does not exist in cache", key)
	}
	networkAddress := prefix.GetPrefix()
	subnetLength := prefix.GetPrefixLength()
	if _, ok := processor.prefixCounts[networkAddress]; !ok {
		return fmt.Errorf("Network %s/%d does not exist in prefix counts", networkAddress, subnetLength)
	}
	if processor.prefixCounts[networkAddress] > 1 {
		processor.log.Debugf("Decrement network %s/%d from announced prefix count", networkAddress, subnetLength)
		processor.prefixCounts[networkAddress]--
		return nil
	}
	processor.log.Debugf("Delete client network %s/%d from cache", networkAddress, subnetLength)
	delete(processor.prefixCounts, networkAddress)
	processor.cache.RemoveClientNetwork(prefix)
	return nil
}

func (processor *PrefixEventProcessor) ProcessPrefixes(prefixes []domain.Prefix) {
	for _, prefix := range prefixes {
		processor.processPrefix(prefix)
	}
}

func (processor *PrefixEventProcessor) HandleEvent(event domain.NetworkEvent) bool {
	switch eventType := event.(type) {
	case *domain.AddPrefixEvent:
		processor.log.Debugln("Received AddPrefixEvent: ", eventType.Prefix.GetKey())
		processor.processPrefix(eventType.Prefix)
	case *domain.DeletePrefixEvent:
		processor.log.Debugln("Received DeletePrefixEvent: ", eventType.GetKey())
		if err := processor.deleteClientNetwork(eventType.GetKey()); err != nil {
			processor.log.Warnln("Error deleting client network from cache: ", err)
		}
	}
	return false
}
