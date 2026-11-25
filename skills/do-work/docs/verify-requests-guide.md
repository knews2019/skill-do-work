# Verify Requests

Use `do-work verify-requests` to check captured REQs against the original input when you want to review the capture. You can start `do-work run` directly after capturing. Use `--against` to check unfinished work after a recorded decision changes.

## When to use

- When you want to check that the capture preserved your requirements, constraints, or intended scope
- After an ADR is superseded or you override a builder decision and need to find queued work that still relies on the old choice

## What it checks

The check compares the original requirements, interaction details, scope, shared constraints, and behavioral proof with the REQs and any folded destinations.

## Gap severity

- **Important** — firm requirements completely dropped or significantly under-captured
- **Minor** — clear details over-summarized or soft preferences missed
- **Nit** — a cosmetic discrepancy with no effect on the build
- **Ambiguous** — the original input itself is unclear (only the user can resolve)

Ambiguous gaps are surfaced as questions with concrete options. You can resolve them on the spot, defer to the builder, or leave them open.

## Output

Each gap includes its severity, the original clause, the REQ or fold destination, what changed or went missing, and a proposed correction or unresolved question. You can approve repairs after reviewing the report. The follow-up tells you what was applied and which questions remain open.

## Decision revalidation

Use `--against` with either a repo-relative superseded decision-record path or an answered `builder_decided: true` follow-up REQ. Repeat the flag to compare several reversals in one queue scan.

The scan reads complete non-terminal REQ files in `do-work/queue/` and reports **Likely affected** or **Possibly affected** only when it can quote the queued requirement and explain its conflict with the replacement. It excludes claimed and archived work, lists claimed ids as excluded, and never edits a REQ or changes `status`.

An explicit scan always shows its queue-file/word estimate and proceeds. When `do-work clarify` triggers the same scan after an override, it runs automatically through 10,000 queued words and asks before spending more. Confirming the builder's choice does not trigger a scan.

## Usage

```
do-work verify-requests          # most recent UR
do-work verify UR-003            # specific UR
do-work check REQ-018            # specific REQ (finds its UR)
do-work verify-requests --against docs/decisions/records/adr-005-pipeline-is-stateful-and-resumable.md
do-work verify-requests --against REQ-025 --against REQ-031
```

For a decision-record source, pass the **superseded** file; its single `superseded-by` relation identifies the replacement. For a REQ source, pass the answered builder-decision follow-up itself — it must retain both the old `Recommended:` choice and the different user answer.
