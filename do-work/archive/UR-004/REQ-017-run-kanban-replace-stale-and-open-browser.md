---
id: REQ-017
title: "`just run-kanban`: replace a stale board server on the port and open the default browser"
status: completed
created_at: 2026-07-01T21:55:23Z
claimed_at: 2026-07-01T21:58:31Z
completed_at: 2026-07-01T22:19:42Z
route: B
commit: 0d5c1f5
user_request: UR-004
domain: backend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
kb_status: pending
kb_entry:
---

# `just run-kanban`: replace a stale board server on the port and open the default browser

## What

Running `just run-kanban` should end with the dashboard visible: (1) if a **previous queue-kanban instance** is still listening on the target port, terminate it before binding (re-running the recipe replaces the old server instead of failing with "address already in use"); (2) once the server is up, open the user's default browser at the board URL. Today the recipe only builds and serves — the user must notice the printed URL and open it by hand, and a leftover server blocks re-runs.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-do-kanban.md`, `general.md`, `karpathy.md`, `backend.md`. **Pre-flight finding:** on read, `tools/queue-kanban/{main,serve,serve_test,generate,generate_test}.go` and `prime-do-kanban.md` already carried *uncommitted, unrelated* changes (a loopback-only default-bind security hardening — `kanbanServeDefaultListenAddress` → `127.0.0.1:8090`, `defaultRecentWindow` const, `--recent-window` dropped from generate/serve) from a concurrent `claude` session on this same working tree (3 `claude` processes were running; the diff was still growing across two polls, then stabilized for 20s). No REQ file in `do-work/queue|working|archive` matches it — it's out-of-band. Not part of REQ-017's scope; did NOT revert it. Treated current on-disk content as the real baseline and layered REQ-017's change on top of it, since blocking entirely on an environment condition I can't resolve would leave REQ-017 undone. Logged as D-01 (ESCALATE) below — the orchestrator/user should know two uncommitted changesets are interleaved in one working tree.
  **Approach:** (1) `serve.go`: add `--open` bool flag (first Bool flag in module) to the existing FlagSet; extract the bind step into `bindServeListenerAndAnnounce(listenAddress, repoRoot, openAfterBind bool, browserOpener func(string)) (net.Listener, error)` — binds via `net.Listen`, and ONLY on success prints the existing banner/warning block then (if `openAfterBind`) calls `browserOpener`; `runServeCommand` calls it, exits(1) with the raw bind error on failure (no banner ever printed), then calls `httpServer.Serve(listener)` instead of `ListenAndServe()`. The `browserOpener` param is the injectable seam for ordering tests. (2) New `openbrowser.go`: `selectBrowserOpenCommand(goos, targetUrl) *exec.Cmd` (pure, goos as a parameter — the seam for selection tests, since `exec.Command` only builds the Cmd, never launches) switching darwin→`open`, windows→`rundll32 url.dll,FileProtocolHandler`, default→`xdg-open`; `openBrowser(targetUrl)` wraps it with `runtime.GOOS` and calls `.Start()` (fire-and-forget), logging (never fatal) on error. targetUrl is always self-built from the resolved listen address, never user input — no shell involved (mirrors `model.go:468`). (3) `main.go`: append `[--open]` to the serve usage doc line. (4) `serve_test.go`: add bind-then-announce ordering tests (bind failure → opener never called; success + open=true → opener called once with the right URL; success + open=false → opener untouched) using a real ephemeral port probe. New `openbrowser_test.go`: table test over `selectBrowserOpenCommand` per goos (darwin/linux/windows/unlisted-defaults-to-xdg-open), asserting argv without starting anything. (5) `actions/install.md`: recipe gets a new kill-stale line (own shell, no `cd` needed) ahead of the existing build+serve line — `lsof -ti tcp:{{port}} -sTCP:LISTEN` → `ps -p PID -o comm=` basename-compared to `queue-kanban` → kill+wait if match, loud `exit 1` if not (aborts the recipe, per `just`'s stop-on-nonzero-line default), skip gracefully if `lsof` missing; `--open` added to the serve invocation; update deliberate-choices prose (now three), Phase 4 report, workflow intro, add an upgrade note. (6) `actions/board.md`: one-line cross-ref near the standing-shortcut paragraph; Step 5 serve bullet untouched (no `--open`).
- [x] **[APPLY]:** Code written per plan above; scope held to the six files declared in Scope (openbrowser.go added as the optional new file).
- [x] **[UNIFY]:** See Verification section of the builder's return report (git diff --stat reviewed file-by-file, `go build`/`go test`/`gofmt`/`go vet` all clean, functional GREEN+RED serve runs, justfile parse check).

## Why (if provided)

"otherwise I run the command and I don't see the dashboard" — the whole point of the shortcut is to land the user on the board, not to print a URL.

## Context

- Current recipe (shipped block in `actions/install.md`, ~lines 300–314): `cd <kanban-dir> && go build -o queue-kanban . && ./queue-kanban serve --repo-root "{{justfile_directory()}}" --port {{port}}` — no kill step, no open step.
- `tools/queue-kanban/serve.go` has no browser-open logic anywhere (verified by grep); `main.go` owns the flag parsing for `serve`.
- The `just-kanban` installer is append-only and **never replaces an existing recipe** (`actions/install.md` Phase 1: "already installed" if a `run-kanban` recipe is present). Shipping a new recipe block therefore does NOT reach projects that already installed it — the deliverable must include an upgrade note (remove the old recipe block, re-run `do-work install just-kanban`).
- The recipe runs on the user's machine interactively; `do-work board` (`actions/board.md`) is typically driven by an agent, where auto-opening a browser is wrong. Keep `do-work board` behavior unchanged (out of scope) — the open behavior must be opt-in at the tool level (flag), with the just recipe being the caller that opts in.
- Cross-platform note: the skill designs for the floor. "the open command" is macOS; Linux is `xdg-open`, Windows `start`/`rundll32`. A platform-aware open inside the Go tool (`--open` flag on `serve`, default off) keeps the recipe one line and portable; the user's own justfile remains project-owned and editable.
- Port-kill safety: "kill the previous one" means the previous **queue-kanban** instance. If a non-queue-kanban process occupies the port, do NOT kill it — fail fast with a clear message naming the squatter. Killing arbitrary port listeners from a shipped recipe is a foot-gun.
- Recipe recipes must keep each `cd … && …` chain on one logical line (`just` runs each line in a fresh shell — see `actions/install.md` "Two deliberate choices").

## Builder Guidance

Firm on the outcome (re-run replaces stale board server; browser opens on the board URL; agents/CI never get a surprise browser). Latitude on mechanics, with a recommended shape:

- **Recommended:** add `--open` to `queue-kanban serve` (default off) that launches the platform opener *after a successful bind* (bind the listener first, then open, then serve — avoids the sleep-race a shell workaround would need). Stale-instance replacement can live either in the tool (detect bind failure → identify listener → kill only if its command is `queue-kanban` → retry once) or in the recipe (a guarded one-liner that kills only matching processes); tool-side is sturdier, recipe-side is simpler — builder's call, but the non-queue-kanban-squatter case must fail loudly either way.
- Update the shipped recipe block in `actions/install.md` (pass `--open`, add the replace-stale behavior) and any docs that show the recipe or serve usage (`actions/board.md` mentions serve mode; `docs/` if the recipe appears there; `README.md` if it does).
- Keep the parser/schema untouched — this is serve/recipe plumbing, no Schema Read Contract impact.
- Mirror the change with tests where runnable (flag parsing, opener-command selection per platform via an injectable seam, bind-then-open ordering); the end-to-end browser launch itself is manual-proof territory (hence `tdd: false` — but write the unit tests the seam makes possible).

## Red-Green Proof
**RED prompt/case:** With a queue-kanban instance already serving the target port, a second `just run-kanban` fails to bind ("address already in use") instead of replacing it; and on any run no browser opens — `grep -rn 'open\|xdg-open' tools/queue-kanban/serve.go main.go` shows no opener, and the shipped recipe block in `actions/install.md` contains no kill or open step.
**Why RED now:** The recipe is build + serve only; the tool has no `--open` capability; a stale server on the port makes the recipe error out.
**GREEN when:** (a) `queue-kanban serve --open` opens the default browser at the board URL after a successful bind (macOS `open` / Linux `xdg-open` / Windows equivalent), and plain `serve` without the flag behaves exactly as today; (b) the shipped `run-kanban` recipe passes `--open` and replaces a stale queue-kanban listener on the target port, so running `just run-kanban` twice in a row succeeds with the second run serving; (c) a non-queue-kanban process on the port is left alive and the run fails with a clear message; (d) `go build` + `go test ./...` pass in `tools/queue-kanban/`; (e) `actions/install.md`'s recipe block is updated and an upgrade note for already-installed projects (remove old recipe, re-run `do-work install just-kanban`) ships with the change.
**Validation:** User confirmed (2026-07-01 — capture summary, including the kill-scope narrowing to "previous queue-kanban instance only" and the recommended `--open` flag shape, was reviewed by the user, who proceeded to verify + run without adjustments).

---
*Source: UR-004 — "just run-kanban should not only start the server (kill the previous one if it is on the same port) but also use the open command to open the default browser with the target kanban, ohterwise I run the command and I don't see the dashboard"*

Think carefully before answering.

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear and user-validated, and the primary sites are named (serve.go, main.go, the recipe block in actions/install.md), but the mechanics carry real latitude (where the replace-stale logic lives, how the opener seam is structured for bind-then-open ordering and testability, which docs show the recipe). Exploration will pin down the serve startup path, flag parsing, and existing test patterns before scope is declared.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Serve startup path:** `main.go:36-37` dispatches to `runServeCommand` (`serve.go:254`); flags parsed via fresh FlagSet (`serve.go:255-261`, no boolean flags in the module yet); port resolution is flag > `QUEUE_KANBAN_PORT` env > `:8090` default (`serve.go:228-247`). **Startup banner prints BEFORE the bind** (`serve.go:272-277`), then `httpServer.ListenAndServe()` does bind+serve in one call (`serve.go:301`) — today a port collision prints a false-positive "live board at http://…" line followed by `log.Fatalf` with `bind: address already in use`. Bind-then-open therefore requires splitting into `net.Listen` + `httpServer.Serve(listener)` and moving the banner (and the new open call) after a successful listen; `Serve` runs alongside the existing SIGINT/SIGTERM shutdown goroutine (`serve.go:287-299`), which a fire-and-forget opener does not interact with.

**Seams/testability:** `serve_test.go` exercises only the `http.Handler` via `httptest` (`serve_test.go:49-126`) — no tests for `runServeCommand`, flag parsing, or bind behavior. The module's one `exec.Command` precedent (`model.go:468`, git lookup) validates untrusted input before argv and runs best-effort. No `runtime.GOOS` usage anywhere — the platform-aware opener introduces that pattern; needs an injectable opener seam to be testable.

**Recipe/doc surface:** primary edit is the shipped recipe block `actions/install.md:300-314` (recipe at 304-305), plus the "two deliberate choices" prose (316-319), the Phase 4 report template (343-349), and the workflow intro (274). `actions/board.md:75` (agent-driven serve) must NOT gain `--open`; `board.md:79` standing-shortcut cross-ref may take a one-line mention. `docs/` has zero recipe references; README/SKILL.md rows describe the agent action and are unaffected. Idempotency greps (`install.md:284-285`, `:329`) key off the recipe name only — unaffected by body changes.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `tools/queue-kanban/serve.go` (modify) — bind-then-serve split (`net.Listen` + `Serve`), banner moved after successful bind, `--open` flag, opener invocation
- `tools/queue-kanban/openbrowser.go` (new, optional — builder may inline in serve.go) — platform-aware `openBrowser(url)` with injectable exec seam
- `tools/queue-kanban/serve_test.go` (modify) — tests for flag parsing, opener selection per GOOS, bind-then-open ordering (opener not called on bind failure)
- `tools/queue-kanban/main.go` (modify) — serve usage doc comment gains `[--open]`
- `actions/install.md` (modify) — recipe block (replace-stale line + `--open`), deliberate-choices prose, Phase 4 report template, workflow intro, upgrade note for already-installed projects
- `actions/board.md` (modify) — one-line cross-ref that the standing shortcut auto-opens while the agent action does not (Step 5 serve bullet stays `--open`-free)

**Files I will NOT touch:** `tools/queue-kanban/model.go` / `generate.go` / `walk.go` / `web/*` (no schema or frontend impact); `README.md` / `SKILL.md` (describe the agent action, unaffected); `docs/` (no recipe references exist). `CHANGELOG.md` + `actions/version.md` are orchestrator-owned commit-time changes, outside builder scope.

**Acceptance criteria (restated from REQ):**
- [x] `queue-kanban serve --open` opens the default browser at the board URL only after a successful bind; plain `serve` behaves exactly as today (no open)
- [x] Shipped `run-kanban` recipe replaces a stale queue-kanban listener on the target port (second `just run-kanban` succeeds) and passes `--open`
- [x] A non-queue-kanban process on the port is left alive; the run fails with a clear message
- [x] `go build` + `go test ./...` pass in `tools/queue-kanban/`
- [x] `actions/install.md` recipe block updated + upgrade note for already-installed projects

## Decisions

- **D-01 [ESCALATE]** — On first read, `tools/queue-kanban/{main,serve,serve_test,generate,generate_test}.go` and `prime-do-kanban.md` already carried *uncommitted* changes from a live, concurrent `claude` session on this same working tree (3 `claude` CLI processes were running on this machine; the diff grew across two ~10s polls before stabilizing). No REQ file in `do-work/queue|working|archive` tracked it. I did not revert or "fix" any of it — treated it as the real on-disk baseline and layered REQ-017's change on top via exact-match `Edit` calls (which fail loudly rather than silently corrupt on a content mismatch). Partway through my implementation the other session committed it as `a2cde95 fix: 0.103.3 The Loopback Latch` — no collision occurred, but the risk was real while it lasted. Value: n/a (process observation, not a design choice). Risk: had either session used a full-file overwrite instead of match-based edits, or had both touched the identical line concurrently, one side's work could have been silently lost; reversible via git if caught, but nothing was watching for it. Flagging so the orchestrator/user knows two independent uncommitted changesets were interleaved in one working tree during this run.
- **D-02 [DECIDE & STATE]** — Ran `gofmt -w` on `serve_test.go`'s pre-existing comment-column misalignment inside `TestResolveServeListenAddressBindsLoopbackByDefault` (introduced by the concurrent session's commit, not by REQ-017). Mechanical, zero semantic change; done because the verification bar requires `gofmt -l .` clean on a file I was already editing for my own new tests.
- **D-03 [DECIDE & STATE]** — Implemented the replace-stale-instance logic in the recipe (`actions/install.md`), not the Go tool, per the REQ's Builder Guidance ("builder's call... recipe-side is simpler"). Keeps `serve.go` scoped to bind/announce/open; the kill-scope verification (`ps` command-name match) lives in shell, consistent with how `install.md` already writes its own guarded one-liners.
- **D-04 [DECIDE & STATE]** — Used parameter-injected seams — `bindServeListenerAndAnnounce(..., browserOpener func(string))` and `selectBrowserOpenCommand(goos string, targetUrl string) *exec.Cmd` (goos as an argument, not read from `runtime.GOOS` internally) — rather than a mutable package-level `var browserOpenCommand func(...)` (the REQ's own illustrative example). Both are equally injectable for tests; the parameter form avoids shared mutable package state and needs no save/restore dance in tests.
- **D-05 [DECIDE & STATE]** — The recipe's kill step sends `SIGTERM` (`kill "$PID"`, not `-9`) and then polls `kill -0` for up to 2s before proceeding, so the killed server's graceful-shutdown path runs and the port is confirmed free before the build+serve line executes. Verified empirically (see builder report Verification section) — without the wait, a fast re-bind could theoretically race the old listener's release.

## Discovered Tasks

- [normal] `tools/queue-kanban/prime-do-kanban.md`'s `## Lessons` section may be worth a line once this REQ archives, alongside the already-committed loopback-binding change (`a2cde95`) — not touched here (out of scope for REQ-017; the prime file wasn't in my declared Scope). *Resolved at archival: the REQ-017 lesson line was written (inline, per the no-archive-links convention); the loopback change needed no Lessons line — its `Traps` entry ("Loopback-only by default") already carries it.*
- [low] Three concurrent `claude` CLI processes were observed running against this same working directory (not isolated via git worktrees) during this REQ's execution — outside my control/scope to fix, but flagging it as a near-miss: see D-01.

## Lessons Learned

**What worked**
- Parameter-injected seams — `browserOpener func(string)` and `goos` as an argument rather than reading `runtime.GOOS` inside — made every selection branch and ordering case testable without launching a browser or doing a save/restore dance on package state (D-04).
- Splitting `ListenAndServe` into `net.Listen` + `Serve(listener)` fixed a latent false-positive discovered during exploration: the old code printed the "live board at …" banner *before* binding, so a port collision announced a server that never came up.

**What didn't**
- Nothing failed hard, but D-01 was a real near-miss: two independent uncommitted changesets from concurrent sessions interleaved in one working tree. Exact-match edits (fail-loud on mismatch) were the only guard; the queue-over-worktrees discipline exists precisely to avoid this.

**Worth knowing**
- The just-kanban installer is append-only — recipe behavior changes never reach already-installed projects; every recipe change must ship an upgrade note (delete block, re-run install).
- A recipe kill step needs SIGTERM + a `kill -0` poll before rebinding: a fast re-bind can race the old listener's port release (D-05).

## Orientation

Now `just run-kanban` lands the user on the board instead of printing a URL: `queue-kanban serve` gained an opt-in `--open` (platform-aware browser open, fired only after a successful bind — new `openbrowser.go` beside the serve path), and the shipped justfile recipe replaces a stale queue-kanban instance on the port before serving. Agent-driven `do-work board` deliberately keeps auto-open off. Lives in the queue-kanban tool + the just-kanban installer recipe (`actions/install.md`); no schema or frontend impact.

## Implementation Summary

*(Appended post-archive by the orchestrating session: the REQ was archived by a concurrent session without this mandatory section — see Qualification note.)*

**Files changed:**
- `tools/queue-kanban/serve.go` (modified) — bind-then-announce split: `bindServeListenerAndAnnounce(...)` binds via `net.Listen` first, only then prints the banner and (with `--open`) fires the browser opener; `runServeCommand` uses `httpServer.Serve(listener)`; `--open` bool flag added
- `tools/queue-kanban/openbrowser.go` (new) — `selectBrowserOpenCommand(goos, url)` (darwin→`open`, windows→`rundll32 url.dll,FileProtocolHandler`, default→`xdg-open`; goos parameter is the test seam) + fire-and-forget `openBrowser` (`.Start()`, warning-only on error)
- `tools/queue-kanban/openbrowser_test.go` (new) — per-GOOS opener-selection table test + pinned windows argv shape
- `tools/queue-kanban/serve_test.go` (modified) — 3 ordering tests: bind failure → opener never invoked; success + `--open` → opener fires exactly once; success without `--open` → opener untouched
- `tools/queue-kanban/main.go` (modified) — serve usage doc comment gains `[--open]`
- `actions/install.md` (modified) — recipe block gains the replace-stale line (kill only a verified `queue-kanban` process; refuse + exit 1 for foreign squatters; graceful without `lsof`) + `--open`; upgrade note for already-installed projects; deliberate-choices prose + Phase 4 report template updated
- `actions/board.md` (modified) — one-line cross-ref: the standing shortcut auto-opens, the agent-driven `do-work board` does not (Step 5 stays `--open`-free)

**What was done:** `queue-kanban serve` now binds before announcing (fixing the false-positive "live board at…" banner on port collisions) and gained an opt-in `--open` flag with a platform-aware, test-seamed browser opener; the shipped `just run-kanban` recipe replaces a stale board instance on the target port and passes `--open`, so the command ends with the dashboard visible.

## Qualification

Passed with a **process flag**: implementation verified against commit `0d5c1f5` (all 7 project files present in the commit, diff matches the Scope declaration; the optional `openbrowser*.go` split was declared). Requirements traced: bind-then-open ordering (serve.go), platform opener seam (openbrowser.go), replace-stale + foreign-squatter refusal (recipe block), upgrade note (install.md), board.md exclusion honored. P-A-U boxes checked with substantive notes.

⚠ **Process anomaly (D-01, ESCALATE):** a concurrent Claude session on the same working tree committed this REQ's implementation (`0d5c1f5`), its own unrelated loopback-hardening fix (`a2cde95`, 0.103.3), and the prime-lessons archival step (`3f74e15`, 0.104.1), and archived this REQ **without** the Implementation Summary/Qualification/Testing/Review sections or the `commit:` writeback. This session repaired the paper trail post-archive. No code collision occurred; the multi-orchestrator hazard is escalated in the session hand-back.

## Testing

**Tests run:** `go test -count=1 ./...` in `tools/queue-kanban/` (re-run by the orchestrator at HEAD `3f74e15` post-commit)
**Result:** ✓ All passing (`ok github.com/knews2019/skill-do-work/queue-kanban`); `gofmt -l .` and `go vet ./...` clean at HEAD

Red-green validation omitted per `tdd: false` — regression + functional evidence traces to the REQ's `## Red-Green Proof`:
- RED reproduced: second serve on an occupied port exited 1 with `bind: address already in use` and **no false-positive banner** (pre-fix the banner printed before the bind)
- GREEN verified by the builder: banner-after-bind + HTTP 200 board render on `:8123`; kill-stale logic exercised against 4 scenarios (queue-kanban listener killed, foreign python listener refused + left alive, empty port no-op, missing `lsof` graceful skip); recipe parsed by `just --list`/`--dump` with tokens intact; one manual `--open` smoke run fired the macOS opener without error
- New tests: 3 bind/open ordering tests + per-GOOS opener selection tests (browser never launched in tests — seam only)

## Review

**Approve** — clean, well-tested fix that delivers exactly what UR-004 asked for; verified end-to-end against HEAD (3f74e15), including the interaction with the concurrently-landed loopback-binding change (a2cde95).
Route B | commit 0d5c1f5 | reviewed post-commit (repair pass — see Qualification)

### What's built
- `serve` binds before announcing — the false-positive "live board at …" banner on port collision is gone (verified live: second serve prints only the raw bind error, exits 1, no banner).
- Opt-in `--open` launches the platform opener only after a successful bind (darwin `open` / linux `xdg-open` / windows `rundll32`) via an injectable seam, fully unit-tested without spawning processes.
- The shipped `just run-kanban` recipe kills a stale **queue-kanban** listener (SIGTERM + poll), refuses loudly and leaves alive any foreign process on the port, degrades gracefully without `lsof`, and passes `--open`.

### Findings

**Important:** None.

**Minor:**
- `openbrowser_test.go` is a new file not itemized in `## Scope`'s file list (it was anticipated in the [PLAN] text) — a bookkeeping gap between REQ sections, not undisclosed drift.
- D-02 (gofmt realignment of a pre-existing misalignment introduced by the concurrent a2cde95) — mechanical, zero-semantic-change, self-disclosed.

**Nit:**
- `actions/install.md` Phase 4 + `actions/board.md` Step 5 still say `http://localhost:8090` while the post-a2cde95 banner prints `http://127.0.0.1:8090`. Functionally harmless; predates REQ-017; low-priority doc cleanup.

### Requirements Checklist

- [x] `--open` opens the browser only after a successful bind; plain `serve` unchanged — verified (ordering unit tests + manual bind-collision run)
- [x] Recipe replaces a stale queue-kanban listener and passes `--open` — verified via built binary + lsof/ps kill-scope test
- [x] Foreign process on the port left alive; run fails with a clear message — verified with a live `python3 -m http.server` listener
- [x] `go build` + `go test ./...` pass — re-verified at HEAD; gofmt + vet clean
- [x] `actions/install.md` recipe block updated + upgrade note for already-installed projects — present (Phase 1)

### Acceptance Testing

**Result: Pass** — full suite + vet + gofmt clean at HEAD; live bind/collision run confirms the RED→GREEN proof; kill-stale verified against real queue-kanban and foreign listeners; the installed root `justfile` parses and matches the shipped block exactly. Kill logic keys off port number only, so it is unaffected by a2cde95's loopback-only binding (confirmed empirically). `--open` not fired during review (no browser popping); path substantiated by seam unit tests + the builder's macOS smoke run.

### Suggested Additional Testing
- Manual `xdg-open`/`rundll32` verification on Linux/Windows (argv shapes are unit-pinned; only macOS smoke-tested live).
- One real back-to-back `just run-kanban` double-run to exercise the full recipe as a unit.
- Doc cleanup: reconcile `localhost:8090` vs `127.0.0.1:8090` references (from a2cde95).

### Scores

**Overall: 94%** — Requirements 100 / Code Quality 95 / Test Adequacy 90 / Scope 90 / Risk: none / Acceptance: Pass.

### Follow-up REQs Created
None — all findings Minor/Nit.

*Generated by review-work agent (post-commit repair pass)*
