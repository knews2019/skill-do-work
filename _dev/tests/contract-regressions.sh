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

# write_set parser lock-step (REQ-032). The dispatch gate and the board parser
# are one contract: work.md may only permit concurrent dispatch on a field the
# shipped parser actually reads, so neither half may be removed alone.
assert_contains \
  "actions/work.md" \
  'pairwise disjoint' \
  'actions/work.md Step 1 must keep the parallel-dispatch gate permitting concurrent REQs only when their write_sets are pairwise disjoint.'

assert_contains \
  "actions/work.md" \
  'serial-only' \
  'actions/work.md Step 1 must keep the serial-only resource-class rule that overrides write_set disjointness for ordered/generated resources.'

assert_contains \
  "tools/queue-kanban/model.go" \
  'fields\["write_set"\]' \
  'tools/queue-kanban/model.go must parse write_set in lock-step with the schema field work.md dispatches on.'

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

# ADR-017 memory engine contracts: the Stop capture must never block a session
# end, destructive consolidation must never be hook-wired, and the engine's
# core guardrails (cap, prompt-injection load, compose-don't-clobber install)
# must stay stated where agents read them.
assert_contains \
  "hooks/memory-stop-capture.sh" \
  'stop_hook_active' \
  'hooks/memory-stop-capture.sh must keep the stop_hook_active loop guard.'

assert_file_not_contains \
  "hooks/memory-stop-capture.sh" \
  '"decision":[[:space:]]*"block"' \
  'hooks/memory-stop-capture.sh must never emit a blocking decision — capture cannot hold a session open.'

for hooks_json_file in "hooks/hooks.json" "hooks/memory-hooks.json"; do
  assert_file_not_contains \
    "$hooks_json_file" \
    'dream' \
    "$hooks_json_file must never wire the destructive dream action to a hook."
done

assert_contains \
  "actions/memory.md" \
  '2,500' \
  'actions/memory.md must state the 2,500-character working-memory hard cap.'

assert_contains \
  "actions/memory.md" \
  'prompt-injection\.md' \
  'actions/memory.md recall must load the prompt-injection guardrail before reading hook-captured log content.'

assert_contains \
  "actions/install.md" \
  'memory-module' \
  'actions/install.md must carry the memory-module target (ADR-017).'

assert_contains \
  "actions/install.md" \
  'settings\.json\.pre-memory-module' \
  'actions/install.md memory-module hook merge must back up settings.json before composing entries.'

assert_contains \
  "hooks/memory-stop-capture.sh" \
  'REDACTED' \
  'hooks/memory-stop-capture.sh must redact credential-shaped text before persisting captures — defense in depth behind the machine-local store.'

# Raw captures and the per-machine ledger must never become committable: the installer
# adds them to .git/info/exclude (machine-local), never the project's .gitignore.
assert_contains \
  "actions/install.md" \
  '\*\*/memory/logs/' \
  'actions/install.md memory-module must add memory/logs/ to .git/info/exclude — verbatim captures must not be committable.'

assert_contains \
  "actions/install.md" \
  '\*\*/memory/usage-ledger\.jsonl' \
  'actions/install.md memory-module must add memory/usage-ledger.jsonl to .git/info/exclude.'

assert_contains \
  "hooks/memory-session-start.sh" \
  'session capture' \
  'hooks/memory-session-start.sh must strip raw session-capture sections from the startup injection — unvetted transcript text must not enter context before a prompt-injection guard can load.'

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

# Common Rationalizations regrowth ratchet (REQ-027). The four "earned" template
# sections (Rules / Common Rationalizations / Red Flags / Verification Checklist)
# drifted from "included when they'd help" to "included because the template listed
# them" — 20 of 42 action files carried all four (24 carry a Common Rationalizations
# table at all), most filled with generic engineering
# advice a capable model already follows. This check catches regrowth in Common
# Rationalizations specifically: a table whose rows carry no do-work-specific noun is
# exactly that generic filler (see CLAUDE.md → Action File Conventions for the full
# omission test). Scoped to files added after REQ-027 — the baseline below grandfathers
# the existing tree (as of REQ-027) so this check lands green without a mass rewrite in
# the same commit; REQ-025/028/029/030/031 clean up the backlog under the new rule.
common_rationalizations_baseline_action_files=(
  abandon.md ai-report.md bkb-reference.md bkb.md board.md capture.md clarify.md
  cleanup.md code-review.md commit.md deep-explore-reference.md deep-explore.md
  dream.md forensics.md help.md inspect.md install.md interview-reference.md
  interview.md kb-lessons-handoff.md note.md pipeline.md present-work.md
  prime-req-reservation.md prime.md prompts.md quick-wins.md reserve.md
  review-work.md roadmap.md sample-archived-req.md scan-ideas.md slop-check.md
  stray-check.md tidy-repo.md tutorial.md ui-review.md validate-feedback.md
  verify-requests.md version.md work-reference.md work.md
)

is_grandfathered_rationalizations_file() {
  local candidate_file_name="$1"
  local baseline_entry
  for baseline_entry in "${common_rationalizations_baseline_action_files[@]}"; do
    if [ "$baseline_entry" = "$candidate_file_name" ]; then
      return 0
    fi
  done
  return 1
}

# Illustrative, not exhaustive (Closed Enumerations Go Stale) — grep for do-work
# machinery vocabulary, not general software-engineering nouns. Split in two so "REQ"
# stays case-sensitive (case-insensitive it collides with "required"/"requires", which
# show up in ordinary generic advice and would silently defeat the check). The memory
# subsystem (working memory / daily log / usage ledger / bootstrap / Stop hook) is
# first-class do-work machinery too, so its nouns are recognized alongside the rest.
rationalization_noun_pattern_case_sensitive='\bREQ-?[0-9]*\b|\bUR-[0-9]+\b'
rationalization_noun_pattern_case_insensitive='queue|frontmatter|pipeline|\barchive|do-work|\bdomain\b|\bblocked\b|\bkb/|\bprime\b|\bclarify|working/|crew-member|\bschema\b|status:|working memory|daily log|\bledger\b|\bbootstrap\b|stop hook'

for action_file_path in "$repo_root"/actions/*.md; do
  action_file_name="$(basename "$action_file_path")"
  if is_grandfathered_rationalizations_file "$action_file_name"; then
    continue
  fi
  if ! grep -q '^## Common Rationalizations' "$action_file_path"; then
    continue
  fi
  rationalization_block="$(awk '/^## Common Rationalizations/{flag=1;next}/^## /{flag=0}flag' "$action_file_path" || true)"
  rationalization_rows="$(grep '^|' <<<"$rationalization_block" | grep -viE "if you're thinking" | grep -vE '^\|[-: |]+\|$' || true)"
  if [ -z "$rationalization_rows" ]; then
    continue
  fi
  if ! grep -qE "$rationalization_noun_pattern_case_sensitive" <<<"$rationalization_rows" \
     && ! grep -qiE "$rationalization_noun_pattern_case_insensitive" <<<"$rationalization_rows"; then
    printf 'FAIL: %s Common Rationalizations table has no do-work-specific noun (REQ, UR, queue, frontmatter, pipeline, archive, domain, blocked, kb/, prime, clarify, working/, crew-member, schema, status, working memory, daily log, ledger, bootstrap, stop hook — illustrative list) in any row — every row reads as generic engineering advice a capable model already follows. Add rows naming a specific do-work failure mode, or omit the section entirely (see CLAUDE.md -> Action File Conventions for the omission test).\n' \
      "$action_file_name" >&2
    fail_count=$((fail_count + 1))
  fi
done

# Stale-lock take-over must re-judge staleness inside the serialized mutex. The
# branch is selected on an unserialized read, and a holder that heartbeats in the
# read-to-write gap must not be overwritten as stale — the take-over nulls the
# holder's claim, so Crash Recovery then strips and re-queues the REQ the live
# session was mid-way through building (validate-feedback 2026-07-28 TOCTOU).
assert_contains \
  "actions/work-reference.md" \
  'Re-validate staleness inside the mutex' \
  'actions/work-reference.md stale take-over must re-validate staleness inside the mutex, not just serialize the write.'

assert_contains \
  "actions/work-reference.md" \
  'recompute the holder.s age from the fresh' \
  'actions/work-reference.md stale take-over must recompute the holder age from the fresh in-mutex read before overwriting.'

# Capture redaction must run on the FULL extracted text before any byte-budget
# cut: every credential pattern needs a complete token shape, so redacting after
# truncation lets a severed token (`ghp_1234567`) persist as an unmatched
# fragment (validate-feedback 2026-07-28, reproduced).
assert_contains \
  "actions/memory-reference.md" \
  'redaction runs BEFORE truncation' \
  'actions/memory-reference.md Stop-Capture spec must order redaction before truncation.'

if command -v jq &>/dev/null; then
  redaction_order_workdir="$(mktemp -d)"
  mkdir -p "$redaction_order_workdir/memory/logs"
  # Pad so the GitHub token straddles the side-truncation cut: with a short
  # assistant side the user side is cut at (1500-15) - 11 - 12 bytes, which
  # lands mid-token for a token starting at byte 1451.
  redaction_order_pad="$(printf 'A%.0s' $(seq 1 1450))"
  printf '{"type":"user","message":{"content":"%sghp_ABCDEFGHIJKLMNOPQRSTUVWXYZ012345 tail"}}\n{"type":"assistant","message":{"content":[{"type":"text","text":"shortreply"}]}}\n' \
    "$redaction_order_pad" > "$redaction_order_workdir/transcript.jsonl"
  printf '{"transcript_path":"%s/transcript.jsonl"}' "$redaction_order_workdir" \
    | CLAUDE_PROJECT_DIR="$redaction_order_workdir" bash "$repo_root/hooks/memory-stop-capture.sh" >/dev/null 2>&1 || true
  if ! ls "$redaction_order_workdir/memory/logs/"*.md >/dev/null 2>&1; then
    printf 'FAIL: hooks/memory-stop-capture.sh redaction-order probe wrote no capture — the probe transcript should be captured, not dropped.\n' >&2
    fail_count=$((fail_count + 1))
  elif grep -rqE 'ghp_[A-Za-z0-9_]+' "$redaction_order_workdir/memory/logs/" 2>/dev/null; then
    printf 'FAIL: hooks/memory-stop-capture.sh persisted a truncation-severed credential fragment — redaction must run on the full extracted text before any byte-budget cut.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  rm -rf "$redaction_order_workdir"
fi

# Worktree dispatch mode (REQ-033). The mode is only safe because three pieces
# hold together: a name that correlates a leftover with its REQ, a merged-ness
# assertion that is never forced, and a verification of the *merged* state before
# archive. Each can be quietly dropped without breaking anything visible, so pin
# them individually — plus cleanup's ownership of the unmerged leftovers, without
# which a crashed parallel run leaves branches nothing in the skill ever removes.
assert_contains \
  "actions/work-reference.md" \
  'worktree-agent-REQ-' \
  'actions/work-reference.md must keep the worktree/branch naming convention that correlates a leftover with its REQ id.'

assert_contains \
  "actions/work-reference.md" \
  'git branch -d' \
  'actions/work-reference.md must keep git branch -d (never -D) as the free merged-ness assertion on worktree cleanup.'

assert_contains \
  "actions/work.md" \
  '[Pp]ost-merge verification' \
  'actions/work.md Step 8 must gate archival on re-running the REQ acceptance checks against the merged tree.'

assert_contains \
  "actions/cleanup.md" \
  'worktree-agent-' \
  'actions/cleanup.md must keep the consent-gated orphaned-worktree pass that owns unmerged builder branches.'

# Multi-claim lock representation (REQ-035). The orchestrator lock's claim became a
# `claimed_reqs` LIST so one orchestrator can hold N concurrent claims, with
# `claimed_req` retained as a derived legacy mirror (claimed_reqs[0]). The gate text,
# lock schema, heartbeat rule, Crash Recovery gate, and cleanup's live-claim gate must
# all tell one story, so pin the canonical field's presence in each of the three files
# that read or write it — drop it from any one and the "one story" breaks silently.
# Also pin that the Crash Recovery gate skips files claimed by ANY fresh claim set,
# this session's own INCLUDED (freshness alone gates): the correctness fix that stops a
# Step 10 -> Step 1 loop from re-queuing a co-dispatched sibling mid-build.
assert_contains \
  "actions/work.md" \
  'claimed_reqs' \
  'actions/work.md must carry the canonical claimed_reqs list (the multi-claim field the parallel-dispatch gate and Step 2/8 bookkeeping depend on).'

assert_contains \
  "actions/work-reference.md" \
  'claimed_reqs' \
  'actions/work-reference.md must carry the canonical claimed_reqs list across the lock schema, heartbeat rule, and Crash Recovery gate.'

assert_contains \
  "actions/cleanup.md" \
  'claimed_reqs' \
  'actions/cleanup.md Pass 0 live-claim gate must exempt every id in a live session claimed_reqs, not just the legacy single claimed_req.'

assert_contains \
  "actions/work-reference.md" \
  'including this session'\''s own' \
  'actions/work-reference.md Crash Recovery gate must skip files claimed by ANY fresh claim set including this session own, so a Step 10 to Step 1 loop does not re-queue a co-dispatched sibling mid-build.'

# The primary action file's compact Crash Recovery summary is a separate restatement of
# the gate above and is NOT covered by the work-reference.md assertion — REQ-035 first
# shipped with work.md's summary still telling the old "another live session" story while
# work-reference.md:205 told the new one (caught in adversarial review). Pin the same-story
# phrasing in work.md too so the two files cannot silently diverge again.
assert_contains \
  "actions/work.md" \
  'this session'\''s own co-dispatched claims' \
  'actions/work.md Step 1 Crash Recovery summary must tell the same story as the work-reference.md gate — it also skips this session own co-dispatched claims on a Step 10 to Step 1 loop, not just another session claims.'

# The proceed-anyway option restates the same gate a THIRD time, and the file-wide
# assertion above cannot see it — a match anywhere in work-reference.md satisfies that
# one, so the Crash Recovery gate's own wording masked this restatement keeping the
# pre-REQ-035 "another live session" story right through REQ-035 (REQ-044). An agent
# reading only its local instruction on that path strips its own co-dispatched
# siblings, so pin the block itself: the current story in, the stale story out.
proceed_anyway_block="$(sed -n '/\*\*(a) Proceed anyway\*\*/,/\*\*(b) Take over\*\*/p' "$repo_root/actions/work-reference.md")"

assert_block_contains \
  "$proceed_anyway_block" \
  'including this session'\''s own' \
  'actions/work-reference.md proceed-anyway option must restate the gate as skipping every fresh claim including this session own, not just another live session claims.'

assert_block_not_contains \
  "$proceed_anyway_block" \
  'skips only files actively claimed by another live session' \
  'actions/work-reference.md proceed-anyway option must not reintroduce the pre-REQ-035 another-live-session-only gate wording, which tells a coexisting session to strip its own co-dispatched siblings.'

# Dispatch re-validation completeness (REQ-045). REQ-036 added the Step 5.5
# disjointness re-check and shipped no guard for it, and its coverage claim had a
# route-shaped hole: `route` is not assigned until Step 3, and Route A never reaches
# Step 5.5, so a co-dispatched Route A builder wrote under an unvalidated capture hint.
# Every co-dispatched REQ now has exactly one post-dispatch validation point — Step 5.5
# for Routes B/C, Step 3 for Route A — plus a named loser and a partition written to
# frontmatter so the re-check compares against the subset the gate actually issued.
# Each piece deletes cleanly without breaking anything visible, and the file-wide
# `pairwise disjoint`/`write_set` assertions above cannot see them (a match anywhere in
# actions/work.md satisfies those — the masking REQ-044 hit), so scope each to its step.
step_5_5_scope_block="$(sed -n '/^### Step 5.5: Scope Declaration/,/^### Step 5.75:/p' "$repo_root/actions/work.md")"
triage_step_block="$(sed -n '/^### Step 3: Triage/,/^### Step 3.5:/p' "$repo_root/actions/work.md")"
parallel_dispatch_block="$(sed -n '/^\*\*Parallel dispatch (optional/,/^\*\*Serial-only resource classes/p' "$repo_root/actions/work.md")"

assert_block_contains \
  "$step_5_5_scope_block" \
  'pairwise disjointness against every other in-flight REQ' \
  'actions/work.md Step 5.5 must keep the REQ-036 re-validation clause — the firmed Scope list is re-checked for pairwise disjointness against every other in-flight REQ before the mirror replaces the field.'

assert_block_contains \
  "$step_5_5_scope_block" \
  'The REQ at this re-check is the loser' \
  'actions/work.md Step 5.5 must name which REQ of an overlapping pair is serialized — with the loser undefined an orchestrator can hold a sibling that has already written under the boundary it was handed.'

assert_block_contains \
  "$triage_step_block" \
  'only post-dispatch validation point' \
  'actions/work.md Step 3 must re-validate a co-dispatched Route A REQ write-set here — Route A skips Step 5.5, so without this the Step 1 gate coverage claim has a route-shaped hole and a Route A builder writes under an unvalidated capture hint.'

assert_block_contains \
  "$parallel_dispatch_block" \
  '[Ww]rite that subset into the REQ.s `write_set` frontmatter' \
  'actions/work.md Step 1 must persist a partition directive into the REQ write_set frontmatter at dispatch — a partition living only in the dispatch prompt is invisible to the later re-checks, which then serialize the partition the gate itself issued.'

# Board overlap annotation stays display-only (REQ-052, follow-up to REQ-034). The
# invariant is structural: `annotateWriteSetOverlap` runs AFTER `bucketColumns` in
# buildBoard, so no column-placement code can read the annotation. That ordering is one
# movable line — a future edit that hoists the call above bucketing (to "reuse" the
# overlap set while placing cards) turns a display badge into scheduling logic, and the
# co-dispatch decision belongs to actions/work.md Step 1's gate, not the viewer. Assert
# the real call-site order inside buildBoard rather than a comment: a comment survives
# the move. Paired with the instruction-side claim in actions/board.md, whose rewording
# is the other way the invariant quietly stops being a promise.
build_board_block="$(sed -n '/^func buildBoard(/,/^}/p' "$repo_root/tools/queue-kanban/model.go")"
bucket_columns_call_line="$(grep -nF 'bucketColumns(board.AllRequests' <<<"$build_board_block" | head -1 | cut -d: -f1 || true)"
overlap_annotation_call_line="$(grep -nF 'annotateWriteSetOverlap(board.AllRequests)' <<<"$build_board_block" | head -1 | cut -d: -f1 || true)"

if [ -z "$bucket_columns_call_line" ] || [ -z "$overlap_annotation_call_line" ]; then
  printf 'FAIL: tools/queue-kanban/model.go buildBoard must call both bucketColumns(board.AllRequests ...) and annotateWriteSetOverlap(board.AllRequests) — one call site was renamed or removed, so the display-only ordering of the write-set overlap annotation can no longer be verified. Fix: restore the call sites (annotation last) or update this anchor in the same commit.\n' >&2
  fail_count=$((fail_count + 1))
elif [ "$overlap_annotation_call_line" -lt "$bucket_columns_call_line" ]; then
  printf 'FAIL: tools/queue-kanban/model.go calls annotateWriteSetOverlap BEFORE bucketColumns in buildBoard — the write-set overlap annotation is display only and must stay structurally unable to affect column placement. Fix: move the annotateWriteSetOverlap call back below bucketColumns; if the board really must schedule on write sets, that decision belongs to actions/work.md Step 1 dispatch gate.\n' >&2
  fail_count=$((fail_count + 1))
fi

board_rules_block="$(sed -n '/^## Rules/,/^## Common Rationalizations/p' "$repo_root/actions/board.md")"

assert_block_contains \
  "$board_rules_block" \
  '`annotateWriteSetOverlap` in `tools/queue-kanban/model\.go` runs \*after\* column bucketing' \
  'actions/board.md Rules must keep the claim that annotateWriteSetOverlap runs after column bucketing — the ordering ratchet above pins the code, this pins the promise a parser-editing agent reads. Fix: restate it in the parser lock-step rule.'

assert_block_contains \
  "$board_rules_block" \
  'never column logic' \
  'actions/board.md Rules must keep the overlap annotation display-only claim (drives the overlaps badge and drawer row, never column logic, never blocking) — without it nothing tells a parser-editing agent that which REQs may co-dispatch stays with actions/work.md Step 1 gate.'

if [ "$fail_count" -gt 0 ]; then
  exit 1
fi

printf 'Contract regression checks passed.\n'
