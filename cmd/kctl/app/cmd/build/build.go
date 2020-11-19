package build

import (
	"fmt"

	"github.com/milosgajdos/kraph/pkg/store"
	"github.com/urfave/cli/v2"
)

var (
	graphStore    string
	graphStoreURL string
	graphStoreID  string
	graphType     string
	graphFormat   string
)

func graphToOut(g store.Graph, format string) (string, error) {
	switch format {
	case "dot":
		dotGraph, ok := g.(store.DOTGraph)
		if !ok {
			return "", fmt.Errorf("unable to convert graph to %s format", format)
		}
		return dotGraph.DOT()
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
}

// New creates new build command and returns it
func New() *cli.Command {
	build := &cli.Command{
		Name:        "build",
		Usage:       "build a graph",
		Subcommands: []*cli.Command{},
	}

	build.Subcommands = append(build.Subcommands, K8s(), GH())

	return build
}
