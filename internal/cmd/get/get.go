package get

import (
	"fmt"
	"log/slog"

	"go.vnbr.de/thumbctl/internal/metadata"
	"go.vnbr.de/thumbctl/internal/thumbnail"
	"go.vnbr.de/thumbctl/internal/thumbnailer"
)

type GetCmd struct {
	Path string `arg:"" help:"Path to the file to check."`
	Size string `enum:"normal,large,x-large,xx-large" default:"large" help:"Choose normal, large, x-large or xx-large."`
}

func (cmd *GetCmd) Run() error {
	meta, err := metadata.FromPath(cmd.Path)
	if err != nil {
		return err
	}

	cached, err := thumbnail.Find(meta, cmd.Size)
	if err != nil {
		return err
	}
	if cached != "" {
		slog.Debug("cache hit")
		fmt.Printf("%s\n", cached)
		return nil
	}
	slog.Debug("cache miss")

	tmp, err := thumbnailer.CreateFile(cmd.Path, cmd.Size)
	if err != nil {
		return err
	}

	dest, err := thumbnail.WriteFile(tmp, meta)

	if err != nil {
		return fmt.Errorf("failed writing thumbnail: %s", err)
	}

	fmt.Println(dest)
	return nil
}
