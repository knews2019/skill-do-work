---
title: "ADR-018: Re-Grain Session Ownership to Claim Anywhere, One Releaser"
type: architecture-decision-record
status: accepted
topic_cluster: workflow-orchestration
decided: 2026-08-04
sources:
  - actions/work-reference.md (Execution Model, Worktree Dispatch Mode, Schema Read Contract, Crash Recovery, In-Progress Record)
  - actions/work.md (Step 1 scan and auto-wave, Step 2 claim, Step 10 recompute, Input flags)
  - tools/queue-kanban/ (model.go assigned_to parse, verify.go two probes, web/ badge)
  - docs/work-guide.md (Several checkouts against one queue)
  - _dev/tests/contract-regressions.sh (invariant retarget + retirement ratchet)
  - do-work/user-requests/UR-018/ (the batch, its input.md and approved plan)
related:
  - page: adr-005-pipeline-is-stateful-and-resumable
    rel: extends
  - page: adr-016-vendor-queue-kanban-into-the-skill
    rel: complements
created: 2026-08-04
updated: 2026-08-04
confidence: high
---

# ADR-018: Re-Grain Session Ownership to Claim Anywhere, One Releaser

Topic cluster: [[_index_workflow-orchestration]] ([topic index](../topics/_index_workflow-orchestration.md))
See also: [[adr-005-pipeline-is-stateful-and-resumable]] (extends), [[adr-016-vendor-queue-kanban-into-the-skill]] (complements)

## Context

**Nothing in `decisions/records/` covered session ownership at all.** The 0.161.0 exclusive-session model — the decision that deleted every piece of orchestrator-lock and parallel-dispatch machinery — was recorded only as an AI report plus its REQ (`do-work/archive/UR-012/REQ-069-exclusive-session-model-removes-concurrency-machinery.md`). So the contract that governed who may claim work had no decision record, while much smaller choices did. That gap is why this ADR exists even though its subject is a *revision*: there was no prior record to amend.

The contract being revised read, in `actions/work-reference.md`: *the pipeline supports one queue owner per checkout … **Two queue owners on one checkout stays outside the contract**, cross-session ownership with it — the pipeline does not detect, coordinate, or recover a second owner.* That bounded claiming to one checkout, and it was load-bearing: the reserve mechanism that had allocated REQs across worktrees and cloud sessions (0.125.0) was deleted at 0.163.0 as a direct consequence — collateral of the exclusive-session cleanup, **not** for defects of its own.

The user's actual working pattern outgrew it. They work the same queue from a local checkout, additional workspaces, clones, and cloud sessions, and asked for that to be sanctioned rather than merely undetected: *"if we already claim a REQ … the other do-work instances should leave it because they should know that this is being worked on. Also local with workspaces should be a valid target too."*

Two things made a naive re-grain dangerous. First, `do-work/CHECKPOINT.md` is a tracked, committed file on any install that commits `do-work/`, and it is **crash recovery's classification input** — so another checkout's live claim arrives by routine `git pull` looking exactly like local crash leftovers. That is the 2026-07-01 incident, and it is deterministic, not a race. Second, the skill has a standing ban by name on locks, leases, heartbeats, refresh intervals, staleness checks, and liveness probes — a ban earned over three patches to REQ-018, one of which reintroduced the very incident it was meant to fix. Any coordination added here had to not be that.

## Decision

**Re-grain ownership from *per checkout* to *per release tail*: any checkout may capture, claim and build; exactly one checkout runs the release tail.** The invariant `one queue owner per checkout` is retired and replaced by `one releaser per queue`, where the release tail is merge integration, the version bump, the `CHANGELOG.md` entry, archive moves, and UR closure. Two releasers against one queue, and two sessions in one working tree, both remain **unspecified** — nothing prevents them, and repair is after the fact and human (`actions/forensics.md`, `actions/cleanup.md`).

The governing philosophy, chosen by the user at every fork: **no prevention machinery anywhere; conflicts and duplicates are fixed at merge; discovery stays cheap via existing probes.** Six decisions follow from it.

1. **The claim marker is an advisory field, and the reserve verb and status stay dead.** `assigned_to: "<session>"` is one optional frontmatter field in the verbatim-read class — no alias map, no normalization, because a session name has no canonical vocabulary to normalize against. The work loop's default scan skips and reports an assigned REQ as a **courtesy, not a gate**; naming it explicitly (`do-work run REQ-NNN`) overrides the skip and clears the field as part of the claim, exactly as explicit naming already overrides `depends_on`. There is deliberately **no `assigned_at`**, no staleness threshold, and no auto-release: an assignment persists until an explicit run or a hand-edit clears it.

   The 0.163.0 forbidden-token ratchet keeps `reserve`/`release` verbs, `status: reserved`, and the `reserved_*` fields dead, unweakened. Bringing back the *field* without the *verb* is the whole shape of this decision: the verb needed a router entry (`SKILL.md` is word-budgeted and the suite enforces it), a status needed a place in the board's column vocabulary, and both invited the staleness clock the liveness ban forbids. A field needs none of that — the board badges it and nothing schedules on it.

2. **Capture happens anywhere; duplicate REQ ids are a merge-time fix.** Rejected: sharding id ranges per checkout. The detector already ships — `queue-kanban verify`'s `duplicate-req-id` probe — and the user named it as the reason: *"everybody can capture, when the branches merged, the duplicate will be fixed. Also we have a go utility that lists all duplicates, so discovery is cheap."*

3. **Crash recovery gains a static writer label, and that is not liveness machinery.** Every `## In Progress (interrupted)` entry now ends `— writer: <hostname>:<absolute-checkout-path>`. Recovery classifies four cases: own-label recovers, **foreign-label is reported and never stripped** (`claim held by <writer>, not touched`, and it never enters the three-hour takeover ladder — age adds nothing to a claim you know the holder of), label-less is own only where the checkpoint is locally uncommitted [since revised — REQ-104 dropped that case: a label-less entry is now always report-only; see Consequences], and unnamed is foreign.

   The `never grow into one` tripwire is amended rather than breached: refresh intervals, staleness checks and liveness probes stay banned **by name**, and a static label is none of them — written once at claim time from two values that never change for this checkout, never refreshed, and never read as evidence that anything is still running. It records *who wrote the entry*, which is the one question classification actually asks. Both halves of the label are load-bearing: two machines can hold the same path, and one machine can hold several checkouts.

4. **Wave dispatch is fully automatic, which supersedes "nothing computes the set".** `do-work run --fan-out [N]` puts the loop in auto-wave mode: it computes the ready set itself — pending, dependency-ready, unclaimed, not `assigned_to` another session — bounds it per `crew-members/background-agents.md`, and dispatches builders with **no confirmation gate**. This directly reverses two sentences that had been deliberate design (`actions/work.md`'s "it does not drive a fan-out wave" and the Fan-Out bullet's "a human picks … nothing computes the set"), on the user's explicit choice of *fully automatic set-picking*.

   Three things do not move. The **floor comes first**, so concurrency is an opt-in flag rather than resident behavior — the simplest agent that can read files and run shell commands must still be able to follow `actions/work.md` end to end. **Integration stays serial** at any builder count. And **`write_set` is not an input**: it stays display-only, is not read by the wave at all, and the reason is that absence of a declaration reads as *unknown* rather than *safe*, so a computed set that consulted it would look proven while under-reporting contention.

5. **The merge is the non-interference proof, and that sentence is what makes automatic set-picking safe.** A computed set asserts that its members are *runnable*, never that they do not overlap. Overlap is caught when the branches meet. This was tested in the falsifying direction: two REQs declaring the identical `write_set` both entered the wave, and the second merge refused with a content conflict.

6. **`do-work/CHECKPOINT.md` stays tracked and committed.** The user's reasoning: *"checkpoints are transient, it's fine to commit them before changing, this way different versions of it already is available in the git history."* Committing it is what makes the writer label necessary — and, it turned out, what makes a double claim detectable at all.

Both halves of the model were proven by acceptance run rather than argued. `do-work/archive/REQ-095-two-clone-acceptance-run.md` reproduced the pre-label poisoning as a silent fast-forward that erased a live claim, showed the shipped rule leaving the same claim byte-identical, and captured the real cross-checkout conflict shapes. `do-work/archive/REQ-100-live-wave-acceptance-run.md` measured two builders running concurrently for 4.109 seconds — the first recorded wall-clock overlap in this skill.

## Alternatives

1. **Keep the exclusive-session model and tell the user to work one checkout at a time.** Rejected — it does not describe how the queue is actually used, and "undetected" is not the same as "safe": the poisoning happens whether or not the contract admits the second checkout exists. Leaving it unsaid meant a second checkout's agent, reading the contract honestly, would refuse to claim while still being able to destroy a claim by syncing.
2. **Restore the reserve verb and `reserved` status as they were at 0.125.0.** Rejected. Nothing was wrong with them functionally — they were deleted as collateral — but restoring them costs a router entry against an enforced word budget, a status in the board's column vocabulary, and, historically, a staleness clock: every version of reserve acquired one, and that is the machinery the liveness ban exists to keep out. The advisory field delivers the cooperative signal without any of the three.
3. **Add real coordination: a lock file, a lease, or a heartbeat in `do-work/`.** Rejected by standing ban and by history — REQ-018 took three patches, and one of them reintroduced the incident. A lock in a *git-synced* directory is worse than none: it arrives stale by construction, and an agent that trusts it is more dangerous than one that knows it has nothing.
4. **Shard REQ id ranges per checkout to prevent duplicate captures.** Rejected — prevention machinery for a problem whose detector already ships, and it makes every checkout's captures visibly non-portable. Fix-at-merge plus `duplicate-req-id` costs nothing and stays correct when someone captures from an unplanned fifth checkout.
5. **Keep the human confirmation gate on the wave (compute the set, then ask).** Rejected on the user's explicit choice. Worth recording that the gate was the safer-feeling option and was declined for a defensible reason rather than convenience: the confirmation would have been asked to approve *runnability*, which the computation already establishes, while the thing a human might actually catch — overlap — is caught by the merge either way. A prompt that cannot change the outcome is friction.
6. **Make `write_set` a scheduling input so the wave avoids overlap.** Rejected, and it is the most tempting of these. Absence of a `write_set` reads as unknown, not safe, and the board's overlaps badge already misses glob-vs-glob, `**`, and directory entries — so a wave that scheduled on it would produce sets that *look* proven. Display-only at any builder count, with the merge as the proof, is the honest arrangement.
7. **Auto-wave on by default whenever the harness supports worktrees.** Rejected — the floor. An implicit trigger makes one `do-work run` behave differently on two harnesses, and a reader of the simplest path would need to understand concurrency to know which they were getting.

## Consequences

A user can now point a workspace, a clone, or a cloud session at one queue and have every instance cooperate: claims are visible, earmarks are honored as a courtesy, and collisions surface as ordinary git conflicts with a documented resolution (**keep every checkpoint entry from both sides** — both one-sided resolves lose a real claim). One checkout still owns the release tail, so version numbers and changelog entries stay serial and uncontended. `do-work run --fan-out N` turns build-phase wall-clock into a knob, with integration unchanged behind it.

Costs and open edges, stated rather than implied. **The unspecified cases are genuinely unspecified** — two releasers, or two sessions in one working tree, will produce damage nothing prevents and only `forensics`/`cleanup` repair. **The checkpoint is now a guaranteed conflict point** on every concurrent claim, including two that overlap in nothing, because two single-line appends land at the same position; the resolution is trivial but must be done right, and a one-sided resolve is exactly the poisoning done by hand. **The label-less recovery bullet shipped known-unsound** under this model — a merge-resolved checkpoint is "locally modified" for reasons unrelated to authorship, which re-opened the strip for entries written before 0.170.0; filed as REQ-104 rather than fixed inside the run that found it, and **since resolved by REQ-104**, which dropped the heuristic outright: a label-less entry is now always report-only, and reclaiming a genuinely-own one is a human act. **Agent behavior under concurrency remains unproven**: REQ-085 exercised real agent builders without overlap, REQ-100 measured real overlap with scripted builders, and no run yet covers both. And on installs where `do-work/` is untracked — the common case — none of the cross-checkout machinery has anything to travel on: no claim syncs, so neither the poisoning nor the fix-at-merge detection exists there, and `duplicate-req-id` is the only detector left.

Two ratchets were tightened rather than loosened while doing this. The retired wording `one queue owner per checkout` now fails the suite if it reappears anywhere in `actions/`, `docs/`, or `SKILL.md` — the same protection its own predecessor got. And `queue-kanban verify` gained two probes for drift this model makes reachable: `assigned-elsewhere-claimed-here` and `ur-archived-with-live-member`, neither marked mechanically fixable, because both resolutions are human calls.

## References

- [actions/work-reference.md](../../actions/work-reference.md) — Execution Model — Claim Anywhere, One Releaser; Worktree Dispatch Mode (builder-tree widening, Auto-wave, claim-conflict resolution); Crash Recovery; In-Progress Record; the `assigned_to` schema line
- [actions/work.md](../../actions/work.md) — Step 1's assigned-elsewhere skip and auto-wave set selection, Step 2's clear-on-claim, Step 10's recompute, the `--fan-out` flag
- [docs/work-guide.md](../../docs/work-guide.md) — "Several checkouts against one queue", the user-facing walkthrough
- [tools/queue-kanban/model.go](../../tools/queue-kanban/model.go) — the display-only `assigned_to` parse, in lock-step with the Schema Read Contract
- [tools/queue-kanban/verify.go](../../tools/queue-kanban/verify.go) — the two new probes
- [do-work/archive/UR-012/REQ-069-exclusive-session-model-removes-concurrency-machinery.md](../../do-work/archive/UR-012/REQ-069-exclusive-session-model-removes-concurrency-machinery.md) — the decision this partially reverses, and the only prior record of it
- [do-work/user-requests/UR-018/input.md](../../do-work/user-requests/UR-018/input.md) — the verbatim user decisions
- [do-work/user-requests/UR-018/assets/approved-plan.md](../../do-work/user-requests/UR-018/assets/approved-plan.md) — the approved plan, including the do-not-build list
- [do-work/archive/REQ-095-two-clone-acceptance-run.md](../../do-work/archive/REQ-095-two-clone-acceptance-run.md) — the poisoning repro and the cross-checkout conflict evidence
- [do-work/archive/REQ-100-live-wave-acceptance-run.md](../../do-work/archive/REQ-100-live-wave-acceptance-run.md) — the measured wall-clock concurrency run
- [[adr-005-pipeline-is-stateful-and-resumable]] — the resumability model whose checkpoint this builds on
- [[adr-016-vendor-queue-kanban-into-the-skill]] — why the board's parser and the schema it tracks live in one repo, which is what let `assigned_to` and its parse land together
