#!/usr/bin/env bash

export PATH="$HOME/.local/bin:${PATH}"  # eb is installed here on linux

export access_id="$huru_access_id"
export access_key="$huru_access_key" #  0fQhN7QjaeD8DPc0AW/owiqvdscvTc/gYk1vVbsB

eb codesource local
eb deploy --staged huru-api2-dev #HuruApi-env-1 #huru-api-dev2 # e-c72mgj2tmr # e-zp2ptw28nv

