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
		time.Sleep(time.Second)
		n.Connect( "localhost", 8081, grpc.WithInsecure())

		reply, err := n.ConnectionPool["localhost:8081"].Client.EchoReplyRPC(
			ctx,
			&pb.PingMessage{
				Info:      "node2 sending message",
				Id:        "node2",
				Timestamp: timestamppb.Now(),
			})
		if err != nil {
			log.Println(err)
		}

		fmt.Println(reply.Info)
		fmt.Println(reply.Id)

		//n.ConnectionPool["node1"].Conn.Close()
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
				n.Connect("localhost", 8081, grpc.WithInsecure())

				reply, err := n.ConnectionPool["localhost:8081"].Client.EchoReplyRPC(
					ctx,
					&pb.PingMessage{
						Info:      "ACK",
						Id:        "node2",
						Timestamp: timestamppb.Now(),
					})
				if err != nil {
					log.Println(err)
				}

				fmt.Println(reply.Info)
				fmt.Println(reply.Id)

		//		n.ConnectionPool["node1"].Conn.Close()

			}
		}

	}
}
func main() {

	/*	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
		if err != nil{
			log.Fatal(err)
		}

		client := pb.NewNodeAgentClient(conn)
	*/

	n := chord.NewDefaultNode("node2", 8082)
	go Ping(n)
	go Acknowledge(n)
	n.Start()

}
