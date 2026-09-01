---
title: "Lessons from REQ-145: Remove the Stateful Pipeline"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-145-remove-the-stateful-pipeline.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-145: Remove the Stateful Pipeline

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Remove the separate resumable pipeline state machine after modular cutover and replace it with a copyable full-cycle prompt.

## Solution summary

Removed the separate stateful pipeline runtime, aliases, state reporting, and installed Stop guard; replaced the public workflow with the exact approved capture → verify → `do-work run` → toolbox presentation prompt. Installer reconciliation now removes only the retired guard while preserving custom Stop hooks in jq and Python paths, and live docs/decisions/tests reflect the prompt-based successor without weakening `do-work run`.

## What worked

- Ratcheting the deleted paths, aliases, state token, exact successor prompt, and installer behavior before implementation made the intended removal observable as RED and protected `do-work run` from accidental weakening.
- Treating the two installer copies and the duplicated crew guidance as lock-step surfaces kept the modular distribution internally consistent.

## What didn't work

- The first installer pass covered executable jq and Python reconciliation but left the human no-JSON-tool fallback saying to preserve every hook. That dead path only surfaced during the independent review, not the GREEN automation.

## Worth knowing

- Retiring a hook requires three aligned paths: fresh configuration, automated migration, and manual fallback. Tests must assert the retired entry disappears while unrelated/custom hooks survive in each path.
- A narrow stale-surface sweep is safer than banning the word “pipeline”; the repository still legitimately uses it for the `do-work run` orchestration, CI, and data flows.

## Back-reference

See `do-work/archive/UR-031/REQ-145-remove-stateful-pipeline.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `c42f228`.
