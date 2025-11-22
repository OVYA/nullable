package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ovya/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TypeTest represents all supported nullable types
type TypeTest struct {
	ID        int64                      `db:"id"`
	StringVal nullable.Of[string]        `db:"string_val"`
	IntVal    nullable.Of[int]           `db:"int_val"`
	Int16Val  nullable.Of[int16]         `db:"int16_val"`
	Int32Val  nullable.Of[int32]         `db:"int32_val"`
	Int64Val  nullable.Of[int64]         `db:"int64_val"`
	FloatVal  nullable.Of[float64]       `db:"float_val"`
	BoolVal   nullable.Of[bool]          `db:"bool_val"`
	UUIDVal   nullable.Of[uuid.UUID]     `db:"uuid_val"`
	TimeVal   nullable.Of[time.Time]     `db:"time_val"`
	JSONVal   nullable.Of[nullable.JSON] `db:"json_val"`
}

func TestAllTypes(t *testing.T) {
	db := getDB(t)
	defer db.Close()

	testUUID := uuid.New()
	testTime := time.Now().UTC().Truncate(time.Second)

	typeTest := TypeTest{
		StringVal: nullable.FromValue("test string"),
		IntVal:    nullable.FromValue(aint),
		Int16Val:  nullable.FromValue(int16(16)),
		Int32Val:  nullable.FromValue(int32(32)),
		Int64Val:  nullable.FromValue(int64(64)),
		FloatVal:  nullable.FromValue(3.14),
		BoolVal:   nullable.FromValue(true),
		UUIDVal:   nullable.FromValue(testUUID),
		TimeVal:   nullable.FromValue(testTime),
		JSONVal:   nullable.FromValue[nullable.JSON](map[string]any{"key": "value"}),
	}

	// Insert
	var insertedID int64
	t.Run("Writing data into database", func(t *testing.T) {
		err := db.QueryRow(`
		INSERT INTO type_test (
			string_val, int_val, int16_val, int32_val, int64_val,
			float_val, bool_val, uuid_val, time_val, json_val
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
			&typeTest.StringVal, &typeTest.IntVal, &typeTest.Int16Val, &typeTest.Int32Val,
			&typeTest.Int64Val, &typeTest.FloatVal, &typeTest.BoolVal, &typeTest.UUIDVal,
			&typeTest.TimeVal, &typeTest.JSONVal,
		).Scan(&insertedID)
		require.NoError(t, err, "Insert types failed")

		t.Logf("Inserted all types with ID: %d", insertedID)
	})

	// Read back
	var readTest TypeTest
	t.Run("Reading data from database", func(t *testing.T) {
		err := db.QueryRow(`
		SELECT id, string_val, int_val, int16_val, int32_val, int64_val,
			float_val, bool_val, uuid_val, time_val, json_val
		FROM type_test WHERE id = $1`,
			insertedID,
		).Scan(
			&readTest.ID, &readTest.StringVal, &readTest.IntVal, &readTest.Int16Val,
			&readTest.Int32Val, &readTest.Int64Val, &readTest.FloatVal, &readTest.BoolVal,
			&readTest.UUIDVal, &readTest.TimeVal, &readTest.JSONVal,
		)

		require.NoError(t, err, "Read types failed")
	})

	t.Run("data maching read <-> write", func(t *testing.T) {
		assert.Equal(t, "test string", *readTest.StringVal.GetValue())
		assert.Equal(t, aint, *readTest.IntVal.GetValue())
		assert.Equal(t, int16(16), *readTest.Int16Val.GetValue())
		assert.Equal(t, int32(32), *readTest.Int32Val.GetValue())
		assert.Equal(t, int64(64), *readTest.Int64Val.GetValue())
		assert.Equal(t, 3.14, *readTest.FloatVal.GetValue())
		assert.True(t, *readTest.BoolVal.GetValue())
		assert.Equal(t, testUUID, *readTest.UUIDVal.GetValue())

		// Time comparison (within 1 second tolerance for timezone handling)
		assert.WithinDuration(t, testTime, *readTest.TimeVal.GetValue(), time.Second)

		// JSON value check
		assert.False(t, readTest.JSONVal.IsNull(), "JSON value should not be null")
	})
}

func TestReadExisting(t *testing.T) {
	db := getDB(t)
	defer db.Close()

	t.Run("Reading existing not null data", func(t *testing.T) {
		rows, err := db.Query("SELECT id, name, date_to, data FROM test WHERE name IS NOT NULL ORDER BY id")
		defer rows.Close()
		t.Run("Querying", func(t *testing.T) {
			require.NoError(t, err, "Query failed")
		})

		var tests []testedStruct[embeddedStruct]
		for rows.Next() {
			var test testedStruct[embeddedStruct]

			t.Run("Scanning", func(t *testing.T) {
				err := rows.Scan(&test.ID, &test.Name, &test.DateTo, &test.Data)
				require.NoError(t, err, "Scan failed")
			})

			t.Run(fmt.Sprintf("Scanned data matching for id %d", test.ID), func(t *testing.T) {
				require.False(t, test.Data.IsNull(), "Data should not be null")
				assert.Contains(t, test.Data.GetValue().String, "value ")
				assert.Equal(t, aint, test.Data.GetValue().Int)

				require.False(t, test.Data.GetValue().Bool.IsNull(), "data.Bool should not be null")
				assert.True(t, *test.Data.GetValue().Bool.GetValue(), "data.Bool should be true")

				require.False(t, test.Name.IsNull(), "Name should not be null")
			})

			tests = append(tests, test)
		}

		require.NotEmpty(t, tests, "Expected at least one record from init.sql")

		t.Logf("Found %d records", len(tests))

		for i, test := range tests {
			b, _ := json.MarshalIndent(test, "", "  ")
			t.Logf("Record %d: %s", i+1, b)
		}
	})

	t.Run("Reading existing null data", func(t *testing.T) {
		rows, err := db.Query("SELECT id, name, date_to, data FROM test WHERE name IS NULL ORDER BY id")
		t.Run("Querying", func(t *testing.T) {
			require.NoError(t, err, "Query failed")
		})
		defer rows.Close()

		for rows.Next() {
			var test testedStruct[embeddedStruct]

			t.Run("Scanning", func(t *testing.T) {
				err := rows.Scan(&test.ID, &test.Name, &test.DateTo, &test.Data)
				require.NoError(t, err, "Scan failed")
			})

			t.Run(fmt.Sprintf("Scanned data matching for id %d", test.ID), func(t *testing.T) {
				assert.True(t, test.Name.IsNull(), "Name should be null")
				assert.True(t, test.DateTo.IsNull(), "DateTo should be null")
			})
		}
	})
}

func TestNullValues(t *testing.T) {
	db := getDB(t)
	defer db.Close()

	test := testedStruct[embeddedStruct]{
		Name:   nullable.Null[string](),
		DateTo: nullable.Null[time.Time](),
		Data:   nullable.Null[embeddedStruct](),
	}

	var insertedID int64

	t.Run("Inserting", func(t *testing.T) {
		err := db.QueryRow(
			"INSERT INTO test (name, date_to, data) VALUES ($1, $2, $3) RETURNING id",
			&test.Name, &test.DateTo, &test.Data,
		).Scan(&insertedID)
		require.NoError(t, err, "Insert NULL failed")

		t.Logf("Inserted NULL record with ID: %d", insertedID)
	})

	// Read back
	var readTest testedStruct[nullable.JSON]

	t.Run("Reading and scanning", func(t *testing.T) {
		err := db.QueryRow(
			"SELECT id, name, date_to, data FROM test WHERE id = $1",
			insertedID,
		).Scan(&readTest.ID, &readTest.Name, &readTest.DateTo, &readTest.Data)
		require.NoError(t, err, "Read NULL failed")
	})

	t.Run("Verify all values are null", func(t *testing.T) {
		assert.True(t, readTest.Name.IsNull(), "Name should be NULL")
		assert.True(t, readTest.DateTo.IsNull(), "DateTo should be NULL")
		assert.True(t, readTest.Data.IsNull(), "Data should be NULL")
	})
}

func TestInsertAndRead(t *testing.T) {
	db := getDB(t)
	defer db.Close()

	// Create test data
	data := getEmbeddedObj()

	test := testedStruct[embeddedStruct]{
		Name:   nullable.FromValue(name),
		DateTo: nullable.FromValue(now),
		Data:   nullable.FromValue(data),
	}

	var insertedID int64

	t.Run("Inserting", func(t *testing.T) {
		err := db.QueryRow(
			"INSERT INTO test (name, date_to, data) VALUES ($1, $2, $3) RETURNING id",
			&test.Name, &test.DateTo, &test.Data,
		).Scan(&insertedID)
		require.NoError(t, err, "Insert failed")
	})

	t.Logf("Inserted record with ID: %d", insertedID)

	var readTest testedStruct[embeddedStruct]

	t.Run("Reading back", func(t *testing.T) {
		err := db.QueryRow(
			"SELECT id, name, date_to, data FROM test WHERE id = $1",
			insertedID,
		).Scan(&readTest.ID, &readTest.Name, &readTest.DateTo, &readTest.Data)
		require.NoError(t, err, "Read failed")
	})

	b, _ := json.MarshalIndent(readTest, "", "  ")
	t.Logf("Read back record:\n%s", b)

	t.Run("Data matching", func(t *testing.T) {
		require.NotNil(t, readTest.Name.GetValue(), "Name should not be null")
		assert.Equal(t, name, *readTest.Name.GetValue())

		require.False(t, readTest.DateTo.IsNull(), "DateTo should not be null")
		assert.WithinDuration(t, now, *readTest.DateTo.GetValue(), time.Millisecond)

		require.NotNil(t, readTest.Data.GetValue(), "Data should not be null")
		assert.Equal(t, astring, readTest.Data.GetValue().String)

		require.NotNil(t, readTest.Data.GetValue().Bool.GetValue(), "Data.Bool should not be null")
		assert.True(t, *readTest.Data.GetValue().Bool.GetValue(), "Data.Bool should be true")

		require.NotNil(t, readTest.Data.GetValue().DateTo.GetValue(), "Data.DateTo should not be null")
		assert.WithinDuration(t, now, *readTest.Data.GetValue().DateTo.GetValue(), time.Millisecond)
	})
}

func TestInsertAndReadWithSqlx(t *testing.T) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5445"),
		getEnv("DB_USER", "testuser"),
		getEnv("DB_PASSWORD", "testpass"),
		getEnv("DB_NAME", "testdb"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := sqlx.Open("pgx", connStr)
	require.NoError(t, err, "Failed to connect to database")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Failed to ping database")

	ctx := context.Background()

	// Create test data
	data := getEmbeddedObj()

	test := testedStruct[embeddedStruct]{
		Name:   nullable.FromValue(name),
		DateTo: nullable.FromValue(now),
		Data:   nullable.FromValue(data),
	}

	var query string
	var args []any
	var insertedID int64

	t.Run("Insert using db.BindNamed ", func(t *testing.T) {
		query, args, err = db.BindNamed(
			"INSERT INTO test (name, date_to, data) VALUES (:name, :date_to, :data) RETURNING id",
			test,
		)
		require.NoError(t, err, "BindNamed failed")

		err = db.QueryRowContext(ctx, query, args...).Scan(&insertedID)
		require.NoError(t, err, "Insert failed")

		t.Logf("Inserted record with ID: %d using BindNamed", insertedID)
	})

	var readTest testedStruct[embeddedStruct]

	t.Run("Reading back using db.GetContext", func(t *testing.T) {

		err = db.GetContext(ctx,
			&readTest,
			"SELECT id, name, date_to, data FROM test WHERE id = $1",
			insertedID,
		)
		require.NoError(t, err, "GetContext failed")

		b, _ := json.MarshalIndent(readTest, "", "  ")
		t.Logf("Read back record with GetContext:\n%s", b)
	})

	t.Run("Data matching", func(t *testing.T) {
		require.NotNil(t, readTest.Name.GetValue(), "Name should not be null")
		assert.Equal(t, name, *readTest.Name.GetValue())

		require.False(t, readTest.DateTo.IsNull(), "DateTo should not be null")
		assert.WithinDuration(t, now, *readTest.DateTo.GetValue(), time.Millisecond)

		require.NotNil(t, readTest.Data.GetValue(), "Data should not be null")
		assert.Equal(t, astring, readTest.Data.GetValue().String)

		require.NotNil(t, readTest.Data.GetValue().Bool.GetValue(), "Data.Bool should not be null")
		assert.True(t, *readTest.Data.GetValue().Bool.GetValue(), "Data.Bool should be true")
	})
}
