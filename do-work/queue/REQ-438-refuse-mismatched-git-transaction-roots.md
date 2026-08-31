---
id: REQ-438
title: '[impact-critical] Refuse mismatched Git transaction roots'
status: pending
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-mechanical
related: [REQ-437, REQ-439, REQ-440, REQ-441, REQ-442, REQ-443, REQ-444]
batch: accepted-feedback-regressions
---

# Refuse Mismatched Git Transaction Roots

## What

Require a mutating command's supplied repository root to identify the same physical worktree root returned by Git. Refuse nested roots before discovery callbacks can mutate files that preflight, rollback, change detection, and `--commit` inspect at a different path.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Verbatim claim / severity:** `[P1] Require the supplied root to match the Git worktree root.`
- **Evidence:** `git_transaction.go:333-349` silently replaces the supplied root with `git rev-parse --show-toplevel`, while cleanup discovery and mutation retain the supplied root.
- **Origin / earned by:** Introduced by transaction foundation `329c55a9`. An isolated nested-root replay returned cleanup success, moved files below `nested/do-work/`, left HEAD unchanged, and left uncommitted changes.
- **Surface-cost:** Earned. The false-success `--commit` incident justifies one transaction-boundary identity check and error; that is cheaper than rebasing every generic callback and target path. Test exact and nested roots, and do not falsely reject a physical alias of the same root.

## Detailed Requirements

- Reject a supplied root that is physically below or otherwise distinct from Git's worktree root.
- Apply the same identity contract to preflight and transaction execution before any mutation callback runs.
- Leave target bytes, index, and HEAD unchanged on refusal and return actionable evidence naming the mismatch.
- Preserve exact-root behavior and legitimate path aliases that resolve to the same physical directory.

## Constraints

- Prefer one shared root-identity check. Do not add per-command target rebasing.
- Preserve the existing installer precedent that mutations require the actual Git project root.

## Red-Green Proof

**RED prompt/case:** In a Git repository with `nested/do-work/`, run cleanup with `--repo-root <repo>/nested --commit` and a mutatable nested REQ.
**Why RED now:** Git guards inspect `<repo>/do-work/...` while cleanup mutates `<repo>/nested/do-work/...`, allowing false success without a commit.
**GREEN when:** The nested-root invocation refuses before its callback, leaves bytes/index/HEAD unchanged, exact-root transactions remain green, and a physical alias of the exact root is handled deliberately.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm on rejection and physical identity; implementation details belong to the builder.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 12 from the validated external feedback.*
