package bblfsh

import (
	"context"
	"errors"
	"io/ioutil"
	"path"

	"gopkg.in/bblfsh/sdk.v1/protocol"
)

// ParseRequest is a placeholder for the parse requests performed by the library
type ParseRequest struct {
	internal protocol.ParseRequest
	client   *BblfshClient
	err      error
}

// Language sets the language of the given source file to parse.
func (req *ParseRequest) Language(lan string) *ParseRequest {
	req.internal.Language = lan
	return req
}

// ReadFile loads a file given a local path and sets the content and the filename of the request.
func (req *ParseRequest) ReadFile(filepath string) *ParseRequest {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		req.err = err
	} else {
		req.internal.Content = string(data)
		req.internal.Filename = path.Base(filepath)
	}
	return req
}

// Content sets the content of the parse request. It should be the source code that wants to be parsed.
func (req *ParseRequest) Content(content string) *ParseRequest {
	req.internal.Content = content
	return req
}

// Filename sets the filename of the content.
func (req *ParseRequest) Filename(filename string) *ParseRequest {
	req.internal.Filename = filename
	return req
}

// Encoding sets the text encoding of the content.
func (req *ParseRequest) Encoding(encoding protocol.Encoding) *ParseRequest {
	req.internal.Encoding = encoding
	return req
}

// Do performs the actual parsing by serializaing the request,
// sending it to the babelfish and waiting for the response (UAST tree).
func (req *ParseRequest) Do() (*protocol.ParseResponse, error) {
	return req.DoWithContext(context.Background())
}

// DoWithContext does the same as Do(), but sopporting cancellation by the use of
// Go contexts.
func (req *ParseRequest) DoWithContext(ctx context.Context) (*protocol.ParseResponse, error) {
	if req.err != nil {
		return nil, req.err
	}
	if req.client == nil {
		return nil, errors.New("request is clientless")
	}
	return req.client.client.Parse(ctx, &req.internal)
}
