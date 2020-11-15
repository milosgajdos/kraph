package star

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/errors"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/query"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

const (
	source = "GitHubStars"
	// GitHub API version
	version = "v3"
	// topicRel is topic relation
	topicRel = "topic"
	// langRel is language relation
	langRel = "lang"
)

type client struct {
	// ctx is client context
	ctx context.Context
	// gh is GitHub gh
	gh *github.Client
	// source is API source
	src api.Source
	// opts are client options
	opts Options
}

// NewClient creates a new GitHub star repository client and returns it
func NewClient(ctx context.Context, gh *github.Client, opts ...Option) *client {
	copts := Options{}
	for _, apply := range opts {
		apply(&copts)
	}

	return &client{
		ctx:  ctx,
		gh:   gh,
		src:  NewSource(source),
		opts: copts,
	}
}

// Discover discovers GH stars API and returns them.
func (g *client) Discover() (api.API, error) {
	a := NewRepoAPI(g.src)

	repo := NewResource("repo", "starred", "repos", version, true, api.Options{Metadata: metadata.New()})
	if err := a.Add(repo, api.AddOptions{}); err != nil {
		return nil, err
	}

	// NOTE: we are "cheating" here as neither lang nor repo are API resources
	// Technically, lang is and it has a dedicated API endpoint, but we're keeping things simple here for now

	lang := NewResource("lang", "repo", "langs", version, false, api.Options{Metadata: metadata.New()})
	if err := a.Add(lang, api.AddOptions{}); err != nil {
		return nil, err
	}

	topic := NewResource("topic", "repo", "topics", version, false, api.Options{Metadata: metadata.New()})
	if err := a.Add(topic, api.AddOptions{}); err != nil {
		return nil, err
	}

	return a, nil
}

func (g *client) mapObjects(top api.Top, to uuid.UID, res api.Resource, rel string, names []string) ([]*Object, error) {
	objects := make([]*Object, len(names))

	for i, name := range names {
		linkMeta := metadata.New()
		linkMeta.Set("relation", rel)

		// NOTE: we are setting uid to the name of the object
		// this is so we avoid duplicating topics with the same name
		uid := uuid.NewFromString(name)

		links := []link{
			{uid: to, opts: api.LinkOptions{Metadata: linkMeta}},
		}

		m := metadata.New()
		objects[i] = NewObject(uid, name, api.NsGlobal, res, api.Options{Metadata: m}, links)

		if err := top.Add(objects[i], api.AddOptions{MergeLinks: true}); err != nil {
			return nil, err
		}
	}

	return objects, nil
}

func getResource(a api.API, name, version string) (api.Resource, error) {
	q := query.Build().
		Name(name, query.StringEqFunc(name)).
		Version(version, query.StringEqFunc(version))

	rx, err := a.Get(q)
	if err != nil {
		return nil, err
	}

	if len(rx) != 1 {
		return nil, fmt.Errorf("expected single %s resource, got: %d", name, len(rx))
	}

	return rx[0], nil
}

// Map builds a map of GH API starred repos and returns their topology.
// It returns error if any of the API calls fails with error.
func (g *client) Map(a api.API) (api.Top, error) {
	top := NewTop(a)

	// NOTE: all API resources have the same version
	// but their version might not be the same in the future

	repoRes, err := getResource(a, "repo", version)
	if err != nil {
		return nil, err
	}

	topicRes, err := getResource(a, "topic", version)
	if err != nil {
		return nil, err
	}

	langRes, err := getResource(a, "lang", version)
	if err != nil {
		return nil, err
	}

	opts := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: g.opts.Paging},
	}

	// NOTE: we are only iterating over the repos resources
	// as topics and langs are merely adjacent nodes of repos objects

	// TODO: make this faster
	for {
		repos, resp, err := g.gh.Activity.ListStarred(g.ctx, g.opts.User, opts)
		if err != nil {
			return nil, fmt.Errorf("failed listing repos: %v", err)
		}

		for _, repo := range repos {
			m := metadata.New()
			m.Set("starred_at", repo.StarredAt)
			m.Set("git_url", repo.Repository.GetURL)

			var owner string
			if repo.Repository.Organization != nil {
				owner = *repo.Repository.Organization.Login
			} else {
				owner = *repo.Repository.Owner.Login
			}

			uid := uuid.NewFromString(*repo.Repository.NodeID)

			topicObjects, err := g.mapObjects(top, uid, topicRes, topicRel, repo.Repository.Topics)
			if err != nil {
				return nil, err
			}

			var langObjects []*Object
			if repo.Repository.Language != nil {
				var err error
				langObjects, err = g.mapObjects(top, uid, langRes, langRel, []string{*repo.Repository.Language})
				if err != nil {
					return nil, err
				}
			}

			var repoLinks []link

			for _, o := range topicObjects {
				linkMeta := metadata.New()
				linkMeta.Set("relation", topicRel)

				l := link{uid: o.UID(), opts: api.LinkOptions{Metadata: linkMeta}}
				repoLinks = append(repoLinks, l)
			}

			for _, o := range langObjects {
				linkMeta := metadata.New()
				linkMeta.Set("relation", langRel)

				l := link{uid: o.UID(), opts: api.LinkOptions{Metadata: linkMeta}}
				repoLinks = append(repoLinks, l)
			}

			obj := NewObject(uid, *repo.Repository.Name, owner, repoRes, api.Options{Metadata: m}, repoLinks)

			if err := top.Add(obj, api.AddOptions{MergeLinks: true}); err != nil {
				return nil, err
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return nil, errors.ErrNotImplemented
}
