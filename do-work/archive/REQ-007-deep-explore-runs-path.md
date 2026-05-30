---
id: REQ-007
title: "Code review: move deep-explore session dirs to do-work/runs/"
status: completed
created_at: 2026-05-29T16:43:59Z
claimed_at: 2026-05-29T18:48:30Z
completed_at: 2026-05-29T18:49:30Z
route: B
review_generated: true
source: code-review
scope: actions/deep-explore.md, CLAUDE.md, crew-members/background-agents.md
---

# Code Review Fix: Move deep-explore Session Dirs to do-work/runs/

## What

`crew-members/background-agents.md:24-31` prescribes that any action fanning work out to background sub-agents creates its run directory at `do-work/runs/<action>-<YYYY-MM-DD-HHMMSS>/`. CLAUDE.md line 162 confirms this is the contract for "code-review, work multi-REQ, pipeline, **and deep-explore**".

But `actions/deep-explore.md:105-110` writes its session directory as `deep-explore-<concept-slug>-<timestamp>/` **at the project root**, not under `do-work/runs/`. The action then claims to follow the durability pattern at line 140 ("follow the durability pattern in `crew-members/background-agents.md`") — which only makes sense if the run dir is in `do-work/runs/`.

Consequences:
- Deep-explore session dirs are invisible to `forensics`, `cleanup`, and the rest of the `do-work/` tooling.
- They pollute the project root.
- The action declares one pattern and implements another, so the documented contract drifts from reality.

## Context

Found during code review of the full repo on 2026-05-29 (run `do-work/runs/code-review-2026-05-29-161332/`).

## Requirements

Pick one:

**Option A — Align to the convention (recommended):**
- Change `actions/deep-explore.md:105-110` to write session dirs at `do-work/runs/deep-explore-<concept-slug>-<ts>/`.
- Add a back-compat search in `continue` mode that also checks the legacy project-root `deep-explore-*-<ts>/` location for one release, with a deprecation note.
- Update `actions/deep-explore-reference.md` if it duplicates the path.
- Ensure `forensics` / `cleanup` are aware of (or naturally pick up) the new location.

**Option B — Carve out:**
- Update CLAUDE.md line 162 to remove `deep-explore` from the background-agents contract list.
- Update `actions/deep-explore.md` to drop the "follow the durability pattern" claim at L140 (or revise it to acknowledge the project-root layout is intentional and why — e.g., per-concept exploratory state, not per-run synthesis state, with a different resume model).

## Acceptance

- `actions/deep-explore.md`'s prescribed session-dir path and the background-agents contract agree.
- The CLAUDE.md `do-work/runs/` documentation matches reality across all fan-out actions.

## Source

Code review run: `do-work/runs/code-review-2026-05-29-161332/`
Finding: `architecture.md` F5

---

## Triage

**Route: B** — Medium

**Reasoning:** Clear outcome (align to the background-agents convention) but two implementation modes were on offer. Light exploration to confirm the actual session-dir-write call site, the keyword-search call site, the schema example in the reference companion, and that no third consumers (forensics, cleanup, roadmap) reference the legacy path.

**Planning:** Not required.

## Decisions

- **D-01:** Chose **Option A** (align to convention) over **Option B** (carve out). Reasoning: the REQ author recommended A; the deep-explore action already declares "follow the durability pattern in `crew-members/background-agents.md`" at L140 — Option A makes that claim true, Option B requires retracting it. Aligning is the cleaner architectural choice; carving out adds nuance for no functional gain.

## Exploration

- `actions/deep-explore.md:67` — `continue` mode keyword search (project root)
- `actions/deep-explore.md:105-110` — session-dir create call (project root)
- `actions/deep-explore.md:140` — declares `background-agents.md` durability pattern
- `actions/deep-explore-reference.md:323` — `session_dir` schema example
- `grep "deep-explore-" actions/forensics.md actions/cleanup.md actions/roadmap.md crew-members/background-agents.md` — empty (no other consumers of the legacy path)
- `.gitignore:2` — `do-work/runs/` is already gitignored (shipped in 0.83.2), so moving deep-explore sessions there keeps them out of the tracked tree

CLAUDE.md L162 already lists deep-explore as a `background-agents.md` caller; no edit needed there.

## Scope

**Files I will touch:**
- `actions/deep-explore.md` (modify) — Continue Mode resolution, session-dir create command, paragraph framing
- `actions/deep-explore-reference.md` (modify) — `session_dir` schema example
- (CLAUDE.md untouched — L162 is already correct for Option A)

**Files I will NOT touch:** `crew-members/background-agents.md` (already correct), `actions/forensics.md` / `actions/cleanup.md` (no current references; if they ever enumerate run directories, they'll naturally pick up the new path because `do-work/runs/` is one shared convention).

**Acceptance criteria (restated from REQ):**
- [x] `actions/deep-explore.md`'s prescribed session-dir path and the background-agents contract agree.
- [x] CLAUDE.md `do-work/runs/` documentation matches reality across all fan-out actions.

## Implementation Summary

**Files changed:**
- `actions/deep-explore.md` (modified) — Continue Mode Step 2 (keyword search) now globs `do-work/runs/` first and falls through to a one-release legacy project-root search with a deprecation warning. Step 2's session-dir creation writes to `do-work/runs/deep-explore-<slug>-<ts>/`. The paragraph framing explains the runs/-as-shared-fan-out-convention reasoning and the gitignored transience.
- `actions/deep-explore-reference.md` (modified) — state-file schema example `session_dir` field now shows `do-work/runs/deep-explore-<slug>-<timestamp>` to match the prescribed path.

**What was done:** Aligned deep-explore's session-directory location with the `crew-members/background-agents.md` durability convention. New sessions write to `do-work/runs/deep-explore-<slug>-<ts>/`; the `continue` mode searches `do-work/runs/` first and back-compat-searches the project root for one release with a deprecation warning, then that branch is scheduled for removal. The action's L140 claim to follow the background-agents durability pattern is now consistent with where state actually lives.

## Qualification

Passed — 2 files modified per scope. The Continue-mode resolution at L67 now starts at `do-work/runs/`, has an explicit legacy back-compat branch (option-A's requirement #2), and includes a removal target. The session-create command at L108 writes under `do-work/runs/`. The schema example at deep-explore-reference.md:323 matches. No drift between the action and the background-agents contract.

## Testing

**Tests run:** Manual cross-reference audit. Verified each occurrence of `deep-explore-` in `actions/deep-explore.md`, `actions/deep-explore-reference.md`, and the other fan-out callers (`actions/forensics.md`, `actions/cleanup.md`, `actions/roadmap.md`, `crew-members/background-agents.md`).
**Result:** ✓ All session-dir paths under `do-work/runs/`. ✓ No other consumers reference the legacy project-root path. ✓ Back-compat search path correctly delimited to a one-release deprecation window.

*Verified by work action*

## Review

**Overall: 88%** | 2026-05-29T18:49Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 85% |
| Test Adequacy | 75% |
| Scope | 95% |
| Risk | low |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — paths align with the background-agents contract; legacy back-compat is gated with a removal target.
**Suggested testing:** A future REQ could add a `forensics` enumeration of `do-work/runs/*` (it doesn't currently list them); that would close the loop the REQ raised about runs-dirs being visible to do-work tooling.
**Follow-ups created:** None — the suggested-testing item is a separate quality-of-life add, not an Important finding.

*Reviewed by work action (Route B self-review)*

## Lessons Learned

**What worked:** Greping the other fan-out actions (`forensics`, `cleanup`, `roadmap`) before editing confirmed there were no other consumers of the legacy path — so the change reduces to two files. The REQ's two-option framing made the trade-off explicit; recording the D-01 decision in the REQ preserves "why Option A" for anyone re-reading later.
**What didn't:** Briefly wrote a misleading line claiming the Writer step copies vision artifacts into `do-work/visions/` — that step doesn't exist. Corrected mid-edit. The fix: when documenting persistence, check the actual code, don't reason from analogy.
**Worth knowing:** The legacy back-compat search is scheduled for removal one release after 0.83.8. A future REQ should drop it explicitly rather than letting it accrete.

