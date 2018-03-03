#!/usr/bin/env bash

if [[ -z "${GOPATH}" ]]; then
  echo $GOPATH
  export GOPATH="${PWD}/../../../.."
fi
go build -o "${GOPATH}/bin/terraform-provider-opsman"
