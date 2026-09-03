---
id: UR-109
title: 'Make heavy verification change-aware and reusable'
created_at: 2026-09-03T22:58:23Z
requests: [REQ-563, REQ-564]
word_count: 108
---

# Make Heavy Verification Change-Aware and Reusable

## Summary

Replace the all-or-nothing heavy-test trigger with lane-aware selection and safe reuse of recent per-lane evidence.

## Extracted Requests

| Request | Description |
| --- | --- |
| REQ-563 | Select only affected heavy-test lanes, explain each selection, and fall back to the complete heavy suite when coverage is uncertain. |
| REQ-564 | Reuse successful per-lane evidence for up to four hours only while its deterministic fingerprint remains valid. |

## Batch Constraints

- Both requests follow REQ-539 (Cut the contract file to the incident core and split the aggregate into fast and heavy).
- Preserve `--heavy` as a force-all override.
- Never use elapsed time alone to authorize evidence reuse.
- Record whether each lane executed or reused evidence.

## Full Verbatim Input

> ```
> Replace the all-or-nothing heavy-test trigger with change-aware lane selection and reusable per-lane evidence. Select only heavy lanes affected by the request's changed paths, explain why each lane was selected, and fall back to the complete heavy suite whenever coverage is uncertain. Cache each successful lane result for at most four hours, but reuse it only when a deterministic fingerprint of its command, test inputs, fixtures, toolchain, and required environment still matches; time alone must never authorize reuse. Allow unaffected lanes to reuse evidence while affected lanes rerun, preserve --heavy as a force-all override, and record whether each lane was executed or reused. Build this after REQ-539's aggregate split.
> ```

---
*Captured: 2026-09-03T22:58:23Z*
