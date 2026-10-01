package main

import (
	"github.com/alecthomas/kong"
)

type Context struct {
	Verbose bool
}

type RootCmd struct {
}

func main() {
	ctx := kong.Parse(&RootCmd{})
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
