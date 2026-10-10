# Hand-back: REQ-683 (interview cadence reads AM and PM)

- Branch: worktree-agent-REQ-683-interview-cadence-am-pm
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-683-interview-cadence-am-pm
- Base: b629e5cd
- Commit: b83101b1 `[REQ-683] read AM and PM in interview cadence clocks` (one commit)

## File manifest
- skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_derivations.go (modified): `interviewClockPattern` now case-insensitive with an optional `am`/`pm` group (space optional); `parseInterviewCadence` converts to 24-hour before `interviewValidClock`; added `fmt` import.
- skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_export_regression_test.go (modified): four rows in the existing `TestInterviewStandingSlotsParseOnlyRecurringTiming` table.

## P-A-U
- [PLAN]: Extend the clock regex with `(?: ?([ap]m))?`, take the submatch, pad as before, convert when a marker exists (reject hour > 12, hour %= 12, add 12 for PM), keep `interviewValidClock` last. Four table rows, no new test.
- [APPLY]: Rows added first and run at base (RED), then parser changed (GREEN), exactly as planned. Only the two named files touched.
- [UNIFY]: `git diff b629e5cd --stat`: 2 files changed, 23 insertions(+), 2 deletions(-) (interview_derivations.go 21 lines, interview_export_regression_test.go 4 lines). Checks, all exit 0: `go test -count=1 ./internal/knowledgecommands/` (whole package, 11.1s wall), probe REQ-683-probe.sh (0.5s), `gofmt -l` (empty), `go vet ./...` (0.65s, clean), `git diff --check`. Files checked: both changed files, no debug artifacts.

## Red-green record
Base (rows added, parser unchanged), `-run TestInterviewStandingSlotsParseOnlyRecurringTiming`:
- daily_5:00_PM: got "daily" "rolling" "05:00" true (want 17:00)
- weekly_Friday_12:30_AM: got "weekly" "Fri" "12:30" true (want 00:30)
- daily_13:00_PM: got "daily" "rolling" "13:00" true (want not ok)
- daily_12:00_PM passed at base (12:00 was already right by luck; it pins the 12 PM rule against a naive "+12" fix).
After the change: all 12 rows PASS, including the 8 existing ones.

## Decisions
- D-01 DECIDE & STATE: the 12:00 PM and 13:00 PM rows are written as "daily 12:00 PM" and "daily 13:00 PM", because a bare clock has no cadence word and returns not-ok today (as the brief said).
- D-02 DECIDE & STATE: a marker with hour > 12 returns not-ok directly inside the conversion. Without this, "13:00 PM" would convert to 13:00 and pass `interviewValidClock`. `interviewValidClock` stays the final check for everything else.
- D-03 DECIDE & STATE: used `fmt.Sscanf` and `fmt.Sprintf` (stdlib, no time library) to read and write the hour and minute; the `Sscanf` result is unchecked because the regex already guarantees `H:MM` digits.
- D-04 DECIDE & STATE: "0:30 AM" is accepted as 00:30 (only hour > 12 is rejected). Not pinned by a test (YAGNI).

## Discovered Tasks
- "5pm" without a colon still yields no time, per the REQ; the spec text (`work-operating-model.md`) does not say that AM/PM is now read. Optional doc mention → report only.

## Lessons read
`_dev/primes/lessons-releases.md` (whole file; no family touches this change). Go satellite `lessons-do-work-cli.md` not read (dropped for budget; no bullet-named code touched).

## Anti-bloat check
`git diff --stat`: 2 files, +23/-2. Added beyond the REQ: nothing. No new function, constant, option, file or test. New local variable `meridiem` and the `fmt` import are part of the one conversion.

## Proposed CHANGELOG entry
Interview cadence reads AM and PM in clock answers. "daily 5:00 PM" was saved as 05:00 and "weekly Friday 12:30 AM" as 12:30, so standing slots were off by 12 hours.
- `parseInterviewCadence` accepts an optional am/pm after the clock (any case, space optional) and converts to 24-hour time: 12 AM is 00, 12 PM is 12, other PM hours add 12.
- Invalid mixes such as "13:00 PM" are rejected; "5pm" without a colon still gives no time.
- Four rows added to the cadence table test.

## Proposed lesson bullet
none

## Integration seams and timings
None. Test wall times: package 11.1s (loaded machine), targeted test 0.7s, probe 0.5s, vet 0.65s.
