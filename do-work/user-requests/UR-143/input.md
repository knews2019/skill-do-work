---
id: UR-143
title: 'Act on the October release review: disk probe to repo root, fan-out prose to a companion reference, operator work rule, activity-correlation report'
created_at: 2026-10-09T16:06:43Z
requests: [REQ-650, REQ-651, REQ-652, REQ-653]
word_count: 2398
---
# Four Accepted Items From the 0.305.60 to 0.305.84 Review Triage

## Summary
The maintainer ran `do-work-toolbox validate-feedback` on a consumer-side review of the 25 October releases and on an upstream proposal titled "operator work is not a REQ", asked for the ask tool to choose how to simplify, answered four questions, approved the plan, and then refined the operator rule twice. Every claim was verified against 0.305.84 and git history. Four items produce work: the disk probe measures the repo root only (F8); an ai-report explains which lines of the board's activity correlation are deterministic git structure and which are scaffold around agent behaviour, so the maintainer can decide its fate afterwards (F7); routine operator work is never captured and a tracked special operator configuration is a blocked, dependency-gated REQ that clarify completes into the Done column (F11, as the maintainer reshaped it); and the fan-out orchestration prose is compressed in place and moved to a companion reference with no size limit (F6). F1 to F5, F9 and F10 produce nothing: F5 is release-granularity policy and F9/F10 are consumer-local.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-650 | The disk probe measures the repo root only |
| REQ-651 | AI report on the board activity correlation separates deterministic git structure from scaffold around agent behaviour |
| REQ-652 | [impact-rule-change] Routine operator work is never a REQ; a tracked operator configuration is a blocked, dependency-gated REQ that clarify completes |
| REQ-653 | [impact-rule-change] Fan-out orchestration prose is compressed in place and moved to a companion reference |

## Batch Constraints
- Order: REQ-650 and REQ-651 are independent of everything. REQ-653 depends on REQ-652 because both edit `actions/work.md` and `actions/work-reference.md` and the compression must include the operator sentences; this is the approved plan order (W1, W4, then W3, then W2).
- Each REQ is its own release per `_dev/primes/prime-releases.md`; do not fold one REQ's change into another. REQ-651 changes no shipped file and is not a release.
- Read `_dev/primes/prime-kanban-board.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` before REQ-650; `_dev/primes/prime-action-files.md` before REQ-652 and REQ-653.
- Pushed back and NOT captured: file mtime as the activity source (F7 remedy), any change to the activity correlation code (decision deferred until REQ-651's report is read), the proposal's "no REQ, report-only line, recommend abandon" shape for operator work (F11), and any action on F1 to F5, F9, F10.
- Surface-cost: REQ-650 and REQ-653 are deletions or simplifications (N/A). REQ-652 adds one capture rule and one clarify option, earned by one consumer REQ that cost 53 commits, six days and three block/clarify round trips.

## Full Verbatim Input
> ```
> do-work validate-feedback: 
> 
> ### The Story
> The "do-work suite"—the set of tools used to track and execute tasks—just underwent a massive update. Between early and mid-October 2026, the maintainers pushed 25 different patch releases. This wasn't one big change, but a series of small tweaks that have now been bundled into the local environment.
> 
> **The Good News:**
> There are some "must-have" wins here. Previously, if you tried to resume a handoff (passing work from one session to another), the system could accidentally wipe out merged claims or delete parts of your plan. Those bugs are gone. Additionally, the system used to completely crash during commits if a single archive file had a typo in its backticks; now, it simply warns you and skips that file. They also added a disk space probe—a safety check that prevents the system from crashing when large-scale browser tests fill up the hard drive.
> 
> **The Noise:**
> Not everything was useful. There is a fair amount of "fluff." For instance, the developers spent five separate releases just to clarify what one single status field means. More concerningly, they added a lot of dense, wordy instructions for the **Orchestrating Model** (the AI that manages the workflow). Because these instructions are written as long prose in files that are already huge, there is a risk that the AI will struggle to remember everything or get bogged down by the sheer volume of text.
> 
> **The Overkill:**
> Some features were built with far more complexity than necessary. The most glaring example is "Activity Correlation." Instead of just looking at the last time a task file was saved, the developers wrote hundreds of lines of code to analyze Git history, merge branches, and track ancestry just to put a date on a dashboard card. Similarly, the disk space check is over-engineered; it treats the computer like a complex network of different volumes when a simple check of the project folder would have sufficed.
> 
> **Local Cleanup:**
> Aside from the suite update, a few small things happened locally: some helpful shortcut commands were accidentally deleted from `CLAUDE.md` (the project's primary AI instruction file), and a stale marker for **REQ-198** (a specific task request ticket) was removed.
> 
> --- Mid-turn message from the maintainer ---
> 
> use the ask tool to help me clarify and choose how to simplify and to refactor the fluff and complexity
> 
> --- Second paste from the maintainer ---
> 
> see aditional context to be evaluated: 
> 
> # Upstream suggestion for `knews2019/skill-do-work` — operator work is not a REQ
> 
> **How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
> clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
> consumer repo it was authored in. Observed against **v0.305.84**; all line numbers are from that tag.
> 
> ---
> 
> ## Request
> 
> Capture and the run orchestrator both know one answer for "the user has to be present as operator": write
> `status: blocked` with `blocked_by` naming the person. That answer is right when the operator is a
> **precondition** (credentials arrive, a service comes up, a person answers) and AI work follows. It is wrong
> when the operator act **is** the work: deploying, publishing, verifying on live hosts, approving content,
> waiting for a calendar day. A REQ of that shape has no path to `completed`. The pipeline can only prepare,
> block, be released, re-block and record receipts, and the user ends the cycle with `do-work abandon`.
> 
> Please add two sentences, no new fields, no new status:
> 
> 1. **Capture, External-condition assessment** (`actions/capture.md:108`): before writing `blocked`, ask
>    whether any AI-buildable work remains once the condition clears. If the remaining work is the operator's
>    own hands (deploy, publish, verify on hosts, approve content, wait for a date), do not mint a REQ. Write
>    the runbook or checklist as a report-only line (`actions/capture-reference.md:34` already lets a
>    maintainer promote a report-only line later) and capture only the tooling gaps the AI must build as REQs.
> 2. **Work, mid-run flips** (`actions/work.md:522-523` and the Environment row at
>    `actions/work-reference.md:788`): add the inverse of the blocked-flip test. If a run finds that nothing
>    AI-buildable remains and the rest is operator hands, do not flip to `blocked`. Say so in the hand-back,
>    point at the runbook or report where receipts should go, and recommend
>    `do-work abandon REQ-NNN <reason>` with the reason pre-filled.
> 
> Nothing here changes queue semantics: `blocked` stays an external precondition, `abandon` stays an explicit
> user act, and the report-only line is the existing promotion path.
> 
> ## What happened, in one consumer repo (timeline)
> 
> Two REQs captured as "operator-owned live rollout" work, both cancelled by the user within a week.
> 
> | When (UTC) | REQ | Event |
> | --- | --- | --- |
> | Oct 3 | A | Captured: "install, enable and verify CMS publishing in production", with an Authorization Boundary saying capture does not authorize deployment and the user is the operator. |
> | Oct 6 | A | First run: prepared a runbook report and blocked on "operator authorization and execution of the runbook". No product code changed. |
> | Oct 7 | A, B | User released both through `do-work clarify`. Sitting 1 with the user at the hosts: receipts recorded, one architecture decision and three tooling gaps captured as their own REQs. |
> | Oct 7, late | A, B | Session earmarked both for "an interactive session with the user", then another session flipped them to `blocked` on "the maintainer as operator". |
> | Oct 8, 08:56 | A, B | User released both via clarify. The earmark stayed, so they landed in Pending → Earmarked; re-blocked by hand at 17:52; released again at 21:49. (This part is already reported in the earlier suggestion "a released REQ must not stay silently earmarked".) |
> | Oct 8, 18:52 | B | User cancelled B via `do-work abandon`, reason recorded as "no reason given". |
> | Oct 8, 19:00 to 22:48 | A | Sitting 2: five install phases and three of nine acceptance checks passed, two partial, four unrun. One unrun check needs an unattended scheduled release on a real calendar day at least three days out. |
> | Oct 9 | A | User cancelled A: "I don't need a REQ for it, I'll publish on my own, these REQs are for things that still need to implemented by ai-coders." |
> 
> Cost: REQ A touched 53 commits over six days, three block/clarify/earmark round trips and two multi-hour
> sittings. Every AI-buildable finding from those sittings (a transport redesign, an nginx routing change, two
> recovery tools) became its own REQ and shipped through the normal pipeline. The container REQ added only
> lifecycle ceremony and a claim that could never close.
> 
> ## Where the behaviour lives today
> 
> - `actions/capture.md:107` (Earmark assessment): "when the user instead has to be present as operator ("I
>   need to be at the keyboard for this one"), capture the REQ `blocked` per the External-condition assessment
>   below". This is the only routing for operator presence, and it assumes AI work follows.
> - `actions/capture.md:108` (External-condition assessment): `blocked` is for a request that "can't start
>   until an external condition is met". The examples ("once LM Studio is running", "after the designer
>   replies", "when the staging creds are provisioned") are all preconditions to AI work. There is no example,
>   and no rule, for the case where the condition clearing leaves nothing for the AI to build.
> - `actions/work.md:523`: "Orchestrator wants to park a queued REQ on the operator (production access,
>   credentials, a person at the keyboard) | Flip it to `status: blocked`". Same assumption.
> - `actions/work-reference.md:788` (Environment row, blocked-flip test): flip to `blocked` when "no
>   substantive implementation edits landed this attempt AND the missing thing is a precondition expected to
>   become available on its own (a service comes up, a person answers, credentials get provisioned)". The
>   test has no branch for "the missing thing is the whole remaining work".
> - `actions/work-reference.md:167`: `blocked_at` has "no enforcement threshold; external conditions
>   legitimately take weeks". True for preconditions; for operator-owned work it means the REQ sits in the
>   claimed or blocked column until someone cancels it.
> - `actions/capture-reference.md:34`: "A maintainer can promote any report-only line later by invoking
>   `do-work capture`". This is the existing path the capture rule should route operator work to.
> - `actions/abandon.md`: the only terminal transition for this shape. It is reachable, but nothing in capture
>   or work tells the user that it is the intended end for an operator-only REQ, so the user discovers it
>   after the round trips.
> 
> ## Suggested wording
> 
> **`actions/capture.md:108`, append to the External-condition assessment:**
> 
> > Before writing `blocked`, check what remains once the condition clears. If the remaining work is the
> > operator's own act (deploy, publish, verify on live hosts, approve content, wait for a date) and no
> > AI-buildable work follows, this is not a REQ: emit it as a report-only line pointing at the runbook or
> > checklist, and capture only the tooling gaps the AI must build. `blocked` is for a precondition to AI
> > work, never for work the pipeline cannot finish.
> 
> **`actions/work.md:523`, append to the "park on the operator" row:**
> 
> > If no AI-buildable work remains after the operator acts, do not park it: report that in the hand-back and
> > recommend `do-work abandon REQ-NNN <reason>` with the receipts' home named.
> 
> **`actions/work-reference.md:788`, Environment row, after the blocked-flip test:**
> 
> > Inverse test: if the missing thing is the whole remaining work (the operator's hands, a calendar day),
> > neither fail nor flip; hand back with an abandon recommendation.
> 
> ## Acceptance
> 
> - A capture whose only remaining work is an operator act produces no REQ file and one report-only line
>   naming where receipts go.
> - A run that discovers the same shape mid-flight ends with a hand-back recommending `do-work abandon`
>   with a pre-filled reason, not a `blocked` flip.
> - Existing `blocked` captures with AI work behind the precondition are unchanged.
> 
> --- Maintainer's answers to the four triage questions (AskUserQuestion, 2026-10-09) ---
> 
> F6, orchestrator prose: "compress in place and create structured companion reference where things remain clear without any size limitation"
> F7, board activity correlation: "create an ai-report and explain this to me, my concerns are that I don't want to build scaffold around not-deterministic behaviour"
> F8, disk probe: "Repo root only (Recommended)"
> F11, the 'operator work is not a REQ' proposal: "operator action needs to show up in the "Needs input · Blocked" column, when the code is ready to be released (but not before, before it should have dependencies that might block it, in which case it should be in pending column)"
> 
> --- Maintainer's refinement after plan approval ---
> 
> also when the operator section is done, the REQ should be moved to the DONE column.
> 
> BTW: don't open operator tasks everytime, I'll deploy on my own time, that is nothing special, unless special configuration needs to be done (again this is prose that mostly this orchestrator skill should not concern itself too much)
> 
> --- Triage verdicts (validate-feedback, 2026-10-09, verified against 0.305.84; plan approved by the maintainer) ---
> 
> F1 · 25 releases, early to mid October · Accurate. 0.305.60 (2026-10-02) to 0.305.84 (2026-10-08), one VERSION bump each.
> F2 · Handoff resume wiped claims / deleted plan parts, now fixed · Accurate with corrections: REQ-629 (0.305.64) reset claims to pending and stripped sections; REQ-635 (0.305.70) then REQ-645 (0.305.81) fixed deletion of an indented sample heading, not the real plan. REQ-645 exists because REQ-635 was incomplete.
> F3 · Commit crashed on a backtick typo, now warns and skips · Accurate, overstated: an exit-2 refusal of `do-work commit`, not a crash; skip-and-warn lives in `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go:312-319` (REQ-634, 0.305.69).
> F4 · Disk probe prevents crashes · Push back on "prevents": the probe only displays a VERIFY finding (`queue-kanban/verify.go:318-340`); the run-loop pause was pushed back in the UR-131 triage and never captured.
> F5 · Five releases to clarify one field · Push back: four distinct defects in four places (capture routing REQ-644, board bucketing in Go REQ-643, clarify release path REQ-647, orchestrator rule REQ-648) plus one tooltip (REQ-649); two were wording-only; one REQ per release is `_dev/primes/prime-releases.md` policy. Nothing to refactor.
> F6 · Dense orchestrator prose in huge files · Accept in part: REQ-639 to REQ-642 and REQ-648 added about 1,115 words to `actions/work.md` (12,711 to 13,826) and 646 to `actions/work-reference.md` (22,109 to 22,755); no size budget exists anywhere; ADR-001 records work.md outgrowing a budget twice before; none of the additions is a shell sequence written as prose. Surface-cost N/A (simplification).
> F7 · Activity correlation is hundreds of lines for a card date · Size accurate (357 production + 429 test lines in `activity_correlation.go`); the remedy (file mtime) pushed back on recorded lesson REQ-284: commits change `.git/`, not `do-work/`, and builder commits land in worktrees. One `git log` plus reads shared with the verify probes, ancestry in memory, no merge-base, five git calls per refresh. Maintainer wants an explanation before deciding.
> F8 · Disk probe treats the machine as many volumes · Accept: `verify.go:279-290` probes the repo root plus each builder worktree and dedupes by device id (`:313-316`); the earning incident (REQ-625 §Why: the repo grew to about 20 GB) was in-repo growth; no incident justifies the multi-device part. Surface-cost N/A (deletion).
> F9 · Shortcut commands deleted from CLAUDE.md · Not in this repo: CLAUDE.md last changed 2026-09-03 (de9a1932); the Aliases section of communication-style.md is intact; the installer rewrites only a marked block that links communication-style.md. Consumer-local.
> F10 · Stale REQ-198 marker removed · Not in this window: `do-work/.req-reservations/REQ-000198` was reaped on 2026-08-17 (0ad25193). Consumer-local.
> F11 · Operator work is not a REQ · All six citations check out at 0.305.84 (`capture.md:107-108`, `work.md:523`, `work-reference.md:788`, `:167`, `capture-reference.md:34`, `abandon.md`). Surface-cost: Earned by one consumer REQ costing 53 commits, six days and three round trips. Maintainer rejected "not a REQ, report-only line, recommend abandon" and chose: routine operator work (deploy, publish) is simply not captured; a special operator configuration the request asks to track is a REQ captured `status: blocked` with `depends_on` on the AI-buildable REQs, shown under Pending while dependencies are unmet (`queue-kanban/model.go:1723`) and under Needs input · Blocked once they are met (`:1735`), and completed from clarify when the operator says it is done, which puts it in the Done column.
> 
> Plan decisions: D1 compress the fan-out prose in place and move it to a companion reference with no size limit; D2 produce an ai-report on the activity correlation before deciding its fate; D3 disk probe measures the repo root only; D4 the operator shape above.
> ```
