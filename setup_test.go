package nullable

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// Setup
func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	code := m.Run()
	os.Exit(code)
}

func TestEmpty(t *testing.T) {
	assert.True(t, true, "C'est pas faux !")
}
