package helper

import (
	"os"
	"strconv"
	"strings"
)

var FlappingThreshold float64 = func() float64 {
	if value, exists := os.LookupEnv("HAWKEYE_FLAPPING_THRESHOLD"); exists {
		if temp, err := strconv.ParseFloat(value, 64); err == nil {
			return temp
		}
	}
	return 0.1
}()

var TwoFactorWeights = func() []float32 {
	if value, exists := os.LookupEnv("HAWKEYE_TWO_FACTOR_WEIGHTS"); exists {
		parts := strings.Split(value, ",")
		weights := make([]float32, len(parts))
		for i, part := range parts {
			if temp, err := strconv.ParseFloat(part, 32); err == nil {
				weights[i] = float32(temp)
			}
		}
		return weights
	}
	return []float32{0.7, 0.3}
}()

var ThreeFactorWeights = func() []float32 {
	if value, exists := os.LookupEnv("HAWKEYE_THREE_FACTOR_WEIGHTS"); exists {
		parts := strings.Split(value, ",")
		weights := make([]float32, len(parts))
		for i, part := range parts {
			if temp, err := strconv.ParseFloat(part, 32); err == nil {
				weights[i] = float32(temp)
			}
		}
		return weights
	}
	return []float32{0.5, 0.3, 0.2}
}()
