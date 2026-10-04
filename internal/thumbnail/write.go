package thumbnail

import (
	"fmt"
	"os"

	"go.vnbr.de/thumbctl/internal/metadata"
	"go.vnbr.de/thumbctl/internal/thumbnail/data"
)

// Embed metadata into PNG and store file in cache. Ensures atomic writes.
func WriteFile(b []byte, meta metadata.Meta) (string, error) {

	if meta.MTime.IsZero() {
		return "", fmt.Errorf("MTime not set.")
	}

	if meta.Software == "" {
		meta.Software = "thumbctl"
	}

	b, err := data.EmbedTextualData(b, meta)
	if err != nil {
		return "", fmt.Errorf("failed embedding metadata: %s", err)
	}

	size, err := PngToSize(b)
	if err != nil {
		return "", err
	}

	dest := getPath(meta, size)

	file, err := os.OpenFile(dest+".part", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return "", err
	}
	defer file.Close()
	defer os.Remove(dest + ".part")

	if _, err := file.Write(b); err != nil {
		return "", err
	}

	if err := os.Rename(dest+".part", dest); err != nil {
		return "", err
	} else {
		return dest, nil
	}
}
