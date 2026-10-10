#!/usr/bin/env bash
# shellcheck disable=SC2016  # backticks in the patterns are literal Markdown, not expansions
# GREEN checks for REQ-689 (integrator test gate): the --coordinate flag, its route row, the run rules,
# the preflight, the stall loop, the coordinated resume command and the teardown additions exist in the
# shipped prose, plain --fan-out is unchanged, fan-out-reference.md has no pid, and the two shipped-prose
# contract tests still pass. Reads paths relative to the tree it runs in. Fails before the build.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
actions_dir="$root/skills/do-work/actions"
work_text="$(<"$actions_dir/work.md")"
fan_out_text="$(<"$actions_dir/fan-out-reference.md")"
handoff_text="$(<"$actions_dir/restart-with-parallel-handoff.md")"
skill_text="$(<"$root/skills/do-work/SKILL.md")"
failures=0
fail() { echo "FAIL: $1"; failures=$((failures + 1)); }

# A1: the flag parses (Input entry, usage line, Step 0 checklist line); plain --fan-out stays as it was.
grep -q -- '^- \*\*`--coordinate`' <<<"$work_text" || fail "work.md ## Input has no --coordinate entry"
usage_line="$(grep -- '^Unrecognized argument(s):' <<<"$work_text" || true)"
grep -q -- '--coordinate' <<<"$usage_line" || fail "work.md usage line lacks --coordinate"
grep -q -- '\[--fan-out \[N\]\]' <<<"$usage_line" || fail "work.md usage line lost [--fan-out [N]]"
checklist_line="$(grep -- 'Step 0: Parse arguments' <<<"$work_text" || true)"
grep -q -- '--coordinate' <<<"$checklist_line" || fail "work.md Step 0 checklist line lacks --coordinate"
grep -q -- 'degrades silently to the serial loop' <<<"$work_text" || fail "work.md lost the plain --fan-out silent degrade"

# A2: the route row sits above the plain run row and points at work.md.
coordinate_row_line="$(grep -n -- 'drive the queue' <<<"$skill_text" | head -n 1 | cut -d: -f1 || true)"
run_row_line="$(grep -n -- '^| `run`, `go`' <<<"$skill_text" | head -n 1 | cut -d: -f1 || true)"
if [ -z "$coordinate_row_line" ] || [ -z "$run_row_line" ]; then
  fail "SKILL.md: coordinate route row or plain run row not found"
elif [ "$coordinate_row_line" -ge "$run_row_line" ]; then
  fail "SKILL.md: coordinate route row is not above the plain run row"
fi
coordinate_row="$(grep -- 'drive the queue' <<<"$skill_text" | head -n 1 || true)"
grep -q -- './actions/work.md' <<<"$coordinate_row" || fail "SKILL.md coordinate row does not route to ./actions/work.md"
grep -q -- 'use the main session as a coordinator' <<<"$coordinate_row" || fail "SKILL.md coordinate row lacks the coordinator phrase"

# A3 and B: mandatory shape, run rules, run-directory rows, manifest columns.
grep -q -- 'Under `--coordinate` this shape is mandatory' <<<"$fan_out_text" || fail "fan-out-reference.md lacks the mandatory-shape sentence"
grep -q -- 'pending-answers' <<<"$fan_out_text" || fail "fan-out-reference.md lacks the pending-answers rule"
grep -q -- '^| .*REQ-NNN-progress\.log' <<<"$fan_out_text" || fail "run-directory table lacks the REQ-NNN-progress.log row"
grep -q -- '^| .*full-gate\.lock' <<<"$fan_out_text" || fail "run-directory table lacks the full-gate.lock row"
manifest_row="$(grep -- '^| `manifest.md`' <<<"$fan_out_text" || true)"
for column_name in 'takeover-to-finalization minutes' 'full gates run' 'stall restarts'; do
  grep -q -- "$column_name" <<<"$manifest_row" || fail "manifest.md row lacks the column: $column_name"
done
grep -q -- 'REQ-069' <<<"$fan_out_text" || fail "run rules do not name the REQ-069/REQ-073 boundary"

# C and D: preflight, run policy, stall loop.
grep -q -- 'low-disk-space' <<<"$fan_out_text" || fail "preflight lacks the low-disk-space level"
grep -q -- 'do-work/run-policy.md' <<<"$fan_out_text" || fail "fan-out-reference.md lacks do-work/run-policy.md"
grep -q -- 'do-work status --watch' <<<"$fan_out_text" || fail "stall loop does not read do-work status --watch"

# E: coordinated resume command and teardown.
grep -q -- 'do-work run --coordinate --fan-out N' <<<"$handoff_text" || fail "handoff resume command lacks do-work run --coordinate --fan-out N"
grep -q -- 'do-work run --fan-out N' <<<"$handoff_text" || fail "handoff lost the plain do-work run --fan-out N resume command"
teardown_text="$fan_out_text
$work_text"
grep -qi -- 'idle background agents' <<<"$teardown_text" || fail "teardown does not stop idle background agents"

# No liveness check by pid in the reference.
if grep -qiw -- 'pid' <<<"$fan_out_text"; then fail "fan-out-reference.md mentions pid"; fi

if [ "$failures" -gt 0 ]; then echo "$failures check(s) failed"; exit 1; fi

bash "$root/_dev/tests/shipped-package-reference-contract.sh"
bash "$root/_dev/tests/action-shell-blocks.sh"
echo "REQ-689 GREEN probe: all checks passed"
