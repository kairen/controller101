#!/bin/bash

# Copyright © 2019 The controller101 Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

REPO_PATH="github.com/cloud-native-taiwan/controller101"

COV_FILE=coverage.txt
COV_TMP_FILE=coverage_tmp.txt

echo "Running go tests..."
cd ${GOPATH}/src/${REPO_PATH}
rm -f out/$COV_FILE || true
echo "mode: count" > out/$COV_FILE
for pkg in $(go list -f '{{ if .TestGoFiles }} {{.ImportPath}} {{end}}' ./...); do
    go test -tags "container_image_ostree_stub containers_image_openpgp" -v $pkg -covermode=count -coverprofile=out/$COV_TMP_FILE
    tail -n +2 out/$COV_TMP_FILE >> out/$COV_FILE || (echo "Unable to append coverage for $pkg" && exit 1)
done
rm -rf out/$COV_TMP_FILE

ignore="vendor\|\_gopath\|assets.go\|\_test.go\|out\|api"

echo "Checking gofmt..."
set +e
files=$(gofmt -l -s . | grep -v ${ignore})
set -e
if [[ $files ]]; then
  echo "Gofmt errors in files: $files"
  exit 1
fi
