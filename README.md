# do-work suite

A four-skill task queue suite for agentic coding tools. Capture requests fast, process them later.

## Installation

Run this exact command from the root of the Git repository where you want the suite installed:

```bash
(
  set -e
  project_root="$(git rev-parse --show-toplevel 2>/dev/null)" || { printf 'do-work bootstrap: run this from inside the target Git repository\n' >&2; exit 1; }
  bootstrap_tmp="$(mktemp -d "${TMPDIR:-/tmp}/do-work-suite-bootstrap.XXXXXX")"
  trap 'rm -rf "$bootstrap_tmp"' EXIT
  archive_file="$bootstrap_tmp/do-work-suite.tar.gz"
  curl -fsSL --retry 3 --retry-delay 2 --retry-max-time 60 -o "$archive_file.download" https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz
  mv "$archive_file.download" "$archive_file"
  mkdir -p "$bootstrap_tmp/source"
  tar xzf "$archive_file" -C "$bootstrap_tmp/source" --strip-components=1
  bash "$bootstrap_tmp/source/tools/install-do-work-suite.sh" --project-root "$project_root" --archive "$archive_file"
)
```

One reviewed archive installs four required sibling skills at one shared version:

- `.claude/skills/do-work/` — capture, queue execution, review, and maintenance
- `.claude/skills/do-work-board/` — Kanban, Testing, summaries, and managed Just recipes
- `.claude/skills/do-work-knowledge/` — BKB, memory, dreams, interviews, and prompts
- `.claude/skills/do-work-toolbox/` — audits, reports, presentation, and repository utilities

The installer validates all four modules before the first managed write, asks once, verifies every installed byte before reporting success, and restores every changed managed module/configuration path after failure. It creates or refreshes the sentinel-owned Just section, links the always-on communication-style crew member from a managed section in the project's `CLAUDE.md` (creating the file when absent, touching nothing outside the markers), enables core hooks, leaves memory capture disabled on fresh installs, and preserves existing optional knowledge hooks. Project `do-work/`, `kb/`, application files, and unrelated Just/settings bytes stay outside the managed plan.

Claude Code can invoke the four skill names directly. In Codex or Gemini, point the agent at the appropriate sibling `SKILL.md` once per session, or add those pointers to the project's agent instructions. Commit all four `.claude/skills/do-work*` directories so each repository carries its suite.

**Updating:** `do-work update` and `just run-do-work-update` call the same installed core engine. Both review one archive, reconcile all four module trees plus the managed Just/settings surfaces behind one confirmation, and either verify the complete resulting suite or recover every managed path. Never delete the repository-root `do-work/` queue or `kb/`; neither update entry point manages them.

### Upgrade an existing installation with an AI agent

Paste this into an AI agent from the root of a repository that already has do-work installed:

> Upgrade this repository's existing project-local do-work installation.
>
> Work only inside the current Git repository. Locate its installed four-skill do-work suite, normally under `.claude/skills/`, and confirm the core updater, manifest validator, and full-suite installer are present. If the current suite installation or updater is missing, stop and tell me what is missing instead of attempting a fresh installation.
>
> Run the installed `tools/do-work-update.sh` with this repository's Git root passed through `--project-root`. Use that updater as the only mutation path. Show me its complete managed-file diff and preserve its built-in confirmation before overwriting anything; do not answer the confirmation automatically.
>
> Do not modify the repository's `do-work/` queue, `kb/` data, application files, or unrelated configuration. When finished, verify the installed version and report the previous version, resulting version, and whether the update completed or was cancelled.

This prompt updates only the repository where it is run, so you can paste it separately into each existing installation whenever its owner is ready.

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

See the [Capture Guide](skills/do-work/docs/capture-guide.md) for folder structure, REQ file format, and full workflow.

### Run (process the queue)

When you're ready to build, start the work loop. The assistant triages each request by complexity and works through them one by one:

```
do-work run
```

- **Simple** (config changes, small fixes) — straight to implementation
- **Medium** (clear goal, unknown location) — explore codebase first
- **Complex** (new features, architectural) — plan, explore, then build

Each completed request gets archived with implementation notes and a git commit. A built-in review runs after each item. The build phase always loads behavioral guardrails (`skills/do-work/crew-members/coding-guardrails.md`) — minimal, surgical changes with verifiable success criteria and names that stay findable by plain-text search, not "it compiles" handwaves.

The [Work Guide's trigger aliases](skills/do-work/docs/work-guide.md#trigger-aliases) are the canonical public list. The guide also covers the full work loop, triage routes, and clarify mode.

### Full cycle without persistent state

`do-work run` already performs implementation, testing, and review for each REQ. A full cycle therefore composes capture, verification, run, and presentation directly; it does not need a second testing stage or a resumable state file.

Copy this prompt and replace the final placeholder:

```text
Use the installed do-work suite to complete this request end to end:

1. Use do-work to capture the request below and record the resulting UR ID.
2. Run do-work verify-requests for that UR. Stop and report if verification fails.
3. Run the UR's REQs through do-work run. Require its built-in tests and review to pass.
4. Use do-work-toolbox ai-report for the same UR.
5. Report the implementation, tests, decisions, and deliverable paths.

Request:
<paste request here>
```

## Other actions

Run each sibling's `help` command for its menu. Core guides live in [`skills/do-work/docs/`](./skills/do-work/docs/); extension guides live with their owning skill.

For presentation, choose the artifact explicitly: `do-work-toolbox ai-report REQ-NNN` creates detailed stakeholder HTML for one completed item, `do-work-toolbox present-work all` refreshes the cross-project portfolio, and `do-work-toolbox present-video REQ-NNN` creates a source-only Remotion walkthrough. `showcase`, `visual report`, and `proof of work` are aliases for `ai-report`; `portfolio` and `work portfolio` are aliases for `present-work`; `remotion` and `video walkthrough` are aliases for `present-video`.

For repository comprehension, `do-work-toolbox architecture-report` writes a new dated, immutable bundle under `ai-reports/` holding one self-contained `index.html` and no Markdown companion. It shares the home of `ai-report`, with a freely redesigned architecture view, rendered diagrams, clickable section navigation, and GitHub source links. Each report opens with an authored account of what changed since the previous HTML report; older Markdown bundles remain untouched and are not baselines. `architecture overview` and `map the repo` are aliases for it.

Common extension calls also include `do-work-board board`, `do-work-knowledge bkb`, `do-work-knowledge memory`, `do-work-toolbox code-review`, and `do-work-toolbox inspect`.

### Queue board (`do-work-board board`)

`do-work-board board` builds and runs the board sibling's Go tool, serving the queue at `http://localhost:8090`. `do-work-board static` writes a self-contained snapshot, `do-work-board summary` prints column counts, and `do-work-board cli` prints the in-flight digest. The board is the suite's only Go-toolchain dependency; updates keep it synchronized with core.

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

Yes. Install it in your project root. Installation and updates write only the managed paths you review and confirm: the four sibling skill trees plus the managed Just/settings surfaces. Durable queue state — requests, working records, and archives — lives under `do-work/`. Project source is written only during explicitly invoked REQ implementation, bounded by the declared `## Scope` for Routes B and C or the focused REQ text for Route A.

### What happens if something goes wrong during processing?

The work loop processes one request at a time. If a request fails, it's marked as failed with notes on what went wrong, and the loop moves to the next one. Nothing is lost — you can fix the issue and re-queue. Run `do-work forensics` to diagnose stuck or failed work.

### Can I edit REQ files manually?

Yes. They're plain markdown with frontmatter. You can change priority, edit requirements, or add context before running the queue. The UR folder's `input.md` preserves your original verbatim input regardless of what you change in the REQ files.

### Do I need to remind it to write lessons, keep working, commit often, or not block when I'm AFK?

Mostly no — those are already built in. `do-work run` appends a `## Lessons Learned` section per REQ, logs out-of-scope finds to `## Discovered Tasks`, commits each finished REQ atomically, loads the YAGNI guardrail on every build, and never blocks on ambiguity (it records a best-judgment decision and files a `pending-answers` follow-up you review later via `do-work clarify`). The two things you *can't* get just by asking — an unbounded "loop until the queue is empty" runner and a backgrounded commit — are deliberately not the default. See the [Standing Preferences](skills/do-work/docs/standing-preferences.md) reference for the full map of common nudges → where each already lives.

## Token efficiency

The suite is designed for selective loading — you don't need everything in context at once.

- Each sibling's **SKILL.md** handles only that package's routing and dispatches to the relevant action file.
- **Action files** are loaded on-demand by the routing decision. Only the active action file needs to be in context.
- Core **crew-members/** are JIT-loaded during implementation based on REQ domain. They never need to be pre-loaded.
- Package **docs/** guides are for human reading, not agent context. Don't load them during work.
- Core **specs/** templates are loaded by the work action after triage, only when a REQ matches.

If your agent has limited context, prioritize: **owning sibling SKILL.md → active action file → relevant crew-member**. Everything else is optional.

## Hooks

The suite installer enables one core Claude Code hook:

- **`skills/do-work/hooks/session-start.sh`** — SessionStart hook that injects the installed version and pending REQ count at the beginning of each session, and also writes to your queue files: it reaps stale REQ-number reservation markers (`skills/do-work/scripts/cleanup-req-reservations.sh`) and mechanically repairs detectably wrong `*_at` stamps in `do-work/queue/` and `do-work/working/` (`skills/do-work/scripts/repair-req-timestamps.sh`).

Fresh installs do not enable memory capture. To opt in later, run `do-work-knowledge setup-memory`; it composes the knowledge hook fragment without clobbering existing settings.

The sample commands are anchored to `$CLAUDE_PROJECT_DIR/.claude/skills/do-work/hooks/…` — Claude Code runs hooks from your project root, not the skill directory, so a project-relative `hooks/…` path wouldn't resolve. This assumes do-work lives at the canonical `.claude/skills/do-work/`; if you installed it elsewhere, change the path in your `.claude/settings.json` to match.

## Designed for agentic coding tools

This skill assumes your tool supports:
- File editing and shell access
- Optional subagent or multi-agent workflows (Plan, Explore, Build)
- Git integration for per-request commits (optional)

Originally written for Claude Code. Works with other tools that can read/write files and run shell commands. If your tool does not support subagents, run Plan, Explore, and Implementation phases sequentially in the same session.

## License

MIT
