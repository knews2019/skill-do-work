# Branch merge audit — 2026-09-07

Audited 3 local branch refs and all 17 live origin branch refs after fetching with pruning. The symbolic origin/HEAD is excluded. Remote heads were checked against `git ls-remote --heads origin`. No local or remote branches were deleted; fetch removed 13 stale remote-tracking refs whose remote branches had already been deleted.

Initial main: `57e6431e`. Initial checkout: `do-work-queue-drain-gemini` at `6b247a49`, clean. This is an integration audit of ancestry, unique changes, conflicts, and current verification, not a line-by-line review of all historical commits.

| Branch ref | Audited tip | Unique commits vs initial main | Resolution |
|---|---|---:|---|
| `backup-main-before-secret-scrub` | `68bc8c5e` | 25 | Already contained in initial main; no merge required. |
| `claude/do-work-queue-drain-4ee2xl` | `7b404684` | 0 | Already contained in initial main; no merge required. |
| `do-work-queue-drain-gemini` | `6b247a49` | 15 | Fast-forwarded 15 commits: release guards, atomic download repair, inventory argument validation, verification fixes, and documentation. |
| `main` | `57e6431e` | 0 | Initial main reference. |
| `origin/claude/add-prime-files-Y4QnT` | `4688d578` | 1 | History merged; current work action already guards missing/unknown domains and absent crew files. Retained current release metadata and removed legacy paths. |
| `origin/claude/add-status-command-Lj02v` | `f9867d45` | 0 | Already contained in initial main; no merge required. |
| `origin/claude/architecture-report-format-1o9ylo` | `84ed75aa` | 0 | Already contained in initial main; no merge required. |
| `origin/claude/blocked-state-kanban-8y9099` | `5a1e5018` | 1 | Ported responsive badge truncation to board-cards.js and board.css. Preserved the shared truncation helper used by other badges and current timestamp formatting. |
| `origin/claude/changelog-discrepancies-0174-olf3ko` | `8ee230c0` | 0 | Already contained in initial main; no merge required. |
| `origin/claude/dependency-aware-selection-vg7mE` | `a68ad594` | 0 | Already contained in initial main; no merge required. |
| `origin/claude/do-work-recovery-reqs-vzfue9` | `eef2b8d2` | 0 | Already contained in initial main; no merge required. |
| `origin/claude/fix-prime-lesson-linkback-Yy56Z` | `ec522055` | 1 | Ported missing prime discovery to general crew guidance and both lesson writers. Preserved current satellites, relative links, deferred archival transaction, and named phase references. |
| `origin/claude/merge-main-branch-cjwyjh` | `cc70f934` | 2 | Normal merge. Its tip tree equals its second parent; retired weekly-signal prompts stay retired. |
| `origin/claude/pr-152-followup-pd1ycz` | `12aecaeb` | 0 | Already contained in initial main; no merge required. |
| `origin/claude/queue-execution-speed-nze6ah` | `a53e34ae` | 0 | Already contained in initial main; no merge required. |
| `origin/claude/req-559-560-515-547-564-g82uxl` | `4cda6fe9` | 1 | History merged; retained the newer September restart record instead of restoring an obsolete August handoff. |
| `origin/claude/req-file-reservation-cleanup-8cg4a8` | `b7d9fca7` | 0 | Already contained in initial main; no merge required. |
| `origin/claude/resolve-merge-conflicts-EfNXt` | `702c4e75` | 0 | Already contained in initial main; no merge required. |
| `origin/main` | `57e6431e` | 0 | Initial main reference. |
| `origin/worktree-agent-REQ-547-finalize-refuses-a-req-with-no-checkpoint-entry` | `b3d25c8f` | 0 | Already contained in initial main; no merge required. |
| `origin/worktree-agent-REQ-564-reuse-matching-per-lane-verification-evidence-for-four-hours` | `19103f50` | 0 | Already contained in initial main; no merge required. |

The two restored fixes ship together as 0.305.34. Historical version/changelog conflicts retain current metadata. The final merge records both old branch tips while applying their still-missing behavior at the current paths.

During verification, an external process committed the pending two-parent integration as `68bc8c5e`, then rewrote local history for a cleanup. Main is now `0d721e1b` and the gemini branch is `afb025dc`. Comparing the tested integration tree with the rewritten main shows only one HTML report differs: `ai-reports/2026-09-06_2113_do-work-cli-guide/index.html`. All source, release metadata, and test files are identical. The rewritten history is preserved.

That process also created `backup-main-before-secret-scrub` at `68bc8c5e`. This newly created backup is deliberately excluded: merging it would make the pre-cleanup history reachable again. All other current local and remote branch tips are ancestors of rewritten main.

Validation passed:

- Full heavy maintainer gate in an isolated worktree at `68bc8c5e`: exit 0, 290 seconds, `Maintainer verification passed.` The post-cleanup source and test trees are byte-identical to that tested tree.
- ShellCheck, gofmt, aggregate contracts (including staged installation, updater, installer), and both Go vet passes succeeded.
- Uncached suites: 481 board tests with strict JavaScript probes, 36 strict browser tests, and 829 CLI tests.
- `TestBrowserBehaviorBlockedBadgeFitsItsCard` passed on the integrated code and failed on the pre-port tree at `c6271e56`, covering narrow and wider cards plus full tooltip content.
- Initial shared-checkout gate was invalidated by concurrent history rewrites and is not counted as evidence. The successful isolated run supersedes it.
- A fresh origin fetch before publication found no new remote work. Every current branch tip except `backup-main-before-secret-scrub` is contained in main. The backup remains local and unmerged.
- The audit report is the only subsequent change; no release or executable files changed after validation.

Publication: local integration and audit are committed on main. HTTPS push failed because Git could not obtain a username; `gh auth status` reports no logged-in hosts. An SSH push also failed with `Permission denied (publickey)`. No remote branch was updated by this task. Once Git authentication is restored, publish with `git push origin main`.
