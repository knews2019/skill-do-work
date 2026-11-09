---
title: "Lessons from REQ-179: Make scope-drift.sh parse Scope headers with trailing annotations"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-179-make-scope-drift-sh-parse-scope-headers-.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-179: Make scope-drift.sh parse Scope headers with trailing annotations

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

`skills/do-work/tools/checks/scope-drift.sh` anchors on the literal header `**Files I will touch:**` — both its path extraction (line 47) and its "header present but unparseable" FAIL guard (line 70). A Scope section whose header carries any annotation (REQ-178 wrote `**Files I will touch (all new, all inside …):**`) matches neither, so the check reports `SKIP` exit 2 as if the section were absent — silently disabling the exact comparison the script's own header comment says must never be silently disabled. Root cause: the rule "recognize the header, not one literal spelling of it" is enforced at zero of its two match sites.

## Solution summary

Fixed the silent self-disable: an annotated `**Files I will touch (…):**` header now either parses (real comparison) or FAILs loudly — never SKIPs; genuinely-absent Scope still SKIPs exit 2. Lock-in probes pin the annotated-comparison, zero-paths-FAIL, and absent-SKIP contracts.

## What worked

Reproducing RED against the real defect artifact (REQ-178's archive) instead of only synthetic fixtures — it also surfaced the bare-filename Scope formatting as a bonus signal.

## What didn't work

N/A — first pass green.

## Worth knowing

`contract-regressions.sh` aborts at line ~1797 on case-sensitive filesystems (`Justfile` vs `justfile`, REQ-180) — everything after never runs there, so a green-looking late-suite check may simply be unreached. The process-tree probe is confirmed flaky in this sandbox (failed at RED, passed at GREEN, no related change).

## Back-reference

See `do-work/archive/UR-040/REQ-179-scope-drift-header-parse-tolerance.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `e530dde`.
