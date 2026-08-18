---
id: REQ-255
title: Give the timestamp repairer shape parity with the read-side detectors
status: pending
created_at: 2026-08-18T17:47:52Z
user_request: UR-056
addendum_to: REQ-246
domain: general
review_generated: true
sweep: true
sweep_key: repairer-detector-shape-parity
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work/scripts/repair-req-timestamps.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Give the Timestamp Repairer Shape Parity With the Read-Side Detectors

## What

`repair-req-timestamps.sh`'s hand-rolled shape recognition is strictly narrower than the read-side detectors it claims parity with (the board's `parseTimestamp` / `splitFrontmatter`), and in one shape it actively corrupts: an unquoted space-separated instant is half-rewritten into an unparseable value. Every instance below is a shape the board detects and the repairer either mangles or silently skips. For each: either repair it, or refuse it and document the refusal next to D-04's numeric-offset entry — never half-rewrite. Lock-in cases either way.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Both findings come from REQ-246's independent review, reproduced by execution against the shipped script (review I1 + I2, findings and fixtures quoted in REQ-246's archived Review section). The shared root cause is the D-01 argument — "the board flags it, so it must be repairable here" — applied to the shapes D-01 missed. **I1 first:** until the mangle is fixed, it is the one path where running the shipped SessionStart hook makes a file worse than it found it.

## Instances

- [ ] **Unquoted space-separated instant is mangled (I1, corrupting).** `created_at: 2093-01-01 00:00:00` → `created_at: 2026-08-10T12:00:00Z 00:00:00`: the extractor tokenizes at the first whitespace, the date fragment matches the date-only pattern and is "repaired", the time-of-day survives as a phantom suffix — unparseable YAML from a board-detectable input (`model.go` includes layout `2006-01-02 15:04:05`). The header's "a space may stand in for the T" fold in `comparison_key_for` is unreachable dead logic. Fix or refuse whole; never half-rewrite. The audit line must report the full old value.
- [ ] **Quoted space-separated instant is silently skipped** — truncates to an unmatched-quote token and passes through with no output.
- [ ] **CRLF-fenced files are invisible (I2).** The awk fence match `$0 == "---"` fails on `---\r`; the board's `frontmatter.go` handles CRLF explicitly. Windows agents — the likeliest source of both CRLF files and wrong local-time stamps — get detection without repair, the exact gap REQ-246 exists to close.
- [ ] **BOM-prefixed files are invisible (I2).** Same fence mismatch; the board strips the BOM.
- [ ] **Report-only riders (fold in only if touching the same lines anyway):** dead `frontmatter_value_for` helper; a lock-in grep tying `future_stamp_skew_seconds=120` to the board constant.

## Requirements

- No board-detectable shape is ever half-rewritten: each instance is either repaired to the canonical form or refused byte-identical, and every refusal is documented in the script header next to the D-04 entry.
- Lock-in cases in `_dev/tests/prescribed-shell-scripts-behavior.sh` for each instance, in whichever direction it resolves — the unquoted space-separated mangle case is mandatory.
- `bash _dev/tests/maintainer-verify.sh` exits 0.
