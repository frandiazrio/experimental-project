// Package internal provides definitions and functions for P2P communication between chordal nodes
package chord

import (

	"log"
	"net"

	pb "github.com/frandiazrio/arca/src/api/node/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	Info             *pb.Node // stores internal node information such as ip address, port, and hash
	SuccessorId    []byte       //TODO
	PredecessorID  []byte       //TODO
	FingerTable    *FingerTable //TODO
	MsgBuffer      chan int
	grpcServer     *grpc.Server
	listener       *net.TCPListener
	connConfig     []grpc.DialOption
	virtualNode    bool
	ConnectionPool map[string]*grpcNodeConn
}

func (node *Node) IsVirtualNode() bool {
	return node.virtualNode
}

func NewNode(Name, IpAddr string, port int32, virtualNode bool, configs ...grpc.DialOption) *Node {

	if ipAddr := net.ParseIP(IpAddr); ipAddr != nil {
		log.Fatalln("Invalid ip address")
	}

	config := createGrpcDialConfig(configs...)
	return &Node{
		Name:           Name,
		Info:             &pb.Node{Address: IpAddr, Port: port}, // TODO
		SuccessorId:    nil, // TODO
		PredecessorID:  nil,
		FingerTable:    nil, //TODO
		grpcServer:     grpc.NewServer(),
		MsgBuffer:      make(chan int),
		listener:       nil,
		connConfig:     config,
		virtualNode:    virtualNode,
		ConnectionPool: make(map[string]*grpcNodeConn),
	}
}

func (node *Node) getServerAddress() string {
	return address(node.Info.Address, int(node.Info.Port))
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



// It immediately closes all connections and listeners from the rpc server
func (node *Node) Kill(){
	node.grpcServer.Stop()
}

// Gracefully closes all connections and listeners from the rpc server
func (node *Node) Stop(){
	node.grpcServer.GracefulStop()
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
	return NewNode(ID, "localhost", int32(port), false)
}









