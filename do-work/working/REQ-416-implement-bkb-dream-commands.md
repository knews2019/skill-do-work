---
id: REQ-416
title: 'Implement deterministic BKB and Dream commands'
status: claimed
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-415]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
claimed_at: 2026-09-01T04:10:12Z
route: C
estimate:
  p50_active_minutes: 95
  confidence: low
  calculated_at: 2026-09-01T04:18:00Z
  basis:
    - Route C
    - 16-file write set
    - 7 new files
    - 4 subsystems involved
    - dependency depth 11
    - cross-route regression gates
    - full-suite verification
---

# Implement Deterministic BKB and Dream Commands

## What
Move deterministic knowledge-base and Dream scans into `do-work-cli` while retaining LLM-only judgment phases in actions.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Accepted a three-task plan for one typed knowledge-command family, safe BKB scaffolding/structural scans, exactly seven read-only Dream scans, narrow recipes, and action consumption with semantic judgment retained.
- [x] **[APPLY]:** Added the four typed commands, exact safe scaffold transaction, BKB/Dream scans, narrow recipes, and action delegation in the frozen 16-file scope.
- [x] **[UNIFY]:** Reviewed all 16 files; focused/race/full Go, vet, exact Go 1.25, recipe/action/contract/staged/install/update, qualification/scope, diff, and final canonical gates pass.

## Detailed Requirements
- Implement `bkb-init`, `bkb-status`, and `bkb-lint-structure` as direct commands and flat recipes.
- Implement Dream’s seven deterministic scans behind `dream-scan`.
- Return typed actionable findings in text and JSON with exact next and verification commands.
- Leave contradiction resolution, synthesis, and cluster design to the existing LLM actions, which must consume canonical command results.

## Constraints
- Preserve knowledge file formats and never convert judgment-heavy phases into brittle heuristics.

## Dependencies
Depends on REQ-415 (hook/runtime migration precedes package-specific knowledge commands).

## Builder Guidance
Certainty level: Firm for deterministic scans; explicitly retain action ownership for semantic judgment.

## Red-Green Proof
**RED prompt/case:** Run BKB and all seven Dream scan fixtures through absent CLI commands, including clean, malformed, and finding cases.
**Why RED now:** Deterministic phases are specified in action prose rather than exposed as one stable no-LLM interface.
**GREEN when:** Every scan is directly runnable, text/JSON agree, findings are actionable, and semantic resolution remains delegated to the LLM action.
**Validation:** User confirmed via the supplied implementation plan.

## Triage

**Route: C** — Complex

**Reasoning:** Four public commands span a mutating scaffold transaction, two read-only scan families, action delegation, recipes, typed result projection, and an explicit deterministic-versus-semantic ownership boundary.

**Estimate:** p50 95 active minutes, low confidence, based on a 16-file write set with seven new files, four subsystems, cross-route contracts, and full-suite verification.

## Plan

1. Add and register a standard-library `internal/knowledgecommands` family with RED fixtures for command availability, option grammar, target confinement, read-only scans, actionable text/JSON parity, and deterministic ordering.
2. Implement safe BKB init/status/structural-lint behavior and make the BKB action consume canonical results while retaining semantic checks and report judgment.
3. Implement exactly seven Dream scans, make Dream consume one canonical worklist while retaining lock/consent/consolidation ownership, add the four required recipes, reconcile the CLI prime/contracts, and run focused/race/full/compatibility/install/update/canonical gates.

## Decisions

- `bkb-lint-structure` and `dream-scan` are read-only; action-owned semantic judgment, reports, consent, locks, consolidation, and repairs remain outside the command family.
- The Dream bonus newer-source probe is not one of the seven required scans.
- `bkb-init` preserves the current exact scaffold/fill-gaps behavior, supports dry-run and meaningful commit, never overwrites, and owns rollback only for invocation-created paths.
- Only the four REQ-416 recipes land now; REQ-419 retains the broad interface/help/hostile-quoting migration.
- Files/findings/pairs are normalized, sorted, and deduplicated for stable text/JSON output.
- REQ-416 must not edit `internal/gittransaction`; REQ-447 exclusively owns that shared seam in this wave.

## Scope

**Files I will touch:**
- `_dev/tests/contract-regressions.sh`
- `justfile`
- `skills/do-work-board/justfile.template`
- `skills/do-work-knowledge/actions/bkb.md`
- `skills/do-work-knowledge/actions/bkb-reference.md`
- `skills/do-work-knowledge/actions/dream.md`
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go` (new)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan.go` (new)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/dream_scan.go` (new)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/dream_scan_test.go` (new)

**Files I will NOT touch:** `internal/gittransaction`, queue-kanban atomic publication, broad help/guides/SKILL interface migration, retained shell utilities, audit-metrics, release metadata, or REQ-417–420 surfaces beyond the four narrowly required recipes.

**Acceptance criteria:**
- [x] All four commands are registered, direct, typed, actionable, and deterministic in text/JSON.
- [x] BKB scaffold/status/structural fixtures preserve exact formats/effects while semantic lint remains action-owned.
- [x] Dream exposes exactly seven read-only scans and its action consumes the canonical worklist without fallback or lock leakage.
- [x] Four flat recipes work in source and installed topology with collision/delegation contracts.
- [x] Exact 16-file scope and focused/race/full/compatibility/install/update/canonical gates pass.

## Implementation Summary

Registered `bkb-init`, `bkb-status`, `bkb-lint-structure`, and `dream-scan` through one typed `internal/knowledgecommands` family. BKB init owns the exact scaffold, fill-gaps, dry-run, exact-path commit, Git/standalone setup, no-overwrite publication, nested unsafe-entry refusal, and invocation-only rollback. Status and structural lint are deterministic/read-only. Dream exposes exactly seven normalized, sorted, read-only finding classes.

Four source/installed recipes now invoke the canonical launcher. BKB and Dream actions consume command JSON without fallback while retaining semantic judgment, reports, lock/consent, consolidation, reindexing, and audit ownership.

**Builder commit:** `fd59155568d304ad37002367c78856ef442c85e3`

**Integration range:** `5e341946..3519b315`

**Files changed:** the exact 16 paths declared in Scope.

## Qualification and Testing

- Mechanical qualification and scope drift passed with the exact 16-file builder range.
- Owner focused race package tests and the full contract regression suite passed after integration.
- Builder full Go, vet, exact Go 1.25, real recipe invocation, staged/install/update, diff hygiene, and final canonical maintainer gate passed; the optional browser lane skipped normally.
- RED evidence begins at absent commands. GREEN fixtures cover exact scaffold/reference bytes, dry-run/preflight/commit/rollback, nested unsafe entries, BKB structural/status cases, all seven Dream classes, deterministic/read-only results, source/installed recipes, action delegation, and lock release on command failure.

## Review — Initial

**Overall: 50%** | 2026-09-01

| Dimension | Score |
|-----------|-------|
| Requirements | 70% |
| Code Quality | 55% |
| Test Adequacy | 50% |
| Scope | 100% |
| Risk | Critical |
| Acceptance | Fail |

**Important findings:**
- F1 — impact-critical — pathname-based scaffold creation and rollback can escape the root or delete replacement objects after parent/object swaps.
- F2 — impact-user-visible — documented absolute BKB/Dream targets fail, and symlink-spelled roots can emit escaping affected paths.
- F3 — impact-user-visible — structural lint silently drops dangling topic-index article links.
- F4 — impact-user-visible — Dream recursively widened the frozen flat `<wiki>/*.md` scan scope.

Two Minor findings cover missing affected paths for dangling Dream index findings and a focused matrix narrower than the handback claims. All findings enter the single remediation pass.

*Reviewed independently; full evidence is in `do-work/runs/work-2026-08-31-165510/REQ-416-review.md`.*

## Remediation

The sole remediation pass committed `1ff0dd7227885ca797203d17df969f1356ff4877` and changed eight existing knowledgecommands files within the accepted 16-file scope.

- F1: `os.Root`-confined exclusive creation, directory revalidation, and inode/device identity-owned cleanup now cover Git and standalone flows, including owned `.git`; deterministic parent/object swap hooks prove no outside write or replacement deletion.
- F2: documented absolute/relative targets work, symlink-spelled roots normalize physically, and next/verification argv retain the user target.
- F3: dangling topic-index targets produce actionable `BKB-TOPIC-DANGLING-ENTRY` findings.
- F4: Dream scans only direct `<wiki>/*.md` pages and refuses duplicate normalized stems deterministically.
- Minor closure: dangling Dream index findings name the index path, and the focused differential matrix now covers the review's target/path/depth/date/vocabulary/clean/malformed/race cases.

**Remediation integration commit:** `684d4d828ecaaf6db190fe6642d56e3e12b41641`

*Generated from `do-work/runs/work-2026-08-31-165510/REQ-416-remediation-handback.md`.*

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
