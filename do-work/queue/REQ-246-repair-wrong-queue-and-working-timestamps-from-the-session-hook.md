---
id: REQ-246
title: Repair detectably wrong queue and working timestamps from the session hook
status: pending
created_at: 2026-08-18T12:38:26Z
user_request: UR-056
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-247, REQ-244, REQ-245]
batch: timestamp-stamping-integrity
---

# Repair Detectably Wrong Queue and Working Timestamps From the Session Hook

## What

A core script that scans REQ files in `do-work/queue/` and `do-work/working/` for detectably wrong `*_at` stamps and rewrites them with a mechanically derived correct value — no agent judgment anywhere in the path. Wired into the SessionStart hook the way `scripts/cleanup-req-reservations.sh` already is, so repair happens before any agent or board render sees the file.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

Detection already exists read-side (board generate/verify, forensics Check 11) but repair depended on the same agent that wrote the bad stamp. This closes that loop mechanically: a fabricated future stamp gets corrected regardless of agent compliance.

## Detailed Requirements

- **Detect only provable wrongness:** stamps later than now plus the 2-minute skew allowance (matching `actions/forensics.md` Check 11 and the board's `futureTimestampSkewAllowance` in `model.go`), and impossible orderings (`completed_at` earlier than `claimed_at`, `claimed_at` earlier than `created_at`).
- **Replacement value, by file state:** file dirty against HEAD → file mtime (the actual write instant); committed content → author time of the commit that introduced the stamp line. Clamp so the repaired set satisfies `created_at ≤ claimed_at ≤ completed_at ≤ now`.
- **Guard style of `tools/checks/record-commit-hash.sh`:** verify before replace, atomic write, a tripped guard leaves the file byte-identical, nonzero exit on any failure.
- **Audit trail:** print each correction (file, field, old value, new value, replacement source) to the hook output. No new frontmatter fields. (Log-only was the capture recommendation; the user did not override it — builder may surface it for confirmation if the hand-back warrants.)
- **Hook wiring:** invoked from the SessionStart hook (`skills/do-work/hooks/`) alongside `cleanup-req-reservations.sh`; also invocable directly.
- **Scope is queue + working only** (user-confirmed). Archived files are REQ-247's territory.

## Constraints

- Not in the board tool: its frontmatter surface is documented read-only (`frontmatter_cli.go` ~line 38) and its write-surface count is pinned in CLAUDE.md § Kanban Board Write Surfaces.
- No LLM, no heuristic text parsing beyond frontmatter field extraction.
- Platform floor applies: POSIX shell, no Go build in the repair path.
- Provenance: replaces the Discuss verdict on validate-feedback Finding 2 (UR-055 triage); Surface-cost — earned by the fabricated-stamp incident report, and cheaper than the alternative (per-site guard prose scattered across write sites, which the Timestamp rule's centralization clause forbids).

## Red-Green Proof

**RED prompt/case:** A `_dev/tests/` lock-in test builds a fixture REQ with `created_at` two hours in the future in a scratch queue, runs the repair script, and asserts the stamp is rewritten to the fixture file's mtime and the correction is logged. Fails today because the script does not exist.
**Why RED now:** Nothing repairs a detected-wrong stamp; the board only warns.
**GREEN when:** The test passes, an ordering-violation fixture is likewise repaired, a clean fixture passes through byte-identical, and `bash _dev/tests/maintainer-verify.sh` exits 0.
**Validation:** User confirmed (scope via ask tool; repair-window and mtime sources are the user's own design)

## Full Context

See `do-work/user-requests/UR-056/input.md` for complete verbatim input.

---
*Source: "basically if a wrong timestamp is detected, it should be automatically mechanically no-llm corrected by the tool" + ask-tool answer "Queue + working (Recommended)"*
