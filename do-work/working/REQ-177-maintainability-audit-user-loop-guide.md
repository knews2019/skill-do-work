---
id: REQ-177
title: Write the maintainability-audit user loop guide
status: claimed
created_at: 2026-08-13T22:35:10Z
claimed_at: 2026-08-14T10:21:48Z
user_request: UR-040
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
related: [REQ-176]
depends_on: [REQ-176]
batch: maintainability-audit
write_set: [skills/do-work-toolbox/docs/maintainability-audit-guide.md, skills/do-work-toolbox/docs/code-review-guide.md, skills/do-work-toolbox/actions/maintainability-audit.md]
maintenance: false
---

# Write the Maintainability-Audit User Loop Guide

## What

Relocate the draft spec's "Loop usage (for the operator, not the agent)" content into a user-facing guide at `skills/do-work-toolbox/docs/maintainability-audit-guide.md`, following the existing docs-guide pattern (e.g. `skills/do-work/docs/capture-guide.md`), and link it from the action's description blockquote as the other guides do.

## Also (folded from REQ-176's review — gate: user-visible)

`skills/do-work-toolbox/docs/code-review-guide.md:80` still shows `do-work-toolbox audit codebase` as a code-review invocation after REQ-176 moved that trigger to the new maintainability-audit action — a reader following the guide invokes a different action than the guide describes. Fix the stale invocation (point it at `code-review` / `review codebase`, and mention `audit codebase` now runs the maintainability audit). Provenance: Important finding in REQ-176's Review, restatement sweep; the retired-triggers fixture's stale attribution (`_dev/tests/fixtures/retired-core-moved-command-triggers.tsv:40`) stays untouched — pinned historical inventory, green as-is.

## Why (if provided)

An action file has no slot for content addressed to a different reader — a floor agent reads everything as instructions. The docs guide is the established home; the action's output footer carries only the immediate next-step commands.

## Context

The loop, using this suite's real vocabulary: run `do-work-toolbox maintainability-audit` → calibration conversation → read `do-work/audits/audit-<date>.md` → paste the Findings section into `do-work-toolbox validate-feedback` → capture Accepts through its handoff, park Discuss items with `do-work-toolbox note` → `do-work run` → re-run the audit (repeat runs ask only "reuse or recalibrate?"); the delta table must move; lock-in limits only ever tighten; classes you choose to live with go in `do-work/audits/waivers.md`. For a push-back citing a documented decision you no longer agree with, change the decision doc — not the code.

## Red-Green Proof
**RED prompt/case:** `ls skills/do-work-toolbox/docs/maintainability-audit-guide.md` fails; the loop narrative exists only inside UR-040's verbatim input.
**Why RED now:** REQ-176 deliberately excludes operator-addressed content from the action file, leaving the loop undocumented for users.
**GREEN when:** The guide exists with the loop above in the suite's vocabulary (lock-in limit, calibration gate, `do-work/audits/`), and the action's description blockquote links it.
**Validation:** Inferred during capture.

## Full Context

See `do-work/user-requests/UR-040/input.md` for complete verbatim input.

---
*Source: UR-040 — pasted maintainability-audit spec, validated via do-work-toolbox validate-feedback*

---

## Triage

**Route: A** - Simple

**Reasoning:** Two named files with fully specified changes: author one docs guide from content already written in this REQ's Context (following the existing docs-guide pattern), and fix one stale invocation line folded in from REQ-176's review. No discovery needed — the loop content, vocabulary, and file targets are all stated.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- **D-01 — write boundary extended to `skills/do-work-toolbox/actions/maintainability-audit.md` (one line).** The REQ's GREEN condition requires "the action's description blockquote links it", which the captured write_set omitted. Orchestrator decision at claim: the blockquote link is part of this REQ's deliverable (sibling precedent: capture.md's "User-facing walkthrough" link), so the file joins the boundary for exactly that edit. DECIDE & STATE — reversible, one line, required by the REQ's own text.
- **D-02 — no markdown links inside the guide.** All nine sibling toolbox guides use backtick code spans exclusively (zero markdown links in `skills/do-work-toolbox/docs/`); the new guide matches, so the both-topologies link-resolution requirement is satisfied trivially. The one required markdown link (action blockquote → guide) uses the exact sibling form from quick-wins.md/code-review.md. DECIDE & STATE.
- **D-03 — code-review-guide fix shape.** Removed the stale `do-work-toolbox audit codebase` line from the Usage block (kept `review codebase`) and added one note line after the block: "`do-work-toolbox audit codebase` invokes the maintainability audit (see `maintainability-audit-guide.md`), not this review." Smallest edit that both stops misdirecting readers and tells them where the phrase went, per the REQ's "Also" section. DECIDE & STATE.
