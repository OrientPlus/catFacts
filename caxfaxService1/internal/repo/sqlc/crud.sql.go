// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: crud.sql

package sqlc

import (
	"context"
	"database/sql"
)

const addFact = `-- name: AddFact :one
INSERT INTO facts (message, length)
VALUES ($1, $2) RETURNING id
`

type AddFactParams struct {
	Message sql.NullString
	Length  sql.NullInt32
}

func (q *Queries) AddFact(ctx context.Context, arg AddFactParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, addFact, arg.Message, arg.Length)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteFactByID = `-- name: DeleteFactByID :exec
DELETE FROM facts
WHERE id = $1
`

func (q *Queries) DeleteFactByID(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteFactByID, id)
	return err
}

const deleteFactByMessage = `-- name: DeleteFactByMessage :exec
DELETE FROM facts
WHERE message = $1
`

func (q *Queries) DeleteFactByMessage(ctx context.Context, message sql.NullString) error {
	_, err := q.db.ExecContext(ctx, deleteFactByMessage, message)
	return err
}

const getFactByID = `-- name: GetFactByID :one

SELECT id, message, length
FROM facts
WHERE id = $1
`

// query.sql
func (q *Queries) GetFactByID(ctx context.Context, id int32) (Fact, error) {
	row := q.db.QueryRowContext(ctx, getFactByID, id)
	var i Fact
	err := row.Scan(&i.ID, &i.Message, &i.Length)
	return i, err
}

const getFactByMessage = `-- name: GetFactByMessage :one
SELECT id, message, length
FROM facts
WHERE message = $1
`

func (q *Queries) GetFactByMessage(ctx context.Context, message sql.NullString) (Fact, error) {
	row := q.db.QueryRowContext(ctx, getFactByMessage, message)
	var i Fact
	err := row.Scan(&i.ID, &i.Message, &i.Length)
	return i, err
}

const updateFactByID = `-- name: UpdateFactByID :exec
UPDATE facts
SET message = $1,
    length = $2
WHERE id = $3
`

type UpdateFactByIDParams struct {
	Message sql.NullString
	Length  sql.NullInt32
	ID      int32
}

func (q *Queries) UpdateFactByID(ctx context.Context, arg UpdateFactByIDParams) error {
	_, err := q.db.ExecContext(ctx, updateFactByID, arg.Message, arg.Length, arg.ID)
	return err
}