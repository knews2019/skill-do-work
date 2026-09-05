---
id: REQ-544
title: '[impact-critical] Anchor every lifecycle gate that reads caller-authored text'
status: completed
priority: now
created_at: 2026-09-03T09:45:00Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-528]
sweep: true
sweep_key: answer-line-marker-position-spoofing
review_generated: true
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/answer.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer_test.go
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go
claimed_at: 2026-09-05T00:33:24Z
route: A
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-05T00:50:53Z
  basis:
    - Route A
    - 4-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
completed_at: 2026-09-05T08:03:31Z
commit: 050a54dd
release_at: 2026-09-05T08:03:31Z
---

# Anchor Every Lifecycle Gate That Reads Caller-Authored Text

## What

REQ-528 fixed one instance of a shape: a lifecycle decision made by searching caller-authored text for a token, anywhere in that text. The shape-grep it ran while fixing that instance found three more, one of which drives a second terminal status in the same file. This sweep closes them all against the same invariant: **evidence a caller writes freely cannot be the anchor for a lifecycle write.**

## Instances

- **`internal/publication/answer.go:268` — a second terminal status, same file.** The stakeholder branch gates `status: completed`, `completed_at`, deletion of `blocked_by`/`blocked_at`, and the archive move on `bytes.Contains(ToLower(blockedHistory), "resolved")` **and** `bytes.Contains(ToLower(implementation), "no code")`. Both payloads are caller-authored prose and both tokens match anywhere, so `still not resolved` satisfies the first and `no code review yet` satisfies the second. The refusal message ("terminal evidence must carry resolved Blocked history and an Implementation no-code marker") names markers, but the code tests for substrings.
- **`internal/publication/answer.go:291` — non-terminal but still a lifecycle write.** The `blocked_by` linkage is gated on `bytes.Contains(reportsHistory, reportPath)` over caller-authored history text, unanchored to any line or position.
- **`internal/cleanup/cleanup_plan.go:236` — adjacent, and destructive.** A `do-work/CHECKPOINT.md` line is deleted when it merely *contains* `- <REQ-ID>:` plus the writer token, with no position anchor. Narrower than the others because both tokens must share one line, but it is the same unanchored-line-selection shape driving a destructive edit.
- **Bounded, listed for completeness rather than repair:** `internal/finalization/finalization_discovery.go:593` admits an unjournaled changelog tail into replay on `bytes.Contains(inserted, requestID)`. `singleInsertion` bounds the searched bytes to one verified diff hunk, so the containment is already anchored by its caller. Confirm that rather than assume it, then leave it.
- **Checked and deliberately not included:** `answer.go:406` (`findQuestionLine`) and `internal/corehelpers/checks.go:232` are unanchored but fail closed — the first requires exactly one matching line and refuses otherwise, the second only suppresses an `OK:` line in a report it composed itself.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the REQ, the archived REQ-528 in full, the CLI prime, the whole lessons satellite, the coding/backend/testing crew members, the orchestrator's `REQ-544-plan.md` and `REQ-544-prep.md`, and the four target files plus `stakeholder-answers.md` Step 5 and the `work-reference.md` stakeholder template. Approach recorded before any edit: read each marker only at a writer-contributed position, bound each scan to its own section, and check `cleanup_plan.go` before assuming the captured defect was still there — it was not.
- [x] **[APPLY]:** Three files, all inside the declared write set; `cleanup_plan.go` deliberately untouched (D-05). RED captured first for both live instances — 19 rows redden under a full revert — then the matchers landed, then two rows were added to pin the two guards the first neuter matrix showed unpinned.
- [x] **[UNIFY]:** `git diff --stat` → 3 files, +620/-4; read the whole `answer.go` diff line by line and both test diffs. No debug artifacts, no `TODO`, no scratch text, no build output; the temporary neutered copies were restored and verified byte-identical with `diff -q` after every experiment. `gofmt -l .` silent, `go vet ./...` clean, `go test -count=1 ./internal/publication ./internal/cleanup` green, `go test -count=1 ./...` green except the pre-existing `internal/heavyverification` failure reproduced at clean HEAD.

## Finding Provenance

- Findings **F-01 through F-04** from REQ-528's implementation shape-grep, which was asked for exactly because the fixed bug had a shape rather than being a one-off. REQ-528 fixed the reported instance only; this sweep owns the rest.
- The `impact-critical` token is carried from the first instance: it decides a terminal status and an archive move on evidence its own caller supplies.

## Detailed Requirements

- Every lifecycle or destructive decision that reads caller-authored text must anchor the token to a position the reader can attribute — a line start, a field boundary, or a separator the writer itself contributes. Containment anywhere in free text is not evidence.
- Where a token is genuinely a marker, test it as a marker: anchored, and at the position whatever wrote it places it. Where the refusal message already says "marker", the code must match that claim.
- An unattributable read must fail toward the non-destructive, non-terminal outcome, as REQ-528 established.
- Preserve every existing refusal code and typed result; these are matching fixes, not policy changes.
- Cover each instance with a test that forges the token in caller text and asserts the lifecycle write does **not** happen, plus a control proving the genuine path still does.

## Constraints

- Do not introduce a new schema field or change any stored format unless a matching fix genuinely cannot anchor the token; say so explicitly if you reach that conclusion.
- `internal/cleanup` and `internal/publication` are separate packages; honor the prime's **Package direction** rule rather than reaching for a shared helper that would violate it.
- Do not weaken `singleInsertion`'s bounding in `finalization_discovery.go`.

## Dependencies

Depends on REQ-528, which establishes the anchoring invariant and the fail-toward-non-terminal rule this applies.

## Red-Green Proof

**RED prompt/case:** Submit a terminal stakeholder disposition whose Blocked-history payload reads `still not resolved` and whose Implementation payload reads `no code review yet`, then inspect the resulting `status`.
**Why RED now:** both gates use `bytes.Contains` over caller-authored prose, so both negations satisfy them and the REQ reaches `completed` with `completed_at` set and the archive move planned.
**GREEN when:** that submission is refused with the existing `ANSWER-STAKEHOLDER-EVIDENCE-INVALID` code; a genuine terminal evidence pair still reaches `completed`; the `blocked_by` linkage and the CHECKPOINT line deletion are likewise anchored with forge-and-control tests; and `singleInsertion`'s bounding is confirmed intact.

---
*Source: REQ-528 implementation shape-grep, findings F-01 through F-04.*

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names all four instances with file and line, states the invariant they violate, and lists its own four-file write set. It also records what was deliberately excluded and why. Nothing needs discovering.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go` (modified)

**What was done:** The three whole-payload `bytes.Contains` gates in the stakeholder publication path were replaced with matchers that read each marker only at a position a writer contributes, and the region a payload is read from now ends at the first Markdown block construct. The checkpoint-deletion instance was found already anchored by earlier work, so it was locked in with tests instead of edited.

Two gates were live defects. The terminal stakeholder gate decides status completed, the completed_at stamp, deletion of blocked_by and blocked_at, and the archive move; before this change a Blocked-history payload reading "still not resolved" and an Implementation payload reading "no code review yet" both satisfied it. The report-linkage gate publishes blocked_by on containment of the report path anywhere in caller history text. Both now read the marker at a history entry's list bullet, at the field after the bracketed date that bullet may carry, or at the field after the entry's single separator, and never anywhere else.

The third instance, the CHECKPOINT line deletion in the cleanup package, was already repaired in the tree by REQ-502, which replaced the containment loop with an anchored shared remover. No production edit was written there. The fourth item, the changelog-replay containment in finalization discovery, was confirmed rather than changed: singleInsertion returns only the bytes the working-tree change inserted, so the caller's bound holds and nothing was weakened.

A reviewer pass then showed the first fix had written the block exclusion as a list of two spellings instead of a condition. The continuation replaced it: a payload's record region ends at the first line that opens a Markdown block construct, meaning an ATX heading or a run of three or more of one ASCII punctuation mark, and a line indented four or more columns records nothing. The punctuation class is reused from the predicate already in that file, so a fence character the code has never seen still ends the region. The same pass taught the report-path reader three skins of the same field, including the Markdown-link form this repository's own archive uses.

The fourth file in the declared write set, cleanup_plan.go, was left byte-identical on purpose. See D-05.

Manifest check against `git diff 5ff8e9c0..050a54dd --stat`: the range covers 49 files, because it carries several other requests as well. The three files above are the only REQ-544 files in it, and their line counts match the range exactly at 341 insertions for answer.go, 430 for answer_test.go and 139 for cleanup_plan_test.go. cleanup_plan.go does not appear in the range, which agrees with the hand-back.

**Implementation range:** `5ff8e9c0..050a54dd`. Builder commit `bd089254` with continuation `6c9ca166`; the range carries them as `61aec3dd` and `050a54dd` with identical stats.

## Decisions

- **D-01 — Markers are read at three writer-contributed positions, never by containment:** a history entry's list bullet, the field after the bracketed date that bullet may carry, and the field after the entry's single separator. Those are the only bytes in a caller-authored payload whose author a reader can name. Everything else is prose, and prose that happened to contain the marker was the whole defect.
- **D-02 — An entry carrying more than one separator has no attributable trailing field and fails the check:** it is REQ-528's rule applied to the same shape one section down, and it is the only rule that resists last-separator anchoring, since the neuter matrix shows the entry with a subject containing an em dash is the sole row that reddens when the count guard is swapped for a last-index search. The risk is that a genuine entry whose subject holds an em dash is refused, which is the required direction: the refusal happens before any mutation, the request stays byte-identical in the queue, and the caller reworks one line. Accepting it would archive a request on ambiguous evidence.
- **D-03 — The Implementation gate accepts two marker spellings, the one place the command was widened:** "no code" is a seven-character English fragment, not a marker, so position anchoring alone could not separate it from the request's own RED case "no code review yet", where the fragment sits at a paragraph start. Testing a real marker literal is the only fix that closes it, and the literal the shipping action prescribes contains no "no code" at all. This deviates from the request's constraint that these are matching fixes and not policy changes, and it is recorded rather than smuggled in: refusing the documented form would have cemented an unreachable terminal path, which the pre-fix binary demonstrates by refusing the writer form the shipped prose prescribes.
- **D-04 — The spelling list was kept to the two forms that exist, with a code comment saying why a stale list is safe here:** the project's closed-enumeration lesson applies to lists standing in for a condition, and this is not that. The marker is a fixed spelling by contract, and a missing spelling causes a visible pre-mutation refusal of a genuine disposition, never acceptance of a forgery. The comment states the direction the list may drift.
- **D-05 — The fourth declared write-set file was left byte-identical, against the declared write set:** the containment loop the request captured no longer exists, because REQ-502 replaced it with an anchored shared remover. Writing a production edit to fill a four-path list would have been machinery, not a fix. The instance is covered instead by a caller-level lock-in whose three subtests all redden under the captured defect.
- **D-06 — Fenced blocks and later section headings are excluded from every evidence scan:** a fenced example of a resolution entry describes a record rather than being one, and a line under a later heading is not Blocked history. Both cost about four lines and both are pinned by the neuter matrix. The reviewer later showed this was written as a spelling list, and D-10 replaced it with a condition.
- **D-07 — The terminal refusal reason names which half failed and where its marker belongs:** REQ-528 found that a caller cannot correct what the refusal does not name, and these two payloads fail for different reasons often enough that one shared sentence sends the caller to the wrong file. It is pinned per row rather than left as decoration, because replacing the evidence with a generic sentence reddens all forgery rows for that gate.
- **D-08 — No schema field, no stored-format change, no new refusal code:** the request asked for an explicit statement if a matching fix could not anchor a token, and that case did not arise; all three places could be anchored. The only stored-format consequence is D-03, which widens what an existing field accepts rather than changing what is written.
- **D-09 — Two test rows were added after the first neuter matrix, not before:** the first matrix showed the word-boundary guard and the list-bullet guard reddening nothing, which is defense nothing earned. Rather than delete them, the forgery each one actually stops was found and pinned, so both now redden their own guard and nothing else.
- **D-10 — The section boundary is a condition over CommonMark's own ingredients, and the punctuation class is reused rather than re-enumerated:** it closes the tilde-fence, indented-code and deeper-heading bypasses at once instead of one spelling at a time, and the run test asks the punctuation predicate already in the file, so there is one enumeration of "punctuation that opens Markdown structure" for both readers. A future fence character is punctuation, so it is already handled. The risk is that a genuine record placed after a thematic break or a subheading inside the same payload is refused, which is fail-closed, pre-mutation, and one line for the caller to move.
- **D-11 — The region ends at a block construct instead of toggling a fence pair:** a toggle must find its close, and an unclosed or mismatched fence then makes everything after it readable again, which is the one failure this treatment cannot have. Ending is strictly safer and needs no pairing state.
- **D-12 — One owner for the indentation condition, measured before any trim:** the same function decides it for the history readers and the Implementation reader together, so they cannot disagree. The measurement happens before trimming because trimming first destroys the whitespace that constitutes the structure, which is REQ-460's recorded trap and the direct cause of the indented-code bypasses.
- **D-13 — A path field is read in three skins and the whole field is compared:** a link-wrapped path, a backticked path and a bare path are the same evidence, with different terminators but the same field. Comparing the whole field rather than prefix-testing keeps a neighbouring bundle whose path merely starts with this one refused in all three skins at once, and the neuter matrix shows that single change reddening all three longer-path rows.
- **D-14 — The rephrasing "No changes were needed." stays refused:** every writer and reader of this marker in the repository agrees on the two accepted spellings, including the doctor scan that suppresses its hollow-completion finding on the documented wording. Accepting the rephrasing here and nowhere else would let a request complete through the answer command and then be reported hollow by doctor, which is the split the marker exists to prevent. It also has no boundary, since once one rephrasing is accepted the neighbouring ones are a word away and a negation sits at the same position. The risk is that a caller who writes the rephrasing is refused, and the refusal now prints both accepted spellings verbatim, so the correction is mechanical.
- **D-15 — Marker spellings stay a list; positions are conditions:** these are opposite problems. A stale position list accepts a forgery, which is why the reviewer's first finding was a real defect. A stale spelling list refuses a genuine disposition, which is visible, pre-mutation and correctable, and the spellings are owned by the writers rather than by this reader.
- **D-16 — The unpinnable escape guard was deleted rather than kept just in case:** the report path has already passed the containment check, so it can never equal a path that climbs out or starts at the root, and the escape branch could not change a verdict. A branch no test can redden is the shape REQ-528 warned about, so it is gone and a comment records that the containment upstream is the real guard. The escaping-link attack row stays and still refuses.

## Qualification

Passed the request-bound advance qualify gate for the cumulative range `5ff8e9c0..050a54dd`. Three of the four declared files changed; `cleanup_plan.go` was left byte-identical because REQ-502 already anchored that instance, verified independently by restoring the old containment loop and watching all three new subtests redden. Independent review ran 63 attack payloads through the built binaries and confirmed the whole-command differential in both directions. The P-A-U boxes were reconciled from the builder hand-back.
## Testing

**Red-green validation:** Before the fix, the terminal stakeholder gate accepted every forged payload. Nine rows of `TestStakeholderTerminalEvidenceRefusesMarkersForgedInCallerProse` failed with `forged terminal evidence accepted: refusal=(*publication.Refusal)(nil)`, covering narrative negation, negation inside a real history entry, the marker hidden in the quoted stakeholder name, the marker only inside a fenced example, the marker only under a later section heading, an entry carrying competing separators, and three implementation-marker forms including `implementation marker buried mid-sentence`. In the same run two rows of the control table `TestStakeholderTerminalEvidenceKeepsGenuineWriterFormsTerminal` failed the other way with `genuine terminal evidence refused: &publication.Refusal{Code:"ANSWER-STAKEHOLDER-EVIDENCE-INVALID", Reason:"terminal evidence must carry resolved Blocked history and an Implementation no-code marker", ...}`, one of them being `the_form_actions/stakeholder-answers.md_prescribes`. The report gate failed the same way: five rows of `TestStakeholderReportLinkageRefusesReportPathForgedInCallerProse` reported `forged report linkage accepted: refusal=(*publication.Refusal)(nil)`, including a path mentioned only in commentary and a longer path that merely starts with the real one.

After the fix, `go test -count=1 -run TestStakeholder ./internal/publication` returned `ok ... 0.024s` with all 22 rows passing. The full-revert differential, with both gates put back to whole-payload containment and nothing else changed, reddens 19 rows: every forgery row plus both canonical writer-form controls, and nothing else in the package.

The reviewer's continuation started from a second RED of 24 failing subtests across the four tables, including seven control rows the first pass wrongly refused. The four tables are now `TestStakeholderTerminalEvidenceRefusesMarkersForgedInCallerProse` (21 rows), `TestStakeholderTerminalEvidenceKeepsGenuineWriterFormsTerminal` (7 rows), `TestStakeholderReportLinkageRefusesReportPathForgedInCallerProse` (12 rows) and `TestStakeholderReportLinkageKeepsGenuineHistoryEntriesLinked` (6 rows), and all pass.

The checkpoint instance had no honest RED, because the defect the request captured had already been repaired by REQ-502 before this work started. `TestCheckpointRemovalIgnoresRequestAndWriterTokensForgedInProse` passed the first time it ran, at `ok ... 0.480s` against unmodified production code. It was proved non-decorative by neutering: with the pre-REQ-502 containment loop temporarily restored, all three subtests fail with `forged prose line was deleted`, `forged prose line was deleted from a headingless checkpoint` and `forged prose planned a checkpoint replacement`. The production file was restored afterwards and verified byte-identical.

Beyond the harness, 46 attack rows were run through the built binary, each a real answer invocation against its own throwaway git repository with committed fixtures. After the continuation, 33 forgeries refuse and 13 controls succeed, with 0 rows wrong. Against the first-pass binary the same table produced 17 wrong rows: 10 bypasses that completed and archived or linked, and 7 false refusals including the archived Markdown-link form, the backticked path and the asterisk bullet. No forgery the first pass refused regressed.

**Controls preserved:** The two stakeholder control tables above were re-run at every step and protect the genuine terminal disposition and the genuine report linkage, which is the path a real caller uses to close a stakeholder request. A 12-mutation neuter matrix was run one mutation at a time, with the production file restored and verified byte-identical after each, and every guard reddens at least one row of its own: the punctuation-run condition 3 rows, the ATX heading branch 2, the indented-code threshold 4, the bullet marker set 1, the bullet-follower check 1, the link-versus-date disambiguation 1, the Markdown-link path field 3, the backticked path field 1, whole-field equality 3, relative resolution 1, the Implementation list-item skin 1, and the readable alternatives in the refusal reason 7. No mutation reddened a row it should not. The full `internal/publication` suite, which carries REQ-528's earlier gates, was re-run green, and every existing refusal code was confirmed unchanged, including `ANSWER-STAKEHOLDER-EVIDENCE-INVALID` and `ANSWER-STAKEHOLDER-REPORT-EVIDENCE-INVALID` whose reasons were extended but whose codes and trigger points are the same. The typed results `PublicationPlan`, `Refusal`, `PlannedMutation`, `StakeholderTerminalEvidence` and `StakeholderReportEvidence` gained, lost and renamed no field.

**Module verification:** All commands were run from the CLI module directory on Go 1.26.1. First pass: `go test -count=1 ./internal/publication ./internal/cleanup` returned `ok ... internal/publication 22.315s` and `ok ... internal/cleanup 3.929s`; `go test -count=1 ./...` returned every package `ok` except `internal/heavyverification`; `gofmt -l .` printed nothing; `go vet ./...` printed nothing and exited 0. Continuation re-run: `go test -count=1 ./internal/publication ./internal/cleanup` returned `ok ... internal/publication 18.002s` and `ok ... internal/cleanup 2.631s`; `go test -count=1 ./...` returned 28 packages `ok` with the same one failing; `gofmt -l .` printed nothing; `go vet ./...` printed nothing and exited 0. `_dev/tests/maintainer-verify.sh` was not run, per instruction.

The one failing package was diagnosed as not belonging to this change. `internal/heavyverification` fails on `TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes` (5 subtests, "default runtime must have a determinable fingerprint") and `TestShippedGitIsolationPreservesGenericLaneInheritance` ("shipped runtime probe did not isolate host Git configuration"). This was verified rather than assumed: the three changed files were moved aside, restored to clean HEAD, and the package was run again, producing identical failures with the same tests and the same messages. That package imports neither `publication` nor `cleanup`.

## Discovered Tasks

- The Blocked history vocabulary is split between two flows. `skills/do-work/actions/stakeholder-answers.md:66` prescribes the resolution wording that now works, while the two real archived entries at `do-work/archive/UR-031/REQ-144-activate-four-skill-distribution.md:72` and `do-work/archive/UR-031/REQ-146-remove-modular-migration-shims.md:83` close with "cleared by user via clarify", which carries no resolution marker. Those came from clarify's blocked-condition confirm rather than the terminal stakeholder path, so nothing is broken today, but the two vocabularies should either converge in prose or the reader should gain the second spelling. → queue as follow-up
- No test pins the marker constants in `skills/do-work/tools/do-work-cli/internal/publication/answer.go` against the literals the prose writes at `skills/do-work/actions/stakeholder-answers.md:66` and `skills/do-work/actions/clarify.md:203`. Three readers now depend on the same spellings, counting `internal/doctor/doctor_scan.go:278`, so a prose/CLI split should fail in CI rather than at review. This is REQ-528's remaining suggested test. → queue as follow-up
- `appendSectionEvidence` in `skills/do-work/tools/do-work-cli/internal/publication/answer.go:~740` trims the heading with a prefix trim, so a payload opening "## Blocked extra" has "## Blocked" shaved off its front and " extra" published into the Blocked section. The new matcher refuses such a payload before that point, so the terminal path is safe, but the partial path's Reports append and any future caller are not. → queue as follow-up
- `internal/heavyverification` is red at clean HEAD on `TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes` (5 subtests) and `TestShippedGitIsolationPreservesGenericLaneInheritance`, reproduced with none of this request's changes present. The messages are environment-shaped rather than logic-shaped, but it is a standing full-module failure someone owns. → queue as follow-up
- `questionSectionBytes` in `skills/do-work/tools/do-work-cli/internal/publication/answer.go:~700` still falls back to the whole body when a request has no "## Open Questions" heading. The stakeholder template at `skills/do-work/actions/work-reference.md:903` deliberately uses "## Questions", so every stakeholder request takes that fallback and its P-A-U checkboxes are judged as resolved questions. Pre-existing, recorded by REQ-528, unchanged here, and now depended on by two branches instead of one. → report only
- A negation sitting at the marker position still reads as the marker, so the payload "blocked on answers from Dana, resolved: nothing was resolved" and an Implementation note whose no-change marker is followed by "except for the seven files we changed" both complete and archive. This is inherent to reading a marker a human writes and needs a second independent witness rather than a better parser. → report only

## Review

**Overall: 92%**
**Acceptance: Pass.** The reviewer ran 63 attack payloads through the built binaries rather than the test harness, and confirmed the whole-command differential: on the pre-fix binary the forged pair completes the request, stamps `completed_at`, deletes `blocked_by` and archives it; on the fixed binary it refuses and leaves the queue file intact.

It also confirmed the pre-existing bug this work exposed. The form `stakeholder-answers.md` Step 5 prescribes was REFUSED by the shipped binary, so a stakeholder request could never be closed by following the shipped prose. A full revert reddens 19 rows including two of the three genuine controls, which is the direct proof.

The review found the fence exclusion was a list of spellings rather than a condition, and demonstrated three live bypasses through the fixed binary — a `~~~` fence, a four-space indented block, and a `###` subheading each still completed and archived a request. It also found legitimate forms newly refused, including the Markdown-link Reports form this repository's own archive uses. Both were closed in the continuation.

One finding stays open by design: a negation sitting at the marker position still completes. That is inherent to reading a marker a human writes, REQ-528's own lesson names it, and it needs a second independent witness rather than a better parser.

## Lessons Learned

When a gate reads text a caller wrote, the position it reads is a condition and must be stated as one, while the literal it matches may stay a short list. The two fail in opposite directions: a stale position rule accepts a forgery silently, and a stale spelling rule refuses a genuine record visibly and before any write. This is why the fence exclusion written as "backtick fence or heading" was a real defect while the two accepted no-change spellings are safe, and it is the test to apply to any future evidence reader in this module.

Before fixing a defect a request captured with a file and a line, check that the defect is still there. One of the three instances here had already been repaired by other work, and the honest response was a lock-in test that reddens under the captured defect rather than a production edit written to fill out a declared write set.

## Orientation

Every lifecycle write in the stakeholder publication path now reads its evidence at a position a writer contributed, inside a region that ends at the first Markdown block construct, so prose that merely contains a marker or a report path can no longer complete, archive or link a request. The documented terminal form also works for the first time, since the shipped binary previously refused the exact wording the shipped prose prescribes.
