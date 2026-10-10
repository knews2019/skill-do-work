# REQ-686 integration report

REQ-686 (short doc and comment lines for the pushed-back upstream findings) is integrated and finalized. UR-152 is closed.

- Release: 0.305.101, "Lessons Record How to Clear a Stuck Finalization Journal After a Revert, and Comments Name Their Real Readers".
- Archive: do-work/archive/UR-152/REQ-686-pushed-back-finding-doc-lines.md. The seven siblings (REQ-679 to REQ-685) and UR-152 input.md moved into do-work/archive/UR-152/ in the same commit.
- Commits: pre 577a4ac7 (run artifacts), merge 93198e39 (full 93198e39c8a4917c21292350e12befffc0338faf), finalization 1fc1e3e6. Builder branch and worktree removed.
- Merge seams (comment only, both coordinator addenda): git_transaction.go HEAD-sentinel comment now names cleanup_apply.go, doctor_repair.go, publication_commands.go (all three check CommitSHA != ""); two stale comments fixed in timeline_scroll_browser_probe_test.go.
- Gate: DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh exit 0 at the merge, wall 141 s, load before 3.03, one run. Probe exit 0. contract-regressions.sh exit 0 after finalization.
- Review: Approve, 97%. Wave-end Restatement Sweep over items a to f: one stale statement, F1. Minor, all report only: F1 lessons-do-work-cli.md:121 still calls --discard-journal "the only supported exit" next to the new hand-delete sentence (reword to "for a journal still in phase prepared"); F2 anti-bloat clean, the gittransaction comment edit is coordinator-ordered; F3 the HEAD comment's caller list is not exhaustive (transaction_findings.go:171, finalization_apply.go:142 also read it); F4 the .git/do-work-finalization/ path is main-checkout only, and REQ-N sits beside REQ-NNN; F5 the heavy_commands.go comment becomes the Handlers() doc comment.
- Heavy lanes (detached checkout, Chrome set, none skipped): queue-kanban-javascript 7 s, queue-kanban-browser 74 s, do-work-cli-integrations 62 s, staged-skills 36 s, updater 65 s, installer 29 s. All exit 0, executed.
- Timing events: builder-work, handback-merge, verification-gate (gate), review, verification-gate (heavy drain); Timing section folded.
- Lessons: builder proposed none. Re-pointed the REQ-679 and REQ-682 bullets in _dev/primes/lessons-kanban-board.md to archive/UR-152/; lessons-index rows refreshed (kanban-board 6122, do-work-cli 19046).
- Link check before the move: git grep for archive/REQ-679..685 outside do-work/ found only those two bullets.
- No new refusals from the merged code; finalization succeeded first try.
- Left dirty under ROOT: only untracked run-directory files of other REQs plus this report (coordinator commits at checkpoint).
- Discovered tasks: none. Report-only follow-up: review F1 wording.
