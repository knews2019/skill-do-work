# Changelog

What's new, what's better, what's different. Most recent stuff on top.

---

## 0.67.4 — The Unified Trunk (2026-04-16)

Merge of the diverging 0.67.x branches. Combines the session-state relocation and export-freshness improvements from "The Settled Tenant" / "The Right Shelf" with the update-safety and schema fixes from "The Gap Sealer".

- **All 0.67.1 features from both branches merged.**
- **0.67.2 and 0.67.3 features preserved.**
-  updated to 0.67.4.

## 0.67.3 — The Right Shelf (2026-04-16)

Moves the 0.67.2 export freshness stamp out of `exports/` and into `session.json.last_exported_at`. The sidecar-file approach would have been picked up by `ingest`'s "for each file" loop and polluted `kb/raw/inbox/` with bogus timestamp documents. Caught in review; the field-on-session.json approach was always the right one.

- `actions/interview.md`: `export` preflight and stamp-write now read/write `session.json.last_exported_at` instead of a sidecar file. Empty session shape gains the new field.
- `actions/interview-reference.md`: `session.json` schema gains `last_exported_at`. Status Vocabulary row updated with a note explaining why the stamp lives on the session, not in `exports/`. `fresh` re-run mode writes `last_exported_at: null` in the new empty session.

## 0.67.2 — The Status Ledger (2026-04-16)

Interview recipe gains a stale-export warning and a consolidated status vocabulary — small operational patches for when an operating model gets re-run in anger. Addresses gaps surfaced by a recent design review of the `work-operating-model` activation path.

- `actions/interview.md`: `export` sub-command now stamps `exports/.exported_at` after each run and does a freshness preflight on the next run — if `session.json.last_activity_at` is newer than the stamp, the user hears about it before exports are regenerated.
- `actions/interview-reference.md`: New Status Vocabulary table consolidates the four independent status fields (session `status`, layer `approved`, entry `status`, export freshness stamp) into a single reference. Explicitly notes that prior runs are archived directories, not `superseded` flags.
- `actions/interview-reference.md`: `update` re-run mode now documents the "empty a layer" path (user can nuke a layer; same approval gate applies, empty layer still counts as approved) and calls out that per-entry edit friction is intentional — the approval gate is the whole point.

## 0.67.1 — The Settled Tenant & The Gap Sealer (2026-04-16)

Combined release addressing both architectural relocation and operational safety.

- **Interview session state moved** from `./interview/<template>/` to `./do-work/interview/<template>/`. It joins `queue/`, `user-requests/`, `archive/`, and `working/` under the canonical workspace and is tracked in git.
- **Templates resolve from `<skill-root>/interviews/`**, fixing discovery when installed via `~/.claude/skills/do-work/`.
- **Widen the auto-update dirty check.** `actions/version.md` now guards `prompts/`, `interviews/`, `specs/`, `docs/`, `decisions/`, `hooks/`, `CLAUDE.md`, `AGENTS.md`, and `next-steps.md`.
- **Pre-clean discoverable dirs on update.** Removes top-level `.md` files in `prompts/` and `interviews/` during update so ghost entries don't persist.
- **Reset review state on meaningful updates.** If an interview update commits a non-zero diff, it now clears `review_completed_at` and `review_runs` to force a fresh review.
- **Layer 1 schema fixes.** Added `days` to `time_windows` and converted `interruptions` to objects in `work-operating-model.md`.
- **Session initialization.** `fresh` and `version` now write `last_activity_at: <now>` on start.

## 0.67.0 — The Open Ear (2026-04-16)

New `interview` action — a generalized elicitation framework that runs prescriptive templates to turn tacit work knowledge into agent-ready operating artifacts. First template `work-operating-model` walks the five-layer Work Operating Model (Nate B. Jones and Jonathan Edwards) across ~45 focused minutes and produces `USER.md` / `SOUL.md` / `HEARTBEAT.md` plus machine-readable exports. Session state is resumable, cross-layer contradictions get surfaced explicitly, and exports flow into BKB via `ingest` for querying.

- `actions/interview.md`: New sub-command dispatcher — `list`, `<template>`, `<template> status`, `<template> review`, `<template> export`, `<template> ingest`, `<template> reset`, `<template> versions`. Session state lives at `./interview/<template>/session.json` and persists across sessions per ADR-005. Export gates on all layers approved + at least one review pass complete. Re-run modes (`fresh`, `update`, `version`) archive prior runs as immutable `versions/v<N>-<date>/` directories.
- `actions/interview-reference.md`: Companion per ADR-001 holding the heavy content — template file format, canonical 11-field entry contract, `session.json` schema (including `review_completed_at` + `review_runs` gate fields), checkpoint format, per-export schemas for the five `work-operating-model` artifacts, re-run mode specifications, versioning scheme, and ingest frontmatter shape.
- `interviews/work-operating-model.md`: First template. Five layers — operating rhythms, recurring decisions, dependencies, institutional knowledge, friction — each with concrete prompt patterns and layer-specific `details` shape. Declares four named cross-layer contradiction checks the `review` sub-command surfaces.

---
