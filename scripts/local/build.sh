#!/usr/bin/env bash

# export GOPATH="$(cd $(dirname "$BASH_SOURCE") && pwd)"

#  go get github.com/gorilla/mux

set -e;

export GOPATH="$PWD";
go clean
go install huru