---
id: REQ-457
title: '[impact-critical] Record cleanup move destinations after exclusive creation'
status: claimed
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
sweep: true
sweep_key: transaction-created-path-rollback-identity
claimed_at: 2026-09-02T23:27:17Z
---

# Make Rollback Ownership Follow the Created Filesystem Object

## What

Make transaction-created-path ownership identify the filesystem object created by this invocation, rather than trusting a pathname that can later resolve to another writer's object. Register cleanup move destinations only after exclusive creation, and keep create/replace/move rollback confined after parent swaps at every later mutation point.

This sweep now owns both premature cleanup-destination recording and REQ-413's post-record parent-swap rollback failure. Both share one invariant: rollback may remove only the same filesystem object this invocation created.

## Instances

- Cleanup records a move destination before exclusive creation, so a losing writer can delete the winner's file during rollback.
- Publication creates and records a repository path, then a later parent swap can redirect pathname-only rollback to an outside same-named file.
- `internal/knowledgecommands/bkb_init.go`: BKB scaffold rollback checks identity and then separately removes by pathname; Git subprocesses ignore the opened root, and incomplete writes can escape ownership recording. Final-boundary replacements can therefore be deleted or mutated despite the recorded identity. (found by REQ-416 / UR-081)

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
- Hold or revalidate rooted object identity for every created path through rollback after each later create, replace, or move mutation.
- Never follow a swapped parent outside the repository or remove a different writer's replacement object.

## Constraints

- Never weaken exclusive creation to a check-then-create sequence.
- Preserve the exact-path transaction contract in `prime-do-work-cli.md`.

## Dependencies

No request prerequisite.

## Builder Guidance

Certainty level: Firm. The ownership event is successful exclusive creation; record exactly after that event.

## Red-Green Proof

**RED prompt/case:** Coordinate two cleanup writers around exclusive creation, then separately swap a created path's parent after it is recorded and fail each later mutation index while protecting same-named outside objects.
**Why RED now:** Cleanup can record ownership before creation, and shared rollback later resolves recorded paths by pathname; both let one transaction delete an object it did not create.
**GREEN when:** Losing-writer rollback preserves the winner, post-record parent swaps never delete outside or replacement objects, and recorder-failure cleanup removes only the object created by the same invocation while preserving its source.
**Validation:** User confirmed after validate-feedback accepted Finding #16.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding #16, captured by UR-085.*
