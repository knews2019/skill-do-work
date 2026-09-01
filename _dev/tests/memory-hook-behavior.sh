#!/usr/bin/env bash
# Consumer probes for the retained memory hook launchers and canonical Go protocols.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
fixture_root="$(mktemp -d)" || exit 1
trap 'rm -rf "$fixture_root"' EXIT
failure_count=0

skills_root="$fixture_root/project/.claude/skills"
core_root="$skills_root/do-work"
knowledge_root="$skills_root/do-work-knowledge"
mkdir -p "$core_root/tools" "$knowledge_root/hooks" "$fixture_root/project/memory/logs"
cp "$repo_root/skills/do-work/tools/do-work-cli.sh" "$core_root/tools/"
cp -R "$repo_root/skills/do-work/tools/do-work-cli" "$core_root/tools/"
rm -f "$core_root/tools/do-work-cli/do-work-cli"
cp "$repo_root/skills/do-work-knowledge/hooks/memory-session-start.sh" "$knowledge_root/hooks/"
cp "$repo_root/skills/do-work-knowledge/hooks/memory-stop-capture.sh" "$knowledge_root/hooks/"

absent_output="$(CLAUDE_PROJECT_DIR="$fixture_root/project" bash "$knowledge_root/hooks/memory-session-start.sh" 2>"$fixture_root/absent.err")"
if [ -n "$absent_output" ] || [ -s "$fixture_root/absent.err" ]; then
  printf 'FAIL: absent memory store was not silent.\n' >&2
  failure_count=$((failure_count + 1))
fi

printf 'standing memory\n' > "$fixture_root/project/memory/working-memory.md"
today="$(date -u +%F)"
printf '## 09:00 UTC note\nkeep\n\n## 10:00 UTC session capture abcdef01\n\n<!-- do-work:capture-body quoted -->\n> raw secret\n\n## 11:00 UTC note\nkeep too\n' > "$fixture_root/project/memory/logs/$today.md"
start_output="$(CLAUDE_PROJECT_DIR="$fixture_root/project" bash "$knowledge_root/hooks/memory-session-start.sh" 2>"$fixture_root/start.err")"
if ! grep -q 'standing memory' <<<"$start_output" || ! grep -q 'keep too' <<<"$start_output" || grep -q 'raw secret' <<<"$start_output" || [ -s "$fixture_root/start.err" ]; then
  printf 'FAIL: memory SessionStart projection/filtering drifted.\n' >&2
  failure_count=$((failure_count + 1))
fi

printf '{"type":"user","message":{"content":"hello ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZ"}}\n{"type":"assistant","message":{"content":[{"type":"text","text":"answer"}]}}\n' > "$fixture_root/project/transcript.jsonl"
printf '{"transcript_path":"%s"}' "$fixture_root/project/transcript.jsonl" \
  | CLAUDE_PROJECT_DIR="$fixture_root/project" bash "$knowledge_root/hooks/memory-stop-capture.sh" >"$fixture_root/stop.out" 2>"$fixture_root/stop.err"
capture_count="$(grep -c 'session capture ' "$fixture_root/project/memory/logs/$today.md" || true)"
if [ -s "$fixture_root/stop.out" ] || [ -s "$fixture_root/stop.err" ] || [ "$capture_count" -ne 2 ] || grep -q 'ghp_' "$fixture_root/project/memory/logs/$today.md"; then
  printf 'FAIL: memory Stop capture output, append, or redaction drifted.\n' >&2
  failure_count=$((failure_count + 1))
fi
printf '{"transcript_path":"%s"}' "$fixture_root/project/transcript.jsonl" \
  | CLAUDE_PROJECT_DIR="$fixture_root/project" bash "$knowledge_root/hooks/memory-stop-capture.sh" >/dev/null 2>&1
duplicate_count="$(grep -c 'session capture ' "$fixture_root/project/memory/logs/$today.md" || true)"
if [ "$duplicate_count" -ne "$capture_count" ]; then
  printf 'FAIL: duplicate Stop event appended the same capture twice.\n' >&2
  failure_count=$((failure_count + 1))
fi

missing_root="$fixture_root/missing/.claude/skills/do-work-knowledge"
mkdir -p "$missing_root/hooks"
cp "$repo_root/skills/do-work-knowledge/hooks/memory-stop-capture.sh" "$missing_root/hooks/"
missing_status=0
CLAUDE_PROJECT_DIR="$fixture_root/missing" bash "$missing_root/hooks/memory-stop-capture.sh" </dev/null >"$fixture_root/missing.out" 2>"$fixture_root/missing.err" || missing_status=$?
if [ "$missing_status" -ne 0 ] || ! grep -q 'launcher is missing' "$fixture_root/missing.err"; then
  printf 'FAIL: Stop launcher did not remain nonblocking and actionable without the CLI.\n' >&2
  failure_count=$((failure_count + 1))
fi

if grep -Eq 'jq|sed |awk |sha256|redact_credentials|CAPTURE_TEXT' "$knowledge_root/hooks/"*.sh; then
  printf 'FAIL: retained memory hook paths still contain domain logic.\n' >&2
  failure_count=$((failure_count + 1))
fi

[ "$failure_count" -eq 0 ] || exit 1
printf 'Memory hook behavior probes passed.\n'
