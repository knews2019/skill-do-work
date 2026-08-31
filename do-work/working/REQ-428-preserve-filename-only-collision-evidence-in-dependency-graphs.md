---
id: REQ-428
title: 'Review fix: Preserve filename-only collision evidence in dependency graphs'
status: claimed
domain: backend
created_at: 2026-08-30T19:21:33Z
claimed_at: 2026-08-31T13:55:10Z
route: A
user_request: UR-081
addendum_to: REQ-408
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
tdd: true
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
---

# Review Fix: Preserve Filename-Only Collision Evidence in Dependency Graphs

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What
Make repository collision evidence authoritative even when two filenames claim an ID that neither file uses in frontmatter. Dependency results must report that target as ambiguous rather than merely missing, while remaining safely blocked.

## Context
Found during review of REQ-408. Two files named for `REQ-021` with frontmatter IDs `REQ-030` and `REQ-031` produce collision evidence but no `REQ-021` node, so the graph's nil-target branch currently discards the more actionable ambiguity evidence.

Fold-first scan found no pending REQ or sweep in any UR that shares this dependency-evidence root cause.

## Requirements
- Check collision evidence before classifying an absent node as missing.
- Preserve deterministic unmet and ambiguity evidence for filename-only collisions.
- Keep the target blocked and unresolved for readiness/depth calculations.

## Red-Green Proof
**RED prompt/case:** Build a snapshot with `REQ-021` claimed by two filenames whose frontmatter IDs are `REQ-030` and `REQ-031`, then depend on `REQ-021`; assert the graph reports an ambiguous target rather than a missing target.
**Why RED now:** The node lookup returns nil before the graph consults `CollisionEntries`.
**GREEN when:** The named fixture is blocked with exact ambiguous-collision evidence and is never labeled merely missing.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---

## Triage

**Route: A** - Simple

**Reasoning:** The failure mode, production decision point, fixture shape, and expected evidence are all named. The change is a focused regression fix in the existing dependency graph.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
