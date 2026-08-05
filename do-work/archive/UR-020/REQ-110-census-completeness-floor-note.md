---
id: REQ-110
title: Name the census's fully-read files so its completeness floor is explicit
status: completed
created_at: 2026-08-05T15:16:52Z
claimed_at: 2026-08-05T15:18:03Z
completed_at: 2026-08-05T15:31:24Z
route: A
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
- [x] **[PLAN]:** Add one bullet to the census's §5 naming the 14 fully-read files and stating the scan method for the rest. `prime_files` is empty and no prime covers `decisions/`. Single-file edit, no code.
- [x] **[APPLY]:** Edited `decisions/audits/2026-08-05-shell-logic-in-prose-census.md` only. Four edits: the §5 note (this REQ's scope) plus three coverage corrections from the PR review, each provenance-marked in the Implementation Summary.
- [x] **[UNIFY]:** Ran `git diff --stat` — 2 files, +94/−3, both expected (the census and this REQ's own log). Grepped the added lines for `console.log`/`debugger`/`TODO`/`FIXME`/`XXX`: none. Verified markdown table integrity on every added row by normalizing escaped pipes and counting columns — 3 for the five new §1 inventory rows, 4 for the two corrected §2 rows, both matching their table headers. No project linter covers `decisions/` (`export-ignore`d, outside the contract suite's shipped-path scope); ran `_dev/tests/contract-regressions.sh` for non-regression instead — 7 failures, identical to the `main` baseline.

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

---

## Triage

**Route: A** - Simple

**Reasoning:** Names the exact file and section to edit, and the content to add was already specified in the validate-feedback finding. No exploration or planning needed.

## Implementation Summary

**What was done:** Added the completeness-floor note to the census's §5, then corrected three coverage verdicts that an external review surfaced in the same file.

**Files changed:**
- `decisions/audits/2026-08-05-shell-logic-in-prose-census.md` (modified)

Four edits, in order:

1. **The REQ's own scope** — §5 gained a third non-claim bullet naming all 14 fully-read action files and stating that the remaining 31 action files and all 18 `prompts/` files were keyword-scanned, so a `NONE` verdict on a grep-derived file means "nothing the pattern found," not "nothing exists."
2. **§1 inventory completed** (Codex review, PR #130, P3) — the section claimed to be the complete inventory of shipped executables but omitted `tools/do-work-update.sh` and the four `hooks/*.sh` scripts, while later rows cited them. Added five rows plus a sentence noting `hooks/` and `tools/` are both shipped paths per `actions/version.md` L73.
3. **Memory-hook row corrected NONE → FULL** (Codex review, PR #130, P2) — verified all three prescribed mechanics ship in `hooks/memory-stop-capture.sh`: `sed -E` redaction before truncation (L84–95, rationale L80–83), the `grep -q "session capture …"` dedupe gate (L195), and the best-effort ledger append (L232). Row now says extraction should extend the hooks, not duplicate them.
4. **`pipeline.md` Step 4 row corrected NONE → PARTIAL** (found by this repo's grep-the-primitive rule, not by the reviewer) — `hooks/pipeline-guard.sh` L46–49 already parses `pipeline.json` and counts `pending`/`in-progress` steps.

**Scope note:** edits 2–4 are outside the REQ's stated scope (the §5 note). They were folded in rather than deferred because they are factual corrections to the *same paragraph-level artifact* this REQ edits, they arrived as review comments on the open PR that carries it, and shipping a known-wrong "complete inventory" alongside a fix for a different honesty gap in the same file would have been incoherent. Each carries its provenance inline.

## Qualification

Passed — 1 file verified present in the diff, all four edits traced to a requirement or a cited review finding, no debug artifacts. Both Codex claims were independently verified against `hooks/memory-stop-capture.sh` and `actions/version.md` before being accepted; neither was taken on the reviewer's word.

## Testing

**Tests run:** `_dev/tests/contract-regressions.sh`

**Result:** Unchanged from baseline — 7 failures in the update-script behavior probes, all pre-existing. Verified by running the same suite on `main` before this work began: identical 7 failures. No new regressions.

**Red-green validation:** RED — before this change, §5 named no files and a reader could not tell which of the 45 action files were read end-to-end. GREEN — §5 now lists all 14 by path and states the scan method for the rest. Confirmed by reading the rendered section.

**Note:** `decisions/` is `export-ignore`d and not covered by the contract suite, so no test exercises this file directly. The suite was run to prove non-regression, not coverage.

## Review

**Approve** — delivers the accepted note and corrects three understated coverage verdicts in the same artifact.

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 80% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Minor:** edits 2–4 exceeded the REQ's declared scope. Justified and provenance-marked inline, but a stricter reading would have split them into a follow-up REQ.
**Acceptance:** Pass — §5 renders with all 14 files named; the three corrected rows cite verified line numbers.
**Follow-ups created:** None.

## Lessons Learned

**What worked:** Adversarially verifying both bot findings against source before accepting them. Both held up, but the P2 finding's specific claim ("most of them" are implemented) turned out to be *stronger* than stated — all three mechanics ship, so the row moved to FULL rather than PARTIAL.

**What didn't:** The original census treated `tools/checks/` and `queue-kanban` as the coverage universe and never asked what else in the shipped tree already runs. `hooks/` was in front of me — `actions/version.md` L73 lists it as a shipped path, and I read that line while auditing `version.md` — and I still built the inventory without it. A census whose baseline is incomplete understates coverage everywhere at once, which is the one error class that actively misdirects the work it exists to inform.

**Worth knowing:** This repo's "grep the same primitive across all actions before calling it fixed" rule pays off on audit documents too, not just prescribed shell. Applying it after Codex's two findings surfaced a third of the same class (`pipeline-guard.sh`) that the reviewer missed. When a review finds an instance of a *class*, the class is the finding.

## Orientation

The shell-logic census now states its own read-depth limit and no longer understates shipped coverage — its §1 inventory covers `hooks/` and the updater alongside `tools/checks/` and `queue-kanban`. Lives in `decisions/audits/`, which never ships. No skill behavior changed.
