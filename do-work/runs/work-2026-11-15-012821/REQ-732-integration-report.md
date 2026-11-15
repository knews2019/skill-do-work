# REQ-732 integration report (allow docs/run-prompts/ in the staged-skills contract and move the run prompts there)

- Release: none (NO_RELEASE=1; the finalizer agreed, no shipped file changed).
- Archive path: `do-work/archive/REQ-732-allow-docs-run-prompts-folder.md` (flat; UR-176 stays open because REQ-727 is blocked on the maintainer).
- Commits: run artifacts `8ad193d8`, merge `ce3cb5dc3b4099574b9d5ca103fd6cffc426957a` (range 8ad193d8..ce3cb5dc, no conflicts, no seams), finalization `35f32e0f6c876536a7f81d7c828af74ceba1bd13` (`cleanup_complete`, empty blocked_paths and reason_codes).
- Empty `_dev/run-prompts` folder: git removed it at the merge, so no `rmdir` was needed.
- Gate: `GATE-EXIT=0 WALL=125s` at the merge, 1-minute load 2.90. GREEN probe `REQ-732 probe PASS`; green gate recorded.
- Review (Route A quick scan, opus): Approve, 95%. F1 `.gitattributes:19` still 141 columns although D-04 says it was reflowed → report only. F2 the auto-memory note `run-prompts-not-in-root-docs.md` still told agents to use `_dev/run-prompts/` → fixed: I deleted the note and its `MEMORY.md` line (outside the repository, no commit). F3 queue Source lines in REQ-725, REQ-726, REQ-727 and REQ-732 still name `_dev/run-prompts/`, kept per triage D-01 → report only. F4 `.gitattributes:57` attribute one column off → report only. Anti-bloat: 0.
- Heavy lanes (planner: coverage uncertain for `.gitattributes`, all six selected), all executed, exit 0, none skipped: queue-kanban-javascript 7 s, queue-kanban-browser 77 s, do-work-cli-integrations 59 s, staged-skills 31 s, updater 65 s, installer 26 s.
- Timing events: builder-work, handback-merge, verification-gate (repository gate), review, verification-gate (heavy drain); folded into `## Timing`.
- Lessons: no satellite bullet (the builder proposed none). `contract-regressions.sh` exit 0 after finalization.
- Report-only discovered task: the `staged-skills` heavy lane coverage in `_dev/tests/heavy-lanes.json` does not include root `docs/`, so a new file there does not trigger the lane.
- Dirty in the main checkout: nothing. This report and the last progress-log line are committed as `[REQ-732] integration report`.
