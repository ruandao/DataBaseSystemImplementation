package BTree

type Item int

type Key interface {
	Less(than Key) bool
}

type Tree struct {
	degree  int
	isLeaf 	bool
	Keys	[]Key
	Children[]interface{}
	parent  *Tree
	Next    *Tree
}

func New(degree int) *Tree {
	t := &Tree{
		degree:degree,
		isLeaf:true,
		Keys:make([]Key,0),
		Children:make([]interface{}, 0),
	}
	return t
}

func (t *Tree)SearchLeafTree(key Key) *Tree {
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
	return t.Children[i].(*Tree).SearchLeafTree(key)
}

func (t *Tree)Insert(key Key, val Item) {
	// 我们设法在适当的叶节点中为新建找到空闲空间, 如果有的话，我们就把健放在那里。
	// 找到适当的叶节点
	tree := t.SearchLeafTree(key)
	// 在叶节点适当的位置插入key,val
	i := 0
	for {
		if i < len(tree.Keys) && tree.Keys[i].Less(key) {
			i ++
		} else {
			break
		}
	}

	tree.Keys = append(tree.Keys[:i], key, tree.Keys[i:]...)
	tree.Children = append(tree.Children[:i], val, tree.Children[i:]...)
	// 如果插入后的叶节点超过限制的大小，那么就将叶节点分裂成两个
	if len(tree.Keys) >= tree.degree {
		newNode := New(tree.degree)
		newNode.Keys = tree.Keys[ceil(tree.degree/2 - 1):]
		newNode.Children = tree.Children[ceil(tree.degree/2 - 1):]
		newNode.Next = tree.Next
		newNode.parent = tree.parent
		tree.Keys = tree.Keys[:ceil(tree.degree/2 - 1)]
		tree.Children = tree.Children[:ceil(tree.degree/2 - 1)]
		tree.Next = newNode
		// 找到要插入到父节点的key
		var rightMostKey Key = tree.Keys[len(tree.Keys) - 1]
		var leftFirstKey Key
		for i := 0; i < len(newNode.Keys); i++ {
			leftFirstKey = newNode.Keys[i]
			if rightMostKey.Less(leftFirstKey) {
				break
			}
		}
		// 如果没有父节点 需要生成一个父节点，然后将左右节点加入
		leftNode := tree
		if tree.parent == nil {
			parent := New(tree.degree)
			parent.isLeaf = false
			tree.parent = parent
			newNode.parent = parent
			parent.Children = append(parent.Children, tree)
			*tree = parent
		}
		newNode.parent.insertNode(rightMostKey, newNode, leftNode)
	}
}
// 插入子节点 subNode, subNode的相邻兄弟节点是node
func (tree *Tree)insertNode(key Key, subNode, node *Tree)  {
	i := 0
	for {
		if i < len(tree.Children) && node != tree.Children[i] {
			i++
		} else {
			break
		}
	}
	tree.Keys = append(tree.Keys[:i], key, tree.Keys[i:]...)
	tree.Children = append(tree.Children[:i+1], subNode, tree.Children[i+1:]...)
	if len(tree.Keys) >= tree.degree {
		// 需要分裂成两个节点
		newNode := New(tree.degree)
		newNode.Keys = tree.Keys[ceil(tree.degree/2 - 1):]
		newNode.Children = tree.Children[ceil(tree.degree/2 - 1):]
		newNode.parent = tree.parent
		tree.Keys = tree.Keys[:ceil(tree.degree/2 - 1)]
		tree.Children = tree.Children[:ceil(tree.degree/2 - 1)]
		// 找到要插入到父节点的key
		var rightMostKey Key = tree.Keys[len(tree.Keys) - 1]
		var leftFirstKey Key
		for i := 0; i < len(newNode.Keys); i++ {
			leftFirstKey = newNode.Keys[i]
			if rightMostKey.Less(leftFirstKey) {
				break
			}
		}
		// 如果没有父节点 需要生成一个父节点，然后将左右节点加入
		leftNode := tree
		if tree.parent == nil {
			parent := New(tree.degree)
			parent.isLeaf = false
			tree.parent = parent
			newNode.parent = parent
			parent.Children = append(parent.Children, tree)
			*tree = parent
		}
		newNode.parent.insertNode(rightMostKey, newNode, leftNode)
	}
}

func ceil(x float32) int {
	x2 := int(x)
	if x > float32(x2) {
		return x2 + 1
	}
	return x2
}