# Builder Brief — REQ-429

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-429-complete-normalized-schema-field-projection`
Branch: `worktree-agent-REQ-429-complete-normalized-schema-field-projection`
Hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-429-handback.md`

## Ownership

Work only in the worktree above and commit on its branch. Treat every `do-work/` path in the worktree as stale and never read or write it. The only main-tree path you may write is the exact hand-back path. Do not change version or changelog files.

## Request

Complete normalized schema-field projection in the shared request record. Add normalized `caveman` value and explicit `FieldResult` evidence to both `RequestRecord` and `TypedRecord`; replace the manually incomplete assertion with a table-driven completeness ratchet covering every field governed by the Schema Read Contract; preserve generic parser evidence beside typed normalized evidence.

RED/GREEN: add `caveman: light` to the every-normalized-field fixture and require typed `lite` plus recognized evidence. The ratchet must fail when any contracted field is omitted.

Route A, `tdd: true`, `domain: backend`. Keep the repair focused and close both recorded sweep instances.

## Required rules and context

Read the general, coding-guardrails, communication-style, backend, and testing crew files, then `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` and its lessons satellite.

## Hand-back

Commit the implementation, then write the hand-back with branch/hash, approach, exact RED and GREEN evidence, every changed project file, checks with direct status, seams, and any `## Decisions` or `## Discovered Tasks`. Return only one status line after the file exists.
