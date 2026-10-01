package info

import (
	"fmt"
	"os"
	"text/tabwriter"

	"go.vnbr.de/thumbctl/internal/thumbnail"

	pngembed "github.com/sabhiram/png-embed"
)

type InfoCmd struct {
	Format string `enum:"human,json" default:"human"`
	Path   string `arg:""`
}

func (cmd *InfoCmd) Run() error {
	thumb, err := thumbnail.FromPath(cmd.Path)
	if err != nil {
		return err
	}

	if !thumb.IsThumbnail {
		return fmt.Errorf("given path is not a thumbnail file")
	}

	bs, err := os.ReadFile(thumb.AbsolutePath)
	if err != nil {
		return fmt.Errorf("failed reading thumbnail file at %s: %s", thumb.AbsolutePath, err)
	}

	tEXt, err := pngembed.Extract(bs)
	if err != nil {
		return fmt.Errorf("failed reading metadata of thumbnail at %s: %s", thumb.AbsolutePath, err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Thumbnail:\t%s\n", thumb.AbsolutePath)
	for k, v := range tEXt {
		fmt.Fprintf(w, "%s:\t%s\n", k, string(v))
	}
	w.Flush()

	return nil
}
