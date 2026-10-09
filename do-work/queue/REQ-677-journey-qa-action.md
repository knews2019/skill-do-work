---
id: REQ-677
title: 'do-work-toolbox journey-qa verifies whole user journeys and classifies each result as passed, product defect, test defect or unresolved'
status: pending
created_at: 2026-10-09T22:53:15Z
user_request: UR-151
domain: testing
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: [_dev/primes/lessons-releases.md]
related: [REQ-675, REQ-676, REQ-678]
batch: portable-verification-actions
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
