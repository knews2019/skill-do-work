# do-work

A task queue skill for agentic coding tools. Capture requests fast, process them later.

## Installation

do-work installs into `.claude/skills/do-work/`, so it never touches your project's own files. Run from your repo root:

```bash
mkdir -p .claude/skills/do-work
curl -sL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz \
  | tar xz -C .claude/skills/do-work --strip-components=1 \
      --exclude='_dev' --exclude='do-work' --exclude='ai-reports' \
      --exclude='.vscode' --exclude='decisions'
```

The entry point is `.claude/skills/do-work/SKILL.md`.

- **Claude Code** auto-discovers it — just say `do-work help`.
- **Codex / Gemini** don't auto-discover skills — point the agent at `.claude/skills/do-work/SKILL.md` once per session (or add a one-line pointer to your `AGENTS.md` / `GEMINI.md`).

Commit `.claude/skills/do-work/` to your repo — each repo carries its own copy of the skill.

### Install with an AI agent

Paste this into Claude Code, Codex, or Gemini — it fetches the instructions and does the install for you:

> Install the **do-work** skill into this repository by fetching and following the Installation section here: https://raw.githubusercontent.com/knews2019/skill-do-work/refs/heads/main/README.md — install into `.claude/skills/do-work/`, don't modify anything outside that folder, then confirm `SKILL.md` exists and tell me how to run it.

The prompt stays short because the raw README it points at carries the command, the install location, the verify step, and the per-tool invocation notes.

**Updating:** The cleanest path is `do-work update` — it checks the upstream version, snapshots your install, pre-cleans the globbed `prompts/` and `interviews/` directories, then extracts (see `actions/version.md`). If you update manually by re-running the install command instead, note that `tar` overwrites files in place but does **not** delete files removed upstream. For directories the skill loads by name (`actions/`, `crew-members/`, `specs/`, `docs/`) leftover files are harmless. But `prompts/` and `interviews/` are *globbed* — `do-work prompts list`/`run` and `do-work interview list` enumerate every `*.md` in them — so a prompt or interview removed upstream stays runnable until you delete it. For a guaranteed-clean manual update, delete the whole `.claude/skills/do-work/` folder and re-extract — it's self-contained, so nothing else is affected. Never delete the repo-root `do-work/` runtime directory (your queue, archives, and deliverables).

## The idea

Separate *thinking of things* from *doing things*. You throw ideas at the queue as they come up. When you're ready, you tell the assistant to work. It picks up each request, triages complexity, and builds until the queue is empty.

## Three core workflows

### Capture

Throw tasks at the queue as they come up — one-liners, multi-feature specs, bug reports, screenshots, meeting notes. Each invocation creates a User Request (UR) folder preserving your verbatim input, plus one or more REQ files that enter the queue.

```
do-work capture-request: add dark mode to the settings page
do-work capture-request: the search is slow, also add an export button, and fix the header alignment
do-work capture-request: [paste meeting notes, specs, or a screenshot]
```

The skill splits compound inputs into separate REQ files automatically. It asks clarifying questions during capture (while you're present) but never starts building — capture and execution are strictly separate. For testable behavioral work, capture also infers and confirms the RED case: how we know it's failing or missing now, and what turns GREEN when the work is done.

See the [Capture Guide](docs/capture-guide.md) for folder structure, REQ file format, and full workflow.

### Run (process the queue)

When you're ready to build, start the work loop. The assistant triages each request by complexity and works through them one by one:

```
do-work run
```

- **Simple** (config changes, small fixes) — straight to implementation
- **Medium** (clear goal, unknown location) — explore codebase first
- **Complex** (new features, architectural) — plan, explore, then build

Each completed request gets archived with implementation notes and a git commit. A built-in review runs after each item. The build phase always loads behavioral guardrails (`crew-members/karpathy.md`) — minimal, surgical changes with verifiable success criteria, not "it compiles" handwaves.

Other trigger words: `go`, `start`, `begin`, `process`, `execute`, `build`, `continue`, `resume`.

See the [Work Guide](docs/work-guide.md) for the full pipeline, triage routes, and clarify mode.

### Pipeline (full end-to-end)

One command, full cycle: investigate → capture → verify → run → review → present. Stateful and resumable — if the session ends mid-pipeline (context limit, crash, closed terminal), re-invoking picks up from the last step.

```
do-work pipeline add dark mode to settings    # initialize with a request
do-work pipeline                              # resume the active pipeline
do-work pipeline status                       # progress without advancing
do-work pipeline abandon                      # deactivate without completing
```

Pipeline state lives at `do-work/pipeline.json`. Each step dispatches to an existing action — the pipeline never re-implements logic.

## Other actions

Run `do-work help` for the full menu. Per-action guides live in [`docs/`](./docs/).

Common ones: `verify-requests`, `review-work`, `validate-feedback`, `clarify`, `code-review`, `ui-review`, `quick-wins`, `scan-ideas`, `deep-explore`, `prime`, `present-work`, `commit`, `inspect`, `cleanup`, `forensics`, `roadmap`, `board`, `stray-check`, `bkb`, `dream`, `interview`, `prompts`, `install ui-design`, `install bowser`, `install last30days`, `install just-kanban`, `version`, `update`, `recap`, `tutorial`, `help`.

### Queue board (`do-work board`)

`do-work board` builds and runs a small Go tool (`tools/queue-kanban/`, shipped with the skill) that renders your `do-work/` queue as a live Kanban board + completion calendar. `do-work board` serves it at `http://localhost:8090`; `do-work board static` writes a self-contained HTML snapshot you can hand off; `do-work board summary` prints column counts. It's a read-only viewer and it's the one part of the skill that needs the **Go toolchain**. Because the tool ships inside the skill, `do-work update` keeps it current — no separate install.

## File structure

```
do-work/
├── queue/                    # Pending requests
│   ├── REQ-018-pending.md
│   └── REQ-019-pending.md
├── user-requests/            # Verbatim input + assets per user request
│   └── UR-003/
│       ├── input.md          # Original user input (source of truth)
│       └── assets/
├── working/                  # Currently being processed
│   └── REQ-020-in-progress.md
└── archive/                  # Completed work (self-contained units)
    ├── UR-001/               # UR folder with its completed REQs inside
    │   ├── input.md
    │   └── REQ-013-done.md
    └── REQ-010-legacy.md     # Legacy REQs archive directly
```

## Q&A

### Why separate capture from processing? Why not just build immediately?

Because capture preserves what you asked for, and processing tracks how it was built — and neither interferes with the other. When you say `do-work capture-request: ...`, your exact words are saved in a UR folder as the permanent source of truth. Nothing is paraphrased, nothing is lost. When `do-work run` picks up that request later, the REQ file tracks every decision: what was planned, what was explored, what was built, what was reviewed. You end up with a clear trail from intent to implementation — what the user wanted, what the builder decided, and why. Without this separation, Claude tends to hear your request, immediately start coding, and leave no trace of what was asked or how decisions were made. The two-phase split means capture is fast and cheap (dump ideas anytime), processing is thorough and auditable (every change is traceable back to a request).

### Why not just let Claude decide what to do?

Claude *does* decide — the skill just raises the floor. Without structure, Claude picks different steps every time. Sometimes it plans, sometimes it dives in. Sometimes it reviews its work, sometimes it ships the first thing that compiles. The skill encodes lessons learned from many sessions into a repeatable baseline: plan before building, review after building, don't lose the original input along the way. Claude still makes all the real decisions within each step. The skill makes sure those decisions happen.

### How is this different from a hardcoded CI pipeline?

It's not a fixed sequence. The triage system (simple/medium/complex) means Claude chooses how much planning each request needs. Simple config changes skip straight to implementation. Complex features get explore, plan, then build. The skill is more like a senior dev's checklist — you still use judgment, but you don't skip steps because you felt confident.

### Do I need Claude Code specifically?

No. The skill works with any agentic coding tool that can read/write files and run shell commands. It was written for Claude Code but the action files are standalone prompts — paste them into any chat interface and they work. Subagent support is a nice-to-have, not a requirement.

### Can I use this with an existing project?

Yes. Install it in your project root. The skill only creates files inside `do-work/` — it doesn't touch your source code structure. Your codebase is read during the build phase, but all skill state (requests, archives, deliverables) lives in `do-work/`.

### What happens if something goes wrong during processing?

The work loop processes one request at a time. If a request fails, it's marked as failed with notes on what went wrong, and the loop moves to the next one. Nothing is lost — you can fix the issue and re-queue. Run `do-work forensics` to diagnose stuck or failed work.

### Can I edit REQ files manually?

Yes. They're plain markdown with frontmatter. You can change priority, edit requirements, or add context before running the queue. The UR folder's `input.md` preserves your original verbatim input regardless of what you change in the REQ files.

## Token efficiency

The skill is designed for selective loading — you don't need everything in context at once.

- **SKILL.md** is the only file loaded initially. It handles routing and dispatches to the relevant action file.
- **Action files** are loaded on-demand by the routing decision. Only the active action file needs to be in context.
- **crew-members/** are JIT-loaded during implementation based on REQ domain. They never need to be pre-loaded.
- **docs/** guides are for human reading, not agent context. Don't load them during work.
- **specs/** templates are loaded by the work action after triage, only when a REQ matches.

If your agent has limited context, prioritize: **SKILL.md → active action file → relevant crew-member**. Everything else is optional.

## Hooks (optional)

Two optional hook scripts for Claude Code users:

- **`hooks/pipeline-guard.sh`** — Stop hook that prevents the agent from stopping mid-pipeline. Install as a `Stop` hook.
- **`hooks/session-start.sh`** — SessionStart hook that injects a status line (version, pending REQs, active pipeline) at the beginning of each session.

To install, merge the hook config from `hooks/hooks.json` into your `.claude/settings.json`. See each script for details.

The sample commands are anchored to `$CLAUDE_PROJECT_DIR/.claude/skills/do-work/hooks/…` — Claude Code runs hooks from your project root, not the skill directory, so a project-relative `hooks/…` path wouldn't resolve. This assumes do-work lives at the canonical `.claude/skills/do-work/`; if you installed it elsewhere, change the path in your `.claude/settings.json` to match.

## Designed for agentic coding tools

This skill assumes your tool supports:
- File editing and shell access
- Optional subagent or multi-agent workflows (Plan, Explore, Build)
- Git integration for per-request commits (optional)

Originally written for Claude Code. Works with other tools that can read/write files and run shell commands. If your tool does not support subagents, run Plan, Explore, and Implementation phases sequentially in the same session.

## License

MIT
