#!/usr/bin/env bash
set -euo pipefail

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "must run inside a git repository" >&2
  exit 1
fi

if ! command -v codex >/dev/null 2>&1; then
  echo "codex CLI is required but was not found in PATH" >&2
  exit 1
fi

diff_context=""
if ! git diff --cached --quiet; then
  diff_context="staged"
  file_list="$(git diff --cached --name-status --)"
  diff_text="$(git diff --cached --no-color --)"
elif ! git diff --quiet; then
  diff_context="unstaged"
  file_list="$(git diff --name-status --)"
  diff_text="$(git diff --no-color --)"
else
  echo "no staged or unstaged changes found" >&2
  exit 1
fi

# Keep prompt size bounded for large diffs.
max_chars=120000
if [ "${#diff_text}" -gt "$max_chars" ]; then
  diff_text="${diff_text:0:$max_chars}
[diff truncated after ${max_chars} characters]"
fi

prompt_file="$(mktemp)"
out_file="$(mktemp)"
err_file="$(mktemp)"
cleanup() {
  rm -f "$prompt_file" "$out_file" "$err_file"
}
trap cleanup EXIT

cat >"$prompt_file" <<EOF
Generate a git commit subject line from the diff below.

Requirements:
- Output exactly one line and nothing else.
- Use Conventional Commits style: type(scope): summary or type: summary.
- Keep it <= 72 characters.
- Use imperative mood.
- Do not wrap in quotes or backticks.

Diff source: ${diff_context}

Changed files:
${file_list}

Patch:
${diff_text}
EOF

if ! codex exec --ephemeral --sandbox read-only -o "$out_file" - <"$prompt_file" >/dev/null 2>"$err_file"; then
  cat "$err_file" >&2
  exit 1
fi

commit_name="$(sed -n '1p' "$out_file" | tr -d '\r' | sed 's/^[[:space:]]*//; s/[[:space:]]*$//')"
commit_name="$(printf '%s' "$commit_name" | sed -E "s/^[\"'\`]+//; s/[\"'\`]+$//")"

if [ -z "$commit_name" ]; then
  echo "failed to generate a commit name" >&2
  exit 1
fi

echo "$commit_name"
