package controller

import (
	"sync"

	"github.com/hawkv6/hawkeye/pkg/calculation"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/hawkv6/hawkeye/pkg/messaging"
	"github.com/sirupsen/logrus"
)

type SessionController struct {
	log             *logrus.Entry
	manager         calculation.Manager
	openSessions    map[string]domain.StreamSession
	pathRequestChan chan domain.PathRequest
	pathResultChan  chan domain.PathResult
	mu              sync.Mutex
	updateChan      chan struct{}
}

func NewSessionController(manager calculation.Manager, messagingChannels messaging.MessagingChannels, updateChan chan struct{}) *SessionController {
	return &SessionController{
		log:             logging.DefaultLogger.WithField("subsystem", Subsystem),
		manager:         manager,
		openSessions:    make(map[string]domain.StreamSession, 0),
		pathRequestChan: messagingChannels.GetPathRequestChan(),
		pathResultChan:  messagingChannels.GetPathResponseChan(),
		mu:              sync.Mutex{},
		updateChan:      updateChan,
	}
}

func (controller *SessionController) watchForContextCancellation(pathRequest domain.PathRequest, serializedPathRequest string) {
	<-pathRequest.GetContext().Done()
	controller.mu.Lock()
	controller.log.Debugf("Context of path request %s has been cancelled", pathRequest)
	delete(controller.openSessions, serializedPathRequest)
	controller.mu.Unlock()
}

func (controller *SessionController) recalculateSessions() {
	if len(controller.openSessions) == 0 {
		controller.log.Debugln("No open sessions to recalculate")
		return
	}
	controller.log.Debugln("Pending updates trigger recalculations of all open sessions")
	for sessionKey, session := range controller.openSessions {
		controller.log.Debugln("Recalculating for session: ", sessionKey)
		result := controller.manager.CalculatePathUpdate(session)
		if result != nil {
			controller.pathResultChan <- *result
		}
	}
}

func (controller *SessionController) handlePathRequest(pathRequest domain.PathRequest) {
	serializedPathRequest := pathRequest.Serialize()
	controller.log.Debugln("Received path request: ", serializedPathRequest)
	if _, ok := controller.openSessions[serializedPathRequest]; ok {
		controller.log.Debugln("Path request already exists")
		controller.pathResultChan <- controller.openSessions[serializedPathRequest].GetPathResult()
	} else {
		pathResult := controller.manager.CalculateBestPath(pathRequest)
		streamSession, err := domain.NewDefaultStreamSession(pathRequest, pathResult)
		if err != nil {
			controller.log.Errorln("Failed to create stream session: ", err)
			return
		}
		controller.openSessions[serializedPathRequest] = streamSession
		go controller.watchForContextCancellation(pathRequest, serializedPathRequest)
		controller.pathResultChan <- pathResult
	}
}

func (controller *SessionController) Start() {
	controller.log.Infoln("Starting controller")
	for {
		select {
		case <-controller.updateChan:
			controller.recalculateSessions()
		case pathRequest := <-controller.pathRequestChan:
			controller.handlePathRequest(pathRequest)
		}
	}
}