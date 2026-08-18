---
id: REQ-241
title: Reconcile the Durations label metrics with the face actually rendered
status: claimed
status_changed_at: 2026-08-18T12:43:06Z
created_at: 2026-08-18T12:09:46Z
user_request: UR-051
addendum_to: REQ-237
domain: general
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: durations-label-metric-constants-disagree-with-rendered-face
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/durations.go
- skills/do-work-board/tools/queue-kanban/durations_test.go
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-18T12:43:06Z
  basis:
    - Route B
    - 3-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - browser evidence
    - cross-route regression gates
claimed_at: 2026-08-18T12:43:06Z
route: B
---

# Reconcile the Durations Label Metrics With the Face Actually Rendered

## What

Two constants that describe the Durations label face disagree with what the browser actually draws. Neither causes a visible collision today, but both are now load-bearing in a way they were not before, because REQ-237 made the label rows actually fill up.

## Context

Found by REQ-237's build and confirmed independently against the merged tree. Both are pre-existing; REQ-237 is what made them reachable, not what caused them.

Until REQ-237, the overflow lane drew two or three labels on any real board, so the slack in these constants never mattered. It now draws as many as fit — measured 21 on a clustered 60-sample fixture — and the slack is what stands between "packed" and "overlapping".

## Instances

- [ ] **`durationsLabelCharacterWidthUnits = 6.2` under-estimates the 11px sans face by ~7%.** Its comment calls the value "deliberately generous"; measurement says otherwise — a 14-character label renders **92.52 user units**, i.e. **6.61 units/char**, not the 86.8 the constant predicts. The 6-unit separation rule absorbs the difference, so nothing collides: the tightest same-row gap in a full 21-label render is 3.08 units, and a direct DOM measurement finds **0 same-row overlaps**. But the real margin is about half what the code claims, and the comment is actively misleading about which direction the error runs.

- [ ] **`DURATIONS_LABEL_ROW_HEIGHT = 12` is smaller than the text box the same file declares.** `DURATIONS_LABEL_TEXT_ASCENT = 11` plus a 2-unit descent is a 13-unit box on a 12-unit pitch; the rendered font box measures 12.83 units. Measured on a densely-populated lane: **20 cross-row bounding-box intersections, each 1.6px deep.** This is line-box padding rather than ink — the render shows two cleanly separated rows, and a screenshot confirms it — but it means no test can honestly assert row-against-row separation the way `TestDurationsLabelRowsClearTheMarkBands` asserts row-against-mark.

## Requirements

- Each constant either matches the face the browser renders, or its comment states the measured value and why the code deliberately differs. A constant whose comment claims a safety margin in the wrong direction is the specific defect here.
- Whatever changes, the same-row separation guarantee holds: **0 same-row label overlaps** at full density, measured from the live DOM rather than computed from the constants under test.
- REQ-231's guarantee holds unchanged: **0 label/mark overlaps** in either band, at any density.
- If the row pitch changes, Panels B and C shift with it and `describeAtPointer` still resolves the same panel for the same pointer position — the same constraint REQ-231 worked under.
- Changing label counts across the view is an accepted consequence of retuning the width model; say so in the REQ trail with before/after counts on a real board rather than only on a fixture.

## Red-Green Proof

**RED prompt/case:** a test asserting each constant against the measured face — that `durationsLabelCharacterWidthUnits` is not below the rendered units-per-character, and that `DURATIONS_LABEL_ROW_HEIGHT` is not below the declared ascent-plus-descent box.
**Why RED now:** 6.2 < 6.61, and 12 < 13.
**GREEN when:** the assertions pass, and a rendered clustered fixture still shows 0 same-row label overlaps and 0 label/mark overlaps measured in the live DOM.
**Validation:** Review finding on REQ-237; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] While reviewing REQ-237 I found two numbers in the Durations chart that describe the label text and are both slightly wrong about it. One says each character is 6.2 units wide when the browser draws 6.61 — and its comment claims the estimate is deliberately generous, which is backwards. The other says a row of label text is 12 units tall when the same file elsewhere describes that text as a 13-unit box. Nothing overlaps on screen today: I measured a fully packed lane and found no two labels on the same row touching, and no label touching a dot. But the spare room absorbing both errors is about half what the code claims it is, and until this session the chart only ever drew two or three labels so it never mattered. It does now — REQ-237 made the rows fill up. Fixing the width number will change how many labels the chart draws on every board, including yours, which is why it is your call and not a quiet tidy-up: it is a visible change to a view you look at, made for correctness reasons rather than because anything is broken. The alternative is to leave the numbers and fix the comments so they stop claiming a margin that is not there. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No — just correct the misleading comments and leave the numbers alone, since nothing overlaps.
  → Confirmed: Yes, add to queue — full fix, retune the constants to match the measured face (user accepted that label counts across the view will visibly change). [2026-08-18, via do-work clarify]

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

- **Worktree path (your working directory):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-241-reconcile-durations-label-metrics-with-the-rendered-face`
- **Branch name:** `worktree-agent-REQ-241-reconcile-durations-label-metrics-with-the-rendered-face`
- **Hand-back file (absolute, main tree — the ONE main-tree path you may write):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-18-124358/REQ-241-handback.md`
- **Repo root of the MAIN tree (read-only for you):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`

## Orchestrator Notes for This REQ

- Your `write_set` is exactly three files. Do **not** touch `generate_test.go` — a sibling builder owns it this wave.
- This REQ **will visibly change the maintainer's own board.** The requirement says record before/after label counts **on a real board** (`--repo-root` pointed at the main checkout's `do-work/`), not only on a synthetic fixture. Do both, and put both numbers in the hand-back.
- Build and render like this (from your worktree):
  ```
  cd skills/do-work-board/tools/queue-kanban && go build -o /tmp/qk-241 .
  /tmp/qk-241 generate --repo-root <your worktree root> --out /tmp/board-241
  ```
  Use a distinct binary and output path (`-241`) so you cannot read a sibling's artifact.
- The two constants may move in opposite directions from what the REQ text assumes. If retuning `DURATIONS_LABEL_ROW_HEIGHT` upward costs a label row, say so and state the trade rather than silently absorbing it.
- REQ-231's guarantee (0 label/mark overlaps) and the same-row guarantee (0 same-row label overlaps at full density) must both be **measured in the live DOM**, not computed from the constants under test. Computing them from the constants you just changed is circular.
- `describeAtPointer` must still resolve the same panel for the same pointer position if the row pitch changes. Check it.
