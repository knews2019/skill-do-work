---
source_type: req_lesson
req_id: REQ-195
req_path: do-work/archive/UR-044/REQ-195-modularize-framework-free-board-client.md
date: 2026-08-15
domain: general
module: _dev/primes
tags: [general, modularize, framework, free, queue]
---

# Lessons from REQ-195: Modularize the framework-free queue board client

## What the REQ was about

Split the framework-free queue-kanban browser client from one approximately 2,524-line `web/board.js` source unit into a private shell and eight ordered closure fragments. Preserve the existing browser runtime and every static/live behavior while making source ownership and review boundaries smaller and explicit.

## Solution summary

**[MAP CHANGED]** The queue-kanban browser runtime is still one private, framework-free classic client, but its source map is now a retained `board.js` shell plus eight ordered responsibility fragments assembled by `generate.go`. Static generation and live serving consume the same assembled bytes, making source ownership smaller without adding browser-visible loading or API seams.

## What worked

Recording the fresh source hash before cutting and comparing the first production assembly before changing factual comments isolated the migration itself from later non-runtime edits. A wildcard embed makes authored inventory observable while a separate literal manifest keeps execution order reviewable and deterministic.

## What didn't work

Counting only the full placeholder line was not equivalent to proving one raw marker token; a second marker without the canonical newline survived both replacement and the first tests. Raw-token uniqueness, canonical placement, and post-assembly absence are three separate invariants and need separate assertions.

## Worth knowing

Keep fragment files as raw statements inside the shell's existing IIFE. Separator blank lines belong to the assembler, so fragment endings and manifest joins are part of the byte contract even though the browser would tolerate many equivalent layouts.

## Back-reference

See `do-work/archive/UR-044/REQ-195-modularize-framework-free-board-client.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `675a7b7`.
