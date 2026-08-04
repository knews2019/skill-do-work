---
id: REQ-035
title: Represent concurrent claims in the orchestrator lock and Crash Recovery gate
status: completed
claimed_at: 2026-07-29T06:59:04Z
completed_at: 2026-07-29T07:51:14Z
commit: fd56267
route: C
kb_status: promoted
kb_entry: REQ-035-represent-concurrent-claims-in-the-orche.md
created_at: 2026-07-28T21:57:21Z
user_request: UR-007
addendum_to: REQ-032
depends_on: [REQ-033]
related: [REQ-032, REQ-033, REQ-036]
batch: parallel-dispatch
domain: general
prime_files: []
tdd: false
review_generated: true
write_set:
  - actions/work.md
  - actions/work-reference.md
  - actions/cleanup.md
  - _dev/tests/contract-regressions.sh
maintenance: false
---

# Represent concurrent claims in the orchestrator lock and Crash Recovery gate

## What

REQ-032's parallel-dispatch gate says every concurrently dispatched REQ runs Steps 2–9 "including the orchestrator lock's `claimed_req` bookkeeping" — but `claimed_req` is a single string per session, and Crash Recovery's per-file gate only protects the one REQ it names (and only against *other* sessions). Give the lock a way to represent N concurrent claims by one orchestrator, and make the Crash Recovery gate honor them (including on the same session's own Step 10 → Step 1 loop iteration).

## Why (if provided)

Review of REQ-032 (confirmed by 2 independent adversarial verifiers per finding): with N co-dispatched REQs, Step 2 rewrites `claimed_req` to the newest claim (erasing the sibling's), Step 8 clears it to null while siblings still build, and the heartbeat rule is undefined for multiple working/ files. Downstream, Crash Recovery and cleanup's live-claim gate exempt only the named REQ — every other concurrently dispatched REQ reads as an abandoned crash artifact and gets stripped-and-re-queued mid-build. Same failure class as the 2026-07-01 incident the lock guard exists to prevent.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]** Represent N concurrent claims as a canonical `claimed_reqs` list on the holder and each coexisting entry; retain `claimed_req` as a derived legacy mirror (`claimed_reqs[0]`, or `null` when empty), never array-shaped. Delete the Crash Recovery "session other than this one" clause so freshness alone gates (skip any id in ANY fresh ≤45m claim set, this session's own included). Scope the heartbeat recompute to *this session's own* dispatched `working/` files. Clear `write_set` on recovery. Add lock-step ratchets. Hit the three plan-missed sites (work.md failure path + blocked-flip, work-reference.md acquire init) and both heartbeat branches.
- [x] **[APPLY]** Edits confined to the declared write_set: `actions/work.md`, `actions/work-reference.md`, `actions/cleanup.md`, `_dev/tests/contract-regressions.sh`.
- [x] **[UNIFY]** Ran `git diff --stat`; ran `bash _dev/tests/contract-regressions.sh` (exit 0, all existing + 4 new ratchets green, each new ratchet mutation-verified); confirmed no debug artifacts added.

## Detailed Requirements

- Choose and specify ONE representation: a list-valued claim field (e.g. `claimed_reqs`), or one `coexisting_sessions`-style entry per dispatched builder — or, if judged simpler, an explicit stated limitation in the gate (concurrent dispatch requires a harness that can track claims out-of-band; the lock protocol itself stays single-claim and the gate must say so plainly). Whichever is chosen, the gate text, the lock schema, the heartbeat-refresh rule, the Crash Recovery per-file gate, and `actions/cleanup.md`'s live-claim gate must all tell the same story.
- Crash Recovery must not re-queue a REQ the *same* session concurrently claimed (today's gate only checks sessions "other than this one").
- Crash Recovery step 2 should also clear `write_set` on a recovered REQ (it strips `## Scope`, the field's source; a sourceless mirror feeds the dispatch gate stale data — Minor finding 3 of REQ-032's review). Absent ⇒ overlaps everything ⇒ serialize is the intended safe default after recovery.
- Keep the serial default untouched; floor agents must see zero behavior change.
- Parser lock-step: if any frontmatter field changes shape, mirror in `tools/queue-kanban/model.go` same commit (none expected — the lock is not a REQ file).

## Constraints

- This is spec-coherence work on `actions/work.md` + `actions/work-reference.md` only; no new action files.
- Draft with the lock schema and Crash Recovery gate open side-by-side (REQ-032 Lessons Learned).

---

## Triage

**Route: C** - Complex

**Reasoning:** Redesigns the orchestrator lock's claim representation — the exact state machine that took REQ-018 two adversarial review rounds to get right, with four coupled read/write sites (gate text, lock schema, heartbeat rule, Crash Recovery per-file gate, cleanup's live-claim gate by reference). Spec-coherence work where a plan and a contradiction hunt have both repeatedly paid off.

**Planning:** Required

## Plan

## REQ-035 Plan — multi-claim representation in the orchestrator lock

### Design decision: **(a) list-valued claim — `claimed_reqs` on the holder and on each `coexisting_sessions` entry, with `claimed_req` retained as a derived legacy mirror (`claimed_reqs[0]`, or `null` when empty).**

**Why (a).** One session = one liveness fact = one heartbeat. `claimed_reqs` is a field wholly owned by one writer, so the serialized-mutex rule ("re-read inside the critical section, change only your own fields") is unchanged — a list doesn't complicate it, provided the value is **recomputed from `do-work/working/` at every write** rather than patched from an earlier read (that also makes any missed transition, e.g. the Step 8 failure path which never clears the claim today, self-heal at the next heartbeat).

**Rejected (b) one pseudo-entry per builder.** Tempting because existing readers already iterate `coexisting_sessions[].claimed_req`, so back-compat and the same-session blind spot would both be fixed for free by *string mismatch* — accidental correctness a future id-normalization breaks. It also lies where it matters: the interactive prompt would tell an arriving orchestrator "3 sessions coexisting" (pushing it toward take-over), cleanup would report a builder as a "live session", and the heartbeat's `holder`/`else-if entry` dispatch misroutes when one session matches four objects. Per-builder heartbeats also let one builder's liveness go stale while its orchestrator is alive — meaningless divergence.

**Rejected (c) stated limitation.** Cheapest, and the REQ permits it, but it makes UR-007's entire dispatch gate unreachable inside the skill's own protocol and leaves REQ-036/037 building on a road that ends. Cost of (a) is one field plus ~six sentences.

**Rejected: making `claimed_req` itself array-shaped.** An old reader comparing a string to an array matches *nothing* and exempts *nothing* — catastrophic. An added field is merely ignored (exempts one). Additive wins.

**Mixed-version guard (new precondition, not optional):** co-dispatch only when this session is the **only live claimant in the lock** (every other holder/coexisting entry stale >45m or absent). A co-resident older session's gate reads only `claimed_req` and would re-queue your siblings. This also matches one-orchestrator-per-queue.

### Same-session rule, stated precisely
`session_id` is generated fresh at acquisition/join, so a restarted session never reuses its predecessor's id. Therefore a claim bearing *this* session's id can only be this session's live in-flight work; a dead predecessor's claim carries a different id and is recovered once its heartbeat ages past 45m. **The "session other than this one" clause is deleted, not extended:** skip any `working/` file whose id appears in *any* fresh (≤45m) claim set — holder's or entry's, this session's or another's. Recover everything stale or unclaimed. Freshness alone is the gate; my own heartbeat is fresh by construction.

### Edit sites (every reader/writer)
`actions/work.md`: **:184** gate bullet (`claimed_req` bookkeeping → `claimed_reqs`, + the only-live-claimant precondition; **keep the strings `pairwise disjoint` and `serial-only` verbatim** — ratcheted); **:239** Step 2 claim (add id to `claimed_reqs`); **:539** post-merge para (per-batch → per-merge default under co-dispatch); **:541** Step 8.1 ("don't clear yet" → "don't remove this REQ from `claimed_reqs` yet"); **:552** Step 8.6 (remove this id after the move, re-derive `claimed_req`).

`actions/work-reference.md`: **:205** Crash Recovery gate (new rule above); **step 2** add `write_set` clearing — "`## Scope` is the field's source; delete the `write_set` frontmatter in the same pass — absent ⇒ overlaps everything ⇒ serialize, the safe post-recovery default"; **:212–219** worktree sweep (delete "does not know what this same session has in flight"; sweep now exempts own fresh claims); **:225** "real co-dispatch waits on the live-claim gate…" → "co-dispatch is representable via `claimed_reqs`; governed by Step 1's gate"; **:239** post-merge verification — **per-merge is the default whenever more than one REQ is in flight** (Step 8 archives+commits per REQ, so a batch has no revertible unit left; per-batch permitted only if every Step 8 is held until the last merge verifies) — this discharges REQ-033's Important 3; **:256/268/274** JSON examples; **:280** field prose (`claimed_reqs` canonical, `claimed_req` derived legacy for older readers); **:354/:363/:382** warn + prompt + refusal templates ("Currently building:" must list *all* claimed ids — take-over nulls the whole set, so the user must see N before authorizing); **:372** proceed-anyway append (`claimed_reqs: []`); **:373** take-over; **:390–391** heartbeat rule (the recompute-from-`working/` sentence); **:396** release-last paragraph's "every freshly claimed `claimed_req`".

`actions/cleanup.md` **:31** Pass 0 gate must gain "or any id in that session's `claimed_reqs`" — it cannot survive unchanged (it names the field). **Flag: extend the declared `write_set` with `actions/cleanup.md` at the Scope step per the mirror rule.** Pass 5 (:123) is by-reference and survives untouched.

### Ratchets (`_dev/tests/contract-regressions.sh`)
REQ-033's four pin `worktree-agent-REQ-`, `git branch -d`, `[Pp]ost-merge verification`, `worktree-agent-` — none is text this plan changes; REQ-032's two must be preserved verbatim in the rewritten gate bullet. Add three: `claimed_reqs` present in each of work.md / work-reference.md / cleanup.md (the "one story" lock-step), and the Crash Recovery gate containing "including this session's own". Mutation-verify each.

### Testing
`bash _dev/tests/contract-regressions.sh`; `go test ./...` and `model.go` untouched (the lock is not a REQ file — no parser lock-step). Cold-agent walkthrough of three cases: (1) two co-dispatched siblings, one at Step 8 → sibling still exempt; (2) same session's Step 10 → Step 1 iteration → own live claims skipped; (3) stale predecessor with two claims → both recovered, `write_set` cleared.

### Not covered
REQ-036 (Step 5.5 mirror revalidation), REQ-037 (merge point/evidence sourcing), REQ-038 (worktree name uniqueness); no filesystem enforcement of write-sets; the mixed-version hazard is *mitigated* by the precondition, not eliminated; take-over's blast radius widens from 1 to N claims (surfaced in the prompt, not otherwise changed).

**Plan validation (orchestrator, session do-work-20260728T211058Z-20017):** Design decision (a) `claimed_reqs` list + derived legacy `claimed_req` is well-argued with three rejected alternatives; the mixed-version co-dispatch precondition and the deleted-not-extended "other than this one" rule are the load-bearing pieces. ⚠ The explorer found THREE write sites the plan's edit list missed — `actions/work.md:577` (Step 8 failure path), `:582` (blocked-flip), and `actions/work-reference.md:352` (acquisition write) — the builder MUST include all three or the one-story requirement fails. `actions/cleanup.md` must join the write_set at Scope time (Pass 0 gate names the field). Per-merge-as-default under co-dispatch discharges REQ-033's dormant per-batch finding — keep that clause.

*Generated by Plan agent*

## Exploration

**Plan verified against disk. All plan line anchors are correct.** Findings:

## 1. Every `claimed_req` occurrence (shipped surface)
Grep of `actions/ crew-members/ docs/ tools/ SKILL.md _dev/ next-steps.md specs/ prompts/ interviews/` returns hits in **only three files** — `docs/`, `tools/`, `SKILL.md`, `next-steps.md`, `_dev/tests/*` have **zero**. Complete list:

| Site | Class |
|---|---|
| `actions/work.md:184` | prose (gate bullet) |
| `actions/work.md:239` | **write** (Step 2 claim) |
| `actions/work.md:541` | **write** (Step 8.1 "don't clear yet") |
| `actions/work.md:552` | **write** (Step 8.6 clear-after-move) |
| `actions/work.md:577` | **write — MISSED BY PLAN'S EDIT-SITE LIST** (Step 8 failure path: "clear `claimed_req` to `null` here too") |
| `actions/work.md:582` | **write — MISSED BY PLAN'S EDIT-SITE LIST** (blocked-flip: "clears `claimed_req` to `null`") |
| `work-reference.md:205` | read/gate (Crash Recovery, both holder + entry) |
| `work-reference.md:256, 268, 274` | example JSON (holder ×2, entry ×1) |
| `work-reference.md:280` | schema definition |
| `work-reference.md:282` | prose (serialized-mutex rationale) |
| `work-reference.md:352` | write (acquire, `claimed_req: null`) — **plan lists 354/363/372/373 but not 352** |
| `work-reference.md:354` | write + warn template (stale take-over) |
| `work-reference.md:363, 382` | prompt/refusal templates |
| `work-reference.md:372` | write (proceed-anyway append) |
| `work-reference.md:373` | write (take-over) |
| `work-reference.md:390, 391` | **write** (heartbeat, two dispatch rules) |
| `work-reference.md:396` | prose ("every freshly claimed `claimed_req`") |
| `actions/cleanup.md:31` | read/gate (Pass 0) |

**Three sites the plan's edit list omits: `work.md:577`, `work.md:582`, `work-reference.md:352`.** 577 and 582 are real writes and are exactly the "Step 8 failure path never clears" case the plan's rationale cites — they must become "remove this id from `claimed_reqs`". 352 is the acquire-time initializer and must gain `claimed_reqs: []`.

## 2–3. Crash Recovery + sweep exact text
- **:205** gate clause to delete: *"check whether a session **other than this one** currently and freshly claims it: the id matches the top-level holder's `claimed_req` (holder `session_id` ≠ this session's) … or it matches any `coexisting_sessions[].claimed_req` (entry `session_id` ≠ this session's)"* — both parenthetical `≠ this session's` clauses go.
- **:209** step 2 strip list ends *"…and `## Discovered Tasks` sections … Leave `## Open Questions` and user-authored content intact."* — `write_set` clearing inserts here (note: `## Scope` is already stripped, so the frontmatter mirror is genuinely orphaned today — the plan's fix is correct and is a latent bug independent of REQ-035).
- **:219** same-session paragraph, full text quoted above; **:214** `**Exempt it**` bullet is by-reference and survives.

## 4–5. work.md / cleanup.md
- `work.md:184` retains neither ratcheted string (`pairwise disjoint` is :180, `serial-only` is :190) — rewriting :184 is ratchet-safe, but the plan's "keep verbatim" caution should point at :180/:190.
- `work.md:707` (Red Flags) and `:745` mention only `coexisting_sessions`/heartbeat — **survive untouched**. Rationalization row 9 (`:723`) says "Check `do-work/orchestrator-lock.json` first" — **survives**. Row 10 (`:724`) and the `-D` row (`:725`) untouched.
- `cleanup.md:31` names the field twice and **cannot survive**: *"it matches the top-level holder's `claimed_req`, or any `coexisting_sessions[].claimed_req`"*, plus the tail *"though by Step 10 its `claimed_req` is already `null`"* — which is **false under co-dispatch and must also change** (the plan flags only the first clause). `cleanup.md:123` Pass 5 is by-reference ("exempt under Pass 0's live-claim gate above") — **survives untouched**.

## 6. Ratchets
`_dev/tests/contract-regressions.sh:126,131` (`pairwise disjoint`, `serial-only` → work.md), `:447` (`worktree-agent-REQ-`), `:452` (`git branch -d`), `:457` (`[Pp]ost-merge verification` → **`actions/work.md`**, not work-reference), `:462` (`worktree-agent-` → cleanup.md). **No assertion pins `claimed_req` or "in-flight REQ per session"** — nothing breaks. Note `:457` targets work.md; work.md:186 says "post-merge verification" lowercase — safe.

## 7. Worktree Dispatch sentences
`:219` (*"Do not co-dispatch several worktree builders from one session until the live-claim gate can exempt more than one in-flight REQ per session."*) and `:225` (*"real co-dispatch waits on the live-claim gate being able to exempt more than one in-flight REQ per session"*). Both become "now representable via `claimed_reqs`; governed by Step 1's gate." **The honesty sentence at :225 — "Nothing here requires concurrency — a single builder in a worktree is the same one-claim loop" — is what REQ-033's Important-1 mitigation depends on; keep it, drop only the "waits on" tail.** `:239` post-merge para currently reads *"Per-batch is the default"* → must become per-merge-when->1-in-flight.

## 8. Concerns
- **Live lock is `{session_id, started_at, heartbeat_at, claimed_req: "REQ-035", coexisting_sessions: []}`** — additive `claimed_reqs` is back-compat; the current session (`do-work-20260728T211058Z-20017`) is the sole live claimant, so the plan's own mixed-version precondition is satisfiable at first write.
- **Real conflict with :282's serialized-mutex rule:** the plan says recompute `claimed_reqs` from `do-work/working/` at every write. `working/` currently holds `baseline.json` **and** the REQ file — and under co-dispatch it holds *other sessions'* claimed files too. A naive "recompute = every REQ in `working/`" would make a holder claim a coexisting session's REQ, which is precisely the field-stealing :282 forbids. The recompute must be scoped to *files this session dispatched*, and the prose must say so explicitly; "whatever's currently in `do-work/working/` under this session's claim" (:390) already carries that qualifier and must survive the rewrite.
- `work-reference.md:391` (coexisting-entry branch) is a second recompute site with the same hazard; plan cites :390–391 as one edit but they are two distinct rules.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `actions/work.md` (modify) — the parallel-dispatch gate bullet (`claimed_req` bookkeeping → `claimed_reqs` + the only-live-claimant co-dispatch precondition; keep `pairwise disjoint` and `serial-only` verbatim), Step 2 claim (add id to `claimed_reqs`), Step 8 post-merge paragraph (per-batch → per-merge default whenever >1 REQ is in flight — discharges REQ-033's dormant per-batch finding), Step 8.1 / 8.6 / the Step 8 failure path (:577) / the blocked-flip (:582) claim-clearing sites (remove *this id* from `claimed_reqs`, re-derive `claimed_req`)
- `actions/work-reference.md` (modify) — Crash Recovery per-file gate (delete the "session other than this one" clause; skip any file whose id is in *any* fresh ≤45m claim set, including this session's own), Crash Recovery step 2 (also clear the recovered REQ's `write_set` frontmatter — its `## Scope` source is already stripped), worktree sweep same-session note, the acquire-time initializer (:352, add `claimed_reqs: []`), both JSON examples, the schema/serialized-mutex prose (`claimed_reqs` canonical, `claimed_req` derived legacy), warn/prompt/refusal templates ("Currently building:" lists all claimed ids), proceed-anyway + take-over writes, the two heartbeat dispatch rules (:390/:391 — recompute scoped to *this session's own* dispatched working/ files), the release-last prose, and the Worktree post-merge-verification default
- `actions/cleanup.md` (modify) — Pass 0 live-claim gate (:31): it names `claimed_req` twice and cannot survive unchanged — add "or any id in that session's `claimed_reqs`", and fix the false "by Step 10 its `claimed_req` is already `null`" tail (false under co-dispatch)
- `_dev/tests/contract-regressions.sh` (modify) — add ratchets pinning the "one story" lock-step: `claimed_reqs` present in each of work.md / work-reference.md / cleanup.md, and the Crash Recovery gate containing the same-session inclusion ("including this session's own" or equivalent)

**Files I will NOT touch:** `tools/queue-kanban/model.go` (the lock is not a REQ file — no frontmatter shape change, so no parser lock-step); `docs/`, `SKILL.md`, `next-steps.md` (zero `claimed_req` occurrences — verified by the exploration grep); the four REQ-033 ratchets already in the test (`worktree-agent-REQ-`, `git branch -d`, `[Pp]ost-merge verification`, `worktree-agent-`) and the two REQ-032 ratchets (`pairwise disjoint`, `serial-only`) — all must stay green

**Scope note (one-directional widening at Scope time):** the REQ's Constraints line says "work.md + work-reference.md only," written before the plan discovered `actions/cleanup.md` Pass 0 *names* the field (so it "cannot survive unchanged") and before the ratchet plan. The validated Plan supersedes: `cleanup.md` and `_dev/tests/contract-regressions.sh` join the declared set here, and `write_set` is mirrored to match. This is the permitted Scope-time widening, not scope creep — recorded so review sees it as intentional.

**Acceptance criteria (restated from REQ):**
- [ ] ONE representation chosen and specified: `claimed_reqs` list on the holder and each `coexisting_sessions` entry, with `claimed_req` retained as a derived legacy mirror (`claimed_reqs[0]`, or `null` when empty)
- [ ] The gate text, lock schema, heartbeat-refresh rule, Crash Recovery per-file gate, and `cleanup.md`'s live-claim gate all tell the same story (the "one story everywhere" requirement)
- [ ] Every `claimed_req` read/write site updated — including the three the plan's edit list missed: `work.md` Step 8 failure path, `work.md` blocked-flip, `work-reference.md` acquire-time initializer
- [ ] Crash Recovery does not re-queue a REQ the *same* session concurrently claimed (the "session other than this one" clause is deleted, not extended; freshness alone gates)
- [ ] Crash Recovery step 2 also clears `write_set` on a recovered REQ (absent ⇒ overlaps everything ⇒ serialize is the safe post-recovery default)
- [ ] Mixed-version co-dispatch precondition stated: co-dispatch only when this session is the only live claimant in the lock
- [ ] The recompute-from-`working/` rule is scoped to *files this session dispatched* (never a bare "every REQ in working/", which would steal a coexisting session's claim)
- [ ] Serial default untouched — floor agents see zero behavior change; all new machinery stays inside the optional parallel-dispatch path
- [ ] `bash _dev/tests/contract-regressions.sh` green (all existing ratchets preserved verbatim) + the new lock-step ratchets pass and mutation-verify

*Scope declared by work action (orchestrator, session do-work-20260729T065754Z-5724)*

## Decisions

- **D-01 (DECIDE & STATE):** Both lock JSON examples use single-element `claimed_reqs` lists (`["REQ-018"]`, `["REQ-042"]`) that match the existing `claimed_req` scalars, rather than showing a multi-element list. Keeps the derived-mirror relationship (`claimed_reqs[0] == claimed_req`) visually obvious in the canonical schema example; the N-claim case is carried by the surrounding prose. Reversible wording choice, no downstream reach.
- **D-02 (DECIDE & STATE):** For the `cleanup.md` Pass 0 false-tail fix I chose a minimal factual reword ("under co-dispatch its `claimed_reqs` may still list sibling REQs in flight, so do not assume that entry is empty") instead of introducing new behavioral rules about whether the running session's own siblings may be swept. Pass 0 only sweeps *terminal-status* REQs, and the task scoped this edit to removing a false assertion — adding sweep behavior would be scope creep beyond the one-story requirement.
- **D-03 (DECIDE & STATE):** The 4th new ratchet pins the literal phrase `including this session's own`, which I wrote into the Crash Recovery gate as the stable anchor (the phrase the task suggested as an example). Chosen because it is load-bearing prose — it is the same-session-inclusion correctness fix — so a future edit that silently drops it is exactly what the ratchet should catch. Reversible; the ratchet and the prose were authored together.

## Implementation Summary

**Files changed:**
- `actions/work.md` (modified)
- `actions/work-reference.md` (modified)
- `actions/cleanup.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Added a canonical list-valued `claimed_reqs` field to the orchestrator lock (holder and each `coexisting_sessions[]` entry), with `claimed_req` retained as an additive derived legacy mirror (`claimed_reqs[0]`, or `null` when empty; never made array-shaped). Every reader/writer now tells one story: the parallel-dispatch gate bullet gained the only-live-claimant mixed-version precondition; Step 2 appends the claimed id (re-deriving the mirror); Step 8 substep 6, the Step 8 failure path, and the mid-run blocked-flip each remove *only* this REQ's id after the move; the Crash Recovery per-file gate now reads each entry's `claimed_reqs` list and skips any file in any fresh (≤45m) claim set *including this session's own* (the "session other than this one" clause deleted, so a Step 10 → Step 1 loop no longer re-queues a co-dispatched sibling); Crash Recovery step 2 additionally clears `write_set`; both heartbeat recompute rules (holder + coexisting-entry) recompute `claimed_reqs` scoped to this session's own dispatched `working/` files (never a bare listing of everything in `working/`); the JSON examples, schema prose, warn/prompt/refusal templates, acquire-time initializer, and cleanup.md Pass 0 gate were all updated in lock-step. Four contract-regression ratchets pin the field's presence in work.md/work-reference.md/cleanup.md and the Crash Recovery gate's same-session-inclusion phrasing. Per-merge post-merge verification became the default whenever >1 REQ is in flight (discharging REQ-033's dormant per-batch finding). No frontmatter shape change, so `tools/queue-kanban/model.go` is untouched; the serial/floor path is behaviorally unchanged.

*Summary written by work action (orchestrator)*

## Qualification

**Passed.** 4 files verified — all present and in the working diff (`qualify.sh` OK), P-A-U boxes all `[x]`, no debug artifacts added. `scope-drift.sh` clean (4 declared = 4 touched). 9 acceptance criteria traced to specific diff changes by orchestrator read: representation chosen (schema prose + JSON examples), one-story coverage (grep confirms only the 3 intentional JSON mirror scalars remain bare), all 3 plan-missed sites hit (work.md failure path + blocked-flip, work-reference.md acquire init), "session other than this one" clause deleted, `write_set` cleared on recovery, mixed-version precondition present, heartbeat recompute scoped to this-session's-own in BOTH branches, serial/floor path behaviorally unchanged, tests green. Judgment check 6 (data flows) N/A — instruction-spec files, no data path.

*Verified by work action (orchestrator)*

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ Contract regression checks passed (exit 0) — all pre-existing ratchets + the 4 new lock-step ratchets green. Matches the Step 5.75 baseline (exit 0); no new regressions.

**Red-green validation:** N/A — non-behavioral change to Markdown instruction files + a shell test. Proof is regression evidence (baseline green → still green) plus the 4 new ratchets, each mutation-verified by the builder (stripping the pinned token/phrase from its file makes the assertion fail).

**New tests added:**
- `_dev/tests/contract-regressions.sh`: `claimed_reqs` presence in each of `actions/work.md`, `actions/work-reference.md`, `actions/cleanup.md` (the "one story" lock-step), and the Crash Recovery gate's `including this session's own` phrasing. A 5th ratchet was added during remediation (below): work.md's Crash Recovery *summary* must also carry `this session's own co-dispatched claims`.

*Verified by work action*

## Review

**Overall: 86%** (as reviewed, pre-remediation) | 2026-07-29T07:44:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 85% |
| Code Quality | 83% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 3 important, 2 minor (0 nit). **All 3 Important survived full adversarial verification** (reviewer + independent contradiction-hunter, then 2 lens-diverse refuters each, 0 refuted). The 3 Important collapse to **2 distinct defects** (the :207 defect was independently found by both the reviewer and the hunter).
**Acceptance:** Pass — `contract-regressions.sh` green, all pre-existing ratchets survive, 4 new ratchets pass. Instruction-spec change; proof is regression + mutation-verified ratchets.
**Suggested testing:** none beyond the regression suite (no runnable co-dispatch path — instruction prose).
**Follow-ups created:** None — every confirmed Important finding and both Minors were **remediated inline** (see `## Remediation`), because each is incomplete execution of *this REQ's own* "one story everywhere" acceptance criterion in files already in its write_set, not separable new work.

*Reviewed by review-work action (pipeline mode, full adversarial rigor per the session's calibration)*

## Remediation

The adversarial review confirmed that REQ-035's plan built its edit-site list by grepping the literal token `claimed_req`, so two restatements of the Crash Recovery gate that phrase it *without* that token ("another live session") were missed — leaving the skill's own instructions self-contradictory on the exact behavior REQ-035 exists to fix. Because these are coherence gaps in the change's own declared criterion (#2, "all tell the same story"), in files already in the declared write_set, and uncommitted, they were fixed in place rather than deferred to a follow-up REQ (a follow-up would only have said "finish REQ-035's coherence").

**Confirmed Important — fixed:**
- **`actions/work.md:119`** (compact Crash Recovery summary) — said the gate "skips any file another live session still actively claims … re-queue everything else" (pre-REQ-035 semantics). Rewritten to "skips any file still freshly claimed by a live session — another session's, or **this session's own co-dispatched claims** on a Step 10 → Step 1 loop … freshness alone gates." This is the primary action file an orchestrator reads; leaving it stale would have re-introduced the re-queue-own-sibling-mid-build failure. **Newly ratcheted** (5th assertion) so it can't silently diverge from the work-reference gate again — the review's root cause was that neither stale site was pinned.
- **`actions/work-reference.md:207`** (recovery loop lead-in, 2 lines below the authoritative gate) — glossed "passes the gate" as "not actively claimed by another live session." Rewritten to "not freshly claimed by any live session — another session's *or this session's own*, per the gate above."

**Minors — folded in (same "one story" criterion, same files, cheap):**
- **`actions/work-reference.md:205`** — the gate said "never gate on `claimed_req` alone" with no fallback, while `cleanup.md:31` was updated to fall back to the legacy scalar for an older single-claim lock. Added the matching older-single-claim-lock fallback so the two gates tell the same mixed-version story (this was the whole point of leaving the derived mirror).
- **`actions/cleanup.md:31`** — Pass 0 named the own-in-flight-siblings hazard but no longer stated why Pass 0 is safe against them. Added the reachability reason: Pass 0 sweeps only terminal-status files, and Step 8 moves each REQ out of `working/` before Step 10 runs cleanup, so no own sibling is terminal-in-`working/` at cleanup time.

**Re-verification:** `bash _dev/tests/contract-regressions.sh` exit 0 (all existing + 5 ratchets); `qualify.sh` and `scope-drift.sh` clean (still exactly the 4 declared files — no scope change); the 5th ratchet mutation-verified (removing the pinned phrase → suite exit 1). Residual audit: every remaining "another session" occurrence across the three files is contextually correct (mutex-reclaim, acquisition-race, `reserved`-status, and the once-only proceed-anyway branch the review's refuters judged defensible). Acceptance criterion #2 now fully delivered.

*Remediated by work action (orchestrator)*

## Lessons Learned

**What worked:** The `claimed_reqs`-list + derived-`claimed_req`-mirror design (additive back-compat, one owning writer per field, recompute-scoped-to-own-claims) handled every reader/writer cleanly and kept the serial/floor path behaviorally identical. The adversarial contradiction-hunter earned its cost here: it found two stale gate restatements that the reviewer's requirements-walk also caught but that the plan's edit-site list structurally could not.

**What didn't:** The plan discovered its edit sites by grepping the literal token `claimed_req`. That grep **cannot** find the places that state the same gate rule *in prose without the token* — "skips any file another live session still actively claims" (work.md:119) and "not actively claimed by another live session" (work-reference.md:207). A token-driven edit list will always miss the behavior-phrased restatements of the same rule. This is the recurring failure mode for "one story everywhere" coherence edits, and it's the second time this class of miss has surfaced in UR-007 (the plan-missed *write* sites were caught by the explorer; these prose restatements slipped past both plan and explorer and needed the review).

**Worth knowing:** When a rule lives in prose in N places, pin BOTH the token presence AND the behavior phrasing in the primary files — a token-only ratchet (`claimed_reqs` present) stays green while a behavior restatement silently tells the old story. The fix added a phrasing ratchet (`this session's own co-dispatched claims` in work.md) alongside the token ratchets. Ratchets are still file-granular presence checks, so they pin that the right phrasing exists *somewhere* in the file, not that a specific line is correct — a second restatement in the same file could still drift; the residual is documented, not eliminated.

## Orientation

The orchestrator lock now represents **N concurrent claims per session** — `claimed_reqs` is the canonical list on the holder and each coexisting entry, with `claimed_req` retained as a derived legacy mirror (`claimed_reqs[0]`/`null`) so older single-claim readers keep working. Crash Recovery and cleanup's Pass 0 gate now exempt every id in any fresh (≤45m) claim set *including the running session's own*, so a Step 10 → Step 1 loop no longer strips and re-queues a co-dispatched sibling mid-build. Lives in the work pipeline's lock guard + Crash Recovery: `actions/work.md` Step 1/2/8 (+ failure and blocked-flip off-ramps), `actions/work-reference.md` (Concurrent-Orchestrator Lock Guard, Crash Recovery, heartbeat rules), and `actions/cleanup.md` Pass 0.

**[MAP CHANGED]** The orchestrator lock's claim field is now list-shaped (`claimed_reqs`), an additive schema change on a runtime state file (not a REQ frontmatter field, so no `tools/queue-kanban/model.go` parser lock-step). This unblocks the REQ-032 parallel-dispatch gate inside the skill's own protocol — real co-dispatch is now representable, gated by the only-live-claimant mixed-version precondition. Follow-ons REQ-036 (Step 5.5 mirror re-validation) and REQ-037 (worktree merge placement) build directly on this.
