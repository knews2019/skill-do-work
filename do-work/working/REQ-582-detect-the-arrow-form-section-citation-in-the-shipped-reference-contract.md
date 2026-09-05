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
