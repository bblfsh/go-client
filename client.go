package bblfsh

import (
	"time"

	"gopkg.in/bblfsh/sdk.v0/protocol"
	"google.golang.org/grpc"
)

// BblfshClient holds the public client API to interact with the babelfish server.
type BblfshClient struct {
	client protocol.ProtocolServiceClient
}

// NewBblfshClient returns a new babelfish client given a server endpoint
func NewBblfshClient(endpoint string) (*BblfshClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTimeout(time.Second*2), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &BblfshClient{
		client: protocol.NewProtocolServiceClient(conn),
	}, nil
}

func (c *BblfshClient) NewParseRequest() *ParseRequest {
	return &ParseRequest{client: c}
}
