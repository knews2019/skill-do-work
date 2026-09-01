---
title: "Lessons from REQ-098: Verify probes: assigned-elsewhere-claimed-here and ur-archived-with-live-member"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-098-verify-probes-assigned-elsewhere-claimed.md]
related:
  - page: concept-contract-verification-gates
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-098: Verify probes: assigned-elsewhere-claimed-here and ur-archived-with-live-member

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

Two new read-only probes in `queue-kanban verify`, extending the report-and-route contract (verify never repairs — fixes belong to `actions/cleanup.md`):

1. `assigned-elsewhere-claimed-here` — a REQ carrying `assigned_to` is sitting in `do-work/working/` (someone claimed work earmarked for another session without clearing the field).
2. `ur-archived-with-live-member` — an archived UR has a member REQ (by `user_request:` scan) still in `queue/` or `working/` (a silent-merge class also reachable from a botched cleanup).

## Solution summary

Added the two read-only probes. Neither is fixable, both name their remedy as something cleanup *asks* about, and the renderer needed no change because it was already generic.

## What worked

- Writing the array-versus-frontmatter divergence test *before* the probe. `RequestIds` is already the right data, so the probe would have been correct by accident; the test is what makes it correct on purpose, and it is the one that would catch a future "optimization" to read the `requests:` array instead.
- Asserting the **negative** in probe 2's firing test — that the terminally-resolved sibling is *not* named. A probe that lists every member of the UR would have passed a test that only checked for the live one.

## What didn't work

- Rebuilding the binary from inside `tools/queue-kanban` and then invoking it by a repo-relative path in the same command line — the `cd` had already moved the shell, so the invocation failed with `No such file or directory` and read like a build failure. Wrap the build in a subshell (`(cd … && go build …)`) so the outer shell never moves. Same class of stale/wrong-path confusion REQ-097 hit twice.

## Worth knowing

- `renderVerifyReport` is generic over findings, so a new probe is three edits (constant, function, runner call) and never a renderer change. The temptation to add per-category formatting is what would break that.
- `Fixable` means *cleanup can resolve this mechanically*, not *this is minor*. Setting it on a human-decision finding would make `do-work cleanup` advertise a repair it must not perform unasked.
- `UserRequestTicket` carries no `TreeSection`. Anything that needs a UR's location reads `FilePath` — and should match `/do-work/archive/` with both separators, not the bare word.

## Back-reference

See `do-work/archive/UR-018/REQ-098-verify-probes-assignment-ur.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `47cd408`.
