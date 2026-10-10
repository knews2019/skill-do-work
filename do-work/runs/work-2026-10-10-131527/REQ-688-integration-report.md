# REQ-688 integration report

- Release 0.305.108, "capture-files Prints a Fill-In Example Manifest With Payload Templates, and the Capture Reference Shows the Fence Capture Accepts". Archive `do-work/archive/REQ-688-capture-files-example-and-fence-fix.md` (flat, UR-153 stays open). Finalization commit `196690b1`, cleanup_complete on the first try.
- Range: pre `008cc74c`. First merge `125c62c2`. Re-merge `f42468da` after the review fix (recorded as `commit:`). My fix on the builder branch: `2213664b`. No conflicts.
- Qualification: qualify and scope-drift satisfied. The QUALIFY-NEW-FILE-UNWIRED warning is a false positive: the new Go file is called from `publication_commands.go`. 5 declared files, 5 touched. I fixed the stale REQ-661 path in the REQ file.
- Review: 94%, Approve. F1 is fixed with the reviewer's exact text. A raw input that contains `UR-NNN` or `REQ-NNN` broke the verbatim block on a global replace, so the header now says to leave the `> ` lines as printed. The delta is that text only, so I did not spawn a second reviewer. Report only: F2 (`clarify.md:106` fence rule has no three-backtick minimum and no no-info-string rule), F3 (`capture.md:136,138` addendum example still has a `text` fence) and the F4-F7 nits.
- Gate: exit 0 at `125c62c2` (125 s, load 2.82). At `f42468da` the first run exited 1 on an unrelated `nextselection` flake ("could not launch: no such process"). The rerun exited 0 (129 s, load 4.51). The probe was GREEN at both merges. Green record at `f42468da`.
- Heavy lanes at `f42468da`, all executed with exit 0: do-work-cli-integrations 64 s, staged-skills 34 s, updater 65 s, installer 26 s.
- Timing: 7 events (merge 2, gate and drain 4, review 1), folded. I skipped the builder-work event because the hand-back had already landed.
- Lessons: one `restated-mechanism-unchecked` bullet in `_dev/primes/lessons-action-files.md`, linked to the flat archive path. REQ-689 re-points it. Index row: 7871 tokens. `contract-regressions.sh` passes.
- Cleanup: builder worktree removed, branch deleted with `-d`, worktrees pruned.
- Left dirty: only this report (untracked).
- Report-only builder task: `capture.md:316` and `:271` spell the capture commit format two different ways.
