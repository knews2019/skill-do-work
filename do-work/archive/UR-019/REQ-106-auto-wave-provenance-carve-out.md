---
id: REQ-106
title: Auto-wave ready set contradicts targeting-token provenance — add the carve-out to work-reference
status: completed
created_at: 2026-08-05T09:43:47Z
claimed_at: 2026-08-05T10:37:10Z
completed_at: 2026-08-05T10:38:39Z
commit: 84d83cc
route: A
user_request: UR-019
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set: [actions/work-reference.md, actions/work.md]
related: [REQ-099, REQ-105, REQ-107]
batch: sync-review-0174
---

# Auto-Wave Ready Set Contradicts Targeting-Token Provenance

## What

`actions/work.md` says `--fan-out` "composes with everything that selects a set, `--wave N` and targeting tokens included" and that an explicitly-named `REQ-NNN` bypasses `depends_on`. But `actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch → **Auto-wave** defines the ready set with condition 2 ("Dependency-ready") stated unconditionally. Reading only the reference, an explicitly-named but dependency-blocked REQ is excluded from the wave; reading work.md, it's included. Make both files state the same rule.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Triage found `actions/work.md` already internally coherent (its Step 1 auto-wave paragraph states every filter — targeting provenance included — applies to the wave, "no separate readiness predicate"). So the fix is reference-side only: append the serial scan's provenance carve-out to Auto-wave condition 2 (mirroring the carve-out sentence condition 4 already carries) and add one paragraph after the list covering targeted-mode pool scoping.
- [x] **[APPLY]:** Both edits landed in `actions/work-reference.md`; `actions/work.md` deliberately untouched (see D-01).
- [x] **[UNIFY]:** `git diff --stat` shows only `actions/work-reference.md` (+4/−2 lines net) beyond do-work/ bookkeeping. `bash _dev/tests/contract-regressions.sh` passes. Restatement sweep run (see Review). No debug artifacts — prose-only change. Verified: condition 2's new sentence reads in the list's voice and cites `actions/work.md` Step 1; the new paragraph sits between the list and the existing "Nothing else enters the computation" paragraph without contradicting it.

## Detailed Requirements

- The intended rule (per `actions/work.md`'s per-token provenance contract, which is the richer statement): in a targeted `--fan-out` run, an **explicitly-named** REQ enters the wave regardless of `depends_on`; a REQ reached by **UR-expansion** still goes through the dependency-ready filter scoped to the UR's member set. Write that as a provenance carve-out in the Auto-wave section, mirroring how condition 4 (`assigned_to`) already carries its own carve-out sentence ("Explicit targeting still overrides it").
- While there, state how the four conditions apply when targeting tokens scope the wave at all — today the Auto-wave list reads as a whole-queue computation and never mentions targeted mode, which is the root of the ambiguity.
- Do not change the serial rule or the default-mode (untargeted) auto-wave predicate — this is an alignment, not a behavior change. If wording in `actions/work.md` needs a pointer adjustment, keep both files agreeing in the same commit.

## Context

Found by a downstream consumer's review of the 0.170.1 → 0.174.3 sync; verified here at triage. Evidence: `actions/work.md:104` (composes-with claim), `actions/work.md:186` and `:200` (named bypasses `depends_on`; UR-expanded does not), versus `actions/work-reference.md` Auto-wave conditions 1–4 where only condition 4 carries a targeting carve-out. Auto-wave shipped in REQ-099 under ADR-018.

## Red-Green Proof
**RED prompt/case:** An agent reading only `actions/work-reference.md`'s Auto-wave list, asked "does `do-work run REQ-042 --fan-out` dispatch REQ-042 when its `depends_on` is unmet?", answers no (condition 2 excludes it) — while an agent reading `actions/work.md` answers yes (explicit naming bypasses `depends_on`, and `--fan-out` composes with targeting tokens).
**Why RED now:** Condition 2 is stated unconditionally; the reference never addresses targeted-mode waves, so the two files give contradictory answers.
**GREEN when:** Both files give the same answer: named → in the wave, UR-expanded → dependency-gated. The Auto-wave section carries the provenance carve-out explicitly.
**Validation:** Inferred during capture (triage-verified against the repo; the carve-out direction follows work.md's existing per-token provenance contract)

## Full Context
See `do-work/user-requests/UR-019/input.md` for complete verbatim input.

## Triage

**Route: A (Direct to Builder)** — both files named in the REQ; the intended rule is already fully stated in `actions/work.md`, so the change reduces to aligning the reference's Auto-wave list with it.

## Decisions

- **D-01** (DECIDE & STATE): `actions/work.md` was declared in `write_set` but left untouched. Its Input section (composes-with claim), Step 1 targeted-mode paragraph, and auto-wave paragraph already state the intended rule consistently — editing it would restate what it already says. The contradiction existed only for a reader of the reference alone; the reference now carries the carve-out. Reversible (a later REQ can still touch work.md); reach is zero (no behavior change to work.md readers).

## Implementation Summary

**What was done:** Aligned the Auto-wave ready-set definition with targeting-token provenance. Condition 2 (Dependency-ready) now carries the serial scan's carve-out — explicitly-named REQs enter the wave regardless of `depends_on`, UR-expanded members stay gated scoped to the UR's set — mirroring the carve-out sentence condition 4 already had. A new paragraph after the four-condition list states that targeting tokens scope the candidate pool and per-token provenance survives into the wave, with "no separate readiness predicate for waves" restated from `actions/work.md`.

**Files changed:**
- `actions/work-reference.md` (modified) — Auto-wave section: condition 2 carve-out sentence + post-list "Targeting tokens scope the pool" paragraph

**Restatement sweep:** `actions/work.md` lines 35/103/186/200/213 already state or summarize the rule consistently — no stale restatement. ADR-018 consequence 4 restates the default-mode computation; as a historical decision record it stays as written (the carve-out refines, does not reverse, the recorded decision). CHANGELOG 0.174.0 likewise historical.

**Tests:** `bash _dev/tests/contract-regressions.sh` passes.

## Qualification

Passed — 1 file verified (`tools/checks/qualify.sh` OK), requirements traced (carve-out → condition 2; targeted-mode statement → new paragraph; both-files-agree → D-01 records why work.md needs no edit), P-A-U confirmed against the diff.

## Testing

- `bash _dev/tests/contract-regressions.sh` — passes.
- **Red-green validation:** RED (captured): reading only the reference's Auto-wave list, "does `do-work run REQ-042 --fan-out` dispatch REQ-042 when its `depends_on` is unmet?" answered no; reading work.md, yes. GREEN: the reference's condition 2 now answers yes for a named REQ and no for a UR-expanded one — the same answer work.md gives. Verified by re-reading both passages post-edit.

## Review

**Acceptance: Pass** (Route A quick scan, calibrated depth). The carve-out direction follows work.md's per-token provenance contract as the REQ required; wording mirrors condition 4's existing carve-out idiom so the list stays internally consistent; the new paragraph does not contradict the adjacent "Nothing else enters the computation" paragraph (provenance is part of *how the four conditions apply*, not a fifth input). Restatement sweep reported above — no stale restatements outside historical records. Scope: `actions/work-reference.md` touched as declared; `actions/work.md` declared-but-untouched, resolved as deliberate in D-01. No findings.

## Orientation

Auto-wave and the serial scan now give the same answer on a targeted, dependency-blocked REQ: named → runs, UR-expanded → gated. Lives in the work pipeline's fan-out contract (work-reference's Auto-wave section); no map change.

---
*Source: downstream sync-review finding 2, verified at triage — see UR-019*
