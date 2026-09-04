# REQ-486 read-only plan and exploration

Extend the existing UR renderer and drawer, using one shared browser summary over the complete `userRequest.requestIds` membership. The current generated payload already carries Go's duration verdict and Timeline's confident medians. Only the saved P50 reader is missing; no second duration or forecast engine is needed.

Prepared against inspected source ending at `f0aeb4a3093c9dd04845aa955316c7e7f3493083`. This is later-step preparation, not a claim, approved Scope, or test result. No tests ran. Re-resolve the request, dependency REQ-510, current source, and status vocabulary at claim time: concurrent REQ-570/571 work can change status-related owners and tests. Retain REQ-486's existing deferral until the canonical dependency gate permits selection.

Read CLAUDE, core skill/work planning rule, REQ-486, archived `do-work/archive/UR-054/REQ-236-add-a-urs-only-lens-to-the-board.md`, both board primes, both whole board lesson satellites, board guide, and general/coding/frontend/testing/shared-principles/security/anti-slop crews. The two partially slugged satellites remain whole-file touch-conditional reads despite their capture-budget drops. No current frontend/testing-specific prime was found in the repository inventory.

## Existing owners and boundaries

- `model.go`: `parseRequestTicket` extracts model fields; `linkRequestsToUserRequests` groups `board.AllRequests`, so generated `requestIds` already includes queue, working, and archive. Use this membership, never the capture-time UR `requests:` array or filtered cards.
- `frontmatter.go`: normal YAML parsing already returns nested maps. Its final lenient salvage path is intentionally flat. Read `estimate.p50_active_minutes` from the normal parsed map; do not add a bespoke YAML scanner or broaden salvage.
- `durations.go`: `measureImplementationSpan` owns parsed stamps, signed minutes, and `dayMedianExclusionReason`. `generate.go` already exports `hasImplementationSpan`, `implementationSpanMinutes`, and `implementationSpanReason` for successful statuses only. Missing spans are distinguishable from valid zero-length spans. Do not consult `CompletionTime`, which can fall back to a git timestamp.
- `timeline.go`: `buildTimelineProjection` owns history confidence and both effort medians; `timelineBucketMedian` already substitutes the overall window median when a bucket is thin. `generate.go` exports `timeline.projection.confident`, `trivialMinutes`, and `normalMinutes`. Browser code consumes these results; it never copies the history minimum, window size, outlier ceiling, or median algorithm. `timelineChainStart` contains the stale no-nested-reader comment, but its scheduling calculation must remain unchanged.
- `web/board-cards.js`: `renderUserRequestLens` derives `shownRequestIds` solely for cards/visibility. Currently only URs only has a fold header and separate Details button. Use that shared path for both readings, changing only the initial open state.
- `web/board-detail.js`: `openUserRequestDetail` appends the grouped count and raw REQ-link list. Add metrics and a separate list disclosure here without hiding `input.md` or changing drawer navigation/copy behavior.
- `web/board-core.js`: `refreshRelativeTimeNodes` captures one `Date.now()` per tick. `web/board.js` already schedules it every second. Add summary-node refresh to this same call and pass its captured instant; do not introduce another interval or re-render whole groups/drawers every second. Reuse elapsed formatting and the existing clock-skew allowance/marker semantics.

## Plan — three tasks

1. **Prove and expose the existing saved estimate.** Add RED parse-to-generated-payload cases from real Markdown, including a nested numeric P50, absence, null/wrong shape, nonfinite/negative values, and lenient-salvage absence. Add an optional numeric model/payload field that preserves unavailable evidence rather than defaulting to zero. Keep the reader tolerant of finite positive numeric values without inventing a second estimator or mutating schema data. Verify that adding saved P50 values leaves the Timeline projection's schedule, medians, and confidence unchanged. Correct only the now-stale reader comments.
2. **Implement the single rollup and both fold surfaces.** Add RED behavior cases over two URs plus empty/synthetic groups. Implement a pure shared summary taking the complete UR and a captured current instant, and one shared metric-node renderer/updater used by header and drawer. Extend the existing URs-only fold path to By UR with `aria-expanded="true"` initially, retain URs-only's false default, and preserve separate Details controls. Add an independent initially-open drawer list button; only its linked IDs hide. Keep open/closed state in the current DOM and reset it on a normal re-render. Update existing tests whose deliberate old contract says By UR must lack `aria-expanded`, and tests whose selectors mistake the head for a drawer trigger.
3. **Verify the actual generated UI and document its limits.** Add a browser fixture that drives the real generated page, normal controls and shared clock, with trusted keyboard input and both themes at 320/768/1280 pixels. Verify folding independence, explicit whole-UR metrics and denominators, filter invariance, clock synchronization, wrapping, and no console errors. Update the board guide's UR reading and active-time/forecast explanation. Run focused RED/GREEN first, then normal board checks, the unpiped maintainer gate, and the canonical selected strict JavaScript/browser/heavy lanes. Record actual lane exits/skips and rendered evidence; this preparation makes no pass claim.

## Proposed exact implementation Scope

All paths below are repository-relative. Declare them before building; any newly required path needs an explicit Scope amendment before editing.

| Path | Change |
|---|---|
| `skills/do-work-board/tools/queue-kanban/model.go` | Optional saved-P50 model field and nested reader |
| `skills/do-work-board/tools/queue-kanban/generate.go` | Optional numeric payload field and projection wiring |
| `skills/do-work-board/tools/queue-kanban/frontmatter.go` | Comment only: flat salvage is a parser limitation, not the model's complete read surface |
| `skills/do-work-board/tools/queue-kanban/timeline.go` | Comment only: retain current Timeline assumption despite new reader |
| `skills/do-work-board/tools/queue-kanban/web/board-core.js` | Shared summary/metric rendering and existing clock hook |
| `skills/do-work-board/tools/queue-kanban/web/board-cards.js` | Both lens defaults, separate Details, whole-UR metric consumer |
| `skills/do-work-board/tools/queue-kanban/web/board-detail.js` | Same metric consumer and independent REQ-list disclosure |
| `skills/do-work-board/tools/queue-kanban/web/board.css` | Wrapping summary, fold controls and narrow-width containment |
| `skills/do-work-board/tools/queue-kanban/user_request_summary_test.go` | New parse/projection/summary behavior fixture tests using existing Go/Node harness |
| `skills/do-work-board/tools/queue-kanban/user_request_summary_browser_test.go` | New real generated-page interaction/render test using existing Chromium transport |
| `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` | REQ-236 fold-contract update and real helper dependencies |
| `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` | Recent-window caller fixture supports fold markup; preserves its visibility assertions |
| `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go` | Open drawer through new dedicated Details selector; retain complete membership/copy assertions |
| `skills/do-work-board/tools/queue-kanban/generate_test.go` | Citation-marker browser probe selects the fold head by UR identity rather than removed drawer-trigger attributes; retain marker/accessibility assertions |
| `skills/do-work-board/docs/board-guide.md` | Both fold defaults, metrics, sources and partial/unavailable meanings |

No change proposed to `durations.go`, forecast algorithms, status classifiers, frontmatter writers, filters, lens state, `web/board.js`, template/fragment manifest, browser harness, dependencies, or generated artifacts. Reuse current embedded fragments; a new production fragment would add assembly work without removing an existing responsibility.

## Formula and test boundaries

Let N be complete grouped membership. Successful is completed plus completed-with-issues; resolved adds cancelled. Failed counts in N and neither numerator. Each label contains count/N and percentage; N=0 prints unavailable. Filters can change cards or group visibility, never summary inputs. If a payload membership ID cannot resolve, retain its place in N and disclose unavailable evidence rather than silently shrinking the group.

Spent time sums only successful members with `hasImplementationSpan` and no Go exclusion reason, plus usable live claimed elapsed time. Track paused, reversed, missing-completed-span and unusable-claim counts separately enough to explain a partial value. Preserve true zero-length accepted spans. Pending/blocked waiting and failed/cancelled lifetime do not manufacture active elapsed spans. A clock-skewed live claim must agree with the existing stopwatch's unusable/skew state, not contribute a misleading clamped zero.

Remaining time ignores resolved members; failed is explicitly unknown even if a saved estimate exists. Other unfinished members use a usable saved P50 first, otherwise the supplied effort median only when `projection.confident`. Blocked and pending-answer members retain estimated active effort irrespective of external wait. For a usable claimed member subtract live elapsed from its selected estimate and floor at zero. A claimed member with no usable claim instant has unknown remaining time because its elapsed subtraction cannot be measured. Present known contributions as approximate and partial whenever any unfinished member is unknown; all-unknown is unavailable, not zero. An empty or fully resolved group can truthfully have zero remaining effort while its zero-denominator progress or historical elapsed remains unavailable.

Tests must exercise these through real parsed records and the generated payload, not a standalone alternative rollup fixture. For outlier agreement, derive cases immediately at and beyond the Go ceiling; for confidence use the actual minimum boundary and verify bucket-to-overall fallback. Include malformed/missing stamps, reversed spans, valid zero spans, saved estimate precedence under thin history, both effort buckets, blocked/pending-answer, failed, and elapsed beyond estimate. Advance a stubbed current instant across formatter boundaries and compare header, drawer and claimed card at the same clock callback. Folding must leave the current instant and metrics untouched.

Keep the existing cross-REQ properties: REQ-236's multiple independent groups and ephemeral folds; REQ-122's Active/All plus recent-window visibility; REQ-368's full membership copy despite a stale capture array; citation marker containment and keyboard accessibility. Update selector or harness assumptions only where the requested behavior supersedes them, naming that in Testing.

## Visual acceptance

Generate from a fixture whose two URs include long titles, long metric labels, sufficient completed history, excluded/unavailable contributors, and an open drawer. Inspect normal and narrow layouts in light/dark; assert positive element bounds and containment/no intersections between title, metrics and Details. Scope card counts to `#user-request-lens`, because hidden Columns retains its cards. Use trusted Tab/Enter/Space for keyboard behavior, not synthetic events or `.focus()` alone. Verify header fold does not open the drawer, Details does not change folds, drawer-list fold leaves metrics and body visible, and multiple groups can remain independently open.

Return `location.href` and browser build in the same call as every measurement; use the existing isolated browser transport and record the exact configured Chromium. Compare text as well as geometry and inspect screenshots. No test framework or persistent fold store is added. Respect the fast-test per-file duration budget; strict browser work remains in its existing heavy lane.
