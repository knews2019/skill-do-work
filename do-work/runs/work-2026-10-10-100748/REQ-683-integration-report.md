# Integration report: REQ-683 (interview cadence parsing reads AM and PM)

- Release: 0.305.96 "Interview Cadence Reads AM and PM, So 5:00 PM Is Stored as 17:00". Archive: do-work/archive/REQ-683-interview-cadence-am-pm.md.
- Merge: cb281775 (builder commit b83101b1, pre-merge 00873158). Finalization commit: 4f62b1a4 (phase cleanup_complete, no blocked paths or reason codes). Worktree and branch removed.
- Gate: DO_WORK_FAST_STAGE_REUSE=off maintainer-verify.sh exit 0 at load 4.35 (134 s wall, slowest file 21.49 s). Green probe REQ-683-probe.sh passed through advance. contract-regressions.sh exit 0 after finalization.
- Review: 95%, Pass (quick scan, Route A). Anti-bloat count: no new helper, option, file or test function; one import, one local, one regex group, one conversion block, four table rows. Restatement sweep: work-operating-model.md lines 104 and 405 still agree; no other restatement. Builder decision D-02 (refuse am/pm on hour above 12) accepted: both callers (deriveStandingSlots, interruption check) use only ok.
- Heavy lanes (detached checkout of the merge, all executed, exit 0): do-work-cli-integrations 78 s, staged-skills 39 s, updater 66 s, installer 28 s. No skips.
- Timing events: builder-work, handback-merge, verification-gate (repository gate), review, verification-gate (heavy drain); Timing section folded.
- Left dirty under ROOT: only other REQs' untracked run files (hand-backs, integrate briefs, integration reports) and REQ-683-integration-report.md (untracked for the coordinator). Sibling working REQ files and baseline.json untouched. Nothing of mine staged.
- Report-only discovered tasks: (1) the spec text (work-operating-model.md) does not say AM/PM is now read; optional doc mention. (2) Review F1: "5:00 p.m." and "5:00  PM" (two spaces) still give 05:00 silently; same bug class, outside the REQ. (3) Review F2: "0:30 PM" is accepted as 12:30 (builder D-04).
- No lesson bullet; lessons-index not touched. No new refusals from the merged CLI.
