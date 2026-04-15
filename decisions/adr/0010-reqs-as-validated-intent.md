---
id: 0010
title: REQs as Validated Intent
status: accepted
decided: 2026-04-08
version: 0.51.3
topic: philosophy
supersedes: []
superseded_by: null
related:
  - adr: 0001
    rel: depends-on
  - adr: 0003
    rel: complements
---

# ADR-0010: REQs as Validated Intent

## Context

For a long time, REQs were implicitly treated as task descriptions — working documents that agents could edit, revise, or "clean up" without much ceremony. The capture action wrote a REQ, the work action read it, and the relationship between what the user asked for and what eventually got built was mostly implicit.

This made the intent trail fragile. A REQ said X on Monday and Y on Tuesday after the agent tweaked the phrasing during a capture follow-up; neither version was the "authoritative" ask. Reviewers couldn't answer "did we build what was asked?" because the ask was a moving target. Addenda would contradict their parents without anyone noticing, and the contradiction only surfaced downstream when the builder got conflicting guidance.

The skill was already producing REQs, URs, addendum chains, and review reports — but there was no *name* for what all of that was collectively producing. Without a name, it's easy to skip the validation step that makes it valuable.

## Decision

Captured REQs are **validated statements of user intent**, not drafts. The skill produces "a trail of intent, not just code," and that framing is load-bearing across the capture, verify, work, and review actions.

Specifically:

1. **Capture produces validated intent.** Open questions are resolved with the user at capture time, not silently decided by the agent. A captured REQ represents something the user has explicitly seen and agreed to.
2. **Coherence is enforced across addendum chains.** When a new addendum contradicts an existing REQ — in its parent or anywhere up the chain — the capture action flags the conflict and asks the user to resolve it before writing. Silent contradictions are forbidden.
3. **Verify checks that validation actually happened.** The verify-requests action's Internal Coherence dimension scores REQs on self-consistency; a REQ whose sections contradict each other fails verification regardless of how technically specified it looks.
4. **Work connects implementation decisions to intent.** The builder's living log records decisions and scope declarations as intent documentation — "I did X because the REQ said Y" — so the trail extends through implementation, not just capture.

## Alternatives Considered

- **REQs as mutable tasks.** Let agents freely edit REQs during capture, verify, and work. Rejected — this was the prior state, and the failure mode (moving targets, unreviewable trails) is what motivated the change.
- **Separate intent and task artifacts.** Keep a rigid "intent" file and a separate "task" file for working notes. Rejected — doubles the surface area; the existing REQ format can carry both with discipline.
- **No formal intent concept.** Treat the skill as just a task queue; let reviewers figure out intent from context. Rejected — this is what it was doing, and reviewers kept getting it wrong because the trail was implicit.

## Consequences

- **The skill has a named output.** "A trail of intent" is what every action collectively produces. Actions that violate the trail (silent REQ edits, dropped addendum chains) are recognizable as defects, not differences in style.
- **Capture asks more questions up front.** Because validation happens at capture time, the capture action can't silently guess at ambiguity. This is a real cost — capture is slower than it could be if the agent just made decisions.
- **Review has a concrete rubric.** "Did we build what was asked?" has a specific meaning: compare the built code against the validated REQ chain, including addenda.
- **Contradictions become first-class events.** Both capture and verify actively look for them. When one appears, the user resolves it; the skill does not choose.
- **Historical REQs hold their value.** Combined with [[0003-immutable-inflight-archived]], the archive preserves the full intent trail for a project. Later readers can reconstruct not just what was built but what was asked and why.

## References

- **CHANGELOG**: v0.51.3 — The Intent Trail (2026-04-08)
- **Documents**: `SKILL.md` ("Trail of Intent" blockquote)
- **Action files**: `actions/capture.md` ("Validated artifacts" philosophy, Coherence Rule), `actions/verify-requests.md` (Internal Coherence evaluation dimension), `actions/work.md` (Living log connected to intent trail)
- **Related ADRs**: [[0001-capture-execute-boundary]] (the boundary exists to create the validation window), [[0003-immutable-inflight-archived]] (immutability preserves validated intent as historical fact)
