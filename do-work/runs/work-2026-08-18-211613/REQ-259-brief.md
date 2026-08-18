# REQ-259 builder brief

**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-259-retire-the-skill-root-reading-outside-backtick-spans` — this is your working directory; everything you do happens here.
**Branch:** `worktree-agent-REQ-259-retire-the-skill-root-reading-outside-backtick-spans` (already checked out in that worktree).
**Hand-back file (write when done):** `/home/user/skill-do-work/do-work/runs/work-2026-08-18-211613/REQ-259-handback.md` — absolute main-tree path, the ONE main-tree path you may write. Never stage or commit it.

## The REQ (verbatim)

```markdown
---
id: REQ-259
title: Retire the skill-root citation reading at its three unbackticked sites
status: claimed
claimed_at: 2026-08-18T21:16:24Z
route: B
created_at: 2026-08-18T18:07:48Z
status_changed_at: 2026-08-18T20:59:31Z
user_request: UR-055
addendum_to: REQ-249
domain: general
review_generated: true
sweep: true
sweep_key: retired-reading-outside-backtick-spans
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- skills/do-work/SKILL.md
- skills/do-work/actions/commit.md
- skills/do-work/crew-members/security.md
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T21:17:28Z
  basis:
    - Route B
    - 3-file write set
    - 3 subsystems involved
    - 3 acceptance criteria
    - full-suite verification
---

# Retire the Skill-Root Citation Reading at Its Three Unbackticked Sites

## What

REQ-249 retired the skill-root-relative citation reading and swept every **backticked** cross-package citation to the literal form — but the retired *reading* survives at three shipped sites that are bare text, which both the sweep and the new checker are structurally blind to. One of them states the retired resolution rule as prose in the core router, now contradicting the prime and the swept corpus.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Found by REQ-249's independent review (Important, gate: rule-change; reproduced by execution — paths verified non-resolving, checker verified passing over them). Not builder scope drift: the user-confirmed rule's letter covers backticked citations, and these are bare text — but the diff retired the reading, and these sites still state or use it. Created `pending-answers` per the generation-≥2 cascade stop (REQ-249 is itself review-generated).

## Instances

- [ ] **`skills/do-work/SKILL.md:16`** — "When a core action names a sibling path, resolve it from the parent directory containing these four skill roots" is the retired resolution rule stated as prose in the core router; applied to the new `../../` form it computes wrong paths. Restate to match the literal rule in `_dev/primes/prime-action-files.md` § Cross-Referencing.
- [ ] **`skills/do-work/actions/commit.md:17`** — bare `../do-work-toolbox/actions/inspect.md`, wrong depth from `actions/`.
- [ ] **`skills/do-work/crew-members/security.md:3`** — bare `../do-work-toolbox/actions/code-review.md` in the JIT_CONTEXT comment; the checker's comment scan covers backticked spans only, so this evades.

## Requirements

- No shipped text states or uses the retired skill-root-relative reading outside REQ-249's documented exemptions (fenced template/example blocks).
- Whether the two bare paths become backticked (and thereby checkable) is builder latitude; the SKILL.md prose must agree with the prime either way.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Open Questions

- [x] REQ-249's review found the retired citation reading surviving at three unbackticked shipped sites (listed under Instances). Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

**Answered [2026-08-18]:** User approved via `do-work clarify`. An external automated review independently flagged the same three sites while this question was open, which corroborates the finding.

---

## Triage

**Route: B** - Medium

**Reasoning:** Three named instances, but the requirement is a class claim (`no shipped text states or uses the retired reading`), so the builder must sweep the corpus for further bare-text sites the backtick-scoped checker is blind to before declaring it closed.

**Planning:** Not required

---

## Exploration

All three instances reproduce exactly as the REQ states:

- `skills/do-work/SKILL.md:16` — "When a core action names a sibling path, resolve it from the parent directory containing these four skill roots. Do not search the core package for an extension action." The first sentence is the retired reading; the second is a separate, still-true instruction.
- `skills/do-work/actions/commit.md:17` — `route to ../do-work-toolbox/actions/inspect.md instead`, bare, wrong depth from `actions/` (needs `../../`).
- `skills/do-work/crew-members/security.md:3` — `../do-work-toolbox/actions/code-review.md` inside the JIT_CONTEXT comment; that same comment already spells a same-package citation in backticks, so the file mixes both forms.

The rule to agree with is `_dev/primes/prime-action-files.md` § Cross-Referencing (line 91): literal relative path from the citing file's own directory, depth per-file, fenced template/example blocks exempt, and `_dev/tests/shipped-package-reference-contract.sh` checks the backticked ones in both topologies.

*Explored inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work/SKILL.md` (modify) — restate line 16's resolution prose to match the prime's literal rule
- `skills/do-work/actions/commit.md` (modify) — line 17's bare sibling path
- `skills/do-work/crew-members/security.md` (modify) — line 3's bare sibling path in the JIT_CONTEXT comment

**Files I will NOT touch:** `_dev/primes/prime-action-files.md` (it is the rule this agrees with, not a site to change), `_dev/tests/shipped-package-reference-contract.sh` (widening the checker to bare text is its own REQ, not this one — record it as a Discovered Task if the sweep argues for it), any other package's files.

**Acceptance criteria (restated from REQ):**
- [ ] No shipped text states or uses the retired skill-root-relative reading outside REQ-249's fenced template/example exemptions — established by sweeping the corpus for the primitive, not by fixing the three listed instances
- [ ] `skills/do-work/SKILL.md`'s prose agrees with `_dev/primes/prime-action-files.md` § Cross-Referencing
- [ ] Any site found beyond the three is either fixed here or reported with its reason
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0
```

## Never touch

- Anything under `do-work/` — the queue, `working/`, `CHECKPOINT.md`, `runs/` — is the orchestrator's. Your worktree carries a stale committed snapshot of it; treat it as absent. The single exception is your own hand-back file at the absolute path above, which you write but never stage or commit.
- `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md` — serial-only files owned by the integrator. Bumping any of them races every sibling builder. Skip the "Before Every Commit" ritual in `CLAUDE.md` entirely; it belongs to the integrating commit, not to yours.
- Any file outside your REQ's `## Scope` "Files I will touch" list. Discovering mid-build that you need one is a **stop-and-report to the orchestrator, never a silent write** — put it in the hand-back and say why.
- The main tree at `/home/user/skill-do-work` — you are not in it and must not write to it.

## Environment

- Go 1.26.1 at `/usr/local/go`, ShellCheck 0.11.0, `just` 1.21.0, Node 22 are all installed and are the versions the gate pins.
- **Never run a bare `go build`** in `skills/do-work-board/tools/queue-kanban/` — it drops an 11 MB gitignored binary into the source tree. Build to a scratch directory with `-o`.
- The canonical gate is `bash _dev/tests/maintainer-verify.sh`, run from your worktree root. **Exit code zero is the only proof it passed** — never pipe it through `tail` or `head`, because the pipeline's exit status hides the failure. Redirect to a file and echo `$?` if you want to read the output.
- Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp anything. Never carry a timestamp forward and never compute one.

## How to work

Follow the do-work crew rules that always load during implementation — read them from your worktree:

- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/communication-style.md`
- plus `skills/do-work/crew-members/testing.md` if this REQ has `tdd: true`, and `skills/do-work/crew-members/maintenance.md` if it has `maintenance: true`.

Read every path in the REQ's `prime_files` before touching code, including its `## Lessons` links.

**P-A-U is mandatory and it is your own note-taking** — the REQ file itself lives in the main tree and is not yours to edit, so write the [PLAN] / [APPLY] / [UNIFY] evidence into your hand-back instead. [PLAN] before any code. [APPLY] stays inside the declared scope. [UNIFY] runs `git diff --stat` against your branch point, runs the native linters, verifies no debug artifacts, and lists each file you checked.

**The failure mode that beat nine of twelve REQs last session:** a mechanism that looks like it closes a class and closes only the instance. The two REQs that beat it grepped or fuzzed the **primitive** before declaring the class closed, and both found the real hole where no instance list pointed. Assume your first fix has that shape. An instance list is a sample, not the class.

**Commit in small, individually-green increments on your branch.** Builders were repeatedly killed mid-run by server-side errors last session and nothing was lost, because each increment was already committed. Do not wait until the end to commit.

## Hand-back format

Write `REQ-NNN-handback.md` at the absolute path given above, with these sections:

- `**Branch:**` the operative name, and the commit it was cut from
- `**Commits (oldest first):**` short hash + subject for each
- `## What I built` — factual, what you actually built, not what the REQ asked for
- `## File manifest` — every file created/modified/deleted with `(new)` / `(modified)` / `(deleted)` and a phrase on what changed in it
- `## P-A-U evidence` — [PLAN] / [APPLY] / [UNIFY], with what you actually ran
- `## Testing evidence` — real observed output, never a prototype or a paraphrase. If the REQ is `tdd: true`, the RED must be a real run against pre-change code, quoted. State the gate's exit code as a number you observed.
- `## Decisions (D-XX)` — each marked DECIDE & STATE or ESCALATE; ESCALATE entries carry `Value:` and `Risk:` lines
- `## Integration seams` — exact lines and where they go, for any shared file you must not edit yourself. `None.` if there are none.
- `## Discovered Tasks` — out-of-scope finds, not fixed inline
- `## Pushback` — where the REQ was wrong, or where you disagree. `None.` if you have none.

Your report is a claim, not evidence: the orchestrator judges from git state. Say plainly what you did not do.
