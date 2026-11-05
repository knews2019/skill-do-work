# Review Work

Post-build quality gate: requirements check, code review, acceptance testing. Runs automatically after each work loop item, or manually on demand.

## Three review phases

### Phase 1: Requirements Check

Extracts all requirements from the REQ and original UR, then verifies each:

- **Delivered** — implemented and visible in the diff
- **Partially delivered** — some aspects implemented, others missing
- **Not delivered** — no evidence in the diff

### Phase 2: Code Review

Evaluates the diff across four dimensions:

| Dimension | What it checks |
|-----------|---------------|
| **Code Quality** (0-100%) | Patterns, naming, readability, error handling, diff hygiene |
| **Test Adequacy** (0-100%) | Meaningful tests exist, right tests run, red-green validation |
| **Scope Discipline** (0-100%) | Stayed focused, no feature creep, files touched match declared scope |
| **Risk** (Critical/Low/None) | Security, performance, data integrity, regression risk |

### Phase 3: Acceptance Testing

Actually runs the code to verify it works end-to-end:

- Happy path first, then edge cases
- For bug fixes, confirms the original bug no longer reproduces
- Checks for regressions in adjacent features

Result: **Pass** / **Partial** / **Fail** / **Untested**

The result covers only the delivery stages the review actually checked: implementation, integration, deployment and live acceptance. A stage that applies but was not checked, such as the served site for a web change, is listed under Suggested Additional Testing as unassessed.

## Scoring

Ordinary review averages the percentage dimensions with qualitative modifiers. In a selected AI-Evaluation-Loop, scores are observations and never acceptance thresholds:

| Range | Meaning |
|-------|---------|
| 90-100% | Ship-ready |
| 75-89% | Minor issues, not worth blocking |
| 50-74% | Needs attention |
| Below 50% | Significant problems |

Critical risk caps at 60%. Acceptance fail caps at 50%.

## Guardrails

Includes anti-rationalization tables, red flags (patterns that trigger extra scrutiny like "summary lists files but diff shows no changes"), and a verification checklist to ensure all scoring dimensions are filled and every requirement is walked against the diff.

## Follow-ups

Every finding gets a recorded impact token (would anyone notice in real use? does the fix change a multi-site rule?). Only `impact-critical` findings generate follow-up REQs automatically. Every other finding stays in the report and its line ends `→ report only`; the review does not append to a sweep, `pending-answers` item, or prose backlog. To promote one later, run `do-work capture` with the complete finding line quoted as the capture source.

## Output

Structured report with scores table, requirements checklist, findings by severity, acceptance result, and suggested additional testing.

## Optional AI-Evaluation-Loop during run

Select it through `do-work run REQ-042 --review-mode sol-loop` or `do-work run UR-011 --review-mode council`; add `--jester` for optional Gemini advice. These action-owned flags are consumed before CLI `advance`. Standalone `review-work` keeps the ordinary scoring and write behavior above.

In a selected loop, reviewers use caller-owned evidence-only review: independent read-only results, no REQ or follow-up writes. The parent validates findings, writes the durable ledger and Review, and owns the decision. Sol requires an actual verified GPT Sol invocation; council requires independent OpenAI, Anthropic and Google invocations with recorded model/identity evidence. A missing or failed required reviewer never approves.

Approval needs the parent and every required reviewer on the same current candidate, including dirty and untracked implementation bytes, with no supported material finding unresolved. Later behavior changes invalidate it. The loop allows five total rounds and stops after two consecutive stalled rounds; only evidence-supported improvement resets stalls, never total rounds. A recoverable best checkpoint and counters survive recovery, and a restored candidate needs review within the remaining budget.

The optional jester defaults to the latest available Gemini resolvable through the harness. It gets the same serious brief without being told its advisory role. Its scores never affect acceptance or council grades in PR descriptions, and a slow/missing reply never delays the loop. Validated material defects still enter the ledger on their evidence.

The parent reports **AI-Evaluation-Loop: approved** or **AI-Evaluation-Loop: unresolved / not approved** with candidate, identities, round/stall counts, remaining issues and checkpoint location. Unresolved selected review uses `completed-with-issues` without calling it approval; this existing status still satisfies dependencies and UR completion. See [installed usage and availability requirements](./work-guide.md#optional-ai-evaluation-loop) and [the full contract](../actions/ai-evaluation-loop-reference.md).

## Usage

```
do-work review-work              # most recent completed REQ
do-work review REQ-005           # specific REQ
do-work review UR-003            # all completed REQs under a UR
do-work review code
```
