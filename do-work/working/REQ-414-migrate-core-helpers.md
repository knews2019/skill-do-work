---
id: REQ-414
title: 'Migrate remaining core checks, publication helpers, Git helpers, and surveys'
status: claimed
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-413]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
claimed_at: 2026-09-01T00:24:42Z
route: C
estimate:
  p50_active_minutes: 135
  confidence: low
  calculated_at: 2026-09-01T00:25:26Z
  basis:
    - Route C
    - 30-file write set
    - 20 new files
    - 10 subsystems involved
    - 4 acceptance criteria
    - dependency depth 1
    - persistence changes
    - cross-route regression gates
    - full-suite verification
---

# Migrate Remaining Core Checks, Publication Helpers, Git Helpers, and Surveys

## What
Move all remaining core utility domain logic into `do-work-cli` subcommands.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Migrate preflight, qualification, scope drift, protected/uncommitted inventories, file association, and commit-hash recording.
- Migrate screenshot and download publication, local Git exclude, blocked checks, commit-diff display, exact deletion staging, timestamp/reservation helpers, and handoff surveys.
- Preserve existing exit status, output, filesystem effects, byte verification, permissions, redaction, and error behavior through characterization fixtures.
- Make findings actionable in both text and JSON without forcing an LLM to rescan.

## Constraints
- Existing `.sh` paths remain for compatibility but ultimately contain no domain logic.
- Target-specific Python project checks remain valid target probes rather than CLI implementation dependencies.

## Dependencies
Depends on REQ-413 (publication and common transaction primitives).

## Builder Guidance
Certainty level: Firm. Inventory every shipped core utility before migration and map each to a named CLI subcommand.

## Red-Green Proof
**RED prompt/case:** Run the current core utility fixture matrix against corresponding absent CLI subcommands and compare status, output, and filesystem effects.
**Why RED now:** Core behavior still lives across numerous shell implementations.
**GREEN when:** Every mapped Go subcommand passes parity fixtures and its compatibility path only builds/executes the canonical binary.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Triage

**Route: C** — Complex

**Reasoning:** This migrates many independently observable shell utilities into one typed command platform, spans multiple mutation and inspection domains, and requires broad parity plus compatibility-shim evidence.

**Planning:** Required

## Plan

1. Add one standard-library `internal/corehelpers` command family for the read-only/read-mostly checks, inventories, association, commit provenance, and handoff survey. Reuse `requestmodel`, `repositorymodel`, `requeststate`, `gittransaction`, and `resultmodel` rather than creating second scanners, editors, or result conventions. Characterization fixtures must pin every legacy status, output fact, Git/filesystem effect, and REQ-390's first-backticked-path scope rule.
2. Add the publication, Git/runtime, timestamp, reservation, and blocked-probe mechanics to the same family. Publication and cleanup own the exact created object rather than a later pathname; blocked probes own the complete child process tree; timestamp and commit provenance factor existing authorities for reuse. Replace the current Go-to-shell dependencies in `archivefetch` and `nextselection` without changing retained shell compatibility paths in this REQ.
3. Register every mapped command at `cmd/do-work-cli/main.go`, prove actionable text/JSON parity through real-command tests, update the CLI prime, and leave hook migration, flat recipes/action aliases, and final thin compatibility shims to REQ-415, REQ-419, and REQ-420 respectively. Run focused/full Go, vet, exact Go 1.25, platform process-control compile/tests, legacy characterization suites, scope/diff hygiene, and the canonical maintainer gate.

**Architecture decisions:** One helper family sits over existing domain authorities; typed ordered evidence drives both renderers; every shipped utility is explicitly mapped either to a new command or an already-delivered authority; runtime and publication ownership are condition-complete; and target-supplied probes may use shell/Python while the helper implementation itself depends only on Go and Git.

**Testing approach:** Start from unknown-command and legacy-characterization REDs, then cover hostile paths, binary/untracked inventory, merge diffs, concurrent publication, parent swaps, timeout descendants, token redaction, timestamp/refusal shapes, reservation races, text/JSON parity, and exact side effects.

**Plan validation:** All four Detailed Requirements map to the three tasks; no task is orphaned. The plan stays at the Route C ceiling of three tasks. Exploration must freeze the exact surface and challenge the initial broad `corehelpers` grouping before code.

*Generated by Plan agent*

## Folded From REQ-390 (2026-08-30)

- **`tools/checks/scope-drift.sh` reads every backticked token in a "Files I will
  touch" bullet as a declared path.** `emit_backticked_paths` splits the whole line
  on backticks and prints each odd-indexed part, so a bullet whose rationale contains
  an ordinary code span — `` - `path/to/file` (modify) — adds one `flex-wrap`
  declaration `` — declares a phantom second path and reports it as "declared in
  ## Scope but never touched". Observed on REQ-390, where it produced a false DRIFT
  line that had to be worked around by rewording the REQ rather than fixing the
  check. The path is the first backticked token on the bullet; the rest is prose.
  Worth closing when this REQ ports the checks to Go, where the extraction can be
  written once and tested.
