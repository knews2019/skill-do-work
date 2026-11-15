# REQ-727 integrator brief — run work-2026-11-15-012821

- REQ: REQ-727 ('Ship a relation-map preset for ai-report second-opinion')
- REQ path: do-work/working/REQ-727-relation-map-second-opinion-preset.md
- UR: UR-176 (do-work/user-requests/UR-176/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/
- Shared guide (read it fully first): do-work/runs/work-2026-11-15-012821/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-15-012821/REQ-727-handback.md
- Operative name / branch: worktree-agent-REQ-727-relation-map-preset
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-727-relation-map-preset
- Dispatch instant (for the builder-work timing event): 2026-11-15T11:10:24Z
- Route: B
- tdd: true
- Closes UR-176: YES (UR_CLOSES=1); archive path do-work/archive/UR-176/REQ-727-relation-map-second-opinion-preset.md
- Release: yes, next patch version; re-read VERSION right before prep-payloads
- Members already finalized in this run: all thirteen earlier members of this run (see manifest.md)
- Set aside: none. Every other claimed member is still building or waiting for its integrator; their working REQ files and worktrees are not yours.
- Last successful integration of the run (run the Restatement Sweep): YES
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-15-012821/REQ-727-probe.sh
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/integrate-REQ-727/
- Write your final report to do-work/runs/work-2026-11-15-012821/REQ-727-integration-report.md and add it, the hand-back, the review file and the progress log to commit_paths so the tree is clean when you finish.
- UR-176 closes with this REQ. In the same finalization transaction, move these six flat-archived siblings into do-work/archive/UR-176/ and put each flat path and each destination in commit_paths, plus do-work/user-requests/UR-176/input.md and its archive destination: do-work/archive/REQ-726-record-no-review-canvas-decision.md, do-work/archive/REQ-728-judge-renders-one-named-page.md, do-work/archive/REQ-729-judge-relative-out-rule.md, do-work/archive/REQ-730-catalog-lists-second-opinion-pages.md, do-work/archive/REQ-731-report-prints-catalog-rebuild-reminder.md, do-work/archive/REQ-732-allow-docs-run-prompts-folder.md. Check with ls that this list matches what is on disk. If the finalizer moves closed-UR siblings itself, believe the command and follow its refusals.
- Link sweep before the manifests: git grep -n -E 'archive/REQ-(726|728|729|730|731|732)' -- ':!do-work/' ; re-point every hit to do-work/archive/UR-176/ and add those files (and do-work/lessons-index.md if a satellite changed) to commit_paths.
- You are the last integration: pass that fact to the reviewer with the sibling ids REQ-733 and REQ-734 (the example pages) and REQ-703, REQ-705 (how presets ship) so it runs the Restatement Sweep. Candidates: every shipped list of the second-opinion presets (the guide, ai-report.md, help text, tests that count presets).
- After the merge, run the ai-report index command on the merged tree and commit the regenerated ai-reports/catalog.json and ai-reports/index.html pair on main as the merge follow-up; the tstack bundle's row must list the new v1c03 page.
- Replace the two glob entries in the REQ's Scope with the real v1c03 file names before qualification (scope-drift needs exact paths).
- Builder decisions for you: D-16 (on a phone, after a click, the map zooms out to 0.2 to 0.33 and titles drop to 4 to 7 px; the panel under the map still lists every relation) and D-17 (the launcher prints the whole preset text above the page: 514 px on desktop, 1,303 px on a phone, so the map starts on the third phone screen). Do not widen this REQ for either: keep them as report-only discovered tasks unless the reviewer shows the REQ's own acceptance criteria fail, and state both plainly in your final message, because the maintainer confirmed the map as a standard option and should see its phone limits.
- Tell the reviewer to open the generated page ai-reports/2026-11-04_1747_proposal-tstack-review-workflows-rev1/index-v1c03-claude-claude-opus-5-5-relation-map-review-loop-map.html at 1440 and 390 width and re-measure the builder's claims (relation words visible at load, no overlap, no clipped or broken label, Tab focus in view, words on lines without JavaScript), then delete its captures.
- Version: minor (BUMP=minor), 0.323.0 to 0.324.0, a new shipped preset.
- CHANGELOG trap from REQ-736: name files in plain words or with paths that resolve from the installed core package; run _dev/tests/shipped-package-reference-contract.sh before finalization.
- The coordinator already committed the hand-back file. Leave it out of commit_paths unless you changed it.
- Never edit anything under /Users/t2/.claude/.
