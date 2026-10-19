---
id: REQ-701
title: 'A failed second-opinion model check keeps the finished page and reports it; a subagent model never makes the backend model ambiguous'
status: completed
created_at: 2026-10-19T14:43:26Z
user_request: UR-158
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-700]
batch: UR-158
route: A
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-10-19T15:29:09Z
  basis:
    - Route A
    - 4-file write set
    - 4 acceptance criteria
claimed_at: 2026-10-19T15:25:12Z
dispatch_at: 2026-10-19T15:35:45Z
builder_handback_at: 2026-10-19T15:43:39Z
integration_at: 2026-10-19T15:44:20Z
review_at: 2026-10-19T16:02:02Z
kb_status: pending
commit: e8b063c1bd6257f5306a53b2d7049f9249891ce1
heavy_verified_at: 2026-10-19T16:26:17Z
heavy_verified_revision: e8b063c1bd6257f5306a53b2d7049f9249891ce1
completed_at: 2026-10-19T16:26:48Z
release_at: 2026-10-19T16:26:48Z
---
# A Failed Second-Opinion Model Check Keeps The Finished Page And Reports It; A Subagent Model Never Makes The Backend Model Ambiguous
## What
After a second-opinion backend finishes writing its page, the launcher's post-run model verification must never delete that page. When the actual model cannot be confirmed, keep the staged page and reply bytes in the bundle under an explicit unverified name, and report a typed finding naming the models seen. For the Claude backend, a subagent's model in the reply's `modelUsage` map is not ambiguity: the backend's own model is the one that was requested, or the one that answered the main turn, and extra entries from subagents are recorded, not treated as a second candidate.
## Why
A 12-minute generation was discarded because the editor used an Opus-based Explore subagent. The user had to recover the page from the transcript and rerun with the Agent tool denied through a temporary local setting. Destroying finished work on a probe verdict is the wrong failure mode; a page with an uncertain model label is still useful and the user can judge it.
## Context
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` line 394: `defer func() { _ = bundleRoot.Remove(stageName) }()` removes the staging page on every early return of `secondOpinionRun`, including a metadata error.
- Same file, `secondOpinionMetadata` for `claude`: `models = envelope.ModelUsage`, then `if len(models) != 1 { return ... "actual model is absent or ambiguous (%d models)" }`. Claude Code's result envelope lists every model that ran, subagents included, so any Explore or Agent dispatch inside the editorial pass trips this.
- The follow-up screenshot (see Assets) shows the agent's own diagnosis: "Found it: the editor spawned an Opus-based Explore subagent, which confused the launcher into thinking two models were ambiguous, causing it to delete the finished page."
## Requirements
- No code path in `secondOpinionRun` deletes a staged page or reply that the backend finished writing. Cleanup of the staging file is limited to the case where the backend never wrote it.
- When model verification fails after a finished page exists, publish the page and reply under a name that says the model is unverified (the builder picks the exact form, consistent with `index-<cli>-<model>[-N].html`), insert no banner, and emit a typed finding listing the requested model and every model seen.
- For `claude`, the actual model is the requested `--model` when it appears in `modelUsage`; otherwise the single entry when there is exactly one. Extra subagent entries are recorded in the run row, not treated as ambiguity. Only a map that is empty, or that lacks the requested model and has several entries, is a verification failure.
- Lock both behaviors in `report_second_opinion_test.go`: a fake backend that writes the page and returns a two-model envelope ends with the page still on disk and a finding; a fake that returns requested model plus one subagent model ends with the normal verified publication.
## Red-Green Proof
**RED prompt/case:** In `report_second_opinion_test.go`, a fake `claude` executable writes the staging page and prints a success envelope whose `modelUsage` has two entries, `claude-fable-5-1` (requested) and `claude-opus-5-5`. Assert the page exists after the run and the run row reports the requested model.
**Why RED now:** `len(models) != 1` returns an error, the deferred `Remove(stageName)` deletes the page, and the run row says "actual model is absent or ambiguous (2 models)".
**GREEN when:** The page is published under the requested model's name with the normal judge and banner flow, the subagent model is recorded in the row, and a separate test with an envelope lacking the requested model and holding two unrelated models ends with the page kept under an unverified name, no banner, and a typed finding.
**Validation:** Inferred during capture from the screenshot and the code at the cited lines.
## Assets
- `do-work/user-requests/UR-158/assets/REQ-701-screenshot-2-second-opinion-probe-deleted-page.png` — Claude Code transcript continuing the REQ-700 run: "Background command 'Run the second-opinion generation with the claude backend' completed (exit code 0)", "The launcher failed because its model probe found two models from the claude CLI and treated that as ambiguous.", "Still diagnosing: the generation itself ran for 12 minutes and then was discarded.", "Found it: the editor spawned an Opus-based Explore subagent, which confused the launcher into thinking two models were ambiguous, causing it to delete the finished page. I'll recover the page from the transcript as a fallback and rerun with the Agent tool blocked for the child process.", "Recovered page saved to the scratchpad as a fallback. Now adding a temporary Agent deny to the gitignored local settings so the child run cannot spawn a second-model subagent, then rerunning.", then the rerun command with `--backend claude --model claude-fable-5-1`.
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 22085 tokens, over the 2000 budget; `slugged: partial` so no targeted entry is legal. Matched on native backend identity and rollback/deletion identity in do-work-cli internals (families native-backend-identity, final-boundary-identity, silent-skip-reads-as-red).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read general, coding-guardrails, shared-principles, communication-style, testing; primes prime-shell-commands and prime-releases; the `native-backend-identity` lesson bullet (REQ-698). Approach exactly as the brief: tests first against a signature-only scaffold (D-02), then the defer rule, the unverified branch, the requested-model-wins rule, the row field, two doc sentences. (from the builder hand-back)
- [x] **[APPLY]:** Code written as planned, confined to the five files. No prepare/finalize work, no flags, no Session-form rewrite (REQ-700 territory untouched). (from the builder hand-back)
- [x] **[UNIFY]:** `git diff cc3e72c2 --stat`: 5 files changed, 182 insertions(+), 40 deletions(-). Checks from the worktree root: RED three-test run exit 1 (expected, 2 s); GREEN family `-run '^TestSecondOpinion'` exit 0, 19 pass, 0 skip (13 s); `go test ./internal/resultmodel/ ./internal/toolboxcommands/` exit 0 (14 s); `go vet` on both packages exit 0; `bash REQ-701-probe.sh` exit 0 (2 s); `gofmt -l skills/do-work/tools/do-work-cli/internal` prints nothing; `git diff --check` exit 0. Files reviewed: the full source diff of `report_second_opinion.go` and `result_model.go` (read line by line; no debug prints, comments updated), the test diff, both doc sentences in place (grep confirmed one match each). (from the builder hand-back)
*Source: [Image #3] <- see additional feedback*

---

## Triage

**Route: A** - Simple

**Reasoning:** Bug fix with a clear reproduction: the REQ names the file (`report_second_opinion.go`), the two defects (`secondOpinionRun`'s deferred `Remove(stageName)` and `secondOpinionMetadata`'s `len(models) != 1` rule for `claude`), and the two test fakes to lock it. Nothing to discover; the existing `secondOpinionFakeJSON` fake and `secondOpinionRows` helper cover the test shape.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)

**What was done:** `secondOpinionRun` now keeps a non-empty staging page on every refusal (the deferred cleanup removes it only when published or empty and otherwise appends "staged page kept at <path>" to the failure reason), and when the post-run model check fails after the page is written it publishes the page and reply as `index-<cli>-<model>-unverified[-N].html` / `<cli>-<model>-unverified[-N]-reply.md` with run state `unverified`, no backlink check, no judge and no banner, so the handler's existing `AI-REPORT-SECOND-OPINION-<STATE>` finding becomes `AI-REPORT-SECOND-OPINION-UNVERIFIED` with the requested and seen models in its evidence. `secondOpinionMetadata` takes the expected model and returns the sorted set of models seen: for `claude` and `gemini` the requested model wins when the native usage map lists it, a single entry wins otherwise, and only an empty map or a multi-entry map lacking the requested model returns the new `errSecondOpinionModelUnverified` sentinel (with the reply and seen set attached). The run row gains `seen_models` (`SeenModels` on `AIReportSecondOpinionRun`, always a JSON array). Tests: helper `secondOpinionFakeClaude`, new `TestSecondOpinionSubagentModelIsNotAmbiguous` and `TestSecondOpinionUnverifiedModelKeepsFinishedPage`, the codex mismatch test renamed to `TestSecondOpinionFreshGenerationModelMismatchPublishesUnverified` with flipped assertions, and a kept-staging assertion in the missing-backlink case. One sentence each in the action file and the guide describes the unverified name, state and `seen_models`.

## Decisions

(from the builder hand-back)

- D-01 (pre-dispatch): the requested-model-wins rule lives in the shared gemini/claude branch of `secondOpinionMetadata`, not behind a `claude` condition. Value: one code path; gemini's `stats.models` has the same shape. Risk: none observed; the gemini empty-map case still errors.
- D-02 DECIDE & STATE: the RED run used a signature-only scaffold (new `secondOpinionMetadata` shape, `SeenModels` field, no behavior change) so the three tests failed on their assertions, not on compilation. A compile failure is not RED-before-GREEN under testing.md. The failure texts match the brief's expectation word for word.
- D-03 DECIDE & STATE: the "staged page kept at" path is repository-relative (`ai-reports/<bundle>/.codex-...-staging-X.html`), the same form as `row.PagePath`, which is why `relative` moved up as the brief said. If the staging file is already gone when the defer runs (stat fails), nothing is removed and nothing is appended, so the reason never names a file that does not exist.
- D-04 DECIDE & STATE: `seen` for codex and agy is derived once after the switch (`if seen == nil && model != ""`) rather than inside each case; same result, one line. The final slug-pattern failure still returns `nil` seen as a non-identity error, per brief outcome (a).
- D-05 DECIDE & STATE: `row.SeenModels = append([]string{}, seen...)` so a `nil` seen from a non-identity error keeps the JSON at `[]`, never `null`.

## Discovered Tasks

(from the builder hand-back)

- impact-low: `secondOpinionRun` sets `row.ActualModel = resolved` from the probe before generation, so a non-identity generation failure (outcome a) reports the probe's model as `actual_model` while `state` is `failed`. Pre-existing, unchanged here → report only.
- impact-low: a kept staging file (`.<cli>-<model>-staging-*.html`) joins the `before` snapshot of the next run, so later runs tolerate it, but nothing ever cleans it up; the user deletes it by hand after reading it → report only.

## Qualification

**Gate records:** `advance REQ-701 --diff-range 6256fd58..5b1536ac` returned one `qualify` record, provenance `merged_range`, state `satisfied`, no findings (no debug artifacts, no unfinished-work markers, all three P-A-U boxes ticked with the hand-back's text).

**Diff read:** `git diff 6256fd58..5b1536ac --stat` is 5 files, 182 insertions, 40 deletions, the same five files the hand-back manifest names; the full source and test diff was read line by line.

**Requirement trace:**
- No code path deletes a finished staged page: the deferred closure after `os.CreateTemp` stats the staging leaf and removes it only when `stagePublished` is true or the file is empty; otherwise it appends `staged page kept at <repo-relative path>` to the failure reason. The missing-backlink test now asserts the kept non-empty staging file and that reason text. Delivered.
- Verification failure after the page exists publishes it unverified, no banner, typed finding with requested and seen models: `unverified := err != nil || actual != resolved` routes to `secondOpinionNames(bundle, backend, resolved+"-unverified")`, sets `row.RunState = "unverified"` and returns before the judge and banner; the handler's existing `AI-REPORT-SECOND-OPINION-<STATE>` finding therefore reads `AI-REPORT-SECOND-OPINION-UNVERIFIED` with evidence `actual model unverified: requested <model>; seen <models>`. `TestSecondOpinionUnverifiedModelKeepsFinishedPage` asserts the finding code, all three model names in its evidence, the page and reply bytes, `banner_state: not-inserted`, `judge_verdict: not-run`, and an unchanged `index.html`. Delivered.
- Claude requested model wins, extras recorded, only empty or multi-without-requested fails: the shared gemini/claude branch of `secondOpinionMetadata` picks `expected` when the usage map lists it, else the single entry, else returns `errSecondOpinionModelUnverified` (empty map falls into this branch because `len(seen)` is 0); every seen model lands in the row's new `seen_models` array. `TestSecondOpinionSubagentModelIsNotAmbiguous` asserts state `ran`, the requested model's page name and `seen_models` of both entries. Delivered.
- Both behaviors locked in `report_second_opinion_test.go`: the two tests above plus the renamed codex mismatch test (now asserts the unverified pair and no banner). Delivered.

**Scope:** Route A, no Scope section; the five touched files are the ones the REQ's Context and Requirements name (the Go source, its test, the result model for the new row field, and one sentence in each shipped doc). No helpers, flags or files beyond the brief's named items; the two new locals and the `fmt` test import are short-lived.

*Verified by work action*

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge 5b1536ac (run 2: exit 0, gate wall 126 s, started 2026-10-19T15:51:18Z at 1-minute load 4.60); `bash do-work/runs/work-2026-10-19-152525/REQ-701-probe.sh` through `advance` (exit 0; the three named tests each printed `--- PASS`).
**Result:** ✓ All passing (do-work-cli 939 tests in 62 s, queue-kanban 427 tests in 39 s, contract regressions and shell audits green; probe 3 of 3)

**Repository gate retry:** run 1 (started 2026-10-19T15:45:42Z, load 2.31) printed `Maintainer verification passed` with gate wall 126 s, but the wrapper around it lost the exit status (zsh treats `status` as read-only), so run 2 was made after waiting for a quiet machine and is the recorded run. Both runs green.

**Re-merge after review (range now 6256fd58..e8b063c1):** review finding M1 (an unwritten page plus a failed model check published a zero-byte unverified page) was fixed on the builder branch (`ebbc3d40`) and re-merged as `e8b063c1`. Gate run 3 at e8b063c1 (started 2026-10-19T16:11:10Z, load 2.55 rising to 5.9) exited 1: `TestCancelledCommitKillsTermDeafHookBeforeReturning` and `TestCancelledCommitThatLandsReportsCommittedRisk` in `internal/gittransaction` failed with "no process id was recorded at .../hook.pid" after 3.09 s, a hook that did not start in time under load; the package does not touch this diff. Under the one-retry rule, run 4 (started 2026-10-19T16:16:52Z, load 3.50) exited 0 with gate wall 137 s (do-work-cli 940 tests in 65 s, queue-kanban 427 in 45 s) and is the recorded run, written with `record-green-gate` because `advance` had moved past the test-gate phase. `bash do-work/runs/work-2026-10-19-152525/REQ-701-probe.sh` at e8b063c1: exit 0, 3 of 3 pass.

**Red-green validation:** (from the hand-back's proof record, traced to `## Red-Green Proof`)
- `TestSecondOpinionSubagentModelIsNotAmbiguous` (the captured RED case: two-entry `modelUsage` holding the requested model): ✗ before ("subagent model treated as ambiguity ... absent or ambiguous (2 models) ... state:failed") → ✓ after
- `TestSecondOpinionUnverifiedModelKeepsFinishedPage` (the captured GREEN second case: envelope lacking the requested model): ✗ before ("unverified model did not keep the finished page ... state:failed") → ✓ after
- `TestSecondOpinionFreshGenerationModelMismatchPublishesUnverified`: ✗ before ("model switch not published as unverified ... substitution refused") → ✓ after
RED ran against a signature-only scaffold (D-02) so all three failed on assertions, not compilation; GREEN family run 19 of 19 pass, 0 skips.

**New tests added:**
- `TestSecondOpinionSubagentModelIsNotAmbiguous`, `TestSecondOpinionUnverifiedModelKeepsFinishedPage` and helper `secondOpinionFakeClaude` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`
- `TestSecondOpinionUnverifiedModelWithoutPageFails` (added at the re-merge for review finding M1: an unverifiable model with no page written is a plain failure and publishes nothing)

**Existing tests updated (cross-REQ impact):**
- `report_second_opinion_test.go` (from REQ-698, the original second-opinion command): `TestSecondOpinionFreshGenerationModelMismatchRefusesPublication` renamed to `...PublishesUnverified` with flipped assertions (intentional behavior change: a model switch now publishes an unverified pair instead of refusing); the six-case identity test passes the new `expected` argument; the missing-backlink case asserts the kept staging file.

**Heavy verification plan:**
- Range: 6256fd58c5729cde66caa67d0fbdbfcfe7cd5ba4..e8b063c1bd6257f5306a53b2d7049f9249891ce1 (re-planned at the re-merge; the same four lanes were selected for 6256fd58..5b1536ac and ran green there first)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the three do-work-cli files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all five changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the three do-work-cli files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the three do-work-cli files matched subtree `skills/do-work/tools/do-work-cli`

*Verified by work action*

## Review

**Overall: 95%** | 2026-10-19T16:02:02Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** M1 `report_second_opinion.go` unverified branch publishes a zero-byte page and reply when the backend wrote nothing and the model check failed, where the old path was a plain `failed` — impact-user-visible → report only; M2 `lessons-do-work-cli.md:5` (`native-backend-identity`, REQ-698) still says mismatched identity is refused before naming, now reconciled only by the REQ-701 bullet the hand-back proposes — impact-negligible → report only; N1 no gemini case pins the shared requested-model-wins branch (D-01) — impact-negligible → report only; N2 summary line prints `model=` empty for unverified rows — impact-negligible → report only. Anti-bloat: no helpers, options, flags, files or decorative tests beyond the named items (unnamed: locals `unverified`, `publishedModel`, test local `names`, `fmt` test import).
**Acceptance:** Pass — Implementation and Integration stages: `go test -run '^TestSecondOpinion'` at 5b1536ac, 19 of 19 pass; full gate per the Testing section not rerun
**Restatement sweep:** redefined the published filename contract (`index-<cli>-<actual-model>[-N].html` gains the `-unverified` form), the run row schema (`seen_models`, state `unverified`), the `AI-REPORT-SECOND-OPINION-<STATE>` code set (`-UNVERIFIED`), and the identity rule (mismatch publishes unverified; a subagent model is not ambiguity); swept `skills/` and `_dev/`: `actions/ai-report.md` and `docs/ai-report-guide.md` updated in this diff, `ai-report-reference.md`, `SKILL.md`, `do-work-cli.sh`, `_dev/tests/do-work-cli-launcher-behavior.sh` and the 0.306.0 changelog entry hold no stale restatement ("never substitutes" still holds); one stale restatement in `lessons-do-work-cli.md:5` recorded as M2 | unread for the sweep: REQ-700
**Suggested testing:** 3 items
**Follow-ups created:** None (4 findings report only)

**Delta re-review (5b1536ac..e8b063c1), 2026-10-19T16:20:34Z:** M1 was taken as a fix on the builder branch (a `len(page) == 0` guard on the unverified path plus `TestSecondOpinionUnverifiedModelWithoutPageFails`); the reviewer read the delta only and reports M1 resolved, no new defect or bloat, Code Quality 90% → 95%, Test Adequacy 90% → 95%, Overall 95% → 97%; remaining findings M2, N1, N2 → report only. Full text in `do-work/runs/work-2026-10-19-152525/REQ-701-review.md`.

*Reviewed by review-work action*

## Heavy Verification Plan

- Base: `6256fd58c5729cde66caa67d0fbdbfcfe7cd5ba4`
- Target: `e8b063c1bd6257f5306a53b2d7049f9249891ce1` (re-planned at the re-merge; the same four lanes were selected for target `5b1536ac936af09711fc009a261bf8c517e063d6` and ran green there at 2026-10-19T15:55:28Z)
- `do-work-cli-integrations`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — reasons: `result_model.go`, `report_second_opinion.go`, `report_second_opinion_test.go` matched subtree `skills/do-work/tools/do-work-cli`
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — reasons: all five changed paths matched subtree `skills`
- `updater`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — reasons: the three do-work-cli paths matched subtree `skills/do-work/tools/do-work-cli`
- `installer`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — reasons: the three do-work-cli paths matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: `e8b063c1bd6257f5306a53b2d7049f9249891ce1`; execution revision: `e8b063c1bd6257f5306a53b2d7049f9249891ce1` (detached checkout `.git/work-run-work-2026-10-19-152525/drain-head-REQ-701-2`, started 2026-10-19T16:22:01Z, ended 2026-10-19T16:25:35Z, `QUEUE_KANBAN_BROWSER` set)
- `do-work-cli-integrations`: executed, exit 0, 72 s
- `staged-skills`: executed, exit 0, 41 s
- `updater`: executed, exit 0, 67 s
- `installer`: executed, exit 0, 29 s
- Earlier drain at the first merge `5b1536ac936af09711fc009a261bf8c517e063d6` (2026-10-19T15:55:28Z to 15:58:58Z): the same four lanes executed, all exit 0 (66, 46, 67, 27 s); superseded by the re-merge.

## Lessons Learned

**What worked:** RED against a signature-only scaffold (D-02) made the three failures assertion failures, so the proof traces cleanly to `## Red-Green Proof`. Putting the requested-model-wins rule in the shared gemini/claude branch kept one code path. The review's one real finding (an empty staging page published as unverified) was a one-line guard plus one pinning test, fixed on the builder branch and re-merged with the same `<pre>`.
**What didn't:** the integrator's gate wrapper assigned `status=$?` under zsh, where `status` is read-only, so gate run 1 lost its exit status and had to be rerun. Gate run 3 went red on two `internal/gittransaction` cancellation tests ("no process id was recorded at .../hook.pid" after 3.09 s) while the 1-minute load passed 5.9; both are hook-start timing under load, not this diff, and run 4 on a quiet machine was green.
**Worth knowing:** the handler derives the finding code from `strings.ToUpper(row.RunState)`, so the new `unverified` state became `AI-REPORT-SECOND-OPINION-UNVERIFIED` without a new literal. An unverified row has an empty `actual_model` and keeps `requested_model`; the models the backend reported are in `seen_models`. A staging page kept after a refusal (`.<cli>-<model>-staging-*.html`) is never cleaned up automatically (Discovered Tasks).

## Orientation

- Now a second-opinion run that finishes its page never loses it: when the backend's model cannot be confirmed afterwards, the page and reply are published as `index-<cli>-<model>-unverified[-N].html` with a typed `AI-REPORT-SECOND-OPINION-UNVERIFIED` finding, and a subagent's model beside the requested one in Claude's usage map is recorded in `seen_models` instead of read as ambiguity. Lives in the do-work-toolbox `ai-report second-opinion` launcher (`do-work-cli` toolboxcommands). Not a map change: one new run state and one new row field, no new module, flow or contract token beyond those. Primes `prime-shell-commands.md` and `prime-releases.md` are untouched by this change and their referenced paths still exist.

## Timing

Observed 2026-10-19T15:35:45Z to 2026-10-19T16:26:17Z: 50m 32s total, 30m 37s attributed across 6 events, 19m 55s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 14m 40s | 3 |
| builder-work | 7m 54s | 1 |
| review | 7m 22s | 1 |
| handback-merge | 41s | 1 |

Slowest stage: builder-work / builder worktree build, 7m 54s, outcome success.
Slowest command: verification-gate / repository gate maintainer-verify runs 3+4 at re-merge e8b063c1 (run 3 red under load, run 4 is the recorded run), 7m 59s, exit 0, .
