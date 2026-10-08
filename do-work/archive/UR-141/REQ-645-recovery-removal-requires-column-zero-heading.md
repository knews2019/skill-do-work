---
id: REQ-645
title: '[impact-critical] recover --take-over keeps a 1-3 space indented generated-name heading; only a column-0 section is removable'
status: completed
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-08T14:12:47Z
created_at: 2026-10-08T14:08:54Z
user_request: UR-141
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-critical
effort_estimate: effort-mechanical
related: [REQ-646]
batch: review-residuals-2026-10-08
write_set: [skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go, skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go, skills/do-work/tools/do-work-cli/internal/requeststate/recovery_markdown_test.go, skills/do-work/CHANGELOG.md]
route: A
claimed_at: 2026-10-08T14:12:16Z
dispatch_at: 2026-10-08T14:14:45Z
builder_handback_at: 2026-10-08T14:28:41Z
integration_at: 2026-10-08T14:29:06Z
review_at: 2026-10-08T14:36:04Z
kb_status: pending
commit: b5e358f1b25c0438e25ce90a52df85b397989d2d
heavy_verified_at: 2026-10-08T14:39:58Z
heavy_verified_revision: b5e358f1b25c0438e25ce90a52df85b397989d2d
completed_at: 2026-10-08T14:40:53Z
release_at: 2026-10-08T14:40:53Z
---
# Recover --take-over Keeps a 1-3 Space Indented Generated-Name Heading; Only a Column-0 Section Is Removable
## What
`stripGeneratedRecoverySections` (`skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go:981-999`) removes every visible section whose name is a generated lifecycle heading, whatever its indent. A user-authored bullet item with a two-space-indented `## Plan` sample under it is therefore deleted by `recover --take-over`. Make removal require a heading at column 0; leave boundary recognition (`VisibleSections`, 0-3 spaces) unchanged.
## Why
Review finding, P2: "When a requirement contains a bullet item with a two-space-indented `## Plan`, this now treats that nested heading as an independent request section. Consequently, `stripGeneratedRecoverySections` deletes the user-authored example during `recover --take-over`." Triage accepted the claim with a different remedy. This is a residual of REQ-635 F1 (recovery deleted user text under a 4-space indented `## Plan`): REQ-635 applied the CommonMark 0-3 space rule and deliberately kept 1-3 space headings as section boundaries (its requirement 1), so the 2-space list case survived. Data loss in a recovery path is `impact-critical`, as REQ-635 was.
## Verified Facts (from triage)
- `visible_sections.go:64-68`: a heading counts after 0-3 spaces (`blockIndentContent`). The REQ-635 test table asserts `" ## Plan"` and `"   ## Plan"` yield a `Plan` section.
- `state_apply.go:981-999`: `generatedRecoveryHeading` matches `"## " + section.Name`, so an indented `Plan` section is removed like a column-0 one.
- 0.305.59 (`96505655`) made 1-3 space indented headings boundaries on purpose: `TestRecoveryPreservesIndentedAndCommentedRequirements` (`recovery_markdown_test.go`) pins that `  ## Requirements` ends a generated `## Plan` so the user bytes under it survive. That test must stay green.
- The REQ-635 record's "Why" already notes that `sectionLineBounds` and `appendSectionEntry` (`state_apply.go`) match only the exact unindented `## X`: every generated writer emits its heading at column 0.
- Memory of REQ-634: never grow a parser grammar for how a model might format a file. The reviewer's "track list/container context" remedy was flagged on surface cost and is not the remedy.
## Detailed Requirements
1. Removal requires column 0: in `stripGeneratedRecoverySections`, a section whose heading line does not start at column 0 is never removed, whatever its name. The simplest carrier is for `VisibleSection` to expose the heading's indent (or a column-0 flag) computed in `VisibleSections`; a byte check at `body[section.Start]` is also acceptable if it reads clearly. Builder's choice.
2. Boundary recognition is unchanged: `VisibleSections` keeps the 0-3 space rule, so an indented user heading still ends the generated section before it. `TestRecoveryPreservesIndentedAndCommentedRequirements` and `TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule` stay green and untouched.
3. Lock-in test, written first and confirmed failing: a body with a bullet item followed by a two-space-indented `## Plan` and sample text (for example `- Example:\n  ## Plan\n  user sample\n\nMUST keep.\n`), then a generated column-0 `## Timing` section. After `stripGeneratedRecoverySections` the bullet, the indented `## Plan`, its sample, and `MUST keep.` survive byte for byte; `## Timing` is removed. Cover `\n` and `\r\n` like the sibling tests.
4. Check the other callers of `VisibleSections` that remove or replace a section by name (`replaceTimingSection` in `lifecycletiming/lifecycle_timing.go`, `advanceSections` in `lifecycleadvance/advance_commands.go`): if one of them would also overwrite or reject an indented user sample under this rule, apply the same column-0 requirement there with one test, or record in the hand-back why it is already safe. Do not widen beyond that.
5. Run the full Go test suite for the do-work CLI module.
6. Release per `_dev/primes/prime-releases.md`. The CHANGELOG entry names this as a data-loss fix in `recover --take-over`, the residual of the 0.305.70 (REQ-635) fix.
7. Lessons: if `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` gains an entry, refresh its token count in `do-work/lessons-index.md` in the same commit. An entry is earned here: the REQ-635 fix pinned the boundary rule but not the removal rule, and the two readers of one section list drifted (family `rule-direction-checked-against-callers` or `paired-predicate-drift`; builder picks the existing one that fits).
## Constraints
- Offsets and line endings stay untouched: `VisibleSections` keeps returning original-body byte offsets for both `\n` and `\r\n`.
- No list-marker or container parsing anywhere. Indent at the heading line is the whole signal.
- Batch constraint: this REQ is its own release and its own commit; do not fold REQ-646's change into it.
## Dependencies
None. REQ-646 (board activity snapshots direct merge matches) is related only by origin; disjoint files, may run in parallel.
## Builder Guidance
Certainty is high: the triage read the code and the test history, and the remedy is one condition in one function plus one test. Latitude is only where the indent is carried (a `VisibleSection` field versus a byte check) and whether requirement 4 needs a second change.
## Red-Green Proof
**RED prompt/case:** Wrap `# Request\n- Example:\n  ## Plan\n  user sample\n\nMUST keep.\n## Timing\ngenerated summary\n` in a REQ document and run `stripGeneratedRecoverySections` on it.
**Why RED now:** `VisibleSections` returns `Plan` and `Timing`; the name match removes both, so the bullet's sample and `MUST keep.` are deleted.
**GREEN when:** Only `## Timing` and its body are removed; the bullet, `  ## Plan`, `  user sample`, and `MUST keep.` remain byte for byte. `TestRecoveryPreservesIndentedAndCommentedRequirements` (indented `## Requirements` as a boundary) and the REQ-635 indent tests stay green.
**Validation:** User confirmed. The maintainer accepted the triage's remedy and RED/GREEN with "capture them and run them".
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18680 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its owning prime governs `requeststate` and `requestmodel`; family `rule-direction-checked-against-callers` names the drift this REQ fixes (boundary rule pinned, removal rule not), and `closed-enumeration-for-a-condition` holds the REQ-635 `TrimLeft` entry.
## Full Context
See `do-work/user-requests/UR-141/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (from the builder hand-back) Carry the heading indent on `VisibleSection` (the parser already measures it in `blockIndentContent`), and require `HeadingIndent == 0` wherever a section is removed, replaced or counted by name. Keep the 0-3 space boundary rule as is. Write the three lock-in tests first and confirm each fails on its assertion. No list or container parsing.
- [x] **[APPLY]:** (from the builder hand-back) One field plus one assignment in `VisibleSections`. One guard each in `stripGeneratedRecoverySections`, `replaceTimingSection` and `advanceSections`, with doc comments kept true. Three focused tests, one per caller, each naming REQ-645 and the failure it pins.
- [x] **[UNIFY]:** (from the builder hand-back) `git diff --stat` (base..880dcb54): 7 files, 62 insertions, 7 deletions; `git diff --check` clean; `gofmt -l .` prints nothing; `go build ./...` and `go vet ./...` pass. Checked `visible_sections.go` (offsets unchanged, `HeadingIndent = len(line) - len(content)` is the stripped 0-3 space count, CRLF unaffected), `state_apply.go` (indented section skipped before the name match, its bytes stay), `lifecycle_timing.go` (indented Timing sample skipped, falls through to a later column-0 Timing or append), `advance_commands.go` (indented headings are neither evidence nor duplicates, both callers `advance_commands.go:133` and `recovery_commands.go:188` behave the same), and the three test files (no debug artifacts, each test pins one caller's real failure).
*Source: review comment "[P2] Exclude list-nested headings from removable request sections — visible_sections.go:64-68", accepted by `do-work-toolbox validate-feedback` on 2026-10-08 with the column-0 removal remedy.*

---

## Triage

**Route: A** - Simple

**Reasoning:** Bug fix with an exact reproduction and a remedy settled at capture: one condition in `stripGeneratedRecoverySections` (plus an indent carrier on `VisibleSection` if the builder prefers), one lock-in test, and a bounded check of two sibling callers. The files are named; exploration would re-read what the triage already read.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/recovery_markdown_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/visible_sections_test.go` (modified)

**What was done:** `VisibleSection` now carries `HeadingIndent`, the 0-3 spaces before a heading's `##`, set by `VisibleSections` without changing boundary recognition. `stripGeneratedRecoverySections`, `replaceTimingSection` and `advanceSections` skip any section with a nonzero indent, so an indented user sample under a list item is never removed, overwritten, or counted as lifecycle evidence; each caller has one new lock-in test.

## Decisions

(from the builder hand-back)

- **D-01 (DECIDE & STATE):** The indent travels as `HeadingIndent int` on `VisibleSection`, not as a byte check at `body[section.Start]`. The parser already measures it, and three callers need it, so one field beats three byte checks. Value: one definition of "column 0". Risk: low, the field is additive and no literal of the struct exists in tests.
- **D-02 (DECIDE & STATE):** The field is an int, not a bool flag. The parser knows the count, and every caller compares it with zero. Reversible and leaf.
- **D-03 (DECIDE & STATE):** Requirement 4 changed both callers. Each one demonstrably overwrote or refused on the indented sample. That is the condition the REQ gave for a change, and each fix has its own lock-in test. The `write_set` in the REQ did not list these four files, but the brief's write boundary allows them when requirement 4 needs a change.

## Discovered Tasks

(from the builder hand-back; impact token added by the integrator)

- **impact-user-visible** `markdownSectionBytes` in `internal/requeststate/state_apply.go` (used by `markdownSectionExists` and `markdownSectionContains`) takes the first section by name with any indent. So an indented user `## Plan` sample can count as lifecycle evidence in requeststate checks, or hide the real column-0 section that comes later. It is read-only, so no data is lost, and requirement 4 excludes it. → report only

## Qualification

**Gate records:** `advance --diff-range c2793dab..b5e358f1` returned the `qualify` gate `satisfied` with provenance `merged_range` and no findings (no debug artifact, no output primitive, P-A-U boxes ticked).

**Requirement trace against the merged diff (7 files, 62 insertions, 7 deletions):**
1. Removal requires column 0: `state_apply.go` skips a section with `HeadingIndent != 0` before the generated-name match in `stripGeneratedRecoverySections`. The indent is carried on `VisibleSection` and set in `VisibleSections` as `len(line) - len(content)`, which is the 0-3 space count `blockIndentContent` stripped.
2. Boundary recognition unchanged: the `VisibleSections` loop and `blockIndentContent` are untouched apart from the new field. The three test files changed by additions only, so `TestRecoveryPreservesIndentedAndCommentedRequirements` and `TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule` are untouched, and both pass at the merge.
3. Lock-in test: `TestRecoveryKeepsAListNestedTwoSpaceIndentedHeading` is the Red-Green Proof body, wrapped in frontmatter, for `\n` and `\r\n`. The integrator re-ran it against the pre-fix source in a throwaway checkout: it failed with "recovery changed user text" and the bullet's sample deleted. It passes at the merge.
4. Sibling callers: both were unsafe and both now require column 0. Integrator red run on pre-fix source: `TestTimingReplacementKeepsAListNestedTwoSpaceIndentedSample` failed in both line endings with the sample overwritten, and `TestAdvanceIgnoresAListNestedTwoSpaceIndentedHeading` failed with the indented sample counted as Plan evidence. Both pass at the merge. `advanceSections` has two callers (`advance_commands.go` and `recovery_commands.go`), and both get the same behaviour.
5. Full module suite: the builder ran `go build/vet/test ./...` green in 58 s; the repository gate below covers it at the merge.
6, 7. Release and lesson: handled at finalization.

**Constraints:** offsets and line endings unchanged (the new field is additive and no offset arithmetic changed). No list-marker or container parsing was added anywhere.

**Scope:** `write_set` names `visible_sections.go`, `state_apply.go`, `recovery_markdown_test.go`, and the CHANGELOG. The four other touched files (`lifecycle_timing.go`, `lifecycle_timing_test.go`, `advance_commands.go`, `lifecycleadvance/visible_sections_test.go`) are requirement 4's bounded allowance: the requirement names both functions and says to apply the column-0 rule there with one test each when a caller would overwrite or reject the sample, which the red runs show. In scope, not drift (D-03).

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at merge `b5e358f1` (wall 121 s; do-work-cli module 879 tests in 59 s), then `advance --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-08-141302/REQ-645-probe.sh`.
**Result:** ✓ Repository gate exit 0 on the first run; `green-gate` and `run-blocked-check` satisfied; the focused probe (`go test ./internal/requeststate/ ./internal/requestmodel/ -run 'Recovery|VisibleSections|StripGenerated'`) exited 0.

**Red-green validation:** (traces to `## Red-Green Proof`; red re-run by the integrator against the pre-fix source with the new tests, green at the merge)
- `TestRecoveryKeepsAListNestedTwoSpaceIndentedHeading` (`requeststate/recovery_markdown_test.go`): ✗ before ("recovery changed user text", the bullet's `  ## Plan`, sample and `MUST keep.` deleted) → ✓ after. This is the captured RED case, wrapped in frontmatter, in `\n` and `\r\n`.
- `TestTimingReplacementKeepsAListNestedTwoSpaceIndentedSample` (`lifecycletiming/lifecycle_timing_test.go`): ✗ before (indented Timing sample overwritten, both line endings) → ✓ after.
- `TestAdvanceIgnoresAListNestedTwoSpaceIndentedHeading` (`lifecycleadvance/visible_sections_test.go`): ✗ before (indented sample counted as Plan evidence) → ✓ after.

**New tests added:**
- The three tests above, one per caller, each naming REQ-645 and the failure it pins.

Pinned boundary tests untouched and green at the merge: `TestRecoveryPreservesIndentedAndCommentedRequirements`, `TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule`, `TestRecoveryKeepsUserTextUnderAFourSpaceIndentedHeading`, `TestTimingReplacementPreservesIndentedAndCommentedRequirements`, `TestTimingReplacementIgnoresIndentedSamplesAndFenceInfoComments`, `TestAdvanceIgnoresExampleAndCommentHeadings`.

**Heavy verification plan:**
- Range: c2793dabdf98a95b3ebf7586afacaf54a8166f03..b5e358f1b25c0438e25ce90a52df85b397989d2d
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — all seven changed paths matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all seven changed paths matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — all seven changed paths matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — all seven changed paths matched subtree `skills/do-work/tools/do-work-cli`

*Verified by work action*

## Review

**Overall: 98%** | 2026-10-08T14:36:04Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- F1. The `VisibleSection` struct comment (`visible_sections.go:7-10`) says by-name counters must skip indented sections, but the read-only `markdownSectionBytes` (`state_apply.go:1015`, used for the Scope and Open Questions checks) still matches any indent. It is the same gap as the builder's discovered task. — impact-user-visible → report only
- Nit F2. One line of the `replaceTimingSection` doc comment (`lifecycle_timing.go:511`, about 150 characters) was not reflowed. — impact-negligible → report only

**Acceptance:** Pass. The four packages are green at b5e358f1. All three new tests fail on the pre-fix source (independent red run) and pass after the fix. The pinned REQ-635 and 0.305.59 tests are untouched and green.
**Restatement sweep:** redefined `VisibleSection` contract (`HeadingIndent`, and column 0 required for by-name writers and evidence readers). The only stale consumer is `markdownSectionBytes` (F1). Prime, lessons, CHANGELOG and guides are not stale. REQ-646 recorded nothing redefined, and a spot check found nothing stale.
**Suggested testing:** 1 item
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Carrying the indent the parser already measured as one field on `VisibleSection` gave three callers one definition of "column 0", instead of three byte checks at `body[section.Start]`.
**What didn't:** REQ-635 pinned the 0-3 space rule as a boundary rule only. Every caller that removed, replaced or counted a section by name reused the same list unchanged, so the boundary test passed while two more callers kept the same data-loss and refusal bugs.
**Worth knowing:** An indented heading is a section boundary but never a generated section. Read-only `markdownSectionBytes` in `state_apply.go` still matches any indent (discovered task and review F1, report only).

## Orientation

Recovery, the Timing writer and `advance` now treat a 1-3 space indented `##` heading as user text: it still ends the section above it, but it is never removed, replaced or counted as lifecycle evidence. Lives in the do-work CLI's request-markdown layer (`prime-do-work-cli.md`); the prime's referenced paths still exist and nothing in it went stale.

## Heavy Verification Plan

- Base: c2793dabdf98a95b3ebf7586afacaf54a8166f03
- Target: b5e358f1b25c0438e25ce90a52df85b397989d2d
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — all seven changed paths matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all seven changed paths matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — all seven changed paths matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — all seven changed paths matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: b5e358f1b25c0438e25ce90a52df85b397989d2d
- Execution revision: b5e358f1b25c0438e25ce90a52df85b397989d2d (detached drain checkout, `QUEUE_KANBAN_BROWSER` set)
- do-work-cli-integrations: executed, exit 0, 62 s
- staged-skills: executed, exit 0, 37 s
- updater: executed, exit 0, 65 s
- installer: executed, exit 0, 26 s

## Timing

Observed 2026-10-08T14:28:51Z to 2026-10-08T14:39:58Z: 11m 07s total, 9m 30s attributed across 4 events, 1m 37s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 6m 03s | 2 |
| review | 3m 10s | 1 |
| handback-merge | 17s | 1 |

Slowest stage: verification-gate / heavy drain, 3m 34s, outcome success.
