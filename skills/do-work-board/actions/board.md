# Board Action

> **Part of the do-work-board skill.** Builds and runs the shipped `queue-kanban` Go tool to render this repo's `do-work/` queue as an interactive HTML board — one page whose view switcher covers the queue's current state, its history, and its timing. Invoked by `do-work-board board` / `do-work-board kanban`.

**Read-only toward the work pipeline.** The board never writes the pipeline's state — it never changes `status`, claims REQs, or moves files. It writes exactly three things: the compiled binary (gitignored); in `static` mode, a throwaway HTML artifact under `build/` (kept out of `git status` via a one-line `.git/info/exclude` entry — see Step 5); and, from the **Testing view** in `serve` mode only, the testing-track placeholders — the `testing_status` / `tested_by` / `testing_updated_at` / `testing_feedback` frontmatter fields of a REQ plus the `do-work/testers.md` profile list (see "Testing view" below). Those testing writes are the point: the Markdown files are the database of who tested what, with git as the history — there is deliberately no locking or concurrency control, because every write lands in the working tree where it can be reviewed and committed.

The tool is a standalone Go module that ships inside the skill at `tools/queue-kanban/` (its module, `go.mod`, and embedded `web/` frontend). It rides do-work version bumps, so `do-work update` from the core sibling carries the latest board into every repo. Because it is compiled, this action needs the **Go toolchain**, as does the core deterministic command platform. It degrades gracefully when Go is absent: it reports and stops, never blocking unrelated natural-language work.

## When to Use

**Use when:**
- The user says "board", "kanban", "show the queue", "queue board", or "visualize the queue".
- The user wants a live board of pending/claimed/blocked/recently-done REQs (serve mode rebuilds from disk on every browser reload — refresh the page to see new state; it does not push updates to an open tab).
- The user wants to track **who tested which finished REQ** — the board's Testing view (serve mode; reachable from the board's view switcher) lets a tester pick their profile, select a finished REQ to test, and mark it in-testing / tested / returned-with-feedback.
- The user wants a shareable static HTML snapshot of queue state (`static` mode).
- The user wants quick column counts without a browser (`summary` mode).
- The user asks what's in flight / in progress / blocked right now and wants it in the terminal, not a browser tab (`open-work` mode — open count, claimed REQ titles, the status parking each needs-input REQ).
- The user wants the board-owned cross-file invariant probes run (`verify` mode).

**Do NOT use when:**
- The user wants a text roadmap or dependency rollup → `../../do-work/actions/roadmap.md`.
- The user wants to *understand* uncommitted changes or REQ contents → `../../do-work-toolbox/actions/inspect.md`.
- The user wants to process the queue (build the work) → `../../do-work/actions/work.md`.

## Input

`$ARGUMENTS` selects the mode (default = `serve`):

| Token | Mode | Effect |
| --- | --- | --- |
| _(empty)_, `serve`, `live` | serve | Live board at `http://localhost:8090` (re-walks the tree per request). |
| `static`, `generate`, `html` | generate | Self-contained static board written to `build/queue-kanban-board/` (opens from `file://`, zero network). |
| `summary`, `status`, `counts` | summary | Prints column counts to the terminal — no browser. |
| `cli`, `open`, `open-work`, `in-flight` | open-work | Prints the open-work digest to the terminal — open count, claimed REQs with titles, needs-input/blocked REQs with their statuses. No browser. |
| `verify`, `check invariants`, `probes` | verify | Builds the tool and runs its read-only invariant probes; exit 1 means findings, not an execution error. See `../../do-work/actions/forensics.md` Check 14 for the canonical probe descriptions. |

An optional trailing `--port N` (serve) or `--out DIR` (static) overrides the default; pass it straight through to the tool.

## Steps

### Step 1: Locate the tool

The skill root is the directory containing `SKILL.md` (this action lives in its `actions/` subdir). The tool is at `<skill-root>/tools/queue-kanban/`. If that directory is missing, report: "queue-kanban tool not found — re-run `do-work update` from the core sibling to fetch it," and stop.

### Step 2: Precondition — Go toolchain

Run `go version`. If `go` is not on `PATH`, stop with:

```
The board needs the Go toolchain (see tools/queue-kanban/go.mod for the required version).
Install it from https://go.dev/dl/ then re-run `do-work-board board`.
```

Do not attempt to install Go. Report the missing prerequisite for this invocation and stop instead of substituting prose for the board command.

### Step 3: Resolve the queue's repo root

Resolve the consuming project root (where `do-work/` lives) with the repo-standard fallback — `git` is optional for the consuming project, matching `justfile.template` and `../../do-work/actions/version.md`:

```bash
REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
```

In a non-git project the invocation directory is the root, so run this from the project root. Pass `--repo-root "$REPO_ROOT"` explicitly so discovery is deterministic regardless of where the binary sits (it's nested under `.claude/skills/do-work-board/`); the tool's own walk-up discovery (`resolveRepoRoot`) is the last resort, not the default. If there is no `do-work/` at `REPO_ROOT`, report that the queue is empty/missing and stop.

### Step 4: Build

Always rebuild (a `do-work update` from the core sibling can leave a stale binary in place):

```bash
cd "<skill-root>/tools/queue-kanban" && go build -o queue-kanban .
```

The first build on a machine whose Go module cache lacks the deps fetches `goldmark` + `yaml.v3` from the module proxy — this needs network once. If the build fails on a Go-version mismatch, surface the required version from `go.mod` and the install link from Step 2.

### Step 5: Run the selected mode

From `<skill-root>/tools/queue-kanban`:

- **serve** — `./queue-kanban serve --repo-root "$REPO_ROOT"` (honor `QUEUE_KANBAN_PORT` or a passed `--port`). Tell the user the URL (`http://localhost:8090` by default), that reloading the page refreshes the data (the server re-walks the tree per request; it does not push updates), and that it's a long-running process — stop it with Ctrl-C. Run it in the background if your environment supports it, so the session isn't blocked.
- **static** — `./queue-kanban generate --out "$REPO_ROOT/build/queue-kanban-board" --repo-root "$REPO_ROOT"`, then point the user at `build/queue-kanban-board/index.html`. This artifact is a throwaway — mention it's safe to delete. After generating to the **default** `--out` (a user-chosen `--out` is theirs to manage), add a local git exclude so the snapshot never sits in `git status` as untracked noise:

  ```bash
  (cd "$(git rev-parse --show-toplevel 2>/dev/null || pwd)" && <skill-root>/../do-work/scripts/add-local-git-exclude.sh build/queue-kanban-board/index.html '/build/queue-kanban-board/')
  ```

  This generated snapshot uses the canonical [Local Git ignore](../../do-work/docs/prescribed-shell-primitives.md#local-git-ignore) primitive with a project-root-specific pattern. The local exclude preserves the action's read-only contract for tracked project files, the guard is idempotent, and a non-git project skips it silently.
- **summary** — `./queue-kanban summary --repo-root "$REPO_ROOT"` and relay the printed counts.
- **open-work** — `./queue-kanban open-work --repo-root "$REPO_ROOT"` and relay the printed digest. Terminal-resolved REQs never appear in it (recently-done is history, not open work — the calendar is a different surface and does carry open work), and parse warnings arrive as a count pointing at `summary` — if the user wants those details, that's the summary mode, not a second command here.
- **verify** — `./queue-kanban verify --repo-root "$REPO_ROOT"`. Relay its findings; exit 1 is the expected findings status, not a launcher failure. The probe set and each invariant's meaning are owned by `../../do-work/actions/forensics.md` Check 14.

**NEEDS INPUT · BLOCKED is the operator-actionable inbox.** A bare `status: blocked` card with any unmet `depends_on` target is not actionable yet, so the board displays it under PENDING → Waiting on dependencies while retaining its blocked status, blocked-by badge, and dependency chips. Once every dependency reaches terminal success, the card moves back to NEEDS INPUT · BLOCKED because its external blocked condition is then the sole gate. `pending-answers`, `pending-heavy-testing`, `blocked-archive-collision`, `blocked-dependency-cycle`, `failed`, and unrecognized statuses keep their existing Needs-input placement regardless of dependencies; their questions, heavy-test permission, repair actions, or visibility warnings remain operator-actionable.

A pending or claimed card carries an `overlaps` badge naming the other pending/claimed REQs whose declared `write_set` could touch the same files — an informational heads-up, never a block (the detail drawer shows the card's own write set and links each contending REQ). The badge schedules nothing at any builder count: under fan-out the declared set is advisory input to the human's pick and the merge is the non-interference proof, never the badge (`../../do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). It just surfaces declared file contention for a human reading the board. No badge means no *declared* overlap: on a REQ that never declared a `write_set` that reads as unknown, not safe. Globs are matched with `path.Match` semantics (OS-independent, `/`-separated): `*` never crosses `/`, `**` is not recursive, and a malformed pattern is treated as no-match for that direction — but literal equality short-circuits first, so two REQs declaring the *identical* malformed pattern still badge each other. The badge can therefore miss real contention; illustrative miss-classes, not a closed list: glob-vs-glob, `**`, a malformed pattern against anything but its own twin, and an entry naming a directory (`actions/` never badges against `actions/board.md` — `path.Match` is false both for `actions/` and for `actions` against it).

### Step 6: Testing view (serve mode, in-browser — nothing more for the agent to run)

The served board's **Testing** view — the only one that writes anything, and the only one with a server API behind it — tracks who tested which finished REQ. It shows every terminal-success REQ (`completed` / `completed-with-issues` — plus any REQ that already carries a testing record, so records never vanish) in four columns: **Ready to test → In testing → Returned with feedback → Tested**. The user picks (or adds) a tester profile in the view's toolbar, then drives per-card actions; the browser POSTs to the live server's `/api/testing/*` endpoints, which write the record into the Markdown itself:

- REQ frontmatter placeholders: `testing_status: in-testing | tested | returned`, `tested_by`, `testing_updated_at`, and (while returned) `testing_feedback` — see `../../do-work/actions/work-reference.md`'s Request File Schema.
- Tester profiles: one `- Name` bullet per profile in `do-work/testers.md` (created on first use; hand-editable).

The main Board view shows a `testing` badge on any card carrying a record, so testing state is visible without switching views. In `static` mode the Testing view renders read-only (no server, no actions). There is no locking: changes land in the working tree and git is the audit trail — when the user asks "who tested REQ-NNN?", the frontmatter (or `git log` on the REQ file) answers.

**Standing shortcut:** the full-suite installer publishes one managed flat recipe surface. Run `just --list` for its live inventory: it includes the four board recipes, direct deterministic core/knowledge/toolbox commands, canonical `just do-work-update`, and compatibility `just run-do-work-update`. Re-running the installer on a project whose managed block drifted offers a diff-and-consent upgrade. One difference: `just run-kanban` auto-opens your default browser at the board URL (a user-initiated shortcut, not an agent action); this action's serve mode (Step 5) never does.

## Output Format

- **serve:** the live URL + how to stop it.
- **static:** the path to `index.html` and a one-line column-count recap.
- **open-work:** the tool's digest — the open total with its pending (ready/waiting) · claimed · needs-input/blocked breakdown, then each claimed REQ as id + title, then each needs-input/blocked REQ as id + status + title (with the `blocked_by` condition when one is named), then a warnings count when the parse raised any.
- **summary:** the tool's column-count block (total REQs, pending — split into ready-to-work and waiting-on-deps — claimed, needs-input/blocked, recently-done, completion anomalies — with the offending REQ ids listed when nonzero — calendar entries (one per REQ, so this equals the total), dependency edges).

## Rules

- Never edit the `do-work/` queue from this action — the agent is strictly a launcher/viewer. The one sanctioned write path is the served Testing view's own `/api/testing/*` endpoints (user-driven, testing placeholders + `do-work/testers.md` only); never write `status` or any other pipeline field, and never hand-edit testing placeholders on the user's behalf from this action.
- Never commit the compiled `queue-kanban` binary (the tool's nested `.gitignore` already excludes it) or the generated `build/queue-kanban-board/` artifact.
- Pass `--repo-root` explicitly (resolved via `git rev-parse --show-toplevel 2>/dev/null || pwd`) — the tool's CWD walk-up is the non-git last resort, not the default.
- Do not vendor or modify the Go source to "make it build" — a build failure is a toolchain/environment issue to report, not a code change.
- If you change the tool's parser, keep it in lock-step with `../../do-work/actions/work-reference.md`'s Schema Read Contract — the `status` vocabulary drives column bucketing, and `depends_on` drives the Ready/Waiting presentation. Pending REQs partition between those groups; exactly a bare `status: blocked` REQ with an unmet dependency also joins Pending → Waiting until the upstream work completes, without changing its on-disk status. A dependency counts as met only when its target reached `completed` or `completed-with-issues`; `cancelled` never satisfies gating, and a `depends_on` id that names no REQ in the tree is treated as unmet **and** raised as a data warning. `domain`, `error`, `error_type`, and the blocked fields (`blocked_by`/`blocked_at`/`blocked_check`) are parsed for display only — failure details appear in the drawer, with `error` read as verbatim scalar text and a present `error_type` normalized under the schema (invalid values retain provenance and raise a data warning); an absent `error_type` stays absent and never fabricates the `code` default. The blocked fields drive the "blocked by" badge and drawer rows on a `status: blocked` card but never route it: routing reads the normalized status plus the already-derived unmet-dependency set. `write_set` (the REQ's declared write surface — display-only at any builder count, per `../../do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch) is read verbatim into the card payload — no normalization, no alias — and feeds exactly one derived, display-only annotation: `annotateWriteSetOverlap` in `tools/queue-kanban/model.go` runs *after* column bucketing and lists, per pending/claimed card, the other pending/claimed REQ ids whose declared set could touch the same files. It drives the `overlaps` badge and a drawer row and nothing else — never column logic, never blocking. Nothing schedules on it: the board just surfaces declared file contention for a human reading it. The testing placeholders (`testing_status` and friends) are the board's own vocabulary — also defined in the Schema Read Contract and mirrored in `tools/queue-kanban/testing.go`; an off-vocabulary `testing_status` renders as not-tested with an invalid flag plus a data warning, never silently.
- `do-work/notes.md` renders as a collapsible Notes strip above the columns (visible in both the board and calendar views), mirroring how `../../do-work/actions/roadmap.md` surfaces it. Notes are plain text, never Markdown, and never tickets: they get no column, no calendar entry, and no detail drawer. The board only reads the file — `do-work-toolbox note` is still the only writer, and the user still deletes lines by hand.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
| --- | --- | --- |
| "Go isn't installed, I'll rewrite the board in shell/JS" | Report the missing toolchain per Step 2 and stop | The tool is the shipped, tested renderer; a one-off reimplementation drifts from the schema and misleads viewers |
| "I'll skip the rebuild, the binary's already there" | Always `go build` first | A `do-work update` from the core sibling overwrites the source but leaves a stale binary — running it renders old logic |
| "I'll just run it from the current directory" | Pass `--repo-root "$(git rev-parse --show-toplevel 2>/dev/null || pwd)"` | From a subdir or the nested skill path, CWD discovery can find the wrong `do-work/` or none |

## Red Flags

- The board renders zero tickets against a repo that clearly has REQs → wrong `--repo-root`, or a `status`-vocabulary drift from the Schema Read Contract.
- A tracked `queue-kanban` binary or `build/queue-kanban-board/` shows up in `git status` → the gitignore contract was bypassed.
- The action blocked another do-work command because Go was missing → the graceful-exit in Step 2 was skipped.

## Verification Checklist

- [ ] `go version` checked before any build; missing Go reported, not worked around.
- [ ] Built fresh via `go build -o queue-kanban .` inside `tools/queue-kanban/`.
- [ ] `--repo-root` resolved from `git rev-parse --show-toplevel 2>/dev/null || pwd` and passed explicitly.
- [ ] Correct mode dispatched (serve / static / summary / open-work) with the user told the URL, artifact path, counts, or digest.
- [ ] Static mode with the default `--out`: `build/queue-kanban-board/` no longer appears untracked in `git status` (the info/exclude entry landed, or was already covered).
- [ ] No binary or generated artifact staged or committed.
