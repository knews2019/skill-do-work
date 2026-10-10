# REQ-655 integration report

- Release 0.305.105, "ai-report Lists Every Report Bundle in One Catalog and Finds Reports on a Topic". Archive `do-work/archive/REQ-655-ai-report-index-and-find-catalog.md` (flat, UR-144 stays open). Finalization commit `424cac54`: cleanup_complete, no blocked paths or reason codes, first try.
- Range: pre 355ea12f. First merge 6a031809. Review-fix re-merge d91312ae (recorded as `commit:`). Integrator fix 851fa107 on the builder branch. No conflicts.
- Qualification: qualify and scope-drift satisfied. One QUALIFY-NEW-FILE-UNWIRED warning on the `_test.go` file, judged a false positive. Touched set = the 8 declared files. D-11 (`--format text`) and D-12 (explicit id boundary) confirmed in the code.
- Review 91%, Approve. I fixed F1: `find` matched the shared `ai-reports/` prefix, so `find ai` returned all 41 bundles. It now matches the folder name, and a new assertion pins this (D-23). Delta re-review 92%, F1 closed. Report only: F2 (stale neighbouring lines at `ai-report.md:30`, `:135` and `ai-report-guide.md:21`; REQ-657/REQ-687 edit these next), F3-F10 (small test gaps, parser nits, anti-bloat count).
- Gate (`DO_WORK_FAST_STAGE_REUSE=off`) exit 0 at 6a031809 (166 s, load 3.07 before) and at d91312ae (160 s, load 3.85 before). Probe exit 0 both times. Green record at d91312ae via `record-green-gate`.
- Heavy lanes at d91312ae, all exit 0 and executed: do-work-cli-integrations 79 s, staged-skills 38 s, updater 68 s, installer 32 s. The earlier drain at 6a031809 was also green.
- Timing events: handback-merge 2, verification-gate 4, review 1, folded into `## Timing`. Builder-work skipped (the hand-back had already landed).
- Lessons: one `regexp-word-boundary-underscore` bullet in `lessons-do-work-cli.md`, linked to the flat archive path. The index row is now 20142 tokens with the new family. `contract-regressions.sh` passes after finalization.
- For the coordinator: when REQ-656 closes UR-144, re-point that bullet to `archive/UR-144/REQ-655-...`.
- Cleanup: builder worktree removed, branch deleted with `-d` and pruned.
- Left dirty: the siblings' untracked `REQ-*-handback.md` files (not mine) and this report (untracked, for the coordinator).
- Report-only discovered tasks (builder): other toolbox verbs that return `exactOutputResult` may be called with `--format json` in their prose and print nothing. The loose `ai-reports/2026-05-28_2335_background-agents-durability.html` is invisible to the catalog.
