---
id: REQ-676
title: 'do-work-toolbox source-audit checks each cited source and returns a claim-to-source table without rewriting the report'
status: completed
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
builder_handback_at: 2026-10-09T23:28:15Z
integration_at: 2026-10-09T23:28:31Z
review_at: 2026-10-09T23:53:30Z
kb_status: pending
commit: 6926c82b5c4e12e1e88cf745d232fda3543c82a8
heavy_verified_at: 2026-10-09T23:54:00Z
heavy_verified_revision: 6926c82b5c4e12e1e88cf745d232fda3543c82a8
completed_at: 2026-10-09T23:54:36Z
release_at: 2026-10-09T23:54:36Z
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
- [x] **[PLAN]:** One action file whose spine is REQ requirements 1-7; four facts kept apart structurally by using two tables (claim table carries judgment and corroboration; source evidence table carries retrieval, page check, passage, date, authority). Error page and unrelated redirect classify as `unavailable` (the cited content was not read); index page as `insufficient`. Judgment precedence for multi-source claims fixed in one sentence. Help served by the router from When to Use and Input, no help code. Integration edits one line or row each, placed next to `slop-check`. (from the builder hand-back)
- [x] **[APPLY]:** Seven files exactly as the Scope list. One correction during apply: example path `docs/research/market-scan.md` was read by the staged-skills runtime-reference scan as a shipped `docs/` citation and failed; changed to `research/market-scan.md` in both files. (from the builder hand-back)
- [x] **[UNIFY]:** `git diff --stat e313e870..HEAD`: 7 files changed, 166 insertions(+), 2 deletions(-). `git diff --check`: clean (exit 0). No em-dashes in the two new files. Checks on the final tree: REQ-676-probe.sh 0, shipped-package-reference-contract.sh 0, core-checks.sh 0, staged-skills-contract.sh (heavy) 0 in 32 s; contract-regressions.sh 1, base-caused by the two run probes piping into `grep -q`, later fixed by the coordinator in c4dda51e. Files reviewed: all seven, plus the installed copies of the two new files (byte-identical to source). (from the builder hand-back)
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

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/source-audit.md` (new)
- `skills/do-work-toolbox/docs/source-audit-guide.md` (new)
- `skills/do-work-toolbox/SKILL.md` (modified)
- `skills/do-work-toolbox/actions/help.md` (modified)
- `skills/do-work/actions/help.md` (modified)
- `README.md` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)

**What was done:** Added the read-only toolbox action source-audit and its user guide. The action extracts claims and citations, fetches each source, and prints a claim table (one judgment per claim: supported, contradicted, insufficient, unavailable) beside a source evidence table (requested and final URL, every retrieval attempt, page check, passage, publication date, authority), then replacement candidates in their own section and the claims that need the author's attention; it never edits the report. The action is wired into the toolbox router (argument-hint and one route row), both help menus, one README paragraph and the staged-skills `toolbox_actions` list.

**Integration seams:** none at this merge. The builder branch merged cleanly onto c33b8994 (REQ-675 had landed as 0.305.90 and touched none of these files). REQ-677 (journey-qa action) edits the same router, menu, README and contract lines and integrates after this REQ, so its integrator keeps both entries.

## Decisions

(from the builder hand-back)

- D-01 DECIDE & STATE: two tables in the claim-to-source section (one row per claim with the judgment, one row per source with the evidence). This keeps one judgment per claim and makes the four facts separate columns instead of prose. Builder Guidance gave latitude on columns.
- D-02 DECIDE & STATE: error page and unrelated redirect judge `unavailable` (the cited content was never read); an index page judges `insufficient` (the cited URL is what loaded, it just does not state the claim). Matches the REQ's GREEN wording.
- D-03 DECIDE & STATE: multi-source precedence is contradicted > supported > insufficient > unavailable, so a single contradicting source is never hidden by a supporting one.
- D-04 DECIDE & STATE: uncited claims and URL-list entries with no claim are `insufficient` with the reason named, so every row still has one of the four judgments.
- D-05 DECIDE & STATE: publication date excludes server headers and archive capture dates as well as fetch date and URL; an archive snapshot is a likely replacement candidate, and its capture date is the easy wrong answer.
- D-06 DECIDE & STATE: the audit is printed only; no saved-file option (the REQ says the action writes no repo file).
- D-07 DECIDE & STATE: README gets a new one-sentence paragraph after the "Common extension calls" line rather than editing that line, to keep the edit local. It still sits next to REQ-677's likely edit, so expect a textual conflict there.
- D-08 DECIDE & STATE: did not add the guide to `toolbox_files` (optional per Exploration; one fewer seam line).

## Discovered Tasks

(from the builder hand-back)

- **impact-user-visible** Base commit e313e870 (`[UR-151] run artifacts`) made `_dev/tests/contract-regressions.sh` exit 1, because the run probes REQ-676-probe.sh:7 and REQ-677-probe.sh:9 piped `sed` into `grep -q`, which `_dev/tests/quiet-grep-pipeline-audit.sh` refuses. Resolved by the coordinator in c4dda51e (the probes now read the argument hint through a here-string); the repository gate at this REQ's merge is the evidence. → report only
- **impact-negligible** The staged-skills runtime-reference scan treats any `docs/...` token in a shipped file as a package citation, including an example user path inside a usage example. Not a bug to fix here; it is this REQ's lesson. → report only

## Qualification

**Gate records (advance --diff-range c33b8994..c93701c3):** `qualify` satisfied (merged range, success: no debug artifacts, no output primitives, P-A-U boxes ticked). `scope-drift` reported three SCOPE-DECLARED-NOT-TOUCHED warnings for the tokens argument-hint, source-audit and toolbox_actions. Judged false: they are backticked words inside the Scope bullet descriptions, which the checker reads as paths. Every declared file was touched and none outside it.

**Scope comparison (Route B):** declared 7 files, touched 7 files, the same set (`git diff --stat c33b8994..c93701c3`: 7 files, 166 insertions, 2 deletions). Nothing under suite/, tools/ or skills/do-work/tools/ changed, so modules.tsv, the installer and the updater are unchanged as the REQ requires. No do-work/ path in the range.

**Requirement trace (against the merged files):**
- R1 extract claims and citations, read actual content: source-audit.md Step 2 (claims C1.. and sources S1.., uncited claims still get a row) and Step 3 (fetch with redirects followed).
- R2 requested URL, final URL, retrieval outcome, passage, "unavailable" when unknown: Step 3 and the source evidence table columns in Output Format.
- R3 error pages (including 200), unrelated redirects, index pages, unavailable sources, unsupported claims: Step 4 page classes plus Step 5 judgments.
- R4 four separate facts: Step 5 bullets (support, corroboration, authority) and the retrieval column; the two tables keep them in separate columns.
- R5 publication date "unknown" unless stated, never from fetch date or URL: Step 5 last bullet (also excludes server headers and archive capture dates).
- R6 original evidence and every attempt kept, candidates separate and verified the same way, report never edited: Step 3 (every attempt), Step 6, Step 1 and Step 7 SHA-256 check, line 5 read-only statement.
- R7 output order: Output Format (claim table, replacement candidates, needs the author's attention).
- R8 toolbox ownership beside slop-check, no queue machinery: description blockquote at line 3.
- Integration: route row and argument-hint (SKILL.md), toolbox menu line, core menu toolbox row, README paragraph, toolbox_actions entry, guide present; prompt-injection loaded in Step 1 and anti-slop in Step 7 (same-package toolbox copies); example lists marked "illustrative, not exhaustive"; per-command help served by the router from When to Use and Input.
- Constraint "no fetch tool means every row unavailable with that reason": Step 3 second paragraph.

**Debug artifacts:** none found by the gate. No em-dashes in the two new files.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at the merge c93701c3 (run directly, unpiped), then `advance --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-09-225703/REQ-676-probe.sh` (advance ran the probe itself).
**Result:** ✓ Repository gate exit 0 (121 s wall, both fast Go stages executed). Green probe exit 0 (BLOCKED-PROBE-SUCCEEDED; the probe checks both new files, the guide link, argument-hint, the route row, both help menus, README, `toolbox_actions`, then runs shipped-package-reference-contract.sh). The green-gate record is satisfied at the merge. The quiet-grep pipeline audit inside contract-regressions passed, which confirms the builder's discovered task is resolved by c4dda51e.

**Red-green validation:** *(one-off behavioral exercise from `## Red-Green Proof`, run by the builder and recorded in its hand-back; nothing committed)*
- Routing: ✗ before (base e313e870 has no `source-audit` in the toolbox router and no action file, so the word falls to "unknown single words print help") → ✓ after (route row `source-audit` → `./actions/source-audit.md`).
- Four-claim scratch report served by a local Python HTTP server (no network): A (error page served with 200) → `unavailable`, error page named; B (302 to an unrelated home page) → `unavailable`, redirect named; C (news index page) → `insufficient`; D (content page, no date) → `supported`, publication date `unknown`. The one replacement candidate (for C, a linked page that returned 404) stayed in its own section and did not change C's row. Report SHA-256 e5bdee3f...5e35c43 unchanged before and after.
- `do-work-toolbox source-audit help`: built from When to Use and Input, 12 lines, no fetch run.
- Packaging check (builder, at branch tip bf37f4ff): install into a disposable git consumer through `tools/install-do-work-suite.sh --archive` exit 0; the action (byte-identical to source) and the guide present under `.claude/skills/do-work-toolbox/`; their relative links and same-package citations resolve; a consumer brief, a `do-work/queue/` REQ and an app file keep their SHA-256. The update leg fetched but skipped with UPDATE-ALREADY-CURRENT because the branch has no version bump; the integrator re-runs the update leg after the release commit (coordinator ruling), with the result in the integration report.
- Not exercised by anyone: the "no fetch tool" path and an injection attempt inside a fetched page (prose rules in Steps 1 and 3). The exercise was run by the agent that wrote the action, so it shows the procedure is followable, not that an independent agent reads it the same way.
- Re-run by the integrator at the merge: the route row resolves to an existing `./actions/source-audit.md` at c93701c3, and the router's per-command help rule (SKILL.md line 39, core help.md "Never execute the command while serving help") serves help without running the action. Everything else above is taken from the hand-back.

**New tests added:**
- None kept (the maintainer chose one-off exercises). `source-audit` added to `toolbox_actions`, so staged-skills now requires the action file to exist and stage.

**Heavy verification plan:**
- Range: c33b8994843eecf858e2c89f9f0836d1d31121ef..c93701c3d39dd897cfaecdb3cf0f98385d87771d
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` (staged-skills-contract.sh matched subtree _dev/tests)
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` (staged-skills-contract.sh matched subtree _dev/tests)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` (staged-skills-contract.sh matched subtree _dev/tests)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` (staged-skills-contract.sh matched subtree _dev/tests; the toolbox SKILL.md, help.md, source-audit.md, source-audit-guide.md and core help.md matched subtree skills)
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` (staged-skills-contract.sh matched subtree _dev/tests)
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` (README.md matched exact path; staged-skills-contract.sh matched subtree _dev/tests)

**Repository gate after the review fix (merge 6926c82b, cumulative range c33b8994..6926c82b):** the fix touched only `skills/do-work-toolbox/actions/source-audit.md` and `skills/do-work-toolbox/docs/source-audit-guide.md`, so the gate ran directly at the new merge (advance was already past the test gate).
**Repository gate retry:** first run exited 1 (92 s, load about 14 right after the heavy drain): only the 30 s per-file budget, `strict_behavior_regression_test.go` 32.45 s in queue-kanban. Rerun exited 1 (143 s, load about 17 from other sessions): `TestBlockedProbeRunsFromSelectedRepositoryRoot` in do-work-cli `internal/nextselection`, probe launch "no such process" (exit 125). Neither file can read the two Markdown files the fix changed; the same test passed 3 of 3 alone (`go test -count=3 -run ...`), and the same code passed the gate at c93701c3. A third run after waiting for load under 5 exited 0 (84 s). Classified as load flakes, not this REQ.
**Heavy verification plan at the fix merge:** `plan-heavy-verification` over c33b8994..6926c82b selected the same six lanes as above.

*Verified by work action*

## Review

**Overall: 93%** | 2026-10-09T23:53:30Z

| Dimension | Score |
|-----------|-------|
| Requirements | 96% |
| Code Quality | 93% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- F1: slop-check gives no pointer to source-audit (`skills/do-work-toolbox/docs/slop-check-guide.md:34` says slop-check cannot verify cited numbers, and `skills/do-work-toolbox/actions/slop-check.md:15-19` has no redirect). The UR asked for next-step guidance; the slop-check files are outside this REQ's Scope. — impact-user-visible → report only
- F2: example lists not marked "illustrative, not exhaustive" in source-audit.md Steps 3 and 4 and the guide's judgment table. Resolved in fbe78182 (merged 6926c82b). — resolved
- F3: the candidate table lacked Final URL and Authority. Resolved in fbe78182: it now carries the same evidence columns as the source table (corroboration stays per claim). — resolved
- F4 (nit): the error-page template row wrote the passage as `none` where Step 3 gives `unavailable`. Resolved in fbe78182. — resolved
- F5 (nit): the attention list said injection attempts are "found in Step 1". Resolved in fbe78182 ("found in the report or a fetched page"). — resolved
- F6 (nit): the template header always printed a report SHA-256. Resolved in fbe78182 ("file input only"). — resolved
- F7 (nit, from the fix re-check): `skills/do-work-toolbox/docs/source-audit-guide.md:41` still lists the report SHA-256 with no file-only condition, the guide's copy of F6. — impact-negligible → report only
- F8 (nit, from the fix re-check): the guide's `insufficient` row says "Content was read" but also lists "no source cited", where nothing was read. Present before the fix. — impact-negligible → report only

**Acceptance:** Pass — implementation and integration: probe, shipped-package citation contract and core-checks re-run green at c93701c3 and the probe again at 6926c82b, routing and per-command help traced (help runs nothing), the no-fetch-tool path walked on a scratch report (every source `unavailable`, report SHA-256 unchanged); the four-claim local-server exercise and the packaging install are from the builder hand-back. Deployment (update leg after the version bump) and live acceptance (a real fetch tool) unassessed at review time.
**Restatement sweep:** redefined the list of toolbox actions (one more). Checked the toolbox router argument-hint and route table, both help menus, README, `toolbox_actions` and the retired-trigger fixture counts, `skills/do-work-toolbox/actions/tutorial.md` (a scenario subset, not a full list), `skills/do-work/next-steps.md`, the anti-slop and prompt-injection caller lists (marked illustrative) and the Go toolbox command registry (CLI commands only). No count of toolbox actions exists. All agree; the one gap is the slop-check next-step pointer (F1).
**Suggested testing:** 4 items (update leg after release, live run with a real fetch tool, injection inside a fetched page, a claim with one supporting and one contradicting source)
**Follow-ups created:** None (3 open findings report only, F2-F6 resolved)

*Reviewed by review-work action* (first pass 91% at c93701c3, re-check of c93701c3..6926c82b 93%; full report `do-work/runs/work-2026-10-09-225703/REQ-676-review.md`)

## Lessons Learned

**What worked:** Two tables (one row per claim with the judgment, one row per source with the evidence) kept retrieval, support, corroboration and authority apart without prose rules. Running the reviewer right after qualify, in parallel with the gate and the heavy drain, found the template gaps early enough that one fix commit and one re-merge closed five findings.
**What didn't:** The builder's first example path `docs/research/market-scan.md` failed the staged-skills runtime-reference scan, which reads any token starting with a shipped package directory as a package citation. The gate at the fix merge went red twice on unrelated load flakes (a per-file budget, then a process-launch race) while other sessions loaded the machine; only a wait for load under 5 gave a clean run.
**Worth knowing:** In an action's Output Format, every example row is copied literally by the agent that runs it, so the template must obey the action's own rules (the error-page passage is `unavailable`, not `none`; the SHA-256 line applies to file input only). Unchanged heavy lanes were `reused` at the fix merge because the fix touched only paths their fingerprints do not cover; only staged-skills and the browser lane (always uncertain) executed again.

## Orientation

Now `do-work-toolbox source-audit <report-or-url-list>` checks whether cited sources say what a report claims and prints a claim-to-source table without editing the report; lives in the toolbox package beside slop-check (`skills/do-work-toolbox/actions/source-audit.md`, guide in `docs/`; prime `_dev/primes/prime-action-files.md`). Leaf addition: one more toolbox action, no change to the suite's shape.

## Heavy Verification Plan

- Base: c33b8994843eecf858e2c89f9f0836d1d31121ef
- Target: 6926c82b5c4e12e1e88cf745d232fda3543c82a8
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`. Reasons: _dev/tests/staged-skills-contract.sh matched subtree _dev/tests.
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`. Reasons: _dev/tests/staged-skills-contract.sh matched subtree _dev/tests.
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`. Reasons: _dev/tests/staged-skills-contract.sh matched subtree _dev/tests.
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. Reasons: _dev/tests/staged-skills-contract.sh matched subtree _dev/tests; skills/do-work-toolbox/SKILL.md, skills/do-work-toolbox/actions/help.md, skills/do-work-toolbox/actions/source-audit.md, skills/do-work-toolbox/docs/source-audit-guide.md and skills/do-work/actions/help.md matched subtree skills.
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`. Reasons: _dev/tests/staged-skills-contract.sh matched subtree _dev/tests.
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`. Reasons: README.md matched exact path README.md; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests.

## Heavy Verification Result

- Target revision: 6926c82b5c4e12e1e88cf745d232fda3543c82a8
- Execution revision: 6926c82b5c4e12e1e88cf745d232fda3543c82a8 (detached checkout `.git/work-run-work-2026-10-09-225703/drain-head-REQ-676`, removed after the run)
- queue-kanban-javascript: exit 0, reused (fingerprint_match; executed at c93701c3, 7 s)
- queue-kanban-browser: exit 0, executed (fingerprint_uncertain), 83 s, QUEUE_KANBAN_BROWSER set, no HEAVY-RUN-LANE-SKIPPED
- do-work-cli-integrations: exit 0, reused (fingerprint_match; executed at c93701c3, 65 s)
- staged-skills: exit 0, executed (fingerprint_mismatch), 34 s
- updater: exit 0, reused (fingerprint_match; executed at c93701c3, 65 s)
- installer: exit 0, reused (fingerprint_match; executed at c93701c3, 35 s)
- An earlier drain at the first merge c93701c3 executed all six lanes, all exit 0 (javascript 7 s, browser 84 s, cli-integrations 65 s, staged-skills 40 s, updater 65 s, installer 35 s), before review findings F2-F6 moved the target.

## Timing

Observed 2026-10-09T23:28:17Z to 2026-10-09T23:53:11Z: 24m 54s total, 30m 50s attributed across 6 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 22m 49s | 4 |
| review | 7m 47s | 1 |
| handback-merge | 14s | 1 |

Slowest stage: verification-gate / repository gate at the fix merge (three runs), 13m 26s, outcome success.
