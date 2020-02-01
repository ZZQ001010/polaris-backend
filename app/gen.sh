#!/bin/sh

echo "Generating code from GraphQL schema..."

cd "$(dirname "$0")"

go run github.com/99designs/gqlgen
