# Gallery handoff contract

Both builders start from main revision `f3840b087913d8501267dda1c1395580304f2ba4`. The coordinator assigns actual worktrees and operative refs at dispatch. Keep the two worktrees independent and integrate REQ-745 before REQ-746 so the second merge can verify the live and static views together.

## Routes and shared files

| Surface | REQ-745 | REQ-746 |
| --- | --- | --- |
| `report_gallery.go` | Add local serving mode and use the existing renderer. Own that dispatch branch. | Extend explicit named lookup with generated sources. Own that resolver branch. |
| `report_gallery_page.go` | Add a live feedback slot and its own bindings. Keep its data and script separate from the clipboard hooks. | Add starter fields, prompt details, and copy targets. Reuse the `data-copy` behavior. |
| Curated gallery | Reuse `preview-manifest.json` rows and real report paths for feedback. The live route renders from current data. | Reuse worked cards and their stable names, fragments and preview links. Restore access to shipped sample starters. Own regeneration of the checked-in `docs/report-gallery/index.html`. |
| Installed consumer | Exercise the installed CLI and ordinary gallery as well as the maintainer Just path. Keep existing static gallery behavior. | Exercise the installed action and CLI against shipped, project and generated sources. Preserve installed-sample discovery and automatic ranking. |

Keep `data-gallery-card`, `data-gallery-name`, and `data-gallery-preview` intact for the existing curated gate. Separate new element IDs and data attributes for feedback and starter copies. Each builder should report the exact template or dispatcher hunk and its anchor when the integrator will have to land it. The integrator applies that hunk inside the responsible merge and regenerates the derived page once if reconciliation changes its static bytes. The preview manifest, captures and frozen baseline/climb/final receipts remain evidence.

## Runtime records

- Feedback uses `do-work/report-gallery-feedback.json` with a versioned entries collection. Coordinate REQ-747 around `target_kind`, `report_path`, optional `figure_id`, `mark` (`like` or `dislike`) and optional `note`. Use one project-relative identity per real report page, including the month folder where present, plus a stable figure identity scoped to that page. The report entry and its figure entries are independently writable. Keep fixture marks in disposable projects.
- Generated design authoring belongs to REQ-746 at `ai-reports/<bundle>/gallery-design.md`, including the optional month directory. Reuse the familiar **Borrow** shape for design moves. Reserve the `report:` namespace for this source class. The published starter must carry the very token the resolver accepts. The starter's target remains the new task; the named source contributes design only. Existing cataloged pages without a sidecar need a supported read-only fallback for the design they actually use. Handle ambiguous identities rather than silently picking a different report.

## Real evidence candidate

Use the phase-by-phase worked report as a shared acceptance fixture. The coordinator selected the separately verified September 6 test-efficiency baseline for the user trial. The report is `ai-reports/2026-11-11_1411_commit-22f82fed-work-procedures-load-by-phase/index.html`, the curated name is `phase-by-phase`, and the source change is `22f82fedff49bb1a2fff59a750a496acca746fd2`. The report includes the real procedure-split diagram under `#how` and replay evidence under `#proof`. Its containing folder is the existing bundle. Check its claims against the source change and existing evidence before creating variations. REQ-747 and REQ-748 own generation and the pause for genuine first-round feedback.
