package controller

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/hawkv6/hawkeye/pkg/api"
	"github.com/hawkv6/hawkeye/pkg/calculation"
	"github.com/hawkv6/hawkeye/pkg/domain"
	"github.com/hawkv6/hawkeye/pkg/graph"
	"github.com/hawkv6/hawkeye/pkg/messaging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewSessionController(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestNewSessionController",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculationManager := calculation.NewMockManager(gomock.NewController(t))
			messagingChannels := messaging.NewPathMessagingChannels()
			sessionController := NewSessionController(calculationManager, messagingChannels, make(chan struct{}))
			assert.NotNil(t, sessionController)
		})
	}
}

func TestSessionController_watchForContextCancellation(t *testing.T) {
	tests := []struct {
		name                   string
		sourceIpv6Address      string
		destinationIpv6Address string
		intents                []domain.Intent
		stream                 api.IntentController_GetIntentPathServer
	}{
		{
			name:                   "TestSessionController_watchForContextCancellation",
			sourceIpv6Address:      "2001:db8::0:1",
			destinationIpv6Address: "2001:db8::0:2",
			intents:                []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})},
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculationManager := calculation.NewMockManager(gomock.NewController(t))
			messagingChannels := messaging.NewPathMessagingChannels()
			sessionController := NewSessionController(calculationManager, messagingChannels, make(chan struct{}))

			ctx, cancel := context.WithCancel(context.Background())
			pathRequest, err := domain.NewDomainPathRequest(tt.sourceIpv6Address, tt.destinationIpv6Address, tt.intents, tt.stream, ctx)
			serializedPathRequest := pathRequest.Serialize()
			assert.NoError(t, err)
			shortestPath := graph.NewMockPath(gomock.NewController(t))
			pathResult, err := domain.NewDomainPathResult(pathRequest, shortestPath, []string{"fc::0:1", "fc::0:2"})
			assert.NoError(t, err)
			session := domain.NewDomainStreamSession(pathRequest, pathResult)
			sessionController.openSessions[serializedPathRequest] = session
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				sessionController.watchForContextCancellation(pathRequest, pathRequest.Serialize())
				wg.Done()
			}()
			cancel()
			wg.Wait()
			assert.Equal(t, 0, len(sessionController.openSessions))
		})
	}
}

func TestSessionController_recalculatePathUpdate(t *testing.T) {
	tests := []struct {
		name                   string
		sourceIpv6Address      string
		destinationIpv6Address string
		intents                []domain.Intent
		stream                 api.IntentController_GetIntentPathServer
		ctx                    context.Context
		wantErr                bool
		want                   domain.PathResult
	}{
		{
			name:                   "TestSessionController_recalculatePathUpdate - no result due to error",
			sourceIpv6Address:      "2001:db8::0:1",
			destinationIpv6Address: "2001:db8::0:2",
			intents:                []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})},
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			wantErr:                true,
			want:                   nil,
		},
		{
			name:                   "TestSessionController_recalculatePathUpdate - no result due to error",
			sourceIpv6Address:      "2001:db8::0:1",
			destinationIpv6Address: "2001:db8::0:2",
			intents:                []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})},
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			wantErr:                false,
			want:                   &domain.DomainPathResult{},
		},
		{
			name:                   "TestSessionController_recalculatePathUpdate - no result due to error",
			sourceIpv6Address:      "2001:db8::0:1",
			destinationIpv6Address: "2001:db8::0:2",
			intents:                []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})},
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
			ctx:                    context.Background(),
			wantErr:                false,
			want:                   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculationManager := calculation.NewMockManager(gomock.NewController(t))
			messagingChannels := messaging.NewPathMessagingChannels()
			sessionController := NewSessionController(calculationManager, messagingChannels, make(chan struct{}))

			if tt.wantErr {
				calculationManager.EXPECT().CalculatePathUpdate(gomock.Any()).Return(nil, assert.AnError).Times(1)
			} else if tt.want != nil {
				calculationManager.EXPECT().CalculatePathUpdate(gomock.Any()).Return(tt.want, nil).Times(1)
			} else {
				calculationManager.EXPECT().CalculatePathUpdate(gomock.Any()).Return(nil, nil).AnyTimes()
			}

			pathRequest, err := domain.NewDomainPathRequest(tt.sourceIpv6Address, tt.destinationIpv6Address, tt.intents, tt.stream, tt.ctx)
			assert.NoError(t, err)
			shortestPath := graph.NewMockPath(gomock.NewController(t))
			pathResult, err := domain.NewDomainPathResult(pathRequest, shortestPath, []string{"fc::0:1", "fc::0:2"})
			assert.NoError(t, err)
			session := domain.NewDomainStreamSession(pathRequest, pathResult)
			go sessionController.recalculatePathUpdate(session)

			if tt.wantErr {
				err = <-messagingChannels.GetErrorChan()
				if (err != nil) != tt.wantErr {
					t.Errorf("SessionController.recalculatePathUpdate() with name '%s' had error = %v, wantErr %v", tt.name, err, tt.wantErr)
				}
			} else if tt.want != nil {
				result := <-messagingChannels.GetPathResponseChan()
				assert.Equal(t, tt.want, result)
			} else {
				select {
				case <-messagingChannels.GetErrorChan():
					t.Errorf("SessionController.recalculatePathUpdate() with name '%s' got unexpected error", tt.name)
				case <-messagingChannels.GetPathResponseChan():
					t.Errorf("SessionController.recalculatePathUpdate() with name '%s' got unexpected result", tt.name)
				default:
				}
			}
		})
	}
}

func TestSessionController_getSessionSnapshot(t *testing.T) {
	tests := []struct {
		name                   string
		sourceIpv6Address      string
		destinationIpv6Address string
		intents                []domain.Intent
		stream                 api.IntentController_GetIntentPathServer
	}{
		{
			name:                   "TestSessionController_getSessionSnapshot",
			sourceIpv6Address:      "2001:db8::0:1",
			destinationIpv6Address: "2001:db8::0:2",
			intents:                []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})},
			stream:                 api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculationManager := calculation.NewMockManager(gomock.NewController(t))
			messagingChannels := messaging.NewPathMessagingChannels()
			sessionController := NewSessionController(calculationManager, messagingChannels, make(chan struct{}))

			pathRequest, err := domain.NewDomainPathRequest(tt.sourceIpv6Address, tt.destinationIpv6Address, tt.intents, tt.stream, context.Background())
			assert.NoError(t, err)
			shortestPath := graph.NewMockPath(gomock.NewController(t))
			pathResult, err := domain.NewDomainPathResult(pathRequest, shortestPath, []string{"fc::0:1", "fc::0:2"})
			assert.NoError(t, err)
			session := domain.NewDomainStreamSession(pathRequest, pathResult)
			sessionController.openSessions[pathRequest.Serialize()] = session
			snapshot := sessionController.getSessionSnapshot()
			assert.Equal(t, 1, len(snapshot))
		})
	}
}

func TestSessionController_recalculateSessions(t *testing.T) {
	sourceIpv6Address := "2001:db8::0:1"
	destinationIpv6Address := "2001:db8::0:2"
	intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
	stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
	calculationManager := calculation.NewMockManager(gomock.NewController(t))
	messagingChannels := messaging.NewPathMessagingChannels()
	ctx := context.Background()

	tests := []struct {
		name       string
		setup      func(*SessionController)
		wantResult bool
	}{
		{
			name: "TestSessionController_recalculateSessions no sessions to recalculate",
		},
		{
			name: "TestSessionController_recalculateSessions with sessions to recalculate",
			setup: func(sessionController *SessionController) {
				shortestPath := graph.NewMockPath(gomock.NewController(t))
				pathRequest, _ := domain.NewDomainPathRequest(sourceIpv6Address, destinationIpv6Address, intents, stream, ctx)
				pathResult, _ := domain.NewDomainPathResult(pathRequest, shortestPath, []string{"fc::0:1", "fc::0:2"})
				calculationManager.EXPECT().CalculatePathUpdate(gomock.Any()).Return(pathResult, nil).AnyTimes()
				session := domain.NewDomainStreamSession(pathRequest, pathResult)
				sessionController.openSessions[pathRequest.Serialize()] = session
			},
			wantResult: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionController := NewSessionController(calculationManager, messagingChannels, make(chan struct{}))
			if tt.setup != nil {
				tt.setup(sessionController)
			}
			go sessionController.recalculateSessions()
			if tt.wantResult {
				<-messagingChannels.GetPathResponseChan()
			}
		})
	}
}

func TestSessionController_sessionExists(t *testing.T) {
	calculationManager := calculation.NewMockManager(gomock.NewController(t))
	messagingChannels := messaging.NewPathMessagingChannels()
	stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
	ctx := context.Background()
	shortestPath := graph.NewMockPath(gomock.NewController(t))
	sourceIpv6Address := "2001:db8::0:1"
	destinationIpv6Address := "2001:db8::0:2"
	intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
	sidAddresses := []string{"fc::0:1", "fc::0:2"}
	tests := []struct {
		name          string
		sessionExists bool
	}{
		{
			name:          "TestSessionController_sessionExists session exists",
			sessionExists: true,
		},
		{
			name:          "TestSessionController_sessionExists session does not exist",
			sessionExists: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionController := NewSessionController(calculationManager, messagingChannels, make(chan struct{}))
			pathRequest, _ := domain.NewDomainPathRequest(sourceIpv6Address, destinationIpv6Address, intents, stream, ctx)
			pathResult, _ := domain.NewDomainPathResult(pathRequest, shortestPath, sidAddresses)
			session := domain.NewDomainStreamSession(pathRequest, pathResult)
			if tt.sessionExists {
				sessionController.openSessions[pathRequest.Serialize()] = session
			}
			snapshot := sessionController.getSessionSnapshot()
			result := sessionController.sessionExists(snapshot, pathRequest)
			assert.Equal(t, tt.sessionExists, result)
		})
	}
}

func TestSessionController_sendExistingPathResult(t *testing.T) {
	calculationManager := calculation.NewMockManager(gomock.NewController(t))
	messagingChannels := messaging.NewPathMessagingChannels()
	stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
	ctx := context.Background()
	shortestPath := graph.NewMockPath(gomock.NewController(t))
	sourceIpv6Address := "2001:db8::0:1"
	destinationIpv6Address := "2001:db8::0:2"
	intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
	sidAddresses := []string{"fc::0:1", "fc::0:2"}
	tests := []struct {
		name          string
		sessionExists bool
	}{
		{
			name:          "TestSessionController_sendExistingPathResult",
			sessionExists: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionController := NewSessionController(calculationManager, messagingChannels, make(chan struct{}))
			pathRequest, _ := domain.NewDomainPathRequest(sourceIpv6Address, destinationIpv6Address, intents, stream, ctx)
			pathResult, _ := domain.NewDomainPathResult(pathRequest, shortestPath, sidAddresses)
			session := domain.NewDomainStreamSession(pathRequest, pathResult)
			sessionController.openSessions[pathRequest.Serialize()] = session
			snapshot := sessionController.getSessionSnapshot()
			go sessionController.sendExistingPathResult(snapshot, pathRequest)
			result := <-messagingChannels.GetPathResponseChan()
			assert.Equal(t, pathResult, result)
		})
	}
}

func TestSessionController_calculateAndCreateSession(t *testing.T) {
	calculationManager := calculation.NewMockManager(gomock.NewController(t))
	messagingChannels := messaging.NewPathMessagingChannels()
	stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
	ctx := context.Background()
	shortestPath := graph.NewMockPath(gomock.NewController(t))
	sourceIpv6Address := "2001:db8::0:1"
	destinationIpv6Address := "2001:db8::0:2"
	intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
	sidAddresses := []string{"fc::0:1", "fc::0:2"}
	tests := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "TestSessionController_calculateAndCreateSession no error",
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionController := NewSessionController(calculationManager, messagingChannels, make(chan struct{}))
			pathRequest, _ := domain.NewDomainPathRequest(sourceIpv6Address, destinationIpv6Address, intents, stream, ctx)
			pathResult, _ := domain.NewDomainPathResult(pathRequest, shortestPath, sidAddresses)
			if tt.wantError {
				calculationManager.EXPECT().CalculateBestPath(gomock.Any()).Return(nil, fmt.Errorf("No path found")).Times(1)
			} else {
				calculationManager.EXPECT().CalculateBestPath(gomock.Any()).Return(pathResult, nil).Times(1)
			}
			result, err := sessionController.calculateAndCreateSession(pathRequest.Serialize(), pathRequest)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, pathResult, result)
			}
		})
	}
}

func TestSessionController_handlePathRequest(t *testing.T) {
	calculationManager := calculation.NewMockManager(gomock.NewController(t))
	messagingChannels := messaging.NewPathMessagingChannels()
	stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
	ctx := context.Background()
	shortestPath := graph.NewMockPath(gomock.NewController(t))
	sourceIpv6Address := "2001:db8::0:1"
	destinationIpv6Address := "2001:db8::0:2"
	intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
	sidAddresses := []string{"fc::0:1", "fc::0:2"}
	tests := []struct {
		name          string
		wantError     bool
		sessionExists bool
	}{
		{
			name:          "TestSessionController_handlePathRequest no error and session does not exist",
			wantError:     false,
			sessionExists: false,
		},
		{
			name:          "TestSessionController_handlePathRequest error and session does not exist",
			wantError:     true,
			sessionExists: false,
		},
		{
			name:          "TestSessionController_handlePathRequest no error and session exists",
			wantError:     false,
			sessionExists: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionController := NewSessionController(calculationManager, messagingChannels, make(chan struct{}))
			pathRequest, _ := domain.NewDomainPathRequest(sourceIpv6Address, destinationIpv6Address, intents, stream, ctx)
			pathResult, _ := domain.NewDomainPathResult(pathRequest, shortestPath, sidAddresses)
			if tt.sessionExists {
				sessionController.openSessions[pathRequest.Serialize()] = domain.NewDomainStreamSession(pathRequest, pathResult)
				go sessionController.handlePathRequest(pathRequest)
				result := <-messagingChannels.GetPathResponseChan()
				assert.Equal(t, pathResult, result)
				return
			}
			if tt.wantError {
				calculationManager.EXPECT().CalculateBestPath(gomock.Any()).Return(nil, fmt.Errorf("No path found")).Times(1)
				go sessionController.handlePathRequest(pathRequest)
				err := <-messagingChannels.GetErrorChan()
				assert.Error(t, err)
			} else {
				calculationManager.EXPECT().CalculateBestPath(gomock.Any()).Return(pathResult, nil).Times(1)
				go sessionController.handlePathRequest(pathRequest)
				result := <-messagingChannels.GetPathResponseChan()
				assert.Equal(t, pathResult, result)
			}
		})
	}
}

func TestSessionController_Start(t *testing.T) {
	calculationManager := calculation.NewMockManager(gomock.NewController(t))
	messagingChannels := messaging.NewPathMessagingChannels()
	stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
	ctx := context.Background()
	sourceIpv6Address := "2001:db8::0:1"
	destinationIpv6Address := "2001:db8::0:2"
	intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
	tests := []struct {
		name string
	}{
		{
			name: "TestSessionController_Start",
		},
	}
	for _, tt := range tests {
		sessionController := NewSessionController(calculationManager, messagingChannels, make(chan struct{}))
		calculationManager.EXPECT().CalculateBestPath(gomock.Any()).Return(nil, fmt.Errorf("No path found")).Times(1)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			sessionController.Start()
			wg.Done()
		}()
		pathRequest, _ := domain.NewDomainPathRequest(sourceIpv6Address, destinationIpv6Address, intents, stream, ctx)
		sessionController.pathRequestChan <- pathRequest
		<-sessionController.errorChan
		sessionController.updateChan <- struct{}{}
		sessionController.quitChan <- struct{}{}
		wg.Wait()
		t.Logf("SessionController.Start() with name '%s' completed", tt.name)
	}
}

func TestSessionController_Stop(t *testing.T) {
	calculationManager := calculation.NewMockManager(gomock.NewController(t))
	messagingChannels := messaging.NewPathMessagingChannels()
	stream := api.NewMockIntentController_GetIntentPathServer(gomock.NewController(t))
	ctx := context.Background()
	sourceIpv6Address := "2001:db8::0:1"
	destinationIpv6Address := "2001:db8::0:2"
	intents := []domain.Intent{domain.NewDomainIntent(domain.IntentTypeLowLatency, []domain.Value{})}
	tests := []struct {
		name string
	}{
		{
			name: "TestSessionController_Stop",
		},
	}
	for _, tt := range tests {
		sessionController := NewSessionController(calculationManager, messagingChannels, make(chan struct{}))
		calculationManager.EXPECT().CalculateBestPath(gomock.Any()).Return(nil, fmt.Errorf("No path found")).Times(1)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			sessionController.Start()
			wg.Done()
		}()
		pathRequest, _ := domain.NewDomainPathRequest(sourceIpv6Address, destinationIpv6Address, intents, stream, ctx)
		sessionController.pathRequestChan <- pathRequest
		<-sessionController.errorChan
		sessionController.updateChan <- struct{}{}
		sessionController.Stop()
		wg.Wait()
		t.Logf("SessionController.Stop() with name '%s' completed", tt.name)
	}
}
