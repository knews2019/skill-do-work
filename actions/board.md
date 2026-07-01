# Board Action

> **Part of the do-work skill.** Builds and runs the shipped `queue-kanban` Go tool to render this repo's `do-work/` queue as a Kanban board + completion calendar. Invoked by `do-work board` / `do-work kanban`.

**Read-only.** The board only *reads* the `do-work/` Markdown tree — it never writes to the queue, claims REQs, or changes state. The one thing it writes is the compiled binary (gitignored) and, in `static` mode, a throwaway HTML artifact under `build/`.

The tool is a standalone Go module that ships inside the skill at `tools/queue-kanban/` (its module, `go.mod`, and embedded `web/` frontend). It rides do-work version bumps, so `do-work update` carries the latest board into every repo. Because it's compiled, this action needs the **Go toolchain** — the one action that does. It degrades gracefully when Go is absent: it reports and stops, never blocking the rest of the skill.

## When to Use

**Use when:**
- The user says "board", "kanban", "show the queue", "queue board", or "visualize the queue".
- The user wants a live, auto-refreshing view of pending/claimed/blocked/recently-done REQs.
- The user wants a shareable static HTML snapshot of queue state (`static` mode).
- The user wants quick column counts without a browser (`summary` mode).

**Do NOT use when:**
- The user wants a text roadmap or dependency rollup → `actions/roadmap.md`.
- The user wants to *understand* uncommitted changes or REQ contents → `actions/inspect.md`.
- The user wants to process the queue (build the work) → `actions/work.md`.

## Input

`$ARGUMENTS` selects the mode (default = `serve`):

| Token | Mode | Effect |
| --- | --- | --- |
| _(empty)_, `serve`, `live` | serve | Live board at `http://localhost:8090` (re-walks the tree per request). |
| `static`, `generate`, `html` | generate | Self-contained static board written to `build/queue-kanban-board/` (opens from `file://`, zero network). |
| `summary`, `status`, `counts` | summary | Prints column counts to the terminal — no browser. |

An optional trailing `--port N` (serve) or `--out DIR` (static) overrides the default; pass it straight through to the tool.

## Steps

### Step 1: Locate the tool

The skill root is the directory containing `SKILL.md` (this action lives in its `actions/` subdir). The tool is at `<skill-root>/tools/queue-kanban/`. If that directory is missing, report: "queue-kanban tool not found — re-run `do-work update` to fetch it," and stop.

### Step 2: Precondition — Go toolchain

Run `go version`. If `go` is not on `PATH`, stop with:

```
The board needs the Go toolchain (see tools/queue-kanban/go.mod for the required version).
Install it from https://go.dev/dl/ then re-run `do-work board`.
```

Do not attempt to install Go, and do not block any other do-work action — this is the only action with a toolchain dependency.

### Step 3: Resolve the queue's repo root

Run `git rev-parse --show-toplevel` to get the consuming repo root (where `do-work/` lives). Call it `REPO_ROOT`. The tool can auto-discover this by walking up for `do-work/`, but pass `--repo-root "$REPO_ROOT"` explicitly so discovery is deterministic regardless of where the binary sits (it's nested under `.claude/skills/do-work/`). If there is no `do-work/` at `REPO_ROOT`, report that the queue is empty/missing and stop.

### Step 4: Build

Always rebuild (a `do-work update` can leave a stale binary in place):

```bash
cd "<skill-root>/tools/queue-kanban" && go build -o queue-kanban .
```

The first build on a machine whose Go module cache lacks the deps fetches `goldmark` + `yaml.v3` from the module proxy — this needs network once. If the build fails on a Go-version mismatch, surface the required version from `go.mod` and the install link from Step 2.

### Step 5: Run the selected mode

From `<skill-root>/tools/queue-kanban`:

- **serve** — `./queue-kanban serve --repo-root "$REPO_ROOT"` (honor `QUEUE_KANBAN_PORT` or a passed `--port`). Tell the user the URL (`http://localhost:8090` by default) and that it's a long-running process — stop it with Ctrl-C. Run it in the background if your environment supports it, so the session isn't blocked.
- **static** — `./queue-kanban generate --out "$REPO_ROOT/build/queue-kanban-board" --repo-root "$REPO_ROOT"`, then point the user at `build/queue-kanban-board/index.html`. This artifact is a throwaway — mention it's safe to delete or gitignore.
- **summary** — `./queue-kanban summary --repo-root "$REPO_ROOT"` and relay the printed counts.

## Output Format

- **serve:** the live URL + how to stop it.
- **static:** the path to `index.html` and a one-line column-count recap.
- **summary:** the tool's column-count block (total REQs, pending, claimed, needs-input/blocked, recently-done, calendar entries, dependency edges).

## Rules

- Never edit the `do-work/` queue from this action — it is strictly a viewer.
- Never commit the compiled `queue-kanban` binary (the tool's nested `.gitignore` already excludes it) or the generated `build/queue-kanban-board/` artifact.
- Pass `--repo-root` explicitly — do not rely on CWD-based discovery.
- Do not vendor or modify the Go source to "make it build" — a build failure is a toolchain/environment issue to report, not a code change.
- If you change the tool's parser, keep it in lock-step with `actions/work-reference.md`'s Schema Read Contract (the `status`/`depends_on`/`domain` vocabularies the board buckets on).

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
| --- | --- | --- |
| "Go isn't installed, I'll rewrite the board in shell/JS" | Report the missing toolchain per Step 2 and stop | The tool is the shipped, tested renderer; a one-off reimplementation drifts from the schema and misleads viewers |
| "I'll skip the rebuild, the binary's already there" | Always `go build` first | A `do-work update` overwrites the source but leaves a stale binary — running it renders old logic |
| "I'll just run it from the current directory" | Pass `--repo-root "$(git rev-parse --show-toplevel)"` | From a subdir or the nested skill path, CWD discovery can find the wrong `do-work/` or none |

## Red Flags

- The board renders zero tickets against a repo that clearly has REQs → wrong `--repo-root`, or a `status`-vocabulary drift from the Schema Read Contract.
- A tracked `queue-kanban` binary or `build/queue-kanban-board/` shows up in `git status` → the gitignore contract was bypassed.
- The action blocked another do-work command because Go was missing → the graceful-exit in Step 2 was skipped.

## Verification Checklist

- [ ] `go version` checked before any build; missing Go reported, not worked around.
- [ ] Built fresh via `go build -o queue-kanban .` inside `tools/queue-kanban/`.
- [ ] `--repo-root` resolved from `git rev-parse --show-toplevel` and passed explicitly.
- [ ] Correct mode dispatched (serve / static / summary) with the user told the URL, artifact path, or counts.
- [ ] No binary or generated artifact staged or committed.
