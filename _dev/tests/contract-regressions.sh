#!/usr/bin/env bash
# shellcheck disable=SC2016
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
fail_count=0
core_root="$repo_root/skills/do-work"
board_root="$repo_root/skills/do-work-board"
knowledge_root="$repo_root/skills/do-work-knowledge"
toolbox_root="$repo_root/skills/do-work-toolbox"

resolve_runtime_file() {
  local relative_path="$1"

  case "$relative_path" in
    SKILL.md|next-steps.md) printf '%s/%s\n' "$core_root" "$relative_path" ;;
    actions/board.md|docs/board-guide.md) printf '%s/%s\n' "$board_root" "$relative_path" ;;
    actions/bkb*|actions/dream.md|actions/interview*|actions/memory*|actions/prompts.md|actions/setup-memory.md|docs/bkb-guide.md|docs/dream-guide.md|docs/interview-guide.md|docs/prompts-guide.md|prompts/*|interviews/*|hooks/memory-*)
      printf '%s/%s\n' "$knowledge_root" "$relative_path" ;;
    actions/ai-report*|actions/code-review.md|actions/deep-explore*|actions/inspect.md|actions/install.md|actions/note.md|actions/present-work.md|actions/prime.md|actions/quick-wins.md|actions/scan-ideas.md|actions/slop-check.md|actions/stray-check.md|actions/tidy-repo.md|actions/tutorial.md|actions/ui-review.md|actions/validate-feedback.md|docs/ai-report-guide.md|docs/code-review-guide.md|docs/inspect-guide.md|docs/present-work-guide.md|docs/prime-guide.md|docs/quick-wins-guide.md|docs/slop-check-guide.md|docs/stray-check-guide.md|docs/ui-review-guide.md)
      printf '%s/%s\n' "$toolbox_root" "$relative_path" ;;
    actions/*|crew-members/*|docs/*|hooks/*|specs/*|tools/checks/*|tools/do-work-update.sh|tools/prime-do-work-update.md)
      printf '%s/%s\n' "$core_root" "$relative_path" ;;
    tools/queue-kanban/*) printf '%s/%s\n' "$board_root" "$relative_path" ;;
    *) printf '%s/%s\n' "$repo_root" "$relative_path" ;;
  esac
}

for cutover_export_path in VERSION suite skills; do
  if git -C "$repo_root" check-attr export-ignore -- "$cutover_export_path" \
    | grep -q 'export-ignore: set'; then
    printf 'FAIL: live modular archive still excludes /%s.\n' "$cutover_export_path" >&2
    fail_count=$((fail_count + 1))
  fi
done

for retired_runtime_path in SKILL.md next-steps.md actions crew-members docs hooks interviews prompts specs tools/checks tools/do-work-update.sh tools/queue-kanban tools/prime-do-work-update.md; do
  if [ -f "$repo_root/$retired_runtime_path" ] \
    || { [ -d "$repo_root/$retired_runtime_path" ] \
      && find "$repo_root/$retired_runtime_path" -type f \
        ! -path "$repo_root/tools/queue-kanban/queue-kanban" -print -quit \
        | grep -q .; }; then
    printf 'FAIL: legacy root runtime still exists after modular cutover: %s.\n' "$retired_runtime_path" >&2
    fail_count=$((fail_count + 1))
  fi
done

assert_contains() {
  local file_path="$1"
  local pattern_text="$2"
  local message_text="$3"

  if ! grep -Eq -- "$pattern_text" "$(resolve_runtime_file "$file_path")"; then
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

  local resolved_file
  resolved_file="$(resolve_runtime_file "$file_path")"
  if grep -Eiq -- "$pattern_text" "$resolved_file"; then
    printf 'FAIL: %s\n' "$message_text" >&2
    grep -Ein -- "$pattern_text" "$resolved_file" >&2 || true
    fail_count=$((fail_count + 1))
  fi
}

skill_dispatch_block="$(sed -n '/^## Routing/,/^## Dispatch/p' "$core_root/SKILL.md")"
work_archive_success_block="$(sed -n '/^### Step 8: Archive/,/^\*\*On failure:/p' "$core_root/actions/work.md")"

# Every phrase the version action documents must be reachable from the always-loaded
# router. Derive the phrase set from the action instead of maintaining a second alias
# inventory, then check the one precedence edge that matters: the exact update-check
# phrase must win before generic request verification's `check` alias.
if ! python3 - "$core_root/SKILL.md" "$core_root/actions/version.md" <<'PY'
import pathlib
import re
import sys

router_file = pathlib.Path(sys.argv[1])
version_action_file = pathlib.Path(sys.argv[2])
router_text = router_file.read_text()
version_action_text = version_action_file.read_text()

routing_block = router_text.split("## Routing", 1)[1].split("## Dispatch", 1)[0]
route_rows = []
for line_number, line in enumerate(routing_block.splitlines(), start=1):
    if not line.startswith("|") or "`./actions/" not in line:
        continue
    cells = [cell.strip() for cell in line.strip().strip("|").split("|")]
    if len(cells) != 2:
        continue
    aliases = re.findall(r"`([^`]+)`", cells[0])
    route_rows.append((line_number, aliases, cells[1]))

documented_phrases = []
for mode_name in ("Version request", "Update check", "Recap"):
    mode_match = re.search(
        rf"^- \*\*{re.escape(mode_name)}\*\* — (.+)$",
        version_action_text,
        flags=re.MULTILINE,
    )
    if mode_match is None:
        raise SystemExit(f"version action has no {mode_name} phrase declaration")
    documented_phrases.extend(re.findall(r'"([^"]+)"', mode_match.group(1)))

alias_routes = {}
for row_number, aliases, route in route_rows:
    for alias in aliases:
        alias_routes.setdefault(alias, []).append((row_number, route))

for phrase in dict.fromkeys(documented_phrases):
    matches = alias_routes.get(phrase, [])
    if len(matches) != 1 or "`./actions/version.md`" not in matches[0][1]:
        raise SystemExit(
            f"documented version phrase {phrase!r} must route exactly once to actions/version.md; "
            f"found {matches or 'no route'}"
        )

update_check_row = alias_routes["check for updates"][0][0]
generic_check_matches = alias_routes.get("check", [])
if len(generic_check_matches) != 1:
    raise SystemExit(f"generic check alias must route exactly once; found {generic_check_matches}")
if update_check_row >= generic_check_matches[0][0]:
    raise SystemExit("check for updates must route before the generic check alias")
PY
then
  printf 'FAIL: core version/update phrases drifted from first-match routing.\n' >&2
  fail_count=$((fail_count + 1))
fi

for retired_pipeline_path in \
  skills/do-work/actions/pipeline.md \
  skills/do-work/actions/pipeline-reference.md \
  skills/do-work/hooks/pipeline-guard.sh
do
  assert_file_missing \
    "$retired_pipeline_path" \
    "stateful pipeline runtime must remain retired: $retired_pipeline_path"
done

assert_block_not_contains \
  "$skill_dispatch_block" \
  '^\|[^|]*`(pipeline|full)`' \
  'Core SKILL.md must not retain a pipeline/full route or compatibility alias.'

assert_file_missing \
  "skills/do-work/actions/moved-command-shim.md" \
  'The one-release moved-command shim must be deleted after client migration.'

assert_block_not_contains \
  "$skill_dispatch_block" \
  'moved-command-shim\.md' \
  'Core routing must not retain compatibility rows for sibling-owned commands.'

assert_file_not_contains \
  "tools/do-work-update.sh" \
  'suite-layout-v2|--capabilities|legacy_shipped_paths|legacy all-in-one skill' \
  'The current updater must not retain bridge capability, monolith, or stale-copy branches.'

assert_file_not_contains \
  "tools/install-do-work-suite.sh" \
  '--migrate-legacy-do-work|\.claude/skills/do-work/hooks/memory-' \
  'The suite installer must not retain exact recipe or old core memory-hook migrations.'

assert_file_not_contains \
  "actions/setup-memory.md" \
  '\.claude/skills/do-work/hooks/memory-' \
  'Knowledge memory setup must describe only the current modular hook paths.'

assert_file_not_contains \
  ".gitignore" \
  'do-work/pipeline\.json' \
  'Root .gitignore must not recreate the retired pipeline state-file lifecycle.'

assert_file_not_contains \
  "hooks/session-start.sh" \
  'pipeline\.json|Pipeline active|do-work pipeline' \
  'SessionStart must report queue state without reading or advertising the retired pipeline state machine.'

assert_file_not_contains \
  "hooks/hooks.json" \
  'pipeline-guard\.sh' \
  'Fresh core hook settings must not install the retired pipeline Stop guard.'

if ! python3 - \
  "$repo_root/README.md" \
  "$core_root/actions/help.md" <<'PY'
import pathlib
import sys

approved_prompt = """Use the installed do-work suite to complete this request end to end:

1. Use do-work to capture the request below and record the resulting UR ID.
2. Run do-work verify-requests for that UR. Stop and report if verification fails.
3. Run the UR's REQs through do-work run. Require its built-in tests and review to pass.
4. Use do-work-toolbox present-work for the same UR.
5. Report the implementation, tests, decisions, and deliverable paths.

Request:
<paste request here>"""
missing = [path for path in sys.argv[1:] if approved_prompt not in pathlib.Path(path).read_text()]
if missing:
    raise SystemExit("approved full-cycle prompt is not byte-identical in: " + ", ".join(missing))
PY
then
  printf 'FAIL: README and core help must carry the approved UR-031 full-cycle prompt byte-for-byte.\n' >&2
  fail_count=$((fail_count + 1))
fi

assert_block_contains \
  "$skill_dispatch_block" \
  '^\| `run`[^|]*\| `\./actions/work\.md`' \
  'Core SKILL.md must route work triggers to actions/work.md so scoped REQ IDs and --wave reach the action input.'

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
  'tools/do-work-update\.sh.*--project-root' \
  'actions/version.md update flow must delegate mutation to the shared updater engine.'

assert_contains \
  "actions/version.md" \
  'do not duplicate its archive or overwrite logic here' \
  'actions/version.md must keep archive review and overwrite logic in the shared updater rather than a second agent-only implementation.'

assert_file_not_contains \
  "actions/version.md" \
  'log -1 --format=%H -- actions/version\.md' \
  'actions/version.md must not use the last version.md-touching commit as the committed-customization baseline.'

# No template may emit an instruction addressed to its reader (REQ-080). The Simple REQ template ended
# with a bare "Think carefully before answering." *inside* the fence, so every REQ capture produced
# inherited it — 25 archived REQs carry it, and one builder correctly flagged it as an instruction-like
# artifact and treated it as data (REQ-012 D-02) without anyone looking at the source that kept
# emitting it. A generated work item must contain the request and nothing that reads as a directive to
# whoever processes it; that shape is precisely what crew-members/prompt-injection.md exists to catch,
# and the skill must not manufacture it. Literal negatives on the known phrasings, deliberately not an
# English-grammar detector — the wordings are illustrative, the condition above is the rule.
for capture_artifact_phrase in \
  'Think carefully before answering' \
  'Reason step by step' \
  'Take your time'
do
  assert_file_not_contains \
    "actions/capture-reference.md" \
    "$capture_artifact_phrase" \
    "actions/capture-reference.md templates must not contain an instruction addressed to the reader — it lands inside the fence and is copied into every generated REQ, where it reads as a directive to the builder rather than as the user's request (REQ-080). Matched: $capture_artifact_phrase"
done

assert_contains \
  "actions/capture.md" \
  'maintenance: false' \
  'actions/capture.md base REQ schema must carry maintenance:false so the marker is discoverable, not documented only for complex requests.'

assert_contains \
  "actions/capture.md" \
  'Maintenance assessment' \
  'actions/capture.md Step 1 must assess skill-instruction removal/narrowing and set the maintenance marker (work.md is marker-only and never infers it).'

# One home for the timestamp command (REQ-078). The Timestamp rule in actions/work-reference.md is
# the only place in actions/ that spells a command for obtaining a stamp; every other site cites the
# rule. Eleven inline copies of `date -u +…` is why 0.167.0's new source-preference order — and the
# Windows form in particular — was unreachable from any site an agent actually follows: the citing
# site already handed it a POSIX command, so it never read the rule. Trigger condition, not a file
# list: any file under actions/ other than the rule's own home that spells a `date -u +…` invocation
# is a regression, whatever the format specifier. hooks/ is out of scope on purpose (executable POSIX
# shell, platform-specific by design), as is tools/ (Go source and user-facing strings, a different
# surface with a different tradeoff).
timestamp_command_copies="$(grep -rlE 'date -u \+%' "$core_root/actions" | grep -v 'work-reference\.md' || true)"
if [ -n "$timestamp_command_copies" ]; then
  printf 'FAIL: only actions/work-reference.md may spell a timestamp command; these action files inline a copy, so an agent following them never reaches the rule preference order and a Windows agent gets a command that does not exist on its box — cite the Timestamp rule instead (REQ-078):\n%s\n' \
    "$timestamp_command_copies" >&2
  fail_count=$((fail_count + 1))
fi

timestamp_rule_block="$(sed -n '/^\*\*Timestamp rule —/,/^```yaml/p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$timestamp_rule_block" \
  'ToUniversalTime' \
  'actions/work-reference.md Timestamp rule must name the .ToUniversalTime() form for Windows — Get-Date -AsUTC is PowerShell 7+ and a stock box ships 5.1 as powershell.exe, where that parameter is unrecognized and the call fails outright (REQ-078).'

assert_block_contains \
  "$timestamp_rule_block" \
  'powershell -NoProfile -Command' \
  'actions/work-reference.md Timestamp rule must give the cmd entry point explicitly — a bare cmdlet is not a command in cmd, which is the shell the Windows clause exists for (REQ-078).'

# write_set's display-only status must never be argued from a one-REQ-at-a-time premise (REQ-075).
# The conclusion is right and permanent — nothing schedules, gates, or dispatches on write_set — but
# REQ-073 falsified that premise: several builders can run at once under a single queue owner. A reader
# who follows the old reasoning concludes the opposite of the contract, since the stated cause is gone.
# The durable reason is that write_set is advisory input to a human's pick and the merge is the
# non-interference proof (actions/work-reference.md → Worktree Dispatch Mode → Fan-Out Dispatch), which
# holds at any builder count.
#
# Granularity is one line, and deliberately so: the premise is only dangerous sitting next to the field
# it purports to explain, and a proximity window false-positives on the canonical Fan-Out Dispatch
# section itself, two lines below the advisory-`write_set` bullet, where "integration runs one REQ at a
# time" is a true statement about integration. Prose files put both on one long line, so the sweep sees
# them; Go and JS comments wrap, so it cannot — those two files are covered by the file-level negatives
# below instead, which is exact because neither has any legitimate reason to state a builder count.
# `|| true` because grep exits 1 on no-match, which is the passing case and must not abort under `set -e`.
#
# TRIGGER CONDITION (REQ-079), not a phrase list: no shipped file may argue write_set's display-only
# status from ANY builder count, in either of the premise's two fingerprints — the thing it *said*
# ("one REQ at a time") and the thing it was *called* ("under the exclusive-session model"). The weak
# form is the more dangerous of the two, because the model it names is still true and only its
# relevance to this conclusion died, so the sentence survives inspection. Both patterns below enumerate
# the wordings seen so far and are ILLUSTRATIVE, NOT EXHAUSTIVE — a new phrasing earns a widened
# pattern, never a second copy of the block. The strong-form pattern is defined once and reused by all
# three of its consumers for the same reason.
#
# tools/queue-kanban/prime-do-kanban.md is deliberately NOT given a file-level negative: its REQ-075
# lesson entry quotes both fingerprints verbatim, which is the one legitimate reason to name them. The
# line sweep's write_set|overlaps filter is what keeps that line out of scope, and it must stay that way.
builder_count_premise_pattern='(one|a single|only one)( [a-z]+){0,2} (REQ|builder|coder|agent)s?( [a-z]+){0,3} (at a time|at once|concurrently)|(one|a single|only one)( [a-z]+){0,2} (REQ|builder|coder|agent)s? is ever (building|running|in flight)'
exclusive_session_premise_pattern='exclusive.session|no other .do-work. session'

stale_write_set_premise_lines="$(
  grep -rIEn "$builder_count_premise_pattern" \
    "$core_root/actions" "$core_root/docs" "$board_root/tools/queue-kanban" 2>/dev/null \
    | grep -E 'write_set|overlaps' || true
)"
if [ -n "$stale_write_set_premise_lines" ]; then
  printf 'FAIL: write_set is display-only at ANY builder count — no shipped file may justify that with a one-REQ-at-a-time premise (REQ-075, pattern widened by REQ-079). Point at actions/work-reference.md → Fan-Out Dispatch instead. Offending lines:\n' >&2
  printf '%s\n' "$stale_write_set_premise_lines" >&2
  fail_count=$((fail_count + 1))
fi

# The weak fingerprint of the same premise (REQ-079). REQ-075 named this form as the more dangerous one
# in its own lesson and then pinned only the strong one, so the guard could not catch the recurrence it
# was written for. Naming the exclusive-session model next to write_set is not automatically wrong —
# the model is still true — but it cannot be the REASON write_set is display-only, because that reason
# has to hold at any builder count and the model says nothing about builders.
stale_write_set_weak_premise_lines="$(
  grep -rIEn "$exclusive_session_premise_pattern" \
    "$core_root/actions" "$core_root/docs" "$board_root/tools/queue-kanban" 2>/dev/null \
    | grep -E 'write_set|overlaps' || true
)"
if [ -n "$stale_write_set_weak_premise_lines" ]; then
  printf 'FAIL: write_set is display-only at ANY builder count — no shipped file may justify that by naming the exclusive-session model either (REQ-079). The model is about queue OWNERS, not builders, so it cannot be the reason; point at actions/work-reference.md → Fan-Out Dispatch instead. Offending lines:\n' >&2
  printf '%s\n' "$stale_write_set_weak_premise_lines" >&2
  fail_count=$((fail_count + 1))
fi

# The wrapped-comment half of the same rule (REQ-075, weak form added by REQ-079). model.go and
# board.js explain write_set in comments and tooltip text that wrap across lines, so the line sweeps
# above cannot see them. Neither file has any business asserting how many REQs run at once, nor any
# business invoking the exclusive-session model at all — so file-level negatives are exact and stable.
assert_file_not_contains \
  "tools/queue-kanban/model.go" \
  "$builder_count_premise_pattern" \
  'tools/queue-kanban/model.go must not explain write_set with a one-REQ-at-a-time premise — it is advisory input to a human pick and the merge is the non-interference proof, at any builder count (REQ-075).'

assert_file_not_contains \
  "tools/queue-kanban/web/board.js" \
  "$builder_count_premise_pattern" \
  'tools/queue-kanban/web/board.js overlaps-badge tooltip must not explain write_set with a one-REQ-at-a-time premise — that reason is false since fan-out dispatch (REQ-075).'

assert_file_not_contains \
  "tools/queue-kanban/model.go" \
  "$exclusive_session_premise_pattern" \
  'tools/queue-kanban/model.go must not invoke the exclusive-session model — it is the weak fingerprint of the retired write_set premise, and the file has no other reason to name it (REQ-079).'

assert_file_not_contains \
  "tools/queue-kanban/web/board.js" \
  "$exclusive_session_premise_pattern" \
  'tools/queue-kanban/web/board.js must not invoke the exclusive-session model — it is the weak fingerprint of the retired write_set premise, and the file has no other reason to name it (REQ-079).'

# write_set parser lock-step (REQ-032, updated by REQ-069). write_set does not gate dispatch; it
# survives as a display-only field the board parser reads for the overlaps badge, so the parser must
# still read it. (The reason it never gates is Fan-Out Dispatch's merge-is-the-proof rule, not any
# builder count — see the REQ-075 assertion above.)
assert_contains \
  "tools/queue-kanban/model.go" \
  'fields\["write_set"\]' \
  'tools/queue-kanban/model.go must keep parsing write_set — the display-only field the board overlaps badge reads.'

# Exclusive-session model (REQ-069). The concurrency machinery — the orchestrator lock, the
# parallel-dispatch gate, the co-dispatch re-validations — was removed for a one-session /
# one-active-REQ / one-coder-context operating rule. None of the lock/claim vocabulary may
# survive in any shipped action file, and the two replacement rules must each be stated exactly
# once: a second copy is drift waiting to diverge, the very failure the removed machinery kept
# re-learning across four REQs.
for removed_concurrency_token in \
  'Concurrent-Orchestrator Lock Guard' \
  'coexisting_sessions' \
  'claimed_reqs' \
  'heartbeat_at' \
  'orchestrator-lock\.json'; do
  concurrency_token_hits="$(grep -rIlE -- "$removed_concurrency_token" "$core_root/actions" 2>/dev/null || true)"
  if [ -n "$concurrency_token_hits" ]; then
    printf 'FAIL: removed concurrency-machinery token "%s" still present in a shipped action file (exclusive-session model, REQ-069):\n%s\n' \
      "$removed_concurrency_token" "$concurrency_token_hits" >&2
    fail_count=$((fail_count + 1))
  fi
done

# Reservation removal. The reserve action allocated REQs to a DIFFERENT worktree/cloud session, which
# the then-current exclusive-session model declared outside the product contract. REQ-096 re-grained
# that model to claim-anywhere, so the IDEA is now in contract — but the machinery is not, and this
# ratchet does not loosen: what stays dead is the reserve VERB, the `reserved` STATUS, and the
# frontmatter fields, because a claim is marked with the advisory `assigned_to` field instead (no
# staleness clock, no router-budget cost). None of the retired vocabulary may survive in shipped prose
# or the board tool. Tokens are underscore/path-shaped on purpose so ordinary English "reserved for"
# never false-positives — and `assigned_to` deliberately matches none of them.
for removed_reservation_token in \
  'status: reserved' \
  'reserved_for' \
  'reserved_at' \
  'do-work reserve' \
  'actions/reserve\.md'; do
  reservation_token_hits="$(grep -rIlE -- "$removed_reservation_token" \
    "$core_root/actions" "$core_root/docs" "$board_root/tools/queue-kanban" \
    "$core_root/SKILL.md" "$core_root/next-steps.md" 2>/dev/null || true)"
  # Per-file exemption (same pattern as the maintainer-doc allowlist): the update
  # flow in actions/version.md must NAME the removed reserve files on its Step 5
  # deletion line — tar extraction never deletes what upstream removed, so without
  # that rm every pre-0.161.0 install keeps the orphaned files forever. Only lines
  # containing `rm -f` are exempt; any other mention in version.md still fails.
  exempt_update_flow_file="$core_root/actions/version.md"
  filtered_reservation_hits=""
  for reservation_hit_file in $reservation_token_hits; do
    if [ "$reservation_hit_file" = "$exempt_update_flow_file" ] && \
       [ -z "$(grep -E -- "$removed_reservation_token" "$reservation_hit_file" | grep -vE -- 'rm -f' || true)" ]; then
      continue
    fi
    filtered_reservation_hits="$filtered_reservation_hits$reservation_hit_file
"
  done
  reservation_token_hits="$(printf '%s' "$filtered_reservation_hits")"
  if [ -n "$reservation_token_hits" ]; then
    printf 'FAIL: removed reservation-workflow token "%s" still present in a shipped file — the reserve verb/status stay dead; the advisory assigned_to field is how a claim is marked (REQ-096/REQ-097):\n%s\n' \
      "$removed_reservation_token" "$reservation_token_hits" >&2
    fail_count=$((fail_count + 1))
  fi
done

assert_file_missing \
  "actions/reserve.md" \
  'the reserve action must stay removed — a claim is marked with the advisory assigned_to field, not a reserve verb or a reserved status (REQ-096 re-grained ownership to claim-anywhere without reviving either).'

assert_file_missing \
  "actions/prime-req-reservation.md" \
  'the reservation prime doc must stay removed along with the reserve action.'

assert_contains \
  "actions/work-reference.md" \
  '^## Execution Model — Claim Anywhere, One Releaser' \
  'actions/work-reference.md must define the Execution Model section that replaced the orchestrator-lock/parallel-dispatch machinery. Renamed from "Exclusive Session" by REQ-096: claiming is no longer exclusive to one checkout, only the release tail is, and a heading naming the retired boundary is exactly the stale-name drift this suite exists to catch.'

# The invariant is about OWNERSHIP, not build count (REQ-073), and REQ-096 moved which
# ownership it names. REQ-069's wording ("one active REQ, one coder context") conflated
# ownership with build count, so lifting the builder cap required rewording it to
# "one queue owner per checkout"; REQ-096 then re-grained the model to claim-anywhere,
# which makes the *claiming* half false and leaves the release tail as the only thing
# that must not run twice. So the invariant is now "one releaser per queue" — still
# exactly once, because a second copy is drift waiting to diverge.
# `|| true` is load-bearing under `set -euo pipefail`: grep exits 1 on no match, and
# with pipefail that aborts the whole suite silently — a missing invariant would read
# as a crash with no FAIL line rather than as the failure it is.
releaser_invariant_count="$( { grep -roh 'one releaser per queue' "$core_root/actions" || true; } | wc -l | tr -d ' ')"
if [ "$releaser_invariant_count" != "1" ]; then
  printf 'FAIL: the ownership invariant ("one releaser per queue") must be stated exactly once across actions/ (found %s) — every other mention is a pointer, not a restatement.\n' \
    "$releaser_invariant_count" >&2
  fail_count=$((fail_count + 1))
fi

# The superseded wording must be gone, not merely outnumbered — same ratchet the
# "one active REQ, one coder context" check below applies to its predecessor. It bounds
# claiming to a single checkout, which is exactly what claim-anywhere makes false.
superseded_owner_invariant_hits="$(grep -rIlE -- 'one queue owner per checkout' "$core_root/actions" "$core_root/docs" "$core_root/SKILL.md" 2>/dev/null || true)"
if [ -n "$superseded_owner_invariant_hits" ]; then
  printf 'FAIL: the superseded invariant wording "one queue owner per checkout" still appears (REQ-096 replaced it with the one-releaser-per-queue formulation — any checkout may claim):\n%s\n' \
    "$superseded_owner_invariant_hits" >&2
  fail_count=$((fail_count + 1))
fi

# The old wording must be gone, not merely outnumbered: it says one active REQ and
# one coder context, which is exactly what fan-out makes false.
retired_invariant_hits="$(grep -rIlE -- 'one active REQ, one coder context' "$core_root/actions" "$core_root/docs" "$core_root/SKILL.md" 2>/dev/null || true)"
if [ -n "$retired_invariant_hits" ]; then
  printf 'FAIL: the retired invariant wording "one active REQ, one coder context" still appears (REQ-073 replaced it with the one-queue-owner formulation):\n%s\n' \
    "$retired_invariant_hits" >&2
  fail_count=$((fail_count + 1))
fi

# Fan-out dispatch (REQ-073). The builder cap was two sentences in the Worktree
# Dispatch Mode opening; both must stay gone, and the section that replaced them
# must keep the three things that make N builders safe without coordination: the
# human picks, the merge is the proof, and a named set of steps never parallelises.
for retired_builder_cap_phrase in \
  'The single active builder' \
  'only one builder is ever in flight'; do
  builder_cap_hits="$(grep -rIlE -- "$retired_builder_cap_phrase" "$core_root/actions" 2>/dev/null || true)"
  if [ -n "$builder_cap_hits" ]; then
    printf 'FAIL: the retired builder-cap phrase "%s" still appears — REQ-073 raised worktree dispatch from one builder to N under a single queue owner:\n%s\n' \
      "$retired_builder_cap_phrase" "$builder_cap_hits" >&2
    fail_count=$((fail_count + 1))
  fi
done

assert_contains \
  "actions/work-reference.md" \
  '\*\*Fan-Out Dispatch' \
  'actions/work-reference.md must define Fan-Out Dispatch inside Worktree Dispatch Mode — several builders under one releaser, with no new coordination state.'

fan_out_block="$(sed -n '/\*\*Fan-Out Dispatch/,/^## Composed Exit Summary/p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$fan_out_block" \
  'Serial-only' \
  'Fan-Out Dispatch must name what never parallelises — queue transitions, REQ id allocation, and the version/changelog files.'

assert_block_contains \
  "$fan_out_block" \
  'CHANGELOG\.md' \
  'the Serial-only list must name CHANGELOG.md explicitly — one entry per REQ written by the owner, because unique version numbers do not make a shared prepend safe.'

assert_block_contains \
  "$fan_out_block" \
  'advisory input|never a gate' \
  'Fan-Out Dispatch must keep write_set advisory input to the humans pick and never a gate — nothing schedules on it under any builder count.'

assert_block_contains \
  "$fan_out_block" \
  'line proximity, not meaning' \
  'Fan-Out Dispatch must state the merge gates honest limit: git detects conflicts by line proximity, so two REQs appending to a shared registry merge cleanly and can still be jointly wrong.'

assert_block_contains \
  "$fan_out_block" \
  'survivable, not prevented' \
  'Fan-Out Dispatch must carry crew-members/background-agents.md own ceiling note — the run-directory pattern makes fan-out failures survivable, never prevented.'

assert_block_contains \
  "$fan_out_block" \
  'absolute main-tree path' \
  'Fan-Out Dispatch must state the brief-delivery trap: a repo-relative path resolves inside the worktree against its own stale copy of do-work/.'

three_attempt_count="$(grep -roh 'consecutive fix attempts' "$core_root/actions" | wc -l | tr -d ' ')"
if [ "$three_attempt_count" != "1" ]; then
  printf 'FAIL: the three-attempt stop condition ("consecutive fix attempts ... in its current context only") must be stated exactly once across actions/ (found %s).\n' \
    "$three_attempt_count" >&2
  fail_count=$((fail_count + 1))
fi

# Claim-respecting crash recovery (REQ-071). Recovery resets frontmatter and strips thirteen
# generated sections, and the pipeline does not commit until Step 9 — so an unconditional recovery
# destroys a finished Plan/Exploration/Scope that exists nowhere else. The premise that licensed
# running it on every working/ file is gone and must stay gone; what replaced it is a classification,
# a human-authorized takeover, and an ask-on-ambiguity default. Each assertion below pins one of those.
#
# The premise sweep (widened by REQ-077). Trigger condition, not a wording list: any sentence in
# either pipeline file telling the reader that every file in do-work/working/ belongs to the current
# session, or that the pipeline keeps no claim record at all, restates that premise — and since
# REQ-077 the pipeline does keep a claim record (CHECKPOINT.md's In Progress list, written at claim
# time by Step 2). The fingerprints below are the ones seen in the tree so far; they are illustrative,
# not exhaustive, so a newly discovered wording earns a generalized pattern rather than a fourth
# literal. Whole-file scope over both files is also deliberate: the two predecessor assertions were
# each narrower than the premise — one read only actions/work.md's Step 1 block while the live
# restatement sat in Step 2, the other was scoped to actions/work-reference.md by argument — so
# broadening the pattern alone would have left both of them green.
for premise_file in actions/work.md actions/work-reference.md; do
  for premise_fingerprint in \
    'no other live session whose in-flight claim a recovery could disturb' \
    "every .working/. file is this session's" \
    'no lock or claim record'
  do
    assert_file_not_contains \
      "$premise_file" \
      "$premise_fingerprint" \
      "$premise_file must not restate the retired premise that every do-work/working/ file is this session's to recover, nor claim the pipeline keeps no claim record at all: Crash Recovery classifies each working/ file against CHECKPOINT.md's In Progress list, which Step 2 writes at claim time (REQ-071, REQ-077). Matched fingerprint: $premise_fingerprint"
  done
done

crash_recovery_block="$(sed -n '/^## Crash Recovery (Step 1)/,/^## Worktree Dispatch Mode/p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$crash_recovery_block" \
  'foreign claim' \
  'actions/work-reference.md Crash Recovery must classify a claimed working/ REQ the checkpoint does not name as a foreign claim, not as this sessions own leftover (REQ-071).'

assert_block_contains \
  "$crash_recovery_block" \
  'absent checkpoint is ambiguous' \
  'actions/work-reference.md Crash Recovery must state that a missing CHECKPOINT.md is ambiguous, never permission to recover — a hard crash usually leaves no checkpoint at all (REQ-071).'

assert_block_contains \
  "$crash_recovery_block" \
  'never authorizes' \
  'actions/work-reference.md Crash Recovery must state that the staleness threshold gates only the offer and never authorizes a takeover by itself, or a later edit simplifies it into an automatic one (REQ-071).'

assert_block_contains \
  "$crash_recovery_block" \
  'unparseable, future-dated, or absent' \
  'actions/work-reference.md Crash Recovery must guard a bad claimed_at toward asking (immediately eligible), or a REQ carrying a corrupt stamp is protected from takeover forever (REQ-071).'

assert_block_contains \
  "$crash_recovery_block" \
  'no human to answer' \
  'actions/work-reference.md Crash Recovery must keep the unattended path non-blocking — a foreign claim is left alone, reported, and the run continues (REQ-071).'

assert_block_contains \
  "$crash_recovery_block" \
  'substep 1 removes it' \
  'actions/work-reference.md Crash Recovery must say claimed_at is read while classifying, before substep 1 discards it — the same ordering trap as the Scope/write_set decision (REQ-071).'

# A label-less in-progress entry is never auto-recovered (REQ-104). The retired bullet made a
# locally modified do-work/CHECKPOINT.md stand in for authorship. REQ-095's two-clone acceptance run
# falsified that premise: F-06 shows the checkpoint conflicting on *every* concurrent claim, even two
# that overlap in nothing, so a merge-resolving checkout holds a dirty checkpoint for reasons unrelated
# to who wrote which entry — and F-07 reproduced the consequence, a label-less *foreign* entry
# classified as an own crash and its live claim stripped. That strip is the 2026-07-01 incident,
# reached through the label-less door. Both pins are needed: the positive one keeps the report-only
# classification stated, the negative one keeps the "locally modified ⇒ mine" inference from being
# reintroduced as an optimization for pre-0.170.0 checkpoints.
assert_block_contains \
  "$crash_recovery_block" \
  'claim of unknown origin, always report-only' \
  'actions/work-reference.md Crash Recovery must classify a label-less In Progress entry as a claim of unknown origin that is always report-only — never recovered, whatever state the local checkpoint is in (REQ-104).'

assert_block_not_contains \
  "$crash_recovery_block" \
  'locally modified or otherwise uncommitted' \
  'actions/work-reference.md Crash Recovery must not treat a locally modified CHECKPOINT.md as evidence this checkout authored a label-less entry: under claim-anywhere every concurrent claim conflicts on that file, so a merge-resolving checkout is dirty for reasons unrelated to authorship and the inference strips a live foreign claim (REQ-104).'

# Recovery stamps the flip instant (REQ-074). The automatic reset had been silently out of compliance
# with status_changed_at's own trigger condition since the field was introduced, and nothing caught it
# for that entire span — the manual reset in actions/forensics.md stamped, this one did not. Since the
# same substep removes claimed_at, an unstamped recovery leaves no trace of the reset at all and the
# board dates a just-recovered REQ from created_at.
assert_block_contains \
  "$crash_recovery_block" \
  'stamp `status_changed_at' \
  'actions/work-reference.md Crash Recovery substep 1 must stamp status_changed_at on the reset — the field is written on any status flip with no dedicated *_at of its own, and this substep also removes claimed_at, so it is the only surviving trace of when recovery happened (REQ-074).'

# The hand-back file's one legal write location (REQ-082). Fan-Out Dispatch makes a per-builder
# REQ-NNN-handback.md mandatory and background-agents.md has the sub-agent write it itself, while
# Sole integrator says a builder never writes the main tree and do-work/ is main-tree-only — so the
# mandatory file had no legal home, and an agent hitting that had three moves, two of which corrupt
# the run (write the main tree anyway, or write the worktree's copy where it lands in the branch and
# the orchestrator reads nothing). The failure this guards is a later maintenance pass reading "the
# builder never writes the main tree" as absolute and deleting the carve-out as redundant, silently
# restoring the contradiction.
worktree_dispatch_block="$(sed -n '/^## Worktree Dispatch Mode (Step 1)/,/^## Composed Exit Summary/p' "$core_root/actions/work-reference.md")"

# actions/work.md is the executable, condensed hand-back path. It must retain the two
# pre-merge guards from the canonical reference: isolate owner bookkeeping from any
# unrelated staged work, and reject builder commits under do-work/ while the branch
# diff can still see them. A trailing pointer to the reference is too late because a
# reader following the numbered commands has already merged by then.
work_action_handback_block="$(sed -n '/^\*\*Hand-back merge/,/^### Step 6\.25:/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_action_handback_block" \
  'git diff --cached --name-only' \
  'actions/work.md hand-back step 0 must inspect the index before committing owner bookkeeping, or a plain commit can hide unrelated staged work below the evidence range.'

assert_block_contains \
  "$work_action_handback_block" \
  'git commit -- do-work/' \
  'actions/work.md hand-back step 0 must commit owner bookkeeping path-limited, never sweep the whole index into the bookkeeping commit.'

assert_block_contains \
  "$work_action_handback_block" \
  'clean index' \
  'actions/work.md hand-back step 0 must end with a clean index so unrelated staged work cannot leak into the merge commit.'

# Hand-back scratch is expected main-tree state, but it must never enter the
# bookkeeping commit. Preserve all three categories in both the canonical and
# condensed instructions: stage owner bookkeeping, tolerate the named scratch
# file without staging it, and stop on everything else.
for handback_staging_block in "$worktree_dispatch_block" "$work_action_handback_block"
do
  assert_block_contains \
    "$handback_staging_block" \
    'stage.*`manifest\.md`.*`REQ-NNN-brief\.md`' \
    'Hand-back step 0 must stage the owner-written manifest and per-builder briefs explicitly instead of admitting the whole run directory.'

  assert_block_contains \
    "$handback_staging_block" \
    'allow but never stage.*`REQ-NNN-handback\.md`' \
    'Hand-back step 0 must treat each expected REQ-NNN-handback.md as allowed scratch that remains unstaged.'

  assert_block_contains \
    "$handback_staging_block" \
    'stop and surface.*every other `do-work/` path' \
    'Hand-back step 0 must stop on do-work paths outside the explicit bookkeeping and scratch categories.'
done

assert_block_not_contains \
  "$worktree_dispatch_block" \
  'the run directory if fan-out created one' \
  'actions/work-reference.md must not restore the whole fan-out run directory to the stage set, because it contains never-staged hand-back scratch.'

assert_block_not_contains \
  "$work_action_handback_block" \
  "this run's directory" \
  'actions/work.md must not restore the whole run directory to the stage set, because it contains never-staged hand-back scratch.'

assert_block_contains \
  "$work_action_handback_block" \
  'git diff --name-only <pre>\.\.\.<operative_name> -- do-work/' \
  'actions/work.md must reject builder queue-state commits before merging, while the three-dot branch diff can still see them.'

assert_block_contains \
  "$worktree_dispatch_block" \
  'exactly one exception' \
  'actions/work-reference.md Sole integrator must carry the bounded hand-back exception — Fan-Out Dispatch mandates REQ-NNN-handback.md, and without a named legal write location the builder must either violate sole-integrator or drop the file the durability pattern exists for (REQ-082).'

assert_block_contains \
  "$worktree_dispatch_block" \
  'never staged, committed, or merged' \
  'actions/work-reference.md Sole integrator must say the hand-back file is never staged/committed/merged — "you may write it" naturally reads as including committing it, which would turn run scratch into branch content and trip the builder-wrote-do-work probe (REQ-082).'

# Scope, not a list: the exception is one path derived from the builder's own REQ id, and do-work/'s
# main-tree-only rule must be stated as a condition. Both were closed enumerations away from silently
# growing — "the run directory" instead of one file, and a three-item list that predated do-work/runs/.
assert_block_contains \
  "$worktree_dispatch_block" \
  'Every path under `do-work/` exists in the main tree only' \
  'actions/work-reference.md State stays home must state the condition rather than enumerate the queue/working//CHECKPOINT.md trio — do-work/runs/ postdated that list and a reader could argue it was out of scope (REQ-082).'

assert_block_not_contains \
  "$worktree_dispatch_block" \
  'the queue, `working/`, `CHECKPOINT.md` — exists in the main tree only' \
  'actions/work-reference.md must not restore State stays home three-item enumeration — a hand-maintained list of what lives under do-work/ goes stale the moment a directory is added (REQ-082).'

# The claim-time write (REQ-077) is the other half of REQ-071's gate. REQ-071 made recovery consume
# the checkpoint's In Progress record, but Step 10 (session end) was its only write site — so a hard
# crash left no record at all, every crashed REQ classified as a foreign claim, and the own-crash
# branch became unreachable by the exact event it exists to handle. Each assertion below pins one
# numbered requirement of the fix: write it at claim time, keep it a list, keep it from becoming a
# lock, and remove it when the REQ leaves working/.
work_step_two_block="$(sed -n '/^### Step 2: Claim the Request/,/^### Step 3: Triage/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_step_two_block" \
  'In Progress \(interrupted\)' \
  'actions/work.md Step 2 must record the claim in CHECKPOINT.md In Progress (interrupted) at claim time — with Step 10 (session end) as the only write site, a hard crash leaves recovery no classification input and the REQ strands in working/ forever (REQ-077).'

in_progress_record_block="$(sed -n '/^## In-Progress Record (Step 2)/,/^## Triage Section Template/p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$in_progress_record_block" \
  'record is a list' \
  'actions/work-reference.md In-Progress Record must specify the record as a list, one entry per claimed REQ — fan-out claims several REQs under one owner, so a singular record classifies every claim but the newest as a foreign claim after a crash (REQ-077).'

assert_block_contains \
  "$in_progress_record_block" \
  'never grow into one' \
  'actions/work-reference.md In-Progress Record must state that the record is a classification input and must never grow into a lock, heartbeat, or liveness check — that is the machinery REQ-069 deleted and REQ-073 declined to revive (REQ-077).'

# The writer label (REQ-094). The checkpoint is a tracked file wherever a consumer commits do-work/,
# so another checkout's live claim arrives by an ordinary git pull looking exactly like a local one —
# recovery then reads it as an own crash and strips a REQ someone is actively building. Each
# assertion below pins one half of the fix: the entry carries the label, and a labeled foreign entry
# is reported rather than aged into the takeover ladder.
assert_block_contains \
  "$in_progress_record_block" \
  'writer: <hostname>:<absolute-checkout-path>' \
  'actions/work-reference.md In-Progress Record must give the entry format a writer: <hostname>:<absolute-checkout-path> label — the path alone collides across machines and the hostname alone collides across checkouts on one machine, so neither half identifies a checkout by itself (REQ-094).'

assert_block_contains \
  "$crash_recovery_block" \
  'claim held by' \
  'actions/work-reference.md Crash Recovery must report a foreign-label entry as a claim held by that writer and leave it untouched — a label is positive evidence of another checkouts claim, so it is classified by the label and never aged into the three-hour takeover ladder (REQ-094).'

# The tripwire had to be reworded to admit a static writer label without admitting liveness
# machinery, and the naive reword drops the ban to make room for the carve-out. Pinning both phrases
# to the SAME paragraph is the point: a carve-out that drifts into its own paragraph reads as a
# general permission rather than as the one exception to the ban standing beside it.
in_progress_tripwire_paragraph="$(printf '%s\n' "$in_progress_record_block" | grep -F 'never grow into one' || true)"
if [ -z "$in_progress_tripwire_paragraph" ] || ! printf '%s\n' "$in_progress_tripwire_paragraph" | grep -qF 'never refreshed'; then
  printf 'FAIL: actions/work-reference.md In-Progress Record must keep the tripwire ("never grow into one") and the static-label carve-out ("never refreshed") in one paragraph — the label is the single exception to the lock/heartbeat/liveness ban, and it only reads as an exception standing next to it (REQ-077, REQ-094).\n' >&2
  fail_count=$((fail_count + 1))
fi

# The two label-destruction paths in actions/work.md Step 10 (REQ-102). Both were scoped to entries
# "carrying another checkout's writer: label", which silently excludes the label-less legacy entry
# that Crash Recovery classifies report-only in a clean, committed checkpoint — the wholesale rewrite
# could drop it and the session-start delete could remove the whole file, after which the next run
# sees a working/ REQ "not named there" and ages it into the three-hour takeover ladder. Nothing
# pinned either clause, so a later "simplify the checkpoint rewrite" pass could reopen the hole with
# the suite green. Both must scope preservation to every entry this checkout did not write.
work_session_checkpoint_block="$(sed -n '/^#### Session Checkpoint/,/^## Clarify Questions/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_session_checkpoint_block" \
  'entry this checkout did not write through verbatim' \
  'actions/work.md Step 10 Session Checkpoint must carry every In Progress entry this checkout did not write through verbatim — scoping the preserve rule to labeled foreign entries lets the wholesale rewrite drop a label-less entry that Crash Recovery had classified report-only (REQ-102).'

assert_block_contains \
  "$work_session_checkpoint_block" \
  'no entry this checkout did not write remains' \
  'actions/work.md session-start step 3 must gate deleting CHECKPOINT.md on no entry this checkout did not write remaining — gating on labeled foreign entries alone authorizes deleting a checkpoint whose only surviving claim is label-less, which the next run then ages into the takeover ladder (REQ-102).'

assert_block_contains \
  "$work_archive_success_block" \
  'In Progress \(interrupted\)' \
  'actions/work.md Step 8 must remove the REQ In Progress entry as part of the archive move — a REQ still listed there after it leaves working/ is the contradiction the next run is told to report, and a report that fires on every normal completion trains readers to ignore it (REQ-077).'

work_step_one_block="$(sed -n '/^### Step 1: Find Next Request/,/^### Step 2\.0/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_step_one_block" \
  "Crash Recovery's input" \
  'actions/work.md Step 1 must name do-work/CHECKPOINT.md as Crash Recoverys input, so the read is understood as a precondition rather than resume convenience (REQ-071).'

# Ordering, not just presence: the checkpoint read has to appear ahead of the Crash Recovery
# paragraph in Step 1. Recovery consumes the checkpoint, so a later edit that moves the read below
# it would leave recovery classifying against a file it has not opened.
checkpoint_read_line="$(printf '%s\n' "$work_step_one_block" | grep -n 'CHECKPOINT\.md' | head -1 | cut -d: -f1)"
crash_recovery_line="$(printf '%s\n' "$work_step_one_block" | grep -n '\*\*Crash Recovery:\*\*' | head -1 | cut -d: -f1)"
if [ -z "$checkpoint_read_line" ] || [ -z "$crash_recovery_line" ] || [ "$checkpoint_read_line" -ge "$crash_recovery_line" ]; then
  printf 'FAIL: actions/work.md Step 1 must read do-work/CHECKPOINT.md BEFORE the Crash Recovery paragraph (checkpoint line: %s, recovery line: %s) — the checkpoint is recovery input, not resume decoration (REQ-071).\n' \
    "${checkpoint_read_line:-none}" "${crash_recovery_line:-none}" >&2
  fail_count=$((fail_count + 1))
fi

# Allocator / verifier subcommands (REQ-072). The tool grew a second write surface
# (one version line, outside do-work/) and three release-ritual subcommands. Two
# things must not drift: the never-write-CHANGELOG boundary, and the three prose
# call sites — a subcommand nothing calls is dead weight, and a call site with no
# stated fallback turns the Go toolchain into a hard dependency of the pipeline.
assert_file_not_contains \
  "tools/queue-kanban/release.go" \
  'os\.(WriteFile|Create|OpenFile)\(.*CHANGELOG' \
  'tools/queue-kanban/release.go must never write CHANGELOG.md — unique version numbers do not make a shared prepend safe, so the changelog stays an owner-only human write.'

assert_file_not_contains \
  "tools/queue-kanban/verify.go" \
  'os\.(WriteFile|Create|OpenFile|Remove|Rename)\(' \
  'tools/queue-kanban/verify.go must stay read-only — verify reports and routes, and repairs belong to actions/cleanup.md, which asks first.'

# The unknown-subcommand message must list every subcommand main() dispatches, or the
# error text lies about what exists. Derived from the dispatch switch rather than pinned
# to a hand-written chain: a frozen list only ever verifies the subcommands that already
# existed when it was written, which is exactly how it would miss the next one added
# (Closed Enumerations Go Stale).
unknown_subcommand_message="$(grep -F 'unknown subcommand %q' "$board_root/tools/queue-kanban/main.go" || true)"
if [ -z "$unknown_subcommand_message" ]; then
  printf 'FAIL: tools/queue-kanban/main.go has no unknown-subcommand message — a mistyped subcommand must name the valid ones, not exit silently.\n' >&2
  fail_count=$((fail_count + 1))
else
  dispatched_subcommands="$(awk '
    /^\tswitch subcommand \{$/ { in_dispatch_switch=1; next }
    in_dispatch_switch && /^\t}$/ { exit }
    in_dispatch_switch && /^\tcase / {
      case_labels=$0
      sub(/^[[:space:]]*case[[:space:]]+/, "", case_labels)
      sub(/:.*/, "", case_labels)
      while (match(case_labels, /"[a-z-]+"/)) {
        print substr(case_labels, RSTART + 1, RLENGTH - 2)
        case_labels=substr(case_labels, RSTART + RLENGTH)
      }
    }
  ' "$board_root/tools/queue-kanban/main.go" | sort -u)"

  advertised_subcommands="$(printf '%s\n' "$unknown_subcommand_message" \
    | sed -n 's/.*(want \(.*\))\\n".*/\1/p' \
    | tr '|' '\n' \
    | sed 's/^[[:space:]]*//; s/[[:space:]]*$//' \
    | sort -u)"
  if [ -z "$dispatched_subcommands" ]; then
    printf 'FAIL: could not read the subcommand dispatch switch in tools/queue-kanban/main.go — this check derives the expected list from it, so a shape change here silently disables it.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  for dispatched_subcommand in $dispatched_subcommands; do
    if ! grep -qxF -- "$dispatched_subcommand" <<<"$advertised_subcommands"; then
      printf 'FAIL: tools/queue-kanban/main.go dispatches %s but its unknown-subcommand message does not list it — the error text lies about what exists.\n' \
        "$dispatched_subcommand" >&2
      fail_count=$((fail_count + 1))
    fi
  done
fi

# Project release instructions must never point next-version at the managed core
# package. That action file is the suite runtime's version source, not the consumer
# application's selected release source; bumping it poisons later update comparisons.
assert_file_not_contains \
  "actions/work.md" \
  'queue-kanban next-version|queue-kanban/queue-kanban next-version|--version-file.*<skill-root>/actions/version\.md' \
  'actions/work.md must leave project version bumps to the Changelog Entry Procedure and never mutate the installed suite version.'

assert_contains \
  "actions/work-reference.md" \
  'Exclude every installed skill, dependency, vendored package, and generated tree.*\.claude/skills/.*installed do-work suite.*VERSION.*actions/version\.md.*must never be selected or bumped' \
  'The canonical project version-source procedure must exclude installed suite metadata, not only remove one unsafe accelerator call site.'

assert_contains \
  "tools/queue-kanban/main.go" \
  'func rejectLeftoverArguments' \
  'tools/queue-kanban/main.go must keep the shared leftover-argument rejection — an unconsumed token is an error for ANY subcommand, not silence, which is how next-version writing the wrong tree stayed invisible (REQ-081).'

for release_subcommand_call_site in \
  'actions/capture.md:queue-kanban next-req' \
  'actions/forensics.md:queue-kanban verify'; do
  call_site_file="${release_subcommand_call_site%%:*}"
  call_site_pattern="${release_subcommand_call_site#*:}"
  assert_contains \
    "$call_site_file" \
    "$call_site_pattern" \
    "$call_site_file must call \`$call_site_pattern\` — REQ-072 wired the allocator/verifier into the three existing actions instead of adding a new action or a SKILL.md routing row."
  # Every call site must name what to do when the toolchain is missing. The board
  # action reports and stops because there the compiler IS the capability; these
  # three must fall back to the procedure they accelerate.
  assert_contains \
    "$call_site_file" \
    'If .go. is absent|when .go. is absent|absent or the build fails' \
    "$call_site_file must state the fallback for a missing Go toolchain next to its queue-kanban call — the compiler is an accelerator here, never a dependency of the pipeline."
done

assert_contains \
  "CLAUDE.md" \
  'three write surfaces' \
  'CLAUDE.md must state the tool has exactly three write surfaces once next-req reserves ids — testing fields, next-version, and reservation markers are the complete set, and nothing but this sentence records the count.'

assert_contains \
  "actions/capture.md" \
  '\.req-reservations/REQ-NNNNNN' \
  'actions/capture.md must stage next-req reservation markers with the UR/REQ capture — an uncommitted marker would reserve only one checkout and pollute git status.'

assert_contains \
  "docs/forensics-guide.md" \
  'Release and queue invariants' \
  'docs/forensics-guide.md must list the release/queue invariants check — a forensics check absent from the user-facing guide is invisible to the person who would run it.'

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
router_word_count="$(wc -w < "$core_root/SKILL.md")"
if [ "$router_word_count" -gt "$router_word_budget" ]; then
  printf 'FAIL: SKILL.md is %s words — over the %s-word router budget. Merge or lazy-load; do not grow the always-loaded router.\n' \
    "$router_word_count" "$router_word_budget" >&2
  fail_count=$((fail_count + 1))
fi

# Hardened checks (REQ-018): the work.md prose pointers and the shipped scripts
# must not drift apart — a pointer at a missing script silently un-hardens the step.
# Each entry is "<script>|<action file that must reference it>". The referencing file is
# per-script because not every hardened check belongs to the work pipeline —
# blanked-req-scan.sh is called from forensics (and, for its restore mode, cleanup). Pinning
# the right caller keeps the assertion strong; hardcoding actions/work.md for all of them
# would have forced either a bogus reference or dropping the check entirely.
hardened_check_scripts=(
  "tools/checks/archive-collision.sh|actions/work.md"
  "tools/checks/preflight.sh|actions/work.md"
  "tools/checks/scope-drift.sh|actions/work.md"
  "tools/checks/qualify.sh|actions/work.md"
  "tools/checks/record-commit-hash.sh|actions/work.md"
  "tools/checks/blanked-req-scan.sh|actions/forensics.md"
  "tools/checks/blanked-req-scan.sh|actions/cleanup.md"
  "tools/checks/uncommitted-inventory.sh|actions/commit.md"
  "tools/checks/uncommitted-inventory.sh|actions/inspect.md"
  "tools/checks/associate-files.sh|actions/commit.md"
  "tools/checks/associate-files.sh|actions/inspect.md"
)

for check_script_entry in "${hardened_check_scripts[@]}"; do
  check_script="${check_script_entry%%|*}"
  referencing_action_file="${check_script_entry##*|}"
  if [ ! -x "$(resolve_runtime_file "$check_script")" ]; then
    printf 'FAIL: %s must exist and be executable (%s points at it).\n' "$check_script" "$referencing_action_file" >&2
    fail_count=$((fail_count + 1))
  fi
  assert_contains \
    "$referencing_action_file" \
    "$(basename "$check_script")" \
    "$referencing_action_file must reference $check_script — the hardened step's pointer was removed without un-hardening."
done

# Pre-flight dirty-tree relevance (feedback 2026-08-04). The serial qualifier and
# reviewer inspect the repository-wide working/staged diff, so changes that predate
# the active REQ can contaminate its evidence even though Step 9 stages explicit
# paths. Keep the broad detection, but keep its warning subordinate to the canonical
# Current-REQ relevance rule: preserve/exclude unexpected state and continue unless
# it prevents this REQ from completing. The old "may stage unrelated files" rationale
# contradicted the explicit staging contract and invited removal of a useful check.
current_req_relevance_block="$(sed -n '/^\*\*Current-REQ relevance\./,/^\*\*Three-attempt stop\./p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$current_req_relevance_block" \
  'only.*prevents the active REQ.*implemented, tested, archived, or committed' \
  'actions/work-reference.md must keep unexpected repository state gated by Current-REQ relevance — otherwise pre-flight warnings can become blockers or cleanup work.'

assert_block_contains \
  "$current_req_relevance_block" \
  'preserve it, exclude it from this REQ.s staging, and continue' \
  'actions/work-reference.md must keep the preserve/exclude/continue response for repository state that does not prevent the active REQ.'

# The rule covers SESSION state too, and says the non-action out loud (feedback
# 2026-08-04). A session enumerated its sibling claude PIDs, noticed a commit that
# landed 20 seconds into its run, and blocked on a four-option "how should I proceed?"
# prompt instead of starting its REQ. Nothing shipped told it to — but the contract
# only said the *pipeline* does not detect a concurrent run, which an agent can read
# as leaving the question open for it to raise. Exclusivity is the user's guarantee;
# asking them to re-confirm it rebuilds the coordination the exclusive-session model
# deleted, one prompt at a time.
assert_block_contains \
  "$current_req_relevance_block" \
  '[Nn]ever probe for a concurrent session' \
  'actions/work-reference.md must keep Current-REQ relevance covering session state with the explicit never-probe/never-ask clause — without it an agent re-derives the concurrency check the exclusive-session model removed and stalls the loop on a prompt.'

# The secret-shaped exclusion is a behavior probe, not a grep, because the bug it
# guards was a glob that LOOKED right. `.env|.env.*` reads as covering ".env*" but
# matches neither `.envrc` (direnv — routinely full of exported secrets) nor
# `.environment`: both are suffixes with no dot, so each fell through to `A` and the
# callers would read and stage them. Codex caught it on PR #134. Assert the tags,
# not the pattern — the next wrong pattern will look right too.
inventory_probe_dir="$(mktemp -d)"
cleanup_inventory_probe() {
  rm -rf -- "$inventory_probe_dir"
}
trap cleanup_inventory_probe EXIT
inventory_probe_setup_error="$inventory_probe_dir/setup-error.txt"
if ! (
  export GIT_CONFIG_GLOBAL=/dev/null
  export GIT_CONFIG_SYSTEM=/dev/null
  cd "$inventory_probe_dir" || exit 1
  git init -q .
  git config user.email probe@test
  git config user.name probe
  mkdir -p nested uppercase
  for probe_name in .env .env.local .envrc .environment production.env \
                    credentials.json server.pem .ENV.PRODUCTION AuthCredentials.json \
                    private.PEM UPPER-SECRET.txt ordinary.js; do
    echo probe > "nested/$probe_name"
  done
  echo probe > uppercase/.ENV
  git add nested/.env.local
  git commit -qm 'seed tracked secret deletion'
  rm nested/.env.local
  echo probe > .env
  git add .env
  git commit -qm 'seed tracked secret rename'
  git mv .env visible-config.txt
) 2>"$inventory_probe_setup_error"; then
  printf 'FAIL: could not set up the uncommitted-inventory behavior probe:\n' >&2
  sed 's/^/  /' "$inventory_probe_setup_error" >&2
  fail_count=$((fail_count + 1))
else
  inventory_probe_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$inventory_probe_dir" 2>/dev/null || true)"
  for must_be_excluded in .env .envrc .environment production.env credentials.json server.pem \
                          .ENV.PRODUCTION AuthCredentials.json private.PEM UPPER-SECRET.txt; do
    if ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'X\tnested/%s' "$must_be_excluded")"; then
      printf 'FAIL: tools/checks/uncommitted-inventory.sh must tag nested/%s as X (secret-shaped) — it is reachable by the advertised exclusion globs.\n' "$must_be_excluded" >&2
      fail_count=$((fail_count + 1))
    fi
  done
  if ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'X\tuppercase/.ENV')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must tag uppercase/.ENV as X (case-insensitive secret-shaped basename).\n' >&2
    fail_count=$((fail_count + 1))
  fi
  if ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'XD\tnested/.env.local')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must tag a deleted secret-shaped path as XD so its deletion can be associated and committed without reading its former contents.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  if ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'XD\t.env')" || ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'X\tvisible-config.txt')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must fail closed when a secret-shaped rename origin moves to an ordinary-looking destination: XD for the source and X for the destination.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  if ! printf '%s\n' "$inventory_probe_output" | grep -qxF "$(printf 'X\tnested/ordinary.js')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must quarantine every A as X when an excluded path makes addition provenance ambiguous.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  # A caller can disable rename detection in repository config. The shipped
  # inventory must override that setting explicitly; otherwise the same staged
  # rename above degrades to an XD + A pair before the action has any provenance
  # to retain.
  git -C "$inventory_probe_dir" config status.renames false
  inventory_renames_disabled_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$inventory_probe_dir" 2>/dev/null || true)"
  if ! printf '%s\n' "$inventory_renames_disabled_output" | grep -qxF "$(printf 'XD\t.env')" || \
      ! printf '%s\n' "$inventory_renames_disabled_output" | grep -qxF "$(printf 'X\tvisible-config.txt')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must force rename detection even when status.renames=false: XD for .env and X for visible-config.txt.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  # commit Step 1 resets a pre-staged X destination. That erases rename
  # provenance from porcelain status and leaves an already-staged source
  # deletion plus an untracked destination. The destination must remain X.
  if ! git -C "$inventory_probe_dir" reset -q -- visible-config.txt; then
    printf 'FAIL: could not reset the secret-rename destination for the re-inventory probe.\n' >&2
    fail_count=$((fail_count + 1))
  else
    inventory_after_reset_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
      "$core_root/tools/checks/uncommitted-inventory.sh" "$inventory_probe_dir" 2>/dev/null || true)"
    if ! printf '%s\n' "$inventory_after_reset_output" | grep -qxF "$(printf 'XD\t.env')" || \
        ! printf '%s\n' "$inventory_after_reset_output" | grep -qxF "$(printf 'X\tvisible-config.txt')"; then
      printf 'FAIL: reset-and-reinventory must fail closed: XD for .env and X, never A, for visible-config.txt.\n' >&2
      fail_count=$((fail_count + 1))
    fi

    staged_deletion_metadata="$inventory_probe_dir/staged-deletion-metadata.bin"
    git -C "$inventory_probe_dir" diff --cached --name-status --no-renames -z -- .env > "$staged_deletion_metadata"
    staged_deletion_status=''
    staged_deletion_path=''
    staged_deletion_extra=''
    {
      IFS= read -r -d '' staged_deletion_status
      IFS= read -r -d '' staged_deletion_path
      IFS= read -r -d '' staged_deletion_extra || true
    } < "$staged_deletion_metadata"
    if [ "$staged_deletion_status" != 'D' ] || [ "$staged_deletion_path" != '.env' ] || [ -n "$staged_deletion_extra" ]; then
      printf 'FAIL: reset probe must leave one exact cached deletion for .env; got status=%s path=%s extra=%s.\n' \
        "$staged_deletion_status" "$staged_deletion_path" "$staged_deletion_extra" >&2
      fail_count=$((fail_count + 1))
    fi
    if git -C "$inventory_probe_dir" add -u -- .env 2>/dev/null; then
      printf 'FAIL: probe no longer reproduces Git rejecting git add -u for an already-staged rename-source deletion.\n' >&2
      fail_count=$((fail_count + 1))
    fi
  fi
fi
cleanup_inventory_probe
trap - EXIT

# An ordinary addition remains readable when no secret-shaped deletion makes
# its provenance ambiguous. Keep this in a separate repository: combining it
# with the XD fixture above would assert the unsafe behavior REQ-128 removes.
ordinary_addition_probe_dir="$(mktemp -d)"
if ! (
  export GIT_CONFIG_GLOBAL=/dev/null
  export GIT_CONFIG_SYSTEM=/dev/null
  cd "$ordinary_addition_probe_dir" || exit 1
  git init -q .
  echo ordinary > ordinary.js
); then
  printf 'FAIL: could not set up the ordinary-addition inventory probe.\n' >&2
  fail_count=$((fail_count + 1))
else
  ordinary_addition_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$ordinary_addition_probe_dir" 2>/dev/null || true)"
  if ! printf '%s\n' "$ordinary_addition_output" | grep -qxF "$(printf 'A\tordinary.js')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must leave an ordinary addition as A when no XD exists.\n' >&2
    fail_count=$((fail_count + 1))
  fi
fi
rm -rf -- "$ordinary_addition_probe_dir"

# REQ-148: awk's common `NR == FNR` first-file discriminator fails when the
# quarantine file is empty: every record from the inventory then also satisfies
# NR == FNR and is swallowed as if it belonged to the exclusion table. Pin the
# portable filename discriminator in every active bridge/modular copy, then
# exercise both an empty quarantine and a populated once-X-always-X set.
association_candidate_action_files=(
  actions/commit.md
  actions/inspect.md
  skills/do-work/actions/commit.md
  skills/do-work-toolbox/actions/inspect.md
)
for association_candidate_action_file in "${association_candidate_action_files[@]}"; do
  assert_contains \
    "$association_candidate_action_file" \
    'FILENAME == ARGV\[1\] \{ excluded\[\$0\] = 1; next \}' \
    "$association_candidate_action_file must distinguish the quarantine by filename so an empty first file cannot swallow every inventory candidate."
  assert_file_not_contains \
    "$association_candidate_action_file" \
    'NR == FNR \{ excluded\[\$0\] = 1; next \}' \
    "$association_candidate_action_file must not use NR == FNR for the possibly-empty quarantine merge."
done

association_candidate_probe_dir="$(mktemp -d)"
association_empty_quarantine="$association_candidate_probe_dir/empty-quarantine.txt"
association_populated_quarantine="$association_candidate_probe_dir/populated-quarantine.txt"
association_inventory="$association_candidate_probe_dir/inventory.txt"
association_expected_empty="$association_candidate_probe_dir/expected-empty.txt"
association_expected_populated="$association_candidate_probe_dir/expected-populated.txt"
association_actual_output="$association_candidate_probe_dir/actual.txt"
: > "$association_empty_quarantine"
printf 'M\tsrc/modified.js\nA\tsrc/added.js\nD\tsrc/deleted.js\nXD\t.env.deleted\nX\tcurrent-secret.txt\n' > "$association_inventory"
printf '%s\n' \
  'src/modified.js' \
  'src/added.js' \
  'src/deleted.js' \
  '.env.deleted' > "$association_expected_empty"

filter_association_candidates() {
  awk -F '\t' '
    FILENAME == ARGV[1] { excluded[$0] = 1; next }
    {
      tag = $1
      sub(/^[^\t]*\t/, "")
      if (tag != "X" && !($0 in excluded)) print
    }
  ' "$1" "$2"
}

filter_association_candidates "$association_empty_quarantine" "$association_inventory" > "$association_actual_output"
if ! cmp -s "$association_expected_empty" "$association_actual_output"; then
  printf 'FAIL: an empty secret quarantine must preserve every safe M/A/D/XD association candidate; expected/actual diff:\n' >&2
  diff -u "$association_expected_empty" "$association_actual_output" >&2 || true
  fail_count=$((fail_count + 1))
fi

printf '%s\n' 'src/added.js' 'previous-secret.txt' > "$association_populated_quarantine"
awk -F '\t' '$1 == "X" { sub(/^[^\t]*\t/, ""); print }' "$association_inventory" >> "$association_populated_quarantine"
printf '%s\n' \
  'src/modified.js' \
  'src/deleted.js' \
  '.env.deleted' > "$association_expected_populated"
filter_association_candidates "$association_populated_quarantine" "$association_inventory" > "$association_actual_output"
if ! cmp -s "$association_expected_populated" "$association_actual_output"; then
  printf 'FAIL: a populated secret quarantine must exclude only retained/current X paths and preserve every other safe candidate; expected/actual diff:\n' >&2
  diff -u "$association_expected_populated" "$association_actual_output" >&2 || true
  fail_count=$((fail_count + 1))
fi
rm -rf -- "$association_candidate_probe_dir"

# Git has no tracked source blob to compare when both a secret-shaped source and
# its copied destination are untracked. The inventory must therefore quarantine
# the ordinary-looking destination too, rather than trusting copy detection
# that cannot exist for this shape.
untracked_secret_copy_probe_dir="$(mktemp -d)"
if ! (
  export GIT_CONFIG_GLOBAL=/dev/null
  export GIT_CONFIG_SYSTEM=/dev/null
  cd "$untracked_secret_copy_probe_dir" || exit 1
  git init -q .
  printf 'fixture-secret\n' > .env
  cp .env application-config.txt
); then
  printf 'FAIL: could not set up the untracked secret-copy inventory probe.\n' >&2
  fail_count=$((fail_count + 1))
else
  untracked_secret_copy_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$untracked_secret_copy_probe_dir" 2>/dev/null || true)"
  if ! printf '%s\n' "$untracked_secret_copy_output" | grep -qxF "$(printf 'X\t.env')" || \
      ! printf '%s\n' "$untracked_secret_copy_output" | grep -qxF "$(printf 'X\tapplication-config.txt')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must quarantine an ordinary-looking untracked copy beside an untracked secret-shaped source.\n' >&2
    fail_count=$((fail_count + 1))
  fi
fi
rm -rf -- "$untracked_secret_copy_probe_dir"

# Copy detection must be requested explicitly rather than inherited from
# repository configuration. Without it, a secret-derived copy is reported as an
# ordinary A and both action callers are allowed to read it. Keep an ordinary
# rename in the same repository to prove copy-aware detection does not change
# the established M classification for non-secret moves.
copy_inventory_probe_dir="$(mktemp -d)"
if ! (
  export GIT_CONFIG_GLOBAL=/dev/null
  export GIT_CONFIG_SYSTEM=/dev/null
  cd "$copy_inventory_probe_dir" || exit 1
  git init -q .
  git config user.email probe@test
  git config user.name probe
  copy_line=1
  while [ "$copy_line" -le 100 ]; do
    printf 'secret-source-%03d\n' "$copy_line"
    copy_line=$((copy_line + 1))
  done > .env.copy-source
  echo ordinary-source > ordinary-source.txt
  git add .env.copy-source ordinary-source.txt
  git commit -qm 'seed copy and rename sources'
  cp .env.copy-source copied-config.txt
  echo changed >> .env.copy-source
  git add .env.copy-source copied-config.txt
  git mv ordinary-source.txt ordinary-destination.txt
  git config status.renames false
); then
  printf 'FAIL: could not set up the copy-aware inventory probe.\n' >&2
  fail_count=$((fail_count + 1))
else
  copy_inventory_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$copy_inventory_probe_dir" 2>/dev/null || true)"
  if ! printf '%s\n' "$copy_inventory_output" | grep -qxF "$(printf 'X\tcopied-config.txt')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must tag a secret-derived copy destination as X, never A, even when status.renames=false.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  if ! printf '%s\n' "$copy_inventory_output" | grep -qxF "$(printf 'M\tordinary-destination.txt')"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must retain M for an ordinary rename while copy-aware detection is active.\n' >&2
    fail_count=$((fail_count + 1))
  fi
fi
rm -rf -- "$copy_inventory_probe_dir"

# The other Step 5 branch remains live: a tracked secret deletion that is not
# cached yet still needs git add -u, followed by deletion-only metadata.
unstaged_deletion_probe_dir="$(mktemp -d)"
if ! (
  export GIT_CONFIG_GLOBAL=/dev/null
  export GIT_CONFIG_SYSTEM=/dev/null
  cd "$unstaged_deletion_probe_dir" || exit 1
  git init -q .
  git config user.email probe@test
  git config user.name probe
  echo probe > .env
  git add .env
  git commit -qm 'seed unstaged secret deletion'
  rm .env
); then
  printf 'FAIL: could not set up the unstaged secret-deletion probe.\n' >&2
  fail_count=$((fail_count + 1))
else
  if [ -n "$(git -C "$unstaged_deletion_probe_dir" diff --cached --name-status --no-renames -- .env)" ]; then
    printf 'FAIL: unstaged secret-deletion probe unexpectedly began with cached metadata.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! git -C "$unstaged_deletion_probe_dir" add -u -- .env || \
       ! git -C "$unstaged_deletion_probe_dir" diff --cached --name-status --no-renames -- .env \
         | grep -qxF "$(printf 'D\t.env')"; then
    printf 'FAIL: an unstaged tracked secret deletion must still stage as one exact cached D entry.\n' >&2
    fail_count=$((fail_count + 1))
  fi
fi
rm -rf -- "$unstaged_deletion_probe_dir"

# The action prose is executable behavior. Both consumers must retain a
# run-level quarantine before their second inventory reaches content or REQ
# association, and the manual fallback must force copy-aware detection too.
for inventory_consumer in actions/commit.md actions/inspect.md; do
  assert_contains \
    "$inventory_consumer" \
    'once X.*always X|once `X`.*always `X`' \
    "$inventory_consumer must preserve a once-X-always-X quarantine across every re-inventory in the action run."
  assert_contains \
    "$inventory_consumer" \
    'git -c status.renames=copies status --porcelain=v1 --untracked-files=all -z' \
    "$inventory_consumer manual fallback must force copy-aware detection instead of inheriting status.renames=false."
done

assert_contains \
  'actions/commit.md' \
  'git rev-parse --git-path do-work-commit-secret-quarantine' \
  'actions/commit.md must re-derive a deterministic Git-private quarantine across separate command blocks.'

assert_contains \
  'actions/inspect.md' \
  'git rev-parse --git-path do-work-inspect-secret-quarantine' \
  'actions/inspect.md must re-derive a deterministic Git-private quarantine across separate command blocks.'

assert_contains \
  'actions/commit.md' \
  'already staged.*exact.*deletion|exact.*deletion.*already staged' \
  'actions/commit.md Step 5 must recognize an exact already-staged XD deletion before deciding whether git add -u is needed.'

assert_contains \
  'actions/commit.md' \
  'skip.*git add -u|without.*git add -u' \
  'actions/commit.md Step 5 must skip git add -u when cached name/status already proves the XD deletion.'

assert_contains \
  'actions/commit.md' \
  'otherwise.*stages.*verif|otherwise it stages.*verif' \
  'actions/commit.md Step 5 must retain the unstaged-deletion branch and verify deletion-only metadata after git add -u.'

# Argument parsing must fail fast. The watchdog keeps a regression from hanging
# the entire suite forever: the original `shift 2` with one argument left made
# the loop reread --repo-root indefinitely on Bash 3.2 and newer alike.
associate_missing_value_output="$(mktemp)"
"$core_root/tools/checks/associate-files.sh" --repo-root </dev/null \
  >"$associate_missing_value_output" 2>&1 &
associate_missing_value_pid=$!
(
  sleep 2
  if kill -0 "$associate_missing_value_pid" 2>/dev/null; then
    kill "$associate_missing_value_pid" 2>/dev/null || true
  fi
) &
associate_missing_value_watchdog_pid=$!
if wait "$associate_missing_value_pid"; then
  associate_missing_value_exit=0
else
  associate_missing_value_exit=$?
fi
kill "$associate_missing_value_watchdog_pid" 2>/dev/null || true
wait "$associate_missing_value_watchdog_pid" 2>/dev/null || true
if [ "$associate_missing_value_exit" -ne 2 ]; then
  printf 'FAIL: tools/checks/associate-files.sh --repo-root with no value must fail promptly with exit 2, got %s.\n' \
    "$associate_missing_value_exit" >&2
  fail_count=$((fail_count + 1))
fi
rm -f "$associate_missing_value_output"

# The status alias belongs to the Schema Read Contract, but this shell reader
# owns its own terminal-success predicate. Exercise the new alias against a
# real REQ fixture so the prose and helper cannot drift back apart.
assert_contains \
  "actions/work-reference.md" \
  '`complete`/`done`/`finished`/`closed`.*`completed`' \
  'actions/work-reference.md must document complete with the other aliases that normalize to completed before terminal-success readers consume the status.'

associate_complete_probe_dir="$(mktemp -d)"
mkdir -p "$associate_complete_probe_dir/do-work/archive/UR-301"
cat > "$associate_complete_probe_dir/do-work/archive/UR-301/REQ-501-legacy-complete.md" <<'EOF'
---
id: REQ-501
status: complete
completed_at: 2026-08-07T12:00:00Z
---

## Implementation Summary

**Files changed:**
- `legacy-file.txt` (modified)
EOF
if ! associate_complete_output="$(printf 'legacy-file.txt\n' | "$core_root/tools/checks/associate-files.sh" --repo-root "$associate_complete_probe_dir")" || ! printf '%s\n' "$associate_complete_output" | grep -qxF "$(printf 'REQ-501\tlegacy-file.txt')"; then
  printf 'FAIL: tools/checks/associate-files.sh must associate a status: complete REQ after normalizing the documented terminal-success alias.\n' >&2
  fail_count=$((fail_count + 1))
fi
rm -rf -- "$associate_complete_probe_dir"

# A git-status failure is not a clean tree. Process substitution used to hide
# the producer's failure, making a bare repository return the clean-tree exit 1.
inventory_failure_probe_dir="$(mktemp -d)"
if ! GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
  git init -q --bare "$inventory_failure_probe_dir"; then
  printf 'FAIL: could not set up the uncommitted-inventory git-status failure probe.\n' >&2
  fail_count=$((fail_count + 1))
else
  if inventory_failure_output="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null \
    "$core_root/tools/checks/uncommitted-inventory.sh" "$inventory_failure_probe_dir" 2>&1)"; then
    inventory_failure_exit=0
  else
    inventory_failure_exit=$?
  fi
  if [ "$inventory_failure_exit" -ne 2 ]; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must return exit 2 when git status fails, got %s.\n' \
      "$inventory_failure_exit" >&2
    fail_count=$((fail_count + 1))
  fi
  if ! grep -qF 'STATUS-FAILED:' <<<"$inventory_failure_output"; then
    printf 'FAIL: tools/checks/uncommitted-inventory.sh must diagnose a git-status failure instead of reporting a clean tree.\n' >&2
    fail_count=$((fail_count + 1))
  fi
fi
rm -rf -- "$inventory_failure_probe_dir"

assert_contains \
  "tools/checks/preflight.sh" \
  'git -c status\.renames=copies status --porcelain=v1 --untracked-files=all -z' \
  'tools/checks/preflight.sh must keep NUL-safe repository-wide dirty-file detection — serial qualification and review read the repository-wide working/staged diff.'

preflight_dirty_warning_line="$(grep -E '^[[:space:]]*echo "WARN: .*uncommitted changes' "$core_root/tools/checks/preflight.sh" || true)"

assert_block_contains \
  "$preflight_dirty_warning_line" \
  'unless they prevent the active REQ' \
  'tools/checks/preflight.sh dirty-file warning must stay subordinate to Current-REQ relevance instead of turning unexpected state into a blocker.'

assert_block_contains \
  "$preflight_dirty_warning_line" \
  'qualification/review evidence' \
  'tools/checks/preflight.sh must explain dirty files as possible qualification/review evidence contamination, not as files the explicit commit step will automatically stage.'

assert_block_not_contains \
  "$preflight_dirty_warning_line" \
  'may stage unrelated files|swept into (the )?commit' \
  'tools/checks/preflight.sh must not claim dirty files are automatically staged — Step 9 stages explicit paths.'

preflight_template_block="$(sed -n '/^## Pre-Flight Template/,/^## Implementation Summary Template/p' "$core_root/actions/work-reference.md")"

assert_block_contains \
  "$preflight_template_block" \
  'preserve/exclude from this REQ.*qualification/review evidence' \
  'actions/work-reference.md Pre-Flight template must frame pre-existing dirty files as qualification/review evidence contamination.'

assert_file_not_contains \
  "actions/work.md" \
  'unrelated dirty files get swept into the commit' \
  'actions/work.md must not teach that pre-existing dirty files are swept into the explicit per-REQ commit.'

# Review regressions: prescribed shell and roadmap classification are runtime
# contracts even though they live in Markdown/just recipes rather than compiled code.
extract_kanban_shutdown_line() {
  awk '
    /^run-kanban \$port=/ {
      getline
      getline
      sub(/^[[:space:]]*/, "")
      print
      exit
    }
  ' "$repo_root/$1"
}

root_kanban_shutdown_line="$(extract_kanban_shutdown_line Justfile)"
installer_kanban_shutdown_line="$(extract_kanban_shutdown_line skills/do-work-board/justfile.template)"
if [ "$root_kanban_shutdown_line" != "$installer_kanban_shutdown_line" ]; then
  printf 'FAIL: Justfile and the board-owned template must carry one identical run-kanban shutdown line.\n' >&2
  fail_count=$((fail_count + 1))
fi

for kanban_recipe_file in "skills/do-work-board/justfile.template" "Justfile"; do
  assert_file_not_contains \
    "$kanban_recipe_file" \
    'case "\$listener_command" in \*queue-kanban\*' \
    "$kanban_recipe_file must not identify a stale board from arbitrary argv text."
  assert_contains \
    "$kanban_recipe_file" \
    'lsof -a -p "\$listener_pid" -d txt -Fn' \
    "$kanban_recipe_file must identify a stale board from its executable, preserving cross-repo binary names without matching unrelated arguments."
  assert_contains \
    "$kanban_recipe_file" \
    'wait_count.*-lt 320.*lsof -a -p "\$listener_pid".*tcp:"\$port".*-sTCP:LISTEN' \
    "$kanban_recipe_file must bound a listener-specific wait at 320 iterations instead of polling process existence."
  assert_contains \
    "$kanban_recipe_file" \
    'remaining_listener_pid=.*lsof -ti tcp:"\$port".*-sTCP:LISTEN' \
    "$kanban_recipe_file must query the port again after the bounded shutdown wait."
  assert_contains \
    "$kanban_recipe_file" \
    'ps -p "\$remaining_listener_pid" -o args=' \
    "$kanban_recipe_file must resolve the remaining listener's full command for the refusal."
  assert_file_not_contains \
    "$kanban_recipe_file" \
    'while kill -0 "\$listener_pid"' \
    "$kanban_recipe_file must wait on listener ownership, not process existence."
  assert_contains \
    "$kanban_recipe_file" \
    '^run-do-work-update:' \
    "$kanban_recipe_file must ship the project-local do-work update shortcut."
done

# Execute the canonical shutdown line with command seams. A queue-kanban PID that remains a
# listener throughout the bounded wait must make the recipe line fail before build+serve, and
# the diagnosis must name both the PID and full command. sleep is a no-op in the fixture, so all
# 320 iterations run without slowing the suite.
if stuck_listener_output="$(
  port=8090
  lsof() {
    if [ "${4:-}" = "-d" ]; then
      printf 'n/tmp/queue-kanban\n'
    else
      printf '4242\n'
    fi
  }
  ps() { printf '/tmp/queue-kanban serve --port 8090\n'; }
  kill() { return 0; }
  sleep() { return 0; }
  eval "$root_kanban_shutdown_line" 2>&1
)"; then
  stuck_listener_status=0
else
  stuck_listener_status=$?
fi
if [ "$stuck_listener_status" -ne 1 ]; then
  printf 'FAIL: run-kanban shutdown must refuse startup when a listener remains after the bounded wait; got exit %s.\n' "$stuck_listener_status" >&2
  fail_count=$((fail_count + 1))
fi
if ! printf '%s\n' "$stuck_listener_output" | grep -qF 'pid 4242'; then
  printf 'FAIL: run-kanban stuck-listener refusal must name pid 4242.\n' >&2
  fail_count=$((fail_count + 1))
fi
if ! printf '%s\n' "$stuck_listener_output" | grep -qF '/tmp/queue-kanban serve --port 8090'; then
  printf 'FAIL: run-kanban stuck-listener refusal must name the listener command.\n' >&2
  fail_count=$((fail_count + 1))
fi

# Preserve the older safety boundary: a foreign executable is refused immediately and never
# passed to kill. This behavior probe complements the executable-identity contract above.
foreign_kill_marker="$(mktemp)"
if foreign_listener_output="$(
  port=8090
  # The shutdown recipe reads this through eval, which ShellCheck cannot trace.
  : "$port"
  lsof() {
    if [ "${4:-}" = "-d" ]; then
      printf 'n/usr/bin/python3\n'
    else
      printf '3131\n'
    fi
  }
  ps() { printf '/usr/bin/python3 -m http.server 8090\n'; }
  kill() { printf 'called\n' > "$foreign_kill_marker"; }
  sleep() { return 0; }
  eval "$root_kanban_shutdown_line" 2>&1
)"; then
  foreign_listener_status=0
else
  foreign_listener_status=$?
fi
# This fixture asserts the refusal status and kill boundary; its captured text is intentionally ignored.
: "$foreign_listener_output"
if [ "$foreign_listener_status" -ne 1 ] || [ -s "$foreign_kill_marker" ]; then
  printf 'FAIL: run-kanban must refuse a foreign listener without calling kill.\n' >&2
  fail_count=$((fail_count + 1))
fi
rm -f "$foreign_kill_marker"

if [ ! -x "$core_root/tools/do-work-update.sh" ]; then
  printf 'FAIL: modular core tools/do-work-update.sh must be executable for the Just shortcut.\n' >&2
  fail_count=$((fail_count + 1))
fi
assert_contains \
  "tools/do-work-update.sh" \
  "--project-root" \
  'tools/do-work-update.sh must derive and validate the consuming project root before updating.'
assert_contains \
  "tools/do-work-update.sh" \
  'tar xzf "\$upstream_tarball" -C "\$fresh_upstream"' \
  'tools/do-work-update.sh must extract only into staging; behavioral probes verify runtime do-work data is outside every managed destination.'
assert_contains \
  "tools/install-do-work-suite.sh" \
  'Install this complete four-skill suite' \
  'The installed suite transaction must require confirmation after showing its reviewed diff.'
assert_file_not_contains \
  "tools/do-work-update.sh" \
  'cp -R "\$skill_root"' \
  'tools/do-work-update.sh must not reintroduce the pre-update rollback copy — git is the undo, and a duplicated tree on every run buys nothing git does not already hold. A mid-update failure reports the partial install instead; see _dev/tests/update-script-behavior.sh.'
assert_contains \
  "tools/install-do-work-suite.sh" \
  'recover_install' \
  'The installed suite transaction must automatically recover its validated managed paths after a destructive-region failure.'

assert_file_not_contains \
  "actions/work.md" \
  'else probe_wrapper=""' \
  'actions/work.md must not drop the blocked-check time limit when timeout/gtimeout is unavailable.'

assert_contains \
  "actions/work.md" \
  'probe_exit=124' \
  'actions/work.md must preserve a bounded portable fallback and report a timed-out blocked check as exit 124.'

blocked_probe_shell_block="$(sed -n '/^# Re-derive paths deterministically/,/^rm -f "\$BLOCKED_CHECK_SCRIPT"/p' "$core_root/actions/work.md")"
if ! bash -n <<<"$blocked_probe_shell_block"; then
  printf 'FAIL: actions/work.md blocked-check shell block must remain syntactically valid.\n' >&2
  fail_count=$((fail_count + 1))
fi

# Target ID Resolution contract (REQ-067). The run tokenizer recognized exactly one token
# shape (REQ- + digits); every other token was residue by construction, so a UR argument was
# handled by whatever the reading agent improvised. The shared contract in work-reference.md is
# the single definition of both token shapes and UR->REQ expansion that run/abandon/
# roadmap cite instead of restating; work.md's Input and usage string must offer the UR- shape.
assert_contains \
  "actions/work-reference.md" \
  '### Target ID Resolution' \
  'actions/work-reference.md must define the shared Target ID Resolution contract (REQ-/UR- token shapes + UR->REQ expansion by user_request: scan) that the id-taking actions cite instead of restating.'

work_input_block="$(sed -n '/^## Input/,/^## Steps/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_input_block" \
  'UR-NNN' \
  'actions/work.md Input must accept UR-NNN targeting tokens (usage string offers UR-NNN and the tokenizer recognizes the UR- shape) per the Target ID Resolution contract — a UR argument must no longer fall through to the unrecognized-argument guard.'

# UR ids accepted by abandon (REQ-068). The action keyed entirely on REQ-NNN tokens; a UR
# argument had no defined handling (abandon's globs substitute a REQ number). Its Input section
# must name the UR- shape and cite the shared Target ID Resolution contract rather than restating
# the resolution rule. Scope the UR- check to the Input section — abandon.md already names
# archive/UR-NNN/ folders elsewhere, so a file-wide grep would pass vacuously without the action
# actually accepting a UR argument.
abandon_input_block="$(sed -n '/^## Input/,/^## Steps/p' "$core_root/actions/abandon.md")"

assert_block_contains \
  "$abandon_input_block" \
  'Target ID Resolution' \
  'actions/abandon.md Input must cite the Target ID Resolution contract so it accepts UR-NNN targeting tokens (a UR cancels its cancellable members) rather than only REQ-NNN.'

# A legacy failed REQ can already sit inside a closed archive/UR-NNN/ folder. The board still
# treats that status as active, so explicitly targeting the REQ must reach abandon's existing
# confirmed in-place failed->cancelled write without moving or reopening the closed UR.
assert_contains \
  "actions/abandon.md" \
  'status: failed.*archive/UR-NNN/.*cancellable in place' \
  'actions/abandon.md must let an explicitly targeted failed REQ inside archive/UR-NNN/ use the confirmed in-place cancellation path; refusing it leaves a permanent active board card.'

assert_file_not_contains \
  "actions/abandon.md" \
  'status: failed.*archive/UR-NNN/.*refuse' \
  'actions/abandon.md must not restore the legacy nested-failure refusal after the in-place resolution path was extended to that shape.'

assert_contains \
  "actions/forensics.md" \
  'archive/UR-NNN/.*cancel.*in place' \
  'actions/forensics.md must route a failed REQ already inside archive/UR-NNN/ to explicit in-place abandon instead of claiming no action is needed.'

assert_contains \
  "actions/roadmap.md" \
  '^-[[:space:]]+\*\*Ready\*\*[[:space:]]+— normalized `status` is `pending`' \
  'actions/roadmap.md must require pending status before classifying a queued REQ as Ready.'

# REQ scoping in roadmap (REQ-070). roadmap accepted UR-NNN as a scope token but not REQ-NNN, so
# `do-work roadmap REQ-067` silently fell through to a whole-queue survey — the inverse of the
# UR-run asymmetry the batch fixed. Its Input must now name a REQ-NNN scope token and cite the
# shared Target ID Resolution contract for the token shape (rather than restating it).
roadmap_input_block="$(sed -n '/^## Input/,/^## Steps/p' "$core_root/actions/roadmap.md")"

assert_block_contains \
  "$roadmap_input_block" \
  'REQ-NNN' \
  'actions/roadmap.md Input must accept a REQ-NNN scope token (single-REQ survey), not only UR-NNN — a recognized REQ id must no longer fall through to the full-survey default.'

assert_block_contains \
  "$roadmap_input_block" \
  'Target ID Resolution' \
  'actions/roadmap.md Input must cite the shared Target ID Resolution contract for the REQ-/UR- token shapes rather than restating them.'

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
  "actions/setup-memory.md" \
  'install memory-module' \
  'knowledge setup-memory must retain the documented memory-module migration alias (ADR-017).'

assert_contains \
  "actions/setup-memory.md" \
  'settings\.json\.pre-memory-module' \
  'knowledge setup-memory hook merge must back up settings.json before composing entries.'

assert_contains \
  "hooks/memory-stop-capture.sh" \
  'REDACTED' \
  'hooks/memory-stop-capture.sh must redact credential-shaped text before persisting captures — defense in depth behind the machine-local store.'

# Raw captures and the per-machine ledger must never become committable: the installer
# adds them to .git/info/exclude (machine-local), never the project's .gitignore.
assert_contains \
  "actions/setup-memory.md" \
  '\*\*/memory/logs/' \
  'knowledge setup-memory must add memory/logs/ to .git/info/exclude — verbatim captures must not be committable.'

assert_contains \
  "actions/setup-memory.md" \
  '\*\*/memory/usage-ledger\.jsonl' \
  'knowledge setup-memory must add memory/usage-ledger.jsonl to .git/info/exclude.'

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

# do-work/ and kb/ are TRACKED in this repo (they are the same Trail of Intent the skill tells
# consumers to commit, and the tracked-path-only data-loss guards in tools/checks/record-commit-hash.sh
# plus cleanup Pass 6 blanked-REQ recovery only work on tracked REQs). That makes these two
# export-ignore lines the only barrier between the maintainer's queue/KB and every consumer
# install — the pre-0.157.0 blanket .git/info/exclude entry no longer backstops them.
assert_contains \
  ".gitattributes" \
  '^/do-work[[:space:]]+export-ignore' \
  '.gitattributes must export-ignore /do-work — this repo tracks its own queue, so without this line the maintainer archive ships to every consumer install.'

assert_contains \
  ".gitattributes" \
  '^/kb[[:space:]]+export-ignore' \
  '.gitattributes must export-ignore /kb — this repo tracks its own knowledge base, so without this line it ships to every consumer install.'

assert_contains \
  "tools/install-do-work-suite.sh" \
  'module_relatives' \
  'The installed suite transaction must construct an explicit managed module plan; behavioral probes verify the project knowledge base is outside it.'

# Shipped files must not cite the skill's own CLAUDE.md/AGENTS.md — those files are absent
# downstream, so a citation dangles. The full rule lives in CLAUDE.md → Action File Conventions.
#
# The check is INVERTED on purpose: it flags ANY mention of CLAUDE.md/AGENTS.md in a shipped
# path, and exempts a short per-file allowlist. It used to enumerate citation idioms
# (`see CLAUDE.md`, `per CLAUDE.md`, `CLAUDE.md →`) and that shape failed exactly as
# "Closed Enumerations Go Stale" predicts: scored against the seven real occurrences in the
# tree it caught ZERO — including the `actions/memory-reference.md` line REQ-088 was filed to
# fix, and six `CLAUDE.md § Before Every Commit` comments in tools/queue-kanban that shipped
# for months. Extending the idiom list would have caught 4 of 6 and false-positived on
# actions/prime.md, where the colon is sentence punctuation. A mention is cheap to detect and
# impossible to phrase around; the judgement lives in the allowlist, where it is visible.
#
# The allowlist is PER-FILE, never per-directory: allowlisting `actions/` wholesale would have
# exempted actions/memory-reference.md, the file that started this. Entries are files whose
# subject genuinely IS a consumer project's own CLAUDE.md/AGENTS.md (prime registries, the KB
# schema file, and tidy-repo's layout rules) — those
# references are correct and must not be "fixed". A new shipped file that mentions the
# maintainer doc fails this check until someone decides which of the two it is; that decision
# is the point.
shipped_citation_paths=(SKILL.md next-steps.md README.md actions crew-members prompts interviews specs docs hooks tools)
maintainer_doc_mention_allowlist=(
  actions/prime.md
  docs/prime-guide.md
  actions/version.md
  actions/tidy-repo.md
  actions/bkb.md
  actions/bkb-reference.md
  docs/bkb-guide.md
  actions/capture.md
  actions/validate-feedback.md
  actions/prompts.md
  README.md
  prompts/README.md
  prompts/prompt-kit-step2-personal-context-doc.md
)
maintainer_doc_mentions="$(cd "$repo_root" && grep -rIn 'CLAUDE\.md\|AGENTS\.md' "${shipped_citation_paths[@]}" 2>/dev/null || true)"
unallowed_maintainer_doc_mentions=""
while IFS= read -r maintainer_doc_hit; do
  [ -n "$maintainer_doc_hit" ] || continue
  maintainer_doc_hit_path="${maintainer_doc_hit%%:*}"
  maintainer_doc_hit_allowed=0
  for allowlisted_citation_file in "${maintainer_doc_mention_allowlist[@]}"; do
    if [ "$maintainer_doc_hit_path" = "$allowlisted_citation_file" ]; then
      maintainer_doc_hit_allowed=1
      break
    fi
  done
  if [ "$maintainer_doc_hit_allowed" -eq 0 ]; then
    unallowed_maintainer_doc_mentions="${unallowed_maintainer_doc_mentions}${maintainer_doc_hit}"$'\n'
  fi
done <<< "$maintainer_doc_mentions"
if [ -n "$unallowed_maintainer_doc_mentions" ]; then
  printf 'FAIL: shipped files must not mention the skill'\''s own CLAUDE.md/AGENTS.md (export-ignored — absent in consumer installs). Restate the rule inline, point at a shipped home, or — if the mention really is about a *consumer project'\''s* CLAUDE.md — add the file to maintainer_doc_mention_allowlist in this script:\n%s' "$unallowed_maintainer_doc_mentions" >&2
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
  interview.md kb-lessons-handoff.md note.md present-work.md
  prime.md prompts.md quick-wins.md
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

for action_file_path in "$repo_root"/skills/do-work*/actions/*.md; do
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
    | CLAUDE_PROJECT_DIR="$redaction_order_workdir" bash "$knowledge_root/hooks/memory-stop-capture.sh" >/dev/null 2>&1 || true
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

# Board overlap annotation stays display-only (REQ-052, follow-up to REQ-034). The
# invariant is structural: `annotateWriteSetOverlap` runs AFTER `bucketColumns` in
# buildBoard, so no column-placement code can read the annotation. That ordering is one
# movable line — a future edit that hoists the call above bucketing (to "reuse" the
# overlap set while placing cards) turns a display badge into scheduling logic, and the
# co-dispatch decision belongs to actions/work.md Step 1's gate, not the viewer. Assert
# the real call-site order inside buildBoard rather than a comment: a comment survives
# the move. Paired with the instruction-side claim in actions/board.md, whose rewording
# is the other way the invariant quietly stops being a promise.
build_board_block="$(sed -n '/^func buildBoard(/,/^}/p' "$board_root/tools/queue-kanban/model.go")"
bucket_columns_call_line="$(grep -nF 'bucketColumns(board.AllRequests' <<<"$build_board_block" | head -1 | cut -d: -f1 || true)"
overlap_annotation_call_line="$(grep -nF 'annotateWriteSetOverlap(board.AllRequests)' <<<"$build_board_block" | head -1 | cut -d: -f1 || true)"

if [ -z "$bucket_columns_call_line" ] || [ -z "$overlap_annotation_call_line" ]; then
  printf 'FAIL: tools/queue-kanban/model.go buildBoard must call both bucketColumns(board.AllRequests ...) and annotateWriteSetOverlap(board.AllRequests) — one call site was renamed or removed, so the display-only ordering of the write-set overlap annotation can no longer be verified. Fix: restore the call sites (annotation last) or update this anchor in the same commit.\n' >&2
  fail_count=$((fail_count + 1))
elif [ "$overlap_annotation_call_line" -lt "$bucket_columns_call_line" ]; then
  printf 'FAIL: tools/queue-kanban/model.go calls annotateWriteSetOverlap BEFORE bucketColumns in buildBoard — the write-set overlap annotation is display only and must stay structurally unable to affect column placement. Fix: move the annotateWriteSetOverlap call back below bucketColumns; if the board really must schedule on write sets, that decision belongs to actions/work.md Step 1 dispatch gate.\n' >&2
  fail_count=$((fail_count + 1))
fi

board_rules_block="$(sed -n '/^## Rules/,/^## Common Rationalizations/p' "$board_root/actions/board.md")"

assert_block_contains \
  "$board_rules_block" \
  '`annotateWriteSetOverlap` in `tools/queue-kanban/model\.go` runs \*after\* column bucketing' \
  'actions/board.md Rules must keep the claim that annotateWriteSetOverlap runs after column bucketing — the ordering ratchet above pins the code, this pins the promise a parser-editing agent reads. Fix: restate it in the parser lock-step rule.'

assert_block_contains \
  "$board_rules_block" \
  'never column logic' \
  'actions/board.md Rules must keep the overlap annotation display-only claim (drives the overlaps badge and drawer row, never column logic, never blocking) — without it nothing tells a parser-editing agent that which REQs may co-dispatch stays with actions/work.md Step 1 gate.'

# Mid-run answer durability. A question asked interactively can satisfy every wording rule
# and still produce an answer that dies with the session: a consumer's long-running
# orchestrator asked two user-owned decisions, got detailed answers, wrote nothing to disk,
# and the next builder re-decided one of them in a fresh context off the stored
# `Recommended:` rationale the user had just rejected (validate-feedback 2026-08-01). Three
# halves hold the fix together — the principle, the pipeline branch that obeys it, and the
# one named format both cite — and each reads as removable prose on its own.
assert_contains \
  "crew-members/clear-questions.md" \
  'outlive the transcript' \
  'crew-members/clear-questions.md must keep the principle that an interactively obtained answer is written into the durable record before it is acted on — wording rules alone let a compliant question produce an answer that dies with the session.'

work_open_questions_block="$(sed -n '/^### Step 3\.5: Open Questions/,/^### Step 3\.7/p' "$core_root/actions/work.md")"

assert_block_contains \
  "$work_open_questions_block" \
  'escalated to the user mid-run and answered' \
  'actions/work.md Step 3.5 must keep the third branch for a decision escalated mid-run and answered (orchestrator writes it in before dispatch) — with only builder-decides and defer-to-clarify, a mid-run answer has nowhere to land.'

assert_block_contains \
  "$work_open_questions_block" \
  '\*\*not\*\* `- \[~\]`' \
  'actions/work.md Step 3.5 must keep a mid-run user answer recorded as - [x] and never - [~] — filing it as a builder decision sends a settled question back through clarify and leaves the rejected Recommended: rationale standing as its reason.'

assert_contains \
  "actions/clarify.md" \
  'Canonical answered-question format' \
  'actions/clarify.md Step 4 must keep the named entry point declaring its - [x] form canonical for any caller that obtains a user answer — cited by name (not step number) from clear-questions.md Principle 8 and work.md Step 3.5.'

# The suite manifest is executable input to both update paths and the fresh installer. Keep
# its path-safety and exact-module contract in one behavioral probe rather than duplicating
# grep assertions for each caller.
suite_manifest_probe="$repo_root/_dev/tests/suite-manifest-contract.sh"
if [ ! -f "$suite_manifest_probe" ]; then
  printf 'FAIL: _dev/tests/suite-manifest-contract.sh is missing — the suite layout has no behavioral coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$suite_manifest_probe"; then
  printf 'FAIL: suite manifest contract probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Published package links have two valid execution contexts: this source tree and the
# manifest-mapped install under .claude/skills/. The local probe parses rendered Markdown
# outside code, validates first-party raw/blob paths without network access, and keeps the
# installed core changelog byte-identical to the release source.
shipped_reference_probe="$repo_root/_dev/tests/shipped-package-reference-contract.sh"
if [ ! -f "$shipped_reference_probe" ]; then
  printf 'FAIL: _dev/tests/shipped-package-reference-contract.sh is missing — shipped package references have no source/install coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$shipped_reference_probe"; then
  printf 'FAIL: shipped package reference contract failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Markdown-prescribed shell is shipped executable guidance, while hook scripts are shipped
# executable code. Exercise both through one attributable syntax/lint probe, including its
# negative fixture so a silently-neutered extractor cannot report a decorative green result.
action_shell_probe="$repo_root/_dev/tests/action-shell-blocks.sh"
if [ ! -x "$action_shell_probe" ]; then
  printf 'FAIL: _dev/tests/action-shell-blocks.sh is missing or not executable — shipped shell guidance has no lint coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$action_shell_probe" --self-test; then
  printf 'FAIL: action shell-block negative self-test failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$action_shell_probe"; then
  printf 'FAIL: shipped shell-block lint failed (see the attributed FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# The SessionStart banner is deliberately fail-soft: malformed or missing runtime inputs must
# produce its unknown/zero defaults, not let shell options terminate the hook before output.
session_start_probe="$repo_root/_dev/tests/session-start-hook-behavior.sh"
if [ ! -x "$session_start_probe" ]; then
  printf 'FAIL: _dev/tests/session-start-hook-behavior.sh is missing or not executable — the startup banner fallback has no behavioral coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$session_start_probe"; then
  printf 'FAIL: SessionStart hook behavior probes failed (see the fixture FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Shared prescribed-shell rationale has one shipped home; callers retain only local intent,
# commands, and explicit pointers so a primitive fix cannot drift across prose copies again.
prescribed_shell_probe="$repo_root/_dev/tests/prescribed-shell-canonicalization.sh"
if [ ! -x "$prescribed_shell_probe" ]; then
  printf 'FAIL: _dev/tests/prescribed-shell-canonicalization.sh is missing or not executable — shell primitive restatements have no ratchet.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$prescribed_shell_probe"; then
  printf 'FAIL: prescribed shell primitive canonicalization failed (see the attributed FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Behavioral probes for tools/checks/record-commit-hash.sh. Kept in their own file because
# they build a throwaway git repo and run the real script rather than grepping prose — but
# invoked from here, since nothing auto-discovers _dev/tests/*.sh and an uninvoked probe file
# is dead weight that reads as coverage.
record_commit_hash_probe="$repo_root/_dev/tests/record-commit-hash-guards.sh"
if [ ! -f "$record_commit_hash_probe" ]; then
  printf 'FAIL: _dev/tests/record-commit-hash-guards.sh is missing — the write-back guards have no behavioral coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash <(sed \
  -e "s|^repo_root=.*|repo_root=\"$repo_root\"|" \
  -e 's|\$repo_root/tools/checks/|\$repo_root/skills/do-work/tools/checks/|g' \
  "$record_commit_hash_probe"); then
  printf 'FAIL: record-commit-hash guard probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Behavioral probes for tools/do-work-update.sh — same reasoning as above, and the same
# no-auto-discovery caveat. These build a synthetic install plus a stubbed upstream fetch
# because the no-rollback-copy contract is a runtime property (no `.bak` left behind, a
# mid-update failure that reports instead of restoring) that no grep can assert.
update_script_probe="$repo_root/_dev/tests/update-script-behavior.sh"
if [ ! -f "$update_script_probe" ]; then
  printf 'FAIL: _dev/tests/update-script-behavior.sh is missing — the updater has no behavioral coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$update_script_probe"; then
  printf 'FAIL: update-script behavior probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# The live modular suite package boundaries need a contract independent from the active
# root bootstrap tools. This checks the staged router, required core
# runtime, hook targets, and the ban on leaking repository maintainer instructions.
staged_skills_probe="$repo_root/_dev/tests/staged-skills-contract.sh"
if [ ! -f "$staged_skills_probe" ]; then
  printf 'FAIL: _dev/tests/staged-skills-contract.sh is missing — staged skill packages have no boundary coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$staged_skills_probe"; then
  printf 'FAIL: staged skills contract probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Fresh installation is a four-module/configuration transaction with its own hermetic
# bootstrap. Keep these probes separate because they build several isolated Git clients,
# including malformed, fallback, rollback, and interruption states.
suite_installer_probe="$repo_root/_dev/tests/install-suite-behavior.sh"
if [ ! -f "$suite_installer_probe" ]; then
  printf 'FAIL: _dev/tests/install-suite-behavior.sh is missing — the full-suite installer has no behavioral coverage.\n' >&2
  fail_count=$((fail_count + 1))
elif ! bash "$suite_installer_probe"; then
  printf 'FAIL: suite installer behavior probes failed (see the FAIL lines above).\n' >&2
  fail_count=$((fail_count + 1))
fi

# Managed Just sections are a byte-preserving ownership boundary, not a prose convention.
# Exercise the real utility across replacement, append, creation, malformed
# markers, filename variants, spaces, modes, idempotence, and Just parsing.
replace_section_tool="$repo_root/tools/replace-text-section.sh"
if [ ! -x "$replace_section_tool" ]; then
  printf 'FAIL: tools/replace-text-section.sh is missing or not executable — managed recipe ownership has no implementation.\n' >&2
  fail_count=$((fail_count + 1))
else
  section_workdir="$(mktemp -d)"
  section_file="$section_workdir/managed-section.just"
  template_file="$section_workdir/complete-template.just"
  printf '# >>> do-work:recipes >>>\nmanaged-probe:\n    echo managed\n# <<< do-work:recipes <<<\n' > "$section_file"
  printf 'set shell := ["bash", "-cu"]\n\n# >>> do-work:recipes >>>\nmanaged-probe:\n    echo managed\n# <<< do-work:recipes <<<\n' > "$template_file"
  chmod 750 "$template_file"

  byte_target="$section_workdir/project with spaces/Justfile"
  mkdir -p "$(dirname "$byte_target")"
  printf 'prefix\000byte\n# >>> do-work:recipes >>>\nold:\n    echo old\n# <<< do-work:recipes <<<\nsuffix\n' > "$byte_target"
  chmod 640 "$byte_target"
  expected_target="$section_workdir/expected-byte-target"
  printf 'prefix\000byte\n# >>> do-work:recipes >>>\nmanaged-probe:\n    echo managed\n# <<< do-work:recipes <<<\nsuffix\n' > "$expected_target"
  if ! "$replace_section_tool" --target "$byte_target" --section-file "$section_file"; then
    printf 'FAIL: replace-text-section could not replace one valid managed section.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$byte_target" "$expected_target"; then
    printf 'FAIL: replace-text-section changed bytes outside the managed section or wrote the wrong replacement.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  target_mode="$(stat -f '%Lp' "$byte_target" 2>/dev/null || stat -c '%a' "$byte_target" 2>/dev/null || true)"
  if [ "$target_mode" != 640 ]; then
    printf 'FAIL: replace-text-section changed the existing target mode (got %s, want 640).\n' "$target_mode" >&2
    fail_count=$((fail_count + 1))
  fi
  cp "$byte_target" "$section_workdir/idempotent-snapshot"
  if ! "$replace_section_tool" --target "$byte_target" --section-file "$section_file" \
     || ! cmp -s "$byte_target" "$section_workdir/idempotent-snapshot"; then
    printf 'FAIL: replace-text-section is not byte-idempotent on repeated execution.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  absent_target="$section_workdir/absent project/.justfile"
  mkdir -p "$(dirname "$absent_target")"
  if ! "$replace_section_tool" --target "$absent_target" --section-file "$section_file" --template-file "$template_file" \
     || ! cmp -s "$absent_target" "$template_file"; then
    printf 'FAIL: replace-text-section did not create an absent target from the complete supplied template.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  absent_mode="$(stat -f '%Lp' "$absent_target" 2>/dev/null || stat -c '%a' "$absent_target" 2>/dev/null || true)"
  if [ "$absent_mode" != 750 ]; then
    printf 'FAIL: replace-text-section did not preserve the complete template mode on create (got %s, want 750).\n' "$absent_mode" >&2
    fail_count=$((fail_count + 1))
  fi

  for justfile_variant in justfile Justfile .justfile; do
    variant_target="$section_workdir/variants/$justfile_variant"
    mkdir -p "$(dirname "$variant_target")"
    printf 'custom-%s:\n    echo untouched\n' "$justfile_variant" > "$variant_target"
    if ! "$replace_section_tool" --target "$variant_target" --section-file "$section_file"; then
      printf 'FAIL: replace-text-section could not append to marker-free %s.\n' "$justfile_variant" >&2
      fail_count=$((fail_count + 1))
    elif [ "$(grep -c '^# >>> do-work:recipes >>>$' "$variant_target")" -ne 1 ] \
      || ! grep -q "^custom-$justfile_variant:" "$variant_target"; then
      printf 'FAIL: replace-text-section duplicated ownership or changed custom content in %s.\n' "$justfile_variant" >&2
      fail_count=$((fail_count + 1))
    fi
  done

  reserved_section_file="$section_workdir/reserved-section.just"
  awk '
    $0 == "# >>> do-work:recipes >>>" { inside=1 }
    inside { print }
    $0 == "# <<< do-work:recipes <<<" { found_end=1; exit }
    END { if (!inside || !found_end) exit 1 }
  ' "$repo_root/skills/do-work-board/justfile.template" > "$reserved_section_file"

  collision_index=0
  for reserved_recipe_name in run-kanban run-kanban-cli kanban-static kanban-summary run-do-work-update; do
    collision_index=$((collision_index + 1))
    collision_target="$section_workdir/collision-$collision_index.just"
    case "$reserved_recipe_name" in
      run-kanban)
        printf 'run-kanban:\n    echo collision\n' > "$collision_target"
        ;;
      run-kanban-cli)
        printf '@run-kanban-cli $view="open:all": dependency-recipe\n    echo collision\n' > "$collision_target"
        ;;
      kanban-static)
        printf 'kanban-static destination="build/output": dependency-recipe\n    echo collision\n' > "$collision_target"
        ;;
      kanban-summary)
        printf 'alias kanban-summary := custom-summary\n' > "$collision_target"
        ;;
      run-do-work-update)
        printf 'run-do-work-update:\r\n    echo collision\r\n' > "$collision_target"
        ;;
    esac
    cp "$collision_target" "$collision_target.before"
    collision_output="$section_workdir/collision-$collision_index.out"
    if "$replace_section_tool" --target "$collision_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$collision_output" 2>&1; then
      printf 'FAIL: replace-text-section accepted external reserved recipe or alias %s.\n' "$reserved_recipe_name" >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fq "reserved Just recipe or alias outside managed section: $reserved_recipe_name" "$collision_output"; then
      printf 'FAIL: replace-text-section collision error did not name %s.\n' "$reserved_recipe_name" >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$collision_target" "$collision_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting collision %s.\n' "$reserved_recipe_name" >&2
      fail_count=$((fail_count + 1))
    fi
  done

  multiline_header_index=0
  for multiline_header_kind in \
    ordinary-single ordinary-double ordinary-backtick \
    triple-single triple-double triple-backtick; do
    multiline_header_index=$((multiline_header_index + 1))
    multiline_header_payload='payload'
    case "$multiline_header_kind" in
      ordinary-single) multiline_header_delimiter="'" ;;
      ordinary-double) multiline_header_delimiter='"' ;;
      ordinary-backtick) multiline_header_delimiter='`' ;;
      triple-single)
        multiline_header_delimiter="'''"
        multiline_header_payload="payload's"
        ;;
      triple-double)
        multiline_header_delimiter='"""'
        multiline_header_payload='payload"s'
        ;;
      triple-backtick)
        multiline_header_delimiter='```'
        multiline_header_payload='payload`s'
        ;;
    esac
    multiline_header_target="$section_workdir/multiline-header-$multiline_header_index.just"
    printf 'run-kanban value=%s\n%s\n%s:\n    echo collision\n' \
      "$multiline_header_delimiter" "$multiline_header_payload" "$multiline_header_delimiter" \
      > "$multiline_header_target"
    if command -v just >/dev/null 2>&1 \
      && ! just --justfile "$multiline_header_target" --list >/dev/null 2>&1; then
      printf 'FAIL: %s multiline-default recipe-header fixture is not valid Just syntax.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
      continue
    fi
    cp "$multiline_header_target" "$multiline_header_target.before"
    multiline_header_output="$section_workdir/multiline-header-$multiline_header_index.out"
    if "$replace_section_tool" --target "$multiline_header_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$multiline_header_output" 2>&1; then
      printf 'FAIL: replace-text-section accepted reserved recipe with a %s multiline default.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fq \
      'reserved Just recipe or alias outside managed section: run-kanban' \
      "$multiline_header_output"; then
      printf 'FAIL: replace-text-section did not name the reserved recipe with a %s multiline default.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$multiline_header_target" "$multiline_header_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting a %s multiline-default collision.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
    fi
  done

  raw_multiline_target="$section_workdir/raw-multiline-string.just"
  {
    printf '%s\n' \
      'custom-before:' \
      '    echo before' \
      "raw_value := '''" \
      'run-kanban:' \
      'alias kanban-summary := ignored' \
      "\\'''" \
      'custom-after:' \
      '    echo after'
  } > "$raw_multiline_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$raw_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: raw multiline-string collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$raw_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking raw multiline-string content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$raw_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting a raw multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  cooked_multiline_target="$section_workdir/cooked-multiline-string-crlf.just"
  {
    printf '%s\r\n' \
      'custom-before:' \
      '    echo before' \
      'cooked_value := """' \
      'run-kanban-cli:' \
      'odd escaped delimiter remains: \"""' \
      'alias kanban-summary := ignored' \
      'even escaped delimiter closes: \\"""' \
      'joined_value := """closed""" + """open' \
      'run-do-work-update:' \
      '"""' \
      'custom-after:' \
      '    echo after'
  } > "$cooked_multiline_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$cooked_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: cooked multiline-string collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$cooked_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking cooked multiline-string content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$cooked_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting a cooked multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  ordinary_and_command_multiline_target="$section_workdir/ordinary-and-command-multiline-literals.just"
  {
    printf '%s\n' \
      'custom-before:' \
      '    echo before' \
      "raw_value := '" \
      'run-kanban:' \
      'alias kanban-summary := ignored' \
      "raw backslash is literal: \\'" \
      'cooked_value := "' \
      'run-kanban-cli:' \
      'escaped double quote remains: \"' \
      'alias kanban-summary := ignored' \
      'even backslashes close: \\"' \
      'command_value := ```' \
      '  printf "%s\n" safe' \
      'run-do-work-update:' \
      'alias kanban-static := ignored' \
      '```' \
      'custom-after:' \
      '    echo after'
  } > "$ordinary_and_command_multiline_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_and_command_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary-quote and triple-backtick collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$ordinary_and_command_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking ordinary-quote or triple-backtick content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_and_command_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting ordinary-quote and triple-backtick literals.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  ordinary_backtick_multiline_target="$section_workdir/ordinary-backtick-multiline-command.just"
  {
    printf '%s\n' \
      'custom-before:' \
      '    echo before' \
      'command_value := `' \
      'run-kanban:' \
      'alias kanban-summary := ignored' \
      'raw backslash closes: \` + `' \
      'printf safe' \
      '`' \
      'custom-after:' \
      '    echo after'
  } > "$ordinary_backtick_multiline_target"
  cp "$ordinary_backtick_multiline_target" "$ordinary_backtick_multiline_target.before"
  ordinary_backtick_multiline_expected="$section_workdir/ordinary-backtick-multiline-command.expected"
  cp "$ordinary_backtick_multiline_target" "$ordinary_backtick_multiline_expected"
  printf '\n' >> "$ordinary_backtick_multiline_expected"
  cat "$reserved_section_file" >> "$ordinary_backtick_multiline_expected"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_backtick_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary multiline-backtick collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$ordinary_backtick_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking ordinary multiline-backtick content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$ordinary_backtick_multiline_target" "$ordinary_backtick_multiline_expected"; then
    printf 'FAIL: replace-text-section wrote unexpected bytes after accepting an ordinary multiline-backtick command.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_backtick_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting an ordinary multiline-backtick command.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  ordinary_backtick_collision_target="$section_workdir/ordinary-backtick-nearby-collisions.just"
  {
    printf '%s\n' \
      '# unmatched ordinary backtick in comment: `' \
      'one_line_command := `printf safe`' \
      'custom-body:' \
      '    echo unmatched-backtick `' \
      'custom-summary:' \
      '    echo custom' \
      'run-kanban:' \
      '    echo real collision before command' \
      'command_value := `' \
      'kanban-static:' \
      'raw backslash closes: \` + `' \
      'alias kanban-static := ignored' \
      '`' \
      'alias kanban-summary := custom-summary' \
      '@run-kanban-cli:' \
      '    echo real collision after command' \
      'run-do-work-update:' \
      '    echo real collision after inactive forms'
  } > "$ordinary_backtick_collision_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_backtick_collision_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary multiline-backtick nearby-collision control is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  else
    cp "$ordinary_backtick_collision_target" "$ordinary_backtick_collision_target.before"
    ordinary_backtick_collision_output="$section_workdir/ordinary-backtick-nearby-collisions.out"
    expected_ordinary_backtick_collision='replace-text-section: target defines reserved Just recipe or alias outside managed section: kanban-summary, run-do-work-update, run-kanban, run-kanban-cli'
    if "$replace_section_tool" --target "$ordinary_backtick_collision_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$ordinary_backtick_collision_output" 2>&1; then
      printf 'FAIL: replace-text-section ignored real reserved definitions around an ordinary multiline-backtick command.\n' >&2
      fail_count=$((fail_count + 1))
    elif [ "$(cat "$ordinary_backtick_collision_output")" != "$expected_ordinary_backtick_collision" ]; then
      printf 'FAIL: replace-text-section did not report only exact sorted collisions around an ordinary multiline-backtick command.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$ordinary_backtick_collision_target" "$ordinary_backtick_collision_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting collisions around an ordinary multiline-backtick command.\n' >&2
      fail_count=$((fail_count + 1))
    fi
  fi

  ordinary_and_command_collision_target="$section_workdir/ordinary-and-command-nearby-collisions.just"
  {
    printf '%s\n' \
      'custom-summary:' \
      '    echo custom' \
      'run-kanban:' \
      '    echo real collision before raw string' \
      "raw_value := '" \
      'alias kanban-static := ignored' \
      "raw backslash is literal: \\'" \
      'alias kanban-summary := custom-summary' \
      'cooked_value := "' \
      'kanban-static:' \
      'escaped double quote remains: \"' \
      'alias kanban-static := ignored' \
      'even backslashes close: \\"' \
      '@run-kanban-cli:' \
      '    echo real collision before command literal' \
      'command_value := ```' \
      '  printf "%s\n" safe' \
      'kanban-static:' \
      '```' \
      'run-do-work-update:' \
      '    echo real collision after command literal'
  } > "$ordinary_and_command_collision_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_and_command_collision_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary-quote and triple-backtick nearby-collision control is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  else
    cp "$ordinary_and_command_collision_target" "$ordinary_and_command_collision_target.before"
    ordinary_and_command_collision_output="$section_workdir/ordinary-and-command-nearby-collisions.out"
    expected_ordinary_and_command_collision='replace-text-section: target defines reserved Just recipe or alias outside managed section: kanban-summary, run-do-work-update, run-kanban, run-kanban-cli'
    if "$replace_section_tool" --target "$ordinary_and_command_collision_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$ordinary_and_command_collision_output" 2>&1; then
      printf 'FAIL: replace-text-section ignored real reserved definitions around ordinary-quote or triple-backtick literals.\n' >&2
      fail_count=$((fail_count + 1))
    elif [ "$(cat "$ordinary_and_command_collision_output")" != "$expected_ordinary_and_command_collision" ]; then
      printf 'FAIL: replace-text-section did not report only exact sorted collisions around ordinary-quote and triple-backtick literals.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$ordinary_and_command_collision_target" "$ordinary_and_command_collision_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting collisions around ordinary-quote or triple-backtick literals.\n' >&2
      fail_count=$((fail_count + 1))
    fi
  fi

  inactive_literal_forms_target="$section_workdir/inactive-literal-forms-nearby-collision.just"
  {
    printf '%s\n' \
      '# delimiter-looking comment: ```' \
      'one_line_command := `printf safe`' \
      'custom-body:' \
      "    printf '%s\\n' '\`\`\`'" \
      'kanban-static:' \
      '    echo real collision after inactive forms'
  } > "$inactive_literal_forms_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$inactive_literal_forms_target" --list >/dev/null 2>&1; then
    printf 'FAIL: inactive multiline-literal opener control is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  else
    cp "$inactive_literal_forms_target" "$inactive_literal_forms_target.before"
    inactive_literal_forms_output="$section_workdir/inactive-literal-forms-nearby-collision.out"
    expected_inactive_literal_forms_collision='replace-text-section: target defines reserved Just recipe or alias outside managed section: kanban-static'
    if "$replace_section_tool" --target "$inactive_literal_forms_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$inactive_literal_forms_output" 2>&1; then
      printf 'FAIL: replace-text-section let an inactive comment, recipe body, or one-line backtick hide a real collision.\n' >&2
      fail_count=$((fail_count + 1))
    elif [ "$(cat "$inactive_literal_forms_output")" != "$expected_inactive_literal_forms_collision" ]; then
      printf 'FAIL: replace-text-section did not report the exact collision after inactive multiline-literal opener forms.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$inactive_literal_forms_target" "$inactive_literal_forms_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting a collision after inactive literal forms.\n' >&2
      fail_count=$((fail_count + 1))
    fi
  fi

  multiline_collision_target="$section_workdir/multiline-nearby-collisions.just"
  {
    printf '%s\n' \
      '# delimiter-looking comment: """' \
      "single_token := '\"\"\"'" \
      "double_token := \"'''\"" \
      'backtick_token := `printf "\"\"\""`' \
      'custom-shell:' \
      "    echo '\"\"\"'" \
      'custom-summary:' \
      '    echo custom' \
      'alias kanban-summary := custom-summary' \
      'payload := """' \
      'run-kanban:' \
      'alias run-kanban-cli := ignored' \
      'even escaped delimiter closes: \\"""' \
      'run-do-work-update:' \
      '    echo real collision'
  } > "$multiline_collision_target"
  cp "$multiline_collision_target" "$multiline_collision_target.before"
  multiline_collision_output="$section_workdir/multiline-nearby-collisions.out"
  expected_multiline_collision='replace-text-section: target defines reserved Just recipe or alias outside managed section: kanban-summary, run-do-work-update'
  if "$replace_section_tool" --target "$multiline_collision_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions >"$multiline_collision_output" 2>&1; then
    printf 'FAIL: replace-text-section ignored real reserved definitions around a multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  elif [ "$(cat "$multiline_collision_output")" != "$expected_multiline_collision" ]; then
    printf 'FAIL: replace-text-section did not report only exact sorted real collisions around a multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$multiline_collision_target" "$multiline_collision_target.before"; then
    printf 'FAIL: replace-text-section changed the target after rejecting collisions around a multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  noncollision_target="$section_workdir/noncollisions.just"
  {
    printf '# run-kanban:\n'
    printf 'reserved_value := "run-kanban-cli:"\n'
    printf "[doc('kanban-static: is reserved')]\n"
    printf 'custom-recipe: kanban-summary\n'
    printf '    echo run-do-work-update:\n'
    printf 'run-kanban-extra:\n    echo prefix\n'
    printf 'alias custom-summary := kanban-summary\n\n'
    cat "$reserved_section_file"
  } > "$noncollision_target"
  if ! "$replace_section_tool" --target "$noncollision_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated comments, variables, attributes, dependencies, bodies, prefixes, aliases to reserved recipes, or managed definitions as collisions.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  external_collision_target="$section_workdir/external-after-managed.just"
  cat "$reserved_section_file" > "$external_collision_target"
  printf '\nrun-kanban:\n    echo external collision\n' >> "$external_collision_target"
  cp "$external_collision_target" "$external_collision_target.before"
  if "$replace_section_tool" --target "$external_collision_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section ignored a reserved recipe outside an existing managed span.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$external_collision_target" "$external_collision_target.before"; then
    printf 'FAIL: replace-text-section mutated a managed target before rejecting its external collision.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  retired_flag_target="$section_workdir/retired-flag.just"
  printf 'custom-only:\n    echo untouched\n' > "$retired_flag_target"
  cp "$retired_flag_target" "$retired_flag_target.before"
  if "$replace_section_tool" --target "$retired_flag_target" --section-file "$section_file" --migrate-legacy-do-work >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section still accepts the retired legacy-migration flag.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$retired_flag_target" "$retired_flag_target.before"; then
    printf 'FAIL: replace-text-section changed the target after rejecting the retired flag.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  malformed_index=0
  for malformed_content in \
    '# >>> do-work:recipes >>>|one|# >>> do-work:recipes >>>|# <<< do-work:recipes <<<' \
    '# <<< do-work:recipes <<<|one|# >>> do-work:recipes >>>' \
    '# >>> do-work:recipes >>>|one' \
    'one|# <<< do-work:recipes <<<' ; do
    malformed_index=$((malformed_index + 1))
    malformed_target="$section_workdir/malformed-$malformed_index.just"
    printf '%s\n' "$malformed_content" | tr '|' '\n' > "$malformed_target"
    cp "$malformed_target" "$malformed_target.before"
    if "$replace_section_tool" --target "$malformed_target" --section-file "$section_file" >/dev/null 2>&1; then
      printf 'FAIL: replace-text-section accepted malformed marker case %s.\n' "$malformed_index" >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$malformed_target" "$malformed_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting malformed marker case %s.\n' "$malformed_index" >&2
      fail_count=$((fail_count + 1))
    fi
  done

  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$section_workdir/variants/justfile" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced a Justfile that does not parse.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  rm -rf "$section_workdir"
fi

assert_contains \
  "Justfile" \
  '^# >>> do-work:recipes >>>$' \
  'root justfile must open the exact managed do-work recipe section.'
assert_contains \
  "Justfile" \
  '^# <<< do-work:recipes <<<$' \
  'root justfile must close the exact managed do-work recipe section.'
assert_contains \
  "tools/install-do-work-suite.sh" \
  'tools/replace-text-section\.sh|replace-text-section\.sh' \
  'suite installer must reconcile recipes through the managed-section utility.'
assert_contains \
  "Justfile" \
  'skills/do-work-board/tools/queue-kanban' \
  'root Justfile must build the canonical board sibling source.'
assert_contains \
  "Justfile" \
  'skill_root="\$project_root/skills/do-work".*\$skill_root/tools/do-work-update\.sh' \
  'root Justfile fallback must invoke the canonical modular core updater.'
assert_contains \
  "tools/replace-text-section.sh" \
  'suffix=.*dir=parent' \
  'replace-text-section must create its temporary file in the target directory for atomic replacement.'
assert_contains \
  "tools/replace-text-section.sh" \
  'os\.replace\(temporary_path, path\)' \
  'replace-text-section must atomically rename the validated temporary over the target.'

if [ "$fail_count" -gt 0 ]; then
  exit 1
fi

printf 'Contract regression checks passed.\n'
