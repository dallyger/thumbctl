package thumbnail

import (
	"fmt"
	"os"
	"path/filepath"

	"go4.org/xdgdir"
)

// Get the thumbnails directory.
//
// If the environment variable $XDG_CACHE_HOME is set and not blank then the
// directory $XDG_CACHE_HOME/thumbnails will be used, otherwise
// $HOME/.cache/thumbnails will be used.
func MustGetRootDir() string {
	if cache := xdgdir.Cache.Path(); cache != "" {
		return filepath.Join(cache, "thumbnails")
	}

	home, err := os.UserHomeDir()

	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: cannot determine home directory: %s\n", err)
		os.Exit(1)
	}

	return filepath.Join(home, ".cache", "thumbnails")
}
