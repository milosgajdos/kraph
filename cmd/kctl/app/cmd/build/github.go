package build

import (
	"fmt"
	"sync"

	"github.com/google/go-github/v32/github"
	"github.com/milosgajdos/kraph"
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gh/star"
	"github.com/milosgajdos/kraph/pkg/store"
	"github.com/milosgajdos/kraph/pkg/store/memory"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

var (
	ghToken  string
	ghUser   string
	ghPaging int
)

// K8s returns K8s subcommand for build command
func GH() *cli.Command {
	return &cli.Command{
		Name:     "github",
		Aliases:  []string{"gh"},
		Category: "build",
		Usage:    "GitHub graph",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "store",
				Aliases:     []string{"s"},
				Value:       "memory",
				Usage:       "graph store",
				Destination: &graphStore,
			},
			&cli.StringFlag{
				Name:        "store-id",
				Aliases:     []string{"id"},
				Value:       "kctl",
				Usage:       "store ID",
				Destination: &graphStoreID,
			},
			&cli.StringFlag{
				Name:        "store-url",
				Aliases:     []string{"surl"},
				Value:       "",
				Usage:       "URL of the graph store",
				EnvVars:     []string{"STORE_URL"},
				Destination: &graphStoreURL,
			},
			&cli.StringFlag{
				Name:        "graph",
				Aliases:     []string{"g"},
				Value:       "star",
				Usage:       "graph type",
				Destination: &graphType,
			},
			&cli.StringFlag{
				Name:        "format",
				Aliases:     []string{"f"},
				Value:       "dot",
				Usage:       "graph output format",
				Destination: &graphFormat,
			},
			&cli.StringFlag{
				Name:        "token",
				Aliases:     []string{"t"},
				Value:       "",
				Usage:       "GitHub API token",
				EnvVars:     []string{"GITHUB_TOKEN"},
				Destination: &ghToken,
			},
			&cli.StringFlag{
				Name:        "user",
				Aliases:     []string{"u"},
				Value:       "",
				Usage:       "GitHub username",
				Destination: &ghUser,
			},
			&cli.IntFlag{
				Name:        "paging",
				Aliases:     []string{"p"},
				Value:       100,
				Usage:       "GitHub API response paging",
				Destination: &ghPaging,
			},
		},
		Action: func(c *cli.Context) error {
			return runGH(c)
		},
	}
}

func runH(ctx *cli.Context) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx.Context, ts)

	ghClient := github.NewClient(tc)

	var gstore store.Store
	var err error

	switch graphStore {
	case "memory":
		gstore, err = memory.NewStore(graphStoreID, store.Options{})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported store: %s", graphStore)
	}

	k, err := kraph.New(gstore)
	if err != nil {
		return fmt.Errorf("failed to create kraph: %w", err)
	}

	// TODO: figure this out
	var filters []kraph.Filter

	var client api.Scraper

	if ghUser == "" {
		return fmt.Errorf("you must specify GitHub username")
	}

	switch graphType {
	case "star":
		client = star.NewScraper(ctx.Context, ghClient, star.Paging(ghPaging), star.User(ghUser))
	default:
		return fmt.Errorf("unsupported graph type: %s", graphType)
	}

	s := newSpinner()

	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Run(ctx.Context, done)
	}()

	cleanup := func() {
		close(done)
		wg.Wait()
		s.Finish()
	}

	if err = k.Build(client, filters...); err != nil {
		cleanup()
		return fmt.Errorf("failed to build kraph: %w", err)
	}

	cleanup()

	// only print the graph if it's an in-memory graph
	if graphStore == "memory" {
		//nodes, err := k.Store().Graph().Nodes()
		//if err != nil {
		//	return err
		//}

		//edges, err := k.Store().Graph().Edges()
		//if err != nil {
		//	return err
		//}

		//fmt.Printf("Nodes: %d, Edges: %d\n", len(nodes), len(edges))

		graphOut, err := graphToOut(k.Store().Graph(), graphFormat)
		if err != nil {
			return err
		}
		fmt.Println(graphOut)
	}

	return nil
}
