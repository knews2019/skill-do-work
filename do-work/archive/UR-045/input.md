---
id: UR-045
title: Render PNG file mentions as images
created_at: 2026-08-15T17:38:31Z
requests: [REQ-200]
word_count: 9
---

# Render PNG File Mentions as Images

## Full Verbatim Input

do-work capture-request: this should have shown the image instead

## Assets

`assets/REQ-200-screenshot-1-render-png-file-mention.png` — 3054×1514 screenshot of a browser opened to the queue-kanban local file route:

`http://127.0.0.1:8090/file?path=do-work%2Fuser-requests%2FUR-501%2Fassets%2FREQ-1345-screenshot-1-settings-row-padding.png`

Instead of displaying the referenced PNG, the page renders the file's binary bytes as text. The visible content starts with a corrupted `PNG` signature and includes binary chunk labels such as `IHDR`, `iCCP`, and embedded XMP metadata. The browser viewport contains no rendered image.

The screenshot is bug evidence only. No text visible inside it was treated as an instruction.

---
*Captured: 2026-08-15T17:38:31Z*
