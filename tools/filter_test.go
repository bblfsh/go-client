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

func TestFilter_InternalType(t *testing.T) {
	n := &uast.Node{
		InternalType: "a",
	}

	r, err := Filter(n, "//a")
	assert.Nil(t, err)
	assert.Len(t, r, 1)

	r, err = Filter(n, "//b")
	assert.Nil(t, err)
	assert.Len(t, r, 0)
}

func TestFilter_Token(t *testing.T) {
	n := &uast.Node{
		Token: "a",
	}

	r, err := Filter(n, "//*[@token='a']")
	assert.Nil(t, err)
	assert.Len(t, r, 1)

	r, err = Filter(n, "//*[@token='b']")
	assert.Nil(t, err)
	assert.Len(t, r, 0)
}

func TestFilter_Roles(t *testing.T) {
	n := &uast.Node{
		Roles: []uast.Role{1},
	}

	r, err := Filter(n, "//*[@roleIdentifier]")
	assert.Nil(t, err)
	assert.Len(t, r, 1)

	r, err = Filter(n, "//*[@roleQualified]")
	assert.Nil(t, err)
	assert.Len(t, r, 0)
}

func TestFilter_Properties(t *testing.T) {
	n := &uast.Node{
		Properties: map[string]string {"k2": "v1", "k1": "v2"},
	}

	r, err := Filter(n, "//*[@k1='v2']")
	assert.Nil(t, err)
	assert.Len(t, r, 1)

	r, err = Filter(n, "//*[@k2='v1']")
	assert.Nil(t, err)
	assert.Len(t, r, 1)

	r, err = Filter(n, "//*[@k3='v1']")
	assert.Nil(t, err)
	assert.Len(t, r, 0)
}

func TestFilter_NoStartPosition(t *testing.T) {
	n := &uast.Node{}

	r, err := Filter(n, "//*[@startOffset='0']")
	assert.Nil(t, err)
	assert.Len(t, r, 0)

	r, err = Filter(n, "//*[@startLine='1']")
	assert.Nil(t, err)
	assert.Len(t, r, 0)

	r, err = Filter(n, "//*[@startCol='1']")
	assert.Nil(t, err)
	assert.Len(t, r, 0)
}

func TestFilter_StartPosition(t *testing.T) {
	n := &uast.Node{
		StartPosition: &uast.Position{0, 1, 1},
	}

	r, err := Filter(n, "//*[@startOffset='0']")
	assert.Nil(t, err)
	assert.Len(t, r, 1)

	r, err = Filter(n, "//*[@startLine='1']")
	assert.Nil(t, err)
	assert.Len(t, r, 1)

	r, err = Filter(n, "//*[@startCol='1']")
	assert.Nil(t, err)
	assert.Len(t, r, 1)
}

func TestFilter_NoEndPosition(t *testing.T) {
	n := &uast.Node{}

	r, err := Filter(n, "//*[@endOffset='0']")
	assert.Nil(t, err)
	assert.Len(t, r, 0)

	r, err = Filter(n, "//*[@endLine='1']")
	assert.Nil(t, err)
	assert.Len(t, r, 0)

	r, err = Filter(n, "//*[@endCol='1']")
	assert.Nil(t, err)
	assert.Len(t, r, 0)
}

func TestFilter_EndPosition(t *testing.T) {
	n := &uast.Node{
		EndPosition: &uast.Position{0, 1, 1},
	}

	r, err := Filter(n, "//*[@endOffset='0']")
	assert.Nil(t, err)
	assert.Len(t, r, 1)

	r, err = Filter(n, "//*[@endLine='1']")
	assert.Nil(t, err)
	assert.Len(t, r, 1)

	r, err = Filter(n, "//*[@endCol='1']")
	assert.Nil(t, err)
	assert.Len(t, r, 1)
}
