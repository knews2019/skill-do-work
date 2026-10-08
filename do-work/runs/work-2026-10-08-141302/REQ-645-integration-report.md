REQ-645 (recover --take-over no longer deletes an indented user heading that shares a generated section's name) is integrated, released as 0.305.81, and archived. UR-141 is closed.

- Archive: do-work/archive/UR-141/REQ-645-recovery-removal-requires-column-zero-heading.md. Finalization also moved REQ-646's archive and UR-141's input.md into do-work/archive/UR-141/.
- Commits: run artifacts c2793dab (pre), merge b5e358f1 (b5e358f1b25c0438e25ce90a52df85b397989d2d), finalization 0aa4a2f7. The record is cleanup_complete with empty blocked_paths and reason_codes.
- Qualification: the qualify gate had no findings. All three new tests failed on the pre-fix source in my own re-run and pass at the merge. The six pinned tests are untouched and green. The four requirement-4 files are in scope.
- Gate: exit 0 on the first run, 121 s. green-gate and the focused probe were satisfied.
- Review: Approve, 98%, Acceptance Pass. No follow-ups were created.
  - F1, impact-user-visible: read-only markdownSectionBytes still matches any indent → report only.
  - F2, impact-negligible: one doc-comment line was not reflowed → report only.
  - The Restatement Sweep found one stale consumer, which is F1.
- Heavy lanes ran in a detached checkout with QUEUE_KANBAN_BROWSER set. All four executed and exited 0, with no skips:

| Lane | Wall seconds |
| --- | --- |
| do-work-cli-integrations | 62 |
| staged-skills | 37 |
| updater | 65 |
| installer | 26 |

- Timing: I recorded handback-merge, the gate plus probe, review, and the heavy drain. All are folded into ## Timing.
- I skipped the builder-work timing event. The hand-back landed at 14:17:56Z and I arrived at 14:28Z. Per work.md, a session that arrives later skips this event, so the builder is not charged with the gap.
- Lesson: I added a rule-direction-checked-against-callers bullet for 0.305.81 to lessons-do-work-cli.md. Its index row went from 18680 to 18987 tokens. contract-regressions.sh passed.
- Cleanup: the builder worktree, its branch (deleted with -d) and both scratch checkouts are gone. No shipped link pointed at REQ-646's old flat archive path.
- Step 0 also staged manifest.md, because it held the coordinator's status edit.
- Left dirty: only this untracked report.
- Discovered tasks: one, report only. It is the same markdownSectionBytes gap as F1.
