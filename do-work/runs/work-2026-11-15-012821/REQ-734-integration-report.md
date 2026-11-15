# REQ-734 integration report (supports-and-limits table page, compared with the map page)

- Release: none (NO_RELEASE=1; finalization agreed). UR-178 is closed.
- Archive: `do-work/archive/UR-178/REQ-734-proposal-table-and-map-comparison.md`, with `do-work/archive/UR-178/input.md`.
- Commits: pre `d31159be`, first merge `c812ea4b`, catalog rebuild `0ccdf98d` (the Tinka row lists v1c04), review fix `e760b5b5` on the builder branch, re-merge `13322db2` (implementation hash), finalization `8c6bb12d`.
- Gate: green at `0ccdf98d` (124 s, load 2.73) and at `13322db2` (122 s, load 2.49). Probe PASS.
- Judge: `ai-report-judge --page` on the new page gives `pass` with no findings, both before and after the fix. The output went to scratch and was deleted.
- Heavy at `13322db2`, all exit 0: kanban-js 6 s, kanban-browser 72 s, cli-integrations 56 s, staged-skills 30 s, updater 66 s, installer 26 s.
- Review (opus): Approve, 94%, then 96% after the delta re-review. Fixed: F1 (the page's catalog correction became false after the merge) and F2 (`ai-report.md:87` moved to line 89). Report only: F3 (two inaccuracies in the hand-back) and F4 (below). Anti-bloat: 0.
- Map vs table (reviewer): at 1440 the table shows 4 of 37 relations on the first screen and the map shows 0. The map shows a line's words only after a click on a box. Its seeds, limits and box texts sit behind clicks or in collapsed groups. Without a click, its relations can be read only in the "Every relation" table, 4 to 6 screens down. The table hides nothing behind a click: all 37 relations are readable within about 4 screens. It leaves out the full box texts and the extra source chips, and links to `index.html` for them. At 390 both pages show 0. The table stacks over about 9.5 screens. The map is a 1910-px drawing in a 331-px window, its relation table needs a sideways swipe, and the page is about 11.5 screens long. No winner is picked.
- Restatement Sweep: one gap. `skills/do-work-toolbox/actions/help.md:18` lists `ai-report judge <dir>` without `--page` (REQ-728). The line leaves out a form but says nothing false (D-20). A fix changes a shipped file, so it needs a follow-up REQ. REQ-729, REQ-726 and REQ-725 are unread for the sweep: their sweep lines record nothing.
- Decisions: D-16 (catalog follow-up), D-17 (Scope entries now name the files the launcher created), D-18 (reply file kept as the editor's record), D-19 (no lessons-satellite bullet), D-20 (help.md judgment).
- Timing: 8 events folded into `## Timing`.
- Report only: help.md:18; the toolbox `coding-guardrails.md` lacks the "Earned defense" paragraph (since REQ-170, not this run); the launcher's direction block has no side gutter; the standalone judge leaves a temp folder when run without `--out`.
- Left dirty under ROOT: nothing. This report and the last progress line are committed as `[REQ-734] integration report`.
