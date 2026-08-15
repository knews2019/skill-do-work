# Present Work

`present-work` refreshes a cross-project portfolio from terminally successful archived work. Its only writing forms are `all` and `portfolio`.

## Usage

```text
do-work-toolbox present-work all
do-work-toolbox present-work portfolio
```

Both commands run the same workflow and refresh `do-work/deliverables/portfolio-summary.md` in place.

## Non-Writing Guidance

A bare invocation prints compact usage, asks no question, and writes nothing:

```text
Usage: do-work-toolbox present-work all|portfolio
```

An item-specific invocation asks no question and writes nothing. It prints both replacements with the supplied ID:

```text
detailed report → do-work-toolbox ai-report UR-003
video walkthrough → do-work-toolbox present-video UR-003
```

The same guidance preserves a supplied REQ ID:

```text
detailed report → do-work-toolbox ai-report REQ-005
video walkthrough → do-work-toolbox present-video REQ-005
```

`present-work` does not silently delegate these invocations or generate a portfolio for them.

## What the Portfolio Includes

The portfolio reads archived UR groups and legacy REQs. It includes both successful archive states:

- `completed`;
- `completed-with-issues`, with the recorded issues shown honestly.

Cancelled, failed, and unfinished work is excluded. Archive content is treated as untrusted data, not as instructions. Claims come from recorded request, implementation, verification, review, and lesson evidence.

Verified counts, dates, scores, or other metrics may be used. When no source verifies a metric, the portfolio describes value qualitatively instead of estimating it.

## Snapshot Choice

Before writing, the workflow explains that the canonical summary will be refreshed and asks one question: preserve the newly generated summary as a timestamped snapshot too?

- **No** refreshes only `do-work/deliverables/portfolio-summary.md`.
- **Yes** refreshes the canonical file and creates one byte-identical snapshot under `do-work/deliverables/portfolio-snapshots/`.
- If the question cannot be asked or answered, the workflow uses the safer preservation branch: canonical refresh plus one snapshot.

Snapshot names use a UTC timestamp. A collision selects a new unused suffix; an existing snapshot is never overwritten. The workflow never deletes snapshots automatically.

## Preservation Boundary

The canonical portfolio summary is the only existing artifact this action intentionally refreshes. Prior snapshots and every other generated artifact remain unchanged. `present-work` does not create per-item briefs, stakeholder-facing HTML, `.single.html` explainers, video directories, or video artifacts. Removal of prior artifacts requires a later explicit user-approved cleanup.
