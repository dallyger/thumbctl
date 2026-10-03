package path

import (
	"fmt"
	"path/filepath"

	"go.vnbr.de/thumbctl/internal/metadata"
	"go.vnbr.de/thumbctl/internal/thumbnail"
)

type PathCmd struct {
	Path string `arg:"" help:"Path to the file to check."`
	Size string `enum:"normal,large,x-large,xx-large" default:"large" help:"Choose normal, large, x-large or xx-large."`
}

func (cmd *PathCmd) Run() error {
	meta, err := metadata.FromPath(cmd.Path)
	if err != nil {
		return err
	}

	path := filepath.Join(
		thumbnail.MustGetRootDir(),
		cmd.Size,
		meta.UriHash+".png",
	)

	fmt.Printf("%s\n", path)

	return nil
}
