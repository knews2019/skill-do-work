---
id: REQ-552
title: '[impact-negligible] Replace two coreutils exec sites with the pure Go the package already has'
status: pending
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: backend
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-550]
related: [REQ-549, REQ-550, REQ-551, REQ-553, REQ-554, REQ-555, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go, skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go, _dev/tests/audit-lockins.sh]
---

# Replace two coreutils exec sites with the pure Go the package already has

## What
Of 90 `exec.Command` sites in do-work-cli, 85 run `git`; two spawn `find` or `cp` for work the same package already does in stdlib Go. Replace the `find` probe in `internal/corehelpers/commands.go` with an `os.Stat`/`WalkDir` readability check reusing the inventory walk, and the `cp` in `internal/toolboxcommands/architecture.go` with the `io.Copy` primitive the package already has.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
The `find` site runs on every `audit-archive-timestamps` invocation only to harvest an error string while the Go walk traverses the same tree again, so the archive is walked twice per run; the `cp` site is the one non-Go copy in a package that already ships a pure-Go recursive copy.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 9, sweep_key `exec-where-pure-go-exists`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -10. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `internal/corehelpers/commands.go` — `exec.Command("find", archiveRoot, "-name", "REQ-*.md", "-print0")`, output discarded, called from the `audit-archive-timestamps` path; `internal/corehelpers/inventory.go` already walks the same tree with `filepath.WalkDir` and the same predicate.
- `internal/toolboxcommands/architecture.go` — `exec.Command("cp", draftPath, stagedPath)` behind `DO_WORK_COMPATIBILITY_SHIM == "1"`; `internal/toolboxcommands/last30days.go` has `copyLast30DaysTree` with the copy primitive.
- Behaviour preserved: the same error string class is returned when the archive is unreadable; the same file lands at `stagedPath` with the same mode.
- Reproduce at dc8a64e3 (prints the two sites): `rg -n 'exec\.Command(Context)?\((ctx, )?"(find|cp|mkdir|grep|sed|ls|rm|mv|cat|touch|head|tail|wc)"' skills/do-work/tools/do-work-cli skills/do-work-board/tools/queue-kanban --glob '!*_test.go'`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- Tests unchanged; the existing package tests are the safety net.
- Lock-in limit: coreutils exec sites in the two Go modules: 0 after this REQ (today 2).

## Dependencies
Depends on REQ-550 (collapse exported delegates) so the `corehelpers` write set has no pending overlap. REQ-557 depends on this REQ.

## Builder Guidance
Firm: pure Go on both sides. Latitude on where the readability check lives inside `corehelpers`.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints two lines (`find` in commands.go, `cp` in architecture.go).
**GREEN when:** It prints nothing; package tests green; the lock-in pins coreutils exec sites in both Go modules at 0.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5660 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "do-work-cli internals".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for exec-where-pure-go-exists.*
