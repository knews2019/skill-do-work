---
id: REQ-378
title: 'Deliver the architecture report as a freeform HTML bundle'
status: pending
created_at: 2026-08-26T17:56:52Z
user_request: UR-075
domain: general
prime_files: ['_dev/primes/prime-action-files.md', '_dev/primes/prime-shell-commands.md']
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
---

# Deliver the Architecture Report as a Freeform HTML Bundle

## What

Change `do-work-toolbox architecture-report` so a run publishes one beautifully rendered, self-contained HTML report — `ai-reports/<yyyy-mm-dd>_<hhmm>_architecture-report/index.html`, the same bundle home `ai-report` uses — and no markdown file at all. The HTML is deliberately freeform: the action states the quality bar and the invariants, never a fixed section-by-section layout, so a more capable future model produces a better architecture view instead of the same template filled in better.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The user wants the architecture view "beautifully rendered so it's easy to understand and so it has all the bells and whistles." The current action mandates the opposite ("Markdown with inline Mermaid only. No HTML."), and its byte-identical carry-forward contract locks every future report to the previous one's authoring — which is exactly what the user does not want: "the HTML file needs to be freeform in the sense that as models evolve, I keep getting a better architecture view into it."

## Context

Captured after an interactive session that resolved the design decisions below; the full record is in UR-075. The action being changed is `skills/do-work-toolbox/actions/architecture-report.md`, whose current contract is markdown-only, carry-forward byte-identical, delta-first — and that contract is pinned by regression tests (`_dev/tests/contract-regressions.sh` pins repo-wide input, dated-immutable publication, carry-forward, and delta-first; `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh` proves the preflight helper), all of which must move with the action rather than be deleted around.

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
**GREEN when:** After this REQ lands, a run of the action publishes `ai-reports/<yyyy-mm-dd>_<hhmm>_architecture-report/index.html` and no markdown file; opening it in a browser shows the diagrams drawn, a clickable section navigation, each `VERIFIED` anchor linking to its file and line on GitHub, and — when a prior HTML report exists — an opening "changed since last report" section.
**Validation:** User adjusted — confirmed the inferred pair, then added: no markdown file at all, and the HTML must stay freeform so future models keep improving the view.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input and the capture-time decisions.

---
*Source: "basically I want it beautifully rendered so it's easy to understand and so it has all the bells and whistles" — see UR-075 for the full record.*
