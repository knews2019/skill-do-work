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

# Retained jq/slurp accepted blank JSONL separators. Pin that protocol plus the
# legacy redirection modes under ordinary and restrictive umasks.
printf '{"type":"user","message":{"content":"blank-separated"}}\n\n{"type":"assistant","message":{"content":"still captured"}}\n' > "$fixture_root/project/transcript-blank.jsonl"
printf '{"transcript_path":"%s"}' "$fixture_root/project/transcript-blank.jsonl" \
  | (umask 022; CLAUDE_PROJECT_DIR="$fixture_root/project" bash "$knowledge_root/hooks/memory-stop-capture.sh") >"$fixture_root/blank.out" 2>"$fixture_root/blank.err"
if [ -s "$fixture_root/blank.out" ] || [ -s "$fixture_root/blank.err" ] || ! grep -q 'blank-separated' "$fixture_root/project/memory/logs/$today.md"; then
  printf 'FAIL: jq-accepted blank-line transcript behavior drifted.\n' >&2
  failure_count=$((failure_count + 1))
fi

restrictive_project="$fixture_root/restrictive-project"
mkdir -p "$restrictive_project/memory/logs"
printf '{"type":"user","message":{"content":"restricted"}}\n{"type":"assistant","message":{"content":"mode"}}\n' > "$restrictive_project/transcript.jsonl"
printf '{"transcript_path":"%s"}' "$restrictive_project/transcript.jsonl" \
  | (umask 077; CLAUDE_PROJECT_DIR="$restrictive_project" bash "$knowledge_root/hooks/memory-stop-capture.sh") >/dev/null 2>"$fixture_root/restrictive.err"
restrictive_log="$restrictive_project/memory/logs/$today.md"
if [ -s "$fixture_root/restrictive.err" ] || [ "$(stat -f '%Lp' "$restrictive_log" 2>/dev/null || stat -c '%a' "$restrictive_log" 2>/dev/null)" != 600 ] || [ "$(stat -f '%Lp' "$restrictive_project/memory/usage-ledger.jsonl" 2>/dev/null || stat -c '%a' "$restrictive_project/memory/usage-ledger.jsonl" 2>/dev/null)" != 600 ]; then
  printf 'FAIL: restrictive-umask capture/ledger modes drifted.\n' >&2
  failure_count=$((failure_count + 1))
fi
ordinary_project="$fixture_root/ordinary-project"
mkdir -p "$ordinary_project/memory/logs"
cp "$restrictive_project/transcript.jsonl" "$ordinary_project/transcript.jsonl"
printf '{"transcript_path":"%s"}' "$ordinary_project/transcript.jsonl" \
  | (umask 022; CLAUDE_PROJECT_DIR="$ordinary_project" bash "$knowledge_root/hooks/memory-stop-capture.sh") >/dev/null 2>"$fixture_root/ordinary.err"
ordinary_log="$ordinary_project/memory/logs/$today.md"
if [ -s "$fixture_root/ordinary.err" ] || [ "$(stat -f '%Lp' "$ordinary_log" 2>/dev/null || stat -c '%a' "$ordinary_log" 2>/dev/null)" != 644 ] || [ "$(stat -f '%Lp' "$ordinary_project/memory/usage-ledger.jsonl" 2>/dev/null || stat -c '%a' "$ordinary_project/memory/usage-ledger.jsonl" 2>/dev/null)" != 644 ]; then
  printf 'FAIL: ordinary-umask capture/ledger modes drifted.\n' >&2
  failure_count=$((failure_count + 1))
fi

# One typed observation must retain the append effects while Stop text remains
# exactly empty. A malformed/invalid transcript and append failure remain no-op.
typed_project="$fixture_root/typed-project"
mkdir -p "$typed_project/memory/logs"
printf '{"type":"user","message":{"content":"typed"}}\n{"type":"assistant","message":{"content":"effects"}}\n' > "$typed_project/transcript.jsonl"
typed_json="$(printf '{"transcript_path":"%s"}' "$typed_project/transcript.jsonl" | CLAUDE_PROJECT_DIR="$typed_project" bash "$core_root/tools/do-work-cli.sh" --repo-root "$typed_project" --format json memory-stop-capture)"
if ! grep -q '"protocol_output": ""' <<<"$typed_json" || [ "$(grep -o '"kind": "appended"' <<<"$typed_json" | wc -l | tr -d ' ')" -ne 2 ]; then
  printf 'FAIL: JSON projection lost typed capture and ledger effects.\n' >&2
  failure_count=$((failure_count + 1))
fi
append_failure_project="$fixture_root/append-failure"
mkdir -p "$append_failure_project/memory/logs/$today.md"
cp "$typed_project/transcript.jsonl" "$append_failure_project/transcript.jsonl"
append_failure_json="$(printf '{"transcript_path":"%s"}' "$append_failure_project/transcript.jsonl" | CLAUDE_PROJECT_DIR="$append_failure_project" bash "$core_root/tools/do-work-cli.sh" --repo-root "$append_failure_project" --format json memory-stop-capture)"
if ! grep -q 'MEMORY-CAPTURE-APPEND-SKIPPED' <<<"$append_failure_json" || ! grep -q '"protocol_output": ""' <<<"$append_failure_json"; then
  printf 'FAIL: append refusal did not remain silent while retaining typed evidence.\n' >&2
  failure_count=$((failure_count + 1))
fi
capture_count="$(grep -c 'session capture ' "$fixture_root/project/memory/logs/$today.md" || true)"
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

# Stop launcher failure matrix: too-old Go and build failure stay nonblocking and
# each produce one actionable diagnostic rather than stacked launcher messages.
run_stop_launcher_failure_case() {
  local case_name="$1" go_body="$2" expected="$3"
  local case_skills="$fixture_root/$case_name/.claude/skills" fake_bin="$fixture_root/$case_name-bin" status=0 line_count=0
  mkdir -p "$case_skills/do-work/tools" "$case_skills/do-work-knowledge/hooks" "$fake_bin"
  cp "$repo_root/skills/do-work/tools/do-work-cli.sh" "$case_skills/do-work/tools/"
  cp -R "$repo_root/skills/do-work/tools/do-work-cli" "$case_skills/do-work/tools/"
  rm -f "$case_skills/do-work/tools/do-work-cli/do-work-cli"
  cp "$repo_root/skills/do-work-knowledge/hooks/memory-stop-capture.sh" "$case_skills/do-work-knowledge/hooks/"
  printf '%b' "$go_body" > "$fake_bin/go"
  chmod +x "$fake_bin/go"
  PATH="$fake_bin:/usr/bin:/bin" CLAUDE_PROJECT_DIR="$fixture_root/$case_name" bash "$case_skills/do-work-knowledge/hooks/memory-stop-capture.sh" </dev/null >"$fixture_root/$case_name.out" 2>"$fixture_root/$case_name.err" || status=$?
  line_count="$(wc -l < "$fixture_root/$case_name.err" | tr -d ' ')"
  if [ "$status" -ne 0 ] || [ -s "$fixture_root/$case_name.out" ] || [ "$line_count" -ne 1 ] || ! grep -q "$expected" "$fixture_root/$case_name.err"; then
    printf 'FAIL: %s Stop boundary status=%s lines=%s stderr=<%s>.\n' "$case_name" "$status" "$line_count" "$(tr '\n' ' ' < "$fixture_root/$case_name.err")" >&2
    failure_count=$((failure_count + 1))
  fi
}
run_stop_launcher_failure_case stop-go-too-old '#!/usr/bin/env bash\nprintf "go version go1.24.0 fixture\\n"\n' 'Go 1.25.0 or newer is required'
run_stop_launcher_failure_case stop-build-failure '#!/usr/bin/env bash\nif [ "${1:-}" = version ]; then printf "go version go1.25.0 fixture\\n"; exit 0; fi\nexit 1\n' 'build failed; stale output was not run'

if grep -Eq 'jq|sed |awk |sha256|redact_credentials|CAPTURE_TEXT' "$knowledge_root/hooks/"*.sh; then
  printf 'FAIL: retained memory hook paths still contain domain logic.\n' >&2
  failure_count=$((failure_count + 1))
fi

[ "$failure_count" -eq 0 ] || exit 1
printf 'Memory hook behavior probes passed.\n'
