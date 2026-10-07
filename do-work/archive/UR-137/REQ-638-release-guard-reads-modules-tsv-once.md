---
id: REQ-638
title: '[impact-negligible] Release guard reads suite/modules.tsv once instead of merging two reads of the same declaration'
status: completed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-07T15:01:53Z
created_at: 2026-10-06T22:45:18Z
user_request: UR-137
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
related: [REQ-635, REQ-636, REQ-637]
batch: review-0305-69-findings
required_lessons: [_dev/primes/lessons-releases.md]
write_set: [skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go, skills/do-work/CHANGELOG.md]
claimed_at: 2026-10-07T15:01:06Z
dispatch_at: 2026-10-07T13:16:29Z
builder_handback_at: 2026-10-07T13:20:23Z
integration_at: 2026-10-07T15:02:13Z
review_at: 2026-10-07T15:06:00Z
commit: 6125e58c49d69acb059f2859eba5f2e661f87f6b
heavy_verified_at: 2026-10-07T15:09:41Z
heavy_verified_revision: 6125e58c49d69acb059f2859eba5f2e661f87f6b
kb_status: pending
completed_at: 2026-10-07T15:09:56Z
release_at: 2026-10-07T15:09:56Z
---
# Release Guard Reads suite/modules.tsv Once Instead of Merging Two Reads of the Same Declaration

## What
Review finding F8 (cosmetic). `releaseShippedChangeError` (`skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:36-55`) calls `DeclaredMaintainerReleaseRoots`, then `DeclaredModuleSources`, then merges the two through a set. `DeclaredMaintainerReleaseRoots` (`releaseownership/release_ownership.go:127-139`) already calls `DeclaredModuleSources` and filters it, so roots are a subset of sources and the merge is redundant; `headReleaseImage` is also built twice. Simplify to one read with identical behaviour.

## Detailed Requirements
1. Call `DeclaredModuleSources` once. If no declared source has a tracked `<source>/VERSION`, return nil (the repository is a consumer project and is not guarded), exactly as today's `len(roots) == 0` branch does.
2. Otherwise the shipped roots are the declared sources plus `suite` and `tools`, sorted. The output must be byte-identical to today's for every input.
3. Build `headReleaseImage(repositoryRoot)` once.
4. Keep the existing tests green and run the full Go test suite for the do-work CLI module. No new test is required: the change is a deletion with identical behaviour and no runnable RED exists.
5. Release per `_dev/primes/prime-releases.md` with a CHANGELOG entry (shipped Go changes).

## Constraints
- Behaviour-preserving only; no change to `DeclaredMaintainerReleaseRoots` or `DeclaredModuleSources` themselves, which have other callers.
- Batch constraint: each REQ in this batch is its own release and its own commit.

## Dependencies
None. Independent of REQ-635, REQ-636 and REQ-637.

## Builder Guidance
Certainty is high. This is a deletion. Builder latitude: whether the tracked-`VERSION` check stays a small loop inline or becomes a one-line helper.

## Red-Green Proof
**RED prompt/case:** Not applicable as a failing test: the change preserves behaviour. Proof is the existing `finalization` test suite green before and after, plus a diff that removes the second `DeclaredModuleSources` read, the second `headReleaseImage` call and the set merge.
**Why RED now:** The function parses `suite/modules.tsv` twice and merges a subset into its superset.
**GREEN when:** One parse, one image build, existing tests green, `go vet` clean.
**Validation:** Inferred during capture. Confirmed by reading both functions during triage.

## Full Context
See `do-work/user-requests/UR-137/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `DeclaredModuleSources` once through one `headReleaseImage`. The consumer-project exit becomes "no declared source has a tracked `<source>/VERSION`", which is exactly `len(DeclaredMaintainerReleaseRoots(...)) == 0` because that function filters the same sources on the same condition. Roots = sources plus `suite` and `tools` when absent, sorted; byte-identical because the sources are already sorted and unique and the old roots were a subset of them. (Ticked by the orchestrator from the builder hand-back.)
- [x] **[APPLY]:** One file edited as planned (`finalization_release_guard.go`, +8/-13); the tracked-VERSION check is an inline `slices.ContainsFunc`. (Ticked by the orchestrator from the hand-back, cross-checked against the merge.)
- [x] **[UNIFY]:** `git diff --stat da3ab5f5..6125e58c`: 1 file, the one in write_set. `gofmt -l` empty, `go vet` on finalization and releaseownership clean, full module `go test -count=1 ./...` green (builder); no debug artifacts; `sort` still used; `DeclaredMaintainerReleaseRoots` keeps its other caller in releaseownership. (Ticked by the orchestrator from the hand-back, cross-checked against the merge range.)
*Source: review finding F8 from the pasted consumer code review of 0.305.58 to 0.305.69, triaged with `do-work validate-feedback` on 2026-10-06 and accepted.*

## Triage

**Route: A** - Simple

**Reasoning:** One function in one file, the REQ names the exact lines, the two readers it merges, and the byte-identical output it must keep. No exploration is needed: the only reach is `releaseShippedChangeError` itself, and both `releaseownership` readers stay unchanged.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go` (modified)

**What was done:** `releaseShippedChangeError` now calls `releaseownership.DeclaredModuleSources` once, with one `headReleaseImage`. It returns nil (consumer project, not guarded) when no declared source has a tracked `<source>/VERSION`, which is the same condition as the old `len(DeclaredMaintainerReleaseRoots(...)) == 0` exit. The shipped roots are the declared sources plus `suite` and `tools` when absent, sorted, so the root list and the refusal text are unchanged for every input. The second declaration read, the second image build and the set merge are gone. `DeclaredMaintainerReleaseRoots` and `DeclaredModuleSources` are unchanged. Merge range da3ab5f5..6125e58c (builder commit 4650565a, branch cut at 13932fc2 before 0.305.70, merged as is).

## Qualification

**Diff range:** da3ab5f5..6125e58c (builder commit 4650565a, merge 6125e58c)
**Gate record:** `advance` qualify satisfied with no findings: P-A-U boxes ticked from the hand-back, no debug artifacts, the one changed file is in write_set.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. 1: one `DeclaredModuleSources` call; the nil return now fires when `slices.ContainsFunc` finds no source with a tracked `<source>/VERSION`, the exact filter `DeclaredMaintainerReleaseRoots` applies, so the consumer-project exit is unchanged. 2: `DeclaredModuleSources` returns `sortedUniqueRoots(sources)`, the old roots were a subset of those sources, and the old set merge only removed duplicates; appending `suite` and `tools` only when absent and sorting gives the same list, so the refusal text is byte-identical. 3: one `headReleaseImage(repositoryRoot)` call. 4: no new test (tdd: false); the existing release-guard tests and the full module are green (hand-back), and the gate runs at the merge. 5 is finalization's release. Error paths are the same: both old reads returned the same `DeclaredModuleSources` error, wrapped the same way.
**Live data flow:** finalization calls `releaseShippedChangeError` for every manifest with `release_manifest_path`; this run's own finalization exercises the new code on the merged tree.
**P-A-U honesty:** the three boxes were ticked by the orchestrator from the builder's hand-back and checked against `git diff --stat da3ab5f5..6125e58c` (1 file).

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at 6125e58c (detached checkout `.git/work-run-2026-10-07-150106/drain-head`)
**Result:** ✓ All passing — exit 0 on the first run, 15:02:48Z to 15:05:08Z, gate wall 140s; stage do-work-cli-fast-tests executed (876 tests, wall 64s, slowest file internal/finalization/finalization_recovery_test.go 23.52s < 30s); queue-kanban-fast-tests executed. Load average about 4 at launch, no other gate running. Green-gate record satisfied by advance.

**Focused tests:** `do-work/runs/work-2026-10-07-150106/REQ-638-probe.sh` (`go test -run '^TestRelease(RefusedWhenTheImplementationShipsNothing|Guard)' ./internal/finalization/`) → exit 0, advance run-blocked-check and green-gate records satisfied. Builder-side: full module `go test -count=1 ./...` green (32 test packages), `gofmt -l` empty, `go vet` on finalization and releaseownership clean.

**Red-green validation:** `tdd: false`; traced to `## Red-Green Proof`, which states no runnable RED exists for a behaviour-preserving deletion. Proof is the existing release-guard tests (7 tests in `finalization_release_guard_test.go`, 9 `TestRelease*` in the package) green before (builder, at 13932fc2) and after (at 6125e58c), plus the diff that removes the second declaration read, the second image build and the set merge.

**New tests added:** none (no runnable RED; the REQ requires none)

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: da3ab5f5..6125e58c (4 of 6 lanes; the one changed path sits under `skills/do-work/tools/do-work-cli`)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — matched subtree skills/do-work/tools/do-work-cli

*Verified by work action*

## Review

**Overall: 96%** | 2026-10-07T15:06:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

Review depth calibrated to an impact-negligible, behaviour-preserving deletion: one focused independent review of the merge range, with the equivalence argued input by input (consumer-project exit, error paths, root list and order, a source named `suite` or `tools`, slice aliasing).

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- M1 `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:44` — `roots := shippedSources` then `append` and `sort` reuse the source slice; safe today because `shippedSources` is not read again, and the old code reused `roots[:0]` the same way. — impact-negligible → report only
- M2 `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go:168` — no test pins the "modules declared, but no declared source has a tracked VERSION" exit, the one condition this REQ rewrote; its equivalence is shown by reading the code. — impact-negligible → report only
- M3 `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go:62` — the refusal tests assert only the `RELEASE-WITHOUT-SHIPPED-CHANGE` token, not the root list in the message. — impact-negligible → report only

**Acceptance:** Pass — the focused probe and all 9 `TestRelease*` tests pass at 6125e58c; `go vet` and `gofmt` clean; the diff is one file inside write_set.
**Suggested testing:** 3 items (in the reviewer report)
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

Full reviewer report: `do-work/runs/work-2026-10-07-150106/REQ-638-review.md`.

## Heavy Verification Plan

- Base revision: da3ab5f5aa8c825ee7a6743c3952b79322a67c3b
- Target revision: 6125e58c49d69acb059f2859eba5f2e661f87f6b (landed in `commit:`)
- do-work-cli-integrations — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — finalization_release_guard.go matched subtree skills/do-work/tools/do-work-cli
- staged-skills — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — matched subtree skills
- updater — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — matched subtree skills/do-work/tools/do-work-cli
- installer — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — matched subtree skills/do-work/tools/do-work-cli

## Heavy Verification Result

- Target revision: 6125e58c49d69acb059f2859eba5f2e661f87f6b
- Execution revision: 6125e58c49d69acb059f2859eba5f2e661f87f6b (detached drain checkout `.git/work-run-2026-10-07-150106/drain-head`, run 15:05:18Z to 15:09:10Z, QUEUE_KANBAN_BROWSER set to Google Chrome)
- do-work-cli-integrations: exit 0, executed, 82s
- staged-skills: exit 0, executed, 46s
- updater: exit 0, executed, 72s
- installer: exit 0, executed, 30s

Green: every selected lane present, exit 0, none skipped, none reused; no HEAVY-RUN-LANE-SKIPPED finding in the log.

## Lessons Learned

**What worked:** Arguing equivalence from the reader's own contract: `DeclaredMaintainerReleaseRoots` is a filter of `DeclaredModuleSources`, and the latter returns sorted unique entries, so one read plus the same predicate gives the same exit and the same root list. A one-file behaviour-preserving deletion needed no exploration, no pre-flight and one focused review.
**What didn't:** Nothing failed. The rewritten consumer-project exit (modules declared, none versioned) has no test of its own (review M2), so its equivalence rests on reading the code.
**Worth knowing:** When one reader is derived from another, call the base reader once and apply the derived filter as a predicate, instead of reading both and merging a subset into its superset. No lesson-satellite entry: there was no incident, and the builder and the orchestrator agree one is not warranted (D-06).

## Orientation

Now the finalization release guard reads `suite/modules.tsv` once, with one HEAD image, and still treats a repository with no versioned declared module as an unguarded consumer project; the shipped-root list and refusal text are unchanged. Lives in the do-work-cli finalization release guard (`skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, `_dev/primes/prime-releases.md`). No map change. Prime spot-check: neither prime describes the old two-read merge.

## Decisions

Builder decisions, copied from the hand-back (`do-work/runs/work-2026-10-07-150106/REQ-638-handback.md`).

- D-01 (DECIDE & STATE): the tracked-VERSION check stays inline as `slices.ContainsFunc`, not a named helper, because it is used in one place (the REQ's Builder Guidance latitude).
- D-02 (DECIDE & STATE): `suite` and `tools` are added only when absent (`slices.Contains`), not through a set, so a declared source literally named `suite` or `tools` still appears once and the output stays identical for every input.
- D-03 (DECIDE & STATE): the builder did not touch `skills/do-work/CHANGELOG.md` (in write_set) because its brief forbade it; finalization writes the release and both changelogs.

Orchestrator decisions:
- D-04 (DECIDE & STATE): dispatch_at (13:16:29Z, the branch creation in the reflog) and builder_handback_at (13:20:23Z, the builder commit) come from the branch, because the coordinator dispatched the builder before this run claimed the REQ at 15:01:06Z. They precede claimed_at and the estimate. No builder-work timing event was recorded, because the recorder times from a start instant to now. The branch was cut at 13932fc2, before 0.305.70, and was merged as is (no rebase); the gate and heavy lanes ran on the merged tree.
- D-05 (DECIDE & STATE): Route A. One file, exact lines named, equivalence settled by reading both readers; so no Exploration, Scope or Pre-Flight, and `do-work/working/baseline.json` was not rewritten by this run.
- D-06 (DECIDE & STATE): review findings M1 to M3 stay report only (all impact-negligible). M1 is safe today and matches the old code's slice reuse; fixing it would cost a re-merge and a second gate and drain for no behaviour change. No lesson-satellite entry and no index refresh, because there was no incident.
- No ESCALATE decision. No pending-answers follow-ups.

## Discovered Tasks

From the review, impact-stamped per review-work Step 10; none is impact-critical, so nothing was queued.

- No test pins the "modules declared, but no declared source has a tracked VERSION" exit of the release guard (review M2). — impact-negligible → report only
- The release-guard refusal tests assert only the `RELEASE-WITHOUT-SHIPPED-CHANGE` token, not the shipped-root list in the message (review M3). — impact-negligible → report only
- `roots := shippedSources` reuses the declared-sources slice before `append` and `sort` (review M1). — impact-negligible → report only

## Timing

Observed 2026-10-07T15:02:00Z to 2026-10-07T15:09:17Z: 7m 17s total, 7m 02s attributed across 3 events, 15s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 6m 32s | 2 |
| handback-merge | 30s | 1 |

Slowest stage: verification-gate / heavy drain, 4 lanes at 6125e58c, 3m 59s, outcome success.
