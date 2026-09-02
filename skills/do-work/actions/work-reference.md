# Work Action — Reference

> Companion file to `work.md`. Holds the heavy templates, tables, and sub-procedures the orchestrator references by name — extracted to keep `actions/work.md` focused on the ten-step skeleton. Each section below is pointed to from the matching step in `actions/work.md`. Loading this file is only necessary when you reach the step that references it — and read only the named section. If this file is already in context from an earlier step this session, reuse it; don't re-read it at every reference site.

---

## Architecture

```
work action (orchestrator - lightweight, stays in loop)
  │
  ├── Read CHECKPOINT.md if exists (crash recovery's classification input, then resume context)
  │
  ├── For each pending request (skip pending-answers):
  │     │
  │     ├── TRIAGE: Assess complexity (no agent, just read & categorize)
  │     │
  │     ├── OPEN QUESTIONS? ── - [ ] items exist ──► Mark - [~], builder decides
  │     │                      (none / all resolved) ──► continue
  │     │
  │     ├── ESTIMATE: Ensure estimate: block exists, print P50 line (never blocks)
  │     │     │
  │     │     ├── Route A (Simple) ──────────────────┐
  │     │     │   Skip plan/explore, direct to build │
  │     │     │                                      │
  │     │     ├── Route B (Medium) ───────┐          │
  │     │     │   Explore, scope declare  │          │
  │     │     │                           ▼          │
  │     │     └── Route C (Complex) ──► Plan ──► Explore ──► Scope declare
  │     │                                            │
  │     │                                            ▼
  │     │                                     Implementation agent
  │     │                                            │
  │     │                                            ▼
  │     │                                  Implementation Summary
  │     │                                            │
  │     │                                            ▼
  │     │                              Qualify (orchestrator verifies)
  │     │                                            │
  │     │                                            ▼
  │     │                                        Testing
  │     │                                            │
  │     │                                            ▼
  │     │                                  Review ◄─── Fail? ──► Remediate ──► Re-review
  │     │                                            │
  │     │                                            ▼
  │     ├── Archive ──► classify discovered tasks ──► queue follow-ups
  │     │                                            │
  │     │                                            ▼
  │     └── Commit (git repos only)
  │
  └── Context wipe → Loop | Write CHECKPOINT.md → cleanup → report
```

## Execution Model — Claim Anywhere, One Releaser

**Any checkout may capture and claim.** A queue can be shared by however many checkouts a user points at it — a spawned worktree, a second local workspace, a clone, a cloud sandbox — and each of them may write REQ files, claim them, and build. Claims and captures reach the other checkouts the way everything else does: by ordinary git sync, with no lock, no lease, and nothing to acquire — **and they travel when bookkeeping commits, not when the claim happens**. Nothing commits a claim at claim time: it reaches git only with the owner's bookkeeping (a session checkpoint, a hand-back's step 0 commit, the release tail), so a claim on work still being built is invisible to every other checkout, and another checkout can claim and build the same REQ in that window. That window is in contract, not a gap to close — the duplicate surfaces when the branches meet (**Cross-checkout collisions are merge artifacts**, below), and publishing claims earlier is exactly the coordination machinery this model deliberately does without.

What the model constrains is not who *works* but who *releases*. The pipeline supports **one releaser per queue** — the single checkout that runs the release tail: merge integration, the version bump, the `CHANGELOG.md` entry, the archive moves, and UR closure. **Two releasers against one queue stays outside the contract**, and so do two sessions in one working tree — the pipeline does not detect, coordinate, or recover either, and spends no durable state on making one safe. Behavior in both cases is unspecified, and the repair path is after the fact and human: `actions/forensics.md` to see the damage, `actions/cleanup.md` to fix it.

**Builders are not owners.** Any number may build at once (**Worktree Dispatch Mode** → Fan-Out Dispatch, below), because a builder writes only its own tree and owns no queue state. What changed is only that the tree may belong to a different checkout than the releaser's.

**Cross-checkout collisions are merge artifacts, not scheduling failures.** Two checkouts claiming the same REQ, or capturing REQs that land on the same id, produce ordinary git conflicts that get fixed when the branches meet — that is the entire coordination mechanism, and prose elsewhere cites this sentence rather than restating it. Both cheap detectors already ship: `queue-kanban verify`'s `duplicate-req-id` probe for colliding captures, and the `writer:` label on checkpoint entries for claims (**In-Progress Record (Step 2)**, below), which is what makes even a byte-identical double claim surface at merge. The label is also why crash recovery now *detects* another checkout's live claim — once that claim's bookkeeping commit has synced here; an uncommitted claim has nothing to detect — and reports it (`claim held by <writer>, not touched` — **Crash Recovery (Step 1)**, below). Detecting is all it does: it still neither coordinates nor recovers one, and it never arbitrates.

**Current-REQ relevance.** Unexpected repository **or session** state — a diff you did not make, a commit that lands mid-run, another live process against this checkout — matters **only** when it prevents the active REQ from being implemented, tested, archived, or committed. Otherwise: preserve it, exclude it from this REQ's staging, and continue — spend no time explaining or repairing it. **Never probe for a concurrent session, and never ask the user to arbitrate one.** Exclusivity is the user's guarantee, not this pipeline's check; a prompt asking them to confirm it rebuilds in conversation the coordination this model removed.

**Three-attempt stop.** The coder counts consecutive fix attempts **in its current context only**. After three failures, stop the local retry loop and report the unresolved blocker; use `pending-answers` only when progress genuinely requires a user decision. A restarted coder session starts with a fresh count — the count lives in coder context and is never written to disk.

## Folder Structure

```
do-work/
├── queue/                         # Pending REQ files (the work queue)
│   └── REQ-018-pending-task.md
├── user-requests/                 # UR folders (verbatim input + assets)
│   └── UR-003/
│       ├── input.md
│       └── assets/
├── working/                       # Currently being processed
│   └── REQ-020-in-progress.md
└── archive/                       # Completed work
    ├── UR-001/                    # Archived as self-contained unit
    │   ├── input.md
    │   ├── REQ-013-feature.md
    │   └── assets/
    ├── REQ-010-legacy-task.md     # Legacy REQs (no UR) archive directly
    └── legacy/                    # Consolidated legacy items
```

- **`queue/`**: The queue — only pending `REQ-*.md` files
- **`working/`**: Claimed requests. Immutable to all actions except the work pipeline.
- **`archive/`**: Completed UR folders (self-contained) and legacy REQs/CONTEXT docs
- **`user-requests/`**: Active UR folders. Moved to `archive/` when all REQs complete.

This tree exists in the **main working tree only**. Worktree dispatch mode's builder checkouts live outside the repo entirely and never carry their own copy of it — see **Worktree Dispatch Mode (Step 1)**.

## Request File Schema — Full Frontmatter

**Timestamp rule — every `*_at` field in this schema, and any timestamp a future field adds:** write the **current UTC instant** as ISO-8601 `YYYY-MM-DDTHH:MM:SSZ`. Write sites that say `<timestamp>` or `<now>` mean exactly this rule.

**This paragraph is the only place in `actions/` that spells a command for obtaining one.** Every other site cites the rule and stops. That is deliberate — while eleven sites each carried their own copy of `date -u`, no agent following one of them ever reached the order below, so a Windows agent kept getting a command that does not exist on its box no matter how the rule was fixed. The cost of the arrangement is that this is now a single point of failure for the whole skill's stamping: an error here reaches every stamp in the pipeline, so change it deliberately and re-read the three sources as a set.

Obtain the instant from the first of these that applies:

1. `<skill-root>/tools/do-work-cli.sh --format text now` — the canonical build-on-demand command. It prints exactly this shape and nothing else.
2. `date -u +%Y-%m-%dT%H:%M:%SZ` — the POSIX emergency floor when the canonical launcher cannot run.
3. **Windows, where option 2's flag form does not exist.** In PowerShell: `(Get-Date).ToUniversalTime().ToString("yyyy-MM-dd\THH:mm:ss\Z")`. From `cmd`, where a cmdlet is not a command at all, invoke it explicitly: `powershell -NoProfile -Command "(Get-Date).ToUniversalTime().ToString('yyyy-MM-dd\THH:mm:ss\Z')"`.

   Three things about that form are load-bearing. **`.ToUniversalTime()`, not `-AsUTC`** — the flag reads better but arrived in PowerShell 7, which ships as `pwsh.exe` and is absent from a stock box; `powershell.exe` is Windows PowerShell 5.1, where `-AsUTC` is an unrecognized parameter and the call fails outright. **The `\T` and `\Z` escapes** — in a .NET custom format string the specifiers are case-sensitive (`t`/`tt` are the AM/PM designators, `z`/`zz`/`zzz` the UTC offset), so bare `T` and `Z` happen to fall through to "copied unchanged"; the backslash says *literal* outright instead of depending on that. **`-NoProfile`** — a user profile can print banners into stdout and corrupt the captured value.

**Never stamp local wall-clock time with a `Z` suffix appended.** In any zone east of UTC that produces a *future* instant, which silently corrupts every elapsed-time reading (queue wait, claim stopwatch) and gets the REQ flagged by `do-work-board board` (a "future stamp" card badge plus a data warning, allowing 2 minutes of clock skew).

**Date-only stamps are a different shape and not part of the rule above.** The condition is what selects them, never a list of the sites that currently have one: **any UTC calendar date written into a durable record** — `YYYY-MM-DD` rather than an instant — is governed here. A log filename, a daily heading, a report header, the `**Answered [YYYY-MM-DD]:**` line `actions/clarify.md` writes into every REQ it touches, and the dated reasoning note beside it are all the same shape and take the same answer; they are illustrative, and a site added tomorrow is governed without editing this sentence. **UTC, not a local date**, so a sweep has one answer everywhere and no site-by-site judgment. It is spelled out here so every such site cites one home, but it is **not** an `*_at` value and must never be written into one. Obtain it as `date -u +%F` on POSIX, or `powershell -NoProfile -Command "(Get-Date).ToUniversalTime().ToString('yyyy-MM-dd')"` on Windows. There is no tool subcommand for it: `queue-kanban now` prints the instant only, and adding a date-only mode would spend the skill's narrow compiled-tooling exception (`../../do-work-board/actions/board.md` is the only capability allowed to need a compiler) on something the floor already covers. **A local-time date is a different thing again** and is correct where it is used deliberately (changelog entry headings, run-directory names, report slugs); those sites are not governed here. **A placeholder inside a template artifact is out of scope as well:** where a shipped prompt tells a model to emit a document for the user, a `Date: [today]` line is a fill-in token addressed to that model — a sibling of the `[name]` and `[owner]` placeholders in the same fenced block — not a stamp any step of this skill writes, so the date it ends up carrying belongs to the artifact's reader and their convention. The **condition** is what excludes it, never a list of which templates happen to carry one, so a template written tomorrow is excluded by this sentence without editing it. A sweep walks past every such placeholder instead of converting it, and the citation checker never reaches one — it masks fenced code before it reads. **A time-of-day heading is out of scope too:** the `## HH:MM UTC` daily-log headings (`../../do-work-knowledge/actions/memory-reference.md` → Daily-Log Entry Conventions) are neither an instant nor a date-only stamp — the log's dated filename already carries the date — so this rule deliberately does not govern them; their write sites are marked, and a sweep walks past every `## HH:MM UTC` heading instead of converting it.

**Named contract — Frontmatter Quoting.** The sibling of `actions/clarify.md` Step 4's **Outside-text containment**, which governs the same text once it lands in a do-work Markdown record *body*; this one governs it in frontmatter, and neither restates the other. **The condition is the rule: whenever a frontmatter value carries text nobody in this pipeline composed — a user's own words, a person's name, a command they supplied — write it so that no character in it can be read as YAML syntax.** `title`, `blocked_by`, `blocked_check`, `stakeholder` and `assigned_to` are today's such fields; they are illustrative, never a list to check against, so a field added tomorrow that carries user text is governed without editing this paragraph. Before choosing a scalar form, apply Outside-text containment's accepted-text preflight; a hand-authored frontmatter writer refuses and reports text that fails it instead of normalizing bytes or inventing its own escape table. Two forms, and only these:

- **The text fits on one line** → a **single-quoted scalar**, with every `'` inside it doubled (`''`). Inside single quotes `"`, `:`, `#`, `[`, `]` and `,` are ordinary characters, so the doubled apostrophe is the only escape and there is no escape table to get wrong.
- **The text contains an LF** → a **literal block scalar**, with every physical content line indented beneath the key (blank lines included). Choose the chomping indicator by the bytes being preserved: `key: |-` for zero terminal LF bytes, `key: |` for exactly one, and `key: |+` for multiple. No quoted scalar on one line can carry an LF, and a fixed `|-` silently strips a terminal one.

**Do not wrap user text in a double-quoted scalar.** A typed `"` ends it early, and the failure is not always loud: `title: "Fix: A " # B"` is *valid* YAML that reads back as `Fix: A`, the rest taken as a comment — no error, no warning, nothing to recover from. When it is loud instead, strict parsing rejects the **whole block**, not the one bad line. (A double-quoted scalar emitted by an *escaping encoder* is a different thing and is correct — the board's Testing view writes its two fields through `../../do-work-board/tools/queue-kanban/testing.go`'s `encodeYamlDoubleQuotedScalar`. This rule addresses a writer composing the line by hand, which is every action in this skill.)

**The board parser's line-based recovery is a last resort, not a contract a writer may aim at.** It exists so one bad line cannot cost a REQ its `status`, UR pointer and dependencies (`../../do-work-board/tools/queue-kanban/frontmatter.go`), and it holds for that. A valid strict parse is still the only target a writer may aim at, because a writer may **not** rely on recovery for any of the following — illustrative of its narrower contract, not the whole of it:

- **It answers for the whole block, not for the bad line.** Recovery is flat and top-level only, so the nested `estimate:` map below is silently dropped: one unquotable `title:` costs that REQ its whole forecast.
- **It rewrites values it does recover, with no warning.** A value that opens `[` and closes `]` is re-read as a YAML flow list, so the bracketed impact tag every REQ now carries loses its commas — worked through in `actions/capture-reference.md` → **REQ Title Convention**.
- **It does not recover block scalars at all**, so the newline form above exists only while the strict parse succeeds.

Pinned by `../../do-work-board/tools/queue-kanban/frontmatter_test.go` → `TestFrontmatterQuotingOfUserText`.

```yaml
---
# Set by capture action
id: REQ-001
title: 'Short descriptive title'   # raw user text — single-quoted per the **Frontmatter Quoting** contract above, apostrophes doubled; when `impact:` is anything other than `impact-user-visible`, lead with a `[<impact token>] ` tag (`actions/capture-reference.md` → **REQ Title Convention**, which is canonical for the shape and for every site that mints one)
status: pending
domain: frontend  # choose one: frontend, backend, ui-design, general, security, testing, or cms
tdd: false       # optional — set true when test-first applies (per capture's TDD heuristic); drives Step 6 testing-crew loading and RED/GREEN mode
caveman: false   # optional — `true` or intensity `lite` | `full` | `ultra`; loads crew-members/caveman.md to compress agent prose
maintenance: false  # optional — set true by capture for a removal/narrowing finding on the skill's OWN instructions (agent/action/crew/prime file); loads crew-members/maintenance.md (delete-before-you-add) in Step 6 alongside coding-guardrails. Not for ordinary app-source dead-code removal.
prime_files: []  # list paths to relevant prime-*.md files, or leave empty
required_lessons: [skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#final-boundary-identity]  # OPTIONAL capture-authored mandatory-read list. Entry forms and the single token budget are owned by `actions/capture-reference.md` → Required Lessons Budget Contract. Bare `path` means the whole satellite; `path#family-slug` means matching family bullets and is valid only for an index row marked `slugged: full`. No index match means absent, never an invented default. Verbatim-read path-list class: no aliases, normalization, case folding, or path canonicalization; every requestmodel/schemanormalization writer preserves the field bytes unchanged. The board does not parse or display it.
created_at: 2025-01-26T10:00:00Z
user_request: UR-001          # May be absent on legacy REQs
addendum_to: REQ-NNN          # optional — present only when this REQ amends an in-flight or completed REQ; set by capture, or by review when creating follow-ups. **Legacy alias:** every read site (Step 8 upstream walk, Step 8 cycle detection, Step 8 follow-up generation, and roadmap Blocked classification) also recognizes `amends:`, `parent:`, and `amendment_to:` as synonyms when `addendum_to` is absent so natural-English glosses don't silently drop the parent linkage; `addendum_to` wins when multiple are present. Capture and follow-up REQs always emit `addendum_to:` — never propagate the alias.
depends_on: []                # optional list of REQ IDs that must reach `completed` or `completed-with-issues` before this REQ runs. Semantically distinct from `addendum_to` ("amends that REQ"): depends_on is "requires that REQ to be done first." A REQ can have both. Honored by Step 1's selection scan and by Step 8's upstream-failure classification. **Legacy alias:** every read site (Step 1 selection, Step 1 cycle detection, Step 1 `--wave` depth, Step 8 upstream walk, roadmap classification) also recognizes a `dependencies:` key as a synonym so muscle-memory typos from Python/Node/Cargo conventions don't silently bypass gating; `depends_on` wins when both are present. Capture and follow-up REQs always emit `depends_on:` — never propagate the alias.
gate_deferred: true           # optional exact boolean marker written only by the canonical repository-gate deferral transaction. The parent remains `pending` and names its repair in `depends_on`; this marker changes selector priority after the dependency succeeds, not status semantics.
repository_gate_repair: true  # optional exact boolean marker on a generated repository-gate repair REQ. Such REQs are pending work under the source parent’s `user_request`, use `related` for parent provenance, and never use `addendum_to`, `blocked`, or `pending-answers` for this lifecycle.
deferred_implementation_base: 0123456789abcdef0123456789abcdef01234567   # optional full commit paired with deferred_implementation_merge for a late deferral; both are absent for a pre-build deferral.
deferred_implementation_merge: 89abcdef0123456789abcdef0123456789abcdef # optional non-empty descendant of deferred_implementation_base whose implementation range is revalidated before late resumption.
write_set: []                 # optional list of repo-relative paths/globs this REQ declares it will write. **Display only, at any builder count** — under fan-out it is advisory input to the human's pick and the merge is the non-interference proof, never this field (**Worktree Dispatch Mode (Step 1)** → Fan-Out Dispatch, below), so nothing schedules, gates, or dispatches on it; it feeds the board's overlaps badge. Seeded by capture when the request names files (`actions/capture-reference.md` → Populating `write_set`), then firmed by Step 5.5: the `## Scope` section's "Files I will touch" list is the source and this field is its mirror, never the reverse. **Crash recovery clears it**, but only when the REQ has a `## Scope` for the mirror to have come from (**Crash Recovery (Step 1)**, substep 1). Parsed for display by `../../do-work-board/tools/queue-kanban/model.go`; **an absent or empty `write_set` gets no overlaps badge at all** — absence reads as *unknown*, not conflict (matching `../../do-work-board/actions/board.md` and `../../do-work-board/tools/queue-kanban/prime-do-kanban.md`). Globs use `path.Match` semantics: `*` never crosses `/`, `**` is not recursive, a malformed pattern is no-match for that direction though literal equality short-circuits first (so two REQs declaring the identical malformed pattern still badge each other), and an entry naming a directory never matches a file inside it (`actions/` vs `../../do-work-board/actions/board.md`) — those board miss-classes are illustrative, not a closed list.

assigned_to: 'cloud-alpha'     # OPTIONAL advisory claim marker: the session this REQ is earmarked for (raw user text — **Frontmatter Quoting** contract above). **Verbatim-read class**, alongside `write_set` — no alias map, no case folding, no canonical vocabulary of session names (Schema Read Contract → the verbatim paragraph, which also states the one shared exception: surrounding whitespace is trimmed, as it is for every field in the class). Seeded by capture when the user earmarks work, or written by a session claiming from another checkout. **Not a lock and not a status:** it grants nothing, nothing waits on it, and it carries no `assigned_at` and no staleness clock — an assignment persists until an explicit run or a hand-edit clears it. Exactly one reader acts on it, as a courtesy rather than a gate: the default work scan skips and reports an assigned REQ, and explicit targeting (`do-work run REQ-NNN`) overrides the skip and clears the field as part of the claim (`actions/work.md` Step 1). Everything else is display: parsed by `../../do-work-board/tools/queue-kanban/model.go` into the card's `assigned` badge and a drawer row, with **no column logic and no scheduling** — keep that parser in lock-step with this line, both changing in the same commit.

review_generated: true          # OPTIONAL exact marker written on review-created follow-ups (`actions/review-work.md` Step 10). The board's parser reads true only when the scalar coercion result is exactly `true`; absent, false, and non-canonical values remain false, with no aliases or normalize-and-warn behavior. The value has no display or scheduling role. Its one board consumer is the read-only archived-UR diagnostic: after terminal queue/working members are left to the stranded-finished probe, a non-terminal exact-true member under its already-archived UR is the legitimate same-UR follow-up shape and does not reopen or move that UR.
sweep: true                    # OPTIONAL marker: this REQ is a consolidation sweep — ONE REQ per root cause, carrying an `## Instances` checklist of every occurrence. Boolean; absent reads as false. Greppable by design (`grep -rl "^sweep: true" do-work/queue/`) so minting flows find the existing sweep instead of judging titles; a `pending` or `pending-answers` sweep is appendable, a claimed one never is. Who writes it and the conversion and escalation rules are stated once in `actions/capture-reference.md` → **Fold-First Rule** — do not re-derive them here. Parsed by `../../do-work-board/tools/queue-kanban/model.go` into a display-only card chip and drawer row carrying the `## Instances` open/done counts — no column logic, no scheduling. Keep that parser in lock-step with this line, both changing in the same commit. The marker is not in the normalize-and-warn table (marker class, exact `true` only).
sweep_key: hardcoded-colors-untokenized   # REQUIRED on sweep REQs (meaningless elsewhere): short kebab-case name for the ROOT CAUSE — the deterministic append discriminator. A minting flow appends to the candidate whose key matches its finding's root cause; when no key matches literally, it compares root-cause statements (same rule = same sweep), never titles (`actions/review-work.md` Step 10). Free slug, no canonical vocabulary — verbatim-read class.
impact: impact-user-visible   # OPTIONAL in the schema, expected on every new REQ: impact-critical | impact-user-visible | impact-rule-change | impact-negligible. Whether anyone would ever notice the work, judged by the two questions in `actions/review-work.md` Step 10. This is the field the user filters on to stop implementing work nobody would notice — the question `effort_estimate` was never able to answer, because it means size. Every value carries the `impact-` prefix, so `grep -rn 'impact-' do-work/queue/` returns every REQ's verdict in one pass. `impact-critical` (security, data loss, a broken production path) skips the consent question and auto-queues at any depth. **Absent or unrecognized reads as `impact-user-visible`** (Schema Read Contract row below) — deliberately never `impact-negligible`, because absence must never be mistakable for the user's stop signal. Written by capture (`actions/capture.md` Step 1's impact assessment) and by automatic follow-up creation (`actions/review-work.md` Step 10; **Discovered Tasks Classification (Step 8)** below) from the finding's recorded token — the token IS the field value, so no translation table sits between the finding line and the frontmatter. **This field is the source of truth and the REQ title is a mirror of it:** a non-default verdict is also written into `title:` as a leading `[<impact token>] ` tag (`actions/capture-reference.md` → **REQ Title Convention**) so the board's title-matching search box can find it, and when the two disagree the field wins — every reader that acts on impact reads the field, never the title. Read as a filter by `actions/work.md` Step 1's `--skip-impact-negligible`, which omits `impact-negligible` REQs from selection; the normalization above is what keeps that flag conservative, since an absent or unrecognized value reads as `impact-user-visible` and is never skipped. *(REQ-228 ruled "no new frontmatter field … `effort_estimate` stays a two-value triage bit." That was a constraint on **persisting a derived forecast** — the board's projection must stay render-time so it never becomes a write surface. `impact:` is not a forecast and nothing derives it; it is an authored judgment, written once at capture or follow-up creation, that the board only reads. REQ-228's other half — that `effort_estimate` must not grow toward t-shirt sizes — is honored exactly: this line renames its two values, it does not widen the enum.)* Display only: parsed by `../../do-work-board/tools/queue-kanban/model.go` into a card chip (rendered for every value except the `impact-user-visible` default) and a drawer row, with no column logic and no scheduling — keep that parser in lock-step with this line, both changing in the same commit.
effort_estimate: effort-substantive   # OPTIONAL triage bit: effort-mechanical | effort-substantive — separates small mechanical fixes from real work so the user can tell at a glance which queued REQs are cheap to approve or batch. **This field is SIZE, judged as size by whoever writes it — never derived from an impact verdict; that judgment is `impact:` above.** Closed two-value enum, deliberately — a triage bit, not an estimation system; do not grow it toward t-shirt sizes (the estimation system lives in the separate `estimate:` block below, and the only bridge between the two is the mechanical-effort short-circuit: `effort-mechanical` ⇒ the floor estimate, no signal extraction — `actions/estimate-reference.md`). Absent or unrecognized reads as `effort-substantive` (Schema Read Contract row below), and the read-only legacy aliases in that row carry every REQ written before the rename unchanged. **Expected on every new REQ:** capture, review follow-up creation, and Discovered Tasks creation judge it by the same three-way standard — judge it, or put the judgment to the user, or leave it absent because neither was possible, never a copied default (`actions/capture.md` Step 1's effort assessment; `actions/review-work.md` Step 10; **Discovered Tasks Classification (Step 8)** below). An unjudged REQ reads as `effort-substantive` and drops out of the selector below for no reason anyone chose. Read as a selection filter by `tools/select-simple-reqs.sh` (backing `actions/run-simple-reqs.md`), which selects only the REQs normalizing to `effort-mechanical` so a cheaper-model session can run the queue's small work; the read-only `trivial` alias above is load-bearing for that selector, because REQs written before the rename spell the value that way and a literal match on the canonical token alone silently drops them. Display only: parsed by `../../do-work-board/tools/queue-kanban/model.go` into a card chip (rendered only when `effort-mechanical` — `effort-substantive` is the default and would be noise) and a drawer row, with no column logic and no scheduling — keep that parser in lock-step with this line, both changing in the same commit.

# OPTIONAL informational forecast — backwards-compatible: a REQ without it is fully
# valid, and no scheduling, gating, or pipeline logic ever reads it. Written by the
# work action's ensure-estimate step (post-triage) or by verify-requests after a
# material repair, then FROZEN once execution begins — never rewritten with knowledge
# gained during implementation. p50_active_minutes is a multiple of 5, never below 5;
# it means roughly a 50% chance of completing within that many ACTIVE agent minutes
# (user wait, paused/suspended sessions, and queue wait are excluded by definition).
# Produced deterministically by tools/estimate-p50.sh from extracted signals; the
# extraction guide, confidence rubric, and presentation formats live in
# actions/estimate-reference.md. No P80 or other percentile fields — by design.
estimate:
  p50_active_minutes: 75
  confidence: medium            # low | medium | high
  calculated_at: 2026-08-16T12:00:00Z   # current UTC instant (Timestamp rule above)
  basis:
    - Route C
    - 12-file write set
    - browser evidence

# Set by work action when claimed
claimed_at: 2025-01-26T10:30:00Z
route: A | B | C

# OPTIONAL observations written only after the named work-pipeline event
# succeeds. Routes that skip an event omit its field; writers never fabricate a
# value for a skipped phase. All eight use the Timestamp rule above.
planning_at: 2025-01-26T10:32:00Z         # Route C plan saved and validated
dispatch_at: 2025-01-26T10:33:00Z         # implementation builder accepted the dispatch
builder_handback_at: 2025-01-26T10:40:00Z # builder returned its completed hand-back
integration_at: 2025-01-26T10:41:00Z      # hand-back integrated (or accepted in serial mode)
review_at: 2025-01-26T10:43:00Z           # first review result recorded
remediation_at: 2025-01-26T10:46:00Z      # one remediation hand-back integrated
re_review_at: 2025-01-26T10:48:00Z        # post-remediation review result recorded
release_at: 2025-01-26T10:50:00Z          # canonical release transaction succeeded

# Set by capture (external-condition task) or by the work pipeline's mid-run blocked flip (Step 8's blocked-flip procedure). Holding state — the REQ stays in do-work/queue/ and the default scan walks past it, exactly like pending-answers.
status: blocked               # waiting on an EXTERNAL condition — not user answers (that's pending-answers), not another REQ (that's depends_on). Cleared to `pending` by a passing blocked_check probe (work Step 1), a `do-work clarify` confirmation, or a manual edit.
blocked_by: 'LM Studio running locally'   # free text naming the condition (raw user text — **Frontmatter Quoting** contract above). Legacy note for the board: an old id-LIST value renders joined for display and is NOT a dependency edge — dependency gating is `depends_on` only.
blocked_at: 2026-07-18T10:00:00Z          # stamped on every flip to blocked — the age anchor the exit summary, board drawer, and forensics read (no enforcement threshold; external conditions legitimately take weeks)
blocked_check: 'curl -sf http://localhost:1234/v1/models'   # OPTIONAL shell probe (raw user text — **Frontmatter Quoting** contract above). User-authored content, run VERBATIM by work Step 1 (exit 0 ⇒ unblock to pending; any non-zero / timeout / unreadable ⇒ stays blocked). Absent ⇒ manual/clarify unblock only.
stakeholder: 'Priya (design)'   # REQUIRED on stakeholder-questions REQs (meaningless elsewhere; raw user text — **Frontmatter Quoting** contract above): the outside person whose confirm-or-override answers this REQ collects — presence is the marker, value is the fold discriminator (actions/capture-reference.md → Fold-First Rule → Stakeholder-audience questions). Verbatim-read class, like assigned_to: no alias map, no case folding, trim-only; greppable by design (grep -rl '^stakeholder: ' do-work/queue/). Always paired with status: blocked + blocked_by naming the person and the latest report bundle path (or "report pending regeneration" until a bundle lands — actions/work.md Step 8) + blocked_at; never with blocked_check (a person is not probeable), and deliberately never with user_request: — UR membership would hold the first source UR open in every closure reader, and nothing waits on this REQ (question provenance lives in per-entry Source: lines). NOT parsed by the board — display rides the existing blocked_by badge, zero parser change. Nothing gates on this REQ: its source REQs completed on the builder's assumptions, and it exists only to route answers back (actions/stakeholder-answers.md); clarify routes it, never yes/no-confirms it (actions/clarify.md Step 5.5).

# Set on ANY status flip that has no dedicated *_at stamp of its own — that
# condition is the rule, the writers are illustrative: answered → pending
# (clarify Step 5), unblock → pending (clarify Step 5.5, work Step 1 probe —
# both REMOVE blocked_at, so this is the only trace of when the flip happened),
# manual/stuck resets back to pending. Flips with a dedicated stamp (claim →
# claimed_at, blocked → blocked_at, terminal →
# completed_at) do NOT write it. Display-only: the board's state timer prefers
# it over created_at/file-mtime for pending-tier cards ("updated … · 3m"); no
# pipeline logic reads it. Timestamp rule applies (current UTC instant).
status_changed_at: 2026-07-22T20:38:00Z

# Set by work action when finished. STAMPING RULE: every flip to a terminal
# status (completed / completed-with-issues / failed / cancelled) MUST stamp
# completed_at with a UTC ISO instant, plus commit with the implementation
# hash in a git repo. These two fields are the ONLY sources the board resolves
# a terminal REQ's completion instant from (no file-mtime fallback); a
# terminal REQ missing both — or carrying an unparseable completed_at or a
# hash git can't resolve — is flagged as a completion anomaly by do-work-board board
# (all three modes: serve, static, summary).
completed_at: 2025-01-26T10:45:00Z   # required on every terminal flip — UTC ISO instant
status: completed | completed-with-issues | failed
commit: abc1234               # required in a git repo — implementation commit hash (see work.md's Commit Phase write-back)
error: "Description"          # Set when a REQ failed; RETAINED verbatim if that failed REQ is later cancelled via do-work abandon — the surviving failure signal on a status: cancelled REQ, NOT drift to strip
error_type: intent|spec|code|environment   # Set with `error` on failure; likewise retained on a failed→cancelled flip

# Set by abandon action (do-work abandon — user-directed won't-do decision).
# Two entry paths: a not-yet-finished REQ (pending / pending-answers / blocked / ...), and an
# already-archived `failed` REQ resolved after the fact — the latter keeps its `error`/`error_type`
# (above) alongside status: cancelled, so error-on-cancelled is valid data, not corruption.
status: cancelled             # terminal, NOT successful; the reason lives in the REQ body's `## Cancelled` section
completed_at: 2025-01-26T10:45:00Z  # stamped (or, on the failed→cancelled path, re-stamped to the cancellation instant) — the terminal timestamp the board's recently-done window reads

# Set by kb-lessons handoff (work.md's Lessons-Capture Phase in pipeline mode / review-work.md's Self-Validation & Lessons Learned step standalone). Optional; absent on REQs that predate the handoff.
kb_status: promoted | pending | declined | skipped
kb_entry: REQ-042-lesson-slug.md   # filename only (survives bkb moves from inbox/ to capture/ to processed/); present only when kb_status: promoted

# Set by the board's Testing view (do-work-board board — `../../do-work-board/actions/board.md` Step 6). Optional; the testing track is orthogonal to `status`: the board never writes `status`, and the work pipeline never writes these. Absent = not tested yet.
testing_status: in-testing | tested | returned   # who-tested-what tracking for finished REQs
tested_by: "Alice"                # tester profile from do-work/testers.md (raw user text; the board emits it through its escaping encoder — **Frontmatter Quoting** contract above)
testing_updated_at: 2026-07-17T10:00:00Z   # stamped by the board server on every transition
testing_feedback: "…"             # present only while testing_status is returned (raw user text; same escaping encoder — one-line double-quoted scalar, newlines as \n escapes)
---
```

## Schema Read Contract

The enum-or-boolean-valued fields above (one table row each, below) are covered by this contract; an audit of `0.76.2`'s `dependencies:` → `depends_on` patch surfaced that several silently swallow natural typo variants from sister conventions (snake_case vs kebab-case YAML, `complete`/`done`/`finished`/`closed` as English glosses of `completed`, lowercase route letters, etc.). Pure silent-alias is risky for enum values because an unknown value should not be quietly remapped — it should leave a footprint. Every read site in this file (and in `actions/roadmap.md`) honors a uniform **normalize-and-warn contract** for these fields:

1. **Normalize first.** Apply the per-field alias map below. If a canonical match results, use it silently.
2. **Warn-on-fallback.** If after normalization the value still doesn't match the canonical enum, emit:

   ```
   ⚠ {field}: '{value}' not recognized — expected one of [{enum}]. Treating as '{default}'.
   ```

   and proceed with the documented default.
3. **Never silently drop.** The warning is the missing feedback channel that allowed `dependencies:` to go unnoticed pre-0.76.2. Warnings render in the queue-status summary block (Step 1) or, for fields read outside Step 1, alongside the operation that triggered the read.

| Field (read sites) | Canonical enum | Normalization | Default on unknown |
|---|---|---|---|
| `domain` (Step 4 Route C plan-agent spawn, Step 6 crew load, Step 7 review-work spawn) | `frontend`, `backend`, `ui-design`, `general`, `security`, `testing`, `cms` | `back-end`/`back_end` → `backend`; `front-end`/`front_end` → `frontend`; `ui_design` → `ui-design`; `sec` → `security`; `test` → `testing`; `content-management`/`content_management` → `cms` | `general` |
| `status` (Step 1 scan + categorization, Step 8 archive trigger, abandon action) | `pending`, `claimed`, `completed`, `completed-with-issues`, `failed`, `cancelled`, `pending-answers`, `blocked`, `blocked-archive-collision`, `blocked-dependency-cycle` | `complete`/`done`/`finished`/`closed` → `completed`; `canceled`/`abandoned`/`wont-do`/`wontfix` → `cancelled` | skip REQ at Step 1 with the warning text — never claim or archive an unrecognized status silently |
| `route` (Step 3 dispatch, Step 5.5 scope declaration, Step 7 scope-drift comparison) | `A`, `B`, `C` | lowercase `a`/`b`/`c` → uppercase | treat as needing re-triage in Step 3 |
| `caveman` (Step 6 crew load) | `false`, `true`, `lite`, `full`, `ultra` | truthy strings (`yes`/`on`) → `true`; `light` → `lite` | `false` |
| `maintenance` (Step 6 crew load) | `true`, `false` (YAML boolean) | truthy strings (`yes`/`on`/`t`) → `true`; `no`/`off`/`f` → `false` | `false` (Step 6 maintenance crew not loaded) |
| `tdd` (Step 6 testing-crew load, Step 6.5 TDD-evidence gate; emission validated in `actions/capture.md`) | `true`, `false` (YAML boolean) | `test_first`/`yes`/`on`/`t` → `true`; `no`/`off`/`f` → `false` | `false` (Step 6 testing crew not loaded; Step 6.5 gate not enforced) |
| `error_type` (Step 8 failure classification, Step 8 upstream-failure short-circuit, forensics) | `intent`, `spec`, `code`, `environment` | (no common typo aliases identified) | `code` |
| `kb_status` (kb-lessons handoff — work.md's Lessons-Capture Phase / review-work.md's Self-Validation & Lessons Learned step; roadmap lessons rollup) | `promoted`, `pending`, `declined`, `skipped` | `skip` → `skipped`; `rejected` → `declined` | `pending` |
| `impact` (capture emission — capture.md Step 1; automatic follow-up creation — review-work.md Step 10, work.md Step 8's Discovered Tasks flow; selection filters — work.md Step 1's `--skip-impact-negligible` and `tools/select-simple-reqs.sh`'s `impact-critical` veto; board display — `../../do-work-board/tools/queue-kanban` parser) | `impact-critical`, `impact-user-visible`, `impact-rule-change`, `impact-negligible` | (no aliases — new prefix-unique vocabulary) | `impact-user-visible` |
| `effort_estimate` (Step 3.6's mechanical-effort short-circuit; selection filter — `tools/select-simple-reqs.sh`, backing `actions/run-simple-reqs.md`; board display — `../../do-work-board/tools/queue-kanban` parser; judged by capture, review follow-up creation, and Discovered Tasks creation on every new REQ) | `effort-mechanical`, `effort-substantive` | `trivial` → `effort-mechanical`; `normal` → `effort-substantive` (read-only legacy aliases, so every REQ written before the rename stays valid unchanged; never propagated on write — `actions/capture-reference.md` § Schema Aliases) | `effort-substantive` |
| `testing_status` (board Testing view — `../../do-work-board/tools/queue-kanban` parser + `/api/testing/status` writes; no work-pipeline read sites) | `in-testing`, `tested`, `returned` | `in_testing`/`in testing`/`testing`/`selected-for-testing`/`selected for testing` → `in-testing`; `returned-with-feedback`/`returned_with_feedback`/`returned with feedback` → `returned` | treat as not-tested (Ready to test) with an invalid flag + data warning |
| `builder_decided` (clarify's confirm routing — `actions/clarify.md` Step 4/Step 5; reversal detection — `actions/clarify.md` Step 4's `overturned_decision_sources` and `actions/verify-requests.md`'s Decision Revalidation Workflow; doctor's `HOLLOW-COMPLETION` no-code-change exception) | exact `true` only | (no aliases — marker class, exactly like `sweep` and `review_generated`) | absent reads as false |
| `gate_deferred` (canonical repository-gate deferral marker; selector priority after dependencies are satisfied) | `true`, `false` | truthy strings (`yes`/`on`/`t`) → `true`; `no`/`off`/`f` → `false` | `false` |
| `repository_gate_repair` (generated repair marker; selector priority) | `true`, `false` | truthy strings (`yes`/`on`/`t`) → `true`; `no`/`off`/`f` → `false` | `false` |
| `deferred_implementation_base` (late-deferral saved range) | any non-empty commit text | none; trim only | absent |
| `deferred_implementation_merge` (late-deferral saved range) | any non-empty commit text | none; trim only | absent |

**Write paths are unaffected.** Step 2 claim, Step 8 archive, Step 8 follow-up generation, the kb-lessons handoff, and capture emission always write the canonical key and canonical enum value — never an alias, never the typo'd input. The normalize-and-warn contract is read-only.

**Fields with no canonical vocabulary are outside this contract and are read verbatim.** `prime_files`, `write_set`, `required_lessons`, and any future path list have no canonical vocabulary to normalize against — no alias map, no case folding, no path canonicalization, no warning. `assigned_to` joins them for the same reason: a session name is whatever the user or the assigning checkout called itself, so there is nothing to normalize *against* and folding case would silently make two distinct sessions look like one. `stakeholder` is the same class again — a person's name is whatever the user called them, and folding case would make two people look like one; its literal-then-same-person fold match is stated in `actions/capture-reference.md` → Fold-First Rule → Stakeholder-audience questions. A reader takes the strings as written, with one narrow exception that applies to **every** field in this class: **surrounding whitespace is trimmed.** YAML already strips it from unquoted scalars, so it only survives explicit quoting (`" cloud-alpha "`), it carries no meaning in a name or a path, and treating `" cloud-alpha "` and `"cloud-alpha"` as two different sessions would break the skip-and-report over a difference nobody intended. Verbatim here means *no alias map, no case folding, no path canonicalization* — not byte-preservation of padding. (`depends_on`, `related`, and `blocked_by` are likewise row-less here; `depends_on`/`related` carry alias keys, documented on their schema lines above, but their *values* are also read verbatim.)

`deferred_implementation_base` and `deferred_implementation_merge` use that same trim-only projection, but the canonical deferral writer accepts them only as a pair, resolves both to commits, requires a non-empty ancestor range, and persists the full commit IDs. A generic reader retains their scalar evidence; resumption owns the later ancestry and path-drift judgment.

**Optional phase-stamp read contract.** `planning_at`, `dispatch_at`, `builder_handback_at`, `integration_at`, `review_at`, `remediation_at`, `re_review_at`, and `release_at` are additive scalar observations. Absence means that phase was not observed and must render no phase, zero, or inferred instant. A present value is trimmed and parsed by the Timestamp rule; an unparseable value is ignored by phase-duration derivation rather than normalized or replaced. The board keeps the declared pipeline order, omits missing/unparseable phases, and measures each displayed interval from the previous parseable observation. None of these fields changes calibration: its only span remains the just-archived REQ's `claimed_at` → `completed_at`.

### Repository-gate deferral transaction

`do-work-cli defer-gate --manifest <path>` is the sole mutation owner for converting a diagnosed unrelated repository-gate failure into repair work. The caller supplies judgment and exact evidence: parent id/path and bytes, checkpoint path and bytes, expected `claimed` state, writer label, structured gate argv, direct non-zero status, diagnostic fingerprint and evidence, stable root-cause `sweep_key`, repair id/path/title/reservation, and optional paired implementation base/merge commits. The command refuses stale preimages, staged owned targets, identity/path collisions, ambiguous fold candidates, unsafe topology, malformed evidence, and invalid merge ranges before publication.

One successful transaction creates or uniquely folds a `pending` `repository_gate_repair: true`, `sweep: true` repair under the parent's `user_request`, projects every affected parent as a canonical open checklist item under `## Instances`, appends each parent through `related`, changes the parent to `status: pending`, adds the repair to canonical `depends_on`, sets `gate_deferred: true`, removes `claimed_at` and only the exact writer's checkpoint entry including its indented detail, appends `## Repository Gate Deferral`, and moves the parent from `working/` to `queue/`. A fingerprint folds only when both `sweep_key` and an exact parsed diagnostic-evidence field match; prefix or substring matches are not identity. The repair never uses `addendum_to`; the parent never enters `blocked` or `pending-answers`, because this lifecycle requires no user choice. Parent, checkpoint, and folded-repair preimages are classified independently as tracked-dirty, tracked-clean, or untracked; every move destination must be absent during planning and remains protected by exclusive creation during apply. Any error after mutation begins restores every parent, repair, reservation, destination, mode, and checkpoint byte to its exact preimage.

Default and UR-expanded queue selection order ready work by three stable classes: `repository_gate_repair: true` first, ready `gate_deferred: true` parents second, and ordinary work third. Existing queue order is preserved inside each class. Explicit REQ tokens retain caller order and bypass dependency readiness as before; `selection_priority` is still projected on selected records and fan-out exclusions so the caller never has to infer the class from display text.

### Terminal-success status set

**After applying the `status` alias map, a REQ counts as *terminally successful* when its status is `completed` or `completed-with-issues`.** This is the canonical set every reader that selects "completed work" must honor — `completed-with-issues` is terminal and counts toward UR completion (it just carries known follow-ups, per `actions/work.md` Step 8), so a filter that accepts only the literal `completed` silently drops remediated-with-issues work. `failed` is terminal but **not** successful — success-readers exclude it.

The trigger is the *condition above*, not the caller list: **any reader that filters for "the completed/most-recent work" inherits this contract.** The known consumers are illustrative, not exhaustive — `actions/cleanup.md` (UR close), `../../do-work-toolbox/actions/completed-work-presentation-reference.md` (shared reader for item-level presentation actions), `actions/review-work.md` (standalone target), and `actions/commit.md` (REQ association); `actions/forensics.md` and `actions/roadmap.md` already honor both. When adding a new reader, normalize status first, accept both canonical values, and point back here — hand-enumerated caller lists go stale silently, which is why the condition, not the list, is the contract.

### Terminal-resolved status set

**After applying the `status` alias map, a REQ counts as *terminally resolved* when its status is `completed`, `completed-with-issues`, or `cancelled`.** This is the set archive-sweep and UR-closure readers honor — any reader deciding whether a REQ still needs work inherits it, so the list that follows is illustrative, not exhaustive: `actions/cleanup.md` Pass 0 + Pass 1, `actions/work.md` Step 8's UR-final check, doctor's `ORPHANED-USER-REQUEST` finding, and this file's Composed Exit Summary. `cancelled` records a deliberate won't-do decision — made via `do-work abandon` — so it archives like finished work and must never hold a UR open the way `failed` does. Three boundaries keep the sets honest:

- `cancelled` is **not** successful. Success-readers (the Terminal-success set above) exclude it — a cancelled REQ is never a review-work target, a completed-work presentation action target, or a commit association.
- `cancelled` does **not** satisfy `depends_on` gating. A dependent presumably needed the cancelled REQ's output; the abandon action surfaces dependents at cancellation time so the user can cascade the cancellation or re-point `depends_on`.
- `failed` stays outside this set: it is terminal and unsuccessful, but it signals work that *should* have happened — Step 8's failure classification may spawn a follow-up REQ. A follow-up does the recovery work but **never flips the original out of `failed`** (nothing does so automatically), so a UR holding a `failed` REQ stays open even after that follow-up completes. The one transition out of `failed` is `do-work abandon REQ-NNN` (`actions/abandon.md`), which flips it to `cancelled` — and therefore into this set — while preserving the failure record (`error`/`error_type` and a `## Cancelled` note). So a UR held open by a `failed` REQ closes only once that REQ is cancelled: after its follow-up has done the needed work, or when no follow-up is wanted at all. A legacy `failed` REQ already inside a closed `archive/UR-NNN/` folder is the same explicit-target transition: abandon cancels it in place so it leaves the board's active view, without moving the file or reopening the UR. `failed` itself never counts as resolved. This is the canonical statement of that resolution **rule**: any reader that decides whether a `failed` REQ still holds its UR open cites this set by reference and must not restate or fork the set, or re-derive the rule, as a competing definition. The condition is the trigger, not a caller list — the known such readers are illustrative, not exhaustive (`actions/cleanup.md` Pass 0 + Pass 1, doctor's `ORPHANED-USER-REQUEST` finding, `actions/work.md` Step 8's UR-final check, and this file's Composed Exit Summary), and `actions/abandon.md` is the rule's sole *writer*. (A user-facing report or finding line that *points* the user at `do-work abandon` as the remedy — and may state the one-line reason why, e.g. that a completed follow-up never resolves the original — is a pointer, not a competing definition, and is expected; `actions/cleanup.md` Pass 1's open-UR report and doctor's `FAILED-UNCLASSIFIED` / `FAILED-WITHOUT-FOLLOW-UP` findings do exactly this. That is not what this prohibition forbids; what it forbids is a reader re-deriving *when a REQ counts as resolved* as its own competing definition.)

### Target ID Resolution

**The shared token grammar for every action that takes id arguments** — `do-work run`, `do-work abandon`, and `do-work roadmap`. A caller cites this contract for token shapes and UR expansion instead of restating them; the condition, not a caller list, is the trigger, so any future id-taking action inherits it.

- **Token shapes.** `REQ-` + digits and `UR-` + digits, **case-insensitive**. **Canonicalize the token to the stored form before any lookup** — uppercase the prefix, and match the digits by **numeric value** against the stored (zero-padded) id, so `req-42`, `REQ-42`, and `REQ-042` all resolve to `REQ-042`, and `Ur-11`/`UR-011` both resolve to `UR-011`. Stored ids are zero-padded and upper-case (`REQ-067`, `UR-011`); callers glob and compare `user_request:` against that canonical form, so a resolver must normalize first and never pass raw user text into a case-sensitive glob or string compare.
- **A `UR-NNN` token expands to its member REQs** by scanning `user_request:` frontmatter across the live locations *the calling action already searches* — which locations those are is the caller's business; the scan key never is. It is **never** the UR's own `requests:` array, which is a capture-time record and explicitly not a membership predicate (`actions/capture.md` → "The `requests:` array is the capture-time record only"). Every prior bug here — REQ-048, REQ-058, REQ-059 — came from reading that array.
- **An id that resolves to nothing** is reported by id and skipped, never silently dropped.
- **An argument list that resolves to an empty set stops the action.** It never falls through to a whole-queue default (a full run, a full survey). Expansion adds a recognized shape; it never makes an unrecognized or empty argument permissive.
- **Expansion widens *which* REQs an action reaches; it never relaxes how any one is treated.** Each caller applies its own per-REQ gates — dependency-readiness, status refusals, confirmations — to every expanded member exactly as if that REQ had been named directly.

## Crash Recovery (Step 1)

**Crash Recovery:** Before checking the queue — but **after** reading `do-work/CHECKPOINT.md`, which is this procedure's input (`actions/work.md` Step 1) — look inside `do-work/working/` for any `REQ-*.md` files. A file there is a claim that outlived the run that made it: this session's own interrupted work, a claim the checkpoint attributes to another checkout, or a claim this session cannot account for at all.

**Recovery is destructive, so it is not the default.** Substeps 1–3 below reset the frontmatter, strip thirteen generated sections (`## Plan`, `## Exploration`, `## Scope` and the rest), and move the file back to `do-work/queue/`. Nothing is committed before Step 9, so those sections usually exist nowhere but that file — which makes substeps 1–3 the right treatment for a crash's half-written leftovers and the wrong treatment for anything else. **Classify each `working/` file before touching it:**

- **Named in the checkpoint's `## In Progress (interrupted)` record under this checkout's own `writer:` label** (the label is defined in **In-Progress Record (Step 2)**, below; derive this checkout's value the same way and compare) → **own crash.** Recover it via substeps 1–3, exactly as before. Only that record counts: `last_completed`, `## Completed This Session`, and `## Still Queued` describe REQs that should not be in `working/` at all, so finding one there is a contradiction to report, not a licence to strip.
- **Named there under a different checkout's `writer:` label** → **foreign claim, and the label says whose.** Leave the file byte-identical and report `claim held by <writer>, not touched`. This case **never enters the three-hour takeover ladder** below: that ladder is for claims nothing accounts for, where age is the only evidence there is, whereas a foreign label is positive evidence that another checkout claimed this REQ — and because the checkpoint travels between checkouts on any install that commits `do-work/`, that other checkout may be building the REQ right now. Clearing an entry that really is dead over there is a human act through `actions/cleanup.md` or an edit by hand, and never recovery's.
- **Named there with no `writer:` label at all** (an entry written before the label existed) → **claim of unknown origin, always report-only.** Leave the file byte-identical and report it as a claim the checkpoint cannot attribute; **no local state of `do-work/CHECKPOINT.md` upgrades it to an own crash.** A dirty checkpoint used to count as evidence that this checkout wrote the entry and has not shared it, and it cannot: under claim-anywhere **every** concurrent claim conflicts on this file (**Worktree Dispatch Mode (Step 1)**, below), so a checkout that merely resolved that merge is holding a modified checkpoint for a reason that has nothing to do with who wrote which entry — and reading it as authorship strips a live foreign claim. Never guess-strip — the whole point of the label is to stop treating "I cannot tell" as "mine." Reclaiming a pre-0.170.0 entry that genuinely is yours stays a human act: the takeover ladder below, or an explicit manual reset — the same path a foreign claim takes.
- **Not named there, or there is no checkpoint at all** → **foreign claim.** Leave the file byte-identical — no frontmatter reset, no section stripping, no move — and report it per *Reporting and takeover* below. An **absent checkpoint is ambiguous, not permission**: a session that died before writing one and a claim made by something else are indistinguishable from here, and only one of the two readings is recoverable when it is wrong. The record is written at **claim time** (**In-Progress Record (Step 2)**, below), so an ordinary crash mid-REQ leaves one and the claims it names under this checkout's label classify as this session's own; a checkpoint missing entirely means the crash preceded this run's first claim, or the claim was never this pipeline's — neither of which is permission.

Recovery consults no lock, because the skill keeps none (**Execution Model — Claim Anywhere, One Releaser**, above). Its two inputs are `do-work/CHECKPOINT.md`'s in-progress record — written at claim time for exactly this purpose (**In-Progress Record (Step 2)**, below) — and each REQ's `claimed_at`, which exists for other reasons. That record is **classification state, not coordination state**: nothing acquires it, nothing waits on it, and another checkout reads it only to classify its entries as foreign and leave them alone.

**Reporting and takeover.** This ladder handles only the claims the checkpoint could not attribute — the unnamed, the no-checkpoint, and the label-less-report-only cases above. A foreign-label entry is already attributed: report `claim held by <writer>, not touched` and go no further, because age adds nothing to a claim you know the holder of. For each remaining foreign claim, compute its age from `claimed_at` (a UTC ISO-8601 instant — Timestamp rule, **Request File Schema — Full Frontmatter** above). **Read `claimed_at` while classifying, not afterwards: substep 1 removes it** — the same ordering trap the `## Scope` / `write_set` decision carries inside that substep.

- **Under three hours** → report and move on, offering nothing: `⚠ REQ-NNN claimed <age> ago, not recorded as this session's work — left untouched.` To reclaim one deliberately, the user resets it by hand; substeps 1–3 below are the reset procedure and run only after the human ownership decision. (Doctor's `STUCK-WORK` one-hour and 24-hour bands are *reporting* severities for a read-only diagnostic, a different purpose from the threshold here, not a second copy of it.)
- **Past three hours, or an unparseable, future-dated, or absent `claimed_at`** → report it and offer takeover. A bad stamp counts as eligible on purpose: a negative or meaningless age has to push toward asking, since the alternative protects a corrupt REQ from takeover forever. Allow **2 minutes of clock skew** before calling a stamp future-dated, matching the skill's other timestamp readers (doctor's `TIMESTAMP-FUTURE` finding, `../../do-work-board/tools/queue-kanban`).

**Three hours bounds how long a dead claim goes unnoticed; it never authorizes anything.** A Route C REQ with a remediation loop can legitimately run longer, so the threshold is not a liveness test and crossing it proves nothing about whether the claim is dead. **The decision to take over is always a human's.** Do not "simplify" this into an automatic takeover at the threshold: an unattended run crossing it on a live claim would then strip exactly the work this classification exists to protect.

Ask with the takeover prompt — `crew-members/clear-questions.md` governs the wording (one decision, options that state their consequence):

```
REQ-042 has been claimed for 4h 20m and is not recorded as this session's own interrupted work.
  (a) Take it over — strip its generated sections (Plan, Exploration, Scope, …) and return it to
      the queue for a clean re-run. Anything the earlier run produced but never committed is lost.
  (b) Leave it claimed — skip it this run and continue with the rest of the queue. Nothing is
      touched; reset it by hand later if it turns out to be dead.
```

Only answer (a) runs substeps 1–3. **With no human to answer** — an unattended or non-interactive run — the outcome is (b): leave the file, report it, continue to the next queue item. Never stall the loop on the prompt, and never resolve a missing answer by stripping.

For each `REQ-*.md` the classification above sent to recovery — an own crash, or a foreign claim a human approved for takeover:
1. Reset frontmatter: set `status` to `pending`, **unless** the REQ file contains a `## Open Questions` section with at least one unresolved `- [ ]` item — in that case, restore `status` to `pending-answers`. (If the `## Open Questions` section exists but all items are already `[x]` or `[~]`, or if no `## Open Questions` section exists at all, set `status` to `pending`.) **Exception — a recovered REQ that already carries `status: blocked` with a `blocked_by` condition stays `blocked`** (the mid-run blocked flip completed its frontmatter write before the crash; its condition is unchanged and it must not be silently promoted to runnable). Remove `claimed_at`, `route`, and all eight optional phase stamps (`planning_at`, `dispatch_at`, `builder_handback_at`, `integration_at`, `review_at`, `remediation_at`, `re_review_at`, `release_at`); a fresh attempt must not inherit observations from the interrupted run. Leave `blocked_by`/`blocked_at`/`blocked_check` intact. **On either flip — to `pending` or to `pending-answers` — stamp `status_changed_at: <now>`** (current UTC instant — Timestamp rule, **Request File Schema — Full Frontmatter** above). Both are status flips with no dedicated `*_at` stamp of their own, so the field's stated trigger condition covers them; and since this substep also removes `claimed_at`, the stamp is the only surviving trace of when recovery happened — without it the board's state timer falls all the way back to `created_at` and dates a just-recovered REQ from the day it was written. **The preserved-`blocked` exception is not a flip and must not stamp** (its `blocked_at` is intact and still correct). **Clear `write_set` only when this REQ actually has a `## Scope` section** — check for it here, while it still exists (substep 2 strips it next). `## Scope` is `write_set`'s only source (**Scope Declaration Template (Step 5.5)**, below), so a mirror that outlives its stripped source is stale; clearing it returns the field to *absent ⇒ unknown*, which the board renders as **no** overlaps badge (not conflict) — the correct post-recovery state, since a recovered REQ has not re-declared its scope. **With no `## Scope`, preserve `write_set`** — it is capture-seeded, not a mirror (the REQ crashed before Step 5.5, or it is a Route A REQ that never runs Step 5.5 at all — `actions/work.md` Step 5.5). Nothing downstream ever re-seeds that field, so clearing it would destroy user- and capture-authored frontmatter.
2. Strip sections generated during the interrupted run: remove `## Triage`, `## Exploration`, `## Plan`, `## Scope`, `## Pre-Flight`, `## Implementation Summary`, `## Qualification`, `## Testing`, `## Review`, `## Lessons Learned`, `## Orientation`, `## Decisions`, and `## Discovered Tasks` sections (and their content) if present — these may be incomplete or stale from the crash. Leave `## Open Questions` and user-authored content intact. (Stripping `## Scope` here is what makes substep 1's `write_set` decision conditional on it — read that decision before this substep runs, never after.)
3. Move the REQ back to `do-work/queue/`

**Worktree sweep — runs once, not per file in the loop above.** Only applies where the run used worktree dispatch (**Worktree Dispatch Mode (Step 1)**, below); skip it entirely when the repo has no `git worktree` support or no `worktree-agent-*` names exist. A leftover branch can outlive its `working/` file (the REQ archived, the branch didn't), which is why this sweeps names rather than iterating the files above. Run `git worktree prune` first so already-deleted directories don't surface as ghosts, then enumerate `git worktree list --porcelain` and `git branch --list 'worktree-agent-*'` — a branch with no worktree and a worktree with no branch each count. For each leftover, read the REQ id out of its `worktree-agent-REQ-NNN-…` name and:

- **Merged into the integration branch** (`git worktree remove <path>` succeeds, then `git branch -d <branch>` succeeds — both run from the integration branch, because `-d` is HEAD-relative; see **Worktree Dispatch Mode (Step 1)**, *Cleanup — happy path*): pure residue. Remove it mechanically, run `git worktree prune` again, and report `Removed merged worktree <name>`.
- **Unmerged or dirty** (either command refuses): **report it and move on — never `-D`, never `--force`.** The branch may hold the only copy of a builder's work; deleting it belongs to `actions/cleanup.md` → **Pass 5: Orphaned Worktrees (consent-gated)**, which asks first. A reported unmerged leftover does **not** block re-dispatch of its REQ — the Naming rule's collision variant (**Worktree Dispatch Mode (Step 1)**, below) dispatches the recovered REQ under a fresh unique variant, so the two coexist until Pass 5 resolves the leftover.
- **Any other worktree name**: not ours. Leave it alone.

Once every `working/` file has been recovered, taken over, or left alone as a reported foreign claim, proceed with finding the next request.

## Repository Gate Deferral and Resumption

This is the full action-owned algorithm behind `actions/work.md`'s baseline and Step 6.5 attribution lanes. The canonical gate is mandatory; deferral changes who owns an unrelated failure and what the selector runs next, never the pass requirement for completion.

### Session state and baseline

At run start hold two session-local sets: **suppressed parents** and **repair closure**. They are scheduling evidence, not REQ fields. A parent enters suppression only from a successful typed `gate_deferral` result and stays there until its repair dependency reaches terminal success. Suppression wins over explicit-REQ provenance, preventing a targeted parent from bypassing the dependency it just gained. Every returned repair id enters the closure even when its `user_request` is outside a targeted UR; this is the only cross-UR widening allowed. Recompute the canonical selector after every deferral and every repair terminal result—never reuse a prior selected record.

After Step 5.75 and before dispatch or source edits, resolve the project-owned canonical gate once as structured argv. Run it directly from the project root and save the current revision, direct status, bounded diagnostic evidence, and a stable semantic fingerprint. Fingerprinting must discard volatile timestamps, scratch roots, and ordering noise but retain the failing command/test identity and normalized diagnostic; use the same procedure everywhere below. A launcher failure produces no comparable fingerprint and stops safely.

The branch table is exhaustive:

| Claimed REQ and baseline | Action |
|---|---|
| Ordinary REQ, exit 0 | Save the green baseline revision; dispatch normally. |
| Ordinary REQ, non-zero | Defer before source edits through the canonical transaction, consume its typed result, suppress the parent, extend repair closure, and select again. |
| `repository_gate_repair: true`, matching recorded red fingerprint | This is the repair's authorized baseline. Implement without recursive deferral. |
| Repair, exit 0 | Complete through the reviewed no-change path; the defect was repaired elsewhere and parents may resume. |
| Repair, different red fingerprint or launch failure | Stop/fail the repair safely; do not manufacture a second repair from it. Parents remain dependency-gated and unrelated ready work continues. |
| Ready `gate_deferred: true` parent without a saved pair | Require a fresh green baseline, then dispatch normally. |
| Ready deferred parent with a saved pair | Run the resume proof below before deciding whether builder work is reusable. |

### Manifest authoring and collision retry

The action authors exact evidence; `do-work-cli defer-gate --manifest <path>` alone mutates. Copy parent/checkpoint preimages to payload files, carry the exact writer label, structured gate argv, direct non-zero status, fingerprint/evidence, stable root-cause `sweep_key`, and optional paired implementation commits. Scan the same request and `.req-reservations/` evidence as capture, propose read-only max+1, and let `defer-gate` exclusively create the unpadded reservation. Do not call a helper that pre-creates a marker. Retry a collision only for one of two typed results: **(a)** pre-mutation `outcome: refused` with the collision finding, an empty `changes` list, and `rollback.status: not_needed`; or **(b)** post-mutation collision with `outcome: rolled_back` and `rollback.status: succeeded`. Rescan the live repository and propose its new max+1 before retrying. An incomplete/failed rollback, `committed_risk`, any non-collision refusal/finding, stale preimage, or non-empty refused-result changes stops without retry.

Fold mode supplies the unique pending repair's exact preimage. A committed repair's SessionStart cleanup may already have removed its reservation, so an absent marker is valid fold topology only when that repair preimage is clean against `HEAD`; a present marker must still match exactly. With a present exact reservation, untracked and manifest-bound tracked-dirty repair folds are supported, but staged repair bytes are always refused. An absent reservation with an untracked or tracked-dirty repair stays refused. An occupied parent queue destination is a pre-mutation planning refusal; exclusive move publication still rechecks absence at the final apply boundary.

On success consume `gate_deferral` fields only: `parent_id`, `parent_path`, `repair_id`, `repair_path`, `repair_outcome`, `repair_dependency`, `diagnostic_fingerprint`, `sweep_key`, command/status, and optional range. Never scrape text rendering or infer the repair from queue order. Validate that the returned parent and fingerprint equal the proposal before updating the session sets.

### Late attribution

The final gate uses the identical argv and fingerprint procedure. Exit zero continues. In worktree dispatch mode, a red current tree is attributed by creating an isolated detached diagnostic worktree at the saved `<pre>` and running the gate there directly. Always remove that diagnostic worktree without force after capturing its status and evidence.

- Active `repository_gate_repair: true`: any same-fingerprint red, different-fingerprint red, missing/malformed fingerprint evidence, or launch-failed final gate is a terminal repair failure. Invoke canonical `fail` to archive it, never invoke `defer-gate`; every parent remains pending behind the failed dependency. Recompute the selector and continue unrelated runnable REQs.
- Base exit 0: current implementation caused the failure; use the bounded remediation loop.
- Base non-zero with the exact current failure fingerprint: unrelated failure; call `defer-gate` with full `<pre>` and `<merge_hash>` commits. After success, the normal archive path will not clean the builder, so immediately run `git worktree remove <path>`, `git branch -d <operative_name>` from the integration branch, then `git worktree prune`. A refusal stops; never add force.
- Base fingerprint mismatch, diagnostic launcher failure, missing/unresolvable range, or inability to isolate the saved base: attribution is unverifiable and stops safely.

Serial mode has no isolated committed implementation range. A late red serial tree therefore keeps the fail-safe stop unless the current diff is demonstrably the cause; it is never submitted as a late deferral with invented commits.

### Saved-range resume proof

Both fields must be present and resolve to commits. Require `base != merge`, base ancestry of merge, and merge ancestry of current `HEAD`. Derive the implementation paths with rename detection from `git diff --name-status -M <base>..<merge>`; for renames, both old and new paths are protected. Reject reuse when any protected path has commit history after merge, or has a current staged, unstaged, untracked, deleted, type-changed, or renamed state. Also reject an unreadable path, ambiguous rename, missing endpoint, side-branch merge, or any other unverifiable evidence.

Valid proof reuses the already-merged implementation but discards old downstream verdicts: rerun qualification over the saved range, focused tests, the canonical gate, and independent review before completion. Drift deletes the two saved pointers from the claimed working REQ, disregards all prior qualification/test/review claims, and returns to Step 6 implementation. An invalid or malformed range stops safely instead of silently rebuilding or archiving; only proven path drift selects the rebuild branch.

### Already-green repair no-op completion

This branch exists only when a claimed `repository_gate_repair: true` REQ's pre-build gate exits 0 before source edits. Append durable evidence using this exact shape, with the gate argv encoded as one JSON array and `Verified at` written by the Timestamp rule:

```markdown
## Repository Gate Repair No-Op

- **Expected diagnostic fingerprint:** <fingerprint recorded by repair intake>
- **Gate command:** ["argv0","argv1"]
- **Direct exit status:** 0
- **Observed result:** green before implementation; repair already satisfied
- **Verified at:** <now> (current UTC instant — Timestamp rule)
```

Write the mandatory summary exactly as:

```markdown
## Implementation Summary

**Files changed:** None — verified repository-gate repair no-op.

**What was done:** Re-ran the repair's recorded canonical repository gate before source edits and confirmed it is already green; no implementation changes were necessary.
```

Qualification is a narrow evidence check, not a vacuous diff pass: require the exact two sections above, rerun the JSON argv directly at exit 0, and prove no project path changed. Do not run the ordinary diff-requiring qualifier. Append `## Qualification\n\nPassed — repository-gate repair no-op; durable gate evidence verified and project diff empty.` Independent review reruns/validates the gate, matches the expected fingerprint to intake, proves the project diff remains empty, checks that no release mutation is planned, and records those facts in the ordinary `## Review`; self-review is insufficient.

After review passes, invoke canonical `complete` normally so the repair archives as terminal success, its checkpoint claim leaves atomically, UR closure/calibration remain canonical, and dependency readiness can resume parents. Skip `release`: write no changelog, version/lock mirror, or `release_at`. Stage only exact lifecycle/archive/calibration and reported UR-move paths; refuse any project, changelog, version, lockfile, or unrelated staged path. The resulting lifecycle-only REQ commit is valid despite the empty project manifest, uses the normal commit format plus `Verified repository gate already green; no implementation changes.`, and becomes the implementation hash recorded by the ordinary separate metadata commit. No other empty implementation receives any of these exceptions.

### Continuation and reporting

A failed, cancelled, or still-gated repair never releases its parents. The selector naturally excludes those parents by `depends_on`; the action continues unrelated selected REQs instead of ending the run. A successful repair causes a fresh selection, where repair priority is exhausted before deferred-parent priority and ordinary work.

Run summaries compose, rather than overwrite: report each deferred parent with its typed repair id and create/fold outcome; each no-change repair completion; each resumed parent as reused or rebuilt-after-drift; each repair failure/cancellation; and every unrelated REQ that continued afterward. The ordinary composed exit summary still reports dependency-gated parents under its blocked-by-dependencies section. No branch writes `blocked` or `pending-answers`, and no branch asks for user confirmation.

## Worktree Dispatch Mode (Step 1)

**Optional, advanced harnesses only.** Each builder runs in its own git worktree on its own branch; the orchestrator stays in the main tree, merges those branches, and remains the only writer of `do-work/` state. Worktree isolation is what makes the ownership boundary hold (**Execution Model — Claim Anywhere, One Releaser**, above): it changes *where* a builder writes — its own tree instead of the main one — so a builder can neither touch queue state nor collide with a sibling, and an interrupted build leaves a branch to merge or discard rather than half-written files in the main tree. **Everything in this section is written per REQ and therefore already holds for any number of concurrent builders** — one `<operative_name>`, one hand-back sequence, one `<pre>..<merge_hash>` range, one cleanup, each per REQ. Fan-Out Dispatch (below) adds only who picks the set and what never parallelises.

**Precondition, then degrade silently.** Probe `git worktree list` (a non-zero exit, or an unrecognized subcommand, means no worktree support) and confirm the harness can run an agent against a working directory you choose. If either is missing, run the serial loop exactly as documented and say nothing further. Unlike `../../do-work-board/actions/board.md`'s Go check — which reports and stops, because there the toolchain *is* the capability — this mode's absence is not an error, so it must never surface as one.

**A builder tree does not have to be a worktree.** The definition a builder has to satisfy is *own tree, own branch, hands back a branch* — a spawned `git worktree`, a second local workspace, a clone, or a remote/cloud sandbox all satisfy it, and everything downstream of the merge is identical for all four because the orchestrator integrates a branch either way. Three deltas are worth stating, because they are the only places the shape matters:

- **The naming and cleanup mechanics below are worktree-specific.** `git worktree add`, `git worktree remove`, `git worktree prune` and the `worktree-agent-REQ-NNN-<suffix>` convention apply to a spawned worktree. A workspace, clone or sandbox the *user* already opened is not the orchestrator's to name or delete; it hands back a branch and the orchestrator merges it under whatever name it arrived with. **Hold that name as this REQ's operative name exactly the same way** (*The name actually created is this REQ's operative name*, below) — the merge, and any reporting, still needs one string.
- **A remote builder's hand-back travels on the branch.** The absolute-main-tree-path hand-back file (*Sole integrator*, below) is a **local-only** mechanism: it works because the builder shares a filesystem with the main tree. A builder that does not — a cloud sandbox, another machine — commits its manifest **on its own branch** and the orchestrator reads it after the merge. Never hand a remote builder a main-tree path; it resolves to nothing or to something else's.
- **A non-releaser checkout's `do-work/` snapshot is potentially stale, and never authoritative.** Where a consumer commits `do-work/`, every non-releaser checkout carries a copy of the queue as of its last sync: a REQ it shows as `pending` may already be claimed elsewhere, and a `status` it shows may be several transitions behind. Read it as a hint about what exists, never as the current state of anything, and resolve disagreements by syncing rather than by writing. This is the same rule *State stays home* (below) applies to a worktree's snapshot, widened to the checkout that owns one.

**Claim conflicts between checkouts are ordinary git conflicts, and `do-work/CHECKPOINT.md` is where they surface.** Two checkouts claiming the same REQ produce a plain **content** conflict on the REQ file — never a rename conflict, because both sides perform the identical `do-work/queue/` → `do-work/working/` rename and git resolves it silently — made entirely of the two claim writes (`claimed_at`, and whatever sections each side generated). With byte-identical claim writes the REQ file does not conflict at all, and the `writer:` label is then the *only* thing that surfaces the double claim. **Expect a `CHECKPOINT.md` conflict on every concurrent claim, including two that overlap in nothing** — two single-line appends land at the same position, so git reports `CONFLICT (add/add)` (`AA`) or `CONFLICT (content)` (`UU`) depending on whether the file already existed. **Resolve it by keeping every entry from both sides.** That is the merge-time reading of the condition Step 10 already carries — *every entry this checkout did not write* survives (**Session Checkpoint Template (Step 10)**, below) — and both one-sided resolutions lose data: taking ours strips another checkout's live claim record by hand, which is the poisoning the label exists to stop; taking theirs discards this checkout's own record and makes its own crash unrecoverable. One id under two labels needs no reconciling: it is the honest record of two checkouts. On the REQ file itself only one claim survives — a human decision, evidenced by which checkout actually has the work. All of that is observed behavior, recorded in `do-work/archive/UR-018/REQ-095-two-clone-acceptance-run.md`'s `## Testing` section, and it holds only where the consumer **commits `do-work/`**; where the directory is untracked nothing syncs, so no conflict surfaces and `duplicate-req-id` is the only detector left.

**Red Flag — a second checkout running the release tail.** Capturing, claiming and building from another checkout are now in contract; the violation to watch for is a second checkout merging integration, bumping the version, prepending to `CHANGELOG.md`, moving files into `archive/`, or closing a UR. Two changelog prepends against one queue is the shape it usually takes, and unique version numbers do not make it safe (**Serial-only**, below).

**Naming.** The worktree directory's basename and the branch name are the **same string**: `worktree-agent-REQ-NNN-<suffix>`. Embedding the REQ id is what lets any later sweep correlate a leftover with its REQ by name alone; sharing one string is what makes a single grep find both the directory and the branch. Derive `<suffix>` from the REQ's filename slug, lowercased and reduced to `[a-z0-9-]`, **as a text operation before you compose any shell command** — never pipe raw REQ text through `tr`/`sed` inside a quoted command line, where an apostrophe breaks the quoting and becomes an injection vector. **On a name collision at creation** — `git worktree add` or branch creation fails because `worktree-agent-REQ-NNN-<suffix>` already exists as a leftover (typically a crash-recovered REQ re-dispatching into the name its own reported-but-not-deleted leftover still holds) — do **not** delete or force. Dispatch under a fresh unique variant: append an **incrementing numeric token** to the suffix — `-2`, then `-3` if that collides too — keeping the `worktree-agent-REQ-NNN-` prefix intact so both names still correlate to the REQ by name (this is what the sweeps grep for). **One scheme, not a choice:** a free pick between a counter and a timestamp token lets two runs shape the same collision differently, so the name stops being readable evidence of what happened. Report the coexistence and leave the original leftover to its owners — the crash sweep above if it turns out merged, `actions/cleanup.md` → **Pass 5: Orphaned Worktrees (consent-gated)** if unmerged.

**The name actually created is this REQ's operative name.** Whatever `git worktree add` succeeded with — the derived name in the common case, the variant after a collision — is the one string every later worktree/branch operation uses: the hand-back merge's `git merge` argument (below), Step 8's `git worktree remove` and `git branch -d` (*Cleanup — happy path*, below), the crash sweep's own-session bookkeeping (**Crash Recovery (Step 1)**, above), and anything reported back to the user. This reference calls it **`<operative_name>`** wherever a command needs it; the worktree path is that same name under the worktrees parent directory. Hold it exactly the way `<pre>`/`<merge_hash>` are held — known from this session's own context and re-typed as a literal into each fresh command, never a shell variable (**Hold both endpoints as re-typed literals**, below). Nothing persists it, because nothing outside this session's run consumes it: the sweeps and `actions/cleanup.md` Pass 5 discover leftover names by enumerating git, never by re-deriving them. **Re-deriving the name from the slug at cleanup time is the failure this closes** — after a variant dispatch the derived string names the *leftover*, so `git worktree remove`/`git branch -d` target unmerged work, refuse, and halt the run on a false "the merge was skipped or lost" alarm while the variant worktree is never cleaned at all. **With no collision the operative name *is* the derived name**, so the common path and serial mode behave exactly as before.

**Where worktrees live: outside the repo working tree.** A sibling directory (`../<repo>-worktrees/worktree-agent-REQ-NNN-…`) or a scratch directory — never nested inside the repo. Inside, it is a second checkout of the repo sitting in the repo: `actions/cleanup.md` Pass 3a scans the filesystem for any `do-work/` directory outside the project root, and where the consumer commits `do-work/` it would find the builder's checkout of one and try to relocate it into the canonical queue. The extra tree also reads as stray untracked residue to every status and stray-file check downstream. That is a corruption path, not just untidiness.

**State stays home.** **Every path under `do-work/` exists in the main tree only and is the orchestrator's** — the queue, `working/`, `CHECKPOINT.md` and `runs/` are examples of the rule, not its extent, so a directory added later is covered the moment it exists rather than when someone remembers to list it. Builders receive their REQ brief in the dispatch prompt and never read or write queue state from inside a worktree. Untracked files do not propagate into a new worktree, so on the common install (where `do-work/` is untracked) a builder simply finds nothing there. The trap is the other install: where a consumer **commits** `do-work/`, the worktree carries a *stale snapshot* of the queue as it stood at the branch point. Treat that snapshot as absent — never read a status from it, never write to it. Every claim, status flip, and archive move happens in the main tree, by the orchestrator.

**Sole integrator.** The builder never writes the main tree or its branch, **with exactly one exception: its own `do-work/runs/work-<YYYY-MM-DD-HHMMSS>/REQ-NNN-handback.md`**, reached by the absolute main-tree path the orchestrator hands it (the repo-relative trap below cuts both directions) and **never staged, committed, or merged** — it is an orchestrator-owned working file, not branch content. That is one path, derived from the builder's own REQ id: a sibling's hand-back, `manifest.md`, anything else under `do-work/runs/`, and every other main-tree path remain violations. The exception exists because the hand-back has to survive a dead transcript (`crew-members/background-agents.md`); without it the mandatory file has nowhere legal to go. It commits on its own branch and hands back its file manifest; the orchestrator merges. A shared file needing one line of wiring — a `<link>` in a shared template, a registry entry, an export barrel — is an **integration seam**: the builder hands back the exact line and where it goes, and the orchestrator applies it in the main tree **inside the merge commit** — step 3 of the hand-back sequence below, which is what keeps the seam inside this REQ's merge range. A builder that edits the seam itself writes the main tree the orchestrator alone owns.

**Merge, never rebase.** Integrate with `git merge --no-ff <branch>` (the full invocation, which adds `--no-commit` so the integration seam has somewhere to go, is step 2 of the hand-back sequence below). Rebasing rewrites the builder's commits, so `git branch -d` no longer recognizes the work as merged and the free merged-ness assertion below is destroyed. `--no-ff` also preserves the merge commit as the "integrated by orchestrator" provenance record.

**When to merge, and the range every evidence step reads.** The orchestrator merges each builder branch **at hand-back — end of Step 6, before Step 6.25 (Implementation Summary)** — so every downstream evidence step (6.25, 6.3, 6.5, 7, 8, 9) observes one integrated main tree, not a clean one. Any position after 6.25 leaves qualify (Step 6.3) and review (Step 7) reading a clean main tree with nothing to check. Run this hand-back sequence on the integration branch — step 0 first, then the four that produce the merge:

0. **Settle the index first, then capture `<pre>` — in that order.** Step 2's claim moves the REQ from `do-work/queue/` into `do-work/working/` and appends to `do-work/CHECKPOINT.md`. Where the consumer **commits `do-work/`** that move is a staged rename sitting in the index, and `git merge` refuses outright: `error: Your local changes to the following files would be overwritten by merge` (exit 2, no merge attempted). First read `git status --short --untracked-files=all -- do-work/` and sort every path into exactly one category: **stage** the Step 2 claim moves, `CHECKPOINT.md`, and the exact owner-written run artifacts (`manifest.md` plus each `REQ-NNN-brief.md`); **allow but never stage** each expected `REQ-NNN-handback.md` named by that run; **stop and surface** every other `do-work/` path as undeclared queue state. The hand-back is scratch the builder writes to the main tree, so excluding it from the stage set must not reclassify it as an error. Check the index too (`git diff --cached --name-only`); any staged path outside the stage category is unrelated work and stops the hand-back. Stage only the exact stage paths with `git add -A -- <exact-bookkeeping-paths>`, re-run the cached-name guard, then use plain `git commit`. **Never use `git commit -- do-work/`:** a path-limited commit takes tracked paths from the working tree and ignores the index content the guard just inspected, so an unstaged edit can be swept in while an untracked bookkeeping file is left behind. If the exact stage set has no changes — normal on a remediation re-merge — skip the commit and continue once the index is empty. Where `do-work/` is untracked, there is no bookkeeping commit. **Step 0 ends with a clean index.** Ordering is load-bearing: this commit must land *below* `<pre>`, so that the next step's capture puts it outside `<pre>..<merge_hash>` — commit after capturing and the owner's own bookkeeping falls inside the merge range, where qualify and review will read it as an undeclared touch of `do-work/`.

1. **Capture `<pre>`** — run `git rev-parse --short HEAD` and read the printed hash. It is the integration tip before this REQ's **first** merge and the lower bound of the merge range. Capture it **once per REQ**; a remediation re-merge does not re-capture it (below). Never recover it afterwards as `HEAD^1` or live `HEAD` — both move as soon as the orchestrator commits the changelog (Step 9) or merges the next sibling.
2. **Guard the queue, then merge without committing** — first run `git diff --name-only <pre>...<operative_name> -- do-work/` (three dots: merge-base to branch tip — which is why this runs **before** the merge; once the branch is merged that diff goes empty). Owner bookkeeping sits below `<pre>` (step 0) and the hand-back file is never committed, so any path printed is queue state committed on the builder's branch — the write *State stays home* forbids. Stop and drop/revert those commits on the branch before integrating (on a remediation re-merge whose fix branch was cut from the integrated tip, owner or sibling commits can surface here too — the cumulative range's safe-direction over-inclusion: judge, don't auto-delete). This guard is the only mechanical check that ever sees them: `tools/checks/qualify.sh` and `tools/checks/scope-drift.sh` both exclude `do-work/` by contract, and `queue-kanban verify`'s committed-queue-state probe reads this same three-dot diff — blind after the merge, and Step 8 deletes the branch before Step 9 runs verify. Then `git merge --no-ff --no-commit <operative_name>` (this REQ's operative name, *Naming* above — the branch the builder was actually dispatched on, which is the collision variant where there was one), then resolve any conflict. `--no-ff` forces a merge commit even where the branch could fast-forward (Merge, never rebase, above); `--no-commit` is what leaves the integration seam somewhere to go. If git answers `Already up to date.` the builder committed nothing: no `MERGE_HEAD` is set, so a `git commit` here would either fail or fabricate a non-merge commit — stop and treat the hand-back as empty instead.
3. **Apply the integration seams, then commit** — stage the handed-back seam lines (Sole integrator, above) and `git commit`. Folding the seam into the merge commit is the only placement that puts it inside the merge range: a seam committed *after* the merge is that merge commit's **child**, hence outside `<pre>..<merge_hash>`, and qualify, review, and Step 9's validation would never see it.
4. **Capture `<merge_hash>`** — `git rev-parse --short HEAD` on the commit just made. It is the upper bound of the merge range and the hash Step 9 writes into the REQ's `commit:` field.

**Hold both endpoints as re-typed literals, never as shell variables.** The canonical [State across command blocks](../docs/prescribed-shell-primitives.md#state-across-command-blocks) rule applies because the consumers sit in later blocks with model round-trips in between — a `"$pre..$merge_hash"` composed in a fresh shell expands to `".."`, which git rejects. Hold both hashes known from this session's own context and re-typed into each fresh command, never a shell variable. `tools/checks/qualify.sh` hard-FAILs on a range it cannot resolve rather than reading an empty diff, so a lost endpoint surfaces as a qualification failure naming the range instead of a vacuous pass.

**The merge range is `<pre>..<merge_hash>`**, and in worktree dispatch mode it replaces the working diff wherever an evidence step reads changes: **Step 6.3** (`tools/checks/qualify.sh` via `DO_WORK_DIFF_RANGE="<pre>..<merge_hash>"`), **Step 7** (review's Get-the-Diff), **Step 8** (post-merge verification, below), and **Step 9's** staged-list validation all consume it. Use the *captured* `<merge_hash>` as the upper bound, never live `HEAD`: HEAD moves the moment the orchestrator commits the changelog (Step 9), so `<pre>..HEAD` would sweep in that commit and misattribute it to this REQ. `<pre>` is `<merge_hash>`'s first-parent ancestor — its direct first parent after a single merge — so `merge-base(<pre>, <merge_hash>) == <pre>` and this is already the merge-base form. Serial mode is unchanged: no range is set and every step reads the working/staged diff exactly as before.

**Remediation re-merges: the range is cumulative.** A failed review sends the REQ back to Step 6 (`actions/work.md`) and the builder's fix branch is merged again. Repeat steps 2–4 above but **not step 1**: keep the first `<pre>` and re-capture only `<merge_hash>` from the newest merge commit, making the range `<pre₁>..<merge_hash₂>` — first pre-merge tip to latest merge — which covers the original work *and* the fix. Re-capturing `<pre>` would instead give `<pre₂>..<merge_hash₂>`, covering only the fix delta: review would read the fix in isolation and every originally-touched file would WARN as listed-but-not-in-the-diff. **Step 9 records the latest `<merge_hash>`** in `commit:`. The cumulative range's one cost is over-inclusion — any orchestrator commit that landed on the integration branch between the two merges falls inside it and surfaces as an undeclared touch for your judgment. That is the safe direction of error: it shows up and gets judged, where under-inclusion silently hides the REQ's own work.

**Post-merge verification, before archive.** The builder verified its own branch; nobody has verified the merged result. Re-run the REQ's acceptance checks — the `## Red-Green Proof` GREEN condition, the project's test command, whatever `## Scope` named — against the merged main tree **before** Step 8 archives anything (`actions/work.md`). A green builder branch must not compose into a red main that archives as done: **the unit you verify is the unit you roll back**, and a red merged state is not an archive-plus-follow-up — stop, revert to the last verified state, and re-dispatch.

**Cleanup — happy path (Step 8).** After the archive move, remove the builder's worktree and branch **by this REQ's operative name** (*Naming*, above — never re-derived from the slug here): `git worktree remove <path>` (no `--force`, where `<path>` is the worktree whose basename is `<operative_name>`), then `git branch -d <operative_name>`, then `git worktree prune`. Run `branch -d` **from the integration branch you merged into**: `-d` tests merged-ness against the current HEAD (or the branch's configured upstream), so from anywhere else a perfectly-merged branch refuses and an unmerged one can pass — "refusal = unmerged" silently becomes "refusal = wrong branch." Both refusals are signal, not friction: `worktree remove` refuses on a dirty worktree (uncommitted builder work you are about to lose), `branch -d` refuses on an unmerged branch (a merge that was skipped or lost). **Never `-D`, never `--force`.** Report the refusal and stop — forcing destroys the only evidence that the integration didn't happen.

**Cleanup — crash path.** Leftovers from an interrupted run are swept by **Crash Recovery (Step 1)** above: already-merged ones are removed mechanically, unmerged ones are reported and never auto-deleted. Discarding unmerged builder work belongs to `actions/cleanup.md` → **Pass 5: Orphaned Worktrees (consent-gated)**, which asks first and only acts when a human can answer.

**Fan-Out Dispatch — several builders, one releaser.** Every guarantee above is already per REQ, so raising the builder count adds no coordination and no durable state.

**Reached two ways, and the default is neither.** `actions/work.md` processes one REQ at a time unless it is asked not to: `do-work run --fan-out [N]` puts it in **auto-wave mode**, where the loop computes the ready set itself and dispatches builders with **no confirmation gate** (*Auto-wave*, below). Without that flag the action runs the serial loop exactly as it always has, and a human or advanced harness can still drive this section by hand — which is what every step here describes, and what auto-wave automates rather than replaces. The floor is why the flag exists instead of the behavior: `actions/work.md` must stay followable by the simplest agent that can read files and run shell commands, so concurrency must be something a reader opts into rather than something sitting in front of them. `--wave N` remains a *scoping* flag — it chooses which dependency depth runs — and composes with `--fan-out`, which chooses how many of the chosen set run at once. The dispatch *mechanism* stays deliberately unspecified (below) in both modes.

What fan-out adds:

- **The set is either picked by a human or computed by auto-wave; `write_set` never gates the first and is not read at all by the second.** In the manual path a human chooses which REQs run together; in auto-wave the loop computes the ready set from `depends_on`, claim state, `assigned_to`, and — when the flag is set — `--skip-impact-negligible` (*Auto-wave*, below, is the canonical condition list; this is a gloss, not a second copy) (*Auto-wave*, below). A REQ's declared `write_set` is **advisory input to that pick, never a gate** in the manual path, and **not read at all** by the computed one — it is display-only, nothing schedules on it, and the board's `overlaps` badge misses glob-vs-glob, `**`, and directory entries, so **absence reads as unknown, not safe** (`../../do-work-board/actions/board.md`). That is why it cannot be a scheduling input: a field whose absence means *unknown* can only ever inform a judgment, never make one.
- **The non-interference proof is the merge, not the pick.** `git merge --no-ff --no-commit` refusing is the only mechanical evidence that two builders' work does not collide. **Its limit is honest: git detects conflicts by line proximity, not meaning.** Two REQs each appending an entry to a shared registry merge cleanly and can still be jointly wrong. The **integration seam** rule (*Sole integrator*, above) is what covers that — and it works only because one integrator applies every seam by hand, inside the merge commit.
- **Integration is serial.** Implementation parallelises; merge → qualify → test → review → changelog → archive runs one REQ at a time. Each merge also invalidates the previous REQ's *Post-merge verification* (above), so those checks re-run per REQ against the tree as it then stands. Expect the wall-clock saving in the build phase only, and say so rather than promising more.
- **A worktree per builder is mandatory, not optional.** Sharing one working tree was considered and ruled out: every test run, qualification check, and review diff would then read a tree carrying the other builder's unfinished edits, so each REQ's evidence steps stop meaning anything and nothing downstream can tell. (The staging race is the lesser problem.) Keep this reason here — without it a later reader re-offers the shared tree as a simplification.

**Auto-wave — what the loop computes, and what it deliberately does not.** `actions/work.md` Step 1 delegates this computation to `tools/do-work-cli.sh --format json next`; its typed `selected` records are the wave and its `excluded` records are the composed reasons. Both collections retain exact `request_path`, original status, and per-record probe/unblock evidence so the owner exact-reads only a returned path for any successful-probe transition before dispatch. The prose below defines that command's contract and must never become a second queue scan. In auto-wave mode (`do-work run --fan-out [N]`) the ready set is every REQ that satisfies **all** of:

1. `status: pending` — every other status is skipped exactly as the serial scan skips it (`actions/work.md` Step 1), holding statuses included.
2. **Dependency-ready** — every id in `depends_on` (or its `dependencies:` alias) resolves to a `completed` or `completed-with-issues` REQ. The same predicate and the same cycle detection as the serial scan; auto-wave adds no second definition of readiness — **including the serial scan's provenance carve-out**: in a targeted run an **explicitly-named** `REQ-NNN` enters the wave regardless of `depends_on` (the user named it outright), while a REQ reached by `UR-NNN` expansion stays gated, scoped to the UR's member set (`actions/work.md` Step 1).
3. **Unclaimed** — not in `do-work/working/`, and carrying no live `claimed_at`. A REQ under another checkout's `writer:` label in `do-work/CHECKPOINT.md` is claimed *there* and is never in this wave (**Crash Recovery (Step 1)**, above).
4. **Not `assigned_to` another session** — the same courtesy read the serial scan performs, skipped and reported the same way (`actions/work.md` Step 1). Explicit targeting still overrides it, and still clears the field on claim.
5. **Not dropped by `--skip-impact-negligible`** — when that flag is set, a REQ whose `impact:` resolves to `impact-negligible` under the Schema Read Contract is out of the wave, skipped and reported exactly as the serial scan skips and reports it (`actions/work.md` Step 1). This is what makes the flag a *which* filter that composes with `--fan-out` rather than one that silently no-ops under it. Explicit targeting still overrides it; `UR-NNN` expansion does not. Without the flag the condition is vacuous and the wave is unchanged.

**Targeting tokens scope the pool, and per-token provenance survives into the wave.** `--fan-out` composes with targeting tokens (`actions/work.md` → Input): in a targeted run the candidate pool is the resolved token set rather than the whole queue, and each of the five conditions applies exactly as the serial scan applies it to that REQ — which is what the carve-outs on conditions 2, 4, and 5 record. The flag never changes *which* REQs run, only how many at once; there is no separate readiness predicate for waves.

**Nothing else enters the computation. In particular, not `write_set`** — it is display-only at any builder count and the wave must not read it, because absence of a declaration reads as *unknown* rather than *safe*, so treating it as a scheduling input would silently under-report contention and produce a wave that looks proven and is not. **The merge is the non-interference proof, not the pick** (above), and that sentence is what makes a computed set safe at all: two REQs that collide are caught when their branches meet, which is the whole fix-at-merge philosophy (**Execution Model — Claim Anywhere, One Releaser**, above). A computed set is therefore *not* a claim that the REQs do not overlap — it is a claim that they are all runnable.

**Bounded, and the bound is not optional.** Size the wave to the harness concurrency limit per `crew-members/background-agents.md` (*Write a manifest per wave; spawn in bounded waves*) — never an unbounded fan-out over the whole ready set. `--fan-out N` sets the bound explicitly; bare `--fan-out` means the harness's own limit, and where that is unknown, **two**, which is the smallest number that is still fan-out. The selector records every ready item outside that bound as `FAN-OUT-LIMIT`; run the selected wave to completion through integration, then invoke `next` again rather than reusing stale records. Recomputing is what lets a REQ whose dependency landed in wave 1 join wave 2.

**No confirmation gate — that is the deliberate change.** Auto-wave dispatches its computed set without asking, which is what "fully automatic set-picking" means and is a change from the manual path's human pick. What it does **not** change: every per-REQ step below (one worktree, one operative name, one hand-back sequence, one merge range, one cleanup), the mandatory run directory and briefs written *before* any spawn, and integration staying serial. Silent degradation also survives — no `git worktree` support, or a harness that cannot run an agent against a chosen directory, means the serial loop runs and nothing is reported as an error (*Precondition, then degrade silently*, above).

**Serial-only — never parallelised, at any builder count:** every `do-work/` queue transition (claim, status flip, archive move); REQ id allocation (`actions/capture.md`); and the project-owned release version file plus `CHANGELOG.md` — one changelog entry per REQ, written by the owner at merge time. Unique version numbers do not make a shared prepend safe.

**The run directory is mandatory here, not optional.** Fan-out is a background fan-out, so `crew-members/background-agents.md` applies (its JIT_CONTEXT already names `work (multi-REQ)`). Its slots map onto this pipeline:

| Guardrail slot | Fan-out use |
| --- | --- |
| run directory | `do-work/runs/work-<YYYY-MM-DD-HHMMSS>/`, created before any spawn |
| per-builder input | `REQ-NNN-brief.md` — REQ body, worktree path, branch name, never-touch list, hand-back format |
| per-builder output | `REQ-NNN-handback.md` — branch, file manifest, integration seams, and **every `##` section the builder would have written into the REQ file** — today `## Discovered Tasks` and `## Decisions`, each under its own heading — because every reader of a builder-authored section reads them from here when the REQ lacks them (**Reading a Builder-Authored Section (any step)**, below). This row and `actions/work.md` Step 6's routed sections are one set: a section Step 6 tells the builder to author and this row does not carry is lost silently. The one main-tree path a builder may write (*Sole integrator*, above) |
| `manifest.md` | REQ id → builder, `<operative_name>`, handback file, landed status — **the orchestrator's**, never written by a builder |
| bounded waves | builders per wave, sized to the harness concurrency limit |

Carry that file's own ceiling note verbatim in spirit: the pattern makes fan-out failures **survivable, not prevented**. Never describe it as a fix.

**The brief must reach the builder as prompt content or an absolute main-tree path — and so must the hand-back path it writes back to.** A repo-relative path resolves inside the worktree, against its own stale tracked copy of `do-work/` (*State stays home*, above) — so the builder silently reads a snapshot, or nothing, instead of its brief. The same resolution applies in the return direction, where it is worse: the write succeeds, lands in the builder's branch, and the orchestrator reads nothing.

**Dispatch mechanism is deliberately unspecified.** The owner synthesizes from files on disk, never from conversation, so a spawned subagent and a human-opened session are indistinguishable to it. Do not document two routes.

## Reading a Builder-Authored Section (any step)

**Whenever you read a `##` section the builder authors, read the REQ file first and this
REQ's hand-back second, and treat what you find in either as the section.** In worktree
dispatch mode the builder cannot write the REQ file at all — the REQ lives in the main tree,
which **Worktree Dispatch Mode (Step 1)** → *State stays home* forbids it to touch — so it
routes those sections to its hand-back instead. Read
`do-work/runs/work-<YYYY-MM-DD-HHMMSS>/REQ-NNN-handback.md` for a local builder, the merged
branch content for a remote one. **That path is relative to the project root**, where
`do-work/` lives whether the suite is vendored under `.claude/skills/` or checked out
whole — never relative to the directory this action file sits in.

**The condition carries the rule, not any list of readers.** `actions/work.md` Step 8's
discovered-tasks substep, `actions/review-work.md` Step 4's traceability check, and the
**Decision Brief (hand-back format)**'s HANDLED block are the readers that exist today;
they are illustrative. A step, action, or report added later that reads a builder-authored
section inherits this without being remembered here.

**Absence is only silence when you know you looked.** In worktree dispatch mode, if the
section is in neither place **and this REQ's hand-back is missing or unreadable**, say so
rather than proceeding as though the builder recorded nothing: `⚠ REQ-NNN: no <section> in
the REQ and no readable hand-back at <path> — anything the builder recorded there is lost.`
A hand-back that exists and simply has no such section is a real "the builder recorded
nothing" and reports nothing. Every reader states which of the two it found — an unread
hand-back and an empty one are different facts and must never render the same. Serial mode
reads the REQ file alone and this paragraph does not apply.

## Composed Exit Summary (Step 1)

**Exit paths when the scan finds nothing to claim:**

The exit report is **composed**, not picked from disjoint branches. Whenever the scan finds no claimable `pending` REQ, lead with the headline that matches the actual queue state — `No pending REQs in queue.` when the queue holds no `pending` REQs at all, `No dependency-ready pending REQs.` when `pending` REQs exist but every one is dependency-blocked, `No claimable pending REQs — every ready one is assigned to another session.` when dependency-ready `pending` REQs exist but every one carries `assigned_to`, `No claimable pending REQs — every ready one is impact-negligible and --skip-impact-negligible is set.` when the flag dropped every otherwise-claimable REQ (in each case the matching section below then enumerates them, so the headline never strands the user). Then append every section that has at least one REQ — **that condition is the rule, and the list below is the set as it stands today**, so a section added later inherits it without anyone re-counting. In this order:

1. **Finished-awaiting-archive section** — applies if any REQ in `do-work/queue/` normalizes under the Schema Read Contract to a terminally resolved status. Read the `user_request` frontmatter field from each to group by UR. Render:

   ```
   ⚠ N finished REQs awaiting archive (UR-137: 3 REQs, UR-138: 1 REQ, ...):
     REQ-351 — [title] (complete → completed)
     REQ-352 — [title] (completed)
     REQ-353 — [title] (cancelled)
     ...

   Run `do-work cleanup` to archive completed work, then `do-work recap` to see full history.
   ```

2. **Pending-answers section** — applies if any REQ has status `pending-answers`. Render from frontmatter only — do not open the REQ body to count `## Open Questions` items at this stage (Step 1 reads frontmatter per the queue scan). The count is deferred to `do-work clarify`, which is the action that reads Open Questions sections:

   ```
   ⚠ N REQs awaiting clarification:
     REQ-NNN — [title] (pending-answers)
     ...

   Run `do-work clarify` to batch-review the open questions; resolved REQs flip to `pending` and re-enter the queue.
   ```

3. **Blocked-on-external-condition section** — applies if any REQ has status `blocked` (waiting on an external condition named in `blocked_by` — a service being up, a person answering, credentials provisioned — not user answers and not another REQ). Render from frontmatter: the `blocked_by` condition, the age from `blocked_at` (now − `blocked_at`), and whether an auto-probe is configured or failed this run. Step 1 already re-ran each `blocked_check` probe before composing this summary, so a REQ that still appears here either has no probe or its probe did not pass this run:

   ```
   ⚠ N REQs blocked on external conditions:
     REQ-NNN — [title] (blocked by: <condition>, since <age>) [probe failed this run | no auto-probe]
     ...

   When a condition is satisfied, re-run `do-work run` (REQs with a `blocked_check` are re-probed automatically and unblock on exit 0),
   or confirm a human-checkable one via `do-work clarify`. To give up on one, `do-work abandon REQ-NNN`.
   ```

   **A blocked REQ carrying `stakeholder:` renders in the stakeholder form instead** of the plain line above, led by `⚠ IRREVERSIBLE` when K > 0:

   ```
     REQ-NNN — questions for <stakeholder> (N open, K irreversible; since <age>) — report: <latest bundle path from blocked_by>
   ```

   and the section's remedy text gains: `To ingest a stakeholder's reply, run do-work stakeholder-answers REQ-NNN — share the report path with them first if you haven't.` Counting N and K is the one bounded exception to this summary's frontmatter-only stance: open the body of stakeholder REQs only — at most one per stakeholder, by construction — never the rest of the queue; if the body cannot be read, fall back to the plain blocked line for that REQ.

4. **Blocked-archive-collision section** — applies if any REQ has status `blocked-archive-collision`. Read the matching archive path from each blocked REQ's frontmatter if recorded; otherwise re-run the Step 2.0 glob (`do-work/archive/**/REQ-NNN-*.md` and `do-work/archive/**/REQ-NNN.md`) to find it. Render:

   ```
   ⚠ N REQs held by archive-collision guard:
     REQ-NNN — [title] (queue file: do-work/queue/REQ-NNN-slug.md)
       already archived at <archive-path>
       recover: rename the queue file (if this is an intentional re-do) or delete it (if it's a stale duplicate), then flip status back to `pending`
     ...
   ```

5. **Blocked-by-dependencies section** — applies if any `pending` REQ has an unmet `depends_on` reference (dependency-blocked) or any REQ has status `blocked-dependency-cycle`. Pending REQs stay `pending` (the gating is dynamic — they become ready as upstream REQs complete); only cycle-detected REQs are flipped to a held status. Render both groups under one heading:

   ```
   ⚠ N REQs blocked by unmet dependencies:
     REQ-NNN — [title] (pending; depends on REQ-MMM, status: <pending|claimed|pending-answers|failed|cancelled>)
     REQ-PPP — [title] (blocked-dependency-cycle; chain: REQ-PPP → REQ-QQQ → REQ-PPP)
     ...

   Resolve the blocking REQs first, then re-run. To force a scoped run that ignores dependency gating for a specific REQ, use `do-work run REQ-NNN`. To break a dependency cycle, edit the REQ's `depends_on` and flip its status back to `pending`. A dependency on a `cancelled` (or `failed`) REQ never self-resolves — re-point the dependent's `depends_on`, or abandon it too (`do-work abandon REQ-NNN`).
   ```

6. **Assigned-elsewhere section** — applies if any `pending` REQ carries a non-empty `assigned_to` (earmarked for another session — **Request File Schema — Full Frontmatter**, above). These stay `pending` and are not a held status: the field is advisory, and the same REQ becomes claimable here the moment a user names it explicitly. Render the value verbatim, never normalized:

   ```
   ⚠ N REQs assigned to another session:
     REQ-NNN — [title] (assigned to <assigned_to, verbatim>)
     ...

   Skipped as a courtesy, not blocked — nothing confirms that session is running. To take one over here, name it explicitly (`do-work run REQ-NNN`), which clears the assignment as part of the claim. To drop an assignment without running it, remove the field by hand.
   ```

7. **Skipped-as-negligible section** — applies only when `--skip-impact-negligible` is set, and then to every otherwise-claimable `pending` REQ whose `impact:` resolves to `impact-negligible` (**Request File Schema — Full Frontmatter**, above). These stay `pending` and are not a held status: the flag reads the field and writes nothing, so re-running without it picks the same REQs straight back up. This section also renders on the targeted exit path (`actions/work.md` Step 1 → Targeted mode) when the flag dropped resolved members — scoped there to the resolved token set, never the whole queue. Render the resolved token, not the raw one:

   ```
   ⚠ N REQs skipped as impact-negligible:
     REQ-NNN — [title] (impact-negligible)
     ...

   Skipped by `--skip-impact-negligible`, not blocked — nothing was written to these REQs. Re-run without the flag to include them, or name one explicitly (`do-work run REQ-NNN --skip-impact-negligible`), which overrides the skip for that REQ. A REQ with no `impact:` reads as `impact-user-visible` and never appears here.
   ```

**After rendering all applicable sections, exit the work loop** — do not proceed to Step 2.0 or beyond. There is no claimable `pending` REQ. Step 1's contract on this path is "render the composed summary, then stop"; the only path that continues is the one where Step 1 finds at least one claimable `pending` REQ (dependency-ready and, in a default scan, not assigned elsewhere).

If **no section applies** (no REQs at all in `do-work/queue/`), report completion and exit. Never silently exit when any section applies — every non-pending or non-ready REQ in the queue is something the user needs to see.

**Composition is deliberate.** A queue with both `pending-answers` and `blocked-archive-collision` REQs (and no completed/done) renders both sections back-to-back. A queue hitting every category renders every section. The user sees the full picture in one report instead of a single branch's slice.

## In-Progress Record (Step 2)

**Canonical lifecycle transaction boundary.** `do-work-cli claim`, `unblock`, `complete`, `fail`, and `cancel` are the sole writers for deterministic request-state transitions. They consume one repository snapshot plus an exact caller-supplied request path, plan every request/checkpoint/archive/UR/calibration target, and apply once through the shared Git transaction. The action supplies provenance and human judgment: confirmation/reason/dependent disposition for cancel, successful-probe or confirmed-human evidence for unblock, terminal status and known implementation hash for complete, and classified error/error type for fail. A command refusal is byte-identical and actionable; a missing or failed command stops the lifecycle operation. The descriptive field/archive rules in this reference define semantics and compatibility, not a free-form fallback implementation.

**Completion/cancellation archive semantics (declarative).** A standalone or legacy request leaves `working/`/`queue/` for archive root. A UR member moves into its archived UR folder only when the transaction's projected snapshot says every member is terminally resolved; the same transaction consolidates the exact sibling/input paths and removes the active UR directory. A failed/nonterminal sibling keeps the UR open and the just-resolved request archives at root. An already-archived failed cancellation stays at its exact path, and an eligible review follow-up may enter an existing archived UR folder without reopening or relocating that folder. Calibration is derived from the just-archived bytes and appended by `complete`, or returned as typed skipped work. These are postconditions for interpreting the command result; actions never execute a second move, collision check, UR consolidation, or calibration append.

**Archived review-follow-up semantics (declarative).** When and only when a completed REQ carries `review_generated: true` and an archived `UR-NNN` folder already exists for the same `user_request`, canonical `complete` plans the REQ into that existing folder in place. The archived UR folder is never moved, reopened, or re-consolidated, and this case bypasses the normal active-UR closure branch. These are planner postconditions, not action-side move instructions.

**Calibration evidence semantics (declarative).**

**At calculation time, read both `claimed_at` and `completed_at` from the just-archived REQ file's frontmatter; never reuse either stamp from a value held in context earlier in the run.**

The canonical `complete` command either appends the planned row once or reports typed skipped work.

`do-work/CHECKPOINT.md`'s `## In Progress (interrupted)` section is **Crash Recovery's classification input** (**Crash Recovery (Step 1)**, above): a `working/` REQ named there **under this checkout's own writer label** is this session's own interrupted work and recovers; **any other entry is left byte-identical, however it fails to match that label.** The own-label condition is the whole test, and the ways an entry can miss it are enumerated once — in Crash Recovery — not restated here. So the record has to exist at the moment a crash happens. Step 10 writes the checkpoint at *session end*, which a hard crash never reaches — leaving that as the only write site made the own-crash branch unreachable by the very event it handles, and every crashed REQ stranded in `working/` (the 0.164.0 regression this procedure closes). **The record is therefore written at claim time**, by `actions/work.md` Step 2, as part of the claim.

**It is a classification input, and nothing else.** It grants no exclusivity and coordinates nothing — the skill keeps no lock, heartbeat, or liveness check (**Execution Model — Claim Anywhere, One Releaser**, above), and this record must never grow into one. Nothing acquires it and nothing waits on it; another checkout reads it only to classify its entries as foreign and leave them alone. Its whole job is to let the *next* run tell this session's leftovers from someone else's. One small file write per claim, one removal per departure; if it ever acquires a **refresh interval**, a **staleness check**, or a **liveness probe**, it has become the machinery this model exists without. **The static `writer:` label below is none of those** — it is written once at claim time from two values that never change for this checkout, is never refreshed, and is never read as evidence that anything is still running; it records *who wrote the entry*, which is the one question classification actually asks.

**Writing it** — immediately after Step 2's frontmatter flip to `status: claimed`, so the record and the claim land together:

- **The record is a list.** Append one entry per claimed REQ — `- REQ-NNN: [title] — claimed <claimed_at> — writer: <hostname>:<absolute-checkout-path>` — and never collapse the section to a single id. Fan-out dispatch claims several REQs concurrently under one owner (**Worktree Dispatch Mode (Step 1)** → *Fan-Out Dispatch*, above), so a singular record would classify every claim but the newest as a foreign claim after a crash — the same silent stranding this record exists to prevent, just narrower.
- **The `writer:` label names the checkout that wrote the entry, and both halves are load-bearing.** `<hostname>` from `hostname -s` (plain `hostname` where `-s` is unsupported), `<absolute-checkout-path>` from `git rev-parse --show-toplevel` or, outside git, the project root's absolute path. Neither half identifies a checkout alone: two machines can both hold `/home/user/repo`, and one machine can hold several checkouts. The label exists because **the checkpoint is an ordinary tracked file wherever a consumer commits `do-work/`** — another checkout's live claim then arrives here by a routine `git pull`, looking exactly like a local one, and without the label recovery reads it as this session's own crash and strips a REQ someone is actively building.
- **Append; never rewrite what is already there.** A sibling claim's entry is not this claim's to restate or reorder, and **an entry carrying another checkout's label is never rewritten, reordered, or removed — only read.** **One entry per REQ id per writer**: if this REQ already has an entry under *this* checkout's label — a re-claim in the same session, after a blocked flip released it or a mover skipped its removal — refresh that entry in place with the new `claimed_at` instead of appending a second. Two own entries for one id would survive the first removal and leave a permanent phantom claim. A foreign-label entry naming the same id is not a duplicate to reconcile: one id under two labels is the honest record of two checkouts, and reconciling it is nobody's job here.
- **No `do-work/CHECKPOINT.md` yet** — the ordinary case on a fresh run, since Step 10's session-start note deletes it once recovery has finished with every `working/` file (`actions/work.md` Step 10). Create the file containing **only** the `# Session Checkpoint` heading and this one section. Step 2 writes no frontmatter and no other section: Step 10 still owns the full checkpoint and rewrites the file wholesale at session end.
- **Remove this checkout's own entry whenever the REQ leaves `do-work/working/` — that condition is the rule, and the movers below are illustrative, not the set.** Anything that relocates a `working/` REQ drops its entry as part of the same move: today that is Step 8's archive move on success and on failure, and the mid-run blocked flip's move back to `do-work/queue/` (`actions/work.md` Step 8), plus the recovery movers — `actions/cleanup.md` Pass 0's terminal-in-`working/` sweep and this section's human-approved reset. A new mover inherits the rule from the condition; do not extend this list and expect it to stay complete. The removal is part of the move, not a later sweep: a REQ listed as in-progress while sitting somewhere other than `working/` is exactly the contradiction the own-crash bullet tells the next run to report (**Crash Recovery (Step 1)**, above), and a report that fires on normal completion is noise that trains readers to ignore it. Removing the last entry leaves the section present and empty — correct, and distinguishable from a session that never wrote one. A foreign-label entry naming the same REQ stays exactly where it is: this move says nothing about what another checkout is holding. **A label-less entry leaves with the REQ when a human reclaims it** — recovery refuses to attribute one (**Crash Recovery (Step 1)**, above), so the decision to reclaim *is* the authorship call, and an entry that outlives the claim it records is a phantom that `actions/work.md` Step 10's delete gate re-reports every session with no exit.
- **Reach the file by its literal path**, `do-work/CHECKPOINT.md` from the project root. Per the canonical [State across command blocks](../docs/prescribed-shell-primitives.md#state-across-command-blocks) rule, never carry it (or a `mktemp` path) in a variable inherited from an earlier block.

## Triage Section Template (Step 3)

```markdown
---

## Triage

**Route: [A/B/C]** - [Simple/Medium/Complex]

**Reasoning:** [1-2 sentences]

**Planning:** [Required/Not required]
```

## Plan Template — Route C (Step 4)

```markdown
## Plan

[Plan agent output]

*Generated by Plan agent*
```

## Plan Skip Note — Routes A/B (Step 4)

```markdown
## Plan

**Planning not required** - [Route A: Direct implementation / Route B: Exploration-guided implementation]

*Skipped by work action*
```

## Scope Declaration Template (Step 5.5)

```markdown
## Scope

**Files I will touch:**
- `src/stores/theme-store.ts` (new) — theme state management
- `src/components/settings/SettingsPanel.tsx` (modify) — add toggle
- `tests/theme-store.test.js` (new) — unit tests

**Files I will NOT touch:** [any files that seem related but are out of scope]

**Acceptance criteria (restated from REQ):**
- [ ] Dark mode toggle visible in settings
- [ ] Theme persists across page reload
- [ ] OS preference respected on first visit
```

**"Files I will touch" is the source of the `write_set` frontmatter field.** After writing this section, the orchestrator mirrors the list into `write_set:` — one direction only, so the prose and the field cannot drift. Never edit `write_set` and expect the Scope list to follow. The mirror feeds the board's overlaps badge only (`write_set` is display, not scheduling, at any builder count — **Worktree Dispatch Mode (Step 1)** → Fan-Out Dispatch, below).

## Pre-Flight Template (Step 5.75)

```markdown
## Pre-Flight

**Git:** ⚠ 3 pre-existing uncommitted files — preserve/exclude from this REQ and account for them in qualification/review evidence (src/temp.ts, .env.local, notes.md)
**Tests baseline:** ✓ All passing (47 tests)
**Dependencies:** ✓ Installed

*Checked by work action*
```

## Implementation Summary Template (Step 6.25)

```markdown
## Implementation Summary

**Files changed:**
- `src/stores/theme-store.ts` (new)
- `src/components/settings/SettingsPanel.tsx` (modified)
- `tests/theme-store.test.js` (new)

**What was done:** [1-2 sentences — what the implementation actually did]
```

## Qualification Anti-Rationalization Table (Step 6.3)

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "The summary says files changed" | Check the file system | The summary is a claim, not evidence |
| "Tests pass so requirements are met" | Compare requirements to diff, word by word | Tests can be incomplete |
| "The builder checked the UNIFY box" | Read the actual diff for debug artifacts | A checked box is a claim, not a fact |
| "This works on my test case" | Test at least 2 additional cases including an edge case | One test case proves nothing about generality |
| "The existing code was already like this" | Flag it in Discovered Tasks | Pre-existing problems are still problems |
| "It's just a small deviation from the plan" | Log it as a Decision (D-XX) | Unlogged deviations break traceability |

## Testing Section Template (Step 6.5)

```markdown
## Testing

**Tests run:** [command]
**Result:** ✓ All passing (X tests)

**Red-green validation:** *(for bug fixes and new features)*
- [test name/file]: ✗ before implementation → ✓ after
- [test name/file]: ✗ before implementation → ✓ after

**New tests added:**
- [list]

**Existing tests updated (cross-REQ impact):**
- [test file] (from REQ-NNN): [what changed and why — intentional behavior change]

*Verified by work action*
```

## Finding-Closure Ratchet (Step 6.5)

**A review- or triage-finding-origin REQ closes only when its captured GREEN names a regression test/check that fails before the fix and passes after, or when the exact named finding surface is deleted.** Closure evidence must match that named check or deletion surface. A bare patch, unrelated green tests, `tdd: false`, and a high review score are not closure evidence.

## Deferred Prime-Link Path Computation (Step 7.5)

**Path computation rule (for use in Step 8):** the lesson bullet is written to the prime's satellite `lessons-<name>.md`, which sits in the prime's own directory, so the link path is relative to that directory — not the repo root. Count how many directories deep the satellite sits (i.e., the number of path components before the filename). Prepend that many `../` steps to the REQ's repo-root-relative archive path. Examples:
- Satellite at `lessons-auth.md` (0 dirs deep) → `do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`
- Satellite at `src/utils/lessons-auth.md` (2 dirs deep: `src/` and `utils/`) → `../../do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`
- Satellite at `web/src/auth/lessons-auth.md` (3 dirs deep) → `../../../do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`

The existence-verify check on the resolved path runs in Step 8 (post-move) — that's the whole reason for deferring. A path that does not resolve because the satellite ships in a package whose consumers never receive `do-work/archive/` takes the canonical repository URL instead; anything else that fails to resolve is reported, never written.

## Builder-Decided Follow-up Template (Step 8)

This template is the *session-user* branch of Step 8's audience fork: a `- [~]` record carrying `Answerer: <name>` never lands here — it routes to that person's stakeholder REQ instead (**Stakeholder REQ Template (Step 8)**, below).

Before creating one, run the fold-first scan (`actions/capture-reference.md` → **Fold-First Rule**) — a pending REQ in any UR sharing the root cause receives the follow-up as an appended instance instead of a new file. Its `created_at` is the current UTC instant (Timestamp rule, above).

   ```markdown
   ---
   id: REQ-NNN
   title: '[<impact token>] Confirm: [brief description of the choice]'   # quoted per **Frontmatter Quoting**, above; omit the tag when impact: is the impact-user-visible default (actions/capture-reference.md → REQ Title Convention)
   status: pending-answers
   created_at: <timestamp>
   user_request: [same UR as the original REQ]
   addendum_to: [original REQ id]
   builder_decided: true
   ---

   # Confirm: [Brief Description]

   ## What
   During [REQ-id], the builder chose [choice] for [question]. This follow-up
   confirms whether that choice matches your intent or if you'd prefer a different approach.

   ## What the Builder Chose
   [Description of the choice and its impact on the implementation]

   ## What Would Change
   [If the user picks a different option, what would need to change]

   ## Open Questions
   - [ ] [Original question]
     Recommended: [builder's choice — already implemented]
     Value: [what this choice buys — copied from the D-NN record]
     Risk: [what breaks if it's wrong, and how reversible — copied from the D-NN record]
     Also: [other alternatives]
   ```

   The `Value:`/`Risk:` lines come from the escalated `D-NN` entry's record (work.md Step 3.5/6). They let `do-work clarify` render the **DECISIONS FOR YOU** section of the Decision Brief so the user can judge in seconds. If the original decision was logged without them (older REQ), omit both lines — `clarify`'s fallback renders `Recommended:`/`Also:` alone.

   All template text the user will read — the question, `Recommended:`/`Also:`/`Value:`/`Risk:`, and the `## What` section — must satisfy `crew-members/clear-questions.md`: self-contained, no spec-internal shorthand or coined labels without a one-line gloss, and the why-this-was-escalated stated in `## What` or the question itself (Principle 7). The builder writing this template has the full spec in context; the user answering it in a later clarify session does not.

## Stakeholder REQ Template (Step 8)

The *stakeholder* branch of Step 8's audience fork mints one of these — but only when no open REQ for the person exists: the Fold-First Rule's **Stakeholder-audience questions** clause (`actions/capture-reference.md`) appends to an existing one first. One open REQ per stakeholder, by construction; when it archives (`actions/stakeholder-answers.md` → **Stakeholder REQ terminal semantics**), the next question for that person starts a fresh one at Q-01.

Q-NN ids are unique within one REQ and never reused — an answered or reclaimed question keeps its line and id forever, because the id is the routing key a reply names (`actions/stakeholder-answers.md` Step 2). The counter comment mirrors the D-XX counter pattern (`actions/work.md` Step 3.5). Both `<timestamp>` stamps below are the current UTC instant (Timestamp rule, above).

   ```markdown
   ---
   id: REQ-NNN
   title: 'Stakeholder questions: Priya (design)'   # stable across folds; no impact tag (the field is absent and reads as the default)
   status: blocked
   created_at: <timestamp>
   # deliberately NO user_request: — UR membership would make this REQ hold the first source UR open (the archive table treats every blocked member as unresolved), and nothing waits on this REQ. Per-question provenance lives in each entry's Source: line.
   stakeholder: 'Priya (design)'   # marker + fold discriminator — verbatim-read class, single-quoted (**Frontmatter Quoting**, Request File Schema above)
   blocked_by: 'answers from Priya (design) — latest report: ai-reports/<bundle-slug>/index.html'
   blocked_at: <timestamp>
   ---

   # Stakeholder Questions: Priya (Design)

   ## What
   Open questions whose answerer is Priya (design), collected from builds that
   proceeded on best-judgment assumptions. **Nothing waits on this REQ** — every
   question below was answered by an assumption and its source REQ completed.
   This REQ exists to route Priya's confirm-or-override answers back
   (`actions/stakeholder-answers.md`).

   ## Questions
   - [ ] **Q-01** — [the question, complete for a cold outside reader]
     Assumed: [the builder's implemented answer]
     Value: [carried from the D-XX record]
     Risk: [carried from the D-XX record, including reversibility]
     Irreversible: yes — [why undoing is expensive]     <- only when the record has it
     Source: REQ-NNN (D-04)

   <!-- Q-NN counter: last used Q-01. Next question: Q-02. -->

   ## Reports
   - [<date>] ai-reports/<bundle-slug>/index.html — 1 open question at generation

   ## Blocked
   - [<date>] blocked on "answers from Priya (design)" — minted by REQ-NNN's archive step
   ```

   Answered forms (Canonical answered-question format, `actions/clarify.md` Step 4): `- [x] **Q-NN** — [question] → [answer]`, `→ Confirmed: [assumed answer]`, or `→ Reclaimed by user via clarify — [answer | moved to REQ-MMM]`. Per-question state lives in the checkbox and its id; REQ-level state lives in `status` alone — the two never stand in for each other. The heading is `## Questions`, deliberately **not** `## Open Questions`: every existing Open-Questions reader (work Step 3.5's scan, clarify's presentation, crash recovery's restore rule) must walk past this REQ untouched.

## Discovered Tasks Classification (Step 8)

   The builder should classify each discovered task when appending them, using the impact vocabulary (**Request File Schema — Full Frontmatter** above) — the same four tokens the review flow records, so the token a builder writes here is the token the follow-up's `impact:` field carries:
   ```
   ## Discovered Tasks
   - **impact-critical** SQL injection vulnerability in user search endpoint
   - **impact-user-visible** The retry banner never clears after a successful reconnect
   - **impact-rule-change** Three adapters hand-roll the retry loop the shared client already provides
   - **impact-negligible** Variable naming inconsistency in auth module
   ```

   If the builder did not classify them, classify each with the two questions in `actions/review-work.md` Step 10 — that step is the single home of the rubric, and a second copy here would drift from it.

   **Fold before creating:** each classified discovery first runs the fold-first scan (`actions/capture-reference.md` → **Fold-First Rule**) — a `pending` or `pending-answers` REQ in any UR sharing the root cause receives it as an appended instance, and a prose-only discovery with no match lands on `do-work/prose-backlog.md`; an append needs no consent question because editing an existing queued REQ is not creation (`actions/review-work.md` Step 10 → the reroute governs creation only). Only a discovery with no destination in the ladder proceeds to the creation paths below.

   **Set `impact:` on every follow-up REQ this classification creates** (any token, either status path) to the discovery's recorded token, verbatim — the token IS the field value, and its title carries the matching `[<impact token>] ` tag under the same rule every other minted title follows (`actions/capture-reference.md` → **REQ Title Convention**; a follow-up that carries the field but not the tag is invisible to the board's title search, which is the whole reason the tag exists). `effort_estimate` is the other axis and is judged separately, as size: judge the fix and emit either `effort-mechanical` or `effort-substantive`; when the size is genuinely unclear, put that judgment to the user. Omit the field only when neither judging nor asking was possible, and never copy a default in either direction. Impact and effort are different axes — an `impact-negligible` style sweep can be substantive effort, and an `impact-user-visible` one-line fix can be mechanical.

   **For `impact-critical` discoveries:** Create follow-up REQ with `status: pending` (not `pending-answers`) — these skip user confirmation and go straight into the work queue. Add a note in Open Questions: `- [x] Auto-approved: critical severity (security/data/production risk). → Added to queue immediately.` Report prominently: `⚠ CRITICAL discovered: [description] — auto-queued as REQ-NNN`

   **Test-hygiene carve-out:** A non-`impact-critical` discovery ALSO auto-queues with `status: pending` (same as critical) when ALL three hold:
   - The fix touches **only test files** (`tests/**`, `*.test.*`, `*.spec.*`, test helpers/fixtures) — zero production-source changes.
   - It is **mechanical hygiene** — silencing warnings/console noise, deflaking, lint/format cleanup in tests — with no behavior or assertion-meaning changes.
   - It is **small** — a single file or a couple of files, no new infrastructure.

   Failing ANY bullet keeps the `pending-answers` flow below. The paper trail mirrors the critical flow: add a note in Open Questions: `- [x] Auto-approved: test-only mechanical hygiene ([impact token]). → Added to queue.` Report visibly: `↺ test-hygiene discovery auto-queued as REQ-NNN`. The user can still discard the REQ from the queue afterwards — that stays the escape hatch.

   **For every non-critical discovery (when the test-hygiene carve-out does not apply):** Use the existing `pending-answers` flow:
   - Set frontmatter: `status: pending-answers`, `user_request: [same UR]`, `addendum_to: [current REQ id]`, `domain: [same domain as current REQ]`.
   - Add an `## Open Questions` section with this checkbox format:
     `- [ ] I discovered this out-of-scope task while working on [current REQ]: [Task Description]. Should I process this as a new task?`
     `  Recommended: Yes, add to queue (will flip to 'pending').`
     `  Also: No, discard it.`
   This ensures non-critical discoveries — other than qualifying test-only hygiene — require the user's explicit permission via `do-work clarify` before execution.

## Failure Classification (Step 8)


This classification runs at any generation: `review_generated: true` on the failed REQ does **not** suppress its failure follow-up — the cascade depth stop (`actions/review-work.md` Step 10 → **Generation ≥ 2**) governs review-*finding* follow-ups only, and a failed generation-≥2 REQ with no successor would die silently.

Before classifying via the symptom table below, **check for upstream failure**. Cascades from a failed prerequisite often present as plausible-looking `code` or `spec` symptoms in the downstream REQ; misclassifying them sends the builder chasing phantom bugs in the wrong domain.

**Upstream-failure short-circuit:**

Read the frontmatter of every REQ this one depends on:
- `addendum_to` (single parent, if set; or `amends`/`parent`/`amendment_to` as the legacy alias if `addendum_to` is absent — same back-compat shape as the `depends_on`/`dependencies:` pair; `addendum_to` wins when multiple are present)
- every entry in `depends_on` (if set, or every entry in the legacy `dependencies:` alias if `depends_on` is absent — same back-compat rule as Step 1; `depends_on` wins when both present)

Resolve each ID by globbing `do-work/archive/**/REQ-NNN-*.md`, `do-work/archive/**/REQ-NNN.md`, `do-work/queue/REQ-NNN-*.md`, and `do-work/working/REQ-NNN-*.md`. If any referenced REQ has `status: failed`, skip the symptom table and short-circuit classification:

- `status: failed`
- `error_type: spec` (the local approach is downstream-correct only if the upstream is correct; with the upstream broken, the local spec is implicitly unsound)
- `error: "Upstream REQ-NNN failed (error_type: <ancestor.error_type>); downstream blocked. Original error: <original error message>"`

Create the follow-up REQ per the Spec row below. It inherits `addendum_to: <this failed REQ>`; the cascade is now visible in the addendum chain and the follow-up's error description names the upstream root cause. The follow-up should also carry the original dependency list so it re-blocks on the same upstream until the upstream's own follow-up lands — and it always emits the canonical `depends_on:` key, even if the failed REQ used the legacy `dependencies:` alias. Don't propagate the alias on follow-ups.

If no upstream REQ is `failed`, fall through to the symptom-based classification table:

| Type | Symptoms | Recovery |
|------|----------|----------|
| **Intent** | Requirements are ambiguous or contradictory; builder couldn't determine what to build | Create a follow-up REQ with `status: pending-answers` containing the specific ambiguities as Open Questions. Archive original as `failed` with `error_type: intent`. |
| **Spec** | Requirements are clear but the technical approach was wrong (wrong files, wrong pattern, wrong architecture) | Create a follow-up REQ with a `## Prior Attempt` section summarizing what was tried and why it failed. Set `status: pending`. Archive original with `error_type: spec`. |
| **Code** | Approach was right but implementation has bugs (tests fail, runtime errors, logic errors) | Create a follow-up REQ targeting the specific code issue. Set `status: pending`. Archive original with `error_type: code`. |
| **Environment** | External dependency unavailable, permissions issue, tooling broken | **First apply the blocked-flip test** (see `actions/work.md`'s mid-run blocked-flip procedure): if *no substantive implementation edits landed this attempt* AND the missing thing is a precondition expected to become available on its own (a service comes up, a person answers, credentials get provisioned), do **not** fail — flip the REQ to `status: blocked` with `blocked_by` + `blocked_at` (non-terminal, stays in the queue) instead. Reserve `error_type: environment` for post-work breakage, a broken/permission-denied environment the user must repair, or a precondition that will not self-resolve. Then: no follow-up REQ — archive with `error_type: environment` and a clear description of what's needed. |

**Anti-rationalization addition.** When checking the symptom table:

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "This REQ failed on a code bug" | Check whether any `addendum_to` or `depends_on` ancestor (or `dependencies:` alias) is also `failed` first | Downstream failures often inherit upstream rot; misclassifying as `code` chases phantom bugs in the wrong domain |

**Procedure:**
1. Run the upstream-failure short-circuit. If it fires, jump to step 3.
2. Otherwise classify using the symptom table above.
3. Update frontmatter: `status: failed`, `error: "description"`, `error_type: [intent|spec|code|environment]`
4. For Intent/Spec/Code failures: create the appropriate follow-up REQ (details above). Set `addendum_to` to the failed REQ's ID so context chains. Preserve the original dependency list on the follow-up when the failure was upstream-driven — always emit it under the canonical `depends_on:` key, even if the failed REQ used the legacy `dependencies:` alias.
5. Move to `archive/` (failed REQs always go to archive root, not into UR folders).
6. Report to user: `[REQ-NNN] failed ([type]): [description]. Follow-up: [REQ-NNN] / None.` When the short-circuit fired, prefix the report with `(upstream cascade — original failure at REQ-NNN)`.

## Changelog Entry Procedure (Step 9)

This procedure owns judgment and payload authoring only. After choosing the source, version, mirrors, existing format, title, voice, and exact replacements below, pass them to one canonical `release` manifest as directed by `actions/work.md` Step 9. `release` is the sole deterministic version/changelog writer and validates monotonicity, preimages, keys, anchors, mirrors, and installed-metadata exclusions in one transaction. Missing/refused tooling stops the release tail; manual edits and package-manager/helper fallbacks are forbidden.

Every successful REQ (`completed` / `completed-with-issues`) gets an entry in the target repo's root `CHANGELOG.md`, written **before** the commit so it ships inside it. Failed and cancelled REQs get no entry — the changelog records delivered change, not attempts. A changelog entry is a human-facing artifact: load `crew-members/anti-slop.md` before writing it (its JIT_CONTEXT condition already covers this — noted here so it isn't skipped).

**Already-green repair exception:** a successfully reviewed repository-gate repair no-op delivered no project change, so it skips this entire procedure: no changelog entry, version/lock mirror, release transaction, or `release_at`. Its exact evidence, archive, and commit rules live in **Repository Gate Deferral and Resumption** → **Already-green repair no-op completion**.

**Precedence check first.** If the repo already has a `CHANGELOG.md` whose entries follow a different convention (keep-a-changelog categories, generated conventional-commit logs, plain dated lists), **match the existing format** — never impose the house voice on a repo with its own. Everything below applies when there is no changelog yet or the existing one already follows this format.

**Bootstrap.** If no root `CHANGELOG.md` exists, create one:

```markdown
# Changelog

What's new, what's better, what's different. Most recent stuff on top.

---
```

**Entry key.** Always `## X.Y.Z — [Short Descriptive Title] (YYYY-MM-DD)` — every entry carries both a version and a date. The title must say what was delivered so a reader scanning only headings knows what changed ("Board View Filters", not a whimsical codename). It must be unique against every existing entry in the file (grep before writing — duplicates have occurred), and the new `X.Y.Z` must be **strictly greater** than the version in the file's first existing entry (duplicate version numbers have occurred).

**Where `X.Y.Z` comes from.** Resolve the version source once per entry, in this order:

1. **A version in a project-owned release file** — `package.json`'s `"version"`, `Cargo.toml`, `pyproject.toml`, a `VERSION` file, or the like (this list is illustrative; any file the project's own release process maintains a version line in qualifies). Exclude every installed skill, dependency, vendored package, and generated tree from candidate discovery regardless of whether Git tracks it; `.claude/skills/`, `.codex/skills/`, `node_modules/`, and `vendor/` are examples, not the boundary. In particular, the installed do-work suite's `VERSION` and `actions/version.md` are runtime metadata and must never be selected or bumped for a consumer REQ. Prepare the selected project-owned old/new bytes for the release manifest. The repo's version and changelog header stay in lock-step; include any committed lockfile mirror below in the same manifest.
2. **Version only in release tags** (no version line in any file) — take the highest release tag as the current version and bump from it, but write the result **only into the changelog header**. do-work never creates a git tag: a tag is a release announcement, and only a human decides when one happens.
3. **No version anywhere** — the changelog is the source of truth. Take the highest `X.Y.Z` across the file's existing entries and bump from it. If there are no entries yet (bootstrap), seed the first entry at `0.1.0`. Nothing outside `CHANGELOG.md` is touched — an unversioned repo stays unversioned, and the header number is a changelog fact, not a claim that a release was cut.

Two guards on source resolution. If **two or more version files disagree** with each other, do not guess which one the release process uses: leave every file untouched, fall back to the changelog counter (source 3), and say so in the Step 9 report. If the resolved source is **behind** the newest changelog entry (someone released or edited out of band), bump from whichever is higher — never emit a version below one already in the file.

**Lockfile mirror (source 1 only).** Some lockfiles record the version of the package you just bumped, not only its dependencies' — so bumping the version file alone leaves a stale copy behind, and the next ordinary build or install rewrites it: the tree reads dirty for a change nobody made, and the fix accretes as ad-hoc "sync the lockfile version" commits. Under npm it is only cosmetic (`npm ci` tolerates a version-only mismatch and leaves it stale, which is exactly why the drift survives unnoticed for many versions), but `cargo check --locked` exits 101 on it and `uv lock --check` exits 1, so elsewhere it hard-fails CI.

The trigger is the condition, not the ecosystem — *the repo commits a lockfile that records this package's own version*. Known instances, not the boundary:

| Lockfile | Where this package's own version is mirrored |
| -------- | -------------------------------------------- |
| `package-lock.json` | the top-level `"version"`, **plus** `packages[""].version` when the file has a `packages` map (`lockfileVersion` 2 or 3; a `lockfileVersion` 1 file has one site, not two). Bumping a **workspace member** is different: its version lives at `packages["<member-path>"].version` in the root lockfile, and the two root-package sites hold the *root's* version and correctly stay put — so read them as "nothing to sync" and you leave the real mirror stale |
| `Cargo.lock` | the `[[package]]` entry whose `name` is this crate — present even in a crate with no dependencies |
| `uv.lock` | the `[[package]]` entry whose `name` is this project (`source = { editable = "." }` or `virtual`) |

Dependency-only lockfiles never trip this: `pnpm-lock.yaml`, `yarn.lock`, `poetry.lock`, and `go.sum` record no version for the root package (yarn Berry writes the sentinel `0.0.0-use.local` on purpose), so a repo with only those has nothing to sync. The split follows the **lock tool, not the manifest** — `pyproject.toml` mirrors under uv and not under poetry — so open the lockfile and look rather than inferring from the manifest's name.

**Prepare the exact mirrored old/new bytes and include them in the release manifest.** It is one or two lines of judged payload content, needs no toolchain and no network, and the command publishes it with the other release targets. Do **not** shell out to the package manager. `npm install --package-lock-only` executes target hooks and may restructure or re-resolve; `cargo generate-lockfile` can drag unrelated dependencies and checksums forward.

Three constraints govern the caller-authored mirror payload. Change **only this package's own entry** — a dependency may legitimately carry the same version string, so a whole-file search-and-replace corrupts resolutions. Include only sites that already exist, and never declare a new lockfile: in a workspace or monorepo the one committed lockfile sits at the root, not beside the member manifest, and the entry inside it is the member's own path per the table above. The payload is a version-line fix, not a dependency resync; unrelated lockfile drift is its own REQ. If the mirrored value is already correct, omit that target.

**Then read the lockfile's diff before you stage it** — it must contain only the mirrored version line(s), and nothing else may ride along in this REQ's commit.

**Bump size.** Read the change the REQ actually delivered, not its wording:

| Bump      | When                                                                                                                                       |
| --------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| **major** | An existing consumer breaks: a public API or CLI flag removed or renamed, an on-disk or wire format changed, a documented default reversed. |
| **minor** | A user-invocable capability exists that didn't before, and nothing existing breaks.                                                          |
| **patch** | Everything else — bug fixes, performance, refactors, tests, docs, internal-only changes.                                                    |

Tie-breakers, in order: a breaking change outranks an additive one in the same REQ (bump major, not minor); when genuinely torn between two levels, pick the **smaller** one. **Below `1.0.0`, a breaking change bumps the minor, not the major** — `0.x` is unstable by semver's own definition, so the first breaking change in a seeded repo must not silently promote it to a `1.0.0` release. A `completed-with-issues` REQ is bumped on what it delivered, exactly like `completed`.

The `CHANGELOG.md` change — and the version file plus any lockfile mirroring it, when source 1 applied — are part of the REQ's lifecycle files. Stage them in the commit below.

**Voice contract (house style).** 1–2 casual sentences leading with *why it matters* — the situation that prompted the change and what's better now — then bullets for the specifics. Lead with value, not implementation; file paths and flags belong in the bullets, not the lead. Keep it brief. Newest on top, one entry per REQ (this matches one-commit-per-request).

```markdown
## 0.4.0 — Clear Questions for Interactive Prompts (2026-07-07)

Agents kept asking questions only they could parse — codenames coined
mid-analysis, options with no stated consequence. Question wording is
now a contract, not a hope.

- New `crew-members/clear-questions.md`, loaded before any interactive ask
- Six principles: one decision per question, decode your own shorthand, say the consequence…
```

## Commit & Metadata-Commit Procedure (Step 9)

```bash
# Stage implementation files + archived REQ + the changelog entry
git add src/stores/theme-store.ts src/components/settings/SettingsPanel.tsx \
  do-work/archive/UR-002/REQ-003-dark-mode.md CHANGELOG.md

# Stage the bumped version file — only when the changelog resolved to a repo
# version file (source 1). Tag-versioned and unversioned repos have none.
git add package.json

# Plus the lockfile that mirrors it, when the repo keeps one — see the Changelog
# Entry Procedure's "Lockfile mirror" note for the condition and for which field(s)
# to edit by hand. package-lock.json here; Cargo.lock / uv.lock in those
# ecosystems. Omit this line entirely when the repo has no such lockfile:
# `git add` on a path that does not exist exits 128 and aborts the commit step.
git add package-lock.json

# Set this guard true only when the successful canonical `complete` result reports
# do-work/calibration-log.tsv among its changes or affected target paths. Otherwise
# leave it false; never infer calibration staging from the filesystem.
if [ "${calibration_changed:-false}" = true ]; then
  git add do-work/calibration-log.tsv
fi

# Stage follow-up REQs created in Step 8 (if any), AND any existing sweep REQs
# the review appended instances to — an append modifies a queue file rather than
# creating one, so "created follow-ups" alone leaves the new ## Instances lines
# unstaged.
git add do-work/queue/REQ-025-confirm-sidebar-palette.md
git add do-work/queue/REQ-021-existing-sweep.md

# Stage the stakeholder REQ Step 8 substep 3 routed questions into (minted or
# folded), plus the fresh report bundle that regeneration published — the
# blocked_by: path must resolve on every checkout, not only this one. Omit both
# lines when this REQ routed no stakeholder questions; omit the bundle line
# where the project ignores ai-reports/ (`git add` on an ignored path fails and
# aborts the commit step).
git add do-work/queue/REQ-030-stakeholder-questions-priya-design.md
git add ai-reports/2026-08-21_1140_REQ-030-questions-priya-design/

# Stage the prose backlog when this REQ touched it — a Step 8 review append
# (actions/review-work.md Step 10) or a drain's ticks. The orchestrator writes
# this file in the main tree, so nothing else stages it. Omit this line when the
# REQ never touched it: `git add` on a path that does not exist exits 128 and
# aborts the commit step.
git add do-work/prose-backlog.md

# Stage UR-folder move (if this was the last REQ and the UR moved to archive/)
# Both the old path (deletion) and new path (addition) must be staged.
# Exception: if the UR was never committed (capture's commit step was skipped,
# or the repo wasn't git at capture time), the old path matches nothing and
# `git add` exits 128 — stage only the new archive path in that case.
git add do-work/user-requests/UR-002/ do-work/archive/UR-002/

git commit -m "$(cat <<'EOF'
[REQ-003] Dark Mode (Route C)

Implements: do-work/archive/UR-002/REQ-003-dark-mode.md

- Created src/stores/theme-store.ts
- Modified src/components/settings/SettingsPanel.tsx

EOF
)"
```

**Format:** `[{id}] {title} (Route {route})` + `Implements:` line + summary bullets. Add a co-author trailer if your platform convention calls for one (e.g., `Co-Authored-By: Agent <agent@example.com>`), otherwise omit.

One commit per request. Stage all files created, modified, moved, or deleted during this request's lifecycle: implementation files (listed in the Implementation Summary), the archived REQ file, the `CHANGELOG.md` entry and the version file it bumped plus any lockfile mirroring that version, if any (successful REQs only — see the Changelog Entry Procedure above), any follow-up REQs created in Step 8 (`pending-answers` files in `do-work/queue/`), any stakeholder REQ Step 8 substep 3 minted or folded into (a `stakeholder:` file in `do-work/queue/`) together with the fresh `ai-reports/` bundle that regeneration published — the path recorded in `blocked_by:` must resolve on every checkout, not only this one; skip the bundle where the project ignores `ai-reports/` — `do-work/prose-backlog.md` when this REQ touched it — a Step 8 review append lands there as well as a drain's ticks, and any UR-folder moves to `archive/`. Stage `do-work/calibration-log.tsv` only when the successful canonical `complete` result reports it among its changes or affected target paths; otherwise do not stage it. If Step 8 substep 7 wrote prime-file lessons links, the modified prime files must also be staged — they are part of the REQ's lifecycle changes even though they aren't listed in the Implementation Summary's `Files changed`. Do not use `git add -A` or `git add .`, and never bypass a commit hook (see `actions/commit.md` § Rules for the full guard). Failed requests get committed too.

**In worktree dispatch mode** the builder's implementation is already committed and merged (Step 6's `--no-ff` merge), so do **not** stage implementation files here — stage only the archived REQ, the `CHANGELOG.md` entry, any bumped version file and its lockfile mirror, follow-up REQs, any stakeholder REQ Step 8 substep 3 minted or folded into plus the fresh `ai-reports/` bundle that regeneration published (skip the bundle where the project ignores `ai-reports/`), UR-folder moves, `do-work/calibration-log.tsv` only when the successful canonical `complete` result reports it among its changes or affected target paths, `do-work/prose-backlog.md` when this REQ touched it — a Step 8 review append lands there as well as a drain's ticks, and prime-file lessons links. **Both `do-work/` files are the orchestrator's writes in the main tree, not the builder's**, so they are never in the merge and this list is the only thing that stages them — omitted, an appended or ticked backlog item and an appended calibration line stay a dirty tree. The `commit:` field gets the `--no-ff` merge commit's hash (`<merge_hash>`, captured in Step 6 — the **latest** merge if remediation re-merged; **Worktree Dispatch Mode (Step 1)** above), not this changelog commit's hash.

**Validation check (successful REQs only):** Before committing, compare the `## Implementation Summary` file list against the staged files (excluding `do-work/` paths). If the Implementation Summary lists files that aren't staged, or if the only staged files are `do-work/` metadata, `CHANGELOG.md`, and/or the version file it bumped together with any lockfile mirroring that version (the changelog entry and the version bump describe the implementation, they aren't the implementation — and a lockfile carrying only the mirrored version is part of the bump, not a deliverable), flag the mismatch — the commit may not contain the actual implementation. Fix the staging or update the Implementation Summary before proceeding. Design-artifact files placed outside `do-work/` satisfy this check — they are project deliverables. **Skip this check for failed REQs** — they may have no Implementation Summary or no project files staged, and that's expected. **In worktree dispatch mode** the implementation files live in the merge commit, not this commit's stage, so validate the `## Implementation Summary` file list against `git diff --name-only <pre>..<merge_hash>` (the merge range, excluding `do-work/` paths) instead of the staged set — a stage of only the changelog/version/`do-work/` metadata is correct here, not a mismatch.

**Already-green repair validation exception:** only the exact no-op evidence and summary shape above may validate with no project files. Its stage must contain only canonical lifecycle/archive/calibration paths and any exact UR move reported by `complete`; the presence of a project, changelog, version, lockfile, or unrelated path is a mismatch rather than something to absorb.

**Write commit hash back to the archived REQ.** After the commit succeeds, resolve the implementation hash — **serially** it is the commit you just made, so read it with `git rev-parse --short HEAD`; **in worktree dispatch mode do NOT rev-parse `HEAD` here**, because HEAD is the changelog commit and the implementation lives in Step 6's `--no-ff` merge — use the `<merge_hash>` literal held since Step 6 (the latest merge, if remediation re-merged). Pass the hash and exact archived path to `complete --record-commit-hash --implementation-hash <hash> --request-path <path> --commit`. This metadata-only lifecycle mode validates the archived terminal-success preimage, edits only `commit:`, creates a **separate metadata commit**, and verifies the committed bytes. Never amend/reset, and never substitute a free-form edit or legacy helper. `OK`, `NOOP`, refusal, rollback, and committed-risk are returned through the same typed result contract as every other lifecycle transition, including exact verification or revert argv.

**If the canonical request-state command is missing, stop.** A consumer on an older tarball must upgrade before performing lifecycle metadata writes; there is no hand-edit or helper fallback.

This ensures the `commit:` field in the archived REQ contains the real implementation commit hash — the commit just made serially, the `--no-ff` merge commit in worktree dispatch mode — which review-work and completed-work presentation actions depend on for traceability. The metadata commit is a lightweight bookkeeping entry — it does not contain implementation changes.

## Session Checkpoint Template (Step 10)

Its `session_ended` is the current UTC instant (Timestamp rule, above).

```markdown
---
session_ended: <timestamp>
last_completed: REQ-NNN
queue_state: [N pending, N pending-answers, N blocked, N blocked-archive-collision, N blocked-dependency-cycle, N in-progress]
reqs_processed_this_session: N
session_depth: light | moderate | heavy
---

# Session Checkpoint

## Completed This Session
- REQ-NNN: [title] (Route [X], [score]%)
- REQ-NNN: [title] (Route [X], [score]%)

## In Progress (interrupted)
- REQ-NNN: [title] — claimed <claimed_at> — writer: <hostname>:<absolute-checkout-path> — stopped at [phase: triage/planning/exploring/implementing/testing/reviewing]
  Last known state: [1-2 sentences]
  Key files being modified: [list]
  Known issues: [any blockers or concerns]
- REQ-NNN: [title] — claimed <claimed_at> — writer: <hostname>:<absolute-checkout-path> — stopped at [phase: ...]
  [one entry per REQ still in do-work/working/ under this checkout's writer label — this section is a
   list, never a single id; every entry this checkout did not write is copied through verbatim, a
   foreign writer label and no label at all alike]

## Still Queued
- REQ-NNN: [title] (pending)
- REQ-NNN: [title] (pending-answers — [N] questions)

## Session Notes
[Environment issues, user preferences expressed, patterns discovered, decisions made outside REQ files]

## Context Summary (heavy sessions only)
[Recap of key decisions (D-XX references), architectural patterns encountered, and prime files
that the next session should re-read before starting. Include this section when 6+ REQs were
processed — at that volume, carried-over assumptions are unreliable.]
```

**`## In Progress (interrupted)` is not written from scratch here.** Step 2 has been appending an entry to it at every claim (**In-Progress Record (Step 2)**, above) and Step 8 removing each one on departure, so at session end the section already holds exactly the REQs still in `do-work/working/`. Step 10 enriches those entries with the phase/state detail above and keeps the list — it never collapses it to a single id, because recovery classifies each `working/` file by name and `writer:` label, and fan-out can leave several. **It enriches only this checkout's own entries, and carries every other one through verbatim** — same id, same `claimed_at`, same `writer:` label, no phase detail invented for work this session never did. Step 10 is the one place that rewrites the whole file, which makes it the one place that can destroy a label: a wholesale rewrite that knows only about its own claims would silently drop another checkout's live claim, and the next run would then have nothing to classify that REQ by.

## Progress Reporting Example

```
Processing REQ-003-dark-mode.md...
  Triage: Complex (Route C)
  Open Questions: 2 found → builder decided (follow-ups queued)
  Planning...     [done]
  Scope...        [done] 4 files declared
  Exploring...    [done]
  Implementing... [done]
  Summary...      [done] 3 files changed
  Qualifying...   [done] ✓ files verified, requirements traced
  Testing...      [done] ✓ 12 tests passing
  Reviewing...    [done] 92% — 0 follow-ups
  Archiving...    [done]
  Committing...   [done] → abc1234

Processing REQ-004-fix-typo.md...
  Triage: Simple (Route A)
  Implementing... [done]
  Summary...      [done] 1 file changed
  Qualifying...   [done] ✓ verified
  Testing...      [done] ✓ 3 tests passing
  Reviewing...    [done] 88% — 0 follow-ups
  Archiving...    [done]
  Committing...   [done] → def5678

All 2 requests completed:
  - REQ-003 (Route C) → abc1234 [review: 92%]
  - REQ-004 (Route A) → def5678 [review: 88%]
```

## Decision Brief (hand-back format)

The canonical shape for handing work back to the user. Used by the end-of-run completion hand-back (work.md Step 10 / Progress Reporting), `actions/clarify.md`'s question presentation, and `actions/review-work.md`'s Step 9 report. Lead with what was built and what needs the user; **never lead with a self-grade**. Applies `crew-members/anti-slop.md` § 8 (lead with the decision) and the decide-vs-escalate gate in `crew-members/coding-guardrails.md` § Think Before Coding. Render only the sections that have content.

```
WHAT'S BEING BUILT            (feature + subsystem altitude — the value)
  • Now you can X — lives in <subsystem>
  • [MAP CHANGED] <new data flow / renamed concept / new contract>   ← only if true

DECISIONS FOR YOU             (escalated by exception — each carries value + risk)
  ▸ <decision, one line>
      Value:  what you gain / why it matters
      Risk:   cost, reversibility, what breaks if the choice is wrong
      → recommend <X>; default if you say nothing: <X>

WAITING ON OTHERS             (FYI — work is done; assumptions await outside confirmation)
  • <name>: N questions (K irreversible) — share the report: <bundle path>

HANDLED  (FYI — spot-check, don't ratify)
  • decided <Y> because <Z>     ← reversible calls made without asking
```

- **WHAT'S BEING BUILT** renders each REQ's `## Orientation` block (work.md Step 7.5) at feature/subsystem altitude — not a file list. Anchor to the touched `prime_files`; flag `[MAP CHANGED]` only when the change alters the system's shape.
- **DECISIONS FOR YOU** renders the **ESCALATE**-tier decisions — the `- [~]` / `D-NN` entries that became `pending-answers` follow-ups — each with the Value/Risk carried from the decision record; entries routed to a stakeholder (`Answerer:` clause) render under WAITING ON OTHERS instead. Source Value/Risk from the touched prime's `## Stakes` when present, else builder-derive.
- **WAITING ON OTHERS** renders the stakeholder-routed questions, one line per stakeholder REQ touched this run, with the latest report bundle path — the session user is the courier, and the share-the-report line is what closes the loop. Nothing here asks the user to decide anything; omit if empty.
- **HANDLED** lists the **DECIDE & STATE** decisions (reversible `D-NN` entries) so the user can spot-check without being asked to ratify. Read each REQ's `## Decisions` per **Reading a Builder-Authored Section (any step)**, above — under fan-out the builder wrote it into its hand-back, and a brief that reads the REQ file alone renders every builder's decisions as an empty list. Omit the block when the sections were read and held no DECIDE & STATE entry; when a REQ's section was in neither place and its hand-back could not be read, render the block with `• REQ-NNN: decisions not recovered — hand-back unread` instead of omitting it, because nothing recorded and nothing readable are different facts.
- **Scale context to reach.** A leaf REQ collapses to a single WHAT'S BEING BUILT line with no DECISIONS and a short HANDLED list; a map-changing REQ earns a short paragraph and a why-it-matters. Review scores never lead — they live under the decision (review-work Step 9) or in the per-REQ progress lines above.
