package chord

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/frandiazrio/arca/src/api/node/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Implement NodeAgentServer interface for the stub on the rpc node

func (node *Node) EchoReplyRPC(ctx context.Context, pingMsg *pb.PingMessage) (*pb.PingMessage, error) {

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
func (node *Node) FindSuccesorRPC(ctx context.Context, targetNode *pb.Node) (*pb.Node, error) {
	return nil, nil
}

//TODO (node needs to broadcast to another node that is exists and update fingertable if necessary)
func (node *Node) AddNodeToFingerTableRPC(ctx context.Context, targetNode *pb.Node) (*pb.Empty, error) {
	return nil, nil
}

func (node *Node) HeartBeatRPC(ctx context.Context, hb *pb.HeartBeat) (*pb.Bool, error) {

	// return True because heartbeat was received
	select {

	case <-ctx.Done():
		log.Println("Error: timeout") // TODO better handle nodes
		return &pb.Bool{
			Value: false,
		}, errors.New("Error: timeout")
	default:
		log.Printf("Info from heartbeat: %s \n %s \n", 
			endpoint(hb.SourceNode.Address,int(hb.SourceNode.Port)),
			hb.Timestamp.String())

		return &pb.Bool{
			Value: true,
		}, nil

	}

}
