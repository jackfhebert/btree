/*
Trying to implement a b-tree.
http://en.wikipedia.org/wiki/B-tree
*/

package BTree

import (
	"fmt"
)

// An item inside of a btree.
type item struct {
	// The key used for sorting items.
	// TODO: int, string? do I have to choose?
	key int
	// The data being placed inside of the tree.
	// TODO: Is this really the right type? 
	value interface{}
}

// Internal node for the tree.
type node struct {
	// Metadata items.
	//
	// Is this node a leaf in the tree?
	isLeaf bool
	// How many items should this leaf hold before it splits?
	maxSize int
	// The number of items currently held in this node.
	currentSize int

	// The parent of this node, possibly nil.
	parent *node

	// The data items inside this node. These should be in sorted order.
	items []item
	// If not a leaf, these are the child nodes.
	// Note that for item[n], items in child[n] are all less than it and
	// items in child[n+1] are all larger than it. This also implies
	// n+1 items in the children list for n items in the items list.
	children []*node
}

// The external interface to the tree.
type BTree struct {
	dimension int
	root *node
}


// Create a new BTree with the given dimension.
func NewBTree(dimension int) *BTree {
	rootNode := &node{true, 2 * dimension, 0, nil, make([]item, 2 * dimension), nil}
	tree := &BTree{dimension, rootNode}
	return tree
}


// Add a key value pair into the tree.
func (tree *BTree) Insert(key int, value interface{}) {
	fmt.Println("Adding value", key, "to tree.")
	tree.root.insert(item{key, value}, nil)
}


func (node *node) insert(value item, child *node) {
	fmt.Println("Adding value", value, "to node.")
	// If this node is a leaf, then clearly we need to insert into the list.
	// If there is a child pointer, then insert as well since this is
	// probably coming back up the tree from a node splitting.
	if (node.isLeaf || child != nil) {
		node.insertItemIntoNode(value, child)
		// If we passed the max size, then split.
		if node.currentSize > node.maxSize {
			node.splitNode()

		}
	} else {
		// Find the correct child node to insert into.
		// Note that we know there is no child pointer
		// to handle since we checked for that above.
		for i := 0; i < node.currentSize; i++ {
			if value.key <= node.items[i].key {
				node.children[i].insert(value, nil)
				return
			}
		}
		// If the item to add is larger than all of the items, then it
		// is handled by the last child node.
		node.children[node.currentSize + 1].insert(value, nil)
		return
	}
}


// Insert the item into the current node.
// This differs from the node.insert() function above in that here we
// always add to the current items list and do not worry about splitting.
// The goal here is just to keep the list of items[] and children[] sorted.
func (node *node) insertItemIntoNode(value item, child *node) {
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

func (currentNode *node) splitNode() {
	// Create a new node for half of these children.
	rightNode := &node{true, currentNode.maxSize, 0, currentNode,
		make([]item, len(currentNode.items)), nil}

	// The median node for the data in this node.
	middleIndex := len(currentNode.items) / 2
	median := currentNode.items[middleIndex]

	//node.items[middleIndex] = nil
	currentNode.currentSize--

	for i := middleIndex + 1; i < len(currentNode.items); i++ {
		rightNode.items[rightNode.currentSize] = currentNode.items[i]
		rightNode.currentSize++
		//node.items[i] = nil;
		currentNode.currentSize--;
	}

	if currentNode.parent != nil {
		currentNode.parent.insert(median, rightNode)
	} else {


	}
}

// Determine the total size of the tree.
// In theory we could track this at the root, but we can also do it this
// way for fun.
func (node *node) size() int {
	totalSize := node.currentSize
	if !node.isLeaf {
		for i := 0; i < node.currentSize; i++ {
			totalSize += node.children[i].size()
		}
		totalSize += node.children[node.currentSize].size()

	}
	return totalSize
}

func (node *node) traversal() []item {
	// TODO: Actually, I think append will handle increasing capacity
	// so I don't need to do this. Which is good because calling size()
	// at each node ends up quadratic.
	results := make([]item, node.size())
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
