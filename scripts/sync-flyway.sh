#!/usr/bin/env bash
set -euo pipefail

repo="external/flyway"
commit_message="${SYNC_FLYWAY_COMMIT_MESSAGE:-chore: update files}"
parent_commit_message="${SYNC_FLYWAY_PARENT_COMMIT_MESSAGE:-$commit_message}"

if [ -n "$(git -C "$repo" ls-files -u)" ]; then
  echo "unresolved conflicts detected in $repo; resolve them first, then rerun sync-flyway"
  exit 1
fi

current_branch="$(git -C "$repo" symbolic-ref --short -q HEAD || true)"
branch="$current_branch"
wip_commit=""

if [ -z "$branch" ]; then
  branch="$(git -C "$repo" symbolic-ref --short refs/remotes/origin/HEAD 2>/dev/null | sed 's@^origin/@@')"
  if [ -z "$branch" ]; then
    branch="master"
  fi

  echo "detached HEAD in $repo; preparing $branch"

  dirty=0
  if ! git -C "$repo" diff --quiet || ! git -C "$repo" diff --cached --quiet || [ -n "$(git -C "$repo" ls-files --others --exclude-standard)" ]; then
    dirty=1
  fi

  if [ "$dirty" -eq 1 ]; then
    git -C "$repo" add -A
    git -C "$repo" commit -m "$commit_message"
    wip_commit="$(git -C "$repo" rev-parse HEAD)"
  fi

  if git -C "$repo" show-ref --verify --quiet "refs/heads/$branch"; then
    git -C "$repo" checkout "$branch"
  else
    git -C "$repo" checkout -b "$branch" --track "origin/$branch"
  fi

  git -C "$repo" pull --rebase --autostash origin "$branch"

  if [ -n "$wip_commit" ]; then
    if ! git -C "$repo" cherry-pick -X theirs "$wip_commit"; then
      echo "cherry-pick conflict in $repo; resolve it, run git -C $repo cherry-pick --continue, then git -C $repo push origin $branch"
      exit 1
    fi
  fi
else
  git -C "$repo" pull --rebase --autostash origin "$branch"
fi

git -C "$repo" add -A
if ! git -C "$repo" diff --cached --quiet; then
  git -C "$repo" commit -m "$commit_message"
fi

git -C "$repo" push origin "$branch"

# Commit only the submodule pointer update in the parent repository.
git add "$repo"
if ! git diff --cached --quiet -- "$repo"; then
  git commit -m "$parent_commit_message" -- "$repo"
fi

git push
