package nullable

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Of[T bool | int | int16 | int32 | int64 | string | uuid.UUID | float64 | JSON] struct {
	//nolint: tagliatelle // Internal use
	Val *T `json:"nullable_value" db:"_"`
}

// IsNull returns true iff the value is nil
func (n *Of[T]) IsNull() bool {
	return n == nil || n.Val == nil
}

// GetValue implements the getter.
func (n *Of[T]) GetValue() *T {
	if n == nil {
		return nil
	}

	return n.Val
}

// SetValue implements the setter.
func (n *Of[T]) SetValue(b T) {
	n.Val = &b
}

// SetValueP implements the setter by pointer.
// If ref is not nil, calls SetValue(*ref)
// If ref is nil, calls SetNull()
func (n *Of[T]) SetValueP(ref *T) {
	if n == nil {
		n = new(Of[T])
	}

	if ref != nil {
		n.SetValue(*ref)
	} else {
		n.SetNull()
	}
}

// SetNull set to null.
func (n *Of[T]) SetNull() {
	if n == nil {
		return
	}

	n.Val = nil
}

// MarshalJSON implements the encoding json interface.
func (n *Of[T]) MarshalJSON() ([]byte, error) {
	if n == nil {
		b, _ := json.Marshal(nil)

		return b, nil
	}

	return marshalJSON[T](n)
}

// UnmarshalJSON implements the decoding json interface.
func (n *Of[T]) UnmarshalJSON(data []byte) error {
	if n == nil {
		return nil
	}

	if n.Val == nil && data != nil {
		n.Val = new(T)
	}

	if data == nil {
		return nil
	}

	err := json.Unmarshal(data, n.Val)
	if err != nil {
		return fmt.Errorf("nullable Unmarshal Error : %w", err)
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (n *Of[T]) Value() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}

	switch value := any(n.Val).(type) {
	case *string, *int16, *int32, *int, *int64, *float64, *bool, *time.Time, *uuid.UUID:
		return *n.Val, nil
	case JSON:
		if value == nil {
			return nil, nil
		}

		b, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("nullable database value error : %w", err)
		}

		return string(b), nil
	}

	return nil, fmt.Errorf("type %T is not supported for value %v", *n.Val, *n.Val)
}

// Scan implements the sql.Scanner interface.
// This method decodes a JSON-encoded value into the struct.
func (n *Of[T]) Scan(v any) error {
	if n == nil {
		return errors.New("calling Scan on nil receiver")
	}

	switch any(n.Val).(type) {
	case *string:
		return n.scanString(v)
	case *uuid.UUID:
		return n.scanUUID(v)
	case *int16, *int32, *int, *int64:
		return n.scanInt(v)
	case *float64:
		return n.scanFloat(v)
	case *bool:
		return n.scanBool(v)
	case *time.Time:
		return n.scanTime(v)
	case *JSON, JSON:
		return n.scanJSON(v)
	}

	return fmt.Errorf("type %T is not handled as nullable", v)
}
