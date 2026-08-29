---
id: REQ-424
title: 'Clone the branch named by the tarball URL'
status: completed
claimed_at: 2026-08-29T20:31:49Z
completed_at: 2026-08-29T20:46:08Z
route: B
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
- [x] **[PLAN]:** Carry the existing parsed branch into an optional shallow single-branch clone argument set and add requested/default/missing-ref fallback coverage.
- [x] **[APPLY]:** Passed the branch parsed by the existing tarball grammar to a shallow single-branch clone, kept default-HEAD behavior only when no branch is derivable, and added divergent-ref regressions.
- [x] **[UNIFY]:** Reviewed the complete diff, verified shell syntax/lint and mirror equality, ran the focused update behavior suite and the canonical maintainer verifier; all required lanes passed.

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

---

## Triage

**Route: B** - Medium

**Reasoning:** The branch selector is localized, but the regression must distinguish requested, default, and missing-ref behavior through the real fallback seam.

**Planning:** Not required; the user supplied an implementation plan.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `tools/fetch-upstream-archive.sh` (modify) — preserve the parsed branch during fallback
- `skills/do-work/tools/fetch-upstream-archive.sh` (modify) — byte-identical shipped mirror
- `_dev/tests/update-script-behavior.sh` (modify) — requested/default/missing-ref regressions
- `_dev/primes/prime-shell-commands.md` (modify) — concise signal/ref guidance shared with REQ-423
- `_dev/primes/lessons-shell-commands.md` (modify) — linked detailed lessons shared with REQ-423

**Files I will NOT touch:** URL parsing grammar, query/fragment handling, or caller interfaces.

**Acceptance criteria (restated from REQ):**
- [x] Canonical branch URLs archive that branch.
- [x] Missing requested refs fail without default fallback.
- [x] Unparseable URLs retain default-HEAD behavior.
- [x] Fetcher mirrors remain byte-identical.

## Implementation Summary

**Files changed:**
- `tools/fetch-upstream-archive.sh` (modified)
- `skills/do-work/tools/fetch-upstream-archive.sh` (modified)
- `_dev/tests/update-script-behavior.sh` (modified)
- `_dev/primes/prime-shell-commands.md` (modified)

**What was done:** The existing parsed branch now selects a shallow single-branch clone before `git archive`; a missing ref fails that route. URLs outside the existing branch grammar still clone default HEAD.

## Testing

**Tests run:** `bash _dev/tests/update-script-behavior.sh`; Bash syntax and ShellCheck; fetcher `cmp`; `bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/maintainer-verify.sh`.

**Result:** All passed; both mirrors compare byte-identical and the canonical verifier exited 0.

**Red-green validation:**
- Divergent default/requested fixture: RED archived the default marker and omitted the requested marker; missing ref returned 0 and published → GREEN selects the requested marker, rejects default substitution, and fails a missing ref without publishing.

**New tests added:**
- Requested non-default branch, missing requested branch, and unparseable-URL default-HEAD cases.

## Lessons Learned

Changing transport must preserve the semantic selector carried by the original source. Naming an archive prefix after a requested branch does not select that branch; a useful regression needs divergent default/requested contents plus a missing-ref case so both accidental defaulting and unconditional branch selection fail visibly.
