# Review: REQ-698 — AI report second opinions

**Approve with follow-ups** — the confirmed process-lifetime defect is fixed and independently verified on merged main. No remaining blocking code defect was found. Actual-session and complete live-batch acceptance remain incomplete and must stay visible.

Route C | Orchestrated re-review | Cumulative range `3b593d22..31e7fcdf` | Remediation commit `2abd7939` | 2026-10-15T13:13:51Z

## What's built

The toolbox routes session, targeted, terminal-default and sequential batch requests to one Go command. It verifies native model identity twice, publishes unique model-named HTML/reply pairs, checks bundle preservation and Git status, runs the existing render judge, and adds original-page banners only after a pass. Codex and pinned Claude artifacts are present in the specified REQ-632 bundle.

The first review read the original UR-156 input as source data, full claimed REQ, plan/exploration, handback and decisions, review-work procedure, required crews and five primes, and the integrated source/tests/docs/artifacts. This re-review read the updated REQ and full remediation addendum and inspected the cumulative range. The only production delta since the first review is the runner and its tests; the intervening owner evidence commit changes do-work metadata. All P-A-U boxes are checked. D-11 documents the supervisor decision. Parent-owned D-12 defers the recurring-process lesson/prime/index update to finalization; it is not an undeclared production expansion.

## Current findings

**Important:**

- **F2 — Full requested live proof remains incomplete.** The REQ's live criterion (`do-work/working/REQ-698-ai-report-second-opinion.md:49`, and Red-Green Proof) requires real session-model and all-installed-backend output. Codex targeted and pinned Claude pass, but the active chat is Codex with its exact model unknown; pinned Claude is not actual Claude-session proof. The live batch excluded Codex truthfully, Gemini failed authentication/availability with exit 55, and AGY produced no published page. The separately authorized AGY retry verified `gemini-3.8-flash-low` and native SUCCESS but placed navigation before the required opening backlink; the launcher correctly refused it, with judge `not-run`, no page/reply/banner, and no outside paths. These are acceptance limits, not evidence that the fake-backend mode contracts fail. No new external runs occurred after remediation. — impact-user-visible → report only

**Minor:** None. **Nit:** None.

## F1 closure and first-review history

The first review on `3b593d22..d0106c31` scored 80.625%, Acceptance Partial. Its Important F1 showed a successful return followed by a background child's write: `Cmd.Wait` reaped the backend leader before `TerminateGroup`, removing the helper's group-identity proof. The parent chose one in-scope remediation cycle. F1 is now **closed**, not a current report-only defect.

The original probe is unchanged at `/private/tmp/req698-review/review_probe_test.go`, with its overlay in the same directory. It runs `/bin/sh -c '(/bin/sleep 0.3; printf late > late-write) >/dev/null 2>&1 &'`, then checks that no delayed file appears after capture returns. Original result: exit 1, `capture returned success but descendant wrote afterward`, test duration 0.66 seconds. Independent merged-main rerun:

```sh
env GOCACHE=/private/tmp/req698-review/go-cache go -C skills/do-work/tools/do-work-cli test -overlay /private/tmp/req698-review/overlay.json ./internal/toolboxcommands -run '^TestReviewSecondOpinionLeaderExitLeavesWriter$' -count=1
```

**Result: exit 0, package 1.082 seconds.** Evidence: `/private/tmp/req698-review/leader-exit-remediated.log`; original RED remains `leader-exit.log` beside it. No repository test source was added by the reviewer.

Inspected the repair at `report_second_opinion.go:565–657`: a POSIX-shell supervisor retains isolated group identity while the backend runs and after it exits. It reports the backend's exact status through a separate bounded channel and waits for cleanup. Backend executable/prompt arguments stay positional through `"$@"`, backend stdin is `/dev/null`, and status fd 3 is closed before backend execution. The caller terminates owned descendants before `Cmd.Wait`, retains cleanup/pipe failures, then reaps its supervisor and drains output. No unsafe numeric group-signal shortcut or shared ownedprocess behavior change was introduced. Unsupported platforms still refuse before process start.

Related new tests (`report_second_opinion_test.go:455` and `:481`) cover redirected and inherited-output children, bounded completion, literal arguments, distinct stdout/stderr, status 0/7 and closed channels. The strengthened cancellation test compares state at return against later state: termination may resume shell statements during teardown, but no backend write may occur after capture returns. That tests the promised boundary without hiding writes occurring before the downstream preservation checks. Builder RED demonstrated both late writers and inherited-pipe delay, and focused GREEN passed in 3.137 seconds. Builder race tests passed in 15.804 seconds; vet and Windows compilation passed.

The implementation proves in-group descendants stopped and reaps its own supervisor. Reaping children already orphaned by the backend belongs to init. Deliberate process-group escape remains the existing ownedprocess boundary, not a new guarantee claimed by this fix. No additional material remediation side effect was found.

## Requirements checklist


Sixteen review criteria consolidate the REQ and the original request without treating individual prompt rules as independent score inflation. Fifteen are fully delivered; one is partial after remediation.

| # | Requirement | Status and evidence |
|---|---|---|
| R1 | Official trigger, natural forms, router and action/guide entry | Delivered — toolbox SKILL frontmatter/router and action lines 61–75. |
| R2 | Session passes known own backend/model; unknown model stops | Delivered in action contract; actual session execution belongs to R16. |
| R3 | Targeted never substitutes; terminal chooses first installed once | Delivered — parser, mode selection, native requested/resolved checks and refusal tests. |
| R4 | Serial codex/gemini/agy batch, caller exclusion, missing/failing summary | Delivered — handler loop and fake executable batch tests; live completeness belongs to R16. |
| R5 | Full clarity-editor substance, rules 1–8, process, source limits and final reply | Delivered — shipped prompt preserves those sections and resolves CDN/offline wording to Tailwind/Mermaid only. |
| R6 | Trusted prompt substitution and install-safe discovery | Delivered — launcher environment handoff, source/install fixture, absolute roots and reserved output/prefix. |
| R7 | Physical direct-child bundle, regular index, latest/current revision | Delivered — catalog walker reuse and boundary validation. |
| R8 | Required native backend flags, actual model identity, no substitution | Delivered — native metadata adapters, generation recheck, explicit failure states and tests. |
| R9 | Launcher chooses unique page/reply names and prefixed assets | Delivered — exclusive publication, paired suffixes, reply path normalization and preservation tests. |
| R10 | Backend completion followed by trustworthy preservation/outside checks | Delivered after remediation — preservation and porcelain checks now run after owned descendant teardown; the unchanged F1 probe passes on merged main. |
| R11 | Shared judge on disposable scratch; pass/fail/skipped/error retained | Delivered — judgeRender reuse and typed findings; all verdicts covered by focused tests. |
| R12 | Only launcher adds original banners after checks/pass; per-file idempotency | Delivered — guarded replacement and real-body marker tests; original integrated diff is exactly two banner blocks. D-01 resolves fresh-run suffix versus rerun-idempotency wording. |
| R13 | Page/reply, clarity changes, corrections and open questions returned | Delivered — typed summaries, prompt final-reply contract and two saved model replies. |
| R14 | Meaningful fake-backend RED/GREEN tests and command/result wiring | Delivered — recorded absent-handler RED, semantic-backlink RED/GREEN, registration and JSON tests. F1 and inherited-pipe regressions now have observed RED/GREEN coverage. |
| R15 | Scoped shipped-file hygiene and required integration gates | Delivered at this review stage — all 18 paths declared, canonical gate and reference/contract checks passed. Release publication remains orchestrator-owned. |
| R16 | Requested live proof in the specified bundle | Partially delivered — targeted Codex and pinned Claude pass; actual session proof and successful complete batch remain incomplete (F2). |

## Acceptance testing

**Result: Partial — local and merged integration checked; prior live acceptance exercised for Codex, pinned Claude and failing batch/AGY paths; consumer deployment remains unassessed.**

- Final merged canonical `bash _dev/tests/maintainer-verify.sh`: exit 0 in 119 seconds, 427 board tests and 937 CLI tests, per-file maxima 17.86/18.49 seconds. Contract/reference checks and vet are included. Request-bound focused and green gates passed. Evidence: `/private/tmp/req698-evidence/remediated-maintainer.log` and `remediated-test-gates.json`. Broad suites were not rerun by this reviewer.
- Independently reran the unchanged F1 overlay against merged main, which passes. Reviewed cancellation, status-channel and inherited-pipe regressions and recorded builder RED/GREEN/race evidence. Existing native adapter tests and full merged gates exercise the repaired runner.
- Live artifact bytes are unchanged by remediation. Codex and pinned Claude passed their prior render checks; parent-observed interactions passed on both pages. This reviewer did not independently repeat browser interaction or a model run. Codex's disclosed verifier source-link amendment separately passed render judging; the amended bytes are not attributed entirely to the model. Original-page changes remain exactly two banners, existing assets preserved.
- Read the final authorized AGY result/invariants: the candidate places sidebar navigation before the backlink, so this retry demonstrates prompt noncompliance. Its native SUCCESS is not launcher acceptance. The original batch's unretained candidate remains indeterminate because its then-current parser had a separately proven false-rejection bug.
- Source/installed launcher topology fixtures are integration evidence, not proof that a real consumer received the release. Heavy delivery lanes are listed in the REQ, not claimed as exercised by this review.

## Restatement sweep

Mandatory sole-member wave-end sweep completed on the original integration and refreshed for the repair. No inherited review members exist. Swept second-opinion forms, original/prior-report immutability, timestamped-bundle claims, model/default/exclusion semantics, filename/reply suffixes, banner rules and typed result consumers across the router, actions, guides, README and tests. The ordinary creation/revision dispatcher and separate architecture/stakeholder writers retain their own collision rules; the second-opinion dispatcher carries its explicit exception. No stale operative restatement was found.

The remediation changes only process-lifetime implementation: all three `secondOpinionCapture` callers inherit the same fix, its explanatory comment and tests now match ownership-before-reap behavior, and the existing CLI prime's bottom-up/proved-group rule agrees. The parent will amplify the recurring trap under D-12. That deferred lesson work does not justify claiming an already-written lesson. No scope drift or undocumented material decision was found.

## Suggested additional testing

1. **Live acceptance, incomplete:** with authorized working credentials and known actual session identity, verify session invocation and remaining Gemini/AGY output, native identities, judge verdicts and banners. Preserve failures rather than substitute or claim all-backend success. This review authorizes no further external calls.
2. **Deployment, unassessed:** after the parent-owned release/update, exercise the installed command in a clean consumer bundle and verify trusted prompt discovery, published relative links and banner-only changes. Run required heavy lanes under the lifecycle owner before claiming those stages.

## Scores

| Dimension | Score | Reason |
|---|---:|---|
| Requirements | 93.75% | 15 of 16 fully delivered after F1 closure. |
| Code Quality | 96% | Narrow ownership repair, truthful status and typed results; no remaining confirmed defect. |
| Test Adequacy | 96% | Meaningful initial and remediation RED/GREEN, unchanged independent closure probe, merged gate and race evidence. |
| Scope | 100% | Two-file remediation; explicit decisions and parent-owned finalization scope. |
| Risk | Low | Live backend compliance/availability limits remain explicit; no critical defect found. |
| Acceptance | Partial | Local/integration pass; actual-session/full-batch acceptance incomplete. |

**Overall: 86.4375%** = `(93.75 + 96 + 96 + 100) / 4 - 10`. No critical-risk or acceptance-fail cap applies. Partial maps to Approve with follow-ups. Informational coding-guardrails check found no additional issue; no concern is penalized twice. Self-validation checked that the successful-exit reproduction stayed unchanged and that cancellation tests measure the post-return boundary. No remaining blocking code finding was confirmed.

**Follow-ups created:** None (1 finding report only). F1 was repaired within this request; the parent owns finalization and any later promotion of F2.

## Review

**Overall: 86.4375%** | 2026-10-15T13:13:51Z

| Dimension | Score |
|-----------|-------|
| Requirements | 93.75% |
| Code Quality | 96% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F2: Requested live acceptance remains incomplete: prior targeted Codex and pinned Claude passed, actual-session proof is absent, Gemini failed exit 55, and the authorized AGY retry correctly refused a backlink following navigation with judge not-run and no publication. No later model run established broader acceptance. — impact-user-visible → report only

**Minor findings:** None.
**Resolved findings:** F1 closed by `2abd7939`, merged at `31e7fcdf`: supervisor retains group ownership through descendant teardown; unchanged independent delayed-writer overlay passes on merged main (exit 0, 1.082 seconds), with inherited-pipe/status/cancellation regressions and final canonical gate passing.
**Acceptance:** Partial — local and merged integration pass after remediation; prior live Codex/pinned-Claude success and batch/AGY refusal evidence retained; actual-session/full-batch live success incomplete; consumer deployment unassessed.
**Restatement sweep:** redefined report publication within existing bundles, original immutability's launcher-banner exception, second-opinion modes/default/exclusion, backend/model-derived filename/reply suffixes and typed result evidence; sole-member wave-end sweep refreshed for owned-process completion semantics, no stale operative restatements found.
**Suggested testing:** 2 items.
**Follow-ups created:** None (1 finding report only).

*Reviewed by review-work action; re-reviewed after the single remediation cycle against cumulative range 3b593d22..31e7fcdf*
