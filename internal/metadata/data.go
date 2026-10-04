package metadata

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"path/filepath"
)

type Meta struct {
	AbsolutePath string
	CanonicalUri string
	OriginalPath string
	UriHash      string
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

	u := url.URL{Path: meta.AbsolutePath}
	meta.CanonicalUri = "file://" + u.EscapedPath()
	meta.UriHash = fmt.Sprintf("%x", md5.Sum([]byte(meta.CanonicalUri)))

	return meta, nil
}
