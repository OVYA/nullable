#!/bin/sh
set -e

# PostgreSQL data directory
PGDATA="/var/lib/postgresql/data"
export PGPORT=5445

echo "==> Starting PostgreSQL setup..."

# Initialize PostgreSQL if needed
if [ ! -f "$PGDATA/PG_VERSION" ]; then
    echo "==> Initializing PostgreSQL database..."
    su-exec postgres initdb -D "$PGDATA"

    # Configure PostgreSQL for local connections and custom port
    echo "host all all 0.0.0.0/0 trust" >>"$PGDATA/pg_hba.conf"
    echo "listen_addresses='*'" >>"$PGDATA/postgresql.conf"
    echo "port=5445" >>"$PGDATA/postgresql.conf"
fi

# Start PostgreSQL in background
echo "==> Starting PostgreSQL server..."
su-exec postgres pg_ctl -D "$PGDATA" -w start

# Wait for PostgreSQL to be ready
echo "==> Waiting for PostgreSQL to be ready..."
until su-exec postgres pg_isready -q; do
    echo "PostgreSQL is not ready yet, waiting..."
    sleep 1
done

echo "==> PostgreSQL is ready!"

# Create database and user if they don't exist
echo "==> Setting up database..."
su-exec postgres psql -v ON_ERROR_STOP=1 <<-EOSQL
    SELECT 'CREATE DATABASE testdb' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'testdb')\gexec
    SELECT 'CREATE USER testuser WITH PASSWORD ''testpass''' WHERE NOT EXISTS (SELECT FROM pg_user WHERE usename = 'testuser')\gexec
    GRANT ALL PRIVILEGES ON DATABASE testdb TO testuser;
EOSQL

# Run initialization SQL
if [ -f /app/tests/init.sql ]; then
    echo "==> Running initialization SQL..."
    su-exec postgres psql -d testdb -f /app/tests/init.sql
fi

# Set environment variables for the Go application
export DB_HOST=localhost
export DB_PORT=5445
export DB_USER=testuser
export DB_PASSWORD=testpass
export DB_NAME=testdb
export DB_SSLMODE=disable

echo "==> Setup complete! Running application..."
echo ""

# Execute the command passed to docker run
exec "$@"
