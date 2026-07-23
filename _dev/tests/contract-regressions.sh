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

assert_contains \
  "actions/capture.md" \
  'maintenance: false' \
  'actions/capture.md base REQ schema must carry maintenance:false so the marker is discoverable, not documented only for complex requests.'

assert_contains \
  "actions/capture.md" \
  'Maintenance assessment' \
  'actions/capture.md Step 1 must assess skill-instruction removal/narrowing and set the maintenance marker (work.md is marker-only and never infers it).'

assert_contains \
  "docs/ai-report-guide.md" \
  'completed-with-issues' \
  'docs/ai-report-guide.md must reflect the terminal-success set (completed | completed-with-issues), not only completed.'

assert_contains \
  "docs/ai-report-guide.md" \
  'DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND' \
  'docs/ai-report-guide.md must document agentic image backends as opt-in via the env flag, not the opportunistic default.'

assert_contains \
  "docs/cleanup-guide.md" \
  'completed-with-issues' \
  'docs/cleanup-guide.md sweep wording must include completed-with-issues in the terminal-status set.'

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

# Router word budget (REQ-020 ratchet). The 2026-07 bloat cleanup cut SKILL.md
# from 5,507 to 2,396 words by deleting duplicate enumerations of the action set
# (the router loads on EVERY invocation — words here tax all 30+ verbs). Budget =
# post-diet count + ~10% headroom. If you hit this limit: the fix is a merge or a
# lazy-load (see actions/help.md for the pattern), not a bigger budget — raise it
# only with an accompanying decisions/ note saying why routing itself had to grow.
router_word_budget=2650
router_word_count="$(wc -w < "$repo_root/SKILL.md")"
if [ "$router_word_count" -gt "$router_word_budget" ]; then
  printf 'FAIL: SKILL.md is %s words — over the %s-word router budget. Merge or lazy-load; do not grow the always-loaded router.\n' \
    "$router_word_count" "$router_word_budget" >&2
  fail_count=$((fail_count + 1))
fi

# Hardened checks (REQ-018): the work.md prose pointers and the shipped scripts
# must not drift apart — a pointer at a missing script silently un-hardens the step.
hardened_check_scripts=(
  "tools/checks/archive-collision.sh"
  "tools/checks/preflight.sh"
  "tools/checks/scope-drift.sh"
  "tools/checks/qualify.sh"
)

for check_script in "${hardened_check_scripts[@]}"; do
  if [ ! -x "$repo_root/$check_script" ]; then
    printf 'FAIL: %s must exist and be executable (work.md points at it).\n' "$check_script" >&2
    fail_count=$((fail_count + 1))
  fi
  assert_contains \
    "actions/work.md" \
    "$(basename "$check_script")" \
    "actions/work.md must reference $check_script — the hardened step's pointer was removed without un-hardening."
done

# Review regressions: prescribed shell and roadmap classification are runtime
# contracts even though they live in Markdown/just recipes rather than compiled code.
for kanban_recipe_file in "actions/install.md" "justfile"; do
  assert_file_not_contains \
    "$kanban_recipe_file" \
    'case "\$listener_command" in \*queue-kanban\*' \
    "$kanban_recipe_file must not identify a stale board from arbitrary argv text."
  assert_contains \
    "$kanban_recipe_file" \
    'lsof -a -p "\$listener_pid" -d txt -Fn' \
    "$kanban_recipe_file must identify a stale board from its executable, preserving cross-repo binary names without matching unrelated arguments."
done

assert_file_not_contains \
  "actions/work.md" \
  'else probe_wrapper=""' \
  'actions/work.md must not drop the blocked-check time limit when timeout/gtimeout is unavailable.'

assert_contains \
  "actions/work.md" \
  'probe_exit=124' \
  'actions/work.md must preserve a bounded portable fallback and report a timed-out blocked check as exit 124.'

blocked_probe_shell_block="$(sed -n '/^# Re-derive paths deterministically/,/^rm -f "\$BLOCKED_CHECK_SCRIPT"/p' "$repo_root/actions/work.md")"
if ! bash -n <<<"$blocked_probe_shell_block"; then
  printf 'FAIL: actions/work.md blocked-check shell block must remain syntactically valid.\n' >&2
  fail_count=$((fail_count + 1))
fi

assert_contains \
  "actions/roadmap.md" \
  '^-[[:space:]]+\*\*Ready\*\*[[:space:]]+— normalized `status` is `pending`' \
  'actions/roadmap.md must require pending status before classifying a queued REQ as Ready.'

# CLAUDE.md/AGENTS.md are the maintainer doc, export-ignored since 0.136.0 so they never
# land in consumer installs (nested CLAUDE.md is auto-loaded into consumer agents' context).
assert_contains \
  ".gitattributes" \
  '^/CLAUDE\.md[[:space:]]+export-ignore' \
  '.gitattributes must export-ignore /CLAUDE.md — the maintainer doc must not ship to consumer installs.'

assert_contains \
  ".gitattributes" \
  '^/AGENTS\.md[[:space:]]+export-ignore' \
  '.gitattributes must export-ignore /AGENTS.md — the redirect stub must not ship to consumer installs.'

# Shipped files must not cite the skill's own CLAUDE.md/AGENTS.md — those files are absent
# downstream, so a citation dangles. The idiom patterns are illustrative, not exhaustive
# (references to a *consumer project's* CLAUDE.md, like capture.md's prime routing, are fine);
# the full rule lives in CLAUDE.md → Action File Conventions.
shipped_citation_paths=(SKILL.md next-steps.md README.md actions crew-members prompts interviews specs docs hooks tools)
self_citation_pattern='(see|per|→) `?CLAUDE\.md|CLAUDE\.md`? *→|(see|per) `?AGENTS\.md'
self_citation_hits="$(cd "$repo_root" && grep -rIEn "$self_citation_pattern" "${shipped_citation_paths[@]}" 2>/dev/null || true)"
if [ -n "$self_citation_hits" ]; then
  printf 'FAIL: shipped files must not cite the skill'\''s own CLAUDE.md/AGENTS.md (export-ignored — absent in consumer installs). Restate the rule inline or point at a shipped home:\n%s\n' "$self_citation_hits" >&2
  fail_count=$((fail_count + 1))
fi

if [ "$fail_count" -gt 0 ]; then
  exit 1
fi

printf 'Contract regression checks passed.\n'
