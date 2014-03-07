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
	// TODO: We don't end up using the dimension anywhere - so maybe
	// drop it? I do like that the tree struct wraps the interface
	// of the nodes themselves.
	dimension int
	root      *node
}

// Create a new BTree with the given dimension.
func NewBTree(dimension int) *BTree {
	// Note that the root starts off as a leaf.
	rootNode := &node{true, 2 * dimension, 0, nil, make([]item, 2*dimension+1), nil}
	tree := &BTree{dimension, rootNode}
	return tree
}

// Add a key value pair into the tree.
func (tree *BTree) Insert(key int, value interface{}) {
	fmt.Println("Adding value", key, "to tree.")
	tree.root.insert(item{key, value}, nil)
}

// Determine the number of items in the tree.
func (tree *BTree) Size() int {
	return tree.root.size()
}

// Find the value of the first item in the tree with the same
// key. If there are multiple items with the same key, the first
// found will be returned. If the key is not found, nil will be
// returned.
func (tree *BTree) Search(key int) interface{} {
	return tree.root.search(key)
}

func (tree *BTree) Remove(key int) interface{} {
	// TODO
	return nil
}

// Function to insert an item into a node.
// This function may call recursively into its child nodes to find the
// correct location.
func (node *node) insert(value item, child *node) {
	fmt.Println("Adding value", value, "child:", child, "to node", node)
	// If this node is a leaf, then clearly we need to insert into the list.
	// If there is a child pointer, then insert as well since this is
	// probably coming back up the tree from a node splitting.
	if node.isLeaf || child != nil {
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
		node.children[node.currentSize].insert(value, nil)
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

			if !node.isLeaf {
				bumpedChild := node.children[i+1]
				node.children[i+1] = child
				child = bumpedChild
			}
		}
	}

	node.items[node.currentSize] = value
	if !node.isLeaf {
		node.children[node.currentSize+1] = child
	}
	node.currentSize += 1
}

// Split a node which has too many items - ie currentSize is larger
// than maxSize. This is done by creating a new leaf node to hold half
// of the items in the current node, then inserting this into the
// parent above this node (which may cause it to split, but that is
// handled by the insertion code). In the case of the root node splitting,
// that must be handled specially.
func (currentNode *node) splitNode() {
	fmt.Println("Splitting:", currentNode)
	// Create a new node for half of these children.
	rightNode := &node{true, currentNode.maxSize, 0, currentNode.parent,
		make([]item, len(currentNode.items)), nil}
	if currentNode.children != nil {
		rightNode.children = make([]*node, 1+cap(currentNode.items))
	}

	// The median node for the data in this node.
	middleIndex := len(currentNode.items) / 2
	median := currentNode.items[middleIndex]
	currentNode.items[middleIndex] = item{0, nil}
	currentNode.currentSize--

	for i := middleIndex + 1; i < len(currentNode.items); i++ {
		rightNode.items[rightNode.currentSize] = currentNode.items[i]
		if currentNode.children != nil {
			rightNode.isLeaf = false
			rightNode.children[rightNode.currentSize] = currentNode.children[i];
			rightNode.children[rightNode.currentSize].parent = rightNode
		}
		rightNode.currentSize++
		currentNode.items[i] = item{0, nil}
		currentNode.currentSize--
	}
	if currentNode.children != nil {
		rightNode.children[rightNode.currentSize] = currentNode.children[len(currentNode.items)];
		rightNode.children[rightNode.currentSize].parent = rightNode
	}

	// If we have a parent node, then insert into it (it might split further)
	// but then we are done here.
	// If there is no parent, then we need to create a new node. The idea is
	// to keep pointers to this node correct, but move half of the children
	// into a new left node.
	if currentNode.parent != nil {
		currentNode.parent.insert(median, rightNode)
		return
	} else {
		leftNode := &node{true, currentNode.maxSize, 0, currentNode,
			make([]item, len(currentNode.items)), nil}
		if currentNode.children != nil {
			leftNode.isLeaf = false
			leftNode.children = make([]*node, 1+cap(currentNode.items))
		}

		for i := 0; i < middleIndex; i++ {
			leftNode.items[i] = currentNode.items[i]
			if currentNode.children != nil {
				leftNode.children[i] = currentNode.children[i];
				leftNode.children[i].parent = leftNode
			}
			currentNode.items[i] = item{0, nil}
			leftNode.currentSize++
		}
		if currentNode.children != nil {
			leftNode.children[middleIndex] = currentNode.children[middleIndex];
			leftNode.children[middleIndex].parent = leftNode
		}
		// The current node now only has one item - this is only
		// allowed at the root of the tree.
		currentNode.currentSize = 1
		currentNode.items[0] = median
		// This node is no longer a leaf.
		if currentNode.isLeaf {
			currentNode.isLeaf = false
			currentNode.children = make([]*node, 1+cap(currentNode.items))
		}
		//
		rightNode.parent = currentNode
		currentNode.children[0] = leftNode
		currentNode.children[1] = rightNode
	}
}

func (n *node) search(key int) interface{} {
	fmt.Println("Searching for", key, "in", n)
	if n.isLeaf {
		// If we are at a leaf node, search through the items list
		// until the end or we have found a key which is larger
		// than the search key.
		for i := 0; i < n.currentSize && key >= n.items[i].key; i++ {
			if n.items[i].key == key {
				return n.items[i].value
			}
		}
	} else {
		// Search through the list to find the first node
		// which is larger than the key which indicates that
		// the data is in the matching child node.
		for i := 0; i < n.currentSize; i++ {
			if key == n.items[i].key {
				return n.items[i].value
			}
			if key < n.items[i].key {
				return n.children[i].search(key)
			}
		}
		return n.children[n.currentSize].search(key)
	}
	// The item is not in the tree.
	return nil
}

// Determine the total size of the tree below this node, including the
// items contained in this node.
// In theory we could track this at the root, but we can also do it this
// way for fun.
func (node *node) size() int {
	totalSize := node.currentSize
	fmt.Println(node)
	if !node.isLeaf {
		for i := 0; i < node.currentSize; i++ {
			totalSize += node.children[i].size()
		}
		totalSize += node.children[node.currentSize].size()

	}
	return totalSize
}

// Return a sorted list of the keys in under this node.
func (node *node) keyTraversal() []int {
	// TODO: Pass slice pointers in such that we can create the initial
	// slice of the right size in a single allocation.
	results := make([]int, 0)
	for i := 0; i < node.currentSize; i++ {
		if !node.isLeaf {
			results = append(results, node.children[i].keyTraversal()...)
		}
		results = append(results, node.items[i].key)
	}
	if !node.isLeaf {
		results = append(results, node.children[node.currentSize].keyTraversal()...)
	}
	return results
}
