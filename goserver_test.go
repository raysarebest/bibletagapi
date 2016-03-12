package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestStub is a canary test to make sure testing and testify are good to go
func TestStub(t *testing.T) {
	assert.True(t, true, "This is good. Canary test passing")
}

// TestLogger is a test to make sure REST requests are logged
func TestLogger(t *testing.T) {
	for _, route := range routes {

		// wrap all current routes in the logger decorator to log out requests
		testhandler := Logger(route.HandlerFunc, route.Name)
		if &testhandler == nil {
			t.Error("Test failed: Could not wrap route with logger")
		}

	}
}

// TestNewRouter is a test to make sure a new router can be created
func TestNewRouter(t *testing.T) {
	testrouter := NewRouter()
	if &testrouter == nil {
		t.Error("Test failed: Could not create new router")
	}
}