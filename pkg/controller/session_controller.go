package controller

import (
	"fmt"
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
	quitChan        chan struct{}
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
		quitChan:        make(chan struct{}),
	}
}

func (controller *SessionController) watchForContextCancellation(pathRequest domain.PathRequest, serializedPathRequest string) {
	<-pathRequest.GetContext().Done()
	controller.mu.Lock()
	defer controller.mu.Unlock()
	controller.log.Debugf("Context of path request %s has been cancelled", serializedPathRequest)
	delete(controller.openSessions, serializedPathRequest)
}

func (controller *SessionController) recalculatePathUpdate(session domain.StreamSession) {
	result, err := controller.manager.CalculatePathUpdate(session)
	if err != nil {
		controller.log.Errorln("Failed to recalculate path update: ", err)
		controller.errorChan <- err
	} else if result != nil {
		controller.pathResultChan <- result
	} else {
		controller.log.Debugln("No path update available")
	}
}

func (controller *SessionController) getSessionSnapshot() map[string]domain.StreamSession {
	controller.mu.Lock()
	defer controller.mu.Unlock()
	sessionsSnapshot := make(map[string]domain.StreamSession, len(controller.openSessions))
	for key, session := range controller.openSessions {
		sessionsSnapshot[key] = session
	}
	return sessionsSnapshot
}

func (controller *SessionController) recalculateSessions() {
	if len(controller.openSessions) == 0 {
		controller.log.Debugln("No open sessions to recalculate")
		return
	}

	controller.log.Debugln("Pending updates trigger recalculations of all open sessions")
	wg := sync.WaitGroup{}
	for sessionKey, session := range controller.getSessionSnapshot() {
		controller.log.Debugln("Recalculating for session: ", sessionKey)
		wg.Add(1)
		go func(sessionKey string, session domain.StreamSession) {
			defer wg.Done()
			controller.log.Debugln("Recalculating path update for session: ", sessionKey)
			controller.recalculatePathUpdate(session)
		}(sessionKey, session)
	}
	wg.Wait()
}

func (controller *SessionController) sessionExists(sessionSnapshot map[string]domain.StreamSession, pathRequest domain.PathRequest) bool {
	serializedPathRequest := pathRequest.Serialize()
	_, ok := sessionSnapshot[serializedPathRequest]
	return ok
}

func (controller *SessionController) sendExistingPathResult(sessionSnapshot map[string]domain.StreamSession, pathRequest domain.PathRequest) {
	serializedPathRequest := pathRequest.Serialize()
	controller.log.Debugln("Path request already exists - returning existing path result for path request: ", serializedPathRequest)
	controller.pathResultChan <- sessionSnapshot[serializedPathRequest].GetPathResult()
}

func (controller *SessionController) calculateAndCreateSession(serializedRequest string, pathRequest domain.PathRequest) (domain.PathResult, error) {
	pathResult, err := controller.manager.CalculateBestPath(pathRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate path result: %w", err)
	}

	controller.mu.Lock()
	controller.openSessions[serializedRequest] = domain.NewDomainStreamSession(pathRequest, pathResult)
	controller.mu.Unlock()

	return pathResult, nil
}

func (controller *SessionController) handleError(err error) {
	controller.log.Warnln(err)
	controller.errorChan <- err
}

func (controller *SessionController) handlePathRequest(pathRequest domain.PathRequest) {
	serializedPathRequest := pathRequest.Serialize()
	controller.log.Debugln("Received path request: ", serializedPathRequest)
	sessionSnapshot := controller.getSessionSnapshot()
	if controller.sessionExists(sessionSnapshot, pathRequest) {
		controller.sendExistingPathResult(sessionSnapshot, pathRequest)
		return
	}

	pathResult, err := controller.calculateAndCreateSession(serializedPathRequest, pathRequest)
	if err != nil {
		controller.handleError(err)
		return
	}

	go controller.watchForContextCancellation(pathRequest, serializedPathRequest)
	controller.pathResultChan <- pathResult
}

func (controller *SessionController) Start() {
	controller.log.Infoln("Starting controller")
	for {
		select {
		case <-controller.quitChan:
			return
		case <-controller.updateChan:
			controller.recalculateSessions()
		case pathRequest := <-controller.pathRequestChan:
			controller.handlePathRequest(pathRequest)
		}
	}
}

func (controller *SessionController) Stop() {
	controller.log.Infoln("Stopping session controller")
	close(controller.quitChan)
}
