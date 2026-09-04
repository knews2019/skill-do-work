---
id: UR-113
title: 'Refuse a release that leaves a version or changelog mirror undeclared'
created_at: 2026-09-04T22:35:42Z
requests: [REQ-569]
word_count: 56
---

# Refuse a Release That Leaves a Version or Changelog Mirror Undeclared

## Summary

The board flagged the version action file as one release behind the changelog after the 0.282.0 release (the read-only next-lifecycle-step command). The release commit wrote only two of the four version and changelog mirrors; a repair commit fixed the rest nineteen minutes later. The user asked why a mechanical sync keeps drifting and accepted the proposed fix: the finalizer should discover every mirror that still carries the old version and refuse a release manifest that does not declare it, instead of trusting a hand-written target list.

## Assets

- `assets/REQ-569-board-version-mismatch-finding.png` — the board finding the user pasted.

## Full Verbatim Input

> ```
> [screenshot of the board's Verify Findings panel: VERSION-CHANGELOG-MISMATCH — version 0.281.0 is behind the newest CHANGELOG.md entry 0.282.0 (Expose the Next Lifecycle Step)] <- why is it so hard to keep this in sync? isn't this a mechanical sync? since it's able to check it, why not update it as well?
> 
> capture it, and build it
> ```

---
*Captured: 2026-09-04T22:35:42Z*
