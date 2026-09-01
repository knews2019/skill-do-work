---
title: "Lessons from REQ-168: Delete-or-test audit of defensive code in shipped skills"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-168-delete-or-test-audit-of-defensive-code-i.md]
related:
  - page: concept-contract-verification-gates
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-168: Delete-or-test audit of defensive code in shipped skills

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

Audit every defensive layer in the shipped `skills/` tree — fallbacks, guards, workarounds, retry/recovery blocks, and warning apparatus in both shell (hooks, prescribed blocks) and prose (Rules/Rationalizations sections that restate hygiene) — and disposition each one: **keep** (traces to a named incident AND is covered by a test), or **delete** (can't name the incident it prevents, or its cost now exceeds the surface it protects).

## Solution summary

Audited the shipped executable and explicit prose defensive surface, tying retained layers to incidents and test evidence while classifying behavior-changing candidates separately. Removed four complete decorative Rationalizations/Red Flags pairs, two duplicate Rationalizations sections, three generic inspect warnings, and one arbitrary commit-size warning—96 shipped action lines—with all action steps, output contracts, permission boundaries, and incident-specific warnings preserved. Added a focused probe that keeps shell/prose audit coverage complete and prevents the safe deletions from returning.

## What worked

- Auditing cohesive mechanisms instead of counting branches made expensive-but-earned transactions distinguishable from cheap decorative prose.
- Starting from the existing steps and verification checklists made safe deletion objective: the removed tables carried no unique permission, output, or execution contract.

## What didn't work

- Keyword counts dramatically overstate defensive surface: changelogs, test descriptions, ordinary domain language, and one atomic rollback transaction all contain many “guard/failure” tokens without representing separate layers.

## Worth knowing

- A numeric commit-size warning is weaker than semantic atomicity; it can split a correct large change and allow an unrelated small one.
- Static reference coverage is valid evidence for prompt-only permission/schema sections, while runtime fallbacks need syntax reachability plus targeted behavioral or contract fixtures.

## Back-reference

See `do-work/archive/UR-036/REQ-168-defensive-code-delete-or-test-audit.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8703b66`.
