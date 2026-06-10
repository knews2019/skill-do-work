# Ultracode Workflow — Fable/Opus/Sonnet Delegation Principles

> Model-tiered delegation policy for code work: cheap models execute, expensive models judge. Plan first, delegate implementation down the ladder, verify with your own test run, gate every diff through a fresh-context Fable review, escalate on repeated failure. The goal is maximum quality per unit of **cost** — not per token, since tokens are priced differently by model.

**Aliases:** `fable-opus-sonnet-workflow-principles`
**When to use:**
- A multi-subtask coding job where implementation can be delegated and you want frontier-model judgment without frontier-model execution costs.
- `do-work run` invoked in ultracode mode, or a REQ carries `ultracode: true` — `actions/work.md` Step 6 loads the **Mode B — Dispatch Policy** section of this file (only that section).
- Ad-hoc: `do-work prompts run ultracode-workflow <task>` applies the full Mode A workflow to the task in args.

**Inputs / flags:**
- **Mode A (standalone):** everything after the prompt name is the task description.
- **Mode B (adopted by `actions/work.md`):** no args — the work action reads only the Mode B section and keeps its own step structure.

**Maintenance note:** the model IDs below (`claude-fable-5`, `claude-opus-4-8`, `claude-sonnet-4-6`, Haiku) are deliberate, verbatim current-generation bindings — that is why the alias names the trio. When the lineup changes, update the IDs in this file in place; the tier *roles* (executor / escalation / judge) and the filename stay stable.

**Name note:** in Claude Code, the bare word "ultracode" is a native harness keyword that opts a session into built-in multi-agent orchestration. This prompt is related in spirit but independent — it is a delegation policy the agent follows on any host, not a harness feature. The compound trigger `ultracode-workflow` keeps the two greppable and distinct.

---

## Philosophy

You are in ultracode mode. The goal is maximum quality per unit of cost — not per token, since tokens are priced differently by model. Cheap models execute, expensive models judge: a run that spends more Sonnet tokens to save Fable tokens is a win. Every delegation decision should be defensible on both axes — never burn a strong model on routine work, never let weak-model output ship unreviewed.

Three invariants hold in every mode and at every degradation level:

1. **Mechanical truth.** The exit code of a test run you executed yourself is ground truth. Prose claims of green are never accepted.
2. **Fresh-context judgment.** The reviewer sees the spec and the diff — never the build transcript. A judge who watched the code being written is anchored, not independent.
3. **Bounded loops.** Fix iterations are capped; failure escalates up the ladder instead of retrying sideways.

## Modes

| Mode | Trigger | What runs |
|---|---|---|
| **A — Standalone** | `do-work prompts run ultracode-workflow <task>` | The full workflow below: capability check → plan → risk-tag → delegate → verify → deep review → Fable gate → report |
| **B — Dispatch policy** | `actions/work.md` Step 6, when the run was invoked in ultracode mode or the REQ's `ultracode` frontmatter is `true` | Only the **Mode B — Dispatch Policy** section. It assigns a model tier to each of work's existing steps. work.md's step structure, 3-attempt test loop, and review gates stay authoritative — never run Mode A's loop inside work.md |

Running both loops on the same work item is a failure mode: two escalation controllers will fight over a failing subtask. Mode B exists precisely so that doesn't happen.

## Step 0: Host Capability Check (both modes)

Before any delegation, determine two facts about the host you are running on:

1. Can it spawn subagents?
2. Can it pin a model per subagent?

Then announce the resulting level to the user — by name, with what is lost — and proceed. Never silently degrade, and never silently skip the gates.

| Level | Host capability | What you do |
|---|---|---|
| **1 — Full tiering** | Subagents + per-subagent model selection (e.g., Claude Code) | The workflow as written below |
| **2 — Role separation** | Subagents, but no model selection (some Codex CLI or other-host setups) | Keep the full role separation — planner, executor, mechanical verify, fresh-context reviewer — on whatever single model the host runs. The quality gates survive; the cost arbitrage is lost. Say so: "Host can't pin models per subagent — running ultracode role separation on one model; gates intact, cost optimization unavailable." If the host has a native equivalent (e.g., per-invocation model selection, profiles), suggest it to the user as an alternative way to recover the tiering. |
| **3 — Single-session sequential** | No subagents at all | Run the phases yourself, in order, one at a time. Before the review phase, re-read **only** the subtask spec and the diff (`git diff`) — not your own build reasoning — and review as adversarially as you can. Disclose in the final report that review independence was reduced. |

## Mode A — Standalone Workflow

### Step 1: Plan First

Analyze the task and break it into discrete subtasks with clear specs before any delegation. Each spec states what the subtask changes, where, and how success is observable. For complex tasks (cross-cutting changes, new architecture), draft the plan with `claude-opus-4-8`; otherwise plan yourself.

For multi-subtask runs, persist the plan as a manifest in a timestamped run directory per `crew-members/background-agents.md` (`do-work/runs/ultracode-<timestamp>/`) — one line of status per subtask. On re-entry after an interruption, read the manifest and resume from the first incomplete subtask instead of restarting; this is what makes the prompt idempotent and resumable.

### Step 2: Risk-Tag Every Subtask

Tag each subtask `standard` or `elevated` at plan time. **Elevated** means a happy-path test can pass while the feature is still broken: persistent state and serialization, hydration/caching, data migrations, concurrency, payments/idempotency, auth/security.

Deep review is exception-based. If more than roughly a quarter of subtasks come out tagged elevated, re-triage — the tag is meaningless if everything carries it.

### Step 3: Delegate Execution

Implementation, test writing, test running, and mechanical refactors go to subagents running `claude-sonnet-4-6`. Trivial mechanical steps (renames, formatting, boilerplate) may use Haiku. The immediate-escalation list in the ladder below overrides this default — those subtasks start on `claude-opus-4-8`.

When subtasks run in parallel, follow the durability pattern in `crew-members/background-agents.md` and never let two executors mutate the same working tree concurrently — use isolated worktrees or strictly disjoint file sets; otherwise run subtasks sequentially.

### Step 4: Verify Mechanically

When an executor reports done, run the project's canonical test command yourself — full suite, not a subset. The exit code is ground truth. Executor reports are never the final word on green, and executors must include the exact command run and the verbatim summary line of its output in their reports.

**If the project has no canonical test command:** either (a) have the executor tier write characterization tests for the touched behavior first and gate on those, or (b) verify by running the program or feature directly and recording what you observed. In both cases, disclose in the final report that no full-suite gate existed.

### Step 5: Deep Review (elevated subtasks only)

After the suite is green and before the Fable gate, elevated subtasks get a `claude-opus-4-8` pass that hunts edge cases — stale or partial state, serialization roundtrips, ordering, failure injection — and specifies the missing tests as concrete cases. A Sonnet subagent implements those tests; you re-run the full suite.

### Step 6: Fable Review Gate

After tests pass, review the diff using `claude-fable-5` — **in a fresh context**. The judge receives the subtask spec, the diff, and the names of the tests that gate it. The judge never receives the build transcript (invariant 2). If you, the orchestrator, are yourself running on Fable, still spawn a fresh-context judge subagent when the host allows it; at Level 3, apply the Step 0 re-read discipline instead.

The review checks: correctness, edge cases, code quality, and whether the tests actually verify the requirement. If the diff touches test files, confirm the changes don't weaken assertions, loosen matchers, or skip tests. For elevated subtasks, Fable verifies the deep-review findings were addressed rather than re-deriving them from scratch.

### Step 7: Iterate or Escalate

If review finds issues, send specific, actionable feedback back to a Sonnet subagent. Maximum 2 fix iterations per subtask **per tier** — escalation hands the subtask to the next tier with a fresh 2-iteration budget, it does not extend the current tier's. Follow the ladder below.

### Step 8: Report

Render the final report per **Output Format** below.

## Escalation Ladder (both modes)

- A Sonnet subagent fails the same subtask twice → escalate that subtask to `claude-opus-4-8`.
- Escalate to Opus immediately (skip Sonnet entirely) for: cross-cutting refactors where a wrong move cascades, concurrency logic, payments/idempotency logic, and security-sensitive code.
- If Opus also fails, Fable may fix it directly — this is the last resort, not a shortcut. **A Fable-authored fix does not review itself:** spawn a fresh-context judge (a `claude-fable-5` subagent fed only the spec and the diff) or run a `claude-opus-4-8` adversarial pass on the diff. At Level 3, where neither is possible, disclose the self-review explicitly in the report.
- If Fable cannot resolve it cleanly, stop and report the blocker with what was tried at each tier.

## Mode B — Dispatch Policy for `actions/work.md`

When `actions/work.md` Step 6 loads this section (run-level ultracode mode, or REQ frontmatter `ultracode: true` per the Schema Read Contract; an explicit canonical `ultracode: false` opts a REQ out of a run-level mode):

**work.md stays the controller.** Its pipeline (triage → plan/explore → implement → qualify → test → review → archive → commit), its 3-attempt test loop, and its review/remediation gates are authoritative. This policy decides only *which model tier runs each step*:

| work.md step | Tier |
|---|---|
| Step 4 — Route C plan agent | `claude-opus-4-8` |
| Step 6 — builder (all routes) | `claude-sonnet-4-6`; Haiku acceptable for trivial Route A one-liners (rename, copy, config value). The immediate-escalation list above overrides: concurrency, payments/idempotency, security-sensitive, cascading refactors start on `claude-opus-4-8` |
| Step 6.5 — test loop | Attempts 1–2: Sonnet (attempt 2 loads the debugging/testing crews per work.md as usual); attempt 3: escalate the builder to `claude-opus-4-8` |
| Step 7 — review (review-work spawn) | `claude-fable-5`, fresh context: the REQ and the diff, never the build transcript |
| Step 7 — remediation after a failed review | Sonnet + debugging crew; Opus if the same review failure repeats |

Additional Mode B rules:

- **Risk tag at triage.** Tag the REQ `standard` or `elevated` using Mode A Step 2's criteria; record the tag in the `## Triage` section. Elevated REQs get the Opus deep-review pass between work's Step 6.5 (suite green) and Step 7 (review).
- **Record the tiers.** Note in the REQ's `## Implementation Summary` and `## Review` sections which model tier handled each phase and any escalations triggered — that is the audit trail work.md's living-log philosophy expects.
- **Degraded hosts.** If the host can't pin models per subagent, announce it per Step 0 and let work.md proceed with its normal dispatch — work's pipeline already provides the role separation.
- **Never import Mode A's loop.** The 2-iteration cap and the subtask loop above do not replace work's 3-attempt loop; the ladder maps onto work's attempts as the table shows.

## Hard Rules

- Never use Fable or Opus for routine implementation.
- Never skip the Fable review gate, even for one-line changes.
- Never accept a prose claim of passing tests — verify via your own test run.
- Executors must include the exact command run and the verbatim summary line of its output in their reports.
- The judge never sees the build transcript — spec and diff only.
- Unbounded loops are a failure mode: respect the 2-iteration cap per tier and the escalation ladder.
- Capability shortfalls are announced at Step 0 and disclosed in the report — never silently degraded around.

## Output Format

At the end, summarize:

- Subtasks completed, which model handled each, and the risk tag assigned to each.
- Escalations and deep reviews triggered, and why.
- Final test result: the exact command and the verbatim summary line of your own run (or the disclosed no-suite fallback).
- Fable's review verdict per subtask.
- **Cost:** per-tier token/cost usage if the host exposes it; otherwise state plainly "cost not measured by host" — never imply the cost optimization was measured when it wasn't.
- **Disclosures:** the Step 0 capability level the run executed at, any no-test-suite fallback, and any self-review at the ladder's last resort.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "This subtask is tiny — I'll just implement it myself on the judge model" | Delegate to Sonnet or Haiku anyway | Routine work on Fable breaks the cost model *and* leaves the diff without an independent reviewer |
| "The executor pasted test output — looks green to me" | Run the canonical command yourself | Executor reports are claims; the exit code of your own run is ground truth |
| "The review gate is overkill for a one-line diff" | Run the Fable gate | One-line diffs in elevated areas (auth, idempotency) are exactly where the gate pays for itself |
| "One more fix iteration at this tier will do it" | Escalate per the ladder | The cap exists because same-tier retries converge on the same blind spot |
| "Everything here feels risky — I'll tag it all elevated" | Re-triage until elevated is roughly a quarter or less | A tag carried by everything selects nothing; the deep-review budget must concentrate where happy-path tests lie |
| "This host can't pin models, so the workflow doesn't apply" | Announce Level 2/3 and keep the gates | Mechanical verification and fresh-context review survive without model tiering; only the cost arbitrage is lost |
| "Fable wrote the fix, so Fable's gate is already satisfied" | Spawn a fresh-context judge or an Opus adversarial pass | Author self-review is the weakest gate in the system, attached to the riskiest path — code two tiers already failed on |

## Red Flags

- A report cites passing tests but contains no exact command + verbatim summary line from the orchestrator's own run.
- The judge's review references implementation reasoning that only exists in the build transcript — context leaked; the review wasn't fresh.
- Every subtask is tagged elevated, or none are.
- More than 2 fix iterations happened at one tier without an escalation.
- A diff touching test files passed review with no explicit statement about assertion weakening, loosened matchers, or skipped tests.
- The run degraded to Level 2/3 but the final report doesn't disclose it.
- (Mode B) Two escalation loops ran on the same REQ — Mode A's subtask loop was applied inside work.md's attempt loop.

## Verification Checklist

- [ ] Step 0 capability level was announced before any delegation
- [ ] Every subtask has a written spec and a `standard`/`elevated` tag recorded before delegation
- [ ] The full-suite test command was run by the orchestrator after each subtask, with the exact command + summary line recorded (or the no-suite fallback disclosed)
- [ ] Every diff passed a fresh-context `claude-fable-5` review (or the Level 3 re-read review, disclosed)
- [ ] Elevated subtasks show an Opus deep-review pass and the resulting tests before the Fable gate
- [ ] Escalations followed the ladder with fresh 2-iteration budgets, each recorded with the reason
- [ ] The final report includes per-subtask model, tags, escalations, test result, Fable verdict, and cost (or "not measured")
- [ ] Mode B: the REQ's sections record which tier handled each phase, and work.md's loop structure was left unchanged
