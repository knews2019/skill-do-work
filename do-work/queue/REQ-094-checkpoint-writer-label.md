---
id: REQ-094
title: Checkpoint writer label — crash recovery ignores foreign entries
status: pending
created_at: 2026-08-04T19:44:17Z
user_request: UR-018
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-095, REQ-096, REQ-097, REQ-098, REQ-099, REQ-100, REQ-101]
batch: parallel-building
write_set: [actions/work-reference.md]
---

# Checkpoint Writer Label — Crash Recovery Ignores Foreign Entries

## What

Give `do-work/CHECKPOINT.md` In-Progress entries a **static writer label** identifying the checkout that wrote them, and scope crash recovery to own-label entries only. Foreign entries are **reported, never stripped**. This defuses a live landmine: the checkpoint is git-tracked, and once two checkouts sync it, checkout A reads checkout B's live claim as its own crash, strips it, and re-queues a REQ someone is actively building — a deterministic replay of the 2026-07-01 collision, no race needed.

## Detailed Requirements

- In `actions/work-reference.md`'s In-Progress Record section (~lines 409–425): each In-Progress entry records `writer: <hostname>:<absolute-checkout-path>` (path alone collides across machines — both sides can be `/home/user/repo`).
- Crash-recovery rule: only entries whose writer label matches **this** checkout are crash-recovery candidates. Foreign entries are listed in the recovery report and left untouched — extend the existing foreign-claim rule (`actions/work-reference.md:239` — "a claim you didn't record is not yours, never touch it") to cover the checkpoint.
- Reword the tripwire at `actions/work-reference.md:413`: refresh intervals, staleness checks, and liveness probes stay **banned by name**; a static writer label is explicitly not liveness machinery (it is written once, never refreshed, never checked for staleness).
- Entries written before this change have no writer label: treat a label-less entry as **own** on the checkout that has it uncommitted/locally modified, otherwise report-only — never guess-strip.
- `CHECKPOINT.md` stays tracked and committed (user decision: "checkpoints are transient, it's fine to commit them before changing, this way different versions of it already is available in the git history").

## Constraints

- Do NOT add: heartbeats, refresh intervals, holder-liveness checks, staleness thresholds, auto-takeover. (Batch-wide do-not-build list — see UR-018.)
- This is a prose/contract change to shipped instruction files; keep wording surgical and in the existing section's voice.

## Red-Green Proof

**RED prompt/case:** Simulate a synced foreign checkpoint: write a CHECKPOINT.md In-Progress entry for a REQ that is in `working/` but was claimed "by another checkout", then follow today's crash-recovery instructions — they classify it as a local crash and strip/re-queue it.
**Why RED now:** In-Progress entries carry no writer identity (deliberately, per the current `:413` tripwire), so recovery cannot distinguish a foreign live claim from an own crash.
**GREEN when:** Following the updated instructions, the same foreign entry is reported ("claim held by <other-writer>, not touched") and the REQ stays claimed; an own-label entry still recovers exactly as before.
**Validation:** User confirmed (ask-tool answer: "Static writer label").

## Full Context

See `do-work/user-requests/UR-018/input.md` and `do-work/user-requests/UR-018/assets/approved-plan.md` (Phase 1).

---
*Source: approved plan `let-s-talk-about-this-kind-robin.md`, Phase 1 item 1*
