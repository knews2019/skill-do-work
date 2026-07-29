# Pipeline Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to run a full end-to-end pipeline: investigate, capture, verify, run, review, present. Manages state across sessions via `do-work/pipeline.json`.

A stateful multi-action orchestration that chains six actions in sequence. Each step dispatches to an existing action. The pipeline tracks progress in a JSON state file, supports resume across sessions, and reports status at each transition.

## Philosophy

- **One command, full cycle.** The user describes what they want; the pipeline handles the rest.
- **Resumable by design.** If the session ends mid-pipeline (context limit, crash, user closes terminal), re-invoking the pipeline picks up from where it left off. The state file is the source of truth.
- **Orchestrator only.** The pipeline never re-implements action logic. It dispatches to existing actions and tracks which ones have completed.
- **Coexists with CHECKPOINT.md.** The pipeline tracks macro-steps (which action to run next). actions/work.md's CHECKPOINT.md tracks micro-state within a single `do-work run` invocation. Both systems operate independently.

## When to Use

**Use when:**
- The user describes a task and wants it taken end-to-end: investigate → capture → verify → build → review → present, without manual hand-offs between steps.
- A prior pipeline was interrupted and the user wants to resume (`do-work pipeline` with no args + active state file).
- The user wants deterministic step-by-step progression with persistent state across sessions.

**Do NOT use when:**
- An active pipeline already exists and the user wants to start an unrelated one — finish or `abandon` the current pipeline first.
- The task is too small to benefit from the full six-step cycle — a one-liner fix is fine with plain `capture` + `run`.
- See `SKILL.md` routing table for sibling action selection.

## Input

`$ARGUMENTS` determines behavior:

| Input | Mode |
|-------|------|
| `do-work pipeline {request text}` | Initialize a new pipeline with the given request |
| `do-work full {request text}` | Alias — same as `pipeline {request text}` |
| `do-work pipeline` (no args, active pipeline exists) | Resume from the first pending step |
| `do-work pipeline status` | Show current pipeline status without advancing |
| `do-work pipeline abandon` | Deactivate the current pipeline without completing it |
| `do-work pipeline` (no args, no active pipeline) | Show pipeline help |

## State File

Pipeline state lives at `do-work/pipeline.json` — its full lifecycle in this action: created at Step 2 (Initialize), read at the top of every subsequent invocation, and rewritten at every step transition (Steps 4, 5, 5a, 6). Nothing in this workflow ever deletes it: Step 5 (Completion) and the Abandon mode (Step 1) both set `active: false` and leave the file on disk as a historical record — it's excluded from git (Step 2 substep 4), not removed from disk.

The full JSON shape, the field-by-field definitions, and the pretty-print invariant that every write to this file must honor live in `actions/pipeline-reference.md` → **State File Schema** — read it before the first write or read of `pipeline.json` this session.

## Steps

### Step 1: Determine Mode

1. Check if `do-work/pipeline.json` exists and has `"active": true`.

2. **Normalize `$ARGUMENTS` and bucket it before consulting the mode table.** Trim leading/trailing whitespace and lowercase the result. The normalized input falls into one of three buckets:
   - **Empty** — no remaining content
   - **Reserved keyword** — exact match (post trim+lowercase) against `status` or `abandon`
   - **Request text** — anything else

3. **Near-miss guard for reserved keywords.** Before treating an unrecognized single-token input as request text, check whether it's a prefix of a reserved keyword. If normalized `$ARGUMENTS` is a single token of ≤8 characters and is a strict prefix of `status` or `abandon` (e.g., `stat`, `statu`, `aban`, `abandn`, `abandon-it`), respond with:

   ```
   Did you mean 'status' / 'abandon'? Pipeline keywords require an exact match — got '{token}'.
   ```

   and stop. **Apply this guard in both active and inactive pipeline states** — do not let the Conflict arm (active) or the Initialize arm (inactive) swallow a near-miss as request text. Same shape as the `install-` / `setup ` normalization in `SKILL.md`'s install routing row: normalize and disambiguate before dispatch, so an "almost match" cannot silently route to the wrong action.

4. Consult the mode table:

| Active pipeline? | `$ARGUMENTS` bucket | Mode |
|-----------------|---------------------|------|
| No | Request text | **Initialize** (Step 2) |
| No | Empty, `status` keyword, or `abandon` keyword | **Help** (show the help menu — `actions/pipeline-reference.md` → **Help Menu** — and stop; there's nothing to inspect or abandon without an active pipeline; both reserved keywords are no-ops here) |
| Yes | Request text | **Conflict** — warn the user. Ask: "A pipeline is already active. Resume it, or abandon it and start fresh?" |
| Yes | Empty | **Resume** (Step 3) |
| Yes | `status` keyword | **Status** (print status block and stop) |
| Yes | `abandon` keyword | **Abandon** — set `active: false`, print final status, stop |

### Step 2: Initialize (new pipeline)

1. Create `do-work/` directory if it doesn't exist
2. Generate `session_id`: today's date + `-001` (the counter is a label for readability — since only one pipeline can be active at a time, incrementing is not required)
3. Write `do-work/pipeline.json` with:
   - `request` set to `$ARGUMENTS` (the request text, stripped of the "pipeline" or "full" keyword)
   - All 6 steps set to `status: "pending"`
   - `active: true`
   - `started_at` set to current ISO 8601 timestamp
4. **Exclude state file from git**: Ensure `do-work/pipeline.json` is ignored regardless of install layout by appending it to the enclosing repo's `.git/info/exclude` (local-only — never committed, never shipped) when it isn't already ignored. Do **not** touch the project's committable `.gitignore`: a root-extract install's shipped `.gitignore` can't reach `do-work/pipeline.json` from a nested `.claude/skills/do-work/` install, and the host project shouldn't carry a committed ignore rule for transient state. Use the exact snippet from `crew-members/background-agents.md` → **Local-ignore snippet (for genuinely-transient paths)**, substituting `do-work/pipeline.json` for `<path>` (appended pattern: `**/do-work/pipeline.json`). Live state is a must-never-commit path, so per that snippet's tracked-path caveat ALSO run `git ls-files -- do-work/pipeline.json` — the ignore rule is powerless over an indexed file, and `check-ignore` passes on one. A non-empty result means an earlier session committed the state file before this guard existed: report the remedy to the user (`git rm --cached do-work/pipeline.json`) — never untrack on their behalf — and continue the pipeline.
5. Print the initial status block (`actions/pipeline-reference.md` → **Status Block**)
6. Proceed to Step 4 (execute first step: `investigate`)

### Step 3: Resume (existing pipeline)

1. Read `do-work/pipeline.json`
2. Find the first step where `status` is `"pending"` or `"in-progress"` or `"failed"`
   - `in-progress` means the previous session ended mid-step — retry it
   - `failed` means a previous attempt failed — retry it
3. Print the status block showing current progress (`actions/pipeline-reference.md` → **Status Block**)
4. Proceed to Step 4 with that step as the current step

### Step 4: Execute Current Step

For the current step:

1. **Update state**: Set the step's `status` to `"in-progress"` in `pipeline.json` — write the file immediately before dispatching
2. **Dispatch** to the corresponding action:

| Pipeline step | Action to dispatch | What to pass | Context from prior steps |
|---------------|-------------------|--------------|--------------------------|
| `investigate` | actions/inspect.md (`do-work inspect`) | No arguments | None — inspects all uncommitted changes. If there are no uncommitted changes, actions/inspect.md will report that and this step completes immediately (it's a pre-flight check, not a blocker). |
| `capture` | actions/capture.md (`do-work capture-request: {request}`) | The `request` field from pipeline.json | None — request text is the input |
| `verify` | actions/verify-requests.md (`do-work verify-requests`) | Target UR from capture artifacts | Pass the UR ID from the capture step's `artifacts` (e.g., `do-work verify UR-018`) |
| `run` | actions/work.md (`do-work run`) | REQ IDs from capture artifacts | Pass the specific REQ IDs from the capture step's `artifacts` (e.g., `do-work run REQ-042`). The sub-agent prompt MUST instruct actions/work.md to process ONLY these REQs, then stop — do NOT process the full queue. |
| `review` | actions/review-work.md (`do-work review-work`) | Target REQ/UR from capture artifacts | Pass the UR ID from the capture step's `artifacts` (e.g., `do-work review UR-018`) so the reviewer knows which work to review |
| `present` | actions/present-work.md (`do-work present-work`) | Target UR from capture artifacts | Pass the UR ID from the capture step's `artifacts` (e.g., `do-work present UR-018`) so the deliverables target this pipeline's work. If the capture step produced no artifacts (empty `artifacts` array), skip this step — mark it `done` with no artifacts and proceed to completion. |

Dispatch each action the same way the main router dispatches actions: subagent if available, inline otherwise. actions/pipeline.md is the orchestrator — it calls the router's dispatch mechanism, not the action files directly.

**Sub-agent context rule:** Sub-agents do not inherit conversation history. When dispatching via sub-agent, always read `pipeline.json` and include in the sub-agent prompt: (1) the pipeline request text, (2) all artifact IDs from completed steps, and (3) any relevant file paths. Without this, the sub-agent won't know which UR was just created or which REQs to target.

**Durability:** When a dispatched action itself fans work out to background or parallel sub-agents, follow the durability pattern in `crew-members/background-agents.md` (disk-durable run directory as source of truth; survives a dead orchestrator session).

**Foreground dispatch override:** All pipeline-dispatched actions run in the foreground (blocking), even if SKILL.md normally marks them as background (e.g., `work`). The pipeline requires synchronous completion of each step before advancing to the next.

3. **After the action completes**:
   - Update the step's `status` to `"done"` and set `completed_at` to current timestamp
   - **Capture step only**: Parse the action's output for created REQ and UR IDs (e.g., `REQ-042`, `UR-018`). Store them in the step's `artifacts` array. If IDs cannot be parsed, leave `artifacts` as an empty array — do not block the pipeline.
   - Write the updated `pipeline.json` immediately
   - Print the updated status block (`actions/pipeline-reference.md` → **Status Block**)

4. **Advance**: If more steps have `status: "pending"`, proceed to the next one (loop back to the top of Step 4). If all steps are `"done"`, proceed to Step 5.

### Step 5: Completion

When all 6 steps are done:

1. Set `active: false` in `pipeline.json`
2. Write the final `pipeline.json`
3. Print the completion status block (all checkmarks) — use the **Completion Status Block** template in `actions/pipeline-reference.md` (includes Duration, Branch, Verdict metadata)
4. **Assemble the Pipeline Completion Report data** (`actions/pipeline-reference.md` → **Pipeline Completion Report — three renderings of one dataset**, which also holds the **Composition rules** every rendering must follow). This is the primary user-education artifact. Pull data from:
   - **Final summary table**: Each completed REQ's frontmatter (`id`, `title`, `commit`, `domain`) and the REQ's `## Implementation Summary` → one-line synthesis. Group rows by domain so related work sits together.
   - **Test state (before → after)**: Each REQ's `## Testing` section records what tests were added/run. Aggregate test-suite counts if actions/work.md logged them (e.g., "Go: 81 → 98"). If no before-baseline was captured, show only the post-state and note "baseline not recorded" rather than inventing numbers.
   - **Cross-REQ coherence highlights**: Pull from the review step's output — the reviewer validates that interacting REQs (shared files, shared symbols, shared subsystems) remained consistent. Include each coherence assertion with the REQ pair. If the review didn't produce coherence notes (single-REQ pipelines, Route A REQs), omit this section.
   - **Carry-forward work (implied, not captured yet)**: Scan for (a) REQs with `status: pending-answers`, (b) `## Lessons Learned` sections mentioning deferred items, (c) TODO/FIXME comments introduced by the pipeline's commits. List them as candidates for a follow-up capture — but **do NOT auto-capture them**; the user decides.
   - **Deliverables**: Paths produced by the `present` step (`do-work/deliverables/{UR-NNN}-client-brief.md`, `-video/`, `-interactive-explainer.single.html`). Read from the present step's artifacts if recorded, or glob `do-work/deliverables/` for matches scoped to this pipeline's UR.
   - **How to verify**: Concrete commands the user can copy-paste — `git show {sha}` for each commit (for a merge-commit SHA — `git rev-parse --verify -q '{sha}^2'` succeeds — emit `git show --first-parent -m {sha}` instead, since plain `git show` on a merge prints a combined diff that is usually empty), the project's test command(s), and the path to open the interactive explainer. This is the validation recipe.
5. **Load `crew-members/anti-slop.md` before rendering.** The three authored formats are the user's primary debrief — apply the principles (lead with the verdict, compress, cite evidence, match depth to pipeline scope, disclose anything unverified). One dataset, three renderings, no padding.
6. **Render the three LLM formats from one dataset, then export the Marp deck to HTML.** Four files total per completion — three authored, one mechanical:

   | File | Format | Audience | Template / Producer |
   |------|--------|----------|---------------------|
   | `do-work/deliverables/{UR-NNN}-pipeline-summary.md`          | Plain markdown             | Developer reading in a terminal or editor — cat / grep / paste into a PR | `actions/pipeline-reference.md` → **Plain Markdown Report** (LLM-authored) |
   | `do-work/deliverables/{UR-NNN}-pipeline-summary.marp.md`     | Marp slide source          | Stakeholder walkthrough via `marp --preview`; also the source for the `.marp.html` export | `actions/pipeline-reference.md` → **Marp Slide Deck** (LLM-authored) |
   | `do-work/deliverables/{UR-NNN}-pipeline-summary.marp.html`   | Marp deck exported to HTML | Stakeholder who can't install `marp-cli` — share via URL or email | Produced mechanically by `npx @marp-team/marp-cli {UR-NNN}-pipeline-summary.marp.md --html` |
   | `do-work/deliverables/{UR-NNN}-pipeline-summary.single.html` | Standalone authored HTML   | Non-technical reader browsing independently — Mermaid + Tailwind via CDN, zero build, cross-links siblings | `actions/pipeline-reference.md` → **Standalone HTML Debrief** (LLM-authored) |

   The three authored files carry the same facts — no format-specific editorializing. The `.marp.html` export inherits its content mechanically from the `.marp.md` source; run the marp-cli command after writing the Marp source. Use `{REQ-id}-pipeline-summary.*` as the filename prefix if no UR was captured.

7. **Print the plain-markdown rendering to stdout** so the user sees the debrief immediately. Reference the other three paths (`.marp.md`, `.marp.html`, `.single.html`) in the closing "Deliverables" section so they can open them next.
8. **Queue continuation check**: Scan `do-work/queue/REQ-*.md` for files with `status: pending` in their frontmatter. Exclude any REQ IDs listed in the current pipeline's `artifacts` array (those should already be completed). If remaining pending REQs exist, proceed to Step 5a. If the queue is empty, suggest next steps and stop.

**Proportional depth:** Single-REQ Route A pipelines (config tweak, docs) get a minimal report in all three authored formats (plus the Marp HTML export) — status block + 1-row Final summary + How to verify. The Marp deck collapses to 3–4 slides; the `.single.html` collapses to a single-screen summary. Multi-REQ or Route B/C pipelines get the full treatment shown above. Don't pad short pipelines with empty sections, and don't truncate long ones. Match the report to the scope of the work.

### Step 5a: Queue Continuation

When the pipeline completes and additional pending REQs remain in the queue (from prior captures, follow-ups created during review, or other sources):

1. Print the continuation notice (`actions/pipeline-reference.md` → **Continuation Notice**) listing each pending REQ ID and its title
2. Record the list of pending REQ IDs about to be processed (e.g., `["REQ-043", "REQ-044"]`) — this is needed for review targeting in step 3
3. Dispatch actions/work.md (`do-work run`) in **standard queue-processing mode** — do NOT scope to pipeline artifacts. Pass the pending REQ IDs (e.g., `do-work run REQ-043 REQ-044`) so actions/work.md processes them.
4. After actions/work.md completes, dispatch actions/review-work.md for each REQ from step 2 individually (e.g., `do-work review REQ-043`, then `do-work review REQ-044`). Always use REQ IDs — never pass a UR ID, since UR-scoped review would re-review all completed REQs under that UR, not just this cycle's batch.
5. **Loop**: Scan `do-work/queue/REQ-*.md` for `status: pending` again. If more pending REQs remain (e.g., follow-ups created during the review step), repeat from step 1. If the queue is empty, print "Queue fully processed." and suggest next steps.

**Max iterations:** The continuation loop runs at most **3 cycles**. If pending REQs still remain after 3 run → review cycles, stop the loop and print the limit-reached message (`actions/pipeline-reference.md` → **Continuation Notice** — the limit-reached variant).

This prevents runaway loops when review steps keep generating follow-up REQs.

**Error handling:** If actions/work.md or actions/review-work.md fails during continuation:

1. Report the error to the user with context about what failed
2. Print how many REQs were successfully processed before the failure
3. Stop the continuation loop — do not retry automatically
4. Suggest the appropriate recovery command:
   - **Run step failed**: Suggest `do-work run` to resume processing pending REQs
   - **Review step failed**: Suggest `do-work review REQ-NNN` for each REQ that was processed but not yet reviewed (those REQs are already `status: completed`, so `do-work run` would be a no-op)

Unlike the main pipeline's error handling (Step 6), continuation errors do not update `pipeline.json` — the formal pipeline is already complete.

**State file note:** The pipeline's `active` field remains `false` during the continuation — the formal pipeline is complete. The continuation is a post-pipeline queue run. If the session ends mid-continuation, the user can resume with `do-work run` to process any remaining pending REQs.

### Step 6: Error Handling

If any step's dispatched action fails (error, exception, or the action reports failure):

1. Set the step's `status` to `"failed"` and add an `error` field with a brief description
2. Write the updated `pipeline.json`
3. Print the status block showing the failure point (`actions/pipeline-reference.md` → **Status Block**)
4. Report the error to the user with context about what failed and why
5. Leave `active: true` — the user can fix the issue and resume with `do-work pipeline`

On resume, the pipeline retries the failed step from scratch.

## Output Format

Every rendering this action prints or writes lives in `actions/pipeline-reference.md`, each already pointed to by name from the step that emits it:

- **Status Block** — printed after every step transition (Steps 2, 3, 4, 6).
- **Completion Status Block** — printed once, when all 6 steps are done (Step 5).
- **Pipeline Completion Report** — the three authored renderings (`.md`, `.marp.md`, `.single.html`) plus the mechanical `.marp.html` export, assembled from the dataset Step 5.4 gathers (Step 5).
- **Continuation Notice** — the initial listing, the "queue fully processed" message, and the limit-reached message (Step 5a).
- **Help Menu** — shown when Step 1 resolves to Help mode.

## Rules

- **Never skip steps.** The pipeline always runs in order: investigate → capture → verify → run → review → present. Steps cannot be reordered or omitted. (Exception: if `capture` produced no artifacts, `present` has nothing to deliver — mark it `done` with no artifacts and proceed.)
- **One pipeline at a time.** If an active pipeline exists, the user must complete, resume, or abandon it before starting a new one.
- **Orchestrator only.** The pipeline dispatches to existing actions. It never re-implements capture, work, verify, review, present, or inspect logic. Each action runs exactly as it would if the user invoked it directly.
- **Write state before dispatch.** Always update `pipeline.json` to `"in-progress"` before dispatching an action, and to `"done"` after it completes. This ensures the state file reflects reality even if the session ends unexpectedly.
- **The `run` step may be long.** actions/work.md processes only this pipeline's captured REQs but may still take significant time for complex requests. When starting this step, note: "Starting queue processing — this may take a while if multiple REQs are pending."
- **Optional stop-hook guard.** If your environment supports stop hooks, install `hooks/pipeline-guard.sh` to prevent accidental stops mid-pipeline — the pipeline works without it.
- **Do not commit the state file.** `do-work/pipeline.json` is transient session state. It tracks a single pipeline run and has no value after completion. It's excluded from git at initialize via the enclosing repo's `.git/info/exclude` (Step 4 of the New Pipeline flow), not the committable `.gitignore`.
- **Pass context to sub-agents explicitly.** Sub-agents have no conversation history. When dispatching a step via sub-agent, always include the pipeline request text and all artifact IDs from completed steps in the sub-agent prompt. Without this, sub-agents cannot target the correct UR/REQs.
- **Scope the `run` step to captured REQs only.** actions/work.md is queue-processing by default. When dispatched from the pipeline, it must only process the REQs created by this pipeline's capture step (listed in `artifacts`). Never process unrelated backlog items during a pipeline run.
- **Process remaining queue after completion.** After the pipeline's 6 steps finish, check for other pending REQs in the queue. If any exist, continue processing them automatically via run + review cycles until the queue is empty. This continuation uses standard queue-processing mode (not scoped to pipeline artifacts) and does not re-run `present` per cycle — the user can run `do-work present all` after the queue is fully processed if they want a portfolio summary. The pipeline state file remains `active: false` — the continuation is a post-pipeline operation. Maximum 3 continuation cycles — if REQs still remain after 3 cycles, stop and let the user continue manually.
- **Suggest next steps on completion.** After the pipeline finishes (including any queue continuation), suggest what the user might want to do next (see `next-steps.md`).
- **Completion is education, not a checkmark.** When all steps finish, produce the full Pipeline Completion Report (Final summary, Test state, Cross-REQ coherence, Carry-forward work, Deliverables, How to verify) in **all three authored formats** — plain markdown (`{UR-NNN}-pipeline-summary.md`), Marp slide source (`{UR-NNN}-pipeline-summary.marp.md`), and standalone HTML (`{UR-NNN}-pipeline-summary.single.html`) — then export the Marp deck to HTML (`{UR-NNN}-pipeline-summary.marp.html`) via `marp-cli`. One dataset, three renderings, one mechanical export, different audiences. A 12-REQ pipeline that prints only "Pipeline complete" — or writes only the markdown and skips the deck and the HTML — wastes the user's opportunity to understand and validate what shipped. Match report depth to pipeline scope — minimal for Route A, full for multi-REQ URs.
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
| "The markdown version is enough — who needs Marp and HTML?" | Write all three renderings — `.md`, `.marp.md`, and `.single.html` — from the same dataset, then export `.marp.html` via marp-cli | Different audiences consume the work on different surfaces: a developer scans the `.md` in a PR, a stakeholder sits through the deck (or views `.marp.html` if they lack marp-cli), a non-technical reader browses the `.single.html`. Shipping only the markdown leaves three audiences unserved. |
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
- Only one or two of the three authored rendering files (`.md`, `.marp.md`, `.single.html`) exist in `do-work/deliverables/`, or the `.marp.html` export is missing
- The three renderings disagree — e.g., the Marp deck lists 11 REQs but the markdown lists 12, or the HTML shows a test delta the markdown doesn't
- The `.single.html` file references external scripts beyond Tailwind + Mermaid, or requires `npm install` to render (the `.marp.html` export is exempt — its assets are whatever marp-cli bundles)
- The `.single.html` debrief caps the whole page at a fixed centered width (`max-w-6xl mx-auto` or similar) — a narrow reading column with empty gutters on a wide monitor; only running prose should cap width, while stat cards, tables, and card grids use the full width and wrap responsively (the Marp deck's fixed-size slides are exempt — full-bleed applies to the scrolling `.single.html` page only)
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
- [ ] Pipeline Completion Report rendered in all three authored formats (`{UR-NNN}-pipeline-summary.md`, `.marp.md`, `.single.html`) and exported to `{UR-NNN}-pipeline-summary.marp.html` via marp-cli — all four present in `do-work/deliverables/`
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
- [ ] The `.single.html` debrief is full-bleed — no page-level `max-w-*` centered cap (only running prose capped ≈74ch), with stat cards, tables, and card grids using the full width and wrapping responsively (Marp slides exempt)
