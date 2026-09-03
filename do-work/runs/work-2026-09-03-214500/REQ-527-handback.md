# REQ-527 Builder Handback

REQ-527 (preventing cleanup from treating merged builder work as finished) is implemented and committed. Pass 5 now removes a leftover automatically only when all three facts are positive: the worktree is clean, the branch or detached head is merged into the current integration `HEAD`, and one exact readable, unambiguous REQ is settled outside `do-work/working/`. Every missing or negative fact uses the existing consent-required finding without mutation. Exact `--discard-worktree` consent still uses the established forceful removal path.

## Branch and Commit

- Branch: `worktree-agent-REQ-527-teach-cleanup-pass-5-that-merged-is-not-finished`
- Commit: `643414128b9282ff5cfdca0e750335e0813a650c`
- Commit subject: `[REQ-527] protect active merged builder lanes`
- Worktree after commit: clean

## Exact Manifest

- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go`
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go`
- `skills/do-work/actions/cleanup.md`
- `skills/do-work/docs/cleanup-guide.md`
- `skills/do-work/actions/work-reference.md`

No `do-work/` path was read or written in the builder worktree. REQ-539, REQ-563, REQ-564, lifecycle metadata, run manifests, release files, and integration state were not touched.

## Implementation

- `ApplyWorktreeRepairs` performs a fresh `repositorymodel.DiscoverRepository` observation when Pass 5 starts, after earlier cleanup passes have applied.
- Worktree attribution now accepts the exact anchored `worktree-agent-REQ-<digits>` prefix and rejects missing digits or a non-hyphen suffix boundary.
- `worktreeRequestState` fails closed across `working`, `absent`, `ambiguous`, `malformed`, and `unreadable`; only one exact parseable REQ in `queue` or `archive`, with no collision evidence, is `settled`.
- The existing `WORKTREE-REQUIRES-CONSENT` result now reports cleanliness, merged state, and request state from the same decision. Its exact recovery argv remains `do-work-cli cleanup --discard-worktree <name>`.
- Cleanup action, user guide, and crash-recovery reference now state the same three-fact rule and preserve unattended automatic removal only for proven residue.

## RED / GREEN Evidence

RED was captured before production edits with:

`go test ./internal/cleanup -run 'TestMergedCleanBuilderWorktree(IsAutomaticButUnmergedNeedsExactConsent|RequiresSettledRequestEvidence)|TestWorktreeEnumerationHandlesNULNewlineDetachedAndAbsentConsent'`

Result: failed in 2.00s. Each literal clean+merged fixture for a working, absent, duplicate, malformed, or unreadable REQ was removed with no finding. This reproduced the acceptance defect.

GREEN after implementation:

`go test ./internal/cleanup -run 'TestMergedCleanBuilderWorktree(IsAutomaticButUnmergedNeedsExactConsent|RequiresSettledRequestEvidence)|TestWorktreeEnumerationHandlesNULNewlineDetachedAndAbsentConsent'`

Result: PASS in 1.84s. Each unproven state remained byte/path intact with one `WORKTREE-REQUIRES-CONSENT` finding; exact named discard then removed it. Existing clean+merged+archived automatic removal, unmerged refusal, explicit discard, and detached/NUL-safe enumeration remained green.

## Verification

- `go test -count=1 ./internal/cleanup` — PASS, final run 12.99s.
- `go test -race -count=1 ./internal/cleanup` — PASS, 14.06s.
- `go vet ./...` in `skills/do-work/tools/do-work-cli` — PASS, 4.44s.
- `go test -count=1 ./...` in `skills/do-work/tools/do-work-cli` — PASS, 102.10s.
- `bash _dev/tests/shipped-package-reference-contract.sh` — PASS, 6.81s.
- `bash _dev/tests/contract-regressions.sh` — PASS, final run 19.37s.
- `git diff --check` and cached diff check — PASS.

The first contract run was intentionally parallel with the full Go module and its SessionStart fixture crossed the per-file timing budget at 47s, with no semantic fixture failure. Standalone `bash _dev/tests/session-start-hook-behavior.sh` passed in 11.65s, and the complete sequential contract rerun passed with that fixture at 15s.

## Prime and Lesson Reads

Read before implementation:

- `CLAUDE.md`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`
- `_dev/primes/prime-action-files.md`
- `_dev/primes/lessons-action-files.md`
- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/communication-style.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/anti-slop.md`

Applied lesson: repository collision evidence is first-class, so a single typed record does not authorize deletion when `CollisionEntries` still proves ambiguous ownership. Applied action-file lesson: the changed rule was swept through the execution action, guide, and crash-recovery reader rather than left only at the primary implementation.

## P-A-U Evidence

- Producer: `repositorymodel.DiscoverRepository` supplies rooted typed request files, exact tree section, parse/read failures, normalized identity, and collision evidence.
- Authority: `cleanup.ApplyWorktreeRepairs` joins that evidence with fresh Git cleanliness and ancestry observations, and owns the automatic-removal versus consent-required verdict.
- Users: `actions/cleanup.md` invokes the authority; `docs/cleanup-guide.md` explains operator behavior; `actions/work-reference.md` consumes the same rule during crash recovery.

## Decisions

- Kept one fresh repository discovery inside Pass 5 instead of passing the command's pre-apply snapshot. This allows a terminal working REQ moved by Pass 0 to qualify and prevents stale lifecycle evidence from deciding deletion.
- Used the existing consent finding and exact discard token for every new refusal state. No liveness signal, age heuristic, new prompt, or force mode was added.
- Treated both `queue` and `archive` as settled locations because the durable contract is exact positive presence outside `do-work/working/`; malformed, unreadable, absent, or multiply claimed identity never qualifies.

## Risks and Seams

- Pass 5 observes repository and Git facts before it issues the existing removal commands. Concurrent external mutation remains outside the cleanup transaction, as before; the new rule only narrows automatic deletion.
- No integration seam is required. Merge the exact commit, then rerun the focused cleanup package and shipped contract checks on the merged tree.

## Discovered Tasks

None.

## Scope Issues

None. The implementation stayed within the declared five-file scope.
