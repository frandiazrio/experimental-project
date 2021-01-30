package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/frandiazrio/arca/src/api/node"

	pb "github.com/frandiazrio/arca/src/api/node/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Ping(n *node.Node) {
	for {
		ctx := context.Background()
		time.Sleep(time.Second)
		n.ConnectionPool["node1"] = n.Connect("node1", "localhost", 8081, grpc.WithInsecure())

		reply, err := n.ConnectionPool["node1"].Client.EchoReply(
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
func Acknowledge(n *node.Node) {
	for {
		select {
		case msg := <-n.MsgBuffer:
			if msg == node.ACK {
				ctx := context.Background()
				time.Sleep(time.Second)
				n.ConnectionPool["node1"] = n.Connect("node1", "localhost", 8081, grpc.WithInsecure())

				reply, err := n.ConnectionPool["node1"].Client.EchoReply(
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

	n := node.NewDefaultNode("node2", 8082)
	go Ping(n)
	go Acknowledge(n)
	n.Start()

}
