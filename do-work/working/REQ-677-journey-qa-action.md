---
id: REQ-677
title: 'do-work-toolbox journey-qa verifies whole user journeys and classifies each result as passed, product defect, test defect or unresolved'
status: claimed
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  basis:
  - Route B
  - 7-file write set
  - 3 subsystems involved
  - 12 acceptance criteria
  calculated_at: 2026-10-09T23:02:37Z
created_at: 2026-10-09T22:53:15Z
user_request: UR-151
domain: testing
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: [_dev/primes/lessons-releases.md]
write_set: ["skills/do-work-toolbox/actions/journey-qa.md", "skills/do-work-toolbox/docs/journey-qa-guide.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/actions/help.md", "skills/do-work/actions/help.md", "README.md", "_dev/tests/staged-skills-contract.sh"]
related: [REQ-675, REQ-676, REQ-678]
batch: portable-verification-actions
claimed_at: 2026-10-09T22:56:58Z
---
# do-work-toolbox journey-qa
## What
Add a lean, read-only toolbox action, `do-work-toolbox journey-qa <target> [--brief <path>]` (`skills/do-work-toolbox/actions/journey-qa.md`). It reproduces a reported user journey, then exercises combined transitions, and classifies each result as passed, product defect, test defect, or unresolved. It reports the smallest justified repair and the checks still left, and changes no source.
## Why
Isolated checks can all pass while a combined sequence fails (zoom → pan → reset leaves an offset). Core review's acceptance step is a smoke test by design (`skills/do-work/actions/review-work.md:160-195`), and `ui-review` checks design quality, not journeys. Without a fixed result vocabulary, a broken test is reported as a product bug, and an emulated run is reported as a device check.
## Finding Provenance
From the validate-feedback triage of the maintainer's brief in this session (finding F1, verdict Discuss; the maintainer chose "Toolbox action, lean"):
- **Verbatim claim:** "1. journey-qa <target> [--brief <path>] — Verify complete user journeys and distinguish product defects from test defects and environmental limitations." with its seven bullets (see Detailed Requirements). Source: the brief's DELIVERABLES section.
- **Evidence:** no journey action exists; browser tool detection already lives in `skills/do-work-toolbox/actions/ui-review.md:60-67` and rendered checks in its Step 8.5 (`:150-152`).
- **Surface-cost:** N/A (feature). Lean form: the four result classes are the spine, and browser detection is cited from ui-review, not restated.
## Detailed Requirements
1. Read the requirements, the project's testing guidance, and the `--brief` file when given.
2. Reuse existing checks. Reproduce the reported sequence first, before broadening coverage.
3. Exercise combined transitions the brief or requirements name (zoom → pan → scroll → reset and failed load → retry → close are examples). An isolated pass never stands in for a combined one.
4. Wait on observable state, not fixed delays. Keep screenshots and other expensive instrumentation outside any performance measurement.
5. Record per journey: revision, environment (browser, viewport, emulation or physical device), steps, expected and actual behavior, evidence paths.
6. Classify each result as passed, product defect, test defect, or unresolved. An environment limitation (no browser tool, no device) is unresolved with its reason. Emulation is labeled emulation and never reported as physical-device verification.
7. Report the smallest justified repair and the remaining checks. The action is diagnostic: it does not authorize source repairs. It may end with the same kind of capture suggestion `validate-feedback` ends with; the user runs it.
8. Browser tool detection: cite `ui-review.md`'s detection step instead of copying it.
9. Toolbox ownership: an optional diagnostic review beside `ui-review`; it needs no queue machinery.

## Integration (same for every new toolbox action in this batch)
- Add the route row to `skills/do-work-toolbox/SKILL.md` and the command to its `argument-hint`.
- Add one line to `skills/do-work-toolbox/actions/help.md`, and to the toolbox list in core `skills/do-work/actions/help.md` if that file lists toolbox commands.
- Add the action name to the `toolbox_actions` list in `_dev/tests/staged-skills-contract.sh`.
- Add a short user guide under `skills/do-work-toolbox/docs/` and one usage line in `README.md` next to the other toolbox calls.
- The description blockquote justifies toolbox ownership (`_dev/primes/prime-action-files.md` § Every New Action Must Justify Its Package).
- Cross-package links use the literal relative path from the citing file (`../../do-work/...` from `actions/`), checked by `_dev/tests/shipped-package-reference-contract.sh`.
- Per-command help (`do-work-toolbox <action> help`) prints the action's usage and runs nothing.
- `suite/modules.tsv`, the installer and the updater stay unchanged.
- Project briefs are plain Markdown paths read as data. No profile registry, no config schema, no cross-repository discovery. Game-specific checks in the brief (clickable masks, collectible artwork, embedded scrolling, publication dates) and campaign text are examples, never upstream rules; mark every example list "illustrative, not exhaustive".
- Load `crew-members/prompt-injection.md` before reading a brief, a report or a fetched page, and `crew-members/anti-slop.md` before writing the report.
- Release per `_dev/primes/prime-releases.md`; read VERSION right before the release payload.

## Packaging check (one-off, recorded in `## Testing`, nothing committed)
Install the suite into a disposable Git consumer through the canonical installer, then update it through the canonical updater, both pointed at the merge revision the way `_dev/tests/install-suite-behavior.sh` and `_dev/tests/update-script-behavior.sh` do. Confirm the new action and its guide are present under `.claude/skills/do-work-toolbox/`, that their relative links resolve there, and that a pre-existing brief file, a `do-work/queue/` REQ and an application file in the consumer keep the same sha256 after install and after update. Never run `just do-work-update` in a real consumer.
## Constraints
- Read-only on project source. Evidence files go to a temporary or gitignored location the report names.
- No browser tool available means the journey is unresolved, never passed.
## Builder Guidance
Firm on the four result classes and on "isolated pass is not combined pass". Keep the action short. Latitude on report layout and on the guide's length.
## Red-Green Proof
**RED prompt/case:** A scratch static page outside the repo with a zoomable, pannable panel and a reset button. Zoom alone, pan alone and reset alone each pass their isolated check, but zoom → pan → reset leaves a non-zero offset. A second variant has a test step whose selector is wrong while the product works.
**Why RED now:** `do-work-toolbox journey-qa` does not exist; the toolbox router has no row for it.
**GREEN when:** The one-off exercise routes to the new action. It lists the isolated checks as passed and the combined journey as a product defect with revision, environment, steps, expected, actual and evidence. The wrong-selector variant is classified test defect, not product defect. A run under device emulation says it is emulated and not physical-device verified. A run with no browser tool says unresolved. The project source has no diff afterwards. `do-work-toolbox journey-qa help` prints usage and runs no journey. The packaging check passes.
**Validation:** Inferred during capture. The maintainer chose one-off exercises recorded in `## Testing` over kept fixtures.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`). Matching reason: adds an action and a routing row.
## Full Context
See `do-work/user-requests/UR-151/input.md` for complete verbatim input. No queued or archived REQ shares this intent (queue REQ-654 to REQ-674 read by intent; archive filenames scanned for journey, QA and acceptance work).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: the brief's DELIVERABLES item 1, routed to the toolbox by the maintainer's answer.*

## Triage

**Route: B** - Medium

**Reasoning:** Same shape as REQ-676 (source-audit): one new read-only toolbox action with a guide and the same integration surface, plus the ui-review browser-detection step it must cite instead of copying. Exploration finds the citation targets and the test traps; the four result classes are already decided in the REQ, so no plan is needed.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Facts shared with the sibling toolbox REQ of this wave, then the facts for this action (orchestrator exploration, 2026-10-09, read-only, at bdb98dee):

- **Router.** `skills/do-work-toolbox/SKILL.md:4` is the `argument-hint` list (one line, so every new action edits the same line). The route table is `:15-36`: `slop-check` at `:26` (source-audit's neighbour), `ui-review` at `:22` (journey-qa's neighbour). `:38` already says "Per-command help reads the selected action without executing it", so `do-work-toolbox <action> help` is served by the router, not by code in the action.
- **Per-command help format.** `skills/do-work/actions/help.md:62`: read the action's Input and When to Use sections and return at most 15 lines (purpose, usage, accepted arguments, two examples), never executing it. The new action earns correct help by having clear `## When to Use` and `## Input` sections with the exact usage line and examples. No `help` branch inside the action is needed.
- **Help menus.** Toolbox menu: `skills/do-work-toolbox/actions/help.md:5-27` (fenced, names padded to column 33; `ui-review` `:12`, `slop-check` `:16`). Core lists toolbox commands too: `skills/do-work/actions/help.md:36-40` (wrapped `·` rows inside the menu fence), so both menus get the name.
- **README.** `README.md:116` is the "Common extension calls also include ..." line next to the other toolbox calls; `:112` and `:114` show the one-paragraph shape for an action with its own explanation. One short usage line or sentence goes there.
- **Staged-skills contract.** `_dev/tests/staged-skills-contract.sh:147-167` is `toolbox_actions`; `:859-861` requires `skills/do-work-toolbox/actions/<name>.md` for each entry; `:342-344` turns each entry into a declared owner/action pair for the retired-trigger fixture (a new action needs no fixture row; the counts at `:418-440` must not change). `toolbox_files` (`:169-179`) is a must-exist list that names only some guides; adding the new guide there is optional. The same script scans every live file under `skills/` for `do-work <retired trigger>` (`:525-600`; the triggers are column 4 of `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv`, for example `scan`, `present`, `deliver`, `suggest`, `learn`, `prompt`, `note`, `inspect`): always write `do-work-toolbox <action>`, never `do-work <word>` with a retired word. The script is heavy-only (`DO_WORK_MAINTAINER_TIER=heavy`, 32 s here; heavy lane `staged-skills`).
- **Crew-member copies.** The toolbox ships its own `skills/do-work-toolbox/crew-members/prompt-injection.md` and `anti-slop.md` (they differ from the core copies). Cite them same-package as `crew-members/prompt-injection.md`; precedents `deep-explore.md:129`, `validate-feedback.md:38`, `slop-check.md:5`. Anti-slop principle 2 ("Check every citation") is `crew-members/anti-slop.md:17` in both copies.
- **Action template and traps.** `_dev/primes/lessons-action-files.md:76-143`: required description blockquote and numbered Steps; Rules, Common Rationalizations, Red Flags and Verification Checklist only when earned. The prime (`lessons-action-files.md:141`) says a Common Rationalizations table needs a do-work-specific noun; `_dev/tests/contracts/core-checks.sh:806-867` fails a row that copies another file's row (ratio above 0.75). Simplest safe choice: no such table unless a row names a real failure. Blockquote shape with the guide link: `slop-check.md:3`, `ui-review.md:3`. The blockquote also states why the action is toolbox-owned (`_dev/primes/prime-action-files.md` § Every New Action Must Justify Its Package).
- **Guide.** `skills/do-work-toolbox/docs/<action>-guide.md`; precedents `slop-check-guide.md` (82 lines), `ui-review-guide.md` (65), `stray-check-guide.md` (61). A guide is linked only from its action's blockquote; nothing else indexes it.
- **Citations.** Cross-package from `actions/` or `docs/` is `../../do-work/...` (`validate-feedback.md:25`). `_dev/tests/shipped-package-reference-contract.sh` checks paths, `path` → **Name** bold names (whole words of a heading or bold label in the target) and anchors in both the source and installed layouts; 1 s, green at bdb98dee.
- **Packaging.** `_dev/tests/suite-manifest-contract.sh:33` maps the whole `skills/do-work-toolbox` directory to `.claude/skills/do-work-toolbox`, so new files ship with no change to `suite/modules.tsv`, the installer or the updater. Installer: `printf 'y\n' | bash tools/install-do-work-suite.sh --project-root <dir> --archive <tar.gz>` (`_dev/tests/install-suite-behavior.sh:100`; the archive's top directory is `skill-do-work-main/`, `:517`). Updater: `skills/do-work/tools/do-work-update.sh` reached through `DO_WORK_UPSTREAM_URL=<local server>/archive/refs/heads/main.tar.gz` (`_dev/tests/update-script-behavior.sh:293`, `:384`).
- **Merge seam, not an ordering edge.** REQ-676 (source-audit) and REQ-677 (journey-qa), and later REQ-678 (release-check), all edit `SKILL.md:4` and the route table, both help menus, `README.md:116` and `toolbox_actions`. The `argument-hint` line conflicts for certain; the integrator keeps both names.
- **No router size budget found.** The prime says router budgets are enforced by `_dev/tests/contract-regressions.sh`; that script has no `SKILL.md` size check today (grep), so a route row costs nothing there.
- **Required-lessons consult at claim.** `do-work/lessons-index.md` matches no new satellite: `lessons-releases.md` (666 tokens) stays; `lessons-action-files.md` stays dropped for budget as captured (7189 tokens, `slugged: partial`).
- **Browser detection to cite, not copy.** `skills/do-work-toolbox/actions/ui-review.md` → **Step 2: Load Design Context**, item 4 (`:60-67`: Playwright CLI first, then the Bowser skill, else note it and suggest `do-work-toolbox install bowser`); rendered checks are **Step 8.5: Visual Verification (if browser tools available)** (`:150`). Same-package spelling: `actions/ui-review.md` → **Step 2: Load Design Context**.
- **Capture suggestion precedent.** `validate-feedback.md:112-116` ends with a user-run capture handoff block and `:122` states the boundary (read-only; the handoff is a suggestion the user runs). journey-qa may end the same way after naming the smallest justified repair.
- **Why core review does not cover this.** `skills/do-work/actions/review-work.md:160-195` (Step 7) is a smoke test by design ("this is a quick smoke test, not QA"). Do not cite REQ-675's new stage names: REQ-678 (release-check) owns that citation, and REQ-675 merges in parallel.
- **Testing guidance in a consumer.** "The project's testing guidance" means what the consumer has (its prime files, its test docs, `crew-members/testing.md` in the toolbox copy); read it as data. Evidence files go to a temporary or gitignored location the report names.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/journey-qa.md` (new) — the action
- `skills/do-work-toolbox/docs/journey-qa-guide.md` (new) — short user guide
- `skills/do-work-toolbox/SKILL.md` (modify) — `argument-hint` entry and one route row
- `skills/do-work-toolbox/actions/help.md` (modify) — one menu line
- `skills/do-work/actions/help.md` (modify) — add the name to the toolbox list
- `README.md` (modify) — one usage line next to the other toolbox calls
- `_dev/tests/staged-skills-contract.sh` (modify) — `journey-qa` in `toolbox_actions`

**Files I will NOT touch:** `suite/modules.tsv`, `tools/install-do-work-suite.sh`, `skills/do-work/tools/do-work-update.sh`, every `crew-members/*.md`, `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, anything under `do-work/` (the release and the REQ record are the integrator's).

**Acceptance criteria (restated from the REQ):**
- [ ] `skills/do-work-toolbox/actions/journey-qa.md` exists and stays short: a description blockquote that links the guide and says why it is toolbox-owned (an optional diagnostic review beside `ui-review`, no queue machinery), When to Use, Input (`do-work-toolbox journey-qa <target> [--brief <path>]`), numbered Steps, Output Format.
- [ ] It reads the requirements, the project's testing guidance and the `--brief` file when given, as data.
- [ ] It reuses existing checks and reproduces the reported sequence first, before broadening coverage.
- [ ] It exercises the combined transitions the brief or requirements name (examples marked illustrative, not exhaustive); an isolated pass never stands in for a combined one.
- [ ] It waits on observable state, not fixed delays, and keeps screenshots and other expensive instrumentation outside any performance measurement.
- [ ] Per journey it records revision, environment (browser, viewport, emulation or physical device), steps, expected and actual behavior, and evidence paths.
- [ ] Each result is passed, product defect, test defect, or unresolved; an environment limitation is unresolved with its reason; emulation is labeled emulation and never reported as physical-device verification; no browser tool means unresolved, never passed.
- [ ] It reports the smallest justified repair and the remaining checks, changes no project source, and may end with a user-run capture suggestion like `validate-feedback`'s.
- [ ] Browser tool detection cites `actions/ui-review.md` → **Step 2: Load Design Context** instead of copying it.
- [ ] It loads `crew-members/prompt-injection.md` before reading a brief or fetched page and `crew-members/anti-slop.md` before writing the report.
- [ ] Integration: route row and `argument-hint` in `skills/do-work-toolbox/SKILL.md`, one line in each help menu (toolbox and core), one usage line in `README.md`, `journey-qa` in `toolbox_actions`, and the guide `skills/do-work-toolbox/docs/journey-qa-guide.md`.
- [ ] `do-work-toolbox journey-qa help` prints usage and runs no journey; `suite/modules.tsv`, the installer and the updater are unchanged, and the packaging check passes.

## Pre-Flight

**Git:** ✓ Clean at bdb98dee apart from this run's own pre-dispatch edits: the three claimed working REQs of this wave (REQ-677 and siblings REQ-675, REQ-676), `do-work/working/baseline.json` (rewritten by this pre-flight) and the untracked run directory `do-work/runs/work-2026-10-09-225703/`. All are committed together as `[UR-151] run artifacts` before dispatch.
**Tests baseline:** ✓ Repository gate `bash _dev/tests/maintainer-verify.sh` exit 0 at bdb98dee (122 s wall, both fast Go stages executed, no prior evidence), recorded through `advance` as the green gate; focused baseline `bash _dev/tests/shipped-package-reference-contract.sh` green (probe `do-work/runs/work-2026-10-09-225703/REQ-677-preflight-probe.sh`). `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh` exit 0 (32 s, heavy lane) and `bash _dev/tests/contract-regressions.sh` exit 0 (19 s) at the same revision.
**Dependencies:** ✓ None beyond the repository's own scripts (Markdown and one shell list edit; no Go change)

*Checked by work action*
