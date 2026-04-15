---
id: 0005
title: Subagent Dispatch with Direct-Read Fallback
status: accepted
decided: 2026-02-24
version: 0.11.0
topic: platform-portability
supersedes: []
superseded_by: null
related:
  - adr: 0004
    rel: extends
---

# ADR-0005: Subagent Dispatch with Direct-Read Fallback

## Context

By v0.11.0 the action files had grown substantial — combined, the corpus exceeded 170KB. When `SKILL.md`'s routing dispatched an action by inlining the full action-file content into the main thread, two things happened:

1. **The main thread drowned.** 200KB of workflow prose flooded the conversation's context window, crowding out the user's actual problem and pushing important context out.
2. **Response latency grew visibly.** Every routed action incurred the cost of loading everything, whether or not the user wanted a running commentary.

The obvious fix was dispatching actions to subagents — fresh context, clean return value, main thread stays lean. But that only works on tools that *have* subagents. Claude Code has the Task tool. A bare chat interface does not. Requiring subagents would have violated [[0004-platform-agnostic-action-files]].

## Decision

`SKILL.md`'s "Action Dispatch" section describes a **two-tier pattern**:

1. **If subagents are available**, dispatch the action to a general-purpose subagent with a prompt that tells it to read the action file and execute it. Long-running actions (`work`, `cleanup`) dispatch as background subagents; short actions (`capture`, `verify`, `version`) dispatch as foreground.
2. **If subagents are not available**, the main thread reads the action file directly and executes inline. The action file is written so this still works — see [[0004-platform-agnostic-action-files]] for the content rules that guarantee it.

The dispatch table in `SKILL.md` identifies which actions run in background vs. foreground (when subagents exist), but the fallback to direct-read is available to every action.

## Alternatives Considered

- **Always load in the main thread.** The simple path. Rejected — it's what we were doing, and the context flooding it caused is exactly what motivated the ADR.
- **Require subagents.** Only run on tools that have them. Rejected — violates the platform-portability contract.
- **Split `SKILL.md` per action.** Give each action its own top-level file; the router just hands off by reference. Rejected — complicates routing without addressing the context-size problem for tools that lack subagents anyway; the direct-read fallback still has to exist somewhere.

## Consequences

- **Main thread stays clean on tools with subagents.** The ~200KB of action-file content stays out of the main conversation's context, letting the user's thread focus on what they care about.
- **The skill still works on tools without subagents.** Bare chat interfaces follow the same action-file content, just in-thread.
- **Foreground vs. background is a performance decision, not a correctness one.** Action files don't rely on whether they're dispatched in foreground or background — they produce the same result either way.
- **Subagents don't share memory with the main thread.** Each dispatched action gets a fresh context, so prompts must be self-contained. This is a constraint on how `SKILL.md` writes dispatch prompts — they must pass along whatever the action needs to know.
- **The pipeline action overrides the default.** All pipeline-dispatched actions run foreground (blocking) to avoid races between sequential steps. This exception is documented in both `SKILL.md` and `actions/pipeline.md`.

## References

- **CHANGELOG**: v0.11.0 — The Delegate (2026-02-24), v0.11.1 — The Soft Landing (fallback path for subagent-less environments), v0.51.6 — The Narrow Pipe (pipeline's foreground override)
- **Action files**: `SKILL.md` (Action Dispatch section), `actions/pipeline.md` (foreground exception)
- **Related ADRs**: [[0004-platform-agnostic-action-files]] (direct-read fallback only works because action files were written to support it)
