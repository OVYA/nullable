package nullable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// JSON permits to handle Postgresl Json[b] type
type JSON = any

type NullableI[T bool | int | int16 | int32 | int64 | string | uuid.UUID | float64 | JSON] interface {
	// IsNull returns true if itself is nil or the value is nil/null
	IsNull() bool
	// GetValue implements the getter.
	GetValue() *T
	// SetValue implements the setter.
	SetValue(T)
	// SetValueP implements the setter by pointer.
	SetValueP(*T)
	// SetNull set to null.
	SetNull()
	// MarshalJSON implements the encoding json interface.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON implements the decoding json interface.
	UnmarshalJSON([]byte) error
	// Value implements the driver.Valuer interface.
	Value() (driver.Value, error)
	// Scan implements the sql.Scanner interface.
	Scan(v any) error
}

// FromValue is a Nullable constructor from the given value thanks to Go generics' inference.
func FromValue[T bool | int | int16 | int32 | int64 | string | uuid.UUID | float64 | JSON](b T) *Of[T] {
	out := Of[T]{}
	out.SetValue(b)

	return &out
}

// Null is a Nullable constructor with Null value.
func Null[T bool | int | int16 | int32 | int64 | string | uuid.UUID | float64 | JSON]() *Of[T] {
	return &Of[T]{}
}

func (n *Of[T]) scanJSON(v any) error {
	null := sql.NullString{}
	err := null.Scan(v)
	if err != nil {
		return fmt.Errorf("nullable database scanning json : %w", err)
	}

	if null.Valid {
		value := new(T)
		err := json.Unmarshal([]byte(null.String), value)
		if err != nil {
			return fmt.Errorf("nullable database unmarshaling json : %w", err)
		}

		n.SetValue(*value)
	} else {
		n.SetNull()
	}

	return nil
}

func (n *Of[T]) scanString(v any) error {
	if n == nil {
		panic("Calling scanString on nil receiver")
	}

	null := sql.NullString{}
	err := null.Scan(v)
	if err != nil {
		return fmt.Errorf("nullable database scanning string : %w", err)
	}

	if null.Valid {
		n.SetValue(any(null.String).(T))
	} else {
		n.SetNull()
	}

	return nil
}

func (n *Of[T]) scanUUID(v any) error {
	if n == nil {
		panic("Calling scanUUID on nil receiver")
	}

	null := sql.NullString{}
	err := null.Scan(v)
	if err != nil {
		return fmt.Errorf("nullable database scanning string : %w", err)
	}

	if null.Valid {
		uid, err := uuid.Parse(null.String)
		if err != nil {
			return fmt.Errorf("UUID parsing failed : %w", err)
		}

		n.SetValue(any(uid).(T))
	} else {
		n.SetNull()
	}

	return nil
}

func (n *Of[T]) scanInt(v any) error {
	switch any(new(T)).(type) {
	case int16, *int16:
		null := sql.NullInt16{}
		err := null.Scan(v)
		if err != nil {
			return fmt.Errorf("nullable database scanning int16 : %w", err)
		}

		if null.Valid {
			n.SetValue(any(null.Int16).(T))
		} else {
			n.SetNull()
		}

		return nil
	case int32, *int32:
		null := sql.NullInt32{}
		err := null.Scan(v)
		if err != nil {
			return fmt.Errorf("nullable database scanning int32 : %w", err)
		}

		if null.Valid {
			n.SetValue(any(null.Int32).(T))
		} else {
			n.SetNull()
		}

		return nil
	case int, *int:
		null := sql.NullInt64{}
		err := null.Scan(v)
		if err != nil {
			return fmt.Errorf("nullable database scanning int : %w", err)
		}

		if null.Valid {
			n.SetValue(any(int(null.Int64)).(T))
		} else {
			n.SetNull()
		}

		return nil
	case int64, *int64:
		null := sql.NullInt64{}
		err := null.Scan(v)
		if err != nil {
			return fmt.Errorf("nullable database scanning int64 : %w", err)
		}

		if null.Valid {
			n.SetValue(any(null.Int64).(T))
		} else {
			n.SetNull()
		}

		return nil
	}

	return fmt.Errorf("type %T is not supported", *new(T))
}

func (n *Of[T]) scanFloat(v any) error {
	null := sql.NullFloat64{}
	err := null.Scan(v)
	if err != nil {
		return fmt.Errorf("nullable database scanning float64 : %w", err)
	}

	if null.Valid {
		n.SetValue(any(null.Float64).(T))
	} else {
		n.SetNull()
	}

	return nil
}

func (n *Of[T]) scanBool(v any) error {
	null := sql.NullBool{}
	err := null.Scan(v)
	if err != nil {
		return fmt.Errorf("nullable database scanning bool : %w", err)
	}

	if null.Valid {
		n.SetValue(any(null.Bool).(T))
	} else {
		n.SetNull()
	}

	return nil
}

func (n *Of[T]) scanTime(v any) error {
	if v == nil {
		n.SetNull()

		return nil
	}

	null := new(sql.NullTime)

	switch t := v.(type) {
	case string:
		var err error
		null.Time, err = time.Parse(t, t)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	case time.Time:
		err := null.Scan(v)
		if err != nil {
			return fmt.Errorf("nullable database scanning Time : %w", err)
		}
	default:
		return fmt.Errorf("canot parse type \"%T\" with value \"%v\" to time", t, t)
	}

	if null.Valid {
		n.SetValue(any(null.Time).(T))
	} else {
		n.SetNull()
	}

	return nil
}

// marshalJSON implements the generic encoding json interface.
func marshalJSON[T any](nullable NullableI[T]) ([]byte, error) {
	b, err := json.Marshal(nullable.GetValue())
	if err != nil {
		return nil, fmt.Errorf("nullable json marshaling %T : %w", nullable, err)
	}

	return b, nil
}
