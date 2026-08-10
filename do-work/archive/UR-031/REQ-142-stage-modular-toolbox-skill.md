---
id: REQ-142
title: "Stage the Modular Toolbox Skill"
status: completed
claimed_at: 2026-08-07T21:47:51Z
completed_at: 2026-08-07T21:56:38Z
commit: df35345
route: C
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: []
tdd: true
suggested_spec: refactor
depends_on: [REQ-136]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
kb_status: pending
kb_entry:
---

# Stage the Modular Toolbox Skill

## What
Create a staged `skills/do-work-toolbox` package preserving the current optional reviews, discovery, reporting, repository utilities, and companion installers.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add RED exact-route and package assertions; copy all retained toolbox actions/references/guides and required guardrails; prune board/memory/core-update ownership from the companion installer; author a toolbox-only router/help; convert queue and sibling action references explicitly.
- [x] **[APPLY]:** Staged all sixteen retained toolbox actions with their references, guides, and guardrails; authored the dedicated router/help; pruned board, knowledge, and updater ownership from companion installation; and made core/board/knowledge dependencies explicit.
- [x] **[UNIFY]:** Reviewed all 46 toolbox package files and the staged-suite contract; verified exact routing, non-empty assets, canonical command examples, cross-package resolution, installer ownership, project-output path exceptions, shell lint/syntax, full contracts, Go checks, Just parsing, and a clean diff check.

## Why
The user wants these capabilities retained, but they should not occupy the core queue router and context.

## Detailed Requirements
- Preserve validate-feedback, code-review, ui-review, present-work, ai-report, slop-check, quick-wins, scan-ideas, deep-explore, prime, inspect, note, stray-check, tidy-repo, tutorial, and companion installers not owned by board or knowledge.
- Move-by-copy every required reference, crew file, template, and guide.
- Add a dedicated router/help contract.
- Preserve links to core queue artifacts where actions read or report on URs/REQs.
- Add consistency tests ensuring every toolbox route exists exactly once and every runtime reference resolves in the staged suite.
- Do not remove or deactivate legacy copies yet.

## Constraints
- Retain all toolbox functionality; this program does not run a usage-pruning phase.
- Board setup belongs to board; memory setup belongs to knowledge; core self-update belongs to core.

## Dependencies
Requires REQ-136's suite contract.

## Builder Guidance
Certainty level: Firm. Preserve behavior while moving ownership; avoid redesigning individual actions during this stage.

## Red-Green Proof
**RED prompt/case:** Invoke each toolbox route through an isolated `skills/do-work-toolbox/SKILL.md` fixture and resolve every referenced runtime file.
**Why RED now:** These optional actions are routed and stored inside the monolithic skill.
**GREEN when:** All toolbox routes dispatch from the staged package, all references resolve, and legacy active behavior remains unchanged.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the retained toolbox action set.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*

## Triage

**Route: C** — sixteen independently useful actions span review, visual/report generation, discovery, repository hygiene, and external installers, with many cross-skill references.

## Scope

**Files I will touch:** `skills/do-work-toolbox/**`, `_dev/tests/staged-skills-contract.sh`, and REQ/release metadata.

**Files I will not touch:** active legacy toolbox files, board/memory/core installers, other staged ownership, application files, or unrelated dirty REQ-134 work.

## Pre-Flight

The exact-route contract fails because the toolbox package does not exist. Core, board, and knowledge siblings are staged and give cross-package references concrete targets.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/SKILL.md` (new) — dedicated sixteen-action router and explicit ownership boundary
- `skills/do-work-toolbox/actions/help.md` (new) — toolbox-only command menu and sibling discovery pointers
- `skills/do-work-toolbox/actions/install.md` (new) — companion-only installer for ui-design, bowser, last30days, and ideation-adhd
- `skills/do-work-toolbox/actions/validate-feedback.md` (new) — retained external-feedback adjudication with explicit core handoffs
- `skills/do-work-toolbox/actions/code-review.md` (new) — retained standalone code review
- `skills/do-work-toolbox/actions/ui-review.md` (new) — retained read-only UI review
- `skills/do-work-toolbox/actions/present-work.md` (new) — retained client-facing presentation workflow
- `skills/do-work-toolbox/actions/ai-report.md` (new) — retained screenshot-anchored report workflow
- `skills/do-work-toolbox/actions/deep-explore.md` (new) — retained multi-round concept exploration
- `skills/do-work-toolbox/actions/inspect.md` (new) — retained read-only working-tree inspection using core checks through sibling paths
- `_dev/tests/staged-skills-contract.sh` (modified) — exact toolbox inventory/routes, installer-boundary checks, and suite-wide runtime-reference resolution

The remaining toolbox actions (slop-check, quick-wins, scan-ideas, prime, note, stray-check, tidy-repo, and tutorial), two action references, nine guides, and seventeen required crew files are staged under the same package. Active legacy copies remain unchanged.

## Qualification

Passed — all 46 package files are non-empty; every retained action has exactly one route and a concrete action file; every runtime reference resolves locally or through a manifest-declared sibling; installer ownership is limited to the four retained companions; project-output documentation paths remain data rather than false runtime dependencies; and the route-to-action-to-reference flows are substantive and connected.

## Testing

**Tests run:** strict RED/GREEN staged-suite contract; full contract regressions; toolbox route-count and runtime-reference scans; warning-level ShellCheck and Bash syntax; Justfile parse; queue-kanban `go test ./...`, `go vet ./...`, and `go build ./...`; `git diff --check`
**Result:** ✓ All passing

**Red-green validation:**
- Toolbox package: ✗ required router/actions absent → ✓ all sixteen retained action routes exist exactly once
- Runtime references: ✗ copied actions resolved core files relative to toolbox → ✓ every core, board, and knowledge dependency names its sibling owner
- Installer ownership: ✗ monolithic installer included board, memory, and updater targets → ✓ only ui-design, bowser, last30days, and ideation-adhd remain

**New tests added:** Toolbox inventory, exact route counts, installer ownership, and generalized four-package runtime-reference assertions in `_dev/tests/staged-skills-contract.sh`.

*Verified by work action*

## Review

**Overall: 99%** | 2026-08-07T21:56:38Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 99% |
| Test Adequacy | 99% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None
**Minor findings:** 0
**Acceptance:** Pass — every retained optional capability is independently routed, dependency-complete, and separated from core, board, and knowledge ownership without changing the active legacy distribution.
**Suggested testing:** REQ-143 must exercise the toolbox package as part of exact four-module fresh-install, reinstall, validation-failure, and rollback fixtures.
**Follow-ups created:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Treating copied prose references as executable dependencies exposed the exact places where an action still assumed a monolithic skill root.
- Keeping command examples canonical to `do-work-toolbox` makes the new ownership boundary visible without deleting any legacy route before cutover.

**What didn't:**
- A verbatim copy left the companion installer owning board recipes, memory setup, and core updating; those workflows needed an explicit ownership prune.

**Worth knowing:**
- Some `docs/...` strings in tidy-repo are intended project destinations, not packaged runtime documents; the reference contract must distinguish outputs from dependencies.
- Toolbox actions still read core queue records, but they do so through explicit sibling references rather than duplicating lifecycle logic.

## Orientation

[MAP CHANGED] Optional reviews, reporting, exploration, repository hygiene, and companion installation now have an independently loadable staged context boundary at `skills/do-work-toolbox`. Core queue actions remain the authority for capture and execution; board and knowledge stay separate siblings.
