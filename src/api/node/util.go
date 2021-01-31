package node

import (
	"fmt"
	"bytes"
	"crypto/sha1"
)



func isEqual(a, b []byte)bool{
	return bytes.Compare(a, b) == 0
}

func isPowerOfTwo(num int) bool{
	return (num !=0 ) && ( (num & (num-1))==0)
}

func address(ipaddr string, port int)string{
	return fmt.Sprintf("%s:%d", ipaddr, port)
}

func getHash(idKey string)[]byte{
	h := sha1.New() // hasher 

	if _, err := h.Write([]byte(idKey)); err != nil{
		return nil
	}

	return h.Sum(nil)

}
