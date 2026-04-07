# Code Review: Last 20 Commits

**Scope**: `64a9efe..9714c2b` (20 commits)  
**Date**: 2026-04-07  

---

## Bugs

### 1. BKB master index line limit contradiction (build-knowledge-base.md)

**Severity: Medium**

Line 85 says `_master_index.md` is "~50 lines max", but lines ~210, ~586, and ~1018 all say "under 80 lines." The 50-line limit in the directory tree comment will cause agents to prematurely restructure the index.

**Fix**: Change line 85 from `~50 lines max` to `~80 lines max` to match the authoritative limit used everywhere else.

### 2. BKB topic index limit vs split threshold are incompatible (build-knowledge-base.md)

**Severity: Medium**

Line ~587: "Each topic index must stay under 60 lines"  
Line ~588: "When a topic index exceeds 80 articles, split it"

If each article is one line (plus headers), 60 lines is hit well before 80 articles. Agents will hit the line limit and not know whether to split or just trim. The architect agent section repeats this same contradiction.

**Fix**: Reconcile the two thresholds — either raise the line limit or lower the article split threshold so they're consistent.

### 3. Cleanup Pass 2 misleading fallback comment (cleanup.md:57)

**Severity: Low-Medium**

Line 57 says: "leave the REQ in archive root for now, Pass 1 will handle it on next run." This is inaccurate — Pass 1 archives entire UR folders when all REQs are complete; it does NOT move individual loose REQs into UR subfolders. The REQ sits in archive root indefinitely until ALL sibling REQs complete and the UR folder is archived as a whole.

**Fix**: Change the comment to: "leave the REQ in archive root for now — it will be consolidated when the UR is fully complete and archived by Pass 1."

### 4. BKB `status` command missing defrag staleness warning (build-knowledge-base.md)

**Severity: Low**

Line ~787 says "The `status` command should note when defrag hasn't run in 14+ days," but the `status` sub-command section (lines ~726-746) doesn't include this check. An agent implementing `status` from that section would omit the warning.

**Fix**: Add a staleness check step to the `status` sub-command section.

---

## Documentation / Consistency Issues

### 5. `build-knowledge-base.md` missing from CLAUDE.md project structure

**Severity: Medium** (affects discoverability)

`CLAUDE.md` lists every action file in the project structure block except `build-knowledge-base.md`. This file was added in commit `485ccc9` but the project structure was never updated.

**Fix**: Add `build-knowledge-base.md` to the `actions/` listing in `CLAUDE.md`.

### 6. BKB help menu in SKILL.md omits 4 sub-commands

**Severity: Low-Medium**

SKILL.md lines 272-280 list 8 BKB sub-commands (`init`, `triage`, `ingest`, `query`, `lint`, `resolve`, `close`, `status`), but the action file defines 12. Missing from help: `defrag`, `garden`, `rollup`, `crew`.

**Fix**: Add the missing sub-commands to the SKILL.md help menu.

### 7. SKILL.md routing table: "code review" ambiguity between priorities 6 and 8

**Severity: Low**

Priority 6 (line 61) lists `do work code review src/` for code-review, but priority 8 (line 63) lists `do work code review` (no scope) for review-work. The routing table at priority 6 also shows `do work code-review` (hyphenated) which is fine, but the proximity of "code review" in both rows could mislead implementors. The explanatory text later clarifies, but the table itself is the primary reference.

**Fix**: Add "(with scope)" annotation to priority 6 examples, or remove the ambiguous examples.

### 8. BKB `CLAUDE.md` reference is ambiguous (build-knowledge-base.md ~line 210)

**Severity: Low**

The architect agent says "CLAUDE.md is the single source of truth for conventions." This refers to the KB's own schema file (`<kb-path>/CLAUDE.md`), but could easily be confused with the project's root `CLAUDE.md`. 

**Fix**: Use `<kb>/CLAUDE.md` or "the KB schema file" instead of bare `CLAUDE.md`.

### 9. BKB defrag vs index rules: split threshold mismatch

**Severity: Low**

Defrag (line ~759) flags "overcrowded clusters" at 40+ articles. Index Size Rules (line ~588) set the split threshold at 80 articles. The 2x gap means defrag will recommend splits that the rules say aren't needed yet.

**Fix**: Align the thresholds, or clarify that defrag is a "soft recommendation" while 80 is the hard limit.

---

## Summary

| Category | Count |
|----------|-------|
| Bugs (logic errors, contradictions) | 4 |
| Documentation / consistency gaps | 5 |
| Portability violations | 0 |
| Security issues | 0 |

The codebase is well-structured overall. The BKB action file (`build-knowledge-base.md`) introduced in this window carries most of the issues — it's a large file (~1145 lines) with several internal contradictions around numeric thresholds. The cleanup pipeline has one misleading comment. No portability or security issues found.
