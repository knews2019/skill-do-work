# Work Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to process the queue. Processes pending requests from the `do-work/queue/` folder in your project. User-facing walkthrough: [`docs/work-guide.md`](../docs/work-guide.md).

An orchestrated build system that processes request files created by actions/capture.md. Uses complexity triage to route simple requests straight to implementation and complex ones through planning and exploration first.

## When to Use

**Use when:**
- The queue has `pending` REQs and the user wants them built (`do-work run`, `start`, `go`, etc.).
- Another orchestrator dispatches actions/work.md as its build step.
- A specific REQ id was named (`do-work run REQ-042`) — the action scopes to it.

**Do NOT use when:**
- The queue is empty — tell the user and stop; suggest `do-work capture-request: [describe]` instead.
- The only REQs left are `pending-answers` — report them and route answers to `do-work clarify`.
- See `SKILL.md` routing table for sibling action selection (inspect, verify-requests, review-work, etc.).

## Phase loading

Read this entry point in full. Follow the phases below in order, loading each linked procedure only when its condition applies. Read a companion once per context; returning to a phase reuses it. These companions are parts of this action, not additional commands. Named steps and anchors here remain routing entry points for existing callers.

`work-reference.md` remains authoritative for schemas, templates, recovery and finalization contracts: read the named section when a phase cites it, not the entire reference upfront. Follow its nested references as needed. Shared crew rules retain their own authoritative homes and conditional loaders. For standalone review, start at [review-work.md](./review-work.md); it does not enter the build pipeline.

The phase pointers replace the duplicate per-request checklist. They are required loads, not optional background reading. Error Handling below applies at every phase, including a phase that fails before its normal successor can load.

## Request Files as Living Logs

Each request file becomes a historical record. As you process a request, append sections documenting each phase: Triage, Plan, Exploration, Implementation Summary (mandatory file manifest), Testing, Review. This ensures full traceability — what was planned vs done, what files were touched, and whether triage was accurate. `<skill-root>/tools/do-work-cli.sh --format text req append-section REQ-NNN --section <name> --from <file>` writes one such section from a file holding only its body: it places the section in canonical order, refuses a second copy with a different body, and exits 0 with no change when the same body is already there; appending by hand in that order stays valid.

This living log is also the **trail of intent**. The REQ starts as a validated statement of what the user wants (written by capture). As actions/work.md processes it, each appended section documents how intent was interpreted and realized: builder decisions (## Decisions) record where the builder exercised judgment beyond stated intent, scope declarations (## Scope) record what the builder committed to, and implementation summaries record what was actually built. The gap between captured intent and realized implementation is visible in a single file.

## Architecture

The per-REQ orchestration pipeline (triage → estimate → plan/explore → implement → qualify → test → review → prepare → finalize) and the command that owns each deterministic mutation are mapped in `actions/work-reference.md` → **Architecture**. The orchestrator owns judgment and authored evidence; canonical commands own deterministic lifecycle, checkpoint, archive, release, and commit mutations.

> **Remember:** Every completed request gets a git commit (Step 9) before looping to the next request.

**Sub-agent note:** This document uses "spawn agent" language. Use your platform's subagent mechanism when available. If your tool doesn't support subagents, run phases sequentially in the same session and label outputs clearly.

**This action processes one REQ at a time unless you ask it not to.** That is the default and the floor. **The simplest agent that can read and write files and run shell commands must be able to follow this file end to end**, which is why concurrency is opt-in rather than resident: machinery in the main path would sit in front of every reader who cannot use it.

**Concurrency is opt-in; isolation is not.** Every run mode, the serial default included, implements each REQ on its own per-REQ branch or worktree and hands that branch back for the orchestrator to merge, so a REQ set aside part-way through implementation cannot reach the next REQ's diff, qualification, tests, staging, or commit. `actions/fan-out-reference.md` → **Worktree Dispatch Mode (Step 1)** is that contract for every mode, *Isolation ladder* included. **Builder and orchestrator are roles, not processes:** where the harness has no subagents the same agent plays both, and the boundary between them is the branch.

**`do-work run --fan-out [N]` opts in.** The loop then dispatches builders, **without a confirmation gate**, for the set queue-mode `advance` computes and claims; `advance` owns the readiness predicate, so do not re-derive one (`actions/fan-out-reference.md` → **Auto-wave**). Per-REQ steps run unchanged, **integration stays serial** (merge → qualify → test → review → changelog → archive, one REQ at a time), and **the dispatch mechanism stays unspecified**: the owner synthesizes from files on disk, never conversation. Without worktree support, or on a harness that cannot run an agent against a directory you choose, the flag **degrades silently to the serial loop**. `--coordinate` adds who integrates: the main session stays the coordinator and hands each REQ's integration to an agent (`actions/fan-out-reference.md` → **Delegated integration — the coordinator shape**).

## Selection and Triage

Read [selection](./work-selection-reference.md) before queue mutation: input validation, schema, recover/claim, triage and Open Questions (Steps 1–3.5).

### Mid-Run Messages (any step)

A user message can arrive at any step and never stops the run: route it, keep going, and say in one progress line where it went. Carry the user's own words, never a paraphrase, to every REQ it affects. What it changes decides where it goes: a progress question is answered from disk; a change to an in-flight REQ becomes a `## Addendum (mid-run)` section before its builder hands back, or a new `addendum_to` REQ after; a queued REQ or new work goes through capture; a change to what the hand-back shows first becomes the run's hand-back emphasis note. Each branch and its writer: `actions/fan-out-reference.md` → **Mid-Run Messages (any step)**.

### Mechanical Evidence-Gate Loop

Whenever per-request `advance` reports a mechanical phase, invoke it again with the judgment-owned inputs named by its typed finding and consume only `advance.gate_records` whose `request_id` and `request_path` match the active REQ. `needs_input` means satisfy the record's tokenized input and retry; `findings` means apply the remaining judgment principle and retry; `failed` stops that evidence boundary; `satisfied` authorizes writing the durable REQ evidence that moves the classifier forward. Then call `advance` again. The command composes estimation, pre-flight, qualification/scope comparison, focused-test baseline comparison, and green-record checks through their existing handlers; do not call those handlers separately or reconstruct their result from display text.

Inputs still require judgment. Extract nontrivial estimator signals with `actions/estimate-reference.md`; choose project test and canonical-gate argv from the prime or native project configuration; judge qualification warnings, semantic failure similarity, TDD validity, retry/failure classification, canonical-gate attribution/deferral, and heavy-lane scheduling. Persist an executed estimate record's `output_lines` as the frozen `estimate:` block with a fresh `calculated_at`, then print its P50, confidence, and basis before planning. A reusable estimate remains byte-identical. Estimation is informational and never gates execution.

### Step 3.7: Spec Loading (optional)

After triage, check if a specification template matches this REQ's domain or task type.

1. **Match by task type:** If the REQ's title or What section clearly indicates a task type (API endpoint, UI component, refactor, bug fix), check `specs/` for a matching template (`specs/api-endpoint.md`, `specs/ui-component.md`, `specs/refactor.md`, `specs/bug-fix.md`).
2. **Match by suggested spec:** If the REQ's frontmatter contains a `suggested_spec` field (set during capture), check `specs/` for that template.
3. **If a matching spec exists**, read it and use it to inform:
   - The implementation checklist order (pass to the planning or implementation agent)
   - Quality standards to verify against (pass to the review step)
   - Common pitfalls to watch for (include in the builder's context)
4. **The spec is guidance, not override** — the REQ's specific requirements always take priority. If the REQ's requirements conflict with a spec's recommendations, follow the REQ.
5. **If no matching spec exists**, proceed normally. Specs are optional — their absence never blocks work.

### Step 4: Planning (Route C only)

Route C: read [planning](./work-planning-reference.md). Routes A/B: append the **Plan Skip Note — Routes A/B (Step 4)** from `actions/work-reference.md` if absent; omit `planning_at`.

### Step 5: Required-Lessons Consult (all routes) and Exploration (Routes B and C)

**Required-lessons consult — every REQ.** Before exploration, read the REQ's existing `required_lessons` entries in order, then consult `do-work/lessons-index.md` even when the field is absent. Existing entries are mandatory reads and consume the budget first. Match additional rows against the request text, Plan/current `write_set`, likely touched paths, and `prime_files`, using the ranking, entry forms, full-coverage targeting rule, mechanical costs, and single limit owned by `actions/capture-reference.md` → **Required Lessons Budget Contract**. Prefer narrowing an eligible bare match before dropping it. Refresh `required_lessons` with the fitting existing-plus-new set; never add an entry without an index match. Preserve and read captured entries when the index itself is absent, but add nothing.

Read the resolved set now: a bare path means the whole satellite; `path#family-slug` means only bullet lines carrying `[family: family-slug]`. A missing listed file never blocks exploration — record the missing entry in `## Exploration` when that section exists and always in the builder hand-back, then continue. Replace any existing `## Required Lessons — Dropped for Budget` body section with the current claim-time drops (entry, cost, matching reason), or remove that section when nothing was dropped; a drop is never silent. Omit `required_lessons` when the refreshed set is empty. This consult closes the capture-time gap for old REQs and for later members of a serial batch whose matching lesson did not exist when the batch was captured.

**Route A:** stop this step after the consult and resolved reads; only exploration is skipped. **Routes B and C:** continue below.

Read [exploration and pre-build evidence](./work-exploration-reference.md) for Routes B/C only; Route A continues to Step 6.

### Step 5.5: Scope Declaration (Routes B and C)

Continue the [exploration and pre-build evidence procedure](./work-exploration-reference.md); Route A skips this step.

### Step 6: Implementation

Read [the implementation procedure](./work-implementation-reference.md).

### Step 6.25: Implementation Summary

Read [the implementation procedure](./work-implementation-reference.md).

### Qualification and Testing Judgment

Read [the testing procedure](./work-testing-reference.md).

### Step 7: Review

Read [the review procedure](./work-review-reference.md).

### Step 7.5: Lessons-Capture Phase

Read [the review procedure](./work-review-reference.md).

### Step 7.7: Heavy-Test Hold and Drain

Read [the heavy-testing procedure](./work-heavy-testing-reference.md). Load only when heavy lanes were selected or recovery reports a held REQ; drain at queue exhaustion.

### Step 8: Prepare Finalization

Read [the finalization procedure](./work-finalization-reference.md).

### Step 9: Commit Phase (Git repos only)

Read [the finalization procedure](./work-finalization-reference.md).

### Step 10: Loop or Exit

Read [the finalization procedure](./work-finalization-reference.md).

## Clarify Questions

The clarify workflow has its own action. Run `do-work clarify` — it handles batch-review of `pending-answers` REQs, where the user confirms, overrides, or discards builder decisions. Resolved REQs flip back to `pending` and re-enter the work queue. Questions routed to an outside stakeholder live on `blocked` REQs carrying `stakeholder:` instead — their answers come back through `do-work stakeholder-answers`, and clarify only routes them.

## Error Handling

| Phase | Action |
|-------|--------|
| `pending-answers` REQs remain after queue is empty | Report them to the user: list each REQ and its unresolved questions. Suggest `do-work clarify` to batch-review. |
| Held claimed REQs remain after the queue is empty | Run the heavy-lane drain in Step 7.7. A REQ still held afterwards (skipped lane, plan drift, dirty tree) is listed in the composed exit summary with its typed finding; not a failed run. |
| `blocked` REQs remain after queue is empty | Report each with its `blocked_by` condition and age (per the composed exit summary). Suggest re-running `do-work run` (auto-probes any `blocked_check`) once a condition is met, or `do-work clarify` to confirm a human-checkable one. For a stakeholder-questions REQ (`stakeholder:` present), suggest sharing its report and then `do-work stakeholder-answers` — never a probe or a yes/no confirm. |
| Plan agent fails (Route C) | Classify failure (Intent/Spec/Code/Environment), create follow-up REQ if applicable, archive as failed |
| Explore agent fails (B/C) | Proceed to implementation with reduced context — builder can explore on its own |
| Builder reports a missing external precondition | Apply the blocked-flip test (`actions/work-reference.md` → **Failure Classification (Step 8)**, Environment row): if no substantive implementation edits landed AND the condition is expected to self-resolve, flip to `status: blocked` (non-terminal, back to the queue) instead of failing. Otherwise classify as Environment and archive as failed. |
| Orchestrator wants to park a queued REQ on the operator (production access, credentials, a person at the keyboard) | Flip it to `status: blocked` with `blocked_by` naming the operator and `blocked_at`, the same mid-run flip as a missing precondition (`actions/work-reference.md` → **Failure Classification (Step 8)**, Environment row). Never write `assigned_to` for this: the earmark is for another session, and the default run skips it, hiding the wait from the user. When the remaining work is a routine operator act (deploy, publish and the like), the REQ instead completes with what the AI built and the hand-back names the operator step: no `blocked` flip, no failure, no abandon recommendation. |
| Implementation fails | Classify failure (Intent/Spec/Code/Environment), create follow-up REQ if applicable, archive as failed |
| Focused tests or an attributed current-diff canonical repository gate fail repeatedly | After 3 fix attempts, classify as Code failure, create a follow-up REQ with the focused-test or attributed current-diff gate failure details, and archive as failed. A matching unrelated gate failure instead uses the repository-gate deferral lifecycle; a mismatch, launch failure, invalid range, or unverifiable attribution remains a fail-safe stop and is never archived as success. |
| Ordinary review: Acceptance = Fail | Return to Step 6 for ONE remediation attempt, then re-review. If still failing: archive as `completed-with-issues` with follow-up REQs |
| Ordinary review work agent fails | Skip review, note it in the REQ file, continue to archive — review failure is not a gate for an unselected run |
| Selected AI-Evaluation-Loop cannot approve (including failed/missing required review) | Persist unresolved / not approved and the stop reason per `actions/ai-evaluation-loop-reference.md`; Step 8 selects `completed-with-issues`. Never use the ordinary skip-review or score fallback. |
| Finalization or commit fails | Consume the typed blocker, correct the judged manifest or underlying hook failure, then rerun the exact `advance` continuation. Never bypass the hook, hand-stage the tail, or infer an archive from partial state. |
| Unrecoverable error | Stop loop, report clearly, leave queue intact for manual recovery |

## Progress Reporting

Keep the user informed with this format:

(keep the user informed in the running per-REQ progress format shown in `actions/work-reference.md` → **Progress Reporting Example**)

When the run finishes or pauses, hand back with the **Decision Brief** (`actions/work-reference.md` → **Decision Brief (hand-back format)**): lead with WHAT'S BEING BUILT (each REQ's `## Orientation`, at subsystem altitude), then DECISIONS FOR YOU (any escalated `pending-answers` follow-ups, each with value + risk), then WAITING ON OTHERS (stakeholder-routed questions, with each report's path to share), then a collapsed HANDLED list, unless the run manifest's hand-back emphasis note orders them otherwise (`actions/fan-out-reference.md` → **Mid-Run Messages (any step)**). For selected AI-Evaluation-Loop work, show approved or unresolved / not approved, current candidate, round/stall counts and remaining concerns. Never lead with review scores — they stay in the per-REQ progress lines, not in front of the hand-back.


## Archived Request File Example

See [sample-archived-req.md](./sample-archived-req.md) for a complete example of what an archived REQ looks like after processing through the full pipeline (Route B). Every section shown there is generated by the steps above.

**Timestamps tell the story:** `created_at` → `claimed_at` = queue wait time. `claimed_at` → `completed_at` = calibration wall span, not active implementation time. Optional phase stamps show the observed pipeline breakdown, and `completed_at` → `release_at` is the release tail. Route + timestamps let you calibrate triage accuracy over time.

## Rules

- Finalization commits only its exact manifest allowlist and never bypasses a commit hook. `supplied_commit` provenance records the supplied merge hash without recommitting implementation paths; a `primary_commit` finalization, which has no implementation commit to record, may require a second metadata commit (`actions/work-reference.md` → **Commit & Metadata-Commit Procedure (Step 9)**).
- `write_set` is display-only at any builder count: **advisory input to a human's fan-out pick, never a gate**, and **not read by auto-wave**. The merge is the non-interference proof (`actions/fan-out-reference.md` → **Fan-Out Dispatch**). It feeds the board's overlaps badge (`actions/work-reference.md` Request File Schema).
- In worktree dispatch mode the orchestrator is the sole integrator: the builder never writes the main tree, and a merged state is not a verified state until the REQ's checks re-run there.

**Common mistakes to avoid:**

- Spawning implementation agent without first moving file to `working/`
- Letting spawned agents handle file management (only the orchestrator moves/archives files)
- Reproducing a queue status transition by hand instead of consuming `advance` or finalization
- Archiving a UR folder before all its REQs are complete
- Forgetting Planning status note for Routes A/B ("Planning not required")
- Committing without validating Implementation Summary file list against staged files
- Implementation Summary that only lists `do-work/` paths (means the REQ wasn't actually implemented — exception: `domain: ui-design` design artifacts placed in project directories like `<project-root>/docs/design/`)
- Creating follow-ups for every `- [~]` item instead of only UX-affecting decisions

**This action does NOT:**

- Create new request files (use actions/capture.md)
- Make architectural decisions beyond what's in the request
- Run without user present (this is supervised automation)
- Modify already-completed requests
- Allow external modification of files in `working/` or `archive/`

## Common Rationalizations

Shared conditions live in [shared-principles.md](../crew-members/shared-principles.md); these rows govern this action’s lifecycle.

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "`recover` refused, so the run has to stop here" | Judge each blocked path and clear it, then re-run the exact `verification_argv` (`actions/work-reference.md` → **Stuck Runs Hand Off to Judgment (any step)**) | The command refuses anything it cannot attribute and has no opinion about what the bytes are; a Finder `.DS_Store` under `do-work/` once parked an entire queue behind that refusal |
| "I'll skip pre-build evidence — the baseline is probably stable" | Run the pre-build evidence gate and inspect its record | Pre-existing failures get misattributed to the builder, and pre-existing dirty files can contaminate the repository-wide diff used for qualification and review |
| "The queue file's twin is already archived, but re-running is harmless" | Consume `advance`'s `blocked-archive-collision` hold | Re-processing a duplicate silently re-commits it and corrupts the archive lineage |
| "The merge went fine, so I'll just `-D` the leftover `worktree-agent-*` branch" | Run `git branch -d` from the integration branch and stop on refusal (Step 9; `do-work-cli worktree cleanup REQ-NNN` runs it) | `-d`'s refusal is the only assertion the integration actually happened; `-D` deletes the evidence that a merge was skipped or lost, along with the work |
| "Cleanup just needs the branch name — I'll rebuild it from the REQ slug" | Merge and clean up by the operative name held since dispatch (Step 6 hand-back, Step 9 cleanup) | After a collision variant, the slug-derived name is the *leftover*: the merge would integrate the wrong branch, `-d` refuses on unmerged work, and the run halts on a false lost-merge alarm while the variant is never cleaned |

## Red Flags

- REQ in `do-work/working/` for >1 hour with no new git commits (builder may be stuck)
- Implementation Summary lists files but `git diff` shows no changes in those files (hollow implementation)
- No Triage section appended to the REQ after processing begins
- Scope section declares 3 files but Implementation Summary lists 12 (scope creep)
- Builder created files only inside `do-work/` and no source files changed (no real work done)

## Verification Checklist

- [ ] All pending REQs processed or explicitly skipped with documented reason
- [ ] Every completed REQ has an Implementation Summary section with file manifest
- [ ] No REQ files remain in `do-work/working/` after the work loop ends — except a claim Step 1 deliberately left intact (Step 1 Crash Recovery)
- [ ] CHECKPOINT.md written if ending mid-session (for resume)
- [ ] Git commit created for each completed REQ
- [ ] Cleanup pass triggered at end of work loop
