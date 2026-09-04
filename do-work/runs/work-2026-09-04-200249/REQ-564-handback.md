# REQ-564 Hand-Back — Reuse Matching Per-Lane Verification Evidence for Four Hours

Branch: `worktree-agent-REQ-564-reuse-matching-per-lane-verification-evidence-for-four-hours`
Commits on the branch (oldest first):

- `27b74c8` `[REQ-564] WIP: heavy-lane evidence reuse, builder interrupted` — not mine; an
  interrupt watchdog committed my working tree mid-build. The content is the same work.
- `4526ab9` `[REQ-564] reuse matching per-lane heavy evidence for at most four hours`
- `19103f5` `[REQ-564] pin heavy-lane evidence reuse at the command seam`

## What Changed, In One Paragraph

A heavy lane that finishes green now stores its result in Git-private state, next to the
green-gate records. A later run reports that record instead of running the lane only when
two independent conditions both hold: the lane's deterministic fingerprint still matches,
and the record is under four hours old. Either condition alone forces the lane to run. The
fingerprint covers the lane's argv, the committed bytes of every path its manifest coverage
declares, the output of every declared toolchain probe, and every declared required
environment variable. Anything that cannot be determined fails closed to running the lane.
Every lane in a run now reports a `disposition` (`executed` or `reused`) plus the exact
condition that decided it.

## File Manifest

Implementation:

- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_evidence.go` — new. The
  Git-private evidence store, the fingerprint document and digest, and the reuse decision.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go` — `RunLanes`
  now takes a `LaneRunRequest`, decides reuse per lane before executing, stamps each lane's
  disposition, and records a green lane's result.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification.go` — the
  manifest lane gains an optional `fingerprint` block, validated at decode.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go` — parses
  `--no-evidence-reuse` and passes the reuse decision through to `RunLanes`.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` — five new
  per-lane fields (`disposition`, `disposition_reason`, `fingerprint_sha256`,
  `evidence_revision`, `evidence_recorded_at`), the matching text rendering, and a normalize
  default of `executed`.
- `_dev/tests/heavy-lanes.json` — every one of the six shipped lanes now declares its
  toolchain probes and required environment variables.

Tests:

- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_evidence_test.go` — new,
  eight tests.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification_test.go` —
  a decode-refusal table for unusable fingerprint declarations, plus an assertion that every
  shipped lane declares toolchain probes.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands_test.go` — the
  `--no-evidence-reuse` parse contract.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` — updated for
  the new rendering, and extended with a reused lane. This is a cross-REQ test change: the
  old assertions pinned a lane line that had no disposition on it. The behavior change is
  intentional (the request asks for the disposition to be recorded), so the assertions were
  updated rather than the rendering reverted.

Prose (the runner's contract is restated in three places outside the code):

- `skills/do-work/actions/work.md` — heavy-lane drain paragraph.
- `skills/do-work/actions/clarify.md` — Step 2.5.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` — the `heavyverification` bullet.

Nothing under `do-work/` in the worktree was written, and the main tree at
`/home/user/skill-do-work` was not written except this file.

## Verification

Commands run from `skills/do-work/tools/do-work-cli` with the Go 1.26.1 toolchain on PATH.

- `go vet ./...` — clean.
- `gofmt -l ./internal` — no output.
- `go test -count=1 ./...` — every package ok.
- `GOOS=windows GOARCH=amd64 go build ./...` and
  `GOOS=windows GOARCH=amd64 go vet ./internal/heavyverification ./internal/ownedprocess` — clean.
- `bash _dev/tests/contract-regressions.sh` from the repository root — exit 0, no FAIL lines
  (see Discovered Tasks for a flake seen on an earlier run of it).
- End-to-end against the real repository: the committed manifest still decodes and plans —
  `plan-heavy-verification --base-revision HEAD~1 --target-revision HEAD` selected four
  lanes with their reasons. I did not run a real heavy lane; each is roughly ten minutes and
  the orchestrator runs the suite after merging.

### Revert and show red

Each behavior has its own mutation, because one blanket revert cannot isolate two conditions
that must hold independently. Every mutation below was applied alone, the named tests were
run, then the file was restored with `git checkout --` and the suite re-run green. The driver
script is `revert_matrix.sh` in the session scratchpad.

| # | Behavior removed | Tests that went RED |
| --- | --- | --- |
| R1 | No lane is ever reused (the behavior before this request) | `TestRunLanesReusesMatchingEvidenceWhileChangedLanesRerun`, `TestRunLanesRerunsExpiredEvidenceAndReuseNeverExtendsTheWindow`, `TestRunLanesWithoutEvidenceReuseExecutesAndStillRefreshesTheRecord` |
| R2 | Age alone authorizes reuse (fingerprint not compared) | `TestRunLanesReusesMatchingEvidenceWhileChangedLanesRerun`, `TestRunLanesRerunsWhenToolchainOrEnvironmentChanges` (both subtests) |
| R3 | A matching fingerprint alone authorizes reuse (no age ceiling) | `TestRunLanesRerunsExpiredEvidenceAndReuseNeverExtendsTheWindow` |
| R4 | A reuse restamps the record, so the four-hour window never ends | `TestRunLanesRerunsExpiredEvidenceAndReuseNeverExtendsTheWindow` |
| R5 | A lane with no declared fingerprint inputs gets an empty fingerprint | `TestRunLanesNeverReusesWithoutADeterminableFingerprint` |
| R6 | The shipped manifest declares no fingerprint inputs | `TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint` |
| R7 | A red or skipped lane is stored as reusable evidence | `TestRunLanesStoresNoEvidenceForRedOrSkippedLanes` |
| R8 | The stored record's own fields are trusted without validation | `TestRunLanesRefusesTamperedEvidenceInsteadOfTrustingIt` |
| R9 | A fingerprint declaration is not validated at manifest decode | `TestDecodeManifestRefusesUnusableFingerprintDeclarations` (all five subtests) |
| R10 | Evidence reuse is off unless the caller asks for it | `TestRunArgumentsDefaultToEvidenceReuse` |
| R11 | The rendered lane line no longer states its disposition | `TestHeavyVerificationRunTextAndJSONCarryTheSameTypedLanes` |
| R12 | The command handler ignores the parsed reuse flag | `TestRunHeavyVerificationCommandReusesEvidenceThroughItsOwnSeam` |

R3 and R4 both redden the same test on purpose: they are the two halves of "four hours is a
ceiling". R3 removes the ceiling; R4 keeps the ceiling but lets each reuse push it forward,
which reaches the same unsafe end by a different route. One differential could not tell them
apart.

After the last restore: `go test -count=1 ./internal/heavyverification ./internal/resultmodel`
ok, and `git status --porcelain` empty.

## Decisions

**D-01 — The manifest stays at `schema_version: 1`; the `fingerprint` block is optional.**
DECIDE & STATE. A version bump would force the decoder to accept a version range, because
the historical-revalidation planner reads the manifest committed at older revisions. An
optional field decodes at every revision with no range logic, and a lane without the block
simply never reuses, which is the correct fail-closed default. Reversible: bumping later is
a two-line change.

**D-02 — The fingerprint's file half is the lane's own declared coverage at the execution
revision, not the whole tree.** DECIDE & STATE. Hashing the whole tree would move every
lane's fingerprint on every commit and reuse would never fire. The coverage rules are already
the manifest's statement of what changes affect a lane, and the selector already trusts them.
Residual risk, stated plainly: a lane whose real inputs exceed its declared coverage can
inherit a green across a change to an undeclared input. That risk is identical to the
selection risk REQ-563 already carries, it is bounded to four hours, and the manifest file
itself is in every lane's coverage, so tightening a declaration invalidates every record.

**D-03 — Toolchain and required environment are declared per lane in the manifest, and a
lane that declares none never reuses.** DECIDE & STATE. The request names toolchain and
environment as fingerprint inputs, and neither can be derived from a lane's argv. A declared
probe list is auditable and deterministic. An absent declaration means the toolchain is
undetermined, which the request says must fail closed — so it does, reported as
`fingerprint_uncertain`, not silently treated as "no toolchain matters".

**D-04 — Reuse is on by default; `--no-evidence-reuse` forces execution.** DECIDE & STATE.
The whole value of this request is the drain skipping lanes nothing changed under, and an
opt-in flag would mean the drain in `work.md` never gets it. The Go-level default is the
opposite (`LaneRunRequest.EvidenceReuse` zero value executes everything), so an in-process
caller can never reuse by forgetting a field.

**D-05 — The reuse fields stay in the run result and are not added to publication's durable
`heavy_testing` evidence.** DECIDE & STATE. Publication's answer manifest decodes with
`DisallowUnknownFields`, and its "confirmed only when every selected command exits zero"
rule would need to be re-judged for an inherited green. That is a separate decision about
durable request evidence, and REQ-547's builder is working nearby. Listed under Discovered
Tasks and Integration Seams. Consequence to know: an archived REQ's heavy evidence still
says only that the lane was green at that execution revision; the disposition lives in the
run output the action reports, not in the REQ file.

**D-06 — A reused lane reports `wall_seconds: 0` and names the inherited record separately
in `evidence_revision` and `evidence_recorded_at`.** DECIDE & STATE. Reporting the original
duration would put a measured-looking number on a run that spent no time, and any duration
roll-up would double-count it.

**D-07 — A record is stamped with the run's evaluation instant (run start), not the lane's
finish time.** DECIDE & STATE. On a long drain this makes records slightly older than they
truly are, which shortens the window rather than extending it, and it makes the age
deterministic in tests.

**D-08 — The browser lane declares fingerprint inputs even though a browser found on PATH by
a well-known name is not version-probed.** DECIDE & STATE. `QUEUE_KANBAN_BROWSER` is in the
lane's declared environment, so the explicitly configured engine is covered. Probing the
PATH fallback would mean copying `maintainer-verify.sh`'s well-known-name list into the
manifest, which is exactly the closed enumeration the prime warns goes stale. Residual risk:
a browser upgrade landing inside a four-hour window could be inherited across. Reversible by
deleting one manifest block if that ever bites.

**D-09 — Records live under the Git common directory (`do-work-heavy-lanes/`), not the
worktree.** DECIDE & STATE. `clarify.md` Step 2.5 runs lanes in a detached worktree of the
same repository, and the common directory is what lets that worktree and the main tree share
one cache. It also matches where `gateevidence` already keeps its records.

Nothing here is an ESCALATE. The one judgment I would flag if you disagree is D-02: it is
the decision that decides whether a false reuse is possible at all.

## Discovered Tasks

- The same evidence-reuse shape could apply to the fast repository gate, which measured about
  half of this run's wall clock across nine invocations. The gate already has a per-argv
  green record in `internal/gateevidence`; what it lacks is an input fingerprint, so it
  currently re-runs whenever HEAD moves for any reason. Not built, per the brief's scope
  instruction. impact-noncritical → report only.
- Carry the per-lane disposition into publication's durable `heavy_testing` evidence, so an
  archived REQ shows which of its greens were inherited. Needs a decision about whether an
  inherited green may still submit `confirmed`. impact-noncritical → report only.
- `_dev/tests/contract-regressions.sh` printed `FAIL: SessionStart hook behavior probes
  failed` on its first of three runs in this session and passed cleanly on the two after,
  with exit 0. I touched no hook code. Looks like a flaky probe, worth a look before it
  wastes someone's afternoon. impact-noncritical → report only.

## Integration Seams

- `internal/publication` (`publication_types.go` `HeavyLaneResult`, `answer.go` line ~663)
  would need `disposition` and `evidence_revision` fields, plus a rule for whether an
  inherited green may submit `confirmed`, to carry reuse into durable REQ evidence. Not
  touched — see D-05.
- `internal/finalization` and the journal image-set code: untouched, as instructed. Nothing
  in this change reaches them.
- `do-work/lessons-index.md`: I wrote no new bullet into
  `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, because the writer of a satellite
  bullet must refresh that index row in the same edit and `do-work/` is a path I may not
  write here. The lesson material is in the next section for whoever lands it.
- `RunLanes`'s signature changed from a parameter list to `LaneRunRequest`. The only
  production caller is `handleRunHeavyVerification` in the same package, so no other package
  needs a change.

## Lesson Evidence

Proposed satellite bullets for `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`,
with their family slugs. Each one has a mutation in the table above that proves it.

- `[family: time-is-not-authorization]` REQ-564: a cached-result window has two halves, and
  they fail independently. Removing the age ceiling (R3) and letting each reuse restamp the
  record (R4) both end with a green that is inherited forever, by different routes, and both
  redden the same test — so a single revert differential cannot isolate them. Write the
  record's stamp from the run that actually executed the lane and never touch it again on
  reuse, and prove the ceiling and the non-extension with separate mutations.
- `[family: fail-closed-declaration]` REQ-564: an absent declaration and an empty declaration
  are different statements. Treating a lane's missing `fingerprint` block as an empty one
  (R5) silently converts "this lane's toolchain is unknown" into "this lane's toolchain does
  not matter", and the resulting digest still looks like a fingerprint. Key reuse on the
  declaration being present and usable, and report the uncertainty as its own reason code.
- `[family: opaque-evidence-projection]` (existing family, second sighting) REQ-564: a reused
  green and a measured green are indistinguishable unless the projection says which is which.
  Every lane carries `disposition` plus the exact condition that decided it — including the
  reason it was *not* reused — so a reader never has to infer the decision from silence.

Related existing lesson honored: `[family: cross-action-exception-closure]` (REQ-518) — a
cached-evidence check bound to `HEAD` cannot serve a caller judging another revision. Here
the fingerprint is computed at the run's own execution revision and the record names the
revision it was measured at, so a reused lane's claim rests on a revision the reader can see.
