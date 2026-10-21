# Review: REQ-704

**Approve**. A second opinion can now build on `index.html` or on any earlier tree page. Each directed pass, and each pass on a tree page, gets a launcher-allocated tree id in both names, a parent link and tree metas in its header, and a checked "Changed in this iteration" list. A passing tree run rebuilds one iteration list on `index.html`. A plain run is unchanged. Nothing blocks. Seven noncritical findings, all report only.
Route C | merge `92443689` (range `87614289..92443689`, 8 files, +468 / -64)

The UR-160 D2 project-folder preset lookup is not scored against this REQ: REQ-705's user answer superseded it (user presets live only inside the bundle).

### What's built
- `ai-report-second-opinion` takes `[<bundle-dir>|<page.html>|latest]`. A page target resolves its bundle through the unchanged `secondOpinionBundle`, must sit in the current revision, and must be a regular `index.html` or a tree page (`report_second_opinion.go:171-186`). Bundle and base resolve before the preset lookup, so a bundle preset still resolves for a page target (pinned).
- Tree ids are allocated at publish time inside the one name builder `secondOpinionNames` (`:1419`): parent + `cNN`, one above the highest number used under the parent by any tree page (descendants count), overflow above 99 refused. Names are `index-<id>-<stem>.html` and `<id>-<stem>-reply.md` with no `-N` suffix.
- Decoration now runs after naming. It adds the `-id`, `-base`, `-model`, `-date` metas and the `Iteration <id>, built on <parent>` paragraph. It strips every copied `ai-report-second-opinion-*` meta in either attribute order and any copied tree block (`:1611-1643`).
- After the judge, `secondOpinionChangeList` (`:1683`) turns a missing list or an unresolved `#rev-*` link into `AI-REPORT-SECOND-OPINION-CHANGE-LIST` findings that turn `pass`/`skipped` into `fail`. The page stays.
- On judge `pass` a tree run calls `secondOpinionTreeList` (`:1650`) instead of the per-filename banner. It lists every regular tree page in byte order (pre-order), with id, preset, model and date.
- Rows gain `tree_id` and `base_page` (omitempty). The prompt, `ai-report.md`, the guide and `ai-report-reference.md:148` describe the page target, the ask-on-unclear-base rule, tree names and the iteration list.
- Not built, and correctly so: the release files and the lesson-link sweep. The integrator writes them at finalization.

### Decisions / risks for you
- None that need a choice before release. F1 (a plain `ls` sorts in tree order only under byte-order collation) is the one user-visible caveat. The UR asked for "a plain sort", and the launcher's own list is correct on every platform.

### Findings

**Important:**
- None.

**Minor:**
- F1. `skills/do-work-toolbox/actions/ai-report.md:73` ("so a plain `ls` lists the tree in order") and `skills/do-work-toolbox/docs/ai-report-guide.md:107` ("a plain file listing shows the tree in order"). Pre-order holds for byte order, because `-` (0x2D) sorts before `c`. I confirmed it on macOS under `C` and `en_US.UTF-8`. glibc `en_US.UTF-8` collation ignores `-` at the first level, so on Linux `index-v1c01c01-codex-...` would sort before its parent `index-v1c01-codex-...` (the comparison is `0` vs `o`). Verdict PLAUSIBLE: reasoned from glibc collation rules, not run on Linux. The iteration list on `index.html` uses Go's byte sort and is correct everywhere. Possible fix: say "a byte-order listing (`LC_ALL=C ls`)". — impact-user-visible → report only
- F2. Plain-run-only sentences that REQ-704 left unqualified. `ai-report.md:75`: "a fresh generation allocates a numeric suffix and therefore gets a new banner". `ai-report-guide.md:107`: "Repeated generation adds matching numeric suffixes to the page and reply". `ai-report-guide.md:111`: "a new generation gets a new suffix and banner", and "after the direction header on a directed page". A tree run gets a new tree id, no numeric suffix and no banner. A tree run with no preset and no direction (a pass on a tree page with no flags) also gets the header block before the backlink. The nearby tree sentences make the meaning recoverable. — impact-negligible → report only
- F3. `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:95` (the REQ-703 bullet) still says "`secondOpinionDecorate` removes a copied meta only when `name` is its first attribute, so a reordered copy survives". REQ-704 D-06 made that false: the strip now uses `reportMetaTag` and `reportMetaAttribute` and handles either order. The integrator already edits this line for the UR-160 link sweep. It could append "(fixed in REQ-704: the strip reads either attribute order)". — impact-negligible → report only
- F4. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_tree_test.go`: no test pins two D-06/D-07 behaviors. (a) The widened strip. Every fake editor writes a fresh page, so no test has a copied content-first meta or a copied tree block on a tree page, which is the REQ-703 F5 closure. (b) The tree-run refusal of a template without `{{BASE_PAGE}}`/`{{CHANGE_LIST_REQUIREMENT}}`. A reviewer probe in a scratch copy confirmed both work. It also confirmed that decoration with tree marks is idempotent and that first-wins metas give the child's values. The new tests also ignore `os.MkdirAll`/`os.WriteFile` errors at :47-48, :98-99, :121, :146 and :168, the same pattern as REQ-705 F4. — impact-negligible → report only
- F5. `report_second_opinion.go:1426` (allocation) to the publish in `secondOpinionPublication`. Two publishes on one base that run at the same moment with different stems (for example a session `--finalize` while a `--backend codex` run publishes) can both pick the same `cNN`. `rootedPublishFile(..., false)` refuses only an identical name. The window is the time between `os.ReadDir` and the page write, and `--all` runs backends one after another, so it does not happen in normal use. A duplicate id would show as two rows with the same id in the iteration list. — impact-negligible → report only

**Nit:**
- F6. `report_second_opinion.go:1696`: two links to the same missing anchor give two identical findings. — impact-negligible → report only
- F7. `report_second_opinion.go:258,655,974,1625`: `cmp.Or(direction.Base, "index.html")` appears four times, and the `base_page` path is built twice (handler and finalize). One `basePage()` method on `secondOpinionDirection` would hold it. — impact-negligible → report only

### Requirements Checklist

- [x] D4 base selection: positional accepts `latest`, a bundle, or a page. A page gives bundle and base, and must be `index.html` or a tree page. The default is `index.html`. Prose: the user usually names the HTML file, and when the base is unclear the agent loads clear-questions and uses "the environment's ask-user prompt", best guess first plus free text. Delivered: `:171-186`, `ai-report.md:21,65-67,71`, guide :72-74. The old non-tree name and a missing tree page are refused (pinned).
- [x] D5 editor input: INPUT names `{{BASE_PAGE}}` as the page to rewrite, keeps `index.html` and the repository as fact sources, and tells the editor to read the base's change list first. Delivered: prompt line 5. A tree prompt carries `Base: <v1c01 page>` (pinned).
- [x] D6 tree ids: `v1` root, two-digit `cNN`, unlimited depth, id in both names, launcher allocates at publish. A session finalize allocates at finalize (`v1c01c02`, pinned). Deleted-number non-reuse and the overflow refusal were confirmed by a reviewer probe.
- [x] Which runs get an id: preset, direction, or a non-`index.html` base (`parentTreeID`, `:119`). A plain run keeps `index-codex-actual-model-1.html`, gets no tree fields and no tree marks (pinned).
- [x] D8 change marks: the prompt constant requires the list, the `#rev-<n>` links and the badge. The launcher checks every link target and a missing list. A failure is a judge finding, the page and reply stay, and `index.html` is untouched (pinned for `#rev-2`, `data-id` look-alike and no list).
- [x] D9 walk-through: one launcher block on `index.html`, rebuilt from tree filenames plus metas on every passing tree publish. It lists id, preset, model, date and link, in pre-order matching `ReadDir`. A script-string marker survives, a later plain run keeps its banner, and the block is unchanged (pinned). Tree pages carry their id and parent link in the header and the `-id`/`-base` metas. They get no per-filename banner.
- [x] D10: the three existing non-tree pages never match the pattern. They are never listed and are refused as bases. `ai-reports/` is untouched in the range.
- [x] Preservation: earlier tree pages and the base are never rewritten. The existing preservation and original identity and bytes checks run unchanged before the `index.html` write.
- [x] Row gains `tree_id` and `base_page` (omitempty). `result_model.go:719-721`. `base_page` is set on prepare, run and finalize rows of tree runs (pinned).
- [x] Docs: `ai-report.md` and the guide document the page target, the ask rule, tree names and the index list. The `## Second-opinion form` heading is unchanged. The release and changelog are the integrator's (not scored).
- [x] Plan D-01 to D-13 and builder D-14 to D-20 each match the code. D-05 order: strip, backlink, variation, names, decorate, publish, unverified return, judge, change list, then `index.html` write. D-12 pins moved to `v1c01`/`v1c02`/`v1c03`, a traced cross-REQ change (REQ-703). D-13: new tests in the new file. D-16 layout has no whitespace outside the markers. D-17 aria-label is kept. D-19 fallback. D-20 `banner_state` reuse and change detail.
- [x] Red-Green Proof GREEN (1) to (6): (1) and (2) `TestSecondOpinionTreeTargetsNamesAndParentHeader` (`v1c01`, `v1c02`; slug form in `TestSecondOpinionDirectedRunNamesAndMarksPage`). (3) the same test, `v1c01c01` visual with the parent link and metas; the resolving change list passes the launcher check in the tree-list test. (4) `TestSecondOpinionTreeListOnIndexFollowsTreeOrder` (listed order equals `ReadDir` order). (5) `TestSecondOpinionChangeListTargetsAreChecked`. (6) the plain-run block of the tree-target test.
- [x] Implementation Summary claims match the diff. All 8 files and every claimed change are present.

### Acceptance Testing

**Result: Pass** (stages covered: implementation and integration)
- At `c5a8aa8d` (code identical to merge `92443689`): the four named tests passed (`-count=1 -v`, 3.7 s). `go test -count=1 -run SecondOpinion ./internal/toolboxcommands ./internal/resultmodel` passed (20.9 s and 0.5 s).
- Reviewer probes in a scratch copy of the module (never in the repo): decorating a copied base page with content-first parent metas, a copied header and a copied tree block gives one fresh set of child metas and no copied block, and decorating twice gives the same bytes. Allocation skips a deleted page's number while its child remains (`v1c02`), and `v1c99` present gives `no free two-digit child id under v1`. A pre-REQ-704 template is refused on a tree run and accepted on a plain run. A change-list link inside a script is ignored.
- `ls` order checked on macOS under `C` and `en_US.UTF-8`: pre-order holds (F1 covers Linux).
- Red-green: the builder's RED at base `94f52975` is recorded for all three new tests. The REQ's stated RED (`bundle must be a physical direct child of ai-reports`) is the tree-list test's RED. Not repeated here, because the review is read-only.
- Test budget: REQ-704 added no runs to `report_second_opinion_test.go` (D-12 edits only). The new file takes about 2.5 s. The 26.49 s gate figure for the old file is not made worse by this REQ.
- The repository gate passed at `c5a8aa8d` (integrator's run, not rerun here).

### Suggested Additional Testing

- Deployment: unassessed. After the release, install into a consumer project and confirm the installed prompt carries `{{BASE_PAGE}}` and `{{CHANGE_LIST_REQUIREMENT}}`. Otherwise every directed run is refused with the new placeholder error.
- Live acceptance: unassessed. Run `ai-report second-opinion --backend codex --preset bold latest`, then `--preset visual <bundle>/index-v1c01-...html`, with a real editor. Confirm that the editor writes a resolving change list (D-08 turns a missing list into `fail`), that the judge passes with the iteration list on `index.html`, and that the parent link opens.
- Edge case: run the session form on a tree page with no `--preset` (header with only the iteration paragraph), and check that the page reads well.
- Environment: on Linux, `ls` a bundle with `v1c01` and `v1c01c01` pages (F1).

### Domain Review

Backend checklist: input is validated at the boundary. The page target must be a regular file, in the current bundle, and match an exact whole-basename pattern, so there is no traversal and no symlinked base. Every launcher-written value in metas, the header and the list is HTML-escaped. The one shared write (`index.html`) keeps the existing identity and bytes guard and the pass gate. Failure keeps the page (REQ-701 rule). No new state file. Tree state is in filenames and display metas only, as the Constraints require. Concurrency: see F5.

### Anti-bloat

Additions the Plan or brief did not name: 8, all kept.
- `secondOpinionTreeMarks` struct: keep. It groups the three display values so that `secondOpinionDecorate` does not grow three more string parameters. The hand-back lists it as named by the brief. It is not: neither the brief nor the Plan names it.
- `secondOpinionChangeList` function: keep. D-08's check, one caller, keeps the tail readable.
- `secondOpinionTreeList` function: keep. D-09 describes the behavior but not the function. The hand-back counts it as named.
- `secondOpinionTreeStart` / `secondOpinionTreeEnd` constants: keep. Used by both the decoration strip and the list builder.
- Test const `secondOpinionValidChangeList`: keep. Used 4 times across 3 tests.
- Test helper `secondOpinionFakeTreeEditor`: keep. 3 callers.
- Test helper `secondOpinionStubJudgePass`: keep. 2 callers. Borderline, but smaller than two inline stubs.
- `cmp` import: keep. Standard library.
- Named by the Plan (not counted): `Base`, `parentTreeID` (D-02 "one method"), `secondOpinionTreePagePattern`, `secondOpinionChangeListRequirement`, `secondOpinionChangeListKind`, `secondOpinionStripBlocks` (D-06 "one helper"), `TreeID`, `BasePage`, the `-model`/`-date` metas (D-06).
- No new flag, environment variable, state file, package or decorative test. The hand-back's list is complete apart from labeling `secondOpinionTreeMarks` and `secondOpinionTreeList` as brief-named.

### Scores (on the record, not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | 10/10 acceptance criteria and GREEN (1) to (6) delivered |
| Code Quality | 92% | Follows the file's masking and naming patterns. Small duplication (F7), duplicate findings (F6), narrow concurrent-allocation window (F5) |
| Test Adequacy | 88% | Real RED for all three new tests and every GREEN case pinned. The widened strip and the old-template refusal are unpinned (F4) |
| Scope | 100% | 8 of 15 declared paths touched. The 7 untouched are release and lesson files for finalization. Nothing undeclared |
| Risk | Low | Locale-dependent `ls` claim (F1). Concurrent-allocation window (F5) |
| Acceptance | Pass | Implementation and integration stages |

### Follow-ups created
None (7 findings report only)

## Review

**Overall: 95%** | 2026-10-21T10:46:04Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1 `ai-report.md:73`, `ai-report-guide.md:107`: "a plain `ls` lists the tree in order" holds for byte order (macOS verified) but plausibly not under glibc `en_US.UTF-8`, which ignores `-`; say "a byte-order listing (`LC_ALL=C ls`)" — impact-user-visible → report only
- F2 `ai-report.md:75`, `ai-report-guide.md:107,111`: plain-run-only sentences (numeric suffix, new banner, "direction header on a directed page") are left unqualified; a tree run gets a new id, no suffix, no banner, and a no-flag tree page also gets the header — impact-negligible → report only
- F3 `lessons-do-work-cli.md:95`: the REQ-703 bullet still says the meta strip misses reordered copies, which REQ-704 D-06 fixed; the integrator could append "(fixed in REQ-704)" during the link sweep — impact-negligible → report only
- F4 `report_second_opinion_tree_test.go`: no test pins the widened strip (copied content-first meta or tree block on a tree page) or the tree-run refusal of a template without the new placeholders (reviewer probe confirmed both); fixture writes ignore errors — impact-negligible → report only
- F5 `report_second_opinion.go:1426`: two simultaneous publishes on one base with different stems can allocate the same `cNN` (ReadDir to page-write window; `--all` is sequential) — impact-negligible → report only
- F6 (Nit) `report_second_opinion.go:1696`: repeated links to one missing anchor give duplicate findings — impact-negligible → report only
- F7 (Nit) `report_second_opinion.go:258,655,974,1625`: `cmp.Or(direction.Base, "index.html")` repeated four times and `base_page` built twice; one method would hold it — impact-negligible → report only

**Anti-bloat:** 8 unnamed additions, all keep: `secondOpinionTreeMarks`, `secondOpinionChangeList`, `secondOpinionTreeList`, `secondOpinionTreeStart`/`End`, test const `secondOpinionValidChangeList`, test helpers `secondOpinionFakeTreeEditor` and `secondOpinionStubJudgePass`, `cmp` import. The hand-back's list is complete but labels `secondOpinionTreeMarks` and `secondOpinionTreeList` as brief-named.
**Acceptance:** Pass. Implementation and integration stages: the four named tests and `-run SecondOpinion` (toolboxcommands, resultmodel) passed at `c5a8aa8d`; scratch-copy probes confirmed strip, idempotence, allocation and template refusal; the builder's RED at `94f52975` is recorded; the gate passed at `c5a8aa8d`.
**Restatement sweep:** redefined: published second-opinion names (directed and base-page runs carry a tree id, `index-v1c01-...`, no `-N`), positional target (accepts a page), `index.html` write (pass-gated banner on plain runs, iteration list on tree runs), prompt placeholder set (`{{BASE_PAGE}}`, `{{CHANGE_LIST_REQUIREMENT}}`) and the DIRECTION exception list (now with CHANGE LIST), run-row fields (`tree_id`, `base_page`), launcher meta set and strip (`-id`, `-base`, `-model`, `-date`; either attribute order). Inherited from REQ-703 (name shape, placeholder set, `--finalize` argument rule, catalog-walk rule, header before the backlink) and REQ-705 (bundle-only user presets, `preset_file` meaning, catalog walk no longer skips the presets folder). Grepped `skills/`, `docs/`, `_dev/`, `suite/` and READMEs, with CHANGELOG entries, `do-work/` records and `ai-reports/` pages left as history. Stale: `ai-report.md:73` and guide :107 `ls` claim (F1), `ai-report.md:75` and guide :107, :111 plain-run sentences (F2), `lessons-do-work-cli.md:95` (F3). Agreeing: the usage string, `ai-report.md:21,65-67,71,73`, guide :72-74, :107, :111 (tree parts), `ai-report-reference.md:148`, the prompt INPUT, CHANGE LIST and DIRECTION lines, `SKILL.md`, `result_model.go:719-721`, the shipped presets (no restatement), launcher test `_dev/tests/do-work-cli-launcher-behavior.sh:252` (an env check, not a fill). REQ-703 and REQ-705 both had a sweep line, so nothing is unread for the sweep.
**Suggested testing:** 4 items
**Follow-ups created:** None (7 findings report only)

*Reviewed by review-work action*
