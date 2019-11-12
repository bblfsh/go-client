package bblfsh

import (
	"context"
	"io/ioutil"

	"google.golang.org/grpc"
	"gopkg.in/src-d/go-errors.v1"
	"gopkg.in/yaml.v2"
)

/*
	MetaClient could be useful when dealing with language-specific parsings on the large scale

	Scenario: we need to parse a ton of go and python files inside the k8s environment, to save time we need to perform this
	parses on the large scale of go- and python-driver containers/instances etc.

	Solution:
	- run two separate Deployments of go- and python-driver container pods
	- provide Horizontal Autoscalers for both of Deployments
	- provide Services with LoadBalancer type
	- during MetaClient initialization provide config with endpoints to Services, this will create two language-oriented clients
	that will be responsible for sending parse request only to a dedicated Service, that will load-balance this request
	between underlying language driver pods
*/

// ErrLangNotSupported is the error that must be returned if requested language driver is not supported/available
var ErrLangNotSupported = errors.NewKind("language %q is not supported in current configuration")

// MetaClient represents an entity that contains:
// 1) MetaClientConfig configuration
// 2) LangClients - map of already initialized clients
type MetaClient struct {
	clients LangClients
	config  MetaClientConfig
}

// LangClients is a map, where key represents language name, value - corresponding client
type LangClients map[string]*Client

// MetaClientConfig is a map, where key represents language name, value - corresponding driver's API endpoint
type MetaClientConfig map[string]string

// NewMetaClient is a MetaClient constructor, used to parse and save MetaClientConfig from a given filepath
func NewMetaClient(configPath string) (*MetaClient, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := make(MetaClientConfig)
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return &MetaClient{
		clients: make(LangClients),
		config:  config,
	}, nil
}

// GetClientContext returns corresponding client for requested language or creates the new one if it's not initialized
func (mc *MetaClient) GetClientContext(
	ctx context.Context,
	language string,
	options ...grpc.DialOption) (*Client, error) {
	if c, ok := mc.clients[language]; ok {
		return c, nil
	}
	if endpoint, ok := mc.config[language]; ok {
		c, err := NewClientContext(ctx, endpoint, options...)
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	return nil, ErrLangNotSupported.New(language)
}
