---
id: REQ-582
title: '[impact-rule-change] Detect the arrow-form section citation in the shipped-package reference contract'
status: completed
created_at: 2026-09-05T01:30:57Z
user_request: UR-119
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
write_set: [_dev/tests/shipped-package-reference-contract.sh, skills/do-work/CHANGELOG.md, skills/do-work/actions/cleanup.md, CHANGELOG.md]
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
builder_handback_at: 2026-09-05T13:19:20Z
integration_at: 2026-09-05T13:19:20Z
review_at: 2026-09-05T14:25:35Z
heavy_verified_at: 2026-09-05T14:25:35Z
heavy_verified_revision: 7b2673b690a671ccb360c26b0c19c56ecc7356b5
commit: 7b2673b690a671ccb360c26b0c19c56ecc7356b5
completed_at: 2026-09-05T14:25:54Z
release_at: 2026-09-05T14:25:54Z
---

# Detect the Arrow-Form Section Citation in the Shipped-Package Reference Contract

## What
`_dev/tests/shipped-package-reference-contract.sh` resolves the path in a citation but never the section name written after it in the `` `path.md` `` → **Named Section** form. Two live shipped citations name sections that do not exist and the check passes. Make that form resolve, and fix the two dangling citations it then reports.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Builder read the shell prime and its lesson satellite, swept all 75 live arrow-form citations, and settled the definition of a cited section against what the corpus does rather than what was cheapest to implement. Recorded under `## P-A-U` in `do-work/runs/work-2026-09-05-120117/REQ-582-handback.md`.
- [x] **[APPLY]:** Two commits on the builder branch (`04f02f10`, `87956ed1`) — the check with its fixtures and the two live fixes, then one false-positive shape the sweep exposed, narrowed and pinned.
- [x] **[UNIFY]:** `git diff --stat` reviewed; ShellCheck clean over the edited script; the check run before and after; the root changelog deliberately left untouched and handed back as a seam, because the same check enforces that the two copies stay byte-identical.

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
- `CHANGELOG.md` (modify) — the same corrected line, re-mirrored by the orchestrator inside the merge commit

**Files I will NOT touch:** any other shipped prose carrying an arrow-form citation — additional dangling citations are reported, not fixed; `VERSION` and the changelog mirror (release paths written by finalization, except the cited line itself which is this request's subject).

**Acceptance criteria (restated from REQ):**
- [ ] The arrow form resolves its section name against the target file and reports a name that is not there, naming citing file, line, resolved target and missing section
- [ ] An in-script fixture case pins the new coverage beside the existing anchor suites
- [ ] Both live dangling citations are corrected, the changelog one under the house rules for historical entries
- [ ] Existing `#fragment` anchor resolution is unchanged and the check's runtime stays in its current range

## Exploration

The builder produced this section in its hand-back (`do-work/runs/work-2026-09-05-120117/REQ-582-handback.md` → `## Exploration`) in place of a separate Explore agent, per the Triage note above: the sweep over all 75 live arrow-form citations is both the exploration and the evidence the design decision rests on.

## Implementation Summary

**Files changed:**
- `_dev/tests/shipped-package-reference-contract.sh` (modified)
- `skills/do-work/actions/cleanup.md` (modified)
- `skills/do-work/CHANGELOG.md` (modified)
- `CHANGELOG.md` (modified)

**What was done:** The check now resolves the section name written after the path in the arrow form, not just the path, and reports the citing file, its line, the resolved target and the missing name. A cited section is an ATX heading or a bold label the target declares, matched by whole-word containment — both choices settled against the live corpus rather than chosen for convenience. An in-script fixture suite pins the new coverage beside the existing anchor suites. The two dangling citations named in the request are fixed: the action file now names the section that exists, and the changelog entry keeps its historical name as prose while gaining a live pointer to where that content moved. The root changelog was re-mirrored by the orchestrator inside the merge commit, because the same check enforces that the two copies stay byte-identical and it is outside the builder's write set. Merge range `43be2af6..7b2673b6`; builder branch head `87956ed1`.

## Qualification

**Passed.** Read from the merge range `43be2af6..7b2673b6`, and the check itself re-run on the merged tree: `bash _dev/tests/shipped-package-reference-contract.sh` exits 0 with `PASS`.

- The class was widened, not narrowed, which is the request's explicit constraint. A heading-only rule would have reported five correct, deliberately-declared citation targets as broken — the request's own reference files name bold contract labels and say in their own text that other files cite them by name.
- Whole-word containment rather than equality is carried by real shapes: a citation may name a heading without its parenthetical qualifier, or a label without its declaration prefix. Containment still catches the case this check exists for — the dangling `In-Progress Record (Step 2)` is not contained in `In-Progress Record (Step 1)`.
- The bold-label allowance is a false *pass* risk, never a false failure, and the builder said so plainly rather than presenting it as free.
- The changelog remedy respects the house rule that history is not rewritten: the old section name stays as prose and a live pointer is added beside it, instead of retargeting the entry as if it had always said something else.

Requirements traced: the arrow form resolves and reports; an in-script fixture suite pins it; both live dangling citations fixed; existing anchor resolution unchanged; runtime unchanged.

*Checked by work action*

## Testing

**The check itself, post-merge on the integrated tree:** `bash _dev/tests/shipped-package-reference-contract.sh` — exit 0, `shipped package reference contract: PASS`.

**Red-green validation** (traced to `## Red-Green Proof`): RED was the state the request describes and the builder confirmed — the check exited 0 while two shipped citations named sections that were not there. GREEN is the same command reporting both, then exiting 0 once they are corrected, with a deliberately planted third dangling arrow-form citation reported the same way. The in-script fixture suite is what keeps that coverage from rotting.

**One thing the merge had to carry:** until the root changelog was re-mirrored, this very check failed with `changelog mirror differs` and would have taken the repository gate down with it. The builder identified that precisely and refused to reach outside its write set to fix it.

## Decisions

- **D-07 — the root changelog was added to this request's write set by the orchestrator, before integration. DECIDE & STATE.** The builder correctly refused to touch it: it is outside the declared write set, and the check this request changes is the very thing that enforces the two changelog copies staying byte-identical. It handed the exact replacement line back as an integration seam instead. Applying that seam inside the merge commit is the orchestrator's job, and the Scope list and `write_set` were extended to match rather than leaving the touch undeclared. Without it, this request's own check fails with `changelog mirror differs` and takes the repository gate down with it.


## Review

**Overall: 90%** | 2026-09-05T13:29:22Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 85% |
| Test Adequacy | 75% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- D-01's bold-label acceptance is wider than its stated risk: any bold run anywhere in the target counts, and these files habitually write their own section names in bold prose. `skills/do-work/actions/work-reference.md` declares 42 headings but yields 288 accepted names, 246 of them bold-only. Measured: renaming the heading `## Worktree Dispatch Mode (Step 1)`, `## Crash Recovery (Step 1)` or `## In-Progress Record (Step 1)` and leaving every citation on the old name keeps the check GREEN in all three cases, because each old name is still written in bold somewhere in that file. Separately, a citation `` `actions/work-reference.md` `` → **Input** resolves against the bold phrase "advisory input to that pick, never a gate" while no section named Input exists. Verified remedy, not applied here: accepting only a *paragraph-leading* bold run keeps all 74 live arrow citations green (zero false failures) and makes all three simulated renames report — impact-rule-change → report only
- Nothing pins the wiring between the two checked halves and a reported finding. Mutating `section_failure_messages` to `return []`, or passing `cited_section=None` at the driver call site (`_dev/tests/shipped-package-reference-contract.sh:2165`), leaves the entire suite green at exit 0 — the same "control that cannot fail" class as REQ-581 in this batch. Both halves themselves are pinned and do fail when mutated. `run_citation_fixtures` already calls `citation_messages` and would take one case carrying `cited_section` + `read_section_names` — impact-rule-change → report only
- Two live shipped citations write the same form with the ASCII arrow and stay invisible: `skills/do-work/actions/forensics.md:54` and `:100` both cite `` `actions/work-reference.md` `` -> **Crash Recovery (Step 1)**. That heading exists (`work-reference.md:302`), so nothing dangles today, but a rename of it escapes silently — the exact failure family this REQ closes. The hand-back census counted `→` only, so this is absent from both the Exploration table and DT-1..DT-4 — impact-rule-change → report only

**Minor findings:**
- The section half checks only the source topology (`resolution_pair[2]` / `resolution_result[4]`), while the adjacent anchor half checks both through `link_topology_targets` + `anchor_failure_messages` and pins the divergence in `run_anchor_topology_fixtures`. The live corpus resolves both topologies to one file every time, so there is no divergence today — impact-rule-change → report only
- Restatement Sweep: `_dev/primes/prime-action-files.md:93` is the canonical home for how to write a cross-reference and still describes the checker as enforcing the path/citation-class condition only; it says nothing about the section name now being resolved. `do-work/prose-backlog.md` already carries a REQ-312-era staleness line against the same sentence — impact-rule-change → report only
- D-04 makes a historical changelog entry's section name permanently checked, so a future section rename can be forced to edit shipped history again. Only one arrow-form citation lives in a changelog today, so the standing cost is near zero — impact-negligible → report only

**Restatement Sweep:** redefined elements were (a) what a citation must satisfy in this contract and (b) the return arity of `cross_package_target_resolves` / `same_package_target_resolves`. Swept: `_dev/primes/prime-action-files.md:93` (stale, above), `_dev/primes/prime-releases.md:9` (mirror rule, still true), `_dev/tests/contracts/probe-lanes.sh:19` (keys on exit status only, unaffected), both resolvers' only consumer is `citation_messages` (grep-verified), and `skills/do-work-board/tools/queue-kanban/citations.go` parses REQ-text citations under a separate contract, unaffected. **No additional dangling citation anywhere in shipped prose:** the check exits 0, and an independent sweep of all 74 live arrow-form citations found 0 dangling and exactly 5 that require the bold-label rule — the same 5 D-01 names.

**Acceptance:** Pass — check green (exit 0, `PASS`); mutating either checked half makes the new fixtures fail (exit 1); end-to-end reporting verified through `citation_messages` with real driver arguments; `bash _dev/tests/action-shell-blocks.sh` exit 0; runtime unchanged (5 runs each: 1.16–1.37 s user before, 1.17–1.26 s user after).
**Suggested testing:** 3 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Settling the design question against the live corpus instead of against what was easiest to implement. The builder swept all 75 arrow-form citations before choosing what a cited section means, and found five deliberately-declared bold contract labels that a heading-only rule would have reported as broken. That is the difference between a check that answers honestly and one that is merely strict.

**What didn't:** The widening it chose is wider than the risk it stated. The reviewer measured it: one reference file declares 42 headings but yields 288 accepted names, 246 of them bold-only, and renaming three real headings leaves the check green because each old name is still written in bold prose somewhere in the same file. So the failure family this request exists to close — a rename escaping — is narrowed, not closed. The reviewer also verified the tighter rule: accepting only a *paragraph-leading* bold run keeps all 74 live citations green and makes all three simulated renames report.

**Worth knowing:** The new check has both halves pinned, and yet mutating the driver call site that joins them leaves the whole suite green at exit 0. That is the same "a control that cannot fail" defect this batch's other request was raised for, in the batch that fixed it. Pinning the pieces is not pinning the wiring, and the wiring is where a check quietly stops running.

## Orientation

The shipped-package reference contract now reads the section name in a citation, not just the path, so a sweep that renames or deletes a section in shipped prose gets a real answer instead of a clean one. Lives in the maintainer gate's shipped-reference contract, governed by `_dev/primes/prime-shell-commands.md`. [MAP CHANGED] — what a citation must satisfy has widened from "the path resolves" to "the path resolves and the named section is findable in it", which every future sweep over shipped prose now depends on. Two live dangling citations were fixed in passing.

## Heavy Verification Result

- **Target revision:** 7b2673b690a671ccb360c26b0c19c56ecc7356b5
- **Execution revision:** 7b2673b690a671ccb360c26b0c19c56ecc7356b5
- **Run at:** 2026-09-05T14:25:35Z, from a detached worktree

| Lane | Exit | Wall | Disposition |
| --- | --- | --- | --- |
| `queue-kanban-javascript` | 0 | 9s | executed |
| `queue-kanban-browser` | 0 | 141s | executed |
| `do-work-cli-integrations` | 0 | 78s | executed |
| `staged-skills` | 0 | 44s | executed |
| `updater` | 0 | 71s | executed |
| `installer` | 0 | 33s | executed |

All six lanes this request selected were present in the run, exited 0, and none was skipped.

