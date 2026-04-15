---
id: 0008
title: Queue Canonical Path — `do-work/queue/`
status: accepted
decided: 2026-04-10
version: 0.60.3
topic: queue-model
supersedes: []
superseded_by: null
related:
  - adr: 0002
    rel: extends
---

# ADR-0008: Queue Canonical Path — `do-work/queue/`

## Context

Pending REQ files originally lived at the root of `do-work/` — as in, `do-work/REQ-042-login-form.md` sitting alongside `do-work/user-requests/`, `do-work/working/`, and `do-work/archive/`. That layout was an artifact of the earliest versions, where there was no clear distinction between "unclaimed" and "in-flight" because there was no working folder yet.

By v0.60, the structure had all four folders (`user-requests/`, `working/`, `archive/`, plus loose REQs at the root), and the loose-at-root placement had become an unambiguous bug magnet. Agents working on the skill — human and LLM alike — kept writing `do-work/queue/` when describing pending REQ paths in action files. Version after version shipped with "queue/" references that didn't match the actual filesystem layout, leading to a recurring class of stale-path bugs (v0.60.2 and v0.60.3 both caught instances of this).

The path wasn't instinct-matching. Everyone's mental model said "the queue is a folder called queue," but on disk it was a bare root directory full of files plus three sibling folders.

## Decision

Pending REQs live at `do-work/queue/`. All queue glob patterns, directory diagrams, REQ placement paths, and `git add` staging across every action file use this canonical path.

The CLAUDE.md "Queue Path Convention" section nails it down: whenever action files reference the queue, they use `do-work/queue/` — not `do-work/` root. The capture action writes new REQs here. The work action's Step 1 scans here. The cleanup action's sweep globs here. The forensics action's stuck-work check looks here.

## Alternatives Considered

- **Keep the root and double-down on documentation.** Make it really clear in CLAUDE.md that REQs live at the root. Rejected — this had already been tried; the stale-path bugs kept coming back. The instinct to write `queue/` was stronger than the documentation.
- **Use status folders like `pending/` or `ready/`.** Alternative names that arguably match the domain better. Rejected — "queue" is what everyone was already writing. Paving the cow path only works if you pave the path people are actually walking.
- **Delete REQs from the queue after processing.** Keep only pending work at the root, moving processed REQs out. Rejected — this was already being done (via `working/`/`archive/` moves). The path name was the bug, not the lifecycle.

## Consequences

- **Instinct matches reality.** The cow path is paved. Since v0.60.3, the class of stale-queue-path bugs has dropped to near zero (later fixes in v0.60.5 caught two lingering references in `scan-ideas.md` and `deep-explore.md` and cleaned them up).
- **Clear boundary between states.** Pending (`queue/`), in-flight (`working/`), done (`archive/`), user-facing UR metadata (`user-requests/`) — each state has its own folder. Status is location.
- **Cost: one-time migration (v0.60.3).** Action files across the skill had to be updated in a single release. The migration touched 13 files but was mechanical — global search and replace on path strings.
- **Ongoing vigilance for docs.** Stale references still slip in occasionally. The `forensics` action and per-release code reviews are the main catch mechanism.

## References

- **CHANGELOG**: v0.60.3 — The Paved Path (2026-04-10); v0.60.5 — The Honest Mirror (caught two lingering references)
- **Documents**: `CLAUDE.md` (Queue Path Convention section)
- **Action files**: `actions/capture.md`, `actions/work.md`, `actions/cleanup.md`, `actions/pipeline.md`, `actions/forensics.md`, `actions/verify-requests.md`, `actions/review-work.md`, `actions/version.md`, `actions/code-review.md`, `actions/clarify.md`
- **Related ADRs**: [[0002-ur-req-pairing]] (the UR+REQ pairing's canonical location is inside `do-work/queue/` pre-claim)
