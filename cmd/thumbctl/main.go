package main

import (
	"github.com/alecthomas/kong"
	"go.vnbr.de/thumbctl/internal/cmd/info"
	"go.vnbr.de/thumbctl/internal/cmd/license"
	"go.vnbr.de/thumbctl/internal/cmd/mimetype"
	"go.vnbr.de/thumbctl/internal/cmd/path"
	"go.vnbr.de/thumbctl/internal/cmd/thumbnailers"
	"go.vnbr.de/thumbctl/internal/cmd/version"
)

type Context struct {
	Verbose bool
}

type RootCmd struct {
	Info         info.InfoCmd                 `cmd:"" help:"Show information about a thumbnail."`
	Path         path.PathCmd                 `cmd:"" help:"Show the path at which a thumbnail would be stored."`
	Mime         mimetype.MimeTypeCmd         `cmd:"" help:"Show MIME type of a file."`
	Thumbnailers thumbnailers.ThumbnailersCmd `cmd:"" help:"Show installed and available thumbnailer tools."`

	License license.LicenseCmd `cmd:"" help:"Show license information and exit."`
	Version version.VersionCmd `help:"Show build information and exit."`
}

func main() {
	ctx := kong.Parse(&RootCmd{})
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
