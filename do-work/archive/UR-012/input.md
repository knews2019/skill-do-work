---
id: UR-012
title: Exclusive-session model replaces the concurrency machinery
created_at: 2026-08-01T12:36:44Z
requests: [REQ-069]
word_count: 12
---

# Exclusive-session model replaces the concurrency machinery

## Summary

The user asked whether coexisting do-work sessions were planned for, then worked through
the answer across a conversation and revised the resulting analysis report directly. The
final position: concurrent sessions on one checkout are **outside the product contract**.
The pipeline should assume exclusive repository use and delete the machinery that models
unsupported concurrent actors.

The design is specified in a report the user edited:
`ai-reports/2026-08-01_1231_lock-simplification-analysis/index.html`

## Extracted Requests

| REQ | Title |
|---|---|
| REQ-069 | Adopt the exclusive-session model and remove the concurrency machinery |

## Batch Constraints

- Removal must leave the pipeline coherent at every commit — the cut blocks reference each
  other (cleanup's live-claim gate reads the lock file the Lock Guard section defines).
- The replacement is capped at ~200 words total. Exceeding that cap defeats the request.
- No new frontmatter fields, lock files, retry artifacts, or other durable state.

## Full Verbatim Input

plan changed a bit, read it again and
capture it as a REQ file:///Users/t2/Desktop/e1-experimental-repos/skill-do-work2/ai-reports/2026-08-01_1231_lock-simplification-analysis/index.html

### Referenced design document (user-edited), verbatim content

The report at the cited path states, in the user's revised framing:

**Verdict (hero):** "Assume one session. Work the REQ in front of you." — At least 6,511
words of concurrency machinery, a sixth of the core pipeline, can be removed and replaced
by one operating rule: a do-work session assumes it is alone. Unexpected repository state
matters only when it prevents the active REQ from being implemented, tested, archived, or
committed; otherwise preserve it and continue.

**The boundary that makes deletion safe:** The pipeline supports one active do-work
session, one active REQ, and one coder context. Parallel sessions, co-dispatched builders,
and cross-session ownership are outside the product contract. The pipeline does not detect,
coordinate, or recover an unsupported concurrent run. "The 2026-07-01 collision explains why
the lock was built; it no longer defines the product boundary. If two do-work sessions are
run against one checkout anyway, behavior is unspecified. The pipeline spends no
instructions or state trying to make that unsupported situation safe."

**The replacement rule — local decision ladder:**
1. The active REQ can proceed — continue. Preserve unrelated repository state, exclude it
   from this REQ's staging, and spend no time explaining or repairing it.
2. The active REQ is blocked — investigate. A failing acceptance check, test, archive move,
   or commit is ordinary work on the current request.
3. Three coder attempts fail — stop the local loop. Count attempts in the coder session
   only. Report the unresolved blocker; use `pending-answers` only when progress genuinely
   requires a user decision.

**No persistence is added.** The coder counts consecutive fix attempts in its current
context. The existing three-attempt limit stops an in-session loop; a restarted coder
session starts fresh. Normal qualification, tests, explicit staging, and commit success
remain the evidence that this REQ landed.

**Resolved decisions (previously open questions):**
- *Attempt count: coder session only.* Count consecutive fix attempts in the current coder
  context. After three failures, stop the local retry loop and report the blocker. Do not
  add frontmatter or create a file; a restarted coder session starts with a fresh count.
- *Git/index failures: investigate only when blocking.* There is no proactive
  `.git/index.lock` retry or concurrency branch. If Git cannot commit the active REQ, that
  is an observed blocker: investigate it through the ordinary failure path. If repository
  state does not prevent this REQ from finishing, preserve it and continue.

**Measured cut list (6,511 words):** Lock Guard section 4,642 · write_set dispatch gate 483 ·
crash-recovery concurrency gate 378 · cleanup live-claim gate 341 · Step 5.5 re-validation
323 · Step 3 Route-A re-validation 251 · checkpoint-delete gate 93.

**Additions (≤200 words):** exclusive-session invariant + current-REQ relevance rule ≤125 ·
coder-context attempt count + three-attempt stop condition ≤75.

**Keep:** crash recovery proper (920 words — a crashed session still leaves files in
`working/`) · qualification, tests, explicit staging, archive and commit checks (existing).

**Supporting evidence:** REQ-018 built the lock guard; REQ-035, REQ-044 and REQ-047 each
fixed a correctness hole in it afterwards. The write_set gate shows the same pattern —
REQ-032 built it, REQ-036 and REQ-045 patched it.

## Capture-time answers

- **`write_set` field disposition:** keep the `write_set:` frontmatter line and the board's
  overlaps badge; delete only the dispatch *rule* that gates co-dispatch on it.
  `tools/queue-kanban/model.go` is untouched.
- **Slicing:** one REQ, not a batch — the cuts are interdependent and splitting them leaves
  broken intermediate states.

---
*Captured: 2026-08-01T12:36:44Z*
