## Review

**Overall: 94%** | 2026-11-09T12:22:01Z

| Dimension | Score |
|-----------|-------|
| Requirements | 96% |
| Code Quality | 93% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict: Approve.** The merge range `7dc47c2d..0c06eb10` delivers every Detailed Requirement and the Red-Green Proof. Ten learning-method presets ship beside `bold` and `visual`. The closing band keeps its layout suggestion and line and adds one teaching-method suggestion with a `--preset` or `--editor-direction` line. Every shipped description of the band now names both suggestions. The new Go test reads the real shipped folder. No finding is critical.

### Requirements Checklist

- [x] Ten preset files under `skills/do-work-toolbox/docs/second-opinion-presets/` with the REQ's names (`eli5`, `80-20`, `sticky-notes`, `feynman`, `mindmap`, `quizme`, `visual-studyguide`, `first-principles`, `socrates`, `schematics`): delivered. Each one is a single paragraph addressed to the editor in the voice of `bold.md`, has no `{{`, and ends with the fact rule.
- [x] Each preset changes how the page teaches, not only how it looks: delivered. `socrates` asks questions before revealing answers, `quizme` ends sections with hidden-answer checks, `schematics` draws the mechanism, and `feynman` adds "where this gets hard" gap notes plus a summary the reader can repeat.
- [x] The closing band keeps the layout suggestion and its line, and adds the method, why it fits, and a ready line: delivered (`skeleton.html:177-191`). The `id="another-layout"` stays, so existing `#another-layout` links and the `#another-layout pre` style still apply.
- [x] The method rule is a condition with an open menu, which is the user's "the model should not be limited by them": delivered (`ai-report-reference.md:91` says "a starting menu, not a closed list". The skeleton comment says the same).
- [x] A method outside the menu uses `--editor-direction` and is written without a single quote: delivered (`ai-report-reference.md:91`).
- [x] "Beside, not replacing" (the user's capture choice): delivered. The layout paragraph and its line are byte-identical.
- [x] Step 8 summary prints the method pick and its line: delivered (`ai-report.md:194`).
- [x] The guide no longer names only `bold` and `visual` and says where the presets live: delivered (`ai-report-guide.md:122`). See F1 for the new hand-maintained list.
- [x] Trigger sentence covers any shipped preset name: delivered (`ai-report.md:22`, "any other shipped or bundle preset name works the same way").
- [x] Go test over the real shipped folder pins the ten names and loads every `.md` file there: delivered (`report_second_opinion_test.go:1004-1039`). The hand-back records RED at 3599f71f (ten "does not ship" errors) and GREEN at 31acf248.
- [x] Fixture tests stay isolated from the real folder, and their `presets found: bold, visual` assertion still holds: delivered. The test suite passes.
- [x] Narrative item 9, the `--kind` sentence, the checklist line and the CSS comments describe both suggestions: delivered.

### Restatement Sweep

The diff redefines two elements: (1) the closing band, from one layout suggestion to a layout suggestion plus a teaching-method suggestion, and (2) the set of shipped presets, from `bold` and `visual` to twelve presets. I searched `skills/`, `suite/` and `tools/` for "another layout", `another-layout`, "layout suggestion", "suggested layout", "one layout", "closing band", `--preset`, "preset bold", "preset visual", `bold.md`, `visual.md`, "shipped preset" and every `preset` word outside the known files. No stale restatement remains. The only old-shape text is in `skills/do-work/CHANGELOG.md:34` and `:156`. Those lines are historical release notes for earlier versions, so they are correct as written. `skills/do-work-toolbox/SKILL.md:3` and `:25` route on "second opinion preset" and name no preset list. `ai-report.md:87` uses `bold` only as a file-name example. No Go code or `_dev` check reads `#another-layout` or counts the band's `<pre>` lines.

### Domain Review (general.md, prime-action-files.md)

- Cross-referencing: `ai-report-reference.md` (in `actions/`) links `../docs/second-opinion-presets/`, and `ai-report-guide.md` (in `docs/`) links `second-opinion-presets/`. Both are same-package paths at the correct depth, and both resolve in the source tree and in an install.
- Closed enumerations: the test's ten-name list is the requirement itself (the curated menu ships), and every other file is found by reading the folder. That matches the REQ. The guide's twelve-name list is a hand-maintained list (F1).
- Naming: the single-word preset tokens follow the existing `bold` and `visual` invocation form. They fall under the single-word-by-design exception.

### Test Review (`TestSecondOpinionShippedPresetsLoadFromTheToolbox`)

- **It pins a real failure.** It catches two real failures: a curated preset that is missing or renamed, and a shipped preset file the launcher refuses (empty, contains `{{`, not a regular file, or an invalid name). It also proves the real shipped prompt fills every placeholder when a preset is applied. Before this test, no test read the shipped files.
- **It works in an installed copy.** Five `..` steps from `internal/toolboxcommands` reach the skills root. In an install, that root is `.claude/skills/`, where `suite/modules.tsv` puts `do-work-toolbox` as a required sibling. The suiteinstall, publication, doctor and bkb_init tests already use the same walk. The test does not read `_dev/`.
- **Fragility:** F2 and F3. Map iteration order varies between runs, which matters only together with F2.

### Skeleton Review

- Tags are balanced (checked with an HTML parser: no unclosed or mismatched tag). The new comment contains no `--`, because the builder wrote "preset" and "editor-direction" without leading dashes (hand-back D-03). One older comment at `skeleton.html:32` contains `--from`, but that comment is outside this diff.
- The new paragraph and the new `<pre>` carry `FILL:` markers. The `<pre>` has two markers, one for the flag choice and one for the bundle folder. The Step 8 "no `FILL:` marker" check therefore catches an unfilled band.
- An author can fill the line correctly. Replacing the first marker gives `--preset feynman` or `--editor-direction '...'`, and replacing the second gives the bundle path. This matches `ai-report-reference.md:91`. The line starts with `ai-report second-opinion` in the same form as the existing layout line.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1 (Minor): `skills/do-work-toolbox/docs/ai-report-guide.md:122` replaces the closed `bold`/`visual` pair with a new hand-maintained list of all twelve preset names. Nothing checks that list against the folder, so the next shipped preset can be left out of the guide. The builder already flagged this in the hand-back's Discovered Tasks. A conditional wording ("the shipped presets live in `second-opinion-presets/`; `bold` and `visual` change the look, the others change how the page teaches") would avoid the drift. — impact-negligible → report only
- F2 (Minor): `report_second_opinion_test.go:1033` calls `secondOpinionRows`, which runs `t.Fatalf("missing typed backend summary")` when the launcher refuses a preset. The first broken shipped file therefore stops the loop. Because map order is random, which file gets reported changes between runs, and the message names the raw result instead of the preset. Checking `result.Outcome` and calling `t.Errorf` with the preset name would report every broken file at once. — impact-negligible → report only
- F3 (Nit): `report_second_opinion_test.go:1035` proves the direction was applied only through the absence of the literal `Direction: None.` from the shipped prompt. If the prompt's label wording changes, that check passes without checking anything. Asserting that the prompt contains the preset file's text would pin the real behavior. — impact-negligible → report only
- F4 (Nit): The REQ's `[UNIFY]` line (`REQ-712-...md:64`) says `skeleton.html` "gained exactly two `FILL:` markers". The file actually gained three `FILL:` tokens (47 before, 50 after) on two lines. This is a record mistake only, and the skeleton itself is correct. — impact-negligible → report only
- F5 (Nit): Only `80-20.md` and `quizme.md` say anything about corrections or open questions, while the prompt's DIRECTION paragraph lets a preset override every RULE, including Rule 2 (the Corrections section). An editor that reads a teaching preset literally could drop the Corrections section. FINAL REPLY still lists corrections, and `bold.md` has the same shape, so the risk is small and predates this change. — impact-negligible → report only

**Acceptance:** Pass — stages covered: Implementation and Integration. `go test ./internal/toolboxcommands -run SecondOpinion -count=1` passed in 16.1s on the merged tree. The real launcher `skills/do-work/tools/do-work-cli.sh ... ai-report-second-opinion --prepare --backend claude --model claude-opus-5-5 --preset <name> latest`, run in a scratch repository, returned `success prepared` with `preset_file` under the shipped folder for `feynman`, `80-20` and `socrates`. Deployment and Live acceptance were not exercised.
**Restatement sweep:** redefined the closing band (layout suggestion plus teaching-method suggestion) and the shipped-preset set (bold and visual to twelve presets). The sweep covered `skills/`, `suite/` and `tools/` and found no stale restatement. The CHANGELOG lines are historical.
**Suggested testing:** 4 items
- Live acceptance: unassessed. Generate a fresh `ai-report` and confirm the closing band shows both suggestions, the method fits the subject, and its line runs as written.
- Live acceptance: unassessed. Run one directed pass with a teaching preset (for example `--preset socrates` or `--preset quizme`) through a real backend. Confirm the page teaches in that shape and still includes Corrections, the backlink and the change list.
- Deployment: unassessed. After the finalization release, confirm an installed copy has all twelve files under `.claude/skills/do-work-toolbox/docs/second-opinion-presets/` and that `go test ./internal/toolboxcommands` passes there.
- Manual check: open a filled skeleton at phone width and confirm the second `<pre>` wraps like the first.
**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*
