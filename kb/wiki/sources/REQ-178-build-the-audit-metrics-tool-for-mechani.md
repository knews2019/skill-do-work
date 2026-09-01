---
title: "Lessons from REQ-178: Build the audit-metrics tool for mechanical audit measurement"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-178-build-the-audit-metrics-tool-for-mechani.md]
related:
  - page: REQ-176-implement-the-maintainability-audit-acti
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-178: Build the audit-metrics tool for mechanical audit measurement

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

A small Go tool, `skills/do-work-toolbox/tools/audit-metrics/`, that produces the maintainability audit's deterministic numbers mechanically — inventory, distributions, band flags, churn — so the action pastes tool output instead of prescribing fragile find/wc/awk pipelines to an LLM. Script what can be scripted; judgment stays in prose.

## Solution summary

Built the vendored audit-metrics Go CLI (zero dependencies, queue-kanban conventions): four subcommands — `inventory`, `folders`, `churn`, `hotspots` — emitting pasteable markdown tables with flag-supplied WATCH/FLAG bands, exclude-prefix filtering, shallow-clone reporting, and rename+copy-normalized churn (`-M -C --find-copies-harder` with dead copy-source reassignment; verified to reproduce `git log --follow`'s 214-touch count for work.md across the 2026-08-08 restructure). 10 lock-in tests including real-git and real-shallow-clone fixtures.

## What worked

Mirroring queue-kanban's conventions wholesale (renderer/io.Writer split, per-subcommand FlagSets, real-git fixtures in t.TempDir()) meant zero design churn; the real-repo spot-check during build caught the biggest correctness bug (staged-copy migration) before review.

## What didn't work

`-M` rename detection alone missed the 2026-08-08 skills/ restructure entirely (8 vs 214 touches) — it was a staged copy-then-delete, invisible to rename detection; only `-C --find-copies-harder` plus dead-copy-source reassignment reproduces `git log --follow`. Also: phrasing the Scope header with a parenthetical silently disabled scope-drift.sh (→ REQ-179).

## Worth knowing

Churn numbers from this tool are only trustworthy because of copy detection — anyone replacing it with a plain `git log --name-only | sort | uniq -c` resurrects the dead-path split. Shallow clones are detected and reported, never silently truncated. The tool is a separate Go module — a repo-root `go build ./...` never reaches it.

## Back-reference

See `do-work/archive/UR-040/REQ-178-audit-metrics-mechanical-tool.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `1afe780`.
