---
title: "Lessons from REQ-094: Checkpoint writer label — crash recovery ignores foreign entries"
type: source-summary
topic_cluster: checkpoint-and-crash-recovery
sources: [raw/processed/2026-09-01/REQ-094-checkpoint-writer-label-crash-recovery-i.md]
related:
  - page: concept-session-checkpoints-and-recovery
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-094: Checkpoint writer label — crash recovery ignores foreign entries

Part of the [[concept-session-checkpoints-and-recovery]] cluster.

## What the REQ was about

Give `do-work/CHECKPOINT.md` In-Progress entries a **static writer label** identifying the checkout that wrote them, and scope crash recovery to own-label entries only. Foreign entries are **reported, never stripped**. This defuses a live landmine: the checkpoint is git-tracked, and once two checkouts sync it, checkout A reads checkout B's live claim as its own crash, strips it, and re-queues a REQ someone is actively building — a deterministic replay of the 2026-07-01 collision, no race needed.

## Solution summary

In-Progress checkpoint entries now carry `— writer: <hostname>:<absolute-checkout-path>`; the canonical contract (In-Progress Record) defines the format and derivation, append/refresh/removal are scoped per-id-per-writer to own-label entries, and Crash Recovery classifies four cases (own-label / foreign-label / label-less legacy / unnamed-or-absent). Foreign-label entries report `claim held by <writer>, not touched` and never enter the 3-hour takeover ladder. The tripwire keeps banning refresh intervals, staleness checks, and liveness probes by name (dropping "holder id" from the ban, carving out the static label); `never grow into one` survives byte-identical. Both label-destruction paths (Step 10 wholesale rewrite, session-start delete) are own-label-scoped with foreign entries preserved verbatim. Three new contract-suite assertions pin the label, the foreign-label report phrase, and the tripwire.

## What worked

**What worked:** A pre-build exploration inventory of every restatement site (3 copies of "no second owner reads it", 2 non-obvious label-destruction paths in Step 10's wholesale rewrite and the session-start delete, 5 pinned contract-test phrases) — the build touched 8 files with zero suite breakage because the collision surface was mapped first. Pinning the tripwire ban and its carve-out to the *same paragraph* via a new assertion (a carve-out that drifts into its own paragraph reads as general permission).

**What didn't:** The builder's Step 10 echoes narrowed "every entry this checkout did not write" to "another checkout's label" — an echo written from memory of the canonical clause, not from it. Echo sites should quote the canonical condition, not paraphrase it (that's REQ-102).

**Worth knowing:** The checkpoint travels between checkouts on any install that commits `do-work/` — every rule about it now has four claim-origin cases (own-label / foreign-label / label-less / unnamed), and the three-hour takeover ladder serves only the last two. `checkpointMentionedRequestIds` in `tools/queue-kanban` extracts ids by regex, so entry-format suffixes are parser-transparent.

## Back-reference

See `do-work/archive/UR-018/REQ-094-checkpoint-writer-label.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `9c305c0`.
