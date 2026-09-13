---
id: REQ-618
title: '[impact-critical] Preserve committed-risk before finalization rollback'
status: claimed
route: A
review_at: 2026-09-13T13:39:09Z
builder_handback_at: 2026-09-13T13:29:13Z
integration_at: 2026-09-13T13:29:13Z
kb_status: pending
estimate:
  p50_active_minutes: 15
  confidence: medium
  basis:
  - Route A
  - 3-file write set
  - 5 acceptance criteria
  calculated_at: 2026-09-13T13:23:14Z
created_at: 2026-09-07T15:44:32Z
user_request: UR-129
domain: backend
impact: impact-critical
effort_estimate: effort-substantive
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
related: ["REQ-615", "REQ-616", "REQ-617", "REQ-619", "REQ-620", "REQ-621"]
batch: validated-finalization-feedback
depends_on: []
write_set: ["skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply_test.go"]
required_lessons: [_dev/primes/lessons-releases.md]
claimed_at: 2026-09-13T13:22:40Z
commit: 0df803498519648f8359b232b1d9921ca8c8610a
---
# Preserve committed-risk before finalization rollback

## What
Preserve a commit that already exists and its unresolved verification risk when exact-path commit reports failure; do not restore pre-commit lifecycle state after HEAD advances.

## Detailed Requirements
- [ ] Record and durably preserve nonempty committed-risk commit evidence before handling failure.
- [ ] Return the committed-risk outcome and actionable SHA/revert evidence.
- [ ] Never perform pre-primary rollback after an actual commit exists.
- [ ] Preserve unresolved verification failure on recovery; merely advancing the phase must not silently bless the bad commit.
- [ ] Recovery must not duplicate the commit or restore the working request over committed archival.

## Finding Provenance
Original severity: P1. Verdict: Accept. Original comments: F7. Duplicate mapping: F7 is a distinct defect.

Source: `do-work/user-requests/UR-129/assets/source-review.txt`, copied byte-for-byte from `/Users/t2/.codex/attachments/87105c57-a948-4858-b928-6e1c45a94733/pasted-text.txt`. External provider not identified. Original installed-path citations map to this repository's canonical `skills/do-work/` tree.

### Original Claim F7 (Verbatim)
> ```
>    - [P1] Preserve committed-risk evidence before handling commit failure — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:134-138
>     If a hook stages an extra path, CommitExactPaths returns both a failure and the SHA of an already-created commit. This branch now returns before recording that SHA, and the deferred rollback restores lifecycle/release files and resets the journal to prepared despite
>     HEAD having advanced. A regression probe reproduced a restored working request alongside committed archival. Record the SHA before handling failure and prohibit pre-primary rollback once a commit exists, preserving the journal contract (.claude/skills/do-work/tools/
>     do-work-cli/lessons-do-work-cli.md#L158).
> 
> ```

## Validation Evidence
Validated against HEAD `d4ea51d7` through source, surrounding checks, independent reviewers, and Git history. No execution probe was run during triage; no staged or unstaged fix was present.

- skills/do-work/tools/do-work-cli/internal/gittransaction/exact_commit.go:81-97 can return a nonempty CommitSHA alongside failure when a hook stages an extra path.
- skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go:1458-1461 preserves CommitSHA and revert argv in committedRisk.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:134-138 returns on Failure before recording SHA; lines 28-33 then enter pre-primary rollback.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:214-245 restores lifecycle/release preimages and resets prepared identity; finalizationFailure at line 645 loses the nested risk outcome/revert evidence.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:530-559 cannot compensate because matchingHeadCommit rejects the outside-allowlist extra path.

History and scope correction: Failure-before-SHA ordering dates to 761d8e6ab47dd940db6d686ef641ea919ee9bb31; the rollback defer arrived in e191b266b9d501ece52d1c9888f7711aaf72f281. The original wording "now" is not supported as a recent regression.

## Surface-cost
Earned — the concrete replay is a hook staging an extra tracked path after preparation. A small extension of existing commit/rollback checks is cheaper than corrupting the lifecycle trail; the GREEN hook/recovery test keeps this condition covered.

## Red-Green Proof
**RED prompt/case:** A pre-commit hook stages one additional tracked path. Git creates the commit, exact-path verification fails, and finalization rolls back the working lifecycle while the archive remains committed.

**Why RED now:** skills/do-work/tools/do-work-cli/internal/gittransaction/exact_commit.go:81-97 can return a nonempty CommitSHA alongside failure when a hook stages an extra path.

**GREEN when:** A focused finalization hook/recovery regression retains the committed archive and lifecycle, records SHA and risk in journal/result, and proves recovery neither duplicates the commit nor silently clears verification failure.

**Validation:** User confirmed — the user explicitly requested capture of the accepted triage findings and their evidence, including the proof cases above. These are targets for the builder to reproduce, not tests already executed.

## Constraints
- Preserve original claims and severity as source data; use the validated scope/history corrections when explaining the fix.
- Limit implementation to this defect and its meaningful regression coverage. Follow the existing journal, release, and output contracts.
- Capture does not execute: this request remains pending for a separate work invocation.

## Dependencies
No prerequisite request. Related requests retain separate acceptance criteria. Shared files alone do not impose sequencing.

## Builder Guidance
Use the existing Go test harness for a genuine test-first regression. The accepted remedy defines the outcome; choose the smallest implementation that satisfies it. Do not claim a recently removed mechanism was restored without additional history evidence.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 15905 indexed tokens; owning prime and relevant release/recovery/evidence failure families match, but the index marks the satellite `slugged: partial`, so targeted loading is ineligible and the whole file exceeds the 2000-token budget.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read listed prime files and agent rules. Write a brief technical approach before coding.
- [x] **[APPLY]:** Implement the agreed scope and test-first regression.
- [x] **[UNIFY]:** Review every changed file and run the appropriate native checks; record what was verified.

## Full Context
See `do-work/user-requests/UR-129/input.md` for the full input and batch mapping.

## Triage

**Route: A** — known remaining typed committed-risk diagnostic gap, localized by read-only recovery exploration.

## Plan

Planning not required. Preserve existing SHA/journal rollback protections from 04d55042, reproduce missing committed-risk outcome/action before edits, and repair the existing finalization/recovery projection with minimal code. Read current required lessons and touch-conditional satellites.

## Decisions

D-01 (DECIDE & STATE): Extend captured write_set to finalization_commands.go and review_regressions_test.go because the existing recovery set-aside projection drops action argv and existing hook regressions miss the typed result contract. This is the same captured diagnostic outcome, not a new lifecycle/schema. Keep recovery aggregate success with actionable per-request set-aside risk.

## Implementation Summary

- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go` (modified) — classify a nonempty primary commit still in release_applied as committed_state_risk; retain the exact revert action in singular and ordered finalization records and finding, include SHA in diagnostic evidence, and require manual resolution. Existing journal/SHA persistence and pre-primary rollback protections remain intact.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modified) — preserve committed-risk action argv when recovery projects a per-request set-aside finding. Ordinary refusal clears its next action as before, and aggregate recovery can continue unrelated requests.
- `skills/do-work/tools/do-work-cli/internal/finalization/review_regressions_test.go` (modified) — use an extra tracked hook path for the captured replay and strengthen existing hook rewrite/deletion checks to require the exact typed outcome and direct/recovery actions, projection agreement, stable HEAD, unresolved journal phase, and no lifecycle rollback.

- `skills/do-work/actions/work-reference.md` (modified) — preserve returned committed-risk actions in the composed exit summary and retain ordinary refusal judgment. Added by review remediation 1d4bf8f4a220be97de60779722925b16b2a0b5c7.

Four files in the final cumulative scope; no release or lifecycle source mutation in builder commits.


## Implementation Evidence


- [x] PLAN: Read CLAUDE.md, do-work skill routing, general/coding guardrails/communication/prompt-injection crew, CLI and release primes, release lessons and entire CLI lessons satellite (touch-conditional even though omitted by claim-time budget). The named implementation.md crew file is absent; continued under existing general and coding-guardrails implementation rules. Plan: expose the missing outcome/action with existing real Git fixtures, derive diagnostics from the already-durable phase/SHA condition, and retain ordinary rollback/set-aside behavior.
- [x] APPLY: Ran genuine RED before production edit, then repaired the two existing projection functions with no new state or abstraction.
- [x] UNIFY: Reviewed every changed file, ran focused real-hook and unaffected-path controls, go vet, gofmt and git diff --check, committed exact three files.

### Observed Red-Green Evidence

Before production edits, `go -C skills/do-work/tools/do-work-cli test -count=1 -timeout 30s ./internal/finalization -run '^TestReviewPrimaryCommit'` failed behaviorally (exit 1, package 5.227s). Both top-level tests failed across extra tracked path, same-path rewrite and deletion. Actual commits existed, but the result was refused; direct next argv merely retried recovery; recovered set-aside records and findings lost the revert argv. These are the precise remaining diagnostic defects, not setup/compiler failures. Existing original safeguards from 04d55042 were already effective and were not reimplemented.

After the minimal production edits, that same focused command passed (exit 0, package 2.519s). Final fixability adjustment to manual was also included in subsequent green review/control run.

Logs retained under `/tmp/do-work-20260913/`: REQ-618-red.log and REQ-618-green.log.


## Builder Decisions

D-01 (already authorized by root): Use finalization_commands.go and review_regressions_test.go outside original captured filenames because they own recovery action projection and existing meaningful hook fixtures.

D-02 (DECIDE & STATE): Reconstruct the exact existing git revert argv from the durable primary SHA and release_applied phase. The same condition already prevents unsafe recovery and exists across direct and replay paths; a journal schema field or alternate finalizer would add unnecessary state. Recovery keeps successful aggregate draining with an actionable per-request exclusion.

D-03 (DECIDE & STATE): Strengthen the original foreign-hook test to modify a tracked file, matching the captured reproduction. The seed commit stages only that fixture file and leaves intended implementation changes pending for finalization.


Implemented in 40afe677f93c360b06301f40420457447b518fef; merged range 1856d52114b24b4cb1f741a428d297d4ff421ac6..6cb141fbd49dc04b5bc5f91b2e0e083bf4d37ed7. Prior safeguards are attributed to original04d55042.

## Initial Qualification

Canonical advance qualification satisfied for exact three-file merged range 1856d52114b24b4cb1f741a428d297d4ff421ac6..6cb141fbd49dc04b5bc5f91b2e0e083bf4d37ed7.

D-04 (DECIDE & STATE): Extend scope to skills/do-work/actions/work-reference.md to repair the independent review finding that the consumer contract says all set-aside actions are empty. Align that narrow description with retained committed-risk revert evidence; ordinary refusal judgment remains unchanged.

## Review Remediation

Independent review identified stale set-aside consumer prose (impact-user-visible). Repaired in 1d4bf8f4, merged at 0df803498519648f8359b232b1d9921ca8c8610a. Recovery-set-aside contract probe passed; no new behavior tests needed for this two-line prose correction. Final cumulative range retains original base.

## Qualification

Canonical qualification satisfied for the final four-file cumulative range 1856d521..0df80349.

## Concurrent Integration Evidence

External session commit 70c169256a827641c2199bc637222624256c0ee4 (release0.305.42) landed after the supplied merge0df80349. Preserved as independent work, including literal pathspec/stable diff-prefix corrections to shared identity code. Our supplied implementation provenance remains0df80349. The143s gate passed but overlapped that commit; a fresh stable-HEAD canonical gate and scoped compatibility review were requested before finalization.

## Testing

Builder genuine RED/GREEN hook cases: 5.227s failure before production edits, 2.519s pass after; all review regressions5.168s; explicit pre-primary/refusal/seven durable-phase controls7.279s; vet0.189s, formatting and diff checks passed. Doc remediation recovery-set-aside contract passed.

Canonical direct repository gates passed at initial merge (140s), at the prose merge with a concurrent external commit during execution (143s), and freshly on stable70c16925 (68s). The final gate explicitly reused matching CLI evidence from13:35:51Z and executed remaining checks; all native per-file budgets passed (final board slowest20.39s; inherited CLI slowest21.76s). No failures excluded. Request-bound focused and green gates both satisfied on70c16925, focused primary/identity/recovery tests8.965s.

Independent reviewer merged hook tests2.733s; external70c16925 compatibility run finalization6.025s and gittransaction1.121s, all exit0. Heavy lanes selected below await shared drain.

## Review

Overall:98.75%; Acceptance:Pass; Risk:Low. Requirements100%, Code95%, Tests100%, Scope100%. Independent reviewer confirmed five requirements, genuine RED/GREEN, direct and ordered result agreement, unchanged ordinary refusal controls, and restatement sweep. One stale action-contract finding was repaired and re-reviewed in0df80349; no open findings or follow-ups. External70c16925 compatibility also passed, without attribution to this run.

## Lessons Learned

The existing opaque-evidence-projection family applies: durable SHA/journal preservation must survive every result consumer. Tests compare the singular result, ordered records and aggregate recovery finding. The existing prime already promotes this rule; no duplicate satellite entry needed. Required release lessons read in full, entire touch-required CLI satellite read; no required lesson missing. The nonexistent implementation.md crew was skipped under the repository missing-file rule.

## Orientation

A verification failure after a primary commit now returns committed_state_risk and the exact revert SHA for manual resolution. Recovery keeps the request set aside and retains that action while allowing unrelated requests to continue. It neither rolls back committed lifecycle state nor approves the bad commit on replay.

## Discovered Tasks

None outstanding. The review's prose finding was fixed in scope.

## Heavy Verification Plan

```json
{
  "manifest_path": "_dev/tests/heavy-lanes.json",
  "base_revision": "1856d52114b24b4cb1f741a428d297d4ff421ac6",
  "target_revision": "0df803498519648f8359b232b1d9921ca8c8610a",
  "forced_all": false,
  "uncertain": false,
  "changed_paths": [
    "skills/do-work/actions/work-reference.md",
    "skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go",
    "skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go",
    "skills/do-work/tools/do-work-cli/internal/finalization/review_regressions_test.go"
  ],
  "uncovered_paths": [],
  "selected_lanes": [
    {
      "lane_id": "do-work-cli-integrations",
      "command_argv": [
        "env",
        "GIT_CONFIG_NOSYSTEM=1",
        "GIT_CONFIG_GLOBAL=/dev/null",
        "bash",
        "_dev/tests/maintainer-verify.sh",
        "--heavy-lane",
        "do-work-cli-integrations"
      ],
      "reasons": [
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go matched subtree skills/do-work/tools/do-work-cli",
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go matched subtree skills/do-work/tools/do-work-cli",
        "skills/do-work/tools/do-work-cli/internal/finalization/review_regressions_test.go matched subtree skills/do-work/tools/do-work-cli"
      ]
    },
    {
      "lane_id": "staged-skills",
      "command_argv": [
        "env",
        "GIT_CONFIG_NOSYSTEM=1",
        "GIT_CONFIG_GLOBAL=/dev/null",
        "bash",
        "_dev/tests/maintainer-verify.sh",
        "--heavy-lane",
        "staged-skills"
      ],
      "reasons": [
        "skills/do-work/actions/work-reference.md matched subtree skills",
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go matched subtree skills",
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go matched subtree skills",
        "skills/do-work/tools/do-work-cli/internal/finalization/review_regressions_test.go matched subtree skills"
      ]
    },
    {
      "lane_id": "updater",
      "command_argv": [
        "env",
        "GIT_CONFIG_NOSYSTEM=1",
        "GIT_CONFIG_GLOBAL=/dev/null",
        "bash",
        "_dev/tests/maintainer-verify.sh",
        "--heavy-lane",
        "updater"
      ],
      "reasons": [
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go matched subtree skills/do-work/tools/do-work-cli",
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go matched subtree skills/do-work/tools/do-work-cli",
        "skills/do-work/tools/do-work-cli/internal/finalization/review_regressions_test.go matched subtree skills/do-work/tools/do-work-cli"
      ]
    },
    {
      "lane_id": "installer",
      "command_argv": [
        "env",
        "GIT_CONFIG_NOSYSTEM=1",
        "GIT_CONFIG_GLOBAL=/dev/null",
        "bash",
        "_dev/tests/maintainer-verify.sh",
        "--heavy-lane",
        "installer"
      ],
      "reasons": [
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go matched subtree skills/do-work/tools/do-work-cli",
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go matched subtree skills/do-work/tools/do-work-cli",
        "skills/do-work/tools/do-work-cli/internal/finalization/review_regressions_test.go matched subtree skills/do-work/tools/do-work-cli"
      ]
    }
  ]
}
```
