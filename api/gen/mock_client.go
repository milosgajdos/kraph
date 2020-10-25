package gen

import "github.com/milosgajdos/kraph/api"

type client struct {
	resPath string
	objPath string
}

func NewMockClient(resPath, objPath string) (api.Client, error) {
	return &client{
		resPath: resPath,
		objPath: objPath,
	}, nil
}

func (c *client) Discover() (api.API, error) {
	return NewMockAPI(c.resPath)
}

func (c *client) Map(a api.API) (api.Top, error) {
	return NewMockTop(c.objPath)
}
