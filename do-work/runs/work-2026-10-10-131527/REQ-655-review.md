## Review: REQ-655

**Approve**: the `ai-report-index` verb (a CLI command that lists every report folder) builds the catalog and runs `find` as the REQ asks. It is deterministic and refuses hand-made files. One small `find` bug is worth a one-line fix.
Route C | merge range `355ea12f..6a031809` (builder commits `9170e02c`, `e56dc378`)

### What's built
- `ai-report index` treats every non-dot folder under `ai-reports/` as a report bundle, whatever its name. It writes `ai-reports/catalog.json` and a static `ai-reports/index.html`, grouped by kind. Superseded rows are greyed and linked to the report that replaced them. A file at either path that the verb did not write is refused and nothing is written.
- `ai-report find <topic>` walks the same bundles in memory and prints paths newest first, superseded ones included. It writes nothing.
- The ai-report action, its guide, the toolbox router row and the help line list both forms. The release half of requirement 10 belongs to finalization and is not yet done.

### Decisions / risks for you
- F1 (below): `find ai` and `find report` match every bundle, because the search text includes the `ai-reports/` path prefix. It is a one-line fix (match the folder name, not the full path). It is cheap to do before release, but it does not block.

### Findings

**Important:**
- None

**Minor:**
- F1 `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go:386`: the `find` search text is the full `bundle.Path`, so any topic that is a substring of `ai-reports/` (`ai`, `report`, `reports`, `port`) matches every bundle. Measured on a temp copy of this repo's `ai-reports/`: `--find ai` 41/41, `--find report` 41/41. The fix is to match `path.Base(bundle.Path)` instead. Recommended before release, not blocking. impact-user-visible → report only
- F2 `skills/do-work-toolbox/actions/ai-report.md:30`, `:135`, and `skills/do-work-toolbox/docs/ai-report-guide.md:21`, `:84` (Restatement Sweep): neighbouring text still describes the action without the new forms. These lines are now stale: `:30` still says "`$ARGUMENTS` is `UR-NNN`, `REQ-NNN`, `most recent`, or blank"; `:135` "The timestamped folder is the only stakeholder artifact this action publishes"; the guide `:21` "This action produces only the report bundle". The guide `:84` "lists every directory" leaves out the dot-folder and symlink skips. The new subsection directly below `:30` limits confusion. REQ-657 and REQ-687 edit these files next and could absorb this. impact-negligible → report only
- F3 `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go:211`: the refusal test pins only a hand-made `catalog.json`. A hand-made `ai-reports/index.html` (the case the source report saw in a consumer repo) is not pinned. Neither is the non-regular-file refusal (D-16). I ran both by hand on a temp copy: a hand-made `index.html`, a symlinked `catalog.json` pointing outside, and a symlink holding the marker were all refused (exit 1), and the outside file stayed unchanged. impact-negligible → report only
- F4 `report_index.go:43` and `:145`: the Markdown title is the first `#` line anywhere, including inside a fenced code block. Measured: a README that opens with a fenced block holding `# not a title` got the title "not a title". impact-negligible → report only

**Nit:**
- F5 `report_index.go:44-45`, meta parsing edges. `\b(name|content)` also matches inside `data-name=`. `<meta\b[^>]*>` cuts a quoted value at `>`. Unquoted attributes (`name=ai-report-kind content=root-cause`) are ignored. Measured: kind or verdict was lost in each case. No writer in this suite emits these shapes. impact-negligible → report only
- F6 `report_index.go:222-223`: the `ai-report-kind` meta is not lowercased, so `Proposal` gets its own page group apart from `proposal`. `find` still matches it. impact-negligible → report only
- F7 `report_index.go:157-163`, linked-id edges. `REQ-0412` and `REQ-412` are stored as two ids. `UR-007-008` links only `UR-007`, so `find UR-008` misses `2026-07-29_1937_UR-007-008-deep-review-batch`. impact-negligible → report only
- F8 `report_index.go:236-239`: when `ai-reports/` is itself a symlink, `index` fails with `AI-REPORT-INDEX-FAILED` (correct, the publish guard refuses it), but `find` follows the link. The failure finding also prints `next: do-work-cli ai-report-index`, and a rerun cannot fix a symlink. impact-negligible → report only
- F9, anti-bloat count (the maintainer asked for this). Items the REQ text did not name: 1 option form (`--find`, plan D-10). 2 finding codes (`AI-REPORT-INDEX-HAND-MADE`, which the REQ asked for only as "a typed finding", and `AI-REPORT-INDEX-FAILED`, D-17). 1 extra catalog field (`entry`). 1 invented meta name (`ai-report-verdict`, D-06: no writer emits it, so `verdict` is empty for every bundle today). 1 known-kind table (D-05). 4 constants. 13 functions. 3 test helpers. No new files beyond the two the plan named. No decorative tests: each of the 4 tests names the failure it pins. There is also cosmetic churn in `commands_test.go:7-12`: three names went one per line and five stayed on one line, a merge-seam choice. Everything above matches the hand-back's Anti-bloat check except the `entry` field, which the hand-back does not list. None is unjustified. impact-negligible → report only

### Requirements Checklist

- [x] R1 Any folder under `ai-reports/` is a bundle, with no closed list of naming styles. Delivered by the `entry.IsDir()` condition (`:90`) and one date condition (`:38`). The fixture covers all six named styles.
- [x] R2 Entry pick order `index.html`, `index.md`, `README.md`, `prompt.md`, `report.md`, then the first `.html`. Delivered (`:31`, `:167-197`), plus a first-`.md` fallback (D-03).
- [x] R3 Title or first heading, `ai-report-kind`, UR/REQ ids from name and title, `ai-report-supersedes`. Delivered. Metas are read from HTML entries only (D-14). See F4 and F5 for edges.
- [x] R4 Fields path, title, date, kind, linked ids, verdict, `supersedes`, `superseded_by`, with `superseded_by` computed. Delivered (`:55-65`, `:101-115`). When several bundles supersede the same one, the newest wins. A self-reference is ignored.
- [x] R5 `catalog.json` and static `index.html`, grouped by kind and date, superseded rows greyed and linked to the successor. Delivered (`:316-373`).
- [x] R6 Loose files skipped. Delivered. Symlinks are skipped too.
- [x] R7 `find` is a substring plus id match, newest first, superseded included. Delivered (`:376-405`). See F1 for the over-match.
- [x] R8 Derived and regenerated, with no bundle file changed. Delivered. Pinned by the byte snapshot (`report_index_test.go:102-120`) and by the byte-identical rerun (`:235-244`).
- [x] R9 Prose wiring in `ai-report.md` (the no-search sentence narrowed), the guide, `SKILL.md` routing phrases and the `help.md` line. Delivered. See F2 for stale neighbours.
- [x] R10 Focused Go tests on a fixture folder. Delivered. The release half is N/A at review: it happens at finalization.
- [x] Red-Green GREEN: every bundle appears once, the loose file is skipped, `superseded_by` is set, `index.html` is written, no other file changes, and `find proposal` prints both proposals newest first. All of this is pinned by the four tests.
- [x] UR-144 Acceptance check (A2 lines): the `index` and `find proposal` lines are met by the same tests and by the smoke run.
- Builder deviation D-11 confirmed. `ExactTextOutput *string \`json:"-"\`` is at `resultmodel/result_model.go:672` (the hand-back cites `:640`, a stale line number). A `--format json ... --find x` run on a temp copy carried no match lines, so the prose's `--format text` call is right.
- Builder deviation D-12 confirmed. Go RE2 `\b` treats `_` as a word character. The real bundle `2026-08-27_1428_req-341-timeline-drag-evidence` got `["REQ-341","REQ-375"]`, and the fixture `2026-09-02_req-77-beta-notes` pins it. `sur-12` is correctly not linked.

### Acceptance Testing

**Result: Pass.** This covers the Implementation and Integration stages.
- `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/toolboxcommands/`: `ok` in 4.37 s. `go vet` is clean, and `gofmt -l` lists nothing.
- Smoke on a `mktemp -d` copy of this repo's `ai-reports/` (temp dir removed afterwards): `--format text ai-report-index` printed `ai-reports/catalog.json: 41 bundles` and exited 0 in 0.15 s. `diff -rq` against the source showed only `catalog.json` and `index.html` added. The loose `.html` was skipped. Kinds: 32 unknown, 7 architecture-report, 2 proposal. Every bundle had an entry. `--find architecture` printed 7 rows newest first, and `--find REQ-0588` printed the 3 REQ-588 bundles.
- Edge fixture (temp): a folder named `weird name #1:x` got href `./2026-10-01_1200_weird%20name%20%231:x/index.html`. The title `<script>` text was escaped on the page. The colon-first folder `10:30-...` got a `./` prefix. An unreadable bundle (mode 000) was listed with an empty entry and no failure (D-22). A 7-digit time fell back to date only. A mutual supersedes pair greyed both rows. An absent `ai-reports/` returned exit 0 with "nothing to index". `--find` with no topic and an unknown argument both returned usage (exit 2).
- The edge fixture also showed F1 and F4 to F8.

### Suggested Additional Testing

- Deployment: unassessed. The release (CHANGELOG, mirror, version) happens at finalization. After it, run `ai-report-index` from an installed consumer copy.
- Live acceptance: unassessed. In the consumer repo the source report describes (about 180 entries, two hand-made catalogs), run `ai-report index` through the toolbox action. Confirm the refusal message is clear enough to act on, and that `find` answers a real "is there a report on X".
- Manual: open the generated `index.html` in light and dark at 1440 px and 390 px. The nowrap date and id columns may push the width at phone size.
- Edge: run `index` while another session publishes a new bundle, and confirm the temp `.catalog.json.publishing-*` file never shows up as a bundle. It cannot, because it is a loose dot-file, but it is unobserved.

### Scores (on the record, not the headline)

**Overall: 91%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All 10 delivered. Release half is N/A until finalization. |
| Code Quality | 85% | Clean, deterministic, confined writes. F1 over-match. Parse edges F4 and F5. |
| Test Adequacy | 85% | Four focused tests pin the RED case, with red-green evidence. Hand-made `index.html` and non-regular refusal unpinned (F3). |
| Scope | 95% | Touched set equals the 8-path write_set. Cosmetic test-list churn. |
| Risk | Low | Writes confined to two root files by `rootedPublishFile`. No git transaction (D-02). Read-only find. |
| Acceptance | Pass | Implementation and integration stages: unit tests plus CLI smoke on a temp copy. |

### Follow-ups created
None (9 findings report only)

## Review

**Overall: 91%** | 2026-10-10T15:39:15Z

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

## Delta re-review

**Approve. F1 is closed.** | 2026-10-10T15:41:40Z | delta `6a031809..d91312ae` (builder commit `851fa107`, merge `d91312ae`)

**Overall: 92%.** Code Quality goes from 85% to 90% now that F1 is fixed. The other dimensions are unchanged: Requirements 100%, Test Adequacy 85%, Scope 95%, Risk Low, Acceptance Pass.

- **F1 closed.** `report_index.go:388` now matches `path.Base(bundle.Path)`. The delta touches only `report_index.go` (+3/-1, including a short why-comment) and `report_index_test.go` (+4).
- **The new assertion fails without the fix.** I extracted `d91312ae` with `git archive` into a temp dir and ran the tests there. `TestAIReportIndex*` passed. I then reverted the one line to `bundle.Path`, and `TestAIReportIndexFindListsMatchesNewestFirstIncludingSuperseded` failed at `report_index_test.go:206` ("find on the shared path prefix"). The temp dir was removed. The full package passes on main (`ok` 4.55 s).
- **Smoke on a temp copy of `ai-reports/`.** Before the fix, `--find ai` and `--find report` each matched 41/41. Now they match 5 and 11, and every match comes from the folder name, title or kind. `--find reports` matches 0, and `--find architecture` still matches 7.
- **Regression check.** A full-path topic such as `ai-reports/2026-08-26_1709_architecture-report` now matches 0 bundles. Before the fix it matched 1. The bare folder name still matches 1. This matters little: `find` exists to search topics, and a caller who already holds a path does not need it. → finding F10.

**New finding:**
- F10 (nit) `report_index.go:388`: a topic written as a full `ai-reports/<folder>` path no longer matches its own bundle. Stripping a leading `ai-reports/` from the topic would restore it if wanted. impact-negligible → report only

**Follow-ups created:** None (10 findings report only)
