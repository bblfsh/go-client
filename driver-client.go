package bblfsh

import (
	"context"

	protocol2 "github.com/bblfsh/sdk/v3/protocol"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
	multipleDriverClient could be useful during language-specific parsings on the large scale

	Examples of endpoint formats:
	- localhost:9432 - casual example there's only one driver or bblfshd server
	- python=localhost:9432,go=localhost:9432 - coma-separated mapping in format language=address
	- %s-driver.bblfsh.svc.example.com - DNS template based on the language
*/

// multipleDriverClient is a DriverClient implementation, contains connection getter and a map[language]connection
type multipleDriverClient struct {
	getConn        getConnFunc
	langConnDriver langConnDriver
}

type langConnDriver map[string]*connDriver

type connDriver struct {
	conn   *grpc.ClientConn
	driver protocol2.DriverClient
}

// multipleDriverHostClient is a DriverHostClient implementation, currently does almost nothing
type multipleDriverHostClient struct{}

// newMultipleDriverClient is a multipleDriverClient constructor
func newMultipleDriverClient(getConn getConnFunc) *multipleDriverClient {
	return &multipleDriverClient{
		getConn:        getConn,
		langConnDriver: make(langConnDriver),
	}
}

// Parse gets connection from a given map, or creates a new connection, then inits driver client and performs Parse
func (c *multipleDriverClient) Parse(
	ctx context.Context,
	in *protocol2.ParseRequest,
	opts ...grpc.CallOption) (*protocol2.ParseResponse, error) {
	lang := in.Language

	connD, ok := c.langConnDriver[lang]
	if !ok {
		gConn, err := c.getConn(ctx, lang)
		if err != nil {
			return nil, err
		}
		connD = &connDriver{
			conn: gConn,
		}
	}

	if connD.driver == nil {
		connD.driver = protocol2.NewDriverClient(connD.conn)
	}
	c.langConnDriver[lang] = connD

	return connD.driver.Parse(ctx, in, opts...)
}

func (c *multipleDriverClient) Close() error {
	var lastErr error
	for k, v := range c.langConnDriver {
		if err := v.conn.Close(); err != nil {
			lastErr = err
		}
		delete(c.langConnDriver, k)
	}
	c.langConnDriver = make(langConnDriver)
	return lastErr
}

func (hc *multipleDriverHostClient) ServerVersion(
	ctx context.Context,
	in *protocol2.VersionRequest,
	opts ...grpc.CallOption) (*protocol2.VersionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ServerVersion is not implemented")
}

func (hc *multipleDriverHostClient) SupportedLanguages(
	ctx context.Context,
	in *protocol2.SupportedLanguagesRequest,
	opts ...grpc.CallOption) (*protocol2.SupportedLanguagesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "SupportedLanguages is not implemented")
}
