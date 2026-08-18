---
id: REQ-250
title: Close the remaining markdown link checker gaps
status: completed
created_at: 2026-08-18T13:55:32Z
claimed_at: 2026-08-18T18:25:40Z
completed_at: 2026-08-18T18:58:44Z
commit: 330797b
kb_status: pending
route: B
status_changed_at: 2026-08-18T13:55:32Z
user_request: UR-042
addendum_to: REQ-243
domain: general
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: markdown-link-checker-unresolved-classes
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- _dev/tests/shipped-package-reference-contract.sh
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T18:26:21Z
  basis:
    - Route B
    - 1-file write set
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Close the Remaining Markdown Link Checker Gaps

## What

REQ-243's anchor checker has four known limits. Three are false negatives — links it accepts that are broken — and one is a path-escape the same change escalated from a `stat` to a `read_text`.

## Context

All four came from REQ-243's independent review, which cleared the checker's core (skip rule condition-keyed, slug rule correct for every real anchor in this repo, corpus failure path non-vacuous) and named these as the remaining edges. They were routed here rather than folded in, because REQ-243's remediation was already carrying two more urgent findings.

## Instances

- [ ] **Same-file bare `#anchor` links are not validated at all.** `[see below](#renamed-section)` is silently broken today. The slugs are already cached by REQ-243's work, so this is roughly three lines.
- [ ] **`os.path.normpath` does not clamp `..` at the repo root, and this path is now read rather than stat'ed.** A link like `../../../../elsewhere.md#a` that happens to exist was previously only tested for existence; the anchor check now calls `read_text()` on it. Pre-existing in origin, escalated in effect. A `relative_to(repo_root)` guard closes it.
- [ ] **HTML tags and entities in heading text diverge from GitHub's slug, in the false-negative direction.** `## <kbd>Ctrl</kbd> and stuff` slugs to `kbdctrlkbd-and-stuff` here and `ctrl-and-stuff` on GitHub; `## Tom &amp; Jerry` gives `tom-amp-jerry` against `tom--jerry`. A link written to the checker's slug passes the suite and is broken in every renderer. No shipped heading trips this today — all 27 with suspicious characters slug correctly — so **check whether this is worth closing at all before closing it**; a stated limitation may be the better answer.
- [ ] **Blockquoted ATX headings are dropped.** `> # Quoted Heading` yields no anchor because the gate tests the code-masked line, which begins with `>`. GitHub does generate one. Fails loudly (spurious FAIL), so low severity — but the limitation is currently unstated.

## Requirements

- Each instance is either closed, or its limitation is stated in the file with the failure direction named. A silently-accepting limit and a loudly-failing one are different risks and the file should say which each is.
- `bash _dev/tests/maintainer-verify.sh` still exits 0, and the existing 27 corpus anchors still resolve.
- **`maintenance: true`: ask what can be removed before adding.** Two of these four may be better closed by deleting a claim than by adding code.

## Implementation Summary

**What was done:** All four checker edges resolved — two closed with code, two documented as stated limitations. (1) Bare `#anchor` links now validate: the fragment resolves as the carrying file's own name through the existing target-and-anchor pipeline; the corpus's 3 real bare-anchor links now get checked. (2) The `..`-escape class is clamped at every normpath-then-probe site — and the hole-hunt found the genuinely silent instance in REQ-249's citation checker, whose consumer-queue probe absorbed interior-`..` tails and could stat outside the repo (demonstrated reaching `/etc/hostname`); both citation sites are clamped and the silent-absorb hole is fixture-locked. (3) HTML-tag/entity slug divergence documented with both failure directions named, pinned by a fixture. (4) Blockquoted-heading drop documented as always-loud, pinned by a fixture.

**Files changed:**
- `_dev/tests/shipped-package-reference-contract.sh` (modified) — the only file (+67/−8): bare-anchor routing, escape clamps at three sites, two limitation statements, three fixtures

*Integrated by orchestrator from builder hand-back; merge range `beb3b7b..330797b`.*

## Decisions

Transcribed from the builder hand-back:

- **D-01 (DECIDE & STATE):** bare fragments resolve through the existing pipeline rather than a parallel anchor branch — dedup, caching and message consistency for ~6 lines.
- **D-02 (DECIDE & STATE):** the clamp is `".." in <normalized>.parts` at every probe site rather than `relative_to(repo_root)` — in repo-relative PurePosixPaths any surviving `..` after normpath IS the escape condition; both topologies guarded so the invariant is stated rather than leaning on destination-root geometry.
- **D-03 (DECIDE & STATE, maintenance):** instance 3 documented, not closed — no live corpus case, and correct closing needs code-span-aware HTML decoding; fails the earned-defense test. Pinning fixture keeps statement and behavior from drifting apart.
- **D-04 (DECIDE & STATE, maintenance):** instance 4 documented, not closed — always loud, zero corpus cases, correct blockquote anchors need nested-quote machinery with no live case.
- **D-05 (DECIDE & STATE):** the citation checker's interior-`..` clamp is in-scope, not a REQ-249 rework — instance 2's exact class at a second copy of the same primitive; all seven prior citation fixtures still pass.

## Qualification

Passed — 1 file verified in merge range `beb3b7b..330797b` (+67/−8), all three requirement clusters traced (each instance closed-or-stated with failure direction; suite green with the 27 pre-existing anchors plus 3 newly validated; maintenance posture honored — two edges resolved by statement, not machinery), P-A-U audited. The builder's pushback on instance 2's framing (the corpus read-risk was overstated; the real silent hole was in the citation checker's consumer-queue probe) is verified by its RED-2a/RED-2b transcripts and accepted as the more accurate account.

## Review

**Overall: 97%** | 2026-08-18T19:12:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Verdict: Approve** — all four instances land as demanded (two closed by execution-verified mechanism, two stated with failure directions and mutation-verified pins); the `..`-escape class claim survives adversarial enumeration; maintenance posture honored (the core fix is a deleted skip, not new machinery).

### Requirements Checklist

- [x] **Instance 1 (bare anchors): closed.** *Reproduced by execution:* injected broken bare anchor fails at HEAD with the right message, silently passes at `beb3b7b`. Instrumented count 27 → 30 anchor-checked links; the 3 new are exactly the corpus's live bare-anchor links, all resolving. Fence/comment masking verified; uppercase fragment fails loudly (correct — GitHub slugs are lowercase); `%`-encoded fragment decodes.
- [x] **Instance 2 (`..`-escape): closed class-wide.** *Reproduced by execution, class verified by enumeration:* escaping link fails pre-probe at HEAD; interior-`..` citation silently absorbed at `beb3b7b` (exit 0) and fails at HEAD; the `/etc/hostname` reach reproduced via probe-recording harness at `beb3b7b` (two out-of-repo stats) with **zero** filesystem probes at HEAD for the same span. Enumeration of every corpus-path-to-filesystem site found no unclamped probe; the citation installed-topology probe is structurally safe (leading-`..` can never `relative_to`-match a destination root) and its branch is fixture-executed.
- [x] **Instance 3 (HTML slug divergence): documented, not closed** — judged sound as the REQ invited: zero live cases, closure needs code-span-aware HTML decoding, fails earned-defense. The statement's claimed GitHub slugs verified accurate by mutation.
- [x] **Instance 4 (blockquoted headings): documented as always-loud** — direction correct, zero corpus cases, same sound reasoning.
- [x] **`maintainer-verify` exit 0, 27 anchors still resolve** — observed on the full un-piped run after all probes reverted (`git status --porcelain` empty).
- [x] **Maintenance clause honored** — instance 1 closes by *deleting* a skip and reusing the pipeline; the +67 is dominated by limitation prose and pins.

### Findings

**Important:** None

**Minor (report only, 2):**
1. **The REQ's D-02 transcription overstated clamp uniformity** — the citation checker's installed-topology probe has no explicit clamp; its safety leans on destination-root geometry (real, fixture-executed). The builder's hand-back scoped the claim correctly ("both **main-loop** topologies"); the qualifier was dropped in the orchestrator's transcription. *(Corrected in the Decisions section below this review.)*
2. **The recorded pushback is itself slightly overstated:** pre-change, a *re-entrant* escape (`../../../../.claude/skills/...`) could silently pass both topologies on a machine with the suite installed at home-level `.claude/` — the standard consumer topology — so the corpus read path had a realistic silent case too. No code consequence: HEAD fails the shape pre-probe (verified).

**Restatement sweep:** every remaining statement of the old behavior lives in frozen records; the only other `startswith("#")` in the codebase is an unrelated TSV filter. No stale live restatement.

**Scope:** `scope-drift.sh` exit 0; one file, matching the write set; the clamp edits inside REQ-249's checker are D-05's documented same-primitive judgment, with all seven prior citation fixtures still passing beside the new eighth.

### Acceptance Testing

**Result: Pass** — four adversarial REDs injected and reverted, each failing at HEAD with the right message and the two claimed silent holes reproducing at `beb3b7b`; all three pinning fixtures proven mutation-falsifiable (none vacuous); full gate exit 0 on a clean tree after reverts.

### Suggested Additional Testing

- Symlink escape is a different, unclamped-by-design class (no shipped symlinks today) — worth a limitation line or a mode probe if one ever appears.
- Consumer-topology spot check of the re-entrant escape shape on a machine with `~/.claude/skills/` installed.

**Follow-ups created:** None · **Sweeps appended to:** None

*Reviewed by review-work action (independent adversarial pass, orchestrated mode; merge range `beb3b7b..330797b`)*

### Orchestrator correction (review Minor 1)

D-02 as transcribed above should read: the `".." in parts` clamp guards **both main-loop topologies** explicitly; the citation checker's consumer-queue and source probes carry their own clamps, and its installed-topology probe is safe by destination-root geometry (leading-`..` cannot `relative_to`-match), executed by the escaping-tail fixture rather than clamped explicitly.

## Lessons Learned

**What worked:** The hole-hunt discipline finally beat the class-vs-instance curse — greping the same primitive (`normpath`-then-probe) across the file before calling the class closed found the genuinely silent hole in a *different* checker's probe, and the review's independent enumeration then found nothing left. Mutation-falsifiable pins: each documented limitation carries a fixture that FAILS if the behavior changes, so statement and behavior cannot drift apart silently.

**What didn't:** Two records-precision slips, both in the orchestrator's transcription rather than the code: D-02's clamp-uniformity claim dropped the builder's "main-loop" qualifier, and the recorded pushback missed the re-entrant-escape case where the pre-change corpus read path was silently unsafe on the standard consumer topology. The builder's own hand-back was the more accurate record both times — transcribe, don't paraphrase.

**Worth knowing:** The escape clamp is `".." in <normalized>.parts` — in repo-relative PurePosixPaths that IS the escape condition. Symlink escapes are a different, deliberately unclamped class (no shipped symlinks exist). Closing the HTML-slug divergence needs code-span-aware entity decoding; the pinning fixture predicts the exact GitHub slugs if anyone attempts it.

## Orientation

Now the reference contract validates bare `#anchor` links, refuses `..`-escaping targets before any filesystem probe (including the citation checker's consumer-queue absorb, which could previously stat `/etc/hostname`), and states its two remaining slug limitations with failure directions and mutation-falsifiable pins. Lives in `_dev/tests/shipped-package-reference-contract.sh`. Leaf change to the checker; map unchanged.

## Red-Green Proof

**RED prompt/case:** per instance — a same-file link to a non-existent anchor; a `..`-escaping relative target; a heading containing an HTML tag with a link written to GitHub's slug; a blockquoted heading with a link to it.
**Why RED now:** the first three pass the suite today and the fourth fails spuriously.
**GREEN when:** each case either fails for the right reason or is documented as an accepted limit, and the suite still exits 0.
**Validation:** Review findings on REQ-243; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---

## Triage

**Route: B** - Medium

**Reasoning:** Four named edges of one checker, each with its instance stated and the fix sketched; the work is judgment about close-vs-document per instance, inside a single test file.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `_dev/tests/shipped-package-reference-contract.sh` (modify) — the four checker edges closed or documented, with fixtures

**Files I will NOT touch:**
- Shipped markdown — this REQ changes the checker, not the corpus.
- The backticked-citation checker REQ-249 just added to the same file — coordinate around it, don't rework it.

**Acceptance criteria (restated from REQ):**
- [ ] Each instance is either closed, or its limitation is stated in the file with the failure direction named.
- [ ] `bash _dev/tests/maintainer-verify.sh` still exits 0 and the existing 27 corpus anchors still resolve.
- [ ] maintenance: true — ask what can be removed before adding; two of the four may close better by deleting a claim.

## Pre-Flight

**Git:** ✓ clean
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 at the branch point (wave-1 integration tip)
**Dependencies:** ✓ Go 1.26.1, ShellCheck 0.11.0, `just`, Node, Chromium all present

*Checked by work action*
