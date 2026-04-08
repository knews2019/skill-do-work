# Pipeline Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to run a full end-to-end pipeline: investigate, capture, verify, run, review. Manages state across sessions via `do-work/pipeline.json`.

A stateful multi-action orchestration that chains five actions in sequence. Each step dispatches to an existing action. The pipeline tracks progress in a JSON state file, supports resume across sessions, and reports status at each transition.

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
    { "name": "review",      "status": "pending", "completed_at": null }
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
| `steps[].name` | string | Step identifier: `investigate`, `capture`, `verify`, `run`, `review` |
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
2. Generate `session_id`: today's date + `-001` (increment if prior pipelines ran today — check existing `pipeline.json` for same-day IDs)
3. Write `do-work/pipeline.json` with:
   - `request` set to `$ARGUMENTS` (the request text, stripped of the "pipeline" or "full" keyword)
   - All 5 steps set to `status: "pending"`
   - `active: true`
   - `started_at` set to current ISO 8601 timestamp
4. Print the initial status block
5. Proceed to Step 4 (execute first step: `investigate`)

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

| Pipeline step | Action to dispatch | What to pass |
|---------------|-------------------|--------------|
| `investigate` | the inspect action (`do work inspect`) | No arguments — inspects all uncommitted changes |
| `capture` | the capture action (`do work capture request: {request}`) | The `request` field from pipeline.json |
| `verify` | the verify requests action (`do work verify requests`) | No arguments — targets most recent UR |
| `run` | the work action (`do work run`) | No arguments — processes the queue |
| `review` | the review work action (`do work review work`) | No arguments — reviews most recent completed work |

Dispatch each action the same way the main router dispatches actions: subagent if available, inline otherwise. The pipeline action is the orchestrator — it calls the router's dispatch mechanism, not the action files directly.

3. **After the action completes**:
   - Update the step's `status` to `"done"` and set `completed_at` to current timestamp
   - **Capture step only**: Parse the action's output for created REQ and UR IDs (e.g., `REQ-042`, `UR-018`). Store them in the step's `artifacts` array. If IDs cannot be parsed, leave `artifacts` as an empty array — do not block the pipeline.
   - Write the updated `pipeline.json` immediately
   - Print the updated status block

4. **Advance**: If more steps have `status: "pending"`, proceed to the next one (loop back to the top of Step 4). If all steps are `"done"`, proceed to Step 5.

### Step 5: Completion

When all 5 steps are done:

1. Set `active: false` in `pipeline.json`
2. Write the final `pipeline.json`
3. Print the completion status block (all checkmarks)
4. Print a completion summary:

```
Pipeline complete.
  Session:    {session_id}
  Request:    {request}
  Duration:   {elapsed time from started_at to now}
  Artifacts:  {list of REQ/UR IDs from capture step}
```

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
```

## Rules

- **Never skip steps.** The pipeline always runs in order: investigate → capture → verify → run → review. Steps cannot be reordered or omitted.
- **One pipeline at a time.** If an active pipeline exists, the user must complete, resume, or abandon it before starting a new one.
- **Orchestrator only.** The pipeline dispatches to existing actions. It never re-implements capture, work, verify, review, or inspect logic. Each action runs exactly as it would if the user invoked it directly.
- **Write state before dispatch.** Always update `pipeline.json` to `"in-progress"` before dispatching an action, and to `"done"` after it completes. This ensures the state file reflects reality even if the session ends unexpectedly.
- **The `run` step may be long.** The work action processes the entire queue and can take significant time. When starting this step, note: "Starting queue processing — this may take a while if multiple REQs are pending."
- **Platform-agnostic.** No tool-specific APIs. Dispatch actions the same way the main router does. If your environment supports stop hooks, you can optionally install `hooks/pipeline-guard.sh` to prevent accidental stops mid-pipeline — but the pipeline works without it.
- **Suggest next steps on completion.** After the pipeline finishes, suggest what the user might want to do next (see the next-steps reference).
