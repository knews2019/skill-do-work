---
id: REQ-635
title: '[impact-critical] recover --take-over deletes user text under an indented heading; VisibleSections must apply the CommonMark 0-3 space rule and not open a fence on a comment-opening line'
status: completed
route: A
estimate:
  p50_active_minutes: 25
  confidence: medium
  basis:
  - Route A
  - 7-file write set
  - 3 subsystems involved
  - 8 acceptance criteria
  calculated_at: 2026-10-07T13:25:40Z
priority: now
created_at: 2026-10-06T22:45:18Z
user_request: UR-137
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-636, REQ-637, REQ-638]
batch: review-0305-69-findings
required_lessons: [_dev/primes/lessons-releases.md]
write_set: [skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go, skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections_test.go, skills/do-work/tools/do-work-cli/internal/requeststate/recovery_markdown_test.go, skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing_test.go, skills/do-work/CHANGELOG.md, skills/do-work/tools/do-work-cli/lessons-do-work-cli.md, do-work/lessons-index.md]
dispatch_at: 2026-10-07T13:18:50Z
builder_handback_at: 2026-10-07T13:25:01Z
integration_at: 2026-10-07T13:25:36Z
review_at: 2026-10-07T13:31:01Z
commit: 2c4dae773e02234fe254f93044193bdc1908515a
heavy_verified_at: 2026-10-07T13:35:00Z
heavy_verified_revision: 2c4dae773e02234fe254f93044193bdc1908515a
kb_status: pending
claimed_at: 2026-10-07T13:16:10Z
completed_at: 2026-10-07T13:34:35Z
release_at: 2026-10-07T13:34:35Z
---
# Recover --take-over Deletes User Text Under an Indented Heading; VisibleSections Must Apply the CommonMark 0-3 Space Rule and Not Open a Fence on a Comment-Opening Line

## What
Two parser defects in `VisibleSections` (`skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go`), both introduced by the 0.305.59 change (commit 96505655). First, review finding F1 (high): line 60 trims `" \t"` before matching `## `, so a heading indented by a tab or by four or more spaces is treated as a section heading, and `stripGeneratedRecoverySections` (`recover --take-over`) deletes user request text sitting under an indented `## Plan` sample. Second, review finding F5: a line that opens an HTML comment can also open a fence in the same pass (lines 33-56), after which the comment never closes and every later heading, including `## Timing`, is hidden. Fix both, with the failing tests written first.

## Why
A consumer project installed 0.305.69 and reviewed the 0.305.58 to 0.305.69 changes. F1 was reproduced by the reviewer and again during triage here: the review's body yields a `Plan` section spanning `"    ## Plan\n    user sample\n\nMUST keep.\n"`, and recovery removes it. The same root cause makes `advanceSections` (`lifecycleadvance/advance_commands.go:344-360`) refuse with `duplicate lifecycle section Review` and lets `replaceTimingSection` (`lifecycletiming/lifecycle_timing.go:514-525`) overwrite an indented `## Timing` sample. It also disagrees with `sectionLineBounds` and `appendSectionEntry` (`requeststate/state_apply.go:1036-1049`), which match only the exact unindented `## X`. F5 was also reproduced during triage: the review's body returns zero sections, so the Timing writer appends a duplicate section. The lessons file already warns about this exact trim (`lessons-do-work-cli.md:66`, "`TrimLeft(" \t")` made indented code unreachable"), so this is a repeat.

## Detailed Requirements
1. F1 heading rule: a `## ` heading is recognised only with 0 to 3 leading spaces. A tab anywhere in the indentation, or 4 or more spaces, is not a heading (CommonMark ATX rule). Keep the existing tests green: 1, 2 and 3 leading spaces at `recovery_markdown_test.go:13-15`, and the 2-space `## Requirements <!-- retained -->` case at `visible_sections_test.go:10`.
2. F1 negative tests, written first and confirmed failing: 4 leading spaces, a leading tab, and the review's recovery example body (the text `MUST keep.` and the indented sample must survive `stripGeneratedRecoverySections`).
3. Same rule for fences: `leadingPunctuationRun` (line 81) trims leading whitespace without bound before measuring a fence, so a 4-space-indented ```` ``` ```` inside an indented code block opens a fence and hides later headings. Apply the same 0-3 space bound there. This was found during triage, not by the review; it is the same CommonMark rule and the same fix.
4. F5: a line that opens a fence must not also be scanned for a comment opener. Evaluate the fence rule before the comment scan. In CommonMark a fence opener wins and `<!--` on that line is only its info string, so this is the spec-faithful choice. With the review's F5 body the result is exactly two sections, `Plan` and `Timing`. Test written first and confirmed failing.
5. Do NOT change: `## Plan <!-- note` with an unclosed comment still yields a zero-length section, because the unclosed region is protected by design.
6. Run the full Go test suite for the do-work CLI module.
7. Release per `_dev/primes/prime-releases.md`. The CHANGELOG entry must name F1 as a data-loss fix in `recover --take-over`.
8. Lessons: add an entry to `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` recording that the `TrimLeft(" \t")`-before-the-test trap from REQ-460 (making indented code unreachable) fired a second time, this time in heading recognition at 0.305.59; refresh the satellite's token count in `do-work/lessons-index.md` in the same commit.

## Constraints
- Offsets and line endings stay untouched: `VisibleSections` must keep returning original-body byte offsets for both `\n` and `\r\n` inputs.
- No normalisation of the heading name; keep the raw name behaviour the 0.305.59 comment at line 58-59 describes.
- Batch constraint: each REQ in this batch is its own release and its own commit; do not fold another REQ's change into this one.

## Dependencies
None. REQ-636 (board activity git reads), REQ-637 (board gap readers) and REQ-638 (release guard simplification) are independent slices of the same review.

## Builder Guidance
Certainty is high; both bugs are reproduced and the intended behaviour is stated by the reviewer. Builder latitude: how to count leading spaces (a small helper shared by the heading check and `leadingPunctuationRun` is the obvious shape), and the exact wording of the lesson entry. The F5 ordering choice (fence before comment) was made at triage; keep it unless a test shows it breaks an existing case.

## Red-Green Proof
**RED prompt/case:** Call `VisibleSections` on the F1 body `"# Request\nExample:\n\n    ## Plan\n    user sample\n\nMUST keep.\n## Timing\nt\n"` and on the F5 body `"# Request\n``` <!--\n## X\n```\n## Plan\np\n## Timing\nt\n"`. Run `stripGeneratedRecoverySections` on the F1 body wrapped in a REQ document.
**Why RED now:** F1 returns a `Plan` section spanning `"    ## Plan\n    user sample\n\nMUST keep.\n"`, and recovery deletes those bytes. F5 returns no sections at all.
**GREEN when:** F1 returns only `Timing`, and recovery leaves `"    ## Plan\n    user sample\n\nMUST keep.\n"` byte-for-byte in place. F5 returns exactly `Plan` and `Timing`. A tab-indented `\t## Plan` returns no `Plan` section. The existing 1-, 2- and 3-space tests and the `## Plan <!-- note` zero-length case stay green.
**Validation:** Inferred during capture. Both RED cases were reproduced on 0.305.69 during triage with a copy of the parser; the reviewer reproduced them independently.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18303 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matched: family `closed-enumeration-for-a-condition` at line 66 names the exact `TrimLeft(" \t")` trap this REQ fixes a second time.

## Full Context
See `do-work/user-requests/UR-137/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** One helper `blockIndentContent` (0-3 spaces; a tab in the indent or a fourth space means indented code) used by the heading check and by `leadingPunctuationRun`; the fence-opener check moves before the comment scan for lines that start outside a comment; the opener scan uses `commentOpenerIndex`, which skips inline code spans on the same line; the `-->` close scan is unchanged. (From the builder hand-back.)
- [x] **[APPLY]:** Done as planned in `visible_sections.go` plus new tests in the three declared test files; doc comments state the conditions. (From the builder hand-back.)
- [x] **[UNIFY]:** `git diff --stat` 4 files, 187 insertions, 10 deletions, each reviewed; `gofmt -l` empty, `go vet ./...` clean, full module `go test -count=1 ./...` 33 packages green (94s wall); no debug artifacts. (From the builder hand-back.)
*Source: review findings F1 and F5 from the pasted consumer code review of 0.305.58 to 0.305.69, triaged with `do-work validate-feedback` on 2026-10-06 and accepted; the review set "Work in priority order" with F1 first, which is the `priority: now` evidence.*

---

## Triage

**Route: A** - Simple

**Reasoning:** Bug fix with exact reproductions in the REQ: one parser file (`visible_sections.go`), the two RED bodies, the intended CommonMark rule, and the F5 ordering choice were all settled at capture. Exploration would re-discover what triage already read; the callers (`replaceTimingSection`, `stripGeneratedRecoverySections`, `advanceSections`) only consume the section list.

**Planning:** Not required

**Run note:** `advance` first refused this REQ at the triage phase with the Triage section present, because Detailed Requirements item 4 quotes `<!--` inside an inline code span and `VisibleSections` opened an unclosed comment there, hiding every later section. The builder was dispatched before the estimate so its fix could make the REQ readable; estimate and the plan skip note were then recorded with the builder's CLI before the merge (D-01, D-02).

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/recovery_markdown_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing_test.go` (modified)

**What was done:** `VisibleSections` now recognises a `## ` heading, and `leadingPunctuationRun` a fence opener or closer, only after zero to three leading spaces; a tab in the indentation or a fourth space makes the line indented code (new helper `blockIndentContent`, which measures the indent before trimming anything). On a line that starts outside a comment the fence-opener rule runs before the comment scan, so a fence line's `<!--` info string no longer opens a comment. The opener scan uses the new `commentOpenerIndex`, which treats `<!--` inside a same-line inline code span as literal; an unmatched backtick run stays literal, so the conservative hidden behaviour holds there. Unclosed `## Plan <!-- note` still yields a zero-length section; offsets and `\r\n` handling are unchanged. New tests pin each rule over both line endings, recovery keeps user text under a four-space-indented `## Plan`, and the Timing writer neither overwrites an indented sample nor appends a duplicate behind a fence info comment.

## Qualification

**Range:** `4337e9ac..2c4dae77` (builder commit 9f8d55e6 merged as 2c4dae77)
**Gate record:** `advance` qualification satisfied with no findings: no debug artifacts, no leftover output primitives, every P-A-U box ticked from the hand-back, no `do-work/` path inside the range.
**Orchestrator judgment:** The diff is substantive and traces to every requirement. R1/R2: the heading check goes through `blockIndentContent`, so a tab or four or more spaces is no heading, and `TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule` plus `TestRecoveryKeepsUserTextUnderAFourSpaceIndentedHeading` pin the negatives and the review's recovery body; the 1/2/3-space and `## Requirements <!-- retained -->` tests are unchanged and green. R3: `leadingPunctuationRun` uses the same helper, pinned by the four-space fence opener and closer cases. R4: the fence check runs before the comment scan only for lines that start outside a comment, pinned by the F5 body (exactly `Plan`, `Timing`) and the Timing writer test. R5: `TestVisibleSectionsKeepAnUnclosedCommentHeadingZeroLength` pins the unchanged zero-length case. D-01's code-span rule is a same-line CommonMark span match; an unmatched run stays literal, which keeps the conservative hidden behaviour. Live flow verified: the merged CLI now reads this REQ, whose own Detailed Requirements item 4 had hidden every later section. R6 (full module) is in the hand-back; R7 and R8 (release, lessons) are orchestrator work at finalization.

## Testing

**Tests run:** `bash do-work/runs/work-2026-10-07-131606/REQ-635-probe.sh` (go test, `-run 'VisibleSections|Recovery|TimingReplacement'` over requestmodel, requeststate and lifecycletiming) through the `advance` test gate: probe status 0, `run-blocked-check` and `green-gate` records satisfied. Builder-side: full module `go test -count=1 ./...` 33 packages green, 94s wall; the three touched packages 0.9s, 8.4s and 4.3s; gofmt and go vet clean.
**Result:** ✓ All passing
**Repository gate:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at 2c4dae77 from the detached worktree `.git/work-run-2026-10-07-131606/drain-head`, exit 0 on the first run, gate wall 149s; both Go stages EXECUTING (reuse_disabled); slowest files strict_behavior_regression_test.go 20.74s and finalization_recovery_test.go 24.48s, both under the 30s budget. Load average about 7 at launch, no other gate running.

**Red-green validation:** (traces to `## Red-Green Proof`)
- `TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule` (visible_sections_test.go): ✗ before, F1 body `got ["Plan" "Timing"], want ["Timing"]`; `\t## Plan`, ` \t## Plan`, 5 spaces `got ["Plan"], want []`; four-space fence opener `got [], want ["Plan"]`; four-space closer `got ["Hidden"], want ["Plan"]` → ✓ after
- `TestVisibleSectionsReadFenceOpenersAndCodeSpansBeforeComments`: ✗ before, F5 body `got [], want ["Plan" "Timing"]`; code-span `<!--` `got [], want ["Plan"]` → ✓ after
- `TestRecoveryKeepsUserTextUnderAFourSpaceIndentedHeading` (recovery_markdown_test.go): ✗ before, the indented sample and `MUST keep.` deleted → ✓ after, bytes kept exactly. This is the captured recovery RED/GREEN pair.
- `TestTimingReplacementIgnoresIndentedSamplesAndFenceInfoComments` (lifecycle_timing_test.go): ✗ before, "indented Timing sample was replaced" and "Timing section behind a fence info comment was not replaced" → ✓ after
- `TestVisibleSectionsKeepAnUnclosedCommentHeadingZeroLength`: new pin, green before and after (requirement 5, unchanged behaviour)
- The reviewer copied the four new tests onto 4337e9ac (all fail) and ran them at 2c4dae77 (all pass).

**New tests added:**
- `TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule`, `TestVisibleSectionsReadFenceOpenersAndCodeSpansBeforeComments`, `TestVisibleSectionsKeepAnUnclosedCommentHeadingZeroLength` (visible_sections_test.go)
- `TestRecoveryKeepsUserTextUnderAFourSpaceIndentedHeading` (recovery_markdown_test.go)
- `TestTimingReplacementIgnoresIndentedSamplesAndFenceInfoComments` (lifecycle_timing_test.go)

**Heavy verification plan:**
- Range: 4337e9ac..2c4dae77 (4 of 6 lanes; every changed path sits under `skills/do-work/tools/do-work-cli`)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the four changed files matched subtree skills/do-work/tools/do-work-cli
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — the four changed files matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the four changed files matched subtree skills/do-work/tools/do-work-cli
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the four changed files matched subtree skills/do-work/tools/do-work-cli

*Verified by work action*

## Review

**Overall: 96%** | 2026-10-07T13:29:46Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- M1: the 0-3 space rule for fences is measured from the page column, not the list-item content column (`visible_sections.go:90`, `:106`). A fence inside a nested list item at 4+ spaces now protects nothing, so an HTML sample inside it (`    <!-- example`) opens an unclosed comment and hides every later section (probe: 0.305.69 parser returns `Plan`, merged parser returns none). A list fence opened at 2 spaces and closed at 4 spaces never closes (D-03), with the same effect. Hide direction only, so no user bytes are lost; the symptom is `advance` missing a section or the Timing writer appending a duplicate `## Timing`. — impact-user-visible → report only
- M2: Restatement Sweep: `VisibleSections` now admits a heading only at 0-3 spaces, but `sectionLineBounds` / `appendSectionEntry` (`requeststate/state_apply.go:1036-1049`) still admit only the exact unindented `## X`, so 1-3 space headings are still read differently by the two readers. The REQ's Why names this pre-existing gap; no requirement asked to close it. — impact-negligible → report only

**Acceptance:** Pass — RED reproduced at 4337e9ac and GREEN at 2c4dae77 for the F1 body (LF and CRLF), the F5 body, tab and space-tab indents; all five new tests fail before and pass after; 4 affected packages green; 22 adversarial cases checked.
**Suggested testing:** 2 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

## Heavy Verification Plan

**Base revision:** 4337e9ac3379c291e4bbe2f7ed0d87397d5a9ccb
**Target revision:** 2c4dae773e02234fe254f93044193bdc1908515a
**Selected lanes** (4 of 6; every changed path sits under `skills/do-work/tools/do-work-cli`):
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the four changed files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — the four changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the four changed files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the four changed files matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

**Target revision:** 2c4dae773e02234fe254f93044193bdc1908515a
**Execution revision:** 2c4dae773e02234fe254f93044193bdc1908515a (detached worktree `.git/work-run-2026-10-07-131606/drain-head`, `QUEUE_KANBAN_BROWSER` set, 2026-10-07T13:30:50Z to 13:34:10Z)
- do-work-cli-integrations: executed, exit 0, 67s
- staged-skills: executed, exit 0, 36s
- updater: executed, exit 0, 65s
- installer: executed, exit 0, 28s

## Lessons Learned

**What worked:** A Route A dispatch with the design fixed in the brief gave a first-try green build (4 files, 187 insertions) with RED shown for every rule. Running the repository gate and the heavy drain from one detached worktree at the merge revision while the review ran.
**What didn't:** The REQ could not be advanced at all: its own Detailed Requirements item 4 quotes `<!--` inside backticks, the pre-fix parser opened a comment there, and `advance` could not see the `## Triage` section. Waiting for the merge was not needed; the builder's worktree CLI read the REQ correctly, so estimate and plan were recorded with it before the merge (D-01, D-02).
**Worth knowing:** The section reader has three independent rules now: block indent (0-3 spaces, measured before any trim), one construct per line (fence opener before comment scan), and code spans quote `<!--`. List-item fences indented four or more spaces are now read as indented code (D-03, review M1); a column model per list item would fix that, but nothing has been lost from it so far, only hidden. `sectionLineBounds` and `appendSectionEntry` still read only unindented headings (review M2).

## Orientation

Now `recover --take-over` keeps user text under an indented `## ` sample, and `advance` and the Timing writer read the same sections a Markdown renderer shows: indented code, fence info strings and inline code spans no longer create or hide sections. Lives in the do-work-cli request model (`prime-do-work-cli.md`), shared by recovery, lifecycle advance and lifecycle timing. No map change: one parser function got stricter rules; no new module, data flow or contract. Prime spot-check: `prime-do-work-cli.md` and `prime-releases.md` name no path this change moved or removed, so neither is stale.

## Decisions

- **D-01** (DECIDE & STATE, orchestrator, triage): added a third rule to this REQ's scope: a `<!--` inside a same-line inline code span is not a comment opener. Reasoning: same function and same defect family (the comment scan opens on bytes that are not a comment), and it blocked this REQ's own lifecycle, because Detailed Requirements item 4 quotes `<!--` in backticks. The reviewer confirmed it is CommonMark-correct and cannot un-hide a heading a renderer hides; an unmatched backtick run stays literal, so the conservative direction holds.
- **D-02** (DECIDE & STATE, orchestrator): the builder was dispatched before the estimate, and estimate plus the plan skip note were recorded with the builder worktree's CLI before the merge, so `dispatch_at` precedes `estimate.calculated_at`. Reasoning: the canonical command still owned each record; only the parser it ran was the fixed one. Estimate is informational.
- **D-03** (DECIDE & STATE, builder): a fence closer indented four or more spaces no longer closes an open fence (CommonMark). Side effect recorded as review M1.
- **D-04** (DECIDE & STATE, builder): the CRLF Timing-append test checks "sample kept, one real `## Timing` appended" instead of exact bytes, because `replaceTimingSection` joins with `\n` (discovered task DT1).

## Discovered Tasks

- DT1: `replaceTimingSection` (`lifecycletiming/lifecycle_timing.go`) joins the new section with `\n` on CRLF documents, leaving mixed line endings. — impact-negligible → report only
- DT2 (review M1): the 0-3 space rule is measured from the page column, not the list-item content column, so a fence inside a nested list item at four or more spaces protects nothing, and a `<!--` sample inside it hides every later section. Hide direction only, no user bytes lost. — impact-user-visible → report only
- DT3 (review M2): `sectionLineBounds` / `appendSectionEntry` (`requeststate/state_apply.go:1036-1049`) still accept only an unindented `## X`, so 1-3 space headings are read differently by the two readers. Pre-existing. — impact-negligible → report only
- DT4: a `<!--` in an indented code block still opens a comment (pre-existing, reviewer probe). Hide direction only. — impact-negligible → report only

## Timing

Observed 2026-10-07T13:18:50Z to 2026-10-07T13:34:23Z: 15m 33s total, 19m 12s attributed across 4 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 8m 00s | 2 |
| builder-work | 6m 38s | 1 |
| review | 4m 34s | 1 |

Slowest stage: builder-work / builder subagent in worktree, commit 9f8d55e6, 6m 38s, outcome success.
