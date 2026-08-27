---
id: REQ-384
title: 'Deliver the architecture report as a freeform HTML bundle'
status: completed
created_at: 2026-08-26T17:56:52Z
user_request: UR-077
domain: general
prime_files: ['_dev/primes/prime-action-files.md', '_dev/primes/prime-shell-commands.md']
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: '2026-08-27T11:45:42Z'
status_changed_at: '2026-08-27T11:45:42Z'
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  basis:
  - Route B
  - 6-file write set
  - 9 acceptance criteria
  - cross-route regression gates
  - full-suite verification
  calculated_at: '2026-08-27T11:45:42Z'
write_set:
- skills/do-work-toolbox/actions/architecture-report.md
- skills/do-work-toolbox/scripts/architecture-report-preflight.sh
- _dev/tests/contract-regressions.sh
- _dev/tests/prescribed-shell-cases/architecture-report-preflight.sh
- skills/do-work-toolbox/actions/help.md
- README.md
completed_at: '2026-08-27T11:57:36Z'
commit:
kb_status: pending
---

# Deliver the Architecture Report as a Freeform HTML Bundle

## What

Change `do-work-toolbox architecture-report` so a run publishes one beautifully rendered, self-contained HTML report — `ai-reports/<yyyy-mm-dd>_<hhmm>_architecture-report/index.html`, the same bundle home `ai-report` uses — and no markdown file at all. The HTML is deliberately freeform: the action states the quality bar and the invariants, never a fixed section-by-section layout, so a more capable future model produces a better architecture view instead of the same template filled in better.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Update existing contract and preflight fixtures first; observe RED. Replace Markdown/carry-forward rules with HTML-only freeform quality and provenance requirements; update helper and discovery prose, then prove GREEN.
- [x] **[APPLY]:** Implemented the HTML action/helper and updated existing tests/discovery prose within the six declared files.
- [x] **[UNIFY]:** Parent reviewed all six files listed in Implementation Summary: action quality/provenance contract, helper publication and scan, both existing test seams, README/help discovery. Targeted helper fixture and diff check passed; full gate recorded below after completion.

## Why

The user wants the architecture view "beautifully rendered so it's easy to understand and so it has all the bells and whistles." The current action mandates the opposite ("Markdown with inline Mermaid only. No HTML."), and its byte-identical carry-forward contract locks every future report to the previous one's authoring — which is exactly what the user does not want: "the HTML file needs to be freeform in the sense that as models evolve, I keep getting a better architecture view into it."

## Context

Captured after an interactive session that resolved the design decisions below; the full record is in UR-077. The action being changed is `skills/do-work-toolbox/actions/architecture-report.md`, whose current contract is markdown-only, carry-forward byte-identical, delta-first — and that contract is pinned by regression tests (`_dev/tests/contract-regressions.sh` pins repo-wide input, dated-immutable publication, carry-forward, and delta-first; `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh` proves the preflight helper), all of which must move with the action rather than be deleted around.

## Detailed Requirements

- **HTML only.** "I don't want the markdown file at all" — a run creates one dated bundle holding `index.html`; no `architecture-report.md` is written anywhere.
- **The ai-report home.** "at the same place as any ai-report lives": a dated bundle under `ai-reports/` following the shared Collision-Safe Publication contract (`skills/do-work-toolbox/actions/completed-work-presentation-reference.md`) — self-contained, relative assets, never touching an existing bundle.
- **Freeform by contract.** The action must not prescribe the HTML's structure. It states goals (a reader understands the architecture; the report is redesigned each run, free to open with an executive overview and tell the story visually) and leaves layout, sectioning, and visual design to the authoring model.
- **Bells and whistles.** Rendered diagrams (drawn, not fenced code), a clickable section navigation, and a considered visual design are the floor, not the ceiling.
- **Verification discipline survives.** Claims stay labeled `VERIFIED` (anchored) or `INFERRED` (basis stated); anchors render as links to the file and line on GitHub; verification still runs against a committed tree, and the run stays unattended.
- **Authored delta, not byte-diff.** Each report opens with a "changed since last report" section written by reading the previous HTML report. The carry-forward byte-identical rule and the diff-two-reports contract are removed, not reworded.
- **Prior-report lookup ignores markdown bundles.** The existing `ai-reports/2026-08-26_1709_architecture-report/` markdown bundle stays as committed history and is never a prior baseline for an HTML run.
- **Immutability survives.** Dated bundles are never edited, deleted, or regenerated; a new run is a new bundle.
- **Capability only.** This REQ lands the changed action, helper, and tests. It does not run the action — the first HTML report is produced by a later invocation.

## Constraints

- Regression tests pinning the old contract (`_dev/tests/contract-regressions.sh`, `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh`, and any staged-skills references) are updated in the same change to pin the new contract; `bash _dev/tests/maintainer-verify.sh` exits 0 at hand-back.
- The action file keeps the agent-compatibility floor (`CLAUDE.md` → Agent Compatibility): generalized language, works as a standalone prompt.
- Routing rows, help menus, and the README's mention of the action stay truthful about the new output.
- How much of `scripts/architecture-report-preflight.sh` survives is the builder's call — the bundle-naming and collision-escalation mechanics likely stay; the markdown-specific parts go.

## Builder Guidance

Certainty: **Firm** on the deliverable (HTML-only, ai-report home, freeform, authored delta, verification labels kept, markdown history kept). **Latitude** on how the action words the quality bar, what the preflight helper becomes, and how the tests pin "freeform" without prescribing structure.

## Red-Green Proof

**RED prompt/case:** Run `do-work-toolbox architecture-report` today. It publishes `ai-reports/<slug>/architecture-report.md` — raw markdown; diagrams are fenced code blocks, anchors are plain text, and the action file itself mandates "Markdown with inline Mermaid only. No HTML."
**Why RED now:** The action's output contract forbids exactly the deliverable the user wants.
**Runnable RED (test-first):** Rewrite the architecture-report pins in `_dev/tests/contract-regressions.sh` and the preflight fixture proofs in `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh` to assert the new contract (bundle holds `index.html`, no markdown file, no carry-forward requirement); confirm they fail against the current action and helper, then make them pass. The browser observations above remain the user-facing GREEN; this is the harness-runnable core the TDD gate checks.
**GREEN when:** After this REQ lands, a run of the action publishes `ai-reports/<yyyy-mm-dd>_<hhmm>_architecture-report/index.html` and no markdown file; opening it in a browser shows the diagrams drawn, a clickable section navigation, each `VERIFIED` anchor linking to its file and line on GitHub, and — when a prior HTML report exists — an opening "changed since last report" section.
**Validation:** User adjusted — confirmed the inferred pair, then added: no markdown file at all, and the HTML must stay freeform so future models keep improving the view. `tdd:` was corrected to true after capture, from review on PR #171: the contract suite supplies the runnable RED above.

## Full Context

See `do-work/user-requests/UR-077/input.md` for complete verbatim input and the capture-time decisions.

---
*Source: "basically I want it beautifully rendered so it's easy to understand and so it has all the bells and whistles" — see UR-077 for the full record.*

## Triage

**Route: B** — Known action/helper change with existing contract and fixture seams; inspect those before a focused test-first rewrite.

## Plan

Planning not required — focused implementation guided by the request and existing patterns.

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/architecture-report.md` (modify)
- `skills/do-work-toolbox/scripts/architecture-report-preflight.sh` (modify)
- `_dev/tests/contract-regressions.sh` (modify)
- `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh` (modify)
- `skills/do-work-toolbox/actions/help.md` (modify)
- `README.md` (modify)

**Acceptance criteria (restated from REQ):**
- All Detailed Requirements and the captured Red-Green Proof are satisfied.

## Exploration

The action mandates Markdown and numbered layout, the helper fixes its report filename and reads a textual watermark, and contract-regressions plus the preflight fixture pin both. README and toolbox help describe the old output. Routing and staged-skills references currently describe only command/path names and need no change.

Scope acceptance: HTML-only index.html bundles; collision-safe immutable publication; freely redesigned layout with rendered diagrams/navigation; committed-tree verification and GitHub anchors; authored opening delta from prior HTML only; no report execution in this REQ.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/architecture-report.md` (modified) — HTML-only freeform composition, authored opening delta, committed GitHub evidence, and offline visual verification.
- `skills/do-work-toolbox/scripts/architecture-report-preflight.sh` (modified) — HTML prior selection and metadata, immutable publication with hidden partial copies.
- `_dev/tests/contract-regressions.sh` (modified) — replaces the previous Markdown/carry-forward expectations with REQ-384's HTML contract.
- `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh` (modified) — HTML publication/metadata/collision cases plus ignored legacy and failed bundles.
- `skills/do-work-toolbox/actions/help.md` (modified) — advertises the HTML architecture map.
- `README.md` (modified) — explains the new artifact and authored comparison.

**What was done:** Changed the architecture-report capability, without invoking it or generating a report. Existing Markdown history remains unchanged; each future report has one self-contained index.html, an authored opening change account, and freely designed diagrams/navigation.

## Decisions

- D-01 (decide and state): Embed presentation assets in one self-contained index.html. This satisfies offline rendering without adding an asset-copying interface.
- D-02 (decide and state): Put the verified commit in an exact HTML meta element. Machine-readable provenance does not impose a visible section template.
- D-03 (decide and state): Keep exclusive directory claiming and numeric suffix ordering. Publish the visible HTML filename only after validating a hidden copy, so interrupted bytes cannot become a prior baseline.

## Qualification

qualify.sh exit0 after the mandatory manifest and P-A-U states were completed. Parent reviewed each declared file and traced substantive changes to the HTML capability. The citation metadata flows from authored HTML to prior lookup; no queue write surface or source/report execution was added.

## Testing

**Test-first RED:** Existing preflight fixture rewritten before production edits exited1, nine cases with12 failures. The full contract suite also exited1 on the old Markdown/carry-forward contract and staged helper behavior. **GREEN:** Helper fixture now exits0, nine cases with0 failures; parent independently reran it. Shell syntax, ShellCheck with sourced harness, and git diff --check passed.

**Existing tests intentionally changed:** REQ-384 supersedes the old architecture-report Markdown/carry-forward expectations in the two existing test files. Repo-wide input, committed-tree provenance, audit separation, numeric collision ordering, and unreadable/unresolvable watermark behavior remain pinned. New cases ignore legacy Markdown and incomplete bundles, and reject partially copied HTML as a prior baseline.

**Capability-only boundary:** No action invocation or new architecture report. Existing ai-reports bytes remain unchanged. Real report visual verification belongs to the later invocation, as requested.

## Orientation

Run do-work-toolbox architecture-report later to author the first HTML report. It will publish a new dated index.html with embedded presentation assets, drawn diagrams, navigation, source evidence, and an opening account of changes since the last HTML baseline.

## Lessons Learned

A filename used as a publication marker should become visible only after the copy is complete and verified. Separating machine-readable metadata from visible structure allows deterministic history lookup without fixing the report's design.

## Review

**Overall: 100%** | Acceptance: Pass for capability scope. Independent reviewer inspected the full six-file diff and original input, reproduced RED against the old helper and GREEN against the new helper, and tested eight simultaneous publications, empty HTML exclusion, corrupt-copy rejection, and safe retry. No actionable findings. Active restatements agree; historical changelog/report content remains unchanged. Real report generation and visual inspection are deferred by the explicit capability-only requirement.

The builder's complete post-implementation contract-regressions.sh run also exited0. The parent runs the canonical maintainer gate separately before archive/commit.

## Repository Verification

`bash _dev/tests/maintainer-verify.sh` completed with exit 0 on the final implementation/release state. Contract suite, Go checks, and strict JavaScript lane passed. The default optional browser lane was explicitly skipped; separate browser evidence is recorded above where applicable.
