// Package internal provides definitions and functions for P2P communication between chordal nodes
package chord

import (
	"context"
	"log"
	"net"
	"time"

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
	Info           *pb.Node     // stores internal node information such as ip address, port, and hash
	HeartBeatDuration time.Duration
	PredecessorID  []byte       //TODO
	FingerTable    *FingerTable //TODO
	HeartBeatTimer       *time.Timer
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

func NewNode(Name, IpAddr string, port int32, virtualNode bool, heartBeatDuration time.Duration, configs ...grpc.DialOption) *Node {

	if ipAddr := net.ParseIP(IpAddr); ipAddr != nil {
		log.Fatalln("Invalid ip address")
	}

	config := createGrpcDialConfig(configs...)
	return &Node{
		Name:           Name,
		Info:           &pb.Node{Address: IpAddr, Port: port}, // TODO
		HeartBeatDuration:       heartBeatDuration,
		PredecessorID:  nil,
		FingerTable:    nil,//TODO
		HeartBeatTimer:      time.NewTimer(heartBeatDuration),
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
func (node *Node) Kill() {
	close(node.MsgBuffer)
	node.grpcServer.Stop()
}

// Gracefully closes all connections and listeners from the rpc server
func (node *Node) Stop() {
	close(node.MsgBuffer)
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
	return NewNode(ID, "localhost", int32(port), false, time.Second*4)
}

// accepts new node to fingertable
func (node *Node) updateFingerTable(newNode *Node) error {
	currentNodeAddrHash := getHash(node.getServerAddress()) // used to compare with key, and other values in the finger table
	newNodeAddrHash := getHash(newNode.getServerAddress())
	for i, v := range *node.FingerTable {
		vAddrHash := v.ID

		if isBetween(newNodeAddrHash, currentNodeAddrHash, vAddrHash) {

			entry, err := node.FingerTable.getIthEntry(i)
			if err != nil {
				log.Println(err)
			}

			entry.UpdateValues(newNodeAddrHash, newNode)

		}

	}

	return nil
}


func (node *Node) HeartBeat(){
	ctx := context.Background()
	for{
		select{
		case <- node.HeartBeatTimer.C:
			node.SendHeartBeat(ctx, nil)
		}
	}
}



func (node *Node)AcknowledgeHeartBeat(targetNodeId string, hb *pb.HeartBeat,  err error){
	if err != nil{
		clnt := node.ConnectionPool[targetNodeId].Client

		if clnt != nil{ // If client is nil attempt to connect and ack

		}
	}	
}
