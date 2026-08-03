---
id: UR-017
title: Board drawer Copy button omits the ticket's metadata
created_at: 2026-08-03T22:49:36Z
requests: [REQ-089]
word_count: 34
---

# Board drawer Copy button omits the ticket's metadata

## Full Verbatim Input

when I hit copy the status, domain, the header itself does not copy over [Image #1] was this already captured but not yet implemented, or do we need to use do-work capture-request to capture it now?

## Assets

`assets/REQ-089-board-drawer-copy-metadata.png` — full-window screenshot of the do-work Kanban board
served at `http://127.0.0.1:8090` for the project `g1w-game-find-the-difference` (a *consumer* of this
skill, not this repo), generated 2026-08-03 22:37 UTC.

Layout: the board header carries the `do-work/queue` breadcrumb, the project name, a filter row
(id/title search, All domains, All statuses), view tabs (Board / Calendar / Testing), a LENS toggle
(Columns / By UR), and a RECENTLY DONE window selector (24h / 48h / 7d). Below it, four columns:
PENDING (1), CLAIMED (1), NEEDS INPUT · BLOCKED (4), RECENTLY DONE (21).

The right half is the open card drawer for **REQ-1322**, which is what the request is about. Its
header row shows a `REQ` kind chip, the id `REQ-1322`, and two buttons — an outlined **Copy** and a
**Close**. Under that is the H2 title `The raw "<label> HTTP <status>" string still reaches editors on
the other 14 CMS pages`, then a grey metadata card with six labelled rows:

| Label        | Value                          |
| ------------ | ------------------------------ |
| STATUS       | claimed                        |
| DOMAIN       | ui-design                      |
| USER REQUEST | UR-323                         |
| CREATED      | Aug 3, 16:14 UTC · 6h ago      |
| CLAIMED      | Aug 3, 22:07 UTC · 30m 52s     |
| TREE         | working                        |

Below the metadata card the rendered REQ body begins with an H1, then `## What`, `## Why this is your
call`, `## Open Questions`, and an `## Answers` heading at the cut-off point.

**Note on a detail visible in the screenshot:** the drawer's H2 title renders the angle-bracket tokens
correctly (`The raw "<label> HTTP <status>" ...`), but the rendered body's H1 immediately below shows
`The raw " HTTP " string still reaches editors on the other 14 CMS pages` — the `<label>` and
`<status>` tokens are gone, because the body H1's unbacktick-quoted angle brackets are parsed as HTML
tags when the Markdown body is rendered to HTML. That is a **separate** observation from this
request's subject (the clipboard payload) and was not part of what the user asked about. It is not
captured here — it needs its own REQ if the user wants it addressed.

## Capture Clarifications

Two questions were asked and answered during capture:

1. **"In the text you pasted, was there a `# REQ-1322: The raw "<label> HTTP <status>" ...` line at
   the very top?"** → **"Yes, the line was there."**

   This rules out a regression of the heading feature shipped in 0.163.3. The "header" the user
   refers to is therefore the grey metadata card (STATUS / DOMAIN / USER REQUEST / CREATED / CLAIMED
   / TREE), not the `# REQ-NNNN: <title>` heading line. The request is scoped to metadata only.

2. **"What should the Copy payload include for the metadata?"** → **"Verbatim frontmatter + body"**
   (the recommended option), with this shape shown and accepted:

   ```
   ---
   id: REQ-1322
   status: claimed
   domain: ui-design
   user_request: UR-323
   created_at: 2026-08-03T16:14:00Z
   claimed_at: 2026-08-03T22:07:00Z
   ---

   # REQ-1322: The raw "<label> HTTP <status>"...

   ## What
   ...
   ```

   The stated rationale for this option over a readable metadata list: the paste round-trips — it can
   be saved back as a valid REQ file.

## Findings Established During Capture

- **Not previously captured.** A sweep of `do-work/queue/` (REQ-083…088), `do-work/archive/`
  (including UR-001…UR-014 and `legacy/`), the two other open URs (UR-015, UR-016), `CHANGELOG.md`,
  and git log found no REQ or UR touching the Copy button or the clipboard payload.
- **Both prior Copy changes shipped without a REQ:** the button itself in 0.119.0 (commit `5e711cd`,
  2026-07-11) and the identifying heading in 0.163.3 (commit `b8b5ea4`, 2026-08-02).
- **Not a bug against spec.** `docs/board-guide.md:50` promises only "the source Markdown under a
  `# REQ-042: <title>` heading," which is exactly the current behavior. This is a feature request.

---
*Captured: 2026-08-03T22:49:36Z*
