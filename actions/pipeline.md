# Pipeline Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to run a full end-to-end pipeline: investigate, capture, verify, run, review, present. Manages state across sessions via `do-work/pipeline.json`.

A stateful multi-action orchestration that chains six actions in sequence. Each step dispatches to an existing action. The pipeline tracks progress in a JSON state file, supports resume across sessions, and reports status at each transition.

## Philosophy

- **One command, full cycle.** The user describes what they want; the pipeline handles the rest.
- **Resumable by design.** If the session ends mid-pipeline (context limit, crash, user closes terminal), re-invoking the pipeline picks up from where it left off. The state file is the source of truth.
- **Orchestrator only.** The pipeline never re-implements action logic. It dispatches to existing actions and tracks which ones have completed.
- **Coexists with CHECKPOINT.md.** The pipeline tracks macro-steps (which action to run next). The work action's CHECKPOINT.md tracks micro-state within a single `do work run` invocation. Both systems operate independently.

## Input

`$ARGUMENTS` determines behavior:

| Input | Mode |
|-------|------|
| `do work pipeline {request text}` | Initialize a new pipeline with the given request |
| `do work full {request text}` | Alias — same as `pipeline {request text}` |
| `do work pipeline` (no args, active pipeline exists) | Resume from the first pending step |
| `do work pipeline status` | Show current pipeline status without advancing |
| `do work pipeline abandon` | Deactivate the current pipeline without completing it |
| `do work pipeline` (no args, no active pipeline) | Show pipeline help |

## State File

Pipeline state lives at `do-work/pipeline.json`. Created on initialize, read on every subsequent invocation.

```json
{
  "session_id": "2026-04-08-001",
  "request": "add dark mode to settings panel",
  "started_at": "2026-04-08T11:00:00Z",
  "active": true,
  "steps": [
    { "name": "investigate", "status": "done",    "completed_at": "2026-04-08T11:01:00Z" },
    { "name": "capture",     "status": "done",    "completed_at": "2026-04-08T11:02:00Z", "artifacts": ["REQ-042", "UR-018"] },
    { "name": "verify",      "status": "pending", "completed_at": null },
    { "name": "run",         "status": "pending", "completed_at": null },
    { "name": "review",      "status": "pending", "completed_at": null },
    { "name": "present",     "status": "pending", "completed_at": null }
  ]
}
```

**Field definitions:**

| Field | Type | Description |
|-------|------|-------------|
| `session_id` | string | `YYYY-MM-DD-NNN` — date + incrementing counter |
| `request` | string | The user's original request text |
| `started_at` | string | ISO 8601 timestamp when the pipeline was initialized |
| `active` | boolean | `true` while pipeline is running, `false` when completed or abandoned |
| `steps` | array | Ordered list of pipeline steps with status tracking |
| `steps[].name` | string | Step identifier: `investigate`, `capture`, `verify`, `run`, `review`, `present` |
| `steps[].status` | string | `pending`, `in-progress`, `done`, or `failed` |
| `steps[].completed_at` | string\|null | ISO 8601 timestamp when step finished, or null |
| `steps[].artifacts` | array\|undefined | REQ/UR IDs produced (capture step only) |
| `steps[].error` | string\|undefined | Error description (failed steps only) |

## Steps

### Step 1: Determine Mode

1. Check if `do-work/pipeline.json` exists and has `"active": true`
2. Check `$ARGUMENTS` for content

| Active pipeline? | $ARGUMENTS | Mode |
|-----------------|------------|------|
| No | Has request text | **Initialize** (Step 2) |
| No | Empty or "status" | **Help** (show help menu and stop) |
| Yes | Has request text | **Conflict** — warn the user. Ask: "A pipeline is already active. Resume it, or abandon it and start fresh?" |
| Yes | Empty | **Resume** (Step 3) |
| Yes | "status" | **Status** (print status block and stop) |
| Yes | "abandon" | **Abandon** — set `active: false`, print final status, stop |

### Step 2: Initialize (new pipeline)

1. Create `do-work/` directory if it doesn't exist
2. Generate `session_id`: today's date + `-001` (the counter is a label for readability — since only one pipeline can be active at a time, incrementing is not required)
3. Write `do-work/pipeline.json` with:
   - `request` set to `$ARGUMENTS` (the request text, stripped of the "pipeline" or "full" keyword)
   - All 6 steps set to `status: "pending"`
   - `active: true`
   - `started_at` set to current ISO 8601 timestamp
4. **Exclude state file from git**: If a `.gitignore` exists in the project root and doesn't already contain `do-work/pipeline.json`, append it. If no `.gitignore` exists, create one containing `do-work/pipeline.json`. The state file is transient session state and should not be committed.
5. Print the initial status block
6. Proceed to Step 4 (execute first step: `investigate`)

### Step 3: Resume (existing pipeline)

1. Read `do-work/pipeline.json`
2. Find the first step where `status` is `"pending"` or `"in-progress"` or `"failed"`
   - `in-progress` means the previous session ended mid-step — retry it
   - `failed` means a previous attempt failed — retry it
3. Print the status block showing current progress
4. Proceed to Step 4 with that step as the current step

### Step 4: Execute Current Step

For the current step:

1. **Update state**: Set the step's `status` to `"in-progress"` in `pipeline.json` — write the file immediately before dispatching
2. **Dispatch** to the corresponding action:

| Pipeline step | Action to dispatch | What to pass | Context from prior steps |
|---------------|-------------------|--------------|--------------------------|
| `investigate` | the inspect action (`do work inspect`) | No arguments | None — inspects all uncommitted changes. If there are no uncommitted changes, the inspect action will report that and this step completes immediately (it's a pre-flight check, not a blocker). |
| `capture` | the capture action (`do work capture request: {request}`) | The `request` field from pipeline.json | None — request text is the input |
| `verify` | the verify requests action (`do work verify requests`) | Target UR from capture artifacts | Pass the UR ID from the capture step's `artifacts` (e.g., `do work verify UR-018`) |
| `run` | the work action (`do work run`) | REQ IDs from capture artifacts | Pass the specific REQ IDs from the capture step's `artifacts` (e.g., `do work run REQ-042`). The sub-agent prompt MUST instruct the work action to process ONLY these REQs, then stop — do NOT drain the full queue. |
| `review` | the review work action (`do work review work`) | Target REQ/UR from capture artifacts | Pass the UR ID from the capture step's `artifacts` (e.g., `do work review UR-018`) so the reviewer knows which work to review |
| `present` | the present work action (`do work present work`) | Target UR from capture artifacts | Pass the UR ID from the capture step's `artifacts` (e.g., `do work present UR-018`) so the deliverables target this pipeline's work. If the capture step produced no artifacts (empty `artifacts` array), skip this step — mark it `done` with no artifacts and proceed to completion. |

Dispatch each action the same way the main router dispatches actions: subagent if available, inline otherwise. The pipeline action is the orchestrator — it calls the router's dispatch mechanism, not the action files directly.

**Sub-agent context rule:** Sub-agents do not inherit conversation history. When dispatching via sub-agent, always read `pipeline.json` and include in the sub-agent prompt: (1) the pipeline request text, (2) all artifact IDs from completed steps, and (3) any relevant file paths. Without this, the sub-agent won't know which UR was just created or which REQs to target.

**Foreground dispatch override:** All pipeline-dispatched actions run in the foreground (blocking), even if SKILL.md normally marks them as background (e.g., `work`). The pipeline requires synchronous completion of each step before advancing to the next.

3. **After the action completes**:
   - Update the step's `status` to `"done"` and set `completed_at` to current timestamp
   - **Capture step only**: Parse the action's output for created REQ and UR IDs (e.g., `REQ-042`, `UR-018`). Store them in the step's `artifacts` array. If IDs cannot be parsed, leave `artifacts` as an empty array — do not block the pipeline.
   - Write the updated `pipeline.json` immediately
   - Print the updated status block

4. **Advance**: If more steps have `status: "pending"`, proceed to the next one (loop back to the top of Step 4). If all steps are `"done"`, proceed to Step 5.

### Step 5: Completion

When all 6 steps are done:

1. Set `active: false` in `pipeline.json`
2. Write the final `pipeline.json`
3. Print the completion status block (all checkmarks) — use the **Completion Status Block** format from the Output Format section (includes Duration, Branch, Verdict metadata)
4. **Assemble the Pipeline Completion Report data**. This is the primary user-education artifact. Pull data from:
   - **Final summary table**: Each completed REQ's frontmatter (`id`, `title`, `commit`, `domain`) and the REQ's `## Implementation Summary` → one-line synthesis. Group rows by domain so related work sits together.
   - **Test state (before → after)**: Each REQ's `## Testing` section records what tests were added/run. Aggregate test-suite counts if the work action logged them (e.g., "Go: 81 → 98"). If no before-baseline was captured, show only the post-state and note "baseline not recorded" rather than inventing numbers.
   - **Cross-REQ coherence highlights**: Pull from the review step's output — the reviewer validates that interacting REQs (shared files, shared symbols, shared subsystems) remained consistent. Include each coherence assertion with the REQ pair. If the review didn't produce coherence notes (single-REQ pipelines, Route A REQs), omit this section.
   - **Carry-forward work (implied, not captured yet)**: Scan for (a) REQs with `status: pending-answers`, (b) `## Lessons Learned` sections mentioning deferred items, (c) TODO/FIXME comments introduced by the pipeline's commits. List them as candidates for a follow-up capture — but **do NOT auto-capture them**; the user decides.
   - **Deliverables**: Paths produced by the `present` step (`do-work/deliverables/{UR-NNN}-client-brief.md`, `-video/`, `-interactive-explainer.html`). Read from the present step's artifacts if recorded, or glob `do-work/deliverables/` for matches scoped to this pipeline's UR.
   - **How to verify**: Concrete commands the user can copy-paste — `git show {sha}` for each commit, the project's test command(s), and the path to open the interactive explainer. This is the validation recipe.
5. **Render the report in all three formats and save each to disk.** The same data powers three deliverables, each tuned to a different audience and surface:

   | File | Format | Audience | Template |
   |------|--------|----------|----------|
   | `do-work/deliverables/{UR-NNN}-pipeline-summary.md`      | Plain markdown          | Developer reading in a terminal or editor — cat / grep / paste into a PR | Output Format → **Plain Markdown Report** |
   | `do-work/deliverables/{UR-NNN}-pipeline-summary.marp.md` | Marp presentation slides | Stakeholder sitting through a walkthrough — viewed with `marp --preview` or exported to PDF | Output Format → **Marp Slide Deck** |
   | `do-work/deliverables/{UR-NNN}-pipeline-summary.html`    | Standalone HTML         | Non-technical reader browsing independently — Mermaid + Tailwind via CDN, zero build | Output Format → **Standalone HTML Debrief** |

   All three files carry the same facts — one rendering per surface, no format-specific editorializing. Use `{REQ-id}-pipeline-summary.*` as the filename prefix if no UR was captured.

6. **Print the plain-markdown rendering to stdout** so the user sees the debrief immediately. Reference the other two paths in the closing "Deliverables" section so they can open them next.
7. **Queue continuation check**: Scan `do-work/queue/REQ-*.md` for files with `status: pending` in their frontmatter. Exclude any REQ IDs listed in the current pipeline's `artifacts` array (those should already be completed). If remaining pending REQs exist, proceed to Step 5a. If the queue is empty, suggest next steps and stop.

**Proportional depth:** Single-REQ Route A pipelines (config tweak, docs) get a minimal report in all three formats — status block + 1-row Final summary + How to verify. The Marp deck collapses to 3–4 slides; the HTML collapses to a single-screen summary. Multi-REQ or Route B/C pipelines get the full treatment shown above. Don't pad short pipelines with empty sections, and don't truncate long ones. Match the report to the scope of the work.

### Step 5a: Queue Continuation

When the pipeline completes and additional pending REQs remain in the queue (from prior captures, follow-ups created during review, or other sources):

1. Print the continuation notice (see Output Format) listing each pending REQ ID and its title
2. Record the list of pending REQ IDs about to be processed (e.g., `["REQ-043", "REQ-044"]`) — this is needed for review targeting in step 3
3. Dispatch the work action (`do work run`) in **standard queue-draining mode** — do NOT scope to pipeline artifacts. Pass the pending REQ IDs (e.g., `do work run REQ-043 REQ-044`) so the work action processes them.
4. After the work action completes, dispatch the review work action for each REQ from step 2 individually (e.g., `do work review REQ-043`, then `do work review REQ-044`). Always use REQ IDs — never pass a UR ID, since UR-scoped review would re-review all completed REQs under that UR, not just this cycle's batch.
5. **Loop**: Scan `do-work/queue/REQ-*.md` for `status: pending` again. If more pending REQs remain (e.g., follow-ups created during the review step), repeat from step 1. If the queue is empty, print "Queue fully drained." and suggest next steps.

**Max iterations:** The continuation loop runs at most **3 cycles**. If pending REQs still remain after 3 run → review cycles, stop the loop and print:

```
Continuation limit reached (3 cycles). {count} REQ(s) still pending:
  {REQ-ID} — {title}
  ...

Run "do work run" to continue processing manually.
```

This prevents runaway loops when review steps keep generating follow-up REQs.

**Error handling:** If the work action or review action fails during continuation:

1. Report the error to the user with context about what failed
2. Print how many REQs were successfully processed before the failure
3. Stop the continuation loop — do not retry automatically
4. Suggest the appropriate recovery command:
   - **Run step failed**: Suggest `do work run` to resume processing pending REQs
   - **Review step failed**: Suggest `do work review REQ-NNN` for each REQ that was processed but not yet reviewed (those REQs are already `status: completed`, so `do work run` would be a no-op)

Unlike the main pipeline's error handling (Step 6), continuation errors do not update `pipeline.json` — the formal pipeline is already complete.

**State file note:** The pipeline's `active` field remains `false` during the continuation — the formal pipeline is complete. The continuation is a post-pipeline queue drain. If the session ends mid-continuation, the user can resume with `do work run` to process any remaining pending REQs.

### Step 6: Error Handling

If any step's dispatched action fails (error, exception, or the action reports failure):

1. Set the step's `status` to `"failed"` and add an `error` field with a brief description
2. Write the updated `pipeline.json`
3. Print the status block showing the failure point
4. Report the error to the user with context about what failed and why
5. Leave `active: true` — the user can fix the issue and resume with `do work pipeline`

On resume, the pipeline retries the failed step from scratch.

## Output Format

### Status Block (printed after every step transition)

```
── Pipeline ─────────────────────────
  ✓ investigate   done
  ✓ capture       done  → REQ-042, UR-018
  ✓ verify        done
  ◎ run           in progress...
  ○ review        pending
  ○ present       pending
─────────────────────────────────────
  Session: 2026-04-08-001
  Request: add dark mode to settings panel
```

**Symbols:**
- `✓` — done
- `◎` — in progress
- `✗` — failed
- `○` — pending

For the `capture` step, append ` → {artifact IDs}` after "done" if artifacts were recorded.

### Completion Status Block (printed when all 6 steps are done)

```
── Pipeline {session_id} — COMPLETE ──
  ✓ investigate   done
  ✓ capture       done  → UR-018, 12 REQs
  ✓ verify        done  (or "skipped ({reason})" if inputs were pre-verified)
  ✓ run           done  → {N} commits
  ✓ review        done  ({verdict: PASS / PASS with caveats / FAIL})
  ✓ present       done  → {N} deliverables
─────────────────────────────────────
  Session:   {session_id}
  Duration:  {elapsed time from started_at to now, e.g. "~10.5h end-to-end"}
  Branch:    {current git branch} ({pushed | local only})
  Verdict:   {PASS | PASS with caveats | FAIL}
```

### Pipeline Completion Report — three renderings of one dataset

The same facts — Final summary, Test state, Coherence, Carry-forward, Deliverables, How to verify — are rendered three ways. One pass over the data, three files on disk. Never author any of the three from scratch if another already exists; re-render from the source data so they stay consistent.

**Composition rules (apply to all three formats):**

- **Serve both audiences in every file.** Each summary opens with a "What got built" narrative for the reader who has no clue, then transitions into the audit data for the reader who wants receipts. A stakeholder landing on the `.html` should understand the feature without opening any other file; a developer scanning the `.md` should reach the commits and test deltas within seconds. Never ship a summary that only audits or only educates.
- **Reuse client-brief content verbatim.** When the `present` step ran, the "What got built" narrative and architecture diagram come from `{UR-NNN}-client-brief.md` — copy the same sentences and the same diagram. Paraphrasing across files introduces drift. If the brief doesn't exist (present skipped or produced nothing), synthesize from the REQ Implementation Summaries.
- **Cite commits, not prose.** Every audit claim should trace to a commit SHA, a REQ ID, or a file path. Tables and bullet lists with pointers beat paragraphs of explanation. (The opening narrative is the exception — it's plain language for the no-clue reader.)
- **Pull from primary sources.** Final summary rows come from REQ frontmatter; coherence notes come from the review step's actual output; test deltas come from what `run` and `review` logged. Do not invent metrics.
- **Be honest about gaps.** If the baseline test count wasn't captured before the pipeline started, write "baseline not recorded" — don't guess. If no cross-REQ coherence was analyzed (single-REQ pipeline), omit that section.
- **Carry-forward ≠ auto-capture.** List candidates clearly with the command the user would run to capture each one, but never capture them automatically.
- **No format-specific editorializing.** The Marp deck must not add facts the markdown lacks; the HTML must not soften or strengthen claims for a broader audience. Format dictates rendering; rendering does not dictate facts.

#### 1. Plain Markdown Report — `{UR-NNN}-pipeline-summary.md`

Developer-facing. Read in a terminal with `cat`, grepped, or pasted into a PR description. No YAML header, no CSS, no slide breaks — just markdown.

```markdown
# Pipeline Completion Report — {UR-NNN}

**Session**: {session_id} · **Duration**: {duration} · **Branch**: {branch} ({pushed|local})
**Verdict**: {PASS | PASS with caveats | FAIL}

## What got built (for the reader who has no clue)

[2-3 plain-language sentences synthesizing the UR title, the REQs' What sections, and the client brief's "What We Built" paragraph — no jargon, no commit SHAs, no REQ IDs. The reader learns what the feature *does* before they see the audit trail. If the `present` step ran, pull this straight from `{UR-NNN}-client-brief.md` so the two files stay in sync; if it didn't run, synthesize from the REQ Implementation Summaries.]

[Optional: reuse the ASCII architecture diagram from the client brief, verbatim. Skip if the work is non-architectural (config tweak, bug fix, docs).]

**Go deeper:** [`{UR-NNN}-client-brief.md`](./{UR-NNN}-client-brief.md) · [`{UR-NNN}-interactive-explainer.html`](./{UR-NNN}-interactive-explainer.html) *(only include links that actually exist on disk)*

## Final summary

| REQ | Commit | Scope | One-line |
|-----|--------|-------|----------|
| REQ-402 | 5ab214d | docs     | 4 lessons-learned files + prime links |
| REQ-410 | 9371a68 | refactor | shared `initializeDatabaseAtPath` — prod/test converged |
| REQ-413 | 9e20bde | backend  | SHA-256 index + O(log N) lookup rewrite |
| ...     | ...     | ...      | ... |

## Test state (before → after the {N}-REQ pipeline)

| Suite         | Before    | After     | Delta |
|---------------|-----------|-----------|-------|
| Go (sa1-server) | 81 tests  | 98 tests  | +17 |
| Frontend      | 1053 tests / 62 suites | 1067 tests / 65 suites | +14 tests / +3 suites |
| `go vet`      | clean     | clean     | — |

## Cross-REQ coherence highlights (verified by the review)

- **REQ-413 ↔ REQ-406**: early-exit preserved at cache-hit + fresh-match. `effectiveLimit` threaded.
- **REQ-413 ↔ REQ-407**: metric_version filter preserved in `loadCachedEdgesForSource`.
- **REQ-411 ↔ REQ-412**: orthogonal, zero file overlap, no shared state.

## Carry-forward work (implied, not captured yet)

- [Deferred item] — capture with `do work capture request: ...`
- [TODO/FIXME introduced and left for a follow-up]
- [`pending-answers` REQs awaiting user input — run `do work clarify`]

## Deliverables

Render each bullet as a relative markdown link to the file (e.g. `[...]({UR-NNN}-client-brief.md)`) so a reader opening the `.md` in GitHub, a PR, or an editor can click through to any sibling artifact. Group by audience so the reader lands on the right surface first.

**For the clueless-reader (start here if you don't know what was built):**

- [`{UR-NNN}-client-brief.md`](./{UR-NNN}-client-brief.md) — plain-language brief with architecture diagram + value prop *(if present ran)*
- [`{UR-NNN}-interactive-explainer.html`](./{UR-NNN}-interactive-explainer.html) — interactive Before/After explainer, open in any browser *(if present ran)*
- [`{UR-NNN}-video/`](./{UR-NNN}-video/) — Remotion video walkthrough (`cd` in, `npm install`, `npm run preview`) *(if present ran)*

**For the developer / reviewer (audit the run):**

- [`{UR-NNN}-pipeline-summary.md`](./{UR-NNN}-pipeline-summary.md) — this report (markdown)
- [`{UR-NNN}-pipeline-summary.marp.md`](./{UR-NNN}-pipeline-summary.marp.md) — Marp slide deck (`marp --preview`)
- [`{UR-NNN}-pipeline-summary.html`](./{UR-NNN}-pipeline-summary.html) — standalone HTML debrief

## How to verify

1. **Check out the branch and pull latest:**
   ```
   git checkout {branch} && git pull
   ```
2. **Inspect each commit** (ordered to show the build-up):
   ```
   git show 5ab214d   # REQ-402 — lessons-learned docs
   git show 9371a68   # REQ-410 — shared init routine
   git show 9e20bde   # REQ-413 — SHA-256 index rewrite
   ```
3. **Run the tests** (matches what the pipeline ran):
   ```
   {project test command — e.g., `go test ./...` and `npm test`}
   ```
4. **Preview the other renderings:**
   ```
   npx @marp-team/marp-cli {UR-NNN}-pipeline-summary.marp.md --preview
   open do-work/deliverables/{UR-NNN}-pipeline-summary.html
   ```
5. **Read the per-REQ archive** for the full trail of intent:
   ```
   do-work/archive/{UR-NNN}/REQ-*.md
   ```
```

#### 2. Marp Slide Deck — `{UR-NNN}-pipeline-summary.marp.md`

Stakeholder-facing. Viewed with `marp --preview` or exported to PDF/HTML. Must start with Marp YAML frontmatter (`marp: true`). Each slide separated by `---`. Keep slides scannable — no slide should fit more than ~8 rows of content; split long Final-summary tables across domain-grouped slides. Use a Mermaid `graph LR` on the coherence slide when there are 2+ cross-REQ links.

Required slide sequence (omit a slide entirely if its section has no data — don't leave empty slides):

1. **Title slide** — UR-NNN, session ID, branch, verdict badge
2. **What got built** — 2-3 plain-language bullets pulled from the client brief's "What We Built" section. No commit SHAs, no REQ IDs. This is the slide a stakeholder who wandered in late needs to orient. Skip only if no `present` step produced a brief AND the UR itself is trivially self-explanatory from its title.
3. **How it works** (when a client brief exists with an architecture diagram) — reuse the ASCII or Mermaid diagram from the brief. Skip for non-architectural changes.
4. **At-a-glance stats** — REQ count, commit count, test delta, duration (big numbers in a 2×2 or 4-column grid)
5. **What shipped — {domain}** — one slide per domain bucket (docs / backend / refactor / frontend / tests). Each is a table of REQ / commit / one-line for that domain only.
6. **Test state (before → after)** — the table, full-width
7. **Cross-REQ coherence** — Mermaid `graph LR` diagram of interacting REQs (skip for single-REQ pipelines)
8. **Coherence assertions** — verbatim review quotes, one bullet per assertion
9. **Carry-forward work** — bullets with capture commands (skip if none)
10. **How to verify** — fenced `bash` block with checkout + git-show + test commands
11. **Deliverables + next steps** — two-column layout: left column "Start here if you want to understand what was built" lists the client brief, interactive explainer, and video (when present ran); right column "Audit the run" lists the markdown and HTML summary siblings. Render each as the bare filename — stakeholders open the deck from `do-work/deliverables/`, so relative filenames are all they need to find the sibling files in the same folder.

Use this Marp frontmatter skeleton and extend the `style:` block as needed — don't invent new themes:

```yaml
---
marp: true
theme: default
paginate: true
size: 16:9
header: '{UR-NNN} — Pipeline Debrief'
footer: 'Session {session_id} · branch {branch}'
style: |
  section { font-family: system-ui, -apple-system, sans-serif; }
  h2 { color: #1e40af; border-bottom: 2px solid #e2e8f0; padding-bottom: 0.25em; }
  code { background: #f1f5f9; padding: 0.1em 0.3em; border-radius: 3px; }
  table { font-size: 0.75em; }
  th { background: #1e40af; color: white; }
  .big { font-size: 3em; font-weight: 700; color: #1e40af; }
  .label { color: #64748b; font-size: 0.9em; }
---
```

#### 3. Standalone HTML Debrief — `{UR-NNN}-pipeline-summary.html`

Non-technical-reader-facing. Single `.html` file, zero build steps. Same content as the markdown, rendered for a browser.

**Stack (CDN only — no npm, no build):**
- Tailwind CSS via `<script src="https://cdn.tailwindcss.com"></script>`
- Mermaid via `<script type="module">` import of `mermaid@10` from jsDelivr
- Vanilla JS only (no React, no framework)

**Required sections (in order):**

1. **Hero** — UR-NNN as H1, one-paragraph description, metadata badges (branch, duration, verdict)
2. **What got built** — a prose block that educates a reader who has no clue: 2-3 sentences explaining what the feature does in plain language, pulled from the client brief's "What We Built" section. No REQ IDs, no commit SHAs. When a client brief exists, this section is the primary educational entry point for a non-technical reader arriving at the HTML. Skip only if no `present` step ran and the UR is self-explanatory from its title.
3. **How it works** (when an architecture diagram exists in the client brief) — a `<div class="mermaid">` with a `graph TD` or `graph LR` rendering of the same components/data flow from the brief's architecture section. Caption each node in plain language. Skip entirely for non-architectural changes (config tweaks, bug fixes).
4. **At-a-glance stat cards** — 4-column grid of big-number stats (REQ count, commits, tests added, suites added)
5. **What shipped** — grouped sections by domain, each with a styled table of REQ / commit / one-line
6. **Test state** — the table, styled with the accent colour for the After column and green for the Delta
7. **How the work holds together** — a `<div class="mermaid">` containing the same `graph LR` from the Marp deck (Mermaid renders on load)
8. **Coherence assertions** — responsive card grid, one card per assertion, with the REQ pair in mono accent and the claim below (skip the whole section for single-REQ pipelines)
9. **Carry-forward work** — cards with a bold title, muted explanation, and the capture command in a `<pre>` block (skip if none)
10. **How to verify** — numbered headings, each followed by a copy-pasteable `<pre><code>` block
11. **Related deliverables** — a navigation card grid **before** the final follow-ups list, splitting cross-links by audience. Left card group "Understand what was built" with real `<a href="./{UR-NNN}-client-brief.md">` / `<a href="./{UR-NNN}-interactive-explainer.html">` / `<a href="./{UR-NNN}-video/">` anchors (only include tiles for artifacts that actually exist on disk — if present ran and produced them). Right card group "Audit the run" linking the markdown (`<a href="./{UR-NNN}-pipeline-summary.md">`) and Marp deck (`<a href="./{UR-NNN}-pipeline-summary.marp.md">`) siblings. The HTML is the most discoverable surface for a non-technical reader — it must point them to the deeper, more educational artifacts.
12. **Footer / next steps** — ordered list with `do work present {UR-NNN}` and other follow-ups

**Design requirements:**

- Light theme default; dark theme via `@media (prefers-color-scheme: dark)` overriding CSS custom properties on `:root`
- Palette: CSS variables for `--bg`, `--surface`, `--text`, `--muted`, `--accent`, `--accent-soft`, `--border`. Light: white/slate-50 / slate-900 / blue-600. Dark: slate-900 / slate-100 / blue-400.
- Font: `system-ui, -apple-system, sans-serif`
- Max content width: `max-w-6xl` centred
- Generous spacing (`py-10` / `py-16` on sections) — readable like a long-form article, not cramped like a dashboard
- Mermaid init: `mermaid.initialize({ startOnLoad: true, theme: 'default', securityLevel: 'loose' })`

**What NOT to do:**

- Don't add charts the source data doesn't support (no fabricated time-series, no fake percentages)
- Don't embed images unless the REQs reference them
- Don't pull in additional CDN scripts beyond Tailwind + Mermaid — the file must work offline once cached
- Don't add interactivity that hides data (collapsible sections are fine; JS-gated sections that require a click to reveal facts are not)

### Continuation Notice (printed when pending REQs remain after pipeline completion)

```
── Queue Continuation ───────────────
  {count} pending REQ(s) remaining:
    {REQ-ID} — {title}
    {REQ-ID} — {title}
    ...

  Processing remaining queue...
─────────────────────────────────────
```

When the continuation loop finishes and the queue is empty:

```
Queue fully drained. All pending requests processed.
```

When the continuation loop hits the max iteration cap (3 cycles):

```
Continuation limit reached (3 cycles). {count} REQ(s) still pending:
  {REQ-ID} — {title}
  ...

Run "do work run" to continue processing manually.
```

### Help Menu (no active pipeline, no arguments)

```
pipeline — full end-to-end orchestration

  Start a new pipeline:
    do work pipeline add dark mode to settings
    do work full add dark mode to settings

  Resume / manage:
    do work pipeline            Resume an active pipeline
    do work pipeline status     Show pipeline progress
    do work pipeline abandon    Deactivate without completing

  Steps (executed in order):
    1. investigate   Inspect uncommitted changes
    2. capture       Capture the request as REQ + UR files
    3. verify        Verify capture quality
    4. run           Process the queue (build, test, review)
    5. review        Post-work code review + acceptance testing
    6. present       Generate client-facing deliverables (brief, diagrams, video, HTML)
```

## Rules

- **Never skip steps.** The pipeline always runs in order: investigate → capture → verify → run → review → present. Steps cannot be reordered or omitted. (Exception: if `capture` produced no artifacts, `present` has nothing to deliver — mark it `done` with no artifacts and proceed.)
- **One pipeline at a time.** If an active pipeline exists, the user must complete, resume, or abandon it before starting a new one.
- **Orchestrator only.** The pipeline dispatches to existing actions. It never re-implements capture, work, verify, review, present, or inspect logic. Each action runs exactly as it would if the user invoked it directly.
- **Write state before dispatch.** Always update `pipeline.json` to `"in-progress"` before dispatching an action, and to `"done"` after it completes. This ensures the state file reflects reality even if the session ends unexpectedly.
- **The `run` step may be long.** The work action processes only this pipeline's captured REQs but may still take significant time for complex requests. When starting this step, note: "Starting queue processing — this may take a while if multiple REQs are pending."
- **Platform-agnostic.** No tool-specific APIs. Dispatch actions the same way the main router does. If your environment supports stop hooks, you can optionally install `hooks/pipeline-guard.sh` to prevent accidental stops mid-pipeline — but the pipeline works without it.
- **Do not commit the state file.** `do-work/pipeline.json` is transient session state. It tracks a single pipeline run and has no value after completion. Ensure it is in `.gitignore`.
- **Pass context to sub-agents explicitly.** Sub-agents have no conversation history. When dispatching a step via sub-agent, always include the pipeline request text and all artifact IDs from completed steps in the sub-agent prompt. Without this, sub-agents cannot target the correct UR/REQs.
- **Scope the `run` step to captured REQs only.** The work action is queue-draining by default. When dispatched from the pipeline, it must only process the REQs created by this pipeline's capture step (listed in `artifacts`). Never process unrelated backlog items during a pipeline run.
- **Drain remaining queue after completion.** After the pipeline's 6 steps finish, check for other pending REQs in the queue. If any exist, continue processing them automatically via run + review cycles until the queue is empty. This continuation uses standard queue-draining mode (not scoped to pipeline artifacts) and does not re-run `present` per cycle — the user can run `do work present all` after the queue drains if they want a portfolio summary. The pipeline state file remains `active: false` — the continuation is a post-pipeline operation. Maximum 3 continuation cycles — if REQs still remain after 3 cycles, stop and let the user continue manually.
- **Suggest next steps on completion.** After the pipeline finishes (including any queue continuation), suggest what the user might want to do next (see the next-steps reference).
- **Completion is education, not a checkmark.** When all steps finish, produce the full Pipeline Completion Report (Final summary, Test state, Cross-REQ coherence, Carry-forward work, Deliverables, How to verify) in **all three formats** — plain markdown (`{UR-NNN}-pipeline-summary.md`), Marp slide deck (`{UR-NNN}-pipeline-summary.marp.md`), and standalone HTML (`{UR-NNN}-pipeline-summary.html`). One dataset, three renderings, different audiences. A 12-REQ pipeline that prints only "Pipeline complete" — or writes only the markdown and skips the deck and the HTML — wastes the user's opportunity to understand and validate what shipped. Match report depth to pipeline scope — minimal for Route A, full for multi-REQ URs.
- **Never author from scratch when re-rendering.** The three report files share one source of truth: the data you extracted in Step 5.4. If you find yourself phrasing the same claim differently across formats, stop and re-render from the data. Divergence between the markdown, Marp, and HTML versions is a bug.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "Skip verify — the capture was clean" | Run every step in sequence, even if you think it's unnecessary | The pipeline's value is consistency — skipping steps creates gaps |
| "The pipeline is stuck — just mark the step done" | Investigate why it's stuck, then fix or escalate | Marking stuck steps as done creates hollow completions |
| "I'll restart the pipeline from scratch" | Resume from the last completed step using pipeline.json state | Restarting loses progress and re-runs already-completed work |
| "This step failed — skip it and continue" | Record the failure, attempt recovery, then decide with the user | Silently skipping failures undermines the pipeline's reliability |
| "Skip present — the user can run it later" | Run present as part of the pipeline | The pipeline closes the loop: code → review → deliverables. Skipping present leaves the work uncommunicated. |
| "Just print 'Pipeline complete' — the user ran it, they know what happened" | Produce the full Pipeline Completion Report with Final summary, Test state, Coherence, Carry-forward, Deliverables, and How to verify | The user kicked off a pipeline that may have run for hours across many REQs. They need a digest they can scan, verify against, and share — not a checkmark. |
| "Invent a test-count baseline — it's probably close" | Write "baseline not recorded" when you lack real numbers | Fabricated metrics in the Completion Report erode trust in every future report |
| "The markdown version is enough — who needs Marp and HTML?" | Write all three renderings — `.md`, `.marp.md`, and `.html` — from the same dataset | Different audiences consume the work on different surfaces: a developer scans the `.md` in a PR, a stakeholder sits through the deck, a non-technical reader browses the HTML. Shipping only the markdown leaves two audiences unserved. |
| "I'll write the Marp deck now and re-author the HTML later" | Render all three files in the same completion pass, from the same extracted data | Sequencing the renderings invites drift — the second render subtly rephrases claims from the first, and the three files stop agreeing |
| "The summary is for devs — the client brief is for stakeholders. They don't need to cross-link." | Put a "What got built" narrative at the top of every summary, and cross-link the client brief, explainer, and summaries from each other's footer | Readers arrive from whichever file a teammate sent them. If the summary only audits and the brief only educates, half the readers bounce instead of drilling in |
| "Writing the 'What got built' section in the summary duplicates the client brief" | Copy the client brief's "What We Built" paragraph verbatim and add a link back to the full brief | Duplication is the point — each file must stand alone for readers who land on it first. Cross-links give the deeper context |

## Red Flags

- Step marked as `in-progress` but no recent git activity (pipeline may be stuck)
- Step marked as `done` but no artifacts recorded (hollow completion)
- Pipeline has been active for >24 hours without step transitions
- pipeline.json shows a step status that contradicts the file system state
- Multi-REQ pipeline finished but the Completion Report is missing sections that should have content (e.g., no Final summary table for a 10-REQ run, no Deliverables list after present)
- Completion Report cites test counts without specifying which suite they came from, or includes "Before" numbers that were never actually measured
- Only one or two of the three rendering files (`.md`, `.marp.md`, `.html`) exist in `do-work/deliverables/`
- The three renderings disagree — e.g., the Marp deck lists 11 REQs but the markdown lists 12, or the HTML shows a test delta the markdown doesn't
- The HTML file references external scripts beyond Tailwind + Mermaid, or requires `npm install` to render
- The Marp deck is missing its `marp: true` frontmatter header or uses a custom theme name (`marp --preview` will fail silently)
- Any of the three summary formats opens straight into the audit (Final summary table, stat cards) with no "What got built" narrative — the clueless reader has no entry point
- The three formats disagree on the "What got built" narrative — summary `.html` says one thing, client brief says another
- Summary files list sibling deliverables as plain text paths instead of rendered links; readers can't click through

## Verification Checklist

- [ ] pipeline.json state matches file system reality
- [ ] Current step status updated before and after execution
- [ ] Artifacts from each completed step recorded in pipeline.json
- [ ] Failed steps have failure reason documented
- [ ] User informed of pipeline progress at each step transition
- [ ] Pipeline Completion Report rendered in all three formats: `{UR-NNN}-pipeline-summary.md`, `{UR-NNN}-pipeline-summary.marp.md`, and `{UR-NNN}-pipeline-summary.html` — all present in `do-work/deliverables/`
- [ ] Plain markdown rendering printed to stdout so the user sees it immediately
- [ ] All three renderings cite the same commit SHAs, same test deltas, and same REQ count (no drift between formats)
- [ ] Report's Final summary cites real commit SHAs from each REQ's frontmatter
- [ ] Report's Test state labels "baseline not recorded" when no before-measurement was captured — no invented numbers
- [ ] Report's How to verify section contains copy-pasteable commands (not abstract instructions)
- [ ] All three summary formats open with a plain-language "What got built" section before any audit data, pulled from the client brief when it exists
- [ ] Summary Deliverables section renders as clickable relative links, grouped by audience (understand-what-was-built vs. audit-the-run)
- [ ] Each summary file's "What got built" narrative matches the client brief word-for-word (no paraphrasing drift)
- [ ] Marp file starts with `marp: true` YAML frontmatter and uses `---` slide separators
- [ ] HTML file is fully standalone: Tailwind + Mermaid via CDN only, no other external dependencies, no build step required
