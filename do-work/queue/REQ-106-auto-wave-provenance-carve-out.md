---
id: REQ-106
title: Auto-wave ready set contradicts targeting-token provenance — add the carve-out to work-reference
status: pending
created_at: 2026-08-05T09:43:47Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

---
*Source: downstream sync-review finding 2, verified at triage — see UR-019*
