package bblfsh

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bblfsh/sdk/v3/driver"
	"github.com/bblfsh/sdk/v3/driver/manifest"
	protocol2 "github.com/bblfsh/sdk/v3/protocol"
	protocol1 "gopkg.in/bblfsh/sdk.v1/protocol"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"gopkg.in/src-d/go-errors.v1"
)

const (
	// defaultConnTimeout is a default connection timeout to bblfshd.
	defaultConnTimeout = 5 * time.Second

	// keepalivePingInterval is a duration after this if the client doesn't see any activity it
	// pings the server to see if the transport is still alive.
	keepalivePingInterval = 2 * time.Minute

	// keepalivePingWithoutStream is a boolean flag.
	// If true, client sends keepalive pings even with no active RPCs.
	keepalivePingWithoutStream = true
)

type ConnFunc func(ctx context.Context, language string) (*grpc.ClientConn, error)

// Client holds the public client API to interact with the bblfsh daemon.
type Client struct {
	closer  io.Closer
	driver2 protocol2.DriverClient
	driver  driver.Driver
}

// NewClientContext returns a new bblfsh client given a bblfshd endpoint.
func NewClientContext(ctx context.Context, endpoint string, options ...grpc.DialOption) (*Client, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                keepalivePingInterval,
			PermitWithoutStream: keepalivePingWithoutStream,
		}),
	}
	opts = append(opts, protocol2.DialOptions()...)
	// user-defined options should go last
	// this allows to override any default option
	opts = append(opts, options...)

	switch {
	case strings.Contains(endpoint, ","):
		endpoints, err := parseEndpoints(endpoint)
		if err != nil {
			return nil, err
		}
		return NewClientWithConnectionsContext(func(ctx context.Context, lang string) (*grpc.ClientConn, error) {
			e, ok := endpoints[lang]
			if !ok {
				return nil, &driver.ErrMissingDriver{Language: lang}
			}
			conn, err := grpc.DialContext(ctx, e, opts...)
			if err != nil {
				return nil, err
			}

			return conn, nil
		})
	case strings.Contains(endpoint, "%s"):
		return NewClientWithConnectionsContext(func(ctx context.Context, lang string) (*grpc.ClientConn, error) {
			conn, err := grpc.DialContext(ctx, fmt.Sprintf(endpoint, lang), opts...)
			if err != nil {
				return nil, err
			}

			return conn, nil
		})
	default:
		conn, err := grpc.DialContext(ctx, endpoint, opts...)
		if err != nil {
			return nil, err
		}
		return NewClientWithConnectionContext(ctx, conn)
	}
}

func NewClientWithConnectionsContext(getConn ConnFunc) (*Client, error) {
	dc := newMultipleDriverClient(getConn)

	return &Client{
		closer:  dc,
		driver2: dc,
		driver:  protocol2.DriverFromClient(dc, &multipleDriverHostClient{}),
	}, nil
}

func parseEndpoints(endpoints string) (map[string]string, error) {
	result := make(map[string]string)
	pairs := strings.Split(endpoints, ",")
	for _, p := range pairs {
		vals := strings.Split(p, "=")
		if len(vals) != 2 {
			return nil, fmt.Errorf("formatting is broken in section: %q", p)
		}
		result[vals[0]] = vals[1]
	}
	return result, nil
}

// NewClient is the same as NewClientContext, but assumes a default timeout for the connection.
//
// Deprecated: use NewClientContext instead
func NewClient(endpoint string, options ...grpc.DialOption) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnTimeout)
	defer cancel()

	return NewClientContext(ctx, endpoint, options...)
}

// NewClientWithConnection returns a new bblfsh client given a grpc connection.
func NewClientWithConnection(conn *grpc.ClientConn) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnTimeout)
	defer cancel()

	return NewClientWithConnectionContext(ctx, conn)
}

func isServiceNotSupported(err error) bool {
	if err == nil {
		return false
	}
	return status.Code(err) == codes.Unimplemented
}

// NewClientWithConnectionContext returns a new bblfsh client given a grpc connection.
func NewClientWithConnectionContext(ctx context.Context, conn *grpc.ClientConn) (*Client, error) {
	host := protocol2.NewDriverHostClient(conn)
	_, err := host.ServerVersion(ctx, &protocol2.VersionRequest{})
	if err == nil {
		// supports v2
		return &Client{
			closer:  conn,
			driver2: protocol2.NewDriverClient(conn),
			driver:  protocol2.AsDriver(conn),
		}, nil
	} else if !isServiceNotSupported(err) {
		return nil, err
	}
	s1 := protocol1.NewProtocolServiceClient(conn)
	return &Client{
		closer:  conn,
		driver2: protocol2.NewDriverClient(conn),
		driver: &driverPartialV2{
			// use only Parse from v2
			Driver: protocol2.AsDriver(conn),
			// use v1 for version and supported languages
			service1: s1,
		},
	}, nil
}

type driverPartialV2 struct {
	driver.Driver
	service1 protocol1.ProtocolServiceClient
}

// Version implements a driver.Host using v1 protocol.
func (d *driverPartialV2) Version(ctx context.Context) (driver.Version, error) {
	resp, err := d.service1.Version(ctx, &protocol1.VersionRequest{})
	if err != nil {
		return driver.Version{}, err
	} else if resp.Status != protocol1.Ok {
		return driver.Version{}, errorStrings(resp.Errors)
	}
	return driver.Version{
		Version: resp.Version,
		Build:   resp.Build,
	}, nil
}

// Languages implements a driver.Host using v1 protocol.
func (d *driverPartialV2) Languages(ctx context.Context) ([]manifest.Manifest, error) {
	resp, err := d.service1.SupportedLanguages(ctx, &protocol1.SupportedLanguagesRequest{})
	if err != nil {
		return nil, err
	} else if resp.Status != protocol1.Ok {
		return nil, errorStrings(resp.Errors)
	}
	out := make([]manifest.Manifest, 0, len(resp.Languages))
	for _, m := range resp.Languages {
		dm := manifest.Manifest{
			Name:     m.Name,
			Language: m.Language,
			Version:  m.Version,
			Status:   manifest.DevelopmentStatus(m.Status),
			Features: make([]manifest.Feature, 0, len(m.Features)),
		}
		for _, f := range m.Features {
			dm.Features = append(dm.Features, manifest.Feature(f))
		}
		out = append(out, dm)
	}
	return out, nil
}

// NewParseRequest is a parsing request to get the UAST.
func (c *Client) NewParseRequest() *ParseRequest {
	return &ParseRequest{ctx: context.Background(), client: c}
}

// NewVersionRequest is a parsing request to get the version of the server.
func (c *Client) NewVersionRequest() *VersionRequest {
	return &VersionRequest{ctx: context.Background(), client: c}
}

// NewSupportedLanguagesRequest is a parsing request to get the supported languages.
func (c *Client) NewSupportedLanguagesRequest() *SupportedLanguagesRequest {
	return &SupportedLanguagesRequest{ctx: context.Background(), client: c}
}

func (c *Client) GetConn() (*grpc.ClientConn, error) {
	if conn, ok := c.closer.(*grpc.ClientConn); ok {
		return conn, nil
	}
	return nil, errors.NewKind("multiple connections").New()
}

func (c *Client) Close() error {
	if c.closer == nil {
		return nil
	}
	return c.closer.Close()
}
