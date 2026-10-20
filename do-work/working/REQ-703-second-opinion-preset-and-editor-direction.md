---
id: REQ-703
title: 'Second opinions accept a named preset and a free-text editor direction'
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
related: [REQ-704]
batch: UR-160
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md", "skills/do-work-toolbox/docs/second-opinion-presets/bold.md", "skills/do-work-toolbox/docs/second-opinion-presets/visual.md", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/SKILL.md", "CHANGELOG.md", "skills/do-work/CHANGELOG.md", "VERSION", "skills/do-work/VERSION", "skills/do-work/actions/version.md", "skills/do-work/tools/do-work-cli/lessons-do-work-cli.md", "do-work/lessons-index.md"]
claimed_at: 2026-10-20T20:02:31Z
route: C
estimate:
  p50_active_minutes: 55
  confidence: low
  calculated_at: 2026-10-20T20:09:15Z
  basis:
    - Route C
    - 18-file write set
    - 3 subsystems involved
    - 9 acceptance criteria
planning_at: 2026-10-20T20:10:46Z
dispatch_at: 2026-10-20T20:17:40Z
builder_handback_at: 2026-10-20T20:27:29Z
---
# Second Opinions Accept A Named Preset And A Free-Text Editor Direction
## What
`ai-report second-opinion` gains `--preset <name>` and `--editor-direction '<text>'`, usable alone or together in the Session (`--prepare`), Targeted and Batch forms. The direction text is placed into the clarity-editor prompt, overrides every prompt rule except the file BOUNDARY, and is visible in the published filename, on the published page, and in the typed result row.
## Why
The maintainer wants to steer second opinions as a director ("add motion, add more sections, go further") instead of always getting the one fixed clarity-editor pass. Today the prompt is fixed and `secondOpinionParse` refuses any extra argument.
## Prior Implementation
- REQ-698 (ai-report second-opinion, 0.306.0 era, commit 31e7fcdf): the launcher `ai-report-second-opinion` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` resolves the bundle, reserves a staging page, fills the shipped prompt `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` (read through env var `DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT`, set by `skills/do-work/tools/do-work-cli.sh`), runs the backend, checks preservation and outside-bundle Git status, runs the render judge, and inserts a per-filename banner on `index.html` only on judge `pass`. Pages are `index-<cli>-<model>[-N].html` with `<cli>-<model>[-N]-reply.md`.
- REQ-701 (0.306.3, commit e8b063c1): a failed model check keeps the finished page as `-unverified`; it is never deleted.
- REQ-700 (0.306.4, commit c70b9b72): Session form is `--prepare` (returns `staging_page`, `staging_reply`, `asset_prefix`, `prompt`) then `--finalize <staging_page>`; the prepare record is the `.json` sidecar (`secondOpinionRecord`).
- REQ-702 (0.306.5, commit 26a20f10): heartbeats and mirrored backend output on stderr.
- The prompt is filled by two duplicated `strings.NewReplacer` calls, in `secondOpinionPrepare` and `secondOpinionRun`.
## Detailed Requirements
- D1: `--preset <name>` and `--editor-direction '<text>'` are accepted by Session prepare, Targeted and Batch runs; both together means the preset text first, then the free text. Neither is accepted with `--finalize` (finalize reads them from the prepare record).
- D2: shipped presets live in `skills/do-work-toolbox/docs/second-opinion-presets/<name>.md`; ship `bold.md` ("add motion, add more sections, go further") and `visual.md` ("diagrams and charts first, prose second; turn tables and lists into pictures where possible"). A project file `ai-reports/second-opinion-presets/<name>.md` with the same name wins. An unknown preset name is a usage error that lists the presets found. The result row records `preset` and `preset_file` (the file actually used).
- D3: the prompt gains a direction paragraph (placeholder such as `{{EDITOR_DIRECTION}}`) that says the direction overrides every rule except BOUNDARY. With no direction it says so plainly ("None"), and a plain run behaves as today.
- Free-text validation: refuse an empty value and a value containing `{{` (it would trip the unresolved-placeholder check with a misleading message).
- Visibility on the page: the launcher writes `<meta name="ai-report-second-opinion-preset">` and `<meta name="ai-report-second-opinion-direction">` into the published page head and a launcher-owned header block at the top of the page body showing the preset name and the full direction text. Launcher-owned blocks use markers, like the existing banner, and are idempotent.
- Visibility in the filename (D6/D7 part owned here): a directed page is named `index-<cli>-<model>[-unverified]-<preset>[-<editor-slug>][-N].html`; free text only uses `custom` in the preset slot. The editor slug is a 2-4 word kebab slug the editor proposes on a `Variation slug: <words>` line in its FINAL REPLY; the launcher lowercases it, keeps `[a-z0-9-]`, and limits its length. The reply file follows the same stem. REQ-704 later adds the tree id after `index-`.
- Visibility in the result: `AIReportSecondOpinionRun` in `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` gains `preset`, `preset_file` and `editor_direction` (omitempty).
- Rule 8 of the prompt gains "honor `prefers-reduced-motion` for any animation".
- Merge the two duplicated prompt fills into one helper that also fills the new placeholder.
- Docs: `skills/do-work-toolbox/actions/ai-report.md` § Second-opinion form gets the flag grammar and a trigger example ("second opinion with bold", "second opinion with '<text>'" means pass the words as one argv token); `skills/do-work-toolbox/docs/ai-report-guide.md` gets the command lines and one sentence; `skills/do-work-toolbox/SKILL.md` description gets the trigger words `second opinion preset`.
## Constraints
- A plain run (no preset, no direction) keeps today's filenames, prompt meaning and behavior.
- REQ-701's rule stays: a finished page is never deleted by a later check.
- The editor still may write only its staging page and prefixed assets.
- Shipped change: this is a release per `_dev/primes/prime-releases.md`, with a descriptive changelog title (for example "Second Opinions Accept Presets and Your Own Direction").
- REQ-704 depends on this REQ and rewrites the same naming function; keep the filename builder in one function so REQ-704 can add the tree id in one place.
## Dependencies
None. REQ-704 (linked iteration tree) depends on this REQ.
## Builder Guidance
Decisions D1-D3 and D7 are user-confirmed and firm. The exact prompt wording, marker names and slug length limit are builder latitude. Keep it simple: no new state file; everything travels in the prepare record, the prompt, meta tags and the result row.
## Red-Green Proof
**RED prompt/case:** `do-work-cli.sh --repo-root . --format json ai-report-second-opinion --prepare --backend claude --model claude-opus-5-5 --preset bold ai-reports/2026-10-15_1731_REQ-698-ai-report-second-opinion` is refused with "unexpected argument \"--preset\"".
**Why RED now:** `secondOpinionParse` knows only `--backend`, `--model`, `--exclude`, `--all`, `--prepare` and `--finalize`, and the prompt has no direction placeholder.
**GREEN when:** (1) the prepare row's `prompt` contains the bold preset text and the row carries `preset: bold` and `preset_file`; (2) a fake-backend Targeted run with `--preset bold --editor-direction 'add a timeline'` receives both texts in its prompt and publishes `index-<cli>-<model>-bold-<slug>.html` with the preset/direction meta tags and header block; (3) a project `ai-reports/second-opinion-presets/bold.md` shadows the shipped one and `preset_file` names it; (4) a run with no flags still publishes `index-<cli>-<model>.html`.
**Validation:** User confirmed (proof split across the two chained REQs at the user's choice).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 22697 tokens, `slugged: partial` so no targeted form; matches native backend identity and the session-is-the-backend families. Read its `native-backend-identity` and `session-is-the-backend` bullets if budget allows.
- `_dev/primes/lessons-action-files.md` — 8398 tokens, `slugged: partial`; matches the action-file and downstream-reader changes.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
## Full Context
See `do-work/user-requests/UR-160/input.md` for complete verbatim input.

*Source: can I get a seccond oppinion with: "add motion, add more sections, go further"? (plus question rounds 1-4 and the capture round in UR-160)*

---

## Triage

**Route: C** - Complex

**Reasoning:** The user-confirmed decisions (D1-D3, D7) fix the outcome, but the change crosses the argument contract, the prepare record, the shared publication tail, the published filename, launcher-owned page markup, the typed result row, the shipped prompt, two new shipped preset files, the report catalog walk and three prose files. The pre-dispatch agent decides the mechanism before the builder starts, so the brief carries one design and REQ-704 (linked iteration tree) inherits one naming function and one header block.

**Planning:** Required

*No `## Open Questions` section exists; Step 3.5 has nothing to resolve. The builder-latitude points (prompt wording, marker names, slug length) are decided in `## Plan` → Decisions as D-01 onwards.*

<!-- D-XX counter: last used D-09. Next decision: D-10. -->

## Plan

**Decided mechanism (pre-dispatch, do not reopen):** two new valued flags on the existing `ai-report-second-opinion` subcommand. The handler resolves them once into one small `secondOpinionDirection` value, which travels in the prepare record (Session) or directly (Targeted and Batch) to one prompt-fill helper and to the shared publication tail. The tail names the page with one name builder and adds launcher-owned meta tags and a header block on directed runs only. A plain run fills `{{EDITOR_DIRECTION}}` with `None.` and is otherwise byte-for-byte today's path.

### Decisions
- **D-01** (DECIDE & STATE): the shipped preset folder is derived from the trusted prompt path: `filepath.Join(filepath.Dir(<DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT>), "second-opinion-presets")`. No new environment variable, so `skills/do-work/tools/do-work-cli.sh` and `_dev/tests/do-work-cli-launcher-behavior.sh` stay untouched and the capture's `write_set` entry for the launcher is dropped. Reason: the presets sit beside the prompt in the same shipped `docs/` folder, and a standalone binary already needs that one variable.
- **D-02** (DECIDE & STATE): a preset name must match `^[a-z0-9][a-z0-9-]*$` (filename-safe, no path traversal). Lookup order: `<repo>/ai-reports/second-opinion-presets/<name>.md`, then the shipped folder. The chosen file must be a regular file (`os.Lstat`; a symlink is refused, like bundle entries). Preset text is the whole file, trimmed. An empty file or one containing `{{` is a usage error naming that file. An unknown name is a usage error `unknown preset "<name>"; presets found: <sorted, de-duplicated stems of *.md in both folders>`. `preset_file` is repository-relative when the file is inside the repository, otherwise absolute.
- **D-03** (DECIDE & STATE): `reportIndexWalk` (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go:81`) skips the `second-opinion-presets` folder through one named constant. Without it, `ai-report index`, `find`, `judge` and `latest` treat the project preset folder as an undated report whose entry is its first `.md` file (`pickReportEntry`, `report_index.go:167`). D2 fixed the folder path, so this follows from it.
- **D-04** (ESCALATE): a project preset file sits outside the bundle, so the unchanged outside-bundle Git check refuses a run while that file is uncommitted. Kept on purpose; the guide says to commit a new preset before using it. Value: the one check that catches backend writes outside the bundle stays whole. Risk: the first run after writing a preset is refused with `outside-bundle changes: ai-reports/second-opinion-presets/<name>.md`; reversible by a later exemption for read-only preset paths.
- **D-05** (DECIDE & STATE): the prompt gains a `DIRECTION.` paragraph after GOAL and before RULES. Its fixed text says the direction below comes from the person directing this pass and overrides ROLE, GOAL, every RULE and PROCESS, everything except BOUNDARY; then `Direction: {{EDITOR_DIRECTION}}`. Fill: `None.` on a plain run; preset text, a blank line, then the free text when both are given. FINAL REPLY gains one sentence: when a direction is given, end with one line `Variation slug: <2 to 4 plain words naming this variation>`. Rule 8 gains "honor `prefers-reduced-motion` for any animation". Exact wording is builder latitude within these points.
- **D-06** (DECIDE & STATE): one helper `secondOpinionPrompt(template, repository, bundle, outputFile, prefix string, direction secondOpinionDirection) (string, error)` replaces the two `strings.NewReplacer` fills (`report_second_opinion.go:482` in `secondOpinionPrepare`, `:648` in `secondOpinionRun`). It fills all six placeholders, keeps the leftover-`{{` refusal, and refuses a directed run when the template has no `{{EDITOR_DIRECTION}}` (a stale standalone template would otherwise drop the direction silently). A plain run with an old template still works.
- **D-07** (DECIDE & STATE): `secondOpinionDirection` holds `Preset`, `PresetFile`, `PresetText`, `EditorDirection` (JSON-tagged, omitempty) with two methods: the combined text (preset text first, then free text; empty when neither) and the filename slot (preset name, or `custom` for free text only; empty when plain). `secondOpinionRecord` gains `Direction secondOpinionDirection` (omitempty), so finalize restores it from the record; a record written before this change reads as a plain run. No new state file.
- **D-08** (DECIDE & STATE): `secondOpinionNames(bundle, backend, publishedModel, variation string)` (`report_second_opinion.go:1202`) stays the only name builder: stem `<backend>-<publishedModel>[-<variation>]`, page `index-<stem>[-N].html`, reply `<stem>[-N]-reply.md`. `variation` is empty on a plain run, else `<slot>[-<slug>]`. The slug comes from the first `Variation slug:` line of the final reply (case-insensitive, optional leading `-` or `*`): lowercased, every run of characters outside `[a-z0-9]` becomes one `-`, trimmed, at most 4 words and 40 characters; empty means no slug segment. It is read only on directed runs. The unverified form falls out as `index-<cli>-<model>-unverified-<slot>[-<slug>][-N].html`. REQ-704 (linked iteration tree) inserts its tree id in this one function.
- **D-09** (DECIDE & STATE): page decoration lives in `secondOpinionPublication` (`report_second_opinion.go:696`). Order: strip any launcher header block and the two launcher meta tags from the editor's page (an editor may copy them from a base page in REQ-704; a no-op on plain pages); the unchanged backlink check; then, on directed runs only, insert `<meta name="ai-report-second-opinion-preset" content="<slot>">` and `<meta name="ai-report-second-opinion-direction" content="<combined text>">` before the first real `</head>` (masked with `secondOpinionRawHTML` like the banner; with no `</head>` the metas go inside the header block), and a header block right after the opening body tag between `<!-- second-opinion-header -->` and `<!-- /second-opinion-header -->`, showing the slot and the full combined direction text, HTML-escaped, with `overflow-wrap:anywhere` so a long token cannot fail the judge's overflow check. Decorating a decorated page gives the same bytes. The judge runs on the decorated page, which is what is published. REQ-704 adds its tree id and parent link to this same block.

### Files to modify
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify): usage string, options and parser, preset resolution in the handler, `secondOpinionDirection`, `secondOpinionPrompt`, record field, finalize restore, names, decoration in the publication tail.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify): fixture template gains `Direction: {{EDITOR_DIRECTION}}`; two new tests.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modify) and `report_index_test.go` (modify): the D-03 skip and one test.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify): three omitempty fields on `AIReportSecondOpinionRun`.
- `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` (modify), `skills/do-work-toolbox/docs/second-opinion-presets/bold.md` and `visual.md` (new).
- `skills/do-work-toolbox/actions/ai-report.md`, `skills/do-work-toolbox/docs/ai-report-guide.md`, `skills/do-work-toolbox/SKILL.md` (modify).
- Release files (integrator only): both CHANGELOGs, both VERSION files, `skills/do-work/actions/version.md`.

### Order of changes
1. **RED: preset resolution and argument rules.** Add `TestSecondOpinionPresetResolutionAndArgumentRules`; at base it fails on `unexpected argument "--preset"` (the REQ's RED).
2. **GREEN: parser, resolution, record, prompt helper, result row.** `--preset` and `--editor-direction` (both `--x value` and `--x=value` through `optionValue`); refuse an empty or whitespace-only direction and one containing `{{`; `--finalize` refuses both flags. The handler resolves the preset after reading the template (D-01, D-02) and passes the direction to `secondOpinionPrepare` and `secondOpinionRun`; `secondOpinionPrompt` replaces both fills (D-06); the record carries the direction (D-07); each row gets `preset`, `preset_file`, `editor_direction`. Test 1 green.
3. **RED then GREEN: names and page marks.** Add `TestSecondOpinionDirectedRunNamesAndMarksPage`; then the slug parse, the `variation` argument to `secondOpinionNames`, and the D-09 decoration.
4. **RED then GREEN: catalog skip.** Add `TestAIReportIndexSkipsSecondOpinionPresetsFolder`; then the D-03 skip.
5. **Shipped text and prose.** Prompt (D-05), the two presets, `ai-report.md`, the guide, `SKILL.md`; propose the changelog entry.

### Testing approach
- Focused, from `skills/do-work/tools/do-work-cli`: `go test -count=1 -v -run '^(TestSecondOpinionPresetResolutionAndArgumentRules|TestSecondOpinionDirectedRunNamesAndMarksPage|TestAIReportIndexSkipsSecondOpinionPresetsFolder)$' ./internal/toolboxcommands`.
- Must stay green: every existing `TestSecondOpinion*` test, every `TestAIReportIndex*` test, `TestHandlersRegisterCanonicalToolboxCommands` (flags, not a new command), `resultmodel.TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence`.
- Package checks: `go vet ./internal/toolboxcommands ./internal/resultmodel`, `gofmt -l ./internal`.
- Tests stay hermetic: the fixture puts its own `second-opinion-presets/` beside its temp template (D-01), never reading the shipped files. The integrator's probe checks the shipped prompt and presets with plain greps.

### Risks
- Header block before the backlink: the backlink check runs on the stripped editor page, before insertion; reversing the order fails every directed run.
- The strip must use the masked surface (real comments and tags only); a marker inside a `<script>` string must survive untouched.
- A project preset test must commit the preset first, or the outside-bundle check refuses (D-04).
- The fake codex reads its prompt from argv; the directed test must make it print the `Variation slug:` line only when the prompt carries a direction, or the plain-run case in the same test stops proving the plain name.

### Consumer field contract (per `runs[]` row)
- Identity and naming: `backend`, `requested_model`, `actual_model`, `page`, `reply` (directed names per D-08).
- Direction provenance: `preset` (the `--preset` value, empty when none), `preset_file` (D-02), `editor_direction` (the free text only); all omitempty, absent on plain runs.
- Session: prepare rows carry the three fields and the filled `prompt`; finalize rows carry them restored from the record.

*Generated by the pre-dispatch agent (planning done directly; the user approved the design as D1-D10 in UR-160)*

### Plan validation (pre-dispatch)
- Requirement coverage: flags in all three forms → tasks 1-2; shipped and project presets, unknown-name error, `preset_file` → tasks 1-2 and 5; direction precedence → task 5 (D-05); free-text validation → task 2; meta tags and header block → task 3; directed filename and editor slug → task 3; result row → task 2; reduced-motion rule → task 5; one prompt helper → task 2; docs → task 5. No uncovered requirement.
- Orphan tasks: task 4 is not named in the REQ; it follows from D2's folder path (D-03).
- Scope sanity: **Plan has 5 tasks — quality degrades past 3.** Not split: tasks 1-4 are one serialized slice of one Go file and its test, and task 5 is prose that must ship with it; REQ-704 is already the split.
- Consumer field contract: present above; the action prose reads exactly those fields.

## Exploration

Read directly by the pre-dispatch agent: `report_second_opinion.go` (all 1383 lines), the fixtures and session tests in `report_second_opinion_test.go`, `report_index.go` (`reportIndexWalk`, `pickReportEntry`), `commands.go` (`optionValue`, `usageResult`), `resultmodel/result_model.go` (`AIReportSecondOpinionRun`), the shipped prompt, `ai-report.md` § Second-opinion form, `ai-report-guide.md` § A Second Opinion and its command block, `SKILL.md`, `do-work-cli.sh:10`, the CLI prime, and the archived REQ-698 and REQ-700 records.

**Argument path.** `secondOpinionParse` (`report_second_opinion.go:195`) matches flags as literal strings; valued flags go through `optionValue` (`commands.go:87`), which refuses `--x=` with nothing after it but returns an empty string for `--x ''`, so the direction needs its own empty and whitespace check. A positional or unknown flag gives `unexpected argument %q` (`:215`), the REQ's RED. The `--finalize` refusal at `:239` gains the two new flags. The usage string is `secondOpinionUsage` (`:29`).

**Template and preset roots.** The handler reads the template from `DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT` (`:115`) only on the non-finalize path, which is also the only path that resolves a preset. `do-work-cli.sh:10` sets it to `<suite>/do-work-toolbox/docs/ai-report-second-opinion-prompt.md`, so the shipped preset folder is that file's sibling directory (D-01). `_dev/tests/do-work-cli-launcher-behavior.sh:252` builds a fake suite with only the prompt file; D-01 leaves it valid.

**The two prompt fills.** `secondOpinionPrepare` (`:482`, absolute stage path as `{{OUTPUT_FILE}}`, because the session's cwd is the repository) and `secondOpinionRun` (`:648`, bare stage name, because the backend's cwd is the bundle). Both are followed by the same `{{` refusal. The helper takes the output value as an argument, so both callers keep their own form.

**Prepare record.** `secondOpinionRecord` (`:310`) is JSON in `<stem>.json` beside the stage; finalize unmarshals it at `:516` and fills the row at `:521`. Adding an omitempty struct field keeps old records readable.

**Publication tail.** `secondOpinionPublication` (`:696`) is shared by the exec run and finalize. It runs the backlink check (`:704`, skipped for the unverified form), names the files (`:708`), publishes page then reply, rewrites the stage name inside the reply (`:727`), removes stage leaves, then runs the judge on `page` (`:745`) and the pass-gated banner. Decoration goes between the backlink check and `secondOpinionNames`, so the judge sees the published bytes. The reply text is available here for the slug.

**Markers.** `secondOpinionBanner` (`:1288`) shows the house pattern: find the body tag on the masked surface (`secondOpinionVisibleHTML`, `secondOpinionBodyTag`), look for a real comment marker after it, insert a marker-wrapped block. The header block follows the same pattern. `secondOpinionBacklink` (`:1242`) allows leading `p|div|aside|nav|header` opening tags before `Original report:` but no text, so a header block placed first must be stripped before that check (D-09). The banner's page-name regex `^index-[A-Za-z0-9][A-Za-z0-9._-]*\.html$` accepts directed names.

**Names.** `secondOpinionNames` (`:1202`) builds `index-<backend>-<model>[-N].html` and `<backend>-<model>[-N]-reply.md`, checking both are free with `os.Lstat`. The unverified form passes `<model>-unverified` as the model (`:701`), so a variation appended after the model gives the REQ's `-unverified-<preset>` order with no special case.

**Catalog walk.** `reportIndexWalk` (`report_index.go:81`) treats every non-dot directory under `ai-reports/` as a bundle; `pickReportEntry` (`:167`) falls back to the first `.md`. It is the only `ai-reports/` walker in the module (D-03). `secondOpinionBundle` already refuses a non-bundle path (no regular `index.html`).

**Tests and fixtures.** `secondOpinionFixture` (`report_second_opinion_test.go:21`) writes a temp template with five placeholders and sets the variable; it must gain `Direction: {{EDITOR_DIRECTION}}` so directed tests can see the fill. `secondOpinionFakeCodex` (`:52`) parses `Output: ` from its last argv and prints a fixed reply; the directed test needs a variant that echoes a `Variation slug:` line. Commit any project preset in the fixture repository before invoking (outside-bundle check, D-04). `TestHandlersRegisterCanonicalToolboxCommands` (`commands_test.go:5`) pins the handler count. `resultmodel.TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence` round-trips without a golden string, so omitempty fields are safe.

**Prose sites.** `ai-report.md:65-67` (Session, Targeted, Batch bullets: add the two flags), `:73` (published names: add the directed form and the header and meta marks), `:21` (trigger bullet: "second opinion with bold", "second opinion with '<text>'" → pass the words as one argv token). `ai-report-guide.md:72-74` (command lines) and § A Second Opinion (one sentence on presets, the project folder, and committing a new preset first; adjust "Each new page starts with a link to `index.html`" for the header block). `SKILL.md:3` description and the `ai-report` routing row at `:24` gain `second opinion preset`. `ai-report-reference.md:148` stays true. History (changelogs, lessons) is not edited.

**Required-lessons consult (re-run at claim).** `required_lessons` stays `_dev/primes/lessons-releases.md` (666 tokens, read: both families concern history links and manifest edits; this REQ adds no archive link). The index still lists `lessons-do-work-cli.md` at 22697 tokens and `lessons-action-files.md` at 8398 tokens, both `slugged: partial`, so the capture's `## Required Lessons — Dropped for Budget` section stands. Read directly as relevant bullets: `native-backend-identity` (REQ-698, REQ-701: identity refusal stops the name, never the bytes; the unverified name gains the variation, nothing else changes), `session-is-the-backend` (REQ-700: everything crosses the prepare/finalize boundary in the record), `writer-anchored-marker-position` (read a writer's marker at its written position; applies to the header block markers and to the editor's `Variation slug:` line, which is editor text and is only cleaned, never trusted as a path).

*Generated by the pre-dispatch agent's direct reading*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify) — flags, preset resolution, direction value, one prompt helper, record field, directed names, page header and meta marks
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify) — fixture placeholder and two new tests
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modify) — catalog walk skips the project preset folder
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (modify) — one test for that skip
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — three omitempty fields on the run row
- `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` (modify) — direction paragraph, slug line, reduced-motion rule
- `skills/do-work-toolbox/docs/second-opinion-presets/bold.md` (new) — shipped preset
- `skills/do-work-toolbox/docs/second-opinion-presets/visual.md` (new) — shipped preset
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — flag grammar, trigger example, directed names
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — command lines and one plain-words sentence
- `skills/do-work-toolbox/SKILL.md` (modify) — trigger words in the description and the ai-report routing row
- `CHANGELOG.md` (modify, integrator's release) — the entry
- `skills/do-work/CHANGELOG.md` (modify, integrator's release) — the byte-identical mirror
- `VERSION` (modify, integrator's release) — the version bump
- `skills/do-work/VERSION` (modify, integrator's release) — the version mirror
- `skills/do-work/actions/version.md` (modify, integrator's release) — the Current version line
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modify, integrator, only if the hand-back proposes a lesson) — one lesson bullet
- `do-work/lessons-index.md` (modify, integrator, only with the lesson) — token refresh for that satellite

**Files I will NOT touch:** `skills/do-work/tools/do-work-cli.sh` and `_dev/tests/do-work-cli-launcher-behavior.sh` (D-01: no new environment variable), `skills/do-work-toolbox/actions/ai-report-reference.md` (its banner sentence stays true), `commands.go` and `commands_test.go` (no new command), the capture supervisor and the codex, gemini, agy and claude adapters, the judge, and the three existing second-opinion pages under `ai-reports/` (UR-160 D10).

**Acceptance criteria (restated from REQ):**
- [ ] `--preset <name>` and `--editor-direction '<text>'` work alone or together in Session prepare, Targeted and Batch runs; together means preset text first, then the free text; `--finalize` refuses both and reads them from the prepare record.
- [ ] Shipped presets `bold` and `visual` live in the toolbox `docs/second-opinion-presets/` folder; a project file `ai-reports/second-opinion-presets/<name>.md` of the same name wins; an unknown name is a usage error listing the presets found; the row records `preset` and `preset_file` (the file actually used).
- [ ] The prompt has a direction paragraph that overrides every rule except BOUNDARY, says "None" on a plain run, and Rule 8 asks to honor `prefers-reduced-motion`.
- [ ] An empty direction and one containing `{{` are refused with a clear message.
- [ ] A directed page carries the preset and direction meta tags in its head and a marker-wrapped launcher header block at the top of the body showing the preset name and the full direction text; the marks are idempotent.
- [ ] A directed page is named `index-<cli>-<model>[-unverified]-<preset>[-<editor-slug>][-N].html` (free text alone uses `custom`), the slug cleaned from the editor's `Variation slug:` line, and the reply follows the same stem.
- [ ] `AIReportSecondOpinionRun` gains `preset`, `preset_file` and `editor_direction` (omitempty).
- [ ] The two prompt fills are one helper that also fills the direction placeholder.
- [ ] A plain run keeps today's names, prompt meaning and behavior; a finished page is never deleted by a later check; the editor still writes only its staging page and prefixed assets; `ai-report.md`, `ai-report-guide.md` and `SKILL.md` document the flags and triggers; the change ships as one release.

## Pre-Flight

**Git:** ✓ clean at a8c0cf11 (main after the REQ-703 claim), the builder's base commit; only do-work/ run artifacts were dirty when the coordinator recorded pre-flight
**Tests baseline:** ✓ `bash do-work/runs/work-2026-10-20-200259/REQ-703-preflight-probe.sh` exit 0, launched (`do-work/working/baseline.json`). Repository gate `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at a8c0cf11: run 1 exit 0 in 130 s at 1-minute load 4.4; green-gate record satisfied at revision a8c0cf11
**Dependencies:** ✓ Go toolchain and ShellCheck present (the gate ran)

*Checked by the coordinator before dispatch (2026-10-20T20:14:34Z to 2026-10-20T20:17:25Z)*
