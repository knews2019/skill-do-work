---
id: UR-055
title: Timestamp fabrication feedback — accepted findings
created_at: 2026-08-18T12:28:33Z
requests: [REQ-244, REQ-245, REQ-249, REQ-251, REQ-253, REQ-254, REQ-259, REQ-261, REQ-262, REQ-263, REQ-264, REQ-269, REQ-270, REQ-273, REQ-299, REQ-309, REQ-310, REQ-312, REQ-317]
word_count: 460
---

# Timestamp fabrication feedback — accepted findings

## Summary

External feedback (triaged via `do-work-toolbox validate-feedback` this session) reported an agent fabricating a `created_at` stamp instead of running a clock command. Triage verdicts: Finding 1 (audit and cite all timestamp write sites) — Accept; Finding 2 (mechanical stamping / write-time guard) — Discuss, deliberately NOT captured here, parked pending a design decision; Finding 3 (broaden the board's future-stamp warning text) — Accept. This capture covers the two accepted findings only.

## Extracted Requests

| REQ | Title | Finding |
|---|---|---|
| REQ-244 | Cite the Timestamp rule at every timestamp write site | Finding 1 |
| REQ-245 | Name fabricated stamps in the board's future-stamp warnings | Finding 3 |

## Batch Constraints

- Finding provenance travels with each REQ (verbatim claim, triage Evidence, Surface-cost result) per the validate-feedback capture handoff.
- The Timestamp rule's centralization clause (work-reference.md ~line 101: "This paragraph is the only place in `actions/` that spells a command") binds REQ-244 — citations only, never per-site command copies.
- Finding 2 is intentionally absent: the proposed `stamp` verb contradicts the frontmatter CLI's documented read-only decision (`frontmatter_cli.go` ~line 38) and CLAUDE.md § Kanban Board Write Surfaces. Parked as a Discuss item, not queued.
- Triage provenance note: the reported incident values (14:20:00Z stamps, commit at 10:46:45 UTC) do not reproduce in this checkout's git state; the root cause (uncited placeholder sites) was verified independently and stands on its own.

## Full Verbatim Input

do-work validate-feedback: Symptom: two review-generated follow-up REQs, created in the same session by the
"Review Fix" template in do-work/actions/review-work.md (Step 10), were both stamped
created_at: 2026-08-18T14:20:00Z — while the commit creating them was authored at
10:46:45 UTC (13:46 local, UTC+3). The value matches neither UTC nor local time and
is a round :20:00 shared by both files: the agent fabricated a plausible timestamp
instead of obtaining one. The board later flagged both cards with future-stamp
warnings, and its diagnosis line ("likely local wall-clock time stamped with a Z
suffix") misdiagnosed this case — the failure mode was invention, not timezone error.

Root cause: the Timestamp rule (do-work/actions/work-reference.md → Request File
Schema) binds write sites that say <timestamp> or <now>, and most stamping sites cite
it inline. But several templates use a bare [timestamp] placeholder with no citation,
so an agent filling the template from context never re-reads the rule and nothing at
the site tells it to run a clock command. Known uncited sites:
- do-work/actions/review-work.md — "Review Fix" follow-up template, created_at (the
  site that produced this bug)
- do-work/actions/work-reference.md — Builder-Decided Follow-up Template created_at
  (~line 627) and session_ended (~line 893)
- do-work-toolbox/actions/code-review.md — follow-up template created_at (~line 301)
- report-body dates: forensics.md "Scan date", roadmap.md "Scan date", clarify.md
  "When:", present-work.md "Generated:"

Requested changes:
1. AUDIT: sweep all four skills for every timestamp write site — grep templates and
   action steps for "[timestamp]", "<timestamp>", "<now>", and every *_at:/date-ish
   placeholder — and bring each under the Timestamp rule with an inline citation,
   normalizing placeholder spelling to the forms the rule recognizes. The list above
   is a starting set, not the full extent.
2. FIX RELIABLY, not just by citation — a citation still lets a model fill the value
   from imagination. Consider making stamping mechanical: (a) extend the frontmatter
   CLI with a `stamp <field>` subcommand that writes the current UTC instant itself,
   and have write sites call it instead of interpolating a value; and/or (b) add a
   write-time guard mirroring the board's check — capture/review/clarify flows reject
   any *_at value later than now + 2min skew, so fabrication is caught when the file
   is written, not when the board renders.
3. Broaden the board's future-stamp warning text: "local wall-clock time with a Z
   suffix" is one cause; a fully fabricated value is a second, now-observed one, and
   the current message sends the reader to the wrong fix.

[Follow-up instruction, same session, after the validate-feedback triage report:]

capture the requests and for #2 restate using @skills/do-work/crew-members/communication-style.md

---
*Captured: 2026-08-18T12:28:33Z*
