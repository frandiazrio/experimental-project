package node

import (
	"errors"
	"math/big"
)

// type for the entries in the fingertable
// Every finger entry will contain the hash id of the node and a pointer to the node the correspoding id belongs to
type fingerEntry struct {
	ID []byte
	*Node
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

func newFingerEntry(id []byte, node *Node)*fingerEntry{
	return &fingerEntry{
		ID: id,
		Node: node,
	}
}


func newFingerTable(tableSize int, n *Node)(*FingerTable, error){
	if tableSize < 0{
		return nil, errors.New("Error creating finger table: Size less than 0")
	}
	fingertable := make(FingerTable, tableSize)

	for i:=0; i< tableSize; i++{
		fingertable[i] = newFingerEntry(fingerID(n.ID, i, tableSize), n)
	}
	return &fingertable, nil
}
