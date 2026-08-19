---
id: REQ-276
title: Give record-commit-hash's readers the same fence guard its writer already has
status: pending-answers
created_at: 2026-08-18T23:38:58Z
status_changed_at: 2026-08-18T23:38:58Z
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

- The read helpers refuse an unterminated fence exactly as the writer does — ideally by **sharing one parser** rather than by growing a second copy of the guard, since the header already claims they share one.
- The header stops claiming a property the code does not have, whichever way the fix goes.
- **Check the class, not this script.** `_dev/tests/prescribed-shell-scripts-behavior.sh` and `blanked-req-scan.sh:91` already require a closing fence; `repair-req-timestamps.sh` gained the same guard in REQ-267. Grep every awk or shell frontmatter parser that ships and confirm each one gates on the closing fence; report any that does not.
- Lock-in cases for the read path, and `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

Discovered by REQ-267's builder while fixing the same defect class in `repair-req-timestamps.sh`, and independently confirmed by REQ-267's reviewer during its Restatement Sweep — which noted the repairer's new buffer-and-gate follows a pattern `record-commit-hash.sh`'s *writer* already established, making the reader's omission an inconsistency inside one file rather than a missing feature.

Classified `[normal]`, so it enters as `pending-answers` under the Discovered Tasks consent flow.

## Open Questions

- [ ] REQ-267's builder found that `record-commit-hash.sh` guards against an unterminated frontmatter fence when writing but not when reading, so `--verify` can read a body `commit:` as the frontmatter one — on the last check every REQ in the pipeline passes through. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the writer's guard means the pipeline never produces such a file itself.
