---
id: REQ-712
title: 'Closing second-opinion suggestion picks the best learning method for the subject'
status: completed
created_at: 2026-11-09T12:02:46Z
user_request: UR-164
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: [_dev/primes/lessons-releases.md]
estimate:
  p50_active_minutes: 60
  confidence: medium
  calculated_at: 2026-11-09T12:05:07Z
  basis:
    - Route C
    - 6-file write set
    - 10 new files
    - 2 subsystems involved
    - 6 acceptance criteria
    - cross-route regression gates
write_set: [skills/do-work-toolbox/docs/second-opinion-presets/eli5.md, skills/do-work-toolbox/docs/second-opinion-presets/80-20.md, skills/do-work-toolbox/docs/second-opinion-presets/sticky-notes.md, skills/do-work-toolbox/docs/second-opinion-presets/feynman.md, skills/do-work-toolbox/docs/second-opinion-presets/mindmap.md, skills/do-work-toolbox/docs/second-opinion-presets/quizme.md, skills/do-work-toolbox/docs/second-opinion-presets/visual-studyguide.md, skills/do-work-toolbox/docs/second-opinion-presets/first-principles.md, skills/do-work-toolbox/docs/second-opinion-presets/socrates.md, skills/do-work-toolbox/docs/second-opinion-presets/schematics.md, skills/do-work-toolbox/docs/ai-report-kit/skeleton.html, skills/do-work-toolbox/docs/ai-report-kit/presentation.css, skills/do-work-toolbox/actions/ai-report-reference.md, skills/do-work-toolbox/actions/ai-report.md, skills/do-work-toolbox/docs/ai-report-guide.md, skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go]
claimed_at: 2026-11-09T12:07:52Z
route: B
dispatch_at: 2026-11-09T12:14:02Z
builder_handback_at: 2026-11-09T12:16:27Z
integration_at: 2026-11-09T12:16:51Z
review_at: 2026-11-09T12:22:01Z
kb_status: pending
commit: 0c06eb10993bf812bfbce4467035b118c5992035
heavy_verified_at: 2026-11-09T12:24:54Z
heavy_verified_revision: 0c06eb10993bf812bfbce4467035b118c5992035
completed_at: 2026-11-09T12:25:03Z
release_at: 2026-11-09T12:25:03Z
---
# Closing Second-Opinion Suggestion Picks the Best Learning Method for the Subject
## What
When `ai-report` finishes a report and its closing band suggests a second opinion on that same report, the suggestion must also name the learning or presentation method that would best help the user understand or deliver that report's subject, with a ready `ai-report second-opinion` line that applies it. The user's ten "accelerated AI prompts" ship as named second-opinion presets so the line can be `--preset <name>`; they are a curated starting menu, not a closed list, and the report author may pick a better method outside it through `--editor-direction`.
## Why
The user's words: "make sure to pick the best prompt that would have help me understand or deliver that subject." The user's own label for the list is "Codes to facilitate learning in Ai bots". In capture the user added that the list "is just a subset of curated lists, there could be better methods and the model should not be limited by them", and that the aim is "to expand it by consider what is the best way to present it".
## Detailed Requirements
- Ship one preset file per method under `skills/do-work-toolbox/docs/second-opinion-presets/`, beside the existing `bold.md` and `visual.md`. Suggested names, each a valid preset token (lowercase letters, digits, hyphens): `eli5`, `80-20`, `sticky-notes`, `feynman`, `mindmap`, `quizme`, `visual-studyguide`, `first-principles`, `socrates`, `schematics`. The user's gloss for Socrates is "teach by asking questions"; the Feynman line reads "help me understand the" (the subject). Each preset body is a direction to the clarity editor in the same voice as `bold.md`: what the rewritten page does for the reader, with every fact still from the report or the repository. Single-word preset tokens follow the existing `bold`/`visual` invocation form (single-word-by-design carve-out in `skills/do-work/crew-members/coding-guardrails.md` § 5).
- The closing band of every report (today "Try another layout" in `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html`, ruled by `skills/do-work-toolbox/actions/ai-report-reference.md` "End with another layout") keeps its layout suggestion and its line, and gains a second suggestion, in the same band or a band of its own (builder's call): the one learning or presentation method that best fits this report's subject and reader, with one or two sentences saying why it fits, and a ready line. The band title and intro text may change to cover both suggestions.
- The rule for choosing the method is a condition, not a list: pick what best helps the reader understand or deliver this subject. The shipped presets are the illustrative menu; the author may name a method outside it. When the pick is a shipped preset, the line is `ai-report second-opinion --preset <name> ai-reports/<bundle>`; otherwise it is `--editor-direction '<the method as one instruction>'`, written without a single quote so it stays one argument (same constraint as the layout line).
- `ai-report` Step 8's printed summary (`skills/do-work-toolbox/actions/ai-report.md`, the "layout suggestion with its line" sentence) also prints the method pick and its line.
- `skills/do-work-toolbox/docs/ai-report-guide.md` today names only `--preset bold` and `--preset visual` in its "A Second Opinion" section. Replace that closed pair with wording that covers every shipped preset, and say where the shipped presets live. The trigger sentence in `skills/do-work-toolbox/actions/ai-report.md` ("Second opinion with bold" adds `--preset bold`) states the same rule for any shipped preset name.
- A Go test reads the real `skills/do-work-toolbox/docs/second-opinion-presets/` folder and pins two things. First, each of the user's ten methods resolves through `--preset <name>` under the name the builder picked; this named list is the requirement itself (the curated menu ships), so it is not a closed enumeration of a condition, and it is what makes the test RED today. Second, every `.md` file in that folder loads (regular file, non-empty, no `{{`, name matches the preset regex), so a future preset is covered without editing the test. Today's tests in `report_second_opinion_test.go` point `DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT` at a temp-dir prompt and fake the preset folder with `secondOpinionFixturePresets`, so no test reads the shipped files. Those fixture tests are isolated from the real folder, and new shipped presets do not break their `presets found: bold, visual` assertion.
## Constraints
- Do not remove the existing layout suggestion or its line (user chose "add beside it").
- Presets live in the installed toolbox, so the uncommitted-project-preset refusal from REQ-705 (confirm uncommitted preset refusal) does not apply to them; do not touch that check.
- The band is prose plus skeleton; no Go code reads `#another-layout`, so the Go change is the shipped-preset test only unless the builder finds a reason otherwise.
- This is a change to shipped files, so it is a release per `_dev/primes/prime-releases.md`.
## Builder Guidance
Certainty: high on the ten presets and the beside-not-replace band. Latitude: the preset names, the band's new title and wording, and whether the method pick gets its own band or a second paragraph in the existing one, are the builder's call; keep the skeleton's FILL markers so the judge still catches an unfilled band. Write each preset so it changes how the page teaches (Socrates asks the reader questions before revealing answers; quizme ends sections with checks; schematics draws the mechanism), not only how it looks.
## Context
Prior work: REQ-703 (second opinions accept a named preset and a free-text editor direction, commit 1d68e403) built `--preset`, `--editor-direction`, the two shipped presets, and the "Try another layout" band. REQ-704 (second-opinion linked iteration tree) made every directed run a tree page under `index.html`. The preset loader reads `<bundle>/second-opinion-presets/<name>.md` first, then the shipped folder beside `ai-report-second-opinion-prompt.md`; a new shipped file needs no loader change. The prompt's DIRECTION paragraph lets a preset override the clarity-editor rules except OUTPUT, BACKLINK, CHANGE LIST, BOUNDARY and FINAL REPLY, so a teaching-shaped preset is legitimate.
## Red-Green Proof
**RED prompt/case:** Run `ai-report second-opinion --prepare --backend claude --model <slug> --preset feynman latest` on any existing bundle, and open the closing band of a freshly generated report.
**Why RED now:** The command is refused with `unknown preset "feynman"; presets found: bold, visual`, and the band names only a layout, never a learning or presentation method.
**GREEN when:** Each of the ten methods resolves as a shipped preset under the name the builder picked (state `prepared`, `preset_file` under the shipped folder); a new report's closing band names one method chosen for this report's subject, says in one or two sentences why it fits, and carries a ready `ai-report second-opinion` line (`--preset <name>` for a shipped preset, `--editor-direction` otherwise) beside the unchanged layout suggestion; Step 8's summary prints the same pick and line; the Go test over the real shipped folder fails before the preset files exist and passes after.
**Validation:** User confirmed, with the adjustment that the pick may be a method outside the shipped list.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8581 tokens, `slugged: partial`; over the 2000-token budget after `_dev/primes/lessons-releases.md` (666), not targetable. Matched because `prime_files` names `_dev/primes/prime-action-files.md` and the REQ edits `ai-report.md` and `ai-report-reference.md`.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24067 tokens, `slugged: partial`; over budget, not targetable. Matched on family `closed-enumeration-for-a-condition` (the hardcoded `presets found: bold, visual` assertion) and `session-is-the-backend` (the Session form this proof uses).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-action-files.md`, `prime-releases.md` and `lessons-releases.md`. Approach: RED test over the real shipped preset folder first, then ten preset files in the voice of `bold.md`, then move every description of the closing band to two suggestions in the same `#another-layout` band.
- [x] **[APPLY]:** RED commit 3599f71f (test only), GREEN commit 31acf248 (presets and prose). Every touched file is in `## Scope`.
- [x] **[UNIFY]:** `git diff --stat aa5d73e0..31acf248` shows 17 files. Checked: the test file passes gofmt; the ten presets contain no `{{` and each ends with the fact rule; `skeleton.html` parses with balanced tags and gained three `FILL:` markers on two new lines (corrected after review F4); `presentation.css`, `ai-report-reference.md`, `ai-report.md` and `ai-report-guide.md` were swept for stale one-suggestion wording. `contract-regressions.sh` exit 0. No debug artifacts.
## Full Context
See `do-work/user-requests/UR-164/input.md` for complete verbatim input.
*Source: Given these accelerated AI prompts, when you make an AI report and then suggest a new second opinion on the same report, make sure to pick the best prompt that would have help me understand or deliver that subject. [ten-method list follows in the UR]*

## Triage

**Route: B** - Medium

**Reasoning:** The REQ names every file and the outcome is clear, but the ten preset texts must match the voice of the two shipped presets and the new Go test must reuse the existing second-opinion test harness, which needs a short exploration. The verified Detailed Requirements already serve as the plan.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Closing band.** `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` lines 177-185 hold the `#another-layout` band: a FILL paragraph for one layout idea and a `<pre><code>` line `ai-report second-opinion --editor-direction 'FILL: ...' ai-reports/FILL: ...`. Six other places describe the band and must move with it: `ai-report-reference.md` line 91 ("End with another layout"), `ai-report.md` line 164 (narrative item 9), line 166 (the `--kind` sentence), line 194 (Step 8 summary), line 220 (checklist), and the `presentation.css` header comment line 17. `presentation.css` line 365 styles `#another-layout pre`, so a second `<pre>` inside the same band is styled with no CSS change. Two published reports under `ai-reports/` and one second-opinion page link `index.html#another-layout`, so the band keeps its `id`.

**Presets.** `skills/do-work-toolbox/docs/second-opinion-presets/` holds `bold.md` and `visual.md`: one paragraph each, written as a direction to the editor, with the fact rule restated in `bold.md`. The loader (`report_second_opinion.go` lines 195-228) reads `<bundle>/second-opinion-presets/<name>.md`, then the folder beside the trusted prompt; names must match `^[a-z0-9][a-z0-9-]*$`; a body must be non-empty and free of `{{`. No loader change is needed. The shipped prompt's DIRECTION paragraph lets a preset override ROLE, GOAL, RULES and PROCESS but never OUTPUT, BACKLINK, CHANGE LIST, BOUNDARY or FINAL REPLY.

**Trigger and guide.** `ai-report.md` line 22 says "Second opinion with bold" adds `--preset bold`. `ai-report-guide.md` line 122 names only `--preset bold` and `--preset visual`.

**Test harness.** `report_second_opinion_test.go`: `secondOpinionFixture` builds a temp repository and points `DO_WORK_AI_REPORT_SECOND_OPINION_PROMPT` at a temp prompt; `secondOpinionInvoke` calls the handler; `secondOpinionRows` decodes the run rows. A new test can call `secondOpinionFixture`, then re-point the env var at the real shipped prompt `../../../../../do-work-toolbox/docs/ai-report-second-opinion-prompt.md` so the loader reads the real shipped folder. That walk stays inside `skills/`, and `suite/modules.tsv` installs `do-work-toolbox` as a required sibling, so the path exists in an installed copy too. This avoids the `shipped-module-test-self-containment` trap, which is about reads of `_dev/` (the same pattern already ships in `suiteinstall` and `publication` tests reading `do-work-board/justfile.template`).

*Explored by work action (orchestrator, inline)*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/docs/second-opinion-presets/eli5.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/80-20.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/sticky-notes.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/feynman.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/mindmap.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/quizme.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/visual-studyguide.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/first-principles.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/socrates.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/schematics.md` (new)
- `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` (modify) — the closing band gains the method pick and its line
- `skills/do-work-toolbox/docs/ai-report-kit/presentation.css` (modify) — header comment that describes the closing band
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modify) — the closing-band design rule
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — trigger sentence, narrative item 9, kind sentence, Step 8 summary, checklist
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — preset list in the second-opinion section
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify) — test over the real shipped presets

**Files I will NOT touch:** `report_second_opinion.go` (the loader already reads a new shipped file), `bold.md` and `visual.md`, the shipped prompt `ai-report-second-opinion-prompt.md`, existing reports under `ai-reports/`, and the REQ-705 uncommitted-preset check.

**Acceptance criteria (restated from REQ):**
- [ ] Each of the user's ten methods ships as a preset under the shipped folder and resolves through `--preset <name>`.
- [ ] Each preset changes how the page teaches, keeps every fact from the report or the repository, and matches the voice of `bold.md`.
- [ ] The closing band keeps its layout suggestion and its line, and adds the one method that best fits this report's subject and reader, why it fits, and a ready line.
- [ ] The method rule is a condition: shipped presets are the illustrative menu, and a method outside it uses `--editor-direction` with no single quote.
- [ ] Step 8's printed summary prints the method pick and its line.
- [ ] The guide no longer names only `bold` and `visual`, and the trigger sentence covers any shipped preset name.
- [ ] A Go test reads the real shipped folder, names the ten methods, loads every `.md` file there, fails before the presets exist, and passes after.

## Pre-Flight

**Git:** ✓ Clean outside `do-work/` at 773bde8c; only this REQ's working file, `baseline.json` and the run directory are changed
**Tests baseline:** ✓ Focused probe green (`bash do-work/runs/2026-11-09-req712/focused-probe.sh`, second-opinion tests, 17s); repository gate `bash _dev/tests/maintainer-verify.sh` exit 0 in 157s at 773bde8c, run from a detached checkout
**Dependencies:** ✓ None added

*Checked by work action*

## Decisions

- D-01 (DECIDE & STATE): Preset names are the REQ's suggestions: `eli5`, `80-20`, `sticky-notes`, `feynman`, `mindmap`, `quizme`, `visual-studyguide`, `first-principles`, `socrates`, `schematics`. They are typed invocation tokens like `bold` and `visual`, the single-word-by-design case in `coding-guardrails.md` § 5.
- D-02 (DECIDE & STATE): The teaching-method suggestion lives in the existing `#another-layout` band, not a new band. The `id` is kept because published reports link `index.html#another-layout` and `presentation.css` styles `#another-layout pre`, so the second line wraps with no CSS change. The heading becomes "Try another layout or teaching method" and the depth-1 jump label "Another layout or method".
- D-03 (DECIDE & STATE): The skeleton's FILL line for the method reads `ai-report second-opinion FILL: --preset and a shipped preset name, or --editor-direction and the method as one instruction in single quotes ai-reports/FILL: this bundle folder`, so the author must choose the flag. The HTML comment says "preset" and "editor-direction" without leading dashes, because `--` inside an HTML comment is a parse hazard.
- D-04 (DECIDE & STATE): Each preset ends with the fact rule ("Every fact still comes from the report or the repository"), like `bold.md`. `quizme.md` extends it to questions and answers. The ELI5 preset asks for concrete repository examples, not invented comparisons, so it cannot pull in unsourced facts.

## Discovered Tasks

- The guide's preset list in "A Second Opinion" names all twelve shipped presets by hand. Nothing checks it against the folder, so a future preset can be missed there. Reviewer F1 found the same. impact-negligible → report only

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/docs/second-opinion-presets/eli5.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/80-20.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/sticky-notes.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/feynman.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/mindmap.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/quizme.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/visual-studyguide.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/first-principles.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/socrates.md` (new)
- `skills/do-work-toolbox/docs/second-opinion-presets/schematics.md` (new)
- `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html` (modified)
- `skills/do-work-toolbox/docs/ai-report-kit/presentation.css` (modified)
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified)

**What was done:** Ten learning-method presets now ship beside `bold` and `visual`, each a one-paragraph direction that changes how the second-opinion page teaches while keeping every fact sourced. The report kit's closing band keeps its layout suggestion and adds the one teaching method that best fits the subject, with a ready `--preset` or `--editor-direction` line; the design rule, the action's narrative, Step 8 summary, checklist, trigger sentence, CSS comments and the guide all describe both suggestions. A new Go test reads the real shipped folder, names the ten methods, and loads every preset file there through the launcher.

## Qualification

**Range:** `7dc47c2d..0c06eb10` (merge of the REQ branch: RED 3599f71f, GREEN 31acf248)
**Mechanical gates:** qualify ✓ satisfied, scope-drift ✓ satisfied (16 declared paths, 16 touched). The first scope-drift run reported every path undeclared because the Scope bullets carried bare paths; the bullets were backticked and the rerun was satisfied, with no file-list change.

**Requirement trace (orchestrator, read against the diff):**
- Ten presets ship under `skills/do-work-toolbox/docs/second-opinion-presets/` with the REQ's names; each is a non-empty paragraph with no `{{`, ends with the fact rule, and says how the page teaches (questions first for `socrates`, revealed answers for `quizme`, labelled drawings for `schematics`).
- The closing band keeps the layout paragraph and its `--editor-direction` line unchanged and adds a FILL paragraph for the teaching method plus a FILL line that makes the author choose `--preset` or `--editor-direction`. The `id="another-layout"` stays, so existing links and the `#another-layout pre` style still apply.
- The condition rule lives in `ai-report-reference.md`: the shipped presets are "a starting menu, not a closed list", and a better method uses `--editor-direction`. Step 8's summary sentence, narrative item 9, the `--kind` sentence and the checklist line in `ai-report.md` name both suggestions.
- The guide replaces the closed `--preset bold`, `--preset visual` pair with `--preset <name>` and lists the shipped presets; the trigger sentence in `ai-report.md` covers any shipped or bundle preset name.
- The Go test reads the real folder, fails RED on the missing ten, passes GREEN, and fails on a broken extra file (negative check in the hand-back).

**Live flow:** the test drives the real launcher handler with the real shipped prompt, so the presets are proven loadable end to end through `--prepare`, not only present on disk.

*Qualified by work action*

## Testing

**Tests run:** `bash do-work/runs/2026-11-09-req712/focused-probe.sh` (the `SecondOpinion` tests in `internal/toolboxcommands`), then `bash _dev/tests/maintainer-verify.sh` at the merge revision 0c06eb10 from a detached checkout
**Result:** ✓ Focused probe ok in 17.2s; repository gate exit 0, gate wall 143s (do-work-cli module 66s, 972 tests; slowest file `finalization_recovery_test.go` 20.1s, under the 30s budget). `_dev/tests/contract-regressions.sh` exit 0 in 25s on the builder branch.

**Red-green validation:**
- `TestSecondOpinionShippedPresetsLoadFromTheToolbox` (`report_second_opinion_test.go`): ✗ at commit 3599f71f, ten errors `learning-method preset "<name>" does not ship` → ✓ at commit 31acf248. This is the captured RED case (`--preset feynman` refused because the preset does not ship) and the GREEN condition "each of the ten methods resolves as a shipped preset" (state `prepared`, `preset_file` under the shipped folder), driven through the real launcher and the real shipped prompt.
- Negative check: a temporary `zz-broken.md` containing `{{OUTPUT_FILE}}` in the shipped folder made the same test fail on that file; removed before commit.
- The band part of GREEN (a new report names one method, why it fits, and a ready line) is prose and skeleton, proven by the review's read of the skeleton and design rule, not by an automated test.

**New tests added:**
- `TestSecondOpinionShippedPresetsLoadFromTheToolbox`

**Heavy verification plan:**
- Range: 7dc47c2d385288876d0196bb7df1959b67c6a2f3..0c06eb10993bf812bfbce4467035b118c5992035
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the test file matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — every changed path matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the test file matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the test file matched subtree `skills/do-work/tools/do-work-cli`

*Verified by work action*

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

## Heavy Verification Plan

- Base revision: 7dc47c2d385288876d0196bb7df1959b67c6a2f3
- Target revision: 0c06eb10993bf812bfbce4467035b118c5992035
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the test file matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — every changed path matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the test file matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the test file matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: 0c06eb10993bf812bfbce4467035b118c5992035
- Execution revision: 0c06eb10993bf812bfbce4467035b118c5992035 (detached checkout)
- do-work-cli-integrations: exit 0, 70s, executed
- staged-skills: exit 0, 38s, executed
- updater: exit 0, 68s, executed
- installer: exit 0, 29s, executed

## Lessons Learned

**What worked:** Writing the test against the real shipped prompt and preset folder, through the launcher handler, instead of the fixture prompt. It proved the presets load end to end, and the RED run showed the real `bold` and `visual` already loading, so the failure was only the missing ten. Launching the reviewer right after qualification, while the gate ran, cost no extra wall time.
**What didn't:** The first scope-drift check reported all sixteen files undeclared, because the `## Scope` bullets held bare paths. The checker only reads backticked paths in "Files I will touch".
**Worth knowing:** The closing band is described in seven places: the skeleton band and its depth-1 jump label, two `presentation.css` comments, the design rule in `ai-report-reference.md`, four spots in `ai-report.md` (trigger, item 9, the `--kind` sentence, Step 8, checklist), and the guide's overview. Change them together. The guide now lists every shipped preset by hand (F1), so a new preset must be added there too. Keep `id="another-layout"`, because published reports link to it.

## Orientation

Now every AI report ends with two second-opinion suggestions: another layout and the one teaching method (ELI5, Feynman, Socrates and the rest, or a better one) that best helps a reader understand or deliver the subject, each with a ready command. Ten teaching presets ship in the toolbox's second-opinion preset folder. This lives in the ai-report kit and second-opinion subsystem of `do-work-toolbox`. Prime staleness check: `prime-action-files.md` and `prime-releases.md` name no path this change moved or removed.

## Timing

Observed 2026-11-09T12:14:02Z to 2026-11-09T12:24:54Z: 10m 52s total, 10m 14s attributed across 3 events, 38s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 7m 44s | 2 |
| builder-work | 2m 30s | 1 |

Slowest stage: verification-gate / repository gate at merge 0c06eb10, 3m 57s, outcome success.
