---
id: UR-013
title: Parallel builds without coordination overhead
created_at: 2026-08-03T11:41:15Z
requests: [REQ-071, REQ-072, REQ-073]
word_count: 547
---

# Parallel builds without coordination overhead

## Summary

The user wants two or more REQs built concurrently, selected manually, without reviving the
orchestrator-lock / heartbeat / co-dispatch machinery that REQ-069 deleted at v0.161.0. The
conversation started from a stale answer in another checkout (which cited the since-removed
`actions/reserve.md`) and converged, over several exchanges, on three separable requests plus two
explicit user decisions.

The user's opening safety argument — each session commits, then verifies the commit landed, and
that is the proof of work — was challenged and revised during capture: the failures that matter
(recovery re-queueing a live REQ, duplicate REQ ids, duplicate version numbers) all leave *both*
commits landing successfully, so commit-landed verification cannot detect them. The user accepted
the correction on recovery, and separately argued down the id-duplication concern as cheap to fix
after the fact — which was accepted.

## Extracted Requests

| REQ | Title | Origin in the input |
| --- | --- | --- |
| REQ-071 | Crash recovery must respect a live claim before stripping and re-queueing | "basically when a ticket is claimed, leave it claimed, if it's older then 3h then ask if it should be taken over" |
| REQ-072 | Go utility allocates REQ ids and version numbers and verifies release consistency | "the golang utility that we talked about, can also update a version number and return the new version number, we can do that with the REQ numbers as well" |
| REQ-073 | Fan-out dispatch: N concurrent builders under one queue owner | "how can I have it have parallel build, but without the anxiety and checking overhead?" |

## Batch Constraints

- **No new durable coordination state.** No locks, heartbeats, claim registries, or liveness
  counters. REQ-069 removed ~6,500 words of exactly that, and it must not return. Reading an
  existing frontmatter field and asking a human is in bounds; writing a new marker file is not.
- **No new action files and no new SKILL.md routing rows.** SKILL.md is at 2,588 of its 2,650-word
  budget.
- **`CHANGELOG.md` stays an owner-only, human-authored write** regardless of how version numbers are
  allocated — unique numbers do not make a shared prepend safe.
- **Gaps in REQ or version numbering are explicitly acceptable** (user's words: "the only problem
  that can appear that this can lead to gaps, but that's fine").
- **REQ-071 and REQ-072 must be justified on their own merits**, independent of parallel builds, so
  they do not read as concurrency machinery smuggled back in.

## User Decisions Made During Capture

Both captured via an interactive option prompt; the user's selection is recorded verbatim as the
chosen option label.

1. **"Keep builders out of it"** — asked whether recovery should be changed so it cannot destroy live
   work, or whether the second session should simply never reach that code. The user chose keeping
   builders out of Step 1 entirely: the second session owns nothing, so there is nothing to
   coordinate. (The user then *additionally* proposed the claim-respecting recovery change, which
   became REQ-071 — belt and braces, not a reversal.)
2. **"Allocate + verify, never write the changelog"** — asked how far the Go utility should reach into
   the release ritual. The user chose: it returns the next REQ number, returns and writes the next
   version, and verifies version-vs-changelog ordering and title reuse — but never writes the
   changelog body.

A third question (how builders are launched — spawned subagents vs. user-opened terminal sessions)
was withdrawn during capture rather than answered: `crew-members/background-agents.md` already
mandates a disk-durable run directory for any fan-out, and because the orchestrator synthesizes from
files rather than from conversation, the two dispatch routes are indistinguishable to it. No decision
was needed.

## Full Verbatim Input

> is this true [Image #1] didn't we just loosened the limits so it does not care if it works in parallel?

(Image #1 was a screenshot of an earlier session in a different checkout answering "which command is
safe to run in parallel?" — it cited `.claude/skills/do-work/docs/work-guide.md:89` and
`actions/reserve.md:5`, listed read-only actions as safe in parallel, and recommended
`do-work reserve REQ-NNN for <label>` plus a separate git worktree as the sanctioned two-builder
path. That answer came from an installed copy at ~0.162.x, before the 0.163.0 reservation removal.)

> how can I have it have parallel build, but without the anxiety and checking overhead?
>
> Basically my plan is to ask manually which one can be run in parallel (with our without a new workspace, is a valid variation) and then execute it without letting the other session know that it is running in parallel.
>
> This is fine because each session must commit it's changes, and then check if the commit landed correctly, that is the proof of work, and if that works, we are good.

> TTS briefing + reference appendix.
>
> Body: flowing prose, no bullets/headers/code/paths. Zero prior context:
> first mention of any ID (ticket, table, service) states what it is in
> the same sentence — "REQ 168, dropping the redundant indexes on the
> results table" — never a bare number, never a bare description. Keep
> real IDs, dates, figures in the prose so they're searchable; write them
> as they should be spoken.
>
> Appendix: marked section I skip when listening — exact paths with line
> numbers, paste-ready commands, URLs, one-line ID→description map.

> we are still talking, work should only start when things are clear.
>
> So for example ticket number, they have the title in their file name, so even if a duplicate ocasionally happens, it will be caught and fixed rather simply, without much if any extra guidance, we could also have a simple command that will ask the current repository  to check if everything ok, and based on that very cheap invocation (golang binary) any fix that needs to be done, can be performed.
>
> What other issues are there?

> "Recovery is the very first thing the pipeline does at step one, before it reads the checkpoint file and before it even looks at the queue. It finds files in the working folder and concludes, reasonably, that a previous run died mid-ticket. For a genuinely dead run that conclusion is right, and so is the response: the half-written plan and scope sections left behind by a crash are stale garbage, so they get stripped and the ticket goes back to pending for a clean re-run." <- basically when a ticket is claimed, leave it claimed, if it's older then 3h then ask if it should be taken over

> the golang utility that we talked about, can also update a version number and return the new version number, we can do that with the REQ numbers as well, the only problem that can appear that this can lead to gaps, but that's fine, the utility otherwise works in ms, and it's very much unlikely that will generate duplicates, so we just eliminated duplicates in effect.

> make sure to run do-work capture-request first to capture the intent and requests

---
*Captured: 2026-08-03T11:41:15Z*
