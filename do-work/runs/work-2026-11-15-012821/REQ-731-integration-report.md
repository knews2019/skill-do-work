# REQ-731 integration report (ai-report prints the catalog rebuild command when a new report is missing from the catalog)

- Release: 0.321.0 "New Reports Remind You to Rebuild the Catalog" (minor bump, D-09: a new read-only flag, as REQ-728 and REQ-730 were bumped).
- Archive: do-work/archive/REQ-731-report-prints-catalog-rebuild-reminder.md (flat; UR-176 stays open).
- Range: pre 4feebfee (run artifacts), merge 15ad4b3e3dae6992bf369b2d8181f9b94711f94f. Finalization commit bddf5c85ebde04eb87516a3cfbd9e1ee68fef2f7 (phase cleanup_complete, no blocked paths, no reason codes).
- Merge: conflicts with REQ-730 (catalog lists second-opinion pages) at the ends of report_index.go and report_index_test.go; kept both sides (D-08). Prose auto-merged; the ai-report.md paragraph and guide bullet sit beside REQ-730's text with no contradiction. Package tests ok (26.9 s).
- Manual check on the merged tree: --check on a listed bundle printed nothing (exit 0); on an unlisted path it printed the rebuild line (exit 0).
- Gate: green at the merge, 141 s, start load 4.53. Contract regressions passed inside it on a quiet machine (core-checks 7 s, replace-text-section 9 s); the builder's reds at load 39 were the per-file budget only. Passed again after finalization.
- Review: Approve, 95% (opus reviewer, Route B standard). F1 --check does not normalize forms like ./ai-reports/x or .../index.html, giving a false reminder, never a missed one → report only. F2 stale D-02 text → amended in the REQ with a superseded note. F3 test nit → report only. Anti-bloat: 0 extra helpers, files or decorative tests.
- Heavy lanes at the merge, all executed, exit 0: do-work-cli-integrations 66 s, staged-skills 37 s, updater 68 s, installer 28 s.
- Timing events: builder-work, handback-merge, verification-gate (gate), verification-gate (heavy drain), review; folded into ## Timing.
- Lessons: builder proposed none; no satellite or lessons-index change.
- Cleanup: worktree removed, branch deleted (was 0e798128), worktree prune done.
- Report-only discovered task: architecture-report.md and stakeholder-report.md write bundles but do not run the new catalog check; architecture-report's output cap needs a wording decision first.
- Left dirty in the main checkout: nothing (this report and the progress log are committed as "[REQ-731] integration report").
