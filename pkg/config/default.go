package config

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type DefaultConfig struct {
	log                  *logrus.Entry
	jagwServiceAddress   string `validate:"required,hostname"`
	jagwRequestPort      uint16 `validate:"required,gte=1,lte=65535"`
	jagwSubscriptionPort uint16 `validate:"required,gte=1,lte=65535"`
	grpcPort             uint16 `validate:"required,gte=1,lte=65535"`
}

func NewDefaultConfig(jagwServiceAddress, jagwRequestPort, jagwSubscriptionPort, grpcPort string) (*DefaultConfig, error) {
	jagwRequestPortInt, err := strconv.ParseInt(jagwRequestPort, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Invalid JAGW Request Port: %v", err)
	}
	jagwRequestPortUint := uint16(jagwRequestPortInt)

	jagwSubscriptionPortInt, err := strconv.ParseInt(jagwSubscriptionPort, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Invalid JAGW Subscription Port: %v", err)
	}
	jagwSubscriptionPortUint := uint16(jagwSubscriptionPortInt)

	grpcPortInt, err := strconv.ParseInt(grpcPort, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Invalid gRPC Port: %v", err)
	}
	grpcPortUint := uint16(grpcPortInt)
	config := &DefaultConfig{
		log:                  logging.DefaultLogger.WithField("subsystem", Subsystem),
		jagwServiceAddress:   jagwServiceAddress,
		jagwRequestPort:      jagwRequestPortUint,
		jagwSubscriptionPort: jagwSubscriptionPortUint,
		grpcPort:             grpcPortUint,
	}
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *DefaultConfig) GetJagwServiceAddress() string {
	return c.jagwServiceAddress
}

func (c *DefaultConfig) GetJagwRequestPort() uint16 {
	return c.jagwRequestPort
}

func (c *DefaultConfig) GetJagwSubscriptionPort() uint16 {
	return c.jagwSubscriptionPort
}

func (c *DefaultConfig) GetGrpcPort() uint16 {
	return c.grpcPort
}
