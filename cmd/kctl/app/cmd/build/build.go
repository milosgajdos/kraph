package build

import (
	"github.com/urfave/cli/v2"
)

// New creates new build command and returns it
func New() *cli.Command {
	build := &cli.Command{
		Name:        "build",
		Usage:       "build a graph",
		Subcommands: []*cli.Command{},
	}

	build.Subcommands = append(build.Subcommands, K8s())

	return build
}
