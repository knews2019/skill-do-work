---
id: REQ-475
title: '[impact-critical] Review fix: Confine all configured Memory tree readers'
status: claimed
priority: now
created_at: 2026-09-01T08:32:57Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-417]
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-417]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-417
sweep: true
sweep_key: memory-configured-tree-readers-not-rooted
claimed_at: 2026-09-03T21:43:31Z
---

# Confine All Configured Memory Tree Readers

## What

Apply one rooted, no-follow, bounded-read contract to every command that reads configured Memory state. Done means no read-only or mutating Memory surface can follow a working-file, log-file, ledger, sentinel, or directory link outside the validated memory tree, and typed output never contains bytes from a refused object.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that owns this configured Memory reader root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] Broad `memory-recall` plain-reads configured working memory and globbed daily logs, unlike the lexical path's rooted refusal.
- [ ] `memory-status` plain-reads configured working memory, daily logs, and ledger evidence.
- [ ] `memory-audit --engine memory` plain-reads configured working memory and daily logs.

## Requirements

- Reuse one repository-rooted, no-follow, regular-file reader for every configured Memory read surface, including directory enumeration.
- Reject linked or special working files, log files, ledgers, sentinels, and log directories before reading their bytes; bound accepted reads to the documented store limits.
- Ensure text and JSON findings name the refused configured path without including target bytes.
- Add adversarial fixtures for each reader and each link position, plus ordinary read-only parity and byte-preservation coverage.

## Red-Green Proof

**RED prompt/case:** Point configured working-memory and daily-log paths at synthetic files outside the memory tree, then run broad recall, status, and memory audit in text and JSON.
**Why RED now:** All three commands succeed and include outside fixture values while lexical recall rejects the equivalent linked log.
**GREEN when:** Every configured reader refuses the linked/special path before reading it, returns no fixture content, and ordinary stores retain deterministic output.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Full Context

See `do-work/runs/work-2026-08-31-165510/REQ-417-rereview.md`.

---
*Source: REQ-417 fresh re-review residual finding 1.*
