// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"database/sql"
)

type Fact struct {
	ID        int32
	Message   sql.NullString
	Length    sql.NullInt32
	TimePoint sql.NullTime
}
