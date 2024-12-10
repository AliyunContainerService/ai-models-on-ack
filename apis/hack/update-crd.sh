#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
"${SCRIPT_ROOT}"/hack/controller-gen-darwin-arm64 crd paths=./... output:crd:dir=./crd