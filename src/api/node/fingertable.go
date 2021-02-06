package chord

import (
	"crypto/sha1"
	"errors"
	"math/big"

	pb "github.com/frandiazrio/arca/src/api/node/proto"
)

// type for the entries in the fingertable
// Every finger entry will contain the hash id of the node and a pointer to the node the correspoding id belongs to
type fingerEntry struct {
	ID []byte
	*pb.Node
}

func (fe *fingerEntry) UpdateValues(ID []byte, node *pb.Node) {
	fe.ID = ID
	fe.Node = node
}

// type for the finger table
type FingerTable []*fingerEntry

// offset =  (n+2^i) mod (2^m)
func fingerID(n []byte, i, m int) []byte {
	idInt := (&big.Int{}).SetBytes(n)

	bigTwo := big.NewInt(2)

	offset := big.Int{}

	// 2^i
	offset.Exp(bigTwo, big.NewInt(int64(i)), nil)

	//(n+2^i)
	sum := big.Int{}

	sum.Add(idInt, &offset)

	// Ceil

	ceil := big.Int{}

	ceil.Exp(bigTwo, big.NewInt(int64(m)), nil)

	idInt.Mod(&sum, &ceil)

	return idInt.Bytes()
}

func newFingerEntry(id []byte, node *pb.Node) *fingerEntry {
	return &fingerEntry{
		ID:   id,
		Node: node,
	}
}

func newFingerTable(tableSize int, n *pb.Node) (*FingerTable, error) {
	if tableSize < 0 {
		return nil, errors.New("Error creating finger table: Size less than 0")
	}
	fingertable := make(FingerTable, tableSize)

	for i := 0; i < tableSize; i++ {
		addr := address(n.Address, int(n.Port))
		fingertable[i] = newFingerEntry(fingerID(hashFunc([]byte(addr), sha1.New()), i, tableSize), n)

	}

	return &fingertable, nil
}

func (ft *FingerTable) getIthEntry(i int) (*fingerEntry, error) {
	if i >= len(*ft) {
		return nil, errors.New("Invalid index")
	}

	return (*ft)[i], nil
}

func (ft *FingerTable) getIthFinger(i int) (*pb.Node, error) {
	entry, err := ft.getIthEntry(i)

	if err != nil {
		return nil, errors.New("Error getting node: Invalid index")
	}

	return entry.Node, nil
}
