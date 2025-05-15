-- name: InsertTaskAssociatedData :exec
INSERT INTO task_associated_data (task_id, top_level_fields)
VALUES ($1, extract_top_level_fields($2));

-- name: InsertTaskAssociatedDatasBatch :batchone
INSERT INTO task_associated_data (task_id, top_level_fields)
VALUES ($1, extract_top_level_fields($2))
RETURNING *;