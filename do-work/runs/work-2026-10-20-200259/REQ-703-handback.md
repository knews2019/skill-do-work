# REQ-703 hand-back (second opinions accept a named preset and a free-text editor direction)

- Branch: `worktree-agent-REQ-703-second-opinion-presets-and-direction`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-703-second-opinion-presets-and-direction`
- Base commit: `c6afcead` (checked before the first edit)
- Commits: `f06e82e0` `[REQ-703] add second-opinion presets and editor direction` (one commit, branch tip)
- Not done by me (integrator's): release (version, both CHANGELOGs, version.md), lesson bullet and lessons-index refresh, gate.

## File manifest

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified): `--preset` / `--editor-direction` parse rules and usage; `secondOpinionDirection` with `combinedText` and `filenameSlot`; preset resolution in the handler (project folder first, then the folder beside the prompt; regular file; trimmed text; unknown-name list); row fields; `secondOpinionRecord.Direction` (`omitzero`) restored by finalize; `secondOpinionPrompt` replaces both `strings.NewReplacer` fills; `secondOpinionNames` gains `variation`; `secondOpinionVariationSlug`; `secondOpinionDecorate` (strip, then metas and header block); publication tail order strip, backlink, slug, decorate, names, publish, judge.
- `.../toolboxcommands/report_second_opinion_test.go` (modified): fixture template gains `Direction: {{EDITOR_DIRECTION}}`; test helper `secondOpinionFixturePresets`; the two new tests.
- `.../toolboxcommands/report_index.go` (modified): constant `secondOpinionPresetsFolder` in the existing const block; `reportIndexWalk` skips it.
- `.../toolboxcommands/report_index_test.go` (modified): `TestAIReportIndexSkipsSecondOpinionPresetsFolder`.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified): `Preset`, `PresetFile`, `EditorDirection` on `AIReportSecondOpinionRun`, all omitempty.
- `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` (modified): DIRECTION paragraph, Rule 8 reduced motion, FINAL REPLY slug line.
- `skills/do-work-toolbox/docs/second-opinion-presets/bold.md` (new), `visual.md` (new): one plain paragraph each.
- `skills/do-work-toolbox/actions/ai-report.md` (modified): trigger bullet, flags in the Session (both the bullet head and step a), Targeted and Batch forms, directed names, metas, header block, row fields, finalize takes no direction flags.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified): three command lines, one sentence on presets, the winning project folder and committing first (D-04), backlink sentence adjusted for directed pages.
- `skills/do-work-toolbox/SKILL.md` (modified): `second opinion preset` in the description and the ai-report routing row. `## Second-opinion form` heading text is unchanged.

## P-A-U

- **[PLAN]:** Follow the brief's decided mechanism D-01 to D-09 exactly. RED first for the three named tests, then implement in one Go file plus the result row and the catalog skip, then the shipped text and prose. No lesson contradicted the plan: native-backend-identity (the unverified name gains the variation after `<model>-unverified`, nothing else changes), session-is-the-backend (the direction crosses prepare/finalize only in the record; finalize never re-reads the preset file), writer-anchored-marker-position (the header markers are read only as real comments on the masked surface; the editor's slug line is cleaned to `[a-z0-9-]`, never trusted).
- **[APPLY]:** As planned. Scope stayed inside the write boundary (11 paths, all listed in the brief).
- **[UNIFY]:** `git diff c6afcead --stat`:

```
 skills/do-work-toolbox/SKILL.md                    |   4 +-
 skills/do-work-toolbox/actions/ai-report.md        |  10 +-
 skills/do-work-toolbox/docs/ai-report-guide.md     |  10 +-
 .../docs/ai-report-second-opinion-prompt.md        |   7 +-
 .../docs/second-opinion-presets/bold.md            |   1 +
 .../docs/second-opinion-presets/visual.md          |   1 +
 .../internal/resultmodel/result_model.go           |   4 +
 .../internal/toolboxcommands/report_index.go       |   8 +-
 .../internal/toolboxcommands/report_index_test.go  |  25 +++
 .../toolboxcommands/report_second_opinion.go       | 238 +++++++++++++++++++--
 .../toolboxcommands/report_second_opinion_test.go  | 118 +++++++++-
 11 files changed, 388 insertions(+), 38 deletions(-)
```

Checks (from the worktree root unless noted, all at the final tree):

| Check | Exit | Wall time |
|---|---|---|
| `bash .../REQ-703-probe.sh` | 0 ("REQ-703 probe: shipped prompt and presets checked; preset, directed-run and catalog-skip tests passed") | 2 s |
| `gofmt -l skills/do-work/tools/do-work-cli/internal` | 0, empty output | <1 s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands ./internal/resultmodel` | 0 | <1 s |
| `go test -count=1 -run 'SecondOpinion\|AIReportIndex\|TestHandlersRegisterCanonicalToolboxCommands' ./internal/toolboxcommands` | 0 | 13-15 s |
| `go test -count=1 -run TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence ./internal/resultmodel` | 0 | <1 s |
| `git diff --cached --check` before commit | 0 | <1 s |

Each changed file was read back in the diff: no debug output, no leftover scratch, comments match the code. No `do-work/`, `ai-reports/`, CHANGELOG or VERSION path touched.

## Proof record

RED at base `c6afcead` (implementation absent, tests present):

1. `TestSecondOpinionPresetResolutionAndArgumentRules`: `report_second_opinion_test.go:848: missing typed backend summary ... "observed_evidence":["unexpected argument \"--preset\"; ai-report-second-opinion [...]"]` (the REQ's RED).
2. `TestSecondOpinionDirectedRunNamesAndMarksPage`: the full test does not compile at base (`undefined: secondOpinionDirection`, the unit calls to `secondOpinionDecorate`). With those two unit-call blocks removed for the run, the assertion RED is `report_second_opinion_test.go:894: missing typed backend summary ... unexpected argument "--preset"`. The test file was restored byte-for-byte afterwards.
3. `TestAIReportIndexSkipsSecondOpinionPresetsFolder`: `report_index_test.go:272: catalog bundles=[ai-reports/2026-10-20_1200_report ai-reports/second-opinion-presets], want only the report`.

GREEN:

- The three named tests: `--- PASS` each (0.02 s, 0.10 s, 0.70 s), no skips.
- Package run `SecondOpinion|AIReportIndex|TestHandlersRegisterCanonicalToolboxCommands`: 32 top-level `--- PASS`, zero FAIL, zero SKIP (plain names, unverified names, session prepare/finalize, banner, judge gate, catalog tests all green).
- `resultmodel.TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence`: PASS.

## Decisions

- **D-10** (DECIDE & STATE): with no opening body tag (possible only on an unverified page, where the backlink check is skipped), the header block goes at offset 0 of the page. Reason: smallest rule that keeps the direction visible; a verified page always has a body because the backlink check requires one.
- **D-11** (DECIDE & STATE): slug limits: the line pattern is `(?im)^[ \t]*[-*]?[ \t]*variation slug:(.*)$`, first match only; the value is lowercased and split on every run of characters outside `[a-z0-9]`; at most the first 4 words, then at most 40 characters with trailing `-` trimmed. Empty means no slug segment.
- **D-12** (DECIDE & STATE): header block markup is `<!-- second-opinion-header -->[metas when no </head>]<aside aria-label="Second-opinion direction" style="overflow-wrap:anywhere"><p><strong>Second-opinion preset:</strong> <slot></p><p style="white-space:pre-line"><strong>Direction:</strong> <combined text></p></aside><!-- /second-opinion-header -->`, with no surrounding newlines, so stripping removes exactly what was inserted and a second decoration gives the same bytes. `pre-line` keeps the blank line between preset and free text visible.
- **D-13** (DECIDE & STATE): `secondOpinionRecord.Direction` uses `json:"direction,omitzero"` (module floor is Go 1.24, which supports it) because `omitempty` has no effect on a struct field. Old records decode to the zero value, which is a plain run.
- **D-14** (DECIDE & STATE): the prompt wording. DIRECTION paragraph: "The direction below comes from the person directing this pass. When it is not "None.", follow it: it overrides ROLE, GOAL, every RULE and PROCESS, everything except BOUNDARY, which still applies in full. When it is "None.", this is a plain clarity pass." then `Direction: {{EDITOR_DIRECTION}}`. Rule 8 ends "and honor prefers-reduced-motion for any animation." FINAL REPLY ends "When a direction is given, end with one line: Variation slug: <2 to 4 plain words naming this variation>."
- **D-15** (DECIDE & STATE): the unknown-preset message says `presets found: none` when both folders are empty, so the message never ends in a bare colon.
- **D-16** (DECIDE & STATE): the preset-name regex and the slug regex are compiled inline where used (the house pattern in `secondOpinionBanner`), so no new package-level names were added.

## Discovered Tasks

- `optionValue` (`commands.go:87`) returns an empty string for `--x ''` for every valued flag; `--backend ''`, `--model ''` and `--finalize ''` therefore read as "not given" instead of a usage error. Only `--editor-direction` and `--preset` check for it now. impact-minor → report only.

## Lessons read

- `_dev/primes/lessons-releases.md` (whole file).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, families only: `native-backend-identity` (REQ-698, REQ-701), `session-is-the-backend` (REQ-700), `writer-anchored-marker-position` (REQ-528), `fixture-cost-is-subprocess-spawning` (REQ-574).
- Primes: `prime-do-work-cli.md` (Traps), `prime-action-files.md` (descriptions are triggers), `prime-releases.md` noted as the integrator's.
- Crew rules: general, coding-guardrails, shared-principles, communication-style, testing, backend.

## Anti-bloat check

`git diff --stat` is above. Added beyond the brief's named items:

- Test helper `secondOpinionFixturePresets` (test file): both new tests need the fixture's `bold.md` and `visual.md` beside the temp template; one helper instead of two copies.
- Option fields `preset`, `editorDirection` on `secondOpinionOptions`: the parsed values of the two named flags, required by the parser.
- Nothing else: no new environment variable, flag, file, config or package-level regex. The launcher `do-work-cli.sh` is untouched (D-01).

## Proposed CHANGELOG entry

**Second Opinions Accept Presets and Your Own Direction**

You can now steer a second opinion instead of always getting the one fixed clarity pass, and every directed page says what it was asked to do, in its name, at its top and in the result row.

- `ai-report second-opinion` accepts `--preset <name>` and `--editor-direction '<text>'`, alone or together, in the session (`--prepare`), targeted and batch forms. Together means preset text first, then your words. `--finalize` reads them from the prepare record.
- Two shipped presets, `bold` (add motion, add more sections, go further) and `visual` (diagrams and charts first). A committed `ai-reports/second-opinion-presets/<name>.md` wins over a shipped preset of the same name; an unknown name lists the presets found.
- The prompt gains a DIRECTION paragraph that overrides every editorial rule except the file boundary, asks for a `Variation slug:` line, and Rule 8 now asks to honor `prefers-reduced-motion`.
- Directed pages are named `index-<cli>-<model>[-unverified]-<preset|custom>[-<slug>][-N].html` with a matching reply, carry `ai-report-second-opinion-preset` and `ai-report-second-opinion-direction` meta tags and a marker-wrapped header block; result rows gain `preset`, `preset_file` and `editor_direction`.
- Plain runs keep today's names and behavior. `ai-report index`, `find`, `judge` and `latest` skip the project preset folder.

## Proposed lesson bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`.

- [family: writer-anchored-marker-position] [REQ-703: a launcher block that is placed ahead of a content check must be stripped before that check and re-inserted after it, on the masked surface, with no surrounding whitespace, so an editor copying a decorated page cannot fail the check and decorating twice yields the same bytes](../../../../do-work/archive/UR-160/REQ-703-second-opinion-preset-and-editor-direction.md)

## Integration seams for REQ-704

- Names: `secondOpinionNames(bundle, backend, model, variation)` is the only name builder; insert the tree id after `index-` there.
- Header block: `secondOpinionDecorate(page, direction)` builds the block between `<!-- second-opinion-header -->` and `<!-- /second-opinion-header -->`; add the tree id and parent link there. Its first half already strips any copied block and the two metas, and `secondOpinionPublication` calls it with the zero direction before the backlink check, so base pages that carry a header block pass.
- Direction value: `secondOpinionDirection` in `secondOpinionRecord.Direction`; the prompt fill is `secondOpinionPrompt`.

## Test wall times

Three named tests about 1 s together; probe 2 s; package selection 13-15 s; resultmodel <1 s. Machine was not otherwise loaded by me; no budget failure, no reruns needed.
