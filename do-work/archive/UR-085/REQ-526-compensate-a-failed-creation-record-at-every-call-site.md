---
id: REQ-526
title: 'Compensate a failed creation record at every create-then-record call site'
status: cancelled
created_at: 2026-09-03T01:05:00Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-457]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-457]
sweep: true
sweep_key: transaction-created-path-rollback-identity
review_generated: true
addendum_to: REQ-457
completed_at: 2026-09-03T20:40:53Z
---

# Compensate a Failed Creation Record at Every Create-Then-Record Call Site

## What

REQ-457 made the successful exclusive create the ownership event, which is correct, but it only implemented the compensating removal at the cleanup move site. Every other create-then-record call site can now publish a file and then fail to record it, leaving an object on disk that rollback never sees because it was never registered.

REQ-457 also widened `RecordCreated`'s failure surface: it now revalidates every *other* created path on each call, so a foreign swap of an unrelated created path makes `RecordCreated` fail at a site whose own create just succeeded. What used to be a near-infallible bookkeeping call is now a real failure point at twelve sites.

## Instances

- `internal/knowledgecommands/interview_commands.go:1183-1193` — reordered by REQ-457 to create then record. A `RecordCreated` failure now leaves the created file orphaned; under the old record-then-create order the same failure left nothing behind. Reviewer probe output: `b.txt ORPHANED on disk after the failed RecordCreated`.
- `internal/publication/publication_commands.go:171` and `:200`
- `internal/requeststate/state_apply.go:67`, `:82`, `:601`
- `internal/toolboxcommands/report_image.go:105`, `:221`; `note.go:45`; `architecture.go:188`; `portfolio.go:74`, `:94`
- `internal/knowledgecommands/bkb_init.go` records from the open handle, so it is already compensated; confirm rather than assume.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding F3** — `impact-rule-change` — from REQ-457's independent review (Approve, 82%). The reviewer reproduced the orphan with a scratch probe and named the twelve affected sites.
- REQ-457's Detailed Requirement 3 ("If recorder registration fails, remove only the destination just created by this process and leave the source intact") is satisfied only at the cleanup site. This REQ generalizes it.

## Detailed Requirements

- A failed creation record must leave no object this invocation just published, at every create-then-record call site — not only cleanup's move.
- The compensating removal must remove only the object this invocation created, proven by the same rooted-identity standard REQ-457 established; it must never remove a replacement.
- Prefer one shared primitive over twelve hand-written compensations. A create-and-record helper that owns both halves is the obvious shape; say so if the call sites are too dissimilar for one.
- Preserve every existing typed result, rollback status, and error wording that is asserted by current tests.

## Constraints

- Do not revert REQ-457's ordering. The successful exclusive create is the ownership event; the fix is compensation on the failure path, not an earlier record.
- Do not weaken `RecordCreated`'s revalidation to reduce its failure surface — that revalidation is REQ-457's acceptance criterion 6.

## Dependencies

Depends on REQ-457, whose ordering change creates this surface.

## Red-Green Proof

**RED prompt/case:** Force `RecordCreated` to fail after a successful create at each listed call site (a swapped unrelated created path is one real trigger) and assert no orphaned object remains.
**Why RED now:** Only the cleanup move site removes the object it just created when registration fails; the other sites return the error with the file still on disk and unregistered, so rollback cannot see it.
**GREEN when:** Every listed site leaves no unregistered object behind on a failed record, a foreign replacement at the same pathname is preserved rather than removed, and all existing typed-result assertions still hold.

## Folded From REQ-497 and REQ-524

Hand triage 2026-09-03, maintainer approved: REQ-497 and REQ-524 are cancelled and their requirements land here as acceptance criteria, so UR-085's remaining review leftovers are one REQ. They are separate root causes from this sweep's `sweep_key`; treat each as its own checklist item with its own lock-in test.

- From REQ-497 (Strictly normalize frontmatter collision identities): repository collision evidence parses frontmatter REQ ids with the selector's strict whole-value numeric grammar while filename-derived claims stay suffix-tolerant. RED today: queue files whose frontmatter values are `REQ-452` and `REQ-452junk` alias each other and explicit selection of `REQ-452` returns `DEPENDENCY-AMBIGUOUS`. GREEN: the malformed value is not merged, explicit selection picks the unique valid record in both discovery orders, and the genuine `REQ-452`/`REQ-0452` collision regressions stay green.
- From REQ-524 (Kill the owned commit process group on cancellation): cancelling a media Git transaction terminates the whole owned process group, not just the direct `git` child. RED today: `internal/toolboxcommands` `TestRemediationCancellationReachesMediaGitCommitAndRollback` fails with `media commit hook survived cancellation` (`report_image_process_test.go:85`). GREEN: that test passes with the hook process gone.

---
*Source: REQ-457 independent review finding F3.*

## Cancelled

- **When:** 2026-09-03T20:40:53Z
- **Why:** review-born below critical, unobserved hardening; under REQ-531's rule (fb0d06ca) this stays a report line. Maintainer's 2026-09-03 triage.
- **Decided by:** user, via `do-work abandon`
