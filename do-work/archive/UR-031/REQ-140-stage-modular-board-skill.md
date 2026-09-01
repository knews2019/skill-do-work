---
id: REQ-140
title: "Stage the Modular Board Skill"
status: completed
claimed_at: 2026-08-07T21:32:45Z
completed_at: 2026-08-07T21:38:54Z
commit: 5e9996f
route: C
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md, tools/prime-do-work-update.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-136, REQ-138]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-141, REQ-142, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
kb_status: promoted
kb_entry: REQ-140-stage-the-modular-board-skill.md
---

# Stage the Modular Board Skill

## What
Create a staged `skills/do-work-board` package that owns the board action, board documentation, Just template, and complete queue-kanban Go module.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add RED package/recipe assertions, copy the committed queue-kanban module and board materials into a dedicated package, author a board-only router/help and managed template, then validate exact source parity and all Go/browser/recipe behavior while leaving root active.
- [x] **[APPLY]:** Staged the committed board module, authored board-only routing/help, adapted cross-skill references, and added the managed client recipe template plus exact package contracts.
- [x] **[UNIFY]:** Reviewed all 50 board-package files; verified committed Go/web parity except documented package-path comments/prime links, non-empty assets, router/reference resolution, exact managed markers and install paths, shutdown/foreign-process/browser contracts, Testing empty-copy coverage, and warning-clean scripts.

## Why
The board is a distinct compiled application with its own UI and regression surface and should not tax or clutter the core router.

## Detailed Requirements
- Add a dedicated `SKILL.md`, board action/help, documentation, Just template, and all `tools/queue-kanban` source/tests/assets.
- Preserve live, static, summary, CLI, Testing-view, and completion-calendar behavior.
- Preserve port validation, foreign-process protection, bounded shutdown, and browser opening.
- Use the managed `do-work:recipes` section.
- Point board recipes at `.claude/skills/do-work-board/tools/queue-kanban`.
- Point `run-do-work-update` at the core updater.
- Do not remove the legacy board files yet.

## Constraints
- The board continues to read core queue data and the shared schema contract.
- The Go tool remains a single source after cutover; no permanent duplicate implementation.

## Dependencies
Requires REQ-136's suite contract and REQ-138's managed recipe mechanism.

## Builder Guidance
Certainty level: Firm. Preserve all validated queue-kanban fixes, including the Testing done-window empty-state behavior.

## Red-Green Proof
**RED prompt/case:** Load and run board commands from `skills/do-work-board` in an isolated staged suite.
**Why RED now:** Board routing, docs, recipes, and compiled source live inside the monolithic core root.
**GREEN when:** The staged board skill builds and passes all Go tests, its recipes resolve the new install path, and legacy active behavior remains green.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the board and installation boundaries.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*

## Triage

**Route: C** — the package owns a compiled Go/server/browser application and safety-sensitive shell recipes, so source parity and behavioral checks are both required.

## Scope

**Files I will touch:** `skills/do-work-board/**`, `_dev/tests/staged-skills-contract.sh`, and REQ/release metadata.

**Files I will not touch:** the active board action/module/Justfile, other staged packages, runtime queue data, or unrelated dirty REQ-134 work.

## Pre-Flight

The focused boundary test fails because all required board-package paths are absent. REQ-139 is archived; REQ-136/138 dependencies are satisfied.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/SKILL.md` (new) — independently loadable board-only router and suite ownership boundary
- `skills/do-work-board/actions/board.md` (new) — complete live/static/summary/CLI/Testing launcher with explicit core/toolbox references
- `skills/do-work-board/actions/help.md` (new) — board modes, toolchain requirement, and installed Just shortcuts
- `skills/do-work-board/docs/board-guide.md` (new) — board, calendar, Testing, drawer, filter, and overlap guidance
- `skills/do-work-board/justfile.template` (new) — exact managed recipe section using the board package and core updater install paths
- `skills/do-work-board/tools/queue-kanban/main.go` (new) — committed queue-kanban command surface copied into the package
- `skills/do-work-board/tools/queue-kanban/web/board.js` (new) — complete embedded board frontend including the Testing done-window empty-state fix
- `skills/do-work/actions/work.md` (modified) — modular next-version call explicitly names core's version file instead of relying on the board tool's repo-root default
- `_dev/tests/staged-skills-contract.sh` (modified) — full board-package inventory, recipe safety/path, router, cross-skill, and runtime-reference assertions

All other committed queue-kanban Go sources, tests, embedded assets, module files, and the subsystem prime were copied into the staged package. Active legacy board sources and recipes remain in place.

## Qualification

Passed — 50 package files are present and non-empty; all seven requirements trace to the router/action/template/module or contract; P-A-U is complete; the compiled module, embedded UI, managed recipe paths, core updater call, and cross-skill schema/version data paths are real and exercised.

## Testing

**Tests run:** focused RED/GREEN staged-suite contract; `go test ./...`, `go vet ./...`, and `go build ./...` from the staged module; warning-level ShellCheck; Bash syntax; Just template parse/list; full contract regressions; committed-source parity diff; Testing-filter and listener/browser source probes; `git diff --check`
**Result:** ✓ All passing

**Red-green validation:**
- Board package: ✗ ten required board entry points absent → ✓ dedicated router/action/help/guide/template and complete tracked module present
- Recipe paths: ✗ active recipes build root `tools/queue-kanban` → ✓ four managed recipes build `.claude/skills/do-work-board/tools/queue-kanban` and update through core
- Validated behavior: ✗ no isolated board module → ✓ staged Go tests cover live/static/summary/CLI, Testing, calendar, browser, parser, and server behavior

**New tests added:** Board inventory, recipe path/safety, router, sibling-reference, and core version-file assertions in `_dev/tests/staged-skills-contract.sh`.

*Verified by work action*

## Review

**Overall: 99%** | 2026-08-07T21:37:55Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 99% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None
**Minor findings:** 0
**Acceptance:** Pass — the staged board preserves every validated Go/UI mode and shell safety contract, while its recipes and cross-skill references resolve to the modular layout.
**Suggested testing:** REQ-143 should execute this exact template inside its hermetic fresh-install and reinstall projects rather than maintaining a second recipe body.
**Follow-ups created:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Exporting the committed Git tree instead of copying the dirty working directory kept unrelated REQ-134 changes out of this package while preserving every accepted board fix.
- Treating the Just template as a board-owned artifact made listener shutdown, foreign-process refusal, browser opening, and install paths testable together.

**What didn't:**
- A direct copy left core-schema and version references looking local to the board package; the runtime-reference contract exposed them and the core next-version seam needed an explicit `--version-file`.

**Worth knowing:**
- The active board module currently has unrelated uncommitted REQ-134 changes. Before cutover, REQ-144 must synchronize any board changes committed after this staging snapshot into the package.
- The queue tool supports `--version-file`; modular core must always pass it because the board package's repo-root default is only correct in the legacy source checkout.

## Orientation

[MAP CHANGED] Queue visualization is now staged as its own compiled application under `skills/do-work-board`, including its launcher, guide, managed Just template, Go server/CLI, and embedded UI. Core supplies queue semantics and version ownership; the board package supplies every visual/testing mode without enlarging core context.
