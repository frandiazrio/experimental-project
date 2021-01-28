package main

import (
	"context"
	"fmt"
	"log"

	"github.com/frandiazrio/arca/src/api/node"

	pb "github.com/frandiazrio/arca/src/api/node/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main(){
	
/*	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil{
		log.Fatal(err)
	}

	client := pb.NewNodeAgentClient(conn)
*/

	node := node.NewDefaultNode("node2", 8082)

	ctx := context.Background()
	reply, err := node.Connect("localhost", 8081, grpc.WithInsecure()).EchoReply(
		ctx, 
		&pb.PingMessage{
			Info: "Message from node2",
			Id: "node2",
			Timestamp: timestamppb.Now(),
		})

	

/*	nodeReply, err := client.EchoReply(ctx, &pb.PingMessage{
		Info: "Message from client",
		Id: "node2",
		Timestamp: timestamppb.Now(),
	})
*/

	if err != nil{
		log.Println(err)
	}


	fmt.Println(reply.Info)
	fmt.Println(reply.Id)

}

