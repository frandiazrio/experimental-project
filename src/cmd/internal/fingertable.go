package internal

import "math/big"

type fingerEntry struct{
	ID []byte 	
	*Node
}

type FingerTable struct{
	fingerTable []fingerEntry
}

func NewFingerEntry(ID []byte, node *Node)*fingerEntry{
	return &fingerEntry{
		ID: ID,
		Node: node,
	}
}


// offset =  (n+2^i) mod (2^m)
func fingerID(n []byte, i, m int)[]byte{
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

func NewFingerTable()FingerTable{

}
