package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/bblfsh/sdk.v1/uast"
)

func nodeTree() *uast.Node {
	child1 := &uast.Node {
		InternalType: "child1",
	}

	subchild21 := &uast.Node {
		InternalType: "subchild21",
	}

	subchild22 := &uast.Node {
		InternalType: "subchild22",
	}

	child2 := &uast.Node {
		InternalType: "child2",
		Children: []*uast.Node{subchild21, subchild22},
	}
	parent := &uast.Node{
		InternalType: "parent",
		Children: []*uast.Node{child1, child2},
	}
	return parent
}

func TestIter_Range(t *testing.T) {
	parent := nodeTree()

	iter, err := NewIterator(parent, PreOrder)
	assert.Nil(t, err)
	assert.NotNil(t, iter)
	defer iter.Dispose()

	count := 0
	for n := range iter.Iterate() {
		assert.NotNil(t, n)
		count++
	}
	assert.Equal(t, 5, count)

	_, err = iter.Next()
	assert.NotNil(t, err)
}

func TestIter_Finished(t *testing.T) {
	parent := nodeTree()

	iter, err := NewIterator(parent, PreOrder)
	defer iter.Dispose()
	for _ = range iter.Iterate() {}

	_, err = iter.Next()
	assert.NotNil(t, err)

	for _ = range iter.Iterate() {
		assert.Fail(t, "iteration over finished iterator")
	}
}

func TestIter_Disposed(t *testing.T) {
	parent := nodeTree()

	iter, err := NewIterator(parent, PreOrder)
	iter.Dispose()

	_, err = iter.Next()
	assert.NotNil(t, err)

	for _ = range iter.Iterate() {
		assert.Fail(t, "iteration over finished iterator")
	}
}

func TestIter_PreOrder(t *testing.T) {
	parent := nodeTree()

	iter, err := NewIterator(parent, PreOrder)
	assert.Nil(t, err)
	assert.NotNil(t, iter)
	defer iter.Dispose()

	//node := IteratorNext(iter)
	node, err := iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "parent")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "child1")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "child2")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "subchild21")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "subchild22")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.Nil(t, node)
}

func TestIter_PostOrder(t *testing.T) {
	parent := nodeTree()

	iter, err := NewIterator(parent, PostOrder)
	assert.Nil(t, err)
	assert.NotNil(t, iter)
	defer iter.Dispose()

	//node := IteratorNext(iter)
	node, err := iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "child1")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "subchild21")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "subchild22")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "child2")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "parent")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.Nil(t, node)
}

func TestIter_LevelOrder(t *testing.T) {
	parent := nodeTree()

	iter, err := NewIterator(parent, LevelOrder)
	assert.Nil(t, err)
	assert.NotNil(t, iter)
	defer iter.Dispose()

	//node := IteratorNext(iter)
	node, err := iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "parent")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "child1")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "child2")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "subchild21")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, node.InternalType, "subchild22")

	//node = IteratorNext(iter)
	node, err = iter.Next()
	assert.Nil(t, err)
	assert.Nil(t, node)
}
