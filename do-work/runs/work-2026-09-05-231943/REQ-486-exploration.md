u# REQ-486 exploration — collapsible UR groups and progress summaries

Repo: `/home/user/skill-do-work` at HEAD `57c0edb9f70a7da1c5b526985393313aaa9194bb`. Read-only pass; nothing edited.
Board tool root: `/home/user/skill-do-work/skills/do-work-board/tools/queue-kanban/`. Paths below are relative to that root unless stated otherwise.

**Route: C (plan first).** Reasoning in §11.

---

## 1. Stale baseline — what changed since the prior exploration

A prior exploration exists at `/home/user/skill-do-work/do-work/runs/work-2026-09-05-170806/REQ-486-exploration.md` (442 lines). Its structural analysis is sound and I reused it as a starting point. Six of its statements no longer hold at HEAD.

| # | Prior statement | Status at HEAD |
|---|---|---|
| S1 | Repo root is `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` | **Wrong tree.** This session's repo is `/home/user/skill-do-work`. Every absolute path in the prior report is unusable. |
| S2 | `renderUserRequestLens` is `web/board-cards.js:819-975` | **Moved.** It is `web/board-cards.js:925-1081` — a uniform +106 shift. Every inner line the prior report cited shifts by the same amount. Structure is unchanged. |
| S3 | "grep for `p50_active_minutes` … returns **zero hits** in `skills/do-work-board/`" | **No longer true as written.** `frontmatter_test.go:341` and `frontmatter_test.go:492` now carry the literal in fixture frontmatter. The claim that matters is still true: **no production reader exists** — no `RequestTicket` field, no read in `parseRequestTicket` (`model.go:707`), no `generatedRequest` JSON key. |
| S4 | "The recorded implementation commit is `456ee9d`" (the REQ says this too) | **Unreachable.** `git rev-parse 456ee9d` → `fatal: unknown revision`. `.git/shallow` exists; `git rev-list --count HEAD` = 103; the earliest reachable commit is `8e58a79` (2026-09-05). The prior report's recipe `git diff 52fb7d72..456ee9d` cannot run here. Use the archived record instead (§3). |
| S5 | The REQ is deferred behind REQ-510 | **Satisfied.** `do-work/archive/UR-098/REQ-510-sweep-work-reference-sections-owned-by-cli-tests.md` carries `status: completed`. The `depends_on: [REQ-510]` edge in REQ-486's frontmatter no longer blocks selection. |
| S6 | Assorted Go and CSS anchors | **Shifted a few lines each.** `isCompletedStatus` `model.go:1002`→**1007**; `isTerminalResolvedStatus` `1034`→**1038**; `isStoppedStatus` `1043`→**1057**; `coerceScalarToString` `2141`→**2131**; `coerceToStringList` `2165`→**2155**; fragment manifest `generate.go:43-56`→**43-55**; `renderTimelineForecast` `board-timeline.js:1487`→**1478**. All of `board.css`'s UR block shifted +188 (see §6). |

Everything the prior report said about `board-detail.js` (`:291-300`, `:614-636`), `board-core.js`, `board.js:68`, `durations.go`, `timeline.go` and the test harness is **unchanged at HEAD**, verified line by line below.

---

## 2. The request's own RED claims, verified

The REQ's "Why RED now" makes four factual claims. All four hold.

| Claim | Evidence at HEAD |
|---|---|
| By UR headers are drawer triggers **without** `aria-expanded` | `web/board-cards.js:991-994` — the non-folded branch sets `head.dataset.detailKind = "ur"` and `head.dataset.detailId`, and never touches `aria-expanded`. `javascript_behavior_a_test.go:1652` currently **fails the build** if it ever does: *"by-UR head carries aria-expanded=%q; the fold must not leak into the always-open lens"*. |
| By UR **always appends every REQ card** | `web/board-cards.js:1038-1041` — `group.appendChild(head); group.appendChild(makeUserRequestCards());`, eager, unconditional. |
| The drawer **always renders its REQ-id list directly** | `web/board-detail.js:625-629` — `appendMetaRow("Grouped REQs", String(requestIds.length))` then `appendMetaRow("REQ ids", makeTicketLinkList(requestIds))`, with no control between them. |
| Generated requests **do not expose the saved P50**, and neither surface has a shared rollup | §5 and §4 below. |

---

## 3. The prior UR-lens implementation (REQ-236)

`456ee9d` is unreachable (S4). The authoritative record at HEAD is `/home/user/skill-do-work/do-work/archive/UR-054/REQ-236-add-a-urs-only-lens-to-the-board.md`. Its declared write set was five files: `web/template.html`, `web/board-cards.js`, `web/board-controls.js`, `web/board.css`, `generate_test.go`.

Decisions from that record that still govern this REQ:

- **D-01** — `URs only` is a **card-fold modifier on the By UR lens**, not a third `viewState.lens` value. Four files ask "is the UR lens on screen?" by naming `"user-request"`; keeping `viewState.lens === "user-request"` true for both readings kept all of them correct unedited. `viewState.lens` is still declared as `"flat" | "user-request"` at `web/board.js:26`.
- **D-02** — the fold flag lives in `board-controls.js`, not `viewState`. Confirmed: `var userRequestCardsFolded = false;` at `web/board-controls.js:8`, set only by `applyLensSelection` (`web/board-controls.js:115-123`).
- **D-03** — expanded/collapsed lives in the DOM, never persisted. A re-render rebuilds rows collapsed. `applyLensSelection` sets `renderedOnce.userRequestLens = false` before `applyLens()` (`web/board-controls.js:119`) so the cached lens is always dropped when the reading changes.
- **D-05** — helpers stay inner functions of `renderUserRequestLens`, closing over what they need. That is what let REQ-236's probe slice one anchor and get a behavioral RED.

### 3.1 `renderUserRequestLens` at HEAD — `web/board-cards.js:925-1081`

```
 925  function renderUserRequestLens() {
 926    var host = document.getElementById("user-request-lens");
 927    host.textContent = "";                       // full rebuild every call
 933    var recentlyDoneIdSet = {};                  // built once per render, not per UR
 938    (boardData.userRequestOrder || []).forEach(function (userRequestId) {
 950      var requestIds = userRequest.requestIds || [];        // FULL membership
 951      var shownRequestIds = requestIds.filter(requestMatchesFilters)  // FILTERED
 954-963  Active-scope gate → hiddenResolvedCount, early return
 964-966  hasActiveFilters() && shownRequestIds.length === 0 → early return
 968      var group = createElement("section", "ur-group");
 974      function makeUserRequestCards() { ... shownRequestIds ... }   // inner
 982      var head = createElement("button", "ur-group-head");
 984      if (userRequestCardsFolded) {
 989        head.setAttribute("aria-expanded", "false");
 990        head.appendChild(createElement("span", "ur-fold-marker", "▸"));
 991      } else {
 992        head.dataset.detailKind = "ur";           // By UR: the head IS the drawer trigger
 993        head.dataset.detailId = userRequestId;
 994      }
 995      head.appendChild(createElement("span", "ur-id", userRequestId));
 996      head.appendChild(createElement("span", "ur-title", ...));
 997-1005 optional "cites" badge, optional "no input.md" chip
1006-1014 head.appendChild(span.ur-count)   →  "12 / 43 REQ"  or  "43 REQ"
1015      if (userRequestCardsFolded) {
1016        var headRow = createElement("div", "ur-group-row");
1018        var detailButton = createElement("button", "ur-group-detail", "Details");
1026        var openCards = null;
1027        head.addEventListener("click", function () { ...toggle... });
1038      } else {
1039        group.appendChild(head);
1040        group.appendChild(makeUserRequestCards());   // always-open By UR reading
1041      }
1043      host.appendChild(group);
1044    });
1046-1080  empty-state text + "n URs hidden — switch URs to All" note
```

### 3.2 How the two readings differ — the whole delta

| | By UR (`userRequestCardsFolded === false`) | URs only (`=== true`) |
|---|---|---|
| head element | `button.ur-group-head`, **is** the drawer trigger (`data-detail-kind="ur"`, :992-993) | `button.ur-group-head`, **is** the fold toggle; no `data-detail-*` |
| `aria-expanded` | **absent entirely** | `"false"` initially (:989), flipped by the click listener at :1031 / :1036 |
| fold marker | none | `span.ur-fold-marker` "▸" (:990), rotated 90° by CSS at `board.css:1576-1579` |
| wrapper | head appended straight into `section.ur-group` (:1039) | head wrapped in `div.ur-group-row` beside `button.ur-group-detail` (:1016-1024) |
| cards | eager at :1040 | lazy on first click (:1034), `group.removeChild` on collapse (:1029) |
| host class | — | `applyLens` adds `.is-folded` to `#user-request-lens` (`board-controls.js:132`) |

**The Details-button collision was already solved once, in the folded branch.** The head cannot be both a fold control and a drawer trigger, so `data-detail-kind`/`data-detail-id` move to a sibling `button.ur-group-detail` inside `div.ur-group-row` (`board-cards.js:1018-1023`). Both are real `<button>`s, so both are keyboard-operable, and the delegated `[data-detail-kind]` handler in `board-controls.js` needs no change — it finds the attribute on a different node. **REQ-486's By UR fold is the same move applied to the other branch**, with the initial state inverted. That is the smallest honest implementation, not a new pattern.

Template at HEAD: `web/template.html:118-136` — three lens buttons in `#lens-group`; the third carries `data-lens-target="user-request" data-ur-cards="folded"`, with an explanatory comment at :126-127.

---

## 4. Where the shared summary function belongs

**`web/board-core.js`.** Three reasons, all checkable:

1. **Manifest position.** The client is assembled from `boardJavaScriptFragmentPaths` — `generate.go:43-55` — in this order: `board-core`, `board-filters`, `board-cards`, `board-calendar`, `board-durations`, `board-timeline`, `board-activity`, `board-testing`, `board-detail`, `board-controls`, `board-clipboard`. `board-core.js` is first; both consumers (`board-cards.js` at position 3, `board-detail.js` at position 9) sit after it.
2. **Scope.** All eleven fragments inline into the single IIFE in `web/board.js` at the `/* INLINE_BOARD_FRAGMENTS */` placeholder (`web/board.js:57`). Everything the summary needs is declared **above** that placeholder: `boardData` (:11), `requestsById` (:18), `userRequestsById` (:19). Function declarations hoist across fragments, so call order at render time is unconstrained.
3. **Neighbours.** `board-core.js` already owns the shared status predicate `isTerminalResolvedStatus` (:261-263), the shared time formatters `formatRelativeTime` (:47) / `formatElapsedDuration` (:128) / `makeElapsedDurationNode` (:150), `createElement` (:3), and the tick refresher `refreshRelativeTimeNodes` (:236-253).

One caveat recorded in `lessons-do-kanban.md:15` (0.295.1): manifest position matters for **document-listener registration order**, not for function visibility. A pure summary function is unaffected. A document listener would not be — so do not add one.

Proposed name, satisfying `coding-guardrails.md` § 5 Naming for Reach (two words minimum, findable by plain-text search): **`summarizeUserRequestProgress(userRequestId, nowMs)`**. The `nowMs` parameter is not decoration — it is what makes the header and the drawer provably read the same instant (§7).

---

## 5. Can the board already read the nested P50? No.

### 5.1 What is missing

- **No `RequestTicket` field.** `model.go:86` opens the struct; the only `estimate`-shaped read is `effort_estimate`, a different top-level scalar, at `model.go:790`.
- **No read in `parseRequestTicket`** (`model.go:707`).
- **No `generatedRequest` JSON key** (`generate.go:135-263`).
- **No numeric coercer.** `model.go` carries only `coerceScalarToString` (:2131) and `coerceToStringList` (:2155). A nested numeric read needs a new one.

### 5.2 What already works in the parser's favour

`parseFrontmatter` returns `map[string]any` from a strict `yaml.Unmarshal`, so `fields["estimate"]` is already a `map[string]any` whenever the block parses strictly. This is **pinned by existing tests**: `frontmatter_test.go:403-404` (*"nested estimate: map is gone — the block fell to the last-resort recovery instead of parsing strictly"*), :432, :449, :506-507.

What cannot reach it is the **salvage path** — `lenientFrontmatterFields` (`frontmatter.go:127-145`), whose contract comment at `frontmatter.go:107-113` says so:

> *"This is a SALVAGE PATH, not a contract a writer may aim at. … it is flat and top-level only, so a nested map (estimate:) or a block scalar beside the bad line is dropped"*

**That comment is accurate and must NOT be changed.** It is a different claim from the stale one in §5.4.

### 5.3 The schema

`skills/do-work/actions/work-reference.md:125-139`:

```yaml
# OPTIONAL informational forecast — never read by scheduling, gating, or pipeline
# logic, and FROZEN once execution begins. p50_active_minutes is a multiple of 5 and
# never below 5 …
estimate:
  p50_active_minutes: 75
  confidence: medium
  calculated_at: 2026-08-16T12:00:00Z
  basis:
    - Route C
```

A display-only board read does not contradict *"never read by scheduling, gating, or pipeline logic"*. But its three sibling display-only fields — `sweep:` (:119), `impact:` (:121), `effort_estimate:` (:123) — each carry an explicit clause: *"parsed by `../../do-work-board/tools/queue-kanban/model.go` into a card chip … and a drawer row … keep that parser in lock-step with this line, both changing in the same commit."* The P50 reader earns the same sentence. This is `lessons-do-kanban.md:35` (REQ-116): *"update the field-specific action explanation whenever a new field gains a board role rather than creating another competing list."*

### 5.4 The stale comment the REQ's Interfaces section means

`timeline.go:392-396`, the doc comment on `timelineChainStart` (`timeline.go:397`):

```go
// timelineChainStart is when the first unstarted REQ can begin: after whatever is
// already running. The in-flight REQ's own `estimate:` block would be the better
// offset for exactly this bar, but the board parses no nested frontmatter blocks,
// and adding that surface for one bar is the sophistication this REQ trades for a
// stated assumption.
```

**Line 394** carries the claim that becomes false. Two notes for whoever edits it:

1. The REQ says *"Do not change the Timeline's scheduling or forecast behavior merely because the board begins exposing the saved P50 value."* The edit is **comment only** — `timelineChainStart` keeps using the projection median. The rewritten comment says the board now reads the block but this bar deliberately does not use it, which is a stronger and still-honest statement.
2. `lessons-do-kanban.md:34` (REQ-117): *"a stated reason in a comment is a factual claim and reviews must check it like any other — 'the board has no warning channel for it' shipped at 98% and was three greps from disproof."* Same shape.

`grep -rn "nested"` across the tool returns only these two claims plus test comments. No third copy.

### 5.5 The presence-vs-zero trap, already solved once next door

`generate.go:224-228` carries the exact precedent:

> *"Deliberately NOT omitempty: a genuine zero-minute span is possible … and omitempty would drop that 0 while hasImplementationSpan still shipped true — leaving the client to multiply undefined and render 'took NaNs'. The flag above is what carries presence; this field carries the value, including zero."*

The schema forbids a P50 below 5, so a genuine zero is impossible here — but the pattern (flag carries presence, field carries value) is the one to copy, and copying it means the reader is right even if the schema floor is later relaxed.

---

## 6. What the summary needs, and where each figure comes from

| Figure | Source at HEAD | Status |
|---|---|---|
| total grouped REQs | `userRequest.requestIds` — `generatedUserRequest.RequestIds` at `generate.go:269`, built Go-side from each REQ's `user_request:` pointer across queue/working/archive | **Ready.** `renderUserRequestLens` already keeps `requestIds` (full, :950) separate from `shownRequestIds` (filtered, :951). |
| active time, completed members | `hasImplementationSpan` / `implementationSpanMinutes` / `implementationSpanReason` — `generate.go:223-230` | **Ready, with two caveats.** See §6.1. |
| active time, live claimed members | `request.claimedAt` (`generate.go:197`) + `formatElapsedDuration` (`board-core.js:128`) | **Ready per-card; nothing aggregates it today.** |
| estimated remaining | saved P50 | **Missing.** §5. |
| fallback median per effort class | `boardData.timeline.projection` — `generatedTimelineProjection` at `generate.go:376-390`, assigned from `buildTimelineProjection` at `generate.go:924` | **Ready, zero Go change.** JSON keys: `confident`, `declinedReason`, `trivialMinutes`, `normalMinutes`, `trivialSamples`, `normalSamples`, `windowSamples`, `minimumSamples`, `windowSize`, `chainStart`, `queueEnd`, `rows`, `excluded`. `web/board-timeline.js:1478-1545` already reads exactly these, including the below-`minimumSamples` borrowed-median disclosure at :1527-1530. |
| successful count / pct | `request.status` | **No JS helper.** §6.3. |
| resolved count / pct | `isTerminalResolvedStatus` — `board-core.js:261-263` | **Ready.** |

### 6.1 The duration outlier rule — and a wording divergence the plan must settle

- **Constant:** `analysisOutlierCeiling = 4 * time.Hour` — `durations.go:32`.
- **The verdict:** `dayMedianExclusionReason(wallSpan)` — `durations.go:323-333`. Returns `"reversed"` (span < 0), `"paused"` (span > ceiling), `""` otherwise.
- **The measurement:** `measureImplementationSpan(ticket)` — `durations.go:222-238`, via `earliestImplementationOrigin` (`durations.go:180`) and the exclusion table `implementationSpanOriginExcludedFields` (`durations.go:162-176`).
- **The struct:** `ImplementationSpan` — `durations.go:144-149`. `ExclusionReason` is `"paused"`, `"reversed"`, or `""`.
- **Human-facing label:** `implementationSpanPausedBadgeText(ceiling)` — `durations.go:38`. **The client receives this string, never the number** — deliberately, per `generate.go:220-222`: *"The client receives only the completed paused badge text above, never a numeric ceiling it could use as a second rule."*

**Divergence D1 — origin vs claim.** REQ-486 says *"the sum of valid completed **claim-to-completion** spans accepted by the existing duration outlier rule."* The existing rule measures **earliest-origin-to-completion**, deliberately. `durations.go:206-211`:

> *"The span runs from the EARLIEST recorded lifecycle stamp to completion, not from claimed_at. A claim stamp rewritten after the work happened would otherwise erase every phase that preceded it: REQ-505 carried planning_at at 16:49 against a 23:00 claim and a 23:01 completion, and the card read 1m 23s for six hours of recorded work."*

Reusing the existing authority — which the REQ's Interfaces section **orders** — means the summary is origin-to-completion: wider where bookkeeping ran out of order, identical where it ran in order. The REQ also forbids re-deriving a competing rule in the browser. The two sentences cannot both be satisfied. Most likely the REQ's wording is loose. **The plan must state which measurement ships and why**, rather than a builder choosing silently.

**Divergence D2 — which members carry a span.** `generate.go:733-738`:

```go
// Terminal SUCCESS only: cancelled work shares the Recently-Done column
// but never took an implementation span worth stating.
implementationSpan := ImplementationSpan{}
if isCompletedStatus(ticket.Status) {
    implementationSpan = measureImplementationSpan(ticket)
}
```

A `cancelled` member counts toward **resolved** but ships `hasImplementationSpan: false`. So does a `failed` one. That is consistent with the REQ's "disclose as excluded or unavailable" — but the disclosure copy must be derived from `hasImplementationSpan` plus `implementationSpanReason` plus status, **never from a status guess about what the rule would have said**.

### 6.2 The Timeline forecast's confidence gate

- `timelineProjectionWindowSize = 60` — `timeline.go:262`. Rolling window of most recent in-rule completions.
- `timelineProjectionMinimumSamples = 5` — `timeline.go:267`. Applied twice: per bucket (`timelineBucketMedian`, `timeline.go:368-377` — a thin bucket borrows the window median rather than inventing one) and to the window as a whole.
- **Confidence gate:** `timeline.go:323-329`. Below the minimum, `Confident` stays false and `DeclinedReason` explains. Otherwise `Confident = true`.
- **Per-REQ pick:** `timelineProjectedSpan(ticket, projection)` — `timeline.go:380-390`. `effortMechanical` → `TrivialMedianMinutes`, everything else → `NormalMedianMinutes`.
- **Window split** skips `sample.ExcludedFromDayMedian()` (`durations.go:78`), so nothing here re-decides what a paused or reversed span is.

Reusable client formatter: `timelineFormatSpanMinutes(minutes)` — `web/board-timeline.js:207-218`, producing `"45 min"` / `"2h 15m"` / `"3d 04h"`. It lives at manifest position 6; declarations hoist, so a caller in `board-cards.js` or `board-detail.js` works. If the summary lands in `board-core.js` (position 1), reaching forward is legal but reads oddly — the plan should decide deliberately, not by accident.

### 6.3 The status sets, and the one predicate that is missing

**Go.** `isCompletedStatus` — `model.go:1007` (exactly `completed` or `completed-with-issues`; the exact match is deliberate so `completed-wth-issues` cannot slip through). `isCancelledStatus` — `model.go:1030`. `isTerminalResolvedStatus` — `model.go:1038`, composed as `isCompletedStatus(s) || isCancelledStatus(s)`. `isStoppedStatus` — `model.go:1057`, resolved **plus `failed`**.

**Instruction side.** `skills/do-work/actions/work-reference.md` § *Terminal-success status set* (:259), § *Dependency-source-ready status set* (:263), § *Terminal-resolved status set* (:269). The resolved section states that any reader *"cites this set by reference and must not restate or fork the set, or re-derive the rule, as a competing definition."*

**JavaScript.** Only `isTerminalResolvedStatus` exists — `web/board-core.js:261-263`, spelled as three literals, with a block comment at :256-259 saying *"Mirrors model.go's isTerminalResolvedStatus / isCompletedStatus."* Callers: `board-core.js:208`, `board-core.js:270`, `board-cards.js:205`, `board-cards.js:526`, `board-filters.js:100`.

**There is no JS `isCompletedStatus`.** `grep -rn isCompletedStatus web/*.js` returns only two comment mentions (`board-calendar.js:10`, `board-core.js:256`). So "successful = completed + completed-with-issues" has no browser-side helper.

The clean move: add `isCompletedStatus(status)` beside `isTerminalResolvedStatus` in `board-core.js` and re-express the resolved predicate as `isCompletedStatus(status) || status === "cancelled"`, mirroring the Go composition rather than adding a fourth literal list. Two lessons converge on this:

- `lessons-do-kanban.md:20` (0.107.0): *"`cancelled` … is terminally resolved but not terminally successful … Keep the two classifiers distinct; success-readers must never widen to the resolved set."*
- `lessons-do-kanban.md:45` (0.294.2, family `paired-predicate-drift`): *"two readers of one contract drift silently, because each one passes its own tests … the fix was a predicate, not a parser: `isDependencySourceReadyStatus` beside `isCompletedStatus`, named after the contract section it implements … a reader that names one arm of a two-arm set is drift even when it compiles and passes."*

Existing probes slice the shipped `isTerminalResolvedStatus` (`javascript_behavior_a_test.go:1457`, `_b:87`, `_c:483`). Adding a sibling keeps those slices working; a new probe should **slice** the new function, not stub it.

### 6.4 The board clock

- **One ticker.** `web/board.js:68` — `setInterval(refreshRelativeTimeNodes, 1000);`, the last line of the boot block.
- **What it does.** `refreshRelativeTimeNodes()` — `web/board-core.js:236-253`. Captures **one** `nowMs` at :237, selects `document.querySelectorAll("[data-instant-ms]")`, and per node recomputes the label — `formatElapsedDuration` when `dataset.tickFormat === "duration"`, else `formatRelativeTime`. Writes only on change; calls `syncClockSkewTitle` for duration nodes.
- **How a node opts in:** by carrying `data-instant-ms`. Builders: `makeRelativeTimeNode` (:68), `makeElapsedDurationNode` (:150-160), composites at :81 and :165.
- **How a node opts out deliberately:** `makeImplementationSpanNode` (`board-cards.js:66`) builds a plain span with no `data-instant-ms` — a finished span must not be rewritten every second. Pinned at `javascript_behavior_b_test.go:1681` and :1749.
- **There is no subscription API.** No callback registry exists. The only extension point is the attribute selector.

For a summary whose live claimed contribution must tick, the attribute route does not work: it ticks a **single** claim's elapsed time, not a sum of a completed base plus N live claims. The requirement — *"Refresh live claimed contributions through the board's existing clock so the header and drawer cannot drift from the claimed card stopwatch"* — is literally a statement about **two `Date.now()` calls in different frames**. The honest fix is a fan-out in `board-core.js`: `registerBoardTickSubscriber(callbackFunction)` plus `refreshTickingSurfaces()` that captures one `nowMs` and passes it to both `refreshRelativeTimeNodes` and every subscriber, with `web/board.js:68` pointed at the fan-out.

**Skew, reused not re-derived.** `futureInstantSkewAllowanceMs = 2 * 60 * 1000` (`board-core.js:105`, mirroring Go's `futureTimestampSkewAllowance` — the comment at :103-104 says keep in lock-step), `clockSkewMarkerText = "⚠ clock skew"` (:106). `formatElapsedDuration` returns the marker for a future instant rather than clamping to `0s` (:129-131). **A claimed member with a future `claimed_at` has no usable live elapsed time and must be disclosed, not counted as zero.** `lessons-do-kanban.md:23` (0.133.0) records the whole rule.

Unrelated second timer, do not confuse: `web/board-timeline.js` runs its own 50ms interval for plot-width detection.

---

## 7. The two surfaces, and the CSS problem

**By UR header** — inside `renderUserRequestLens`, `web/board-cards.js:925`. Today's only summary is one `span.ur-count` at :1006-1014, which says `"12 / 43 REQ"` when filtered and `"43 REQ"` when not.

**UR drawer** — `openUserRequestDetail`, `web/board-detail.js:614-636`. That is the `GROUPED REQS 43` plus long REQ-id list from the screenshot. `appendMetaRow` is `board-detail.js:321-333`; `makeTicketLinkList` is `board-detail.js:409-422`, one `span.detail-dep` per id reusing `.detail-dep-list` (`board.css:2065`).

**Decision D-A the plan must make.** `ur-count` states a **filter-dependent** shown/total; the new summary states a **filter-independent** total. Side by side under an active filter they will look like a contradiction. `lessons-do-kanban.md:46` (REQ-588, family `subject-not-restated-in-detail`): *"when a payload gains a structured field, the prose that carried the same fact must stop carrying it, and every test that read the fact out of the prose must move to the field."* Good news for whoever changes it: **no test reads `ur-count`** — `grep -rn "ur-count" *_test.go` returns nothing, and the REQ-236 probe decodes `ur-id`, `aria-expanded`, card ids and drawer triggers only (`generate_test.go:2791-2796`, `javascript_behavior_a_test.go:1552-1557`). The wording is free; the decision still has to be made.

### 7.1 The layout is a real change, not just markup

`web/board.css:1482-1496`:

```css
.ur-group-head {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  ...
  padding: 14px 18px;
}
```

**No `flex-wrap`.** And `web/board.css:1515-1524`:

```css
.ur-title { color: var(--ink-strong); font-weight: 520; font-size: 14px;
  flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
```

Adding metric spans without `flex-wrap: wrap` (or a dedicated second row) will **squeeze the ellipsised title**, not wrap — exactly what the REQ forbids at narrow widths. Current UR block map at HEAD: `.user-request-lens` :1469, `.ur-group` :1475, `.ur-group-head` :1482, `.ur-id` :1505, `.ur-title` :1515, `.ur-count` :1525, `.ur-synthetic` :1531, `.ur-group-cards` :1541, REQ-236's folded block :1552-1601, `.ur-lens-empty` :1602, `.ur-lens-hidden-note` :1608.

**Colour warning.** `.ur-count` uses `--ink-faint`. `prime-kanban-board.md:29` records: *"`--ink-faint` against `--ink-soft` measures 1.29:1 light and 1.82:1 dark, which is not a channel and cannot carry a row distinction."* Approximate / unavailable / partial markers need a real channel — a word, a symbol, a border — not a second ink tone.

---

## 8. Test surface

### 8.1 Node behavior lane (harness in `generate_test.go`)

- `generateLiveSiteInDir(t)` :335 / `generateLiveSite(t)` :361 — build the board against the real `do-work/` tree, cached `index.html` string.
- `sliceBalancedBlockAfter(t, sourceText, anchorToken)` :1998 — brace-matches one function out of the assembled page, so the probe drives the **shipped** source, not a copy.
- `runJavaScriptBehaviorProbe(t, probeName, probe)` :276-293 — pipes the probe to `node -` on **stdin** (an `-e` argument exceeds Linux's 128 KiB per-arg limit; macOS would pass and CI would fail). Increments `javaScriptBehaviorProbeCount` (:29) for the strict lane's zero-probe guard (:203).

**There is no synthetic-queue fixture builder for the JS probes.** Each probe hand-writes a `var boardData = {...}` literal plus a DOM stub. The UR probe's stub is `makeNode()` at `javascript_behavior_a_test.go:1495-1520` — `childNodes`, `dataset`, `attributes`, `listeners`, `appendChild`, `removeChild`, `setAttribute`, `getAttribute`, `addEventListener`, and a synchronous `dispatch(eventName)`. `document` is a two-method stub at :1522-1525; `makeRequestCard` is stubbed to `{ className: "req-card", requestId: requestId }` at :1526.

**Probes covering `renderUserRequestLens`:**

| Test | File:line | What it pins |
|---|---|---|
| `TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened` | `javascript_behavior_a_test.go:1453` | The full REQ-236 contract. Slices twelve functions (:1456-1467). Fixture: 3 URs, 4 REQs, `Date.now` stubbed to `2026-08-15T12:00:00Z` (:1470). **Its last assertion, :1652, currently demands the By UR head carry no `aria-expanded` — this REQ inverts it.** |
| `TestJavaScriptBehaviorByUserRequestLensCountsRecentlyDoneAsActive` | `javascript_behavior_b_test.go:75` | Active scope + recently-done window. |
| `TestJavaScriptBehaviorByUserRequestLensEmptyState` | `javascript_behavior_b_test.go:139` | The three empty-state branches. |
| `TestJavaScriptBehaviorTestingStatusUpdateInvalidatesUserRequestLens` | `javascript_behavior_b_test.go:180` | `renderedOnce.userRequestLens` invalidation. |
| `TestJavaScriptBehaviorByUserRequestLensUsesRecentWindowAtCaller` | `javascript_behavior_c_test.go:479` | Window read at the caller. |
| `TestGenerateOffersThreeLensButtons` | `generate_test.go:2773` | Source tokens only: three `data-lens-target` values, `data-ur-cards="folded"`, `URs&nbsp;only`, the `applyLensSelection(...)` call site verbatim. |

Shared decode type `renderedUserRequestRow` — `generate_test.go:2791-2796`; `userRequestIdsOf` at :2803.

**The drawer has no probe at all.** `grep -rn openUserRequestDetail *_test.go` → zero hits. A new probe needs a DOM stub covering the ten ids read at `web/board-detail.js:291-300`: `detail-resizer`, `detail-drawer`, `detail-kind`, `detail-id`, `detail-drawer-title`, `detail-meta`, `detail-body`, `detail-glossary`, `detail-copy`, `detail-copy-all`. It must drive the shipped `setDetailTarget` writer (`board-detail.js:315-319`), never assign `currentDetailKind`/`currentDetailId` by hand — `lessons-do-kanban.md:15` (0.295.1): *"A probe that assigns those variables by hand in the order the listener assumed cannot see this class of bug at all — drive the shipped writer."*

### 8.2 Browser lane

Harness: `browser_probe_test.go` — DevTools pipe straight to the browser binary, result written into `#queue-kanban-probe-result` (:49). `QUEUE_KANBAN_BROWSER` overrides the binary (:44); `QUEUE_KANBAN_BROWSER_PROBES` (:39) and `QUEUE_KANBAN_STRICT_BROWSER_BEHAVIOR` (:38) gate the lane. Entry points at :98, :111, :124.

`TestBrowserBehaviorUserRequestCopyAllIncludesGroupedRequests` (`user_request_clipboard_browser_probe_test.go:44`) is **the one browser probe that drives the real By UR UI and opens a UR drawer.** Its fixture builder `userRequestCopyFixture` (:28) composes `*RequestTicket` structs across `queue`, `working` and `archive` sections with `pending`/`claimed`/`completed` statuses (:52-61). Closest existing thing to a synthetic UR fixture at the Go level.

Go-side synthetic tree: `createSyntheticDoWorkTree(t)` / `syntheticBoard(t)` — `board_synthetic_test.go:25` / :68. Right base for Go-level fixtures carrying `estimate:` blocks; it does **not** feed the Node probes.

### 8.3 Gap analysis against the REQ's RED case

The REQ asks for members covering: completed, completed-with-issues, cancelled, pending, claimed, blocked, failed, missing timestamps, an outlier span, a saved P50, a missing P50 with confident history, and insufficient-history fallback.

- **Statuses** — trivially more object literal in the inline `boardData.requests`. No new helper.
- **Missing timestamps / outlier span** — the client never re-derives the outlier verdict; it reads `hasImplementationSpan` / `implementationSpanMinutes` / `implementationSpanReason`. So the probe expresses these as **payload shapes**, and the outlier rule itself is exercised on the Go side (`durations_test.go`). **Two lanes, two questions** — say so in the plan.
- **Saved P50** — needs the new payload field first, plus the strict-vs-salvage parse pair.
- **Missing-P50-with-confident-history / insufficient-history** — stub `boardData.timeline.projection` with `confident: true` plus medians and sample counts, then a second run with `confident: false` and a `declinedReason`. No new machinery.
- **Stubbed clock** — the idiom is at `javascript_behavior_a_test.go:1470`. Advancing it is a reassignment between calls, but the probe **must call the tick entry point itself** — `setInterval` never runs there.
- **Both surfaces report identical values** — the cleanest RED slices both shipped renderers and replaces the summary with a spy. `lessons-kanban-board.md:56` (REQ-305): *"a probe that calls the function under test directly cannot hold its call site — five mutations of the copy were caught and the one that reverted the real defect passed clean."*
- **Two-theme / narrow-width render** — browser lane. `prime-kanban-board.md:24`: *"A chart's correctness is partly a claim about pixels — generate a board and look at it … REQ-226, REQ-231, REQ-237 and REQ-240 each shipped a defect that every assertion passed over and a render made obvious."* `prime-kanban-board.md:25`: return `location.href` alongside every measurement. The lane **skips silently** without `QUEUE_KANBAN_BROWSER`; a skip is not a pass.
- **Keyboard operability** — `lessons-kanban-board.md:9` (REQ-338): *"Tab cannot be tested with synthetic events at all, since its focus movement is a trusted-input default action."* A synthetic-event probe proves `aria-expanded` flips; **both fold controls being real `<button>`s is what actually delivers Tab order.**
- **Disclosure assertions** — `lessons-kanban-board.md:26` (REQ-245): *"asserting a phrase is absent is not a guard — it passes when the whole string is replaced."* Assert the **presence** of the qualifier, never the absence of a number.
- **Threshold fixtures** — `lessons-kanban-board.md:8` (REQ-374): *"a fixture that spans a threshold widely does not test it … only a pair straddling the real boundary, derived from the constant, catches a second definition."* Directly on the four-hour outlier ceiling and the five-sample confidence floor.
- **Missing-branch fixtures** — `lessons-kanban-board.md:55` (REQ-304): *"a missing-branch fix needs a fixture that can fail in both directions."* Directly on the three-arm remaining-time forecast.

---

## 9. Documentation the change owes

1. **`skills/do-work-board/docs/board-guide.md:25`** — *"a **Lens** toggle (flat Columns vs. grouped **By UR**)"*. Names two lenses; three ship. `grep -rn "URs only"` across `docs/`, `SKILL.md`, `actions/` and `_dev/primes/prime-kanban-board.md` returns **nothing** — the guide has been stale since REQ-236, before this REQ adds anything. Fix it here, or say why not.
2. **`board-guide.md:52-54`, "## Card drawer"** — describes the REQ drawer only. The UR drawer's rows are undocumented; this REQ adds a summary and a fold to them.
3. **`skills/do-work/actions/work-reference.md`, the `estimate:` block at :125-139** — earn the same "parsed by `model.go` … keep that parser in lock-step, both changing in the same commit" clause its display-only siblings carry at :119, :121, :123.
4. **`skills/do-work-board/actions/board.md:120`** — the field-by-field board-role list. Add the P50 read there (`lessons-do-kanban.md:35`, REQ-116).
5. **`timeline.go:394`** — the false claim. §5.4.

**What this REQ does NOT touch:** `_dev/primes/prime-kanban-board.md:14` pins the tool at exactly **three write surfaces** and `_dev/tests/contract-regressions.sh` enforces the count. REQ-486 adds none; the sentence stays as written.

---

## 10. Assumptions and decisions the plan must settle

- **D-A.** `ur-count` vs the summary's total — one of the two stops carrying it. §7.
- **D-B.** Origin-to-completion vs claim-to-completion for "active time spent". §6.1. Recommend origin-to-completion, stated openly, because the Interfaces section forbids a second measurement.
- **D-C.** Where `summarizeUserRequestProgress` lives (`board-core.js` recommended, §4) and whether `timelineFormatSpanMinutes` is reached forward from position 6 or moved.
- **D-D.** Tick fan-out shape. Recommend `registerBoardTickSubscriber` + `refreshTickingSurfaces` in `board-core.js`, one captured `nowMs`, `board.js:68` repointed. §6.4.
- **D-E.** The rewrite of `javascript_behavior_a_test.go:1652`. That assertion is green today and this REQ makes it wrong. A builder who "fixes" the resulting failure by not setting `aria-expanded` has silently dropped a requirement.
- **D-F.** How the summary renders at narrow widths — `flex-wrap: wrap` on `.ur-group-head`, or a dedicated second row. §7.1.
- **A1.** `omitempty` choice for the new payload fields — flag carries presence, value field carries the number, per `generate.go:224-228`.
- **A2.** The drawer summary must deregister its tick subscriber (or no-op) on close, or a closed drawer keeps recomputing.

---

## 11. Route: C

Not Route B. `skills/do-work/actions/work.md:61-69` puts Route C at *"New feature requiring multiple components, architectural changes … touches multiple systems, long request (100+ words) with many requirements."* REQ-486 is all four:

- **Two languages, two packages.** A Go frontmatter reader for a nested block plus a payload field plus a new numeric coercer; a JavaScript rollup, a new status predicate, and a clock fan-out.
- **Two independent fold surfaces in two different fragments**, one of which (`board-detail.js`) has no behavior probe at all.
- **A CSS layout change** to a non-wrapping flex row whose correctness is a claim about pixels in two themes at two widths.
- **Roughly 14 files**, three test lanes, five documentation edits.
- **Two unresolved contract questions** (D-A, D-B) that a builder must not settle silently — and one green assertion this REQ inverts (D-E).

Route B would hand a builder an exploration with two open semantic questions, no agreed edit order, and a shipped test that fails for the right reason but looks like a regression.

---

*Generated by Explore agent*