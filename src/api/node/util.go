package node

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"hash"
)

func isEqual(a, b []byte) bool {
	return bytes.Compare(a, b) == 0
}

//Compares if key is between (a,b)
func isBetween(key, a, b []byte) bool {
	switch bytes.Compare(a, b) {
	case 0: // a, b are Equal, key cannot be the same as a b by definition
		return false
	case 1: // a > b, or b < a implies key > b and key < a
		return bytes.Compare(key, b) == 1 && bytes.Compare(key, a) == -1
	case -1: // a < b implies key < b and key > a
		return bytes.Compare(key, b) == -1 && bytes.Compare(key, a) == 1
	}

	return false
}
func isPowerOfTwo(num int) bool {
	return (num != 0) && ((num & (num - 1)) == 0)
}

func address(ipaddr string, port int) string {
	return fmt.Sprintf("%s:%d", ipaddr, port)
}

func hashFunc(key []byte, h hash.Hash) []byte {
	if _, err := h.Write(key); err != nil {
		return nil
	}

	return h.Sum(nil)
}
func getHash(idKey string) []byte {

	return hashFunc([]byte(idKey), sha1.New())

}
