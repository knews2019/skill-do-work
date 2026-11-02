## Review: REQ-706

**Approve**: the ai-report catalog, find, second-opinion and the revise rule now read report bundles inside month folders, and a repository with no month folders gets byte-identical output. One reader the REQ did not name, the architecture-report pre-flight, still looks only one level deep (F1).
Route B | merge range `6e98d229..3edfd2c9` (builder commit `c004cba0`)

### What's built
- `ai-report index` and `find` list each bundle inside a month folder (`YYYY-MM`, `YYYY-MM_1-report`, `YYYY-MM_<N>-reports`) with path `ai-reports/<month folder>/<bundle>`. The month folder itself is never listed. The supersedes chain resolves between a flat bundle and a nested one.
- `ai-report second-opinion` accepts a nested bundle, a nested page target and `latest`. A symlinked month folder is still refused.
- The revise and second-opinion rules in `ai-report.md` and the user guide describe the month-folder layout. New reports are still written directly under `ai-reports/`.
- Still missing (outside the REQ's named paths): `architecture-report-preflight --scan` does not see a prior architecture report that sits in a month folder (F1).

### Decisions / risks for you
- F1 becomes common once REQ-707 (the `ai-report organize` command) moves finished months into month folders. The cheapest fix is an addendum to REQ-707 before it runs: reuse `reportMonthFolder` in `architectureScan` and add one focused test. Value: the next architecture report keeps its "changed since last report" baseline. Risk if skipped: after the first organize run, every architecture report says it is the first one.

### Findings

**Important:**
- F1: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go:62-85`. `architectureScan` reads only the direct children of the reports folder. A prior architecture report filed in a month folder is not found, so `prior_report=` is empty, and `skills/do-work-toolbox/actions/architecture-report.md:83` then tells the author to write "no prior HTML baseline". I reproduced this in a scratch repository: the nested `2026-10_1-report/2026-10-20_1000_architecture-report` gives an empty `prior_report`, and the same bundle placed directly under `ai-reports/` is found. This is a code consumer of the redefined contract that still uses the old one. Suggested fix: descend one level into children of `reportMonthFolder` folders in `architectureScan` and add a focused test, preferably folded into REQ-707. — impact-user-visible → report only

**Minor:**
- F2: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:387-389`. The predicate accepts `ai-reports/<month folder>` itself, because its parent is `ai-reports`. Second-opinion therefore prepares on a month folder that holds a regular `index.html`, even though the catalog never lists that folder and the REQ says a month folder is never a bundle. The bundle boundary, including the outside-bundle Git check, would then cover every nested bundle. Reproduced: `--prepare --backend claude --model claude-fable-5-1 ai-reports/2026-05_1-report` returned `prepared`. Suggested fix: also refuse when `parent == reportsFolder && reportMonthFolder(filepath.Base(relative))`, with one assertion in `TestSecondOpinionBundleInsideMonthFolder`. — impact-negligible → report only
- F3: `skills/do-work-toolbox/actions/ai-report.md:52` and `:203` still call the revision a "sibling" bundle. Revise step 1 (`:54`) now says the revision is written directly under `ai-reports/` even when the prior bundle sits in a month folder, and the guide already dropped the word "sibling". Suggested fix: "writes a new bundle directly under `ai-reports/`" and "a new `-rev<N>` bundle". — impact-negligible → report only

**Nit:**
- F4: `report_second_opinion.go:389` mixes `||` and `&&` without parentheses (`err != nil || parent != reportsFolder && !inMonthFolder || ...`). Go's operator order makes it correct, but it is easy to misread. Suggested fix: put parentheses around the middle clause. — impact-negligible → report only
- F5: `report_index.go:47`. The regexp accepts `_<N>-reports` for any N, including `0` and `1`, but accepts `_N-report` only when N is 1. So a hand-named `2026-07_2-report` is read as one broken bundle and its reports are hidden (scratch run). The REQ's exact-name rule allows this. However, REQ-707 reuses the helper and cannot recognise such a folder in order to recount it. Suggested fix: none needed now. Decide in REQ-707 whether to accept `_[0-9]+-reports?`. — impact-negligible → report only

**Anti-bloat count:** 0 unrequested helpers, 0 options or flags, 0 new files, 0 decorative tests. `reportMonthFolder` and `reportMonthPattern` are the one helper and the one regexp that exploration decision D-01 asked for. Each of the two new tests has a `// Pins:` comment that names the failure it guards. `writeReportFiles` is a local helper inside one test and is used twice. The guide's new layout paragraph is part of the declared Scope.

### Requirements Checklist

- [x] A direct child named `YYYY-MM`, `YYYY-MM_1-report` or `YYYY-MM_<N>-reports` is a month folder. Its non-dot child directories are bundles with path `ai-reports/<month>/<bundle>`, and the folder itself is never a catalog bundle. Delivered (see F2 for second-opinion).
- [x] Flat bundles keep working. Loose files and symbolic links are never bundles. A symlinked month folder stays refused. Delivered.
- [x] Second-opinion accepts a bundle whose parent is `ai-reports` or a month folder, and `latest` resolves nested bundles. Delivered.
- [x] Revise accepts `ai-reports/<folder>`, `ai-reports/<month>/<folder>` or a bare folder, and the supersedes chain resolves across flat and nested bundles. Delivered: prose rule plus a catalog test.
- [x] `find` matches nested bundles by folder name. Delivered.
- [x] The action text and the guide no longer say a bundle must sit directly under `ai-reports/`. Delivered (F3 is a leftover "sibling" wording, not a direct-child rule).
- [x] A repository with no month folders gets a byte-identical catalog. Delivered and independently verified.
- [x] UR constraints: writers unchanged, no write-side change, nothing moves. Delivered (no writer file touched).

### Acceptance Testing

**Result: Pass** (Implementation and Integration stages)
- `go test ./internal/toolboxcommands/ -count=1 -run 'AIReportIndex|SecondOpinion'` at `3edfd2c9`: ok (16.1 s).
- Scratch CLI build of `3edfd2c9` on a fixture outside ROOT. Catalog: 8 bundles. Nested bundles were listed under `2026-09`, `2026-09_1-reports`, `2026-06_01-reports` and `2026-05_1-report`. The empty `2026-09_0-reports` and `2026-10` added nothing. `2026-13` and `2026-07_2-report` were read as broken bundles (F5). The double-nested `2026-04/2026-03` was listed as a bundle, as decision D-02 says. The `index.html` that sits inside `2026-05_1-report` itself was not listed. Page links read `2026-09/2026-09-27_2200_a/index.html`. `--find 2026-09-27` found the nested bundle.
- An unreadable month folder (mode 000) fails the whole index with `AI-REPORT-INDEX-FAILED`, as decision D-07 says.
- Second-opinion `--prepare`: the nested page target `ai-reports/2026-09/2026-09-27_2200_a/index.html` was prepared, and so was `latest`. A month folder without `index.html` was refused ("must hold a regular index.html"). A nested bundle under a non-month folder (`2026-13/...`) was refused with the new message. The month folder itself with an `index.html` was prepared (F2).
- Byte-identical check: I built `6e98d229` and `3edfd2c9` separately and ran both on the same flat fixture (8 bundles, a hidden folder, a loose file, a symlink, a supersedes pair). `catalog.json`, `index.html` and the `--find proposal` output matched byte for byte (`cmp`).
- F1 was reproduced with `architecture-report-preflight --scan ai-reports`.

### Suggested Additional Testing

- Deployment: unassessed. In the consumer repository from the UR-161 inbox note, after updating the suite, run `ai-report index` and confirm about 248 bundles and no bundle named `2026-03` to `2026-10`.
- Live acceptance: unassessed. Open the regenerated `ai-reports/index.html` in a browser and click a nested report link. Then run one real second-opinion backend end to end on a nested bundle, including `--finalize`, the render judge and banner insertion. I exercised only `--prepare`; `--finalize` relies on the same `secondOpinionBundle` check.
- Heavy lanes not run, by the user's test decision: `do-work-cli-integrations`, `staged-skills`, `updater`, `installer`.
- Edge case: duplicate bundle names (flat plus nested, or in two month folders) make supersedes and a bare `<folder>` revise target ambiguous. The builder already recorded this as a Discovered Task for REQ-707.

### Scores (on the record, not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All 7 REQ criteria and the UR constraints delivered |
| Code Quality | 90% | Small, follows the file's patterns; F2 boundary gap, F4 readability |
| Test Adequacy | 90% | RED-to-GREEN shown, two focused tests that guard named failures, mutation check; no refusal test for the month folder itself |
| Scope | 100% | 6 declared files, 6 touched; decisions D-01 to D-11 recorded |
| Risk | Low | Symlink guard kept on both levels; F1 is a sibling reader that misses month folders, not a regression |
| Acceptance | Pass | Implementation and Integration stages |

### Follow-ups created
- None (5 findings report only)

## Review

**Overall: 95%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1: `architecture.go:62-85` `architectureScan` reads only the direct children of `ai-reports/`, so a prior architecture report in a month folder is missed, `prior_report` is empty, and the next report claims no prior baseline (reproduced). Suggested fix: reuse `reportMonthFolder` there with a focused test, folded into REQ-707. — impact-user-visible → report only

**Minor findings:** F2: `report_second_opinion.go:387-389` accepts the month folder itself as a bundle when it holds a regular `index.html` (reproduced `prepared`), although the catalog never lists it; refuse a month-named direct child. — impact-negligible → report only. F3: `ai-report.md:52` and `:203` still call the revision a "sibling" bundle, which contradicts step 1's directly-under-`ai-reports/` rule; reword. — impact-negligible → report only. Nit F4: `report_second_opinion.go:389` mixes `||` and `&&` without parentheses. — impact-negligible → report only. Nit F5: `report_index.go:47` accepts `_N-reports` for any N but `_N-report` only for N=1, so a hand-named `_2-report` folder hides its reports. — impact-negligible → report only
**Acceptance:** Pass — Implementation and Integration stages: focused `AIReportIndex|SecondOpinion` tests green at `3edfd2c9`; a scratch CLI build confirmed the nested catalog, find, page links, second-opinion prepare on a nested bundle, a page target and `latest`, and a byte-identical flat catalog against a `6e98d229` build. Deployment and Live acceptance unassessed.
**Restatement sweep:** redefined where a report bundle may sit under `ai-reports/` (directly, or inside a month folder that is never itself a bundle). Swept `skills/`, `suite/` and `README.md` (no root `docs/` exists) for direct-child, every-directory-is-a-bundle, sibling and one-level readers: `ai-report.md:42/54/72` and `ai-report-guide.md:74/95/97/105/117` agree. Stale: `architecture.go:62-85` (F1) and `ai-report.md:52/203` "sibling" (F3).
**Suggested testing:** 4 items
**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*

## Re-check (3edfd2c9..ba902c63)

*Written by the reviewer at 2026-11-02T17:34Z; its own append command stalled (it began with a scratch-directory removal outside the repository that never ran), so the orchestrator copied the section verbatim from the reviewer's transcript at 2026-11-02T17:55Z.*

**Approve**: F1 to F5 are all fixed. The delta (builder commit `a7a2fe6b`, 7 files, +60/-12) adds nothing beyond the five findings. The one new item is a Nit. Overall now **97%**, acceptance **Pass** (Implementation and Integration stages).

| Finding | Status | Evidence |
|---|---|---|
| F1 (architecture pre-flight one level) | fixed | `architectureScan` (`architecture.go:62-105`) collects direct folders plus the child directories of `reportMonthFolder` folders, then matches on the base name. New `TestArchitectureScanFindsPriorReportInMonthFolder`, RED shown. Scratch run: a nested `-2` same-minute prior is picked over an older one; a symlinked month folder holding a newer report is skipped; an absolute `--scan` path works; an unreadable month folder fails with `ARCHITECTURE-PREFLIGHT-FAILED` (same choice as decision D-07). |
| F2 (month folder itself accepted) | fixed | `flatBundle` excludes month-folder names (`report_second_opinion.go:388-391`). The new test assertion was RED before the fix. Scratch run: `ai-reports/2026-07_2-report`, its `index.html` page target, and a month-named `2026-09` holding `index.html` are refused. The nested bundle, a flat bundle and `latest` are still prepared. |
| F3 ("sibling" wording) | fixed | `ai-report.md:52` and `:203` now say the new bundle goes directly under `ai-reports/`. |
| F4 (unparenthesized condition) | fixed | Replaced by the named `!(flatBundle \|\| inMonthFolder)` (decision D-16). |
| F5 (`_N-report` only for N=1) | fixed | Pattern `^[0-9]{4}-(0[1-9]\|1[0-2])(_[0-9]+-reports?)?$` and an updated doc comment (decision D-13). New `2026-09_2-report` subtest, RED shown. The prose keeps naming the canonical forms only (decision D-14); I agree, because organize (REQ-707) renames to those. |

**New finding:**
- N1 (Nit): `report_second_opinion.go:392`. Naming a month folder itself returns "bundle must be a physical report bundle directly under ai-reports or under a month folder", but the folder *is* directly under `ai-reports/`, so the message does not say why it was refused. Suggested fix: a separate message such as "a month folder is a container, not a bundle; name a report inside it". — impact-negligible → report only

**Checks run:** `go test ./internal/toolboxcommands/ -count=1 -run 'SecondOpinion|AIReportIndex|Architecture'` at `ba902c63`: ok (16.0 s). `go vet` and `gofmt -l` were clean. Scratch builds of `6e98d229` and `ba902c63` on one flat fixture gave identical `catalog.json`, `index.html`, `--find` output and architecture `prior_report`/`prior_hash` output (`cmp`). So a repository without month folders is still unchanged, including the architecture pre-flight.

**Second sweep for one-level readers of `ai-reports/` (Go and shell under `skills/`):** none left.
- All 59 `ReadDir`/`Glob`/`WalkDir` sites in shipped Go were checked. Only `report_index.go` (catalog walk) and `architecture.go` (pre-flight scan) read the reports folder, and both now descend into month folders.
- The other readers inside `report_second_opinion.go` read inside one bundle (presets, tree pages, snapshot).
- The recursive repository walkers (cleanup, inventory, memory) see nested bundles by construction.
- `report_judge.go` takes any directory that holds `index.html`.
- No `.sh`, `.py` or `.js` file under `skills/` reads `ai-reports/`. The launcher only exports the prompt path.
- No prose "newest/prior report" lookup exists outside the already-updated `ai-report.md` and guide lines. The `blocked_by`/`## Reports` paths are exact links that organize (REQ-707) must rewrite.

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% (extension to `architecture.go`/`architecture_test.go` recorded as decision D-12) |
| Risk | Low |
| Acceptance | Pass |

**Restatement sweep (re-check):** redefined the month-folder name rule (`_<N>-report` with any N); its prose restatements name only canonical forms and still agree. Nothing else redefined.
**Follow-ups created:** None (1 finding report only)
