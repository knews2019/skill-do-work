---
id: REQ-565
title: '[impact-critical] Review fix: Close residual workspace release identity gaps'
status: pending
created_at: 2026-09-03T23:26:06Z
user_request: UR-110
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-512]
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
review_generated: true
addendum_to: REQ-512
sweep: true
sweep_key: legacy-finalization-workspace-identity-residual
---

# Close Residual Workspace Release Identity Gaps

## What

Finish the fail-closed workspace ownership boundary left by REQ-512's post-remediation re-review: prove Cargo and uv source identity is unique, and require every structurally present npm root version mirror when the root source changes.

The fold-first scan found no pending or pending-answers REQ in any UR that shares this exact residual finalization workspace identity root cause.

## Context

REQ-512 completed its one permitted remediation pass and closed changed-source-first selection, exact shared-lock replacement, malformed missing/unquoted identity, malformed npm JSON, and bounded fold termination. Independent re-review found the two remaining cases below, so this successor owns them instead of silently taking a second remediation pass.

## Instances

- [ ] `finalization_discovery.go:tomlSectionScalar` accepts the first Cargo/uv `name` without proving uniqueness or TOML identity validity.
- [ ] `npmRootVersionCopies` counts only root lock values already equal to `oldVersion`, so parseable stale or mismatched root copies can disappear from the required mirror set.

## Requirements

- Refuse a changed Cargo or uv source when its applicable package/project identity is duplicated, competing, ambiguous, or otherwise not uniquely parseable.
- When a changed npm root has a tracked parseable lock, treat every structurally present root `version` and `packages[""].version` copy as an obligation; a stale or mismatched copy must refuse rather than be omitted.
- Preserve member-only releases, pre-existing target-version neighbors, multiple changed members sharing one lock, typed enumeration/ownership failures, protected-path refusal, and public recovery-to-claim success.
- Add committed RED/GREEN tests for duplicate/competing Cargo and uv names plus one-stale and both-stale npm root copies.

## Red-Green Proof

**RED prompt/case:** Run strict finalization discovery for (a) a changed Cargo/uv source with duplicate or competing applicable `name` declarations and (b) a changed npm root whose parseable lock has one or both structurally present root copies stale.

**Why RED now:** The TOML helper returns the first name, while npm root-copy counting ignores present values that do not already equal the expected old version; both can silently reduce required mirror ownership.

**GREEN when:** Every ambiguous source identity and stale structurally present npm root copy refuses before mutation with typed path evidence, exact valid mirrors still finalize, and the full REQ-512/recovery matrix remains green.

**Validation:** REQ-512 post-remediation re-review; successor required by the one-remediation rule.

## Full Context

See `do-work/user-requests/UR-110/input.md` and REQ-512's `## Re-Review` section.

## Open Questions

- [x] Auto-approved: critical severity (release/finalization ownership risk). → Added to queue immediately.
