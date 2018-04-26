#!/usr/bin/env bash
set -e

# Get the base repository path
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

# Move into our base repository path
cd "$DIR"

if [[ $(uname) == "Linux" ]]; then
  cp bin/* $GOPATH/bin/.
elif [[ $(uname) == "Darwin" ]]; then
  cp bin/* $GOPATH/bin/.
else
  echo "Unable to install on $(uname). Use Linux or Darwin."
  exit 1
fi
