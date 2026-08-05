---
session_ended: 2026-08-05T20:04:00Z
last_completed: REQ-113
queue_state: 1 pending (REQ-114), 0 pending-answers, 0 blocked, 0 in-progress
reqs_processed_this_session: 6
session_depth: moderate
---

# Session Checkpoint

## In Progress (interrupted)

## Completed This Session

- REQ-110: Name the census's fully-read files so its completeness floor is explicit (Route A) — commit `66a7fcc` (review: Pass, 95%; closed UR-020)
- REQ-111: Implement the seven missing Schema Read Contract field normalizers (Route A) — v0.174.15, commit `e77383a` (review: Pass, 98%). Renumbered from 0.174.14 when main's own 0.174.14 (lockfile mirror) landed first; the merge kept both entries.
- REQ-112: Give frontmatter.go a CLI surface so prose can stop reimplementing it (Route A) — v0.175.0, commit `a6560bc` (review: Pass, 97%)
- REQ-115: status/testing_status normalization suppressed its own contract warning (Route A) — v0.175.1, commit `c976d64` (review: Pass, 98%; closed UR-023). From Codex's review of PR #130.
- REQ-113: Migrate the first prose frontmatter read site onto queue-kanban (Route A) — v0.175.2, commit `d2576a7` (review: Pass, 98%; closed UR-021)

## Still Queued

- REQ-114: `pending` — the three residual shell-logic extraction candidates (merge-aware diff, uncommitted inventory + REQ association, writer-label claim classification), restated as greps rather than line numbers. **Not approved work**: each needs its own floor decision before it becomes a change, and each candidate's grep must be re-run because the as-of-census site counts are explicitly untrusted.
- UR-021 closed to `archive/UR-021/` (3 REQs). UR-022 stays open in `user-requests/` — its only member is REQ-114.

## Session Notes

- **The census's two durable findings are now implemented.** `decisions/audits/2026-08-05-shell-logic-in-prose-census.md` was the evidence; REQ-111 added the seven missing Schema Read Contract normalizers, REQ-112 exposed the parser as `queue-kanban frontmatter get`. The audit's own §4b records that these two were the non-perishable findings and that candidates 3–5 need their greps re-run before anyone queues them.
- **Version discipline correction made this session:** the census itself got two version bumps and two changelog entries before being reverted — `decisions/` is `export-ignore`d, so nothing shipped and the changelog records delivered change only. REQ-111/112 *are* shipped changes (`tools/` ships), so they carry entries legitimately.
- **Lesson from REQ-111, worth carrying:** a RED/GREEN pair proves the stated behaviour and does not bound the change. Wiring `domain` through the new normalizer made an *absent* domain resolve to `general`, which would have given every domain-less board card a badge and a filter entry. The suite went green over it; the catch came from UNIFY asking who consumes the changed field (`grep '\.Domain'`, then `grep domain web/*.js`).
- **Lesson from REQ-112:** unit tests with injected writers prove a command's logic and nothing about its dispatch wiring — a missing `case "frontmatter"` in `main.go` would leave every test green while the binary reported `unknown subcommand`. Build the binary and exercise it for real; that step is mandatory for any future subcommand.
- REQ-110/111/112 all carry no `kb_status` (the lessons handoff was not run in this session). Offer via `do-work bkb` triage when convenient, alongside the REQ-104/108/109 backlog.
- `_dev/tests/contract-regressions.sh` has 7 pre-existing failures in its update-script probes, reproduced identically on `main`. Unrelated to this session's work and still untracked — nothing has queued them.
