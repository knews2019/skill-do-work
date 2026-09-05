# Builder brief — REQ-582

## Where you work

- **Your worktree (cd here first):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1201/worktree-agent-REQ-582-arrow-citations`
- **Your branch (already checked out there):** `worktree-agent-REQ-582-arrow-citations`
- **Route:** B
- **Base commit:** 961697bb

You are the builder. The orchestrator runs in the main checkout at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` and is the only writer of `do-work/`. Commit your work on your own branch in your own worktree and hand back a manifest; the orchestrator merges.

## Never touch

- Anything under `do-work/` — with exactly one exception, the hand-back file named below, which you write by its absolute main-tree path and never stage or commit.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION` — release paths owned by finalization.
- Any file outside the write set declared in the REQ below. If you need one, stop and report it in the hand-back instead of writing it, unless the REQ's own requirements already demand that file class (then flag the contradiction and proceed).
- Do not run `bash _dev/tests/maintainer-verify.sh` (the repository gate). The orchestrator owns it and concurrent runs corrupt each other's timing budgets. Run only the focused tests named below.

## Rules to load and follow (read these first, from your worktree)

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/shared-principles.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/testing.md` (the REQ is `tdd: true`)

Also read every path in the request's `prime_files`, and the `lessons-<name>.md` satellite beside each prime whose Read-first or Traps entries your change touches.

## P-A-U phasing (mandatory, reported in the hand-back)

The REQ file is the orchestrator's, so report your P-A-U record under a `## P-A-U` heading in the hand-back instead of ticking boxes in the REQ:
- **[PLAN]** — brief technical approach, written before code.
- **[APPLY]** — code exactly as planned, strictly inside the declared write set.
- **[UNIFY]** — run `git diff --stat`, run the native linters (`gofmt -l .`, `go vet ./...` for Go changes, `node --check` for changed client files), verify no debug artifacts in added lines, and list each file you checked and what you checked.

## Focused tests

Every test-file invocation must finish in under 30 seconds. Use:
- `bash _dev/tests/shipped-package-reference-contract.sh` (the check you are changing — run it before and after)
- `bash _dev/tests/action-shell-blocks.sh` (ShellCheck lint over shipped shell, which your edit must keep passing)

## Hand-back (write this file, then stop)

Write **`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-120117/REQ-582-handback.md`** using that absolute path — it is the one main-tree path you may write, and you must never stage or commit it.

It must contain, each under its own `##` heading:
- `## Branch` — the branch name and the head commit you left on it.
- `## File manifest` — every source file created/modified/deleted with the verb, plus tests touched.
- `## P-A-U` — the three phases above.
- `## Test evidence` — every command you ran, its exit status, the RED observation (test name + failure text) and the GREEN observation.
- `## Lesson evidence` — each lesson satellite you read and any listed path that was missing.
- `## Decisions` — significant choices as `D-NN`, each with reasoning. Mark a reversible low-reach choice DECIDE & STATE; mark an irreversible, taste-dependent or contestable one ESCALATE and add `Value:` and `Risk:` lines.
- `## Discovered Tasks` — out-of-scope findings, each stamped with one of exactly these impact tokens: `impact-critical`, `impact-user-visible`, `impact-rule-change`, `impact-negligible`. Do not invent a token outside that set and do not fix the items inline.
- `## Integration seams` — any exact line that belongs in a file outside your write set, with where it goes. The orchestrator applies it.
- `## Exploration` — what the sweep over all 75 arrow-form citations found, including every additional dangling citation. The orchestrator folds this into the request.


This is Route B and you owe an `## Exploration` section in the hand-back — the sweep over all 75 arrow-form citations is the exploration.

The real decision is the one the request names: does "section" mean an ATX heading only, or also a bold contract label that leads a paragraph? Both forms ship live. Settle it against what the 75 citations actually do, record it as a `D-NN` decision with reasoning, and expect the first honest pass to surface more dangling citations than the two already named. **Report the extra ones, do not fix prose outside this request's scope.**

`skills/do-work/CHANGELOG.md` is in your write set for one line only, and it needs care: a changelog entry describes what shipped on the day it shipped, so silently retargeting it at a renamed heading rewrites history. Read `_dev/primes/prime-releases.md` before touching it and say in your hand-back which remedy you chose and why. Note that the root `CHANGELOG.md` and `skills/do-work/CHANGELOG.md` are required to stay byte-identical by this very check — if your fix changes one, the other needs the same change, and the root copy is NOT in your write set, so hand it back as an integration seam.

You are editing shipped shell. Read `_dev/primes/prime-shell-commands.md` first; it is the request's own prime and it holds the trap list.

---

# The request

---
id: REQ-582
title: '[impact-rule-change] Detect the arrow-form section citation in the shipped-package reference contract'
status: claimed
created_at: 2026-09-05T01:30:57Z
user_request: UR-119
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
write_set: [_dev/tests/shipped-package-reference-contract.sh, skills/do-work/CHANGELOG.md, skills/do-work/actions/cleanup.md]
tdd: true
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
claimed_at: 2026-09-05T12:40:40Z
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-05T12:41:13Z
  basis:
    - Route B
    - 3-file write set
    - 4 acceptance criteria
route: B
dispatch_at: 2026-09-05T12:41:13Z
---

# Detect the Arrow-Form Section Citation in the Shipped-Package Reference Contract

## What
`_dev/tests/shipped-package-reference-contract.sh` resolves the path in a citation but never the section name written after it in the `` `path.md` `` → **Named Section** form. Two live shipped citations name sections that do not exist and the check passes. Make that form resolve, and fix the two dangling citations it then reports.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Every sweep that deletes or renames a section in a shipped action or reference file relies on this check to catch dangling inbound references. For the arrow form it reports clean regardless, so a sweep gets false confidence exactly where it most needs a real answer. It already happened: during REQ-510's RED run (sweeping work-reference sections whose contract is now a CLI behavior test) the check passed while two live callers named a heading that had just been deleted.

## Context
Claim verified before capture, on the main tree at commit `a55f24ce`. `bash _dev/tests/shipped-package-reference-contract.sh` exits 0 and prints `shipped package reference contract: PASS`, while both of these are dangling:

- `skills/do-work/CHANGELOG.md:279` cites `actions/work-reference.md` → **Recovery Refusals (Step 1)**. No such heading exists; the content now lives under `## Stuck Runs Hand Off to Judgment (any step)`.
- `skills/do-work/actions/cleanup.md:48` cites `actions/work-reference.md` → **In-Progress Record (Step 2)**. The only matching heading is `## In-Progress Record (Step 1)` at line 482.

Why it cannot see them: the script contains no occurrence of the arrow character and no bold-name parsing at all. `citation_shape` (`^(?:\.\./)*(?:[A-Za-z0-9._-]+/)+[A-Za-z0-9._-]*`) matches path-shaped tokens, and `heading_anchor_slugs` is reached only when the token itself carries a `#fragment`. The section name in the arrow form is separate bold text, so nothing looks at it.

75 arrow-form citations ship across the three skill packages, so the form is load-carrying prose, not an occasional flourish.

This is the same failure family as REQ-249, REQ-269 and REQ-312, which each widened this checker after a dangling citation escaped through a spelling the class did not cover. REQ-269 in particular fixed the class boundary for *which tokens are paths*. This one is the next axis: the *section name* beside the path.

## Red-Green Proof
**RED prompt/case:** On the current tree, run `bash _dev/tests/shipped-package-reference-contract.sh`.
**Why RED now:** It exits 0 while `skills/do-work/CHANGELOG.md:279` and `skills/do-work/actions/cleanup.md:48` both name sections of `actions/work-reference.md` that are not there. A citation form the checker cannot see is a citation form no sweep can trust.
**GREEN when:** The same command exits non-zero and reports both citations, each naming the citing file, its line, the resolved target and the missing section name. After the two citations are corrected it exits 0 again, and a deliberately planted third arrow-form citation to a non-existent section is reported the same way.
**Validation:** User confirmed — the maintainer asked for the claim to be verified before writing this request; the verification above is the result and it holds.

## Detailed Requirements
- Resolve the section name in the `` `path.md` `` → **Named Section** form against the target file, and report a name that is not present there.
- Add an in-script fixture case for the form, next to the existing `run_anchor_slug_fixtures` and `run_anchor_topology_fixtures` suites, so the new coverage is itself pinned.
- Fix the two live dangling citations named in Context.
- Keep the existing `#fragment` anchor resolution working unchanged.

## Constraints
- Do not narrow the citation class to close the gap. The rule is that a reader must be able to follow the citation; the fix widens what the checker sees, it does not shrink what shipped prose may write.
- Keep the check's runtime in its current range.

## Builder Guidance
The real decision is what "section" means here. Many arrow-form citations name an ATX heading (`## In-Progress Record (Step 1)`), but many name a bold contract label that is a paragraph lead and not a heading at all — `actions/work-reference.md` → **Frontmatter Quoting** and `actions/capture-reference.md` → **Populating `write_set`** are both live examples. Decide whether the checker resolves both forms or only headings, record the decision in this REQ, and expect a first pass over all 75 citations to surface more than the two already named. Report any additional dangling citations you find rather than silently fixing prose outside this REQ's scope.

The changelog citation may deserve a different remedy from the action-file one: a historical entry describes what shipped at the time, so retargeting it at the renamed heading and correcting the entry are not the same choice. `_dev/primes/prime-releases.md` holds the changelog house rules; read it before editing `skills/do-work/CHANGELOG.md`.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget in `actions/capture-reference.md` → Required Lessons Budget Contract, and `slugged: partial` so no targeted entry is legal. Matches through its owning prime, `_dev/primes/prime-shell-commands.md`, which governs the shell this REQ rewrites.

## Full Context
See `do-work/user-requests/UR-119/input.md` for complete verbatim input.

*Source: discovered task from the REQ-510 builder (sweeping work-reference sections whose contract is now a CLI behavior test), work run `do-work/runs/work-2026-09-05-003420/`.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear and the two dangling citations are named, but the real decision — whether a "section" means an ATX heading only, or also a bold contract label that leads a paragraph — has to be settled against the 75 arrow-form citations that actually ship. That is discovery.

**Planning:** Not required — the requirements are already a design.

**Exploration:** Delegated to the builder rather than a separate agent, because the discovery here *is* the first sweep over all 75 citations, which the builder has to run anyway to implement the check. The builder writes `## Exploration` into its hand-back and the orchestrator folds it in. Recorded as an orchestrator judgment so the skipped step is visible rather than silent.

## Plan

**Planning not required** - Route B: exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `_dev/tests/shipped-package-reference-contract.sh` (modify) — resolve the section name in the arrow form, plus an in-script fixture suite pinning the new coverage
- `skills/do-work/CHANGELOG.md` (modify) — the dangling citation at line 279, remedied per the changelog house rules
- `skills/do-work/actions/cleanup.md` (modify) — the dangling citation at line 48

**Files I will NOT touch:** any other shipped prose carrying an arrow-form citation — additional dangling citations are reported, not fixed; `VERSION` and the changelog mirror (release paths written by finalization, except the cited line itself which is this request's subject).

**Acceptance criteria (restated from REQ):**
- [ ] The arrow form resolves its section name against the target file and reports a name that is not there, naming citing file, line, resolved target and missing section
- [ ] An in-script fixture case pins the new coverage beside the existing anchor suites
- [ ] Both live dangling citations are corrected, the changelog one under the house rules for historical entries
- [ ] Existing `#fragment` anchor resolution is unchanged and the check's runtime stays in its current range
