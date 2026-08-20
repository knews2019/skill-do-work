#!/usr/bin/env bash
# select-simple-reqs.sh — list the pending REQs a cheaper model can be trusted
# with, plus a P50 estimate for the batch. Backs actions/run-simple-reqs.md.
#
# Usage: tools/select-simple-reqs.sh [--repo-root DIR] [--skip-impact-negligible]
#                                    [--ids-only]
# Output (default): a human report — one row per selected REQ with its p50, the
#         dropped REQs with the reason each was dropped, any Schema Read
#         Contract warnings, then the batch totals. The final line is always
#         "run_set: <ids>" (empty when none qualify), which is the line
#         actions/run-simple-reqs.md reads to build its handoff command.
#        (--ids-only): the selected ids, one per line, and nothing else.
# Exit 0: report emitted, INCLUDING when nothing qualifies — an empty queue of
#         mechanical work is a normal answer, not an error, and the caller must
#         be able to tell the two apart without parsing stderr.
# Exit 2: usage error, no do-work/queue directory under --repo-root, or a scan
#         that failed outright. A failed scan must NEVER be reported as an empty
#         selection: "nothing qualifies" and "the scan died" are different
#         answers and the caller acts differently on each.
#
# THE PREDICATE IS ONE CONDITION: is there an objective gate on the result?
# A REQ is selected when it is pending, dependency-ready, unclaimed, unassigned,
# and its `effort_estimate` normalizes to `effort-mechanical` — and it is dropped
# when nothing objectively gates the result, or when the cost of a miss is
# unbounded. The markers below are the current such cases, illustrative and not
# a closed list (_dev/primes/prime-shell-commands.md § Closed Enumerations Go
# Stale):
#   maintenance: true  — edits the skill's own rule prose; nothing tests a rule.
#   domain: security   — a gate exists, but the cost of a miss is unbounded.
#   impact-critical    — same.
#
# `standing: true` is excluded for a different reason, and it is not a judgment
# call: actions/work.md Step 1 states that the default scan NEVER selects a
# standing sweep, because batching instances is the whole economy of one. It
# drains only when a human names it or opportunistically. Since this selector
# hands its set to `do-work run REQ-NNN ...`, auto-selecting one would name it
# explicitly — the exact drain trigger — and turn a batching REQ into an
# unrequested drain. Same lesson as --skip-impact-negligible below: a rule that
# must survive the handoff has to be applied before it.
#
# `tdd: true` is DELIBERATELY NOT a veto, and re-adding one is a regression
# (_dev/tests/select-simple-reqs-behavior.sh pins it). A TDD REQ carries an
# objective pass/fail gate, often with a captured RED case in its
# `## Red-Green Proof` section — a STRONGER gate than the qualification-plus-
# review a non-TDD REQ gets. By the condition above that makes test-first work a
# positive signal, so excluding it would invert the rule it looks like it serves.
#
# Dependency-readiness is not optional here. The caller hands these ids to
# `do-work run REQ-NNN ...`, and an explicitly-named REQ bypasses `depends_on`
# by design (actions/work.md → Input). A selector that skipped this check would
# turn its caller into a silent dependency-gate bypass.
#
# `--skip-impact-negligible` is applied HERE for the same reason, not forwarded
# by the caller: explicit naming also overrides that flag downstream, so a
# forwarded copy would be silently inert and build the very work the user asked
# to omit. Every filter that must survive the handoff has to run before it.
#
# This script implements actions/work.md Step 1's selection scan (pending,
# unclaimed, unassigned, dependency-ready) plus the two filters above. Step 1
# stays canonical for that predicate; this file changes when it does.
#
# Field normalization follows the Schema Read Contract
# (actions/work-reference.md), including its read-only legacy aliases —
# `trivial` for effort-mechanical is not optional trivia: REQs written before
# the rename still carry it, and a literal match on the canonical token alone
# silently drops them. It also honors the contract's warn-on-fallback leg: a
# PRESENT but unrecognized enum value is reported, never silently defaulted,
# because a typo'd `effort_estimate` would otherwise make a qualifying REQ
# vanish from both the selected and the held-back list with nothing to diagnose.
#
# Bash 3.2 compatible, because stock macOS /bin/bash still is and the scripts in
# tools/ hold that floor. That rules out `mapfile` and `declare -A`, so the
# frontmatter parse and the predicate are resolved in awk — which has had
# associative arrays all along — rather than in shell. It also rules out GNU-only
# tool flags: the scans use `find -exec ... {} +` rather than `xargs -r`, whose
# `-r` is a GNU extension absent from the BSD xargs macOS ships. Determinism
# comes from sorting awk's output, not from sorting find's input.
#
# Read-only. This script selects and reports; claiming, running, and every
# queue write stay with the caller.
set -uo pipefail

repository_root="."
ids_only="false"
skip_impact_negligible="false"

while [ $# -gt 0 ]; do
  case "$1" in
    --repo-root)
      if [ $# -lt 2 ]; then
        printf 'select-simple-reqs: --repo-root needs a directory\n' >&2
        exit 2
      fi
      repository_root="$2"
      shift 2
      ;;
    --ids-only)
      ids_only="true"
      shift
      ;;
    --skip-impact-negligible)
      skip_impact_negligible="true"
      shift
      ;;
    -h|--help)
      sed -n '2,/^set -uo/p' "$0" | sed 's/^# \{0,1\}//; $d'
      exit 0
      ;;
    *)
      printf 'select-simple-reqs: unrecognized argument %s\n' "$1" >&2
      printf 'usage: select-simple-reqs.sh [--repo-root DIR] [--skip-impact-negligible] [--ids-only]\n' >&2
      exit 2
      ;;
  esac
done

work_root="$repository_root/do-work"
queue_directory="$work_root/queue"

if [ ! -d "$queue_directory" ]; then
  printf 'select-simple-reqs: no queue at %s\n' "$queue_directory" >&2
  exit 2
fi

script_path="${BASH_SOURCE[0]}"
script_directory="${script_path%/*}"
if [ "$script_directory" = "$script_path" ]; then
  script_directory='.'
fi
estimator="$script_directory/estimate-p50.sh"

scratch_directory="$(mktemp -d 2>/dev/null || mktemp -d -t selectsimple)"
if [ -z "$scratch_directory" ] || [ ! -d "$scratch_directory" ]; then
  printf 'select-simple-reqs: could not allocate a scratch directory\n' >&2
  exit 2
fi
trap 'rm -rf "$scratch_directory"' EXIT

status_index="$scratch_directory/status-index.tsv"
selection_rows="$scratch_directory/selection.tsv"
raw_rows="$scratch_directory/selection-unsorted.tsv"

# The status index answers "is this dependency done?" for every REQ the
# repository knows about — queued, in flight, or archived. One awk over every
# file, so the predicate below is a lookup rather than a per-dependency hunt.
if ! find "$work_root" -name 'REQ-*.md' -type f -exec awk '
      function trim(value) { sub(/^[ \t]+/, "", value); sub(/[ \t]+$/, "", value); return value }
      function flush() {
        if (identifier == "") { return }
        resolved = state
        if (resolved == "complete" || resolved == "done" || resolved == "finished" || resolved == "closed") {
          resolved = "completed"
        }
        printf "%s\t%s\n", identifier, resolved
        identifier = ""
      }
      FNR == 1 { flush(); delimiters = 0; identifier = ""; state = "" }
      /^---[ \t]*$/ { delimiters++; next }
      delimiters == 1 {
        if ($0 ~ /^id:/)          { identifier = trim(substr($0, index($0, ":") + 1)) }
        else if ($0 ~ /^status:/) { state = tolower(trim(substr($0, index($0, ":") + 1))) }
      }
      END { flush() }
    ' {} + > "$status_index" 2>/dev/null; then
  printf 'select-simple-reqs: could not scan %s for REQ statuses\n' "$work_root" >&2
  exit 2
fi

# One awk pass over the queue applies the whole predicate. Emitting the drops
# with their reason is deliberate: a selector that silently narrows the queue is
# indistinguishable from an empty queue, and `--skip-impact-negligible` already
# set the precedent of reporting what it dropped (actions/work.md → Input).
if ! find "$queue_directory" -maxdepth 1 -name '*.md' -type f -exec awk \
    -v index_file="$status_index" \
    -v skip_negligible="$skip_impact_negligible" '
    function trim(value) { sub(/^[ \t]+/, "", value); sub(/[ \t]+$/, "", value); return value }
    function value_of(line) { return trim(substr(line, index(line, ":") + 1)) }
    function truthy(value) {
      value = tolower(trim(value))
      return (value == "true" || value == "yes" || value == "on" || value == "t")
    }
    function normalize_effort(value) {
      value = tolower(trim(value))
      if (value == "trivial" || value == "effort-mechanical") { return "effort-mechanical" }
      return "effort-substantive"
    }
    function normalize_impact(value) {
      value = tolower(trim(value))
      if (value == "impact-critical" || value == "impact-user-visible" \
          || value == "impact-rule-change" || value == "impact-negligible") { return value }
      return "impact-user-visible"
    }
    function normalize_domain(value) {
      value = tolower(trim(value))
      if (value == "back-end" || value == "back_end") { return "backend" }
      if (value == "front-end" || value == "front_end") { return "frontend" }
      if (value == "ui_design") { return "ui-design" }
      if (value == "sec") { return "security" }
      if (value == "test") { return "testing" }
      if (value == "content-management" || value == "content_management") { return "cms" }
      return value
    }
    # The Schema Read Contract requires a PRESENT but unrecognized enum value to
    # be reported rather than silently defaulted. Absence is not a fallback and
    # never warns — every one of these fields has a documented absent-reads-as.
    function effort_recognized(value) {
      value = tolower(trim(value))
      return (value == "" || value == "trivial" || value == "normal" \
              || value == "effort-mechanical" || value == "effort-substantive")
    }
    function impact_recognized(value) {
      value = tolower(trim(value))
      return (value == "" || normalize_impact(value) == value)
    }
    function domain_recognized(value) {
      resolved = normalize_domain(value)
      return (resolved == "" || resolved == "frontend" || resolved == "backend" \
              || resolved == "ui-design" || resolved == "general" || resolved == "security" \
              || resolved == "testing" || resolved == "cms")
    }
    function boolean_recognized(value) {
      value = tolower(trim(value))
      return (value == "" || value == "true" || value == "false" || value == "yes" \
              || value == "on" || value == "t" || value == "no" || value == "off" || value == "f")
    }
    function status_recognized(value) {
      value = tolower(trim(value))
      return (value == "" || value == "pending" || value == "claimed" || value == "completed" \
              || value == "completed-with-issues" || value == "failed" || value == "cancelled" \
              || value == "pending-answers" || value == "blocked" \
              || value == "blocked-archive-collision" || value == "blocked-dependency-cycle" \
              || value == "complete" || value == "done" || value == "finished" || value == "closed" \
              || value == "canceled" || value == "abandoned" || value == "wont-do" || value == "wontfix")
    }
    function emit() {
      if (identifier == "") { return }

      if (!effort_recognized(effort))        { printf "WARN\t%s\teffort_estimate\t%s\teffort-substantive\n", identifier, effort }
      if (!impact_recognized(impact))        { printf "WARN\t%s\timpact\t%s\timpact-user-visible\n", identifier, impact }
      if (!domain_recognized(domain))        { printf "WARN\t%s\tdomain\t%s\tgeneral\n", identifier, domain }
      if (!boolean_recognized(maintenance))  { printf "WARN\t%s\tmaintenance\t%s\tfalse\n", identifier, maintenance }
      if (!status_recognized(state))         { printf "WARN\t%s\tstatus\t%s\tskipped\n", identifier, state }

      # `depends_on` wins outright when both keys are present — the canonical
      # canonical contract rule (actions/work.md → Input). Merging both lists would
      # let a stale id under the legacy key hold back an already-ready REQ.
      if (depends_on_seen) { dependency_count = canonical_count; for (i = 1; i <= canonical_count; i++) { dependencies[i] = canonical[i] } }
      else                 { dependency_count = legacy_count;    for (i = 1; i <= legacy_count;    i++) { dependencies[i] = legacy[i] } }

      reason = ""
      if (state != "pending")                                   { reason = "status " state }
      else if (normalize_effort(effort) != "effort-mechanical")  { reason = "not mechanical" }
      else if (claimed != "")                                   { reason = "already claimed" }
      else if (assigned != "")                                  { reason = "assigned to " assigned }
      else if (truthy(standing))                                { reason = "standing sweep: drains only when named or opportunistically" }
      else if (truthy(maintenance))                             { reason = "maintenance: rule prose has no test" }
      else if (normalize_domain(domain) == "security")           { reason = "security: cost of a miss is unbounded" }
      else if (normalize_impact(impact) == "impact-critical")     { reason = "impact-critical" }
      else if (skip_negligible == "true" && normalize_impact(impact) == "impact-negligible") {
        reason = "impact-negligible (--skip-impact-negligible)"
      }
      else {
        for (i = 1; i <= dependency_count; i++) {
          dependency_state = known_status[dependencies[i]]
          if (dependency_state != "completed" && dependency_state != "completed-with-issues") {
            reason = "waits on " dependencies[i]
            break
          }
        }
      }
      # Only "not mechanical" and a non-pending status are silent: every REQ in
      # the queue trips one of those, so reporting them would bury the drops the
      # user can actually act on.
      if (reason == "") { printf "SELECT\t%s\t%s\t%s\n", identifier, (frozen_p50 == "" ? "-" : frozen_p50), title }
      else if (reason !~ /^(status |not mechanical)/) { printf "DROP\t%s\t%s\t%s\n", identifier, reason, title }
    }
    BEGIN {
      FS = "\n"
      while ((getline line < index_file) > 0) {
        split(line, parts, "\t")
        known_status[parts[1]] = parts[2]
      }
      close(index_file)
    }
    FNR == 1 {
      emit()
      delimiters = 0; identifier = ""; state = ""; effort = ""; impact = ""
      maintenance = ""; domain = ""; title = ""; claimed = ""; assigned = ""; standing = ""
      frozen_p50 = ""; dependency_count = 0; collecting = ""
      canonical_count = 0; legacy_count = 0; depends_on_seen = 0
      delete dependencies; delete canonical; delete legacy
    }
    /^---[ \t]*$/ { delimiters++; next }
    delimiters != 1 { next }
    {
      # A block-form list continues until the next key, so the collector has to
      # be cleared by any key line rather than only by a recognized one.
      if (collecting != "" && $0 ~ /^[ \t]+-[ \t]*/) {
        entry = $0
        sub(/^[ \t]+-[ \t]*/, "", entry)
        entry = trim(entry); gsub(/^"|"$/, "", entry)
        if (entry != "") {
          if (collecting == "canonical") { canonical[++canonical_count] = entry }
          else                           { legacy[++legacy_count] = entry }
        }
        next
      }
      if ($0 ~ /^[A-Za-z_][A-Za-z0-9_]*:/) { collecting = "" }

      if ($0 ~ /^id:/)              { identifier  = value_of($0) }
      else if ($0 ~ /^status:/)     { state       = tolower(value_of($0)) }
      else if ($0 ~ /^effort_estimate:/) { effort = value_of($0) }
      else if ($0 ~ /^impact:/)     { impact      = value_of($0) }
      else if ($0 ~ /^maintenance:/){ maintenance = value_of($0) }
      else if ($0 ~ /^standing:/)   { standing    = value_of($0) }
      else if ($0 ~ /^domain:/)     { domain      = value_of($0) }
      else if ($0 ~ /^claimed_at:/) { claimed     = value_of($0) }
      else if ($0 ~ /^assigned_to:/){ assigned    = value_of($0); gsub(/^"|"$/, "", assigned) }
      else if ($0 ~ /^title:/)      { title       = value_of($0); gsub(/^"|"$/, "", title) }
      else if ($0 ~ /p50_active_minutes:/) { frozen_p50 = value_of($0) }
      else if ($0 ~ /^depends_on:/ || $0 ~ /^dependencies:/) {
        which = ($0 ~ /^depends_on:/) ? "canonical" : "legacy"
        if (which == "canonical") { depends_on_seen = 1 }
        raw = value_of($0)
        if (raw == "") { collecting = which }
        else if (raw ~ /^\[/) {
          gsub(/^\[|\]$/, "", raw)
          count = split(raw, entries, ",")
          for (i = 1; i <= count; i++) {
            entry = trim(entries[i]); gsub(/^"|"$/, "", entry)
            if (entry != "") {
              if (which == "canonical") { canonical[++canonical_count] = entry }
              else                      { legacy[++legacy_count] = entry }
            }
          }
        }
      }
    }
    END { emit() }
  ' {} + > "$raw_rows" 2>/dev/null; then
  printf 'select-simple-reqs: could not scan %s\n' "$queue_directory" >&2
  exit 2
fi

# Determinism comes from sorting the rows, not the scan order: `find -exec ... +`
# hands files over in filesystem order, which differs between checkouts.
LC_ALL=C sort "$raw_rows" > "$selection_rows"

selected_ids=""
estimator_arguments=""
selected_count=0

# Every selected REQ is effort-mechanical, so the estimator's floor mode answers
# for all of them and answers identically — resolve it once, never per REQ, and
# never by signal extraction (actions/estimate-reference.md § The
# Mechanical-Effort Short-Circuit).
floor_minutes=""
if [ -x "$estimator" ]; then
  floor_minutes="$(bash "$estimator" --trivial 2>/dev/null | awk '/^p50_active_minutes:/ { print $2 }')"
fi
[ -n "$floor_minutes" ] || floor_minutes="5"

: > "$scratch_directory/selected.tsv"
while IFS="$(printf '\t')" read -r kind identifier detail title; do
  [ "$kind" = "SELECT" ] || continue
  minutes="$detail"
  if [ "$minutes" = "-" ]; then
    minutes="$floor_minutes"
  fi
  selected_count=$((selected_count + 1))
  selected_ids="$selected_ids $identifier"
  # No dependency edges are passed because none can exist inside the set: a
  # selected REQ's dependencies are all terminal, and a terminal REQ is never
  # itself pending, so no selected REQ can depend on another.
  estimator_arguments="$estimator_arguments $identifier:$minutes"
  printf '%s\t%s\t%s\n' "$identifier" "$minutes" "$title" >> "$scratch_directory/selected.tsv"
done < "$selection_rows"

selected_ids="${selected_ids# }"

if [ "$ids_only" = "true" ]; then
  # shellcheck disable=SC2086
  # Word splitting is intended: selected_ids is a space-separated id list.
  for identifier in $selected_ids; do
    printf '%s\n' "$identifier"
  done
  exit 0
fi

if [ "$selected_count" -eq 0 ]; then
  printf 'No pending REQ currently qualifies for a cheaper model.\n'
else
  printf 'Mechanical REQs a cheaper model can take (%s):\n\n' "$selected_count"
  while IFS="$(printf '\t')" read -r identifier minutes title; do
    printf '  %-10s %4s min  %s\n' "$identifier" "$minutes" "$title"
  done < "$scratch_directory/selected.tsv"
  printf '\n'
  if [ -x "$estimator" ]; then
    # shellcheck disable=SC2086
    # Word splitting is the interface: each REQ is a separate ID:MINUTES argument
    # (actions/estimate-reference.md § Multi-REQ Totals and Critical Path).
    bash "$estimator" critical-path $estimator_arguments 2>/dev/null \
      | awk -F': ' '
          /^total_estimated_effort_minutes:/ { printf "  Total estimated effort: %s active minutes\n", $2 }
          /^critical_path_minutes:/          { printf "  Estimated critical path: %s active minutes\n", $2 }
        '
    printf '\n'
  fi
fi

if grep -q '^DROP' "$selection_rows" 2>/dev/null; then
  printf 'Mechanical but held back:\n\n'
  while IFS="$(printf '\t')" read -r kind identifier detail title; do
    [ "$kind" = "DROP" ] || continue
    printf '  %-10s %s\n' "$identifier" "$detail"
  done < "$selection_rows"
  printf '\n'
fi

# The Schema Read Contract's warn-on-fallback leg, in its prescribed wording.
if grep -q '^WARN' "$selection_rows" 2>/dev/null; then
  while IFS="$(printf '\t')" read -r kind identifier field written_value fallback; do
    [ "$kind" = "WARN" ] || continue
    printf '⚠ %s %s: %s not recognized — treating as %s.\n' \
      "$identifier" "$field" "'$written_value'" "'$fallback'" >&2
  done < "$selection_rows"
  printf '\n' >&2
fi

printf 'run_set: %s\n' "$selected_ids"
