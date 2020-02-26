#!/usr/bin/env bash
set -euo pipefail

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
golangci-lint run
echo "done."


echo -n "Scanning for securirty issues... "
gosec ./...
echo "done."
