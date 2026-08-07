---
id: REQ-148
title: "Addendum: preserve association candidates with empty quarantine"
status: completed
completed_at: 2026-08-07T22:55:24Z
commit:
claimed_at: 2026-08-07T22:51:17Z
status_changed_at: 2026-08-07T22:51:17Z
route: A
kb_status: pending
kb_entry:
created_at: 2026-08-07T22:40:46Z
user_request: UR-034
addendum_to: REQ-128
domain: general
prime_files: []
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-121, REQ-128]
maintenance: false
effort_estimate: trivial
write_set: [actions/commit.md, actions/inspect.md, skills/do-work/actions/commit.md, skills/do-work-toolbox/actions/inspect.md, _dev/tests/contract-regressions.sh]
---

# Addendum: Preserve Association Candidates with Empty Quarantine

## What

Correct REQ-128's commit and unscoped-inspect quarantine merge so an empty run-level secret quarantine preserves every safe inventory candidate for REQ association.

## Context

The merge currently reads the quarantine file and the fresh `uncommitted-inventory.sh` output with an `NR == FNR` first-file rule. When the quarantine is empty, `NR` and `FNR` start together and remain equal throughout the inventory file, so every inventory row enters the exclusion-table branch. This is broader than dropping only the first candidate and explains the live reproduction where six legitimate candidates became an empty stdin stream for `associate-files.sh`.

This addendum corrects REQ-128 without weakening its once-X-always-X quarantine. The active bridge actions and staged modular actions both contain the faulty merge.

## Prior Implementation

REQ-121 introduced `tools/checks/uncommitted-inventory.sh` and `tools/checks/associate-files.sh` in commit `167b0ae`, then rewired `actions/commit.md` and `actions/inspect.md` around those shared tools. REQ-128, itself an addendum to REQ-121, added the deterministic run-level quarantine and the current two-file merge in commit `7bb03d2`. The association helper already uses the portable `FILENAME == ARGV[1]` discriminator for its own two-file awk join.

## Requirements

- Preserve every `M`, `A`, `D`, and `XD` candidate when the quarantine file is empty.
- Continue excluding current `X` rows and every exact path retained in a populated run-level quarantine.
- Replace the unsafe `NR == FNR` discriminator with the repository's existing portable `FILENAME == ARGV[1]` form.
- Update every checked-in bridge or modular copy of the commit and inspect merge that still exists when this request runs.
- Keep `tools/checks/associate-files.sh`'s interface unchanged and do not add a new micro-helper for this one-condition correction.
- Add regression coverage for empty and populated quarantine files using multiple safe candidates and an `X` row.

## Constraints

- Preserve REQ-128's once-X-always-X safety behavior and never expose secret-shaped contents.
- Keep the awk logic portable to the repository's supported shell and awk floor.
- Do not fold unrelated modular-migration work into this fix.

## Red-Green Proof

**RED prompt/case:** Run the candidate filter with an empty quarantine and an inventory containing multiple safe `M`/`A`/`D` rows. The current `NR == FNR` rule emits no candidates, so `associate-files.sh` receives empty stdin and exits 1.

**Why RED now:** An empty first file contributes no records, leaving `NR` and `FNR` equal for every record in the second file; the first-file branch consumes the complete inventory.

**GREEN when:** An empty quarantine emits every safe candidate; a populated quarantine excludes only retained paths; current `X` rows remain excluded; and both bridge and modular action copies are covered by regression checks.

**Validation:** User confirmed after the live six-candidate reproduction and the independent validate-feedback triage.

## Assets

None.

---
*Source: `do-work capture-request for accepted issue`, referring to the accepted empty-quarantine validate-feedback finding in the preceding conversation.*

## Triage

**Route: A** — this is a reproduced one-condition shell bug with exact affected files, an established repository-local correction, and explicit regression cases.

## Pre-Flight

The archive collision check passed. All four active bridge/modular action copies still use `NR == FNR`; `tools/checks/associate-files.sh` already demonstrates the portable `FILENAME == ARGV[1]` form. Unrelated dirty modular-cutover planning files will be preserved and excluded.

## Root Cause

`NR == FNR` identifies the first input only while that input contributes at least one record. An empty quarantine contributes none, so `NR` and `FNR` remain equal throughout the second file and every inventory row is consumed by the exclusion-table branch. The merge needs to identify its first input by filename, not record counters.

## Implementation Summary

- `actions/commit.md` (modified) — uses `FILENAME == ARGV[1]` for the commit candidate merge.
- `actions/inspect.md` (modified) — uses the same safe discriminator for unscoped inspection.
- `skills/do-work/actions/commit.md` (modified) — mirrors the bridge fix in staged core.
- `skills/do-work-toolbox/actions/inspect.md` (modified) — mirrors the inspect fix in staged toolbox.
- `_dev/tests/contract-regressions.sh` (modified) — ratchets all four sources and executes empty/populated quarantine cases with multiple safe rows and a current X row.
- `VERSION` (modified) — bumps the bridge release to v0.183.23.
- `actions/version.md` (modified) — reports v0.183.23.
- `CHANGELOG.md` (modified) — records the empty-quarantine association fix.
- `skills/do-work/VERSION` (modified) — synchronizes the staged core release.
- `skills/do-work/actions/version.md` (modified) — synchronizes staged version reporting.
- `skills/do-work/CHANGELOG.md` (modified) — mirrors the release notes.

## Testing

**RED:** The new source ratchets failed for all four action copies, each still carrying the unsafe counter discriminator. The behavioral fixture independently demonstrated the intended empty and populated quarantine outputs.

**GREEN:** Replacing exactly those four conditions makes the full contract suite pass. The empty fixture emits every M/A/D/XD path, while the populated fixture excludes a prior quarantined path and the current X path without dropping unrelated candidates.

Warning-level ShellCheck, Bash syntax, Just parsing, root and staged queue-kanban Go tests/vet, all repository contract suites, queue-kanban verification, and `git diff --check` pass.

## Qualification

Passed. The four production changes are the exact one-condition correction requested; the regression covers the bridge and modular sources plus both data states, and no helper or interface changed.

## Review

**Acceptance: Pass.** `FILENAME == ARGV[1]` is already the repository's working portable two-file join pattern and remains correct whether its first file has zero, one, or many records. The once-X-always-X overlay and current-X exclusion are unchanged. No Important or Minor findings remain.

## Lessons Learned

- `NR == FNR` is not a safe first-file test when an input may be empty; identify that input explicitly through `FILENAME` and `ARGV`.
- A regression for a two-input merge must cover the zero-record first input, not only populated joins.

## Orientation

[MAP UNCHANGED] Commit and inspect still own their candidate filtering inline; only the first-input discriminator changed, with no new tool or interface.
