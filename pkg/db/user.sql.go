// Code generated by sqlc. DO NOT EDIT.
// source: user.sql

package DB

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const getUserUUIDByNick = `-- name: GetUserUUIDByNick :one
SELECT id FROM twitch_user WHERE nick=$1
`

func (q *Queries) GetUserUUIDByNick(ctx context.Context, nick string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getUserUUIDByNick, nick)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const insertUser = `-- name: InsertUser :one
INSERT INTO twitch_user
(nick,name, accesslevel,created_at,updated_at)
VALUES ($1,$2,$3,$4,$5) RETURNING id
`

type InsertUserParams struct {
	Nick        string
	Name        string
	Accesslevel int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, insertUser,
		arg.Nick,
		arg.Name,
		arg.Accesslevel,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}