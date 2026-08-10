---
id: REQ-139
title: "Stage the Modular Core Skill"
status: completed
claimed_at: 2026-08-07T21:20:27Z
completed_at: 2026-08-07T21:31:33Z
commit: 9ba534e
route: C
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [tools/prime-do-work-update.md, tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-137, REQ-138]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-140, REQ-141, REQ-142, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
kb_status: pending
kb_entry:
---

# Stage the Modular Core Skill

## What
Create a self-contained staged `skills/do-work` package while the repository-root all-in-one distribution remains active.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Stage a reduced core router plus copied orchestration actions and their guardrail/spec/check dependencies; add an isolated package validator that resolves relative runtime references and bans maintainer-doc citations, while leaving every root source active.
- [x] **[APPLY]:** Staged the curated core runtime, authored its reduced router/help boundary, rewired board/knowledge/toolbox calls to declared siblings, and added the package contract without changing the active root router or actions.
- [x] **[UNIFY]:** Reviewed the full 69-file staged tree plus both test changes; verified all files are non-empty, route/hook/runtime references resolve, sibling boundaries are manifest-declared, scripts are warning-clean, the active root board route remains, and no debug artifacts were added.

## Why
Core needs an independently discoverable context boundary without weakening the orchestration behavior the user relies on.

## Detailed Requirements
- Stage capture, run, verify-requests, review-work, clarify, abandon, cleanup, commit, roadmap, forensics, version/update/recap, help, required references, queue schema, crew guardrails, specs, core hooks, checks, and updater.
- Preserve the current feature-rich work orchestrator.
- Keep `actions/kb-lessons-handoff.md` in core because `work` and `review-work` invoke its consent, status, and archival flow; resolve its knowledge-storage and follow-on command references through sibling `do-work-knowledge`.
- Keep pipeline temporarily for behavioral parity; REQ-145 removes it after cutover.
- Resolve eventual queue-kanban calls from sibling `do-work-board` and knowledge handoffs from sibling `do-work-knowledge`.
- Add package-boundary tests proving every staged runtime reference resolves and no shipped file cites repository-root maintainer instructions.
- Do not activate this package or alter the active root router yet.

## Constraints
- This is a move-by-copy staging step. Legacy active files remain until cutover.
- Avoid permanent duplicated sources; REQ-144 removes the old active copies.

## Dependencies
Requires the bridge updater and managed-section foundation in REQ-137 and REQ-138.

## Builder Guidance
Certainty level: Firm. Preserve behavior first; simplification is limited to functional ownership boundaries.

## Red-Green Proof
**RED prompt/case:** Point an isolated skill loader at `skills/do-work/SKILL.md` and validate all runtime references in a staged four-skill fixture.
**Why RED now:** No independently loadable core package exists and current references assume one monolithic skill root.
**GREEN when:** Core loads from its staged package, its required references resolve across the staged suite, and the active legacy installation remains unchanged.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the core behavior and module-boundary decisions.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*

---

## Triage

**Route: C** - Complex

**Reasoning:** This establishes the first runtime package boundary and must preserve orchestration behavior while converting board/knowledge calls into sibling-skill boundaries.

## Plan

1. Add a failing staged-package contract for an independently loadable `skills/do-work` router, required core artifacts, reference resolution, and maintainer-doc exclusions.
2. Copy the active core actions, references, crew guardrails, specs, hooks, checks, updater, and supporting docs; author a core-only router and explicit sibling handoffs.
3. Validate package/reference integrity, unchanged root runtime, shell/contract/Go surfaces, and exact updater/managed-section parity.

## Scope

**Files I will touch:** `skills/do-work/**`, `_dev/tests/staged-skills-contract.sh`, `_dev/tests/contract-regressions.sh`, and the REQ/release metadata.

**Files I will NOT touch:** Active root router/actions, staged sibling implementations, runtime data, or unrelated REQ-134 files.

**Acceptance criteria:** Core router loads independently; required actions/guards/specs/hooks/checks/updater exist; internal runtime references resolve or name an explicit sibling; no staged file cites root maintainer instructions; root runtime stays unchanged.

## Pre-Flight

REQ-137/138 are archived. Existing unrelated dirty work is preserved and the shared contract hunk will be staged selectively.

## Implementation Summary

**Files changed:**
- `skills/do-work/SKILL.md` (new) — independently loadable core-only router with explicit four-skill ownership and temporary pipeline compatibility
- `skills/do-work/actions/help.md` (new) — compact core menu plus named extension entry points
- `skills/do-work/actions/work.md` (new) — feature-complete work orchestrator copy with board and toolbox sibling references
- `skills/do-work/actions/work-reference.md` (new) — complete queue schema, execution contracts, and explicit board/knowledge/toolbox boundaries
- `skills/do-work/actions/kb-lessons-handoff.md` (new) — consent-preserving core handoff into the sibling knowledge skill
- `skills/do-work/actions/pipeline.md` (new) — temporary stateful orchestration with toolbox-owned investigate/present steps
- `skills/do-work/tools/do-work-update.sh` (new) — suite-aware updater staged inside the core package
- `skills/do-work/hooks/hooks.json` (new) — core SessionStart and temporary pipeline Stop hooks
- `_dev/tests/staged-skills-contract.sh` (new) — package contents, routes, hooks, sibling declarations, runtime references, and maintainer-citation contract
- `_dev/tests/contract-regressions.sh` (modified) — invokes the staged-suite boundary contract from the repository acceptance suite

The remaining staged files are curated copies of the required core actions/references, all crew guardrails, all implementation specs, core hooks, shell checks, updater helpers, guides, release history, version, and next-step guidance. The legacy root runtime remains active and unchanged.

## Qualification

Passed — 69 staged runtime files are non-empty; ten implementation entry points were verified against the diff and Scope; all seven detailed requirements trace to package files or the boundary contract; P-A-U is complete; routers, hooks, board tooling, knowledge handoff, and toolbox compatibility paths flow through real or manifest-declared destinations.

## Testing

**Tests run:** focused RED/GREEN staged-skill contract; warning-level ShellCheck; Bash syntax; `/bin/bash _dev/tests/contract-regressions.sh`; every standalone `_dev/tests/*.sh`; `go test ./...`; `go vet ./...`; `go build ./...`; `git diff --check`
**Result:** ✓ All passing

**Red-green validation:**
- Package load: ✗ all 32 required staged-core paths missing → ✓ independently loadable `skills/do-work/SKILL.md` and complete required runtime
- Runtime references: ✗ monolithic local board/knowledge/toolbox assumptions → ✓ local paths resolve and extension paths name manifest-declared siblings
- Activation safety: ✗ no parallel staged package to validate → ✓ staged package passes while the active root router retains its board route and action

**New tests added:** `_dev/tests/staged-skills-contract.sh`, invoked by the main contract suite.

*Verified by work action*

## Review

**Overall: 99%** | 2026-08-07T21:29:57Z

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
**Acceptance:** Pass — the core package loads independently, preserves the complete queue orchestrator and schemas, resolves local and sibling runtime paths, and leaves the active monolithic runtime in place.
**Suggested testing:** REQ-140 through REQ-142 should extend the same boundary probe so declared sibling references become exact on-disk reference checks as each package appears.
**Follow-ups created:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Copying the feature-rich actions first and changing only their ownership-boundary references preserved behavior while making the split auditable.
- A staged-suite contract can validate not-yet-created siblings safely by requiring their names in the exact suite manifest, then automatically tighten to file existence when each sibling directory appears.

**What didn't:**
- A naive ban on every textual `CLAUDE.md` mention incorrectly treated historical changelog text and consumer-project discovery guidance as live citations; the contract needed to target actual links and directive citations.

**Worth knowing:**
- The temporary pipeline necessarily crosses into toolbox for inspect and present; REQ-145 removes that stateful compatibility path only after activation.
- Core still carries all domain crew files because `work` selects them dynamically; moving those guardrails would weaken the orchestrator even though several toolbox actions also consume them.

## Orientation

[MAP CHANGED] The modular suite now has a staged core context boundary at `skills/do-work`: queue orchestration remains feature-complete, while board, knowledge, and toolbox capabilities are named sibling dependencies instead of hidden local files. The repository-root all-in-one skill remains the live distribution until the migration gate.
