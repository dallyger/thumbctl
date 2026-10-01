package thumbnail

import (
	"path/filepath"
	"testing"
)

func TestAbsolutePath(t *testing.T) {
	t.Chdir("/tmp")
	cache := MustGetRootDir()

	tests := []struct {
		name string
		path string
		abs  string
		uri  string
		hash string
	}{
		{
			"absolute-path",
			"/home/user/Documents/foo.txt",
			"/home/user/Documents/foo.txt",
			"file:///home/user/Documents/foo.txt",
			"de3fbcbf97bda23a345e220096c185cc",
		},
		{
			"relative-path",
			"foo.txt",
			"/tmp/foo.txt",
			"file:///tmp/foo.txt",
			"de3fa011b3c05d1f7a89e0e66c7577ff",
		},
		{
			"escaped-path",
			"foo bat.txt",
			"/tmp/foo bar.txt",
			"file:///tmp/foo%20bar.txt",
			"c329efdb2167dd14ca2d591c80ed35ea",
		},
		// As specified in the spec, any thumbnail file must be used directly.
		{
			"thumbnail-self-reference",
			filepath.Join(cache, "dfc930eb4552e2baf77ca4d89c917a3b.png"),
			filepath.Join(cache, "dfc930eb4552e2baf77ca4d89c917a3b.png"),
			"file://" + filepath.Join(cache, "dfc930eb4552e2baf77ca4d89c917a3b.png"),
			"dfc930eb4552e2baf77ca4d89c917a3b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			thumb, err := FromPath(tt.abs)

			if err != nil {
				t.Fatalf("failure: %s\n", err)
			}

			if thumb.OriginalPath != tt.abs {
				t.Logf("mismatching absolute path.\nexpected: %s\nactual:   %s\n", tt.abs, thumb.AbsolutePath)
				t.Fail()
			}

			if thumb.CanonicalUri != tt.uri {
				t.Logf("mismatching canonical URI.\nexpected: %s\nactual:   %s\n", tt.uri, thumb.CanonicalUri)
				t.Fail()
			}

			if thumb.UriHash != tt.hash {
				t.Logf("mismatching URI hash.\nexpected: %s\nactual:   %s\n", tt.hash, thumb.UriHash)
				t.Fail()
			}
		})
	}
}
