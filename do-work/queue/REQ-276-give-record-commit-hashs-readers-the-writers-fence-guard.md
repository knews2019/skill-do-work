---
id: REQ-276
title: Give record-commit-hash's readers the same fence guard its writer already has
status: pending
created_at: 2026-08-18T23:38:58Z
status_changed_at: 2026-08-19T13:45:20Z
user_request: UR-056
addendum_to: REQ-267
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work/tools/checks/record-commit-hash.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Give record-commit-hash's Readers the Same Fence Guard Its Writer Already Has

## What

`skills/do-work/tools/checks/record-commit-hash.sh` guards against an unterminated frontmatter fence **on the write path only**. Its rewrite awk (`:472-487`) buffers and gates on `frontmatter_closed`. Its **read** helpers (`:108-122`, called at `:247` and `:318-338`) still scan to end of file — so on a file whose opening `---` is never closed, `--verify` can read a body `commit:` line as the frontmatter one.

The script's own header at `:103-106` claims both halves use one parser. They do not.

**Why this is worth more than its severity suggests:** `record-commit-hash.sh --verify` is the last check every REQ in this pipeline passes through, and it exists because free-form frontmatter edits once truncated six archived REQs to zero bytes in a consumer repo. A guard that is asymmetric between its reader and its writer is exactly the shape that lets a corrupted file report as verified.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

Scope was narrowed by the user during clarify (see `## Decision Record` below). Refuse the file at
the door rather than teaching each reader a guard.

- Add one `require_closed_frontmatter` helper — an awk that fails when a file opens `---` on line 1
  and never closes it — and call it at exactly **two** sites:
  - on `$request_file` at startup, beside the existing CRLF refusal (`:310`), which covers every
    reader call on that file (`:318`, `:322`, `:328`, `:336-338`, `:566-570`) without editing any of them;
  - on the parent blob fetched by `--verify` (`:247`), the only reader input that is not `$request_file`.
- The header comment at `:103-106` stops claiming the readers and the writer share one guard, and says
  what is actually true: an unterminated fence is refused up front, so no reader ever sees one.
- Two lock-in cases: a file whose fence never closes and whose body carries a column-0 `commit:` line is
  refused on the write path and on the `--verify` path. Each names the failure it pins.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

**Explicitly out of scope** — do not do these under this REQ:

- Refactoring the three readers to share one parser with the writer. A single precondition is smaller
  than threading a guard through each awk, and buys the same property.
- The repo-wide sweep of every shipped awk or shell frontmatter parser. Different scope, different
  write set; the user declined to capture it as a separate REQ during this clarify.

## Context

Discovered by REQ-267's builder while fixing the same defect class in `repair-req-timestamps.sh`, and independently confirmed by REQ-267's reviewer during its Restatement Sweep — which noted the repairer's new buffer-and-gate follows a pattern `record-commit-hash.sh`'s *writer* already established, making the reader's omission an inconsistency inside one file rather than a missing feature.

Classified `[normal]`, so it enters as `pending-answers` under the Discovered Tasks consent flow.

## Open Questions

- [x] REQ-267's builder found that `record-commit-hash.sh` guards against an unterminated frontmatter fence when writing but not when reading, so `--verify` can read a body `commit:` as the frontmatter one — on the last check every REQ in the pipeline passes through. Should I process this as a new task? → Yes, add to queue — with the scope simplified to one up-front guard (see the Decision Record).
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the writer's guard means the pipeline never produces such a file itself.

## Decision Record

- **[2026-08-19] Scope narrowed to one up-front refusal.** The user reviewed this REQ during
  `do-work clarify`, judged the original requirements over-complicated, and chose the simplified
  form now in `## Requirements`. Reasoning: the three readers already run the same four-line awk;
  the defect is that none of them gates on the closing fence, and a single precondition checked once
  at startup fixes every reader at once. Sharing one parser with the writer would be a larger diff
  for the same property.
- **Put out of scope by the same decision:** the reader/writer parser-sharing refactor, and the
  repo-wide sweep of other shipped frontmatter parsers. The user was offered the sweep as its own
  REQ and declined it, so it is not captured anywhere — reopen deliberately if wanted.
- **Decided by:** user, via `do-work clarify`.
