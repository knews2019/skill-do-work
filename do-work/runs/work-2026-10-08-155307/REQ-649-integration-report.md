REQ-649 (the board's Earmarked badge tooltip now names Needs input · Blocked) is finalized as release 0.305.84. It closed UR-142.

- **Archive:** do-work/archive/UR-142/REQ-649-board-earmarked-tooltip-names-blocked-column.md. REQ-647, REQ-648 and UR-142's input.md moved into the same folder.
- **Commits:** run artifacts ce5cfbe7 (the pre-merge tip), merge 5eb615be (full 5eb615beeff876fb99353eba441469b14304a2f5), finalization 5bbfd367.
- **Qualification:** the qualify gate was satisfied with no findings. The diff is one file, exactly the write_set.
- **Gates:** maintainer-verify.sh exited 0 in 118 s on the first run. The probe ran through advance and green-gate was satisfied.
- **Review:** Approve, 100%, acceptance Pass. F1: board-guide.md:18 does not mention the operator case. The file is outside the wave, impact-negligible → report only. F2: a blocked REQ with an unmet depends_on shows under Pending → Waiting. The wording was accepted as quoted, impact-negligible → report only. No re-merge was needed.
- **Restatement Sweep:** this REQ redefines nothing. These files agree: work-reference.md:114, capture.md:107, work.md:523, work-guide.md:119, clarify.md:192 and board-cards.js:240. board.md, prime-do-kanban.md, model.go and board.css describe placement only and do not conflict. The only gap is F1.
- **Heavy lanes:** run at 5eb615be from a detached checkout with QUEUE_KANBAN_BROWSER set. All three executed and exited 0: queue-kanban-javascript 7 s, queue-kanban-browser 72 s, staged-skills 33 s.
- **Lessons:** no new bullet. I moved the REQ-648 and REQ-647 bullets in lessons-action-files.md to sit right after the REQ-644 bullet and pointed their links at archive/UR-142/. The lessons-index token count went from 6852 to 6855. contract-regressions.sh passed after finalization and again after cleanup.
- **Timing events:**
  - builder-work: 1279 s. It includes the wait behind the two sibling integrations.
  - handback-merge: 18 s.
  - repository gate: 131 s.
  - heavy drain: 128 s.
  - review: not recorded, because I did not capture the start time.
- **Cleanup:** the builder worktree is removed and its branch deleted with -d. The drain checkout is removed.
- **Left dirty under ROOT:** three untracked files, all left for the coordinator to commit at checkpoint: REQ-647-integration-report.md, REQ-648-integration-report.md and this report.
- **Discovered tasks:** the builder reported none. Review findings F1 and F2 are the only report-only items.
