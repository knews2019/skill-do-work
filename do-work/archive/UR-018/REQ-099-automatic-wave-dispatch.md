---
id: REQ-099
title: Automatic wave dispatch — the work loop computes and dispatches the ready set
status: completed
created_at: 2026-08-04T19:44:17Z
claimed_at: 2026-08-04T21:16:51Z
completed_at: 2026-08-04T21:30:00Z
commit: 0cf9420
kb_status: pending
user_request: UR-018
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-096]
maintenance: false
related: [REQ-096, REQ-100]
batch: parallel-building
write_set: [actions/work.md, actions/work-reference.md, docs/work-guide.md, SKILL.md, _dev/tests/contract-regressions.sh]
---

# Automatic Wave Dispatch

## What

Give the work pipeline a fan-out mode where **the loop computes the wave itself** and dispatches builders without a confirmation gate. This is a deliberate contract change: today `actions/work.md:33` says the action "does not drive a fan-out wave" and `actions/work-reference.md:320` says "a human picks which REQs run together — nothing computes the set." Both sentences get rewritten, per the user's explicit choice of fully automatic set-picking over a human-confirmed set.

## Detailed Requirements

- **Wave computation:** ready = pending REQs whose `depends_on` are satisfied, unclaimed, and not `assigned_to` another session. Wave size bounded per `crew-members/background-agents.md:53` (builders per wave sized to the harness concurrency limit). `--wave N` keeps its existing meaning (dependency-depth scoping) — document how the two compose or that auto-wave supersedes it when active.
- **Dispatch:** for each REQ in the wave, follow the existing Worktree Dispatch Mode per-REQ flow unchanged (worktree per builder mandatory, run directory + briefs before any spawn, hand-back merge sequence). Silent degradation stays: no `git worktree` support → serial, no error.
- **Integration stays serial and load-bearing:** merge → qualify → test → review → changelog → archive one REQ at a time; `actions/work-reference.md:321` ("the non-interference proof is the merge, not the pick") survives unchanged and becomes the safety argument — overlapping picks are caught at merge, which is the batch philosophy. `write_set` stays display-only; the wave computation must NOT read it as a scheduling input.
- **Mode entry:** define how auto-wave is invoked (e.g., a `do-work run` fan-out flag or harness-capability trigger) — floor-first: the default single-REQ loop remains the baseline for the simplest agent; auto-wave is the advanced-harness path, consistent with the existing "Optional, advanced harnesses only" gate at `:277`.
- Update every echo of the old "human picks / nothing computes the set" claim across shipped files (Closed Enumerations Go Stale rule).

## Constraints

- No `write_set` scheduling; no computed-set gate on anything other than `depends_on`, claim state, and `assigned_to`.
- Wall-clock saving is build-phase only — keep the honest-expectations sentence (`:322`).

## Red-Green Proof

**RED prompt/case:** Ask the work action to run the queue in parallel today: `actions/work.md:33` instructs it to process one REQ at a time and provides no wave-selection or launch-before-wait path.
**Why RED now:** Fan-out is an owner-driven manual procedure by explicit design (confirmed at `work-reference.md:316`).
**GREEN when:** The rewritten instructions let an advanced harness compute a bounded ready set and dispatch builders unattended, while the floor path (serial single-REQ) is unchanged and the merge remains the non-interference proof.
**Validation:** User confirmed (ask-tool answer: "Fully automatic set-picking").

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 3, item 8).

---
*Source: approved plan, Phase 3*

## Triage

**Route: B** - Medium

**Reasoning:** The wave predicate, the bound, what must survive, and the do-not-read list are all specified. What needed discovery was which suite assertions pin the fan-out block's phrases (five of them, one directly on the sentence being rewritten), where the old "nothing computes the set" claim is echoed, and where a new `do-work run` flag has to be declared so the unrecognized-argument guard does not reject it.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided contract rewrite

*Skipped by work action*

## Exploration

**The sentence being rewritten is itself suite-pinned.** `_dev/tests/contract-regressions.sh` extracts a `fan_out_block` and asserts five phrases inside it; one is `'advisory input|never a gate'`, which lives in the very bullet whose first half ("A human picks... Nothing computes the set") this REQ inverts. The rewrite therefore had to keep both `advisory input` and `never a gate` while changing the claim around them. The other four pins — `Serial-only`, `CHANGELOG.md`, `line proximity, not meaning`, `survivable, not prevented`, `absolute main-tree path` — sit in paragraphs this REQ does not touch, but they bound where new prose can go: anything inserted must stay *inside* the block (before `## Composed Exit Summary`) or the extraction window changes.

**A new flag has three declaration sites, not one.** `actions/work.md`'s **Input** list, its unrecognized-argument strip list, and its usage string. Miss the strip list and the guard rejects the flag it just documented. `SKILL.md`'s routing row and dispatch table name the accepted arguments too, and `SKILL.md` is word-budgeted (2650), so additions there have to be tight.

**Echoes of the old claim:** `docs/work-guide.md:91`'s "you pick the set" was the only user-facing one (the phrase "owner-driven" and "nothing computes the set" appear nowhere else in shipped files — the `boundary note` hits in `actions/ai-report*.md` are an unrelated trust boundary). Also found: the fan-out heading assertion's *message* still said "one queue owner", a leftover from REQ-096's rename — fixed here since it is one line in a file this REQ edits anyway.

**The bound's source:** `crew-members/background-agents.md` guardrail 3, *Write a manifest per wave; spawn in bounded waves* — "sized to the harness concurrency limit — not one unbounded fan-out", plus "update the manifest as each wave's files land before launching the next wave", which is what makes recompute-per-wave the natural loop shape rather than an extra rule.

## Decisions

- **D-01 — mode entry is an explicit `--fan-out [N]` flag, not a harness-capability trigger.** DECIDE & STATE (reversible, low reach; the REQ offered both). The user's decision was *fully automatic set-picking* — that is about there being no confirmation gate on the **set**, not about the **mode** turning itself on. Floor-first settles the rest: an implicit trigger would mean the same `do-work run` behaves differently on two harnesses, and a reader of the simplest agent's path would have to understand concurrency to know which one they were getting. The flag keeps the default serial and makes the change visible in the command the user typed. Bare `--fan-out` defaults to the harness limit, or **two** where unknown, because a flag that silently means "one" would read as broken.
- **D-02 — `--fan-out` composes with everything, rather than being mutually exclusive with `--wave` or targeting tokens.** DECIDE & STATE. `--wave` already documents itself as selecting *which* REQs run and explicitly not *how many at once*; that sentence names the missing knob, so making the two exclusive would leave the gap it points at permanently open. The rule is one line — `--fan-out` changes how many of the selected set run at once, never which — which is why it needs no per-combination table.

## Scope

**Files I will touch:**
- `actions/work-reference.md` — the owner-driven paragraph, the human-picks bullet, and a new **Auto-wave** subsection inside Fan-Out Dispatch
- `actions/work.md` — the `:33` stance, the Input flag declaration, strip list and usage string, the `--wave` companion note, Step 1's set selection, Step 10's recompute rule, the `write_set` rule, and the orchestrator checklist
- `docs/work-guide.md` — the user-facing bullet and the `--wave` sentence
- `SKILL.md` — routing row and dispatch table arguments
- `_dev/tests/contract-regressions.sh` — one stale assertion message left over from REQ-096

**Acceptance criteria (restated from REQ):**
- [ ] Wave computation defined: pending + dependency-ready + unclaimed + not `assigned_to` another session
- [ ] Wave size bounded per `crew-members/background-agents.md`; never unbounded
- [ ] `--wave N` keeps its meaning, and how the two compose is documented
- [ ] Per-REQ dispatch flow unchanged; silent degradation to serial preserved
- [ ] Integration stays serial; the non-interference-proof sentence survives unchanged and becomes the safety argument
- [ ] `write_set` stays display-only and is NOT a scheduling input
- [ ] Mode entry defined, floor-first: default single-REQ loop unchanged
- [ ] Every echo of the old "human picks / nothing computes the set" claim updated
- [ ] Honest-expectations sentence (build-phase-only saving) kept

## Pre-Flight

- **WARN — baseline suite red before any change:** the same 8 `chmod 500`-versus-root failures inherited by every REQ in this batch.
- Working tree clean outside `do-work/` at claim time.

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified) — the owner-driven paragraph replaced by **Reached two ways, and the default is neither**, naming `--fan-out [N]`, the no-confirmation-gate change, and the floor as the reason the flag exists instead of the behavior. The human-picks bullet rewritten to cover both paths while keeping `advisory input` and `never a gate` intact, and adding *why* `write_set` cannot be a scheduling input (a field whose absence means unknown can only inform a judgment). A new **Auto-wave** subsection: the four-clause ready predicate, the explicit exclusion of `write_set` with the merge-is-the-proof argument, the bound and its recompute-per-wave loop, and the list of what the flag does *not* change.
- `actions/work.md` (modified) — `:33` rewritten from "does not drive a fan-out wave" to "one REQ at a time unless you ask it not to", with the flag, the three unchanged things, and the `write_set` non-input rule. `--fan-out [N]` declared in **Input**, added to the strip list, added to the usage string. The `--wave` note now names `--fan-out` as the other knob. Step 1 gains the set-selection paragraph (every existing filter applies unchanged — no second readiness predicate). Step 10 gains the recompute rule. The `write_set` entry under `## Rules` now says auto-wave does not read it. The orchestrator checklist gains a Step 0 argument-parse line.
- `docs/work-guide.md` (modified) — the user-facing bullet now describes `--fan-out`, and the `--wave` sentence names it as the concurrency companion.
- `SKILL.md` (modified) — routing row and dispatch table accept `--fan-out [N]`, with the composes-with-either rule stated in the row.
- `_dev/tests/contract-regressions.sh` (modified) — the Fan-Out heading assertion's message updated from "one queue owner" to "one releaser" (REQ-096 leftover).

**What was done:** Inverted the "nothing computes the set" contract into an opt-in auto-wave mode with a computed, bounded ready set and no confirmation gate, while leaving the serial floor path, the per-REQ dispatch flow, serial integration, and every `write_set` display-only pin exactly as they were.

## Testing

### The safety sentence survives unchanged

```
$ git diff -U0 actions/work-reference.md | grep -E "^[-+].*(non-interference proof|line proximity|Serial-only —|survivable, not prevented|wall-clock saving|build phase only)"
```

Empty except for the new Auto-wave paragraph's *pointer* to it. `**The non-interference proof is the merge, not the pick.**`, the `line proximity, not meaning` limit, the Serial-only list, the run-directory ceiling note, and the honest-expectations sentence are all byte-identical — which matters more here than anywhere else in this batch, because the merge being the proof is the entire reason a computed set is safe.

### `write_set` is not a scheduling input

Every pre-existing pin still holds, and the suite's `advisory input|never a gate` assertion passes against the rewritten bullet:

```
$ bash _dev/tests/contract-regressions.sh 2>&1 | grep -c '^FAIL'
8
```

The pre-existing eight, name-for-name — including all five `fan_out_block` assertions and both `write_set`-premise sweeps (`builder_count_premise_pattern` and `exclusive_session_premise_pattern` filtered to lines mentioning `write_set|overlaps`). The new prose mentions `write_set` three times and none of those lines argues its display-only status from a builder count or from the ownership model — they argue it from *absence reads as unknown*, which is the reason that holds at any count.

### The flag cannot be rejected by the guard it must pass

`--fan-out` is declared in all three places `actions/work.md` requires: the **Input** list, the unrecognized-argument strip list ("After stripping `--wave N` and `--fan-out [N]`"), and the usage string. Missing the second would have made the action document a flag and then reject it — the failure mode this check exists for.

### SKILL.md word budget

```
$ bash _dev/tests/contract-regressions.sh 2>&1 | grep -i 'budget'
```

No output — the router is still under its 2650-word budget with the two additions.

### Ripple check

```
$ grep -rn "nothing computes the set\|Nothing computes\|human picks\|you pick the set\|owner-driven" actions/ docs/ crew-members/ SKILL.md tools/queue-kanban/*.md README.md
```

No hits. The only user-facing echo (`docs/work-guide.md`'s "you pick the set") is rewritten; `actions/ai-report*.md`'s "boundary note" hits are an unrelated prompt-injection trust boundary, checked and left alone.

### Floor path unchanged

`git diff actions/work.md` touches no step's mechanics: Step 2 through Step 9 are untouched, and the only additions to Step 1 and Step 10 are paragraphs scoped to auto-wave mode that state existing filters apply unchanged. A reader who never passes `--fan-out` reads the same instructions as before, plus one paragraph telling them a flag exists.

## Lessons Learned

**What worked:**
- Reading the suite's `fan_out_block` extraction *before* writing. The window runs from `**Fan-Out Dispatch` to `^## Composed Exit Summary`, so new prose had to land inside it — and one of the five pinned phrases sits in the exact bullet being inverted. Discovering that after writing would have meant rewriting the rewrite.
- Keeping `write_set`'s exclusion argued from *absence reads as unknown* rather than from a builder count. That is the form the two suite sweeps are built to catch, and it is also the only version of the argument that stays true under a computed set.

**What didn't:**
- The first draft of the Auto-wave predicate restated dependency-readiness in its own words. That is a second definition of readiness — the exact drift the Closed-Enumerations rule warns about, and it would have diverged from Step 1's the first time either changed. Replaced with "the same predicate and the same cycle detection as the serial scan; auto-wave adds no second definition of readiness."
- Writing multi-paragraph prose through a shell heredoc broke twice on the `→` character and on an unterminated quote inside the here-document. Prose edits of this size belong in a script file invoked by path, not inlined into a shell command.

**Worth knowing:**
- A new `do-work run` flag needs **three** edits in `actions/work.md` (Input list, strip list, usage string) plus **two** in `SKILL.md` (routing row, dispatch table). The strip list is the one that bites: omit it and the unrecognized-argument guard rejects the flag the same file just documented.
- `SKILL.md` is word-budgeted at 2650 by the suite. Router additions must be phrases, not sentences.
- The fan-out section's five pinned phrases are extracted as a *block* between two markers. Anything inserted after `## Composed Exit Summary` is outside the window and will not satisfy them, however correct it reads.

## Orientation

`do-work run --fan-out [N]` now exists: the work loop computes its own ready set — pending, dependencies met, unclaimed, not earmarked elsewhere — bounds it to N, and dispatches builders with no confirmation prompt. `[MAP CHANGED]` — this inverts a stated contract ("nothing computes the set") and adds a new argument to the queue's entry point, so the pipeline now has two execution modes where it had one. What did not move is the load-bearing part: the serial floor path, the per-REQ dispatch flow, serial integration, and the merge as the non-interference proof — which is precisely what makes a computed set defensible, since a computed set claims its REQs are runnable, never that they do not overlap. Lives in `actions/work.md` (Input, Step 1, Step 10) and `actions/work-reference.md`'s Fan-Out Dispatch -> Auto-wave, with the user-facing description in `docs/work-guide.md`. `prime_files` is empty; no prime went stale.

## Qualification

**Passed** — 5 files verified, 9 acceptance criteria traced, and the claims that matter checked mechanically rather than by reading: the survival of the safety sentences by `git diff` grep, the `write_set` pins by the suite's own sweeps, the flag's three declaration sites by inspection of each, and the router budget by the suite.

- **Substantive:** the Auto-wave subsection is new contract text, not a restatement; the two rewritten paragraphs both invert their claim rather than reword it.
- **Requirements traced:** all nine criteria map to a diff hunk or a `## Testing` check.
- **No hollow prose:** every new rule names its consumer (Step 1's filters, Step 10's loop, the bound's source in `crew-members/background-agents.md`) rather than asserting a principle with nowhere to land.

## Review

**Overall: 91%** | 2026-08-04T21:30:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | Medium |
| Acceptance | Pass |

**Findings:** 0 important, 3 minor
**Acceptance:** Pass — the contract is inverted as the user chose, the bound and the exclusions are explicit, every pinned phrase and every `write_set` display-only assertion survives, and the floor path is provably untouched.
**Suggested testing:** 2 items
**Follow-ups created:** None

**Requirements checklist:** all nine `## Scope` criteria delivered; evidence per criterion in `## Testing`.

**Minor:**
- **Test Adequacy is the weakest dimension, and honestly so.** Everything verified here is *prose consistency* — pins held, echoes gone, flag declared, floor untouched. Nothing in this REQ can demonstrate that an agent reading the rewritten instructions actually dispatches a bounded concurrent wave. That is REQ-100's entire job (it must prove overlapping builder timestamps), and until it runs, auto-wave is specified and unexercised. Scored rather than filed, because the REQ that closes it is already queued and next.
- Bare `--fan-out` falling back to **two** when the harness limit is unknown is a chosen number, not a derived one. It is the smallest value that is still fan-out, which makes it the safest default, but it is a decision a reader could reasonably want justified in the flag's own line rather than inferred. Stated in both the Input entry and the Auto-wave subsection.
- Scope grew by three files beyond the REQ's `write_set` (`docs/work-guide.md`, `SKILL.md`, `_dev/tests/contract-regressions.sh`), all declared before editing. The first two are mandatory — a new flag that the router does not accept is unreachable, and the user-facing guide asserting "you pick the set" would have been left false. The third is one stale assertion *message* from REQ-096's rename, fixed in passing.

**Scope drift:** none against the declaration — five files declared, five touched.

**Restatement sweep (MUST):** run, and it is central here: this REQ redefines a contract claim ("nothing computes the set") that other text restated. Swept `actions/`, `docs/`, `crew-members/`, `SKILL.md`, `tools/queue-kanban/*.md` and `README.md` for the claim in five phrasings; the one shipped echo (`docs/work-guide.md`) is rewritten, and one suite assertion *message* carrying REQ-096's superseded wording was found and fixed. Also swept the *converse* — text that would become false if `write_set` were read by the wave — and confirmed all such pins are intact and now have an explicit "not read at all by auto-wave" statement in `actions/work.md`'s Rules to bind them to the new mode.

**Risk note (why Medium):** the contract now licenses unattended concurrent dispatch, and the safety argument rests entirely on one sentence that must keep being true — the merge is the non-interference proof. If a later change weakened `--no-ff --no-commit` integration or made a builder write the main tree, a computed set would become genuinely unsafe rather than merely unproven. Mitigated by that sentence being suite-pinned and by integration staying serial, but the residual risk is real and belongs on the record.

**Suggested additional testing:**
- REQ-100 must prove overlapping wall-clock builder timestamps; if it cannot, the honest outcome is that auto-wave is documented and unexercised, not that it works.
- A deliberate two-REQ wave with overlapping `write_set`s, to confirm the merge refuses and the computed set does *not* pretend the pick was safe — the negative case for the safety argument above.

*Reviewed by review-work action (pipeline mode, in-session)*
