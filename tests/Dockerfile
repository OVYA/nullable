FROM golang:1.24-alpine

# Install PostgreSQL server and client
RUN apk add --no-cache \
    postgresql \
    postgresql-contrib \
    su-exec

# Set working directory
WORKDIR /app

# Install gotestsum for better test output
RUN go install gotest.tools/gotestsum@latest

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Create directories for PostgreSQL
RUN mkdir -p /var/lib/postgresql/data /run/postgresql && \
    chown -R postgres:postgres /var/lib/postgresql /run/postgresql

# Copy entrypoint script
COPY tests/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# Expose PostgreSQL port (for debugging)
EXPOSE 5445

# Run the entrypoint script
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["gotestsum", "--format", "testdox", "--", "-v", "./tests/..."]
