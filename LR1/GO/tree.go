package main

import (
	"fmt"
	"io"
	"strings"
)

type CBNode struct {
	value int
	left  *CBNode
	right *CBNode
}

type CompleteBinaryTree struct {
	root *CBNode
	size int
}

func (t *CompleteBinaryTree) TInsert(v int) {
	t.root = insertAt(t.root, t.size, 0, v)
	t.size++
}

func insertAt(root *CBNode, target, idx, v int) *CBNode {
	if root == nil {
		if idx == target {
			return &CBNode{value: v}
		}
		return nil
	}
	if 2*idx+1 <= target {
		root.left = insertAt(root.left, target, 2*idx+1, v)
	}
	if 2*idx+2 <= target {
		root.right = insertAt(root.right, target, 2*idx+2, v)
	}
	return root
}

func (t *CompleteBinaryTree) TFind(v int) (int, bool) {
	return findDepth(t.root, v, 0)
}

func findDepth(root *CBNode, v, depth int) (int, bool) {
	if root == nil {
		return 0, false
	}
	if root.value == v {
		return depth, true
	}
	if d, ok := findDepth(root.left, v, depth+1); ok {
		return d, true
	}
	return findDepth(root.right, v, depth+1)
}

func (t *CompleteBinaryTree) TPrint(w io.Writer) {
	fmt.Fprint(w, "T[")
	inorder(t.root, w)
	fmt.Fprintln(w, "]")
}

func inorder(n *CBNode, w io.Writer) {
	if n == nil {
		return
	}
	inorder(n.left, w)
	fmt.Fprint(w, n.value, " ")
	inorder(n.right, w)
}

func (t *CompleteBinaryTree) TPrintTree(w io.Writer) {
	printTree(t.root, 0, w)
}

func printTree(n *CBNode, level int, w io.Writer) {
	if n == nil {
		return
	}
	printTree(n.right, level+1, w)
	fmt.Fprint(w, strings.Repeat("\t", level), n.value, "\n")
	printTree(n.left, level+1, w)
}

func (t *CompleteBinaryTree) TIsComplete() bool {
	count := 0
	return isCompleteRec(t.root, 0, t.size, &count) && count == t.size
}

func isCompleteRec(n *CBNode, idx, size int, count *int) bool {
	if n == nil {
		return true
	}
	if idx >= size {
		return false
	}
	*count++
	return isCompleteRec(n.left, 2*idx+1, size, count) && isCompleteRec(n.right, 2*idx+2, size, count)
}

func (t *CompleteBinaryTree) TIsFull() bool {
	return isFullRec(t.root)
}

func isFullRec(n *CBNode) bool {
	if n == nil || (n.left == nil && n.right == nil) {
		return true
	}
	if n.left == nil || n.right == nil {
		return false
	}
	return isFullRec(n.left) && isFullRec(n.right)
}
