// Code generated by sqlc. DO NOT EDIT.

package DB

import (
	"time"

	"github.com/google/uuid"
)

type Msg struct {
	ID        uuid.UUID
	Msg       string
	MsgID     string
	Event     string
	FkUser    uuid.UUID
	CreatedAt time.Time
}

type TwitchUser struct {
	ID          uuid.UUID
	Name        string
	Nick        string
	Accesslevel int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
