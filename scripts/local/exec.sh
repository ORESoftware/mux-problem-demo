#!/usr/bin/env bash
# export GOPATH="$(cd $(dirname "$BASH_SOURCE") && pwd)"

export GOPATH="$PWD";
export huru_api_port="3000"
"$GOPATH/bin/huru"