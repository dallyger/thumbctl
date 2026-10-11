# thumbctl

View, generate and manage thumbnails following the freedesktop.org specs.

## Usage

The main usage is to generate a preview for any file and return its location:

```sh
thumbctl get ./path/to/file/to/preview
```

You can use it to generate thumbnails in bulk:

```sh
find . -maxdepth 3 -type f | parallel thumbctl get "{}"
```

For additional functionality, check the help messages:

```sh
thumbctl --help
```

## Installation

```sh
go install golang.vnbr.de/thumbctl/cmd/thumbctl@latest
```

## Build from source

Requirements: `just`, `go`

Clone the project from one of the following locations:

```sh
git clone https://codeberg.org/dallyger/thumbctl.git
git clone https://github.com/dallyger/thumbctl.git
```

Compile binary and run it:

```
just build && target/thumbctl --help
```

## Thanks To

- https://jimmyg.org/blog/2023/thumbnailers/index.html
- https://tecnocode.co.uk/2013/10/21/writing-a-gnome-thumbnailer/

## Related Links

- Thumbnail standard: <https://specifications.freedesktop.org/thumbnail/latest/>
- GnomeDesktopThumbnailFactory: <https://gnome.pages.gitlab.gnome.org/gnome-desktop/html/gnome-desktop3/gnome-desktop3-GnomeDesktopThumbnailFactory.html>
- Shared MIME-info standard: <https://specifications.freedesktop.org/shared-mime-info/latest/>
- MD5 hash algorithm: <https://www.ietf.org/rfc/rfc1321.txt>
- PNG standard: <http://www.w3.org/TR/REC-png>
- URI standard: <https://www.ietf.org/rfc/rfc2396.txt>
