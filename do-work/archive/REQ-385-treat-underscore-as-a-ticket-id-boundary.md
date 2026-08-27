---
id: REQ-385
title: '[impact-user-visible] Treat an underscore as a ticket-id boundary on both surfaces'
status: completed
created_at: 2026-08-26T23:05:00Z
status_changed_at: '2026-08-27T11:57:50Z'
user_request: UR-075
addendum_to: REQ-383
review_generated: true
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-378, REQ-379, REQ-383]
batch: ticket-id-autocomplete
write_set:
- skills/do-work-board/tools/queue-kanban/citations.go
- skills/do-work-board/tools/queue-kanban/citations_test.go
- skills/do-work-board/tools/queue-kanban/web/board-detail.js
claimed_at: '2026-08-27T11:57:50Z'
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: '2026-08-27T11:57:50Z'
completed_at: '2026-08-27T12:11:02Z'
commit: 259b1479
kb_status: pending
---

# Treat An Underscore As A Ticket-Id Boundary On Both Surfaces

## What

`\b` counts `_` as a word character, so the mention pattern behaves differently around an underscore
than a reader expects. Change the boundary on both the Go and the client side together, in one
commit, so the agreement test stays green.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Extend the existing agreement corpus and real mention regression first, prove RED, then use the same post-match non-alphanumeric boundary rule in Go and JavaScript. Consume compound candidates before checking boundaries so they never fall back to a UR prefix. Keep other follow-ups unchanged.
- [x] **[APPLY]:** Added matching post-match ASCII alphanumeric boundary checks in Go/JavaScript and regression coverage for underscore emphasis and compound consumption.
- [x] **[UNIFY]:** Reviewed all three source files listed below: candidate matching/post-filter, actual drawer consumption, Go offsets and clipboard behavior. gofmt, go vet, node syntax and diff check passed. Full gate and review recorded below.

## Why

Two symptoms, one cause. Both were found by adversarial review of REQ-383 and reproduced.

**A ticket id in underscore emphasis is silently dropped.** `Fixed in _REQ-1679_ last week.` yields no
mention at all — no title, no appendix line — because `_REQ-1679_` fails the pattern's `\b` anchors in
the SOURCE bytes. The drawer scans the RENDERED text, where emphasis is already consumed, so it
resolves the id, expands the title and adds a glossary entry. `*REQ-1679*` works and `_REQ-1679_` does
not, for the same rendered output. That is a drawer/clipboard divergence, and REQ-383's stated rule is
that the two must say the same thing about the same body.

**A compound id followed by an underscore corrupts on paste.** RE2 explores all alternatives, so when
the compound alternative's trailing `\b` fails against a following word character, the shorter
`UR-\d+` alternative succeeds at the same start. `_tracked under UR-003-REQ-077_` therefore emits a
mention for the six characters `UR-003` with `Expand: true`, and the client splices the UR's title
into the middle of the compound id:

```
_tracked under UR-003 (Ship The Widget)-REQ-077_
```

The glossary then lists `UR-003` instead of `UR-003-REQ-077`.

## Context

**Not a regression.** The pattern is byte-identical to the one REQ-379 shipped and the one
`web/board-detail.js` still carries; REQ-383 moved block classification into Go and left match
semantics alone on purpose. This REQ is the follow-up that changes them.

**The corruption is latent; the dropped mention is not.** Zero of this board's 485 `id:` fields are
compound — every id in the tree is the flat `REQ-nnn` form, and the flat form degrades to NO match
rather than a corrupt splice, which is the benign direction. Underscore emphasis around an id needs no
special setup at all.

**Why both files move together.** `bodyTicketMentionPattern` (`citations.go`) and `bodyMentionPattern`
(`web/board-detail.js`) are pinned to each other by
`TestJavaScriptBehaviorTicketMentionPatternAndResolverAgreeWithGo`, which drives both over one corpus
in both directions. Changing one alone fails that test — correctly. This is why REQ-383 could not do
it: the drawer was explicitly out of its scope.

## Detailed Requirements

- **A ticket id's boundary is a non-alphanumeric character, not a non-word character.** `_` ends an
  id the way a space or a bracket does. RE2 has no lookaround, so the boundary cannot simply be
  rewritten as a lookahead — either capture and restore the boundary characters, or post-filter each
  match against the bytes on either side. The client must use whichever shape the Go side uses.
- **The compound alternative must win wherever it matches at all.** A trailing word character must not
  demote `UR-003-REQ-077` to `UR-003`; if the compound form cannot match cleanly, the correct answer
  is no match, never a shorter one at the same start.
- **Both surfaces change in one commit**, keeping the agreement test green throughout.

## Constraints

- **The agreement test is the gate, not a formality.** Its corpus must gain the underscore cases in
  the same commit; a fix that passes the current corpus has not been tested.
- **Do not widen into the other REQ-383 follow-ups.** REQ-386 (the restating H1), REQ-387 (unescaped
  titles) and REQ-388 (the remaining divergences) are separate.

## Dependencies

None — this is the head of the chain. It shares `citations.go` and `citations_test.go` with REQ-381,
REQ-386 and REQ-388, and `web/board-detail.js` with REQ-382 and REQ-388; every one of those waits on
this REQ, directly or transitively.

**This edge serializes a shared file; it is not a need for the other's output.** `write_set` gates nothing — root `CLAUDE.md` § Glossary calls it "never a safety guarantee" — so only `depends_on` keeps two writers of one file apart under `do-work run --fan-out`. The whole batch is one chain, **REQ-385 → REQ-381 → REQ-386 → REQ-388 → REQ-382 → REQ-387**, because `citations.go` alone is claimed by four of the six. That is ONE valid total order: reordering the queue means recomputing every edge, since a chain is only correct as a whole. `queue-kanban verify`'s `ungated-write-set-overlap` probe reports any pair this misses.

## Red-Green Proof

**RED prompt/case:** Copy a ticket whose body contains `Fixed in _REQ-1679_ last week.` — the paste
carries no title and no appendix line, while the drawer beside it shows both. Second case: a body
containing `_tracked under UR-003-REQ-077_` on a board holding both `UR-003` and `UR-003-REQ-077` —
the paste reads `UR-003 (…)-REQ-077`.

**Why RED now:** `\b` treats `_` as a word character, so the first case fails the anchor and the
second falls through to the shorter alternative.

**GREEN when:** both cases behave as the drawer does, the agreement corpus carries them, and no
existing assertion needed rewriting.

**Validation:** Reproduced by adversarial review of REQ-383, both cases run against a scratch copy of
the package.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: REQ-383's independent review, findings S1 and S2.*

## Triage

**Route: A** — Known shared-pattern bug with explicit reproduction; fix both surfaces in one verified increment.

## Plan

Planning not required — focused implementation guided by the request and existing patterns.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/citations.go` (modified) — consumes complete candidates before checking non-alphanumeric boundaries.
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modified) — matches the Go boundary rule and keeps all ticket candidates consumed, including intentionally suppressed compounds.
- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modified) — extends the shared agreement corpus and tests Go mention offsets, actual drawer fragment consumption, copied text, and reference lists.

**What was done:** Underscores now delimit ticket IDs. Boundary checks run after candidate matching, preventing the regular expression from falling back to a shorter UR alternative. Ticket candidates remain consumed even when resolution intentionally suppresses them; only non-ticket runs keep the previous retry behavior.

## Decisions

- D-01 (decide and state): Use ASCII alphanumeric post-filters in both languages. This preserves the previous ASCII matching model while treating underscores as punctuation, and avoids unavailable RE2 lookaround.
- D-02 (decide and state): Keep the existing shared test as a candidate/resolver agreement check and add a caller-level comparison of accepted mentions. That exercises the real production filters instead of reimplementing them in test drivers.

## Testing

**Test-first RED:** TestJavaScriptBehaviorTicketMentionUnderscoreBoundaries failed before production edits on missing underscore mentions, corrupt compound splices, and incorrect reference lists. **GREEN:** The final focused Go/JavaScript run exited0. Cases include an emoji prefix (UTF16 offsets), adjacent underscore-separated IDs, optional suffixes, invalid compound suffixes, and a valid mention after a rejected compound.

Existing assertions were preserved. The shared resolver fixture adds UR-003 so the corrupt UR-prefix fallback resolves visibly instead of remaining a hidden missing ID. The raw Markdown payload/offset tests still pass. go vet, gofmt cleanliness, node --check and git diff --check passed.

## Orientation

Underscore emphasis now receives the same ticket annotation as other prose. Invalid compound IDs stay unchanged instead of receiving a title in the middle.

## Review Remediation

Independent review found an underscore-triggered regression in code: an unresolved compound could be intentionally suppressed and then rescanned as a resolvable inner REQ. This was not accepted as completed. The narrow repair consumes every ticket candidate even when resolution suppresses it; file-path retry remains unchanged for REQ-388. A new code-context regression must fail before that repair and pass after; final verification and re-review follow the correction.

## Qualification

qualify.sh exit0. Parent inspected all three actual source diffs, including the review correction, and traced boundary checks and consumption to the captured failures. No third scanner, new dependency, or unrelated path/title policy was introduced.

## Browser Acceptance

The existing dump-dom drawer probe timed out after75s in Chrome141 before any assertion result; no pass is claimed for it or the full browser suite. A focused acceptance test in an isolated module copy used the existing trusted CDP harness instead. Parent verified all three changed files byte-identical to the main module, inspected the probe, then independently reran it on Chrome for Testing141.0.7390.37: exit0.

The actual generated board drawer displayed the titles for REQ-1679 and UR-003-REQ-077, listed exactly those two glossary entries, left invalid compound suffixes unchanged, and suppressed the unresolved compound in a fenced block. Clicking the real Copy button produced the matching annotated body and appendix without a spurious UR-003 entry or broken compound. The result recorded page URL and browser in the same measurement, with no captured page errors.

## Lessons Learned

Compound-first alternation does not guarantee compound-first behavior when a failing boundary permits fallback. Consume before checking boundaries, and keep intentionally suppressed ticket candidates consumed too: a caller's retry can recreate a fallback the regular expression no longer performs.

## Review

Initial review: R1 identified the newly exposed suppressed-compound retry, so completion was held. The added code-context regression failed before the repair (exit1) and passed after it (focused run exit0). Final independent re-review: **100% | Acceptance: Pass**, no open findings. Reviewer reran focused tests and additional production-function probes for suppressed compounds, recovery to later valid IDs, prose missing references, and unchanged path retry.

The first canonical run passed but preceded the review repair; it is not the final-state gate. A fresh canonical run on the corrected source is required and recorded below only after completion.

## Repository Verification

`bash _dev/tests/maintainer-verify.sh` completed with exit 0 on the final implementation/release state. Contract suite, Go checks, and strict JavaScript lane passed. The default optional browser lane was explicitly skipped; separate browser evidence is recorded above where applicable.
