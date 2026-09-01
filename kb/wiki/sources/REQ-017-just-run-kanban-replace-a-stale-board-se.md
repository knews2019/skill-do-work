---
title: "Lessons from REQ-017: `just run-kanban`: replace a stale board server on the port and open the default browser"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-017-just-run-kanban-replace-a-stale-board-se.md]
related:
  - page: concept-kanban-board-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-017: `just run-kanban`: replace a stale board server on the port and open the default browser

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

Running `just run-kanban` should end with the dashboard visible: (1) if a **previous queue-kanban instance** is still listening on the target port, terminate it before binding (re-running the recipe replaces the old server instead of failing with "address already in use"); (2) once the server is up, open the user's default browser at the board URL. Today the recipe only builds and serves — the user must notice the printed URL and open it by hand, and a leftover server blocks re-runs.

## Solution summary

`queue-kanban serve` now binds before announcing (fixing the false-positive "live board at…" banner on port collisions) and gained an opt-in `--open` flag with a platform-aware, test-seamed browser opener; the shipped `just run-kanban` recipe replaces a stale board instance on the target port and passes `--open`, so the command ends with the dashboard visible.

## What worked

- Parameter-injected seams — `browserOpener func(string)` and `goos` as an argument rather than reading `runtime.GOOS` inside — made every selection branch and ordering case testable without launching a browser or doing a save/restore dance on package state (D-04).
- Splitting `ListenAndServe` into `net.Listen` + `Serve(listener)` fixed a latent false-positive discovered during exploration: the old code printed the "live board at …" banner *before* binding, so a port collision announced a server that never came up.

## What didn't work

- Nothing failed hard, but D-01 was a real near-miss: two independent uncommitted changesets from concurrent sessions interleaved in one working tree. Exact-match edits (fail-loud on mismatch) were the only guard; the queue-over-worktrees discipline exists precisely to avoid this.

## Worth knowing

- The just-kanban installer is append-only — recipe behavior changes never reach already-installed projects; every recipe change must ship an upgrade note (delete block, re-run install).
- A recipe kill step needs SIGTERM + a `kill -0` poll before rebinding: a fast re-bind can race the old listener's port release (D-05).

## Back-reference

See `do-work/archive/UR-004/REQ-017-run-kanban-replace-stale-and-open-browser.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0d5c1f5`.
