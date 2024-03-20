package controller

import (
	"reflect"
	"sync"

	"github.com/hawkv6/hawkeye/pkg/calculation"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/hawkv6/hawkeye/pkg/messaging"
	"github.com/sirupsen/logrus"
)

type DefaultController struct {
	log             *logrus.Entry
	calculator      calculation.Calculator
	openSessions    map[string]domain.StreamSession
	pathRequestChan chan domain.PathRequest
	pathResultChan  chan domain.PathResult
	mu              sync.Mutex
	updateChan      chan struct{}
}

func NewDefaultController(calculator calculation.Calculator, messagingChannels messaging.MessagingChannels, updateChan chan struct{}) *DefaultController {
	return &DefaultController{
		log:             logging.DefaultLogger.WithField("subsystem", Subsystem),
		calculator:      calculator,
		openSessions:    make(map[string]domain.StreamSession, 0),
		pathRequestChan: messagingChannels.GetPathRequestChan(),
		pathResultChan:  messagingChannels.GetPathResponseChan(),
		mu:              sync.Mutex{},
		updateChan:      updateChan,
	}
}

func (controller *DefaultController) watchForContextCancellation(pathRequest domain.PathRequest, serializedPathRequest string) {
	<-pathRequest.GetContext().Done()
	controller.mu.Lock()
	controller.log.Debugf("Context of path request %s has been cancelled", pathRequest)
	delete(controller.openSessions, serializedPathRequest)
	controller.mu.Unlock()
}

func (controller *DefaultController) recalculateSessions() {
	if len(controller.openSessions) == 0 {
		controller.log.Debugln("No open sessions to recalculate")
		return
	}
	controller.log.Debugln("Pending updates trigger recalculations of all open sessions")
	for sessionKey, session := range controller.openSessions {
		controller.log.Debugln("Recalculating session: ", sessionKey)
		result := controller.calculator.HandlePathRequest(session.GetPathRequest())
		newSidList := result.GetIpv6SidAddresses()
		oldSidList := session.GetPathResult().GetIpv6SidAddresses()
		if !reflect.DeepEqual(newSidList, oldSidList) {
			controller.openSessions[sessionKey].SetPathResult(result)
			controller.log.Debugln("Path result has changed for session: ", sessionKey)
			controller.log.Debugf("Sid List changed from: %s to %s", oldSidList, newSidList)
			controller.pathResultChan <- result
		} else {
			controller.log.Debugf("No change in path result for session %s", sessionKey)
		}
	}
}

func (controller *DefaultController) handlePathRequest(pathRequest domain.PathRequest) {
	serializedPathRequest := pathRequest.Serialize()
	controller.log.Debugln("Received path request: ", serializedPathRequest)
	if _, ok := controller.openSessions[serializedPathRequest]; ok {
		controller.log.Debugln("Path request already exists")
		controller.pathResultChan <- controller.openSessions[serializedPathRequest].GetPathResult()
	} else {
		pathResult := controller.calculator.HandlePathRequest(pathRequest)
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

func (controller *DefaultController) Start() {
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
