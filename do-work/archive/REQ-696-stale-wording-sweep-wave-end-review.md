---
id: REQ-696
title: '[impact-rule-change] Sweep stale wording left by the wave-end review: stuck routing, fence rule, addendum example, sibling-route test text'
status: completed
created_at: 2026-10-10T19:19:53Z
user_request: UR-154
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-693, REQ-694, REQ-695]
batch: review-followups-ur154
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work/actions/forensics.md", "README.md", "skills/do-work/actions/clarify.md", "skills/do-work/actions/capture.md", "_dev/tests/fixtures/retired-core-moved-command-triggers.tsv", "_dev/tests/staged-skills-contract.sh"]
claimed_at: 2026-10-10T19:22:01Z
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-10-10T19:23:55Z
  basis:
    - trivial short-circuit
builder_handback_at: 2026-10-10T19:27:28Z
integration_at: 2026-10-10T19:30:16Z
review_at: 2026-10-10T19:39:31Z
kb_status: pending
commit: 92a2dca959b9561ff39362b0a861837deb0c6173
heavy_verified_at: 2026-10-10T19:39:51Z
heavy_verified_revision: 92a2dca959b9561ff39362b0a861837deb0c6173
completed_at: 2026-10-10T19:40:15Z
release_at: 2026-10-10T19:40:15Z
---
# Sweep Stale Wording Left by the Wave-End Review

## What

Fix five stale wording sites the last run's reviews recorded as report only. Each one restates a rule that changed in the same run, so two shipped texts now disagree.

## Detailed Requirements

- REQ-690 (run-status action) F9: `skills/do-work/actions/forensics.md:10` ("User suspects something is stuck, broken, or producing confusing results") and `README.md:176` ("Run `do-work forensics` to diagnose stuck or failed work.") still send "stuck" questions to forensics. `actions/status.md:10` now owns "is it stuck". Point stuck questions at `do-work status` and keep forensics for broken or failed work.
- REQ-688 (capture-files example and fence fix) F2: `skills/do-work/actions/clarify.md:106` says "open a code fence longer than the longest backtick run anywhere in the text". The Go writers and `capture-reference.md:196` apply the stricter rule. Reviewer's replacement: "open a code fence one backtick longer than the longest backtick run anywhere in the text, never shorter than three, with no info string."
- REQ-688 F3: the queued-addendum example at `skills/do-work/actions/capture.md:136` and `:138` uses a four-backtick `text` fence. Replace both lines with a bare three-backtick fence (`> ` followed by three backticks).
- REQ-692 (validate-feedback capture run) M2: the header on line 1 of `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` and the fail message at `_dev/tests/staged-skills-contract.sh:792` still state the old "core routes no sibling action" rule. Reviewer's fixes: append to the tsv header " The one exception is the core forward row for validate-feedback / triage feedback (UR-153, 2026-10-10)."; change the message to `fail "core must not route sibling-owned action $public_action through ./actions/ (a ../$sibling_owner/ forward row is allowed)"`.

## Red-Green Proof

**RED prompt/case:** `grep -n "stuck" skills/do-work/actions/forensics.md README.md`, `grep -n "longer than the longest backtick run" skills/do-work/actions/clarify.md`, and `grep -n "text$" skills/do-work/actions/capture.md` around line 136.
**Why RED now:** each grep shows the stale wording at the lines above.
**GREEN when:** forensics and README no longer route "stuck" to forensics, clarify states the three-backtick minimum and no info string, the capture example uses a bare three-backtick fence, and the two test texts name the forward-row exception. `bash _dev/tests/staged-skills-contract.sh` still passes.
**Validation:** Inferred during capture

## Constraints

- Prose only. No new test: the only check that covers these sites is the existing staged-skills contract, which must still pass.
- Re-verify each site at build time; tick one off with evidence if an unrelated commit already fixed it.
- No `depends_on` edges in this batch; no other REQ touches these files.

## Builder Guidance

Certainty: high (all five re-checked at HEAD `0e8e0ef9`; the line numbers in the input still match). Use the reviewers' exact text where given.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 8128 tokens, bare only (`slugged: partial`); matches the action-file edits (family `restated-mechanism-unchecked`).

## Full Context

See `do-work/user-requests/UR-154/input.md` for complete verbatim input. Sources: `do-work/runs/work-2026-10-10-131527/REQ-690-review.md` (F9), `do-work/runs/work-2026-10-10-131527/REQ-688-review.md` (F2, F3), `do-work/runs/work-2026-10-10-131527/REQ-692-review.md` (M2), the archived REQ-688, REQ-690 and REQ-692 under `do-work/archive/UR-153/`, and `do-work/runs/work-2026-10-10-131527/REQ-689-integration-report.md`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the brief, the REQ (Triage D-01 to D-04), the three review excerpts, `shared-principles.md`, `prime-releases.md`, `prime-action-files.md` (Cross-Referencing: same-package citation `actions/status.md` is the correct form from `actions/forensics.md`), `lessons-releases.md`. Approach: six exact-string replacements with each old string asserted to occur once, then probe, `bash -n`, `git diff --check`. (from the builder hand-back)
- [x] **[APPLY]:** Applied the six replacements with the reviewers' exact text via one Python script (each `count(old) == 1` asserted). No other edits. (from the builder hand-back)
- [x] **[UNIFY]:** `git diff fd9a2378 --stat`: 6 files changed, 7 insertions(+), 7 deletions(-) (README.md 2, tsv 2, staged-skills-contract.sh 2, capture.md 4, clarify.md 2, forensics.md 2). `bash do-work/runs/work-2026-10-10-192201/REQ-696-probe.sh` from the worktree root: exit 0, `REQ-696 GREEN probe: ok`. `bash -n _dev/tests/staged-skills-contract.sh`: exit 0. `git diff fd9a2378 --check`: exit 0. Each changed line read once in the diff: no stray whitespace; capture.md fences are exactly three backticks with no info string; tsv line 3 keeps its tabs. Not run (per brief and D-04): `staged-skills-contract.sh` itself (heavy-only lane, drained by the integrator) and the repository gate. (from the builder hand-back)
*Source: UR-154 R5 — "stale wording the wave-end sweep left: `forensics.md:10` and `README.md:176` still send \"stuck\" questions to forensics (REQ-690 F9); `clarify.md:106` ...; the addendum example at `capture.md:136,138` ...; the retired-trigger fixture header and the message at `_dev/tests/staged-skills-contract.sh:792` ... (REQ-692 M2)."*

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names all six files, the exact lines, and the reviewers' exact replacement text. Prose plus one test message and one fixture comment, no code path, so nothing needs discovery.

**Planning:** Not required

### Open Questions

No `## Open Questions` section exists. Choices made unattended at triage:

- D-01: forensics.md:10 and README.md:176 use the REQ-690 review's exact F9 fix text (`do-work/runs/work-2026-10-10-131527/REQ-690-review.md:85`), because the REQ says to use the reviewers' exact text where given and F9 gives one. DECIDE & STATE.
- D-02: forensics.md:5 ("feels broken, stuck, or produces confusing results") and `skills/do-work/docs/forensics-guide.md:3` ("detects stuck work") stay. They describe what forensics detects (stuck REQs are one of its checks, forensics.md:49), not where a "stuck" question routes, and the REQ-690 restatement sweep named only :10 and README:176. DECIDE & STATE.
- D-03: capture.md:136 and :138 both become `> ` plus three backticks. The enclosing example fence is a three-backtick markdown fence, but a line that starts with `> ` cannot close it (a closing fence may start only with up to three spaces), so the example still renders whole. DECIDE & STATE.
- D-04: The builder does not run `_dev/tests/staged-skills-contract.sh`. It is a heavy-only lane (it refuses without `DO_WORK_MAINTAINER_TIER=heavy`), and the integrator's heavy drain runs the `staged-skills` lane from `_dev/tests/heavy-lanes.json`, which covers `skills/` and `_dev/tests/`. The builder runs `bash -n` on it and the focused probe `do-work/runs/work-2026-10-10-192201/REQ-696-probe.sh`. DECIDE & STATE.
<!-- D-XX counter: last used D-04. Next decision: D-05. -->

## Plan

**Planning not required** - Route A: Direct implementation

Required-lessons consult (claim time): the captured set stays `_dev/primes/lessons-releases.md` (666 tokens, read; neither family applies, because this REQ moves no archived record and edits no manifest). Index re-check: `_dev/primes/lessons-action-files.md` still matches the action-file edits but is `slugged: partial` at 8128 tokens, so it cannot be narrowed and stays dropped (the `## Required Lessons — Dropped for Budget` section above is unchanged). No other index row matches six prose and test-message lines.

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/forensics.md` (modified)
- `README.md` (modified)
- `skills/do-work/actions/clarify.md` (modified)
- `skills/do-work/actions/capture.md` (modified)
- `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)

**What was done:** Replaced six stale wording sites with the reviewers' exact text: forensics.md:10 and README.md:176 now send "how long has in-flight work been quiet" to `actions/status.md` / `do-work status` and keep forensics for failed or broken work; clarify.md:106 states the fence rule in full (one backtick longer, never shorter than three, no info string); capture.md:136 and :138 use a bare three-backtick quoted fence; the retired-trigger fixture header and the staged-skills contract fail message name the validate-feedback forward-row exception. Merged as `8636a267..92a2dca9` (builder commit `e8381dc0`).

## Decisions

(from the builder hand-back)

- D-05: Applied all six edits in one Python pass that asserts each old string occurs exactly once, instead of six separate Edit calls, so a drifted site would stop the run instead of editing the wrong place. All six sites were still stale at base; none had been fixed by an unrelated commit. DECIDE & STATE.
- D-06: One commit for all six sites (brief prefers one commit; the change is one coherent sweep). DECIDE & STATE.

## Discovered Tasks

(from the builder hand-back)

- `skills/do-work/actions/capture-reference.md:196` still restates the full fence rule after citing clarify.md Step 4. Now that clarify.md:106 states the complete rule, line 196 could cite it and drop the restatement (REQ-688 review F2 called this optional; the REQ does not name it). → report only
- `skills/do-work/docs/forensics-guide.md:3` ("detects stuck work") and `forensics.md:5` keep "stuck" as something forensics detects, per D-02. If a later reader finds them confusing next to the new routing, a status pointer could be added there. → report only

## Qualification

**Gate records:** `advance REQ-696 --diff-range 8636a267..92a2dca9` returned the `qualify` gate `satisfied`, provenance `merged_range`, no findings (no debug artifacts, P-A-U boxes ticked).

**Requirement trace against `git diff 8636a267..92a2dca9` (6 files, 7+/7-) and the merged files:**
- REQ-690 F9: `skills/do-work/actions/forensics.md:10` now reads "User suspects something is broken or producing confusing results (for how long in-flight work has been quiet, see `actions/status.md`)"; `README.md:176` now names `do-work status` for quiet in-flight work and `do-work forensics` for failed or broken work. `actions/status.md:10` lists "is it stuck" among its triggers, so the new pointer is true. forensics.md:5, :49, :159 keep "stuck" as something forensics detects (Triage D-02).
- REQ-688 F2: `skills/do-work/actions/clarify.md:106` carries the reviewer's exact text ("one backtick longer than the longest backtick run anywhere in the text, never shorter than three, with no info string").
- REQ-688 F3: `skills/do-work/actions/capture.md:136` and `:138` are `> ` plus a bare three-backtick fence. The `> ` prefix keeps them from closing the enclosing three-backtick `markdown` example fence (Triage D-03).
- REQ-692 M2: the fixture header (line 1) gains the forward-row exception sentence; `cat -vet` shows the tab-separated header row on line 3 unchanged. `_dev/tests/staged-skills-contract.sh:792` fail message names `./actions/` and the `../$sibling_owner/` forward row; `$sibling_owner` is set by the loop's `read` at line 776, and the check logic is unchanged.

**GREEN probe:** `bash do-work/runs/work-2026-10-10-192201/REQ-696-probe.sh` on main at `92a2dca9`: exit 0, `REQ-696 GREEN probe: ok`. The staged-skills contract run is owned by the heavy drain (`staged-skills` lane).

**Scope:** Route A has no `## Scope` section; the six touched files equal the `write_set` frontmatter exactly. No file outside it changed.

*Qualified by work action (integrator)*

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` on main at `92a2dca9` (1-minute load 2.26 before the run, no other gate running), recorded through `advance REQ-696 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-192201/REQ-696-probe.sh`; plus the GREEN probe `bash do-work/runs/work-2026-10-10-192201/REQ-696-probe.sh`.
**Result:** ✓ Repository gate exit 0 in 139 s ("Maintainer verification passed."); probe exit 0 (`REQ-696 GREEN probe: ok`); `advance` recorded the `green-gate` and `run-blocked-check` gates satisfied.

**Red-green validation:** non-behavioral prose change (tdd: false); the GREEN probe stands in as the regression check:
- `do-work/runs/work-2026-10-10-192201/REQ-696-probe.sh`: ✗ at base `fd9a2378` (exit 1, 11 FAIL lines, from the builder hand-back) → ✓ at merge `92a2dca9` (exit 0)

**New tests added:**
- None (Constraints: prose only, no new test).

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/staged-skills-contract.sh` (from REQ-692): fail message text only; the check logic is unchanged. `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv`: header comment only.

**Heavy verification plan:**
- Range: 8636a26756353215b2a6c4e8452d64f719494963..92a2dca959b9561ff39362b0a861837deb0c6173
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests; skills/do-work/actions/capture.md matched subtree skills; skills/do-work/actions/clarify.md matched subtree skills; skills/do-work/actions/forensics.md matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — README.md matched exact path README.md; _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests

*Verified by work action*

## Review

**Overall: 100%** | 2026-10-10T19:39:31Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1 `skills/do-work-toolbox/actions/tutorial.md:260-261` routes "Something seems stuck" to `do-work forensics`, while core SKILL.md:40 routes `stuck` to status. Fix: point it at `do-work status` — impact-negligible → report only; F2 `skills/do-work/docs/roadmap-guide.md:48` routes "broken or stuck" to `do-work forensics`. Fix: broken → forensics, stuck in-flight work → status — impact-negligible → report only; F3 (Nit, builder find 1) `capture-reference.md:196` restates the fence rule after citing clarify.md. It now agrees with clarify.md:106 and both Go writers, so it is a duplicate only and not stale — impact-negligible → report only; F4 (Nit, builder find 2) `forensics-guide.md:3` and `forensics.md:5` keep "stuck" as something forensics detects (doctor `STUCK-WORK`, forensics.md:49). This is true, D-02 stands, and no change is needed — impact-negligible → report only. Anti-bloat: 0 helpers, options, files or tests that the REQ did not name.
**Acceptance:** Pass — implementation and integration stages: the GREEN probe passes on main, `bash -n` and `git diff --check` are clean, and the repository gate exited 0 at 92a2dca9. The staged-skills heavy lane is left to the integrator's drain.
**Restatement sweep:** redefined the containment fence rule (clarify.md:106; consistent: capture-reference.md:196, publication_manifest.go:107, state_apply.go:1052) and the "stuck" routing wording (stale: tutorial.md:260-261 F1, roadmap-guide.md:48 F2; consistent detection-only wording: forensics.md:5,49, forensics-guide.md:3,13, roadmap.md:5,124,278, stray-check.md:16, stray-check-guide.md:5,59)
**Suggested testing:** 1 item
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Taking the reviewers' exact replacement text and asserting each old string occurs once (D-05) made a six-site sweep a one-pass edit with no drift risk.
**What didn't:** The source reviews' restatement sweeps were narrower than the rule they checked: the review of this REQ still found two more "stuck → forensics" lines (`skills/do-work-toolbox/actions/tutorial.md:260-261`, `skills/do-work/docs/roadmap-guide.md:48`), one of them in a sibling package.
**Worth knowing:** "Stuck" has two meanings in the shipped prose: a question about quiet in-flight work (routes to `do-work status`) and doctor's `STUCK-WORK` finding (reported by forensics). Only the first is a routing claim; a sweep must grep the sibling packages too, not just `skills/do-work/`.

## Orientation

Shipped prose and two test messages now agree with the previous run's rule changes: "is it stuck" routes to status, and the outside-text containment fence rule is stated in full (action files, `_dev/primes/prime-action-files.md` area; release per `_dev/primes/prime-releases.md`). Primes spot-checked: both still exist and name no moved path.

## Heavy Verification Plan

- Base: 8636a26756353215b2a6c4e8452d64f719494963
- Target: 92a2dca959b9561ff39362b0a861837deb0c6173
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests; skills/do-work/actions/capture.md matched subtree skills; skills/do-work/actions/clarify.md matched subtree skills; skills/do-work/actions/forensics.md matched subtree skills
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — README.md matched exact path README.md; _dev/tests/fixtures/retired-core-moved-command-triggers.tsv matched subtree _dev/tests; _dev/tests/staged-skills-contract.sh matched subtree _dev/tests

## Heavy Verification Result

- Target revision: 92a2dca959b9561ff39362b0a861837deb0c6173
- Execution revision: 92a2dca959b9561ff39362b0a861837deb0c6173 (detached drain checkout, `QUEUE_KANBAN_BROWSER` set to Google Chrome; no `HEAVY-RUN-LANE-SKIPPED`)
- queue-kanban-javascript: exit 0, executed, 9 s
- queue-kanban-browser: exit 0, executed, 89 s
- do-work-cli-integrations: exit 0, executed, 71 s
- staged-skills: exit 0, executed, 40 s
- updater: exit 0, executed, 66 s
- installer: exit 0, executed, 28 s

## Timing

Observed 2026-10-10T19:29:58Z to 2026-10-10T19:39:51Z: 9m 53s total, 14m 10s attributed across 4 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 8m 14s | 2 |
| review | 5m 38s | 1 |
| handback-merge | 18s | 1 |

Slowest stage: verification-gate / heavy drain, 5m 45s, outcome success.

Note: no builder-work event. The hand-back had already landed when the integrator started, so per fan-out-reference "Landed hand-back" the builder wait was not recorded (dispatch 2026-10-10T19:26:48Z, builder commit 2026-10-10T19:27:28Z).
