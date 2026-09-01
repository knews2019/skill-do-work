---
id: REQ-094
title: Checkpoint writer label — crash recovery ignores foreign entries
status: completed
created_at: 2026-08-04T19:44:17Z
claimed_at: 2026-08-04T19:48:36Z
completed_at: 2026-08-04T20:08:59Z
commit: 9c305c0
route: B
kb_status: promoted
kb_entry: REQ-094-checkpoint-writer-label-crash-recovery-i.md
user_request: UR-018
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-095, REQ-096, REQ-097, REQ-098, REQ-099, REQ-100, REQ-101]
batch: parallel-building
write_set: [actions/work-reference.md, actions/work.md, actions/cleanup.md, actions/forensics.md, docs/work-guide.md, tools/queue-kanban/verify.go, tools/queue-kanban/verify_test.go, _dev/tests/contract-regressions.sh, do-work/CHECKPOINT.md]
---

# Checkpoint Writer Label — Crash Recovery Ignores Foreign Entries

## What

Give `do-work/CHECKPOINT.md` In-Progress entries a **static writer label** identifying the checkout that wrote them, and scope crash recovery to own-label entries only. Foreign entries are **reported, never stripped**. This defuses a live landmine: the checkpoint is git-tracked, and once two checkouts sync it, checkout A reads checkout B's live claim as its own crash, strips it, and re-queues a REQ someone is actively building — a deterministic replay of the 2026-07-01 collision, no race needed.

## Detailed Requirements

- In `actions/work-reference.md`'s In-Progress Record section (~lines 409–425): each In-Progress entry records `writer: <hostname>:<absolute-checkout-path>` (path alone collides across machines — both sides can be `/home/user/repo`).
- Crash-recovery rule: only entries whose writer label matches **this** checkout are crash-recovery candidates. Foreign entries are listed in the recovery report and left untouched — extend the existing foreign-claim rule (`actions/work-reference.md:239` — "a claim you didn't record is not yours, never touch it") to cover the checkpoint.
- Reword the tripwire at `actions/work-reference.md:413`: refresh intervals, staleness checks, and liveness probes stay **banned by name**; a static writer label is explicitly not liveness machinery (it is written once, never refreshed, never checked for staleness).
- Entries written before this change have no writer label: treat a label-less entry as **own** on the checkout that has it uncommitted/locally modified, otherwise report-only — never guess-strip.
- `CHECKPOINT.md` stays tracked and committed (user decision: "checkpoints are transient, it's fine to commit them before changing, this way different versions of it already is available in the git history").

## Constraints

- Do NOT add: heartbeats, refresh intervals, holder-liveness checks, staleness thresholds, auto-takeover. (Batch-wide do-not-build list — see UR-018.)
- This is a prose/contract change to shipped instruction files; keep wording surgical and in the existing section's voice.

## Red-Green Proof

**RED prompt/case:** Simulate a synced foreign checkpoint: write a CHECKPOINT.md In-Progress entry for a REQ that is in `working/` but was claimed "by another checkout", then follow today's crash-recovery instructions — they classify it as a local crash and strip/re-queue it.
**Why RED now:** In-Progress entries carry no writer identity (deliberately, per the current `:413` tripwire), so recovery cannot distinguish a foreign live claim from an own crash.
**GREEN when:** Following the updated instructions, the same foreign entry is reported ("claim held by <other-writer>, not touched") and the REQ stays claimed; an own-label entry still recovers exactly as before.
**Validation:** User confirmed (ask-tool answer: "Static writer label").

## Full Context

See `do-work/user-requests/UR-018/input.md` and `do-work/user-requests/UR-018/assets/approved-plan.md` (Phase 1).

---
*Source: approved plan `let-s-talk-about-this-kind-robin.md`, Phase 1 item 1*

---

## Triage

**Route: B** - Medium

**Reasoning:** The target file and sections are named, but the change redefines a contract token (the no-holder-id tripwire) that other shipped text restates — exploration must find every restatement site (crash-recovery prose, Step 2 claim notes, checkpoint templates, contract tests) before editing, or the edit ships a stale-echo defect.

**Planning:** Not required

## Plan

Planning not required for Route B — proceeding with exploration.

## Exploration

Edit-site inventory (Explore agent, full sweep of shipped files + parsers + tests):

**Canonical sites (work-reference.md):** `:417` entry format string (gains the writer field); `:413` tripwire (pinned by contract test on `never grow into one` — keep the phrase); `:411` "named there" → "named there under this checkout's writer label"; `:418` "one entry per REQ id" → per id *per writer*, refresh scoped to own-label; `:420` removal scoped to own-label; `:796-802` Step 10 checkpoint template must carry the label; `:817` Step 10 "enriches and keeps the list" needs an explicit preserve-foreign-entries clause (wholesale rewrite is a label-destruction path).

**Crash Recovery (work-reference.md:232-273):** `:238` own-crash bullet → own-label; `:239` foreign-claim bullet gains a **third case**: entry present with foreign label = positive evidence of another checkout's live claim — report `claim held by <writer>, not touched`, never enters the 3-hour takeover ladder (that ladder stays for label-less/unaccounted claims); a **fourth case**: legacy label-less entry = own iff the checkpoint is locally modified/uncommitted here, else report-only. Keep `absent checkpoint is ambiguous` verbatim (contract-test-pinned). `:241` + `:413` + work.md:230 — all three copies of "no second owner reads it" become literally false; reword consistently (a second checkout reads it only to classify-and-skip; nothing acquires or waits on it).

**work.md:** `:118`, `:120` condensed restatements; `:230` Step 2 claim write gains label; `:555` Step 8 removal own-label-scoped; `:642-647` session-start — `:645` whole-file delete is dangerous: delete own entries only, whole-file delete only when no foreign entries remain; `:656-672` checklist lines (light).

**Other actions:** `cleanup.md:40` (Pass 0 drop → own-label); `forensics.md:39` (Check 1 → own-label); `forensics.md:200` ghost-id caveat (foreign entries can name REQs archived elsewhere).

**Parsers/tests:** `tools/queue-kanban/verify.go:250` remedy string "edit or delete the checkpoint" now recommends the destructive act — reword; id-regex parsing unaffected by a `writer` suffix. `_dev/tests/contract-regressions.sh:428-562` pins exact phrases (`never grow into one`, `absent checkpoint is ambiguous`, `foreign claim`, `Crash Recovery's input`, line-order, `one queue owner per checkout` exactly once) — preserve all; add new assertions pinning the writer label + foreign-label branch.

**Docs:** `docs/work-guide.md:66-68` light updates (claim-time write + own-classification sentence).

**Live state:** `do-work/CHECKPOINT.md:5` (this REQ's own entry) is the legacy label-less shape — gains the label when the change lands.

*Generated by Explore agent*

## Implementation Summary

**What was done:** In-Progress checkpoint entries now carry `— writer: <hostname>:<absolute-checkout-path>`; the canonical contract (In-Progress Record) defines the format and derivation, append/refresh/removal are scoped per-id-per-writer to own-label entries, and Crash Recovery classifies four cases (own-label / foreign-label / label-less legacy / unnamed-or-absent). Foreign-label entries report `claim held by <writer>, not touched` and never enter the 3-hour takeover ladder. The tripwire keeps banning refresh intervals, staleness checks, and liveness probes by name (dropping "holder id" from the ban, carving out the static label); `never grow into one` survives byte-identical. Both label-destruction paths (Step 10 wholesale rewrite, session-start delete) are own-label-scoped with foreign entries preserved verbatim. Three new contract-suite assertions pin the label, the foreign-label report phrase, and the tripwire.

**Files changed:**
- `actions/work-reference.md` (modified) — In-Progress Record entry format + derivation, tripwire reword, Crash Recovery four-case classification, ladder preamble, Step 10 template + preserve-foreign clause
- `actions/work.md` (modified) — Steps 1/2/8/10 echoes, session-start own-label-scoped delete, checklist lines
- `actions/cleanup.md` (modified) — Pass 0 entry-drop scoped to own-label
- `actions/forensics.md` (modified) — Check 1 remediation scoping, ghost-id foreign-entry caveat
- `docs/work-guide.md` (modified) — checkpoint section wording
- `tools/queue-kanban/verify.go` (modified) — ghost-finding remedy string no longer recommends whole-file deletion
- `_dev/tests/contract-regressions.sh` (modified) — +3 assertions (writer label, `claim held by`, tripwire phrase)
- `do-work/CHECKPOINT.md` (modified) — live REQ-094 entry relabeled with this checkout's writer label

*Written by orchestrator from builder manifest*

## Review

**Overall: 90%** | 2026-08-04

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 88% |
| Test Adequacy | 80% |
| Scope | 97% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 1 important, 3 minor. Important: `actions/work.md` Step 10's two preserve rules are scoped to "another checkout's `writer:` label", which excludes the label-less report-only case — a label-less entry in a clean checkpoint survives Crash Recovery's report-only branch but is then deleted by the session-start whole-file delete, re-entering the takeover ladder next run. Fix: scope both to "every entry this checkout did not write" (matches `actions/work-reference.md`'s canonical Step 10 clause) + pin both label-destruction paths with a contract assertion. → follow-up REQ-102. Minor: `actions/work-reference.md:55` / `docs/work-guide.md:91` "does not detect… a second owner" now partly false → folded into queued REQ-096 (owns those lines); "classifies by name" self-corrects in-paragraph; verify.go remedy string unasserted.
**Acceptance:** Pass — contract suite exit 0, go test ok, live writer label matches `hostname -s`:`git rev-parse --show-toplevel`, 3 new assertions proven RED before / GREEN after against HEAD.
**Follow-ups created:** REQ-102 (Step 10 preserve-rule scoping + destruction-path assertions); REQ-096 addendum (Execution Model restatement).

*Reviewed by review-work action*

## Lessons Learned

**What worked:** A pre-build exploration inventory of every restatement site (3 copies of "no second owner reads it", 2 non-obvious label-destruction paths in Step 10's wholesale rewrite and the session-start delete, 5 pinned contract-test phrases) — the build touched 8 files with zero suite breakage because the collision surface was mapped first. Pinning the tripwire ban and its carve-out to the *same paragraph* via a new assertion (a carve-out that drifts into its own paragraph reads as general permission).

**What didn't:** The builder's Step 10 echoes narrowed "every entry this checkout did not write" to "another checkout's label" — an echo written from memory of the canonical clause, not from it. Echo sites should quote the canonical condition, not paraphrase it (that's REQ-102).

**Worth knowing:** The checkpoint travels between checkouts on any install that commits `do-work/` — every rule about it now has four claim-origin cases (own-label / foreign-label / label-less / unnamed), and the three-hour takeover ladder serves only the last two. `checkpointMentionedRequestIds` in `tools/queue-kanban` extracts ids by regex, so entry-format suffixes are parser-transparent.

## Orientation

Crash recovery in the work pipeline now identifies *which checkout* wrote each in-progress checkpoint entry (static `writer:` label) and only ever strips its own — another checkout's synced live claim is reported and left alone. [MAP CHANGED] — first piece of the UR-018 multi-checkout contract: the checkpoint file is now shared state with per-checkout ownership semantics, groundwork for claim-anywhere/one-releaser (REQ-096).

Passed — 8 files verified in the working diff (mechanical check OK), all 5 REQ requirements traced (label format + derivation, own-label recovery scoping, foreign-label report line outside the takeover ladder, tripwire reword with pinned phrases intact, legacy label-less handling), P-A-U confirmed. Judgment checks: hunks are substantive and voice-matched; no hollow edits. `one queue owner per checkout` still appears exactly once across `actions/`.

## Testing

- `bash _dev/tests/contract-regressions.sh` → exit 0 (orchestrator-run, independent of builder's evidence). The 3 new assertions were verified non-vacuous by the builder via manual sed-block extraction (1 hit writer label, 2 hits `claim held by`, tripwire + carve-out on one line).
- `cd tools/queue-kanban && go test ./...` → ok (orchestrator-run).
- Red-green validation (describable proof, `tdd: false`): RED — pre-change prose classified any checkpoint-named entry as own crash (strip/re-queue); GREEN — post-change Crash Recovery classifies a foreign-label entry as `claim held by <writer>, not touched`, never enters the takeover ladder, and the new suite assertion pins that phrase. Baseline: suite was green pre-change (Pre-Flight), green post-change with 3 additional assertions.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** `prime_files` is empty; loaded `crew-members/general.md` and `coding-guardrails.md`.
  Approach: state the contract **once** at its canonical site — `actions/work-reference.md` →
  **In-Progress Record (Step 2)** — where the entry format gains `— writer: <hostname>:<absolute-checkout-path>`,
  the label's derivation is defined, and the tripwire is reworded to keep banning refresh intervals,
  staleness checks and liveness probes *by name* while explicitly carving the static label out
  (`never grow into one` stays byte-identical). Crash Recovery's classification bullet list then grows
  from two cases to four (own-label / foreign-label / label-less legacy / unnamed-or-no-checkpoint),
  with the foreign-label case reported as `claim held by <writer>, not touched` and routed *around* the
  three-hour ladder rather than into it. Every other site is a pointer or a light echo: `actions/work.md`
  Steps 1/2/8/10 + checklist, `actions/cleanup.md` Pass 0, `actions/forensics.md` Check 1 and the verify
  table, `docs/work-guide.md`. Two rewrite paths are the real hazards and both become own-label-scoped:
  Step 10's wholesale checkpoint rewrite (must copy foreign entries through verbatim) and Step 10's
  session-start delete (own entries only; whole-file delete only when no foreign entry remains).
  `tools/queue-kanban/verify.go` changes one remedy string — no parser change, since a trailing
  `writer:` suffix does not affect id extraction. Three new assertions pin the label, the foreign-label
  report phrase, and the tripwire's ban-plus-carve-out living in one paragraph.
- [x] **[APPLY]:** Eight files, all declared in `## Scope`, no others. `tools/queue-kanban/verify_test.go`
  was in Scope conditionally and stayed untouched: the ghost-finding test asserts only `Detail`
  (`verify_test.go:246-257`), never `Remedy`, so the reworded remedy string is unasserted. No parser,
  schema, or logic change anywhere — `writer:` is a trailing suffix on a free-text checkpoint line and
  `checkpointMentionedRequestIds` extracts ids by regex, so id parsing is untouched by construction.
- [x] **[UNIFY]:** `git diff --stat` → 6 implementation files (`actions/work-reference.md` 30±,
  `actions/work.md` 18±, `_dev/tests/contract-regressions.sh` +25, `actions/forensics.md` 4±,
  `docs/work-guide.md` 4±, `actions/cleanup.md` 2±, `tools/queue-kanban/verify.go` 2±) plus
  `do-work/CHECKPOINT.md` and this REQ. Verified per file: **`work-reference.md`** — read all seven
  hunks; confirmed `never grow into one`, `absent checkpoint is ambiguous`, and `foreign claim` are
  byte-identical, and that the four classification bullets sit inside the existing
  `crash_recovery_block` sed range so no existing assertion's haystack shifted. **`work.md`** — read
  all seven hunks; `Crash Recovery's input` intact at :118 and still ahead of the `**Crash Recovery:**`
  paragraph (the suite's line-order probe passes). **`forensics.md`** — table row still has exactly two
  cells (no `|` introduced). **`contract-regressions.sh`** — `bash -n` parses; the three new checks
  proved non-vacuous by extracting both blocks by hand and matching (1/2/1 hits). **`verify.go`** —
  `go build`, `go vet`, `go test ./...` all clean. `one queue owner per checkout` still occurs exactly
  once across `actions/`. No debug artifacts in added lines (`console.log`/`debugger`/`TODO`/`FIXME`
  grep over the `+` side: none).

## Decisions

- **D-01 — The report line is emitted at classification, not from the takeover ladder.** The ladder's
  first act is computing an age from `claimed_at`, which is the wrong question for an attributed
  claim: knowing *who* holds it makes *how old* it is irrelevant, and any ladder entry point invites a
  later edit to reconnect the 3-hour branch. So `claim held by <writer>, not touched` is emitted by the
  foreign-label bullet itself, and *Reporting and takeover* opens by naming the three cases it does
  handle (unnamed, no-checkpoint, label-less-report-only) and explicitly excluding the fourth.
- **D-02 — Legacy label-less entries are classified by `git status` on the checkpoint, not by
  guessing.** The REQ specifies own-iff-locally-modified. Stated as evidence rather than as a rule
  ("this checkout wrote it and has not shared it"), because the reasoning is what stops a maintainer
  from later "simplifying" it into treating any label-less entry as own — which is exactly the
  pre-REQ-094 behavior. The bullet ends by naming that failure: never treat "I cannot tell" as "mine."
- **D-03 — Step 10's wholesale rewrite gets the preserve-foreign rule stated twice, once as rationale
  and once as an instruction.** Normally an echo would be a pointer only, but Step 10 is the single
  place in the pipeline that rewrites the whole checkpoint, which makes it the single place a label can
  be destroyed — and `actions/work.md` Step 10 is where an agent actually performs the write, while
  `actions/work-reference.md`'s template paragraph is where the reason lives. The `work.md` echo is one
  sentence and defers to the reference for the why.
- **D-04 — The tripwire's ban list drops "a holder id" and gains "a liveness probe."** A writer label
  *is* a holder id, so the old list would have banned the thing this REQ adds. What the list is
  actually protecting against is machinery that has to be kept current — refresh intervals, staleness
  checks, liveness probes — so the ban is restated in those terms and the label is carved out by the
  property that distinguishes it: written once, never refreshed, never read as evidence anything is
  alive. The new suite assertion pins the ban and the carve-out to the same paragraph so the exception
  cannot drift into reading as a general permission.
- **D-05 — `verify`'s ghost remedy routes to editing one entry out, not to deleting the file.** The old
  remedy ("edit or delete the checkpoint") is now a data-loss instruction: the file may hold another
  checkout's live claims, and a ghost id is now an *expected* reading of a foreign entry whose REQ was
  archived elsewhere. `actions/forensics.md`'s verify table carries the same caveat so the reader sees
  it before acting. Remedy text only — the check itself still fires identically, per the REQ's
  no-parser-change constraint.
- **D-06 — Both `writer:` halves are justified in prose at the canonical site.** Not decoration: the
  next maintenance pass looking to shorten the entry format will drop the hostname as redundant unless
  the file says why it isn't (two machines can both hold `/home/user/repo`). The same sentence carries
  the collision in the other direction, so neither half can be dropped by that argument.

## Discovered Tasks

- [low] The checkpoint's **frontmatter** (`session_ended`, `last_completed`, `queue_state`,
  `reqs_processed_this_session`) carries no writer identity, so it stays ambiguous after this REQ. On
  a synced checkout, `actions/work.md` Step 10 → *On session start* reports "Resuming from previous
  session. Last completed: REQ-NNN" from values another checkout wrote. Harmless today — it is a
  resume banner, and nothing classifies or strips anything on those fields — but it is the same
  ambiguity the `writer:` label just closed for the In Progress list, one section up in the same file.
  Out of scope here: this REQ's contract is scoped to the In-Progress entries recovery consumes.

## Scope

**Files I will touch:**
- `actions/work-reference.md` (In-Progress Record, Crash Recovery, Session Checkpoint Template)
- `actions/work.md` (Step 1/2/8/10, checklist)
- `actions/cleanup.md` (Pass 0 removal scoping — one line)
- `actions/forensics.md` (Check 1 removal scoping, ghost-id caveat — two lines)
- `docs/work-guide.md` (lines ~66-68, light)
- `tools/queue-kanban/verify.go` (ghost-finding remedy string only)
- `tools/queue-kanban/verify_test.go` (only if the remedy string is asserted there)
- `_dev/tests/contract-regressions.sh` (new writer-label assertions; existing pinned phrases preserved)
- `do-work/CHECKPOINT.md` (live entry gains this checkout's label)

**Acceptance criteria (restated from REQ):**
- In-Progress entries carry `writer: <hostname>:<absolute-checkout-path>`; recovery treats only own-label entries as candidates; foreign-label entries reported (`claim held by <writer>, not touched`), never stripped, never in the takeover ladder.
- Label-less legacy entries: own iff checkpoint locally modified/uncommitted here; else report-only. Never guess-strip.
- Tripwire still bans refresh intervals, staleness checks, liveness probes by name; static writer label explicitly carved out. `never grow into one` and `absent checkpoint is ambiguous` survive verbatim.
- CHECKPOINT.md stays tracked/committed; no heartbeat/staleness/auto-takeover machinery anywhere.
- `bash _dev/tests/contract-regressions.sh` green; `go test ./...` in tools/queue-kanban green.
