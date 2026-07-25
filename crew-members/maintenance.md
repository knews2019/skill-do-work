# The Subtractor — Maintenance Crew Member

<!-- JIT_CONTEXT: Loaded before a deliberate *maintenance pass* on the skill's own operating instructions — an agent/action/crew/prime file that has drifted, bloated, or started teaching the wrong thing — where removing or narrowing is a candidate fix. This is the maintenance-time window where removal IS the task, distinct from feature implementation (coding-guardrails.md governs that). Auto-loaded by `actions/work.md` Step 6 when a REQ carries the `maintenance: true` marker, which capture sets for a removal/narrowing finding on the skill's own instructions — e.g. one surfaced by `actions/quick-wins.md` (read-only — it *surfaces* removable rules/config). The marker is deliberately the *only* trigger: there is no description heuristic, because one would misfire on ordinary implementation REQs (which routinely touch adjacent dead code) and load the opposite posture from coding-guardrails'. The trigger is the condition above, not any one caller — callers are illustrative, not the boundary; a future dedicated `maintain` action sets the same marker. -->

> Maintenance is the deliberate window where removal is the task. Before you add an instruction to fix a drift, try deleting or narrowing one — most harnesses rot because every fix is one more rule.

## Principles

### 1. Delete before you add

When an agent, action, or instruction has drifted, the first move is subtraction, not addition. A new rule is the most expensive fix: it is permanent, it competes for attention with every other rule, and it compounds. Reach for removal or narrowing first; add only when subtraction provably can't do the job.

### 2. Ask the deletion questions

Before writing a new instruction, ask what existing thing is *causing* the drift:

- Is a **stale source** feeding it? (an out-of-date example, a wrong default, a dead link) → fix or remove the source.
- Is a **bad example** teaching it? (the model is pattern-matching on something you'd never want repeated) → delete the example.
- Is a **tool too broad**? (an option, flag, or path that invites the wrong behavior) → narrow its scope.
- Is the **job too vague**? (open-ended enough to drift) → constrain the job, don't pile on caveats.

A "yes" to any of these is a fix *by removal*. Only when all four are "no" is a new instruction the right tool.

### 3. Prove any addition against a replay pack

If you must add, earn it. An addition is justified only when a concrete case **fails without it and passes with it** — the maintenance analogue of a `## Red-Green Proof` (`actions/capture.md`). No replay case, no addition: an instruction that fixes nothing reproducible is bloat that will outlive the problem it imagined.

## Persistence

Active for the duration of a maintenance pass — from the moment removal or narrowing is on the table as a fix until the pass is committed. Re-engage for each drifting file in the same pass. Drops when the work transitions from maintenance into feature implementation, where `coding-guardrails.md`'s surgical-changes rule takes over.

## Boundaries

- **Complementary to `coding-guardrails.md`, never in conflict.** coding-guardrails § 3 "Surgical Changes" says *don't* delete adjacent or pre-existing dead code **while implementing a feature** — there, removal is scope creep. This file governs the opposite window: a deliberate maintenance pass where removal IS the assigned task. Implementation-time = leave it alone; maintenance-time = removal is the point. If you can't tell which window you're in, you're implementing — default to coding-guardrails.
- **Points at YAGNI, doesn't restate it.** coding-guardrails § 2 "Simplicity First" is the canonical YAGNI home (don't add speculative code). This file is its maintenance-time twin: don't add speculative *instructions*. Read them together.
- **Subtraction is not vandalism.** coding-guardrails' "Simplify ≠ strip" still holds: if removing something would have to be restored next week, it was foundation, not bloat. Delete what causes the drift, not what holds the thing up. When in doubt, narrow rather than delete — and record the call.
- **Loaded alongside other crew rules**, not instead of them — general, coding-guardrails, anti-slop, and domain rules still apply.

## What this looks like in practice

- **A recurring correction** (e.g. surfaced by forensics' Recurring Corrections check): the same lesson keeps reappearing across REQs. The reflex is to add a guardrail. The maintenance move is to ask *why* the harness keeps teaching the wrong thing — usually a stale source or a bad example — and remove that, so the correction stops recurring at the root.
- **A stale enumeration**: a hand-maintained list of a closed set has drifted. Don't add a "remember to update this list" note — generalize the list to a trigger condition so there's nothing to keep in sync.
- **An action that grew caveats**: an action file has accreted five edge-case clauses. Before adding a sixth, check whether narrowing the action's input, or fixing the upstream source, makes some of the existing five unnecessary.
