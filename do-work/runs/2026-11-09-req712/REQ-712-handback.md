# REQ-712 builder hand-back

Builder: the orchestrator session (Claude Opus 5.5), on branch `worktree-agent-REQ-712-second-opinion-suggestion-picks-best-learning-method`.

## File manifest

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified) — new `TestSecondOpinionShippedPresetsLoadFromTheToolbox`
- `skills/do-work-toolbox/docs/second-opinion-presets/{eli5,80-20,sticky-notes,feynman,mindmap,quizme,visual-studyguide,first-principles,socrates,schematics}.md` (new, ten files)
- `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` (modified) — closing band gains the teaching-method suggestion and its line; jump label renamed
- `skills/do-work-toolbox/docs/ai-report-kit/presentation.css` (modified) — two comments describing the closing band
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified) — the closing-band design rule
- `skills/do-work-toolbox/actions/ai-report.md` (modified) — trigger sentence, narrative item 9, kind sentence, Step 8 summary, checklist line
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified) — overview sentence and the preset list

## Red-green evidence

- RED (commit 3599f71f, test only): `go test ./internal/toolboxcommands -run TestSecondOpinionShippedPresetsLoadFromTheToolbox` failed with ten errors, `learning-method preset "eli5" does not ship in .../second-opinion-presets` and the same for the other nine. The real `bold` and `visual` loaded through the real shipped prompt.
- GREEN (commit 31acf248): the same test passes (0.26s).
- Negative check: a temporary `zz-broken.md` containing `{{OUTPUT_FILE}}` in the shipped folder made the test fail on that file; the file was removed before commit.
- Focused probe `bash do-work/runs/2026-11-09-req712/focused-probe.sh`: ok in 17.2s. `bash _dev/tests/contract-regressions.sh`: exit 0 in 25s.

## Decisions

- D-01 (DECIDE & STATE): Preset names are the REQ's suggestions: `eli5`, `80-20`, `sticky-notes`, `feynman`, `mindmap`, `quizme`, `visual-studyguide`, `first-principles`, `socrates`, `schematics`. They are typed invocation tokens like `bold` and `visual`, the single-word-by-design case in `coding-guardrails.md` § 5.
- D-02 (DECIDE & STATE): The teaching-method suggestion lives in the existing `#another-layout` band, not a new band. The `id` is kept because published reports link `index.html#another-layout` and `presentation.css` styles `#another-layout pre`, so the second line wraps with no CSS change. The heading becomes "Try another layout or teaching method" and the depth-1 jump label "Another layout or method".
- D-03 (DECIDE & STATE): The skeleton's FILL line for the method reads `ai-report second-opinion FILL: --preset and a shipped preset name, or --editor-direction and the method as one instruction in single quotes ai-reports/FILL: this bundle folder`, so the author must choose the flag. The HTML comment says "preset" and "editor-direction" without leading dashes, because `--` inside an HTML comment is a parse hazard.
- D-04 (DECIDE & STATE): Each preset ends with the fact rule ("Every fact still comes from the report or the repository"), like `bold.md`. `quizme.md` extends it to questions and answers. The ELI5 preset asks for concrete repository examples, not invented comparisons, so it cannot pull in unsourced facts.

## Discovered Tasks

- The guide's preset list in "A Second Opinion" names all twelve shipped presets by hand. Nothing checks it against the folder, so a future preset can be missed there. → report only

## Lesson evidence

- `_dev/primes/lessons-releases.md` read whole (bare entry). No entry was missing.
