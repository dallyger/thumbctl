package main

import (
	"github.com/alecthomas/kong"
	"go.vnbr.de/thumbctl/internal/cmd/info"
	"go.vnbr.de/thumbctl/internal/cmd/path"
	"go.vnbr.de/thumbctl/internal/cmd/version"
)

type Context struct {
	Verbose bool
}

type RootCmd struct {
	Info info.InfoCmd `cmd:"" help:"Show information about a thumbnail."`
	Path path.PathCmd `cmd:"" help:"Show the path at which a thumbnail would be stored."`

	Version version.VersionCmd `help:"Show build information and exit."`
}

func main() {
	ctx := kong.Parse(&RootCmd{})
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
