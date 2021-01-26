// Package internal provides definitions and functions for P2P communication between chordal nodes
package node

import (
	"context"
	"fmt"

	"log"
	"net"

	pb "github.com/frandiazrio/arca/src/api/node/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//Implements NodeAgentServerInterface.
// Provides the server rpc API
type nodeAgentServerAPI struct {
	
}
// Node type provides definition of a chord node
type Node struct {
	ID             string
	IpAddr         string
	nodeAgent      nodeAgentServerAPI
	grpcServer     *grpc.Server
	listener       *net.TCPListener
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

	nodeAgentServer := nodeAgentServerAPI{}
	config := createGrpcDialConfig(configs...)
	return &Node{
		ID:             ID,
		IpAddr:         IpAddr,
		nodeAgent:      nodeAgentServer,
		grpcServer:     grpc.NewServer(),
		listener:       nil,
		port:           port,
		connConfig:     config,
		virtualNode:    virtualNode,
		connectionPool: make(map[string]grpc.ClientConn),
	}
}
func (node *Node) getServerAddress() string{
	return fmt.Sprintf("%s:%d", node.IpAddr, node.port)
}
// Start node service by binding to the assigned address to the node
func (node *Node) Start() *Node {
	var err error
	
	listener, err := net.Listen("tcp", node.getServerAddress())

	if err != nil {
		// error creating listener on address
		log.Println("Error creating listener")
		log.Fatal(err)
	}

	node.listener = listener.(*net.TCPListener)
	// Using the grpc server in node
	// Node implements the NodeAgentServer interface so we can use directly in here to start the service
	pb.RegisterNodeAgentServer(node.grpcServer, &node.nodeAgent)
	fmt.Printf("Starting node server on %s", node.getServerAddress())
	if err = node.grpcServer.Serve(node.listener); err != nil {
		log.Println("Error starting server")
		log.Fatal(err)
	}
	
	return node
}

func (node *Node) Connect(targetNode *Node) {
	grpc.Dial(targetNode.IpAddr, node.connConfig...)
}

// Creates a grpc Dial Options slice
func createGrpcDialConfig(configs ...grpc.DialOption) []grpc.DialOption {
	config := []grpc.DialOption{}

	for _, cfg := range configs {
		config = append(config, cfg)
	}
	return config
}

func NewDefaultNode(ID string, port int) *Node {
	return NewNode(ID, "localhost", port, false)
}

func (nodeAgent *nodeAgentServerAPI) EchoReply(ctx context.Context, pingMsg *pb.PingMessage) (*pb.PingMessage, error) {
	log.Println("RCV - CONTENT: ", pingMsg.Info)
	return &pb.PingMessage{
		Info:      "Echo Message",
		Timestamp: timestamppb.Now(),
	}, nil
}

func (node *Node) EchoReply(ctx context.Context, pingMsg *pb.PingMessage)(*pb.PingMessage, error){
	msg, err := node.nodeAgent.EchoReply(ctx, pingMsg)

	if err != nil{
		log.Println(err)
	}


	return msg, err
} 
