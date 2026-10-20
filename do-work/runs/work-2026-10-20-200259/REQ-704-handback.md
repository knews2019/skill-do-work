# REQ-704 hand-back (second opinions build a linked iteration tree on any base page)

- Branch: `worktree-agent-REQ-704-second-opinion-linked-iteration-tree`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-704-second-opinion-linked-iteration-tree`
- Base commit: `94f52975` (main after REQ-705 0.306.7 and the coordinator note)
- Commits: `50555523` `[REQ-704] build a linked second-opinion iteration tree` (one commit, branch tip)
- Nothing under `do-work/`, `ai-reports/`, changelogs, VERSION or mirrors was touched. `report_index.go` was not touched (only `readReportMetas`, `reportMetaTag`, `reportMetaAttribute` are called).

## File manifest

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified): usage string takes `<page.html>`; page target resolves bundle and base before the preset lookup (D-01); `secondOpinionDirection.Base` + `parentTreeID()` (D-02); `secondOpinionTreePagePattern` (D-03); `secondOpinionNames` allocates the child id and returns it (D-04); publication order strip → backlink → variation → names → decorate → publish → unverified return → judge → change-list check → banner or iteration list (D-05); `secondOpinionTreeMarks` + widened strip + shared `secondOpinionStripBlocks` (D-06); `{{BASE_PAGE}}` and `{{CHANGE_LIST_REQUIREMENT}}` fills with `secondOpinionChangeListRequirement` (D-07); `secondOpinionChangeList` + `secondOpinionChangeListKind` (D-08); `secondOpinionTreeList` (D-09); `TreeID` / `BasePage` set on rows, including prepare and finalize (D-10).
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified): `TreeID` (`tree_id,omitempty`) and `BasePage` (`base_page,omitempty`).
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified): fixture template gains `Base: {{BASE_PAGE}}` and `Change list: {{CHANGE_LIST_REQUIREMENT}}`; D-12 pins (`v1c01`, `v1c02`, `v1c03` in the directed names; `secondOpinionDecorate` unit calls pass `secondOpinionTreeMarks{}`).
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_tree_test.go` (new, D-13): the three REQ-704 tests plus two small shared helpers.
- `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` (modified): INPUT names `{{BASE_PAGE}}`, keeps `index.html` and the repository as fact sources, says a non-index base opens with a "Changed in this iteration" list to read first; new `CHANGE LIST. {{CHANGE_LIST_REQUIREMENT}}` paragraph after BACKLINK; DIRECTION exceptions are OUTPUT, BACKLINK, CHANGE LIST, BOUNDARY and FINAL REPLY.
- `skills/do-work-toolbox/actions/ai-report.md` (modified): three forms take `[<bundle-dir>|<page.html>|latest]`; trigger bullet maps "second opinion on <page>.html" to the page target; target paragraph gives the page-target rule and the ask-on-unclear-base rule (clear-questions + "the environment's ask-user prompt", best guess first plus free text); names paragraph: tree ids, `index-v1c01-codex-gpt-5-bold-add-a-timeline.html`, parent header, tree metas, change list and its finding, the `index.html` iteration list, row fields `tree_id` and `base_page`. Heading `## Second-opinion form` unchanged.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified): three command lines take `<page.html>`; § A Second Opinion explains the tree in plain words (`v1`, `v1c01`, `v1c01c01`, names, parent link, change list), the change list joins the direction's exception list, and the pass-gated link paragraph names the iteration list.
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified): "a pass-gated banner or iteration list".

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Follow D-01 to D-13 as decided. Read: builder brief incl. coordinator note, REQ-704 (all sections), REQ-703 hand-back seams, general.md, coding-guardrails.md, shared-principles.md, testing.md, backend.md, communication-style.md, prime-do-work-cli.md (Traps), prime-action-files.md, prime-releases.md not re-read beyond the brief (release is the integrator's), lessons-releases.md (whole), and the five lesson families below. Write the three new tests first in a new file (D-13) against base code, run RED, then implement in one Go file plus the row fields, then prose. No lesson contradicted the plan; `exact-basename-authority` drives the whole-basename tree pattern, `writer-anchored-marker-position` drives stripping the copied header, tree block and every `ai-report-second-opinion-*` meta (either attribute order) on the masked surface before the backlink check.
- [x] **[APPLY]:** Implemented as planned, scope limited to the write boundary. One plan change during APPLY: I first added a finalize-time check that the prepare record's `base` is a tree page name, then removed it (the record is launcher-written, and a bad value only yields a plain-run page with no tree marks; not earned defense).
- [x] **[UNIFY]:** `git diff --stat 94f52975 HEAD`:

```
 .../do-work-toolbox/actions/ai-report-reference.md |   2 +-
 skills/do-work-toolbox/actions/ai-report.md        |  12 +-
 skills/do-work-toolbox/docs/ai-report-guide.md     |  10 +-
 .../docs/ai-report-second-opinion-prompt.md        |   6 +-
 .../internal/resultmodel/result_model.go           |   3 +
 .../toolboxcommands/report_second_opinion.go       | 266 +++++++++++++++++----
 .../toolboxcommands/report_second_opinion_test.go  |  18 +-
 .../report_second_opinion_tree_test.go             | 215 +++++++++++++++++
 8 files changed, 468 insertions(+), 64 deletions(-)
```

Checks (load average 1.65 / 2.24 / 3.38 at the end):

| Check | Exit | Wall |
|---|---|---|
| `bash .../REQ-704-probe.sh` from worktree root | 0 | 3.3 s |
| `gofmt -l skills/do-work/tools/do-work-cli/internal` | 0, no output | <1 s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands ./internal/resultmodel` | 0 | ~2 s |
| `git diff --check` (tree and `94f52975..HEAD`) | 0 | <1 s |
| `go test -count=1 -run 'SecondOpinion' ./internal/toolboxcommands` | 0 (ok 16.9 s) | 17.3 s |
| `go test -count=1 -run TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence ./internal/resultmodel` | 0 | 0.4 s |

Files checked by reading the diff: all eight above. No debug artifacts; a temporary view test (`zz_scratch_view_test.go`) used to eyeball the rendered `index.html` and child page was deleted before commit and never staged.

## Proof record

RED at base `94f52975` (the three new tests written first, fixture placeholders not yet added, no implementation):

- `TestSecondOpinionTreeTargetsNamesAndParentHeader`: `run on index.html, want v1c01: ... page:ai-reports/2026-10-15_1200-example/index-codex-actual-model-1-bold.html` (directed name without a tree id; first assertion to fail).
- `TestSecondOpinionChangeListTargetsAreChecked`: `v1c01: change list not checked, want findings [#rev-2]: ... banner_state:inserted judge_findings:[] judge_verdict:pass ...` (missing `#rev-2` passed and got a banner).
- `TestSecondOpinionTreeListOnIndexFollowsTreeOrder`: `missing typed backend summary: ... "observed_evidence":["bundle must be a physical direct child of ai-reports"]` (the page target refused; this is the REQ's stated RED).

GREEN after implementation:

- Four probe tests (`TestSecondOpinionTreeTargetsNamesAndParentHeader`, `TestSecondOpinionChangeListTargetsAreChecked`, `TestSecondOpinionTreeListOnIndexFollowsTreeOrder`, `TestSecondOpinionDirectedRunNamesAndMarksPage`): ok, 3.2 s package time, first try.
- `report_second_opinion_test.go` alone (all 27 top-level tests, `-v`): 27 PASS, 0 FAIL, 0 SKIP; ok 13.1 s (wall 13.3 s).
- `report_second_opinion_tree_test.go` alone (3 tests, `-v`): 3 PASS, 0 FAIL, 0 SKIP; ok 2.5 s (wall 2.7 s).
- `-run SecondOpinion` package selection: ok 16.9 s; covers plain names, unverified names, session prepare/finalize, banner, judge gate, capture, and REQ-705's preset tests.
- resultmodel `TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence`: PASS.
- Coordinator's extra assertion: in the tree-target test, a bundle preset (`<bundle>/second-opinion-presets/visual.md`) resolves when the target is the `v1c01` page path (`preset_file` is the bundle path and the prompt carries its text).

## Decisions

- **D-14** (DECIDE & STATE): change-list constant wording: `Right after the backlink, add a list titled "Changed in this iteration". Each item names one change you made to the base page and links to #rev-<n> (rev-1, rev-2 and so on). Give each changed section the matching id="rev-<n>" and a visible "changed" badge. Every link must resolve to a section on this page, because the launcher checks it.`
- **D-15** (DECIDE & STATE): prompt INPUT wording: `The page to rewrite is {{BASE_PAGE}}. Read it and every asset it references before you start, and read index.html too when it is a different page. index.html and the repository are the sources of fact. A page to rewrite other than index.html is an earlier iteration that opens with a "Changed in this iteration" list: read that list first, so you know what the last pass changed.` On a plain run this reads "rewrite index.html, read it and its assets", the same meaning as before.
- **D-16** (DECIDE & STATE): tree-list layout on `index.html`: `<!-- second-opinion-tree -->` + `<nav aria-label="Second-opinion iterations">` with a title line `Second-opinion iterations (tree id · preset · model · date)` and one `<li style="margin-left:<1.5 × (depth-1)>em">` per tree page: `<a href="<page>"><id></a> · <preset or plain> · <model> · <date>`, then `<!-- /second-opinion-tree -->`. No whitespace outside the markers, so rebuilding never accumulates bytes. Depth comes from the id length.
- **D-17** (DECIDE & STATE): parent header paragraph `<p><strong>Iteration <id></strong>, built on <a href="<parent>"><parent></a></p>`, placed first inside REQ-703's existing `<aside aria-label="Second-opinion direction">`; the aria-label is kept unchanged even when a tree page has no direction (session finalize on a base page with no flags), to avoid churn in REQ-703's block.
- **D-18** (DECIDE & STATE): the model meta is `<backend> <publishedModel>`, so an unverified tree page shows `codex gpt-5-unverified`; the date is the UTC publish date (`time.DateOnly`).
- **D-19** (DECIDE & STATE): `parentTreeID` derives the parent from `Base` by the tree pattern; a `Base` that does not match (only possible through an edited prepare record) falls back to the plain/`v1` rule. I removed a finalize-time validation of the record's `base` I had first added: not earned (see APPLY).
- **D-20** (DECIDE & STATE): in the tree list, `banner_state` is `inserted` when the rebuilt bytes differ from the original, else `already-present`; the change detail reads `second-opinion iteration list only`.

## Discovered Tasks

- `reportMetaAttribute` in `report_index.go` matches `\bname`, so a `data-name="..."` attribute is read as `name` by `readReportMetas` (and by the widened second-opinion meta strip, where the effect is only that such a tag is also stripped). impact-negligible → report only
- `lessons-do-work-cli.md:95` (REQ-703 bullet) links the flat path `do-work/archive/REQ-703-second-opinion-preset-and-editor-direction.md`; closing UR-160 moves it under `do-work/archive/UR-160/`. This is the one such link in `skills/` (grep `blob/main/do-work/archive/REQ-70[0-9]`); the integrator's sweep, as the brief says. impact-minor → report only

## Lessons read

- `_dev/primes/lessons-releases.md` (whole).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, families only: `exact-basename-authority` (REQ-485), `writer-anchored-marker-position` (REQ-528, REQ-703), `native-backend-identity` (REQ-698, REQ-701), `session-is-the-backend` (REQ-700), `fixture-cost-is-subprocess-spawning` (REQ-574).
- Prime traps: `prime-do-work-cli.md` § Traps, `prime-action-files.md` (whole; its `alternate-writer-contract-drift` trap led me to grep every shipped reader: only `SKILL.md`'s trigger description and `ai-report-reference.md:148` mention the form, and only the latter needed a change).

## Anti-bloat check

Added and named by the brief: `Base` field, `parentTreeID`, `secondOpinionTreePagePattern`, `secondOpinionChangeListRequirement`, `secondOpinionChangeListKind`, `secondOpinionStripBlocks` (shared marker-strip helper), `secondOpinionTreeList` (tree-list builder), `secondOpinionTreeMarks` (optional struct), `TreeID` and `BasePage`.

Added and not named by the brief:
- `secondOpinionChangeList(page)` (function): the D-08 check, kept out of `secondOpinionPublication` so the tail stays readable; one caller.
- `secondOpinionTreeStart`, `secondOpinionTreeEnd` (two constants): the marker pair is used by both the decoration strip and the tree-list builder, so one spelling.
- Test helpers `secondOpinionValidChangeList` (const), `secondOpinionFakeTreeEditor`, `secondOpinionStubJudgePass`: shared by the three new tests (the brief's "use `secondOpinionFakeCodex`'s extra argument" runs through `secondOpinionFakeTreeEditor`).
- Import `cmp` (standard library) for `cmp.Or(direction.Base, "index.html")`.

No new flag, environment variable, state file, config or package.

## Proposed CHANGELOG entry (integrator writes the version)

**Second Opinions Build a Linked Iteration Tree**

You can now build a second opinion on top of an earlier one, and walk through every iteration from the original report. Directed pages now carry a tree id in their names.

- `ai-report second-opinion` accepts a page path as its target: `index.html` or a tree page such as `index-v1c01-codex-gpt-5-bold.html`; the page names both the bundle and the page the editor rewrites.
- Every directed run (`--preset` or `--editor-direction`) and every run on a tree page gets a tree id: `index.html` is `v1`, its passes are `v1c01`, `v1c02`, a pass on `v1c01` is `v1c01c01`. The id sits in both names (`index-v1c01-codex-gpt-5-bold.html`, `v1c01-codex-gpt-5-bold-reply.md`), so a plain `ls` shows the tree in order. Plain runs keep their old names.
- Each tree page's header links its parent page, and its head carries `ai-report-second-opinion-id`, `-base`, `-model` and `-date` metas.
- The editor must open a tree page with a "Changed in this iteration" list linking `#rev-<n>` sections; a link without its target is an `AI-REPORT-SECOND-OPINION-CHANGE-LIST` judge finding that keeps the page and withholds the link.
- After judge `pass`, one launcher-owned iteration list on `index.html` shows every tree page (id, preset, model, date) and replaces the per-filename banner for tree pages.
- Result rows gain `tree_id` and `base_page`.

## Proposed lesson bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`

- [family: writer-anchored-marker-position] [REQ-704: a launcher block the launcher rebuilds in place must carry no whitespace outside its markers, and every copied launcher mark the editor can bring from a base page (header, tree block, any `ai-report-second-opinion-*` meta in either attribute order) must be stripped before the backlink check, because a first-wins meta reader would otherwise show the base page's values](../../../../do-work/archive/UR-160/REQ-704-second-opinion-linked-iteration-tree.md#lessons-learned)

## Test wall times

- `report_second_opinion_test.go`: 13.1 s (ok line), 13.3 s wall, quiet machine (load 1.65). Before this change the gate measured 27.2 s for this file under gate load; the D-12 edits add no runs.
- `report_second_opinion_tree_test.go`: 2.5 s (ok line), 2.7 s wall.
- Probe: 3.3 s. `-run SecondOpinion` package selection: 16.9 s. No budget failure, no reruns.
