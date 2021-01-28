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
type grpcNodeConn struct{
	targetID string
	targetIPAddr string
	targetPort int
	conn *grpc.ClientConn
	client pb.NodeAgentClient
	lastTimeStamp *timestamppb.Timestamp
	active bool
}
// Node type provides definition of a chord node
type Node struct {
	ID             string
	IpAddr         string
	MsgBuffer 	   chan int
	grpcServer     *grpc.Server
	listener       *net.TCPListener
	port           int
	connConfig     []grpc.DialOption
	virtualNode    bool
	ConnectionPool map[string] *grpcNodeConn
}

func (node *Node) IsVirtualNode() bool {
	return node.virtualNode
}

func NewNode(ID, IpAddr string, port int, virtualNode bool, configs ...grpc.DialOption) *Node {

	if ipAddr := net.ParseIP(IpAddr); ipAddr != nil {
		log.Fatalln("Invalid ip address")
	}

	config := createGrpcDialConfig(configs...)
	return &Node{
		ID:             ID,
		IpAddr:         IpAddr,
		grpcServer:     grpc.NewServer(),
		MsgBuffer: 		make(chan int),
		listener:       nil,
		port:           port,
		connConfig:     config,
		virtualNode:    virtualNode,
		ConnectionPool: make(map[string]*grpcNodeConn),
	}
}

func address(ipaddr string, port int)string{
	return fmt.Sprintf("%s:%d", ipaddr, port)
}

func (node *Node) getServerAddress() string {
	return address(node.IpAddr, node.port)
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
	pb.RegisterNodeAgentServer(node.grpcServer, node)
	fmt.Printf("Starting node server on %s \n", node.getServerAddress())
	if err = node.grpcServer.Serve(node.listener); err != nil {
		log.Println("Error starting server")
		log.Fatal(err)
	}

	return node
}


// For now accept a string, with the fingertable implementation this will change
// TODO lookup Node on fingertable
// connect to targetNode ip and port and returns a client that can perform the grpc calls 

func (node *Node) Connect(targetID, targetIPAddr string, targetPort int, config ...grpc.DialOption)*grpcNodeConn {
	conn , err := grpc.Dial(address(targetIPAddr, targetPort) , config...)

	if err != nil{
		// Error establising connection
		log.Fatal(err)
	}
		
	client := pb.NewNodeAgentClient(conn)
	return &grpcNodeConn{
		targetID: targetID, 
		targetIPAddr: targetIPAddr,
		targetPort: targetPort, 
		client: client, 
		conn: conn,
		lastTimeStamp: timestamppb.Now(), 
		active: true,

	}
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

func (node *Node) EchoReply(ctx context.Context, pingMsg *pb.PingMessage) (*pb.PingMessage, error) {

	log.Println("RCV - CONTENT: ", pingMsg.Info)
	node.MsgBuffer <- CONNECT
	node.MsgBuffer <- ACK
	return &pb.PingMessage{
		Info:      fmt.Sprintf("Message from %v", node.ID),
		
		Timestamp: timestamppb.Now(),
	}, nil


}
