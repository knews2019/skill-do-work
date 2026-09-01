---
source_type: req_lesson
req_id: REQ-247
req_path: do-work/archive/UR-056/REQ-247-archive-timestamp-audit-tool-driven-by-git-commit-times.md
date: 2026-08-18
domain: general
module: _dev/primes
tags: [general, archive, timestamp, audit, tool]
---

# Lessons from REQ-247: Archive timestamp audit tool driven by git commit times

## What the REQ was about

A deliberately-invoked audit tool that scans `do-work/archive/` for detectably wrong `*_at` stamps and repairs them, deriving every replacement from git commit times — the author time of the commit that introduced the stamp. Never run from a hook: repairing the archive is an exception to the immutability rule and stays a conscious invocation.

## Solution summary

A deliberately-invoked archive auditor, `scripts/audit-archive-timestamps.sh`, scanning `do-work/archive/**/REQ-*.md` at any depth for REQ-246's detection predicate, deriving every replacement from the introducing commit's author time (`git blame --line-porcelain`; mtime disabled in this path; unanswerable blame reports and leaves bytes untouched). Report-only by default with exit 1 on findings, `--fix` writes through the repairer's full guard set. Sharing is by **sourcing**: the repairer became a sourceable library (two pre-source switches, report-only bail, return guard), so predicate, shape recognizer, clamp, and atomic-write guards are one code body — REQ-255's widening reaches both tools in one edit. `capture.md` § Immutability Rule gained the mechanical-timestamp-repair exception in the same commit as the tool, with two co-located restatement fixes; the orchestrator applied the builder's offered seam scoping `present-work-guide.md`'s immutability restatement to its own workflow. Never hook-wired.

## What worked

Sharing by sourcing rather than a copied or third-file library — the reviewer verified by grep that the auditor holds zero predicate code, which is what turns REQ-255's future fix into a single edit that reaches both tools. The blameless lock-in that proves mtime is *never* a fallback (file left byte-identical although a valid mtime existed) is the strongest kind of negative evidence this suite has.

## What didn't work

The ordering-clamp requirement's two clauses ("clamped to the anchor" and "no later than the introducing commit's time") are unsatisfiable together when the derived commit time precedes the anchor — the REQ text carried a latent contradiction nobody caught at capture or verify. Resolved in favor of ordering with a transparent audit annotation; a REQ stating two ceilings should state which one wins.

## Worth knowing

The library switches (`timestamp_repair_apply_mode`, `timestamp_repair_git_only`) read `${var:-default}`, so exported environment variables of those names override them — inert in practice, worth knowing before renaming. A symlinked archive REQ file is silently skipped and counted clean (review Minor 1). Mixed `--fix` runs write what they can and exit 1.

## Back-reference

See `do-work/archive/UR-056/REQ-247-archive-timestamp-audit-tool-driven-by-git-commit-times.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `4035ddc`.
