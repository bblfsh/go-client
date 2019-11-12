package bblfsh

import (
	"context"

	protocol2 "github.com/bblfsh/sdk/v3/protocol"
	"google.golang.org/grpc"
)

/*
	MultipleDriverClient could be useful during language-specific parsings on the large scale

	Scenario: we need to parse a ton of go and python files inside the k8s environment, to save time we need to perform this
	parses on the large scale of go- and python-driver containers/instances etc.

	Solution:
	- run two separate Deployments of go- and python-driver container pods
	- provide Horizontal Autoscalers for both of Deployments
	- provide Services with LoadBalancer type
	- during client initialization provide Services endpoints configuration, this will create two language-oriented connections
	that will be responsible for sending parse request only to a dedicated Service, that will load-balance this request
	between underlying language driver pods

	Examples of endpoint formats:
	- localhost:9432 - casual example there's only one driver or bblfshd server
	- python=localhost:9432,go=localhost:9432 - coma-separated mapping in format language=address
	- %s-driver.bblfsh.svc.example.com - DNS template based on the language
*/

// MultipleDriverClient is a DriverClient implementation, contains connection getter and a map[language]connection
type MultipleDriverClient struct {
	getConn   getConnFunc
	Languages map[string]*grpc.ClientConn
}

// MultipleDriverHostClient is a DriverHostClient implementation, currently does almost nothing
type MultipleDriverHostClient struct{}

// NewMultipleDriverClient is a MultipleDriverClient constructor
func NewMultipleDriverClient(getConn getConnFunc) *MultipleDriverClient {
	return &MultipleDriverClient{
		getConn:   getConn,
		Languages: make(map[string]*grpc.ClientConn),
	}
}

// Parse gets connection from a given map, or creates a new connection, then inits driver client and performs Parse
func (c *MultipleDriverClient) Parse(
	ctx context.Context,
	in *protocol2.ParseRequest,
	opts ...grpc.CallOption) (*protocol2.ParseResponse, error) {
	lang := in.Language

	conn, ok := c.Languages[lang]
	if !ok {
		gConn, err := c.getConn(ctx, lang)
		if err != nil {
			return nil, err
		}
		conn = gConn
	}

	return protocol2.NewDriverClient(conn).Parse(ctx, in, opts...)
}

func (hc *MultipleDriverHostClient) ServerVersion(
	ctx context.Context,
	in *protocol2.VersionRequest,
	opts ...grpc.CallOption) (*protocol2.VersionResponse, error) {
	return nil, ErrNotImplemented.New()
}

func (hc *MultipleDriverHostClient) SupportedLanguages(
	ctx context.Context,
	in *protocol2.SupportedLanguagesRequest,
	opts ...grpc.CallOption) (*protocol2.SupportedLanguagesResponse, error) {
	return nil, ErrNotImplemented.New()
}
