# Builder brief — REQ-641 (Wave-end consistency check: the last review in a fan-out wave also sweeps for elements earlier members redefined)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-641-wave-end-consistency-check
- Branch: worktree-agent-REQ-641-wave-end-consistency-check (commit here; never merge, never rebase, never push)
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-641-wave-end-consistency-check.md — read it fully: What, Why, Detailed Requirements, Constraints, Builder Guidance, Red-Green Proof, and the orchestrator's `## Exploration` and `## Scope` (exact line numbers for every insertion point, the one mechanism claim checked against Go, and every restatement found so far).
- Hand-back file (outside both trees, never commit it): /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/87062c87-2ca6-4782-9a66-d1dde45f9c4a/scratchpad/REQ-641-handback.md
- Route B, tdd: false, impact-rule-change. Prose-only. Depends on REQ-640 (archived at do-work/archive/REQ-640-mid-run-messages-steer-a-run-in-progress.md, shipped as 0.305.75), which depends on REQ-639 (do-work/archive/REQ-639-delegated-integration-coordinator-shape.md, 0.305.74). Their words are already in your base, so the new text must fit them exactly as written: work-reference.md **Delegated integration — the coordinator shape** (line 470; the integrator brief carries "the wave membership"), the guardrail table's `manifest.md` row (line 480: REQ id → builder, operative name, handback file, landed status, held dispatch instant, emphasis note), and review-work.md's **Append to REQ File** template (lines 380-408).

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md (domain: general, so no extra domain crew). Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md (the REQ's prime_files; read it before editing). Lessons in its satellite /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-action-files.md (`slugged: partial` and over the lesson budget, so it is not in required_lessons):
- the `[family: alternate-writer-contract-drift]` bullets, lines 52-60 and 62. Rule: a changed rule is not done until every writer and every action-bearing reader that restates it is swept, by ownership condition, not only the cited lines. REQ-640's bullet (line 62): grep the restated phrase, not only the cited line, and read every sibling branch for a shared clause.
- the REQ-639 `[family: restated-mechanism-unchecked]` bullet, line 61. Rule: before shipped prose restates what a command does, read the command. The orchestrator checked the REQ's one mechanism claim (see `## Exploration`); check any new mechanism claim you add the same way.

## The change (decided at capture, do not reopen)
Copied verbatim from the REQ's `## Detailed Requirements`:

1. `skills/do-work/actions/review-work.md` Step 6 Restatement Sweep:
   - Record the result in the appended `## Review` block as one line: `**Restatement sweep:** redefined <elements> | nothing redefined`.
   - Add the line to the "Append to REQ File" template, and point the existing Verification Checklist item at it.
   - Go checks section presence only, so the extra line is safe.
   - New condition in the sweep's trigger (step 1): when the REQ under review is the wave's last successful integration (wave membership from the run manifest; a set-aside member does not skip the check), the trigger set also includes every element the wave's earlier members recorded as redefined, read from their archived `## Review` lines.
   - One-sentence reason in the text: an earlier member's review ran before the later member merged, so only the last review sees both.
   - Findings route exactly as today (impact token, `→ report only` unless `impact-critical`).
2. `skills/do-work/actions/work-reference.md` Fan-Out Dispatch bullet "The non-interference proof is the merge, not the pick": one sentence pointing at the wave-end sweep as the semantic half the merge cannot provide.

Constraints (verbatim from the REQ):

- Prose-only changes to shipped action files. No Go changes, no new flags, no new action files, no SKILL.md routing rows.
- Keep every existing heading in `work-reference.md` unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings that name them.
- Sweep every restatement of a changed rule across `work.md`, `work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md` (lesson family alternate-writer-contract-drift).
- Read `_dev/primes/prime-action-files.md` before editing.
- No new test: no existing test exercises review prose.

Builder Guidance from the REQ: certainty is high; latitude is limited to wording; finding routing does not change.

Orchestrator additions from exploration (they make the decided text true; they do not reopen it):
- **Exact insertion points (at base 2291812d, before your edits):**
  - review-work.md **Restatement Sweep** paragraph: heading line 130, step 1 (Trigger) line 134, step 2 line 135, step 3 (routing, leave as is) line 136, step 4 (Skip it when nothing was redefined) line 137, Origin line 139.
  - review-work.md **Append to REQ File**: line 380; the fenced template runs lines 384-408. Put the new bold line with the other bold record lines (Minor findings, Acceptance, Suggested testing, Follow-ups created at lines 402-405).
  - review-work.md **Verification Checklist**: line 484; the sweep item is line 489 and quotes the old N/A wording ("nothing redefined, sweep N/A"). Point it at the template line and its `nothing redefined` form.
  - work-reference.md Fan-Out Dispatch bullet "The non-interference proof is the merge, not the pick": line 456. It already says git detects conflicts "by line proximity, not meaning" and points at the integration-seam rule. Add the one sentence there.
- **Step 4 must agree with the new condition.** "Skip it when nothing was redefined" (line 137) would let a last-of-wave review skip when its own diff redefines nothing even though earlier members did. Make step 4 say the skip applies only when the whole trigger set, including inherited elements, is empty.
- **The record line holds this REQ's own redefinitions.** A later member reads it as "what this member redefined", so inherited elements a last-of-wave review swept do not belong in it; their stale restatements are ordinary findings. State this in one clause, or record your alternative as a D-XX.
- **work.md joins Scope (the coordinator already mirrored it into write_set).** Step 7 **Restatement sweep (MUST)** (work.md line 368) restates the trigger as "If this REQ's diff redefines something other text restates"; bring it in line with the widened trigger in one clause and cite review-work.md, do not copy the condition. Step 7 **How to run it** (line 370) lists what the review agent receives and gives it no run manifest; in a fan-out run the reviewer of the last member needs the run directory (manifest path) to apply the condition. Add it there, condition-keyed ("when the run has a run manifest").
- **Open choices; state each as a D-XX with Value and Risk (DECIDE & STATE, reversible prose):**
  1. *How a review knows it is the wave's last successful integration.* Recommendation: every other member of its wave in the run manifest is already finalized or set aside, so no member still waits to integrate. Wave membership comes from the run manifest; under delegated integration the integrator brief also carries it (work-reference.md line 470). One manifest may span several waves (background-agents.md step 3 says "a manifest per wave", but a run like this one keeps one manifest for a serial chain), so say the wave is the set the manifest or the integrator brief names, not the whole run.
  2. *Where the condition does not apply.* Recommendation: a run with no run manifest (the serial in-session loop, standalone `do-work review`) has no wave, and the trigger stays as today. State it as the condition, not as a list of modes.
  3. *An earlier member with no recorded line* (archived before this shipped, or its line missing). Recommendation: say it in the record or findings as not recorded rather than treating it as nothing redefined, because absence is unknown, not empty; do not add machinery to derive it.
  4. *Where earlier members' reviews live.* They are archived before the next member integrates (Fan-Out Dispatch "Integration is serial", line 457), flat at `do-work/archive/REQ-NNN-*.md` or under `do-work/archive/UR-NNN/` when the UR closed; a set-aside member was never archived and has no line. Name "the earlier members' archived REQ files" without hardcoding one path form.
- **Same rule, other readers (judge each; change only when it would otherwise be wrong; list each in your restatement accounting):**
  - work-reference.md line 462 (Auto-wave: "The merge is the non-interference proof, not the pick (above)"): points back at line 456; no change expected.
  - work-reference.md line 470 (Delegated integration): carries the wave membership already; no change expected.
  - docs/work-guide.md line 134 ("collisions surface when the branches merge"): user-facing, still true for textual collisions. Outside your write boundary; if you judge a sentence on meaning collisions is needed, report it as a discovered task.
  - docs/review-work-guide.md (Phase 2 table, Follow-ups line 55): never states the trigger; no change expected.
  - crew-members/background-agents.md: nothing on review; no change expected.
  - review-work.md line 471 (Common Rationalizations on out-of-scope stale restatements): routing unchanged.

## Write boundary
Only these three files, as the REQ's `## Scope` declares: skills/do-work/actions/review-work.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/work.md. Never touch VERSION, CHANGELOG.md, skills/do-work/VERSION, skills/do-work/CHANGELOG.md, skills/do-work/actions/version.md, _dev/primes/lessons-action-files.md, do-work/lessons-index.md, docs/work-guide.md, docs/review-work-guide.md, crew-members/background-agents.md, any Go source or test, SKILL.md, or anything under do-work/. Never change a heading's text in work-reference.md. The integrator writes the release and the lessons entry. If the sweep shows another file must change, stop and say so in the hand-back; do not write it.

## Verify before hand-back
- From the worktree root, with wall times: `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh`. Both must pass. Baseline at 2291812d: green (pre-flight probe do-work/runs/work-2026-10-07-191445/REQ-641-probe.sh).
- A cold-reader prose read of every changed section: read each one as an agent that has never seen this REQ, and confirm it can act on it without the REQ open (the condition stated, where the manifest and the earlier reviews are found, what the record line holds, no forward references to text that does not exist, every citation names a real heading or bold label).
- A restatement grep for each changed rule across review-work.md, work.md, work-reference.md, docs/review-work-guide.md and docs/work-guide.md, at least for: Restatement Sweep / restatement sweep / Restatement sweep, redefin, nothing redefined / sweep N/A, non-interference, line proximity / not meaning, wave membership / last successful / set aside / set-aside, run manifest / manifest.md, `## Review` / Append to REQ File, report only. Paste each grep command and account for every hit (changed, or why it stays).
- `git diff --stat` and `git diff --check` on the branch.

## Hand-back (/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/87062c87-2ca6-4782-9a66-d1dde45f9c4a/scratchpad/REQ-641-handback.md)
- Branch name and every commit hash.
- File manifest: each file with (modified), and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY] (the orchestrator ticks the REQ's boxes from it); [UNIFY] lists git diff --stat and each file checked.
- `## Decisions` (D-01 onwards; DECIDE & STATE or ESCALATE with Value and Risk). At least the four open choices above and the record-line content.
- `## Discovered Tasks` (out-of-scope finds; do not fix them). Candidate: a user-facing sentence in docs/work-guide.md beside line 134 saying the last review in a wave also checks for meaning collisions.
- Lessons read (the family bullets above, and any satellite you opened).
- A proposed lesson entry for _dev/primes/lessons-action-files.md (one bullet, with a `[family: <slug>]` marker) and proposed CHANGELOG text (descriptive title, what shipped and why, plain words). The integrator writes both.
- Integration seams (none expected; this is the last REQ of the UR-138 chain).
- The two test wall times and the restatement-grep accounting.
Final chat reply under 3000 characters.
