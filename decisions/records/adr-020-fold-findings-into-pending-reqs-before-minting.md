---
title: "ADR-020: Fold Findings Into Pending REQs Before Minting New Ones"
type: architecture-decision-record
status: accepted
topic_cluster: workflow-orchestration
decided: 2026-08-20
sources:
  - ai-reports/2026-08-20_1131_remaining-work-and-the-req-treadmill/index.html
  - do-work/user-requests/UR-060/input.md
  - do-work/archive/REQ-294-make-captures-impact-guard-symmetric.md
  - do-work/archive/REQ-297-report-skipped-negligible-reqs-in-targeted-mode.md
  - do-work/archive/REQ-306-reserve-new-req-creation-for-behavior-changes.md
  - do-work/archive/REQ-307-standing-prose-reconciliation-sweep.md
  - skills/do-work/actions/capture-reference.md (Fold-First Rule)
  - skills/do-work/actions/review-work.md (Step 10 Sweep consolidation)
related:
  - page: adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles
    rel: refines
  - page: adr-015-load-maintenance-crew-via-req-marker
    rel: complements
created: 2026-08-20
updated: 2026-08-20
confidence: high
---

# ADR-020: Fold Findings Into Pending REQs Before Minting New Ones

Topic cluster: [[_index_workflow-orchestration]] ([topic index](../topics/_index_workflow-orchestration.md))
See also: [[adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles]] (refines — its bounded-cycles brake was retired with the stateful pipeline), [[adr-015-load-maintenance-crew-via-req-marker]] (complements — the marker-only precedent the `sweep` marker follows)

## Context

The 2026-08-20 queue analysis measured the runaway: 145 of 298 lifetime REQs (49%) were spawned by a prior REQ's review or build, 80 REQs were second generation or deeper, and the queue refilled at the rate it drained for a month straight. The volume class was prose reconciliation — stale counts, wrong cross-reference numbers, comments describing superseded mechanisms — each minted as its own file.

Two brakes already existed and neither engaged. UR-060's `impact:` field defaulted to `impact-user-visible` on every lazy path (a pre-filled template value, the Schema Read Contract default, and a one-directional anti-invention guard), so `--skip-impact-negligible` removed one REQ from a 31-REQ queue. And Step 10's sweep consolidation filtered candidate sweeps to the discovering REQ's own UR, so a cross-cutting root cause accumulated one sweep per UR instead of one sweep.

The maintainer's directive, verbatim intent: "instead of creating a new req, check if there is another pending req that was not yet claimed to fold it into it" — and keep judging `impact:` on every REQ so negligible work stays cancellable.

## Decision

**Any flow about to mint a REQ from a finding first scans the whole queue for an existing pending, unclaimed REQ sharing the root cause and folds the finding into it; creating a new REQ file is the justified exception.** The contract lives once, in `actions/capture-reference.md` → **Fold-First Rule**, and the four minting flows (capture, review Step 10, builder follow-ups/discovered tasks, toolbox code-review) cite it — the same condition-carried pattern as the REQ Title Convention, so a new flow inherits it without a list being re-counted.

Load-bearing properties:

- **Cross-UR, whole queue.** A root cause is a property of the codebase, not of the UR that surfaced it. The old same-UR filter let sibling URs mint duplicates for one root cause; it is removed.
- **Queue residency is the unclaimed check.** Only `pending`/`pending-answers` REQs in `do-work/queue/` are fold targets; a claimed REQ lives in `working/` and stays immutable to other writers, unchanged.
- **Non-sweep pending REQs convert once.** On a root-cause match a plain pending REQ gains `sweep: true`, a `sweep_key:`, and an `## Instances` checklist seeded with its own original instance. This one-time frontmatter edit happens in the same window capture already edits queued REQs; afterwards the append-only rule governs. Instance lines carry `(found by REQ-NNN / UR-NNN)` as the cross-UR attribution.
- **Impact judgment stays per-REQ and becomes symmetric.** A folded instance inherits its sweep's `impact:`, with one escalation-only exception: an instance whose judged impact outranks the sweep's promotes the sweep's `impact:` to the instance's verdict, and an `impact-critical` instance also flips a `pending-answers` sweep to `pending` — the critical pierce applies to folds exactly as to creation, so folding can never bury an urgent finding behind a consent gate or the skip filter. A created REQ carries a judged verdict. Capture's guard now blocks unjudged assertion in both directions — asserting `impact-user-visible` because it is the default is as wrong as inventing `impact-negligible` — and the template's pre-filled verdict is commented out. The Schema Read Contract default (absence reads as `impact-user-visible`) is deliberately unchanged: absence must never be mistakable for the user's stop signal.
- **Consumer-report findings get no exemption.** What matters is whether one fix closes the class, not who found the instance.
- **A prose-only discrepancy never mints a REQ at all.** *(Amended 2026-08-20 by [[adr-021-keep-prose-only-findings-out-of-the-pipeline]]: the destination is now `do-work/prose-backlog.md`, a plain file outside the pipeline, and `standing: true` no longer exists. Everything else in this bullet stands.)* Prose-only — the fix changes no behavior, no checker's predicate, no rule's stated condition — routes to the queue's standing prose-reconciliation sweep (`standing: true`, keyed `prose-only-discrepancy-reconciliation`, never closes; this repository's instance is REQ-307) when no root-cause match exists, with three explicit exemptions: critical findings are never deferred; a contradiction between two shipped instructions changes behavior and stays first-class; user-facing artifact-contract prose is judged on its reader. Absorbed from UR-063's REQ-306, which specified this boundary independently and was completed by the same change that merged the two designs. The sweep is identified by key, created on demand when a queue lacks it (the one sweep capture may mint from scratch — a consumer install has none until its first prose-only finding), skipped by selection while empty, returned to `pending` after every drain, and excluded from its UR's terminal-resolution membership so it never holds the UR open.
- **Simplification refactors stay legitimate.** A refactor that leaves the system both simpler and more robust judges as `impact-rule-change` on its merits; the brake targets unjudged defaults and one-REQ-per-prose-discrepancy, not pragmatic quality work.

## Alternatives

1. **A single standing prose-reconciliation sweep drained on a cadence, alone** (the queue analysis's original proposal). Adopted as one component rather than the whole answer: UR-063 captured it independently as REQ-307 and it became the prose-only boundary's guaranteed destination, but fold-first is the general rule — any root cause, not only prose — so the standing sweep is where prose-only findings land when nothing more specific matches, not a second ledger format beside the queue.
2. **Drop prose-only findings entirely.** Rejected: the Restatement Sweep's findings are real (a wrong cross-reference does send readers to the wrong place); the fix is where they land, not whether they are recorded.
3. **Keep the same-UR filter and rely on `--skip-impact-negligible` alone.** Rejected: the filter governs *processing*, not *arrival*. The maintainer's stated cost is thinking about queue complexity; only fewer files reduces it.

## Consequences

The unit of queue growth changes from one file per facet to one checklist line per facet. The user approves or cancels root causes, not instances. Accepted trade-off: a discovering UR can now close while its finding pends inside another UR's sweep — the review report's "sweeps appended to" line and the instance's `(found by …)` citation are the trace, and this is judged acceptable because the finding survives either way and UR closure never claimed to mean "no known issues anywhere." The `sweep`/`sweep_key` markers gain a second writer (the one-time conversion) recorded in the schema; the board still does not parse them.

## References

- [actions/capture-reference.md](../../skills/do-work/actions/capture-reference.md) — Fold-First Rule (canonical home)
- [actions/review-work.md](../../skills/do-work/actions/review-work.md) — Step 10 Sweep consolidation, Restatement Sweep routing
- [actions/work-reference.md](../../skills/do-work/actions/work-reference.md) — schema `sweep`/`sweep_key`/`impact` rows, Discovered Tasks Classification
- [actions/capture.md](../../skills/do-work/actions/capture.md) — symmetric impact guard, fold-first duplicate check
- `ai-reports/2026-08-20_1131_remaining-work-and-the-req-treadmill/index.html` — the measurements this decision answers
