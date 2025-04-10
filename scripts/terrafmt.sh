#!/usr/bin/env bash

echo "==> Checking terraform blocks are formatted..."

files=$(find ./azuredevops -type f -name "*_test.go")
error=false

for f in $files; do
  if command -v terrafmt; then
    terrafmt diff -c -q -f "$f" || error=true
  else command -v $(GOPATH)/bin/terrafmt;
   $(GOPATH)/bin/terrafmt diff -c -q -f "$f" || error=true
  fi
done

if ${error}; then
  echo "------------------------------------------------"
  echo ""
  echo "The preceding files contain terraform blocks that are not correctly formatted or contain errors."
  echo "You can fix this by running 'make tools' and then terrafmt on them."
  echo ""
  echo "to easily fix all terraform blocks:"
  echo "$ make terrafmt"
  echo ""
  echo "format a single test file:"
  echo "$ terrafmt fmt -f ./azuredevops/internal/acceptancetests/data_area_test.go"
fi

exit 0