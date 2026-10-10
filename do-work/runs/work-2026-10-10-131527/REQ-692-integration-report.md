# REQ-692 integration report

- Release 0.305.109, "validate-feedback Can Capture, Verify and Run the Accepted Findings, and Core do-work Routes to It". Archive `do-work/archive/REQ-692-validate-feedback-capture-run.md` (flat, UR-153 stays open). Finalization commit `08759290`, cleanup_complete on the first try.
- Range: pre `8f4bda05` (run-artifacts commit). First merge `c07fcb13`. Re-merge `b5d65c91` after the review fixes (recorded as `commit:`). My fix on the builder branch: `46b4a62e`. No conflicts.
- Qualification: qualify and scope-drift satisfied. The first scope-drift call refused because the pre-dispatch `## Scope` listed paths without backticks; I added backticks only.
- Review: 90%, Approve. I fixed I1 (a "capture and run" phrase inside a pasted review could switch on write mode), I2 (capture-reference.md Fold-First destination 3 and 4 did not name the validate-feedback hand-off) and M1 (one combined report), each with the reviewer's exact text. I2 touched `skills/do-work/actions/capture-reference.md`, outside the write set; I added it to `## Scope` with the reason. The delta is the reviewer's text only, so I did not spawn a second reviewer. Report only: M2 (retired-trigger fixture header and staged-skills-contract.sh:792 message still state the old "core routes no sibling action" rule), M3 (small unnamed additions), M4 (Step 6 "print no line" wording), N1-N4.
- Gate: exit 0 at `c07fcb13` (124 s, load 3.28) and exit 0 at `b5d65c91` (119 s, load 2.99). GREEN probe passed at both. Green record at `b5d65c91` via `record-green-gate`.
- Heavy lanes: staged-skills only, executed exit 0: 47 s at `c07fcb13`, 33 s at `b5d65c91`.
- Timing: 7 events (merge 2, gate and drain 4, review 1), folded. Builder-work event skipped because the hand-back had landed long before.
- Lessons: one `alternate-writer-contract-drift` bullet in `_dev/primes/lessons-action-files.md`, linked to the flat archive path; REQ-689 re-points it. Index row: 8028 tokens. `contract-regressions.sh` passes.
- Cleanup: builder worktree removed, branch deleted with `-d`, worktrees pruned. The detached drain checkouts were removed with `--force` because of the seeded `test-durations.tsv`.
- Left dirty: only this report (untracked).
- Report-only builder tasks: `journey-qa.md:111` and the maintainability-audit loop footer describe only the no-flag mode; core `help.md` does not mention the new forward.
