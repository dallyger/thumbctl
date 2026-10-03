package mimetype

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"go.vnbr.de/thumbctl/internal/mimetype"
)

type MimeTypeCmd struct {
	Path        string `arg:"" help:"Path to the file to check."`
	Format      string `aliases:"fmt" default:"human" enum:"human,json,short" help:"Choose human, short or json."`
	ShowAliases bool   `name:"aliases" aliases:"a" negatable:"" default:"0" help:"Show mime type aliases."`
}

func (cmd *MimeTypeCmd) Run() error {
	m, err := mimetype.FromPath(cmd.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	switch cmd.Format {
	case "json":
		outputJson(cmd, m)
	case "short":
		outputShort(cmd, m)
	default:
		outputHuman(cmd, m)
	}

	return nil
}

func outputHuman(cmd *MimeTypeCmd, m string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, byte(' '), 0)
	fmt.Fprintf(w, "Mime:\t%s\n", m)
	if cmd.ShowAliases {
		fmt.Fprintf(w, "Aliases:\n")
		for _, a := range mimetype.Aliases(m) {
			fmt.Fprintf(w, "- %s\n", a)
		}
	}
	w.Flush()
}

func outputJson(cmd *MimeTypeCmd, m string) {
	type t struct {
		Mime    string   `json:"mime"`
		Aliases []string `json:"aliases,omitempty"`
	}
	o := t{Mime: m}

	if cmd.ShowAliases {
		o.Aliases = mimetype.Aliases(m)
	}

	j, err := json.Marshal(o)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: encoding JSON: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(j))
}

func outputShort(cmd *MimeTypeCmd, m string) {
	fmt.Println(m)
	if cmd.ShowAliases {
		for _, a := range mimetype.Aliases(m) {
			fmt.Println(a)
		}
	}
}
