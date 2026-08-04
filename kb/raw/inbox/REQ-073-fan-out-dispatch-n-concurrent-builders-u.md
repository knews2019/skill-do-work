---
source_type: req_lesson
req_id: REQ-073
req_path: do-work/archive/UR-013/REQ-073-fan-out-dispatch-n-builders-one-owner.md
date: 2026-08-03
domain: general
module: actions
tags: [actions, fan-out, dispatch, concurrent, builders]
---

# Lessons from REQ-073: Fan-out dispatch: N concurrent builders under one queue owner

## What the REQ was about

Raise Worktree Dispatch Mode from one builder to several concurrent builders under a single queue
owner. The exclusive-session contract is about **who owns `do-work/`**, not about how many builds run,
so a builder that owns nothing shared needs no coordination at all.

## Solution summary

Raised worktree dispatch from one builder to N under a single queue owner, by rewording rather than building. The Execution Model invariant became **one queue owner per checkout** — the session that claims, flips status, and archives — with builders explicitly not owners, any number allowed to build concurrently, and the two-queue-owners ban kept explicit so none of REQ-069's deleted machinery has a reason to return; the paragraph came in at 198 of its 200-word cap. Both capping sentences in Worktree Dispatch Mode's opening are gone, replaced by the fact that makes the lift safe: everything in that section is already written per REQ (one `<operative_name>`, one hand-back, one `<pre>..<merge_hash>`, one cleanup — each per REQ), so it already holds for any builder count. A new **Fan-Out Dispatch** subsection adds only what fan-out genuinely introduces: a human picks the set; `write_set` is advisory input to that pick and never a gate; the merge is the non-interference proof, with its honest limit stated (git detects conflicts by line proximity, not meaning, so two REQs appending to a shared registry merge cleanly and can still be jointly wrong — the integration-seam rule covers that, and only under one integrator); integration is serial and each merge invalidates the previous post-merge verification; a worktree per builder is mandatory, with the ruled-out shared-tree rationale kept in shipped prose; a serial-only list naming queue transitions, REQ id allocation, and `actions/version.md` + `CHANGELOG.md`; the mandatory run-directory mapping from `crew-members/background-agents.md` with its survivable-not-prevented ceiling carried; the brief-delivery path trap; and deliberate silence on dispatch mechanism. The contract suite's exactly-once invariant assertion was repointed at the new token and hardened against a silent-abort bug, plus eight new assertions. `CLAUDE.md`'s Before-Every-Commit ritual is now scoped to the integrating commit so a worktree builder skips it, and `docs/work-guide.md`'s "One REQ at a time" bullet is true for ownership and no longer false for builders.

## What worked

- **Reading the target section for what was already per-REQ before planning any addition.** Eleven requirements looked like a large build; the section turned out to already specify `<operative_name>`, the hand-back, the merge range, verification, and cleanup *per REQ*, so the "cap" was literally two sentences. That observation became the replacement sentence and is why this landed in 29 net lines. The generalizable move: before adding a concurrency story, check whether the existing prose is already written per-unit — if it is, the cap is a claim, not an architecture.
- **Measuring the word budget before drafting.** The Execution Model cap is ≤200 and the section was at 158, so the invariant rewrite had ~40 words. Knowing that first produced a boundary sentence; discovering it after would have produced a good paragraph that then had to be mangled. Two trim passes were still needed (201 → 201 → 198), which is a fair sign the cap is doing real work.

## What didn't work

- **The first version of my own assertion was worse than no assertion.** `grep -roh <pattern> | wc -l` under `set -euo pipefail` aborts the suite when the pattern is absent — so the check for "the invariant exists" exited 1 with **no output whatsoever** in precisely the case it was written to catch. It read as a crash, not a finding. Any counter-style assertion in a `pipefail` script needs `{ … || true; }`, and the tell is that the assertion's own failure mode is silence.
- **Writing the decision down is not the same as updating the sweep's input** — the same lesson REQ-072 hit from the other side. Here the conditional Scope entry ("`actions/work.md` — only if the sweep finds a restatement") worked, because the condition was written into the declaration where `scope-drift.sh` could see the file listed. A conditional declared *only* in prose would have failed the check.
- **My first sweep pass looked for the invariant's string and stopped.** The eight-site `write_set` family only surfaced on a second pass that searched for *consumers of the reasoning* rather than copies of the wording. A sweep that greps the changed token finds restatements; a sweep that asks "what argued from this premise?" finds the dangerous ones. The five in REQ-075 share no phrase with anything I edited.

## Worth knowing

- **The premise, not the conclusion, is what goes stale.** Every one of the eight sweep findings kept a true conclusion ("nothing schedules on `write_set`") attached to a premise that had just become false ("one REQ runs at a time"). That shape is more dangerous than a plainly wrong sentence, because a careful reader *reasons forward* from the dead premise and lands on the opposite of the contract. When a change falsifies a premise, grep for what was justified by it, not just for restatements of it.
- **`CLAUDE.md` auto-loads in a worktree of this repo.** That is why requirement 10 insisted the *instruction* be fixed rather than each builder's brief. It also means this file's rules reach agents nobody briefed — worth remembering before adding any imperative here.
- **The one-owner/N-builders split is now the load-bearing distinction, and the heading doesn't say it.** "Execution Model — Exclusive Session" is still accurate (the session is exclusive; the builds are not), and the paragraph makes that explicit, but the mismatch is the first thing a fresh reader will trip on. If it comes up twice, rename the heading and sweep its seven pointers — the cost is known and small, and I chose not to spend it without a requirement.
- **REQ-071's takeover prompt and this REQ's ownership boundary now read as one position** — which the batch's REQ-071 flagged as needing a check. Recovery spends prose on a claim it cannot account for, and the invariant now says why that is coherent rather than contradictory: a foreign claim is evidence the *ownership* rule was violated, and the pipeline's answer is to report and ask rather than to coordinate. No edit was needed; the reworded invariant absorbed it.

## Back-reference

See `do-work/archive/UR-013/REQ-073-fan-out-dispatch-n-builders-one-owner.md` for the full REQ — triage, implementation, review, and lessons. Commit `9ba2cda`.
