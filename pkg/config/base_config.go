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
	jagwServiceAddress string
	jagwRequestPort    uint16
}

type BaseConfigInput struct {
	JagwServiceAddress string `validate:"required,hostname|ip"`
	JagwRequestPort    uint16 `validate:"required,gte=1,lte=65535"`
}

func NewBaseConfig(jagwServiceAddress, jagwRequestPort string) (*BaseConfig, error) {
	requestPortInt, err := strconv.ParseInt(jagwRequestPort, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Invalid Jalapeno API GW Request Port: %v", err)
	}
	requestPortUint := uint16(requestPortInt)

	baseConfigInput := &BaseConfigInput{
		JagwServiceAddress: jagwServiceAddress,
		JagwRequestPort:    requestPortUint,
	}

	validate := validator.New()
	if err := validate.Struct(baseConfigInput); err != nil {
		return nil, err
	}
	config := &BaseConfig{
		log:                logging.DefaultLogger.WithField("subsystem", Subsystem),
		jagwServiceAddress: baseConfigInput.JagwServiceAddress,
		jagwRequestPort:    baseConfigInput.JagwRequestPort,
	}
	return config, nil
}

func (config *BaseConfig) GetJagwServiceAddress() string {
	return config.jagwServiceAddress
}

func (config *BaseConfig) GetJagwRequestPort() uint16 {
	return config.jagwRequestPort
}

func (config *BaseConfig) GetJagwSubscriptionPort() uint16 {
	return 0
}

func (config *BaseConfig) GetGrpcPort() uint16 {
	return 0
}
