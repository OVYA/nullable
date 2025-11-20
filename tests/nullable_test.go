package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ovya/nullable"
	"github.com/stretchr/testify/assert"
)

func TestMarshalUnmarshal(t *testing.T) {
	obj := getTestObjs(getEmbeddedObj())
	toObj := []testedStruct[embeddedStruct]{{}, {}}

	b, err := json.Marshal(obj)

	if !assert.Nil(t, err, "Marshaling Nullable data failed") {
		panic(err)
	}

	err = json.Unmarshal(b, &toObj)

	if !assert.Nil(t, err, "Unmarshaling into Nullable data failed") {
		panic(err)
	}

	for i := range 2 {
		assert.Equal(t, obj[i].Name.GetValue(), toObj[i].Name.GetValue(), fmt.Sprintf("Marshaling and Unmarshaling does not return the same value. i = %d ; s = '%v'", i, toObj[i].Name.GetValue()))
		dte := toObj[i].DateTo.GetValue()
		if dte == nil {
			assert.Nil(t, obj[i].DateTo.GetValue(), fmt.Sprintf("Marshaling and Unmarshaling nil value mismatch. i = %d", i))
		} else {
			assert.Equal(t, time.Duration(0), dte.Sub(now), fmt.Sprintf("Marshaling and Unmarshaling does not return the same value. i = %d", i))
		}

		data := obj[i].Data.GetValue()
		if data == nil {
			assert.Nil(t, toObj[i].Data.GetValue(), fmt.Sprintf("Marshaling and Unmarshaling nil value mismatch. i = %d", i))
		} else {
			assert.Equal(t, data.Bool.GetValue(), toObj[i].Data.GetValue().Bool.GetValue(), "Marshaling and Unmarshaling does not return the same value")
			assert.Equal(t, data.Int, toObj[i].Data.GetValue().Int, "Marshaling and Unmarshaling does not return the same value")
			assert.Equal(t, data.String, toObj[i].Data.GetValue().String, "Marshaling and Unmarshaling does not return the same value")
			assert.Equal(t, data.String, toObj[i].Data.GetValue().String, "Marshaling and Unmarshaling does not return the same value")
		}
	}
}

func TestNullableEdgeCases(t *testing.T) {
	t.Run("SetValueP with nil pointer", func(t *testing.T) {
		test := testedStruct[embeddedStruct]{
			Name: nullable.Of[string]{},
		}
		test.Name.SetValueP(nil)

		if !test.Name.IsNull() {
			t.Error("SetValueP(nil) should result in NULL value")
		}
	})

	t.Run("SetValueP with value pointer", func(t *testing.T) {
		value := "test value"
		test := testedStruct[embeddedStruct]{
			Name: nullable.Of[string]{},
		}
		test.Name.SetValueP(&value)

		if test.Name.IsNull() {
			t.Error("SetValueP(&value) should not result in NULL")
		}
		if *test.Name.GetValue() != value {
			t.Errorf("Expected '%s', got '%s'", value, *test.Name.GetValue())
		}
	})

	t.Run("IsNull on zero value", func(t *testing.T) {
		var test testedStruct[embeddedStruct]
		if !test.Name.IsNull() {
			t.Error("Zero value should be NULL")
		}
	})
}
