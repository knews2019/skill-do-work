# REQ-745 builder brief

In your assigned worktree at `f3840b087913d8501267dda1c1395580304f2ba4`, implement durable likes, dislikes and optional short notes for the gallery's report and figure targets. The user opens it with `just run-report-gallery`. Reopening must restore a report like with no note and an independent dislike and note on one figure. The agent generating the next round reads the same stored entries.

Read the working request, the shared gallery handoff and RED before editing. The source scope and acceptance checklist are in the working REQ. Start with the repo's Step 6 crew rules, the listed action prime, the shell prime for launch commands, and the UI and API specs as relevant.

## Build and verify

1. Start with an automated test that fails on the base for the captured behavior. The preparation probe demonstrates the missing Just recipe and live mode and the lack of controls/store. It is supporting evidence; save your own test-first failure and then the passing run in the existing harness.
2. Reuse the existing CLI renderer and its curated report inventory, and handle generated reports in an ordinary installed consumer with no preview manifest. Show report and figure controls on the opened gallery page. Resolve real figure identities from the report, including the phase-by-phase diagram, and keep feedback data bound to the report and figure across reload.
3. Keep live feedback local and project-scoped. Cover note round-trip, optional empty note, independent update, and one-record-per-target behavior. Verify the Just launch in a disposable source checkout and the same installed command path in a fresh consumer. Read the stored JSON independently of the page. A missing or failed save should produce actionable UI state.
4. Keep edits to the feedback seam only in the shared dispatcher and page renderer. REQ-746 owns starter prompt markup and checked-in curated gallery regeneration. Reuse the existing gallery hooks, preview links and generated-page ownership guard. Run focused tests on your branch. The integrator runs the full gate.

## Handback

Use the worktree path and actual branch supplied by the coordinator. Preferred operative name follows the REQ slug, `worktree-agent-REQ-745-persistent-report-figure-feedback`, subject to the dispatcher's actual allocation. Prefix each implementation commit with `[REQ-745]`.

Work on your own branch. Your two main-checkout outputs are the absolute handback and progress paths below. Keep lifecycle and run-manifest edits, the other REQ's files, the curated derivative, archived reports and frozen verification receipts outside your write set. Keep synthetic marks in fixtures.

The handback must give operative branch/path and commit, exact new/modified/deleted paths, any integration seam as an exact insertion and anchor, test-first RED and GREEN, installed and curated results, facts about what you implemented, and `## Decisions` / `## Discovered Tasks` where applicable. Include P-A-U proof for the owner to update the request. Log `<UTC> <phase> <command>` before and after commands expected to take over a minute. Finish all commands before handing back.

- Working REQ. `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-745-persistent-report-figure-feedback.md`
- Shared gallery handoff. `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-18-191126/shared-gallery-contract.md`
- Preparation RED. `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-18-191126/REQ-745-red.log`
- Handback. `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-18-191126/REQ-745-handback.md`
- Progress. `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-18-191126/REQ-745-progress.log`
