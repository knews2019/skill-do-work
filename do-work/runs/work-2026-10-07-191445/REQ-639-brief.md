# Builder brief — REQ-639 (Delegated integration: in fan-out mode the main session may hand each REQ's integration to one agent at a time and stay a coordinator)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-639-delegated-integration-coordinator-shape
- Branch: worktree-agent-REQ-639-delegated-integration-coordinator-shape (commit here; never merge, never rebase, never push)
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-639-delegated-integration-coordinator-shape.md — read it fully: What, Why, Decision, Verified Facts, Detailed Requirements, Constraints, Builder Guidance, Red-Green Proof, and the orchestrator's `## Exploration` and `## Scope` (exact line numbers for every insertion point and every restatement found so far).
- Hand-back file (outside both trees, never commit it): /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/87062c87-2ca6-4782-9a66-d1dde45f9c4a/scratchpad/REQ-639-handback.md
- Route B, tdd: false, impact-rule-change. Prose-only.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md (domain: general, so no extra domain crew). Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md (the REQ's prime_files; read it before editing). Lesson family to honour: alternate-writer-contract-drift in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-action-files.md (the satellite is `slugged: partial` and over the lesson budget, so it is not in required_lessons; read the six `[family: alternate-writer-contract-drift]` bullets, lines 52-60). Its rule: a changed rule is not done until every writer and every action-bearing reader that restates it is swept, by ownership condition, not only the cited lines.

## The change (decided at capture, do not reopen)
Copied verbatim from the REQ's `## Detailed Requirements`:

1. `skills/do-work/actions/work.md` Step 6: add a new condition before "Spawn a general-purpose agent". If this REQ's row in the run manifest records a landed hand-back and the hand-back file exists, consume it: do not dispatch, take the dispatch instant from that manifest row for record-timing-event, stamp `builder_handback_at` only if absent, and go straight to the hand-back merge, which already proves the branch state. Never re-dispatch a builder for a landed hand-back. State it as a condition so it also serves a fresh session after a crash (background-agents.md: agents whose findings file exists are done).
2. `work.md` Step 6 dispatch paragraph: the manifest row for a REQ carries the held dispatch instant, so an integrator can record the builder-work timing event without reading `dispatch_at` back from the file.
3. `work.md` Step 10, one sentence: under delegated integration the integrator never runs `advance --checkpoint` or the loop, because the checkpoint preserves only foreign records and would drop same-writer sibling claims; the coordinator runs it, then cleanup, after the last integrator returns.
4. `skills/do-work/actions/work-reference.md`, Worktree Dispatch Mode, Fan-Out Dispatch: a new paragraph "Delegated integration — the coordinator shape" stating, in this order:
   - (a) After the wave's builders are dispatched, the orchestrator may hand each REQ's integration (hand-back merge through Step 9) to one agent at a time, in series, in the same checkout. Release and changelog stay serial-only (cite the existing Serial-only paragraph).
   - (b) The integrator enters through plain `do-work run REQ-NNN`: recover reports the live claim and `advance REQ-NNN` continues it; the existing section steps skip what is already written; Step 6's landed-hand-back condition takes it to the merge. Never run-with-recovery, `--assume-sole-authority` or `--take-over` for an integrator: those reset every sibling claim and requeue the whole wave.
   - (c) One writer under the project root at a time: while an integrator runs, the coordinator writes nothing under the project root (no REQ sections, no manifest edits, no captures); otherwise this is the "two sessions in one working tree" case the Execution Model leaves unspecified. The coordinator's turn to write is the gap between integrators; its live notes go to its own session scratch until then.
   - (d) Why: the integration span is long and context-heavy, and running it in the main session blocks the conversation the user steers from.
   - (e) The integrator brief, `REQ-NNN-integrate.md` in the run directory, written by the coordinator in the gap before that integrator starts: REQ id, run directory, hand-back path, operative name, wave membership, any mid-run addendum text, and the instruction not to run Step 10.
   - Add one row to the guardrail-slot table: `per-integrator input | REQ-NNN-integrate.md — written by the coordinator between integrators, never by a builder`.
   - Keep "Dispatch mechanism is deliberately unspecified" as is. The new paragraph must not say subagent versus session.
5. `skills/do-work/docs/work-guide.md`, "Building several REQs at once": two sentences saying the session that runs `--fan-out` can stay a coordinator by handing each REQ's integration to one agent at a time, and that the saving is still build-phase only.

Constraints (verbatim from the REQ):

- Prose-only changes to shipped action files plus one user-guide paragraph. No Go changes, no new flags, no new action files, no SKILL.md routing rows.
- Keep every existing heading in `work-reference.md` unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings that name them.
- Sweep every restatement of a changed rule across `work.md`, `work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md` (lesson family alternate-writer-contract-drift).
- Read `_dev/primes/prime-action-files.md` before editing.
- No flag: the integrator is a placement, and the dispatch mechanism stays unspecified.

Orchestrator additions from exploration (inside the declared files, not new scope):
- The guardrail-slot table row for `manifest.md` in work-reference.md ("REQ id → builder, `<operative_name>`, handback file, landed status") restates what the manifest row carries. Requirement 2 makes the row carry the held dispatch instant, so name it in that table row too.
- Place the new Fan-Out Dispatch paragraph after the **Serial-only** paragraph and before **The run directory is mandatory here**, so 4(a)'s citation points backwards. Put the new table row next to the per-builder rows.
- Judge, and report in the hand-back, whether the work.md Orchestrator Checklist lines for Step 6 and Step 10, the work.md Rules bullet "the orchestrator is the sole integrator", and work-guide.md's **What isn't specified** paragraph ("two `do-work` sessions in the same working tree") need a word changed so they do not contradict the new rule. Change one only when it would otherwise be wrong; record each choice as a D-XX.
- Do not add a heading anywhere in work-reference.md; the new paragraph opens with a bold lead, like its neighbours.

## Write boundary
Only these three files, as the REQ's `## Scope` declares: skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/docs/work-guide.md. Never touch VERSION, CHANGELOG.md, skills/do-work/VERSION, skills/do-work/CHANGELOG.md, skills/do-work/actions/version.md, _dev/primes/lessons-action-files.md, do-work/lessons-index.md, crew-members/background-agents.md, review-work.md, any Go source, SKILL.md, or anything under do-work/. The integrator writes the release and the lessons entry. If the sweep shows another file must change, stop and say so in the hand-back; do not write it.

## Verify before hand-back
- From the worktree root, with wall times: `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh`. Both must pass. Baseline at 0efba253: 1s and 22s, both green.
- A cold-reader prose read of every changed section: read each one as an agent that has never seen this REQ, and confirm it can act on it without the REQ open (conditions stated, no forward references to text that does not exist, every citation names a real heading or bold label).
- A restatement grep for each changed rule across work.md, work-reference.md, docs/work-guide.md and crew-members/background-agents.md, at least for: landed hand-back / re-dispatch, dispatch instant / dispatch_at, advance --checkpoint / session-end writer, integrator / sole integrator / integration serial, run-with-recovery / --assume-sole-authority / --take-over, two sessions in one working tree, the guardrail-slot table rows. Paste each grep command and account for every hit (changed, or why it stays).
- `git diff --stat` and `git diff --check` on the branch.

## Hand-back (/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/87062c87-2ca6-4782-9a66-d1dde45f9c4a/scratchpad/REQ-639-handback.md)
- Branch name and every commit hash.
- File manifest: each file with (modified), and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY] (the orchestrator ticks the REQ's boxes from it); [UNIFY] lists git diff --stat and each file checked.
- `## Decisions` (D-01 onwards; DECIDE & STATE or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds; do not fix them).
- Lessons read (the family bullets above, and any satellite you opened).
- A proposed lesson entry for _dev/primes/lessons-action-files.md (one bullet, with a `[family: <slug>]` marker) and proposed CHANGELOG text (descriptive title, what shipped and why, plain words). The integrator writes both.
- Integration seams (none expected).
- The two test wall times and the restatement-grep accounting.
Final chat reply under 3000 characters.

## Coordinator note (observed on this run, 2026-10-07T19:14:45Z; record your handling as a D-XX in the hand-back)
Requirement 4(b) says the integrator enters through `recover` reporting the live claim. On this very run, plain `recover` refused with `FINALIZATION-DISCOVERY-AMBIGUOUS` on the untracked run directory (`do-work/runs/work-2026-10-07-191445/manifest.md`), and the REQ-637 recipe records the same refusal on the dirty `do-work/working/baseline.json` that pre-flight leaves. Both are the wave's own in-flight artifacts, which a mid-wave `recover` cannot attribute. The bare targeted `advance REQ-NNN` succeeded on the same tree and continues the claim. So write 4(b) as: the integrator does not run `recover`; the coordinator ran it at the queue boundary, and a mid-wave `recover` reads the wave's own in-flight artifacts as unattributed dirt. The integrator enters with the read-only `advance REQ-NNN`, which continues the live claim. Keep the never-rwr / never-take-over sentence. This is a DECIDE & STATE change (reversible, evidence above); do not reopen anything else.
