---
name: do-work
description: Task queue - add requests or process pending work
argument-hint: "pipeline [request] | capture-request: (describe a task) | run | verify-requests | review-work | code-review | ui-review | validate-feedback [paste findings] | present-work | ai-report [target] | slop-check [target] | dream [path] | clarify | cleanup | commit | inspect | quick-wins | scan-ideas [focus] | deep-explore [concept] | prime [create|audit] | forensics | roadmap [scope] | note [text] | stray-check [path] | bkb [subcommand] | interview [template] | prompts [subcommand] | install [target] | version | recap | tutorial [mode] | help"
---

# Do-Work Skill

A unified entry point for task capture and processing.

**Actions:**

- **pipeline**: Full end-to-end orchestration → investigate, capture, verify, run, review, present in sequence with persistent state tracking
- **capture-requests**: Capture new tasks/requests → creates UR folder (verbatim input) + REQ files (queue items), always paired; for testable behavior, infer and confirm the RED case/GREEN proof target during capture. This is essential, high-value work.
- **verify-requests**: Evaluate captured REQs against original input → quality check
- **work**: Process pending requests → executes the queue
- **clarify questions**: Batch-review Open Questions from completed work → user answers, confirms, or skips
- **review-work**: Post-work review → requirements check, code review, acceptance testing, and testing suggestions
- **validate-feedback**: Triage external review feedback / audit findings → per item, verify against the real code + git history and recommend Already done / Accept / Push back / Discuss. Read-only; offers a capture handoff for accepted items
- **present-work**: Client-facing deliverables → briefs, architecture diagrams, value propositions, Remotion videos
- **ai-report**: Self-contained folder report of a completed UR/REQ (`index.html` + `screenshots/`, plus `generated/` for AI images) — live screenshots with SVG callouts, before/after toggles, optional AI-generated diagrams, Mermaid fallback. Output to `ai-reports/`
- **cleanup**: Consolidate archive → moves loose REQs into UR folders, closes completed URs
- **code-review**: Standalone codebase review scoped by prime files and/or directories → consistency, patterns, security, performance, architecture, and risk-driven test coverage
- **quick-wins**: Scan a target directory for obvious refactoring opportunities and low-hanging tests to add
- **scan-ideas**: Generate ideas for what to build, improve, or explore next — grounded in codebase analysis and project history
- **deep-explore**: Multi-round structured exploration of a concept — spawns divergent/convergent subagent dialogue, produces vision documents and idea briefs
- **ui-review**: Validate UI quality against design best practices — read-only audit with structured findings report
- **slop-check**: Validate a human-facing artifact (brief, report, summary) against the anti-slop principles before it ships — read-only by default, optional rewrite on confirmation
- **dream**: Manual four-phase consolidation of a plain-text memory directory — orient, lint, heal, prune + reindex. Destructive; explicit invocation only.
- **install**: Install companion skills/tooling into the current project. Targets: `ui-design` (Anthropic's `frontend-design` skill) and `bowser` (Playwright CLI + Bowser skill for browser automation, screenshots, and visual UI verification).
- **forensics**: Pipeline diagnostics → detects stuck work, hollow completions, orphaned URs, scope contamination, recurring corrections (read-only)
- **roadmap**: Queue survey → classifies pending REQs (ready / needs-clarification / blocked / stale), reports TDD posture, rolls up in-progress and recently-completed work (read-only)
- **note**: Append a lightweight, dated next-step hint to `do-work/notes.md` → surfaced at the top of `do-work roadmap`. Not a REQ; no capture, no schema, no implementation. User deletes lines by hand when resolved
- **stray-check**: Repo-wide orphan/junk scan → temp/backup/OS files, committed build artifacts, should-be-gitignored, misplaced/duplicate/empty files, large blobs, AI scratch, best-effort dead code (report-only by default; fixes on confirmation)
- **prime**: Create and audit prime files — AI context documents that index utility codebases
- **build knowledge base (bkb)**: LLM Knowledge Base builder → initialize, triage, ingest, query, lint, and maintain a persistent Markdown wiki compiled from raw sources
- **interview**: Run a structured elicitation interview against a prescriptive template → produces agent-ready operating artifacts (USER.md, SOUL.md, HEARTBEAT.md, plus machine-readable exports). First template: `work-operating-model`.
- **prompts**: Run reusable prompts from the library → `list`, `show <name>`, `run <name>`; extensible collection of battle-tested prompts for recurring jobs
- **commit**: Commit uncommitted files → analyzes, groups atomically, traces to REQs
- **inspect**: Explain uncommitted changes — what changed, why, and whether it's ready to commit (read-only)
- **version**: Show current version, last 5 releases, or check for upstream updates
- **recap**: Summary of last 5 completed user requests with their REQs
- **tutorial**: Interactive tutorials — quick start, concepts, workflow recipes, or guided tour

> **Core concept:** The capture-requests action always produces both a UR folder (preserving the original input) and REQ files (the queue items). Each REQ links back to its UR via `user_request` frontmatter. This pairing is mandatory for all requests — simple or complex.

> **Trail of Intent.** This skill doesn't just produce code — it produces a documented trail of *intent*. Every REQ records what the user wanted, why, and what "done" looks like — validated with the user during capture. As the REQ moves through the pipeline, builder decisions, scope declarations, and implementation notes are appended. The result is a living document tracing from original intent through every decision to final implementation. The UR preserves the verbatim request; the REQ preserves the validated, structured intent; the appended sections preserve how intent was realized. This trail is the skill's primary value — code is the side effect. Builder decisions during Step 6 implementation are guided by always-loaded behavioral guardrails (`crew-members/karpathy.md`): think before coding, simplicity first, surgical changes, goal-driven execution.

> **Capture ≠ Execute.** The capture-requests action captures requests. The work action executes them. These are strictly separate operations. After capture finishes writing files and reporting back, **STOP**. Do not start processing the queue, do not begin implementation, do not "helpfully" transition into the work action. The user decides when to execute — always. The only exception is if the user explicitly says something like "add this and then run it" or "capture this and start working" in the same invocation.

> **Human time has two optimal windows.** The system is designed to maximize the value of human attention:
>
> 1. **Capture phase** (capture-requests action) — The user is present, actively thinking about the request. This is the best time for back-and-forth: clarifying ambiguities, resolving contradictions, making scope decisions. Use the ask tool if your environment provides one; otherwise use your environment's normal ask-user prompt/tool. Every question must present concrete options — never open-ended "what do you mean?" prompts.
>
> 2. **Batch question review** (clarify questions action) — After the build phase completes everything it can without feedback, any remaining `pending-answers` REQs are surfaced as a batch. The user reviews all builder-decided questions together, confirms or adjusts, and resolved REQs re-enter the queue.
>
> Between these windows, the build phase runs autonomously. Builders never block on Open Questions — they mark them `- [~]` with best-judgment reasoning and create `pending-answers` follow-ups when they return via `do-work clarify`.

## Routing Decision

### Step 1: Parse the Input

Examine what follows "do-work":


Check these patterns **in order** — first match wins:

| Priority | Pattern                  | Example                                                                                                                            | Route                         |
| -------- | ------------------------ | ---------------------------------------------------------------------------------------------------------------------------------- | ----------------------------- |
| 1        | Empty, bare, or help     | `do-work`, `do-work help`                                                                                                          | → help menu                   |
| 2        | Version exact phrases    | `do-work check for updates`, `do-work check for update`                                                                            | → version                     |
| 3        | Pipeline keywords        | `do-work pipeline`, `do-work full`, `do-work pipeline add dark mode`, `do-work full add dark mode`, `do-work pipeline status`, `do-work pipeline abandon` | → pipeline                    |
| 4        | Action verbs (± REQ IDs ± flags) | `do-work run`, `do-work go`, `do-work start`, `do-work run REQ-042`, `do-work run --wave 0`                            | → work                        |
| 5        | Verify keywords          | `do-work verify`, `do-work verify-requests`, `do-work check REQ-018`, `do-work evaluate`, `do-work audit`, `do-work review requests` | → verify-requests              |
| 6        | Clarify keywords         | `do-work clarify`, `do-work questions`, `do-work pending`                                                                          | → clarify questions            |
| 7        | Code-review keywords     | `do-work code-review`, `do-work code review`, `do-work code-review prime-auth`, `do-work code review src/`, `do-work audit codebase`, `do-work review codebase`, `do-work codebase review` | → code-review                  |
| 8        | UI-review keywords       | `do-work ui-review`, `do-work ui-review src/`, `do-work review ui`, `do-work design review`, `do-work validate ui`, `do-work ui audit`, `do-work design audit` | → ui-review                    |
| 9        | Validate-feedback kw     | `do-work validate-feedback`, `do-work validate feedback`, `do-work triage findings`, `do-work triage feedback`, `do-work feedback review`, `do-work review feedback`, `do-work assess feedback`, `do-work should we push back` | → validate-feedback            |
| 10       | Review keywords          | `do-work review`, `do-work review-work`, `do-work review code`, `do-work audit code`                                               | → review-work                  |
| 11       | Present keywords         | `do-work present`, `do-work present-work`, `do-work showcase`, `do-work deliver`                                                   | → present-work                 |
| 12       | Cleanup keywords         | `do-work cleanup`, `do-work clean up`, `do-work tidy`, `do-work consolidate`, `do-work organize archive` — **excluding** any phrase that names memory/wiki/notes (`consolidate memory`, `clean up wiki`, `memory cleanup`, `lint and merge notes` → dream, priority 29) or that names stray/orphan/junk files or repo pollution (`clean up junk files`, `remove stray files`, `find orphan files` → stray-check, priority 30) | → cleanup                     |
| 13       | Commit keywords          | `do-work commit`, `do-work commit changes`, `do-work commit files`, `do-work save changes`, `do-work save work`                    | → commit                      |
| 14       | Inspect keywords         | `do-work inspect`, `do-work inspect REQ-005`, `do-work inspect UR-003`, `do-work explain changes`, `do-work what changed`, `do-work show changes` | → inspect                     |
| 15       | Version keywords         | `do-work version`, `do-work update`, `do-work what's new`, `do-work release notes`, `do-work what's changed`, `do-work updates`, `do-work history` | → version                     |
| 16       | Recap keywords           | `do-work recap`                                                                                                                    | → version                     |
| 17       | Forensics keywords       | `do-work forensics`, `do-work diagnose`, `do-work health check`, `do-work health`                                                  | → forensics                   |
| 18       | Roadmap keywords         | `do-work roadmap`, `do-work queue-status`, `do-work where are we`, `do-work what's left`, `do-work what's feasible`, `do-work what should I work on next` | → roadmap                     |
| 19       | Prime keywords           | `do-work prime`, `do-work prime create src/auth/`, `do-work prime audit`, `do-work create prime`, `do-work audit primes`           | → prime                       |
| 20       | BKB keywords             | `do-work bkb`, `do-work bkb init`, `do-work bkb ingest`, `do-work build knowledge base`, `do-work knowledge base`                 | → build knowledge base        |
| 21       | Interview keywords       | `do-work interview`, `do-work interview list`, `do-work interview work-operating-model`, `do-work interview <template> export`, `do-work interview <template> ingest`, `do-work elicit`, `do-work operating model` | → interview                    |
| 22       | Prompts keywords         | `do-work prompts`, `do-work prompts list`, `do-work prompts run architecture-decisions-log`, `do-work prompts show <name>`, `do-work prompt <name>` | → prompts                     |
| 23       | Quick-wins keywords      | `do-work quick-wins`, `do-work quick wins`, `do-work low-hanging`, `do-work scan`, `do-work scan src/`                             | → quick-wins                  |
| 24       | Scan-ideas keywords      | `do-work scan-ideas`, `do-work scan-ideas performance`, `do-work scan-ideas src/api/`, `do-work ideas`, `do-work brainstorm`, `do-work what should I build`, `do-work suggest`, `do-work ideate` | → scan-ideas                    |
| 25       | Deep-explore keywords    | `do-work deep-explore`, `do-work deep-explore performance`, `do-work explore concept`, `do-work deep dive`, `do-work develop idea`, `do-work deep-explore continue` | → deep-explore                  |
| 26       | Install keywords         | `do-work install ui-design`, `do-work install-ui-design`, `do-work install ui design`, `do-work install ui`, `do-work install frontend-design`, `do-work setup ui design`, `do-work setup design skill`, `do-work install bowser`, `do-work install-bowser`, `do-work install playwright`, `do-work install playwright-cli`, `do-work setup bowser`, `do-work setup playwright` | → install (target = `ui-design` or `bowser`) |
| 27       | Tutorial keywords        | `do-work tutorial`, `do-work tutorial quick-start`, `do-work tutorial concepts`, `do-work tutorial recipes`, `do-work tutorial tour` | → tutorial                      |
| 28       | Slop-check keywords      | `do-work slop-check`, `do-work slop check`, `do-work anti-slop`, `do-work slop-check do-work/deliverables/UR-003-client-brief.md`, `do-work slop-check REQ-042` | → slop-check                    |
| 29       | Dream keywords           | `do-work dream`, `do-work dream ./memory`, `do-work consolidate memory`, `do-work clean up wiki`, `do-work lint and merge notes`, `do-work memory cleanup` | → dream                         |
| 30       | Stray-check keywords     | `do-work stray-check`, `do-work stray-check src/`, `do-work find stray files`, `do-work find orphan files`, `do-work junk`, `do-work what doesn't belong`, `do-work file hygiene` | → stray-check                   |
| 31       | AI-report keywords       | `do-work ai-report`, `do-work ai-report UR-246`, `do-work ai report`, `do-work make-report`, `do-work make report`, `do-work screenshot-report`, `do-work visual report`, `do-work proof of work` | → ai-report                     |
| 32       | Note keywords            | `do-work note investigate xyz`, `do-work note "check Y before running"`, `do-work note add revisit after Z`, `do-work add note revisit after Z` | → note                          |
| 33       | Descriptive content      | `do-work capture-request: add dark mode`, `do-work [meeting notes]`, `do-work the button is broken`                                | → capture-requests              |


### Step 2: Preserve Payload

**Critical rule**: Never lose the user's content.

**Single-word rule**: A single word is either a known keyword or ambiguous — it is never "descriptive content."

- **Matches a keyword** in the routing table (e.g., "version", "verify", "cleanup") → route to that action directly.
- **Doesn't match any keyword** (e.g., "refactor", "optimize") → ambiguous. Ask: "Do you want to add '`{word}`' as a new request, or did you mean something else?"

Only route to **capture-requests** when the input is clearly descriptive — multiple words, a sentence, a feature request, etc.

If routing is genuinely unclear AND multi-word content was provided:

- Default to **capture-requests** (adding a task)
- Hold onto $ARGUMENTS
- If truly ambiguous, ask: "Add this as a request, or start the work loop?"
- User replies with just "add" or "work" → proceed with original content

### Verb Reference

| Route | Trigger verbs | Notes |
|-------|--------------|-------|
| **pipeline** | pipeline, full, full pipeline | Everything after keyword → `$ARGUMENTS` (request text). "pipeline status" → status mode. "pipeline abandon" → abandon mode. If no args and active pipeline exists, resume. |
| **work** | run, go, start, begin, work, process, execute, build, continue, resume | Optional REQ IDs after keyword (e.g., `do-work run REQ-042`) → process only those REQs. No args → work through full queue (dependency-aware order). Flag (default mode only): `--wave N` runs only REQs at dependency depth N. Strip the flag before extracting REQ IDs. `--wave` and explicit REQ IDs are mutually exclusive — reject with a clear error. An argument that is neither a `REQ-NNN` ID nor `--wave N` is **rejected** with a usage error — it is **not** treated as "no args" (which would silently build the whole queue). See `actions/work.md` Input → "Unrecognized arguments are rejected, not ignored." |
| **clarify** | clarify, answers, questions, pending, pending answers, blocked, what's blocked, what needs answers | Routes to `actions/clarify.md` |
| **verify-requests** | verify, verify-requests, verify requests, check, evaluate, review requests, review reqs, audit | "check" alone → verify; "check for updates" → version (priority 2); "audit" alone → verify; "audit codebase" → code-review; "audit primes" → prime |
| **code-review** | code-review, code review, review codebase, audit codebase, codebase review | Both hyphenated and unhyphenated forms route here, with or without scope. Scope args: prime file refs, directory paths, or combined |
| **ui-review** | ui-review, review ui, design review, validate ui, ui audit, design audit | Do NOT use "check ui" — consumed by verify at priority 5. Scope args: file paths, directory paths, prime file refs |
| **validate-feedback** | validate-feedback, validate feedback, triage findings, triage feedback, feedback review, review feedback, assess feedback, should we push back | Read-only triage of pasted external feedback / audit findings. Everything after the verb → `$ARGUMENTS` (the pasted findings). Per item verdict: Already done / Accept / Push back / Discuss, verified against code + git. Placed above review-work (priority 9) so feedback phrasings win before the bare "review" verb; triggers require feedback/findings/triage/push-back, so plain "review-work" is unaffected. |
| **review-work** | review, review-work, review work, review code, audit code, audit implementation, review REQ-NNN | "review requests" / "review reqs" → verify (priority 5), not here. "code review" → code-review (priority 7), not here. "review feedback" / "feedback review" → validate-feedback (priority 9), not here |
| **present-work** | present, present-work, present work, showcase, deliver, pitch, client brief | No target → most recent UR. "present all" → portfolio mode |
| **ai-report** | ai-report, ai report, make-report, make report, screenshot-report, visual report, proof of work | Self-contained folder report (`index.html` + `screenshots/`/`generated/`) with screenshots + SVG callouts + before/after + optional AI-generated diagrams; output to `ai-reports/`. Target arg: `UR-NNN`, `REQ-NNN`, "most recent", or empty. Distinct from present-work (educational explainer) and pipeline's `.single.html` (multi-REQ debrief) |
| **cleanup** | cleanup, clean up, tidy, consolidate, organize archive, fix archive | **Archive consolidation only.** A `clean up` / `consolidate` / `cleanup` that names **memory / wiki / notes** is the dream action (priority 29); one that names **stray / orphan / junk files** or repo pollution is the stray-check action (priority 30). Evaluate those carve-outs first; a bare verb with no memory/wiki/notes/junk target stays here. |
| **commit** | commit, commit changes, commit files, save changes, save work | |
| **inspect** | inspect, inspect changes, explain changes, what changed, show changes, describe changes | "what changed" (no apostrophe) → inspect; "what's changed" → version |
| **recap** | recap | Routes to version action with `mode: recap` |
| **version** | version, update, check for updates, what's new, release notes, what's changed, updates, history | "updates" (plural) shows last 5 releases; "update" (singular) triggers update check |
| **forensics** | forensics, diagnose, health check, health | |
| **roadmap** | roadmap, queue-status, where are we, what's left, what's feasible, what should I work on next | Optional scope: `pending`, `in-progress`, `done`, `UR-NNN`, or `since <date>`. Read-only. `"<action> status"` → that action (e.g., "pipeline status" → pipeline, "bkb status" → bkb, "interview <template> status" → interview). |
| **note** | note, note add, add note | Everything after the keyword → `$ARGUMENTS` (the note text). Strips a leading `add `; appends `- [YYYY-MM-DD] text` to `do-work/notes.md` (creates it on first use). Not a capture — no UR/REQ, no work loop, no commit. |
| **stray-check** | stray files, strays, orphan files, orphans, junk, what doesn't belong, file hygiene, stray-check | Repo-wide file hygiene (NOT do-work's own files — that's cleanup). Optional path scope; `fix` to apply fixes on confirmation, `report` to force read-only. Bare "clean up" / "consolidate" stays cleanup (priority 12). |
| **prime** | prime, prime create, prime audit, create prime, audit primes, primes | Everything after verb → `$ARGUMENTS`. "audit primes" → prime; plain "audit" → verify |
| **bkb** | bkb, build knowledge base, knowledge base, kb | Everything after verb → `$ARGUMENTS` (sub-command + params) |
| **interview** | interview, elicit, operating model | Everything after verb → `$ARGUMENTS` (`list`, `<template>`, or `<template> <sub-command>`). No args → help menu |
| **prompts** | prompts, prompt | Everything after verb → `$ARGUMENTS` (sub-command + prompt name + args). `prompts` (plural) and `prompt` (singular) both route here. First arg is the sub-command (`list`, `show`, `run`) or a prompt name (shorthand for `run`) |
| **quick-wins** | quick-wins, quick wins, low-hanging, low hanging fruit, scan, opportunities, what can we improve | "scan" alone or with a bare directory path → quick-wins; bare path = last meaningful token (any text after it is descriptive content → capture) |
| **install** | install, install ui-design, install-ui-design, install ui design, install ui, install frontend-design, setup ui design, setup design skill, install bowser, install-bowser, install playwright, install playwright-cli, setup bowser, setup playwright | **Normalize the invocation first**, then extract the target token: (a) strip a leading `install-` from hyphenated forms (`install-ui-design` → `ui-design`, `install-bowser` → `bowser`); (b) strip a leading `setup ` (`setup bowser` → `bowser`, `setup ui design` → `ui design`); (c) for `install X` with a space, take everything after `install `. Then map the normalized token: `ui-design`/`ui design`/`ui`/`frontend-design`/`design skill` → target `ui-design`; `bowser`/`playwright`/`playwright-cli` → target `bowser`. Pass the target as `$ARGUMENTS`. Bare `install` with no token (or unknown target after normalization) → help block. |
| **scan-ideas** | scan-ideas, ideate, ideas, brainstorm, what should I build, suggest, what's next, what could we improve | Everything after keyword → `$ARGUMENTS` (focus topic or directory). No args → open exploration |
| **deep-explore** | deep-explore, explore concept, deep dive, develop idea, explore idea | Everything after keyword → `$ARGUMENTS` (concept, file path, topic, or "continue"). No args → ask user what to explore |
| **tutorial** | tutorial, tutorial quick-start, tutorial concepts, tutorial recipes, tutorial tour, learn, getting started, how does this work | Everything after "tutorial" → `$ARGUMENTS` (mode). No args → ask user which mode |
| **slop-check** | slop-check, slop check, anti-slop | Everything after the verb → `$ARGUMENTS` (file path, REQ/UR ID, "most recent", or empty). Triggers must be distinctive — any `check ...` form (e.g., "check slop", "check draft", "check for slop") collides with verify priority 5 and is intentionally not listed here. |
| **dream** | dream, consolidate memory, clean up wiki, lint and merge notes, memory cleanup | Everything after the verb → `$ARGUMENTS` (memory directory path or empty). No path → auto-resolve default (`./memory`, `./wiki`, `./kb/wiki`, `./knowledge-base/wiki`). Single-word `dream` is a known keyword, not descriptive content. These memory/wiki/notes phrases take precedence over the generic cleanup verbs (priority 12) even though Dream sits lower in the table. |
| **capture-requests** | `capture-request:` prefix (the spaced `capture request:` form is also accepted), descriptive text, feature requests, bug reports, "add", "create", "I need", "we should" | Default for multi-word descriptive content that doesn't match any keyword |

## Examples

### Help Menu (bare invocation)

When invoked with no arguments or with `help` (`do-work`, `do-work help`), show a help menu with available actions and example prompts:

```
do-work — task queue for agentic coding tools

  Capture & pipeline:
    do-work capture-request: add dark mode to settings
    do-work pipeline add dark mode      End-to-end: investigate → capture → verify → run → review → present
    do-work pipeline status             Show progress / resume active pipeline

  Process the queue:
    do-work run                         Triage, build, test, review — one REQ at a time
    do-work clarify                     Review pending questions from completed work

  Verify & review:
    do-work verify-requests             Check capture quality against original input
    do-work review-work                 Review completed work (requirements + code + acceptance)
    do-work validate-feedback [findings] Triage external review feedback — accept / push back / already done
    do-work code-review [scope]         Standalone codebase review (prime refs, dirs, or both)
    do-work ui-review [scope]           Read-only UI quality validation
    do-work slop-check [target]         Validate a draft against the anti-slop principles before it ships

  Present & inspect:
    do-work present-work                Client brief, architecture, video, HTML explainer
    do-work ai-report [target]          Pixel-anchored HTML report: screenshots + SVG callouts + before/after
    do-work inspect                     Explain uncommitted changes (what, why, readiness)

  Scan & improve:
    do-work quick-wins [dir]            Refactoring opportunities and low-hanging tests
    do-work scan-ideas [focus]          Generate ideas for what to build next
    do-work deep-explore [concept]      Multi-round structured exploration of a concept
    do-work prime create src/auth/      Generate a prime file via interactive Q&A
    do-work prime audit                 Audit prime files for staleness and broken links

  Knowledge base:
    do-work bkb [sub]                   Sub-commands: init | triage | ingest | query | lint |
                                        resolve | close | status | defrag | garden | rollup | crew
    do-work dream [path]                Manual four-phase consolidation of a plain-text memory
                                        directory (orient, lint, heal, prune + reindex) — destructive

  Interviews:
    do-work interview                   Help menu
    do-work interview list              List available templates
    do-work interview <template>        Start or resume a structured elicitation interview
    do-work interview <template> review Run the cross-layer contradiction pass
    do-work interview <template> export Produce agent-ready operating artifacts

  Prompt library:
    do-work prompts                     Help menu
    do-work prompts list                List every available prompt
    do-work prompts show <name>         Print a prompt (read-only)
    do-work prompts run <name> [args]   Execute a prompt (e.g. architecture-decisions-log)

  Setup:
    do-work install ui-design           Frontend-design skill for production-grade UI
    do-work install bowser              Playwright CLI + Bowser for browser automation

  Maintenance & info:
    do-work cleanup                     Consolidate the archive
    do-work commit                      Analyze and commit files atomically
    do-work forensics                   Pipeline diagnostics — stuck work, orphaned URs
    do-work roadmap [scope]             Queue survey — ready/blocked/stale + TDD posture
    do-work queue-status                Alias for roadmap
    do-work note "investigate xyz"      Jot a dated next-step hint (surfaces atop roadmap)
    do-work stray-check [path]          Find orphan/junk files polluting the repo
    do-work version                     Version + last 5 releases
    do-work update                      Check for upstream updates
    do-work recap                       Last 5 completed URs with their REQs
    do-work tutorial                     Learn the skill (quick-start, concepts, recipes, tour)
    do-work help                        Show this menu

  Tip: add "help" to any command for details — e.g. do-work commit help
```

Do not ask "Start the work loop?" — just print the help menu and wait.

### Per-Command Help

When any action is invoked with `help` as its sole argument (e.g., `do-work commit help`, `do-work inspect help`), show a brief usage summary instead of executing the action.

**Actions with built-in help** (`pipeline`, `prime`, `bkb`): dispatch normally — they handle `help` internally.

**All other actions**: read the action file and present a compact summary:

```
<action-name> — <description from the blockquote>

  Usage:
    do-work <action> [args]       <brief description>

  Arguments:
    <list accepted arguments/modes from the action file's Input section>

  Examples:
    <2-3 example invocations>
```

Keep it short — no more than 15 lines. The goal is quick orientation, not a tutorial. After showing the summary, stop — do not execute the action.

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
| pipeline           | `./actions/pipeline.md`         | `$ARGUMENTS` (request text, "status", or "abandon") |
| capture-requests   | `./actions/capture.md`          | Full user input text           |
| work               | `./actions/work.md`             | `$ARGUMENTS` (REQ IDs, `--wave`, or empty) |
| clarify questions  | `./actions/clarify.md`          | (none needed)                  |
| verify-requests    | `./actions/verify-requests.md`  | Target UR/REQ or "most recent" |
| review-work        | `./actions/review-work.md`      | Target REQ/UR or "most recent" |
| validate-feedback  | `./actions/validate-feedback.md`| `$ARGUMENTS` (the pasted feedback / findings) |
| present-work       | `./actions/present-work.md`     | Target REQ/UR, "most recent", or "all" |
| ai-report          | `./actions/ai-report.md`        | Target REQ/UR, "most recent", or empty |
| cleanup            | `./actions/cleanup.md`          | (none needed)                  |
| commit             | `./actions/commit.md`           | (none needed)                  |
| inspect            | `./actions/inspect.md`          | Target REQ/UR or (none = all)  |
| code-review        | `./actions/code-review.md`      | Prime file refs and/or directory paths |
| ui-review          | `./actions/ui-review.md`        | File/directory paths and/or prime file refs |
| quick-wins         | `./actions/quick-wins.md`       | Target directory               |
| scan-ideas         | `./actions/scan-ideas.md`       | `$ARGUMENTS` (focus topic, directory, or empty) |
| deep-explore       | `./actions/deep-explore.md`     | `$ARGUMENTS` (concept, file path, topic, "continue", or empty) |
| install            | `./actions/install.md`          | `$ARGUMENTS` (target: `ui-design` or `bowser`) |
| forensics          | `./actions/forensics.md`        | (none needed)                  |
| roadmap            | `./actions/roadmap.md`          | `$ARGUMENTS` (optional scope: `pending`, `in-progress`, `done`, `UR-NNN`, `since <date>`) |
| note               | `./actions/note.md`             | `$ARGUMENTS` (the note text)   |
| stray-check        | `./actions/stray-check.md`      | `$ARGUMENTS` (optional path scope; `fix` / `report` mode token) |
| prime              | `./actions/prime.md`            | `$ARGUMENTS` (sub-command + params) |
| build knowledge base | `./actions/bkb.md`              | `$ARGUMENTS` (sub-command + params) |
| interview          | `./actions/interview.md`        | `$ARGUMENTS` (`list`, `<template>`, or `<template> <sub-command>`) |
| prompts            | `./actions/prompts.md`          | `$ARGUMENTS` (sub-command + prompt name + args) |
| version            | `./actions/version.md`          | `$ARGUMENTS`                   |
| recap              | `./actions/version.md`          | `mode: recap`                  |
| tutorial           | `./actions/tutorial.md`         | `$ARGUMENTS` (mode name or empty) |
| slop-check         | `./actions/slop-check.md`       | `$ARGUMENTS` (file path, REQ/UR ID, "most recent", or empty for newest deliverable) |
| dream              | `./actions/dream.md`            | `$ARGUMENTS` (memory directory path or empty for default resolution) |

### If subagents are available

Dispatch each action to a subagent. The subagent reads the action file and executes it — the main thread only sees the routing decision and the returned summary.

- **`work`**: Run in the background if your environment supports it. Print a status line (e.g., "Work queue processing in background...") and return control to the user immediately.
- **`cleanup`**: Run in the background if your environment supports it. Print a status line and return control to the user immediately.
- **Exception — pipeline dispatch**: When the pipeline action dispatches `work`, it runs in the **foreground** (blocking). The pipeline requires each step to complete before advancing. This override applies only when the pipeline is the caller.
- **`pipeline`, `capture-requests`, `clarify questions`, `verify-requests`, `review-work`, `validate-feedback`, `code-review`, `ui-review`, `present-work`, `ai-report`, `slop-check`, `dream`, `quick-wins`, `scan-ideas`, `deep-explore`, `prime`, `forensics`, `roadmap`, `note`, `stray-check`, `commit`, `inspect`, `install`, `version`, `recap`, `tutorial`, `prompts`, `interview`**: Run in the foreground (blocking). These need user interaction or produce small immediate output.
- **Screenshots (`capture-requests` only):** Subagents can't see images from the main conversation. Before dispatching, save screenshots to `do-work/user-requests/.pending-assets/screenshot-{n}.png`, write a text description of each, and include the paths + descriptions in the subagent prompt.

### If subagents are not available

Read the action file directly and follow its instructions in the current session. The action files are designed to work as standalone prompts — no subagent infrastructure required.

### On failure

Report the error to the user. Do not retry automatically.

## Suggest Next Steps

After every action completes, suggest the next logical prompts the user might want to run. See [`next-steps.md`](./next-steps.md) for the full per-action reference (what to suggest after each action, formatting rules, and constraints).
