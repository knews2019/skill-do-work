REQ-675 (core review names four delivery stages) is integrated and released as 0.305.90.

- Archive: do-work/archive/REQ-675-review-delivery-stages.md.
- Range: pre 623f00d5 (run artifacts commit). First merge f7dad3b1. Final merge fec8d220 (recorded as commit:). Finalization commit 19ec46fd (cleanup_complete, nothing blocked).
- Review: 94%, Acceptance Pass, verdict Approve. I fixed F2 on the builder branch (67e0ed86) and re-merged with the same pre. F2: Step 7 did not say which line gets the stage names; it now also names the persisted **Acceptance:** line. Delta re-review after re-gate and re-drain: 95%, Pass, no new finding. Report only: F1 (impact-rule-change), shared-principles.md:15 "Acceptance cannot be exercised, Record Untested" does not cover the partial case. F3 (impact-negligible), the sample-archived-req.md:99 example line names no stage. F4 (impact-negligible), review-work.md:460 says Route A suggested testing is "usually empty or 1 item".
- Behavioral exercise (reviewer, real run against the edited text): GREEN. The "publish the updated rules page" scenario gives Acceptance Pass scoped to implementation and integration, plus "Deployment: unassessed" and "Live acceptance: unassessed" lines. The base text gives a bare Pass (RED).
- Gate: maintainer-verify exit 0 at f7dad3b1 (133 s) and at fec8d220 (144 s). Probe exit 0 at both. contract-regressions after finalization: exit 0 (21 s), so the lesson link resolves.
- Heavy: staged-skills only. Executed exit 0 at f7dad3b1 (43 s) and at fec8d220 (35 s).
- Timing events (8, folded into ## Timing): builder-work 237 s, handback-merge 13 s and 7 s, verification-gate 171 s and 151 s (gate plus probe), heavy drain 59 s and 44 s, review 200 s. The delta re-review was not timed.
- Lesson: one bullet added to _dev/primes/lessons-action-files.md (family alternate-writer-contract-drift, the rule must name the persisted line). lessons-index row tokens 7189 -> 7330. REQ-678's integrator must re-point this link to do-work/archive/UR-151/ when UR-151 closes.
- Discovered task (report only, impact-negligible): the builder flagged the REQ-676/677 run probes failing the quiet-grep audit. The coordinator already fixed this in c4dda51e.
- Changelog date is 2026-10-09 (UTC); the 0.305.89 entry says 2026-10-10.
- Left dirty under ROOT: untracked REQ-676-handback.md and REQ-677-handback.md (sibling builders, not mine), and this report (untracked on purpose). Builder worktree and branch removed. Nothing pushed.
