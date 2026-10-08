---
id: REQ-647
title: 'Clarify tells the user when a released REQ stays earmarked and clears it on request'
status: completed
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-08T15:54:07Z
route: A
created_at: 2026-10-08T15:50:17Z
user_request: UR-142
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-648, REQ-649]
batch: earmark-release-visibility
write_set: [skills/do-work/actions/clarify.md, skills/do-work/docs/work-guide.md]
claimed_at: 2026-10-08T15:53:15Z
dispatch_at: 2026-10-08T15:56:02Z
builder_handback_at: 2026-10-08T16:09:05Z
integration_at: 2026-10-08T16:09:22Z
review_at: 2026-10-08T16:14:24Z
kb_status: pending
commit: be1f0ea4e7dc36144728ea0c9aca64fed0bcc7e2
heavy_verified_at: 2026-10-08T16:14:49Z
heavy_verified_revision: be1f0ea4e7dc36144728ea0c9aca64fed0bcc7e2
completed_at: 2026-10-08T16:15:22Z
release_at: 2026-10-08T16:15:22Z
---
# Clarify Tells the User When a Released REQ Stays Earmarked and Clears It on Request
## What
`do-work clarify` releases a `blocked` REQ through the `unblock` transaction, which removes `blocked_by` and `blocked_at` and leaves `assigned_to` in place. `actions/clarify.md` Step 5.5 then tells the user "The REQ re-enters the queue for the next `do-work run`", which is false when the REQ is earmarked: the default run skips it and the board shows it under Pending → Earmarked. Make clarify say so at the moment of release, offer to clear the earmark (a hand edit of the field, the documented clear path), name every released-but-earmarked REQ in the report, and fix the false sentence.
## Why
Consumer incident, Oct 8 08:56 UTC on 0.305.79: two REQs that needed the operator carried both `status: blocked` and `assigned_to: 'user-interactive'`. The user released both through clarify ("yes, unblock"). They returned to `pending`, landed under Pending → Earmarked, and every default run skipped them without saying so. The user waited, then asked why REQs that needed their input sat under Pending. The skill handed back a REQ it would not run and did not say so. Triage (validate-feedback, 2026-10-08) accepted the prose fix (Option A) and pushed back on a typed `--clear-assigned-to` flag: the hand edit is already the documented clear path, so a CLI flag adds surface for one advisory field.
## Verified Facts (from triage)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go:604-621` (`TransitionUnblock`): sets `status: pending`, stamps `status_changed_at`, deletes `blocked_by` and `blocked_at`, appends the Blocked history line; never reads `assigned_to`. The only writer that deletes the field is an explicit claim (`state_apply.go:562-565`). Both stay as they are.
- `skills/do-work/actions/clarify.md:188` ("Yes → unblock") ends "The REQ re-enters the queue for the next `do-work run`." `grep -n assigned_to skills/do-work/actions/clarify.md` finds nothing.
- `skills/do-work/actions/work-reference.md:114` (`assigned_to` schema line): "an assignment persists until an explicit run or a hand-edit clears it". `work-reference.md:511` (exit-summary row 6): "To drop an assignment without running it, remove the field by hand." The hand edit is the documented clear path, not a free-form reproduction of a transaction.
- Clarify already hand-edits queued REQ files in the same flow (Step 4 item 2 reclaim lines, the Builder Was Right / Discarded fast paths).
- `skills/do-work/docs/work-guide.md:120` (Earmarking with `assigned_to`): names the two ways to take an earmarked REQ (run it by name, or delete the field by hand) and says operator-wait work is captured `blocked` and released by clarify; it does not say what happens when a released REQ is also earmarked.
- No `_dev/tests` lock-in pins the sentence at `clarify.md:188` (grep for "re-enters the queue" across `_dev/tests/` finds nothing).
## Detailed Requirements
1. `actions/clarify.md` Step 5.5, the "Yes → unblock" bullet: keep the transaction sentence and the "do not reproduce the mutation free-form" rule unchanged. After the transaction, add the condition: when the released REQ carries `assigned_to`, tell the user on the spot that the REQ is still earmarked for `<session name, verbatim>` and that the default run will skip it, and ask one question with two concrete options (`crew-members/clear-questions.md` is already loaded by Step 3): clear the earmark now (remove the `assigned_to` line by hand edit, the documented clear path in `actions/work-reference.md` → Request File Schema `assigned_to` and the exit summary's assigned-elsewhere row), or keep it and run the REQ by name with `do-work run REQ-NNN`. Replace the closing sentence with: "The REQ re-enters the queue for the next `do-work run` unless it carries `assigned_to`, in which case the default run skips it until it is named explicitly or the field is cleared."
2. Say how the hand edit is committed: the same way clarify's other hand edits in the same invocation are committed today (read the file; do not invent a new commit step or a new transaction).
3. `actions/clarify.md` Step 6 (Report): the summary lists every released REQ that stays earmarked, with the verbatim session name and the run-by-name command. One sentence keyed on the condition, not a new subsection.
4. `docs/work-guide.md` Earmarking paragraph: one sentence saying that when clarify releases a blocked REQ that is also earmarked, it says so and offers to clear the earmark. Keep the paragraph's existing sentences.
5. Verification Checklist in `clarify.md`: one line that every unblocked REQ still carrying `assigned_to` was named to the user and in the report. Only if the checklist's existing shape takes it without ceremony; otherwise skip and say so in the hand-back.
6. Run `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh`; both green.
7. Release per `_dev/primes/prime-releases.md` (the integrator writes the CHANGELOG entry and bump; the title says what shipped in plain words).
## Constraints
- Prose only. No Go change, no new CLI flag, no new transaction, no `verify` probe (both were pushed back at triage; see the UR). `TransitionUnblock` and the default scan stay as they are.
- `assigned_to` stays verbatim-read: the question quotes the session name as written, no normalization.
- Clearing stays an explicit act: never clear the field silently on unblock, and never clear it when the user picks "keep".
- Do not edit `actions/work-reference.md` (REQ-648 owns its `assigned_to` schema line; the exit-summary row 6 is unchanged) or `actions/capture.md`.
- Read `_dev/primes/prime-action-files.md` before editing; every cross-reference uses that prime's spelling (`actions/work-reference.md` → section name).
- Batch constraint: this REQ is its own release and its own commit; do not fold REQ-648 or REQ-649 into it.
## Dependencies
None. REQ-648 (the schema line and the work action state the operator rule) and REQ-649 (board tooltip) touch different files; the three may run in parallel.
## Builder Guidance
Certainty is high: the triage read the transaction and the action text, and the remedy is two conditions in one action file plus one sentence in the user guide. Latitude: exact wording of the question and the report line, and whether the checklist line is worth adding.
## Red-Green Proof
**RED prompt/case:** A queued REQ with `status: blocked`, `blocked_by: 'the operator'`, `blocked_at`, and `assigned_to: 'user-interactive'`. Run `do-work clarify` and answer "Yes, unblock it".
**Why RED now:** `clarify.md` Step 5.5 runs the unblock transaction and says the REQ re-enters the queue. Nothing reads `assigned_to`, so the user is not told the default run will skip it, is not offered the clear, and the Step 6 report does not name it.
**GREEN when:** A cold read of `clarify.md` Step 5.5 routes that case to a two-option question (clear the earmark by hand edit, or keep it and run by name), the closing sentence carries the `assigned_to` exception, and Step 6 names every released-but-earmarked REQ with its verbatim session name. `grep -n assigned_to skills/do-work/actions/clarify.md` finds the Step 5.5 condition and the Step 6 line. `shipped-package-reference-contract.sh` and `contract-regressions.sh` stay green.
**Validation:** User confirmed. The maintainer approved the triage's Option A remedy and said "capture and run".
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (6587 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its owning prime governs the action files this REQ edits; family `alternate-writer-contract-drift` (REQ-644: a rule added to one action stales its restatements elsewhere) is the trap this REQ's work-guide sentence guards against.
## Full Context
See `do-work/user-requests/UR-142/input.md` for complete verbatim input (the consumer incident timeline and the full triage). No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the brief, REQ-647, crew members (general, coding-guardrails, shared-principles, communication-style), prime-action-files, prime-releases, and every `alternate-writer-contract-drift` lesson. Read clarify.md in full to find its commit discipline: Step 4/5 answers go through `answer --manifest … --commit`; the Step 5.5 `unblock` invocation carries no `--commit`; the Step 5.5 reclaim line and the Builder Was Right / Discarded fast paths state no commit step at all. So requirement 2 takes the brief's fallback: say clarify states no commit step for hand edits and point at `actions/commit.md` in one clause. Edits: (1) extend the "Yes → unblock" bullet in place after the "do not reproduce" sentence, (2) one sentence appended to Step 6, (3) one checklist line beside the existing blocked-REQ line, (4) one sentence appended to the work-guide Earmarking paragraph. No new bold labels or headings; cross-references in the prime's `actions/work-reference.md` → **section** spelling.
- [x] **[APPLY]:** Applied exactly the four planned edits to the two write-set files. Option wording: "Clear the earmark now" (hand edit, documented clear path cited to Request File Schema `assigned_to` and the Composed Exit Summary assigned-elsewhere row) and "Keep it, and run the REQ by name with `do-work run REQ-NNN`". Added "Never clear the field without that answer, and leave it untouched on keep." The closing sentence is the REQ's text verbatim.
- [x] **[UNIFY]:** `git diff --stat`: `clarify.md | 9 +++++++--`, `work-guide.md | 2 +-`, 2 files changed, 8 insertions(+), 3 deletions(-). `git diff --check` clean. `bash _dev/tests/shipped-package-reference-contract.sh` exit 0 ("shipped package reference contract: PASS"). `bash _dev/tests/contract-regressions.sh` exit 0 ("Contract regression checks passed."). Files checked: clarify.md (Step 5.5 bullet reads cold, the transaction sentence and "do not reproduce the mutation free-form" are unchanged, Step 6 line, checklist line, all bold labels and headings unchanged); work-guide.md (paragraph intact, new sentence last). No debug artifacts.

*Source: consumer feedback "unblock keeps assigned_to, and clarify says nothing about it" (Gap 1 / the ask §1, Option A), accepted by `do-work-toolbox validate-feedback` on 2026-10-08 as F1 and F2; Option B and the verify probe pushed back.*

---

## Triage

**Route: A** - Simple

**Reasoning:** Names the two files (actions/clarify.md Step 5.5 and Step 6, docs/work-guide.md Earmarking paragraph) and quotes the sentences to add; prose only, no exploration needed.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/clarify.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)

**What was done:** Clarify's Step 5.5 "Yes → unblock" bullet now checks `assigned_to` after the unblock transaction, tells the user the REQ is still earmarked and that the default run skips it, and asks one two-option question (clear the earmark by hand edit, committed through `do-work commit`, or keep it and run the REQ by name); its closing sentence carries the `assigned_to` exception. Step 6 names every released REQ that stays earmarked, the Verification Checklist gains one matching line, and the work guide's Earmarking paragraph gains one closing sentence. Merged from builder commit fe2c1878 as be1f0ea4.

## Decisions

(from the builder hand-back)

- **D-01 (DECIDE & STATE) — commit discipline for the earmark hand edit (requirement 2).** Clarify has no stated commit discipline for hand edits: only Step 5's `answer` carries `--commit`; the Step 5.5 `unblock` invocation has no `--commit`, and the reclaim line and Builder Was Right / Discarded fast paths name no commit step. Per the brief's fallback, the new clause says so and routes the edit through `do-work commit` (`actions/commit.md`). Value: no invented commit step or transaction. Risk: low; the edit can sit uncommitted until the user commits, which is already true of the unblock itself.
- **D-02 (DECIDE & STATE) — checklist line (requirement 5).** Added. The checklist already has a per-condition line for blocked REQs ("`blocked` REQs the user confirmed satisfied flipped to `pending` …"), so one more line beside it fits the existing shape without ceremony. It also pins "cleared only on the user's answer", the never-silent constraint.
- **D-03 (DECIDE & STATE) — "keep" option names the claim side effect.** The keep option says `do-work run REQ-NNN` "claims it and clears the field", matching work-guide.md's existing sentence, so the user knows naming the REQ also drops the earmark.
- **D-04 (DECIDE & STATE) — placeholder spelling.** Used `<assigned_to, verbatim>` rather than `<session name, verbatim>` so it matches the exit-summary row 6 field cell ("assigned to <assigned_to, verbatim>") exactly.

## Discovered Tasks

(from the builder hand-back)

- impact-user-visible (re-stamped from the builder's non-canonical `impact-low`, review N3): `actions/clarify.md` Step 5.5 invokes `unblock` without `--commit` (the CLI accepts `--commit`, `internal/requeststate/state_commands.go`), while Step 5 commits its `answer` transaction; a clarify session that only unblocks leaves an uncommitted lifecycle change. → report only
- impact-rule-change (re-stamped from the builder's non-canonical `impact-low`, review N3): `actions/clarify.md` "Builder Was Right / Discarded" section describes hand edits (status flip, archive, appended notes) while Step 5 says the `answer` command derives the disposition and "There is no hand-edit, helper, or manual Git fallback"; the two sections may disagree about who writes the fast-path result. → report only

## Qualification

**Gate records:** `advance --diff-range ffc1a7bb..be1f0ea4` returned the `qualify` gate `satisfied` (provenance `merged_range`) with no findings: no debug artifacts, no output primitives, P-A-U boxes ticked from the hand-back.

**Diff:** `git diff ffc1a7bb..be1f0ea4 --stat` shows `skills/do-work/actions/clarify.md` (+7/-2) and `skills/do-work/docs/work-guide.md` (+1/-1). Read in full against the files at the merge.

**Requirement trace:**
1. Step 5.5 "Yes → unblock": the transaction sentence and "do not reproduce the mutation free-form" are byte-identical. The added condition tells the user the REQ is still earmarked for `<assigned_to, verbatim>` and that the default `do-work run` skips it, then asks one question with the two required options. The clear option cites `actions/work-reference.md` → **Request File Schema** `assigned_to` and → **Composed Exit Summary (Step 1)** assigned-elsewhere row; both targets say what the clause claims (schema line 114 "persists until an explicit run or a hand-edit clears it", row 6 "remove the field by hand"). The closing sentence is the REQ's text verbatim. The placeholder follows row 6's spelling, not the REQ's `<session name, verbatim>` (builder D-04, accepted).
2. Commit discipline: clarify states no commit step for its hand edits (the `unblock` call itself has no `--commit`), so the clause says so and routes the edit through `do-work commit`. No new step or transaction was invented (builder D-01, accepted).
3. Step 6 gains one condition-keyed sentence naming each released-but-earmarked REQ with the verbatim session name and `do-work run REQ-NNN`. Met.
4. Work guide Earmarking paragraph: every existing sentence kept; one closing sentence added. Met.
5. Verification Checklist: one line beside the existing blocked-REQ line, same shape. Met.
6. Contract scripts: green in the builder worktree; the repository gate re-runs them at the merge below.

**Never-silent constraint:** "Never clear the field without that answer, and leave it untouched on keep." The keep option names the claim side effect (builder D-03), which matches the existing work-guide sentence.

**Scope:** Route A has no `## Scope`; touched files equal `write_set` exactly. No Go, CLI flag, transaction, or probe added. `actions/work-reference.md` and `actions/capture.md` untouched.

**Note for review:** the schema citation uses the short heading form `**Request File Schema**`; the heading is "Request File Schema — Full Frontmatter". One other file uses the full form. The shipped-package reference contract passes, so this is wording, not a broken link.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at merge be1f0ea4 (the repository gate; it includes `contract-regressions.sh`), recorded through `advance --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0` with the focused probe `do-work/runs/work-2026-10-08-155307/REQ-647-probe.sh` (grep for `assigned_to` in `clarify.md`, then `shipped-package-reference-contract.sh`).
**Result:** ✓ Gate exit 0 on the first run (wall 117 s); probe exit 0; `test-gate`, `run-blocked-check`, and `green-gate` records satisfied.

**Red-green validation:** `tdd: false`, prose-only change, so no harness test. The `## Red-Green Proof` case is verified by reading the diff: at ffc1a7bb `grep -n assigned_to skills/do-work/actions/clarify.md` finds nothing and Step 5.5 says the released REQ re-enters the queue; at be1f0ea4 the grep finds the Step 5.5 condition (line 188), both options (189-190), the closing exception (192), the Step 6 line (198), and the checklist line (275). A cold read of the blocked + `assigned_to: 'user-interactive'` case released with "Yes" now reaches the two-option question.

**New tests added:** none (the REQ forbids a `verify` probe and names no lock-in test).

**Heavy verification plan:**
- Range: ffc1a7bbba79cefa73a208360b44bd92fa757527..be1f0ea4e7dc36144728ea0c9aca64fed0bcc7e2
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — `skills/do-work/actions/clarify.md` and `skills/do-work/docs/work-guide.md` matched subtree skills

*Verified by work action*

## Review

**Overall: 98%** | 2026-10-08T16:14:24Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** N1 `clarify.md:189` cites the schema by the short heading "Request File Schema", which is the house spelling at about 15 other sites; no change needed — impact-negligible → report only. N2 `clarify.md:188-190` two-option question names no recommended option although clearing is the evident answer for the incident shape — impact-negligible → report only. N3 builder Discovered Tasks carried non-canonical `impact-low`; re-stamped DT1 (unblock without `--commit`) impact-user-visible and DT2 (Builder Was Right hand edits vs Step 5 no-hand-edit wording) impact-rule-change, both noncritical — impact-negligible → report only
**Acceptance:** Pass — cold read of Step 5.5 routes the blocked + `assigned_to: 'user-interactive'` release to the two-option question and the closing exception; Step 6 names it.
**Restatement sweep:** redefined the Step 5.5 closing "re-enters the queue" sentence (now conditional on `assigned_to`); restatements read and consistent: `clarify.md:13`, `work-guide.md:102`, `work-guide.md:120`, `work-reference.md` exit-summary rows 3 and 6
**Suggested testing:** 0 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Reading the whole action for its commit discipline before writing the clear option. Clarify commits only Step 5's `answer`, so the honest clause says so and points at `do-work commit` instead of inventing a commit step.
**What didn't:** Nothing failed. The builder's Discovered Tasks used `impact-low`, which is not a Step 10 impact token; review re-stamped them.
**Worth knowing:** A transaction that clears part of a REQ's state can leave a sibling field that still steers the scheduler. `unblock` removes `blocked_by`/`blocked_at` but keeps `assigned_to`, so "re-enters the queue" was false for an earmarked REQ. When an action reports a state change, check every field the next reader filters on, not only the fields the transaction wrote.

## Orientation

Now clarify warns when a REQ it releases from `blocked` is still earmarked, offers to clear the earmark or run it by name, and names it in its report; lives in the clarify action (action-file prose, `prime-action-files`). No staleness in the touched primes.

## Heavy Verification Plan

- Base: ffc1a7bbba79cefa73a208360b44bd92fa757527
- Target: be1f0ea4e7dc36144728ea0c9aca64fed0bcc7e2
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — `skills/do-work/actions/clarify.md` and `skills/do-work/docs/work-guide.md` matched subtree skills

## Heavy Verification Result

- Target revision: be1f0ea4e7dc36144728ea0c9aca64fed0bcc7e2
- Execution revision: be1f0ea4e7dc36144728ea0c9aca64fed0bcc7e2 (detached checkout `.git/work-run-work-2026-10-08-155307/drain-head-REQ-647`, `QUEUE_KANBAN_BROWSER` set)
- staged-skills: exit 0, disposition executed, 43 s

## Timing

Observed 2026-10-08T15:56:02Z to 2026-10-08T16:12:09Z: 16m 07s total, 15m 21s attributed across 3 events, 46s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 13m 03s | 1 |
| verification-gate | 2m 08s | 1 |
| handback-merge | 10s | 1 |

Slowest stage: builder-work / builder worktree build, 13m 03s, outcome success.
