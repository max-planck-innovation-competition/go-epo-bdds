package epo_bbds

import (
	"os"
	"testing"
)

func init() {
	// os.Setenv("TEST", "true")
}

func skipTest(t *testing.T) {
	if os.Getenv("TEST") == "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestNewFeature(t *testing.T) {
	skipTest(t)
}
