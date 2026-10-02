---
id: UR-132
title: 'Board pages get URLs, and the Testing page shows free disk space'
created_at: 2026-10-02T18:02:37Z
requests: [REQ-626, REQ-627]
word_count: 40
---
# Board pages get URLs, and the Testing page shows free disk space

## Summary
Two requests about the Kanban board (`do-work-board board`, served at http://127.0.0.1:8090 in the user's session). Captured after REQ-625 (the low-disk-space verify probe) shipped in 0.305.60 and the user asked where the free space is shown: it is not, because a verify finding appears only below the thresholds.

## Extracted Requests
- R1 → REQ-626 (Give every board page and lens its own URL): each page (Board, Activity, Calendar, Timeline, Durations, Testing) and each Board lens opens from its own URL, and switching pages updates the address bar so the link can be copied. Filters stay out of the URL (user chose the page-and-lens scope).
- R2 → REQ-627 (Show free disk space on the Testing page): an always-visible line with the repo root's free and total space, fed by the same measurement REQ-625 added, coloured by the same 10 GiB / 3 GiB thresholds. User chose the Testing page.

## Batch Constraints
- Both are board-tool changes under `skills/do-work-board/tools/queue-kanban`; the two REQs are independent and may run in either order or together.
- "Already collected" is only half true: REQ-625 measures free space on every verify run but keeps only a finding below threshold. R2 must carry the healthy measurement into the board payload too.
- Shipped files change, so each REQ is a release entry.

## Clarifications (capture time)
- Page for the free-space line: Testing page (user picked the recommended option).
- RED/GREEN for links: confirmed as captured (user picked the recommended option; filters excluded).

## Full Verbatim Input
> ```
> there are a lot of pages in that kanban (BTW: capture-request, I need a link for each page so I can point to them by URL) in one of the pages list the free space too since it's already collected
> ```
