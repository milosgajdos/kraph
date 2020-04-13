package cmd

import (
	"github.com/milosgajdos/kraph/cmd/kctl/app/cmd/build"
	"github.com/urfave/cli/v2"
)

// Build builds app commands and reutrns them
func Build() []*cli.Command {
	cmds := make([]*cli.Command, 0)

	cmds = append(cmds, build.New())

	return cmds
}
