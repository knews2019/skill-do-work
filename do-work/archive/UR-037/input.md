---
id: UR-037
title: validate-feedback should flag unearned defensive surface
created_at: 2026-08-11T12:00:13Z
requests: [REQ-169]
word_count: 32
---

# validate-feedback should flag unearned defensive surface

## Full Verbatim Input

also do-work validate-feedback should also flag the following:

"""
For each incident check what earned this, and is the fix still cheaper than the surface it added?
"""

### Conversation context (same session, condensed)

Sent immediately after UR-036 (stabilization batch) was captured. The quoted question is the audit rubric from that discussion (now REQ-168's disposition rule): defensive layers must name the incident that earned them, and the fix must cost less than the surface it adds. The user wants `do-work-toolbox validate-feedback` — the read-only triage of pasted review findings — to apply the same rubric when judging findings, so review remedies that propose adding guards/fallbacks/rules without a named incident get flagged instead of accepted wholesale. Motivation ties to the stated UR-036 goal: reviews keep producing 3–5 findings; unchallenged "add more defense" remedies are one way surface (and future findings) accrete.

---
*Captured: 2026-08-11T12:00:13Z*
