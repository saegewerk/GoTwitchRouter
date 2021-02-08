package twitchrouter

import (
	"github.com/google/uuid"
)

type App struct {
	CmdRegister CmdRegister //this only get's written once
	Queue *chan CmdMessage
	Closed *bool
}
type CmdMessage struct{
	Cmd string
	Event string
	Msg string
	MsgId string
	Uuid uuid.UUID
}
type CmdRegister struct{
	Cmd string
	Help string
	AccessLevel int32
}

type AppsQueued struct{
	Queue chan App
	Apps []App //remove apps that close transport
}
