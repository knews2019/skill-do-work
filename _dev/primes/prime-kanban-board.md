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

## Lessons

- [REQ-185: separate ordinary optional-tool skips from a counted maintainer-strict behavior lane](../../do-work/archive/REQ-185-javascript-behavior-reachability.md#lessons-learned)
- [REQ-194: retain canonical structured detector evidence and test the source seam directly](../../do-work/archive/UR-043/REQ-194-forward-stray-reqs-through-forensics.md#lessons-learned)
