# Builder brief — REQ-640 (Mid-run messages: route a user message that arrives during a run to the right REQ in the user's own words without stopping the run)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-640-mid-run-messages-steer-a-run-in-progress
- Branch: worktree-agent-REQ-640-mid-run-messages-steer-a-run-in-progress (commit here; never merge, never rebase, never push)
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-640-mid-run-messages-steer-a-run-in-progress.md — read it fully: What, Why, Verified Facts, Detailed Requirements, Constraints, Builder Guidance, Red-Green Proof, and the orchestrator's `## Exploration` and `## Scope` (exact line numbers for every insertion point, each Verified Fact checked against the Go code, and every restatement found so far).
- Hand-back file (outside both trees, never commit it): /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/87062c87-2ca6-4782-9a66-d1dde45f9c4a/scratchpad/REQ-640-handback.md
- Route B, tdd: false, impact-rule-change. Prose-only. Depends on REQ-639 (archived at do-work/archive/REQ-639-delegated-integration-coordinator-shape.md, shipped as 0.305.74): its words are already in your base, so the new text must fit the **Delegated integration — the coordinator shape** paragraph in work-reference.md (line 470), Step 6's landed-hand-back condition and Step 10 in work.md, and the integrator brief `REQ-NNN-integrate.md` exactly as written there.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md (domain: general, so no extra domain crew). Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md (the REQ's prime_files; read it before editing). Lessons in its satellite /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-action-files.md (`slugged: partial` and over the lesson budget, so it is not in required_lessons):
- the six `[family: alternate-writer-contract-drift]` bullets, lines 52-60. Rule: a changed rule is not done until every writer and every action-bearing reader that restates it is swept, by ownership condition, not only the cited lines.
- the REQ-639 `[family: restated-mechanism-unchecked]` bullet, line 61. Rule: before shipped prose restates what a command does, read the command. The orchestrator already checked this REQ's three Verified Facts (see `## Exploration`); check any new mechanism claim you add the same way.

## The change (decided at capture, do not reopen)
Copied verbatim from the REQ's `## Detailed Requirements`:

1. `skills/do-work/actions/work.md`: new subsection `### Mid-Run Messages (any step)` directly after Step 3.5, beside its existing "question escalated mid-run and answered" paragraph. Condition-keyed, not a list of message types:
   - A question about progress or state is answered from disk (REQ files, run manifest, hand-backs), never from memory of the conversation, and the run does not stop.
   - A change that affects an in-flight REQ whose builder has not handed back: the durable record is `## Addendum (mid-run)` on the working REQ, holding the user's words under `actions/clarify.md` Step 4's Outside-text containment plus one framing line (extends, narrows, or corrects which requirement). In the serial loop the orchestrator writes it now. Under delegated integration (REQ-639, the coordinator shape) the coordinator puts the text in that REQ's integrator brief and the integrator writes the section before the merge (one writer under the project root). Forward it to the builder where the harness can message a running agent; otherwise say in the progress output that the builder will not see it and review will judge it. The section is user intent, so the takeover reset keeps it.
   - A change that affects an in-flight REQ whose hand-back already landed: capture's existing in-flight path (new REQ with `addendum_to`), reported as queued for the next loop.
   - A change to a queued REQ, or new work: capture's existing paths, run in the coordinator's writing gap.
   - A change of emphasis for the hand-back only (what to show first): write it to the run manifest, not memory; the Decision Brief reads it.
   - Two rules across all branches: forward the user's words, never a paraphrase, because meaning is lost in restatement; route the message to every REQ it affects, not only the one the user named.
2. `skills/do-work/actions/review-work.md` Step 5 item 1: extract requirements from What/Detailed Requirements, any `## Addendum` section, and the UR.
3. `skills/do-work/actions/capture.md` Step 2 table, `do-work/working/` row: keep "NEVER modify" for capture; add one pointer clause: a message that arrives during a live run in this session is routed by `actions/work.md` -> Mid-Run Messages first.
4. `work.md` Orchestrator Checklist: one line "Mid-run message: route per Mid-Run Messages; never stop the run for it."

Constraints (verbatim from the REQ):

- Prose-only changes to shipped action files. No Go changes, no new flags, no new action files, no SKILL.md routing rows.
- Keep every existing heading in `work-reference.md` unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings that name them.
- Sweep every restatement of a changed rule across `work.md`, `work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md` (lesson family alternate-writer-contract-drift).
- Read `_dev/primes/prime-action-files.md` before editing.
- The subsection is keyed on conditions, not a closed list of message types.

Orchestrator additions from exploration (they make the decided text true; they do not reopen it):
- **Exact insertion points (at base 31a6e35f):** work.md Step 3.5 is lines 137-163, the escalated-answer paragraph is line 155, and the new `### Mid-Run Messages (any step)` goes after line 163 and before `### Mechanical Evidence-Gate Loop` (line 165). review-work.md Step 5 item 1 is line 87. capture.md's `do-work/working/` table row is line 124. The work.md Orchestrator Checklist is lines 476-500; put the new line next to the Step 3.5 line (483) or wherever a reader scanning per step finds it, your call.
- **One section, never two.** The advance classifier refuses any section name that appears twice (`internal/lifecycleadvance/advance_commands.go` line 356, "duplicate lifecycle section"). Say that a later mid-run message for the same REQ appends inside the existing `## Addendum (mid-run)` section; never a second heading with the same name. Write it as the condition, not as a reference to the Go file.
- **"The Decision Brief reads it" is not true yet; make it true.** `actions/work-reference.md` → **Decision Brief (hand-back format)** (lines 846-868) names no run-manifest source today, and the guardrail-slot table row for `manifest.md` (line 480) lists what the manifest carries without an emphasis note. Add one clause to the Decision Brief section saying the hand-back reads the run manifest's emphasis note and orders what it shows by it, and name that note in the `manifest.md` table row. Keep both headings and the table's other rows as they are.
- **A run with no run manifest.** A run directory exists only for fan-out or a background builder (`crew-members/background-agents.md`; work-reference.md "The run directory is mandatory here"). State what the emphasis branch does when the run has none, as a condition. Recommendation: the run manifest when one exists; otherwise create the run directory and its `manifest.md` for that note, since background-agents.md already treats it as an ordinary committable path. Record your choice as a D-XX with Value and Risk.
- **A message for a REQ whose integrator is already running.** Under delegated integration the coordinator writes nothing under the project root while an integrator runs, and that REQ's brief already exists. State where the text goes then without breaking the one-writer rule. Recommendation: forward the user's words to the running integrator where the harness can message it (the integrator writes the section before its merge); otherwise the coordinator holds the words in its session scratch, and once the hand-back has merged the landed-hand-back branch applies. Record as a D-XX.
- **Citation form.** capture.md's pointer must cite `actions/work.md` → **Mid-Run Messages (any step)** in the house arrow form with the exact heading text, so the citation contract test resolves it. The REQ's ASCII "->" is the capture shorthand, not the shipped form.
- **Same rule, other readers in your Scope (judge each; change only when it would otherwise be wrong; record each choice as a D-XX):**
  - review-work.md Step 2 line 55 ("What was requested — the What/Detailed Requirements sections") restates requirement 2's source list.
  - capture.md **Immutability Rule** (line 63) says an addition to an in-flight request always becomes a new addendum REQ; a capture reader in a live run reads it before the table.
- **Swept and consistent (no change expected, list them in your restatement accounting):** work-reference.md line 51 ("working/: Immutable to all actions except the work pipeline"), work-reference.md line 470 (delegated paragraph; the brief already carries "any mid-run addendum text"), work.md line 135 and Progress Reporting (lines 519-527), docs/work-guide.md lines 72 and 101, docs/capture-guide.md line 76, crew-members/background-agents.md.
- **Containment detail you can rely on:** clarify.md lines 101-106 put a body passage in a blockquote whose lines open a code fence longer than the longest backtick run. The section parser skips fenced lines, so a `## ` line inside the user's words cannot open a section. Cite the contract by name; do not restate its mechanics.

## Write boundary
Only these four files, as the REQ's `## Scope` declares: skills/do-work/actions/work.md, skills/do-work/actions/review-work.md, skills/do-work/actions/capture.md, skills/do-work/actions/work-reference.md. Never touch VERSION, CHANGELOG.md, skills/do-work/VERSION, skills/do-work/CHANGELOG.md, skills/do-work/actions/version.md, _dev/primes/lessons-action-files.md, do-work/lessons-index.md, docs/work-guide.md, docs/capture-guide.md, crew-members/background-agents.md, actions/clarify.md, any Go source or test, SKILL.md, or anything under do-work/. The integrator writes the release and the lessons entry. If the sweep shows another file must change, stop and say so in the hand-back; do not write it.

## Verify before hand-back
- From the worktree root, with wall times: `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh`. Both must pass. Baseline at 31a6e35f: 1s and 21s, both green.
- A cold-reader prose read of every changed section: read each one as an agent that has never seen this REQ, and confirm it can act on it without the REQ open (conditions stated, not a list of message types; no forward references to text that does not exist; every citation names a real heading or bold label).
- A restatement grep for each changed rule across work.md, work-reference.md, review-work.md, capture.md, clarify.md and docs/work-guide.md, at least for: mid-run / mid run, `## Addendum` / Addendum section / addendum_to, NEVER modify / immutable / Immutability Rule, integrator brief / REQ-NNN-integrate.md / mid-run addendum text, run manifest / manifest.md, Decision Brief, What/Detailed Requirements, paraphrase / verbatim. Paste each grep command and account for every hit (changed, or why it stays).
- `git diff --stat` and `git diff --check` on the branch.

## Hand-back (/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/87062c87-2ca6-4782-9a66-d1dde45f9c4a/scratchpad/REQ-640-handback.md)
- Branch name and every commit hash.
- File manifest: each file with (modified), and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY] (the orchestrator ticks the REQ's boxes from it); [UNIFY] lists git diff --stat and each file checked.
- `## Decisions` (D-01 onwards; DECIDE & STATE or ESCALATE with Value and Risk). At least the no-run-manifest choice and the running-integrator choice.
- `## Discovered Tasks` (out-of-scope finds; do not fix them). Candidate: a user-facing sentence in docs/work-guide.md beside line 72 saying a message sent during a run reaches the in-flight REQ.
- Lessons read (the family bullets above, and any satellite you opened).
- A proposed lesson entry for _dev/primes/lessons-action-files.md (one bullet, with a `[family: <slug>]` marker) and proposed CHANGELOG text (descriptive title, what shipped and why, plain words). The integrator writes both.
- Integration seams (none expected; REQ-641, the wave-end consistency check, comes next and touches the same files).
- The two test wall times and the restatement-grep accounting.
Final chat reply under 3000 characters.

## Coordinator note (answers to the two open choices; record your handling as D-XX in the hand-back)
1. No run manifest: key it on the condition, not on a mode. Where a run directory with a manifest exists (every fan-out run, and any serial run that created one), the hand-back emphasis note goes into the manifest. Where none exists, the same session that received the message renders the Decision Brief, so it holds the note itself and says so in one progress line; do not create machinery (no new file, no new directory) for that case.
2. A message for a REQ whose integrator is already running: that REQ's hand-back has landed, so it is the "hand-back already landed" branch: capture's existing in-flight path (new REQ with addendum_to), executed by the coordinator in its next writing gap, never while the integrator runs. Until the gap the coordinator keeps the user's words in its own session scratch outside the project root. Say this in the subsection in one sentence so a reader does not invent a third path.
3. Second message for the same in-flight REQ: append inside the existing `## Addendum (mid-run)` section, each entry under its own timestamp line, because advance refuses a duplicated section heading (your pre-dispatch finding). State it as the condition "one section per REQ, entries inside".
These are DECIDE & STATE (reversible prose, evidence in the REQ and this run).
