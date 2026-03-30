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

### 13. Install companion skills

```
do work install-ui-design   # Anthropic's frontend-design skill for production-grade UI
do work install-bowser      # Playwright CLI + Bowser skill for browser automation
```

### 14. Version and history

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

## Designed for agentic coding tools

This skill assumes your tool supports:
- File editing and shell access
- Optional subagent or multi-agent workflows (Plan, Explore, Build)
- Git integration for per-request commits (optional)

Originally written for Claude Code. Works with other tools that can read/write files and run shell commands. If your tool does not support subagents, run Plan, Explore, and Implementation phases sequentially in the same session.

## License

MIT
