package service

type ServiceMonitor interface {
	StartMonitoring()
	StopMonitoring()
}
