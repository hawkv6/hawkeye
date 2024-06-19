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
	errorChan       chan error
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
		errorChan:       messagingChannels.GetErrorChan(),
		mu:              sync.Mutex{},
		updateChan:      updateChan,
	}
}

func (controller *SessionController) watchForContextCancellation(pathRequest domain.PathRequest, serializedPathRequest string) {
	<-pathRequest.GetContext().Done()
	controller.mu.Lock()
	defer controller.mu.Unlock()
	controller.log.Debugf("Context of path request %s has been cancelled", pathRequest.Serialize())
	delete(controller.openSessions, serializedPathRequest)
}

func (controller *SessionController) recalculateSessions() {
	if len(controller.openSessions) == 0 {
		controller.log.Debugln("No open sessions to recalculate")
		return
	}
	controller.mu.Lock()
	defer controller.mu.Unlock()
	controller.log.Debugln("Pending updates trigger recalculations of all open sessions")
	for sessionKey, session := range controller.openSessions {
		controller.log.Debugln("Recalculating for session: ", sessionKey)
		result, err := controller.manager.CalculatePathUpdate(session)
		if err != nil {
			controller.log.Errorln("Failed to recalculate path update: ", err)
			controller.errorChan <- err
		} else if result != nil {
			controller.pathResultChan <- *result
		} else {
			controller.log.Debugln("No path update available")
		}
	}
}

func (controller *SessionController) handlePathRequest(pathRequest domain.PathRequest) {
	serializedPathRequest := pathRequest.Serialize()
	controller.log.Debugln("Received path request: ", serializedPathRequest)
	if _, ok := controller.openSessions[serializedPathRequest]; ok {
		controller.log.Debugln("Path request already exists - returning existing path result")
		controller.pathResultChan <- controller.openSessions[serializedPathRequest].GetPathResult()
	} else {
		pathResult, err := controller.manager.CalculateBestPath(pathRequest)
		if err != nil {
			controller.log.Warnln("Failed to calculate path result: ", err)
			controller.errorChan <- err
			return
		}
		streamSession, err := domain.NewDefaultStreamSession(pathRequest, pathResult)
		if err != nil {
			controller.log.Warnln("Failed to create stream session: ", err)
			controller.errorChan <- err
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
