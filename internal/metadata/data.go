package metadata

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"go.vnbr.de/thumbctl/internal/thumbnail"
)

type Meta struct {
	AbsolutePath string
	CanonicalUri string
	OriginalPath string
	UriHash      string
	IsThumbnail  bool
}

func FromPath(path string) (Meta, error) {
	meta := Meta{
		OriginalPath: path,
	}

	if path, err := filepath.Abs(path); err == nil {
		meta.AbsolutePath = path
	} else {
		return meta, err
	}

	// Check if file is already a thumbnail and return itself.
	// As per the freedesktop thumbnail spec:, we "must load and use these files directly."
	if thumbnail.IsThumbnail(meta.AbsolutePath) {
		meta.CanonicalUri = "file://" + meta.AbsolutePath
		meta.UriHash = strings.Split(filepath.Base(meta.AbsolutePath), ".")[0]
		meta.IsThumbnail = true
		return meta, nil
	}

	u := url.URL{Path: meta.AbsolutePath}
	meta.CanonicalUri = "file://" + u.EscapedPath()
	meta.UriHash = fmt.Sprintf("%x", md5.Sum([]byte(meta.CanonicalUri)))

	return meta, nil
}
