package config

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
)

type FullConfig struct {
	*BaseConfig
	jagwSubscriptionPort uint16 `validate:"required,gte=1,lte=65535"`
	grpcPort             uint16 `validate:"required,gte=1,lte=65535"`
}

func NewFullConfig(jagwServiceAddress, jagwRequestPort, jagwSubscriptionPort, grpcPort string) (*FullConfig, error) {
	baseConfig, err := NewBaseConfig(jagwServiceAddress, jagwRequestPort)
	if err != nil {
		return nil, fmt.Errorf("Error creating base config: %v", err)
	}

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

	config := &FullConfig{
		BaseConfig:           baseConfig,
		jagwSubscriptionPort: jagwSubscriptionPortUint,
		grpcPort:             grpcPortUint,
	}
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *FullConfig) GetJagwServiceAddress() string {
	return c.jagwServiceAddress
}

func (c *FullConfig) GetJagwRequestPort() uint16 {
	return c.jagwRequestPort
}

func (c *FullConfig) GetJagwSubscriptionPort() uint16 {
	return c.jagwSubscriptionPort
}

func (c *FullConfig) GetGrpcPort() uint16 {
	return c.grpcPort
}
