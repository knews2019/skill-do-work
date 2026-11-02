# Upstream suggestion for `knews2019/skill-do-work` — let ai-report tools find bundles inside `ai-reports/YYYY-MM/`

**How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
consumer repo it was authored in. Observed against **v0.307.0**; all line numbers are from that tag.

---

## Request

Today every ai-report tool treats a bundle as a directory directly under `ai-reports/`. That works while a
repo has a few dozen reports. A consumer repo with about 250 bundles found the flat folder unreadable in an
editor's file tree, and it filed every bundle into a month folder: `ai-reports/2026-10/2026-10-28_0103_x/`.
After that move, three suite paths no longer see the bundles. The index would catalogue each month folder as
one bundle. Revise and second-opinion refuse any bundle that is not a direct child. The consumer now avoids
`ai-report index` and cannot revise or second-opinion a filed report.

Please recognize an optional month layer, a directory named exactly `YYYY-MM` (4 digits, `-`, 2 digits), as a
container of bundles and never a bundle itself:

1. **Catalog walk** (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go:79-92`,
   `reportIndexWalk`): when a direct child of `ai-reports/` is a month folder, walk its children as bundles,
   with `Path` = `ai-reports/YYYY-MM/<folder>`. Flat bundles keep working.
2. **Second-opinion bundle check** (`report_second_opinion.go:386-389`): accept a bundle whose parent is
   `ai-reports` or `ai-reports/<YYYY-MM>`.
3. **Revise target rule** (`skills/do-work-toolbox/actions/ai-report.md:53`): allow `ai-reports/YYYY-MM/<folder>`
   and a bare `<folder>` found inside a month folder. Use the same wording for the second-opinion rule at
   `ai-report.md:71`.

Optional, in a separate step: new bundles could be written straight into `ai-reports/<YYYY-MM>/` when the
repo already has month folders. Collision-Safe Publication would then check that path.

Nothing here changes the bundle contract: a bundle is still one folder holding `index.html`, the catalog
keeps its fields, and folder names keep their `yyyy-mm-dd_hhmm_<slug>` form. The month folder is only a
parent directory.

## What happened, in one consumer repo (timeline)

The consumer ran suite v0.307.0 the whole time. Its own tooling did the move. The suite was not edited.

| When (UTC) | Suite | Event |
| --- | --- | --- |
| Oct 28 | 0.307.0 | The user asked to "organize per month": `ai-reports/` held 248 bundle folders (Mar 1, Jun 6, Jul 56, Aug 48, Sep 28, Oct 109). |
| Oct 28 | 0.307.0 | A repo-local script moved every bundle into `ai-reports/YYYY-MM/` and rewrote 660 files of references to them. |
| Oct 28 | 0.307.0 | The consumer froze `ai-reports/catalog.json` at its rewritten state instead of regenerating it, because the walk would list six "bundles" named `2026-03` to `2026-10`. |

Cost: new suite bundles still land flat, so a repo-local fold script must run after each one. Revise and
second-opinion on a filed bundle need a manual move back to the top level first.

## Where the behaviour lives today

- `report_index.go:79-80` (`reportIndexWalk`): "treats every directory directly under ai-reports/ whose name
  does not start with "." as a bundle, whatever its naming style." It reads one `os.ReadDir` level (line 82),
  so `ai-reports/2026-10/` becomes one bundle with no `index.html` and kind `unknown`.
- `report_index.go` supersedes resolution matches by `path.Base` of the folder. This already works with
  nested paths, as long as the walk returns them.
- `report_second_opinion.go:387-388`: `filepath.Dir(relative) != reportsFolder` → "bundle must be a physical
  direct child of ai-reports".
- `skills/do-work-toolbox/actions/ai-report.md:53` (revise, step 1): "`<dir>` is `ai-reports/<folder>` or a bare
  `<folder>` and must be a directory directly under `ai-reports/`; anything else stops with one line."
- `skills/do-work-toolbox/actions/ai-report.md:71` (second-opinion): "It must be a physical direct child of
  `ai-reports/` with a regular `index.html`."

## Proposed direction

Option A (recommended): read-side support only, items 1 to 3. A small change, and repos that never create
month folders see no difference. Option B adds the optional write-side step, which changes where new bundles
land and so touches Collision-Safe Publication and every action that names `ai-reports/<report-slug>/`.

**`report_index.go:79-80`, replace the comment and extend the walk:**

> reportIndexWalk treats every directory directly under ai-reports/ whose name does not start with "." as a
> bundle, whatever its naming style, except a month folder (exactly YYYY-MM): its non-dot child directories
> are bundles instead. Loose files and symbolic links are never bundles.

**`skills/do-work-toolbox/actions/ai-report.md:53`, replace the `<dir>` sentence:**

> `<dir>` is `ai-reports/<folder>`, `ai-reports/<YYYY-MM>/<folder>` or a bare `<folder>`, and must be a bundle
> the catalog lists (directly under `ai-reports/` or inside a month folder); anything else stops with one line.

**`skills/do-work-toolbox/actions/ai-report.md:71`, replace "It must be a physical direct child of `ai-reports/`":**

> It must be a physical bundle directory directly under `ai-reports/` or under a month folder
> `ai-reports/<YYYY-MM>/`, with a regular `index.html`.

## Acceptance check

- In a fixture with `ai-reports/2026-09/2026-09-27_2200_a/index.html` and `ai-reports/2026-10-28_0900_b/index.html`,
  `ai-report index` lists exactly two bundles, `b` first. `catalog.json` paths are
  `ai-reports/2026-10-28_0900_b` and `ai-reports/2026-09/2026-09-27_2200_a`, and no bundle is named `2026-09`.
- `ai-report second-opinion ai-reports/2026-09/2026-09-27_2200_a` passes the bundle check. A symlinked month
  folder is still refused (`validateNoLinkedAncestors`).
- A revise of the nested bundle resolves its supersedes chain across a flat and a nested revision.
- A repo with no month folder produces a byte-identical catalog to v0.307.0.

## Out of scope

- Moving existing bundles: the consumer does that with its own script.
- Any change to bundle naming, catalog fields, or the CMS-side readers a consumer may have.
