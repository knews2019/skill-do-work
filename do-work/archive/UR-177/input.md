---
id: UR-177
title: 'Rebuild the relation map page with a readable layout'
created_at: 2026-11-14T23:21:00Z
requests: [REQ-733]
word_count: 68
---
## Summary
The maintainer looked at the first relation map page of the Tinka proposal and asked whether it is the best way to present the content. The assistant named four problems and offered one line to rebuild the page. The maintainer asked to capture that as a request, so it runs with the rest of the queue.

## Extracted Requests

| REQ | Request |
| --- | --- |
| REQ-733 | Build an improved map page beside the first one: left to right, nine boxes first, readable line labels, light dimming. |

## Folded Requests

- REQ-727 (ship a relation-map preset for ai-report second-opinion) — the improved design now also defines what the preset asks for, and the preset waits for REQ-733

## Batch Constraints
- Nothing under `skills/` changes in REQ-733.
- REQ-727 stays blocked until the maintainer has looked at the page REQ-733 builds.

## Full Verbatim Input
> ```
> convert it to a captured request, so we can run everything using do-work run
> 
> Context the message refers to. "It" is the line the assistant offered one message earlier, after the maintainer asked "is this the best way to present it?" about the map page:
> 
> Rebuild the map page by hand, nothing shipped: left to right, nine boxes first, seeds and limits open on click, readable line labels
> ```
