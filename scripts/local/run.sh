#!/usr/bin/env bash

set -e;
# export GOPATH="$(cd $(dirname "$BASH_SOURCE") && pwd)"


export GOPATH="$PWD";
export huru_api_port="3000"

go clean
go install huru
"$GOPATH/bin/huru"