---
id: REQ-655
title: 'ai-report index and find build a derived catalog of every report bundle in any naming style'
status: completed
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
builder_handback_at: 2026-10-10T13:31:11Z
integration_at: 2026-10-10T15:29:59Z
review_at: 2026-10-10T15:40:38Z
kb_status: pending
commit: d91312aeaefb3600cf6fc3c7f4fa156451ebaf17
heavy_verified_at: 2026-10-10T15:49:00Z
heavy_verified_revision: d91312aeaefb3600cf6fc3c7f4fa156451ebaf17
completed_at: 2026-10-10T15:49:44Z
release_at: 2026-10-10T15:49:44Z
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
- [x] **[PLAN]:** (from the builder hand-back) I followed the REQ `## Plan` and the brief step by step. I read all of the REQ, the brief, UR-144 A2 (as data), the crew rules, the action-files, releases and shell primes, `lessons-releases.md`, and the three named lesson families. The tests were written first. Then I built one Go file holding a walk shared by both forms, then the prose. The deviations are D-11 to D-22 below.
- [x] **[APPLY]:** (from the builder hand-back) Code was written as planned, and only the eight files in Scope changed. `architecture.go` was not changed (D-08). There is no `--dry-run`, `--commit`, `--reports` or Just recipe. No `runTransaction`. No "may be stale" mark.
- [x] **[UNIFY]:** (from the builder hand-back) `git diff bd56c4b0 --stat`:
  ```
   skills/do-work-toolbox/SKILL.md                    |   2 +-
   skills/do-work-toolbox/actions/ai-report.md        |  14 +-
   skills/do-work-toolbox/actions/help.md             |   1 +
   skills/do-work-toolbox/docs/ai-report-guide.md     |  10 +
   .../internal/toolboxcommands/commands.go           |   2 +
   .../internal/toolboxcommands/commands_test.go      |  11 +-
   .../internal/toolboxcommands/report_index.go       | 418 +++++++++++++++++++++
   .../internal/toolboxcommands/report_index_test.go  | 245 ++++++++++++
   8 files changed, 697 insertions(+), 6 deletions(-)
  ```
  Checks (all run from the worktree root, on a loaded machine):

  | Check | Exit | Wall time |
  |---|---|---|
  | `bash .../REQ-655-probe.sh` (gofmt, vet, five named tests with `--- PASS:`) | 0 | about 1 s (warm cache) |
  | `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/toolboxcommands/` | 0 (`ok ... 4.461s`) | 5 s |
  | `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands/` | 0 | under 1 s |
  | `bash _dev/tests/shipped-package-reference-contract.sh` | 0 (`PASS`) | 1 s |
  | `git diff bd56c4b0 --check` | 0, clean | n/a |

  Every check passed on its first run, so nothing needed a rerun. I checked each changed file in the diff: the Go is gofmt-clean, vets, and has no debug output. The prose edits stay local, with no reflow of neighbouring lines. The help and guide columns line up with their neighbours (description at column 34 in help, column 43 in the guide).
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

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work-toolbox/SKILL.md` (modified)
- `skills/do-work-toolbox/actions/help.md` (modified)

**What was done:** Added the toolbox CLI verb `ai-report-index`. It treats every non-dot directory under `ai-reports/` as a bundle (loose files and links skipped), picks the entry file (`index.html`, `index.md`, `README.md`, `prompt.md`, `report.md`, first `.html`, first `.md`), reads title, date, kind, linked UR/REQ ids, and the `ai-report-kind`, `ai-report-supersedes` and `ai-report-verdict` metas, computes `superseded_by`, and publishes a deterministic `ai-reports/catalog.json` and a static `ai-reports/index.html` (grouped by kind, superseded rows greyed and linked to the successor). A file at either output without the generator marker is refused with `AI-REPORT-INDEX-HAND-MADE` and nothing is written. `--find <topic>` walks the same bundles in memory and prints matching paths newest first, superseded ones included, and writes nothing. The verb is registered in `commands.go` (toolbox handler count 7 to 8), and the ai-report action, guide, toolbox router row and help list the `index` and `find` forms. The hand-back merge `355ea12f..6a031809` needed no conflict resolution and no seam edits.

After review (builder-branch commit `851fa107` by the integrator, re-merged with the same `<pre>` as `d91312ae`): `--find` matches the bundle folder name (`path.Base`) instead of the full path, because every path starts with `ai-reports/` and topics such as `ai` or `reports` matched every bundle (review F1). The find test gained one assertion that `--find reports` prints the no-match line.

## Decisions

(from the builder hand-back, verbatim; D-01 to D-10 are in `## Plan`)

- **D-11** DECIDE & STATE: the action prose runs the verb with `--format text`, not the brief's `--format json`. The reason is that `ExactTextOutput` is tagged `json:"-"` (`resultmodel/result_model.go:640`). With JSON, the find matches and the bundle count would not appear in the output at all. Text prints the exact lines on success. On a refusal it prints `ai-report-index: refused` plus the finding and exits 1. This needs no extra Go. If the run wants JSON, the alternative is an info finding that carries the lines, the way `ARCHITECTURE-SCAN` does.
- **D-12** DECIDE & STATE: the linked-id regex is `(?i)(?:^|[^A-Za-z0-9])((?:UR|REQ)-\d+)`, not `\b...\b`. Go's `\b` counts `_` as a word character, so the brief's regex misses `2026-08-27_1428_req-341-...`, a real bundle in this repo. The fixture `2026-09-02_req-77-beta-notes` pins this. No trailing boundary is needed because `\d+` is greedy.
- **D-13** DECIDE & STATE: kind is matched over the kebab words of the whole folder name, not only "the slug after the date". Kind words contain no digits, so dated names give the same result. Names in the `slug-yyyy-mm-dd` style keep their slug, which comes before the date.
- **D-14** DECIDE & STATE: meta tags are read only from HTML entries. A Markdown entry that quotes `<meta name="ai-report-kind" ...>` in prose or code would otherwise set the kind falsely.
- **D-15** DECIDE & STATE: the help line sits under `ai-report [REQ|UR]`, which is `help.md:14` at base. The brief's `:13` is `journey-qa` at base. The REQ describes `:13` as the ai-report line, so the line numbers drifted by one.
- **D-16** DECIDE & STATE: the hand-made check is a byte-substring test for the generator marker in each existing file. For JSON the marker is `"generator": "do-work-cli ai-report-index"`. An existing path at either output that is not a regular file (a directory or a link) also counts as "not written by this verb" and gets the same refusal code.
- **D-17** DECIDE & STATE: one failure code, `AI-REPORT-INDEX-FAILED`, covers an unreadable `ai-reports/` and a publish error. The brief named no code for these, and both are manual fixes.
- **D-18** DECIDE & STATE: `--find` with a blank topic returns usage. Otherwise it would match every bundle.
- **D-19** DECIDE & STATE: on the page, kind groups are sorted alphabetically with `unknown` last.
- **D-20** DECIDE & STATE: page hrefs go through `(&url.URL{Path: ...}).String()` before `html.EscapeString`. Without this, a folder name with a space, `#` or `:` would produce a broken link.
- **D-21** DECIDE & STATE: the catalog uses a JSON encoder with `SetEscapeHTML(false)`, so `&` and `<` in titles stay readable. Output is still deterministic, and the byte-identical rerun is pinned.
- **D-22** DECIDE & STATE: if a bundle folder or its entry file cannot be read, the bundle is still listed, with an empty entry or the folder-name title. A lookup walk never drops a record and never fails because of one (lesson family `silent-skip-reads-as-red`).
- **D-23** (integrator, after review, DECIDE & STATE): review F1 showed that `find` matched the shared `ai-reports/` prefix, so `find ai` and `find report` returned all 41 bundles of this repo. The integrator changed the matched field from the full path to the folder name and pinned it in `TestAIReportIndexFindListsMatchesNewestFirstIncludingSuperseded` (fails at `report_index_test.go:206` without the fix, per the delta re-review). Side effect, accepted: a topic written as a full `ai-reports/<folder>` path no longer matches (review F10, report only); the bare folder name still does. The plan's task 3 wording "substring over path" now reads "over the folder name".

## Discovered Tasks

(from the builder hand-back; impact tokens are the builder's)

- impact-low: other toolbox verbs whose prose uses `--format json` may return only `exactOutputResult`. If so, the JSON caller sees no payload, the trap D-11 avoided here. Audit the `exactOutputResult` callers against their action prose. → report only
- impact-low: this repo's own `ai-reports/2026-05-28_2335_background-agents-durability.html` is a loose file, so `ai-report index` and `find` skip it by design. If it is a real report, it is invisible to the catalog until someone moves it into a bundle folder. → report only

## Qualification

**Gate records (`advance --diff-range 355ea12f..6a031809`):** `qualify` satisfied (success), `scope-drift` satisfied (success). One `QUALIFY-NEW-FILE-UNWIRED` warning on `report_index_test.go`, judged a false positive: a `_test.go` file is found by the Go test runner by convention. `report_index.go` was not flagged because `commands.go` registers `handleAIReportIndex`. No debug-artifact, P-A-U or output-primitive findings.

**Scope:** the declared `write_set` (8 paths) equals the touched set in `git diff 355ea12f..6a031809 --stat` (8 files, 697 insertions, 6 deletions). The queue guard printed nothing before the merge. The merge auto-resolved against REQ-658, REQ-659 and REQ-660 with no conflict (none of them touch these files), and the merged tree builds and vets (`go build ./...`, `go vet ./internal/toolboxcommands/`).

**Requirement trace (read against the merged files):**
1. Any directory is a bundle: `reportIndexWalk` (`report_index.go:81`) lists every non-dot directory under `ai-reports/`. There is no list of naming styles; the date is read by one condition (`reportDatePattern`, `yyyy-mm-dd` plus an optional `_hhmm`/`_hhmmss`). The fixture covers the six named styles.
2. Entry pick order: `reportEntryNames` (`:31`) then the first `.html`, then the first `.md` (pre-dispatch D-03), regular files only (`pickReportEntry`, `:167`).
3. Per bundle: `<title>` else the first `<h1>`-`<h6>` for HTML, the first `#` heading for Markdown, else the folder name; the three metas from HTML entries only (D-14); linked ids from folder name plus title with the explicit non-alphanumeric lead (D-12).
4. Catalog fields: `reportBundle` (`:55`) holds path, entry, title, date, kind, linked_ids, verdict, supersedes, superseded_by. `superseded_by` is computed in the walk from the other bundles' `supersedes`, matched by exact last path segment; with several successors the newest wins (bundles are sorted newest first).
5. Outputs: `reportIndexWrite` publishes both files with `rootedPublishFile`; the page groups by kind (alphabetical, `unknown` last, D-19), rows newest first, superseded rows get `class="superseded"` (opacity .5) and a link to the successor's entry.
6. Loose files skipped: only `entry.IsDir()` entries are bundles; a symbolic link reports as a link, not a directory.
7. Find: case-insensitive substring over path, title, kind, verdict and linked ids, plus an id compare with leading zeros ignored (`normalizeReportID`); output order is the walk's newest-first order and superseded rows are kept with a `(superseded by ...)` suffix.
8. No bundle changes: only the two root outputs are written. A file at either output without the generator marker (or not a regular file, D-16) is refused with `AI-REPORT-INDEX-HAND-MADE` before any write. Output carries no timestamp, so a rerun is byte-identical (pinned in the refusal test).
9. Prose: `ai-report.md` gains `### Catalog forms: index and find` after Input; the Step 6 no-search sentence and the checklist item name the catalog forms as the one exception. The guide gains two Input lines and `## Report Catalog`. `SKILL.md:24` gains `report index`, `is there a report on` at the end of the ai-report cell. `help.md` gains one line under `ai-report [REQ|UR]`.
10. Tests and release: four focused fixture tests plus the updated registration test; the release is made at finalization.

**Builder deviations judged:** D-11 holds: `ExactTextOutput` is tagged `json:"-"` (`resultmodel/result_model.go:672`) and `exactOutputResult` sets only that field, so `--format json` would drop the find lines and the bundle count; the action's `--format text` call is right. D-12 holds: Go RE2 `\b` treats `_` as a word character, so `\bREQ-` misses `_req-341`; the fixture `2026-09-02_req-77-beta-notes` pins it. D-13 to D-22 are local and reversible, consistent with the code.

**Constraints:** no new queue field or status; no bundle renamed; `architecture.go` unchanged (D-08); no Just recipe (D-09); own release.

**Review-fix re-merge, cumulative range `355ea12f..d91312ae`:** `advance` takes no qualify input at this phase, so the integrator ran the same `qualify --request-path <P> --diff-range 355ea12f..d91312ae` handler directly: success, only the judged `QUALIFY-NEW-FILE-UNWIRED` warning on `report_index_test.go`. The touched set is still the 8 declared files (703 insertions, 6 deletions); the delta touches `report_index.go` and its test only. The queue guard before the re-merge printed nothing. Requirement 7 now reads: the substring match runs over the folder name, title, kind, verdict and linked ids (D-23).

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge `6a031809` (machine quiet before launch: 1-minute load 3.07, no other gate running; load 6.14 at the end), then `advance REQ-655 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-131527/REQ-655-probe.sh`.
**Result:** ✓ Repository gate passed on the first run (exit 0, gate wall 166 s; do-work-cli 901 tests in 78 s, slowest file `internal/finalization/finalization_recovery_test.go` 25.58 s under the 30 s limit; queue-kanban 421 tests in 55 s, slowest 20.90 s). Probe exit 0 (gofmt, vet, five named tests each `--- PASS:`). Gate records `green-gate`, `scope-drift` and `run-blocked-check` satisfied.

**Red-green validation:** (from the builder hand-back, traced to `## Red-Green Proof`)
- `TestAIReportIndexCatalogsEveryBundleNamingStyleOnce`, `TestAIReportIndexLinksSupersededProposalToSuccessor`, `TestAIReportIndexFindListsMatchesNewestFirstIncludingSuperseded`, `TestAIReportIndexRefusesHandMadeCatalog` (`internal/toolboxcommands/report_index_test.go`): ✗ before the production code (build failure: `undefined: CommandAIReportIndex`, `undefined: reportCatalog`, `undefined: reportIndexGenerator`, `undefined: reportBundle`; at base the launcher answered `ai-report-index` with `UNKNOWN-COMMAND`) → ✓ after (0.03 s, 0.03 s, 0.01 s, 0.06 s). The fixture is the captured RED case: one bundle per naming style, a `README.md`-only bundle, a proposal superseded by a later one, and loose `.patch` and `.md` files; `find proposal` prints both proposals newest first.
- `TestHandlersRegisterCanonicalToolboxCommands`: updated for the new verb (count 8), passes.

**New tests added:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (four tests on one shared fixture)

**Existing tests updated:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (registration list and count 7 to 8)

**Heavy verification plan:**
- Range: 355ea12f6c4485c8fd3de1bfae65062bbe64d0cc..d91312aeaefb3600cf6fc3c7f4fa156451ebaf17 (re-planned after the review-fix re-merge; same four lanes as the first plan at `6a031809`)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — changed files under `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — changed files under `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — changed files under `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — changed files under `skills/do-work/tools/do-work-cli`

**After the review-fix re-merge (`d91312ae`):** the same gate argv ran again (load 3.85 before launch, 8.08 at the end): exit 0, gate wall 160 s, do-work-cli 901 tests in 74 s, slowest file `internal/finalization/finalization_recovery_test.go` 23.12 s. `advance` refuses gate input past this phase, so the green record was written with `record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh` at `d91312ae`, and `REQ-655-probe.sh` was run directly: exit 0. Red-green for the fix: the new `--find reports` assertion fails with the old full-path match (`report_index_test.go:206`) and passes with the fix. The heavy plan was recomputed for `355ea12f..d91312ae`: the same four lanes.

*Verified by work action*

## Review

**Overall: 91%** | 2026-10-10T15:40:38Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 85% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1 `report_index.go:386`: the find search text includes the `ai-reports/` prefix, so `find ai` and `find report` match all 41 bundles. Match `path.Base(bundle.Path)` instead. impact-user-visible → report only
- F2 `ai-report.md:30`, `:135` and `ai-report-guide.md:21`, `:84`: neighbouring input, output and "only the bundle" lines are stale against the new index/find forms (Restatement Sweep). impact-negligible → report only
- F3 `report_index_test.go:211`: a hand-made `index.html` and a non-regular output (D-16) are not pinned. Both were verified refused by hand. impact-negligible → report only
- F4 `report_index.go:43`, `:145`: the Markdown title can come from a `#` line inside a fenced code block. impact-negligible → report only
- F5 (nit) `report_index.go:44-45`: meta parsing misses `data-name=`, a `>` inside a quoted value, and unquoted attributes. impact-negligible → report only
- F6 (nit) `report_index.go:222-223`: the kind meta is not lowercased, so the page splits its groups. impact-negligible → report only
- F7 (nit) `report_index.go:157-163`: `REQ-0412` and `REQ-412` are not deduped, and `UR-007-008` links only `UR-007`. impact-negligible → report only
- F8 (nit) `report_index.go:236-239`: a symlinked `ai-reports/` makes index fail but find follow the link, and the failure suggests a rerun. impact-negligible → report only
- F9 (nit) anti-bloat: 1 option form, 2 finding codes, 1 extra field (`entry`), 1 invented meta (`ai-report-verdict`, which no writer emits), 1 kind table, 4 constants, 13 functions, 3 test helpers, 0 decorative tests, and cosmetic `commands_test.go` list churn. All are plan-traced except `entry`, which the hand-back omits. impact-negligible → report only

**Acceptance:** Pass. Implementation and integration stages: the package tests pass, and the CLI smoke on a temp copy of `ai-reports/` gave 41 bundles, only two new files, correct find order, and hand-made, symlink and absent-folder refusals as specified.
**Restatement sweep:** redefined the ai-report action's input forms and its "only the bundle, no search" boundary. Stale neighbours found in `ai-report.md:30`, `:135` and `ai-report-guide.md:21` (F2). The bundle naming `ai-reports/yyyy-mm-dd_hhmm_<slug>/` was not redefined (the verb reads all styles and writes no bundle). The entry-file pick order is new and not restated in shipped prose. `architectureScan` stays `index.html`-only by D-08. `command-line-guide.md:26` lists Just recipes, and none was added (D-09), so it is not stale.
**Suggested testing:** 4 items
**Follow-ups created:** None (9 findings report only)

*Reviewed by review-work action*

**Delta re-review** | 2026-10-10T15:48:47Z | delta `6a031809..d91312ae` (builder-branch commit `851fa107`)

**Overall: 92%** (Code Quality 85% → 90%; others unchanged). Acceptance: Pass. F1 closed: `find` matches the folder name; the new assertion fails without the fix; `find ai` and `find report` dropped from 41/41 to 5 and 11 on a copy of this repo's `ai-reports/`, and `find architecture` still finds 7.
- F10 (nit) `report_index.go:388`: a topic written as a full `ai-reports/<folder>` path no longer matches its own bundle; the bare folder name does. impact-negligible → report only

**Follow-ups created:** None (10 findings report only)

## Lessons Learned

**What worked:** One in-memory walk shared by `index` and `find` kept the two forms consistent and let `find` stay read-only. Running the verb on a scratch copy of this repo's real `ai-reports/` (41 bundles) found the underscore-joined `req-341` id that a `\b` regex misses (D-12).
**What didn't:** Matching the topic against the full stored path: every bundle path shares the `ai-reports/` prefix, so short topics matched everything (review F1). The fixture tests did not catch it because no test searched for a word inside the shared prefix.
**Worth knowing:** Toolbox verbs that return `exactOutputResult` print nothing useful under `--format json` (`ExactTextOutput` is `json:"-"`), so their action prose must call them with `--format text` (D-11). Go RE2 `\b` counts `_` as a word character. The catalog field names `path`, `date`, `supersedes` and `superseded_by` are read by REQ-656 (ai-report revise).

## Orientation

Now `ai-report index` writes a derived catalog (`ai-reports/catalog.json` plus a static `index.html`) of every report bundle in any naming style, and `ai-report find <topic>` lists matching bundles newest first; the walk lives in the toolbox CLI (`toolboxcommands/report_index.go`), the prose in `skills/do-work-toolbox/actions/ai-report.md`. [MAP CHANGED]: a new toolbox CLI verb and a new derived catalog contract that REQ-656 consumes. Touched primes (`prime-action-files.md`, `prime-releases.md`) name no path this change removed.

## Heavy Verification Plan

- Base: 355ea12f6c4485c8fd3de1bfae65062bbe64d0cc
- Target: d91312aeaefb3600cf6fc3c7f4fa156451ebaf17
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — changed files under `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — changed files under `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — changed files under `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — changed files under `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: d91312aeaefb3600cf6fc3c7f4fa156451ebaf17
- Execution revision: d91312aeaefb3600cf6fc3c7f4fa156451ebaf17 (detached checkout under `.git/work-run-work-2026-10-10-131527/`, `QUEUE_KANBAN_BROWSER` set)
- do-work-cli-integrations: exit 0, executed, 79 s
- staged-skills: exit 0, executed, 38 s
- updater: exit 0, executed, 68 s
- installer: exit 0, executed, 32 s
- The earlier drain at `6a031809` was also green (78 s, 42 s, 71 s, 33 s, all executed).

## Timing

Observed 2026-10-10T15:29:30Z to 2026-10-10T15:48:18Z: 18m 48s total, 20m 37s attributed across 7 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 13m 50s | 4 |
| review | 6m 00s | 1 |
| handback-merge | 47s | 2 |

Slowest stage: review / review, 6m 00s, outcome success.
Slowest command: verification-gate / heavy drain, 3m 57s, exit 0, .

Notes: the builder-work event was skipped because the hand-back had already landed when integration began (`actions/fan-out-reference.md` → Landed hand-back); `builder_handback_at` is the builder commit's committer date.
