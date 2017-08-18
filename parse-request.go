package bblfsh

import (
	"context"
	"io/ioutil"
	"path"

	"github.com/bblfsh/sdk/protocol"
)

type ParseRequest struct {
	internal protocol.ParseRequest
	client   *BblfshClient
	err      error
}

func (req *ParseRequest) Language(lan string) *ParseRequest {
	req.internal.Language = lan
	return req
}

func (req *ParseRequest) ReadFile(filename string) *ParseRequest {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		req.err = err
		return req
	}
	req.internal.Content = string(data)
	req.internal.Filename = path.Base(filename)
	return req
}

func (req *ParseRequest) Content(content string) *ParseRequest {
	req.internal.Content = content
	return req
}

func (req *ParseRequest) Filename(filename string) *ParseRequest {
	req.internal.Filename = filename
	return req
}

func (req *ParseRequest) Encoding(encoding protocol.Encoding) *ParseRequest {
	req.internal.Encoding = encoding
	return req
}

func (req *ParseRequest) Do() (*protocol.ParseResponse, error) {
	return req.client.client.Parse(context.TODO(), &req.internal)
}

func (req *ParseRequest) DoWithContext() (*protocol.ParseResponse, error) {
	return req.client.client.Parse(context.TODO(), &req.internal)
}
