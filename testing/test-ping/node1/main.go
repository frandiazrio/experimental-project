package main

import (
	"context"
	"fmt"
	"log"
	"time"

	chord "github.com/frandiazrio/arca/src/api/node"
	pb "github.com/frandiazrio/arca/src/api/node/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)
func Ping(n *chord.Node) {
	for {
		ctx := context.Background()
		time.Sleep(3*time.Second)
		n.Connect("node2", "localhost", 8082, grpc.WithInsecure())
		reply, err := n.ConnectionPool["localhost:8082"].Client.EchoReplyRPC(
			ctx,
			&pb.PingMessage{
				Info:      "node1 sending message",
				Id:        "node1",
				Timestamp: timestamppb.Now(),
			})
		if err != nil {
			log.Println(err)
		}

		fmt.Println(reply.Info)
		fmt.Println(reply.Id)
		//n.ConnectionPool["node2"].Conn.Close()
		/*	nodeReply, err := client.EchoReply(ctx, &pb.PingMessage{
				Info: "Message from client",
				Id: "node2",
				Timestamp: timestamppb.Now(),
			})
		*/

	}
}
func Acknowledge(n *chord.Node) {
	for {
		select {
		case msg := <-n.MsgBuffer:
			if msg == chord.ACK {
				ctx := context.Background()
				time.Sleep(time.Second)
				n.Connect("node2", "localhost", 8082, grpc.WithInsecure())

				reply, err := n.ConnectionPool["localhost:8082"].Client.EchoReplyRPC(
					ctx,
					&pb.PingMessage{
						Info:      "ACK",
						Id:        "node1",
						Timestamp: timestamppb.Now(),
					})
				if err != nil {
					log.Println(err)
				}

				fmt.Println(reply.Info)
				fmt.Println(reply.Id)

				//n.ConnectionPool["node2"].Conn.Close()

			}
		}

	}
}
func main() {
	//Acting as server
	n := chord.NewDefaultNode("node1", 8081)
	go Ping(n)
	go Acknowledge(n)
	n.Start()

}
