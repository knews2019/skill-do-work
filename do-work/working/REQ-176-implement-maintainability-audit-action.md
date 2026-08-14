---
id: REQ-176
title: Implement the maintainability-audit action in do-work-toolbox
status: claimed
created_at: 2026-08-13T22:35:10Z
claimed_at: 2026-08-14T09:51:18Z
user_request: UR-040
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
related: [REQ-177, REQ-178]
depends_on: [REQ-178]
batch: maintainability-audit
write_set: [skills/do-work-toolbox/actions/maintainability-audit.md, skills/do-work-toolbox/actions/maintainability-audit-reference.md, skills/do-work-toolbox/SKILL.md, skills/do-work-toolbox/actions/help.md]
maintenance: false
---

# Implement the Maintainability-Audit Action

## What

Author a new do-work-toolbox action, `maintainability-audit`, from the validated draft spec in UR-040: a grounded, interactive, read-only codebase maintainability audit (measured metrics with user-calibrated bands → hotspot-scoped judgment → root-cause classes → persistent report with cross-run deltas), whose findings feed `do-work-toolbox validate-feedback`. Ship it as an action + reference companion, routed in SKILL.md and listed in help.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
