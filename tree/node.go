package tree

type Node struct {
	key   string
	value interface{}

	left   *Node
	right  *Node
	parent *Node
}

func NewNode(p *Node, k string, v interface{}) *Node {
	return &Node{
		parent: p,
		key:    k,
		value:  v,
	}
}

func (n *Node) IsLeap() bool {
	if n.left == nil || n.right == nil {
		return true
	}
	return false
}
