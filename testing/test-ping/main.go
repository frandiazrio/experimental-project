package main


import (
	"github.com/frandiazrio/arca/src/api/node"
)

func main(){
	n := node.NewDefaultNode("localhost", 8081)
	n.Start()

}
