package star

import (
	"context"
	"strings"
	"sync"

	"github.com/google/go-github/v32/github"
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/query"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

const (
	// source is the API source
	source = "GitHubStars"
	// version is GitHub API version
	version = "v3"
	// topicRel is topic relation
	topicRel = "topic"
	// langRel is language relation
	langRel = "lang"
	// resKind is api.Resource kind
	resKind = "starred"
	// workers is the default number of workers
	workers = 5
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

	if copts.Workers <= 0 {
		copts.Workers = workers
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

	repo := NewResource("repo", resKind, "repos", version, true, api.Options{Metadata: metadata.New()})
	if err := a.Add(repo, api.AddOptions{}); err != nil {
		return nil, err
	}

	// NOTE: we are "cheating" here as neither lang nor repo are actual API resources per se
	// Technically, lang is and it has a dedicated API endpoint, but for now we're keeping things simple

	lang := NewResource("lang", resKind, "langs", version, false, api.Options{Metadata: metadata.New()})
	if err := a.Add(lang, api.AddOptions{}); err != nil {
		return nil, err
	}

	topic := NewResource("topic", resKind, "topics", version, false, api.Options{Metadata: metadata.New()})
	if err := a.Add(topic, api.AddOptions{}); err != nil {
		return nil, err
	}

	return a, nil
}

func (g *client) mapObjects(top *Top, linkTo uuid.UID, res api.Resource, rel string, names []string) ([]*Object, error) {
	objects := make([]*Object, len(names))

	for i, name := range names {
		linkMeta := metadata.New()
		linkMeta.Set("relation", rel)

		// NOTE: we are setting uid to the name of the object
		// this is so we avoid duplicating topics with the same name
		uid := uuid.NewFromString(strings.ToLower(name + "-" + res.Name()))

		links := []link{
			{uid: linkTo, opts: api.LinkOptions{Metadata: linkMeta}},
		}

		objects[i] = NewObject(uid, strings.ToLower(name), api.NsGlobal, res, api.Options{Metadata: metadata.New()}, links)

		if err := top.Add(objects[i], api.AddOptions{MergeLinks: true}); err != nil {
			return nil, err
		}
	}

	return objects, nil
}

func (g *client) fetchRepos(reposChan chan<- []*github.StarredRepository, done <-chan struct{}) error {
	defer close(reposChan)

	opts := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: g.opts.Paging},
	}

	for {
		repos, resp, err := g.gh.Activity.ListStarred(g.ctx, g.opts.User, opts)
		if err != nil {
			return err
		}

		select {
		case reposChan <- repos:
		case <-done:
			return nil
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return nil
}

func (g *client) mapRepos(reposChan <-chan []*github.StarredRepository, top *Top, res resources) error {
	for repos := range reposChan {
		// NOTE: we are only iterating over the repos resources
		// since topics and langs are merely adjacent nodes of repos objects
		// and do not have any API endpoint for querying them further
		for _, repo := range repos {
			m := metadata.New()
			m.Set("starred_at", repo.StarredAt)
			m.Set("git_url", repo.Repository.GetURL)

			owner := *repo.Repository.Owner.Login
			if repo.Repository.Organization != nil {
				owner = *repo.Repository.Organization.Login
			}

			uid := uuid.NewFromString(*repo.Repository.NodeID)

			topicObjects, err := g.mapObjects(top, uid, res.topic, topicRel, repo.Repository.Topics)
			if err != nil {
				return err
			}

			var langObjects []*Object
			if repo.Repository.Language != nil {
				var err error
				langObjects, err = g.mapObjects(top, uid, res.lang, langRel, []string{*repo.Repository.Language})
				if err != nil {
					return err
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

			obj := NewObject(uid, *repo.Repository.Name, owner, res.repo, api.Options{Metadata: m}, repoLinks)

			if err := top.Add(obj, api.AddOptions{MergeLinks: true}); err != nil {
				return err
			}
		}
	}

	return nil
}

func getResource(a api.API, name, version string) (api.Resource, error) {
	q := query.Build().
		Name(name, query.StringEqFunc(name)).
		Version(version, query.StringEqFunc(version))

	rx, err := a.Get(q)
	if err != nil {
		return nil, err
	}

	return rx[0], nil
}

type resources struct {
	repo  api.Resource
	topic api.Resource
	lang  api.Resource
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

	res := resources{
		repo:  repoRes,
		topic: topicRes,
		lang:  langRes,
	}

	reposChan := make(chan []*github.StarredRepository, g.opts.Workers)
	errChan := make(chan error)
	done := make(chan struct{})

	var wg sync.WaitGroup

	// launch repo processing workers
	// these are building the graph
	for i := 0; i < g.opts.Workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			select {
			case errChan <- g.mapRepos(reposChan, top, res):
			case <-done:
			}
		}(i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- g.fetchRepos(reposChan, done)
	}()

	select {
	case err = <-errChan:
		close(done)
	case <-done:
	}

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return top, nil
}
