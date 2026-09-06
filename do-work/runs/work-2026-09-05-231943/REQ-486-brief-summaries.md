# Builder brief — REQ-486, increment 2 of 2: the progress summaries

**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-486-collapsible-ur-progress-summaries`
Increment 1 (the folds) is merged and the worktree is fast-forwarded onto it. Build on top.
**Route:** C. **TDD: yes.** **Impact: user-visible.**

## Read first

1. `do-work/runs/work-2026-09-05-231943/REQ-486-plan.md` in the MAIN checkout — **tasks T3, T4 and T5 are yours**, with
   decisions D1, D4, D5, D7, D8, D10, D11, D12, D13, D14 and D15 and the per-lane testing argv. It is
   the authority; this brief is a pointer.
2. `do-work/runs/work-2026-09-05-231943/REQ-486-handback-folds.md` — what increment 1 actually built and the three
   things it deliberately left for you.
3. `do-work/runs/work-2026-09-05-231943/REQ-486-exploration.md` — the evidence with file:line anchors.
4. `do-work/working/REQ-486-collapsible-ur-progress-summaries.md` in the worktree — the request, its
   `## Plan`, `## Scope` and acceptance criteria.
5. `_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`,
   and `_dev/primes/prime-releases.md` (T5 includes the release).
6. Crew rules: general.md, shared-principles.md, coding-guardrails.md, communication-style.md,
   testing.md, ui-design.md, frontend.md, anti-slop.md.

## The three tasks

**T3 — read `estimate.p50_active_minutes` into the Go model and payload.** No reader exists today.
`parseFrontmatter` already returns the whole `estimate` block as a map when it parses strictly, so
this is a lookup, not a parser change. Ship presence and value as a **pair** (`hasEstimateP50ActiveMinutes`
with `omitempty`, `estimateP50ActiveMinutes` without), copying the pinned precedent at
`generate.go:224-228` — the salvage path drops nested maps by design, so the tests must assert
**absence, not zero**. Also correct `timeline.go`'s comment that gives "the board parses no nested
frontmatter blocks" as its reason; that reason stops being true. Do not change Timeline behaviour.

**T4 — the shared rollup, its two call sites, and one clock instant.** One function computes the five
figures; both surfaces render from it and must agree byte-for-byte. New fragment
`web/board-user-request-summary.js` at manifest position 7 (the manifest is pinned twice in
`generate_test.go` — miss either copy and the fast stage fails loudly). Predicates and the clock
fan-out stay in `board-core.js`. **No tick-subscriber registry** (D5): one `refreshTickingSurfaces()`
captures one `nowMs`, calls the existing `refreshRelativeTimeNodes(nowMs)` **first**, then re-renders
`[data-ur-summary-id]` nodes.

**T5 — the browser lane, the docs and the release.** Wrap/collision at 320, 768 and 1280 CSS px in
both themes; drawer containment with a 43-member UR; contrast measured against
`getComputedStyle(document.body).backgroundColor`, driven through the colour-scheme flag or Chromium
resolves `prefers-color-scheme` to light and nothing automated ever sees the dark palette; real-button
tab order. Then `board-guide.md`, `board.md`, `work-reference.md`, and the release trio.

## The single biggest risk, and it is not a coverage gap

`web/board.js:68` is the board's **only** ticker — every claim stopwatch, every relative time, every
state timer, the clock-skew tooltip. T4 points it at a function that runs a second pass. If anything
in that pass throws, the interval callback dies: the board renders perfectly on load and then nothing
updates again, which looks exactly like a board full of very young claims. The current suite cannot
see it — `setInterval` never runs inside a Node probe, every probe calls render functions directly,
and the browser lane does not wait a second and re-measure.

Three containments, in order: the existing refresh runs **first** with the captured instant, so an
unguarded throw in the new pass cannot cost the existing surfaces their tick; keep the rollup total by
**narrowing, not by try/catch** — a swallowed exception hides the bug instead of the freeze, which is
worse; and write the positive freeze-guard assertion the plan names, driving the tick against an
incomplete UR payload and asserting an existing stopwatch node still shows the new label.

**Second risk (D1).** Active time is **origin-to-completion**, summing the existing
`implementationSpanMinutes`, never claim-to-completion arithmetic in the browser. A cancelled member
counts toward resolved but ships no span; the temptation is to subtract `claimedAt` from `completedAt`
in JavaScript. That produces a plausible number with no outlier rule and no origin correction, for
exactly the REQs whose bookkeeping was worst. **If your diff contains arithmetic on a raw timestamp
inside the summary fragment, that is the failure**, regardless of whether the number looks right.

## Two things increment 1 left for you, on purpose

- `.ur-count` still shows the group total. The plan has it dropping to a filter-only "n of m shown"
  (D11), but the summary strip that owns the total is yours — make that change **in the same commit
  that adds the strip**, never before.
- Nothing yet proves tab order or narrow-width layout. Both are T5's browser probe.

## Testing — three lanes, and a lane that skips reports success

Record each lane's own exit line, not the gate's. Increment 1's greens, for comparison: fast stage
`ok … 60.824s`; strict JavaScript `ok … 7.691s` with 65 PASS and 0 SKIP; strict browser
`ok … 94.899s` with 34 PASS and 0 SKIP on Chromium 141.

```
cd skills/do-work-board/tools/queue-kanban && \
  QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off \
  DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior go test -count=1 ./...

QUEUE_KANBAN_BROWSER_PROBES=off QUEUE_KANBAN_JAVASCRIPT_PROBES=on \
  QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 go test -count=1 -run '^TestJavaScriptBehavior' -v .

QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=on \
  QUEUE_KANBAN_STRICT_BROWSER_BEHAVIOR=1 QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  go test -count=1 -run '^TestBrowserBehavior' -v .
```

Closing gate: `bash _dev/tests/maintainer-verify.sh --heavy` with `QUEUE_KANBAN_BROWSER` exported,
exit 0, with both heavy-lane lines present rather than their skip lines.

**A known flake, not yours:** increment 1 saw one heavy gate run fail a single case in
`_dev/tests/prescribed-shell-cases/qualify.sh` (untracked-file TODO scan). Run alone it passes at both
revisions, and the second full gate run passed. If you see it, re-run once and say so.

## Environment

```
env -u NODE_OPTIONS \
  -u GIT_CONFIG_COUNT -u GIT_CONFIG_KEY_0 -u GIT_CONFIG_KEY_1 -u GIT_CONFIG_KEY_2 \
  -u GIT_CONFIG_VALUE_0 -u GIT_CONFIG_VALUE_1 -u GIT_CONFIG_VALUE_2 \
  GIT_CONFIG_GLOBAL=/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/gitconfig-gate \
  QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  <command>
```

Capture exit status from `$?` — never pipe a gate to `tail`. 4 CPUs, other builders are running: run
the full heavy gate at most once, at the end.

## The release is yours (T5)

This change touches shipped files under `skills/`, so `_dev/primes/prime-releases.md` makes it a
release. **Do not invent the version number**: `VERSION`, `skills/do-work/VERSION` and
`skills/do-work/actions/version.md:5` all read `0.303.10` today, but four other requests are held
unreleased in this same run and their finalization may bump it first. Write the changelog entry and
the mirror, bump from whatever the files say **at the moment you write them**, and state in the
hand-back exactly what you bumped from and to, so the orchestrator can reconcile it at finalization.

## Hand-back

Commit on your branch, message ending `Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>`. Write
your report to `/home/user/skill-do-work/do-work/runs/work-2026-09-05-231943/REQ-486-handback-summaries.md` in the MAIN
checkout. Do not stage or commit it, do not push, do not touch any other worktree.
