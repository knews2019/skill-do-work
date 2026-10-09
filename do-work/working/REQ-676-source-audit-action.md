---
id: REQ-676
title: 'do-work-toolbox source-audit checks each cited source and returns a claim-to-source table without rewriting the report'
status: claimed
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  basis:
  - Route B
  - 7-file write set
  - 3 subsystems involved
  - 11 acceptance criteria
  calculated_at: 2026-10-09T23:02:37Z
created_at: 2026-10-09T22:53:15Z
user_request: UR-151
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: [_dev/primes/lessons-releases.md]
write_set: ["skills/do-work-toolbox/actions/source-audit.md", "skills/do-work-toolbox/docs/source-audit-guide.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/actions/help.md", "skills/do-work/actions/help.md", "README.md", "_dev/tests/staged-skills-contract.sh"]
related: [REQ-675, REQ-677, REQ-678]
batch: portable-verification-actions
claimed_at: 2026-10-09T22:56:58Z
dispatch_at: 2026-10-09T23:08:40Z
---
# do-work-toolbox source-audit
## What
Add a read-only toolbox action, `do-work-toolbox source-audit <report-or-url-list>` (`skills/do-work-toolbox/actions/source-audit.md`), that extracts the claims and citations from a research report or a URL list, reads each source, and returns one claim-to-source table with a judgment per claim: supported, contradicted, insufficient, or unavailable. It never rewrites the report.
## Why
Nothing in the suite checks that a cited page says what a report claims. The only related rule is anti-slop principle 2, "check every citation" (`skills/do-work/crew-members/anti-slop.md:17`). Fetching a URL is not the same as the page supporting the claim: error pages that return success, redirects to an unrelated home page, and index pages pass a naive check.
## Finding Provenance
From the validate-feedback triage of the maintainer's brief in this session (finding F3, verdict Accept):
- **Verbatim claim:** "3. source-audit <report-or-url-list> — Verify sources supporting research claims." with its seven bullets (see Detailed Requirements). Source: the brief's DELIVERABLES section.
- **Evidence:** no action in `skills/` verifies citations; `skills/do-work-toolbox/actions/deep-explore.md:129` is the precedent for loading prompt-injection before fetching.
- **Surface-cost:** N/A (a feature, not a defensive layer).
## Detailed Requirements
1. Extract the claims and their citations, then read the actual source content.
2. Per source, record the requested URL, the final URL after redirects, the retrieval outcome, and the passage that supports or contradicts the claim. When the tools cannot establish a value, write "unavailable" for it.
3. Detect error pages (including ones served with a success status), redirects to unrelated pages, index or listing pages, unavailable sources, and claims the source does not support.
4. Keep four things apart: retrieval success, claim support, corroboration by an independent source, and authority of the source.
5. A publication date that the page does not state stays "unknown". Never infer it from the fetch date or the URL.
6. Keep the original evidence and the failure history (every attempt and its outcome). Replacement candidates go in a separate section, each verified the same way, and the report file itself is never edited.
7. Output: the claim-to-source table with one judgment per claim (supported, contradicted, insufficient, unavailable), then replacement candidates, then a short list of claims that need the author's attention.
8. Toolbox ownership: an optional review/reporting action beside `slop-check`; it needs no queue machinery.

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
- Read-only: the audited report keeps its sha256; the action writes no repo file.
- No fetch tool available means every row is "unavailable" with that reason, never a guess.
## Builder Guidance
Medium certainty on the table columns; the brief lists the facts each row must carry. Latitude on column order, on file split (one action file is expected), and on the guide's length.
## Red-Green Proof
**RED prompt/case:** A scratch report outside the repo with four claims: A cites a page that returns success but whose body is "Page not found"; B cites a URL that redirects to an unrelated home page; C cites a site index page that does not contain the claim; D cites a page that supports the claim and states no publication date. Serve them locally (for example a small local HTTP server with one redirect) so the exercise needs no network.
**Why RED now:** `do-work-toolbox source-audit` does not exist; the toolbox router has no row for it.
**GREEN when:** The one-off exercise routes to the new action and prints a table where A and B are not supported (unavailable or insufficient, with the error page and the redirect named), C is insufficient, and D is supported with publication date "unknown". The report file's sha256 is unchanged. A replacement candidate, if offered for A, appears only in the separate section. `do-work-toolbox source-audit help` prints usage and audits nothing. The packaging check passes.
**Validation:** Inferred during capture. The maintainer chose one-off exercises recorded in `## Testing` over kept fixtures.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`). Matching reason: adds an action and a routing row.
## Full Context
See `do-work/user-requests/UR-151/input.md` for complete verbatim input. No queued or archived REQ shares this intent (queue REQ-654 to REQ-674 read by intent; archive filenames scanned for source, citation and audit work).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: the brief's DELIVERABLES item 3, accepted in the validate-feedback triage.*

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear (one new read-only toolbox action plus its guide), but the integration surface had to be found: the router row and argument-hint, two help menus, the README line, the staged-skills contract list, the per-command help mechanism, the crew-member copies the toolbox carries, and the test traps a new action file meets. Exploration records those facts; no plan is needed.

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
- **No citation checker exists.** Nothing under `skills/` verifies that a cited page supports a claim; the only rule is anti-slop principle 2 (`crew-members/anti-slop.md:17`). The action is new ground, so its spine is the REQ's own requirements 1-7, not an existing action.
- **Neighbour and boundary.** `slop-check` (`skills/do-work-toolbox/actions/slop-check.md`) checks a draft's prose against anti-slop and is read-only with an optional rewrite; source-audit checks the sources behind claims and never rewrites. Say so in When to Use (Do NOT use when) with a redirect each way.
- **Fetch tools vary by harness.** The action names the fallback, not a tool: when no fetch tool can read a source, every affected row is `unavailable` with that reason (REQ Constraints). Load `crew-members/prompt-injection.md` before reading the report and every fetched page (`deep-explore.md:129` is the precedent wording).

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/source-audit.md` (new) — the action
- `skills/do-work-toolbox/docs/source-audit-guide.md` (new) — short user guide
- `skills/do-work-toolbox/SKILL.md` (modify) — `argument-hint` entry and one route row
- `skills/do-work-toolbox/actions/help.md` (modify) — one menu line
- `skills/do-work/actions/help.md` (modify) — add the name to the toolbox list
- `README.md` (modify) — one usage line next to the other toolbox calls
- `_dev/tests/staged-skills-contract.sh` (modify) — `source-audit` in `toolbox_actions`

**Files I will NOT touch:** `suite/modules.tsv`, `tools/install-do-work-suite.sh`, `skills/do-work/tools/do-work-update.sh`, every `crew-members/*.md`, `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, anything under `do-work/` (the release and the REQ record are the integrator's).

**Acceptance criteria (restated from the REQ):**
- [ ] `skills/do-work-toolbox/actions/source-audit.md` exists: a description blockquote that links the guide and says why the action is toolbox-owned (an optional review/reporting action beside `slop-check`, no queue machinery), When to Use, Input (`do-work-toolbox source-audit <report-or-url-list>`), numbered Steps, Output Format.
- [ ] It extracts the claims and their citations, then reads each source's actual content.
- [ ] Per source it records the requested URL, the final URL after redirects, the retrieval outcome, and the supporting or contradicting passage; a value the tools cannot establish is written "unavailable".
- [ ] It detects error pages (including ones served with a success status), redirects to unrelated pages, index or listing pages, unavailable sources, and claims the source does not support.
- [ ] It keeps retrieval success, claim support, independent corroboration and source authority as four separate facts.
- [ ] A publication date the page does not state stays "unknown"; it is never inferred from the fetch date or the URL.
- [ ] It keeps the original evidence and every attempt with its outcome; replacement candidates go in a separate section, each verified the same way; the report file is never edited (its sha256 is unchanged) and the action writes no repo file.
- [ ] Output: the claim-to-source table with one judgment per claim (supported, contradicted, insufficient, unavailable), then replacement candidates, then the claims that need the author's attention.
- [ ] It loads `crew-members/prompt-injection.md` before reading the report or a fetched page and `crew-members/anti-slop.md` before writing the report; with no fetch tool every row is "unavailable" with that reason.
- [ ] Integration: route row and `argument-hint` in `skills/do-work-toolbox/SKILL.md`, one line in each help menu (toolbox and core), one usage line in `README.md`, `source-audit` in `toolbox_actions`, and the guide `skills/do-work-toolbox/docs/source-audit-guide.md`; every example list is marked "illustrative, not exhaustive".
- [ ] `do-work-toolbox source-audit help` prints usage and audits nothing; `suite/modules.tsv`, the installer and the updater are unchanged, and the packaging check passes.

## Pre-Flight

**Git:** ✓ Clean at bdb98dee apart from this run's own pre-dispatch edits: the three claimed working REQs of this wave (REQ-676 and siblings REQ-675, REQ-677), `do-work/working/baseline.json` (rewritten by this pre-flight) and the untracked run directory `do-work/runs/work-2026-10-09-225703/`. All are committed together as `[UR-151] run artifacts` before dispatch.
**Tests baseline:** ✓ Repository gate `bash _dev/tests/maintainer-verify.sh` exit 0 at bdb98dee (122 s wall, both fast Go stages executed, no prior evidence), recorded through `advance` as the green gate; focused baseline `bash _dev/tests/shipped-package-reference-contract.sh` green (probe `do-work/runs/work-2026-10-09-225703/REQ-676-preflight-probe.sh`). `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh` exit 0 (32 s, heavy lane) and `bash _dev/tests/contract-regressions.sh` exit 0 (19 s) at the same revision.
**Dependencies:** ✓ None beyond the repository's own scripts (Markdown and one shell list edit; no Go change)

*Checked by work action*
