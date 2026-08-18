---
id: REQ-245
title: Name fabricated stamps in the board's future-stamp warnings
status: claimed
created_at: 2026-08-18T12:28:33Z
user_request: UR-055
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-244]
batch: timestamp-stamping-integrity
effort_estimate: trivial
write_set: ["skills/do-work-board/tools/queue-kanban/model.go", "skills/do-work-board/tools/queue-kanban/verify.go", "skills/do-work-board/tools/queue-kanban/*_test.go"]
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-18T12:43:06Z
  basis:
    - trivial short-circuit
claimed_at: 2026-08-18T12:43:06Z
route: A
---

# Name Fabricated Stamps in the Board's Future-Stamp Warnings

## What

The board's future-stamp diagnosis messages name exactly one cause — "likely local wall-clock time stamped with a Z suffix" — but a fully fabricated value is a second, now-observed cause, and the current wording sends that reader to the wrong fix. Reword the diagnosis clauses to name both causes; keep the fix instruction (rewrite with the current UTC instant per the Timestamp rule) unchanged, since it is correct for both.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

Sibling messages found at capture — update together so they don't drift:

- `skills/do-work-board/tools/queue-kanban/model.go:379` — generate-time data warning: "…likely local wall-clock time stamped with a Z suffix; fix: rewrite with the current UTC instant…"
- `skills/do-work-board/tools/queue-kanban/model.go:1232` — reversed-span message: "…one stamp is usually local wall-clock time written with a Z suffix…"
- `skills/do-work-board/tools/queue-kanban/verify.go:371` — verify-time future `claimed_at`: "…usually local wall-clock time written with a Z suffix"

Comments asserting the single-cause story (e.g. `timestamp_test.go:42`, `completion_anomaly_test.go:227`, `model.go:1338`) should be brought in line where they would otherwise contradict the new wording.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this change — versioning, parser lock-step, build outputs. Read it before touching the tool.
- Message-text change only: no new checks, no threshold changes, the 2-minute skew allowance stays as is.
- Finding provenance (validate-feedback triage, this session): verdict Accept; Surface-cost N/A — text accuracy fix to an existing warning, no new surface.

## Red-Green Proof

**RED prompt/case:** A Go test asserting the future-stamp warning message names fabrication as a possible cause (alongside the wall-clock/Z-suffix cause) fails against the current strings.
**Why RED now:** All three diagnosis messages assert the timezone cause alone; a fabricated stamp — the observed incident — is misdiagnosed by the rendered warning.
**GREEN when:** The three messages name both causes with the fix instruction unchanged, the new assertion passes, and `go test ./...` in the tool directory exits 0.
**Validation:** Inferred during capture

## Full Context

See `do-work/user-requests/UR-055/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding 3 — "Broaden the board's future-stamp warning text: 'local wall-clock time with a Z suffix' is one cause; a fully fabricated value is a second, now-observed one, and the current message sends the reader to the wrong fix."*

---

# Builder Guardrails (orchestrator-issued — binding)

## Your tree

- Work **only** inside your worktree (path below). It is a full checkout on your own branch.
- **Never write anywhere in the main tree** except the single hand-back file named below. That is the only main-tree path you may touch.
- **Never touch `do-work/`** — not the queue, not `working/`, not `CHECKPOINT.md`, not `archive/`. Queue state is the orchestrator's alone. Your branch must contain **zero** commits touching `do-work/`; the orchestrator runs `git diff --name-only <pre>...<your-branch> -- do-work/` and a single path there stops your hand-back.
- **Never touch `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, or `CHANGELOG.md`.** Those are serial-only integrator-owned files. A bump on your branch races every sibling.
- **Scratch files go in `/tmp` or inside your worktree — never the main tree root.** A previous builder left a PNG in the repo root; that is a write-set violation. Screenshots, fixtures, generated boards: `/tmp`.

## Commit on your branch

Commit your implementation on your own branch before handing back. Message body only — no version bump, no changelog entry. End the message with:

```
Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
```

## The P-A-U loop is yours to fill

The REQ body contains an `## AI Execution State (P-A-U Loop)` section with three checkboxes, or the orchestrator will add one. **You must tick all three and write the required content into each**, in your worktree's copy of nothing — instead, put the filled P-A-U block **into your hand-back file** verbatim under a `## P-A-U` heading, since you may not write `do-work/`. `qualify.sh` FAILs on unticked boxes and the orchestrator will otherwise have to fill them from your evidence.

- **[PLAN]** — read the listed `prime_files` and agent rules, then write the technical approach. No code yet.
- **[APPLY]** — code written exactly as planned, scope strictly limited to planned files.
- **[UNIFY]** — run `git diff --stat`, review every changed file, run the project's linters/tests, confirm no debug artifacts. List each file you verified and what you checked.

## Evidence rules — every one of these was learned by getting it wrong

1. **Two REDs when the first is a reference error.** A test that fails because a constant or function does not exist yet proves nothing. Put the code in place, break exactly one rule, and let the assertion fail *for the reason it exists*. Report both RED outputs.
2. **`git stash push` on a clean file stashes nothing** — and the resulting green run reads as proof when it is vacuous. To reproduce RED against pre-change code, check out the pre-change blob by hash (`git show <hash>:<path>`) instead.
3. **Assert page identity inside the same call that reads the DOM.** If you drive a browser, return `location.href` (and, where relevant, the page's own rule text) from the *same* `evaluate` call as every measurement. A shared browser instance can be navigated by a sibling between your navigate and your evaluate, and the numbers come back confident, well-formed, and about somebody else's page. A URL checked *before* navigating is not the same claim. Prefer an isolated browser context.
4. **A programmatic `.focus()` does not trigger `:focus-visible` in Chrome.** Use a real `Tab` keypress if focus styling is in question.
5. **Generate the artifact and look at it.** For anything that changes what appears on screen, a passing assertion is not evidence about two glyphs sharing a coordinate. Measure `getBoundingClientRect()` intersections in the live DOM when the question is "do two things overlap"; read the rendered text when the question is "what does this say".
6. **Push back if the brief is wrong.** If a requirement contradicts an existing test, or a piece of code you wrote turns out unneeded, say so in the hand-back rather than quietly editing the test or keeping dead code. Two builders pushed back last session and both were right.

## Verification bar

`bash _dev/tests/maintainer-verify.sh` from your worktree root. **Exit code 0 is the only proof.** Never pipe it through `tail`/`head` — the pipeline's exit status hides the failure. Run it, then `echo $?` on its own line, and paste that.

## Hand-back

Write **one** file, at the absolute path given below, containing:

1. `## Branch` — your branch name.
2. `## P-A-U` — the three filled, ticked checkboxes with their content.
3. `## Files Changed` — `git diff --stat` against your branch's merge base, plus one line per file saying what changed and why.
4. `## Red-Green Evidence` — the RED output(s) and the GREEN output, quoted.
5. `## Verification` — the `maintainer-verify.sh` tail and its `echo $?` line.
6. `## Integration Seams` — anything the orchestrator must apply by hand in the merge commit (shared registries, cross-REQ text). Say "none" if none.
7. `## Decisions` — numbered D-01, D-02… for choices with reach beyond this REQ.
8. `## Lessons Learned` — what a future session should know. Omit if genuinely nothing.
9. `## Pushback` — anything in this brief you think is wrong. Omit if none.

Your final message back should be a short summary; the hand-back file is the real deliverable.

## Your Assignment

- **Worktree path (your working directory):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-245-name-fabricated-stamps-in-the-boards-future-stamp-warnings`
- **Branch name:** `worktree-agent-REQ-245-name-fabricated-stamps-in-the-boards-future-stamp-warnings`
- **Hand-back file (absolute, main tree — the ONE main-tree path you may write):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-18-124358/REQ-245-handback.md`
- **Repo root of the MAIN tree (read-only for you):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`

## Orchestrator Notes for This REQ

- Your `write_set` says `*_test.go`. **Narrow it yourself to exactly `timestamp_test.go` and `completion_anomaly_test.go`** (plus `model.go` and `verify.go`). A sibling builder owns `generate_test.go` in the next wave — do not touch it.
- This is a message-text change. No new checks, no threshold changes, the 2-minute skew allowance stays exactly as is.
- The three sibling messages must be updated **together** so they cannot drift: `model.go:379`, `model.go:1232`, `verify.go:371`. Also bring the single-cause comments at `timestamp_test.go:42`, `completion_anomaly_test.go:227`, `model.go:1338` in line where the new wording contradicts them.
- Keep the fix instruction (rewrite with the current UTC instant per the Timestamp rule) unchanged — it is correct for both causes.
- **Context for why this REQ exists:** last session three `completed_at` stamps were written by extrapolating the clock forward instead of reading it. The board's own future-stamp check caught them. The current message told that reader "your timezone is wrong", which was the wrong fix. That is the second cause you are naming.
