---
id: UR-141
title: 'Recovery strips a list-nested generated-name heading, and board activity attributes main commits through an inner merge'
created_at: 2026-10-08T14:08:54Z
requests: [REQ-645, REQ-646]
word_count: 780
---
# Two Accepted Review Findings: Recovery Removal Indent Rule and Board Merge-Ancestry Snapshot

## Summary
The maintainer ran `do-work-toolbox validate-feedback` on two P2 review comments and both were accepted, then asked to capture and run them. Finding 1: `recover --take-over` deletes a user-authored `  ## Plan` sample nested under a bullet item, because removal is keyed on the section name alone while 1-3 space indented headings count as sections (on purpose since 0.305.59). The triage replaced the reviewer's remedy (track list context) with a narrower one: only a column-0 heading is removable. Finding 2: the Kanban activity correlation reads the attribution map inside its ancestry loop, so a merge of main nested inside a tagged hand-back merge inherits the REQ id and attributes main's commits to that REQ, against the prime's commit-evidence rule.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-645 | [impact-critical] recover --take-over keeps a 1-3 space indented generated-name heading; only a column-0 section is removable |
| REQ-646 | Board activity snapshots direct merge matches before expanding ancestry, so an inner merge of main attributes nothing |

## Batch Constraints
- No dependency between the two REQs: disjoint modules (`do-work-cli` and `queue-kanban`), disjoint files, may run in parallel.
- Each REQ is its own release per `_dev/primes/prime-releases.md`; do not fold one REQ's change into the other.
- Both fixes are direct (no new guard, fallback, or validation layer): Surface-cost N/A per the triage. The reviewer's list-context remedy for finding 1 was flagged and is not the remedy.
- Read `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` before REQ-645; `_dev/primes/prime-kanban-board.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` before REQ-646.

## Full Verbatim Input
> ```
> capture them and run them
> 
> [Capture source: the two findings accepted by `do-work-toolbox validate-feedback` in this session, pasted by the maintainer as "do-work validate-feedback: is the following important enough to be fixed?", then the triage verdicts. Original feedback verbatim first, then the accepted-finding blocks with Evidence and Surface-cost.]
> 
> --- Original feedback (verbatim) ---
> 
> Full review comments:
> 
> - [P2] Exclude list-nested headings from removable request sections — .claude/skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go:64-68
>   When a requirement contains a bullet item with a two-space-indented `## Plan`, this now treats that nested heading as an independent request section. Consequently, `stripGeneratedRecoverySections` deletes the user-authored example during `recover --take-over`. A targeted regression test reproduces the deletion. Track list/container context before recognizing removable sections; indentation alone cannot distinguish these from request-level headings.
> - [P2] Freeze direct matches before expanding merge ancestry — .claude/skills/do-work-board/tools/queue-kanban/activity_correlation.go:142-150
>   If a builder merges main before its tagged hand-back merge, processing the outer merge attributes its request ID to that inner merge. This loop then reads the mutated map as `directIds` and attributes unrelated main-branch commits through the inner merge, falsely shortening the reported activity gap. A targeted regression test reproduces this. Snapshot the direct matches before expanding ancestry to preserve the request-owned commit evidence rule (.claude/skills/do-work-board/tools/queue-kanban/prime-do-kanban.md#L30).
> 
> --- Triage (validate-feedback, 2026-10-08) ---
> 
> Finding 1: list-nested `  ## Plan` is stripped by `recover --take-over` · P2
> - Verdict: Accept, with a different remedy than the reviewer proposes.
> - Evidence: `skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go:64-68` admits a heading after 0-3 spaces, and `requeststate/state_apply.go:981-999` strips any section whose name matches the generated list, with no indent check. The REQ-635 test table itself asserts `" ## Plan"` and `"   ## Plan"` yield a `Plan` section, so the deletion needs no reproduction run.
> - Reasoning: The claim is correct, but "this now treats" is wrong on history. One-to-three-space headings have counted as sections since 0.305.59 (`96505655`), on purpose: `TestRecoveryPreservesIndentedAndCommentedRequirements` pins that an indented `  ## Requirements` is a boundary so user bytes under it survive. REQ-635 narrowed the rule to CommonMark's 0-3 spaces and explicitly kept those cases green (its requirement 1), so the 2-space list case is a residual of the same incident (REQ-635 F1), not a regression. The REQ-635 record's "Why" already notes the other writers (`sectionLineBounds`, `appendSectionEntry`) match only the exact unindented `## X`. That asymmetry is the fix: generated sections are always written at column 0, so removal should require column 0 while boundary recognition keeps the 0-3 space rule. That keeps the 0.305.59 tests green and preserves the reviewer's example without the parser learning list structure.
> - Surface-cost: Flagged for the proposed remedy, N/A for the alternative. "Track list/container context" grows the Markdown grammar (list markers, continuation indent, nesting), the exact trap the REQ-634 lesson warns about. Requiring column 0 for the removable set is a narrowing, not a new layer.
> - Remedy: In `stripGeneratedRecoverySections`, skip any section whose heading line does not start at column 0 (or have `VisibleSection` carry its indent). Lock-in test: a bullet item with `  ## Plan` under it survives, with the existing 1-3 space boundary tests unchanged.
> 
> Finding 2: inner merge inherits a request id and attributes main's commits · P2
> - Verdict: Accept.
> - Evidence: `skills/do-work-board/tools/queue-kanban/activity_correlation.go:142-150` reads `requestIdsByHash[commit.hash]` inside the loop after earlier iterations have called `attribute()` on range hashes. The comment at lines 139-141 says ancestry runs "off the DIRECT matches only, collected before any range is attributed". The code does the opposite, and git log order (newest first) processes the outer hand-back merge before the inner merge, so the inner merge is already tagged when its turn comes. Its second-parent range (main's commits since the branch point) then gets the REQ id. The prime rule at `prime-do-kanban.md:30` ("commits that REQ's work owns") is violated.
> - Reasoning: The existing test (`activity_correlation_test.go:145`) covers one matched merge and one unmatched merge, never a merge nested inside a matched range, so nothing pins the stated rule. On this repo the trigger is rare: the suite's builders are never told to merge main into their branch, and the last 60 merges show none. In consumer repos where humans or other agents merge main before hand-back it fires routinely, and the drawer's "Largest gap" and a claimed card's last activity are the only consumers, so the damage is a wrong display number. P2 holds for the shipped board; nothing writes state from it.
> - Surface-cost: N/A. Direct fix.
> - Remedy: Collect the merges with direct ids (or copy the hash→ids map) before the ancestry loop, then attribute from that snapshot. Lock-in test: outer `[REQ-N]` merge whose second parent contains an inner merge of main; the main-side commits must stay unattributed.
> 
> Summary: Accept 2 (#1 with column-0 removal, not list tracking; #2). Already done 0. Push back 0. Discuss 0.
> ```