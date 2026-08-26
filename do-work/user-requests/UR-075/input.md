---
id: UR-075
title: 'Flag dead cross-references, and let the filter box find tickets that cite an id'
created_at: 2026-08-26T13:24:45Z
requests: [REQ-377]
word_count: 105
---

# Flag Dead Cross-References, And Let The Filter Box Find Tickets That Cite An Id

## Summary

Two gaps in the id-autocomplete work captured as UR-074, raised while REQ-374 was claimed but before
any code had been written. Both were on UR-074's "adjacent improvements raised, deliberately not
captured" list; this UR is the user asking for them.

1. **Silent failure.** An id that resolves to no board record renders as ordinary prose, so a typo or
   a reference to never-captured work is invisible until someone hunts for it. It should render as
   broken — red, with a tooltip.
2. **Search limitation.** Titles are being added to the display but not to the index, so you still
   cannot ask the filter box for "tickets that cite REQ-1679".

## Folded Requests

- REQ-374 (title-bearing-ticket-links-and-a-drawer-glossary) — the silent-failure gap on the
  **display** surface. REQ-374 was claimed and pre-dispatch when this arrived, and its captured
  requirement said the opposite ("unresolved ids stay bare"), so the requirement was rewritten in
  place and the reversal recorded as its D-01 decision rather than shipped and then undone.
- REQ-375 (copy-carries-titles-and-a-referenced-requests-glossary) — the same gap on the **copy**
  surface. A paste has no colour, so the equivalent is a `not found in this queue` glossary line.
  Recorded as an Addendum on that still-pending REQ.

## Extracted Requests

| REQ | Title | Covers |
|---|---|---|
| REQ-377 | Index cited ticket ids and let the filter box match them | "you still cannot search for tickets that cite REQ-1679 because the filter box only matches titles, not references" |

## Batch Constraints

- **The mention pattern must not become a third copy.** `bodyMentionPattern` in
  `web/board-detail.js` and `repoFileMentionPattern` in `filementions.go` already carry a
  documented lock-step obligation. A Go-side citation extractor is a third place the same shape
  would live, and drift there silently under-indexes rather than failing loudly.
- **Never guess.** An ambiguous REQ segment resolves to nothing today and must keep doing so — in
  the index as on the display. Ambiguous is not missing, and neither is a citation.

## Full Verbatim Input

> ````text
> address the **The "Silent Failure" Gap:** <- now show the placeholder as broken (red with tooltip)
> The text explicitly notes that if a REQ ID is mentioned but doesn't actually exist in the system, it simply renders as ordinary prose. There is no warning or "broken link" indicator for dead references, meaning typos in IDs remain invisible until someone manually tries to find them.
> 
> address the **Search Limitation:**
> While titles are being added to the *display*, they aren't being added to the *index*. The user notes that you still cannot search for "tickets that cite REQ-1679" because the filter box only matches titles, not references.
> ````

---
*Captured: 2026-08-26T13:24:45Z*
