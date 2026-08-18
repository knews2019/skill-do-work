---
id: REQ-247
title: Archive timestamp audit tool driven by git commit times
status: claimed
created_at: 2026-08-18T12:38:26Z
claimed_at: 2026-08-18T18:25:40Z
route: C
user_request: UR-056
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-246]
maintenance: false
related: [REQ-246, REQ-244, REQ-245]
batch: timestamp-stamping-integrity
write_set:
- skills/do-work/scripts/audit-archive-timestamps.sh
- skills/do-work/actions/capture.md
- _dev/tests/prescribed-shell-scripts-behavior.sh
- skills/do-work/scripts/repair-req-timestamps.sh
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-08-18T18:26:21Z
  basis:
    - Route C
    - 4-file write set
    - 1 new file
    - 2 subsystems involved
    - 5 acceptance criteria
    - cross-route regression gates
    - full-suite verification
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

---

## Triage

**Route: C** - Complex

**Reasoning:** A new audit tool with git-blame-driven derivation, an ordering clamp, a same-commit immutability-rule amendment across every statement of the rule found in a sweep, and shared logic with REQ-246's script — multiple components and a co-location constraint.

**Planning:** Required

## Plan

**Approach**

1. **Share REQ-246's logic by extraction, not duplication.** `repair-req-timestamps.sh` already owns the detection predicate (future beyond 120s, impossible orderings) and the canonical-stamp comparison machinery. Factor what the auditor needs into functions the repairer keeps (or a sourced lib file if that stays ShellCheck-clean and suite-green); the auditor must not re-implement the predicate — that shared logic is why this REQ depends on REQ-246.
2. **Replacement source is git only.** For each defective stamp line in `do-work/archive/**`, derive the author time of the commit that *introduced* that line (`git blame --line-porcelain` is the proven mechanism from REQ-246's committed-file path). Never mtime — checkout resets archive mtimes, so they carry no signal. A blame that cannot answer skips with a report line, never invents.
3. **Clamp** `created_at ≤ claimed_at ≤ completed_at`, each no later than its introducing commit's time.
4. **Dry-run by default, `--fix` to write** (the REQ names this acceptable latitude; the Red-Green Proof refers to the fixing mode). Reuse the repairer's guard set on the write path: verify-before-replace, atomic rename, tripped guard leaves the file byte-identical, nonzero on failure, one audit line per correction with the sourcing commit hash.
5. **Amend the archive-immutability rule in the same commit.** `skills/do-work/actions/capture.md` § Immutability Rule (~line 63) gains a mechanical-timestamp-repair exception alongside the review-annotation precedent (`actions/review-work.md` ~line 448 wording). Sweep for any *other* statement of archive immutability in shipped text and amend those in the same commit — the co-location rule.
6. **Never hook-wired.** No `hooks/` edit anywhere; manual invocation only.
7. **Lock-in cases** in `_dev/tests/prescribed-shell-scripts-behavior.sh`: future stamp in a committed archive fixture repaired to the introducing commit's author time under `--fix`; dry-run reports and leaves bytes; ordering clamp; unreachable-blame skip. Fixtures are scratch git repos, never this repo's own archive.

**Watch for:** the class-vs-instance shape — REQ-246's review found the repairer's shape recognition narrower than the read-side detectors (REQ-255 pends on that). Do not copy that hole into the auditor: reuse the same extraction code so the two tools stay one shape-set, and say in the header that REQ-255's resolution will widen both at once.

*Written inline by the orchestrator (no separate Plan agent) — Route C.*

## Scope

**Files I will touch:**
- `skills/do-work/scripts/audit-archive-timestamps.sh` (new) — the deliberately-invoked archive auditor: git-author-time derivation, ordering clamp, dry-run default, audit trail
- `skills/do-work/actions/capture.md` (modify) — § Immutability Rule gains the mechanical-timestamp-repair exception, same commit as the tool
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — lock-in cases from the Red-Green Proof
- `skills/do-work/scripts/repair-req-timestamps.sh` (modify, only if sharing requires extraction) — shared detection/derivation logic factored rather than duplicated

**Files I will NOT touch:**
- Any hook file — this tool is never wired into a hook (REQ constraint).
- `skills/do-work-board/tools/queue-kanban/**` — read-only frontmatter decision, pinned write-surface count.
- `do-work/archive/**` itself — the tool operates on consumer archives at runtime; this repo's archive is not a test fixture (build fixtures in scratch).

**Acceptance criteria (restated from REQ):**
- [ ] Replacement values derive only from git commit author times — file mtimes never consulted.
- [ ] Same detection predicate as REQ-246, shared rather than duplicated.
- [ ] Ordering clamp holds across the repaired file: created ≤ claimed ≤ completed, each no later than its introducing commit's time.
- [ ] The archive-immutability rule is amended in the same commit that ships the tool, at every statement of the rule the sweep finds.
- [ ] Never wired into any hook; dry-run default with --fix is acceptable latitude.
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0.

## Pre-Flight

**Git:** ✓ clean
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 at the branch point (wave-1 integration tip)
**Dependencies:** ✓ Go 1.26.1, ShellCheck 0.11.0, `just`, Node, Chromium all present

*Checked by work action*
