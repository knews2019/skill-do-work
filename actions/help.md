# Help Action

> **Part of the do-work skill.** Invoked on a bare `do-work` / `do-work help` (routing priority 1) or when any action gets `help` as its sole argument. Prints orientation and stops — never executes an action. This file exists so the menu text loads only on help routes instead of taxing every invocation via SKILL.md.

## Steps

### Step 1: Mode selection

- Empty `$ARGUMENTS` or `help` → print the full menu (Step 2), then wait. Do not ask "Start the work loop?" — just print and wait.
- `<action> help` → per-command summary (Step 3). Exception: `pipeline`, `prime`, and `bkb` have built-in help — SKILL.md dispatches those to their own action files, so they never arrive here.

### Step 2: Full menu (bare invocation)

```
do-work — task queue for agentic coding tools

  Capture & pipeline:
    do-work capture-request: add dark mode to settings
    do-work pipeline add dark mode      End-to-end: investigate → capture → verify → run → review → present
    do-work pipeline status             Show progress / resume active pipeline

  Process the queue:
    do-work run                         Triage, build, test, review — one REQ at a time
    do-work clarify                     Review pending questions from completed work
    do-work abandon REQ-042 [why]       Mark a REQ won't-do — cancelled + archived, shows with done work

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
    do-work prime audit                 Health-check primes + refresh their Stakes (writes Stakes)

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
    do-work install last30days          Engagement-ranked social-research engine (vendored, keyless)
    do-work install just-kanban         Justfile recipes for the queue-kanban board (just run-kanban)

  Maintenance & info:
    do-work cleanup                     Consolidate the archive
    do-work commit                      Analyze and commit files atomically
    do-work forensics                   Pipeline diagnostics — stuck work, orphaned URs
    do-work roadmap [scope]             Queue survey — ready/blocked/stale + TDD posture
    do-work queue-status                Alias for roadmap
    do-work board [mode]                Kanban board of the queue — live (serve) / static / summary (needs Go)
    do-work note "investigate xyz"      Jot a dated next-step hint (surfaces atop roadmap)
    do-work stray-check [path]          Find orphan/junk files polluting the repo
    do-work tidy-repo [path]            Tidy the repo layout safely (plan → approve → move → verify)
    do-work version                     Version + last 5 releases
    do-work update                      Check for upstream updates
    do-work recap                       Last 5 completed URs with their REQs
    do-work tutorial                     Learn the skill (quick-start, concepts, recipes, tour)
    do-work help                        Show this menu

  Tip: add "help" to any command for details — e.g. do-work commit help
```

### Step 3: Per-command help

Read the target action's file and present a compact summary:

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

## Rules

- Help never executes an action, never touches the queue, never commits.
- Keep the menu in sync with SKILL.md's routing table and Action Dispatch table — a menu line for a verb that no longer routes (or a missing line for one that does) is drift; `do-work tutorial`'s recipe spot-check applies here too.
