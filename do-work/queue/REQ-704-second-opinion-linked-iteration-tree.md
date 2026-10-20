---
id: REQ-704
title: 'Second opinions build a linked iteration tree on any base page'
status: pending
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
write_set: ["skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md"]
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
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 22697 tokens, `slugged: partial` so no targeted form; matches exact filesystem-basename authority and native backend identity. Read its `exact-basename-authority` and `writer-anchored-marker-position` bullets if budget allows.
- `_dev/primes/lessons-action-files.md` — 8398 tokens, `slugged: partial`; matches the action-file and downstream-reader changes.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
## Full Context
See `do-work/user-requests/UR-160/input.md` for complete verbatim input.

*Source: direction wins all, because I'm the director and I want to steer reports as well while keeping them linked togheter so I can walk through all the iterations that where made (plus question rounds 2-4 in UR-160)*
