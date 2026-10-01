package version

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"text/tabwriter"

	"github.com/alecthomas/kong"
)

type VersionCmd bool

// may be set at build time via -ldflags
var (
	version    = ""
	buildDate  = ""
	commit     = ""
	commitDate = ""
)

func (cmd *VersionCmd) BeforeReset(app *kong.Kong) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	if info, ok := debug.ReadBuildInfo(); ok {

		if version == "" {
			version = info.Main.Version
		}

		for _, kv := range info.Settings {
			switch kv.Key {
			case "vcs.revision":
				if commit == "" {
					commit = kv.Value
				}
			case "vcs.time":
				if commitDate == "" {
					commitDate = kv.Value
				}
			}
		}
	}

	if version == "" {
		version = "(devel)"
	}

	fmt.Fprintf(w, "Version:\t%s\n", version)

	if commit != "" {
		fmt.Fprintf(w, "Commit:\t%.10s\n", commit)
	}

	if commitDate != "" {
		fmt.Fprintf(w, "Commit Date:\t%s\n", commitDate)
	}

	if buildDate != "" {
		fmt.Fprintf(w, "Build Date:\t%s\n", buildDate)
	}

	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Compiler:\t%s\n", runtime.Compiler)
	fmt.Fprintf(w, "Go Version:\t%s\n", runtime.Version())
	fmt.Fprintf(w, "Platform:\t%s/%s\n", runtime.GOOS, runtime.GOARCH)

	err := w.Flush()

	if err != nil {
		return nil
	}

	os.Exit(0)
	return nil
}
