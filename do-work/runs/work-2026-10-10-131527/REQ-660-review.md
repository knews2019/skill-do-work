# REQ-660 review (do-work-cli worktree new, status, merge and cleanup)

Orchestrated mode, Route B, standard review. Diff read from the merge range `bcf12cf3..08f005e2` (6 files, 789 insertions, 3 deletions). The reviewer edited no file except this one, committed nothing, and ran no command against the repository root. Every acceptance run used throwaway fixture repositories under the session scratchpad. They are removed.

## Review: REQ-660

**Approve with follow-ups.** The four subcommands work as specified and never force. One data-loss path needs a fix before the link feature is advertised to consumers.
Route B | merge `08f005e2` (builder commit `cd557b36`)

### What's built
- `do-work-cli worktree new|status|merge|cleanup` runs the fan-out-reference.md lifecycle: name derivation with the `-2` collision suffix, a sibling `<repo>-worktrees/` directory, optional `do-work/worktree-links` symlinks, the queue guard before `git merge --no-ff --no-commit`, a `[REQ-N]` merge commit, and cleanup with `git worktree remove`, `git branch -d` and `git worktree prune`. No subcommand passes `--force` or `-D`.
- What is still missing: `merge` does not guard against a builder that committed one of the links. This was reproduced, and it deletes the main tree's ignored directory (F1).

### Decisions / risks for you
- F1 is `impact-critical`. The integrator should create the follow-up REQ. The smallest fix is one extra guard in `merge`, plus one warning in `new`, plus one sentence in the docs. Each is described under F1.

### Findings

**Important:**
- F1. `worktree merge` merges a builder commit that contains a configured link. This destroys the main tree's copy of that path. Reproduced in a fixture: `worktree-links` has `node_modules`, `.gitignore` has `node_modules/` (the usual form, and the REQ's own example), the builder runs `git add -A`, then `worktree merge REQ-7` exits 0. The real ignored `node_modules/` directory in the main tree is deleted and replaced by a committed symlink that points at itself ("Too many levels of symbolic links"). Cause: a `dir/` ignore pattern does not match a symlink, so the link stays untracked and `git add -A` stages it. Git treats ignored files as expendable on merge. The REQ assumed "linked paths are expected to be git-ignored in the worktree", but the command never checks this. The qualify and scope-drift gates read the range only after the merge, so the deletion happens before any gate runs. Suggested fix: (a) in `mergeBuilderBranch`, next to the queue guard, refuse when `git diff --name-only <pre>...<name> -- <each configured link path>` prints anything (`worktree_lifecycle.go:291-298`). (b) In `createWorktree`, warn when `git -C <worktree> check-ignore -q <link>` says a new link is not ignored. (c) In `fan-out-reference.md` → Where worktrees live, add one sentence: the ignore line must match the link itself (`node_modules`, not `node_modules/`). RED case: the fixture sequence above. — impact-critical → follow-up REQ for the integrator to create (the reviewer created none, per the brief)
- F2. A remediation re-merge reports the wrong `<pre>`. Reproduced: the second `worktree merge` on the same branch returned `pre` = the first merge's hash, not the first `<pre>`. `fan-out-reference.md` → Remediation re-merges says to keep the first `<pre>` (and to append a reason to the merge subject). The new paragraph after step 4 (`fan-out-reference.md:83`) says the command "runs steps 1 to 4" with no remediation caveat. An integrator who trusts the reported `pre` gets the fix-only range, which is the failure that section describes. Fix: one clause in that paragraph ("on a remediation re-merge keep the first `<pre>`; the reported `pre` is this merge's parent only"), or a flag. Restatement-sweep finding. — impact-user-visible → report only
- F3. Other readers count the links that `worktree new` creates as uncommitted work. This extends the builder's Discovered Tasks 1 and 2 with two code readers the builder did not list: `queue-kanban verify` `worktreeHasUncommittedWork` (`skills/do-work-board/tools/queue-kanban/verify.go:1109-1116`), whose message "holds uncommitted changes that are in no commit" becomes false for a linked leftover, and the restart-handoff survey (`skills/do-work/tools/do-work-cli/internal/corehelpers/handoff.go:87`). Also affected: cleanup Pass 5 `worktreeClean` (`internal/cleanup/cleanup_git.go:345`) and the Crash Recovery prose. All four fail in the safe direction: a linked leftover is reported and needs consent, and nothing is deleted. Restatement-sweep finding. — impact-user-visible → report only

**Minor:**
- F4. `worktree status` fails outright (exit 2, `WORKTREE-GIT-FAILED`, no `next_argv`) when any `worktree-agent-*` record is prunable, meaning its directory is gone. One stale leftover hides the table for every builder. `cleanup` of that REQ fails the same way until a manual `git worktree prune`. Reproduced. Fix: skip the record as a row marked missing, or run `git worktree prune` first as cleanup Pass 5 does (`worktree_lifecycle.go:241-248`). — impact-user-visible → report only
- F5. `cleanup` removes the owned links before `git worktree remove`. A git refusal that the preflight does not check, such as a locked worktree (reproduced: outcome failure, links already gone), leaves the worktree without its links. This is a side effect before a refusal. Fix: check `locked` in the porcelain record during preflight, or remove the links last and retry the remove once (`worktree_lifecycle.go:354-367`). — impact-negligible → report only
- F6. The link config accepts `do-work` or `do-work/...` lines. Reproduced in a consumer whose `do-work/` is untracked: `new` symlinks the live queue into the builder tree, which breaks "State stays home". Fix: refuse such lines in `readLinkConfig` with the existing `WORKTREE-LINK-CONFIG-INVALID` code. — impact-negligible → report only

**Nit:**
- F7. `dirtyPaths` assumes a rename's source path follows its entry (`index++`), but `splitZero` sorts the entries. For a staged rename, the evidence paths can be shifted or cut (for example `Makefile` reported as `efile`). The dirty/clean answer stays correct, because the rename entry itself is always counted (`worktree_lifecycle.go:464-475`). — impact-negligible → report only
- F8. Anti-bloat count. Beyond what the REQ named: 20 finding codes, a 12-field result block plus a 6-field row type, a 24-line text renderer, the MERGE_HEAD precheck (D-13) and the link-config validation (D-15). Flags are only `--from` and `--name`, which the Assumptions named. No extra files. Five tests, each tied to a named GREEN condition, none decorative. Judged unneeded: `links_skipped` repeats the `WORKTREE-LINK-SKIPPED` findings, `empty_hand_back` repeats the `WORKTREE-EMPTY-HAND-BACK` code, and `WORKTREE-LINK-CONFIG-UNREADABLE`/`-INVALID` plus `WORKTREE-DISCOVERY-FAILED`/`WORKTREE-LINK-FAILED` could be two codes instead of four. Small cost. No Constraint forbids any of it. — impact-negligible → report only

**Builder's Discovered Tasks, judged:** Discovered Task 1 (Pass 5 reads links as dirt) is report-only. It does not make the new command unusable, because `worktree status` and `worktree cleanup` filter their own links (D-11), and the fixture run `new` → commit → `merge` → `cleanup` left nothing behind. The other readers fail toward asking for consent, never toward deleting. F3 widens the list. Discovered Task 2 is the same family. Discovered Task 3 (empty `<repo>-worktrees/` parent) is correctly negligible, because the hand path leaves it too.

### Requirements Checklist

- [x] 1. `new` derives the name from the filename slug as a text operation, adds `-2` on a collision without deleting, creates the worktree from the main tree's branch in the sibling directory, links configured paths, and reports the operative name. Delivered. The brief skeleton is skipped under the REQ's own Assumption (D-03), since there is no fixed template.
- [x] 2. Optional `do-work/worktree-links`, one path per line, no links when the file is absent. Delivered (F6 is a gap).
- [x] 3. `status` rows show ahead/behind, dirty or clean, and last-commit age. Delivered (F4 is a robustness gap).
- [x] 4. `merge`: clean index, `<pre>`, queue guard, `--no-ff --no-commit`, `[REQ-N]` commit, `<pre>` and `<merge_hash>` reported. Empty hand-back refuses with no commit. A conflict stops with the merge in progress and the paths listed. Delivered (F1, F2).
- [x] 5. `cleanup`: links, `worktree remove`, `branch -d` from the integration branch, `prune`. Unmerged or dirty refuses first. Delivered (F5).
- [x] 6. Never `--force`, never `-D`. Every refusal reports and stops. Delivered: no `"--force"`, `"-D"` or `"-f"` argv in the file, and `"-d"` appears only in `branch -d`.
- [x] 7. `fan-out-reference.md` names each command next to its hand sequence, and the hand path stays valid. Delivered (F2 is a missing caveat).
- [ ] 8. Release. N/A at review time; the integrator's finalization does it.
- [x] Constraints: rules unchanged, no new frontmatter field or status, queue guard enforced. UR-145 acceptance item 3 passes.

### Acceptance Testing

**Result: Pass.** Stages covered: Implementation and Integration (merged tree at `08f005e2`, fixture repositories, built binary). Not covered: Deployment, Live acceptance.
- `go test -count=1 -run 'TestWorktree' ./internal/cleanup/`: exit 0, all six matching tests PASS (five new tests plus the existing enumeration test).
- Fixture: `new` (slug `Demo Name's Slug!` became `worktree-agent-REQ-7-demonamesslug`; `node_modules` and `env.vars` linked; a missing path was skipped with a warning) → builder commit → `status` (ahead 1, clean) → `merge` (two-parent `[REQ-7]` merge, `pre`/`merge_hash` match git) → `cleanup` (exit 0; no worktree, branch or link left; the main-tree originals intact).
- Refusals, each exit 1 with no side effect: empty hand-back, queue guard (`do-work/working/x.md`), unmerged cleanup, dirty cleanup (links, worktree and branch kept), staged index, several matching names without `--name`, detached main tree, unknown subcommand, `--from --force`, `--from -D`, `--name --force`. A conflict returns outcome `findings` with `conflicted_paths: [README.md]` and the merge left in progress.
- Every refusal `next_argv` seen is read-only: `git status --short [--branch]`, `git log`, `git diff --name-only`, `git diff --cached --name-only`, `git branch --list`. Usage refusals carry none.
- Reproduced defects: F1 (data loss), F2, F4, F5, F6.

### Suggested Additional Testing

- Deployment: unassessed. Run `tools/do-work-cli.sh --repo-root <consumer> worktree new REQ-N` from an installed consumer package to confirm the wrapper routes the new verb.
- Live acceptance: unassessed. A JS consumer run where a coordinator uses `new`/`merge`/`cleanup` for a real wave with `node_modules` linked, after F1 is fixed.
- Environment: Windows, where `os.Symlink` needs privileges or developer mode. Expect `WORKTREE-LINK-FAILED` after the worktree was already created.
- Edge case: `worktree new` started from inside a linked builder tree (D-16 says it resolves to the main tree; no test pins this).

### Scores (on the record, not the headline)

**Overall: 60%** (dimension average 87.5%, capped at 60% by Risk = Critical)

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | All seven in-review requirements delivered; F2 is a missing doc caveat |
| Code Quality | 80% | Clean reuse of the cleanup helpers and preflight order; F1, F4, F5, F7 |
| Test Adequacy | 80% | Five focused tests with red-green and a queue-guard mutation check; no conflict, prunable or link-commit test |
| Scope | 95% | 6 declared = 6 touched; small duplicate fields and codes (F8) |
| Risk | Critical | F1: a reproduced data-loss path in the advertised new → merge flow |
| Acceptance | Pass | Implementation and Integration stages, fixture repositories |

### Follow-ups created
- None by the reviewer. F1 is impact-critical, and the integrator will create its follow-up REQ.

Self-validation: re-checked that F1 needs no command misuse beyond a builder's `git add -A`, which no action forbids (`actions/work.md` builder rules name no staging form). Confirmed that the hand path has the same hazard when links are made by hand, but the command is now the documented way to create them. Confirmed that the P-A-U boxes are all `[x]`.

## Review

**Overall: 60%** | <INTEGRATOR-TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 80% |
| Test Adequacy | 80% |
| Scope | 95% |
| Risk | Critical |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `worktree merge` merges a builder commit containing a configured link: with a `node_modules/` ignore line and `git add -A`, the main tree's ignored directory is deleted and replaced by a self-looping committed symlink (reproduced; fix = link-path guard beside the queue guard at `worktree_lifecycle.go:291-298`, a `check-ignore` warning in `new`, one doc sentence) — impact-critical → follow-up REQ for the integrator to create (reviewer created none, per brief)
- F2 the remediation re-merge reports this merge's parent as `pre`, and the new `fan-out-reference.md:83` paragraph omits the keep-the-first-`<pre>` rule — impact-user-visible → report only
- F3 linked worktrees read as dirty in queue-kanban `verify.go:1109-1116` (false "uncommitted changes" message), `corehelpers/handoff.go:87`, cleanup Pass 5 `cleanup_git.go:345`, and the Crash Recovery prose; all fail toward consent — impact-user-visible → report only

**Minor findings:** F4 `status`/`cleanup` exit 2 on a prunable record and hide every row — impact-user-visible → report only; F5 `cleanup` removes links before a locked-worktree refusal from `git worktree remove` — impact-negligible → report only; F6 link config accepts `do-work` paths and links the live queue into a builder tree — impact-negligible → report only; F7 (nit) `dirtyPaths` rename skip is wrong after `splitZero` sorts, so evidence paths only — impact-negligible → report only; F8 (nit) anti-bloat: `links_skipped`/`empty_hand_back` repeat findings, four failure codes could be two — impact-negligible → report only
**Acceptance:** Pass — Implementation and Integration stages: focused tests exit 0, and the fixture new → commit → status → merge → cleanup, empty hand-back, queue guard, dirty and unmerged cleanup, conflict, and bad-argument refusals all behave as specified, with read-only `next_argv` only
**Restatement sweep:** redefined the builder worktree add / hand-back merge / cleanup hand sequences (now wrapped by `do-work-cli worktree`), the new `do-work/worktree-links` file, and what "dirty" means for a linked builder worktree; stale: `fan-out-reference.md:83` (remediation `<pre>`, F2), Where worktrees live (ignore-line form, F1), dirt readers in `verify.go`, `handoff.go`, `cleanup_git.go` and Crash Recovery (F3); consistent: `crew-members/background-agents.md:186`, `actions/cleanup.md` Pass 5, `actions/work.md:574-575`, `actions/restart-with-parallel-handoff.md:55`, `docs/cleanup-guide.md`
**Suggested testing:** 4 items
**Follow-ups created:** None by the reviewer. F1 is impact-critical and goes to the integrator to create (7 other findings report only)

*Reviewed by review-work action*
