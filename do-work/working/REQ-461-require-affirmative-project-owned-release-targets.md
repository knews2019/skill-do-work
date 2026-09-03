---
id: REQ-461
title: '[impact-user-visible] Require affirmative project-owned release targets'
status: claimed
created_at: 2026-09-01T00:12:38Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
kb_status: pending
estimate:
  p50_active_minutes: 50
  confidence: medium
  calculated_at: 2026-09-03T09:45:12Z
  basis:
    - Route C
    - 6-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - persistence changes
    - cross-route regression gates
    - full-suite verification
related: [REQ-413]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-413
status_changed_at: 2026-09-03T11:10:08Z
claimed_at: 2026-09-03T11:32:55Z
---

# Require Affirmative Project-Owned Release Targets

## What

Replace convention-based installed/generated path exclusions with condition-complete evidence that every release target is a project-owned source or declared maintainer mirror.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read both prime files, both lesson satellites, and required crew rules. Implement the planned exact normalized partition in publication, close legacy discovery through tracked repository evidence plus declared maintainer topology, then update direct fixtures and shipped contract prose; verify RED before production edits, then focused/full Go and contract gates.
- [x] **[APPLY]:** Added RED fixtures proving nested npm/Cargo/uv workspaces and unrooted chains self-authorized, then restricted initial manifest ownership to repository-root and declared maintainer topology and recursively promoted only members of already-owned workspace manifests. Workspace-lock association now uses that same proven owner.
- [x] **[UNIFY]:** Reviewed the remediation diff and confirmed root-only initial manifest ownership, recursive promotion through already-owned workspace parents, and proven-owner lock association. Focused finalization ownership/workspace tests, full CLI tests, vet, contract regressions, and `git diff --check` pass with no debug artifacts.

## Finding Provenance

REQ-413's fresh re-review found that release exclusion recognizes named conventions such as `vendor`, `node_modules`, `.codex/skills`, and literal generated directories, but accepts other dependency or generated locations such as `third_party/do-work`, `dist/skills/do-work`, and cache-owned skill trees.

## Detailed Requirements

- Require affirmative, verifiable project-owned classification for every consumer release target instead of inferring safety from the absence of known bad directory names.
- Permit declared maintainer mirrors only through the existing explicit mirror contract and retain byte-identity validation for changelog mirrors.
- Refuse installed skills, dependencies, vendored packages, caches, distribution outputs, and generated trees regardless of their directory spelling.
- Add fixtures for non-example paths including `third_party/do-work`, `dist/skills/do-work`, and an arbitrarily named cache tree.
- Keep refusal actionable by identifying the target and the ownership evidence that is missing or inconsistent.

## Constraints

- Do not replace one directory-name denylist with a larger denylist.
- Preserve caller-selected semantic version and changelog-content judgment.

## Dependencies

No request prerequisite.

## Builder Guidance

Certainty level: Firm on the invariant, flexible on the proof mechanism. Prefer an explicit manifest or repository-owned declaration that the planner can verify mechanically.

## Red-Green Proof

**RED prompt/case:** Supply otherwise valid release targets beneath `third_party`, `dist`, and an unrecognized cache/generated subtree.
**Why RED now:** The current exclusion is a finite convention list, so targets outside those spellings are accepted.
**GREEN when:** Undeclared or non-project-owned targets refuse independent of path spelling, while verified project sources and declared maintainer mirrors apply atomically.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-413-rereview.md` for source context and independent evidence.

---
*Source: REQ-413 fresh re-review finding F7.*

---

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3636 tokens; matches the shipped release action contract and its downstream readers, but the partial-slug satellite cannot be narrowed within the 2000-token budget.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2786 tokens; matches publication topology classification and structured refusal evidence, but the partial-slug satellite cannot be narrowed within the 2000-token budget.

## Re-Review

**Overall: 98%** | 2026-09-03T10:56:08Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 97% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings:** None — the prior impact-critical self-authorization finding is closed and was not appended to REQ-512.

**Minor findings:** 0 (report only)
**Acceptance:** Pass — nested npm, Cargo, and uv roots and unrooted chains refuse, while a repository-rooted recursive workspace chain and its lock association pass.
**Suggested testing:** 0 items beyond the canonical maintainer gate immediately before integration.
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action after remediation*

