package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/hawkv6/hawkeye/pkg/cache"
	"github.com/hawkv6/hawkeye/pkg/helper"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type ConsulServiceMonitor struct {
	log               *logrus.Entry
	stopChan          chan struct{}
	client            *api.Client
	services          map[string]map[string]*ConcreteService
	monitoredServices map[string]context.CancelFunc
	cache             cache.Cache
	updapteChan       chan struct{}
	needsUpdate       bool
	mu                sync.RWMutex
}

func NewConsulServiceMonitor(cache cache.Cache, updateChan chan struct{}) (*ConsulServiceMonitor, error) {
	config := api.DefaultConfig()
	config.Address = helper.ConsulServerAddress
	config.Scheme = "https"
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	consultServiceMonitor := &ConsulServiceMonitor{
		log:               logging.DefaultLogger.WithField("subsystem", Subsystem),
		stopChan:          make(chan struct{}),
		client:            client,
		services:          make(map[string]map[string]*ConcreteService, 0),
		monitoredServices: make(map[string]context.CancelFunc, 0),
		cache:             cache,
		updapteChan:       updateChan,
		needsUpdate:       false,
		mu:                sync.RWMutex{},
	}
	return consultServiceMonitor, nil
}

func (monitor *ConsulServiceMonitor) createServiceEntry(serviceId, serviceType, prefixSid string, healthy bool) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	if _, ok := monitor.services[serviceType]; !ok {
		monitor.services[serviceType] = make(map[string]*ConcreteService)
	}
	if _, ok := monitor.services[serviceType][serviceId]; !ok {
		monitor.services[serviceType][serviceId] = NewConcreteService(serviceType, serviceId, prefixSid, healthy)
		if healthy {
			monitor.storeServiceSidInCache(serviceType, prefixSid)
		}
		monitor.log.Infof("Service %s from type %s with sid %s created - healthy: %t", serviceId, serviceType, prefixSid, healthy)
	}
}

func (monitor *ConsulServiceMonitor) deleteServiceEntry(serviceId, serviceType string) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	monitor.cache.Lock()
	monitor.cache.RemoveServiceSid(serviceType, serviceId)
	monitor.cache.Unlock()
	delete(monitor.services[serviceType], serviceId)
	monitor.log.Debugf("Service %s deleted", serviceId)
	monitor.needsUpdate = true
}

func (monitor *ConsulServiceMonitor) deleteServiceEntries(knownServices map[string]bool, serviceType string) {
	for serviceId := range knownServices {
		if !knownServices[serviceId] {
			monitor.deleteServiceEntry(serviceId, serviceType)
		}
	}
}

func (monitor *ConsulServiceMonitor) getKnownServices(serviceEntries []*api.ServiceEntry) (map[string]bool, error) {
	if len(serviceEntries) == 0 {
		return nil, fmt.Errorf("No service entries found")
	}
	knownServices := make(map[string]bool, len(serviceEntries))
	serviceType := serviceEntries[0].Service.Service
	for serviceId := range monitor.services[serviceType] {
		knownServices[serviceId] = false
	}
	return knownServices, nil
}

func (monitor *ConsulServiceMonitor) validateServiceEntries(serviceEntries []*api.ServiceEntry) {
	knownServices, err := monitor.getKnownServices(serviceEntries)
	if len(serviceEntries) == 0 {
		return
	}
	if err != nil {
		return
	}
	for _, serviceEntry := range serviceEntries {
		prefixSid := serviceEntry.Service.Meta["sid"]
		serviceId := serviceEntry.Service.ID
		if prefixSid == "" {
			monitor.log.Warnf("SID not found for service %s", serviceId)
		} else {
			healthState := serviceEntry.Checks.AggregatedStatus() == api.HealthPassing
			monitor.createServiceEntry(serviceId, serviceEntry.Service.Service, prefixSid, healthState)
			knownServices[serviceId] = true
		}
	}

	serviceType := serviceEntries[0].Service.Service
	monitor.deleteServiceEntries(knownServices, serviceType)
}

func (monitor *ConsulServiceMonitor) storeServiceSidInCache(serviceType, prefixSid string) {
	monitor.cache.Lock()
	monitor.cache.StoreServiceSid(serviceType, prefixSid)
	monitor.cache.Unlock()
	monitor.needsUpdate = true
}

func (monitor *ConsulServiceMonitor) removeServiceSidInCache(serviceType, prefixSid string) {
	monitor.cache.Lock()
	monitor.cache.RemoveServiceSid(serviceType, prefixSid)
	monitor.cache.Unlock()
	monitor.needsUpdate = true
}

func (monitor *ConsulServiceMonitor) markServiceUnhealthy(service *ConcreteService) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	monitor.log.Infof("Service %s turned unhealthy", service.serviceId)
	service.healthy = false
	monitor.removeServiceSidInCache(service.serviceType, service.prefixSid)
}

func (monitor *ConsulServiceMonitor) markServiceHealthy(service *ConcreteService) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	monitor.log.Infof("Service %s is healthy again", service.serviceId)
	service.healthy = true
	monitor.storeServiceSidInCache(service.serviceType, service.prefixSid)
}

func (monitor *ConsulServiceMonitor) checkServiceHealth(serviceEntries []*api.ServiceEntry) {
	for _, serviceEntry := range serviceEntries {
		healthy := serviceEntry.Checks.AggregatedStatus() == api.HealthPassing
		service := monitor.services[serviceEntry.Service.Service][serviceEntry.Service.ID]
		if service.healthy != healthy {
			if healthy {
				monitor.markServiceHealthy(service)
			} else {
				monitor.markServiceUnhealthy(service)
			}
		} else {
			monitor.log.Debugf("Service %s health state unchanged - healthy: %t ", service.serviceId, service.healthy)
		}
	}
}

func (monitor *ConsulServiceMonitor) monitorServiceHealth(ctx context.Context, serviceType string) {
	lastIndex := uint64(0)
	for {
		select {
		case <-ctx.Done():
			monitor.log.Infoln("Stopping monitoring health state for service type:", serviceType)
			return
		default:
			monitor.needsUpdate = false
			options := &api.QueryOptions{
				WaitIndex: lastIndex,
				WaitTime:  helper.ConsulQueryWaitTime,
			}
			serviceEntries, meta, err := monitor.client.Health().Service(serviceType, "", false, options)
			if err != nil {
				monitor.log.Errorf("Error checking service health: %v", err)
				continue
			}
			monitor.validateServiceEntries(serviceEntries)
			monitor.checkServiceHealth(serviceEntries)
			if monitor.needsUpdate {
				monitor.log.Info("Services or Service health changed - sending update message")
				monitor.updapteChan <- struct{}{}
			}
			lastIndex = meta.LastIndex
		}
	}
}

func (monitor *ConsulServiceMonitor) servicesChanged(services map[string][]string) bool {
	monitor.mu.RLock()
	defer monitor.mu.RUnlock()
	for service := range services {
		if service == "consul" {
			continue
		}
		if _, exists := monitor.monitoredServices[service]; !exists {
			return true
		}
	}
	return false
}

func (monitor *ConsulServiceMonitor) stopMonitoringRemovedServices(newServiceTypes map[string][]string) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	for serviceType, cancel := range monitor.monitoredServices {
		if _, exists := newServiceTypes[serviceType]; !exists {
			cancel()
			delete(monitor.monitoredServices, serviceType)
		}
	}
}
func (monitor *ConsulServiceMonitor) startServiceHealthMonitoring(services map[string][]string) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	for service := range services {
		if _, exists := monitor.monitoredServices[service]; !exists && service != "consul" {
			ctx, cancel := context.WithCancel(context.Background())
			monitor.monitoredServices[service] = cancel
			go monitor.monitorServiceHealth(ctx, service)
		}
	}
}

func (monitor *ConsulServiceMonitor) updateMonitoredServices(services map[string][]string) {
	if !monitor.servicesChanged(services) {
		return
	}
	monitor.startServiceHealthMonitoring(services)
	monitor.stopMonitoringRemovedServices(services)
}

func (monitor *ConsulServiceMonitor) StartMonitoring() {
	monitor.log.Infoln("Starting monitoring services")
	lastIndex := uint64(0)
	for {
		select {
		case <-monitor.stopChan:
			for _, cancel := range monitor.monitoredServices {
				cancel()
			}
			return
		default:
			options := &api.QueryOptions{
				WaitIndex: lastIndex,
				WaitTime:  helper.ConsulQueryWaitTime,
			}
			services, meta, err := monitor.client.Catalog().Services(options)
			if err != nil {
				monitor.log.Errorf("Error fetching services: %v", err)
				continue
			}
			monitor.updateMonitoredServices(services)
			lastIndex = meta.LastIndex
		}
	}
}

func (monitor *ConsulServiceMonitor) StopMonitoring() {
	monitor.log.Infoln("Stopping monitoring services")
	close(monitor.stopChan)
}
