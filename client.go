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

// import "github.com/bblfsh/go-client"

// func main() {
// 	client, _ := bblfsh.NewBblfshClient("0.0.0.0:9413")
// 	root, _ := client.Parse().ReadFile("hola").Do()
// 	// 	root, _ := client.Parse().Content("C = A+B").Language("Python").Do()

// 	for _, node := range bblfsh.Find(root, "//Module/*") {
// 		fmt.Println(node)
// 	}
// }
