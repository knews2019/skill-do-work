---
id: REQ-655
title: 'ai-report index and find build a derived catalog of every report bundle in any naming style'
status: claimed
created_at: 2026-10-09T21:10:00Z
user_request: UR-144
domain: backend
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-654, REQ-656, REQ-657]
batch: ai-report-modes
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/**", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/actions/help.md"]
claimed_at: 2026-10-10T12:52:19Z
---
# ai-report Index and Find Build a Derived Catalog of Every Report Bundle in Any Naming Style
## What
Two new forms of `ai-report`. `ai-report index` walks `ai-reports/`, reads every bundle whatever its naming style, and writes a derived catalog as `ai-reports/catalog.json` plus a static `ai-reports/index.html`. `ai-report find <topic>` searches that catalog and prints matching bundle paths, newest first. The mechanical walk is a new toolbox CLI verb next to `architecture-report-preflight`, not prose.
## Why
Report A2 (UR-144 input): 3 sessions in 3 repos asked "is there an ai-report on X" or "which recent ai-reports have design proposals, include rejected ones", and one ended with a hand-built HTML index. One consumer repo, counted 2026-10-09, has about 180 entries under `ai-reports/` in at least five folder-naming styles, loose files beside the bundles, 45 bundles without `index.html`, and two hand-made catalogs. Today "which report is current" is answered by reading folder names.
## Verified Facts (checked at 0.305.87)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go:45` (`architectureScan`) already scans `ai-reports/`, but `:25` matches only `^(\d{4}-\d{2}-\d{2}_\d{4})_architecture-report(?:-(\d+))?$`, and (report cites `:80-83`) only bundles with a nonempty `index.html`.
- `skills/do-work-toolbox/actions/ai-report.md:109`: "It also does not publish, host, or search for distribution targets." No search of existing reports exists in the toolbox.
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md:71`: existing artifacts are immutable.
## Detailed Requirements
1. New toolbox CLI verb (report suggests `ai-report-index`) that walks `ai-reports/` and accepts every naming style found, including at least `yyyy-mm-dd_hhmm_slug`, `yyyy-mm-dd_slug`, `yyyy-mm-dd-slug`, `REQ-NNNN-slug`, `slug-yyyy-mm-dd` and `yyyy-mm-dd_hhmmss_slug`. These are examples of the condition "a directory under `ai-reports/` is a bundle"; do not hard-code them as a closed list.
2. Entry file pick order: `index.html`, then `index.md`, `README.md`, `prompt.md`, `report.md`, then the first `.html`.
3. Per bundle, read: the `<title>` or first heading, the `ai-report-kind` meta when present (REQ-654), any UR/REQ id in the folder name or title, and the `ai-report-supersedes` meta when present (REQ-656).
4. Catalog fields per bundle: path, title, date, kind, linked UR/REQ, verdict or decision, `supersedes`, `superseded_by`. `superseded_by` is computed from other bundles' `supersedes`.
5. Write `ai-reports/catalog.json` and a static `ai-reports/index.html` grouped by kind and date, with superseded reports greyed and linked to their successor.
6. Loose files directly under `ai-reports/` (for example `.md`, `.patch`, `.html`) are skipped, not listed as bundles.
7. `find <topic>` is a substring match plus a UR/REQ id match over the catalog, printing paths newest first, including superseded bundles.
8. The catalog is regenerated, never hand-edited, so it does not break the immutability rule. No existing bundle file changes.
9. Wire both forms into `ai-report.md` (relax the `:109` "does not ... search" sentence for this form only), `skills/do-work-toolbox/docs/ai-report-guide.md`, the routing phrases in `skills/do-work-toolbox/SKILL.md:23` (add a phrase such as "is there a report on") and `skills/do-work-toolbox/actions/help.md:13`.
10. Focused Go tests on a fixture `ai-reports/` folder; release per `_dev/primes/prime-releases.md`.
## Constraints
- No new queue fields and no new statuses.
- Do not rename or migrate existing report folders; `index` reads them as they are.
- Existing bundles are immutable; only `catalog.json` and `index.html` at the `ai-reports/` root are written.
- Leave `architectureScan` behaviour unchanged unless sharing the walk is a clear simplification; if it is shared, its existing tests stay green.
- This REQ is its own release; do not fold REQ-654, REQ-656 or REQ-657 into it.
## Assumptions (recorded at capture, no questions asked)
- `find` regenerates the catalog before searching (the catalog is derived and cheap to rebuild), so it is never stale against the folder. If the builder measures this as slow on a large tree, reading an existing catalog is acceptable; record the choice.
- If `ai-reports/catalog.json` or `ai-reports/index.html` already exists and was not written by this verb (no generator marker), the verb refuses with a typed finding instead of overwriting a hand-made file. A file it wrote itself is overwritten freely.
- Kind falls back in this order: the meta, then a `<kind>-` prefix in the slug after the date part, then a known slug family such as `architecture-report`, then `unknown`.
- "Verdict or decision" is read only when the page marks it in a way the verb can find without guessing (a heading or meta); otherwise the field is empty. Do not build a prose parser for it.
- The report's "Proposed direction" adds a "may be stale" mark for a report whose linked REQs have all closed since it was written. It is not in the Request section, so it is optional: build it only if the REQ status lookup is cheap, and record the decision either way.
- The verb name follows the toolbox CLI's existing naming; `ai-report-index` is the report's suggestion, not a requirement.
## Dependencies
None upstream. REQ-656 (revise) depends on this REQ because it regenerates the catalog and relies on `superseded_by`.
## Builder Guidance
Certainty is high on the contract (the report gives the acceptance check). Latitude: verb name, where the walk lives, the HTML layout, and how the fallbacks are ordered beyond what is listed.
## Red-Green Proof
**RED prompt/case:** A fixture `ai-reports/` with one bundle per naming style above, one bundle holding only `README.md`, one proposal bundle superseded by a later one, and a loose `.patch` file; run the new verb, then `find proposal`.
**Why RED now:** No such verb or form exists; `architectureScan` sees only `_architecture-report` bundles with `index.html`.
**GREEN when:** `catalog.json` lists every bundle exactly once, skips the loose file, and shows the older proposal's `superseded_by` pointing at the newer one; `index.html` is written; `git status` shows no other file changed; `find proposal` prints both proposal bundles, newest first, including the superseded one.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18987 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers changing do-work-cli internals and condition-based classifiers; family `closed-enumeration-for-a-condition` is the risk when matching naming styles.
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over budget; `slugged: partial`). Matching reason: its owning prime governs the action and guide edits in requirement 9.
## Full Context
See `do-work/user-requests/UR-144/input.md` for complete verbatim input (sections Request A2, What happened, Where the behaviour lives today, Proposed direction A2, Acceptance check). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-ai-report-kinds-and-index.md`, Request item A2: "`ai-report find <topic>` and `ai-report index`. One command answers "is there a report on X" and "which reports have design proposals, including rejected ones"."*
