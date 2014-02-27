/*
Trying to implement a b-tree.
http://en.wikipedia.org/wiki/B-tree
*/

package BTree

import (
	"fmt"
)

type Item struct {
	key int
	value interface{}
}

type Node struct {
	// Metadata items.
	isLeaf bool
	maxSize int
	currentSize int

	// The parent of this node, possibly nil.
	parent *Node


	items []Item
	children []*Node
}

type BTree struct {
	dimension int
	root *Node
}



func NewBTree(dimension int) *BTree {
	rootNode := &Node{true, 2 * dimension, 0, nil, make([]Item, 2 * dimension), nil}
	tree := &BTree{dimension, rootNode}
	return tree
}

func (tree *BTree) Insert(key int, value interface{}) {
	fmt.Println("Adding value", key, "to tree.")
	tree.root.Insert(Item{key, value}, nil)
}

func (node *Node) Insert(value Item, child *Node) bool {
	fmt.Println("Adding value", value, "to node.")
	if (node.isLeaf || child != nil) {
		node.insertItemIntoNode(value, child)
		if node.currentSize > node.maxSize {
			node.splitNode()

		}
	} else {
		for i := 0; i < node.currentSize; i++ {
			if value.key <= node.items[i].key {
				return node.children[i].Insert(value, nil)
			}
		}
		return node.children[node.currentSize + 1].Insert(value, nil)
	}
	return true
}

func (node *Node) insertItemIntoNode(value Item, child *Node) {
	for i := 0; i < node.currentSize; i++ {
		if value.key < node.items[i].key {
			bumpedItem := node.items[i]
			node.items[i] = value
			value = bumpedItem

			if (!node.isLeaf) {
				bumpedChild := node.children[i + 1]
				node.children[i + 1] = child
				child = bumpedChild
			}
		}
	}

	node.items[node.currentSize] = value
	if !node.isLeaf {
		node.children[node.currentSize + 1] = child
	}
	node.currentSize += 1
}

func (node *Node) splitNode() {
	// Create a new node for half of these children.
	rightNode := &Node{true, node.maxSize, 0, node,
		make([]Item, len(node.items)), nil}

	// The median node for the data in this node.
	middleIndex := len(node.items) / 2
	median := node.items[middleIndex]

	//node.items[middleIndex] = nil
	node.currentSize--

	for i := middleIndex + 1; i < len(node.items); i++ {
		rightNode.items[rightNode.currentSize] = node.items[i]
		rightNode.currentSize++
		//node.items[i] = nil;
		node.currentSize--;
	}

	if node.parent != nil {
		node.parent.Insert(median, rightNode)
	} else {


	}
}

func (node *Node) size() int {
	totalSize := node.currentSize
	if !node.isLeaf {
		for i := 0; i < node.currentSize; i++ {
			totalSize += node.children[i].size()
		}
		totalSize += node.children[node.currentSize].size()

	}
	return totalSize
}

func (node *Node) traversal() []Item {
	// TODO: Actually, I think append will handle increasing capacity
	// so I don't need to do this. Which is good because calling size()
	// at each node ends up quadratic.
	results := make([]Item, node.size())
	for i := 0; i < node.currentSize; i++ {
		if (!node.isLeaf) {
			results = append(results, node.children[i].traversal()...)
		}
		results = append(results, node.items[i])
	}
	if (!node.isLeaf) {
		results = append(results, node.children[node.currentSize].traversal()...)
	}
	return results
}
