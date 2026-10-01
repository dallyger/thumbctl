package main

import (
	"github.com/alecthomas/kong"
	"go.vnbr.de/thumbctl/internal/cmd/path"
)

type Context struct {
	Verbose bool
}

type RootCmd struct {
	Path path.PathCmd `cmd:"" help:"Show the path at which a thumbnail would be stored."`
}

func main() {
	ctx := kong.Parse(&RootCmd{})
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
