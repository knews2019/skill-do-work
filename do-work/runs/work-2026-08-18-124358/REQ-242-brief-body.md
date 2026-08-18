---
id: REQ-242
title: Stop Panel B's slowest-day annotation colliding with its own title
status: claimed
status_changed_at: 2026-08-18T13:05:12Z
created_at: 2026-08-18T12:09:46Z
user_request: UR-051
addendum_to: REQ-237
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T13:05:12Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - browser evidence
claimed_at: 2026-08-18T13:05:12Z
route: B
---

# Stop Panel B's Slowest-Day Annotation Colliding With Its Own Title

## What

In the Durations view, Panel B's slowest-day annotation is drawn at `y = 355` while Panel B's own title sits at `y = 350` (`DURATIONS_MEDIAN_TITLE_Y`). The two overlap: on a synthetic fixture the annotation `209 min` renders directly through the words "paused and broken spans excluded".

## Context

Found while reviewing REQ-237 by rendering a dense fixture and looking at it. **It is not a REQ-237 regression** — the same annotation sits at the identical `x = 357.2, y = 355.0` on a board built from the pre-REQ-237 binary, checked side by side. It is pre-existing and was simply never looked at on a fixture whose slowest day lands under the title text.

It is invisible on this repository's own board because the annotation's x-position depends on which day is slowest, and here that day falls clear of the title's width. That is luck, not design — which is why it wants pinning rather than nudging.

The annotation reuses the `durations-mark-label` class, so it is not part of either label band's row packing and is not covered by REQ-231's mark-band geometry test or by REQ-237's row-fill test. Nothing in the suite looks at it at all.

## Requirements

- Panel B's slowest-day annotation does not overlap Panel B's title at any x-position the annotation can take, including when the slowest day is the leftmost one.
- The annotation stays associated with the bar it describes — moving it somewhere it no longer reads as belonging to that day is not a fix.
- No change to Panel A, Panel C, or the label bands; no change to `describeAtPointer`'s panel boundary.
- A test pins the separation, so the next person to move `DURATIONS_MEDIAN_TITLE_Y` finds out.

## Red-Green Proof

**RED prompt/case:** an assertion that the slowest-day annotation's text box and Panel B's title text box do not intersect, read from the renderer's own constants the way `TestDurationsLabelRowsClearTheMarkBands` reads them — evaluated at the annotation's worst-case x, not at whichever x this repository's data happens to produce.
**Why RED now:** the title's baseline is 350 and the annotation's is 355, so their boxes intersect wherever their x-ranges do; reproduced on a fixture as `209 min` drawn through the title.
**GREEN when:** the assertion passes and a rendered fixture whose slowest day sits under the title shows the two clear of each other.
**Validation:** Review finding on REQ-237; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] While reviewing REQ-237 I rendered the Durations chart on a test board and found the little "209 min" note that marks the slowest day printed straight through the heading above it — the note sits five units below a heading that is taller than five units. It has been like this since before any of today's work; it does not show on your own board only because the slowest day happens to fall to the right of where the heading's text ends, which is chance rather than design. The fix is small — move the note, or move the heading, and add an assertion so the next person who shifts either one finds out. I am asking rather than doing it because "move which one, and where" is a look-and-feel choice about a chart you read regularly, and there is more than one reasonable answer. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — it only shows on contrived data and the note is legible enough where it is.
  → Confirmed: Yes, add to queue (builder picks placement, pinned by a test). [2026-08-18, via do-work clarify]

---

# Builder Guardrails (orchestrator-issued — binding)

## Your tree

- Work **only** inside your worktree (path below). It is a full checkout on your own branch.
- **Never write anywhere in the main tree** except the single hand-back file named below. That is the only main-tree path you may touch.
- **Never touch `do-work/`** — not the queue, not `working/`, not `CHECKPOINT.md`, not `archive/`. Queue state is the orchestrator's alone. Your branch must contain **zero** commits touching `do-work/`; the orchestrator runs `git diff --name-only <pre>...<your-branch> -- do-work/` and a single path there stops your hand-back.
- **Never touch `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, or `CHANGELOG.md`.** Those are serial-only integrator-owned files. A bump on your branch races every sibling.
- **Scratch files go in `/tmp` or inside your worktree — never the main tree root.** A previous builder left a PNG in the repo root; that is a write-set violation. Screenshots, fixtures, generated boards: `/tmp`.

## Commit on your branch

Commit your implementation on your own branch before handing back. Message body only — no version bump, no changelog entry. End the message with:

```
Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
```

## The P-A-U loop is yours to fill

The REQ body contains an `## AI Execution State (P-A-U Loop)` section with three checkboxes, or the orchestrator will add one. **You must tick all three and write the required content into each**, in your worktree's copy of nothing — instead, put the filled P-A-U block **into your hand-back file** verbatim under a `## P-A-U` heading, since you may not write `do-work/`. `qualify.sh` FAILs on unticked boxes and the orchestrator will otherwise have to fill them from your evidence.

- **[PLAN]** — read the listed `prime_files` and agent rules, then write the technical approach. No code yet.
- **[APPLY]** — code written exactly as planned, scope strictly limited to planned files.
- **[UNIFY]** — run `git diff --stat`, review every changed file, run the project's linters/tests, confirm no debug artifacts. List each file you verified and what you checked.

## Evidence rules — every one of these was learned by getting it wrong

1. **Two REDs when the first is a reference error.** A test that fails because a constant or function does not exist yet proves nothing. Put the code in place, break exactly one rule, and let the assertion fail *for the reason it exists*. Report both RED outputs.
2. **`git stash push` on a clean file stashes nothing** — and the resulting green run reads as proof when it is vacuous. To reproduce RED against pre-change code, check out the pre-change blob by hash (`git show <hash>:<path>`) instead.
3. **Assert page identity inside the same call that reads the DOM.** If you drive a browser, return `location.href` (and, where relevant, the page's own rule text) from the *same* `evaluate` call as every measurement. A shared browser instance can be navigated by a sibling between your navigate and your evaluate, and the numbers come back confident, well-formed, and about somebody else's page. A URL checked *before* navigating is not the same claim. Prefer an isolated browser context.
4. **A programmatic `.focus()` does not trigger `:focus-visible` in Chrome.** Use a real `Tab` keypress if focus styling is in question.
5. **Generate the artifact and look at it.** For anything that changes what appears on screen, a passing assertion is not evidence about two glyphs sharing a coordinate. Measure `getBoundingClientRect()` intersections in the live DOM when the question is "do two things overlap"; read the rendered text when the question is "what does this say".
6. **Push back if the brief is wrong.** If a requirement contradicts an existing test, or a piece of code you wrote turns out unneeded, say so in the hand-back rather than quietly editing the test or keeping dead code. Two builders pushed back last session and both were right.

## Verification bar

`bash _dev/tests/maintainer-verify.sh` from your worktree root. **Exit code 0 is the only proof.** Never pipe it through `tail`/`head` — the pipeline's exit status hides the failure. Run it, then `echo $?` on its own line, and paste that.

## Hand-back

Write **one** file, at the absolute path given below, containing:

1. `## Branch` — your branch name.
2. `## P-A-U` — the three filled, ticked checkboxes with their content.
3. `## Files Changed` — `git diff --stat` against your branch's merge base, plus one line per file saying what changed and why.
4. `## Red-Green Evidence` — the RED output(s) and the GREEN output, quoted.
5. `## Verification` — the `maintainer-verify.sh` tail and its `echo $?` line.
6. `## Integration Seams` — anything the orchestrator must apply by hand in the merge commit (shared registries, cross-REQ text). Say "none" if none.
7. `## Decisions` — numbered D-01, D-02… for choices with reach beyond this REQ.
8. `## Lessons Learned` — what a future session should know. Omit if genuinely nothing.
9. `## Pushback` — anything in this brief you think is wrong. Omit if none.

Your final message back should be a short summary; the hand-back file is the real deliverable.

## Your Assignment

- **Worktree path (your working directory):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-242-stop-panel-b-annotation-colliding-with-its-title`
- **Branch name:** `worktree-agent-REQ-242-stop-panel-b-annotation-colliding-with-its-title`
- **Hand-back file (absolute, main tree — the ONE main-tree path you may write):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-18-124358/REQ-242-handback.md`
- **Repo root of the MAIN tree (read-only for you):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`

## Orchestrator Notes for This REQ — READ THESE, THE GROUND MOVED UNDER THIS REQ

**REQ-241 landed first, in this same session, and it moved the row pitch.** Your worktree already contains it. What that means for you:

- `DURATIONS_LABEL_ROW_HEIGHT` is now **13**, not 12, and `durationsLabelCharacterWidthUnits` is now **6.75**, not 6.2. Any coordinate arithmetic in the REQ body predating that is stale — **re-measure, do not trust the REQ's numbers.**
- **The collision still reproduces.** REQ-241's builder rendered dense fixtures before and after its change and saw `233 min` overprinting `…spans excluded` in both. So this REQ is still live; it was not resolved by REQ-241. Confirm that yourself first, on your own fixture, before building anything — if it does not reproduce for you, stop and say so rather than fixing an absent defect.
- **REQ-241 D-03 is a hard constraint you inherit.** Panel B's title (`DURATIONS_MEDIAN_TITLE_Y = 350`) is the sole input to `describeAtPointer`'s A/B boundary (`pointerY <= DURATIONS_MEDIAN_TITLE_Y - 12`). REQ-241 deliberately did **not** move it, because moving it moves the pointer boundary. It also measured that the reversed band's last row now clears the title by only **1.364 units** — that is the entire remaining headroom above the title, and a row pitch of 15 would consume it.

  **Therefore: strongly prefer moving the annotation, not the title.** Moving `DURATIONS_MEDIAN_TITLE_Y` costs you the `describeAtPointer` requirement and eats REQ-241's headroom. If you conclude the title must move anyway, that is allowed — but you must then prove `describeAtPointer` resolves the same panel for the same pointer position at probe positions spanning the boundary, and say why the annotation could not move instead.

- Your `write_set` is `web/board-durations.js` and `generate_test.go`. **You own `generate_test.go` this wave** — no sibling will touch it. Do not touch `durations.go` or `durations_test.go` (REQ-241's files, already merged) unless you have a reason you state.
- The annotation reuses the `durations-mark-label` class and is outside both label bands' row packing, so REQ-231's mark-band test and REQ-237's row-fill test do not cover it. Nothing in the suite looks at it. That is why it wants a test, not just a nudge.
- The test must be evaluated at the annotation's **worst-case x** — including when the slowest day is the leftmost one — not at whichever x this repository's data happens to produce. That is the whole reason the defect was invisible here.
- **Generate a board and look at it.** Use a fixture whose slowest day falls under the title. Measure `getBoundingClientRect()` intersections in the live DOM; return `location.href` from the same evaluate call. Use distinct scratch paths (`/tmp/qk-242`, `/tmp/board-242`).
