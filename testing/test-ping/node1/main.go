package main


import (
	"github.com/frandiazrio/arca/src/api/node"
)

func main(){
	//Acting as server	
	n := node.NewDefaultNode("node1", 8081)
	n.Start()

}
