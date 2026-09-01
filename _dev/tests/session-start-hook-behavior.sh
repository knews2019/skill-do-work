#!/usr/bin/env bash
# Behavioral regression probes for the retained core SessionStart launcher and Go authority.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
fixture_root="$(mktemp -d)" || exit 1
trap 'rm -rf "$fixture_root"' EXIT
failure_count=0

install_core_fixture() {
  local case_root="$1"
  mkdir -p "$case_root/skill/hooks" "$case_root/skill/actions" "$case_root/skill/tools"
  cp "$repo_root/skills/do-work/hooks/session-start.sh" "$case_root/skill/hooks/"
  cp "$repo_root/skills/do-work/tools/do-work-cli.sh" "$case_root/skill/tools/"
  cp -R "$repo_root/skills/do-work/tools/do-work-cli" "$case_root/skill/tools/"
  rm -f "$case_root/skill/tools/do-work-cli/do-work-cli"
}

run_banner_case() {
  local case_name="$1" version_text="$2" queue_count="$3" expected="$4"
  local case_root="$fixture_root/$case_name" request_index=0 output='' status=0
  install_core_fixture "$case_root"
  mkdir -p "$case_root/project/do-work/queue"
  printf '%b' "$version_text" > "$case_root/skill/actions/version.md"
  while [ "$request_index" -lt "$queue_count" ]; do
    request_index=$((request_index + 1))
    printf -- '---\nid: REQ-%03d\n---\n' "$request_index" > "$case_root/project/do-work/queue/REQ-$(printf '%03d' "$request_index")-fixture.md"
  done
  output="$(CLAUDE_PROJECT_DIR="$case_root/project" bash "$case_root/skill/hooks/session-start.sh" 2>"$case_root/stderr")" || status=$?
  if [ "$status" -ne 0 ] || [ "$output" != "$expected" ] || [ -s "$case_root/stderr" ]; then
    printf 'FAIL: %s status=%s output=<%s> stderr=<%s>\n' "$case_name" "$status" "$output" "$(tr '\n' ' ' < "$case_root/stderr")" >&2
    failure_count=$((failure_count + 1))
  fi
}

run_banner_case valid '**Current version**: 9.8.7\n' 2 "do-work v9.8.7 loaded. 2 pending REQ(s). Say 'do-work help' for commands."
run_banner_case missing '' 0 "do-work vunknown loaded. 0 pending REQ(s). Say 'do-work help' for commands."
run_banner_case reformatted '**Version**: 9.8.7\n' 0 "do-work vunknown loaded. 0 pending REQ(s). Say 'do-work help' for commands."
run_banner_case multiple '**Current version**: 1\n**Current version**: 2\n' 0 "do-work v1
2 loaded. 0 pending REQ(s). Say 'do-work help' for commands."

missing_root="$fixture_root/missing-launcher"
mkdir -p "$missing_root/skill/hooks" "$missing_root/project"
cp "$repo_root/skills/do-work/hooks/session-start.sh" "$missing_root/skill/hooks/"
missing_status=0
CLAUDE_PROJECT_DIR="$missing_root/project" bash "$missing_root/skill/hooks/session-start.sh" >"$missing_root/stdout" 2>"$missing_root/stderr" || missing_status=$?
if [ "$missing_status" -eq 0 ] || ! grep -q 'canonical launcher is missing' "$missing_root/stderr"; then
  printf 'FAIL: core SessionStart did not stop actionably when the canonical launcher was absent.\n' >&2
  failure_count=$((failure_count + 1))
fi

if grep -Eq 'cleanup-req-reservations\.sh|repair-req-timestamps\.sh|sed |awk |find ' "$repo_root/skills/do-work/hooks/session-start.sh"; then
  printf 'FAIL: retained SessionStart path contains domain logic instead of a thin launcher.\n' >&2
  failure_count=$((failure_count + 1))
fi

[ "$failure_count" -eq 0 ] || exit 1
printf 'SessionStart hook behavior probes passed.\n'
