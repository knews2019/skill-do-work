---
id: REQ-457
title: '[impact-critical] Record cleanup move destinations after exclusive creation'
status: pending
created_at: 2026-08-31T20:49:21Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-450, REQ-451, REQ-452, REQ-453, REQ-454, REQ-455, REQ-456]
batch: accepted-validate-feedback-root-causes
---

# Record Cleanup Move Destinations After Exclusive Creation

## What

Register a cleanup move destination as transaction-owned only after this process creates it exclusively, and before the later fallible source deletion. If transaction recording then fails, remove only the file this process just created. A losing writer must never roll back another writer's winning file.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this premature cleanup-destination ownership root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding #16 — P1 — source:** `internal/cleanup/cleanup_apply.go:296-298`

> ````text
> [P1] Record move destinations only after creating them — [prj].claude/skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go:296-298
> When two cleanup processes, or another writer, race to the same absent destination, this records the path before CreateExclusiveAt succeeds. The losing create returns EEXIST, but transaction rollback then treats the winner's file as its own creation and removes it.
> Register the destination only after this process publishes it, before later fallible source deletion. See prime-do-work-cli.md:19-20 (.claude/skills/do-work/tools/do-work-cli/prime-do-work-cli.md#L19-L20).
> ````

- **Evidence:** `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go:293-299` calls `RecordCreated` before `moveWithoutOverwrite`; the exclusive destination create occurs at lines `359-362`. Rollback deletes every recorded created path at `internal/gittransaction/git_transaction.go:463-475`. In the race, writer A publishes, writer B receives `EEXIST`, and B's rollback deletes A's file.
- **Surface-cost result:** N/A — this is a direct ownership and ordering correction at the existing exclusive-create seam.

## Detailed Requirements

- Call the exclusive destination create before registering the destination as created by this transaction.
- Register the destination immediately after successful create and before deleting the source.
- If recorder registration fails, remove only the destination just created by this process and leave the source intact.
- If exclusive creation fails with `EEXIST`, never record or remove that destination during rollback.
- Preserve no-overwrite semantics and exact transaction rollback behavior for paths the process genuinely owns.

## Constraints

- Never weaken exclusive creation to a check-then-create sequence.
- Preserve the exact-path transaction contract in `prime-do-work-cli.md`.

## Dependencies

No request prerequisite.

## Builder Guidance

Certainty level: Firm. The ownership event is successful exclusive creation; record exactly after that event.

## Red-Green Proof

**RED prompt/case:** Coordinate two cleanup writers so both observe an absent destination, writer A wins exclusive creation, and writer B loses with `EEXIST` and rolls back.
**Why RED now:** Writer B records ownership before it attempts exclusive creation, so its rollback can delete writer A's file.
**GREEN when:** Writer A's bytes remain after writer B fails and rolls back; recorder-failure coverage removes only a destination created by the same process while preserving its source.
**Validation:** User confirmed after validate-feedback accepted Finding #16.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding #16, captured by UR-085.*
