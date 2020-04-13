package app

import (
	"sort"

	"github.com/milosgajdos/kraph/cmd/kctl/app/cmd"
	"github.com/urfave/cli/v2"
)

const (
	name = "kctl"
)

// New creates kctl cli app and returns it
func New() *cli.App {
	app := &cli.App{
		Name:     name,
		Usage:    "build and query API object graphs",
		Commands: []*cli.Command{},
	}

	app.Commands = append(app.Commands, cmd.Build()...)

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
