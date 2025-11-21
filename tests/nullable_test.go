package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ovya/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshal(t *testing.T) {
	obj := getTestObjs(getEmbeddedObj())
	toObj := []testedStruct[embeddedStruct]{{}, {}}

	b, err := json.Marshal(obj)
	require.NoError(t, err, "Marshaling Nullable data failed")

	err = json.Unmarshal(b, &toObj)
	require.NoError(t, err, "Unmarshaling into Nullable data failed")

	for i := range 2 {
		assert.Equal(t, obj[i].Name.GetValue(), toObj[i].Name.GetValue(), "Name mismatch at index %d", i)

		dte := toObj[i].DateTo.GetValue()
		if dte == nil {
			assert.Nil(t, obj[i].DateTo.GetValue(), "DateTo nil value mismatch at index %d", i)
		} else {
			assert.Equal(t, time.Duration(0), dte.Sub(now), "DateTo value mismatch at index %d", i)
		}

		data := obj[i].Data.GetValue()
		if data == nil {
			assert.Nil(t, toObj[i].Data.GetValue(), "Data nil value mismatch at index %d", i)
		} else {
			assert.Equal(t, data.Bool.GetValue(), toObj[i].Data.GetValue().Bool.GetValue(), "Data.Bool mismatch")
			assert.Equal(t, data.Int, toObj[i].Data.GetValue().Int, "Data.Int mismatch")
			assert.Equal(t, data.String, toObj[i].Data.GetValue().String, "Data.String mismatch")
		}
	}
}

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
