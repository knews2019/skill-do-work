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
    - 3-file write set
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
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
3. Keep standalone `DownloadAtomic` no-overwrite, update the CLI prime to distinguish the two policies, and run focused/full Go, vet, exact Go 1.25, update/contract, scope/diff, and canonical gates.

## Decisions

### D-01: Restore validated replacement for unchanged regular archive targets

**Decision:** DECIDE & STATE — `FetchArchive` may replace an initially regular non-symlink target only after candidate validation and final identity/content revalidation through one held root. Absent targets remain exclusive-create; changed, missing, symlinked, or special targets refuse unchanged.

**Reasoning:** The request's accepted target set is absent or regular. REQ-414 fixed destructive pre-validation publication but its absent-only policy left successful refresh unimplemented.

### D-02: Keep generic download no-overwrite and candidate mode private

**Decision:** DECIDE & STATE — public `DownloadAtomic` remains no-overwrite; only `FetchArchive` owns validated refresh. A successful archive replacement publishes the private candidate's 0600 mode.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch.go` (modified) — shared rooted target snapshot, candidate preparation, and publication.
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch_test.go` (modified) — cross-route replacement, preservation, unsafe-target, and race fixtures.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — distinguish generic no-overwrite from validated archive refresh.

**Files I will NOT touch:** retained fetch/download shell paths, `corehelpers`, `atomicfile`, suiteinstall callers, update fixtures, actions, hooks, recipes, release metadata, or REQ-414 residual packages. Expansion requires a focused RED and owner approval.

**Acceptance criteria:**
- [ ] Valid HTTP can replace an unchanged regular archive only after validation, with final mode 0600 and no scratch.
- [ ] Invalid HTTP followed by valid Git applies the same validated replacement contract and reports the winning route truthfully.
- [ ] Total failure preserves pre-existing bytes/mode; absent targets remain exclusive and unsafe/changed targets refuse unchanged.
- [ ] Parent swaps and competing creations/replacements cannot redirect or clobber publication.
- [ ] Focused/full Go, compatibility, update/contract, scope/diff, and canonical gates pass within exactly three files.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 20 from the validated external feedback.*
