# Pre-exploration for REQ-486 — collapsible UR groups + progress summaries

Repo: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` (read-only pass; nothing edited).
Board tool root: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/`
All paths below are relative to that root unless stated otherwise.

---

## 1. The prior implementation (REQ-236, commit `456ee9d`)

### Archived record

`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/UR-054/REQ-236-add-a-urs-only-lens-to-the-board.md`

`commit: 456ee9d` is a **merge commit** (`Merge branch 'worktree-agent-REQ-236-urs-only-lens'`, parents `52fb7d72` + `266ec091`). `git show 456ee9d -- <path>` prints an empty combined diff; use `git diff 52fb7d72..456ee9d -- <path>` to see the change. Merged diff: 5 files, +395/−14 — `generate_test.go` (+247), `web/board-cards.js` (+59/−14), `web/board-controls.js` (+39), `web/board.css` (+53), `web/template.html` (+11).

Key decisions recorded on the archived REQ:

- **D-01** — `URs only` is a *card-fold modifier* on the by-UR lens, not a third `viewState.lens` value. Three files ask "is the UR lens on screen?" by naming `"user-request"` (`board-cards.js`, `board-testing.js`, `board-controls.js`, `board-filters.js`), two of them outside the write set. Keeping `viewState.lens === "user-request"` true for both readings kept all of them correct unedited.
- **D-02** — the fold flag lives in `board-controls.js`, not in `viewState` (which is declared in `web/board.js`, outside the write set). All fragments inline into one IIFE so this is a module-level `var`.
- **D-03** — expanded/collapsed lives in the DOM, not a persisted set. A re-render rebuilds rows collapsed.
- **D-05** — helpers stay inner functions of `renderUserRequestLens`, closing over the state they need; this is what let the probe slice one anchor and get a behavioural RED.

Lessons-learned line worth carrying: *"a document-wide DOM count is only evidence when its root is the thing under test"* — `document.querySelectorAll('.req-card')` returns 25 in a correctly-folded lens because the hidden Columns lens keeps its DOM. Scope every card count to the lens host.

### `renderUserRequestLens` today — `web/board-cards.js:819-975`

Structure (line numbers exact):

```
819  function renderUserRequestLens() {
820    var host = document.getElementById("user-request-lens");
821    host.textContent = "";                       // full rebuild every call
822    var hiddenResolvedCount = 0;
823    var hiddenResolvedFilterMatchCount = 0;
827    var recentlyDoneIdSet = {};                  // built once per render, not per UR
828    recentlyDoneIds(viewState.windowHours).forEach(...)
832    (boardData.userRequestOrder || []).forEach(function (userRequestId) {
833      var userRequest = userRequestsById[userRequestId];
842      var groupMatchesSearch = ...searchMatchesUserRequest(...)
844      var requestIds = userRequest.requestIds || [];     // <- FULL membership
845      var shownRequestIds = requestIds.filter(requestMatchesFilters)  // <- FILTERED
848      if (filterState.userRequestActivity === "active" && !userRequestHasOpenOrRecentWork(...)) { … return; }
857      if (hasActiveFilters() && shownRequestIds.length === 0) { return; }
862      var group = createElement("section", "ur-group");
868      function makeUserRequestCards() {                  // inner, closes over shownRequestIds
869        var cardsNode = createElement("div", "ur-group-cards");
870        shownRequestIds.forEach(cardsNode.appendChild(makeRequestCard(requestId, { showCompleted: true })))
873        return cardsNode;
874      }
876      var head = createElement("button", "ur-group-head");
877      head.type = "button";
878      if (userRequestCardsFolded) {
883        head.setAttribute("aria-expanded", "false");
884        head.appendChild(createElement("span", "ur-fold-marker", "▸"));
885      } else {
886        head.dataset.detailKind = "ur";           // By UR: the head IS the drawer trigger
887        head.dataset.detailId = userRequestId;
888      }
889      head.appendChild(createElement("span", "ur-id", userRequestId));
890      head.appendChild(createElement("span", "ur-title", userRequest.title || "(no input.md title)"));
891-896  optional `badge citation-match` ("cites")
897-899  optional `ur-synthetic` chip ("no input.md")
900-908  head.appendChild(createElement("span", "ur-count",
              shownRequestIds.length < requestIds.length
                ? shownRequestIds.length + " / " + requestIds.length + " REQ"
                : requestIds.length + " REQ"))
909      if (userRequestCardsFolded) {
910        var headRow = createElement("div", "ur-group-row");
911        headRow.appendChild(head);
912        var detailButton = createElement("button", "ur-group-detail", "Details");
913-916    detailButton.type/dataset.detailKind="ur"/dataset.detailId/aria-label
917        headRow.appendChild(detailButton);
918        group.appendChild(headRow);
920        var openCards = null;
921        head.addEventListener("click", function () {
922          if (openCards) { group.removeChild(openCards); openCards = null;
925            head.setAttribute("aria-expanded", "false"); return; }
928          openCards = makeUserRequestCards();
929          group.appendChild(openCards);
930          head.setAttribute("aria-expanded", "true");
931        });
932      } else {
933        group.appendChild(head);
934        group.appendChild(makeUserRequestCards());     // always-open By UR reading
935      }
937      host.appendChild(group);
938    });
940-975  empty-state text + "n URs hidden — switch URs to All" note
```

**How the two readings differ (this is the whole delta):**

| | By UR (`userRequestCardsFolded === false`) | URs only (`=== true`) |
|---|---|---|
| head element | `button.ur-group-head`, **is** the drawer trigger (`data-detail-kind="ur"`) | `button.ur-group-head`, **is** the fold toggle; no `data-detail-*` |
| `aria-expanded` | **absent entirely** | `"false"` initially, flipped to `"true"`/`"false"` by the click listener at :925 / :930 |
| fold marker | none | `span.ur-fold-marker` "▸", rotated 90° by CSS on `[aria-expanded="true"]` |
| wrapper | head appended straight into `section.ur-group` | head wrapped in `div.ur-group-row` beside `button.ur-group-detail` ("Details") |
| cards | `makeUserRequestCards()` appended eagerly at :934 | built lazily on first click (:928), `group.removeChild` on collapse (:923) |
| host class | — | `applyLens` adds `.is-folded` to `#user-request-lens` |

**Details-button collision resolution:** the head cannot be both, so in the folded reading `data-detail-kind`/`data-detail-id` move to a sibling `button.ur-group-detail` inside `div.ur-group-row`. Both are real `<button>`s, so both are keyboard-operable. The delegated `[data-detail-kind]` click handler in `board-controls.js` needs no change — it just finds the attribute on a different node.

**Where fold state lives:** two places, both ephemeral.
- The *lens-level* flag `userRequestCardsFolded` — `web/board-controls.js:8`, a module-level `var` inside the IIFE. Set only by `applyLensSelection` (`board-controls.js:101-109`).
- The *per-group* open state — the `openCards` closure variable at `board-cards.js:920` plus the `aria-expanded` attribute on the head. Nothing outside the DOM; the next `renderUserRequestLens()` call rebuilds everything collapsed. No listener teardown needed because the nodes are discarded whole.

`applyLensSelection` sets `renderedOnce.userRequestLens = false` before `applyLens()` so the cached lens is always dropped when the reading changes (`board-controls.js:105`).

### `makeRequestCard` — `web/board-cards.js:110-413`

`button.req-card` with `data-detail-kind="req"` / `data-detail-id`. Unknown id → bare id + `disabled`. Then `req-card-top` (id + status dot + status text), `h3.req-card-title`, `div.req-card-badges` (citation-match, domain, UR, priority, route, blocked-by, …), and further down the state timer and the implementation span. Options seen in callers: `{ showCompleted: true }` from the UR lens (`:871`) and `{ ... }` at `:464`. Not otherwise relevant to this REQ — the UR summary reads the payload, not the card.

### CSS — `web/board.css`

- `.user-request-lens` :1281 — `flex; column; gap: 18px`
- `.ur-group` :1287, `.ur-group-head` :1294 — `display:flex; align-items:center; gap:12px; width:100%; padding:14px 18px`. **No `flex-wrap`.**
- `.ur-id` :1317, `.ur-title` :1327 (`flex:1; white-space:nowrap; overflow:hidden; text-overflow:ellipsis`), `.ur-count` :1337 (`flex:none`, mono 12px, `--ink-faint`), `.ur-synthetic` :1343
- `.ur-group-cards` :1353 — `grid; repeat(auto-fill, minmax(260px,1fr))`
- REQ-236's block :1364-1412 — `.user-request-lens.is-folded` (gap 8px), `.ur-group-row`, `.ur-fold-marker`, `.ur-group-detail`
- `.ur-lens-empty` :1414, `.ur-lens-hidden-note` :1420

Relevant to the REQ's "must remain readable when the header wraps at narrow widths": the head is a non-wrapping flex row with a `nowrap` ellipsised title at `flex:1`. Adding metric spans without `flex-wrap: wrap` will squeeze the title, not wrap. This is a real CSS change, not just markup.

### Template — `web/template.html:93-107`

Three lens buttons in `#lens-group`; the third carries `data-lens-target="user-request" data-ur-cards="folded"`.

---

## 2. The two surfaces that must show the same summary

**By UR header** — `web/board-cards.js`, inside `renderUserRequestLens()` (:819). The count today is one `span.ur-count` built at :900-908. That is the only summary the header carries.

**UR detail drawer** — `web/board-detail.js:614-636`:

```js
614  function openUserRequestDetail(userRequestId) {
615    var userRequest = userRequestsById[userRequestId];
616    if (!userRequest) { return; }
619    drawerKind.textContent = "UR";
620    drawerId.textContent = userRequestId;
621    drawerTitle.textContent = userRequest.title || "(no input.md title)";
623    drawerMeta.textContent = "";
624    clearDetailGlossary();
625    var requestIds = userRequest.requestIds || [];
626    appendMetaRow("Grouped REQs", String(requestIds.length));
627    if (requestIds.length > 0) {
628      appendMetaRow("REQ ids", makeTicketLinkList(requestIds));
629    }
630    appendMetaRow("input.md", userRequest.inputFilePresent ? "present" : "synthesized from REQ pointers");
632    drawerBody.innerHTML = userRequest.bodyHtml || "<p>(no input.md body)</p>";
633    renderDetailGlossary(linkifyDetailBody(drawerBody, userRequest.title));
634    setDetailTarget("ur", userRequestId);
635    showDrawer();
636  }
```

That is the `GROUPED REQS 43` + long REQ-id list from the screenshot. `appendMetaRow` is `board-detail.js:321-333` (`<dt>`/`<dd>` pair into `#detail-meta`). `makeTicketLinkList` is `board-detail.js:409-422` — one `span.detail-dep` per id, reusing `makeDependencyDetailList`'s flex-column row styles (`.detail-dep-list` at `board.css:1877`).

**Where a shared summary function should live.** The client is assembled from an ordered fragment manifest, `generate.go:43-56`:

```
web/board-core.js, web/board-filters.js, web/board-cards.js, web/board-calendar.js,
web/board-durations.js, web/board-timeline.js, web/board-activity.js,
web/board-testing.js, web/board-detail.js, web/board-controls.js, web/board-clipboard.js
```

All eleven inline into the one IIFE in `web/board.js` at the `/* INLINE_BOARD_FRAGMENTS */` placeholder (`board.js:57`), so function declarations hoist across fragments and call order at render time is unconstrained. **`web/board-core.js` is the natural home**: it is first in the manifest, it already owns the shared status predicates (`isTerminalResolvedStatus` :261), the shared time formatters (`formatElapsedDuration` :128, `formatRelativeTime` :47, `makeElapsedDurationNode` :150), and the tick refresher (`refreshRelativeTimeNodes` :236). Both consumers (`board-cards.js` at manifest position 3 and `board-detail.js` at position 9) sit after it.

Note the ordering trap already recorded in lessons (`lessons-do-kanban.md`, 0.295.1): manifest position matters for *document listener registration order*, not for function visibility. A summary function is pure; a document listener would not be.

---

## 3. The data the summary needs

| Figure | Where the board gets it today |
|---|---|
| **total grouped REQs** | `userRequest.requestIds.length` — payload field `requestIds` on `generatedUserRequest` (`generate.go:265-272`). Built Go-side from each REQ's `user_request:` upward pointer across queue/working/archive (`model.go:~970-982`), so it is already the complete membership the REQ demands, independent of filters. `renderUserRequestLens` already keeps `requestIds` (full) separate from `shownRequestIds` (filtered) at `board-cards.js:844-845`. |
| **active time spent (completed members)** | Per-request payload fields `hasImplementationSpan` / `implementationSpanMinutes` / `implementationSpanReason` (`generate.go:223-231`). Go already applies the outlier verdict; the client receives no numeric ceiling. **Caveat:** this span is *earliest origin-eligible lifecycle stamp → completed_at*, **not** claim→completion. See §3b. Also gated to terminal SUCCESS only (`generate.go:734-737` — `if isCompletedStatus(ticket.Status)`), so a cancelled member never carries a span. |
| **active time spent (live claimed members)** | `request.claimedAt` (`generate.go:198`, `json:"claimedAt"`) plus `formatElapsedDuration(instantMs, Date.now())` from `board-core.js:128`. Nothing aggregates it today. |
| **estimated active time remaining** | **Missing.** See §3a. |
| **fallback median per effort class** | Already on the payload: `boardData.timeline.projection` (`generate.go:368`, struct at :376-390). See §3c. |
| **successful count/pct** | Derivable client-side from `request.status`, but there is **no JS `isCompletedStatus` helper** — see §4. |
| **resolved count/pct** | `isTerminalResolvedStatus` exists in JS (`board-core.js:261`). |

### 3a. `estimate.p50_active_minutes` — CONFIRMED ABSENT from the board

Grep across the whole repo for `p50_active_minutes` / `P50ActiveMinutes` / `p50ActiveMinutes` in `*.go`, `*.js`, `*.md` returns **only REQ frontmatter in `do-work/`** — zero hits in `skills/do-work-board/`. There is:

- **no** field on `RequestTicket` (`model.go:86-283`);
- **no** read in `parseRequestTicket` (`model.go:707-...`; the only `estimate`-shaped read is `effort_estimate`, a different top-level scalar, at :790);
- **no** JSON key on `generatedRequest` (`generate.go:~120-260`).

The frontmatter schema is `skills/do-work/actions/work-reference.md:125-141`:

```yaml
estimate:
  p50_active_minutes: 75
  confidence: medium            # low | medium | high
  calculated_at: 2026-08-16T12:00:00Z
  basis:
    - Route C
```
with the header comment: *"OPTIONAL informational forecast — never read by scheduling, gating, or pipeline logic, and FROZEN once execution begins. p50_active_minutes is a multiple of 5 and never below 5"*.

**The parser can already reach it.** `parseFrontmatter` returns `map[string]any` from a strict `yaml.Unmarshal`, so `fields["estimate"]` is a `map[string]any` whenever the block parses strictly. What *cannot* reach it is the **salvage path** — `lenientFrontmatterFields` (`frontmatter.go:118-144`), which recovers flat top-level scalars/lists only. Its comment says so at `frontmatter.go:107-113`:

> *"This is a SALVAGE PATH … it is flat and top-level only, so a nested map (`estimate:`) or a block scalar beside the bad line is dropped"*

That comment is **accurate and must stay**. `frontmatter_test.go:357/404/411/432/460` already pins that the strict path preserves the nested `estimate:` map and the salvage path drops it.

Adding the reader needs: a numeric coercion helper (only `coerceScalarToString` at `model.go:2141-2160` and `coerceToStringList` at :2165 exist — **no numeric coercer**), a `RequestTicket` field, a `generatedRequest` field with `omitempty` semantics chosen deliberately (see the `implementationSpanMinutes` comment at `generate.go:216-222`: a genuine zero must not be dropped by `omitempty` while a presence flag ships true — the same trap applies here, though the schema forbids a P50 below 5).

### 3b. The duration outlier rule

- **Constant:** `analysisOutlierCeiling = 4 * time.Hour` — `durations.go:32`.
- **The rule:** `dayMedianExclusionReason(wallSpan time.Duration) string` — `durations.go:323-333`. Returns `"reversed"` (span < 0), `"paused"` (span > ceiling), `""` otherwise.
- **The measurement:** `measureImplementationSpan(ticket *RequestTicket) ImplementationSpan` — `durations.go:222-241`, using `earliestImplementationOrigin` (`durations.go:180-203`) and the exclusion table `implementationSpanOriginExcludedFields` (`durations.go:162-176`: `created_at`, `completed_at`, `release_at`, `status_changed_at`, `blocked_at`, `testing_updated_at`).
- **Human-facing label:** `implementationSpanPausedBadgeText(ceiling)` — `durations.go:38-51`, producing `"over 4h · assumed pause"`. The client receives *this string*, never the number.
- **Where the rule is stated once:** `skills/do-work/actions/estimate-reference.md` → Calibration. `durations.go:20-29` says explicitly *"The rule is stated once … this is its second reader, not a second definition."*
- Per-sample verdict on the day median: `DurationSample.ExcludedFromDayMedian()` — `durations.go:78`.

**Divergence to flag to the planner:** REQ-486 says *"the sum of valid completed **claim-to-completion** spans accepted by the existing duration outlier rule."* The existing rule measures **earliest-origin-to-completion**, deliberately, and `durations.go:229-238` explains why (REQ-505 carried `planning_at` at 16:49 against a 23:00 claim and a 23:01 completion; the card read 1m 23s for six hours of recorded work). Reusing the existing authority means the summary is origin-to-completion, which is *wider* than claim-to-completion where the bookkeeping ran out of order and *identical* where it ran in order. The REQ also forbids re-deriving competing rules in the browser. Either the REQ's wording is loose (most likely) or a second measurement is being asked for; worth naming as an assumption rather than silently choosing.

**Second divergence:** the payload span is emitted only for terminal SUCCESS (`generate.go:734-737`). A `cancelled` member counts toward *resolved* but has no span, and a `failed` member has none either — consistent with the REQ's "disclose as excluded or unavailable", but the reason has to be phrased from what the payload says, not from a status guess.

### 3c. The Timeline forecast median per effort class, and what makes it confident

- **Builder:** `buildTimelineProjection(tickets, aggregate DurationAggregate, now) TimelineProjection` — `timeline.go:310-341`.
- **Constants:** `timeline.go:258-268`
  - `timelineProjectionWindowSize = 60` — the rolling window of most recent in-rule completions.
  - `timelineProjectionMinimumSamples = 5` — *"Below this many samples a median is a coincidence."* Applied **twice**: per bucket (fall back to the window's overall median) and to the window as a whole (decline entirely).
- **Confidence gate:** `timeline.go:321-329` — if `!hasWindowMedian || len(windowMinutes) < timelineProjectionMinimumSamples`, `Confident` stays false and `DeclinedReason` is set to *"only %d completed REQ%s inside the read-time rule; %d are needed before a median means anything"*. Otherwise `projection.Confident = true`.
- **Per-class median:** `timelineBucketMedian(bucketMinutes, windowMedian)` — `timeline.go:368-377`. A bucket with fewer than 5 samples borrows the window's overall median rather than inventing one.
- **Window split:** `timelineProjectionWindow(aggregate)` — `timeline.go:346-364`, skipping `sample.ExcludedFromDayMedian()` so nothing here re-decides what a paused/reversed span is.
- **Per-REQ pick:** `timelineProjectedSpan(ticket, projection)` — `timeline.go:380-390`; `effortMechanical` → `TrivialMedianMinutes`, everything else → `NormalMedianMinutes` (absent `effort_estimate` reads as substantive).

**Already on the payload** — `generatedTimelineProjection` (`generate.go:376-390`, assigned at :940-950), JSON keys: `confident`, `declinedReason`, `trivialMinutes`, `normalMinutes`, `trivialSamples`, `normalSamples`, `windowSamples`, `minimumSamples`, `windowSize`, `chainStart`, `queueEnd`, `rows`, `excluded`. The client already reads exactly these in `renderTimelineForecast` (`web/board-timeline.js:1487-1545`) including the `< minimumSamples` "borrowed median" disclosure at :1536-1540. So the confident-fallback rule is available to the UR summary with **zero Go change** — read `boardData.timeline.projection`.

Reusable formatter: `timelineFormatSpanMinutes(minutes)` — `web/board-timeline.js:207-218` (`"45 min"` / `"2h 15m"` / `"3d 04h"`). It lives in `board-timeline.js` (manifest position 6), which is fine for a caller in `board-cards.js`/`board-detail.js` since declarations hoist across fragments; but if the shared summary moves to `board-core.js` a maintainer may prefer moving/duplicating the formatter deliberately rather than reaching across fragments.

---

## 4. The terminal status sets

### Go

- `isCompletedStatus(normalizedStatus)` — **`model.go:1002-1009`**. Exactly `completed` or `completed-with-issues`. Comment names it *"the Terminal-success status set from actions/work-reference.md's Schema Read Contract"* and explains the exact match is deliberate (a prefix match would let `completed-wth-issues` through).
- `isCancelledStatus` — `model.go:1030-1032`.
- `isTerminalResolvedStatus(normalizedStatus)` — **`model.go:1034-1040`**: `isCompletedStatus(s) || isCancelledStatus(s)`. *"gates completion-time resolution, Recently-done bucketing, and the calendar."*
- `isStoppedStatus` — `model.go:1043-1058`: resolved **plus `failed`**. The comment records that conflating the two cost the Timeline a regression.
- CLI exposure: `frontmatter_cli.go:42` maps the set name `"terminal-resolved"` → `isTerminalResolvedStatus`.
- `normalizeStatus` — `model.go:989-999` (alias collapse; `completed-with-issues` deliberately left as-is).

### Instruction side

`skills/do-work/actions/work-reference.md`:
- **§ Terminal-success status set** — line **259** (heading), body at 261.
- **§ Dependency-source-ready status set** — line 263.
- **§ Terminal-resolved status set** — line **269**, body at 271, boundaries at 273-277.

The resolved section says explicitly: *"any reader that decides whether a `failed` REQ still holds its UR open cites this set by reference and **must not restate or fork the set, or re-derive the rule**, as a competing definition."* That is the rule the REQ's Interfaces section is invoking.

### JavaScript — what already exists to reuse

- **`isTerminalResolvedStatus(status)` — `web/board-core.js:261-263`**, spelled as three literals:
  ```js
  return status === "completed" || status === "completed-with-issues" || status === "cancelled";
  ```
  Its block comment at :256-259 says *"Mirrors model.go's isTerminalResolvedStatus / isCompletedStatus."* Existing callers: `board-core.js:270` (`activeDependentIds`), `board-filters.js:100` (`userRequestHasOpenOrRecentWork`), `board-cards.js:205`, `board-cards.js:526` (recently-done filter).
- **There is NO JS `isCompletedStatus`.** Confirmed by `grep -rn "function isCompletedStatus" web/*.js` → no matches. So "successful = completed + completed-with-issues" has no browser-side helper today. The clean move is to add `isCompletedStatus(status)` beside `isTerminalResolvedStatus` in `board-core.js` and re-express the resolved predicate as `isCompletedStatus(status) || status === "cancelled"` — mirroring the Go composition rather than adding a fourth literal list. That satisfies "do not re-derive competing rules" while keeping the mirror in the one place the comment already points at.
- The Go/JS mirror is tested by slicing the shipped function into probes (`javascript_behavior_a_test.go:1457`, `_b:87`, `_c:483` all call `sliceBalancedBlockAfter(t, indexHtml, "function isTerminalResolvedStatus(")`), plus a stub at `javascript_behavior_b_test.go:1651`. Adding a sibling function means those slices keep working, but a new probe should slice the new function rather than stub it.

---

## 5. The board clock

- **The interval:** `web/board.js:68` — `setInterval(refreshRelativeTimeNodes, 1000);`, at the end of the boot block inside the IIFE. One shared 1-second ticker for the whole page.
- **What it updates:** `refreshRelativeTimeNodes()` — `web/board-core.js:236-253`. It selects `document.querySelectorAll("[data-instant-ms]")` and, per node, recomputes the label from `Number(node.dataset.instantMs)` against a single `Date.now()` captured at the top of the call. `dataset.tickFormat === "duration"` picks `formatElapsedDuration` (stopwatch, `board-core.js:128`), otherwise `formatRelativeTime` ("6min ago", `board-core.js:47`). Writes only when the text changed, and calls `syncClockSkewTitle` for duration nodes.
- **How a node opts in:** by carrying `data-instant-ms` (and optionally `data-tick-format="duration"`). `makeRelativeTimeNode` (`board-core.js:68-77`) and `makeElapsedDurationNode` (`board-core.js:150-160`) are the two builders that set it. Composites: `makeInstantWithRelativeNode` (:81), `makeInstantWithStopwatchNode` (:165).
- **How a node opts out deliberately:** `makeImplementationSpanNode` (`board-cards.js:66`) builds a *plain* `span.elapsed-duration` with no `data-instant-ms`, and the comment at `board-cards.js:54-59` explains why — a finished span must not be rewritten every second.
- **There is no subscription API.** No registry of tick callbacks exists; the only extension point is the attribute selector. A UR summary whose *live claimed contribution* must tick has two honest options:
  1. Emit the live part as its own `[data-instant-ms][data-tick-format="duration"]` node, so the existing ticker rewrites it with no new machinery — but that only ticks a *single* claim's elapsed time, not a *sum* of a completed base plus N live claims.
  2. Add a small tick-subscriber list in `board-core.js` that `refreshRelativeTimeNodes` fans out to (or a sibling `refreshDerivedSummaries()` called from the same `setInterval`), so header and drawer recompute from one `Date.now()`. This is the shape the REQ's "refresh through the board's existing clock so the header and drawer cannot drift from the claimed card stopwatch" is asking for — the drift risk is exactly two `Date.now()` calls in different frames.

  Option 2 is the one that satisfies the requirement literally. Keep the single captured `nowMs` per tick, as `refreshRelativeTimeNodes:237` already does.
- Skew handling to reuse rather than re-derive: `futureInstantSkewAllowanceMs = 2 * 60 * 1000` (`board-core.js:105`, mirrors Go's `futureTimestampSkewAllowance`), `clockSkewMarkerText = "⚠ clock skew"` (:106). `formatElapsedDuration` returns the marker for a future instant instead of a clamped `0s` (:129-131). A claimed member with a future `claimed_at` therefore has **no usable live elapsed time** and must be disclosed, not counted as zero.
- Unrelated second timer, do not confuse: `web/board-timeline.js:2738-2739` runs its own 50ms `setInterval` for plot-width change detection.

---

## 6. Test surface

### The Node behaviour lane (the one that matters here)

Harness, all in `generate_test.go`:
- `generateLiveSiteInDir(t)` :335-360 — builds the board against the **real `do-work/` tree** with a stubbed git lookup, writes the static site to a temp dir.
- `generateLiveSite(t)` :361-374 — cached (`liveIndexOnce`) `index.html` string.
- `sliceBalancedBlockAfter(t, sourceText, anchorToken)` :1998-2019 — brace-matches one function out of the assembled page so the probe drives the **shipped** source, not a copy.
- `runJavaScriptBehaviorProbe(t, probeName, probe)` :276-293 — pipes the probe to `node -` on **stdin** (an `-e` argument exceeds Linux's 128 KiB per-arg limit; macOS would pass and CI would fail). Increments `javaScriptBehaviorProbeCount` for the strict lane's zero-probe guard.

**There is no synthetic-queue fixture builder for the JS probes.** Each probe hand-writes a `var boardData = {...}` object literal plus a minimal DOM stub inside the probe string. The stub in the UR probe is `makeNode()` — `javascript_behavior_a_test.go:1495-1519` — a plain object with `childNodes`, `dataset`, `attributes`, `listeners`, `appendChild`, `removeChild`, `setAttribute`, `getAttribute`, `addEventListener`, and a `dispatch(eventName)` that invokes registered handlers synchronously. `document` is a two-method stub (:1521-1524). `makeRequestCard` is stubbed to `{ className: "req-card", requestId: requestId }` (:1525).

**Probes covering `renderUserRequestLens`:**

| Test | File:line | What it pins |
|---|---|---|
| `TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened` | `javascript_behavior_a_test.go:1453-1650` | The REQ-236 contract in full: folded renders one row per UR with zero cards and `aria-expanded="false"`; exactly one drawer trigger per row; opening reveals only that UR's filtered cards; siblings stay empty; re-activating collapses; Active + `status:pending` hides the same URs in both readings; **the By UR head carries no `aria-expanded` at all** (:1645-1649). Fixture: 3 URs, 4 REQs (`pending`, `completed`, `claimed`, `completed`), a 2-entry `calendar`, `Date.now` stubbed to `2026-08-15T12:00:00Z`. |
| `TestJavaScriptBehaviorByUserRequestLensCountsRecentlyDoneAsActive` | `javascript_behavior_b_test.go:75` | The Active scope + recently-done window. Its preamble declares `var userRequestCardsFolded = false;`. |
| `TestJavaScriptBehaviorByUserRequestLensEmptyState` | `javascript_behavior_b_test.go:139` | The three empty-state branches. |
| `TestJavaScriptBehaviorByUserRequestLensUsesRecentWindowAtCaller` | `javascript_behavior_c_test.go:479` | The window is read at the caller, not baked in. Slices `isTerminalResolvedStatus` at :483. |
| `TestJavaScriptBehaviorTestingStatusUpdateInvalidatesUserRequestLens` | `javascript_behavior_b_test.go:180` | `renderedOnce.userRequestLens` invalidation. |
| `TestGenerateOffersThreeLensButtons` | `generate_test.go:2773-2786` | Markup/source tokens only: the three `data-lens-target` values, `data-ur-cards="folded"`, `URs&nbsp;only`, and the `applyLensSelection(...)` call site verbatim. |

Shared decode type: `renderedUserRequestRow` — `generate_test.go:2791-2796` (`userRequestId`, `expanded`, `cardIds`, `drawerTriggers`), plus `userRequestIdsOf` at :2803.

**Probes covering the drawer:**

- **`openUserRequestDetail` has NO JavaScript behaviour probe.** `grep -rn "openUserRequestDetail" *_test.go` → zero hits. Nothing pins the UR drawer's meta rows.
- The drawer probes that exist drive the REQ path: `TestJavaScriptBehaviorDetailRendersOnlyObservedPhaseBreakdown` (`javascript_behavior_a_test.go:2231`) slices `appendPhaseBreakdownRows` and stubs `appendMetaRow` at :2266; `TestJavaScriptBehaviorDrawerHeadingDeduplication` (`javascript_behavior_b_test.go:13`).
- `generate_test.go:1324` asserts the source token `appendMetaRow("Error", request.error)` — a source-token style check, not behaviour.

**Browser probes.** Harness in `browser_probe_test.go` (DevTools pipe straight to the browser binary; result written into `#queue-kanban-probe-result`; `QUEUE_KANBAN_BROWSER` overrides the binary, `QUEUE_KANBAN_BROWSER_PROBES` / `QUEUE_KANBAN_STRICT_BROWSER_BEHAVIOR` gate the lane). Existing:
- `TestBrowserBehaviorUserRequestCopyAllIncludesGroupedRequests` — `user_request_clipboard_browser_probe_test.go:44`. **This is the one browser probe that drives the real By-UR UI and opens a UR drawer.** Its fixture builds `*UserRequestTicket` / `*RequestTicket` structs directly (`userRequestCopyFixture` at :28-37; `boardColumnCopyFixtureTicket`) with members in `queue`, `working`, and `archive` `TreeSection`s and statuses `pending` / `claimed` / `completed`. Closest existing thing to a synthetic UR fixture at the Go level.
- `TestBrowserBehaviorDrawerTicketTitlesAndGlossary` — `browser_probe_test.go:794`.
- `TestBrowserBehaviorMarkLabelTextExtent` :655, `TestBrowserBehaviorProbeResultKeepsLiteralText` :1162.

Go-side synthetic tree: `createSyntheticDoWorkTree(t)` / `syntheticBoard(t)` — `board_synthetic_test.go:25` / :68. Builds a deterministic repo-independent `do-work/` tree with UR linkage, both archive shapes, and >=900 tickets. This is the right base for a **Go-level** fixture that carries `estimate:` blocks and varied statuses; it does **not** feed the Node probes.

### Gap analysis against the REQ's RED case

The RED case asks for a fixture with members covering: completed, completed-with-issues, cancelled, pending, claimed, blocked, failed, missing-timestamp, outlier-span, saved-P50, missing-P50-with-confident-history, insufficient-history.

- **Statuses:** trivially expressible in the existing inline `boardData.requests` literal — the probes already carry `pending`/`completed`/`claimed`. Adding `completed-with-issues`, `cancelled`, `blocked`, `failed` is more object literal, no new helper.
- **Missing timestamps / outlier span:** the client never re-derives the outlier verdict; it reads `hasImplementationSpan` / `implementationSpanMinutes` / `implementationSpanReason`. So the probe fixture expresses these as payload shapes (`hasImplementationSpan: false`, or `implementationSpanReason: "paused"` / `"reversed"`), not as raw stamps. That means **the outlier rule itself is exercised on the Go side** (`durations_test.go`), and the JS probe only pins the disclosure branches. Worth being explicit about in the plan: two lanes, two questions.
- **Saved P50:** needs the new payload field to exist first — a Go-side test that a REQ carrying `estimate: {p50_active_minutes: N}` reaches `generatedRequest`, plus a strict-vs-salvage parse test (extend the existing `frontmatter_test.go:357-460` family rather than starting a new one).
- **Missing-P50-with-confident-history / insufficient-history:** expressed in the probe by stubbing `boardData.timeline.projection` with `confident: true` + `normalMinutes`/`trivialMinutes`/`normalSamples`/`minimumSamples`, and a second run with `confident: false` + a `declinedReason`. No new fixture machinery.
- **Stubbed clock:** already the idiom — `Date.now = function () { return Date.parse("…"); };` (`javascript_behavior_a_test.go:1471`). Advancing it mid-probe is just a reassignment between calls, but the probe must also call whatever tick entry point the implementation adds, since `setInterval` never runs in the probe.
- **The drawer half needs a new probe from scratch** — there is no `openUserRequestDetail` probe to extend, and the drawer's DOM stub needs `#detail-meta`, `#detail-body`, `#detail-glossary`, `#detail-kind`, `#detail-id`, `#detail-drawer-title`, `#detail-drawer`, `#detail-resizer`, `#detail-copy`, `#detail-copy-all` (all read at `board-detail.js:291-300`). That, plus the module-level `var`s, is the real cost of the drawer RED case.
- **Both surfaces report identical values:** the cleanest RED is a probe that slices the shared summary function once and asserts both consumers call it with the same UR id and render the same numbers — rather than two independent assertions that could pass while disagreeing.
- **Two-theme / narrow-width render evidence** is browser-lane work; note the maintainer's memory that the browser lane **skips silently** unless `QUEUE_KANBAN_BROWSER` is set, and a skip is not a pass.

---

## 7. Board guide and stale comments

### The board guide

`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/docs/board-guide.md` — 72 lines.

- **Line 25** describes the toolbar: *"a **Lens** toggle (flat Columns vs. grouped **By UR**)"* and the Active/All scope. It says **two** lens choices.
- **`grep -n "URs only"` over the guide, `SKILL.md`, `actions/board.md`, and `_dev/primes/prime-kanban-board.md` → zero matches.** The third lens shipped in REQ-236 and the guide was never updated. So the guide is already stale by one feature before this REQ adds two more fold surfaces and a summary.
- **Line 52-54, "## Card drawer"** describes the REQ drawer's rows and Copy behaviour. It does not describe the UR drawer's `Grouped REQs` / `REQ ids` rows at all.

### The stale nested-estimate comment

**`skills/do-work-board/tools/queue-kanban/timeline.go:392-396`**, the doc comment on `timelineChainStart`:

```go
// timelineChainStart is when the first unstarted REQ can begin: after whatever is
// already running. The in-flight REQ's own `estimate:` block would be the better
// offset for exactly this bar, but the board parses no nested frontmatter blocks,
// and adding that surface for one bar is the sophistication this REQ trades for a
// stated assumption.
```

Line **394** carries the exact claim: *"the board parses no nested frontmatter blocks"*. This is the comment the REQ's Interfaces section means. It becomes false the moment the P50 reader lands.

Two notes for whoever edits it:
1. The REQ says *"Do not change the Timeline's scheduling or forecast behavior merely because the board begins exposing the saved P50 value."* So the edit is to the **comment only** — `timelineChainStart` keeps using the projection median. The rewritten comment should say the board now reads the block but that this bar deliberately does not use it, which is a stronger statement than the current one and still honest.
2. **`frontmatter.go:109` is a different claim and must NOT be changed:**
   ```go
   // top-level only, so a nested map (estimate:) or a block scalar beside the bad
   ```
   That sentence is about the **salvage path** for malformed frontmatter and stays true. It is pinned by `frontmatter_test.go:266`, `:404`, `:411`, `:432`.

`grep -rn "nested"` over `*.go`/`*.js`/`*.md` in the tool returns only these two plus test comments — no third copy of the claim, and nothing in the guide or either prime.

---

## 8. Primes and lessons

### `_dev/primes/prime-kanban-board.md` (maintainer-side; export-ignored; nothing shipped may cite it)

Entries bearing on this REQ:

> **"Keep the parser in lock-step with the schema."** *The board buckets tickets by the `status` vocabulary in `skills/do-work/actions/work-reference.md`; its parsed display fields must stay aligned with `…/model.go`, and the Testing placeholders with `…/testing.go`.*

> **"Write surfaces are counted here."** *The tool has exactly three write surfaces, and none touches pipeline state… Adding a fourth write surface means amending this sentence in the same commit; `_dev/tests/contract-regressions.sh` pins the count to this file.* — REQ-486 adds none; the count stays three.

> **"A chart's correctness is partly a claim about pixels — generate a board and look at it."** *A passing suite is not evidence about two glyphs sharing a coordinate: REQ-226, REQ-231, REQ-237 and REQ-240 each shipped a defect that every assertion passed over and a render made obvious. Measure `getBoundingClientRect()` intersections in the live DOM when the question is "do two things overlap"; read the rendered text when the question is "what does this say".* — directly the REQ's "no collisions or clipped metrics at normal and narrow widths".

> **"Render evidence must name the page it measured, in the same call that measures it."** *Where several agents share one browser instance, a sibling can navigate it between your navigate and your evaluate… Return `location.href` alongside every measurement, or drive an isolated browser; a URL checked *before* navigating is not the same claim.*

> **"The surface behind this board's SVG is `<body>`, not any `--surface-*` token."** *…two ink tokens picked for a text hierarchy can be nearly indistinguishable as fills: `--ink-faint` against `--ink-soft` measures 1.29:1 light and 1.82:1 dark, which is not a channel and cannot carry a row distinction.* — relevant if the summary uses `--ink-faint` (as `.ur-count` does) to distinguish approximate/unavailable markers from real values. It cannot carry that distinction alone.

> **Browser support:** *the strict browser lane targets current stable Chromium. Chrome 141 is deprecated and is not a compatibility target (REQ-375).*

### `_dev/primes/lessons-kanban-board.md` (a link index; `slugged: partial`)

Relevant entries, quoted:

> - [REQ-236: a document-wide DOM count is only evidence when its root is the thing under test](../../do-work/archive/UR-054/REQ-236-add-a-urs-only-lens-to-the-board.md#lessons-learned) — line 31. **The direct predecessor.**
> - [REQ-233: a programmatic `.focus()` cannot answer a `:focus-visible` question, and will say the ring is broken when it is not] — line 32. Keyboard-affordance testing.
> - [REQ-338: a roving tabindex over a VIRTUALIZED list needs a clamp into the rendered range, not an exact match — matching alone marks nothing tabbable and takes the whole list out of the Tab order; every other row needs an explicit `tabindex="-1"` because a focusable element without the attribute is still a stop; and **Tab cannot be tested with synthetic events at all**, since its focus movement is a trusted-input default action] — line 10. Bears on "the fold control is keyboard-operable": a synthetic-event probe can prove `aria-expanded` flips on click/Enter dispatch, but not that the control is in the Tab order. Both fold controls being real `<button>`s is what actually delivers that.
> - [REQ-374: a fixture that spans a threshold widely does not test it — 40min/18h/−3h against a four-hour ceiling let a SECOND ceiling of six hours pass the agreement test silently… only a pair straddling the real boundary, derived from the constant, catches a second definition. Also: a property argued at length in a comment is not a property under test] — line 8. **Directly on the outlier rule this REQ reuses.**
> - [REQ-245: asserting a phrase is absent is not a guard — it passes when the whole string is replaced] — line 26. Bears on the "never present a partial sum as complete" assertions: assert the *presence* of the qualifier, not the absence of a number.
> - [REQ-304: a missing-branch fix needs a fixture that can fail in both directions — a reversed-only case passes against `if (true)`; and a mirrored branch anchors to its own span, not to the one it mirrors] — line 55.
> - [REQ-305: a probe that calls the function under test directly cannot hold its call site — five mutations of the copy were caught and the one that reverted the real defect passed clean; and changing a parameter's type silently re-points every existing call, because `[]` is truthy] — line 56. **Directly on "one shared summary function consumed by two surfaces"**: a probe that calls the summary function alone proves nothing about the header and drawer calling it.
> - [REQ-291: `getBBox()` returns zeros for an unrendered element, so a browser probe's default failure is a successful-looking measurement of nothing — write the result node last and assert positive-finite-and-known-font-size, never merely "no error"] — line 51.
> - [REQ-228: a forecast may be wrong about timing, never about order] — line 42. The Timeline projection's own lesson.
> - [REQ-323: …also, a guard no mutation can break is dead before it is untested, which is how the collapse function's `isFinite` branch turned out to be unreachable] — line 13.

### `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (ships; `slugged: partial`)

Relevant entries, quoted (trimmed where marked):

> - **0.295.1:** *a surface that reads the drawer's state back has to be TOLD when that state changes; it cannot guess when. The Activity row highlight read `currentDetailKind`/`currentDetailId` from a document click listener in `board-activity.js` — manifest position 7, ahead of the `[data-detail-kind]` delegation in `board-controls.js` at position 10 — so on every click the drawer still named the previous ticket… The fix was a deletion: `setDetailTarget` beside the two variables as their only writer, called by every open and close, and both document listeners gone. **A probe that assigns those variables by hand in the order the listener assumed cannot see this class of bug at all — drive the shipped writer.*** — **directly on the drawer**, and on manifest-position ordering for anything listener-shaped.
> - **0.107.0:** *`cancelled` (won't-do, written by `do-work abandon`) is terminally resolved but not terminally successful — it shares Recently-done via `isTerminalResolvedStatus`, while `isCompletedStatus` stays the success-only gate. **Keep the two classifiers distinct; success-readers must never widen to the resolved set.*** — exactly the successful-vs-resolved split this REQ needs.
> - **0.218.0:** *the calendar carries EVERY REQ, one entry each… `boardData.calendar` is no longer a list of completions: `recentlyDoneIds` (`web/board-cards.js`) filters it on `isTerminalResolvedStatus(entry.status)`, and **any new consumer must gate the same way** or it will count claimed and failed work as done.*
> - **0.133.0:** *a `*_at` stamp parsing past the board's `now` + 2min skew (`futureTimestampSkewAllowance`, mirrored by `futureInstantSkewAllowanceMs` in `web/board-core.js` — **change together**) marks the ticket `FutureTimestampFields` → "future stamp" badge + data warning; **the stopwatch renders "⚠ clock skew" instead of a clamped-frozen "0s"**…* — the live-claimed contribution inherits this: a future `claimed_at` is unavailable, not zero.
> - **Failure-detail display:** *a schema default is not an instruction to invent missing data — `error_type` resolves an invalid **present** value to `code` while retaining its original value, unrecognized flag, and data warning, but **an absent field remains empty through parse, projection, and drawer rendering**… Keep both assertions at parse level and projection level: a normalizer-only test cannot prove that the board's actual read path is wired, and a source-token frontend test pins the display seam without coupling Go tests to a browser runtime.* — the template for the P50 reader's tests.
> - **REQ-032/034:** *…**An empty set never overlaps — on the board absence means unknown, which must not render as conflict.***
> - **REQ-117:** *`board.Warnings` is a **free UI channel**… **a stated *reason* in a comment is a factual claim and reviews must check it like any other** — "the board has no warning channel for it" shipped at 98% and was three greps from disproof.* — the same shape as `timeline.go:394`'s soon-to-be-false claim.
> - **[family: paired-predicate-drift] 0.294.2 (2026-09-05):** *two readers of one contract drift silently, because each one passes its own tests… The board already parsed the field the rule needs, so the fix was a predicate, not a parser: `isDependencySourceReadyStatus` beside `isCompletedStatus`, **named after the contract section it implements**. When a shared vocabulary section in `work-reference.md` gains an arm, grep for every predicate that spells out the old arm list; **a reader that names one arm of a two-arm set is drift even when it compiles and passes.***
> - **[family: subject-not-restated-in-detail] REQ-588 (2026-09-05):** *when a payload gains a structured field, the prose that carried the same fact must stop carrying it, and every test that read the fact out of the prose must move to the field.* — applies to `ur-count` if the new summary subsumes the `N REQ` string.
> - **REQ-116:** *…**update the field-specific action explanation whenever a new field gains a board role** rather than creating another competing list.* — `work-reference.md`'s `estimate:` block comment currently says the field is "never read by scheduling, gating, or pipeline logic"; a *display-only* board read does not contradict that, but the sibling display-only fields (`impact:`, `effort_estimate:`, `sweep:`) each carry an explicit "parsed by `…/model.go` into a card chip and a drawer row… keep that parser in lock-step with this line, both changing in the same commit" sentence. The P50 reader should earn the same sentence.
> - **REQ-089 / 0.193.4 / REQ-447 / REQ-458 / 0.275.3** — not bearing on this REQ.

---

## Assumptions and open points a planner should settle

- **A1.** "Claim-to-completion span" in the REQ vs. the shipped rule's earliest-origin-to-completion (`durations.go:222-241`). Reusing the existing authority as the Interfaces section demands means origin-to-completion. Recommend stating that in the plan rather than building a second measurement.
- **A2.** The completed-span payload is emitted for terminal SUCCESS only (`generate.go:734-737`). Cancelled and failed members carry no span at all — which is fine for the sum, but the "excluded / unavailable" copy must be derived from `hasImplementationSpan` plus status, not from a re-derived rule.
- **A3.** The board clock has no subscriber API. A tick fan-out in `board-core.js` beside `refreshRelativeTimeNodes`, sharing its one `nowMs`, is the smallest change that satisfies "cannot drift from the claimed card stopwatch".
- **A4.** No JS `isCompletedStatus` exists. Add it in `board-core.js` beside `isTerminalResolvedStatus` and recompose the resolved predicate from it, matching the Go composition and the 0.294.2 lesson.
- **A5.** `.ur-group-head` is a nowrap flex row with an ellipsised `flex:1` title. Metrics need `flex-wrap: wrap` (or a second row) or they will squeeze the title rather than wrap.
- **A6.** The board guide's Lens sentence (line 25) is already stale by one lens. This REQ's guide edit should fix that too, or say why it did not.
- **A7.** No probe exists for `openUserRequestDetail`. The drawer half of the RED case is a new probe with a ten-element DOM stub, not an extension of an existing one.
