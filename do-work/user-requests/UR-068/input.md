---
id: UR-068
title: User-typed text is written into REQ files without neutralization
created_at: 2026-08-23T22:35:07Z
requests: [REQ-342, REQ-343, REQ-344, REQ-345]
word_count: 447
---

# User-Typed Text Is Written Into REQ Files Without Neutralization

## Summary

Three findings from one fixture session, all on the same surface: text the user types is written
into a REQ file — into the body by clarify's answered-question format, into frontmatter values by
capture and work — with no stated rule for neutralizing it first, and the board's own `verify`
cannot see the resulting damage. Two of the three ask for a rule stated once at a named entry point
so every caller inherits it; the third asks for the detection that would have caught the damage.

## Extracted Requests

| REQ | What | Domain |
|---|---|---|
| REQ-342 | State how user text is neutralized before it is written into a REQ body, at the Canonical answered-question format's named entry point | security |
| REQ-343 | Give `queue-kanban verify` a structural-anomaly probe, and lift unrecognized-status warnings into findings | testing |
| REQ-344 | State a quoting rule for writing user text into any frontmatter value, named once where the schema is defined | security |
| REQ-345 | Found while capturing this UR, not requested: adding these three REQs to the queue fails the timeline landing probe and red-lights the canonical gate | testing |

## Batch Constraints

- REQ-342 and REQ-344 both add a rule at a named entry point rather than a local fix, so each must be
  stated once and cited by its callers — not restated per caller.
- REQ-343 must keep the parser's leniency: the point is to report damage, not to start rejecting files.
- The three are independent and can land in any order. REQ-343 is the safety net for the damage
  REQ-342 and REQ-344 prevent, so landing it first makes the other two's fixtures easier to trust.

## Full Verbatim Input

do-work capture-request: A user's typed answer is written verbatim into a REQ body as "- [x] question → answer", and arbitrary typed text can forge the body's own delimiters. Proven on a fixture: a typed answer containing a "- [ ]" line becomes a real unchecked open question, so clarify Step 5's "any remaining - [ ] wins" pins that REQ in pending-answers forever and re-presents the pasted fragment at every session; an unclosed triple-backtick fence swallows every following section into a code block on the board while the actions' prose greps still see them; and a paste landing above the opening fence loses the entire frontmatter, leaving status, title and user_request empty so the REQ silently drops out of its UR. Make the Canonical answered-question format in actions/clarify.md Step 4 state how user text is neutralized before it is written — at minimum a leading "- [ ]" or "- [x]", a bare "---" line, a "## " heading and an unbalanced code fence — and state it once at the named entry point so every caller that records an answer obeys it, not just clarify.

do-work capture-request: queue-kanban verify is blind to a structurally damaged REQ file. On a fixture where six of seven REQs carried delimiter damage — including one whose opening fence was broken so its status, title and user_request all parsed empty — verify printed "OK: no findings" and exited 0. The parser is deliberately lenient and recovers almost everything, so the damage shows up as empty fields rather than a parse error, and buildBoard's unrecognized-status warning is the only trace: verify lifts duplicate-id warnings into findings but not status warnings. Add a structural-anomaly probe that fails the mechanical check on a REQ file with no leading frontmatter fence, an empty or unrecognized status, an empty id, or a missing user_request pointer, and lift the existing unrecognized-status warnings into findings the same way appendDuplicateRequestIdFindings does. Keep the parser's leniency — the point is to report the damage, not to start rejecting files.

do-work capture-request: Frontmatter values that carry user-typed text (title, blocked_by, blocked_check, stakeholder) are written as double-quoted YAML scalars, and a typed quote or colon makes the block strictly invalid — it survives today only because parseFrontmatterFields falls back to a line-based recovery parser, which is a last resort rather than a contract. State a quoting rule for writing user text into any frontmatter value (single-quoted scalar with internal quotes doubled, or a block scalar when the text contains a newline), name it once where the schema is defined in actions/work-reference.md so capture, work and clarify all cite the same rule, and add a lock-in test that a title carrying a double quote, a colon and a hash round-trips through the board parser unchanged.

---
*Captured: 2026-08-23T22:35:07Z*
