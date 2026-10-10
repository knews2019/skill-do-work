#!/usr/bin/env bash
# GREEN checks for REQ-683 (integrator test gate): the existing cadence table test passes with the AM and PM rows ("daily 5:00 PM" -> 17:00 and "weekly Friday 12:30 AM" -> 00:30 must be present as passing rows). Fails before the build.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
output="$(go test -C "$root/skills/do-work/tools/do-work-cli" -count=1 -v -run '^TestInterviewStandingSlotsParseOnlyRecurringTiming$' ./internal/knowledgecommands/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for row_name in 'daily_5:00_PM' 'weekly_Friday_12:30_AM'; do
  grep -q -- "--- PASS: TestInterviewStandingSlotsParseOnlyRecurringTiming/$row_name " <<<"$output" || { printf '%s\n' "$output"; echo "table row $row_name is missing or did not pass"; exit 1; }
done
