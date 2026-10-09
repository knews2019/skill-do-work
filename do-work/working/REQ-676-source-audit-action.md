---
id: REQ-676
title: 'do-work-toolbox source-audit checks each cited source and returns a claim-to-source table without rewriting the report'
status: claimed
created_at: 2026-10-09T22:53:15Z
user_request: UR-151
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: [_dev/primes/lessons-releases.md]
related: [REQ-675, REQ-677, REQ-678]
batch: portable-verification-actions
claimed_at: 2026-10-09T22:56:58Z
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
