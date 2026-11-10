# Work phase: planning

> Companion to [work.md](./work.md); core queue procedures. Read when its phase routing selects this file. Existing step numbers remain stable.

### Step 4: Planning (Route C only)

**Route C:** Spawn a **Plan agent** with the request content, project context, the `crew-members/[domain].md` file (normalize `domain` per the Schema Read Contract first; if the resolved domain is missing, falls back to `general` for an unknown value, or the file doesn't exist, skip loading it), `crew-members/supply-chain.md` when the REQ creates a new service, app, package, or deployable, or changes dependencies or install, build, or image configuration (ask the plan to fill that file's **When Loaded at Planning** list), and any files listed in the `prime_files` array. Instruct it to use the prime files as the strict index for discovering the source of truth. Do not load global architecture. Ask it to produce a specific implementation plan (files to modify, order of changes, architectural decisions, testing approach). If a `## Plan` section does not already exist, append the output:

(append the plan per the **Plan Template — Route C (Step 4)** in `actions/work-reference.md`)

**Plan validation (Route C only):** After the Plan agent returns, run a quick quality check before proceeding:

1. **Requirement coverage:** Re-read the REQ's What/Detailed Requirements and any `## Addendum` section. Every requirement should map to at least one planned task. Flag uncovered requirements.
2. **No orphan tasks:** Every planned task should trace back to at least one requirement. Tasks that don't address any requirement suggest scope creep.
3. **Scope sanity:** Count the planned tasks. If 5+, flag: "Plan has [N] tasks — quality degrades past 3. Consider splitting this REQ into multiple smaller REQs."
4. **Consumer field contract:** For every planned command whose output drives an action-owned mutation, verify that the plan identifies the exact per-record identity, provenance, state, and outcome fields required by its consumer, as applicable. Flag plans that omit the consumer's required fields. This is the plan-time counterpart of `actions/review-work.md`'s **Restatement Sweep**, not a replacement; the sweep still runs at review.

Append validation findings to the `## Plan` section (if any issues found). These are **warnings, not blockers** — the builder can adapt. But flag them visibly so the orchestrator and review step are aware.

After the Route C plan is saved and this validation finishes, stamp `planning_at: <now>` using the current UTC instant (Timestamp rule, `actions/work-reference.md`) — only if the field is absent, because a REQ re-planned after a recovery keeps the first attempt's observation (`frontmatter set` enforces this; **Stamps are append-only**, `actions/work-reference.md`). Stamp only this successful observed event. Routes A and B omit the field.
