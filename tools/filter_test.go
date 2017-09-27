package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/bblfsh/sdk.v1/uast"
)

func TestFilter(t *testing.T) {
	n := &uast.Node{}

	r, err := Filter(n, "")
	assert.Len(t, r, 0)
	assert.Nil(t, err)
}

func TestFilter_All(t *testing.T) {
	n := &uast.Node{}

	_, err := Filter(n, "//*")
	assert.Nil(t, err)
}
