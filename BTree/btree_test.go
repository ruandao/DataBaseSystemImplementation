package BTree

import (
	"testing"
	"fmt"
)

func TestNew(t *testing.T) {
	tree := New(3)
	if tree.isLeaf != true {
		t.Errorf("New出来的对象必须是叶子节点")
	}
}

type X struct {
	x int
	//tree *BTree
}

func (x X)Less(x2 interface{}) bool {
	x3 := x2.(X)
	return x.x < x3.x
}

func TestBTree_Insert(t *testing.T) {
	tree := New(3)
	tree.Insert(X{3}, 4)
	tree.Insert(X{8}, 9)
	tree.Insert(X{8}, 9)
	tree.Insert(X{8}, 9)
	tree.Insert(X{10}, 20)
	tree.Insert(X{11}, 21)
	fmt.Printf("%v\n", tree)
}

