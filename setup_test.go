package nullable

import (
	"os"
	"testing"
)

// Setup
func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
