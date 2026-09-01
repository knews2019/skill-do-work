---
id: REQ-420
title: 'Replace shell implementations with shims and prove whole-suite parity'
status: claimed
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-419, REQ-478]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419]
batch: go-no-llm-command-platform
sweep: true
sweep_key: canonical-command-shim-parity-and-authority
write_set:
  - skills/do-work-knowledge/hooks/memory-session-start.sh
  - skills/do-work-knowledge/hooks/memory-stop-capture.sh
  - skills/do-work-knowledge/scripts/install-memory-hooks.sh
  - skills/do-work-knowledge/scripts/lexical-memory-recall.sh
  - skills/do-work-toolbox/scripts/architecture-report-preflight.sh
  - skills/do-work-toolbox/scripts/generate-report-image-batch.sh
  - skills/do-work-toolbox/scripts/generate-report-image.sh
  - skills/do-work-toolbox/scripts/install-last30days.sh
  - skills/do-work-toolbox/scripts/publish-portfolio-summary.sh
  - skills/do-work/hooks/session-start.sh
  - skills/do-work/scripts/add-local-git-exclude.sh
  - skills/do-work/scripts/atomic-download.sh
  - skills/do-work/scripts/audit-archive-timestamps.sh
  - skills/do-work/scripts/capture-screenshot.sh
  - skills/do-work/scripts/cleanup-req-reservations.sh
  - skills/do-work/scripts/handoff-state-survey.sh
  - skills/do-work/scripts/protected-inventory.sh
  - skills/do-work/scripts/repair-req-timestamps.sh
  - skills/do-work/scripts/run-blocked-check.sh
  - skills/do-work/scripts/show-commit-diff.sh
  - skills/do-work/scripts/stage-exact-deletion.sh
  - skills/do-work/tools/checks/archive-collision.sh
  - skills/do-work/tools/checks/associate-files.sh
  - skills/do-work/tools/checks/blanked-req-scan.sh
  - skills/do-work/tools/checks/preflight.sh
  - skills/do-work/tools/checks/qualify.sh
  - skills/do-work/tools/checks/record-commit-hash.sh
  - skills/do-work/tools/checks/scope-drift.sh
  - skills/do-work/tools/checks/uncommitted-inventory.sh
  - skills/do-work/tools/do-work-cli.sh
  - skills/do-work/tools/do-work-update.sh
  - skills/do-work/tools/estimate-p50.sh
  - skills/do-work/tools/fetch-upstream-archive.sh
  - skills/do-work/tools/install-do-work-suite.sh
  - skills/do-work/tools/replace-text-section.sh
  - skills/do-work/tools/select-simple-reqs.sh
  - skills/do-work/tools/validate-suite-manifest.sh
  - tools/fetch-upstream-archive.sh
  - tools/install-do-work-suite.sh
  - tools/replace-text-section.sh
  - tools/validate-suite-manifest.sh
  - _dev/tests/fixtures/shipped-shell-command-map.tsv
  - _dev/tests/shipped-shell-parity.sh
  - _dev/tests/shipped-shell-thinness.sh
  - _dev/tests/contract-regressions.sh
  - _dev/tests/maintainer-verify.sh
  - _dev/tests/prescribed-shell-canonicalization.sh
  - _dev/tests/prescribed-shell-cases/install-memory-hooks.sh
  - _dev/tests/staged-skills-contract.sh
  - skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/checks_test.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/git_helpers.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/git_helpers_test.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/handoff.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/handoff_test.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/publication.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/publication_test.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/reservations.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/reservations_test.go
  - skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands.go
  - skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan.go
  - skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan_test.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init_test.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan_test.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands_test.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_churn_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_inventory_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_metrics_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/last30days.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/last30days_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/portfolio.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/portfolio_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_test.go
  - README.md
  - _dev/primes/prime-kanban-board.md
  - skills/do-work/actions/capture.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work-board/tools/queue-kanban/frontmatter_cli.go
  - skills/do-work-board/tools/queue-kanban/frontmatter_cli_test.go
  - skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
  - skills/do-work-board/tools/queue-kanban/timestamp.go
  - skills/do-work-board/tools/queue-kanban/timestamp_test.go
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
  - skills/do-work-toolbox/tools/audit-metrics/.gitignore
  - skills/do-work-toolbox/tools/audit-metrics/churn.go
  - skills/do-work-toolbox/tools/audit-metrics/churn_test.go
  - skills/do-work-toolbox/tools/audit-metrics/distribution.go
  - skills/do-work-toolbox/tools/audit-metrics/distribution_test.go
  - skills/do-work-toolbox/tools/audit-metrics/git_support.go
  - skills/do-work-toolbox/tools/audit-metrics/go.mod
  - skills/do-work-toolbox/tools/audit-metrics/inventory.go
  - skills/do-work-toolbox/tools/audit-metrics/inventory_test.go
  - skills/do-work-toolbox/tools/audit-metrics/main.go
  - skills/do-work-toolbox/tools/audit-metrics/prime-audit-metrics.md
  - _dev/tests/prescribed-shell-cases/generate-report-image.sh
  - _dev/tests/prescribed-shell-cases/generate-report-image-batch.sh
  - _dev/tests/prescribed-shell-cases/repair-req-timestamps.sh
  - _dev/tests/prescribed-shell-cases/add-local-git-exclude.sh
  - _dev/tests/prescribed-shell-cases/architecture-report-preflight.sh
  - _dev/tests/prescribed-shell-cases/atomic-download.sh
  - _dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh
  - _dev/tests/prescribed-shell-cases/capture-screenshot.sh
  - _dev/tests/prescribed-shell-cases/cleanup-req-reservations.sh
  - _dev/tests/prescribed-shell-cases/install-last30days.sh
  - _dev/tests/prescribed-shell-cases/lexical-memory-recall.sh
  - _dev/tests/prescribed-shell-cases/protected-inventory.sh
  - _dev/tests/prescribed-shell-cases/publish-portfolio-summary.sh
  - _dev/tests/prescribed-shell-cases/qualify.sh
  - _dev/tests/prescribed-shell-cases/run-blocked-check.sh
  - _dev/tests/prescribed-shell-cases/show-commit-diff.sh
  - _dev/tests/prescribed-shell-cases/stage-exact-deletion.sh
claimed_at: 2026-09-01T15:16:17Z
route: C
estimate:
  p50_active_minutes: 210
  confidence: low
  calculated_at: 2026-09-01T15:26:01Z
  basis:
    - Route C
    - 110-file write set
    - 3 new files
    - 12 subsystems involved
    - 20 acceptance criteria
    - dependency depth 1
    - persistence changes
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
---

# Replace Shell Implementations with Shims and Prove Whole-Suite Parity

## What
Complete the migration by making every retained shell path a thin launcher and enforcing full suite parity mechanically.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Preserve the existing shell behavior suite as the legacy-first oracle, add a complete path-to-command inventory and launcher-only ratchet, close the characterized Go gaps, then cut over the retained paths before deleting audit-metrics and reconciling authority text. Verify each phase with focused shell and Go tests; leave the complete maintainer gate to orchestration.
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

---

## Triage

**Route: C** - Complex

**Reasoning:** This is a repository-wide compatibility migration across 41 public shell paths, several Go command families, Git/filesystem effects, board behavior, authority documentation, and the canonical verification gate. Behavioral characterization must precede the destructive cutover.

**Planning:** Required

## Plan

1. Capture a non-vacuous legacy oracle before deleting shell implementations: inventory all 41 shipped paths and add differential fixtures that compare status, ordered output/facts, argv, Git/index/private state, and filesystem effects, including the full 17-command core-helper matrix and mutation dimensions.
2. Close Go compatibility gaps before cutover: add missing compatibility commands and correct actionable findings, portfolio refusal output, usage text, SessionStart ownership/remedies, BKB outside-root initialization, master-index status counts, and committed BKB audit coverage.
3. Convert all 41 retained paths in dependency-safe order to argv-preserving build-and-exec launchers; keep only the documented pre-install bootstrap seam where the CLI cannot exist yet.
4. Delete the standalone audit-metrics implementation after its behavior and tests are owned by `do-work-cli`, and reconcile every live authority/compiler/fallback restatement.
5. Run thinness, differential parity, mutation, no-Python/no-jq, installer/updater, Just, Go, board, and full canonical maintainer verification against the final state.

**Plan validation:** Every detailed and folded requirement maps to at least one task, and no task is orphaned. Warning: the plan has 5 tasks; quality normally degrades past 3, but splitting would make the required pre-cutover oracle and final whole-suite cutover independently vacuous, so the builder must preserve this order within one REQ.

*Generated by Plan agent*

## Exploration

- The shipped inventory is exactly 41 tracked shell paths: 37 below `skills/` and 4 root installer-tool mirrors.
- Existing Go owners already cover most shell behavior under `skills/do-work/tools/do-work-cli/internal/`; the missing compatibility seams are memory-hook installation, P50 estimation, archive-collision and blanked-REQ checks, plus a few core-helper inputs that remain shell-derived.
- The 17 core helpers already share result rendering, but current tests are smoke/renderer checks rather than exact differential parity. New fixtures must be recorded before the shim cutover so the oracle cannot accidentally compare Go to itself.
- The shell cutover should use `skills/do-work/tools/do-work-cli.sh` as the single build-on-demand launcher and preserve argv byte boundaries with `exec`; the two installer entry points may need their static bootstrap-print behavior before installation.
- The standalone `skills/do-work-toolbox/tools/audit-metrics/` module is redundant with `internal/toolboxcommands` and can be removed only after its focused behavior is represented in the retained Go tests and maintainer gate.
- Board parity work is localized to queue-kanban command wiring, timestamp/frontmatter authority wording, verify findings/remedies, and their focused tests. Live documentation restatements span the root README plus do-work, board, knowledge, and toolbox guidance.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-knowledge/hooks/memory-session-start.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work-knowledge/hooks/memory-stop-capture.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work-knowledge/scripts/install-memory-hooks.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work-knowledge/scripts/lexical-memory-recall.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work-toolbox/scripts/architecture-report-preflight.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work-toolbox/scripts/generate-report-image-batch.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work-toolbox/scripts/generate-report-image.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work-toolbox/scripts/install-last30days.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work-toolbox/scripts/publish-portfolio-summary.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/hooks/session-start.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/scripts/add-local-git-exclude.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/scripts/atomic-download.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/scripts/audit-archive-timestamps.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/scripts/capture-screenshot.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/scripts/cleanup-req-reservations.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/scripts/handoff-state-survey.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/scripts/protected-inventory.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/scripts/repair-req-timestamps.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/scripts/run-blocked-check.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/scripts/show-commit-diff.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/scripts/stage-exact-deletion.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/checks/archive-collision.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/checks/associate-files.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/checks/blanked-req-scan.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/checks/preflight.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/checks/qualify.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/checks/record-commit-hash.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/checks/scope-drift.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/checks/uncommitted-inventory.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/do-work-cli.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/do-work-update.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/estimate-p50.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/fetch-upstream-archive.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/install-do-work-suite.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/replace-text-section.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/select-simple-reqs.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/validate-suite-manifest.sh` (modify) — preserve path as a thin compatibility launcher
- `tools/fetch-upstream-archive.sh` (modify) — preserve path as a thin compatibility launcher
- `tools/install-do-work-suite.sh` (modify) — preserve path as a thin compatibility launcher
- `tools/replace-text-section.sh` (modify) — preserve path as a thin compatibility launcher
- `tools/validate-suite-manifest.sh` (modify) — preserve path as a thin compatibility launcher
- `_dev/tests/fixtures/shipped-shell-command-map.tsv` (new) — parity/thinness evidence
- `_dev/tests/shipped-shell-parity.sh` (new) — parity/thinness evidence
- `_dev/tests/shipped-shell-thinness.sh` (new) — parity/thinness evidence
- `_dev/tests/contract-regressions.sh` (modify) — preserve path as a thin compatibility launcher
- `_dev/tests/maintainer-verify.sh` (modify) — preserve path as a thin compatibility launcher
- `_dev/tests/prescribed-shell-canonicalization.sh` (modify) — preserve path as a thin compatibility launcher
- `_dev/tests/prescribed-shell-cases/install-memory-hooks.sh` (modify) — preserve path as a thin compatibility launcher
- `_dev/tests/staged-skills-contract.sh` (modify) — preserve path as a thin compatibility launcher
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/git_helpers.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/git_helpers_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/handoff.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/handoff_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/publication.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/publication_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_churn_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_inventory_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_metrics_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/last30days.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/last30days_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/portfolio.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/portfolio_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `README.md` (modify) — close the corresponding Go, test, board, or authority contract
- `_dev/primes/prime-kanban-board.md` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/actions/capture.md` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work/actions/work-reference.md` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work-board/tools/queue-kanban/frontmatter_cli.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work-board/tools/queue-kanban/frontmatter_cli_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work-board/tools/queue-kanban/timestamp.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work-board/tools/queue-kanban/timestamp_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modify) — close the corresponding Go, test, board, or authority contract
- `skills/do-work-toolbox/tools/audit-metrics/.gitignore` (delete) — retire standalone audit implementation
- `skills/do-work-toolbox/tools/audit-metrics/churn.go` (delete) — retire standalone audit implementation
- `skills/do-work-toolbox/tools/audit-metrics/churn_test.go` (delete) — retire standalone audit implementation
- `skills/do-work-toolbox/tools/audit-metrics/distribution.go` (delete) — retire standalone audit implementation
- `skills/do-work-toolbox/tools/audit-metrics/distribution_test.go` (delete) — retire standalone audit implementation
- `skills/do-work-toolbox/tools/audit-metrics/git_support.go` (delete) — retire standalone audit implementation
- `skills/do-work-toolbox/tools/audit-metrics/go.mod` (delete) — retire standalone audit implementation
- `skills/do-work-toolbox/tools/audit-metrics/inventory.go` (delete) — retire standalone audit implementation
- `skills/do-work-toolbox/tools/audit-metrics/inventory_test.go` (delete) — retire standalone audit implementation
- `skills/do-work-toolbox/tools/audit-metrics/main.go` (delete) — retire standalone audit implementation
- `skills/do-work-toolbox/tools/audit-metrics/prime-audit-metrics.md` (delete) — retire standalone audit implementation
- `_dev/tests/prescribed-shell-cases/generate-report-image.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/generate-report-image-batch.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh` (modify) — assert behavior instead of shell-owned domain constants
- `_dev/tests/prescribed-shell-cases/add-local-git-exclude.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/atomic-download.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/capture-screenshot.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/cleanup-req-reservations.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/install-last30days.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/lexical-memory-recall.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/protected-inventory.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/publish-portfolio-summary.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/qualify.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/run-blocked-check.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/show-commit-diff.sh` (modify) — preserve behavioral assertions through the canonical launcher
- `_dev/tests/prescribed-shell-cases/stage-exact-deletion.sh` (modify) — preserve behavioral assertions through the canonical launcher

**Files I will NOT touch:** Version/changelog release files, unrelated queue items, and the pre-existing untracked REQ-418 run artifacts.

**Acceptance criteria (restated from REQ):**
- [ ] All 41 tracked shell paths remain present but contain launcher/bootstrap mechanics only, with no Python, jq, or domain implementation branches.
- [ ] A legacy-first differential oracle proves exact status, ordered facts/output, argv, Git/index/private state, and filesystem parity across earned fixtures and mutation dimensions.
- [ ] Every non-clean core-helper condition emits typed per-path findings with condition-specific next and verification argv.
- [ ] Portfolio refusal text, trailing-option usage families, memory-hook installation, P50 estimation, archive-collision, blanked-REQ scan, and other compatibility seams match retained contracts.
- [ ] BKB outside-root init, master-index status counts, and the committed audit matrix pass in text/JSON with read-only guarantees.
- [ ] SessionStart, board compiler/fallback, and recovery-command authority restatements identify the Go owner and emit runnable canonical argv.
- [ ] The standalone audit-metrics module is deleted only after its behavior and tests are absorbed into do-work-cli; the maintainer gate has no separate audit lane.
- [ ] Final install/update succeeds without Python or jq, every managed Just command runs without an LLM, aliases remain unchanged, and actionable commands avoid rescans.
- [ ] Focused Go/board tests, shell thinness/parity suites, contract regressions, and the unpiped canonical maintainer gate all pass.

## Decisions

### D-01 — DECIDE & STATE: move isolated legacy fixtures to the canonical runtime

The prescribed-shell fixture family copied retained scripts into isolated directories, intercepted shell utilities such as `cp`/`mv`, or required domain constants to remain literally shell-owned. REQ-420 explicitly requires launcher-only retained paths, so those implementation-shaped contracts contradicted the accepted behavior. Expand Scope and `write_set` to all 18 fixture owners, retain their behavioral assertions through the canonical CLI, and do not make Go shell out or preserve domain code merely to satisfy obsolete injection seams.

### D-02 — DECIDE & STATE: serialize the overlapping REQ-478 authority edit

REQ-478 was claimed by another live writer during this build and modifies `skills/do-work/actions/capture.md` and `skills/do-work/actions/work-reference.md`, both declared by REQ-420. Add REQ-478 as a dependency so a restarted default run cannot resume REQ-420's authority sweep until that foreign claim completes; fan-out remains 1 while the shared checkout is dirty.

## Decisions

- **D1 — Preserve the legacy oracle in executable fixtures.** The existing per-script behavior cases are retained as the pre-cutover expected status/output/effect oracle. The new parity gate runs those fixtures through every retained path after cutover and separately proves every path maps to a registered canonical command. This avoids comparing the Go command to a shim that already invokes the same Go command.
- **D2 — Treat bootstrap text as the sole pre-install exception.** The two installer mirrors may retain the static quoted bootstrap heredoc because the canonical binary cannot exist yet. The thinness gate otherwise rejects Python, jq, and shell domain implementations.
