---
id: REQ-704
title: 'Second opinions build a linked iteration tree on any base page'
status: claimed
created_at: 2026-10-20T19:59:08Z
user_request: UR-160
addendum_to: REQ-698
domain: backend
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: [REQ-703]
related: [REQ-703]
batch: UR-160
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/actions/ai-report-reference.md", "CHANGELOG.md", "skills/do-work/CHANGELOG.md", "VERSION", "skills/do-work/VERSION", "skills/do-work/actions/version.md", "skills/do-work/tools/do-work-cli/lessons-do-work-cli.md", "do-work/lessons-index.md"]
claimed_at: 2026-10-20T21:19:46Z
route: C
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-10-20T21:26:46Z
  basis:
    - Route C
    - 13-file write set
    - 3 subsystems involved
    - 10 acceptance criteria
planning_at: 2026-10-20T21:28:43Z
---
# Second Opinions Build A Linked Iteration Tree On Any Base Page
## What
A second opinion can build on any earlier second-opinion page, not only on `index.html`. Each published page gets a sortable tree id in its filename (`index-v1c01c02-...`), a launcher-owned header linking to its parent, and a "Changed in this iteration" list with `#rev-*` anchors. One launcher-owned list on the original `index.html` shows the whole tree.
## Why
The maintainer wants to stack directions ("bold", then "visual" on top of bold) and "keep them linked together so I can walk through all the iterations that were made". Today every second opinion rewrites only `index.html`, and pages carry no parent link or change record.
## Prior Implementation
- REQ-698 / REQ-700 / REQ-701 / REQ-702: see REQ-703's Prior Implementation. Key facts: `secondOpinionNames` builds `index-<cli>-<model>[-N].html`; `secondOpinionPreserved` already refuses any change to any pre-existing bundle file; `secondOpinionBanner` inserts a per-filename marker banner after `<body>` on `index.html` only after judge `pass`; `secondOpinionBacklink` requires the page to link `index.html`.
- REQ-703 (this batch, must land first): `--preset`, `--editor-direction`, the preset/direction meta tags and header block, and directed filenames `index-<cli>-<model>[-unverified]-<preset>[-<editor-slug>][-N].html`.
## Detailed Requirements
- D4 base selection: the positional target accepts `latest`, a bundle directory, or an HTML page inside a bundle; a page path gives both the bundle and the base. The base must be `index.html` or a tree page (`index-v1c...html`) in that bundle. Default base is `index.html`. Action prose: the user usually names the HTML file; when the user's choice of base is unclear, the agent asks with AskUserQuestion, best guess first plus a free-text entry.
- D5 editor input: the prompt names the base page as the page to rewrite and the original `index.html` plus the repository as fact sources, and tells the editor to read the base page's "Changed in this iteration" list.
- D6 tree ids: `index.html` is `v1`; a child of page P has id `P` + `cNN` with two-digit NN = next free sibling number; depth is unlimited (one-to-many). Filename: `index-<tree-id>-<cli>-<model>[-unverified]-<preset>[-<editor-slug>].html`; reply `<tree-id>-<cli>-<model>...-reply.md` on the same stem. A plain `ls` sorts the pages in tree (pre-order) order. The id is allocated by the launcher at publish time, never by the editor, so two runs on one base cannot collide.
- Which runs get a tree id: every run that has a preset, a direction, or a non-`index.html` base. A plain run (no preset, no direction, base `index.html`) keeps today's name.
- D8 change marks: the prompt requires a "Changed in this iteration" list at the top of the page; each item links to `#rev-<n>` on the changed section, and that section carries a visible `changed` badge. The launcher checks every listed `#rev-*` target exists; a failure is a judge finding (page kept and published, banner withheld), never a deletion.
- D9 walk-through: one launcher-owned marked block on `index.html` lists the whole tree (tree id, preset, model, date, link), rebuilt from a scan of `index-v1c*.html` filenames plus their meta tags on every publish, and written only after judge `pass` like today's banner. Each tree page's launcher header (from REQ-703) gains its tree id and a link to its parent page; meta tags gain `ai-report-second-opinion-id` and `ai-report-second-opinion-base`. Tree pages do not use the per-filename banner.
- D10: the existing non-tree second-opinion pages (`index-claude-claude-opus-5-5.html`, `index-codex-gpt-6-astra.html` in the REQ-632 bundle; `index-agy-gemini-4-argon-high.html` in the REQ-698 bundle) are ignored by the tree code and left untouched.
- Preservation: earlier tree pages and the base are never rewritten; only the launcher's own marked block on `index.html` changes.
- Result row gains `tree_id` and `base_page` (omitempty).
- Docs: `skills/do-work-toolbox/actions/ai-report.md` § Second-opinion form gets the page-path target, the ask-on-unclear-base rule, and one line on tree names and the index list; `skills/do-work-toolbox/docs/ai-report-guide.md` matches.
## Constraints
- Depends on REQ-703 (same Go file, same naming function).
- No separate state file: tree state lives in filenames and meta tags only.
- A plain run keeps today's behavior; REQ-701's keep-finished-page rule stays.
- Shipped change: a release per `_dev/primes/prime-releases.md`, with a descriptive changelog title (for example "Second Opinions Build a Linked Iteration Tree").
## Dependencies
Depends on REQ-703 (presets and free-text direction).
## Builder Guidance
D4-D10 are user-confirmed and firm, including two-digit children and the tree list living only on `index.html`. Marker names, the exact list layout, and the badge style are builder latitude. Keep it simple: derive the tree from filenames; meta tags only add display fields.
## Red-Green Proof
**RED prompt/case:** after REQ-703, run `ai-report second-opinion --preset visual <bundle>/index-<cli>-<model>-bold-<slug>.html`: the page target is refused ("bundle must be a physical direct child of ai-reports"), and pages carry no tree id, no parent link and no tree list on `index.html`.
**Why RED now:** the launcher always edits `index.html`, `secondOpinionBundle` accepts only bundle directories, and `secondOpinionNames` has no tree id.
**GREEN when:** with fake backends: (1) `--preset bold` on `index.html` publishes `index-v1c01-<cli>-<model>-bold-<slug>.html`; (2) a second run on `index.html` publishes `index-v1c02-...`; (3) `--preset visual` targeting the `v1c01` page publishes `index-v1c01c01-...-visual-...html`, whose header links `index-v1c01-...html` and whose "Changed in this iteration" list links resolve to badged `#rev-*` sections; (4) `index.html`'s tree list shows `v1c01`, `v1c01c01`, `v1c02` in that order, the same order `ls` gives; (5) a page whose list names a missing `#rev-2` is kept and gets a judge finding; (6) a plain run still publishes `index-<cli>-<model>.html`.
**Validation:** User confirmed.
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 22858 tokens, `slugged: partial` so no targeted form; matches exact filesystem-basename authority and native backend identity. Read its `exact-basename-authority` and `writer-anchored-marker-position` bullets if budget allows.
- `_dev/primes/lessons-action-files.md` — 8398 tokens, `slugged: partial`; matches the action-file and downstream-reader changes.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
## Full Context
See `do-work/user-requests/UR-160/input.md` for complete verbatim input.

*Source: direction wins all, because I'm the director and I want to steer reports as well while keeping them linked togheter so I can walk through all the iterations that where made (plus question rounds 2-4 in UR-160)*

## Triage

**Route: C** - Complex

**Reasoning:** The user-confirmed decisions (D4-D10) fix the outcome, but the change crosses the positional target contract, the prepare record, the shared publication tail, the one name builder, the launcher header block, a new launcher block on `index.html`, a new launcher check that feeds the judge verdict, the typed result row, the shipped prompt and two prose files. It also changes REQ-703's directed names (a directed run now gets a tree id), so REQ-703's test pins move. The pre-dispatch agent decides the mechanism so the builder gets one design.

**Planning:** Required

*No `## Open Questions` section exists; Step 3.5 has nothing to resolve. The builder-latitude points (marker names, list layout, badge style) and the CHANGE LIST precedence the coordinator asked about are decided in `## Plan` → Decisions as D-01 onwards.*

<!-- D-XX counter: last used D-12. Next decision: D-13. -->

## Plan

**Decided mechanism (pre-dispatch, do not reopen):** the base page joins the direction value, so it crosses the prepare record and reaches the prompt fill and the publication tail with no new parameters. A run is a tree run when it has a preset, a free-text direction, or a base other than `index.html`. On a tree run the one name builder allocates the child id from the bundle's filenames at publish time, the one decoration function adds the tree marks and the parent link, a new launcher check reads the page's `#rev-*` change list after the judge, and a judge `pass` rebuilds one launcher block on `index.html` from the tree filenames instead of inserting the per-filename banner. A plain run (no preset, no direction, base `index.html`) takes today's path byte for byte. Function names below are current main (REQ-703 released as 0.306.6); REQ-705 (bundle-local presets) lands first and changes only preset lookup.

### Decisions
- **D-01** (DECIDE & STATE): page target. A positional target ending in `.html` is a page: its directory goes through the unchanged `secondOpinionBundle`, and the result must equal that directory, else usage error `base page must sit in a current report bundle` (the same rule `--finalize` uses for its staging page, so a page in a superseded bundle is never silently moved to the successor). The page name must be `index.html` or a tree page name (D-03) and a regular file (`os.Lstat`), else usage error `base page must be index.html or a tree page (index-v1c...html) in the bundle`. The three existing non-tree pages are refused as bases (UR-160 D10). A directory or `latest` keeps base `index.html`. Bundle and base resolve before the preset lookup, because REQ-705 looks presets up inside the bundle.
- **D-02** (DECIDE & STATE): `secondOpinionDirection` gains `Base string` (`json:"base,omitempty"`, empty means `index.html`). The person directing the pass chose it, it must survive prepare/finalize exactly like the preset, and `secondOpinionRecord.Direction` (`omitzero`) already carries the struct, so finalize restores the base with no new record field. One method returns the parent tree id: empty on a plain run, `v1` when the base is `index.html`, otherwise the base page's id. A record written before this change reads as base `index.html`.
- **D-03** (DECIDE & STATE): tree id grammar, one exact pattern on the whole basename, never trimmed (lesson `exact-basename-authority`): `^index-(v1(?:c[0-9]{2})+)-[A-Za-z0-9][A-Za-z0-9._-]*\.html$`. The id is capture 1. `c` followed by a letter (codex, claude) cannot extend the id, so the CLI segment never merges into it. The tree structure (parent, depth, order) comes from these filenames only; metas are display fields.
- **D-04** (DECIDE & STATE): child allocation lives in `secondOpinionNames`, which gains the parent id and also returns the allocated id. NN = 1 + the highest two-digit number used directly under the parent by any tree page name in the bundle root, counting descendants' ids too, so a deleted page's number is never reused while its children remain. A taken page or reply name bumps NN. NN above 99 is a failure `no free two-digit child id under <parent>`. Tree names: page `index-<id>-<backend>-<model>[-<variation>].html`, reply `<id>-<backend>-<model>[-<variation>]-reply.md`, no `-N` suffix. A plain run keeps today's names and `-N` loop. The unverified form needs no special case (`publishedModel` is already `<model>-unverified`).
- **D-05** (DECIDE & STATE): publication order in `secondOpinionPublication`: strip (zero marks) → backlink check (skipped when unverified, as today) → variation slug → `secondOpinionNames` (allocates the id) → decorate with the direction and the tree marks → publish page and reply → unverified return → judge → change-list check (D-08) → pass-gated `index.html` write (D-09 on tree runs, today's banner on plain runs). Decoration moves after naming because the header needs the id.
- **D-06** (DECIDE & STATE): tree marks. `secondOpinionDecorate` gains the tree id plus two display values, the model label `<backend> <publishedModel>` and the UTC publish date `YYYY-MM-DD`; an empty id means no tree marks. On a tree page it adds metas `ai-report-second-opinion-id`, `ai-report-second-opinion-base` (parent filename), `ai-report-second-opinion-model`, `ai-report-second-opinion-date` beside REQ-703's preset and direction metas, and the header block (shown when there is a tree id or a direction slot) gains a first paragraph `Iteration <id>, built on <a href="<parent file>"><parent file></a>`. Every value is HTML-escaped. The strip half widens in two ways: it removes every `<meta>` whose name starts with `ai-report-second-opinion-`, in either attribute order (reuse `reportMetaTag` and `reportMetaAttribute` from `report_index.go`), and it removes copied `<!-- second-opinion-tree -->` blocks as well as header blocks. Reason: the tree list reads metas first-wins, so a copy the editor took from its base page must never survive (closes REQ-703 review F5 for these fields), and a copied tree block ahead of the backlink would fail the backlink check. The marker-strip loop becomes one helper shared by decoration and the tree list. Decorating twice still gives the same bytes.
- **D-07** (DECIDE & STATE): prompt. INPUT names `{{BASE_PAGE}}` (the bare base filename) as the page to rewrite, keeps `index.html` and the repository as the fact sources, and tells the editor that a base other than `index.html` opens with a "Changed in this iteration" list to read. A new `CHANGE LIST. {{CHANGE_LIST_REQUIREMENT}}` paragraph follows BACKLINK; its text is a Go constant beside `secondOpinionBacklinkRequirement`, filled on tree runs and replaced by `None for this pass.` on plain runs. The constant says: right after the backlink, a list titled "Changed in this iteration"; each item names one change from the base and links to `#rev-<n>`; each changed section has `id="rev-<n>"` and a visible `changed` badge; every link must resolve because the launcher checks it. The DIRECTION paragraph's exception list becomes OUTPUT, BACKLINK, CHANGE LIST, BOUNDARY and FINAL REPLY: the direction may not override the change list, because the launcher checks it (answers the coordinator's question; same reason as REQ-703's D-17). `secondOpinionPrompt` refuses a tree run whose template lacks `{{BASE_PAGE}}` or `{{CHANGE_LIST_REQUIREMENT}}`; a plain run with an older template still works, as REQ-703 did for `{{EDITOR_DIRECTION}}`.
- **D-08** (DECIDE & STATE): change-list check, tree runs only, after the judge (an unverified page gets no judge and no check, as today). On the visible surface (`secondOpinionVisibleHTML`), collect every `href="#rev-…"` target and every real `id="rev-…"` attribute (an `id` attribute preceded by whitespace, so `data-id` does not count). No `#rev-` link at all is one finding; each link without a matching id is one finding naming the target. One finding kind, a new constant `AI-REPORT-SECOND-OPINION-CHANGE-LIST`, with `target` set to the anchor. Any finding is appended to `judge_findings` and turns verdict `pass` or `skipped` into `fail` (`error` stays `error`), so the page and reply stay published and `index.html` is untouched (REQ-701 keep-finished-page rule). Requiring at least one link is a small step past D8's wording: without it an editor that skips the list passes a check that exists to keep iterations walkable.
- **D-09** (DECIDE & STATE): the tree list. On a tree run with judge `pass`, after the same outside-bundle and original identity/bytes checks the banner uses today, the launcher rewrites `index.html` with one block between `<!-- second-opinion-tree -->` and `<!-- /second-opinion-tree -->`, found as real comments on the masked surface: any old block is removed and a fresh one goes right after the opening body tag. The block is rebuilt on every passing tree publish from `os.ReadDir` of the bundle root: every regular file matching D-03, sorted by name, which is pre-order (`-` sorts before `c`). Each item shows the id (indented by depth), preset (meta, or `plain`), model and date (metas, read with `readReportMetas(secondOpinionVisibleHTML(page))`), and a link to the page; all values escaped; a page missing a meta shows the field empty. Every tree page on disk is listed, including unverified and judge-failed ones, so no child is orphaned. Tree pages never get the per-filename banner; plain runs keep it; existing banners stay. The row's `banner_state` reuses `inserted` / `already-present` for this block (no new state value).
- **D-10** (DECIDE & STATE): result row. `AIReportSecondOpinionRun` gains `TreeID` (`json:"tree_id,omitempty"`, set at publish) and `BasePage` (`json:"base_page,omitempty"`, the repository-relative base path, set on every tree-run row including prepare). Plain rows are unchanged.
- **D-11** (DECIDE & STATE): docs. `ai-report.md` § Second-opinion form: the three forms take `[<bundle-dir>|<page.html>|latest]`; the trigger bullet maps "second opinion on <page>.html" to the page path as the target; when it is unclear which page is the base, load `../../do-work/crew-members/clear-questions.md` and ask with the environment's ask-user prompt, best guess first plus a free-text entry (shipped files say "ask-user prompt", as `present-work.md:119` and `prime.md:70` do, never a tool name); one line on tree names, the parent header, the change list and the `index.html` list; the directed-name example gains its tree id; rows gain `tree_id` and `base_page`. `ai-report-guide.md` matches (command lines, § A Second Opinion). `ai-report-reference.md:148` says "a pass-gated banner" and becomes "a pass-gated banner or iteration list". The `## Second-opinion form` heading text stays (anchor `#second-opinion-form`). `SKILL.md` is not touched.
- **D-12** (DECIDE & STATE): REQ-703's `TestSecondOpinionDirectedRunNamesAndMarksPage` pins directed names without a tree id. Every directed run is a tree run now (UR-160 D6), so its expected page and reply names gain `v1c01`, `v1c02`, `v1c03` in run order; every other assertion stays, and its `secondOpinionDecorate` unit calls pass empty tree marks. That is a contract change the REQ asks for, not a weakened test.

### Files to modify
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify): usage string, target parse, base in the direction value, tree pattern, names, decoration, prompt fill, change-list check, tree list, row fields.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify): fixture template placeholders, three new tests, REQ-703 directed-test name pins.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify): two omitempty fields.
- `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` (modify): INPUT, CHANGE LIST, DIRECTION exception list.
- `skills/do-work-toolbox/actions/ai-report.md`, `skills/do-work-toolbox/docs/ai-report-guide.md`, `skills/do-work-toolbox/actions/ai-report-reference.md` (modify).
- Release files (integrator only): both CHANGELOGs, both VERSION files, `skills/do-work/actions/version.md`; lesson satellite and index only with a lesson.

### Order of changes
1. **RED: targets, names, parent header.** `TestSecondOpinionTreeTargetsNamesAndParentHeader` fails at base on `bundle must be a physical direct child of ai-reports` for a page target (the REQ's RED).
2. **GREEN: D-01 to D-06 and D-10**, plus the fixture template placeholders and D-12's name pins.
3. **RED then GREEN: change list.** `TestSecondOpinionChangeListTargetsAreChecked`, then D-07's constant and fill and D-08's check.
4. **RED then GREEN: tree list.** `TestSecondOpinionTreeListOnIndexFollowsTreeOrder`, then D-09.
5. **Shipped text and prose** (D-07 prompt, D-11 docs); propose the changelog entry.

### Testing approach
- Focused, from `skills/do-work/tools/do-work-cli`: `go test -count=1 -v -run '^(TestSecondOpinionTreeTargetsNamesAndParentHeader|TestSecondOpinionChangeListTargetsAreChecked|TestSecondOpinionTreeListOnIndexFollowsTreeOrder|TestSecondOpinionDirectedRunNamesAndMarksPage)$' ./internal/toolboxcommands`.
- Must stay green: every existing `TestSecondOpinion*` test (plain names, unverified names, session prepare/finalize, banner, judge gate, capture) and `resultmodel.TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence`.
- The judge is stubbed through `secondOpinionJudge` where a test needs `pass` (as `TestSecondOpinionJudgeVerdictsGateOriginal` does); with no browser it is `skipped`.
- Package checks: `go vet` on `./internal/toolboxcommands ./internal/resultmodel`, `gofmt -l ./internal`.

### Risks
- Decoration must move after `secondOpinionNames`; the judge must still see the decorated bytes it publishes.
- The fake editor's default page has no change list, so a tree-run test that expects `pass` must write one, or D-08 turns the verdict to `fail`.
- Child allocation reads the bundle at publish time; a session finalize must allocate then, not at prepare.
- The tree list strip and insert must run on the masked surface; a marker inside a script string must survive.
- Closing UR-160 moves the REQ-703 record into `do-work/archive/UR-160/`, which breaks the flat-path link at `lessons-do-work-cli.md:95` (integrator, lesson `canonical-link-outlives-its-target`).

### Consumer field contract (per `runs[]` row)
- Identity and naming: `backend`, `requested_model`, `actual_model`, `page`, `reply` (tree names per D-04).
- Tree provenance: `tree_id` (allocated id, publish rows), `base_page` (repository-relative base path, every tree-run row), plus REQ-703's `preset`, `preset_file`, `editor_direction`; all omitempty, absent on plain runs.
- Outcome: `judge_verdict` and `judge_findings` carry D-08's change-list findings; `banner_state` reports the `index.html` write (banner or tree list).

*Generated by the pre-dispatch agent (planning done directly; the user approved the design as D4-D10 in UR-160)*

### Plan validation (pre-dispatch)
- Requirement coverage: D4 base selection → D-01, D-02, task 1-2, D-11 prose; D5 editor input → D-07; D6 tree ids and names → D-03, D-04; which runs get an id → D-02; D8 change marks and check → D-07, D-08; D9 tree list, header, metas → D-06, D-09; D10 old pages ignored → D-01, D-03; preservation → unchanged checks, D-09; row fields → D-10; docs → D-11. No uncovered requirement.
- Orphan tasks: none. D-12 follows from D6; the `ai-report-reference.md` line follows from D-09.
- Scope sanity: **Plan has 5 tasks — quality degrades past 3.** Not split: tasks 1-4 are one serialized slice of one Go file and its test, task 5 is prose that ships with it, and UR-160 is already split into REQ-703 and this REQ.
- Consumer field contract: present above.

## Exploration

Read directly by the pre-dispatch agent at main `11780ac0`: `report_second_opinion.go` (all 1579 lines), the fixture helpers and the session, judge-gate, banner and REQ-703 tests in `report_second_opinion_test.go`, `readReportMetas` and its regexes in `report_index.go`, the judge codes in `report_judge.go`, `AIReportSecondOpinionRun` in `resultmodel/result_model.go`, the shipped prompt, `ai-report.md` § Second-opinion form, `ai-report-guide.md` § Input and § A Second Opinion, `ai-report-reference.md:148`, the archived REQ-703 record (Plan, Exploration, Review, D-17), its hand-back, and REQ-705's answered question.

**Target and bundle.** `secondOpinionParse` (`report_second_opinion.go:264`) stores the one positional in `options.bundle` (`:298-300`). `secondOpinionBundle` (`:327`) accepts only a direct child of `ai-reports` and follows the supersession chain, which is why a page path is refused today (the REQ's RED). The finalize branch (`:132-142`) already shows the page-in-bundle pattern: resolve `filepath.Dir`, then require the resolved bundle to equal it (D-01). The handler reads the template and resolves presets only on the non-finalize path (`:144-188`); REQ-705 moves preset lookup into the bundle, so bundle and base resolve first.

**Direction value and record.** `secondOpinionDirection` (`:95`) with `combinedText` and `filenameSlot`; the handler builds it at `:131`, every row copies its fields at `:215`, `secondOpinionRecord.Direction` is `omitzero` (`:398`), and finalize restores it at `:606`. A `Base` field rides all of these unchanged (D-02).

**Prompt fill.** `secondOpinionPrompt` (`:891`) fills six placeholders and refuses a directed run on a template without `{{EDITOR_DIRECTION}}`; `secondOpinionBacklinkRequirement` (`:85`) is the pattern for a Go-owned requirement text (D-07). Callers: `secondOpinionPrepare` (`:566`, absolute stage path) and `secondOpinionRun` (`:734`, bare stage name). The shipped prompt's INPUT is line 5, BACKLINK line 11, DIRECTION line 17 (exception list from REQ-703 D-17).

**Publication tail.** `secondOpinionPublication` (`:783`): strip at `:787`, backlink at `:793`, variation at `:798`, decorate at `:805`, `secondOpinionNames` at `:806`, publish `:811-830`, unverified return `:838`, judge `:843`, then outside check, original identity and bytes (`:852-868`) and `secondOpinionBanner` (`:869`), written with `rootedPublishFile(..., true)` (`:880`). D-05 swaps decorate and names; D-08 sits after `:843`; D-09 replaces the banner call on tree runs and reuses `:852-868` and `:880`.

**Names.** `secondOpinionNames` (`:1339`) builds the stem and loops `-N` while the page or reply exists (`os.Lstat`). It is the only name builder (REQ-703 D-08), so D-04 lives there.

**Marks.** `secondOpinionDecorate` (`:1463`) strips header blocks in a loop over real comments on the raw-text-masked surface (`:1469-1490`), strips the two metas with a name-first regex (`:1491-1495`), and inserts metas before the real `</head>` and the block after the body tag. `secondOpinionBanner` (`:1428`) shows the same masking. `readReportMetas` (`report_index.go:204`) reads metas in either attribute order, first wins, through `reportMetaTag` and `reportMetaAttribute` (`:47-48`); D-06 reuses both for the widened strip and D-09 for the list.

**Judge.** `secondOpinionJudge` (`:38`) is a test seam; `secondOpinionRender` copies the bundle and renders the page as `index.html`. Finding kinds are constants such as `judgeCodeBrokenLink` (`report_judge.go:31`); D-08 adds one constant in the same style.

**Tests and fixtures.** `secondOpinionFixture` (`report_second_opinion_test.go:21`) writes a temp template; it must gain `{{BASE_PAGE}}` and `{{CHANGE_LIST_REQUIREMENT}}` lines so tree runs see their fill. `secondOpinionFakeCodex` (`:52`) writes a fixed page and runs `extra` after it, which a test can use to rewrite the page with a change list. `TestSecondOpinionJudgeVerdictsGateOriginal` (`:202`) shows the judge stub. `TestSecondOpinionDirectedRunNamesAndMarksPage` (`:880`) pins directed names without a tree id (D-12). `secondOpinionFixturePresets` (`:828`) writes shipped presets beside the template; REQ-705 keeps the shipped lookup.

**Prose sites.** `ai-report.md:21` (trigger bullet), `:65-67` (three forms), `:73` (names, marks, row fields). `ai-report-guide.md:72-74` (command lines), `:107` and `:111` (§ A Second Opinion). `ai-report-reference.md:148` (banner sentence, D-11). Shipped files never name AskUserQuestion; the house phrase is "the environment's ask-user prompt" (`present-work.md:119`, `prime.md:70`).

**Required-lessons consult (claim time).** `required_lessons` stays `_dev/primes/lessons-releases.md` (666 tokens, read whole): `canonical-link-outlives-its-target` applies to this REQ's integrator, because closing UR-160 moves the REQ-703 record and the new REQ-703 bullet links its flat path (`lessons-do-work-cli.md:95`). The index lists `lessons-do-work-cli.md` at 22858 tokens and `lessons-action-files.md` at 8398 tokens, both `slugged: partial`, so the Dropped section stands (cost refreshed). Bullets read directly as relevant: `exact-basename-authority` (REQ-485: the tree id is parsed from the exact basename, never trimmed, D-03), `writer-anchored-marker-position` (REQ-528, REQ-703: launcher blocks are stripped before the backlink check and read only as real comments; REQ-703's note that the meta strip misses reordered copies is why D-06 widens it), `native-backend-identity` (REQ-698, REQ-701: the unverified name gains the tree id and nothing else changes), `session-is-the-backend` (REQ-700: the base crosses prepare/finalize in the record).

*Generated by the pre-dispatch agent's direct reading*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify) — page target, base in the direction value, tree id grammar, child allocation in the name builder, tree marks and widened strip, prompt fill, change-list check, tree list on index.html, row fields
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify) — fixture placeholders, three new tests, REQ-703 directed-test name pins
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — tree_id and base_page on the run row
- `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` (modify) — base page input, change-list paragraph, direction exception list
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — page target, ask-on-unclear-base rule, tree names, row fields
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — command lines and the plain-words tree description
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modify) — banner sentence covers the iteration list
- `CHANGELOG.md` (modify, integrator's release) — the entry
- `skills/do-work/CHANGELOG.md` (modify, integrator's release) — the byte-identical mirror
- `VERSION` (modify, integrator's release) — the version bump
- `skills/do-work/VERSION` (modify, integrator's release) — the version mirror
- `skills/do-work/actions/version.md` (modify, integrator's release) — the Current version line
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modify, integrator) — the REQ-703 link sweep after UR-160 closes, plus a lesson bullet only if the hand-back proposes one
- `do-work/lessons-index.md` (modify, integrator, only with a lesson) — token refresh for that satellite

**Files I will NOT touch:** `skills/do-work-toolbox/SKILL.md` (no new trigger family), `report_index.go` and its test (REQ-705 owns the catalog-skip removal; this REQ only calls `readReportMetas` and its two regexes), `report_judge.go` (the change-list check is a launcher check, not a render finding), the codex, gemini, agy and claude adapters and the capture supervisor, `skills/do-work/tools/do-work-cli.sh`, REQ-705's preset files, and the three existing second-opinion pages under `ai-reports/` (UR-160 D10).

**Acceptance criteria (restated from REQ):**
- [ ] The positional target accepts `latest`, a bundle directory, or an HTML page inside a bundle; a page gives both bundle and base; the base must be `index.html` or a tree page; default base `index.html`; the action prose says the user usually names the HTML file and the agent asks (best guess first plus free text) when the base is unclear.
- [ ] The prompt names the base page as the page to rewrite, `index.html` and the repository as fact sources, and tells the editor to read the base page's "Changed in this iteration" list.
- [ ] Tree ids: `index.html` is `v1`; a child of P is P + `cNN`, two digits, next free sibling number, unlimited depth; page `index-<tree-id>-<cli>-<model>[-unverified]-<preset>[-<editor-slug>].html`, reply on the same stem; a plain sort gives pre-order; the launcher allocates the id at publish time.
- [ ] Every run with a preset, a direction or a non-`index.html` base gets a tree id; a plain run keeps today's name.
- [ ] The prompt requires a "Changed in this iteration" list linking `#rev-<n>` on badged sections; the launcher checks every listed target; a failure is a judge finding, the page is kept and published, and no banner or list is written.
- [ ] One launcher-owned block on `index.html` lists the whole tree (id, preset, model, date, link), rebuilt from the tree filenames plus metas on every publish, written only after judge `pass`; tree pages carry their id and a parent link in the launcher header and metas `ai-report-second-opinion-id` and `ai-report-second-opinion-base`; tree pages get no per-filename banner.
- [ ] The existing non-tree second-opinion pages are ignored by the tree code and left untouched.
- [ ] Earlier tree pages and the base are never rewritten; only the launcher block on `index.html` changes.
- [ ] The result row gains `tree_id` and `base_page` (omitempty).
- [ ] `ai-report.md` and `ai-report-guide.md` document the page target, the ask-on-unclear-base rule, tree names and the index list; the change ships as one release with a descriptive changelog title.
