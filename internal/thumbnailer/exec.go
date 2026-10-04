package thumbnailer

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"go.vnbr.de/thumbctl/internal/metadata"
	"go.vnbr.de/thumbctl/internal/mimetype"
	"go.vnbr.de/thumbctl/internal/thumbnail"
)

type RunExecOpts struct {
	Output string
	Path   string
	Size   int
	Uri    string
}

type MalformedExecError struct {
	Err  string
	Exec string
	Pos  int
}

func (e MalformedExecError) Error() string {
	return fmt.Sprintf("malformed exec command at position %d: %s", e.Pos, e.Err)
}

func CreateFile(filepath, sizePreset string) (img []byte, err error) {
	meta, err := metadata.FromPath(filepath)
	if err != nil {
		return
	}

	mime, err := mimetype.FromPath(filepath)
	slog.Debug("thumbnailer: detect mimetype", "mime", mime, "path", meta.AbsolutePath)
	if err != nil {
		return
	}

	dim, err := thumbnail.SizeToInt(sizePreset)
	if err != nil {
		return
	}

	dir, err := os.MkdirTemp("", "thumbctl-*")
	defer os.RemoveAll(dir)
	if err != nil {
		return img, fmt.Errorf("failed creating temporary directory: %s", err)
	}
	slog.Debug("thumbnailer: created temp directory", "path", dir)

	o := RunExecOpts{
		Output: path.Join(dir, fmt.Sprintf("%s.png", meta.UriHash)),
		Path:   meta.AbsolutePath,
		Size:   dim,
		Uri:    meta.CanonicalUri,
	}

	if err := RunThumbnailer(mime, o); err != nil {
		return img, err
	}

	slog.Debug("thumbnailer: created temp thumbnail", "path", o.Output)

	handle, err := os.Open(o.Output)
	if err != nil {
		return img, fmt.Errorf("failed opening thumbnail: %s", err)
	}

	content, err := io.ReadAll(handle)
	if err != nil {
		return img, fmt.Errorf("failed reading thumbnail: %s", err)
	}

	return content, nil
}

func RunThumbnailer(mime string, o RunExecOpts) error {
	avail, _ := GetEntries()
	var t Entry
	for t = range avail.MatchMimeType(mime) {
		if err := t.RunTryExec(); err != nil {
			slog.Warn("thumbnailer skipped", "error", err, "cmd", t.TryExec)
			continue
		}

		if err := t.RunExec(o); err != nil {
			slog.Warn("thumbnailer call failed", "error", err, "cmd", t.Exec)
			continue
		}

		return nil
	}

	return fmt.Errorf("no thumbnailer available for type %s\n", mime)
}

func (e *Entry) RunExec(o RunExecOpts) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	args, err := e.PrepareExec(o)
	if err != nil {
		return err
	}

	bin := args[0]
	args = args[1:]

	cb := exec.CommandContext(ctx, bin, args...)
	cb.WaitDelay = 2 * time.Second
	cb.Cancel = func() error { return cb.Process.Signal(syscall.SIGTERM) }

	return cb.Run()
}

func (e *Entry) RunTryExec() error {
	if e.TryExec != "" {
		if path, _ := exec.LookPath(e.TryExec); path == "" {
			return fmt.Errorf("program not found")
		}
	}
	return nil
}

func (e *Entry) PrepareExec(o RunExecOpts) ([]string, error) {
	pos := 0
	var args []string
	for arg := range strings.FieldsSeq(e.Exec) {
		if s, err := substitute(arg, o); err != nil {
			return []string{}, MalformedExecError{
				Err:  err.(MalformedExecError).Err,
				Exec: e.Exec,
				Pos:  pos + err.(MalformedExecError).Pos,
			}
		} else {
			args = append(args, s)
			// assumption: args are split by a single whitespace char.
			pos = pos + 1 + len(s)
		}
	}
	return args, nil
}

func substitute(arg string, o RunExecOpts) (res string, err error) {
	var pos int
	for pos < len(arg) {
		curr := arg[pos]

		if curr == '%' {
			isLast := pos == len(arg)-1
			if isLast {
				return "", MalformedExecError{
					Err: "invalid substitution: %",
					Pos: pos + 1,
				}
			}

			next := arg[pos+1]

			switch next {
			case '%':
				res = res + "%"
				pos = pos + 2
			case 'u':
				res = res + o.Uri
				pos = pos + 2
			case 'i':
				res = res + o.Path
				pos = pos + 2
			case 'o':
				res = res + o.Output
				pos = pos + 2
			case 's':
				res = res + strconv.Itoa(o.Size)
				pos = pos + 2
			default:
				return "", MalformedExecError{
					Err: "invalid substitution: %" + string(next),
					Pos: pos + 1,
				}
			}

		} else {
			res = res + string(curr)
			pos++

		}
	}

	return
}
