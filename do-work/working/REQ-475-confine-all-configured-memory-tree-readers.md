---
id: REQ-475
title: '[impact-critical] Review fix: Confine all configured Memory tree readers'
status: claimed
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
claimed_at: 2026-09-03T21:43:31Z
route: C
planning_at: 2026-09-03T21:53:38Z
exploration_at: 2026-09-03T21:53:38Z
dispatch_at: 2026-09-03T21:55:13Z
implementation_at: 2026-09-03T22:12:45Z
builder_handback_at: 2026-09-03T22:13:52Z
write_set:
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands_test.go
  - skills/do-work-knowledge/actions/memory.md
  - skills/do-work-knowledge/actions/memory-reference.md
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
