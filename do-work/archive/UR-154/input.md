---
id: UR-154
title: 'Review follow-ups R1 to R5 from run work-2026-10-10-131527'
created_at: 2026-10-10T19:19:53Z
requests: [REQ-693, REQ-694, REQ-695, REQ-696]
word_count: 254
---
## Summary

Five report-only review findings from the twelve-REQ coordinated run, promoted to queue work by the maintainer. Each finding was re-checked against the code and prose at HEAD `0e8e0ef9` before capture. All five are still true, so none was dropped. They are folded by owner into four REQs, one coherent fix each.

## Extracted Requests

| REQ | Finding | Owner | Re-check at HEAD |
| --- | --- | --- | --- |
| REQ-693 | R1: REQ-660 review F4 and F3 (also REQ-658 F7) | `do-work-cli worktree` and the worktree dirt readers | Still true. `reportStatus` and `removeBuilderWorktree` run `git status` inside a missing folder and stop with `WORKTREE-GIT-FAILED`. `worktreeClean` (`cleanup_git.go:345`), `handoff.go` and queue-kanban `worktreeHasUncommittedWork` (`verify.go:1109`) still read the whole porcelain status with no `do-work/worktree-links` filter. |
| REQ-694 | R2 + R3: REQ-659 re-review N1 and review F3 | `frontmatter set` and `req append-section` (`request_writers.go`) | Still true. The post-insert check counts only visible copies of the new section, not the sections after it. `handleFrontmatterSet` has no guard for `status` or `id`. |
| REQ-695 | R4: REQ-690 review F3 and F4 | `do-work-cli run-status` (`run_status.go`) | Still true. `resolveRunDirectory` does not stat `--run` (D-19), and the text never names the run directory. The C3 arm returns `do-work run` with no check of `last_activity_phase`. |
| REQ-696 | R5: REQ-690 F9, REQ-688 F2 and F3, REQ-692 M2 | Shipped prose plus one contract-test message | Still true at every site. Line numbers moved: `forensics.md:10`, `README.md:176`, `clarify.md:106`, `capture.md:136,138` (unchanged), `staged-skills-contract.sh:792` (unchanged), fixture header line 1. |

## Folded Requests

None. The queue was empty at capture time, so no fold-first candidate existed.

## Batch Constraints

- YAGNI: each REQ names the failure, the smallest fix, and the one test that pins it. A prose-only part pins nothing new.
- No `depends_on` edges: the four REQs touch disjoint lines.
- R4 keeps the D-19 coordinator ruling (no `--run` refusal). The fix is display only.
- Capture decision (unattended, no question asked): R1 to R4 carry `addendum_to` pointing at the archived REQ whose review raised them, because they correct shipped work. R5 spans three archived REQs, so it carries no `addendum_to` and names its sources in the body. The commit uses the plain multi-REQ capture message, because this capture is not one single addendum.

## Full Verbatim Input
> ```
> Report-only follow-ups from the reviews of run work-2026-10-10-131527 (UR-144, UR-145, UR-153), queued at the maintainer's "keep going, don't wait on me" after the run summary listed them as R1 to R5:
> 
> - R1: REQ-660 (worktree new/status/merge/cleanup): `status` and `cleanup` exit 2 when a builder worktree folder is gone (review F4). Cleanup Pass 5, the board's verify check, the handoff survey and crash recovery read the worktree command's own links as uncommitted work and ask for consent (review F3; also REQ-658 F7). Source: do-work/archive/UR-145/REQ-660-worktree-lifecycle-command.md ## Review, do-work/runs/work-2026-10-10-131527/REQ-660-integration-report.md.
> - R2: REQ-659: `req append-section` accepts a body that leaves a code fence or comment open, which hides every later section from `advance` (review N1; suggested fix: refuse unless the visible section count grows by exactly one). Source: do-work/archive/UR-145/REQ-659-frontmatter-set-and-req-append-section.md ## Review.
> - R3: REQ-659: `frontmatter set` accepts any value for `status` (review F3, impact-rule-change). Same source.
> - R4: REQ-690: a mistyped `--run` silently turns C3 (hand-back landed) rows into C7 (claim past threshold) rows (review F3); C3 recommends `do-work run` even while an integrator runs (F4). Source: do-work/archive/UR-153/REQ-690-run-status-action.md ## Review, do-work/runs/work-2026-10-10-131527/REQ-690-integration-report.md.
> - R5: stale wording the wave-end sweep left: `forensics.md:10` and `README.md:176` still send "stuck" questions to forensics (REQ-690 F9); `clarify.md:106` states the fence rule without the three-backtick minimum or the no-info-string rule (REQ-688 F2); the addendum example at `capture.md:136,138` uses a four-backtick `text` fence (REQ-688 F3); the retired-trigger fixture header and the message at `_dev/tests/staged-skills-contract.sh:792` still state that core does not route sibling actions (REQ-692 M2). Sources: the REQ-688, REQ-690, REQ-692 archived REQs under do-work/archive/UR-153/ and do-work/runs/work-2026-10-10-131527/REQ-689-integration-report.md.
> ```
