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

## Request Files as Living Logs

Each request file becomes a historical record. As you process a request, append sections documenting each phase: Triage, Plan, Exploration, Implementation Summary (mandatory file manifest), Testing, Review. This ensures full traceability — what was planned vs done, what files were touched, and whether triage was accurate.

This living log is also the **trail of intent**. The REQ starts as a validated statement of what the user wants (written by capture). As actions/work.md processes it, each appended section documents how intent was interpreted and realized: builder decisions (## Decisions) record where the builder exercised judgment beyond stated intent, scope declarations (## Scope) record what the builder committed to, and implementation summaries record what was actually built. The gap between captured intent and realized implementation is visible in a single file.

## Architecture

The per-REQ orchestration pipeline (triage → estimate → plan/explore → implement → qualify → test → review → prepare → finalize) and the command that owns each deterministic mutation are mapped in `actions/work-reference.md` → **Architecture**. The orchestrator owns judgment and authored evidence; canonical commands own deterministic lifecycle, checkpoint, archive, release, and commit mutations.

> **Remember:** Every completed request gets a git commit (Step 9) before looping to the next request.

**Sub-agent note:** This document uses "spawn agent" language. Use your platform's subagent mechanism when available. If your tool doesn't support subagents, run phases sequentially in the same session and label outputs clearly.

**This action processes one REQ at a time unless you ask it not to.** That is the default and the floor: Step 1 advances one claim, Step 6 waits for its builder before Step 6.25 reads the output, and the loop advances again after finalization. **The simplest agent that can read and write files and run shell commands must be able to follow this file end to end**, which is why concurrency is opt-in rather than resident — machinery in the main path would sit in front of every reader who cannot use it.

**Concurrency is opt-in; isolation is not.** They are separate axes and only the first is a flag. Every run mode — the serial default included — implements each REQ on its own per-REQ branch or worktree and hands that branch back for the orchestrator to merge, so a REQ set aside part-way through implementation cannot reach the next REQ's diff, qualification, tests, staging, or commit. `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** is that contract for every mode, and its *Isolation ladder* defines what isolation means on a harness with no `git worktree` support and names the single rung where none is possible. **Builder and orchestrator are roles, not processes:** where the harness has no subagents the same agent plays both, and the boundary between them is the branch — never a second process.

**`do-work run --fan-out [N]` opts in.** In that mode the loop **computes the ready set itself and dispatches builders without a confirmation gate** — pending, dependency-ready, unclaimed, not `assigned_to` another session, and not dropped by a selection filter you passed (`--skip-impact-negligible` is today's), bounded to N (or the harness limit, or two). **That list is a gloss of the conditions as they stand, never the definition:** queue-mode `advance` owns the predicate and returns every exclusion typed, and `actions/work-reference.md` → **Worktree Dispatch Mode** → *Fan-Out Dispatch* → **Auto-wave** states the policy around it rather than a second predicate. Do not re-derive one here either. Three things the flag does **not** change: every per-REQ step below runs unchanged per REQ (one worktree, one hand-back merge, one merge range, one cleanup); **integration stays serial** — merge → qualify → test → review → changelog → archive, one REQ at a time, which is where the wall-clock saving stops; and **the dispatch mechanism stays unspecified**, so a spawned subagent and a separately-opened session are indistinguishable to the owner, which must synthesize from files on disk rather than conversation. Without worktree support, or on a harness that cannot run an agent against a directory you choose, the flag **degrades silently to the serial loop** — not an error.

**`write_set` is not an input to the wave.** It is display-only at any builder count, and the merge — not the pick — is what proves two builders did not collide (`actions/work-reference.md` → Fan-Out Dispatch). A computed set asserts that its REQs are all *runnable*, never that they do not overlap.

## Complexity Triage

Before spawning any agents, assess the request to determine the right route.

### Route A: Direct to Builder (Simple)

Skip planning and exploration entirely.

**Indicators:** Bug fix with clear steps, value/config change, single UI element add/remove, styling tweak, request names specific files, well-specified with obvious scope (<~50 words), copy changes, feature flag toggle.

**Examples:** "Change button color from blue to green", "Fix crash when clicking submit with empty form", "Update API timeout from 30s to 60s"

### Route B: Explore then Build (Medium)

Skip planning, run exploration. The "what" is clear, the "where" or existing patterns need discovery.

**Indicators:** Clear outcome but unspecified location, "like the existing X", need to find similar implementations, modifying something at an unknown location.

**Examples:** "Add form validation like we have on the login page", "Create a new API endpoint following our existing patterns"

### Route C: Full Pipeline (Complex)

Plan, then explore, then implement.

**Indicators:** New feature requiring multiple components, architectural changes, ambiguous scope ("improve", "refactor"), touches multiple systems, external service integration, long request (100+ words) with many requirements, user explicitly asked for a plan.

**Examples:** "Add user authentication with OAuth", "Implement dark mode across the app", "Refactor state management to use Zustand"

### Decision Flow

```
Read the request
  ├── Names specific files AND has clear changes? → Route A
  ├── Bug fix with clear reproduction? → Route A
  ├── Simple value/config/copy change? → Route A
  ├── Clear outcome but location/pattern unknown? → Route B
  ├── Ambiguous, multi-system, or architectural? → Route C
  └── Default: Route B (builder can request planning if needed)
```

**When uncertain, prefer Route B.** Under-planning is recoverable; over-planning is wasted time.

## Folder Structure

The `do-work/` folder layout is described in `actions/work-reference.md` → **Folder Structure**. Briefly: `queue/` holds pending REQs, `working/` holds the claimed REQ, `archive/` holds completed work (UR folders + legacy REQs), and `user-requests/` holds active UR folders until all their REQs finish.

## Request File Schema

The full annotated frontmatter schema and the **Schema Read Contract** — the normalize-and-warn rules every read site honors for fields with a documented canonical vocabulary — live in `actions/work-reference.md` → **Request File Schema — Full Frontmatter** and **Schema Read Contract**. Every reference below to "the Schema Read Contract" points there.

**Status flow (frontmatter values):** `pending` → `claimed` → `completed` / `completed-with-issues` / `failed`. The heavy-test hold is a phase of `claimed`, marked by the REQ's `## Heavy Verification Plan` section and its `commit:`, per the section-tracking rule in the next paragraph — it is not a status.

The intermediate phases (planning, exploring, implementing, testing, reviewing) are tracked by which `##` sections exist in the REQ file, not by frontmatter status changes. Only two status transitions are written on the ordinary path: `advance` commits `pending` → `claimed`, then finalization commits the terminal status. Exception paths write their documented hold or remediation status through their canonical command.

**Special statuses — these REQs stay in the queue but are not claimable while the condition remains:**
- `pending-answers` — a follow-up REQ whose Open Questions need user input before it can be worked. These accumulate in the queue and get batch-reviewed when the user runs `do-work clarify`.
- `blocked` — waiting on an **external condition** named in `blocked_by`; `advance` runs any recorded probe and atomically unblocks a success before claim. Human-confirmable conditions may also resolve through `do-work clarify`.
- `blocked-archive-collision` — set by `advance` when recursive discovery finds the queue REQ id in the archive. Non-destructive holding state; the user flips it back to `pending` (or removes/renames the duplicate) after deciding what to do.
- `blocked-dependency-cycle` — set by Step 1 when a REQ's `depends_on` graph contains a cycle (e.g., REQ-A depends on REQ-B which depends on REQ-A). Non-destructive holding state; the user edits the `depends_on` chain to break the cycle, then flips the status back to `pending`.

## Input

`$ARGUMENTS` may contain:

- **Targeting tokens — specific `REQ-NNN` or `UR-NNN` ids** (e.g., `REQ-042`, `REQ-042 REQ-043`, `UR-011`, or a mix) — process only the resolved REQs and stop (do not process the full queue). This is how a caller scopes work to a specific batch. Token shapes and UR→REQ expansion follow the **Target ID Resolution** contract in `actions/work-reference.md`. **Provenance decides dependency gating:** an **explicitly-named `REQ-NNN`** bypasses `depends_on` (the user named it directly); a REQ reached by **`UR-NNN` expansion** does **not** — it goes through the normal dependency-ready filter, scoped to the UR's member set (naming a batch is a weaker signal than naming each member, and capture wrote those edges expecting them honored). A mixed `do-work run REQ-042 UR-011` is the deduped union, each member keeping its own provenance.
- **`--fan-out [N]`** (optional integer) — enter **auto-wave mode**: compute the ready set and dispatch builders concurrently, bounded to N. Bare `--fan-out` uses the harness concurrency limit, or **two** where that is unknown. It changes *how many* of the selected set run at once, never *which* — so it composes with everything that selects a set, `--wave N` and targeting tokens included. Requires worktree support and a harness that can run an agent against a chosen working directory; without either it degrades silently to the serial loop. The predicate belongs to queue-mode `advance`; the policy around it is `actions/work-reference.md` → **Worktree Dispatch Mode** → *Fan-Out Dispatch* → **Auto-wave**.
- **`--wave N`** (integer flag, default mode only) — run only REQs at dependency depth N. Roots (no `depends_on`, or all `depends_on` resolve to archived REQs) are depth 0; depth grows by one per dependency layer. Mutually exclusive with **any** targeting token (`REQ-` or `UR-`) — reject the combination with an error.
- **`--skip-impact-negligible`** (boolean flag) — omit every REQ whose `impact:` resolves to `impact-negligible` from Step 1's selection, and report the ones it dropped. It changes *which* REQs are selected, never how many run at once, so it **composes with `--fan-out` and stacks with `--wave`**: `--wave` chooses the depth, this flag subtracts the negligible REQs from it, `--fan-out` chooses how many of the remainder run concurrently. **An explicitly-named `REQ-NNN` overrides it** — the user named the REQ outright, exactly as explicit naming overrides `depends_on` and `assigned_to` — while a REQ reached by `UR-NNN` expansion does **not** override it, per the same per-token provenance rule. Deliberately one boolean and not a general `--impact <token>` selector: stopping negligible work is the only filter anyone has asked for, and a second use would be the time to generalize.

**Unrecognized arguments are rejected, not ignored.** After stripping `--wave N`, `--fan-out [N]`, and `--skip-impact-negligible` and extracting targeting tokens (`REQ-`/`UR-` followed by digits, case-insensitive, per the Target ID Resolution contract), any non-empty token still left in `$ARGUMENTS` is an error. Stop and report:

```
Unrecognized argument(s): <tokens>. Usage: do-work run [REQ-NNN|UR-NNN ...] [--skip-impact-negligible] [--fan-out [N]] | do-work run --wave N [--skip-impact-negligible] [--fan-out [N]] | do-work run
```

Do **not** fall through to full-queue processing. A leftover token almost always means the user meant to *scope* the run — a typo'd REQ ID (`REG-042` instead of `REQ-042`), or dead muscle memory (a retired mode word) — so silently building the entire queue is the wrong, hard-to-undo default. This generalizes the existing `--wave`-plus-targeting-token rejection to all unrecognized residue; both are parse-time guards.

When `$ARGUMENTS` is empty — no targeting tokens, no flags, no other tokens — process all pending REQs in dependency-aware order (default behavior).

## Steps

**actions/work.md is an orchestrator.** You own semantic judgment and action-authored evidence. Canonical commands own deterministic file management; spawned agents handle implementation work only.

### Step 1: Select and Claim

At each queue boundary, run canonical `recover`, then queue-mode `advance` with the user's targeting and selection arguments, and treat its typed result as the sole authority: a REQ is claimable when its normalized queue state, targeting provenance, dependency state, assignment, impact filter, external probe, and collision evidence permit it. **Recovery is read per REQ, never as one pass/fail gate for the run:** a finalization record whose `reason_codes` carry `FINALIZATION-SET-ASIDE` excludes that one REQ from this run's selection, its other reason codes say what refused, and the run continues with what remains, while a typed refusal — dirt no REQ owns — stops the run and hands off to judgment for the resolving verb (`actions/work-reference.md` → **Commit & Metadata-Commit Procedure (Step 9)** for the per-record contract, → **Stuck Runs Hand Off to Judgment (any step)** for the verb). The command owns selection, holds, unblock, claim, frozen membership, dispatch bounds, and tokenized continuation; the action dispatches only committed `queue_advance.claimed` members, follows the returned continuation unchanged after integration, and otherwise renders **Composed Exit Summary (Step 1)** from the same result without reproducing any mechanic.
### Step 3: Triage

Read the request, apply the decision flow, update frontmatter with `route`. If a `## Triage` section does not already exist, append to the request file:

(append per the **Triage Section Template (Step 3)** in `actions/work-reference.md`)

Report the triage decision briefly to the user.

**Addendum REQs:** If the REQ has `addendum_to` in frontmatter, read the original REQ before building. If the original includes a `## Prior Implementation` section, use it. If it doesn't (e.g., the original was in-flight when the addendum was captured but has since completed), find the original in `do-work/archive/` and read it to understand what was already built — key files, patterns, and approach. This prevents duplicating or conflicting with existing work.

### Step 3.5: Open Questions — Best Judgment, Not a Gate

After triage, scan the REQ for a `## Open Questions` section with `- [ ]` items. Open Questions are **not a blocker** — the builder proceeds with its best judgment and completes the REQ.

Open Questions use checkbox syntax:
- `- [ ]` — **Unresolved**: has `Recommended:` and `Also:` choices from capture
- `- [x]` — **Resolved**: user answered (answer follows `→`)
- `- [~]` — **Deferred**: builder used its best judgment (reasoning follows `→`)

**If unresolved `- [ ]` items exist:**

1. Note them. Read the `Recommended:` default and `Also:` alternatives for each.
2. Mark each as `- [~]` with a numbered decision and the builder's reasoning: `- [~] [question] → **D-01**: Builder chose: [choice]. Reasoning: [why]`. An Open-Questions item is a deferred ambiguity, so it is almost always an **ESCALATE** decision under the decide-vs-escalate gate (`crew-members/coding-guardrails.md` § Think Before Coding) — append its **value** and **risk** so they carry into the follow-up the user reviews — or into the stakeholder REQ, when the record names an outside `Answerer:` (Step 8's audience fork): `... Reasoning: [why]. Value: [what this choice buys]. Risk: [what breaks if it's wrong, and how reversible]`. When the REQ's `prime_files` cover this area, source the value/risk from the prime's `## Stakes` section rather than re-deriving it.

   Two optional clauses extend the record, each keyed on its condition: append `Answerer: <name>` when the person who can really answer is a named **outside stakeholder** — someone other than the session user, named in the question's own capture-time `Answerer:` line or in the REQ's text, verbatim and never invented — and append `Irreversible: yes — [why undoing is expensive]` when the assumption would be expensive or impossible to undo. **Neither clause blocks anything.** The builder proceeds to completion exactly as before; an irreversible assumption is flagged prominently (Step 8's stakeholder routing and the composed exit summary), never gated on. The clauses only route the confirm-or-override at Step 8: no `Answerer:` → the user's `pending-answers` follow-up; `Answerer:` present → that person's stakeholder REQ (`do-work stakeholder-answers` ingests their reply later).
3. Number decisions sequentially per REQ (D-01, D-02, D-03...). Open Questions decisions and Implementation Decisions (Step 6) share the same D-XX ID space — if Open Questions uses D-01 through D-03, the first implementation decision is D-04. After resolving all `- [ ]` items, append a counter comment immediately after the `## Open Questions` section so Step 6 knows the next available ID: `<!-- D-XX counter: last used D-03. Next decision: D-04. -->` If no decisions were made in this step, write `<!-- D-XX counter: none used. Next decision: D-01. -->` These IDs can be referenced by future REQs.
4. Proceed with implementation using those decisions.

**If a question was escalated to the user mid-run and answered:** a long-running orchestration does legitimately surface a decision between units of work, and that answer lives only in the asking session's context until something writes it down (`crew-members/clear-questions.md` Principle 8). The **orchestrator** (never the builder — all file management is the orchestrator's) writes it into the REQ *before dispatch*, in the **Canonical answered-question format** (`actions/clarify.md`): the item becomes `- [x] [question] → [the user's answer]`, **not** `- [~]`, and carries no D-XX — those record *builder* choices, and recording a user answer as one would leave the stored `Recommended:` rationale standing as the reason and send the settled decision back through `do-work clarify` for confirmation. Append the reasoning too, including anything the answer put out of scope, so a fresh builder reads the decision instead of re-deriving it from `Recommended:`/`Also:`. Any *new* work the answer implies is captured as its own REQ, not left as a sentence in the hand-back. Do not flip `status` here: this REQ is already in flight, so the `pending-answers` → `pending` transition (`actions/clarify.md` Step 5) does not apply.

The follow-up REQs for builder-decided questions are created during **Step 8 (Archive)** — not here. Step 3.5 just records the decisions; the archive step handles the paperwork after the REQ is fully complete.

**Why not block?** Human time is the bottleneck. The optimal windows for user interaction are: (1) capture time, when the user is actively fleshing out requests, and (2) batch-review time, when the user returns to answer accumulated questions. Blocking mid-build wastes builder capacity on idle waiting.

**`pending-answers` REQs:** These accumulate in the queue. When the user returns, they run `do-work clarify` to review all `pending-answers` REQs at once, answer the questions, and flip the status to `pending` so the next work run picks them up. The work loop skips `pending-answers` REQs — it only processes `pending` ones.

If all `- [ ]` items are already `[x]` or `[~]`, or no Open Questions section exists, skip this step entirely.

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

**Route C:** Spawn a **Plan agent** with the request content, project context, the `crew-members/[domain].md` file (normalize `domain` per the Schema Read Contract first; if the resolved domain is missing, falls back to `general` for an unknown value, or the file doesn't exist, skip loading it), and any files listed in the `prime_files` array. Instruct it to use the prime files as the strict index for discovering the source of truth. Do not load global architecture. Ask it to produce a specific implementation plan (files to modify, order of changes, architectural decisions, testing approach). If a `## Plan` section does not already exist, append the output:

(append the plan per the **Plan Template — Route C (Step 4)** in `actions/work-reference.md`)

**Plan validation (Route C only):** After the Plan agent returns, run a quick quality check before proceeding:

1. **Requirement coverage:** Re-read the REQ's What/Detailed Requirements. Every requirement should map to at least one planned task. Flag uncovered requirements.
2. **No orphan tasks:** Every planned task should trace back to at least one requirement. Tasks that don't address any requirement suggest scope creep.
3. **Scope sanity:** Count the planned tasks. If 5+, flag: "Plan has [N] tasks — quality degrades past 3. Consider splitting this REQ into multiple smaller REQs."
4. **Consumer field contract:** For every planned command whose output drives an action-owned mutation, verify that the plan identifies the exact per-record identity, provenance, state, and outcome fields required by its consumer, as applicable. Flag plans that omit the consumer's required fields. This is the plan-time counterpart of `actions/review-work.md`'s **Restatement Sweep**, not a replacement; the sweep still runs at review.

Append validation findings to the `## Plan` section (if any issues found). These are **warnings, not blockers** — the builder can adapt. But flag them visibly so the orchestrator and review step are aware.

After the Route C plan is saved and this validation finishes, stamp `planning_at: <now>` using the current UTC instant (Timestamp rule, `actions/work-reference.md`) — only if the field is absent, because a REQ re-planned after a recovery keeps the first attempt's observation (**Stamps are append-only**, `actions/work-reference.md`). Stamp only this successful observed event. Routes A and B omit the field.

**Routes A and B:** Append a skip note (if not already present):

(append the skip note per the **Plan Skip Note — Routes A/B (Step 4)** in `actions/work-reference.md`)

### Step 5: Required-Lessons Consult (all routes) and Exploration (Routes B and C)

**Required-lessons consult — every REQ.** Before exploration, read the REQ's existing `required_lessons` entries in order, then consult `do-work/lessons-index.md` even when the field is absent. Existing entries are mandatory reads and consume the budget first. Match additional rows against the request text, Plan/current `write_set`, likely touched paths, and `prime_files`, using the ranking, entry forms, full-coverage targeting rule, mechanical costs, and single limit owned by `actions/capture-reference.md` → **Required Lessons Budget Contract**. Prefer narrowing an eligible bare match before dropping it. Refresh `required_lessons` with the fitting existing-plus-new set; never add an entry without an index match. Preserve and read captured entries when the index itself is absent, but add nothing.

Read the resolved set now: a bare path means the whole satellite; `path#family-slug` means only bullet lines carrying `[family: family-slug]`. A missing listed file never blocks exploration — record the missing entry in `## Exploration` when that section exists and always in the builder hand-back, then continue. Replace any existing `## Required Lessons — Dropped for Budget` body section with the current claim-time drops (entry, cost, matching reason), or remove that section when nothing was dropped; a drop is never silent. Omit `required_lessons` when the refreshed set is empty. This consult closes the capture-time gap for old REQs and for later members of a serial batch whose matching lesson did not exist when the batch was captured.

**Route A:** stop this step after the consult and resolved reads; only exploration is skipped. **Routes B and C:** continue below.

Spawn an **Explore agent** to find relevant files, existing patterns, types/interfaces, and testing conventions.

- **Route C**: Give it the plan and ask it to find files mentioned in the plan plus similar implementations
- **Route B**: Give it the request and ask it to find where the change should go and what patterns to follow
- **Both routes**: If the REQ's `prime_files` reference primes with a `## Lessons` section, include them in the explore context. Previous failed approaches and gotchas from this codebase area save the explorer from repeating dead ends.

If an `## Exploration` section does not already exist, append the output:

```markdown
## Exploration

[Explore agent findings — key files, patterns, concerns]

*Generated by Explore agent*
```

### Step 5.5: Scope Declaration (Routes B and C)

Before the builder starts coding, declare intent. This prevents scope drift from being discovered only at review time, after the code is already written.

**Route A:** Skip — scope is inherently constrained (single file, single change).

**Routes B and C:** Based on the plan (Route C) or exploration output (Route B), write a `## Scope` section into the REQ file:

(write the `## Scope` section per the **Scope Declaration Template (Step 5.5)** in `actions/work-reference.md` — declared file list + restated acceptance criteria. The review step compares the Implementation Summary's file list against this declaration; any undeclared touch or unused declaration is scope drift.)

**Mirror the file list into `write_set`.** Immediately after writing `## Scope`, copy its "Files I will touch" paths into the REQ's `write_set:` frontmatter field, replacing whatever capture seeded. The sync runs in **one direction only** — Scope is the source, `write_set` is the mirror — which is what makes drift between the two impossible. The field is display-only at any builder count (it feeds the board's overlaps badge — `actions/work-reference.md` Request File Schema); Route A skips Step 5.5, so a Route A REQ's `write_set` stays as captured.

The Scope section serves two purposes:
1. The builder commits to a file list before writing code — drift becomes measurable.
2. The acceptance criteria, restated from the REQ, become the word-by-word comparison target for review.

Scope-drift protection enforces **YAGNI**: only declared files get touched, and undeclared exploratory work becomes a discovered task (Step 8) rather than speculative scope creep. See `crew-members/coding-guardrails.md` § Simplicity First.

For Routes B and C, `advance` returns scope-drift evidence with qualification and again before testing. Review consumes that record and judges each declared-but-untouched or touched-but-undeclared path; Route A has no Scope comparison.

### Pre-Build Evidence Judgment

Routes B and C satisfy the pre-flight gate before dispatch through the **Mechanical Evidence-Gate Loop**. Repository dirt and dependency warnings are advisory under **Current-REQ relevance**; a test baseline with `launched: false` is unusable and excludes nothing. Resolve the canonical repository gate to exact argv. A matching green record authorizes dispatch; otherwise run that argv directly and return its status through the gate record — a non-zero exit is rerun once before any classification, and everything downstream reads the second run's status, output, and fingerprint (`actions/work-reference.md` → **One retry before classification**). The action retains failure fingerprinting and the deferral/no-op/reuse decisions below; `advance` owns checking and recording the evidence.

- **Ordinary parent:** call `do-work-cli defer-gate --manifest <path>` with exact parent/checkpoint preimages and writer identity. Allocate the proposed repair id by capture's read-only max-request plus reservation scan; the transaction alone creates the reservation. On a fully rolled-back id collision, rescan and retry with the new max+1. Consume only the typed `gate_deferral` result. Add its `repair_id` to this session's runnable repair closure even when the repair belongs to another UR, add its `parent_id` to session-local suppression even when explicitly targeted, recompute selection, and continue.
- **Repository-gate repair:** compare the failure to the repair's exact recorded fingerprint. The expected matching red baseline authorizes implementation and never recursively defers. A green baseline takes the durable reviewed no-op completion path in `actions/work-reference.md` → **Already-green repair no-op completion**, including canonical finalization so parents can resume. A mismatch or launch failure prepares a terminal `fail` finalization manifest, leaves parents dependency-gated, recomputes selection, and continues unrelated work; it never converts parents to `blocked` or `pending-answers`.
- **Deferred parent:** after every repair terminal result, recompute selection. Failed, cancelled, or dependency-gated repairs leave parents pending behind `depends_on`; continue unrelated runnable work. A repaired parent with no saved implementation range starts Step 6 normally after a fresh green baseline. A paired saved range follows the resume proof below.

For `gate_deferred: true` with paired `deferred_implementation_base` / `deferred_implementation_merge`, resolve both commits, require a non-empty base-to-merge range, require the merge in current `HEAD` ancestry, derive every implementation path rename-aware from that range, and reject any later path history or current index/worktree change on either side of a rename. Missing, malformed, non-ancestor, or otherwise unverifiable evidence stops safely. Drift deletes the stale pair from the claimed working REQ, treats all prior qualification/test/review evidence as stale, and returns to Step 6. No drift reuses the implementation and skips builder dispatch, but it must still rerun qualification, focused tests, the canonical gate, and independent review before completion. The full command/evidence contract lives in `actions/work-reference.md` → **Repository Gate Deferral and Resumption**.

### Step 6: Implementation

**Agent rules loading:** Before spawning the implementation agent, load domain-specific rules:

1. **Always load** `crew-members/general.md` — cross-domain rules and PRIME Files Philosophy
2. **Always load** `crew-members/coding-guardrails.md` — the always-on implementation guardrails (think before coding, simplicity, surgical changes, goal-driven execution, naming for reach). That file is authoritative for the current set; this gloss is illustrative and must not be read as closed.
2a. **Always load** `crew-members/shared-principles.md` — condition-keyed principles shared by implementation and review.
2b. **Always load** `crew-members/communication-style.md` — the always-on communication contract for agent prose (plain specific language, answer-first replies, reference codes, banned filler patterns). Governs status updates and hand-backs; artifacts stay under `anti-slop.md`, question wording under `clear-questions.md`.
3. **Conditionally load** `crew-members/[domain].md` — normalize the REQ's `domain` frontmatter per the Schema Read Contract first (e.g., `back-end` → `backend`, `ui_design` → `ui-design`), then load if the resolved domain is set AND the file exists (e.g., `domain: ui-design` → `ui-design.md`). An unknown value after normalization emits the contract's warning and falls back to `general` — no additional domain-specific crew loads (the always-loaded `general.md` from step 1 is the base).
4. **Conditionally load** `crew-members/testing.md` — if the REQ's `tdd` frontmatter normalizes to `true` per the Schema Read Contract (accepts `test_first`/`yes`/`on`/`t` as truthy aliases), or `domain: testing`
4a. **Conditionally load** `crew-members/security.md` — if the REQ's normalized `domain` is `security`, OR if the REQ description references authentication, authorization, session handling, cryptography, secrets handling, input validation/sanitization, or any OWASP-category surface. The "OR" clause is heuristic — when in doubt, load it; the cost of loading a checklist when not needed is low, the cost of skipping it on real security work is high.
5. **Conditionally load** `crew-members/caveman.md` — if the REQ's `caveman` frontmatter normalizes to a non-`false` value per the Schema Read Contract (any of `true`, `lite`, `full`, `ultra`, plus `yes`/`on` → `true`, `light` → `lite`). Compresses agent prose ~65-75% while keeping code and technical terms exact.
5a. **Conditionally load** `crew-members/maintenance.md` — if the REQ's `maintenance` frontmatter normalizes to `true` per the Schema Read Contract. This marks the REQ as a deliberate maintenance pass on the skill's *own* operating instructions (a drifting agent/action/crew/prime file) where removing or narrowing is a candidate fix; it loads the delete-before-you-add discipline **alongside** `coding-guardrails.md`, not instead of it. **Marker-only — do not infer it from the description.** A plain dead-code removal in application source is not a maintenance pass and stays under `coding-guardrails.md`'s implementation-time surgical-changes rule; only the explicit `maintenance: true` marker (set by capture for a removal/narrowing finding on the skill's own instructions) triggers the load. Unlike the security heuristic above, there is deliberately **no** description-based fallback here — a heuristic trigger would misfire on ordinary implementation REQs (which routinely touch adjacent dead code) and load the opposite posture from the one coding-guardrails wants.
6. **If a rules file is missing**, proceed without it — never block on a missing rules file

**Durability (background builder):** When the builder runs as a background or detached sub-agent, follow the durability pattern in `crew-members/background-agents.md` (disk-durable run directory as source of truth; survives a dead orchestrator session).

**Overlapping parallel writers:** If implementation is manually split among concurrent agents and their explicitly declared file lists or globs overlap, put each overlapping writer in its own worktree and branch before any write, then hand every completed branch back for serial reconciliation and merged-state verification. Follow `crew-members/background-agents.md` → **Worktree isolation is a separate axis** for the shared trigger and unsafe-branch policy, and `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** for this action's canonical hand-back sequence. The shared rule leaves read-only and declared-disjoint parallel work unisolated; `do-work run` is stronger still: it uses one worktree per builder regardless of overlap.

Spawn a **general-purpose agent** with the loaded rules, any files listed in the `prime_files` array, and context appropriate to the route:

- **Route A**: Request content only — "triaged as simple, aim for a focused minimal change"
- **Route B**: Request + exploration output — "follow existing patterns identified above"
- **Route C**: Request + plan + exploration output — "implement according to the plan"

Once the implementation builder has accepted that dispatch, take the current UTC instant (Timestamp rule) and hold it. Stamp `dispatch_at` with it only if the field is absent — after a recovery the field still names the first attempt's dispatch and that observation stays (**Stamps are append-only**, `actions/work-reference.md`). If dispatch fails before a builder accepts it, leave the field absent. When the builder returns its completed hand-back, stamp `builder_handback_at: <now>` on the same condition, before the hand-back merge begins. Then record that delegated wait once through `record-timing-event` (`--category builder-work --started-at <the instant you held>`), which owns every timestamp, duration and redaction mechanic so no step derives its own. **Pass the held instant, never `dispatch_at` read back out of the file:** on a retry that field belongs to the earlier attempt, and the recorder would charge this builder with every hour between the two.

All routes include these instructions to the agent (pointers — the underlying rules live in the loaded crew-members files and in the REQ frontmatter the orchestrator already wrote):

**Required-lesson regime — three additive layers, never substitutes.** Before implementation, the builder reads every current `required_lessons` entry unconditionally (captured stamps plus Step 5 claim-time matches), using whole-satellite or matching-family-bullets semantics. Independently, the touch-conditional Lessons Discipline rule still applies to every REQ, stamped or not, so a relevant satellite excluded by the budget can still be required by the touched prime. If any required path is missing, proceed without it and name the missing entry in the hand-back; never turn missing lesson context into a build blocker.

- **Crew rules govern behavior:** `crew-members/general.md` (always loaded) carries the Prime Files philosophy, Lessons-discipline, test-writing posture, cross-REQ test-break rules, and Discovered-Tasks contract. `crew-members/coding-guardrails.md` (always loaded) carries the implementation-time guardrails — that file is authoritative for which ones, and is deliberately not re-enumerated here. Domain/testing/caveman crews layer on top per Step 6's loading order. The builder reads these — do not re-state their contents inline.
- **Prime files come first:** Read every path in `prime_files` before touching code. If the primary utility you are modifying has no prime, investigate and create one (`prime-[name].md`), then update REQ frontmatter. Each prime's `lessons-[name].md` satellite encodes prior mistakes in that area — the unconditional `required_lessons` reads above are additive to reading it whenever the change touches code the prime's Read-first or Traps entries name (`crew-members/general.md` → Lessons Discipline).
- **P-A-U phasing is mandatory:** Edit the REQ's "AI Execution State (P-A-U Loop)" checkboxes in real time. [PLAN] writes a brief technical approach. [APPLY] stays in declared scope. [UNIFY] runs `git diff --stat`, runs native linters, verifies no debug artifacts, and lists each file checked (the orchestrator audits this during Qualification and Testing Judgment).
- **TDD mode when `tdd: true`:** Follow RED → GREEN → REFACTOR. Anchor RED on the REQ's `## Red-Green Proof` section if present — it arrived with the REQ and is not yours to write. Report the red-green evidence (test name, failure-before, pass-after) — Qualification and Testing Judgment verifies it.
- **Captured proof first:** If `## Red-Green Proof` is present, its RED prompt/case and GREEN outcome are the primary behavior tests must prove. Only adapt with documented reason.
- **Log Decisions as D-XX:** Significant implementation choices not dictated by plan/requirements become numbered entries in a `## Decisions` section. Continue numbering from the `<!-- D-XX counter: ... -->` comment Step 3.5 left behind; if none, start at D-01. Each decision needs reasoning — without it, the intent trail breaks. Sort each by the decide-vs-escalate gate (`crew-members/coding-guardrails.md` § Think Before Coding): a reversible, low-reach choice is **DECIDE & STATE** (reasoning only — it surfaces later as a *handled* item); a choice that's irreversible/expensive, taste-dependent, or genuinely contestable is **ESCALATE** — add `Value:` and `Risk:` lines so the hand-back can surface them. An ESCALATE entry whose real answerer is a named outside stakeholder additionally carries `Answerer: <name>` (and `Irreversible: yes — [why undoing is expensive]` when that holds) — the same clauses as Step 3.5, and Step 8 routes it to that person's stakeholder REQ instead of a `pending-answers` follow-up. **In worktree dispatch mode that section goes in your hand-back**, under the same heading, because the REQ file is one of the main-tree paths you may not write.
- **Write only inside the declared scope:** the REQ's `## Scope` "Files I will touch" list (mirrored to `write_set` at Step 5.5) is the builder's write boundary — read it; it is not yours to write. Needing a file outside it has two paths. When the REQ's own requirements or completion proof already require that file class, the declaration contradicts the REQ: flag it before editing, proceed with the required class, and report the contradiction plus the actual files. Otherwise it is a scope expansion: **stop and report to the orchestrator; never silently write it.** The orchestrator records the request and its resolution in the REQ trail as a `## Decisions` D-XX entry (it is a scope judgment) and extends both the Scope list and `write_set` before a serial builder continues, or from an unattended/worktree builder's handback before integration. An **absent or empty `write_set` is not a write prohibition** — but it is not a stated boundary either, so it is not licence: it means this REQ never ran Step 5.5 (a Route A REQ, or one whose `## Scope` was stripped by crash recovery), so derive the boundary from the REQ's own text and keep the change inside it — do not read the empty field as full-scope freedom.
- **Write only on your own branch, never the orchestrator's state:** commit on this REQ's branch and hand back the manifest — the orchestrator is the sole integrator and merges. On the worktree rung that means never writing the main tree at all; on the branch rung the tree is shared, so it means every commit you make stays on this REQ's branch and every `do-work/` path stays the orchestrator's (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** → *Isolation ladder*). A shared file that needs one line of wiring is an **integration seam**: hand back the exact line and where it goes rather than editing the shared file yourself.
- **Out-of-scope finds go to `## Discovered Tasks`** (a separate section, not nested inside Implementation Summary) — do not fix inline. Step 8 classifies and queues them. **In worktree dispatch mode that section goes in your hand-back**, under the same heading, because the REQ file is one of the main-tree paths the bullet above forbids you to write.
- **Report back the file manifest:** list every source file created/modified/deleted with the action verb, plus tests touched. The orchestrator writes the formal `## Implementation Summary` from your report — that section is not yours to write.
- **Report lesson evidence:** list each `required_lessons` entry read, whether it was whole-satellite or family-targeted, and every missing entry that was skipped. This evidence belongs in the hand-back even when the implementation file manifest is otherwise short.
- **Standard freedoms and obligations:** Full file/shell access. Escalate to explore or plan if the work proves harder than triaged. Document blockers explicitly. Identify and run related existing tests; honor any test-command map in the prime file (takes precedence over generic detection).

**Hand-back merge (the orchestrator's job, not the builder's).** When the builder returns its manifest, integrate here, at the end of Step 6, **before** Step 6.25: every evidence step from 6.25 onward reads the merged tree, so a merge deferred past that point leaves qualify and review with nothing to check. Before running any hand-back command, read the canonical full sequence in `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** → **When to merge, and the range every evidence step reads.** On the integration branch, its condensed sequence is:

0. Settle the owner's run artifacts and the index first. Queue-mode `advance` already committed the claim move and `do-work/CHECKPOINT.md`; do not stage either here. Read `git status --short --untracked-files=all -- do-work/` and use three categories: **stage** exact owner-written run artifacts (`manifest.md` and `REQ-NNN-brief.md`); **allow but never stage** each expected `REQ-NNN-handback.md`; **leave alone and name** every other `do-work/` path, including claim or checkpoint dirt this REQ did not produce. The hand-back is builder scratch, not an undeclared queue mutation. A third-category path is another session's or the user's own: print it once in the progress output, keep it out of the stage set, and continue — never stop on it and never commit it, because the pipeline commits only what this REQ declares. Only a dirty path this REQ owns still stops. Run `git diff --cached --name-only`; every commit here takes the whole index, so take any staged path outside the stage set out of the index with `git restore --staged -- <path>` (the working-tree bytes are untouched; a `MM` path also loses the distinct snapshot it had staged, so say so in that line) and name it too. Where the consumer commits `do-work/`, stage only the exact run-artifact paths (`git add -A -- <exact-run-artifact-paths>`), re-run the cached-name guard, then use plain `git commit`. Never use `git commit -- do-work/`: path-limited commit reads tracked paths from the working tree instead of committing the index just inspected. If the stage set has no changes, skip the run-artifact commit; that is normal when no durable run artifacts changed. Where `do-work/` is untracked, nothing is staged and there is nothing to commit — skip it. **Step 0 ends with the index holding nothing but this REQ's own stage set**, not with a clean working tree. Do this **before** step 1, so any run-artifact commit lands below `<pre>` and stays outside the merge range.
1. Run `git rev-parse --short HEAD` and hold the printed hash as **`<pre>`** — this REQ's pre-merge integration tip and the lower bound of its merge range. Capture it once per REQ; a remediation re-merge keeps the first one.
2. Guard the queue before merging: run `git diff --name-only <pre>...<operative_name> -- do-work/`. Any path printed is queue state committed on the builder's branch; stop and drop or revert those commits before integrating. Only when the guard is empty, run `git merge --no-ff --no-commit <operative_name>` — the branch the builder was actually dispatched on, which is the collision variant where there was one and the derived `worktree-agent-REQ-NNN-<suffix>` otherwise; never re-derive it here. If git says `Already up to date.`, stop and treat the hand-back as empty. Otherwise resolve any conflict, apply and stage the builder's handed-back integration-seam lines, then `git commit` — folding the seam into the merge commit is what puts it inside the merge range. Re-type `<pre>` and `<merge_hash>` wherever they are consumed, following the canonical [State across command blocks](../docs/prescribed-shell-primitives.md#state-across-command-blocks) rule.
3. Run `git rev-parse --short HEAD` again and hold that hash as **`<merge_hash>`** — the upper bound of the range and the supplied-provenance hash that finalization records.

Hold both hashes — and `<operative_name>` — as literals you re-type into each later command (shell variables do not survive between command blocks), and pass the range `<pre>..<merge_hash>` to Steps 6.3, 7, 8, and 9. The canonical reference carries the full rationale, remediation re-merge handling, and the queue guard's safe-direction over-inclusion caveat.

After the hand-back is successfully merged, stamp `integration_at: <now>` using the Timestamp rule, only if the field is absent (**Stamps are append-only**, `actions/work-reference.md`). A failed or empty hand-back does not create an integration observation.

### Step 6.25: Implementation Summary

After implementation completes, write a manifest of what changed to the request file. This is the primary auditability artifact — without it, there's no way to verify the REQ was implemented without digging through git history.

**If a `## Implementation Summary` section already exists** (e.g., from a re-qualification or remediation loop), **replace it entirely** with the new content. Do not append a second copy. The most recent implementation is the one that matters.

Append (or replace) in the request file:

(write the manifest per the **Implementation Summary Template (Step 6.25)** in `actions/work-reference.md`)

**Rules:**
- **Mandatory for all routes.** Route A gets a short list. Route C gets a detailed list.
- List all project files that changed — source code, config (`package.json`, `Dockerfile`, CI YAML), documentation, etc. Exclude only `do-work/` metadata files.
- Mark files as `(new)`, `(modified)`, or `(deleted)`.
- The "What was done" summary should be factual, not aspirational — describe what you built, not what the REQ asked for.
- This section is the primary auditability artifact. If `Files changed` only lists `do-work/` paths or is empty, the REQ was not implemented.
- **Already-green repository-gate repair exception:** only the exact durable no-op evidence in `actions/work-reference.md` → **Already-green repair no-op completion** permits `**Files changed:** None — verified repository-gate repair no-op.` The exception proves that implementation was unnecessary; it never blesses an empty ordinary REQ.
- **Design-artifact exception:** For `domain: ui-design` requests that produce design deliverables rather than code (wireframes, IA specs, visual specs, interaction specs), the artifact files themselves count as project files. Place them in the project's design docs directory (e.g., `<project-root>/docs/design/`) — not inside `do-work/`. The Implementation Summary lists these files normally.

### Qualification and Testing Judgment

After the builder returns and the Implementation Summary is written, run the **Mechanical Evidence-Gate Loop** with the exact request path, the saved merge range, and the selected focused-test argv. Treat `advance.gate_records` as evidence, not as the decision: the orchestrator still reads the actual diff, proves that changes are substantive, traces every requirement, verifies live data flow, and applies the **Qualification Anti-Rationalization Table** in `actions/work-reference.md`. Static warnings remain judgment inputs; conventions, dynamic imports, side-effect imports, and standalone entry points are not dead code merely because a grep cannot see their consumers.

An already-green repository-gate repair must present the exact no-op evidence and empty project diff described in `actions/work-reference.md`; any project change or malformed evidence returns to ordinary qualification. Record `## Qualification` from the observed gate records and the orchestrator's judgment. On failure, return the concrete findings to the builder for at most two re-qualification attempts before review owns the remaining judgment.

Select focused tests from the touched primes' valid test maps, then cover unmapped files with repository-native detection. Every test-file invocation must remain under 30 seconds. The focused-test gate distinguishes green, matching recorded red, new red, timeout, launch failure, and an unusable `launched: false` baseline; only green or an exact matching red clears mechanically. Whether two failures are semantically the same, or a warning represents a real regression, remains the orchestrator's judgment. New regressions return to implementation for at most three attempts, loading the debugging and testing crew guidance on attempt two onward.

Run the declared canonical repository gate directly, unpiped, against the final tree. To attribute that gate's wall time, launch the same argv through `run-timed-command`, which satisfies *directly* as `actions/work-reference.md` → **One retry before classification** defines that condition, so the retry rule and the argv and exit status reported to `advance` are unchanged. A non-zero exit is rerun once before any classification, and everything downstream — repair-failure handling, diagnostic attribution, deferral — reads the second run's status, output, and fingerprint (`actions/work-reference.md` → **One retry before classification**). Return its exact argv and exit status through `advance` so the typed evidence is checked or recorded; focused-test exclusions never waive this gate. Repository-gate repair attribution, base-revision diagnosis, deferral, and cleanup remain action-owned judgments governed by `actions/work-reference.md`. Plan affected heavy lanes with the repository's typed planner; selected lanes are recorded in the Testing section and held at Step 7.7 after review.

Append to the request file:

(append per the **Testing Section Template** in `actions/work-reference.md`; omit Red-green validation for non-behavioral changes, and trace it back to `## Red-Green Proof` when present)

Omit `Red-green validation` if no request-specific tests were written or identified, or if the change is non-behavioral (refactor, config, docs, cleanup) — use regression evidence instead. Omit `Existing tests updated` if no prior tests were modified.

When the REQ includes `## Red-Green Proof`, the `Red-green validation` entries should trace back to that captured RED/GREEN pair. If the implemented test uses a nearby equivalent instead of the exact captured prompt/case, explain why.

**Already-green repository-gate repair exception:** Before applying ordinary TDD verification, capture `<now>` (current UTC instant — Timestamp rule) once for validation and finalization, then invoke `<skill-root>/tools/do-work-cli.sh --repo-root <project-root> --format json validate-already-green-repair --request-path <exact working REQ path> --writer <exact finalization writer> --at <now>`. Permit the missing RED/GREEN pair only when the command has typed success and `already_green_repair.tdd_allowed: true`. The validator is the sole decision authority: it reads the claimed repair marker, exact `## Repository Gate Repair Intake`, `## Repository Gate Repair No-Op`, `## Implementation Summary`, and `## Qualification` shapes, matches the intake fingerprint and argv to no-op evidence, verifies the recorded past-revision gate record, and proves the no-op project diff is empty. Do not reconstruct any part of that predicate from prose. Missing, failed, malformed, false, ordinary, or non-empty results remain subject to the ordinary RED/GREEN requirement below.

**TDD verification:** If the REQ has `tdd: true`, the `Red-green validation` section is mandatory — the builder must show test-first evidence that they used RED/GREEN TDD. Qualifying evidence is a runnable test in the project's existing automated test harness, written before implementation, observed failing before the change and passing after it, and re-runnable by another agent. A repeatable check outside that harness is regression proof, not `tdd: true` evidence. If qualifying evidence is missing, treat it as a test failure: return to implementation (same path as step 4 above) with explicit instructions to provide red/green evidence — write the failing test first, confirm it fails, then make it pass.

### Step 7: Review

Run actions/review-work.md in **orchestrated mode** against this REQ.

For an already-green repository-gate repair, independent review reads the durable no-op evidence instead of an implementation diff and freshly invokes the same validator with the exact request path, finalization writer, and completion timestamp. It may accept the no-diff path only from typed success with `already_green_repair.review_allowed: true`; record the typed fingerprint, gate-evidence, canonical completion paths, staged paths, and decision in the ordinary `## Review`. Without that result the no-op may not complete.

The review reads the REQ (in `do-work/working/`), the original UR, and the current diff (`git diff` or `git diff --staged`) to evaluate the implementation: requirements check (did we build what was asked?), code review (is it solid?), and acceptance testing (does it actually work?). **In worktree dispatch mode** the working tree is clean after the merge, so the review reads this REQ's merge range `<pre>..<merge_hash>` instead (`actions/review-work.md` Step 4, Get the Diff).

**Restatement sweep (MUST).** If this REQ's diff redefines something other text restates — a contract token, a schema field's semantics, a gate's wording, a prescribed command's output shape — the review runs the sweep defined in `actions/review-work.md` Step 6 (**Restatement Sweep**) and reports every stale restatement as a finding, including ones in files outside this REQ's declared Scope. Only an `impact-critical` result routes through Step 10's Fold-First Rule; every other result stays in the report and ends `→ report only`.

**How to run it:** Spawn an agent with actions/review-work.md file, the REQ path, and the `crew-members/[domain].md` file (normalize `domain` per the Schema Read Contract first; if the resolved domain has a matching file, load it; otherwise skip) — in worktree dispatch mode, also pass this REQ's merge range `<pre>..<merge_hash>` so the review reads the merged diff rather than the clean working tree. Or read actions/review-work.md file and follow its orchestrated-mode instructions in the current session.

**What happens next depends on the review result:**

- **Acceptance = Pass AND overall ≥ 75%**: Append the Review section to the REQ and continue to archive as `completed`. Minor findings go in the report only.
- **Acceptance = Partial OR overall 50-74%**: Append Review and continue to archive as `completed`. Every finding is impact-stamped per `actions/review-work.md` Step 10. Only `impact-critical` findings create automatic follow-up work; every noncritical line ends `→ report only` and stays in the Review section.
- **Acceptance = Fail OR overall < 50%**: **Do NOT archive as completed.** Instead:
  1. Append the Review section to the REQ.
  2. Return to Step 6 (Implementation) with the review findings as context for the builder. Load `crew-members/debugging.md` for the remediation attempt — the builder needs structured debugging methodology, not just "try again."
  3. The builder gets **ONE remediation attempt**.
  4. Re-run Steps 6.25 through 7 (Summary → Qualification → Testing → Review) on the remediated code.
  5. If still failing after remediation: update frontmatter to `status: completed-with-issues`, `completed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`), append a `## Remediation` section documenting both attempts, and route the remaining findings through `actions/review-work.md` Step 10. Only `impact-critical` findings auto-queue; the rest stay in the Review section. Then proceed to archive (Step 8) — the frontmatter is already set, so Step 8 should not overwrite it.

After the first review result is recorded, stamp `review_at: <now>` using the Timestamp rule, regardless of its verdict. If remediation runs, stamp `remediation_at: <now>` only after that builder hand-back is successfully integrated, then stamp `re_review_at: <now>` only after the post-remediation review result is recorded. All three are written only if the field is absent (**Stamps are append-only**, `actions/work-reference.md`). A passing first review leaves both remediation fields absent.

The status `completed-with-issues` means the REQ was archived but has known unresolved problems. It counts toward UR completion for archiving purposes. Any critical follow-up remains queued; noncritical findings remain visible in the archived Review until a maintainer explicitly captures one. This status remains visible to recap and every completed-work presentation action; those readers inherit the Terminal-success status set from `actions/work-reference.md` rather than defining a caller-specific filter.

**Automatic review follow-ups are based on impact, not score or severity.** Record an impact token on every finding per `actions/review-work.md` Step 10; the judgment never re-scores it. Only `impact-critical` findings run the Fold-First Rule and auto-queue with `status: pending`, `user_request: [same UR]`, `addendum_to: [reviewed REQ id]`, `domain: [same domain]`, `review_generated: true`, and `impact: impact-critical`. Every other finding stays in the report, with its line ending `→ report only`; it never appends to a REQ or sweep and never enters `pending-answers` or the prose backlog. The report tells a maintainer to promote one explicitly with `do-work capture`, quoting the complete finding line as the capture source. Failure classification, builder-decided follow-ups, and stakeholder-requested follow-ups remain separate and unchanged. Cycle detection (Step 8, substep 5) still applies before a critical follow-up is created.

**Calibrate depth to route:** Route A gets a quick scan (skip dimensions that don't apply). Route B gets a standard review. Route C gets a thorough review comparing against the plan.

Append to the request file:

(append per the **Append to REQ File** template in `actions/review-work.md` — the file dispatched above, so it is already in context; review-work.md owns the Review section format)

### Step 7.5: Lessons-Capture Phase

> **Named entry point.** Other actions reference this as **work.md's Lessons-Capture Phase** (not by step number) — e.g. `actions/kb-lessons-handoff.md` and `actions/review-work.md`. The `7.5` is for internal navigation only; callers must use the phase name so they don't break if steps are renumbered.

Before archiving, capture what's worth remembering. This section is the institutional memory — when someone revisits this code in six months, the REQ file tells them what happened, what was tried, and why things ended up the way they did.

Append to the request file:

```markdown
## Lessons Learned

**What worked:** [1-2 bullets — approaches, patterns, or tools that paid off]
**What didn't:** [1-2 bullets — dead ends, failed approaches, and *why* they failed]
**Worth knowing:** [Anything the next person touching this code should know — gotchas, edge cases, non-obvious dependencies]
```

**Rules:**
- Keep it concise — pointers to code, not walls of text. The code is the source of truth.
- **Required for Routes B and C** — there's always something worth recording when exploration or planning was involved. **Optional for Route A** — skip if the change was straightforward with no unexpected discoveries, no failed approaches, and no gotchas worth noting. If anything surprised you (undocumented behavior, unexpected test failures, a file that wasn't where you expected), record it.
- "What didn't work" is the most valuable part — it prevents repeating mistakes.
- File lists are no longer needed here — they're covered by the mandatory Implementation Summary (Step 6.25).

**Write the `## Orientation` block (the hand-back's "what's being built"):** After the Lessons Learned section, append a short `## Orientation` section reporting the change at feature/subsystem altitude — "Now you can X; lives in Y subsystem" — not a file list. Use the REQ's `prime_files` to name the subsystem; flag `[MAP CHANGED]` only when the change alters the system's shape (new module, data flow, contract, or a renamed concept). Run a narrowed staleness spot-check on each touched prime per `../../do-work-toolbox/actions/prime.md` Step 2 / Step 6 (do its referenced paths still exist?) and flag any prime the change made stale. **Scale to reach:** a leaf REQ is one line; a map-changing REQ gets a short paragraph and a why-it-matters. When `prime_files` is empty, derive a one-line feature-altitude summary from the What / Implementation Summary instead — never a file list. This block feeds the **WHAT'S BEING BUILT** section of the Decision Brief (`actions/work-reference.md` → **Decision Brief (hand-back format)**). Crash recovery strips `## Orientation` on re-queue (it's orchestrator-generated).

**Update prime lesson satellites (deferred to Step 8):** After writing the Lessons Learned section, check the REQ's `prime_files` frontmatter. For each listed prime file relevant to this lesson, **collect a pending lesson write** — do NOT execute the write here. The REQ is still in `do-work/working/`, so any link pointing to its eventual archive location would either be broken or tempt a link to the transient working path.

The write target is the prime's **lesson satellite**, `lessons-<name>.md` beside `prime-<name>.md` — never the prime itself. A prime is read in full every time its area is touched; appending to it on every REQ is what turns a routing index into an archive (`crew-members/general.md` → PRIME Files Philosophy).

Choose one short kebab-case **failure-family slug** for each pending lesson. Record each pending write as a tuple: `{ primeFilePath, satellitePath, relativeLinkText, lessonSummary, familySlug }`. Hold them in memory (or a small scratch file under `do-work/working/`) until Step 8. The literal marker written to the satellite is `[family: <family-slug>]`; use the same slug for recurrence checks and the lessons index.

Compute each deferred prime-link path relative to the prime file's location (not the repo root) per `actions/work-reference.md` → **Deferred Prime-Link Path Computation (Step 7.5)**; Step 8 validates the resolved target against the exact archive path in the lifecycle plan before adding the prime or satellite to the finalization allowlist.

Only add a link when the lesson is relevant to that prime file's scope — don't spray every lesson into every prime file. If the REQ has no `prime_files` or the lessons aren't relevant to any prime file, skip this and clear the pending list.

**Knowledge-base handoff.** After the Lessons Learned section is written and prime-file links are in place, follow `actions/kb-lessons-handoff.md` to offer dropping a structured source document into `kb/raw/inbox/` so the next `do-work-knowledge bkb triage` + `do-work-knowledge bkb ingest` cycle compiles the lessons into the wiki. The handoff asks the user before writing and records `kb_status` (plus `kb_entry` on success) back onto the REQ. In unattended work runs with no human in the loop, the handoff defaults to `kb_status: pending` — it never writes to the KB without consent. If the project has no `kb/` directory, the handoff points the user at `do-work-knowledge bkb init` and defers; it never blocks archival.

### Step 7.7: Heavy-Test Hold and Drain

**Heavy-test hold:** after Steps 7 and 7.5 pass, a REQ whose Testing section recorded selected heavy lanes does **not** run those lanes and does not hold the queue loop open. Record the fast commands and durations plus the typed lane plan, then:

1. land the implementation commit and record its full hash in the REQ's canonical `commit:` field;
2. append `## Heavy Verification Plan` with the exact base/target revisions and each selected lane's id, argv, and reasons; and
3. recompute selection and continue every unrelated runnable REQ.

The REQ stays `claimed` in `do-work/working/`: no status change, no `## Open Questions` machine line, no move to `do-work/queue/`, and no checkpoint edit. The landed `commit:` is what makes its source ready for dependent REQs while it waits.

**Heavy-lane drain (at queue exhaustion):** when no claimable pending REQ remains and at least one held claimed REQ exists — including one routed here by `recover`'s `RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES` finding — this loop runs the held lanes itself. No permission is asked. Per held REQ, recompute `plan-heavy-verification` at its stored base and `commit:` and refuse any drift against the stored `## Heavy Verification Plan`. Plan drift, or a stored `historical-revalidation` plan, leaves that REQ held and is a typed finding for a human, never a hand edit. Union the selected lane ids across every remaining REQ and run them once at HEAD, repeating `--lane` for each unioned id:

```bash
<skill-root>/tools/do-work-cli.sh --repo-root "<project-root>" --format json \
  run-heavy-verification --manifest _dev/tests/heavy-lanes.json --lane "<lane-id>"
```

The runner refuses an unknown lane and a dirty tracked tree outside `do-work/` before anything executes; `HEAVY-RUN-DIRTY-TREE` is a typed condition to resolve, never a hand edit. It reports each lane's exit status, skip state, wall seconds, and `disposition`, and a red or skipped lane is a warning finding rather than a command failure.

A lane whose deterministic fingerprint — its argv, its covered and globally unclassified committed paths, regular untracked input bytes, bounded toolchain probes, and the entire inherited environment — still matches a successful record no older than four hours is reported `reused` instead of executed; every other lane is `executed` with the exact reason it was not reused (`fingerprint_mismatch`, `evidence_expired`, `no_prior_evidence`, `fingerprint_uncertain`). Execution revokes the previous success before launch; an invalidation failure refuses the run. Fingerprint uncertainty (including an opaque browser runtime or a symlink input) always executes. Shipped lane argv explicitly isolates system/global Git configuration; custom manifest commands retain their declared environment. Report each lane's disposition alongside its result: a `reused` green was measured by an earlier run, not this one. Add `--no-evidence-reuse` when a lane must execute regardless — a suspected flake, or an input outside its declared fingerprint.

Then dispose of each held REQ from the subset of lanes it selected, with no `answer` transaction in any outcome:

- **Green** — every lane that REQ selected is present in the run, exited 0, and was not skipped. Write `heavy_verified_at: <now>` (current UTC instant — Timestamp rule, `actions/work-reference.md`), `heavy_verified_revision: <runner's execution revision>`, and a `## Heavy Verification Result` section naming the target revision, the execution revision, and one line per lane onto the claimed record, then run Steps 8 and 9 for it in this same turn.
- **Red** — any selected lane exited non-zero. Delete `commit:` and the `## Heavy Verification Plan` section from the record so dependents stop building against withdrawn work, then enter Step 7's remediation path now and re-hold at this step after a fresh review.
- **Skipped** — any selected lane was skipped. The REQ stays claimed and held, its `HEAVY-RUN-LANE-SKIPPED` finding names the lane, the next drain retries it, and the composed exit summary names it if the run ends while it is still held.

A `commit:` drained once this run is never drained again.

### Step 8: Prepare Finalization

Keep judgment here and fold before minting. Read builder-authored decisions and discoveries from the REQ or its durable hand-back per `actions/work-reference.md` → **Reading a Builder-Authored Section (any step)**; an unreadable hand-back is a finding, not silence.

Run `fold-timing-summary` for this REQ, with its run identity and current request path, before the manifest below is authored: it writes one compact `## Timing` section from the run's Git-private timing stream and removes that stream, so the summary travels with the request instead of earning a commit or an agent turn of its own. A run that recorded no timing events changes nothing and is not an error.

1. **Fold-First follow-ups:** route each builder-decided question by its `Answerer:` clause, fold it into an existing open user or stakeholder REQ when the root cause/audience matches, and mint only when no valid fold exists. Preserve the cold-reader wording, source decision, value/risk, irreversible warning, route-back, and stakeholder-report judgment defined by the **Builder-Decided Follow-up Template (Step 8)**, **Stakeholder REQ Template (Step 8)**, and Fold-First Rule. Detect an existing `addendum_to` cycle before either action.
2. **Sweep consolidation and impact stamping:** consolidate the builder's `## Discovered Tasks` with the current sweep, stamp every item by impact, queue only `impact-critical` work, and keep every noncritical item on this REQ ending `→ report only`; follow **Discovered Tasks Classification (Step 8)**.
3. **Terminal judgment:** preserve an already-selected `completed-with-issues`, otherwise choose `completed`; on failure apply **Failure Classification (Step 8)**, including its upstream-cascade and external-precondition blocked-flip judgments. A blocked flip returns the REQ to the queue through its canonical lifecycle command instead of entering finalization.
4. **Release and lesson judgment:** decide whether this successful REQ delivered a release, then choose affirmative project ownership, source, bump, mirrors, changelog title/voice, and exact payload bytes per **Changelog Entry Procedure (Step 9)**. An already-green repair carries no release. Finish any deferred lesson content, family promotion, and links, and include their exact paths in the manifest; the canonical transaction validates and commits them.
5. **Finalization intent:** author exactly one strict manifest containing the chosen terminal/failure values, writer and timestamp, exact request/checkpoint preimages, exact lifecycle/release/follow-up/lesson allowlist, commit message, and provenance judgment. A merged isolated build uses its retained merge hash as `supplied_commit`; a finalization with no implementation commit to record uses `primary_commit`. Retain the operative worktree identity until typed finalization success.
### Step 9: Commit Phase (Git repos only)

> **Named entry point.** Other actions reference this as **work.md's Commit Phase**; the number is only for navigation.

Run the exact `advance` continuation with the selected request path and the single action-authored finalization manifest, continue only when its global outcome is success and exactly one ordered `finalizations` record matches the REQ/path with `phase: cleanup_complete` and empty `blocked_paths`/`reason_codes`, report that record's archive and settled/created commit hashes, then remove any retained worktree by its operative name without force.
### Step 10: Loop or Exit

After integration, replay the canonical selector required by the current run mode and loop while it selects work; when it does not, run `<skill-root>/tools/do-work-cli.sh --repo-root <project-root> --format json advance --checkpoint`, then cleanup and render Step 1's composed exit summary. The checkpoint command is the sole session-end writer and preserves every live foreign or unlabelled in-progress record.

**Context-wipe principle:** every REQ starts with fresh agents and a fresh canonical selection. Carry only durable REQ evidence, and treat unexplained overlap with the previous Implementation Summary as qualification drift.

## Clarify Questions

The clarify workflow has its own action. Run `do-work clarify` — it handles batch-review of `pending-answers` REQs, where the user confirms, overrides, or discards builder decisions. Resolved REQs flip back to `pending` and re-enter the work queue. Questions routed to an outside stakeholder live on `blocked` REQs carrying `stakeholder:` instead — their answers come back through `do-work stakeholder-answers`, and clarify only routes them.

## Orchestrator Checklist (per request)

```
□ Step 0: Parse arguments (targeting tokens, --wave N, --fan-out [N], --skip-impact-negligible); reject unrecognized residue
□ Step 1: Run recover FIRST; consume its typed finalization and working-claim decisions, then validate frontmatter and select
□ Step 1: Consume advance's committed claim result and frozen continuation argv
□ Step 3: Triage (decide route, append ## Triage, read original if addendum)
□ Step 3.5: Handle Open Questions (mark - [~] with D-XX numbered decisions; a user answer obtained mid-run is written in as - [x] before dispatch — never - [~], no D-XX)
□ Evidence gates: run `advance` at each mechanical phase with the exact request path and consume request-bound `gate_records`
□ Step 4: Plan (Route C: spawn Plan agent + validate plan / Routes A & B: note skipped)
□ Step 5: Consult required lessons (all routes); then Explore (Routes B & C only)
□ Scope judgment (Routes B & C: declare files + acceptance criteria in REQ)
□ Pre-build evidence judgment (Routes B & C: interpret repository, baseline, and dependency evidence)
□ Step 6: Implement (spawn agent with lessons + TDD mode if set, log decisions as D-XX)
□ Step 6.25: Implementation Summary (append file manifest — mandatory for all routes)
□ Qualification judgment (orchestrator verifies substantive changes, live flow, requirement coverage, and warnings using advance's mechanical records)
□ Testing judgment (measure every test file against the <30s budget; run the direct repository gate; plan affected heavy lanes and record the selected ones in the Testing section; load debug rules on attempt 2+; verify TDD evidence if tdd:true)
□ Step 7: Review (spawn actions/review-work.md — gate on acceptance: Pass→archive, Fail→remediate with debug rules)
□ Step 7.5: Lessons Learned + Orientation (append sections at subsystem altitude, update prime files, skip lessons for Route A if no surprises)
□ Step 7.7: Heavy hold after review (held requests stay claimed); drain at exhaustion, finalize green in the same turn
□ Step 8: Prepare finalization intent (choose terminal status, classify failures, route questions/tasks, collect deferred lessons, and preserve exact lifecycle/release inputs without mutating the tail)
□ Step 9: Finalize once (pass the strict manifest to `advance`; canonical finalization owns archive/checkpoint/UR/calibration/release, exact commit, provenance, verification, and cleanup)
□ Step 10: Loop or Exit (fresh selection if looping, else advance --checkpoint + cleanup)
```


## Error Handling

| Phase | Action |
|-------|--------|
| `pending-answers` REQs remain after queue is empty | Report them to the user: list each REQ and its unresolved questions. Suggest `do-work clarify` to batch-review. |
| Held claimed REQs remain after the queue is empty | Run the heavy-lane drain in Step 7.7. A REQ still held afterwards (skipped lane, plan drift, dirty tree) is listed in the composed exit summary with its typed finding; not a failed run. |
| `blocked` REQs remain after queue is empty | Report each with its `blocked_by` condition and age (per the composed exit summary). Suggest re-running `do-work run` (auto-probes any `blocked_check`) once a condition is met, or `do-work clarify` to confirm a human-checkable one. For a stakeholder-questions REQ (`stakeholder:` present), suggest sharing its report and then `do-work stakeholder-answers` — never a probe or a yes/no confirm. |
| Plan agent fails (Route C) | Classify failure (Intent/Spec/Code/Environment), create follow-up REQ if applicable, archive as failed |
| Explore agent fails (B/C) | Proceed to implementation with reduced context — builder can explore on its own |
| Builder reports a missing external precondition | Apply the blocked-flip test (Step 8 → **Mid-run blocked flip**): if no substantive implementation edits landed AND the condition is expected to self-resolve, flip to `status: blocked` (non-terminal, back to the queue) instead of failing. Otherwise classify as Environment and archive as failed. |
| Implementation fails | Classify failure (Intent/Spec/Code/Environment), create follow-up REQ if applicable, archive as failed |
| Focused tests or an attributed current-diff canonical repository gate fail repeatedly | After 3 fix attempts, classify as Code failure, create a follow-up REQ with the focused-test or attributed current-diff gate failure details, and archive as failed. A matching unrelated gate failure instead uses the repository-gate deferral lifecycle; a mismatch, launch failure, invalid range, or unverifiable attribution remains a fail-safe stop and is never archived as success. |
| Review: Acceptance = Fail | Return to Step 6 for ONE remediation attempt, then re-review. If still failing: archive as `completed-with-issues` with follow-up REQs |
| Review work agent fails | Skip review, note it in the REQ file, continue to archive — review failure is not a gate |
| Finalization or commit fails | Consume the typed blocker, correct the judged manifest or underlying hook failure, then rerun the exact `advance` continuation. Never bypass the hook, hand-stage the tail, or infer an archive from partial state. |
| Unrecoverable error | Stop loop, report clearly, leave queue intact for manual recovery |

## Progress Reporting

Keep the user informed with this format:

(keep the user informed in the running per-REQ progress format shown in `actions/work-reference.md` → **Progress Reporting Example**)

When the run finishes or pauses, hand back with the **Decision Brief** (`actions/work-reference.md` → **Decision Brief (hand-back format)**): lead with WHAT'S BEING BUILT (each REQ's `## Orientation`, at subsystem altitude), then DECISIONS FOR YOU (any escalated `pending-answers` follow-ups, each with value + risk), then WAITING ON OTHERS (stakeholder-routed questions, with each report's path to share), then a collapsed HANDLED list. Never lead with review scores — they stay in the per-REQ progress lines, not in front of the hand-back.


## Archived Request File Example

See [sample-archived-req.md](./sample-archived-req.md) for a complete example of what an archived REQ looks like after processing through the full pipeline (Route B). Every section shown there is generated by the steps above.

**Timestamps tell the story:** `created_at` → `claimed_at` = queue wait time. `claimed_at` → `completed_at` = calibration wall span, not active implementation time. Optional phase stamps show the observed pipeline breakdown, and `completed_at` → `release_at` is the release tail. Route + timestamps let you calibrate triage accuracy over time.

## Rules

- The orchestrator owns semantic judgment and action-authored sections; deterministic lifecycle, archive, release, commit, provenance, verification, and cleanup mutations go through the canonical finalization engine. Spawned agents do implementation work only.
- Only two frontmatter status transitions are written on the normal path: `advance` commits `pending` → `claimed`, then finalization commits the terminal status. Intermediate phases are tracked by which `##` sections exist, not by status.
- Finalization commits only its exact manifest allowlist and never bypasses a commit hook. `supplied_commit` provenance records the supplied merge hash without recommitting implementation paths; a `primary_commit` finalization, which has no implementation commit to record, may require a second metadata commit (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**).
- `write_set` is display-only, and stays display-only at any builder count: under fan-out it is **advisory input to the human's pick, never a gate**, it is **not read at all by auto-wave's computation**, and the merge — not the field — is the non-interference proof (`actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). Nothing schedules, gates, or dispatches on it; it feeds the board's overlaps badge (`actions/work-reference.md` Request File Schema).
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
| "P-A-U is bookkeeping — I'll just tick the boxes" | Do each phase; qualification audits the diff against the checked boxes | A checked `[UNIFY]` over a diff containing `console.log` is a false claim the qualifier will catch |
| "The queue file's twin is already archived, but re-running is harmless" | Consume `advance`'s `blocked-archive-collision` hold | Re-processing a duplicate silently re-commits it and corrupts the archive lineage |
| "The merge went fine, so I'll just `-D` the leftover `worktree-agent-*` branch" | Run `git branch -d` from the integration branch and stop on refusal (Step 8 substep 8) | `-d`'s refusal is the only assertion the integration actually happened; `-D` deletes the evidence that a merge was skipped or lost, along with the work |
| "Cleanup just needs the branch name — I'll rebuild it from the REQ slug" | Merge and clean up by the operative name held since dispatch (Step 6 hand-back, Step 8 substep 8) | After a collision variant, the slug-derived name is the *leftover*: the merge would integrate the wrong branch, `-d` refuses on unmerged work, and the run halts on a false lost-merge alarm while the variant is never cleaned |

## Red Flags

- REQ in `do-work/working/` for >1 hour with no new git commits (builder may be stuck)
- Implementation Summary lists files but `git diff` shows no changes in those files (hollow implementation)
- All P-A-U checkboxes marked complete but diff contains `console.log`, `debugger`, or `TODO` (debug artifacts)
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
