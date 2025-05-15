-- CreateTable
CREATE TABLE
    tasks (
        id BIGINT GENERATED ALWAYS AS IDENTITY,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        args JSONB,
        idempotency_key TEXT,
        PRIMARY KEY (id)
    );

CREATE UNIQUE INDEX
    tasks_idempotency_key_idx ON tasks (idempotency_key);

-- CreateTable
CREATE TABLE
    task_associated_data (
        task_id BIGINT NOT NULL,
        top_level_fields TEXT[],
        PRIMARY KEY (task_id)
    );

CREATE OR REPLACE FUNCTION extract_top_level_fields(args_json JSONB)
RETURNS TEXT[] AS $$
DECLARE
    result TEXT[];
BEGIN
    -- Extract all top-level keys from the JSONB object and convert them to an array
    SELECT array_agg(key)
    INTO result
    FROM jsonb_object_keys(args_json) AS key;
    
    -- Handle the case when args_json is null or empty
    IF result IS NULL THEN
        result := '{}'; -- Empty array
    END IF;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;