# REQ-725 integration report (sort the local memory notes by destination)

- Release: none (NO_RELEASE=1; only `ai-reports/` changed, which is export-ignored).
- Archive: `do-work/archive/UR-175/REQ-725-sort-local-memory-notes-by-destination.md`. UR-175 closed: its `input.md` and REQ-724 (moved from the flat archive by the closure) are now in `do-work/archive/UR-175/`. No file outside `do-work/` cited REQ-724's flat path, so nothing was re-pointed.
- Merges: pre `e5972dc5`; first merge `fbb5d529`; catalog commit `27e77d53`; review-fix re-merge `406d32aa` (implementation hash). Finalization commit `69579c22`.
- Report corrections on the builder branch: `189da332`, before the first merge (D-11, D-12). The row for `run-prompts-not-in-root-docs.md` is removed (the REQ-732 integrator deleted that note), the counts are 46 notes and 28 deletes, and one sentence records the deletion. The REQ-724 and REQ-732 rows were re-checked on main and now cite what landed. Two line numbers that moved were fixed. Then `dd970843` fixed review F1 and F2 (D-14).
- Final count: 5 move to the skill, 13 move to the repository, 28 delete, 46 in total. Page: `ai-reports/2026-11-15_0137_proposal-REQ-725-memory-note-destinations/index.html`.
- Catalog: the `ai-report-index` rebuild gave 50 bundles, committed as `27e77d53`. `--check` printed nothing for the bundle afterwards. The rebuild at `406d32aa` changed nothing.
- Gate: green at `27e77d53` (126 s, load 4.84) and at `406d32aa` (123 s, load 2.60). The probe passes with 46 notes.
- Heavy: six lanes were selected, because no lane covers `ai-reports/`. All six were green at `406d32aa`: kanban-js 6 s, kanban-browser 74 s, cli-integrations 57 s, staged-skills 31 s, updater 66 s, installer 26 s. They were also green at `27e77d53`.
- Review: opus, 92%, Pass, 16 rows sampled. F1 (wrong line for the Lessons phase) and F2 (a report folder that moved) were fixed anyway. F3 (a REQ-724 path that is true only after the closure) is now true. Anti-bloat: 0.
- Timing: 8 events recorded and folded into `## Timing`.
- Lessons: no satellite bullet, because no satellite covers memory notes (D-15).
- Report-only discovered tasks: 11 of 18 toolbox crew-member copies differ from core. `commit.md` examples commit without a pathspec. Bare `release --manifest` reportedly writes the tree.
- Left dirty under ROOT: nothing.
