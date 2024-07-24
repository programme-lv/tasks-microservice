#! /bin/bash

set -ex # exit on error, print each command

SCRIPT_DIR=$(dirname $0)

pushd "$SCRIPT_DIR"/..

sam build
sam package

popd