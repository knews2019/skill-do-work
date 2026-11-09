---
id: REQ-712
title: 'Closing second-opinion suggestion picks the best learning method for the subject'
status: pending
created_at: 2026-11-09T12:02:46Z
user_request: UR-164
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set: [skills/do-work-toolbox/docs/second-opinion-presets/*.md, skills/do-work-toolbox/docs/ai-report-kit/skeleton.html, skills/do-work-toolbox/actions/ai-report-reference.md, skills/do-work-toolbox/actions/ai-report.md, skills/do-work-toolbox/docs/ai-report-guide.md, skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go]
---
# Closing Second-Opinion Suggestion Picks the Best Learning Method for the Subject
## What
When `ai-report` finishes a report and its closing band suggests a second opinion on that same report, the suggestion must also name the learning or presentation method that would best help the user understand or deliver that report's subject, with a ready `ai-report second-opinion` line that applies it. The user's ten "accelerated AI prompts" ship as named second-opinion presets so the line can be `--preset <name>`; they are a curated starting menu, not a closed list, and the report author may pick a better method outside it through `--editor-direction`.
## Why
The user's words: "make sure to pick the best prompt that would have help me understand or deliver that subject." In capture the user added that the list "is just a subset of curated lists, there could be better methods and the model should not be limited by them", and that the aim is "to expand it by consider what is the best way to present it".
## Detailed Requirements
- Ship one preset file per method under `skills/do-work-toolbox/docs/second-opinion-presets/`, beside the existing `bold.md` and `visual.md`. Suggested names, each a valid preset token (lowercase letters, digits, hyphens): `eli5`, `80-20`, `sticky-notes`, `feynman`, `mindmap`, `quizme`, `visual-studyguide`, `first-principles`, `socrates`, `schematics`. The user's gloss for Socrates is "teach by asking questions"; the Feynman line reads "help me understand the" (the subject). Each preset body is a direction to the clarity editor in the same voice as `bold.md`: what the rewritten page does for the reader, with every fact still from the report or the repository. Single-word preset tokens follow the existing `bold`/`visual` invocation form (single-word-by-design carve-out in `skills/do-work/crew-members/coding-guardrails.md` § 5).
- The closing band of every report (today "Try another layout" in `skills/do-work-toolbox/docs/ai-report-kit/skeleton.html`, ruled by `skills/do-work-toolbox/actions/ai-report-reference.md` "End with another layout") keeps its layout suggestion and its line, and gains a second suggestion: the one learning or presentation method that best fits this report's subject and reader, with one or two sentences saying why it fits, and a ready line. The band title and intro text may change to cover both suggestions.
- The rule for choosing the method is a condition, not a list: pick what best helps the reader understand or deliver this subject. The shipped presets are the illustrative menu; the author may name a method outside it. When the pick is a shipped preset, the line is `ai-report second-opinion --preset <name> ai-reports/<bundle>`; otherwise it is `--editor-direction '<the method as one instruction>'`, written without a single quote so it stays one argument (same constraint as the layout line).
- `ai-report` Step 8's printed summary (`skills/do-work-toolbox/actions/ai-report.md`, the "layout suggestion with its line" sentence) also prints the method pick and its line.
- `skills/do-work-toolbox/docs/ai-report-guide.md` lists the shipped presets and states that "second opinion with <preset>" works for any shipped preset name, not only `bold`.
- A Go test pins that every shipped preset file under the real `skills/do-work-toolbox/docs/second-opinion-presets/` loads through `--preset <name>` (regular file, non-empty, no `{{`, name matches the preset regex). Today's test in `report_second_opinion_test.go` fakes that folder with `secondOpinionFixturePresets`, so the shipped files are not covered. Prefer deriving the expected set from the folder listing over hardcoding ten names; the existing `presets found: bold, visual` assertion is a closed enumeration that must not break on every new preset.
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
**GREEN when:** Each of the ten shipped presets resolves (state `prepared`, `preset_file` under the shipped folder); a new report's closing band names one method chosen for this report's subject, says in one or two sentences why it fits, and carries a ready `ai-report second-opinion` line (`--preset <name>` for a shipped preset, `--editor-direction` otherwise) beside the unchanged layout suggestion; Step 8's summary prints the same pick and line; the Go test over the real shipped folder passes.
**Validation:** User confirmed, with the adjustment that the pick may be a method outside the shipped list.
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24067 tokens, `slugged: partial`, over the 2000-token budget and not targetable; matched on family `closed-enumeration-for-a-condition` (the hardcoded `presets found: bold, visual` assertion) and `session-is-the-backend` (the Session form this proof uses).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
## Full Context
See `do-work/user-requests/UR-164/input.md` for complete verbatim input.
*Source: Given these accelerated AI prompts, when you make an AI report and then suggest a new second opinion on the same report, make sure to pick the best prompt that would have help me understand or deliver that subject. [ten-method list follows in the UR]*
