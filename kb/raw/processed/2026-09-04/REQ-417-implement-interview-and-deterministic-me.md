---
source_type: req_lesson
req_id: REQ-417
req_path: do-work/archive/REQ-417-implement-interview-memory-commands.md
date: 2026-09-01
domain: general
module: _dev/primes
tags: [general, implement, interview, deterministic]
---

# Lessons from REQ-417: Implement interview and deterministic memory store commands

## What the REQ was about

Expose all deterministic interview and memory operations through `do-work-cli`.

## Solution summary

`do-work-cli` now exposes six deterministic Interview commands and six deterministic Memory commands through typed text/JSON results, matching source and installed flat recipes, and action delegation that retains consent and semantic judgment. Interview templates remain data; export, ingest, reset, and version operations use deterministic multi-file plans. Memory remember, forget, recall, status, bootstrap, and audit preserve the documented store formats, with tracked/private transaction separation and exact commit behavior.

## What worked

- Applying a rooted no-follow helper to mutators is insufficient when read-only status, recall, or audit paths independently enumerate and read the same configured store.
- Manual acceptance evidence cannot replace a committed regression for a `tdd: true` requirement; coverage must exercise every named engine, not only its sibling mode.
- *Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Back-reference

See `do-work/archive/REQ-417-implement-interview-memory-commands.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `ecf77a3da1751d170c22ae94b782e1354337c67b`.
