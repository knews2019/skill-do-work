---
id: REQ-552
title: '[impact-negligible] Replace two coreutils exec sites with the pure Go the package already has'
status: claimed
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
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-09-06T00:40:32Z
  basis:
    - Route B
    - 5-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
    - full-suite verification
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go, skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go, _dev/tests/audit-lockins.sh, _dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh, _dev/tests/prescribed-shell-cases/architecture-report-preflight.sh]
route: B
claimed_at: 2026-09-06T00:38:56Z
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

## Triage

**Route: B** — Explore then build.

**Reasoning:** The two code edits are mechanical and the request names both files and both
replacements. What needed discovery is that the request's own Constraints cannot be satisfied: it says
"Tests unchanged; the existing package tests are the safety net" and "no test files beyond the
lock-in", and two heavy-tier fixture cases drive the very behaviour being changed by putting a fake
binary on `PATH`. The moment either call becomes in-process, the shim is inert and the case asserts
the opposite of what happens. That is discovery, not design — Builder Guidance already fixes the
approach — so exploration runs and planning does not.

**Planning:** Skipped.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Explore agent, read-only, re-verified against HEAD. Full report in the run directory as
`do-work/runs/work-2026-09-05-231943/REQ-552-exploration.md`.

**Both exec sites are still exactly where the request says**, at `commands.go:724` and
`architecture.go:133`, same line numbers as the earlier exploration. **The request's headline count is
stale in one half:** do-work-cli has 103 non-test `exec.Command` sites at HEAD, not 90; the "85 run
git" half is exact. The two coreutils sites, which are the only number the request acts on, are exact.

**Neither replacement the request names is usable as named.** `internal/corehelpers/inventory.go`
already walks the same tree with `filepath.WalkDir`, but it is not a reusable readability probe — it
parses every match and builds an ownership map, and it walks `do-work/working` too. `copyLast30DaysTree`
in `last30days.go:251-294` cannot be called for the `cp` site at all: `WalkDir` on a regular file
returns at `relative == "."` copying nothing, and its `O_EXCL` open fails against the temp file
`os.CreateTemp` already made. So both replacements are written fresh from the standard library, which
is what the request's Builder Guidance actually asks for.

**The behaviour to preserve was measured rather than assumed.** With GNU coreutils, `cp src dst` onto
a pre-created 0600 file leaves the mode at 0600 and copies the bytes; `os.CreateTemp` creates at 0600,
so the draft's own mode is never copied today. `os.WriteFile(stagedPath, data, 0o600)` reproduces that
exactly.

**One semantic difference is accepted and stated:** on a partly unreadable tree, `find` reports every
unreadable subdirectory and still exits non-zero, while the `WalkDir` replacement stops at the first.
Both produce a non-empty evidence string, so the gate behaves identically; only the evidence text
differs.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (modify) — the `find` probe becomes a `filepath.WalkDir` readability check that records the first traversal error as its evidence string; the `os.Stat` guard above it is untouched
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go` (modify) — the `cp` becomes `os.WriteFile(stagedPath, data, 0o600)`, keeping the `draft copy failed: ` evidence prefix
- `_dev/tests/audit-lockins.sh` (modify) — one new assertion block in the file's existing shape, appended after the Finding 2 block
- `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh` (modify) — the fake-`find` case is rewritten to drive the in-process failure; see D-01
- `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh` (modify) — the fake-`cp` case is rewritten the same way; see D-01

**Files I will NOT touch:** `internal/corehelpers/inventory.go` and `internal/toolboxcommands/last30days.go`
— both were considered as sources of a reusable primitive and neither fits; the reasons are in the
exploration. The pre-existing wart at `commands.go:721-722`, where a non-directory `do-work/archive`
yields the evidence string `<nil>`: real, out of scope, recorded as a discovered task. Every other
`exec.Command` site in either module.

**Acceptance criteria (restated from REQ):**
- [ ] Neither `find` nor `cp` is spawned by shipped code in either module
- [ ] The archive-walk failure is still detected and still produces non-empty evidence
- [ ] The draft copy still lands at the staged path with the same bytes and the same mode
- [ ] A lock-in assertion fails if a coreutils exec site returns to shipped code in either module
- [ ] The heavy tier is green, not just the fast gate — both rewritten fixture cases still detect the
      failure they name
