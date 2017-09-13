package bblfsh

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cli, err := NewBblfshClient("localhost:9432")
	require.Nil(t, err)

	python := "import foo"

	res, err := cli.NewParseRequest().Language("python").Content(python).Do()
	require.Nil(t, err)

	fmt.Println(res.Errors)
	fmt.Println(res.UAST)

	require.Equal(t, len(res.Errors), 0)
	require.NotNil(t, res.UAST)

}
