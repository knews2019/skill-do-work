---
id: REQ-577
title: 'Repository gate repair: remove the launcher fixture single-iteration loop'
status: claimed
commit: cd179d584c602be956fb6cc7b0bda5c0490f6c87
review_at: 2026-09-05T00:09:59Z
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-05T00:05:17Z
  basis: [trivial short-circuit]
created_at: 2026-09-04T23:55:41Z
user_request: UR-098
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-shell-commands.md]
write_set: [_dev/tests/do-work-cli-launcher-behavior.sh]
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: do-work-cli-launcher-single-iteration-loop
depends_on: []
related: [REQ-506]
integration_at: 2026-09-05T00:06:17Z
builder_handback_at: 2026-09-05T00:06:17Z
dispatch_at: 2026-09-05T00:04:33Z
claimed_at: 2026-09-04T23:58:14Z
---

# Repository gate repair: remove the launcher fixture single-iteration loop

## What

Repair the repository-gate failure recorded below so dependency-gated requests can resume.

## Instances

- [ ] repository gate: shellcheck:sc2043:do-work-cli-launcher-behavior:single-iteration-loop affecting REQ-506 (found by REQ-506 / UR-098)

## Repository Gate Repair Intake

- **Parent:** REQ-506
- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** shellcheck:sc2043:do-work-cli-launcher-behavior:single-iteration-loop
- **Repair dependency:** REQ-577
- **Diagnostic evidence:** "ShellCheck warning SC2043 in _dev/tests/do-work-cli-launcher-behavior.sh: for command_name in bash; do; This loop will only ever run once. Bad quoting or missing glob/expansion?"
- **Diagnostic evidence:** "The required second direct canonical gate run exited 1 on this warning. The launcher test bytes now belong to committed release7cceea12 and remain unchanged since that failure. First attempt had only Go per-file timing overruns; those are not the deciding fingerprint."

## Triage

**Route: A — Direct to Builder.** The captured ShellCheck failure is a three-line single-iteration loop in one named test fixture. Replace it with the equivalent direct symlink command. No production behavior or new test surface is needed.

## Plan

Skipped for RouteA. Existing ShellCheck must fail before the edit and pass afterward; the launcher behavior script must preserve its complete fake-toolchain and no-Go fixture behavior.

## Acceptance Criteria

- Remove SC2043 without a suppression or weaker lint threshold.
- Keep the same bash symlink target and destination, with shell quoting intact.
- Existing launcher behavior checks and the full canonical gate pass.

## Decisions

- **D-01 — narrow repair:** Preserve all other release0.286.0 launcher work. Edit only the test fixture's redundant loop; no launcher, module, release file or new test changes.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Replace the single-iteration loop with its one direct command, after existing ShellCheck RED.
- [x] **[APPLY]:** Integrated builder67d7e14e; replaced the redundant loop in one file.
- [x] **[UNIFY]:** One-file diff inspected; builder lint/launcher/whitespace passed. Main focused and full gate follow.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, matched shipped shell/argv fixture behavior. Partial family means no targeted form; exceeds2000-token budget. Whole touch-conditional read remains mandatory.

## Pre-Build Gate Evidence

Two direct `bash _dev/tests/maintainer-verify.sh` executions exited 1 before source edits with the captured SC2043 warning at launcher fixture line56. The second identical fingerprint authorizes this repair. RouteA has no pre-flight phase. Estimation initially refused because the skip-plan note was authored before its field; temporarily demoting that note allowed the estimator, whose observed result is recorded above.

## Implementation Summary

**What was done:** Removed the launcher fixture single-iteration loop while preserving the bash source and destination, quoting, and every assertion.

**Files changed:**
- `_dev/tests/do-work-cli-launcher-behavior.sh` (modified) — one direct symlink command replaces the redundant loop.

**Implementation range:** `89ade961fefafa4deb0d5e485736b9f387ad76f0..cd179d584c602be956fb6cc7b0bda5c0490f6c87`. Builder `67d7e14e2ec40677010c5716fca90e7161de2106`.

**D-02 — existing proof:** Existing ShellCheck reproduced SC2043 before the edit and passed afterward; no duplicate test added. All mandatory crew and complete shell prime/satellite reads confirmed in handback.

## Qualification

Passed request-bound advance qualify gate for `89ade961fefafa4deb0d5e485736b9f387ad76f0..cd179d584c602be956fb6cc7b0bda5c0490f6c87`. Exactly the declared fixture changed, one insertion/three deletions. Quoted command substitution resolves the same bash and identical destination; no suppression, assertion deletion, source/API or release edit. Existing lint RED/GREEN proves the captured finding closed.

## Testing

**Red-green validation:** Existing ShellCheck failed before implementation on exact SC2043, then passed on the final source. Independent review reproduced the baseline failure and final success. Launcher behavior script passed unchanged assertions in 1.002 seconds during review. Main request-bound focused gate: exit0, launchedtrue, timed_outfalse, green.

**Canonical repository gate:** Direct unpiped `bash _dev/tests/maintainer-verify.sh`, with `GOMAXPROCS=2`, exited0 on integrated source. All contract/lint checks passed; board382 tests, 19-second module wall, slowest file9.22 seconds; CLI728 tests, 103-second module wall, slowest file20.42 seconds. Every test file below30 seconds. Canonical advance recorded focused and green-gate states satisfied at HEAD11fde0b7fcd471c34795a88307618ef4ee3a2334. The fast tier intentionally excludes separately planned heavy lanes, so no heavy success is claimed.

**Heavy selection:** Six lanes selected by exact typed plan below, all because the changed fixture matches the shared `_dev/tests` subtree. Execution is pending at the queue-exhaustion drain.

## Review

**Overall: 100%**
**Acceptance: Pass.** Independent review approved the captured repair and complete launcher fixture. It was written while the full gate was running; the orchestrator subsequently observed the full direct gate exit0 and both canonical records satisfied as documented above. No findings, no follow-ups.

### Independent Review Record

**Approve** — the direct symlink command removes SC2043 and preserves the launcher fixture. The orchestrator's full canonical gate remains pending and is required before completion.

Route A | Merge `cd179d584c602be956fb6cc7b0bda5c0490f6c87`

#### What's built

`_dev/tests/do-work-cli-launcher-behavior.sh` (modified): one quoted symlink statement replaces the three-line, single-iteration loop. Reviewed exact range `89ade961fefafa4deb0d5e485736b9f387ad76f0..cd179d584c602be956fb6cc7b0bda5c0490f6c87`; it contains exactly this one insertion and three deletions.

#### Decisions / risks for you

No implementation decision needs user input. The old loop bound `command_name` to the literal `bash`; substituting that literal preserves source resolution, destination, and quoting. All assertions are unchanged. The now-removed loop variable has no subsequent consumer. D-01 and builder D-02 accurately explain the narrow scope and use of existing lint proof. Both the REQ and handback have all P-A-U boxes checked.

#### Findings

Important: None. Minor: None. Nit: None.

#### Requirements Checklist

- [x] Remove the captured SC2043 diagnostic without suppression or weaker lint settings.
- [x] Preserve the same bash symlink source and destination, with quoting intact.
- [x] Preserve and pass the complete existing launcher behavior fixture.
- [ ] Full canonical `bash _dev/tests/maintainer-verify.sh`: running under the orchestrator; `.git/work-run-20260905/REQ-577/final-gate.json` was absent when this review was written. This report does not attest that gate.

UR-098 supplies the parent migration context; this dependency repair moves no lifecycle mechanics and changes no production contract, so its four-part migration write-set constraint does not apply to this fixture-only repair.

#### Acceptance Testing

**Result: Pass** for the captured repair and launcher behavior. Repository-wide completion remains contingent on the separate pending gate.

Independent review checks:

| Check | Result | Duration |
|---|---|---|
| `git show 89ade961fefafa4deb0d5e485736b9f387ad76f0:_dev/tests/do-work-cli-launcher-behavior.sh`, then pass those bytes to `shellcheck -s bash -` | Expected exit 1, SC2043 at line 56 | 0.155 s lint |
| `shellcheck _dev/tests/do-work-cli-launcher-behavior.sh` | Exit 0, no diagnostics | 0.048 s |
| `bash _dev/tests/do-work-cli-launcher-behavior.sh` | Exit 0, launcher behavior tests passed | 1.002 s |
| `git diff --check 89ade961fefafa4deb0d5e485736b9f387ad76f0..cd179d584c602be956fb6cc7b0bda5c0490f6c87` | Exit 0 | 0.013 s |

Finding closure is proven by the same existing lint check failing on the pre-change bytes and passing on the merged fixture, with the exact single-iteration finding surface removed. No checks were skipped. The existing fixture still exercises argv transport, tool status propagation, old and missing Go refusals, build failure refusal, module-directory cleanliness, and invocation from a changed caller directory. No new assertion was necessary for this equivalent statement simplification.

#### Domain Review

Backend checks concerning endpoints, authentication, data exposure, and write idempotency are inapplicable to the one-line test setup change. The shell prime and complete lesson satellite were read. Quoting is preserved; no new API, dependency, branch, failure fallback, or process-management behavior is introduced. Restatement sweep: nothing redefined, sweep N/A.

#### Suggested Additional Testing

Finish and record the already-running canonical gate; no additional smoke suite is warranted for this change.

#### Scores (on the record)

**Overall: 100%**

| Dimension | Score | Notes |
|---|---|---|
| Requirements | 100% | Implementation requirements delivered; aggregate gate separately pending |
| Code Quality | 100% | Equivalent, simpler statement |
| Test Adequacy | 100% | Existing lint RED/GREEN plus complete fixture acceptance |
| Scope | 100% | Exactly the declared file and repair |
| Risk | None | No changed production behavior or test assertions |
| Acceptance | Pass | Captured failure closed; canonical completion gate not yet attested |

Overall is the arithmetic mean of the four percentage dimensions. Self-validation checked the exact merge range, original UR, REQ, builder claims and decisions, full fixture assertions, baseline failure, merged success, and the distinction between focused acceptance and pending repository-wide evidence.

#### Follow-ups created

None (0 findings report only).


## Lessons Learned

Existing warning-level lint already supplies a regression check for this simple redundant-loop repair. No new reusable lesson or satellite write is warranted.

## Orientation

The launcher test fixture now passes warning-level ShellCheck while retaining its existing behavior checks.

## Heavy Verification Plan

**Base revision:** `89ade961fefafa4deb0d5e485736b9f387ad76f0`

**Target revision:** `cd179d584c602be956fb6cc7b0bda5c0490f6c87`

- **queue-kanban-javascript** — argv JSON: `["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "queue-kanban-javascript"]`; reasons: `["_dev/tests/do-work-cli-launcher-behavior.sh matched subtree _dev/tests"]`.
- **queue-kanban-browser** — argv JSON: `["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "queue-kanban-browser"]`; reasons: `["_dev/tests/do-work-cli-launcher-behavior.sh matched subtree _dev/tests"]`.
- **do-work-cli-integrations** — argv JSON: `["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "do-work-cli-integrations"]`; reasons: `["_dev/tests/do-work-cli-launcher-behavior.sh matched subtree _dev/tests"]`.
- **staged-skills** — argv JSON: `["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "staged-skills"]`; reasons: `["_dev/tests/do-work-cli-launcher-behavior.sh matched subtree _dev/tests"]`.
- **updater** — argv JSON: `["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "updater"]`; reasons: `["_dev/tests/do-work-cli-launcher-behavior.sh matched subtree _dev/tests"]`.
- **installer** — argv JSON: `["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "installer"]`; reasons: `["_dev/tests/do-work-cli-launcher-behavior.sh matched subtree _dev/tests"]`.

The claimed request is held for these lanes under work Step7.7. No heavy_verified_at/revision has been written. After all selected lanes pass unskipped, finalization is test-only with no release bump. Preserve the full supplied implementation merge.
