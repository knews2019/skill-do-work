# Builder Brief — REQ-730 (the report catalog lists the second-opinion pages of each bundle)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-730-catalog-lists-second-opinion-pages`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-730-catalog-lists-second-opinion-pages`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its decisions D-01 to D-05, Exploration and Scope: read `do-work/working/REQ-730-catalog-lists-second-opinion-pages.md` from the MAIN tree (your worktree may hold an older copy).
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-730-handback.md`
- Up to ten builders run at the same time in this run, all from the same main HEAD. Integration is serial afterwards. Siblings that share your files are listed under **Integration seams** below.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.
2. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, `_dev/primes/prime-action-files.md` (before you edit `ai-report.md`) and `_dev/primes/prime-releases.md`.
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add, focused tests that name the failure they pin).
4. Required lesson: `_dev/primes/lessons-releases.md`. Dropped only for budget, read these bullets in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (find them with a fixed-string grep): `[family: catalog-tie-order]`, `[family: closed-enumeration-for-a-condition]`, `[family: assertion-passes-when-feature-is-dead]`.

## Decided (do not reopen; full text in the REQ)
- **D-01 — which files are pages:** every regular file at the top level of the bundle whose name matches `^index-[A-Za-z0-9][A-Za-z0-9._-]*\.html$`, except the bundle's own entry file. Tree pages (`index-v1c01-….html`) and plain pages (`index-claude-claude-opus-5-5.html`) both match; staging pages start with a dot and never match; `*-reply.md` files are not pages.
- **D-02 — JSON shape:** a new per-bundle field `second_opinions`, always an array (`[]`, never `null`), items `{ "page", "tree_id", "preset", "model", "date" }`. `page` is the repository-relative path, like `entry` (`ai-reports/<bundle path>/<file>`). The other four are the raw values of the `ai-report-second-opinion-id`, `-preset`, `-model`, `-date` metas, empty when absent. Order: byte sort of file names (tree pre-order). Add the field last in `reportBundle`; change no existing field.
- **D-03 — real data:** after the code passes, regenerate the tracked `ai-reports/catalog.json` and `ai-reports/index.html` in your worktree with the new code (`bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format text ai-report-index`) and commit only those two files in their own commit. If `ai-report-index` refuses or prints an `organize` hint, commit nothing extra and report it.
- **D-04 — HTML shape:** inside the bundle's title cell, a nested `<ul class="second-opinions">`, one `<li>` per page: a link to the page (href relative to `ai-reports/index.html`, built the way `reportPageLink` builds the entry link), then `tree id · preset · model · date` with empty parts skipped. With no tree id, the link text is the file name. A bundle without pages renders byte-identical to today. Escape every value.
- **D-05 — one page-name pattern:** declare it once as a named package variable in `report_index.go` (for example `reportSecondOpinionPagePattern`) and make `secondOpinionBanner` (`report_second_opinion.go:1533`) use it in place of its inline regex. No other change in that file.

## What to build
Three commits.

**Commit 1 (RED, tests only).** Add two new test functions at the END of `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go`, each with its own small `t.TempDir()` fixture (do not extend the shared `reportIndexFixture`; the existing naming-style test counts its outputs):
- `TestAIReportIndexListsSecondOpinionPagesUnderTheirBundle` — a bundle with `index.html` and one `index-v1c01-<model>-<slot>.html` page carrying the four metas (write them the way `secondOpinionDecorate` does, `report_second_opinion.go:1626-1630`), plus a second bundle with no pages. Assert: the first bundle's `second_opinions` has exactly that page with the four values; the second bundle's field is an empty array, not null (check the raw JSON bytes, since `json.Unmarshal` turns both into a nil or empty slice); `ai-reports/index.html` links the page inside the first bundle's row. Pin the test's purpose in one comment line: the catalog had one entry per bundle, so second-opinion pages could not be found through it.
- `TestAIReportIndexListsPlainSecondOpinionPageByFileName` — a bundle with `index.html`, a plain `index-claude-claude-opus-5-5.html` without metas, a meta written only inside an HTML comment or a `<script>` (it must not count: read metas through `secondOpinionVisibleHTML`, as `secondOpinionTreeList` does at `report_second_opinion.go:1670-1676`), a `.x-staging-1.html` file and a `v1c01-x-reply.md` file. Assert: exactly the plain page is listed, all four values empty, and the HTML link text is its file name. Add one bundle without `index.html` whose entry falls back to an `index-*.html` file, and assert the entry is not also listed as a page.
Run both and record the failures (the field does not exist) in the hand-back.

**Commit 2 (GREEN).**
1. `report_index.go`: the D-05 pattern variable; a small helper that lists the pages of one bundle folder (it may call `os.ReadDir` again; do not change `pickReportEntry`'s return shape); fill the field in `readReportBundle` (`report_index.go:156-203`) and initialise it to an empty slice beside `LinkedIDs` (`:158`); the nested list in `renderReportIndexPage` (row at about `:410-413`) plus one CSS line if needed.
2. `report_second_opinion.go`: D-05 only.
3. `skills/do-work-toolbox/actions/ai-report.md`, section "Catalog forms: index and find" (`:40-48`): one sentence that the catalog also lists each bundle's second-opinion pages. Do not touch Step 8 or the closing summary (REQ-731 owns them).
4. `skills/do-work-toolbox/docs/ai-report-guide.md`, "Report Catalog" (`:94-102`): add one new bullet after the meta-tags bullet: the catalog lists each bundle's second-opinion pages (`index-*.html`) under the bundle, with tree id, preset, model and date read from the `ai-report-second-opinion-*` metas, and a page without them by its file name. Add it as its own bullet; do not reword neighbours.

**Commit 3 (derived data, D-03).** `ai-reports/catalog.json` and `ai-reports/index.html` only. The bundle `ai-reports/2026-11-14_1547_proposal-tinka-review-surface/` must now show `v1c01`, `v1c02` and `v1c03`. No file inside any bundle changes.

**Not yours:** `VERSION`, `CHANGELOG.md` (the integrator writes the release; propose the entry in the hand-back), anything under `do-work/`, every file inside an `ai-reports/` bundle, `report_judge.go` and its test (REQ-728, REQ-729), `reportIndexFind` behaviour (find still matches bundles only), `completed-work-presentation-reference.md`.

## Rules
- Write only inside the Scope's file list. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-730]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. In `skills/do-work/tools/do-work-cli`: `gofmt -l ./internal/toolboxcommands` prints nothing; `go vet ./internal/toolboxcommands/`; `go test -count=1 ./...` (the REQ's GREEN line; record the time).
2. `bash _dev/tests/contract-regressions.sh` and `bash _dev/tests/audit-lockins.sh` (they read `ai-report.md`).
3. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-730-probe.sh` with your worktree as the current directory prints `REQ-730 probe PASS`. It FAILS at base by design (the new tests do not exist and the tracked catalog has no `second_opinions` field). It assumes the D-02 field names `second_opinions` and `page`, the D-04 href form, and the two test names above; if you must deviate, say so under "Scope changes" so the integrator can adjust the probe.
4. Open the regenerated `ai-reports/index.html` source and check the tinka bundle's row by eye; quote that row in the hand-back.
5. `git diff --name-only $(git merge-base HEAD main) -- ai-reports/` lists only the catalog pair.

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file created, modified or deleted, with the verb.
- `## Scope changes`
- `## Red-green evidence` — the two new tests failing after commit 1 and passing after commit 2; the GREEN probe failing at base and passing at your tip.
- `## Verification` — each command above with its result; for a rerun under load, both runs.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-06 onward (D-01 to D-05 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams` — every shared file below that you touched, with the function or heading you changed and line ranges after your edit.
- `## Proposed CHANGELOG entry` — a descriptive title (for example "Report Catalog Lists Each Bundle's Second-Opinion Pages") and the line "What changes for an installed copy": `ai-report index` now writes a `second_opinions` list per bundle into `ai-reports/catalog.json` (additive; every existing field is unchanged) and shows those pages under their bundle's row in `ai-reports/index.html`; the next `ai-report index` run in a project picks them up, and no bundle is changed.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
All siblings build from the same main HEAD and integrate serially in hand-back order.
| Shared file | REQ-730 edits | Siblings that also edit it |
| --- | --- | --- |
| `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` | `reportBundle`, `readReportBundle`, `renderReportIndexPage`, one new helper and one new pattern variable | none expected |
| `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` | `secondOpinionBanner`, one line (D-05) | none expected; REQ-728 and REQ-729 edit `report_judge.go`, not this file |
| `skills/do-work-toolbox/actions/ai-report.md` | "Catalog forms: index and find" only | REQ-728, REQ-729 (judge usage), REQ-731 (Step 8 and the closing summary: the catalog rebuild reminder) |
| `skills/do-work-toolbox/docs/ai-report-guide.md` | "Report Catalog", one new bullet | REQ-728, REQ-729 (judge), REQ-731 (likely the same "Report Catalog" section) |
| `ai-reports/catalog.json`, `ai-reports/index.html` | regenerated (D-03) | REQ-725 adds a report bundle under `ai-reports/`, and REQ-733 rebuilds a page inside one bundle; neither should regenerate the pair. On any conflict in these two derived files, the integrator reruns `ai-report-index` on the merged tree and never merges them by hand. |

Your final message: under 2500 characters. State the branch tip hash, the verification results, and anything the integrator must decide.
