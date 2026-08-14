---
id: REQ-176
title: Implement the maintainability-audit action in do-work-toolbox
status: completed
created_at: 2026-08-13T22:35:10Z
claimed_at: 2026-08-14T09:51:18Z
completed_at: 2026-08-14T10:21:01Z
commit: f845e1c
kb_status: pending
user_request: UR-040
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
related: [REQ-177, REQ-178]
depends_on: [REQ-178]
batch: maintainability-audit
write_set: [skills/do-work-toolbox/actions/maintainability-audit.md, skills/do-work-toolbox/actions/maintainability-audit-reference.md, skills/do-work-toolbox/SKILL.md, skills/do-work-toolbox/actions/help.md, skills/do-work-toolbox/actions/code-review.md, skills/do-work/actions/help.md, _dev/tests/staged-skills-contract.sh]
maintenance: false
---

# Implement the Maintainability-Audit Action

## What

Author a new do-work-toolbox action, `maintainability-audit`, from the validated draft spec in UR-040: a grounded, interactive, read-only codebase maintainability audit (measured metrics with user-calibrated bands → hotspot-scoped judgment → root-cause classes → persistent report with cross-run deltas), whose findings feed `do-work-toolbox validate-feedback`. Ship it as an action + reference companion, routed in SKILL.md and listed in help.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read general.md, coding-guardrails.md, both primes, the REQ plan, UR-040 verbatim spec, audit-metrics CLI contract (main.go file-top comment + prime), and precedents (deep-explore split, validate-feedback footer/guardrail phrasing, capture.md accelerator block). Approach: (1) write the reference first (bands table w/ measured-by, calibration procedure, EXCLUDE defaults, command reference, manual fallbacks, lock-in-limit guidance, finding-class template, report template), (2) write the action (steps 1-8 per the plan skeleton, earned sections with do-work nouns), (3) routing edits in the four listed files, (4) add `maintainability-audit` to staged-skills-contract.sh toolbox_actions, (5) run the four suites + the ratchet/validator-name greps. Verified suite mechanics up front: route-count contract counts routed `./actions/*.md` paths against toolbox_actions; retired-trigger scanner bans `do-work <legacy-trigger>` text on live surfaces (so the validator is always written `do-work-toolbox validate-feedback`); action-shell-blocks lints every bash fence with `<placeholder>` neutralization matching capture.md's shape; Common Rationalizations noun check applies to the new action.
- [x] **[APPLY]:** Written exactly as planned: two new action files (action 2,446 w; reference 2,378 w), four routing/help edits, one test-array edit. Scope held to the 7 declared files + this REQ.
- [x] **[UNIFY]:** `git diff --stat`: 6 modified + 2 new files, 11 insertions / 8 deletions across the modified set — every path is in the declared write set (plus this REQ file). Suites: `action-shell-blocks.sh` exit 0 (67 fenced blocks; shellcheck unavailable in this sandbox — bash -n only, matching pre-flight); `shipped-package-reference-contract.sh` exit 0; `staged-skills-contract.sh` exit 1 and `contract-regressions.sh` exit 2 — in both, the ONLY failing lines are the recorded environmental probe (`run-blocked-check process-tree case left the descendant alive`, prescribed-shell-scripts-behavior.sh:132) plus its wrapper relay line; isolated re-run of `prescribed-shell-scripts-behavior.sh` confirms the probe is the sole root cause. No NEW failures. Per-file verification: `maintainability-audit.md` — steps 1-8 match the plan skeleton, crew loads at the right steps, blockquote is trigger-style with package justification, no "Loop usage" section, Common Rationalizations rows carry do-work nouns (queue, REQ, write_set, validator), both shell fences mirror capture.md's placeholder shape and pass lint; `maintainability-audit-reference.md` — bands table has measured-by column with the jscpd no-fallback statement, calibration/EXCLUDE/command-reference/fallback/lock-in/finding-template/report-template sections all present, manual churn block is shallow-check-first with the rename caveat and write_set cross-check; `SKILL.md` — code-review row dropped `audit codebase`, new row routes `./actions/maintainability-audit.md` exactly once, argument-hint gains the verb after code-review; toolbox `help.md` — entry after code-review, description column aligned at col 34; `code-review.md` — only the one Use-when phrase changed; core `help.md` — toolbox roster reflowed with the new verb, no other lines touched; `staged-skills-contract.sh` — one array element added matching format. Greps: `ratchet` → zero hits in both new files; bare `do-work validate-feedback` → 0 in both; `do-work-toolbox validate-feedback` → 4 each. No debug artifacts in the diff.

## Why (if provided)

No existing action covers this: code-review is qualitative and single-run, quick-wins eyeballs the same impact formula without measurement, and nothing persists reports with deltas. The validated spec's loop (audit → validate-feedback → capture → work → re-audit) closes a gap in the suite.

## Context

The draft spec is preserved verbatim in `do-work/user-requests/UR-040/input.md`. The full validation triage — what's confirmed, all corrections and seams with evidence — is `ai-reports/2026-08-13_2200_maintainability-audit-spec-validation/index.html` (committed on this branch). The builder starts from the verbatim spec and applies every requirement below; where they conflict, the requirements below win (they encode the validation and the user's decisions).

## Detailed Requirements

**Decisions already made (encode, don't reopen):**
1. Reports and waivers live in **`do-work/audits/`** (`audit-YYYY-MM-DD.md`, `waivers.md`) — not root `audits/`. The action creates the directory on first use; the EXCLUDE defaults must exclude the audit's own output from its metrics.
2. Prime context coverage is a **judgment finding only** — REQ-touch counts as impact evidence, no machine-enforced coverage assertion, no structured-scope change to prime files.
3. Terminology: the spec's "ratchet" becomes **lock-in limit** everywhere — a single number or zero-hit assertion pinned at the current worst observed value, tightening as fixes land. Aligns with the glossary term "lock-in test".
4. Lock-in limits are **proposals only**: the report proposes each limit with its red case; accepted ones flow through validate-feedback → capture and land as lock-in tests in the project's test suite (`_dev/tests/` here) or CI where one exists. The audit never installs enforcement itself.
5. Packaging: action + reference companion in do-work-toolbox (the established pair pattern; nearest analog deep-explore). The SKILL.md trigger phrase `audit codebase` **moves from code-review to this action**; code-review keeps `code-review` and `review codebase`.

**Corrections to the draft (validated against commit b1fad27):**
6. The validator command is `do-work-toolbox validate-feedback` (not `do-work validate-feedback`) — fix everywhere the spec references it.
7. Churn method must survive real repos: check `git rev-parse --is-shallow-repository` first (deepen/unshallow, or mark churn NOT-MEASURED); exclude release-ceremony files named by the project's commit ritual (here CHANGELOG*, VERSION, version.md — measured at 24% of all file-touches); normalize renames (`--follow` or old→new mapping — the 2026-08-08 skills/ restructure splits counts across dead paths). Cross-check hotspots against the do-work record: archived REQ `write_set:` frontmatter is rename-immune, shallow-proof churn.
8. do-work-record vocabulary: the Lessons heading is `## Lessons Learned`; the greppable scope fields are frontmatter `write_set:` and body `## Scope` / `**Files I will touch:**` ("REQ Scope" does not exist); coverage is incomplete (legacy REQs lack both) and any REQ-touch count must state that caveat; drop "PRDs" — the URs are the requirements history; a missing `do-work/queue/` directory means an empty queue.
9. The naming dimension **points to** `crew-members/coding-guardrails.md` § 5 (canonical copy — do not restate it; the draft's restatement had already drifted and dropped the single-word-by-design exemption that protects `do-work run`). The one audit-specific addition kept, framed as additive: names should match the project's own vocabulary (URs, primes, REQ titles).
10. Surface-cost pre-classification uses § 2's actual elements: concrete incident or replay case, **the surface that will remain**, the cost judgment, and the regression test.
11. `effort_estimate` uses the canonical enum `trivial | normal` only (no `large`); drop `Score = impact / effort` division — rank by impact, effort as tie-break (or define an explicit numeric mapping in the action; pick one, state it).
12. Severity comes from impact alone: ≥4 → P1, 3 → P2, ≤2 → P3. `gate` stays a separate routing field with the review-flow tokens (`user-visible` | `rule-change` | `trivial`); never let "trivial" double as both gate and severity.
13. First-run branch: when `do-work/audits/` has no prior report, say "no baseline" in the header, skip deltas, proceed. Rephrase Rule 2 honestly: "writes only under `do-work/audits/`". Drop brew/apt from the default install plan (contradicts user-local-only); state that duplication has no wc/grep fallback so a declined jscpd means NOT-MEASURED.
14. Default EXCLUDE adds, for do-work repos: the do-work record, knowledge stores (`kb/`), report outputs (`ai-reports/`, `do-work/audits/`), changelog archives. All confirmed at the calibration gate.
15. `## Instances` coordinates are grep-pattern-first: `path — greppable pattern`, line number an optional annotation valid only at the audited SHA (recorded lesson: path:line tables rot within hours on an active repo).
16. Rename the interactive gate: **calibration gate**, not "Checkpoint" (`do-work/CHECKPOINT.md` is a live pipeline artifact). The person is the **user**, not "operator". Question budget: one bundled approval (installs + bands + scope as a single editable proposal) plus at most three domain questions.
17. `sweep_key` travels in finding blocks and captured REQ bodies as provenance only — capture's schema is unchanged (sweep markers stay review-flow-only).

**Shipping apparatus (per prime-action-files.md):**
18. Action file gets the description blockquote, a trigger-style description, When to Use (positioned against code-review, quick-wins, forensics, and prime's audit subcommand), an `$ARGUMENTS` contract, and crew loads: `crew-members/clear-questions.md` before the calibration gate, `crew-members/anti-slop.md` before writing the report, `crew-members/prompt-injection.md` before reading archived REQ/Lessons prose (the forensics Check 10 precedent).
19. The "Loop usage" section does not ship inside the action: loop steps become the report's output footer (validate-feedback handoff pattern); the narrative goes to REQ-177's docs guide.
20. SKILL.md routing + help entry updated; the router word-budget lock-in test (2,650 words, `_dev/tests/contract-regressions.sh`) must stay green.
21. Agent-compatibility floor: generalized language only; every external tool (scc, lizard, jscpd, shellcheck, gocyclo) has either a declared fallback or an honest NOT-MEASURED path — none of the "Always" trio can be assumed present.
22. **Mechanical measurement goes through the shipped tool:** Phase 0/1's inventory, distributions, band flags, churn, and hotspot join invoke REQ-178's `audit-metrics` tool (build-on-demand, queue-kanban pattern) and paste its output; the reference file keeps the manual command fallback for when `go` is absent. Scripts over LLM calls for anything deterministic — cheaper and more robust (user's explicit instruction).

## Constraints

- The action itself is read-only outside `do-work/audits/` — no fixes, no installs into the repo tree.
- Never estimate what a tool can measure; pasted command output is the number.
- The two adopted-at-capture defaults (requirements 4 and 5) may be revisited at review if the user objects — note them in the hand-back.

## Dependencies

Depends on REQ-178 (`audit-metrics` tool) — this action's measurement steps invoke it. REQ-177 (docs guide) builds on this action's final vocabulary and routing; it depends on this REQ.

## Builder Guidance

Certainty: Firm on requirements 1–17 (validated + user-decided); Exploratory on exact section ordering and the action/reference split point — follow prime-action-files.md's earned-sections discipline and keep the action lean. Channel YAGNI: the reference holds the bands table, calibration detail, and report template; the action holds the flow.

## Open Questions

- [~] Lock-in limit enforcement model → default adopted at capture (proposals → REQs → lock-in tests); user dismissed the prompt. Flag in the hand-back so review can overturn cheaply.
- [~] `audit codebase` trigger takeover from code-review → default adopted at capture; same review note.

## Red-Green Proof
**RED prompt/case:** `grep -n 'maintainability-audit' skills/do-work-toolbox/SKILL.md` returns nothing; `skills/do-work-toolbox/actions/maintainability-audit.md` does not exist.
**Why RED now:** The capability exists only as a validated draft in UR-040; no action, routing, or help entry ships it.
**GREEN when:** The action + reference files exist and encode requirements 1–22 — explicitly including requirement 22's tool integration (the measurement steps invoke `audit-metrics`, with the manual fallback in the reference); SKILL.md routes `maintainability-audit` (and `audit codebase`); help lists it; `bash _dev/tests/contract-regressions.sh` and `bash _dev/tests/shipped-package-reference-contract.sh` exit 0.
**Validation:** Inferred during capture (user pre-authorized capture; decisions recorded above).

## Full Context

See `do-work/user-requests/UR-040/input.md` for complete verbatim input and `ai-reports/2026-08-13_2200_maintainability-audit-spec-validation/index.html` for the validation evidence behind every requirement.

---
*Source: UR-040 — pasted maintainability-audit spec, validated via do-work-toolbox validate-feedback*

---

## Triage

**Route: C** - Complex

**Reasoning:** Multiple coupled components (action + reference companion + SKILL.md routing + help entry), a 22-requirement spec, and a real architectural decision (the action/reference split point under the router word budget). The upstream tool (REQ-178) landed with a concrete CLI contract the prescribed shell blocks must match exactly.

**Planning:** Required

## Plan

Full plan by Plan agent (Route C), validated: 22/22 requirements traced, no orphan tasks, effort ≈3 core tasks (at threshold — flagged, kept as one coherent deliverable).

**File split:** action `maintainability-audit.md` (~2,300-2,500 w — flow: blockquote w/ package justification, write-boundary banner, Philosophy, When to Use positioned vs code-review/quick-wins/forensics/prime-audit, $ARGUMENTS contract, Steps 1-8, Output Format, earned Rules/Common Rationalizations/Red Flags/Verification Checklist) + reference `maintainability-audit-reference.md` (~2,200-2,600 w — default bands table with measured-by column, calibration procedure incl. FLAG=max(floor,p95) + bundled-proposal gate shape, do-work EXCLUDE defaults, audit-metrics command reference with real flags, manual fallbacks incl. shallow-check-first churn with rename caveat, lock-in-limit guidance w/ dimension-5 exception, finding-class template, full report template). Deep-explore split precedent.

**Steps skeleton:** 1 baseline/first-run branch → 2 grounding (prompt-injection crew load before record reads; accelerator block, no band flags) → 3 calibration gate (clear-questions load; bundled proposal + ≤3 questions) → 4 metrics (agreed bands; external tools or NOT-MEASURED; write_set cross-check) → 5 judgment (5 dims; naming points to §5 + additive vocabulary rule; dim 5 judgment-only) → 6 consolidate (sweep_key provenance; grep-pattern Instances; severity from impact alone ≥4/3/≤2; DECIDED: rank by impact desc, effort tie-break, no Score division; surface-cost w/ §2's four elements; lock-in-limit proposals w/ red case) → 7 report (anti-slop load; mkdir -p do-work/audits) → 8 debrief + loop footer (no Loop-usage section).

**Router/help:** code-review row drops `audit codebase`; new row `maintainability-audit`, `audit codebase`, `audit maintainability`; argument-hint insert; toolbox help.md entry; code-review.md:13 Use-when phrase edit; skills/do-work/actions/help.md toolbox roster line.

**Contract-suite awareness (verified by plan agent):** `_dev/tests/staged-skills-contract.sh` `toolbox_actions` array (line ~127) MUST gain `maintainability-audit` (route-count contract); shipped-package-reference auto-discovers (use backtick spans, no link to REQ-177's not-yet-existing guide); action-shell-blocks lints every fence (mirror capture.md's placeholder shape); Common Rationalizations noun check applies (not grandfathered); retired-triggers fixture stays green untouched (stale owner attribution noted for hand-back only).

**Prescribed block (Step 2 accelerator, capture.md:78-84 pattern):** subshell build + absolute-path run of `inventory`/`folders`/`churn` with do-work/kb/ai-reports excludes, NO band flags at grounding; Step 4 adds agreed `--watch-*`/`--flag-*` + ceremony excludes + `--since-window`. Fallback wording: "If `go` is absent or the build fails, fall back to the manual commands in the reference — accelerator, never a dependency."

*Generated by Plan agent*

## Exploration

Folded into planning: the Plan agent verified every target against the live tree (audit-metrics CLI flags from main.go, both contract suites' exact assertions with line numbers, router word counts, help formats, deep-explore split sizes, action-shell-blocks lint scope). No separate exploration pass needed — recorded as D-01 in Decisions.

*Generated by Plan agent (verification-grade)*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/maintainability-audit.md` (new) — the action
- `skills/do-work-toolbox/actions/maintainability-audit-reference.md` (new) — the companion
- `skills/do-work-toolbox/SKILL.md` (modify) — routing rows + argument-hint
- `skills/do-work-toolbox/actions/help.md` (modify) — menu entry
- `skills/do-work-toolbox/actions/code-review.md` (modify) — one Use-when line (drop `audit codebase`)
- `skills/do-work/actions/help.md` (modify) — toolbox roster line
- `_dev/tests/staged-skills-contract.sh` (modify) — `toolbox_actions` array + one route line

**Files I will NOT touch:** the audit-metrics tool (REQ-178, done), capture/work pipeline files, `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` (stays green; stale owner attribution noted in hand-back), CHANGELOG/VERSION (integrator's).

**Acceptance criteria (restated from REQ):**
- [x] Requirements 1-22 all encoded (traceability table in Plan); measurement steps invoke audit-metrics with the manual fallback in the reference
- [x] SKILL.md routes `maintainability-audit` + `audit codebase`; help lists it; code-review no longer claims `audit codebase`
- [x] `bash _dev/tests/contract-regressions.sh`, `shipped-package-reference-contract.sh`, `staged-skills-contract.sh`, `action-shell-blocks.sh` — no NEW failures vs the recorded environmental baseline
- [x] "lock-in limit" terminology throughout; never "ratchet"
- [x] Two capture-adopted defaults flagged in hand-back as revisitable (see Open Questions — lock-in enforcement model and the `audit codebase` trigger takeover)

## Decisions

- **D-01** (pre-assigned): Exploration folded into planning — Plan agent verified targets against live tree.
- **D-02** — **FILE_WORDS band row added with no numeric default.** The tool measures per-file words (`--watch-words`/`--flag-words`) and the draft spec names word budgets for prose-heavy repos, but gives no number; inventing one would ship a fake default. The row says "no universal default — propose from the repo p95", which is exactly the calibration procedure. DECIDE & STATE — reversible, consistent with the gate's own rule.
- **D-03** — **Instances heading is `#### Instances`, not the draft's `## Instances`.** A literal `##` heading inside a finding block would terminate the pasteable `## Findings` section at extraction time, breaking the self-contained handoff to the validator (the section "may be seen alone"). Heading demoted below `### Finding N`; the grep-pattern-first coordinate rule (requirement 15) kept verbatim, and the reference states the deviation and why. DECIDE & STATE.
- **D-04** — **Lock-in-limit landing place phrased generically** ("the project's own test suite or CI"), not `_dev/tests/` — `_dev/` is export-ignored, so a shipped citation would dangle in every consumer install (same rule that bans citing the primes). Requirement 4's `_dev/tests/` mention describes this repo, not shippable text. DECIDE & STATE.
- **D-05** — **Ceremony excludes prescribed as a condition, not a list.** Step 4's block carries a `--exclude-path <ceremony-file-prefix>` placeholder and the reference states the condition ("files the project's commit ritual touches") with the measured 24% figure as evidence — per Closed Enumerations Go Stale and the tool's caller-owns-the-list contract. DECIDE & STATE.

## Discovered Tasks

- `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` line for `audit codebase` attributes the retired trigger to `code-review` (`canonical_action`); after this REQ the live owner is `maintainability-audit`. Confirmed stale — the fixture is a pinned historical inventory and stays green untouched, but the attribution no longer matches live routing.
- `skills/do-work-toolbox/docs/code-review-guide.md:80` shows `do-work-toolbox audit codebase` in code-review's alias example block; that trigger now routes to maintainability-audit. Outside this REQ's write set — natural home is REQ-177's docs work.

## Pre-Flight

**Git:** ✓ clean outside do-work/ (claim bookkeeping only)
**Tests baseline:** ⚠ Environmental, pre-existing on surfaces byte-identical to origin/main: `contract-regressions.sh` exit 2 and `staged-skills-contract.sh` exit 1, both solely on the run-blocked-check process-tree probe (`prescribed-shell-scripts-behavior.sh:132` — sandbox process semantics). `shipped-package-reference-contract.sh` exit 0, `action-shell-blocks.sh` exit 0. Gate for this REQ: no NEW failing lines beyond that probe.
**Dependencies:** ✓ Go 1.24.7 present (audit-metrics builds; toolchain auto-fetches 1.26)

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/maintainability-audit.md` (new, 2,446 words)
- `skills/do-work-toolbox/actions/maintainability-audit-reference.md` (new, 2,378 words)
- `skills/do-work-toolbox/SKILL.md` (modified — code-review row drops `audit codebase`; new maintainability-audit row; argument-hint)
- `skills/do-work-toolbox/actions/help.md` (modified — menu entry)
- `skills/do-work-toolbox/actions/code-review.md` (modified — Use-when phrase edit)
- `skills/do-work/actions/help.md` (modified — toolbox roster line)
- `_dev/tests/staged-skills-contract.sh` (modified — `maintainability-audit` in `toolbox_actions` route contract)

**What was done:** Authored the maintainability-audit action + reference companion encoding all 22 validated requirements (grounding → calibration gate → measured metrics via audit-metrics with manual fallbacks → hotspot-scoped judgment → root-cause classes → persistent do-work/audits/ report with deltas → loop footer into do-work-toolbox validate-feedback), routed it in the toolbox with the `audit codebase` trigger takeover, and extended the route-count contract.

## Testing

**Tests run:** all four contract suites, orchestrator re-run
**Result:** ✓ No new failures — `action-shell-blocks.sh` exit 0, `shipped-package-reference-contract.sh` exit 0; `staged-skills-contract.sh` exit 1 and `contract-regressions.sh` exit 2 solely on the pre-existing environmental process-tree probe recorded in Pre-Flight (filtered FAIL diff vs baseline: empty). Route-count contract green with the new `toolbox_actions` entry.

**Red-green validation:** non-behavioral (instruction authoring, `tdd: false`) — regression evidence used instead: REQ's GREEN condition verified (routing grep hits SKILL.md:18, both help entries present, `audit codebase` absent from code-review.md, zero "ratchet" hits, validator name toolbox-prefixed 4×/file).

*Verified by work action*

## Review

**Overall: 95%** | 2026-08-14

| Dimension | Score |
|-----------|-------|
| Requirements | 98% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):**
- `docs/code-review-guide.md:80` still shows `do-work-toolbox audit codebase` as a code-review invocation after the trigger takeover (restatement sweep; outside this REQ's write set) — gate: user-visible → folded into queued REQ-177 (docs guide), whose write_set now includes the file; builder pre-recorded it in Discovered Tasks.

**Minor findings:** 3 — (1) Step 4 prescribed block lacked `--watch-words`/`--flag-words` placeholders (fixed by the integrator pre-commit; shell-block lint re-run green); (2) stale fixture attribution in `retired-core-moved-command-triggers.tsv:40` — pinned historical inventory, stays green untouched; (3) crew-citation form inconsistency (bare vs `../do-work/` prefixed) — cosmetic, both forms have sibling precedent.

**Nit:** the reference's command table restates the tool's caller contract with citation to `prime-audit-metrics.md` — acceptable (consumer installs need the caller-facing copy) but a drift-watch pair for future tool changes.

**Requirements: 22/22 delivered** (reviewer walked each against file content; §2 surface-cost elements verbatim-equivalent incl. "the surface that will remain"; severity-from-impact mapping in both files; zero "ratchet" hits; validator toolbox-prefixed 4×/file; template section order exact per prime-action-files).

**Acceptance:** Pass — all four suites reproduce the Pre-Flight environmental baseline exactly (only the run-blocked-check process-tree probe fails); reviewer built the tool and executed both prescribed command shapes against this repo (1,165 commits scanned; ceremony-exclude condition empirically confirmed incl. the mirrored-copy clause); routing takeover verified by grep; word counts match the Implementation Summary exactly.

**Scope drift:** None — diff is exactly the 7 declared files + the REQ file (hand comparison; scope-drift.sh gap is REQ-179).

**Restatement sweep:** trigger-ownership grep complete; the builder's Discovered Tasks list is the complete stale set (docs guide line 80 → REQ-177; fixture line 40 → untouched by design).

**Suggested testing:** one interactive end-to-end run (calibration gate + repeat-run branch); confirm FILE_WORDS bands appear when agreed; re-grep `audit codebase` in docs/ after REQ-177.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Plan-first with the traceability table — all 22 requirements landed on the first build pass and the reviewer confirmed 22/22 with zero remediation. Having the Plan agent pre-verify the contract suites' exact assertions (route-count array, noun checks, link parser) meant no suite surprises after authoring.
**What didn't:** The capture-seeded write_set missed three files the plan surfaced (`code-review.md`, core `help.md`, `staged-skills-contract.sh`) — routing takeovers always touch the OLD owner's Use-when text and the route-count contract, not just the router. Prescribed blocks with per-metric band flags are easy to leave incomplete (words flags omitted where lines flags were present) — bands-only-from-flags means an omitted placeholder silently loses a whole metric's bands.
**Worth knowing:** The `## Instances` heading in finding templates must ship demoted (`#### Instances`) or it terminates the pasteable `## Findings` section (D-03). The environmental process-tree probe failure (sandbox-only) is the recorded baseline for suite runs in this session — surfaces byte-identical to origin/main.

## Orientation

Now you can run a measured, calibrated, delta-tracked maintainability audit: `do-work-toolbox maintainability-audit` (also `audit codebase`, which moved here from code-review). [MAP CHANGED] — new action + reference pair in the toolbox's findings family (code-review/quick-wins produce, validate-feedback receives, the audit measures-then-judges); reports persist in `do-work/audits/`; measurement runs through the REQ-178 tool. Primes checked: `prime-action-files.md` and `prime-shell-commands.md` still accurate; `prime-audit-metrics.md` gains a caller (noted as drift-watch pair).
