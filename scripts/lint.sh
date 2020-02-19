#!/usr/bin/env bash
set -euo pipefail


export GOBIN=$(dirname `which go`)

export GO111MODULE=off

echo -n "Formatting golang code... "
gofmtout=$(go fmt ./...)
if [ "$gofmtout" ]; then
  printf "\n\n"
  echo "Files with formatting errors:"
  echo "${gofmtout}"
  exit 1
fi
echo "done."

echo -n "Linting golang code... "
# TODO: Make sure golangci-lint is installed and ready to be run
#GOLANG_LINT_CMD=golangci-lint
#
#if [ ! "$(command -v ${GOLANG_LINT_CMD})" ]; then
#  echo -n "installing ${GOLANG_LINT_CMD}... "
#  go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
#fi

go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
golangci-lint run
echo "done."


echo -n "Scanning for securirty issues... "
GOSEC_CMD=gosec

if [ ! "$(command -v ${GOSEC_CMD})" ]; then
  echo -n "installing ${GOSEC_CMD}... "
  go get -u github.com/securego/gosec/cmd/gosec
fi

gosec ./...
echo "done."

#GOSEC_VERSION=v2.2.0
#curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b . $GOSEC_VERSION
#./gosec ./...
#rm gosec