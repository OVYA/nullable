# PostgreSQL Integration Tests

This directory contains integration tests for the nullable library with PostgreSQL.

## Test Files

- **postgres_test.go** - PostgreSQL integration tests

## Running Tests

### With Docker (Recommended)

From the project root:

```bash
# Run all tests (including database tests)
make test

# Or manually
docker build -t nullable-postgres-test -f tests/Dockerfile .
docker run --rm nullable-postgres-test
```

### Locally (Requires PostgreSQL)

1. Start PostgreSQL and create the test database:
```bash
createdb testdb
psql testdb < init.sql
```

2. Set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5445
export DB_USER=testuser
export DB_PASSWORD=testpass
export DB_NAME=testdb
export DB_SSLMODE=disable
```

3. Run the tests:
```bash
# Run all tests
go test -v ./tests

# Run specific test
go test -v ./tests -run TestInsertAndRead

# Run with short flag to skip slow tests (if implemented)
go test -v -short ./tests
```

## Test Coverage

The test suite includes:

### TestReadExisting
Verifies reading pre-seeded data from the database, ensuring:
- Database connection works
- Nullable fields are properly scanned
- JSON data is correctly deserialized

### TestInsertAndRead
Tests the full write-read cycle:
- Insert records with nullable values
- Read them back
- Verify data integrity
- Check nested nullable fields

### TestNullValues
Ensures proper NULL handling:
- Insert NULL values
- Verify NULL preservation
- Check IsNull() works correctly

### TestAllTypes
Comprehensive type coverage test:
- string, int, int16, int32, int64
- float64, bool
- uuid.UUID, time.Time
- JSON/JSONB for complex types

### TestNullableEdgeCases
Edge cases and special scenarios:
- SetValueP with nil pointer
- SetValueP with value pointer
- IsNull on zero values

## Database Schema

Tests expect the following tables (created by `init.sql`):

```sql
CREATE TABLE test (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    date_to TIMESTAMP,
    data JSONB
);

CREATE TABLE type_test (
    id SERIAL PRIMARY KEY,
    string_val VARCHAR(255),
    int_val INTEGER,
    int16_val SMALLINT,
    int32_val INTEGER,
    int64_val BIGINT,
    float_val DOUBLE PRECISION,
    bool_val BOOLEAN,
    uuid_val UUID,
    time_val TIMESTAMP,
    json_val JSONB
);
```

## Environment Variables

The tests use these environment variables (with defaults):

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5445 - non-standard to avoid conflicts)
- `DB_USER` - Database user (default: testuser)
- `DB_PASSWORD` - Database password (default: testpass)
- `DB_NAME` - Database name (default: testdb)
- `DB_SSLMODE` - SSL mode (default: disable)

## Tips

1. Tests require a running PostgreSQL instance
2. Use Docker setup for isolated testing
3. Each test is independent and uses transactions where possible
4. Tests create their own data (don't rely on specific database state)
5. The `getDB(t)` helper manages connections and logs connectivity

## CI/CD Integration

These tests are designed to run in CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run PostgreSQL integration tests
  run: make docker-test
```

The Docker setup ensures:
- Consistent test environment
- No external dependencies
- Fast startup (~5-10 seconds)
- Automatic cleanup
