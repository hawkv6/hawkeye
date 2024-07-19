package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-playground/validator"
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
	updateChan        chan struct{}
	needsUpdate       bool
	mu                sync.RWMutex
	wg                sync.WaitGroup
}

type ConsulServiceMonitorInput struct {
	Address string `validate:"required,hostname|ip"`
}

func NewConsulServiceMonitor(cache cache.Cache, updateChan chan struct{}, address string) (*ConsulServiceMonitor, error) {
	input := &ConsulServiceMonitorInput{
		Address: address,
	}
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return nil, err
	}
	config := api.DefaultConfig()
	config.Address = address
	config.Scheme = "https"
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConsulServiceMonitor{
		log:               logging.DefaultLogger.WithField("subsystem", Subsystem),
		stopChan:          make(chan struct{}),
		client:            client,
		services:          make(map[string]map[string]*ConcreteService, 0),
		monitoredServices: make(map[string]context.CancelFunc, 0),
		cache:             cache,
		updateChan:        updateChan,
		needsUpdate:       false,
		mu:                sync.RWMutex{},
		wg:                sync.WaitGroup{},
	}, nil
}

func (monitor *ConsulServiceMonitor) storeServiceSidInCache(serviceType, prefixSid string) {
	monitor.cache.Lock()
	monitor.cache.StoreServiceSid(serviceType, prefixSid)
	monitor.cache.Unlock()
	monitor.needsUpdate = true
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

func (monitor *ConsulServiceMonitor) removeServiceSidInCache(serviceType, prefixSid string) {
	monitor.cache.Lock()
	monitor.cache.RemoveServiceSid(serviceType, prefixSid)
	monitor.cache.Unlock()
	monitor.needsUpdate = true
}

func (monitor *ConsulServiceMonitor) deleteServiceEntry(serviceId, serviceType string) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	service := monitor.services[serviceType][serviceId]
	monitor.removeServiceSidInCache(serviceType, service.prefixSid)
	delete(monitor.services[serviceType], serviceId)
	monitor.log.Debugf("Service %s deleted", serviceId)
}

func (monitor *ConsulServiceMonitor) deleteServiceEntries(knownServices map[string]bool, serviceType string) {
	for serviceId := range knownServices {
		if !knownServices[serviceId] {
			monitor.deleteServiceEntry(serviceId, serviceType)
		}
	}
}

func (monitor *ConsulServiceMonitor) getKnownServices(serviceEntries []*api.ServiceEntry) (map[string]bool, string, error) {
	if len(serviceEntries) == 0 {
		return nil, "", fmt.Errorf("No service entries found")
	}
	knownServices := make(map[string]bool, len(serviceEntries))
	serviceType := serviceEntries[0].Service.Service
	for serviceId := range monitor.services[serviceType] {
		knownServices[serviceId] = false
	}
	return knownServices, serviceType, nil
}

func (monitor *ConsulServiceMonitor) processServiceEntries(serviceEntries []*api.ServiceEntry, knownServices map[string]bool, serviceType string) {
	for _, serviceEntry := range serviceEntries {
		prefixSid := serviceEntry.Service.Meta["sid"]
		serviceId := serviceEntry.Service.ID
		if prefixSid == "" {
			monitor.log.Warnf("SID not found for service %s", serviceId)
		} else {
			healthState := serviceEntry.Checks.AggregatedStatus() == api.HealthPassing
			if _, exists := knownServices[serviceId]; !exists {
				monitor.createServiceEntry(serviceId, serviceType, prefixSid, healthState)
			} else {
				knownServices[serviceId] = true
			}
		}
	}
}

func (monitor *ConsulServiceMonitor) validateServiceEntries(serviceEntries []*api.ServiceEntry) {
	knownServices, serviceType, err := monitor.getKnownServices(serviceEntries)
	if err == nil {
		monitor.processServiceEntries(serviceEntries, knownServices, serviceType)
		monitor.deleteServiceEntries(knownServices, serviceType)
	}
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
		currentHealth := serviceEntry.Checks.AggregatedStatus() == api.HealthPassing
		serviceType := serviceEntry.Service.Service
		serviceId := serviceEntry.Service.ID
		service := monitor.services[serviceType][serviceId]
		if service.healthy == currentHealth {
			monitor.log.Debugf("Service %s health state unchanged - healthy: %t ", service.serviceId, service.healthy)
		} else {
			if currentHealth {
				monitor.markServiceHealthy(service)
			} else {
				monitor.markServiceUnhealthy(service)
			}
		}
	}
}

func (monitor *ConsulServiceMonitor) sendUpdate() {
	if monitor.needsUpdate {
		monitor.log.Info("Services or service health changed - sending update message")
		monitor.updateChan <- struct{}{}
	}
}

func (monitor *ConsulServiceMonitor) getQueryOptions(lastIndex *uint64) *api.QueryOptions {
	return &api.QueryOptions{
		WaitIndex: *lastIndex,
		WaitTime:  helper.ConsulQueryWaitTime,
	}
}
func (monitor *ConsulServiceMonitor) queryAndUpdateServiceHealth(serviceType string, lastIndex *uint64) bool {
	monitor.needsUpdate = false
	serviceEntries, meta, err := monitor.client.Health().Service(serviceType, "", false, monitor.getQueryOptions(lastIndex))
	if err != nil {
		monitor.log.Errorf("Error checking service health for %s: %v", serviceType, err)
		return true // continue polling
	}
	monitor.validateServiceEntries(serviceEntries)
	monitor.checkServiceHealth(serviceEntries)
	monitor.sendUpdate()
	*lastIndex = meta.LastIndex
	return true
}

func (monitor *ConsulServiceMonitor) monitorServiceHealth(ctx context.Context, serviceType string) {
	lastIndex := uint64(0)
	for {
		select {
		case <-ctx.Done():
			monitor.log.Infoln("Stopping monitoring health state for service type:", serviceType)
			return
		default:
			monitor.queryAndUpdateServiceHealth(serviceType, &lastIndex)
		}
	}
}

func (monitor *ConsulServiceMonitor) unmonitoredServicesExist(services map[string][]string) bool {
	monitor.mu.RLock()
	defer monitor.mu.RUnlock()
	for service := range services {
		if _, exists := monitor.monitoredServices[service]; !exists && service != "consul" {
			return true
		}
	}
	return false
}

func (monitor *ConsulServiceMonitor) stopMonitoringRemovedServices(services map[string][]string) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	for serviceType, cancel := range monitor.monitoredServices {
		if _, exists := services[serviceType]; !exists {
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
			monitor.wg.Add(1)
			go func(service string) {
				monitor.monitorServiceHealth(ctx, service)
				monitor.wg.Done()
			}(service)
		}
	}
}

func (monitor *ConsulServiceMonitor) updateMonitoredServices(services map[string][]string) {
	if monitor.unmonitoredServicesExist(services) {
		monitor.startServiceHealthMonitoring(services)
		monitor.stopMonitoringRemovedServices(services)
	}
}

func (monitor *ConsulServiceMonitor) stopMonitoredServices() {
	monitor.log.Infof("Stopping monitoring services, can take up to %v", 2*helper.ConsulQueryWaitTime)
	monitor.mu.RLock()
	defer monitor.mu.RUnlock()
	for _, cancel := range monitor.monitoredServices {
		cancel()
	}
}

func (monitor *ConsulServiceMonitor) queryAndUpdateMonitoredServices(lastIndex *uint64) {
	services, meta, err := monitor.client.Catalog().Services(monitor.getQueryOptions(lastIndex))
	if err != nil {
		monitor.log.Errorf("Error fetching services: %v", err)
		return
	}
	monitor.updateMonitoredServices(services)
	*lastIndex = meta.LastIndex
}

func (monitor *ConsulServiceMonitor) Start() {
	monitor.log.Infoln("Starting monitoring services")
	lastIndex := uint64(0)
	for {
		select {
		case <-monitor.stopChan:
			monitor.stopMonitoredServices()
			return
		default:
			monitor.queryAndUpdateMonitoredServices(&lastIndex)
		}
	}
}

func (monitor *ConsulServiceMonitor) Stop() {
	close(monitor.stopChan)
	monitor.wg.Wait()
}
