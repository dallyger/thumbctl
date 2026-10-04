package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"go.vnbr.de/thumbctl/internal/cmd/get"
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
	Get          get.GetCmd                   `cmd:"" help:"Show thumbnail path and generate if it is stale or missing."`
	Info         info.InfoCmd                 `cmd:"" help:"Show information about a thumbnail."`
	Path         path.PathCmd                 `cmd:"" help:"Show the path at which a thumbnail would be stored."`
	Mime         mimetype.MimeTypeCmd         `cmd:"" help:"Show MIME type of a file."`
	Thumbnailers thumbnailers.ThumbnailersCmd `cmd:"" help:"Show installed and available thumbnailer tools."`

	License license.LicenseCmd `cmd:"" help:"Show license information and exit."`
	Verbose int                `short:"v" type:"counter" help:"Increase logging verbosity."`
	Version version.VersionCmd `help:"Show build information and exit."`
}

func (cmd *RootCmd) logLevel() slog.Level {
	switch cmd.Verbose {
	case 0:
		return slog.LevelWarn
	case 1:
		return slog.LevelInfo
	default:
		return slog.LevelDebug
	}
}

func main() {
	cmd := RootCmd{}
	ctx := kong.Parse(&cmd)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: cmd.logLevel(),
	})))

	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
