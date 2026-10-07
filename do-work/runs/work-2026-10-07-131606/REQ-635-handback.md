# Hand-back — REQ-635 (VisibleSections: 0-3 space indent rule, fence before comment, no comment opener inside a code span)

- Branch: worktree-agent-REQ-635-visible-sections-indented-heading-and-fence-comment-rules
- Commit: 9f8d55e6 (one commit on top of 13932fc2; not merged, not pushed)

## File manifest (all modified, none new)
- skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go
- skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections_test.go
- skills/do-work/tools/do-work-cli/internal/requeststate/recovery_markdown_test.go
- skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing_test.go

## Red-Green evidence (every case runs over both `\n` and `\r\n`)
- TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule. RED: F1 body `got sections ["Plan" "Timing"], want ["Timing"]`; `\t## Plan`, ` \t## Plan`, 5 spaces: `got ["Plan"], want []`; `    ```` ` then `## Plan`: `got [], want ["Plan"]`; 4-space closer inside a fence: `got ["Hidden"], want ["Plan"]`. GREEN: pass.
- TestVisibleSectionsReadFenceOpenersAndCodeSpansBeforeComments. RED: F5 body `got [], want ["Plan" "Timing"]`; `` Use `<!--` here. `` and the double-backtick span: `got [], want ["Plan"]`. The unmatched-backtick case (hidden), the "comment after a code span on the same line" case and the "backticks inside a comment close" case were already green (they pin unchanged behaviour). GREEN: pass.
- TestVisibleSectionsKeepAnUnclosedCommentHeadingZeroLength: new pin, green before and after (no test pinned `## Plan <!-- note` before).
- TestRecoveryKeepsUserTextUnderAFourSpaceIndentedHeading. RED: `got "---\nid: REQ-501\n---\n# Request\nExample:\n\n"` (the indented sample and `MUST keep.` were deleted). GREEN: bytes kept exactly.
- TestTimingReplacementIgnoresIndentedSamplesAndFenceInfoComments. RED: `indented Timing sample was replaced` (sample overwritten) and `Timing section behind a fence info comment was not replaced` (duplicate `## Timing` appended). GREEN: pass.
- Existing 1/2/3-space and `## Requirements <!-- retained -->` tests stay green.

## P-A-U
- [x] **[PLAN]:** One helper `blockIndentContent` (0-3 spaces; a tab in the indent or a 4th space means indented code) used by the heading check and by `leadingPunctuationRun`. Move the fence-opener check before the comment scan for lines that start outside a comment. Replace the `<!--` search in the opener scan with `commentOpenerIndex`, which skips CommonMark inline code spans on the same text. The `-->` close scan is unchanged.
- [x] **[APPLY]:** Done as planned in visible_sections.go only, plus the four test files' new tests. Doc comments on `VisibleSections`, `leadingPunctuationRun`, and the two new helpers state the conditions.
- [x] **[UNIFY]:** `git diff --stat HEAD~1`: 4 files, 187 insertions, 10 deletions. `gofmt -l .` empty. `go vet ./...` clean. Full module `go test -count=1 ./...`: 33 packages, 0 failures, wall 94s. The three touched packages: requestmodel 0.9s, requeststate 8.4s, lifecycletiming 4.3s. Untouched packages over 30s under full-module parallel load: finalization 90s, heavyverification 57s, lifecycleadvance 43s, publication 40s (pre-existing, not changed here). No debug artifacts in the diff.
- `advance REQ-635` (read-only, worktree CLI against the main tree) now reports `"phase": "estimate-p50"` (mechanical), past `agent judgment: triage and open questions`. Main tree `git status` was the same before and after.

## Decisions
- D-03 (DECIDE & STATE): a fence closer indented 4+ spaces no longer closes an open fence. This follows CommonMark (a closer allows 0-3 spaces) and the brief's "opener or closer" rule. Value: an indented code sample inside a fence cannot end the fence early. Risk: a fence inside a list item, indented 4+ spaces, now opens nothing; its contents are then only hidden if they are themselves 4+ indented, which they are in such a list. Easy to revert (one helper).
- D-04 (DECIDE & STATE): the timing append test for the indented sample checks "sample kept as a prefix, one real `## Timing` appended at the end" instead of exact bytes. The appender joins with `\n\n` after trimming only `\n`, so on CRLF input it writes `\r\n\n`. That join is outside this REQ's write boundary.

## Discovered Tasks
- `replaceTimingSection` (lifecycletiming/lifecycle_timing.go) joins with `\n` on CRLF documents (`...\r\n\n## Timing\n...`), producing mixed line endings. Low impact; candidate for a small follow-up if CRLF REQ files matter.

## Lessons read
- lessons-do-work-cli.md line 66 (REQ-460, closed-enumeration-for-a-condition: a trim before the test destroys the structure). This REQ is that trap again; the new helper measures the indent before any trim, and its doc comment says why.
- lessons-do-work-cli.md line 27 (lifecycle-section-evidence: hidden text and quoted text are separate classes). Code spans are quoted text, so a comment opener inside one is not hidden text.
- Line 24 (`#` excluded from enclosing fences): unchanged.

## Integration seams
None. Callers (`stripGeneratedRecoverySections`, `replaceTimingSection`, `advanceSections`) consume the same section list type.
