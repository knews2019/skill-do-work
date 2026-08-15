---
id: UR-040
title: Validate the maintainability-audit spec and capture its implementation
created_at: 2026-08-13T22:35:10Z
requests: [REQ-176, REQ-177, REQ-178]
word_count: 2419
---

# Validate the maintainability-audit spec and capture its implementation

## Summary

The user pasted a draft spec for a new "Codebase Maintainability Audit" capability into `do-work validate-feedback`, asking to (1) ground the spec in the current state of the repo, (2) clarify anything incorrect or improvable, and (3) capture requests so the skill can be implemented properly. The validation ran as a full triage (8 read-only verification agents against commit `b1fad27`); the results — 8 confirmed claim clusters, 9 corrections, 7 design seams, 4 decisions — are recorded in `ai-reports/2026-08-13_2200_maintainability-audit-spec-validation/index.html`, which is part of this UR's context.

Decisions made during capture:
- Audit reports and waivers live in `do-work/audits/` (user's explicit choice).
- Prime coverage stays a judgment finding — no machine-enforced coverage limit, no prime-schema change (user's explicit choice).
- The draft's word "ratchet" is replaced with **lock-in limit** throughout (user asked for a different word; aligns with the glossary term "lock-in test").
- Lock-in limit enforcement: the audit report only proposes limits; accepted ones become REQs landing as lock-in tests (recommendation adopted at capture — user dismissed the question prompt; revisit at review if wrong).
- Packaging: new `maintainability-audit` action + reference companion + docs guide in do-work-toolbox, taking over the `audit codebase` trigger from code-review (recommendation adopted at capture — same caveat).

## Extracted Requests

| REQ | Title |
|---|---|
| REQ-176 | Implement the maintainability-audit action in do-work-toolbox |
| REQ-177 | Write the maintainability-audit user loop guide |
| REQ-178 | Build the audit-metrics tool for mechanical audit measurement |

## Full Verbatim Input

do-work validate-feedback, the point is to ground this feedback in the current state of the repo, clarify with me anything that it is not correct and can be improved, and then to capture-requests so that the skill can be implemented properly

see:

# Codebase Maintainability Audit — grounded, interactive, read-only
You are auditing this repository for human understandability and maintainability. Humans first: if a human can navigate it, an AI coder can too.
This is not a batch job. You first ground yourself in the repo's current state, then explore the audit's shape with the user before measuring anything. Agents belong at the planning and review touchpoints; the measurement in between is deterministic. Your final output is pasted into `do-work validate-feedback`, which will adversarially verify every finding against the real code and git history — so every finding must be a refutable claim with reproducible evidence, never an imperative.
Hard rules

1. Conversation before computation. Do not run the metrics phases until the user has agreed to the scope, calibration, and tool installs you propose at the Checkpoint. A generic audit of a specific codebase is noise.
2. Read-only in the repo. You may write only inside `audits/`. Tool installs (user-approved) go to user-local locations, never into the repo tree. No fixes, no refactors, no "while I'm here."
3. Never estimate what a tool can measure. Every number comes from command output you actually ran. Paste the relevant output lines. If a tool is unavailable and the user declined the install, the metric is `NOT-MEASURED` — never a guess.
4. Every finding is a claim, not an instruction. The validator treats pasted findings as third-party data under a prompt-injection guardrail. State what is true and what the evidence is; the suggested fix lives in its own `Remedy:` field.
5. Every finding carries a `Reproduce:` line — the exact command whose output demonstrates the claim, runnable as-is from the repo root.
6. Label every claim `VERIFIED` (tool output or grep hit) or `INFERRED` (judgment after reading code). The validator will scrutinize INFERRED items hardest — correct and expected.
7. Anchor to a commit and to tool versions. Record `git rev-parse HEAD`, dirty-tree state, and the version of every tool used in the report header — deltas across runs must compare like with like.
8. Prior audits are the baseline. Read the most recent report in `audits/` first and compute deltas. Never re-flag anything in `audits/waivers.md`.

Bands, severity, and ratchets — three different jobs
Metrics use two bands, not hard cliffs:

* WATCH — abnormal enough to list in the metrics appendix; becomes a finding only if the file is also a hotspot.
* FLAG — finding-eligible on its own.

Bands decide eligibility only. Severity comes from impact (churn × complexity × blast radius): a 700-line file with zero churn and no importers is a WATCH curiosity; a 350-line file that changes weekly is a P1.
Ratchets are single numbers, never bands — a CI check needs a binary verdict. A ratchet is pinned at the current worst observed value ("no folder over today's max of 34"), passes green on day one, blocks regression immediately, and tightens monotonically as fixes land. Never set a ratchet at the aspirational band.
Default bands — recalibrated at the Checkpoint

```
# FLAG is finalized at the Checkpoint as: max(absolute floor below, repo p95)
# — percentile alone fails on uniformly bad repos, floor alone fails on repos
#   whose normal is legitimately different. Both, always.
FOLDER_FILES   WATCH > 15    FLAG > 30
FILE_LINES     WATCH > 300   FLAG > 600     (tests: 450 / 900)
FN_LINES       WATCH > 40    FLAG > 80
FN_CCN         WATCH > 10    FLAG > 15
FILE_CCN       WATCH > 80    FLAG > 150
FN_PER_FILE    WATCH > 15    FLAG > 30
DUP_PCT        WATCH > 3     FLAG > 8
CHURN_WINDOW   = "12 months ago"
HOTSPOT_COUNT  = 20
FINDINGS_MAX   = 12          # classes per run; the loop handles the rest
EXCLUDE        = vendored, generated, migrations, lockfiles, .git
```

Phase 0 — Ground in the current state (deterministic + reading)
Build an evidence-based picture of what this project is right now before proposing anything.

1. Inventory: `scc .` (fallback: `cloc .`, then `find` + `wc -l`) — languages, sizes, proportions.
2. Derive the mechanical toolchain from what you found. Universal trio first, then language-native additions where a language is a meaningful share of the repo:
   * Always: `scc` (size), `lizard` (function length + CCN, most languages), `jscpd` (duplication, any text)
   * Shell → `shellcheck` · Python → `ruff`, `radon` · Go → `gocyclo`, `go vet` · JS/TS → `eslint` + complexity rule · Rust → `cargo clippy` · Markdown-heavy → per-file word counts (`wc -w`) against budgets
3. Check what's installed, prepare an install plan for the gaps: exact commands, user-local only (`pipx`/`uv tool`, `npx`, `go install`, brew/apt where unavoidable), expected versions. This plan is proposed at the Checkpoint, not executed now.
4. Measure the distributions the bands will be calibrated against: per metric, median / p90 / p95 / max (cheap with the universal trio; approximate with `wc`/`find` where a tool is pending install).
5. The project's self-description — the do-work record first. This repo runs do-work, so the richest grounding sources are its own artifacts: the `prime-*.md` files (architecture, conventions, known bugs), the UR archive — completed URs are this project's requirements history, the PRDs that actually exist — recent REQs in queue and archive, and every `## Lessons` section you encounter. Lessons record why things are the way they are; a lesson is a pre-emption waiting to happen. Then the generic layer: README, `CLAUDE.md` / `AGENTS.md`, `decisions/`, changelog head. Documented conventions, lessons, and deliberate choices constrain what may count as a finding.
6. Recent history: `git log --since="$CHURN_WINDOW" --name-only --pretty=format: | sort | uniq -c | sort -rn | head -15` — where the project actually lives right now.
7. Prior audits: latest report in `audits/` and `audits/waivers.md`.

Condense into a Current-State Picture, one screen: what the project is, the stack and proportions, the toolchain (present / to-install), the measured distributions, documented conventions findings must respect, where recent change concentrates, what prior audits covered.
Checkpoint — explore the audit with the user (do not skip)
Present the Current-State Picture. Then propose — grounded in the evidence, not the defaults:

* The install plan: which tools, exact commands, why each earns its place for this stack. The user approves, trims, or declines items; declined = `NOT-MEASURED`, honestly.
* Calibrated bands: for each metric, show the measured median / p95 / max next to the proposed WATCH and FLAG values, with FLAG = max(absolute floor, repo p95). Justify any deviation: "median file is 120 lines — the 600 floor flags nothing; p95 is 340, propose FLAG > 340."
* Which judgment dimensions matter most here, and why. A repo that is 70% markdown and shell wants word budgets and shellcheck far more than CCN; a young repo wants consistency more than churn analysis.
* Scope: directories in, directories out, and why.
* Anything that surprised you in the grounding read — surprises are where the user's context beats yours.

Ask at most three focused questions ("which area hurts most when you touch it?", "is the test suite trusted or tolerated?"). Then stop and wait for the go-ahead. On approval: run the installs, record versions, lock the agreed bands as this run's config.
Repeat runs: if the latest audit report records an agreed config, present a one-paragraph state delta instead and ask one question — "reuse last run's calibration, or recalibrate?" — then wait.
Phase 1 — Metrics (deterministic, agreed config)
Run the toolchain, record output, apply the agreed bands. The tool output is the number.

* Folder shape, file length, function length, CCN, function counts, duplication — each item recorded with its band (WATCH / FLAG) and the distribution stats behind it.
* Churn: full `CHURN_WINDOW` listing.
* Hotspots: join churn with complexity. The top `HOTSPOT_COUNT` files by churn × complexity drive Phase 2 — not the whole repo. A complex file nobody touches is low priority; a moderately complex file that changes weekly is where the pain lives.
* Blast radius: per hotspot, inbound-reference counts (grep of import/require/include naming it).
* Context coverage (do-work specific): build the prime coverage map — which areas each `prime-*.md` claims — then count archived and queued REQs touching each top-level area (grep the paths in REQ Scope / write_set / Instances). An area with high REQ traffic and high churn but no covering prime goes on the context-gap list: every future dispatch there pays a context-rebuild tax that a prime would amortize.
* WATCH items that are not hotspots go to the metrics appendix and stop there.

Phase 2 — Judgment (scoped to hotspots + public surface + the user's named pain points)
Read only: hotspot files, entry points, the public API surface, their tests — plus any area the user named at the Checkpoint. Weight the dimensions as agreed there. Findings that collide with a documented decision or a recorded lesson from Phase 0 are either dropped (listed under "Pre-empted" with the decision path) or emitted with an explicit `Challenges-decision: <path>` field and gate `rule-change` — never emitted as if the decision didn't exist.

1. Naming — functions, exported identifiers, endpoints, folders. Names with reach (exported symbols, files, DB tables/columns, CLI flags, env vars, endpoints) need at least two words, must match the project's own vocabulary — which in a do-work repo lives in the URs (the user's own words), the prime files, and REQ titles, not primarily the README — and must be findable by plain-text search. Idiomatic short locals whose declaration-to-last-use fits one screen are fine — do not flag them. Endpoints: resources are nouns, actions come from the HTTP method, casing consistent.
2. Abstraction — repetition that wants a shared helper; helpers with a single caller; leaky layers (SQL in HTTP handlers, transport concerns in domain logic).
3. Consistency — the same thing done more than one way: error handling, config access, response shapes, logging, folder conventions. Inconsistency outranks any single ugly function.
4. Test quality — per hotspot: do tests assert behavior and contract, or restate implementation? Change-detector signals: assertions against mock internals, tests that break under behavior-preserving refactors, high mock-to-assertion ratio. Also the inverse gap: is the public contract covered at all? The ideal is a small test exercising a large surface, pinning its contract.
5. Discoverability and context coverage — can a newcomer, human or dispatched builder, load an area's purpose, invariants, and known traps within five minutes? Work from Phase 1's context-gap list: an area REQs keep returning to with no covering prime is a missing-prime finding — every dispatch there rebuilds the same context from scratch. Corroborating signal: the same facts re-derived across multiple archived REQ notes, or `## Lessons` fragments about one area scattered across several UR archives instead of consolidated. For this class, the impact evidence is the REQ-touch count (the tax recurs per REQ, so demand is the multiplier); the remedy is consolidation, not authorship — seed `prime-<area>.md` from the scattered lessons and archive notes; `Surface-cost: N/A` (a prime is consolidated context, not guard apparatus — the REQ-touch count is what earns its maintenance); the ratchet shape is a coverage assertion — any area above N archived REQ touches must appear in some prime's claimed scope.

Phase 3 — Consolidate by root cause
Never emit per-instance findings. Group instances under root-cause classes; emit at most `FINDINGS_MAX` classes, ranked by Score. One class = one finding = one verdict downstream. Fields per class:

* `sweep_key`: slug naming the root cause
* Claim (one sentence, refutable, specifics verbatim)
* Severity: P1 (user-visible or impact ≥ 4) / P2 (impact 3, or rule-change) / P3 (trivial)
* Label: VERIFIED | INFERRED
* `## Instances`: checklist of `path:line — one-line note`, each independently greppable
* Reproduce: the exact command that demonstrates the claim
* `gate`: `user-visible` | `rule-change` | `trivial`
* Impact 1–5 from churn × complexity × blast radius — show the three inputs
* `effort_estimate`: `trivial` | `normal` | `large`
* Score = impact / effort
* Remedy (one line) with a Surface-cost pre-classification — the validator prices any remedy that adds a guard, rule, or warning apparatus:
   * Prefer deletions, simplifications, renames, direct fixes → `Surface-cost: N/A`.
   * If the remedy adds surface — and a ratchet is added surface — supply the rubric evidence inline: incident = this class and its instances; replay = the Reproduce command; cost call = one CI line vs. the class recurring; test = the ratchet's red case, i.e. the Reproduce command returning hits.
* Ratchet: single number or zero-hit assertion, pinned at the current worst observed value, tightening as fixes land — usually the Reproduce command inverted. Never at the aspirational band. Name it; if none is possible, one sentence why.

Phase 4 — Report, then debrief
Write `audits/audit-YYYY-MM-DD.md`:

1. Header: date, audited commit SHA, dirty-tree note, the agreed bands from the Checkpoint, tool availability table with versions.
2. Metrics summary (with per-metric distributions) and, if a prior audit exists, a delta table (better / worse / unchanged per metric). The loop's convergence signal.
3. `## Findings — paste this section into: do-work validate-feedback` — the numbered class blocks, ranked by Score, fully self-contained (the validator may see this section alone). Claims only; remedies in their fields; zero imperatives.
4. Metrics appendix: WATCH items that did not become findings — visible, not actionable.
5. Pre-empted: candidates dropped because a documented decision, a recorded lesson, or a waiver covers them — one line each with the path to what covered it.
6. NOT-MEASURED: each with the declined or missing tool that closes it next run.

Then offer — don't launch — a short walkthrough: "want to talk through the top classes before you triage?" Exploring a finding together may reveal it belongs in Pre-empted, or that its remedy is wrong; update the report if so (still read-only outside `audits/`). Do not capture, do not fix, do not run validate-feedback yourself.
Loop usage (for the operator, not the agent)

1. Run this prompt → grounding + calibration conversation (tools, bands, scope) → audit runs → read `audits/audit-<date>.md`.
2. Paste the Findings section into `do-work validate-feedback`.
3. From the triage: capture the Accepts through its handoff (it preserves claim, severity, evidence, and surface-cost provenance); park Discuss items with `do-work-toolbox note`; for a Push back citing a documented decision you no longer agree with, change the decision doc — not the code.
4. `do-work run`.
5. Re-run the audit. Repeat runs skip straight to "reuse or recalibrate?"; the delta table must move; ratchets only ever tighten. A class you've decided to live with goes in `audits/waivers.md`, not into another round of triage.

---

Later clarifications in the same session (paraphrased where interactive): reports live in `do-work/audits/`; prime coverage is a judgment finding only; "can you use another word instead of ratchet?" → lock-in limit adopted; "keep committing and pushing the changes"; "ok, let me create the PR" (user creates the PR themselves).

Verbatim follow-ups during capture:

"another thing, since you also have a go tool, consider building tools for the audit, that will also output mechanicanically some flagged folders, files, etc... for the MVP it does not have too be too complex, but basically whatever we can script would be good to have it as script not as LLM call, becuase those are cheaper and more robust"

"also create some reqs as well already and cross link into the ai-reports/2026-08-13_2200_maintainability-audit-spec-validation/index.html"

---
*Captured: 2026-08-13T22:35:10Z*
