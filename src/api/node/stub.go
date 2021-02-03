package chord

import (
	"context"
	"fmt"
	"log"

	pb "github.com/frandiazrio/arca/src/api/node/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Implement NodeAgentServer interface for the stub on the rpc node

func (node *Node) EchoReply(ctx context.Context, pingMsg *pb.PingMessage) (*pb.PingMessage, error) {

	log.Println("RCV - CONTENT: ", pingMsg.Info)
	if pingMsg.Info != "ACK" {
		go func(n *Node) {
			n.MsgBuffer <- ACK
		}(node)
	}

	return &pb.PingMessage{
		Info: fmt.Sprintf("Sending message to %s", pingMsg.Id),

		Timestamp: timestamppb.Now(),
	}, nil

}

//TODO
func (node *Node) FindSuccesor(ctx context.Context, targetNode *pb.Node) (*pb.Node, error) {
	return nil, nil
}

//TODO (node needs to broadcast to another node that is exists and update fingertable if necessary)
func (node *Node) AddNodeToFingerTable(ctx context.Context, targetNode *pb.Node) (*pb.Empty, error) {
	return nil, nil
}

func (node *Node) ConfirmHeartBeat(ctx context.Context,  hb *pb.HeartBeat) (*pb.Bool, error) {
	if hb != nil{
		return &pb.Bool{Value:true}, nil 
	}
	return &pb.Bool{Value:false}, nil
}
func (node *Node)SendHeartBeat(ctx context.Context, void *pb.Empty) (*pb.HeartBeat, error) {
	return &pb.HeartBeat{
		SourceNode: node.Info,	
		Timestamp: timestamppb.Now(),
	}, nil
}
