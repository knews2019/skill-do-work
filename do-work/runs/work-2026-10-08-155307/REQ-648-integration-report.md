# REQ-648 integration report (operator-gated work is blocked, never earmarked)

REQ-648 is released as 0.305.82 and archived. Finalization reached cleanup_complete with no blocked paths and no reason codes.

- **Archive path:** do-work/archive/REQ-648-schema-and-work-action-state-operator-rule.md (flat, because UR-142 still has REQ-647 and REQ-649 open).
- **Commits:** run artifacts 336cbb52 (this is `<pre>`). Merge 4c50ccde (4c50ccdecf64fc63ec429ca252768252e32faf83). Finalization 3aa5f079 (3aa5f0798f34297be084166d0b8e646b480debfd).
- **Gates:** the repository gate exited 0 on the first run, in 115s. The qualify gate, the green-gate record and the focused probe were all satisfied.
- **Review:** Approve, 98%, Acceptance Pass. Two findings, both impact-negligible and report only. F1: `work-reference.md:788` still cites a "work.md mid-run blocked-flip procedure" that does not exist, and work.md now points back at that row (same as the builder's discovered task). F2: the new row says the default run "hides" an earmarked REQ, while other text says "skips and reports". The maintainer accepted that wording. The restatement check found the capture-side text agrees. I applied no fixes.
- **Heavy lanes:** staged-skills ran at 4c50ccde in a detached checkout. Exit 0, executed, 43s wall. The first attempt went red in 0s because the duration log header I wrote was wrong. I copied ROOT's header line and reran it. That was an environment problem, not a code failure.
- **Timing events recorded:** builder-work, handback-merge, verification-gate twice (gate and heavy drain), and review. They are folded into the REQ's Timing section.
- **Lesson:** a new alternate-writer-contract-drift bullet in _dev/primes/lessons-action-files.md links to the flat archive path. The lessons-index row now says 6716 tokens. contract-regressions.sh exited 0 after finalization.
- **Cleanup:** I removed the builder worktree and deleted its branch with -d. I removed the drain checkout and pruned worktrees.
- **Left dirty under ROOT:** the REQ-647 and REQ-649 hand-back and integrate files, all untracked. They belong to the sibling integrators, so I left them alone. This report is also untracked.
- **Report-only discovered tasks:** (1) The loop in F1 above. A later REQ could drop the "see work.md" clause at work-reference.md:788. (2) Line 21 of _dev/primes/prime-action-files.md names a `pipeline.md` action, and that file no longer exists in any package. This predates REQ-648.
