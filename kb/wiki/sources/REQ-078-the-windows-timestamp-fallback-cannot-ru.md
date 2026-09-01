---
title: "Lessons from REQ-078: The Windows timestamp fallback cannot run on stock Windows in either shell it names"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-01/REQ-078-the-windows-timestamp-fallback-cannot-ru.md]
related:
  - page: concept-timestamp-and-metadata-governance
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-078: The Windows timestamp fallback cannot run on stock Windows in either shell it names

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

REQ-076 (v0.167.0) added a Windows fallback to the Timestamp rule at `actions/work-reference.md:91`.
As shipped it fails on a stock Windows box in both shells it mentions: the cmdlet flag it uses requires
PowerShell 7+, and the command is offered as the remedy *for `cmd`*, where a bare cmdlet is not a
command at all. The correct PowerShell 5.1 form exists only in the session's own notes, never in the
shipped file.

## Solution summary

Rewrote the Timestamp rule as a lead paragraph, a three-item numbered list of sources, and two labelled trailing paragraphs, instead of the single run-on paragraph that already carried four ideas. Option 3 now names `(Get-Date).ToUniversalTime().ToString("yyyy-MM-dd\THH:mm:ss\Z")` — which runs on Windows PowerShell 5.1, the version a stock box actually ships — and gives the `cmd` entry point as an explicit `powershell -NoProfile -Command` wrapper, since a bare cmdlet is not a command in `cmd`. `-AsUTC` is named only as the 7-plus shorthand it is, in the paragraph explaining why the longer form is the one written down. Both literal characters are backslash-escaped rather than left to .NET's copy-unrecognized-characters behaviour.

## What worked

- **Running three grep shapes before believing the REQ's site list**, as its own constraint demanded. Shape 1 reproduced the REQ's eleven; shape 3 found three more outside `actions/`. The inventory in a sweep REQ is a floor — this is now the second batch in a row where that held.
- **`git grep -c "date -u" HEAD -- actions/` as the before/after measure.** A per-file count is a better regression artifact than a total: it names which files changed and survives a re-run.

## What didn't work

- **The first instinct on `actions/memory-reference.md`'s bash block was to write an exception into the assertion** — a fence-aware grep, or a hardcoded allowlist. Both would have shipped a hand-maintained list into the very check whose job is to prevent hand-maintained drift. Reading the block's own preamble ("`$safe_query` and `$hit_count` are already-sanitized values") produced a fix that needed no exception at all. When a rule seems to need a carve-out, re-read what the code already claims about itself.
- **A first draft of the rule cited `CLAUDE.md`** for the compiled-tooling carve-out. `CLAUDE.md` is `export-ignore`d, so that citation dangles in every consumer install. Caught by grepping the touched files before qualification, not by the suite — the existing guard greps for other idioms.

## Worth knowing

- **The REQ's diagnosis was right and its remedy line was wrong in the same way it accused REQ-076 of being.** Requirement 1 proposes `(Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")` — unescaped — while requirement 3 asks for the escaped spelling. Two requirements of one REQ disagreeing is easy to miss when each reads fine alone; requirement 3 is the later and more specific one, so it won.
- **PowerShell version reality, in one line:** `powershell.exe` is Windows PowerShell 5.1 and is on every box; `pwsh.exe` is PowerShell 7+ and is on none by default. Any Windows prescription that needs a 7-only feature is a prescription that fails by default.
- **`-NoProfile` is not optional in a prescribed one-liner.** A user profile that prints a banner writes it to stdout, and a captured stamp becomes a banner plus a stamp. This is the same class as the repo's existing "prescribed commands must emit what the next step consumes" traps.

## Back-reference

See `do-work/archive/UR-015/REQ-078-windows-timestamp-fallback-cannot-run.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `7998740`.
