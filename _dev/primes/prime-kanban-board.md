# Prime: Kanban Board Tool

> Read this before touching `skills/do-work-board/tools/queue-kanban/`. Prime files are
> low noise, high value: they point at the code that is the source of truth, and new
> lessons in this domain get added here so future sessions find them. This file is
> maintainer-side (`_dev/` is export-ignored) — nothing shipped may cite it.

`skills/do-work-board/tools/queue-kanban/` is a standalone Go module (its own `go.mod`, embedded `web/` frontend) that renders the `do-work/` queue as a Kanban board. It ships in the board module and is invoked by `do-work-board board` (`skills/do-work-board/actions/board.md`).

## Conventions

- **Versioning is folded into the skill.** The tool has no independent changelog — its changes get entries in the root `CHANGELOG.md` and a normal skill version bump, exactly like any action. (It was independently versioned through 1.1.0 before being vendored in; that history lives in `decisions/records/adr-016-*`.)
- **Keep the parser in lock-step with the schema.** The board buckets tickets by the `status` vocabulary in `skills/do-work/actions/work-reference.md`; its parsed display fields must stay aligned with `skills/do-work-board/tools/queue-kanban/model.go`, and the Testing placeholders with `skills/do-work-board/tools/queue-kanban/testing.go`.
- **Write surfaces are counted, and the count lives in root `CLAUDE.md`.** The canonical sentence — the tool has exactly three write surfaces, none touching pipeline state — lives in `CLAUDE.md` § Kanban Board Write Surfaces, and adding a surface means amending that sentence in the same commit. Do not restate the count here; a second copy would drift.
- **`verify` is the mechanical half of "Before Every Commit."** It is wired into `skills/do-work/actions/forensics.md`; it reports and routes, while repairs belong to `skills/do-work/actions/cleanup.md`.
- **Toolchain exception to "design for the floor."** The board is the one capability that needs Go (`skills/do-work-board/tools/queue-kanban/go.mod`); `skills/do-work-board/actions/board.md` degrades gracefully when it is absent. Core may use an already-built sibling binary only where a shell-portable fallback remains documented.
- **Never commit build outputs.** The compiled `queue-kanban` binary is gitignored by `skills/do-work-board/tools/queue-kanban/.gitignore`; the `do-work-board static` artifact lands in `build/` at the repo root.
- **Active repository HTML gets a folder origin, never the board origin.** The live `/file` route may render `.html` / `.htm` only by redirecting to a lazily bound loopback preview server rooted at that file's containing directory. Keep the preview read-only, authority-guarded, traversal-safe, and shut down with the board; ordinary text and active formats such as SVG stay inert on `/file`.

## Lessons

- [REQ-185: separate ordinary optional-tool skips from a counted maintainer-strict behavior lane](../../do-work/archive/UR-041/REQ-185-javascript-behavior-reachability.md#lessons-learned)
- [REQ-194: retain canonical structured detector evidence and test the source seam directly](../../do-work/archive/UR-043/REQ-194-forward-stray-reqs-through-forensics.md#lessons-learned)
- [REQ-195: separate raw-marker uniqueness, canonical placement, and post-assembly absence](../../do-work/archive/UR-044/REQ-195-modularize-framework-free-board-client.md#lessons-learned)
- [REQ-200: allowlist byte-detected inline formats without weakening the inert-text fallback](../../do-work/archive/UR-045/REQ-200-render-png-file-mentions-as-images.md#lessons-learned)
- [REQ-207: isolate active HTML by folder origin, capture addenda before implementation, and commit before hand-back](../../do-work/archive/UR-046/REQ-207-render-html-file-mentions-as-folder-aware-previews.md#lessons-learned)
- [REQ-219: ship a rule's verdict in the payload so a second reader cannot become a second definition](../../do-work/archive/UR-050/REQ-219-board-durations-view.md#lessons-learned)
- [REQ-226: decide a chart's label placement once, on the Go side, and verify the result by rendering it](../../do-work/archive/REQ-226-stop-durations-chart-overprinting-and-clipping.md#lessons-learned)
- [REQ-227: a view that binds to persistent nodes owns its listener teardown; draw in pixels so a zoom has no scale to invalidate](../../do-work/archive/REQ-227-timeline-view-with-two-segment-req-bars.md#lessons-learned)
