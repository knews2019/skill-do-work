# do-work

A task queue skill for agentic coding tools. Capture requests fast, process them later.

## Installation

```bash
# Run from the directory where you want the skill installed
curl -sL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz | tar xz --strip-components=1 --exclude='_dev'
```

**Updating:** Re-run the same command to update. Note that tar extraction overwrites but does not delete files removed upstream — stale files from older versions may linger (generally harmless). For a fully clean update, delete only the known skill paths (`actions/`, `crew-members/`, `SKILL.md`, `CHANGELOG.md`, `README.md`) before re-extracting — never delete `do-work/` or other project files.

## The idea

Separate *thinking of things* from *doing things*. You throw ideas at the queue as they come up. When you're ready, you tell the assistant to work. It picks up each request, triages complexity, and builds until the queue is empty.

## Usage scenarios

### 1. Capture requests

Throw tasks at the queue as they come up — one-liners, multi-feature specs, bug reports, screenshots, meeting notes. Each invocation creates a User Request (UR) folder preserving your verbatim input, plus one or more REQ files that enter the queue.

```
do work capture request: add dark mode to the settings page
do work capture request: the search is slow, also add an export button, and fix the header alignment
do work capture request: [paste meeting notes, specs, or a screenshot]
```

The skill splits compound inputs into separate REQ files automatically. It asks clarifying questions during capture (while you're present) but never starts building — capture and execution are strictly separate.
For testable behavioral work, capture should also infer and confirm the RED case: how we know it's failing or missing now, and what turns GREEN when the work is done.

See the [Capture Guide](docs/capture-guide.md) for folder structure, REQ file format, and full workflow.

### 2. Process the queue

When you're ready to build, start the work loop. The assistant triages each request by complexity and works through them one by one:

```
do work run
```

- **Simple** (config changes, small fixes) — straight to implementation
- **Medium** (clear goal, unknown location) — explore codebase first
- **Complex** (new features, architectural) — plan, explore, then build

Each completed request gets archived with implementation notes and a git commit. A built-in review runs after each item.

Other trigger words: `go`, `start`, `begin`, `process`, `execute`, `build`, `continue`, `resume`.

See the [Work Guide](docs/work-guide.md) for the full pipeline, triage routes, and clarify mode.

### 3. Run the full pipeline

One command, full cycle: investigate → capture → verify → run → review. Stateful and resumable — if the session ends mid-pipeline (context limit, crash, closed terminal), re-invoking picks up from the last step.

```
do work pipeline add dark mode to settings    # initialize with a request
do work pipeline                              # resume the active pipeline
do work pipeline status                       # progress without advancing
do work pipeline abandon                      # deactivate without completing
```

Pipeline state lives at `do-work/pipeline.json`. Each step dispatches to an existing action — the pipeline never re-implements logic.

### 4. Verify captured requests

Quality-check your captured REQs against the original input before building. Catches missed requirements, lost UX details, or intent drift.

```
do work verify requests
do work verify UR-003
do work check REQ-018
```

See the [Verify Requests Guide](docs/verify-requests-guide.md) for scoring dimensions and gap severity.

### 5. Review completed work

Post-build review: requirements check, code review, acceptance testing, and suggested testing. Also runs automatically after each work loop item.

```
do work review work
do work review REQ-005
do work review UR-003
```

See the [Review Work Guide](docs/review-work-guide.md) for the three review phases and scoring.

### 6. Answer pending questions

During the build phase, the assistant makes best-judgment calls on ambiguities instead of blocking. After work completes, review those decisions as a batch — confirm, override, or skip.

```
do work clarify
do work questions
```

### 7. Code review (standalone)

Review the actual codebase independent of the task queue. Scoped by prime files (architectural reference docs), directories, or both.

```
do work code-review                        # interactive — lists prime files, asks
do work code-review prime-auth             # everything prime-auth.md touches
do work code-review src/api/               # all source files in a directory
do work code-review prime-auth src/utils/  # combined scope
do work audit codebase
```

See the [Code Review Guide](docs/code-review-guide.md) for scoping modes, review dimensions, and health ratings.

### 8. UI review (read-only)

Validate UI quality against design best practices — structure, aesthetics, accessibility, UX copy, interaction patterns. Does not modify code.

```
do work ui-review                          # interactive — lists UI files, asks
do work ui-review src/components/          # validate a directory
do work ui-review prime-dashboard          # validate everything a prime file touches
do work design review
```

See the [UI Review Guide](docs/ui-review-guide.md) for review dimensions and severity levels.

### 9. Scan for quick wins

Find obvious improvements in a directory — dead code, duplication, complexity, missing tests.

```
do work quick-wins
do work quick-wins src/
do work scan src/api/
```

See the [Quick Wins Guide](docs/quick-wins-guide.md) for what it looks for and ranking criteria.

### 10. Generate ideas

Brainstorm what to build, improve, or explore next — grounded in codebase analysis and project history, not generic advice. Every idea points at something concrete and comes with a size estimate.

```
do work scan-ideas                # open exploration of the whole project
do work scan-ideas performance    # focused brainstorm on a theme
do work scan-ideas src/api/       # ideas scoped to a directory
do work ideas
do work brainstorm
```

Output is a ranked list of ideas, sized S/M/L, ready to paste into `do work capture request:`.

### 11. Explore a concept in depth

Multi-round structured exploration of a concept through specialized subagents — a Free Thinker (divergent generation), a Grounder (convergent evaluation), and a Writer (neutral synthesis). Produces idea briefs and a consolidated vision document.

```
do work deep-explore                       # ask what to explore
do work deep-explore "streaming rendering" # seed concept
do work deep-explore src/renderer/         # explore what a directory suggests
do work deep-explore continue              # resume an in-progress session
```

Use when you have a seed idea and want to develop it, not when you just want a list of tasks (use scan-ideas for that).

### 12. Manage prime files

Create and audit prime files — AI context documents that help an AI coder navigate a utility in minimum tokens. Prime files index entry points, traps, and exclusions so the AI doesn't waste tool calls rediscovering architecture.

```
do work prime create src/auth/    # interactive Q&A to generate a prime file
do work prime audit               # read-only health check of all prime files
do work prime                     # show prime sub-command help
```

See the [Prime Guide](docs/prime-guide.md) for the create workflow, audit checks, and prime file format.

### 13. Present work to clients

Generate client-facing deliverables from completed work: briefs, architecture diagrams, value propositions, Remotion videos, interactive HTML explainers.

```
do work present work
do work present UR-003
do work present all         # portfolio summary of all completed work
do work showcase
```

Artifacts are saved to `do-work/deliverables/`.

See the [Present Work Guide](docs/present-work-guide.md) for detail mode, portfolio mode, and artifact types.

### 14. Commit changes

Analyze uncommitted files, group them by REQ, and create atomic git commits with traceability.

```
do work commit
do work save work
```

See the [Commit Guide](docs/commit-guide.md) for grouping logic and commit message format.

### 15. Inspect changes

Read-only examination of uncommitted changes — explains what changed, why, traces to REQs, and assesses commit readiness.

```
do work inspect             # all uncommitted changes
do work inspect REQ-005     # changes for a specific REQ
do work inspect UR-003      # changes for all REQs under a UR
```

See the [Inspect Guide](docs/inspect-guide.md) for readiness signals and verdict definitions.

### 16. Cleanup the archive

Consolidate completed work — close finished UR folders, move loose REQs into their URs, organize legacy files.

```
do work cleanup
do work tidy
```

Also runs automatically at the end of every work loop.

See the [Cleanup Guide](docs/cleanup-guide.md) for the four consolidation passes.

### 17. Diagnostics

Pipeline health check — detects stuck work, hollow completions, orphaned URs, scope contamination. Read-only.

```
do work forensics
do work diagnose
do work health check
```

See the [Forensics Guide](docs/forensics-guide.md) for the full list of checks and severity levels.

### 18. Build a knowledge base

Build and maintain an LLM-friendly Markdown wiki from raw sources (PDFs, articles, notes). Aliases: `bkb`, `kb`, `build knowledge base`, `knowledge base`.

```
do work bkb init              # initialize a new knowledge base at ./kb
do work bkb init ~/research   # initialize at a custom path
do work bkb triage            # sort inbox items into capture directories
do work bkb ingest            # compile sources into wiki pages
do work bkb query [question]  # search the wiki and synthesize an answer
do work bkb lint              # quick health check
do work bkb lint full         # full structural check
do work bkb resolve           # walk through open contradictions
do work bkb close             # finalize daily log, refresh overview
do work bkb rollup            # monthly summary
do work bkb status            # show KB stats and pending items
```

For the full folder structure, file lifecycle, and wiki page format, see the [BKB Guide](docs/bkb-guide.md).

### 19. Install companion skills

```
do work install-ui-design   # Anthropic's frontend-design skill for production-grade UI
do work install-bowser      # Playwright CLI + Bowser skill for browser automation
```

### 20. Version and history

```
do work version             # current version + last 5 releases
do work update              # check for upstream updates
do work recap               # last 5 completed user requests
```

See the [Version Guide](docs/version-guide.md) for update behavior and recap format.

### 21. Learn the skill

Interactive tutorials for users new to do-work. Four modes:

```
do work tutorial                # ask which mode
do work tutorial quick-start    # hands-on walkthrough
do work tutorial concepts       # how the pieces fit together
do work tutorial recipes        # common workflow patterns
do work tutorial tour           # guided tour of the whole system
```

### Help

Run `do work help` at any point to get a refresher on all available commands.

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

Because capture preserves what you asked for, and processing tracks how it was built — and neither interferes with the other. When you say `do work capture request: ...`, your exact words are saved in a UR folder as the permanent source of truth. Nothing is paraphrased, nothing is lost. When `do work run` picks up that request later, the REQ file tracks every decision: what was planned, what was explored, what was built, what was reviewed. You end up with a clear trail from intent to implementation — what the user wanted, what the builder decided, and why. Without this separation, Claude tends to hear your request, immediately start coding, and leave no trace of what was asked or how decisions were made. The two-phase split means capture is fast and cheap (dump ideas anytime), processing is thorough and auditable (every change is traceable back to a request).

### Why not just let Claude decide what to do?

Claude *does* decide — the skill just raises the floor. Without structure, Claude picks different steps every time. Sometimes it plans, sometimes it dives in. Sometimes it reviews its work, sometimes it ships the first thing that compiles. The skill encodes lessons learned from many sessions into a repeatable baseline: plan before building, review after building, don't lose the original input along the way. Claude still makes all the real decisions within each step. The skill makes sure those decisions happen.

### How is this different from a hardcoded CI pipeline?

It's not a fixed sequence. The triage system (simple/medium/complex) means Claude chooses how much planning each request needs. Simple config changes skip straight to implementation. Complex features get explore, plan, then build. The skill is more like a senior dev's checklist — you still use judgment, but you don't skip steps because you felt confident.

### Do I need Claude Code specifically?

No. The skill works with any agentic coding tool that can read/write files and run shell commands. It was written for Claude Code but the action files are standalone prompts — paste them into any chat interface and they work. Subagent support is a nice-to-have, not a requirement.

### What if I only have one or two tasks?

The queue still helps. Even a single request benefits from the triage → build → review → commit pipeline. The overhead is near zero — capture is instant, and `do work run` processes whatever is there.

### Can I use this with an existing project?

Yes. Install it in your project root. The skill only creates files inside `do-work/` — it doesn't touch your source code structure. Your codebase is read during the build phase, but all skill state (requests, archives, deliverables) lives in `do-work/`.

### What happens if something goes wrong during processing?

The work loop processes one request at a time. If a request fails, it's marked as failed with notes on what went wrong, and the loop moves to the next one. Nothing is lost — you can fix the issue and re-queue. Run `do work forensics` to diagnose stuck or failed work.

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

## Designed for agentic coding tools

This skill assumes your tool supports:
- File editing and shell access
- Optional subagent or multi-agent workflows (Plan, Explore, Build)
- Git integration for per-request commits (optional)

Originally written for Claude Code. Works with other tools that can read/write files and run shell commands. If your tool does not support subagents, run Plan, Explore, and Implementation phases sequentially in the same session.

## License

MIT
