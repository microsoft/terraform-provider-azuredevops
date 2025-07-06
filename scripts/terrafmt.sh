#!/usr/bin/env bash

echo "==> Checking terraform blocks are formatted..."

files=$(find ./azuredevops -type f -name "*_test.go")
error=false

for f in $files; do
  terrafmt diff -c -q -f "$f" || error=true
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
  exit 1
fi

exit 0
