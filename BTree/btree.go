package BTree

import (
	"fmt"
)

type Item int

type Key interface {
	Less(than interface{}) bool
}

type BTree struct {
	degree  int
	isLeaf 	bool
	Keys	[]Key
	Children[]interface{}
	parent  *BTree
	Next    *BTree
}

func New(degree int) *BTree {
	t := &BTree{
		degree:degree,
		isLeaf:true,
		Keys:make([]Key, 0),
		Children:make([]interface{}, 0),
	}
	return t
}

func (bt *BTree)Copy() *BTree {
	tree := &BTree{
		degree:bt.degree,
		isLeaf:bt.isLeaf,
		Keys:bt.Keys,
		Children:bt.Children,
		parent:bt.parent,
		Next:bt.Next,
	}
	return tree
}

func (tree *BTree)String() string {
	s := "《"
	//s = fmt.Sprintf("degree: %v", tree.degree)
	s = fmt.Sprintf("%s isLeaf: %v", s, tree.isLeaf)
	s = fmt.Sprintf("%s Keys: %v", s, tree.Keys)
	s = fmt.Sprintf("%s Children: \n\t%v", s, tree.Children)
	s = fmt.Sprintf("%s 》", s)
	return s
}

func (t *BTree)SearchLeafTree(key Key) *BTree {
	// 找到适当的叶节点
	if t.isLeaf {
		return t
	}
	i := 0
	for  {
		if i < len(t.Keys) && t.Keys[i].Less(key) {
			i++
		} else {
			break
		}
	}
	return t.Children[i].(*BTree).SearchLeafTree(key)
}

// 在叶节点插入键值
func insertKVtoLeaf(btree *BTree, key Key, val interface{}) {
	i := 0
	for {
		if i < len(btree.Keys) && btree.Keys[i].Less(key) {
			i ++
		} else {
			break
		}
	}
	keys := btree.Keys
	btree.Keys = append(keys[:i], key)
	btree.Keys = append(btree.Keys, keys[i:]...)
	//tree.Keys = append(tree.Keys[:i], key, tree.Keys[i:]...)
	children := btree.Children
	btree.Children = append(children[:i], val)
	btree.Children = append(btree.Children, children[i:]...)
}

// 将现有的节点分裂出新的节点
func splitOutNewNode(btree *BTree) *BTree {
	newNode := New(btree.degree)
	newNode.Keys = btree.Keys[ceil(float32(btree.degree)/2 - 1):]
	newNode.Children = btree.Children[ceil(float32(btree.degree)/2 - 1):]
	newNode.parent = btree.parent

	btree.Keys = btree.Keys[:ceil(float32(btree.degree)/2 - 1)]
	btree.Children = btree.Children[:ceil(float32(btree.degree)/2 - 1)]
	return newNode
}

// 依照旧的节点，在新的节点上找到新的key
func findoutNewKey(oldTree, newTree *BTree) Key {
	var rightMostKey Key = oldTree.Keys[len(oldTree.Keys) - 1]
	var leftFirstKey Key
	for i := 0; i < len(newTree.Keys); i++ {
		leftFirstKey = newTree.Keys[i]
		if rightMostKey.Less(leftFirstKey) {
			break
		}
	}
	return leftFirstKey
}

// 生成新的父节点
func newFather(child *BTree) *BTree {
	parent := New(child.degree)
	parent.isLeaf = false
	parent.Children = append(parent.Children, child.Copy())
	return parent
}
func (t *BTree)Insert(key Key, val Item) {
	// 我们设法在适当的叶节点中为新建找到空闲空间, 如果有的话，我们就把健放在那里。
	// 找到适当的叶节点
	btree := t.SearchLeafTree(key)

	// 在叶节点适当的位置插入key,val
	insertKVtoLeaf(btree, key, val)

	// 如果插入后的叶节点超过限制的大小，那么就将叶节点分裂成两个
	if len(btree.Keys) >= btree.degree {
		// 分裂出新的节点
		newNode := splitOutNewNode(btree)
		// 叶节点需要处理Next链接
		newNode.Next = btree.Next
		btree.Next = newNode

		// 找到要插入到父节点的key
		leftFirstKey := findoutNewKey(btree, newNode)

		leftNode := btree
		if btree.parent == nil {
			// 如果没有父节点 需要生成一个父节点，然后将左右节点加入
			parent := newFather(btree)
			btree.parent = parent
			newNode.parent = parent
			parent.insertNode(leftFirstKey, newNode, leftNode)
			*btree = *parent
		} else {
			// 有父节点，直接加
			newNode.parent.insertNode(leftFirstKey, newNode, leftNode)
		}
	}
}

func insertKVtoNode(btree *BTree, key Key, val interface{}, siblingNode *BTree)  {
	i := 0
	for {
		if i < len(btree.Keys) && siblingNode != btree.Children[i] {
			i++
		} else {
			break
		}
	}

	keys := btree.Keys
	fmt.Printf("has %d keys\n", len(keys))

	btree.Keys = append(keys[:i], key)
	btree.Keys = append(btree.Keys, keys[i:]...)

	fmt.Printf("has %d keys\n", len(btree.Keys))

	children := btree.Children
	fmt.Printf("has %d children\n", len(children))

	btree.Children = append(children[:i+1], val)
	btree.Children = append(btree.Children, children[i+1:]...)

	fmt.Printf("has %d children\n", len(btree.Children))
}
// 插入子节点 subNode, subNode的相邻兄弟节点是node
func (btree *BTree)insertNode(key Key, subNode, siblingNode *BTree)  {
	// 插入子节点到兄弟节点旁
	insertKVtoNode(btree, key, subNode, siblingNode)

	if len(btree.Keys) >= btree.degree {
		fmt.Printf("go into")
		// 分裂出新的节点
		newNode := splitOutNewNode(btree)

		// 找到要插入到父节点的key
		leftFirstKey := findoutNewKey(btree, newNode)

		// 如果没有父节点 需要生成一个父节点，然后将左右节点加入
		leftNode := btree
		if btree.parent == nil {
			// 如果没有父节点 需要生成一个父节点，然后将左右节点加入
			parent := newFather(btree)
			btree.parent = parent
			newNode.parent = parent
			parent.insertNode(leftFirstKey, newNode, leftNode)
			*btree = *parent
		} else {
			// 有父节点，直接加
			newNode.parent.insertNode(leftFirstKey, newNode, leftNode)
		}

	}
}

func ceil(x float32) int {
	x2 := int(x)
	if x > float32(x2) {
		return x2 + 1
	}
	return x2
}