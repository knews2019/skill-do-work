---
id: REQ-130
title: Malformed request recovery walks to valid history
status: completed
claimed_at: 2026-08-07T10:58:08Z
completed_at: 2026-08-07T11:00:10Z
route: B
created_at: 2026-08-07T08:45:11Z
user_request: UR-029
addendum_to: REQ-064
domain: general
prime_files: []
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-063, REQ-129, REQ-131, REQ-132]
batch: audited-safety-fixes
effort_estimate: normal
write_set: [tools/checks/blanked-req-scan.sh, _dev/tests/record-commit-hash-guards.sh, actions/version.md, CHANGELOG.md]
---

# Addendum: Malformed Request Recovery Walks to Valid History

## What

Treat every empty or non-empty malformed request blob as damaged history, recover only from the first older blob with parseable frontmatter, and refuse to replace a file unless the recovered temporary content independently passes the same structural validation.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

The current scanner correctly recognizes a malformed filesystem file as damaged, but its historical resolver accepts the first non-empty blob without parsing it. Its recorded-hash resolver also treats that malformed blob as the end of the damage chain. In the reproduced valid-then-malformed history, porcelain mode named the malformed commit as its own recovery source and lost the recorded hash; restore exited 0, claimed full repair, and left a rescan failing.

## Prior Implementation

REQ-063 introduced detection and historical resolution in commit `d91c567`. REQ-064 added consent-gated restore in commit `069c943`, writing through a same-directory temporary file and delegating the `commit:` update to the guarded writer. Preserve those shared detector/restore boundaries and their existing empty-file, partial-repair, and provenance reporting behavior.

## Detailed Requirements

- Replace the early-exiting frontmatter check with a stream-compatible validator that deliberately consumes through end of file, preventing `git cat-file` from receiving SIGPIPE under `set -o pipefail`.
- Use the same structural criterion for filesystem files and historical Git blobs.
- Require every recovery candidate to be non-empty and to have parseable frontmatter.
- While resolving the implementation hash, keep empty and malformed blobs inside the damage chain; inspect their commit subjects and stop only at the first valid historical version.
- For a malformed commit whose subject contains `record commit hash <valid-short-sha>`, report that hash while selecting the earlier valid commit as the recovery source.
- After materializing recovered content into the temporary file, independently verify its frontmatter before moving it over the damaged target.
- A restore must never exit 0 or count a file as fully repaired while that file remains malformed.
- Preserve existing exit-code and porcelain record contracts except for correcting the recovery source and recorded hash values.

## Constraints

- Keep history reads stream-safe under Bash 3.2 and `set -uo pipefail`.
- Preserve atomic same-directory replacement and reuse `record-commit-hash.sh`; do not hand-edit `commit:`.
- Add the malformed-history regression to the existing guard fixture and retain all empty-file and partial-repair cases.
- Run relevant regressions, `bash -n`, ShellCheck, Bash 3.2 checks, and the repository's broader required validation.
- Do not stage, commit, or push without separate user authorization.

## Builder Guidance

Firm on the recovery semantics and false-success prevention. Keep the existing definition of parseable frontmatter unless a test proves it insufficient; the required change is that the validator can consume either a file or a complete blob stream without terminating its producer early.

## Red-Green Proof

**RED prompt/case:** Commit a valid archived request, then commit a non-empty version without parseable frontmatter using a subject containing `record commit hash <valid-short-sha>`. Run porcelain mode, restore mode, and porcelain mode again.
**Why RED now:** The first non-empty historical blob wins even when malformed, and restore validates only that its temporary file is non-empty.
**GREEN when:** The first scan exits 1 and reports the earlier valid commit plus the recorded short hash; restore reinstates the valid body, reapplies the hash, and exits 0; the final scan exits 0 with no damaged request remaining.
**Validation:** User confirmed by requesting capture of the audited findings after reviewing the exact reproduced output.

## Dependencies

No queued dependency. This addendum corrects completed REQ-064 and is related to REQ-063's detector contract.

## Full Context

See `do-work/user-requests/UR-029/input.md` for the complete capture context.

---
*Source: UR-029 - "run do-work capture-request on these issues"*

---

## Triage

**Route: B** - Medium

**Reasoning:** One recovery script owns the behavior, but source selection, damage-chain traversal, replacement validation, and end-to-end restore evidence must change together.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The current validator accepts only a path and exits as soon as it sees a closing delimiter. Recovery chooses the first non-empty blob, hash traversal stops at any non-empty blob, and restore checks only non-emptiness before replacement. A stream-compatible validator that consumes EOF can become the shared validity predicate for the working file, historical blobs, and restore temp files without triggering producer SIGPIPE under pipefail.

## Scope

**Files I will touch:**
- `tools/checks/blanked-req-scan.sh` (modify) - validate history candidates and restore content as parseable frontmatter
- `_dev/tests/record-commit-hash-guards.sh` (modify) - reproduce a valid-to-malformed damage chain and clean repair
- `actions/version.md` (modify) - bump the integration version
- `CHANGELOG.md` (modify) - record the recovery safety fix

**Files I will NOT touch:** The guarded commit-hash writer, cleanup prose, or archived antecedents.

**Acceptance criteria (restated from REQ):**
- [ ] Frontmatter validation accepts a file path or stdin and consumes the complete stream.
- [ ] Recovery selects the newest parseable historical blob, not merely the newest non-empty blob.
- [ ] Malformed and empty blobs remain in the damage chain while resolving the recorded hash.
- [ ] Restore validates its temp file before replacement.
- [ ] The reproduced chain reports the earlier valid SHA and recorded hash, restores valid content, exits zero only after repair, and rescans clean.

## Implementation Summary

**Files changed:**
- `tools/checks/blanked-req-scan.sh` (modified)
- `_dev/tests/record-commit-hash-guards.sh` (modified)
- `actions/version.md` (modified)
- `CHANGELOG.md` (modified)

**What was done:** Reused one EOF-consuming frontmatter validator for working files, historical blobs, and restore temp files; traversed malformed damage chains for source and hash recovery; and added an end-to-end valid-to-malformed repair regression.

## Qualification

Passed - 4 files verified, 5 acceptance requirements traced, P-A-U confirmed.

## Testing

**Tests run:**
- /bin/bash _dev/tests/record-commit-hash-guards.sh - passed
- /bin/bash _dev/tests/contract-regressions.sh - passed

**Red-green validation:**
- RED: The malformed-chain fixture selected the malformed SHA, lost the recorded hash, restored malformed bytes, and remained damaged on rescan.
- GREEN: The focused guard and full contract suites passed after parseability became the shared source, chain, and replacement predicate.

**Existing tests updated:** _dev/tests/record-commit-hash-guards.sh adds the reproduced valid-to-malformed history and verifies exact source SHA, recorded hash, restored content, zero repair exit, and clean rescan.

## Root Cause

Recovery equated non-empty with valid. That let malformed metadata damage terminate both history walks and pass the replacement gate, while an early-exiting validator was unsafe to place behind git cat-file under pipefail.

## Review

**Overall: 100%** | 2026-08-07

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition - this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass - valid recovery selection, damage-chain hash traversal, guarded replacement, and clean rescan are proven end to end.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** One EOF-consuming validator kept the detector, history traversal, and restore gate on the same definition of validity.

**What didn't:** Treating blob size as validity allowed non-empty corruption to masquerade as recovery content and terminate provenance traversal.

**Worth knowing:** A validator used in a pipefail pipeline must consume the complete producer stream; early success can turn a valid large blob into a git cat-file SIGPIPE failure.
