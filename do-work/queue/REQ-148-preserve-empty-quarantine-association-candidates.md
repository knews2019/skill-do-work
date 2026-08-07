---
id: REQ-148
title: "Addendum: preserve association candidates with empty quarantine"
status: pending
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
