---
title: "Lessons from REQ-101: Docs + ADR — multi-checkout guide and the session-ownership decision record"
type: source-summary
topic_cluster: worktree-and-parallel-dispatch
sources: [raw/processed/2026-09-01/REQ-101-docs-adr-multi-checkout-guide-and-the-se.md]
related:
  - page: REQ-094-checkpoint-writer-label-crash-recovery-i
    rel: complements
  - page: REQ-096-execution-model-re-grain-claim-anywhere-
    rel: depends-on
  - page: REQ-097-assigned-to-advisory-field-schema-line-s
    rel: depends-on
  - page: REQ-099-automatic-wave-dispatch-the-work-loop-co
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-101: Docs + ADR — multi-checkout guide and the session-ownership decision record

Part of the [[concept-worktree-isolation-and-parallelism]] cluster.

## What the REQ was about

User-facing documentation for the new model and a decision record for the re-grain. No ADR currently covers session ownership at all — the 0.161.0 exclusive-session decision was recorded only as an AI report plus REQ-069.

## Solution summary

Wrote the decision record for the whole UR-018 re-grain, the user-facing multi-checkout guide, the log line, and the two index updates a new ADR requires to be reachable.

## What worked

- Writing the ADR's Consequences from the batch's *findings* rather than from its plan. Four of the five costs listed (guaranteed checkpoint conflict, the unsound label-less bullet, unproven agent concurrency, nothing travels on untracked installs) are things REQ-095 and REQ-100 discovered; none of them were foreseeable when the plan was approved, and an ADR that omitted them would read as more confident than the evidence supports.
- Recording the two *nearly-taken* alternatives (keep the confirmation gate; schedule on `write_set`) with real reasons rather than listing only the obviously-wrong ones. A rejected alternative that nobody was tempted by teaches a later reader nothing.
- Checking every relative path in the References section by loop. Nothing validates links in `decisions/`, so a dangling reference is permanent and silent.

## What didn't work

- Nearly shipping the ADR without touching either index. `decisions/records/` is not a scanned directory — the topic index and master index are hand-maintained lists, so a record absent from them is unreachable by navigation even though the file exists. Caught by grepping the indexes for `adr-017` to see whether *it* was listed.
- The master index's counts turned out to be **three behind** (header claimed 15 in-force decisions against 18 records, and the cluster row's "4 ADRs + 1 declined" predated nothing but was about to). Adding another hand-maintained number would have made the same failure again on ADR-019, so both counts were replaced with a scope phrase and a "read `records/`" pointer.

## Worth knowing

- `decisions/` is `export-ignore`d and therefore not a shipped path: it may cite `CLAUDE.md`, and it never ships to consumers. `docs/` **is** shipped and may not.
- A new ADR is **five** file changes, not one: the record, the topic index (twice — frontmatter and body), the master index row, and `decisions/log.md`.
- `decisions/log.md` was stale from 2026-07-01 to this entry. The REQ said append, don't backfill, which is right — a hand-written history reconstructed after the fact is worth less than an honest gap.

## Back-reference

See `do-work/archive/UR-018/REQ-101-docs-and-adr-session-ownership.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `e452989`.
