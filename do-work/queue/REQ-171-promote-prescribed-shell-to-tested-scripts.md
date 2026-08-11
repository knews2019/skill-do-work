---
id: REQ-171
title: "Addendum: promote prescribed shell primitives to shipped, fixture-tested scripts"
status: pending
created_at: 2026-08-11T13:58:25Z
user_request: UR-038
addendum_to: REQ-165
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-167, REQ-170]
batch: stabilization-audit
---

# Addendum: Promote Prescribed Shell Primitives to Shipped, Fixture-Tested Scripts

## What

Graduate the canonical prescribed-shell primitives from documented prose to real, shipped script files with fixture-repo execution tests. Each *multi-line* primitive in `skills/do-work/docs/prescribed-shell-primitives.md` (and the remaining multi-line blocks in action files, e.g. capture.md Step 4's screenshot copy/verify/link block) becomes a `.sh` file under a per-package `scripts/` directory (core first: `skills/do-work/scripts/`); call sites keep a one-line intent statement plus the invocation; `_dev/tests/` gains a fixture-repo scaffold (mktemp repo, git init, seeded queue/version fixtures) that *executes* each script and asserts output and exit codes. The dividing line is "does this block contain logic that can be wrong" — one-liners and illustrative fragments stay inline, covered by the existing lint harness as residue.

## Context

Addendum to REQ-165 (completed). The plan this implements was approved by the user *while the original batch was being built* — the queue ran between plan approval and this capture, so the delta arrives as an addendum rather than a reshape. Lint (what shipped) catches syntax and quoting; it structurally cannot catch the semantic trap classes that motivated the batch — pipefail-kills-the-fallback, porcelain collapsing untracked dirs, curl partial files, merge-commit `git show` — which only surface when the commands *execute*. Only execution tests close that class; REQ-166's `session-start-hook-behavior.sh` is the proof of shape and the pattern to generalize.

## Prior Implementation

- **REQ-165** (commit `a45d5c4`, Route C): built `_dev/tests/action-shell-blocks.sh` (214 lines) — extracts fenced `bash`/`sh` blocks from the shipped `skills/` tree, runs `bash -n` always and shellcheck when present; wired into contract-regressions. This harness **stays** — its role narrows to the inline residue once multi-line blocks are promoted.
- **REQ-167** (commit `1a27c07`): created the canonical prose home `skills/do-work/docs/prescribed-shell-primitives.md` with eight primitive sections (per-file untracked inventory, merge-aware commit diff, commit file listing, local git ignore, atomic download publication, raw text before shell quoting, diff output filtering, state across command blocks), pointed consuming action files at it, and ratcheted the arrangement with `_dev/tests/prescribed-shell-canonicalization.sh` plus an audit record in `decisions/audits/2026-08-11-prescribed-shell-primitives.md`.
- **REQ-166** (commit `6538bdd`, Route A): simplified the session-start hook and added `_dev/tests/session-start-hook-behavior.sh` — the fixture-execution pattern this addendum generalizes.

## Detailed Requirements

- Build the shared fixture scaffold first (helpers to create a throwaway repo with seeded queue/version/untracked-file state), following `_dev/tests/` conventions; REQ-166's hook test may be migrated onto it if that is a simplification, not a rewrite.
- Promote each qualifying primitive to a script; the canonical guide's sections then document intent and *point at the script* as the normative implementation — update `prescribed-shell-canonicalization.sh`'s contract accordingly rather than deleting it (the single-home ratchet must survive the move, now enforcing script-as-home).
- Every promoted script gets fixture tests covering the trap it exists to avoid (e.g. the download script's mid-transfer-failure case; the untracked-inventory script against a wholly-untracked directory).
- Go-owned capabilities (atomic REQ reservation) get no shell twin — script layer is shell-portable primitives only.
- Call sites keep one sentence of intent so action files still work as standalone pasted prompts; the floor (read/write files, run shell) is respected because scripts ship inside the packages.
- Cross-package: a primitive used by board/knowledge/toolbox actions lives in core `scripts/` and is referenced with explicit sibling paths, same direction rules as prose cross-references.
- Net-surface accounting in the report: lines of prescribed shell removed from prose vs. script+test lines added — the prose side must shrink.

## Builder Guidance

Certainty: Firm on the architecture (approved plan); exploratory on which primitives qualify — produce the promotion inventory (primitive → script, or stays-inline rationale) as the first artifact, then execute it. Migrate incrementally with the suite green after each promotion; this is one REQ because the scaffold is shared, but the promotions are independent and abortable partway with value retained.

## Red-Green Proof

**RED prompt/case:** Reintroduce a semantic trap into a canonical primitive — e.g. change the download primitive back to plain `curl -o` without temp-and-rename, or give a version-parse pipeline the `set -euo pipefail` dead-fallback shape. Today `action-shell-blocks.sh` (bash -n + shellcheck) passes both, and nothing executes the primitive to notice.
**Why RED now:** The shipped harness is lint-only; the canonical guide is prose-only. The exact bug class that motivated UR-036 (demonstrated live in the session-start hook) remains undetectable for every primitive except the hook itself.
**GREEN when:** The promoted scripts exist with fixture tests; seeding either regression above into its script makes the suite fail naming the script and case; the canonicalization ratchet enforces script-as-home; `action-shell-blocks.sh` still covers inline residue; net prescribed-shell prose shrank.
**Validation:** User confirmed (plan approved with "capture"; delta re-anchored against the shipped batch by the capturing agent).

## Full Context

See `do-work/user-requests/UR-038/input.md` for complete verbatim input.

---
*Source: "capture" — approving the stabilization plan v2 discussed in-session (UR-038)*
