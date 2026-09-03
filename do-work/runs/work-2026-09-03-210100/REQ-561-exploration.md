# REQ-561 exploration

Explored against main-tree commit `056d54d68fa7ead5e370cb184c1c5505038b4b5a`. The implementation plan is sound at the contract level. Three scope corrections are required before build:

1. Add `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go` to the expected manifest. It is the existing semantic home for `resolveTargets`, mixed explicit/UR token order, and UR-member ordering. `next_selection_test.go` should still cover the final selected/excluded results, but it should not be the only lock on the UR comparator.
2. Use a dedicated board browser test such as `skills/do-work-board/tools/queue-kanban/priority_browser_probe_test.go`. Do not add priority behavior to `clipboard_browser_probe_test.go` unless the clipboard fixture itself must change. That test is a focused REQ-367 copy-all contract with hand-authored column order and exact clipboard payloads; coupling the priority UI to it would enlarge an unrelated contract. The shared browser harness stays in `browser_probe_test.go` and needs no production edit.
3. `capture-reference.md` needs the optional priority line in the Addendum REQ Template as well as the Simple/Complex REQ template, or the prose beside that template must explicitly say that Step 1's priority assessment may add it. The request explicitly says addenda may change priority. For a queued addendum, `capture.md` must say set/change/remove `priority` in the existing pending REQ frontmatter from the user's new ranking words. For a new addendum REQ created for working/archived work, the ordinary Step 1 assessment applies to the new REQ.

## Confirmed implementation seams

### Core schema and typed request model

- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go`
  - `fieldContracts` is the authoritative executable table. Add `priority` with canonical values `now`, `next`, `later`, no aliases, and default `next`.
  - `NormalizeField` already trims, lowercases, recognizes canonical values case-insensitively, treats absence as recognized/defaulted, and emits the exact fallback warning for invalid present values. No new normalization function is needed.
  - `SchemaFieldNames()` automatically includes the new row and will force typed projection through the existing inventory test.
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go`
  - Extend `TestNormalizeFieldAppliesAliasesDefaultsAndExactWarnings` or add a compact priority table covering absent, all three canonical values, case/whitespace normalization, and one invalid token. Assert `ResolvedValue`, `IsRecognized`, `IsDefaulted`, and exact warning behavior.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go`
  - `RequestRecord` carries one explicit value/evidence pair for every contracted field. Add `RequestPriorityValue string` and `RequestPriorityEvidence schemanormalization.FieldResult`.
  - `RequestDocument.TypedRecord()` must call `NormalizeField("priority", document.scalarValue("priority"))` and populate both fields. `FieldEvidenceByName` already retains the generic raw field when present; absence correctly means no generic map entry while the typed value still resolves to `next`.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go`
  - Add a `priority:` line and a projection row to `TestTypedRecordCarriesEveryNormalizedSchemaFieldAndGenericEvidence`. That test compares its projection names to `SchemaFieldNames()`, so adding the schema row without both typed fields will fail for the intended reason.

### Result projection and selector

- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
  - `SelectionRecord` and `SelectionExclusion` currently expose `SelectionPriority string \`json:"selection_priority"\`` for the scheduling class. Keep it unchanged.
  - Add `RequestPriority string \`json:"priority"\`` to both structs. In `NormalizeResult`, default an empty value to `next` for both selected and excluded records. This gives target-not-found/ambiguous exclusions, which have no source request record, a deterministic wire value.
  - The current text renderer prints only `selection_priority`. The captured acceptance specifically requires typed/JSON projection and says callers must not infer the authored value from display text, so changing text output is not required. If the builder elects to show both values in text, update `TestSelectionTextAndJSONCarryTheSameTypedCommands` deliberately rather than silently changing the established bracket format.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
  - Extend `TestSelectionTextAndJSONCarryTheSameTypedCommands` for JSON projection and add/extend normalization coverage so empty request priority becomes `next` without altering `selection_priority`.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go`
  - Keep the existing class constants `PriorityRepositoryGateRepair`, `PriorityDeferredParent`, and `PriorityOrdinary` separate from authored priority. New names should say `RequestPriorityNow`, `RequestPriorityNext`, and `RequestPriorityLater`; a separate ranking helper avoids conflating the two axes.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
  - `Select` evaluates dependency/status/assignment/impact eligibility first, then stable-sorts eligible default-scan records at lines 42-46, then applies the fan-out bound at lines 48-59. Extend that comparator to class rank first and authored priority only when both records are ordinary. This preserves gate-repair and deferred-parent precedence and places priority before fan-out bounding.
  - `evaluateCandidate` has one `newExclusion` closure used by every real-request exclusion. Set `RequestPriority` there from `record.RequestPriorityValue`, and set it on the selected record at its construction site. This covers invalid requests with a parsed record, already-claimed, blocked, status, assigned, impact, simple-only, wave, and dependency exclusions.
  - `copySelectionEvidenceToExclusion` is the fan-out conversion seam. Copy `RequestPriority` there with `SelectionPriority`.
  - `appendSchemaWarnings` is the selector warning inventory. Add priority evidence once. `Select` invokes it once per candidate before accepting or excluding the candidate, so an invalid present priority yields one `SCHEMA-FALLBACK`; absence remains recognized and quiet.
  - `exclusionFor` should initialize authored priority to `next`. `targetNotFound` and `targetAmbiguous` in `next_targets.go` bypass `exclusionFor`, but `NormalizeResult` supplies their final wire default. Direct `Select` tests that inspect those synthetic exclusions should not pretend an authored value existed.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go`
  - Default `queueCandidates` already establishes numeric queue order. The later stable selector sort therefore preserves numeric order for equal class/priority keys.
  - Explicit REQ tokens append in caller order and `Select` deliberately does not globally sort when any target token exists. Leave that path unchanged.
  - UR members are currently stable-sorted by class rank, then dependency depth, then request id inside `resolveTargets`. Insert authored priority after class rank and only when both members are ordinary, before depth/id. Mixed explicit and UR tokens must keep their token anchors; only each UR expansion is internally ordered.
- Tests:
  - `next_selection_test.go`: default ordering with conflicting ids, class precedence, equal-priority stability, dependency gating (`later` prerequisite selected while `now` dependent is excluded), fan-out exclusions retaining authored priority, invalid warning/default, and both priority fields remaining distinct.
  - `next_targets_test.go`: UR-expanded comparator and mixed explicit/UR anchors; explicit REQ order must win even when authored priorities conflict.
  - `next_commands_test.go`: extend `commandSelectionRecord` with `RequestPriority string \`json:"priority"\`` and assert selected plus fan-out-excluded JSON through the real `go run` seam under `DO_WORK_HEAVY_TESTS=1`.

### Board parser, ordering, payload, and cards

- `skills/do-work-board/tools/queue-kanban/model.go`
  - `schemaReadContractFields` is an intentionally independent board parser table. Add the same `priority` row and constants.
  - Add `Priority`, `OriginalPriority`, and `PriorityUnrecognized` to `RequestTicket` near impact/effort. Unlike those display-only fields, priority must always call `resolveSchemaField("priority", raw)` even when absent, because the effective default `next` participates in sorting. Preserve the raw scalar and recognition flag.
  - Add priority to `collectSchemaFieldWarnings`. The generic `schemaFieldWarningText` already supplies exact fallback wording; absence resolves recognized and produces no warning.
  - `buildBoard` sorts `AllRequests` numerically before dependency annotation and `bucketColumns`, so the current inputs to each pending slice are already in the desired tie order.
  - `bucketColumns` is the narrow ordering seam. After bucketing, stable-sort `PendingReady` and `PendingWaiting` independently by authored priority rank, then rebuild `Pending` as Ready followed by Waiting. Do not sort client-side. This also covers a bare `status: blocked` ticket with unmet dependencies because it is intentionally in Pending Waiting.
  - Leave Claimed, Needs input / Blocked, Recently Done, completion anomalies, calendar, and `AllRequests` ordering untouched.
- `skills/do-work-board/tools/queue-kanban/model_test.go`
  - Existing schema tests around `TestNormalizeSchemaFieldCoversContractAliases`, `TestResolveSchemaFieldFallsBackWithoutSilentRemap`, and the parse-level unrecognized-field tests provide the pattern. Add parse-level cases for absent, all canonical values, case/whitespace, and invalid, with effective/raw/unrecognized/warning assertions.
  - Add an ordering fixture with ids deliberately opposed to `now`/`next`/`later` in both Ready and Waiting and an equal-priority pair. Assert both group slices and the rebuilt `Pending == Ready + Waiting` union.
- `skills/do-work-board/tools/queue-kanban/generate.go`
  - Add non-omitempty `Priority string \`json:"priority"\`` plus raw/unrecognized fields to `generatedRequest`, and project them in `buildGeneratedBoardDataWithMentions`.
  - Generated columns already copy `board.Columns.Pending`, `PendingReady`, and `PendingWaiting` without sorting. Keep that shared path; `serve.go` calls the same `buildBoard + buildGeneratedBoardDataWithMentions` projection.
- `generate_test.go` and `serve_test.go`
  - Static generation must start from real fixture files and decode `board-data.js`, not a hand-built `Board`, to cover parse, warning, sort, and projection together.
  - `serve_test.go` already has `fetchServedBoardData` against `newLiveBoardServer`. Use the same priority fixture and compare live priority payload plus all three pending-order slices against static output. This is the real static/live seam.
- `web/board-cards.js`
  - `buildRequestCard` constructs `req-card-badges` through `makeBadge`. Add one badge only for effective `request.priority === "now"` or `"later"`; no badge for `next`, absent legacy payloads, or an invalid raw value resolved to `next`. The warning banner remains the invalid-value footprint.
  - Give the badge a stable class/data hook so the browser probe can select it without relying on all badge text.
- `web/board.css`
  - Use the existing `.badge` geometry and theme variables. Add a priority base/value modifier only as needed; the palette variables already switch under `prefers-color-scheme`, so duplicate dark/light declarations are unnecessary. Verify contrast against the actual card background, wrapping, positive rects, and horizontal containment.
- Browser harness
  - `browser_probe_test.go` owns `startTrustedInputBrowserSession`, `runBrowserBehaviorProbeInDirectory`, Chromium discovery through `QUEUE_KANBAN_BROWSER`, strict counting, DevTools media emulation, and same-evaluation `location.href`. Reuse it without editing the harness.
  - A dedicated `TestBrowserBehaviorPriorityOrderAndBadges` should generate the real static fixture, open that directory, emulate light and dark with `Emulation.setEmulatedMedia`, and return in one evaluation: `location.href`, `navigator.userAgent`, Ready ids, Waiting ids, badge text per card, computed colors/backgrounds, and badge/card rectangles. Assert no `next` badge and no overflow/overlap.
  - The current harness opens `file://.../probe.html`. For live HTTP evidence, start `newLiveBoardServer` on loopback and use the same DevTools session to navigate to its exact URL, or perform the live exact-URL evidence manually. Do not call a `file://` result “live.”

### Action contracts

- `skills/do-work/actions/work-reference.md`: add the optional frontmatter field, Schema Read Contract row, and selection-order contract. Name both output axes: `selection_priority` is scheduling class; `priority` is authored rank. State that dependency readiness filters first, ordinary priority sorts next, fan-out bounds last, and explicit REQ caller order is unchanged.
- `skills/do-work/actions/capture.md`: add a Step 1 priority assessment keyed only to explicit user ranking/timing language. Do not infer from impact, effort, dependency depth, age, numeric id, or source list order. In queued-addendum handling, explicitly allow set/change/remove from new user words. Add final-reread traceability for each emitted priority.
- `skills/do-work/actions/capture-reference.md`: add the commented optional line to the base template and cover the Addendum REQ Template as noted above. Add `priority` to the normalize-and-warn enum inventory.
- `skills/do-work-board/actions/board.md`: document effective default, invalid warning, Ready/Waiting-only ordering, and `now`/`later` badge behavior. State that priority does not change status or dependency routing.
- `skills/do-work/actions/work.md` needs no edit. It already consumes canonical selector order. Keeping it out avoids another semantic restatement.

## Revised builder manifest

Expected tracked source/test paths:

```text
skills/do-work/actions/capture-reference.md
skills/do-work/actions/capture.md
skills/do-work/actions/work-reference.md
skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go
skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go
skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go
skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go
skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go
skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go
skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go
skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go
skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go
skills/do-work-board/actions/board.md
skills/do-work-board/tools/queue-kanban/model.go
skills/do-work-board/tools/queue-kanban/model_test.go
skills/do-work-board/tools/queue-kanban/generate.go
skills/do-work-board/tools/queue-kanban/generate_test.go
skills/do-work-board/tools/queue-kanban/serve_test.go
skills/do-work-board/tools/queue-kanban/priority_browser_probe_test.go
skills/do-work-board/tools/queue-kanban/web/board-cards.js
skills/do-work-board/tools/queue-kanban/web/board.css
```

`clipboard_browser_probe_test.go` is removed from the expected manifest unless implementation proves that priority legitimately changes its copy-all fixture. `browser_probe_test.go` is reused, not edited. No builder change belongs under tracked `do-work/`, version/changelog, generated `build/`, screenshots, primes, lessons, or `work.md`.

## Integration and conflict findings

- The REQ-561 builder worktree exists on branch `worktree-agent-REQ-561-add-a-three-value-priority-field-the-selector-orders-by-and-the-board-shows` at claim commit `4e765265`. Main has advanced through planning and the REQ-502 checkpoint-cleanup hand-back. Before final integration, verify the builder branch manifest against its merge base and resolve from current main; do not assume the claim-time snapshot is current.
- REQ-502 (checkpoint cleanup) changed `internal/cleanup/cleanup_plan.go`, its test, and `internal/requeststate/state_apply.go`. It does not overlap REQ-561's proposed files, but it may advance the release baseline before REQ-561 lands.
- `work-reference.md` and `internal/nextselection` are likely overlap points with the REQ-503-to-REQ-510 lifecycle chain. Integrate REQ-561 serially and rerun both sides' tests after any merge resolution. In particular, REQ-505/REQ-510 are likely to alter selection/reference ownership and must not drop either priority axis or reintroduce duplicated prose.
- Current suite version is `0.272.1`. This feature is a minor user-visible release, but calculate the next minor from the actual integration state rather than hard-coding `0.273.0`. Keep root and installed changelogs byte-identical.
- The 23:20 report currently resolves to 27 still-pending `now` files: REQ-475, 483, 485, 490, 496, 503-510, 512, 514, 515, 527, 534-536, 539, 542, 544, 545, 547, 559, 560. It resolves to 12 still-pending `later` files: REQ-482, 486, and REQ-549 through REQ-558. Re-read at stamp time.
- REQ-502 currently sits in `do-work/queue/` as `pending-heavy-testing`, owned by its separate lane. Do not stamp it. The stamp authority is still “capture-editable pending files only,” not “anything under queue.”
- REQ-530 (newest-ready ordering) is already archived cancelled. Do not run cancel again. The orchestrator owns the post-landing provenance amendment/refusal record.
- Land executable core schema, selector, board parser/UI, and queue stamps in one integrating release increment. Priority-bearing queue files must not appear at a revision whose readers do not recognize them. Stage every stamped queue path explicitly.

## Verification focus

The plan's focused Go, heavy command, vet, full-module, JavaScript, strict Chromium, and repository gates are appropriate. The highest-risk assertions are:

1. both `selection_priority` and `priority` exist on selected and excluded JSON records and never overwrite each other;
2. invalid present priority emits exactly one core `SCHEMA-FALLBACK` and one board warning, while absence emits none;
3. repair and ready deferred-parent classes stay above ordinary `now`;
4. dependency gating happens before priority and fan-out happens after priority;
5. explicit REQ order stays caller-authored, while each UR expansion receives internal priority ordering;
6. Ready and Waiting remain separate, stable priority-sorted groups, and `Pending` is their exact concatenation;
7. static and live payloads agree, only `now`/`later` badges render, and light/dark measurements name the exact inspected URL and browser build.

Guidance read: the request, brief, plan, `CLAUDE.md`, action and board primes plus full lesson satellites, do-work-cli and queue-kanban shipped primes plus full lesson satellites, release prime, and general, coding-guardrails, backend, testing, and communication-style crew files.
