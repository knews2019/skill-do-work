REQ-647 (clarify says when a released REQ is still earmarked and offers to clear it) is released as 0.305.83 and archived flat.

- Archive: do-work/archive/REQ-647-clarify-reports-released-req-still-earmarked.md
- Run artifacts commit ffc1a7bb (<pre>). Merge be1f0ea4 (full be1f0ea4e7dc36144728ea0c9aca64fed0bcc7e2), clean, no seams. Finalization commit fc920b2c, phase cleanup_complete, no blocked paths or reason codes.
- Qualify gate satisfied with no findings. Repository gate exit 0 on the first run, 117 s. Probe and green-gate records satisfied.
- Review: Pass, 98%. Three nits, all impact-negligible, report only. N1: the schema citation uses the short heading "Request File Schema", the house form at about 15 sites. N2: the two-option question names no recommended option. N3: the builder stamped its Discovered Tasks impact-low, which is not a valid token; I re-stamped them in the REQ as the review said. No fix needed.
- Heavy lane: staged-skills, executed, exit 0, 43 s, from a detached checkout at be1f0ea4 with ROOT's durations header and QUEUE_KANBAN_BROWSER set. Drain checkout removed.
- Timing events recorded: builder-work, handback-merge, verification-gate. I skipped the review and heavy-drain events. The recorder has no end-instant flag and times to now, so recording after the fact would charge them with idle gap.
- Lesson bullet appended to _dev/primes/lessons-action-files.md below REQ-648's bullet. The lessons-index row is refreshed to 6852 tokens. contract-regressions.sh passed after finalization, so the link resolves.
- Builder worktree removed and branch deleted with -d.

Left dirty under ROOT, none of it mine:
- do-work/runs/work-2026-10-08-155307/manifest.md has the coordinator's row edits. I did not change my row.
- REQ-648-integration-report.md, REQ-649-handback.md and REQ-649-integrate.md belong to the sibling REQs.
- This report is untracked, as the guide asks.

Report-only discovered tasks:
- impact-user-visible: clarify Step 5.5 calls unblock without --commit, so a clarify session that only unblocks leaves an uncommitted lifecycle change. This was already true before this change.
- impact-rule-change: the Builder Was Right / Discarded section of clarify describes hand edits, while Step 5 says answer has no hand-edit fallback.
- Placement: REQ-648's lesson bullet, and mine beside it, sit at the end of lessons-action-files.md under the Agent Compatibility section. The other alternate-writer-contract-drift bullets are near line 65.
