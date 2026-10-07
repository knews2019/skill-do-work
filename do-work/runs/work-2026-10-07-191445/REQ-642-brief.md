# Builder brief — REQ-642 (Addendum to REQ-641: the wave-end decision reads the set-aside list, and a late set-aside is named in the Decision Brief)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-642-wave-end-set-aside-gaps
- Branch: worktree-agent-REQ-642-wave-end-set-aside-gaps, created from HEAD d2d3c316 with `git worktree add -b` (commit here; never merge, never rebase, never push).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-642-wave-end-set-aside-gaps.md. Read it fully: What, Why, Prior Implementation, Detailed Requirements, Constraints, Builder Guidance, Red-Green Proof.
- Original REQ (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/UR-138/REQ-641-wave-end-consistency-check.md, shipped as 0.305.76. Its words are already in your base.
- Hand-back file (outside both trees, never commit it): /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/87062c87-2ca6-4782-9a66-d1dde45f9c4a/scratchpad/REQ-642-handback.md
- Route A, tdd: false, impact-rule-change, effort-mechanical. Prose-only, two shipped files.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md (domain: general, so no extra domain crew). Prime: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md (read it before editing). Lessons: the `[family: alternate-writer-contract-drift]` bullets in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-action-files.md (lines 52-60 and 62-63). Rule: a changed rule is not done until every writer and every action-bearing reader that restates it is swept; grep the restated phrase, not only the cited line. REQ-641's bullet (line 63) adds: a rule that tells an agent to read a state must name the writer of that state.

## The change (decided at capture, do not reopen)
From the REQ's `## Detailed Requirements`, condition-keyed, no new mechanism:
1. `skills/do-work/actions/work-reference.md`, paragraph **Delegated integration — the coordinator shape** (line 470). Its last sentence lists what the integrator brief `REQ-NNN-integrate.md` carries ("the REQ id, the run directory, the hand-back path, the operative name, the wave membership, any mid-run addendum text, and the instruction not to run Step 10"). The brief also names every wave member the coordinator has set aside, because a set-aside member keeps its claim in `do-work/working/` and otherwise looks like one still building. One sentence or one clause in that list.
2. `skills/do-work/actions/work.md` Step 7, **How to run it** (line 370), the wave-end clause: "The run manifest's rows name the wave's members, but they do not track who is done: the archive shows which members finalized, and the orchestrator knows which it set aside." Make it say the orchestrator decides "last successful integration" from the archive plus the set-aside list it holds, which under delegated integration is the list the integrator brief carries. One clause.
3. Same paragraph, F2, say it rather than build for it: when a member is set aside after the wave's last review already ran, the wave-end sweep did not run for that wave. The orchestrator names that in the run's Decision Brief under HANDLED (`actions/work-reference.md` → **Decision Brief (hand-back format)**) so the gap is visible, and no review is re-run. One sentence.
4. Sweep every restatement of the wave-end rule; change none unless it contradicts the new sentences.

Constraints: prose only, only the two files below; keep every existing heading and bold label unchanged (`_dev/tests/shipped-package-reference-contract.sh` pins citation strings). Builder Guidance: latitude is limited to wording.

## Orchestrator notes (they make the decided text true; they do not reopen it)
- **The HANDLED bullet** (work-reference.md line 871) says HANDLED "lists the **DECIDE & STATE** decisions (reversible `D-NN` entries)". Requirement 3 puts a set-aside gap there. Judge whether the new sentence contradicts that bullet. The orchestrator's choice not to re-run a review is a reversible call made without asking, so phrasing the note as that call should fit HANDLED's existing rule without editing the bullet. If you conclude the bullet must change, stop and say so in the hand-back instead of editing it (it is outside the REQ's named clauses). Record your judgment as a D-XX.
- **Where the requirement-2 clause lands.** Line 368 (**Restatement sweep (MUST)**) says "(passed below)" and points at **How to run it**; keep the decision in **How to run it** and do not restate it at 368.
- **Restatement sweep targets (requirement 4), judged at HEAD d2d3c316:**
  - review-work.md Step 6 **Restatement Sweep**, step 1 "At a wave end, the trigger set also inherits" (line 135): the reviewer is passed the fact; it never decides "last". No set-aside wording. No change expected.
  - work-reference.md Fan-Out Dispatch bullet "The non-interference proof is the merge, not the pick" (line 456): names the wave-end sweep. No change expected.
  - work-reference.md guardrail table row "per-integrator input" (line 479): points at the coordinator-shape paragraph. No change expected.
  - work.md **Mid-Run Messages (any step)** (line 170): mentions the integrator brief carrying addendum text. No change expected.
  - docs/work-guide.md line 134: user-facing wave sentence. No change expected; it is outside your write boundary.
  Grep at least: set aside / set-aside, last successful, wave membership, integrate.md / integrator brief, Decision Brief / HANDLED, wave-end / wave end. Account for every hit.

## Write boundary
Only skills/do-work/actions/work-reference.md and skills/do-work/actions/work.md (the REQ's `write_set`; Route A has no `## Scope`). Never touch VERSION, CHANGELOG.md, skills/do-work/VERSION, skills/do-work/CHANGELOG.md, skills/do-work/actions/version.md, review-work.md, docs/, crew-members/, the lessons satellite, do-work/lessons-index.md, any Go source or test, SKILL.md, or anything under do-work/. If the sweep shows another file must change, stop and say so in the hand-back; do not write it.

## Verify before hand-back
- From the worktree root, with wall times: `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh`. Both must pass. Baseline at d2d3c316: green on tonight's REQ-641 run.
- A cold-reader read of each changed paragraph: an agent that has never seen this REQ can act on it without the REQ open (who writes the set-aside list, who reads it, where the F2 note goes, that no review is re-run).
- The restatement grep above, each command pasted, every hit accounted for.
- `git diff --stat` and `git diff --check` on the branch.

## Hand-back
- Branch name, worktree path and every commit hash.
- File manifest: each file with (modified) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY] (the orchestrator ticks the REQ's boxes from it); [UNIFY] lists git diff --stat and each file checked.
- `## Decisions` (D-01 onwards; DECIDE & STATE or ESCALATE with Value and Risk), including the HANDLED judgment.
- `## Discovered Tasks` (out-of-scope finds; do not fix them).
- Lessons read.
- A proposed lesson bullet for _dev/primes/lessons-action-files.md with a `[family: <slug>]` marker, or "none" with a reason, and a proposed CHANGELOG entry (descriptive title saying what shipped, plain words). The orchestrator writes both.
- Integration seams (none expected).
- The two test wall times and the restatement-grep accounting.
