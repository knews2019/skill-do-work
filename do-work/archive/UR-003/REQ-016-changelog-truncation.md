---
id: REQ-016
title: "CHANGELOG: keep newest 20 entries live, archive the rest with a no-git pointer"
status: completed
created_at: 2026-07-15T17:33:04Z
user_request: UR-003
claimed_at: 2026-07-15T17:55:00Z
completed_at: 2026-07-15T18:02:00Z
route: A
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
related: []
batch: harness-bloat-cleanup
maintenance: false
commit: 4bd7342
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
- [x] Live CHANGELOG.md has exactly the newest 20 entries + header pointer.
- [x] Archive file holds entries 21..162 verbatim; nothing lost (entry-count and
      word-count reconciliation recorded in the REQ).
- [x] Archive is export-ignored; pointer resolves without `.git`.
- [x] actions/version.md "last 5 releases" parse still works: first 5 `## ` blocks
      sit within the first ~80 lines of the truncated file.

## Open Questions
(none)

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Split at `## ` boundaries; new 0.123.1 entry authored first so "newest 20" includes it; header pointer per 0.76.1 pattern (GitHub blob URL, works with no .git).
- [x] **[APPLY]:** CHANGELOG.md → 20 entries (2,686 words); CHANGELOG-archive.md → 144 entries (23,799 words) verbatim; .gitattributes export-ignore; version.md pointer note updated; version 0.123.1.
- [x] **[UNIFY]:** Reconciliation: 163 pre-split entries + 1 new = 20 + 144 ✓; word totals 2,686 + 23,799 = 26,485 vs 26,193 pre-split + new entry+headers ✓; `git check-attr export-ignore CHANGELOG-archive.md` → set; first 80 lines contain 9 entries ≥ 5 needed by version.md; SHIPPED_PATHS in version.md unchanged (archive deliberately not a shipped path, so the update-flow diff loop never sees it).

## Triage

Route A — mechanical split with exact boundaries; no exploration needed.

## Implementation Summary

**What was done:** Live changelog truncated to newest 20 entries; 144 older entries moved verbatim to an export-ignored archive with a tarball-safe pointer.

Files changed:
- `CHANGELOG.md` (modified) — 26,193 → 2,686 words; new header pointer covers both the new archive and the pre-0.65 bf15fe2 archives.
- `CHANGELOG-archive.md` (new) — 144 entries verbatim + provenance header.
- `.gitattributes` (modified) — export-ignore for the archive.
- `actions/version.md` (modified) — pointer note updated; version 0.123.1.

## Testing

- Entry/word reconciliation (see UNIFY) — nothing lost.
- `git check-attr export-ignore CHANGELOG-archive.md` → set.
- version.md parse precondition: first ~80 lines hold 9 entries (needs 5).
- Red-green validation: omitted — non-behavioral (payload/docs); reconciliation is the regression evidence.

## Lessons Learned

**What worked:** Writing the release entry for the truncation BEFORE splitting, so "newest 20" self-includes its own provenance record.
**Worth knowing:** The archive must stay OUT of version.md's SHIPPED_PATHS — it is export-ignored, so the update-flow diff would otherwise flag it as a "local customization" on every git-clone install.
