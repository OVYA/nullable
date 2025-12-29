package tests

import (
	"testing"

	"github.com/ovya/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullableEdgeCases(t *testing.T) {
	t.Run("SetValueP with nil pointer", func(t *testing.T) {
		test := testedStruct[embeddedStruct]{
			Name: nullable.Of[string]{},
		}
		test.Name.SetValueP(nil)

		assert.True(t, test.Name.IsNull(), "SetValueP(nil) should result in NULL value")
	})

	t.Run("SetValueP with value pointer", func(t *testing.T) {
		value := "test value"
		test := testedStruct[embeddedStruct]{
			Name: nullable.Of[string]{},
		}
		test.Name.SetValueP(&value)

		require.False(t, test.Name.IsNull(), "SetValueP(&value) should not result in NULL")
		assert.Equal(t, value, *test.Name.GetValue())
	})

	t.Run("IsNull on zero value", func(t *testing.T) {
		var test testedStruct[embeddedStruct]
		assert.True(t, test.Name.IsNull(), "Zero value should be NULL")
	})
}
