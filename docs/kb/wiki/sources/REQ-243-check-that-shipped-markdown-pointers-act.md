---
title: "Lessons from REQ-243: Check that shipped markdown pointers actually resolve"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-04/REQ-243-check-that-shipped-markdown-pointers-act.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-243: Check that shipped markdown pointers actually resolve

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

`_dev/tests/prescribed-shell-canonicalization.sh` proves a restatement is **absent** from shipped markdown. Nothing proves that the pointer which replaced it **resolves** — not the relative path, not the `#anchor`. Add that check.

## Solution summary

**Files changed:**
- `_dev/tests/shipped-package-reference-contract.sh` (modified)

## What worked

- **The gap a REQ names and the gap a REQ has are not always the same gap.** REQ-243's Red-Green
- Proof asserted "no check reads markdown links at all, so the wrong depth passes the entire suite
- today." Running the stated mutation against the unmodified tree took one command and disproved it:
- `shipped-package-reference-contract.sh` has resolved relative markdown targets from the citing
- file's own directory, in two topologies, for some time. Run the RED against pre-change code
- *before* building — not to satisfy a ritual, but because a "RED" that was already red is the
- cheapest possible signal that the work is half-done already.
- **Two hand-fixes of the same class is good evidence a checker is needed; it is not evidence that
- no checker exists.** REQ-230 and REQ-238 both verified the path by hand and the anchor by hand.
- The path half was machine-checked the whole time and nobody knew. Before writing the third
- instance of a manual check, grep for the check as well as for the defect.
- **"0 checks ran" from a relocated script is usually the harness.** The instrumented probe reported
- zero anchors checked when copied to `/tmp`, which reads exactly like a vacuous pass. The cause was
- `repo_root` derived from `BASH_SOURCE`. Any suite that locates itself relative to its own file
- cannot be probed from outside the tree — copy it *into* the tree, run, delete.
- **Where masking helps and where it hurts, in one file.** `strip_markdown_code` is what makes link
- extraction trustworthy and what makes heading extraction wrong, for the same reason. The
- offset-preserving property — already fixture-locked — is what lets both be right at once.

## Back-reference

See `do-work/archive/UR-042/REQ-243-check-that-shipped-markdown-pointers-resolve.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `37d7729`.
