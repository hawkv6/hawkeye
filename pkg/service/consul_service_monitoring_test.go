package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewConsulServiceMonitor(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{
			name:    "Test NewConsulServiceMonitor with empty address",
			address: "",
			wantErr: true,
		},
		{
			name:    "Test NewConsulServiceMonitor with valid ipv4 address",
			address: "127.0.0.1",
			wantErr: false,
		},
		{
			name:    "Test NewConsulServiceMonitor with valid hostname",
			address: "localhost",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConsulServiceMonitor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, serviceMonitor)
			}
		})
	}
}

func TestConsulServiceMonitoring_storeServiceSidInCache(t *testing.T) {
	serviceType := "fw"
	prefixSid := "fc:0:2f::"
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitoring_storeServiceSidInCache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			cacheMock.EXPECT().StoreServiceSid(serviceType, prefixSid).Return()
			serviceMonitor.storeServiceSidInCache(serviceType, prefixSid)
		})
	}
}

func TestConsulServiceMonitor_createServiceEntry(t *testing.T) {
	serviceId := "SERA-1"
	serviceType := "fw"
	prefixSid := "fc:0:2f::"
	healthy := true
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitor_createServiceEntry healthy service does not exist",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			cacheMock.EXPECT().StoreServiceSid(serviceType, prefixSid).Return().AnyTimes()
			serviceMonitor.createServiceEntry(serviceId, serviceType, prefixSid, healthy)
		})
	}
}

func TestConsulServiceMonitor_removeServiceSidInCache(t *testing.T) {
	serviceType := "fw"
	prefixSid := "fc:0:2f::"
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitor_removeServiceSidInCache",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			cacheMock.EXPECT().RemoveServiceSid(serviceType, prefixSid).Return().AnyTimes()
			serviceMonitor.removeServiceSidInCache(serviceType, prefixSid)
		})
	}
}

func TestConsulServiceMOnitor_deleteServiceEntry(t *testing.T) {
	serviceId := "SERA-1"
	serviceType := "fw"
	serviceSid := "fc:0:2f::"
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMOnitor_deleteServiceEntry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			cacheMock.EXPECT().RemoveServiceSid(gomock.Any(), gomock.Any()).Return().AnyTimes()
			serviceMonitor.services = make(map[string]map[string]*ConcreteService)
			serviceMonitor.services[serviceType] = make(map[string]*ConcreteService)
			serviceMonitor.services[serviceType][serviceId] = NewConcreteService(serviceType, serviceId, serviceSid, true)
			serviceMonitor.deleteServiceEntry(serviceId, serviceType)
			assert.Empty(t, serviceMonitor.services[serviceType])
		})
	}
}

func TestConsulServiceMonitor_deleteServiceEntries(t *testing.T) {
	serviceType := "fw"
	knownServices := map[string]bool{
		"SERA-1": true,
		"SERA-2": false,
	}
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitor_deleteServiceEntries",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			cacheMock.EXPECT().RemoveServiceSid(gomock.Any(), gomock.Any()).Return().AnyTimes()
			serviceMonitor.services = make(map[string]map[string]*ConcreteService)
			serviceMonitor.services["fw"] = make(map[string]*ConcreteService)
			serviceMonitor.services["fw"]["SERA-1"] = NewConcreteService(serviceType, "SERA-1", "fc:0:2f::", true)
			serviceMonitor.services["fw"]["SERA-2"] = NewConcreteService(serviceType, "SERA-2", "fc:0:2f::", true)
			serviceMonitor.deleteServiceEntries(knownServices, serviceType)
			assert.Len(t, serviceMonitor.services[serviceType], 1)
		})
	}
}

func TestConsulServiceMonitor_getKnownServices(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestConsulServiceMonitor_getKnownServices no error",
			wantErr: false,
		},
		{
			name:    "TestConsulServiceMonitor_getKnownServices error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			if tt.wantErr {
				_, _, err := serviceMonitor.getKnownServices(nil)
				assert.Error(t, err)
			} else {
				serviceMonitor.services = make(map[string]map[string]*ConcreteService)
				serviceMonitor.services["fw"] = make(map[string]*ConcreteService)
				serviceMonitor.services["fw"]["SERA-1"] = NewConcreteService("fw", "SERA-1", "fc:0:2f::", true)
				knownServices, serviceType, err := serviceMonitor.getKnownServices([]*api.ServiceEntry{
					{
						Service: &api.AgentService{
							Service: "fw",
						},
					},
				})
				assert.NoError(t, err)
				assert.Equal(t, "fw", serviceType)
				assert.Len(t, knownServices, 1)
			}
		})
	}
}

func TestConsulServiceMonitor_processServiceEntries(t *testing.T) {
	serviceName := "SERA-1"
	serviceType := "fw"
	prefixSid := "fc:0:2f::"
	tests := []struct {
		name         string
		prefixSid    string
		knownService bool
	}{
		{
			name:         "TestConsulServiceMonitor_processServiceEntries no prefixSid",
			prefixSid:    "",
			knownService: false,
		},
		{
			name:         "TestConsulServiceMonitor_processServiceEntries with prefixSid and known service",
			prefixSid:    prefixSid,
			knownService: true,
		},
		{
			name:         "TestConsulServiceMonitor_processServiceEntries with prefixSid and known service",
			prefixSid:    prefixSid,
			knownService: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			cacheMock.EXPECT().StoreServiceSid(gomock.Any(), gomock.Any()).Return().AnyTimes()
			var knownServices map[string]bool
			if tt.knownService {
				knownServices = map[string]bool{
					serviceName: true,
				}
			} else {
				knownServices = make(map[string]bool)
			}
			serviceEntries := []*api.ServiceEntry{
				{
					Service: &api.AgentService{
						ID:      serviceName,
						Service: "fw",
						Meta: map[string]string{
							"sid": tt.prefixSid,
						},
					},
					Checks: api.HealthChecks{
						{
							Status: api.HealthPassing,
						},
					},
				},
			}
			serviceMonitor.processServiceEntries(serviceEntries, knownServices, serviceType)
			if tt.prefixSid == "" || tt.knownService {
				assert.Len(t, serviceMonitor.services["fw"], 0)
			} else {
				assert.Len(t, serviceMonitor.services["fw"], 1)
			}
		})
	}
}

func TestConsulServiceMonitor_validateServiceEntries(t *testing.T) {
	tests := []struct {
		name           string
		serviceEntries bool
	}{
		{
			name:           "TestConsulServiceMonitor_validateServiceEntries no service entries",
			serviceEntries: false,
		},
		{
			name:           "TestConsulServiceMonitor_validateServiceEntries with service entries",
			serviceEntries: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			if tt.serviceEntries {
				serviceEntries := []*api.ServiceEntry{
					{
						Service: &api.AgentService{
							Service: "fw",
						},
					},
				}
				knownServices, serviceType, err := serviceMonitor.getKnownServices(serviceEntries)
				assert.NoError(t, err)
				serviceMonitor.validateServiceEntries(serviceEntries)
				assert.Equal(t, "fw", serviceType)
				assert.Len(t, knownServices, 0)
			} else {
				serviceMonitor.validateServiceEntries(nil)
			}
		})
	}
}

func TestConsulServiceMonitor_markServiceUnhealty(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitor_markServiceUnhealty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			cacheMock.EXPECT().RemoveServiceSid(gomock.Any(), gomock.Any()).Return().AnyTimes()
			serviceMonitor.services = make(map[string]map[string]*ConcreteService)
			serviceMonitor.services["fw"] = make(map[string]*ConcreteService)
			service := NewConcreteService("fw", "SERA-1", "fc:0:2f::", true)
			serviceMonitor.services["fw"]["SERA-1"] = service
			serviceMonitor.markServiceUnhealthy(service)
			assert.False(t, serviceMonitor.services["fw"]["SERA-1"].healthy)
			assert.False(t, service.healthy)

		})
	}
}

func TestConsulServiceMonitor_markServiceHealthy(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitor_markServiceHealthy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			cacheMock.EXPECT().StoreServiceSid(gomock.Any(), gomock.Any()).Return().AnyTimes()
			serviceMonitor.services = make(map[string]map[string]*ConcreteService)
			serviceMonitor.services["fw"] = make(map[string]*ConcreteService)
			service := NewConcreteService("fw", "SERA-1", "fc:0:2f::", false)
			serviceMonitor.services["fw"]["SERA-1"] = service
			serviceMonitor.markServiceHealthy(service)
			assert.True(t, serviceMonitor.services["fw"]["SERA-1"].healthy)
			assert.True(t, service.healthy)
		})
	}
}

func TestConsulServiceMonitor_checkServiceHealth(t *testing.T) {
	tests := []struct {
		name                string
		healthy             bool
		healthStatusChanged bool
	}{
		{
			name:                "TestConsulServiceMonitor_checkServiceHealth healthy and health status changed",
			healthy:             true,
			healthStatusChanged: true,
		},
		{
			name:                "TestConsulServiceMonitor_checkServiceHealth unhealthy and health status changed",
			healthy:             false,
			healthStatusChanged: true,
		},
		{
			name:                "TestConsulServiceMonitor_checkServiceHealth healthy and health status changed",
			healthy:             true,
			healthStatusChanged: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			cacheMock.EXPECT().RemoveServiceSid(gomock.Any(), gomock.Any()).Return().AnyTimes()
			cacheMock.EXPECT().StoreServiceSid(gomock.Any(), gomock.Any()).Return().AnyTimes()
			serviceMonitor.services = make(map[string]map[string]*ConcreteService)
			serviceMonitor.services["fw"] = make(map[string]*ConcreteService)
			var service *ConcreteService
			if tt.healthStatusChanged {
				service = NewConcreteService("fw", "SERA-1", "fc:0:2f::", !tt.healthy)
			} else {
				service = NewConcreteService("fw", "SERA-1", "fc:0:2f::", tt.healthy)
			}
			serviceMonitor.services["fw"]["SERA-1"] = service
			var status string
			if tt.healthy {
				status = api.HealthPassing
			} else {
				status = api.HealthCritical
			}
			serviceEntries := []*api.ServiceEntry{
				{
					Service: &api.AgentService{
						ID:      "SERA-1",
						Service: "fw",
					},
					Checks: api.HealthChecks{
						{
							Status: status,
						},
					},
				},
			}
			serviceMonitor.checkServiceHealth(serviceEntries)
			assert.Equal(t, tt.healthy, serviceMonitor.services["fw"]["SERA-1"].healthy)
			assert.Equal(t, tt.healthy, service.healthy)
		})
	}
}

func TestConsulServiceMonitor_sendUpdate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitor_sendUpdate",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			serviceMonitor.needsUpdate = true
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				<-serviceMonitor.updateChan
				wg.Done()
			}()
			serviceMonitor.sendUpdate()
			wg.Wait()
		})
	}
}

func TestConsulServiceMonitor_queryAndUpdateServiceHealth(t *testing.T) {
	serviceType := "fw"
	lastIndex := uint64(0)
	tests := []struct {
		name        string
		throwsError bool
	}{
		{
			name:        "TestConsulServiceMonitor_queryAndUpdateServiceHealth throws error",
			throwsError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			cacheMock.EXPECT().Lock().Return().AnyTimes()
			cacheMock.EXPECT().Unlock().Return().AnyTimes()
			cacheMock.EXPECT().RemoveServiceSid(gomock.Any(), gomock.Any()).Return().AnyTimes()
			cacheMock.EXPECT().StoreServiceSid(gomock.Any(), gomock.Any()).Return().AnyTimes()
			serviceMonitor.services = make(map[string]map[string]*ConcreteService)
			serviceMonitor.services["fw"] = make(map[string]*ConcreteService)
			serviceMonitor.services["fw"]["SERA-1"] = NewConcreteService("fw", "SERA-1", "fc:0:2f::", true)
			assert.True(t, serviceMonitor.queryAndUpdateServiceHealth(serviceType, &lastIndex))
		})
	}
}

func TestConsulServiceMonitor_monitorServiceHealth(t *testing.T) {
	serviceType := "fw"
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitor_monitorServiceHealth",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			stopChan := make(chan struct{})
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, stopChan, "localhost")
			assert.Nil(t, err)
			ctx, cancel := context.WithCancel(context.Background())
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				serviceMonitor.monitorServiceHealth(ctx, serviceType)
				wg.Done()
			}()
			cancel()
			wg.Wait()
		})
	}
}

func TestConsulServiceMonitor_unmonitoredServicesExist(t *testing.T) {
	serviceType := "fw"
	tests := []struct {
		name                     string
		unmonitoredServicesExist bool
	}{
		{
			name:                     "TestConsulServiceMonitor_unmonitoredServicesExist unmonitored services exist",
			unmonitoredServicesExist: true,
		},
		{
			name:                     "TestConsulServiceMonitor_unmonitoredServicesExist unmonitored services do not exist",
			unmonitoredServicesExist: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			services := make(map[string][]string)
			services["consul"] = []string{}
			if tt.unmonitoredServicesExist {
				serviceMonitor.monitoredServices = make(map[string]context.CancelFunc)
				services[serviceType] = []string{}
			}
			assert.Equal(t, tt.unmonitoredServicesExist, serviceMonitor.unmonitoredServicesExist(services))
		})
	}
}

func TestConsulServiceMonitor_stopMonitoringRemovedServices(t *testing.T) {
	serviceType := "fw"
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitor_stopMonitoringRemovedServices",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			serviceMonitor.monitoredServices = make(map[string]context.CancelFunc)
			serviceMonitor.monitoredServices[serviceType] = func() {}
			newServiceTypes := make(map[string][]string)
			serviceMonitor.stopMonitoringRemovedServices(newServiceTypes)
			assert.Empty(t, serviceMonitor.monitoredServices)
		})
	}
}

func TestConsulServiceMonitor_startServiceHealthMonitoring(t *testing.T) {
	const FW = "fw"
	const IDS = "ids"
	tests := []struct {
		name         string
		serviceTypes []string
	}{
		{
			name:         "TestConsulServiceMonitor_startServiceHealthMonitoring multiple services",
			serviceTypes: []string{FW, IDS},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			serviceMonitor.monitoredServices = make(map[string]context.CancelFunc)
			services := make(map[string][]string)
			for _, serviceType := range tt.serviceTypes {
				services[serviceType] = []string{}
			}
			serviceMonitor.startServiceHealthMonitoring(services)
			serviceMonitor.mu.RLock()
			for _, cancel := range serviceMonitor.monitoredServices {
				cancel()
			}
			serviceMonitor.mu.RUnlock()
			serviceMonitor.wg.Wait()
		})
	}
}

func TestConsulServiceMonitor_updateMonitoredServices(t *testing.T) {
	const FW = "fw"
	const IDS = "ids"
	tests := []struct {
		name         string
		serviceTypes []string
	}{
		{
			name:         "TestConsulServiceMonitor_updasteMonitoredServices only fw service",
			serviceTypes: []string{FW},
		},
		{
			name:         "TestConsulServiceMonitor_updasteMonitoredServices no services",
			serviceTypes: []string{},
		},
		{
			name:         "TestConsulServiceMonitor_updasteMonitoredServices multiple services",
			serviceTypes: []string{FW, IDS},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			serviceMonitor.monitoredServices = make(map[string]context.CancelFunc)
			services := make(map[string][]string)
			for _, serviceType := range tt.serviceTypes {
				services[serviceType] = []string{}
			}
			serviceMonitor.updateMonitoredServices(services)
			serviceMonitor.mu.RLock()
			for _, cancel := range serviceMonitor.monitoredServices {
				cancel()
			}
			serviceMonitor.wg.Wait()
		})
	}
}

func TestConsulServiceMonitor_queryAndUpdateMonitoredServices(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitor_queryAndUpdateMonitoredServices ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			serviceMonitor.monitoredServices = make(map[string]context.CancelFunc)
			lastIndex := uint64(0)
			serviceMonitor.queryAndUpdateMonitoredServices(&lastIndex)
		})
	}
}

func TestConsulServiceMonitor_Start(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitor_Start",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.Nil(t, err)
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				serviceMonitor.Start()
				wg.Done()
			}()
			time.Sleep(100 * time.Millisecond)
			close(serviceMonitor.stopChan)
			wg.Wait()
		})
	}
}

func TestConsulServiceMonitor_Stop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestConsulServiceMonitor_Stop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := cache.NewMockCache(gomock.NewController(t))
			serviceMonitor, err := NewConsulServiceMonitor(cacheMock, make(chan struct{}), "localhost")
			assert.NoError(t, err)
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				serviceMonitor.Start()
				wg.Done()
			}()
			serviceMonitor.Stop()
			wg.Wait()
		})
	}
}
