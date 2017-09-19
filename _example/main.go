package main

import (
	"flag"
	"fmt"

	"gopkg.in/bblfsh/client-go.v1"
)

var endpoint = flag.String("e", "localhost:9432", "endpoint of the babelfish server")
var filename = flag.String("f", "", "file to parse")
var query = flag.String("q", "", "xpath expression")

func main() {
	flag.Parse()
	if *filename == "" {
		fmt.Println("filename was not provided. Use the -f flag")
		return
	}

	client, err := bblfsh.NewBblfshClient(*endpoint)
	if err != nil {
		panic(err)
	}

	res, err := client.NewParseRequest().ReadFile(*filename).Do()
	if err != nil {
		panic(err)
	}

	if *query == "" {
		fmt.Println(res.UAST)
		return

	}

	results, _ := bblfsh.Filter(res.UAST, *query)
	for i, node := range results {
		fmt.Println("-", i+1, "----------------------")
		fmt.Println(node)
	}

}
