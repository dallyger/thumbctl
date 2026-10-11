package thumbnailer

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
)

// Get directories in which thumbnailers are installed to.
func GetLocations() []string {
	var dirs []string
	for _, dir := range getDataDirs() {
		dirs = append(dirs, filepath.Join(dir, "thumbnailers"))
	}

	return dirs
}

func GetEntryFiles() []string {
	var files []string
	for _, dir := range GetLocations() {
		matches, _ := filepath.Glob(path.Join(dir, "*.thumbnailer"))
		for _, match := range matches {
			files = append(files, match)
		}
	}
	return files
}

func getDataDirs() []string {
	dataHome := xdg.DataHome

	if dataHome == "" {
		home, err := os.UserHomeDir()

		if err != nil {
			slog.Warn("cannot determine home directory", "error", err)
			return xdg.DataDirs
		}

		dataHome = path.Join(home, ".local", "share")
	}

	return append([]string{dataHome}, xdg.DataDirs...)
}
