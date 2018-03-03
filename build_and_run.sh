#!/usr/bin/env bash

./build.sh

export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform.log

TF="${PWD}/../../../../bin/terraform"
if [[ ! -f "${TF}" ]]; then
  TF="terraform"
fi

pushd examples/opsman  
  $TF init
  $TF apply
popd
