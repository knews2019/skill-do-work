# HANDDOWN — UR-018: Parallel Building Across Checkouts (claim anywhere, one releaser)

> **SUPERSEDED 2026-08-04T21:38Z — the batch is built.** REQ-095 through REQ-101 all shipped
> (v0.170.2 → v0.174.2); this file's *Where the batch stands* table and per-REQ notes describe work
> that is now archived, and its "pending" column is wrong. Read `do-work/CHECKPOINT.md` for current
> state and `decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md`
> for what was decided. Two REQs remain: **REQ-104** (a known safety hole this batch's own acceptance
> run found — see the checkpoint) and **REQ-103** (`pending-answers`, waits for `do-work clarify`).
> Kept for its *Traps this session already hit* section, which is still accurate and still worth
> reading before touching these files.

**Written:** 2026-08-04T20:16Z, at v0.170.1 (HEAD after REQ-102's metadata commit).
**Batch spec:** `do-work/user-requests/UR-018/assets/approved-plan.md` is the authoritative plan; `do-work/user-requests/UR-018/input.md` records every user decision verbatim (read its "Ask-tool answers" block before overriding anything).

## Where the batch stands

| REQ | Status | Note |
|---|---|---|
| REQ-094 | **completed** v0.170.0, commit `9c305c0` | Writer label on checkpoint entries; 4-case crash-recovery classification |
| REQ-102 | **completed** v0.170.1, commit `44d4563` | Review-generated fix: Step 10 preserve rules cover label-less entries; both destruction paths suite-pinned |
| REQ-095 | pending (next in numeric order; deps satisfied) | Two-clone acceptance run — see per-REQ notes |
| REQ-096 | pending (deps satisfied) | Execution-model re-grain — **carries 2 addendum fold-ins** (see its `## Addendum`) |
| REQ-097 | pending, depends_on 096 | `assigned_to` field — schema + scan skip + `model.go` in the SAME commit |
| REQ-098 | pending, depends_on 097 | Two verify probes |
| REQ-099 | pending, depends_on 096 | Automatic wave dispatch (contract change to "nothing computes the set") |
| REQ-100 | pending, depends_on 099 | Live wave run — must prove overlapping builder timestamps |
| REQ-101 | pending, depends_on 096+097+099 | Docs + ADR (next ADR number after 017) |
| REQ-103 | **pending-answers** | Checkpoint frontmatter writer identity — waits for `do-work clarify`, do not build |

Queue is otherwise empty; `do-work/working/` is empty; `CHECKPOINT.md` has an empty In-Progress section. UR-018 stays in `do-work/user-requests/` until all members resolve (REQ-094/102 are in `archive/` root awaiting consolidation).

## What already shipped (and what it changed under you)

Checkpoint In-Progress entries now END with ` — writer: <hostname>:<absolute-checkout-path>` (`hostname -s` + `git rev-parse --show-toplevel`). **When you claim a REQ, write your entry in that format** — the contract suite pins it. Crash recovery classifies four cases (own-label / foreign-label / label-less / unnamed); foreign-label = report `claim held by <writer>, not touched`, never the takeover ladder. Step 10's rewrite and the session-start delete preserve every entry you did not write. If this handdown is being read on a *different checkout* (cloud): any surviving entries written by `t2s-Virtual-Machine:...` are foreign to you — leave them; there are none at the time of writing.

## Per-REQ notes for the next session

- **REQ-095 (two-clone run):** REQ-085's precedent artifact lives at `do-work/archive/UR-016/REQ-085-run-the-live-two-builder-acceptance-test.md` (its `## Testing` section, findings F-01…) — mirror that recording form inside the REQ file itself. Clones go in scratch space OUTSIDE the repo; dummy REQs only; never touch this repo's `do-work/`. The poisoning repro must test the OLD rule from `git show 9c305c0^:actions/work-reference.md` (pre-label prose) vs the NEW shipped rule. Evidence beats reasoning: if observed conflict shapes contradict shipped prose, fix the prose and say so.
- **REQ-096 (re-grain):** its `## Addendum` carries two fold-ins from REQ-094/102 reviews: (1) `actions/work-reference.md:55` + `docs/work-guide.md:91` "does not detect… a second owner" is already partly false (detection-by-label shipped in 0.170.0); (2) the Session Checkpoint Template's inline comment (~`work-reference.md:806`) still has the labeled-only scoping bug — reword to "every entry this checkout did not write". Contract constraints: `one queue owner per checkout` currently appears **exactly once** across `actions/` and a suite assertion counts it — the re-grain rewrites that very sentence, so UPDATE THE ASSERTION in the same commit (that is the sanctioned way; do not weaken other pins). `:57` "never probe / never arbitrate" survives verbatim (also pinned).
- **REQ-097 (`assigned_to`):** verbatim-read class (no normalization, alongside `write_set` at the `:206` paragraph). Board parse in `tools/queue-kanban/model.go` + badge + drawer row + test in the SAME commit (lock-step rule). Run the forbidden-token part of the suite early: `assigned_to` must not trip the reservation patterns (`reserved_for`/`reserved_at`/`status: reserved`/`do-work reserve`) — it doesn't by inspection, but prove it.
- **REQ-099 (auto wave):** rewrites `actions/work.md:33` ("does not drive a fan-out wave") and `work-reference.md` fan-out bullet "A human picks… Nothing computes the set" (`:320`-ish). Wave = pending, deps satisfied, unclaimed, not `assigned_to` someone else; bounded per `crew-members/background-agents.md:53`. `write_set` stays display-only — the wave computation must NOT schedule on it (multiple pins assert display-only language; keep them). Floor path (serial single-REQ) unchanged.
- **REQ-100:** overlapping wall-clock timestamps from ≥2 builders is the deliverable; REQ-085 got only Partial — beating that is the point.
- **REQ-101:** ADR next number = 018 (records/ stops at adr-017); also append one line to `decisions/log.md` (stale since 2026-07-01 — just append).

## Traps this session already hit (don't re-hit)

- `tools/checks/record-commit-hash.sh`: flag order is `--verify <file> <hash>` — `<file> <hash> --verify` errors after the fact.
- `tools/checks/preflight.sh`: pass the suite as `"bash _dev/tests/contract-regressions.sh"` (direct path hits Permission denied — not executable).
- The compiled `queue-kanban` binary can be STALE — rebuild (`cd tools/queue-kanban && go build -o queue-kanban .`) before trusting `verify` findings (a stale binary reported a ghost-REQ false positive fixed in 0.169.9).
- Contract-suite pinned phrases near this batch's files: `never grow into one`, `absent checkpoint is ambiguous`, `foreign claim`, `Crash Recovery's input`, `claim held by`, `writer: <hostname>:<absolute-checkout-path>`, `entry this checkout did not write through verbatim`, `no entry this checkout did not write remains`, Step-1 line order (checkpoint read before **Crash Recovery:**). Reword AROUND them.
- Builders paraphrasing a canonical condition at echo sites caused both follow-ups so far — make builders QUOTE the canonical wording at echo sites.
- Serial-only, integrator-only: `actions/version.md` + `CHANGELOG.md` bumps happen once per REQ at Step 9, by the queue owner. Changelog titles say what shipped (no codenames). Verify title-uniqueness + version-monotonicity with `queue-kanban verify` after writing.

## Process reminders

Full pipeline per REQ (`actions/work.md`): claim → checkpoint entry (labeled!) → triage → explore (B/C) → scope+write_set mirror → preflight → builder (fresh agent per REQ; crew: general + coding-guardrails always) → Implementation Summary → qualify.sh → tests (orchestrator-run, not builder-claimed) → review (follow-ups for Important findings) → lessons/orientation → archive → changelog+version → commit → record-commit-hash + metadata commit. KB handoff: unattended default is `kb_status: pending` (REQ-094/102 both carry it; a later `bkb triage` sweeps).
