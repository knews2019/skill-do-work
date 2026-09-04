# REQ-527 Independent Review

## Decision Brief

## Review: REQ-527

**Approve with follow-ups** — the destructive path now fails closed for active or unproven worktrees, but dry-run and board verification still disagree with the new cleanup authority in reproducible cases.

Route C | `9ccbc9f2928a50ad155e7e3101880feff7393df0` (range `f57d9176..9ccbc9f2928a50ad155e7e3101880feff7393df0`)

### What's built

- Cleanup Pass 5 now removes a builder lane automatically only when a fresh request-tree read proves one exact request outside `do-work/working/`, the worktree is clean, and its branch or detached head is merged.
- Working, absent, ambiguous, malformed, unreadable, dirty, and unmerged states use the existing `WORKTREE-REQUIRES-CONSENT` path; exact `--discard-worktree` consent remains the only force path.
- The cleanup action's detailed Pass 5 instructions, cleanup guide, and crash-recovery reference state the three-fact rule.

### Decisions / risks for you

- No unsafe default deletion was found. The remaining risk is contradictory operator evidence: dry-run can omit a deletion the real run performs, and board verify can advertise an ambiguous lane as mechanically fixable even though cleanup refuses it.
- The builder's D-01 through D-03 decisions were recovered from both the REQ and readable handback. They are reversible and fit the request: rediscover after Passes 0–4, fail closed on identity evidence, and reuse the existing consent path.

### Findings

**Important:**

- A dry run does not preview the same Pass 5 deletion as a real run when Pass 0 moves a terminal REQ out of `do-work/working/`. `handleCleanup` dry-runs `ApplyPlan`, then `ApplyWorktreeRepairs` rediscovers the unchanged tree (`skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands.go:40-43`; `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go:215,261,268`). In a CLI fixture with a completed REQ under `working/`, `cleanup --dry-run` reported `request_state=working` and required consent while planning the archive move; the real command applied that move and automatically removed the worktree. This violates the guide's promise that dry-run inspects the exact ordered changes and makes destructive preview incomplete — `impact-user-visible` → report only
- Board verify and cleanup disagree for duplicate request identity. Cleanup explicitly returns `ambiguous` from collision evidence (`skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go:373-383`), but the board discards duplicate copies into a single winner (`skills/do-work-board/tools/queue-kanban/model.go:667-701`) and then classifies any winning copy outside `working/` as settled (`skills/do-work-board/tools/queue-kanban/verify.go:912-920`). A built-CLI fixture with one queue and one archive copy of REQ-527 produced both `duplicate-req-id` and `merged-worktree-leftover [fixable]`; cleanup on the same tree produced `WORKTREE-REQUIRES-CONSENT ... request_state=ambiguous`. That fails the requirement that verify and cleanup agree on finishedness — `impact-user-visible` → report only
- The changed cleanup action retains a stale safety restatement: `skills/do-work/actions/cleanup.md:128` still says Pass 5 “only discards unmerged work with the user's explicit consent,” although the new rule requires consent for dirty, still-working, absent, ambiguous, malformed, and unreadable cases too. The detailed Pass 5 section is correct, but this earlier operational summary can still send a reader down the retired merged-only rule — `impact-rule-change` → report only
- Two pre-existing cleanup tests were behaviorally rewritten without the required cross-REQ traceability: `TestMergedCleanBuilderWorktreeIsAutomaticButUnmergedNeedsExactConsent` and `TestWorktreeEnumerationHandlesNULNewlineDetachedAndAbsentConsent` now seed settled REQs (`skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go:141-177,313-333`), but neither the REQ Testing section nor handback labels them as existing tests updated from REQ-409. The behavior is intentional and the assertions remain meaningful; the missing provenance is the review-work traceability defect — `impact-negligible` → report only

**Minor:**

- The live restart handoff action still defines `REMOVABLE` as requiring an archived claim (`skills/do-work/actions/restart-with-parallel-handoff.md:51-55`), while the new cleanup authority defines settled as either queue or archive, i.e. any exact request outside `working/`. This is conservative rather than destructive, but it is another current operational restatement of the redefined finishedness contract — `impact-rule-change` → report only

**Nit:** None.

### Requirements Checklist

- [x] A clean merged worktree whose REQ remains in `do-work/working/` is preserved and reported as unfinished.
- [~] Cleanup and verify use the same anchored REQ id and request-tree evidence — cleanup is anchored and collision-aware, but board verify accepts a deduplicated winner and its name regex does not enforce cleanup's suffix boundary.
- [x] No heartbeat, lock, PID probe, mtime heuristic, claim registry, or time threshold was introduced.
- [x] Clean merged residue with one positively settled REQ still removes mechanically.
- [~] Cleanup action, guide, crash recovery, and adjacent current restatements agree — the detailed action, guide, and work reference agree, but `cleanup.md:128`, board verify, and restart handoff do not fully match.
- [x] Automatic removal remains non-forced; `-D`/`--force` remains behind exact discard consent.
- [x] The existing consent gate was extended rather than relocated.
- [x] The finding-closure ratchet is satisfied for REQ-458 F2: the named clean+merged+working regression was RED before the implementation and GREEN after it.
- [x] All P-A-U boxes are checked, and the five changed files exactly match the declared Scope.

### Acceptance Testing

**Result: Partial**

- `go test -count=1 ./internal/cleanup` — PASS (`10.372s`).
- `go test -race -count=1 ./internal/cleanup` — PASS (`12.153s`).
- Focused board worktree consumer tests covering dirty, working, settled, absent, missing-id, and unreadable states — PASS (`2.663s`).
- `bash _dev/tests/shipped-package-reference-contract.sh` — PASS.
- Sequential `bash _dev/tests/contract-regressions.sh` — PASS; all per-file budgets stayed below 30s.
- Built-CLI active-run fixture — expected finding and preservation: `WORKTREE-REQUIRES-CONSENT ... request_state=working`; the worktree remained present.
- Built-CLI settled fixture — automatic removal remained functional.
- Built board/cleanup duplicate-id fixture — reproduced the Important board disagreement: board printed `[fixable]`, cleanup refused with `request_state=ambiguous`.
- Built cleanup terminal-in-working fixture — reproduced the Important dry-run/apply mismatch: dry-run required consent; the real run moved the REQ and removed the worktree.
- Exact range checks: five declared files only; `git diff --check` passed; no added TODO/FIXME/debug-print artifact was found.

### Restatement Sweep

- Swept the live `skills/` tree for the redefined worktree finishedness/removal rule using `Pass 5`, `merged-worktree-leftover`, merged/clean/settled phrasing, removal/fixability phrasing, and `worktree-agent-*` consumers.
- Confirmed aligned: `skills/do-work/actions/cleanup.md` detailed Pass 5, `skills/do-work/docs/cleanup-guide.md`, and `skills/do-work/actions/work-reference.md` crash and happy paths.
- Recorded stale or divergent live consumers above: `cleanup.md:128`, queue-kanban duplicate/name identity handling and its current prime gloss, and `restart-with-parallel-handoff.md:55`.
- Historical changelogs, archived REQs, and processed KB sources were treated as historical evidence, not current instructions; the older queue-kanban lesson explicitly marks its merged-only model superseded.

### Suggested Additional Testing

- Add a command-level parity test where Pass 0 plans/applies a terminal `working/` move and assert dry-run reports the same Pass 5 removal as the real run.
- Add one shared board/cleanup acceptance fixture for duplicate identity and a prefix-like malformed worktree name; both surfaces should refuse mechanical fixability from the same evidence.

### Scores (on the record — not the headline)

**Overall: 73%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 75% | Core deletion policy is delivered; verify agreement and all-restatement consistency are partial. |
| Code Quality | 80% | Fail-closed authority is clear and surgical; dry-run sequencing and consumer agreement remain wrong. |
| Test Adequacy | 78% | Strong RED/GREEN state matrix and race coverage; misses command-level dry-run parity, board collision agreement, and cross-REQ test provenance. |
| Scope | 100% | Exact declared five-file range, no drift. |
| Risk | Low | No unsafe default deletion reproduced; misleading preview/fixability can still misdirect operators. |
| Acceptance | Partial | Core paths pass end-to-end; two integration contradictions reproduce. |

The raw percentage average is 83.25%; the documented 10-point Acceptance Partial penalty yields 73% after rounding.

### Directive and Guardrail Checks

- No approach directive was assigned.
- Think Before Coding: D-01 through D-03 were recorded.
- Simplicity First: the implementation reuses repository discovery and the existing consent finding; no speculative subsystem was added.
- Surgical Changes: exact Scope match.
- Goal-Driven Execution: literal RED/GREEN evidence and focused/race/module/contract runs were recorded.
- Naming for Reach: no new exported or broadly consumed single-word name was introduced.

### Self-validation

- Rechecked the happy path, working-state refusal, collision state, terminal Pass 0 interaction, detached/branch-only logic, explicit discard bypass, and both prose and executable consumers.
- The second pass is what surfaced the dry-run/apply divergence; it is not inferred solely from code and was reproduced with the built command.
- No additional security, data-loss, force-path, or scope finding was found.

### Follow-ups created

None (5 findings report only).

## Append-ready durable Review block

```markdown
## Review

**Overall: 73%** | 2026-09-04T00:00:44Z

| Dimension | Score |
|-----------|-------|
| Requirements | 75% |
| Code Quality | 80% |
| Test Adequacy | 78% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Dry-run can report consent-required for a terminal REQ still in `working/` while the real run first moves that REQ and then automatically deletes its worktree (`cleanup_commands.go:40-43`; `cleanup_git.go:215,261,268`) — `impact-user-visible` → report only
- Board verify deduplicates ambiguous REQ identity into one settled winner and advertises the worktree `[fixable]`, while cleanup refuses the same candidate as `request_state=ambiguous` (`queue-kanban/model.go:667-701`; `queue-kanban/verify.go:912-920`; `cleanup_git.go:373-383`) — `impact-user-visible` → report only
- `actions/cleanup.md:128` still says only unmerged work needs explicit consent, omitting every newly consent-gated state — `impact-rule-change` → report only
- Existing REQ-409 cleanup tests were behaviorally updated without the required cross-REQ provenance in Testing (`cleanup_git_test.go:141-177,313-333`) — `impact-negligible` → report only

**Minor findings:** `actions/restart-with-parallel-handoff.md:51-55` still equates removable with an archived claim rather than the new exact outside-`working/` settled state — `impact-rule-change` → report only
**Acceptance:** Partial — active-run preservation and settled cleanup pass, but built fixtures reproduce dry-run/apply and board/cleanup contradictions.
**Suggested testing:** 2 items
**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*
```
