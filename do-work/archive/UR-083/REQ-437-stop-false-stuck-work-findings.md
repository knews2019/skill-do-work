---
id: REQ-437
title: 'Stop false stuck-work findings for active and terminal REQs'
status: completed
route: B
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-09-01T21:22:13Z
  basis:
    - Route B
    - 4-file write set
    - 2 subsystems involved
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go
  - skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go
  - skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan.go
  - skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan_test.go
related: [REQ-438, REQ-439, REQ-440, REQ-441, REQ-442, REQ-443, REQ-444]
batch: accepted-feedback-regressions
claimed_at: 2026-09-01T21:20:23Z
completed_at: 2026-09-01T21:33:40Z
commit: 406cfc5d
---

# Stop False Stuck-Work Findings for Active and Terminal REQs

## What

Make doctor reserve `STUCK-WORK` for nonterminal working REQs that lack recent activity. A terminal file stranded under `working/` must receive the terminal-location finding without a contradictory ownership warning, and a recently edited old claim must not be called abandoned.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add one stable modification timestamp to each discovered request, then use it only to suppress STUCK-WORK for terminal or recently touched working records. Verify the captured terminal, recent-edit, and inactive-claim replays with focused Go tests.
- [x] **[APPLY]:** Added the verified snapshot timestamp and a narrow STUCK-WORK eligibility predicate. The change remains within the declared repository-model and doctor files.
- [x] **[UNIFY]:** Ran `git diff --stat`, `git diff --check`, and `go vet ./...`; reviewed the four declared code and test files for scope, formatting, and debug artifacts. The request log and checkpoint are lifecycle evidence only.

## Finding Provenance

- **Verbatim claim / severity:** `[P2] Doctor reports actively edited and terminal working REQs as “stuck”.`
- **Evidence:** `doctor_scan.go:101-103` invokes stuck detection for every working record; `doctor_scan.go:202-227` considers only `claimed_at`. `forensics.md` says a recently modified file is likely active.
- **Origin / earned by:** Commit `210d1459` centralized a longstanding contradiction between the timestamp-only detector and the action's recent-activity rule.
- **Surface-cost:** Terminal filtering is N/A. Modification-time evidence is Earned by the old-claim/recent-edit replay; one snapshot field and focused regressions are cheaper than recurring false ownership warnings.

## Detailed Requirements

- Exclude terminal statuses from `STUCK-WORK` while retaining the appropriate stranded-terminal finding.
- Carry trustworthy file modification-time evidence through the repository snapshot and use the action's recent-activity boundary consistently.
- Preserve existing age severity behavior for genuinely inactive nonterminal claims.
- Keep diagnosis deterministic and do not introduce timestamp recovery or automatic ownership changes.

## Constraints

- This root cause is distinct from REQ-435 (Complete the Doctor-Forensics Delegation Contract); coordinate overlapping files without merging the intents.
- The activity check must be one evidence field and predicate, not a general liveness subsystem.

## Red-Green Proof

**RED prompt/case:** Scan one terminal REQ stranded under `working/` and one nonterminal REQ claimed more than a day ago but modified moments ago.
**Why RED now:** Doctor evaluates every working file and its stuck predicate sees only `claimed_at`.
**GREEN when:** The terminal fixture receives no `STUCK-WORK`, the recently edited fixture receives no `STUCK-WORK`, a genuinely old inactive claim still does, and the terminal-location finding remains intact.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm. Use the smallest snapshot evidence and predicate change that closes both named false-positive replays.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 6 from the validated external feedback.*

## Triage

**Route: B** - Medium

**Reasoning:** The bug has a precise behavioral replay, but its evidence must travel from repository discovery into the doctor predicate. The change is bounded to the shared snapshot and doctor scan paths.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The repository snapshot reads contained regular files under a rooted handle and currently discards their modification time. Doctor invokes STUCK-WORK for every working record and uses only claimed_at, while terminal-location detection runs separately. The smallest repair is a snapshot timestamp field plus one predicate that excludes terminal records and timestamps touched within the existing one-hour activity boundary.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go`
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go`
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan.go`
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan_test.go`

**Acceptance criteria:** terminal records left under working retain STRANDED-TERMINAL-REQUEST without STUCK-WORK; a recently modified old nonterminal claim has no STUCK-WORK; a genuinely inactive old nonterminal claim keeps its existing severity; diagnosis remains read-only and deterministic.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (modified) — retains the modification time verified during a contained request-file read.
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (modified) — proves the snapshot exposes a UTC, second-precision request timestamp.
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan.go` (modified) — limits STUCK-WORK to nonterminal records with no activity in the preceding hour.
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan_test.go` (modified) — replays terminal, recently modified, and genuinely inactive working claims.

**What was done:** Doctor now consumes trusted snapshot modification evidence rather than deriving activity from claimed_at alone. Terminal working records still receive the location diagnostic while active and terminal records avoid contradictory ownership warnings.

## Testing

**Red-green validation:**
- RED: `go test ./internal/doctor -run '^TestStuckWorkSkipsTerminalAndRecentlyModifiedClaims$' -count=1` failed before the change because both REQ-130 (terminal) and REQ-131 (recently edited) emitted STUCK-WORK.
- GREEN: the same focused test passes after the predicate change, retaining an error-severity STUCK-WORK finding for inactive REQ-132.
- Snapshot evidence: `go test ./internal/repositorymodel ./internal/doctor -count=1` passes.
- Module regression: `go test ./... -count=1` passes.

## Qualification

Passed — 4 declared files verified, 4 acceptance criteria traced, and P-A-U confirmed. `qualify.sh` passed the file, diff, wiring, and debug-artifact checks; `scope-drift.sh` confirmed the Implementation Summary matches the declared Scope.

## Review

**Overall: 100%** | 2026-09-01T21:28:57Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.

**Minor findings:** 0 (report only)
**Acceptance:** Pass — the focused repository-snapshot scan proves terminal and recently modified claims suppress STUCK-WORK while an inactive claim retains its error-severity finding and terminal-location diagnosis remains present.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*
