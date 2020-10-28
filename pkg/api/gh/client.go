package gh

import (
	"context"

	"github.com/google/go-github/v32/github"
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/errors"
)

const (
	source = "gh"
)

type repo struct {
	// ctx is client context
	ctx context.Context
	// client is GitHub client
	client *github.Client
}

// NewRepoClient creates a new GitHub repository client and returns it
func NewRepoClient(ctx context.Context, client *github.Client) *repo {
	return &repo{
		ctx:    ctx,
		client: client,
	}
}

// Discover discovers GH reposity API and returns them.
func (g *repo) Discover() (api.API, error) {
	api := NewRepoAPI(source)

	return api, errors.ErrNotImplemented
}

// Map builds a map of GH API resources and returns their topology.
// It returns error if any of the API calls fails with error.
func (g *repo) Map(a api.API) (api.Top, error) {
	return nil, errors.ErrNotImplemented
}
