# Re-review: REQ-498

**Request changes** — the remediation closes the journal, phase, protected-state, multi-group, and typed-result defects, but legacy discovery can still finalize incomplete release state or absorb foreign changes from an already-tracked follow-up. The central “recover only unambiguous state” acceptance condition therefore remains open.

Route C | cumulative merge range `e8e5a79d..1249e856` | reviewed 2026-09-02T16:22:06Z

## What's built

- Planned and discovered finalizations now share durable `verified` and `cleanup_complete` phases, exact image integrity, created-versus-settled commit evidence, and full descendant commit recognition.
- The remediated behavior suite proves a complete REQ-494-shaped fixture, enriched checkpoint removal, calibration/UR/follow-up/release grouping, protected staged-versus-unstaged behavior, two-group ordering, all durable phase interruptions, primary/supplied provenance, failure completion, corrupt-image refusal, and normalized success/refusal records.
- Startup and ordinary commit preflight invoke the same recovery authority before queue selection or loose-change association.

## Decisions / risks for the orchestrator

- Do not mark the semantic legacy-discovery finding closed. `associateReleaseMetadata` validates only the release files that happen to be dirty and knows no required-mirror set; a partially applied release can be committed and cleaned up while configured mirrors remain stale.
- The one-remediation limit has been consumed. Route the remaining semantic-ownership and acceptance-matrix gaps through the work action's post-remediation failure path rather than attempting another in-place remediation.

## Findings

### Important

1. **Legacy semantic ownership still accepts incomplete or foreign shared state.** `finalization_discovery.go:373-445` derives a release group only from currently dirty paths, while `releaseMetadataPath` at `:459-461` hard-codes the root/suite files and never proves the publication contract's `required_mirrors`. A crash with only `CHANGELOG.md` and `VERSION` updated therefore satisfies the one-version/one-insertion test, is committed as complete, and leaves configured mirrors stale; consumer release sources such as `package.json`, `Cargo.toml`, or lockfile mirrors are not release-semantic paths at all. Separately, `followupPathProves` at `:317-324` accepts any dirty tracked `do-work/**` document whose current `addendum_to` matches, without proving that it was newly created or that its whole diff is the originating follow-up mutation, so foreign edits in an existing follow-up can ride in the exact commit. Both violate the no-guess/whole-diff boundary that motivated remediation. — `impact-critical`

2. **Operative action prose still restates the retired direct tail.** `work.md:120` tells the orchestrator that it handles all moves, frontmatter updates, and archiving, despite Step 9 assigning those mutations exclusively to finalization. `commit.md:7` still says work's Commit Phase delegates to the request-state command and commits the planned paths itself, including a separate serial metadata commit. These are active instructions, not historical commentary, and contradict the single-finalizer contract established later in both actions. — `impact-user-visible`

3. **The explicit acceptance/TDD matrix remains incomplete.** The new tests are substantial, but no finalization test simulates a real pre-commit hook failure and resumes the retained journal; no test covers the already-green repair no-op path or a planned release manifest; and no negative fixture covers partial required mirrors or foreign edits in a tracked follow-up. The git-transaction package's hook test validates the lower-level authority, not the requested end-to-end finalization recovery boundary. — `impact-user-visible`

### Minor

- `work.md` still titles Step 8 “Archive” and says the substeps below archive work even though the step now only prepares intent. This reinforces finding 2 but is not scored separately.

### Nit

None.

## Prior Important Finding Closure

| Prior finding | Result | Evidence |
|---|---|---|
| Complete semantic legacy ownership and foreign-hunk refusal | **Open** | The full suite-specific fixture and changelog/project multi-hunk refusals pass, but partial required mirrors and dirty tracked follow-ups remain admissible. |
| Multi-group earlier-commit recovery | **Closed** | `matchingHeadCommit` scans every ancestry-path descendant and verifies cumulative exact diff plus candidate-local allowlisted paths; `TestMatchingHeadCommitSearchesPastEarlierIndependentGroup` passes. |
| Protected unstaged versus staged behavior | **Closed** | Discovery excludes unstaged `X`/`XD` rows, distinguishes protected-staged from foreign-staged reason codes, and the focused behavior test passes. |
| Complete ordered terminal/refusal evidence | **Closed** | Durable terminal phases, created/settled hashes, normalized slices, singular/plural refusal records, and exact argv fields are implemented and tested. |
| Acceptance/TDD matrix | **Partially closed** | Phase, provenance, corruption, semantic-tail, ordering, and idempotence coverage landed; hook/no-op/planned-release and the remaining ownership negatives did not. |
| Stale action restatements | **Partially closed** | The originally cited Step 8/9/checklist passages were corrected, but contradictory operative statements remain in `work.md:120` and `commit.md:7`. |

## Requirements Checklist

- [x] Strict `finalize --manifest` and `recover-finalization --discover` commands share one Git-private journal/phase engine.
- [x] Journal phases include prepared, lifecycle, release, primary, metadata, verified, cleanup-complete, and typed discovery refusal.
- [x] Canonical request-state, publication, protected-inventory, and exact Git-transaction authorities are reused.
- [~] Retry/idempotence is proved across durable phases and the complete fixture, but incomplete release mirrors can be accepted as settled state.
- [x] Primary provenance and validated supplied-worktree provenance are distinct and tested.
- [x] Startup recovery precedes checkpoint/working recovery, selection, and claim; ordinary commit preflight delegates first.
- [ ] Legacy discovery admits only fully proven semantic shared state and refuses every foreign hunk.
- [~] Unrelated ordinary/protected unstaged changes are preserved and staged state refuses, but a dirty tracked follow-up can be over-associated.
- [x] Ordered typed records expose terminal phases, created/settled commits, blockers/reasons, and non-null command arrays.
- [~] Work and commit actions delegate operationally, but their operative overview/restatement text still contradicts that ownership.
- [x] Existing lifecycle/release commands remain registered and the one-releaser model is unchanged.
- [~] The REQ-494-shaped positive flow is demonstrated; the original safety and hook/no-op acceptance matrix is not complete.

## Acceptance Testing

**Result: Fail**

- The focused remediation matrix passed independently: `go test -count=1 ./internal/finalization -run '<12 remediation cases>'` — exit 0 in 8.365s.
- The full finalization package passed independently: `go test -count=1 ./internal/finalization` — exit 0 in 10.774s.
- Source-level negative acceptance failed: the release association has no required-mirror completeness check, and tracked follow-up ownership has no preimage/diff-shape check. Both paths flow into `effectivePaths`, journal creation, and exact commit as successful discovery.
- An independent `bash _dev/tests/contract-regressions.sh` rerun reached the finalization/action predicates but exited 1 in an unrelated update-script behavior probe (`printf: write error: Broken pipe`; expected “four-module suite” text was present in the captured output). The orchestrator's previously recorded clean contract and maintainer-gate runs remain the authoritative repository evidence; this rerun does not explain or excuse the acceptance defects above.

## Suggested Additional Testing

- Add a partial-release fixture where only a subset of required version/changelog/lockfile mirrors is dirty; require byte-identical refusal rather than successful cleanup.
- Modify an already-tracked follow-up carrying `addendum_to: REQ-NNN` with an unrelated hunk and require refusal; retain a positive case for a newly created originating follow-up or an exactly proved fold.
- Install a failing pre-commit hook during primary finalization, then remove/fix it and prove the same journal resumes without bypass, duplicate commits, or lost implementation bytes.
- Exercise the already-green/no-release manifest and a planned release manifest through every interruption boundary, not only discovered release state.

## Scores (on the record — not the headline)

**Overall: 50%**

| Dimension | Score | Notes |
|---|---:|---|
| Requirements | 73% | Most transaction/recovery mechanics are complete; the core unambiguous legacy ownership boundary is still violated. |
| Code Quality | 68% | The phase engine and projections are strong, but semantic association contains unsafe incomplete proofs. |
| Test Adequacy | 72% | Broad remediation coverage, with explicit hook/no-op/planned-release and ownership negatives still missing. |
| Scope | 100% | All 14 implementation files match the declared scope; cumulative `do-work/` artifacts are orchestrator bookkeeping, not builder drift. |
| Risk | Critical | Automatic recovery can canonically commit partial release state or foreign tracked follow-up edits. |
| Acceptance | Fail | Safe legacy-tail recovery is a central requirement, and its negative boundary still fails by inspection. |

Raw percentage average: 78%. Critical risk caps at 60%; Acceptance Fail caps the final score at 50%.

## Follow-up Disposition

No queue, request, checkpoint, or implementation state was modified. Per the re-review brief, this reviewer wrote only this report; the orchestrator owns post-remediation follow-up routing.

## Self-validation

- Re-read UR-096, the complete REQ, initial failed review, original and remediation hand-backs, and cumulative range `e8e5a79d..1249e856`.
- Walked every detailed requirement and every prior Important finding against production predicates and behavior fixtures rather than relying on the hand-back.
- Applied the restatement sweep across `work.md`, `work-reference.md`, `commit.md`, the CLI prime, and executable contract predicates.
- Verified all P-A-U boxes are checked and the cumulative implementation-only diff contains exactly the 14 declared files.

*Re-reviewed independently with the review-work action*
