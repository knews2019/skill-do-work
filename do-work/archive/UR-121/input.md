---
id: UR-121
title: 'Make the board top bar one line: single-line identity, Touched-in chips move into the Activity view'
created_at: 2026-09-05T12:40:00Z
requests: [REQ-586]
word_count: 44
---

## Summary

Follow-up to UR-120 (fix the double scrolling on the Activity view). Looking at the mockups, the user pointed at the board's top bar: when its control groups wrap, the identity block on the left wraps to four lines and the bar grows from 68 px to about 150 px, which costs vertical space the Activity view needs for REQ-573 (open the detail drawer from an Activity row and highlight every row of the same REQ). Three options were offered; the user chose O1 (one-line identity) and O2 (move the Touched-in window chips out of the top bar and into the Activity view's summary line) and asked for the capture before any build.

## Full Verbatim Input

> ```
> <- this part of the header is still taking up too much vertical space, and that is precious when I want to click a req and I want to highlight all of it's occurances
> 
> ok, do o1 and o2 capture it first do-work capture-request
> ```
