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

## Reviewer notes

### Requirements checklist
- [x] R1 heading 0-3 space rule; tab or 4+ spaces is not a heading. Existing 1/2/3-space tests and `## Requirements <!-- retained -->` stay green.
- [x] R2 negative tests (4 spaces, tab, review recovery body), confirmed failing at 4337e9ac.
- [x] R3 same bound in `leadingPunctuationRun`.
- [x] R4 fence opener before comment scan; F5 body yields exactly `Plan`, `Timing`.
- [x] R5 `## Plan <!-- note` still zero-length (new pin test).
- [x] R6 full module run per hand-back (33 packages); reviewer re-ran requestmodel, requeststate, lifecycletiming, lifecycleadvance: all ok; `gofmt -l` empty; `go vet` clean.
- [ ] R7 release + CHANGELOG naming F1 as data-loss fix: not in range, orchestrator finalization work (N/A at review time).
- [ ] R8 lessons entry + lessons-index token refresh: same, finalization work (N/A at review time).
- [x] Constraints: original byte offsets kept for `\n` and `\r\n` (probe D); raw heading name unchanged.

### Finding-Closure Ratchet (F1 and F5 are review findings)
The four new tests were copied onto 4337e9ac and each failed there (`TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule`, `TestVisibleSectionsReadFenceOpenersAndCodeSpansBeforeComments`, `TestRecoveryKeepsUserTextUnderAFourSpaceIndentedHeading`, `TestTimingReplacementIgnoresIndentedSamplesAndFenceInfoComments`); all pass at 2c4dae77. They match the REQ's GREEN statement.

### Orchestrator decision D-01 (`<!--` inside a same-line code span is not a comment opener)
Verdict: correct, and it does not weaken take-over protection.
- A block-level HTML comment (CommonMark HTML block type 2) can only start at the beginning of a line after 0-3 spaces. No code span can come before it there, so `commentOpenerIndex` returns 0 for it, the same as before.
- An inline `<!--` in CommonMark cannot hide a following `## ` line, because an ATX heading interrupts the paragraph. So the only headings D-01 un-hides are headings a renderer shows anyway.
- A span that is unmatched on the line, or that continues onto the next line, stays literal, and the `<!--` after it still opens a comment (probes H and V). This is the conservative direction.
- The `-->` close scan ignores backticks, which is CommonMark-correct inside a comment (existing test plus probe K).
- One divergence from CommonMark: a backslash-escaped backtick is treated as a span opener (probe O). The outcome still matches rendering, for the reason in the second point.

### Adversarial probes (0.305.69 parser vs merged parser)
Changed as intended: F1 LF/CRLF, F5, `\t## Plan`, `  \t## Plan`, double-backtick span with `<!--`, `-->` then a code span with `<!--` on one line, `~~~ <!--` info string.
Unchanged and correct: `<!--` after a closed span opens a comment; fence inside a comment; comment closing and reopening on one line; fence closed by a 3-space closer; `## Plan <!-- note` stays zero-length; comment opener inside an open fence ignored; mismatched backtick run lengths stay literal.
Changed in the hide direction (M1): top-level fence with a 4-space closer (CommonMark-correct), list fence opened at 2 spaces and closed at 4 (CommonMark closes it), nested-list fence at 4 spaces holding `<!--`.
Pre-existing, not changed: `   ##\tPlan` (tab after `##`) is not a heading; a `<!--` in an indented code block opens a comment.

### Callers
- `advanceSections` (`lifecycleadvance/advance_commands.go:344`): same section list type; fewer false `duplicate lifecycle section` refusals for indented samples. Package tests green.
- `stripGeneratedRecoverySections` (`state_apply.go:989`): no longer deletes text under a 4-space or tab-indented generated heading. Recovery tests green.
- `markdownSectionBytes` (`state_apply.go:1012`): now finds sections behind a fence info `<!--` or a code-span `<!--`. No contract change.
- `replaceTimingSection` (`lifecycle_timing.go:515`): appends instead of overwriting an indented sample; replaces instead of duplicating behind a fence info comment.

### Decisions D-03 and D-04
Both are recorded in the hand-back. D-03 is spec-faithful for top-level fences; its list-container risk is M1. D-04 is a test-shape choice. The CRLF `\r\n\n` join it works around is pre-existing and is the builder's discovered task, outside this REQ.

### Suggested testing
1. After finalization, check that the CHANGELOG entry names F1 as a data-loss fix in `recover --take-over`, and that the `lessons-index.md` token count matches `wc -c / 4` for the satellite (R7, R8).
2. Run `advance` on a REQ whose body has a nested-list fenced HTML sample (the M1 shape) to see whether the hide-direction change surfaces in a real REQ.
