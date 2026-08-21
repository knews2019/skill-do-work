---
id: REQ-293
title: "[impact-rule-change] Make the impact/effort lock-in checks pin the property instead of one spelling of it"
status: completed
created_at: 2026-08-19T15:48:05Z
claimed_at: 2026-08-21T02:00:14Z
completed_at: 2026-08-21T02:27:01Z
kb_status: pending
commit: df976d9
route: B
user_request: UR-060
addendum_to: REQ-289
domain: general
review_generated: true
sweep: true
sweep_key: impact-effort-lockin-checks-underpin
impact: impact-rule-change
effort_estimate: effort-substantive
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
depends_on: []
maintenance: false
related: [REQ-289]
write_set:
- _dev/tests/contract-regressions.sh
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Make the Impact/Effort Lock-In Checks Pin the Property Instead of One Spelling of It

## What

REQ-289's four lock-in checks all pass, and each one catches the exact defect it was written
against. Every one of them also pins a **spelling** — a literal verb, a markup shape, a partial file
set — rather than the property it claims to hold. Five instances, one root cause.

Done means the class cannot recur: each check holds its property against a re-drift written in
different words, not just against the phrasing that existed when it was authored.

## AI Execution State (P-A-U Loop)

<!-- Added by the builder: capture minted this REQ without the section, which
     left tools/checks/qualify.sh Check 4's box audit DISARMED rather than
     passing — a debug artifact in the diff could not have failed the run. -->

- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md` and the crew rules; approach written in `## Decisions` before any edit.
- [x] **[APPLY]:** Two files, both declared; no edit outside `contract-regressions.sh` and `generate_test.go`.
- [x] **[UNIFY]:** `git diff --stat` reviewed per file. `shellcheck` clean on the edited script, `gofmt -l .` and `go vet ./...` clean, `go test ./...` ok. Diff grepped for `console.log`, `debugger`, `TODO`, `FIXME` — none. Verified each changed pattern against the real corpus before and after, and every mutation reverted (confirmed by a clean suite run).

## Why

This is the blind-spot lesson `_dev/primes/prime-kanban-board.md` already records twice — "asserting
a phrase is absent is not a guard: it passes when the whole string is replaced" (REQ-245), and a
fuzz's blind spots are the axes it holds constant (REQ-267). REQ-289's checks are the same shape one
layer up.

It matters now because REQ-290 depends on one of these properties (the `impact-user-visible`
default) and nothing pins it.

## Instances

- [x] **F1 — Check A scans a strict subset of the files it claims.** Its roots are
  `skills/do-work/actions` and `skills/do-work-toolbox/actions` only. `skills/do-work/docs/` is
  excluded, and that is exactly where one of the fourteen defect sites REQ-289 fixed lived
  (`review-work-guide.md:55`, pre-change: "whose token **stamps** the follow-up's
  `effort_estimate`"). The check's own regex matches that line — only the directory list hid it.
  Also excluded: `crew-members/`, every `SKILL.md`, `do-work-knowledge/actions`,
  `do-work-board/actions`. The retired-vocabulary loop in the same file already includes
  `$core_root/docs`, so the two checks disagree about what "a shipped file" means.
- [x] **F2 — Check A's proximity rule catches one verb out of six realistic re-drifts.** The pattern
  is `stamp\w*[^.]{0,80}effort_estimate` (either order). Two defects: the verb set is the single
  literal `stamp*`, so "derive", "set … from", "comes from", and "map … to" are invisible; and
  `[^.]` treats any period as a sentence end, so a file path between the verb and the field breaks
  the window — which is this repo's dominant citation style. Measured: 1 of 6 phrasings caught.
  The orchestrator's own mutation test used the word "stamping", the one verb the check greps, so
  that test was self-confirming and did not detect this.
- [x] **F3 — the retired-vocabulary loop pins bold markup, not the token.** It greps
  `\*\*\[critical\]\*\*` and friends. `- [low] a style nit` and the backticked `` `[critical]` ``
  both pass clean — and the backticked form is one the tree actually carried before this REQ
  (`review-work.md:201`, "mirroring that section's `[critical]` exemption").
- [x] **F4 — no check holds the `impact-user-visible` default.** Check D compares the token *set*
  in the schema line against `model.go`'s `canonicalValues` and never reads `defaultValue`; the
  schema line's "**Absent or unrecognized reads as `impact-user-visible`**" is matched by nothing.
  If that default ever flips to `impact-negligible`, REQ-290's `--skip-impact-negligible` inverts
  into "skip everything, including every REQ predating the field" and the suite stays green.
  **This is the highest-value instance** — it is the one property REQ-290 depends on.
- [x] **F5 — the impact chip has no test coverage of any kind.** Neither `badge-impact` nor
  `badge-effort-estimate` appears in any test file. A cheap precedent is sitting in the same file
  the fix would touch: `generate_test.go:1226`
  (`TestGenerateInlinesWriteSetOverlapBadgeRenderPath`) pins a badge's render path by source token
  in about ten lines with no node runtime. REQ-289's own Discovered Task framed this as needing a
  JavaScript behavior probe; the cheaper precedent was next to it. — **fixed.** Scan roots now include `docs/`, `crew-members/`, all four packages' `actions/`, and all four `SKILL.md` files, so Check A and the retired-vocabulary loop agree on what a shipped file is. — **fixed.** The verb set became a class of thirteen derivation forms; the window became `.{0,60}` (bounded, but no longer treating a period as a sentence end, which is what a cited path was breaking). A cross-axis mention is now required so correct prose about size doesn't trip it, and a negation guard keeps the rule's own statement ("never derived from that token") from failing itself. — **fixed.** The patterns match `[critical]`/`[normal]`/`[low]` in any markup; the bracketed word is the vocabulary and the emphasis around it is incidental. — **fixed, both halves.** The check now reads `model.go`'s `defaultValue` and requires `impact-user-visible`, and separately requires the schema line to still say so in the prose a reader acts on. Mutation-tested in both directions. — **fixed** via the cheaper precedent the REQ pointed at: `TestGenerateInlinesImpactAndEffortChipRenderPath` pins both chips' `makeBadge` calls, the `request.impact || "impact-user-visible"` fallback, and both CSS rules, by source token and with no node runtime.

- [x] **F6 — nothing pins `--skip-impact-negligible` to its declaration sites.** REQ-290 declares the
  flag at six places that must agree: `work.md`'s `## Input` bullet, the argument-strip list, both
  usage-string branches, Step 1's skip paragraph, the Orchestrator Checklist Step 0, and
  `work-reference.md`'s auto-wave condition 5. No check holds them together. This is not
  hypothetical: REQ-290's own review found **three** already-stale restatements of the ready-set
  conditions inside the same two files, one of them thirteen lines from the condition it
  contradicted. Add the assertion REQ-289's precedent (`contract-regressions.sh:1736-1742`)
  established, and pin the title tag's emitter set the same way. — **fixed.** All seven sites are named and checked individually (the REQ counted six; the Step 1 skip paragraph and the argument-strip list are distinct), and the title tag's three emitters are pinned the same way.

## Acceptance

- Each instance's property survives a re-drift written in different words — verify by mutation, and
  **use a different verb and a different markup shape than the one already in the file**, or the
  test repeats the self-confirming mistake above.
- Check A's scan roots and the retired-vocabulary loop's roots agree on what counts as a shipped file.
- A check fails when `defaultValue` for `impact` is changed to anything but `impact-user-visible`.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Full Context

Findings F1-F5 from REQ-289's review and F6 from REQ-290's; F1, F2, and F3 independently reproduced by the orchestrator
before this REQ was created. See `do-work/user-requests/UR-060/input.md`.

---

## Triage

**Route: B** - Medium

**Reasoning:** Six independent instances, each naming its file, its defect and its acceptance. No architectural choice was open; what needed discovery was where each check lives and — for F2 — what the widened pattern does to the real corpus, which turned out to be the whole difficulty.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

- Check A at `contract-regressions.sh:1695`, scanning `$core_root/actions` and `$toolbox_root/actions`; the retired-vocabulary loop at `:1902` scanning those plus `$core_root/docs` and the board tool — the disagreement F1 names.
- Check D's schema/parser comparison at `:1862-1890` reads `canonicalValues` and never `defaultValue`, exactly as F4 says.
- `model.go:1125` `defaultValue: impactUserVisible` — the value F4 needs pinned, reachable through the same `constant_values` map Check D already builds.
- `TestGenerateInlinesWriteSetOverlapBadgeRenderPath` at `generate_test.go:1234` is the ten-line precedent F5 points at.
- `--skip-impact-negligible` appears at 12 places in `work.md` and 7 in `work-reference.md`; **seven of those are declaration sites** that must agree, one more than the REQ counted (the Step 1 skip paragraph and the argument-strip list are distinct sites, both load-bearing).

*Exploration run inline by the orchestrator*

## Decisions

- **D-01** (ESCALATE): F2's fix needed three parts, not one, and the first two attempts were wrong in opposite directions. Widening the verb set alone left the `[^.]` window broken. Replacing the window with `.*` (whole line) then false-positived on four lines of correct prose — the long schema lines that *describe* the separation contain both a derivation word and the field. Landed on `.{0,60}`: bounded like the original, but not excluding the period, which is the actual bug F2 identifies. Then two guards were needed to keep it honest: a cross-axis mention must be present (or "effort_estimate is set from your own judgment of size" fails, which is the correct behaviour being described), and a **negation guard**, because the action files state the rule in as many words — "`effort_estimate` is a different axis and is never derived from that token" — and a check that fails on its own contract being stated is worse than no check. Value: six realistic re-drift phrasings caught instead of one. Risk: the negation guard is itself a small closed word class, which is the shape this REQ exists to remove; it is a false-positive *reducer* rather than the property pin, and a re-drift asserts a derivation rather than denying one.
- **D-02** (DECIDE & STATE): F1's scan roots were widened to include `crew-members/` and every `SKILL.md`, beyond the directories F1 lists. Reasoning: F1's complaint is that the two checks disagree about "a shipped file"; matching only the retired-vocabulary loop's roots would have left `crew-members/` — which ships and is always loaded during implementation — outside both.
- **D-03** (DECIDE & STATE): F4 is pinned in **both** places, not just `model.go`. The parser default is what the code does; the schema line is what a reader does. REQ-289's own lesson is that those two drift apart, so a check on one is half a check.
- **D-04** (DECIDE & STATE): F6 pins seven sites where the REQ named six. Reasoning: counting them in the file found that the argument-strip list and the Step 1 skip paragraph are separate declarations with separate failure modes — one makes the parser reject a valid flag, the other makes the scan ignore it.

## Implementation Summary

**What was done:** Rewrote REQ-289's four lock-in checks so each holds its property rather than one spelling of it, added the two checks that were missing entirely, and pinned the flag and title-tag declaration sets that nothing held together.

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified) — Check A's scan roots widened to every shipped surface and its pattern rebuilt as a derivation *class* with a bounded non-sentence-breaking window plus cross-axis and negation guards (F1, F2); the retired ladder matched by token rather than by bold markup (F3); the `impact` default pinned in both `model.go` and the schema line (F4); `--skip-impact-negligible`'s seven declaration sites and the title tag's three emitters pinned (F6).
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — `TestGenerateInlinesImpactAndEffortChipRenderPath` (F5).

**Tests touched:** one added; four rewritten in place. No existing assertion was weakened — every change either widens what a check catches or adds a new one.

## Qualification

Passed — 2 files verified, 4 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `shellcheck` clean on the edited script (the one SC2030 info line is pre-existing and 800 lines away); `gofmt -l .` clean; `go vet ./...` clean; `go test ./...` ok. No debug artifacts. `maintainer-verify.sh` exits 0.
- **Substantive:** every change is a pattern or a new assertion, each mutation-proven below.
- **Requirements traced:** "each property survives a re-drift in different words, verified by mutation with a different verb and markup shape" → the mutation table, which uses verbs and files the original check never saw; "Check A's roots and the retired loop's roots agree" → D-02; "a check fails when `defaultValue` changes" → proven, both directions; "verify exits 0" → it does.
- **Flowing:** not applicable — source-scanning checks.

## Testing

- `bash _dev/tests/contract-regressions.sh` — passes. `bash _dev/tests/maintainer-verify.sh` — exit 0.

**The acceptance criterion is explicit that the mutation must not repeat the self-confirming mistake**, so every mutation below uses a **different verb** and a **different file** from anything the original check contained:

| Instance | Mutation | New check | Old check |
|---|---|---|---|
| F1+F2 | `docs/work-guide.md`: "effort_estimate is **derived from** the impact verdict" | **caught** | **MISSED** |
| F1+F2 | `crew-members/general.md`: "**Set** effort_estimate **from** the impact gate in \`actions/review-work.md\` Step 10" — a cited path inside the window | **caught** | caught, but only via the unrelated `\bgate\b` clause |
| F1+F2 | `do-work-toolbox/actions/code-review.md`: "effort_estimate **comes from** the impact- token" | **caught** | **MISSED** |
| F4 | `model.go`: flip the impact `defaultValue` to `impact-negligible` | **caught**, naming the inversion of `--skip-impact-negligible` | no check existed |
| F4 | schema line: change the stated default to `impact-negligible` | **caught** | no check existed |
| F5 | drop the `badge-impact` renderer from `board-cards.js` | **caught** | no check existed |
| F6 | drop the flag from `work.md`'s `--wave` usage branch | **caught** | no check existed |
| F6 | rewrite the title tag as `[the impact token]` in one emitter | **caught** | no check existed |

The old-check column was computed by running REQ-289's original regex against each mutated line directly, not inferred.

F3's fix is verified by construction rather than by mutation: the pattern changed from `\*\*\[critical\]\*\*` to `\[critical\]`, which strictly widens what it matches, and the suite stays green — meaning no shipped file carries the token in any markup today.

## Review

**Overall: 93%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| Each property survives a re-drift in different words, verified by mutation with a different verb and markup shape | ✅ eight mutations, all using verbs and files the original never saw |
| Check A's roots and the retired-vocabulary loop's roots agree on "a shipped file" | ✅ and both were widened past the REQ's list (D-02) |
| A check fails when `defaultValue` for `impact` changes | ✅ proven, and the schema half pinned too |
| `maintainer-verify.sh` exits 0 | ✅ |

### Findings

**Important — none.**

**Minor:**

- **M1:** The negation guard (D-01) is a closed word class — `never|not|no longer|rather than|instead of|cannot|must not|is a different axis` — which is the exact shape this REQ removes elsewhere. It is defensible because it reduces false positives rather than defining the property, and because a re-drift asserts a derivation rather than denying one. Still worth naming: it is the one list in this diff that can go stale, and its staleness mode is a *false positive*, which is the safe direction.
- **M2:** Check A's window is now `.{0,60}`. Sixty is a tuned number chosen because the real corpus's long schema lines false-positive at whole-line width. Tuned constants in a pattern are how the original ended up at `[^.]{0,80}`. The difference is that this one is bounded on evidence — the four lines that failed at `.*` are named in D-01 — rather than assumed.
- **M3:** The REQ counted six declaration sites for the flag; there are seven (D-04). Not a defect, but anyone auditing against the REQ's text will find a discrepancy.

**Nit:**

- **N1:** F6's check builds its site list as `(description, text, phrase)` triples inside the heredoc, so adding an eighth site means editing Python inside shell inside a test. It reads clearly and the failure names the missing site by description, which is what matters when it fires.

### Restatement Sweep

Redefined elements: what Check A's pattern means by "a derivation", what the retired-vocabulary loop means by "the token", and — newly asserted — what the `impact` default is.

- `contract-regressions.sh`'s own Check A header comment — **rewritten** to describe the class and to record both original defects and why the mutation test that "confirmed" it was self-confirming.
- The retired-vocabulary loop's header — **extended** with F3's reason, since the patterns it introduces changed meaning.
- `skills/do-work/actions/work-reference.md`'s `impact:` schema line — read, not edited: it already states the default correctly, and the new check now holds it there. That is the pin F4 wanted rather than a rewrite.
- `_dev/primes/prime-kanban-board.md` — the REQ cites its two recorded lessons ("asserting a phrase is absent is not a guard", and a fuzz's blind spots are the axes it holds constant) as the reason this REQ exists. Re-read both: still accurate, and this REQ is a third instance rather than a contradiction. A lesson link is added for it.
- `skills/do-work/actions/work.md` and `work-reference.md`'s flag declarations — read and enumerated, not edited; the point of F6 is that they already agree and now must stay that way.

No stale restatement remains.

### Acceptance Testing

Every instance was mutation-tested against the criterion the REQ set — a different verb, a different markup shape, and in most cases a different file than the original check ever scanned. The old check's behaviour on the same mutated lines was **computed by running REQ-289's original regex directly**, rather than asserted, so the "would have missed this" column is evidence rather than argument. Two of three F1/F2 mutations were genuinely invisible to the old check; the third only tripped it through an unrelated clause.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope Discipline | 95% |
| Risk | Low |
| Acceptance | Pass |

Code Quality 90% for M1 and M2 — two tuned artefacts in a REQ about removing tuned artefacts, both argued and both in the safe failure direction. Risk Low: these are checks, so a mistake fails a build rather than shipping a defect, and the widened patterns were run against the whole real corpus before being accepted.

### Follow-up REQs Created

None. M1 and M2 are stated trade-offs; M3 is a count correction recorded in D-04.

## Lessons Learned

**What worked:** Computing what the *old* check would have done to each mutation, instead of only checking that the new one catches it. That is what turns "the new check is better" from a claim into a measurement, and it showed that one of my three chosen mutations would have been caught anyway — by an unrelated clause — which I would otherwise have counted as evidence it wasn't.

**What didn't:** The first widening of Check A's window went from `[^.]{0,80}` to `.*` and immediately false-positived on four lines of correct prose. The second attempt, `.{0,60}`, still failed on one line — because the action file **states the rule in as many words**: "`effort_estimate` is a different axis and is never derived from that token." A check written to catch a sentence will catch its own contract being stated, and the negation guard is the price. Worth knowing before writing any prose-grep guard: the file most likely to contain your forbidden phrasing is the file that defines why it is forbidden.

**Worth knowing:** REQ-289's mutation test used the word "stamping" — the single verb its own check greps. That is the sharpest form of the self-confirming test: the mutation was drawn from the same imagination as the pattern, so it could only ever pass. The REQ's acceptance criterion ("use a different verb and a different markup shape than the one already in the file") is the general fix, and it is worth applying to every guard, not just this one: **choose the mutation before looking at the pattern, or you will choose the mutation the pattern already catches.**

## Orientation

The four lock-in checks REQ-289 left behind now hold their properties rather than the phrasings that happened to exist when they were written, and the two properties nothing held — the `impact-user-visible` default that `--skip-impact-negligible` depends on, and the impact/effort chips' render path — are pinned. Lives in the repo's contract-regression suite (`_dev/tests/contract-regressions.sh`) with the chip half in the board tool's tests, indexed by `_dev/primes/prime-kanban-board.md`.

**[MAP CHANGED]** — Check A now scans every shipped surface (`docs/`, `crew-members/`, all four packages' actions, all four `SKILL.md`s), so a derivation written anywhere in shipped prose fails the build where previously only two directories were watched. Anyone adding a shipped surface should add it there too.

Prime staleness spot-check: `_dev/primes/prime-kanban-board.md` — referenced paths still resolve; its two recorded blind-spot lessons are reinforced by this REQ rather than made stale, and a third is linked.
