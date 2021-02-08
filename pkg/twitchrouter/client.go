package twitchrouter

import (
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
)

type Message struct{
	Msg *string
	Command *string
	Uuid *string
	MsgId *string
}
func Send(){

}
func Client(cmd string,help string,accessLevel int32,onMessage func(request *Message, send func(*MessageRequest) error)){

	println("try to connect")
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa( int(8765)), grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	println("connected")
	defer conn.Close()
	c := NewTwitchRouterClient(conn)

	// Contact the server and print out its response.
	ctx := context.Background()


	stream, err := c.Message(ctx)
	if err!= nil{
		println(err.Error())
		return
	}
	waitc := make(chan struct{})
	go func() {
		stream.Send(&MessageRequest{
			Command: &cmd,
			Help:    &help,
			AccessLevel: &accessLevel,
		})
		for {

			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}

			onMessage( &Message{
				Msg:     in.Msg,
				MsgId: in.Msgid,
				Command: in.Command,
				Uuid:    in.Uuid,
			},stream.Send)
		}
	}()
	<-waitc
	stream.CloseSend()
}
