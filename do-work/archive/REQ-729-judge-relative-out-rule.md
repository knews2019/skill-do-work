---
id: REQ-729
title: '[impact-negligible] ai-report-judge resolves or documents a relative --out'
status: completed
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
required_lessons: ["_dev/primes/lessons-releases.md"]
related: [REQ-726, REQ-727, REQ-728, REQ-730, REQ-731, REQ-732]
batch: tinka-follow-ups
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go", "skills/do-work-toolbox/actions/ai-report.md", "CHANGELOG.md", "VERSION"]
claimed_at: 2026-11-14T23:28:18Z
route: A
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-11-14T23:30:17Z
  basis:
    - Route A
    - 4-file write set
    - 2 subsystems involved
    - 3 acceptance criteria
dispatch_at: 2026-11-14T23:32:29Z
builder_handback_at: 2026-11-14T23:36:35Z
integration_at: 2026-11-14T23:36:52Z
review_at: 2026-11-14T23:49:23Z
commit: 572e8b940f1d98e3d3524e4309cc691f5a16a6e1
heavy_verified_at: 2026-11-14T23:59:09Z
heavy_verified_revision: 572e8b940f1d98e3d3524e4309cc691f5a16a6e1
kb_status: pending
completed_at: 2026-11-15T00:00:21Z
release_at: 2026-11-15T00:00:21Z
---
# ai-report-judge Resolves or Documents a Relative --out
## What
Make the rule for a relative `--out` of `ai-report-judge` visible: either resolve it against the current folder, or state in the usage text and in the action that it is resolved against `--repo-root`.
## Why
On 2026-11-14, with `--repo-root` pointing at a scratch copy, a relative `--out` put the captures inside the copy and not where the caller expected. The usage text does not say which folder a relative path starts from.
## Context
- The option is parsed in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` (around lines 94 and 95) and the result is written with `filepath.Join(output, ...)`.
## Detailed Requirements
- Choose one: resolve a relative `--out` against the current folder, or keep today's behavior and state the rule in the usage text and in the judge form of `skills/do-work-toolbox/actions/ai-report.md`.
- An absolute `--out` behaves exactly as today.
- When the behavior changes, a Go test pins it. When only the text changes, no test is needed.
## Constraints
- No caller that passes an absolute path changes.
- The judge never writes inside the bundle.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** Run the judge with `--repo-root <scratch copy>` and `--out captures` from another folder, and read the usage text.
**Why RED now:** The captures land under the scratch copy, and nothing in the usage text or the action says that a relative `--out` starts at `--repo-root`.
**GREEN when:** Either the captures land under the current folder and a Go test pins it, or the usage text and the action state the rule in one sentence each.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Builder Guidance
- The choice is the builder's, after reading the existing callers. If a shipped caller depends on today's behavior, keep the behavior and document it.
- Record the choice as a decision in the request.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24500 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes do-work-cli Go code for reports.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Per D-01, append the rule clause to `judgeUsage` and add one paragraph to ai-report.md Step 7; no code behavior change. (from the builder hand-back)
- [x] **[APPLY]:** Two text edits exactly as briefed, gofmt clean. (from the builder hand-back)
- [x] **[UNIFY]:** `git diff --stat` showed 2 files, 3 insertions, 1 deletion; checked report_judge.go changed only the constant string and ai-report.md gained only the one paragraph plus blank line; no debug artifacts. (from the builder hand-back)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-judge-renders-one-named-page.md*

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names the two files and the exact change. D-01 below keeps today's behavior, so the work is one sentence in the usage text of `report_judge.go` and one sentence in the judge step of `ai-report.md`.

**Planning:** Not required

## Decisions

- **D-01** (DECIDE & STATE): Keep today's behavior and document it. A relative `--out` stays resolved against `--repo-root`. Reasoning: the same handler resolves the relative `<bundle-dir>` against `--repo-root` too (`report_judge.go:110` and `:132`, both through `judgeAbsolutePath` at `:285`), and `ai-report-second-opinion` uses that helper for its own relative paths (`report_second_opinion.go:161`, `:173`, `:385`). Changing only `--out` would give one command two rules for relative paths. The shipped callers (`skills/do-work-toolbox/actions/ai-report.md:181`, `skills/do-work-toolbox/actions/stakeholder-report.md:44`) pass `--repo-root <project-root>` and no `--out`, so documenting changes nothing for them. The rule sentence names both paths: a relative `<bundle-dir>` or `--out` starts at `--repo-root`, not the current folder. Value: no behavior change and no new test. Risk: none for absolute paths; a caller who expected the current folder still has to read the text.
- **D-02** (DECIDE & STATE, from the builder hand-back): No new Go test. The change is text-only and the REQ says text-only needs none; the GREEN probe pins the wording.
<!-- D-XX counter: last used D-02. Next decision: D-03. -->

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` (modified: the judgeUsage constant only)
- `skills/do-work-toolbox/actions/ai-report.md` (modified: one paragraph after the Step 7 judge command block)

**What was done:** Kept today's behavior (D-01) and stated the rule in both places: the judge's usage text and the ai-report judge step now say that a relative `<bundle-dir>` or `--out` starts at `--repo-root`, not the current folder.

## Discovered Tasks

(from the builder hand-back)
- `skills/do-work-toolbox/docs/ai-report-guide.md` may want the same sentence if it describes the judge's `--out`. Integrator check: the guide's only judge paragraph (line 59) never mentions `--out` or a relative path, so nothing to add → report only.

## Qualification

**Gate record:** `advance REQ-729 --diff-range cfa45037..572e8b94` returned the `qualify` gate `satisfied` with no findings (first call refused `QUALIFY-CLAIMED-PATH-MISSING` because the file list carried a backticked constant name; the backticks were removed, nothing else changed).

**Diff:** `git diff cfa45037..572e8b94 --stat`: 2 files, 3 insertions, 1 deletion. Both files are in the REQ's write set; no scope drift (Route A, no declared Scope).

**Requirement trace:**
- Choose one (resolve against the current folder, or keep and state the rule): D-01 keeps today's behavior. `report_judge.go` `judgeUsage` now ends with "; a relative <bundle-dir> or --out starts at --repo-root, not the current folder". `ai-report.md` Step 7 gains one paragraph directly after the judge command block: "A relative `<bundle-dir>` or `--out <dir>` starts at `--repo-root`, not the current folder." One sentence in each place, as GREEN asks.
- The stated rule is true: `handleAIReportJudge` resolves both the bundle (line 110) and `--out` (line 132) through `judgeAbsolutePath(ctx.RepositoryRoot, ...)`. Checked by running the judge from an unrelated folder with `--repo-root` at a scratch root: a relative `--out captures` wrote `judge.json` and the four captures under `<scratch root>/captures`, and nothing in the current folder.
- An absolute `--out` behaves as today: an absolute scratch `--out` received `judge.json` at exactly that path (extra case 1). No code path changed.
- The judge never writes inside the bundle: a relative `--out bundle/inside` was still refused with "--out must be outside the bundle" (extra case 2, edge).
- Text-only change, so no Go test (REQ rule, D-02). The GREEN probe pins the usage wording.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` (repository gate, stage reuse off) at merge `572e8b940f1d98e3d3524e4309cc691f5a16a6e1` in the run's detached checkout, started 2026-11-14T23:45:08Z after the 1-minute load fell to 4.90 with no other gate running; GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-729-probe.sh` through `advance`.
**Result:** ✓ Gate exit 0 in 141 s (single run, no retry). Probe exit 0; `advance` recorded the green-gate and run-blocked-check records as satisfied.

**Red-green validation:**
- `REQ-729-probe.sh`: ✗ at base `df1bf423` ("judgeUsage in report_judge.go does not state the relative --out rule", exit 1) → ✓ at builder tip `82eb9936` and at merge `572e8b94` (from the hand-back, re-run by `advance`).
- Builder checks (from the hand-back): `gofmt -l` clean, `go vet ./internal/toolboxcommands/` clean, `go test -count=1 -run 'TestAIReportJudge'` ok, `bash _dev/tests/action-shell-blocks.sh` exit 0.

**New tests added:**
- None (text-only change, D-02).

**Heavy verification plan:**
- Range: cfa4503728630b65893a4c24b6c5606cc6928f74..572e8b940f1d98e3d3524e4309cc691f5a16a6e1
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — report_judge.go matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — ai-report.md and report_judge.go matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — report_judge.go matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — report_judge.go matched subtree skills/do-work/tools/do-work-cli

*Verified by work action*

## Review

**Overall: 98%** | 2026-11-14T23:49:23Z

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

## Heavy Verification Plan

- Base: cfa4503728630b65893a4c24b6c5606cc6928f74
- Target: 572e8b940f1d98e3d3524e4309cc691f5a16a6e1
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — report_judge.go matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — ai-report.md and report_judge.go matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — report_judge.go matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — report_judge.go matched subtree skills/do-work/tools/do-work-cli

## Heavy Verification Result

- Target: 572e8b940f1d98e3d3524e4309cc691f5a16a6e1
- Execution revision: 572e8b940f1d98e3d3524e4309cc691f5a16a6e1 (run's detached drain-head checkout, `QUEUE_KANBAN_BROWSER` set, 1-minute load 4.78 at start)
- do-work-cli-integrations: executed, exit 0, 68 s
- staged-skills: executed, exit 0, 49 s
- updater: executed, exit 0, 73 s
- installer: executed, exit 0, 39 s

## Lessons Learned

**What worked:** Reading every caller of `judgeAbsolutePath` before choosing (D-01) kept one rule for every relative path in the command, so the fix was two sentences and no behavior change.

**What didn't:** The qualifier read a backticked constant name in the Implementation Summary file list as a path and refused; plain text fixed it.

**Worth knowing:** "Relative to `--repo-root`" equals "relative to the current folder" when `--repo-root` is omitted, because the root defaults to the working directory (review F1). No lesson bullet: the builder proposed none.

## Orientation

`ai-report-judge` (do-work-cli toolbox commands) resolves every relative path argument against `--repo-root` through `judgeAbsolutePath`; the usage text and the ai-report judge step now say so.

## Timing

Observed 2026-11-14T23:32:29Z to 2026-11-14T23:59:10Z: 26m 41s total, 12m 33s attributed across 5 events, 14m 08s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 6m 49s | 2 |
| builder-work | 4m 06s | 1 |
| review | 1m 24s | 1 |
| handback-merge | 14s | 1 |

Slowest stage: verification-gate / heavy drain, 4m 12s, outcome success.
