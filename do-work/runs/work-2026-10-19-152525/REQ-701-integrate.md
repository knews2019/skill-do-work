# Integrator brief: REQ-701 (a failed second-opinion model check keeps the finished page; a subagent model never makes the backend model ambiguous)

Read `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-19-152525/INTEGRATOR-GUIDE.md` in full first; this file holds only the facts for your REQ.

- REQ id: REQ-701. Working file: `do-work/working/REQ-701-second-opinion-model-check-keeps-finished-page.md`. UR: UR-158 (`do-work/user-requests/UR-158/input.md`).
- Run directory: `do-work/runs/work-2026-10-19-152525/` (ROOT-relative). Progress log: `REQ-701-progress.log` there.
- Hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-19-152525/REQ-701-handback.md` (landed; consume it, never re-dispatch).
- Builder brief (what the builder was told): `REQ-701-brief.md` in the run directory. GREEN probe for the test gate: `REQ-701-probe.sh` there (Go tests only, about 10 s).
- Operative name / branch: `worktree-agent-REQ-701-keep-finished-page-unverified-model`. Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-701-keep-finished-page-unverified-model`. Base commit cc3e72c2.
- Dispatch instant (for the builder-work timing event): 2026-10-19T15:35:45Z. `dispatch_at` is already stamped; stamp `builder_handback_at` with your own `date -u` if absent.
- Route A: no Scope section, no scope-drift comparison, no pre-flight. Section order after the Plan skip note: Implementation Summary → Qualification → Testing → Review → Lessons Learned → Orientation (Decisions and Discovered Tasks by hand after Implementation Summary; Heavy Plan + Result by hand between Review and Lessons Learned).
- Wave membership: REQ-701 (you) and REQ-700. Set aside: none. Mid-run addendum text: none. Hand-back emphasis note: none.
- This REQ is NOT the wave's last successful integration; REQ-700's integrator carries the wave-end Restatement Sweep. Tell the reviewer it is a Route A quick scan with the anti-bloat count.
- Release: yes (a shipped Go change plus two shipped docs). VERSION is 0.306.2 at run start; re-read it right before `prep-payloads.py`.
- UR_CLOSES=0: REQ-700 is still in flight, so the archive path is the flat `do-work/archive/REQ-701-second-opinion-model-check-keeps-finished-page.md` and your lesson link (if any) points there; REQ-700's integrator repoints it to `archive/UR-158/` later.
- Lessons satellite for a Go lesson: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (family `native-backend-identity` likely), then refresh its row in `do-work/lessons-index.md`.
- Sibling state to leave alone and name in one progress line: `do-work/working/REQ-700-*.md` (pre-dispatch sections, possibly dirty), `do-work/working/baseline.json`.
- Never run Step 10 (`advance --checkpoint`), never `recover` in any form, never claim another REQ. When finalization reports `cleanup_complete`, clean the worktree and branch (`worktree cleanup REQ-701 --name <operative>` or the git commands in the guide), run `bash _dev/tests/contract-regressions.sh` once more, write `REQ-701-integration-report.md` in the run directory (untracked), and report back.
