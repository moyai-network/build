#!/usr/bin/env bash
set -euo pipefail

mapfile -t submods < <(git config --file .gitmodules --get-regexp '^submodule\..*\.path$' 2>/dev/null | awk '{print $2}')

if [ "${#submods[@]}" -eq 0 ]; then
  cloc_args=(.)
else
  # Escape submodule paths for regex and match both "pkg/x" and "./pkg/x".
  pattern="$(
    printf '%s\n' "${submods[@]}" \
      | sed 's/[.[\*^$()+?{|]/\\&/g' \
      | paste -sd'|' -
  )"
  cloc_args=(. --fullpath --not-match-d "^(\\./)?($pattern)(/|$)")
fi

# Per-language table.
cloc "${cloc_args[@]}" "$@"

# Explicit totals line (blank + comment + code).
csv="$(cloc --quiet --csv "${cloc_args[@]}" "$@")"

sum_row="$(printf '%s\n' "$csv" | awk -F, '$1=="SUM" || $2=="SUM"{print $3 "," $4 "," $5; exit}')"
if [ -n "$sum_row" ]; then
  IFS=, read -r blank comment code <<<"$sum_row"
else
  read -r blank comment code < <(
    printf '%s\n' "$csv" \
      | awk -F, 'NR > 1 && NF >= 5 && $1 != "SUM" && $2 != "SUM" {blank += $3; comment += $4; code += $5} END {printf "%d %d %d\n", blank, comment, code}'
  )
fi

total_lines="$((blank + comment + code))"
printf '\nTotal lines: %d (blank=%d, comment=%d, code=%d)\n' "$total_lines" "$blank" "$comment" "$code"
