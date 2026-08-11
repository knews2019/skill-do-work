---
id: REQ-167
title: Deduplicate copy-pasted shell primitives across action files
status: completed
created_at: 2026-08-11T11:46:50Z
user_request: UR-036
domain: general
prime_files: []
write_set: [decisions/audits/2026-08-11-prescribed-shell-primitives.md, skills/do-work/docs/prescribed-shell-primitives.md, _dev/tests/prescribed-shell-canonicalization.sh, _dev/tests/contract-regressions.sh, skills/do-work/actions/commit.md, skills/do-work/actions/review-work.md, skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/crew-members/background-agents.md, skills/do-work-toolbox/actions/inspect.md, skills/do-work-toolbox/actions/ai-report.md, skills/do-work-toolbox/actions/present-work.md, skills/do-work-toolbox/actions/install.md, skills/do-work-toolbox/actions/stray-check.md, skills/do-work-toolbox/crew-members/background-agents.md, skills/do-work-knowledge/actions/memory-reference.md, skills/do-work-knowledge/crew-members/background-agents.md, skills/do-work-board/actions/board.md]
tdd: false
suggested_spec:
depends_on: [REQ-165]
maintenance: true
related: [REQ-165, REQ-166, REQ-168]
batch: stabilization-audit
claimed_at: 2026-08-11T12:30:27Z
route: C
completed_at: 2026-08-11T12:43:36Z
commit:
kb_status: pending
kb_entry:
---

# Deduplicate Copy-Pasted Shell Primitives Across Action Files

## What

Sweep the shipped `skills/` tree for shell primitives that are restated in multiple action files (the CLAUDE.md trap-list primitives are the starting inventory: untracked-file enumeration, merge-commit-safe `git show`, root-anchored ignore patterns, curl download-and-rename, `git diff-tree` file listing, etc.) and give each one exactly one canonical home; other sites reference it instead of restating it.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Inventory the documented trap primitives first; create one core shipped canonical guide; replace only duplicated rationale with short site intent plus cross-package-safe pointers while retaining executable commands/fallbacks; add a regression probe that rejects known stale restatement phrases outside the guide; verify every rewritten shell fence and full contracts.
- [x] **[APPLY]:** Added the canonical guide and durable inventory, replaced former rationale copies with package-safe pointers while retaining their commands and local gates, and added the canonicalization ratchet through the aggregate contract seam.
- [x] **[UNIFY]:** Reviewed the complete diff/stat and all 18 scoped files; verified the guide's commands and anchors, every former site's retained caller-specific instructions, the audit/test correspondence, executable bit, absence of stale phrases/debug artifacts, shell-fence lint, Bash syntax, and whitespace.

## Why (if provided)

One bug in a copy-pasted primitive becomes N review findings — the documented incident: the untracked-files trap "had been copy-pasted into four action files; the audit only flagged one." Dedup converts N future findings into 1.

## Context

- `maintenance: true` — this is a deliberate narrowing pass on the skill's own operating instructions (delete-before-you-add applies).
- Canonical-home options, builder's choice per primitive: (a) one action/reference file owns the full block and others cross-reference it by local or explicit sibling path per CLAUDE.md's cross-reference convention; (b) a small shipped helper script the prose invokes. Respect agent compatibility: each action file must still work as a standalone prompt, so a reference must carry enough inline context to act on (one-line intent + pointer), not a bare link.
- Package boundaries matter: core cannot depend on board/knowledge/toolbox files. A primitive used across packages needs its home in core or a per-package copy with a single-source note — flag these cases rather than silently duplicating.
- Depends on REQ-165: the lint harness validates rewritten blocks and guards the survivors.

## Detailed Requirements

- Build the inventory first: grep each trap-list primitive across `skills/` and record every restatement site (the audit artifact).
- For each primitive: pick the canonical home, rewrite other sites as references, and verify no site kept a stale variant of the block.
- Do not change primitive semantics in this pass — behavior changes are separate REQs; this is consolidation.
- Where two sites had *diverged* copies (one fixed, one stale), the fixed variant wins; note the divergence in the audit artifact.

## Builder Guidance

Certainty: Firm on the goal, exploratory on per-primitive home choice. Scope cue: consolidation only — resist improving the primitives while moving them.

## Red-Green Proof

**RED prompt/case:** `grep -rn 'untracked-files=all\|ls-files --others' skills/*/actions/` (and equivalents for the other trap primitives) returns multiple files each restating the full block.
**Why RED now:** The traps were fixed where found, but the copies were never consolidated — the CLAUDE.md rule "grep the same primitive across all actions before calling it fixed" exists because divergence already happened.
**GREEN when:** Each inventoried primitive has one canonical statement in the shipped tree; every other former restatement site is a reference; the REQ's audit artifact lists primitive → home → former sites.
**Validation:** Inferred during capture (plan discussed and endorsed in-session).

## Full Context

See `do-work/user-requests/UR-036/input.md` for complete verbatim input.

---
*Source: "do-work capture-request for audit and fix to simplify and make it robust" (UR-036)*

---

## Triage

**Route: C** - Complex

**Reasoning:** This is a cross-package instruction-maintenance sweep requiring a durable inventory, per-primitive ownership decisions, many standalone-prompt rewrites, package-safe references, and a ratchet against future rationale copies.

**Planning:** Required

## Plan

1. Record a grep-based audit under `decisions/audits/` mapping each CLAUDE.md trap primitive to its shipped canonical home, duplicated rationale sites, execution-only uses, and any divergent variants; do not use perishable line-number coordinates.
2. Add `skills/do-work/docs/prescribed-shell-primitives.md` as the single shipped rationale/fallback home. Preserve exact semantics for per-file untracked inventory, merge-aware diffs and commit file lists, local excludes, atomic downloads, raw-text sanitization, output filtering, and cross-block state.
3. Replace former rationale copies with one-line local intent plus explicit pointers to the guide, using same-package paths in core and `../../do-work/docs/...` from sibling packages. Keep every caller's executable command, caller-specific gate, and documented fallback behavior.
4. Add `_dev/tests/prescribed-shell-canonicalization.sh`, wire it into the aggregate contracts, and verify canonical headings/pointers, absence of the old high-risk rationale phrases outside the guide, shell-fence lint, reference contracts, and the full suite.

**Plan validation:** Every Detailed Requirement maps to an ordered task, and every planned rewrite traces to a restatement found by the inventory. Four tasks keep the work grouped by audit, canonical source, subtractive rewrites, and ratchet verification; no primitive behavior change or helper-script extraction is planned.

*Generated by Plan phase after repository inspection*

## Exploration

- REQ-121 already extracted the shared uncommitted inventory into `tools/checks/uncommitted-inventory.sh` and `associate-files.sh`; the remaining duplication is the long manual-fallback rationale in `actions/commit.md` and toolbox `actions/inspect.md`, not the executable mechanism.
- Merge-aware commit reading is restated in core `review-work.md`, toolbox `ai-report.md`, and three places in toolbox `present-work.md`. The command stays at each execution/template site; the explanation of merge detection, first-parent diffing, and quoted `^2` belongs once.
- The local-ignore section is copied at length into the core, knowledge, and toolbox `background-agents.md` crew files with only package paths changed. Board and install actions also restate parts of its root-anchoring/worktree rationale. A core guide is reachable from every sibling without reversing the forbidden core dependency direction.
- Atomic-download commands in shipped scripts and install tables are real uses, not prose duplication to delete. The duplicated explanation/red-flag language can point to one invariant while the commands remain self-contained.
- `git diff-tree` file listing, raw query sanitization, `diff -x` output filtering, and cross-block state have fewer full copies; the audit will distinguish their canonical statement from necessary caller-specific commands so grep count is not mistaken for duplication.
- Existing regression style favors standalone `_dev/tests/*.sh` probes explicitly invoked by `contract-regressions.sh`; REQ-165 now additionally parses every surviving shell fence.

*Generated by Exploration phase*

## Scope

**Files I will touch:**
- `decisions/audits/2026-08-11-prescribed-shell-primitives.md` (new) — durable primitive/home/former-site inventory
- `skills/do-work/docs/prescribed-shell-primitives.md` (new) — sole shipped rationale and fallback contract
- `_dev/tests/prescribed-shell-canonicalization.sh` (new) — canonical-home/restatement regression probe
- `_dev/tests/contract-regressions.sh` (modified) — invoke the new probe
- `skills/do-work/actions/commit.md` (modified) — replace duplicated inventory/state rationale with core-guide pointers
- `skills/do-work/actions/review-work.md` (modified) — reference the canonical merge-diff rule
- `skills/do-work/actions/work.md` (modified) — reference the canonical cross-block state rule
- `skills/do-work/actions/work-reference.md` (modified) — reference the canonical cross-block state rule
- `skills/do-work/crew-members/background-agents.md` (modified) — replace the long local-ignore copy with a core-guide pointer
- `skills/do-work-toolbox/actions/inspect.md` (modified) — replace duplicated inventory/state rationale with core-guide pointers
- `skills/do-work-toolbox/actions/ai-report.md` (modified) — reference canonical file-listing and merge-diff rules
- `skills/do-work-toolbox/actions/present-work.md` (modified) — replace three merge-diff rationale copies with pointers
- `skills/do-work-toolbox/actions/install.md` (modified) — reference canonical atomic-download and local-ignore invariants
- `skills/do-work-toolbox/actions/stray-check.md` (modified) — reference canonical per-file untracked enumeration rationale
- `skills/do-work-toolbox/crew-members/background-agents.md` (modified) — replace the long local-ignore copy with the core guide
- `skills/do-work-knowledge/actions/memory-reference.md` (modified) — reference canonical raw-text sanitization
- `skills/do-work-knowledge/crew-members/background-agents.md` (modified) — replace the long local-ignore copy with the core guide
- `skills/do-work-board/actions/board.md` (modified) — reference canonical local-ignore rationale after the retained command

**Files I will NOT touch:** shipped helper-script behavior, queue-kanban code/templates, action routing/help, package manifests, or primitive command semantics

**Acceptance criteria (restated from REQ):**
- [x] The durable audit maps every starting primitive to one home, every former rationale site, and divergence/execution-only disposition.
- [x] Exactly one shipped guide owns each primitive's full rationale and fallback constraints.
- [x] Former copies carry enough site-specific intent to act plus an explicit package-safe guide pointer, with executable commands preserved.
- [x] No primitive semantics change; fixed variants remain the canonical wording when prior copies diverged.
- [x] The regression probe and REQ-165 shell harness pass, and no stale high-risk rationale phrase remains outside the guide.

## Decisions

- **D-01 — Canonicalize rationale in a core document, not commands in a new helper.** The affected actions still need independently runnable commands, while a core guide is reachable from every sibling package without introducing a forbidden core-to-sibling dependency.
- **D-02 — Treat execution count separately from duplication count.** Required command blocks remain at their callers; only shared failure-mode explanations and fallback rationale move to the canonical home.
- **D-03 — Preserve the fixed variant when copies diverged.** The guide owns `-uall`/NUL-aware inventory, quoted merge-parent detection with first-parent show, worktree-safe Git excludes, and temp-download publication with failure-preserving cleanup.

## Implementation Summary

**Files changed:**
- `decisions/audits/2026-08-11-prescribed-shell-primitives.md` (new)
- `skills/do-work/docs/prescribed-shell-primitives.md` (new)
- `_dev/tests/prescribed-shell-canonicalization.sh` (new)
- `_dev/tests/contract-regressions.sh` (modified)
- `skills/do-work/actions/commit.md` (modified)
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/crew-members/background-agents.md` (modified)
- `skills/do-work-toolbox/actions/inspect.md` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/actions/present-work.md` (modified)
- `skills/do-work-toolbox/actions/install.md` (modified)
- `skills/do-work-toolbox/actions/stray-check.md` (modified)
- `skills/do-work-toolbox/crew-members/background-agents.md` (modified)
- `skills/do-work-knowledge/actions/memory-reference.md` (modified)
- `skills/do-work-knowledge/crew-members/background-agents.md` (modified)
- `skills/do-work-board/actions/board.md` (modified)

**What was done:** Established a single shipped core guide for eight prescribed shell primitives and a durable primitive-to-home/former-site audit. Replaced duplicated cross-package rationale with explicit pointers while preserving executable commands and caller-specific policy, then added a regression probe that enforces the guide headings, all former-site pointers, and absence of known stale rationale phrases.

## Qualification

Passed — the mechanical qualifier verified the 18-file scope and requirement trace, and scope-drift found exact agreement between the declaration and Implementation Summary. The qualifier's one judgment warning is expected: the new `decisions/audits/` document is itself the requested durable audit artifact, not an executable/documentation dependency that shipped code should import. Manual review confirmed a substantive net-subtractive consolidation, package-safe data flow from caller pointers to the core guide, and retained command semantics at every execution site.

## Testing

**Tests run:**
- `_dev/tests/prescribed-shell-canonicalization.sh`
- `_dev/tests/action-shell-blocks.sh`
- `bash -n _dev/tests/prescribed-shell-canonicalization.sh _dev/tests/contract-regressions.sh`
- `shellcheck --severity=warning _dev/tests/prescribed-shell-canonicalization.sh _dev/tests/contract-regressions.sh`
- `skills/do-work/tools/checks/qualify.sh do-work/working/REQ-167-dedupe-prescribed-shell-primitives.md`
- `skills/do-work/tools/checks/scope-drift.sh do-work/working/REQ-167-dedupe-prescribed-shell-primitives.md`
- `bash _dev/tests/contract-regressions.sh`
- `git diff --check`
- stale-rationale phrase sweep across shipped Markdown

**Result:** ✓ All passing; the aggregate contract suite exited 0, shipped package references passed, all 59 fenced shell blocks plus 15 shipped shell files passed the REQ-165 harness, and the stale-phrase sweep returned zero matches outside changelogs and the canonical guide.

**Red-green validation:**
- Before rewrites, the new probe attributed missing canonical pointers at every former site and found the captured stale rationale phrases.
- After consolidation, the same probe passes with all required guide headings/pointers and zero forbidden restatements.

**New tests added:**
- `_dev/tests/prescribed-shell-canonicalization.sh` — enforces canonical homes, former-site pointers, and a high-risk stale-phrase ratchet

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/contract-regressions.sh` — invokes the new canonicalization probe; no existing assertion was weakened

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-11T12:43:01Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — the audit covers all eight starting primitives, the guide is the sole full-rationale home, callers retain actionable commands/gates, and both focused and aggregate regressions pass.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Separating “command executions” from “rationale copies” kept standalone actions runnable while still removing the maintenance multiplier that caused prior partial fixes.
- A phrase ratchet plus required pointers catches both obvious copy-back and silent removal of the canonical trail; the shell-fence harness independently protects the commands that remain.

**What didn't:**
- Grep count alone over-reported duplication because legitimate caller commands and shared explanation use the same tokens. The inventory needed a disposition column before deletion decisions were safe.

**Worth knowing:**
- Cross-package shell policy belongs in core because every sibling can depend on core, while moving a shared primitive into toolbox/board/knowledge would reverse the allowed dependency direction.
- Durable audit artifacts are evidence entry points, not runtime dependencies; a no-static-reference qualifier warning is expected for this requested document class.

## Orientation

[MAP CHANGED] Prescribed shell safety rationale now lives in `skills/do-work/docs/prescribed-shell-primitives.md`; caller actions keep only their command, local intent, and package-safe pointer, with `_dev/tests/prescribed-shell-canonicalization.sh` preventing drift.
