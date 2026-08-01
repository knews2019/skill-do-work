---
id: UR-011
title: UR ids should be accepted wherever REQ ids are
created_at: 2026-08-01T12:31:45Z
requests: [REQ-067, REQ-068]
word_count: 48
---

# UR ids should be accepted wherever REQ ids are

## Summary

`do-work run UR-059` is rejected as an unrecognized argument. The guard at `actions/work.md:101`
extracts tokens of shape `REQ-` + digits and errors on everything else — correct for a typo'd
`REG-042`, too narrow for a UR. Seven actions already take either prefix (verify-requests,
review-work, present-work, ai-report, inspect, slop-check, roadmap); only **run**, **abandon**, and
**reserve/release** reject URs.

The fix is one shared resolution contract plus its three consumers: a `UR-NNN` token resolves to its
member REQs by scanning `user_request:` frontmatter, and each action then applies its own existing
per-REQ gates to the expanded set.

## Extracted Requests

| REQ | Title | Depends on |
|---|---|---|
| REQ-067 | Target ID Resolution contract + `do-work run UR-NNN` | — |
| REQ-068 | `do-work abandon` and `do-work reserve/release` accept UR ids | REQ-067 |

## Batch Constraints

- **Membership is derived, never read from `requests:`.** The UR's `requests:` array is a
  capture-time record only (`actions/capture.md:210`); every expansion scans `user_request:`
  frontmatter across the live locations the calling action already searches.
- **One contract, three consumers.** The token shapes and the expansion rule live once in
  `actions/work-reference.md` next to the Terminal-* status sets. Each action cites it; none
  restates it. Adding a tenth independent restatement is the failure mode this batch exists to avoid.
- **A UR-expanded set honors `depends_on`** — unlike explicitly-named REQ ids, which bypass gating
  (`actions/work.md:178`). Decided with the user at capture time.
- **Expansion never widens a gate.** Each action's existing per-REQ status rules apply unchanged to
  every expanded member; a UR argument reaches more REQs, never more permissive handling of any one
  of them.
- `SKILL.md`'s word budget is enforced by `_dev/tests/contract-regressions.sh` — the routing and
  dispatch rows must absorb `UR-NNN` without growing the file materially.

## Full Verbatim Input

do-work capture-request: executing a UR should be just as valid as executing a REQ, they are the same familly, at the moment I get a warning
"""
1. do-work run UR-059 isn't a valid argument — the action takes REQ IDs. I resolved it rather than erroring.
"""

## Capture-Time Answers

Two decisions were resolved with the user during capture (recorded in the REQs as `- [x]`):

1. **Scope** → run + abandon + reserve/release (all three ID-taking actions that reject URs today).
2. **Dependency gating for a UR-expanded set** → honor `depends_on`.

---
*Captured: 2026-08-01T12:31:45Z*
