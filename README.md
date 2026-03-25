# do-work

A task queue skill for agentic coding tools. Capture requests fast, process them later.

## Installation

```bash
# Run from the directory where you want the skill installed
curl -sL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz | tar xz --strip-components=1 --exclude='_dev'
```

**Updating:** Re-run the same command to update. Note that tar extraction overwrites but does not delete files removed upstream — stale files from older versions may linger (generally harmless). For a fully clean update, delete only the known skill paths (`actions/`, `agent-rules/`, `SKILL.md`, `CHANGELOG.md`, `README.md`) before re-extracting — never delete `do-work/` or other project files.

## Welcome to your new work loop

This skill gives you a two-phase workflow:

1. **Capture**: Throw ideas, bugs, and feature requests at your assistant as they come up. Each one becomes a structured request file in `do-work/`.

2. **Process**: When you're ready, tell the assistant to work. It picks up pending requests one by one, triages complexity, and builds until the queue is empty.

The idea: separate *thinking of things* from *doing things*. Capture is instant. Processing is thorough.

## Quick usage

**Add a request:**
```
do work add dark mode to the settings page
```
Creates `do-work/REQ-001-dark-mode.md`

**Add multiple at once:**
```
do work the search is slow, also add an export button, and fix the header alignment
```
Creates three separate request files.

**Process the queue:**
```
do work run
```
Starts the work loop. The assistant triages each request by complexity:
- **Simple** (config changes, small fixes) → straight to implementation
- **Medium** (clear goal, unknown location) → explore codebase first
- **Complex** (new features, architectural) → plan, explore, then build

Each completed request gets archived with its implementation notes and a git commit.

## How it works

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

Every capture invocation creates a User Request (UR) folder preserving the verbatim input. REQ files in the queue reference their UR. When all REQs from a UR are completed, the UR folder moves to archive as a self-contained unit.

Legacy REQs (created before the UR system) work the same as before — they archive directly without a UR folder.

## Designed for Agentic Coding Tools

This skill assumes your tool supports:
- File editing and shell access
- Optional subagent or multi-agent workflows (Plan, Explore, Build)
- Git integration for per-request commits (optional)

It was originally written for Claude Code and should work with other tools that provide similar capabilities. If your tool does not support subagents, run the Plan, Explore, and Implementation phases sequentially in the same session.

## Actions

### Capture Requests

Invoked when you provide descriptive content. Optimized for speed:
- Minimal questions — capture what was said, don't interrogate
- Handles simple one-liners and complex multi-feature specs
- Always creates a UR folder preserving the full verbatim input
- Checks for duplicates against existing requests

See [actions/capture.md](./actions/capture.md) for the full capture logic.

### Work (process)

Invoked when you say "run", "go", or "start". Runs the build loop:
- Triages each request to determine the right amount of planning
- Spawns specialized agents only when needed
- Archives completed work with implementation notes
- Creates atomic git commits per request

See [actions/work.md](./actions/work.md) for the full processing logic.

### Clarify Questions

Invoked when you say "clarify", "questions", or "answers". Batch-reviews Open Questions from completed work:
- Presents all `pending-answers` REQs and their unresolved questions
- User can answer, confirm the builder's choice, or skip
- Confirmed choices resolve without re-entering the work loop
- Answered questions flip the REQ to `pending` for the next work run

See [actions/work.md](./actions/work.md) "Clarify Questions" section for the full workflow.

### Verify Requests

Invoked when you say "verify requests", "verify", "check", "evaluate", or "review requests". Quality gate for captured requirements:
- Reads the original user input from the UR folder
- Compares against extracted REQ files for completeness
- Scores coverage, UX detail capture, intent signal preservation
- Optionally fixes identified gaps

See [actions/verify-requests.md](./actions/verify-requests.md) for the full evaluation logic.

### Review Work

Invoked when you say "review work", "review", "review code", or "code review". Also runs automatically after each work loop item completes. Comprehensive post-work review:
- **Requirements check** — walks through every requirement to confirm it was built
- **Code review** — evaluates code quality, scope discipline, and risk
- **Acceptance testing** — actually runs/tests the feature to verify it works
- **Suggested testing** — recommends additional checks the user should perform
- **Human UAT** — interactively collects your manual testing feedback and lessons learned (Standalone mode only)
- Creates follow-up REQs for Important findings

See [actions/review-work.md](./actions/review-work.md) for the full review logic.

### Code Review (standalone)

Invoked when you say "code-review", "audit codebase", or "review codebase". Standalone codebase review — not tied to the REQ/UR queue:
- **Scoped by prime files**: `do work code-review prime-auth` reviews everything that prime file touches
- **Scoped by directories**: `do work code-review src/api/` reviews all source files in that directory
- **Combined**: `do work code-review prime-auth src/utils/` reviews the union of both scopes
- **Interactive**: `do work code-review` (no scope) lists available prime files and asks
- Evaluates consistency, naming, error handling, architecture, security, and test coverage
- Optionally creates REQ files for Critical/Important findings

See [actions/code-review.md](./actions/code-review.md) for the full review logic.

### Present Work

Invoked when you say "present work", "present", "showcase", or "deliver". Generates client-facing deliverables from completed work:
- **Client Brief** — what was built, how it works (architecture + data flow), why it matters
- **Value Proposition** — business impact, revenue opportunities, competitive advantage
- **Remotion Video** — browser-previewable Remotion video project (React components, no mp4) when the feature is user-facing
- **Interactive Explainer** — a zero-dependency, clickable HTML file demonstrating the before/after state and data flow
- **Portfolio Summary** — `do work present all` aggregates all completed work into a cumulative overview

Artifacts are saved to `do-work/deliverables/` for reuse.

See [actions/present-work.md](./actions/present-work.md) for the full presentation logic.

### Cleanup (consolidate)

Invoked when you say "cleanup", "tidy", or "consolidate". Also runs automatically at the end of every work loop. Keeps the archive organized:
- Closes UR folders in `user-requests/` when all their REQs are complete
- Moves loose REQ files from `archive/` root into their UR folders
- Moves legacy REQs (no UR reference) into `archive/legacy/`
- Fixes misplaced folders (e.g., `archive/user-requests/UR-NNN` → `archive/UR-NNN`)

See [actions/cleanup.md](./actions/cleanup.md) for the full consolidation logic.

### Quick-Wins (scan)

Invoked when you say "quick-wins", "scan", or "low-hanging". Scans a target directory for obvious improvements:
- Refactoring opportunities (dead code, duplication, complexity)
- Low-hanging tests to add
- Optionally targets a specific directory (`do work quick-wins src/`)

See [actions/quick-wins.md](./actions/quick-wins.md) for the full scan logic.

### Commit

Invoked when you say "commit", "commit changes", or "save work". Atomic git commits traced to REQs:
- Analyzes uncommitted files and groups them by REQ
- Creates one commit per REQ with traceability in the commit message
- Validates Implementation Summary against staged files

See [actions/commit.md](./actions/commit.md) for the full commit logic.

### UI Review (validate)

Invoked when you say "ui-review", "review ui", "design review", or "validate ui". Read-only UI quality audit — does not modify code:
- **Scoped by files/directories**: `do work ui-review src/components/` validates all UI files in that directory
- **Scoped by prime files**: `do work ui-review prime-dashboard` validates everything that prime file touches
- **Interactive**: `do work ui-review` (no scope) lists UI-relevant files and asks
- Evaluates structure/IA, visual aesthetics, component consistency, UX copy, interaction/accessibility, and implementation patterns
- Produces a severity-rated findings report with file:line references and concrete fix suggestions
- Leverages both `rules-ui-design.md` and the `frontend-design` skill (if installed)
- Uses Playwright CLI or the Bowser skill when available for rendered-page validation
- Optionally captures findings as `domain: ui-design` REQs in the queue

See [actions/ui-review.md](./actions/ui-review.md) for the full validation logic.

### Install UI Design

Invoked when you say "install-ui-design" or "install ui design". Installs Anthropic's `frontend-design` skill for production-grade UI design capabilities:
- Checks if already installed
- Installs the skill into `.claude/skills/frontend-design/`
- Works alongside `domain: ui-design` rules for a complete design workflow

See [actions/install-ui-design.md](./actions/install-ui-design.md) for the full installation logic.

### Version / Recap

Invoked when you say "version", "update", or "what's new". Also supports "recap" for work history:
- **Version**: Shows current version and last 5 skill releases
- **Update**: Checks for upstream updates and applies them
- **Recap**: Summary of last 5 completed user requests

See [actions/version.md](./actions/version.md) for the full version logic.

## License

MIT
