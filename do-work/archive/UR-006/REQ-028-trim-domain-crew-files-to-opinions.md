---
id: REQ-028
title: Trim domain crew files to opinions
status: completed
created_at: 2026-07-27T07:34:50Z
claimed_at: 2026-07-27T07:40:27Z
completed_at: 2026-07-27T08:06:57Z
commit: f95b1cf
user_request: UR-006
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
related: [REQ-025, REQ-026, REQ-027, REQ-029, REQ-030, REQ-031]
batch: context-engineering-alignment
---

# Trim Domain Crew Files to Opinions

## What

The six domain crew files — `crew-members/{backend,frontend,security,testing,ui-design,debugging}.md`, ~6,737 words total — are largely generic engineering advice ("use appropriate HTTP status codes", "don't chase a coverage number") that any capable model already follows, with real content overlap between `frontend.md`↔`security.md` and `backend.md`↔`security.md`. Trim to ~3,500 words, keeping only opinions:

- **Resolve the duplication first.** `crew-members/security.md` owns all security content. `frontend.md` and `backend.md` drop their overlapping security bullets (`dangerouslySetInnerHTML`, token storage, API keys in bundles, resource-level authorization, rate limiting) and point at `crew-members/security.md`.
- **Per-file keep test:** a line survives only if it is (a) a deliberate preference of this project, or (b) tied to do-work machinery (REQ markers, `## Red-Green Proof`, prime-file test mappings). Lines that restate standard practice go.
- **Preserve every `JIT_CONTEXT` comment and its trigger wording verbatim** — those are load contracts, not prose.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read general.md, coding-guardrails.md, maintenance.md, all six target files, and every actions/ + crew-members/ pointer into them (`grep -rn 'crew-members/\(backend\|frontend\|security\|testing\|ui-design\|debugging\)\.md'`) to find hard section-level dependents before cutting anything. Found: `actions/ui-review.md` quotes ui-design.md's 6 phases, Heuristic Review Criteria, Severity definitions, Accessibility Baseline, and Implementation Patterns by name/concept; `actions/code-review.md` cites security.md's "OWASP Top 10 checklist and framework-specific patterns" headers. No other file has section-level (vs. file-level) pointers into the six. Plan: apply the per-line keep test (preference / machinery, else cut), move the two frontend security bullets not already duplicated in security.md (Vue/Angular innerHTML naming, `rel=noopener`) into security.md before deleting them from frontend.md, trim security.md's own generic prose (Static Analysis Catches/Misses, redundant Anti-Patterns, vague A04 row) while preserving every hard-cited header, and trim ui-design.md's prose while preserving every phase/section ui-review.md depends on.
- [x] **[APPLY]:** Rewrote backend.md, frontend.md, testing.md, debugging.md; surgically edited security.md (dedup absorption + trims) and ui-design.md (prose trim, structure preserved) — the six named files only.
- [x] **[UNIFY]:** `git diff --stat` reviewed (six files, all in scope). `bash _dev/tests/contract-regressions.sh` passed clean. Verified: JIT_CONTEXT lines byte-identical before/after (diffed each of the six); dedup grep shows each named topic in exactly one substantive location (security.md) with pointer-only mentions in frontend.md/backend.md; pointer-resolution grep confirms every `crew-members/{these-six}.md` citation in actions/ and crew-members/ still resolves to a section that exists under the same name/concept. No debug artifacts — all diffs are prose.

## Why (if provided)

The first of Anthropic's five shifts for Claude 5 generation models is rules → judgment: a capable model given a clear goal outperforms one given a long list of rules it already knows, and the rules cost context on every load. These files load JIT per domain, so the cost lands on exactly the REQs doing the most work.

The user chose this scope over the safer "keep as-is" option with the tradeoff stated explicitly: trimming to opinions removes a safety net if this skill is ever run on a weaker model. That tradeoff was accepted at plan approval — do not re-litigate it during implementation, and do not hedge by keeping generic lines "just in case."

## Context

- These files load conditionally per their `JIT_CONTEXT` (domain match, `tdd`/`caveman` flags, security surface, debugging retries). The comment in each file is the canonical statement of when it loads.
- `crew-members/general.md` and `crew-members/coding-guardrails.md` (the always-loaded pair, 1,524 w) are **out of scope** except for removing anything this trim makes redundant. `coding-guardrails.md` holds the canonical decide-vs-escalate gate and the canonical YAGNI statement that other files point at — breaking either breaks those pointers.
- Maintenance pass on the skill's own instructions — `crew-members/maintenance.md` (delete-before-you-add) loads via `maintenance: true`.

## Detailed Requirements

- **Word target:** ~3,500 total across the six files, from 6,737. Report per-file before/after `wc -w`. The target is a direction, not a quota — do not pad a file back up to hit it, and do not delete a genuine opinion to get under it.
- **The keep test is per line, and the answer must be defensible.** For each file, the Implementation Summary states what was kept and the one-word reason (preference / machinery). If the reason for a line is "it's true", that line goes.
- **Deduplication is a move, not a delete.** Any security content unique to `frontend.md`/`backend.md` moves into `crew-members/security.md` before those bullets are removed. Nothing that was said once should end up said zero times.
- **Pointer check after the move:** `grep -rn 'crew-members/' crew-members/ actions/` — every pointer must resolve to a section that still exists.
- **Do not touch `JIT_CONTEXT` comments.** Verify with `grep -n 'JIT_CONTEXT' crew-members/*.md` before and after; the lines must be byte-identical.
- **No new files.** This is a trim, not a reorganization.

## Constraints

- `bash _dev/tests/contract-regressions.sh` must pass clean.
- Crew files ship — they must cite other files by shipped path and must never cite `CLAUDE.md`/`AGENTS.md`.
- Version bump + descriptive `CHANGELOG.md` entry.

## Dependencies

None upstream. Independent of the other REQs; touches no file they touch.

## Builder Guidance

**Certainty: Firm on the test, exploratory on each line's verdict.** Expect roughly half of each file to go. When genuinely torn about a line, apply `crew-members/maintenance.md` — delete-before-you-add means the default is delete, and the git history keeps it recoverable. Do surface, in Discovered Tasks rather than inline, anything you find that belongs in a different crew file altogether.

## Red-Green Proof

- **RED now:** `wc -w crew-members/{backend,frontend,security,testing,ui-design,debugging}.md` totals ~6,737. `grep -rn 'dangerouslySetInnerHTML\|token storage\|API key' crew-members/` returns hits in more than one file — the same security guidance stated in two places.
- **GREEN when:** the six-file total is ~3,500; each overlapping security topic appears in exactly one file (`crew-members/security.md`), with `frontend.md`/`backend.md` pointing at it; every `JIT_CONTEXT` line is unchanged.
- **Validation:** per-file `wc -w` receipt; the dedup grep before/after; the `JIT_CONTEXT` byte-identical check; the crew-pointer resolution grep; `bash _dev/tests/contract-regressions.sh` clean.

## Open Questions

None — the tradeoff was presented and accepted at capture.

## Full Context

See `do-work/user-requests/UR-006/input.md` for complete verbatim input.

---
*Source: "compare with the current skill, is there something that we need to update?" — resolved into the approved seven-REQ plan.*

Think carefully before answering.

## Triage

**Route A** — surgical line-by-line trim across six known files, no exploration needed. The keep test and the overlap list were both fully specified in the REQ; the only judgment calls were per-line verdicts, which the Builder Guidance already anticipated ("expect roughly half of each file to go").

## Scope

Files committed to before editing (exactly the REQ's six, no others):
- `crew-members/backend.md`
- `crew-members/frontend.md`
- `crew-members/security.md`
- `crew-members/testing.md`
- `crew-members/ui-design.md`
- `crew-members/debugging.md`

Acceptance criteria restated from the REQ:
- Six-file total trimmed toward ~3,500 words (direction, not quota) from 6,737 — no padding to hit it, no deleting a genuine opinion to get under it.
- Per-line keep test applied and defensible: survives only as (a) a deliberate preference of this project or (b) tied to do-work machinery; "it's true" is not a reason.
- Security dedup resolved first: `security.md` owns all security content; `frontend.md`/`backend.md` drop overlapping bullets (`dangerouslySetInnerHTML`, token storage, API keys in bundles, resource-level authorization, rate limiting) and point at it. Nothing said once ends up said zero times — unique content moves before deletion.
- Every `JIT_CONTEXT` comment byte-identical before/after.
- No new files.
- `bash _dev/tests/contract-regressions.sh` passes clean.
- Version bump + CHANGELOG entry (orchestrator's job, not mine).

## Implementation Summary

**What was done:** Applied the per-line keep test to all six files. `backend.md` and `frontend.md` lost nearly all of their generic-engineering-advice bulk (API design boilerplate, state-management conventions, dependency hygiene, scope-discipline restatements already covered by the always-loaded `coding-guardrails.md`) and kept only machinery-tied or genuinely specific lines (Decision D-XX versioning, Discovered-Tasks-tied caching note, UNIFY quality-check tables, the compositor-thread animation rule). Their security bullets were removed; the two that weren't already duplicated in `security.md` (Vue `v-html`/Angular `[innerHTML]` naming, `rel="noopener noreferrer"` tabnabbing protection) were moved into `security.md`'s framework-specific section (renamed `React / Frontend` → `React / Vue / Angular` since it now covers all three) before the frontend bullets were deleted — nothing lost. `security.md` itself was trimmed of its most generic prose (the "What Static Analysis Catches/Misses" textbook explanation, condensed into one clause; a vague "Threat modeling" row; two Anti-Patterns bullets already covered by the OWASP table rows) but its OWASP checklist and framework-specific patterns were kept close to intact — each row names a specific algorithm/threshold/function, which is exactly the "opinion, not platitude" bar the keep test sets, and both are hard-cited by name from `actions/code-review.md`. `testing.md` lost the AAA-structure/naming/assertion-style prose (standard practice), the Flaky Test Prevention causes table (textbook), and the Coverage Expectations section (the REQ's own example of generic-cut content) — kept the caller-seam test, production-faithful fixtures, the TDD Red-Green workflow (machinery-tied), and four non-obvious Anti-Patterns. `debugging.md` condensed the Investigation Techniques and Cognitive Bias Guards from prose blocks to one-liners (cutting an exact internal duplicate: the Anchoring-bias countermeasure restated "When to Escalate"'s 2-attempt rule word-for-word) while keeping the Scientific Method, Confidence Levels, and When-to-Escalate sections, all of which are machinery-tied (REQ Lessons Learned, `error_type` classification, prime-file lesson capture). `ui-design.md` had the least room to cut — `actions/ui-review.md` cites its 6 phases, Heuristic Review Criteria, Severity definitions, Accessibility Baseline, and Implementation Patterns by name — so it was prose-trimmed (verbose numbered sub-steps collapsed to single sentences, the Scope Discipline subsection folded into one closing sentence in Quality Checks) rather than gutted; every phase and cited section still exists under its original name.

**Files changed:**
- `crew-members/backend.md` (modified) — 801 → 202 words.
- `crew-members/frontend.md` (modified) — 547 → 213 words.
- `crew-members/security.md` (modified) — 1,674 → 1,534 words (absorbed 2 moved bullets from frontend.md; trimmed generic prose elsewhere).
- `crew-members/testing.md` (modified) — 1,473 → 587 words.
- `crew-members/ui-design.md` (modified) — 1,082 → 849 words.
- `crew-members/debugging.md` (modified) — 1,160 → 852 words.

**Word-count receipt:**

| File | Before | After |
|---|---|---|
| `crew-members/backend.md` | 801 | 202 |
| `crew-members/frontend.md` | 547 | 213 |
| `crew-members/security.md` | 1,674 | 1,534 |
| `crew-members/testing.md` | 1,473 | 587 |
| `crew-members/ui-design.md` | 1,082 | 849 |
| `crew-members/debugging.md` | 1,160 | 852 |
| **Six-file total** | **6,737** | **4,237** |

Six-file total landed at 4,237 (37% cut), above the ~3,500 direction. Two files carry real structural floors that block further cutting without breaking a hard citation: `security.md` (1,534w) is the named dedup owner and its OWASP/framework-pattern content is specific/concrete — exactly what the keep test calls an opinion, not a platitude — and both sections are cited by name from `actions/code-review.md`; `ui-design.md` (849w) has five section-level citations from `actions/ui-review.md` (6 phases, Heuristic Review Criteria, Severity definitions, Accessibility Baseline, Implementation Patterns) that would break if gutted further. The other four files absorbed the shortfall — backend (−75%), frontend (−61%), testing (−60%), debugging (−27%) — landing the six-file average cut near the "roughly half" the Builder Guidance anticipated, even though the distribution across files is uneven. Per the REQ's own instruction ("the target is a direction, not a quota"), no further lines were cut to force the number down, and no line was kept just to pad it.

Always-read floor (unchanged — not in scope): `SKILL.md` (2,557w) + `next-steps.md` (1,741w) = 4,298w before and after.
Orchestrator load (unchanged — not in scope): `SKILL.md` + `actions/work.md` + `actions/work-reference.md` + `crew-members/general.md` + `crew-members/coding-guardrails.md` + `next-steps.md` = 2,557 + 11,983 + 7,426 + 678 + 846 + 1,741 = 25,231w before and after.

## Testing

- `bash _dev/tests/contract-regressions.sh` → `Contract regression checks passed.` (clean, run after all six edits).
- Per-file `wc -w crew-members/{backend,frontend,security,testing,ui-design,debugging}.md` before and after — see receipt table above.
- Dedup grep before: `grep -rn 'dangerouslySetInnerHTML\|token storage\|API key' crew-members/` hit both `frontend.md` and `security.md`. After: hits only `security.md` for substantive content, plus one pointer-only mention of "token storage" in `frontend.md`'s one-line redirect to `security.md` (not a restatement of the guidance itself).
- `JIT_CONTEXT` byte-identical check: `grep -n 'JIT_CONTEXT' crew-members/{backend,frontend,security,testing,ui-design,debugging}.md` before and after — all six lines identical, character for character.
- Crew-pointer resolution grep: `grep -rn 'crew-members/\(backend\|frontend\|security\|testing\|ui-design\|debugging\)\.md' crew-members/ actions/` — every hit (`actions/ui-review.md` ×6, `actions/work.md` ×5, `actions/code-review.md` ×1, plus the two new pointer lines in `backend.md`/`frontend.md`) resolves to a section/concept that still exists in the target file under the same name.
- No-CLAUDE/AGENTS-citation check: `grep -rniE 'CLAUDE\.md|AGENTS\.md' crew-members/{backend,frontend,security,testing,ui-design,debugging}.md` → no matches.

## Lessons Learned

**What worked:** Running the pointer-dependency grep *before* touching any content (rather than trimming and checking afterward) meant every cut for `ui-design.md` and `security.md` was made with the hard constraints already known, so no rework was needed. Deciding the dedup moves (Vue/Angular naming, `rel=noopener`) before deleting the source bullets kept "nothing said once ends up said zero times" mechanically checkable via the before/after dedup grep.

**What didn't:** My first pass toward backend.md and frontend.md undershot the "roughly half" expectation by cutting closer to 75%/61% — the REQ's own worked examples ("use appropriate HTTP status codes", "don't chase a coverage number") calibrate the keep-test bar stricter than a first read suggests, and once applied consistently, most of those two files' generic-engineering-advice content had no defensible "preference or machinery" justification left.

**Worth knowing:** The keep test's two buckets (preference / machinery) don't just filter platitudes — they also surface *specificity* as a proxy: a line naming a concrete algorithm, threshold, or function (`bcrypt cost 12+`, `AES-256-GCM`, `169.254.169.254`) reads as an opinion even when the general topic is well-known, while a line with no concrete parameter ("use appropriate HTTP status codes") reads as generic regardless of topic. That's why `security.md`'s OWASP checklist survived largely intact while `backend.md`'s API Design section didn't.

## Discovered Tasks

- `crew-members/frontend.md` (Animation & Rendering Performance / kept "Opinions" section) and `crew-members/ui-design.md` (Phase 6: Interaction & Motion) both state the `prefers-reduced-motion` rule and transition-duration guidance. This overlap wasn't in this REQ's named scope (only frontend↔security and backend↔security were called out) and I didn't touch it, but it's the same shape of duplication this REQ just resolved elsewhere — worth a follow-up REQ if a similar frontend↔ui-design consolidation pass is wanted.
