package main

import (
	"context"
	"database/sql"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	_ "github.com/lib/pq"
	"github.com/saegewerk/BotZtwitch/pkg/config"
	DB "github.com/saegewerk/BotZtwitch/pkg/db"
	"github.com/saegewerk/BotZtwitch/pkg/twitchchat"
	"github.com/saegewerk/BotZtwitch/pkg/twitchrouter"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"time"
)



var (
	twitchConfig *twitchchat.Configuration
	twitch       *twitchchat.Chat
)

func initConfiguration(nickname,channel string) {
	oauth := "oauth:" + os.Getenv("TWITCH")
	twitchConfig = twitchchat.NewConfiguration(nickname, oauth, channel)
}



type server struct {
	twitchrouter.UnimplementedTwitchRouterServer
	apps     twitchrouter.AppsQueued
	messages chan string
}

func (s *server) Message(stream twitchrouter.TwitchRouter_MessageServer) error {
	ctx := stream.Context()
	Queue := make(chan twitchrouter.CmdMessage)
	Closed := false
	Registered := false
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		// receive data from stream
		req, err := stream.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			Closed = true
			log.Printf("receive error %v", err)
			continue
		}
		if req.Msg != nil {
			twitch.SendMessage(*req.Msg)
		}
		if !(Registered) {

			go func() {
				for {
					r := <-Queue
					uuid := r.Uuid.String()
					err = stream.Send(&twitchrouter.MessageResponse{
						Msg:     &r.Msg,
						Command: &r.Cmd,
						Uuid:    &uuid,
					})
					if err != nil {
						println(err.Error())
						Closed = true
						break
					}
				}
			}()
			s.apps.Queue <- twitchrouter.App{
				CmdRegister: twitchrouter.CmdRegister{
					Cmd:         *req.Command,
					Help:        *req.Help,
					AccessLevel: *req.AccessLevel,
				},
				Queue:  &Queue,
				Closed: &Closed,
			}
			Registered = true
		}

	}
	return nil
}

func main() {
	c, err := config.YAML()
	if err != nil {
		println(err.Error())
		return
	}
	initConfiguration(c.Twitch.Nickname,c.Twitch.Channel)
	db, err := sql.Open("postgres", c.SprintfDBConfig())
	if err != nil {
		println(err.Error())
		return
	}
	queries := DB.New(db)
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8765))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	server := &server{apps: twitchrouter.AppsQueued{
		Queue: make(chan twitchrouter.App),
		Apps:  make([]twitchrouter.App, 0),
	}}
	twitchrouter.RegisterTwitchRouterServer(grpcServer, server)
	twitch = twitchchat.NewChat(twitchConfig)
	go server.runWithCallbacks(twitch, queries)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) runWithCallbacks(twitch *twitchchat.Chat, queries *DB.Queries) {
	stop := make(chan struct{})
	defer close(stop)

	err := twitch.ConnectWithCallbacks(
		func() {
			fmt.Println("Connected to Twitch IRC")
		},
		func() {
			fmt.Println("Disconnected from Twitch IRC")
			stop <- struct{}{}
		},
		func(message *irc.Line, event string) {
			ctx := context.Background()
			userUUID, err := queries.GetUserUUIDByNick(ctx, message.Nick)
			if err != nil {
				if err == sql.ErrNoRows {
					userUUID, err = queries.InsertUser(ctx, DB.InsertUserParams{
						Nick:        message.Nick,
						Name:        message.Args[0][1:],
						Accesslevel: 0,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					})
					if err != nil {
						println(err.Error())
					}
				} else {
					println(err.Error())
				}
			}
			msgID := ""
			if event == "USERNOTICE" {
				msgID = message.Tags["msg-id"]
			}
			msgUUID, err := queries.InsertMessage(ctx, DB.InsertMessageParams{
				Msg:    message.Args[1],
				MsgID:  msgID,
				Event:  event,
				FkUser: userUUID,
			})
			for i := 0; i < len(s.apps.Apps); {
				if *s.apps.Apps[i].Closed {
					s.apps.Apps = append(s.apps.Apps[:i], s.apps.Apps[i+1:]...)
				} else {
					i++
				}
			}
			//check app queue
			select {
			case app, ok := <-s.apps.Queue:
				if ok {
					s.apps.Apps = append(s.apps.Apps, app)
				} else {
					fmt.Println("Channel closed!")
				}
			default:
			}
			log.Printf("Event: %s, Nick: %s, MsgId: %s, Msg: %s\n", event, message.Nick, msgID,message.Args[1])
			if message.Args[1][0] == '!' {

				if len(message.Args[1]) == 1 {
					res := ""
					for _, app := range s.apps.Apps {
						tmp := res + "!" + app.CmdRegister.Cmd + ": " + res + app.CmdRegister.Help + "; "
						//check if msg is too long
						if len(tmp)+len(res) >= 500 {
							twitch.SendMessage(res)
							res = tmp
						} else {
							res = res + tmp
						}
					}
					twitch.SendMessage(res)
				}
				for _, app := range s.apps.Apps {
					msg := message.Args[1][len(app.CmdRegister.Cmd)+1:]
					if message.Args[1][1:len(app.CmdRegister.Cmd)+1] == app.CmdRegister.Cmd &&
						"join" != app.CmdRegister.Cmd && "part" != app.CmdRegister.Cmd && "usernotice" != app.CmdRegister.Cmd {

						select {
						case *app.Queue <- twitchrouter.CmdMessage{
							Cmd:   app.CmdRegister.Cmd,
							Msg:   msg,
							MsgId: msgID,
							Uuid:  msgUUID,
						}:
						default:
						}
						//twitch.SendMessage(*res.Msg)
					}else if "join" == app.CmdRegister.Cmd && event=="join"{
						select {
						case *app.Queue <- twitchrouter.CmdMessage{
							Cmd:   app.CmdRegister.Cmd,
							Msg:   msg,
							MsgId: msgID,
							Uuid:  msgUUID,
						}:
						default:
						}
					}else if "part" == app.CmdRegister.Cmd && event=="part"{
						select {
						case *app.Queue <- twitchrouter.CmdMessage{
							Cmd:   app.CmdRegister.Cmd,
							Msg:   msg,
							MsgId: msgID,
							Uuid:  msgUUID,
						}:
						default:
						}
					}else if "usernotice" == app.CmdRegister.Cmd && event=="usernotice"{
						select {
						case *app.Queue <- twitchrouter.CmdMessage{
							Cmd:   app.CmdRegister.Cmd,
							Msg:   msg,
							MsgId: msgID,
							Uuid:  msgUUID,
						}:
						default:
						}
					}
				}

			}
		},
	)

	if err != nil {
		return
	}

	<-stop
}
