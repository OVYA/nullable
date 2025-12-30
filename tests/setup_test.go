package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ovya/nullable"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var now = time.Now().UTC()
var astring = "a string"
var name = "PLOP"
var aint = 42

// Package-level variables for shared test infrastructure
var (
	testDB        *sql.DB                      // Shared database connection
	testContainer *postgres.PostgresContainer  // Container reference for cleanup
)

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

// TestMain sets up the PostgreSQL container before tests and tears it down after
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Create PostgreSQL container with testcontainers
	container, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithInitScripts("init.sql"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	testContainer = container

	// Get connection string and connect to database
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to get connection string: %v", err)
	}

	testDB, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Verify connection
	if err := testDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✓ PostgreSQL container started and database connected")

	// Run tests
	code := m.Run()

	// Cleanup
	if testDB != nil {
		testDB.Close()
	}
	if testContainer != nil {
		if err := testContainer.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate container: %v", err)
		}
	}

	os.Exit(code)
}

// getDB returns the shared test database connection
func getDB(t *testing.T) *sql.DB {
	t.Helper()

	if testDB == nil {
		t.Fatal("testDB is nil - TestMain may not have run")
	}

	return testDB
}

// cleanupTables truncates specified tables to reset state between tests
func cleanupTables(t *testing.T, db *sql.DB, tables ...string) {
	t.Helper()

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		_, err := db.Exec(query)
		if err != nil {
			t.Fatalf("Failed to cleanup table %s: %v", table, err)
		}
	}
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
