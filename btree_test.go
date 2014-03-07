package BTree

import (
	"fmt"
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
	root := node{false, 5, 1, nil, make([]item, 5), make([]*node, 5)}
	// Start it off with some initial data.
	root.items[0] = item{0, "initial"}
	root.children[0] = &node{true, 5, 0, nil, make([]item, 5), nil}
	root.children[0].insert(item{-1, "left child"}, nil)
	root.children[1] = &node{true, 5, 0, nil, make([]item, 5), nil}
	root.children[1].insert(item{1, "right child"}, nil)
	if root.children[1].size() != 1 {
		t.Error("wrong total size", root.children[1])
	}
	if root.size() != 3 {
		t.Error("wrong total size", root)
	}

	lowNode := &node{true, 5, 0, nil, make([]item, 5), nil}
	lowNode.insert(item{3, "new right child"}, nil)
	root.insert(item{2, "foo"}, lowNode)
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

	highNode := &node{true, 5, 0, nil, make([]item, 5), nil}
	highNode.insert(item{12, "high right child"}, nil)
	root.insert(item{10, "bar"}, highNode)
	if root.currentSize != 3 {
		t.Error("wrong size on root node", root)
	}
	if root.items[2].key != 10 {
		t.Error("wrong third item", root.items)
	}
	if root.children[3].items[0].key != 12 {
		t.Error("wrong fourth child", root.children)
	}

	midNode := &node{true, 5, 0, nil, make([]item, 5), nil}
	midNode.insert(item{7, "mid right child"}, nil)
	root.insert(item{5, "baz"}, midNode)
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
	root := node{false, 5, 1, nil, make([]item, 6), make([]*node, 7)}
	// Start it off with some initial data.
	root.items[0] = item{0, "initial"}
	root.children[0] = &node{true, 3, 0, nil, make([]item, 4), nil}
	root.children[0].parent = &root
	root.children[0].insert(item{-1, "left child"}, nil)

	root.children[1] = &node{true, 3, 0, nil, make([]item, 4), nil}
	root.children[1].parent = &root
	root.children[1].insert(item{1, "right child"}, nil)

	// Now add some nodes.
	root.insert(item{2, "right child (2)"}, nil)
	root.insert(item{3, "right child (3)"}, nil)
	root.insert(item{4, "right child (4)"}, nil)
	root.insert(item{5, "right child (5)"}, nil)

	if root.currentSize != 2 {
		t.Error("root has wrong size", root)
	}
	// 6 inserts called above- and 1 item manually inserted to start
	// it off.
	if root.size() != 7 {
		t.Error("root has wrong total size", root.size())
	}
	if root.children[0].currentSize != 1 {
		t.Error("first child has the wrong size", root.children[0])
	}
	if root.children[1].currentSize != 2 {
		t.Error("second child has the wrong size", root.children[1])
	}
}

// Test that we can split at the parent level.
func Test_SplitRoot(t *testing.T) {
	tree := NewBTree(2)
	tree.Insert(5, "first item")
	tree.Insert(4, "second item")
	tree.Insert(8, "third item")
	tree.Insert(7, "fourth item")
	tree.Insert(1, "fifth item")
	tree.Insert(1, "sixth item")
	tree.Insert(1, "seventh item")
	if tree.Size() != 7 {
		t.Error("tree has wrong size:", tree.Size(), tree)
	}
	expectedKeys := []int{1, 1, 1, 4, 5, 7, 8}
	keys := tree.root.keyTraversal()
	if len(keys) != len(expectedKeys) {
		t.Error("weird keys length: ", keys)
	} else {
		for i := 0; i < len(keys); i++ {
			if keys[i] != expectedKeys[i] {
				t.Error("weird key values: ", keys)
				break
			}

		}
	}
	if tree.root.currentSize != 1 {
		t.Error("root node has weird size:", tree.root)
	}

	// Confirm that search returns valid items from the tree.
	if tree.Search(5) != "first item" {
		t.Error("Search for key 5 failed")
	}
	if tree.Search(4) != "second item" {
		t.Error("Search for key 4 failed")
	}
	if tree.Search(8) != "third item" {
		t.Error("Search for key 8 failed")
	}
	if tree.Search(7) != "fourth item" {
		t.Error("Search for key 7 failed")
	}
	// Some values not expected to be found in the tree.
	if tree.Search(6) != nil {
		t.Error("search for key 6 failed")
	}
	if tree.Search(-2) != nil {
		t.Error("search for key 6 failed")
	}
	if tree.Search(100) != nil {
		t.Error("search for key 6 failed")
	}
}

func Test_AddManyAllIncreasing(t *testing.T) {
	tree := NewBTree(2);
	for i := 0; i < 50; i++ {
		tree.Insert(i, fmt.Sprintf("foo: %d", i));
		if tree.Size() != i+1 {
			t.Error("break in i:", i);
			break;
		}
	}
	if tree.Size() != 50 {
		t.Error("tree has wrong size:", tree.Size(), tree.root)
	}
}

func Test_AddManyAlternating(t *testing.T) {
	tree := NewBTree(2);
	for i := 0; i < 50; i++ {
		tree.Insert(i, fmt.Sprintf("foo: %d", i));
		tree.Insert(-i, fmt.Sprintf("foo: %d", -i));
		if tree.Size() != 2*i+1 {
			t.Error("break in i:", i);
			break;
		}
	}
	if tree.Size() != 100 {
		t.Error("tree has wrong size:", tree.Size(), tree.root)
	}
}