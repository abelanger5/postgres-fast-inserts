// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: copyfrom.go

package dbsqlc

import (
	"context"
)

// iteratorForInsertTasksCopyFrom implements pgx.CopyFromSource.
type iteratorForInsertTasksCopyFrom struct {
	rows                 []InsertTasksCopyFromParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertTasksCopyFrom) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertTasksCopyFrom) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Args,
		r.rows[0].IdempotencyKey,
	}, nil
}

func (r iteratorForInsertTasksCopyFrom) Err() error {
	return nil
}

func (q *Queries) InsertTasksCopyFrom(ctx context.Context, db DBTX, arg []InsertTasksCopyFromParams) (int64, error) {
	return db.CopyFrom(ctx, []string{"tasks"}, []string{"args", "idempotency_key"}, &iteratorForInsertTasksCopyFrom{rows: arg})
}
