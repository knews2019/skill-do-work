---
title: "ADR-021: Keep Prose-Only Findings in a Plain File Outside the Pipeline"
type: architecture-decision-record
status: accepted
topic_cluster: workflow-orchestration
decided: 2026-08-20
sources:
  - do-work/archive/REQ-307-standing-prose-reconciliation-sweep.md
  - skills/do-work/actions/capture-reference.md (Fold-First Rule, destination 3)
  - skills/do-work/tools/checks/record-commit-hash.sh
related:
  - page: adr-020-fold-findings-into-pending-reqs-before-minting
    rel: refines
created: 2026-08-20
updated: 2026-08-20
confidence: high
---

# ADR-021: Keep Prose-Only Findings in a Plain File Outside the Pipeline

Topic cluster: [[_index_workflow-orchestration]] ([topic index](../topics/_index_workflow-orchestration.md))
See also: [[adr-020-fold-findings-into-pending-reqs-before-minting]] (refines — the Fold-First Rule stands; only its destination 3 changes)

## Context

ADR-020 gave prose-only findings a destination: a permanent queue REQ carrying `standing: true`, which never closes. Within a day of shipping, the cost of that choice arrived in installments. Three branches independently had to teach a new reader that this one REQ is special — main's `0.218.0` lifecycle rework, `0.221.0`'s exit-summary section, and `0.222.0`'s cheaper-model selector, which needed its own veto so naming the sweep in a handoff would not become an unrequested drain. Two branches also designed the *same* decision independently — where a drain's commit hash goes — and git merged both cleanly, because they were prose in different paragraphs.

That fork is the sharp end. `commit:` is terminal-only provenance and a standing sweep is never terminal, so a drain could not use the field. Both branches landed on the same workaround: write the hash to a body `## Drains` line by hand. `tools/checks/record-commit-hash.sh` exists because free-form edits at exactly this step truncated six archived REQs to 0 bytes, 9 KB to 26 KB of decision trail each; its guard covers the frontmatter `commit:` line only, and its own header names a *body* `commit:` line as "precisely the confusion a body `commit:` exploits". The single write in the whole pipeline that bypassed that guard was the one the guard's comments already called out.

Counted across the tree, the mechanism reached the selection scan, the terminal-resolved status set, UR closure, `cleanup`, Step 8's substep skipping, Step 9's staging substitution and hash surface, the composed exit summary, `select-simple-reqs.sh`, and the board's Go parser, JS chips, and tests.

## Decision

**A prose-only finding is appended to `do-work/prose-backlog.md`, a plain checklist file that no part of the pipeline reads. Draining it is an ordinary REQ.** The Fold-First Rule is unchanged as a rule — scan the queue, append to a matching sweep, convert a matching plain REQ, and mint a new file only as the justified exception. Only destination 3 moves: from a never-closing REQ to a file outside the queue.

Load-bearing properties:

- **No pipeline reader learns a special case.** Nothing selects the backlog, no status tracks it, it belongs to no UR and can never hold one open. Its one reader is the work scan's queue status summary, which counts `- [ ]` lines for display.
- **A drain is an ordinary REQ, end to end.** It claims, builds, archives, and records its hash through `record-commit-hash.sh` into `commit:` like every other REQ. The unguarded body-line write is deleted rather than made safer, which is the only fix that cannot be forgotten.
- **The drain deletes the lines it resolved.** Git holds the record; the file stays a live to-do list rather than a growing ledger. Every item is re-verified at drain time, because an item recorded weeks ago may already have been fixed by an unrelated commit.
- **Opportunistic folding stays the cheapest path.** When a claimed REQ's scope already touches a file the backlog names, fixing that item there and deleting its line costs no dispatch at all — the behavior ADR-020 wanted, now with nothing to reset afterwards.
- **Visibility narrows deliberately.** The count survives on the queue status summary, printed on every run including the ones that exit with nothing claimable. It leaves the Kanban board, which renders REQs; prose debt is no longer a REQ.

## Alternatives

1. **Keep the standing sweep and extend `record-commit-hash.sh` to guard the `## Drains` line.** Rejected: it hardens the one symptom and leaves every other carve-out, each of which a future reader still has to be taught. The guard would also gain a body-line write surface it was written to forbid.
2. **Drop prose-only findings entirely.** Rejected for the same reason ADR-020 rejected it: a wrong cross-reference does send readers to the wrong place. The question is where such a finding lands, never whether it is recorded.
3. **Keep it as a REQ but let it close after each drain, re-created on the next prose finding.** Rejected: churn without a gain. Creation-on-demand already existed and the carve-outs came from the REQ-ness, not from the never-closing part alone.
4. **Wait until it bites a third time.** It had, before this was written — the `0.222.0` selector veto was the third branch inside one day. Holding longer would have kept paying the cost in installments rather than once.

## Consequences

Prose debt is visible on the queue status summary and not on the board. `standing: true` is no longer a schema field, no longer parsed, and its lock-in tests are removed with it — including the board's archived-UR live-member carve-out, whose absence restores the plain rule that any non-terminal queue member of an archived UR is an anomaly. `sweep: true` and `sweep_key:` are untouched: consolidation sweeps are ordinary REQs that do close, and nothing about them needed a carve-out.

REQ-307, the repository's own standing sweep, is archived `completed`: its two instances moved verbatim to `do-work/prose-backlog.md` after re-verification, and its purpose — giving destination 3 a home on day one — is what shipped, in a different shape.

## References

- [actions/capture-reference.md](../../skills/do-work/actions/capture-reference.md) — Fold-First Rule, destination 3 (canonical home)
- [actions/work.md](../../skills/do-work/actions/work.md) — Step 1's queue status summary line
- [tools/checks/record-commit-hash.sh](../../skills/do-work/tools/checks/record-commit-hash.sh) — the guard the deleted write bypassed
- [adr-020-fold-findings-into-pending-reqs-before-minting.md](./adr-020-fold-findings-into-pending-reqs-before-minting.md) — the rule this refines
