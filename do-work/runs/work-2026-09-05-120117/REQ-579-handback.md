# Hand-back — REQ-579 (render verify findings and skipped probes as compact rows in one list)

## Branch

- Branch: `worktree-agent-REQ-579-finding-rows`
- Head commit: `1de64c4c` — `[REQ-579] clear both findings hosts before hiding the strip, and pin the ids applyView reads`
- First commit: `1dc13ef7` — `[REQ-579] render verify findings and skipped probes as compact rows in one list`
- Base: `5f4821ab`
- Two commits. Working tree clean; nothing under `do-work/` was written from the worktree.

## File manifest

Source (all under `skills/do-work-board/tools/queue-kanban/`):

- `verify.go` (modified) — `Subject` field on `VerifyFinding`; 23 probe sites set it; new `releaseFindingSubject` constant; the undetermined-worktree remedy no longer points at the deleted disclosure footer.
- `generate.go` (modified) — `Subject` on `generatedVerifyFinding` with `json:"subject,omitempty"`, filled through `reduceAbsolutePaths` in `attachVerifyFindings`.
- `web/template.html` (modified) — the `board-anomalies-cards` grid and the `<details>` disclosure are gone; one `#board-findings-rows` list holding two pass-through hosts.
- `web/board.css` (modified) — every `.board-finding*` and `.board-findings-skipped*` rule deleted, replaced by the row classes.
- `web/board-cards.js` (modified) — `renderVerifyFindingsStrip` rewritten plus four new helpers: `formatFindingsSummary`, `groupFindingsBySubject`, `makeFindingRow`, `makeSkippedProbeRow`. Both hosts are now cleared before the empty check (see D-07).

Tests:

- `verify_test.go` (modified) — added `TestVerifyNamesTheSubjectEachFindingIsAbout` and `TestVerifyNamesTheChangelogAsTheReleaseFindingSubject`; added the `encoding/json` import.
- `javascript_behavior_c_test.go` (modified) — added `TestJavaScriptBehaviorVerifyFindingsRenderAsOneRowList` and two helpers (`hasClassName`, `sliceFindingsStripMarkup`); edited REQ-578's `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip` (see Cross-REQ test changes below).

Not touched: `web/board-controls.js`, the completion-anomalies strip, `CHANGELOG.md`, `VERSION`, everything under `do-work/` except this file.

## P-A-U

**[PLAN]**

1. Go producer: add `Subject` to `VerifyFinding`, set it at every probe that has a natural subject (worktree name / REQ id / `CHANGELOG.md`), leave it empty where there is none. Mirror it into `generatedVerifyFinding` as `subject,omitempty` and run it through `reduceAbsolutePaths` like the sibling strings, because the stray-file subject is a path and a static snapshot is shareable.
2. Client: replace the card loop and the disclosure with one flat row list. Group findings by exact string match on `subject`, print each non-empty subject once as a heading, keep subjectless findings in producer order after the grouped ones, then the skipped probes. Two weights only, both from the payload.
3. Markup: delete `board-anomalies-cards` and the `<details>` block. Keep `#board-findings` and its `hidden` attribute exactly as REQ-578 expects.
4. CSS: replace `.board-finding*` / `.board-findings-skipped*` with row classes, no card border, no grid.
5. RED first in both lanes, then implement, then render a board and look at it.

**[APPLY]** — coded as planned, strictly inside the declared write set. The one plan change made during the build is D-01 below (two hosts instead of one, forced by REQ-578's own test). One thing the plan did not anticipate and the render did: D-04.

**[UNIFY]**

`git diff --stat` on the branch:

```
 generate.go                  |  13 +-
 javascript_behavior_c_test.go| 287 ++++++++++++++--
 verify.go                    |  35 ++-
 verify_test.go               |  85 ++++++
 web/board-cards.js           | 139 ++++++++--
 web/board.css                | 101 +++++---
 web/template.html            |  32 ++-
 7 files changed, 594 insertions(+), 98 deletions(-)
```

Linters and per-file checks:

| File | Checked | Result |
|---|---|---|
| `verify.go` | `gofmt -l .`, `go vet ./...` | clean; struct literal field alignment reformatted by gofmt |
| `generate.go` | `gofmt -l .`, `go vet ./...` | clean |
| `verify_test.go` | `gofmt -l .`, `go vet ./...` | clean |
| `javascript_behavior_c_test.go` | `gofmt -l .`, `go vet ./...` | clean |
| `web/board-cards.js` | `node --check` | clean |
| `web/board.css` | rendered in Chromium and inspected | every value is an existing `var(--…)` token; no new colour literals, so both themes follow the tokens |
| `web/template.html` | page generated and rendered | strip renders; no `<details>`, no card grid |

Debug artifacts: `git diff -U0 | grep '^+'` searched for `console.log`, `debugger`, `fmt.Print`, `TODO`, `FIXME`, `XXX` — no matches.

## Test evidence

All commands run from the worktree root.

**RED — Go lane.** `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run 'TestVerifyNamesThe' ./...` → exit 1.

First run failed to compile (`finding.Subject undefined`), which is not an assertion. The field was declared with no producer wiring it, and the run was repeated so the failure is the assertion itself:

```
--- FAIL: TestVerifyNamesTheSubjectEachFindingIsAbout
    verify_test.go:2387: unmerged-worktree-leftover subject = "", want the worktree name "worktree-agent-REQ-506-focused-evidence"
    verify_test.go:2387: worktree-wrote-queue-state subject = "", want the worktree name "worktree-agent-REQ-506-focused-evidence"
    verify_test.go:2407: the board payload carries no "subject":"worktree-agent-REQ-506-focused-evidence" — the page cannot group rows it never received
--- FAIL: TestVerifyNamesTheChangelogAsTheReleaseFindingSubject
    verify_test.go:2429: release finding subject = "", want "CHANGELOG.md"
```

**RED — Node lane.** `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehaviorVerifyFindingsRenderAsOneRowList$' ./...` → exit 1:

```
--- FAIL: TestJavaScriptBehaviorVerifyFindingsRenderAsOneRowList
    javascript_behavior_c_test.go:2772: header count = "2", want "2 findings · 1 probe not checked"
    javascript_behavior_c_test.go:2782: the list holds 0 children, want a subject heading plus three rows: []
```

(The first attempt of this probe died on `undefined` rather than asserting, because it read the row host out of the stub map instead of through `getElementById`. Fixed in the probe before the RED above was taken.)

**GREEN — Go lane.** `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...` → exit 0, `tests=394 wall=59s slowest-file=generate_test.go:10.90s limit=<30s`. The whole module passes, not only the new tests.

**GREEN — Node lane.** `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...` → exit 0, `tests=58 wall=7s slowest-file=javascript_behavior_c_test.go:2.18s limit=<30s`. Includes REQ-578's hide-on-Activity test, which drives the real `applyView` against the new markup.

**Mutation check on the two assertions added in `1de64c4c`** — each was reverted in the source and the test failed with the message that names the cause, then the source was restored:

- Removed `skippedHost.textContent = ""` → `stale rows survived the empty re-render (findings=0, skipped=1) — applyView reads those counts and would show the strip again`
- Renamed `id="board-findings-skipped-list"` to `…-rows` in the template → `the strip dropped id="board-findings-skipped-list" — applyView reads it to decide visibility and would dereference null`

The repository gate (`_dev/tests/maintainer-verify.sh`) was NOT run, per the brief.

**Render evidence.** `go run . generate --out <scratch> --repo-root <worktree>` from the queue-kanban module, served over `http://127.0.0.1:8791/` (Playwright blocks `file:`), Chromium via the Playwright MCP, light theme, page URL returned with every measurement.

Two renders were read:

1. The worktree's real queue (561 REQs). Five findings, five distinct worktree subjects, no skipped probes — subjects present and correct, one heading each.
2. The same page with a crafted payload injected into `board-data.js` (six findings across three subjects, one of them subjectless; two fixable; two skipped probes), which is the only way to see grouping, muting and skipped rows together.

What the crafted render measured, at `http://127.0.0.1:8791/index.html`: `stripHidden=false`, header `"6 findings · 2 probes not checked"`, 8 rows, 3 subject headings (`worktree-agent-REQ-506-focused-evidence`, `CHANGELOG.md`, `REQ-540`), `details` elements = 0, `.board-finding` cards = 0, both hosts direct children of `#board-findings-rows`, strip height 398px for 8 rows.

Screenshots kept in the session scratchpad (not in the repo):
`/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/679469bb-0bb5-481e-b46d-c032081867d7/scratchpad/req579-strip-final.png` (final), `…/req579-strip-after.png` (before the spacing fix), `…/req579-findings-rows.png` (full board).

**Browser check of REQ-578's rule against the new markup** (same page, `location.href` returned with the measurement). Driving the real view switcher (`[data-view-target]`), not the stub: `#board-findings.hidden` is `false` on Board, `true` on Activity, `false` again back on Board, `false` on Timeline — with `applyView` reading 9 children from `#board-findings-cards` (3 subject headings + 6 rows) and 2 from `#board-findings-skipped-list` throughout. This is the evidence the stub lane cannot give, since the stub manufactures a node for any id.

**What the render looked like.** A flat list of one-line rows on the strip's tinted background. Each row: a small grey uppercase pill on the left, then the detail, then a green `cleanup can fix` where the producer set the flag, then `→ remedy` in muted text, wrapping to a second indented line when long. Above each group a small bold grey heading with the subject. Skipped probes are two `NOT CHECKED` rows at the bottom, muted, same shape. No borders, no columns, no disclosure. Six findings and two skipped probes now occupy less vertical space than the two cards in the capture screenshot did.

**Cross-REQ test changes** (general.md § Cross-REQ Test-Break Rules — both intentional):

- REQ-578's `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip` lost its `skippedDisclosureOnTimeline` assertion. That assertion read `nodesById["board-findings-skipped"].hidden` — the `<details>` element this REQ deletes. Against the stub `document` it would have auto-created an empty node and passed vacuously forever, which is worse than removing it. The rest of that test is untouched and still passes.
- The same test's `functionBlocks` gained the four new renderer helpers, because `renderVerifyFindingsStrip` now calls them and the probe slices functions individually.

## Lesson evidence

Read in full, both listed in the REQ as dropped for budget:

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — present. Directly applied: REQ-117's point that `board.Warnings` is a free UI channel (not used here, but it is why no new warning was added for the subject field); REQ-083/REQ-458's rule that `Fixable` means exactly "cleanup resolves it" and must never be widened, which is why `fixable` is the only weight the client reads.
- `_dev/primes/lessons-kanban-board.md` — present. Directly applied: REQ-285's "for a rendering change the screenshot IS the test" (two defects there passed `node --check`, `go test` and the strict JS lane and were obvious on sight) — which is what caught D-04 below; REQ-265's warning that a handed-back integration seam is the one part a builder cannot test, which fed D-01; REQ-237's "when a REQ removes a rule, every test asserting that rule's shape is in scope — say so loudly", which is why the REQ-578 test edit is spelled out above.
- `_dev/primes/prime-kanban-board.md` — present. Applied: "generate a board and look at it"; "render evidence must name the page it measured, in the same call that measures it" (every measurement returns `location.href`); "never commit build outputs" (the generated board lives in the scratchpad).

No listed lesson path was missing.

## Decisions

**D-01 (ESCALATE) — the two row hosts keep their old element ids instead of one honestly-named host plus a `board-controls.js` seam.**

The brief offered both. I built the single-host version first and the run proved it unsafe: REQ-578's `applyView` reads `#board-findings-cards` and `#board-findings-skipped-list` to decide whether the strip has content, and with one host named `#board-findings-rows` its own test fails until someone applies a seam by hand. Handing back a knowingly-red test is worse than a naming wart. The markup now nests two hosts, both `display: contents`, inside one `#board-findings-rows` list — one host per payload array, which is the split the renderer already worked in — and the page is one list to the reader.

- *Value:* my branch is green as handed back, no cross-file coordination, no window where a forgotten seam breaks every view switch.
- *Risk:* two element ids now name things they no longer describe (`-cards` holds rows, `-skipped-list` is not a list). Fully reversible — the rename is three lines across two files and is filed under Integration seams as optional. Also `display: contents` outranks the `hidden` attribute's own `display: none`, so neither host may ever be given `hidden`; that trap is written into the template comment next to the ids.

One correction to the constraint as it reached me: the Node behavior lane **does** catch a removed host id, just not by throwing. The stub `document.getElementById` manufactures an empty node, so `findingsStripHasContent` computes `false` and REQ-578's `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip` fails with `two findings rendered but the strip is hidden on the Board view`. That is the failure I actually observed on the single-host build, and it is what made me switch. The markup-level assertion was still added (`1de64c4c`) — it names the cause directly instead of through a visibility symptom, and it holds even if that stub is ever replaced by one that returns `null`.

**D-07 (DECIDE & STATE) — both hosts are cleared before the empty check, not after it.** The old renderer returned early on the empty payload without clearing, and again after hiding an empty skipped list, so a re-render with nothing to report left stale children behind. Latent before this REQ, real after it: `applyView` decides visibility from those two child counts, so stale rows under a hidden strip would put the strip back on screen at the next view switch. Clearing first also removes the second early return, so the renderer now has one clearing point and one exit for the empty case.

**D-02 (DECIDE & STATE) — `Subject` is set on every probe with a natural subject, including the four categories the board suppresses.** `boardRenderedVerifyCategories` keeps completion anomalies, duplicate ids, stray files and unrecognized statuses out of the payload, so their subjects are never rendered today. Setting them anyway is one field per site and makes the rule uniform ("the probe that knows, sets it") instead of a list someone has to maintain. Duplicate-id findings are the exception: they are lifted from warning prose with no parsed id, so they stay empty, as do malformed calibration rows.

**D-03 (DECIDE & STATE) — the header shows `"2 findings"`, not a bare `"2"`, and the count pill loses its pill styling in this strip.** The REQ asked for the skipped count to join the header when non-zero, and `"2 · 1 probe not checked"` inside a numeric pill reads as nothing. Nouns on both halves, no pill, and each half appears only when its count is above zero. The completion-anomalies strip above keeps its numeric pill unchanged.

**D-04 (DECIDE & STATE) — a subjectless row after a group, and the first skipped row, get extra top margin (`.board-findings-row-detached`).** Found by looking at the render, not by a test: D3's ordering puts subjectless findings after the grouped ones, and on the page the first of them sat under the previous heading and read as that group's third row. A subject heading says nothing about rows that are not its own. The renderer marks exactly those two rows; the node test pins both the ordering and the class.

**D-05 (DECIDE & STATE) — `groupFindingsBySubject` uses `Object.create(null)` for its lookup.** A subject spelled `constructor` or `toString` would otherwise hit `Object.prototype` and be treated as an already-open group. Subjects come from file paths and branch names, so this is cheap insurance in one word rather than a guard.

**D-06 (DECIDE & STATE) — the release probes share the subject `"CHANGELOG.md"`.** The REQ allowed "the changelog path or version". The path groups all three release findings under one heading; a version would split them and change every release.

Silent: `.board-findings-detail` carries no CSS rule of its own — it labels the row's main text beside the classed `-fixable` and `-remedy` spans, and an unclassed span among classed ones reads as an oversight.

## Discovered Tasks

- The two row-host element ids (`board-findings-cards`, `board-findings-skipped-list`) no longer describe what they hold, and `applyView` in `web/board-controls.js` is the only reason they cannot be renamed in this REQ's write set. A three-line follow-up renames both and updates that read. `impact-negligible` → report only.
- `do-work/user-requests/UR-118/assets/REQ-579-screenshot-1-verify-findings-two-cards-and-skipped-probes.png` and the REQ's own Assets description now describe a board that no longer exists (two cards plus a `2 probe(s) could not run` disclosure). Any future reader of UR-118 will see a capture of the pre-REQ state with nothing saying so. `impact-negligible` → report only.
- REQ-580's own hand-back flagged the same asset staleness from the other direction (its T-02). The two notes are about one artifact. `impact-negligible` → report only.

## Integration seams

None required. The change is self-contained and green on the branch.

Optional, only if the orchestrator wants the follow-up rename from Discovered Tasks done now rather than later — in `skills/do-work-board/tools/queue-kanban/web/board-controls.js`, replace:

```js
    var findingsStripHasContent =
      document.getElementById("board-findings-cards").children.length > 0 ||
      document.getElementById("board-findings-skipped-list").children.length > 0;
```

with:

```js
    var findingsStripHasContent =
      document.getElementById("board-findings-rows").children.length > 0;
```

That form is only correct if the two inner hosts are also removed from `web/template.html` and the renderer is repointed at `#board-findings-rows` — three files, so it is a separate REQ, not a seam to apply on top of this merge. Do not apply the snippet alone: with the hosts still nested, `#board-findings-rows` always has exactly two children and the strip would never hide.

## Exploration

What the request's capture-time findings named, checked against the code:

- **Producer** — `VerifyFinding` in `verify.go:75`. Correct. What the capture did not say: there are **23 finding-construction sites** across 15 `append*Findings` functions, not a handful, and one of them (`routeWorktreeLeftover`) supplies its category through a variable rather than a constant. `collectVerifyFindings` (`verify.go:149`) is the ordered list of all of them.
- **Payload site** — `generatedVerifyFinding` in `generate.go:610` and `attachVerifyFindings` just below it. Correct, and it is the single point both `generate` and `serve` go through, so one field addition covers both modes.
- **Renderer** — `renderVerifyFindingsStrip` in `web/board-cards.js`. Correct.
- **Classes being replaced** — `.board-finding`, `-head`, `-category`, `-fixable`, `-detail`, `-remedy`, `.board-findings-skipped`, `-summary`, `-list`, `-item`. Correct and complete; a repo-wide grep found no other consumer of any of them outside the strip's own three files.
- **Node-lane test pattern** — `javascript_behavior_c_test.go`. Correct: `generateLiveSite` + `sliceBalancedBlockAfter` per function + a stub `document`.

What the capture got wrong or missed:

1. **The strip's markup was never card-specific to begin with.** `#board-findings` reuses the completion-anomalies section wholesale — `class="board-anomalies board-findings"`, `board-anomalies-head`, `board-anomalies-title`, `board-anomalies-count`, `board-anomalies-hint`, `board-anomalies-cards`. Only the cards host and the disclosure needed replacing; the header markup and its rules are shared with a strip this REQ must not touch, which is why the header changes are two CSS overrides scoped under `.board-findings` rather than edits to the shared rules.
2. **The `board-controls.js` coupling is stronger than "keep the same element".** REQ-578 does not just read `#board-findings.hidden`; it recomputes emptiness on every view switch from the *children counts of two specific hosts*. Any restructuring that removes either id breaks the strip on every view, and it fails loudly in REQ-578's own test rather than silently. This is what decided D-01.
3. **One shipped sentence went stale.** `routeWorktreeLeftover`'s undetermined-worktree remedy told the reader the failed read is listed "under the `probe(s) could not run` footer on the board" — the exact footer this REQ deletes. Fixed in the same commit; `verify.go` was already in the write set. A grep for that phrase found no other copy in shipped files.
4. **The `fixable` flag is rarer than the design assumes.** Of the 23 sites, only two set `Fixable: true` (stranded terminal REQs, and finished worktree residue). The muted weight is therefore an uncommon state, which is why the skipped rows carry it too — otherwise the lighter weight would almost never appear and would read as a rendering accident.
5. **The real queue produces no skipped probes and no shared subjects right now.** The five findings on this worktree's own board are five different worktrees, one row each. Seeing grouping, muting and skipped rows on a page required injecting a crafted payload into a generated `board-data.js`; the plain render alone would have shown none of the three behaviours the REQ is about.
6. **`display: contents` and the `hidden` attribute conflict.** The old renderer set `skippedHost.hidden = true/false`. With a `display` rule on the same element, `hidden`'s own `display: none` loses and the element stays visible. The new renderer never sets `hidden` on either host — an empty host renders nothing — and the template comment says why so the next editor does not reintroduce it.
