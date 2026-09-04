---
id: REQ-475
title: '[impact-critical] Review fix: Confine all configured Memory tree readers'
status: pending
priority: now
created_at: 2026-09-01T08:32:57Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-417]
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-417]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-417
sweep: true
sweep_key: memory-configured-tree-readers-not-rooted
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-09-03T21:44:31Z
  basis:
    - Route C
    - 8-file write set
    - 3 subsystems involved
    - 4 acceptance criteria
    - persistence changes
    - cross-route regression gates
    - full-suite verification
route: C
planning_at: 2026-09-03T21:53:38Z
exploration_at: 2026-09-03T21:53:38Z
dispatch_at: 2026-09-03T21:55:13Z
implementation_at: 2026-09-03T22:12:45Z
builder_handback_at: 2026-09-03T22:13:52Z
integration_at: 2026-09-03T22:20:33Z
testing_at: 2026-09-03T22:24:22Z
review_at: 2026-09-03T22:24:22Z
status_changed_at: 2026-09-04T12:36:12Z
commit: c6d457473d24cdb188070709100884f019323ebc
write_set:
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands_test.go
  - skills/do-work-knowledge/actions/memory.md
  - skills/do-work-knowledge/actions/memory-reference.md
heavy_verified_at: 2026-09-04T12:36:12Z
heavy_verified_revision: c0d8ce1cb44cc1830b167214c018d76ba87baffc
---

# Confine All Configured Memory Tree Readers

## What

Apply one rooted, no-follow, bounded-read contract to every command that reads configured Memory state. Done means no read-only or mutating Memory surface can follow a working-file, log-file, ledger, sentinel, or directory link outside the validated memory tree, and typed output never contains bytes from a refused object.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that owns this configured Memory reader root cause.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Inventoried every configured Memory-tree reader and designed one handler-scoped, bounded, rooted/no-follow acquisition contract with exact-path refusal evidence.
- [x] **[APPLY]:** Migrated recall, status, audit, mutator, ledger, sentinel, and log-directory reads; added adversarial link/special/limit/root fixtures and documented finite transport ceilings.
- [x] **[UNIFY]:** Reviewed the exact four-file builder diff; focused, race, vet, full-module, contract, formatting, and diff checks pass with no debug, BKB, lifecycle, generated, or release drift.

## Instances

- [ ] Broad `memory-recall` plain-reads configured working memory and globbed daily logs, unlike the lexical path's rooted refusal.
- [ ] `memory-status` plain-reads configured working memory, daily logs, and ledger evidence.
- [ ] `memory-audit --engine memory` plain-reads configured working memory and daily logs.

## Requirements

- Reuse one repository-rooted, no-follow, regular-file reader for every configured Memory read surface, including directory enumeration.
- Reject linked or special working files, log files, ledgers, sentinels, and log directories before reading their bytes; bound accepted reads to the documented store limits.
- Ensure text and JSON findings name the refused configured path without including target bytes.
- Add adversarial fixtures for each reader and each link position, plus ordinary read-only parity and byte-preservation coverage.

## Red-Green Proof

**RED prompt/case:** Point configured working-memory and daily-log paths at synthetic files outside the memory tree, then run broad recall, status, and memory audit in text and JSON.
**Why RED now:** All three commands succeed and include outside fixture values while lexical recall rejects the equivalent linked log.
**GREEN when:** Every configured reader refuses the linked/special path before reading it, returns no fixture content, and ordinary stores retain deterministic output.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Full Context

See `do-work/runs/work-2026-08-31-165510/REQ-417-rereview.md`.

---
*Source: REQ-417 fresh re-review residual finding 1.*

## Triage

**Route: C** - Complex

**Reasoning:** This security-sensitive confinement change spans every configured Memory reader, shared rooted/no-follow primitives, enumeration and size limits, text/JSON redaction, and adversarial parity tests. Planning and exploration are required before declaring the exact reader set.

**Planning:** Required

## Plan

1. Add adversarial RED coverage for linked and special working, log, ledger, sentinel, directory, and outside-root paths across broad/lexical recall, status, and Memory audit, including typed/text/JSON canary absence and byte preservation.
2. Introduce one handler-scoped, repository-contained Memory reader with bounded directory enumeration and bounded, identity-checked rooted regular-file reads that return metadata and distinguish missing from refusal.
3. Migrate recall, status, and Memory audit to that reader, split ledger parsing from acquisition, and preserve ordering, scoring, classifications, and BKB acquisition behavior.
4. Document the finite Memory-store read limits and repository-root policy, inventory away remaining direct configured-tree pathname reads, then run focused, race, vet, contract, and full-module verification.

**Plan validation:** Every requirement maps to the reader primitive, its three consumer families, adversarial output proof, or the documented bounds. No board, hook, git-transaction, generated artifact, or screenshot surface is included.

*Generated from delegated exploration; full evidence: `do-work/runs/work-2026-09-03-214500/REQ-475-exploration.md`.*

## Exploration

The defect is one authority split: rooted/no-follow helpers protect mutation and lexical seams, while broad recall, status, and Memory audit still use direct `os.ReadFile`, `os.Stat`, or `filepath.Glob` paths. A single opened-root reader must own both enumeration and final-file acquisition; preflight `Lstat` calls alone would preserve pathname races. Missing objects remain ordinary, while links, special objects, oversize objects, identity changes, and outside-repository configured roots must become path-bearing refusals without partial evidence or ledger writes.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands_test.go` only if the real-runtime text/JSON projection needs a separate assertion
- `skills/do-work-knowledge/actions/memory.md`
- `skills/do-work-knowledge/actions/memory-reference.md`

**Acceptance criteria:** Every configured Memory-tree reader uses one bounded rooted observation; linked, special, oversize, changed, and outside-root objects are refused by configured path without target bytes in typed, text, or JSON output; ordinary recall/status/audit semantics and BKB behavior remain stable.

## Pre-Flight

**Git:** The wave baseline was clean at `c27d349a` after the three claims, estimates, run manifest, and briefs were committed.

**Tests:** Direct canonical fast gate passed and was recorded at the shared wave baseline before any builder dispatch.

**Dependencies:** REQ-417 is completed. Its referenced scratch re-review artifact is absent from the current tree and Git object inventory, so the durable archive summary and remediation diff are the recoverable source evidence.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands_test.go`
- `skills/do-work-knowledge/actions/memory.md`
- `skills/do-work-knowledge/actions/memory-reference.md`

**What was done:** Added one opened-root Memory reader with bounded identity-checked file and directory operations, migrated every configured Memory consumer to it, separated ledger parsing from acquisition, propagated exact configured child failures without partial/canary evidence, enforced repository-contained configured roots, and documented the finite transport ceilings. Builder commit: `0f288b7ccbf454c7c73935a8dd6aa3b8f211932b`.

**Builder verification:** Handler-level RED reproduced outside-byte disclosure and outside-root acceptance. Focused, race, vet, full-module, contract, inclusive-limit, and diff checks are green; full evidence is in `do-work/runs/work-2026-09-03-214500/REQ-475-handback.md`.

## Review

**Verdict:** Request changes — one remediation pass required before final testing/finalization.

- **Critical — impact-critical:** `appendMemoryLedger` retained a direct `os.OpenRoot(memoryAbsolute)` after configured recall. A post-scan root swap to an outside symlink could therefore redirect the best-effort ledger mutation outside the validated Memory tree, bypassing the new root identity check and contradicting the no-follow contract for mutating surfaces.

## Remediation

Builder commit `1ace19970e242d4c61409a53e81ab78800fb8065` routes ledger append acquisition through `openMemoryRoot`. A deterministic pre-open seam swaps the already-scanned root to an outside symlink: RED appended a recall event to the outside canary ledger; GREEN leaves those bytes unchanged. Focused/race/vet/full-module/contracts and diff hygiene pass.

## Re-Review

**Verdict:** Pass. The cumulative `63045c9e..c6d45747` implementation now applies the same verified-root boundary to the post-recall ledger mutation as to every configured read. All configured Memory file/directory reads are bounded and rooted, refused target bytes stay out of typed/text/JSON output, ordinary recall/status/audit semantics and BKB acquisition remain intact, and the documented limits match executable constants and inclusive-boundary tests.

## Qualification

Passed — mechanical qualification accepted the cumulative `63045c9e..c6d45747` range. The intervening REQ-485 commits are independently scoped and recorded; REQ-475's implementation is limited to its four declared files and its two-file remediation subset, with no builder-authored `do-work/` state.

## Testing

**Tests run:** focused and race `internal/knowledgecommands`; full do-work-cli `go vet ./...` and `go test ./...`; contract regressions; direct `bash _dev/tests/maintainer-verify.sh`.

**Result:** Focused, race, full-module, contracts, and the fast canonical gate pass on final merge `c6d457473d24cdb188070709100884f019323ebc`; the exact gate revision is recorded. `memory_commands_test.go` matches the declared CLI heavy-test surface, so finalization requires exact-revision heavy permission.

**Red-green validation:** The outside-object/root disclosure matrix, special/limit boundaries, and deterministic ledger-root swap all failed in their pre-fix states and pass on the merged implementation with byte-preservation and canary-absence assertions.

## Heavy Verification Plan

- `mode`: `historical-revalidation`
- `source_ranges`:
  - `63045c9e03062858ef74c65ca631bcba8b5c832f..a0207eaa0285fcbdfa66db168d093d78d0c5737f`
  - `11061016a89b88ae3b14ee83ec748ca03cc1d132..c6d457473d24cdb188070709100884f019323ebc`
- `manifest_revision`: `c0d8ce1cb44cc1830b167214c018d76ba87baffc`
- `execution_revision`: `c0d8ce1cb44cc1830b167214c018d76ba87baffc`
- `manifest_path`: `_dev/tests/heavy-lanes.json`
- `forced_all`: `false`
- `uncertain`: `false`
- `uncovered_paths`: `[]`
- `changed_paths`:
  - `skills/do-work-knowledge/actions/memory-reference.md`
  - `skills/do-work-knowledge/actions/memory.md`
  - `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go`
  - `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands_test.go`
- `selected_lanes`:
  - `do-work-cli-integrations`: `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
  - `staged-skills`: `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
  - `updater`: `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
  - `installer`: `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## Open Questions

- [x] Run the selected heavy lane commands at execution revision `c0d8ce1cb44cc1830b167214c018d76ba87baffc`; did every command exit 0? → Confirmed: Yes
  Recommended: Yes
  Also: No — report the failing lane


## Answer Notes

- 2026-09-03 - [ ] Run `bash _dev/tests/maintainer-verify.sh --heavy` at `c6d457473d24cdb188070709100884f019323ebc`; did it exit 0?: No — exit 1; update-script-behavior.sh took 55s, exceeding the under-30s limit
- 2026-09-04 - [ ] Run the selected heavy lane commands at execution revision `c0d8ce1cb44cc1830b167214c018d76ba87baffc`; did every command exit 0?: Confirmed: Yes

## Heavy Verification Result

Target revision: `c6d457473d24cdb188070709100884f019323ebc`
Execution revision: `c0d8ce1cb44cc1830b167214c018d76ba87baffc`

- do-work-cli-integrations: exit 0 — `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- staged-skills: exit 0 — `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- updater: exit 0 — `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- installer: exit 0 — `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`
