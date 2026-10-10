# Restart With Parallel Handoff Action

> **Part of the do-work skill.** Finishes the current REQ to a restartable point, encodes the remaining plan as queue state, and writes `do-work/RESTART-PROMPT.md` so a fresh session resumes from one pasted block. It lives in **core** because every input and output is core pipeline state: the queue, `do-work/CHECKPOINT.md`, `do-work/working/` claims, dependency gates, and `pending-answers`.

**State is binding; prose is advisory.** The queue is what the next session actually reads — `do-work run` never opens `RESTART-PROMPT.md`. So the handoff is not where the plan lives. The plan goes into the queue, and the prompt is the paste-ready entry point plus a reference trail for humans.

## When to Use

**Use when:**

- Context is running out, or the session is ending with work still in the queue.
- A fresh session should pick up with different parallelism than this one ran with.
- The user types the `phandoff` shorthand (`crew-members/communication-style.md` → Aliases).
- In a `--coordinate` run, the harness reports context usage above a threshold: the coordinator follows this action on its own. The harness supplies the reading; with no reading there is no automatic trigger. This automatic handoff waits for the next writing gap (no integrator running), runs once per run, writes and commits the handoff only, and does not end the session. In that turn the coordinator's turn-status line comes just before the two announcement lines, which still end the message (Step 6).

**Do NOT use when:**

- You are mid-merge or mid-archive on a REQ. Finish that first — a half-integrated REQ has no honest description.
- You only need a session checkpoint. Run canonical `advance --checkpoint`; this action adds a restart prompt on top of it and does not replace that command.

Stopping before review is fine, and often better: fresh eyes review more carefully than tired ones. Restartable means merged and verified green, not reviewed.

## Input

No arguments. Reads the queue, `do-work/CHECKPOINT.md`, `do-work/working/`, and git state. Writes `do-work/RESTART-PROMPT.md`, plus whatever queue fields Step 1 sets.

## Steps

### Step 1: Encode the plan as state, not prose

Holds, ordering, and collision constraints go into the queue itself: dependency gates, status fields, `pending-answers`. Apply the Schema Read Contract and the Timestamp rule in `actions/work-reference.md` to every field you touch.

The test: **`do-work run` with no other reading must do the right thing for queued work**, because it reads state. Claims left in `do-work/working/` are the one exception, and Step 4 names them in the paste block. If the right behavior depends on someone reading your prose first, the constraint is not encoded yet — go back and encode it.

### Step 2: Survey from git, not memory

```
<skill-root>/scripts/handoff-state-survey.sh
```

Read-only. It prints recent history, every worktree, which `worktree-agent-*` branches are merged and unmerged, and the dirty state of every checkout with `--untracked-files=all`. Pass an integration branch name if the repo does not use `main`.

### Step 3: Classify what the survey found

For each in-flight REQ, state:

- Merged or not merged, with the merge commit and the **full merge range**.
- What remains, in order.
- Uncommitted files **by name**. "0 commits" reads as "nothing there" and is often wrong.
- The canonical `recover` result for each working claim, including its structural writer evidence, and the claim's continue command: bare `advance REQ-NNN`, which changes nothing and names its next phase. Do not classify ownership again from checkpoint prose.

For each worktree, give one verdict:

- **ACTIVE** — my in-flight REQ, detailed above.
- **FOREIGN** — another session's claim; leave byte-identical.
- **REMOVABLE** — all three hold: the branch is in the survey's merged list, its status is clean, and its claim is archived. Write out the `git worktree remove <path>` command; **do not run it.** The next session re-checks all three conditions, then removes.

### Step 4: Write and commit `do-work/RESTART-PROMPT.md`

Exactly two sections, in this order.

**The paste block** — first thing in the file, one fenced code block, nothing above it. It is the complete restart prompt and must work with zero other reading. Write it as instructions addressed to the next session, not as a status document for a human. The resume command is one of:

- to build — `do-work run --fan-out N` (pick N per Step 5)
- to build, when this run was started with `--coordinate` — `do-work run --coordinate --fan-out N`
- to answer questions — `do-work clarify`, included only if some REQ is at `pending-answers`

When this session leaves a REQ claimed in `do-work/working/`, put one line per claim above the resume command: `advance REQ-NNN`, then do the phase it names, so claimed work finishes before new work starts. The resume command's selection reads only the queue, so it does not continue a claim's remaining phases. Never write `recover --take-over` for these claims: it resets a claim (returns it to the queue and strips its orchestrator sections, merged evidence included). In a handoff from a `--coordinate` run, write each claim line as `hand REQ-NNN to an integrator (it enters with advance REQ-NNN)`: the next session coordinates from its first line and never runs a claim's phase itself (`actions/fan-out-reference.md` → **Delegated integration — the coordinator shape**).

Immediately after the commands, write: `These commands are sufficient; everything below them is context.` **If you cannot honestly write that sentence, return to Step 1 until you can.**

**Reference** — after a `---`. Paths, merge ranges, evidence, worktree verdicts, the parallelism analysis from Step 5, and a heads-up list: anything that will bite the next session in its first ten minutes (uncommitted edits from another session in a shared checkout, a half-applied migration, a test that only passes on retry). One line each, naming who should act. This section is for humans and debugging; the next session must not need it to start.

Then commit the file. Before the commit run `git diff --cached --name-only` and stop if anything you did not stage is there — another session shares this index.

### Step 5: Decide parallelism, then mirror it into the queue

You know which REQs collide from having built them. Write it down:

- Which REQs are safe to run concurrently and which must not, each with the reason: overlapping `write_set` paths, dependency gates, shared spec sections.
- The critical path, so the new session starts there rather than on leaves.
- Any REQ to hold back, and what unblocks it.

Set `--fan-out N` in the paste block to match. **Every "must not" is mirrored into queue gates by Step 1** — otherwise the resume command violates your own plan. `write_set` is display-only and gates nothing on its own (`actions/fan-out-reference.md` → **Fan-Out Dispatch**).

### Step 6: Announce it

End the final chat message with exactly two lines:

```
Handoff: do-work/RESTART-PROMPT.md (committed as <sha>)
Resume:  <the resume command(s)>
```

No one should have to ask where the handoff is or what to paste.

## Output Format

One committed `do-work/RESTART-PROMPT.md`, any queue-state edits Step 1 made, and the two announcement lines. No worktree is removed and no foreign claim is touched.

## Rules

- **Never ask the user anything in the main thread.** Per `actions/work.md` Step 3.5, mark the item `- [~] ... D-NN: chose X. Reasoning. Value. Risk.` and proceed. Per Step 8, queue user-facing questions as a follow-up REQ with `status: pending-answers` (the user answers them with `do-work clarify`); a question carrying an outside `Answerer:` routes to that person's stakeholder REQ instead (`do-work stakeholder-answers` ingests their reply).
- **Do any more work after writing the handoff and you rewrite it before handing over.**
- **Refresh `do-work/CHECKPOINT.md` only through `advance --checkpoint`.** Its transaction preserves every live foreign or unlabelled In Progress record.
- Explorers, planners, and reviewers write full output to `do-work/runs/<run>/` and reply with a short summary, so the handoff's context stays clear.

## Common Rationalizations

| If you're thinking...                                                        | STOP. Instead...                                                                             | Because...                                                                                                                                        |
| ---------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| "I finished one more REQ after writing the prompt; the rest still reads fine" | Rewrite the handoff before handing over                                                      | A stale handoff is worse than none. One saying "committed but unmerged" about a merged REQ makes the next session run `git merge`, get "Already up to date", and read the hand-back as empty. |
| "The ordering is explained clearly in the Reference section"                  | Put the gate in the REQ's dependency field                                                   | `do-work run` reads queue state and never opens `RESTART-PROMPT.md`. Prose ordering is invisible to the resume command you just told them to paste. |
| "That worktree is merged and clean, I'll just remove it"                      | Write the `git worktree remove` command into Reference and leave it                          | The third condition is whether its claim is archived. A REQ still in `do-work/working/` with a removed worktree is unrecoverable.                    |
| "`git status` shows 0 commits ahead, so there is nothing there"               | List uncommitted files by name from the survey                                               | Uncommitted work does not move the commit count. That is the reading that loses a builder's whole diff.                                              |

## Red Flags

- The paste block cites a path, a REQ body, or the Reference section. It must stand alone.
- `RESTART-PROMPT.md` names an ordering constraint that no REQ's frontmatter encodes.
- A worktree marked REMOVABLE whose REQ is still in `do-work/working/`.
- The final message does not end with the two announcement lines.

## Verification Checklist

- [ ] `do-work run` with no other reading would do the right thing for queued work — every hold is a queue field
- [ ] Paste block is the first thing in the file, one fence, carrying one `advance REQ-NNN` line per working claim, then the resume command, and the sufficiency sentence
- [ ] Every in-flight REQ states merged/not-merged with its merge range, and uncommitted files by name
- [ ] Every worktree in the survey has exactly one verdict; no worktree was removed
- [ ] Every foreign claim is byte-identical to before this action ran
- [ ] `git diff --cached --name-only` showed only what you staged, and the prompt is committed
- [ ] Final message ends with the `Handoff:` and `Resume:` lines
