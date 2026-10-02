package thumbnailers

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"go.vnbr.de/thumbctl/internal/thumbnailer"
)

var outputHumanMaxMimeTypes = 5

type ThumbnailersCmd struct {
	Format string `enum:"human,json" default:"human"`
}

type Entries struct {
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Path     string   `json:"path"`
	Exec     string   `json:"cmd"`
	TryExec  string   `json:"check"`
	MimeType []string `json:"mimetypes"`
}

func (cmd *ThumbnailersCmd) Run() error {

	thumbnailers, errs := thumbnailer.GetEntries()
	for _, err := range errs {
		fmt.Fprintf(os.Stderr, "thumbnailer: error parsing file: %s", err)
	}

	var entries Entries
	for path, entry := range thumbnailers {
		entries.Entries = append(entries.Entries, Entry{
			Path:     path,
			Exec:     entry.Exec,
			TryExec:  entry.TryExec,
			MimeType: entry.MimeType,
		})
	}

	switch cmd.Format {
	case "json":
		return outputJson(entries)
	default:
		return outputHuman(entries)
	}
}

func outputHuman(e Entries) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, entry := range e.Entries {
		fmt.Fprintf(w, "Path:\t%s\n", entry.Path)
		fmt.Fprintf(w, "Check:\t%s\n", entry.TryExec)
		fmt.Fprintf(w, "Command:\t%s\n", entry.Exec)
		fmt.Fprintf(w, "Mime Types:\n")
		for i, mime := range entry.MimeType[:min(outputHumanMaxMimeTypes+1, len(entry.MimeType))] {
			if i == outputHumanMaxMimeTypes {
				fmt.Fprintf(w, "\t- (%d more)\n", len(entry.MimeType)-outputHumanMaxMimeTypes)
			} else {
				fmt.Fprintf(w, "\t- %s\n", mime)
			}
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	return nil
}

func outputJson(e Entries) error {
	if out, err := json.Marshal(e); err != nil {
		return err
	} else {
		fmt.Println(string(out))
		return nil
	}
}
