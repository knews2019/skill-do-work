# REQ-729 Hand-back (ai-report-judge documents its relative --out rule)

Branch: `worktree-agent-REQ-729-judge-relative-out-rule`, base `df1bf423`, tip `82eb9936`.

## Commits
- `82eb9936` [REQ-729] ai-report-judge states that a relative bundle path or --out starts at --repo-root

## Files
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` — modified (only the `judgeUsage` constant, line 39).
- `skills/do-work-toolbox/actions/ai-report.md` — modified (one new paragraph in Step 7, after the command block, new lines 184-185).

## Scope changes
None.

## Red-green evidence
- Probe at base `df1bf423`: `REQ-729 probe FAIL: judgeUsage in report_judge.go does not state the relative --out rule` (exit 1).
- Probe at tip `82eb9936`: `REQ-729 probe PASS`.
- Usage before: `Usage: ai-report-judge <bundle-dir> [--out <dir>] [--browser <path>]`
- Usage after: `Usage: ai-report-judge <bundle-dir> [--out <dir>] [--browser <path>]; a relative <bundle-dir> or --out starts at --repo-root, not the current folder`
  (from `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format text ai-report-judge`, exit 2, finding TOOLBOX-USAGE).

## Verification
- `gofmt -l ./internal/toolboxcommands` — no output.
- `go vet ./internal/toolboxcommands/` — clean.
- `go test -count=1 -run 'TestAIReportJudge' ./internal/toolboxcommands/` — ok. Ran: `TestAIReportJudgeReportsSkippedWithoutBrowser` (PASS). Skipped (no browser engine): `TestAIReportJudgeFailsOnPhoneWidthOverflow`, `TestAIReportJudgeFailsOnBrokenRelativeImage`, `TestAIReportJudgePassesCleanBundleAndStopsServer`.
- `bash _dev/tests/action-shell-blocks.sh` — "Shell-block lint passed: 88 fenced blocks and 35 shipped shell files; ShellCheck enabled." exit 0. Single run, no rerun.
- GREEN probe — PASS (above).

## P-A-U
- PLAN: per D-01, append the rule clause to `judgeUsage` and add one paragraph to ai-report.md Step 7; no code behavior change.
- APPLY: two text edits exactly as briefed, gofmt clean.
- UNIFY: `git diff --stat` showed 2 files, 3 insertions, 1 deletion; checked report_judge.go changed only the constant string and ai-report.md gained only the one paragraph plus blank line; no debug artifacts.

## Decisions
- **D-02** (DECIDE & STATE): No new Go test. The change is text-only and the REQ says text-only needs none; the GREEN probe pins the wording.

## Discovered Tasks
- `skills/do-work-toolbox/docs/ai-report-guide.md` may want the same sentence if it describes the judge's `--out`; not checked in depth, outside the write set.

## Integration seams
- `report_judge.go`: only line 39 (`judgeUsage`). If REQ-728 adds a flag to this line, keep its flag and keep this trailing clause `; a relative <bundle-dir> or --out starts at --repo-root, not the current folder` at the end.
- `ai-report.md`: inserted lines 184-185 (the `Pass --out <dir> ...` paragraph and a blank line) between the Step 7 command block closing fence (line 182) and "The command serves the bundle ...". REQ-728, REQ-730, REQ-731 may edit nearby; keep the paragraph directly after the judge command block.

## Proposed CHANGELOG entry
**ai-report-judge States Where a Relative --out Starts**

What changes for an installed copy: the judge's usage text and the ai-report judge step now say that a relative bundle path or `--out` starts at `--repo-root`; behavior is unchanged.

## Proposed lesson
none

## Lessons read
- `_dev/primes/lessons-releases.md` (head).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — only the `[family: rule-direction-checked-against-callers]` bullets (0.305.81 REQ-645, 0.305.30).
