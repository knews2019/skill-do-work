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
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/actions/help.md"]
claimed_at: 2026-10-10T12:52:19Z
route: C
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-10-10T13:20:51Z
  basis:
    - Route C
    - 8-file write set
    - 2 subsystems involved
    - 9 acceptance criteria
planning_at: 2026-10-10T13:21:45Z
required_lessons: ["_dev/primes/lessons-releases.md"]
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
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (19046 tokens at the claim-time consult, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers changing do-work-cli internals and condition-based classifiers; family `closed-enumeration-for-a-condition` is the risk when matching naming styles.
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens at the claim-time consult, over budget; `slugged: partial`). Matching reason: its owning prime governs the action and guide edits in requirement 9.
## Full Context
See `do-work/user-requests/UR-144/input.md` for complete verbatim input (sections Request A2, What happened, Where the behaviour lives today, Proposed direction A2, Acceptance check). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-ai-report-kinds-and-index.md`, Request item A2: "`ai-report find <topic>` and `ai-report index`. One command answers "is there a report on X" and "which reports have design proposals, including rejected ones"."*

---

## Triage

**Route: C** - Complex

**Reasoning:** A new feature with several parts: a new toolbox CLI verb with two forms (write the catalog, find in it), a catalog schema with a computed `superseded_by`, a static HTML render, a typed refusal for hand-made files, focused Go fixture tests, and wiring into four toolbox prose files. Ten numbered requirements plus six recorded assumptions, so the shape needs a plan before exploration.

**Planning:** Required

## Plan

**Approach.** One new toolbox CLI verb, `ai-report-index`, in a new file `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`, registered in `commands.go` beside `architecture-report-preflight`. One in-memory walk builds the catalog; the two forms differ only in what they do with it.

- `ai-report-index` walks `ai-reports/`, then writes `ai-reports/catalog.json` and `ai-reports/index.html`.
- `ai-report-index --find <topic...>` runs the same walk in memory and prints matches newest first. It neither reads nor writes `catalog.json`.

**Tasks, in order**

1. **Walk and catalog** (requirements 1 to 4, 6). A bundle is any directory directly under `ai-reports/` whose name does not start with `.`. Loose files and symbolic links are never bundles. For each bundle:
   - Entry file (first match): `index.html`, `index.md`, `README.md`, `prompt.md`, `report.md`, the first `.html` by name, then the first `.md` by name (D-03). With no entry file, the bundle is still listed, with an empty entry and the folder name as its title.
   - Title: for HTML, the `<title>` text, else the first `<h1>`-`<h6>`. For Markdown, the first `#` heading line. Else the folder name. Strip tags, unescape entities, collapse whitespace.
   - Date: the first `yyyy-mm-dd` in the folder name, plus a `_hhmm` or `_hhmmss` that directly follows it. This is a condition on the name, not a list of naming styles. Stored as `2026-10-09T19:16`, `2026-10-09T19:16:05` or `2026-10-09`. Empty when the name carries no date (D-04).
   - Kind: the `ai-report-kind` meta first. Then a known kind that appears as a whole kebab word run in the slug after the date part. Then `unknown` (D-05).
   - Linked ids: every `UR-n` / `REQ-n` in the folder name or title, matched case-insensitively and stored uppercase with no duplicates.
   - Verdict: the `ai-report-verdict` meta only, else empty (D-06).
   - Supersedes: the `ai-report-supersedes` meta, matched to a bundle by its last path segment. Stored as that bundle's path, or as the raw value when nothing matches.
   - `superseded_by`: computed from the other bundles' `supersedes`. When several bundles name the same one, the newest wins.
   - Sort order: date descending, undated last, ties by path ascending.
   - Meta parsing reads `<meta>` tags with either attribute order and either quote style.
2. **Write the outputs** (requirements 5 and 8). Before writing anything, check both outputs. If `ai-reports/catalog.json` exists without the top-level `"generator": "do-work-cli ai-report-index"` field, or `ai-reports/index.html` exists without `<meta name="generator" content="do-work-cli ai-report-index">`, refuse with outcome `refused`, finding code `AI-REPORT-INDEX-HAND-MADE`, the path, fixability manual, and write nothing. Otherwise publish each file with `rootedPublishFile` (create, or replace when present), the way `portfolio.go` publishes its canonical file. No git transaction and no `--commit` (D-02). The output is deterministic: no generation timestamp, so a rerun on an unchanged tree is byte-identical. `index.html` is static, with inline CSS, a light and dark scheme, and every value HTML-escaped. It groups bundles by kind, then by date newest first. Superseded rows are greyed and link to their successor's entry. A missing `ai-reports/` directory is a success that writes nothing and says so.
3. **Find** (requirement 7). The match is a case-insensitive substring over path, title, kind, verdict and linked ids. When the topic is a `UR-n` or `REQ-n` id, linked ids are also compared with leading zeros ignored. Each output line starts with the bundle path, then a tab and the title. Superseded rows end with `(superseded by <path>)`. Superseded bundles are always included. No match is a success with a one-line "no report matches" message.
4. **Tests** (requirement 10), in a new file `report_index_test.go` built on one shared fixture: the six naming styles, a `README.md`-only bundle, a bundle with only a named `.md`, a proposal superseded by a later proposal, and loose `.patch` and `.md` files. Four focused tests (names fixed so the probe can find them):
   - `TestAIReportIndexCatalogsEveryBundleNamingStyleOnce`: every bundle exactly once, loose files absent, entry, title, date and ids correct, and no fixture file's bytes changed.
   - `TestAIReportIndexLinksSupersededProposalToSuccessor`: `superseded_by` is set, and `index.html` greys the row and links to the successor.
   - `TestAIReportIndexFindListsMatchesNewestFirstIncludingSuperseded`: `--find proposal` prints both proposal bundles, newer first.
   - `TestAIReportIndexRefusesHandMadeCatalog`: a hand-made `catalog.json` is refused with nothing written, and a file the verb wrote itself is overwritten on the next run.
   - Update `TestHandlersRegisterCanonicalToolboxCommands` for the new verb.
5. **Prose wiring** (requirement 9):
   - `skills/do-work-toolbox/actions/ai-report.md`: Input gains `index` and `find <topic>`. Add one short catalog-forms branch before Step 1 that runs the CLI verb, prints its output, and stops (Steps 1 to 8 do not apply). Narrow the Step 6 sentence at `:109` and the Verification Checklist item at `:140` so the catalog forms are the one exception to "does not search".
   - `skills/do-work-toolbox/docs/ai-report-guide.md`: the Input block gains the two forms, plus a short "Report catalog" section.
   - `skills/do-work-toolbox/SKILL.md:23`: add the routing phrases `report index` and `is there a report on`.
   - `skills/do-work-toolbox/actions/help.md`: add one help line for the two forms.

**Decisions (pre-dispatch, best judgment; no question was put to the user)**
- **D-01**: `find` walks the folder in memory on every call and never writes or reads `catalog.json`. Reasoning: this follows the capture assumption that the result is never stale, keeps `find` read-only, and means `find` never trips the hand-made-file refusal. Value: `find` has no side effects. Risk: slower on a very large tree. The walk reads one entry file per bundle, so 180 bundles is cheap. Reversible.
- **D-02**: outputs are published directly with `rootedPublishFile`, not through `runTransaction`. Reasoning: the git transaction refuses dirty targets, so a second `index` run before the first catalog was committed would be refused. `portfolio.go` is the precedent for a derived file published without a transaction. Committing stays the user's or the run's decision, as the source report asks. Value: a rerun just works. Risk: no automatic rollback if the second file fails after the first was written. A rerun repairs that, because both files are derived.
- **D-03**: after "first `.html`", the pick order adds "first `.md`". Reasoning: this repo's own `ai-reports/2026-08-26_1709_architecture-report/` holds only `architecture-report.md`, and the source report counts named Markdown entry files. Builder Guidance gives latitude on fallbacks beyond the listed ones.
- **D-04**: undated bundles (for example `REQ-NNNN-slug`) get an empty date and sort last. There is no file-time or git-date fallback, because file times are unreliable in a checkout. Value: deterministic output. Risk: an undated bundle cannot be ordered against dated ones.
- **D-05**: the known-kind table holds only the kinds that writers produce today: `proposal` and `root-cause` (REQ-687 (ai-report proposal and root-cause kinds)) and `architecture-report`. A code comment marks it illustrative, and the meta always wins. One table, matched as a whole kebab word run anywhere in the slug after the date, covers both the capture's "`<kind>-` prefix" and "known slug family" fallbacks with one rule. Risk: a slug that only mentions a kind word gets that kind. Acceptable for a derived, regenerated catalog.
- **D-06**: verdict comes from an `ai-report-verdict` meta only. No heading or prose parsing, per the capture assumption. No writer is changed to emit it in this REQ. The guide documents the meta name.
- **D-07**: the optional "may be stale" mark is not built. Reasoning: it needs REQ status lookups across `do-work/archive/` and is outside the Request section. Recorded here, as the capture assumption asks.
- **D-08**: `architectureScan` is left unchanged. Its regex, watermark and `index.html`-only rules serve a different job, so sharing the walk is not a clear simplification.
- **D-09**: no Just recipe for the new verb. Not every CLI verb has one (`advance` does not), and a recipe would add `justfile`, `skills/do-work-board/justfile.template` and the command-line-guide list to this REQ for no requirement. The action calls the launcher directly, as `architecture-report.md` does.
- **D-10**: verb name `ai-report-index`, with `--find` as the read-only form, following the `--scan` / `--publish` option style of `architecture-report-preflight`.
<!-- D-XX counter: last used D-10. Next decision: D-11. -->

**Plan validation**
- Requirement coverage: requirements 1 to 4 and 6 are task 1; 5 and 8 are task 2; 7 is task 3; 10 (tests) is task 4, and its release half belongs to the integrator; 9 is task 5. Constraints: no queue fields or statuses; no bundle is renamed or changed (task 4's byte snapshot); `architectureScan` is unchanged (D-08); own release. All covered.
- No orphan tasks.
- Scope sanity: five tasks, over the three-task comfort line. Splitting is not worth it: tasks 1 to 4 are one Go file and its tests, and task 5 is four short prose edits. Kept as one REQ, flagged here.
- Consumer field contract: REQ-656 (ai-report revise) consumes `catalog.json`. It needs, per bundle, `path` (identity), `date` (for `latest`), `supersedes` and `superseded_by`. All are in the schema above. `find` output starts each line with the bundle path, so a later reader can split on the tab.

*Generated by the pre-dispatch agent (Route C plan; no separate Plan agent was spawned)*

## Exploration

**Required-lessons consult (claim time).** I read `do-work/lessons-index.md`. Added `_dev/primes/lessons-releases.md` (666 tokens, `slugged: full`). It matches through `prime_files` (`prime-releases.md`), and I read it: its one family, `canonical-link-outlives-its-target`, matters only if the integrator's release links an archived record. The two over-budget satellites stay dropped, with refreshed costs in the dropped section above. No listed lesson file was missing.

**Code (paths under `skills/do-work/tools/do-work-cli/internal/toolboxcommands/`)**
- `commands.go:11-19` holds the command-name constants and `:21-31` the `Handlers()` map. `cmd/do-work-cli/main.go:77` merges that map into the CLI, so a new constant plus a map entry is the whole registration. `commands.go:33-48` (`usageResult`, `toolboxFinding`) and `:76-78` (`exactOutputResult`) are the result helpers to reuse.
- `commands_test.go:5-15` (`TestHandlersRegisterCanonicalToolboxCommands`) lists the names and asserts `len(handlers) != 7`. This is a merge seam (see below).
- `architecture.go:25` (`architectureName`) and `:45-102` (`architectureScan`) show the existing `ai-reports/` read. It treats a missing directory as empty history (`:65-67`) and resolves a relative reports path against `ctx.RepositoryRoot` (`:59-62`). The new walk follows both conventions and leaves this file unchanged (D-08).
- `portfolio.go:141-154` is the precedent for D-02: it checks that the existing target is a regular file and not a link, calls `rootedPublishFile(root, rel, data, 0o644, existed)`, and records a `resultmodel.RecordedChange{Kind: "created"|"modified"}`. No git transaction.
- `mutation.go` provides `validateNoLinkedAncestors`, `rootedPublishFile` (exclusive create, or inode-and-bytes-checked replace) and `rootedReadFile`. `runTransaction` (`mutation.go:54`) is not used: `gittransaction.ExecuteTransaction` (`gittransaction/git_transaction.go:448`) refuses dirty targets, which would block a second `index` run before the first catalog is committed.
- Test helpers: `report_image_test.go:67` (`toolboxTestRepository`, a git-initialized temp dir). The new verb needs no git, so a plain `t.TempDir()` with an `ai-reports/` fixture is enough. `architecture_test.go:21-45` shows the fixture style.
- The resultmodel refusal shape: `resultmodel.OutcomeRefused`. A refusal finding whose `next_argv` names its own verb has that argv cleared (`resultmodel/result_model.go:946-951`). Fixability manual is correct for "move the hand-made file aside", since no verb resolves it.

**Real-world sample.** This repo's own `ai-reports/` holds 41 bundles plus one loose `.html`. Most are `yyyy-mm-dd_hhmm_slug`. `2026-08-26_1709_architecture-report/` holds only `architecture-report.md` (the reason for D-03). `2026-08-27_1428_req-341-timeline-drag-evidence` writes the id in lowercase (the reason ids are matched case-insensitively). `2026-07-29_1937_UR-007-008-deep-review-batch` shows a name with two numbers. The builder may run the verb against a scratch copy of this folder as a smoke check, never against the worktree's own `ai-reports/`.

**Prose (paths under `skills/do-work-toolbox/`)**
- `actions/ai-report.md:28-30` (Input), `:109` (Step 6 "does not publish, host, or search"), `:125` (Output Format) and `:140` (checklist "No ... search artifact"). Narrow `:109` and `:140` for the catalog forms only.
- `docs/ai-report-guide.md:63-70` (Input block) and `:74-78` (Evidence Safety, the last section).
- `SKILL.md:23` (the ai-report routing row) and `:3` (`argument-hint`). The hint lists actions, not forms, so it is unchanged.
- `actions/help.md:13` (the ai-report help line).
- `_dev/tests/staged-skills-contract.sh` pins the retired-trigger counts (`:432`, `("do-work-toolbox", "ai-report"): 7`) from `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv`. These are retired core triggers, not the toolbox routing row, so new routing phrases need no fixture change. No contract test sizes the toolbox router.
- `skills/do-work/docs/command-line-guide.md:26` lists toolbox Just recipes. It is unchanged because there is no new recipe (D-09).

**Merge seams with siblings in this run**
- REQ-657 (ai-report judge render check) may add its own toolbox CLI verb. If it does, it edits the same `commands.go` constant block, the `Handlers()` map and the `commands_test.go` name list and count. Both builders should add their lines next to their own family, and the second integrator sets the count to the sum.
- REQ-687 (ai-report proposal and root-cause kinds) and REQ-657 both edit `actions/ai-report.md` (Input `:30`, Step 1, Step 6), `docs/ai-report-guide.md` (`:47`, Input block), `SKILL.md:23` and `actions/help.md:13`. REQ-655 keeps its edits local:
  - its own short section right after `## Input`;
  - two new lines at the end of the guide's Input block;
  - a new guide section after Evidence Safety;
  - a new help line under `:13` rather than a rewrite of `:13`;
  - phrases appended at the end of the `SKILL.md:23` trigger cell.

*Explored by the pre-dispatch agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (new): the walk, the catalog, the writer and find for the ai-report-index verb.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (new): four focused fixture tests.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` (modify): constant and handler registration.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (modify): registration list and count.
- `skills/do-work-toolbox/actions/ai-report.md` (modify): index and find forms; narrow the no-search sentence and the checklist item.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify): Input lines and a Report catalog section.
- `skills/do-work-toolbox/SKILL.md` (modify): routing phrases on the ai-report row.
- `skills/do-work-toolbox/actions/help.md` (modify): one help line.

**Files I will NOT touch:** `architecture.go` and its tests (D-08), `justfile`, `skills/do-work-board/justfile.template`, `skills/do-work/docs/command-line-guide.md` (D-09), `actions/ai-report-reference.md`, `actions/completed-work-presentation-reference.md`, any bundle under `ai-reports/`, and the release files (`CHANGELOG.md`, its mirror and version files belong to the integrator).

**Acceptance criteria (restated from REQ):**
- [ ] A new toolbox CLI verb walks `ai-reports/` and treats every directory under it as a bundle whatever its naming style (at least `yyyy-mm-dd_hhmm_slug`, `yyyy-mm-dd_slug`, `yyyy-mm-dd-slug`, `REQ-NNNN-slug`, `slug-yyyy-mm-dd`, `yyyy-mm-dd_hhmmss_slug`), with no closed list of styles.
- [ ] Entry file pick order: `index.html`, `index.md`, `README.md`, `prompt.md`, `report.md`, then the first `.html`.
- [ ] Per bundle it reads the title or first heading, the `ai-report-kind` meta, any UR/REQ id in the folder name or title, and the `ai-report-supersedes` meta.
- [ ] The catalog records path, title, date, kind, linked UR/REQ, verdict or decision, `supersedes` and `superseded_by`, with `superseded_by` computed from other bundles' `supersedes`.
- [ ] `ai-reports/catalog.json` and a static `ai-reports/index.html` are written, grouped by kind and date, with superseded reports greyed and linked to their successor.
- [ ] Loose files directly under `ai-reports/` are skipped.
- [ ] `find <topic>` is a substring plus UR/REQ id match that prints paths newest first, including superseded bundles.
- [ ] No existing bundle file changes. A hand-made `catalog.json` or `index.html` without the generator marker is refused, not overwritten.
- [ ] Both forms are wired into `ai-report.md` (the `:109` sentence relaxed for this form only), `docs/ai-report-guide.md`, the `SKILL.md:23` routing phrases and the `help.md` help line. Focused Go tests run on a fixture `ai-reports/` folder.
