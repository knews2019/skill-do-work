# REQ-508 Exploration — Capture Templates and Schema-Backed Authoring

## Executive finding

The requested boundary is implementable without a new manifest discriminator and without a shell predicate. `capture-files` already receives the final UR/REQ bytes and parses them through the lossless `requestmodel` before planning any mutation. The parser exposes exact source keys, raw scalar spelling, decoded scalar values, list-vs-scalar evidence, duplicate counts, and source lines. That is enough to enforce canonical keys/values, list shape, UTC whole-second timestamps, safe title encoding, blocked-field pairing, and exact UR membership at the publication boundary.

The machine cannot determine whether user intent is Simple or Complex, whether TDD/maintenance/impact/priority/effort/assignment applies, which dependencies or lessons are relevant, or which optional body sections are warranted. Those judgments must remain in `capture.md` and the named prose contracts. The shortened examples may demonstrate the result, but `capture-files` must not infer it.

The captured `_dev/tests/contract-regressions.sh` write-set entry is stale. That file is now a 76-line dispatcher. No active shell predicate under `_dev/tests/` quotes the four examples or their comments. `_dev/tests/contracts/core-checks.sh:517-518` mentions `capture-reference.md` only as literal fixture data for the unrelated multi-path scope-drift parser. Neither shell file should change.

## Verified current surfaces

- Baseline focused packages pass: `go test ./internal/schemanormalization ./internal/requestmodel ./internal/publication` (exit 0; publication completed in 17.739s).
- Current example fences in `capture-reference.md` are exactly:
  - Simple REQ: lines 84-136, 53 lines including fence delimiters.
  - Complex REQ additions: lines 146-179, 34 lines including delimiters.
  - UR `input.md`: lines 220-239, 20 lines including delimiters.
  - Addendum REQ: lines 258-287, 30 lines including delimiters.
  - Total: 137 lines. “Under half” therefore means at most 68 total fenced lines.
- The five documented key-alias families are `addendum_to`, `depends_on`, `batch`, `related`, and `suggested_spec`. Only `depends_on` and `addendum_to` are currently projected through `requestmodel`; their alias lists are restated at `request_model.go:282-283`. `related` reads only its canonical key, while `batch` and `suggested_spec` are not projected.
- `schemanormalization.NormalizeField` already owns the canonical enum values, value aliases, defaults, warnings, and read-time status predicates. It does not report whether a recognized value was the exact canonical authored spelling. For example, `domain: back_end`, `route: b`, and `tdd: yes` are all recognized reads.
- `requestmodel.FieldEvidence` already carries `FieldName`, `LineNumber`, `RawValue`, `ScalarValue`, `ListValues`, `NestedValues`, and `DuplicateCount`. `ListValues != nil` distinguishes a flow/indented list (including `[]`) from a scalar. `RawValue` retains single quotes or a literal-block indicator, while `ScalarValue` contains decoded text. `ParseTimestamp` deliberately accepts legacy date-only, local, and offset shapes; `CanonicalTimestamp` emits UTC whole seconds.
- `BuildCapturePlan` already validates manifest IDs, canonical destination paths, payload regularity, parseable frontmatter fences, UR/REQ linkage, one-way UR membership, reservation spelling/collisions, raw-input containment, asset location, fold preimages, topology, and create collisions. These checks happen before apply, although the UR mutation is appended to the in-memory plan before all later validation completes.

## Smallest schema/requestmodel API change

Keep ordinary reads backward-compatible and add only evidence that the publication validator needs:

1. Add one schema-owned key-alias registry for all five families and expose a defensive-copy accessor such as `FieldKeyAliases(canonicalKey string) []string`.
2. Make `requestmodel` obtain aliases from that registry instead of spelling them at call sites. Project canonical/alias-backed `related`, `batch`, and `suggested_spec` values plus their selected source key, following the existing `DependsOn`/`DependencySource` pattern. Canonical-key presence must continue to win even when its value is empty.
3. Extend `schemanormalization.FieldResult` with exact-authoring evidence (for example `IsCanonicalValue`). It is true only when a present decoded enum/boolean value already equals the canonical spelling; it is false for aliases, case normalization, invalid/defaulted values, and absent defaults. `NormalizeField` must otherwise preserve its existing resolved values and warnings byte-for-byte.

Publication can use `RequestDocument.FieldValue` directly for key presence, raw scalar style, duplicates, and list shape, and `TypedRecord`/schema evidence for normalized values. No YAML dependency, manifest `record_kind`, or parser rewrite is required.

## Exact capture-publication validation boundary

Run these validations after both payloads are read and parsed but before the returned plan can be runnable:

### UR payload

- Require canonical `id`, `title`, `created_at`, `requests`, and `word_count` keys; reject duplicate effective keys and any documented read-only alias key.
- Require `id` to match the manifest (already enforced), `title` to be a non-empty safely encoded user-text scalar, `created_at` to parse as exactly `YYYY-MM-DDTHH:MM:SSZ`, `requests` to be a list, and `word_count` to be a non-negative integer scalar.
- Require the `requests` list to equal the manifest request-ID set exactly: no missing member, phantom member, or duplicate. Preserve `requests: []` for a legitimate fold-only capture.
- Keep raw-input byte containment and control-byte validation under their existing specific refusal codes. Do not infer or validate Summary, Extracted Requests, Batch Constraints, Folded Requests, or any prose quality.

### REQ payload

- Require canonical `id`, `title`, `status`, `created_at`, `user_request`, `domain`, `prime_files`, `tdd`, and `maintenance` keys. Keep `impact`, `priority`, `effort_estimate`, `assigned_to`, `suggested_spec`, `depends_on`, `required_lessons`, `related`, `batch`, `write_set`, `addendum_to`, and blocked metadata optional because absence is meaningful under the prose contracts.
- Reject every registered alias key in a newly authored payload while preserving all aliases for ordinary reads.
- Require any present schema enum/boolean field to carry a canonical value. In particular, refuse value aliases such as `back_end`, `yes`, `test_first`, `trivial`, or lowercase route values even though normal readers continue accepting them.
- Require list shape, not scalar coercion, for any present `prime_files`, `required_lessons`, `depends_on`, `related`, and `write_set`. Preserve empty lists. Do not validate path relevance, lesson budgets, dependency correctness, or write-set completeness here.
- Require `created_at` and any present `blocked_at` to be exact UTC whole-second instants. Do not enforce “current” wall-clock proximity; that is not stable authoring shape and future-stamp diagnosis already has a separate owner.
- If normalized status is `blocked`, require non-empty `blocked_by` and `blocked_at`; if status is not blocked, refuse capture-authored `blocked_by`, `blocked_at`, or `blocked_check` as incoherent. `blocked_check` remains optional for blocked records.
- Require `title` and any present user-text fields (`blocked_by`, `blocked_check`, `assigned_to`) to use the Frontmatter Quoting contract: a lexically valid single-quoted scalar with doubled internal apostrophes, or a supported literal block for LF-bearing text. Reject plain and hand-written double-quoted user text. The existing `RawValue` plus decoded `ScalarValue` is sufficient evidence; a local publication helper should check valid quote pairing rather than treating matching first/last quote bytes as sufficient.
- When `impact` is present, validate only the emitted field/title mirror: non-default impact requires the matching leading `[impact-*] ` tag, while the default/absent value must not carry a contradictory impact tag. Do not choose the impact verdict.
- Keep existing path/linkage/reservation/raw/fold refusal codes more specific than the new schema refusal. New record-shape failures should use stable `CAPTURE-UR-SCHEMA-INVALID` or `CAPTURE-REQ-SCHEMA-INVALID` codes with record ID/path and one concrete field reason.

Do not require or interpret Complex-only body sections, Red-Green quality, Open Questions, Prior Implementation applicability, a dependency graph, a lesson match, `impact`, `priority`, or `effort_estimate`. A common non-empty Markdown body may be retained as template shape, but section-level enforcement would encode the missing Simple/Complex/Addendum judgment and should not be added in this REQ.

## Removable template-rule map

| Current fence annotation | After removal | Behavior proof |
|---|---|---|
| Simple title quoting and impact prefix (line 87); UR title quoting (223); Addendum title/tag composition (261) | One pointer outside the fences to Frontmatter Quoting and REQ Title Convention; capture validation owns emitted syntax/mirror | RED: plain/double-quoted/invalid-apostrophe title and mismatched impact tag refuse. GREEN: doubled-apostrophe single quote and literal-block user text pass. Existing board quoting test remains compatibility evidence, not the capture validator. |
| `created_at` comments (89, 224, 263) and `blocked_at` comment (105) | One pointer outside the fences to the Timestamp rule; publication owns canonical authored syntax | RED: date-only, local-with-`Z`, offset, fractional, and malformed values refuse. GREEN: valid whole-second UTC passes. Existing `ParseTimestamp` tests continue proving legacy reads. |
| Domain enum (91), TDD value shape (94), maintenance value shape (97), status/impact/priority/effort value spellings (98-103) | Schema registry is the sole mechanical value source; capture Step 1 retains every applicability/verdict judgment | Table-test every canonical value/alias/default and exact-authoring evidence; publication refuses aliases/invalid values while ordinary reads preserve resolution/warnings. |
| `prime_files`, `required_lessons`, `depends_on` list comments (92-93, 96) | Minimal examples show list shape; named Required Lessons/depends_on contracts remain prose | Publication tests scalar refusal and empty/populated-list acceptance. No test infers relevant primes, lesson selection/budget, dependencies, or TDD write-set completeness. |
| `suggested_spec` comment (95) | Optional minimal field may be omitted; Step 1 retains the matching judgment | Alias registry/requestmodel tests canonical precedence for `suggested_spec`/`spec_hint`/`suggested-spec`; publication rejects alias keys on write. No test decides a spec. |
| `impact`, `priority`, `effort_estimate`, `assigned_to`, external-condition and `blocked_check` applicability comments (98-106, 266-267) | Consolidated judgment pointer outside fences plus Step 1 remains authoritative | Only emitted shape is tested: canonical optional values, safe user text, blocked pairing. Absence stays valid; no machine verdict or assignment/probe invention. |
| Addendum arrows/comments for new `user_request` and `addendum_to` (264-265) | Unannotated canonical example plus Addendum prose immediately above the fence | Existing linkage tests plus alias-key refusal; positive Addendum fixture with exact manifest/UR/addendum linkage. Do not prove archived-vs-in-flight Prior Implementation judgment mechanically. |
| Complex fence instructions for Detailed Requirements, Constraints, Dependencies, Builder Guidance, Open Questions, and Full Context (147-178) | Short copyable headings/placeholders; rules remain in the prose immediately after/outside the example | Example-fixture parse succeeds. No Go test judges completeness, wording, recommended answers, or whether sections apply. |
| UR sample's user-text/raw-input commentary and Addendum body guidance | One prose pointer to Outside-text containment/Addendum contract | Retain existing byte-containment tests; positive ordinary, fold-only, and Addendum capture fixtures. |

There are no active comment-pinning predicates to delete. Record the predicate leg as an evidence-backed no-op rather than adding a replacement sentence grep.

## Compatibility and fixture hazards

- Do not make the new check a general `requestmodel.ParseDocument` rule. Queue/archive readers deliberately accept legacy aliases, malformed-but-recoverable titles, old timestamps, scalar dependency shapes, and noncanonical enum aliases. Canonicality is a `capture-files` write-boundary property only.
- Do not change `ParseTimestamp` to reject legacy shapes. Add a separate canonical authored-timestamp predicate or compare a strict parse/format at publication.
- Do not infer canonicality from `FieldResult.IsRecognized`: aliases currently return recognized=true. Exact-authoring evidence must distinguish alias/case normalization from an exact canonical value.
- Do not inspect only `TypedRecord.DependsOn` for list shape; `requestmodel.listValue` intentionally coerces a scalar to a one-element slice. Use `FieldEvidence.ListValues != nil`.
- Do not compare UR membership only as a set. A set catches phantoms/missing members but hides duplicate IDs; reject duplicates explicitly.
- Do not require optional judgment fields merely because the current annotated example shows them. `impact` and `effort_estimate` are expected but explicitly allowed absent when judgment was impossible; `priority`, assignment, spec hints, dependencies, lessons, and blocked probes are conditional.
- Do not require all `*_at` fields on a newly captured REQ or validate the timestamp against wall-clock now. Phase stamps belong to later lifecycle writers.
- `BuildCapturePlan` is also called by `BuildAnswerPlan` for `override_capture`. The structured override fixture at `publication/answer_test.go:125-129` currently uses minimal UR/REQ frontmatter and will fail strict capture validation. Update that fixture only if the shared validator makes it fail; retain the answer behavior and prefixed refusal-code assertion.
- Every existing `capture_files_test.go` happy-path fixture is similarly underspecified today (typically only UR `id`/`requests` and REQ `id`/`status`/`user_request`). Since that test file is already primary scope, centralize canonical fixture builders there and keep each existing test focused on its original path/reservation/raw/fold behavior.
- `publication_manifest_test.go` tests strict JSON decoding only and does not call `BuildCapturePlan`; it should not need edits. `publication_types.go` and `publication_manifest.go` need no new `record_kind` or manifest field.

## Exact write set

Primary:

1. `skills/do-work/actions/capture-reference.md`
2. `skills/do-work/actions/capture.md`
3. `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go`
4. `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go`
5. `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go`
6. `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go`
7. `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go`
8. `skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go`

Conditional fixture-only scope:

9. `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` — only to canonicalize the existing structured override payload after the shared `BuildCapturePlan` validator makes the fixture fail.

Explicitly excluded: `_dev/tests/contract-regressions.sh`, `_dev/tests/contracts/core-checks.sh`, publication manifest/types production files, `work-reference.md`, `docs/capture-guide.md`, and the board package. Any newly discovered need outside the nine files requires owner-recorded scope expansion before implementation.

## Public RED/GREEN cases

1. **Alias keys:** new REQ payloads using each of `amends`/`parent`/`amendment_to`, `dependencies`, `batch_name`, `related_reqs`, `spec_hint`, or `suggested-spec` refuse; `requestmodel` ordinary reads still select each alias and canonical presence still wins.
2. **Alias/invalid values:** `domain: back_end`, `route: b`, `tdd: yes`, `maintenance: on`, and `effort_estimate: trivial` refuse as capture writes; normalization tests still resolve them for legacy reads with unchanged warnings/defaults.
3. **List shape:** scalar `prime_files`, `required_lessons`, `depends_on`, `related`, `write_set`, or UR `requests` refuses; `[]`, flow lists, and indented lists pass.
4. **Timestamp shape:** date-only, offset, local, fractional, calendar-invalid, and malformed authored timestamps refuse; `2026-09-04T19:00:00Z` passes. Legacy `ParseTimestamp` acceptance remains green.
5. **User-text scalar:** plain or hand-double-quoted title, an undoubled apostrophe, and unsafe blocked/assignment text refuse; single-quoted apostrophes and literal blocks pass.
6. **Blocked pairing:** `status: blocked` without `blocked_by` or `blocked_at` refuses; an ordinary pending REQ requires none; a blocked REQ without `blocked_check` remains valid.
7. **Impact mirror:** a non-default emitted impact with a missing/wrong title tag refuses; default/absent impact with no tag passes. The validator never selects impact.
8. **Membership:** missing manifest member, phantom UR member, and duplicate UR member refuse; exact ordinary membership and `requests: []` fold-only membership pass.
9. **Four examples:** canonical manifests constructed from the shortened Simple, Complex, UR, and Addendum examples pass, while Simple/Complex distinction is not inferred.
10. **Existing safeguards:** path, linkage, reservation, raw containment, assets, fold stale-preimage, mutation order, collision, topology, and rollback tests keep their specific outcomes after fixture canonicalization.

## Commands

Focused RED/GREEN:

```bash
cd skills/do-work/tools/do-work-cli && go test ./internal/schemanormalization ./internal/requestmodel ./internal/publication
```

Concurrency-sensitive focused check:

```bash
cd skills/do-work/tools/do-work-cli && go test -race ./internal/schemanormalization ./internal/requestmodel ./internal/publication
```

Module suite:

```bash
cd skills/do-work/tools/do-work-cli && go test ./...
```

Action/shipped-contract checks:

```bash
bash _dev/tests/contract-regressions.sh
bash _dev/tests/staged-skills-contract.sh
```

Canonical repository gate, direct and unpiped:

```bash
bash _dev/tests/maintainer-verify.sh
```

Final manual audit: recount the four exact fenced spans and require a combined total of at most 68 lines; inspect only those spans for explanatory inline `#`/arrow annotations; confirm all four named examples remain copyable; and rerun a repository search showing no new shell sentence predicate was introduced.

*Generated by Explore agent*
