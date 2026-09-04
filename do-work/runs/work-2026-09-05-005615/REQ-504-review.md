# Review: REQ-504

**Request changes** — the public recovery and checkpoint commands work for canonical records, but a supported legacy checkpoint can become unrecoverable and then lose its visible ownership evidence during refresh.

Route C. Reviewed `773787b74acddfdfc4c16498a89d99a5cc3ab716..f412a8411057d0a833df5584657161008f315b84`. This independent orchestrated review used the queue request as read-only evidence; the orchestrator must resolve its exact path and evidence again after claim. No request, follow-up, or lifecycle record was written.

## What's built

The migration adds finalization-first `recover`, explicit takeover authority, structural checkpoint evidence, all-entry claim removal, and checkpoint-only `advance --checkpoint`. It removes the split shell recovery fixture and large action algorithms, replacing them with command boundaries and behavioral tests. The ordinary saved-revision advance remains read-only; later requests deliberately extend it.

## Findings

**Important F-01 — supported legacy checkpoint evidence disagrees across discovery, recovery, and refresh.** `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go:347` deliberately discovers heading-less claim entries. `internal/requeststate/state_apply.go:854` (current HEAD; line 832 at the reviewed revision) requires the canonical section before removing any entry, so `recover --assume-sole-authority` refuses the same record with `RECOVER-CLAIM-CHECKPOINT-EVIDENCE`. Next, `internal/lifecycleadvance/checkpoint_commands.go:99` appends an empty canonical section below the legacy record. Discovery subsequently selects that empty section and stops returning the retained record, although refresh reported `preserved_claims: 1`. The bytes survive, but the live ownership evidence disappears from the public result. The direct public replay below confirmed both failures; the implicated functions remain unchanged at current HEAD. This fails the requirement to recover every supported checkpoint shape and preserve structurally observed live claims. — **impact-critical**; return to the orchestrator for the single remediation and fold-first decision, with no follow-up created by this reviewer.

No additional current findings. The saved revision's `work-reference.md:703` still restated automatic own-label recovery despite the new explicit-authority policy; later work removed that operative paragraph. It is historical scope context, not a remaining repair request.

## Requirements checklist

- [x] Public finalization-first recovery, coherent claim-only topology, typed ownership decisions, and constant authority argv are implemented and exercised.
- [ ] Every supported checkpoint shape recovers and remains visible after refresh: canonical labelled, unlabelled, duplicate, and hostile-label cases pass; heading-less legacy evidence fails F-01.
- [x] Authorized canonical recovery removes all same-request entries, including aliases and continuations, preserving unrelated entries and project dirt.
- [x] The saved ordinary advance is read-only; checkpoint mode changes only its checkpoint target in the passing canonical case.
- [x] Step 10 is one loop paragraph plus a context-wipe principle; Crash Recovery and session-checkpoint mechanics collapse to principles and command boundaries.
- [x] Inherited finalization-ownership, Step 8 naming, commit, handoff, and guide changes appear in the declared scope. The saved selection exception is explicitly deferred to REQ-505.
- [x] All four migration parts are present: owning commands, deleted prose, retired shell predicate/fixture registration, and replacement Go behavior tests. The original sentence-predicate RED premise was stale, and the Plan documents the replacement public-boundary RED cases.
- [x] The diff contains the declared 26 implementation paths plus its own REQ evidence record; no unrelated implementation expansion. Decisions D-01 through D-07 explain the substantive choices.

## Acceptance testing

**Acceptance: Fail.** Existing suites pass, but the independently exercised supported legacy flow fails.

Fresh test command at the exact target, from `skills/do-work/tools/do-work-cli`:

```text
go test -count=1 ./internal/lifecycleadvance ./internal/finalization ./internal/requeststate ./internal/repositorymodel ./internal/resultmodel
```

All five passed: lifecycleadvance 3.575s, finalization 31.723s, requeststate 4.729s, repositorymodel 1.122s, resultmodel 1.340s. These include public recover → next → fresh claim, explicit takeover versus observation, hostile labels, multiple and unlabelled entries, claim-only finalization topology, all-entry removal, and checkpoint-only mutation tests. Fresh compilation of the public CLI also passed.

The six saved heavy lanes remain valid exact-revision evidence: JavaScript 8s, browser 97s, CLI integrations 60s, staged skills 23s, updater 51s, installer 23s; every lane exit 0 with no skips. The earlier timing-header refusal and canonical-helper retry are recorded transparently in the request and preparation artifact. Heavy lanes were not rerun by this review.

### Exact independent replay

Build the target CLI with `go build -o <absolute-review-cli> ./cmd/do-work-cli`. Create a temporary Git repository, configure fixture author identity, and commit exactly these two files:

`do-work/working/REQ-713-fixture.md`:

```markdown
---
id: REQ-713
title: Legacy checkpoint replay
status: claimed
claimed_at: 2026-09-04T12:00:00Z
---

# Legacy replay
```

`do-work/CHECKPOINT.md`:

```markdown
# Session Checkpoint

- REQ-713: Legacy checkpoint replay — claimed 2026-09-04T12:00:00Z — writer: other:/checkout
  Keep detail.
```

Each following invocation used `<absolute-review-cli> --repo-root <absolute-fixture-root> --format json` followed by the shown argv:

| Argv | Observed result |
|---|---|
| `recover` | Exit 0, success; one claim with one checkpoint evidence record, writer `other:/checkout`, source line 3. |
| `recover --assume-sole-authority` | Exit 1, refused; `RECOVER-CLAIM-CHECKPOINT-EVIDENCE`; claim remains working. |
| `advance --checkpoint` | Exit 0, success; `preserved_claims: 1`; appends empty `## In Progress (interrupted)` below the old entry. |
| `recover` | Exit 0, success; the same working claim now carries `checkpoint_evidence: []`. |

The fixture ran inside the detached review checkout and was removed by its temporary-directory owner. No main queue was used as a test fixture.

## Remediation boundary and ratchet

The minimum implicated implementation paths are `internal/repositorymodel/repository_model.go`, `internal/requeststate/state_apply.go`, and `internal/lifecycleadvance/checkpoint_commands.go`, all beneath `skills/do-work/tools/do-work-cli/`. Their tests belong with those owners or the public lifecycle acceptance tests. Reuse a consistent structural interpretation of supported legacy records; avoid merely silencing the refusal or reporting the old pre-write count as preservation.

Suggested named checks for RED-before-GREEN remediation: `TestRecoverLegacyCheckpointClaimsThroughPublicCommand` and `TestAdvanceCheckpointPreservesLegacyClaimDiscovery`. The first should reproduce the refusal before the fix and prove authorized recovery removes every matching legacy entry while leaving unrelated records untouched. The second should compare public discovery before/after checkpoint refresh and prove writer, request identity, and continuation evidence remain attributable. Canonical-section, labelled/unlabelled, alias, duplicate, and no-claim controls should remain green. These names describe tests to add; this review did not install tests or fixes.

REQ-544 is an existing pending sweep about position-anchored caller-authored lifecycle evidence. It owns publication substring gates and cleanup checkpoint line selection. The checkpoint writer's substring heading check is related structural territory, but F-01 also fails on an ordinary heading-less file with no forged token. Attribution and fold belong to the orchestrator after claim; this review did not change REQ-544.

## Restatement and guardrail review

Swept core and sibling Markdown for removed `Session Checkpoint Template`, old own-label/three-hour recovery wording, and `recover-finalization --discover` consumers. Current remaining finalization references describe the still-supported primitive; changelog history is context. The current saved-to-HEAD comparison confirms that later queue/gate/finalization extensions are intentional, while the F-01 functions retain the failed logic. REQ-515's set-aside preservation is a separate later policy and is not incorrectly charged to this migration.

No speculative dependency, unrelated refactor, or naming concern was found in the new functions/files. The three new mechanics use existing typed result, repository snapshot, state plan, and transaction owners. The key test weakness is semantic composition: each canonical fixture passes while the explicitly supported legacy shape crosses inconsistent readers. Self-validation checked that this was an actual public result difference, not a complaint based only on helper code or current-version drift.

## Scores

**Overall: 50%** — the acceptance-failure cap applies to the 85% arithmetic average.

| Dimension | Score |
|---|---:|
| Requirements | 80% |
| Code quality | 85% |
| Test adequacy | 75% |
| Scope | 100% |
| Risk | Critical |
| Acceptance | Fail |

**Suggested additional testing:** the two remediation checks above and their canonical controls. No browser or external-service uncertainty remains for this finding.

**Follow-ups created:** None; the orchestrator requested read-only review and owns remediation/fold routing.

**Cleanup confirmed:** all test/build/replay processes completed. The owned binary and fixture were removed; `git status --porcelain=v1 --untracked-files=all` in the detached checkout was empty. `git worktree remove .git/work-run-20260905/review-504` then completed without force. No background work remains.
