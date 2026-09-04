---
id: REQ-417
title: 'Implement interview and deterministic memory store commands'
status: completed-with-issues
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-416]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
claimed_at: 2026-09-01T06:08:08Z
route: C
estimate:
  p50_active_minutes: 105
  confidence: low
  calculated_at: 2026-09-01T06:08:08Z
  basis:
    - Route C
    - 17-file write set
    - 4 new files
    - 12 command groups
    - persistence, privacy, and rollback changes
completed_at: 2026-09-01T08:34:36Z
commit: ecf77a3da1751d170c22ae94b782e1354337c67b
kb_status: promoted
kb_entry: REQ-417-implement-interview-and-deterministic-me.md
---

# Implement Interview and Deterministic Memory Store Commands

## What
Expose all deterministic interview and memory operations through `do-work-cli`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Characterized both stores and froze a 17-file Route C plan with explicit semantic/action and deterministic/CLI ownership.
- [x] **[APPLY]:** Added twelve commands, twelve recipes, action delegation, deterministic store operations, and private-untracked transaction support; applied the single allowed remediation.
- [x] **[UNIFY]:** Reviewed the exact builder and remediation scopes; focused, race, full, vet, compatibility, Windows, command parity, contract/install/update, and canonical gates pass.

## Detailed Requirements
- Implement interview list, status, export, ingest, reset, and versions commands.
- Implement memory remember, forget, recall, status, bootstrap, and audit commands.
- Preserve store formats, ordering, redaction, deduplication, version semantics, and atomic mutations.
- Provide direct commands, flat Just recipes, text/JSON parity, dry-run where meaningful, optional commit, and actionable findings.

## Constraints
- Natural-language knowledge actions remain aliases and delegate deterministic phases to the CLI.

## Dependencies
Depends on REQ-416 (knowledge command conventions and scan result integration).

## Builder Guidance
Certainty level: Firm. Characterize rendering and store mutations before migration.

## Red-Green Proof
**RED prompt/case:** Exercise each interview and memory operation against representative stores, malformed data, duplicate records, redaction cases, dry-run, and rollback.
**Why RED now:** These operations are not uniformly available as stable direct Go commands.
**GREEN when:** Every operation has fixture-proven deterministic output/effects, matching text/JSON, and its action alias delegates without free-form fallback.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

## Implementation Summary

`do-work-cli` now exposes six deterministic Interview commands and six deterministic Memory commands through typed text/JSON results, matching source and installed flat recipes, and action delegation that retains consent and semantic judgment. Interview templates remain data; export, ingest, reset, and version operations use deterministic multi-file plans. Memory remember, forget, recall, status, bootstrap, and audit preserve the documented store formats, with tracked/private transaction separation and exact commit behavior.

The implementation also adds a private-untracked target seam to the shared Git transaction layer. The remediation binds transaction snapshots and publication identities, prevents remembered/forgotten text from entering ledger query fields, JSON-escapes template substitutions, restores reset backlinks and complete traversal failures, completes Memory audit probes, and expands malformed/duplicate/collision coverage.

**Files changed:**
- `_dev/tests/contract-regressions.sh`
- `justfile`
- `skills/do-work-board/justfile.template`
- `skills/do-work-knowledge/actions/interview.md`
- `skills/do-work-knowledge/actions/interview-reference.md`
- `skills/do-work-knowledge/actions/memory.md`
- `skills/do-work-knowledge/actions/memory-reference.md`
- `skills/do-work-knowledge/actions/memory-value.md`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go`
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands_test.go`

**Builder commit:** `01d920ddf6c71e8b93683fcabf60702b61e499c7`

**Initial integration:** `440682f2f3ce77b97b1902dd6347cdb0655b8afb`

**Remediation commit:** `aa5e38bd32e8847ab90956d45d98a0066a177832`

**Final integration:** `ecf77a3da1751d170c22ae94b782e1354337c67b`

## Qualification and Testing

- Focused and race tests for `knowledgecommands` and `gittransaction`, full uncached Go tests, vet, exact Go 1.25 compatibility, and Windows cross-compiles pass.
- All twelve real commands dispatch with text/JSON parity; shipped Interview export rendering and retained-script lexical recall differential fixtures pass.
- Source/installed recipes, staged-skill, install/update, aggregate contracts, exact scope, diff hygiene, and canonical maintainer verification pass. The browser lane had its standard no-browser skip.

## Remediation

The initial independent review scored 50%, Critical risk, Acceptance Fail with seven Important findings. The one allowed remediation added adversarial regressions and closed transaction preimage/final identity, ledger privacy, renderer null/escaping, reset completeness/backlinks, audit semantics, and most matrix gaps. The final re-review confirms those closures but found the same no-follow reader class on broad recall, status, and Memory audit, plus missing durable BKB-audit coverage.

## Review

**Overall: 50%** | **Risk: Critical** | **Acceptance: Fail**

REQ-417 completes with issues because its single remediation is exhausted. Every residual Important finding is routed fold-first:

- `impact-critical` configured Memory reader confinement → REQ-475.
- `impact-negligible` durable BKB audit matrix → REQ-476.

Full evidence is in `do-work/runs/work-2026-08-31-165510/REQ-417-review.md` and `REQ-417-rereview.md`.

## Lessons Learned

- Applying a rooted no-follow helper to mutators is insufficient when read-only status, recall, or audit paths independently enumerate and read the same configured store.
- Manual acceptance evidence cannot replace a committed regression for a `tdd: true` requirement; coverage must exercise every named engine, not only its sibling mode.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
