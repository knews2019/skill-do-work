---
id: REQ-377
title: '[impact-negligible] Stop preflight scratch showing up as untracked'
status: pending-answers
created_at: 2026-08-26T14:40:00Z
user_request: UR-074
addendum_to: REQ-374
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
---

# Stop Preflight Scratch Showing Up as Untracked

## What

`skills/do-work/tools/checks/preflight.sh` writes `do-work/working/baseline.json` and `baseline-failures.txt`, which nothing in this repository excludes, so every work run leaves untracked files behind until something removes them.

## Context

Discovered while working on REQ-374. The suite's own changelog records that these files are meant to be locally excluded alongside the orchestrator lock and its mutex files; that exclusion is installed into a consumer project, and the installer never runs against this repository, so the maintainer checkout does not have it.

They are genuinely transient — Step 5.75 writes them and Step 6.5 is their only reader — so the fix is an exclusion, never tracking them.

## Red-Green Proof

**RED prompt/case:** run `do-work run` through a Route B or C REQ, then `git status --porcelain --untracked-files=all`. `?? do-work/working/baseline.json` appears.
**Why RED now:** neither `.gitignore` nor `.git/info/exclude` names the preflight scratch files in this repository.
**GREEN when:** the same sequence reports nothing untracked, and the files are still written and still read by Step 6.5 — an exclusion, not a change to preflight's behaviour.
**Validation:** Inferred during capture.

## Open Questions
- [ ] I discovered this out-of-scope task while working on REQ-374: preflight's scratch files are untracked and unexcluded in this repo, so every work run leaves an untracked file. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

## Full Context
See `do-work/user-requests/UR-074/input.md` and REQ-374's `## Discovered Tasks`.
