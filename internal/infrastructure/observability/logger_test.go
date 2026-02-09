package observability

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	// Ensure Logger is nil before test
	Logger = nil

	InitLogger()

	assert.NotNil(t, Logger)
}

func TestGetLogger(t *testing.T) {
	t.Run("Initialize if nil", func(t *testing.T) {
		Logger = nil
		logger := GetLogger()
		assert.NotNil(t, logger)
		assert.Equal(t, Logger, logger)
	})

	t.Run("Return existing", func(t *testing.T) {
		InitLogger()
		existingLogger := Logger
		logger := GetLogger()
		assert.Equal(t, existingLogger, logger)
	})
}
