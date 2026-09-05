---
id: REQ-418
title: 'Migrate toolbox commands and absorb audit-metrics into do-work-cli'
status: completed-with-issues
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, skills/do-work-toolbox/tools/audit-metrics/prime-audit-metrics.md]
tdd: true
suggested_spec:
depends_on: [REQ-417]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-419, REQ-420]
batch: go-no-llm-command-platform
claimed_at: 2026-09-01T08:41:56Z
completed_at: 2026-09-01T12:01:48Z
commit: 82534d36
---

# Migrate Toolbox Commands and Absorb Audit-Metrics into Do-Work-CLI

## What
Move deterministic toolbox utility behavior into the shared CLI and consolidate the separate audit-metrics module into it.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Worktree builder plan in run scratch (`REQ-418-plan.md`); toolbox family as shared-CLI handlers plus audit-metrics port under the frozen 30-file scope.
- [x] **[APPLY]:** Builder `a7c975c5` (merged `94560fde`), sole remediation `a43b2587` (merged `82534d36`); 21 remediation paths within the documented 32-path ceiling.
- [x] **[UNIFY]:** `git diff --check` clean over `b713fa8b..a43b2587`; go vet, full uncached tests, race lanes, Windows cross-compile, and canonical `maintainer-verify.sh` all exit 0 (verified independently by the fresh re-review). Checkboxes marked by the orchestrator from the handback and gate record per the remediation handback's no-lifecycle instruction.

## Detailed Requirements
- Implement `do-work-note`, architecture-report preflight/publication, report-image lifecycle, portfolio publication, and last30days install/check.
- Preserve report/media process handling, publication atomicity, failure cleanup, and existing user-visible behavior.
- Absorb audit-metrics inventory, folders, churn, and hotspots behavior into `do-work-cli` with compatible flags/output plus JSON.
- Keep last30days’ Python prerequisite as a target-tool requirement, not a do-work implementation dependency.

## Constraints
- `queue-kanban` remains separate; only audit-metrics is consolidated.
- Characterization parity is required before retiring the separate audit-metrics source tree.

## Dependencies
Depends on REQ-417 (knowledge/store command migration precedes the toolbox family).

## Builder Guidance
Certainty level: Firm. Port existing audit-metrics tests rather than re-deriving its mature Git behavior.

## Red-Green Proof
**RED prompt/case:** Run current toolbox and audit-metrics fixture suites against absent shared CLI equivalents, including media cancellation and target-Python absence.
**Why RED now:** Toolbox behavior lives in shell and audit metrics lives in a separate Go module without the common result contract.
**GREEN when:** Shared commands match current status/output/effects, add actionable JSON, and preserve target-specific dependency reporting without Python/jq implementation branches.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Implementation Summary

Seven shared-CLI handlers (`do-work-note`, architecture scan/publication, single and
batch report images, portfolio publication, last30days install/check) plus four
audit-metrics modes (inventory, folders, churn, hotspots) in
`skills/do-work/tools/do-work-cli/internal/toolboxcommands/`, with command/result
registration and typed JSON. The remediation added a confined-publication family
(`os.Root` ancestor validation, pinned parent handles, private staging, exclusive
claims), identity-bound tracked rollback in `internal/gittransaction/`, one invocation
context with owned-process-group TERM/grace/KILL escalation, snapshot-first portfolio
publication, separated display/filesystem scan paths, a leading-only option region
with `--` escape, and option-presence tracking for audit-flag parity — plus thirteen
committed `TestRemediation*` tests. Standalone audit sources and retained toolbox
scripts are unchanged; retirement stays with REQ-420. Builder `a7c975c5` merged
`94560fde`; remediation `a43b2587` merged `82534d36`.

## Remediation

- **Initial review** (`REQ-418-review.md`, against `94560fde`): Request changes,
  Overall 35%, Acceptance Fail — nine Important findings (four `impact-critical`:
  linked-ancestor escape, pre-eligibility last30days mutation, rollback overwriting a
  concurrent writer, incomplete media cancellation) plus two Minor.
- **Sole remediation** (`a43b2587`, merged `82534d36`): closed all nine in named
  tests; passed the full focused/race/full/vet/Go 1.25/Windows/differential/contract/
  install/update/canonical stack. The allowance is exhausted.
- **Fresh independent re-review** (`REQ-418-rereview.md`, at `d924b2fc`, product
  paths byte-identical to `82534d36`): Approve with follow-ups, Overall 69%,
  Acceptance Partial. Eight of nine original counterexamples re-reproduced and
  confirmed closed; finding 6 partially closed; three new Important defects in
  `architecture.go`, two of them reproduced regressions against `a7c975c5`. With the
  sole remediation spent, this REQ is terminal as `completed-with-issues` and every
  residual Important finding is durably routed below.

## Review

Fresh re-review verdict: **Approve with follow-ups** — 69%, Acceptance Partial.
Gates at review time, all exit 0, judged directly: `go test -count=1 ./...`,
`go vet ./...`, `go test -race ./internal/toolboxcommands ./internal/gittransaction`,
`bash _dev/tests/maintainer-verify.sh` (declared no-browser skip only), thirteen
`TestRemediation*` tests run verbosely, Windows cross-compile.

Residual finding routing (fold-first):

- Finding 6 residual: portfolio snapshot-first failure branch sets `ExactTextOutput`
  on a non-success outcome, hiding the canonical refusal and snapshot-retained
  warning in text mode (`portfolio.go:99-105`, `resultmodel.go:384-385`) —
  `impact-user-visible` → appended to REQ-420 (maintainer-approved disposition,
  2026-09-01: goal-shaped output-parity criterion for the not-yet-authoritative CLI).
- N1: unbounded CPU-spinning bundle-claim loop in `architecture-report-preflight
  --publish` on any persistent claim error (`architecture.go:156-164`); reproduced
  regression against `a7c975c5` — `impact-critical` → REQ-483 created.
  ⚠ CRITICAL review finding auto-queued as REQ-483.
- N2: `architecture-report-preflight --commit` dead in every argument position
  (`architecture.go:144` passes `dryRun=true` alongside the caller's `commit`);
  reproduced regression — `impact-user-visible` → REQ-483 created. Challenge noted:
  the approved disposition routed N2 to REQ-420, but N1's fix edits the same ten
  lines of the same file; splitting them would leave a one-line regression waiting
  behind the whole spine while a builder works beside it. Resolution: combined into
  REQ-483, recorded here and in that REQ.
- N3: five usage strings still advertise the trailing `[--dry-run|--commit]` form the
  finding-8 fix stopped accepting (`architecture.go:38`, `report_image.go:38`,
  `report_image.go:103`, `portfolio.go:20`, `last30days.go:33`) —
  `impact-user-visible` → appended to REQ-420.
- Four new Minor findings are report-only in `REQ-418-rereview.md` (dead
  `ensureLast30DaysExclude`; swallowed `restore()` error under a deferred backup
  delete with a misleading success message; decorative audit comparator test;
  P-A-U boxes — the last resolved by this file's orchestrator-marked checkboxes).
