---
id: REQ-441
title: '[impact-critical] Validate HTTP archives before publication'
status: claimed
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-414]
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-437, REQ-438, REQ-439, REQ-440, REQ-442, REQ-443, REQ-444]
batch: accepted-feedback-regressions
claimed_at: 2026-09-01T02:34:54Z
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-09-01T02:42:22Z
  basis:
    - Route B
    - 2-file builder write set
    - 5 acceptance criteria
    - dependency depth 1
    - asynchronous publication behavior
    - cross-route regression gates
    - full-suite verification
---

# Validate HTTP Archives Before Publication

## What

Download HTTP archive candidates to a private sibling, validate them there, and only then publish onto an absent or regular non-symlink target. A failed HTTP validation followed by failed Git fallback must leave every pre-existing target byte-identical.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reconciled the superseded shell seam and accepted a shared cross-route publication plan for validated replacement of an unchanged regular target.
- [x] **[APPLY]:** Implemented one rooted cross-route target snapshot, candidate preparation, final revalidation, and absent-versus-regular publication policy in the frozen two-file builder scope.
- [x] **[UNIFY]:** Reviewed both changed files; race/focused/full Go, vet, exact Go 1.25, update/contract, canonical maintainer, and diff hygiene gates pass on the builder branch.

## Finding Provenance

- **Verbatim claim / severity:** `[P1] Validate HTTP downloads before replacing the target.`
- **Evidence:** `archive_fetch.go:117-124` passes `ArchiveTargetPath` directly to `atomic-download.sh` and validates only after the helper has published it.
- **REQ-414 seam reconciliation:** REQ-414 removes the shell-publication call and validates the HTTP body inside Go before no-overwrite publication. Once REQ-414 is terminal, re-triage this REQ against the remaining cross-route target-shape/replacement contract; do not replay the superseded `atomic-download.sh` evidence.
- **Origin / earned by:** The Go port `f27f564d` inherited the behavior from shell commit `0e8cf0d9`. A successful HTTP transfer containing an unreadable tar followed by failed Git fallback replaced a pre-existing archive despite `FetchArchive`'s preservation contract.
- **Surface-cost:** Earned. The destructive replay and symlink/special-target class justify one private stage lifecycle and target-mode guard; tests must preserve an old regular target and refuse non-regular targets.

## Detailed Requirements

- Allocate HTTP staging beside the public target and pass only that private path to the download primitive.
- Validate the staged candidate as a readable archive before publication.
- Before publication, require the public target to be absent or an existing regular non-symlink file.
- On HTTP invalidity and Git fallback failure, preserve a pre-existing target byte-for-byte and remove only invocation-owned scratch.
- Apply a consistent target-shape rule to both HTTP and Git publication paths.

## Constraints

- Preserve the two-route fallback contract and truthful route evidence.
- Do not reimplement archive parsing; retain the existing readability check.
- Keep private staging same-directory so final publication can remain atomic on supported platforms.

## Red-Green Proof

**RED prompt/case:** Seed a regular archive target, make HTTP return success with a non-tar body, and make Git fallback fail; separately use a symlink or special file as the public target.
**Why RED now:** HTTP publishes onto the public target before Go validates the candidate.
**GREEN when:** Total failure preserves the old regular target byte-for-byte with no scratch, non-regular targets are refused unchanged, and a valid HTTP archive still publishes successfully.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm on staging and target preservation. Reuse existing publication mechanics where that reduces, rather than expands, surface.

## Triage

**Route: B** — Medium

**Reasoning:** REQ-414 delivered private staging and validation but narrowed both routes to absent-only targets. The remaining delta is one archivefetch package policy with a bounded cross-route race matrix.

**Planning:** Exploration-guided plan completed in `do-work/runs/work-2026-08-31-165510/REQ-441-plan.md`; repository evidence is in `REQ-441-exploration.md`.

## Plan

1. Replace the occupied-target refusal lock-in with RED fixtures for valid HTTP refresh and invalid-HTTP/valid-Git refresh of a seeded regular target, plus preservation, unsafe-shape, concurrent-change, parent-swap, absent-target, and final-mode cases.
2. Open the target parent once, snapshot absence or regular identity/content, let HTTP and Git prepare validated private candidates under that root, and publish through one shared helper: exclusive link for initial absence or rooted atomic rename after final identity/digest revalidation for an initially regular target.
3. Keep standalone `DownloadAtomic` no-overwrite and hand the exact prime restatement to the REQ-415 builder, which exclusively owns the shared document in this fan-out wave. Run focused/full Go, vet, exact Go 1.25, update/contract, scope/diff, and canonical gates.

## Decisions

### D-01: Restore validated replacement for unchanged regular archive targets

**Decision:** DECIDE & STATE — `FetchArchive` may replace an initially regular non-symlink target only after candidate validation and final identity/content revalidation through one held root. Absent targets remain exclusive-create; changed, missing, symlinked, or special targets refuse unchanged.

**Reasoning:** The request's accepted target set is absent or regular. REQ-414 fixed destructive pre-validation publication but its absent-only policy left successful refresh unimplemented.

### D-02: Keep generic download no-overwrite and candidate mode private

**Decision:** DECIDE & STATE — public `DownloadAtomic` remains no-overwrite; only `FetchArchive` owns validated refresh. A successful archive replacement publishes the private candidate's 0600 mode.

### D-03: Transfer the shared prime edit to REQ-415

**Decision:** DECIDE & STATE — this builder does not edit `prime-do-work-cli.md`; its brief supplies the accepted archive-refresh wording to the REQ-415 builder, the sole owner of that shared file in this wave.

**Reasoning:** Explicit single ownership removes the only fan-out collision while keeping the living authority synchronized before either REQ is reviewed.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch.go` (modified) — shared rooted target snapshot, candidate preparation, and publication.
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch_test.go` (modified) — cross-route replacement, preservation, unsafe-target, and race fixtures.

**Files I will NOT touch:** `prime-do-work-cli.md` (owned by REQ-415 for this wave), retained fetch/download shell paths, `corehelpers`, `atomicfile`, suiteinstall callers, update fixtures, actions, hooks, recipes, release metadata, or REQ-414 residual packages. Expansion requires a focused RED and owner approval.

**Acceptance criteria:**
- [x] Valid HTTP can replace an unchanged regular archive only after validation, with final mode 0600 and no scratch.
- [x] Invalid HTTP followed by valid Git applies the same validated replacement contract and reports the winning route truthfully.
- [x] Total failure preserves pre-existing bytes/mode; absent targets remain exclusive and unsafe/changed targets refuse unchanged.
- [x] Parent swaps and competing creations/replacements cannot redirect or clobber publication.
- [x] Focused/full Go, compatibility, update/contract, scope/diff, and canonical gates pass within exactly two builder files, with the coordinated prime restatement present before review.

## Implementation Summary

`FetchArchive` now holds one opened parent, snapshots either absence or an existing regular target's identity and SHA-256 digest, prepares and validates private HTTP/Git candidates, and sends either winning route through one publication helper. Initial absence remains exclusive-link publication. An initially regular target is atomically renamed over only after final identity/content revalidation; unsafe, changed, replaced, or removed targets refuse unchanged. Standalone `DownloadAtomic` remains no-overwrite.

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch.go` (modified) — rooted target snapshot, candidate preparation, final validation, and shared publication.
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch_test.go` (modified) — HTTP/Git replacement, total-failure preservation, unsafe-shape, race, parent-swap, mode, and generic no-overwrite fixtures.

**Builder commit:** `bab4bbae470d19d150bf790a8292133893d5374b`

**Cumulative builder range:** `eec9ea36..bab4bbae`

**Integrated merge commits:** `1a0d447f`, `e0d66f1f`

The coordinated `prime-do-work-cli.md` restatement remains owned by the in-flight REQ-415 builder and must land before independent review.

*Generated by work action from `do-work/runs/work-2026-08-31-165510/REQ-441-handback.md`.*

## Qualification

- `DO_WORK_DIFF_RANGE=eec9ea36..bab4bbae ... qualify.sh` — passed with no warnings after the final unsafe-target evidence amendment.
- `DO_WORK_DIFF_RANGE=eec9ea36..bab4bbae ... scope-drift.sh` — passed with the exact two-file builder match.

## Testing

- Merged `go test -race -count=1 ./internal/archivefetch` and focused vet — passed after the final correction.
- Merged focused archivefetch, suiteinstall, and corehelpers tests — passed.
- Builder full Go, vet, exact Go 1.25, update behavior, contract regressions, canonical maintainer gate, and diff hygiene — passed per the durable handback.
- REQ-415 merge `a18bf17a` landed the coordinated `prime-do-work-cli.md` restatement before independent review.

## Review — Initial

**Overall: 84%** | 2026-09-01

| Dimension | Score |
|-----------|-------|
| Requirements | 94% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important finding:**
- F1 — impact-user-visible — a target-parent open failure returns a bare error instead of the required truthful HTTP/Git not-attempted evidence envelope and actionable escape-hatch shape.

The rooted publication, regular refresh, preservation, unsafe-target, race, mode, and standalone no-overwrite contracts passed. F1 enters the single remediation pass.

*Reviewed independently; full evidence is in `do-work/runs/work-2026-08-31-165510/REQ-441-review.md`.*

## Remediation

The sole remediation pass committed `0fc171d5781006516007806e6be9699fdcc0c4e3` and remained within the frozen two-file scope. Target-parent open failures now use the shared archive-fetch failure envelope, preserving the concrete error while reporting both routes as not attempted and retaining the standard escape hatch. A real-command text/JSON RED-GREEN fixture proves both route outcomes, zero HTTP requests, and no parent mutation.

**Remediation integration commit:** `7101e21878e4681a4cfd4991cd454adea5777679`

*Generated from `do-work/runs/work-2026-08-31-165510/REQ-441-remediation-handback.md`.*

## Review — Fresh Re-review

**Overall: 98%** | 2026-09-01

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

Initial F1 is closed. Missing and non-directory parent failures now independently prove truthful text/JSON route evidence, the concrete local cause, the standard escape hatch, zero network attempts, and no mutation. No Important, Minor, or Nit findings remain.

*Re-reviewed independently; full evidence is in `do-work/runs/work-2026-08-31-165510/REQ-441-rereview.md`.*

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 20 from the validated external feedback.*
