# Review: REQ-729 (ai-report-judge states where a relative --out starts)

**Approve** — the documented rule is true against the code, both sentences sit in the right place, and nothing that matches the old usage string breaks.
Route A | range `cfa45037..572e8b94` (builder commit `82eb9936`, merge `572e8b94`)

## What's built

- The judge's usage text now ends with "a relative <bundle-dir> or --out starts at --repo-root, not the current folder". Behavior is unchanged (Decision D-01, keep and document).
- `ai-report.md` Step 7 gains one paragraph right after the judge command block. It names `--out <dir>` and states the same rule.

## Decisions / risks for you

None. D-01 (keep today's behavior) is sound: `handleAIReportJudge` sends the bundle path (`report_judge.go:110`) and `--out` (`report_judge.go:132`) through the same `judgeAbsolutePath(ctx.RepositoryRoot, ...)` helper (`report_judge.go:285`). Changing only `--out` would give one command two rules for relative paths.

## Checks run

1. Rule against code. `judgeAbsolutePath` joins a non-absolute path onto `ctx.RepositoryRoot` and cleans it, and leaves an absolute path as it is. Both the bundle and `--out` use it. The rule is true.
2. Usage text inside error messages. `report_judge.go:104` builds `err.Error()+"; "+judgeUsage`. Ran `do-work-cli.sh --repo-root <root> --format text ai-report-judge a b`; the finding reads `unexpected argument "b"; Usage: ai-report-judge <bundle-dir> [--out <dir>] [--browser <path>]; a relative <bundle-dir> or --out starts at --repo-root, not the current folder`. It reads as three clauses and stays clear.
3. Callers of the exact old string. `grep` over the repo (outside run artifacts and REQ files) finds the old form only in the 0.x changelog entries (`CHANGELOG.md:459`, `skills/do-work/CHANGELOG.md:459`), which describe a past release and must not change. No test matches `Usage:` or the usage string in `report_judge_test.go`, and no shell test or other Go file names it.
4. Placement. The new paragraph sits between the command block (line 181) and "The command serves the bundle ..." in Step 7. It does not contradict the nearby text, with one small wording overlap (F2 below).
5. `ai-report-guide.md` (named in the UR). Its judge paragraph (line 59) never mentions `--out` or relative paths, so nothing needs adding. The integrator recorded the same check.

## Findings

**Important:** None.

**Minor:**
- F1 `report_judge.go:39` and `ai-report.md:184`: "not the current folder" is exact only when `--repo-root` is passed. When it is omitted, `--repo-root` defaults to the working directory (`command_runtime.go:117`), so the two folders are the same. Every shipped caller passes `--repo-root`, so no one is misled in practice. — impact-negligible → report only

**Nit:**
- F2 `ai-report.md:184-186`: the new paragraph says `--out <dir>` chooses the output directory, and the next paragraph still says the captures go "to a fresh directory outside the bundle". "Fresh" is true only without `--out`. — impact-negligible → report only
- F3 `stakeholder-report.md:44` runs the same judge command with a `<bundle-dir>` that may be relative, without the rule sentence. The REQ named only `ai-report.md`, and that caller passes no `--out`, so this is outside scope. — impact-negligible → report only

## Requirements Checklist

- [x] Choose one: resolve against the current folder, or keep and state the rule — delivered (D-01 keeps behavior).
- [x] Rule stated in the usage text — delivered (`judgeUsage`, one clause).
- [x] Rule stated in the judge form of `ai-report.md` — delivered (one rule sentence, plus one sentence that names `--out`, which the action did not mention before).
- [x] An absolute `--out` behaves exactly as today — delivered (no code path changed; the integrator's extra case 1 confirms).
- [x] No test needed when only text changes — N/A by the REQ's own rule (D-02).
- [x] Constraint: no caller with an absolute path changes — met.
- [x] Constraint: the judge never writes inside the bundle — met (`report_judge.go:133-134` unchanged; integrator's extra case 2 confirms the refusal).
- [x] UR item: park the choice if callers depend on today's behavior — the callers pass no `--out`, so no parking was needed.

`CHANGELOG.md` and `VERSION` are in the write set but not in this range. The hand-back proposes the entry, and the run's finalization owns the release, so this is not a gap.

**Anti-bloat count:** 0 helpers, 0 options, 0 files, 0 tests that the REQ did not name.

## Acceptance Testing

**Result: Pass**
- Usage text checked live through the CLI (check 2 above).
- The relative-`--out` behavior was confirmed by the integrator (`## Qualification` in the REQ): a relative `--out captures` with `--repo-root` at a scratch root wrote under `<scratch root>/captures`.

## Suggested Additional Testing

None.

## Scores (on the record — not the headline)

**Overall: 98%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All requirements and constraints met |
| Code Quality | 95% | One constant string; F1 wording is approximate when `--repo-root` is omitted |
| Test Adequacy | N/A | Text-only change; the REQ waives a test |
| Scope | 100% | Two files, both in the write set |
| Risk | None | No behavior change |
| Acceptance | Pass | Usage text verified live |

## Follow-ups created
None (3 findings report only)

## Review

**Overall: 98%** | 2026-11-15T00:05:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1 "not the current folder" is exact only when `--repo-root` is passed; omitted, it defaults to the working directory (`command_runtime.go:117`) — impact-negligible → report only
- F2 (nit) `ai-report.md` next paragraph still says captures go to "a fresh directory", true only without `--out` — impact-negligible → report only
- F3 (nit) `stakeholder-report.md:44` runs the same judge command without the rule sentence; outside the REQ's scope — impact-negligible → report only
- Anti-bloat: 0 unnamed helpers, options, files or tests.

**Acceptance:** Pass — usage text verified live through the CLI; the rule matches `judgeAbsolutePath` for both the bundle and `--out`.
**Restatement sweep:** not run (this REQ is not the wave's last successful integration)
**Suggested testing:** 0 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*
