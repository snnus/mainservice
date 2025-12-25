CREATE SCHEMA shard_1;
CREATE SCHEMA shard_2;
CREATE SCHEMA shard_3;
CREATE SCHEMA shard_4;

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE shard_1.service_points (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    short_name VARCHAR(10) NOT NULL,
    office_number VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_service_points_updated_at_shard_1 BEFORE UPDATE
    ON shard_1.service_points FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE shard_2.service_points (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    short_name VARCHAR(10) NOT NULL,
    office_number VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_service_points_updated_at_shard_2 BEFORE UPDATE
    ON shard_2.service_points FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE shard_3.service_points (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    short_name VARCHAR(10) NOT NULL,
    office_number VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_service_points_updated_at_shard_3 BEFORE UPDATE
    ON shard_3.service_points FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE shard_4.service_points (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    short_name VARCHAR(10) NOT NULL,
    office_number VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_service_points_updated_at_shard_4 BEFORE UPDATE
    ON shard_4.service_points FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();