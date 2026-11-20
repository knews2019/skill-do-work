# REQ-754 implementation handback

Implemented and committed. Final focused tests and the final repository gate are green. No merge, release bump, claim/finalization command or do-work run was performed. The prepared checkout remains available and clean.

## Branch and commits

- Checkout: `/Users/t2/.t3/worktrees/skill-do-work2/t3-018643a3-a57e084b2c7761c3`
- Branch: `t3/second-opinion-isolation-cleanup`
- Base: `8b7aa44c912dcb5ed86e2db581f726d8d381ae5c`
- Implementation: `60186a38040b63681065ec18c66a60604c14391a` — [REQ-754] Isolate second-opinion subprocesses and clean up worktrees
- Final HEAD: `e7bd242b3718fd522acc424e76f042a9d10f790c` — [REQ-754] Retain finished pages when real bundle revalidation fails
- Integrate both commits, or merge this branch. No push or PR was requested/performed.

## Full file manifest

| Action | File | Verified change |
|---|---|---|
| M | `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` | Isolate probes/generation; separate real and isolated inode snapshots; retain finished stages; publish through the original real-bundle tail; remove real-checkout outside-status checks from subprocess, session prepare/finalize and the shared tail. |
| A | `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go` | Own detached temporary checkout; materialize exact current bundle; rooted exclusive copy-back; independent bounded cleanup after owned-process reaping; report cleanup failures. |
| M | `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` | Remove the prior REQ-698 expectation that pre-existing outside dirt must refuse. Its replacement is the dirty tracked README positive/negative regression. |
| A | `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree_test.go` | Dirty tracked README, detached HEAD, isolated outside writes, prompt paths, exact additions/edits/deletions/empty directories, allowed assets, refusals, concurrent real changes/collisions, session final publication, cancellation/deadline/SIGINT, setup cancellation/failure, descendant reaping, cleanup fault reporting and disk/registration assertions. |
| M | `skills/do-work-toolbox/actions/ai-report.md` | Describe isolated subprocess paths and outside_paths, session preservation without outside status, retained output, cleanup and unchanged backend permissions. |
| M | `skills/do-work-toolbox/docs/ai-report-guide.md` | Document current-bundle copying, isolated outside_paths, cleanup and worktree limitations; remove the clean-real-checkout prerequisite. |

The extra test file is explicitly permitted by the task. Tests were split to honor the repository's per-file runtime budget. Final diff from base: 6 files, 612 insertions, 42 deletions. No release metadata or do-work/ path is in either commit. No do-work/ file in this checkout was written; this handback is the sole parent-tree write.

## RED → GREEN

Commands below ran from `/Users/t2/.t3/worktrees/skill-do-work2/t3-018643a3-a57e084b2c7761c3/skills/do-work/tools/do-work-cli` unless noted otherwise.

### Required dirty tracked README regression

Before implementation, added the fake-codex regression and ran:

```sh
go test ./internal/toolboxcommands -run '^TestSecondOpinionWorktreeDirtyReadme$' -count=1
```

RED (exit 1), before source changes:

```text
--- FAIL: TestSecondOpinionWorktreeDirtyReadme (0.25s)
    --- FAIL: TestSecondOpinionWorktreeDirtyReadme/unrelated_dirt (0.13s)
        report_second_opinion_worktree_test.go:60: dirty README blocked publication: map[actual_model: backend:codex banner_state:not-inserted judge_findings:[] judge_verdict:not-run outside_paths:[README.md] page: reason:outside-bundle changes: README.md reply: requested_model: seen_models:[] state:failed]
    --- FAIL: TestSecondOpinionWorktreeDirtyReadme/backend_outside_write (0.12s)
        report_second_opinion_worktree_test.go:57: finished refusal page not retained: map[actual_model: backend:codex banner_state:not-inserted judge_findings:[] judge_verdict:not-run outside_paths:[README.md] page: reason:outside-bundle changes: README.md reply: requested_model: seen_models:[] state:failed]
FAIL
FAIL	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	0.587s
FAIL
```

The same command against final HEAD is GREEN (exit 0):

```text
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	1.453s
```

The success case proves detached HEAD, successful page/banner publication despite dirty tracked README, byte-identical real README, no extra Git worktree registration and no scratch directory. The outside-write case writes `../../README.md`, reports README.md in outside_paths, keeps the finished page in the real bundle, preserves real README bytes and removes the isolated worktree.

### Final-review retention regression

Strengthened the existing real-change cases to require the finished page to survive a concurrent real index replacement or asset collision. Before the follow-up fix:

```sh
go test ./internal/toolboxcommands -run '^TestSecondOpinionWorktreeCurrentBundleAndRefusals/real' -count=1
```

RED (exit 1):

```text
--- FAIL: TestSecondOpinionWorktreeCurrentBundleAndRefusals (0.95s)
    --- FAIL: TestSecondOpinionWorktreeCurrentBundleAndRefusals/real_original_replacement (0.49s)
        report_second_opinion_worktree_test.go:137: finished page lost: map[actual_model:actual-model-1 backend:codex banner_state:not-inserted judge_findings:[] judge_verdict:not-run outside_paths:[] page: reason:backend changed existing bundle path: index.html reply: requested_model: seen_models:[] state:failed]
    --- FAIL: TestSecondOpinionWorktreeCurrentBundleAndRefusals/real_asset_collision (0.47s)
        report_second_opinion_worktree_test.go:137: finished page lost: map[actual_model:actual-model-1 backend:codex banner_state:not-inserted judge_findings:[] judge_verdict:not-run outside_paths:[] page: reason:backend created an unprefixed or non-file asset: codex-actual-model-1-chart.svg reply: requested_model: seen_models:[] state:failed]
FAIL
FAIL	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	1.290s
FAIL
```

After the fix, normal import/publication still refuses real-bundle changes. Only the exclusive recovery staging leaf is retained in this case. The final full focused run covers these tightened assertions:

```sh
go test ./internal/toolboxcommands -run '^TestSecondOpinion' -count=1
```

```text
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	33.293s
```

## Verification

- Final focused second-opinion tests: exit 0, output above. Existing Codex/Gemini/AGY/Claude adapters, flags, model identity, unverified pages, session preservation, tree publication, judge/banner gating and process ownership remain covered.
- Whole toolbox package race check: `go test -race ./internal/toolboxcommands -count=1`, exit 0: `ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	42.809s`. After the final retention change, reran its affected matrix with `go test -race ./internal/toolboxcommands -run '^TestSecondOpinionWorktreeCurrentBundleAndRefusals$' -count=1`, exit 0: `ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	7.004s`.
- Go floor: `GOTOOLCHAIN=go1.24.0 go build ./...`, exit 0, including a final-HEAD rebuild. `GOTOOLCHAIN=go1.24.0 go test ./internal/toolboxcommands -run '^TestSecondOpinion' -count=1`, exit 0: `ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	33.164s` (before the final retention-only follow-up).
- Windows compatibility: `GOOS=windows GOARCH=amd64 go build ./...`, exit 0, including a final-HEAD rebuild. `GOOS=windows GOARCH=amd64 go vet ./internal/ownedprocess ./internal/gittransaction ./internal/toolboxcommands ./internal/heavyverification`, exit 0.
- Native `go vet ./...`, gofmt and diff whitespace checks: pass. Final gate includes native vet and formatting on final code.
- Final repository gate: exit 0. Exact invocation from checkout root:

```sh
REQ754_REPOSITORY="/Users/t2/.t3/worktrees/skill-do-work2/t3-018643a3-a57e084b2c7761c3" bash /tmp/REQ-754-maintainer-verify.sh > /tmp/REQ-754-gate-final.log 2>&1
```

This is a temporary copy of `_dev/tests/maintainer-verify.sh`, with only its repository-root resolution set to this checkout, timing output redirected to `/tmp/REQ-754-durations.tsv`, and `run_stage_with_evidence` overridden to run every stage freshly without reading/revoking/recording do-work/ stage evidence. The gate normally hardcodes writes under do-work/, which the parent explicitly prohibited. No checks were removed: shell lint, formatting, aggregate contracts, both modules' vet and uncached budgeted test stages ran. This is not a claim that the unmodified gate was invoked in this checkout. The parent can run the canonical gate after integration.

```text
maintainer-verify: stage queue-kanban-fast-tests: EXECUTING (builder, no state writes)
go-test budget: module=/Users/t2/.t3/worktrees/skill-do-work2/t3-018643a3-a57e084b2c7761c3/skills/do-work-board/tools/queue-kanban wall=44s tests=427 slowest-file=strict_behavior_regression_test.go:21.70s limit=<30s
maintainer-verify: do-work-cli go vet
maintainer-verify: do-work-cli uncached tests
maintainer-verify: stage do-work-cli-fast-tests: EXECUTING (builder, no state writes)
go-test budget: module=/Users/t2/.t3/worktrees/skill-do-work2/t3-018643a3-a57e084b2c7761c3/skills/do-work/tools/do-work-cli wall=66s tests=1011 slowest-file=internal/finalization/finalization_recovery_test.go:17.85s limit=<30s
maintainer-verify: gate wall 141s
Maintainer verification passed.
```

Final per-file second-opinion durations (30-second limit):

```text
20261119T235429Z-34083	skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go	17.34	0
20261119T235429Z-34083	skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_tree_test.go	3.63	0
20261119T235429Z-34083	skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree_test.go	14.06	0
```

The default gate explicitly skips the heavy staged-skills/updater/installer lanes; no heavy-tier completion is claimed. Full raw logs are `/tmp/REQ-754-red.log`, `/tmp/REQ-754-dirty-green.log`, `/tmp/REQ-754-green-final.log`, `/tmp/REQ-754-retention-red.log`, `/tmp/REQ-754-race.log`, `/tmp/REQ-754-retention-race.log`, `/tmp/REQ-754-go124.log` and `/tmp/REQ-754-gate-final.log`. Git's NUL-delimited worktree inventory appears in captured stderr in failure logs, so use a binary-aware reader or `rg -a` when inspecting those logs.

## P-A-U verification

- [x] PLAN: Followed the authoritative brief: detached owned worktree, exact current bundle, isolated snapshot, worktree-only outside check, exclusive import, real publication tail, independent bounded cleanup.
- [x] APPLY: Only the six manifest files changed. Test split was authorized. Backend argv flags and the shared per-file publication sequence remain unchanged.
- [x] UNIFY: Reviewed every manifest file and the complete diff/stat. Checked the real/isolated snapshot boundary, deferred cleanup ordering, all session/publication entry points, regular-file allowlist, recovery stage retention, exclusive collision behavior, native/model flags and both shipped prose surfaces. Final formatting, vet, focused tests, affected race tests and fresh gate pass. No debug artifacts, release changes or do-work/ changes are in the diff; checkout is clean.

## Decisions

- D1: Replace the isolated bundle rather than overlay it, so deleted tracked bundle files stay deleted; snapshot the isolated files separately because inode identity cannot cross the copy.
- D2: Reuse rootedfs and rootedPublishFile for confined materialization and exclusive copy-back. Keep the original real-bundle publication tail and recoverable per-file results.
- D3: Reuse the existing owned-process capture for Git setup/cleanup as well as backend generation. Git checkout/filter children must be reaped too. Disable checkout hooks for the launcher-owned worktree.
- D4: Canonicalize the temporary directory before comparing it with Git registration. The first cleanup regression exposed macOS `/var` versus `/private/var` aliasing; it is fixed and pinned by disk/registration assertions.
- D5: Cleanup is deferred before setup starts and uses a new 30-second background context. Never prune unrelated worktrees. A cleanup error changes the run to failed and is reported, even if publication already happened.
- D6: A real-bundle preservation failure still allows exclusive retention of the finished staging page, but no normal asset import or publication. This closes the final-review data-loss case without overwriting a concurrent writer.
- D7: Replace the REQ-698 pre-existing-dirt refusal assertion with the requested positive regression. The old session lesson's outside-status guarantee is intentionally superseded by REQ-754, while its digest-based preservation rule remains.

## Guidance and lessons

Read CLAUDE.md; the general, coding-guardrails, shared-principles, communication-style, backend, testing and anti-slop crews; CLI, action-files, releases and shell-commands primes; and the toolbox routing SKILL.md. Required `_dev/primes/lessons-releases.md` was read in full. Read the relevant CLI lesson entries for native backend identity, session ownership, process reaping, worktree cleanup, test-file budgets and toolchain/rootedfs behavior, plus its Verify and Scoped command contracts sections. Read the relevant action lesson material on alternate-writer drift, templates, accepted variants and agent compatibility. The oversized CLI/action satellites were not loaded in full. No listed guidance file was missing.

Lessons for the parent: Git registers canonical worktree paths; cleanup must compare canonical paths. Teardown must own an uncancelled context after process reaping, including Git setup failures. A preservation refusal should not erase finished new page bytes when an exclusive recovery leaf can still be created.

## Limitations and handoff

No live paid/model backends were invoked; native protocol fixtures exercise the adapters. New tests stub the render judge where a pass is needed; existing judge behavior is unchanged. Windows was cross-built/vetted, not executed. Worktrees are isolation for ordinary relative writes, not a security boundary against absolute-path writes or shared Git refs/stash. Existing backend permission flags remain in force.

If Git or the filesystem refuses teardown, cleanup is reported as failed and may leave a stale registration; the injected-failure test confirms the failure is not hidden. If the real destination cannot accept an exclusive staging file, retention itself can fail and reports that error; existing files are never overwritten. No transaction-wide rollback was added.

Parent owns merge, canonical verification, release/finalization and checkout cleanup. No unrelated discovered task was introduced.
