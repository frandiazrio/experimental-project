package main

import (
	"context"

	"github.com/frandiazrio/arca/src/api/node"
	"google.golang.org/grpc"
)

func main(){
	//Acting as server	
	n := node.NewDefaultNode("node1", 8081)

	go func(nd *node.Node){
		ctx := context.Background()
		for {
			select{
				case msg := <- nd.MsgBuffer:
					switch msg {
						case msg == node.CONNECT{
							
							nd.ConnectionPool["node2"] = nd.Connect("localhost", 8082, grpc.WithInsecure())
					}	
			

				}
			}
		}
		// Send ACK to another to confirm message was received
		
	}(n)




	n.Start()

}
