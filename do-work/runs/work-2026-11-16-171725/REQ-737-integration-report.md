# REQ-737 integration report (run work-2026-11-16-171725)

REQ-737 (the relation-map preset keeps phone titles readable after a click, and the map starts in the top half of the first phone screen) is finalized and released as 0.324.2. UR-180 is closed.

- Release: 0.324.1 → 0.324.2, a patch (D-29). The entry claims only the phone-click limit of 0.324.0, because 0.324.1 removed the other one.
- Archive: do-work/archive/UR-180/REQ-737-relation-map-phone-readable-titles.md. REQ-738 and UR-180's input.md moved into do-work/archive/UR-180/ in the same commit.
- Commits: run artifacts 38b278f2, builder merge 8e6d94e9 (implementation hash), catalog follow-up 3f6ebf38, finalization d6fad8a7 (cleanup_complete on the first call). The worktree was removed and the branch deleted with -d.
- Review: 95%, Pass, with the Restatement Sweep (no stale shipped text). The reviewer's rerun matched the hand-back: frame top y = 398 at 390 x 844, 12.32 px titles in 33 of 33 phone clicks, 70 of 70 Show steps, 0 words on arrowheads or cut at the frame edge, header fold works, 0 extra files. Every finding is report only:
  - F1 (user-visible): after a phone click, 88 to 223 px of empty space sits above the column title.
  - F2: the draft changelog title. Fixed in the entry.
  - F3 (REQ-738's): the header has no side gutter.
  - F4 (Nit): wording of check 15.
- Gate at 3f6ebf38: green on the first run (164 s, load 2.82), and the GREEN probe passed.
- Heavy lanes, all exit 0: queue-kanban-javascript 10 s, queue-kanban-browser 93 s, do-work-cli-integrations 79 s, staged-skills 42 s, updater 76 s, installer 31 s.
- After finalization, contract-regressions.sh and shipped-package-reference-contract.sh pass, so the lesson links resolve.
- Timing events: handback-merge, gate, review, heavy drain. Builder-work was skipped because the build had three passes with maintainer waits (D-26).
- Left dirty: nothing. I restored do-work/working/baseline.json after the pre-flight re-pass (D-27). This report was committed separately (D-30).
- Report-only task: REQ-738's header has no side gutter at 390 px. The coordinator captures it as a new request.
