---
title: "Lessons from REQ-280: Probe timestamp ordering, and point Check 12 at the archive repair that already ships"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-01/REQ-280-probe-timestamp-ordering-and-point-check.md]
related:
  - page: REQ-281-reconcile-the-calibration-log-against-th
    rel: complements
  - page: REQ-284-emit-every-verify-finding-from-the-board
    rel: complements
  - page: REQ-304-draw-a-reversed-wait-as-a-break-not-as-a
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-280: Probe timestamp ordering, and point Check 12 at the archive repair that already ships

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

Two gaps in the same loop, both on the read side.

1. **No ordering probe.** `skills/do-work-board/tools/queue-kanban/model.go:1261` is the suite's entire time-consistency surface: one comparison, `completed_at < claimed_at`, inside `detectCompletionAnomaly`. Nothing anywhere checks `created_at <= claimed_at <= completed_at`. Add that ordering probe to `queue-kanban verify`'s probe set and to `actions/forensics.md` Check 12, reported per violated pair.

## Solution summary

Added a timestamp-ordering probe to `queue-kanban verify` covering queue, working and archive, reporting one non-fixable finding per violated pair with both field names, both raw values, and a remedy routed by where the file lives. Rewrote forensics Check 12 to carry the ordering condition alongside its future-stamp condition and to point both halves at the repair scripts that already ship instead of prescribing hand git archaeology, and added the probe to Check 14's table.

## What worked

Running the captured RED before writing code, again. It confirmed the fixture reproduces the blind spot and — more usefully — the *shape* of the REQ's proof was already right, so the test could be the fixture rather than an invention. Second: running the finished probe against the repository's own 250+ REQs. A new failing condition on a gate everyone runs is only safe if you have measured its false-positive rate on real data, and "zero on 250" is the sentence that makes it safe to ship.

## What didn't work

Writing the Go source through an **unquoted** shell heredoc. Backticks inside a Go comment were command-substituted away, leaving `Not Fixable —  does not rewrite stamps` — which compiles, passes vet, passes gofmt, and passes every test. Only re-reading the written block caught it. For any generated file with backticks or `$` in it, quote the heredoc delimiter (`<<'PY'`); a build is not evidence that a comment survived.

Also: the REQ's two named pairs were not the whole rule. Implementing them exactly left `created_at > completed_at` with an absent `claimed_at` passing — and my own carve-out fixture was the thing that walked into it. The instance list was again narrower than the class, which is the third time this session.

## Worth knowing

The ordering predicate now exists in two languages — `scripts/repair-req-timestamps.sh` (repair) and `verify.go` (read). Nothing holds them together mechanically; the Go comment names the shell file and its two boundary decisions (strict comparison, absent stamp is other checks' territory) precisely because that is the seam most likely to drift. If a third spelling is ever proposed, that is the moment to build the shared-fixture harness instead.

## Back-reference

See `do-work/archive/UR-057/REQ-280-probe-timestamp-ordering-and-point-check-12-at-the-shipped-repair.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5e180d0`.
