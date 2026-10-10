---
id: REQ-677
title: 'do-work-toolbox journey-qa verifies whole user journeys and classifies each result as passed, product defect, test defect or unresolved'
status: completed
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
dispatch_at: 2026-10-09T23:08:40Z
builder_handback_at: 2026-10-09T23:57:36Z
integration_at: 2026-10-09T23:58:06Z
review_at: 2026-10-10T00:08:45Z
kb_status: pending
commit: ec1746ee12268fde2e59beee4208694f1c1ffd2a
heavy_verified_at: 2026-10-10T00:13:33Z
heavy_verified_revision: ec1746ee12268fde2e59beee4208694f1c1ffd2a
completed_at: 2026-10-10T00:14:03Z
release_at: 2026-10-10T00:14:03Z
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
- [x] **[PLAN]:** Write a short action on the template (blockquote, When to Use, Input, numbered Steps, Output Format; no earned sections, because the firm rules fit in two short paragraphs and a class table). Cite ui-review Step 2 item 4 for browser detection, cite same-package crew-members for prompt-injection, testing and anti-slop, cite core review-work by literal relative path. Make one-line local edits on the seven integration surfaces. Then the behavioral exercise and packaging check outside both trees. (from the builder hand-back)
- [x] **[APPLY]:** As planned, seven files exactly as the Scope list. One correction after the first heavy run: the staged-skills runtime-reference scan read the example brief path under `docs/qa/` as a toolbox `docs/` citation and failed; the example became `qa/map-zoom-brief.md` in both files. Usage is the same in the action Input, guide Usage, both help menus and README. (from the builder hand-back)
- [x] **[UNIFY]:** `git diff --stat e313e870..HEAD`: 7 files changed, 156 insertions(+), 2 deletions(-). Checks on the committed content: REQ-677-probe.sh 0, shipped-package-reference-contract.sh 0, core-checks.sh 0, staged-skills-contract.sh (heavy) 0 in 31 s, `git diff --check` 0; contract-regressions.sh 1, caused by the two run probes piping into `grep -q`, later fixed by the coordinator in c4dda51e. All seven files read in full after edit; help column alignment verified (column 33); no em-dashes in the new files. (from the builder hand-back)
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

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/journey-qa.md` (new)
- `skills/do-work-toolbox/docs/journey-qa-guide.md` (new)
- `skills/do-work-toolbox/SKILL.md` (modified)
- `skills/do-work-toolbox/actions/help.md` (modified)
- `skills/do-work/actions/help.md` (modified)
- `README.md` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)

**What was done:** Added the read-only toolbox action journey-qa and its user guide. The action reads the requirements, the project's testing guidance and an optional brief as data, reproduces the reported sequence first, then the combined transitions the brief or requirements name, waits on observable state, and records revision, environment, steps, expected, actual and evidence per journey. Each result is passed, product defect, test defect or unresolved; emulation is labeled emulation, and no browser tool or no device means unresolved. It ends with the smallest justified repair, the remaining checks and a capture line the user runs. Browser detection cites ui-review Step 2. The action is wired into the toolbox router (argument-hint and one route row), both help menus, one README paragraph and the staged-skills `toolbox_actions` list.

**Integration seams:** the merge with REQ-676 (source-audit, released as 0.305.91) conflicted in two places and both were resolved by keeping both entries: the toolbox `argument-hint` (`ui-review | journey-qa | ...` and `slop-check | source-audit | ...`) and the core help toolbox row (`present-video · slop-check · source-audit · journey-qa`, which stays within the width of the longest row). The route rows, toolbox menu lines, README paragraphs and `toolbox_actions` entries merged cleanly with both names present. `REQ-677-probe.sh`, `REQ-676-probe.sh` and `shipped-package-reference-contract.sh` exit 0 on the resolved tree before the merge commit.

## Decisions

(from the builder hand-back)

- D-01 DECIDE & STATE: no Rules, Common Rationalizations, Red Flags or Verification Checklist sections. The firm rules (four classes, isolated is not combined, emulation is not a device, no browser means unresolved) sit in the opening paragraphs, Step 4 and the Step 5 table. Keeps it short, and avoids the core-checks row-similarity ratchet.
- D-02 DECIDE & STATE: browser detection cites ui-review Step 2 item 4 and adds one clause: "a browser automation tool the session already provides also counts". Without it, an agent that has Playwright MCP but no `playwright-cli` would report a false unresolved. Value: honest results in MCP sessions. Risk: small widening beyond the cited step; revert the clause if the maintainer wants detection strictly identical to ui-review.
- D-03 DECIDE & STATE: the action tells the agent to run the browser tool from the evidence directory. The exercise showed `playwright-cli` writes `.playwright-cli/` into its current directory, which would dirty the project.
- D-04 DECIDE & STATE: route aliases `journey qa` and `user journey` beside `journey-qa`; none collides with an existing toolbox or core trigger.
- D-05 DECIDE & STATE: core help gets the name appended to the short `present-video · slop-check` row, not inserted next to `ui-review`, so no other row re-wraps.
- D-06 DECIDE & STATE: README gets a one-sentence paragraph before line 116 instead of editing line 116, so REQ-676 and REQ-678 conflicts stay one-line.
- D-07 DECIDE & STATE: the new guide is not added to `toolbox_files` (optional per Exploration; `toolbox_actions` already requires the action file).

## Discovered Tasks

(from the builder hand-back)

- **impact-user-visible** Base commit e313e870 (`[UR-151] run artifacts`) made `_dev/tests/contract-regressions.sh` exit 1, because the run probes REQ-676-probe.sh and REQ-677-probe.sh piped `sed` into `grep -q`, which `_dev/tests/quiet-grep-pipeline-audit.sh` refuses. Resolved by the coordinator in c4dda51e (the probes now read the argument hint through a here-string); the repository gate at this REQ's merge is the evidence. → report only
- **impact-negligible** `skills/do-work-toolbox/docs/ui-review-guide.md` Usage block repeats `do-work-toolbox ui-review` four times with no arguments, and `skills/do-work-toolbox/docs/stray-check-guide.md` has a `do-work-toolbox stray-check           Same thing` line; both look like leftovers of a retired-alias sweep. → report only

## Qualification

**Gate records (advance --diff-range 68bd8bd7..1cb75f9f):** `qualify` satisfied (merged range, success: no debug artifacts, no output primitives, P-A-U boxes ticked). `scope-drift` reported three SCOPE-DECLARED-NOT-TOUCHED warnings for the tokens argument-hint, journey-qa and toolbox_actions. Judged false: they are backticked words inside the Scope bullet descriptions, which the checker reads as paths (the same three-token pattern REQ-676 recorded). Every declared file was touched and none outside it.

**Scope comparison (Route B):** declared 7 files, touched 7 files, the same set (`git diff --stat 68bd8bd7..1cb75f9f`: 7 files, 156 insertions, 2 deletions). Nothing under suite/, tools/ or skills/do-work/tools/ changed, so modules.tsv, the installer and the updater are unchanged as the REQ requires. No do-work/ path in the range (queue guard printed nothing before the merge).

**Requirement trace (against the merged files):**
- R1 read requirements, testing guidance and the brief as data: journey-qa.md Step 1 (prompt-injection loaded first; testing guidance includes same-package `crew-members/testing.md`, which exists in the toolbox).
- R2 reuse existing checks, reported sequence first: Step 2 second paragraph (list and reuse existing checks), Step 3 ("Before broadening coverage").
- R3 combined transitions, isolated pass never stands in: Step 4 (isolated then combined rows reported separately; examples marked "illustrative not exhaustive") and line 7.
- R4 observable state, no fixed delays; instrumentation outside timed spans: Step 3 and Step 4 second paragraph.
- R5 per-journey revision, environment, steps, expected, actual, evidence: Step 2 (revision), Step 4 (environment) and the Output Format journey block.
- R6 four classes, environment limit unresolved with reason, emulation labeled: Step 5 table, Step 2 (no browser tool), Step 4 third paragraph (emulation, no device).
- R7 smallest justified repair and remaining checks, diagnostic only, user-run capture suggestion: Step 6, Output Format, line 5 and line 111 (validate-feedback precedent).
- R8 browser detection cited, not copied: Step 2 cites `actions/ui-review.md` → **Step 2: Load Design Context** item 4 and adds one clause (D-02) that a browser tool the session already provides also counts.
- R9 toolbox ownership beside ui-review, no queue machinery: description blockquote at line 3.
- Integration: route row and argument-hint (SKILL.md:4 and :23), toolbox menu line, core menu toolbox row, README paragraph, `toolbox_actions` entry, guide present; anti-slop loaded in Step 6 (same-package toolbox copy); per-command help served by the router from When to Use and Input.
- Constraints: read-only on project source with evidence in a temporary or gitignored directory (line 5, Step 2 third paragraph, Step 6 status comparison); no browser tool means unresolved (Step 2).

**Debug artifacts:** none found by the gate. No em-dashes in the two new files.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` run directly, unpiped, at the first merge 1cb75f9f and again at the fix merge ec1746ee, each after waiting for a 1-minute load under 5 with no other gate running; then `advance --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-09-225703/REQ-677-probe.sh` at ec1746ee (advance ran the probe itself).
**Result:** ✓ Repository gate exit 0 at 1cb75f9f (124 s wall) and exit 0 at ec1746ee (126 s wall); no red run, so no retry. Green probe exit 0 (BLOCKED-PROBE-SUCCEEDED; the probe checks both new files, the guide link, the ui-review citation, the three result classes, argument-hint, the route row, both help menus, README, `toolbox_actions`, then runs shipped-package-reference-contract.sh). The green-gate record is satisfied. The builder's red `contract-regressions.sh` (the two run probes piping `sed` into `grep -q`) was already fixed by the coordinator in c4dda51e; the gate at both merges includes that audit and passed.

**Merge seam checks (before the merge commit, resolved tree):** `REQ-677-probe.sh` 0, `REQ-676-probe.sh` 0, `bash _dev/tests/shipped-package-reference-contract.sh` 0, `git diff --check --cached` clean.

**Red-green validation:** *(one-off behavioral exercise from `## Red-Green Proof`; nothing committed)*
- Routing: ✗ before (base e313e870 has no `journey` in the toolbox router, so the word falls to "unknown single words print help") → ✓ after (route row `journey-qa` → `./actions/journey-qa.md`, which exists).
- Builder exercise (hand-back, scratch zoom-app served locally, playwright-cli 0.1.14, headless Chrome 155): zoom alone, pan alone, reset alone and pan, zoom, reset passed; zoom, pan, reset was a product defect (expected x 0, actual x 50) with revision, environment, steps, expected, actual and evidence; the phone run was labeled emulation, not physical-device verification, and the physical-phone row stayed unresolved; the wrong-selector variant was a test defect, proved by reaching reset through its accessible name; with playwright-cli hidden every rendered journey was unresolved (simulated: the session also had browser MCP tools); `journey-qa help` printed usage and ran nothing; fixture projects clean afterwards.
- Re-run by the reviewer at 1cb75f9f: the zoom-app copy in a fresh temp directory reproduced the product defect (zoom, pan, reset gave x 50) while each step alone and pan, zoom, reset passed; playwright-cli wrote `.playwright-cli/` only into the evidence directory; the fixture tree stayed clean.
- Re-run by the integrator at the merge: the route row resolves to an existing action file, and per-command help is served by the router rule (toolbox SKILL.md "Per-command help reads the selected action without executing it"; core help.md "Never execute the command while serving help") from the action's When to Use and Input sections, so it runs no journey.
- Packaging install leg: the builder (branch tip 4b2c7a5c) and the reviewer (merge 1cb75f9f) each installed into a disposable Git consumer through the canonical installer: exit 0, both new files present under `.claude/skills/do-work-toolbox/`, every relative link resolves, and a brief, a `do-work/queue/` REQ and an app file keep their sha256. The builder's update leg stopped with UPDATE-ALREADY-CURRENT because the branch has no version bump; the integrator re-runs the update leg after the release commit (coordinator ruling), with the result in the integration report.
- Not exercised by anyone: a physical device, and a real agent with no browser tool at all.

**New tests added:**
- None kept (the maintainer chose one-off exercises). `journey-qa` added to `toolbox_actions`, so staged-skills now requires the action file to exist and stage.

**Heavy verification plan:**
- Range: 68bd8bd78551e80c1541bd5713340bde65b37e14..ec1746ee12268fde2e59beee4208694f1c1ffd2a (the plan at the first merge 1cb75f9f selected the same six lanes)
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` (staged-skills-contract.sh matched subtree _dev/tests)
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` (staged-skills-contract.sh matched subtree _dev/tests)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` (staged-skills-contract.sh matched subtree _dev/tests)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` (staged-skills-contract.sh matched subtree _dev/tests; the toolbox SKILL.md, help.md, journey-qa.md, journey-qa-guide.md and core help.md matched subtree skills)
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` (staged-skills-contract.sh matched subtree _dev/tests)
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` (README.md matched exact path; staged-skills-contract.sh matched subtree _dev/tests)

*Verified by work action*

## Review

**Overall: 95%** | 2026-10-10T00:08:45Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- F1: the guide said browser detection is only Playwright CLI or Bowser, while the action also counts a browser tool the session already provides. Resolved in 84738e35 (merged ec1746ee): the guide's "What it needs" line now matches the action. — resolved
- F2: builder decision D-02 gives journey-qa a wider browser-detection rule than `skills/do-work-toolbox/actions/ui-review.md` Step 2 item 4 and Step 8.5. The clause follows the Agent Compatibility rule (generalized, no tool API), but one capability check now has two meanings in one package. Options: widen ui-review the same way (recommended by the reviewer) or drop the clause. ui-review is outside this REQ's Scope. — impact-rule-change → report only
- F3: Step 6 said to remove "anything the browser tool wrote", which could revert a tracked file or delete a user's file, and plain `git status --porcelain` collapses untracked directories. Resolved in 84738e35: Steps 2 and 6 use `--untracked-files=all`, remove only new untracked tool output, and never revert a tracked file. — resolved
- F4: `skills/do-work-toolbox/actions/ui-review.md` "Do NOT use when" never points users to journey-qa, though UR-151 asks for next-step guidance (the same gap REQ-676's review recorded for slop-check). Outside this REQ's Scope. — impact-user-visible → report only

**Acceptance:** Pass. Implementation and integration assessed: probes, heavy staged-skills and core-checks re-run green at 1cb75f9f and ec1746ee, a real zoom, pan, reset re-run by the reviewer reproduced the product defect while the isolated steps passed, routing and per-command help traced (help runs nothing). Deployment assessed only as a fresh install at 1cb75f9f at review time (the update leg runs after the release commit). Live acceptance unassessed.
**Restatement sweep:** redefined the list of toolbox actions (one more). Checked the toolbox argument-hint and route table, both help menus, README, `toolbox_actions`, `skills/do-work-toolbox/actions/tutorial.md` (a sample, not a list), `skills/do-work/next-steps.md` (derived from the router), the toolbox crew-member caller lists (illustrative), the Go toolbox command registry (deterministic CLI phases only) and the retired-trigger fixture (does not apply). No count of toolbox actions exists. All agree; the one gap is the ui-review next-step pointer (F4).
**Suggested testing:** 4 items (update leg after release, live help run in a real consumer, a real agent with no browser tool, a physical device)
**Follow-ups created:** None (2 open findings report only, F1 and F3 resolved)

*Reviewed by review-work action* (first pass 94% at 1cb75f9f, re-check of 68bd8bd7..ec1746ee 95%; full report `do-work/runs/work-2026-10-09-225703/REQ-677-review.md`)

## Lessons Learned

**What worked:** Resolving the REQ-676 seam by keeping both names in the order each file already used left only two textual conflicts (argument-hint and the core help row); the route rows, menu lines, README paragraphs and `toolbox_actions` entries merged cleanly because both builders added new lines instead of editing a shared one. Running both run probes plus the citation contract on the resolved tree before the merge commit proved the seam before any gate time was spent. Waiting for a 1-minute load under 5 gave two clean gate runs with no load reruns.
**What didn't:** The builder's first example brief path began with the docs directory name and failed the heavy staged-skills scan (the same trap REQ-676 recorded as a lesson). The first Step 6 cleanup wording ("remove anything the browser tool wrote") was broader than the read-only promise: it allowed reverting a tracked file, and plain `git status --porcelain` hides new files inside an untracked directory.
**Worth knowing:** Browser CLIs such as playwright-cli write session files into their current directory, so a read-only action that drives one must say where to run it (the evidence directory) and limit any cleanup to new untracked paths the tool wrote. journey-qa counts a browser tool the session already provides while ui-review does not (review F2); widening ui-review is open, report only.

## Orientation

Now `do-work-toolbox journey-qa <target> [--brief <path>]` reproduces a reported user journey in a browser, tries the combined transitions around it, and classifies each result as passed, product defect, test defect or unresolved; lives in the toolbox package beside ui-review (`skills/do-work-toolbox/actions/journey-qa.md`, guide in `docs/`; prime `_dev/primes/prime-action-files.md`). Leaf addition: one more toolbox action, no change to the suite's shape. Prime spot-check: the paths `_dev/primes/prime-action-files.md` names for toolbox actions still exist.

## Heavy Verification Plan

- Base: 68bd8bd78551e80c1541bd5713340bde65b37e14
- Target: ec1746ee12268fde2e59beee4208694f1c1ffd2a
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`. Reasons: _dev/tests/staged-skills-contract.sh matched subtree _dev/tests.
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`. Reasons: _dev/tests/staged-skills-contract.sh matched subtree _dev/tests.
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`. Reasons: _dev/tests/staged-skills-contract.sh matched subtree _dev/tests.
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. Reasons: _dev/tests/staged-skills-contract.sh matched subtree _dev/tests; skills/do-work-toolbox/SKILL.md, skills/do-work-toolbox/actions/help.md, skills/do-work-toolbox/actions/journey-qa.md, skills/do-work-toolbox/docs/journey-qa-guide.md and skills/do-work/actions/help.md matched subtree skills.
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`. Reasons: _dev/tests/staged-skills-contract.sh matched subtree _dev/tests.
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`. Reasons: README.md matched exact path README.md; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests.

## Heavy Verification Result

- Target revision: ec1746ee12268fde2e59beee4208694f1c1ffd2a
- Execution revision: ec1746ee12268fde2e59beee4208694f1c1ffd2a (detached checkout `.git/work-run-work-2026-10-09-225703/drain-head-REQ-677`, removed after the run)
- queue-kanban-javascript: exit 0, executed (fingerprint_mismatch), 8 s
- queue-kanban-browser: exit 0, executed (fingerprint_uncertain), 85 s, QUEUE_KANBAN_BROWSER set, no HEAVY-RUN-LANE-SKIPPED
- do-work-cli-integrations: exit 0, executed (fingerprint_mismatch), 63 s
- staged-skills: exit 0, executed (fingerprint_mismatch), 44 s
- updater: exit 0, executed (fingerprint_mismatch), 68 s
- installer: exit 0, executed (fingerprint_mismatch), 28 s
- The drain ran once, at the fix merge; no drain ran at the first merge 1cb75f9f.

## Timing

Observed 2026-10-09T23:57:36Z to 2026-10-10T00:13:24Z: 15m 48s total, 16m 10s attributed across 5 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 9m 52s | 3 |
| review | 5m 48s | 1 |
| handback-merge | 30s | 1 |

Slowest stage: review / review at first merge, 5m 48s, outcome success.
Slowest command: verification-gate / heavy drain, 5m 11s, exit 0, .
