// Package internal provides definitions and functions for P2P communication between chordal nodes
package internal

import (
	"context"

	"log"
	"net"

	pb "github.com/frandiazrio/arca/src/api/node"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NodeAgentServer struct{}

// Node type provides definition of a chord node
type Node struct {
	ID             string
	IpAddr         string
	nodeAgent      NodeAgentServer
	grpcServer     *grpc.Server
	listener       *net.Listener
	port           int
	connConfig     []grpc.DialOption
	virtualNode    bool
	connectionPool map[string]grpc.ClientConn
}

func (node *Node) IsVirtualNode() bool {
	return node.virtualNode
}

func NewNode(ID, IpAddr string, port int, virtualNode bool, configs ...grpc.DialOption) *Node {

	if ipAddr := net.ParseIP(IpAddr); ipAddr != nil {
		log.Fatalln("Invalid ip address")
	}
	
	config := make([]grpc.DialOption, 2)
	for _,cfg := range configs{
		config = append(config, cfg)
	}
	return &Node{
		ID:             ID,
		IpAddr:         IpAddr,
		grpcServer:     grpc.NewServer(),
		listener:       nil,
		port:           port,
		connConfig:    	config, 
		virtualNode:    virtualNode,
		connectionPool: make(map[string]grpc.ClientConn),
	}
}

// Start node service by binding to the assigned address to the node
func (node *Node) Start()*Node {
	var err error
	*node.listener, err = net.Listen("tcp", node.IpAddr)

	if err != nil {
		// error creating listener on address
		log.Fatal(err)
	}

	// Using the grpc server in node
	// Node implements the NodeAgentServer interface so we can use directly in here to start the service
	pb.RegisterNodeAgentServer(node.grpcServer, &node.nodeAgent)
	if err = node.grpcServer.Serve(*node.listener); err != nil {
		log.Fatal("failed to serve %v", err)
	}

	return node
}

func (node *Node) Connect(targetNode *Node){
	grpc.Dial(targetNode.IpAddr,node.connConfig...)
}



func NewDefaultNode(ID string, port int) *Node {
	return NewNode(ID, "localhost", port, false)
}

func (nodeAgent *NodeAgentServer) EchoReply(ctx context.Context, pingMsg *pb.PingMessage) (*pb.PingMessage, error) {
	log.Println("RCV - CONTENT: ", pingMsg.Info)
	return &pb.PingMessage{
		Info:      "Echo Message",
		Timestamp: timestamppb.Now(),
	}, nil
}
