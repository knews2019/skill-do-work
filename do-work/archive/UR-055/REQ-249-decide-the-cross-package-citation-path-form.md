---
id: REQ-249
title: Decide the cross-package citation path form and sweep to match
status: completed
created_at: 2026-08-18T13:54:59Z
claimed_at: 2026-08-18T16:09:27Z
completed_at: 2026-08-18T18:08:42Z
commit: cc1083c
kb_status: pending
route: B
status_changed_at: 2026-08-18T14:12:05Z
user_request: UR-055
addendum_to: REQ-244
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: true
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-08-18T16:10:30Z
  basis:
    - Route B
    - 15-file write set
    - 4 subsystems involved
    - 3 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set:
- _dev/primes/prime-action-files.md
- skills/do-work/**/*.md
- skills/do-work-board/**/*.md
- skills/do-work-knowledge/**/*.md
- skills/do-work-toolbox/**/*.md
- _dev/tests/shipped-package-reference-contract.sh
---

# Decide the Cross-Package Citation Path Form and Sweep to Match

## What

Two incompatible readings of the cross-package citation path coexist in shipped markdown, and nothing can tell them apart. Pick one and sweep.

## Context

Raised by REQ-244's builder as pushback rather than decided mid-REQ, verified independently by the orchestrator, and confirmed again by REQ-244's review. REQ-244 added eleven more citations in the prescribed form while the question was open, so the count grows until this is settled.

- `\`../do-work/actions/work-reference.md\`` is what `_dev/primes/prime-action-files.md:91` prescribes and what `actions/memory.md` uses. From `skills/do-work-toolbox/actions/` it resolves to `skills/do-work-toolbox/do-work/actions/...`, **which does not exist.** It coheres only read as *skill-root*-relative rather than relative to the citing file's own directory.
- `\`../../do-work/actions/work-reference.md\`` is what `skills/do-work-toolbox/actions/present-work.md:37` already uses. It resolves literally, in both the source and installed topologies.

Counts at capture, verified by REQ-244's review: 30 local-form citations in core, 3 sibling-form in do-work-board, 11 sibling-form in do-work-knowledge and do-work-toolbox — consistent by package, so whichever reading is intended, the corpus is at least internally regular.

**REQ-243's checker cannot arbitrate this.** It resolves Markdown `[text](target)` link syntax; every one of these citations is a backticked path, which that checker never sees. So the correctness of cross-package pointers currently rests on convention alone.

## Open Questions

- [x] Shipped action files cite each other across skill packages with a backticked path, and there are two spellings in the tree that mean different things. One (`../do-work/actions/...`) is what the prime file tells writers to use, and it only makes sense if you read the `../` as "up to the skills folder" rather than as a real relative path — typed literally into a terminal from the citing file's folder, it points at nothing. The other (`../../do-work/actions/...`) is a real path that works from where the file actually sits. Both are in use today, the first far more than the second. Nothing checks either, because our new link checker only understands Markdown links and these are backticks. Which reading do you want to be the rule — the skill-root one that most files already use, or the literal one that a reader could paste and follow? Whichever you pick, the other spelling gets swept to match and the prime file gets updated to say so. → Confirmed: literal paths (`../../`), so a citation a reader pastes actually resolves and a future checker can verify it mechanically.
  Recommended: literal paths (`../../`), so a citation a reader pastes actually resolves and a future checker can verify it mechanically.
  Also: keep skill-root-relative (`../`), sweep `present-work.md` and `completed-work-presentation-reference.md` to match, and state the convention explicitly in the prime so nobody reads it as a filesystem path.

**Answered [2026-08-18]:** User confirmed the recommended literal form via `do-work clarify`. The rule is: every backticked cross-package citation must resolve as a real relative path from the citing file's own directory. Sweep every citation that doesn't (the Context counts are a starting set, not the extent — re-derive them), update `_dev/primes/prime-action-files.md:91` to prescribe the literal form, and answer the Requirements question in favor of mechanical checkability: backticked cross-package pointers should become checkable now that they resolve literally.

## Requirements

- One spelling is the documented rule, stated in `_dev/primes/prime-action-files.md`.
- Every shipped citation matches it — the sweep is the requirement, and the counts above are a starting set rather than the extent.
- Whether backticked cross-package pointers become mechanically checkable is part of the answer: if the literal form wins, say whether a checker should now read them.

---

## Triage

**Route: B** - Medium

**Reasoning:** The decision is already made; what is unknown is the extent — the capture counts are explicitly a starting set, so every citation across four shipped packages has to be found before any of them is rewritten.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `_dev/primes/prime-action-files.md` (modify) — § Cross-Referencing states the literal form as the rule
- Shipped markdown under `skills/do-work/`, `skills/do-work-board/`, `skills/do-work-knowledge/`, `skills/do-work-toolbox/` (modify) — every backticked cross-package citation swept to the spelling that resolves from the citing file's own directory
- `_dev/tests/shipped-package-reference-contract.sh` (modify) — the mechanical check for backticked cross-package pointers, if the builder judges one is warranted

**Files I will NOT touch:**
- `skills/do-work-board/tools/queue-kanban/**` Go sources and `web/*.js` — REQ-248 is building there in this same wave.
- `CHANGELOG.md`, `VERSION`, `skills/do-work/actions/version.md` — serial-only, owned by the integrator.
- Markdown `[text](target)` links — REQ-243's checker already governs those; this REQ is about backticked paths.

**Acceptance criteria (restated from REQ):**
- [ ] One spelling is the documented rule, stated in `_dev/primes/prime-action-files.md`.
- [ ] Every shipped backticked cross-package citation resolves literally from its own file's directory — the capture counts are a starting set, not the extent (measured 126 local-form / 14 sibling-form at claim time, and package-root `SKILL.md` citations are already literal-correct).
- [ ] The answer states whether backticked cross-package pointers become mechanically checkable, and implements the check if the answer is yes.
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0.

## Implementation Summary

**What was done:** Every backticked cross-package citation in shipped markdown now resolves as a real relative path from the citing file's own directory, in both source and installed topologies — 122 citations rewritten across 46 files, 19 already-correct ones untouched, depth derived per file (`../../` from `actions/`/`crew-members/`/`docs/`, `../../../` from `tools/queue-kanban/`). `_dev/primes/prime-action-files.md` § Cross-Referencing states the literal form as the rule, retires the skill-root-relative shorthand by name, exempts fenced template/example blocks, and points at the checker. The mechanical-checkability question is answered **yes and implemented**: `_dev/tests/shipped-package-reference-contract.sh` now verifies backticked cross-package citations in both topologies, reusing the existing CommonMark walk via a behavior-identical extraction refactor (`mask_block_code` / `inline_code_regions`), locked by the existing parser fixtures; against the pre-sweep tree it reports exactly the 122 swept citations.

**Files changed (48, +319/−122):**
- `_dev/primes/prime-action-files.md` (modified) — the literal rule stated in § Cross-Referencing
- `_dev/tests/shipped-package-reference-contract.sh` (modified) — backticked-citation check + parser refactor + fixtures
- `skills/do-work-board/actions/board.md` (modified) — citations swept to the literal form
- `skills/do-work-board/docs/board-guide.md` (modified) — citations swept to the literal form
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified) — citations swept to the literal form
- `skills/do-work-knowledge/actions/bkb.md` (modified) — citations swept to the literal form
- `skills/do-work-knowledge/actions/dream.md` (modified) — citations swept to the literal form
- `skills/do-work-knowledge/actions/interview-reference.md` (modified) — citations swept to the literal form
- `skills/do-work-knowledge/actions/interview.md` (modified) — citations swept to the literal form
- `skills/do-work-knowledge/actions/memory-reference.md` (modified) — citations swept to the literal form
- `skills/do-work-knowledge/actions/memory.md` (modified) — citations swept to the literal form
- `skills/do-work-knowledge/crew-members/anti-slop.md` (modified) — citations swept to the literal form
- `skills/do-work-knowledge/crew-members/background-agents.md` (modified) — citations swept to the literal form
- `skills/do-work-knowledge/crew-members/clear-questions.md` (modified) — citations swept to the literal form
- `skills/do-work-knowledge/crew-members/maintenance.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/code-review.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/deep-explore-reference.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/deep-explore.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/inspect.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/install.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/maintainability-audit.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/note.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/present-work.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/quick-wins.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/stray-check.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/tidy-repo.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/ui-review.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/actions/validate-feedback.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/crew-members/anti-slop.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/crew-members/background-agents.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/crew-members/clear-questions.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/crew-members/coding-guardrails.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/crew-members/debugging.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/crew-members/maintenance.md` (modified) — citations swept to the literal form
- `skills/do-work-toolbox/crew-members/security.md` (modified) — citations swept to the literal form
- `skills/do-work/actions/capture-reference.md` (modified) — citations swept to the literal form
- `skills/do-work/actions/capture.md` (modified) — citations swept to the literal form
- `skills/do-work/actions/commit.md` (modified) — citations swept to the literal form
- `skills/do-work/actions/forensics.md` (modified) — citations swept to the literal form
- `skills/do-work/actions/kb-lessons-handoff.md` (modified) — citations swept to the literal form
- `skills/do-work/actions/review-work.md` (modified) — citations swept to the literal form
- `skills/do-work/actions/roadmap.md` (modified) — citations swept to the literal form
- `skills/do-work/actions/work-reference.md` (modified) — citations swept to the literal form
- `skills/do-work/actions/work.md` (modified) — citations swept to the literal form
- `skills/do-work/crew-members/general.md` (modified) — citations swept to the literal form
- `skills/do-work/crew-members/maintenance.md` (modified) — citations swept to the literal form
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified) — citations swept to the literal form
- `skills/do-work/docs/standing-preferences.md` (modified) — citations swept to the literal form

No serial-only file touched. Deliberately exempt: fenced template/example blocks (D-01), Markdown link syntax (REQ-243's checker), `<skill-root>/../` invocation forms.

*Integrated by orchestrator from builder hand-back; merge range `e8c6f80..cc1083c`.*

## Decisions

Transcribed by the orchestrator from the builder hand-back:

- **D-01 (DECIDE & STATE): fenced blocks are exempt from sweep and checker** — fenced occurrences are content destined for a consumer's files, not citations from the citing file; no single literal path is correct from their eventual location. Matches REQ-243's checker. Stated in the prime.
- **D-02 (DECIDE & STATE): HTML comments are checked** — JIT_CONTEXT headers carry real citations; masking comments would leave a 12-citation hole exactly where the always-loaded instruction surface lives (the standing class-vs-instance warning's shape).
- **D-03 (DECIDE & STATE): the core-package/queue-root name collision is a documented skip** — `do-work` names both the core package and the consumer queue root, so a span whose tail names no real package content is skipped; residual: a citation to a *deleted core* file is indistinguishable from a consumer-state example. The three other packages still catch deletions. Keyed on the collision condition, not a hand list.
- **D-04 (DECIDE & STATE): first whitespace token is the path** — invocation spans resolve only the leading token, so argument placeholders never false-positive.
- **D-05 (DECIDE & STATE): parser refactor over a second parser** — extracting `mask_block_code`/`inline_code_regions` from `strip_markdown_code` (fixtures prove identity) beats a parallel fence-walker that would drift.

## Qualification

Passed — 48 files verified in merge range `e8c6f80..cc1083c`, requirements traced (rule stated in `prime-action-files.md` § Cross-Referencing with per-file depth examples and the retired shorthand named; spot-resolved swept citations from `present-work.md` land on real files; checker present in the reference contract), P-A-U audited (sweep scripts stayed in session scratch, no stray binary, no debug artifacts in the merged diff). The 122-vs-140 count difference from the brief is explained and verified by the builder's condition-derived extent (fenced content and `<skill-root>/../` invocations are not citations under the decided rule).

## Review

**Overall: 97%** | 2026-08-18T18:05:59Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

### Requirements Checklist

- [x] **One spelling is the documented rule, stated in `_dev/primes/prime-action-files.md`** — delivered: § Cross-Referencing states the literal rule with per-file depth examples, retires the skill-root shorthand by REQ number, states the fenced-block exemption, points at the checker. *(reproduced-by-reading)*
- [x] **Every shipped backticked citation matches it** — delivered, independently confirmed. The reviewer's own scanner (deliberately not derived from the checker) over all tracked `skills/**/*.md`: 146 citation-shaped spans; 141 resolve in **both** the source topology and the installed four-sibling topology; the 5 non-resolving are verified-legitimate non-citations (a generic refactor example, a templated worktrees path, three consumer-state lesson-link computation examples). **Zero spans the scanner flags that the shipped checker passes.** Counts reconcile exactly: 122 rewritten + 19 already-correct = 141. *(reproduced-by-execution)*
- [x] **Mechanical-checkability answered and implemented** — against the pre-sweep tree the checker reports **exactly 122** failures, reproducing the Implementation Summary's claim verbatim. *(reproduced-by-execution)*
- [x] **`bash _dev/tests/maintainer-verify.sh` exits 0** — reproduced un-piped after all mutation probes were restored.

### Mutation testing of the checker (all restored; final `git status` clean)

Wrong-depth prose citation → FAIL with file:line ✓ · deleted non-core file → FAIL ✓ · broken citation in a fenced block → PASS (D-01) ✓ · broken backticked citation in an HTML comment → FAIL (D-02) ✓ · D-03 residual probed deliberately: a prose citation to a deleted **core** file passes silently — real, and honestly documented, condition-keyed on the name collision and pinned by a fixture ✓ · `<skill-root>/../` invocation → skipped ✓ · wrong-depth citation with arguments → FAIL on first token only (D-04) ✓ · bogus `#anchor` on a correct citation → passes (Minor 1).

**Parser refactor (D-05):** read line-by-line against the pre-image — behavior-identity confirmed; pre-existing parser fixtures pass; live corpus produces identical link-check results. *(argued-from-reading, fixture-locked)*

### Findings

**Important:**

- **The retired skill-root-relative reading survives outside backtick spans — three shipped sites the sweep and the checker are both structurally blind to.** (1) `skills/do-work/SKILL.md:16` — the retired resolution rule stated as prose in the core router, now contradicting the prime; applied to the new `../../` form it computes wrong paths. (2) `skills/do-work/actions/commit.md:17` — bare `../do-work-toolbox/actions/inspect.md`, wrong depth from `actions/`. (3) `skills/do-work/crew-members/security.md:3` — bare `../do-work-toolbox/actions/code-review.md` in the JIT_CONTEXT comment (the checker's comment scan covers backticked spans only). A Restatement Sweep finding, **not** builder scope drift: the user-confirmed rule's letter covers backticked citations, and these are bare text — but the diff retired the *reading*, and these sites still state or use it. *(reproduced-by-execution)* — gate: **rule-change** → sweep REQ-259 (`sweep_key: retired-reading-outside-backtick-spans`), created `pending-answers` per the generation-≥2 reroute.

**Minor:** 3 (report only)
1. Anchor fragments on backticked citations are discarded (`.split("#")[0]`), never validated. Latent: no live backticked cross-package citation carries an anchor today.
2. The fenced-yaml Request File Schema comments (`work-reference.md:130/132/137`) carry live parser-lock-step citations in the retired form, permanently checker-exempt via D-01 — the exemption reasoning fits true templates but sits awkwardly on schema documentation a maintainer follows.
3. `scope-drift.sh` exits 1 flagging all 46 swept files: the Scope section declares package trees as prose and `write_set` covers them as globs (all 48 changed files inside it), but the tool only reads enumerated paths. Tool limitation surfaced by sweep-shaped REQs, not builder drift.

**Confirmed exemptions (verified, not re-litigated):** the `code-review.md`/`validate-feedback.md` fenced templates (D-01's reasoning holds); `install-last30days.sh:133`'s shell path is runtime-anchored and already resolves in both topologies. No serial-only file in the merge range.

### Acceptance Testing

**Result: Pass** — checker catches every mutation class it claims (7 live probes), passes the post-sweep corpus, reports exactly 122 on the pre-sweep tree, independent scan finds zero holes in the backticked class, verify exit 0 after restoration.

**Suggested testing:** (1) run the checker in a real consumer install once to confirm the modeled installed topology matches a deployed one; (2) add anchor validation before backticked citations grow anchors.

### Reviewer-Recommended Disposition

**Approve with follow-ups.** The decided backticked-citation class is genuinely closed (sweep complete, checker sound, both topologies modeled, residuals honestly documented); the retired *reading* survives at three unbackticked sites — exactly the instance-vs-class shape this review hunts — landing as one root-cause sweep, rerouted `pending-answers` per the depth cap.

**Important findings (audit record):**
- Retired reading survives at three unbackticked shipped sites — gate: rule-change → sweep REQ-259 (pending-answers)

**Acceptance:** Pass · **Follow-ups created:** REQ-259 (by orchestrator, orchestrated mode) · **Sweeps appended to:** none (no matching candidate under UR-055)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Deriving the extent mechanically (condition over the corpus) instead of trusting the capture counts — the brief's grep said 140, the condition said 122+19, and the review's independent scanner reconciled to the same 141 exactly. Making the sweep's condition BE the checker's condition means the two cannot disagree by construction. Refactoring the existing CommonMark walk instead of writing a second parser kept the fixtures as the identity proof.

**What didn't:** The class boundary was drawn at backtick spans — the letter of the decided rule — and the retired *reading* survived at three bare-text sites, including the core SKILL.md stating the old resolution rule as prose (REQ-259). Seventh consecutive occurrence of the instance-vs-class shape; the class this time was "the reading", not "the spelling".

**Worth knowing:** `do-work` names both the core package and the consumer queue root, so the checker deliberately skips core-package spans whose tail names nothing real — a citation to a *deleted core* file is invisible to it (documented, fixture-pinned; the other three packages still catch deletions). Anchor fragments on backticked citations are discarded, not validated. Fenced blocks are exempt by design.

## Orientation

Now every backticked cross-package citation in shipped markdown is a real relative path a reader can paste, in both source and installed topologies, and the reference contract enforces it mechanically. Lives in the shipped-docs convention layer (`_dev/primes/prime-action-files.md` § Cross-Referencing) plus `_dev/tests/shipped-package-reference-contract.sh`. [MAP CHANGED] — the citation convention flipped from skill-root-relative to literal, and a checker now guards the class; writers of new action files inherit the rule from the prime. Prime staleness spot-check: `prime-action-files.md` paths still resolve; not stale.

## Pre-Flight

**Git:** ✓ clean outside `do-work/`
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 (recorded in `do-work/working/baseline.json`)
**Dependencies:** ⚠ this checkout needed Go 1.26.1, ShellCheck 0.11.0 and `just` installed before the baseline could run at all, and one pre-existing Linux-only test failure had to be fixed first (0.212.8) — see the REQ brief.

*Checked by work action*
