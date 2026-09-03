---
id: UR-105
title: 'Capture the maintainability-audit plan REQs and update the velocity report where speed is involved'
created_at: 2026-09-03T19:45:35Z
requests: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-555, REQ-556, REQ-557, REQ-558]
word_count: 19
---

# Capture the maintainability-audit plan REQs and update the velocity report where speed is involved

## Summary

The maintainer asked to capture "the requests": the eleven paste-ready `capture-request:` lines in `do-work/audits/audit-2026-09-03.md` §Plan (audited commit dc8a64e3, report committed at 83594c5e), handed back by the audit session moments earlier. Ten become REQs below in the plan's order (pure deletions, then consolidations, then behaviour-preserving refactors); the eleventh is an addendum to queued REQ-509 and is recorded under Folded Requests. The second half of the input (update the velocity report where speed is involved) was carried out in the capturing session, not queued: the audit's one runtime speed item (REQ-552, the archive walked twice per `audit-archive-timestamps` run) and the gate-guard observation (the launcher-thinness ratchet deleted in 0.266.9, re-homed by REQ-551) were appended to `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/index.html` as lesson L14 and a handoff paragraph.

## Extracted Requests

| REQ | Title | sweep_key | Tag | depends_on |
|---|---|---|---|---|
| REQ-549 | Drop the eight dead path tokens from the decision indexes and the lessons prime | `dead-path-pointers-in-records` | SIMPLE | — |
| REQ-550 | Collapse four exported one-line Go delegates into their targets | `exported-delegate-no-production-caller` | SIMPLE | — |
| REQ-551 | Delete the five caller-less toolbox shell shims and re-point their fixtures at do-work-cli | `toolbox-shims-no-callers` | JUDGMENT | — |
| REQ-552 | Replace two coreutils exec sites with the pure Go the package already has | `exec-where-pure-go-exists` | JUDGMENT | REQ-550 |
| REQ-553 | Source one do-work-cli launcher preamble instead of hand-rolling it in every launcher | `cli-launcher-preamble-copied` | JUDGMENT | REQ-551 |
| REQ-554 | Move the 46 lines commit.md and inspect.md share into the prescribed-shell guide | `commit-inspect-shared-body` | JUDGMENT | — |
| REQ-555 | Rewrite the prescribed-shell guide executable-homes table to the do-work-cli route form | `stale-shell-ownership-prose` | JUDGMENT | REQ-554 |
| REQ-556 | Cut the debug-artifact rule prose that do-work-cli qualify already enforces | `qualify-debug-artifact-prose-restated` | JUDGMENT | — |
| REQ-557 | Deduplicate six Go helper names defined fourteen times across do-work-cli | `per-req-duplicate-go-helpers` | JUDGMENT | REQ-550, REQ-552 |
| REQ-558 | Keep one nil-root guard in git_transaction.go and delete the other eight | `nil-root-guards-git-transaction` | JUDGMENT | REQ-557 |

## Batch Constraints

- Batch `maintainability-audit-2026-09-03`. Scope for every REQ: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- Every REQ carries its sweep_key verbatim (provenance only, not a sweep marker) and a lock-in limit pinned at today's worst value, landing in `_dev/tests/audit-lockins.sh` in the fast tier.
- Order is by dependency, not by number: REQ-552 after REQ-550; REQ-553 after REQ-551; REQ-555 after REQ-554; REQ-557 after REQ-550 and REQ-552; REQ-558 after REQ-557. The rest are independent.
- Overlaps with queued REQ-510 (work-reference.md) and REQ-509 (rationalization tables) are recorded on the REQs and the addendum; shared files are not dependencies.

## Folded Requests

- REQ-509 (merge-common-rationalizations-into-one-crew-member) — the R11 addendum: its Why clause and RED case rest on repeated rows that measure zero; restate as one loading point or cancel

## Full Verbatim Input

> ```
> capture the requests and if there is anything that is related to speed improvement also update this report ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/index.html.
> ```

---
*Captured: 2026-09-03T19:45:35Z*
