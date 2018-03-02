#!/usr/bin/env bash

./build.sh

pushd examples
  export TF_LOG=DEBUG
  terraform init
  terraform apply
popd
