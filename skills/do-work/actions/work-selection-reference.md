# Work phase: selection

> Companion to [work.md](./work.md); core queue procedures. Read when its phase routing selects this file. Existing step numbers remain stable.

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
- **`--fan-out [N]`** (optional integer) — enter **auto-wave mode**: dispatch builders for the selected set concurrently, bounded to N. Bare `--fan-out` uses the harness concurrency limit, or **two** where that is unknown. It changes *how many* of the selected set run at once, never *which*, so it composes with `--wave N` and targeting tokens. Without worktree support it degrades silently to the serial loop (**Architecture**; policy: `actions/fan-out-reference.md` → **Auto-wave**).
- **`--coordinate`** (boolean flag) — the main session coordinates and never builds or integrates: the coordinator shape becomes mandatory, with its run rules, preflight and stall loop (`actions/fan-out-reference.md` → **Delegated integration — the coordinator shape**). It composes with `--fan-out [N]`, `--wave N` and targeting tokens; bare `--coordinate` implies `--fan-out`. The action consumes the flag and never passes it to queue-mode `advance`, which refuses any token it does not know; when the user gave no `--fan-out`, it passes a bare `--fan-out`. `advance`'s continuation does not carry the flag, so the session keeps the mode for the whole run. Without worktree or agent-dispatch support it prints one line, `--coordinate: no worktree or agent dispatch here, running the serial loop`, and runs the serial loop.
- **`--review-mode sol-loop|council`** — select the optional AI-Evaluation-Loop for each selected REQ; `--jester` adds an optional Gemini advisory reviewer and requires a review mode. Read `actions/ai-evaluation-loop-reference.md` in full when selected or when a claimed REQ already carries its persisted selection. Reject missing/unknown/duplicate mode values and bare `--jester` before queue mutation. These are action-owned options: consume them before every CLI `advance` call, retain them across continuations, and persist them in each claimed REQ before dispatch. Omitted options preserve the ordinary review unless this REQ is resuming a selected loop.
- **`--wave N`** (integer flag, default mode only) — run only REQs at dependency depth N. Roots (no `depends_on`, or all `depends_on` resolve to archived REQs) are depth 0; depth grows by one per dependency layer. Mutually exclusive with **any** targeting token (`REQ-` or `UR-`) — reject the combination with an error.
- **`--skip-impact-negligible`** (boolean flag) — omit every REQ whose `impact:` resolves to `impact-negligible` from Step 1's selection, and report the ones it dropped. It changes *which* REQs are selected, never how many run at once, so it **composes with `--fan-out` and stacks with `--wave`**: `--wave` chooses the depth, this flag subtracts the negligible REQs from it, `--fan-out` chooses how many of the remainder run concurrently. **An explicitly-named `REQ-NNN` overrides it** — the user named the REQ outright, exactly as explicit naming overrides `depends_on` and `assigned_to` — while a REQ reached by `UR-NNN` expansion does **not** override it, per the same per-token provenance rule. Deliberately one boolean and not a general `--impact <token>` selector: stopping negligible work is the only filter anyone has asked for, and a second use would be the time to generalize.

**Unrecognized arguments are rejected, not ignored.** After consuming `--review-mode sol-loop|council` and `--jester`, stripping `--wave N`, `--fan-out [N]`, `--coordinate`, and `--skip-impact-negligible` and extracting targeting tokens (`REQ-`/`UR-` followed by digits, case-insensitive, per the Target ID Resolution contract), any non-empty token still left in `$ARGUMENTS` is an error. Stop and report:

```
Unrecognized argument(s): <tokens>. Usage: do-work run [REQ-NNN|UR-NNN ... | --wave N] [--skip-impact-negligible] [--coordinate] [--fan-out [N]] [--review-mode sol-loop|council [--jester]]
```

Do **not** fall through to full-queue processing. A leftover token almost always means the user meant to *scope* the run — a typo'd REQ ID (`REG-042` instead of `REQ-042`), or dead muscle memory (a retired mode word) — so silently building the entire queue is the wrong, hard-to-undo default. This generalizes the existing `--wave`-plus-targeting-token rejection to all unrecognized residue; both are parse-time guards.

When `$ARGUMENTS` is empty — no targeting tokens, no flags, no other tokens — process all pending REQs in dependency-aware order (default behavior).

## Steps

**actions/work.md is an orchestrator.** You own semantic judgment and action-authored evidence. Canonical commands own deterministic file management; spawned agents handle implementation work only.

### Step 1: Select and Claim

At each queue boundary, run canonical `recover`, then queue-mode `advance` with the user's targeting and selection arguments, and treat its typed result as the sole authority: a REQ is claimable when its normalized queue state, targeting provenance, dependency state, assignment, impact filter, external probe, and collision evidence permit it. **Recovery is read per REQ, never as one pass/fail gate for the run:** a finalization record whose `reason_codes` carry `FINALIZATION-SET-ASIDE` excludes that one REQ from this run's selection, its other reason codes say what refused, and the run continues with what remains, while a typed refusal — dirt no REQ owns — stops the run and hands off to judgment for the resolving verb (`actions/work-reference.md` → **Commit & Metadata-Commit Procedure (Step 9)** for the per-record contract, → **Stuck Runs Hand Off to Judgment (any step)** for the verb). The command owns selection, holds, unblock, claim, frozen membership, dispatch bounds, and tokenized continuation; the action dispatches only committed `queue_advance.claimed` members, follows the returned continuation unchanged after integration, and otherwise renders **Composed Exit Summary (Step 1)** from the same result without reproducing any mechanic.

After a committed claim, restore or initialize its selected AI-Evaluation-Loop record before any builder/reviewer dispatch, per `actions/ai-evaluation-loop-reference.md`. Pass that durable selection to a delegated integrator; CLI continuations do not carry action-owned options.

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
