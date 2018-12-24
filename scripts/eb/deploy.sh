#!/usr/bin/env bash

export PATH="$HOME/.local/bin:${PATH}"  # eb is installed here on linux

export access_id="$huru_access_id"
export access_key="$huru_access_key" #

eb codesource local
eb deploy --staged huru-api2-dev

