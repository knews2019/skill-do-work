---
id: REQ-110
title: Name the census's fully-read files so its completeness floor is explicit
status: pending
created_at: 2026-08-05T15:16:52Z
user_request: UR-020
domain: general
prime_files: []
tdd: false
depends_on: []
maintenance: false
---

# Name the Census's Fully-Read Files So Its Completeness Floor Is Explicit

## What

Add a paragraph to `decisions/audits/2026-08-05-shell-logic-in-prose-census.md` §5 ("What this census does not claim") that names the 14 action files read end-to-end and states that every other file's rows are grep-derived. The census's per-row `VERIFIED` guarantee is sound as worded — every cited line was read — but its *completeness* claim ("every step") rests on a grep pattern for the files that were not read in full, and it does not say so.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

Raised as finding #6 (self-raised, verdict Accept) during the `do-work validate-feedback` pass over the census request. A completeness claim that rests on a search pattern should say so, or a later reader will treat the 31 grep-scanned files as having the same coverage guarantee as the 14 fully-read ones.

## Context

The 14 files read end-to-end during the census were:

`actions/work.md`, `actions/work-reference.md`, `actions/forensics.md`, `actions/cleanup.md`, `actions/version.md`, `actions/commit.md`, `actions/inspect.md`, `actions/board.md`, `actions/tidy-repo.md`, `actions/stray-check.md`, `actions/ai-report.md`, `actions/review-work.md`, `actions/memory-value.md`, `actions/validate-feedback.md`.

The remaining 31 action files and all 18 `prompts/` files were scanned with a broad pattern (backticked shell commands, `glob`, `frontmatter`, `scan`/`parse`/`compare`/`filter`/`count`) and only the matching lines were read. A mechanic phrased without any of those tokens would have been missed.

The note belongs in §5 because that section already exists to hold the census's stated non-claims — this is a third one, not a new section. Keep it short: §5's two existing bullets are each a short paragraph, and the census is already long.

Scope boundary: this REQ adds the note only. It does **not** re-scan the 31 grep-derived files to close the floor — that is separate work, and doing it here would silently turn a three-sentence honesty fix into a full second census pass.

## Red-Green Proof

**RED prompt/case:** Open `decisions/audits/2026-08-05-shell-logic-in-prose-census.md` §5 and ask "which files were read end-to-end?" — the section names none, so the 45 action files all appear to carry the same verification depth.
**Why RED now:** §5 currently states exactly two non-claims (it does not propose the extractions; it does not judge which mechanics should stay prose). Neither mentions read depth or the grep-derived rows, and no other section does either — the header at L5–9 describes the citation method but not its completeness limit.
**GREEN when:** §5 contains a paragraph that lists all 14 fully-read files by path and states that the remaining action files and every `prompts/` file were grep-derived, so a reader can locate the completeness floor without reconstructing it.
**Validation:** User confirmed — the capture request quotes the remedy proposed in the validate-feedback report's finding #6.

## Assets

None.

---
*Source: add a completeness-floor note to the shell-logic census naming the 14 fully-read files*
