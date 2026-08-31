# Run Simple REQs Action

> **Part of the do-work skill.** Lists the pending REQs whose implementation a cheaper model can be trusted with, estimates the batch, and hands them to `do-work run` after one confirmation. It belongs in core because it reads core's queue and dispatches core's pipeline — it completes `actions/work.md` rather than extending any sibling package.

**You supply the model; this action supplies the selection.** Nothing here names, chooses, or switches a model, and nothing detects which harness is running. The user has already launched this session in whatever cheaper environment they wanted; the only question left is *which queued work is safe to do here*, and that is the whole job.

## When to Use

**Use when:**

- This session runs on a cheaper or faster model than usual and you want the queue's mechanical work done in it.
- You want to see which REQs qualify, and what the batch costs in active minutes, before committing to anything.

**Do NOT use when:**

- You want the whole queue — that is `do-work run` (`actions/work.md`).
- You want to skip work nobody would notice — that is `do-work run --skip-impact-negligible`, which filters the **impact** axis. This action filters **effort**, and the two are independent (`actions/work-reference.md` → Request File Schema).
- You only want the list and the estimate — run the Step 1 script and stop; it writes nothing.

## Input

Two flags are accepted, and they are **not** interchangeable in how they travel:

- **`--fan-out [N]`** — forwarded verbatim to `do-work run` in Step 4. It changes how many of the set run at once, never which, so the handoff is the right place for it.
- **`--skip-impact-negligible`** — passed to the **Step 1 selector**, never forwarded to Step 4. Forwarding it would be silently inert: Step 4 names REQs explicitly, and an explicitly-named `REQ-NNN` overrides that flag by design (`actions/work.md` → Input), so the negligible work the user asked to omit would be built anyway. Any filter that must survive the handoff has to run before it.

**`--wave N` is rejected.** It is mutually exclusive with every targeting token (`actions/work.md` → Input), and Step 4 always passes targeting tokens, so the combination could only ever fail at the handoff. Stop and report: `run-simple-reqs selects its own set, so --wave cannot apply — use do-work run --wave N for a depth-scoped run.`

**Targeting tokens (`REQ-NNN` / `UR-NNN`) are rejected here.** This action computes its own set, so naming a REQ would either duplicate that set or contradict it. Stop and report: `run-simple-reqs computes its own set — use do-work run REQ-NNN to target a specific REQ.`

Any other unrecognized argument is rejected rather than ignored, exactly as `actions/work.md` rejects residue.

## Steps

### Step 1: Select and estimate

Run the canonical typed selector from the project root:

```bash
<skill-root>/tools/do-work-cli.sh --repo-root <project-root> next --simple [--skip-impact-negligible]
```

Pass `--skip-impact-negligible` through when the user gave it — that is the only place it can take effect.

It prints one typed `selected` row per qualifying REQ with that REQ's p50 estimate, every excluded REQ with a stable code and actionable reason, exact next and verification commands, and a final `run_set: <ids>` line. The same records are available with global `--format json`; that last line is the text contract Step 4 runs, so read it rather than re-deriving ids from the table.

The command owns the predicate and consumes the same typed snapshot and dependency graph as ordinary work selection. Do not restate, re-scan, or re-judge it here. If the canonical command is missing or fails, stop with its actionable output; never fall back to a hand-built selection. `tools/select-simple-reqs.sh` remains only a compatibility entry point to this command.

### Step 2: Show the list

Present the script's output as it stands. Keep the held-back section: a selector that silently narrows the queue is indistinguishable from an empty queue, and the reasons are what let the user correct a mis-tagged REQ instead of wondering where it went.

When `run_set:` is empty, say so plainly, name what was held back and why, and stop. An empty set is a normal answer, not a failure.

### Step 3: Confirm once

Ask whether to run the listed REQs, quoting the count and the total estimate. Load `crew-members/clear-questions.md` before asking. One gate, then commit to the answer — this is the only confirmation in the flow.

### Step 4: Hand off to the work pipeline

Run `do-work run` with the `run_set` ids and any forwarded flags:

```
do-work run REQ-NNN REQ-NNN [--fan-out N]
```

`--fan-out` is the only flag forwarded here. `--skip-impact-negligible` was already applied by the selector and must not be repeated, and `--wave` was rejected in Input.

Every later step — claim, triage, build, qualify, test, review, archive — is `actions/work.md` unchanged. Do not reimplement any of it here, and do not pass a REQ the selector did not list.

## Output Format

The selector's report, then either the confirmation question or the empty-set statement, then `actions/work.md`'s own per-REQ progress output. This action adds no report of its own and writes no files.

## Rules

- **Never widen the set by hand.** If a REQ looks like it belongs and the selector disagrees, the fix is that REQ's `effort_estimate` or its markers, not an extra id on the Step 4 command line.
- **The estimate is informational.** It is the same P50 forecast `actions/estimate-reference.md` defines — never a deadline, and never a reason to stop a REQ that runs long.

## Common Rationalizations

| If you're thinking...                                                              | STOP. Instead...                                                                                 | Because...                                                                                                                                                            |
| ---------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| "I'll add the REQ the selector skipped for `depends_on` — it looks independent"     | Run its dependency first, or leave it for the next pass                                          | Step 4 names REQs explicitly, and an explicitly-named REQ **bypasses `depends_on`** by design (`actions/work.md` → Input). The selector's readiness check is the only thing standing between this action and a silent dependency-gate bypass |
| "This REQ is `maintenance: true` but the diff is one line, so it's mechanical"      | Leave it for a full-strength session                                                             | `crew-members/maintenance.md` work edits the skill's own operating instructions, where nothing tests a rule. Instruction REQs are always small in diff, so effort alone selects the queue's highest-judgment work first |
| "`tdd: true` means extra ceremony — hold it back too"                               | Select it; test-first work is a **better** fit, not a worse one                                   | A `tdd: true` REQ carries an objective pass/fail gate, often a captured RED case in `## Red-Green Proof` — a stronger check than the qualification-plus-review a non-TDD REQ gets                                    |
| "The user passed `--skip-impact-negligible`, so I'll forward it to `do-work run`" | Pass it to the Step 1 selector instead                                                           | Step 4 names REQs explicitly, and explicit naming overrides that flag (`actions/work.md` → Input). Forwarded, it is inert, and the queue builds the negligible REQs the user asked to skip |
| "Nothing qualified, so I'll relax the filter and pick the closest few"              | Report the empty set and stop                                                                    | The held-back list already names what to fix; loosening the predicate at read time makes the queue's own `effort_estimate` and marker fields meaningless                                                             |

## Verification Checklist

- [ ] The ids passed to `do-work run` are exactly the selector's `run_set` line — no additions, no removals
- [ ] The held-back REQs and their reasons were shown to the user, not summarized away
- [ ] Exactly one confirmation was asked before any REQ was claimed
