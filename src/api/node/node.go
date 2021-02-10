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
	Name              string
	Active            bool
	Info              *pb.Node //
	HeartBeatDuration time.Duration
	PredecessorID     []byte       //TODO
	FingerTable       *FingerTable //TODO
	HeartBeatTicker   *time.Ticker
	MsgBuffer         chan int
	grpcServer        *grpc.Server
	listener          *net.TCPListener
	connConfig        []grpc.DialOption
	virtualNode       bool
	ConnectionPool    map[string]*grpcNodeConn
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
		Name:              Name,
		Active:            false,
		Info:              &pb.Node{Address: IpAddr, Port: port}, // TODO
		HeartBeatDuration: heartBeatDuration,
		PredecessorID:     nil,
		FingerTable:       nil, //TODO
		HeartBeatTicker:   time.NewTicker(heartBeatDuration),
		grpcServer:        grpc.NewServer(),
		MsgBuffer:         make(chan int),
		listener:          nil,
		connConfig:        config,
		virtualNode:       virtualNode,
		ConnectionPool:    make(map[string]*grpcNodeConn),
	}
}

func (node *Node) EndPoint() string {
	return endpoint(node.Info.Address, int(node.Info.Port))
}

// Start node service by binding to the assigned address to the node
func (node *Node) Start() *Node {
	var err error

	listener, err := net.Listen("tcp", node.EndPoint())

	if err != nil {
		// error creating listener on address
		log.Println("Error creating listener")
		log.Fatal(err)
	}

	node.listener = listener.(*net.TCPListener)
	// Using the grpc server in node
	// Node implements the NodeAgentServer interface so we can use directly in here to start the service
	pb.RegisterNodeAgentServer(node.grpcServer, node)
	log.Printf("Starting node server on %s \n", node.EndPoint())
	if err = node.grpcServer.Serve(node.listener); err != nil {
		log.Println("Error starting server")
		log.Fatal(err)
	}
	node.Active = true
	return node
}

// It immediately closes all connections and listeners from the rpc server
func (node *Node) Kill() {
	node.HeartBeatTicker.Stop()
	close(node.MsgBuffer)
	node.grpcServer.Stop()
}

// Gracefully closes all connections and listeners from the rpc server
func (node *Node) Stop() {
	node.HeartBeatTicker.Stop()
	close(node.MsgBuffer)
	node.grpcServer.GracefulStop()
}

// For now accept a string, with the fingertable implementation this will change
// TODO lookup Node on fingertable
// connect to targetNode ip and port and returns a client that can perform the grpc calls

func (node *Node) Connect(targetIPAddr string, targetPort int, config ...grpc.DialOption) *grpcNodeConn {
	ctx, cancel := context.WithTimeout(context.Background(), 4*node.HeartBeatDuration) // Use 3 HeartBeats before cancelling
	defer cancel()

	var conn *grpc.ClientConn
	var err error
	if len(config) == 0 { // Use some default grpc DialOption (meh)
		conn, err = grpc.DialContext(ctx, endpoint(targetIPAddr, targetPort), grpc.WithInsecure())
	} else {

		conn, err = grpc.DialContext(ctx, endpoint(targetIPAddr, targetPort), config...)
	}
	if err != nil {
		// Error establising connection
		log.Fatal(err)
	}

	addr := endpoint(targetIPAddr, targetPort)
	client := pb.NewNodeAgentClient(conn)
	grpcConn := grpcNodeConn{
		TargetID:      addr, // We can later hash this to get a hash id
		TargetIPAddr:  targetIPAddr,
		TargetPort:    targetPort,
		Client:        client,
		Conn:          conn,
		LastTimeStamp: timestamppb.Now(),
		Active:        true,
	}

	node.ConnectionPool[endpoint(targetIPAddr, targetPort)] = &grpcConn

	return &grpcConn

}

// Hashes the nodes ip address to get the node id
func (node *Node) GetNodeId() []byte {
	return getHash(node.EndPoint())
}

func NewDefaultNode(ID string, port int) *Node {
	return NewNode(ID, "localhost", int32(port), false, time.Second*4, grpc.WithInsecure())
}

// accepts new node to fingertable
func (node *Node) updateFingerTable(newNode *Node) error {
	currentNodeAddrHash := getHash(node.EndPoint()) // used to compare with key, and other values in the finger table
	newNodeAddrHash := getHash(newNode.EndPoint())
	for i, v := range *node.FingerTable {
		vAddrHash := v.ID

		if isBetween(newNodeAddrHash, currentNodeAddrHash, vAddrHash) {

			entry, err := node.FingerTable.getIthEntry(i)
			if err != nil {
				log.Println(err)
			}

			entry.UpdateValues(newNodeAddrHash, newNode.Info)

		}

	}

	return nil
}

// Check to see if node has a connection with a targetNode
func (node *Node) Connected(targetNode *pb.Node) bool {
	// first check to see if it is in the connection pool

	nodeConn := node.ConnectionPool[endpoint(targetNode.Address, int(targetNode.Port))]

	if validConnState(nodeConn.Conn) {
		return true
	}

	return false
}

// Send HeartBeat to the ipaddress of a targetNode
// this communicates that current node (node) is active
func (node *Node) SendHeartBeat(targetNode *pb.Node) {

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*node.HeartBeatDuration)
		defer cancel()
		for {
			select {
			case <-node.HeartBeatTicker.C:
				if node.Connected(targetNode) {
					node.ConnectionPool[endpoint(targetNode.Address, int(targetNode.Port))].Client.HeartBeatRPC(
						ctx,
						&pb.HeartBeat{
							SourceNode: node.Info,
							Timestamp:  timestamppb.Now(),
						})
				} else { //attempt to connect
					node.Connect(targetNode.Address, int(targetNode.Port))
				}
			}
		}
	}()
}
