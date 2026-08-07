---
id: REQ-131
title: Commit-hash idempotency rejects body changes
status: completed
claimed_at: 2026-08-07T11:01:20Z
completed_at: 2026-08-07T11:03:42Z
commit: 2eb5252
route: B
created_at: 2026-08-07T08:45:11Z
user_request: UR-029
addendum_to: REQ-062
domain: general
prime_files: []
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-129, REQ-130, REQ-132]
batch: audited-safety-fixes
effort_estimate: normal
write_set: [tools/checks/record-commit-hash.sh, _dev/tests/record-commit-hash-guards.sh, actions/version.md, CHANGELOG.md]
---

# Addendum: Commit-Hash Idempotency Rejects Body Changes

## What

Keep the interrupted-metadata-edit success path, but allow it only when removing the top-level frontmatter `commit:` field from HEAD and the working file makes the remaining content identical. Reject every pending difference outside that field.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

The current idempotent branch checks only whether the requested value is already present and whether the path differs from HEAD. The deep audit reproduced two runs: metadata-only pending content exited 0 correctly, while the same metadata edit plus an unrelated body line also exited 0 and printed `Stage and commit it now.`

## Prior Implementation

REQ-062 introduced the guarded writer in commit `62a4188`. It scopes edits to the first frontmatter block, makes body `commit:` examples structurally unreachable, validates exact edit arithmetic, and preserves an interrupted metadata edit as an idempotent success. This addendum narrows only that success condition; retain all other input, size-floor, write, and `--verify` guards.

## Detailed Requirements

- When the working frontmatter already contains the requested commit value and the tracked file differs from HEAD, compare normalized HEAD and working content.
- Normalize by removing only the column-zero `commit:` field inside the first frontmatter block from each side.
- Preserve every other byte for comparison, including body prose, fenced examples, later delimiter blocks, indented nested keys, and body lines beginning with `commit:`.
- If normalized content is identical, retain exit 0 and the interrupted-metadata-edit instruction.
- If normalized content differs, exit 1, explain that pending content extends beyond the frontmatter field, and do not print `Stage and commit it now.`
- Do not modify the file, index, or repository on the rejected path.
- Keep behavior for untracked files, non-Git use, committed no-ops, and actual write operations unchanged.

## Constraints

- Use a comparison shape that cannot hide producer failures or SIGPIPE under `set -o pipefail`.
- Preserve Bash 3.2 compatibility and existing exit-code vocabulary.
- Add both the allowed metadata-only case and rejected metadata-plus-body case to the existing guard fixture. Make a body `commit:` change comparison-significant.
- Run relevant regressions, `bash -n`, ShellCheck, Bash 3.2 checks, and the repository's broader required validation.
- Do not stage, commit, or push without separate user authorization.

## Builder Guidance

Firm on the normalization boundary. Reuse the writer's existing definition of the first frontmatter block so the idempotency guard and actual editor cannot drift about which `commit:` field is metadata.

## Red-Green Proof

**RED prompt/case:** Commit a request with an old or absent frontmatter hash, then leave the requested hash on disk together with a body edit, including a changed body or fenced-example `commit:` line. Invoke the writer with the already-present requested hash.
**Why RED now:** Any tracked difference enters the stranded-edit success branch without checking what changed.
**GREEN when:** A metadata-only pending edit exits 0 with the staging instruction; the metadata-plus-body case exits 1, preserves the file exactly, explains the extra content diff, and contains no staging instruction.
**Validation:** User confirmed by requesting capture of the audited findings after reviewing both reproduced exit-zero runs.

## Dependencies

No queued dependency. This addendum corrects completed REQ-062.

## Full Context

See `do-work/user-requests/UR-029/input.md` for the complete capture context.

---
*Source: UR-029 - "run do-work capture-request on these issues"*

---

## Triage

**Route: B** - Medium

**Reasoning:** The behavior is localized to one idempotency branch, but safe normalization must preserve every body byte and needs multiple real-git regressions.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The existing branch treats every tracked worktree difference as a stranded metadata edit. The fixture already proves metadata-only recovery, and its fenced YAML example is the exact trap for file-wide commit-line removal. Comparing temporary normalized HEAD and worktree blobs after removing only a column-zero commit field inside the first frontmatter block preserves all body and fenced-example text.

## Scope

**Files I will touch:**
- `tools/checks/record-commit-hash.sh` (modify) - guard idempotency with normalized content comparison
- `_dev/tests/record-commit-hash-guards.sh` (modify) - retain metadata-only recovery and reject body differences
- `actions/version.md` (modify) - bump the integration version
- `CHANGELOG.md` (modify) - record the idempotency safety fix

**Files I will NOT touch:** Recovery scanning, commit procedure prose, or archived antecedents.

**Acceptance criteria (restated from REQ):**
- [ ] Metadata-only interruption recovery still exits zero and prints the staging instruction.
- [ ] Any body difference exits one and does not print the staging instruction.
- [ ] A fenced-example commit line remains significant and cannot be normalized away.
- [ ] Only the first frontmatter block's column-zero top-level commit field is removed for comparison.
- [ ] Temporary-file and producer failures are surfaced under Bash 3.2 and pipefail.

## Implementation Summary

**Files changed:**
- `tools/checks/record-commit-hash.sh` (modified)
- `_dev/tests/record-commit-hash-guards.sh` (modified)
- `actions/version.md` (modified)
- `CHANGELOG.md` (modified)

**What was done:** Compared normalized HEAD and worktree content after removing only the first frontmatter block's column-zero commit field, rejected every other difference, and added body and fenced-example regressions while retaining metadata-only recovery.

## Qualification

Passed - 4 files verified, 5 acceptance requirements traced, P-A-U confirmed.

## Testing

**Tests run:**
- /bin/bash _dev/tests/record-commit-hash-guards.sh - passed
- /bin/bash _dev/tests/contract-regressions.sh - passed

**Red-green validation:**
- RED: Body and fenced-example changes both exited zero and printed Stage and commit it now.
- GREEN: Both differences exit one without the staging instruction, while the existing metadata-only probe still exits zero.

**Existing tests updated:** _dev/tests/record-commit-hash-guards.sh adds explicit negative output assertions for ordinary body and fenced commit text.

## Root Cause

The idempotency branch used only git diff presence and never checked whether the pending difference was actually the metadata field it claimed to recover.

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
**Acceptance:** Pass - normalized metadata-only recovery succeeds and every tested non-frontmatter difference is rejected without staging guidance.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Normalizing both HEAD and worktree into files makes the allowed delta explicit and keeps comparison independent of diff formatting.

**What didn't:** A dirty-path check alone could not distinguish the intended interrupted metadata edit from arbitrary archived-request changes.

**Worth knowing:** The normalization boundary must be both structural and positional: first frontmatter block, column zero, top-level commit field only. Body and fenced examples often contain the same token.
