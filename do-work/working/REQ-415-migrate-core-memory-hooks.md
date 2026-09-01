---
id: REQ-415
title: 'Migrate the core SessionStart and memory hooks into Go subcommands'
status: claimed
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-414]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
claimed_at: 2026-09-01T02:34:48Z
route: C
estimate:
  p50_active_minutes: 105
  confidence: low
  calculated_at: 2026-09-01T02:42:22Z
  basis:
    - Route C
    - 21-file write set
    - 9 new files
    - 7 subsystems involved
    - 5 acceptance criteria
    - dependency depth 1
    - persistence changes
    - asynchronous hook lifecycle behavior
    - cross-route regression gates
    - full-suite verification
---

# Migrate the Core SessionStart and Memory Hooks into Go Subcommands

## What
Replace hook domain logic with canonical `do-work-cli` hook subcommands.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Accepted a three-part plan for exact hook protocol projection, typed Go event owners, thin retained launchers, and differential consumer fixtures.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Migrate the core SessionStart hook plus memory SessionStart and Stop-capture hooks into Go subcommands.
- Preserve exact hook stdin/stdout protocols, event behavior, redaction, deduplication, timestamp/reservation repair, and failure handling.
- Keep each existing hook `.sh` path as a thin build-and-exec shim.
- Add fixture tests for malformed input, secrets, duplicate content, repair cases, and exact output.

## Constraints
- Hook launchers must remain safe in installed projects and must stop actionably when the canonical binary cannot build or run.

## Dependencies
Depends on REQ-414 (remaining core utility primitives used by SessionStart).

## Builder Guidance
Certainty level: Firm. Freeze byte-level hook protocols before moving any logic.

## Triage

**Route: C** — Complex

**Reasoning:** The three hook paths jointly own byte-consumed startup output, unattended mutation, transcript ingestion, redaction, UTF-8 budgets, deduplication, concurrent append behavior, best-effort instrumentation, and a nonblocking Stop contract.

**Planning:** Required and completed in `do-work/runs/work-2026-08-31-165510/REQ-415-plan.md`; repository exploration is in `REQ-415-exploration.md`.

## Plan

1. Add an optional exact protocol payload to the shared typed result so text mode emits hook bytes unchanged while JSON retains the same observation. Register `session-start`, `memory-session-start`, and `memory-stop-capture`, then freeze retained status/stdout/stderr/effect behavior in real-command and consumer fixtures before cutover.
2. Implement one standard-library `internal/hookcommands` family. Reuse doctor/repository timestamp authorities and corehelpers reservation cleanup; implement memory injection/capture with exact sentinel filtering, redaction-before-truncation, UTF-8 budgeting, hash deduplication, one append write, and best-effort ledger rows.
3. Replace the three retained hook bodies with thin sibling-CLI launchers, reconcile the memory and CLI ownership documents, and run focused/full Go, vet, exact Go 1.25, both hook behavior suites, contract/install/update gates, scope/diff hygiene, and the canonical maintainer gate.

## Decisions

### D-01: Preserve exact hook protocol through the typed result

**Decision:** DECIDE & STATE — use an optional protocol payload on the shared result; only text rendering short-circuits to its exact bytes, while JSON carries the same payload with findings/effects.

### D-02: Make the intended jq semantics universal

**Decision:** DECIDE & STATE — Go implements the full intended transcript selection behavior on every host. The degraded no-jq shell branch is not a second canonical protocol.

### D-03: Keep Stop nonblocking

**Decision:** DECIDE & STATE — domain behavior always succeeds; if the canonical launcher cannot build/run, the Stop shim reports an actionable stderr diagnostic but exits zero. SessionStart shims propagate launcher failures.

### D-04: Close REQ-463 at the consumed authority

**Decision:** EXPAND SCOPE — include the two existing reservation implementation/test files and satisfy REQ-463's unborn-Git and final-eligibility REDs. REQ-463 now depends on REQ-415 to prevent a parallel collision and will be dispositioned from this implementation's reviewed evidence.

**Reasoning:** Migrating SessionStart onto a known unsafe reservation primitive would violate both exact hook parity and the critical cleanup authority contract.

### D-05: Preserve observed compatibility quirks

**Decision:** DECIDE & STATE — retain raw loop-guard detection, legacy capture suppression to EOF, exact ledger source strings, current empty-side handling, redaction-before-truncation, hash-before-quoting, and observable newline/version parsing behavior.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/hookcommands/commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/hookcommands/commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/hookcommands/session_start.go` (new)
- `skills/do-work/tools/do-work-cli/internal/hookcommands/session_start_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/hookcommands/memory_start.go` (new)
- `skills/do-work/tools/do-work-cli/internal/hookcommands/memory_start_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/hookcommands/memory_capture.go` (new)
- `skills/do-work/tools/do-work-cli/internal/hookcommands/memory_capture_test.go` (new)
- `_dev/tests/memory-hook-behavior.sh` (new)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations_test.go` (modified)
- `skills/do-work/hooks/session-start.sh` (modified)
- `skills/do-work-knowledge/hooks/memory-session-start.sh` (modified)
- `skills/do-work-knowledge/hooks/memory-stop-capture.sh` (modified)
- `_dev/tests/session-start-hook-behavior.sh` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `skills/do-work-knowledge/actions/memory-reference.md` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**Files I will NOT touch:** hook JSON fragments, `install-memory-hooks.sh`, setup/memory action entry points, `commandruntime`, Git/publication primitives, Justfiles, non-hook shims, release metadata, or REQ-416/417 surfaces. Expansion requires a focused RED and owner approval.

**Acceptance criteria:**
- [ ] All three hook domains are registered Go commands, and retained paths are thin installed-layout launchers.
- [ ] Captured fixtures preserve exact stdin/stdout/stderr/status and filesystem/ledger/request effects.
- [ ] Memory capture preserves redaction, valid UTF-8 budgeting, sentinel isolation, hash deduplication, and one-write append behavior.
- [ ] Core SessionStart preserves timestamp/reservation housekeeping while closing REQ-463's committed-fresh-authority defect.
- [ ] Focused, full, compatibility, consumer, contract, install/update, and canonical gates pass with exact 21-file scope.

## Red-Green Proof
**RED prompt/case:** Replay captured valid, malformed, redacted, duplicate, and repair hook events against missing Go hook subcommands.
**Why RED now:** Hook behavior currently resides in three shipped shell implementations.
**GREEN when:** Go subcommands produce equivalent status/stdout/stderr/effects for every fixture and the original hook paths are launch-only shims.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
