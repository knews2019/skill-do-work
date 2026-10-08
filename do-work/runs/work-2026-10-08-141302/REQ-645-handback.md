# REQ-645 hand-back (builder)

- Branch: `worktree-agent-REQ-645-recovery-removal-requires-column-zero-heading`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-645-recovery-removal-requires-column-zero-heading`
- Base: `b38928de`
- Commit (single): `880dcb54` — `[REQ-645] only a column-0 section is removable or replaceable by name`

## File manifest

All paths under `skills/do-work/tools/do-work-cli/`.

- `internal/requestmodel/visible_sections.go` (modified): `VisibleSection` gains `HeadingIndent int` (the 0-3 spaces before `##`), set by `VisibleSections` from the already measured `blockIndentContent` result; the struct comment states the column-0 contract for writers. Boundary recognition unchanged.
- `internal/requeststate/state_apply.go` (modified): `stripGeneratedRecoverySections` skips a section whose `HeadingIndent != 0` before the generated-name match.
- `internal/requeststate/recovery_markdown_test.go` (modified): new `TestRecoveryKeepsAListNestedTwoSpaceIndentedHeading` (the REQ's Red-Green Proof body, `\n` and `\r\n`).
- `internal/lifecycletiming/lifecycle_timing.go` (modified, requirement 4): `replaceTimingSection` replaces only a column-0 `## Timing`; doc comment updated.
- `internal/lifecycletiming/lifecycle_timing_test.go` (modified, requirement 4): new `TestTimingReplacementKeepsAListNestedTwoSpaceIndentedSample`.
- `internal/lifecycleadvance/advance_commands.go` (modified, requirement 4): `advanceSections` counts only column-0 sections; new doc comment.
- `internal/lifecycleadvance/visible_sections_test.go` (modified, requirement 4): new `TestAdvanceIgnoresAListNestedTwoSpaceIndentedHeading`.

`internal/requestmodel/visible_sections_test.go` was not touched: no test builds a `VisibleSection` literal, so the new field needed no fixture change.

## P-A-U

- **[PLAN]:** Carry the heading indent on `VisibleSection` (the parser already measures it in `blockIndentContent`), and require `HeadingIndent == 0` wherever a section is removed, replaced or counted by name. Keep the 0-3 space boundary rule as is. Write the three lock-in tests first and confirm each fails on its assertion. No list or container parsing.
- **[APPLY]:** One field plus one assignment in `VisibleSections`. One guard each in `stripGeneratedRecoverySections`, `replaceTimingSection` and `advanceSections`, with doc comments kept true. Three focused tests, one per caller, each naming REQ-645 and the failure it pins.
- **[UNIFY]:**
  - `git diff --stat` (base..880dcb54): 7 files, 62 insertions, 7 deletions. `git diff --check` clean.
  - `gofmt -l .` prints nothing. `go build ./...` and `go vet ./...` pass.
  - Checked `visible_sections.go`: offsets unchanged, `HeadingIndent = len(line) - len(content)` is the stripped space count (0-3) and is only set on lines that passed `blockIndentContent`. The CRLF case is unaffected because the indent is leading spaces only.
  - Checked `state_apply.go`: the indented section is skipped before the name match. Its bytes stay, and the earlier section still ends at it, so `TestRecoveryPreservesIndentedAndCommentedRequirements` stays green.
  - Checked `lifecycle_timing.go`: an indented Timing sample is skipped and the loop falls through to a later column-0 Timing or to append.
  - Checked `advance_commands.go`: indented headings are not evidence and not duplicates. Both callers, `advance_commands.go:133` and `recovery_commands.go:188`, get the same behaviour.
  - Checked the three test files: no debug artifacts, and each test pins one caller's real failure.

## Red-green evidence

Before the fix, from `go test ./internal/requeststate/ ./internal/lifecycletiming/ ./internal/lifecycleadvance/ -run ListNested -count=1`:

```
--- FAIL: TestRecoveryKeepsAListNestedTwoSpaceIndentedHeading (0.00s)
    recovery_markdown_test.go:105: recovery changed user text: error=<nil>, got "---\nid: REQ-645\n---\n# Request\n- Example:\n", want "---\nid: REQ-645\n---\n# Request\n- Example:\n  ## Plan\n  user sample\n\nMUST keep.\n"
--- FAIL: TestTimingReplacementKeepsAListNestedTwoSpaceIndentedSample (0.00s)
    lifecycle_timing_test.go:63: indented Timing sample was replaced instead of a section appended: got "# Request\n- Example:\n## Timing\nnew summary\n", want "# Request\n- Example:\n  ## Timing\n  sample summary\n\n## Timing\nnew summary\n"
    lifecycle_timing_test.go:63: indented Timing sample was replaced instead of a section appended: got "# Request\r\n- Example:\r\n## Timing\nnew summary\n", want "# Request\r\n- Example:\r\n  ## Timing\r\n  sample summary\r\n\n## Timing\nnew summary\n"
--- FAIL: TestAdvanceIgnoresAListNestedTwoSpaceIndentedHeading (0.00s)
    visible_sections_test.go:23: indented sample became Plan evidence: map[string]lifecycleadvance.sectionEvidence{"Plan":lifecycleadvance.sectionEvidence{start:11, end:36, count:1}},
```

After the fix, same command with `-v`:

```
--- PASS: TestRecoveryKeepsAListNestedTwoSpaceIndentedHeading (0.00s)
--- PASS: TestTimingReplacementKeepsAListNestedTwoSpaceIndentedSample (0.00s)
--- PASS: TestAdvanceIgnoresAListNestedTwoSpaceIndentedHeading (0.00s)
```

Pinned boundary tests stay green and untouched: `TestRecoveryPreservesIndentedAndCommentedRequirements`, `TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule`, `TestRecoveryKeepsUserTextUnderAFourSpaceIndentedHeading`, `TestTimingReplacementPreservesIndentedAndCommentedRequirements`, `TestTimingReplacementIgnoresIndentedSamplesAndFenceInfoComments`, `TestAdvanceIgnoresExampleAndCommentHeadings`.

## Requirement 4 verdicts

- **`replaceTimingSection` was unsafe, now fixed.** The line read was `if visible.Name != "Timing" { continue }` at `lifecycle_timing.go:516`. The first `Timing` section of any indent was replaced, so a two-space `## Timing` sample under a bullet was overwritten with the generated summary. The red output above shows this. It now also requires `visible.HeadingIndent == 0`.
- **`advanceSections` was unsafe, now fixed.** The line read was `for _, visible := range requestmodel.VisibleSections(body) { section := sections[visible.Name]; section.count++ ...` at `advance_commands.go:346`. An indented user `## Plan` counted as Plan evidence on its own. Beside the generated column-0 Plan it made `count == 2`, so advance returned "duplicate lifecycle section Plan" and refused. Indented sections are now skipped.

## Decisions

- **D-01 (DECIDE & STATE):** The indent travels as `HeadingIndent int` on `VisibleSection`, not as a byte check at `body[section.Start]`. The parser already measures it, and three callers need it, so one field beats three byte checks. Value: one definition of "column 0". Risk: low, the field is additive and no literal of the struct exists in tests.
- **D-02 (DECIDE & STATE):** The field is an int, not a bool flag. The parser knows the count, and every caller compares it with zero. Reversible and leaf.
- **D-03 (DECIDE & STATE):** Requirement 4 changed both callers. Each one demonstrably overwrote or refused on the indented sample. That is the condition the REQ gave for a change, and each fix has its own lock-in test. The `write_set` in the REQ did not list these four files, but the brief's write boundary allows them when requirement 4 needs a change.

## Discovered Tasks

- `markdownSectionBytes` in `internal/requeststate/state_apply.go` (used by `markdownSectionExists` and `markdownSectionContains`) takes the first section by name with any indent. So an indented user `## Plan` sample can count as lifecycle evidence in requeststate checks, or hide the real column-0 section that comes later. It is read-only, so no data is lost, and requirement 4 excludes it. → report only

## Lessons read

- Satellite `lessons-do-work-cli.md`: the `[family: closed-enumeration-for-a-condition]` bullets. These were the 0.305.70 REQ-635 entry (this bug's parent), 0.303.7, the two REQ-460 bullets, REQ-461 and 0.283.0.
- Family `[family: rule-direction-checked-against-callers]`: the 0.305.30 REQ-604 entry.
- Prime `prime-do-work-cli.md` Traps, and the crew rules general, coding-guardrails, testing and communication-style.

## Proposed lesson bullet (integrator writes it)

Family chosen: `rule-direction-checked-against-callers`. `paired-predicate-drift` has no entry in this satellite, and the REQ asks for an existing family. This one fits because REQ-635 set the indent rule's direction for boundaries and checked it only against the boundary callers, never against the callers that remove, replace or count by name.

```
- [family: rule-direction-checked-against-callers] <version> (REQ-645, 2026-10-08): **a boundary is not ownership.** REQ-635 kept CommonMark's 0-3 space rule so an indented user heading ends the generated section above it, and pinned that with a test, but every reader that removed, replaced or counted a section by name reused the same section list unchanged. A two-space `## Plan` sample under a bullet was therefore deleted by `recover --take-over`, a two-space `## Timing` sample was overwritten by `replaceTimingSection`, and `advanceSections` refused the request as a duplicate Plan. Generated writers emit only column-0 headings, so `VisibleSection` now carries `HeadingIndent` and every by-name writer or evidence reader requires zero. When one parser answers two questions, pin each question with its own test at each caller; a test of the first does not cover the second. Read-only `markdownSectionBytes` still matches any indent (reported, not fixed).
```

## Proposed CHANGELOG entry (integrator writes it)

```
## <version> — Recover --take-over Keeps an Indented Heading Sample Under a List Item

Data-loss fix in `recover --take-over`, finishing the 0.305.70 (REQ-635) fix. A heading indented by one to three spaces, such as a two-space `## Plan` sample under a bullet, still ends the section above it. It is now never removed as a generated section, because generated sections always start at column 0. The Timing writer no longer overwrites an indented `## Timing` sample, and `advance` no longer counts an indented heading as lifecycle evidence or as a duplicate section.
```

## Integration seams

None expected. The change is additive, and the only struct construction site is inside `VisibleSections`. REQ-646 touches disjoint files.

## Test wall times

| Run | Wall time |
| --- | --- |
| RED run, three packages, `-run ListNested` | 1.4 s |
| Whole module: `go build ./... && go vet ./... && go test ./... -count=1` | 58 s |
| Probe body: `go test ./internal/requeststate/ ./internal/requestmodel/ -run 'Recovery|VisibleSections|StripGenerated' -count=1` | 0.45 s |
