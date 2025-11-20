-- Create test table
CREATE TABLE IF NOT EXISTS test (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    date_to TIMESTAMP,
    data JSONB
);

-- Create additional test table for various types
CREATE TABLE IF NOT EXISTS type_test (
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

-- Insert some test data
INSERT INTO test (name, date_to, data) VALUES
    ('Test 1', NOW(), '{"string": "value 1", "bool": true, "int": 42}'::jsonb),
    (NULL, NULL, NULL),
    ('Test 3', '2025-12-31 23:59:59', '{"string": "value 2", "bool": true, "int": 42}'::jsonb);

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO testuser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO testuser;
