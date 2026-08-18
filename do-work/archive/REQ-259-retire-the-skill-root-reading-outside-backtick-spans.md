---
id: REQ-259
title: Retire the skill-root citation reading at its three unbackticked sites
status: completed
completed_at: 2026-08-18T21:45:31Z
claimed_at: 2026-08-18T21:16:24Z
route: B
created_at: 2026-08-18T18:07:48Z
status_changed_at: 2026-08-18T20:59:31Z
user_request: UR-055
addendum_to: REQ-249
domain: general
review_generated: true
sweep: true
sweep_key: retired-reading-outside-backtick-spans
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- skills/do-work/SKILL.md
- skills/do-work/actions/commit.md
- skills/do-work/crew-members/security.md
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T21:17:28Z
  basis:
    - Route B
    - 3-file write set
    - 3 subsystems involved
    - 3 acceptance criteria
    - full-suite verification
---

# Retire the Skill-Root Citation Reading at Its Three Unbackticked Sites

## What

REQ-249 retired the skill-root-relative citation reading and swept every **backticked** cross-package citation to the literal form — but the retired *reading* survives at three shipped sites that are bare text, which both the sweep and the new checker are structurally blind to. One of them states the retired resolution rule as prose in the core router, now contradicting the prime and the swept corpus.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `general.md`, `coding-guardrails.md`, `communication-style.md` and `maintenance.md` from the worktree (`tdd: false`, so `testing.md` was not loaded), `_dev/primes/prime-action-files.md` in full including its `## Lessons` links, the archived REQ-249 lesson this REQ descends from, and `_dev/tests/shipped-package-reference-contract.sh`'s `backticked_citation_messages` so the sweep could be strictly wider than the checker. Plan: sweep the primitive in three passes first, fix only what the sweep confirms is in scope, re-sweep. Deletion of the `SKILL.md` sentence was weighed first per `maintenance.md` § 1 and rejected — see D-02.
- [x] **[APPLY]:** One scripted edit pass, three literal string replacements, each asserted to match exactly once in its file. No file outside the write set was opened for writing. `git diff --stat` after apply: 3 files, 3 insertions, 3 deletions.
- [x] **[UNIFY]:** `git diff --stat 662788c..HEAD` → the three write-set files, 3 insertions / 3 deletions, no others. Read the full diff line by line: three single-line text hunks, no collateral edits, no debug artifacts, no stray whitespace. `git status --porcelain` empty — no scratch file left in the worktree. Each file verified individually: `SKILL.md`'s new sentence checked against the prime's § Cross-Referencing; both corrected paths resolved by hand from their own directory and confirmed to exist, then confirmed checker-covered now that they are backticked; `crew-members/security.md` diffed against the toolbox copy to confirm the two are proper mirror images. Gate run un-piped with the exit code captured directly: **`GATE_EXIT=0`**.

## Context

Found by REQ-249's independent review (Important, gate: rule-change; reproduced by execution — paths verified non-resolving, checker verified passing over them). Not builder scope drift: the user-confirmed rule's letter covers backticked citations, and these are bare text — but the diff retired the reading, and these sites still state or use it. Created `pending-answers` per the generation-≥2 cascade stop (REQ-249 is itself review-generated).

## Instances

- [x] **`skills/do-work/SKILL.md:16`** — "When a core action names a sibling path, resolve it from the parent directory containing these four skill roots" is the retired resolution rule stated as prose in the core router; applied to the new `../../` form it computes wrong paths. Restate to match the literal rule in `_dev/primes/prime-action-files.md` § Cross-Referencing.
- [x] **`skills/do-work/actions/commit.md:17`** — bare `../do-work-toolbox/actions/inspect.md`, wrong depth from `actions/`.
- [x] **`skills/do-work/crew-members/security.md:3`** — bare `../do-work-toolbox/actions/code-review.md` in the JIT_CONTEXT comment; the checker's comment scan covers backticked spans only, so this evades.

## Requirements

- No shipped text states or uses the retired skill-root-relative reading outside REQ-249's documented exemptions (fenced template/example blocks).
- Whether the two bare paths become backticked (and thereby checkable) is builder latitude; the SKILL.md prose must agree with the prime either way.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Open Questions

- [x] REQ-249's review found the retired citation reading surviving at three unbackticked shipped sites (listed under Instances). Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

**Answered [2026-08-18]:** User approved via `do-work clarify`. An external automated review independently flagged the same three sites while this question was open, which corroborates the finding.

---

## Triage

**Route: B** - Medium

**Reasoning:** Three named instances, but the requirement is a class claim (`no shipped text states or uses the retired reading`), so the builder must sweep the corpus for further bare-text sites the backtick-scoped checker is blind to before declaring it closed.

**Planning:** Not required

---

## Exploration

All three instances reproduce exactly as the REQ states:

- `skills/do-work/SKILL.md:16` — "When a core action names a sibling path, resolve it from the parent directory containing these four skill roots. Do not search the core package for an extension action." The first sentence is the retired reading; the second is a separate, still-true instruction.
- `skills/do-work/actions/commit.md:17` — `route to ../do-work-toolbox/actions/inspect.md instead`, bare, wrong depth from `actions/` (needs `../../`).
- `skills/do-work/crew-members/security.md:3` — `../do-work-toolbox/actions/code-review.md` inside the JIT_CONTEXT comment; that same comment already spells a same-package citation in backticks, so the file mixes both forms.

The rule to agree with is `_dev/primes/prime-action-files.md` § Cross-Referencing (line 91): literal relative path from the citing file's own directory, depth per-file, fenced template/example blocks exempt, and `_dev/tests/shipped-package-reference-contract.sh` checks the backticked ones in both topologies.

*Explored inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work/SKILL.md` (modify) — restate line 16's resolution prose to match the prime's literal rule
- `skills/do-work/actions/commit.md` (modify) — line 17's bare sibling path
- `skills/do-work/crew-members/security.md` (modify) — line 3's bare sibling path in the JIT_CONTEXT comment

**Files I will NOT touch:** `_dev/primes/prime-action-files.md` (it is the rule this agrees with, not a site to change), `_dev/tests/shipped-package-reference-contract.sh` (widening the checker to bare text is its own REQ, not this one — record it as a Discovered Task if the sweep argues for it), any other package's files.

**Acceptance criteria (restated from REQ):**
- [ ] No shipped text states or uses the retired skill-root-relative reading outside REQ-249's fenced template/example exemptions — established by sweeping the corpus for the primitive, not by fixing the three listed instances
- [ ] `skills/do-work/SKILL.md`'s prose agrees with `_dev/primes/prime-action-files.md` § Cross-Referencing
- [ ] Any site found beyond the three is either fixed here or reported with its reason
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0

---

## Implementation Summary

**Files changed:**
- `skills/do-work/SKILL.md` (modified)
- `skills/do-work/actions/commit.md` (modified)
- `skills/do-work/crew-members/security.md` (modified)

**What was done:** Three one-line text edits, backed by a three-pass corpus sweep for the primitive. `SKILL.md:16` no longer states the retired skill-root-relative resolution rule — it now says a sibling path is literal and resolves from the citing file's own directory at the depth the path spells, matching `_dev/primes/prime-action-files.md` § Cross-Referencing. The two bare sibling paths in `actions/commit.md:17` and `crew-members/security.md:3` were corrected from one `../` to `../../` and backticked, which moves both into `_dev/tests/shipped-package-reference-contract.sh`'s enforcement so they cannot silently rot again. The sweep (every `(../)+<package>/` token in every file type regardless of backticks; prose stating a skill-root resolution rule with no path in it; the same reading spelled with no `../`) found nothing further in the `../`-led class — 10 non-resolving unfenced tokens before, 8 after, and all 8 survivors are `<skill-root>/`-anchored invocations or consumer-queue example paths that are correct as written.

---

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` against the merged tree (range `662788c..4081c50`), run un-piped with the exit code read directly
**Result:** ✓ `GATE_EXIT=0` — "Maintainer verification passed." This run is both Step 6.5's testing and Step 8's post-merge verification: the builder verified its own branch, this verifies the integrated result.

**Red-green validation:** no new test was written, and the Finding-Closure Ratchet is satisfied by the two other routes it allows rather than by a new one. The orchestrator verified both by execution rather than accepting the hand-back's claim:

- **Instances 2 and 3 — an existing check now fails if the fix is reverted.** Reverting `skills/do-work/actions/commit.md:17` to the one-`../` spelling and running `bash _dev/tests/shipped-package-reference-contract.sh` produced, observed: `FAIL: skills/do-work/actions/commit.md:17: backticked citation does not resolve in source and installed topology: ../do-work-toolbox/actions/inspect.md` / `shipped package reference contract: FAIL (1 broken reference(s), changelog mirror matches)`, exit 1. The probe was reverted and the tree confirmed clean. Backticking is what created this coverage — bare text is the checker's blind spot, which is why both sites survived REQ-249.
- **Instance 1 — the named finding surface is deleted.** `SKILL.md:16`'s retired resolution sentence no longer exists in shipped text; there is nothing left to regress. No mechanical check covers prose that states a rule, which is exactly D-T2's point and why that follow-up matters.

**New tests added:** none — see above. A test asserting "these three lines say X" would pin the instances, not the class, which is the failure shape this REQ exists to correct.

**Existing tests updated (cross-REQ impact):** none.

*Verified by work action*

---

## Discovered Tasks

Transcribed by the orchestrator from `do-work/runs/work-2026-08-18-211613/REQ-259-handback.md`. **A worktree builder cannot write this REQ file** — the brief routes its out-of-scope finds to the hand-back instead — so Step 8 would otherwise read an absent section and queue nothing. That gap is itself a finding of this REQ's review (Important 2) and is queued as REQ-270.

- **[normal] D-T1 — a cross-package citation spelled with no `../` at all, in three copies of `crew-members/prompt-injection.md`.** Line 3 of the core, toolbox and knowledge copies each reads `through do-work-toolbox/actions/completed-work-presentation-reference.md` inside the JIT_CONTEXT comment. Not the retired reading by spelling, but it resolves from nowhere but `skills/`, and it is invisible to the checker because it is bare text. The toolbox copy is wrong twice over: from `skills/do-work-toolbox/crew-members/` the correct citation is the same-package `../actions/...`.
- **[normal] D-T2 — the fence exemption is keyed on the fence character, not on its own rationale.** REQ-249 exempted fenced blocks because their text lands in *some other file*. That holds for a Red-Green Proof template pasted into a REQ; it does not hold for the Schema Read Contract yaml block in `skills/do-work/actions/work-reference.md`, whose `#` annotations are documentation for the agent reading that file. Four citations there (lines 130, 132, 137, 204) use the retired one-`../` reading under that shield; correct depth from `actions/` is `../../`.
- **[normal] D-T3 — widening the checker past backticked spans.** Two of this REQ's three sites existed only because the checker's condition is "backticked span" while the rule's condition is "cross-package citation". Backticking moved two sites into coverage; it did not close the hole.

All three share one root cause and consolidate into a single sweep REQ (REQ-269) rather than three follow-ups.

---

## Review

**Overall: 91%** | 2026-08-18T21:44:13Z

| Dimension | Score |
|-----------|-------|
| Requirements | 88% |
| Code Quality | 90% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- The sweep's condition is narrower than the rule's condition: the `../`-led spelling is closed, but the same skills-folder base spelled *without* `../` survives at three unfenced shipped sites, and the checker is still blind to every bare-text spelling. One root cause with the builder's D-T1/D-T2/D-T3 — the class is bounded by punctuation (a `../` lead, a backtick, a fence) instead of by the thing (a citation that must resolve from this file's directory). — gate: user-visible → REQ-269 created (sweep, consolidating all four instances)
- In worktree dispatch mode a builder's `## Discovered Tasks` never reach Step 8: `work.md` Step 8 substep 4 reads the section from the REQ file, which a worktree builder may not write, so the finds silently vanish. Reproduced on this REQ — the section was absent until the orchestrator transcribed it by hand. — gate: rule-change → REQ-270 created

**Minor findings:** 2 (report only) — the router sentence is unqualified against the two conforming-but-not-literal forms the corpus uses (`<skill-root>/`-anchored invocations, fenced consumer examples); and the archived REQ's P-A-U evidence lives in the hand-back with only a transcription here, correct for worktree mode but worth not reading as self-attested.
**Acceptance:** Pass — gate exit 0 re-run by the reviewer; the sweep's counts (190 tokens, 10→8 non-resolving unfenced, 28→26 unbackticked) reproduced exactly from an independently written script; mutation tests confirm the enforcement lock is real (wrong depth with backticks → FAIL naming both sites) and that the pre-REQ state passes (bare + wrong depth → PASS), which is the hard evidence behind finding 1.
**Suggested testing:** 3 items
**Follow-ups created:** REQ-269, REQ-270; **sweeps appended to:** None

*Reviewed by review-work action*

---

## Lessons Learned

**What worked:** Sweeping for the *primitive* in three passes — including a pass for prose that states the rule with no path in it — is the only reason `SKILL.md:16` was reachable at all; it contains no `../` token, so any token-based sweep would have missed it. Backticking the two corrected paths converted them from "fixed today" into "checked forever" at zero cost, and a mutation probe proved the lock rather than assuming it.

**What didn't:** The sweep still drew its boundary at punctuation. Three passes wide, and pass 3 *found* the no-`../` sites — then the hand-back reasoned them out of the class on the grounds that "there is no `../`, so it never claimed to be a relative path." That is a spelling test standing in for a semantic one, which is the exact substitution the builder's own Pushback names one section later. Finding the instances and then dismissing them by the marker is a distinct failure from not finding them.

**Worth knowing:** The eighth consecutive REQ in this area has been bounded by a marker (a `../` lead, a backtick, a fence) rather than by the thing being governed (a citation that must resolve from the citing file's directory). Until the checker's condition *is* the rule's condition, every fix in this area closes a spelling and leaves the class open — that is what REQ-269 exists to end. Separately: a worktree builder cannot write its own REQ file, so any Step 8 substep that reads a builder-authored section from the REQ is silently disarmed in fan-out mode (REQ-270).

## Orientation

The core router now states the cross-package resolution rule correctly, so a consumer agent reading `skills/do-work/SKILL.md` no longer derives paths one directory too high; lives in the do-work core package's citation surface, governed by `_dev/primes/prime-action-files.md` § Cross-Referencing. Two sibling citations moved from unchecked prose into `_dev/tests/shipped-package-reference-contract.sh`'s enforcement. No map change — no module, data flow, or contract was added or renamed, and the prime's referenced paths all still exist.

