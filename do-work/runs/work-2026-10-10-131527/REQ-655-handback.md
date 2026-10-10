# Hand-back: REQ-655 (ai-report index and find build a derived catalog of every report bundle)

- Branch: `worktree-agent-REQ-655-ai-report-index-and-find`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-655-ai-report-index-and-find`
- Base commit: `bd56c4b0`
- Commits:
  - `9170e02c` [REQ-655] add ai-report-index catalog and find verb (the four Go files)
  - `e56dc378` [REQ-655] wire ai-report index and find into the toolbox prose (the four prose files)
- Nothing under `do-work/` was staged or committed. Nothing was written into the worktree's `ai-reports/`.

## File manifest

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (new): the `ai-report-index` verb. It walks the bundles, builds the catalog, checks for hand-made files, publishes `catalog.json` and `index.html` with `rootedPublishFile`, and handles `--find`.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (new): one shared fixture and the four named tests.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` (modified): `CommandAIReportIndex` constant and its `Handlers()` entry, each on its own line under `CommandArchitecture`.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (modified): the name list now has one name per line around the new name, and the count is 8.
- `skills/do-work-toolbox/actions/ai-report.md` (modified): a new `### Catalog forms: index and find` subsection after Input. The `:109` sentence and the `:140` checklist item now name the catalog forms as the one exception.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified): two lines added at the end of the Input block, and a new `## Report Catalog` section at the end of the file.
- `skills/do-work-toolbox/SKILL.md` (modified): `` `report index`, `is there a report on` `` appended to the end of the ai-report trigger cell.
- `skills/do-work-toolbox/actions/help.md` (modified): one new line under the `ai-report [REQ|UR]` line.

## P-A-U

- **[PLAN]:** I followed the REQ `## Plan` and the brief step by step. I read all of the REQ, the brief, UR-144 A2 (as data), the crew rules, the action-files, releases and shell primes, `lessons-releases.md`, and the three named lesson families. The tests were written first. Then I built one Go file holding a walk shared by both forms, then the prose. The deviations are D-11 to D-22 below.
- **[APPLY]:** Code was written as planned, and only the eight files in Scope changed. `architecture.go` was not changed (D-08). There is no `--dry-run`, `--commit`, `--reports` or Just recipe. No `runTransaction`. No "may be stale" mark.
- **[UNIFY]:** `git diff bd56c4b0 --stat`:
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

## Proof record

- **RED** (tests written, verb absent): `go test -C skills/do-work/tools/do-work-cli -count=1 -run 'TestAIReportIndex' ./internal/toolboxcommands/` gave a build failure: `undefined: CommandAIReportIndex`, `undefined: reportCatalog`, `undefined: reportIndexGenerator`, `undefined: reportBundle`, then `FAIL ... [build failed]`. At base, `bash skills/do-work/tools/do-work-cli.sh --repo-root "$(mktemp -d)" --format json ai-report-index` returned `"outcome": "failure"` with finding `UNKNOWN-COMMAND` ("command \"ai-report-index\" is not available").
- **GREEN:**
  ```
  --- PASS: TestHandlersRegisterCanonicalToolboxCommands (0.00s)
  --- PASS: TestAIReportIndexCatalogsEveryBundleNamingStyleOnce (0.03s)
  --- PASS: TestAIReportIndexLinksSupersededProposalToSuccessor (0.03s)
  --- PASS: TestAIReportIndexFindListsMatchesNewestFirstIncludingSuperseded (0.01s)
  --- PASS: TestAIReportIndexRefusesHandMadeCatalog (0.06s)
  ```
- **Smoke** (this repo's `ai-reports/` copied into a `mktemp -d` tree, which was removed afterwards):
  - `--format text ai-report-index` exited 0 in 3 s, including the launcher build. It printed `ai-reports/catalog.json: 41 bundles` and `ai-reports/index.html`. The loose `2026-05-28_2335_background-agents-durability.html` was skipped.
  - `2026-08-26_1709_architecture-report` got entry `.../architecture-report.md`, title "Architecture Report — skill-do-work", date `2026-08-26T17:09` and kind `architecture-report`.
  - `2026-08-27_1428_req-341-timeline-drag-evidence` got `linked_ids` `["REQ-341","REQ-375"]`, with the underscore-joined lowercase id found.
  - Kinds: 32 unknown, 7 architecture-report, 2 proposal. No bundle lacked an entry.
  - `--find architecture` printed 7 rows newest first, from 2026-09-06_1740 down to 2026-08-26_1709.
  - `--find REQ-0341` matched `REQ-341`.

## Decisions

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

## Discovered Tasks

- impact-low: other toolbox verbs whose prose uses `--format json` may return only `exactOutputResult`. If so, the JSON caller sees no payload, the trap D-11 avoided here. Audit the `exactOutputResult` callers against their action prose. → report only
- impact-low: this repo's own `ai-reports/2026-05-28_2335_background-agents-durability.html` is a loose file, so `ai-report index` and `find` skip it by design. If it is a real report, it is invisible to the catalog until someone moves it into a bundle folder. → report only

## Lessons read

- `_dev/primes/lessons-releases.md` (whole file): family `canonical-link-outlives-its-target`. Not triggered, because no release links were written.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, the bullets for:
  - `closed-enumeration-for-a-condition` (lines 7, 28, 67, 70, 72, 86). Applied: a bundle is a condition (any non-dot directory), no list of naming styles exists anywhere, and the known-kind table is marked illustrative with the meta winning.
  - `exact-basename-authority` (line 78). Applied: the supersedes match uses the exact last path segment, with no trimming beyond a trailing `/`.
  - `silent-skip-reads-as-red` (lines 9, 76). Applied as D-22.
- Primes: `prime-action-files.md` (Cross-Referencing, Descriptions Are Triggers), `prime-releases.md`, `prime-shell-commands.md` (one launcher line copied from the `architecture-report.md:41` pattern).

## Anti-bloat check

The diff stat is in [UNIFY] above. These are the items I added that the brief does not name, each with its reason:

- Constants `reportIndexMarker`, `reportsFolder`, `reportCatalogFile`, `reportIndexPageFile`: each string is used in two or more places (writer, refusal check, page).
- Functions:
  - `reportIndexWalk`, `readReportBundle`, `pickReportEntry`, `readReportMetas`, `cleanReportText`, `reportKind`: the walk rules in brief step 2, one function per rule.
  - `reportIndexWrite`, `renderReportCatalog`, `renderReportIndexPage`, `reportPageLink`: brief steps 3 to 5.
  - `reportIndexFind`, `normalizeReportID`: brief step 6.
  - `reportIndexFailure`: D-17.
- Test helpers `runReportIndex`, `readReportCatalog`, `snapshotReportTree`: shared by the four tests. `snapshotReportTree` is the in-test form of "git status shows no other file changed".
- Fixture extras:
  - `.hidden-bundle/`, the `linked-bundle` symlink and `2026-09-09_empty-bundle/` pin the dot-skip, link-skip and no-entry rules.
  - `2026-09-06_1000_architecture-report/architecture-report.md` is the named-`.md` bundle and pins the slug-kind rule.
  - The `aaa.md` decoy pins pick order.
  - `REQ-0412-delta-fix` uses an `<h2>`, which pins the heading fallback.
- Two extra assertions inside the find test: an id-topic match that ignores leading zeros (`req-412` finds `REQ-0412`), and the no-match line. Both are behaviours the brief specifies, with no new test functions.
- None of these is an option, a flag or a new file.

## Proposed CHANGELOG entry (integrator adds the version)

**ai-report Index and Find for Every Report Bundle**

`ai-report index` and `ai-report find <topic>` answer "is there a report on X" in one command, across every folder naming style a project has used.

- A new toolbox CLI verb, `ai-report-index`, treats every directory under `ai-reports/` as a report bundle, whatever its naming style. It writes a derived `ai-reports/catalog.json` and a static `ai-reports/index.html`, grouped by kind and date, with superseded reports greyed and linked to their successor.
- Each bundle records path, entry file, title, date, kind (from `ai-report-kind` or a known slug kind), linked UR/REQ ids, verdict (`ai-report-verdict`), `supersedes` and `superseded_by`.
- `--find <topic>` prints matching bundles newest first, superseded ones included, and writes nothing.
- A hand-made `catalog.json` or `index.html` is refused (`AI-REPORT-INDEX-HAND-MADE`), never overwritten. Rerunning over the verb's own output is byte-identical.
- The toolbox router, help and ai-report guide list the two new forms.

## Proposed lesson bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`

- [family: regexp-word-boundary-underscore] [REQ-655: Go RE2 `\b` counts `_` as a word character, so `\b(UR|REQ)-\d+\b` silently misses ids joined by `_` in folder names (`2026-08-27_1428_req-341-...`); match ids with an explicit `(?:^|[^A-Za-z0-9])` lead and pin an underscore-joined fixture](../../../../do-work/archive/UR-144/REQ-655-ai-report-index-and-find-catalog.md)

## Integration seams hit

- `commands.go`: the new constant and map entry each sit on their own line directly under `CommandArchitecture`. `commands_test.go`: the name list is now one name per line for `CommandNote`, `CommandArchitecture` and `CommandAIReportIndex`, and the count is `8`. The REQ-657 integrator sets the final sum.
- Prose: all edits are in the locations the brief named. Expect line-level conflicts on `SKILL.md:23` (the cell grew at its end) and on the `ai-report.md` Step 6 no-search sentence (old `:109`, now `:119`), which had a clause appended.
- Test wall times: the four new tests take 0.01 to 0.06 s each. The whole `toolboxcommands` package runs in 4.5 s.
