package nullable_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ovya/nullable"
	"github.com/stretchr/testify/assert"
)

type Test[T any] struct {
	ID     int64                   `json:"id"`
	Name   *nullable.Of[string]    `json:"name"`
	DateTo *nullable.Of[time.Time] `json:"dateTo"`
	Data   *nullable.Of[T]         `json:"data"`
}

type testedType = struct {
	String string                      `json:"string"`
	Bool   *nullable.Of[bool]          `json:"bool"`
	Int    int64                       `json:"int"`
	JSON   *nullable.Of[nullable.JSON] `json:"json"`
}

func TestMarshalUnmarshal(t *testing.T) {
	data1 := testedType{
		String: "a string",
		Bool:   nullable.FromValue(true),
		Int:    768,
	}

	data1.JSON = nullable.FromValue[nullable.JSON](data1)

	name := "PLOP"
	now := time.Now()
	obj1 := Test[testedType]{
		Name:   nullable.FromValue(name),
		DateTo: nullable.FromValue(now),
		Data:   nullable.FromValue(data1),
	}

	obj2 := Test[testedType]{
		Name:   nullable.Null[string](),     // Null value
		DateTo: nullable.Null[time.Time](),  // Null value
		Data:   nullable.Null[testedType](), // Null value
	}

	obj := []Test[testedType]{obj1, obj2}

	b, err := json.Marshal(obj)

	if !assert.Nil(t, err, "Marshaling Nullable data failed") {
		panic(err)
	}

	toObj := []Test[testedType]{{}, {}}
	err = json.Unmarshal(b, &toObj)

	if !assert.Nil(t, err, "Unmarshaling into Nullable data failed") {
		panic(err)
	}

	for i := 0; i < 2; i++ {
		assert.Equal(t, obj[i].Name.GetValue(), toObj[i].Name.GetValue(), "Marshaling and Unmarshaling does not return the same value")
		dte := toObj[i].DateTo.GetValue()
		if dte == nil {
			assert.Nil(t, obj[i].DateTo.GetValue(), "Marshaling and Unmarshaling nil value mismatch.")
		} else {
			assert.Equal(t, time.Duration(0), dte.Sub(now), "Marshaling and Unmarshaling does not return the same value")
		}

		data := obj[i].Data.GetValue()
		if data == nil {
			assert.Nil(t, toObj[i].Data.GetValue(), "Marshaling and Unmarshaling nil value mismatch")
		} else {
			assert.Equal(t, data.Bool.GetValue(), toObj[i].Data.GetValue().Bool.GetValue(), "Marshaling and Unmarshaling does not return the same value")
			assert.Equal(t, data.Int, toObj[i].Data.GetValue().Int, "Marshaling and Unmarshaling does not return the same value")
			assert.Equal(t, data.String, toObj[i].Data.GetValue().String, "Marshaling and Unmarshaling does not return the same value")
			assert.Equal(t, data.String, toObj[i].Data.GetValue().String, "Marshaling and Unmarshaling does not return the same value")
		}
	}
}
