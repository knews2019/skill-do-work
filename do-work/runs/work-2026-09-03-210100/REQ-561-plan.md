# REQ-561 implementation plan

## Contract and scope decisions

1. Add `priority` as an optional Schema Read Contract enum with canonical values `now`, `next`, and `later`, no aliases, and fallback `next`. Absence is recognized and resolves to `next`; a present off-vocabulary value resolves to `next` and emits the canonical schema warning. Preserve the original scalar and recognition evidence wherever other contracted fields do.
2. Keep the existing `selection_priority` output exactly as the selection-class signal (`repository-gate-repair`, `deferred-parent`, or `ordinary`). Do not overload it with the authored request priority. Add a separate `priority`/`RequestPriority` projection so callers can consume both dimensions without parsing display text.
3. Apply authored priority only inside the ordinary ready class: class rank remains repository-gate repair, then dependency-ready deferred parent, then ordinary; only ordinary records receive the `now`/`next`/`later` secondary key. Preserve the selector's pre-existing queue order as the final stable tie-break. Explicit REQ targets preserve caller order. UR-expanded ordering keeps its current dependency-depth/id tie-break after the class and authored-priority keys.
4. Dependency readiness remains a filter before ordering and fan-out bounding. A `now` dependent with a pending `later` prerequisite is excluded and the ready prerequisite is selected; priority never changes dependency state or promotes the blocked dependent.
5. The board uses the same effective values but retains its intentionally independent parser table. It sorts only the Pending Ready and Pending Waiting groups, independently, by `now`, `next`, `later`, retaining the model's existing numeric/queue order for ties. It does not reorder Claimed, Needs input / Blocked, Recently done, or completion anomalies. `Pending` remains the compatibility union of the two dependency groups in Ready-then-Waiting order.
6. Cards show a compact `now` or `later` badge. The default/fallback `next` has no badge, including an invalid raw value that resolved to `next`; invalid provenance remains visible through the existing board warning banner.
7. `work.md` needs no behavioral edit: it already commands the caller to process canonical `next` records in returned stable order. The selection contract belongs in `work-reference.md`; duplicating the sort algorithm in the procedure would create drift. Likewise no prime or lesson file needs a product edit unless implementation uncovers a genuinely reusable new lesson.
8. Generated sites, browser screenshots, fixture repositories, and build outputs are evidence only and remain untracked. The builder does not edit any tracked `do-work/` file in its worktree.

## Requirement-to-file map

### Schema, typed model, and selector

- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go`
  - Add the `priority` contract row: canonical `now`, `next`, `later`; empty aliases; default `next`.
  - Let the existing `NormalizeField`, `SchemaFieldNames`, and warning machinery remain the single implementation of absence/invalid behavior.
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go`
  - Add table rows for absent, `now`, `next`, `later`, case/whitespace normalization if supported by the common contract, and an invalid token. Assert resolved value, recognition bit, original value, and warning presence/absence.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go`
  - Add `RequestPriorityValue string` and `RequestPriorityEvidence schemanormalization.FieldResult` to `RequestRecord` (the `Request` prefix avoids confusion with the selector's class priority).
  - Normalize `document.scalarValue("priority")` in `TypedRecord` and populate both fields. This is required by the typed-projection inventory guard; relying only on `FieldEvidenceByName` would leave a contracted field unprojected.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go`
  - Add a `priority:` fixture value and the reflection-table row mapping schema field `priority` to `RequestPriorityValue`/`RequestPriorityEvidence`. Keep the equality assertion against `SchemaFieldNames()` so later schema additions cannot bypass typed projection.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
  - Add `RequestPriority string \`json:"priority"\`` to both `SelectionRecord` and `SelectionExclusion`; keep `SelectionPriority` and its JSON name unchanged.
  - In `NormalizeResult`, default an empty request priority to `next`, matching the schema default and giving typed target-resolution exclusions a deterministic value when no request record exists.
  - Copy request priority anywhere result records are normalized or converted; do not infer it in renderers from class names, reason strings, or IDs.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
  - Assert selected and excluded records normalize/project `priority`, default to `next` when absent, and leave `selection_priority` unchanged.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go`
  - Add request-priority constants/ranking separate from the existing selection-class constants. Naming must make the distinction explicit (`RequestPriorityNow`, `RequestPriorityNext`, `RequestPriorityLater`).
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
  - Populate `SelectionRecord.RequestPriority` from `candidate.RequestFile.TypedRecord.RequestPriorityValue` in `evaluateCandidate` before any exclusion can be produced.
  - Copy it through `copySelectionEvidenceToExclusion`; initialize direct exclusions through `exclusionFor` to the schema default. Every exclusion created from a real request must carry that request's normalized value, including dependency, assignment, status, skip, and `FAN-OUT-LIMIT` paths.
  - Add priority evidence to `appendSchemaWarnings`, producing `SCHEMA-FALLBACK` for invalid present values through the existing warning message.
  - Extend the default eligible stable comparator to class rank first, then request-priority rank only when both records are ordinary. Return equality for equal keys so the pre-sorted queue order remains the tie-break. Apply the fan-out limit only after this sort.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go`
  - Extend UR-member stable ordering with request priority only inside the ordinary class, before the existing dependency-depth and request-id keys. Leave explicit REQ token accumulation/deduplication order untouched.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go`
  - Add the table-driven selector cases described below. Include class-precedence, stable tie, dependency, fan-out-limit, explicit-target, and UR-expansion assertions so every selection path has a pinned ordering rule and both output priority fields are distinguished.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go`
  - Extend the command-level decoded JSON fixture struct with `RequestPriority string \`json:"priority"\`` and assert the actual CLI JSON projects it on selected records and fan-out exclusions. Run this caller seam under its existing `DO_WORK_HEAVY_TESTS=1` opt-in.

### Board model, payload, and client

- `skills/do-work-board/tools/queue-kanban/model.go`
  - Add board constants for `now`, `next`, `later` and the matching row in `schemaReadContractFields`.
  - Add `Priority`, `OriginalPriority`, and `PriorityUnrecognized` to `RequestTicket`.
  - At the request parse site, always call `resolveSchemaField("priority", raw)` so absence becomes the effective `next` needed for sorting; unlike older display-only fields, do not leave the resolved value empty when absent. Retain original/invalid evidence.
  - Add priority to `collectSchemaFieldWarnings` so an invalid present scalar enters `Board.Warnings` with the shared canonical wording.
  - After dependency annotation and bucketing, stable-sort `PendingReady` and `PendingWaiting` independently by priority rank. Rebuild `Pending` as Ready followed by Waiting so its documented compatibility-union invariant and the UI grouping agree. Do not put the comparator in client JavaScript.
- `skills/do-work-board/tools/queue-kanban/model_test.go`
  - Add a table-driven parser test for absent, each canonical value, and invalid. Assert effective/original/unrecognized fields plus warnings.
  - Add a model-order fixture whose numeric order deliberately conflicts with authored priority in both Ready and Waiting groups. Include equal-priority pairs to pin the existing queue-order tie, and assert the full `Pending` union remains Ready then Waiting.
- `skills/do-work-board/tools/queue-kanban/generate.go`
  - Add `Priority`, `OriginalPriority`, and `PriorityUnrecognized` to `generatedRequest` with JSON names `priority`, `originalPriority`, and `priorityUnrecognized`; project the fields from `RequestTicket`.
  - Project the already sorted model slices into `columns.pendingReady`, `columns.pendingWaiting`, and `columns.pending` without a second sort. This keeps static generation and the live server on their shared `buildBoard` -> `buildGeneratedBoardData` path.
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
  - Add an end-to-end generation test that builds real request files, decodes the generated `board-data.js`, and asserts effective priority values, invalid provenance/warning, and conflicting-ID order in all three Pending slices. This prevents a test-only hand-built model from hiding parse/projection drift.
- `skills/do-work-board/tools/queue-kanban/serve_test.go`
  - Build the same kind of fixture through the live server and fetch `/board-data.js`; assert the priority payload and Pending Ready/Waiting ordering equal the static projection. Keep the existing shared-client/static-live asset identity assertion intact.
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js`
  - In the existing card badge/chip construction path, append one semantic priority badge only for `request.priority === "now"` or `"later"`; emit nothing for `next` or missing legacy payloads. Use the effective value, while data warnings continue to expose invalid raw input.
- `skills/do-work-board/tools/queue-kanban/web/board.css`
  - Add a small `.badge-priority` base and value-specific modifiers for now/later using existing theme variables/tokens. Verify both themes for contrast, wrapping, alignment, and no card-width regression.
- `skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go`
  - Add or extend a strict Chromium probe over a generated static fixture. Assert DOM card order separately in Ready and Waiting, exactly one visible `now` and `later` badge per applicable card, and zero `next` badges. Return `location.href` in the same browser evaluation as the DOM evidence. Keep clipboard expectations coherent if the existing fixture's order is changed; prefer a dedicated priority fixture to avoid coupling unrelated clipboard assertions.

### Action contract and capture behavior

- `skills/do-work/actions/work-reference.md`
  - Add `priority: next` to Full Frontmatter as optional authored ranking, with closed values/default/warning and its selector/board roles.
  - Add the Schema Read Contract row with no aliases and default `next`.
  - Rewrite Selection Order to state class order, ordinary priority order, queue-order tie, unchanged dependency gate, unchanged explicit caller order, UR-expanded behavior, and the two distinct typed fields: `selection_priority` (class) and `priority` (authored rank).
- `skills/do-work/actions/capture-reference.md`
  - Add only a commented optional template line such as `# priority: now` with all three values and the “user's words only” constraint. Add `priority` to the list of enum fields governed by normalize-and-warn.
- `skills/do-work/actions/capture.md`
  - Add a Step 1 priority assessment: emit `priority` only when the user's words explicitly rank timing/order, map that language to `now`/`next`/`later`, and otherwise omit it. Never infer from impact, effort, dependency depth, recency/REQ number, list position, or the default.
  - In the queued-addendum path, allow later user wording to set, change, or remove priority while preserving the same evidence rule; do not silently retain an explicitly withdrawn rank.
  - Add a checklist item requiring traceability from every emitted priority to user ranking words.
- `skills/do-work-board/actions/board.md`
  - Update the parser lock-step paragraph to name priority's effective default/warning, Ready/Waiting-only ordering role, and now/later/no-next rendering. State explicitly that it does not change status or dependency routing.
- `skills/do-work/actions/work.md`
  - No edit planned. Its “Process selected REQs in the returned stable order” sentence already delegates authority to the Go selector.

## Ordered RED/GREEN implementation

### 1. RED: pin the unknown field, ignored selector order, and missing board behavior

Write tests first without referencing not-yet-created Go struct fields directly where that would create a compile failure. Initial RED must be an assertion failure against present behavior:

1. In `schema_normalization_test.go`, call `NormalizeField("priority", ...)` table-wise and assert the new contract. Today absence/invalid handling and the contract inventory fail.
2. In `request_model_test.go`, add the projection inventory row by reflection. Today the typed fields are absent and the inventory differs; the test reports behavioral projection omissions instead of failing compilation.
3. In `next_selection_test.go`, create ordinary ready records in numeric order `later`, `next`, `now` and request fan-out large enough to observe the whole sequence. Assert order `now`, `next`, `later` and inspect records through JSON/reflection or maps for `priority`. Today stable numeric order wins and the field is absent. Also add:
   - absent/each canonical/invalid table cases, invalid `SCHEMA-FALLBACK`, and normalized `next` projection;
   - older `next` and newer `now` so the fixture cannot accidentally pass under existing queue order;
   - repair and ready deferred-parent records with `later` still preceding an ordinary `now`;
   - equal-priority records preserving numeric queue order;
   - explicit tokens preserving caller order even when priority conflicts;
   - UR-expanded members ordered by class, ordinary priority, then existing depth/id;
   - fan-out exclusions carrying `priority` after bounding.
4. Add the dependency case with a pending `later` prerequisite and a `now` dependent. Assert the prerequisite is the selected ready record and the dependent is excluded `DEPENDENCIES-UNMET`; assert both typed records retain their authored priorities. This test should already pass the gating part but fail projection until the feature exists, proving the new sort never bypasses readiness.
5. In board tests, use real fixture files with IDs deliberately opposite their priorities. Assert effective parse, Ready/Waiting order, generated JSON, and DOM tags. Today the board retains numeric order, has no priority payload, and renders no tags.

Run and retain exact failing assertion output for the hand-back:

```sh
cd skills/do-work/tools/do-work-cli
go test ./internal/schemanormalization ./internal/requestmodel ./internal/resultmodel ./internal/nextselection

cd ../../../../do-work-board/tools/queue-kanban
go test -count=1 -run 'Test.*Priority' .
QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 go test -count=1 -run 'TestBrowserBehavior.*Priority' .
```

### 2. GREEN: implement from the schema outward

1. Add the core schema row and typed `RequestRecord` projection; make the normalization/requestmodel tables green first.
2. Add the separate result-model field/default and nextselection projection on every selected/excluded path. Make projection tests green before changing order.
3. Add the ordinary-class comparator and UR-expansion comparator, preserving class/dependency/explicit-target contracts; make focused selector and command-level JSON tests green.
4. Add the independent board schema row, parse evidence, warnings, and stable Pending group ordering; make model tests green.
5. Add generated payload fields, then the card badge and CSS. Make static projection, live server, and strict browser behavior tests green.
6. Update the action/reference prose only after executable behavior is green, so documentation records the proven contract.
7. Review every changed file, run formatting (`gofmt` on changed Go files), and confirm no generated artifact or worktree `do-work/` path entered the branch.

## Verification matrix

### Focused and package checks

```sh
cd skills/do-work/tools/do-work-cli
go test -count=1 ./internal/schemanormalization ./internal/requestmodel ./internal/resultmodel ./internal/nextselection
DO_WORK_HEAVY_TESTS=1 go test -count=1 ./internal/nextselection
go vet ./...
go test -count=1 ./...

cd ../../../../do-work-board/tools/queue-kanban
go test -count=1 .
go vet ./...
QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 go test -count=1 -run 'TestBrowserBehavior.*Priority' .
```

For the actual-browser lane, use the repository's supported pinned/current stable Chromium executable and strict markers already used by browser probes:

```sh
QUEUE_KANBAN_BROWSER='<exact-chromium-path>' QUEUE_KANBAN_BROWSER_PROBES=on QUEUE_KANBAN_STRICT_BROWSER_BEHAVIOR=1 go test -count=1 -run 'TestBrowserBehavior.*Priority' .
```

### Repository checks after serial integration

Run from the main-tree repository root, unpiped where the action requires direct gate attribution:

```sh
bash _dev/tests/contract-regressions.sh
bash _dev/tests/shipped-package-reference-contract.sh
bash _dev/tests/maintainer-verify.sh
```

The orchestrator decides whether its earlier one-off `--heavy` evidence remains valid after this feature lands or whether the canonical final gate must rerun; the builder must report, not invent, that integration decision.

## Visual and DOM verification

1. Create an untracked temporary fixture repository containing at least Ready and Waiting cards for all three priorities, an absent field, and an invalid field. Choose IDs so numeric order conflicts with desired order and include long titles/badges to exercise wrapping.
2. Generate a static board into a temporary output directory and open its exact `file://.../index.html` URL in the supported Chromium build. Capture a screenshot in light mode and another in dark mode. Inspect tag text, contrast, spacing, wrapping, and badge alignment; confirm no overlap or horizontal clipping.
3. In the same DOM evaluation that gathers evidence, return `location.href`, the card ID order under Pending Ready and Pending Waiting, priority badge text by card, and computed style/rect data. Required result: each group is now/next/later with stable ties, only now/later have tags, and the returned URL exactly names the inspected page.
4. Serve the same fixture over loopback using the live board, inspect the exact `http://127.0.0.1:<port>/...` URL, and return `location.href` plus the same DOM evidence. Confirm static and live pages use identical ordering and tag behavior.
5. After the queue stamp lands, the orchestrator regenerates the real repository board and repeats exact-URL DOM evidence against the actual Pending groups, checking all stamped `now` records precede default `next`, which precede stamped `later`, without moving waiting work into Ready. Do not commit screenshots or generated board files.

## Builder branch manifest expected

The builder may change only the source/action/test files justified above. Expected branch paths are:

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
skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go
skills/do-work-board/actions/board.md
skills/do-work-board/tools/queue-kanban/model.go
skills/do-work-board/tools/queue-kanban/model_test.go
skills/do-work-board/tools/queue-kanban/generate.go
skills/do-work-board/tools/queue-kanban/generate_test.go
skills/do-work-board/tools/queue-kanban/serve_test.go
skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go
skills/do-work-board/tools/queue-kanban/web/board-cards.js
skills/do-work-board/tools/queue-kanban/web/board.css
```

`work.md`, the prime/lesson files, version/changelogs, all queue/archive/working records, the run artifacts, generated sites, and screenshots stay out of the builder commit. The hand-back must state the actual full manifest and explain any deviation.

## Orchestrator-only integration seams

### Serial merge and release

1. Re-read the hand-back, verify its branch/commit and full manifest, and reject any tracked `do-work/` mutation from the builder worktree. Merge this feature serially; do not expose stamped queue records before schema, selector, and board readers exist.
2. Resolve overlaps with other in-flight changes (`work-reference.md` and `internal/nextselection` are known shared seams) against the current integration tree, then rerun focused tests and the canonical repository gate.
3. Treat this as a user-visible feature release. The current suite version is `0.272.1`; absent a newer concurrent release, bump `skills/do-work/actions/version.md` to `0.273.0`. Add a changelog entry whose title says the shipped behavior (for example, “Priority-ordered queue selection and board tags”) and keep root `CHANGELOG.md` and `skills/do-work/CHANGELOG.md` byte-identical. If integration has advanced the version, calculate the next appropriate release from that state rather than forcing this number.
4. Keep schema/parser/selector/board/stamp/release in one integrating commit, or in the canonical finalization transaction's one coherent landed increment, so no revision has priority-bearing queue files with old readers.

### Queue stamp from the latest velocity report

Immediately before stamping, re-read `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/index.html` and resolve each named REQ's current location/status. Edit only capture-editable files still under `do-work/queue/`; never edit `do-work/working/`, a foreign claim, or an already terminal record. Stage explicit paths, never `do-work/queue/` as a directory.

At plan time the report/current tree resolve to these pending files for `priority: now` (REQ-502 is currently claimed in its own lane and is intentionally absent):

```text
REQ-475 REQ-483 REQ-485 REQ-490 REQ-496
REQ-503 REQ-504 REQ-505 REQ-506 REQ-507 REQ-508 REQ-509 REQ-510
REQ-512 REQ-514 REQ-515 REQ-527 REQ-534 REQ-535 REQ-536
REQ-539 REQ-542 REQ-544 REQ-545 REQ-547 REQ-559 REQ-560
```

At plan time these pending files receive `priority: later`:

```text
REQ-482 REQ-486 REQ-549 REQ-550 REQ-551 REQ-552
REQ-553 REQ-554 REQ-555 REQ-556 REQ-557 REQ-558
```

All other records remain untouched. If REQ-502 returns to `do-work/queue/` before the stamp and is no longer claimed, it belongs to the report's build-now set and may then be stamped `now`; otherwise its lane owns it. Recompute from the report/current tree rather than treating this plan-time list as authority.

### REQ-530 cancellation provenance

REQ-530 is already archived with `status: cancelled` at `do-work/archive/UR-101/REQ-530-select-the-newest-ready-req-first.md` (cancelled at commit `f23c1c7f`). Therefore do not invoke the cancellation transaction a second time: the terminal archive is no longer a cancelable queue record. After the REQ-561 landing hash exists, the orchestrator must reconcile the request's “cancel against this landing hash” provenance by recording that exact hash in the existing REQ-530 cancellation reason/metadata using the canonical terminal-record amendment path, and stage that archive path explicitly. If the lifecycle refuses post-terminal amendment, record the refusal and exact landing hash in REQ-561's finalization/hand-back rather than fabricating a second cancellation; this discrepancy cannot be delegated to the builder.

### Final checks

- Verify JSON selected and excluded records contain both `selection_priority` and `priority` with their distinct meanings.
- Verify invalid priority produces one core `SCHEMA-FALLBACK` and one board data warning, while absence produces none.
- Verify gate repairs and ready deferred parents remain above ordinary `now` work.
- Verify priority changes ordering only after dependency gating and before fan-out bounding.
- Verify static and live board payload/order match, actual queue stamp order is visible, and no `next` tag renders.
- Verify explicit-path staging contains no foreign claimed file, generated output, screenshot, or unrequested queue record.
- Verify the builder hand-back contains branch, commit, full manifest, RED then GREEN outputs, browser evidence with exact URLs, lesson-read evidence, seams, decisions, and discovered tasks.

## Guidance read

Planning incorporated `CLAUDE.md`; the REQ brief and authoritative main-tree request; `_dev/primes/prime-action-files.md`; `_dev/primes/lessons-action-files.md`; `_dev/primes/prime-kanban-board.md`; `_dev/primes/lessons-kanban-board.md`; `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`; `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`; `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`; `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`; `_dev/primes/prime-releases.md`; and the general, coding-guardrails, backend, testing, and communication-style crew files. In particular, the plan preserves the typed schema-projection inventory, exercises real caller seams, mirrors the independent board parser in the same change, retains invalid raw provenance for warnings, uses ordering-conflict fixtures, and requires pixel plus same-evaluation exact-URL DOM evidence.
