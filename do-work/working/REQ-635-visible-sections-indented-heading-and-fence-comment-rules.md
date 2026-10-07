---
id: REQ-635
title: '[impact-critical] recover --take-over deletes user text under an indented heading; VisibleSections must apply the CommonMark 0-3 space rule and not open a fence on a comment-opening line'
status: claimed
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
claimed_at: 2026-10-07T13:16:10Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: review findings F1 and F5 from the pasted consumer code review of 0.305.58 to 0.305.69, triaged with `do-work validate-feedback` on 2026-10-06 and accepted; the review set "Work in priority order" with F1 first, which is the `priority: now` evidence.*
