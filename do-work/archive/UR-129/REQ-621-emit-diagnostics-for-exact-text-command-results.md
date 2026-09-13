---
id: REQ-621
title: 'Emit diagnostics for exact-text command results'
status: completed
route: A
review_at: 2026-09-13T13:45:42Z
kb_status: pending
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-09-13T13:44:05Z
created_at: 2026-09-07T15:44:32Z
user_request: UR-129
domain: backend
impact: impact-user-visible
effort_estimate: effort-mechanical
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-shell-commands.md"]
tdd: true
maintenance: false
related: ["REQ-615", "REQ-616", "REQ-617", "REQ-618", "REQ-619", "REQ-620"]
batch: validated-finalization-feedback
depends_on: []
write_set: ["skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go", "skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_test.go"]
claimed_at: 2026-09-13T13:44:04Z
completed_at: 2026-09-13T13:45:42Z
commit: 04d550423ab7ec89b1270869b09386c43a11ec41
---
# Emit diagnostics for exact-text command results

## What
Expose existing warning/error findings on stderr for exact-text results while preserving compatibility stdout and intentional success/fallback semantics.

## Detailed Requirements
- [ ] Project existing warning/error findings to stderr for text-mode exact-text results without inventing a new warning schema.
- [ ] Keep stdout byte-compatible: the generated directory on partial success, empty output when every image fails.
- [ ] Keep intentional success status and SVG/Mermaid fallback semantics.
- [ ] Preserve JSON results and their existing findings.
- [ ] Cover both partial and all-failed image batches.

## Finding Provenance
Original severity: P2. Verdict: Accept. Original comments: F11. Duplicate mapping: F11 is a distinct defect.

Source: `do-work/user-requests/UR-129/assets/source-review.txt`, copied byte-for-byte from `/Users/t2/.codex/attachments/87105c57-a948-4858-b928-6e1c45a94733/pasted-text.txt`. External provider not identified. Original installed-path citations map to this repository's canonical `skills/do-work/` tree.

### Original Claim F11 (Verbatim)
> ```
>   - [P2] Retain stderr diagnostics for exact-text results — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go:94-95
>     For text-mode generate-report-image-batch, ExactTextOutput contains only the generated directory, or an empty string when every image fails. The command intentionally returns success with per-image warning findings. Removing stderr rendering therefore makes partial
>     failures invisible and all-failed runs exit zero without any explanation. Keep compatibility stdout unchanged while emitting warning/error findings to stderr, preserving the output-evidence contract (.claude/skills/do-work/tools/do-work-cli/prime-do-work-cli.md#L36).
> ```

## Validation Evidence
Validated against HEAD `d4ea51d7` through source, surrounding checks, independent reviewers, and Git history. No execution probe was run during triage; no staged or unstaged fix was present.

- skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image.go:232-243 adds per-image warning findings but sets ExactTextOutput to only the directory, or an empty string when all images fail.
- skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go:1027-1028 returns ExactTextOutput before rendering findings.
- skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go:94-95 writes only that rendered output; main.go wires the runtime to stdout.
- skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_test.go:42-56 already asserts all-failed success, empty exact text, and two warning findings.

History and scope correction: The old shell implementation at a7c975c5^:skills/do-work-toolbox/scripts/generate-report-image-batch.sh prints MISSING fallback warnings to stderr. No recent removal of a Go stderr writer is established; the gap arose during shell-to-Go migration (a7c975c5). The shipped report action uses JSON and retains findings, limiting impact to text compatibility callers.

## Surface-cost
Earned — partial/all-failed image generation produces real typed warnings that text callers cannot see. A small stderr projection of existing findings restores useful fallback guidance; tests pin both diagnostic presence and compatibility output.

## Red-Green Proof
**RED prompt/case:** One or all image backends fail. Typed warning findings exist, but text-mode callers receive only a directory or empty stdout and exit zero with no failure explanation.

**Why RED now:** skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image.go:232-243 adds per-image warning findings but sets ExactTextOutput to only the directory, or an empty string when all images fail.

**GREEN when:** Focused commandruntime and image-batch regressions assert per-image stderr diagnostics alongside unchanged stdout and exit status for partial/all-failed results, with unchanged JSON findings.

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
- `_dev/primes/lessons-shell-commands.md` — 8480 indexed tokens; argv/quoting or migration parity matches, but `slugged: partial` prevents narrowing and the whole satellite exceeds budget.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read listed prime files and agent rules. Write a brief technical approach before coding.
- [x] **[APPLY]:** Implement the agreed scope and test-first regression.
- [x] **[UNIFY]:** Review every changed file and run the appropriate native checks; record what was verified.

## Full Context
See `do-work/user-requests/UR-129/input.md` for the full input and batch mapping.

## Triage

**Route: A** — reconcile already-delivered focused fix against its historical implementation and existing regression tests.

## Plan

Planning not required. Verify existing implementation commit 04d550423ab7ec89b1270869b09386c43a11ec41 and close the stale request without duplicate production edits or a duplicate release.

## Implementation Summary

**Files changed by the existing implementation:**
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go` (modified or added in 04d55042).
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_review_test.go` (modified or added in 04d55042).
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_review_test.go` (modified or added in 04d55042).

**What was done:** Exact-text output retains warning/error stderr diagnostics with compatible stdout, success status, JSON findings and image fallback behavior.

Implementation predates this run: 04d550423ab7ec89b1270869b09386c43a11ec41, already released as 0.305.39. Current run verifies and reconciles the durable request only. Its original multi-fix commit includes unrelated changes, which are not attributed to this request.

## Decisions

D-01 (DECIDE & STATE): Preserve the existing implementation and reuse its supplied commit provenance; no redundant fix or release. Historical RED/GREEN verifies the unchanged regression against the fix parent, and does not claim this run authored tests before code.

## Implementation Evidence

Required lessons consulted from current lessons index. Release lesson entries retained/read where captured; large partial CLI satellite remains budget-dropped but was read by independent verifier under touch-conditional discipline. Original builder handback unavailable, so no authorship/timing reconstructed. Current P-A-U means plan existing-fix reconciliation, apply no source edits, unify by reviewing actual diff and executing focused/historical evidence.

## Qualification

Canonical qualification satisfied for original implementation range 0c0d9933..04d55042. Only the request-relevant files/hunks are attributed here.

## Reconciliation Heavy Evaluation

No new implementation delta: canonical plan-heavy-verification for ad2c63bb41dd23647e40768ce00e92171a4361cf..ad2c63bb41dd23647e40768ce00e92171a4361cf selects no lanes. Historical range above is used for behavioral verification and qualification.

## Testing

Request-bound focused test and green-gate records satisfied. Canonical `_dev/tests/maintainer-verify.sh` was executed directly for this revision; its successful status was supplied to advance. All per-test-file budgets passed.

**Red-green validation (historical differential, not current test authorship):** TestExactTextWarningsAndErrorsRemainVisibleOnStderr; TestReportImageBatchExactTextRetainsFailureDiagnostics. Existing tests copied unchanged to04d55042 parent (0c0d9933) fail on behavior (0.00s /0.44s); current passes (0.00s /0.45s). Historical text stderr is empty. Current preserves warning/error findings for partial/all-failed image generation, compatible stdout and success, with unchanged JSON control.

No new implementation changes were made by reconciliation. The typed heavy plan for this reconciliation range contains no selected lanes. The active batch also verifies all new CLI implementation through its selected shared heavy lanes.

## Review

Overall: 97.5%
Acceptance: Pass
Independent review confirms detailed requirements and historical finding closure. Requirements100%, Code95%, Tests95%, Scope100%; Risk Low. Existing implementation and tests belong to04d550423ab7ec89b1270869b09386c43a11ec41, released0.305.39. No new source is attributed to this run. Original handback unavailable; current reconciliation P-A-U and historical evidence are explicitly distinguished from authorship. No unresolved acceptance findings or follow-up requests.

## Lessons Learned

**Worth knowing:** A queued finding can outlive an independently delivered fix. Compare the existing regression with its actual pre-fix revision before adding code; source history is stronger than stale queue wording. No new subsystem lesson beyond existing finalization/evidence contracts.

## Orientation

Text-mode image commands show failure diagnostics while preserving output compatibility and fallback behavior. This run verifies and closes the preexisting fix.

## Discovered Tasks

None.

## Timing

Observed 2026-09-13T13:44:11Z to 2026-09-13T13:45:19Z: 1m 08s total, 1m 08s attributed across 1 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 1m 08s | 1 |

Slowest command: verification-gate / maintainer-verify.sh (1 argv tokens), 1m 08s, exit 0, maintainer-verify.sh (1 argv tokens).
