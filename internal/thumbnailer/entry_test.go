package thumbnailer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFromFile(t *testing.T) {
	type test struct {
		name   string
		config string
		expect Entry
	}

	tests := []test{
		{
			name: "generic",
			config: `
[Thumbnailer Entry]
Exec=evince-thumbnailer -s %s %u %o
MimeType=application/pdf;application/x-bzpdf;application/x-gzpdf;
`,
			expect: Entry{
				Exec:    "evince-thumbnailer -s %s %u %o",
				TryExec: "",
				MimeType: []string{
					"application/pdf",
					"application/x-bzpdf",
					"application/x-gzpdf",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			entry, err := FromFile([]byte(tt.config))
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.expect, entry); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}

		})
	}
}
