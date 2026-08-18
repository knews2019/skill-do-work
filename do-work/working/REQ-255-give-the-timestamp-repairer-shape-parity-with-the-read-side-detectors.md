---
id: REQ-255
title: Give the timestamp repairer shape parity with the read-side detectors
status: claimed
created_at: 2026-08-18T17:47:52Z
claimed_at: 2026-08-18T20:08:45Z
route: B
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
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-18T20:08:45Z
  basis:
    - Route B
    - 2-file write set
    - 3 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Give the Timestamp Repairer Shape Parity With the Read-Side Detectors

## What

`repair-req-timestamps.sh`'s hand-rolled shape recognition is strictly narrower than the read-side detectors it claims parity with (the board's `parseTimestamp` / `splitFrontmatter`), and in one shape it actively corrupts: an unquoted space-separated instant is half-rewritten into an unparseable value. Every instance below is a shape the board detects and the repairer either mangles or silently skips. For each: either repair it, or refuse it and document the refusal next to D-04's numeric-offset entry — never half-rewrite. Lock-in cases either way.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read prime (incl. REQ-246/244/250 lessons), crew rules, both scripts, the board's parseTimestamp/splitFrontmatter/duplicate-key recovery, and the suite. Reproduced all six shapes against the shipped script in scratch BEFORE writing code. Approach: one primitive fix — full-value comment-aware extraction, calendar-validated comparison key, span-exact rewrite, last-occurrence dedup, fence tolerance. (Transcribed from builder hand-back.)
- [x] **[APPLY]:** Five commits, each shape RED-first; scope limited to the two write-set files throughout (one premature fence-tolerance edit backed out pre-commit to preserve RED-first). (Transcribed from builder hand-back.)
- [x] **[UNIFY]:** 2 files, +313/−58, full diff reviewed (span logic consistency checked; LC_ALL=C byte arithmetic both sides); ShellCheck clean; idempotency spot-check byte-stable; maintainer-verify exit 0. (Transcribed from builder hand-back.)

## Context

Both findings come from REQ-246's independent review, reproduced by execution against the shipped script (review I1 + I2, findings and fixtures quoted in REQ-246's archived Review section). The shared root cause is the D-01 argument — "the board flags it, so it must be repairable here" — applied to the shapes D-01 missed. **I1 first:** until the mangle is fixed, it is the one path where running the shipped SessionStart hook makes a file worse than it found it.

## Instances

- [x] **Unquoted space-separated instant is mangled (I1, corrupting).** `created_at: 2093-01-01 00:00:00` → `created_at: 2026-08-10T12:00:00Z 00:00:00`: the extractor tokenizes at the first whitespace, the date fragment matches the date-only pattern and is "repaired", the time-of-day survives as a phantom suffix — unparseable YAML from a board-detectable input (`model.go` includes layout `2006-01-02 15:04:05`). The header's "a space may stand in for the T" fold in `comparison_key_for` is unreachable dead logic. Fix or refuse whole; never half-rewrite. The audit line must report the full old value.
- [x] **Quoted space-separated instant is silently skipped** — truncates to an unmatched-quote token and passes through with no output.
- [x] **CRLF-fenced files are invisible (I2).** The awk fence match `$0 == "---"` fails on `---\r`; the board's `frontmatter.go` handles CRLF explicitly. Windows agents — the likeliest source of both CRLF files and wrong local-time stamps — get detection without repair, the exact gap REQ-246 exists to close.
- [x] **BOM-prefixed files are invisible (I2).** Same fence mismatch; the board strips the BOM.
- [x] **Shape-valid but calendar-impossible values are erased instead of left for diagnosis (external review, reproduced).** `created_at: 9999-99-99T99:99:99Z` matches the glob shape, compares lexicographically later than the horizon, and gets rewritten to the file mtime — contrary to the script's guarantee that unparseable values stay untouched, and it destroys the malformed evidence the board would have diagnosed. Range-check calendar components before using a string as a comparison key. (PR #145 Codex finding, orchestrator-reproduced 2026-08-18.)
- [x] **Duplicate anchor keys read the first occurrence while the board reads the last (external review, reproduced).** With `claimed_at: <later>` followed by a duplicate `claimed_at: <earlier>`, the board's effective ordering is reversed (defect) while the repairer checks the first occurrence and reports the file clean. Use the last occurrence to match the read side, or refuse duplicate anchors without claiming clean. (PR #145 Codex finding, orchestrator-reproduced 2026-08-18.)
- [x] **Report-only riders (fold in only if touching the same lines anyway):** dead `frontmatter_value_for` helper; a lock-in grep tying `future_stamp_skew_seconds=120` to the board constant.

## Implementation Summary

**What was done:** All six shape-parity instances fixed at the shared primitive in `repair-req-timestamps.sh` — no per-symptom patches, auditor untouched (parity arrives by sourcing). The extractor now reads the full comment-aware value and the rewrite splices by the old value's byte length carried in the plan, so the I1 mangle class is gone (space-separated instants repair whole, quoted variants land canonical-unquoted, audit lines report the full old value). CRLF fences and BOM prefixes are scanned like the board's `splitFrontmatter` with the bytes preserved through repair. Calendar-impossible values are refused byte-identical via real per-month/leap-year validation (a genuine leap-day future stamp still repairs — pinned both directions). Duplicate `_at` keys follow the last occurrence like the board's YAML dedup; shadowed lines are never touched. Riders folded: dead `frontmatter_value_for` deleted, skew-constant lock-in grep added. Header documents every refusal beside the numeric-offset entry. Suite 55 → 64 named cases, including two archive-scope parity pins proving the sourced fix reaches `audit-archive-timestamps.sh`.

**Files changed (2, +313/−58):**
- `skills/do-work/scripts/repair-req-timestamps.sh` (modified) — extractor, span-exact rewrite, calendar validation, last-occurrence dedup, header docs, dead helper removed
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified) — 9 new named cases across both scan scopes

*Integrated by orchestrator from builder hand-back; merge range `2fe04c6..84add20`.*

## Decisions

Transcribed from the builder hand-back:

- **D-01:** fix at the extractor/rewrite pair, span-exact — the rewrite re-guessing a token boundary IS the mangle mechanism; carrying the measured span removes the class while every existing guard holds.
- **D-02:** space-separated instants repaired, not refused — the board parses that layout, so refusal would be detection-without-repair.
- **D-03:** calendar-impossible refused whole with real calendar arithmetic — a naive range check would still erase `2093-04-31`; refusing real leap days would open a new gap; pinned in both directions. Co-location is with REQ-246's in-script refusal text, not its non-shipping decision log.
- **D-04:** duplicate keys repaired last-occurrence, not refused — matches the read side byte-for-byte, keeps the hook non-blocking on board-renderable files; the refusal branch would fail every session until a human intervened. Surfaced for the maintainer: the dedup loop is the one place to flip if loud refusal is preferred.
- **D-05:** riders folded — both touched lines already under edit.
- **D-06:** archive-scope pins grouped into two combined cases — same coverage, less suite ballast; per-shape RED lives on the repairer side plus the old-vs-new archive contrast.

## Qualification

Passed — 2 files in merge range `2fe04c6..84add20`, all three acceptance criteria traced (no shape half-rewritten — the span-exact splice is structural; per-instance lock-ins incl. the mandatory mangle case; both scan scopes pinned per REQ-247's review requirement), P-A-U audited. The builder's D-04 judgment (repair over refusal for duplicates) is accepted: the REQ offered both branches and the chosen one keeps an unattended hook non-blocking, with the flip point named.

## Requirements

- No board-detectable shape is ever half-rewritten: each instance is either repaired to the canonical form or refused byte-identical, and every refusal is documented in the script header next to the D-04 entry.
- Lock-in cases in `_dev/tests/prescribed-shell-scripts-behavior.sh` for each instance, in whichever direction it resolves — the unquoted space-separated mangle case is mandatory.
- Each widened or refused shape gets a fixture through BOTH scan scopes (queue/working via the repairer, archive via `audit-archive-timestamps.sh`), so the shared-fix-reaches-both-tools property is pinned rather than assumed (REQ-247 review).
- `bash _dev/tests/maintainer-verify.sh` exits 0.

---

## Triage

**Route: B** - Medium

**Reasoning:** Six instances of one root cause in a shared shell library with a corrupting case first; the fix direction per shape (repair vs documented refusal) is builder judgment inside two files, with both-scan-scopes pinning required.

**Planning:** Not required

## Plan

**Planning not required** - Route B

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (modify) — shape-parity fixes in the shared library (reaches both tools by sourcing)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — lock-ins per shape, through both scan scopes

**Files I will NOT touch:**
- `skills/do-work/scripts/audit-archive-timestamps.sh` — it sources the library; parity arrives without editing it (edit only if a switch's contract must change, and then say so).
- `skills/do-work/hooks/session-start.sh` — wiring unchanged.

**Acceptance criteria (restated from REQ):**
- [ ] No board-detectable shape is ever half-rewritten: each instance repaired to canonical form or refused byte-identical with the refusal documented next to D-04's entry.
- [ ] Lock-in cases per instance; the unquoted space-separated mangle case is mandatory; each widened/refused shape pinned through BOTH scan scopes.
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0.

## Pre-Flight

**Git:** ✓ clean
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 at the branch point (0.212.17 tip)
**Dependencies:** ✓ toolchain present

*Checked by work action*
