package config

const Subsystem = "config"

type Config interface {
	GetJagwServiceAddress() string
	GetJagwRequestPort() uint16
	GetJagwSubscriptionPort() uint16
	GetGrpcPort() uint16
}
