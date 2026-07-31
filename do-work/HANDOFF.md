# Session Handoff — 2026-07-29 (session do-work-20260728T211058Z-20017)

Continuation document for the next session, written at user request mid-queue. Companion: `do-work/CHECKPOINT.md` (the work loop's native resume file — keep both in sync if you edit one).

## TL;DR

`do-work run` processed UR-007's three captured REQs (032/033/034) — all built, adversarially reviewed, archived, and committed (HEAD `9c7e3be`, version **0.143.0**). Reviews confirmed real spec holes and queued five follow-up REQs (035–038, 040), all `pending` and ready; two `pending-answers` REQs (039, 041) wait on the user. Orchestrator lock **released**, `working/` empty, tree clean. Restart with `do-work run` — it picks up REQ-035 first.

## State at handoff

- **Queue:** REQ-035, 036, 037, 038, 040 `pending` (no unmet deps — 035/036 depended on REQ-033, now completed); REQ-039, 041 `pending-answers` (→ `do-work clarify`).
- **Archive (root):** REQ-032, 033, 034 completed with full trails + `commit:` hashes. UR-007 folder stays in `user-requests/` until all its REQs finish.
- **Lock:** deleted (clean release — holder with no coexisting sessions). CHECKPOINT.md was rewritten at exit, so no stale-checkpoint confusion.
- **Commits this session** (local only — never push from this repo): `97f65b9` REQ-032 (0.141.0) · `849b2a5` REQ-033 (0.142.0) · `931bf2c` changelog date fix · `9c7e3be` REQ-034 (0.143.0).

## Head start on REQ-035 (do not redo)

REQ-035's file already carries **validated Triage (Route C), Plan, and Exploration sections** — the pipeline's per-section idempotence means the next run resumes at **Step 5.5 (Scope)**. Key plan facts:

- Chosen representation: **`claimed_reqs` list** on the holder and on each `coexisting_sessions` entry; `claimed_req` kept as a derived legacy mirror (`claimed_reqs[0]` / `null`). Additive for back-compat (an old reader ignores the new field; an array-shaped `claimed_req` would break every old reader).
- **Mixed-version co-dispatch precondition:** co-dispatch only when this session is the only live claimant in the lock.
- **Crash Recovery rule change:** the "session *other than this one*" clause is **deleted, not extended** — skip any `working/` file whose id appears in any fresh (≤45m) claim set; freshness alone gates (a restarted session gets a new session_id, so a dead predecessor's claims age out and recover).
- ⚠ The explorer found **three write sites the plan's edit list missed**: `actions/work.md:577` (Step 8 failure path), `:582` (blocked flip), `actions/work-reference.md:352` (acquisition). The builder must cover all `claimed_req` sites — the REQ's requirement is "one story everywhere."
- `actions/cleanup.md` Pass 0 names the field → **extend REQ-035's write_set with it at Scope time** (the one-directional mirror rule allows the Scope step to widen the set).
- Per-merge verification becomes the default whenever >1 REQ is in flight — this deliberately discharges REQ-033's dormant "per-batch rollback" finding; keep that clause.

## Run configuration the user chose (calibrated rigor)

- **Full adversarial review** (review agent + independent contradiction-hunter + 2 refuters per Important+ finding, as a Workflow) for **REQ-035, 036, 037** — the concurrency-spec surface where this pass confirmed 4 real root causes across two REQs today.
- **Standard single-reviewer pass** for **REQ-038 and 040** (small, well-scoped).
- Follow-up REQs for confirmed Important findings only; Minors stay in the report. Frequent per-REQ commits (background commit agent; version bump + descriptive changelog title each time; stage explicit files only; `Co-Authored-By: Claude Fable 5 <noreply@anthropic.com>` trailer).

## Environment facts (verified this session)

- `do-work/` is git-excluded via `.git/info/exclude` — **stage nothing under it.** This note deliberately carries **no local variant of the `commit:` write-back** (it used to say "write hashes directly into archived REQs," which is exactly the free-form edit `tools/checks/record-commit-hash.sh` was built to replace, and which blanked six archived REQs in a consumer repo). Follow the standard procedure as written — `actions/work-reference.md` → **Commit & Metadata-Commit Procedure (Step 9)** — and treat it as the only instruction. Two consequences fall out of the exclusion, and only these two: the script runs normally on an untracked path but prints `INFO: … is not tracked by git` and skips its size-floor and diff guards (every content guard still runs), and there is **no metadata commit to make** afterward, since `git commit` has nothing to stage. For the same reason, don't run the script's `--verify` mode here — it reports `FAIL: … is not tracked by git` by design, not because anything went wrong.
- Baselines all green at handoff: `bash _dev/tests/contract-regressions.sh`, `cd tools/queue-kanban && go test ./...`, `go vet`, `gofmt -l` clean.
- SKILL.md word budget: **2648/2650** — two words of headroom; add no routing prose.
- Changelog dates use the **local** date (2026-07-29), not `date -u` (which was still 07-28 during this session — bit us once, fixed in `931bf2c`).
- Transient subagent "Login expired" API failure hit once (REQ-034's builder, mid-run). Recovery pattern that worked: check the partial diff compiles → relaunch a continuation builder against the diff + REQ trail → have it reconcile P-A-U/Decisions truthfully. If it recurs repeatedly, the user needs `/login`.
- One review Minor worth remembering when touching REQ-035's surface: the absent-`write_set` semantics are described with clashing adverbs in three places (board.md "same reading" vs model.go/prime "opposite reading") — unify the framing if you touch them ("same reading, opposite rendering").

## Suggested next prompts

- `do-work run` — continues the queue (REQ-035 first, plan pre-baked).
- `do-work clarify` — answer REQ-039 (background-agents.md worktree mention) and REQ-041 (four board-hardening approvals).
- `do-work cleanup` — after the queue drains, to consolidate UR-007's archive folder.
