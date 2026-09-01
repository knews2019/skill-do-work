---
title: "Lessons from REQ-076: Go utility emits the canonical UTC timestamp, preferred over date -u when built"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-01/REQ-076-go-utility-emits-the-canonical-utc-tim.md]
related:
  - page: concept-timestamp-and-metadata-governance
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-076: Go utility emits the canonical UTC timestamp, preferred over date -u when built

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

The shipped Go board utility already allocated the two other things the pipeline must get right and cannot guess — ticket numbers and version numbers. Timestamps were the third, still obtained by every action shelling out to `date -u +%Y-%m-%dT%H:%M:%SZ`. That command is a POSIX-ism with no Windows `cmd` equivalent, so the prescribed stamp was silently unobtainable there. The constraint that made this a REQ rather than a one-liner: the Go toolchain is optional by design, and timestamps are written by nearly every action, so making the tool the required source would promote a compiler to a hard dependency of the whole pipeline.

## Solution summary

Added a `now` subcommand that prints the current UTC instant in exactly the schema's shape and nothing else — read-only, and the only subcommand that takes no `--repo-root` because it reads a clock rather than the tree. The Timestamp rule was amended in one place with a three-option preference order: the binary if it is *already built*, else `date -u`, else PowerShell on Windows `cmd`; never build the tool to get a stamp. The maintainer-side prohibition on reaching for a compiled tool now states the exception's gate instead of contradicting a shipped instruction.

## What worked

- Testing the writer against the project's own reader (`parseTimestamp`) rather than against a hand-written expectation. That is the assertion that would actually catch a future "let's just use `time.RFC3339`" edit, because it fails on the offset form a bare shape test would accept.

## What didn't work

- The first instinct was `time.RFC3339` for both directions — it is the obvious constant and it is wrong on the emit side (numeric offset for non-UTC, sub-second digits). Caught while writing the non-UTC test, not while writing the code, which is the argument for writing that test first.

## Worth knowing

- Adding a *preferred* source to a rule does not make it used — the sites that inline the old command keep teaching it. When a canonical rule grows a preference order, the citing sites either have to stop restating the mechanism or they quietly pin the old one. That is the same failure shape REQ-075 fixed for `write_set`, arriving from the opposite direction: there the restatements carried a stale *reason*, here they carry a stale *mechanism*.
- `time.RFC3339` is right to parse with and wrong to emit with: it writes a numeric offset for a non-UTC instant and can carry sub-second digits, neither of which the schema's `*_at` shape accepts. Truncate sub-second precision rather than rounding, so a stamp can never round forward past the board's own `now` plus its skew allowance.

## Back-reference

See `do-work/archive/UR-014/REQ-076-go-utility-emits-canonical-utc-timestamp.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `3fb5938`.
