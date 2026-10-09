# REQ-653 exploration: citations of the fan-out prose to be moved

Scope searched: every git-tracked file under skills/ (all four packages), _dev/ (tests, primes), CLAUDE.md, suite/. Excluded: CHANGELOG*.md, do-work/archive/, do-work/runs/, decisions/records/, decisions/log.md. Soft-wrapped matches (phrase split across one newline) are included. Search is case-sensitive as given ("Fan-out" capital F only).

Moved ranges assumed: work-reference.md 394-488; work.md 165-174, 288, 296, 368. Rows marked "INSIDE moved range" travel with the text (internal cross-references to re-point to the new file or to "below/above").

## Must-fix findings (would break tests or leave dangling pointers)

F1. prescribed-shell-canonicalization.sh:86-99 requires `skills/do-work/actions/work-reference.md` to contain the literal `../docs/prescribed-shell-primitives.md`. Its ONLY occurrence is line 437 (Hold both endpoints, link `../docs/prescribed-shell-primitives.md#state-across-command-blocks`), inside the moved range. Moving it fails this test unless (a) work-reference.md keeps another pointer to the guide, or (b) the test's core_site list adds `skills/do-work/actions/fan-out-reference.md` and drops/keeps work-reference.md appropriately.
F2. Outside-range citations inside the two source files that name moved sections (must be re-pointed to `actions/fan-out-reference.md`): work-reference.md:19, :55, :112, :491 (Reading a Builder-Authored Section cites Worktree Dispatch Mode (Step 1) -> State stays home), :848 (cites work.md Mid-Run Messages); work.md:35, :37, :39, :106, :286, :309, :315, :370, :475, :478, :494, :537, :550, :551, :583.
F3. Section-name check (shipped-package-reference-contract.sh) resolves `path` -> **Name** against the TARGET file's headings and bold labels. Every external `actions/work-reference.md` -> **Worktree Dispatch Mode...** bold citation will FAIL once the heading leaves work-reference.md unless re-pointed (or work-reference.md keeps a heading/bold label containing the name as whole words). Bold-form external sites: review-work.md:72, :74; cleanup.md:39, :134; restart-with-parallel-handoff.md:82; capture.md:124 (work.md -> **Mid-Run Messages (any step)**); background-agents.md:187 in do-work, do-work-toolbox, do-work-knowledge (wrapped "**Worktree Dispatch\nMode (Step 1)**" in the two siblings). Undelimited (non-bold) citations are NOT checked but go stale: capture-reference.md:166, work-guide.md:99, board.md:94, :120, board-guide.md:61, lessons-do-kanban.md:37, model.go:1950, verify.go:1025, :1079, :1275, :1392.
F4. core-checks.sh near_identical_cross_file_pairs: any "## Common Rationalizations" table in the new file must not copy rows from another file (ratio > 0.75 fails).

## Full citation table
| file:line | matched | quoted citation text | note |
|---|---|---|---|
| _dev/tests/shipped-package-reference-contract.sh:679 | Worktree Dispatch Mode | …xt inside it. # # The name may wrap: three shipped citations write "**Worktree Dispatch Mode\n(Step 1)**", # and stopping at the newline would… |  |
| _dev/tests/shipped-package-reference-contract.sh:1011 | Fan-Out Dispatch | …section first",             "`actions/work.md` → **Dispatch Mode** → *Fan-Out Dispatch*",             "actions/work.md",             "Di… |  |
| _dev/tests/shipped-package-reference-contract.sh:1039 | Worktree Dispatch Mode | …break is still the cited section",             "`actions/work.md` → **Worktree Dispatch Mode\n(Step 1)** for the flow",             "actions/w… |  |
| _dev/tests/shipped-package-reference-contract.sh:1041 | Worktree Dispatch Mode | …(Step 1)** for the flow",             "actions/work.md",             "Worktree Dispatch Mode\n(Step 1)",         ),         (             "an … |  |
| skills/do-work-board/actions/board.md:94 | Worktree Dispatch Mode | …e proof, never the badge (`../../do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). It just surfaces declared fi… |  |
| skills/do-work-board/actions/board.md:94 | Fan-Out Dispatch | …(`../../do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). It just surfaces declared file contention for a… |  |
| skills/do-work-board/actions/board.md:120 | Worktree Dispatch Mode | …at any builder count, per `../../do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch) is read verbatim into the car… |  |
| skills/do-work-board/actions/board.md:120 | Fan-Out Dispatch | … `../../do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch) is read verbatim into the card payload — no norm… |  |
| skills/do-work-board/docs/board-guide.md:61 | Worktree Dispatch Mode | … builders didn't collide (`../../do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). It just surfaces declared fi… |  |
| skills/do-work-board/docs/board-guide.md:61 | Fan-Out Dispatch | …(`../../do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). It just surfaces declared file contention.  It … |  |
| skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md:37 | Worktree Dispatch Mode | …on-interference proof (`../../../do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch); the `overlaps` badge is pure… |  |
| skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md:37 | Fan-Out Dispatch | …./../../do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch); the `overlaps` badge is purely an informational… |  |
| skills/do-work-board/tools/queue-kanban/model.go:1950 | Worktree Dispatch Mode | …e merge is the // non-interference proof (actions/work-reference.md → Worktree Dispatch Mode → // Fan-Out Dispatch). func annotateWriteSetOver… |  |
| skills/do-work-board/tools/queue-kanban/model.go:1951 | Fan-Out Dispatch | …erence proof (actions/work-reference.md → Worktree Dispatch Mode → // Fan-Out Dispatch). func annotateWriteSetOverlap(tickets []*Request… |  |
| skills/do-work-board/tools/queue-kanban/verify.go:1025 | Worktree Dispatch Mode | …t branch -d`'s own trap, documented at // actions/work-reference.md → Worktree Dispatch Mode, "Cleanup — happy path": // asked from an unrelat… |  |
| skills/do-work-board/tools/queue-kanban/verify.go:1025 | Cleanup — happy path | …documented at // actions/work-reference.md → Worktree Dispatch Mode, "Cleanup — happy path": // asked from an unrelated checkout, a perfectl… |  |
| skills/do-work-board/tools/queue-kanban/verify.go:1079 | Worktree Dispatch Mode | …ree-agent-REQ-NNN-<suffix> convention (actions/work-reference.md → // Worktree Dispatch Mode, "Naming"). Anchored at the prefix on purpose: a … |  |
| skills/do-work-board/tools/queue-kanban/verify.go:1275 | Worktree Dispatch Mode | … builder must never // write queue state (actions/work-reference.md → Worktree Dispatch Mode, "state // stays home" and "sole integrator"). //… |  |
| skills/do-work-board/tools/queue-kanban/verify.go:1392 | Cleanup — happy path | …appens to be (actions/work-reference.md → Worktree // Dispatch Mode, "Cleanup — happy path"). A detached repo-root checkout has no // branch… |  |
| skills/do-work-knowledge/crew-members/background-agents.md:187 | Worktree Dispatch Mode | …r `do-work run`, follow `../../do-work/actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** for the canonical hand-back, merge-ran… |  |
| skills/do-work-toolbox/crew-members/background-agents.md:187 | Worktree Dispatch Mode | …r `do-work run`, follow `../../do-work/actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** for the canonical hand-back, merge-ran… |  |
| skills/do-work/actions/capture-reference.md:166 | Worktree Dispatch Mode | … schedules on it at any builder count — `actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch), so an invented set is strict… |  |
| skills/do-work/actions/capture-reference.md:166 | Fan-Out Dispatch | …uilder count — `actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch), so an invented set is strictly worse than absen… |  |
| skills/do-work/actions/capture.md:124 | Mid-Run Messages | … during a live run in this session is routed by `actions/work.md` → **Mid-Run Messages (any step)** first \| `do-work/queue/` — work loop… |  |
| skills/do-work/actions/cleanup.md:39 | Worktree Dispatch Mode | … run finished — plural under fan-out (`actions/work-reference.md` → **Worktree Dispatch Mode** → Fan-Out Dispatch) — so no in-flight REQ sits … |  |
| skills/do-work/actions/cleanup.md:39 | Fan-Out Dispatch | …r fan-out (`actions/work-reference.md` → **Worktree Dispatch Mode** → Fan-Out Dispatch) — so no in-flight REQ sits terminal-in-`working/… |  |
| skills/do-work/actions/cleanup.md:121 | Fan-out | …/` alongside legacy REQs  ### Pass 4: Sweep Consumed Run Directories  Fan-out actions (code-review, deep-explore, multi-REQ wor… |  |
| skills/do-work/actions/cleanup.md:134 | Worktree Dispatch Mode | …orktree dispatch mode is interrupted (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**). Only `worktree-agent-REQ-NNN-*` name… |  |
| skills/do-work/actions/restart-with-parallel-handoff.md:82 | Worktree Dispatch Mode | …ay-only and gates nothing on its own (`actions/work-reference.md` → **Worktree Dispatch Mode** → *Fan-Out Dispatch*).  ### Step 6: Announce it… |  |
| skills/do-work/actions/restart-with-parallel-handoff.md:82 | Fan-Out Dispatch | … its own (`actions/work-reference.md` → **Worktree Dispatch Mode** → *Fan-Out Dispatch*).  ### Step 6: Announce it  End the final chat m… |  |
| skills/do-work/actions/review-work.md:72 | Worktree Dispatch Mode | …sses: `git diff <pre>..<merge_hash>` (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**).  **Standalone mode:** Read the commi… |  |
| skills/do-work/actions/review-work.md:74 | Worktree Dispatch Mode | …o-ff` merge is the normal merge case (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**).  Read the diff carefully. For large … |  |
| skills/do-work/actions/review-work.md:135 | Restatement sweep | …igger set also includes every element those members recorded on the **Restatement sweep:** line of their `## Review`, read from their arc… |  |
| skills/do-work/actions/review-work.md:405 | Restatement sweep | …e"] **Acceptance:** [Pass/Partial/Fail/Untested] — [1-line summary] **Restatement sweep:** redefined [this diff's own elements only, neve… |  |
| skills/do-work/actions/review-work.md:491 | Restatement sweep | …s at a wave end) was grepped and verified, and its result is on the **Restatement sweep:** line of the Append to REQ File template - [ ] … |  |
| skills/do-work/actions/work-reference.md:19 | Worktree Dispatch Mode | …fix it.  **Builders are not owners.** Any number may build at once (**Worktree Dispatch Mode** → Fan-Out Dispatch, below), because a builder w… |  |
| skills/do-work/actions/work-reference.md:19 | Fan-Out Dispatch | … owners.** Any number may build at once (**Worktree Dispatch Mode** → Fan-Out Dispatch, below), because a builder writes only its own tr… |  |
| skills/do-work/actions/work-reference.md:55 | Worktree Dispatch Mode | …utside the repo entirely and never carry their own copy of it — see **Worktree Dispatch Mode (Step 1)**.  ## Request File Schema — Full Frontm… |  |
| skills/do-work/actions/work-reference.md:112 | Worktree Dispatch Mode | …pick and the merge is the non-interference proof, never this field (**Worktree Dispatch Mode (Step 1)** → Fan-Out Dispatch, below), so nothing… |  |
| skills/do-work/actions/work-reference.md:112 | Fan-Out Dispatch | …erence proof, never this field (**Worktree Dispatch Mode (Step 1)** → Fan-Out Dispatch, below), so nothing schedules, gates, or dispatch… |  |
| skills/do-work/actions/work-reference.md:394 | Worktree Dispatch Mode | …` or `pending-answers`, and no branch asks for user confirmation.  ## Worktree Dispatch Mode (Step 1)  **Every run mode, the serial default in… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:396 | Isolation ladder | …Q branch, in its own git worktree wherever the harness supports one (*Isolation ladder*, below); the orchestrator merges those branches … | INSIDE moved range |
| skills/do-work/actions/work-reference.md:396 | Fan-Out Dispatch | …sequence, one `<pre>..<merge_hash>` range, one cleanup, each per REQ. Fan-Out Dispatch (below) adds only who picks the set and what neve… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:398 | Isolation ladder | …h (below) adds only who picks the set and what never parallelises.  **Isolation ladder — three rungs, and only the last one has no isola… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:408 | operative name | …t under whatever name it arrived with. **Hold that name as this REQ's operative name exactly the same way** (*The name actually create… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:408 | operative name | …name exactly the same way** (*The name actually created is this REQ's operative name*, below) — the merge, and any reporting, still ne… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:409 | Sole integrator | …travels on the branch.** The absolute-main-tree-path hand-back file (*Sole integrator*, below) is a **local-only** mechanism: it works … | INSIDE moved range |
| skills/do-work/actions/work-reference.md:410 | State stays home | …sagreements by syncing rather than by writing. This is the same rule *State stays home* (below) applies to a worktree's snapshot, widene… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:418 | operative name | …sent-gated)** if unmerged.  **The name actually created is this REQ's operative name.** Whatever `git worktree add` succeeded with — t… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:418 | Cleanup — happy path | …rgument (below), Step 8's `git worktree remove` and `git branch -d` (*Cleanup — happy path*, below), the crash sweep's own-session bookkeepi… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:418 | Hold both endpoints | …typed as a literal into each fresh command, never a shell variable (**Hold both endpoints as re-typed literals**, below). Nothing persists … | INSIDE moved range |
| skills/do-work/actions/work-reference.md:418 | operative name | …the variant worktree is never cleaned at all. **With no collision the operative name *is* the derived name**, so the common path behav… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:422 | State stays home | … check downstream. That is a corruption path, not just untidiness.  **State stays home.** **Every path under `do-work/` exists in the ma… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:424 | Sole integrator | …n a builder branch — step 2's queue guard below is what proves it.  **Sole integrator.** The builder never writes the main tree or its … | INSIDE moved range |
| skills/do-work/actions/work-reference.md:428 | When to merge | …erge commit as the "integrated by orchestrator" provenance record.  **When to merge, and the range every evidence step reads.** The o… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:433 | State stays home | …printed is queue state committed on the builder's branch — the write *State stays home* forbids. Stop and drop/revert those commits on t… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:433 | operative name | …fy. Then `git merge --no-ff --no-commit <operative_name>` (this REQ's operative name, *Naming* above — the branch the builder was actu… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:434 | Sole integrator | … integration seams, then commit** — stage the handed-back seam lines (Sole integrator, above) and `git commit`. Folding the seam into t… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:437 | Hold both endpoints | …nance hash that finalization records in the REQ's `commit:` field.  **Hold both endpoints as re-typed literals, never as shell variables.**… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:445 | Cleanup — happy path | …low-up — stop, revert to the last verified state, and re-dispatch.  **Cleanup — happy path (Step 9, after typed finalization success).** Aft… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:445 | operative name | …p_complete`, remove the builder's worktree and branch **by this REQ's operative name** (*Naming*, above — never re-derived from the sl… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:449 | Fan-Out Dispatch | …-gated)**, which asks first and only acts when a human can answer.  **Fan-Out Dispatch — several builders, one releaser.** Every guarant… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:451 | Auto-wave | …dy set itself and dispatches builders with **no confirmation gate** (*Auto-wave*, below). Without that flag the action runs the s… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:455 | Auto-wave | …e computation to queue-mode `advance` and dispatches what it claims (*Auto-wave*, below, states the policy around that delegation… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:456 | Sole integrator | …eanly and can still be jointly wrong. The **integration seam** rule (*Sole integrator*, above) is what covers that — and it works only … | INSIDE moved range |
| skills/do-work/actions/work-reference.md:460 | Auto-wave | …t it a later reader re-offers the shared tree as a simplification.  **Auto-wave — what the loop computes, and what it deliberatel… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:466 | Auto-wave | … dispatch.  **No confirmation gate — that is the deliberate change.** Auto-wave dispatches its computed set without asking, which… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:466 | operative name | …t it does **not** change: every per-REQ step below (one worktree, one operative name, one hand-back sequence, one merge range, one cle… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:466 | Isolation ladder | …EQ at a time on the branch rung and nothing is reported as an error (*Isolation ladder*, above).  **Serial-only — never parallelised, at… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:470 | Delegated integration | …ge time. Unique version numbers do not make a shared prepend safe.  **Delegated integration — the coordinator shape.** After a wave's builder… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:470 | coordinator shape | …ers do not make a shared prepend safe.  **Delegated integration — the coordinator shape.** After a wave's builders are dispatched, the or… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:470 | operative name | …run directory: the REQ id, the run directory, the hand-back path, the operative name, the wave membership, every wave member the coord… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:472 | Fan-out | …10 says why).  **The run directory is mandatory here, not optional.** Fan-out is a background fan-out, so `crew-members/backgro… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:474 | Fan-out | … (multi-REQ)`). Its slots map onto this pipeline:  \| Guardrail slot \| Fan-out use \| \| --- \| --- \| \| run directory \| `do-work/ru… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:478 | Sole integrator | … carry is lost silently. The one main-tree path a builder may write (*Sole integrator*, above) \| \| per-integrator input \| `REQ-NNN-inte… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:479 | Delegated integration | …written by the coordinator between integrators, never by a builder (**Delegated integration — the coordinator shape**, above) \| \| `manifest.m… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:479 | coordinator shape | …etween integrators, never by a builder (**Delegated integration — the coordinator shape**, above) \| \| `manifest.md` \| REQ id → builder, `… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:480 | Mid-Run Messages | …hand-back emphasis note when the user sent one (`actions/work.md` → **Mid-Run Messages (any step)**) — **the orchestrator's**, never wri… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:485 | State stays home | …side the worktree, against its own stale tracked copy of `do-work/` (*State stays home*, above) — so the builder silently reads a snapsh… | INSIDE moved range |
| skills/do-work/actions/work-reference.md:491 | Worktree Dispatch Mode | …t write the REQ file at all — the REQ lives in the main tree, which **Worktree Dispatch Mode (Step 1)** → *State stays home* forbids it to tou… |  |
| skills/do-work/actions/work-reference.md:491 | State stays home | … lives in the main tree, which **Worktree Dispatch Mode (Step 1)** → *State stays home* forbids it to touch — so it routes those section… |  |
| skills/do-work/actions/work-reference.md:848 | Mid-Run Messages | …run manifest carries a hand-back emphasis note (`actions/work.md` → **Mid-Run Messages (any step)**), read it first and order the sectio… |  |
| skills/do-work/actions/work.md:35 | Worktree Dispatch Mode | …ification, tests, staging, or commit. `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** is that contract for every mode, and i… |  |
| skills/do-work/actions/work.md:35 | Isolation ladder | …ee Dispatch Mode (Step 1)** is that contract for every mode, and its *Isolation ladder* defines what isolation means on a harness with n… |  |
| skills/do-work/actions/work.md:37 | Worktree Dispatch Mode | …nd returns every exclusion typed, and `actions/work-reference.md` → **Worktree Dispatch Mode** → *Fan-Out Dispatch* → **Auto-wave** states the… |  |
| skills/do-work/actions/work.md:37 | Fan-Out Dispatch | …yped, and `actions/work-reference.md` → **Worktree Dispatch Mode** → *Fan-Out Dispatch* → **Auto-wave** states the policy around it rath… |  |
| skills/do-work/actions/work.md:37 | Auto-wave | …k-reference.md` → **Worktree Dispatch Mode** → *Fan-Out Dispatch* → **Auto-wave** states the policy around it rather than a secon… |  |
| skills/do-work/actions/work.md:39 | Fan-Out Dispatch | …at proves two builders did not collide (`actions/work-reference.md` → Fan-Out Dispatch). A computed set asserts that its REQs are all *r… |  |
| skills/do-work/actions/work.md:106 | Worktree Dispatch Mode | …de `advance`; the policy around it is `actions/work-reference.md` → **Worktree Dispatch Mode** → *Fan-Out Dispatch* → **Auto-wave**. - **`--wa… |  |
| skills/do-work/actions/work.md:106 | Fan-Out Dispatch | …und it is `actions/work-reference.md` → **Worktree Dispatch Mode** → *Fan-Out Dispatch* → **Auto-wave**. - **`--wave N`** (integer flag,… |  |
| skills/do-work/actions/work.md:106 | Auto-wave | …k-reference.md` → **Worktree Dispatch Mode** → *Fan-Out Dispatch* → **Auto-wave**. - **`--wave N`** (integer flag, default mode o… |  |
| skills/do-work/actions/work.md:165 | Mid-Run Messages | …`, or no Open Questions section exists, skip this step entirely.  ### Mid-Run Messages (any step)  A user message can arrive while the r… | INSIDE moved range |
| skills/do-work/actions/work.md:170 | Delegated integration | …ion now. Under delegated integration (`actions/work-reference.md` → **Delegated integration — the coordinator shape**) the coordinator keeps … | INSIDE moved range |
| skills/do-work/actions/work.md:170 | coordinator shape | …egration (`actions/work-reference.md` → **Delegated integration — the coordinator shape**) the coordinator keeps the words in its own ses… | INSIDE moved range |
| skills/do-work/actions/work.md:286 | Worktree Dispatch Mode | …trigger and unsafe-branch policy, and `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** for this action's canonical hand-back … |  |
| skills/do-work/actions/work.md:288 | Landed hand-back | …ger still: it uses one worktree per builder regardless of overlap.  **Landed hand-back — consume it, never re-dispatch.** When the hand-… | INSIDE moved range |
| skills/do-work/actions/work.md:309 | Worktree Dispatch Mode | …work/` path stays the orchestrator's (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** → *Isolation ladder*). A shared file t… |  |
| skills/do-work/actions/work.md:309 | Isolation ladder | …(`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** → *Isolation ladder*). A shared file that needs one line of wiring is… |  |
| skills/do-work/actions/work.md:315 | Worktree Dispatch Mode | …read and execute the full sequence in `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** → **When to merge, and the range every… |  |
| skills/do-work/actions/work.md:315 | When to merge | …`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** → **When to merge, and the range every evidence step reads.** Run i… |  |
| skills/do-work/actions/work.md:368 | Restatement sweep | …ge_hash>` instead (`actions/review-work.md` Step 4, Get the Diff).  **Restatement sweep (MUST).** If this REQ's diff redefines something … | INSIDE moved range |
| skills/do-work/actions/work.md:370 | Delegated integration | …he list the integrator brief carries (`actions/work-reference.md` → **Delegated integration — the coordinator shape**). A run with no manifes… |  |
| skills/do-work/actions/work.md:370 | coordinator shape | … carries (`actions/work-reference.md` → **Delegated integration — the coordinator shape**). A run with no manifest has no wave. When a me… |  |
| skills/do-work/actions/work.md:475 | operative name | …ttled/created commit hashes, then remove any retained worktree by its operative name without force. ### Step 10: Loop or Exit  After i… |  |
| skills/do-work/actions/work.md:478 | Fan-Out Dispatch | …ss record. Under delegated integration (`actions/work-reference.md` → Fan-Out Dispatch → **Delegated integration — the coordinator shape… |  |
| skills/do-work/actions/work.md:478 | Delegated integration | …gated integration (`actions/work-reference.md` → Fan-Out Dispatch → **Delegated integration — the coordinator shape**) an integrator stops af… |  |
| skills/do-work/actions/work.md:478 | coordinator shape | …work-reference.md` → Fan-Out Dispatch → **Delegated integration — the coordinator shape**) an integrator stops after Step 9 and never run… |  |
| skills/do-work/actions/work.md:494 | Mid-Run Messages | … before dispatch — never - [~], no D-XX) □ Mid-run message: route per Mid-Run Messages; never stop the run for it □ Evidence gates: run … |  |
| skills/do-work/actions/work.md:537 | Mid-Run Messages | …s the run manifest's hand-back emphasis note orders them otherwise (**Mid-Run Messages (any step)**, above). Never lead with review scor… |  |
| skills/do-work/actions/work.md:550 | Worktree Dispatch Mode | …may require a second metadata commit (`actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)**). - `write_set` is display-only, and s… |  |
| skills/do-work/actions/work.md:551 | Worktree Dispatch Mode | … field — is the non-interference proof (`actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). Nothing schedules, gates, or… |  |
| skills/do-work/actions/work.md:551 | Fan-Out Dispatch | …ference proof (`actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). Nothing schedules, gates, or dispatches on it; … |  |
| skills/do-work/actions/work.md:583 | operative name | …name — I'll rebuild it from the REQ slug" \| Merge and clean up by the operative name held since dispatch (Step 6 hand-back, Step 8 sub… |  |
| skills/do-work/crew-members/background-agents.md:187 | Worktree Dispatch Mode | …e --force`. For `do-work run`, follow `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** for the canonical hand-back, merge-ran… |  |
| skills/do-work/docs/work-guide.md:99 | Worktree Dispatch Mode | …ifies, and archives them one at a time (`actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). The saving is in the build p… |  |
| skills/do-work/docs/work-guide.md:99 | Fan-Out Dispatch | …one at a time (`actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). The saving is in the build phase; everything af… |  |

Not in table (outside requested scope, informational): kb/wiki/** and ai-reports/** carry ~150 historical mentions; live queue files citing these names: do-work/HANDOFF-2026-08-18-queue-241-245.md:54 ("`skills/do-work/actions/work-reference.md` → **Worktree Dispatch Mode** is the contract"), do-work/user-requests/UR-143/input.md, do-work/queue/REQ-653-*.md.

## Test hits, which test and what it asserts

| test:line | asserts |
|---|---|
| _dev/tests/shipped-package-reference-contract.sh:679 | Comment only: explains the bold section name may wrap, citing the three shipped "**Worktree Dispatch Mode\n(Step 1)**" citations. No corpus pin. |
| _dev/tests/shipped-package-reference-contract.sh:1011 | Synthetic parse fixture "`actions/work.md` → **Dispatch Mode** → *Fan-Out Dispatch*": a chain's first bold name is the cited section. Parser-only; never reads the real file. |
| _dev/tests/shipped-package-reference-contract.sh:1039-1041 | Synthetic parse fixture: a name wrapped across a newline ("Worktree Dispatch Mode\n(Step 1)") is still parsed as the section. Parser-only. |
| _dev/tests/prescribed-shell-canonicalization.sh:86-99 | work-reference.md (and work.md, review-work.md, commit.md, capture.md, background-agents.md) must contain `../docs/prescribed-shell-primitives.md`. See F1. |
| _dev/tests/audit-lockins.sh:264 | rg over `**/actions/*.md` for 'If the script is missing or will not run' — any new action file is scanned automatically, must not contain that sentence. |
| _dev/tests/audit-lockins.sh:314-345 | Debug-artifact mention pin = 3 across work.md, review-work.md, work-reference.md. None of the 3 mentions are in work.md/work-reference.md (all in review-work.md), so the move does not change the count. |
| _dev/tests/contracts/core-checks.sh:719-790 | work.md must keep '### Mechanical Evidence-Gate Loop' (line 175, not moved) and Step 8/9 tokens; work-reference.md must keep Step 6.3/6.5 headings and Changelog/Commit procedure tokens. None in moved ranges. |
| _dev/tests/contracts/core-checks.sh:806-835 | work.md Step 6 must carry an active "Read/Always load `crew-members/shared-principles.md`" line; lines 288/296 are not that line. |
| _dev/tests/contracts/recovery-set-aside.sh:46-61 | Required phrases in work-reference.md ('one record at a time', 'Set-aside-by-recovery section', ...). Verified none sit only in 394-488. |

No _dev test pins any of the 16 heading/phrase strings against the live corpus.

## (a) Companion-reference citation pattern

Same package: backticked package-root-relative path, Unicode arrow, bold section name: `actions/<file>.md` → **Section Name** (chains add → *Subsection*). Examples:
- skills/do-work/actions/work.md:217 — `actions/capture-reference.md` → **Required Lessons Budget Contract**
- skills/do-work/actions/work-reference.md:117 — `actions/capture-reference.md` → **Fold-First Rule**
- sibling package: skills/do-work-knowledge/actions/setup-memory.md:68 — `actions/memory-reference.md` → **Hook Install Internals**
Cross-package uses the literal relative path from the citing file's dir: `../../do-work/actions/work-reference.md` → ... (board.md, background-agents.md in siblings).
Companion header convention: H1 "# <Name> — Reference" then a blockquote "> **Companion file to `actions/work.md`.** Holds ... Load only when you reach the step that references it — and read only the named section. If already in context, reuse it." (capture-reference.md:1-3, estimate-reference.md:1-3).

## (b) What shipped-package-reference-contract.sh checks

Scans every tracked *.md under each suite/modules.tsv source root (so the new file is scanned as soon as it is tracked). For each path-shaped token (prose, inline code, HTML comments, fence annotations; not fence payloads):
1. Same-package token (first segment is a present content dir like actions/, docs/, crew-members/): leadless resolves from the package root, `../`-led from the citing dir; must exist and stay inside the package in BOTH source and installed (.claude/skills/...) topologies.
2. Cross-package token (first segment do-work-board/-knowledge/-toolbox/do-work after `../`): must resolve in both topologies.
3. Section half: if the token is followed by `→ **Name**` (soft wraps allowed, not across a blank line), Name, normalized (inline code unwrapped, emphasis dropped, whitespace collapsed, trailing .,:; stripped, casefolded), must appear as whole words inside some ATX heading or bold run/label of the target file. Undelimited names are not checked. Dynamic tokens containing { } < > $ are skipped.
4. Markdown link targets: file exists, `#anchor` matches a GitHub heading slug in the target.
So the new file's citations pass if each cited path exists and each bold cited name is a heading/bold label (whole-word containment) in the target; and every external bold citation to the moved headings must point at the file that now holds them.

## (c) Registration of a new file under skills/do-work/actions/

None required.
- skills/do-work/SKILL.md routes only public actions (## Routing); companion references (capture-, estimate-, work-reference) are not listed there.
- _dev/tests/staged-skills-contract.sh:25-55 `core_files` is a must-exist list (includes work-reference.md, not estimate-reference.md); adding the new file is optional, not enforced.
- staged-skills-contract.sh:775 finds `actions/<public_action>.md` only for sibling-owned public route names; not affected.
- skills/do-work/tools/do-work-cli/internal/suitemanifest/suite_manifest.go:192 only checks actions/version.md.
- Globs that pick the file up automatically: audit-lockins.sh `**/actions/*.md`; core-checks.sh:806 and staged-skills-contract.sh:1034 rglob all skills *.md (path-existence of `actions/...` references; rationalization-row duplication).
