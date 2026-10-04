package thumbnailer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPrepareExec(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  []string
		err  error
		opts RunExecOpts
	}{
		{
			name: "no-substitutions",
			in:   "/bin/true",
			out:  []string{"/bin/true"},
			err:  nil,
			opts: RunExecOpts{},
		},
		{
			name: "substitutions",
			in:   "/bin/true %% %u %i %s %o",
			out: []string{
				"/bin/true",
				"%",
				"file:///tmp/in%20(1).png",
				"/tmp/in (1).png",
				"256",
				"/tmp/out.png",
			},
			opts: RunExecOpts{
				Output: "/tmp/out.png",
				Path:   "/tmp/in (1).png",
				Uri:    "file:///tmp/in%20(1).png",
				Size:   256,
			},
		},
		{
			name: "invalid-substitution",
			in:   "/bin/true %x",
			out:  []string{},
			err: MalformedExecError{
				Err:  "invalid substitution: %x",
				Exec: "/bin/true %x",
				Pos:  11,
			},
		},
		{
			name: "substitution-cutoff",
			in:   "/bin/true %",
			out:  []string{},
			err: MalformedExecError{
				Err:  "invalid substitution: %",
				Exec: "/bin/true %",
				Pos:  11,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Entry{
				Exec: tt.in,
			}

			out, err := e.PrepareExec(tt.opts)

			if diff := cmp.Diff(tt.err, err); diff != "" {
				t.Errorf("error mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.out, out); diff != "" {
				t.Errorf("output mismatch (-want +got):\n%s", diff)
			}

		})
	}
}
