---
id: UR-110
title: 'Close residual workspace release identity gaps'
created_at: 2026-09-03T23:26:06Z
requests: [REQ-565]
word_count: 77
---

# Close Residual Workspace Release Identity Gaps

## Full Verbatim Input

> ```
> REQ-512 re-review found two remaining fail-closed workspace ownership gaps after its one permitted remediation pass. Duplicate or competing Cargo and uv package name declarations can satisfy the first-match identity parser. Parseable npm root lock versions that differ from the release old version can be omitted from a changed root's required mirror count, allowing stale root copies to remain untouched. Create a critical review follow-up that requires unique source identity and every structurally present npm root version copy.
> ```

---
*Captured: 2026-09-03T23:26:06Z*
