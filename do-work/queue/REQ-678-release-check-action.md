---
id: REQ-678
title: 'do-work-toolbox release-check reports each delivery stage with evidence and a readiness verdict that stale or empty content cannot pass'
status: pending
created_at: 2026-10-09T22:53:15Z
user_request: UR-151
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: [_dev/primes/lessons-releases.md]
depends_on: [REQ-675]
related: [REQ-675, REQ-676, REQ-677]
batch: portable-verification-actions
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
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`). Matching reason: adds an action and a routing row, and cites a core rule (family `alternate-writer-contract-drift`).
## Full Context
See `do-work/user-requests/UR-151/input.md` for complete verbatim input. No queued or archived REQ shares this intent (queue REQ-654 to REQ-674 read by intent; archive filenames scanned for release, deploy and readiness checks).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: the brief's DELIVERABLES item 2, routed to the toolbox by the maintainer's answer.*
