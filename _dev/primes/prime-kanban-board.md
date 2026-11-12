# Prime: Kanban Board Tool

> Read this before touching `skills/do-work-board/tools/queue-kanban/`. Prime files are
> low noise, high value: they point at the code that is the source of truth, and new
> lessons in this domain get added here so future sessions find them. This file is
> maintainer-side (`_dev/` is export-ignored) — nothing shipped may cite it.

`skills/do-work-board/tools/queue-kanban/` is a standalone Go module (its own `go.mod`, embedded `web/` frontend) that renders the `do-work/` queue as a Kanban board. It ships in the board module and is invoked by `do-work-board board` (`skills/do-work-board/actions/board.md`).

## Conventions

Before changing rendering or clipboard behavior, read [rendering and clipboard conventions](lessons-kanban-board.md#rendering-and-clipboard-conventions).

- **Versioning is folded into the skill.** The tool has no independent changelog — its changes get entries in the root `CHANGELOG.md` and a normal skill version bump, exactly like any action. (It was independently versioned through 1.1.0 before being vendored in; that history lives in `docs/decisions/records/adr-016-*`.)
- **Keep the parser in lock-step with the schema.** The board buckets tickets by the `status` vocabulary in `skills/do-work/actions/work-reference.md`; its parsed display fields must stay aligned with `skills/do-work-board/tools/queue-kanban/model.go`, and the Testing placeholders with `skills/do-work-board/tools/queue-kanban/testing.go`.
- **Write surfaces are counted here.** The tool has exactly three write surfaces, and none touches pipeline state: (1) the board's Testing view writes only the testing placeholders plus `do-work/testers.md`; (2) `queue-kanban next-version` rewrites the single `**Current version**: X.Y.Z` line in one version file (`actions/version.md` by default, `--version-file` to point elsewhere); (3) `queue-kanban next-req` atomically creates one durable number marker under `do-work/.req-reservations/` — queue coordination metadata, not a REQ or pipeline-field edit; stale markers are reaped mechanically by core's `skills/do-work/scripts/cleanup-req-reservations.sh` (SessionStart hook), never by an agent. Everything else the tool does is read-only; the rule is the count, not any list of subcommands. Nothing in the tool writes `status`, any other pipeline field, or `CHANGELOG.md` — the changelog stays an owner-only, human-authored write. Adding a fourth write surface means amending this sentence in the same commit; `_dev/tests/contract-regressions.sh` pins the count to this file. Do not restate it elsewhere; a second copy would drift.
- **`verify` is the mechanical half of the release ritual (`_dev/primes/prime-releases.md`).** It is wired into `skills/do-work/actions/forensics.md`; it reports and routes, while repairs belong to `skills/do-work/actions/cleanup.md`.
- **Board rendering has its own module, but deterministic scalar commands do not depend on it.** Canonical timestamp/frontmatter reads belong to the build-on-demand core CLI; queue-kanban keeps its read-only spellings for compatibility with board callers.
- **A tree fingerprint invalidates parsed inputs, not wall-clock outputs.** A clean live-cache hit may reuse tickets and Markdown, but it must rebuild every payload field derived from the current instant from that cached parsed state, using one fresh instant for the whole response.
- **Never commit build outputs.** The compiled `queue-kanban` binary is gitignored by `skills/do-work-board/tools/queue-kanban/.gitignore`; the `do-work-board static` artifact lands in `build/` at the repo root.
- **Active repository HTML gets a folder origin, never the board origin.** The live `/file` route may render `.html` / `.htm` only by redirecting to a lazily bound loopback preview server rooted at that file's containing directory. Keep the preview read-only, authority-guarded, traversal-safe, and shut down with the board; ordinary text and active formats such as SVG stay inert on `/file`.

- **Browser support:** the strict browser lane targets current stable Chromium. Chrome 141 is deprecated and is not a compatibility target, as confirmed in REQ-375 (Restore the strict browser lane on current Chromium). Record the exact build with browser evidence; a failed trial does not establish behavior across other browser releases.

## Traps

- [family: exact-basename-authority] A filesystem basename is authority before it is presentation: never trim or normalize it before deciding whether it reserves an ID. Canonical writers format the stored minimum-three-digit ID; compatibility readers accept only whole positive-numeric aliases, and documentation or fixtures that construct marker paths count as writers.

## Stakes

Owned by the shipped prime for the same tool: `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` § Stakes (status-driven column bucketing; the embedded client shared by `generate` and `serve`). Read that file's **Read first** and **Traps** with this one — it is the routing index, this file is the maintainer conventions. Not restated here; a second copy would drift.

## Lessons

See [`lessons-kanban-board.md`](lessons-kanban-board.md) — read it before changing what **Read first** or **Traps** name above.
