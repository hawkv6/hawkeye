package service

const Subsystem = "service"

type Service interface {
	GetType() string
	GetId() string
	GetSid() string
	IsHealty() bool
}

type ConcreteService struct {
	serviceType string
	serviceId   string
	prefixSid   string
	healthy     bool
}

func NewConcreteService(serviceType, id, sid string, healthy bool) *ConcreteService {
	return &ConcreteService{
		serviceType: serviceType,
		serviceId:   id,
		prefixSid:   sid,
		healthy:     healthy,
	}
}

func (service *ConcreteService) GetType() string {
	return service.serviceType
}

func (service *ConcreteService) GetId() string {
	return service.serviceId
}

func (service *ConcreteService) GetSid() string {
	return service.prefixSid
}
func (service *ConcreteService) IsHealty() bool {
	return service.healthy
}
