help:
    @just --list

run +args="":
    go run cmd/thumbctl/main.go {{args}}

build bin="target/thumbctl":
    go build -ldflags="\
        -X 'go.vnbr.de/thumbctl/internal/cmd/version.version=$(git describe --tags --always --dirty)' \
        -X 'go.vnbr.de/thumbctl/internal/cmd/version.buildDate=$(TZ=UTC date +"%Y-%m-%dT%H:%M:%SZ")' \
        " -o "{{justfile_dir()}}/{{bin}}" cmd/thumbctl/main.go

# Check for incompatible licenses and save them to a file.
licenses out="THIRD_PARTY_LICENSES":
    #!/usr/bin/sh
    set -eu
    go-licenses check ./... --allowed_licenses=MIT,BSD-3-Clause,BSD-2-Clause,Apache-2.0,ISC
    go-licenses save --force \
        --add_dir_header \
        --save_path target/third_party/ \
        --ignore go.vnbr.de/thumbctl \
        ./...
    {
      echo "This project includes third-party software. Their licenses follow."
      find target/third_party/ -type f \( -iname 'LICENSE*' -o -iname 'NOTICE*' -o -iname 'COPYING*' \) | sort |
      while read -r f; do
        printf '\n================================================================================\n'
        printf '%s\n' "$(dirname "${f#target/third_party/}")"
        printf '================================================================================\n\n'
        cat "$f"
      done
    } > "{{out}}"
