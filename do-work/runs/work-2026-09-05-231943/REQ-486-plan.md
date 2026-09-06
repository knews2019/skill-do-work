## Plan

### Angle

Fold first, summary second. The two features in this REQ share only the group markup, and the fold half is a net deletion: REQ-236 already solved every problem the By UR fold has (head cannot be both fold control and drawer trigger, so `data-detail-*` moves to a sibling `button.ur-group-detail` inside `div.ur-group-row`), and REQ-486 is that same move with the initial state inverted. Everything Go-side, every disclosure rule, all the layout risk and the browser lane belong to the summary half.

Task order puts the two riskiest assumptions early inside their own halves: T1 carries the markup change that silently breaks two existing browser probes, T4 carries the rewiring of the board's only clock.

---

### T1 — Fold the By UR group by merging the head's two branches into one

**Files:** `skills/do-work-board/tools/queue-kanban/web/board-cards.js`, `web/board.css`, `javascript_behavior_a_test.go`, `generate_test.go`, `user_request_clipboard_browser_probe_test.go`

`renderUserRequestLens` (`web/board-cards.js:925-1081`) has two shapes of head. The By UR branch sets `head.dataset.detailKind = "ur"` at `:992-993` and appends cards eagerly at `:1039-1040`. The folded branch wraps the head in `div.ur-group-row` beside `button.ur-group-detail`, sets `aria-expanded="false"` at `:989`, and builds cards lazily at `:1026-1037`.

Delete the branch. Always build the `div.ur-group-row` wrapper, always put `data-detail-kind`/`data-detail-id` on the sibling `button.ur-group-detail`, always attach the same click listener, always append the fold marker. `userRequestCardsFolded` (`web/board-controls.js:8`) then decides exactly two things: the initial `aria-expanded` value, and whether `makeUserRequestCards()` runs eagerly or lazily on first open. `makeUserRequestCards()` already closes over `shownRequestIds` — reuse it unchanged. The delegated `[data-detail-kind]` handler in `board-controls.js` needs no edit; it finds the attribute on whatever node carries it.

`.ur-count` (`board-cards.js:1006-1014`) stops carrying the total. It renders only when a filter hides members ("12 of 43 shown") and is omitted otherwise. The denominator moves to the summary strip in T4.

**CSS.** `board.css:1552-1601` scopes REQ-236's rules under `.user-request-lens.is-folded`, but `.ur-group-row`, `.ur-group-row .ur-group-head`, `.ur-group-row + .ur-group-cards`, `.ur-group-detail` and `.ur-fold-marker` are already unscoped and apply as written. Keep `.user-request-lens.is-folded`'s tighter `gap: 8px` and head padding for URs only.

**Two existing browser probes select the drawer trigger by the markup this task changes, and both must be repaired in this same commit:**

- `user_request_clipboard_browser_probe_test.go:142` — `.ur-group-head[data-detail-id="..."]`
- `generate_test.go:1037` — `.ur-group-head[data-detail-id="UR-075"]`

Both become `.ur-group-detail[data-detail-id="..."]`. Every other assertion in both probes stays as written. The fast stage excludes `TestBrowserBehavior*`, so these stay green until the heavy lane runs — a builder who checks only the fast stage ships a broken lane.

**`javascript_behavior_a_test.go:1652` is inverted, not deleted.** It currently reads:

```go
if result.ScopedByUserRequest[0].Expanded != "" {
    t.Fatalf("by-UR head carries aria-expanded=%q; the fold must not leak into the always-open lens", ...)
}
```

The replacement demands the positive pair in one assertion: the By UR head starts `aria-expanded="true"` with cards present, the URs-only head starts `"false"` with none. Name REQ-486 in the failure message as the inversion so a later reader cannot mistake it for a regression. `renderedUserRequestRow` (`generate_test.go:2791-2796`) gains one field for the Details button's presence.

**Proves it (heavy JavaScript lane):** first render shows cards and `aria-expanded="true"`; activating the head removes the cards and flips to `"false"`; activating again restores the same card ids; collapsing UR-401 leaves UR-402's cards listed; in both readings the head carries no `data-detail-kind` and a sibling button carries `data-detail-kind="ur"`; URs only still starts collapsed.

---

### T2 — Fold the drawer's REQ-id list and cap its height

**Files:** `web/board-detail.js`, `web/board.css`, `javascript_behavior_d_test.go` (new)

`openUserRequestDetail` (`web/board-detail.js:614-636`) calls `appendMetaRow("Grouped REQs", ...)` then `appendMetaRow("REQ ids", makeTicketLinkList(requestIds))` with no control between them. That is the 43-row wall in the screenshot.

Add one sibling helper beside `appendMetaRow` (`:321-333`): `appendFoldableMetaRow(label, valueNode)` puts a real `<button class="detail-fold" type="button" aria-expanded="true" aria-controls="...">` carrying the label inside the `<dt>`, and toggles the `<dd>`'s child list through `el.hidden`. `drawerMeta` is a `<dl>`, so the control goes in the `dt` to keep `dt`/`dd` pairing valid. Use `el.hidden`, never node removal: the list is already built, nothing needs rebuilding on reopen, and there is no teardown to get wrong. `drawerMeta.textContent = ""` at `:623` rebuilds the drawer on every open, so fold state resets per open — that is the ephemeral-state constraint, not a bug.

Only the "REQ ids" row folds. "Grouped REQs", "input.md", the body and the glossary stay visible.

Add `max-height: 40vh; overflow-y: auto` to `.detail-dep-list` inside this row. The user confirmed the list starts open, so a fold alone reproduces the reported problem on every single drawer open. The cap keeps `input.md` and the rendered body reachable with zero clicks; the fold stays as the hide-it-entirely escape.

`.detail-fold` gets the same appearance reset and `:focus-visible` outline `.ur-group-detail` already carries. Do not invent a new focus style.

**This is the first probe ever to drive `openUserRequestDetail`** — `grep -rn openUserRequestDetail *_test.go` returns zero hits. Its DOM stub must cover the ten ids read at `board-detail.js:291-300` (`detail-resizer`, `detail-drawer`, `detail-kind`, `detail-id`, `detail-drawer-title`, `detail-meta`, `detail-body`, `detail-glossary`, `detail-copy`, `detail-copy-all`), extending the `makeNode()` idiom at `javascript_behavior_a_test.go:1495-1525`. It must drive the shipped `setDetailTarget` writer (`board-detail.js:315-319`), never assign `currentDetailKind`/`currentDetailId` by hand.

**Proves it:** first open shows `aria-expanded="true"` and a visible list; activating sets `"false"` and hides the list node while the "Grouped REQs" row, the "input.md" row and `drawerBody` are all still present in the same pass.

> **Commit boundary and split seam.** T1 plus T2 satisfy every fold sentence in the GREEN block with no Go change, no payload field, no clock change and no forecast. Commit them as their own increment before T3 starts, whether or not the REQ is formally split.

---

### T3 — Read the nested P50 into the payload and correct the comment that says the board cannot

**Files:** `model.go`, `generate.go`, `timeline.go`, `frontmatter_test.go`, `model_test.go`, `generate_test.go`, `skills/do-work/actions/work-reference.md`, `skills/do-work-board/actions/board.md`

No production reader exists: no `RequestTicket` field, no read in `parseRequestTicket` (`model.go:707`), no `generatedRequest` JSON key, and only `coerceScalarToString` (`model.go:2131`) and `coerceToStringList` (`:2155`) as coercers.

`parseFrontmatter` already returns `fields["estimate"]` as a `map[string]any` whenever the block parses strictly — pinned by `frontmatter_test.go:403-404`, `:432`, `:449`, `:506-507`. The read is a map lookup, not a parser change.

Add `coerceNestedScalarToFloat(value any, nestedKey string) (float64, bool)` beside the two existing coercers, returning `ok=false` for nil, non-map, missing key, non-numeric, non-finite and non-positive. Add `HasEstimateP50ActiveMinutes bool` and `EstimateP50ActiveMinutes float64` to `RequestTicket` beside `EffortEstimate` (`model.go:239`), read them in `parseRequestTicket` next to the `effort_estimate` block at `:790`, and ship them on `generatedRequest` as `hasEstimateP50ActiveMinutes` (omitempty) plus `estimateP50ActiveMinutes` (deliberately **not** omitempty). That copies the precedent stated verbatim at `generate.go:224-228`: the flag carries presence, the field carries the value including zero, because omitempty would drop a genuine zero while the flag still shipped true and leave the client rendering `NaN`. The schema floors P50 at 5 so a real zero is impossible today; copying the pattern keeps the reader correct if that floor is relaxed.

**Do not touch `frontmatter.go`.** Its salvage-path comment at `:107-113` says the lenient path is flat and drops a nested `estimate:` map. That is still exactly true — the new reader uses the strict path — and it is a different claim from the false one below.

**`timeline.go:394` becomes false with this commit.** The doc comment on `timelineChainStart` says "the board parses no nested frontmatter blocks". Rewrite the comment only: the board now reads the block, and this bar deliberately still does not use it. `timelineChainStart` keeps the projection median; the REQ forbids changing Timeline behaviour merely because the board began exposing the value. `grep -rn nested *.go` returns only this claim, `frontmatter.go`'s accurate one, and test comments — no third copy.

**Docs in the same commit.** The `estimate:` block at `work-reference.md:125-139` earns the lock-step clause its display-only siblings carry at `:119` (`sweep`), `:121` (`impact`) and `:123` (`effort_estimate`): parsed by `model.go` into the UR progress summary, display only, keep that parser in lock-step, both changing in the same commit. Add the P50's board role to `board.md:120`'s field-by-field list rather than starting a competing list.

**Proves it (fast stage only — this stage excludes both probe lanes):**
1. Parse level on a strict fixture: `EstimateP50ActiveMinutes == 75` and the flag true.
2. The salvage twin: a fixture whose malformed line forces `lenientFrontmatterFields` asserts the flag **false**, never `== 0`.
3. Projection level: the generated JSON carries `estimateP50ActiveMinutes` for the strict fixture; for the salvage fixture the flag key is absent and the value key still serialises, so no client can multiply undefined.
4. Timeline invariance: `estimate:` blocks added to a synthetic tree move no projection row, median, sample count or confidence flag.

A parse-level assertion alone cannot prove the read path is wired; both altitudes are required.

---

### T4 — One shared rollup, both surfaces, one clock instant

**Files:** `web/board-user-request-summary.js` (new), `web/board-core.js`, `web/board-cards.js`, `web/board-detail.js`, `web/board.js`, `web/board.css`, `generate.go`, `generate_test.go`, `javascript_behavior_d_test.go`

**Where the code lives.** New fragment `web/board-user-request-summary.js`, registered in `boardJavaScriptFragmentPaths` (`generate.go:43-55`) at position 7, after `web/board-timeline.js` and before `web/board-activity.js`. Every dependency is then declared above it: `createElement`, `formatElapsedDuration`, `futureInstantSkewAllowanceMs` from `board-core.js` at position 1, and `timelineFormatSpanMinutes` (`board-timeline.js:207-218`) at position 6. Nothing calls forward. The manifest is pinned twice — `generate_test.go:48` (authored inventory) and `:76` (execution order) — and both copies change here; a miss is a loud fast-stage failure. Two things stay in `board-core.js` because they belong to it: the status predicates, and the clock fan-out the boot block calls.

**Predicates.** There is no JS `isCompletedStatus` today; `grep -rn isCompletedStatus web/*.js` finds only comment mentions at `board-calendar.js:10` and `board-core.js:256`. Add `isCompletedStatus(status)` beside `isTerminalResolvedStatus` (`board-core.js:261-263`) and re-express the resolved one as `isCompletedStatus(status) || status === "cancelled"`, mirroring the Go composition at `model.go:1007`/`:1030`/`:1038` instead of writing a fourth literal list for "successful". Three existing probes slice the shipped `isTerminalResolvedStatus` (`javascript_behavior_a_test.go:1457`, `_b:87`, `_c:483`) and must slice the new sibling too, or they throw `ReferenceError` — a loud failure in the JavaScript lane.

**The rollup.** `summarizeUserRequestProgress(userRequestId, nowMs)` is pure: no DOM, no formatting, `nowMs` a parameter and never a `Date.now()` call inside. It reads `userRequestsById[id].requestIds` — the full membership, never `shownRequestIds` — so no filter can move the denominator. It returns one object: `totalCount`, `successfulCount`, `resolvedCount`, `activeMinutes`, `excludedSpanCount`, `unmeasuredCount`, `skewedClaimCount`, `liveClaimCount`, `remainingMinutes`, `unknownForecastCount`, `remainingIsPartial`.

- **Active time** sums `implementationSpanMinutes` only where `hasImplementationSpan` is true and `implementationSpanReason` is empty, plus live elapsed for claimed members. Never re-measures a span in the browser. A non-empty reason increments `excludedSpanCount`; a false flag with an empty reason increments `unmeasuredCount`. A claimed member whose `claimedAt` is unparseable or more than `futureInstantSkewAllowanceMs` (`board-core.js:105`) ahead of `nowMs` contributes nothing and increments `skewedClaimCount` — disclosed, never clamped to zero, matching `formatElapsedDuration`'s skew marker at `:129-131`.
- **Remaining time**, three arms per unfinished member: saved `estimateP50ActiveMinutes` when `hasEstimateP50ActiveMinutes`; otherwise `boardData.timeline.projection`'s `trivialMinutes` for effort-mechanical and `normalMinutes` otherwise, but only while `projection.confident` is true; otherwise unknown, incrementing `unknownForecastCount` and setting `remainingIsPartial`. A claimed member subtracts its live elapsed and floors at zero. Failed members are unknown even with a saved estimate. Pending, pending-answer and blocked members keep their full estimate.
- **Zero members** yields the unavailable token and performs no division.
- The function is total: an unknown id returns a well-formed result rather than reaching into undefined. Narrowing, not `try`/`catch` — a swallowed exception hides the bug instead of the freeze.

**Formatting** uses `timelineFormatSpanMinutes` (minutes and up), not `formatElapsedDuration` (seconds below an hour). A whole-UR budget stated to the second is false precision, and it makes the 1 Hz tick nearly free because the label rarely changes. Render "under a minute" rather than "0 min" for a sub-minute live claim.

**Disclosure vocabulary, fixed here because the assertions pin it:** `~` prefix for anything forecast, "at least" prefix for a known-partial sum, the literal word "unavailable" for a percentage with no denominator, and explicit `N excluded` / `N unmeasured` / `N unknown` / `⚠ clock skew` suffixes. Carried by words and symbols only — `prime-kanban-board.md` records `--ink-faint` against `--ink-soft` at 1.29:1 light and 1.82:1 dark, which is not a channel, and `.ur-count` already uses `--ink-faint`.

**The two call sites.** `makeUserRequestSummaryStrip(summary, options)` builds the DOM once; the header and the drawer pass different density options and read the same returned object.

- **By UR header:** the strip is a **sibling** of `div.ur-group-row` inside `section.ur-group`, never a child of the head. `.ur-group-head` is a `<button>` (`board-cards.js:982`), so anything inside it joins its accessible name — the fold control already announces "UR-081 alpha request 43 REQ" and five more metrics would make it announce a paragraph. The strip also gets its own `flex-wrap: wrap` row, so it can never squeeze `.ur-title` (`board.css:1515-1524`: `flex: 1; min-width: 0; white-space: nowrap; text-overflow: ellipsis`). **Do not add `flex-wrap` to `.ur-group-head`** — with `flex: 1` on the title and `flex: none` on chips the line does not break, the title just compresses to an ellipsis. The strip stays visible when the group is collapsed; folding removes only the card grid.
- **UR drawer:** the metrics become meta rows above the foldable REQ-ids row from T2.
- **Not in the URs-only reading.** The REQ names two surfaces and separately constrains URs only to keep its existing behaviour; that reading exists to see many URs at once (its own CSS comment at `board.css:1547-1551`) and its rows are deliberately tighter.

**The clock.** `web/board.js:68` is the board's only ticker: `setInterval(refreshRelativeTimeNodes, 1000)`. The attribute route cannot carry this — a UR's live contribution is a fixed base plus N live claims, which is not one instant for N ≥ 2. Give `refreshRelativeTimeNodes` a `nowMs` parameter (no test slices it; only comment mentions at `javascript_behavior_b_test.go:1681`/`:1749`) and add `refreshTickingSurfaces()` beside it: capture one `nowMs`, call `refreshRelativeTimeNodes(nowMs)` **first**, then re-render every node matching `[data-ur-summary-id]` from that same instant. Point `board.js:68` at it. Only URs with a live claimed contribution carry the attribute, copying the deliberate opt-out at `makeImplementationSpanNode` (`board-cards.js:66`).

No subscriber registry. `renderUserRequestLens` rebuilds `host` from scratch on every call (`board-cards.js:927`), so a registry would accumulate references to dead nodes and would need deregistration on drawer close. An attribute pass has neither problem: a node not in the document is not selected. Add no document listener from `board-core.js` — manifest position governs listener registration order.

**Proves it (heavy JavaScript lane, `javascript_behavior_d_test.go`, every probe slicing the shipped source with `sliceBalancedBlockAfter`):**
1. **One instant:** stub `Date.now` to a counter; drive `refreshTickingSurfaces` once; assert it was called exactly once for the whole tick.
2. **Freeze guard:** drive `refreshTickingSurfaces` against a board whose UR payload is missing or incomplete and assert a `[data-instant-ms]` stopwatch node **still shows the new label**. This is the assertion that protects every existing ticking surface.
3. **Immovable denominator:** with a status filter hiding 3 of 5 members, `totalCount` is still 5 and both percentages are unchanged from the unfiltered run.
4. **Status split:** completed, completed-with-issues, cancelled, failed, pending → successful 2, resolved 3, failed in neither.
5. **Zero members:** unavailable token, no division.
6. **Live sum:** two claimed members 30 and 90 minutes old → 120 min; re-call at `nowMs + 60000` → exactly 122 min.
7. **Skew:** a claim stamped 10 minutes ahead → `skewedClaimCount` 1 and the qualifier **present** in rendered text. Assert presence of the qualifier, never absence of a number.
8. **Three forecast arms, each able to fail in both directions:** saved P50 75 → 75; no P50 with `confident: true`, `normalMinutes: 60` → 60; no P50 with `confident: false` plus a `declinedReason` → unknown and partial. Plus a claimed member with P50 90 and 100 minutes elapsed contributing 0, not −10.
9. **The central claim:** in one run with one stubbed clock, render the By UR header **and** open the drawer for the same UR, then assert the five metric strings are byte-identical **between the two rendered surfaces**. Do not call `summarizeUserRequestProgress` twice and compare its own output — a probe that calls the function under test directly cannot hold its call site.
10. **Tick reaches both:** advance the stub, call `refreshTickingSurfaces` again, assert the header figure, the drawer figure and a claimed card's stopwatch all moved by the same delta.
11. **Source token** in `generate_test.go`: the boot line reads `setInterval(refreshTickingSurfaces, 1000)`.

The four-hour outlier ceiling (`durations.go:32`) and the five-sample confidence floor (`timeline.go:267`) are **not** re-tested here. The client reads `hasImplementationSpan` / `implementationSpanMinutes` / `implementationSpanReason` and never re-derives either rule, so this lane expresses paused, reversed and unmeasured members as payload shapes in the inline `boardData` literal. Both rules stay covered where they live, in `durations_test.go` and `timeline_test.go`.

---

### T5 — Browser evidence, board guide, release

**Files:** `user_request_progress_browser_probe_test.go` (new), `skills/do-work-board/docs/board-guide.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `skills/do-work/actions/version.md`

**The lane skips by default in this environment and a skip is not a pass.** `maintainer-verify.sh:50-63`'s `browser_engine_available` returns 0 immediately when `QUEUE_KANBAN_BROWSER` is set, and otherwise checks PATH for `google-chrome`, `google-chrome-stable`, `chromium`, `chromium-browser`, `chrome` — none of which resolve here. `/opt/pw-browsers/chromium` is a symlink outside PATH. **Every browser-lane invocation must export `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium`**, or the pixel half of this REQ silently never happens and the gate still exits 0.

New probe built on `generateLiveSiteInDir` plus `runBrowserBehaviorProbeInDirectory`, so it drives the real generated page with `board-data.js` beside `index.html`. Model the Go-level fixture on `userRequestCopyFixture` (`user_request_clipboard_browser_probe_test.go:28`), the one existing builder that composes `*RequestTicket` values across queue, working and archive; extend it with the statuses the RED case names plus one member carrying `estimate.p50_active_minutes` and one without.

Return `location.href` in the **same** evaluate call as every measurement, and record the exact Chromium build beside every number. `prime-kanban-board.md` records 141 as deprecated and not a compatibility target, so a green on 141 is evidence about this run, not the compatibility claim the prime reserves for current stable.

**What only this lane can prove:**
1. The strip's `getBoundingClientRect().top` is at or below the head row's `.bottom` at 320, 768 and 1280 CSS px — the metrics are on their own row at every width.
2. No pairwise rect intersection among `.ur-title`, the strip's metrics and `.ur-group-detail`; every metric box is contained in its `.ur-group` box.
3. Rendered summary text read back is not truncated at any of the three widths.
4. With a 43-member UR open, the drawer's "input.md" row stays inside the panel with the id list expanded.
5. Contrast of every metric and qualifier measured against `getComputedStyle(document.body).backgroundColor` — not a `--surface-*` token, because the surface behind this board's content is `<body>` — at ≥ 4.5:1 in both themes, driven through the colour-scheme flag (`runBrowserBehaviorProbeWithFlags`). Without that flag Chromium resolves `prefers-color-scheme` to light and nothing automated ever sees the dark palette.
6. Both fold controls and the Details control are `<button>` elements with `tabIndex >= 0`. Trusted Tab reaches the fold then Details in that order; Enter and Space toggle the fold. Synthetic events cannot test Tab at all — being real buttons is what delivers the tab order, and the JavaScript lane only proves `aria-expanded` flips.

**Board guide.** `board-guide.md:25` calls the Lens toggle "flat Columns vs. grouped **By UR**" — two lenses, three ship. `grep -rn "URs only"` across `docs/`, `SKILL.md`, `actions/` and the prime returns nothing, so the guide has been stale since REQ-236, before this REQ touches it. Name the third lens, state both fold defaults, and say what the summary reports and what it refuses to report when evidence is missing. `board-guide.md:52-54` ("Card drawer") documents the REQ drawer only; the UR drawer's rows are undocumented and now gain a summary and a fold. State plainly that active time is measured from the first recorded lifecycle stamp and that remaining time is approximate and can be partial.

**Release** per `_dev/primes/prime-releases.md`: one entry in the root `CHANGELOG.md` whose title says what was delivered, newest on top, title not already used by an earlier entry; a version bump from 0.303.10 in `skills/do-work/actions/version.md`; a byte-identical copy of the root changelog to `skills/do-work/CHANGELOG.md`, enforced by `_dev/tests/shipped-package-reference-contract.sh`. The board tool has no independent changelog.

This REQ adds no write surface, so `prime-kanban-board.md`'s three-write-surface sentence and `_dev/tests/contract-regressions.sh`'s count stay exactly as written.

**Discovered, out of scope:** `board-guide.md` documents the `took ...` badge as a wall-clock span from `claimed_at` to `completed_at`, but `measureImplementationSpan` (`durations.go:222-238`) measures from the earliest origin-eligible stamp. That line was already wrong before this REQ. Correcting it belongs to a follow-up REQ, not to this diff.

---

### Decisions

**D1 — Active time spent is origin-to-completion (`implementationSpanMinutes`), not the REQ's literal "claim-to-completion".** DECIDE & STATE. The REQ's Detailed Requirements say claim-to-completion; its Interfaces section orders "Keep the existing duration outlier verdict … Do not copy their constants or re-derive competing rules in the browser." Both cannot hold, and the Interfaces sentence is the one that names an authority. `durations.go:204-211` records why the existing rule measures from the earliest stamp: REQ-505 carried `planning_at` at 16:49 against a 23:00 claim and a 23:01 completion, and the card read 1m 23s for six hours of recorded work. Origin-to-completion means the header can never disagree with the `took …` badges on the cards below it. Cost: for a UR whose claim stamps were rewritten late, the figure is wider than the REQ's literal words promise. Reversible — it is which existing payload field the client sums. Write the divergence into this REQ file and the commit message.

**D2 — The By UR fold is implemented by deleting the two-branch head and parameterizing only the initial state.** DECIDE & STATE. Net deletion, the two readings stay structurally identical, and neither reading's fold can drift from the other's.

**D3 — `javascript_behavior_a_test.go:1652` is inverted in the same commit that sets `aria-expanded`, never deleted or relaxed.** DECIDE & STATE. A builder who resolves the red by not setting `aria-expanded` has silently dropped Detailed Requirement 3, and a builder who deletes it removes the only thing watching the seam between the two readings.

**D4 — The rollup lives in a new fragment `web/board-user-request-summary.js` at manifest position 7; the predicates and the clock fan-out stay in `board-core.js`.** DECIDE & STATE. Position 7 makes `timelineFormatSpanMinutes` (position 6) a backward reference instead of a forward call from position 1, and keeps roughly 200 lines of domain arithmetic out of the shared helpers file. The manifest is pinned twice in `generate_test.go`; missing either copy fails the fast stage loudly.

**D5 — No tick subscriber registry. `refreshTickingSurfaces()` captures one `nowMs`, calls `refreshRelativeTimeNodes(nowMs)` first, then re-renders `[data-ur-summary-id]` nodes.** DECIDE & STATE. This overrules the exploration's `registerBoardTickSubscriber`. `renderUserRequestLens` rebuilds its host on every call, so a registry accumulates references to detached nodes and needs deregistration on drawer close — the exploration's own open item A2 is the registry's problem, not the requirement's. An attribute pass dissolves it: a node out of the document is not selected.

**D6 — The summary strip is a sibling of the head row inside `section.ur-group`, not a child of `button.ur-group-head`; do not add `flex-wrap` to `.ur-group-head`.** DECIDE & STATE. Two independent reasons. The head is a `<button>`, so contents join its accessible name. And with `flex: 1` on an ellipsised title, `flex-wrap` does not break the line — the title compresses, which the REQ forbids at narrow widths. A dedicated row behaves the same at every width. Cost: one extra row of vertical space per group.

**D7 — The strip stays visible while a By UR group is collapsed; folding removes only the card grid.** DECIDE & STATE. A collapsed group showing only a title is a worse board than one showing title plus progress, and the fold and the summary are then genuinely independent.

**D8 — The summary renders on the By UR header and in the UR drawer only, not in the URs-only reading.** DECIDE & STATE. The REQ names exactly two surfaces and separately requires URs only to keep its existing behaviour. That reading tightens gaps to see many URs at once; five metrics per group would undo it. One conditional to reverse.

**D9 — The drawer id list starts expanded and gets `max-height: 40vh; overflow-y: auto`.** DECIDE & STATE. The user confirmed the open default, so a fold alone reproduces the reported problem on every open. The cap fixes it with zero clicks; the fold stays as the full-hide escape. Two CSS declarations to reverse.

**D10 — No segmented progress meter.** DECIDE & STATE. One of the three input plans proposed one and escalated it. Rejecting it: the REQ asks for five text figures and "a compact metric layout that can wrap", not a new visual channel, and a meter makes colour carry meaning, which needs its own contrast proof in both themes for something nobody asked for. Purely additive later if the maintainer wants it — the text figures stand alone and nothing depends on it.

**D11 — `.ur-count` renders only while a filter hides members ("12 of 43 shown") and is omitted otherwise; the strip owns the total.** DECIDE & STATE. A filter-dependent "12 / 43 REQ" beside a filter-independent "30/43 successful" reads as a contradiction. `grep -rn "ur-count" *_test.go` returns nothing, so the wording is free.

**D12 — Payload is `hasEstimateP50ActiveMinutes` (omitempty) plus `estimateP50ActiveMinutes` (not omitempty); the salvage path keeps dropping the block.** DECIDE & STATE. Copies the pinned precedent at `generate.go:224-228`. `frontmatter.go:107-113` is accurate and is not edited.

**D13 — UR rollup times use `timelineFormatSpanMinutes`, never `formatElapsedDuration`.** DECIDE & STATE. Seconds on a whole-UR budget is false precision, and it makes the 1 Hz tick almost free. Sub-minute live claims render "under a minute", not "0 min".

**D14 — Disclosure is carried by words and symbols, never by a second ink tone.** DECIDE & STATE. `--ink-faint` against `--ink-soft` measures 1.29:1 light and 1.82:1 dark. `.ur-count` already uses `--ink-faint`, so the obvious fainter-tone move is the one that has already failed twice on this board.

**D15 — New probes go in new files: `javascript_behavior_d_test.go` and `user_request_progress_browser_probe_test.go`.** DECIDE & STATE. Only `javascript_behavior_a_test.go:1652` is edited in place, because that assertion belongs to the REQ-236 contract test and must be inverted where it lives.

Nothing here is escalated. The two questions the input plans wanted to escalate (D1 and D8) are both settled by sentences already in the REQ — the Interfaces order for D1, the two named surfaces plus the URs-only preservation constraint for D8. An escalation there would be a deferred question wearing a decision's clothes.

---

### Files I will touch

- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify)
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modify)
- `skills/do-work-board/tools/queue-kanban/web/board-core.js` (modify)
- `skills/do-work-board/tools/queue-kanban/web/board-user-request-summary.js` (new)
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modify)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify)
- `skills/do-work-board/tools/queue-kanban/model.go` (modify)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify)
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modify — comment only)
- `skills/do-work-board/tools/queue-kanban/frontmatter_test.go` (modify)
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modify)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` (modify)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go` (modify — selector)
- `skills/do-work-board/tools/queue-kanban/user_request_progress_browser_probe_test.go` (new)
- `skills/do-work-board/docs/board-guide.md` (modify)
- `skills/do-work-board/actions/board.md` (modify)
- `skills/do-work/actions/work-reference.md` (modify)
- `skills/do-work/actions/version.md` (modify)
- `CHANGELOG.md` (modify)
- `skills/do-work/CHANGELOG.md` (modify — byte-identical mirror)

Paths are repo-relative because `annotateWriteSetOverlap` compares them with `path.Match` against other REQs' declared sets; absolute paths would silently never overlap.

---

### Testing approach

Three lanes, three different questions. Record the actual exit line for each — a silent skip reads as a pass.

**Fast stage.** Owns the Go payload and nothing else. Runs with `DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior`, so neither client lane executes here.

```
cd skills/do-work-board/tools/queue-kanban && \
  QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off \
  DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior go test -count=1 ./...
```

The nested read at parse level and again at projection level; the strict-versus-salvage pair asserting absence rather than zero; Timeline invariance; the fragment manifest at both pinned copies. Where a threshold fixture is added, use a pair straddling the real boundary derived from the constant, not one that spans it widely.

**Heavy JavaScript lane.** Owns every semantic claim, because it slices the shipped source rather than a copy. Node v22 is on PATH at `/opt/node22/bin/node`; probes pipe to `node -` on stdin, never `-e`.

```
QUEUE_KANBAN_BROWSER_PROBES=off QUEUE_KANBAN_JAVASCRIPT_PROBES=on \
  QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 go test -count=1 -run '^TestJavaScriptBehavior' -v .
```

Covers both folds with opposite defaults, the immovable denominator, the successful/resolved split with failed in neither, the zero-member case, all three forecast arms in both directions, the live sum advancing by exactly the clock delta, the skew disclosure, the one-`Date.now`-per-tick proof, the freeze guard, and the byte-identical agreement between the two rendered surfaces. The prior run recorded this lane green at HEAD in 7.588s; take a fresh baseline before the first edit so a later red has a known-green predecessor.

**Heavy browser lane.** Owns pixels and trusted input only.

```
QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=on \
  QUEUE_KANBAN_STRICT_BROWSER_BEHAVIOR=1 QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  go test -count=1 -run '^TestBrowserBehavior' -v .
```

Wrap and collision at 320/768/1280 in both themes, drawer containment with a 43-member UR, contrast against `document.body`, and real-button tab order. The prior run recorded this lane green at HEAD in 97.100s on Chromium 141.0.7390.37.

**What proves what, one line each:**
- Parser lock-step → fast stage, parse level and projection level. A normalizer-only test cannot hold this line.
- Presence versus zero → fast stage, the salvage fixture asserting the flag false, never the value `== 0`.
- Header and drawer agree → JavaScript lane, comparing rendered text from both shipped call sites in one run under one stubbed clock.
- The clock did not freeze → JavaScript lane, the positive assertion that an existing `[data-instant-ms]` stopwatch still updates when the UR pass has nothing to work with.
- Layout at narrow widths → browser lane only. No JavaScript assertion can see it.
- Tab order → nothing proves it. Focus movement is a trusted-input default action; both controls being real `<button>` elements is what delivers it, and the probes only prove `aria-expanded` flips. Say so in the hand-back rather than implying coverage.

**Closing gate:** `bash _dev/tests/maintainer-verify.sh --heavy` with `QUEUE_KANBAN_BROWSER` exported, exit zero, with both heavy-lane lines present in the output rather than their skip lines. Plus `bash _dev/tests/contract-regressions.sh` (write-surface count unchanged at three) and `bash _dev/tests/shipped-package-reference-contract.sh` (changelog mirror byte-identical).

---

### Biggest risk

Repointing `web/board.js:68` can freeze every live surface on the board, and nothing in the current suite would fail.

That line is the only ticker. It drives every claim stopwatch, every relative time, every state timer, and the clock-skew tooltip sync. T4 points it at `refreshTickingSurfaces`, which runs a second pass over `[data-ur-summary-id]` nodes. If anything in that pass throws — an unknown UR id, a missing `requestIds` array, a summary called before `boardData.timeline` exists — the interval callback dies. The board renders perfectly on load, every number is correct, and then nothing updates again. Every stopwatch silently freezes at its first value, which looks exactly like a board full of very young claims.

The existing suite cannot see it: `setInterval` never runs inside a Node probe, every probe calls render functions directly, and the browser lane does not wait a second and re-measure. The failure passes all three lanes and reaches a user.

Three things contain it, in order of value. First, the freeze-guard assertion in T4 item 2 — drive `refreshTickingSurfaces` against a board with a missing or incomplete UR payload and assert an existing stopwatch node still shows the new label. Second, call order: `refreshRelativeTimeNodes(nowMs)` runs first, so even an unguarded throw in the new pass cannot cost the existing surfaces their current tick. Third, keep `summarizeUserRequestProgress` total by narrowing, not by `try`/`catch` — a swallowed exception hides the bug instead of the freeze, which is worse.

**Second risk, contained by D1.** The summary becomes a browser-side second authority on how long work took. The path is concrete: a cancelled member counts toward resolved but ships `hasImplementationSpan: false`, because `generate.go:733-738` measures a span for terminal success only. A builder sees the sum come out low, notices `claimedAt` and `completedAt` sitting in the payload, and subtracts one from the other in JavaScript. That produces a plausible number with no outlier rule and no origin correction, passes its own probe, and states a figure the Durations view refuses to state — for exactly the REQs whose bookkeeping was worst, which is where a reader is most likely to be checking. If a review finds arithmetic on a raw timestamp inside the summary fragment, that is the failure, regardless of whether the number looks right.

**Third risk, contained by T5's mandatory environment variable.** The browser lane skips and the gate exits 0. One forgotten export and the pixel half of this REQ is never proven.

---

### Split recommendation

Split at the seam between T2 and T3 — between the folding and the summaries. These are two requests wearing one REQ number.

**REQ-486 keeps the folds** (T1, T2). By UR groups fold, the drawer's REQ-id list folds and is height-capped, both start expanded, both are real buttons announcing `aria-expanded`, Details stays a separate control, URs only is unchanged. Five files. No Go change, no payload field, no clock change, no forecast, no new authority to keep in lock-step. It is a net deletion of code and it ships the thing the screenshot actually complains about.

**An addendum REQ takes the summaries** (T3, T4, T5). The nested P50 reader and its payload field, the shared rollup, the new status predicate, the clock fan-out, both surfaces rendering identical values, the strip layout, the browser lane, the stale `timeline.go` comment, and four of the five doc edits. Seventeen files, both languages, all three lanes.

Why this seam and not the other obvious one: splitting by surface (By UR versus drawer) would be wrong, because the REQ's central claim is that both surfaces report identical values, and that claim cannot be asserted until the second surface lands. Fold-versus-summary gives each half a self-contained GREEN and runs in one direction — the fold half creates the `div.ur-group-row` structure and the collapsible drawer region, the summary half decorates both. Rework cost is a few lines of CSS.

The fold half also carries the one inverted assertion and the two repaired browser selectors, which are far easier to review as a five-file commit than buried in a twenty-two-file one.

If the split is refused, the ordering above still holds and T1 plus T2 must be committed as their own increment before T3 starts.

---

### Plan validation (orchestrator)

**Requirement coverage.** Every Detailed Requirement, Interface and Constraint maps to a task:

| REQ requirement | Task |
|---|---|
| By UR header folds its grid independently; starts expanded; multiple at once | T1 |
| Fold control keyboard-operable, exposes `aria-expanded` | T1 (real buttons, ARIA flip), T5 (trusted Tab) |
| Details stays a separate control | T1 |
| Drawer REQ-id list collapsible, starts expanded, hides only ids | T2 |
| URs only preserved: collapsed default, same filters and scope, DOM-only state | T1 (initial-state parameter), probe unchanged apart from the inversion |
| Both surfaces show total, active time, remaining, successful %, resolved % | T4 |
| Complete membership; filters never move the summary or its denominator | T4 |
| Successful = completed + completed-with-issues; resolved adds cancelled; failed neither | T4 (`isCompletedStatus` composition) |
| Percentage with count and total; zero members unavailable | T4 |
| Active time = accepted spans + live claimed | T4, under D1 |
| Excluded / unmeasured / unusable-claim disclosed, never counted as zero | T4, D14 |
| Remaining from saved P50, else confident Timeline median | T3 (payload), T4 (three arms) |
| Claimed member subtracts live elapsed, floors at zero | T4 |
| Pending / pending-answer / blocked keep estimated effort | T4 |
| Failed and no-source members unknown; known forecast explicitly partial | T4 |
| Forecast marked approximate; labels readable when the header wraps | T4 (copy, D6 layout), T5 (measured) |
| Live contributions refresh from the board's existing clock, no drift | T4 (one captured instant) |
| Read nested `estimate.p50_active_minutes` into model and payload | T3 |
| One shared summary function, no per-surface formulas | T4 |
| Keep Go duration and projection authorities; no browser re-derivation | D1, T4 |
| Update the board guide and the stale nested-blocks comment | T3 (`timeline.go:394`), T5 (guide) |
| Do not change Timeline scheduling or forecast behaviour | T3, comment only; fast-stage invariance assertion |
| No write surface; fold state stays ephemeral | T1, T2; contract-regressions count unchanged |
| Explicit unavailable / excluded / partial / approximate formatting | T4, D14 |
| GREEN gate: queue-kanban tests and `maintainer-verify.sh` exit zero; both themes, both widths, no collisions or clipping | T5 |

**Orphan tasks.** None. One piece of work is not a REQ sentence and needs naming: the two browser-probe selector repairs in T1 (`user_request_clipboard_browser_probe_test.go:142`, `generate_test.go:1037`). They are forced by the requirement that Details stays a separate control, and they must land in the same commit as the markup change or the heavy browser lane breaks while the fast stage stays green. Two of the three input plans missed them entirely.

**Task count.** Five, at the limit, and the honest seam is named rather than used to pad the count: T1-T2 are the folds, T3-T5 are the summaries, and T2 is a commit boundary either way. The input plans ran to eight, seven and five tasks; the eight-task version split the payload, the pure function, the wiring and the tick into four separate tasks, which is revert granularity rather than task granularity — T4 here is one coherent unit because the function, its two call sites and the clock instant are only correct together, and the freeze guard cannot be written without all three.

**Consumer field contract.** `write_set` paths are repo-relative, matching `path.Match` semantics in `annotateWriteSetOverlap`; the "Files I will touch" list above is the source and `write_set:` is mirrored from it in Step 5.5, one direction only. Twenty-two entries, all literal paths, no globs and no directory entries — a directory entry never matches a file inside it. Two new source files (`web/board-user-request-summary.js`, and the fragment's manifest entry in `generate.go`) plus three new test files are declared. The release trio (`CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `skills/do-work/actions/version.md`) is in the set because this change touches shipped files under `skills/`.

**Where the three plans disagreed, and what I took.**

- **Summary location.** Took the new fragment at manifest position 7 over `board-core.js`. It makes `timelineFormatSpanMinutes` a backward reference instead of a commented forward call from position 1, and keeps domain arithmetic out of the shared helpers file. The predicates and the clock fan-out stay in `board-core.js` because the boot block calls them.
- **Tick mechanism.** Took the attribute pass over the subscriber registry the exploration proposed. The registry creates a node-lifetime problem the requirement does not have, since `renderUserRequestLens` rebuilds its host every render and the drawer would need deregistration on close. A node out of the document is simply not selected.
- **Header layout.** Took the sibling strip over `flex-wrap` on `.ur-group-head`. The head is a `<button>`, so metrics inside it join its accessible name, and `flex-wrap` does not save an ellipsised `flex: 1` title from being squeezed.
- **Active time measurement.** Decided rather than escalated. Two plans wanted to escalate origin-versus-claim; the REQ's own Interfaces section orders reuse of the existing authority, which settles it. The divergence from the REQ's literal wording goes into the REQ file and the commit message.
- **Segmented meter.** Rejected. One plan escalated it; it is new visual machinery, with colour as a meaning channel, that the REQ did not ask for and that would need its own two-theme contrast proof. Purely additive later.
- **Drawer height cap.** Took it from the reader-first plan. The user confirmed the list starts open, so a fold alone reproduces the screenshot's problem on every open.
- **Minute-resolution formatting.** Took it. A whole-UR budget stated to the second is noise, and it makes the 1 Hz tick almost free.
- **Summary in URs only.** Took "no" over "yes". The REQ names two surfaces and separately constrains the third to keep its behaviour.