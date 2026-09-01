---
source_type: req_lesson
req_id: REQ-255
req_path: do-work/archive/UR-056/REQ-255-give-the-timestamp-repairer-shape-parity-with-the-read-side-detectors.md
date: 2026-08-18
domain: general
module: _dev/primes
tags: [general, give, timestamp, repairer, shape]
---

# Lessons from REQ-255: Give the timestamp repairer shape parity with the read-side detectors

## What the REQ was about

`repair-req-timestamps.sh`'s hand-rolled shape recognition is strictly narrower than the read-side detectors it claims parity with (the board's `parseTimestamp` / `splitFrontmatter`), and in one shape it actively corrupts: an unquoted space-separated instant is half-rewritten into an unparseable value. Every instance below is a shape the board detects and the repairer either mangles or silently skips. For each: either repair it, or refuse it and document the refusal next to D-04's numeric-offset entry — never half-rewrite. Lock-in cases either way.

## Solution summary

All six shape-parity instances fixed at the shared primitive in `repair-req-timestamps.sh` — no per-symptom patches, auditor untouched (parity arrives by sourcing). The extractor now reads the full comment-aware value and the rewrite splices by the old value's byte length carried in the plan, so the I1 mangle class is gone (space-separated instants repair whole, quoted variants land canonical-unquoted, audit lines report the full old value). CRLF fences and BOM prefixes are scanned like the board's `splitFrontmatter` with the bytes preserved through repair. Calendar-impossible values are refused byte-identical via real per-month/leap-year validation (a genuine leap-day future stamp still repairs — pinned both directions). Duplicate `_at` keys follow the last occurrence like the board's YAML dedup; shadowed lines are never touched. Riders folded: dead `frontmatter_value_for` deleted, skew-constant lock-in grep added. Header documents every refusal beside the numeric-offset entry. Suite 55 → 64 named cases, including two archive-scope parity pins proving the sourced fix reaches `audit-archive-timestamps.sh`.

## What worked

**What worked:** Fixing at the primitive rather than per symptom — one span-exact extractor/rewrite pair closed four shapes at once and, under deliberate sabotage, the pre-existing size guards still refused to write a corrupted file. Reproducing all six shapes against the shipped script *before* writing code made every RED honest.

**What didn't:** The fuzz found two shapes the six-instance list never contemplated, one of which can wedge the SessionStart hook into permanent failure. An instance list — even a six-item one assembled from two independent reviews — is still a sample; the value space is what needed enumerating.

**Worth knowing:** The board treats an unterminated fence as *no frontmatter*; the repairer scans to EOF. Padded-inside-quotes stamps are board-parseable and refused here. The archive auditor's "clean" is not yet trustworthy (REQ-268). The forensics check number cited in the script header is off by one.

## Back-reference

See `do-work/archive/UR-056/REQ-255-give-the-timestamp-repairer-shape-parity-with-the-read-side-detectors.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `84add20`.
