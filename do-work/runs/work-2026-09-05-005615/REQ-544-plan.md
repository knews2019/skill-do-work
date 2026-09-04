# REQ-544 plan and exploration

Proceed with a narrow publication matching fix and cleanup caller coverage. The two stakeholder gates remain unanchored; cleanup already delegates to the shared structural remover. No schema change, production cleanup edit or finalization edit is justified by the observed instances.

Read-only preparation, no tests run. Inspected the queued REQ, REQ-544-prep.md, the CLI prime and whole lesson satellite, shared/general/testing/coding/backend/security/communication/anti-slop crews, and the current source. Latest observed main HEAD: `eabc29842dc537eb2cbc43adcfd2ae294bfbc92b`, which contains REQ-504's checkpoint bounds repair. Re-resolve after canonical claim; preparation is not gate evidence.

## Exact proposed scope

- `skills/do-work/tools/do-work-cli/internal/publication/answer.go`
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go`
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go`

The source edit stays package-local in publication. Command-runtime seam tests can live in the declared test files and reuse the packages' handlers and fixture helpers. Remove the captured `cleanup_plan.go` production path from the prospective Scope with a recorded decision: its old containment loop no longer exists. If an actual uncovered failure requires a further production path, report the exact failing case before expanding Scope.

## Writer positions are the contract

The role of these bytes is an action-authored attestation, not natural-language sentiment analysis. Require the marker at the writer's structural position within the appropriate evidence block. Matching a word somewhere in a caller's narrative cannot establish that role. The parser should not add a free-standing list of synonyms or negative words; derive the accepted positions from the actual writer forms below.

| Evidence | Actual writer or existing compatibility form | Required anchor |
|---|---|---|
| Blocked resolution | `stakeholder-answers.md:66`: `- [<date>] blocked on "answers from <name>" — resolved: all questions answered` | A real top-level Blocked history entry, with the resolution disposition in its trailing writer-supplied field; a marker inside the quoted name or before the separator is not the disposition. Ambiguous competing separators must refuse. |
| Existing Blocked fixtures | `answer_test.go`: `- Resolved 2026-09-01 after stakeholder answer.` and `- Resolved after stakeholder override.` | The explicit `Resolved` marker begins the history entry, with a whole marker boundary; context following it does not license finding the marker later in an arbitrary sentence. Preserve this accepted legacy form rather than rewriting unrelated positive tests. |
| Implementation disposition | `stakeholder-answers.md:66` and clarify's confirmed-choice path: `**No changes needed.**` followed by explanatory prose | The exact no-change marker begins an active Implementation paragraph; its sentence/format boundary must be complete. This shipped form currently lacks the substring `no code`, so a genuine writer-form positive control is essential. |
| Existing Implementation fixtures | `No code changes.` | Preserve this explicit legacy no-code marker at the same paragraph-start role; `no code review yet` is not that marker. |
| Report linkage | `work-reference.md` Stakeholder template: `- [<date>] ai-reports/<bundle>/index.html — <count> open question at generation`; existing fixture omits the date | The history entry's path field equals the normalized manifest report path in full, followed by the field delimiter or end of line. A path appearing in commentary, another path's prefix/suffix, a quoted example or another section does not qualify. Preserve dated and existing undated forms. Do not tokenize valid paths solely on spaces. |

Evidence payloads can retain their currently supported optional exact heading. Bind interpretation to that heading's actual section, not `## Blocked extra`, an unrelated section or a substring of prose. Fenced/quoted examples are not active records. Keep original LF/CRLF bytes for publication; matching must not rewrite unrelated history. The parser should validate the bytes it will actually append. Avoid turning this narrow fix into a repository-wide Markdown parser replacement.

## Task 1 — Prove the two publication forgeries and genuine controls before changing code

In `answer_test.go`, add focused table-driven tests for terminal and partial stakeholder evidence. Exercise the public command-runtime/handler seam with real manifests and payload files for at least the captured terminal forgery and the report-linkage forgery; use existing planner tests for the remaining grammar rows when they improve failure localization.

- Terminal RED: submit the recorded `still not resolved` / `no code review yet` pair. Test each independently with the other side genuinely valid, so one rejection cannot hide the other missing check. Require existing `ANSWER-STAKEHOLDER-EVIDENCE-INVALID`, unchanged queue bytes and blocked metadata, no `completed_at`, and no archive/override/report publication.
- Add displaced/quoted/fenced/wrong-section and false marker-boundary cases, plus a `resolved:` token in a stakeholder name. Cover multiple separator ambiguity. Keep the canonical resolution field control and the legacy prefix control separate.
- Genuine terminal controls: the shipped dated resolution plus `**No changes needed.**` pair, and the existing `- Resolved …` plus `No code changes.` pair. Require completed status, correct timestamp, blocked-field removal and actual archive bytes.
- Partial RED: the exact report path appears only in commentary, in a longer path, or in a different section. Require existing `ANSWER-STAKEHOLDER-REPORT-EVIDENCE-INVALID`, unchanged original `blocked_by`, no fresh report publication and no request mutation. Genuine dated and undated history fields must update the linkage and publish the immutable report together. Keep `ANSWER-STAKEHOLDER-REPORT-LINKAGE-INVALID` for the separate manifest `blocked_by` mismatch.

Run and record RED before implementation, including the canonical writer-form control's current refusal. Existing happy paths must remain green; do not change their evidence just to match a new recognizer.

## Task 2 — Replace only the matching gates

In `answer.go`, replace whole-payload containment with small package-local matchers for the exact evidence roles above. Keep the typed manifest, question-resolution authority, refusal codes, terminal status policy and publication transaction unchanged. Determine completion from the existing question verdict first, then validate the required attestation in the exact block/field position; prose found outside that position cannot supply authority.

The existing field roles and canonical/legacy marker spellings are enough for these listed forgeries; no new stored schema is currently necessary. Do not weaken matching merely to accept arbitrary previous caller prose. If preserving an actually documented format makes the position ambiguous, record the concrete collision and seek a scope decision rather than quietly inventing a new marker language. Explain the distinction between canonical and compatibility forms in the implementation decision and test names.

## Task 3 — Pin the already-repaired cleanup caller and verify the bounded exclusion

`cleanup_plan.go:ownedCheckpointRemoval` already calls `requeststate.RemoveOwnedCheckpointClaim`. The latter uses `CheckpointClaimBounds`, a request ID anchored at entry start and a writer anchored at the end. REQ-504 repaired canonical-versus-headingless range selection. The existing cleanup test proves own/foreign enriched-entry handling and a prose heading token, but does not explicitly pin the originally captured two-token forgery.

Extend `cleanup_plan_test.go` with a caller control containing both `- REQ-NNN:` and the exact current writer inside a narrative line. It must survive cleanup byte-for-byte. Pair it with a genuine own entry that is removed including its continuation, while foreign and unlabeled records survive. Exercise real cleanup application through its handler/runtime on a committed fixture, not just an isolated matcher; confirm the completed request moves as expected and no checkpoint replacement is planned when the only matching-looking line is forged. Include one headingless legacy control to cover the REQ-504 seam. These should already pass on the pre-build source; record that as existing protection, not newly fixed RED.

The finalization exclusion was read directly: `associateReleaseMetadata` calls `singleInsertion(before.Bytes, after.Bytes)` before request/version containment. `singleInsertion` requires the unchanged prefix before the first existing heading and the complete original suffix, returning only the inserted bytes. Association also requires the configured release members and exactly one candidate. Keep this bounded request-token check and its tests untouched.

After GREEN, run `go test -count=1 ./internal/publication ./internal/cleanup` and the corresponding focused race tests once if the normal workflow calls for them. The orchestrator owns canonical preflight, post-merge gate and heavy-lane planning. Sweep the named action writers and other uses of the removed primitives for attribution, but keep unrelated `findQuestionLine`, report-only OK suppression, generic section append behavior and finalization ownership outside this change unless a listed forgery demonstrably depends on them.

## Expected handback

Three changed paths, exact RED/GREEN output, existing-green cleanup evidence, the genuine writer-form controls, preserved refusal codes, and an explicit note that the bounded finalization case was inspected and left unchanged. No arbitrary token whitelist, new schema field, synthetic cleanup source edit or unexplained positive-fixture rewrite.
