---
id: UR-135
title: 'Evidence-based activity indicator replaces the assumed-pause badge on the board'
created_at: 2026-10-05T20:11:59Z
requests: [REQ-632, REQ-633]
word_count: 1300
---
# Evidence-Based Activity Indicator Replaces the Assumed-Pause Badge on the Board

## Summary
A design brief for the queue-kanban board. Done cards show `over 4h · assumed pause` whenever the earliest-stamp-to-completion span exceeds `analysisOutlierCeiling = 4h` (`skills/do-work-board/tools/queue-kanban/durations.go`). The rule protects the Durations view's Panel B day medians, but the card presents a statistical exclusion as a fact about the REQ, and continuous 4h+ Route C runs get labelled paused. The brief asks for (A) a correlated-activity indicator on cards, (B) an evidence-based idle judgement from the largest gap between activity events instead of the raw span, and (C) the Panel B calibration exclusion decided from the same gap evidence in both readers (`skills/do-work/actions/estimate-reference.md` and durations.go). Nine design questions (Q1 to Q9) were answered in the session and the maintainer picked the recommended option for each decision put to them.

Decisions recorded in the session (2026-10-05):
- D1 Two REQs: the card fix first, the Panel B rule change second.
- D2 Done cards show nothing new. The largest-gap evidence goes to the detail drawer. The badge is deleted, not reworded.
- D3 One threshold constant, 2h, phase named in the label. The implementation checks the archive's per-phase p95 before fixing the number.
- D4 A commit belongs to a REQ when any of three signals matches: it touches the REQ's own file or run artefacts, its subject carries the `[REQ-NNN]` prefix, or it is reachable through a matched merge commit's second parent or the live `worktree-agent-REQ-NNN-*` branch tip.

Facts established in the session that the brief did not know: builder commits stay reachable from main after a `--no-ff` merge, so `git log --all` and reflog are not needed; the verdict string `paused` has readers in the timeline forecast, the UR progress summary and the board guide, not only Panel B; `do-work/calibration-log.tsv` has no stamp or gap column, so part C needs a new column; the verify probes already run git per request outside the board's mtime cache and are the right home for the new git read.

## Extracted Requests
- A, the card half of B, and the minimal set from Q9 → REQ-632 (Board cards show last correlated activity and drop the assumed-pause badge)
- The threshold rule from B, C, and the migration from Q8 → REQ-633 (Panel B and the calibration log exclude by largest stamp gap, not raw span)

## Batch Constraints
- REQ-633 depends on REQ-632: the badge and its text field must be gone before the verdict string is renamed, and the gap computation REQ-632 adds is what REQ-633's drawer and Panel B wording refer to.
- Both change shipped files, so each is a release.
- The timing stream under the git common dir is out of scope for both; it is not a liveness source.
- No change to the board's mtime fingerprint; git reads live in the per-request slot beside the verify probes.

## Full Verbatim Input
> ```
>  I want to discuss a design change to the queue-kanban board before anything is built.
>  Read this whole brief, then answer the questions at the end with options and a
>  recommendation. Do not implement yet — I will pick an option with you first.
> 
>  ## The problem
> 
>  Done cards carry the badge "OVER 4H · ASSUMED PAUSE" whenever the wall span (earliest
>  lifecycle stamp → completed_at) exceeds 4 hours. In my repo three of the last five finished
>  REQs carry it, and all three were worked continuously. The badge is a guess from a single
>  number, while the repo already holds the evidence that would disprove it.
> 
>  Where the rule lives today:
> 
>  - do-work-board/tools/queue-kanban/durations.go — `analysisOutlierCeiling = 4 * time.Hour`,
>    `dayMedianExclusionReason` returns "paused" for any span over the ceiling,
>    `implementationSpanPausedBadgeText` builds the badge, `measureImplementationSpan`
>    reads only frontmatter stamps.
>  - do-work/actions/estimate-reference.md → Calibration: "exclude spans > 4h or negative
>    (assumed user pauses / broken stamps)". durations.go calls itself "its second reader".
>  - do-work-board/tools/queue-kanban/web/board-cards.js — `makeImplementationSpanNode`
>    prints "wall time …" and the badge; the tooltip says "Duration-quality marker only …
>    excluded from duration medians".
> 
>  The 4h rule was written to keep one paused session from inflating the Panel B day medians.
>  That is a fair goal. The badge then leaked a statistical exclusion rule onto the card as
>  if it were a fact about the REQ. I want the card to show what actually happened.
> 
>  ## Evidence from one real REQ (REQ-2248, Route C, wall time 4h 21m)
> 
>  Frontmatter stamps (UTC): claimed 11:25, planning 11:36, dispatch 12:11,
>  builder_handback 13:16, integration 14:23, remediation 15:14, review/re_review 15:40,
>  completed 15:46. Largest gap between consecutive stamps: ~67 min. Nothing was paused.
> 
>  `buildPhaseBreakdown` in durations.go already computes exactly these per-phase gaps for
>  the detail drawer. The badge never consults them.
> 
>  Git history over the same window (local time +03:00):
> 
>  - `[REQ-2248] claim request lifecycle` 14:25, `[REQ-2248][REQ-2261] Record Route C
>    triage …` 14:29, then no `[REQ-2248]`-prefixed commit until `[REQ-2248] Integrate owned
>    phone answer …` 17:22 — a 2h53m "gap" if you grep by message prefix.
>  - But in that same window main received seven commits with other prefixes
>    (`docs(do-work): validate phone completion …`, `Record wave 15 builder plans …`,
>    `Record neighbour handback …`) that each touch
>    `do-work/working/REQ-2248-reveal-rings-….md`. Activity never stopped.
>  - The builder worked in a `worktree-agent-REQ-2248-*` worktree; its commits arrived on
>    main as `Integrate …` merge commits, so main-only history under-reports builder activity
>    during the dispatch → handback phase.
> 
>  Lesson for the design: correlating by `[REQ-NNNN]` message prefix is lossy. The robust
>  key is the commit's touched paths — the REQ file (`do-work/queue|working|archive/REQ-NNNN-*.md`)
>  and its run artefacts (`do-work/runs/*/REQ-NNNN*`) — with the prefix as a second signal,
>  and worktree branches (`worktree-agent-REQ-NNNN-*`, codex `req-nnnn-*`) as a third.
> 
>  ## What the board already has that this can build on
> 
>  - Phase stamps per REQ (`claimed_at`, `planning_at`, `dispatch_at`, `builder_handback_at`,
>    `integration_at`, `review_at`, `remediation_at`, `re_review_at`, `completed_at`),
>    parsed in model.go `lifecycleTimestampFields`. Append-only by contract
>    (do-work/actions/work-reference.md).
>  - `buildPhaseBreakdown` (durations.go) — per-phase elapsed time between stamps.
>  - One existing git read: `lookupGitCommitDate` in model.go runs
>    `git log -1 --format=%cI <hash>` as the completion-time fallback, injectable for tests.
>    verify.go runs more git probes (`rev-parse`, `worktree list`, `branch --list
>    worktree-agent-*`, `status --porcelain`) on every request, outside the mtime cache.
>  - `staleClaimThreshold = 3h` in verify.go, explicitly "NOT a liveness test".
>  - The Activity view (activity.go) — one row per lifecycle stamp, newest first.
>  - A per-REQ timing event stream under `<git-common-dir>/do-work-timing/` (never
>    committed) folded into a `## Timing` body section at archive time. REQ-2248's says
>    "3h 45m 40s total, 1h 44m attributed … 2h 01m unattributed". The board reads neither.
>  - Serving model: the server rebuilds only when a file mtime changes (serve.go);
>    the page has no polling or auto-reload; a 1s `setInterval` only re-renders relative
>    time text ("52min ago").
> 
>  ## What I want
> 
>  A. An activity indicator on each REQ card — at minimum "last commit N min ago" where the
>     commit is correlated to the REQ (see the lossy-prefix lesson above). On in-progress
>     cards this is the headline signal: a big number means something is stuck. On done
>     cards it is the honest replacement for the assumed-pause guess.
> 
>  B. The pause judgement, if the card keeps one at all, should be evidence-based: "largest
>     gap between any two activity events (stamps ∪ correlated commits) exceeded X" rather
>     than "total span exceeded 4h". A continuous 4h21m REQ shows no badge; a REQ with a
>     2h hole in both stamps and commits shows "idle 2h between dispatch and handback".
> 
>  C. The Panel B calibration exclusion can keep its own rule, but it should be decided from
>     the same gap evidence, not from the raw span, so the medians stop throwing away long
>     continuous sessions. Both readers (estimate-reference.md and durations.go) must change
>     together — the doc calls itself the single definition.
> 
>  ## Questions to discuss with me (answer each with options and a recommendation)
> 
>  Q1. Definition of "activity": stamps only, commits only, or the union? What counts as a
>      commit "belonging" to a REQ — touched-path match, message prefix, worktree branch
>      name, or a scored union? Give the false-positive / false-negative trade-off of each.
> 
>  Q2. Where does the git read live and how expensive is it? One `git log --since=<claimed_at>
>      --format=%cI --name-only` per board rebuild, filtered in Go, versus one `git log -1 --
>      <paths>` per REQ, versus reading the uncommitted timing stream under the git common
>      dir. Consider the mtime-cache: commits change `.git/` not `do-work/`, so a rebuild may
>      not trigger; what is the cheapest correct invalidation?
> 
>  Q3. Which cards show the indicator, and in what form? In-progress cards (claimed / planning
>      / dispatched / integrating): "last activity 12 min ago", ticking like the relative
>      time already does. Done cards: "largest idle gap 1h07m (dispatch → handback)" or
>      nothing when below threshold. Blocked / pending-answers: probably nothing. Push back
>      if you think the indicator belongs in the detail drawer instead.
> 
>  Q4. What threshold replaces 4h, and is it one number? Phase-aware (a dispatch → handback
>      gap of 90 min is normal for Route C; a planning → dispatch gap of 90 min is not)
>      versus one global gap ceiling. I lean to one number plus the phase name in the
>      label, and I want you to argue the other side.
> 
>  Q5. Should the badge text change class? "assumed pause" asserts a cause. "idle 2h01m
>      (dispatch → handback)" states an observation. Does the calibration exclusion reason
>      ("paused") need renaming to match, and what else reads that string?
> 
>  Q6. Worktree builders: during dispatch → handback the builder commits on a branch that may
>      be deleted after integration. Should the board read `git log --all` / reflog, or should
>      the handback stamp (`builder_handback_at`) be accepted as the activity proof for that
>      phase? What breaks if the branch is gone?
> 
>  Q7. The timing stream already attributes/un-attributes minutes per REQ. Is it the better
>      source than git for "active vs idle", and if so why has the board not read it so far
>      (never committed, lives outside the tree)? Would making the board read it be simpler
>      than git correlation, or worse because it is invisible to a fresh clone?
> 
>  Q8. Migration: existing archived REQs with spans > 4h currently sit outside Panel B. If the
>      rule changes to gap-based, their medians shift. Say whether that is acceptable or
>      whether the calibration log needs a version marker.
> 
>  Q9. Attack the premise: is there a simpler fix than any of the above — for example drop the
>      badge from cards entirely, keep the 4h rule as a Panel-B-only exclusion, and add just
>      "last commit N min ago" on in-progress cards? Tell me if that is enough.
> 
>  Format: answer Q1–Q9 in order, options with short labels, one recommendation per question,
>  then a single proposed design in under 30 lines. Cite upstream file paths for every claim
>  about current behaviour. Then stop and wait for my decisions.
> ```
