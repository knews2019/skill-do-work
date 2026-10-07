---
id: UR-137
title: 'Fix the accepted findings from the consumer code review of do-work 0.305.58 to 0.305.69'
created_at: 2026-10-06T22:45:18Z
requests: [REQ-635, REQ-636, REQ-637, REQ-638]
word_count: 1028
---
# Fix the Accepted Findings From the Consumer Code Review of do-work 0.305.58 to 0.305.69

## Summary
A consumer project installed do-work 0.305.69 and reviewed the 0.305.58 to 0.305.69 changes. The review was pasted into `do-work validate-feedback` on 2026-10-06 and every finding was checked against this repository at 0.305.69. All seven accepted findings checked out: F1 and F5 were reproduced with a copy of `VisibleSections`; F3, F4, F7, F8 and F6 were confirmed by reading the cited code. None was already fixed (no commit after 0.305.69). The maintainer then asked to capture the accepted findings as work.

Triage decisions carried into the REQs:
- F1 and F5 are one REQ because both edit the same function in `visible_sections.go`. F1 is `impact-critical` (user text deleted by `recover --take-over`) and carries `priority: now` because the review says "Work in priority order" with F1 first. The same CommonMark bound is applied to fence detection too, a sibling defect found during triage.
- F5's defined behaviour: a line that opens a fence is not scanned for a comment opener (fence before comment, the CommonMark reading). The `## Plan <!-- note` zero-length case does not change.
- F3 and F7 are one REQ because both touch the agent-branch listing and the review says a combined read must still exclude unowned tips.
- F4 and F6 are one REQ: both are about how the board reads activity gaps. F6 is resolved as a relabel, not by including the open gap, to keep one gap meaning across cards and avoid a third definition beside Panel B and the calibration column.
- F8 is its own REQ: a behaviour-preserving deletion in the CLI release guard, `tdd: false`.
- The review's "Do NOT change" items are honoured: the no-phase-stamp exclusion rule stays; only the `durations.go` comment is corrected. Struct field names and `board-controls.js` are untouched.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-635 | [impact-critical] recover --take-over deletes user text under an indented heading; VisibleSections must apply the CommonMark 0-3 space rule and not open a fence on a comment-opening line |
| REQ-636 | Board activity counts only commits the builder branch owns, and each served board response lists worktrees and agent branches once |
| REQ-637 | [impact-rule-change] Board Panel B compares the stamp gap in whole minutes like the calibration log, and the drawer row is relabelled to mean the gap between events |
| REQ-638 | [impact-negligible] Release guard reads suite/modules.tsv once instead of merging two reads of the same declaration |

## Batch Constraints
- Each REQ adds its failing test first (REQ-638 excepted: behaviour-preserving deletion), then fixes, then runs the full Go test suite of every module it touches.
- Each REQ is its own release and commit per `_dev/primes/prime-releases.md`; REQ-635's CHANGELOG entry names F1 as a data-loss fix in `recover --take-over`; REQ-637's entry states why the drawer row was relabelled.
- Lessons entries plus the `do-work/lessons-index.md` token refresh land in the same commit as the fix for REQ-635, REQ-636 and REQ-637.
- Do not change the no-phase-stamp exclusion rule, the single-word private struct fields in `activity_correlation.go`, or `.replace(/^#/, "")` in `board-controls.js`.

## Full Verbatim Input
> `````
> # Fix review findings in do-work suite 0.305.69
> 
> You are working in the do-work suite repository (knews2019/skill-do-work). A consumer project
> installed 0.305.69, and a code review of the 0.305.58 → 0.305.69 changes found the bugs below.
> All of them were confirmed against the code, and F1 and F5 were reproduced. Paths below are as installed
> (`.claude/skills/<module>/...`), so find the matching file in this repo's module layout.
> 
> Work in priority order. Add a failing test first for each bug, then fix it. Run the full Go test
> suite for every module you touch. Follow the repo's own rules for releases, CHANGELOG and lessons.
> 
> ## F1 (high): VisibleSections treats indented `## X` lines as headings
> 
> File: `do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go` (around line 60)
> 
> The 0.305.6x change replaced `strings.CutPrefix(line, "## ")` with
> `strings.CutPrefix(strings.TrimLeft(line, " \t"), "## ")`. That accepts any indentation,
> including 4+ spaces and tabs. In CommonMark that is an indented code block or list
> content, not a heading.
> 
> Impact: `stripGeneratedRecoverySections` (`recover --take-over`) matches names like `Plan`,
> `Review` or `Timing` and DELETES user request text. Example body:
> 
>     # Request
>     Example:
> 
>         ## Plan
>         user sample
> 
>     MUST keep.
>     ## Timing
>     t
> 
> Today this gives a section `Plan` spanning `"    ## Plan ... MUST keep.\n"`, and recovery removes it.
> The same root cause makes `advanceSections` refuse with `duplicate lifecycle section Review`, and lets
> `replaceTimingSection` overwrite an indented `    ## Timing` sample. It also now disagrees with
> `sectionLineBounds` / `appendSectionEntry` in `internal/requeststate/state_apply.go`, which match only
> the exact unindented `## X`.
> 
> Intended behaviour, taken from the new tests in `recovery_markdown_test.go`: accept 0–3 leading SPACES
> only (CommonMark ATX rule). A tab or 4+ spaces is not a heading. Keep the existing tests passing
> (1, 2 and 3 spaces, and an inline `<!-- retained -->` after the heading). Add negative tests for 4 spaces,
> a tab, and the recovery example above (the user text must survive). Note that the lessons file already
> warns about this trap ("`TrimLeft(" \t")` made indented code unreachable"), so add a lesson entry
> saying it happened again.
> 
> ## F5: same-line comment opener plus fence opener hides every later heading
> 
> Same file, lines ~33–56. Removing the old `hadComment` skip means a line that OPENS a comment can
> also open a fence in the same pass. With a line like ```` ``` <!-- ````, `inComment` becomes true and
> the fence opens. Fence lines are skipped without looking for `-->`. After the fence closes,
> `inComment` is still true, so every later heading (including `## Timing`) is hidden, and the Timing
> writer appends a duplicate section.
> 
> Repro body:
> 
>     # Request
>     ``` <!--
>     ## X
>     ```
>     ## Plan
>     p
>     ## Timing
>     t
> 
> Today this returns no sections at all. Fix: a line that opens an unclosed comment must not also
> open a fence. Choose and test a defined behaviour. Do NOT change this: `## Plan <!-- note` with an
> unclosed comment gives a zero-length section by design, because the unclosed region is protected.
> 
> ## F3: liveBranchTipInstants credits a commit the branch does not own
> 
> File: `do-work-board/tools/queue-kanban/activity_correlation.go` (around line 186)
> 
> `git log -1 --format=%cI <worktree-agent-REQ-NNN-*>` returns the branch tip even when the branch
> has no commits of its own. In that case the tip is the integration commit the branch was cut from,
> which may belong to another request. That fake "commit" event splits real idle gaps and can set "last
> activity" on the card. Fix: only count commits unique to the branch (for example
> `git log -1 <integration>..<branch>`, or use the merge-base). Add a test for a branch with no
> own commits.
> 
> ## F4: calibration log and board disagree for gaps of 120–121 minutes
> 
> `do-work/tools/do-work-cli/internal/requestmodel/calibration_row.go` (~52) writes
> `max_stamp_gap_minutes` rounded down to whole minutes, which is the documented format.
> `do-work-board/tools/queue-kanban/durations.go` (`dayMedianExclusionReason`, ~358) compares the
> exact `time.Duration > 2h`. For a 2h00m40s gap the board excludes the request, but a re-fit reading
> the log (`> 120`) keeps it. The docs say both readers must agree exactly. Fix: make the board compare
> whole minutes the same way. Add a boundary test.
> 
> ## F7: duplicate git work on every live board request
> 
> File: `do-work-board/tools/queue-kanban/verify.go` (~202–215, ~1281–1287) and `activity_correlation.go`.
> Each served board response runs:
> - `git worktree list --porcelain` twice (`appendWorktreeFindings` and the disk-space probe). List it once
>   and pass the map to both.
> - the `worktree-agent-*` branch listing twice (verify and `liveBranchTipInstants`), then one `git log -1`
>   per branch. Replace these with one `git for-each-ref --format='%(refname:short) %(committerdate:iso-strict)' refs/heads/worktree-agent-*`.
>   If you do this together with F3, it must still exclude tips the branch does not own.
> Leave the `since` bound in `attachRequestActivity` as it is unless you find a cheap cap that keeps
> stale-claim activity correct.
> 
> ## F8 (cosmetic): suite/modules.tsv is parsed twice
> 
> `do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go` (~37–55) calls
> `DeclaredMaintainerReleaseRoots` and then `DeclaredModuleSources`, then merges them through a set.
> The first already calls the second, and roots are a subset of sources. Simplify it: call
> `DeclaredModuleSources` once, return nil if no source has a tracked `<source>/VERSION`, otherwise use
> the sources plus `suite` and `tools`, sorted. The behaviour must stay identical, so keep the existing
> tests green.
> 
> ## Optional: F6 drawer label
> 
> `activity_correlation.go` (~247–256): for a claimed request, "Largest idle gap" measures only between past
> events and ignores the open gap from the last event to now. The card already shows "last activity Nh
> ago". Either include the open gap or relabel the row so it clearly means between events. Pick one and
> state why in the CHANGELOG.
> 
> ## Do NOT change (reviewed and rejected)
> 
> - No-phase-stamp requests are excluded when their single claimed→completed gap is over 2h. This is the
>   documented rule in `do-work/actions/estimate-reference.md`. You may optionally fix the comment
>   in `durations.go` that says "a long continuous session still counts" so it mentions this case.
> - Single-word private struct fields in `activity_correlation.go`, and `.replace(/^#/, "")` in
>   `board-controls.js`. These are style preferences, not bugs.
> 
> ## Done when
> 
> - Every test above was red before its fix and is green after.
> - The full Go test suites pass for the do-work CLI and queue-kanban modules.
> - A release is cut per this repo's release rules, with a CHANGELOG entry that names F1 as a
>   data-loss fix in `recover --take-over`.
> 
> Paste it into the AI coder's session in the upstream repo. It is self-contained and includes reproduction inputs for F1 and F5.
> `````
