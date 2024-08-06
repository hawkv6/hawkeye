package helper

import (
	"os"
	"strconv"
	"strings"
	"time"
)

var FlappingThreshold float64 = func() float64 {
	if value, exists := os.LookupEnv("HAWKEYE_FLAPPING_THRESHOLD"); exists {
		if temp, err := strconv.ParseFloat(value, 64); err == nil {
			return temp
		}
	}
	return 0.1
}()

var TwoFactorWeights = func() []float64 {
	if value, exists := os.LookupEnv("HAWKEYE_TWO_FACTOR_WEIGHTS"); exists {
		parts := strings.Split(value, ",")
		weights := make([]float64, len(parts))
		for i, part := range parts {
			if temp, err := strconv.ParseFloat(part, 64); err == nil {
				weights[i] = temp
			}
		}
		return weights
	}
	return []float64{0.7, 0.3}
}()

var ThreeFactorWeights = func() []float64 {
	if value, exists := os.LookupEnv("HAWKEYE_THREE_FACTOR_WEIGHTS"); exists {
		parts := strings.Split(value, ",")
		weights := make([]float64, len(parts))
		for i, part := range parts {
			if temp, err := strconv.ParseFloat(part, 64); err == nil {
				weights[i] = temp
			}
		}
		return weights
	}
	return []float64{0.7, 0.2, 0.1}
}()

var ConsulQueryWaitTime time.Duration = func() time.Duration {
	if value, exists := os.LookupEnv("HAWKEYE_CONSUL_QUERY_WAIT_TIME"); exists {
		if temp, err := strconv.ParseInt(value, 10, 64); err == nil {
			return time.Duration(temp) * time.Second
		}
	}
	return 5 * time.Second
}()

var RollingWindowSize uint8 = func() uint8 {
	if value, exists := os.LookupEnv("HAWKEYE_ROLLING_WINDOWS_SIZE"); exists {
		if temp, err := strconv.ParseUint(value, 10, 8); err == nil {
			return uint8(temp)
		}
	}
	return 5
}()

var NetworkProcessorHoldTime time.Duration = func() time.Duration {
	if value, exists := os.LookupEnv("HAWKEYE_NETWORK_PROCESSOR_HOLD_TIME"); exists {
		if temp, err := strconv.ParseInt(value, 10, 64); err == nil {
			return time.Duration(temp) * time.Second
		}
	}
	return 1 * time.Second
}()
