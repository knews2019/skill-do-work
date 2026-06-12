# Ultracode Fable Workflow — Fable/Opus/Sonnet/Haiku Batch Orchestration

> Batch orchestration policy for the work queue: the session model takes two turns per batch — launch one background Opus orchestrator that owns the batch end-to-end, then audit the result. Cheap models execute, mid models orchestrate and judge, the session model audits. The goal is maximum quality per unit of **cost** — not per token, since tokens are priced differently by model.

**Aliases:** `ultracode-fable` (short form), `ultracode-workflow`, `ultracode` (both pre-0.90 names), `fable-opus-sonnet-workflow-principles`
**When to use:**
- `do-work run ultracode-fable-workflow` (or `ultracode-fable`) — `actions/work.md`'s Input diverts the session here instead of running the work loop inline. Optional REQ IDs or `--wave N` scope the batch.
- `do-work prompts run ultracode-fable-workflow` — same architecture, invoked directly; operates on the queue exactly as above.
- NOT for ad-hoc free-text tasks — the former standalone mode is retired. Capture the task first (`do-work capture request: ...`), then run a scoped batch (`do-work run REQ-NNN ultracode-fable`).

**Inputs / flags:** no free-text task. Scope arrives via the work action's Input: targeted REQ IDs or `--wave N` select which pending REQs form the batch; with neither, the batch is the full dependency-ready queue and queue draining applies (see the Orchestrator Contract).

**Maintenance note:** the model IDs below (`claude-fable-5`, `claude-opus-4-8`, `claude-sonnet-4-6`, Haiku) are deliberate, verbatim current-generation bindings — the `-fable` in the name marks the Claude family (Fable/Opus/Sonnet/Haiku) this policy's Tier Table is tuned for. When the lineup changes within the family, update the IDs in this file in place; the tier *roles* (executor / orchestrator-judge / auditor) and the filename stay stable.

**Name note:** in Claude Code, the bare word "ultracode" is a native harness keyword that opts a session into built-in multi-agent orchestration. This prompt is related in spirit but independent — it is a delegation policy the agent follows on any host, not a harness feature. The compound trigger `ultracode-fable-workflow` keeps the two greppable and distinct, and names the model family the tiers bind to (`ultracode-fable` is the accepted short form). `ultracode-workflow` and bare `ultracode` are the pre-0.90 names, kept as aliases.

**Migration note:** this file replaces the former Mode A / Mode B split (≤ 0.88.x). Batch orchestration is the only ultracode shape. The per-REQ `ultracode:` frontmatter field is retired and ignored — scope the batch with REQ IDs or `--wave` instead (see `actions/work-reference.md` → Schema Read Contract → Retired Fields).

---

## Philosophy

You are in ultracode mode. The goal is maximum quality per unit of cost — not per token, since tokens are priced differently by model. Cheap models execute, mid models orchestrate and judge, the session model audits. One structural fact dominates all others: every turn the session model takes re-reads the entire conversation as input, almost always cache-missed because batch workflows outlast the cache window. Therefore the session model's turn count is the primary cost control, ahead of any per-token optimization.

Four invariants hold at every degradation level:

1. **Mechanical truth.** The exit code of a test run you executed yourself is ground truth. Prose claims of green are never accepted.
2. **Fresh-context judgment.** The reviewer sees the spec and the diff — never the build transcript. A judge who watched the code being written is anchored, not independent.
3. **Bounded loops.** Fix iterations are capped; failure escalates up the ladder instead of retrying sideways.
4. **Touch economy.** The session model takes exactly two turns per batch — launch and audit. Everything in between belongs to the orchestrator.

## Step 0: Host Capability Check

Before launching anything, determine three facts about the host you are running on:

1. Can it run a detached/background agent?
2. Can it pin a model per agent?
3. Can it spawn subagents at all?

Then announce the resulting level to the user — by name, with what is lost — and proceed. Never silently degrade, and never silently skip the gates. The level is disclosed again in the digest and the audit report.

| Level | Host capability | What you do |
|---|---|---|
| **1 — Full batch** | Background agents + per-agent model selection (e.g., Claude Code) | The architecture as written below |
| **2 — Foreground batch** | Subagents, but no detached background agents | Run the orchestrator as ONE synchronous subagent. The touch economy survives as two decision points — launch, then audit when the subagent returns its digest — but the session turn stays open for the duration. The conversation stays lean because only the ≤20-line digest comes back. Say so: "Host can't run detached background agents — running the batch orchestrator as one synchronous subagent; launch + audit contract unchanged, the session waits instead of yielding." |
| **3 — Role separation** | No per-agent model selection (composes with Level 1 or 2) | Keep every role boundary — orchestrator, executors, fresh-context reviewers, auditor — on whatever single model the host runs. The quality gates survive; the cost arbitrage is lost. Say so: "Host can't pin models per subagent — running ultracode role separation on one model; gates intact, cost optimization unavailable." If the host has a native equivalent (per-invocation model selection, profiles), suggest it to the user as a way to recover the tiering. |
| **4 — Single-session sequential** | No subagents at all | There is no batch architecture to run. Announce it, then run `actions/work.md`'s normal loop yourself with the discipline overlay: risk-tag at triage, full canonical suite before every commit, and before each Step 7 review re-read **only** the REQ and the diff (`git diff`) — never your own build reasoning. Keep the run directory and manifest for resumability, and disclose the reduced review independence in the final report. |

## Session Model Contract — Two Touches Per Batch

The session model (`claude-fable-5`) takes exactly two turns per batch:

### Touch 1: Launch

1. Read the queue (frontmatter scan only — no REQ bodies into the conversation) and apply the scope: targeted REQ IDs, `--wave N`, or the full dependency-ready queue.
2. **Pre-flight clarity sweep.** This is the only human-attention window before the audit. A REQ that is genuinely unworkable as written (contradictory requirements, missing decision only the user can make) gets flipped to `pending-answers` now, with the question recorded — the user resolves it later via `do-work clarify`. Do not launch a batch whose REQs need answers mid-build.
3. Create the run directory `do-work/runs/ultracode-<timestamp>/` with a manifest (one status line per REQ) per `crew-members/background-agents.md` — disk is the source of truth, not this conversation.
4. Start ONE background `claude-opus-4-8` batch orchestrator with a launch packet: the batch scope, the run-directory path, and pointers to this file's **Batch Orchestrator Contract** and to `actions/work.md`. The packet never contains the `ultracode-fable-workflow` mode word (in any form) — the orchestrator runs the work pipeline plainly (see Anti-recursion in the Orchestrator Contract).
5. Yield immediately.

### Touch 2: Audit

When the orchestrator's digest arrives, run the **Final Audit** below, then report. Between the two touches: no polling, no per-REQ check-ins, no inline fixes, no paperwork. The session model wakes only for the orchestrator's digest or a user message. Queue draining belongs to the orchestrator, not the session model.

## Batch Orchestrator Contract (`claude-opus-4-8`, one background agent)

**work.md stays the controller.** Per REQ, you execute `actions/work.md`'s pipeline exactly — its triage routes, 3-attempt test loop, review and remediation gates, archive and commit steps are authoritative. You never invent a parallel per-REQ pipeline; two escalation controllers fighting over one failing subtask is a known failure mode. **Anti-recursion:** you were launched without the `ultracode-fable-workflow` mode word, so `actions/work.md`'s Input never diverts back here — you run its normal Steps 1–10, applying the Tier Table below at each spawn.

Each batch duty lives inside work.md's existing steps:

| Duty | Where it lives |
|---|---|
| Plan subtasks with specs, risk-tag `standard`/`elevated` | work Step 3 (tag recorded in `## Triage`) + Step 4 |
| Delegate execution | work Step 6 builder spawn, tiered per the Tier Table |
| Verify mechanically | work Steps 6.3 + 6.5 — targeted tests for the touched area inside fix loops; the **full canonical suite, run by you, before each Step 9 commit**. The exit code is ground truth; executor reports are never the final word on green |
| Fresh-context review per diff | work Step 7, pinned to `claude-opus-4-8`. If the diff touches test files, confirm it doesn't weaken assertions, loosen matchers, or skip tests |
| Iterate or escalate | work's attempt loop + remediation gate; the Escalation Ladder maps onto them |
| Paperwork | work Steps 6.25 / 7.5 — a `claude-sonnet-4-6` (or Haiku) drafter may turn your structured notes into the Implementation Summary / Testing / Review / Lessons prose; you spot-check and remain the writer-of-record, since work.md gives ALL file management to the orchestrator |
| Archive + scoped local commit per REQ | work Steps 8–9 unchanged. Do NOT push — commits stay local for the session model's audit |
| Drain the queue | work Step 10's loop — **default full-queue mode only**. After finishing the batch, re-scan; newly dependency-ready REQs form a new sub-batch under the same run directory. In targeted or `--wave` mode, never process beyond the given scope — report newly-ready REQs in the digest's Queue state line instead |

Batch-level duties work.md has no opinion on:

- **Parallelism triage.** Before building, map file-overlap between the batch's REQs. Disjoint REQs may run as parallel lanes — concurrent work.md pipelines over strictly disjoint file sets or isolated worktrees; never let two executors mutate the same working tree concurrently. Serialize the full-suite gates and Step 9 commits across lanes (the commit procedure assumes a quiet index). Never run overlapping REQs concurrently — sequence them.
- **Manifest upkeep.** Update the run directory's manifest as each REQ changes state. A dead orchestrator must be resumable from disk alone.
- **Clarify filing.** A REQ that turns out to need user input mid-batch follows work Step 3.5: best-judgment notes, `pending-answers` follow-up, move on. Never block the batch and never wake the session model for it.
- **The digest.** One report to the session model at the end, in the Digest Format below.

## Tier Table

| work.md step | Tier |
|---|---|
| Step 4 — Route C plan agent | `claude-opus-4-8` |
| Step 6 — builder (all routes) | `claude-sonnet-4-6`; Haiku acceptable for trivial Route A one-liners (rename, copy, config value). The immediate-escalation list below overrides: concurrency, payments/idempotency, security-sensitive, cascading refactors start on `claude-opus-4-8` |
| Step 6.5 — test loop | Attempts 1–2: Sonnet (attempt 2 loads the debugging/testing crews per work.md as usual); attempt 3: escalate the builder to `claude-opus-4-8` |
| Step 7 — review (review-work spawn) | `claude-opus-4-8`, fresh context: the REQ and the diff, never the build transcript. Never review your own orchestration inline |
| Step 7 — remediation after a failed review | Sonnet + debugging crew; Opus if the same review failure repeats |
| Deep review (elevated REQs, between Step 6.5 and Step 7) | fresh `claude-opus-4-8` |
| Paperwork drafting (Steps 6.25 / 7.5) | Sonnet; Haiku for mechanical sections |

## Escalation Ladder

- A Sonnet subagent fails the same subtask twice → escalate that subtask to a **fresh** `claude-opus-4-8` executor with a fresh 2-iteration budget — escalation hands the subtask up, it does not extend the current tier's budget.
- Escalate to Opus immediately (skip Sonnet entirely) for: cross-cutting refactors where a wrong move cascades, concurrency logic, payments/idempotency logic, and security-sensitive code.
- If the Opus executor also fails, mark the subtask **BLOCKED** in the digest with what was tried at each tier. Do not loop. The session model decides at audit whether to fix directly (last resort — and a session-model-authored fix gets a fresh-context `claude-opus-4-8` review, never self-review) or return the REQ to the queue.
- The ladder maps onto work.md's attempt loop and remediation gate — never run a second iteration cap alongside them.

## Risk Tags and Deep Review

- Tag every REQ's subtasks `standard` or `elevated` at triage; record the tag in the `## Triage` section. **Elevated** means a happy-path test can pass while the feature is still broken: persistent state and serialization, hydration/caching, data migrations, concurrency, payments/idempotency, auth/security.
- Elevated subtasks get a **deep review** after the suite is green and before the commit: a fresh `claude-opus-4-8` pass hunting edge cases — stale or partial state, serialization roundtrips, ordering, failure injection — specifying the missing tests as concrete cases. A Sonnet subagent implements them; re-run the suite.
- Deep review is exception-based. If more than roughly a quarter of subtasks come out elevated, re-triage — the tag is meaningless if everything carries it.

## Final Audit (session model, one turn per batch)

- Run the full canonical suite yourself — independent of the orchestrator's runs.
- Review `git log` for the batch: scope, message quality, nothing unexpected.
- Read diffs by priority: ALL elevated, escalated, BLOCKED, or partial-verdict items in full; sample roughly a third of standard items. Hunt specifically for the known failure seams: fixtures that misrepresent production shape, internally contradictory fix recommendations, weakened tests (`crew-members/testing.md` is the canonical statement of the first and last).
- Spot-check the archive and paperwork against the digest.
- **PASS → declare the batch releasable.** All commits stay local — this skill never pushes anything anywhere (see `actions/commit.md`); releasing is the user's call.
- **FAIL → file fix-REQs back to the queue** for the orchestrator (normal capture conventions, `addendum_to` set). Fix directly only if the orchestrator is demonstrably stuck.

## Run Directory & Resumability

The run directory (`do-work/runs/ultracode-<timestamp>/`, per `crew-members/background-agents.md`) is the source of truth for the whole batch: manifest, per-REQ status, and the digest written as a file with its path returned. Recovery after an interrupted session: re-invoke, find the newest incomplete run directory, resume from the manifest plus work.md's own crash recovery — never restart finished REQs. This is what makes the prompt idempotent and resumable.

## Context Hygiene (binding on every agent)

- No diffs, logs, or paperwork bodies in the main conversation. Write to files; report paths plus one-line summaries. (This is the same rule as `crew-members/background-agents.md`: disk is the source of truth, the transcript is not.)
- Test evidence is always: the exact command run + the verbatim summary line of its output. Nothing else from the output.
- Subagent reports use the digest format, hard cap ~20 lines. Full detail goes to files.
- The session model's conversation must stay lean enough that a wake-up is cheap even cache-cold.

## Hard Rules

- Never use Fable or Opus for routine implementation.
- Every diff gets a fresh-context `claude-opus-4-8` review (work Step 7), even one-line changes, and every batch gets the session-model audit. The pair replaces the old per-diff Fable gate — neither half substitutes for the other.
- Never accept a prose claim of passing tests — verify via your own run at your own gate. Executors must include the exact command run and the verbatim summary line of its output in their reports.
- The reviewer never sees the build transcript — REQ and diff only.
- No `git stash` as part of any workflow.
- Tests exercise the caller seam, not just the unit in isolation; fixtures must be production-faithful — a fixture that lies about store shape is a defect (canonical detail in `crew-members/testing.md`).
- RED labels must be truthful: a test claimed as failing-first must actually have failed (enforced by work.md's TDD-evidence gate).
- Respect the 2-iteration cap and the Escalation Ladder, mapped onto work.md's attempts; unbounded loops are a failure mode.
- Commits are local only. No push, ever.
- Capability shortfalls are announced at Step 0 and disclosed in the digest and the audit report — never silently degraded around.

## Digest Format (orchestrator → session model)

```
BATCH: <id>  REQs: <list>  Run dir: <path>  Capability level: <Step 0 level>
Per REQ: <REQ-id> | <verdict: SHIPPED/BLOCKED/CLARIFY> | <risk tags> | <escalations/deep-reviews + 1-line why> | <commit SHA> | <test command + verbatim summary line> | <paperwork path>
Anomalies: <anything a reviewer should look at first, or "none">
Queue state: <remaining / clarifies pending / newly-ready outside scope>
Cost: <per-tier usage if the host exposes it; otherwise "cost not measured by host">
```

## Output Format

The audit's user-facing report:

- Per-REQ verdicts (SHIPPED / BLOCKED / CLARIFY), which tier handled each phase, risk tags, and escalations or deep reviews triggered, with why.
- The auditor's own test result: the exact command and the verbatim summary line of the session model's independent run.
- Audit verdict: **RELEASABLE (commits local)** or **FIX-REQS FILED** (list them).
- **Cost:** per-tier token/cost usage if the host exposes it; otherwise state plainly "cost not measured by host" — never imply the cost optimization was measured when it wasn't.
- **Disclosures:** the Step 0 capability level the run executed at, any no-test-suite fallback, and any session-model-authored fix at the ladder's last resort.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "I'll just check on the orchestrator mid-run" | Wait for the digest | Every session-model turn re-reads the whole conversation, almost always cache-missed — polling is the cost model's single worst move |
| "This subtask is tiny — I'll implement it myself on the orchestrator/judge model" | Delegate to Sonnet or Haiku anyway | Routine work on an expensive model breaks the cost model *and* leaves the diff without an independent reviewer |
| "The executor pasted test output — looks green to me" | Run the gate command yourself | Executor reports are claims; the exit code of your own run is ground truth |
| "The review gate is overkill for a one-line diff" | Run the fresh-context Opus review | One-line diffs in elevated areas (auth, idempotency) are exactly where the gate pays for itself |
| "The digest looks complete — I'll skip reading the diffs" | Read all elevated/escalated/BLOCKED diffs in full, sample the standard ones | The digest is the orchestrator's claim about its own work; the audit exists because claims aren't evidence |
| "One more fix iteration at this tier will do it" | Escalate per the ladder | The cap exists because same-tier retries converge on the same blind spot |
| "Everything here feels risky — I'll tag it all elevated" | Re-triage until elevated is roughly a quarter or less | A tag carried by everything selects nothing; the deep-review budget must concentrate where happy-path tests lie |
| "This REQ needs the user — I'll ask now" | Best-judgment notes + `pending-answers` follow-up, move on | Human time is the bottleneck; questions batch into `do-work clarify`, and waking the session model mid-batch breaks the touch economy |
| "Everything passed — I'll push" | Declare releasable; commits stay local | This skill never pushes; releasing is the user's call |
| "This host can't background or pin models, so the workflow doesn't apply" | Announce Level 2/3/4 and keep the gates | Mechanical verification and fresh-context review survive without tiering; only the cost arbitrage is lost |

## Red Flags

- The session model took more than two turns for one batch with no documented reason (a user message is a reason; curiosity is not).
- A digest line cites passing tests but contains no exact command + verbatim summary line.
- The review references implementation reasoning that only exists in the build transcript — context leaked; the review wasn't fresh.
- More than 2 fix iterations happened at one tier without an escalation, or an Opus failure looped instead of going BLOCKED.
- Every subtask is tagged elevated, or none are.
- A diff touching test files passed review with no explicit statement about assertion weakening, loosened matchers, or skipped tests.
- REQ files are missing work.md's sections (Triage / Implementation Summary / Testing / Review) — the observable symptom that the orchestrator invented its own pipeline instead of running work.md.
- A push happened, or `git stash` appeared anywhere in the run.
- The run degraded to Level 2/3/4 but the digest or audit report doesn't disclose it.
- A multi-REQ batch ran with no run directory or a stale manifest.

## Verification Checklist

- [ ] Step 0 capability level was announced before launch and disclosed in the digest and audit report
- [ ] The run directory and manifest exist and reflect the final per-REQ states
- [ ] Every REQ file shows work.md's sections, with the handling tier and any escalations recorded in Implementation Summary / Review
- [ ] Every diff passed a fresh-context `claude-opus-4-8` review (REQ + diff only)
- [ ] Elevated REQs show the deep-review pass and the resulting tests before their commit
- [ ] Escalations followed the ladder with fresh 2-iteration budgets, each recorded with the reason; Opus failures went BLOCKED, not looped
- [ ] The orchestrator ran the full canonical suite before each commit, and the session model ran it independently at audit — both with exact command + verbatim summary line recorded
- [ ] The digest matches the archive and paperwork on spot-check
- [ ] All commits are local; nothing was pushed
- [ ] Cost is reported per tier or stated as "cost not measured by host"
