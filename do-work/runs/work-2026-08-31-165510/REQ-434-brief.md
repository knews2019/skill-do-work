# REQ-434 Builder Brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-434-refuse-unsupported-timestamp-ordering-anchors`
Branch: `worktree-agent-REQ-434-refuse-unsupported-timestamp-ordering-anchors`
REQ: `do-work/working/REQ-434-refuse-unsupported-timestamp-ordering-anchors.md`
Integration owner: main checkout; do not edit `do-work/`, release metadata, changelogs, or version files.

## Goal

Prevent offset and fractional timestamp shapes that are parseable but unsupported for canonical repair from acting as ordering anchors for supported successor fields. Preserve supported ordering and the current-time ceiling.

## Required method

- Read the REQ, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, its lessons satellite, and the applicable backend/testing/general/coding guardrails.
- Locate the canonical doctor timestamp planning seam and its existing mixed-shape tests.
- Write the exact mixed unsupported/supported regression first and capture a real failing result before production changes.
- Make the smallest coherent production correction, keeping unsupported field bytes unchanged and explicit refusal evidence intact.
- Run the focused regression, affected package tests, `go vet ./...`, full module tests, exact Go 1.25 compatibility, and `git diff --check`.
- Review every changed file, commit the implementation on this branch, and write the durable handback to `do-work/runs/work-2026-08-31-165510/REQ-434-handback.md` in the main checkout.

## Handback

Record branch and full commit hash, approach, RED/GREEN output, changed files, all commands with exit results, decisions, seams, and any discovered tasks. Do not create or edit queue files.
