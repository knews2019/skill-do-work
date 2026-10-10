# REQ-682 integration report

REQ-682 (a hidden Timeline view no longer redraws its rows when the user scrolls another view) is integrated and finalized.

- Release: 0.305.100, "A Hidden Timeline No Longer Redraws When You Scroll Another Board View".
- Archive: do-work/archive/REQ-682-hidden-timeline-scroll-guard.md
- Commits: pre a7b78f62 (run artifacts), merge 32552b5b (full 32552b5badd349ff7fd0baac9056e9ba31192a7e), probe fix bc7de192, finalization 747b7b4b. Builder branch and worktree removed.
- Probe fix: REQ-682-probe.sh now sets QUEUE_KANBAN_JAVASCRIPT_PROBES=on, fails on any SKIP line and requires a PASS line for the Node-lane step. Advance runs a probe through `sh -c`, and the probe's process substitution failed there (BLOCKED-PROBE-FAILED); replaced with a here-string loop. Rerun passed.
- Gate: `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` exit 0 at the merge, wall 120 s, load 3.18, one run. Probe exit 0 (about 8 s). contract-regressions.sh exit 0 after finalization.
- Review: Approve, 95%. Guard is one null-safe early return, no listener release, state or observer. Only one scroll listener exists on #board-main, so the arrival dispatch triggers nothing else. Bloat: buildTimelineScrollProbeSite (verbatim move, three callers) and two result structs; all used. Restatement sweep: no stale sentence in board docs, primes or lessons.
- Minor findings, all report only: F2 two stale comments in the moved helper (timeline_scroll_browser_probe_test.go:309-310 says "Both" for three callers; :335 says "below" for a test now above). F3 the synchronous dispatch can run one wasted render on a filter-reset arrival (under 1 ms). F4 the probe's Node filter does not cover TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip; the heavy lane does.
- Scope: javascript_behavior_c_test.go (one stub line) was undeclared but coordinator-approved (D-02). Two SCOPE-DECLARED-NOT-TOUCHED warnings (":43-44", "renderVisibleRows") are parser false positives.
- Heavy lanes (detached checkout, Chrome set): queue-kanban-javascript 10 s, queue-kanban-browser 84 s, staged-skills 44 s. All exit 0, executed, none skipped.
- Timing events: builder-work, handback-merge, verification-gate (gate), review, verification-gate (heavy drain); Timing section folded.
- Lesson added: family hidden-panel-listener-guard in _dev/primes/lessons-kanban-board.md (flat archive link), lessons-index row refreshed (6118 tokens).
- No new refusals from the merged code; finalization succeeded first try (cleanup_complete, no blocked paths).
- Left dirty under ROOT: only untracked run-directory files of other REQs plus this report (coordinator commits at checkpoint).
- Discovered tasks: both from the hand-back, report only (vacuous Node step, now fixed in the probe; hidden-panel geometry removed by the guard).
