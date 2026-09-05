# Hand-back — REQ-578 (hide the Verify Findings strip on the Activity view)

## Branch

- Branch: `worktree-agent-REQ-578-findings-strip`
- Head commit: `d52af05f` — `[REQ-578] Hide the verify-findings strip on the Activity view`
- Base: `09a13839` (main)

## File manifest

- Modified: `skills/do-work-board/tools/queue-kanban/web/board-controls.js` — `applyView` now sets `#board-findings.hidden` (11 added lines).
- Modified: `skills/do-work-board/tools/queue-kanban/web/template.html` — the strip's placement comment now says the Activity view is the exception and names `applyView` as the owner (13 changed lines, net +8).
- Modified: `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` — added `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip` (139 added lines).

No file outside the REQ's `write_set` was touched. `git diff --stat` on the branch: 3 files, 158 insertions, 5 deletions.

## P-A-U

**[PLAN]**

The strip is rendered once at boot by `renderVerifyFindingsStrip` (`web/board-cards.js`), which owns one rule: hidden when there are no findings and no skipped probes. `applyView` (`web/board-controls.js`) owns every other visibility toggle in the topbar and the view panels, and runs on every view switch. So the Activity rule goes in `applyView`, beside the existing `hidden =` toggles, as `hidden = view is activity || the strip has nothing to say`.

The second half needed a source. Three candidates:

1. re-read `boardData.verifyFindings` / `verifySkipped` lengths — a second copy of the renderer's emptiness rule;
2. cache what the renderer decided in a module-level variable — correct only while the renderer runs exactly once, and unreachable from the Node probe lane, which can slice functions out of the built page but not module-scope `var`s;
3. read back the DOM the renderer drew.

Picked (3): the cards host and the skipped list are the renderer's own output, so the emptiness rule keeps exactly one definition and cannot go stale if the renderer is ever re-run.

Test plan: one Node behavior probe slicing `createElement`, `renderVerifyFindingsStrip` and `applyView` out of the generated page, driving them over stub DOM nodes. Four cases — findings present on Board, on Activity, back on Board and on Timeline; nothing to report on Board and Activity; skipped probes only across an Activity round trip.

**[APPLY]**

Written exactly as planned, inside the three declared files. The rule is inline in `applyView` rather than in a helper function, on purpose: the Node lane slices whole named functions out of the built page, so a helper would have made the RED run fail with a `ReferenceError` instead of a real assertion.

**[UNIFY]**

- `git diff --stat` → 3 files changed, 158 insertions(+), 5 deletions(-). No unexpected path.
- `gofmt -l .` in the module → no output (clean).
- `go vet ./...` in the module → exit 0.
- `node --check web/board-controls.js` → exit 0.
- Added lines scanned for `console.log`, `debugger`, `TODO`, `FIXME` → none.
- `web/board-controls.js` — read the whole diff hunk: the block sits with the other `hidden =` toggles, uses two-word names (`findingsStrip`, `findingsStripHasContent`), and adds no state. No null guard, matching the rest of `applyView`, which reads every element id straight (see D-02).
- `web/template.html` — read the whole comment: the "survives every view switch" claim is gone, the Activity exception and its owner (`applyView`) are named, and the rest of the comment (filters, recently-done window, Go-side suppression) is unchanged.
- `javascript_behavior_c_test.go` — read the added test: stub node factory mirrors the one in the neighbouring Activity test, every assertion carries a failure message naming the real defect, no assertion is vacuous (each was observed failing in at least one run).

## Test evidence

RED (before the `board-controls.js` change, test only):

```
QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 \
  bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban \
  -run '^TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip$' ./...
```

exit 1 —
`javascript_behavior_c_test.go:2641: the Verify Findings strip is still visible on the Activity view`
`--- FAIL: TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip (1.83s)`

The failure is the assertion itself, not a load error: every other case in the same probe (strip visible on Board, empty strip hidden) passed in that run.

Intermediate RED worth recording: the first implementation read `!document.getElementById("board-findings-skipped").hidden` and the probe failed with
`a strip with nothing to report became visible on a view switch (afterRender=true, board=false, activity=true)`.
That branch depended on the skipped disclosure's initial `hidden` attribute from the template, which the stub node did not carry. Reading the skipped **list's children** instead removed the dependency on any initial attribute and made both the probe and the browser read the same thing — the renderer's output.

GREEN (same command, after the change): exit 0 —
`go-test budget: module=skills/do-work-board/tools/queue-kanban wall=3s tests=1 slowest-file=javascript_behavior_c_test.go:1.95s limit=<30s`

Whole Node behavior lane:

```
QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 \
  bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban \
  -run '^TestJavaScriptBehavior' ./...
```

exit 0 — `wall=7s tests=57 slowest-file=javascript_behavior_c_test.go:2.30s limit=<30s`.

Whole module suite:

```
bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...
```

exit 0 — `wall=49s tests=392 slowest-file=generate_test.go:10.20s limit=<30s`. Per-file budget respected; the 49s is the module wall, not one file.

Untested: the rendered pixels. This is a visibility change, and the Verify Findings strip needs a live board with real verify findings to look at. The prime's "generate a board and look at it" rule is about geometry and overlap; here the assertion is a single boolean attribute the browser honours natively, and the probe drives the real `applyView`. The orchestrator may still want a render pass with findings present.

## Lesson evidence

The REQ's `required_lessons` were dropped for budget, but both were read in full (they are the touch-conditional set for this area, and `general.md` makes that additive):

- `_dev/primes/prime-kanban-board.md` — read in full (the REQ's only `prime_files` entry).
- `_dev/primes/lessons-kanban-board.md` — read in full. Directly relevant: REQ-232 ("describe the board by what its switcher covers, not by naming some of its tabs") pushed the template comment toward naming the rule's owner rather than re-listing views; REQ-235 ("derive state instead of storing it") is why the module-level cache option was dropped.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — read in full. The 0.275.3 entry ("when a status changes owner, walk every classifier that names it") is the same shape as the stale comment in `board-cards.js` reported under Integration seams below.
- No listed path was missing.

## Decisions

- **D-01 — DECIDE & STATE. Read the strip's emptiness back off the rendered DOM (`#board-findings-cards` children, `#board-findings-skipped-list` children) rather than re-reading `boardData`.** Both hosts are `renderVerifyFindingsStrip`'s own output, so the "is there anything to say" rule keeps one definition. Re-reading `boardData` in `applyView` would be a second copy that drifts the first time the producer's payload shape changes. Reversible in one line.
- **D-02 — DECIDE & STATE. No `if (findingsStrip)` guard, unlike the renderer's.** `applyView` reads eleven element ids with no null checks; the strip ships in the same `template.html` as those elements, so a tree where the strip is missing is a tree where `applyView` was already broken. The renderer's guard exists for a different caller shape (an older static `index.html` beside a newer `board-data.js` payload), which cannot happen for the embedded template. Matching the surrounding style, per the surgical-changes rule.
- **D-03 — DECIDE & STATE. The rule is inline in `applyView`, not extracted into a named helper.** The Node behavior lane slices whole functions out of the generated page; a helper would not have been sliced, so the RED run would have died on a `ReferenceError` rather than the assertion — not RED-before-GREEN evidence. Eleven lines with a comment is also below the threshold where extraction pays.
- **D-04 — DECIDE & STATE. The test covers a skipped-probes-only strip as its own case.** A finding count of zero with skipped probes present is the exact state where the naive `children.length` check on the cards host alone would wrongly hide a strip that says "1 probe(s) could not run — unverified, not clean". The repo's standing rule is that a skipped probe must never read as clean, so the case is pinned.

## Discovered Tasks

- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` line 640 comment now states a half-truth: "it lives outside the view panels, ignores the recently-done window, and ignores the shared filters" reads as an absolute view-switch exemption, which the Activity rule breaks. The file is outside this REQ's `write_set`; the exact replacement line is under Integration seams. → report only (impact-cosmetic).
- `skills/do-work-board/tools/queue-kanban/web/board.css` line 535 carries the same claim for the **completion anomalies** strip ("the strip stays visible in every view"). That one is still true — REQ-578 leaves the anomalies strip alone — so it needs no change. Recorded only so a later sweep for the phrase does not "fix" a correct comment. → report only (impact-cosmetic).
- The board's changelog entry for REQ-285 (`CHANGELOG.md` and `skills/do-work/CHANGELOG.md`, both line 1905) says "no view switch … can hide a finding". Changelog history is not rewritten, so nothing to do; the release entry for this REQ should state the new exception plainly. → report only (impact-cosmetic).

## Integration seams

One line, in a file outside the write set. `skills/do-work-board/tools/queue-kanban/web/board-cards.js`, in the `renderVerifyFindingsStrip` header comment (currently lines 639-643):

Replace:

```
  // Same exemptions as the anomalies strip above, for the same reason: it lives
  // outside the view panels, ignores the recently-done window, and ignores the
  // shared filters, because a finding must not be hideable by a filter
  // combination. Every string is set with textContent — a detail or remedy is
  // producer text that can carry any punctuation and must never become markup.
```

with:

```
  // Nearly the same exemptions as the anomalies strip above, for the same
  // reason: it lives outside the view panels, ignores the recently-done window,
  // and ignores the shared filters, because a finding must not be hideable by a
  // filter combination. The one exception is the Activity view, which hides the
  // strip from applyView (board-controls.js, REQ-578) — this renderer decides
  // emptiness and nothing else. Every string is set with textContent — a detail
  // or remedy is producer text that can carry any punctuation and must never
  // become markup.
```

Verified before handing it over: the quoted "replace" block is the exact current text of those five lines on `main` at `09a13839`, and the replacement changes only the comment (no code moves, no line the tests read).
