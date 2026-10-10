---
id: REQ-678
title: 'do-work-toolbox release-check reports each delivery stage with evidence and a readiness verdict that stale or empty content cannot pass'
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
  calculated_at: 2026-10-10T00:17:36Z
created_at: 2026-10-09T22:53:15Z
user_request: UR-151
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: [_dev/primes/lessons-releases.md]
write_set: ["skills/do-work-toolbox/actions/release-check.md", "skills/do-work-toolbox/docs/release-check-guide.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/actions/help.md", "skills/do-work/actions/help.md", "README.md", "_dev/tests/staged-skills-contract.sh"]
depends_on: [REQ-675]
related: [REQ-675, REQ-676, REQ-677]
batch: portable-verification-actions
claimed_at: 2026-10-10T00:15:34Z
dispatch_at: 2026-10-10T00:23:20Z
builder_handback_at: 2026-10-10T00:30:56Z
---
# do-work-toolbox release-check
## What
Add a read-only toolbox action, `do-work-toolbox release-check <target> [--brief <path>]` (`skills/do-work-toolbox/actions/release-check.md`). It traces intended content from source to generated package to serving environment to consumer-visible behavior, reports each delivery stage with evidence, and returns a readiness verdict with concrete gaps. It uses the four stage names REQ-675 (core review delivery stages) defines in `skills/do-work/actions/review-work.md`, by citation.
## Why
File presence, a successful decode or a matching identifier is not proof that the consumer sees the right thing. A hit-mask image can exist, decode and be empty, so nothing is clickable. A served copy can be an older build than the source. No action in the suite checks the path from source to consumer.
## Finding Provenance
From the validate-feedback triage of the maintainer's brief in this session (finding F2, verdict Discuss; the maintainer chose "Core owns stages, toolbox cites"):
- **Verbatim claim:** "2. release-check <target> [--brief <path>] — Verify that intended content or functionality reached its consumer correctly." with its six bullets (see Detailed Requirements). Source: the brief's DELIVERABLES section.
- **Evidence:** not present in `skills/`. The 0.305.85 rule keeps routine operator acts (deploy, publish, verify on live hosts) out of the queue; a read-only check the operator invokes does not conflict with it.
- **Surface-cost:** N/A (feature).
## Detailed Requirements
1. Trace the intended source, the generated package, the serving environment, and the consumer-visible behavior.
2. Test what the consumer does with the content. File presence, successful decoding or matching identifiers alone never count as consumer evidence.
3. When relevant, investigate package-version mixing, stale caches, missing dependencies, and interrupted activation (examples, not a checklist).
4. Report implementation, integration, deployment and live acceptance separately, with evidence for each. Cite the stage definitions in `../../do-work/actions/review-work.md` instead of restating them. Mark a stage not assessed as "unassessed".
5. Label each piece of evidence as current run or historical (with its date or revision). Historical evidence never verifies a stage in the current run.
6. Return a readiness verdict (ready, not ready, or unknown) and the concrete gaps. "Ready" needs current-run live acceptance evidence. Deployment, rollback, content repair and recurring monitoring are separate operator actions; the report names them as next steps and does none of them.
7. Toolbox ownership: an optional diagnostic review; it needs no queue machinery and mints no REQ for operator work.

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
- Read-only: no deploy, rollback, cache purge, content repair or monitor setup.
- Depends on REQ-675 because it cites the stage names that REQ adds to core review.
## Builder Guidance
Firm on the per-stage table and on the readiness rule. Latitude on layout and on how the action probes a serving environment (whatever the brief and the available tools allow).
## Red-Green Proof
**RED prompt/case:** Two scratch fixtures outside the repo, served locally. A: the package's clickable-region image exists and decodes, but every pixel is transparent, so a click hits nothing. B: the served bundle carries an older version string than the source revision, and an old note says "deployed and verified" last week.
**Why RED now:** `do-work-toolbox release-check` does not exist; the toolbox router has no row for it.
**GREEN when:** The one-off exercise routes to the new action. A is not ready: live acceptance fails because the click reaches nothing, even though the file exists and decodes. B is not ready: deployment fails on the version mismatch, and last week's note is labeled historical and does not verify the stage. A third fixture where every stage has current-run evidence is ready. Unexercised stages read "unassessed". `do-work-toolbox release-check help` prints usage and checks nothing. The packaging check passes.
**Validation:** Inferred during capture. The maintainer chose one-off exercises recorded in `## Testing` over kept fixtures.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7624 tokens at the claim-time consult, over the 2000 budget; `slugged: partial`, so no targeted entry is legal). Matching reason: adds an action and a routing row, and cites a core rule (family `alternate-writer-contract-drift`).
## Full Context
See `do-work/user-requests/UR-151/input.md` for complete verbatim input. No queued or archived REQ shares this intent (queue REQ-654 to REQ-674 read by intent; archive filenames scanned for release, deploy and readiness checks).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: the brief's DELIVERABLES item 2, routed to the toolbox by the maintainer's answer.*

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear (one new read-only toolbox action plus its guide, the third in the pattern REQ-676 source-audit and REQ-677 journey-qa set), but the integration surface and the citation of the core delivery stages had to be located: the router row and argument-hint, two help menus, the README paragraph, the staged-skills contract list, the per-command help rule, and the test traps this wave already hit. Exploration records those facts; no plan is needed.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Orchestrator exploration, 2026-10-10, read-only, at 458afc0f. REQ-676 (source-audit, 0.305.91) and REQ-677 (journey-qa, 0.305.92) already shipped the same integration pattern, so most facts below point at their shipped lines.

- **Stage definitions to cite, not restate.** `skills/do-work/actions/review-work.md:160` is the heading `### Step 7: Acceptance Testing`. The four stages are `:188-192`: ":188 **Delivery stages.** A change can be checked at four stages. Which ones apply is a judgment, and many REQs have no deployment stage at all." then ":189 **Implementation**: the diff does what the REQ asks." ":190 **Integration**: the change works in the merged tree with the rest of the system, for example the test suite and adjacent flows." ":191 **Deployment**: the built or packaged result reached its serving environment or consumer install." ":192 **Live acceptance**: the consumer-visible behavior is right in the real environment." `:202` says a result scores only the stages it exercised. Step 8 `:210` is the "unassessed" wording ("each applicable Step 7 stage this review did not exercise, by stage name, marked unassessed, with the check that would cover it"). release-check uses the four names and links the definitions; it does not copy the four definition sentences.
- **Literal citation from the new action.** From `skills/do-work-toolbox/actions/` the link is `../../do-work/actions/review-work.md`; precedents `skills/do-work-toolbox/actions/journey-qa.md:18` and `validate-feedback.md:26`. The checked section form is `` `../../do-work/actions/review-work.md` → **Step 7: Acceptance Testing** `` (a real heading at `:160`); `_dev/tests/shipped-package-reference-contract.sh` checks the path, the bold name after the arrow (whole words of a heading or bold label in the target) and anchors in both the source and the installed layout. From `skills/do-work-toolbox/docs/` the same link is also `../../do-work/actions/review-work.md`.
- **Router.** `skills/do-work-toolbox/SKILL.md:4` is the one-line `argument-hint` (contains `journey-qa` and `source-audit` now). Route table `:15-38`: `journey-qa` row `:23`, `slop-check` `:27`, `source-audit` `:28`. Add one row; `:28` (after source-audit) or `:23` (after journey-qa) are the natural neighbours. `:40` says "Per-command help reads the selected action without executing it", so `do-work-toolbox release-check help` is served by the router; no help branch inside the action.
- **Per-command help format.** `skills/do-work/actions/help.md:62`: read the action's Input and When to Use sections, return at most 15 lines (purpose, usage, accepted arguments, two examples), never execute. Clear `## When to Use` and `## Input` with the exact usage line and two examples are what make help correct.
- **Help menus.** Toolbox menu `skills/do-work-toolbox/actions/help.md:5-29` (fenced, descriptions start at column 33): `journey-qa <target> [--brief]` `:13`, `source-audit <report|urls>` `:18`. `release-check <target> [--brief]` is 32 characters and would push its description past column 33; a shorter left column such as `release-check <target>` keeps the alignment (builder's choice, alignment must hold). Core menu toolbox list `skills/do-work/actions/help.md:36-40` (wrapped `·` rows, 76 to 82 characters wide); `:38` ends `· source-audit · journey-qa`. Appending ` · release-check` to `:38` is the smallest edit (94 characters; the same menu already has rows up to 113).
- **README.** `README.md:116` (journey-qa paragraph) and `:120` (source-audit paragraph) are one-sentence paragraphs naming the usage; `:118` is the "Common extension calls" line. One new one-sentence paragraph next to them, with `do-work-toolbox release-check <target> [--brief <path>]`.
- **Staged-skills contract.** `_dev/tests/staged-skills-contract.sh:147-169` is `toolbox_actions` (`journey-qa` `:153`, `source-audit` `:158`); `:861-863` requires `skills/do-work-toolbox/actions/<name>.md` per entry; `:344` turns each entry into an owner/action pair for the retired-trigger fixture, and no fixture row is needed. `toolbox_files` `:171-181` is optional for the guide (REQ-676 and REQ-677 did not add theirs). Heavy-only: `DO_WORK_MAINTAINER_TIER=heavy`, 32 to 44 s this run.
- **Retired-trigger trap.** The same script scans every live file under `skills/` for `do-work ` followed by a retired trigger (column 4 of `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv`, regex at `:527-528`), with a right boundary that rejects letters. Retired words relevant here, illustrative not exhaustive: `deliver`, `present`, `scan`, `inspect`, `note`, `ideas`. So "do-work deliver" fails and "do-work delivery" passes; always write `do-work-toolbox release-check`.
- **Example-path trap (this run's lesson, REQ-676).** `_dev/tests/staged-skills-contract.sh:1030-1034` reads any token starting with `actions/`, `tools/`, `hooks/`, `crew-members/`, `specs/` or `docs/` (not preceded by a path character) in any shipped `.md` as a package citation and fails when it does not exist. A brief example such as a path beginning with the docs directory name fails; pick example brief paths whose first segment is not a package directory (REQ-677 used one starting with `qa/`). Lesson bullet `_dev/primes/lessons-action-files.md:76`.
- **Read-only tool side effects (this run's lesson, REQ-677).** `_dev/primes/lessons-action-files.md:77`: an action that drives an external tool says where the tool runs and what cleanup may touch. journey-qa's shape: `journey-qa.md:46-48` (record `git status --porcelain --untracked-files=all`, evidence in `mktemp -d`) and `:77` (remove only new untracked tool output). release-check may fetch, serve or open pages, so the same applies.
- **Template and neighbours.** Template `_dev/primes/lessons-action-files.md:79-147` (required description blockquote and numbered Steps; Rules, Common Rationalizations, Red Flags, Verification Checklist only when earned; `_dev/tests/contracts/core-checks.sh` fails a Common Rationalizations row copied from another file). Blockquote shape with guide link and ownership reason: `journey-qa.md:3`, `source-audit.md:3`; the ownership justification rule is `_dev/primes/prime-action-files.md:9` (§ Every New Action Must Justify Its Package). Sizes: `journey-qa.md` 111 lines, `source-audit.md` 109; guides `journey-qa-guide.md` 38, `source-audit-guide.md` 50, each with a "Not to be confused with" note. Neighbours to name in Do NOT use when: `journey-qa` (checks a journey works in a browser, no stage or readiness verdict), `source-audit` (checks cited sources), `../../do-work/actions/review-work.md` (post-build REQ review; its Acceptance scores only the stages it exercised).
- **Crew-member copies.** The toolbox ships its own `skills/do-work-toolbox/crew-members/prompt-injection.md` and `anti-slop.md`; cite them same-package as `crew-members/prompt-injection.md` (precedent `journey-qa.md:38`, `:79`).
- **Operator boundary.** The REQ forbids deploy, rollback, cache purge, content repair and monitor setup. journey-qa ends with a suggested capture line the user runs (`journey-qa.md:106-111`); release-check names operator next steps and mints no REQ for operator work, so a capture line, if any, is only for a product or content defect, never for a deploy.
- **Packaging.** `_dev/tests/suite-manifest-contract.sh:33` maps the whole `skills/do-work-toolbox` directory to `.claude/skills/do-work-toolbox`, so new files ship with no change to `suite/modules.tsv`, the installer or the updater. Installer argv shape `_dev/tests/install-suite-behavior.sh:100` (archive top directory `skill-do-work-main/`, `:517`); updater through `DO_WORK_UPSTREAM_URL` to a local archive server (`_dev/tests/update-script-behavior.sh:289-296`). In this run the builder runs the install leg at its branch tip; the update leg skips with UPDATE-ALREADY-CURRENT until the version bump, so the integrator runs it after the release commit (REQ-676 and REQ-677 integration reports).
- **Fixture tooling on this machine.** `python3` with Pillow 11.3.0 (can write and read a fully transparent PNG and count opaque pixels) and `curl`; no ImageMagick. A local `python3 -m http.server` in a `mktemp -d` directory serves the three fixtures with no network.
- **Required-lessons consult (claim time).** `do-work/lessons-index.md` matches `_dev/primes/lessons-releases.md` (666 tokens, kept) and `_dev/primes/lessons-action-files.md` (7624 tokens, `slugged: partial`, so bare or nothing; dropped for budget, recorded above). The builder still reads that file's Template section and the three bullets at `:69`, `:76`, `:77` (families alternate-writer-contract-drift, example-path-read-as-citation, read-only-action-tool-side-effects) as context, not as a required set.
- **Lesson-link sweep owed at this REQ's finalization.** REQ-678 closes UR-151, so finalization moves REQ-675, REQ-676, REQ-677 and REQ-678 under `do-work/archive/UR-151/`. The three lesson links at `_dev/primes/lessons-action-files.md:69`, `:76`, `:77` point at flat archive paths and must be re-pointed in the same finalization (`_dev/primes/lessons-releases.md:7`, family canonical-link-outlives-its-target). Integrator work, not builder work.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/release-check.md` (new): the action
- `skills/do-work-toolbox/docs/release-check-guide.md` (new): short user guide
- `skills/do-work-toolbox/SKILL.md` (modify): argument-hint entry and one route row
- `skills/do-work-toolbox/actions/help.md` (modify): one menu line
- `skills/do-work/actions/help.md` (modify): the name in the toolbox list
- `README.md` (modify): one usage paragraph next to the other toolbox calls
- `_dev/tests/staged-skills-contract.sh` (modify): the name in the toolbox action list

**Files I will NOT touch:** `suite/modules.tsv`, `tools/install-do-work-suite.sh`, `skills/do-work/tools/do-work-update.sh`, `skills/do-work/actions/review-work.md` (cited, never edited), every `crew-members/*.md`, `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `_dev/primes/lessons-action-files.md`, anything under `do-work/` (the release, the lesson-link sweep and the REQ record are the integrator's).

**Acceptance criteria (restated from the REQ):**
- [ ] `skills/do-work-toolbox/actions/release-check.md` exists with a description blockquote that links `docs/release-check-guide.md` and says why it is toolbox-owned (an optional diagnostic review, no queue machinery, mints no REQ for operator work), When to Use, Input (`do-work-toolbox release-check <target> [--brief <path>]`), numbered Steps and an Output Format.
- [ ] It traces the intended source, the generated package, the serving environment and the consumer-visible behavior.
- [ ] It tests what the consumer does with the content; file presence, a successful decode or a matching identifier alone never count as consumer evidence.
- [ ] When relevant it investigates package-version mixing, stale caches, missing dependencies and interrupted activation, marked as examples, not a checklist.
- [ ] It reports implementation, integration, deployment and live acceptance separately with evidence for each, citing `../../do-work/actions/review-work.md` Step 7 for the stage definitions instead of restating them, and writes "unassessed" for a stage it did not assess.
- [ ] Every piece of evidence is labeled current run or historical (with its date or revision); historical evidence never verifies a stage in the current run.
- [ ] It returns a readiness verdict (ready, not ready, or unknown) with the concrete gaps; "ready" requires current-run live acceptance evidence.
- [ ] Deployment, rollback, content repair and recurring monitoring appear only as named operator next steps; the action does none of them and changes no project file.
- [ ] It loads `crew-members/prompt-injection.md` before reading a brief, a report or a fetched page and `crew-members/anti-slop.md` before writing the report; the brief is a plain Markdown path read as data, and game-specific checks in it are examples marked "illustrative, not exhaustive".
- [ ] Integration: route row and argument-hint in `skills/do-work-toolbox/SKILL.md`, one line in each help menu (toolbox and core), one usage paragraph in `README.md`, the name in the staged-skills toolbox action list, and the guide `skills/do-work-toolbox/docs/release-check-guide.md`.
- [ ] `do-work-toolbox release-check help` prints usage and checks nothing; `suite/modules.tsv`, the installer and the updater are unchanged.
- [ ] The one-off exercise gives: fixture A (transparent hit-mask) not ready on live acceptance, fixture B (stale served bundle plus a historical "deployed and verified" note) not ready on deployment with the note labeled historical, fixture C (all stages verified in the current run) ready; the packaging check passes.

## Pre-Flight

**Git:** ✓ Clean at 458afc0f apart from this REQ's own pre-dispatch edits: this working REQ, `do-work/working/baseline.json` (rewritten by this pre-flight) and the two untracked probes `do-work/runs/work-2026-10-09-225703/REQ-678-preflight-probe.sh` and `REQ-678-probe.sh`. All are committed together as `[UR-151] run artifacts: REQ-678 pre-dispatch` before the worktree is cut.
**Tests baseline:** ✓ Repository gate `bash _dev/tests/maintainer-verify.sh` exit 0 at 458afc0f (122 s wall, 1-minute load 2.36 at start, both fast Go stages executed on fingerprint mismatch, contract-regressions and the quiet-grep pipeline audit passed inside it), recorded through `advance` as the green gate; focused baseline probe `do-work/runs/work-2026-10-09-225703/REQ-678-preflight-probe.sh` exit 0 (the four Step 7 stage names and the Step 7 heading exist in `skills/do-work/actions/review-work.md`, and `bash _dev/tests/shipped-package-reference-contract.sh` passes). The GREEN probe `REQ-678-probe.sh` exits 1 today, as expected before the build.
**Dependencies:** ✓ REQ-675 (core review delivery stages) is archived and released in 0.305.90, so the cited stage definitions exist. No tool beyond the repository's own scripts (Markdown and one shell list edit; no Go change).

*Checked by work action*
