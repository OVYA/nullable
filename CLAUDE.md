# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library (`github.com/ovya/nullable`) that provides generic nullable types for any data type, with special focus on database operations and JSON marshaling/unmarshaling. The library uses Go generics to wrap values in a nullable container (`Of[T]`) that can represent SQL NULL values while maintaining type safety.

### Core Architecture

**Main library files (root directory):**
- `nullable.go` - Core interface `NullableI[T]`, helper functions (`FromValue`, `Null`), and type-specific scanning methods (`scanJSON`, `scanString`, `scanUUID`, `scanInt`, `scanFloat`, `scanBool`, `scanTime`)
- `of.go` - Generic `Of[T]` struct implementation with methods for SQL scanning (`Scan`), SQL value conversion (`Value`), JSON marshaling/unmarshaling, and value management (getters/setters)
- `doc.go` - Package documentation

**Key design patterns:**
1. **Generic nullable wrapper**: `Of[T]` wraps any supported type `T` with an internal pointer `val *T` where `nil` represents NULL
2. **Type dispatch in Scan/Value**: The `Scan` and `Value` methods use type switches to route to specialized handlers for primitive types vs JSON types
3. **Custom type support**: Types implementing `sql.Scanner` or `driver.Valuer` interfaces are automatically supported without JSON marshaling
4. **Dual module structure**: Main module at root, separate test module in `tests/` directory with `replace` directive

### Supported Types

The library constrains type parameter `T` to: `bool | int | int16 | int32 | int64 | string | uuid.UUID | float64 | JSON`

- **Primitive types** are stored directly in the database
- **JSON type** (alias for `any`) is marshaled to JSON for database storage unless the type implements custom `sql.Scanner`/`driver.Valuer`
- Custom types can implement `sql.Scanner` and `driver.Valuer` to control their own serialization (see README.md example with `PhoneNumber`)

## Development Commands

### Running Tests

**Run all tests (including PostgreSQL integration tests):**
```bash
make test
```
This builds a Docker container with PostgreSQL and runs the complete test suite.

**Run only unit tests (no database):**
```bash
cd tests
go test -run TestMarshalUnmarshal -v
go test -run TestNullableEdgeCases -v
```

**Run specific database test:**
```bash
cd tests
# Requires Docker to be running
go test -run TestAllTypes -v
go test -run TestInsertAndRead -v
```

**Run a single test from the tests directory:**
```bash
cd tests
go test -run TestName -v
```

### Code Quality

**Lint the code:**
```bash
golangci-lint run
```

The project uses extensive linting (see `.golangci.yml`) with 30+ enabled linters including gosec, govet, errcheck, and revive.

**Tidy dependencies:**
```bash
go mod tidy
cd tests && go mod tidy
```

### Docker Development

**Build Docker test image:**
```bash
make docker-build
```

**Open shell in test container (for debugging):**
```bash
make docker-shell
```

**Run tests interactively:**
```bash
make docker-run
```

## Test Organization

The test suite is located in `tests/` directory with its own `go.mod` that uses a `replace` directive to reference the parent module.

**Test files:**
- `nullable_test.go` - Unit tests for JSON marshaling/unmarshaling and edge cases
- `postgres_test.go` - Integration tests with PostgreSQL database (requires Docker)
- `setup_test.go` - Test fixtures and database connection helpers

**Test coverage verification:**
```bash
cd tests
# Count test files
find . -name "*_test.go" | wc -l
# Should match the number of test files reported by `go test -v`
```

## Working with MarshalJSON/UnmarshalJSON

The `Of[T]` type implements custom JSON marshaling:

**MarshalJSON (of.go:67-73):**
- Returns `[]byte("null")` if value is null
- Otherwise delegates to generic `marshalJSON` helper (nullable.go:260-267)

**UnmarshalJSON (of.go:76-97):**
- Handles `null` JSON values by calling `SetNull()`
- For non-null values, unmarshals directly into the wrapped value
- Allocates new `T` if needed before unmarshaling

**Key invariant:** JSON `null` maps to Go `nil` pointer, not a special Valid/Invalid flag like `database/sql` types.

## Database Integration

The library integrates with `database/sql` through two interfaces:

1. **`driver.Valuer` (of.go:100-132)**: Converts Go values to database values
   - Primitive types return their dereferenced value
   - JSON types check for custom `driver.Valuer` first, then marshal to JSON string

2. **`sql.Scanner` (of.go:136-159)**: Converts database values to Go values
   - Routes to type-specific scan methods based on the wrapped type
   - Each scan method (in nullable.go) handles SQL NULL properly

## Go Version and Dependencies

- **Go version:** 1.24.10
- **Dependencies:**
  - `github.com/google/uuid` - UUID type support
  - Test dependencies: `pgx/v5`, `sqlx`, `testify`

## Common Gotchas

1. **Module structure**: Root module (`github.com/ovya/nullable`) and test module (`github.com/ovya/nullable/tests`) are separate. Always run `go mod tidy` in both directories after dependency changes.

2. **Test execution**: Integration tests require Docker. Use `make test` for full suite, or run unit tests directly with `go test` in the `tests/` directory.

3. **Type constraints**: The generic constraint limits supported types. Adding new primitive types requires updating the constraint in both `NullableI` interface and `Of[T]` struct definition.

4. **Time precision**: PostgreSQL tests truncate time to seconds (`Truncate(time.Second)`) to match database precision.
