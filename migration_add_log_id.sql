-- Migration to add log_id column for cursor-based pagination
-- This improves performance for log retrieval by using cursor-based pagination instead of offset

-- Add log_id column as BIGSERIAL (auto-incrementing)
ALTER TABLE log_statements ADD COLUMN log_id BIGSERIAL;

-- Create index on log_id for efficient cursor queries
CREATE INDEX idx_log_statements_log_id ON log_statements(log_id);

-- Create composite index for deployment_id and log_id for optimal query performance
CREATE INDEX idx_log_statements_deployment_log_id ON log_statements(deployment_id, log_id DESC);