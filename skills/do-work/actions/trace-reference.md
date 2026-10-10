# Trace Reference

> **Part of the do-work skill.** Lookup material for `actions/trace.md`: source handling, search order, the commit evidence call, the coverage table, and the go-ahead payload.

## Source Handling

- **UR id:** read that UR's `input.md` and every file in its `assets/`, from `do-work/user-requests/UR-NNN/` or `do-work/archive/UR-NNN/`.
- **Image:** read it directly and write a thorough text description (what it shows, visible text, layout, problems visible), as `actions/capture.md` Step 4 does for screenshots. Describe it before the split; the asks come from the description.
- **Pasted text:** use it as is.
- **URL:** use the platform's web fetch first. When the result is a shell (only the site name, or almost no text), say so and ask the user for a paste or a screenshot. Never judge coverage from it; a row that rests on it is `unverified`.

Write the transcription to a new `mktemp -d` directory outside the working tree, for example as `pm-ticket-source.md`, and print its path. When the source has an external ticket id, the file's first lines are:

```text
Ticket: <ticket id>
Source: <source URL, or "pasted text", "screenshot", "UR-NNN">
Fetched: <UTC date of the read>
```

The file stays outside the repository until go, when it becomes capture's verbatim input and asset.

## Search Order

Search each ask in this order and record every candidate, not only the first:

1. UR `input.md` files in `do-work/user-requests/` and `do-work/archive/UR-*/` (the same source may have been captured before, often under its ticket id).
2. REQ title, heading, `## What` and `## Detailed Requirements` in `do-work/queue/`, `do-work/working/` and `do-work/archive/`, compared by intent, not by file name. A fixed-string search (`grep -rilF`) over titles and `## What` may narrow the candidates before you read whole files. `grep` exit 1 means no match; exit 2 means the search failed, and an ask that rests on it is `unverified`.
3. Each matched REQ's `## Implementation Summary` and its `commit:` value.
4. The commits the board credits to each matched REQ (Commit Evidence below).
5. Code and tests, by the names the ask uses. A test whose assertion matches the REQ's `**GREEN when:**` line is evidence for that criterion.

## Commit Evidence

Build and run the board tool's read-only `request-commits` subcommand with every matched REQ id:

```bash
(cd <suite-root>/do-work-board/tools/queue-kanban && go build -o queue-kanban .) 2>/dev/null \
  && <suite-root>/do-work-board/tools/queue-kanban/queue-kanban request-commits --repo-root <project-root> <REQ ids>
```

The build writes only the gitignored `queue-kanban` binary in the board tool directory. Exit 0 prints a tab-separated header (`request_id`, `commit`, `committed_at`, `paths_outside_do_work`, `subject`) and one row per commit the board credits to a named REQ, newest first. The header alone means no commit credits any named REQ. A `paths_outside_do_work` of 0 is bookkeeping (claim, archive, run artifacts) or a merge; above 0 touched files outside `do-work/`: code, docs, or, where the completion commit also releases, the changelog and version files. Cite the newest row above 0 in the Commit column.

Any other exit (no `go`, a failed build, exit 1 for a git failure, exit 2 for a bad argument) means the git search was skipped. Say so in the report, use each matched REQ's `commit:` value as its commit evidence, and mark a row `unverified`, never `not started`, when its verdict rests on the skipped search.

## Coverage Table

| Ask | REQ (plain title, status) | Commit | Criteria evidenced | Verdict (as of) | Gap | Open question |
| --- | --- | --- | --- | --- | --- | --- |
| A1 | REQ-210 (retry on HTTP timeout), completed | `ab12cd3`, 2026-09-30 | 3 of 3 + GREEN yes | complete (2026-10-09) | none | none |
| A2 | REQ-214 (phone album layout), completed | `ef45ab6`, 2026-10-03 | 2 of 4 + GREEN no | partial (2026-10-09) | landscape, empty state | Is landscape in scope? |
| A3 | none found (searched: "stamp", "first day") | none | n/a | not started (2026-10-09) | all | none |

- **Criteria evidenced** counts the matched REQ's `## Detailed Requirements` bullets (`N of M`) plus its `**GREEN when:**` line (`GREEN yes` or `GREEN no`), each backed by a commit, a test, or a file and line. A REQ with neither section counts its `## What` as the one criterion. When an ask covers only part of its REQ, count only the criteria inside the ask and name them (`1 of 1, bullet 1 only`).
- **As of** is the UTC date this trace ran.
- **Verdicts:**
  - `complete`: every acceptance criterion of the known scope has evidence (commit, test, or file and line).
  - `partial`: some criteria have evidence, and the unfinished ones are named.
  - `not started`: no code evidence, and either a pending REQ or no REQ at all. The search terms are the evidence.
  - `unverified`: the evidence is not enough either way, including any row built from an incomplete fetch.

## Go-Ahead Payload

On go, invoke capture once, as `capture-request:`, for the `partial` and `not started` rows only. Its verbatim input is those rows of the table plus the transcription file, which capture attaches as an asset. Per row, carry:

- the ask id and its verdict;
- for a `partial` row: the matched REQ id and its archive path, so capture writes an addendum REQ with `addendum_to` (`actions/capture.md` Step 2), and the cited commit hash for its `## Prior Implementation`;
- when the source has one: the external ticket id, with the instruction that the UR title starts with it (capture has no ticket rule of its own), and the source URL and fetch date (already the transcription's first lines).

Capture picks the UR and REQ numbers when it writes. `complete` rows stay in the printed table only.
