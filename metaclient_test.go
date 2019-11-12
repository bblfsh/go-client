package bblfsh

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const timeout = 5 * time.Second

func TestMetaClient(t *testing.T) {
	mc, err := NewMetaClient("_testdata/endpoints.yml")
	require.NoError(t, err)

	for _, mct := range metaClientTests {
		mct := mct
		t.Run(mct.name, func(t *testing.T) {
			mct.test(t, mc)
		})
	}
}

var metaClientTests = []struct {
	name string
	test func(t *testing.T, mc *MetaClient)
}{
	{name: "GetClientContextAndParse", test: testGetClientContextAndParse},
	{name: "GetClientContextTwiceAndParse", test: testGetClientContextTwiceAndParse},
	{name: "GetNotSupportedClient", test: testGetNotSupportedClient},
}

func testGetClientContextAndParse(t *testing.T, mc *MetaClient) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	getClientContextAndParse(t, ctx, mc, "python", "import foo")
	getClientContextAndParse(t, ctx, mc, "go", "package main")
}

func testGetClientContextTwiceAndParse(t *testing.T, mc *MetaClient) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	getClientContextAndParse(t, ctx, mc, "python", "import foo")
	getClientContextAndParse(t, ctx, mc, "python", "import foo")
}

func testGetNotSupportedClient(t *testing.T, mc *MetaClient) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := mc.GetClientContext(ctx, "lang")
	if err == context.DeadlineExceeded {
		t.Skip("bblfshd is not running")
	}
	require.True(t, ErrLangNotSupported.Is(err))
}

func getClientContextAndParse(
	t *testing.T,
	ctx context.Context,
	mc *MetaClient,
	language,
	content string) {
	c, err := mc.GetClientContext(ctx, language)
	if err == context.DeadlineExceeded {
		t.Skip("bblfshd is not running")
	}
	require.NoError(t, err)

	res, err := c.NewParseRequest().Mode(Native).Language(language).Content(content).Do()
	require.NoError(t, err)

	require.Len(t, res.Errors, 0)
	require.NotNil(t, res)
}
