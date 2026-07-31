---
id: REQ-033
title: Worktree dispatch mode — orchestrator-managed builders, integration, and cleanup
status: completed
commit: 849b2a5
completed_at: 2026-07-28T22:40:47Z
claimed_at: 2026-07-28T21:59:03Z
created_at: 2026-07-28T20:52:05Z
user_request: UR-007
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-032]
related: [REQ-032, REQ-034]
batch: parallel-dispatch
write_set:
  - actions/work-reference.md
  - actions/work.md
  - actions/cleanup.md
  - docs/cleanup-guide.md
  - _dev/tests/contract-regressions.sh
maintenance: false
---

# Worktree dispatch mode — orchestrator-managed builders, integration, and cleanup

## What

Document an advanced-harness dispatch mode in the work pipeline: parallel builders each run in an orchestrator-created git worktree on their own branch, the orchestrator is the sole writer to the main tree and merges builder branches in dependency order, and worktree lifecycle (creation, naming, cleanup, crash recovery) has defined owners. Canonizes the pattern a consumer runner already demonstrated (`worktree-agent-*` branches, "integrated by orchestrator" merges).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Spec lives in `actions/work-reference.md` as a new `## Worktree Dispatch Mode (Step 1)` section between Crash Recovery and the Lock Guard (its two callers). Crash Recovery gains a worktree sweep placed *after* numbered step 3 but *outside* the per-REQ-file loop — the sweep is keyed on `worktree-agent-*` leftovers, which can outlive their `working/` file, so nesting it in the loop would both misapply per-file and miss orphans (D-01). `actions/work.md` gets hooks only (~230 words): a Step 1 bullet inside the Parallel-dispatch subsection (after line 186), a Step 6 builder-brief bullet (after the `write_set` bullet), a `**Post-merge verification**` paragraph under Step 8's `**On success:**`, a new substep 8 for happy-path cleanup (after substep 7 — substeps are not named entry points), one Rules bullet, one Common Rationalizations row. `actions/cleanup.md` gets `### Pass 5: Orphaned Worktrees (consent-gated)` before the Repoint section, plus every pass-count / "one place cleanup deletes" / reporting / staging / checklist ripple. `docs/cleanup-guide.md` moves in lock-step. Four `assert_contains` ratchets pin the name convention, `git branch -d`, post-merge verification, and cleanup's ownership. Latitude exercised per the plan's recommendations and logged as D-02..D-04 (merge --no-ff never rebase; per-batch verification default; worktrees outside the repo tree). Live-claim gates are referenced, never restated (REQ-035's surface); Step 5.5's mirror text untouched (REQ-036's surface).
- [x] **[APPLY]:** Code written exactly as planned. Scope strictly limited to the five declared files.
- [x] **[UNIFY]:** `git diff --stat` reviewed; every changed file re-read; contract-regression suite green with all four new ratchets mutation-tested; `go test ./...` green. See `## Testing` notes in the Implementation Summary hand-back.

## Why (if provided)

Worktree isolation replaces mutual exclusion entirely (making timed file locks moot), keeps one claim trail and one integrator, and parallelizes beneath the single loop rather than beside it — but without defined ownership of integration, verification, and cleanup, worktrees leak (observed: 37-minute-old orphan branches) and merged-but-untested main states can be archived as done.

## Detailed Requirements

- **Naming convention:** worktree and branch names embed the REQ id — `worktree-agent-REQ-NNN-<suffix>` — so recovery can correlate any leftover with its REQ (same move the orchestrator lock makes with `claimed_req`).
- **Sole integrator:** builders never write the main tree or its branch; the orchestrator merges builder branches in dependency order. Shared files needing one-line wiring (e.g., a `<link>` insertion) are integration seams: builders hand back the line, the orchestrator applies it.
- **Post-merge verification gate:** each builder's checks validated its own branch; nobody has tested the merged state. After integration, re-run the REQ's acceptance/verification checks on the merged main tree **before** the archive step (per merge, or once per batch with the batch as the rollback unit). Four individually-green builders must not compose into a red main that archives as done.
- **Cleanup ownership — happy path:** the archive step runs `git worktree remove <path>` and `git branch -d worktree-agent-REQ-NNN-…`. Use `-d`, not `-D`: its refusal on an unmerged branch is a free assertion the integration actually happened; a `-d` failure means a skipped or lost merge — stop and investigate, don't force.
- **Cleanup ownership — crash path:** Crash Recovery (Step 1) additionally sweeps `worktree-agent-*` leftovers: branches already merged into main are removed mechanically (worktree + branch + `git worktree prune`); **unmerged** ones are reported, never auto-deleted — deletion of unmerged work belongs to `do-work cleanup` behind its existing consent gate. The sweep must skip worktrees whose REQ is freshly claimed by another live session (heartbeat-fresh, same per-file gate as the `working/` sweep).
- **State stays home:** `do-work/` (queue, working, lock, checkpoint) lives in the main tree only; builders receive their brief in the dispatch prompt and never carry their own copy of queue state (untracked files don't propagate into worktrees anyway — make the constraint explicit).
- **Preconditions + floor:** the mode requires git worktree support and a harness with parallel subagents; precondition-check and degrade gracefully to the serial loop (same pattern as `actions/board.md`'s Go check). Generalized language — no tool-specific APIs.

## Constraints

- Scheduling still consults REQ-032's write-sets: disjoint sets merge trivially; heavy overlap → serialize or pre-partition. Worktrees demote write-set violations from corruption to merge cost — they do not repeal REQ-032's serial-only resource classes (semantically colliding migrations merge textually clean).
- No new action file; lands in `actions/work.md` + `actions/work-reference.md` (and `actions/cleanup.md` for the unmerged-orphan pass).
- Do not suggest user-managed parallel sessions as part of this mode — this is orchestrator-managed parallelism beneath one loop.

## Dependencies

Depends on REQ-032 (write-set vocabulary used by the dispatch decision).

## Builder Guidance

Certainty: Firm on ownership rules, naming, `-d`-as-assertion, post-merge gate, and the consent-gated unmerged deletion; latitude on merge mechanics (merge vs rebase-then-ff) and per-merge vs per-batch verification granularity. `- [~]` Where worktree directories live (repo-sibling dir vs system temp) → deferred to builder; recommended outside the repo tree so they can't pollute status/stray checks.

## Red-Green Proof

**RED prompt/case:** `grep -i "worktree" actions/work.md actions/work-reference.md actions/cleanup.md` finds no operative dispatch mode — no naming convention, no post-merge verification requirement, no cleanup owner; a crashed parallel run leaves `worktree-agent-*` branches nothing in the skill will ever remove.
**Why RED now:** The pattern exists only as one runner's improvisation; the skill neither enables nor bounds it.
**GREEN when:** The work action documents the mode end-to-end (dispatch → integrate → verify merged state → archive-time cleanup), Crash Recovery + cleanup cover orphans with the merged/unmerged distinction, and the naming convention is stated.
**Validation:** User confirmed (design agreed in conversation, including cleanup ownership Q&A; "run the capture").

## Full Context

See `do-work/user-requests/UR-007/input.md` for complete verbatim input.

---
*Source: "so we also have worktrees as well … who will delete the merged worktrees?" (UR-007)*

Think carefully before answering.

---

## Triage

**Route: C** - Complex

**Reasoning:** End-to-end dispatch-mode specification across `actions/work.md`, `actions/work-reference.md` (Crash Recovery), and `actions/cleanup.md` (consent-gated unmerged deletion), with ownership rules that must stay coherent with the orchestrator lock, the per-file recovery gate, and REQ-032's fresh dispatch gate. Exactly the class of state-machine prose where REQ-032's review found confirmed contradictions — plan first.

**Planning:** Required

## Plan

## REQ-033 Implementation Plan — Worktree Dispatch Mode

### 1. Where the spec lives

**New `actions/work-reference.md` section `## Worktree Dispatch Mode (Step 1)`, inserted at line 211** (after Crash Recovery, before the Lock Guard — it sits between its two callers). Full spec ~500 words lives there; `actions/work.md` (13,485 words, loaded every run) gains **≤200 words** of hooks only. Rationale: REQ-032's gate is a scheduling *rule* (needed inline); this is a *procedure* (cited on demand), same split as Crash Recovery and the Lock Guard.

### 2. Requirement → edit site

| Requirement | Site |
|---|---|
| Naming `worktree-agent-REQ-NNN-<suffix>` | work-reference new section. Suffix = sanitized `[a-z0-9-]` token derived from the REQ **filename slug** as a text op before any shell quoting (CLAUDE.md injection rule). Branch name **==** worktree dir basename, so one grep correlates both. |
| Sole integrator + seams | work-reference new section; **one bullet in `actions/work.md` Step 6's builder instruction list** (~line 395, beside the `write_set` bullet): builders never write the main tree; a shared file needing one-line wiring is handed back as the line, the orchestrator applies it. |
| Post-merge verification gate | **`actions/work.md` Step 8, new paragraph immediately under `**On success:**` (line 535), before substep 1** — literal phrase "post-merge verification". No new numbered step (Step 9 is a named entry point; renumbering is a hazard). Procedure in work-reference. |
| Happy-path cleanup | work-reference new section, invoked from Step 8 substep 6 by one appended sentence: after the archive move, `git worktree remove` (no `--force`) then `git branch -d`. `-d`/`remove` refusal = unmerged or dirty = **stop and report**, never force. |
| Crash-path sweep | **work-reference `## Crash Recovery (Step 1)` only** (new numbered step 4 after line 208) — `actions/work.md:119` already points there, so hot-path cost is zero. |
| State stays home | work-reference new section, one bullet. |
| Preconditions + degrade | work-reference new section, opening paragraph (`git worktree list` probe + harness parallel-subagent capability; absent ⇒ silently run the serial loop), mirroring `actions/board.md` Step 2's report-and-stop precedent in generalized language. |

Also: one `## Rules` bullet in `actions/work.md` (~line 683) and one Common Rationalizations row ("the branch merged fine, I'll `-D` the leftover" → `-d` refusing *is* the assertion). `docs/work-guide.md` — no change (user-facing surface is unchanged; flagged below).

### 3. Interlock with REQ-032 / room for REQ-035 / REQ-036

The rule: **delegate to the gate by reference; never restate its mechanics.** Concretely —

- work.md Step 1 hook (appended to the existing Parallel-dispatch subsection at line 194): *"Worktree isolation changes the **consequence** of an overlap — merge cost instead of interleaved writes — not the gate: disjointness and serial-only classes still decide what may be co-dispatched."* True before and after REQ-035/036.
- Say nothing about how N claims are recorded. Both the crash sweep and cleanup's pass say only "skip any worktree whose REQ is exempted by the live-claim gate this pass already applies" — pointing at `work-reference.md` → Crash Recovery's concurrency gate / `cleanup.md` Pass 0's live-claim gate, which REQ-035 rewrites in one place.
- Add one sentence: *"Nothing here requires concurrency — a single builder in a worktree is the same one-claim loop."* This makes the mode shippable and correct today, while REQ-035 gates real co-dispatch.

### 4. `actions/cleanup.md`

**New `### Pass 5: Orphaned Worktrees (consent-gated)`** after Pass 4 (line 115). Not folded into Pass 4 — different object class (git branches, not `do-work/` artifacts) and the only pass needing interactive consent. Contents: enumerate `git worktree list --porcelain` + `git branch --list 'worktree-agent-*'`; `git worktree prune` first; attempt `worktree remove` → `branch -d`; refusal ⇒ unmerged ⇒ **list it and ask** (load `crew-members/clear-questions.md` before the prompt — mandatory contract); **non-interactive run (including Step 10's automatic cleanup) reports only, never deletes**; live-claim gate by reference.

Ripple edits (Closed-Enumerations discipline): `Five passes, in order` (line 25) → six; Pass 4's "This is the one place cleanup deletes" (line 115) must be narrowed; **What This Action Does NOT Do** (line 224); Reporting block (line 141); Verification Checklist "All 5 cleanup passes" (line 253).

### 5. Ratchet (recommended — cheap, REQ-032 precedent)

Four `assert_contains` in `_dev/tests/contract-regressions.sh`, anchored on text that can't survive gutting (D-03 lesson): `actions/work-reference.md` ⊃ `worktree-agent-REQ-` and `git branch -d`; `actions/work.md` ⊃ `post-merge verification`; `actions/cleanup.md` ⊃ `worktree-agent-`.

### 6. Testing

`bash _dev/tests/contract-regressions.sh` green; `cd tools/queue-kanban && go test ./...` (untouched, but the suite is the repo's floor). Mutation-test each new ratchet (delete the guarded text, watch it fire, restore). GREEN grep from the Red-Green Proof: `grep -i "worktree" actions/work.md actions/work-reference.md actions/cleanup.md` must show naming convention, post-merge gate, and both cleanup owners.

### 7. Decisions, latitude, gaps

- **Merge, never rebase** (latitude exercised) — rebase rewrites builder commits so `git branch -d` no longer sees them merged, destroying the free assertion the REQ makes load-bearing. `--no-ff` preserves "integrated by orchestrator" provenance.
- **Verification granularity: per-batch default** (latitude) — rule stated as *the unit you verify is the unit you roll back*; per-merge allowed when checks are cheap.
- **Worktrees outside the repo tree** (latitude, REQ-recommended) — a worktree inside the repo would be seen by cleanup Pass 3a's "misplaced `do-work/`" scan and by `-uall` status scans: a concrete corruption path, not just tidiness.
- **`-d`/`remove` refusal is the only merged-ness test** — no parallel `--merged` computation to drift.
- **Extra hazard stated beyond the REQ:** where a consumer *commits* `do-work/`, a worktree carries a **stale snapshot** of the queue at branch point — builders treat it as absent, never write it.
- **Not covered:** `docs/work-guide.md` (no user-facing behavior change until REQ-035); `SKILL.md` (2-word headroom — no routing change); no board/parser change (no new frontmatter field).

**Plan validation (orchestrator):** Every Detailed Requirement maps to an edit site; no orphan tasks (the ratchet traces to REQ-032's precedent and the Red-Green GREEN condition). ⚠ Scope spans ~5 files with heavy ripple edits in `actions/cleanup.md` — inherent to adding a cleanup pass; explorer confirmed every ripple site plus one the plan missed (`actions/cleanup.md:16` carve-out, `docs/cleanup-guide.md:7` pass list). Corrected anchors: Step 1 hook goes after line 186 (inside the Parallel-dispatch subsection), Step 8 happy-path cleanup becomes substep 8 (substeps are not named entry points). No file conflicts (nothing else in flight).

*Generated by Plan agent*

## Exploration

## REQ-033 Explore — edit sites, patterns, seams

### 1. `actions/work.md` (13,485 words; no word-budget test — only SKILL.md has one at 2650, currently **2648** = 2 words headroom, plan's claim confirmed)

- **Crash-Recovery pointer**: line **119** (`**Crash Recovery:** … per actions/work-reference.md → **Crash Recovery (Step 1)**`). Already a pure pointer — a worktree sweep added to the reference section costs zero hot-path words, as planned. Lock guard is line **117**, above it.
- **REQ-032 parallel-dispatch gate**: heading paragraph **180**, bullets **182–186**; `**Serial-only resource classes.**` **188–193**; blank **194**; `**Queue status summary:**` **195**. ⚠️ The plan says "appended … at line 194" — that lands *after* the Serial-only block, not inside the Parallel-dispatch subsection. Insert after **186** to keep the hook inside the subsection it qualifies.
- **Step 8 archive**: `**On success:**` **535**, substep 1 **537** → post-merge-verification paragraph goes at **536**. Substep 6 is prose at **548** followed by the archive-behavior table **550–554**; substep 7 (deferred prime links) **556–567**. A sentence appended to 548 would sit *before* the table that decides where the file goes — prefer a short **substep 8** after 567 (substeps are not named entry points; only Step 9 is, per line **581**).
- **Step 6 builder brief**: instruction bullets **389–398**; `write_set` write-boundary bullet at **395** (already says out-of-set need is "stop-and-report … never a silent write" and cites the Step 1 gate) — the "builders never write the main tree / hand back the line" bullet belongs immediately after **395**. Related: Step 6 heartbeat **379**, `background-agents.md` durability **377**, approach directives **375**.
- **Rules** bullet after **683** (existing write_set rule). **Common Rationalizations** rows **708–717**; append after **717** (row 717 is the serial-only row — same shape).

### 2. `actions/work-reference.md`

- `## Crash Recovery (Step 1)` **199–210**: gate **203**, numbered steps **206/207/208**, closing paragraph **210**. New step **4** inserts after **208**.
- `## Concurrent-Orchestrator Lock Guard (Step 1)` starts **212** → new `## Worktree Dispatch Mode (Step 1)` at **211**. Confirmed clean seam.
- `## Folder Structure` **53–77** — describes only the `do-work/` tree; nothing there contradicts worktrees, but the "worktrees live outside the repo tree" decision is worth one line near **75–77** or in the new section.
- No existing text contradicts "builders never write the main tree." Closest is **75** ("`working/` … Immutable to all actions except the work pipeline") — compatible.

### 3. `actions/cleanup.md`

- `Five passes, in order:` **25**. Pass 4 **106–115**; line **115** is the "This is the one place cleanup deletes" claim. `### Repoint Documentation Links` **117** runs *after all passes* — **Pass 5 must be inserted at line 116**, before it.
- **Consent-gate wording to reuse**: there is **no existing interactive prompt in cleanup.md** — this pass would be its first. Reuse the loader idiom from `actions/reserve.md:30` / `actions/capture.md:151` ("Load `crew-members/clear-questions.md` and ask with your environment's ask-user prompt"). `crew-members/clear-questions.md`'s JIT_CONTEXT already states the trigger condition and lists callers as *illustrative* — no update strictly required, but adding cleanup keeps the list honest.
- **Live-claim gate to cite by reference**: Pass 0's gate, line **31** (also mirrored at `work-reference.md:363`, `work-reference.md:203`).
- **Ripple sites (all confirmed)**: **16** (the "lone exception is Pass 4" carve-out — plan missed this one), **25**, **115**, Reporting block **141–151**, staging prose **220** + commit block **187–201**, **224**, **226**, **239** ("Run all 5 passes"), **253** ("All 5 cleanup passes … Passes 0–4"), plus a Verification-Checklist line for the new pass.

### 4. `actions/board.md` — degrade pattern to imitate

Line **7** ("Because it's compiled, this action needs the **Go toolchain** … It degrades gracefully when Go is absent: it reports and stops, never blocking the rest of the skill"), `### Step 2: Precondition — Go toolchain` **41–47**, Red Flag **130**, checklist **134**. Note the difference: board *reports and stops*; the worktree precondition must **silently fall through to the serial loop** (opt-in advanced mode) — imitate the shape, not the stop.

### 5. `docs/work-guide.md`

REQ-032's bullet is line **87** ("One REQ at a time by default; parallel dispatch is opt-in"), inside `## What run does (and does not) do`. Plan's "no change" is defensible; if a sentence is added it goes at the end of **87**.

### 6. `_dev/tests/contract-regressions.sh`

- REQ-032 template block: **121–137** (comment **121–123**, then three `assert_contains`). Copy that shape verbatim.
- Tail: last check ends **436** (`rm -rf "$redaction_order_workdir"`); `if [ "$fail_count" -gt 0 ]` at **439**, final printf **443**. New asserts append at **437–438**.
- `assert_contains` greps with `grep -Eq` → **regex, not fixed string**: `git branch -d` is safe, but anchor patterns carefully (`worktree-agent-REQ-` fine).
- Existing cleanup-guide assert at **150–152** is the precedent for guarding a docs file.

### 7. Closed-enumeration sweep — must update

| File / line | Enumeration made stale |
|---|---|
| `actions/cleanup.md` 16, 25, 115, 141–151, 220, 224, 226, 239, 253 | pass count, "only place cleanup deletes", reporting lines, staging list |
| `docs/cleanup-guide.md` **7** ("Five passes, in order (matching actions/cleanup.md)") + new section after **21** | user-facing pass list — must move in lock-step |
| `actions/work-reference.md` **207** | Crash-Recovery strip list (REQ-035 also touches this; see below) |
| `crew-members/background-agents.md` **186–191** | Harness tiers 1–3 ("manual parallel/background spawns") — worktree dispatch is a *fourth* capability axis; JIT_CONTEXT line 3 is already condition-stated, so a mention is optional, not required |

**Consciously skip**: `SKILL.md` (no routing change; 2-word headroom), `next-steps.md`, `actions/forensics.md` (39, 127 name Pass 0 only), `crew-members/general.md` (no parallel-builder claims), `tools/queue-kanban/*` (no new frontmatter field), `actions/reserve.md` (reservations are a different worktree concept — worth one disambiguating clause so readers don't conflate `reserved` with `worktree-agent-*`).

### 8. Concerns

1. **Live-claim gate by reference is currently a broken promise.** `claimed_req` is a **single string** (`work-reference.md:223, 247`). Both Crash Recovery's gate (**203**) and cleanup Pass 0's gate (**31**) exempt only *one* REQ, and only for *other* sessions. REQ-035 (`do-work/queue/REQ-035-lock-multi-claim-representation.md`, `depends_on: [REQ-033]`) exists precisely to fix this. So the plan's "skip any worktree whose REQ is exempted by the live-claim gate this pass already applies" is correct *by construction* but today exempts almost nothing under real co-dispatch. The plan's own mitigation — "Nothing here requires concurrency" — is what keeps this honest; make sure that sentence actually ships, or the crash sweep can delete a live sibling builder's worktree.
2. **Same-session blind spot.** Crash Recovery's gate explicitly checks sessions "other than this one" (**203**). A worktree sweep inheriting that gate would, on a Step 10 → Step 1 loop, sweep the *same* session's own in-flight worktrees. REQ-035 owns the fix; REQ-033 should state the single-builder scope rather than re-derive the gate.
3. **Write-set collision with REQ-035/036.** Both declare `write_set: [actions/work.md, actions/work-reference.md]` — identical to REQ-033's surface. They `depends_on: [REQ-033]`, so serialization is already correct, but REQ-033 must not pre-empt their text: don't restate lock mechanics (035) and don't touch Step 5.5's mirror (036).
4. **Non-interactive cleanup.** Step 10 (`work.md:592`) calls cleanup automatically inside an unattended loop; cleanup.md has no interactivity test today (the lock guard's, `work.md:117`, is the only precedent). The "report only, never delete when non-interactive" rule needs a stated test, not self-assessment.
5. **`git branch -d` merged-ness is branch-relative.** `-d` checks merged into the *current* HEAD (or upstream); cleanup Pass 5 may run from any branch. Prescribe `git branch -d` from the integration branch, or state the assumption — otherwise "refusal = unmerged" silently becomes "refusal = you're on a different branch."
6. **Worktrees outside the repo tree** interacts well with cleanup Pass 3a (**77–90**, scans for any `do-work/` directory in the repo) and `-uall` scans — the plan's rationale checks out on disk.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `actions/work-reference.md` (modify) — new `## Worktree Dispatch Mode (Step 1)` section (line ~211, between Crash Recovery and the Lock Guard); Crash Recovery gains numbered step 4 (worktree sweep: merged → mechanical removal, unmerged → report only; live-claim gate by reference); one Folder Structure line (worktrees live outside the repo tree)
- `actions/work.md` (modify) — Step 1 hook inside the Parallel-dispatch subsection (after line 186): worktree isolation changes the consequence of overlap, not the gate; Step 6 builder bullet (builders never write the main tree; integration seams hand back the line); Step 8 post-merge verification paragraph under **On success:** + happy-path cleanup as substep 8; one Rules bullet; one Common Rationalizations row (`-d` refusal is the assertion — never `-D`)
- `actions/cleanup.md` (modify) — new `### Pass 5: Orphaned Worktrees (consent-gated)` before the Repoint Documentation Links section; ripple edits at lines 16, 25, 115, Reporting 141–151, staging/commit 187–220, 224, 226, 239, 253
- `docs/cleanup-guide.md` (modify) — pass list moves in lock-step ("Five passes" → six) + short Pass 5 section
- `_dev/tests/contract-regressions.sh` (modify) — four ratchets: `worktree-agent-REQ-` and `git branch -d` in work-reference.md; `post-merge verification` in work.md; `worktree-agent-` in cleanup.md

**Files I will NOT touch:** `SKILL.md` (2-word budget headroom, no routing change), `docs/work-guide.md` (no user-facing behavior change until REQ-035 — conscious skip), `crew-members/clear-questions.md` (JIT caller list is explicitly illustrative), `actions/reserve.md` (speculative disambiguation — YAGNI), `tools/queue-kanban/*` (no new frontmatter field), Step 5.5 mirror text (REQ-036's surface), lock-schema/claim mechanics (REQ-035's surface).

**Acceptance criteria (restated from REQ):** *(builder note: every box below is implemented — left unticked for the orchestrator's Step 6.3 qualification)*
- [ ] Naming convention `worktree-agent-REQ-NNN-<suffix>` stated; branch name == worktree dir basename; suffix sanitized as a text op before any shell use
- [ ] Sole integrator: builders never write the main tree; one-line wiring handed back as integration seams; orchestrator merges in dependency order
- [ ] Post-merge verification gate before archive (per-merge or per-batch; the unit verified = the unit rolled back)
- [ ] Happy-path cleanup at archive: `git worktree remove` + `git branch -d`; `-d` refusal = lost/skipped merge = stop and investigate, never force
- [ ] Crash-path sweep: merged leftovers removed mechanically (+ `git worktree prune`); unmerged reported, never auto-deleted — deletion only via `do-work cleanup` behind consent; sweep skips freshly-claimed live REQs per the existing per-file gate (by reference)
- [ ] State stays home: `do-work/` lives in the main tree only; builders get their brief in the dispatch prompt; untracked-files-don't-propagate stated explicitly
- [ ] Preconditions + graceful degrade to the serial loop (git worktree support + parallel-subagent harness); generalized language, no tool-specific APIs
- [ ] "Nothing here requires concurrency — a single builder in a worktree is the same one-claim loop" ships (keeps the mode honest until REQ-035)
- [ ] GREEN grep: `grep -i worktree actions/work.md actions/work-reference.md actions/cleanup.md` shows naming, post-merge gate, cleanup owners; `bash _dev/tests/contract-regressions.sh` green

## Decisions

**D-01 — Crash Recovery's worktree sweep is a paragraph after step 3, not a numbered step 4 inside the per-file loop.** DECIDE & STATE. The plan and Exploration both said "new numbered step 4 after line 208," but steps 1–3 sit under `For each REQ file that passes the gate:` — a per-`working/`-file loop. The sweep is keyed on `worktree-agent-*` *names*, and a leftover branch routinely outlives its `working/` file (the REQ archived, the branch didn't), so nesting it in that loop would both misapply per-file and miss exactly the orphan class the REQ exists to catch. Placed at the same anchor (immediately after step 3, before the closing "once every file … has been recovered" paragraph), explicitly labelled "runs once, not per file in the loop above."

**D-02 — Merge with `git merge --no-ff`, never rebase.** DECIDE & STATE (latitude the REQ granted on merge mechanics; followed the plan's recommendation). Rebasing rewrites the builder's commits, so `git branch -d` stops recognizing the work as merged — which destroys the free merged-ness assertion the REQ makes load-bearing in three separate places (Step 8 cleanup, Crash Recovery's sweep, cleanup Pass 5). `--no-ff` additionally preserves the merge commit as the "integrated by orchestrator" provenance record the consumer runner already produced.

**D-03 — Verification granularity: per-batch by default, stated as a rule rather than a fixed number.** DECIDE & STATE (latitude granted; followed the plan). The shipped rule is *the unit you verify is the unit you roll back* — per-batch verifies once after the batch's last merge and reverts the batch on red; per-merge is called out as better when checks are cheap or one REQ is unusually risky. Framing it as the rollback contract rather than a policy means it stays correct if a consumer picks the other granularity.

**D-04 — Worktrees live outside the repo working tree.** DECIDE & STATE (the REQ's `- [~]` open question, deferred to the builder with a recommendation; followed it). A worktree nested inside the repo is a second checkout that `actions/cleanup.md` Pass 3a's filesystem scan for stray `do-work/` directories would find and try to relocate into the canonical queue wherever the consumer commits `do-work/` — a concrete corruption path, not tidiness. Recorded in the new work-reference section and cross-noted in Folder Structure.

**D-05 — Step 8's `**On failure:**` path gets no worktree-cleanup substep.** DECIDE & STATE. A failed REQ's builder branch is unmerged by definition, so `git branch -d` would refuse; prescribing cleanup there would either be a dead branch of prose or an invitation to `-D`. The existing machinery already covers it end-to-end: the next Step 1's Crash Recovery sweep finds the leftover, classifies it unmerged, reports it, and hands deletion to cleanup Pass 5's consent gate.

**D-06 — Two additions beyond the plan's edit list, both inside declared files.** DECIDE & STATE. (a) `actions/cleanup.md` gained a Common Rationalizations row and a Red Flag for Pass 5 — the REQ's own Why cites observed 37-minute-old orphan branches, and "it's ancient, obviously abandoned, just `-D` it" is the specific failure that produces silently-lost builder work; both rows name Pass 5 and the `-d` assertion, so they pass the earned-not-generic test. (b) The Reporting section's "Archive is clean" line is now gated on Pass 5 also being empty — a clean archive says nothing about worktrees, and printing it over an unreported leftover is the exact staleness the Closed-Enumerations discipline targets.

**D-07 — The `post-merge verification` ratchet pattern is `[Pp]ost-merge verification`.** SILENT/DECIDE & STATE. `assert_contains` uses `grep -Eq` (case-sensitive), and the phrase ships sentence-initial in `actions/work.md` Step 8. Case-tolerant pattern rather than contorting the prose to hold a lowercase copy.

## Discovered Tasks

- [low] `crew-members/background-agents.md` (lines ~186–191) describes harness capability tiers 1–3 in terms of manual parallel/background spawns. Worktree dispatch is a fourth, orthogonal axis (isolated working directory per builder) that the tiering doesn't mention. Its `JIT_CONTEXT` is already condition-stated so nothing is *wrong*, but a reader sizing a harness against those tiers won't learn worktree isolation exists. Out of this REQ's declared write_set.

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified) — new `## Worktree Dispatch Mode (Step 1)` section (precondition + silent degrade to serial, `worktree-agent-REQ-NNN-<suffix>` naming with branch == dir basename and text-op sanitization, worktrees outside the repo tree, state-stays-home incl. the stale-snapshot hazard, sole integrator + integration seams, merge --no-ff never rebase, post-merge verification procedure, happy-path cleanup, crash-path pointer); Crash Recovery gains a `worktree-agent-*` sweep (merged → mechanical removal + prune; unmerged → report only; live-claim gate by reference; same-session co-dispatch caveat until REQ-035); one Folder Structure line
- `actions/work.md` (modified) — Step 1 Parallel-dispatch bullet (worktree isolation changes the consequence of overlap, not the gate); Step 6 builder bullet (builders never write the main tree; seams hand back the line); Step 8 **Post-merge verification** paragraph under **On success:** + new substep 8 (happy-path worktree cleanup, `git branch -d` from the integration branch); one Rules bullet; one Common Rationalizations row (`-d` refusal is the assertion, never `-D`)
- `actions/cleanup.md` (modified) — new `### Pass 5: Orphaned Worktrees (consent-gated)` (prune → remove → `-d`; refusal = unmerged = list-and-ask; interactivity test cited from the Lock Guard; non-interactive incl. Step 10 reports only, never deletes; live-claim gate by reference) + all pass-count/reporting/staging/checklist ripple edits
- `docs/cleanup-guide.md` (modified) — pass list in lock-step (six passes) + Pass 5 section + key-rules deletion line
- `_dev/tests/contract-regressions.sh` (modified) — four ratchets: naming convention + `git branch -d` in work-reference.md, `[Pp]ost-merge verification` in work.md, `worktree-agent-` in cleanup.md (all mutation-tested)

**What was done:** Documented the worktree dispatch mode end-to-end in the work pipeline — orchestrator-created worktrees/branches embedding the REQ id, builders isolated from the main tree with the orchestrator as sole integrator merging in dependency order, a post-merge verification gate before archive (per-batch default: the unit verified is the unit rolled back), and defined cleanup ownership (archive substep on the happy path; Crash Recovery sweep for merged leftovers; consent-gated cleanup Pass 5 for unmerged ones), with graceful degradation to the serial loop and explicit single-builder honesty until REQ-035 lands.

## Qualification

Passed — 5 files verified in `git diff --stat` (exact match to Scope; `tools/checks/scope-drift.sh`: zero drift), all eight acceptance criteria traced to shipped text (naming + sanitization, sole integrator + seams, post-merge gate with rollback-unit rule, `-d`-from-integration-branch cleanup, crash sweep with merged/unmerged split + gate-by-reference + same-session caveat, state-stays-home + stale-snapshot hazard, silent degrade contrasted with board.md's report-and-stop, the single-builder honesty sentence verbatim). P-A-U confirmed against the diff; no debug artifacts; no CLAUDE.md/AGENTS.md citations in added lines. `tools/checks/qualify.sh`: OK. Contamination check vs REQ-032: overlap on `actions/work.md`/`actions/work-reference.md`/`_dev/tests/contract-regressions.sh` is declared batch overlap (`related: [REQ-032]`, same UR), not contamination.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`; `cd tools/queue-kanban && go test ./...` (both re-run by the orchestrator, independent of the builder's runs)
**Result:** ✓ All passing (contract regressions: passed; go: ok)

**Red-green validation:**
- `grep -i worktree actions/work.md actions/work-reference.md actions/cleanup.md`: ✗ no operative dispatch mode before (captured in Red-Green Proof) → ✓ naming convention, post-merge verification gate, and both cleanup owners present
- Four new ratchets mutation-tested by the builder: each fired alone when its guarded text was removed/reworded, then restored green (transcript in builder report)

**New tests added:**
- `_dev/tests/contract-regressions.sh`: four ratchet assertions (naming convention, `git branch -d` assertion, `[Pp]ost-merge verification`, consent-gated worktree pass)

*Verified by work action*

## Review

**Overall: 73%** | 2026-07-28T22:24:01Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 75% |
| Test Adequacy | 75% |
| Scope | 92% |
| Risk | Low |
| Acceptance | Partial |

**Findings:** 3 important, 5 minor
**Acceptance:** Partial — every runnable check is green (contract regressions, `go test ./...`, the REQ's GREEN grep, all four ratchets mutation-fired individually, every named pointer resolves), but a cold-agent walkthrough finds the crash sweep unreachable on its own headline case and the merge point never placed in the Step sequence.
**Suggested testing:** 5 items
**Follow-ups created:** None — read-only review; 3 would-be follow-ups are listed under Findings → Important.

*Reviewed by review-work action*

---

### Verdict

**Approve with follow-ups** — the mode is specified end-to-end and the prose is genuinely good (dense, reasoned, non-generic), but three seams with the surrounding state machine are open: the crash sweep can't be reached in the orphan case it was written for, nothing says at which step the merge happens (so the three diff-consuming steps go blind), and the per-batch rollback default can't roll back through a per-REQ Step 8.

Route C | uncommitted (5 tracked files)

### What's built

- `actions/work-reference.md` gains a complete `## Worktree Dispatch Mode (Step 1)` section — precondition + silent degrade, `worktree-agent-REQ-NNN-<suffix>` naming with branch == dir basename and text-op sanitization, worktrees outside the repo tree, state-stays-home incl. the committed-`do-work/` stale-snapshot trap, sole integrator + integration seams, merge `--no-ff` never rebase, post-merge verification, happy-path cleanup — plus a `worktree-agent-*` sweep in Crash Recovery and a Folder Structure line.
- `actions/work.md` gets five hooks (Step 1 gate bullet, Step 6 builder bullet, Step 8 post-merge paragraph + substep 8, one Rules bullet, one Rationalizations row); `actions/cleanup.md` gets `### Pass 5: Orphaned Worktrees (consent-gated)` with every pass-count / carve-out / reporting / staging / Red-Flag / checklist ripple, mirrored in `docs/cleanup-guide.md`; four ratchets land in `_dev/tests/contract-regressions.sh`.
- Still missing: where in the pipeline the merge happens, how Steps 6.3/6.5/7 obtain evidence once the work is on a merged branch, and a reachable trigger for the crash sweep when `working/` is empty.

### Decisions / risks for you

- **D-01 (sweep outside the per-file loop) — accept.** The reasoning is right and better than the Plan's "numbered step 4": a leftover branch outlives its `working/` file, so a per-file loop would misapply and miss orphans. It just needs the reachability fix in Important 1 to actually deliver on that reasoning.
- **D-05 (no cleanup substep on the `On failure:` path) — accept.** A failed REQ's branch is unmerged by definition; `-d` would refuse and prescribing it there invites `-D`. Deferring to the crash sweep + Pass 5 consent gate is the coherent choice.
- **D-06 (two additions beyond the plan's edit list) — accept.** Both are inside declared files, both name Pass 5 and the `-d` assertion, and both pass the earned-not-generic test in `CLAUDE.md`. The Reporting gate on "Archive is clean" is exactly the Closed-Enumerations discipline applied correctly.
- **D-02/D-03/D-04 — accept as latitude**, with one caveat: D-03's *per-batch default* is the source of Important 3.

### Findings

**Important:**

1. **The Crash Recovery worktree sweep is unreachable in the case it exists for** (`actions/work-reference.md:212`, entry condition at `actions/work.md:119`). The sweep's own justification is "a leftover branch can outlive its `working/` file (the REQ archived, the branch didn't)" — but the only instruction that sends an orchestrator into that reference section is `actions/work.md:119`: *"if `do-work/working/` contains any `REQ-*.md` … Reset and re-queue everything else per `actions/work-reference.md` → **Crash Recovery (Step 1)**"*. When `working/` is empty (the exact orphan class named), the hot-path pointer never fires and the sweep never runs. The reference section's own opening (`work-reference.md:203`, "If any exist, a previous run may have been interrupted") reinforces the same conditional read. `actions/work.md` is in the declared `write_set`, so a one-clause fix ("…and, whether or not `working/` has files, run the worktree sweep once") was in scope. Would-be follow-up.

2. **The merge point is never placed in the Step sequence, and the three evidence-consuming steps are left assuming uncommitted main-tree changes.** The mode says builders "commit on their own branch" and the orchestrator "merges, in dependency order" (`actions/work-reference.md:235`, `:237`), and Step 8's gate assumes a merge already happened (`actions/work.md:539`) — but no shipped text says *when*, and the steps in between all read the working diff:
   - `tools/checks/qualify.sh:46` computes `changed_file_list` from `git diff` + `git diff --staged` only, and its debug-artifact scan at `:112` uses the same source. After a merge the main tree is clean, so Step 6.3's check 1 degrades to WARNs and check 4 scans nothing.
   - `actions/review-work.md:30`/`:46` (pipeline mode) gets the diff via `git diff`, then: *"no implementation changes (pipeline mode), report that there's nothing to review and exit."* Step 7 silently skips.
   - Step 9's Commit Phase (`actions/work.md:590`) stages "implementation files from the Implementation Summary" and then "validate the staged file list against the Implementation Summary" — those files are already committed, so nothing stages and the validation mismatches on every worktree-dispatched REQ. The `## Rules` bullet "One commit per request" is also now false (builder commits + a `--no-ff` merge commit, exactly the "integrated by orchestrator" shape the UR screenshot shows).
   None of this is mentioned in the new text. Would-be follow-up: state the merge point and re-point the diff-based checks at the merge range (e.g. `<merge-base>..HEAD`) when the mode is on.

3. **The per-batch verification default cannot roll back through a per-REQ Step 8.** `actions/work-reference.md:239` makes per-batch the default ("verify once after the batch's last merge, and a red result reverts the batch") while `actions/work.md:539` requires verification "**before** the substeps below archive anything" — but Step 8 runs per REQ. For every REQ except the batch's last, those two are unsatisfiable together: either the orchestrator archives before the batch verification (violating the gate the REQ exists to create), or it must hold all Step 8s until the last merge — which nothing says. And once REQ #1 is archived, committed (Step 9) and its branch deleted (substep 8), "reverts the batch" has no clean unit left to revert. Dormant while the mode ships single-builder (batch of one == per-merge), but it becomes live the moment REQ-035 enables co-dispatch. Would-be follow-up (or an explicit clause folded into REQ-035).

**Minor:**

4. **The four ratchets pin strings, not the sections they claim to protect** (`_dev/tests/contract-regressions.sh:439-464`). Verified by mutation in a scratch copy: each assertion fires when its exact string is reworded, but deleting the entire `### Pass 5: Orphaned Worktrees (consent-gated)` section from `actions/cleanup.md` leaves the suite **green** (`worktree-agent-` survives in the line-16 carve-out, the Rationalizations row, the Red Flag and the checklist), and deleting the entire `## Worktree Dispatch Mode (Step 1)` section from `actions/work-reference.md` also leaves it **green** (the crash sweep keeps both guarded strings alive). The Plan's own criterion was "anchored on text that can't survive gutting"; anchoring on the headings (`Pass 5: Orphaned Worktrees`, `## Worktree Dispatch Mode`) would meet it.

5. **`actions/work.md` gained 413 words against a stated hook budget of ≤200** (Plan §1) / "~230 words" (P-A-U `[PLAN]`), and `[APPLY]` claims "Code written exactly as planned." The overrun is real (13,485 → 13,898 words in a file loaded every run) and undisclosed — no `D-XX` covers it. The bulk is substep 8 (`actions/work.md:573`) restating work-reference's whole command sequence including the `-d`-is-HEAD-relative rationale, which the Plan explicitly placed in the reference ("Procedure in work-reference"). Related drift surface: that rationale now exists in four near-identical copies (`work-reference.md:215`, `:241`, `work.md:573`, `cleanup.md:127`), so a future correction must find all four — the exact copy-paste pattern `CLAUDE.md` warns about.

6. **"Integration branch" is load-bearing but never defined, and worktree *creation* is never prescribed.** The term carries the merged-ness assertion in six places, and Pass 5 even branches on it ("If you cannot establish which branch that is, delete nothing"), yet nothing says what it is (the branch the orchestrator's main tree is on? a named integration branch? the default branch?). Symmetrically, `git worktree add` appears **zero** times across `actions/` and `docs/` — removal is prescribed three times, creation not once, and the base branch a builder forks from is unstated. The REQ's What says "orchestrator-created", so the owner is implied, but a cold agent gets exact removal commands and has to improvise the creation.

7. **Cleanup Pass 5's scope sentence contradicts its own enumeration step.** `actions/cleanup.md:119` says "Only `worktree-agent-REQ-NNN-*` names are in scope … every other worktree is the user's own and is never touched", but step 2 (`:126`) enumerates `git branch --list 'worktree-agent-*'` and step 3 acts on what step 2 produced — with no "any other name: not ours" bullet like the one Crash Recovery has (`work-reference.md:217`). Practical consequence: the UR's own motivating evidence is `worktree-agent-<hash>` branches (no REQ id), which the scope sentence excludes and the enumeration includes.

8. **A dangling "same words" instruction.** `actions/work-reference.md:214` says "Report the exemption in the same words the gate uses" — but Crash Recovery's gate (`:205`) defines no report wording; only `actions/cleanup.md:31` (Pass 0) does. Pass 5's identical phrasing (`cleanup.md:123`) resolves correctly; the work-reference copy does not.

**Nit:**

9. The mode's only `actions/work.md` entry point (`:187`) sits inside the **Parallel dispatch** subsection, whose opening says "Floor agents ignore this subsection entirely" — while the mode itself ships explicitly single-builder ("Nothing here requires concurrency"). A single-builder user never meets the mode from `work.md`.
10. `actions/work-reference.md:231` describes Pass 3a as scanning "for any `do-work/` directory outside the project root"; Pass 3a actually scans *inside the repo, excluding* the project root (`cleanup.md:81`). Reads as "outside the repo" on first pass.

### Requirements Checklist

- [x] Naming convention `worktree-agent-REQ-NNN-<suffix>`; branch name == worktree dir basename; suffix sanitized as a text op before shell use — delivered (`work-reference.md:229`), injection rule stated explicitly
- [x] Sole integrator: builders never write the main tree; one-line wiring handed back as integration seams; orchestrator merges in dependency order — delivered (`work-reference.md:235`, `work.md:397`)
- [~] Post-merge verification gate before archive (per-merge or per-batch; unit verified = unit rolled back) — delivered as text (`work-reference.md:239`, `work.md:539`), but the per-batch default is unsatisfiable against per-REQ Step 8 (Important 3)
- [x] Happy-path cleanup at archive: `git worktree remove` + `git branch -d`; refusal = stop and investigate, never force — delivered (`work-reference.md:241`, `work.md:573` substep 8), with the HEAD-relative caveat correctly captured
- [~] Crash-path sweep: merged removed mechanically + prune; unmerged reported never auto-deleted; deletion only via `do-work cleanup` behind consent; skips freshly-claimed live REQs by reference — content fully delivered (`work-reference.md:212-219`) but unreachable when `working/` is empty (Important 1)
- [x] State stays home: `do-work/` in the main tree only; brief comes via the dispatch prompt; untracked-files-don't-propagate stated — delivered (`work-reference.md:233`, `:79`), plus the committed-`do-work/` stale-snapshot hazard beyond the REQ's ask
- [x] Preconditions + graceful degrade to the serial loop; generalized language, no tool-specific APIs — delivered (`work-reference.md:227`), correctly contrasted with `actions/board.md`'s report-and-stop
- [x] "Nothing here requires concurrency — a single builder in a worktree is the same one-claim loop" ships — delivered verbatim (`work-reference.md:225`), reinforced at `:219`
- [x] GREEN grep shows naming, post-merge gate, both cleanup owners; contract regressions green — verified by running both
- [x] Constraint: REQ-032's write-set gate referenced, never repealed — `work.md:187` is precise ("changes the consequence… not the gate")
- [x] Constraint: no new action file; no user-managed parallel sessions suggested — held
- [x] Non-pre-emption of REQ-035 (lock/claim mechanics) and REQ-036 (Step 5.5 mirror) — verified: no lock-schema text added, Step 5.5 untouched, and `work-reference.md:219` states the single-builder limit by reference instead of re-deriving the gate

### Acceptance Testing

**Result: Partial**

- `bash _dev/tests/contract-regressions.sh` → `Contract regression checks passed.`
- `cd tools/queue-kanban && go test ./...` → `ok github.com/knews2019/skill-do-work/queue-kanban`
- REQ's GREEN grep → naming convention (5 hits across the three files), `Post-merge verification` at `work.md:539` and `work-reference.md:239`, both cleanup owners (Step 8 substep 8 + Crash Recovery sweep + cleanup Pass 5) all present
- **Mutation-tested all four new ratchets independently** in a scratch copy: each fires alone with the correct message and restores green. **Additional mutation the builder did not run:** deleting the whole `### Pass 5` section, and separately the whole `## Worktree Dispatch Mode` section, both leave the suite green → Minor 4
- **Cold-agent pointer resolution:** all 17 named pointers resolve — `Worktree Dispatch Mode (Step 1)` (×7), `Pass 5: Orphaned Worktrees (consent-gated)` (×2), `Concurrent-Orchestrator Lock Guard` → *Interactivity test* (`work-reference.md:344`, and its two conditions are environment-general, so the citation is portable to cleanup), `Crash Recovery (Step 1)`, Pass 0's live-claim gate, `actions/board.md`'s Go check, Pass 3a, `crew-members/clear-questions.md`, `work.md` Step 10, Step 8 substep 8
- **Closed-enumeration sweep:** no stale "Five passes" / "Passes 0–4" / "all 5 passes" survives anywhere in shipped files; `docs/cleanup-guide.md` moved in lock-step
- **Diff hygiene:** `git diff --check` clean; no TODO/FIXME/debug artifacts; no `CLAUDE.md`/`AGENTS.md` citations in added lines; P-A-U all `[x]`; scope exactly the 5 declared files, zero drift
- **Why Partial, not Pass:** the deliverable is executable-by-an-agent prose. Walking it as a cold agent, the crash sweep has no reachable trigger on its headline case, and following the mode leaves Steps 6.3 / 6.5 / 7 / 9 with no defined source of evidence. Both are gaps in what an agent would actually do, not stylistic.

### Suggested Additional Testing

- **Live dry-run on a scratch clone:** `git worktree add -b worktree-agent-REQ-999-demo ../tmp-wt`, commit there, `git merge --no-ff`, then walk Step 8 substep 8 — confirm `branch -d` succeeds from the integration branch and refuses from an unrelated branch (the whole `-d`-as-assertion contract rests on this).
- **The reachability case for Important 1:** repo with an orphan `worktree-agent-*` branch and an **empty** `do-work/working/` — hand the current text to a fresh agent at Step 1 and see whether it ever opens the sweep.
- **Pass 5 both branches:** one merged and one unmerged leftover, run `do-work cleanup` interactively (expect mechanical removal + list-and-ask) and again dispatched as a subagent (expect report-only, nothing deleted).
- **Evidence path for Important 2:** run `tools/checks/qualify.sh` against a REQ whose implementation arrived by merge with a clean working tree, and confirm what Step 6.3 and Step 7 actually see.
- **Human read of `docs/cleanup-guide.md` Pass 5** — it is the user-facing consent promise ("asks before deleting"); worth one pair of eyes that it can't be read as "cleanup deletes my branches".

**Adversarial verification (orchestrator addendum):** 12 serious findings (3 review + 9 contradiction-hunt) each attacked by 2 independent refuters. **Confirmed (2 root causes):** (A) the merge point is never placed in the step sequence and the evidence-consuming steps assume uncommitted main-tree work — `tools/checks/qualify.sh` reads `git diff`/`--staged` only, `actions/review-work.md` pipeline mode exits on an empty diff, Step 9 stages already-committed files and its validation mismatches, and the "one commit per request" rule is false in this mode → follow-up REQ-037; (B) the deterministic `worktree-agent-REQ-NNN-<suffix>` name has no uniqueness component, so a crash-recovered REQ re-dispatches into the name its surviving unmerged leftover occupies (creation fails; report-only sweep never frees it) → follow-up REQ-038. **Refuted (killed):** crash-sweep unreachability (the sweep text runs once-per-entry, not per-file, and both refuters showed the entry path holds), per-batch-rollback unsatisfiability (co-dispatch is explicitly banned until REQ-035, so batch>1 is unreachable — becomes live with REQ-035 and is noted there), substep-8 stranding, integration-branch unobtainability, Pass 5 interactivity contradictions, missing-creation-trigger. **Split (judged Minor, report only):** Pass 5's lead scope sentence. The review's would-be follow-up 1 (sweep reachability) was NOT created — refuted; 2→REQ-037; 3→noted in REQ-035's body via its existing "same-session co-dispatch" scope rather than a new REQ.

## Lessons Learned

**What worked:** The exploration's closed-enumeration sweep — it caught `docs/cleanup-guide.md`'s pass list and the cleanup.md line-16 carve-out the plan missed, so no stale "five passes" text shipped. Shipping the mode honestly single-builder ("Nothing here requires concurrency…") kept the spec true instead of promising bookkeeping that doesn't exist yet — verifiers repeatedly leaned on that sentence to kill would-be contradictions.
**What didn't:** Placing the merge everywhere except in the step sequence. The prose says who merges and how (`--no-ff`, dependency order) but never *when*, and every downstream evidence consumer (qualify diff, review diff, Step 9 staging) still assumes uncommitted main-tree work. A mode that changes *where work lives* has to walk every consumer of "the diff" — the review found this, the builder and orchestrator both missed it.
**Worth knowing:** REQ-037 (merge placement + evidence re-pointing) and REQ-038 (name uniqueness on re-dispatch) gate real worktree use; until they land the mode is documented but should be exercised with care. The four new ratchets pin phrases, not sections — deleting a whole section can leave the suite green (review Minor 4); anchor on headings when ratcheting a section's existence.

## Orientation

The work pipeline now documents worktree dispatch end-to-end: orchestrator-created `worktree-agent-REQ-NNN-*` worktrees/branches, builders isolated from the main tree with the orchestrator as sole integrator, a post-merge verification gate before archive, and defined cleanup ownership across the happy path (Step 8 substep 8), crash path (Crash Recovery sweep), and consent path (cleanup's new Pass 5 — the action's first consent-gated pass). Lives in the work-pipeline reference (`actions/work-reference.md` → Worktree Dispatch Mode) with hooks in `actions/work.md` and `actions/cleanup.md`. **[MAP CHANGED]** — new dispatch mode beneath the Step 1 gate; cleanup grew from five passes to six; two coherence follow-ups (REQ-037/038) gate real use. No prime files listed; no board/parser surface touched.
