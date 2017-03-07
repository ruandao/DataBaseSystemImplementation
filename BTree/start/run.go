package main

import (
	"github.com/ruandao/DataBaseSystemImplementation/BTree"
	"fmt"
)


type X struct {
	x int
}

func (x X)Less(x2 interface{}) bool {
	x3 := x2.(X)
	return x.x < x3.x
}

func main() {
	tree := BTree.New(3)
	tree.Insert(X{3}, 4)
	tree.Insert(X{8}, 9)
	tree.Insert(X{8}, 9)
	//tree.Insert(X{8}, 9)
	//tree.Insert(X{10}, 20)
	//tree.Insert(X{11}, 21)
	fmt.Printf("%v", tree)
}
