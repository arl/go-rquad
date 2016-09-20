#!/bin/bash

GOLINT="${HOME}/gopath/bin/golint"
EXCLUDES='vendor|_string.go'

find . -name '*.go' -print | egrep -v "${EXCLUDES}" | xargs $GOLINT -set_exit_status
