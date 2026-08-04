---
id: REQ-101
title: Docs + ADR — multi-checkout guide and the session-ownership decision record
status: pending
created_at: 2026-08-04T19:44:17Z
user_request: UR-018
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-096, REQ-097, REQ-099]
maintenance: false
related: [REQ-096, REQ-097, REQ-099]
batch: parallel-building
write_set: [docs/work-guide.md, decisions/records/*, decisions/log.md]
---

# Docs + ADR — Multi-Checkout Guide and Session-Ownership Decision Record

## What

User-facing documentation for the new model and a decision record for the re-grain. No ADR currently covers session ownership at all — the 0.161.0 exclusive-session decision was recorded only as an AI report plus REQ-069.

## Detailed Requirements

- **`docs/work-guide.md`:** a "several checkouts against one queue" section — how to claim from a workspace/clone/cloud session, the `assigned_to` earmark, the one-releaser rule, what happens on a double claim (ordinary merge conflict, fixed at merge, `queue-kanban verify` finds duplicates), and the automatic wave. Update the guide's existing "one queue owner, one REQ at a time" summary (~line 91) to match the new contract.
- **ADR in `decisions/records/`** (next number after adr-017) recording: the checkout→queue re-grain and claim-anywhere/one-releaser model; why the reserve *verb/status* stays dead while the advisory *field* returns (ratchet + router budget); capture-anywhere with fix-at-merge instead of id sharding; the auto-wave contract change (supersedes "nothing computes the set"); the static writer label vs. the liveness-machinery ban; links to UR-012/REQ-069 (the decision being partially reversed), UR-018, and the two acceptance-run artifacts (REQ-095, REQ-100).
- **`decisions/log.md`:** one-line entry (its last entry is 2026-07-01 — the log is stale; just append, don't backfill).
- Cross-references by file path per convention; no CLAUDE.md/AGENTS.md citations in shipped files.

## Red-Green Proof

**RED prompt/case:** `docs/work-guide.md:91` still tells users "one queue owner, and by default one REQ at a time" with no multi-checkout path; `decisions/records/` has no session-ownership ADR.
**Why RED now:** The docs describe the 0.161.0–0.166.0 contract; the re-grain decision lives only in plan files and this UR.
**GREEN when:** The guide documents the multi-checkout workflow accurately against the shipped contract text, and the ADR exists with the five decision points and links above.
**Validation:** User confirmed (approved plan, Phase 4).

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 4).

---
*Source: approved plan, Phase 4*
