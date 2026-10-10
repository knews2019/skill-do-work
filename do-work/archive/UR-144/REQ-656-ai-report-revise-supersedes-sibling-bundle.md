---
id: REQ-656
title: 'ai-report revise writes a new sibling bundle that supersedes the prior one'
status: completed
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
builder_handback_at: 2026-10-10T16:04:10Z
integration_at: 2026-10-10T16:29:26Z
review_at: 2026-10-10T16:39:00Z
kb_status: pending
commit: 65547570bdbe7e38a008b4dc8cd507fafa4d3eb0
heavy_verified_at: 2026-10-10T16:46:45Z
heavy_verified_revision: 65547570bdbe7e38a008b4dc8cd507fafa4d3eb0
completed_at: 2026-10-10T16:47:34Z
release_at: 2026-10-10T16:47:34Z
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
- [x] **[PLAN]:** Prose-only change in five toolbox files, following the brief's decided shape (D-01 to D-15). Reuse existing contracts by name (catalog command fence, Safety Load Order, Terminal-Success Target Resolution, Evidence Honesty, Collision-Safe Publication, `architecture-report.md` Step 3, Report Design Rules) instead of restating them. No Go, no tests, no new bash fence. Proof: RED probe at base, GREEN walkthrough in a `mktemp -d` scratch repo for rev1 and rev2, GREEN probe, reference contract, contract-regressions. (from the builder hand-back)
- [x] **[APPLY]:** Edits exactly as planned, inside the five-file write set; each edit kept local to its anchor so the REQ-657 and REQ-687 merges stay small. (from the builder hand-back)
- [x] **[UNIFY]:** `git diff 7fbeff30 --stat`: 5 files, 26 insertions, 5 deletions (`SKILL.md`, `ai-report-reference.md`, `ai-report.md`, `help.md`, `ai-report-guide.md`). Checks from the worktree root: `REQ-656-probe.sh` exit 0 (1.2 s); `shipped-package-reference-contract.sh` exit 0 (1.1 s); `contract-regressions.sh` exit 0 (24 s); `git diff 7fbeff30 --check` clean; `staged-skills-contract.sh` refused as heavy-only (exit 2), left to the integrator's heavy drain. All five files read in the full diff: probe tokens inside `### Revise form`, no `#` line inside the subsection, no `CLAUDE.md` or `_dev/` citation, `completed-work-presentation-reference.md` and `architecture-report.md` untouched. (from the builder hand-back)
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

<!-- D-XX counter: last used D-22. Next decision: D-23. -->

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

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work-toolbox/SKILL.md` (modified)
- `skills/do-work-toolbox/actions/help.md` (modified)

**What was done:** `ai-report.md` gained a Use-when bullet and a `### Revise form` subsection after the catalog forms: seven steps that resolve `<dir|latest>` from the regenerated catalog (following `superseded_by` forward), read the prior bundle as untrusted data, re-check its claims, name the new folder `yyyy-mm-dd_hhmm_<slug>-rev<N>` through Collision-Safe Publication, write it with a rev-N block and an `ai-report-supersedes` meta, regenerate the catalog, and leave committing to the user. `## Output Format` now ends every bundle-writing run with the bundle path and a `file://` link; the Step 6 paragraph and the last checklist line were widened for the revise form's catalog write, and one checklist line was added. `ai-report-reference.md` Report Design Rules gained a length-keyed table-of-contents bullet. The guide gained an Input line, a widened `:21` sentence and a `## Revising a Report` section; the toolbox routing row and the help list name the form.

Hand-back merge `c44d9f29..3e535501` (integrator): four files conflicted with REQ-657 (ai-report judge, already on main) and were resolved by keeping both texts: the toolbox `SKILL.md` keeps one `ai-report` row holding the judge and the revise phrases; `help.md` and the guide's Input block keep both new lines; the `ai-report.md` checklist keeps REQ-657's `ai-report-judge` render line, the builder's widened catalog-exception line and the new revise line. Seam: the rev block's relative link to the prior bundle is a parent-directory link, which `ai-report-judge` (it serves only the bundle folder) always reports as `AI-REPORT-JUDGE-BROKEN-LINK`; revise step 5 now names that finding as expected and says to check the link once the bundle is written, the same rule REQ-657 gave `architecture-report.md` Step 5.

Re-merge `c44d9f29..6c837fab` (integrator, after Qualification found requirement 10 unmet): builder-branch commit `8a56cd5b` names the revise form on the `ai-report.md` `$ARGUMENTS` line; the conflict with REQ-657's version of that line was resolved by keeping its `judge`/`index`/`find` clause and appending the revise clause.

Review fixes (integrator, re-merged with the same pre): the builder branch was fast-forwarded to `6c837fab` because the lines involved exist only after the REQ-657 merge, then `c434be7c` (review F1, F2, F3, F5, F6) merged as `85395d8e` and `49f9b215` (re-review N1-N3) merged as `65547570`. A revise now ends Step 7 with the one expected prior-bundle link finding, and Step 7, Step 8, the render checklist line and the guide's render paragraph accept it; a revise of a report that is not a completed-work report keeps the prior report's section order (revise step 5 and the narrative checklist line); the step 3 re-check is the provenance ledger for every revise; the Do-NOT-use "unfinished" bullet, the `$ARGUMENTS` "one invocation" sentence and the guide's "every invocation creates a bundle" sentence were made true for the revise and catalog forms. Final range `c44d9f29..65547570`: 5 files, 34 insertions, 13 deletions.

## Decisions

(from the builder hand-back, verbatim; D-01 to D-15 are in `## Exploration`)

- **D-16 (DECIDE & STATE)** N counting when the catalog command refused and the user gave an explicit `<dir>`: step 4 says to read the hop count from each bundle's `ai-report-supersedes` meta instead of the catalog. The brief only covered `latest` under refusal. Without this fallback an explicit-folder revise would have no defined N. Reversible: one parenthetical.
- **D-17 (DECIDE & STATE)** Step 1's chain-forward rule says "the resolved bundle (from `latest` or `<dir>`)" instead of "the named bundle". Reason: the walkthrough showed a same-minute tie sorts the older revision first, so `latest` alone can pick a superseded bundle. Following `superseded_by` forward gives the right bundle with no Go change.
- **D-18 (DECIDE & STATE)** The Output Format path-last rule is written as one sentence with a semicolon joining the HTTP note, to keep it as "one sentence" per the brief.
- **D-19 (DECIDE & STATE)** Guide Input line: the command is 58 characters, longer than the column-43 description slot, so the description follows two spaces (per the brief).
- **D-20 (DECIDE & STATE)** No example relative link like `../<prior folder>/index.html` in shipped prose. "a relative link to the prior bundle that works from the new folder" is enough, and a `../` token risks the citation scan.
- **D-15 sweep verdicts:**
  - `ai-report.md` Step 6, "creates only the report bundle": false for revise (it also regenerates the catalog). Widened to "..., plus the regenerated catalog for the revise form."
  - Same paragraph, "the catalog forms above are the one exception to searching": false for revise (it reads the catalog to resolve `latest` and the chain). Widened to "the catalog forms and the revise form above are the only exceptions to searching".
  - `ai-report.md` Output Format, "The timestamped folder is the only stakeholder artifact this action publishes": still true. The catalog is derived navigation data, not a stakeholder artifact, and the catalog forms already wrote it under this same sentence. Unchanged. The new path-last sentence follows it in the same paragraph.
  - `ai-report.md` last checklist line, "(the catalog forms' `catalog.json` and `index.html` excepted)": false for revise. Widened to "(the catalog forms' and the revise form's ...)".
  - `ai-report-guide.md:21`, "This action produces only the report bundle": false for revise. Widened to "... (a revise also regenerates the report catalog)."
- **D-21 (DECIDE & STATE, integrator)** Merge seam with REQ-657: revise step 5 names the judge's `AI-REPORT-JUDGE-BROKEN-LINK` finding for the prior-bundle link as expected, so a revise can finish Step 7. Without it, Step 7's "rerun until the verdict is `pass`" could never be met for any revision. Wording follows `architecture-report.md` Step 5 and avoids a `../` token (D-20).
- **D-22 (DECIDE & STATE, integrator, from review F2)** A revise of a report that is not a completed-work report (the UR's deploy guide, "findings so far") keeps the prior report's own section order instead of Step 5's shipped-work narrative, and the step 3 re-check is the provenance ledger for every revise, because Step 1, which builds the ledger, does not run for a revise. Value: the UR's real cases revise cleanly. Risk: one more clause in the revise step and the narrative checklist line; reversible.

## Discovered Tasks

(from the builder hand-back; impact token is the builder's)

- impact-minor: `ai-report-index` breaks same-date ties by path ascending, so of two bundles written in the same minute the older `-rev1` sorts ahead of `-rev2` and becomes `latest`. The revise prose works around it by following `superseded_by` (D-17). The Go sort could instead put a bundle that supersedes another first on a tie, or sort ties by path descending → report only

(integrator sweep at Step 8)

- impact-negligible: review F4: `latest` of any kind can pick an architecture-report bundle, and its `-rev1` revision is invisible to the architecture scanner's folder regex; follows the capture assumption → report only
- impact-negligible: review B9: the routing phrase `update the report` may also catch "update the architecture report" → report only
- impact-negligible: closing UR-144 moves REQ-654 to `do-work/archive/UR-144/`, so the queued `do-work/working/REQ-687-ai-report-kind-proposal-root-cause.md:81` cites a flat path that no longer exists; the REQ-687 integrator can re-point it → report only

## Qualification

**Gate records (`advance --diff-range c44d9f29..3e535501`):** `qualify` satisfied (P-A-U boxes ticked from the hand-back, no debug-artifact or output-primitive finding). `scope-drift` reported one warning, `SCOPE-DECLARED-NOT-TOUCHED` for the path `### Revise form`: the scope parser read the backticked heading name in the `ai-report.md` Scope bullet as a declared path. It is not a file; judged a false positive.

**Cumulative range after the re-merges:** `c44d9f29..6c837fab` at qualification; final `c44d9f29..65547570` after the review fixes (same five files).

**Scope comparison:** declared five files (`ai-report.md`, `ai-report-reference.md`, `ai-report-guide.md`, toolbox `SKILL.md`, `help.md`); `git diff c44d9f29..3e535501 --stat` touches exactly those five (26 insertions, 5 deletions, merge seam line included). No Go, test, CHANGELOG or version file. `git diff --check` clean.

**Requirement trace (read in the merged files, not the hand-back):**
1. Resolve `<dir|latest>`, read the prior bundle as untrusted data: `ai-report.md` Revise form steps 1 and 2 (catalog first, `latest` = first catalog entry, chain followed forward, prompt-injection then anti-slop loaded before reading). Met.
2. Fresh `yyyy-mm-dd_hhmm_<same-slug>-rev<N>` through Collision-Safe Publication: step 4 (N from supersedes hops, slug stripping rule, example path under `ai-reports/`). Met.
3. rev-N block with date, Changed, Still to do: step 5. Met.
4. `ai-report-supersedes` meta in `<head>`: step 5, bare folder name as the guide's catalog section already documents. Met.
5. Catalog regenerated so `superseded_by` points forward, prior bytes unchanged: step 6 plus step 5's last sentence; the builder's scratch walkthrough confirmed the chain original -> rev1 -> rev2 with `git diff --exit-code` 0 on prior bundles. Met.
6. Readable-style baseline: step 5 applies the full Report Design Rules even where the prior style differs. Met.
7. Length-keyed table of contents in Report Design Rules: new last bullet in `ai-report-reference.md`. Met.
8. Path and `file://` link as the last line: `## Output Format`, keyed on every bundle-writing run. Met.
9. Committing stays out: step 7. Met.
10. Guide, routing row, help: met (guide Input line plus `## Revising a Report`; one `ai-report` routing row with `revise the report`, `update the report`; help line). The `$ARGUMENTS` line (`ai-report.md` `## Input`, the `:30` line the REQ names) still lists only `UR-NNN`, `REQ-NNN`, `most recent` or blank and does not mention the revise form. Fixed by the integrator on the builder branch (`8a56cd5b`) and re-merged with the same pre as `6c837fab`: the line now keeps REQ-657's `judge`/`index`/`find` clause and adds that the `revise` form writes a new revision of an existing report. Met.

**Restatement check by the builder (D-15):** four sentences widened, one judged still true, each verdict recorded in `## Decisions`.

**Merge seam judged:** REQ-657's `ai-report-judge` serves only the bundle folder, so a rev block's parent-directory link would fail Step 7 forever; D-21 adds one sentence to step 5. Passed to the reviewer as a meaning-changing seam.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge `6c837fab` (machine quiet before launch: 1-minute load 2.39, no other gate running), then `advance REQ-656 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-131527/REQ-656-probe.sh`.
**Result:** ✓ Repository gate passed on the first run (exit 0, gate wall 150 s; do-work-cli 902 tests in 71 s, slowest file `internal/finalization/finalization_recovery_test.go` 22.27 s under the 30 s limit; queue-kanban 420 tests in 48 s, slowest 19.95 s; contract regressions passed). Gate records `green-gate` and `run-blocked-check` satisfied (probe exit 0, "REQ-656 probe passed."); `scope-drift` carries only the `### Revise form` false positive judged in `## Qualification`. A direct probe rerun at `6c837fab` also passed.

**After the review-fix re-merges:** the same gate argv ran at `85395d8e` (load 2.35 before launch): exit 0, gate wall 142 s; and at the final merge `65547570` (load 4.20 before launch): exit 0, gate wall 143 s, do-work-cli 902 tests in 69 s, slowest file 22.95 s. `advance` refuses gate input past this phase, so the green record was written with `record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh` at `65547570`. `REQ-656-probe.sh` run directly at `65547570`: exit 0.

**Red-green validation:** non-behavioral prose change (tdd: false); regression evidence from the builder hand-back, traced to `## Red-Green Proof`:
- `REQ-656-probe.sh`: ✗ at base `7fbeff30` (exit 1, 23 failures: no revise form or its tokens, no `file://` in Output Format, no TOC rule, no guide/routing/help lines) → ✓ after (exit 0).
- GREEN walkthrough in a `mktemp -d` scratch repo, prose followed literally for `ai-report revise latest "sitemaps re-checked"`: new `2026-10-10_1903_deploy-guide-rev1` with the rev block first after the title, the supersedes meta, a 7-link TOC with no missing ids; regenerated catalog shows the prior bundle `superseded_by` rev1; `git diff --exit-code` on the prior bundle exit 0; last line the path plus `file://` link. A second revise gave `-rev2` (N = 2) and the chain original -> rev1 -> rev2.

**New tests added:** none (prose REQ; the run probe is the regression check).

**Heavy verification plan:**
- Range: `c44d9f2941c34ea5c30c24e5bf9bba9bc06273e8..65547570bdbe7e38a008b4dc8cd507fafa4d3eb0` (re-planned after the review fixes; same single lane as at `6c837fab`)
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all five changed paths match subtree `skills`

*Verified by work action*

## Review

**Overall: 79%** | 2026-10-10T16:38:59Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 80% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `ai-report.md:150`, `:170`, `ai-report-guide.md:55` still require a final `pass` verdict, but a revise always ends `fail` on the expected prior-bundle link (reproduced: verdict=fail, exit 1); widen Step 8, the checklist and the guide to accept that one finding — impact-user-visible → report only
- F2 `ai-report.md:53` "Steps 2 to 8 apply ... as for its kind": no kind is defined at this merge, so a non-completed-work report (the UR's deploy guide) is forced into Step 5's shipped-work narrative and Step 6's provenance-ledger check; keep the prior section order and use the step 3 re-check as the ledger — impact-user-visible → report only

**Minor findings:** F3 `ai-report.md:27` Do-NOT-use "target is unfinished" bullet can block a revise of a findings-so-far report — impact-user-visible → report only; F4 `latest` of any kind can revise an architecture-report into a `-rev1` folder its scanner never sees (follows the capture assumption) — impact-negligible → report only; F5 (nit) D-21 "once the bundle is written" adds no timing in ai-report — impact-negligible → report only; F6 (nit) guide `:9` "Every invocation creates a fresh bundle" and `ai-report.md:31` "One invocation covers one UR or one REQ" are stale, from REQ-655/657 — impact-negligible → report only
**Acceptance:** Partial — implementation and integration stages: probe and reference contract pass, scratch catalog chain and tie-forward verified; judge returns fail on the expected link, so Step 8 is unreachable for a revise
**Restatement sweep:** redefined the revise input form, `latest`, the path-last Output Format rule, the length-keyed TOC rule, the "only the bundle / only exceptions to searching" boundary (now also revise), and (via D-21) the end verdict of Step 7 for a revise. Stale: F1 (`ai-report.md:150`, `:170`, guide `:55`), F2 (`:53`, `:169`), F3 (`:27`), F6 (guide `:9`, `:31`). Inherited from REQ-655 (input forms and boundary): `:31` now names revise, Output Format still true, guide `:21` widened. Inherited from REQ-657 (judge verdicts): `:144`, `:150`, `:170`, guide `:55` → F1. Checked, not stale: README.md:112, toolbox help.md:14, command-line-guide.md:26 (no new verb), architecture-report.md, stakeholder-report.md (TOC rule applies by length, as intended), completed-work-presentation-reference.md:7 (architecture-report named as an example only)
**Suggested testing:** 3 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

**Integrator disposition:** Approve at 79% (Acceptance Partial) continues to completion. Agreed and fixed F1, F2, F3, F5 and F6 on the builder branch (`c434be7c`, merged as `85395d8e`); F4 stays report only (it follows the capture assumption that `latest` is of any kind). Delta re-review of `6c837fab..85395d8e`: 92%, Approve, F1/F2/F3/F5/F6 closed, three new findings N1 (narrative checklist line, impact-user-visible), N2 (ledger only for non-completed-work revises, impact-negligible), N3 (Step 7 rerun rule, impact-negligible), each fixed with the reviewer's proposed text (`49f9b215`, merged as `65547570`); a second delta check of `85395d8e..65547570` closed N1-N3 with no new finding. Anti-bloat B1-B9 accepted as small and within the Constraints; B9 (`update the report` routing phrase may catch "update the architecture report") → report only. Reports: `do-work/runs/work-2026-10-10-131527/REQ-656-review.md`, `REQ-656-rereview.md`.

## Lessons Learned

**What worked:** The builder's literal scratch walkthrough (two revisions in one minute) found the same-minute catalog tie before review, and the chain-forward rule fixed it in prose with no Go change (D-17). Fast-forwarding the builder branch to the merge commit let review fixes on main-only lines (REQ-657's judge text) go through the builder branch and re-merge with the same pre without conflicts.
**What didn't:** The walkthrough skipped Step 7 (render check), so neither the builder nor the first merge saw that a parent-directory link back to the prior bundle always fails `ai-report-judge`, which serves only the bundle folder. The first seam fix (D-21) widened only the revise step; Step 7, Step 8, two checklist lines and the guide still demanded `pass` (review F1, re-review N3). When a step gets an expected-failure exception, grep every line that restates that step's success condition.
**Worth knowing:** "Steps 2 to 8 apply as for its kind" assumed kinds that do not exist yet (REQ-654 was cancelled, REQ-687 is unmerged); a revise of a non-completed-work report needed its own section-order rule (D-22). REQ-687 (`--kind proposal|root-cause`) edits the same Do-NOT-use bullet and checklist area, so its integrator meets these lines as a seam. `ai-report-index` sorts same-minute ties by path ascending (report-only discovered task).

## Orientation

Now a stakeholder report can be updated as a linked revision: `ai-report revise <dir|latest>` writes a new `-rev<N>` bundle that supersedes the old one in the report catalog, and every report run ends with its path and a `file://` link. Lives in the do-work-toolbox ai-report action (prose; `_dev/primes/prime-action-files.md` area). Primes spot-checked: `prime-action-files.md` and `prime-releases.md` paths unchanged by this REQ, not stale.

## Heavy Verification Plan

- Base: `c44d9f2941c34ea5c30c24e5bf9bba9bc06273e8`
- Target: `65547570bdbe7e38a008b4dc8cd507fafa4d3eb0`
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — reasons: `skills/do-work-toolbox/SKILL.md`, `actions/ai-report-reference.md`, `actions/ai-report.md`, `actions/help.md`, `docs/ai-report-guide.md` each matched subtree `skills`

## Heavy Verification Result

- Target revision: `65547570bdbe7e38a008b4dc8cd507fafa4d3eb0`
- Execution revision: `65547570bdbe7e38a008b4dc8cd507fafa4d3eb0` (detached checkout `.git/work-run-work-2026-10-10-131527/drain-head-REQ-656`, removed afterwards)
- `staged-skills`: exit 0, executed (not reused), 37 s

## Timing

Observed 2026-10-10T15:51:30Z to 2026-10-10T16:46:45Z: 55m 15s total, 30m 00s attributed across 9 events, 25m 15s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 11m 48s | 4 |
| review | 9m 53s | 2 |
| planning | 3m 49s | 1 |
| exploration-preflight | 3m 38s | 1 |
| handback-merge | 52s | 1 |

Slowest stage: verification-gate / repository gate after review fixes, 5m 23s, outcome success.

Timing notes (integrator): no builder-work event was recorded, because the hand-back had already landed when integration began (`builder_handback_at` comes from the builder commit's committer date). The two post-review gate events were recorded after both gates finished, so each spans from its own start to the recording instant; the gate walls themselves were 142 s and 143 s (see `## Testing`).
