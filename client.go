package bblfsh

import (
	"time"

	"github.com/bblfsh/sdk/protocol"
	"google.golang.org/grpc"
)

type BblfshClient struct {
	client protocol.ProtocolServiceClient
}

func NewBblfshClient(endpoint string) (*BblfshClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTimeout(time.Second*2), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &BblfshClient{
		client: protocol.NewProtocolServiceClient(conn),
	}, nil
}

func (c *BblfshClient) Parse() *ParseRequest {
	return &ParseRequest{client: c}
}
