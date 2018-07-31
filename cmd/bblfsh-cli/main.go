// Copyright 2018 Sourced Technologies SL
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	bblfsh "gopkg.in/bblfsh/client-go.v2"
	"gopkg.in/bblfsh/client-go.v2/tools"
	"gopkg.in/bblfsh/sdk.v1/uast"
)

func main() {
	var opts struct {
		Host     string `short:"h" long:"host" description:"Babelfish endpoint address" default:"localhost:9432"`
		Language string `short:"l" long:"language" required:"true" description:"Language"`
		Query    string `short:"q" long:"query" description:"XPath query applied to the resulting UAST"`
	}
	args, err := flags.Parse(&opts)
	if err != nil {
		fatalf("couldn't parse flags: %v", err)
	}
	filename := args[0]

	if len(args) == 0 {
		fatalf("missing file to parse")
	} else if len(args) > 1 {
		fatalf("couldn't parse more than a file at a time")
	}

	client, err := bblfsh.NewClient(opts.Host)
	if err != nil {
		fatalf("couldn't create client: %v", err)
	}

	res, err := client.NewParseRequest().
		Language(opts.Language).
		Filename(filename).
		ReadFile(filename).
		Do()
	if err != nil {
		fatalf("couldn't parse %s: %v", args[0], err)
	}

	nodes := []*uast.Node{res.UAST}
	if opts.Query != "" {
		nodes, err = tools.Filter(res.UAST, opts.Query)
		if err != nil {
			fatalf("couldn't apply query %q: %v", opts.Query, err)
		}
	}

	b, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		fatalf("couldn't encode UAST: %v", err)
	}
	fmt.Printf("%s\n", b)
}

func fatalf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
