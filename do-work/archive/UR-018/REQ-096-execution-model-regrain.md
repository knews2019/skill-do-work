---
id: REQ-096
title: "Execution-model re-grain: claim anywhere, one releaser; dispatch widened to any tree"
status: completed
created_at: 2026-08-04T19:44:17Z
claimed_at: 2026-08-04T20:55:49Z
completed_at: 2026-08-04T21:35:00Z
commit: 7024c4a
kb_status: promoted
kb_entry: REQ-096-execution-model-re-grain-claim-anywhere-.md
user_request: UR-018
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-094]
maintenance: false
related: [REQ-094, REQ-097, REQ-099, REQ-101]
batch: parallel-building
write_set: [actions/work-reference.md, actions/work.md, actions/cleanup.md, docs/work-guide.md, _dev/tests/contract-regressions.sh]
---

# Execution-Model Re-Grain: Claim Anywhere, One Releaser

## What

Rewrite the Execution Model contract (`actions/work-reference.md:53–61`) from "one queue owner per checkout, cross-session ownership unsupported" to the user's chosen model: **any checkout may capture and claim/build; exactly one designated releaser checkout runs the release tail** (merge integration, version bump, `CHANGELOG.md` entry, archive moves, UR closure). Widen Worktree Dispatch (`:275–341`) so a builder tree may be a spawned worktree, a user workspace, a clone, or a remote/cloud sandbox.

## Detailed Requirements

**Execution Model rewrite (`:53–61`):**
- Any checkout captures and claims; claims and captures travel between checkouts via ordinary git sync.
- One releaser per queue owns the release tail. Two releasers = unspecified. Two sessions in one working tree = unspecified (unchanged). No prevention machinery for either — repair path is `actions/forensics.md` / `actions/cleanup.md`.
- The `:57` rule survives **verbatim**: never probe for a concurrent session, never ask the user to arbitrate one.
- Cross-checkout conflicts (double claims, duplicate REQ ids from concurrent capture) are ordinary merge artifacts, fixed when the branches meet; `queue-kanban verify` (`duplicate-req-id` probe) is the cheap detector. This philosophy sentence belongs in the contract so downstream prose can cite it.

**Worktree Dispatch widening (`:275–341`):**
- Builder tree generalization: worktree, user workspace, clone, or remote sandbox all satisfy the builder definition (own tree, own branch, hands back a branch).
- Remote hand-back travels on the branch itself — the absolute-main-tree-path handback mechanism is local-only.
- A non-releaser checkout treats its synced `do-work/` snapshot as potentially stale.
- New Red Flag: a second checkout running the **release tail** is the violation to watch — claiming/building/capturing elsewhere is now in contract.
- Keep intact: merge --no-ff hand-back sequence, serial integration, "the merge is the non-interference proof" (`:321`), worktrees-outside-the-repo (`:285`), run-directory pattern.

**Ripple check:** grep shipped files for restatements of the old exclusive-session claim ("one queue owner per checkout", "cross-session ownership") — per the Closed Enumerations Go Stale rule, update every echo (`docs/work-guide.md`'s summary is a known one; REQ-101 owns the full docs pass but inline one-liners that would become *false* get fixed here).

## Constraints

- Serial-only list (`:325`) keeps every item — the release tail stays serial and single-checkout.
- No reservation vocabulary (`reserve` verb, `reserved` status) — the forbidden-token ratchet stays green.
- Do-not-build list per UR-018 Batch Constraints.

## Red-Green Proof

**RED prompt/case:** Today `actions/work-reference.md:55` declares cross-session ownership outside the contract — an agent in a second clone reading it must refuse to claim.
**Why RED now:** The exclusive-session model (0.161.0) bounds ownership to one checkout.
**GREEN when:** The rewritten section licenses capture+claim from any checkout, names the single-releaser rule, and the `:57` never-probe sentence is unchanged; no shipped file still asserts the old boundary.
**Validation:** User confirmed (ask-tool answers: "Claim anywhere, one releaser"; "collisions are fixed by the agent as needed").

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 2, items 3–4).

## Addendum (2026-08-04)

From REQ-094's review (Minor finding, folded here because this REQ owns the lines): as of REQ-094, `actions/work-reference.md:55`'s "the pipeline does not detect, coordinate, or recover a second owner" and `docs/work-guide.md:91`'s "does not coordinate a second owner" are **partly false for a reason unrelated to this REQ's own scope** — crash recovery now *detects* another checkout's live claim by its `writer:` label and reports it (`claim held by <writer>, not touched`); it still doesn't coordinate or recover one. The Execution Model rewrite must account for that already-shipped behavior rather than rediscovering it. Also reword `actions/work-reference.md`'s Step-10 template line "recovery classifies each `working/` file by name" → by name *and label* if this REQ touches that paragraph.

Second fold-in (REQ-102 discovered task): the Session Checkpoint Template's inline comment (~`actions/work-reference.md:806`) still says "any entry carrying another checkout's label is copied through verbatim" — the labeled-only scoping REQ-102 fixed at both `actions/work.md` echo sites. Reword to the canonical condition ("every entry this checkout did not write" / "carries every other one through verbatim") so the template comment stops contradicting the prose 15 lines below it.

Third fold-in (REQ-095's two-clone acceptance run, evidence-backed — see its `## Testing` findings F-06
and F-04): where a consumer commits `do-work/`, **`do-work/CHECKPOINT.md` conflicts on every concurrent
claim, including two fully disjoint ones** — two single-line appends land at the same position, so git
reports `CONFLICT (add/add)` / `AA` (or `UU` where the file already existed) while the REQ files
themselves merge cleanly. The widened-dispatch prose must state the resolution, because both one-sided
resolves lose data: taking ours drops another checkout's live claim record (the poisoning, by hand),
taking theirs makes this checkout's own crash unrecoverable. The rule is **keep every entry from both
sides** — the merge-time reading of the same condition Step 10 already carries ("every entry this
checkout did not write"). Two shapes worth naming from the same run: a double claim on one REQ is a
plain **content** conflict, never a rename conflict (both sides perform the identical
`queue/` → `working/` rename, which git resolves silently); and with byte-identical claim writes the REQ
file does not conflict at all, so the `writer:` label in `CHECKPOINT.md` is the *only* thing that
surfaces the double claim. Do not write prose predicting a rename conflict.

---
*Source: approved plan, Phase 2*

## Triage

**Route: B** - Medium

**Reasoning:** The what is fully specified — the two target regions, the verbatim survivals, the three
fold-ins, and the do-not-build list are all named in the REQ. What needed discovery was *where* the old
boundary is restated (five shipped sites plus five citations of a section heading that this rewrite makes
false) and *which* suite assertions pin the very sentences being rewritten. No architectural choice is
open; the user settled the model.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided contract rewrite

*Skipped by work action*

## Exploration

**The invariant and its pin.** `actions/work-reference.md:55` was the single statement of
`one queue owner per checkout`, and `_dev/tests/contract-regressions.sh:337` counted that literal across
`actions/` asserting exactly one. The re-grain rewrites that very sentence, so the assertion had to move
with it in the same commit — the sanctioned change, per this batch's handdown.

**A second pin nobody had flagged:** `_dev/tests/contract-regressions.sh:326` pins the section *heading*
`^## Execution Model — Exclusive Session`. Renaming the heading — which the re-grain makes necessary,
since claiming is no longer exclusive — therefore touches two assertions, not one.

**Five citations of the heading by name:** `actions/work.md:118`, `actions/work-reference.md:243`,
`:279`, `:415`, `docs/work-guide.md:91`. All resolve by anchor text, so all five dangle on a rename.

**Four more restatements that go stale but are not the invariant** (found by grepping the *condition*,
not a phrase list, per the Closed-Enumerations rule): `actions/cleanup.md:31` ("that model is about queue
*owners*"), `actions/work.md:33` ("several builders under one queue owner"), `actions/work.md:555` ("the
exclusive-session model keeps none"), `actions/work-reference.md:332` ("Fan-Out Dispatch — several
builders, one queue owner").

**Two suite *rationales* that become false without any check changing.** The reservation ratchet's
comment and one of its messages justify the ban as "cross-session ownership is unsupported" and
"cross-session REQ allocation is outside the exclusive-session product contract"
(`_dev/tests/contract-regressions.sh:279`, `:318`). Claim-anywhere makes both false, and REQ-097 lands
`assigned_to` next — a reviewer reading the old rationale would read that REQ as a ratchet violation. The
checks themselves must not loosen; only the reasoning is wrong.

**What must survive untouched:** `**Current-REQ relevance.**` including its never-probe/never-arbitrate
clause and `**Three-attempt stop.**` (both extracted by the suite via their bold markers, not the
heading — so the rename is safe for them, but their text is pinned); the Fan-Out `A human picks … Nothing
computes the set` bullet (REQ-099 owns it — not this REQ's to touch); `the non-interference proof is the
merge, not the pick`; the whole Serial-only list; worktrees-outside-the-repo; the hand-back sequence.

## Scope

**Files I will touch:**
- `actions/work-reference.md` — Execution Model section rewritten and renamed; Worktree Dispatch widened; three citation/fold-in edits; Fan-Out sub-heading term
- `actions/work.md` — Step 1's exclusive-session paragraph (becomes false), plus two stale-term one-liners
- `actions/cleanup.md` — Pass 0's dangling reference to the renamed model
- `docs/work-guide.md` — the user-facing one-queue-owner bullet (becomes false)
- `_dev/tests/contract-regressions.sh` — heading pin, invariant count, a new retirement ratchet, and two false rationales

**Acceptance criteria (restated from REQ):**
- [ ] Any checkout captures and claims; claims travel by ordinary git sync (req: Execution Model)
- [ ] One releaser per queue owns the release tail; two releasers and two sessions in one tree are both unspecified, with the forensics/cleanup repair path named (req: Execution Model)
- [ ] The `:57` never-probe / never-arbitrate rule survives **verbatim** (req: Execution Model)
- [ ] Fix-at-merge philosophy stated in the contract with `duplicate-req-id` named as the detector (req: Execution Model)
- [ ] Builder tree generalized to worktree / workspace / clone / remote sandbox (req: Dispatch)
- [ ] Remote hand-back travels on the branch; the absolute-path mechanism is marked local-only (req: Dispatch)
- [ ] A non-releaser checkout treats its synced `do-work/` snapshot as potentially stale (req: Dispatch)
- [ ] New Red Flag: a second checkout running the **release tail** (req: Dispatch)
- [ ] Hand-back sequence, serial integration, `:321` non-interference sentence, worktrees-outside-the-repo, run directory all intact (req: Dispatch)
- [ ] Ripple check: no shipped file still asserts the old boundary (req: Ripple)
- [ ] Fold-in 1: detection-by-label accounted for; Step-10 template line says by name *and* label
- [ ] Fold-in 2: the Session Checkpoint Template's inline comment states the canonical non-own condition
- [ ] Fold-in 3: the checkpoint-conflict resolution rule and the observed conflict shapes are written from REQ-095's evidence
- [ ] Serial-only list keeps every item; no reservation vocabulary; suite green with no ratchet weakened

## Pre-Flight

- **WARN — baseline suite is red before any change:** 8 FAIL lines, all from
  `_dev/tests/update-script-behavior.sh`'s `chmod 500` injection, which root ignores. Pre-existing and
  unrelated; recorded at REQ-095's Step 5.75 and unchanged since. The gate for this REQ is "still exactly
  those 8, name-for-name".
- Working tree clean outside `do-work/` at claim time; dependencies present (`go` toolchain resolves).

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified) — `## Execution Model — Exclusive Session` renamed to `## Execution Model — Claim Anywhere, One Releaser` and its single invariant paragraph replaced by four: any checkout captures and claims over ordinary git sync; **one releaser per queue** owns the release tail with two releasers and two-sessions-in-one-tree both unspecified and repaired after the fact by `actions/forensics.md` / `actions/cleanup.md`; builders are not owners; cross-checkout collisions are merge artifacts with `duplicate-req-id` and the `writer:` label as the shipped detectors. Worktree Dispatch Mode gained two paragraphs and a Red Flag — the builder-tree generalization with its three deltas (worktree-specific mechanics, branch-borne remote hand-back, stale non-releaser snapshot), the claim-conflict/checkpoint-conflict resolution rule written from REQ-095's evidence, and the release-tail Red Flag. Three further edits: two citations of the renamed section, the Session Checkpoint Template's inline scoping comment, the Step-10 "by name" line, and the Fan-Out sub-heading's stale term.
- `actions/work.md` (modified) — Step 1's `**Exclusive session.**` paragraph rewritten to `**One releaser, and this session assumes it is that one.**`; two stale-term one-liners at `:33` and `:555`.
- `actions/cleanup.md` (modified) — Pass 0's "do not argue from the exclusive-session model" now names the Execution Model's ownership rule and its release-tail framing, so the pointer resolves.
- `docs/work-guide.md` (modified) — the user-facing bullet rewritten to "Claim from anywhere; release from one place", naming the merge-conflict and duplicate-id detectors.
- `_dev/tests/contract-regressions.sh` (modified) — heading pin retargeted; invariant count switched from `one queue owner per checkout` to `one releaser per queue`; a **new** retirement ratchet asserting the superseded wording is gone from `actions/`, `docs/` and `SKILL.md`; the reservation ratchet's comment and one message rewritten so they justify the ban by the retired verb/status rather than by cross-session ownership.

**What was done:** Re-grained the ownership contract from one-owner-per-checkout to claim-anywhere/one-releaser, widened builder trees to any own-tree-own-branch shape, and folded in all three addenda. Every echo of the old boundary was found by grepping the condition rather than a phrase list, and both suite assertions that pinned the rewritten sentences moved with them in this commit. No check was loosened: one was retargeted, one was retitled, and one was added.

## Qualification

**Mechanical check:** `tools/checks/qualify.sh` FAILs on the same no-code-change shape REQ-095 hit — this REQ's five files are all prose/shell, and the script's `` - `path` `` bullet parse plus its only-`do-work/`-paths rule are written for source manifests. Verified by hand instead:

```
$ git diff --stat
 _dev/tests/contract-regressions.sh | 53 ++++++++++++++--------
 actions/cleanup.md                 |  2 +-
 actions/work-reference.md          | 33 ++++++++++----
 actions/work.md                    |  6 +--
 docs/work-guide.md                 |  2 +-
```

- **Substantive:** every one of the five files carries a semantic change, not a rename sweep — the largest is the Execution Model section, replaced paragraph-for-paragraph.
- **Requirements traced:** all fourteen Scope criteria map to a diff hunk; see the Testing section's per-criterion evidence.
- **No scope drift:** five files declared, five touched, none extra. The two files the REQ's original `write_set` did not name (`actions/cleanup.md`, `_dev/tests/contract-regressions.sh`) were added to Scope and `write_set` *before* editing, per Step 5.5's one-direction mirror.
- **P-A-U:** this REQ carries no P-A-U checkbox block, so there is nothing to audit; no debug artifacts are possible in a prose diff, and `git diff` confirms none.

## Testing

### Verbatim survivals (the highest-risk requirement)

```
$ git diff -U0 actions/work-reference.md | grep -E "^[-+].*(Current-REQ relevance|Three-attempt stop|never probe|non-interference proof|Serial-only|A human picks)"
+**Red Flag — a second checkout running the release tail.** … (**Serial-only**, below).
```

The only match is the new Red Flag's *pointer* to Serial-only. `**Current-REQ relevance.**` with its
never-probe/never-arbitrate clause, `**Three-attempt stop.**`, the Fan-Out `A human picks … Nothing
computes the set` bullet (REQ-099's to change, not this REQ's), `the non-interference proof is the merge,
not the pick`, and the whole Serial-only list are all absent from the diff — byte-identical.

### Suite: green against baseline, no ratchet weakened

```
$ bash _dev/tests/contract-regressions.sh 2>&1 | grep -c '^FAIL'
8
```

Name-for-name the same eight as the pre-existing baseline (five `mid-update failure:`, two
`dirty install:`, one roll-up) — the `chmod 500`-versus-root environment artifact, not code. **Zero new
regressions.** The three assertions that moved were each exercised deliberately:

- Heading pin: passes against the new heading; would FAIL if the section were renamed again without it.
- Invariant count: `grep -c 'one releaser per queue' actions/work-reference.md` → `1`; the assertion
  demands exactly one across `actions/`.
- New retirement ratchet: `grep -rl 'one queue owner per checkout' actions/ docs/ SKILL.md` → no hits,
  which is what it asserts. This is a **tightening** — the superseded wording can no longer creep back
  in alongside the new one, the same protection the `one active REQ, one coder context` check gives its
  predecessor.

### Ripple check (requirement: no shipped file still asserts the old boundary)

```
$ grep -rniE "exclusive.session|second owner|cross-session ownership|only .do-work. session|one queue owner" \
    actions/ docs/ crew-members/ SKILL.md tools/queue-kanban/*.md README.md
tools/queue-kanban/prime-do-kanban.md:60: - REQ-075: … the thing it was *called* ("under the exclusive-session model") …
```

One hit, deliberately left: it is a REQ-075 **lesson entry quoting the historical fingerprint verbatim**,
which the suite's own comment names as the single legitimate reason to write that phrase
(`_dev/tests/contract-regressions.sh:194`, the reason that file gets no file-level negative). Rewriting a
historical lesson to match today's model would destroy the record it exists to keep. Every other site is
fixed.

### Fold-ins

- **Fold-in 1** — the Execution Model's collision paragraph now states that recovery *detects* another
  checkout's claim by label and reports it, and that detection is all it does; the Step-10 template line
  reads `by name and \`writer:\` label`.
- **Fold-in 2** — the template's inline comment now reads `every entry this checkout did not write is
  copied through verbatim, a foreign writer label and no label at all alike`, so it states the same
  condition as the prose fifteen lines below instead of contradicting it.
- **Fold-in 3** — written from REQ-095's transcripts: plain content conflict never a rename conflict,
  byte-identical claims leaving the label as the only detector, `AA`/`UU` on `CHECKPOINT.md` for every
  concurrent claim including disjoint ones, resolve by keeping every entry from both sides, and the
  committed-`do-work/`-only caveat. The paragraph cites the archived REQ so the evidence is reachable.

### Do-not-build list

```
$ grep -rniE "reserved_for|reserved_at|status: reserved|do-work reserve|assigned_at|heartbeat|refresh interval|staleness check|liveness prob" actions/ docs/ SKILL.md
(no hits outside the ban-list definitions themselves)
```

No lock, lease, heartbeat, refresh interval, staleness check or liveness probe was added; no
`assigned_at`; no reserve verb or `reserved` status; nothing schedules on `write_set` (this REQ does not
touch it); no auto `git pull`/`push` was prescribed — the new prose describes syncs the *user* performs.

## Lessons Learned

**What worked:**
- Grepping for the *condition* rather than the phrase. `one queue owner per checkout` had one hit;
  the boundary it states had nine, across five files and two suite rationales. A phrase-list sweep would
  have shipped a renamed section still cited by five dangling anchors.
- Checking the suite for pins **before** writing, not after. The heading pin at `:326` was not in this
  REQ's brief or the batch handdown — only the invariant count was — and finding it first turned a
  surprise failure into a planned two-assertion change.

**What didn't:**
- The first instinct was to leave the section heading alone to avoid the five-citation ripple. That would
  have shipped `## Execution Model — Exclusive Session` above a paragraph beginning "Any checkout may
  capture and claim" — precisely the stale-name drift this repo's maintenance rules exist to stop. The
  ripple turned out to be five mechanical edits verifiable by one grep.
- Editing the suite's reservation *rationale* looked out of scope at first. It is not: with claim-anywhere
  in contract, "cross-session REQ allocation is outside the product contract" reads as banning exactly
  what REQ-097 is about to ship, and a reviewer following the message would reject a valid change. A
  ratchet's justification is part of the ratchet.

**Worth knowing:**
- The suite extracts `Current-REQ relevance` by its **bold markers**
  (`sed -n '/^\*\*Current-REQ relevance\./,/^\*\*Three-attempt stop\./p'`), not by the section heading —
  which is why renaming the heading was safe for it. Both bold lead-ins must keep their exact text and
  their order, or that whole block of assertions silently extracts nothing.
- `tools/queue-kanban/prime-do-kanban.md` is the one shipped file allowed to name the retired
  exclusive-session premise, because its lesson entry quotes it as history. Any future sweep of that
  phrase must exempt it, and the suite comment says so.
- `one releaser per queue` is now a counted invariant: state it once, point at it everywhere else.

## Orientation

The ownership contract now says what the user chose: **any checkout may capture, claim and build against
a shared queue; exactly one checkout runs the release tail.** Builder trees widen from spawned worktrees
to any own-tree-own-branch shape, including clones and remote sandboxes, with the local-only hand-back
mechanism marked as such. `[MAP CHANGED]` — this renames a contract concept (`Execution Model — Exclusive
Session` → `Claim Anywhere, One Releaser`) and changes its counted invariant from `one queue owner per
checkout` to `one releaser per queue`, so every downstream file cites a new anchor and the suite pins a
new phrase. Lives in `actions/work-reference.md`'s Execution Model and Worktree Dispatch Mode sections,
with the user-facing summary in `docs/work-guide.md`. `prime_files` is empty; no prime staleness check
applied.

## Review

**Overall: 92%** | 2026-08-04T21:35:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Medium |
| Acceptance | Pass |

**Findings:** 0 important, 3 minor
**Acceptance:** Pass — all fourteen Scope criteria are met, every verbatim survival is proven byte-identical by diff, the suite matches baseline name-for-name, and the ripple grep leaves one deliberate historical quote.
**Suggested testing:** 2 items
**Follow-ups created:** None

**Requirements checklist:** all fourteen `## Scope` criteria delivered — see `## Testing` for the evidence
behind each. Anchor resolution verified separately: every `Execution Model` citation across `actions/`,
`docs/`, `crew-members/` and `SKILL.md` now reads `Execution Model — Claim Anywhere, One Releaser`; no
dangling anchor remains.

**Minor:**
- `**Current-REQ relevance.**` was kept byte-identical as required, and its closing sentence still reads
  "Exclusivity is the user's guarantee, not this pipeline's check." Under the new model the guarantee is
  *single-releaser*, not exclusivity, so the word is now loose — the rule it protects (never probe, never
  arbitrate) is unaffected and the paragraph is suite-pinned, which is why it was not touched. Worth a
  one-word follow-up in a REQ that owns that paragraph; not worth breaking a verbatim requirement for.
- `actions/work.md:33`'s wave paragraph got only its stale term fixed (`one queue owner` → `one
  releaser`). REQ-099 owns that paragraph and rewrites it wholesale, so a fuller edit here would collide.
  Deliberate hand-off, recorded so it does not read as an oversight.
- The suite's `exclusive_session_premise_pattern` still matches the literal `exclusive.session`, which no
  shipped file now contains outside the allowed historical quote. The pattern is a *negative* guard for
  `model.go` / `board.js` and remains correct as written, but its name and comments now describe a retired
  model — cosmetic, and deliberately left rather than widening a diff that already touches the file.

**Scope drift:** none. Five files declared in `## Scope`, five touched, none extra. The two beyond the
REQ's original `write_set` (`actions/cleanup.md`, `_dev/tests/contract-regressions.sh`) were added to both
Scope and `write_set` before any edit, keeping Scope the source and `write_set` the mirror.

**Restatement sweep (MUST):** run, and it is the substance of this REQ rather than an addendum to it. The
diff redefines a contract token (the ownership invariant) and renames a cited section, so every restatement
was swept: nine sites across five shipped files plus two suite rationales, each either rewritten or
deliberately exempted with the reason recorded. Two of them — `actions/cleanup.md` and the reservation
ratchet's justification — were outside this REQ's original `write_set` and are exactly what the sweep is
for. The one surviving mention is a historical lesson quote the suite's own comment protects.

**Risk note (why Medium, not Low):** a renamed anchor plus a re-pointed invariant is the shape that leaves
silent dangling references. Mitigated by two mechanical checks rather than by reading: the anchor grep
above, and a **new** suite ratchet that fails if the superseded wording reappears anywhere in `actions/`,
`docs/` or `SKILL.md`. The residual risk is a citation of the old heading in a file outside those trees.

**Suggested additional testing:**
- After REQ-101's docs pass, re-run the ripple grep — that REQ adds a user-facing multi-checkout section
  and is the most likely place to re-introduce owner-per-checkout phrasing.
- Have REQ-099 confirm the Fan-Out `A human picks … Nothing computes the set` bullet is still byte-identical
  when it starts, so the two REQs' edits to that region stay attributable.

*Reviewed by review-work action (pipeline mode, in-session)*
