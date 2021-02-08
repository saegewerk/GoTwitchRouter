package main

import (
	"github.com/saegewerk/BotZtwitch/pkg/twitchrouter"
)

func onMessage(in *twitchrouter.Message, send func(request *twitchrouter.MessageRequest) error){
	response:="a test response"
	help:=""
	accessLevel:=int32(1)
	err:=send(&twitchrouter.MessageRequest{
		Command: in.Command,
		Help:    &help,
		Msg:     &response,
		Uuid:    in.Uuid,
		AccessLevel: &accessLevel,
	})
	if err!=nil{
		println(err.Error())
	}
}

func main(){
	twitchrouter.Client("cmd","this is the help",3,onMessage)
}
