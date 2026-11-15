# Review: REQ-731 (ai-report prints the catalog rebuild command when a new report is missing from the catalog)

**Approve** — the reminder works for the path forms Step 8 passes, writes nothing, and the merge with REQ-730 lost nothing.
Route B | merge `15ad4b3e` (range `4feebfee..15ad4b3e`, 4 files, +106/-1)

## What's built

- `ai-report-index --check <bundle>` reads only `ai-reports/catalog.json`. It prints one rebuild line when the bundle is not listed, the catalog is missing, or the catalog cannot be parsed. It prints nothing when the bundle is listed.
- Step 8 of `skills/do-work-toolbox/actions/ai-report.md` runs it after the organize check, prints its line with the organize line just before the final bundle-path line, and skips it for a revise.
- The guide's "Report Catalog" section has one new bullet: a new report enters the catalog at the next rebuild, and the closing summary reminds you.

## Decisions / risks for you

- D-04 (a hand-made catalog that lists the bundle prints nothing) is acceptable. The REQ's GREEN condition is "prints nothing when the bundle is in `catalog.json`", and a hand-made catalog that lists it satisfies that literally. When a hand-made catalog does not list the bundle, the line still prints and `ai-report index` then refuses with `AI-REPORT-INDEX-HAND-MADE` and its move-aside fix (checked by hand). The only cost is the record drift in F2.
- D-03 leaves `architecture-report.md:122` and `stakeholder-report.md:54` without the reminder. Those bundles stay missing from the catalog until a rebuild, as today. The REQ names only `ai-report.md` Step 8, so this is not a finding against this REQ.

## Findings

**Important:** None.

**Minor:**
- F1. `report_index.go:112-113` compares the raw argument after trimming one trailing `/`, with no `filepath.Clean` or repo-root strip. `./ai-reports/<folder>`, an absolute path, `ai-reports/<folder>//` and `ai-reports/<folder>/index.html` (the form of the Output Format's final line) all print the rebuild line for a listed bundle (probed in a scratch repo). The forms Step 8 prescribes (`ai-reports/<folder>`, with or without one trailing `/`) work. The failure is a false reminder, never a missed one, and the rebuild it suggests is idempotent. — impact-user-visible → report only
- F2. The REQ's D-02 entry (`do-work/working/REQ-731-report-prints-catalog-rebuild-reminder.md:60`) still says a hand-made catalog "counts as not listing the bundle, so the line still prints". D-04 and the shipped code do the opposite, and D-02 was not amended to point at D-04. The REQ now holds two contradicting records of one decision. The brief (`REQ-731-brief.md:21`) had already changed the rule, so the builder followed its instructions. — impact-negligible → report only

**Nit:**
- F3. `report_index_test.go:474-488` reaches "catalog exists, new bundle not in it" by deleting bundle `b`, rebuilding, and writing `b` again. Writing `b` after the first rebuild would give the same state in fewer lines. — impact-negligible → report only

## Merge check (REQ-730 seam)

- For each of the 4 files, the added and removed lines of `git diff 4feebfee..15ad4b3e` match the builder's own `git diff 1de00505..0e798128` exactly. The merge therefore added only REQ-731's lines on top of main, and no REQ-730 line was removed or changed.
- `report_index.go`: REQ-730's `renderReportSecondOpinionList` stays intact, and `reportIndexCheck` follows it at the end of the file. `reportIndexCheck` decodes into `reportCatalog`, and REQ-730's `second_opinions` field decodes without effect. The usage string `[--find <topic...> | --check <bundle>]` merged cleanly.
- `report_index_test.go`: REQ-730's `TestAIReportIndexListsPlainSecondOpinionPageByFileName` ends intact, and the new test follows.
- `ai-report-guide.md`: the new bullet sits between the hand-made-refusal bullet and the meta-tags bullet. REQ-730's second-opinion bullet follows the meta-tags bullet unchanged. No contradiction.
- `ai-report.md`: the new paragraph follows the organize `--check` block, before `## Output Format`. REQ-730's "The catalog also lists each bundle's second-opinion pages" line (line 48) is untouched.

## Requirements Checklist

- [x] Step 8 checks whether the new bundle is in the catalog and prints the `ai-report index` line when it is not — delivered (`ai-report.md:208-212`).
- [x] The check is a program step on an existing command, like `ai-report-organize --check` — delivered (`--check` branch on `ai-report-index`).
- [x] The line prints just before the final bundle-path line, like the organize line — delivered ("Print its line, if any, with the organize line, just before the final bundle-path line.").
- [x] The guide's "Report Catalog" says a new report enters at the next rebuild and the summary reminds you — delivered (`ai-report-guide.md:102`).
- [x] Creating a report does not rebuild or commit the catalog — delivered. Step 8 says "Creating a report never rebuilds the catalog", the commit rule still says "Never stage anything outside the bundle", and Step 6 still allows the catalog only for the revise form.
- [x] No bundle is edited — delivered (`reportIndexCheck` only reads, and the test snapshots the `ai-reports/` tree around every call).
- [x] When the catalog files do not exist yet, the same line prints — delivered (test case 1, and probed by hand).
- [x] GREEN: line when absent, nothing when listed — delivered (test plus manual probe).
- [x] UR-176 intent (reminder, not rebuild on every report) — delivered.

## Acceptance Testing

**Result: Pass** (implementation and integration stages)
- `go test -count=1 -run TestAIReportIndex ./internal/toolboxcommands/` at the merged tree → `ok` (0.39 s).
- Built the CLI into the scratchpad and ran it in a scratch Git repo:
  - Generated catalog with bundle `a` listed: `--check ai-reports/<a>` and `--check ai-reports/<a>/` → no output, exit 0.
  - New bundle `b` after the catalog → the rebuild line, exit 0.
  - Hand-made catalog listing `b` → no output. The same hand-made catalog with `a` not listed → the line, and the rebuild then refuses with `AI-REPORT-INDEX-HAND-MADE`.
  - Unparseable catalog → the line. Missing catalog with `--format json` → `outcome: success`, `changes: []`.
  - F1 forms (`./`, absolute, `//`, `/index.html`) → the line for a listed bundle.
- Gate evidence from the REQ: `maintainer-verify` green at `15ad4b3e` (141 s). Not re-run, as instructed.

## Suggested Additional Testing

- Live acceptance: unassessed. Run an ordinary `ai-report` in a project whose catalog exists, and confirm the summary prints the rebuild line just before the final bundle-path line, then prints nothing after `ai-report index`.
- Deployment: unassessed. After the next release, confirm an installed copy's `do-work-cli.sh ... ai-report-index --check ai-reports/<folder>` returns the line (the installed binary carries the new flag).

## Anti-bloat count

- Helpers the REQ did not name: 0 in production. `reportIndexCheck` is the function behind the flag that D-01 chose. The test adds 2 local closures (`writeBundle`, `check`), which stay inside the test.
- Options: 1 (`--check`), the flag the REQ's Builder Guidance allows.
- Files: 0 new, 0 undeclared (4 of 6 declared files touched; `CHANGELOG.md` and `VERSION` wait for finalization).
- Decorative tests: 0. One test, which pins the named failure ("a new report missing from the catalog with no reminder"). Its usage-error and after-rebuild cases go beyond the GREEN condition but pin D-07 and the "listed → empty" side, so they are not decorative.

## Scores (on the record, not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Every requirement and constraint delivered |
| Code Quality | 90% | Small, mirrors organize `--check`. F1: no path normalization |
| Test Adequacy | 90% | Focused, red-green shown (RED `a6ac1be6`, GREEN `0e798128`). Hand-made and unparseable cases not pinned |
| Scope | 100% | Declared files only, D-01..D-09 recorded |
| Risk | Low | Read-only, false-positive only |
| Acceptance | Pass | Implementation and integration stages |

## Follow-ups created

None (3 findings report only)

## Review

**Overall: 95%** | 2026-11-15T01:45:50Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1. `--check` does not normalize its argument (`report_index.go:112-113`), so `./ai-reports/<folder>`, an absolute path, `//` or `ai-reports/<folder>/index.html` print a false rebuild line for a listed bundle. The forms Step 8 prescribes work. — impact-user-visible → report only
- F2. The REQ's D-02 entry (line 60) still says a hand-made catalog counts as not listing the bundle. D-04 and the code do the opposite, and D-02 was not amended. — impact-negligible → report only
- F3 (Nit). The test reaches "catalog without the new bundle" by deleting and rewriting bundle `b` (`report_index_test.go:474-488`). Writing `b` after the first rebuild is shorter. — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages: focused Go tests ok at `15ad4b3e`, plus a scratch-repo run of missing, listed, new, hand-made and unparseable catalogs. Merge hunks match the builder's diff byte for byte.
**Restatement sweep:** nothing redefined (a new read-only flag only; the usage string has no restatement outside `report_index.go`, checked with grep)
**Suggested testing:** 2 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*
