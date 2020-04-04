package mock

import "github.com/milosgajdos/kraph/api"

type client struct{}

func NewClient() (api.Client, error) {
	return &client{}, nil
}

func (m *client) Discover() (api.API, error) {
	return NewAPI(), nil
}

func (m *client) Map(a api.API) (api.Top, error) {
	return NewTop(), nil
}
