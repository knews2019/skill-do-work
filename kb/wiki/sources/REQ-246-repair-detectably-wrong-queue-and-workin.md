---
title: "Lessons from REQ-246: Repair detectably wrong queue and working timestamps from the session hook"
type: source-summary
topic_cluster: checkpoint-and-crash-recovery
sources: [raw/processed/2026-09-01/REQ-246-repair-detectably-wrong-queue-and-workin.md]
related:
  - page: REQ-247-archive-timestamp-audit-tool-driven-by-g
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-246: Repair detectably wrong queue and working timestamps from the session hook

Part of the [[concept-session-checkpoints-and-recovery]] cluster.

## What the REQ was about

A core script that scans REQ files in `do-work/queue/` and `do-work/working/` for detectably wrong `*_at` stamps and rewrites them with a mechanically derived correct value — no agent judgment anywhere in the path. Wired into the SessionStart hook the way `scripts/cleanup-req-reservations.sh` already is, so repair happens before any agent or board render sees the file.

## Solution summary

Built `skills/do-work/scripts/repair-req-timestamps.sh`, a POSIX-floor mechanical repairer for detectably wrong `*_at` stamps in `do-work/queue/` and `do-work/working/`: future stamps beyond the shared 2-minute skew allowance (any top-level `*_at` key — the suffix is the rule, not a field list) and impossible orderings among `created_at`/`claimed_at`/`completed_at`. Replacement source is file mtime when the file is dirty against HEAD (or untracked / no git), otherwise the introducing commit's author time via `git blame --line-porcelain`; derived values are clamped to `created_at ≤ claimed_at ≤ completed_at ≤ now`. Guard style follows `tools/checks/record-commit-hash.sh` (verify-before-replace, atomic same-directory rename, tripped guard leaves the file byte-identical, nonzero exit, one audit line per correction). Wired into `hooks/session-start.sh` alongside the reservation cleanup, presence-guarded so a partial install still prints the banner; also directly invocable.

## What worked

Cloning the guard architecture of `record-commit-hash.sh` wholesale (verify-before-replace, atomic rename, byte-identical on trip) survived every adversarial mutant the review threw at it — including a fabricated future mtime. Small commits (script+RED, GREEN, hook wiring) made two transport-level interruptions nearly free to resume.

## What didn't work

The repair-side parser was hand-rolled instead of derived from the read-side detectors it claims parity with — so it recognizes strictly fewer shapes than the board (space-separated instants, CRLF/BOM fences), and in one shape (unquoted space-separated) it half-rewrites and corrupts. The session's standing class-vs-instance warning fired anyway, one layer deeper than the builder looked: D-01 closed quoted stamps and the review found the shapes D-01 missed (REQ-255).

## Worth knowing

The hook's `2>/dev/null` is safe only because the repairer deliberately prints failure lines to stdout (D-03) — anyone adding stderr output to the script will silently lose it in the banner. `comparison_key_for`'s space-fold is dead code until REQ-255 resolves it. The 120s skew constant now has a fourth hand-kept copy.

## Back-reference

See `do-work/archive/UR-056/REQ-246-repair-wrong-queue-and-working-timestamps-from-the-session-hook.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `270a2d0`.
