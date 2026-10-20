# REQ-703 integrator brief — run work-2026-10-20-200259

- REQ: REQ-703 (second opinions accept a named preset and a free-text editor direction)
- REQ path: do-work/working/REQ-703-second-opinion-preset-and-editor-direction.md
- UR: UR-160 (do-work/user-requests/UR-160/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-20-200259/
- Shared guide (read it fully first): do-work/runs/work-2026-10-20-200259/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-10-20-200259/REQ-703-handback.md
- Operative name / branch: worktree-agent-REQ-703-second-opinion-presets-and-direction
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-703-second-opinion-presets-and-direction
- Dispatch instant (for the builder-work timing event): 2026-10-20T20:17:40Z
- Route: C (full pipeline; thorough review against the Plan)
- Wave membership: REQ-703, REQ-704 (serial). Set aside: none.
- Closes UR-160: no (UR_CLOSES=0; flat archive path do-work/archive/REQ-703-second-opinion-preset-and-editor-direction.md)
- Last successful integration of the wave (run the Restatement Sweep): no (REQ-704 is the last member and carries the wave-end sweep)
- Lesson satellite for a Go lesson: skills/do-work/tools/do-work-cli/lessons-do-work-cli.md (action-file lessons: _dev/primes/lessons-action-files.md)
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- Pre-flight was recorded by the coordinator at a8c0cf11 (gate green, 130 s); `## Pre-Flight` is in the REQ. Never re-run pre-flight.
- The REQ is tdd: true: the Testing section must carry the RED/GREEN evidence for TestSecondOpinionPresetResolutionAndArgumentRules, TestSecondOpinionDirectedRunNamesAndMarksPage and TestAIReportIndexSkipsSecondOpinionPresetsFolder from the hand-back.
- D-04 in the Plan is ESCALATE (an uncommitted project preset is refused by the outside-bundle Git check). At Step 8, route it per the Builder-Decided Follow-up rules (a pending-answers follow-up for the user unless you find an existing home); the coordinator decided to keep the behavior for this run.
- Release: VERSION is 0.306.5 at run start; this REQ ships 0.306.6 (patch) unless you judge minor; re-read VERSION right before prep-payloads.
- Reviewer anti-bloat: ask it to count helpers, options, files or tests the brief did not name.
