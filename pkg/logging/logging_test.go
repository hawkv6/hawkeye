package logging

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestIntializeDefaultLogger(t *testing.T) {
	assert.NotNil(t, InitializeDefaultLogger())
	assert.Equal(t, logrus.InfoLevel, InitializeDefaultLogger().Level)
}
func TestIntializeDefaultLoggerDebug(t *testing.T) {
	os.Setenv("HAWKV6_DEBUG", "true")
	assert.NotNil(t, InitializeDefaultLogger())
	assert.Equal(t, logrus.DebugLevel, InitializeDefaultLogger().Level)
	os.Unsetenv("HAWKV6_DEBUG")
}
