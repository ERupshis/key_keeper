package configutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetEnvToParamIfNeed(t *testing.T) {
	// Test case 1: Set int64 parameter
	var intValue int64
	assert.NoError(t, SetEnvToParamIfNeed(&intValue, "123"))
	if intValue != 123 {
		t.Errorf("Expected intValue to be 123, got %d", intValue)
	}

	// Test case 2: Set string parameter
	var stringValue string
	assert.NoError(t, SetEnvToParamIfNeed(&stringValue, "testString"))
	if stringValue != "testString" {
		t.Errorf("Expected stringValue to be 'testString', got '%s'", stringValue)
	}

	// Test case 3: Empty value, should not modify parameters
	assert.NoError(t, SetEnvToParamIfNeed(&intValue, ""))
	if intValue != 123 {
		t.Errorf("Expected intValue to remain 123, got %d", intValue)
	}

	assert.NoError(t, SetEnvToParamIfNeed(&stringValue, ""))
	if stringValue != "testString" {
		t.Errorf("Expected stringValue to remain 'testString', got '%s'", stringValue)
	}

	// Test case 4: Wrong input type, should panic with an error message
	assert.Error(t, SetEnvToParamIfNeed(42, "test"))

	// Test case 5: Set timeDuration parameter
	var timeValue time.Duration
	assert.NoError(t, SetEnvToParamIfNeed(&timeValue, "1h"))
	if intValue != 123 {
		t.Errorf("Expected intValue to be 123, got %s", timeValue)
	}

	// Test case 6: Set timeDuration parameter wrong
	assert.Error(t, SetEnvToParamIfNeed(&timeValue, "123"))
	if intValue != 123 {
		t.Errorf("Expected intValue to be 123, got %s", timeValue)
	}

}
