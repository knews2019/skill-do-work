#!/usr/bin/env bash
# select-simple-reqs-behavior.sh — lock-in suite for the shipped cheaper-model
# selector (skills/do-work/tools/select-simple-reqs.sh) behind
# skills/do-work/actions/run-simple-reqs.md.
#
# Every probe here pins a failure that actually cost something, not a shape:
#
#   T1  The `trivial` alias. When this selector was written, 4 of the 7
#       mechanical REQs in the live queue spelled `effort_estimate: trivial`,
#       not the canonical token — a literal match on `effort-mechanical` alone
#       found 3 and silently dropped the majority. The Schema Read Contract
#       alias map is what makes this selector see them.
#   T2  The dependency-gate bypass. The action hands these ids to
#       `do-work run REQ-NNN ...`, and an explicitly-named REQ bypasses
#       `depends_on` by design (actions/work.md → Input). If the selector does
#       not filter for readiness itself, the whole verb becomes a silent way to
#       run REQs whose prerequisites never landed.
#   T3  The maintenance trap. Instruction-maintenance REQs are always small in
#       diff, so an effort-only filter selects the repository's HIGHEST-judgment
#       work first — 3 of those same 7 live REQs were maintenance REQs. Nothing
#       tests a rule, so nothing downstream would catch a plausible-wrong edit.
#   T4  `tdd: true` is deliberately NOT a veto, because a test-first REQ carries
#       an objective pass/fail gate. Adding the veto looks like a safety
#       improvement and inverts the rule it appears to serve, so it is pinned.
#   T5  An empty selection must exit 0 with a stated answer. A selector that
#       failed here would be indistinguishable from a broken one, and its caller
#       cannot tell "nothing qualifies" from "the scan died" out of stderr.
#   T6  --skip-impact-negligible has to be applied HERE. Forwarded to the
#       handoff it is inert, because explicit REQ naming overrides it
#       downstream — so the flag would appear to work while building exactly
#       the REQs the user asked to omit.
#   T7  `depends_on` wins outright when the legacy `dependencies` key is also
#       present. Merging both lists lets a stale id under the legacy key hold
#       back a REQ whose real dependencies have all landed.
#   T8  A PRESENT but unrecognized enum value must be reported. Silently
#       defaulting a typo'd `effort_estimate` to substantive makes a qualifying
#       REQ vanish from BOTH the selected and the held-back list, leaving the
#       user nothing to diagnose — the exact failure the Schema Read Contract's
#       warn-on-fallback leg exists to prevent.
#   T9  No GNU-only tool flags. `xargs -r` is a GNU extension absent from the
#       BSD xargs macOS ships, and this script claims the stock-macOS floor.
#       Worse than unportable: with the pipeline broken and no `set -e`, the
#       scan would report an empty selection and exit 0.
#  T11  The warn legs for `impact` and `domain` fail in the PROMOTION
#       direction, which T8 does not cover. A typo'd `impact-critcal` or
#       `securty` normalizes to the documented default (`impact-user-visible`,
#       `general`), so neither veto fires and the REQ lands in the run set a
#       cheaper model then builds — the opposite of T8's silent omission and the
#       reason those two warn legs are not narrowable to `effort_estimate`
#       alone. The warning is the only signal either case produces.
#  T12  The compatibility path must delegate selection to `do-work-cli next
#       --simple`; retaining a second parser would let readiness drift again.
#  T13  Duplicate dependency statuses must aggregate conservatively. A
#       last-row-wins index lets filesystem traversal order decide whether an
#       explicitly handed-off dependent bypasses its unresolved prerequisite.
#
# Exit 0: every probe passed. Exit 1: at least one FAIL line above.
set -uo pipefail

script_path="${BASH_SOURCE[0]}"
script_directory="${script_path%/*}"
if [ "$script_directory" = "$script_path" ]; then
  script_directory='.'
fi
repo_root="$(cd "$script_directory/../.." && pwd)"
selector="$repo_root/skills/do-work/tools/select-simple-reqs.sh"

fail_count=0

report_failure() {
  printf 'FAIL: %s\n' "$1" >&2
  fail_count=$((fail_count + 1))
}

if [ ! -x "$selector" ]; then
  report_failure "selector missing or not executable at skills/do-work/tools/select-simple-reqs.sh"
  printf 'select-simple-reqs suite: %s probes failed.\n' "$fail_count" >&2
  exit 1
fi

fixture_root="$(mktemp -d)"
trap 'rm -rf "$fixture_root"' EXIT
mkdir -p "$fixture_root/do-work/queue" "$fixture_root/do-work/archive/UR-001"

write_req() {
  # write_req <relative-path> <frontmatter-body>
  {
    printf -- '---\n'
    printf '%s\n' "$2"
    printf -- '---\n\nFixture body.\n'
  } > "$fixture_root/$1"
}

# An archived, completed dependency — the satisfied side of the T2 pair.
write_req 'do-work/archive/UR-001/REQ-900-done.md' 'id: REQ-900
title: "Landed dependency"
status: completed
domain: general'

# T1 — legacy alias, and a cased/padded variant of it.
write_req 'do-work/queue/REQ-101-legacy-alias.md' 'id: REQ-101
title: "Legacy alias spelling"
status: pending
domain: general
effort_estimate: trivial'
write_req 'do-work/queue/REQ-102-cased-alias.md' 'id: REQ-102
title: "Cased and padded alias"
status: pending
domain: general
effort_estimate:   TRIVIAL  '

# T2 — one REQ waiting on an absent dependency, one whose dependency is archived
# completed. Only the second may be selected.
write_req 'do-work/queue/REQ-103-unmet-dependency.md' 'id: REQ-103
title: "Waits on a REQ that never landed"
status: pending
domain: general
effort_estimate: effort-mechanical
depends_on: [REQ-999]'
write_req 'do-work/queue/REQ-104-met-dependency.md' 'id: REQ-104
title: "Dependency already archived"
status: pending
domain: general
effort_estimate: effort-mechanical
depends_on: [REQ-900]'
# Block-form list, unmet — the inline form above must not be the only shape parsed.
write_req 'do-work/queue/REQ-105-block-form-dependency.md' 'id: REQ-105
title: "Block-form dependency list"
status: pending
domain: general
effort_estimate: effort-mechanical
depends_on:
  - REQ-999'

# T3 — the vetoes.
write_req 'do-work/queue/REQ-106-maintenance.md' 'id: REQ-106
title: "Rule prose, small diff"
status: pending
domain: general
maintenance: true
effort_estimate: effort-mechanical'
write_req 'do-work/queue/REQ-107-security.md' 'id: REQ-107
title: "Small security fix"
status: pending
domain: security
effort_estimate: effort-mechanical'
write_req 'do-work/queue/REQ-108-critical.md' 'id: REQ-108
title: "Small but critical"
status: pending
domain: general
impact: impact-critical
effort_estimate: effort-mechanical'

# T4 — test-first work is a positive signal, never a veto.
write_req 'do-work/queue/REQ-109-tdd.md' 'id: REQ-109
title: "Test-first mechanical fix"
status: pending
domain: general
tdd: true
effort_estimate: effort-mechanical'

# Controls: substantive effort, absent effort, and a non-pending status.
write_req 'do-work/queue/REQ-110-substantive.md' 'id: REQ-110
title: "Real work"
status: pending
domain: general
effort_estimate: effort-substantive'
write_req 'do-work/queue/REQ-111-absent-effort.md' 'id: REQ-111
title: "No effort judgment recorded"
status: pending
domain: general'
write_req 'do-work/queue/REQ-112-claimed.md' 'id: REQ-112
title: "Already in flight"
status: claimed
domain: general
effort_estimate: effort-mechanical
claimed_at: 2026-08-20T10:00:00Z'
write_req 'do-work/queue/REQ-113-assigned.md' 'id: REQ-113
title: "Earmarked for another session"
status: pending
domain: general
effort_estimate: effort-mechanical
assigned_to: "cloud-alpha"'

# T6 — mechanical, ready, and negligible. Selected by default; dropped only when
# the flag is passed to the selector.
write_req 'do-work/queue/REQ-114-negligible.md' 'id: REQ-114
title: "Nobody would notice"
status: pending
domain: general
impact: impact-negligible
effort_estimate: effort-mechanical'

# T7 — canonical depends_on (met) alongside a stale legacy dependencies key.
write_req 'do-work/queue/REQ-115-both-dependency-keys.md' 'id: REQ-115
title: "Corrected REQ that kept the legacy key"
status: pending
domain: general
effort_estimate: effort-mechanical
depends_on: [REQ-900]
dependencies: [REQ-999]'

# T8 — present but unrecognized effort value.
write_req 'do-work/queue/REQ-116-typo-effort.md' 'id: REQ-116
title: "Typo in the effort token"
status: pending
domain: general
effort_estimate: effort-mechanial'

# T11 — typo'd vetoes. Both were MEANT to be held back; both normalize to the
# documented default instead, so the report warns and selects them anyway.
write_req 'do-work/queue/REQ-118-typo-impact.md' 'id: REQ-118
title: "Typo in the impact-critical token"
status: pending
domain: general
impact: impact-critcal
effort_estimate: effort-mechanical'

write_req 'do-work/queue/REQ-119-typo-domain.md' 'id: REQ-119
title: "Typo in the security domain"
status: pending
domain: securty
effort_estimate: effort-mechanical'

# T13 — conflicting duplicate statuses in both traversal orders, all-success
# duplicates, and a unique pending control. The existing REQ-900/REQ-104 pair
# is the unique completed control.
write_req 'do-work/archive/UR-001/REQ-901-a-pending.md' 'id: REQ-901
title: "Duplicate pending first"
status: pending'
write_req 'do-work/archive/UR-001/REQ-901-z-completed.md' 'id: REQ-901
title: "Duplicate completed second"
status: completed'
write_req 'do-work/archive/UR-001/REQ-902-a-completed.md' 'id: REQ-902
title: "Duplicate completed first"
status: completed'
write_req 'do-work/archive/UR-001/REQ-902-z-pending.md' 'id: REQ-902
title: "Duplicate pending second"
status: pending'
write_req 'do-work/archive/UR-001/REQ-903-a-completed.md' 'id: REQ-903
title: "Duplicate success one"
status: completed'
write_req 'do-work/archive/UR-001/REQ-903-z-completed-with-issues.md' 'id: REQ-903
title: "Duplicate success two"
status: completed-with-issues'
write_req 'do-work/archive/UR-001/REQ-904-pending.md' 'id: REQ-904
title: "Unique pending dependency"
status: pending'
for dependency_case in \
  'REQ-120 REQ-901' \
  'REQ-121 REQ-902' \
  'REQ-122 REQ-903' \
  'REQ-123 REQ-904'; do
  dependent_id="${dependency_case%% *}"
  dependency_id="${dependency_case#* }"
  write_req "do-work/queue/$dependent_id-duplicate-status-dependent.md" "id: $dependent_id
title: \"Duplicate status dependent\"
status: pending
domain: general
effort_estimate: effort-mechanical
depends_on: [$dependency_id]"
done

# The report's final `run_set:` line IS the selector's machine-readable
# contract (actions/run-simple-reqs.md reads exactly this line), so the probes
# below read it rather than a second output mode: an ids-only path that exits
# before the warning block would hide every T8 and T11 diagnostic from the very
# assertions meant to catch them.
run_set_ids() {
  bash "$selector" --repo-root "$fixture_root" "$@" 2>/dev/null \
    | sed -n 's/^run_set: //p' | tr ' ' '\n'
}

selected_ids="$(run_set_ids)"

assert_selected() {
  if ! printf '%s\n' "$selected_ids" | grep -qx "$1"; then
    report_failure "$2 — $1 was not selected; selected set was: $(printf '%s' "$selected_ids" | tr '\n' ' ')"
  fi
}

assert_not_selected() {
  if printf '%s\n' "$selected_ids" | grep -qx "$1"; then
    report_failure "$2 — $1 was selected and must not be"
  fi
}

# T1
assert_selected REQ-101 "T1 legacy alias: effort_estimate: trivial must normalize to effort-mechanical (Schema Read Contract)"
assert_selected REQ-102 "T1 legacy alias: a cased, whitespace-padded 'TRIVIAL' must normalize too"

# T2
assert_not_selected REQ-103 "T2 dependency bypass: an unmet depends_on must exclude the REQ, because explicit naming bypasses the gate downstream"
assert_selected     REQ-104 "T2 dependency bypass: a dependency archived as completed must not block selection"
assert_not_selected REQ-105 "T2 dependency bypass: a block-form depends_on list must be parsed, not ignored"

# T3
assert_not_selected REQ-106 "T3 maintenance trap: maintenance: true must be held back — nothing tests a rule"
assert_not_selected REQ-107 "T3 veto: domain: security must be held back"
assert_not_selected REQ-108 "T3 veto: impact-critical must be held back"

# T4
assert_selected REQ-109 "T4 tdd is not a veto: test-first work carries an objective gate and must stay selected"

# T6 — default keeps the negligible REQ; the flag removes it.
assert_selected REQ-114 "T6 default: an impact-negligible REQ is mechanical work and stays selected without the flag"
negligible_filtered_ids="$(run_set_ids --skip-impact-negligible)"
if printf '%s\n' "$negligible_filtered_ids" | grep -qx REQ-114; then
  report_failure "T6 --skip-impact-negligible: REQ-114 must be dropped by the selector, since forwarding the flag to the handoff is inert"
fi
if ! printf '%s\n' "$negligible_filtered_ids" | grep -qx REQ-104; then
  report_failure "T6 --skip-impact-negligible: the flag must remove only negligible REQs, not narrow the set further"
fi
negligible_filtered_report="$(bash "$selector" --repo-root "$fixture_root" --skip-impact-negligible 2>/dev/null)"
if ! printf '%s' "$negligible_filtered_report" | grep -q 'REQ-114'; then
  report_failure "T6 --skip-impact-negligible: the dropped REQ must still be reported with its reason"
fi

# T7 — depends_on wins over the legacy alias.
assert_selected REQ-115 "T7 alias precedence: depends_on wins when both keys are present, so a stale legacy dependencies id must not hold the REQ back"

# T8 — warn on a present-but-unrecognized value, never silently default.
typo_warnings="$(bash "$selector" --repo-root "$fixture_root" 2>&1 >/dev/null)"
if ! printf '%s' "$typo_warnings" | grep -q 'REQ-116'; then
  report_failure "T8 warn-on-fallback: an unrecognized effort_estimate must be reported naming the REQ, not silently defaulted"
fi
if ! printf '%s' "$typo_warnings" | grep -q 'effort-mechanial'; then
  report_failure "T8 warn-on-fallback: the warning must quote the value as written, so the typo is diagnosable"
fi
assert_not_selected REQ-116 "T8 warn-on-fallback: the documented default still applies — an unrecognized value reads as effort-substantive"

# T11 — the promotion direction: a typo'd veto field selects rather than drops,
# so the warning is the whole diagnosis. Narrowing the warn legs to
# `effort_estimate` would make both of these silent.
promotion_warnings="$(bash "$selector" --repo-root "$fixture_root" 2>&1 >/dev/null)"
assert_selected REQ-118 "T11 promotion: a typo'd impact normalizes to impact-user-visible, so the impact-critical veto cannot fire — the documented default, pinned so the warning below is understood as the only signal"
assert_selected REQ-119 "T11 promotion: a typo'd domain normalizes to general, so the security veto cannot fire"
for promotion_case in 'REQ-118 impact impact-critcal' 'REQ-119 domain securty'; do
  set -- $promotion_case
  for promotion_needle in "$1" "$3"; do
    if ! printf '%s' "$promotion_warnings" | grep -q -- "$promotion_needle"; then
      report_failure "T11 promotion: the $2 warn leg must name '$promotion_needle' — with the veto silently not firing, this warning is the only thing standing between a typo and a cheaper model building $1"
    fi
  done
done

# T9 — no GNU-only tool flags, against the stated stock-macOS floor.
selector_code_only="$(sed 's/[[:space:]]*#.*$//' "$selector")"
if printf '%s' "$selector_code_only" | grep -qE 'xargs.*(-r\b|--no-run-if-empty)'; then
  report_failure "T9 portability: xargs -r is a GNU extension absent from the BSD xargs macOS ships; with the pipeline broken and no set -e the scan would report an empty selection and exit 0"
fi

# T12 — one canonical selector, no shell-side request parser.
if ! printf '%s' "$selector_code_only" | grep -q 'do-work-cli.sh.*next --simple'; then
  report_failure "T12 canonical delegation: the compatibility script must invoke do-work-cli next --simple"
fi
if printf '%s' "$selector_code_only" | grep -qE 'find .*REQ-|function normalize_|known_status\['; then
  report_failure "T12 canonical delegation: the compatibility script must not retain its own queue parser or readiness index"
fi

# T13 — conflicting rows stay held in either order, while every-success and
# unique completed dependencies are ready. A unique pending dependency stays held.
assert_not_selected REQ-120 "T13 duplicate status order: pending then completed must remain unresolved"
assert_not_selected REQ-121 "T13 duplicate status order: completed then pending must remain unresolved"
assert_selected REQ-122 "T13 duplicate status aggregation: every successful copy satisfies the dependency"
assert_not_selected REQ-123 "T13 unique pending dependency remains unresolved"

# Controls
assert_not_selected REQ-110 "control: effort-substantive must never be selected"
assert_not_selected REQ-111 "control: an absent effort_estimate defaults to effort-substantive and must not be selected"
assert_not_selected REQ-112 "control: a claimed REQ must never be selected"
assert_not_selected REQ-113 "control: a REQ assigned to another session must not be selected"

# The held-back REQs must be REPORTED, not silently dropped: the reason is what
# lets a user fix a mis-tagged REQ instead of wondering where it went.
report_output="$(bash "$selector" --repo-root "$fixture_root" 2>/dev/null)"
for held_back_id in REQ-103 REQ-106 REQ-107 REQ-108; do
  if ! printf '%s' "$report_output" | grep -q "$held_back_id"; then
    report_failure "held-back $held_back_id must appear in the report with its reason, not be dropped silently"
  fi
done
if ! printf '%s' "$report_output" | grep -q '^run_set: '; then
  report_failure "report must end with a run_set: line — it is the contract actions/run-simple-reqs.md reads"
fi

# T5 — an empty selection is a normal answer.
empty_root="$(mktemp -d)"
mkdir -p "$empty_root/do-work/queue"
write_req_empty_case() {
  {
    printf -- '---\n'
    printf 'id: REQ-201\ntitle: "Real work only"\nstatus: pending\ndomain: general\neffort_estimate: effort-substantive\n'
    printf -- '---\n\nBody.\n'
  } > "$empty_root/do-work/queue/REQ-201-substantive.md"
}
write_req_empty_case
if empty_output="$(bash "$selector" --repo-root "$empty_root")"; then
  if ! printf '%s' "$empty_output" | grep -q 'No pending REQ currently qualifies'; then
    report_failure "T5 empty selection: must state that nothing qualifies"
  fi
  if ! printf '%s' "$empty_output" | grep -qx 'run_set: '; then
    report_failure "T5 empty selection: run_set: line must still be emitted, and empty"
  fi
else
  report_failure "T5 empty selection: must exit 0 — an empty queue of mechanical work is an answer, not an error"
fi
rm -rf "$empty_root"

# A missing queue is a usage error, distinct from an empty one.
missing_root="$(mktemp -d)"
missing_queue_status=0
bash "$selector" --repo-root "$missing_root" >/dev/null 2>&1 || missing_queue_status=$?
if [ "$missing_queue_status" -ne 2 ]; then
  report_failure "a missing do-work/queue must exit 2 (got $missing_queue_status), so a caller can tell it apart from an empty selection"
fi
rm -rf "$missing_root"

if [ "$fail_count" -ne 0 ]; then
  printf 'select-simple-reqs suite: %s probes failed.\n' "$fail_count" >&2
  exit 1
fi

printf 'select-simple-reqs suite: all probes passed.\n'
