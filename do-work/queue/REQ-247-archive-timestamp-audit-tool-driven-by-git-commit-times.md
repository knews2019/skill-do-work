---
id: REQ-247
title: Archive timestamp audit tool driven by git commit times
status: pending
created_at: 2026-08-18T12:38:26Z
user_request: UR-056
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-246]
maintenance: false
related: [REQ-246, REQ-244, REQ-245]
batch: timestamp-stamping-integrity
---

# Archive Timestamp Audit Tool Driven by Git Commit Times

## What

A deliberately-invoked audit tool that scans `do-work/archive/` for detectably wrong `*_at` stamps and repairs them, deriving every replacement from git commit times — the author time of the commit that introduced the stamp. Never run from a hook: repairing the archive is an exception to the immutability rule and stays a conscious invocation.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- **Replacement source is git only** (user-specified): the author time of the commit that introduced the stamp line (`git log -L`-style lookup or equivalent). File mtimes are never consulted here — checkout resets them, so they carry no signal for committed archive content.
- **Same detection predicate as REQ-246:** future beyond the 2-minute skew, and impossible field orderings. Share the detection/derivation logic with REQ-246's script rather than duplicating it — that shared logic is why this REQ depends on REQ-246.
- **Ordering clamp** across the repaired file: `created_at ≤ claimed_at ≤ completed_at`, each no later than its introducing commit's time.
- **Amend the archive-immutability rule in the same commit** that ships the tool: `skills/do-work/actions/capture.md` § Immutability Rule (~line 63) gains a stated mechanical-timestamp-repair exception alongside the existing review-annotation exception (`actions/review-work.md` ~line 448 is the precedent wording). Co-location rule: the exception and the tool land together, and any other statement of archive immutability found in the sweep is amended in that commit too.
- **Audit trail:** print each correction (file, field, old value, new value, sourcing commit hash) to stdout. No new frontmatter fields.
- **Dry-run by default is acceptable builder latitude** (report what would change, `--fix` to write); if implemented, the RED/GREEN below refers to the fixing mode.

## Constraints

- Not in the board tool (read-only frontmatter decision, pinned write-surface count).
- Never wired into any hook. Manual invocation only — the user runs it as an audit.
- `record-commit-hash.sh` guard style: verify before replace, atomic write, tripped guard leaves the file byte-identical.
- Repairs are ordinary commits: the tool edits the working tree and reports; committing the repaired archive files follows the normal commit flow.
- Provenance: second half of the Finding 2 replacement (UR-055 triage → UR-056); requested verbatim by the user in the ask-tool answer.

## Red-Green Proof

**RED prompt/case:** A `_dev/tests/` lock-in test creates a scratch git repo with an archived fixture REQ whose `completed_at` is future-dated relative to its committing instant, runs the audit tool in fixing mode, and asserts the stamp is rewritten to the introducing commit's author time with the correction logged. Fails today because the tool does not exist.
**Why RED now:** Archived wrong stamps are permanent; the board warns on every render and nothing can correct them.
**GREEN when:** The test passes, a clean archived fixture passes through byte-identical, the immutability-rule amendment ships in the same commit as the tool, and `bash _dev/tests/maintainer-verify.sh` exits 0.
**Validation:** User confirmed ("please make an audit tool that will fix all the archive, but there it needs to take the timestamp of the git commits where it was commited")

## Full Context

See `do-work/user-requests/UR-056/input.md` for complete verbatim input.

---
*Source: ask-tool answer — "so please make an audit tool that will fix all the archive, but there it needs to take the timestamp of the git commits where it was commited"*
