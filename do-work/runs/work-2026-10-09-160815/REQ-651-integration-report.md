# REQ-651 integration report (AI report on the board activity correlation)

- Not a release: no version bump, no changelog entry (NO_RELEASE=1). VERSION stays 0.305.86.
- Archive: do-work/archive/REQ-651-ai-report-on-board-activity-correlation-determinism.md (flat, UR-143 still has REQ-653 open).
- pre 5e996eed (run-artifact commit). First merge 0240a46f, then re-merge bd4d754c (full bd4d754c186baf68f3862b838f4536d270f039df), recorded as commit:. Finalization commit 8173895f, phase cleanup_complete, no blocked paths or reason codes.
- D-08 (integrator fix before the gate): REQ-650 shortened verify.go by 47 lines after the report's base 2eb14357, so all eleven verify.go citations in index.html pointed at the wrong code on main. Shifted each by 47 on the builder branch (7a04e2ec), checked each against the merge, and rewrote the drift note. activity_correlation.go and its test are unchanged since the base.
- D-07: the builder's Playwright captures went to the gitignored .playwright-mcp/ and were deleted. The folder holds no file from today and no tracked path was touched.
- Gate: maintainer-verify passed on the first run at bd4d754c, 137 s. advance recorded green-gate and the probe as satisfied.
- Review: 98%, Pass, one reviewer agent. F1: index.html:325 says "Thirteen have the shape [work run] REQ-NNN step N", but the evidence holds 12 such commits plus one remediation commit (cfcf63e4). F2: D-01's deletion of activity_correlation.go:111-113 would leave line 110's comment mid-sentence, and the count stays 54. Both are impact-negligible, report only, and left unfixed. Restatement sweep: nothing redefined.
- Heavy drain at bd4d754c, every lane executed and exited 0: queue-kanban-javascript 8 s, queue-kanban-browser 87 s, do-work-cli-integrations 85 s, staged-skills 44 s, updater 71 s, installer 27 s.
- Timing events: handback-merge, verification-gate x2 (gate, heavy drain), review, all folded into ## Timing. The builder-work event was skipped because dispatch (16:11:57Z) was 35 minutes before this integrator arrived, and recording it would charge the builder with the gap.
- Cleanup: builder worktree removed, branch deleted with -d, drain checkout removed, contract-regressions passed.
- Dirty under ROOT: only this untracked report. The REQ-653 worktree and branch exist and are not mine.
- Report-only discovered tasks (builder): the call-count test pins that the count stays constant, not that it equals five (activity_correlation_test.go:379-406). A detached repo-root checkout costs a sixth git command (verify.go:1357 at main).
