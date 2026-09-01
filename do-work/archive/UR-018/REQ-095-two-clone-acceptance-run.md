---
id: REQ-095
title: Two-clone acceptance run — checkpoint poisoning repro and claim-conflict evidence
status: completed
created_at: 2026-08-04T19:44:17Z
claimed_at: 2026-08-04T20:38:51Z
completed_at: 2026-08-04T21:20:00Z
commit: 0526e44
kb_status: promoted
kb_entry: REQ-095-two-clone-acceptance-run-checkpoint-pois.md
user_request: UR-018
domain: testing
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-094]
maintenance: false
write_set: []
related: [REQ-094, REQ-096]
batch: parallel-building
---

# Two-Clone Acceptance Run — Poisoning Repro and Claim-Conflict Evidence

## What

Prove the cross-checkout model with a real two-clone experiment, the way REQ-085's fan-out acceptance run proved worktree dispatch (and found the index-settling bug). Two parts: (1) reproduce the checkpoint-poisoning failure against the **pre-REQ-094** instructions and confirm the writer label stops it; (2) claim the same REQ in two clones, merge, and capture the actual conflict text git produces.

## Detailed Requirements

- Set up two throwaway clones of this repo (scratch space, not inside the repo tree).
- **Poisoning repro:** in clone A claim a dummy REQ (checkpoint entry written, committed); sync to clone B; run clone B's crash-recovery reading per the old rule to show the strip would fire, then per the REQ-094 rule to show the foreign entry is reported and left alone. Record both transcripts.
- **Claim conflict:** both clones claim the same queued dummy REQ (move to `working/`, edit frontmatter, commit); merge one into the other; capture the real conflict git reports (expected: same-path content or rename conflict). Document the observed fix-at-merge resolution.
- Correct any failure-mode claims in the shipped prose **from this evidence, not from reasoning** — if the observed conflict shape differs from what REQ-096's widened dispatch prose predicts, fix the prose.
- Record the run like REQ-085's fan-out run (find where that artifact lives and mirror its form); clean up the throwaway clones afterward.

## Constraints

- The experiment must not touch this repo's own `do-work/` state — dummy REQs live only in the throwaway clones.

## Red-Green Proof

**RED prompt/case:** Under pre-REQ-094 rules, clone B's recovery strips clone A's live claim (deterministic, first sync).
**Why RED now:** No writer identity on checkpoint entries; the model has never been exercised across two clones (fan-out concurrency itself was only ever recorded as Partial, REQ-085).
**GREEN when:** With the writer label, the same sequence leaves clone A's claim intact and reports it; the double-claim merge conflict is captured verbatim and the documented behavior matches the evidence.
**Validation:** User confirmed (approved plan, Phase 1 item 2).

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 1).

---
*Source: approved plan, Phase 1 item 2*

## Triage

**Route: B** - Medium

**Reasoning:** The procedure is fully specified by the REQ's five Detailed Requirements and its Red-Green Proof — it must not be re-derived. What needed discovery was where REQ-085's precedent artifact lives and what form it takes, the exact pre-REQ-094 prose to reproduce the RED against, and the exact shipped rule to test GREEN against. The deliverable is an execution and its record, not a code change.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided execution

*Skipped by work action*

## Exploration

**Precedent artifact (recording form to mirror):** `do-work/archive/UR-016/REQ-085-run-the-live-two-builder-acceptance-test.md`. Its `## Testing` section opens with a *Run parameters* table (one row per operative value: paths, branches, hashes, ranges, dispatch mechanism), then one `###` sub-section per numbered requirement carrying real command output, then an explicit *What this run did not cover* section, then findings as `F-01`-style entries promoted to `## Discovered Tasks` rather than fixed inline. Its `## Scope` records "none as *implementation*" — the REQ's own `## Testing` record is the deliverable. Mirror all of that.

**The RED rule (pre-REQ-094), extracted verbatim from `git show 9c305c0^:actions/work-reference.md` → Crash Recovery (Step 1):**

- `**Named in the checkpoint's `## In Progress (interrupted)` record** → **own crash.** Recover it via substeps 1–3, exactly as before.`
- Only two classes existed: "either this session's own interrupted work, or a claim this session cannot account for." There is no writer identity anywhere in the old section, so *named at all* was sufficient to authorize the destructive path.

**The GREEN rule (shipped, 0.170.0+), `actions/work-reference.md` → Crash Recovery (Step 1):** four classes — own-label (recover), foreign-label (`claim held by <writer>, not touched`, never the takeover ladder), label-less (own only where the checkpoint is locally modified/uncommitted, else report-only), unnamed/no-checkpoint (foreign). Label format per **In-Progress Record (Step 2)**: `- REQ-NNN: [title] — claimed <claimed_at> — writer: <hostname>:<absolute-checkout-path>`, `hostname -s` plus `git rev-parse --show-toplevel`.

**Why the poisoning is reachable at all:** the checkpoint is an ordinary tracked file (`do-work/CHECKPOINT.md`, committed here), so clone A's live claim entry arrives in clone B by a routine `git pull`/`git merge`. Under the RED rule clone B cannot distinguish it from its own crash leftovers.

**Experiment shape:** two throwaway clones in scratch space outside the repo tree, sharing a bare origin so the sync is a real fetch/merge rather than a file copy. Dummy REQs only (`REQ-900`), and the clones get their own throwaway `do-work/` fixture so this repo's queue is never involved.

## Scope

**Files I will touch:** none as *implementation*. This REQ's deliverable is its own `## Testing` record plus the two throwaway clones and their bare origin, all in scratch space outside the repo and deleted afterwards. If the observed evidence contradicts shipped prose, the prose fix is added to this list before it is written.

**Acceptance criteria (restated from REQ):**
- [ ] Two throwaway clones set up in scratch space, not inside the repo tree (req 1)
- [ ] Poisoning repro: clone A claims a dummy REQ (checkpoint entry written and committed), synced to clone B; clone B's recovery read recorded under the OLD rule (strip would fire) and under the REQ-094 rule (foreign entry reported, left alone) — both transcripts recorded (req 2)
- [ ] Claim conflict: both clones claim the same queued dummy REQ, merge one into the other, real git conflict text captured verbatim; the fix-at-merge resolution documented (req 3)
- [ ] Any shipped failure-mode claim corrected from the observed evidence, not from reasoning (req 4)
- [ ] Recorded in REQ-085's form; throwaway clones cleaned up (req 5)
- [ ] This repo's own `do-work/` state untouched by the experiment (constraint)

## Implementation Summary

**Files changed:** none — this REQ is an execution. Its deliverable is the `## Testing` record below,
plus one fold-in appended to REQ-096's `## Addendum` (evidence handed to the REQ that will write the
prose) and one new queue REQ for the defect found. The six throwaway fixtures under
`<scratch>/req095/run1..run6/` were created, consumed, and deleted.

**What was done:** Ran the two-clone acceptance test for the first time. Six fixtures, each a bare
origin plus two clones with their own throwaway `do-work/` tree, exercised: the pre-REQ-094 poisoning
(reproduced deterministically), the shipped writer-label rule against byte-identical input, the
double-claim merge, the byte-identical double-claim merge, two disjoint concurrent claims, and the
label-less entry under a merge resolution. Two defects were found and filed rather than fixed.

## Testing

### Run parameters (requirement 1)

| Item | Value |
| --- | --- |
| Scratch root (`<scratch>` below) | `/tmp/claude-0/-home-user-skill-do-work/d296c685-7470-5207-9ec5-d1b06940ccd5/scratchpad` — outside the repo tree, per the REQ's constraint |
| Fixture shape | `<scratch>/req095/runN/{origin.git,clone-a,clone-b}` — a **bare origin** plus two real clones, so every sync is a genuine `git fetch`/`git merge`, not a file copy |
| Dummy REQs | `REQ-900` (poisoning target), `REQ-901` (double-claim target) — throwaway, `user_request: UR-900`, never in this repo's queue |
| Clone A writer label | `vm:<scratch>/req095/runN/clone-a` |
| Clone B writer label | `vm:<scratch>/req095/runN/clone-b` |
| RED rule source | `git show 9c305c0^:actions/work-reference.md` → Crash Recovery (Step 1) |
| GREEN rule source | `actions/work-reference.md` → Crash Recovery (Step 1), as shipped at v0.170.1 |
| Runs | run1 poisoning with a live local edit · run2 deterministic poisoning replay · run3 shipped rule + double claim · run4 byte-identical double claim · run5 two disjoint claims · run6 label-less entry under merge resolution |
| Scripts | `setup.sh`, `claim-in-a.sh`, `recover-old-rule.sh`, `consequence-in-a.sh`, `run2-clean-strip.sh`, `run3-new-rule.sh`, `run3-double-claim.sh`, `run4-identical-claim.sh`, `run5-different-reqs.sh`, `run6-labelless-backdoor.sh` — transcripts captured per run |
| This repo's `do-work/` | untouched by the experiment; verified below |

The two halves of the writer label are worth reading off that table: both clones report hostname `vm`,
so on a single machine the **path half is the sole discriminator**. That is the shipped rationale
("Neither half identifies a checkout alone", **In-Progress Record (Step 2)**) confirmed rather than
contradicted — the hostname half only starts working across machines.

### Poisoning repro — RED under the old rule (requirement 2)

*(Quoted `## In Progress (interrupted)` lines inside the transcripts below are indented two
spaces so they cannot read as this file's own headings to the line-based section parsers —
`tools/checks/qualify.sh`'s awk and crash recovery's `^## ` strip. Nothing else is altered.)*

Clone A claimed `REQ-900` exactly as Step 2 prescribes (move to `working/`, `status: claimed`,
`claimed_at`, one labelled checkpoint entry), appended a `## Triage` and a `## Scope` that exist
nowhere but that file, committed its bookkeeping and pushed. Clone B synced and classified:

```
=== B3. classification under the OLD rule ===
REQ-900 IS named in the '## In Progress (interrupted)' record.
OLD rule verdict: own crash -> recover via substeps 1-3 (DESTRUCTIVE).
```

The old rule had only two classes ("either this session's own interrupted work, or a claim this
session cannot account for") and no writer identity anywhere, so *named at all* authorized substeps
1–3. Clone B executed them, committed and pushed. **Clone A's next routine sync (run 2), with clone A
sitting between two pipeline steps — its claim committed, nothing modified since:**

```
=== R2-3. clone A performs a routine sync ===
Updating 86ba65a..94f62b0
Fast-forward
 do-work/CHECKPOINT.md                             |  1 -
 do-work/queue/REQ-900-dummy-poisoning-target.md   | 16 ++++++++++++
 do-work/working/REQ-900-dummy-poisoning-target.md | 31 -----------------------
 3 files changed, 16 insertions(+), 32 deletions(-)
merge exit status: 0

=== R2-4. clone A after the sync ===
--- do-work/working/ ---
(empty)
--- clone A's own CHECKPOINT.md ---
# Session Checkpoint

  ## In Progress (interrupted)

```

**F-01 — the poisoning is real, deterministic, and silent.** Exit 0, a fast-forward, no conflict and
no warning. Clone A's live claim left `working/`, its `## Triage` and `## Scope` were destroyed, the
REQ went back to `pending` in the queue, and clone A's own claim record was erased from its own
checkpoint. This confirms `actions/work-reference.md`'s stated mechanism verbatim — "another
checkout's live claim then arrives here by a routine `git pull`, looking exactly like a local one,
and without the label recovery reads it as this session's own crash and strips a REQ someone is
actively building."

**F-02 — the window is narrower than "any sync", and it is exactly the routine one.** Run 1 ran the
same sequence with clone A holding *local* edits to the claimed REQ, and git protected it both ways:

```
=== A2. sub-case (a): A has UNCOMMITTED mid-run edits to the claimed REQ ===
error: Your local changes to the following files would be overwritten by merge:
	do-work/working/REQ-900-dummy-poisoning-target.md
Please commit your changes or stash them before you merge.
Aborting

=== A3. sub-case (b): A commits its bookkeeping first, then syncs ===
CONFLICT (modify/delete): do-work/working/REQ-900-dummy-poisoning-target.md deleted in origin/main
and modified in HEAD.  Version HEAD of do-work/working/REQ-900-dummy-poisoning-target.md left in tree.
merge exit status: 1
```

So the strip is silent only while the claiming checkout is *idle between steps* — which is the state
it is in every time it pauses, and the state a crashed run is in permanently. Sub-case (b) also
leaves the REQ in **both** `working/` and `queue/` at once, a duplicate-id state no probe was looking
for at the time of the run.

### Poisoning repro — GREEN under the shipped rule (requirement 2)

Run 3 replayed the identical input state (`sha256 57e71ae2a99f…` for the claimed REQ, byte-for-byte
what run 2 received) against the shipped rule:

```
=== R3-2. clone B derives its OWN writer label the same way Step 2 does ===
clone B's own label: vm:<scratch>/req095/run3/clone-b
entry's label:      vm:<scratch>/req095/run3/clone-a

=== R3-3. classification under the SHIPPED rule ===
foreign-label -> report and DO NOT TOUCH
report line: claim held by vm:<scratch>/req095/run3/clone-a, not touched
takeover ladder: not entered (a foreign label is positive attribution; age adds nothing)

=== R3-5. byte-identity check on the foreign claim ===
sha256 after clone B's recovery pass: 57e71ae2a99f81664c9c860bd18dcfd345e61726988374e15fd486f739a27993
IDENTICAL — the foreign claim was not modified
(empty above = clone B wrote nothing)
```

**F-03 — GREEN.** Byte-identical, nothing staged, the `claim held by <writer>, not touched` report
emitted verbatim, the three-hour ladder not entered, and clone B moved on to `REQ-901` as the next
pending REQ. Clone A's subsequent sync was `Already up to date.` with its claim, its `claimed_at`, all
three `##` sections and its checkpoint entry intact.

### Claim conflict — the real merge text (requirement 3)

Both clones claimed `REQ-901` independently from a common base, each writing its own `claimed_at` and
its own `## Triage`. Clone B pushed; clone A merged:

```
=== P2-3. clone A syncs — the fix-at-merge moment ===
Auto-merging do-work/CHECKPOINT.md
CONFLICT (content): Merge conflict in do-work/CHECKPOINT.md
Auto-merging do-work/working/REQ-901-dummy-double-claim-target.md
CONFLICT (content): Merge conflict in do-work/working/REQ-901-dummy-double-claim-target.md
Automatic merge failed; fix conflicts and then commit the result.
merge exit status: 1

=== P2-4. git status --porcelain ===
UU do-work/CHECKPOINT.md
UU do-work/working/REQ-901-dummy-double-claim-target.md
```

```
=== P2-6. conflicted do-work/CHECKPOINT.md, verbatim ===
  ## In Progress (interrupted)

- REQ-900: … — writer: vm:<scratch>/req095/run3/clone-a
<<<<<<< HEAD
- REQ-901: Dummy double-claim target — claimed 2026-08-04T20:55:40Z — writer: vm:…/run3/clone-a
=======
- REQ-901: Dummy double-claim target — claimed 2026-08-04T20:55:10Z — writer: vm:…/run3/clone-b
>>>>>>> origin/main
```

```
=== P2-7. conflicted do-work/working/REQ-901-*.md, verbatim ===
created_at: 2026-08-04T20:00:00Z
<<<<<<< HEAD
claimed_at: 2026-08-04T20:55:40Z
=======
claimed_at: 2026-08-04T20:55:10Z
>>>>>>> origin/main
user_request: UR-900
…
<<<<<<< HEAD
**Reasoning:** Clone A's triage of REQ-901.
=======
**Reasoning:** Clone B's triage of REQ-901.
>>>>>>> origin/main
```

**F-04 — the shape is a plain content conflict, never a rename conflict.** This REQ predicted "same-path
content or rename conflict"; the rename never surfaces, because both sides perform the *identical*
rename (`queue/REQ-901…` → `working/REQ-901…`) and git resolves it silently, leaving only the content
inside to conflict. Two files conflict, and the REQ file's conflict is entirely made of the two claim
writes — `claimed_at` and the generated sections. **Fix-at-merge resolution:** keep one claim on the REQ
file (a human/releaser decision — whichever checkout is actually building it, evidenced by its
worktree), and keep **both** lines in `CHECKPOINT.md`, which is what the shipped prose already requires
("one id under two labels is the honest record of two checkouts", **In-Progress Record (Step 2)**).

**F-05 — with byte-identical claim writes the REQ file does not conflict at all, and the writer label
is the only thing that catches the double claim.** Run 4 gave both clones the same `claimed_at` and no
generated sections yet — the realistic case when two checkouts claim within the same second:

```
=== R4-1. both clones make a byte-identical claim write ===
clone B's REQ-901 sha256: 35ec100dbffe0453ac82e9b478a9e31a71e943fafec9d2356e1d197c2d273c36
clone A's REQ-901 sha256: 35ec100dbffe0453ac82e9b478a9e31a71e943fafec9d2356e1d197c2d273c36

=== R4-2. clone B pushes, clone A merges ===
Auto-merging do-work/CHECKPOINT.md
CONFLICT (add/add): Merge conflict in do-work/CHECKPOINT.md

=== R4-3. git status --porcelain ===
AA do-work/CHECKPOINT.md

=== R4-5. and the REQ file itself (did it conflict?) ===
unmerged: do-work/CHECKPOINT.md
```

The REQ file merged clean. The label — which differs between checkouts *by construction* — is what
turns an invisible double claim into an `AA` conflict a human must resolve. That is a second,
unstated job the writer label is doing, and it is the one that makes "duplicates are fixed at merge"
true for this case rather than merely hoped.

**F-06 — `do-work/CHECKPOINT.md` conflicts on *every* concurrent claim, including fully disjoint
ones.** Run 5 had clone A claim `REQ-900` and clone B claim `REQ-901` — no overlap of any kind, the
normal sanctioned case under claim-anywhere:

```
=== R5-2. clone B pushes, clone A merges ===
CONFLICT (add/add): Merge conflict in do-work/CHECKPOINT.md
merge exit status: 1

=== R5-3. git status --porcelain ===
AA do-work/CHECKPOINT.md
R  do-work/queue/REQ-901-dummy-double-claim-target.md -> do-work/working/REQ-901-dummy-double-claim-target.md
```

Both REQ files merged cleanly (the rename shows as a plain `R`); the checkpoint did not, because two
single-line appends land at the same position. So under claim-anywhere with a committed `do-work/`,
this file is a guaranteed conflict point per concurrent claim. The conflict is trivial to resolve —
**keep every entry from both sides** — but *nothing shipped says so*, and both one-sided resolutions
lose data: taking ours drops another checkout's live claim record (the strip, by hand), taking theirs
drops our own and makes our own crash unrecoverable. Handed to REQ-096 as a third fold-in, since its
widened-dispatch prose is where "claim conflicts between checkouts are ordinary git conflicts fixed at
merge" gets written.

### Defect found: the label-less bullet is unsound under claim-anywhere (requirement 4)

The shipped label-less rule reads: *"own only where `do-work/CHECKPOINT.md` is locally modified or
otherwise uncommitted in this checkout, which is evidence this checkout wrote it and has not shared
it; recover it as an own crash."* F-06 breaks its premise. Run 6 gave clone A an old-version,
**label-less** entry and had clone B claim concurrently, so B had to resolve the guaranteed checkpoint
conflict:

```
=== R6-3. clone B syncs — guaranteed AA conflict, resolved by keeping BOTH entries ===
AA do-work/CHECKPOINT.md
--- conflict resolved (both entries kept), merge NOT yet committed ---
M  do-work/CHECKPOINT.md

=== R6-4. clone B now runs Step 1 ===
working/ contents:
REQ-900-dummy-poisoning-target.md
REQ-901-dummy-double-claim-target.md
M  do-work/CHECKPOINT.md
 do-work/CHECKPOINT.md | 1 +
REQ-900's entry is label-less. The checkpoint IS uncommitted here. Verdict by the rule:
  -> OWN CRASH -> recover via substeps 1-3 -> clone A's live claim is STRIPPED.
```

**F-07 (critical) — "locally modified" is not evidence of authorship once merges are routine.** A
checkout that resolved a checkpoint conflict — which claim-anywhere makes it do on *every* concurrent
claim — holds a modified checkpoint for a reason that has nothing to do with who wrote which entry, so
the heuristic classifies a foreign label-less entry as an own crash and strips a live claim. That is
the 2026-07-01 incident, reachable again through the label-less door. Filed as **REQ-104**, not fixed
here: the fix is a behaviour change to a suite-pinned rule with a real trade-off (narrow the heuristic
to "modified *and* no merge in progress", or drop it so label-less is always report-only and
pre-0.170.0 crashes lose auto-recovery), and this REQ's Scope declares no implementation.

### What this run did not cover

- **Installs where `do-work/` is untracked — the common case — are entirely outside every result
  above.** Nothing syncs, so no claim, no checkpoint entry and no conflict ever travels between
  checkouts: the poisoning cannot happen, and neither can any fix-at-merge detection. `queue-kanban
  verify`'s `duplicate-req-id` probe and REQ-098's new probes are the only detectors there. Every
  conflict shape recorded here is a **committed-`do-work/` result**.
- **One machine, so the hostname half of the writer label was never exercised as a discriminator** —
  only the path half was (see Run parameters).
- **No real builder ran.** Claims were made by the prescribed Step 2 file operations; no
  implementation, worktree, or hand-back was involved. Wall-clock concurrency of builders is REQ-100's
  deliverable, not this one's.
- **Push races were not tested** — clone B always pushed first and clone A always merged. A rejected
  non-fast-forward push is a different (and much louder) shape.
- **Three or more checkouts** were not tested; nothing here says whether the checkpoint conflict stays
  a two-way resolve at higher counts.

### Fixture cleanup (requirement 5)

```
$ git status --porcelain --untracked-files=all
 M do-work/CHECKPOINT.md
 M do-work/queue/REQ-096-execution-model-regrain.md
RM do-work/queue/REQ-095-two-clone-acceptance-run.md -> do-work/working/REQ-095-two-clone-acceptance-run.md
?? do-work/queue/REQ-104-labelless-entry-authorship-heuristic.md
?? do-work/working/baseline-failures.txt
?? do-work/working/baseline.json
$ git worktree list
/home/user/skill-do-work  e97826e [claude/ur-018-do-work-batch-cbmfeb]
$ git branch --list 'worktree-agent-*'
$ find . -name 'REQ-90*.md' -not -path './.git/*'
(none)
$ find /home/user /tmp/claude-0 -maxdepth 6 -type d -name 'do-work' -not -path '/home/user/skill-do-work/*'
(none)
```

All six fixtures and their bare origins removed. The only dirt in this repo is this REQ's own
bookkeeping — the claim move, the checkpoint entry, REQ-096's fold-in, REQ-104, and Step 5.75's
baseline files. No `REQ-900`/`REQ-901` file ever existed inside the repo, no stray `do-work/`
directory outside the project root, no leftover worktree or branch. The constraint holds.

### Regression check (Step 6.5)

No shipped file changed, so the gate is "the suite is exactly where Step 5.75 left it":

```
$ bash _dev/tests/contract-regressions.sh 2>&1 | grep -c '^FAIL'
8
```

All eight FAIL lines match `do-work/working/baseline-failures.txt` name-for-name — the five
`mid-update failure:` probes, the two `dirty install:` probes, and their roll-up. **Zero new
regressions.**

Those eight are a **pre-existing environment failure, not a code regression**, and every REQ in this
batch inherits them: `_dev/tests/update-script-behavior.sh` injects its failure with
`chmod 500 "$failing_install/docs"`, and this container runs as **root**, which ignores the permission
bit — so the updater completes where the probe expects it to fail. `baseline.json` records
`{"exit_status": 1, "launched": true}`, so the comparison is against a real test run, not an
un-launched one.

## Qualification

**Mechanical check: FAIL, accepted — `tools/checks/qualify.sh` cannot express a no-code-change REQ.**

```
$ tools/checks/qualify.sh do-work/working/REQ-095-two-clone-acceptance-run.md
FAIL: no '## Implementation Summary' file list — run after Step 6.25
```

Identical to REQ-085's accepted FAIL, for the identical reason: the script requires at least one
`` - `path` `` bullet under `## Implementation Summary`, and separately rejects a list of only
`do-work/` paths (its design-artifact exception covers artifacts placed *outside* `do-work/`, the
opposite shape). This REQ changed no project file. Its deliverable is the `## Testing` record above;
the only files it wrote are `do-work/` paths — REQ-096's fold-in and REQ-104 — and both are pipeline
bookkeeping the orchestrator owns, not implementation.

**Judgment checks, run by hand in the script's place:**

- **Requirements traced (5/5).** req 1 → *Run parameters*; req 2 → *Poisoning repro* RED and GREEN
  sub-sections; req 3 → *Claim conflict* with verbatim `UU`/`AA` output; req 4 → *Defect found* plus
  the F-04 correction of this REQ's own predicted conflict shape; req 5 → *Fixture cleanup* with
  command output. Constraint → same section.
- **Substantive, not a placeholder.** Every quoted block is real command output captured to a
  transcript at run time, including the two exit statuses and three sha256 digests that carry the
  byte-identity and byte-difference claims.
- **Evidence, not reasoning.** Each finding names the run that produced it. F-02, F-04, F-05 and F-06
  all *contradict* what the plan or this REQ predicted, which is the check that the record was written
  from output rather than from the plan.

## Discovered Tasks

- [critical] The label-less crash-recovery bullet classifies a foreign entry as an own crash whenever
  this checkout's `do-work/CHECKPOINT.md` is uncommitted — which claim-anywhere makes routine, because
  every concurrent claim forces a checkpoint merge resolution (F-06/F-07). Reintroduces the
  2026-07-01 strip through the label-less door. → queued as **REQ-104**.
- [normal] `actions/work-reference.md`'s widened-dispatch prose needs the checkpoint-conflict
  resolution rule stated (keep every entry from both sides; never resolve one-sided) — F-06. → folded
  into **REQ-096**'s `## Addendum` rather than queued separately, because REQ-096 rewrites exactly
  that prose.

## Lessons Learned

**What worked:**
- A **bare origin plus two real clones** rather than a copied directory. Every sync in the record is a
  genuine `git fetch`/`git merge`, which is why the exact conflict strings (`CONFLICT (add/add)`, `AA`,
  `CONFLICT (modify/delete)`, `Fast-forward … merge exit status: 0`) are quotable at all. A file-copy
  fixture would have produced a plausible narrative and no evidence.
- Replaying the **identical input state** (same sha256) against the old rule and the shipped rule, in
  two separate fixtures rather than by resetting one. It makes the A/B airtight and the transcripts
  independently re-runnable.
- Pushing past the two requirements into the adversarial variants (byte-identical claims, disjoint
  claims, a label-less entry after a merge resolution). Both defects came from those three runs; the two
  required runs alone would have produced a clean bill.

**What didn't:**
- The REQ's own predicted conflict shape — "same-path content or **rename** conflict" — is wrong, and
  reasoning would have kept it. Git resolves the `queue/` → `working/` rename silently *because both
  sides perform the identical rename*, so only the content inside ever conflicts. Prose predicting a
  rename conflict would have shipped a claim no run can reproduce.
- The first framing of the poisoning as "any routine sync strips the claim" is also wrong. Git protects
  a claiming checkout that has local edits (refuses the merge) or committed divergent ones (raises
  `modify/delete`). The silent strip needs the claimant to be *idle between steps* — which is both
  narrower than assumed and the state every paused or crashed run is actually in.
- `tools/checks/qualify.sh` and `tools/checks/scope-drift.sh` both refuse a no-code-change REQ
  (`FAIL: no '## Implementation Summary' file list`, `SKIP: … exit 2`). REQ-085 hit the identical wall
  a batch earlier. Accepted both times with the reasoning written down; still a script gap, not a REQ
  defect.

**Worth knowing:**
- **Every conflict result here is a committed-`do-work/` result.** On the common install, where
  `do-work/` is untracked, nothing syncs between checkouts: the poisoning cannot happen and neither can
  any fix-at-merge detection. Do not generalize these transcripts to that install.
- **`do-work/CHECKPOINT.md` is a guaranteed conflict point per concurrent claim, disjoint claims
  included** — two single-line appends land at the same position. That is load-bearing twice over: it is
  the only detector of a byte-identical double claim (F-05), and it is what dirties the checkpoint and
  breaks the label-less authorship heuristic (F-07).
- On one machine every clone reports the same `hostname -s`, so the **path half of the writer label is
  the sole discriminator** locally. Any future test of the hostname half needs two machines.
- The eight suite FAILs in this container are `chmod 500` injections defeated by running as **root** —
  environment, not code. Compare against `do-work/working/baseline-failures.txt` name-for-name rather
  than expecting green.

## Orientation

Crash recovery's writer label now has an acceptance record behind it: the pre-0.170.0 poisoning is
reproduced as a silent fast-forward that erases a live claim, and the shipped foreign-label rule is
shown leaving that claim byte-identical. The run also establishes what fix-at-merge actually looks like
across checkouts — content conflicts on the REQ file, an `add/add` conflict on `CHECKPOINT.md` for every
concurrent claim — and found that the label-less recovery bullet is unsound once merges are routine
(filed as REQ-104). This is evidence about the `do-work/CHECKPOINT.md` + Crash Recovery subsystem
described in `actions/work-reference.md`; no code or contract changed. `prime_files` is empty, so no
prime staleness spot-check applied.

## Review

**Overall: 93%** | 2026-08-04T21:20:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | n/a (no code changed) |
| Test Adequacy | 95% |
| Scope | 90% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 2 minor
**Acceptance:** Pass — all five requirements and the constraint are evidenced by real command output; the two required runs plus three adversarial variants produced two filed defects and four corrections to predicted behavior.
**Suggested testing:** 3 items
**Follow-ups created:** REQ-104

**Requirements checklist:**
- [x] req 1 — two throwaway clones in scratch space outside the repo tree (six fixtures, bare origin each)
- [x] req 2 — poisoning repro under the old rule and the shipped rule, both transcripts recorded (runs 1–3)
- [x] req 3 — double-claim merge run, real conflict text captured verbatim, fix-at-merge resolution documented (run 3, P2-3…P2-8)
- [x] req 4 — one shipped claim confirmed verbatim, one contradicted and filed as REQ-104, this REQ's own predicted conflict shape corrected from output (F-04)
- [x] req 5 — recorded in REQ-085's form; all six fixtures removed with cleanup output
- [x] constraint — this repo's `do-work/` untouched beyond this REQ's own bookkeeping

**Minor:**
- `## Scope` said "the prose fix is added to this list before it is written." No prose was edited: the
  one contradicted rule became REQ-104 (a behaviour change with a real trade-off, and this batch's
  REQ-096 owns the surrounding lines) and the conflict-resolution rule became REQ-096's third fold-in.
  Defensible, and it matches REQ-085's file-a-defect-don't-fix-it precedent, but it is a declared-and-not-
  taken path worth naming rather than leaving to inference. Scored as the 10% off Scope.
- Two transcript lines quoted `## In Progress (interrupted)` at column 0 inside code fences, where the
  line-based section parsers (`qualify.sh`'s awk, crash recovery's `^## ` strip) would read them as this
  file's own headings. Indented two spaces with a note; worth watching for in any future REQ that quotes
  checkpoint content.

**Scope drift:** `tools/checks/scope-drift.sh` exits 2 (`SKIP: no '## Implementation Summary' file
list`) — the same no-code-change gap `qualify.sh` hits. Compared by hand: `## Scope` declared no
implementation files and none were written. The two `do-work/` writes (REQ-096's fold-in, REQ-104) are
orchestrator bookkeeping, which Scope's file list does not cover in either direction. No drift.

**Restatement sweep:** not applicable — this REQ redefined no contract token, schema field, gate wording
or prescribed-command output shape. Its two writes are a queue REQ and a fold-in paragraph, neither of
which any other text restates. The sweep's *findings* were handed forward instead: F-06's
conflict-resolution rule and F-04's corrected conflict shape go to REQ-096, which does rewrite the prose
they belong to.

**Suggested additional testing:**
- Two machines, to exercise the hostname half of the writer label as a discriminator.
- Three or more checkouts, to see whether the checkpoint conflict stays a two-way resolve.
- A rejected non-fast-forward push (this run always had clone B push first and clone A merge).

*Reviewed by review-work action (pipeline mode, in-session)*
