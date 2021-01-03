#!/usr/bin/env bash

set -euo pipefail

echo 'Dearmor GPG Signature'
file ./*

FILE_NAME=""
files=$(ls)

for filename in $files; do
  echo "$filename"
  if [ "${filename##*.}"x = "sig"x ]; then
    echo "Found signature file"
    FILE_NAME="$filename"
    break
  fi
done

if [ ! "${FILE_NAME}" ]; then
  echo "Signature file not found"
  exit 1
fi

#cat "${FILE_NAME}"
cp "${FILE_NAME}" "${FILE_NAME}.bak"
gpg --dearmor "${FILE_NAME}"
mv "${FILE_NAME}.gpg" "${FILE_NAME}"
#rm "${FILE_NAME}"

echo "Print file info"
ls -al
