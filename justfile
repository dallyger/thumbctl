help:
    @just --list

run +args="":
    go run cmd/thumbctl/main.go {{args}}

build:
    go build -o thumbctl cmd/thumbctl/main.go
