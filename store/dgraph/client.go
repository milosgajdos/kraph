package dgraph

import (
	dgo "github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

// Client is dgraph client
type Client struct {
	*dgo.Dgraph
	conn *grpc.ClientConn
}

// NewClient creates new dgraph client and returns it
func NewClient(target string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, err
	}

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	return &Client{
		Dgraph: dg,
		conn:   conn,
	}, nil
}

// Close closes dgraph connection
func (c *Client) Close() error {
	return c.conn.Close()
}
