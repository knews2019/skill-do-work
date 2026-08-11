---
id: REQ-167
title: Deduplicate copy-pasted shell primitives across action files
status: pending
created_at: 2026-08-11T11:46:50Z
user_request: UR-036
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-165]
maintenance: true
related: [REQ-165, REQ-166, REQ-168]
batch: stabilization-audit
---

# Deduplicate Copy-Pasted Shell Primitives Across Action Files

## What

Sweep the shipped `skills/` tree for shell primitives that are restated in multiple action files (the CLAUDE.md trap-list primitives are the starting inventory: untracked-file enumeration, merge-commit-safe `git show`, root-anchored ignore patterns, curl download-and-rename, `git diff-tree` file listing, etc.) and give each one exactly one canonical home; other sites reference it instead of restating it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
