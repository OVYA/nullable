package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ovya/nullable"
)

var now = time.Now()
var astring = "a string"
var name = "PLOP"
var aint = 42

type embeddedStruct struct {
	ID     int64                      `json:"id" db:"id"`
	String string                     `json:"string" db:"string"`
	Int    int                        `json:"int" db:"int"`
	Bool   nullable.Of[bool]          `json:"bool" db:"bool"`
	DateTo nullable.Of[time.Time]     `json:"dateTo" db:"date_to"`
	JSON   nullable.Of[nullable.JSON] `json:"json" db:"json"`
}

type testedStruct[T nullable.JSON] struct {
	ID     int64                  `json:"id" db:"id"`
	Name   nullable.Of[string]    `json:"name" db:"name"`
	DateTo nullable.Of[time.Time] `json:"dateTo" db:"date_to"`
	Data   nullable.Of[T]         `json:"data" db:"data"`
}

// Setup
// func TestMain(m *testing.M) {
// 	code := m.Run()
// 	os.Exit(code)
// }

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDB creates a database connection for testing
func getDB(t *testing.T) *sql.DB {
	t.Helper()

	// Build connection string from environment variables
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5445"),
		getEnv("DB_USER", "testuser"),
		getEnv("DB_PASSWORD", "testpass"),
		getEnv("DB_NAME", "testdb"),
		getEnv("DB_SSLMODE", "disable"),
	)

	// Connect to database
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	t.Log("✓ Connected to PostgreSQL successfully")
	return db
}

func getEmbeddedObj() embeddedStruct {
	obj := embeddedStruct{
		String: astring,
		Int:    aint,
		Bool:   nullable.FromValue(true),
		DateTo: nullable.FromValue(now),
	}

	obj.JSON = nullable.FromValue[nullable.JSON](obj)

	return obj
}

func getNullTestObj[T nullable.JSON]() testedStruct[T] {
	return testedStruct[T]{
		ID:     1,
		Name:   nullable.Null[string](),    // Null value
		DateTo: nullable.Null[time.Time](), // Null value
		Data:   nullable.Null[T](),         // Null value
	}
}

func getTestObjs[T nullable.JSON](data T) []testedStruct[T] {
	obj1 := testedStruct[T]{
		Name:   nullable.FromValue(name),
		DateTo: nullable.FromValue(now),
		Data:   nullable.FromValue(data),
	}

	obj2 := getNullTestObj[T]()

	return []testedStruct[T]{obj1, obj2}
}
