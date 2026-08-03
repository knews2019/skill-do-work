# Work Action

> **Part of the do-work skill.** Invoked when routing determines the user wants to process the queue. Processes pending requests from the `do-work/queue/` folder in your project. User-facing walkthrough: [`docs/work-guide.md`](../docs/work-guide.md).

An orchestrated build system that processes request files created by actions/capture.md. Uses complexity triage to route simple requests straight to implementation and complex ones through planning and exploration first.

## When to Use

**Use when:**
- The queue has `pending` REQs and the user wants them built (`do-work run`, `start`, `go`, etc.).
- The pipeline dispatches actions/work.md as its build step.
- A specific REQ id was named (`do-work run REQ-042`) — the action scopes to it.

**Do NOT use when:**
- The queue is empty — tell the user and stop; suggest `do-work capture-request: [describe]` instead.
- The only REQs left are `pending-answers` — route to `do-work clarify` so the user can resolve them first.
- See `SKILL.md` routing table for sibling action selection (inspect, verify-requests, review-work, etc.).

## Request Files as Living Logs

Each request file becomes a historical record. As you process a request, append sections documenting each phase: Triage, Plan, Exploration, Implementation Summary (mandatory file manifest), Testing, Review. This ensures full traceability — what was planned vs done, what files were touched, and whether triage was accurate.

This living log is also the **trail of intent**. The REQ starts as a validated statement of what the user wants (written by capture). As actions/work.md processes it, each appended section documents how intent was interpreted and realized: builder decisions (## Decisions) record where the builder exercised judgment beyond stated intent, scope declarations (## Scope) record what the builder committed to, and implementation summaries record what was actually built. The gap between captured intent and realized implementation is visible in a single file.

## Architecture

The per-REQ orchestration pipeline (triage → plan/explore → implement → qualify → test → review → archive → commit, with the orchestrator handling all file management) is diagrammed in `actions/work-reference.md` → **Architecture**.

> **Remember:** Every completed request gets a git commit (Step 9) before looping to the next request.

**Sub-agent note:** This document uses "spawn agent" language. Use your platform's subagent mechanism when available. If your tool doesn't support subagents, run phases sequentially in the same session and label outputs clearly.

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

The full annotated frontmatter schema and the **Schema Read Contract** — the normalize-and-warn rules every read site honors for the enum/boolean fields `domain`, `status`, `route`, `caveman`, `tdd`, `error_type`, `kb_status` — live in `actions/work-reference.md` → **Request File Schema — Full Frontmatter** and **Schema Read Contract**. Every reference below to "the Schema Read Contract" points there.

**Status flow (frontmatter values):** `pending` → `claimed` → `completed` / `completed-with-issues` / `failed`

The intermediate phases (planning, exploring, implementing, testing, reviewing) are tracked by which `##` sections exist in the REQ file, not by frontmatter status changes. Only two status transitions are written to frontmatter on the normal path: `pending` → `claimed` (Step 2), then `claimed` → final status (Step 8). Exception paths write their own statuses: the special holding statuses listed below (Step 1's `blocked-dependency-cycle`, Step 2.0's `blocked-archive-collision`), the mid-run blocked flip (`claimed` → `blocked` when a builder hits a missing external precondition — see Step 8's blocked-flip procedure), and Step 7's early `completed-with-issues` write after a failed remediation (which Step 8 must not overwrite). One terminal status is never written by this action: `cancelled` — a user-directed won't-do decision made via `do-work abandon` (`actions/abandon.md`); the scan treats it like any other terminal status (never claim it).

**Special statuses — these REQs stay in the queue but Step 1 won't pick them up (they're not `pending`, so the "find next pending REQ" scan walks right past them):**
- `pending-answers` — a follow-up REQ whose Open Questions need user input before it can be worked. These accumulate in the queue and get batch-reviewed when the user runs `do-work clarify`.
- `blocked` — waiting on an **external condition** named in `blocked_by` (a service being up, a person answering, credentials provisioned) — not user answers (`pending-answers`) and not another REQ (`depends_on`). Set at capture or by the mid-run blocked flip (Step 8 / `actions/work-reference.md` Failure Classification). Step 1 walks past it, but **re-probes it first** when a `blocked_check` command is present (see the Blocked-condition re-probe paragraph in Step 1); it also unblocks via `do-work clarify` (human-confirmable conditions) or a manual edit back to `pending`.
- `blocked-archive-collision` — set by Step 2.0 when a queue file's REQ id is already archived. Non-destructive holding state; the user flips it back to `pending` (or removes/renames the duplicate) after deciding what to do.
- `blocked-dependency-cycle` — set by Step 1 when a REQ's `depends_on` graph contains a cycle (e.g., REQ-A depends on REQ-B which depends on REQ-A). Non-destructive holding state; the user edits the `depends_on` chain to break the cycle, then flips the status back to `pending`.

## Input

`$ARGUMENTS` may contain:

- **Targeting tokens — specific `REQ-NNN` or `UR-NNN` ids** (e.g., `REQ-042`, `REQ-042 REQ-043`, `UR-011`, or a mix) — process only the resolved REQs and stop (do not process the full queue). This is how actions/pipeline.md scopes work to a specific batch. Token shapes and UR→REQ expansion follow the **Target ID Resolution** contract in `actions/work-reference.md`. **Provenance decides dependency gating:** an **explicitly-named `REQ-NNN`** bypasses `depends_on` (the user named it directly); a REQ reached by **`UR-NNN` expansion** does **not** — it goes through the normal dependency-ready filter, scoped to the UR's member set (naming a batch is a weaker signal than naming each member, and capture wrote those edges expecting them honored). A mixed `do-work run REQ-042 UR-011` is the deduped union, each member keeping its own provenance.
- **`--wave N`** (integer flag, default mode only) — run only REQs at dependency depth N. Roots (no `depends_on`, or all `depends_on` resolve to archived REQs) are depth 0; depth grows by one per dependency layer. Mutually exclusive with **any** targeting token (`REQ-` or `UR-`) — reject the combination with an error.

**Unrecognized arguments are rejected, not ignored.** After stripping `--wave N` and extracting targeting tokens (`REQ-`/`UR-` followed by digits, case-insensitive, per the Target ID Resolution contract), any non-empty token still left in `$ARGUMENTS` is an error. Stop and report:

```
Unrecognized argument(s): <tokens>. Usage: do-work run [REQ-NNN|UR-NNN ...] | do-work run --wave N | do-work run
```

Do **not** fall through to full-queue processing. A leftover token almost always means the user meant to *scope* the run — a typo'd REQ ID (`REG-042` instead of `REQ-042`), or dead muscle memory (a retired mode word) — so silently building the entire queue is the wrong, hard-to-undo default. This generalizes the existing `--wave`-plus-targeting-token rejection to all unrecognized residue; both are parse-time guards.

When `$ARGUMENTS` is empty — no targeting tokens, no flags, no other tokens — process all pending REQs in dependency-aware order (default behavior).

## Steps

**actions/work.md is an orchestrator.** You handle ALL file management (moving files, updating frontmatter, appending sections, archiving). Spawned agents handle implementation work only.

### Step 1: Find Next Request

**Exclusive session.** This pipeline assumes it is the only `do-work` session running against this checkout (`actions/work-reference.md` → **Execution Model — Exclusive Session**). It does not detect or coordinate a concurrent run, and acquires no lock. On session start, **read `do-work/CHECKPOINT.md` first** (Step 10 → **On session start**): it is Crash Recovery's input, not just resume context, so the read is a precondition of the next paragraph rather than a convenience.

**Crash Recovery:** if `do-work/working/` contains any `REQ-*.md`, a claim outlived the run that made it. **A claimed REQ is not automatically this session's to reclaim** — recovery resets frontmatter and strips thirteen generated sections that nothing has committed yet, so it runs only on a REQ the checkpoint records as this session's own interrupted work. Any other claimed REQ is left byte-identical and reported, and is offered for takeover only once it is stale — where a human, never the threshold, authorizes the takeover. Full classification, the staleness threshold and its timestamp guard, and the unattended path: `actions/work-reference.md` → **Crash Recovery (Step 1)**. Once every `working/` file is recovered, taken over, or left alone, proceed with finding the next request.

Glob for `do-work/queue/REQ-*.md`. Sort by number. Read the frontmatter of each (in number order) to check `status`. Don't read the full body at this stage.

**Blocked-condition re-probe.** Determine the probe set by run mode: in **default mode** (empty `$ARGUMENTS`), the set is every REQ whose `status` is `blocked`. In **targeted mode** (`$ARGUMENTS` names targeting tokens), the set is **only the named blocked REQs plus the blocked REQs reached by `UR-NNN` expansion** — do NOT probe the rest of the queue here (that's handled by the Targeted-mode paragraph below). `blocked_check` is verbatim shell that can run up to 30s or touch local services, so a scoped run must never fire probes for REQs the user didn't name. For each REQ in the probe set, check for a non-empty `blocked_check` field. If present, run it as a machine probe **before** categorizing the queue, so a condition that has since become true unblocks the REQ and it participates in *this* run's selection. The `blocked_check` value is repo-authored content (same trust level as the skill's own hook scripts) and is executed **verbatim** — but never interpolate it into a quoted command line (an apostrophe in the condition breaks the quoting and is an injection vector). Instead, write the field value byte-for-byte to a scratch file and run that file:

```bash
# Re-derive paths deterministically — do not carry a variable across blocks.
mkdir -p do-work/working   # may not exist yet — capture.md never pre-creates it, and a queue of only
                           # captured blocked REQs would otherwise have no working/ dir, so the write below
                           # would fail-to-launch and fail-closed would wrongly keep a satisfiable REQ blocked
BLOCKED_CHECK_SCRIPT="do-work/working/.blocked-check-${REQ_ID}.sh"   # REQ_ID is the sanitized REQ id token, e.g. REQ-042
# Write the blocked_check field value to $BLOCKED_CHECK_SCRIPT exactly as read (no quoting, no echo -e), then:
# `timeout` is GNU coreutils — stock macOS ships neither `timeout` nor `gtimeout` (brew coreutils
# adds the latter). Use either when present; otherwise poll the background probe with stock shell
# primitives so every platform keeps the same 30-second bound.
if command -v timeout >/dev/null 2>&1; then
  timeout 30 sh "$BLOCKED_CHECK_SCRIPT"; probe_exit=$?
elif command -v gtimeout >/dev/null 2>&1; then
  gtimeout 30 sh "$BLOCKED_CHECK_SCRIPT"; probe_exit=$?
else
  sh "$BLOCKED_CHECK_SCRIPT" &
  probe_process_id=$!
  probe_wait_ticks=0
  while kill -0 "$probe_process_id" 2>/dev/null && [ "$probe_wait_ticks" -lt 300 ]; do
    sleep 0.1
    probe_wait_ticks=$((probe_wait_ticks + 1))
  done
  if kill -0 "$probe_process_id" 2>/dev/null; then
    kill "$probe_process_id" 2>/dev/null || true
    sleep 0.1
    kill -9 "$probe_process_id" 2>/dev/null || true
    wait "$probe_process_id" 2>/dev/null || true
    probe_exit=124
  else
    wait "$probe_process_id"; probe_exit=$?
  fi
fi
rm -f "$BLOCKED_CHECK_SCRIPT"
```

**Fail closed.** Only `probe_exit == 0` unblocks. Any non-zero exit, a timeout (124), an unreadable/absent field, or a failure to launch means the condition is still unmet — leave the REQ `blocked` and note "probe failed this run" for the exit summary. A probe never halts the work loop and never raises an error.

On exit 0, unblock the REQ: set `status: pending`, stamp `status_changed_at: <now>` (Timestamp rule, `actions/work-reference.md` — the only trace of the flip instant once blocked_at is gone; the board's state timer reads it), **remove `blocked_by` and `blocked_at`** (a stale condition on a runnable REQ would mislead the board and every reader), keep `blocked_check` (harmless while pending, useful if the REQ re-blocks later), and append one history line to a `## Blocked` body section — `- [<date>] blocked on "<condition>" — cleared by probe`. A `blocked` REQ with **no** `blocked_check` is never probed here; it clears only via `do-work clarify` or a manual edit.

**Dependency-aware selection.** For each `pending` REQ, evaluate its `depends_on` field (or its legacy alias `dependencies:` — recognized for back-compat; `depends_on` wins when both present). A REQ is **dependency-ready** when every ID in the resolved dependency list reaches a REQ with `status: completed` or `status: completed-with-issues`. Resolve each dependency ID by globbing `do-work/archive/**/REQ-NNN-*.md`, `do-work/archive/**/REQ-NNN.md`, `do-work/queue/REQ-NNN-*.md`, and `do-work/working/REQ-NNN-*.md`. Cache resolution within a single Step 1 invocation — a 20-REQ queue with 3 deps each is 60 globs; cache hits keep the cost flat. A REQ with unmet dependencies is **dependency-blocked** and is skipped by the scan; it surfaces in the composed exit summary if no other pending REQ is dependency-ready. Process dependency-ready REQs in numeric ID order.

**REQs with neither `depends_on` nor `dependencies:` are roots** and are always dependency-ready. Existing REQs (captured before the field existed) behave exactly as before.

**Cycle detection for `depends_on`.** Before evaluating a REQ's dependencies, walk its `depends_on` graph (or `dependencies:` if `depends_on` is absent — same alias rule) collecting visited IDs into a seen set. If you encounter the current REQ's ID during the walk, the graph contains a cycle — set the REQ's `status` to `blocked-dependency-cycle`, report it, and skip. Mirrors the `addendum_to` cycle-detection approach already used in Step 8 substep 5. Non-destructive — the user breaks the cycle by editing the dependency list and flips status back to `pending`.

**Wave execution (`--wave N`).** If the `--wave N` flag is set, compute each pending REQ's dependency depth before the dependency-ready filter:

- Depth 0: REQs with no dependency list (neither `depends_on` nor the legacy `dependencies:` alias), or whose dependency members are all already archived (completed/completed-with-issues).
- Depth K (K > 0): `max(depth of each dependency member in the current pending set) + 1`.
- A dependency member that is neither archived (completed/completed-with-issues) nor in the current pending set — i.e. it sits in `pending-answers`, `blocked-archive-collision`, `blocked-dependency-cycle`, `claimed`, `failed`, or `cancelled` — contributes depth 0 to this computation. Depth is only about ordering waves; the member's own gating is handled separately by the dependency-ready filter below, which holds the dependent REQ until every member reaches `completed`/`completed-with-issues`.

Filter the pending list to REQs whose depth equals N, then apply the dependency-ready filter normally. If no REQ at depth N is dependency-ready (or none exists at that depth), render the composed exit summary with a leading `No REQs at wave N (depth-N set is empty or fully gated).` line and exit. `--wave` and any targeting token (`REQ-`/`UR-`) are mutually exclusive — reject the combination at parse time with a clear error.

**Explicitly-named REQ ids bypass dependency gating.** When `$ARGUMENTS` names a `REQ-NNN` directly, process it in the given order regardless of `depends_on` (or its `dependencies:` alias). The user named it explicitly. **REQs reached by `UR-NNN` expansion do not bypass** — they go through the dependency-ready filter, scoped to the UR's member set, and run in dependency order (per-token provenance: a mixed run keeps each member's own rule). See the Targeted-mode paragraph below.

**Queue status summary:** After reading all REQ frontmatter, categorize every REQ by status and print a summary before proceeding:

```
Queue: N pending | N completed/done (awaiting archive) | N pending-answers | N blocked | N blocked-archive-collision
```

Count `completed`, `completed-with-issues`, `cancelled`, and `done` statuses together as "completed/done (awaiting archive)." Count `blocked` (external-condition holds) and `blocked-archive-collision` separately so held REQs don't disappear into the silence between "no pending" and "no REQs at all." When the Blocked-condition re-probe unblocked any REQs this run, append `(M probed, K unblocked)` to the `N blocked` figure. If any completed/done REQs exist in `do-work/queue/`, add:

```
⚠ N completed REQs across M URs awaiting archive. Run `do-work cleanup` after this session.
```

**Targeted mode:** If `$ARGUMENTS` contains targeting tokens, first **resolve them per the Target ID Resolution contract** (`actions/work-reference.md`): expand each `UR-NNN` to its member REQs in `do-work/queue/` by scanning `user_request:` frontmatter (never the UR's `requests:` array), dedupe the union with any explicitly-named `REQ-NNN`, and tag each resolved REQ with its provenance — *named* or *UR-expanded*. **A REQ reached both ways — explicitly named *and* a member of a named UR (e.g. `do-work run UR-011 REQ-068`) — is tagged *named*: explicit naming wins.** It is the stronger signal, so the REQ takes the named branch below (bypasses `depends_on`); the more-permissive disposition is the correct one because the user named it outright. **Announce the expansion before the first claim:** a UR argument is the one case where the user cannot see what the run claimed from what they typed, so name the UR and list its resolved REQs in execution order — e.g. `UR-011 → REQ-067, REQ-068 (dependency-ordered).` Without it, a silently-skipped member is indistinguishable from one that was never in the UR. A UR (or a whole argument list) that resolves to zero runnable REQs reports and exits — e.g. `UR-011: no runnable REQs (2 completed, 0 pending).` — and must never fall through to a full-queue run.

Then handle each resolved REQ by its provenance, applying every per-REQ gate exactly as if it had been named directly:

- **Explicitly-named REQ:** verify it exists in `do-work/queue/` and has `status: pending`. Process it in the given order regardless of `depends_on`.
- **UR-expanded REQ:** run it through the dependency-ready filter scoped to the UR's set, in dependency order — it does **not** bypass `depends_on`.
- **`status: blocked` (either provenance):** not claimed on naming alone: run its `blocked_check` probe (per the Blocked-condition re-probe procedure above) — on exit 0 it unblocks to `pending` and is then claimed; on a failing or absent probe, report its `blocked_by` condition and skip it, because naming does not make an unmet external condition true.
- **Missing or any other status:** report the issue and skip it.

Process only the resolved REQs, then stop after the last one completes (skip the loop-or-exit logic in Step 10).

**Default mode (empty `$ARGUMENTS`):** Scan for the first REQ with `status: pending` (skip `pending-answers` — those wait for user input). Reaching default mode requires `$ARGUMENTS` to be genuinely empty — any `REQ-`/`UR-` targeting token routes to targeted mode above, and the unrecognized-argument guard in **Input** has already rejected any other non-flag token, so neither a targeting token nor a fluffed argument (nor a UR that expands to zero runnable REQs) ever silently lands here as a full-queue run.

**Exit paths when no dependency-ready `pending` REQ is found:** render the *composed* exit summary — lead with the dependency-aware headline (`No pending REQs in queue.` when the queue holds no `pending` REQs at all, or `No dependency-ready pending REQs.` when `pending` REQs exist but every one is dependency-blocked), then append every applicable section (completed-awaiting-archive, pending-answers, blocked-on-external-condition, blocked-archive-collision, blocked-by-dependencies) in that order — per `actions/work-reference.md` → **Composed Exit Summary (Step 1)**, then exit the work loop. Only continue past Step 1 when at least one dependency-ready `pending` REQ exists.

**REQ validation:** When reading each REQ's frontmatter, verify it has the required fields (`id`, `status`, `title`). If a REQ file has missing or unparseable frontmatter, skip it and report: `⚠ Skipping [filename]: missing required frontmatter ([field]).` Do not let a single malformed REQ block the entire work loop — skip it and continue to the next.

**Exact glob pattern:** `do-work/queue/REQ-*.md` — if this returns no results, do NOT conclude the queue is empty. Verify by listing `do-work/queue/` contents to rule out a bad pattern.

### Step 2.0: Pre-Claim Archive Collision Check

Before claiming the queue file, verify it isn't a duplicate of an already-archived REQ — rerunning against a file whose twin was archived in a prior run silently re-processes and re-commits it. Run the shipped check:

```bash
<skill-root>/tools/checks/archive-collision.sh REQ-NNN
```

Exit 0 → no collision, proceed to Step 2. Exit 1 (matching archive paths printed) → **bail without moving or claiming**: set the queue file's frontmatter to `status: blocked-archive-collision` (non-destructive — the user flips it back to `pending` after deciding; it also prevents Step 10 → Step 1 livelock), report `REQ-NNN already archived at <path>; remove the duplicate from do-work/queue/ or rename if this is a re-do.`, and continue to the next pending REQ. Never delete the queue file — stale-duplicate vs intentional re-do is the user's call. If the script is missing, glob `do-work/archive/**/REQ-NNN-*.md` and `do-work/archive/**/REQ-NNN.md` (both forms) yourself — same decision rule.

**Scope (minimal):** archive-only; no post-move or pre-commit collision guards (parallel-orchestrator concerns, out of scope).

### Step 2: Claim the Request

1. `mkdir -p do-work/working` and move the REQ file there. If this session is interrupted after the move, the file stays in `working/` and Crash Recovery (Step 1) reclaims it on the next run — every `working/` file is this session's to recover under the exclusive-session model, so no lock or claim record is needed.
2. Update frontmatter: `status: claimed`, `claimed_at: <timestamp>` — the current **UTC** instant `YYYY-MM-DDTHH:MM:SSZ` from `date -u +%Y-%m-%dT%H:%M:%SZ`, exactly like `completed_at` (Timestamp rule, `actions/work-reference.md`). Never local wall-clock time with a `Z` suffix — a future-dated stamp freezes the board's claim stopwatch and flags the card with a clock-skew warning.
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
2. Mark each as `- [~]` with a numbered decision and the builder's reasoning: `- [~] [question] → **D-01**: Builder chose: [choice]. Reasoning: [why]`. An Open-Questions item is a deferred ambiguity, so it is almost always an **ESCALATE** decision under the decide-vs-escalate gate (`crew-members/coding-guardrails.md` § Think Before Coding) — append its **value** and **risk** so they carry into the follow-up the user reviews: `... Reasoning: [why]. Value: [what this choice buys]. Risk: [what breaks if it's wrong, and how reversible]`. When the REQ's `prime_files` cover this area, source the value/risk from the prime's `## Stakes` section rather than re-deriving it.
3. Number decisions sequentially per REQ (D-01, D-02, D-03...). Open Questions decisions and Implementation Decisions (Step 6) share the same D-XX ID space — if Open Questions uses D-01 through D-03, the first implementation decision is D-04. After resolving all `- [ ]` items, append a counter comment immediately after the `## Open Questions` section so Step 6 knows the next available ID: `<!-- D-XX counter: last used D-03. Next decision: D-04. -->` If no decisions were made in this step, write `<!-- D-XX counter: none used. Next decision: D-01. -->` These IDs can be referenced by future REQs.
4. Proceed with implementation using those decisions.

**If a question was escalated to the user mid-run and answered:** a long-running orchestration does legitimately surface a decision between units of work, and that answer lives only in the asking session's context until something writes it down (`crew-members/clear-questions.md` Principle 8). The **orchestrator** (never the builder — all file management is the orchestrator's) writes it into the REQ *before dispatch*, in the **Canonical answered-question format** (`actions/clarify.md`): the item becomes `- [x] [question] → [the user's answer]`, **not** `- [~]`, and carries no D-XX — those record *builder* choices, and recording a user answer as one would leave the stored `Recommended:` rationale standing as the reason and send the settled decision back through `do-work clarify` for confirmation. Append the reasoning too, including anything the answer put out of scope, so a fresh builder reads the decision instead of re-deriving it from `Recommended:`/`Also:`. Any *new* work the answer implies is captured as its own REQ, not left as a sentence in the hand-back. Do not flip `status` here: this REQ is already in flight, so the `pending-answers` → `pending` transition (`actions/clarify.md` Step 5) does not apply.

The follow-up REQs for builder-decided questions are created during **Step 8 (Archive)** — not here. Step 3.5 just records the decisions; the archive step handles the paperwork after the REQ is fully complete.

**Why not block?** Human time is the bottleneck. The optimal windows for user interaction are: (1) capture time, when the user is actively fleshing out requests, and (2) batch-review time, when the user returns to answer accumulated questions. Blocking mid-build wastes builder capacity on idle waiting.

**`pending-answers` REQs:** These accumulate in the queue. When the user returns, they run `do-work clarify` to review all `pending-answers` REQs at once, answer the questions, and flip the status to `pending` so the next work run picks them up. The work loop skips `pending-answers` REQs — it only processes `pending` ones.

If all `- [ ]` items are already `[x]` or `[~]`, or no Open Questions section exists, skip this step entirely.

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

Append validation findings to the `## Plan` section (if any issues found). These are **warnings, not blockers** — the builder can adapt. But flag them visibly so the orchestrator and review step are aware.

**Routes A and B:** Append a skip note (if not already present):

(append the skip note per the **Plan Skip Note — Routes A/B (Step 4)** in `actions/work-reference.md`)

### Step 5: Exploration (Routes B and C)

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

The review step (Step 7) **MUST** run the scope-drift comparison (Routes B and C only): `<skill-root>/tools/checks/scope-drift.sh <req-file>` computes both set-differences (touched-but-undeclared, declared-but-untouched); severity stays your judgment — Important if significant, Minor if trivial like a forgotten import update. Exit 2 means a section is missing (Route A REQs have no Scope declaration — skip the comparison, exactly as the script reports). If the script is missing, compare the two file lists by hand.

### Step 5.75: Pre-Flight Check (Routes B and C)

Quick environment sanity check before the builder starts coding. All checks are **warnings, not blockers.** Append findings to REQ as `## Pre-Flight` section only if issues are found — skip the section entirely if clean.

**Route A:** Skip pre-flight — too lightweight to justify the overhead.

**Routes B and C:** resolve the project's test command first (the prime file's testing section is primary; else `package.json` test scripts, `pytest.ini`, etc. — that resolution is your judgment), then run the shipped check:

```bash
<skill-root>/tools/checks/preflight.sh [test-command ...]
```

It performs the three checks (git clean with `-uall`, test baseline, dependencies present), prints WARN/OK lines, always exits 0, and — when a test command was given — records `do-work/working/baseline.json` + `baseline-failures.txt` so Step 6.5 can separate pre-existing failures from new regressions. A test command that could not be launched at all records `"launched": false` and no failures file — Step 6.5 must refuse to compare against that rather than treat it as a red baseline. Pass the command as separate words or as one quoted string (the quoted form runs via `sh -c`, so `cd app && npm test` works). If the script is missing, run the same three checks by hand (`git status --porcelain --untracked-files=all`; run the test command; check `node_modules`/venv presence).

(append findings per the **Pre-Flight Template (Step 5.75)** in `actions/work-reference.md`, only if issues are found — all checks are warnings, not blockers)

### Step 6: Implementation

**Agent rules loading:** Before spawning the implementation agent, load domain-specific rules:

1. **Always load** `crew-members/general.md` — cross-domain rules and PRIME Files Philosophy
2. **Always load** `crew-members/coding-guardrails.md` — behavioral guardrails (think before coding, simplicity, surgical changes, goal-driven execution)
3. **Conditionally load** `crew-members/[domain].md` — normalize the REQ's `domain` frontmatter per the Schema Read Contract first (e.g., `back-end` → `backend`, `ui_design` → `ui-design`), then load if the resolved domain is set AND the file exists (e.g., `domain: ui-design` → `ui-design.md`). An unknown value after normalization emits the contract's warning and falls back to `general` — no additional domain-specific crew loads (the always-loaded `general.md` from step 1 is the base).
4. **Conditionally load** `crew-members/testing.md` — if the REQ's `tdd` frontmatter normalizes to `true` per the Schema Read Contract (accepts `test_first`/`yes`/`on`/`t` as truthy aliases), or `domain: testing`
4a. **Conditionally load** `crew-members/security.md` — if the REQ's normalized `domain` is `security`, OR if the REQ description references authentication, authorization, session handling, cryptography, secrets handling, input validation/sanitization, or any OWASP-category surface. The "OR" clause is heuristic — when in doubt, load it; the cost of loading a checklist when not needed is low, the cost of skipping it on real security work is high.
5. **Conditionally load** `crew-members/caveman.md` — if the REQ's `caveman` frontmatter normalizes to a non-`false` value per the Schema Read Contract (any of `true`, `lite`, `full`, `ultra`, plus `yes`/`on` → `true`, `light` → `lite`). Compresses agent prose ~65-75% while keeping code and technical terms exact.
5a. **Conditionally load** `crew-members/maintenance.md` — if the REQ's `maintenance` frontmatter normalizes to `true` per the Schema Read Contract. This marks the REQ as a deliberate maintenance pass on the skill's *own* operating instructions (a drifting agent/action/crew/prime file) where removing or narrowing is a candidate fix; it loads the delete-before-you-add discipline **alongside** `coding-guardrails.md`, not instead of it. **Marker-only — do not infer it from the description.** A plain dead-code removal in application source is not a maintenance pass and stays under `coding-guardrails.md`'s implementation-time surgical-changes rule; only the explicit `maintenance: true` marker (set by capture for a removal/narrowing finding on the skill's own instructions) triggers the load. Unlike the security heuristic above, there is deliberately **no** description-based fallback here — a heuristic trigger would misfire on ordinary implementation REQs (which routinely touch adjacent dead code) and load the opposite posture from the one coding-guardrails wants.
6. **If a rules file is missing**, proceed without it — never block on a missing rules file

**Durability (background builder):** When the builder runs as a background or detached sub-agent, follow the durability pattern in `crew-members/background-agents.md` (disk-durable run directory as source of truth; survives a dead orchestrator session).

Spawn a **general-purpose agent** with the loaded rules, any files listed in the `prime_files` array, and context appropriate to the route:

- **Route A**: Request content only — "triaged as simple, aim for a focused minimal change"
- **Route B**: Request + exploration output — "follow existing patterns identified above"
- **Route C**: Request + plan + exploration output — "implement according to the plan"

All routes include these instructions to the agent (pointers — the underlying rules live in the loaded crew-members files and in the REQ frontmatter the orchestrator already wrote):

- **Crew rules govern behavior:** `crew-members/general.md` (always loaded) carries the Prime Files philosophy, Lessons-discipline, test-writing posture, cross-REQ test-break rules, and Discovered-Tasks contract. `crew-members/coding-guardrails.md` (always loaded) enforces think-before-code, surgical scope, and goal-driven execution. Domain/testing/caveman crews layer on top per Step 6's loading order. The builder reads these — do not re-state their contents inline.
- **Prime files come first:** Read every path in `prime_files` before touching code. If the primary utility you are modifying has no prime, investigate and create one (`prime-[name].md`), then update REQ frontmatter. Lessons sections in those primes encode prior mistakes — heed them.
- **P-A-U phasing is mandatory:** Edit the REQ's "AI Execution State (P-A-U Loop)" checkboxes in real time. [PLAN] writes a brief technical approach. [APPLY] stays in declared scope. [UNIFY] runs `git diff --stat`, runs native linters, verifies no debug artifacts, and lists each file checked (the orchestrator audits this in Step 6.3).
- **TDD mode when `tdd: true`:** Follow RED → GREEN → REFACTOR. Anchor RED on the REQ's `## Red-Green Proof` section if present. Report the red-green evidence (test name, failure-before, pass-after) — Step 6.5 verifies it.
- **Captured proof first:** If `## Red-Green Proof` is present, its RED prompt/case and GREEN outcome are the primary behavior tests must prove. Only adapt with documented reason.
- **Log Decisions as D-XX:** Significant implementation choices not dictated by plan/requirements become numbered entries in a `## Decisions` section. Continue numbering from the `<!-- D-XX counter: ... -->` comment Step 3.5 left behind; if none, start at D-01. Each decision needs reasoning — without it, the intent trail breaks. Sort each by the decide-vs-escalate gate (`crew-members/coding-guardrails.md` § Think Before Coding): a reversible, low-reach choice is **DECIDE & STATE** (reasoning only — it surfaces later as a *handled* item); a choice that's irreversible/expensive, taste-dependent, or genuinely contestable is **ESCALATE** — add `Value:` and `Risk:` lines so the hand-back can surface them.
- **Write only inside the declared scope:** the REQ's `## Scope` "Files I will touch" list (mirrored to `write_set` at Step 5.5) is the builder's write boundary. Discovering mid-build that you need a file outside it is a **stop-and-report to the orchestrator, never a silent write.** The orchestrator records the request and its resolution in the REQ trail as a `## Decisions` D-XX entry (it is a scope judgment) and extends both the Scope list and `write_set`. An **absent or empty `write_set` is not a write prohibition** — but it is not a stated boundary either, so it is not licence: it means this REQ never ran Step 5.5 (a Route A REQ, or one whose `## Scope` was stripped by crash recovery), so derive the boundary from the REQ's own text and keep the change inside it — do not read the empty field as full-scope freedom.
- **In worktree dispatch mode, never write the main tree:** commit on your own branch and hand back the manifest — the orchestrator is the sole integrator and merges. A shared file that needs one line of wiring is an **integration seam**: hand back the exact line and where it goes rather than editing the shared file yourself (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**).
- **Out-of-scope finds go to `## Discovered Tasks`** (a separate section, not nested inside Implementation Summary) — do not fix inline. Step 8 classifies and queues them.
- **Report back the file manifest:** list every source file created/modified/deleted with the action verb, plus tests touched. The orchestrator writes the formal `## Implementation Summary` from your report.
- **Standard freedoms and obligations:** Full file/shell access. Escalate to explore or plan if the work proves harder than triaged. Document blockers explicitly. Identify and run related existing tests; honor any test-command map in the prime file (takes precedence over generic detection).

**Hand-back merge (worktree dispatch mode only — the orchestrator's job, not the builder's).** When the builder returns its manifest, integrate here, at the end of Step 6, **before** Step 6.25: every evidence step from 6.25 onward reads the merged tree, so a merge deferred past that point leaves qualify and review with nothing to check. On the integration branch:

1. Run `git rev-parse --short HEAD` and hold the printed hash as **`<pre>`** — this REQ's pre-merge integration tip and the lower bound of its merge range. Capture it once per REQ; a remediation re-merge keeps the first one.
2. `git merge --no-ff --no-commit <operative_name>` — the branch the builder was actually dispatched on, which is the collision variant where there was one and the derived `worktree-agent-REQ-NNN-<suffix>` otherwise; never re-derive it here. Resolve any conflict, apply the builder's handed-back integration-seam lines, then `git commit` — folding the seam into the merge commit is what puts it inside the merge range.
3. Run `git rev-parse --short HEAD` again and hold that hash as **`<merge_hash>`** — the upper bound of the range, and the hash Step 9 writes into `commit:`.

Hold both hashes — and `<operative_name>` — as literals you re-type into each later command (shell variables do not survive between command blocks), and pass the range `<pre>..<merge_hash>` to Steps 6.3, 7, 8, and 9. Full sequence, the `Already up to date.` guard, and the cumulative range across remediation re-merges: `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**. Serial mode has no merge and no range — skip this entirely.

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
- **Design-artifact exception:** For `domain: ui-design` requests that produce design deliverables rather than code (wireframes, IA specs, visual specs, interaction specs), the artifact files themselves count as project files. Place them in the project's design docs directory (e.g., `docs/design/`) — not inside `do-work/`. The Implementation Summary lists these files normally.

### Step 6.3: Qualify Implementation

After the builder returns and the Implementation Summary is written, the **orchestrator** (not the builder) independently verifies the builder's claims before proceeding. This is not self-reporting — the orchestrator reads actual output, not the builder's description of it.

**Mechanical checks (run the shipped script):**

```bash
<skill-root>/tools/checks/qualify.sh <req-file>
```

**In worktree dispatch mode** the builder's work is already committed and merged (end of Step 6), so the working tree is clean — pass this REQ's merge range so the mechanical checks read the merged diff instead of an empty working diff:

```bash
DO_WORK_DIFF_RANGE="<pre>..<merge_hash>" <skill-root>/tools/checks/qualify.sh <req-file>
```

where `<pre>` and `<merge_hash>` are the two hashes held since Step 6's hand-back merge — re-type both as literals (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**). The script **FAILs** (exit 1, naming the range) if the range does not resolve, rather than reading an empty diff and passing vacuously — so a `..` from a lost endpoint surfaces here as a qualification failure. Serial mode omits the variable and reads the working/staged diff.

It verifies checklist items **1 (files exist / show in diff)**, **4 (P-A-U box audit + debug artifacts in the diff)**, and the grep half of **5 (wiring)** — plus Step 6.25's "only `do-work/` paths ⇒ not implemented" rule. FAIL lines are qualification failures; WARN lines are evidence handed to your judgment — in particular, an unreferenced `(new)` file is only dead code if it isn't an **exception**: entry points, config files, test files, standalone scripts, framework-convention files discovered by file-system routing (Next.js `pages/`/`app/`, SvelteKit/Remix `routes/`, Nuxt/Astro `pages/`), barrel re-exports, side-effect-only imports (CSS modules, polyfills), and dynamic-import-only files that static grep can't see. If the script is missing, run items 1/4/5 by hand per its header comment.

**Judgment checks (yours, not the script's):**

2. **Changes are substantive:** For each `(new)` file, verify it is not a placeholder (more than boilerplate/empty exports/TODO comments — minimum 10 meaningful lines for source files, 3 for config). For `(modified)` files, verify the diff contains changes related to the REQ's requirements, not just whitespace or import shuffling.
3. **Requirements traced:** Re-read the REQ's What/Detailed Requirements section. For each stated requirement, confirm at least one file in the Implementation Summary plausibly addresses it (by filename and diff content). Flag any requirement with no corresponding file change.
6. **Flowing:** For files that handle data (API endpoints, data stores, handlers, services), verify the data path isn't hardcoded or stubbed. Check for: hardcoded empty arrays `return []`, placeholder strings like `"TODO"` or `"placeholder"`, `return null` in data-fetching functions, commented-out database calls. If found, flag as hollow implementation — the file exists and is wired but doesn't actually do anything.

**Anti-rationalization rules** (apply when evaluating the above):

Apply the qualification anti-rationalization table in `actions/work-reference.md` → **Qualification Anti-Rationalization Table (Step 6.3)** (e.g., "the summary says files changed" → check the file system; "the builder checked UNIFY" → read the diff for debug artifacts).

**If qualification fails on any check:**
1. Append a `## Qualification` section to the REQ noting what failed and why.
2. Return to Step 6 — spawn the builder again with the specific failures as context.
3. Maximum **2 re-qualification attempts**. After that, note remaining issues and proceed to Testing (Step 6.5). The review step will catch what remains.

**If qualification passes:**
- Append a brief `## Qualification` section: "Passed — [N] files verified, [N] requirements traced, P-A-U confirmed."
- Proceed to Step 6.5.

### Step 6.5: Testing

Before marking complete, verify tests pass:

1. **Check the prime file for test guidance** — if the REQ's `prime_files` reference a prime with a testing section (test commands, code-area-to-test mappings), use that as the primary source for what to run. **Before running, verify each listed command still exists**: for npm scripts check it's present in `package.json`; for other tools verify the config file exists (`jest.config.*`, `pytest.ini`, `Cargo.toml`, etc.). If a prime test command is no longer valid, fall back to generic detection for that command and note: `Prime test command '[cmd]' not found — falling back to generic detection.` Prime test maps are project-specific knowledge that generic detection can't replicate (e.g., "changes to `lib/inpainting.js` require `npm run test:api`" or "`npm test` is always safe but `npm run test:e2e` costs money").
2. **Fall back to generic detection for unmapped files** — if the prime has no testing section, or if you changed files the prime's test map doesn't cover, fall back to generic detection for those files: look for `package.json` test scripts, `jest.config.*`, `pytest.ini`, `Cargo.toml`, `*_test.go`, etc. A partial prime map is not an excuse to skip tests — matched files use the prime's commands, unmatched files use generic detection. If neither source yields test commands for a file, skip testing for it and note it.
3. **Run relevant tests** — target tests related to changed code, not the full suite (unless it's fast). If the prime specifies different commands for different code areas, run only the commands relevant to the files you changed. For unmapped files, run whatever generic detection found.
4. **If tests fail** — check whether the failures were already recorded as baseline failures in Step 5.75 (Pre-Flight); if `do-work/working/baseline.json` / `baseline-failures.txt` exist (written by `tools/checks/preflight.sh`), compare against those records mechanically. **First read `baseline.json`'s `launched` field: `launched: false` means the test command never ran, so there is no baseline — exclude nothing, treat every failing test as a candidate regression, and say so rather than comparing against a record that describes no test run.** Otherwise: if a failing test matches a pre-existing baseline failure (same test name/file, same failure mode), exclude it from the pass/fail gate — the builder should not be blamed for pre-existing failures. Only **new regressions** (tests that passed at baseline but fail after implementation) require fixing. Return to implementation to fix new regressions. On attempt 2+, load `crew-members/debugging.md` and `crew-members/testing.md` for the builder to follow the structured debugging methodology and review test quality. Loop until passing or mark as failed after 3 attempts.
5. **If new tests are needed** — spawn a general-purpose agent to write them following existing patterns, then run them.

Append to the request file:

(append per the **Testing Section Template (Step 6.5)** in `actions/work-reference.md`; omit Red-green validation for non-behavioral changes, and trace it back to `## Red-Green Proof` when present)

Omit `Red-green validation` if no request-specific tests were written or identified, or if the change is non-behavioral (refactor, config, docs, cleanup) — use regression evidence instead. Omit `Existing tests updated` if no prior tests were modified.

When the REQ includes `## Red-Green Proof`, the `Red-green validation` entries should trace back to that captured RED/GREEN pair. If the implemented test uses a nearby equivalent instead of the exact captured prompt/case, explain why.

**TDD verification:** If the REQ has `tdd: true`, the `Red-green validation` section is mandatory — the builder must show test-first evidence that they used RED/GREEN TDD (test written before implementation, failed, then passed after). If this evidence is missing, treat it as a test failure: return to implementation (same path as step 4 above) with explicit instructions to provide red/green evidence — write the failing test first, confirm it fails, then make it pass.

### Step 7: Review

Run actions/review-work.md in **pipeline mode** against this REQ.

The review reads the REQ (in `do-work/working/`), the original UR, and the current diff (`git diff` or `git diff --staged`) to evaluate the implementation: requirements check (did we build what was asked?), code review (is it solid?), and acceptance testing (does it actually work?). **In worktree dispatch mode** the working tree is clean after the merge, so the review reads this REQ's merge range `<pre>..<merge_hash>` instead (`actions/review-work.md` Step 4, Get the Diff).

**Restatement sweep (MUST).** If this REQ's diff redefines something other text restates — a contract token, a schema field's semantics, a gate's wording, a prescribed command's output shape — the review runs the sweep defined in `actions/review-work.md` Step 6 (**Restatement Sweep**) and reports every stale restatement as a finding, including ones in files outside this REQ's declared Scope.

**How to run it:** Spawn an agent with actions/review-work.md file, the REQ path, and the `crew-members/[domain].md` file (normalize `domain` per the Schema Read Contract first; if the resolved domain has a matching file, load it; otherwise skip) — in worktree dispatch mode, also pass this REQ's merge range `<pre>..<merge_hash>` so the review reads the merged diff rather than the clean working tree. Or read actions/review-work.md file and follow its pipeline mode instructions in the current session.

**What happens next depends on the review result:**

- **Acceptance = Pass AND overall ≥ 75%**: Append the Review section to the REQ and continue to archive as `completed`. Minor findings go in the report only.
- **Acceptance = Partial OR overall 50-74%**: Append Review, continue to archive as `completed`, but the review **MUST** create follow-up REQs for every Important finding. These are not optional — they enter the queue and block the UR from being considered "done" until addressed.
- **Acceptance = Fail OR overall < 50%**: **Do NOT archive as completed.** Instead:
  1. Append the Review section to the REQ.
  2. Return to Step 6 (Implementation) with the review findings as context for the builder. Load `crew-members/debugging.md` for the remediation attempt — the builder needs structured debugging methodology, not just "try again."
  3. The builder gets **ONE remediation attempt**.
  4. Re-run Steps 6.25 through 7 (Summary → Qualification → Testing → Review) on the remediated code.
  5. If still failing after remediation: update frontmatter to `status: completed-with-issues`, `completed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`), append a `## Remediation` section documenting both attempts, and create follow-up REQs for all remaining Important findings. Then proceed to archive (Step 8) — the frontmatter is already set, so Step 8 should not overwrite it.

The status `completed-with-issues` means the REQ was archived but has known unresolved problems. It counts toward UR completion for archiving purposes, but the follow-up REQs must be processed before the work is considered ship-ready. This status is visible in the recap and present-work actions.

**Follow-up REQs are created based on finding severity, not score.** The review creates follow-up REQs for each **Important** finding (regardless of overall score). Minor and Nit findings go in the report only. The follow-up REQs enter the queue and get processed in a future loop iteration. Follow-up REQs created by the review step must include: `status: pending`, `user_request: [same UR as the reviewed REQ]`, `addendum_to: [reviewed REQ id]`, `domain: [same domain]`, and `review_generated: true`. Place them in `do-work/queue/`. Cycle detection (Step 8, substep 5) applies to these follow-ups — check the `addendum_to` chain before creating.

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

**Write the `## Orientation` block (the hand-back's "what's being built"):** After the Lessons Learned section, append a short `## Orientation` section reporting the change at feature/subsystem altitude — "Now you can X; lives in Y subsystem" — not a file list. Use the REQ's `prime_files` to name the subsystem; flag `[MAP CHANGED]` only when the change alters the system's shape (new module, data flow, contract, or a renamed concept). Run a narrowed staleness spot-check on each touched prime per `actions/prime.md` Step 2 / Step 6 (do its referenced paths still exist?) and flag any prime the change made stale. **Scale to reach:** a leaf REQ is one line; a map-changing REQ gets a short paragraph and a why-it-matters. When `prime_files` is empty, derive a one-line feature-altitude summary from the What / Implementation Summary instead — never a file list. This block feeds the **WHAT'S BEING BUILT** section of the Decision Brief (`actions/work-reference.md` → **Decision Brief (hand-back format)**). Crash recovery strips `## Orientation` on re-queue (it's orchestrator-generated).

**Update prime files (deferred to Step 8):** After writing the Lessons Learned section, check the REQ's `prime_files` frontmatter. For each listed prime file relevant to this lesson, **collect a pending prime-link write** — do NOT execute the write here. The REQ is still in `do-work/working/`, so any link pointing to its eventual archive location would either be broken or tempt a link to the transient working path.

Record each pending write as a tuple: `{ primeFilePath, relativeLinkText, lessonSummary }`. Hold them in memory (or a small scratch file under `do-work/working/`) until Step 8.

Compute each deferred prime-link path relative to the prime file's location (not the repo root) per `actions/work-reference.md` → **Deferred Prime-Link Path Computation (Step 7.5)**; the existence-verify on the resolved path runs in Step 8 (post-move), which is why the write is deferred.

Only add a link when the lesson is relevant to that prime file's scope — don't spray every lesson into every prime file. If the REQ has no `prime_files` or the lessons aren't relevant to any prime file, skip this and clear the pending list.

**Knowledge-base handoff.** After the Lessons Learned section is written and prime-file links are in place, follow `actions/kb-lessons-handoff.md` to offer dropping a structured source document into `kb/raw/inbox/` so the next `bkb triage` + `bkb ingest` cycle compiles the lessons into the wiki. The handoff asks the user before writing and records `kb_status` (plus `kb_entry` on success) back onto the REQ. In unattended pipeline runs with no human in the loop, the handoff defaults to `kb_status: pending` — it never writes to the KB without consent. If the project has no `kb/` directory, the handoff points the user at `do-work bkb init` and defers; it never blocks archival.

### Step 8: Archive

**On success:**

**Post-merge verification (worktree dispatch mode only).** This REQ's work was merged at hand-back — end of Step 6, before Step 6.25 (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**) — so nobody has verified the merged state; the builder's checks ran on its own branch. Re-run the REQ's acceptance checks against the merged tree **before** the substeps below archive anything; the unit you verify is the unit you roll back. A red merged state stops the archive: revert to the last verified state and re-dispatch, never archive-and-follow-up. Procedure: `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**.

1. Update frontmatter: if the current status is already `completed-with-issues` (set by Step 7 after a failed remediation), preserve `completed-with-issues` and ensure `completed_at: <timestamp>` is present. Otherwise set `status: completed`, `completed_at: <timestamp>`. **`completed_at` (current UTC instant — `date -u +%Y-%m-%dT%H:%M:%SZ`, per the Timestamp rule in `actions/work-reference.md`) is mandatory on every terminal flip — never skip the stamp.** It and the `commit:` hash (written back in the Commit Phase) are the only sources the board resolves a completion instant from; a terminal REQ with neither surfaces as a completion anomaly on `do-work board` (see `actions/work-reference.md`'s Full Frontmatter stamping rule).
2. Verify `## Implementation Summary` is present (written in Step 6.25). If missing, append it now — this should not happen in normal flow, but crash recovery may skip it.
3. **Create follow-ups for builder-decided questions:** If the REQ has any `- [~]` items in Open Questions where the builder's choice affects what the user sees or interacts with, create a follow-up REQ for each. **Create follow-ups for:** UX decisions (interaction behavior, visibility, layout), scope boundaries (what's included/excluded), data representation choices. **Skip follow-ups for:** purely technical decisions (caching strategy, algorithm choice, internal naming, DB indexes) that don't change user-facing behavior.

   Create each follow-up per the **Builder-Decided Follow-up Template (Step 8)** in `actions/work-reference.md`; these go in `do-work/queue/` with `status: pending-answers`, and the user reviews them via `do-work clarify`.

   **Whenever authoring Open Questions text a user will answer via clarify** — here, the intent-failure follow-ups in the Failure Classification table, or any other `pending-answers` REQ — load `crew-members/clear-questions.md` and write for a cold reader: gloss every coined label or spec §-reference, and state why the decision is the user's rather than yours (Principle 7). You have the whole spec in your head right now; the reader answering in a later clarify session has none of it. Don't rely on clarify to repair density at presentation time.
4. **Queue Discovered Tasks:** Check the REQ file for a `## Discovered Tasks` section (appended by the implementation agent as a separate section — not inside `## Implementation Summary`). For every item listed, classify by severity and create follow-up REQs accordingly.

   Classify each by severity and queue follow-ups per `actions/work-reference.md` → **Discovered Tasks Classification (Step 8)**: `[critical]` → `status: pending`, auto-queued + prominent report; `[normal]`/`[low]` → `status: pending-answers` via the Open-Questions consent flow — except test-only mechanical-hygiene discoveries meeting that section's carve-out (all three bullets), which auto-queue as `status: pending` with an auto-approved note and a `↺` report line.
5. **Cycle detection:** Before creating any follow-up REQ, verify the current REQ's own `addendum_to` chain is not already circular. Algorithm: walk `addendum_to` links (honoring the `amends`/`parent`/`amendment_to` alias per the Schema Read Contract when the canonical key is absent) starting from the current REQ, collecting each visited ID into a seen set. If you encounter the current REQ's ID again during the walk, the chain is already circular — do not create any follow-ups. Report: `⚠ Cycle detected in addendum_to chain: REQ-NNN → REQ-MMM → ... → REQ-NNN. Skipping follow-up — manual resolution needed.` This handles chains of any length.
6. Archive based on REQ type. The physical move out of `working/` is the archive; no lock or claim record is updated (the exclusive-session model keeps none).

| REQ has... | Archive behavior |
|------------|-----------------|
| `user_request: UR-NNN` | Check if ALL REQs in the UR are **terminally resolved** — see `actions/work-reference.md`'s Schema Read Contract → Terminal-resolved status set; that set is canonical, don't restate or fork it here. Any status outside it holds the UR open, **`failed` included** (how a `failed` REQ is resolved so it leaves this held-open state is defined at that canonical statement — do not re-derive it here). Check `do-work/queue/`, `do-work/working/`, `do-work/archive/` root, and `do-work/archive/UR-NNN/` for REQs belonging to this UR. If all resolved: move them into the UR folder (`cancelled` REQs are resolved-by-decision and consolidate like completed ones), move entire UR folder to `archive/`. If any REQ is **not** terminally resolved — e.g. `failed`, `pending`, `pending-answers`, `claimed`, `blocked` (waiting on an external condition), or a `blocked-*` holding status: move this REQ to `archive/` root; UR stays in `user-requests/` until every REQ resolves. |
| `context_ref` (legacy) | Move REQ to `archive/`. If all related REQs are now archived, move the CONTEXT doc too. |
| Neither (standalone legacy) | Move directly to `archive/`. |

7. **Execute deferred prime-link writes (from Step 7.5):** Now that the REQ is at its final archive path, walk the `pendingPrimeLinkWrites` collected during Step 7.5. For each pending entry:
   - **Check for the inline-only marker first.** Look at the prime file's `## Lessons` section (if it already exists). If it opens with an HTML comment containing the phrase "inlined, not linked" (the pattern is `<!-- Lessons are inlined, not linked: ... -->` — see `tools/queue-kanban/prime-do-kanban.md`'s `## Lessons` header for the exact wording), the prime has declared itself inline-only: skip the link steps below and instead append a plain bullet with the lesson summary — `- REQ-NNN: 1-line summary of the lesson` (no link, matching the prime's existing inlined entries) — then continue to the next pending entry.
   - Otherwise (no marker present), proceed as before:
     - Compute the relative path from the prime file to the REQ's actual archived location (UR folder if the UR was just consolidated, or `archive/` root if the UR is incomplete).
     - Verify the resolved path points to an existing file. If it doesn't, report the broken link and skip — do NOT silently write a broken link.
     - Append the link to a `## Lessons` section in the prime file (create the section if it doesn't exist):
       ```markdown
       ## Lessons

       - [REQ-NNN: 1-line summary of the lesson](<relative-path-to-archived-req>#lessons-learned)
       ```
   - Stage the prime file along with the implementation files in Step 9.

   This is the post-move execution that makes the existence-verify meaningful — Step 7.5 only collected; the writes happen here.

8. **Worktree cleanup (worktree dispatch mode only):** now that the REQ file has reached its final path, remove the builder's leftovers **by this REQ's operative name** — the name the worktree was actually created with (the collision variant where there was one), held since dispatch and never re-derived from the slug: `git worktree remove <path>` (no `--force`, `<path>` being the worktree whose basename is `<operative_name>`), then `git branch -d <operative_name>` **run from the integration branch you merged into** (`-d` tests merged-ness against the current HEAD, so from anywhere else the refusal means the wrong branch, not an unmerged one), then `git worktree prune`. Never `-D` or `--force`: a refusal means a dirty worktree or a merge that was skipped or lost — stop and report, don't force away the evidence (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**).

**On failure:**

Classify the failure and queue the right follow-up per `actions/work-reference.md` → **Failure Classification (Step 8)**. Run the **upstream-failure short-circuit first** (if any `addendum_to`/`depends_on` ancestor is `failed`, short-circuit to `error_type: spec` with an upstream-cascade error), then fall through to the Intent/Spec/Code/Environment symptom table. Set `status: failed`, `completed_at: <timestamp>` (mandatory on every terminal flip, same stamping rule as success), `error`, `error_type`; create the follow-up (Intent/Spec/Code) with `addendum_to` chained and the original dependency list preserved; move to `archive/` root.

**Mid-run blocked flip (external precondition):** Before classifying an Environment failure as terminal, apply this test — it is the non-terminal alternative to `error_type: environment` for a precondition that will simply become true later:

- **Both must hold to flip:** (1) *No substantive implementation edits landed this attempt* — the orchestrator confirms via `git status --porcelain -- . ':(exclude)do-work/'` / `git diff -- . ':(exclude)do-work/'` that the builder made no repo changes for this REQ **outside `do-work/`**. The exclusion is load-bearing: the REQ's own bookkeeping (the move to `working/`, appended Triage/Plan sections) is always dirty mid-run, so an unscoped porcelain check can never read clean and would silently defeat this flip (triage/plan/explore or the first implementation action failed on the missing dependency with an otherwise-clean tree). (2) The missing thing is an **external precondition expected to become available on its own** — a service coming up (LM Studio, a DB), a person answering (a designer's mockup), credentials getting provisioned — not a broken toolchain or a permission the user must repair, and not a transient crash (retry those in-loop first, then classify normally).
- **In worktree dispatch mode, clause (1) is judged from the builder's branch — the main tree cannot answer it.** The builder commits in its own worktree and never writes the main tree (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**), so the porcelain/diff check above reads clean however much work landed; taken at face value it would flip every worktree failure to `blocked`, including ones that already dropped a pile of committed work. Clause (2) is unchanged; substitute this evidence for clause (1) only, picking the case by whether Step 6's hand-back reached a merge:
  - **A hand-back merge completed for this REQ** (you are holding a `<merge_hash>`) — edits landed: do not flip. A merge commit can only exist here if the builder committed something, because the hand-back stops on `Already up to date.` rather than fabricating one, and `--no-ff` produces a merge commit even where the branch could fast-forward.
  - **No merge completed** (the failure preceded hand-back, or the hand-back was empty) — probe for the branch before counting anything: if `git rev-parse --verify -q '<operative_name>'` fails, the branch was never created (dispatch did not get that far) and nothing landed, so flip to `blocked`. Do not reach for the count here — `rev-list` on a missing branch exits fatal and prints no number at all, so it cannot decide the very case the flip exists for. Only once the branch resolves: edits landed if and only if `git rev-list --count HEAD..<operative_name>` is greater than zero — `HEAD` because the orchestrator runs this from the integration branch it dispatched from — meaning the builder's branch carries commits the integration branch does not contain (integration is by merge and never rebase, so they stay recognizable as the builder's). A count of `0` is the other genuine nothing-happened-this-attempt case and still flips to `blocked`.

  Judge from the branch, not from the handed-back manifest: the manifest is the builder's claim about its own work, and the orchestrator reads actual git state rather than the builder's description of it (same stance as Step 6.3). Uncommitted edits sitting in the builder's worktree do not count as landed — the main tree is pristine, so a re-dispatch after the block starts clean, and the stray worktree is swept by `actions/work-reference.md` → **Crash Recovery (Step 1)** or `actions/cleanup.md` → **Pass 5: Orphaned Worktrees (consent-gated)**. Serial mode ignores this bullet entirely and uses the working-tree check above.
- **If both hold**, do NOT fail. The orchestrator (never the builder — all file management is the orchestrator's) sets `status: blocked`, `blocked_by: "<condition>"`, `blocked_at: <now>` (current UTC instant — Timestamp rule, `actions/work-reference.md`); removes `claimed_at` and `route`; appends a `## Blocked` section recording what's missing, how it was discovered, and — only if the user supplied or confirmed one — a `blocked_check:` probe command; then moves the file **back to `do-work/queue/`** (it is a hold, not an archive), reports `[REQ-NNN] blocked on: <condition> — released, continuing`, and continues to the next REQ. The REQ re-enters selection on a future run via its `blocked_check` probe, `do-work clarify`, or a manual edit.
- **If either fails** (real edits already landed, environment the user must fix, or retries exhausted), fall through to the Environment classification above and archive as `failed` with `error_type: environment`.

### Step 9: Commit Phase (Git repos only)

> **Named entry point.** Other actions reference this as **work.md's Commit Phase** (not by step number) — e.g. `actions/commit.md` and `actions/review-work.md`. The `9` is for internal navigation only; callers must use the phase name so they don't break if steps are renumbered.

Check for git with `git rev-parse --git-dir 2>/dev/null`. If not a git repo, skip.

Before committing a successful REQ, write a changelog entry in the target repo's root `CHANGELOG.md` per `actions/work-reference.md` → **Changelog Entry Procedure (Step 9)** — create the file if it's missing, match the repo's existing format if it has one.

House-format entries are keyed `## X.Y.Z — [Short Descriptive Title] (YYYY-MM-DD)`; the version is bumped by change type from the repo's own version file (which gets bumped and staged too), its release tags, or — for an unversioned repo — the changelog's own counter. Then: one commit per request, format `[{id}] {title} (Route {route})` + `Implements:` line + summary bullets. Stage only the explicit files (implementation files from the Implementation Summary, the archived REQ, the `CHANGELOG.md` entry, the version file it bumped, any follow-up REQs, UR-folder moves, and any prime files touched in Step 8 substep 7) — see `## Rules` below for the staging/hook guard. Validate the staged file list against the Implementation Summary (successful REQs only). After the commit, write the real short hash back into the archived REQ's `commit:` field with `<skill-root>/tools/checks/record-commit-hash.sh <req-file> <hash>` — it guards the one-line frontmatter edit and refuses to write when the rewrite would change anything else — then record it in a **separate metadata commit** (never amend), and confirm what landed with the script's `--verify` mode. A `FAIL:` line means stop, not retry: free-form edits at this step truncated six archived REQs to 0 bytes in a consumer repo. Full bash + metadata-commit procedure: `actions/work-reference.md` → **Commit & Metadata-Commit Procedure (Step 9)**.

**Optional accelerator for the version bump and the two hand-checked invariants.** Decide the bump size yourself (the procedure's Bump-size table — the tool must not infer it), then:

```bash
(cd <skill-root>/tools/queue-kanban && go build -o queue-kanban .) 2>/dev/null \
  && <skill-root>/tools/queue-kanban/queue-kanban next-version <patch|minor|major> --repo-root <project-root>
```

It bumps the repo's `**Current version**: X.Y.Z` line, reads it back to confirm the write landed, and prints the new number — and it **never writes `CHANGELOG.md`**, which stays a human-authored write. It exits non-zero without writing when `go` is absent or the repo keeps no such line (the common case outside this skill's own checkout): **fall back to the Changelog Entry Procedure's own version-source resolution**, which handles `package.json`, tags, and unversioned repos. After writing the entry, `queue-kanban verify --repo-root <project-root>` checks the version-vs-changelog agreement and the title-reuse rule mechanically — the two checks the procedure otherwise asks you to perform by eye.

**In worktree dispatch mode the implementation commit already exists** — the builder committed on its branch and Step 6's `--no-ff` merge integrated it — so Step 9 does **not** stage implementation files. Stage only the changelog entry, the version file it bumped, the archived REQ, any follow-up REQs, UR-folder moves, and any prime files touched in Step 8 substep 7, then commit those. Write **`<merge_hash>`** — the `--no-ff` merge commit captured in Step 6, the latest one if remediation re-merged, and **never** `git rev-parse --short HEAD`, which here names this changelog commit — into the archived REQ's `commit:` field (the merge commit is the implementation's provenance record), passing it as the second argument to the same `record-commit-hash.sh` guard the serial path uses; the guards are identical on both paths, and the script reports a merge-commit hash back so consumers know to read it with `git show --first-parent -m`. Validate the Implementation Summary file list against `git diff --name-only <pre>..<merge_hash>` (the merge range), not the stage — the implementation files are in the merge, not in this commit's staged set.

### Step 10: Loop or Exit

Re-check `do-work/queue/` for `REQ-*.md` files (fresh check, not cached).

- **Dependency-ready `pending` REQs found**: **CONTEXT WIPE** (see below), then loop to Step 1.
- **No dependency-ready `pending` REQs remain** (queue may still have dependency-blocked or held REQs): Write a **Session Checkpoint** (see below), run actions/cleanup.md, then report the final summary using the **same composed structure** as Step 1's "Exit paths when no `pending` REQs found" — render the completed/done section, the pending-answers section, the blocked-on-external-condition section, the blocked-archive-collision section, and the blocked-by-dependencies section in that order, including only those that have at least one REQ. If none of the five sections applies (queue is fully empty), report completion and exit. Mixed cases render all applicable sections in one summary.

#### Context Wipe — Verified

Before looping to Step 1 for the next REQ:

1. **Fresh agents:** Spawn a NEW agent for the next REQ. Do not reuse the previous builder/explorer/planner agent. Each REQ gets clean agents with no carried-over context.
2. **Explicit declaration:** State in your progress message: `Context wipe: previous REQ was [REQ-NNN] working on [files]. Now starting fresh for next REQ.`
3. **Contamination check:** When the next REQ's builder returns its Implementation Summary (Step 6.25), compare the file list against the *previous* REQ's Implementation Summary. Unexpected overlap — files from the previous REQ appearing without an explicit `addendum_to` or `related` link — is a scope contamination signal. Flag it in the Qualification step (Step 6.3).

#### Session Checkpoint

At the end of every work session (whether all REQs completed, user stops, or session is ending), write `do-work/CHECKPOINT.md`. Scale the checkpoint to how much happened:

(write `do-work/CHECKPOINT.md` per the **Session Checkpoint Template (Step 10)** in `actions/work-reference.md`, scaled to session depth: light / moderate / heavy)

**Session depth guide:**
- **light** (1-2 REQs): Minimal checkpoint — Completed + Still Queued sections are sufficient
- **moderate** (3-5 REQs): Add Session Notes with patterns observed and environment quirks
- **heavy** (6+ REQs): Add Context Summary recapping key decisions and recommending the next session re-read prime files fresh rather than trusting carried-over assumptions

**On session start (Step 1 addition):** Before crash recovery, check for `do-work/CHECKPOINT.md`. If it exists:
1. Read it and report a brief summary: `Resuming from previous session. Last completed: REQ-NNN. [N] REQs still queued.`
2. Its `## In Progress (interrupted)` section is what crash recovery classifies against — a `working/` REQ named there is this session's own to recover, and one that isn't is a foreign claim recovery must not strip.
3. **Do not delete yet.** Keep the checkpoint until crash recovery has finished with every `working/` file, then delete it. Deleting only after recovery is done still prevents losing resume context if the session crashes again mid-recovery.

This is NOT a blocking gate. With no checkpoint, the session starts normally — crash recovery still runs, it just has nothing to match against, so every claimed `working/` REQ is treated as a foreign claim and left intact.

## Clarify Questions

The clarify workflow has its own action. Run `do-work clarify` — it handles batch-review of `pending-answers` REQs, where the user confirms, overrides, or discards builder decisions. Resolved REQs flip back to `pending` and re-enter the work queue.

## Orchestrator Checklist (per request)

```
□ Step 1: Read CHECKPOINT.md FIRST (it is recovery's input), crash recovery (recover own-crash files; report foreign claims, ask before takeover), validate frontmatter, pick first pending
□ Step 2: Claim request (mkdir -p working/ + move, update status & claimed_at)
□ Step 3: Triage (decide route, append ## Triage, read original if addendum)
□ Step 3.5: Handle Open Questions (mark - [~] with D-XX numbered decisions; a user answer obtained mid-run is written in as - [x] before dispatch — never - [~], no D-XX)
□ Step 4: Plan (Route C: spawn Plan agent + validate plan / Routes A & B: note skipped)
□ Step 5: Explore (Routes B & C: spawn Explore agent, include prime file lessons)
□ Step 5.5: Scope Declaration (Routes B & C: declare files + acceptance criteria in REQ)
□ Step 5.75: Pre-Flight Check (Routes B & C: git clean, test baseline, dependencies)
□ Step 6: Implement (spawn agent with lessons + TDD mode if set, log decisions as D-XX)
□ Step 6.25: Implementation Summary (append file manifest — mandatory for all routes)
□ Step 6.3: Qualify (orchestrator verifies: files exist, substantive, wired, flowing, requirements traced, P-A-U audit)
□ Step 6.5: Test (run relevant tests, load debug rules on attempt 2+, verify TDD evidence if tdd:true)
□ Step 7: Review (spawn actions/review-work.md — gate on acceptance: Pass→archive, Fail→remediate with debug rules)
□ Step 7.5: Lessons Learned + Orientation (append sections at subsystem altitude, update prime files, skip lessons for Route A if no surprises)
□ Step 8: Archive (update status, classify failures, triage discovered tasks, cycle-check follow-ups, queue follow-ups, move to archive/)
□ Step 9: Commit (stage explicit files, commit if git repo, write hash to REQ via tools/checks/record-commit-hash.sh, then a separate metadata commit)
□ Step 10: Loop or Exit (context wipe + contamination check if looping, else write CHECKPOINT.md with depth + cleanup)
```


## Error Handling

| Phase | Action |
|-------|--------|
| `pending-answers` REQs remain after queue is empty | Report them to the user: list each REQ and its unresolved questions. Suggest `do-work clarify` to batch-review. |
| `blocked` REQs remain after queue is empty | Report each with its `blocked_by` condition and age (per the composed exit summary). Suggest re-running `do-work run` (auto-probes any `blocked_check`) once a condition is met, or `do-work clarify` to confirm a human-checkable one. |
| Plan agent fails (Route C) | Classify failure (Intent/Spec/Code/Environment), create follow-up REQ if applicable, archive as failed |
| Explore agent fails (B/C) | Proceed to implementation with reduced context — builder can explore on its own |
| Builder reports a missing external precondition | Apply the blocked-flip test (Step 8 → **Mid-run blocked flip**): if no substantive implementation edits landed AND the condition is expected to self-resolve, flip to `status: blocked` (non-terminal, back to the queue) instead of failing. Otherwise classify as Environment and archive as failed. |
| Implementation fails | Classify failure (Intent/Spec/Code/Environment), create follow-up REQ if applicable, archive as failed |
| Tests fail repeatedly | After 3 fix attempts, classify as Code failure, create follow-up REQ with test failure details, archive as failed |
| Review: Acceptance = Fail | Return to Step 6 for ONE remediation attempt, then re-review. If still failing: archive as `completed-with-issues` with follow-up REQs |
| Review work agent fails | Skip review, note it in the REQ file, continue to archive — review failure is not a gate |
| Commit fails | Investigate the error (usually a pre-commit hook failure). Fix the underlying issue, re-stage, and retry as a **new** commit (never bypass — see `## Rules`). If unfixable, report the error to the user and continue to next request — changes remain uncommitted but archived. |
| Unrecoverable error | Stop loop, report clearly, leave queue intact for manual recovery |

## Progress Reporting

Keep the user informed with this format:

(keep the user informed in the running per-REQ progress format shown in `actions/work-reference.md` → **Progress Reporting Example**)

When the run finishes or pauses, hand back with the **Decision Brief** (`actions/work-reference.md` → **Decision Brief (hand-back format)**): lead with WHAT'S BEING BUILT (each REQ's `## Orientation`, at subsystem altitude), then DECISIONS FOR YOU (any escalated `pending-answers` follow-ups, each with value + risk), then a collapsed HANDLED list. Never lead with review scores — they stay in the per-REQ progress lines, not in front of the hand-back.


## Archived Request File Example

See [sample-archived-req.md](./sample-archived-req.md) for a complete example of what an archived REQ looks like after processing through the full pipeline (Route B). Every section shown there is generated by the steps above.

**Timestamps tell the story:** `created_at` → `claimed_at` = queue wait time. `claimed_at` → `completed_at` = implementation time. Route + timestamps let you calibrate triage accuracy over time.

## Rules

- The orchestrator handles ALL file management (moving files, updating frontmatter, appending sections, archiving). Spawned agents do implementation work only.
- Only two frontmatter status transitions are written on the normal path: `pending` → `claimed` (Step 2), then `claimed` → final status (Step 8); exception paths (Steps 1, 2.0, the Step 8 mid-run blocked flip to `blocked`, and 7's failed-remediation write) set the documented special statuses. Intermediate phases are tracked by which `##` sections exist, not by status.
- One commit per request; stage explicit files only — never `git add -A`/`.` or bypass a commit hook (see `actions/commit.md` § Rules for the full staging/hook guard). This governs serial mode; in worktree dispatch mode the implementation's commit boundary is the merge (builder commits + the `--no-ff` merge commit), and Step 9's commit carries only the changelog/version/metadata (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**).
- `write_set` is display-only, and stays display-only at any builder count: under fan-out it is **advisory input to the human's pick, never a gate**, and the merge — not the field — is the non-interference proof (`actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). Nothing schedules, gates, or dispatches on it; it feeds the board's overlaps badge (`actions/work-reference.md` Request File Schema).
- In worktree dispatch mode the orchestrator is the sole integrator: the builder never writes the main tree, and a merged state is not a verified state until the REQ's checks re-run there.

**Common mistakes to avoid:**

- Spawning implementation agent without first moving file to `working/`
- Letting spawned agents handle file management (only the orchestrator moves/archives files)
- Forgetting to update status in frontmatter (normal path has only two transitions: `claimed` at Step 2, final status at Step 8)
- Archiving a UR folder before all its REQs are complete
- Forgetting Planning status note for Routes A/B ("Planning not required")
- Committing without validating Implementation Summary file list against staged files
- Implementation Summary that only lists `do-work/` paths (means the REQ wasn't actually implemented — exception: `domain: ui-design` design artifacts placed in project directories like `docs/design/`)
- Creating follow-ups for every `- [~]` item instead of only UX-affecting decisions

**This action does NOT:**

- Create new request files (use actions/capture.md)
- Make architectural decisions beyond what's in the request
- Run without user present (this is supervised automation)
- Modify already-completed requests
- Allow external modification of files in `working/` or `archive/`

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "I'll skip Pre-Flight — the baseline is probably stable" | Run `git status` and the test baseline anyway (Step 5.75) | Pre-existing failures get misattributed to the builder, and unrelated dirty files get swept into the commit |
| "I wrote the test after the code but it fails without it, so this counts as TDD" | For `tdd: true`, write the failing test first and show it RED before the code | Post-hoc tests encode the implementation's quirks; the RED-before-GREEN ordering is the evidence Step 6.5 gates on |
| "P-A-U is bookkeeping — I'll just tick the boxes" | Do each phase; Step 6.3 audits the diff against the checked boxes | A checked `[UNIFY]` over a diff containing `console.log` is a false claim the qualifier will catch |
| "This file change is small — it doesn't need to go in the Scope section" | Declare every file before coding (Step 5.5) | Undeclared touches are exactly what the scope-drift check flags at review; "small" is judged after the fact, not before |
| "Tests still fail on attempt 2, but I'll just try the same fix again" | Load `crew-members/debugging.md` and `testing.md` before retrying | Unstructured retries repeat the same dead end; the debugging methodology exists for the 2nd+ attempt |
| "The Implementation Summary is too detailed — I'll just write 'updated logic'" | List every changed file with its action verb + a factual one-liner | The Summary is the primary auditability artifact; "updated logic" is unverifiable and reads as a hollow completion |
| "I'll fix this out-of-scope thing inline while I'm here" | Record it in `## Discovered Tasks`; Step 8 classifies and queues it | Inline scope creep escapes triage, review, and the per-REQ commit boundary |
| "The queue file's twin is already archived, but re-running is harmless" | Stop — Step 2.0 sets `blocked-archive-collision` for exactly this | Re-processing a duplicate silently re-commits it and corrupts the archive lineage |
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
- [ ] No REQ files remain in `do-work/working/` after the work loop ends — except a reported foreign claim this run deliberately left intact (Step 1 Crash Recovery)
- [ ] CHECKPOINT.md written if ending mid-session (for resume)
- [ ] Git commit created for each completed REQ
- [ ] Cleanup pass triggered at end of work loop
