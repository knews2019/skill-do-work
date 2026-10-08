## Review: REQ-645

**Approve** — `recover --take-over` no longer deletes a two-space `## Plan` sample under a bullet, and the two sibling callers that had the same bug are fixed too.
Route A | merge range c2793dab..b5e358f1 (builder commit 880dcb54)

### What's built
- `VisibleSection` now carries `HeadingIndent` (the 0-3 spaces before `##`). Recovery (`stripGeneratedRecoverySections`), the Timing writer (`replaceTimingSection`) and lifecycle evidence (`advanceSections`) act only on column-0 sections. An indented heading still ends the section above it, so the REQ-635 boundary rule is unchanged.
- Still open, by design: the read-only `markdownSectionBytes` matches any indent. The builder reported it as a discovered task.

### Decisions / risks for you
- None. The field is additive, the only place that builds the struct is `VisibleSections`, and byte offsets did not change.

### Findings

**Important:**
- None.

**Minor:**
- F1. The new `VisibleSection` struct comment (`internal/requestmodel/visible_sections.go:7-10`) says every reader that "counts a section by name must skip one with a nonzero indent". `markdownSectionBytes` (`internal/requeststate/state_apply.go:1015`, used by `markdownSectionExists("Scope")` and `markdownSectionContains("Open Questions")`) does not skip it. An indented `  ## Scope` or `  ## Open Questions` sample can count as evidence or hide the real column-0 section that comes later. Nothing is deleted. Requirement 4 kept this out of scope, and it matches the builder's discovered task. — impact-user-visible → report only

**Nit:**
- F2. The `replaceTimingSection` doc comment was not reflowed. Line 511 of `internal/lifecycletiming/lifecycle_timing.go` is about 150 characters long, while the lines around it are about 80. — impact-negligible → report only

### Requirements Checklist

- [x] 1. Removal requires column 0: `stripGeneratedRecoverySections` skips `HeadingIndent != 0` before the name match. Delivered.
- [x] 2. Boundary recognition unchanged: the `VisibleSections` loop and `blockIndentContent` are untouched apart from the new field. All three test files changed by additions only. `TestRecoveryPreservesIndentedAndCommentedRequirements`, `TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule` and `TestRecoveryKeepsUserTextUnderAFourSpaceIndentedHeading` are not in the diff and pass. Delivered.
- [x] 3. Lock-in test `TestRecoveryKeepsAListNestedTwoSpaceIndentedHeading` covers both `\n` and `\r\n`. The bullet, `  ## Plan`, `  user sample` and `MUST keep.` survive byte for byte, and the column-0 `## Timing` is removed. Delivered.
- [x] 4. Sibling callers. Both verdicts are true:
  - `replaceTimingSection`: before the fix it replaced the first `Timing` section of any indent, so a two-space sample was overwritten. Reproduced. It now requires column 0 and falls through to a later column-0 section or appends.
  - `advanceSections`: before the fix an indented `## Plan` counted as Plan evidence on its own, and next to the generated Plan it caused a "duplicate lifecycle section" refusal. Reproduced. Both callers (`advance_commands.go:133` and `recovery_commands.go:188` `heldForHeavyLanes`) now ignore indented headings, which is correct because generated writers emit only column-0 headings.
  - The four files outside `write_set` (`lifecycle_timing.go`, `lifecycle_timing_test.go`, `advance_commands.go`, `lifecycleadvance/visible_sections_test.go`) fall inside requirement 4's allowance. They are in scope, not drift.
- [x] 5. Full module suite: the builder ran it green in 58 s, and the repository gate passed at the merge.
- [ ] 6, 7. Release (CHANGELOG) and lessons entry: N/A at review time. Finalization handles them after this review. The builder's hand-back proposes the text for both.
- [x] Constraint: no list-marker or container parsing was added. The whole signal is `len(line) - len(content)` on a line that already passed `blockIndentContent`.
- [x] Constraint: offsets and line endings are unchanged.

### Acceptance Testing

**Result: Pass**
- At the merge, `go test ./internal/requeststate/ ./internal/requestmodel/ ./internal/lifecycletiming/ ./internal/lifecycleadvance/ -count=1` passes in all four packages. HEAD == b5e358f1, and there is no diff under `skills/`.
- With `-v` and the new tests plus the pinned tests: all PASS.
- I ran an independent red check. I exported b5e358f1 to the scratchpad, restored the four pre-fix source files from c2793dab and kept the new tests. All three new tests fail with the expected messages: "recovery changed user text" (the sample and `MUST keep.` deleted), "indented Timing sample was replaced" (in both line endings), and "indented sample became Plan evidence".

### Doc comments
- `VisibleSection` struct comment: accurate for the three fixed callers. It is overstated against `markdownSectionBytes` (F1).
- `VisibleSections` doc comment: unchanged and still true. It describes recognition, which did not change.
- `replaceTimingSection`: accurate, but one line was not reflowed (F2).
- `advanceSections`: the new comment is accurate.
- `stripGeneratedRecoverySections`: the new inline comment is accurate.

### Restatement Sweep
- This REQ redefined: the `VisibleSection` contract (`HeadingIndent`, and column 0 required for by-name writers and evidence readers).
- I grepped for `VisibleSection`, `HeadingIndent`, the caller names, "generated section", "by name", "0-3 space", "zero to three" and recovery-strip prose. Inside Go, the only stale consumer is `markdownSectionBytes` (F1). Outside Go:
  - `prime-do-work-cli.md` does not restate the contract.
  - The `lessons-do-work-cli.md` 0.305.70 entry and the CHANGELOG 0.305.70 entry are history and correct about that release.
  - `work-guide.md:130` and `work-reference.md:310` say take-over strips "orchestrator sections". That is still true, because those sections are column 0.
- REQ-646 (board activity snapshots direct merge matches) recorded "nothing redefined" on its Restatement sweep line, so the inherited trigger set is empty. As a spot check, `correlateCommitsToRequests`'s doc comment, the snapshot comment, `board-guide.md:53` ("a builder commit inside its merge") and the `git-history-evidence` lessons all agree with the fix. Nothing is stale.

### Suggested Additional Testing
- Edge case: a generated column-0 section whose own text holds an indented generated-name heading (for example `  ## Testing` inside a Plan). After take-over that tail now stays in the file instead of being stripped. That is leftover text, not data loss, and is worth one manual look.

### Scores (on the record, not the headline)

**Overall: 98%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Reqs 1-5 delivered. 6-7 belong to finalization |
| Code Quality | 95% | Minimal, one shared field. One comment line not reflowed, and the struct comment is overstated (F1) |
| Test Adequacy | 100% | One red-green test per caller. Red reproduced independently |
| Scope | 100% | The extra files are requirement 4's bounded allowance |
| Risk | Low | Additive struct field, offsets unchanged |
| Acceptance | Pass | Four packages green, red reproduced on the pre-fix source |

### Follow-ups created
- None (2 findings report only)

## Review

**Overall: 98%** | <timestamp>

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
