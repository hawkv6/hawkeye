package processor

import (
	"reflect"

	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type Dispatcher interface {
	Dispatch(event domain.NetworkEvent) bool
}

type EventDispatcher struct {
	log           *logrus.Entry
	eventHandlers map[reflect.Type]EventHandler
}

func NewEventDispatcher(nodeEventHandler *NodeEventProcessor, linkEventHandler *LinkEventProcessor, prefixEventHandler *PrefixEventProcessor, sidEventHandler *SidEventProcessor) *EventDispatcher {
	return &EventDispatcher{
		log: logging.DefaultLogger.WithField("subsystem", Subsystem),
		eventHandlers: map[reflect.Type]EventHandler{
			reflect.TypeOf(&domain.AddNodeEvent{}):      nodeEventHandler,
			reflect.TypeOf(&domain.UpdateNodeEvent{}):   nodeEventHandler,
			reflect.TypeOf(&domain.DeleteNodeEvent{}):   nodeEventHandler,
			reflect.TypeOf(&domain.AddLinkEvent{}):      linkEventHandler,
			reflect.TypeOf(&domain.UpdateLinkEvent{}):   linkEventHandler,
			reflect.TypeOf(&domain.DeleteLinkEvent{}):   linkEventHandler,
			reflect.TypeOf(&domain.AddPrefixEvent{}):    prefixEventHandler,
			reflect.TypeOf(&domain.DeletePrefixEvent{}): prefixEventHandler,
			reflect.TypeOf(&domain.AddSidEvent{}):       sidEventHandler,
			reflect.TypeOf(&domain.DeleteSidEvent{}):    sidEventHandler,
		},
	}
}

func (dispatcher *EventDispatcher) Dispatch(event domain.NetworkEvent) bool {
	if handler, ok := dispatcher.eventHandlers[reflect.TypeOf(event)]; ok {
		return handler.HandleEvent(event)
	} else {
		dispatcher.log.Warnf("No handler found for event: %v", event)
	}
	return false
}
