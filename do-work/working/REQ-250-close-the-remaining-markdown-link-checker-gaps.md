---
id: REQ-250
title: Close the remaining markdown link checker gaps
status: claimed
created_at: 2026-08-18T13:55:32Z
claimed_at: 2026-08-18T18:25:40Z
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
