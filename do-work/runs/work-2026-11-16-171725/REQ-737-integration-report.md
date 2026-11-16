# REQ-737 integration report (run work-2026-11-16-171725)

REQ-737 (the relation-map preset keeps box titles readable after a phone click, and the map starts in the top half of the first phone screen) is finalized and released as 0.324.2. UR-180 is closed.

- Release: 0.324.1 → 0.324.2 (patch, D-29). Entry title: "Relation-Map Preset Keeps Titles Readable After a Phone Click". It claims only the phone-click limit of 0.324.0, because 0.324.1 removed the other one (review F2).
- Archive: do-work/archive/UR-180/REQ-737-relation-map-phone-readable-titles.md. REQ-738 moved from the flat archive to do-work/archive/UR-180/, and UR-180's input.md moved to do-work/archive/UR-180/input.md. All of this is in the finalization commit.
- Commits: run artifacts 38b278f2, builder merge 8e6d94e9 (the implementation hash), catalog follow-up 3f6ebf38 (the tstack bundle row lists the final v1c04 page), finalization d6fad8a7 (cleanup_complete on the first call, no refusals). The worktree was removed and the branch deleted with -d.
- Review: 95%, Pass (standard review plus the Restatement Sweep). The reviewer's rerun of measure-relation-map.js matched the hand-back: frame top y = 398 at 390 x 844, 12.32 px titles in 33 of 33 phone clicks, 70 of 70 Show steps, 0 words on arrowheads or cut at the frame edge, the header fold works, 0 extra files or helpers. The sweep found no stale shipped restatement. Every finding is report only:
  - F1 (impact-user-visible): after a phone click or Show tap, 88 to 223 px of empty space sits above the column title in the frame. This confirms the coordinator's 95 px observation.
  - F2 (impact-negligible): the builder's draft changelog title and limits. Handled in the entry text.
  - F3 (REQ-738's, impact-user-visible): the short header has no side gutter on a phone.
  - F4 (Nit): check 15 leaves the overridden rule implicit.
- Gate: maintainer-verify at 3f6ebf38, green on the first run (164 s, load 2.82). The GREEN probe passed under advance.
- Heavy lanes at 3f6ebf38, all executed with exit 0: queue-kanban-javascript 10 s, queue-kanban-browser 93 s, do-work-cli-integrations 79 s, staged-skills 42 s, updater 76 s, installer 31 s.
- Post-finalization checks: contract-regressions.sh and shipped-package-reference-contract.sh both pass, so the lesson links to archive/UR-180/ resolve. The REQ-738 lesson link was re-pointed and the REQ-737 lesson was added.
- Timing events: handback-merge, verification-gate (gate), review, verification-gate (heavy drain). The builder-work event was skipped on purpose, because the build had three passes with maintainer waits (D-26).
- Left dirty under ROOT: nothing. The pre-flight re-pass rewrote do-work/working/baseline.json, and I restored it at once (D-27). This report is committed separately (D-30).
- Report-only discovered task: REQ-738's short header has no side gutter at 390 px. The coordinator captures it as a new request.
