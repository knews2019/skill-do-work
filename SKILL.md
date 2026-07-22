---
name: do-work
description: Task queue - add requests or process pending work
argument-hint: "pipeline [request] | capture-request: (describe a task) | run | verify-requests | review-work | code-review | ui-review | validate-feedback [paste findings] | present-work | ai-report [target] | slop-check [target] | dream [path] | clarify | abandon [REQ-NNN] | reserve [REQ-NNN for label] | release [REQ-NNN|label] | cleanup | commit | inspect | quick-wins | scan-ideas [focus] | deep-explore [concept] | prime [create|audit] | forensics | roadmap [scope] | board [mode] | note [text] | stray-check [path] | tidy-repo [path] | bkb [subcommand] | interview [template] | prompts [subcommand] | install [target] | version | recap | tutorial [mode] | help"
---

# Do-Work Skill

A unified entry point for task capture and processing. Each action's own file (see Action Dispatch) describes what it does and when to use it — this file only routes.

> **Core concept:** The capture-requests action always produces both a UR folder (preserving the original input) and REQ files (the queue items). Each REQ links back to its UR via `user_request` frontmatter. This pairing is mandatory for all requests — simple or complex.

> **Trail of Intent.** This skill doesn't just produce code — it produces a documented trail of *intent*. Every REQ records what the user wanted, why, and what "done" looks like — validated with the user during capture. As the REQ moves through the pipeline, builder decisions, scope declarations, and implementation notes are appended. The UR preserves the verbatim request; the REQ preserves the validated, structured intent; the appended sections preserve how intent was realized. This trail is the skill's primary value — code is the side effect. Builder decisions during implementation are guided by always-loaded behavioral guardrails (`crew-members/coding-guardrails.md`).

> **Capture ≠ Execute.** The capture-requests action captures requests. The work action executes them. After capture finishes writing files and reporting back, **STOP** — do not start processing the queue or "helpfully" transition into the work action. The user decides when to execute — always. The only exception: the user explicitly says something like "add this and then run it" in the same invocation.

> **Human time has two optimal windows.** (1) **Capture** — the user is present; this is the time for clarifying questions, always with concrete options, worded per `crew-members/clear-questions.md` (the contract for question wording, loaded before any interactive ask). (2) **Batch question review** — `do-work clarify` surfaces remaining `pending-answers` REQs for review together. Between these windows the build phase runs autonomously: builders never block on Open Questions — they mark them `- [~]` with best-judgment reasoning and create `pending-answers` follow-ups.

## Routing Decision

Examine what follows "do-work". Check the patterns below **in order** — first match wins. Trigger-verb lists are exhaustive for routing; argument handling beyond the Notes column lives in each action file's Input section.

| # | Trigger patterns | Route | Notes |
|---|---|---|---|
| 1 | empty, bare, `help` | help | Print the menu per `actions/help.md` and wait — never ask "start the work loop?" |
| 2 | `check for updates`, `check for update` (exact) | version | Wins before "check" hits verify (priority 5) |
| 3 | `pipeline`, `full` (± request text) | pipeline | Rest → `$ARGUMENTS`. `pipeline status` → status mode; `pipeline abandon` → abandon mode; no args + active pipeline → resume |
| 4 | `run`, `go`, `start`, `begin`, `work`, `process`, `execute`, `build`, `continue`, `resume` (± REQ IDs ± `--wave N`) | work | REQ IDs scope the run; `--wave` and IDs are mutually exclusive → reject; any other leftover token → usage error, never a silent full-queue run (`actions/work.md` Input) |
| 5 | `verify`, `verify-requests`, `verify requests`, `check`, `evaluate`, `audit`, `review requests`, `review reqs` | verify-requests | Bare `check`/`audit` land here; `audit codebase` → code-review (7); `audit primes` → prime (19) |
| 6 | `clarify`, `questions`, `answers`, `pending`, `pending answers`, `blocked`, `what's blocked`, `what needs answers` | clarify | |
| 7 | `code-review`, `code review`, `review codebase`, `audit codebase`, `codebase review` (± scope) | code-review | Scope args: prime file refs and/or directory paths |
| 8 | `ui-review`, `review ui`, `design review`, `validate ui`, `ui audit`, `design audit` (± scope) | ui-review | `check ui` is NOT a trigger — consumed by verify (5) |
| 9 | `validate-feedback`, `validate feedback`, `triage findings`, `triage feedback`, `feedback review`, `review feedback`, `assess feedback`, `should we push back` | validate-feedback | Rest → `$ARGUMENTS` (the pasted findings). Sits above review-work so feedback phrasings win before the bare `review` verb |
| 10 | `review`, `review-work`, `review work`, `review code`, `audit code`, `audit implementation`, `review REQ-NNN` | review-work | `review requests` → verify (5); `code review` → code-review (7); `review feedback` → validate-feedback (9) |
| 11 | `present`, `present-work`, `present work`, `showcase`, `deliver`, `pitch`, `client brief` | present-work | No target → most recent UR; `present all` → portfolio mode |
| 12 | `cleanup`, `clean up`, `tidy`, `consolidate`, `organize archive`, `fix archive` | cleanup | **Archive consolidation only.** Phrases naming memory/wiki/notes → dream (29); stray/orphan/junk files → stray-check (30); repo layout/structure/root → tidy-repo (34). Evaluate those carve-outs first; a bare verb stays here |
| 13 | `commit`, `commit changes`, `commit files`, `save changes`, `save work` | commit | |
| 14 | `inspect` (± REQ/UR id), `explain changes`, `what changed`, `show changes`, `describe changes` | inspect | `what changed` (no apostrophe) → inspect; `what's changed` → version (15) |
| 15 | `version`, `update`, `what's new`, `release notes`, `what's changed`, `updates`, `history` | version | `updates` (plural) shows last 5 releases; `update` (singular) triggers the update check |
| 16 | `recap` | version | Dispatch with `mode: recap` |
| 17 | `forensics`, `diagnose`, `health check`, `health` | forensics | |
| 18 | `roadmap`, `queue-status`, `where are we`, `what's left`, `what's feasible`, `what should I work on next` | roadmap | Optional scope: `pending`, `in-progress`, `done`, `UR-NNN`, `since <date>`. `"<action> status"` → that action (e.g. `pipeline status`, `bkb status`) |
| 19 | `prime`, `prime create …`, `prime audit`, `create prime`, `audit primes`, `primes` | prime | Rest → `$ARGUMENTS`; plain `audit` → verify (5) |
| 20 | `bkb`, `build knowledge base`, `knowledge base`, `kb` (± sub-command) | bkb | Rest → `$ARGUMENTS` |
| 21 | `interview`, `elicit`, `operating model` (± template ± sub-command) | interview | Rest → `$ARGUMENTS` (`list`, `<template>`, `<template> <sub-command>`); no args → its help menu |
| 22 | `prompts`, `prompt` (± sub-command ± name) | prompts | Rest → `$ARGUMENTS`; first arg is `list`/`show`/`run` or a prompt name (shorthand for `run`) |
| 23 | `quick-wins`, `quick wins`, `low-hanging`, `low hanging fruit`, `scan`, `opportunities`, `what can we improve` | quick-wins | `scan` alone or with a bare directory path lands here; bare path = last meaningful token — any text after it is descriptive content → capture (37) |
| 24 | `scan-ideas`, `ideas`, `ideate`, `brainstorm`, `what should I build`, `suggest`, `what's next`, `what could we improve` (± focus) | scan-ideas | Rest → `$ARGUMENTS` (topic or directory); no args → open exploration |
| 25 | `deep-explore`, `explore concept`, `deep dive`, `develop idea`, `explore idea` (± concept/path/`continue`) | deep-explore | Rest → `$ARGUMENTS`; no args → ask what to explore |
| 26 | `install <target>`, `install-<target>`, `setup <target>` | install | Normalize first: strip leading `install-`/`setup `; then map: `ui-design`/`ui design`/`ui`/`frontend-design`/`design skill` → `ui-design`; `bowser`/`playwright`/`playwright-cli` → `bowser`; `last30days`/`last 30 days` → `last30days`; `just-kanban`/`justfile`/`run-kanban` → `just-kanban`; `adhd`/`adhd mode`/`adhd-mode` → `adhd`. Pass target as `$ARGUMENTS`. Bare `install` or unknown target → help block. Canonical target list: `actions/install.md` |
| 27 | `tutorial` (± mode), `learn`, `getting started`, `how does this work` | tutorial | Rest → `$ARGUMENTS` (mode); no args → ask which mode |
| 28 | `slop-check`, `slop check`, `anti-slop` (± target) | slop-check | Rest → `$ARGUMENTS` (path, REQ/UR ID, `most recent`, or empty). `check …` forms are deliberately NOT triggers (collide with verify, 5) |
| 29 | `dream` (± path), `consolidate memory`, `clean up wiki`, `lint and merge notes`, `memory cleanup` | dream | Rest → `$ARGUMENTS` (memory dir; empty → auto-resolve `./memory`, `./wiki`, `./kb/wiki`, `./knowledge-base/wiki`). Single-word `dream` is a keyword, not descriptive content. These memory/wiki/notes phrases beat the generic cleanup verbs (12) |
| 30 | `stray-check` (± path), `stray files`, `strays`, `orphan files`, `orphans`, `junk`, `what doesn't belong`, `file hygiene` | stray-check | Repo junk, NOT do-work's own files (that's cleanup). Optional path scope; `fix` applies on confirmation, `report` forces read-only. Relocating legitimate files → tidy-repo (34) |
| 31 | `ai-report`, `ai report`, `make-report`, `make report`, `screenshot-report`, `visual report`, `proof of work` (± target) | ai-report | Target: `UR-NNN`, `REQ-NNN`, `most recent`, or empty. Distinct from present-work (explainer) and pipeline's debrief |
| 32 | `note <text>`, `note add <text>`, `add note <text>` | note | Rest → `$ARGUMENTS` (strips a leading `add `); appends `- [YYYY-MM-DD] text` to `do-work/notes.md`. Not a capture — no UR/REQ, no work loop, no commit |
| 33 | `board` (± mode), `kanban`, `kanban board`, `queue board`, `visualize queue`, `show the board` | board | Mode: empty/`serve`/`live` → live board at `:8090`; `static`/`generate`/`html` → static HTML; `summary`/`status` → column counts. Needs the Go toolchain |
| 34 | `tidy-repo`, `tidy repo`, `file-reorg` (legacy), `reorg`, `reorganize`, `restructure`, `declutter`, `tidy layout`, `fix the layout`, `clean up the root` (± path ± `plan`) | tidy-repo | Junk deletion → stray-check (30); do-work bookkeeping → cleanup (12); code architecture changes → the work pipeline |
| 35 | `abandon`, `cancel`, `wont-do`, `won't do` — only with empty args or a `REQ-NNN` ID | abandon | Rest → `$ARGUMENTS` (REQ IDs + optional reason). `abandon`/`cancel` + ID-less prose is descriptive content → capture (37). `pipeline abandon` already matched (3). Bare verb → list cancellable REQs and ask |
| 36 | `reserve` (± REQ IDs ± `for <label>`), `release`, `unreserve` (± REQ IDs or label) | reserve | Rest → `$ARGUMENTS` — but for the `release`/`unreserve` triggers pass `release <rest>` so the action enters release mode. Allocates pending REQs to another worktree/cloud session; bare `reserve` lists reservations with staleness flags |
| 37 | `capture-request:` / `capture request:` prefix, or descriptive multi-word content (feature requests, bug reports, "add …", "I need …") | capture-requests | The default for multi-word descriptive content that matches no keyword |

**Single-word rule:** a single word is either a known keyword (route it) or ambiguous — never "descriptive content." For an unknown single word (e.g. "refactor"), ask: "Do you want to add '`{word}`' as a new request, or did you mean something else?"

## Payload Preservation Rules

**Never lose the user's content.** When clarification is needed but content was provided:

1. Keep the full `$ARGUMENTS` payload in context — never ask the user to re-paste.
2. Ask a simple question: "Add this as a request, or start the work loop?"
3. Accept minimal replies ("add" / "work") and apply the chosen action to the stored content.

## Action Dispatch

Each action has an action file with full instructions. How you execute it depends on your environment's capabilities.

| Action             | Action file                     | Context to pass                |
|--------------------|---------------------------------|--------------------------------|
| help               | `./actions/help.md`             | `$ARGUMENTS` (empty, or `<action> help`) |
| pipeline           | `./actions/pipeline.md`         | `$ARGUMENTS` (request text, "status", or "abandon") |
| capture-requests   | `./actions/capture.md`          | Full user input text           |
| work               | `./actions/work.md`             | `$ARGUMENTS` (REQ IDs, `--wave`, or empty) |
| clarify questions  | `./actions/clarify.md`          | (none needed)                  |
| abandon            | `./actions/abandon.md`          | `$ARGUMENTS` (REQ IDs + optional reason) |
| reserve            | `./actions/reserve.md`          | `$ARGUMENTS` (REQ IDs + `for <label>`, `release …`, or empty to list) |
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
| install            | `./actions/install.md`          | `$ARGUMENTS` (install target — canonical list in `./actions/install.md`) |
| forensics          | `./actions/forensics.md`        | (none needed)                  |
| roadmap            | `./actions/roadmap.md`          | `$ARGUMENTS` (optional scope: `pending`, `in-progress`, `done`, `UR-NNN`, `since <date>`) |
| note               | `./actions/note.md`             | `$ARGUMENTS` (the note text)   |
| stray-check        | `./actions/stray-check.md`      | `$ARGUMENTS` (optional path scope; `fix` / `report` mode token) |
| tidy-repo          | `./actions/tidy-repo.md`        | `$ARGUMENTS` (optional path scope; `plan` / `--plan-only`; `file-reorg` is a legacy alias) |
| prime              | `./actions/prime.md`            | `$ARGUMENTS` (sub-command + params) |
| build knowledge base | `./actions/bkb.md`              | `$ARGUMENTS` (sub-command + params) |
| interview          | `./actions/interview.md`        | `$ARGUMENTS` (`list`, `<template>`, or `<template> <sub-command>`) |
| prompts            | `./actions/prompts.md`          | `$ARGUMENTS` (sub-command + prompt name + args) |
| version            | `./actions/version.md`          | `$ARGUMENTS`                   |
| recap              | `./actions/version.md`          | `mode: recap`                  |
| tutorial           | `./actions/tutorial.md`         | `$ARGUMENTS` (mode name or empty) |
| slop-check         | `./actions/slop-check.md`       | `$ARGUMENTS` (file path, REQ/UR ID, "most recent", or empty for newest deliverable) |
| dream              | `./actions/dream.md`            | `$ARGUMENTS` (memory directory path or empty for default resolution) |
| board              | `./actions/board.md`            | `$ARGUMENTS` (mode: `serve` / `static` / `summary`) |

**Per-command help:** any action invoked with `help` as its sole argument (e.g. `do-work commit help`) routes to `actions/help.md` per-command mode instead of executing — except `pipeline`, `prime`, and `bkb`, which handle `help` internally (dispatch those normally).

### If subagents are available

Dispatch each action to a subagent. The subagent reads the action file and executes it — the main thread only sees the routing decision and the returned summary.

- **`work`**, **`cleanup`**: run in the background if supported; print a status line and return control to the user.
- **`board`** (`serve` mode): run in the background if supported — it's a long-running local server; print the URL (`http://localhost:8090`). The `static` and `summary` modes run in the foreground.
- **Exception — pipeline dispatch**: when the pipeline action dispatches `work`, it runs in the **foreground** (blocking); the pipeline requires each step to complete before advancing.
- **All other actions**: run in the foreground (blocking) — they need user interaction or produce small immediate output.
- **Screenshots (`capture-requests` only):** subagents can't see images from the main conversation. Before dispatching, save screenshots to `do-work/user-requests/.pending-assets/screenshot-{n}.png`, write a text description of each, and include the paths + descriptions in the subagent prompt.

### If subagents are not available

Read the action file directly and follow its instructions in the current session — action files are designed to work as standalone prompts.

### On failure

Report the error to the user. Do not retry automatically.

## Suggest Next Steps

After every action completes, suggest the next logical prompts the user might want to run. See [`next-steps.md`](./next-steps.md) for the full per-action reference (what to suggest after each action, formatting rules, and constraints).
