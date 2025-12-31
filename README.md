# Go Nullable

[![golangci-lint](https://github.com/ovya/nullable/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/ovya/nullable/actions/workflows/golangci-lint.yml)
[![mod-verify](https://github.com/OVYA/nullable/actions/workflows/mod-verify.yml/badge.svg)](https://github.com/OVYA/nullable/actions/workflows/mod-verify.yml)
[![gosec](https://github.com/OVYA/nullable/actions/workflows/gosec.yaml/badge.svg)](https://github.com/OVYA/nullable/actions/workflows/gosec.yaml)
[![staticcheck](https://github.com/OVYA/nullable/actions/workflows/staticcheck.yaml/badge.svg)](https://github.com/OVYA/nullable/actions/workflows/staticcheck.yaml)
[![test](https://github.com/OVYA/nullable/actions/workflows/test.yml/badge.svg)](https://github.com/OVYA/nullable/actions/workflows/test.yml)

## This library was developed by [pivaldi](https://github.com/pivaldi/) and is now maintained [HERE](https://github.com/pivaldi/nullable).

A type-safe nullable value library for Go using generics, designed for seamless JSON marshaling and database operations.

## Features

- **Type-safe nullable values** for any supported type using Go generics
- **Database-friendly** with built-in `sql.Scanner` and `driver.Valuer` implementations
- **JSON marshaling** that uses standard `null` instead of `{Valid: true, Value: ...}`
- **PostgreSQL JSON/JSONB support** for storing complex types
- **UUID support** with `github.com/google/uuid`
- **Zero external dependencies** (except `google/uuid`)
- **Fully tested** with comprehensive unit and integration tests

## Installation

```bash
go get github.com/ovya/nullable
```

## Quick Start

```go
import "github.com/ovya/nullable"

// Create nullable values
name := nullable.FromValue("John Doe")
age := nullable.FromValue(30)
email := nullable.Null[string]() // Explicitly null

// Check if null
if name.IsNull() {
    // Handle null case
}

// Get value
if !age.IsNull() {
    fmt.Println(*age.GetValue()) // 30
}
```

## Supported Types

The library supports the following types through the `Of[T]` generic wrapper:

- **Integers**: `int`, `int16`, `int32`, `int64`
- **Floating point**: `float64`
- **Boolean**: `bool`
- **String**: `string`
- **UUID**: `uuid.UUID` (from `github.com/google/uuid`)
- **JSON**: `nullable.JSON` (alias for `any`) - for complex types stored as JSON in database

## Why Use This Library?

### Standard `database/sql` Approach

```go
type User struct {
    Name sql.NullString `json:"name"`
    Age  sql.NullInt64  `json:"age"`
}

// JSON output:
// {"name":{"String":"John","Valid":true},"age":{"Int64":30,"Valid":true}}
```

### With `nullable`

```go
type User struct {
    Name nullable.Of[string] `json:"name"`
    Age  nullable.Of[int]    `json:"age"`
}

// JSON output:
// {"name":"John","age":30}
// or with null values:
// {"name":null,"age":null}
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/ovya/nullable"
)

type User struct {
    ID       nullable.Of[int]    `json:"id"`
    Name     nullable.Of[string] `json:"name"`
    Email    nullable.Of[string] `json:"email"`
    Age      nullable.Of[int]    `json:"age"`
    IsActive nullable.Of[bool]   `json:"isActive"`
}

func main() {
    // Create user with some null fields
    user := User{
        ID:       nullable.FromValue(1),
        Name:     nullable.FromValue("John Doe"),
        Email:    nullable.Null[string](), // Null email
        Age:      nullable.FromValue(30),
        IsActive: nullable.FromValue(true),
    }

    // Marshal to JSON
    data, _ := json.Marshal(user)
    fmt.Println(string(data))
    // Output: {"id":1,"name":"John Doe","email":null,"age":30,"isActive":true}

    // Unmarshal from JSON
    jsonStr := `{"id":2,"name":"Jane Doe","email":"jane@example.com","age":null,"isActive":false}`
    var user2 User
    json.Unmarshal([]byte(jsonStr), &user2)

    fmt.Println(*user2.Name.GetValue())    // "Jane Doe"
    fmt.Println(user2.Age.IsNull())        // true
}
```

### Database Operations

#### Insert

```go
import (
    "database/sql"
    "time"
    "github.com/ovya/nullable"
    _ "github.com/jackc/pgx/v5/stdlib"
)

type Article struct {
    ID          int64                  `db:"id"`
    Title       nullable.Of[string]    `db:"title"`
    Content     nullable.Of[string]    `db:"content"`
    PublishedAt nullable.Of[time.Time] `db:"published_at"`
    AuthorID    nullable.Of[int64]     `db:"author_id"`
}

func insertArticle(db *sql.DB) error {
    article := Article{
        Title:       nullable.FromValue("My Article"),
        Content:     nullable.FromValue("Article content here..."),
        PublishedAt: nullable.FromValue(time.Now()),
        AuthorID:    nullable.Null[int64](), // Anonymous article
    }

    query := `
        INSERT INTO articles (title, content, published_at, author_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

    return db.QueryRow(
        query,
        article.Title,
        article.Content,
        article.PublishedAt,
        article.AuthorID,
    ).Scan(&article.ID)
}
```

#### Query

```go
func getArticle(db *sql.DB, id int64) (*Article, error) {
    var article Article

    query := `
        SELECT id, title, content, published_at, author_id
        FROM articles
        WHERE id = $1
    `

    err := db.QueryRow(query, id).Scan(
        &article.ID,
        &article.Title,
        &article.Content,
        &article.PublishedAt,
        &article.AuthorID,
    )

    if err != nil {
        return nil, err
    }

    return &article, nil
}
```

### Working with JSON/JSONB (PostgreSQL)

Store complex Go types as JSON in PostgreSQL:

```go
type Metadata struct {
    Tags       []string          `json:"tags"`
    Properties map[string]string `json:"properties"`
    Version    int               `json:"version"`
}

type Document struct {
    ID       int64                      `db:"id"`
    Title    nullable.Of[string]        `db:"title"`
    Metadata nullable.Of[nullable.JSON] `db:"metadata"` // Stored as JSONB
}

func insertDocument(db *sql.DB) error {
    meta := Metadata{
        Tags:       []string{"golang", "database"},
        Properties: map[string]string{"type": "article", "lang": "en"},
        Version:    1,
    }

    doc := Document{
        Title:    nullable.FromValue("Go Nullable Guide"),
        Metadata: nullable.FromValue[nullable.JSON](meta),
    }

    query := `INSERT INTO documents (title, metadata) VALUES ($1, $2) RETURNING id`
    return db.QueryRow(query, doc.Title, doc.Metadata).Scan(&doc.ID)
}
```

### Nested Structures

```go
type Address struct {
    Street  nullable.Of[string] `json:"street"`
    City    nullable.Of[string] `json:"city"`
    ZipCode nullable.Of[string] `json:"zipCode"`
}

type Profile struct {
    Bio     nullable.Of[string]        `json:"bio"`
    Website nullable.Of[string]        `json:"website"`
    Address nullable.Of[nullable.JSON] `json:"address"`
}

type User struct {
    Username nullable.Of[string]        `json:"username"`
    Email    nullable.Of[string]        `json:"email"`
    Profile  nullable.Of[nullable.JSON] `json:"profile"`
}

func main() {
    user := User{
        Username: nullable.FromValue("johndoe"),
        Email:    nullable.FromValue("john@example.com"),
        Profile: nullable.FromValue[nullable.JSON](Profile{
            Bio:     nullable.FromValue("Software Developer"),
            Website: nullable.FromValue("https://johndoe.com"),
            Address: nullable.FromValue[nullable.JSON](Address{
                Street:  nullable.FromValue("123 Main St"),
                City:    nullable.FromValue("New York"),
                ZipCode: nullable.FromValue("10001"),
            }),
        }),
    }

    data, _ := json.MarshalIndent(user, "", "  ")
    fmt.Println(string(data))
}
```

### Custom Types with Scanner/Valuer

For custom primitive types that should be stored as their underlying type (not JSON):

```go
import (
    "database/sql/driver"
    "errors"
    "fmt"
    "strconv"
)

type PhoneNumber string

// Value implements driver.Valuer to store as string in database
func (pn PhoneNumber) Value() (driver.Value, error) {
    return string(pn), nil
}

// Scan implements sql.Scanner to read from database
func (pn *PhoneNumber) Scan(v any) error {
    switch val := v.(type) {
    case int, int64, uint64:
        *pn = PhoneNumber(strconv.Itoa(val.(int)))
    case string:
        *pn = PhoneNumber(val)
    default:
        return errors.New(fmt.Sprintf("cannot scan phone number from type %T", val))
    }
    return nil
}

// Now PhoneNumber will be stored as string, not JSON
type Contact struct {
    Email nullable.Of[string]      `db:"email"`
    Phone nullable.Of[PhoneNumber] `db:"phone"` // Stored as string, not JSON
}
```

## API Reference

### Creating Nullable Values

```go
// From a value
name := nullable.FromValue("John")

// Explicitly null
email := nullable.Null[string]()

// From a pointer (nil pointer becomes null)
var ptr *string = nil
value := nullable.Of[string]{}
value.SetValueP(ptr) // Sets to null
```

### Checking and Accessing Values

```go
// Check if null
if value.IsNull() {
    // Handle null
}

// Get value (returns *T)
if !value.IsNull() {
    v := value.GetValue()
    fmt.Println(*v)
}
```

### Setting Values

```go
var value nullable.Of[string]

// Set a value
value.SetValue("hello")

// Set from pointer
str := "world"
value.SetValueP(&str)

// Set to null
value.SetNull()
```

### JSON Operations

```go
// Marshal
data, err := json.Marshal(value)

// Unmarshal
var value nullable.Of[string]
err := json.Unmarshal([]byte(`"hello"`), &value)

// Unmarshal null
err := json.Unmarshal([]byte(`null`), &value)
// value.IsNull() == true
```

## Testing

Run all tests including PostgreSQL integration tests:

```bash
cd tests
go test -v ./...
```

Or from the root:

```bash
make test
```

**Requirements:**
- Docker must be running (testcontainers uses Docker to spin up PostgreSQL)
- No manual database setup needed - testcontainers handles everything

**First run:** Tests will download the PostgreSQL 18 image (~80MB), subsequent runs use cached image.

Run only unit tests (no database required):

```bash
cd tests
go test -run 'TestMarshal|TestUnmarshal|TestNullableEdgeCases' -v
```

## Comparison with Alternatives

| Feature | `nullable` | `database/sql.Null*` | `gopkg.in/guregu/null.v4` |
|---------|-----------|---------------------|--------------------------|
| Type-safe generics | ✅ | ❌ (separate type per kind) | ❌ (separate type per kind) |
| Clean JSON output | ✅ `null` | ❌ `{"Valid":false}` | ✅ `null` |
| PostgreSQL JSON/JSONB | ✅ | ❌ | ⚠️ Limited |
| UUID support | ✅ | ❌ | ❌ |
| Custom types | ✅ via Scanner/Valuer | ✅ via Scanner/Valuer | ✅ via Scanner/Valuer |
| Zero dependencies* | ✅ | ✅ | ❌ |

*Except `google/uuid` for UUID support

### Detailed Comparison with `aarondl/opt`

The [`opt` package](https://github.com/aarondl/opt) is another modern approach to nullable values in Go, but with a fundamentally different philosophy.

#### Core Difference: 2-State vs 3-State Model

**`nullable` (2-state model):**
```go
type User struct {
    Name nullable.Of[string]  // Can be: null OR "John"
}
// Zero value is null
// No distinction between "field not provided" and "field set to null"
```

**`opt` (3-state model):**
```go
import "github.com/aarondl/opt/omitnull"

type User struct {
    Name omitnull.Val[string]  // Can be: unset OR null OR "John"
}
// Zero value is unset (omitted)
// Distinguishes: not provided vs explicitly null vs actual value
```

#### When the 3-State Model Matters

The distinction between "unset" and "null" is crucial for **partial API updates**:

```json
// Request 1: Update name, clear age
{"name": "John", "age": null}

// Request 2: Update name only, don't touch age
{"name": "John"}

// Request 3: Clear both fields
{"name": null, "age": null}
```

**With `nullable`:** Cannot distinguish between Request 1 and Request 2 (both result in `IsNull() == true`)

**With `opt`:** Can distinguish all three scenarios:
```go
if req.Name.IsUnset() {
    // Don't update (Request 2)
} else if req.Name.IsNull() {
    // Set to NULL (Request 3)
} else {
    // Update with value (Request 1)
}
```

#### Feature Comparison

| Feature | `nullable` | `opt` |
|---------|-----------|-------|
| **State Model** | 2-state (null/value) | 3-state (unset/null/value) |
| **Zero Value** | `null` | `unset` |
| **Clean JSON** | ✅ | ✅ |
| **Database Operations** | ✅ | ✅ |
| **Partial Updates** | ❌ | ✅ |
| **Distinguish unset vs null** | ❌ | ✅ |
| **Type Constraints** | ✅ (safer) | ❌ (any type) |
| **PostgreSQL JSON/JSONB** | ✅ Optimized | ✅ Generic |
| **UUID Support** | ✅ Built-in | ✅ Any type |
| **Functional Operations** | ❌ | ✅ `Map()`, etc. |
| **Package Structure** | Single type | 3 sub-packages |
| **Maturity** | Stable | Pre-1.0 |

#### API Comparison

**Creating Values:**
```go
// nullable
name := nullable.FromValue("John")
email := nullable.Null[string]()

// opt
import "github.com/aarondl/opt/omitnull"
name := omitnull.From("John")
email := omitnull.FromNull[string]()
unset := omitnull.Val[string]{}  // unset state
```

**Checking State:**
```go
// nullable - 2 checks
if value.IsNull() {
    // null or zero value
}

// opt - 3 distinct checks
if value.IsUnset() {
    // field omitted
} else if value.IsNull() {
    // explicitly null
} else if value.IsValue() {
    // has value
}
```

**Getting Values:**
```go
// nullable
if !value.IsNull() {
    v := value.GetValue()  // *T
    fmt.Println(*v)
}

// opt - more options
v, ok := value.Get()           // (T, bool)
v := value.GetOr("default")    // with fallback
v := value.MustGet()           // panics if not set
ptr := value.Ptr()             // *T or nil
```

#### Real-World Example: PATCH Endpoint

```go
// With opt - perfect for partial updates
type UpdateUserRequest struct {
    Name  omitnull.Val[string] `json:"name"`
    Email omitnull.Val[string] `json:"email"`
    Age   omitnull.Val[int]    `json:"age"`
}

func UpdateUser(req UpdateUserRequest) {
    query := "UPDATE users SET "
    var updates []string
    var args []any

    if req.Name.IsValue() {
        updates = append(updates, "name = ?")
        args = append(args, req.Name.MustGet())
    } else if req.Name.IsNull() {
        updates = append(updates, "name = NULL")
    }
    // else: IsUnset() - don't touch this field

    // Same pattern for Email and Age...
}

// Handles all these requests correctly:
// {}                              → no updates
// {"name": "John"}                → update name only
// {"name": "John", "age": null}   → update name, clear age
// {"name": null, "age": null}     → clear both
```

#### When to Choose Each

**Choose `nullable` when:**
- Building typical CRUD applications
- You need clean JSON marshaling for database types
- Simpler 2-state model (null/value) fits your needs
- You want type safety with constraints
- You prefer a simpler API

**Choose `opt` when:**
- Building REST APIs with PATCH endpoints
- You need to distinguish "omitted" from "null"
- Implementing partial update semantics
- Working with GraphQL (handles optional/nullable distinction)
- You need support for any type (not just specific types)
- You want functional operations like `Map()`

Both packages solve the "clean JSON marshaling" problem well. The key difference is whether you need 2 states (null/value) or 3 states (unset/null/value).

## Similar Projects

This project was inspired by [gonull](https://github.com/lomsa-dev/gonull), which had issues with PostgreSQL types like `enum`, `timestamp`, and `json`/`jsonb`. This library addresses those limitations while providing a cleaner API.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See [LICENSE](LICENSE) file for details.
