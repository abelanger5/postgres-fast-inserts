-- name: InsertTaskSingleton :one
INSERT INTO tasks (args, idempotency_key)
VALUES ($1, $2)
RETURNING *;

-- name: InsertTasksBatch :batchone
INSERT INTO tasks (args, idempotency_key)
VALUES ($1, $2)
RETURNING *;

-- name: InsertTasksWithUnnest :many
WITH input AS (
    SELECT
        UNNEST(@args::JSONB[]) AS args,
        UNNEST(@keys::TEXT[]) AS idempotency_key
)
INSERT INTO tasks (args, idempotency_key)
SELECT
    args,
    idempotency_key
FROM input
RETURNING *;

-- name: InsertTasksCopyFrom :copyfrom
INSERT INTO tasks (args, idempotency_key) VALUES ($1, $2);