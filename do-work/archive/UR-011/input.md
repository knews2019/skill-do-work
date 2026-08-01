---
id: UR-011
title: UR ids should be accepted wherever REQ ids are
created_at: 2026-08-01T12:31:45Z
requests: [REQ-067, REQ-068, REQ-070]
word_count: 48
---

# UR ids should be accepted wherever REQ ids are

## Summary

`do-work run UR-059` is rejected as an unrecognized argument. The guard at `actions/work.md:101`
extracts tokens of shape `REQ-` + digits and errors on everything else — correct for a typo'd
`REG-042`, too narrow for a UR. Seven actions already take either prefix (verify-requests,
review-work, present-work, ai-report, inspect, slop-check, roadmap); only **run**, **abandon**, and
**reserve/release** reject URs.

The fix is one shared resolution contract plus its consumers: a `UR-NNN` token resolves to its
member REQs by scanning `user_request:` frontmatter, and each action then applies its own existing
per-REQ gates to the expanded set.

The asymmetry runs both ways. `actions/roadmap.md` is the mirror defect — it accepts `UR-NNN` and not
`REQ-NNN`, and silently surveys the whole queue on an unrecognized token. Surfaced by
`do-work verify-request`; the user chose full symmetry, so REQ-070 covers it.

## Extracted Requests

| REQ | Title | Depends on |
|---|---|---|
| REQ-067 | Target ID Resolution contract + `do-work run UR-NNN` | — |
| REQ-068 | `do-work abandon` and `do-work reserve/release` accept UR ids | REQ-067 |
| REQ-070 | `do-work roadmap` accepts REQ ids (the inverse asymmetry) | REQ-067 |

## Batch Constraints

- **Membership is derived, never read from `requests:`.** The UR's `requests:` array is a
  capture-time record only (`actions/capture.md:210`); every expansion scans `user_request:`
  frontmatter across the live locations the calling action already searches.
- **One contract, every consumer.** The token shapes and the expansion rule live once in
  `actions/work-reference.md` next to the Terminal-* status sets. Each action cites it; none
  restates it. Adding a tenth independent restatement is the failure mode this batch exists to avoid.
- **Both prefixes everywhere, no exceptions.** After this batch, every id-taking action accepts
  `REQ-NNN` and `UR-NNN`. Actions that take no id at all (`clarify`, `board`) are a different design,
  not an asymmetry — explicit non-goal.
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

Decisions resolved with the user (recorded in the REQs as `- [x]`):

1. **Scope** → run + abandon + reserve/release (all three ID-taking actions that reject URs today).
2. **Dependency gating for a UR-expanded set** → honor `depends_on`. A `--force` bypass was offered
   and declined.
3. *(verify)* **The inverse asymmetry** → full symmetry; `roadmap` gains `REQ-NNN` scoping (REQ-070)
   rather than only losing its silent unrecognized-token default.

## Verification

`do-work verify-request` run 2026-08-01, overall confidence **88% → 95%** after fixes.

- **Important, fixed:** REQ-067's RED case asserted a hard `Unrecognized argument(s)` stop, but the
  verbatim input reports the opposite — the agent warned and *"resolved it rather than erroring."*
  The RED now names the undefined, reader-dependent behavior and explicitly warns the builder not to
  read an improvised resolution as GREEN.
- **Minor, fixed:** no success-path output was specified; targeted mode now announces
  `UR-NNN → REQ-…` before claiming.
- **Nit, fixed:** REQ-068's RED described abandon's glob imprecisely.
- **Ambiguous, resolved:** the `roadmap` inverse asymmetry → REQ-070.
- No prompt-injection content in the verbatim input; the triple-quoted block is a prior agent's
  status line, treated as data.

---
*Captured: 2026-08-01T12:31:45Z*
