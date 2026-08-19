# REQ-248 builder brief

**Route B.** Estimated 30 active minutes (P50, medium confidence).

**This is the highest-value item in the queue and it is a live visual defect on this repository's own board.** Panel B's leftmost bar sits in the axis gutter, and a one- or two-day board renders Panel B entirely off-canvas.

**Generate a board and look at it** — at one, two and many active days — and measure in the live DOM. A passing assertion is not evidence about two glyphs sharing a coordinate. `_dev/primes/prime-kanban-board.md` is your entry point and carries the board conventions, including this one: a measured face is per-browser, so record the build beside any number you measure.

## How this build runs

You are a **worktree builder** dispatched by the do-work work pipeline (`skills/do-work/actions/work.md`). Read that file's Step 6 expectations if you need them; everything binding is repeated here.

**Your tree, your branch.** Work only inside `/home/user/skill-do-work-worktrees/worktree-agent-REQ-248-anchor-durations-day-buckets-to-utc-midnight`. That is a full checkout of this repository on branch `worktree-agent-REQ-248-anchor-durations-day-buckets-to-utc-midnight`, cut from the integration tip `67dae6b`.

- Never write anything under `/home/user/skill-do-work` — the main tree belongs to the orchestrator. The one exception is your hand-back file, named below.
- Never read or write anything under `do-work/` in your own worktree. Your worktree carries a **stale snapshot** of the queue; the live queue lives in the main tree and is the orchestrator's alone. Your REQ body is inlined in this brief — that is your copy of it.
- Commit your work on your own branch, as many commits as the work naturally splits into. Do not bump `VERSION`, do not touch `CHANGELOG.md`, do not touch `skills/do-work/actions/version.md` — those are serial-only files the integrator owns, and a builder bumping them races every sibling.
- If you need one line of wiring in a file outside your write set — a shared registry entry, a pointer in someone else's doc — **do not edit it**. Hand back the exact line and where it goes as an *integration seam*; the orchestrator applies it inside the merge commit.
- If you discover you need a file outside your declared scope, stop and report it in your hand-back rather than writing it silently.

**Crew rules load from your own worktree** (they ship, so they are there at the same paths): read `skills/do-work/crew-members/general.md`, `skills/do-work/crew-members/coding-guardrails.md`, and `skills/do-work/crew-members/communication-style.md` before you write code. This REQ is `tdd: true`, so also read `skills/do-work/crew-members/testing.md` and follow RED → GREEN → REFACTOR. Read every path in the REQ's `prime_files` too — those primes encode prior mistakes in this exact area.

**P-A-U phasing is mandatory.** Your REQ body carries an `AI Execution State (P-A-U Loop)` block. Work it: [PLAN] a brief technical approach, [APPLY] inside declared scope only, [UNIFY] run `git diff --stat`, run the linters, verify no debug artifacts, and list each file you checked. The orchestrator audits the checked boxes against the diff — a checked [UNIFY] over a diff containing stray instrumentation is a qualification failure. Record the P-A-U evidence in your hand-back, not in a `do-work/` file.

**Log significant decisions as D-XX** in your hand-back with reasoning. A reversible, low-reach choice is DECIDE & STATE (reasoning only); an irreversible, taste-dependent, or genuinely contestable one is ESCALATE — add `Value:` and `Risk:` lines.

**Out-of-scope finds** go in a `## Discovered Tasks` list in your hand-back. Do not fix them inline.

## Environment notes for this checkout

This is a Linux container, not the maintainer's machine. Before you start, know:

- `bash _dev/tests/maintainer-verify.sh` is the canonical pass/fail gate and it **exits 0 at your branch point** — that is your baseline. Exit code zero is the only proof; never pipe it through `| tail` or judge it from a summary. It takes a few minutes.
- The toolchain was installed for this session: Go 1.26.1, ShellCheck 0.11.0, `just` 1.21.0. They are on `PATH`.
- **Never run a bare `go build`** in `skills/do-work-board/tools/queue-kanban/` — it drops an 11 MB `queue-kanban` binary into the source tree, which is gitignored (so nothing warns you) and multiplies through the installer probe's copies. Build with `go build -o /tmp/<name> .`.
- Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp anything. Never carry a timestamp forward and never compute one.
- Browser tooling may drop a `.playwright-mcp/` directory into the main tree without you issuing a write. It is gitignored and holds sibling agents' evidence — if you clean it, remove only your own files.

## Hand-back

When you are done, write your report to exactly this absolute path:

```
/home/user/skill-do-work/do-work/runs/work-2026-08-18-162355/REQ-248-handback.md
```

That is the one main-tree path you may write. Never stage it, never commit it — it is an orchestrator working file, not branch content.

Structure it as:

```markdown
# REQ-NNN hand-back

**Branch:** worktree-agent-REQ-248-anchor-durations-day-buckets-to-utc-midnight
**Commits:** <short hashes, oldest first>

## What I built
## File manifest
- `path` (new|modified|deleted) — one factual line each
## P-A-U evidence
## Testing evidence
Red-green: the test name, the failure text before, the pass after. Quote real output — never a transcript from a prototype or from memory.
Full gate: the `maintainer-verify.sh` exit code you actually observed.
## Decisions (D-XX)
## Integration seams
Exact lines and where they go, or "none".
## Discovered Tasks
## Pushback
Anything in this brief you think is wrong, with your evidence.
```

**One standing warning, from the previous session's own record.** Five consecutive REQs shipped a mechanism that looked like it closed a class and closed only the instance — and in three of five the remaining hole was exactly where the real data lives. Assume your first fix has that shape and go looking for the hole before a reviewer does. A passing assertion is not evidence about the thing you did not sample.

---

# Your REQ (verbatim copy — the live one lives in the main tree)

---
id: REQ-248
title: Anchor the Durations day buckets to UTC midnight so Panel B stays on canvas
status: claimed
created_at: 2026-08-18T13:54:59Z
claimed_at: 2026-08-18T16:09:27Z
route: B
status_changed_at: 2026-08-18T13:54:59Z
user_request: UR-051
addendum_to: REQ-242
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-18T16:10:30Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - browser evidence
    - cross-route regression gates
---

# Anchor the Durations Day Buckets to UTC Midnight So Panel B Stays on Canvas

## What

Panel B's bars are placed with `xOfEpoch`, which maps each day bucket's **midnight**, while `timeStart` is the **first completion instant**. The two disagree by however far into its first day the earliest sample falls, so the leftmost bar is drawn to the left of the plot area — and on a board with one or two active days the disagreement dominates the whole span and the panel renders off-canvas entirely.

## Context

Found by REQ-242's builder as an unrelated pre-existing quirk, then confirmed and extended by REQ-242's independent review. This is not cosmetic at low day counts.

## Instances

- [ ] **Leftmost bar sits in the axis gutter on the real board.** `x=37.1 width=12` spans 37.1–49.1, entirely left of `DURATIONS_MARGIN_LEFT` (54), and the render shows it struck through by the "0" axis tick. Visible on this repository's own board today.
- [ ] **One active day: Panel B renders empty.** `timeSpan` collapses to the intra-day sample span (about 3 hours), so `xOfEpoch(midnight)` maps to roughly minus three plot-widths — measured annotation at `x=-3330`, bar at `x=-3342`. Both completely off-canvas.
- [ ] **Two active days: same failure, smaller magnitude.** Annotation at `x=-336.5`, bar at `x=-348.5`.

## Requirements

- Every Panel B bar renders inside the plot area at every day count, including one and two active days.
- The slowest-day annotation renders on canvas at every day count — it exists to state a value a clipped bar cannot, and cannot do that from off-screen.
- No change to `DURATIONS_MEDIAN_TITLE_Y` or `describeAtPointer`'s A/B boundary.
- REQ-241's and REQ-242's guarantees hold unchanged: 0 same-row label overlaps, 0 label/mark overlaps, the annotation clear of every neighbour in its strip.

## Builder Guidance

The suggested root fix from the review is to floor `timeStart` to its UTC midnight and ceil `timeEnd` to the following midnight before computing `timeSpan`, so the axis domain and the day buckets share one origin. Verify that against the other panels before adopting it — Panels A and C read the same domain.

**Generate a board and look at it**, at one, two and many active days. Measure in the live DOM.

## Red-Green Proof

**RED prompt/case:** a test asserting every Panel B bar's x-range and the annotation's x-range fall inside the plot area, evaluated on one-day, two-day and many-day fixtures.
**Why RED now:** measured `x=-3330` on a one-day board and `x=37.1` against a left margin of 54 on the real board.
**GREEN when:** the assertion passes at every day count and a render at one and two days shows Panel B populated.
**Validation:** Review finding on REQ-242; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---

## Triage

**Route: B** - Medium

**Reasoning:** The defect is reproduced and the suggested root fix is named, but it has to be checked against Panels A and C which read the same domain, and the day-count evidence needs live measurement — the what is clear, the blast radius is not.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modify) — anchor the axis domain and the day buckets to one origin
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — the day-count lock-in probe

**Files I will NOT touch:**
- `durations.go` and `durations_test.go` — REQ-252 owns the measured-face provenance work and is gated behind this REQ.
- `DURATIONS_MEDIAN_TITLE_Y` and `describeAtPointer`'s A/B boundary — named unchanged by the REQ.

**Acceptance criteria (restated from REQ):**
- [ ] Every Panel B bar renders inside the plot area at every day count, including one and two active days.
- [ ] The slowest-day annotation renders on canvas at every day count.
- [ ] No change to `DURATIONS_MEDIAN_TITLE_Y` or `describeAtPointer`'s A/B boundary.
- [ ] REQ-241's and REQ-242's guarantees hold unchanged: 0 same-row label overlaps, 0 label/mark overlaps, the annotation clear of every neighbour in its strip.

## Pre-Flight

**Git:** ✓ clean outside `do-work/`
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 (recorded in `do-work/working/baseline.json`)
**Dependencies:** ⚠ this checkout needed Go 1.26.1, ShellCheck 0.11.0 and `just` installed before the baseline could run at all, and one pre-existing Linux-only test failure had to be fixed first (0.212.8) — see the REQ brief.

*Checked by work action*
