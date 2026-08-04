# Parallel Building — Final Plan (all decisions made)

## Context

The user wants sanctioned parallel building across every instance shape — one session with parallel builders, multiple local sessions/workspaces, clones, cloud sessions — unified by one cooperative rule: a claimed/assigned REQ is visible to every other do-work instance and left alone. The reserve mechanism that once did this (0.125.0) was deleted at 0.163.0 as collateral of the exclusive-session cleanup (0.161.0), not for defects. This plan partially re-grains that contract along the optimistic philosophy the user chose at every fork: **no prevention machinery anywhere; conflicts and duplicates are fixed at merge; discovery stays cheap via existing probes.**

## Decisions (user-confirmed, 2026-08-04)

1. **Claim marker**: `assigned_to` advisory frontmatter field only. No reserve verb, no `reserved` status, no staleness clock — the 0.163.0 forbidden-token ratchet and SKILL.md router budget stay untouched.
2. **Instance shapes**: own tree per instance (spawned worktree, user workspace, clone, cloud sandbox). Same-tree dual sessions stay *unspecified* — no prevention built, damage repaired via forensics/cleanup as needed.
3. **Ownership — claim anywhere, one releaser**: any checkout may capture AND claim/build. Exactly one designated releaser checkout runs the release tail: merge integration, version bump, `CHANGELOG.md` entry, archive moves, UR closure.
4. **Capture everywhere**: duplicate REQ ids are a merge-time fix; `queue-kanban verify`'s `duplicate-req-id` probe is the cheap detector.
5. **Wave dispatch — fully automatic set-picking**: the work loop computes the ready set and dispatches builders without a confirmation gate. Deliberate contract change to "a human picks / nothing computes the set."
6. **Checkpoint**: `CHECKPOINT.md` stays tracked and committed (history = version trail). Crash recovery gains a **static writer label** — entries record which checkout wrote them; recovery ignores foreign entries (fixes the deterministic cross-checkout claim-stripping replay of the 2026-07-01 incident). The `:413` tripwire is amended: refresh intervals, staleness checks, liveness probes stay banned by name; a static writer label is not liveness machinery.

## Build (each phase lands as normal REQs via capture → work, versioned per convention)

### Phase 1 — Checkpoint writer label (the landmine, independent value)
1. `actions/work-reference.md` In-Progress Record (~:409–425): entries carry `writer: <hostname>:<checkout-path>` (path alone collides across machines). Recovery rule: only own-label entries are crash-recovery candidates; foreign entries are **reported, never stripped** (extends the `:239` foreign-claim rule to the checkpoint). Reword the `:413` tripwire accordingly.
2. Two-clone acceptance run: reproduce the poisoning (clone B strips clone A's live claim), confirm the writer label stops it; also claim the same REQ in two clones and capture the real merge-conflict text. Record like REQ-085's fan-out run.

### Phase 2 — Contract re-grain + satellite generalization
3. `actions/work-reference.md:53–61` Execution Model rewrite: any checkout captures/claims; **one releaser per queue** owns the release tail; two releasers = unspecified; same-tree dual sessions = unspecified; ":57 never probe / never arbitrate" survives verbatim.
4. `actions/work-reference.md:275–341` Worktree Dispatch widened: a builder tree may be a spawned worktree, user workspace, clone, or remote sandbox. Deltas: remote hand-back travels on the branch (absolute-main-tree-path handback is local-only); a non-releaser checkout treats its synced `do-work/` snapshot as potentially stale; claim conflicts between checkouts are ordinary git conflicts fixed at merge. Red Flag: a second checkout running the *release tail* is the violation to watch — claiming/building elsewhere is now in contract.
5. `assigned_to` schema line in the Schema Read Contract (verbatim-read class, alongside `write_set`), plus skip-and-report in `actions/work.md` Step 1 default scan; explicit `do-work run REQ-NNN` overrides and clears it.
6. `tools/queue-kanban/model.go`: parse `assigned_to` display-only (badge + drawer row) — **same commit** as the schema line (lock-step rule). ~15 lines + test.
7. `tools/queue-kanban/verify.go`: probe `assigned-elsewhere-claimed-here` (an `assigned_to` REQ sitting in `working/`); probe `ur-archived-with-live-member`. ~30 lines each + tests.

### Phase 3 — Automatic wave dispatch
8. Amend `actions/work.md:33` (one-REQ-at-a-time stance) and `actions/work-reference.md:320` ("nothing computes the set"): the loop now computes the wave — pending REQs with dependencies satisfied, unclaimed, not `assigned_to` another session — sized per `crew-members/background-agents.md:53` bounded-wave guidance, then dispatches builders and integrates **serially** (merge → qualify → review → changelog → archive stays one at a time; `:321` "the merge is the non-interference proof" becomes the load-bearing sentence and survives unchanged; `write_set` stays display-only).
9. Live auto-wave acceptance run — real wall-clock concurrency has never been proven (REQ-085: Partial); this is the proof, not ceremony.

### Phase 4 — Docs + ADR
10. `docs/work-guide.md`: "several checkouts against one queue" user-facing section.
11. ADR in `decisions/records/`: the re-grain; why the reserve *verb/status* stays dead while the *field* returns; capture-anywhere with fix-at-merge; the auto-wave contract change; the writer label vs. the liveness ban. (No ADR currently covers session ownership at all.)

## Do NOT build (named failures behind each)

- Locks, leases, heartbeats, refresh intervals, staleness checks, liveness probes (REQ-018 history: 3 patches, one reintroduced the incident).
- Auto-release/takeover on any staleness threshold; no `assigned_at`.
- `reserve`/`release` verbs or `reserved` status (ratchet + router budget).
- `write_set`-based scheduling (display-only at any builder count).
- Sharded REQ-id ranges (fix-at-merge + `duplicate-req-id` probe instead).
- Auto `git pull/push` inside any action (mid-run pull changes the queue under a live claim).

## Verification

- Phase 1 two-clone run and Phase 3 live wave run are the model's proof artifacts; record both.
- `_dev/tests/contract-regressions.sh` stays green with **no ratchet weakening** (no reservation tokens reintroduced; `assigned_to` must not trip the underscore-token patterns — verify).
- `go test ./...` in `tools/queue-kanban/` (model parsing + two new probes).
- Board smoke: `assigned_to` badge renders, no bucket change; legacy `reserved` card behavior unchanged.
