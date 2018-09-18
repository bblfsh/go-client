package tools

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"gopkg.in/bblfsh/sdk.v2/protocol/v1"
	"gopkg.in/bblfsh/sdk.v2/uast/yaml"
)

const fixture = `./testdata/json.go.sem.uast`

func BenchmarkXPathV1(b *testing.B) {
	data, err := ioutil.ReadFile(fixture)
	require.NoError(b, err)
	node, err := uastyml.Unmarshal(data)
	require.NoError(b, err)

	node1, err := uast1.ToNode(node)
	require.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		arr, err := Filter(node1, `//Identifier`)
		if err != nil {
			b.Fatal(err)
		} else if len(arr) != 2292 {
			b.Fatal("wrong result:", len(arr))
		}
	}
}

func BenchmarkXPathV2(b *testing.B) {
	data, err := ioutil.ReadFile(fixture)
	require.NoError(b, err)
	node, err := uastyml.Unmarshal(data)
	require.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		it, err := FilterXPath(node, `//uast:Identifier`)
		if err != nil {
			b.Fatal(err)
		}
		cnt := 0
		for it.Next() {
			cnt++
			_ = it.Node()
		}
		if cnt != 2292 {
			b.Fatal("wrong result:", cnt)
		}
	}
}
