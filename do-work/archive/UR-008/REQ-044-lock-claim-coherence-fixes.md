---
id: REQ-044
title: "Lock claim coherence: stale proceed-anyway restatement, Step 2 recompute ordering, write_set clear-on-recovery, self-reinforcing claim exemption"
status: completed
route: C
created_at: 2026-07-29T09:30:45Z
claimed_at: 2026-07-29T12:38:08Z
completed_at: 2026-07-29T12:49:09Z
commit: bb7f8c4
user_request: UR-008
addendum_to: REQ-035
domain: general
prime_files: []
tdd: false
depends_on: []
related: [REQ-043, REQ-045, REQ-047]
batch: deep-review-followups
write_set: [actions/work.md, actions/work-reference.md, _dev/tests/contract-regressions.sh]
maintenance: true
---

# Lock Claim Coherence Fixes (REQ-035 Follow-Up)

## What

REQ-035 made `claimed_reqs` canonical and rewrote the Crash Recovery gate to freshness-only ("including this session's own"), but four coherence defects survive around the new claim semantics.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `crew-members/general.md`, `crew-members/coding-guardrails.md`, `crew-members/maintenance.md` (`maintenance: true`). Approach recorded in `## Plan` below; restatement sweep in `## Exploration`; design call for item 4 in `## Decisions` (D-01…D-05).
- [x] **[APPLY]:** Three files changed, exactly the declared `write_set`. No new sections invented, no behavior REQ-035 established re-litigated.
- [x] **[UNIFY]:** `git diff --stat` → `_dev/tests/contract-regressions.sh | 18 ++++`, `actions/work-reference.md | 19 +++---`, `actions/work.md | 2 +-` (3 files, 32 insertions, 7 deletions). Verified each: **`actions/work-reference.md`** — proceed-anyway restatement now matches the Crash Recovery gate; recovery substep 1 carries the conditional `write_set` clear and substep 2 no longer clears it; heartbeat bullets re-sourced to the dispatch record; new *Dispatch bookkeeping* block; lock-schema `claimed_reqs` definition reconciled to the Step 2 window; `claimed_req` still scalar everywhere (`grep -c 'claimed_reqs\[0\]'` = 2, unchanged) and never array-shaped. **`actions/work.md`** — one appended clause on Step 2 substep 1 pointing at the record; no other line touched. **`_dev/tests/contract-regressions.sh`** — block extraction plus two assertions in the file's existing `assert_block_contains`/`assert_block_not_contains` idiom. Red-green on the new ratchet: reintroducing the stale sentence made both assertions FAIL (exit 1), restoring it returned `Contract regression checks passed.` (exit 0). Final run: **`Contract regression checks passed.`** No debug artifacts, no stray files, no `TODO`/`XXX` in the diff.

## Prior Implementation

REQ-035 (archived, commit `fd56267`, v0.144.0): `claimed_reqs` list canonical with `claimed_req` as derived legacy mirror; all lock-write paths updated; recovery gate rewritten at `actions/work-reference.md` (Crash Recovery, ~:205) to skip any REQ under a fresh claim including this session's own; cleanup Pass 0 reads the whole list; five ratchets added to `_dev/tests/contract-regressions.sh`.

## Detailed Requirements

1. **Fix the stale proceed-anyway restatement.** `actions/work-reference.md` option (a) "Proceed anyway" (~:377) still says the per-file gate "skips only files actively claimed by *another* live session" — pre-REQ-035 semantics contradicting the rewritten gate's "including this session's own." An agent reading its local instruction on that path strips its own co-dispatched siblings — the exact bug REQ-035 fixed. Align the phrasing, and **extend the ratchet** to pin this restatement (the existing ratchet pins only the ~:205 gate wording).
2. **Carve out Step 2 from the heartbeat recompute.** The heartbeat rule names Step 2 as a refresh touchpoint and says "recompute `claimed_reqs` = whatever's currently in `working/` under this session's own claim" — but Step 2 appends the claim *before* moving the file into `working/` (the 0.140.1 claim-before-move ordering). Applied literally at Step 2, the recompute erases the just-appended claim, reopening the unclaimed-crash-artifact window. State that at Step 2 the freshly appended id is part of the set even though its file hasn't moved yet.
3. **Stop destroying capture-seeded `write_set` on recovery.** Recovery clears `write_set` unconditionally, justified as "`## Scope` (its source) is stripped just above" — but capture-seeded sets on REQs that crashed before Step 5.5, and Route A REQs, have no Scope; the set is user/capture-authored frontmatter with no source to re-derive from, and its loss permanently serializes the REQ. Clear only when `## Scope` actually existed (the mirror case); preserve capture-seeded sets. Also move the clear from recovery's strip-sections step to the frontmatter-reset step where it belongs.
4. **Break the self-reinforcing claim exemption.** `claimed_reqs` is recomputed *from* `working/` contents and the recovery gate skips whatever is in `claimed_reqs` — so a co-dispatched builder that dies without reaching Step 8 leaves its REQ re-asserted on every heartbeat and unreclaimable for the life of the session (pre-REQ-035, the own-session case was recovered). Define how a claim leaves the set when its builder is known-dead (e.g. the orchestrator removes the id when a dispatched builder returns failed/aborted, so the recompute is bounded by dispatch bookkeeping, not raw `working/` listing).

## Constraints

- The mixed-version safety story (legacy `claimed_req` scalar readers) must survive every change — re-derive the mirror on any write path touched.
- Grep every restatement of the tokens you change (`claimed_req`, `claimed_reqs`, gate phrasing) across `actions/` before calling an item fixed — this REQ exists because one restatement was missed.

## Red-Green Proof
**RED prompt/case:** `grep -n "another live session" actions/work-reference.md` surfaces the proceed-anyway gate restatement contradicting the Crash Recovery gate's "including this session's own" wording; reading the heartbeat rule at Step 2 literally erases the claim Step 2 just appended (append happens before the file reaches `working/`).
**Why RED now:** Four claim-semantics sites were left un-updated or newly self-conflicting by REQ-035.
**GREEN when:** All gate restatements agree (and a ratchet pins the proceed-anyway one); Step 2's recompute exception is stated; recovery preserves capture-seeded `write_set` (clears only Scope-mirrored sets); a dead builder's claim has a defined in-session removal path so its REQ is recoverable without waiting out the 45-minute post-session staleness.
**Validation:** User confirmed (approved capture of the reviewed finding set)

## Full Context
See `do-work/user-requests/UR-008/input.md`.

## Triage

**Route:** C (Full Pipeline)
**Reasoning:** Four interlocking coherence defects in the orchestrator-lock / crash-recovery claim semantics — the concurrency surface REQ-035 rewrote and where prior full-adversarial passes confirmed real root causes. Touches `actions/work.md` (Step 2 / heartbeat rule), `actions/work-reference.md` (Crash Recovery gate, proceed-anyway, recovery reset/strip steps), and a ratchet in `_dev/tests/contract-regressions.sh`. Semantically coupled, self-referential, mixed-version-safety-critical → Route C.
**Complexity indicators:** 4 requirements, one requiring a new design decision (how a dead builder's claim leaves the set — item 4); mixed-version mirror invariant must survive; repo-wide restatement grep required (this REQ exists because REQ-035 missed one).
**Rigor:** Full adversarial-grade review at Step 7 (main-context, against the full diff + independent restatement sweep) — lock/recovery semantics are unforgiving of a missed restatement.

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## Plan

Reconcile-not-redesign pass across the claim-semantics restatements. Per `crew-members/maintenance.md`, each item was first tried as a *narrowing* of existing prose; only item 4 needed genuinely new contract text, and even there the fix removes a source of truth rather than adding a mechanism.

1. **Proceed-anyway restatement (item 1).** Rewrite the option (a) tail so the gate it describes is the gate that actually runs: skips *every* fresh claim — holder, other coexisting entries, and this session's own co-dispatched siblings — freshness alone gating, with a pointer to **Crash Recovery (Step 1)** as the canonical statement. Then extend the ratchet. The existing REQ-035 ratchet is file-wide `assert_contains "actions/work-reference.md" 'including this session'\''s own'`, which the Crash Recovery gate's own wording satisfies — that is *why* it never caught this. So the new pin must be **block-scoped**: extract the proceed-anyway option with `sed` (same idiom as `skill_dispatch_block`, `blocked_probe_shell_block`) and assert on the block, in both directions (current story in, stale sentence out).
2. **Step 2 carve-out (item 2)** and **4. dead-builder claim removal (item 4)** — one fix, not two. Both defects come from the same root: the recompute's *source* is wrong. "The ids currently in `working/` under this session's own claim" tries to express "mine" with a directory listing plus a qualifier a listing cannot answer, which (a) finds nothing at Step 2 (claim appended before the move) and (b) in practice sends the agent back to the lock's own `claimed_reqs` for the "mine" half, closing a self-reinforcing loop with the recovery gate. Re-source the recompute to **this session's dispatch record**, then state the record's lifecycle: enters at Step 2's append (so the Step 2 window needs no exception — the id is in the record before the file moves), leaves at Step 8's release or on a known-dead builder. Add the matching local pointer in `actions/work.md` Step 2 substep 1 so an agent reading only that step does not apply the old literal recompute.
3. **Conditional `write_set` clear (item 3).** Move the clear from recovery substep 2 (strip sections) to substep 1 (frontmatter reset) and gate it on `## Scope` actually being present. The move is what makes the gate *possible*: substep 1 runs before `## Scope` is stripped, so the mirror-vs-capture-seeded distinction is still observable there. Substep 2 keeps one back-pointer sentence so the ordering dependency is visible from the stripping side.

**Design decision for item 4** (full reasoning in D-02/D-03): the removal trigger is *builder known-dead with no Step 8 left to run* — the dispatch aborted/errored or the harness reports the agent gone. It is deliberately **not** "the build failed": an attempt that returns keeps its claim through retries, remediation, and Step 8's failure classification, which already owns removal on that path. Releasing a claim on an ordinary failure would invite another live session to strip a REQ mid-remediation — a worse bug than the one being fixed. Recovery of the released REQ reuses existing machinery: the file is still in `working/`, so the next Step 10 → Step 1 Crash Recovery pass resets and re-queues it. No new mechanism, no new field.

## Exploration

**Restatement sweep.** REQ-035 shipped with a missed restatement, so every token and behavior phrase touched here was grepped repo-wide across `actions/`, plus `crew-members/`, `docs/`, `decisions/`, `SKILL.md`, `README.md`, `tools/`, and `_dev/` to bound the surface.

| Grep | Hits | Disposition |
| --- | --- | --- |
| `claimed_req` / `claimed_reqs` | `actions/work.md`:185, 240, 560, 571, 596, 601; `actions/work-reference.md`:205, 219, 225, 269–270, 282–283, 289–290, 296, 298, 368, 370, 379, 388, 389, 398, 406, 407, 412; `actions/cleanup.md`:31; `_dev/tests/contract-regressions.sh`:466–487 | Touched: work-reference 296 (schema definition), 388 (proceed-anyway), 406–407 (recompute); work.md 240 (Step 2). All others read the field or re-derive the mirror correctly — see below. |
| `another live session` | `actions/work-reference.md`:388 only (plus the ratchet's own comment at `_dev/tests/contract-regressions.sh`:496) | **Fixed** — the single stale site, exactly the one REQ-035 missed. |
| `including this session's own` | `actions/work-reference.md`:205; ratchet :491 | Canonical gate wording; now also present in the proceed-anyway block and pinned there block-scoped. |
| `recompute` | `actions/work-reference.md`:339, 370, 406, 407; ratchet :406 | 339/370 are the *staleness age* recompute (unrelated to claims) — untouched. 406/407 re-sourced. |
| `own claim` / `per-file (concurrency) gate` | `actions/work-reference.md`:205, 298, 370, 388, 406, 407, 410; `actions/work.md`:119, 643, 654 | work.md:119/643/654 already tell the post-REQ-035 story ("this session's own co-dispatched claims", "every file *this session is allowed to touch*") — correct, untouched. |
| `write_set` | `actions/work.md`:180, 182, 184, 338, 340, 399, 710; `actions/work-reference.md`:185, 209, 550; `actions/capture-reference.md`:107, 113, 115; `actions/capture.md`:213; `actions/board.md`:92; ratchet :121–136 | Only work-reference :208/:209 (the recovery clear) touched. work.md:338 ("Route A skips Step 5.5, so a Route A REQ's `write_set` stays as captured") and `actions/capture-reference.md`:113 (capture seeds it, "a hint, never a commitment") are the two sites that *establish* item 3's premise — both already correct, and now consistent with recovery instead of contradicted by it. |
| `currently has in do-work/working` / `ids of whatever's currently in` | `actions/work-reference.md`:296, 406 pre-change | Both re-sourced; post-change grep returns nothing. |

**Sites deliberately not touched, with reasons:**

- **`actions/work-reference.md`:383–384 — the user-facing option menu copy** ("(a) Proceed anyway — … will NOT touch whatever the holder (or another coexisting session) has actively claimed"). Accurate *at the decision instant*: a joining session's entry is written with `claimed_reqs: []`, so it has no own claims to describe yet. It is a promise about what will not be touched, not a claim that own claims are exempt from the gate. Left as-is to keep consent copy short and readable; the mechanics line below it (:388) is where the gate is actually specified, and that is what was fixed and pinned.
- **`actions/cleanup.md`:31 — Pass 0's live-claim gate.** Different gate, deliberately different semantics ("a session *other than the one running this cleanup*"), and REQ-035 already updated it to read the whole `claimed_reqs` list. Re-checked against both new windows: a dead builder's REQ is non-terminal in `working/`, and a Step-2-claimed REQ is `status: claimed` in `queue/` — Pass 0 sweeps only terminal-status files, so neither becomes a sweep candidate. No change needed, and the file is outside this REQ's `write_set`.
- **`actions/work.md`:560, 571, 596, 601 — Step 8's success / archive / failure / blocked paths.** Each already removes just this id and re-derives `claimed_reqs[0]` (or `null`). They are the "Step 8 released it" event; the record's lifecycle is stated once in `actions/work-reference.md` rather than restated four times — the closed-enumeration trap this REQ family keeps hitting.
- **`actions/work-reference.md`:219, 225 (worktree sweep / dispatch mode), :412 (Release last), :370 (stale take-over).** All *readers* of `claimed_reqs`; none defines the set's source, so re-sourcing the recompute leaves them correct verbatim.
- **`tools/queue-kanban/`** — no hits for any claim token; the board never parses the lock file.

**Mixed-version safety re-verified.** `claimed_req` remains a scalar on every path: `grep -c 'claimed_reqs\[0\]'` is 2 in each of `actions/work.md` and `actions/work-reference.md` (unchanged by this REQ), the two lock-shape JSON samples still show a string, and the schema line's "*never* itself made array-shaped" clause was not touched. The one write path this REQ altered (the heartbeat recompute) kept its explicit re-derivation sentence.

## Scope

**Files I will touch:**

- `actions/work-reference.md` — proceed-anyway restatement (item 1); Crash Recovery substeps 1–2 for the conditional `write_set` clear (item 3); heartbeat recompute bullets + new *Dispatch bookkeeping is the recompute's source* block (items 2 and 4); lock-schema `claimed_reqs` definition reconciled to the Step 2 window.
- `actions/work.md` — one clause on Step 2 substep 1 pointing the local reader at the dispatch record (item 2).
- `_dev/tests/contract-regressions.sh` — block-scoped ratchet on the proceed-anyway restatement (item 1).

**Acceptance criteria (restated from `## Red-Green Proof`):**

- [x] No restatement of the per-file gate anywhere in `actions/` still says the pre-REQ-035 "another live session" story; `grep -rn "skips only files actively claimed"` over `actions/` returns nothing.
- [x] A ratchet pins the proceed-anyway restatement specifically — it fails when the stale sentence returns, and it cannot be satisfied by a match elsewhere in the file.
- [x] Step 2's recompute cannot erase the claim Step 2 just appended; the reason is stated where the recompute is defined and pointed to from Step 2 itself.
- [x] Recovery preserves a capture-seeded `write_set` and clears only a `## Scope`-mirrored one, with the clear living in the frontmatter-reset substep.
- [x] A dead builder's claim has a defined in-session removal path, so its REQ is recoverable on the next Step 10 → Step 1 pass rather than after the session's 45-minute staleness.
- [x] `bash _dev/tests/contract-regressions.sh` passes.
- [x] `claimed_req` is still a scalar mirror on every write path touched.

## Implementation Summary

- `actions/work-reference.md` (modified) — proceed-anyway gate restatement aligned to freshness-only semantics; recovery's `write_set` clear moved to the frontmatter-reset substep and gated on `## Scope` existing; heartbeat recompute re-sourced from the `working/` listing to this session's dispatch record, with a new block defining that record's lifecycle (Step 2 append → Step 8 release / known-dead-builder drop); lock-schema `claimed_reqs` definition reconciled to the Step 2 claim-before-move window.
- `actions/work.md` (modified) — Step 2 substep 1 now states that its append enters the id into the dispatch record the heartbeat recompute reads, so the refresh carrying that write cannot erase it.
- `_dev/tests/contract-regressions.sh` (modified) — added `proceed_anyway_block` extraction plus one `assert_block_contains` (current story present) and one `assert_block_not_contains` (stale sentence absent).

**What was done.** Four defects that all surfaced around REQ-035's new claim semantics turned out to be three edits, because items 2 and 4 shared a root cause: the heartbeat recompute was sourced from a `working/` directory listing, which cannot answer "whose claim is this" and finds nothing at all during Step 2's claim-before-move window. Re-sourcing it to the orchestrator's own dispatch record closes both — the freshly appended id is in the record before the file moves, and a known-dead builder's id can *leave* the record, which a listing-plus-lock-read source could never allow. The stale proceed-anyway restatement was one sentence, and the reason REQ-035's ratchet missed it was structural rather than accidental: that ratchet is file-wide, and the correct wording elsewhere in the same file satisfied it. The new pin is block-scoped, so the two restatements can no longer mask each other. Item 3 was a straight narrowing — recovery still clears a `write_set` that mirrors a `## Scope` it is about to strip, but no longer destroys a capture-seeded set that has no source to re-derive from, which had been permanently serializing every Route A REQ and every REQ that crashed before Step 5.5.

## Decisions

- **D-01 — Items 2 and 4 fixed as one re-sourcing, not two exceptions.** The REQ frames item 2 as "state the Step 2 exception" and item 4 as "define how a claim leaves the set." Both were symptoms of the recompute reading the wrong source. Adding a Step 2 exception *and* a dead-builder exception on top of a listing-based recompute would have been two caveats bolted to a broken definition (`crew-members/maintenance.md`: an action that grew caveats — fix the source, don't add the sixth clause). Re-sourcing to the dispatch record makes the Step 2 case fall out for free and makes item 4 expressible at all. Net effect: one removed source of truth, one added lifecycle statement.
- **D-02 — The dead-builder trigger is "no Step 8 left to run," not "the build failed."** Chosen deliberately. A build that returns — badly, red, or rejected at review — still has an orchestrator step that owns its release (Step 6 retries, Step 7 remediation, Step 8's failure classification). Dropping the claim there would expose a REQ to another live session's recovery *while it is being remediated*, a strictly worse failure than the one item 4 fixes. The trigger therefore names the only condition where nothing else will ever release the claim: the dispatch itself aborted/errored, or the harness reports the agent gone. **Flagged for review scrutiny** — this boundary is the one genuinely new judgment in the REQ.
- **D-03 — A released claim recovers through existing Crash Recovery, not a new path.** When the id is dropped, the REQ file is still sitting in `working/`; it simply stops being exempt from the per-file gate, so the next Step 10 → Step 1 pass resets and re-queues it. Considered and rejected: routing the dead builder straight into Step 8's failure classification. That archives the REQ with an `error_type` instead of returning it to the queue, which is a behavior change beyond this REQ's reconciling mandate and asserts a classification nobody observed. Considered and rejected: a new "orphaned claim" lock field — new state to keep in sync, for a condition the existing gate already handles once the claim is gone.
- **D-04 — Recovery's `write_set` decision lives in substep 1, gated on `## Scope`; substep 2 keeps a one-line back-pointer.** The move is load-bearing, not cosmetic: substep 1 runs before `## Scope` is stripped, so it is the only place the mirror-vs-capture-seeded distinction is still observable. The back-pointer in substep 2 exists because the stripping is what creates the ordering dependency, and an agent editing substep 2 alone would otherwise not know it had broken substep 1.
- **D-05 — The new ratchet is block-scoped and asserts in both directions.** File-wide `assert_contains` is exactly what let this defect ship: the Crash Recovery gate's correct wording satisfied the assertion while the proceed-anyway restatement contradicted it. The block extraction uses the file's existing `sed`-range idiom, and the bolded `**(a) Proceed anyway**` / `**(b) Take over**` markers were verified to match only the mechanics lines, never the un-bolded menu copy above them. `assert_block_contains` catches a rewrite that silently drops the semantics; `assert_block_not_contains` catches a literal regression to the old sentence. Red-green proven by reintroducing the stale sentence (both FAIL, exit 1) and restoring it (pass, exit 0).

## Review

**Acceptance: Pass — overall ~96%.** Main-context adversarial review against the full 3-file diff + an independent restatement sweep.

**Requirements (all 4 met):**
1. Proceed-anyway restatement now matches the Crash Recovery gate ("skips every file under a fresh claim … including this session's own"); ratchet extended with a **block-scoped** `assert_block_contains`/`assert_block_not_contains` pair — the pre-existing file-wide assertion couldn't catch this restatement (a match anywhere satisfied it). Red-green verified; contract-regressions passes.
2 & 4 (solved together via a new **dispatch-record** concept): the heartbeat recompute now sources `claimed_reqs` from the orchestrator's in-memory dispatch record — never a `working/` listing, never the lock's previous `claimed_reqs`. Step 2's append enters the id into that record before the file moves (closes the erase-the-just-made-claim window, req 2); a known-dead builder's id is dropped from the record (breaks the self-reinforcing exemption, req 4).
3. Recovery clears `write_set` only when a `## Scope` section exists (mirror case) and preserves capture-seeded / Route-A sets; the clear moved from the strip-sections step to the frontmatter-reset step, checked while `## Scope` still exists.

**Mixed-version safety:** `claimed_req` scalar mirror re-derived on every touched path, never array-shaped.
**Restatement sweep (mine):** grepped recompute / "another live session" / "skips only" / working-listing phrasing across `actions/` — every site consistent; only the unrelated age-recompute in the take-over path remains (correctly untouched).

No Important/Critical findings. No follow-ups queued.

## Lessons Learned
**What worked:** An explicit **dispatch record** as the recompute's source unified two separate defects (the Step-2 window + the self-reinforcing exemption) into one contract — cleaner than patching each "listing of `working/`" phrasing.
**What didn't:** The pre-existing file-wide `assert_contains` ratchet gave false comfort — a match anywhere in the file satisfied it, so the stale proceed-anyway wording stayed invisible. Block-scoped assertions (extract the option's block, then assert in/out) are needed to pin a specific restatement.
**Worth knowing:** "Under this session's own claim" was always in-memory orchestrator state, never derivable from a `working/` listing (a listing can't tell whose claim a file is under). A known-dead builder ≠ a failed build — the former's claim is dropped so the REQ can be re-queued; the latter keeps its claim through remediation.

## Orientation
The orchestrator-lock `claimed_reqs` recompute now has a single named source of truth — the session's in-memory **dispatch record** — replacing the ambiguous "ids in `working/` under this session's own claim." Step 2's claim-before-move window and a dead builder's claim removal are both defined against that record. Lives in `actions/work-reference.md` (Concurrent-Orchestrator Lock Guard + Crash Recovery) and `actions/work.md` Step 2; a contract-regression ratchet pins the proceed-anyway restatement. No map change — hardens the REQ-035 lock subsystem.
