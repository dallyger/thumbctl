package path

import (
	"fmt"
	"path/filepath"

	"go.vnbr.de/thumbctl/internal/thumbnail"
)

type PathCmd struct {
	Size string `enum:"normal,large,x-large,xx-large" default:"large"`
	Path string `arg:""`
}

func (cmd *PathCmd) Run() error {
	thumb, err := thumbnail.FromPath(cmd.Path)
	if err != nil {
		return err
	}

	path := filepath.Join(
		thumbnail.MustGetRootDir(),
		cmd.Size,
		thumb.UriHash+".png",
	)

	fmt.Printf("%s\n", path)

	return nil
}
