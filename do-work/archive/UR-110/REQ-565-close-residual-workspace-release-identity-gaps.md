---
id: REQ-565
title: '[impact-critical] Review fix: Close residual workspace release identity gaps'
status: completed
created_at: 2026-09-03T23:26:06Z
user_request: UR-110
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-512]
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
review_generated: true
addendum_to: REQ-512
sweep: true
sweep_key: legacy-finalization-workspace-identity-residual
claimed_at: 2026-09-05T00:41:15Z
route: A
estimate:
  p50_active_minutes: 10
  confidence: medium
  calculated_at: 2026-09-05T00:50:53Z
  basis:
    - Route A
    - 2-file write set
    - 4 acceptance criteria
completed_at: 2026-09-05T08:05:36Z
commit: 532183cc
release_at: 2026-09-05T08:05:36Z
---

# Close Residual Workspace Release Identity Gaps

## What

Finish the fail-closed workspace ownership boundary left by REQ-512's post-remediation re-review: prove Cargo and uv source identity is unique, and require every structurally present npm root version mirror when the root source changes.

The fold-first scan found no pending or pending-answers REQ in any UR that shares this exact residual finalization workspace identity root cause.

## Context

REQ-512 completed its one permitted remediation pass and closed changed-source-first selection, exact shared-lock replacement, malformed missing/unquoted identity, malformed npm JSON, and bounded fold termination. Independent re-review found the two remaining cases below, so this successor owns them instead of silently taking a second remediation pass.

## Instances

- [ ] `finalization_discovery.go:tomlSectionScalar` accepts the first Cargo/uv `name` without proving uniqueness or TOML identity validity.
- [ ] `npmRootVersionCopies` counts only root lock values already equal to `oldVersion`, so parseable stale or mismatched root copies can disappear from the required mirror set.

## Requirements

- Refuse a changed Cargo or uv source when its applicable package/project identity is duplicated, competing, ambiguous, or otherwise not uniquely parseable.
- When a changed npm root has a tracked parseable lock, treat every structurally present root `version` and `packages[""].version` copy as an obligation; a stale or mismatched copy must refuse rather than be omitted.
- Preserve member-only releases, pre-existing target-version neighbors, multiple changed members sharing one lock, typed enumeration/ownership failures, protected-path refusal, and public recovery-to-claim success.
- Add committed RED/GREEN tests for duplicate/competing Cargo and uv names plus one-stale and both-stale npm root copies.

## Red-Green Proof

**RED prompt/case:** Run strict finalization discovery for (a) a changed Cargo/uv source with duplicate or competing applicable `name` declarations and (b) a changed npm root whose parseable lock has one or both structurally present root copies stale.

**Why RED now:** The TOML helper returns the first name, while npm root-copy counting ignores present values that do not already equal the expected old version; both can silently reduce required mirror ownership.

**GREEN when:** Every ambiguous source identity and stale structurally present npm root copy refuses before mutation with typed path evidence, exact valid mirrors still finalize, and the full REQ-512/recovery matrix remains green.

**Validation:** REQ-512 post-remediation re-review; successor required by the one-remediation rule.

## Full Context

See `do-work/user-requests/UR-110/input.md` and REQ-512's `## Re-Review` section.

## Open Questions

- [x] Auto-approved: critical severity (release/finalization ownership risk). → Added to queue immediately.

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names both gaps by function, states the fail-closed rule each must satisfy, and lists the four RED cases its tests must cover. The predecessor REQ-512 supplies the surrounding design.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req565_test.go` (new)

**What was done:** Two release-identity helpers in finalization discovery were made fail-closed: Cargo and uv source identity must now be declared exactly once across the accepted sections or the release refuses, and a changed npm root now owes every structurally present root version copy instead of only the copies that already hold the old version. A new test file locks in both gaps with committed RED/GREEN cases.

The first gap was a first-match guess. The TOML scalar helper returned the first name line it met inside an accepted section and reported success, so a manifest that declared the key twice, repeated its table, or disagreed between the two pyproject dialects released under whichever identity was listed first. The other declared name's lock entry stayed on the old version and nothing refused. The helper was renamed to `tomlUniqueSectionScalar`, now proves uniqueness, and returns an error. `releasePackageName` returns that error, and its one caller wraps it.

The second gap was counting by value instead of by presence. `npmRootVersionCopies` incremented only when a root lock value already equalled the old version, so a present root `version` or `packages[""].version` holding any other value lowered the required mirror count instead of refusing. Two shapes were observed. With one stale copy, the required count fell from two to one, the single matching copy satisfied the mirror check, and the release landed with the top-level value still on 0.9.0. With both copies stale the count was zero, so the lock never entered the mirror set or the committed path list at all: the obligation vanished rather than failing.

Both new failure paths surface through the existing typed release-enumeration refusal, before any mutation, under reason code `FINALIZATION-DISCOVERY-RELEASE-ENUMERATION`. The refusal renders as `Outcome: refused`, `Phase: discovery_refused`, error severity, `Fixability: refused`, with the dirty release metadata paths as affected paths and the reason text as the first evidence line. Identity refusals read either as a competing-declarations count, as a named unquoted or malformed declaration, or as the pre-existing absent-declaration message. Lock refusals name which root copy is stale and what value it holds against the released version. Every listed condition returns a copy count of zero together with the error, so no partial count can escape.

Nothing outside the finalization package was touched, and the worktree was clean after the commit.

One correction to the hand-back: its closing note reports the diff as 69 insertions and 23 deletions in the source file plus one new test file. Git reports 92 insertions and 23 deletions in the source file and 155 insertions in the new test file, 224 insertions total across the two files. The file manifest itself matches exactly.

**Implementation range:** `3e988f39..532183cc`. Builder commit `2d715fd940b56a86000385bdfb0765c6fdb458b2`.

## Decisions

- **D-01 — "duplicated" means refuse, even when the two names agree:** The request lists "duplicated, competing, ambiguous, or otherwise not uniquely parseable". The builder implemented the count rule literally: exactly one declaration across the accepted sections, or refuse. A pyproject that declares the same name in both the project and the poetry sections therefore refuses too. The reasoning is that TOML forbids a duplicate key outright, so the document is invalid either way, and a fail-closed refusal on a transitional poetry config is cheap to fix by hand while a first-match guess is silent. The builder flagged this as the one place a reviewer may want a different call.
- **D-02 — refuse on any declaration that is present but not a quoted scalar, not only on an absent one:** The old code skipped a non-matching malformed name line and reported "no name", which happened to refuse anyway. A present-but-malformed declaration now refuses on its own terms with the offending line quoted, so the operator sees which line is wrong. This keeps REQ-512's malformed uv project name test green while making the evidence actionable.
- **D-03 — renamed the TOML scalar helper to `tomlUniqueSectionScalar`:** The contract changed from "first match wins" to "exactly one, or an error", and a reader meeting the old name elsewhere would assume the old semantics. It had exactly one caller, so the rename cost nothing and makes the contract findable by search. It is consistent with `projectLockVersion`, which already proves uniqueness through `exactSingleVersion`.
- **D-04 — `releasePackageName` returns an error, not a boolean:** The caller previously flattened every identity problem into one message. The project's lessons file records the trap of one shared remedy line serving conditions that need different remedies, from REQ-461 (requiring affirmative project-owned release targets). The wrapped error names the actual condition — competing, malformed, or absent — inside the same typed reason code, so no consumer contract changed.
- **D-05 — presence, not value, is the npm root obligation, decoded through `json.RawMessage`:** Decoding into a string cannot tell an absent key from an empty one, which is exactly the confusion that produced the bug. `json.RawMessage` distinguishes them, so "structurally present" is decided by the key existing and the value is validated separately. `workspaceMirrorReplacement` needed no change, because every counted copy is now guaranteed to equal the old version, so its replacement-count guard still balances.
- **D-06 — kept one passing positive control in the new npm test:** The subtest `only_the_top-level_root_copy_is_present` passed in RED and still passes in GREEN. It exists to fail if someone later simplifies the rule to "a changed npm root always owes two copies", which would falsely refuse a lock with no `packages[""]` entry.
- **D-07 — RED fixtures were built so the broken code succeeds, not so it refuses differently:** The builder's first sketch would have produced a refusal with the wrong reason code in RED, which proves little. Each fixture instead updates only the mirror the broken code selects, so RED reaches a successful outcome with a real commit. That is the actual hazard the request describes.

## Qualification

Passed the request-bound advance qualify gate for `3e988f39..532183cc`. Two files, both inside `internal/finalization/`, nothing in the packages other builders held concurrently. Independent review reproduced the RED in a scratch copy and confirmed the sharpest claim: on the pre-change code all six refusal fixtures returned success with a real commit, and in the both-stale npm case the lockfile dropped out of the committed path list entirely. The P-A-U boxes were reconciled from the builder hand-back.
## Testing

**Red-green validation:** RED was captured before any source change, from `go test -count=1 -run 'TestREQ565' ./internal/finalization`. All six refusal subtests reported the same thing: the broken code committed the release instead of refusing.

```
--- FAIL: TestREQ565AmbiguousWorkspaceSourceIdentityFailsClosed (0.00s)
    --- FAIL: .../cargo_declares_the_package_name_twice (2.03s)
        finalization_req565_test.go:76: release identity gap did not fail closed: Outcome:"success" (finalization committed)
    --- FAIL: .../cargo_repeats_the_package_section (1.76s)
    --- FAIL: .../uv_declares_the_project_name_twice (2.06s)
    --- FAIL: .../uv_project_and_poetry_sections_compete (1.93s)
--- FAIL: TestREQ565ChangedNPMRootRequiresEveryPresentRootVersionCopy (0.00s)
    --- FAIL: .../top-level_root_copy_is_stale (2.31s)
        finalization_req565_test.go:132: release identity gap did not fail closed: Outcome:"success" (finalization committed)
    --- FAIL: .../both_root_copies_are_stale (1.77s)
```

The unabbreviated dumps show what "did not fail closed" meant in practice. For `cargo_declares_the_package_name_twice` the result carried `Outcome:"success"`, `Phase:"cleanup_complete"`, `PrimaryCommit:"cc48558d…"`, and committed both the manifest and the lock while the second declared name's lock entry stayed on 1.0.0. For `top-level_root_copy_is_stale` the committed path list contained the npm lock with one copy still stale. For `both_root_copies_are_stale` the committed path list omitted the lock entirely, which is why both npm cases exist: the two failure shapes are different.

GREEN is the same command on the committed tree, covering both new test functions:

```
$ go test -count=1 -run 'TestREQ565' ./internal/finalization
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/finalization	2.602s
```

`TestREQ565AmbiguousWorkspaceSourceIdentityFailsClosed` pins the four identity cases. `TestREQ565ChangedNPMRootRequiresEveryPresentRootVersionCopy` pins the two stale-copy cases plus `only_the_top-level_root_copy_is_present`, which passed in RED and still passes: it is the guard against over-tightening, since a lock with no `packages[""]` entry owes exactly one copy, not two.

**Controls preserved:** Nine named must-not-regress tests were re-run at the post-fix commit and all passed, in one run reported as `ok … 14.929s`:

- `TestREQ512WorkspaceMembersSelectChangedSourcesBeforeEqualVersionRoots` and `TestRecoverFinalizationRequiresWorkspaceMemberLockMirrors` protect member-only releases across npm, cargo and uv; the first asserts the unchanged root is absent from the committed path list and byte-identical.
- `TestREQ512SharedWorkspaceLocksReplaceOnlyMultipleChangedMembers` protects two things at once: pre-existing neighbours already at the target version, and multiple changed members sharing one lock.
- `TestRecoverFinalizationReleaseEnumerationFailureIsTypedAndFailClosed` and `TestREQ512MalformedWorkspaceIdentityAndNPMRootLockFailClosed` protect typed enumeration failures; the three REQ-512 malformed identity and lock cases still emit the same reason code through the reworded errors.
- `TestRecoverFinalizationRefusesReleaseMetadataWithoutProjectOwnership` protects the typed ownership failure.
- `TestRecoverFinalizationAssumeSoleReleaserStillRefusesStagedProtectedPath` protects protected-path refusal.
- `TestREQ512ChangedWorkspaceRootsStillRequireRootLockMirrors` protects the rule that changed roots still owe their root lock mirrors.
- `TestRecoverFinalizationRefusesPartialConfiguredReleaseMirrors` protects refusal on partial configured mirrors.

The tenth control is the heavy lane. `TestPublicRecoverFinalizationMovesURThenAllowsRealClaim` protects public recovery-to-claim success and was run separately with the heavy-test environment variable set, passing in 3.70s (`ok … 3.743s`).

**Module verification:** Four commands, run from the CLI module directory inside the worktree.

```
$ go test -count=1 ./internal/finalization
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/finalization	56.203s

$ gofmt -l .
(no output)

$ go vet ./...
(no output, exit 0)

$ go test -count=1 ./...
FAIL	github.com/knews2019/skill-do-work/do-work-cli/internal/heavyverification	41.273s
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/finalization	77.598s
FAIL
```

Every other package in the full run was `ok`, including corehelpers (44.666s), doctor (5.712s), gateevidence (6.726s), gittransaction (7.931s), knowledgecommands (20.973s), lifecycleadvance (26.677s), nextselection (7.100s), publication (34.699s), requeststate (5.064s), suiteinstall (8.344s) and toolboxcommands (4.427s).

The `internal/heavyverification` failure was diagnosed as not belonging to this change. Two tests failed there, `TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes` (subtests GOOS, GOARCH, CGO_ENABLED, GOENV and DO_WORK_TEST_DO_WORK_CLI_BINARY, reporting "default runtime must have a determinable fingerprint") and `TestShippedGitIsolationPreservesGenericLaneInheritance` ("shipped runtime probe did not isolate host Git configuration"). That package does not import `internal/finalization`, proven by a grep that returned nothing. The builder then stashed the source change and re-ran only those two tests on the untouched baseline, getting identical failures in 4.182s before restoring the change. It reads as a sandbox limitation of this worktree, not a regression from this work.

## Discovered Tasks

- **DT-1 — the version twin of the gap just closed is still open.** `tomlSectionVersion` at `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go:776` still returns the first version it meets in an accepted section. A Cargo or uv manifest with two version declarations, or with the project and poetry sections declaring different versions, releases from whichever is listed first. It is narrower than the identity gap, because `semanticVersionReplacement` independently requires exactly one occurrence of the old version string in the manifest bytes for the default branch, which catches many but not all shapes. Rated impact-medium by the builder. → queue as follow-up
- **DT-2 — the npm identity branch does not prove uniqueness.** The npm branch of `releasePackageName` at `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go:1005` decodes the name with the standard JSON decoder, so duplicate name keys in a package manifest silently resolve to the last one rather than refusing. Cargo and uv now prove uniqueness; npm does not. `jsonObjectValueLocation` already refuses duplicate keys at replacement time, so the practical exposure is limited to which lock entry gets selected. Rated impact-low by the builder. → report only

## Review

**Overall: 92%**
**Acceptance: Pass, approve with follow-ups.** The reviewer reproduced the RED independently in a scratch copy and confirmed the sharpest claim: on the pre-change code all six refusal fixtures returned `Outcome:"success"` with a real commit, and in the both-stale npm case the lockfile dropped out of the committed path list entirely — the release obligation vanished rather than failing.

It verified the six must-not-regress behaviours by reading each named test before running it, rather than accepting the table, and confirmed the container's `heavyverification` failures are environmental by running the same command through an environment-stripping wrapper and getting `ok`.

One Important finding: identical duplicate declarations render as "2 competing name declarations (alpha, alpha)", where "competing" contradicts the evidence beside it and no remedy is given. The project's own lessons file records this exact trap, and the builder cited that lesson without applying it here.

A second-order effect worth knowing: a `pyproject.toml` declaring the same name in both `[project]` and `[tool.poetry]` is now refused even though the names agree. That shape appears in projects mid-migration. It is on-spec and fails in the safe direction, but a consumer with a transitional manifest cannot release until they hand-edit it.

## Lessons Learned

Two rules came out of this, both above the level of the files touched.

When a rule says "every present copy is an obligation", decode presence and value separately. Decoding a JSON field straight into a string cannot tell an absent key from an empty or wrong one, so a present-but-wrong copy silently lowers the required count instead of refusing, and an obligation can leave the required set entirely without any failure. Decode into a raw message first, decide presence from whether the key exists, then validate the value. The same reasoning applies to any counter that increments only on a match: ask whether a non-match should lower the count or refuse.

A RED fixture for a fail-closed rule has to make the broken code succeed, not refuse for a different reason. A refusal with the wrong reason code in RED proves almost nothing, because the code was already refusing. Build the fixture so the broken path reaches a successful outcome with a real commit, which is the actual hazard, and assert on the outcome rather than on the message.

## Orientation

Finalization discovery now refuses before any mutation when a changed Cargo or uv source declares its identity more than once, and a changed npm root now owes every structurally present root version copy rather than only the ones already holding the old version. The failure mode where a release quietly committed under the wrong identity, or dropped a lockfile out of its own obligation set, is now a typed refusal with the offending path and value named.
