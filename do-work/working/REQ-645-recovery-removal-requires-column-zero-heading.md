---
id: REQ-645
title: '[impact-critical] recover --take-over keeps a 1-3 space indented generated-name heading; only a column-0 section is removable'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: review comment "[P2] Exclude list-nested headings from removable request sections — visible_sections.go:64-68", accepted by `do-work-toolbox validate-feedback` on 2026-10-08 with the column-0 removal remedy.*

---

## Triage

**Route: A** - Simple

**Reasoning:** Bug fix with an exact reproduction and a remedy settled at capture: one condition in `stripGeneratedRecoverySections` (plus an indent carrier on `VisibleSection` if the builder prefers), one lock-in test, and a bounded check of two sibling callers. The files are named; exploration would re-read what the triage already read.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
