# REQ-508 Technical Plan

## Outcome

Replace the four annotated capture templates with four coherent minimal examples, while moving their mechanical field contract behind the existing Go schema/publication boundary. Keep capture-time judgments in prose. Do not turn judgment calls into enum checks, and do not replace deleted sentence predicates with new sentence predicates.

The four record examples are:

1. Simple REQ.
2. Complex REQ (the Simple shape plus complex-only frontmatter and body sections).
3. UR `input.md`.
4. Addendum REQ.

The current fenced examples occupy 137 lines in total (Simple 53, Complex 34, UR 20, Addendum 30). The edited examples should total at most 68 lines, contain no per-field explanatory comments, and remain complete enough to copy as their named record shape.

## Current-Tree Findings

- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` owns enum values, value aliases, defaults, warnings, canonical status predicates, and a generic canonical-key-precedence helper. It does not yet own the five field-key alias sets printed in the reference table, and it has no write-time “canonical spelling only” predicate.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` restates the `depends_on` and `addendum_to` alias lists at its call sites. It projects `related` only from the canonical key and does not project `batch` or `suggested_spec`, so the reference table is not currently backed by one schema source.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go` already refuses unsafe/colliding destinations, noncanonical IDs and reservation paths, mismatched UR/REQ linkage, absent UR membership, invalid frontmatter fences, unsafe raw input, missing byte containment, invalid asset paths, and stale fold preimages. It currently accepts a capture-authored read alias such as `domain: back_end` or `dependencies: [...]`, accepts noncanonical timestamps and unsafe scalar spellings, and checks that each manifest REQ is in `requests` without rejecting phantom extra membership.
- The captured `_dev/tests/contract-regressions.sh` path is stale. It is now a 77-line dispatcher. A repository-wide search found no active shell predicate that quotes any of the four templates or their per-field comments. `_dev/tests/contracts/core-checks.sh` mentions `capture-reference.md` only inside an unrelated scope-drift fixture. Therefore this REQ should not edit the dispatcher or invent a replacement prose predicate; the deletion leg is already satisfied in the current tree and should be recorded as an evidence-backed no-op.

## Boundary: Mechanics Versus Judgment

### Schema/publication mechanics

These rules can leave the templates because the Go layer can enforce or expose them:

- canonical field keys and the read-only key aliases for `addendum_to`, `depends_on`, `batch`, `related`, and `suggested_spec`;
- canonical enum/boolean spellings, value aliases, and defaults already held by `schemanormalization`;
- required identity/linkage fields for a newly published UR/REQ, canonical REQ/UR IDs, exact manifest linkage, canonical queue/UR/reservation paths, and exact UR membership;
- list-versus-scalar shapes for list fields used by the minimal examples;
- canonical UTC whole-second syntax for `created_at` and, when present on a blocked capture, `blocked_at`;
- a safely encoded title/user-text scalar shape, using the existing single-quoted or literal-block forms rather than relying on parser recovery;
- capture-only canonical writing: a read alias is accepted by ordinary readers but refused in newly published payloads;
- paired blocked-record metadata when a new record is authored with `status: blocked`;
- raw-input byte containment and atomic publication, which already have behavior coverage.

### Prose judgments

These stay in `skills/do-work/actions/capture.md` or in the named prose contracts of `capture-reference.md`; the examples should point to them once instead of restating them per field:

- deciding Simple versus Complex and extracting all requirements;
- deciding whether TDD is realistic and selecting the RED/GREEN proof;
- judging `impact`, `effort_estimate`, `priority`, `maintenance`, `assigned_to`, and whether an external condition truly blocks work;
- choosing relevant prime/lesson paths and honoring the Required Lessons Budget Contract;
- deciding dependency edges, `write_set`, sweep conversion/folding, open-question wording, stakeholder audience, and whether archived prior implementation context is needed;
- composing titles and descriptions from user intent, including the non-default impact mirror;
- deciding whether optional sections such as Assets, Why, Red-Green Proof, Open Questions, and Prior Implementation apply.

The machine may validate the shape of a judgment once the author emits it; it must not make the judgment.

## Field-Rule Disposition

| Template annotation/rule | Durable owner after this change | Required proof |
|---|---|---|
| title quoting and safe user-text scalar | `capture-files` record validation, with `capture.md` retaining the judgment/composition rule | unsafe/unquoted syntax refuses; quoted apostrophes and literal-block text pass |
| non-default impact title tag | `capture.md` Step 1 + REQ Title Convention; publication validates an emitted field/title pair without guessing the verdict | default has no required tag; a present non-default impact cannot disagree with its title mirror |
| `created_at` / blocked timestamp shape | schema/publication canonical timestamp validation; timestamp-generation procedure remains prose | date-only, offset, local-with-Z, fractional, or malformed authored stamps refuse; UTC whole-second passes |
| domain, TDD, maintenance, impact, priority, effort values/defaults | `schemanormalization` | every canonical value/default/alias remains table-tested; capture rejects aliases/noncanonical authored values |
| key aliases | `schemanormalization` key-alias registry used by `requestmodel`; aliases remain read-only | all five alias families and canonical precedence are table-tested; capture rejects alias keys |
| `prime_files`, `required_lessons`, `depends_on`, `related`, `write_set`, `requests` list shape | `capture-files` structure validation; path relevance and budgets remain prose | scalar where a list is required refuses; empty and populated canonical lists pass |
| REQ/UR identity, linkage, membership, queue path, reservation | existing `capture-files` validation | retain existing public tests and add exact-membership/phantom-member refusal |
| `status: blocked` metadata | `capture-files` structure validation; deciding that the condition is blocking remains prose | blocked without `blocked_by`/`blocked_at` refuses; ordinary pending payload does not require them |
| TDD/maintenance/priority/assignment/external-condition applicability | `capture.md` Step 1 | no new mechanical inference; concise pointer outside each fence |
| RED/GREEN, complex sections, Open Questions, Prior Implementation | record examples plus existing capture prose | preserve the named shapes; do not add a brittle sentence grep |

## Likely File Scope

Primary files:

- `skills/do-work/actions/capture-reference.md` — replace the four fences, remove the duplicated Schema Aliases value table, and add one concise schema/publication pointer plus one judgment pointer.
- `skills/do-work/actions/capture.md` — update the companion/Step 5 wording so examples are copyable shapes and Go is the mechanical acceptance authority; keep all judgment rules.
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` — centralize key aliases and expose canonical-write evidence without changing ordinary normalize-and-warn reads.
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go` — table-test all canonical values/defaults/value aliases/key aliases and canonical-write rejection.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` — consume the schema-owned key aliases rather than restating them, and project the canonical/alias-backed fields needed by publication validation.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` — prove canonical precedence and all centralized alias projections.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go` — validate new UR/REQ payload shape and canonical authoring before building any mutation.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go` — public RED/GREEN cases for malformed/canonical records and the four example shapes.

Conditional fixture-only scope, to use only if stricter shared `BuildCapturePlan` validation makes its structured override fixture invalid:

- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` — update the existing override-capture fixture to the canonical minimum; no answer behavior change.

Explicit exclusions:

- `_dev/tests/contract-regressions.sh` and `_dev/tests/contracts/core-checks.sh` — there is no active template-comment predicate to delete.
- `skills/do-work/actions/work-reference.md` and `skills/do-work/docs/capture-guide.md` — they are separate schema/reference surfaces and belong to the later reference sweep; edit only if exploration proves a direct broken link caused by this REQ, then record the scope expansion before building.
- publication manifest types — do not add a `record_kind`; use fields already present in the payload to validate common capture shape, leaving Simple/Complex selection as judgment.

## Ordered Implementation Plan

1. **Write public RED tests for the missing capture-authoring contract.** Add a reusable canonical UR/REQ fixture in `capture_files_test.go`, then show that `BuildCapturePlan` currently accepts at least: a read-only key alias, a value alias/noncanonical enum, malformed required list/scalar metadata, a noncanonical authored timestamp/title scalar, and a UR with phantom `requests` membership. Add positive cases for the four minimal record shapes and preserve the existing raw-containment, collision, topology, and rollback cases. Run the publication package and record the intended failures before implementation.

2. **Make the schema layer the single field-rule source.** Move all five key-alias families into `schemanormalization`, add read-only access/selection APIs and a canonical-write check that distinguishes canonical values from accepted read aliases/default fallbacks, and table-test the entire registry. Update `requestmodel` to consume those APIs and expose exact source/evidence for every capture field publication needs. Keep ordinary read normalization byte-compatible and warning-compatible.

3. **Enforce canonical new-record shape at the publication boundary.** Before `BuildCapturePlan` appends mutations, validate UR/REQ required identity and field shapes, canonical authored keys/values, safe scalar form, canonical timestamps, blocked-field pairing, exact manifest-to-UR membership, and all applicable list fields. Return a stable typed `CAPTURE-UR-SCHEMA-INVALID` or `CAPTURE-REQ-SCHEMA-INVALID` refusal with record id/path and the concrete field reason. Preserve existing specific linkage/path/raw/fold refusal codes where those checks already own the failure. Do not infer impact, priority, effort, TDD, maintenance, assignment, dependency, or record complexity.

4. **Reduce the action reference without weakening judgments.** Replace the Simple REQ, Complex REQ, UR input, and Addendum REQ fences with coherent uncommented examples totaling at most 68 lines. Replace the Schema Aliases table with a direct pointer to the schema registry and publication validator. In `capture.md`, make Step 5 say the examples define record shape, the Go command checks mechanics, and Step 1/named reference contracts own judgments. Retain Fold-First, Required Lessons Budget, dependency/write-set, outside-text, open-question, and archived-addendum judgments. Re-run the repository search proving there is no active comment predicate to delete.

5. **GREEN, compatibility, and full verification.** Run gofmt, focused package tests, the race package check, the module suite, action/contract checks, then the canonical maintainer gate. Re-read every changed fence as a copyable artifact; verify the aggregate line target, all four named examples, no per-field comments, no new prose predicate, no alias/default drift, and no unexpected source file outside the declared scope.

## RED/GREEN Cases

- **RED — alias key:** a new REQ payload using `dependencies:` (and table variants for the other schema key aliases) currently publishes; GREEN refuses it while ordinary `requestmodel` reads still resolve the alias.
- **RED — alias/invalid value:** `domain: back_end`, `tdd: yes`, or an invalid normalized field currently reaches the plan; GREEN requires a canonical authored value and includes the field in the refusal.
- **RED — unsafe/canonical metadata:** an unquoted user-derived title or non-UTC/non-whole-second `created_at` currently publishes; GREEN refuses it, while single-quoted apostrophes/literal-block text and `YYYY-MM-DDTHH:MM:SSZ` pass.
- **RED — record shape:** a scalar in a required list field or a blocked record missing its paired metadata currently publishes; GREEN refuses before any mutation is planned.
- **RED — membership:** a UR `requests` list may contain an ID absent from the manifest; GREEN requires exact set equality, while an empty list remains valid for a fold-only capture.
- **GREEN — examples:** manifests built from the four shortened examples pass the public capture validator, including Simple, Complex, fold-only/ordinary UR, and Addendum linkage shapes.
- **Regression:** canonical-key precedence and every legacy read alias remain accepted by non-capture readers; existing capture path/reservation/raw-input/asset/fold safeguards retain their specific findings.

## Test Commands

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

Action and shipped-contract checks:

```bash
bash _dev/tests/contract-regressions.sh
bash _dev/tests/staged-skills-contract.sh
```

Canonical repository gate (run directly and unpiped):

```bash
bash _dev/tests/maintainer-verify.sh
```

Template-size and residue audit (verification only, not a permanent sentence predicate): count the four fenced example spans before/after, require a combined post-edit count of at most 68, and search the edited fences for inline `#`/arrow explanations while allowing Markdown headings in the record bodies.

## Plan Validation

- Requirement coverage: every Detailed Requirement maps to tasks 1–4; task 5 supplies the requested test and full-suite evidence.
- No orphan tasks: schema centralization is necessary to make the alias table deletable; requestmodel changes remove the current duplicate alias source; publication validation supplies the missing machine enforcement; action edits perform the requested deletion.
- Scope sanity: five ordered tasks, with one conditional fixture-only file and explicit exclusions.
- Consumer field contract: publication validates exact record identity (`id`, canonical path), provenance (`user_request`, manifest linkage/reservation), state (`status` plus blocked pairing), and outcome (typed refusal code/reason/path or a mutation plan). Ordinary readers retain source key, raw value, resolved value, recognition/default state, and warning evidence from schema/requestmodel.

*Generated by Plan agent*
