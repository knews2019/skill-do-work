# do-work

A task queue skill for agentic coding tools. Capture requests fast, process them later.

## Installation

```bash
# Run from the directory where you want the skill installed
curl -sL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz | tar xz --strip-components=1 --exclude='_dev'
```

**Updating:** Re-run the same command to update. Note that tar extraction overwrites but does not delete files removed upstream — stale files from older versions may linger (generally harmless). For a fully clean update, delete only the known skill paths (`actions/`, `agent-rules/`, `SKILL.md`, `CHANGELOG.md`, `README.md`) before re-extracting — never delete `do-work/` or other project files.

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

### 3. Verify captured requests

Quality-check your captured REQs against the original input before building. Catches missed requirements, lost UX details, or intent drift.

```
do work verify requests
do work verify UR-003
do work check REQ-018
```

### 4. Review completed work

Post-build review: requirements check, code review, acceptance testing, and suggested testing. Also runs automatically after each work loop item.

```
do work review work
do work review REQ-005
do work review UR-003
```

### 5. Answer pending questions

During the build phase, the assistant makes best-judgment calls on ambiguities instead of blocking. After work completes, review those decisions as a batch — confirm, override, or skip.

```
do work clarify
do work questions
```

### 6. Code review (standalone)

Review the actual codebase independent of the task queue. Scoped by prime files (architectural reference docs), directories, or both.

```
do work code-review                        # interactive — lists prime files, asks
do work code-review prime-auth             # everything prime-auth.md touches
do work code-review src/api/               # all source files in a directory
do work code-review prime-auth src/utils/  # combined scope
do work audit codebase
```

### 7. UI review (read-only)

Validate UI quality against design best practices — structure, aesthetics, accessibility, UX copy, interaction patterns. Does not modify code.

```
do work ui-review                          # interactive — lists UI files, asks
do work ui-review src/components/          # validate a directory
do work ui-review prime-dashboard          # validate everything a prime file touches
do work design review
```

### 8. Scan for quick wins

Find obvious improvements in a directory — dead code, duplication, complexity, missing tests.

```
do work quick-wins
do work quick-wins src/
do work scan src/api/
```

### 9. Present work to clients

Generate client-facing deliverables from completed work: briefs, architecture diagrams, value propositions, Remotion videos, interactive HTML explainers.

```
do work present work
do work present UR-003
do work present all         # portfolio summary of all completed work
do work showcase
```

Artifacts are saved to `do-work/deliverables/`.

### 10. Commit changes

Analyze uncommitted files, group them by REQ, and create atomic git commits with traceability.

```
do work commit
do work save work
```

### 11. Cleanup the archive

Consolidate completed work — close finished UR folders, move loose REQs into their URs, organize legacy files.

```
do work cleanup
do work tidy
```

Also runs automatically at the end of every work loop.

### 12. Diagnostics

Pipeline health check — detects stuck work, hollow completions, orphaned URs, scope contamination. Read-only.

```
do work forensics
do work diagnose
do work health check
```

### 13. Build a knowledge base

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

### 14. Install companion skills

```
do work install-ui-design   # Anthropic's frontend-design skill for production-grade UI
do work install-bowser      # Playwright CLI + Bowser skill for browser automation
```

### 15. Version and history

```
do work version             # current version + last 5 releases
do work update              # check for upstream updates
do work recap               # last 5 completed user requests
```

### Help

Run `do work help` at any point to get a refresher on all available commands.

## File structure

```
do-work/
├── REQ-018-pending.md       # Queue (pending requests)
├── REQ-019-pending.md
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

## Designed for agentic coding tools

This skill assumes your tool supports:
- File editing and shell access
- Optional subagent or multi-agent workflows (Plan, Explore, Build)
- Git integration for per-request commits (optional)

Originally written for Claude Code. Works with other tools that can read/write files and run shell commands. If your tool does not support subagents, run Plan, Explore, and Implementation phases sequentially in the same session.

## License

MIT
