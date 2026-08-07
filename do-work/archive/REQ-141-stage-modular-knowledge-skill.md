---
id: REQ-141
title: "Stage the Modular Knowledge Skill"
status: completed
claimed_at: 2026-08-07T21:41:06Z
completed_at: 2026-08-07T21:46:13Z
route: C
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: []
tdd: true
suggested_spec: refactor
depends_on: [REQ-136]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-142, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
kb_status: pending
kb_entry:
---

# Stage the Modular Knowledge Skill

## What
Create a staged `skills/do-work-knowledge` package for BKB, dream, memory, interview, prompts, knowledge assets, and knowledge hooks.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add RED inventory/privacy/hook-migration assertions; copy the cohesive knowledge actions, references, assets, guides, hooks, and required guardrails; author a knowledge-only router/help and explicit memory setup action; adapt core/toolbox boundaries without enabling hooks by default.
- [x] **[APPLY]:** Staged the cohesive knowledge runtime and assets, authored its router/help and explicit memory setup owner, rewired cross-skill references, and moved the optional hook fragment/scripts to deterministic knowledge paths without default activation.
- [x] **[UNIFY]:** Reviewed all 44 package files; verified complete prompt/interview inventory, non-empty actions/references/guides, route and runtime-reference resolution, privacy and raw-store rules, exact hook migrations, memory-free core defaults, JSON validity, hook code parity, Bash syntax, and warning-clean shell.

## Why
Knowledge management and session memory are self-contained concerns that accumulated far outside the core task-queue mission.

## Detailed Requirements
- Add a dedicated `SKILL.md` and move-by-copy BKB, dream, memory, interview, prompts, references, templates, docs, hooks, and interviewer guidance.
- Make this package own memory-module setup and hook configuration.
- Keep memory capture disabled on a fresh full-suite install.
- Allow explicit setup to enable memory hooks.
- Define deterministic migration targets from `.claude/skills/do-work/hooks/memory-*.sh` to `.claude/skills/do-work-knowledge/hooks/`.
- Preserve memory privacy, machine-local raw stores, bootstrap semantics, KB behavior, and optional core lessons handoff.
- Do not delete or deactivate legacy copies yet.

## Constraints
- Existing knowledge behavior must remain unchanged through staging.
- Do not enable memory transcript capture merely because the suite is installed.

## Dependencies
Requires REQ-136's suite contract.

## Builder Guidance
Certainty level: Firm. The package is broad but cohesive around retained knowledge and memory functionality.

## Red-Green Proof
**RED prompt/case:** Load BKB, memory, dream, interview, and prompt actions from an isolated `skills/do-work-knowledge` package.
**Why RED now:** These actions and their hooks/templates are coupled to the monolithic skill root.
**GREEN when:** Every knowledge route resolves within the staged package, existing behavioral tests pass, and fresh-install fixtures prove memory hooks remain disabled.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for knowledge ownership and hook policy.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*

## Triage

**Route: C** — broad package staging with persistent-memory privacy, hook composition, prompt trust, and interview state contracts.

## Scope

**Files I will touch:** `skills/do-work-knowledge/**`, `_dev/tests/staged-skills-contract.sh`, and REQ/release metadata.

**Files I will not touch:** active legacy knowledge files/hooks, core default hooks, client memory/KB data, other staged packages, or unrelated dirty REQ-134 work.

## Pre-Flight

The focused boundary test fails on all required knowledge-package paths. Core default hooks are currently memory-free, satisfying the pre-change privacy baseline.

## Implementation Summary

**Files changed:**
- `skills/do-work-knowledge/SKILL.md` (new) — independent BKB/memory/dream/interview/prompts router and privacy boundary
- `skills/do-work-knowledge/actions/help.md` (new) — package command menu and explicit fresh-install hook policy
- `skills/do-work-knowledge/actions/setup-memory.md` (new) — user-invoked scaffolding, local-ignore safety, exact legacy migration, hook composition, rollback, and verification
- `skills/do-work-knowledge/actions/memory.md` (new) — complete hookless memory command surface with core capture/schema boundaries
- `skills/do-work-knowledge/actions/memory-reference.md` (new) — store schemas, cap/consolidation, recall, ledger, redaction, and hook internals
- `skills/do-work-knowledge/actions/bkb.md` (new) — full knowledge-base lifecycle with optional core lessons inbox flow
- `skills/do-work-knowledge/actions/interview.md` (new) — structured interview state and export/ingest lifecycle
- `skills/do-work-knowledge/actions/prompts.md` (new) — trusted shipped-library resolver and execution gate
- `skills/do-work-knowledge/hooks/memory-hooks.json` (new) — optional fragment targeting only `do-work-knowledge`, never a default hook set
- `_dev/tests/staged-skills-contract.sh` (modified) — knowledge inventory, default-disable, modular hook target, and deterministic migration assertions

The package also contains dream, all companion references/auditor, four guides, six required guardrails, every shipped prompt, the interview template, and both memory hook scripts. Legacy root knowledge files and hooks remain active and unchanged.

## Qualification

Passed — 44 package files are non-empty; every detailed requirement traces to an action, asset, hook, or boundary assertion; P-A-U is complete; core lessons-to-inbox, prompt trust, interview state, memory store, and optional hook paths remain substantive and connected.

## Testing

**Tests run:** focused RED/GREEN staged-suite contract; full contract regressions; complete tracked prompt/interview inventory; JSON parse; staged hook Bash syntax and warning-level ShellCheck; executable hook-code parity with legacy sources; runtime-reference resolution; `git diff --check`
**Result:** ✓ All passing

**Red-green validation:**
- Package routes/assets: ✗ eighteen required knowledge entry points absent → ✓ BKB, memory, dream, interview, prompt, hook, and guidance paths resolve
- Fresh privacy: ✓ core default hooks contain no memory capture before and after staging; optional fragment is inert until explicit setup
- Hook migration: ✗ fragment targeted legacy core paths → ✓ exact old-to-new pairs and both modular commands are contract-pinned

**New tests added:** Knowledge contents, asset inventory, default-disable, hook target, and migration assertions in `_dev/tests/staged-skills-contract.sh`.

*Verified by work action*

## Review

**Overall: 99%** | 2026-08-07T21:45:20Z

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
**Acceptance:** Pass — knowledge behavior and privacy contracts are preserved in an independent package, with memory capture opt-in and deterministic legacy hook migration.
**Suggested testing:** REQ-143 must exercise enabled, declined, no-JSON-tool, partial-hook, legacy-migration, invalid-settings, and rollback states against this fragment.
**Follow-ups created:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Keeping memory and BKB together preserves their shared ledger/value audit while removing both from core context.
- Separating the optional fragment (`memory-hooks.json`) from core's default `hooks.json` makes “installed but disabled” mechanically testable.

**What didn't:**
- Copying the actions verbatim retained an obsolete rationale that memory belonged inside the monolith and left setup references pointing at the all-purpose installer; both had to become explicit package ownership.

**Worth knowing:**
- Exact legacy hook strings are data-migration keys. Broader path substitution risks rewriting unrelated commands in client settings.
- The raw store remains plaintext even though it is machine-local; redaction-before-truncation and never-track checks are independent defenses, not substitutes.

## Orientation

[MAP CHANGED] Retained knowledge now has its own staged context boundary at `skills/do-work-knowledge`. BKB synthesis, lightweight memory, consolidation, interviews, prompts, and their privacy-sensitive optional hooks live together; core can hand off consented lessons without loading or enabling the knowledge engine.
