package tests

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ovya/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSON_NullValues(t *testing.T) {
	t.Run("null string", func(t *testing.T) {
		n := nullable.Null[string]()
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
	})

	t.Run("null int", func(t *testing.T) {
		n := nullable.Null[int]()
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
	})

	t.Run("null bool", func(t *testing.T) {
		n := nullable.Null[bool]()
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
	})

	t.Run("null float64", func(t *testing.T) {
		n := nullable.Null[float64]()
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
	})

	t.Run("null UUID", func(t *testing.T) {
		n := nullable.Null[uuid.UUID]()
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("null"), data)
	})
}

func TestMarshalJSON_PrimitiveTypes(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		n := nullable.FromValue("hello world")
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`"hello world"`), data)
	})

	t.Run("empty string", func(t *testing.T) {
		n := nullable.FromValue("")
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`""`), data)
	})

	t.Run("int value", func(t *testing.T) {
		n := nullable.FromValue(42)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("42"), data)
	})

	t.Run("int16 value", func(t *testing.T) {
		n := nullable.FromValue(int16(123))
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("123"), data)
	})

	t.Run("int32 value", func(t *testing.T) {
		n := nullable.FromValue(int32(456))
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("456"), data)
	})

	t.Run("int64 value", func(t *testing.T) {
		n := nullable.FromValue(int64(789))
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("789"), data)
	})

	t.Run("zero int", func(t *testing.T) {
		n := nullable.FromValue(0)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("0"), data)
	})

	t.Run("negative int", func(t *testing.T) {
		n := nullable.FromValue(-42)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("-42"), data)
	})

	t.Run("bool true", func(t *testing.T) {
		n := nullable.FromValue(true)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("true"), data)
	})

	t.Run("bool false", func(t *testing.T) {
		n := nullable.FromValue(false)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("false"), data)
	})

	t.Run("float64 value", func(t *testing.T) {
		n := nullable.FromValue(3.14159)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("3.14159"), data)
	})

	t.Run("float64 zero", func(t *testing.T) {
		n := nullable.FromValue(0.0)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("0"), data)
	})

	t.Run("float64 negative", func(t *testing.T) {
		n := nullable.FromValue(-2.5)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte("-2.5"), data)
	})
}

func TestMarshalJSON_SpecialFloatValues(t *testing.T) {
	t.Run("float64 NaN", func(t *testing.T) {
		n := nullable.FromValue(math.NaN())
		_, err := n.MarshalJSON()
		// JSON doesn't support NaN, so this should error
		assert.Error(t, err)
	})

	t.Run("float64 positive infinity", func(t *testing.T) {
		n := nullable.FromValue(math.Inf(1))
		_, err := n.MarshalJSON()
		// JSON doesn't support Inf, so this should error
		assert.Error(t, err)
	})

	t.Run("float64 negative infinity", func(t *testing.T) {
		n := nullable.FromValue(math.Inf(-1))
		_, err := n.MarshalJSON()
		// JSON doesn't support -Inf, so this should error
		assert.Error(t, err)
	})
}

func TestMarshalJSON_UUID(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		testUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		n := nullable.FromValue(testUUID)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`"550e8400-e29b-41d4-a716-446655440000"`), data)
	})

	t.Run("zero UUID", func(t *testing.T) {
		n := nullable.FromValue(uuid.UUID{})
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`"00000000-0000-0000-0000-000000000000"`), data)
	})
}

func TestMarshalJSON_JSONType(t *testing.T) {
	t.Run("simple map", func(t *testing.T) {
		obj := map[string]any{"key": "value", "number": 42}
		n := nullable.FromValue[nullable.JSON](obj)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		// Can't assert exact JSON due to map ordering, so unmarshal and compare
		var result map[string]any
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)
		assert.Equal(t, "value", result["key"])
		assert.Equal(t, float64(42), result["number"])
	})

	t.Run("nested structure", func(t *testing.T) {
		obj := map[string]any{
			"nested": map[string]any{
				"inner": "value",
			},
		}
		n := nullable.FromValue[nullable.JSON](obj)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		var result map[string]any
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)
		nested := result["nested"].(map[string]any)
		assert.Equal(t, "value", nested["inner"])
	})

	t.Run("array", func(t *testing.T) {
		obj := []any{1, 2, 3, "four"}
		n := nullable.FromValue[nullable.JSON](obj)
		data, err := n.MarshalJSON()
		require.NoError(t, err)
		assert.JSONEq(t, `[1,2,3,"four"]`, string(data))
	})
}

func TestUnmarshalJSON_NullValues(t *testing.T) {
	t.Run("null keyword to string", func(t *testing.T) {
		var n nullable.Of[string]
		err := n.UnmarshalJSON([]byte("null"))
		require.NoError(t, err)
		assert.True(t, n.IsNull())
		assert.Nil(t, n.GetValue())
	})

	t.Run("null keyword to int", func(t *testing.T) {
		var n nullable.Of[int]
		err := n.UnmarshalJSON([]byte("null"))
		require.NoError(t, err)
		assert.True(t, n.IsNull())
	})

	t.Run("nil byte slice", func(t *testing.T) {
		var n nullable.Of[string]
		err := n.UnmarshalJSON(nil)
		require.NoError(t, err)
		assert.True(t, n.IsNull())
	})

	t.Run("null to previously set value", func(t *testing.T) {
		n := nullable.FromValue("previous value")
		err := n.UnmarshalJSON([]byte("null"))
		require.NoError(t, err)
		assert.True(t, n.IsNull())
		assert.Nil(t, n.GetValue())
	})
}

func TestUnmarshalJSON_PrimitiveTypes(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		var n nullable.Of[string]
		err := n.UnmarshalJSON([]byte(`"hello world"`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, "hello world", *n.GetValue())
	})

	t.Run("empty string", func(t *testing.T) {
		var n nullable.Of[string]
		err := n.UnmarshalJSON([]byte(`""`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, "", *n.GetValue())
	})

	t.Run("string with quotes", func(t *testing.T) {
		var n nullable.Of[string]
		err := n.UnmarshalJSON([]byte(`"say \"hello\""`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, `say "hello"`, *n.GetValue())
	})

	t.Run("int value", func(t *testing.T) {
		var n nullable.Of[int]
		err := n.UnmarshalJSON([]byte("42"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, 42, *n.GetValue())
	})

	t.Run("int16 value", func(t *testing.T) {
		var n nullable.Of[int16]
		err := n.UnmarshalJSON([]byte("123"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, int16(123), *n.GetValue())
	})

	t.Run("int32 value", func(t *testing.T) {
		var n nullable.Of[int32]
		err := n.UnmarshalJSON([]byte("456"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, int32(456), *n.GetValue())
	})

	t.Run("int64 value", func(t *testing.T) {
		var n nullable.Of[int64]
		err := n.UnmarshalJSON([]byte("789"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, int64(789), *n.GetValue())
	})

	t.Run("zero int", func(t *testing.T) {
		var n nullable.Of[int]
		err := n.UnmarshalJSON([]byte("0"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, 0, *n.GetValue())
	})

	t.Run("negative int", func(t *testing.T) {
		var n nullable.Of[int]
		err := n.UnmarshalJSON([]byte("-42"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, -42, *n.GetValue())
	})

	t.Run("bool true", func(t *testing.T) {
		var n nullable.Of[bool]
		err := n.UnmarshalJSON([]byte("true"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, true, *n.GetValue())
	})

	t.Run("bool false", func(t *testing.T) {
		var n nullable.Of[bool]
		err := n.UnmarshalJSON([]byte("false"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, false, *n.GetValue())
	})

	t.Run("float64 value", func(t *testing.T) {
		var n nullable.Of[float64]
		err := n.UnmarshalJSON([]byte("3.14159"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, 3.14159, *n.GetValue())
	})

	t.Run("float64 zero", func(t *testing.T) {
		var n nullable.Of[float64]
		err := n.UnmarshalJSON([]byte("0.0"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, 0.0, *n.GetValue())
	})

	t.Run("float64 negative", func(t *testing.T) {
		var n nullable.Of[float64]
		err := n.UnmarshalJSON([]byte("-2.5"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, -2.5, *n.GetValue())
	})

	t.Run("float64 scientific notation", func(t *testing.T) {
		var n nullable.Of[float64]
		err := n.UnmarshalJSON([]byte("1.23e-4"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.InDelta(t, 0.000123, *n.GetValue(), 0.0000001)
	})
}

func TestUnmarshalJSON_UUID(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		var n nullable.Of[uuid.UUID]
		err := n.UnmarshalJSON([]byte(`"550e8400-e29b-41d4-a716-446655440000"`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		expected := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		assert.Equal(t, expected, *n.GetValue())
	})

	t.Run("zero UUID", func(t *testing.T) {
		var n nullable.Of[uuid.UUID]
		err := n.UnmarshalJSON([]byte(`"00000000-0000-0000-0000-000000000000"`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, uuid.UUID{}, *n.GetValue())
	})

	t.Run("invalid UUID format", func(t *testing.T) {
		var n nullable.Of[uuid.UUID]
		err := n.UnmarshalJSON([]byte(`"not-a-uuid"`))
		assert.Error(t, err)
	})
}

func TestUnmarshalJSON_InvalidJSON(t *testing.T) {
	t.Run("invalid JSON for string", func(t *testing.T) {
		var n nullable.Of[string]
		err := n.UnmarshalJSON([]byte(`not valid json`))
		assert.Error(t, err)
	})

	t.Run("invalid JSON for int", func(t *testing.T) {
		var n nullable.Of[int]
		err := n.UnmarshalJSON([]byte(`"not a number"`))
		assert.Error(t, err)
	})

	t.Run("invalid JSON for bool", func(t *testing.T) {
		var n nullable.Of[bool]
		err := n.UnmarshalJSON([]byte(`"not a bool"`))
		assert.Error(t, err)
	})

	t.Run("invalid JSON for float", func(t *testing.T) {
		var n nullable.Of[float64]
		err := n.UnmarshalJSON([]byte(`"not a float"`))
		assert.Error(t, err)
	})

	t.Run("number overflow for int16", func(t *testing.T) {
		var n nullable.Of[int16]
		err := n.UnmarshalJSON([]byte("100000"))
		assert.Error(t, err)
	})

	t.Run("number overflow for int32", func(t *testing.T) {
		var n nullable.Of[int32]
		err := n.UnmarshalJSON([]byte("10000000000"))
		assert.Error(t, err)
	})
}

func TestUnmarshalJSON_JSONType(t *testing.T) {
	t.Run("simple map", func(t *testing.T) {
		var n nullable.Of[nullable.JSON]
		err := n.UnmarshalJSON([]byte(`{"key":"value","number":42}`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		result := (*n.GetValue()).(map[string]any)
		assert.Equal(t, "value", result["key"])
		assert.Equal(t, float64(42), result["number"])
	})

	t.Run("nested structure", func(t *testing.T) {
		var n nullable.Of[nullable.JSON]
		err := n.UnmarshalJSON([]byte(`{"nested":{"inner":"value"}}`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		result := (*n.GetValue()).(map[string]any)
		nested := result["nested"].(map[string]any)
		assert.Equal(t, "value", nested["inner"])
	})

	t.Run("array", func(t *testing.T) {
		var n nullable.Of[nullable.JSON]
		err := n.UnmarshalJSON([]byte(`[1,2,3,"four"]`))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		result := (*n.GetValue()).([]any)
		assert.Len(t, result, 4)
		assert.Equal(t, float64(1), result[0])
		assert.Equal(t, "four", result[3])
	})
}

func TestMarshalUnmarshal_RoundTrip(t *testing.T) {
	t.Run("string round trip", func(t *testing.T) {
		original := nullable.FromValue("test value")
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored nullable.Of[string]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.Equal(t, *original.GetValue(), *restored.GetValue())
	})

	t.Run("int round trip", func(t *testing.T) {
		original := nullable.FromValue(42)
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored nullable.Of[int]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.Equal(t, *original.GetValue(), *restored.GetValue())
	})

	t.Run("bool round trip", func(t *testing.T) {
		original := nullable.FromValue(true)
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored nullable.Of[bool]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.Equal(t, *original.GetValue(), *restored.GetValue())
	})

	t.Run("float64 round trip", func(t *testing.T) {
		original := nullable.FromValue(3.14159)
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored nullable.Of[float64]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.Equal(t, *original.GetValue(), *restored.GetValue())
	})

	t.Run("UUID round trip", func(t *testing.T) {
		original := nullable.FromValue(uuid.New())
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored nullable.Of[uuid.UUID]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.Equal(t, *original.GetValue(), *restored.GetValue())
	})

	t.Run("null round trip", func(t *testing.T) {
		original := nullable.Null[string]()
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored nullable.Of[string]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		assert.True(t, restored.IsNull())
	})

	t.Run("JSON type round trip", func(t *testing.T) {
		obj := map[string]any{"key": "value", "number": float64(42)}
		original := nullable.FromValue[nullable.JSON](obj)
		data, err := original.MarshalJSON()
		require.NoError(t, err)

		var restored nullable.Of[nullable.JSON]
		err = restored.UnmarshalJSON(data)
		require.NoError(t, err)
		result := (*restored.GetValue()).(map[string]any)
		assert.Equal(t, "value", result["key"])
		assert.Equal(t, float64(42), result["number"])
	})
}

func TestMarshalUnmarshal_InStructs(t *testing.T) {
	type TestStruct struct {
		Name   nullable.Of[string]  `json:"name"`
		Age    nullable.Of[int]     `json:"age"`
		Active nullable.Of[bool]    `json:"active"`
		Score  nullable.Of[float64] `json:"score"`
	}

	t.Run("struct with all values", func(t *testing.T) {
		original := TestStruct{
			Name:   nullable.FromValue("John"),
			Age:    nullable.FromValue(30),
			Active: nullable.FromValue(true),
			Score:  nullable.FromValue(95.5),
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var restored TestStruct
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		assert.Equal(t, *original.Name.GetValue(), *restored.Name.GetValue())
		assert.Equal(t, *original.Age.GetValue(), *restored.Age.GetValue())
		assert.Equal(t, *original.Active.GetValue(), *restored.Active.GetValue())
		assert.Equal(t, *original.Score.GetValue(), *restored.Score.GetValue())
	})

	t.Run("struct with null values", func(t *testing.T) {
		original := TestStruct{
			Name:   nullable.Null[string](),
			Age:    nullable.Null[int](),
			Active: nullable.Null[bool](),
			Score:  nullable.Null[float64](),
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)
		assert.JSONEq(t, `{"name":null,"age":null,"active":null,"score":null}`, string(data))

		var restored TestStruct
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		assert.True(t, restored.Name.IsNull())
		assert.True(t, restored.Age.IsNull())
		assert.True(t, restored.Active.IsNull())
		assert.True(t, restored.Score.IsNull())
	})

	t.Run("struct with mixed null and non-null", func(t *testing.T) {
		original := TestStruct{
			Name:   nullable.FromValue("Jane"),
			Age:    nullable.Null[int](),
			Active: nullable.FromValue(false),
			Score:  nullable.Null[float64](),
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var restored TestStruct
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		assert.Equal(t, "Jane", *restored.Name.GetValue())
		assert.True(t, restored.Age.IsNull())
		assert.Equal(t, false, *restored.Active.GetValue())
		assert.True(t, restored.Score.IsNull())
	})
}

func TestUnmarshalJSON_OverwritingExisting(t *testing.T) {
	t.Run("overwrite value with new value", func(t *testing.T) {
		n := nullable.FromValue("original")
		err := n.UnmarshalJSON([]byte(`"new value"`))
		require.NoError(t, err)
		assert.Equal(t, "new value", *n.GetValue())
	})

	t.Run("overwrite value with null", func(t *testing.T) {
		n := nullable.FromValue(42)
		err := n.UnmarshalJSON([]byte("null"))
		require.NoError(t, err)
		assert.True(t, n.IsNull())
	})

	t.Run("overwrite null with value", func(t *testing.T) {
		n := nullable.Null[int]()
		err := n.UnmarshalJSON([]byte("123"))
		require.NoError(t, err)
		assert.False(t, n.IsNull())
		assert.Equal(t, 123, *n.GetValue())
	})
}

func TestMarshalUnmarshal_ComplexStructures(t *testing.T) {
	// Define complex nested structures
	type Address struct {
		Street     nullable.Of[string]  `json:"street"`
		City       nullable.Of[string]  `json:"city"`
		PostalCode nullable.Of[string]  `json:"postalCode"`
		Country    nullable.Of[string]  `json:"country"`
		Verified   nullable.Of[bool]    `json:"verified"`
		Lat        nullable.Of[float64] `json:"lat"`
		Lng        nullable.Of[float64] `json:"lng"`
	}

	type ContactInfo struct {
		Email       nullable.Of[string]        `json:"email"`
		Phone       nullable.Of[string]        `json:"phone"`
		Address     nullable.Of[nullable.JSON] `json:"address"`
		IsPrimary   nullable.Of[bool]          `json:"isPrimary"`
		LastUpdated nullable.Of[int64]         `json:"lastUpdated"`
	}

	type Metadata struct {
		Tags        nullable.Of[nullable.JSON] `json:"tags"`
		Properties  nullable.Of[nullable.JSON] `json:"properties"`
		Version     nullable.Of[int]           `json:"version"`
		IsActive    nullable.Of[bool]          `json:"isActive"`
		CreatedBy   nullable.Of[string]        `json:"createdBy"`
		CreatedByID nullable.Of[uuid.UUID]     `json:"createdById"`
	}

	type Profile struct {
		Bio         nullable.Of[string]        `json:"bio"`
		Website     nullable.Of[string]        `json:"website"`
		AvatarURL   nullable.Of[string]        `json:"avatarUrl"`
		Contacts    nullable.Of[nullable.JSON] `json:"contacts"`
		Preferences nullable.Of[nullable.JSON] `json:"preferences"`
		Metadata    nullable.Of[nullable.JSON] `json:"metadata"`
		Score       nullable.Of[float64]       `json:"score"`
		Level       nullable.Of[int32]         `json:"level"`
	}

	type User struct {
		ID          nullable.Of[uuid.UUID]     `json:"id"`
		Username    nullable.Of[string]        `json:"username"`
		Email       nullable.Of[string]        `json:"email"`
		FirstName   nullable.Of[string]        `json:"firstName"`
		LastName    nullable.Of[string]        `json:"lastName"`
		Age         nullable.Of[int]           `json:"age"`
		IsActive    nullable.Of[bool]          `json:"isActive"`
		Balance     nullable.Of[float64]       `json:"balance"`
		Profile     nullable.Of[nullable.JSON] `json:"profile"`
		Roles       nullable.Of[nullable.JSON] `json:"roles"`
		Permissions nullable.Of[nullable.JSON] `json:"permissions"`
		CreatedAt   nullable.Of[int64]         `json:"createdAt"`
	}

	t.Run("deeply nested structure with all values", func(t *testing.T) {
		// Create deeply nested structure
		address := Address{
			Street:     nullable.FromValue("123 Main St"),
			City:       nullable.FromValue("New York"),
			PostalCode: nullable.FromValue("10001"),
			Country:    nullable.FromValue("USA"),
			Verified:   nullable.FromValue(true),
			Lat:        nullable.FromValue(40.7128),
			Lng:        nullable.FromValue(-74.0060),
		}

		contact := ContactInfo{
			Email:       nullable.FromValue("user@example.com"),
			Phone:       nullable.FromValue("+1-555-0100"),
			Address:     nullable.FromValue[nullable.JSON](address),
			IsPrimary:   nullable.FromValue(true),
			LastUpdated: nullable.FromValue(int64(1234567890)),
		}

		metadata := Metadata{
			Tags:        nullable.FromValue[nullable.JSON]([]string{"premium", "verified"}),
			Properties:  nullable.FromValue[nullable.JSON](map[string]any{"theme": "dark", "language": "en"}),
			Version:     nullable.FromValue(3),
			IsActive:    nullable.FromValue(true),
			CreatedBy:   nullable.FromValue("admin"),
			CreatedByID: nullable.FromValue(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
		}

		profile := Profile{
			Bio:         nullable.FromValue("Software developer"),
			Website:     nullable.FromValue("https://example.com"),
			AvatarURL:   nullable.FromValue("https://example.com/avatar.jpg"),
			Contacts:    nullable.FromValue[nullable.JSON]([]ContactInfo{contact}),
			Preferences: nullable.FromValue[nullable.JSON](map[string]any{"notifications": true, "theme": "dark"}),
			Metadata:    nullable.FromValue[nullable.JSON](metadata),
			Score:       nullable.FromValue(98.5),
			Level:       nullable.FromValue(int32(42)),
		}

		user := User{
			ID:          nullable.FromValue(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")),
			Username:    nullable.FromValue("johndoe"),
			Email:       nullable.FromValue("john@example.com"),
			FirstName:   nullable.FromValue("John"),
			LastName:    nullable.FromValue("Doe"),
			Age:         nullable.FromValue(30),
			IsActive:    nullable.FromValue(true),
			Balance:     nullable.FromValue(1234.56),
			Profile:     nullable.FromValue[nullable.JSON](profile),
			Roles:       nullable.FromValue[nullable.JSON]([]string{"admin", "user"}),
			Permissions: nullable.FromValue[nullable.JSON](map[string]bool{"read": true, "write": true, "delete": false}),
			CreatedAt:   nullable.FromValue(int64(1609459200)),
		}

		// Marshal
		data, err := json.Marshal(user)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Unmarshal
		var restored User
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		// Verify top-level fields
		assert.Equal(t, *user.ID.GetValue(), *restored.ID.GetValue())
		assert.Equal(t, *user.Username.GetValue(), *restored.Username.GetValue())
		assert.Equal(t, *user.Email.GetValue(), *restored.Email.GetValue())
		assert.Equal(t, *user.Age.GetValue(), *restored.Age.GetValue())
		assert.Equal(t, *user.Balance.GetValue(), *restored.Balance.GetValue())

		// Verify roles (array)
		roles := (*restored.Roles.GetValue()).([]any)
		assert.Len(t, roles, 2)
		assert.Equal(t, "admin", roles[0])
		assert.Equal(t, "user", roles[1])

		// Verify permissions (map)
		permissions := (*restored.Permissions.GetValue()).(map[string]any)
		assert.Equal(t, true, permissions["read"])
		assert.Equal(t, true, permissions["write"])
		assert.Equal(t, false, permissions["delete"])

		// Verify nested profile
		profileData := (*restored.Profile.GetValue()).(map[string]any)
		assert.Equal(t, "Software developer", profileData["bio"])
		assert.Equal(t, "https://example.com", profileData["website"])
		assert.Equal(t, float64(98.5), profileData["score"])
		assert.Equal(t, float64(42), profileData["level"])

		// Verify deeply nested metadata
		metadataData := profileData["metadata"].(map[string]any)
		assert.Equal(t, float64(3), metadataData["version"])
		assert.Equal(t, true, metadataData["isActive"])
		assert.Equal(t, "admin", metadataData["createdBy"])
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", metadataData["createdById"])

		// Verify deeply nested tags
		tags := metadataData["tags"].([]any)
		assert.Len(t, tags, 2)
		assert.Equal(t, "premium", tags[0])
		assert.Equal(t, "verified", tags[1])
	})

	t.Run("deeply nested structure with mixed null values", func(t *testing.T) {
		// Create structure with some null values at different levels
		contact := ContactInfo{
			Email:       nullable.FromValue("contact@example.com"),
			Phone:       nullable.Null[string](),        // Null phone
			Address:     nullable.Null[nullable.JSON](), // Null address
			IsPrimary:   nullable.FromValue(false),
			LastUpdated: nullable.FromValue(int64(9876543210)),
		}

		metadata := Metadata{
			Tags:        nullable.FromValue[nullable.JSON]([]string{"new"}),
			Properties:  nullable.Null[nullable.JSON](), // Null properties
			Version:     nullable.FromValue(1),
			IsActive:    nullable.Null[bool](), // Null isActive
			CreatedBy:   nullable.FromValue("system"),
			CreatedByID: nullable.Null[uuid.UUID](), // Null UUID
		}

		profile := Profile{
			Bio:         nullable.Null[string](), // Null bio
			Website:     nullable.FromValue("https://site.com"),
			AvatarURL:   nullable.Null[string](), // Null avatar
			Contacts:    nullable.FromValue[nullable.JSON]([]ContactInfo{contact}),
			Preferences: nullable.Null[nullable.JSON](), // Null preferences
			Metadata:    nullable.FromValue[nullable.JSON](metadata),
			Score:       nullable.FromValue(75.0),
			Level:       nullable.Null[int32](), // Null level
		}

		user := User{
			ID:          nullable.FromValue(uuid.MustParse("abcd1234-e89b-12d3-a456-426614174000")),
			Username:    nullable.FromValue("janedoe"),
			Email:       nullable.Null[string](), // Null email
			FirstName:   nullable.FromValue("Jane"),
			LastName:    nullable.Null[string](), // Null last name
			Age:         nullable.FromValue(25),
			IsActive:    nullable.FromValue(true),
			Balance:     nullable.Null[float64](), // Null balance
			Profile:     nullable.FromValue[nullable.JSON](profile),
			Roles:       nullable.Null[nullable.JSON](), // Null roles
			Permissions: nullable.FromValue[nullable.JSON](map[string]bool{"read": true}),
			CreatedAt:   nullable.FromValue(int64(1609459200)),
		}

		// Marshal
		data, err := json.Marshal(user)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Unmarshal
		var restored User
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		// Verify null fields
		assert.True(t, restored.Email.IsNull())
		assert.True(t, restored.LastName.IsNull())
		assert.True(t, restored.Balance.IsNull())
		assert.True(t, restored.Roles.IsNull())

		// Verify non-null fields
		assert.Equal(t, *user.ID.GetValue(), *restored.ID.GetValue())
		assert.Equal(t, *user.Username.GetValue(), *restored.Username.GetValue())
		assert.Equal(t, *user.Age.GetValue(), *restored.Age.GetValue())

		// Verify nested profile with nulls
		profileData := (*restored.Profile.GetValue()).(map[string]any)
		assert.Nil(t, profileData["bio"])
		assert.Equal(t, "https://site.com", profileData["website"])
		assert.Nil(t, profileData["avatarUrl"])
		assert.Nil(t, profileData["preferences"])
		assert.Equal(t, float64(75.0), profileData["score"])
		assert.Nil(t, profileData["level"])

		// Verify deeply nested metadata with nulls
		metadataData := profileData["metadata"].(map[string]any)
		assert.Equal(t, float64(1), metadataData["version"])
		assert.Nil(t, metadataData["isActive"])
		assert.Nil(t, metadataData["properties"])
		assert.Nil(t, metadataData["createdById"])
	})

	t.Run("deeply nested structure with all null values", func(t *testing.T) {
		user := User{
			ID:          nullable.Null[uuid.UUID](),
			Username:    nullable.Null[string](),
			Email:       nullable.Null[string](),
			FirstName:   nullable.Null[string](),
			LastName:    nullable.Null[string](),
			Age:         nullable.Null[int](),
			IsActive:    nullable.Null[bool](),
			Balance:     nullable.Null[float64](),
			Profile:     nullable.Null[nullable.JSON](),
			Roles:       nullable.Null[nullable.JSON](),
			Permissions: nullable.Null[nullable.JSON](),
			CreatedAt:   nullable.Null[int64](),
		}

		// Marshal
		data, err := json.Marshal(user)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Verify all fields are null in JSON
		var jsonMap map[string]any
		err = json.Unmarshal(data, &jsonMap)
		require.NoError(t, err)
		for key, value := range jsonMap {
			assert.Nil(t, value, "Field %s should be null", key)
		}

		// Unmarshal
		var restored User
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		// Verify all fields are null
		assert.True(t, restored.ID.IsNull())
		assert.True(t, restored.Username.IsNull())
		assert.True(t, restored.Email.IsNull())
		assert.True(t, restored.FirstName.IsNull())
		assert.True(t, restored.LastName.IsNull())
		assert.True(t, restored.Age.IsNull())
		assert.True(t, restored.IsActive.IsNull())
		assert.True(t, restored.Balance.IsNull())
		assert.True(t, restored.Profile.IsNull())
		assert.True(t, restored.Roles.IsNull())
		assert.True(t, restored.Permissions.IsNull())
		assert.True(t, restored.CreatedAt.IsNull())
	})

	t.Run("array of complex structures", func(t *testing.T) {
		contact1 := ContactInfo{
			Email:       nullable.FromValue("contact1@example.com"),
			Phone:       nullable.FromValue("+1-555-0101"),
			IsPrimary:   nullable.FromValue(true),
			LastUpdated: nullable.FromValue(int64(1000000)),
		}

		contact2 := ContactInfo{
			Email:       nullable.FromValue("contact2@example.com"),
			Phone:       nullable.Null[string](),
			IsPrimary:   nullable.FromValue(false),
			LastUpdated: nullable.Null[int64](),
		}

		contact3 := ContactInfo{
			Email:       nullable.Null[string](),
			Phone:       nullable.Null[string](),
			IsPrimary:   nullable.Null[bool](),
			LastUpdated: nullable.Null[int64](),
		}

		contacts := []ContactInfo{contact1, contact2, contact3}

		// Marshal
		data, err := json.Marshal(contacts)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Unmarshal
		var restored []ContactInfo
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		assert.Len(t, restored, 3)

		// Verify first contact (all fields present)
		assert.Equal(t, "contact1@example.com", *restored[0].Email.GetValue())
		assert.Equal(t, "+1-555-0101", *restored[0].Phone.GetValue())
		assert.Equal(t, true, *restored[0].IsPrimary.GetValue())

		// Verify second contact (some nulls)
		assert.Equal(t, "contact2@example.com", *restored[1].Email.GetValue())
		assert.True(t, restored[1].Phone.IsNull())
		assert.Equal(t, false, *restored[1].IsPrimary.GetValue())
		assert.True(t, restored[1].LastUpdated.IsNull())

		// Verify third contact (all nulls)
		assert.True(t, restored[2].Email.IsNull())
		assert.True(t, restored[2].Phone.IsNull())
		assert.True(t, restored[2].IsPrimary.IsNull())
		assert.True(t, restored[2].LastUpdated.IsNull())
	})

	t.Run("map with complex nullable values", func(t *testing.T) {
		data := map[string]nullable.Of[nullable.JSON]{
			"user1": nullable.FromValue[nullable.JSON](map[string]any{
				"name":   "Alice",
				"age":    30,
				"active": true,
			}),
			"user2": nullable.FromValue[nullable.JSON](map[string]any{
				"name": "Bob",
				"age":  25,
			}),
			"user3": nullable.Null[nullable.JSON](),
		}

		// Marshal
		jsonData, err := json.Marshal(data)
		require.NoError(t, err)
		require.NotEmpty(t, jsonData)

		// Unmarshal
		var restored map[string]nullable.Of[nullable.JSON]
		err = json.Unmarshal(jsonData, &restored)
		require.NoError(t, err)

		assert.Len(t, restored, 3)

		// Verify user1
		user1Val := restored["user1"]
		user1 := (*user1Val.GetValue()).(map[string]any)
		assert.Equal(t, "Alice", user1["name"])
		assert.Equal(t, float64(30), user1["age"])
		assert.Equal(t, true, user1["active"])

		// Verify user2
		user2Val := restored["user2"]
		user2 := (*user2Val.GetValue()).(map[string]any)
		assert.Equal(t, "Bob", user2["name"])
		assert.Equal(t, float64(25), user2["age"])

		// Verify user3 is null
		user3Val := restored["user3"]
		assert.True(t, user3Val.IsNull())
	})

	t.Run("extreme nesting with 5 levels", func(t *testing.T) {
		// Level 5 (deepest)
		level5 := map[string]any{
			"value":    "deep value",
			"level":    5,
			"isDeep":   true,
			"metadata": []string{"tag1", "tag2", "tag3"},
		}

		// Level 4
		level4 := map[string]any{
			"data":    level5,
			"level":   4,
			"count":   42,
			"numbers": []int{1, 2, 3, 4, 5},
		}

		// Level 3
		level3 := map[string]any{
			"nested": level4,
			"level":  3,
			"active": true,
			"items":  []map[string]any{{"id": 1}, {"id": 2}},
		}

		// Level 2
		level2 := map[string]any{
			"inner":       level3,
			"level":       2,
			"description": "second level",
			"tags":        []string{"a", "b", "c"},
		}

		// Level 1 (top)
		level1 := nullable.FromValue[nullable.JSON](map[string]any{
			"root":  level2,
			"level": 1,
			"name":  "top level",
		})

		// Marshal
		data, err := json.Marshal(level1)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		// Unmarshal
		var restored nullable.Of[nullable.JSON]
		err = json.Unmarshal(data, &restored)
		require.NoError(t, err)

		// Navigate through all levels
		l1 := (*restored.GetValue()).(map[string]any)
		assert.Equal(t, float64(1), l1["level"])
		assert.Equal(t, "top level", l1["name"])

		l2 := l1["root"].(map[string]any)
		assert.Equal(t, float64(2), l2["level"])
		assert.Equal(t, "second level", l2["description"])

		l3 := l2["inner"].(map[string]any)
		assert.Equal(t, float64(3), l3["level"])
		assert.Equal(t, true, l3["active"])

		l4 := l3["nested"].(map[string]any)
		assert.Equal(t, float64(4), l4["level"])
		assert.Equal(t, float64(42), l4["count"])

		l5 := l4["data"].(map[string]any)
		assert.Equal(t, float64(5), l5["level"])
		assert.Equal(t, "deep value", l5["value"])
		assert.Equal(t, true, l5["isDeep"])

		// Verify deeply nested array
		metadata := l5["metadata"].([]any)
		assert.Len(t, metadata, 3)
		assert.Equal(t, "tag1", metadata[0])
		assert.Equal(t, "tag2", metadata[1])
		assert.Equal(t, "tag3", metadata[2])
	})
}

func TestMarshalUnmarshal(t *testing.T) {
	obj := getTestObjs(getEmbeddedObj())
	toObj := []testedStruct[embeddedStruct]{{}, {}}

	b, err := json.Marshal(obj)
	t.Run("Marshal nested structs test", func(t *testing.T) {
		require.NoError(t, err, "Marshaling Nullable data failed")
	})

	t.Run("Unmarshal tests suite", func(t *testing.T) {
		err = json.Unmarshal(b, &toObj)
		require.NoError(t, err, "Unmarshaling into Nullable data failed")

		for i := range 2 {
			t.Run("Simple string matching", func(t *testing.T) {
				assert.Equal(t, obj[i].Name.GetValue(), toObj[i].Name.GetValue(), "Name mismatch at index %d", i)
			})

			t.Run("Simple datetime matching", func(t *testing.T) {
				dte := toObj[i].DateTo.GetValue()
				if dte == nil {
					assert.Nil(t, obj[i].DateTo.GetValue(), "DateTo nil value mismatch at index %d", i)
				} else {
					assert.Equal(t, time.Duration(0), dte.Sub(now), "DateTo value mismatch at index %d", i)
				}
			})

			data := obj[i].Data.GetValue()
			if data == nil {
				t.Run("nil data into nil object checking", func(t *testing.T) {
					assert.Nil(t, toObj[i].Data.GetValue(), "Data nil value mismatch at index %d", i)
				})
			} else {
				t.Run("non nil data into non nil object matching", func(t *testing.T) {
					assert.Equal(t, data.Bool.GetValue(), toObj[i].Data.GetValue().Bool.GetValue(), "Data.Bool mismatch")
					assert.Equal(t, data.Int, toObj[i].Data.GetValue().Int, "Data.Int mismatch")
					assert.Equal(t, data.String, toObj[i].Data.GetValue().String, "Data.String mismatch")
				})
			}
		}
	})
}
