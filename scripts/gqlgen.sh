#!/bin/bash
printf "\nRegenerating gqlgen files\n"
go get -u github.com/99designs/gqlgen/internal/imports@v0.13.0
go get -u github.com/99designs/gqlgen/internal/code@v0.13.0
go get -u github.com/99designs/gqlgen/cmd@v0.13.0
go get -u github.com/vektah/gqlparser/v2@v2.1.0

rm -f internal/gql/generated/exec.go \
    internal/gql/models/generated.go \
    internal/gql/resolvers/generated/generated.go
time go run -v github.com/99designs/gqlgen $1
printf "\nDone.\n\n"