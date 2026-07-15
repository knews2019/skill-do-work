---
id: REQ-016
title: "CHANGELOG: keep newest 20 entries live, archive the rest with a no-git pointer"
status: pending
created_at: 2026-07-15T17:33:04Z
user_request: UR-003
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
related: []
batch: harness-bloat-cleanup
maintenance: false
---

# CHANGELOG truncation + archive

## What
CHANGELOG.md is 26,193 words / 162 entries; runtime reads only the first ~80 lines
(actions/version.md:35). Truncate the live file to the newest 20 entries; move the
older entries verbatim to `CHANGELOG-archive.md`; add `/CHANGELOG-archive.md
export-ignore` to `.gitattributes`; put a pointer in the live header that works for
tarball installs with no `.git` (0.76.1 pattern: GitHub URL). Update the
version.md glob note if it enumerates archive filenames.

## Why
~23k words of shipped payload with no runtime reader. Precedent: 0.76.0 removed two
earlier archives; 0.76.1 restored discoverability via a GitHub commit pointer for
tarball installs. Audit §1f, DELETE bucket.

## Acceptance criteria
- [ ] Live CHANGELOG.md has exactly the newest 20 entries + header pointer.
- [ ] Archive file holds entries 21..162 verbatim; nothing lost (entry-count and
      word-count reconciliation recorded in the REQ).
- [ ] Archive is export-ignored; pointer resolves without `.git`.
- [ ] actions/version.md "last 5 releases" parse still works: first 5 `## ` blocks
      sit within the first ~80 lines of the truncated file.

## Open Questions
(none)

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:**
- [ ] **[APPLY]:**
- [ ] **[UNIFY]:**
