# REQ-704 integration report (second opinions build a linked iteration tree on any base page)

REQ-704 is released as **0.306.8** "Second Opinions Build a Linked Iteration Tree". UR-160 (directed second opinions) is closed.

- Archive: `do-work/archive/UR-160/REQ-704-second-opinion-linked-iteration-tree.md`. REQ-703, REQ-705 and the UR `input.md` moved to `do-work/archive/UR-160/` in the same commit.
- Merge: the first merge was `92443689`. The re-merge after the review fix is `30666c9d`, with the same pre `87614289`. Range `87614289..30666c9d`. Finalization commit: `6a53259e`. Worktree and branch removed (`git branch -d`). `contract-regressions.sh` passed after finalization.
- Gate: run 1 at `c5a8aa8d` (same code as `92443689`) exit 0, 146 s, load 4.99. Run 2 at `30666c9d` exit 0, 142 s, load 4.94. Slowest file: `report_second_opinion_test.go` at 26.49 s (limit 30 s). The probe passed. Green-gate recorded at `30666c9d`.
- Review: 95%, Approve, Route C, with the wave-end Restatement Sweep over REQ-703 and REQ-705.
  - F1 (the "plain `ls` order" claim fails on Linux locales) and F2 (sentences true only for plain runs) were fixed in prose on the builder branch (`cd29e7f3`): 4 sentences in `ai-report.md` and `ai-report-guide.md`.
  - F3 (the REQ-703 lesson says the meta strip misses reordered copies) was handled in the link sweep: the bullet now says "REQ-704 closed that gap".
  - F4-F7 are report only: missing tests for the widened strip and the old-template refusal, a possible id race between simultaneous publishes, duplicate findings, and `cmp.Or` written four times.
  - Anti-bloat: 8 helpers the plan did not name. The reviewer kept all 8.
- Heavy drain at `30666c9d`, all 6 lanes. The re-plan was uncertain because the range now holds `do-work/` run-artifact commits, so every lane was selected.
  - Green: queue-kanban-javascript 9 s, queue-kanban-browser 101 s (a real run, not skipped), staged-skills 45 s, updater 75 s, installer 31 s.
  - do-work-cli-integrations failed on run 1 (80 s, load above 5). The failing test was `TestCancelledCommitThatLandsReportsCommittedRisk`, which is outside REQ-704's files: "no process id recorded at hook.pid". The lane alone passed on rerun at load 4.19 (97 s).
- Timing events recorded: builder-work, handback-merge, verification-gate (gate), review, verification-gate (heavy drain). The review event (12 min) also includes the fix and gate run 2, because `record-timing-event` has no end-time flag. `## Timing` folded.
- Lesson: a REQ-704 bullet in `lessons-do-work-cli.md` (family writer-anchored-marker-position). The REQ-703 link was re-pointed to `archive/UR-160/`. The index row is refreshed to 23054 tokens.
- Left dirty: none. `git status` is clean. This report is untracked.
- Report-only discovered tasks:
  - `reportMetaAttribute` matches `data-name` (negligible).
  - `TestCancelledCommitThatLandsReportsCommittedRisk` can fail under load (its PID wait is timing-bound).
  - Review F4 and F5.
- Next: the coordinator runs `advance --checkpoint`.
