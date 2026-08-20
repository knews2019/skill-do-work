---
id: REQ-300
title: "Resync the text that still plans around the pre-split shell behavior suite"
status: pending
created_at: 2026-08-20T08:37:00Z
status_changed_at: 2026-08-20T08:37:00Z
user_request: UR-056
addendum_to: REQ-258
domain: general
review_generated: true
sweep: true
impact: impact-user-visible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- do-work/RESTART-PROMPT.md
- do-work/queue/REQ-263-tighten-qualifys-ownership-probe-and-warn-legibility.md
- do-work/queue/REQ-264-make-a-disarmed-p-a-u-audit-visible.md
- do-work/queue/REQ-271-make-the-layout-lock-step-see-every-spelling.md
- decisions/audits/2026-08-11-defensive-surface.md
---

# Resync the Text That Still Plans Around the Pre-Split Shell Behavior Suite

## What

REQ-258 split `_dev/tests/prescribed-shell-scripts-behavior.sh` into one case file per script under `_dev/tests/prescribed-shell-cases/`. Every *code* consumer was fine — there is one caller and it reads only the exit status. What went stale is the text that **plans around** the old single-file layout, and it is all the same root cause, so it is one sweep rather than four REQs.

The live consequence: the restart prompt tells the next session that this file is a scheduling bottleneck forcing four more waves. It is not one any more. REQ-263, REQ-264 and REQ-271 now write three different case files (`qualify.sh`, `qualify.sh`, `repair-req-timestamps.sh` respectively — so 263 and 264 still overlap each other, but neither overlaps 271), which changes how they can be batched.

## Instances

- [ ] **`do-work/RESTART-PROMPT.md:33`** — states `_dev/tests/prescribed-shell-scripts-behavior.sh` "is written by REQ-258, 263, 264, 268 and 271 — at most one per wave" and that REQ-258 will dissolve it. REQ-258 has shipped. Rewrite to the post-split reality: the runner is no longer written by case-adding REQs, REQ-263/264 share `prescribed-shell-cases/qualify.sh`, REQ-271 writes `prescribed-shell-cases/repair-req-timestamps.sh`, and 271 is therefore disjoint from both.
- [ ] **`do-work/queue/REQ-263` `write_set`** — names the runner; should name `_dev/tests/prescribed-shell-cases/qualify.sh`.
- [ ] **`do-work/queue/REQ-264` `write_set`** — same.
- [ ] **`do-work/queue/REQ-271` `write_set`** — names the runner; should name `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh`. Its `## Red-Green Proof` command (`bash _dev/tests/prescribed-shell-scripts-behavior.sh`) is still correct and must not be changed; only the observed case count in that paragraph is stale, and it was stale before REQ-258.
- [ ] **`decisions/audits/2026-08-11-defensive-surface.md`, eleven Coverage rows** — each cites "`_dev/tests/prescribed-shell-scripts-behavior.sh` <name> case", and the named case now lives in `_dev/tests/prescribed-shell-cases/<script>.sh`. **Decide before editing:** this is a dated decision record, and REQ-234 set the precedent that dated history is left alone. The pointer is one hop from the case either way. If the call is to leave it, tick this box with that reasoning rather than editing.

## Context

From REQ-258's review, Important finding 1 and Minor finding 2, consolidated by root cause. `write_set` is display-only and Step 5.5 overwrites it from the fresh Scope declaration at claim time, so **no build is at risk** — the stale values mislead the board's overlaps badge and a human planning waves, which is exactly what `RESTART-PROMPT.md` exists for.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
