---
id: UR-018
title: Parallel building across checkouts — claim anywhere, one releaser
created_at: 2026-08-04T19:44:17Z
requests: [REQ-094, REQ-095, REQ-096, REQ-097, REQ-098, REQ-099, REQ-100, REQ-101]
word_count: 420
---

# Parallel building across checkouts — claim anywhere, one releaser

## Summary

The user wants sanctioned parallel building across every instance shape — one session with parallel builders, multiple local sessions/workspaces, clones, cloud sessions — unified by one cooperative rule: a claimed/assigned REQ is visible to every other do-work instance and left alone. Decisions were made interactively across two planning sessions (plan files `~/.claude/plans/ethereal-cooking-wren.md` and `~/.claude/plans/let-s-talk-about-this-kind-robin.md` — the latter is the approved, authoritative plan). Governing philosophy chosen at every fork: **no prevention machinery anywhere; conflicts and duplicates are fixed at merge; discovery stays cheap via existing probes.**

## Extracted Requests

| REQ | Title | Phase |
|---|---|---|
| REQ-094 | Checkpoint writer label — recovery ignores foreign entries | 1 |
| REQ-095 | Two-clone acceptance run — poisoning repro + claim-conflict evidence | 1 |
| REQ-096 | Execution-model re-grain: claim anywhere, one releaser; dispatch widened to any tree | 2 |
| REQ-097 | `assigned_to` advisory field — schema, scan skip, board parse (lock-step) | 2 |
| REQ-098 | Verify probes: assigned-elsewhere-claimed-here, ur-archived-with-live-member | 2 |
| REQ-099 | Automatic wave dispatch — the loop computes the ready set | 3 |
| REQ-100 | Live auto-wave acceptance run — prove real wall-clock concurrency | 3 |
| REQ-101 | Docs + ADR — multi-checkout guide, session-ownership decision record | 4 |

## Batch Constraints

- **Do NOT build** (named failures behind each — from the approved plan): locks, leases, heartbeats, refresh intervals, staleness checks, liveness probes; auto-release/takeover on staleness (no `assigned_at`); `reserve`/`release` verbs or `reserved` status; `write_set`-based scheduling; sharded REQ-id ranges; auto `git pull`/`push` inside any action.
- The 0.163.0 forbidden-token ratchet in `_dev/tests/contract-regressions.sh` must stay green with **no weakening**; `assigned_to` must not trip its token patterns.
- `tools/queue-kanban` lock-step rule: any schema-contract change and its `model.go` mirror land in the **same commit**.
- Serial release tail stays serial: merge → qualify → review → changelog → archive one REQ at a time, by the one releaser checkout.

## Full Verbatim Input

> [Session 1, with screenshot of scenario question] that is a weird understanding of what we want to do, I want to explore ways to build things in parallel

> [Scenario answer, all four options checked plus typed:] basically all of the above are potentially correct ways to continue working, and my thinking is that if we already claim a REQ, block it or reserve it, then the other do-work instances should leave it because they should know that this is being worked on. Also local with workspaces should be a valid target too

> [Ask-tool answers, session 2:]
> - Claims: "Yes, cooperative-only"
> - Topology: "collisions are fixed by the agent as needed"
> - Integration: "depends when an agent merges for example, it fixes merge conflicts, if it is in the same git instance, then it is fixed when the feature is commited"
> - Fan-out automation: "Yes, as its own REQ"
> - Claim shape: "assigned_to field only"
> - Capture: "everybody can capture, when the branches merged, the duplicate will be fixed. Also we have a go utility that lists all duplicates, so discovery is cheap."
> - Wave shape: "Fully automatic set-picking"
> - Checkpoint: "checkpoints are transient, it's fine to commit them before changing, this way different versions of it already is available in the git history"
> - Ownership: "Claim anywhere, one releaser"
> - Recovery: "Static writer label"

> also read ~/.claude/plans/ethereal-cooking-wren.md and let's see if there is anything else that need to be cleared — before asking me, make sure to explain to me using the TTS briefing method presented before

> after everything clear use the ask tool to help me choose

Approved plan: `~/.claude/plans/let-s-talk-about-this-kind-robin.md` (copied to `do-work/user-requests/UR-018/assets/approved-plan.md`).

---
*Captured: 2026-08-04T19:44:17Z*
