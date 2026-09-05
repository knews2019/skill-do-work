---
id: REQ-578
title: 'Hide the verify-findings strip on the Activity view'
status: claimed
created_at: 2026-09-04T23:58:59Z
user_request: UR-117
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-05T12:16:57Z
  basis:
    - trivial short-circuit
related: [REQ-573]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
claimed_at: 2026-09-05T12:00:56Z
dispatch_at: 2026-09-05T12:06:43Z
route: A
builder_handback_at: 2026-09-05T12:16:29Z
integration_at: 2026-09-05T12:16:29Z
review_at: 2026-09-05T12:24:37Z
commit: 09aaa9a443f8bb6191b162393403e60b4f8fa6f4
---

# Hide the Verify-Findings Strip on the Activity View

## What

The Verify Findings strip (`#board-findings`, added by REQ-285) sits outside the view panels so it stays visible on every view. On the Activity view it pushes the transitions table down and is not what that view is for. Hide the strip while the Activity view is active and show it again when the reader switches to any other view. The strip's content and its behavior on the other views do not change.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Builder read the kanban prime, both lesson satellites and the crew rules, then chose to put the rule in the view switcher beside its other visibility toggles and to read the strip's emptiness back off the renderer's own output. Three sources for that second half were weighed and two rejected. Recorded under `## P-A-U` in `do-work/runs/work-2026-09-05-120117/REQ-578-handback.md`.
- [x] **[APPLY]:** One commit on the builder branch (`d52af05f`) touching exactly the three Scope files; the rule is inline rather than in a helper so the Node lane's function slicing can drive it.
- [x] **[UNIFY]:** `git diff --stat` reviewed (3 files, +158/-5); `gofmt -l .` empty; `go vet ./...` exit 0; `node --check web/board-controls.js` exit 0; debug-artifact scan over added lines empty; each of the three files read line by line, with the checks per file listed in the hand-back.

## Detailed Requirements

- With the Activity view selected, `#board-findings` is hidden (the `hidden` attribute, matching how the strip already hides itself when there are no findings) even when findings exist.
- Switching to Board, Calendar, Durations, Timeline or Testing shows the strip again exactly as today; the "probe(s) could not run" disclosure under it follows the strip.
- Only the Verify Findings strip is affected. The completion-anomalies strip above it is not part of this request.
- The rule lives in the view-switching code (`board-controls.js`), not in the Activity renderer, so a re-render of the Activity table never touches the strip. Update the template comment that says the strip stays visible in every view.

## Red-Green Proof
**RED prompt/case:** In the Node behavior lane, render with two verify findings in `boardData`, switch the view to `activity`, and read `document.getElementById("board-findings").hidden`.
**Why RED now:** The strip is outside the view panels by design and nothing in the view switch touches it, so it stays visible on the Activity view (screenshot 3).
**GREEN when:** `hidden` is true while the Activity view is active and false again after switching back to the Board view with the same findings.
**Validation:** User request from the live board; proof inferred during capture.

## Builder Guidance

The user is certain about the outcome. Keep it to the view switch plus one test; do not restructure the strips.

## Assets

- `do-work/user-requests/UR-117/assets/REQ-578-screenshot-3-activity-view-with-verify-strip.png`: the Activity view at 24h with "175 transitions across 38 REQs in the last 24 hours" and rows for REQ-576, 575, 574, 572 (four rows: work merged, builder handed back, builder dispatched, captured), 506, 570, 573, 505 and others; above the table a Verify Findings strip with two cards (WORKTREE-MERGE-STATE-UNDETERMINED for a REQ-506 worktree, WORKTREE-PRESENT-RUN-IN-FLIGHT for the REQ-570 worktree) and a "1 probe(s) could not run" disclosure.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

*Source: "remove verify finding from this view"*

---

## Triage

**Route: A** - Simple

**Reasoning:** One view-switch rule in a named file (`board-controls.js`), one template comment, one Node behavior test, all three declared in the write set, with a captured RED/GREEN pair. No discovery needed.

**Planning:** Not required

## Plan

**Planning not required** - Route A: direct to builder

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)

The view rule lives in the controls file, the strip's placement comment moved to name the Activity exception and its owner, the new Node-lane behavior test is `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip`, and the cards file carries the builder's integration seam, applied by the orchestrator inside the merge commit.

**What was done:** `applyView` now sets `#board-findings.hidden` to "the Activity view is active, or the strip has nothing to say", beside the other `hidden =` toggles it already owns. Emptiness is read back off the renderer's own output (`#board-findings-cards` children and `#board-findings-skipped-list` children) rather than re-read from `boardData`, so the rule that decides whether the strip has anything to say keeps exactly one definition. Merge range `7dbb2756..09aaa9a4`; builder branch head `d52af05f`. Builder-authored `## Decisions` (D-01 to D-04) and `## Discovered Tasks` live in `do-work/runs/work-2026-09-05-120117/REQ-578-handback.md`.

## Qualification

**Passed.** `qualify` returned `OK: mechanical qualification passed` over the merge range `7dbb2756..09aaa9a4` after two orchestrator bookkeeping fixes: the three P-A-U boxes were ticked from the hand-back (the builder cannot write the REQ file) and the file list was reduced to bare paths, because the checker reads backticked words in that list as paths.

Read from the diff, not from the summary:
- `board-controls.js` — the new block sits with `applyView`'s other `hidden =` toggles and sets the strip hidden when the active view is `activity` or when the renderer's own output is empty. Emptiness is read from `#board-findings-cards` children and `#board-findings-skipped-list` children, so the "does the strip have anything to say" rule keeps one definition and cannot drift from the renderer. No new module state.
- A skipped probe still shows: the empty test covers a strip with zero findings but a live "probe(s) could not run" disclosure, which is the case a naive children-count on the cards host alone would wrongly hide. A skipped probe reading as clean is the failure this pins.
- `template.html` and `board-cards.js` carry comment changes only; the `board-cards.js` line was the builder's integration seam, applied inside the merge commit so it sits within this REQ's range.
- Only the Verify Findings strip changed. The completion-anomalies strip above it is untouched, and its "stays visible in every view" comment in `board.css` is still true.

Requirements traced: hidden on Activity even with findings present (the `activity` arm), shown again on Board/Calendar/Durations/Timeline/Testing (the same expression on every other view, covered by the round-trip cases), the disclosure follows the strip (it is the strip's own child), the anomalies strip untouched, and the rule living in the view-switching code rather than in the Activity renderer so a re-render never touches the strip.

*Checked by work action*

## Testing

**Focused tests (post-merge, main tree at `09aaa9a4`):**
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...` — exit 0, 57 tests, wall 15s, slowest file `citations_test.go` at 4.36s against the 30s per-file budget.

**Builder lanes (in its worktree at `d52af05f`):** whole Node behavior lane exit 0 (57 tests, wall 7s); whole module suite exit 0 (392 tests, wall 49s, slowest file `generate_test.go` at 10.20s).

**Red-green validation** (traced to `## Red-Green Proof`): RED with the test in place and the `board-controls.js` change absent — exit 1, `javascript_behavior_c_test.go:2641: the Verify Findings strip is still visible on the Activity view`, with every other case in the same probe passing, so the failure is the assertion and not a load error. GREEN after the change — exit 0. One intermediate RED is recorded in the hand-back: an implementation that read the skipped disclosure's `hidden` attribute made a strip with nothing to report become visible on a view switch; reading the skipped list's children instead removed the dependency on a template-initial attribute so the probe and the browser read the same thing.

**Not covered:** the rendered pixels. This is a single boolean attribute the browser honours natively and the probe drives the real `applyView`, so no render pass was made.

## Review

**Overall: 94%** | 2026-09-05T12:22:16Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Restatement Sweep: `do-work/queue/REQ-579-render-verify-findings-and-skipped-probes-as-compact-rows.md:45` (the queued REQ that rebuilds this strip as compact rows) states "REQ-578 keeps working against the same `#board-findings` element and `hidden` attribute; do not restructure the strip in a way that breaks it". That is now an incomplete statement of this REQ's contract: `applyView` (`web/board-controls.js:55-59`) also hard-depends on `#board-findings-cards` and `#board-findings-skipped-list` and reads `.children.length` on both. The same REQ's D1 explicitly plans to "delete the `<details>` disclosure", which removes `#board-findings-skipped-list`; in a browser that makes `applyView` throw a TypeError on `null.children` and kills every view switch, not just this strip. The Node behavior lane will not catch it (see the Minor finding on the stub DOM). Stale restatement in a file this REQ never declared, reported per Step 6.3 — `impact-rule-change` → report only

**Minor findings:**
- `renderVerifyFindingsStrip` (`web/board-cards.js:655-658` and `:682-685`) returns early on its two empty paths **without** clearing `#board-findings-cards` / `#board-findings-skipped-list`, so D-01's rationale ("cannot go stale if the renderer is ever re-run") does not hold. Demonstrated in a sandbox copy of the module with a probe built from this REQ's own test harness: render with one finding + one skipped probe, then re-render with an empty payload → the renderer hides the strip (`hidden=true`) but both hosts keep their stale children (1 and 1), and the next `applyView` sets `hidden=false`, putting a stale finding card back on screen. Latent today, and only today: `renderVerifyFindingsStrip` is called exactly once, at boot (`web/board.js:76`), and new data arrives only through a full page reload. The one-line fix is to move `cardsHost.textContent = ""` / `skippedList.textContent = ""` above the early returns. Named here because REQ-579 rewrites this renderer — `impact-negligible` → report only
- The Node probe's stub `document.getElementById` (`javascript_behavior_c_test.go:2555-2559`) creates a fresh node for *any* id it is asked for, so a renamed or deleted element reads as "present with zero children" instead of failing. This test therefore cannot prove that the three ids `applyView` now depends on exist in `template.html`. Verified by hand instead: `#board-findings`, `#board-findings-cards` and `#board-findings-skipped-list` are all present (`web/template.html:193`, `:201`, `:209`). The blind spot is a property of the whole behavior lane, not of this test alone — `impact-rule-change` → report only

**Nit findings:**
- `web/board.css:533-536` comment says "the strip stays visible in every view" directly above the `.board-anomalies` rules. True for the completion-anomalies strip it names, but the findings section carries `class="board-anomalies board-findings"`, so the sentence sits above styling that the Activity rule now contradicts for its second user. The builder recorded this deliberately as correct-and-do-not-fix; recorded here with the same verdict plus the class-sharing caveat — `impact-negligible` → report only
- `skills/do-work-board/docs/board-guide.md:20` says "Two strips sit above the columns and stay visible in every view" and lists Notes and Completion anomalies only. The Verify Findings strip has been missing from the shipped guide since REQ-285 (the REQ that added it), so the sentence is not stale in itself — but whoever finally documents the strip must not extend that "every view" clause to it. Pre-existing, not caused by this diff. That file is currently modified by another live session, so nothing was touched — `impact-negligible` → report only

**Acceptance:** Pass — RED reproduced independently (module copied to a scratch tree, the `applyView` block removed, `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip` failed on exactly `javascript_behavior_c_test.go:2641: the Verify Findings strip is still visible on the Activity view` and on no other assertion), GREEN confirmed in the main checkout (`QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 run-go-tests-with-budget.sh … -run '^TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip$'` exit 0, wall 22s, inside the 30s per-file budget). Boot order checked: `renderVerifyFindingsStrip()` runs at `web/board.js:76` before `applyView()` at `:79`, so the first paint reads populated hosts, not empty ones. D-02 (no null guard) is sound and its reasoning is stronger than stated: the client JS is inlined into `index.html` at generate time (`generate.go:19`, `:442`), and only the code-free `board-data.js` is a separate file, so the mixed-version tree the renderer's own guard was written for cannot pair an old template with this new `applyView`. A skipped-probes-only strip behaves correctly on every non-Activity view (`board-findings-skipped-list` children keep the strip visible), and on Activity the whole strip is hidden by the user's explicit request rather than a skipped probe being rendered as clean.

**Suggested testing:** 4 items — (1) open a live board with real verify findings, switch Board → Activity → Board and confirm the transitions table starts at the top and the strip returns with its cards intact; (2) the same round trip on a board whose only content is "N probe(s) could not run", confirming the disclosure comes back; (3) a board with zero findings, confirming no view switch makes an empty strip appear; (4) when REQ-579 lands, re-run this test *in a browser*, not only in the Node lane, because the lane cannot see a removed element id.

**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Reading the strip's emptiness back off the renderer's own DOM output instead of re-reading `boardData` kept the "is there anything to say" rule in one place, and the review confirmed the no-null-guard choice was safe for a stronger reason than the builder gave: the client JavaScript is inlined into `index.html` at generate time, so the mixed-version tree the renderer's own guard was written for cannot occur.

**What didn't:** The first implementation read the skipped disclosure's `hidden` attribute and made an empty strip visible on a view switch, because the stub node in the test lane carries no template-initial attributes. Reading the skipped list's children removed the dependency on any initial attribute and made the probe and the browser read the same thing.

**Worth knowing:** The Node behavior lane's stub `document.getElementById` manufactures a fresh node for any id it is asked for, so a renamed or deleted element reads as "present with zero children" rather than failing. Any rule that depends on an element existing needs a separate assertion against the generated page text, or it is untested. That matters immediately: REQ-579 rewrites this same strip and plans to delete the disclosure that holds `#board-findings-skipped-list`, which `applyView` now dereferences.

## Orientation

The board's Verify Findings strip now steps out of the way on the Activity view and comes back everywhere else. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`), in the view switcher rather than in the Activity renderer, so re-rendering the transitions table never touches the strip. No prime was made stale.

## Heavy Verification Plan

- **Base revision:** 7dbb27562f58b5ede067a478453edd7fbe70c3c8
- **Target revision:** 09aaa9a443f8bb6191b162393403e60b4f8fa6f4
- **Planned at:** 2026-09-05T12:24:37Z, from `_dev/tests/heavy-lanes.json`

| Lane | Argv | Why it was selected |
| --- | --- | --- |
| `queue-kanban-javascript` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` | every changed path matched subtree `skills/do-work-board/tools/queue-kanban` |
| `queue-kanban-browser` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` | same subtree match |
| `staged-skills` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` | same subtree match |

No path was left uncovered by the manifest. The browser lane is where the reviewer's suggested manual checks land: a Board → Activity → Board round trip with real findings, the same round trip with only skipped probes, and a zero-finding board where no view switch may make an empty strip appear. The request stays `claimed` with its `commit:` landed until the queue-exhaustion drain.

