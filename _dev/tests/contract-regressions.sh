#!/usr/bin/env bash
# shellcheck disable=SC2016
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
fail_count=0

assert_contains() {
  local file_path="$1"
  local pattern_text="$2"
  local message_text="$3"

  if ! grep -Eq "$pattern_text" "$repo_root/$file_path"; then
    printf 'FAIL: %s\n' "$message_text" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_block_contains() {
  local block_text="$1"
  local pattern_text="$2"
  local message_text="$3"

  if ! grep -Eq "$pattern_text" <<<"$block_text"; then
    printf 'FAIL: %s\n' "$message_text" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_block_not_contains() {
  local block_text="$1"
  local pattern_text="$2"
  local message_text="$3"

  if grep -Eq "$pattern_text" <<<"$block_text"; then
    printf 'FAIL: %s\n' "$message_text" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_file_missing() {
  local file_path="$1"
  local message_text="$2"

  if [ -e "$repo_root/$file_path" ]; then
    printf 'FAIL: %s\n' "$message_text" >&2
    fail_count=$((fail_count + 1))
  fi
}

assert_file_not_contains() {
  local file_path="$1"
  local pattern_text="$2"
  local message_text="$3"

  if grep -Eiq "$pattern_text" "$repo_root/$file_path"; then
    printf 'FAIL: %s\n' "$message_text" >&2
    grep -Ein "$pattern_text" "$repo_root/$file_path" >&2 || true
    fail_count=$((fail_count + 1))
  fi
}

skill_dispatch_block="$(sed -n '/^## Action Dispatch/,/^## Suggest Next Steps/p' "$repo_root/SKILL.md")"
work_archive_success_block="$(sed -n '/^### Step 8: Archive/,/^\*\*On failure:/p' "$repo_root/actions/work.md")"

assert_block_contains \
  "$skill_dispatch_block" \
  '^\| work[[:space:]]*\| `\./actions/work\.md`[[:space:]]*\| `\$ARGUMENTS`' \
  'SKILL.md must pass work arguments through so scoped REQ IDs and --wave are not dropped.'

assert_block_contains \
  "$work_archive_success_block" \
  'already (set to )?`completed-with-issues`|status is already `completed-with-issues`|preserve[^[:cntrl:]]*`completed-with-issues`' \
  'actions/work.md Archive success path must explicitly preserve completed-with-issues from failed remediation.'

assert_block_not_contains \
  "$work_archive_success_block" \
  '^1\. Update frontmatter: `status: completed`, `completed_at: <timestamp>`$' \
  'actions/work.md Archive success path must not unconditionally overwrite status with completed.'

assert_contains \
  "actions/ai-report.md" \
  'DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND' \
  'actions/ai-report.md must keep sandbox-bypassed agentic image generation behind an explicit opt-in.'

assert_contains \
  "actions/ai-report.md" \
  'mktemp -d' \
  'actions/ai-report.md must run any agentic image fallback from a locked temporary directory, not the repo cwd.'

assert_contains \
  "actions/ai-report.md" \
  'chmod 700' \
  'actions/ai-report.md must lock down the temporary image-generation directory before invoking an agentic backend.'

assert_contains \
  "actions/version.md" \
  'fresh upstream tarball|fresh upstream tree' \
  'actions/version.md update flow must compare against a freshly extracted upstream tarball before overwriting.'

assert_contains \
  "actions/version.md" \
  'diff -ru' \
  'actions/version.md update flow must prescribe a recursive pre-extraction diff against the fresh upstream tree.'

assert_file_not_contains \
  "actions/version.md" \
  'log -1 --format=%H -- actions/version\.md' \
  'actions/version.md must not use the last version.md-touching commit as the committed-customization baseline.'

assert_file_missing \
  "prompts/ultracode-fable-workflow.md" \
  'retired ultracode/fable prompt file must be removed from the active prompt library.'

active_runtime_docs=(
  "SKILL.md"
  "README.md"
  "next-steps.md"
  "actions/work.md"
  "actions/work-reference.md"
  "prompts/README.md"
)

for runtime_doc in "${active_runtime_docs[@]}"; do
  if [ -f "$repo_root/$runtime_doc" ]; then
    assert_file_not_contains \
      "$runtime_doc" \
      'ultracode|fable' \
      "active runtime doc $runtime_doc must not mention the retired ultracode/fable workflow."
  fi
done

if [ "$fail_count" -gt 0 ]; then
  exit 1
fi

printf 'Contract regression checks passed.\n'
