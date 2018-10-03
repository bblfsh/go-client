# client-go [![GoDoc](https://godoc.org/gopkg.in/bblfsh/client-go.v3?status.svg)](https://godoc.org/gopkg.in/bblfsh/client-go.v3) [![Build Status](https://travis-ci.org/bblfsh/client-go.svg?branch=master)](https://travis-ci.org/bblfsh/client-go) [![Build status](https://ci.appveyor.com/api/projects/status/github/bblfsh/client-go?svg=true)](https://ci.appveyor.com/project/mcuadros/client-go) [![codecov](https://codecov.io/gh/bblfsh/client-go/branch/master/graph/badge.svg)](https://codecov.io/gh/bblfsh/client-go)

[Babelfish](https://doc.bblf.sh) Go client library provides functionality to both
connecting to the Babelfish server for parsing code
(obtaining an [UAST](https://doc.bblf.sh/uast/specification.html) as a result)
and for analysing UASTs with the functionality provided by [libuast](https://github.com/bblfsh/libuast).

## Installation

The recommended way to install *client-go* is:

```sh
go get -u gopkg.in/bblfsh/client-go.v3/...
```

## Example

This small example illustrates how to retrieve the [UAST](https://doc.bblf.sh/uast/specification.html) from a small Python script.

If you don't have a bblfsh server installed, please read the [getting started](https://doc.bblf.sh/using-babelfish/getting-started.html) guide, to learn more about how to use and deploy a bblfsh server. 

Go to the [quick start](https://github.com/bblfsh/bblfshd#quick-start) to discover how to run Babelfish with Docker.

```go
package main

import (
	"fmt"

	bblfsh "gopkg.in/bblfsh/client-go.v3"
	"gopkg.in/bblfsh/client-go.v3/tools"
	"gopkg.in/bblfsh/sdk.v2/uast/nodes"
)

func main() {
	client, err := bblfsh.NewClient("0.0.0.0:9432")
	if err != nil {
		panic(err)
	}

	python := "import foo"

	req := client.NewParseRequest().Language("python").Content(python)
	node, _, err := req.UAST()
	if err != nil {
		panic(err)
	}

	it, err := tools.Filter(node, "//*[@role='Import']")
	if err != nil {
		panic(err)
	}

	for it.Next() {
		nodes.WalkPreOrderExt(it.Node(), func(e nodes.External) bool {
			switch val := e.(type) {
			case nodes.Object:
				fmt.Printf("Object.Kind: %v\tObject.Keys: %v\n", val.Kind(), val.Keys())
			case nodes.Array:
				fmt.Printf("Array.Kind: %v\tArray.Values:", val.Kind())
				for i := 0; i < val.Size(); i++ {
					fmt.Printf(" %v", val.ValueAt(i))
				}
				fmt.Println()
			case nodes.Value:
				fmt.Printf("Value.Kind: %v\tValue.Value: %v\n", val.Kind(), val.Value())
			default:
				return false
			}

			return true
		})

	}
}
```

```
Object.Kind: Object	Object.Keys: [@pos @role @token @type names]

Object.Kind: Object	Object.Keys: [@type start]

Value.Kind: String	Value.Value: uast:Positions

Object.Kind: Object	Object.Keys: [@type col line offset]

Value.Kind: String	Value.Value: uast:Position

Value.Kind: Uint	Value.Value: 1

Value.Kind: Uint	Value.Value: 1

Value.Kind: Uint	Value.Value: 0

Array.Kind: Array	Array.Values: Import Declaration Statement
Value.Kind: String	Value.Value: Import

Value.Kind: String	Value.Value: Declaration

Value.Kind: String	Value.Value: Statement

Value.Kind: String	Value.Value: import

Value.Kind: String	Value.Value: Import

Object.Kind: Object	Object.Keys: [@role @type name_list]

Array.Kind: Array	Array.Values: Import Pathname Identifier Incomplete
Value.Kind: String	Value.Value: Import

Value.Kind: String	Value.Value: Pathname

Value.Kind: String	Value.Value: Identifier

Value.Kind: String	Value.Value: Incomplete

Value.Kind: String	Value.Value: ImportFrom.names

Array.Kind: Array	Array.Values: map[All:false Names:<nil> Path:map[@pos:map[@type:uast:Positions] @type:uast:Alias Name:map[@type:uast:Identifier Name:foo] Node:map[]] Target:<nil> @type:uast:RuntimeImport]
Object.Kind: Object	Object.Keys: [@type All Names Path Target]

...
```

Please read the [Babelfish clients](https://doc.bblf.sh/using-babelfish/clients.html) guide section to learn more about babelfish clients and their query language.

## License

Apache License 2.0, see [LICENSE](LICENSE)
