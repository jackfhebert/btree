package BTree

import (
	"testing"
)

func Test_BTreeConstructor(t *testing.T) {
	tree := NewBTree(5)
	if tree.root == nil {
		t.Error("no root node")
	}
	// Dimension is not max size of the internal nodes.
	if tree.root.maxSize != 10 {
		t.Error("node has wrong max size")
	}
	if tree.root.currentSize != 0 {
		t.Error("node started with wrong initial size.")
	}
}

// Test adding elements to a node when no splitting is required.
// This is the basic case of just keeping the items in sorted order.
func Test_AddAFewElementsNoSplitting(t *testing.T) {
	// Build the tree to hold enough items that we won't have to split nodes.
	tree := NewBTree(5)

	// Add a single item just fine.
	tree.Insert(2, "foo")
	if tree.root.items[0].key != 2 {
		t.Error("item not inserted where expected.", tree.root.items)
	}
	if tree.root.items[0].value != "foo" {
		t.Error("item not inserted where expected.")
	}
	if tree.root.currentSize != 1 {
		t.Error("node has the wrong size", tree.root)
	}

	// Add an item which should be sorted after the initial one.
	tree.Insert(3, "bar")
	if tree.root.items[1].key != 3 {
		t.Error("item not inserted where expected.", tree.root.items)
	}
	if tree.root.items[1].value != "bar" {
		t.Error("item not inserted where expected.")
	}
	if tree.root.currentSize != 2 {
		t.Error("node has the wrong size", tree.root)
	}

	// This one should be inserted in the middle to keep the sorted order.
	tree.Insert(2, "baz")
	if tree.root.items[1].key != 2 {
		t.Error("item not inserted where expected.", tree.root.items)
	}
	if tree.root.items[2].key != 3 {
		t.Error("item not inserted where expected.", tree.root.items)
	}
	if tree.root.items[1].value != "baz" {
		t.Error("item not inserted where expected.")
	}
	if tree.root.items[2].value != "bar" {
		t.Error("item not inserted where expected.")
	}
	if tree.root.currentSize != 3 {
		t.Error("node has the wrong size", tree.root)
	}
}

// Test inserting an item into a node when it also has a pointer to a
// right child. This is part of the logic for a node splitting.
func Test_InsertWithChildren(t *testing.T) {
	// The parent node for the tree. Set the initial size to 1 since
	// we setup these manually.
	root := Node{false, 5, 1, nil, make([]Item, 5), make([]*Node, 5)}
	// Start it off with some initial data.
	root.items[0] = Item{0, "initial"}
	root.children[0] = &Node{true, 5, 0, nil, make([]Item, 5), nil}
	root.children[0].Insert(Item{-1, "left child"}, nil)
	root.children[1] = &Node{true, 5, 0, nil, make([]Item, 5), nil}
	root.children[1].Insert(Item{1, "right child"}, nil)
	if root.children[1].size() != 1 {
		t.Error("wrong total size", root.children[1])
	}
	if root.size() != 3 {
		t.Error("wrong total size", root)
	}

	lowNode := &Node{true, 5, 0, nil, make([]Item, 5), nil}
	lowNode.Insert(Item{3, "new right child"}, nil)
	root.Insert(Item{2, "foo"}, lowNode);
	if root.currentSize != 2 {
		t.Error("wrong size on root node", root)
	}
	if root.items[1].key != 2 {
		t.Error("wrong second item", root.items)
	}
	if root.children[0].items[0].key != -1 {
		t.Error("wrong first child", root.children)
	}
	if root.children[1].items[0].key != 1 {
		t.Error("wrong second child", root.children[1])
	}
	if root.children[2].items[0].key != 3 {
		t.Error("wrong third child", root.children)
	}

	highNode := &Node{true, 5, 0, nil, make([]Item, 5), nil}
	highNode.Insert(Item{12, "high right child"}, nil)
	root.Insert(Item{10, "bar"}, highNode);
	if root.currentSize != 3 {
		t.Error("wrong size on root node", root)
	}
	if root.items[2].key != 10 {
		t.Error("wrong third item", root.items)
	}
	if root.children[3].items[0].key != 12 {
		t.Error("wrong fourth child", root.children)
	}

	midNode := &Node{true, 5, 0, nil, make([]Item, 5), nil}
	midNode.Insert(Item{7, "mid right child"}, nil)
	root.Insert(Item{5, "baz"}, midNode);
	if root.currentSize != 4 {
		t.Error("wrong size on root node", root)
	}
	if root.items[2].key != 5 {
		t.Error("wrong third item", root.items)
	}
	if root.children[3].items[0].key != 7 {
		t.Error("wrong fourth child", root.children[2])
	}
}

// Test splitting a node when the parent node has enough space such that
// further splitting is not required.
func Test_SplitNoParentHasRoom(t *testing.T) {
	root := Node{false, 5, 0, nil, make([]Item, 5), make([]*Node, 5)}
	// Start it off with some initial data.
	root.items[0] = Item{0, "initial"}
	root.children[0] = &Node{true, 3, 0, nil, make([]Item, 3), nil}
	root.children[0].parent = &root
	root.children[0].Insert(Item{-1, "left child"}, nil)

	root.children[1] = &Node{true, 3, 0, nil, make([]Item, 3), nil}
	root.children[1].parent = &root
	root.children[1].Insert(Item{1, "right child"}, nil)

	// Now add some nodes.
}
