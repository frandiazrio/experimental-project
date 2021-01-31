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

const (
	M_SHA1 = 160 // number of bits for sha-1 or 20-byte hash value
)

type grpcNodeConn struct {
	TargetID      string
	TargetIPAddr  string
	TargetPort    int
	Conn          *grpc.ClientConn
	Client        pb.NodeAgentClient
	LastTimeStamp *timestamppb.Timestamp
	Active        bool
}

// Node type provides definition of a chord node
type Node struct {
	Name           string
	ID             []byte
	SuccessorId    []byte       //TODO
	PredecessorID  []byte       //TODO
	FingerTable    *FingerTable //TODO
	IpAddr         string       // A node's identifier is chosen by hashing the node's IpAddr
	MsgBuffer      chan int
	grpcServer     *grpc.Server
	listener       *net.TCPListener
	port           int
	connConfig     []grpc.DialOption
	virtualNode    bool
	ConnectionPool map[string]*grpcNodeConn
}

func (node *Node) IsVirtualNode() bool {
	return node.virtualNode
}

func NewNode(Name, IpAddr string, port int, virtualNode bool, configs ...grpc.DialOption) *Node {

	if ipAddr := net.ParseIP(IpAddr); ipAddr != nil {
		log.Fatalln("Invalid ip address")
	}

	config := createGrpcDialConfig(configs...)
	return &Node{
		Name:           Name,
		ID:             getHash(address(IpAddr, port)),
		SuccessorId:    nil, // TODO
		PredecessorID:  nil,
		FingerTable:    nil, //TODO
		IpAddr:         IpAddr,
		grpcServer:     grpc.NewServer(),
		MsgBuffer:      make(chan int),
		listener:       nil,
		port:           port,
		connConfig:     config,
		virtualNode:    virtualNode,
		ConnectionPool: make(map[string]*grpcNodeConn),
	}
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
	log.Printf("Starting node server on %s \n", node.getServerAddress())
	if err = node.grpcServer.Serve(node.listener); err != nil {
		log.Println("Error starting server")
		log.Fatal(err)
	}

	return node
}

// For now accept a string, with the fingertable implementation this will change
// TODO lookup Node on fingertable
// connect to targetNode ip and port and returns a client that can perform the grpc calls

func (node *Node) Connect(targetID, targetIPAddr string, targetPort int, config ...grpc.DialOption) *grpcNodeConn {
	conn, err := grpc.Dial(address(targetIPAddr, targetPort), config...)

	if err != nil {
		// Error establising connection
		log.Fatal(err)
	}

	client := pb.NewNodeAgentClient(conn)
	return &grpcNodeConn{
		TargetID:      targetID,
		TargetIPAddr:  targetIPAddr,
		TargetPort:    targetPort,
		Client:        client,
		Conn:          conn,
		LastTimeStamp: timestamppb.Now(),
		Active:        true,
	}
}

// Hashes the nodes ip address to get the node id
func (node *Node) GetNodeId() []byte {
	return getHash(node.getServerAddress())
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
