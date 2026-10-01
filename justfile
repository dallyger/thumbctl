help:
    @just --list

run +args="":
    go run cmd/thumbctl/main.go {{args}}

build bin="target/thumbctl":
    go build -ldflags="\
        -X 'go.vnbr.de/thumbctl/internal/cmd/version.version=$(git describe --tags --always --dirty)' \
        -X 'go.vnbr.de/thumbctl/internal/cmd/version.buildDate=$(TZ=UTC date +"%Y-%m-%dT%H:%M:%SZ")' \
        " -o "{{justfile_dir()}}/{{bin}}" cmd/thumbctl/main.go
