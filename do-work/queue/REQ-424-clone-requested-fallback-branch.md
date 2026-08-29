---
id: REQ-424
title: 'Clone the branch named by the tarball URL'
status: pending
created_at: 2026-08-29T20:26:10Z
user_request: UR-082
domain: backend
prime_files: ['_dev/primes/prime-shell-commands.md']
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-421, REQ-422, REQ-423]
batch: accepted-review-fixes
write_set: ['tools/fetch-upstream-archive.sh', 'skills/do-work/tools/fetch-upstream-archive.sh', '_dev/tests/update-script-behavior.sh', '_dev/primes/prime-shell-commands.md', '_dev/primes/lessons-shell-commands.md']
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
---

# Clone the Branch Named by the Tarball URL

## What
Make Git fallback select the branch parsed from the existing canonical tarball URL grammar rather than silently archiving the remote's default HEAD.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Preserve the branch parsed by the existing canonical `/archive/refs/heads/<branch>.tar.gz` grammar.
- Pass that branch to a shallow single-branch clone during Git fallback.
- Fail when the requested ref is missing rather than substituting default HEAD.
- Retain default-HEAD fallback behavior for URLs from which the current grammar cannot derive a branch.
- Test forced HTTP failure with distinct default-branch and requested-branch markers.

## Constraints
- Do not add URL grammars or support query strings/fragments in this change.
- Keep root and shipped fetcher scripts byte-identical.
- Add concise ref-preservation guidance to the shell prime and a detailed linked lesson entry.

## Dependencies
None. It shares fetcher and shell-test files with REQ-423 and is implemented as one shell slice in this batch.

## Builder Guidance
Certainty: Firm. Select the parsed branch only when the existing grammar yields one; otherwise leave the existing default clone route unchanged.

## Context
No pending or unassigned queue candidate shares this root cause. Provenance: accepted review finding `[P2] Clone the branch named by the tarball URL` against `skills/do-work/tools/fetch-upstream-archive.sh:90-91`. The review states that `upstream_branch` is extracted but the fallback clones default HEAD and packages it beneath the requested branch prefix.

## Red-Green Proof
**RED prompt/case:** Force HTTP failure for a canonical non-default branch URL against a repository whose default and requested branches contain different markers.
**Why RED now:** The fallback clone archives default HEAD regardless of the parsed branch.
**GREEN when:** The archive contains the requested marker, excludes the default marker, and a missing requested ref fails; an unparseable URL still clones default HEAD.
**Validation:** User accepted the review finding and supplied the implementation plan.

## Full Context
See `do-work/user-requests/UR-082/input.md` for the approved plan and batch constraints.

---
*Source: accepted review finding [P2] on fallback branch selection, followed by the user-approved plan.*
