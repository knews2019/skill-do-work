# Review: REQ-730 (the report catalog lists the second-opinion pages of each bundle)

**Approve.** The catalog now lists every second-opinion page under its bundle, in both `catalog.json` and `index.html`. The new tests fail when any part of the feature is removed. No existing field changed, and no bundle is written.
Route B | range `72f72c3a..afa0153f` (merge `afa0153f`)

### What's built
- `ai-report index` writes a new `second_opinions` array on every bundle in `ai-reports/catalog.json`. Each item holds `page`, `tree_id`, `preset`, `model` and `date`, read from the page's `ai-report-second-opinion-*` meta tags.
- `ai-reports/index.html` shows the pages as a nested list in the bundle's title cell. A page without a tree id is linked by its file name.
- The regenerated real catalog lists the Tinka bundle with four pages: v1c01, v1c02, v1c03 and v1c03c01 (`index-v1c03c01-claude-claude-opus-5-5-custom-left-to-right-map.html`).
- Still to come at finalization: `CHANGELOG.md` and `VERSION` (the integrator writes them).

### Decisions / risks for you
- None. D-01 (every `index-<name>.html` except the entry counts as a page) can list a hand-made `index-foo.html`. The REQ already accepted that risk, and it is visible and reversible.

### Findings

**Important:**
- None.

**Minor:**
- F1. `report_index_test.go:318` adds `writeReportIndexFiles`, which has the same body as `writeOrganizeFiles` (`report_organize_test.go:25`, same package) and as the loop inside `reportIndexFixture` (`report_index_test.go:39`). The new tests could call `writeOrganizeFiles` with `ai-reports/` key prefixes. — impact-negligible → report only
- F2. The "Report Catalog" bullet (`skills/do-work-toolbox/docs/ai-report-guide.md:103`) does not name the JSON field `second_opinions` or its item keys. The REQ's Scope said this section covers "the new field and what it reads". A script author who reads `catalog.json` must read the Go code to find the name. — impact-user-visible → report only

**Nit:**
- F3. Two descriptions are slightly wider than the code. The guide says `index-*.html`, but the bundle's own entry file is excluded and the character after `index-` must be a letter or a digit. The comment on `readReportSecondOpinionPages` (`report_index.go:486`) says byte order is "tree pre-order". That holds for tree pages only, because plain pages sort among them by name. — impact-negligible → report only

### Requirements Checklist
- [x] `catalog.json` gains a per-bundle list with tree id, preset, model and date. The change is additive: `SecondOpinions` is the last `reportBundle` field, and no existing tag changed (`report_index.go:80-91`).
- [x] `ai-reports/index.html` shows the pages under their bundle's row. Five rows carry a `<ul class="second-opinions">` list, and the Tinka row lists all four pages.
- [x] A plain page is listed by its file name. `TestAIReportIndexListsPlainSecondOpinionPageByFileName` pins it, and the real REQ-632 row shows `index-claude-claude-opus-5-5.html` as link text.
- [x] Every existing field keeps its name and meaning. `reportIndexFind`, `secondOpinionBundle` (`latest`, which reads only `catalog[0].Path`) and organize do not read the new field. The only `DisallowUnknownFields` in the package is in `periodic_check.go`, which has nothing to do with the catalog. No existing test changed.
- [x] No bundle is edited. `readReportSecondOpinionPages` calls only `os.ReadDir` and `os.ReadFile`. A regeneration in a scratch copy changed only the catalog pair (checksum compare of every tracked `ai-reports/` file).
- [x] RED then GREEN. The hand-back quotes both new tests failing after `f00c54d5` and passing after `359bb427`. The gate ran green at `afa0153f`.
- [ ] `CHANGELOG.md` and `VERSION`: not in the range by design. The integrator writes them at finalization.

### Code review notes
- Field always an array. Both the `bundle :=` literal and the helper start from an empty slice, and the helper returns that slice on a read error. A missing bundle folder therefore still gives `[]`.
- Entry exclusion. `name == entryName` uses the name returned by `pickReportEntry`, so a bundle without `index.html` whose entry is an `index-*.html` page does not list that page twice.
- Metas. The code reads them with `readReportMetas(secondOpinionVisibleHTML(data))`, the same call `secondOpinionTreeList` uses. A meta inside a comment or a script does not count. `readReportMetas` unescapes the value, and the renderer escapes it once, so no value is double-escaped.
- Href escaping. The code applies `url.URL{Path: …}.String()` and then `html.EscapeString`, the same form as `reportPageLink`. The page-name pattern allows only `[A-Za-z0-9._-]`, so only the bundle folder part can need escaping.
- Byte order. `os.ReadDir` returns entries sorted by file name. The real catalog shows v1c03 before v1c03c01, because `-` sorts before `c`.
- Symlinks and staging pages. `IsRegular()` skips symlinks, the same rule `pickReportEntry` uses. The pattern requires a name to start with `index-`, so dot-prefixed staging pages never match.
- `secondOpinionBanner` (`report_second_opinion.go:1533`). The regex text is byte-identical. The only change is that the code compiles it once at package load instead of on every call, so behaviour does not change. `TestSecondOpinionBanner*` passes.
- Docs. `ai-report.md:48` adds one true sentence in "Catalog forms: index and find", after the command block. REQ-728's `--page` text sits at lines 36 and 190, in the judge paragraph and Step 7, and does not overlap it. The guide bullet (`ai-report-guide.md:103`) sits in "Report Catalog", after the meta-tags bullet. REQ-728's text is at lines 59 and 80, in other sections. Neither new line contradicts the merged text.

### Anti-bloat count (items the REQ did not name)
- 1 test helper: `writeReportIndexFiles` (F1).
- 1 rendering function: `renderReportSecondOpinionList`. It is not named by the REQ, but D-04 implies it. It keeps the row `Fprintf` readable. Accepted.
- 1 CSS line: `ul.second-opinions{…}`. The hand-back records it. It changes the page header bytes for every catalog, while rows without pages stay byte-identical. Accepted.
- 0 new options, 0 new files, 0 decorative tests. The pattern variable, the item type and `readReportSecondOpinionPages` are named by D-05, D-02 and the Exploration.

### Acceptance Testing

**Result: Pass**
- Ran `go test -count=1 -run 'TestAIReportIndex|TestSecondOpinionBanner|TestAIReportOrganize' ./internal/toolboxcommands/`: ok.
- Regenerated the catalog from the tracked `ai-reports/` files at `afa0153f` in a scratch git repository with `go run ./cmd/do-work-cli --repo-root <scratch> --format text ai-report-index`. Output: "49 bundles". Both `catalog.json` and `index.html` are byte-identical to the committed pair, so the pair is the real output of the merged code, not a hand-merge.
- Ran mutations in a scratch copy of the module. All five failed at least one new test:
  - M1: removing the entry exclusion fails the plain-page test.
  - M2: dropping the visible-HTML filter fails it (it reads `tree_id v1c09` from a comment).
  - M3: a nil slice fails both tests (`null`).
  - M4: dropping the HTML list fails both tests.
  - M5: dropping the field fill fails both tests.
  The two new tests are not decorative.

### Suggested Additional Testing
- Live acceptance: unassessed. Open `ai-reports/index.html` in a browser at wide and phone widths, in light and dark mode, and check that the nested list reads well and does not overflow.

### Scores (on the record, not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 97% | All met. Changelog and version come at finalization. |
| Code Quality | 95% | Follows the existing catalog and tree-list patterns. One duplicate test helper. |
| Test Adequacy | 95% | Two tests. Five mutations each fail at least one of them. |
| Scope | 95% | Declared files only. One CSS line and one test helper beyond the REQ. |
| Risk | None | Read-only walk. Additive JSON field. |
| Acceptance | Pass | Committed pair reproduced byte-identically. Tinka row lists four pages. |

### Follow-ups created
- None (3 findings report only)

## Review

**Overall: 95%** | 2026-11-15T01:17:20Z

| Dimension | Score |
|-----------|-------|
| Requirements | 97% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 95% |
| Risk | None |
| Acceptance | Pass |

**Verdict:** Approve. No finding must be fixed in this REQ.

**Important findings (each with its recorded impact token, the durable audit record the judgment requires):**
- None

**Minor findings:**
- F1 (Minor). `writeReportIndexFiles` (`report_index_test.go:318`) duplicates `writeOrganizeFiles` (`report_organize_test.go:25`, same package). — impact-negligible → report only
- F2 (Minor). The "Report Catalog" bullet (`ai-report-guide.md:103`) does not name the `second_opinions` field or its item keys. — impact-user-visible → report only
- F3 (Nit). The guide's `index-*.html` and the "tree pre-order" comment (`report_index.go:486`) are slightly wider than the code. The entry file is excluded, and plain pages sort among the tree pages by name. — impact-negligible → report only

**Acceptance:** Pass. The focused tests pass, and the committed catalog pair regenerates byte-identically from `afa0153f`. The Tinka row lists v1c01, v1c02, v1c03 and v1c03c01, and five mutations each fail at least one new test.
**Restatement sweep:** nothing redefined. The diff adds the `second_opinions` field and moves the page-name pattern into a shared variable, and no existing contract changes meaning. Per the dispatcher, this is not the wave's last integration, so no inherited sweep was run.
**Suggested testing:** 1 item
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*
