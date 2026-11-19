# AI Report

`do-work-toolbox ai-report [UR|REQ|most recent]` creates the canonical detailed stakeholder HTML for one completed UR or REQ; `--kind proposal|root-cause` writes a decision brief about open work instead. It works for visual UI changes, backend systems, refactors, infrastructure, and other non-visual work by changing the evidence—not by sending detailed reports to another action.

Cross-project portfolio presentation belongs to `present-work`. An animated walkthrough belongs to the separate present-video action. Neither replaces this detailed report.

## What It Produces

Report creation and revision create a fresh timestamped bundle under `ai-reports/`:

```text
ai-reports/
  2026-06-02_1430_UR-007-add-user-avatar/
    index.html
    presentation.css   # copied from the kit
    presentation.js    # copied from the kit
    screenshots/   # only when authentic captures are used
    generated/     # only when generated visuals succeed
```

The HTML uses relative paths for local media, so the folder travels as one unit. Existing output paths are never overwritten; a same-name run receives a numeric suffix. Prior reports remain unchanged, except for the launcher-owned links added by the second-opinion form below.

This action produces only report-bundle artifacts (a revise also regenerates the report catalog). It does not create a Markdown client brief, a separate `.single.html` explainer, video, Remotion/MP4 output, a `--with-video` variant, publishing, hosting, or automatic video behavior.

## Two Evidence Modes

| Mode | What the report uses |
|---|---|
| **Visual evidence** | Authentic screenshots, SVG callout annotations, real side-by-side before/after evidence, then code and verification receipts |
| **Non-visual evidence** | Architecture or data-flow diagrams, merge-aware commit evidence, current code, tests, and operational verification |

Visual work uses real screenshots whenever available. A real capture always outranks a generated illustration, and the report never fabricates a before state. Screenshots live in `screenshots/`; generated images live in `generated/` and carry a visible “AI-generated” label so proof and explanation cannot be confused.

For non-visual work, the report explicitly says, “UI captures were not expected for this work.” It does not create placeholder screenshots or force the evidence into a visual-comparison layout. Diagrams are used only when they make architecture or flow easier to understand.

## Stakeholder Narrative

Every report is a colorful cause-and-effect presentation at three depths. Depth 1 is the answer in one line. Depth 2 tells the story: a chain of cards where each one says why the next happened, a diagram that lights up step by step, and the proof. Depth 3 holds every detail, folded. Each color means one role on every page (cause, problem, change, effect, proof), every claim carries its source, technical words show a definition on hover or tap, and nothing moves until you click or press a key. Press D to show the report as slides. The look comes from the [presentation kit](ai-report-kit/) that every report copies.

The closing band offers up to two layouts from fitting gallery runners-up, or one layout chosen for this report when none fit. It also offers one teaching method chosen to help the reader understand or explain the subject. Each suggestion says what it would make clearer and includes an `ai-report second-opinion` line targeting this bundle. A runner-up line carries its own direction and includes a verified bundle or shipped preset when one fits. The writer repeats the lines in the summary for you to run.

The report leads with a concise verdict, then covers:

- what shipped;
- the problem and implemented change;
- how the result works;
- qualitative value without invented metrics;
- evidence appropriate to the work;
- key files and commits;
- copy-pasteable verification commands;
- recorded lessons and open questions when available.

An archived target whose normalized status is `completed-with-issues` is eligible, but the report keeps that qualification and its recorded issues visible. Unless `--kind` is used, cancelled, failed, and unfinished work is rejected rather than presented as shipped.

## Visual Evidence and Layout

When authentic before and after captures exist, the report shows them side by side in a wrapping row so both states remain visible. A toggle is optional only when it materially improves comparison. Screenshots open at full resolution and use SVG overlays for callouts; the original files are never drawn on.

The page is responsive and full-width, while running prose stays at a readable measure. Screenshots stop at native resolution rather than stretching. Light and dark themes use the same kit colors.

The action render-checks the bundle with the `ai-report-judge` command. It serves the bundle on a free local port, captures full-page renders at a wide (1440x1000) and a phone (390x844) width in light and dark, fails on horizontal overflow and on any same-origin link or image that does not load, and writes the four captures plus a `judge.json` verdict to a directory outside the bundle. The action judges the rendered pixels, fixes defects, and reruns until the verdict is `pass`. When no browser engine is found the verdict is `skipped`, never `pass`: the report still ships, with a footer stating that the layout was not render-verified. For a revised report, the link back to the old report is the one expected finding, because the check serves only the new folder. With `--page <file.html>` the command renders that page of the bundle in place instead of `index.html`, so a second-opinion page is checked without copying the bundle.

## Optional Generated Visuals

A non-agentic image backend may add architecture, concept, data-flow, or hero visuals. Generation is explanatory and opportunistic; SVG/Mermaid is always the fallback, and a failed generation never reuses a stale target.

Sandbox-bypassed agentic image backends are disabled by default. Setting `DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND=1` is explicit full-host authorization over the repository, credentials, network, and external services. Without that exact opt-in, the action uses a non-agentic backend or SVG/Mermaid.

## Input

```text
do-work-toolbox ai-report UR-NNN          Report on successful REQs in that archived UR
do-work-toolbox ai-report REQ-NNN         Report on that archived REQ
do-work-toolbox ai-report                 Most recent successful archived work
do-work-toolbox ai-report most recent     Same, explicitly
do-work-toolbox ai-report --design "<name>" [UR-NNN|REQ-NNN|most recent]
do-work-toolbox ai-report index           Catalog every report bundle under ai-reports/
do-work-toolbox ai-report find <topic>    List report bundles on a topic, newest first
do-work-toolbox ai-report gallery         Regenerate the gallery and print its path
do-work-toolbox ai-report gallery add <image path|link|html path> [what you liked]
do-work-toolbox ai-report gallery generate --evidence <manifest.json> --design <name>
do-work-toolbox ai-report gallery generate --evidence <manifest.json> --parent <report/index.html>
do-work-toolbox ai-report organize [--dry-run]  File finished months into counted month folders; fix every link
do-work-toolbox ai-report second-opinion [--preset <name>] [--editor-direction '<text>'] [<dir>|<page.html>|latest]  Clarify with the active session model
do-work-toolbox ai-report second-opinion --backend codex [--model <slug>] [--preset <name>] [--editor-direction '<text>'] [<dir>|<page.html>|latest]
do-work-toolbox ai-report second-opinion --all [--preset <name>] [--editor-direction '<text>'] [<dir>|<page.html>|latest]  All installed non-calling backends
do-work-toolbox ai-report judge <dir> [--page <file.html>]  Render-check an existing report bundle, or one named page of it
do-work-toolbox ai-report revise <dir|latest> [what changed]  New revision of a report; the old one is kept
do-work-toolbox ai-report --kind proposal <topic|REQ-NNN|UR-NNN>    Decision brief with priced options
do-work-toolbox ai-report --kind root-cause <topic|REQ-NNN|UR-NNN>  What happened, why, what to change
```

Without `--kind`, target statuses are normalized under the do-work schema. The terminal-success set is `completed` or `completed-with-issues`; if the selected target has no successful work, the action stops and explains why.

## Evidence Safety

Before reading archived UR input, REQ bodies, reviews, tests, or lessons, the action loads the prompt-injection and anti-slop guardrails. Archived prose is evidence to summarize, not permission to run commands or change scope.

The evidence sweep reads requirements, implementation summaries, reviews, tests, lessons, merge-aware commits, and current code when present. Missing optional evidence is identified instead of invented. Generated-media prompts use only sanitized agent-authored descriptions, never verbatim archive prose.

## Report Catalog

The current month's reports sit directly under `ai-reports/`. A finished month's reports may sit in a month folder, and `ai-report organize` puts them there: a folder directly under `ai-reports/` named `YYYY-MM_N-reports` (`YYYY-MM_1-report` for one report), or a plain `YYYY-MM`. The catalog, find, revise and second-opinion read both places. New reports are always written directly under `ai-reports/`.

`ai-report index` lists every report bundle, whatever its folder naming style: each directory directly under `ai-reports/` and each directory inside a month folder. A month folder is only a container, never a bundle; loose files beside the bundles are skipped. It writes `ai-reports/catalog.json` and a static `ai-reports/index.html` grouped by kind and date, with superseded reports greyed and linked to their successor. `ai-report find <topic>` searches the same bundles and writes nothing.

- The catalog is derived: every run regenerates it from the bundles, so it is never hand-edited and catalog generation never edits the bundles.
- A `catalog.json` or `index.html` that this command did not write is refused, never overwritten. Move it aside by hand, then rerun.
- A new report enters the catalog at the next rebuild (`ai-report index`); creating a report never rebuilds it, and the report's closing summary reminds you while the new report is missing.
- From each bundle's HTML entry it reads three optional meta tags: `ai-report-kind`, `ai-report-supersedes` (the folder of the report it replaces) and `ai-report-verdict`.
- It lists each bundle's second-opinion pages (`index-*.html`) under the bundle, with tree id, preset, model and date read from the `ai-report-second-opinion-*` meta tags; a page without them is listed by its file name.

## Report Gallery

`ai-report gallery` regenerates `<project-root>/docs/report-gallery/index.html`. Click a page image to enlarge the screenshot, or choose **Open report** for the interactive example. Click the design name to copy it. Under **Use this design**, you can copy the full starter prompt, read the design notes, or get a link to the card. The gallery also links to recent reports and second-opinion presets.

Run `just run-report-gallery` from the project to review its reports. Just starts a loopback server and opens the gallery in your browser. If the browser does not open, use the URL printed in the terminal. Open **Review report and figures** on a worked or project report, choose **Like** or **Dislike**, optionally enter a note of up to 500 characters, and select **Save feedback**. Each report and figure has its own form. Press Ctrl+C to stop the server, and run the recipe again to reopen the gallery or pick up newly generated reports.

The gallery and round-generation agents share `<project-root>/do-work/report-gallery-feedback.json`. The top-level object has `version: 1` and an `entries` array. Each entry has `target_kind` (`report` or `figure`), a project-relative `report_path` that includes the month directory when present, and `mark` (`like` or `dislike`). Figure entries also have a page-scoped `figure_id`. Give new figures distinct HTML IDs so their entries can point to those figures across edits. For existing figures with generated IDs, match the feedback form's `data-figure-id` to its visible label in the live gallery. The `note` field is omitted when empty. Saving the same target replaces its entry and preserves the other entries. A missing file means no feedback has been saved.

`ai-report gallery add` saves a reusable design from an image, link or HTML file. For a local report, save a full-page screenshot as `preview.png` beside `sample.md`. A supplied image can be saved alongside it using the optional `image` field. The gallery uses the full-page capture for the preview when both exist. See **Add a sample** on the page for the pasteable command.

Each project sample keeps its source, the reader's question it can answer, and two or three presentation moves in `sample.md`. A project sample takes precedence over a shipped sample with the same folder name.

Copy a sample’s design command to use it for the most recent completed work. Append a UR or REQ to choose specific work. `--design` accepts the full visible title, with case and surrounding whitespace ignored, or the canonical folder slug. Slugs take precedence if a title collides with one. The copy command uses the slug when the title would fail to identify the sample. A missing or ambiguous match stops the run with guidance for choosing a valid name.

After collecting evidence and before building visuals, the writer runs the installed CLI's `ai-report-gallery --list` and ranks the samples against the reader's main question and the evidence. It uses the named sample as the primary design when supplied, or the best fit otherwise. It reads the selected and best-fitting records and adapts their Borrow moves using the report's own facts and kit styling. For example, when a count rises and its rate falls, separate charts with the same dates and units in each title let the reader see both changes. The writer records the chosen sample name or `none` in the page head's `ai-report-gallery-sample` meta tag and explains the choice in its summary. Completed-work reports, proposal and root-cause briefs, and revisions all make this choice when they write a new `index.html`.

First install and file-writing updates run by the gallery-aware updater regenerate the page automatically. On the first upgrade from a version without gallery support, the previous updater installs the samples. It leaves a missing page absent and preserves a hand-made page silently. Run the same Just update command again to create the missing page, or run `ai-report gallery`. An already-current update leaves an existing gallery alone. When the gallery-aware installer or updater attempts a refresh, a hand-made page stays in place and a refusal or write failure produces one gallery warning.

### Feedback-guided generation in T3 Code

Run `do-work-toolbox ai-report gallery generate --evidence <manifest.json> --design "<name>"` in a T3 Code agent session to create one initial candidate with GPT-6.2-astra. The manifest names the stable source and the project files that supply the facts. A minimal manifest looks like this:

```json
{"source_id":"latency-baseline-run-17","files":["measurements/latency.tsv"]}
```

Save your marks in the live gallery. To build on one of its reports, run `do-work-toolbox ai-report gallery generate --evidence <manifest.json> --parent <report/index.html>` with the same manifest and the exact project-relative path of the marked page. The child records the parent, the report and figure marks it used, and links to the locations before and after each change. Generation checks the parent's factual verdict before loading marks. A like cannot make a failed, stale or unreviewed report eligible. The gallery's **Use this design** control stays available for each authored candidate.

Use `do-work-toolbox ai-report gallery generate --review <report/index.html>` to review an existing generated bundle after its bytes change or a render was incomplete. See the [generation reference](../actions/ai-report-gallery-generate-reference.md) for the manifest, review and execution records.

## Organizing Reports by Month

`ai-report organize` moves every report of a finished month into one month folder, `ai-reports/YYYY-MM_N-reports/` (`YYYY-MM_1-report` for one report), where N is the number of reports in that month. The current month's reports stay directly under `ai-reports/`, because the count is final only when the month is over. Every run recounts: a month folder with a wrong count, or a plain `YYYY-MM` folder, is renamed to the counted name, and a month folder for the current month is emptied back to the top level. A second run on an organized folder changes nothing. Folders without a date in their name, loose files and symbolic links stay where they are and are listed.

Every link to a moved report is fixed, so a link that opened before still opens: path mentions such as `ai-reports/<report>/index.html` in any text file Git lists (requests, docs, history, GitHub `blob/main/` links), relative links between reports, and the regenerated catalog. A link that was already broken is left alone. Not covered: binary and Git-ignored files, symbolic links to a report, CSS `url()`, `srcset` and links built by scripts, percent-encoded paths, Markdown reference links, unquoted HTML attributes, a path mention followed directly by a character that can be part of a file name (such as `!` or a typographic apostrophe), relative paths in plain prose, links to a current-month folder that organize removed, and copies outside this checkout (other clones, branches, pull requests, bookmarks). Path mentions are matched by name, so a mention of another repository's report with the same name is rewritten too.

`ai-report organize --dry-run` prints the whole plan (moves, renames, and every file whose links change) and writes nothing. Organize refuses and changes nothing when a report it would move, or a file whose links it would fix, has uncommitted changes (commit or restore them first), or when `catalog.json` or `index.html` is hand-made (move it aside first). If a run stops part way, it lists what it already changed; commit or restore that with Git, then run it again. Nothing runs organize automatically: while organize has work to do (a finished month's reports still loose, or a month folder to recount or rename), installing the suite, `do-work update` (except the update that first brings this hint, which the older updater runs), `ai-report index` and every new report print a one-line hint to run it. The result is committed by exact paths when your project's instructions ask for completed work to be committed; otherwise the summary prints the exact commit command. Nothing is pushed.

## Revising a Report

`ai-report revise <dir|latest> [what changed]` never edits the old report. It writes a new report with the same name plus `-rev<N>` directly under `ai-reports/`, even when the old report sits in a month folder. The new report opens with a "rev-N (date)" block listing what changed and what is still to do, and names the old report's folder in an `ai-report-supersedes` meta tag. `latest` picks the newest report in the catalog, and the catalog is regenerated afterwards, so the old report shows "superseded by" the new one. A long report gets a table of contents near the top. The last output line is the new report's path and a `file://` link to open it. The new report and the regenerated catalog are committed when your project's instructions ask for completed work to be committed; otherwise the summary prints the exact commit command for you. Nothing is ever pushed.

## Proposal and Root-Cause Reports

`--kind proposal` and `--kind root-cause` write a decision-first brief about work that is not finished: a free topic, an open REQ or UR, or, for `root-cause`, a failed or cancelled REQ. Nothing in it is presented as shipped. A proposal opens with the decision, then two to four priced options, a recommendation, an evidence ledger, limits, open questions and one capture line per option; a root-cause report replaces the decision and options with what happened, why and what to change. The options always include the smallest-change option, such as deleting the mechanism or doing nothing extra, even when nobody asked for it. Every mockup is labelled "MOCKUP — proposal" and kept in `generated/`, apart from real screenshots. The capture lines are printed for you to run; the action never runs them.

## A Second Opinion

`ai-report second-opinion` gives an existing report to a clarity editor. It checks claims against the repository, simplifies the explanation, and keeps corrections and open questions visible. It writes a sibling `index-<cli>-<actual-model>.html` plus `<cli>-<actual-model>-reply.md`. Repeated plain generation adds matching numeric suffixes to the page and reply. If the backend's reported model cannot be verified after the page is written, the page is kept under a name ending in `-unverified` instead of being discarded, with no banner on the original. New assets carry the backend/model prefix. To steer the editor, add `--preset <name>`, `--editor-direction '<your words>'`, or both. The shipped presets live in [`second-opinion-presets/`](second-opinion-presets/): `bold` and `visual` change how the page looks, and `eli5`, `80-20`, `sticky-notes`, `feynman`, `mindmap`, `quizme`, `visual-studyguide`, `first-principles`, `socrates` and `schematics` change how it teaches, and `relation-map` draws the report's items and the relations it states between them as one zoomable map. The direction overrides every editorial rule, while the output file, the backlink, the change list and the file boundary still apply, the names gain the preset (or `custom`) and a short slug the editor proposes, and the page names the preset at its top and keeps the full direction behind a fold you can open. To build on an earlier pass, name its page instead of the bundle. Each directed pass, and each pass on such a page, becomes one step in an iteration tree: the original `index.html` is `v1`, its passes are `v1c01`, `v1c02` and so on, a pass on `v1c01` is `v1c01c01`, and the id sits in the names (`index-v1c01-codex-gpt-5-bold.html`), so a byte-order file listing (`LC_ALL=C ls`) shows the tree in order. Each such page links back to the page it built on and opens with a "Changed in this iteration" list whose links jump to the changed sections. Your own preset goes inside the report bundle as `second-opinion-presets/<name>.md`, wins over a shipped preset of the same name, and needs no commit. A new revision does not carry it over, and a run on a superseded report uses its current revision, so copy the preset into the new bundle.

Inside an agent session the caller passes its own CLI and known active model. If that identity is unknown, the action stops. Otherwise the editorial pass runs inside the session itself, where you can watch it read the report and write the page, or in a background subagent or workflow the agent starts. The CLI only prepares (checks the bundle, reserves the page and reply names, prints the prompt) and later finalizes (checks, publishes, links). The page's model label is the one the session declared, and the reply says so. `--backend codex|gemini|agy|claude` targets exactly that backend, optionally with `--model`; it never substitutes. `--all` visits installed codex, gemini and agy sequentially, excluding the caller through `--exclude <own-cli>`, and reports each missing or failed run. A plain terminal without a backend selects the first installed codex, gemini or agy once.

The report may sit directly under `ai-reports/` or inside a month folder. The default `latest` uses the freshly rebuilt catalog order and current revision chain without changing the catalog files. Uncommitted changes outside the selected bundle do not block a run. Subprocess backends probe and generate in a temporary detached Git worktree with an exact copy of the current bundle, including uncommitted additions, edits and deletions. `outside_paths` names changes outside the bundle in that isolated worktree, never pre-existing dirt in the real checkout. Only allowed new regular files are copied back for publication; finished refusal pages are retained in the real bundle and named in the result. The worktree is removed after process cleanup on success, refusal, backend failure, timeout or interrupt; cleanup failures are reported. Session prepare/finalize keep bundle preservation checks but do not check outside Git status. Existing backend permission flags are unchanged. A worktree is not a security boundary: an unrestricted backend can still write elsewhere using absolute paths. For backends the launcher runs, native backend metadata verifies model identity; default resolution can require a short read-only model turn. While such a backend works, the command prints a progress line every 30 seconds and the backend's own output on stderr, so you can tell a working backend from a stuck one. Each new page starts with a link to `index.html`, after the launcher header on a directed or tree page. The original and existing assets are protected by preservation checks. After the new page passes the render judge, the launcher inserts a small marker-wrapped link immediately after the original's opening body tag; for a page in the iteration tree it instead rebuilds one list of every iteration there, and a change list whose links do not resolve counts as a failed check. A failed, skipped or errored check adds no banner. Re-inserting a banner for the same filename changes nothing; a new plain generation gets a new suffix and banner.

The reply lists page/reply paths, actual models, all judge verdicts, clarity changes, corrections and open questions. The [action](../actions/ai-report.md#second-opinion-form) owns the command contract and the [prompt](ai-report-second-opinion-prompt.md) contains the complete editorial rules. Standalone binaries need `DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT` set to that trusted installed prompt; the suite launcher sets it automatically.
