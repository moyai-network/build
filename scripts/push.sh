#!/usr/bin/env bash
set -euo pipefail

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "must run inside a git repository" >&2
  exit 1
fi

git add -A

if git diff --cached --quiet; then
  echo "no local changes to commit; pushing current branch"
  git push
  exit 0
fi

commit_name="$(bash scripts/generate-commit-name.sh)"
echo "$commit_name"

git commit -m "$commit_name"
git push
