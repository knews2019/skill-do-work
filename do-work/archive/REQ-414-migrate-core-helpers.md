---
id: REQ-414
title: 'Migrate remaining core checks, publication helpers, Git helpers, and surveys'
status: completed-with-issues
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
completed_at: 2026-09-01T02:26:12Z
commit: 0689970c
kb_status: pending
---

# Migrate Remaining Core Checks, Publication Helpers, Git Helpers, and Surveys

## What
Move all remaining core utility domain logic into `do-work-cli` subcommands.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Accepted a three-task plan for leaf helper commands over existing domain authorities, condition-complete runtime/publication ownership, and real-command typed parity.
- [x] **[APPLY]:** Implemented all 17 typed helper commands, reused the two existing authorities, and removed archivefetch/nextselection shell implementation dependencies in the approved 35-file scope.
- [x] **[UNIFY]:** Reviewed all 35 changed files; focused/full Go, vet, exact Go 1.25, Windows compile, legacy characterization, scope/diff hygiene, and the canonical maintainer gate pass on the builder branch.

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

## Exploration

The remaining surface is 19 shipped core utilities: 17 require public CLI commands and two already map to typed collision/damaged-record authorities. The smallest honest design keeps `internal/corehelpers` as the leaf handler/mechanics package while shared blocked-process execution remains in `nextselection`, HTTP transport remains in `archivefetch`, timestamp policy remains in `doctor`, and commit provenance remains in `requeststate`.

Exploration identified four contract choices. The defaults are recorded below: new CLI commands keep the suite-wide 0–4 envelope while carrying raw legacy statuses as typed evidence; parity means ordered facts, paths, statuses, modes, and side effects rather than byte-identical legacy prose; screenshot/download retain their observed private 0600 publication mode; and new-target publication stays outside pathname-only `gittransaction` rollback, with publication as the final commit point. HTTP behavior pins curl's initial request plus up to three eligible retries, while raw transfer status remains evidence under the typed envelope.

With those decisions, the initial 30-file boundary below was sufficient for the command family. The full Go gate then proved that removing archivefetch's shell dependency requires four suiteinstall caller/fixture edits; D-06 records the focused expansion to 34 files. The legacy contract gate then proved that its update behavior harness still asserted and injected the deleted shell seam; D-07 adds that one characterization file for a final 35-file boundary. Retained shell paths remain behaviorally unchanged for REQ-420, and no `resultmodel`, `commandruntime`, `gittransaction`, `atomicfile`, hook, action, Justfile, or release-metadata expansion is authorized without another focused RED and an owner-side scope revision before code.

*Generated by Explore agent*

## Decisions

### D-01: Keep publication outside pathname rollback

**Decision:** DECIDE & STATE — screenshot and download publication use rooted/private preparation and a final no-overwrite publish, with no later pathname-based rollback removal.

**Reasoning:** This satisfies REQ-414 without depending on unresolved REQ-457 and preserves the invariant that a later parent swap cannot redirect cleanup to another object's path.

### D-02: Preserve raw probe status as evidence

**Decision:** DECIDE & STATE — the in-process blocked runner preserves raw probe/timeout/isolation status, while the public CLI command uses the standard 0–4 result envelope.

**Reasoning:** One command cannot simultaneously use `commandruntime`'s canonical exit mapping and expose arbitrary process exit codes. Retained shell compatibility remains unchanged until REQ-420 owns shim translation.

### D-03: Define parity by observations and effects

**Decision:** DECIDE & STATE — parity requires identical ordered facts, paths, classifications, modes, bytes, Git/filesystem effects, and semantic statuses; new CLI text uses the standard typed envelope rather than byte-identical legacy prose.

**Reasoning:** The generic renderer cannot stream raw `git show`, TAB inventories, and surveys without expanding shared result/runtime scope. Legacy paths remain byte-compatible in this REQ.

### D-04: Retain private publication modes

**Decision:** DECIDE & STATE — screenshot and download destinations retain the observed private mode, normally 0600.

**Reasoning:** The current helpers publish their private temporary files and do not preserve source mode; keeping 0600 honors the captured permissions contract.

### D-05: Pin HTTP retry semantics

**Decision:** DECIDE & STATE — HTTP publication performs the initial request plus up to three eligible transient/429 retries, retains token precedence/redaction, and records raw transfer outcome in typed evidence.

**Reasoning:** This matches `curl --retry 3`'s attempt count without making a transport-specific status a second CLI exit authority.

### D-06: Remove stale suiteinstall download-script callers

**Decision:** EXPAND SCOPE — add the three suiteinstall call sites and their update fixture to the write set so they call the Go archivefetch contract directly and serve test archives over HTTP.

**Reasoning:** The full Go gate is a focused RED: four update tests still install an offline fake `atomic-download.sh`, while production suiteinstall still calls the compatibility locator. Keeping a no-op locator only to satisfy those callers would contradict the accepted requirement that archivefetch stop locating/executing the shell helper.

### D-07: Migrate the update characterization seam

**Decision:** EXPAND SCOPE — add only `_dev/tests/update-script-behavior.sh` so its 17 affected cases exercise the Go HTTP archivefetch path instead of injecting and asserting the removed curl/atomic-download seam.

**Reasoning:** Focused/full Go, vet, exact Go 1.25, and Windows compile pass after D-06; the legacy contract gate alone fails because its harness still demands the implementation dependency this REQ explicitly removes. Updating the characterization authority is required evidence, not production scope growth.

### D-08: Reconcile the blocked-probe orchestrator restatement

**Decision:** EXPAND SCOPE — add `skills/do-work/actions/work.md` so its selector contract describes the new in-process blocked probe instead of directing agents to a removed temporary-file and shell-helper seam.

**Reasoning:** The independent restatement sweep found a live orchestrator instruction that contradicts the migrated runtime. This is part of making the replacement authoritative, not a separate feature.

### D-09: Reconcile the blocked-probe lock-in assertion

**Decision:** EXPAND SCOPE — add `_dev/tests/contract-regressions.sh` so the canonical contract gate requires the selector-owned in-process blocked probe rather than the removed shell runner.

**Reasoning:** The canonical maintainer gate supplied the focused RED after D-08: its assertion still mandates the obsolete implementation seam and fails on the corrected orchestrator text.

## Scope

**Files I will touch:**
- `_dev/tests/contract-regressions.sh` (modified) — lock the selector-owned in-process blocked-probe contract.
- `skills/do-work/actions/work.md` (modified) — reconcile the selector's blocked-probe restatement with in-process ownership.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (new) — handler family and typed projections.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go` (new) — real command text/JSON and registration coverage.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go` (new) — preflight, qualify, and scope drift.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks_test.go` (new) — characterization and REQ-390 scope cases.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go` (new) — uncommitted/protected inventories and association.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go` (new) — NUL/rename/secret/quarantine fixtures.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/publication.go` (new) — screenshot and download command seams.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/publication_test.go` (new) — bytes, modes, collision, and parent-swap fixtures.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/git_helpers.go` (new) — exclude, commit display, and deletion staging.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/git_helpers_test.go` (new) — worktree, merge, and exact-index fixtures.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations.go` (new) — identity-revalidated reservation cleanup.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations_test.go` (new) — landed, aged, race, and unsafe-entry fixtures.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/handoff.go` (new) — structured handoff survey.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/handoff_test.go` (new) — branch/worktree/dirty-inventory survey fixtures.
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go` (new) — shared process runner contract.
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` (new) — status, timeout, descendant, and unrelated-process cases.
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go` (new) — verified process-group ownership.
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go` (new) — standard-library fail-closed platform boundary.
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified) — public command registration.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go` (modified) — in-process blocked runner integration.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go` (modified) — selector probe parity.
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch.go` (modified) — shared Go HTTP publication primitive.
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch_test.go` (modified) — HTTP retry/redaction/publication parity.
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_timestamps.go` (modified) — reusable policy-driven timestamp engine.
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_timestamps_test.go` (modified) — source-policy and shape fixtures.
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_repair.go` (modified) — shared timestamp policy consumption.
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_repair_test.go` (modified) — guarded repair parity.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified) — shared guarded provenance operation.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modified) — standalone and complete provenance guards.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (modified) — remove obsolete download-script request wiring.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go` (modified) — call Go archivefetch without a shell-helper locator.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go` (modified) — remove the remaining command-side locator call.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go` (modified) — replace the fake shell downloader with an HTTP archive fixture.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — helper ownership and verification map.
- `_dev/tests/update-script-behavior.sh` (modified) — exercise direct Go HTTP archive fetching without a fake shell downloader.

**Files I will NOT touch:** `resultmodel`, `commandruntime`, `gittransaction`, `atomicfile`, retained shell utilities, hooks, actions, Justfiles, release metadata, or later migration REQs. Any expansion requires a focused failing fixture and an owner-approved scope revision before implementation.

**Acceptance criteria (restated from REQ):**
- [x] All 17 remaining helper surfaces are registered and every shipped utility is explicitly mapped to a new or existing typed authority.
- [x] Preflight/check/inventory/association/provenance/survey semantics preserve ordered facts, statuses, and Git/filesystem effects, including REQ-390's first-path scope fix.
- [x] Publication, blocked-process, Git, timestamp, reservation, and download mechanics preserve exact safety, byte, mode, retry, redaction, and failure behavior without shell implementation dependencies.
- [x] Text and JSON findings are actionable from one typed observation set, retained compatibility scripts remain unchanged, and the full characterization/gate matrix passes.

## Implementation Summary

Added 17 registered typed commands for the remaining core checks, inventories, publication, Git, process, timestamp, reservation, provenance, and survey surfaces. Archive collision and damaged-record recovery remain on their existing typed authorities. `archivefetch` now performs direct Go HTTP publication and Git fallback without locating Bash/curl helpers; `nextselection` owns the blocked process tree in-process; timestamp policy remains shared through doctor; and commit provenance remains shared through requeststate.

The implementation preserves the recorded 0–4/raw-status boundary, ordered-fact parity, private 0600 publication mode, initial-plus-three HTTP retry behavior, and final no-overwrite publication with no pathname rollback. D-06 and D-07 removed stale suiteinstall callers and migrated the update characterization harness to local HTTP fixtures.

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified)
- `_dev/tests/update-script-behavior.sh` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/git_helpers.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/git_helpers_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/handoff.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/handoff_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/publication.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/publication_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations.go` (new)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_repair.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_repair_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_timestamps.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_timestamps_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go` (new)
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go` (new)
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go` (new)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**Builder commits:** `bf739a0073f11ced4fd14b920484ca46c76eafb6`, `1612f20091aeb98183826b34cb25db3d02daf5ac`, `b369e06aa4e4a34982f9faef0d759b305f5fde9a`, `65e1de4a68acced0d4daf5490fcce4b946e93d02`

**Cumulative implementation range:** `baf05854..65e1de4a`

**Integrated merge commits:** `4af7bcb5`, `ec8e80ea`, `0689970c`

## Qualification

- `DO_WORK_DIFF_RANGE=fd243be9..ec8e80ea bash skills/do-work/tools/checks/qualify.sh do-work/working/REQ-414-migrate-core-helpers.md` — passed; static-reference warnings were limited to Go package/test/platform files covered by the qualification exception categories.
- `DO_WORK_DIFF_RANGE=fd243be9..ec8e80ea bash skills/do-work/tools/checks/scope-drift.sh do-work/working/REQ-414-migrate-core-helpers.md` — passed with an exact 35-file match.
- The initial integrated qualifier exposed a self-match in the detector source; builder commit `b369e06aa4e4a34982f9faef0d759b305f5fde9a` removed the literal signature and the cumulative range passed on re-run.

## Testing

- Focused Go packages, full `go test -count=1 ./...`, and `go vet ./...` — passed on the integrated range.
- `bash _dev/tests/do-work-cli-go125-compatibility.sh` — passed.
- `bash _dev/tests/update-script-behavior.sh` — passed in an isolated authoritative run after a concurrent invocation contended during fixture startup.

## Review — Initial

**Overall: 50%** | 2026-09-01

| Dimension | Score |
|-----------|-------|
| Requirements | 48% |
| Code Quality | 52% |
| Test Adequacy | 30% |
| Scope | 100% |
| Risk | Critical |
| Acceptance | Fail |

**Important findings:**
- F1 — impact-critical — provenance verification does not prove the metadata commit changed only the frontmatter `commit:` line.
- F2 — impact-critical — protected inventory drops durable secret-quarantine persistence, rows, and status semantics.
- F3 — impact-user-visible — preflight and empty association omit retained findings/statuses.
- F4 — impact-user-visible — mutating helpers reject the required `--dry-run` surface.
- F5 — impact-user-visible — qualification loses relocation semantics and deterministic multi-file ordering.
- F6 — impact-user-visible — generic remediation placeholders and opaque handoff blobs are not actionable typed projections.
- F7 — impact-rule-change — the required real-command characterization matrix is absent.
- F8 — impact-rule-change — the blocked-probe orchestrator restatement and pending REQ-441 evidence still describe removed seams.
- F9 — impact-user-visible — HTTP retry delay/budget behavior is not compatible with the retained primitive.

**Acceptance:** Fail. Direct CLI replays reproduced empty-association success and dry-run usage failures; source-contract comparison reproduced the two critical safety regressions and the remaining parity gaps.

*Reviewed independently; full evidence is in `do-work/runs/work-2026-08-31-165510/REQ-414-review.md`. The sole remediation pass must close all nine findings before fresh re-review.*

## Remediation

The sole remediation pass committed `65e1de4a68acced0d4daf5490fcce4b946e93d02` and closed F1-F9 plus both adjacent Minor findings within the owner-approved 37-file boundary.

- F1 now proves an exact one-field metadata commit against `HEAD^..HEAD`, rejects re-staged body rewrites, and fails closed when Git object/delta guards cannot answer.
- F2-F3 restore durable protected-quarantine state, typed rows/statuses, exact `X` classification, dependency preflight findings, and empty-association refusal.
- F4 gives every migrated mutating surface a truthful side-effect-free `--dry-run`, exercised through the real binary.
- F5-F7 restore qualifier relocation/count/order semantics, finding-specific actions, structured handoff facts, and the all-17 text/JSON characterization matrix.
- F8 reconciles the live orchestrator text and its canonical lock-in assertion with the selector-owned in-process blocked runner; the owner separately blocked and annotated REQ-441's superseded evidence.
- F9 restores fixed two-second retries under a 60-second budget, removes the shorter whole-request timeout, and streams each attempt into a private rooted stage.
- The D-07 test seam now names REQ-414, and screenshot staging cleanup emits a typed warning without rolling back successful publication.

Focused/full Go tests, vet, exact Go 1.25, update behavior, contract regressions, the canonical maintainer gate, mechanical qualification, diff hygiene, and exact 37/37 scope matching pass. The cumulative implementation-only range is `baf05854..65e1de4a`; remediation merge commit is `0689970c`.

*Generated by work action from `do-work/runs/work-2026-08-31-165510/REQ-414-remediation-handback.md`.*

## Review — Fresh Re-review

**Overall: 50%** | 2026-09-01

| Dimension | Score |
|-----------|-------|
| Requirements | 70% |
| Code Quality | 68% |
| Test Adequacy | 52% |
| Scope | 100% |
| Risk | Critical |
| Acceptance | Fail |

F1-F5, F8-F9, and both original Minor findings are closed. F6 and F7 remain partial, and fresh direct replays found two retained-contract divergences.

**Remaining Important findings and routing:**
- impact-user-visible — combined Git porcelain `AD` is classified as readable `A` instead of deletion-first `D` → REQ-462 created.
- impact-critical — reservation cleanup trusts uncommitted requests in an unborn repository and does not revalidate final eligibility → REQ-463 created.
- impact-user-visible — family-wide fallback argv and opaque dirty-path blobs remain non-specific and non-structured → REQ-464 created.
- impact-rule-change — the all-17 matrix remains smoke-only rather than retained-versus-Go semantic/effect parity → REQ-465 created.

**Acceptance:** Fail. The staged-new-then-deleted inventory fixture and unborn-repository reservation fixture diverge from retained behavior, and a two-path handoff remains one multiline evidence blob.

**Lifecycle result:** `completed-with-issues` after the sole remediation and fresh re-review. All residual Important findings are preserved as dependency-linked sweep REQs.

*Re-reviewed independently; full evidence is in `do-work/runs/work-2026-08-31-165510/REQ-414-rereview.md`.*

## Lessons Learned

- Registration, renderer agreement, and broad green gates do not prove a migration preserved behavior; differential fixtures must compare exact retained statuses, ordered facts, paths, actions, and side effects across combined-state and authority-boundary cases.
- Destructive cleanup needs fresh authorization as well as object identity: revalidate the evidence that permits deletion and every eligibility predicate immediately before the mutation.

## Orientation

- **Implementation commit:** `0689970c`
- **Lifecycle result:** `completed-with-issues` after one remediation and one fresh re-review.
- **Residuals:** REQ-462 through REQ-465 retain inventory, reservation-authority, structured-projection, and differential-parity closure work.

*Generated by work action from the builder hand-back*

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
