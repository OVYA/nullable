# Go Nullable

[![golangci-lint](https://github.com/ovya/nullable/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/ovya/nullable/actions/workflows/golangci-lint.yml)
[![mod-verify](https://github.com/OVYA/nullable/actions/workflows/mod-verify.yml/badge.svg)](https://github.com/OVYA/nullable/actions/workflows/mod-verify.yml)
[![gosec](https://github.com/OVYA/nullable/actions/workflows/gosec.yaml/badge.svg)](https://github.com/OVYA/nullable/actions/workflows/gosec.yaml)
[![staticcheck](https://github.com/OVYA/nullable/actions/workflows/staticcheck.yaml/badge.svg)](https://github.com/OVYA/nullable/actions/workflows/staticcheck.yaml)
[![test](https://github.com/OVYA/nullable/actions/workflows/test.yml/badge.svg)](https://github.com/OVYA/nullable/actions/workflows/test.yml)

* Provide Go database null value for *any* data type as JSON thanks to the Golang generic features.
* Support data type `github.com/google/uuid.UUID`.
* Make possible to scan and store any structs' type into json and jsonb Postgresql data type.
* Support JSON marshaling and unMarshaling with conventional
  Javascript/Typescript value `null` instead of `Valid:true/false, Type:value`
  as `database/sql` does.

## Usages

### From the Go test

```go
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

data1 := testedType{
	String: "a string",
	Bool:   nullable.FromValue(true),
	Int:    768,
}

data1.JSON = nullable.FromValue[nullable.JSON](data1)

obj1 := Test[testedType]{
	Name:   nullable.FromValue("PLOP"),
	DateTo: nullable.FromValue(time.Now()),
	Data:   nullable.FromValue(data1),
}

obj2 := Test[testedType]{
	Name:   nullable.Null[string](),     // Null value
	DateTo: nullable.Null[time.Time](),  // Null value
	Data:   nullable.Null[testedType](), // Null value
}
```

### Database Insertion

Comes from a test on Postgresql database with Time and JSON insertion :

```
// The model of the db table test
type Test[T nullable.JSON] struct {
	ID int64 `json:"id"`
	// This string can be null
	Name *nullable.Of[string] `json:"name"`
	// This timestamp can be null
	DateTo *nullable.Of[time.Time] `json:"dateTo"`
	// You can any interface you want in data
	Data *nullable.Of[T] `json:"data"`
}

type dataType = struct {
	String string `json:"string"`
	Bool   bool   `json:"bool"`
	Int    int64  `json:"int"`
}

data := dataType{
	String: "This is a string",
	Bool:   true,
	Int:    768,
}

obj := Test[dataType]{
	Name:      nullable.FromValue("My name"),
	DateTo:    nullable.FromValue(time.Now()),
	Data:      nullable.FromValue(data),
}

_, err = daoService.NamedExec("INSERT INTO test (name, date_to, data) VALUES (:name, :date_to, :data)", obj)

if err != nil…
```

### Custom primitive data type

The library `Nullable` persist and read the structured types as `JSON`
in/from the database but it is possible to workaround this feature.

Suppose you need the custom type `type PhoneNumber string`.
If you use this type as is :
```
type Coordinate struct {
	Email      *nullable.Of[string]            `json:"email,omitempty"`
	Phone      *nullable.Of[phonenum.PhoneNum] `json:"phone"`
}
```

persisting a `Coordinate` in database will persist the field `Phone` in `JSON`, what we don't want.
In order to persist (resp. read) the type `Phone` as a `string`, you
must implement the native Go database `driver.Valuer` (reps. `sql.Scanner`) :

```
// Value implements the driver.Valuer interface.
func (pn *PhoneNumb) Value() (driver.Value, error) {
	return driver.Value(string(*pn)), nil
}

// Scan implements the sql.Scanner interface.
func (pn *PhoneNum) Scan(v any) error {
	switch val := v.(type) {
	case int, int64, uint64:
		*pn = PhoneNum(strconv.Itoa(val.(int)))
	case string:
		*pn = PhoneNum(val)
	default:
		return errors.New(fmt.Sprintf("can not scan phone number from type %T", val))
	}

	return nil
}
```

This `Scanner` can even read an integer from the database to fulfill the phone number as a `string`.

## Notes

### Similar Project

This project is inspired from
[gonull](https://github.com/lomsa-dev/gonull) that fails
scanning/storing some Postgresql type like `enum`, `timestamp` and `json`/`jsonb`.

### Go Tests for Postgresql

Go tests storing and scanning data from/to Postgresql database are in
progress and will be available soon in an other `git` repository.
