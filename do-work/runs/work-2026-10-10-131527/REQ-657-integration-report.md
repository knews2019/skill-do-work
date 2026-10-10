# REQ-657 integration report

- Release 0.305.106 ("ai-report Render-Checks a Report Bundle With One Command at Phone and Wide Widths in Light and Dark"). Archive `do-work/archive/REQ-657-ai-report-judge-render-check-command.md` (flat; UR-144 stays open). Finalization commit `bbc03236`, cleanup_complete on the first try.
- Range: pre `aea1c4d3`, first merge `3f8fd77c`, re-merge `fe0b19e2` (the recorded `commit:`). My fixes on the builder branch: `4021963c` and `1ae7a03b`.
- Five files conflicted with REQ-655. I kept both texts and set the handler count to 9. I also fixed the `ai-report.md` `$ARGUMENTS` line (REQ-655 review F2).
- The heavy drain at `3f8fd77c` failed on staged-skills: "route ai-report exactly once (found 2)". I fixed it by moving the judge phrases into the single ai-report routing row.
- Review: 80%, Approve. I fixed F1 (architecture-report could not reach `pass` because the `../` link to the prior bundle always returns 404), F2 (an exit-2 error is now disclosed like a skip) and F3 (the command-line guide names both verbs). Delta re-review: 91%, Pass. Report only: F4-F7, the code part of F8, and F9.
- Gate: passed at `3f8fd77c` (150 s). At `fe0b19e2` the first run failed because the disk was full (131 MB free). I trimmed Go build-cache entries unused for more than 1 day (about 2.3 GB). The rerun passed (170 s) and the probe is GREEN with no skips.
- Heavy lanes at `fe0b19e2`: staged-skills ran (40 s). The other three reused their exit-0 runs from `3f8fd77c` (79/72/35 s).
- Timing: 7 events folded into the REQ. I skipped the builder-work event because the hand-back had already landed.
- Lesson: one bullet in `lessons-do-work-cli.md`, linked to the flat archive path. `contract-regressions.sh` passes.
- Coordinator: when REQ-656 closes UR-144, re-point this bullet. Disk space is low (about 3 GB free).
- Builder worktree removed; branch deleted with `-d`.
- Left untracked: `REQ-656-handback.md` (the REQ-656 builder's file) and this report.
- Report-only discovered tasks: `command-line-guide.md` mixes recipe and recipe-less verbs, and `architecture-report.md:79` does not name the draft-folder rule.
