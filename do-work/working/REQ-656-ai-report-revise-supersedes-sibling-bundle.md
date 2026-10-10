---
id: REQ-656
title: 'ai-report revise writes a new sibling bundle that supersedes the prior one'
status: claimed
created_at: 2026-10-09T21:10:00Z
user_request: UR-144
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-654, REQ-655, REQ-657]
batch: ai-report-modes
depends_on: [REQ-655]
write_set: ["skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/actions/ai-report-reference.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/actions/help.md"]
claimed_at: 2026-10-10T15:51:14Z
route: B
required_lessons: ["_dev/primes/lessons-releases.md"]
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-10-10T15:55:08Z
  basis:
    - Route B
    - 5-file write set
    - 8 acceptance criteria
---
# ai-report Revise Writes a New Sibling Bundle That Supersedes the Prior One
## What
`ai-report revise <dir|latest> [what changed]` updates a report without breaking the immutability contract. It never edits the old bundle. It writes a new sibling bundle that opens with a "rev-N (date): changed / still to do" block, carries a `supersedes` link to the previous bundle, and regenerates the catalog (REQ-655) so the old bundle's `superseded_by` points at the new one. The revised report gets the readable-style baseline and a table of contents when it is long, and the action prints the final path at the end.
## Why
Report A3 (UR-144 input): 3 sessions in 1 consumer repo asked to update an existing report, then each quality step arrived as its own prompt on the same report: "update it ... with a new rev-* tag what changed and what we need to still do", "add a TOC to the ai-report so it's easy to navigate", "improve the style for the HTML report so it is more readable", "give me the link after you commit it". Nothing in the action offers them.
## Verified Facts (checked at 0.305.87)
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md:71`: "Existing artifacts are immutable: never delete, truncate, merge into, rename, migrate, or overwrite them." `:73` holds the `-2`, `-3` sibling rule.
- `skills/do-work-toolbox/actions/architecture-report.md:12` ("A prior report is never edited") and `:119` ("Never edit, delete, or regenerate a prior report"). Its Step 3 (report cites `:71`, `:83`) already reads the prior report, re-checks it, and writes an authored "changed since last report" opening. This REQ copies that pattern.
- `skills/do-work-toolbox/actions/ai-report-reference.md` Report Design Rules (report cites `:70`, `:77`, `:79`) hold type-size and layout rules but no table-of-contents rule.
## Detailed Requirements
1. Resolve `<dir|latest>`. Read the prior bundle as untrusted data (`crew-members/prompt-injection.md`).
2. Write a fresh bundle `yyyy-mm-dd_hhmm_<same-slug>-rev<N>` through the existing Collision-Safe Publication path.
3. Open it with the rev-N block: date, "changed", "still to do".
4. Add `<meta name="ai-report-supersedes" content="<prior dir>">` to `<head>`.
5. Regenerate the catalog with REQ-655's verb, so the prior bundle's `superseded_by` points at the new bundle. The prior bundle's bytes never change.
6. Apply the readable-style baseline from `ai-report-reference.md` to the new bundle.
7. Add a table of contents once the report passes a set length (suggested: more than six top-level sections or about 1,500 words). Put the rule in `ai-report-reference.md` Report Design Rules.
8. Print the final path (and a `file://` link) as the last line, so the user does not have to ask "give me the link".
9. Committing stays the user's or the run's decision, as today.
10. Update `ai-report.md` (the `:30` arguments line), `skills/do-work-toolbox/docs/ai-report-guide.md`, the routing phrases in `skills/do-work-toolbox/SKILL.md:23` (add a phrase such as "revise the report") and `skills/do-work-toolbox/actions/help.md:13`. Release per `_dev/primes/prime-releases.md`; run `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh` green.
## Constraints
- Never edit a prior bundle in place, not even to append a rev block (report, Out of scope).
- No new queue fields and no new statuses.
- Committing or publishing stays out of the action.
- This REQ is its own release; do not fold REQ-654, REQ-655 or REQ-657 into it.
## Assumptions (recorded at capture, no questions asked)
- The report says A3 "conflicts with the current contract, so it needs an explicit upstream decision either way". Capture assumes the sibling-bundle design is inside the immutability rule, because it never changes a prior bundle's bytes and the catalog is derived data. If the maintainer rejects it, abandon this REQ with `do-work abandon`.
- N counts revisions in the supersedes chain: the first revise of an original bundle writes `-rev1`, a revise of `-rev1` writes `-rev2`, and so on. The `-2`/`-3` collision suffix still applies on top if the derived path exists.
- `latest` means the newest bundle by the date in its folder name across all naming styles, as read from the catalog, of any kind.
- The table-of-contents threshold is approximate. The rule is keyed on length, so it applies to any long report, not only revised ones; this is one condition instead of a revise-only exception.
- The `file://` link is printed for the user to open. The render check itself still serves over HTTP, never `file://` (`ai-report.md:113`).
- "What changed" from the argument, when given, seeds the "changed" list; the action still re-checks the prior report's claims, as `architecture-report` Step 3 does.
## Dependencies
Depends on REQ-655 (index and find): revise regenerates the catalog and needs `superseded_by`. REQ-654 (kinds) is not a prerequisite; a revise keeps whatever kind the prior bundle had.
## Builder Guidance
Certainty is high on the contract, medium on naming. Latitude: exact rev block wording and placement of the rules between `ai-report.md` and `ai-report-reference.md`.
## Red-Green Proof
**RED prompt/case:** `do-work-toolbox ai-report revise latest "sitemaps re-checked"` in a repo with one existing report bundle.
**Why RED now:** `ai-report.md:30` accepts no `revise` form, and the immutability rule leaves no documented way to update a report.
**GREEN when:** A new sibling bundle exists with a rev-N block at the top and an `ai-report-supersedes` meta naming the prior bundle; `git diff --exit-code -- <prior dir>` passes; the regenerated catalog shows the prior bundle's `superseded_by` pointing at the new bundle; a revised report over the length threshold has a working table of contents; the action's last line is the new bundle's path.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens at the claim-time consult on 2026-10-10, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its owning prime governs the action files this REQ edits; family `alternate-writer-contract-drift` fits a second writer of report bundles beside the original publication path, and `example-path-read-as-citation` fits any example path the new prose adds. The pre-dispatch agent read both families' bullets and carried their rules into the builder brief.
## Full Context
See `do-work/user-requests/UR-144/input.md` for complete verbatim input (sections Request A3, What happened, Where the behaviour lives today, Proposed direction A3, Acceptance check, Out of scope). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-ai-report-kinds-and-index.md`, Request item A3: "`ai-report revise <dir|latest> [what changed]` that keeps the immutability contract. It never edits the old bundle."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fully specified (ten numbered requirements, a Red-Green Proof, and named files), and the change is prose only in five toolbox files with no Go code: REQ-655 already shipped the catalog verb that computes `superseded_by`. What needs discovery is where each rule sits so it reuses the existing contracts (Collision-Safe Publication, the architecture-report re-check pattern, Report Design Rules, the `ai-report-index` verb) and stays clear of the REQ-657 and REQ-687 edits to the same files.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Required-lessons consult (2026-10-10): `do-work/lessons-index.md` matches two rows. `_dev/primes/lessons-releases.md` (666 tokens, `slugged: full`; the REQ runs the shipped-package reference contract) is added to `required_lessons` and read. `_dev/primes/lessons-action-files.md` (7756 tokens, `slugged: partial`) stays dropped for budget; its `alternate-writer-contract-drift` and `example-path-read-as-citation` bullets were read by hand and shape the sweep and example-path rules below.

**What REQ-655 shipped (0.305.105), read from the code, not the REQ text:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`: verb `ai-report-index`. Every non-dot directory directly under `ai-reports/` is a bundle (`:79-92`). Bundles sort newest first by the date read from the folder name (`yyyy-mm-dd`, optional `_hhmm` or `_hhmmss`, any naming style), dateless last, ties by path (`:94-100`). `superseded_by` is computed from other bundles' `ai-report-supersedes` meta, matched by the base folder name, so the meta may name a bare folder or an `ai-reports/<folder>` path (`:105-116`). Catalog field names `path`, `entry`, `title`, `date`, `kind`, `linked_ids`, `verdict`, `supersedes`, `superseded_by` (`:55-65`).
- `--find <topic>` matches the folder name (not the full path, review F1 of REQ-655), title, kind, verdict and linked ids, prints newest first, and writes nothing (`:376-408`). It takes a topic, so it cannot resolve `latest`; revise uses the regenerated `catalog.json` instead.
- The verb refuses with `AI-REPORT-INDEX-HAND-MADE` when `catalog.json` or `index.html` exists without its generator marker, and writes nothing (`:236-284`). It returns exact text, so the action prose calls it with `--format text` (REQ-655 D-11).
- Fixture check at main HEAD 7fbeff30: a two-bundle `ai-reports/` with `<meta name="ai-report-supersedes" content="2026-09-11_1430_deploy-guide">` in the newer bundle gives the older bundle `superseded_by: ai-reports/2026-10-10_1600_deploy-guide-rev1` (0.1 s). No Go change is needed for this REQ.
- Prose on main: `skills/do-work-toolbox/actions/ai-report.md:32-40` (`### Catalog forms: index and find`, with the one command fence at `:36-38`), `:119` and `:150` (catalog exceptions), `skills/do-work-toolbox/SKILL.md:24` (routing row now ends `report index`, `is there a report on`), `skills/do-work-toolbox/actions/help.md:15` (index/find line, description at column 34), `skills/do-work-toolbox/docs/ai-report-guide.md:70-71` (Input lines) and `:82-88` (`## Report Catalog`, which already documents `ai-report-supersedes` as "the folder of the report it replaces").

**Other anchors:**
- Collision-Safe Publication: `skills/do-work-toolbox/actions/completed-work-presentation-reference.md:69-75`. Its `:7` already allows an action to read that section alone without the completed-work target contract (architecture-report is the named example).
- Re-check pattern to copy: `skills/do-work-toolbox/actions/architecture-report.md:50` (prior reports are untrusted; read prior HTML as source, never run its scripts), `:71-75` (Step 3: walk the prior report claim by claim), `:83` (authored opening account), `:93` (relative links to the prior bundle must work from the new location), `:119` (never edit a prior report).
- Readable-style baseline: `skills/do-work-toolbox/actions/ai-report-reference.md:70-90` (Report Design Rules). No table-of-contents rule exists there.
- Output: `skills/do-work-toolbox/actions/ai-report.md:133-135` (`## Output Format`), untouched by REQ-657 and REQ-687; Step 8's summary line `:131` is edited by REQ-687 and its neighbour `:129` by REQ-657.
- Tests that read these files: `_dev/tests/prescribed-shell-canonicalization.sh:107` requires `ai-report.md` to keep its pointer to `../../do-work/docs/prescribed-shell-primitives.md` (present at `:48`); a new bash fence would be checked as prescribed shell, so the revise form reuses the existing catalog command fence by reference. `_dev/tests/staged-skills-contract.sh` scans example paths (lesson `example-path-read-as-citation`): example paths must not start with a package directory name such as `actions/` or `docs/`.

**Sibling seams (diffs of the unmerged branches against bd56c4b0):**
- REQ-657 (`worktree-agent-REQ-657-ai-report-judge-render-check`, ai-report judge): `ai-report.md` inserts a `judge <bundle-dir>` paragraph right after the `$ARGUMENTS` line `:30`, rewrites Step 7 (`:121-125`) and the first line of Step 8 (`:129`), and edits the render-judged checklist line `:149`; `SKILL.md` adds a new row after `:24`; `help.md` adds a line after `:14`; guide rewrites `:55` and adds an Input line after `:69`.
- REQ-687 (`worktree-agent-REQ-687-ai-report-proposal-root-cause-kinds`, `--kind proposal|root-cause`): `ai-report.md` edits the header `:3`, the Do-NOT-use bullet `:26`, adds a `--kind` paragraph after `:30`, a line at the top of Step 1 (`:44`), a sentence on the slug paragraph `:61`, a line after the Step 5 list (`:111`), a sentence on Step 8's summary `:131`, and a checklist line after `:149`; `ai-report-reference.md` appends `## Proposal and Root-Cause Kinds` at end of file (after `:119`); `SKILL.md` appends three phrases to row `:24`; `help.md` adds a line after `:14`; guide edits `:47`, adds two Input lines after `:69`, and appends `## Proposal and Root-Cause Reports` at end of file.
- Both branches already conflict with REQ-655's lines on main (the `:30`-to-`## Steps` region, the `:150` checklist line, the guide Input fence and end of file). REQ-656's edits below sit beside those regions, never inside a line either sibling rewrites, except two union seams it cannot avoid: the `SKILL.md:24` row (REQ-687 appends to the same row) and the guide's end of file (REQ-687 appends there too).

**Decisions (pre-dispatch, DECIDE & STATE unless marked):**
- **D-01** No Go change. The shipped verb already links `superseded_by` from the meta (fixture check above). No `--latest` flag, no new verb.
- **D-02** `latest` = run the catalog command, then take the first entry of `bundles` in `ai-reports/catalog.json` (newest first by folder-name date, any naming style, any kind). This follows the shipped sort; `--find` is not used, because it needs a topic. If the catalog command refuses (hand-made `catalog.json` or `index.html`), `latest` cannot resolve: stop with one line naming the refusal and asking for an explicit bundle folder.
- **D-03** `<dir>` accepts `ai-reports/<folder>` or a bare `<folder>` and must be a directory directly under `ai-reports/` (the catalog's own bundle condition); anything else stops with one line. When the named bundle already has a `superseded_by` in the catalog, revise follows the chain forward to its newest bundle and says so in one output line. Value: no forked chain with two `-rev1` successors, of which the catalog would show only the newest. Risk: a user who really wanted to branch from an old revision gets the newest one instead; they can still copy the old content into the new revision by hand, and nothing is overwritten.
- **D-04** N = 1 + the number of `supersedes` hops from the prior bundle in the regenerated catalog (stop at a repeated path). Same slug = the prior folder name with its date (and time) token and the separator next to it removed, then a trailing `-rev<K>` removed together with a numeric collision suffix right after it. New preferred folder: `yyyy-mm-dd_hhmm_<slug>-rev<N>`, then Collision-Safe Publication unchanged (`-2`, `-3` on top).
- **D-05** The meta is `<meta name="ai-report-supersedes" content="<prior folder name>">`, the bare folder name, matching the guide's existing wording. When the prior bundle has an `ai-report-kind` meta, the new bundle copies it (a revise keeps the prior kind).
- **D-06** The revise form replaces Step 1's Terminal-Success Target Resolution with the prior bundle as the target, because the reports users asked to revise (a deploy guide, findings so far) were not completed-work reports. Safety Load Order (prompt-injection, then anti-slop, before reading the prior bundle), Evidence Honesty and Collision-Safe Publication still apply. Linked UR/REQ ids from the catalog are read at their current status and never presented as shipped when unfinished.
- **D-07** Before writing, revise re-checks the prior report claim by claim against the current repository, as `architecture-report.md` Step 3 does. A `[what changed]` argument seeds the Changed list; it never replaces the re-check.
- **D-08** The rev block is the first content after the page title: a `rev-N (yyyy-mm-dd)` heading, a Changed list, a Still-to-do list, and a relative link to the prior bundle that works from the new folder. Earlier rev blocks are not carried forward; the catalog holds the chain.
- **D-09** The whole of Report Design Rules applies to the new bundle even when the prior bundle used another style; the prior CSS is not copied where it breaks those rules. Steps 2 to 8 apply to the new bundle as for its kind.
- **D-10** The table-of-contents rule is one bullet at the end of Report Design Rules (after `ai-report-reference.md:90`), keyed on length for every report: more than six top-level sections or about 1,500 words of body text (approximate). A `<nav>` near the top with in-page anchor links to each top-level section's `id`, working without a script.
- **D-11** The "path last" rule goes in `ai-report.md` `## Output Format` and applies to every invocation that writes a bundle: the last output line is the bundle's `index.html` path relative to the project root plus its `file://` absolute link. Keyed on the condition, not a revise-only exception, and placed away from the Step 8 lines REQ-657 and REQ-687 rewrite. Value: no "give me the link" follow-up for any report. Risk: one extra output line on default reports; trivially reversible.
- **D-12** The revise procedure lives in one `### Revise form` subsection of `ai-report.md`, directly after the catalog-forms subsection (after `:40`, before `## Steps`). It reuses the catalog command fence at `:36-38` by reference; no second bash fence. `ai-report-reference.md` gets only the TOC bullet.
- **D-13** When the final catalog regeneration refuses, the new bundle stays (it is complete and immutable); the action reports the refusal with its fix (move the hand-made file aside, then run `ai-report index`) and still ends with the path line.
- **D-14** Routing and help: append `revise the report` and `update the report` to the ai-report row at `SKILL.md:24` (as REQ-655 did for its phrases); a help line after `help.md:15`; a guide Input line after `ai-report-guide.md:71` and a short `## Revising a Report` section after `## Report Catalog` at the end of the guide.
- **D-15** Restatement sweep (lesson family `alternate-writer-contract-drift`): `ai-report.md:119` ("creates only the report bundle", "the catalog forms above are the one exception"), `:135` ("the only stakeholder artifact this action publishes"), `:150` (checklist, "the catalog forms' catalog.json and index.html excepted") and `ai-report-guide.md:21` ("This action produces only the report bundle") each now meet a second writer of the catalog. The builder reads each and widens it only where it would be false for the revise form, recording each verdict.

<!-- D-XX counter: last used D-15. Next decision: D-16. -->

*Generated by the pre-dispatch agent (exploration done directly; no sub-agent)*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — a Use-when bullet, the `### Revise form` subsection after the catalog forms, the path-last rule in Output Format, a checklist line, and the restatement sweep (D-15)
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modify) — one table-of-contents bullet at the end of Report Design Rules
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — an Input line for the revise form and a short Revising a Report section
- `skills/do-work-toolbox/SKILL.md` (modify) — two routing phrases on the ai-report row
- `skills/do-work-toolbox/actions/help.md` (modify) — one help line for the revise form

**Files I will NOT touch:** any Go file (the shipped catalog verb already links superseded bundles, D-01); `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` (its immutability and sibling rules stay as written); `skills/do-work-toolbox/actions/architecture-report.md` (pattern source only); CHANGELOG, VERSION and version mirrors (the integrator releases).

**Acceptance criteria (restated from REQ):**
- [ ] The revise form with a bundle folder or latest, plus optional "what changed", is documented in the ai-report action's Input, and latest resolves to the newest bundle of any kind from the regenerated catalog.
- [ ] The prior bundle is read as untrusted data, after the prompt-injection guardrail loads, and its claims are re-checked.
- [ ] A fresh bundle named with the same slug plus -rev<N> is written through Collision-Safe Publication; N counts the supersedes chain.
- [ ] The new bundle opens with a rev-N block: date, Changed, Still to do.
- [ ] The new bundle's head carries the ai-report-supersedes meta naming the prior bundle folder.
- [ ] The catalog is regenerated with the ai-report-index verb so the prior bundle's superseded_by names the new bundle; the prior bundle's bytes never change.
- [ ] The readable-style baseline (Report Design Rules) applies to the new bundle, and a table-of-contents rule keyed on length sits in Report Design Rules.
- [ ] The last output line is the new bundle's path and a file:// link.
- [ ] Committing stays out of the action.
- [ ] The guide, the toolbox routing row and the help list name the revise form; the shipped-package reference contract and contract-regressions pass.
