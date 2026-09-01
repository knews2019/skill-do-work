---
id: REQ-420
title: 'Replace shell implementations with shims and prove whole-suite parity'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-419]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419]
batch: go-no-llm-command-platform
sweep: true
sweep_key: canonical-command-shim-parity-and-authority
---

# Replace Shell Implementations with Shims and Prove Whole-Suite Parity

## What
Complete the migration by making every retained shell path a thin launcher and enforcing full suite parity mechanically.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Replace domain logic in all 41 shipped shell utilities and hooks with thin build-and-exec compatibility shims while preserving every path.
- Add a mechanical contract that retained `.sh` files are launchers only and contain no embedded Python or jq implementation branches.
- Require old shim and new subcommand parity for exit status, output, and filesystem effects on characterized fixtures.
- Remove retired audit-metrics sources after its tests and behavior move into `do-work-cli`.
- Extend `_dev/tests/maintainer-verify.sh` with `go vet` and uncached `do-work-cli` tests, retain queue-kanban verification, and replace the separate audit-metrics lane.
- Prove final installation/update without Python/jq, every Just command without an LLM, unchanged skill aliases, and actionable findings that avoid repository rescans.

## Constraints
- Keep target-specific Python checks only for Python targets such as Python project preflight and last30days.
- Run the complete canonical maintainer gate and whole-suite parity verification before acceptance.

## Dependencies
Depends on REQ-419 (complete command and recipe surface).

## Builder Guidance
Certainty level: Firm. Retire implementations only after their parity fixtures pass against the Go engine.

## Red-Green Proof
**RED prompt/case:** Run the shell-thinness contract and full parity suite while any shipped shell still contains domain logic or Python/jq branches, and run install with those tools absent.
**Why RED now:** All 41 scripts/hooks still contain or route to shell domain implementations, and audit-metrics has a separate verification lane.
**GREEN when:** The mechanical contract proves launcher-only shell, parity fixtures pass, audit-metrics is consolidated, Go/board maintainer lanes pass uncached, and final no-Python/no-jq acceptance succeeds.
**Validation:** User confirmed via the supplied implementation plan.

## Instances

- [ ] All retained shell utilities/hooks and the standalone audit-metrics implementation: replace domain logic with launcher-only shims and prove whole-suite behavioral parity. (original REQ-420 / UR-081)
- [ ] `_dev/primes/prime-kanban-board.md:16`, `skills/do-work/actions/work-reference.md:105,113`, `skills/do-work-board/tools/queue-kanban/timestamp.go:19-22`, and `frontmatter_cli.go:24-32`: retire the board-only compiler and mandatory shell-fallback contract now contradicted by build-on-demand `do-work-cli`. (found by REQ-419 / UR-081)

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Folded From REQ-406 (2026-08-30)

- **This REQ's first clause is already done.** REQ-406 added the `do-work-cli go vet`
  and `do-work-cli uncached tests` lanes to `_dev/tests/maintainer-verify.sh` and
  registered `_dev/tests/do-work-cli-launcher-behavior.sh` in
  `_dev/tests/contract-regressions.sh`, moving all four of maintainer-verify's
  self-test enumerations in lock-step. The orchestrator accepted the overlap
  deliberately: an unwired module is an unverified module, and leaving it unwired
  would have meant fourteen REQs of unrun code. **What remains here is the
  destructive half** — replacing the separate audit-metrics lane — which REQ-406
  deliberately did not touch.
- **`maintainer-verify.sh`'s self-test EXIT trap reads a function-local.**
  `trap 'rm -rf -- "$self_test_root"' EXIT` is set inside `run_self_test`, where
  `self_test_root` is declared `local`. On any self-test failure the function returns
  before `trap - EXIT`, so the trap fires at process exit with the variable out of
  scope and prints an `unbound variable` line after the real `FAIL:`, leaving the
  fixture directory behind. Pre-existing and cosmetic — it does not mask the failure
  or change the exit status — but this REQ is already editing that file.

## Folded From REQ-464 and REQ-465 (2026-09-01)

Both were REQ-414 fresh re-review findings (3 and 4; see
`do-work/runs/work-2026-08-31-165510/REQ-414-rereview.md`) spawned as standalone
REQs only because this REQ was dependency-gated at capture time and therefore not
fold-eligible. The maintainer folded them here on 2026-09-01: their whole substance
is parity/actionability acceptance criteria for the shim conversion this REQ owns,
and building either as a separate framework first would be duplicated work.

**From REQ-465 (core-helper differential parity)** — the parity requirement in this
REQ's third bullet is satisfied for the 17 core helper commands only when:

- Every retained helper runs beside its Go command on the same characterized
  fixtures in text and JSON modes; exact semantic status, ordered facts, affected
  paths, recovery/verification argv, and filesystem/index/private-state effects are
  compared mechanically — not the current registration/smoke matrix in
  `internal/corehelpers/commands_test.go`, which accepts any exit status up to one
  and checks renderer agreement only.
- Fixtures cover happy, non-clean, hostile-path, combined-state, dry-run, refusal,
  and concurrent-change cases earned by each helper's retained contract.
- Mutation tests prove the parity adapter cannot silently accept a divergence
  across status, fact, action, path, and effect dimensions (RED case: a controlled
  `AD` inventory misclassification or wrong recovery argv must fail the matrix).

**From REQ-464 (specifically actionable structured findings)** — the
actionable-findings acceptance bullet is satisfied only when:

- Each non-clean core-helper finding carries next/verification argv whose success
  specifically resolves or proves that exact condition — no family-wide `git
  status`/`git diff`/rerun fallbacks (`internal/corehelpers/commands.go`).
- Each dirty path/state is one typed handoff record (spaces, newlines, renames,
  missing/prunable worktrees, simultaneous rows) rather than one newline-bearing
  opaque evidence string (`internal/corehelpers/handoff.go`), with deterministic
  text/JSON parity from the same observation set.
- Mutation tests that replace a specific action or collapse structured dirty rows
  fail the characterization gate.

## Folded From REQ-418 fresh re-review (2026-09-01)

Two residual findings from REQ-418's fresh re-review (`completed-with-issues`,
remediation allowance exhausted), folded here by maintainer disposition 2026-09-01:
both are goal-shaped output/documentation criteria for the toolbox family whose
parity this REQ proves, not shipped-behavior defects — the retained scripts remain
the shipped implementations until this REQ converts them. The critical loop and dead
`--commit` from the same re-review went to standalone REQ-483, not here.

**Finding 6 residual (portfolio refusal visibility)** — the parity requirement is
satisfied for `publish-portfolio-summary` only when a snapshot-first canonical
refusal is visible: `portfolio.go:99-105` sets `ExactTextOutput` on the failure
branch (every sibling guards it on `OutcomeSuccess`), and `resultmodel.go:384-385`
returns it in place of the rendered block, so text mode shows only the snapshot path
and exit 2 — no refusal reason, no snapshot-retained warning. The retained script
reports both; the parity fixtures must compare that failure-branch output.

**N3 (usage-string restatements)** — the shim conversion's acceptance includes
self-consistent command documentation: five usage strings still advertise the
trailing `[--dry-run|--commit]` form the leading-only option region deliberately
stopped accepting (`architecture.go:38`, `report_image.go:38`, `report_image.go:103`,
`portfolio.go:20`, `last30days.go:33`; verified: trailing `--dry-run` is a
wrong-argument-count error for portfolio and becomes the `[source-repository]`
positional for last30days). Sweep wording variants, not only these five sites.

Provenance: `do-work/runs/work-2026-08-31-165510/REQ-418-rereview.md` while the run
scratch survives; durably, the `## Review` section of
`do-work/archive/REQ-418-migrate-toolbox-absorb-audit-metrics.md`.

## Folded From REQ-467 (2026-09-01)

REQ-415's fresh re-review residual finding 2, minted standalone only because this REQ
was dependency-gated at capture time ("REQ-420 adjacent but not dependency-ready" —
its own body). The maintainer folded it here on 2026-09-01: reconciling SessionStart
authority guidance is intrinsic to this REQ's shim conversion of the very scripts the
guidance misattributes. The shim-conversion bullet is satisfied for
`cleanup-req-reservations.sh` and `repair-req-timestamps.sh` only when:

- Every live shipped restatement of SessionStart reservation-cleanup/timestamp-repair
  ownership — `README.md`, `skills/do-work/actions/capture.md`,
  `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`, and
  `skills/do-work-board/tools/queue-kanban/verify.go` findings/remedies/tests —
  identifies the registered Go hook/core owner, matching wording variants, not only
  filenames; retained scripts are described only by their thin-launcher role, and no
  guidance recommends the pre-REQ-463 fail-open cleanup path.
- Board findings and recovery commands emit runnable canonical CLI argv instead of
  retired script paths.
- A contract test fails when live shipped guidance reintroduces the old ownership claim.

Provenance: `do-work/archive/REQ-415-migrate-core-memory-hooks.md` (the re-review
scratch it cited was pruned with the run directory).

## Folded From REQ-473 (2026-09-01)

REQ-416's fresh re-review residual F2, a documented-behavior divergence whose stated
requirement is a differential parity proof — this REQ's parity suite is where that
proof lives. The parity requirement is satisfied for `bkb-init` only when:

- The characterized fixtures include the documented outside-target contract
  (`bkb init ~/research` and a parent-relative equivalent) with the pre-migration
  action as oracle: absolute and parent-relative targets outside the invocation Git
  root route to the standalone initialization flow without weakening path safety, and
  the invocation repository is left untouched.
- Exact user-supplied targets are preserved in text/JSON next argv, Just recipes, and
  verification argv for success, dry-run, and every refusal.
- Fixtures cover invocation inside Git, target inside the same root, target outside
  it, symlink-spelled paths, spaces, and ambiguous/unavailable Git evidence, with all
  rooted publication/rollback protections retained.

Provenance: `do-work/archive/REQ-416-implement-bkb-dream-commands.md`.

## Folded From REQ-474 (2026-09-01)

REQ-416's fresh re-review residual finding 3, a characterized parity break. The
parity requirement is satisfied for `bkb-status` only when:

- Article and topic-cluster counts parse from `wiki/_master_index.md` exactly as the
  pre-migration status action defined them; disk inventories appear only where
  separately named and are never relabeled as master-index counts.
- The fixture matrix includes a master index declaring 17 articles / 3 topic clusters
  while disk inventory differs (today's implementation reports 0/0), plus missing,
  malformed, duplicate, and inconsistent count declarations returning actionable
  typed findings.
- Text/JSON parity and byte-for-byte read-only behavior are proven.

Provenance: `do-work/archive/REQ-416-implement-bkb-dream-commands.md`.

## Folded From REQ-476 (2026-09-01)

REQ-417's fresh re-review residual finding 2 (sweep
`memory-bkb-audit-tdd-coverage-missing`, impact-negligible, effort-mechanical): the
plan-required committed BKB audit matrix, folded because this REQ's whole-suite
characterization matrix must contain those rows regardless. The characterization
requirement is satisfied for the BKB audit surface only when:

- Real-command fixtures cover explicit, default, absolute, and parent BKB discovery;
  repository shape; committed history/authors; inbound references; absent and
  malformed ledger evidence; pre-ledger fairness; and exact classification boundaries.
- `auditBKBEngine` and `countBKBInboundReferences` execute under the committed test
  suite (both report 0.0% coverage today) with focused coverage evidence retained.
- Text/JSON parity and byte-for-byte read-only behavior are proven.

Provenance: `do-work/archive/REQ-417-implement-interview-memory-commands.md`.
