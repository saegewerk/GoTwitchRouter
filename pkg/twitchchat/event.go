package twitchchat

import irc "github.com/fluffle/goirc/client"

type (
	Connected    func()
	Disconnected func()
	NewMessage   func(message *irc.Line,event string)
)
