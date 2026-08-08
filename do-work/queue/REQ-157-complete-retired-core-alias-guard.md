---
id: REQ-157
title: "Review fix: Complete the retired core alias guard"
status: pending-answers
domain: general
created_at: 2026-08-08T20:13:29Z
user_request: UR-031
addendum_to: REQ-153
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: retired-core-alias-guard-completeness
---

# Review Fix: Complete the Retired Core Alias Guard

## What
Make the shipped-surface recurrence guard cover every former moved-command trigger, not only canonical sibling action names and a small sample of natural-language aliases.

## Context
Found during review of REQ-153. All current stale occurrences were repaired, but the guard still permits many equivalent retired core invocations, so the same contract can recur under an alias without failing distribution tests.

## Requirements
- Recover the complete former moved-command trigger set from the deleted core router/shim history and represent it as one auditable contract source.
- Reject every retired core trigger on live root/module surfaces with exact command boundaries.
- Keep the trigger list test-only or historical; do not republish legacy aliases as user guidance or runtime routes.
- Add negative controls for branding/noun phrases, generic pipeline prose, historical changelogs/archives, and explicit negative fixtures.
- Preserve current unique sibling ownership/routes, prime transition fingerprints, all 15 repaired live surfaces, and full distribution tests.

## Instances
- [ ] Guard `do-work kanban` and every former board trigger, including natural-language board phrases.
- [ ] Guard direct knowledge aliases such as `do-work recall` and every former memory/dream/interview/prompt/setup trigger.
- [ ] Guard toolbox aliases such as `do-work code review` and `do-work describe changes`, including former install targets.
- [ ] Prove the complete former trigger set is covered by a table-driven mutation fixture.

## Open Questions

- [ ] The live retired commands are repaired, but the new recurrence test covers only part of the former alias set, so an equivalent old command could return unnoticed. The cascade-depth rule requires your consent before automatically working a follow-up created by the review of another review-generated task. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
