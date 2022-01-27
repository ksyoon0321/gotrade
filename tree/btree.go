package tree

/*
샘플 프로젝트는 날짜베이스 TID를 사용하는 구조라서 Btree는 구조적으로 맞지 않음.
시계열 로그를 사용하는게 맞지만 이건 연습이라 그냥 Go로 구현해 보는것의 의의를 둔다.

btree효율을 높이려면 root키를 12시로 잡으면 분배는 좌우 비슷하게 될듯함.
*/
import (
	"strings"
)

type BTree struct {
	root *Node
}

func NewBTree() *BTree {
	btree := &BTree{}

	return btree
}

func (t *BTree) Put(k string, v interface{}) {
	if t.root == nil {
		t.root = NewNode(nil, k, v)
	} else {
		t.putLeap(t.root, k, v)
	}

}

func (t *BTree) Get(k string) interface{} {
	node := t.find(t.root, k)
	if node == nil {
		return nil
	}

	return node.value
}

// func (t *BTree) Remove(k string) {
// 	node := t.find(t.root, k)
// 	if node == nil {
// 		return
// 	}

// 	pnode := node.parent
// 	//중복 부분을 변수로 쪼개는게 더 이해하기 힘든 코드라서 중복작성하는게 좋을듯
// 	if pnode.left == node {
// 		//left 삭제
// 		if node.right != nil {
// 			pnode.left = node.right
// 			pnode.left.parent = pnode
// 		} else {
// 			pnode.left = node.left
// 			if pnode.left != nil {
// 				pnode.left.parent = pnode
// 			}
// 		}
// 	} else {
// 		//right 삭제
// 		if node.left != nil {
// 			pnode.right = node.left
// 			pnode.right.parent = pnode

// 			if node.right != nil {
// 				node.left = node.right
// 				node.right = nil
// 			}
// 		} else {
// 			pnode.right = node.right
// 			if pnode.right != nil {
// 				pnode.right.parent = pnode
// 			}
// 		}
// 	}

// 	node = nil
// }

func (t *BTree) PrintTree() string {
	return "\n\n" + t.print(t.root, "[ROOT]", 0)
}

func (t *BTree) print(node *Node, prefix string, depth int) string {
	if node == nil {
		return ""
	}

	retv := prefix + node.key
	ndepth := depth + 1

	ntab := ""
	for ii := 0; ii < ndepth; ii++ {
		ntab += "\t"
	}
	nprefixL := "\n" + ntab + "->[L]"
	nprefixR := "\n" + ntab + "->[R]"
	retv += t.print(node.left, nprefixL, ndepth)
	retv += t.print(node.right, nprefixR, ndepth)

	return retv
}

func (t *BTree) less(nodek, k string) bool {
	return nodek < k
}

func (t *BTree) putLeap(node *Node, k string, v interface{}) {
	if node.key == k {
		node.value = v
	} else {
		newnode := NewNode(node, k, v)
		if t.less(k, node.key) {
			if node.left != nil {
				t.putLeap(node.left, k, v)
			} else {
				node.left = newnode
			}
		} else {
			if node.right != nil {
				t.putLeap(node.right, k, v)
			} else {
				node.right = newnode
			}
		}
	}
}

func (t *BTree) find(node *Node, k string) *Node {
	if node == nil {
		return nil
	}

	if node.key == k {
		return node
	} else {
		//k > left.key
		if strings.Compare(k, node.left.key) == 1 {
			return t.find(node.right, k)
		} else {
			return t.find(node.left, k)
		}
	}
}
