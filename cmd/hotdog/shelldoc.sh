#!/usr/bin/env bash
if [[ $# -lt 1 ]]; then
  command cat <<EOF 1>&2
$(basename "$0") [input files]...
EOF
  exit 1
fi

for arg in "$@"; do
  cat $arg | hotdog -first='^\s*#{30}' -middle='^\s*#' -last='^\s*#{30}' | while IFS= read -r line || [[ -n "$line" ]]; do
    printf '%s' "$line" | tr "\036" '\n' | sed -E 's/^[ \f\n\r\t]*#*//g' | yq .
  done
done | jq -s .
