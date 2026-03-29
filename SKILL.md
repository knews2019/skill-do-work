---
name: do-work
description: Task queue - add requests or process pending work
argument-hint: (describe a task) | run | verify requests | review work | code-review | ui-review | present work | clarify | cleanup | quick-wins | install-ui-design | install-bowser | version | recap
upstream: https://raw.githubusercontent.com/knews2019/skill-do-work/main/SKILL.md
---

# Do-Work Skill

A unified entry point for task capture and processing.

**Actions:**

- **capture requests**: Capture new tasks/requests → creates UR folder (verbatim input) + REQ files (queue items), always paired
- **verify requests**: Evaluate captured REQs against original input → quality check
- **work**: Process pending requests → executes the queue
- **clarify questions**: Batch-review Open Questions from completed work → user answers, confirms, or skips
- **review work**: Post-work review → requirements check, code review, acceptance testing, and testing suggestions
- **present work**: Client-facing deliverables → briefs, architecture diagrams, value propositions, Remotion videos
- **cleanup**: Consolidate archive → moves loose REQs into UR folders, closes completed URs
- **code-review**: Standalone codebase review scoped by prime files and/or directories → consistency, patterns, security, architecture
- **quick-wins**: Scan a target directory for obvious refactoring opportunities and low-hanging tests to add
- **ui-review**: Validate UI quality against design best practices — read-only audit with structured findings report
- **install-ui-design**: Install the `frontend-design` Claude skill for production-grade UI design capabilities
- **install-bowser**: Install Playwright CLI + Bowser skill for browser automation, screenshots, and visual UI verification
- **forensics**: Pipeline diagnostics → detects stuck work, hollow completions, orphaned URs, scope contamination (read-only)
- **commit**: Commit uncommitted files → analyzes, groups atomically, traces to REQs
- **version**: Show current version, last 5 releases, or check for upstream updates
- **recap**: Summary of last 5 completed user requests with their REQs

> **Core concept:** The capture requests action always produces both a UR folder (preserving the original input) and REQ files (the queue items). Each REQ links back to its UR via `user_request` frontmatter. This pairing is mandatory for all requests — simple or complex.

> **Capture ≠ Execute.** The capture requests action captures requests. The work action executes them. These are strictly separate operations. After capture finishes writing files and reporting back, **STOP**. Do not start processing the queue, do not begin implementation, do not "helpfully" transition into the work action. The user decides when to execute — always. The only exception is if the user explicitly says something like "add this and then run it" or "capture this and start working" in the same invocation.

> **Human time has two optimal windows.** The system is designed to maximize the value of human attention:
>
> 1. **Capture phase** (capture requests action) — The user is present, actively thinking about the request. This is the best time for back-and-forth: clarifying ambiguities, resolving contradictions, making scope decisions. Use `AskUserQuestion` with concrete choices here. Every question must present options — never open-ended "what do you mean?" prompts.
>
> 2. **Batch question review** (clarify questions action) — After the build phase completes everything it can without feedback, any remaining `pending-answers` REQs are surfaced as a batch. The user reviews all builder-decided questions together, confirms or adjusts, and resolved REQs re-enter the queue.
>
> Between these windows, the build phase runs autonomously. Builders never block on Open Questions — they mark them `- [~]` with best-judgment reasoning and create `pending-answers` follow-ups when they return via `do work clarify`.

## Routing Decision

### Step 1: Parse the Input

Examine what follows "do work":


Check these patterns **in order** — first match wins:

| Priority | Pattern                  | Example                                                                                                                            | Route                         |
| -------- | ------------------------ | ---------------------------------------------------------------------------------------------------------------------------------- | ----------------------------- |
| 1        | Empty or bare invocation | `do work`                                                                                                                          | → help menu                   |
| 2        | Version exact phrases    | `do work check for updates`, `do work check for update`                                                                            | → version                     |
| 3        | Action verbs only        | `do work run`, `do work go`, `do work start`                                                                                       | → work                        |
| 4        | Verify keywords          | `do work verify`, `do work verify requests`, `do work check REQ-018`, `do work evaluate`                                           | → verify requests              |
| 5        | Clarify keywords         | `do work clarify`, `do work questions`, `do work pending`                                                                          | → clarify questions            |
| 6        | Code-review keywords     | `do work code-review`, `do work code-review prime-auth`, `do work code review src/`, `do work audit codebase`, `do work review codebase`, `do work codebase review` | → code-review                  |
| 7        | UI-review keywords       | `do work ui-review`, `do work ui-review src/`, `do work review ui`, `do work design review`, `do work validate ui`                 | → ui-review                    |
| 8        | Review keywords          | `do work review`, `do work review work`, `do work review code`, `do work code review` (no scope), `do work audit code`             | → review work                  |
| 9        | Present keywords         | `do work present`, `do work present work`, `do work showcase`, `do work deliver`                                                   | → present work                 |
| 10       | Cleanup keywords         | `do work cleanup`, `do work tidy`, `do work consolidate`                                                                           | → cleanup                     |
| 11       | Commit keywords          | `do work commit`, `do work commit changes`, `do work save work`                                                                    | → commit                      |
| 12       | Version keywords         | `do work version`, `do work update`, `do work what's new`, `do work release notes`, `do work what's changed`, `do work updates`, `do work history` | → version                     |
| 13       | Recap keywords           | `do work recap`                                                                                                                    | → version                     |
| 14       | Forensics keywords       | `do work forensics`, `do work diagnose`, `do work health check`, `do work health`                                                  | → forensics                   |
| 15       | Quick-wins keywords      | `do work quick-wins`, `do work quick wins`, `do work low-hanging`                                                                  | → quick-wins                  |
| 16       | Install keywords         | `do work install-ui-design`, `do work install ui design`, `do work install-bowser`, `do work install bowser`, `do work install playwright` | → install-ui-design / install-bowser |
| 17       | Descriptive content      | `do work add dark mode`, `do work [meeting notes]`, `do work capture request [the request]`                                        | → capture requests              |


### Step 2: Preserve Payload

**Critical rule**: Never lose the user's content.

**Single-word rule**: A single word is either a known keyword or ambiguous — it is never "descriptive content."

- **Matches a keyword** in the routing table (e.g., "version", "verify", "cleanup") → route to that action directly.
- **Doesn't match any keyword** (e.g., "refactor", "optimize") → ambiguous. Ask: "Do you want to add '`{word}`' as a new request, or did you mean something else?"

Only route to **capture requests** when the input is clearly descriptive — multiple words, a sentence, a feature request, etc.

If routing is genuinely unclear AND multi-word content was provided:

- Default to **capture requests** (adding a task)
- Hold onto $ARGUMENTS
- If truly ambiguous, ask: "Add this as a request, or start the work loop?"
- User replies with just "add" or "work" → proceed with original content

### Action Verbs (→ Work)

These signal "process the queue":
run, go, start, begin, work, process, execute, build, continue, resume

### Clarify Verbs (→ Clarify Questions)

These signal "review pending questions":
clarify, answers, questions, pending, pending answers, blocked, what's blocked, what needs answers

Note: This routes to the work action with `mode: clarify` — see work.md "Clarify Questions" section.

### Verify Verbs (→ Verify Requests)

These signal "check request quality":
verify, verify requests, check, evaluate, review requests, review reqs, audit

Note: "check" routes to verify requests ONLY when used alone or with a target (e.g., "do work check UR-003"). "check for updates" is intercepted at priority 2 and routes to version — it never reaches verify. When followed by other descriptive content it routes to capture requests (e.g., "do work check if the button works" → capture requests).

Note: "audit" alone routes to verify requests. "audit code" and "audit implementation" route to review work (see Review Verbs below). "audit codebase" routes to code-review (see Code-Review Verbs below).

### Code-Review Verbs (→ Code Review)

These signal "standalone codebase review":
code-review, code review [scope], review codebase, audit codebase, review codebase [scope], codebase review

Note: "code-review" (hyphenated) always routes to code-review (priority 6). "code review" followed by a **prime file reference or directory path** routes to code-review (priority 6). "codebase review" always routes to code-review (priority 6). Plain "code review" (no scope) falls through to **review work** (priority 8) for backwards compatibility. "audit codebase" and "review codebase" always route to code-review. The key distinction: review work reviews completed REQ/UR work items; code-review reviews the actual source code independent of the queue.

Scope arguments are passed through as `$ARGUMENTS`:
- Prime file references: `prime-auth`, `prime-auth.md`, `src/prime-auth.md`
- Directory paths: `src/`, `src/api/ src/utils/`
- Combined: `prime-auth src/utils/`
- No scope: interactive — lists available prime files and asks

### UI-Review Verbs (→ UI Review)

These signal "validate UI quality (read-only)":
ui-review, review ui, design review, validate ui, ui audit, design audit

Note: "ui-review" (hyphenated) always routes to ui-review. "review ui" and "design review" route to ui-review. "validate ui" routes to ui-review. Do NOT use "check ui" — "check" is consumed by verify-requests at priority 4 before reaching this rule. The key distinction from code-review: ui-review evaluates visual design, UX, accessibility, and component consistency against design best practices. code-review evaluates code patterns, architecture, and security.

Scope arguments are passed through as `$ARGUMENTS`:
- File paths: `src/components/Header.tsx`
- Directory paths: `src/pages/`
- Prime file references: `prime-dashboard`
- Combined: `prime-auth src/components/`
- No scope: interactive — lists UI-relevant files and asks

### Review Verbs (→ Review Work)

These signal "review the completed work":
review, review work, review code, code review, audit code, audit implementation, review REQ-NNN, review UR-NNN

Note: "review requests" and "review reqs" route to **verify requests** (priority 4), not review work. "review" alone or followed by a target/code-related word routes to review work (priority 8). The review work action also runs automatically as part of the work pipeline — see `work.md` Step 7.

### Present Verbs (→ Present Work)

These signal "generate client-facing deliverables":
present, present work, showcase, deliver, pitch, client brief

Note: `do work present` (no target) presents the most recent completed UR. `do work present all` or `do work present portfolio` enters portfolio mode. `do work present UR-003` or `do work present REQ-005` targets specific work.

### Cleanup Verbs (→ Cleanup)

These signal "consolidate the archive":
cleanup, clean up, tidy, consolidate, organize archive, fix archive

### Commit Verbs (→ Commit)

These signal "commit uncommitted files atomically":
commit, commit changes, commit files, save changes, save work

### Recap Verbs (→ Version)

These signal "show recent work summary":
recap

### Version / Release Info Verbs (→ Version)

These signal "show version or release info":
version, update, check for updates, what's new, release notes, what's changed, updates, history

Note: "updates" (plural) and "what's new" show version + last 5 releases. "update" (singular) triggers the update check flow. Both are handled by the version action.

### Quick-Wins Verbs (→ Quick-Wins)

These signal "scan for improvement opportunities":
quick-wins, quick wins, low-hanging, low hanging fruit, scan, opportunities, what can we improve

Note: "scan", "opportunities", and "what can we improve" route to quick-wins ONLY when used alone or with a directory path (e.g., "do work scan", "do work scan src/"). When followed by descriptive content they route to capture requests (e.g., "do work scan the checkout logs for 500s" → capture requests).

### Install Verbs (→ Install UI Design / Install Bowser)

Two install actions exist. Route based on the keyword after "install":

**UI Design** — these signal "install the frontend-design skill":
install-ui-design, install ui design, install ui, install frontend-design, setup ui design, setup design skill

**Bowser** — these signal "install Playwright CLI + Bowser skill":
install-bowser, install bowser, install playwright, install playwright-cli, setup bowser, setup playwright

### Content Signals (→ Capture Requests)

These signal "add a new task":

- Descriptive text beyond a single verb
- Feature requests, bug reports, ideas
- Screenshots or context
- "add", "create", "I need", "we should"

## Examples

### Help Menu (bare invocation)

When invoked with no arguments (`do work`), show a help menu with available actions and example prompts:

```
do-work — task queue for agentic coding tools

  Capture requests:
    do work add dark mode to settings
    do work the search is slow and the header is misaligned
    do work [paste meeting notes, specs, or a screenshot]

  Process the queue:
    do work run

  Verify & review:
    do work verify requests     Check capture quality against original input
    do work review work         Review completed work (requirements + code + acceptance)

  Code review (standalone):
    do work code-review                   Review codebase (interactive scope selection)
    do work code-review prime-auth        Review everything prime-auth.md touches
    do work code-review src/api/          Review a directory
    do work code-review prime-auth src/   Review prime file scope + directory combined

  Present to client:
    do work present work        Generate client brief, architecture, video, and interactive HTML explainer
    do work present all         Portfolio summary of all completed work

  Scan for improvements:
    do work quick-wins          Scan cwd for refactoring and test opportunities
    do work quick-wins src/     Scan a specific directory

  UI review (read-only):
    do work ui-review                     Validate UI quality (interactive scope selection)
    do work ui-review src/components/     Validate a directory
    do work ui-review prime-dashboard     Validate everything a prime file touches

  Setup:
    do work install-ui-design   Install the frontend-design skill for production-grade UI
    do work install-bowser      Install Playwright CLI + Bowser skill for browser automation

  Diagnostics:
    do work forensics           Pipeline diagnostics — stuck work, hollow completions, orphaned URs

  Other actions:
    do work clarify             Answer pending questions from completed work
    do work cleanup             Consolidate the archive
    do work commit              Analyze and commit uncommitted files atomically
    do work version             Check version + last 5 skill releases
    do work update              Check for upstream updates
    do work recap               Last 5 completed URs with their REQs
```

Do not ask "Start the work loop?" — just print the help menu and wait.

### Routes to Work

- `do work run` → Starts work action immediately
- `do work go` → Starts work action immediately

### Routes to Clarify Questions

- `do work clarify` → Presents all pending-answers REQs for batch review
- `do work questions` → Same as clarify
- `do work answers` → Same as clarify
- `do work pending` → Same as clarify
- `do work what's blocked` → Same as clarify

### Routes to Verify Requests

- `do work verify requests` → Evaluates most recent UR's REQs
- `do work verify` → Evaluates most recent UR's REQs
- `do work verify UR-003` → Evaluates specific UR
- `do work check REQ-018` → Evaluates the UR that REQ-018 belongs to
- `do work evaluate` → Evaluates most recent UR's REQs
- `do work review requests` → Evaluates most recent UR's REQs

### Routes to Code Review

- `do work code-review` → Interactive scope selection (lists available prime files)
- `do work code-review prime-auth` → Reviews all files referenced by prime-auth.md
- `do work code-review prime-auth.md` → Same (explicit extension)
- `do work code-review src/prime-auth.md` → Same (explicit path)
- `do work code-review prime-auth prime-checkout` → Reviews union of both prime file scopes
- `do work code-review src/` → Reviews all source files in src/
- `do work code-review src/api/ src/utils/` → Reviews multiple directories
- `do work code-review prime-auth src/utils/` → Combined: prime file scope + directory
- `do work audit codebase` → Same as code-review (no scope → interactive)
- `do work review codebase` → Same as code-review
- `do work review codebase src/` → Reviews src/ directory
- `do work codebase review` → Same as code-review

### Routes to UI Review

- `do work ui-review` → Interactive scope selection (lists UI-relevant files)
- `do work ui-review src/components/` → Validates all UI files in directory
- `do work ui-review prime-dashboard` → Validates all files referenced by prime-dashboard.md
- `do work ui-review prime-auth src/components/` → Combined: prime file scope + directory
- `do work review ui` → Same as ui-review (no scope → interactive)
- `do work design review` → Same as ui-review
- `do work validate ui` → Same as ui-review
- `do work design review src/pages/` → Validates specific directory
- `do work ui audit` → Same as ui-review
- `do work design audit` → Same as ui-review

### Routes to Review Work

- `do work review work` → Reviews the most recently completed REQ
- `do work review` → Reviews the most recently completed REQ
- `do work review REQ-005` → Reviews a specific completed REQ
- `do work review UR-003` → Reviews all completed REQs under that UR
- `do work code review` → Reviews the most recently completed REQ
- `do work review code` → Reviews the most recently completed REQ

### Routes to Present Work

- `do work present work` → Generates deliverables for most recently completed UR
- `do work present` → Same as present work
- `do work present UR-003` → Generates deliverables for specific UR
- `do work present REQ-005` → Generates deliverables for specific REQ
- `do work present all` → Portfolio summary of all completed work
- `do work present portfolio` → Same as present all
- `do work showcase` → Same as present work

### Routes to Cleanup

- `do work cleanup` → Consolidates archive, closes completed URs
- `do work tidy` → Same as cleanup
- `do work consolidate` → Same as cleanup

### Routes to Commit

- `do work commit` → Analyzes and commits all uncommitted files atomically
- `do work commit changes` → Same as commit
- `do work save work` → Same as commit

### Routes to Version

- `do work version` → Reports version + last 5 skill releases
- `do work update` → Checks for upstream updates
- `do work check for updates` → Same as update
- `do work what's new` → Same as version (shows releases)
- `do work release notes` → Same as version
- `do work updates` → Same as version
- `do work history` → Same as version

### Routes to Recap (via Version)

- `do work recap` → Last 5 completed URs with their REQs

### Routes to Forensics

- `do work forensics` → Full pipeline diagnostics (read-only)
- `do work diagnose` → Same as forensics
- `do work health check` → Same as forensics
- `do work health` → Same as forensics

### Routes to Quick-Wins

- `do work quick-wins` → Scans current working directory
- `do work quick wins` → Same
- `do work quick-wins src/` → Scans specific directory
- `do work low-hanging` → Same
- `do work scan` → Scans current working directory
- `do work scan src/` → Scans specific directory
- `do work scan the checkout logs for 500s` → Routes to capture requests (descriptive content)
- `do work opportunities` → Scans current working directory

### Routes to Install UI Design

- `do work install-ui-design` → Installs the frontend-design skill
- `do work install ui design` → Same
- `do work install ui` → Same
- `do work setup design skill` → Same

### Routes to Install Bowser

- `do work install-bowser` → Installs Playwright CLI + Bowser skill
- `do work install bowser` → Same
- `do work install playwright` → Same
- `do work setup bowser` → Same
- `do work setup playwright` → Same

### Routes to Capture Requests

- `do work add dark mode` → Creates REQ file + UR folder
- `do work the button is broken` → Creates REQ file + UR folder
- `do work [400 words]` → Creates REQ files + UR folder with full verbatim input

## Payload Preservation Rules

When clarification is needed but content was provided:

1. **Do not lose $ARGUMENTS** - keep the full payload in context
2. **Ask a simple question**: "Add this as a request, or start the work loop?"
3. **Accept minimal replies**: User says just "add" or "work"
4. **Proceed with original content**: Apply the chosen action to the stored arguments
5. **Never ask the user to re-paste content**

This enables a two-phase commit pattern:

1. Capture intent payload
2. Confirm action

## Action Dispatch

Each action has an action file with full instructions. How you execute it depends on your environment's capabilities.

| Action             | Action file                     | Context to pass                |
|--------------------|---------------------------------|--------------------------------|
| capture requests   | `./actions/capture.md`          | Full user input text           |
| work               | `./actions/work.md`             | (none needed)                  |
| clarify questions  | `./actions/work.md`             | `mode: clarify`                |
| verify requests    | `./actions/verify-requests.md`  | Target UR/REQ or "most recent" |
| review work        | `./actions/review-work.md`      | Target REQ/UR or "most recent" |
| present work       | `./actions/present-work.md`     | Target REQ/UR, "most recent", or "all" |
| cleanup            | `./actions/cleanup.md`          | (none needed)                  |
| commit             | `./actions/commit.md`           | (none needed)                  |
| code-review        | `./actions/code-review.md`      | Prime file refs and/or directory paths |
| ui-review          | `./actions/ui-review.md`        | File/directory paths and/or prime file refs |
| quick-wins         | `./actions/quick-wins.md`       | Target directory               |
| install-ui-design  | `./actions/install-ui-design.md`| (none needed)                  |
| install-bowser     | `./actions/install-bowser.md`   | (none needed)                  |
| forensics          | `./actions/forensics.md`        | (none needed)                  |
| version            | `./actions/version.md`          | `$ARGUMENTS`                   |
| recap              | `./actions/version.md`          | `mode: recap`                  |

### If subagents are available

Dispatch each action to a subagent. The subagent reads the action file and executes it — the main thread only sees the routing decision and the returned summary.

- **`work` and `cleanup`**: Run in the background if your environment supports it. Print a status line (e.g., "Work queue processing in background...") and return control to the user immediately.
- **`capture requests`, `clarify questions`, `verify requests`, `review work`, `code-review`, `ui-review`, `present work`, `quick-wins`, `forensics`, `version`, `recap`**: Run in the foreground (blocking). These need user interaction or produce small immediate output.
- **Screenshots (`capture requests` only):** Subagents can't see images from the main conversation. Before dispatching, save screenshots to `do-work/user-requests/.pending-assets/screenshot-{n}.png`, write a text description of each, and include the paths + descriptions in the subagent prompt.

### If subagents are not available

Read the action file directly and follow its instructions in the current session. The action files are designed to work as standalone prompts — no subagent infrastructure required.

### On failure

Report the error to the user. Do not retry automatically.

## Suggest Next Steps

After every action completes, suggest the next logical prompts the user might want to run. Use fully qualified action names so the user can copy-paste directly.

**After capture requests:**
```
Next steps:
  do work verify requests     Check capture quality before building
  do work run                 Start processing the queue
```

**After work (queue processing):**
```
Next steps:
  do work review work         Review the completed work
  do work present work        Generate client-facing deliverables
  do work clarify             Answer any pending questions
```

**After verify requests:**
```
Next steps:
  do work run                 Start processing the queue
  do work [describe changes]  Capture additional requests
```

**After review work:**
```
Next steps:
  do work present work        Generate client-facing deliverables
  do work ui-review [scope]   Validate UI quality (if domain: ui-design)
  do work run                 Process follow-up REQs (if any were created)
```

**After code-review:**
```
Next steps:
  do work run                   Process follow-up REQs (if any were created)
  do work ui-review [scope]     Validate UI quality for the same scope
  do work quick-wins [dir]      Scan for additional improvements
```

**After ui-review:**
```
Next steps:
  do work [describe fix]        Capture findings as requests
  do work run                   Process follow-up REQs (if any were created)
  do work install-bowser        Install Playwright CLI + Bowser skill for visual verification (if not installed)
```

**After present work:**
```
Next steps:
  do work present all         Generate portfolio summary (if multiple URs completed)
  do work [describe changes]  Capture new requests
```

**After forensics:**
```
Next steps:
  do work cleanup               Fix orphaned URs and misplaced files
  do work run                   Process stuck or pending REQs
  do work [describe fix]        Capture a specific finding as a request
```

**After quick-wins:**
```
Next steps:
  do work [describe fix]        Capture a finding as a request
  do work run                   Process the queue
```

**After commit:**
```
Next steps:
  do work review work         Review the committed changes
  do work [describe changes]  Capture new requests
```

**After clarify questions:**
```
Next steps:
  do work run                 Process answered questions
  do work clarify             Continue answering (if skipped any)
```

**After version / recap:**
```
Next steps:
  do work run                 Start processing the queue
  do work [describe changes]  Capture new requests
```

**Rules:**
- Only suggest prompts that provide value given the current state (e.g., don't suggest `do work run` if the queue is empty)
- Use the full action name (`verify requests`, not just `verify`; `review work`, not just `review`)
- Keep it to 2-3 suggestions max — don't overwhelm
- Format as a simple list the user can scan and copy

