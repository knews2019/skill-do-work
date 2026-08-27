---
id: REQ-377
title: '[impact-negligible] Stop preflight scratch showing up as untracked'
status: cancelled
status_changed_at: 2026-08-27T11:27:38Z
completed_at: 2026-08-27T11:27:38Z
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
- [x] I discovered this out-of-scope task while working on REQ-374: preflight's scratch files are untracked and unexcluded in this repo, so every work run leaves an untracked file. Should I process this as a new task? → Discarded: already addressed in this checkout.
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  - **Answered 2026-08-27:** The maintainer chose to discard after read-only inspection showed `.git/info/exclude` already excludes both `**/do-work/working/baseline.json` and `**/do-work/working/baseline-failures.txt`; `git check-ignore -v` confirmed both patterns. At inspection, only `baseline.json` existed and ordinary Git status was clean. This specific decision supersedes the concurrent work session's broader approval. Preserve the scratch files and preflight behavior; shared exclusions for other checkouts were not requested. Date obtained under the Timestamp rule's date-only paragraph in `skills/do-work/actions/work-reference.md`.

## Full Context
See `do-work/user-requests/UR-074/input.md` and REQ-374's `## Discovered Tasks`.

## Cancelled

- **When:** 2026-08-27T11:27:38Z
- **Why:** The maintainer discarded this during clarify because existing local exclusions already address the reported issue in this checkout. No implementation was needed or performed.
- **Decided by:** maintainer, via `do-work clarify`
