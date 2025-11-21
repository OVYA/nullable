package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ovya/nullable"
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
		IntVal:    nullable.FromValue(42),
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
	err := db.QueryRow(`
		INSERT INTO type_test (
			string_val, int_val, int16_val, int32_val, int64_val,
			float_val, bool_val, uuid_val, time_val, json_val
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		&typeTest.StringVal, &typeTest.IntVal, &typeTest.Int16Val, &typeTest.Int32Val,
		&typeTest.Int64Val, &typeTest.FloatVal, &typeTest.BoolVal, &typeTest.UUIDVal,
		&typeTest.TimeVal, &typeTest.JSONVal,
	).Scan(&insertedID)

	if err != nil {
		t.Fatalf("Insert types failed: %v", err)
	}

	t.Logf("✓ Inserted all types with ID: %d", insertedID)

	// Read back
	var readTest TypeTest

	err = db.QueryRow(`
		SELECT id, string_val, int_val, int16_val, int32_val, int64_val,
			float_val, bool_val, uuid_val, time_val, json_val
		FROM type_test WHERE id = $1`,
		insertedID,
	).Scan(
		&readTest.ID, &readTest.StringVal, &readTest.IntVal, &readTest.Int16Val,
		&readTest.Int32Val, &readTest.Int64Val, &readTest.FloatVal, &readTest.BoolVal,
		&readTest.UUIDVal, &readTest.TimeVal, &readTest.JSONVal,
	)

	if err != nil {
		t.Fatalf("Read types failed: %v", err)
	}

	t.Log("✓ All types read successfully")

	// Verify values
	if *readTest.StringVal.GetValue() != "test string" {
		t.Errorf("String mismatch: expected 'test string', got '%s'", *readTest.StringVal.GetValue())
	}
	if *readTest.IntVal.GetValue() != 42 {
		t.Errorf("Int mismatch: expected 42, got %d", *readTest.IntVal.GetValue())
	}
	if *readTest.Int16Val.GetValue() != 16 {
		t.Errorf("Int16 mismatch: expected 16, got %d", *readTest.Int16Val.GetValue())
	}
	if *readTest.Int32Val.GetValue() != 32 {
		t.Errorf("Int32 mismatch: expected 32, got %d", *readTest.Int32Val.GetValue())
	}
	if *readTest.Int64Val.GetValue() != 64 {
		t.Errorf("Int64 mismatch: expected 64, got %d", *readTest.Int64Val.GetValue())
	}
	if *readTest.FloatVal.GetValue() != 3.14 {
		t.Errorf("Float mismatch: expected 3.14, got %f", *readTest.FloatVal.GetValue())
	}
	if !*readTest.BoolVal.GetValue() {
		t.Error("Bool should be true")
	}
	if *readTest.UUIDVal.GetValue() != testUUID {
		t.Errorf("UUID mismatch: expected %s, got %s", testUUID, *readTest.UUIDVal.GetValue())
	}

	// Time comparison (within 1 second tolerance for timezone handling)
	timeDiff := readTest.TimeVal.GetValue().Sub(testTime)
	if timeDiff < -time.Second || timeDiff > time.Second {
		t.Errorf("Time mismatch: expected %v, got %v (diff: %v)", testTime, *readTest.TimeVal.GetValue(), timeDiff)
	}

	// JSON value check
	if readTest.JSONVal.IsNull() {
		t.Error("JSON value should not be null")
	}

	t.Log("✓ Type verification passed")
}

func TestReadExisting(t *testing.T) {
	db := getDB(t)
	defer db.Close()

	rows, err := db.Query("SELECT id, name, date_to, data FROM test WHERE name IS NOT NULL ORDER BY id")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	var tests []testedStruct[embeddedStruct]
	for rows.Next() {
		var test testedStruct[embeddedStruct]

		err := rows.Scan(&test.ID, &test.Name, &test.DateTo, &test.Data)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}

		if test.Data.IsNull() {
			t.Fatal("Data should not be null")
		} else {
			t.Log("Data is not null => ok")
			if !strings.Contains(test.Data.GetValue().String, "value ") {
				t.Fatalf("data.String should contains 'value ', get: '%s'", test.Data.GetValue().String)
			} else {
				t.Log("data.String contains 'value ' => OK")
			}

			if test.Data.GetValue().Int != 42 {
				t.Fatalf("Bad value for data.Int, get: %d", test.Data.GetValue().Int)
			} else {
				t.Log("data.Int == 42 => OK")
			}

			if test.Data.GetValue().Bool.IsNull() {
				t.Fatal("data.Bool should not be null")
			} else {
				if !*test.Data.GetValue().Bool.GetValue() {
					t.Fatal("data.Bool should be true")
				} else {
					t.Log("data.Bool is true => OK")
				}
			}
		}

		if test.Name.IsNull() {
			t.Fatal("Name should not be null")
		} else {
			t.Log("data.Name is not null => OK")
		}

		tests = append(tests, test)
	}

	if len(tests) == 0 {
		t.Fatal("Expected at least one record from init.sql, got none")
	}

	t.Logf("✓ Found %d records", len(tests))

	for i, test := range tests {
		b, _ := json.MarshalIndent(test, "", "  ")
		t.Logf("Record %d: %s", i+1, b)
	}

	rows, err = db.Query("SELECT id, name, date_to, data FROM test WHERE name IS NULL ORDER BY id")
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var test testedStruct[embeddedStruct]

		err := rows.Scan(&test.ID, &test.Name, &test.DateTo, &test.Data)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}

		if !test.Name.IsNull() {
			t.Fatal("Name should be null")
		} else {
			t.Log("Null check for Name => OK")
		}

		if !test.DateTo.IsNull() {
			t.Fatal("DateTo should be null")
		} else {
			t.Log("Null check for DateTo => OK")
		}
	}

}

func TestInsertAndRead(t *testing.T) {
	db := getDB(t)
	defer db.Close()

	name := "plop name"

	// Create test data
	data := getEmbeddedObj()

	test := testedStruct[embeddedStruct]{
		Name:   nullable.FromValue(name),
		DateTo: nullable.FromValue(now),
		Data:   nullable.FromValue(data),
	}

	// Insert
	var insertedID int64
	err := db.QueryRow(
		"INSERT INTO test (name, date_to, data) VALUES ($1, $2, $3) RETURNING id",
		&test.Name, &test.DateTo, &test.Data,
	).Scan(&insertedID)

	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	t.Logf("✓ Inserted record with ID: %d", insertedID)

	// Read back
	var readTest testedStruct[embeddedStruct]

	err = db.QueryRow(
		"SELECT id, name, date_to, data FROM test WHERE id = $1",
		insertedID,
	).Scan(&readTest.ID, &readTest.Name, &readTest.DateTo, &readTest.Data)

	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	b, _ := json.MarshalIndent(readTest, "", "  ")
	t.Logf("✓ Read back record:\n%s", b)

	// Verify data
	if readTest.Name.GetValue() == nil {
		t.Fatal("Name should not be null")
	}
	if *readTest.Name.GetValue() != name {
		t.Fatalf("Name mismatch: expected '%s', got '%s'", name, *readTest.Name.GetValue())
	}

	if readTest.DateTo.IsNull() {
		t.Fatal("DateTo should not be null")
	}

	if readTest.DateTo.GetValue().Sub(now) > time.Nanosecond {
		t.Fatalf("DateTo mismatch: expected '%s', got '%s'", now, *readTest.DateTo.GetValue())
	}

	if readTest.Data.GetValue() == nil {
		t.Fatal("Data should not be null")
	}
	if readTest.Data.GetValue().String != astring {
		t.Fatalf("Data.String mismatch: expected '%s', got '%s'", astring, readTest.Data.GetValue().String)
	}

	if readTest.Data.GetValue().Bool.GetValue() == nil {
		t.Fatal("Data.Bool should not be null")
	}
	if !*readTest.Data.GetValue().Bool.GetValue() {
		t.Fatal("Data.Bool should be true")
	}

	t.Log("✓ Data verification passed")
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
	err := db.QueryRow(
		"INSERT INTO test (name, date_to, data) VALUES ($1, $2, $3) RETURNING id",
		&test.Name, &test.DateTo, &test.Data,
	).Scan(&insertedID)

	if err != nil {
		t.Fatalf("Insert NULL failed: %v", err)
	}

	t.Logf("✓ Inserted NULL record with ID: %d", insertedID)

	// Read back
	var readTest testedStruct[nullable.JSON]

	err = db.QueryRow(
		"SELECT id, name, date_to, data FROM test WHERE id = $1",
		insertedID,
	).Scan(&readTest.ID, &readTest.Name, &readTest.DateTo, &readTest.Data)

	if err != nil {
		t.Fatalf("Read NULL failed: %v", err)
	}

	// Verify all are null
	if !readTest.Name.IsNull() {
		t.Fatal("Name should be NULL")
	}
	if !readTest.DateTo.IsNull() {
		t.Fatal("DateTo should be NULL")
	}
	if !readTest.Data.IsNull() {
		t.Fatal("Data should be NULL")
	}

	t.Log("✓ NULL values correctly preserved")
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
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	ctx := context.Background()
	name := "plop name sqlx"

	// Create test data
	data := getEmbeddedObj()

	test := testedStruct[embeddedStruct]{
		Name:   nullable.FromValue(name),
		DateTo: nullable.FromValue(now),
		Data:   nullable.FromValue(data),
	}

	// Insert using BindNamed
	query, args, err := db.BindNamed("INSERT INTO test (name, date_to, data) VALUES (:name, :date_to, :data) RETURNING id",
		test,
	)
	if err != nil {
		t.Fatalf("BindNamed failed: %v", err)
	}

	var insertedID int64
	err = db.QueryRowContext(ctx, query, args...).Scan(&insertedID)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	t.Logf("✓ Inserted record with ID: %d using BindNamed", insertedID)

	// Read back using GetContext
	var readTest testedStruct[embeddedStruct]

	err = sqlx.GetContext(ctx, db,
		&readTest,
		"SELECT id, name, date_to, data FROM test WHERE id = $1",
		insertedID,
	)
	if err != nil {
		t.Fatalf("GetContext failed: %v", err)
	}

	b, _ := json.MarshalIndent(readTest, "", "  ")
	t.Logf("✓ Read back record with GetContext:\n%s", b)

	// Verify data
	if readTest.Name.GetValue() == nil {
		t.Fatal("Name should not be null")
	}
	if *readTest.Name.GetValue() != name {
		t.Fatalf("Name mismatch: expected '%s', got '%s'", name, *readTest.Name.GetValue())
	}

	if readTest.DateTo.IsNull() {
		t.Fatal("DateTo should not be null")
	}

	if readTest.DateTo.GetValue().Sub(now) > time.Nanosecond {
		t.Fatalf("DateTo mismatch: expected '%s', got '%s'", now, *readTest.DateTo.GetValue())
	}

	if readTest.Data.GetValue() == nil {
		t.Fatal("Data should not be null")
	}
	if readTest.Data.GetValue().String != astring {
		t.Fatalf("Data.String mismatch: expected '%s', got '%s'", astring, readTest.Data.GetValue().String)
	}

	if readTest.Data.GetValue().Bool.GetValue() == nil {
		t.Fatal("Data.Bool should not be null")
	}
	if !*readTest.Data.GetValue().Bool.GetValue() {
		t.Fatal("Data.Bool should be true")
	}

	t.Log("✓ Data verification passed with sqlx")
}
