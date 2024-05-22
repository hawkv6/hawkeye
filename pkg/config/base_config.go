package config

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/hawkv6/hawkeye/pkg/logging"
	"github.com/sirupsen/logrus"
)

type BaseConfig struct {
	log                *logrus.Entry
	jagwServiceAddress string `validate:"required,hostname"`
	jagwRequestPort    uint16 `validate:"required,gte=1,lte=65535"`
}

func NewBaseConfig(jagwServiceAddress, jagwRequestPort string) (*BaseConfig, error) {
	requestPortInt, err := strconv.ParseInt(jagwRequestPort, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Invalid Jalapeno API GW Request Port: %v", err)
	}
	requestPortUint := uint16(requestPortInt)

	config := &BaseConfig{
		log:                logging.DefaultLogger.WithField("subsystem", Subsystem),
		jagwServiceAddress: jagwServiceAddress,
		jagwRequestPort:    requestPortUint,
	}
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *BaseConfig) GetJagwServiceAddress() string {
	return c.jagwServiceAddress
}

func (c *BaseConfig) GetJagwRequestPort() uint16 {
	return c.jagwRequestPort
}

func (c *BaseConfig) GetJagwSubscriptionPort() uint16 {
	return 0
}

func (c *BaseConfig) GetGrpcPort() uint16 {
	return 0
}
