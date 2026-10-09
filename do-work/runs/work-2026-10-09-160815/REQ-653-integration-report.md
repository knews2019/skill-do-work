REQ-653 (moving the fan-out orchestration prose into a companion reference) is integrated and released as 0.305.87. UR-143 is closed.

- **Archive:** `do-work/archive/UR-143/REQ-653-fan-out-orchestration-prose-companion-reference.md`. REQ-650, REQ-651, REQ-652 and the UR input moved into `do-work/archive/UR-143/` in the same commit.
- **Commits:** pre ad74c131. First merge 35803b1e. Review-fix re-merge ce1a14b6, the recorded `commit:`, with cumulative range ad74c131..ce1a14b6. Finalization 741b5e5b ended `cleanup_complete` on the first attempt, with empty blocked_paths and reason_codes.
- **Review:** 96%, Pass, after the delta re-review (the first pass scored 94%). There were three Minor findings, all impact-negligible:
  - F1 was a wrong 0.305.73 size (22,109 instead of 22,162) in `decisions/log.md`. Fixed on the builder branch, and the companion size was corrected to 5,732.
  - F2 was the dropped clause "an unmerged one can pass" in fan-out-reference Cleanup. Restored.
  - F3: D-04 does not list two cuts, both stated elsewhere → report only.
- **Wave-end sweep:** REQ-650 (disk probe measures the repo root only) is clean. REQ-651 recorded nothing redefined. REQ-652 (operator-blocked rule) is clean. REQ-653's own citations all land on existing headings. No member was unread.
- **Gate:** `maintainer-verify.sh` with stage reuse off exited 0 at 35803b1e (120 s) and at ce1a14b6 (130 s). The probe exited 0 both times.
- **Heavy lanes at ce1a14b6 (detached checkout, QUEUE_KANBAN_BROWSER set):**

| Lane | Result | Wall time |
| --- | --- | --- |
| queue-kanban-javascript | reused | 8 s, from its executed run at 35803b1e |
| queue-kanban-browser | executed, exit 0 | 84 s |
| staged-skills | executed, exit 0 | 37 s |

  The earlier drain at 35803b1e executed all three green (8 s, 80 s and 35 s). No lane was skipped.
- **Timing events:** hand-back merge, two repository gates, heavy drain and review. They were folded into `## Timing`. The builder-work event was skipped. The hand-back landed about 16:46Z and this integrator started at 17:00Z, so timing from dispatch would charge the wait to the builder.
- **Lessons:** the condition-preserving-prose-extraction bullet was added. The REQ-652 and REQ-650 lesson links were repointed to `archive/UR-143/`, and both rows in `lessons-index.md` were refreshed. `contract-regressions.sh` passed after finalization.
- **Cleanup:** the builder worktree, the builder branch (`-d`) and the drain checkout were removed.
- **Dirty under root:** only this report, which is untracked on purpose.
- **Report-only discovered task (from the builder):** stale plain-text mentions of work-reference.md → Worktree Dispatch Mode remain in kb/wiki, ai-reports, one HANDOFF file, CHANGELOG entries and ADR-018. These are historical records and were left alone.
