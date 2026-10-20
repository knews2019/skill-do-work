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
write_set: ["skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md", "skills/do-work-toolbox/docs/second-opinion-presets/bold.md", "skills/do-work-toolbox/docs/second-opinion-presets/visual.md", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work/tools/do-work-cli.sh", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/SKILL.md"]
claimed_at: 2026-10-20T20:02:31Z
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
