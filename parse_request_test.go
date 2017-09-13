package bblfsh

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/bblfsh/sdk.v0/protocol"
)

func TestParseRequestConfiguration(t *testing.T) {
	req := ParseRequest{}
	req.Filename("file.py").Language("python").Encoding(protocol.UTF8).Content("a=b+c")

	require.Equal(t, "file.py", req.internal.Filename)
	require.Equal(t, "python", req.internal.Language)
	require.Equal(t, protocol.UTF8, req.internal.Encoding)
	require.Equal(t, "a=b+c", req.internal.Content)
}

func TestParseRequestClientlessError(t *testing.T) {
	req := ParseRequest{}
	_, err := req.Content("a=b+c").Language("python").Do()
	require.Errorf(t, err, "request is clientless")
}

func TestParseRequestReadFileError(t *testing.T) {
	req := ParseRequest{}
	_, err := req.ReadFile("NO_EXISTS").Do()
	require.Errorf(t, err, "open NO_EXISTS: no such file or directory")
}
