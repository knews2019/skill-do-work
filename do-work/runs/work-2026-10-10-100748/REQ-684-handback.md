# REQ-684 hand-back (qualify tells a missing summary section from one that names no files)

- Branch: worktree-agent-REQ-684-qualify-summary-messages
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-684-qualify-summary-messages
- Base: b629e5cd. Commit: 2a1a1789 (one commit, subject starts `[REQ-684]`).

## File manifest
- skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go (modified): `handleQualify` evidence is now "Implementation Summary section not found" when `!found`, "Implementation Summary lists no backticked file paths" otherwise; `parseError` branch, code, severity, fixability unchanged.
- skills/do-work/tools/do-work-cli/internal/corehelpers/checks_test.go (modified): two tests added.

## P-A-U
- [PLAN]: Write one test per message, see both RED on the old string, then split the single evidence string by the existing `found` value.
- [APPLY]: Two tests, then a 4-line change in `handleQualify` (default string changed, one `else if !found` added). No helpers, files or options.
- [UNIFY]: `git diff --stat`: checks.go 4 (+3/-1), checks_test.go 32 (+32), 35 insertions, 1 deletion. Checks: `go test -count=1 ./internal/corehelpers/` exit 0 (11.6s wall; whole package, existing qualify tests included); REQ-684-probe.sh exit 0 (0.7s); `gofmt -l` empty; `go vet ./...` exit 0; `git diff --check` exit 0. Files checked: checks.go, checks_test.go. No debug artifacts, nothing under `do-work/` staged.

## Red-green record
- TestQualifyImplementationSummarySectionNotFound: base output evidence `["Implementation Summary is missing or empty"]` (FAIL); after change PASS.
- TestQualifyImplementationSummaryListsNoBacktickedPaths: base output the same old string (FAIL); after change PASS.
- Each test asserts exactly one finding, code QUALIFY-SUMMARY-MISSING, and the exact evidence string.

## Decisions
- D-01 (DECIDE & STATE): the "lists no backticked file paths" string is the fallthrough default inside the existing `if`, because the only remaining way to reach it is found-and-empty. This avoids a second condition.
- D-02 (DECIDE & STATE): tests build their REQ file inline instead of reusing `writeQualificationRequest`, since that helper always writes a summary with paths and extending it would add a parameter.

## Discovered Tasks
- None.

## Lessons read
- _dev/primes/lessons-releases.md (whole). The Go satellite was not read; no code named by its bullets was touched.

## Anti-bloat check
```
 .../do-work-cli/internal/corehelpers/checks.go     |  4 ++-
 .../internal/corehelpers/checks_test.go            | 32 ++++++++++++++++++++++
 2 files changed, 35 insertions(+), 1 deletion(-)
```
Added that the REQ did not name: none (two tests and one `else if` branch, all named).

## Proposed CHANGELOG entry
Qualify says which Implementation Summary problem it found

A builder told "missing or empty" when the section exists but says "None, verification only" was misled. Qualify now names the real case.
- `QUALIFY-SUMMARY-MISSING` evidence reads "Implementation Summary section not found" when the section is absent.
- It reads "Implementation Summary lists no backticked file paths" when the section exists but names no files.
- Finding code, severity and fixability are unchanged; two focused tests pin the messages.

## Proposed lesson bullet
none

## Integration seams and wall times
None. Package test 11.6s, probe 0.7s, single run each, no load retries needed.
