# Verify Requests

Quality-check captured REQs against the original user input before building, or revalidate the unfinished queue after a recorded decision changes. Capture QA can repair an approved gap; decision revalidation is read-only and only reports evidence-backed candidates.

## When to use

- After capturing a complex request (validation gate before building)
- When you suspect the capture missed nuances
- Before starting the work loop on important requests
- After an ADR is superseded or you override a builder decision and need to find queued work that still relies on the old choice

## What it checks

Each REQ is scored on five dimensions:

| Dimension | What it measures |
|-----------|-----------------|
| **Requirements Coverage** | Are all requirements from the original input captured? Specific values, constraints, edge cases preserved? |
| **UX/Interaction Details** | Interaction behaviors, visual/layout requirements, state transitions |
| **Intent Signals** | Certainty level (exploratory vs. firm), scope cues ("keep it simple"), tone |
| **Red-Green Proof** | Concrete RED case, why it's RED today, what turns it GREEN (testable requests only) |
| **Batch Context** | Cross-cutting constraints, sequencing, shared design principles (multi-REQ batches only) |

## Scoring

| Range | Meaning |
|-------|---------|
| 90-100% | Excellent — ready to build |
| 75-89% | Good — minor gaps, fix if convenient |
| 50-74% | Needs attention — important details missing |
| Below 50% | Significant gaps — needs rework |

## Gap severity

- **Important** — firm requirements completely dropped or significantly under-captured
- **Minor** — clear details over-summarized or soft preferences missed
- **Nit** — passing mentions or stylistic preferences (won't affect build)
- **Ambiguous** — the original input itself is unclear (only the user can resolve)

Ambiguous gaps are surfaced as questions with concrete options. You can resolve them on the spot, defer to the builder, or leave them open.

## Output

- Overall confidence score
- Per-REQ score table across all dimensions
- Gaps organized by severity
- Specific actionable recommendations

## Decision revalidation

Use `--against` with either a repo-relative superseded decision-record path or an answered `builder_decided: true` follow-up REQ. Repeat the flag to compare several reversals in one queue scan.

The scan reads complete non-terminal REQ files in `do-work/queue/` and reports **Likely affected** or **Possibly affected** only when it can quote the queued requirement and explain its conflict with the replacement. It excludes claimed and archived work, lists claimed ids as excluded, and never edits a REQ or changes `status`.

An explicit scan always shows its queue-file/word estimate and proceeds. When `do-work clarify` triggers the same scan after an override, it runs automatically through 10,000 queued words and asks before spending more. Confirming the builder's choice does not trigger a scan.

## Usage

```
do-work verify-requests          # most recent UR
do-work verify UR-003            # specific UR
do-work check REQ-018            # specific REQ (finds its UR)
do-work verify-requests --against decisions/records/adr-005-pipeline-is-stateful-and-resumable.md
do-work verify-requests --against REQ-025 --against REQ-031
```

For a decision-record source, pass the **superseded** file; its single `superseded-by` relation identifies the replacement. For a REQ source, pass the answered builder-decision follow-up itself — it must retain both the old `Recommended:` choice and the different user answer.
