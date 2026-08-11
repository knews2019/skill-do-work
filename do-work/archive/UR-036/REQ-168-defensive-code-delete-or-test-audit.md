---
id: REQ-168
title: Delete-or-test audit of defensive code in shipped skills
status: completed
created_at: 2026-08-11T11:46:50Z
user_request: UR-036
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-165]
maintenance: true
related: [REQ-165, REQ-166, REQ-167]
batch: stabilization-audit
claimed_at: 2026-08-11T12:46:09Z
route: C
write_set: [decisions/audits/2026-08-11-defensive-surface.md, _dev/tests/defensive-surface-audit.sh, _dev/tests/contract-regressions.sh, skills/do-work/actions/commit.md, skills/do-work/actions/verify-requests.md, skills/do-work-toolbox/actions/code-review.md, skills/do-work-toolbox/actions/inspect.md, skills/do-work-toolbox/actions/quick-wins.md, skills/do-work-toolbox/actions/ui-review.md]
completed_at: 2026-08-11T12:55:57Z
commit:
kb_status: pending
kb_entry:
---

# Delete-or-Test Audit of Defensive Code in Shipped Skills

## What

Audit every defensive layer in the shipped `skills/` tree — fallbacks, guards, workarounds, retry/recovery blocks, and warning apparatus in both shell (hooks, prescribed blocks) and prose (Rules/Rationalizations sections that restate hygiene) — and disposition each one: **keep** (traces to a named incident AND is covered by a test), or **delete** (can't name the incident it prevents, or its cost now exceeds the surface it protects).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Inventory executable guards and earned prose sections separately; trace keeps to incident-bearing tests, delete only generic restatement tables/rows whose underlying action contract already exists in Steps/verification, record every disposition in one durable audit, and ratchet the deliberate deletions plus audit coverage through the aggregate suite.
- [x] **[APPLY]:** Recorded executable and prose dispositions, deleted the six planned generic defensive surfaces without changing their underlying action contracts, and wired the coverage/deletion ratchet into aggregate contracts.
- [x] **[UNIFY]:** Reviewed the complete diff/stat and all 9 scoped files; verified every current explicit defensive heading and all 15 shipped shell sources map into the audit, all retained action steps/checklists/rules remain intact, planned generic rows are absent, the probe is executable and ShellCheck-clean, and no debug/whitespace artifacts entered the diff.

## Why (if provided)

User: "many things got more complex than needed." Untested defensive code is negative-value — it adds review surface without adding safety; the session-start hook proved robustness machinery can itself be the bug (46 lines of defense around 2 lines of logic, and the defense was the defect). Deleted lines can't come back as review findings.

## Context

- `maintenance: true` — removal/narrowing pass on the skill's own instructions; `crew-members/maintenance.md` (delete-before-you-add) loads.
- The audit question per layer, verbatim from the discussed plan: "what incident earned this, and is the fix still cheaper than the surface it added?"
- Known dispositions from the conversation: the hook's path-anchoring comment **keeps** (real regression behind it); the hook's dead fallback is REQ-166's scope, not this REQ's. CLAUDE.md's trap-list entries each trace to real incidents — they're the *standard* for "earned," not deletion targets (and CLAUDE.md is repo-only anyway; this audit's surface is the shipped tree).
- The earned-section test in CLAUDE.md ("can I name the specific failure this row prevents, and where it happened?") is the same rubric — apply it to shipped Rules / Common Rationalizations / Red Flags sections too, per the existing convention that generic tables are worse than none.
- Depends on REQ-165: "keep" dispositions in shell rely on the harness (or a targeted `_dev/tests/` case) for their test coverage.

## Detailed Requirements

- Produce the audit artifact: each defensive layer → location → incident it traces to (or "none found") → disposition → evidence (test name for keeps, diff for deletes).
- Apply deletions in this REQ; deletions that would change *observable* behavior of an action (not just its robustness theater) get flagged as follow-up candidates instead.
- Surviving shell fallbacks must be exercised by REQ-165's harness or a targeted test — an unreached fallback is presumed dead until proven live (the hook bug is the precedent).
- Serial-only files (`CHANGELOG.md`, `actions/version.md`) stay untouched by the builder; integrator owns them.

## Builder Guidance

Certainty: Firm on the rubric, exploratory on the inventory. Scope cue: bias toward deletion when the incident can't be named — that's the point — but log every delete in the artifact so review can restore cheaply. Not a rewrite pass: layers that survive stay textually as-is unless a test exposes them as broken.

## Red-Green Proof

**RED prompt/case:** Pick any shipped fallback (exemplar: session-start.sh's pre-fix "unknown" path) and ask (a) which incident earned it, (b) which test exercises it. Today the general answer to (b) is "none," and for some layers the answer to (a) is also "none."
**Why RED now:** Defensive layers accreted one reasonable fix at a time with no test-or-delete discipline, so reviews keep harvesting them as findings.
**GREEN when:** The audit artifact exists covering the shipped tree; every surviving defensive layer names its incident and its test; deletions are applied and committed; a subsequent review pass over the audited surface yields no findings of the "untested/unearned defensive code" class.
**Validation:** Inferred during capture (plan discussed and endorsed in-session).

## Full Context

See `do-work/user-requests/UR-036/input.md` for complete verbatim input.

---
*Source: "do-work capture-request for audit and fix to simplify and make it robust" (UR-036)*

## Triage

**Route: C** - Complex

**Reasoning:** The request spans all four shipped packages and requires judgment across executable fallbacks, action boundaries, safety guidance, historical incidents, and multiple test seams; deletion must be separated from observable behavior changes.

**Planning:** Required

## Plan

1. Inventory all 15 shipped shell sources and every shipped Rules/Common Rationalizations/Red Flags/Warnings/Recovery surface, grouping only when locations share the same incident and test seam.
2. Retain runtime guards that trace to recorded corruption, recovery, credential, install, parsing, or cross-shell incidents and name their behavioral/contract coverage; rely on REQ-165 for syntax reachability of every surviving shell source and fence.
3. Delete decorative review-hygiene sections from code-review, quick-wins, ui-review, and verify-requests where every row fails the named-incident test and duplicates the action's steps/checklist. Delete generic duplicate rows from commit and inspect while preserving their incident-specific secret, status-alias, and committed-state warnings.
4. Add a focused audit ratchet that requires the durable inventory, proves those unearned sections/rows stay deleted, confirms all shipped shell files remain in the REQ-165 harness surface, and invoke it from aggregate contracts.
5. Run the focused ratchet, shell-fence harness, scope/qualification checks, stale generic-surface sweep, and full aggregate suite; review the complete subtraction for any lost action behavior.

**Plan validation:** The inventory and artifact satisfy the complete-audit requirement; keep/delete criteria map directly to the user rubric; deletions are limited to redundant warning apparatus, with behavior-changing candidates recorded rather than implemented; the ratchet and full suite cover the requested no-regression proof.

*Generated by Plan phase after repository inspection*

## Exploration

- The shipped tree contains 15 shell files. All are parsed by `_dev/tests/action-shell-blocks.sh`; behavioral seams additionally exist for SessionStart, memory capture/startup, uncommitted inventory, association, preflight, commit-hash recovery, blanked-REQ recovery, suite manifests/install/update, and atomic section replacement.
- The largest executable guard surface is `install-do-work-suite.sh`, but its rollback, path validation, settings reconciliation, and post-write checks are exercised by `install-suite-behavior.sh`, `suite-manifest-contract.sh`, `staged-skills-contract.sh`, and aggregate fixtures. This is expensive but earned by named partial-install/index/settings incidents, so it is not a deletion target here.
- Action-specific rationalizations around archive collisions, X/XD secret quarantine, destructive memory/KB operations, local-vs-global installs, worktree cleanup, and prompt injection name concrete failure shapes and have adjacent contract or behavioral coverage.
- Four general review actions contain decorative Common Rationalizations and Red Flags tables made entirely of generic thoroughness advice already encoded in their steps, output schemas, and verification checklists: `code-review.md`, `quick-wins.md`, `ui-review.md`, and `verify-requests.md`.
- `commit.md` carries an arbitrary “>20 files” Red Flag with no incident and a generic two-row rationalization table; semantic grouping already owns that behavior. `inspect.md` similarly repeats generic “actually inspect / trace / catch debug” hygiene around three earned state-specific warnings.
- Rules sections and prompt-kit Red Flags are not blanket-deletion candidates: most define permission boundaries, schemas, destructive-write gates, or untrusted-input behavior. Removing those would change observable action behavior and belongs in a separately captured behavior-change REQ if review later challenges one.

*Generated by Exploration phase*

## Scope

**Files I will touch:**
- `decisions/audits/2026-08-11-defensive-surface.md` (new) — complete executable/prose defensive-layer inventory and dispositions
- `_dev/tests/defensive-surface-audit.sh` (new) — keep/delete audit ratchet
- `_dev/tests/contract-regressions.sh` (modified) — invoke the focused audit probe
- `skills/do-work/actions/commit.md` (modified) — delete generic duplicate rationale and arbitrary file-count warning
- `skills/do-work/actions/verify-requests.md` (modified) — delete decorative generic Rationalizations/Red Flags sections
- `skills/do-work-toolbox/actions/code-review.md` (modified) — delete decorative generic Rationalizations/Red Flags sections
- `skills/do-work-toolbox/actions/inspect.md` (modified) — delete generic duplicate rationale/warnings while retaining state-specific warnings
- `skills/do-work-toolbox/actions/quick-wins.md` (modified) — delete decorative generic Rationalizations/Red Flags sections
- `skills/do-work-toolbox/actions/ui-review.md` (modified) — delete decorative generic Rationalizations/Red Flags sections

**Files I will NOT touch:** executable shell/Go behavior, SessionStart (REQ-166), the canonical shell guide (REQ-167), serial-owned version/changelog files during implementation, or action semantics whose removal would change observable behavior

**Acceptance criteria (restated from REQ):**
- [x] The audit covers every shipped executable guard source and prose defensive section with location, incident (or none), disposition, and test/diff evidence.
- [x] Every surviving executable fallback names an incident and is reached by the REQ-165 harness or a targeted behavioral/contract test.
- [x] Every unearned redundant layer identified for safe deletion is removed without altering the underlying action steps, output, or permissions.
- [x] Potential behavior-changing deletions are retained and explicitly classified for follow-up rather than silently changed.
- [x] A focused ratchet plus the full contracts pass, and review finds no remaining untested/unearned defensive layer in the audited surface.

## Decisions

- **D-01 — Audit cohesive layers, not individual branches.** A transaction such as suite installation contains many cooperating guards but represents one recovery layer; counting every `if` separately would obscure the incident and test that justify the whole transaction.
- **D-02 — Delete only redundant apparatus in this REQ.** Rules that define permissions, schemas, destructive-write gates, or security outcomes were retained because removal changes observable behavior; generic review advice duplicated by existing steps/checklists was deleted now.
- **D-03 — Semantic grouping wins over a file-count heuristic.** The `>20 files` commit warning had no named incident and could incorrectly split one atomic generated/schema change; the existing one-logical-change rule is the cheaper and stronger contract.

## Implementation Summary

**Files changed:**
- `decisions/audits/2026-08-11-defensive-surface.md` (new)
- `_dev/tests/defensive-surface-audit.sh` (new)
- `_dev/tests/contract-regressions.sh` (modified)
- `skills/do-work/actions/commit.md` (modified)
- `skills/do-work/actions/verify-requests.md` (modified)
- `skills/do-work-toolbox/actions/code-review.md` (modified)
- `skills/do-work-toolbox/actions/inspect.md` (modified)
- `skills/do-work-toolbox/actions/quick-wins.md` (modified)
- `skills/do-work-toolbox/actions/ui-review.md` (modified)

**What was done:** Audited the shipped executable and explicit prose defensive surface, tying retained layers to incidents and test evidence while classifying behavior-changing candidates separately. Removed four complete decorative Rationalizations/Red Flags pairs, two duplicate Rationalizations sections, three generic inspect warnings, and one arbitrary commit-size warning—96 shipped action lines—with all action steps, output contracts, permission boundaries, and incident-specific warnings preserved. Added a focused probe that keeps shell/prose audit coverage complete and prevents the safe deletions from returning.

## Qualification

Passed — mechanical qualification and scope drift verified all 9 declared files and exact Implementation Summary agreement. The expected new-document warning was judged safe: `decisions/audits/2026-08-11-defensive-surface.md` is the requested evidence entry point and is consumed by the focused repository test, not a shipped runtime dependency. Manual qualification confirmed the change is substantive and subtractive, every keep/delete trace reaches its evidence, and no action data flow or permission boundary was removed.

## Testing

**Tests run:**
- `_dev/tests/defensive-surface-audit.sh`
- `_dev/tests/action-shell-blocks.sh`
- `bash -n _dev/tests/defensive-surface-audit.sh _dev/tests/contract-regressions.sh`
- `shellcheck --severity=warning _dev/tests/defensive-surface-audit.sh _dev/tests/contract-regressions.sh`
- `skills/do-work/tools/checks/qualify.sh do-work/working/REQ-168-defensive-code-delete-or-test-audit.md`
- `skills/do-work/tools/checks/scope-drift.sh do-work/working/REQ-168-defensive-code-delete-or-test-audit.md`
- `bash _dev/tests/contract-regressions.sh`
- `go test ./...` and `go vet ./...` in `skills/do-work-board/tools/queue-kanban`
- `git diff --check`
- explicit deleted-heading/phrase sweep

**Result:** ✓ All passing; aggregate contracts exited 0, the audit probe covers all current explicit defensive headings and 15 shipped shell sources, all shell fences/sources lint, and board Go tests/vet pass.

**Red-green validation:**
- Before deletion, the focused probe attributed every planned generic table/row as a failure, including all four Rationalizations/Red Flags pairs, both duplicate Rationalizations sections, the arbitrary commit-size threshold, and three generic inspect warnings.
- After deletion, the same probe passes while still requiring a disposition for every surviving explicit prose section and shell source.

**New tests added:**
- `_dev/tests/defensive-surface-audit.sh` — completeness and deletion ratchet for the audit

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/contract-regressions.sh` — invokes the focused audit probe; no prior assertion changed

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-11T12:55:15Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — all shipped defensive layers are dispositioned, surviving executable fallbacks have incident/test evidence, safe redundant prose is deleted, behavior-changing candidates remain explicit, and focused/aggregate validation is green.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Auditing cohesive mechanisms instead of counting branches made expensive-but-earned transactions distinguishable from cheap decorative prose.
- Starting from the existing steps and verification checklists made safe deletion objective: the removed tables carried no unique permission, output, or execution contract.

**What didn't:**
- Keyword counts dramatically overstate defensive surface: changelogs, test descriptions, ordinary domain language, and one atomic rollback transaction all contain many “guard/failure” tokens without representing separate layers.

**Worth knowing:**
- A numeric commit-size warning is weaker than semantic atomicity; it can split a correct large change and allow an unrelated small one.
- Static reference coverage is valid evidence for prompt-only permission/schema sections, while runtime fallbacks need syntax reachability plus targeted behavioral or contract fixtures.

## Orientation

[MAP CHANGED] Defensive-surface ownership is recorded in `decisions/audits/2026-08-11-defensive-surface.md`; `_dev/tests/defensive-surface-audit.sh` keeps the inventory complete and prevents six classes of generic warning apparatus from returning.
