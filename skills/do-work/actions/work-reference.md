# Work Action — Reference

> Companion file to `work.md`. Holds the heavy templates, tables, and sub-procedures the orchestrator references by name — extracted to keep `actions/work.md` focused on the ten-step skeleton. Each section below is pointed to from the matching step in `actions/work.md`. Loading this file is only necessary when you reach the step that references it — and read only the named section. If this file is already in context from an earlier step this session, reuse it; don't re-read it at every reference site.

---

## Architecture

The per-REQ pipeline is `actions/work.md`'s numbered steps; this section is the ownership map beside them, not a second drawing of the sequence.

- **The loop.** `recover` → select and claim → triage → open questions → mechanical evidence gate → route → build → Implementation Summary → qualification → focused/gate evidence → review (one remediation pass) → prepare finalization → finalize → checkpoint → cleanup. The orchestrator stays in the loop and stays light.
- **The three routes converge.** Route A goes straight to the build; Route B explores and declares scope first; Route C plans before both. All three meet again at the pre-build evidence record and run every evidence step from there identically.
- **Who owns what.** The orchestrator owns judgment and every authored `##` section. Canonical commands own the deterministic mutations: `advance` (selection, claim, mechanical evidence gates, checkpoint, finalization hand-off), `defer-gate` (repository-gate deferral), and the finalization engine (archive, release, commit, provenance, verification, cleanup).

## Execution Model — Claim Anywhere, One Releaser

**Any checkout may capture and claim.** A queue can be shared by however many checkouts a user points at it — a spawned worktree, a second local workspace, a clone, a cloud sandbox — and each of them may write REQ files, claim them, and build. Claims and captures reach the other checkouts the way everything else does: by ordinary git sync, with no lock, no lease, and nothing to acquire. The canonical claim transaction commits its queue-to-working move plus checkpoint entry immediately, in every mode, so a successful claim has one durable footprint before implementation begins. This narrows duplicate work to checkouts that have not synced; it does not add coordination or arbitration, and any collision still surfaces when branches meet (**Cross-checkout collisions are merge artifacts**, below).

What the model constrains is not who *works* but who *releases*. The pipeline supports **one releaser per queue** — the single checkout that runs the release tail: merge integration, the version bump, the `CHANGELOG.md` entry, the archive moves, and UR closure. **Two releasers against one queue stays outside the contract**, and so do two sessions in one working tree — the pipeline does not detect, coordinate, or recover either, and spends no durable state on making one safe. Behavior in both cases is unspecified, and the repair path is after the fact and human: `actions/forensics.md` to see the damage, `actions/cleanup.md` to fix it.

**Builders are not owners.** Any number may build at once (**Worktree Dispatch Mode** → Fan-Out Dispatch, below), because a builder writes only its own tree and owns no queue state. What changed is only that the tree may belong to a different checkout than the releaser's.

**Cross-checkout collisions are merge artifacts, not scheduling failures.** Two unsynced checkouts claiming the same REQ, or capturing REQs that land on the same id, produce ordinary git conflicts that get fixed when the branches meet — that is the entire coordination mechanism, and prose elsewhere cites this sentence rather than restating it. Both cheap detectors already ship: `queue-kanban verify`'s `duplicate-req-id` probe for colliding captures, and the `writer:` label on checkpoint entries for claims (**In-Progress Record (Step 1)**, below), which makes even a byte-identical double claim surface when the committed footprints meet. The label also lets crash recovery detect and report another checkout's synced live claim as `claim held by <writer>, not touched`. Detecting is all it does: it still neither coordinates nor recovers one, and it never arbitrates.

**One broken pipe does not stop the rest of the factory.** Whenever a failure is local to one REQ, set that REQ aside with a typed finding and keep draining the runnable queue; a failed build, an unrelated gate failure, and a recoverable finalization tail are illustrative, not a closed list. Only shared-target dirt stops the loop, because every next claim would write through the same target. That stop's typed finding names `do-work run-with-recovery` as the resolving verb when the user can assert sole queue authority.

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

**Stamps are append-only.** Once an `*_at` field exists on a REQ, no lifecycle transition deletes it or writes a different instant over it — a phase re-entered after a recovery, a hold, or a gate deferral writes its stamp only when the field is **absent**, so `claimed_at` keeps naming the claim that started the work and the wall time an archived REQ reports is the real one. **The suffix is the condition**, never a list of today's fields, so a stamp this schema gains tomorrow is governed without editing this paragraph. Four fields are documented exceptions because they carry *current state* rather than a phase observation, and each is marked where it is defined below: `status_changed_at` (the most recent status change, by definition), `completed_at` on the `failed` → `cancelled` path, `blocked_at` (removed together with `blocked_by` when an unblock clears the condition), and `heavy_verified_at` / `heavy_verified_revision` (withdrawn together with the `commit` they verified when a REQ is re-claimed). Repairing a stamp that is *detectably wrong* is a different operation and stays allowed: `doctor --repair-timestamps` derives the instant from git history and is never part of a transition.

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
priority: next                # OPTIONAL authored queue rank: now | next | later. Capture writes it only when the user's words explicitly rank timing or order; absence and an unrecognized value both resolve to `next`, with the latter producing a Schema Read Contract warning. It orders ordinary dependency-ready work before fan-out bounding and sorts Pending Ready and Pending Waiting independently on the board. It never outranks repository-gate repairs, ready deferred parents, or `depends_on`; equal priorities keep existing queue order. The board shows compact `now`/`later` badges and no badge for the default `next`.
effort_estimate: effort-substantive   # OPTIONAL triage bit: effort-mechanical | effort-substantive — separates small mechanical fixes from real work so the user can tell at a glance which queued REQs are cheap to approve or batch. **This field is SIZE, judged as size by whoever writes it — never derived from an impact verdict; that judgment is `impact:` above.** Closed two-value enum, deliberately — a triage bit, not an estimation system; do not grow it toward t-shirt sizes (the estimation system lives in the separate `estimate:` block below, and the only bridge between the two is the mechanical-effort short-circuit: `effort-mechanical` ⇒ the floor estimate, no signal extraction — `actions/estimate-reference.md`). Absent or unrecognized reads as `effort-substantive` (Schema Read Contract row below), and the read-only legacy aliases in that row carry every REQ written before the rename unchanged. **Expected on every new REQ:** capture, review follow-up creation, and Discovered Tasks creation judge it by the same three-way standard — judge it, or put the judgment to the user, or leave it absent because neither was possible, never a copied default (`actions/capture.md` Step 1's effort assessment; `actions/review-work.md` Step 10; **Discovered Tasks Classification (Step 8)** below). An unjudged REQ reads as `effort-substantive` and drops out of the selector below for no reason anyone chose. Read as a selection filter by `tools/select-simple-reqs.sh` (backing `actions/run-simple-reqs.md`), which selects only the REQs normalizing to `effort-mechanical` so a cheaper-model session can run the queue's small work; the read-only `trivial` alias above is load-bearing for that selector, because REQs written before the rename spell the value that way and a literal match on the canonical token alone silently drops them. Display only: parsed by `../../do-work-board/tools/queue-kanban/model.go` into a card chip (rendered only when `effort-mechanical` — `effort-substantive` is the default and would be noise) and a drawer row, with no column logic and no scheduling — keep that parser in lock-step with this line, both changing in the same commit.

# OPTIONAL informational forecast — never read by scheduling, gating, or pipeline
# logic, and FROZEN once execution begins. p50_active_minutes is a multiple of 5 and
# never below 5: roughly a 50% chance of finishing within that many ACTIVE agent
# minutes (user wait, paused sessions, and queue wait are excluded by definition).
# Written by the work action's ensure-estimate step or by verify-requests after a
# material repair, produced deterministically by tools/estimate-p50.sh. Extraction
# guide, confidence rubric, and presentation formats: actions/estimate-reference.md.
estimate:
  p50_active_minutes: 75
  confidence: medium            # low | medium | high
  calculated_at: 2026-08-16T12:00:00Z   # current UTC instant (Timestamp rule above)
  basis:
    - Route C
    - 12-file write set
    - browser evidence

# Set by work action when claimed — written once (Stamps are append-only, above):
# a REQ recovered and claimed again keeps this first instant, and crash recovery
# clears `route` without touching it.
claimed_at: 2025-01-26T10:30:00Z
route: A | B | C

# OPTIONAL observations written only after the named work-pipeline event succeeds.
# Routes that skip an event omit its field; writers never fabricate a skipped phase.
# No transition removes one afterwards (Stamps are append-only, above).
planning_at: 2025-01-26T10:32:00Z         # Route C plan saved and validated
dispatch_at: 2025-01-26T10:33:00Z         # implementation builder accepted the dispatch
builder_handback_at: 2025-01-26T10:40:00Z # builder returned its completed hand-back
integration_at: 2025-01-26T10:41:00Z      # hand-back integrated
review_at: 2025-01-26T10:43:00Z           # first review result recorded
remediation_at: 2025-01-26T10:46:00Z      # one remediation hand-back integrated
re_review_at: 2025-01-26T10:48:00Z        # post-remediation review result recorded
release_at: 2025-01-26T10:50:00Z          # canonical release transaction succeeded

# Set by capture (external-condition task) or by the work pipeline's mid-run blocked flip (Step 8's blocked-flip procedure). Holding state — the REQ stays in do-work/queue/ and the default scan walks past it, exactly like pending-answers.
status: blocked               # waiting on an EXTERNAL condition — not user answers (that's pending-answers), not another REQ (that's depends_on). Cleared to `pending` by a passing blocked_check probe (work Step 1), a `do-work clarify` confirmation, or a manual edit.
blocked_by: 'LM Studio running locally'   # free text naming the condition (raw user text — **Frontmatter Quoting** contract above). Legacy note for the board: an old id-LIST value renders joined for display and is NOT a dependency edge — dependency gating is `depends_on` only.
blocked_at: 2026-07-18T10:00:00Z          # stamped on every flip to blocked — the age anchor the exit summary, board drawer, and forensics read (no enforcement threshold; external conditions legitimately take weeks). One of the four documented exceptions to **Stamps are append-only** (above): an unblock removes it together with blocked_by, because the pair describes a condition that no longer holds.
blocked_check: 'curl -sf http://localhost:1234/v1/models'   # OPTIONAL shell probe (raw user text — **Frontmatter Quoting** contract above). User-authored content, run VERBATIM by work Step 1 (exit 0 ⇒ unblock to pending; any non-zero / timeout / unreadable ⇒ stays blocked). Absent ⇒ manual/clarify unblock only.
stakeholder: 'Priya (design)'   # REQUIRED on stakeholder-questions REQs (meaningless elsewhere; raw user text — **Frontmatter Quoting** contract above): the outside person whose confirm-or-override answers this REQ collects — presence is the marker, value is the fold discriminator (actions/capture-reference.md → Fold-First Rule → Stakeholder-audience questions). Verbatim-read class, like assigned_to: no alias map, no case folding, trim-only; greppable by design (grep -rl '^stakeholder: ' do-work/queue/). Always paired with status: blocked + blocked_by naming the person and the latest report bundle path (or "report pending regeneration" until a bundle lands — actions/work.md Step 8) + blocked_at; never with blocked_check (a person is not probeable), and deliberately never with user_request: — UR membership would hold the first source UR open in every closure reader, and nothing waits on this REQ (question provenance lives in per-entry Source: lines). NOT parsed by the board — display rides the existing blocked_by badge, zero parser change. Nothing gates on this REQ: its source REQs completed on the builder's assumptions, and it exists only to route answers back (actions/stakeholder-answers.md); clarify routes it, never yes/no-confirms it (actions/clarify.md Step 5.5).

# Set on ANY status flip that has no dedicated *_at stamp of its own — that condition
# is the rule and the writers are illustrative (answered → pending, unblock → pending,
# a manual or stuck reset). An unblock REMOVES blocked_at, so this is the only trace of
# when that flip happened. Flips with a dedicated stamp (claimed_at, blocked_at,
# completed_at) do NOT write it. Display-only: the board's pending-tier state timer
# prefers it over created_at/file-mtime. Timestamp rule applies (current UTC instant).
# This is the one stamp every transition overwrites on purpose — it means "most
# recent status change" by definition (**Stamps are append-only**, above).
status_changed_at: 2026-07-22T20:38:00Z

# Set by work action when finished. STAMPING RULE: every flip to a terminal status
# (completed / completed-with-issues / failed / cancelled) MUST stamp completed_at with
# a UTC ISO instant, plus commit with the implementation hash in a git repo. These two
# are the ONLY sources the board resolves a terminal REQ's completion instant from (no
# file-mtime fallback); missing both, an unparseable stamp, or a hash git cannot resolve
# is a completion anomaly in all three do-work-board board modes.
completed_at: 2025-01-26T10:45:00Z   # required on every terminal flip — UTC ISO instant
status: completed | completed-with-issues | failed
commit: abc1234               # required in a git repo — implementation commit hash (also recorded at the heavy hold so dependents can build against the landed source)
error: "Description"          # Set when a REQ failed; RETAINED verbatim if that failed REQ is later cancelled via do-work abandon — the surviving failure signal on a status: cancelled REQ, NOT drift to strip
error_type: intent|spec|code|environment   # Set with `error` on failure; likewise retained on a failed→cancelled flip

# Written by work.md Step 7.7's drain onto the claimed record before finalization.
# They prove which execution revision the heavy lanes actually checked. A
# documented exception to **Stamps are append-only** (above): a re-claim deletes
# both together with the `commit` they verified, so dependents never build
# against withdrawn work.
heavy_verified_at: 2026-09-03T12:00:00Z
heavy_verified_revision: def5678

# Set by abandon action (do-work abandon — user-directed won't-do decision). Two entry
# paths: a not-yet-finished REQ, and an already-archived `failed` REQ resolved after the
# fact — the latter keeps its error/error_type, so error-on-cancelled is valid data.
status: cancelled             # terminal, NOT successful; the reason lives in the REQ body's `## Cancelled` section
completed_at: 2025-01-26T10:45:00Z  # stamped (or, on the failed→cancelled path, re-stamped to the cancellation instant — a documented exception to **Stamps are append-only**, above; the failure instant survives in the `## Cancelled` section's `Previously:` line) — the terminal timestamp the board's recently-done window reads

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
| `status` (Step 1 scan + categorization, Step 9 finalization, abandon action) | `pending`, `claimed`, `completed`, `completed-with-issues`, `failed`, `cancelled`, `pending-answers`, `blocked`, `blocked-archive-collision`, `blocked-dependency-cycle` | `complete`/`done`/`finished`/`closed` → `completed`; `canceled`/`abandoned`/`wont-do`/`wontfix` → `cancelled` | skip REQ at Step 1 with the warning text — never claim or archive an unrecognized status silently |
| `route` (Step 3 dispatch, Step 5.5 scope declaration, Step 7 scope-drift comparison) | `A`, `B`, `C` | lowercase `a`/`b`/`c` → uppercase | treat as needing re-triage in Step 3 |
| `caveman` (Step 6 crew load) | `false`, `true`, `lite`, `full`, `ultra` | truthy strings (`yes`/`on`) → `true`; `light` → `lite` | `false` |
| `maintenance` (Step 6 crew load) | `true`, `false` (YAML boolean) | truthy strings (`yes`/`on`/`t`) → `true`; `no`/`off`/`f` → `false` | `false` (Step 6 maintenance crew not loaded) |
| `tdd` (Step 6 testing-crew load, Qualification and Testing Judgment's TDD-evidence gate; emission validated in `actions/capture.md`) | `true`, `false` (YAML boolean) | `test_first`/`yes`/`on`/`t` → `true`; `no`/`off`/`f` → `false` | `false` (Step 6 testing crew not loaded; TDD-evidence gate not enforced) |
| `error_type` (Step 8 failure classification, Step 8 upstream-failure short-circuit, forensics) | `intent`, `spec`, `code`, `environment` | (no common typo aliases identified) | `code` |
| `kb_status` (kb-lessons handoff — work.md's Lessons-Capture Phase / review-work.md's Self-Validation & Lessons Learned step; roadmap lessons rollup) | `promoted`, `pending`, `declined`, `skipped` | `skip` → `skipped`; `rejected` → `declined` | `pending` |
| `impact` (capture emission — capture.md Step 1; automatic follow-up creation — review-work.md Step 10, work.md Step 8's Discovered Tasks flow; selection filters — work.md Step 1's `--skip-impact-negligible` and `tools/select-simple-reqs.sh`'s `impact-critical` veto; board display — `../../do-work-board/tools/queue-kanban` parser) | `impact-critical`, `impact-user-visible`, `impact-rule-change`, `impact-negligible` | (no aliases — new prefix-unique vocabulary) | `impact-user-visible` |
| `priority` (capture emission/addenda — capture.md Steps 1–2; ordinary ready ordering before fan-out — work.md Step 1 selector; Pending Ready/Waiting order and badges — `../../do-work-board/tools/queue-kanban` parser) | `now`, `next`, `later` | (no aliases) | `next` |
| `effort_estimate` (Mechanical Evidence-Gate Loop's effort short-circuit; selection filter — `tools/select-simple-reqs.sh`, backing `actions/run-simple-reqs.md`; board display — `../../do-work-board/tools/queue-kanban` parser; judged by capture, review follow-up creation, and Discovered Tasks creation on every new REQ) | `effort-mechanical`, `effort-substantive` | `trivial` → `effort-mechanical`; `normal` → `effort-substantive` (read-only legacy aliases, so every REQ written before the rename stays valid unchanged; never propagated on write — `actions/capture-reference.md` § Schema Aliases) | `effort-substantive` |
| `testing_status` (board Testing view — `../../do-work-board/tools/queue-kanban` parser + `/api/testing/status` writes; no work-pipeline read sites) | `in-testing`, `tested`, `returned` | `in_testing`/`in testing`/`testing`/`selected-for-testing`/`selected for testing` → `in-testing`; `returned-with-feedback`/`returned_with_feedback`/`returned with feedback` → `returned` | treat as not-tested (Ready to test) with an invalid flag + data warning |
| `builder_decided` (clarify's confirm routing — `actions/clarify.md` Step 4/Step 5; reversal detection — `actions/clarify.md` Step 4's `overturned_decision_sources` and `actions/verify-requests.md`'s Decision Revalidation Workflow; doctor's `HOLLOW-COMPLETION` no-code-change exception) | exact `true` only | (no aliases — marker class, exactly like `sweep` and `review_generated`) | absent reads as false |
| `gate_deferred` (canonical repository-gate deferral marker; selector priority after dependencies are satisfied) | `true`, `false` | truthy strings (`yes`/`on`/`t`) → `true`; `no`/`off`/`f` → `false` | `false` |
| `repository_gate_repair` (generated repair marker; selector priority) | `true`, `false` | truthy strings (`yes`/`on`/`t`) → `true`; `no`/`off`/`f` → `false` | `false` |
| `deferred_implementation_base` (late-deferral saved range) | any non-empty commit text | none; trim only | absent |
| `deferred_implementation_merge` (late-deferral saved range) | any non-empty commit text | none; trim only | absent |

**Write paths are unaffected.** Step 1 queue advance, Step 9 finalization, Step 8 follow-up generation, the kb-lessons handoff, and capture emission always write the canonical key and canonical enum value — never an alias, never the typo'd input. The normalize-and-warn contract is read-only.

**Fields with no canonical vocabulary are outside this contract and are read verbatim.** `prime_files`, `write_set`, `required_lessons`, and any future path list have no canonical vocabulary to normalize against — no alias map, no case folding, no path canonicalization, no warning. `assigned_to` joins them for the same reason: a session name is whatever the user or the assigning checkout called itself, so there is nothing to normalize *against* and folding case would silently make two distinct sessions look like one. `stakeholder` is the same class again — a person's name is whatever the user called them, and folding case would make two people look like one; its literal-then-same-person fold match is stated in `actions/capture-reference.md` → Fold-First Rule → Stakeholder-audience questions. A reader takes the strings as written, with one narrow exception that applies to **every** field in this class: **surrounding whitespace is trimmed.** YAML already strips it from unquoted scalars, so it only survives explicit quoting (`" cloud-alpha "`), it carries no meaning in a name or a path, and treating `" cloud-alpha "` and `"cloud-alpha"` as two different sessions would break the skip-and-report over a difference nobody intended. Verbatim here means *no alias map, no case folding, no path canonicalization* — not byte-preservation of padding. (`depends_on`, `related`, and `blocked_by` are likewise row-less here; `depends_on`/`related` carry alias keys, documented on their schema lines above, but their *values* are also read verbatim.)

`deferred_implementation_base` and `deferred_implementation_merge` use that same trim-only projection, but the canonical deferral writer accepts them only as a pair, resolves both to commits, requires a non-empty ancestor range, and persists the full commit IDs. A generic reader retains their scalar evidence; resumption owns the later ancestry and path-drift judgment.

**Optional phase-stamp read contract.** `planning_at`, `dispatch_at`, `builder_handback_at`, `integration_at`, `review_at`, `remediation_at`, `re_review_at`, and `release_at` are additive scalar observations. Absence means that phase was not observed and must render no phase, zero, or inferred instant. A present value is trimmed and parsed by the Timestamp rule; an unparseable value is ignored by phase-duration derivation rather than normalized or replaced. The board keeps the declared pipeline order, omits missing/unparseable phases, and measures each displayed interval from the previous parseable observation. None of these fields changes calibration: its only span remains the just-archived REQ's `claimed_at` → `completed_at`. The same tolerance covers the body: `## Timing` is written by `fold-timing-summary` at Step 8 only when the run recorded lifecycle timing events, so a REQ without that section is complete rather than incomplete and no reader may require it.

### Terminal-success status set

**After applying the `status` alias map, a REQ counts as *terminally successful* when its status is `completed` or `completed-with-issues`.** This is the canonical set every reader that selects "completed work" must honor — `completed-with-issues` is terminal and counts toward UR completion (it just carries known follow-ups, per `actions/work.md` Step 8), so a filter that accepts only the literal `completed` silently drops remediated-with-issues work. `failed` is terminal but **not** successful — success-readers exclude it.

### Dependency-source-ready status set

**A dependency is source-ready when it is terminally successful, or when its normalized status is `claimed` and it carries a nonblank `commit:`.** The latter is a request held for heavy lanes after review: it has landed implementation and passed the fast gate, so dependents may build against the exact committed source while the hold is open. `claimed` without a commit fails closed, and a red drain withdraws `commit:` so the dependency is unmet again. This set is scheduling authority only: it does not make the held REQ terminal, successful, archivable, or eligible for completed-work readers.

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

## Stuck Runs Hand Off to Judgment (any step)

The canonical commands are deterministic on purpose: each decides exactly what it can prove and refuses the rest. A refusal, a typed blocker, or a gate that returns the same finding on retry means the deterministic layer has said everything it can. From there the run is stuck, and what unsticks it is the orchestrator's judgment about *why* the command refused, not another retry of the same command and not a mechanical rule looked up in this file. The "no free-form fallback" sentences elsewhere in the pipeline forbid reproducing a command's mutation by hand: a claim, a finalization, a state transition. They do not forbid thinking about the cause of a refusal and clearing it so the command can succeed. A queue is a factory; a stray file or an unattributed byte must never park it.

Under `do-work run-with-recovery` the handoff is strongest. The user chose that verb to assert sole authority over this queue, so ownership is already decided and the only open question is what the blocking bytes are. Under that verb, stop only for an action that would be destructive or irreversible; everything else is the orchestrator's to decide and do, and to report afterwards.

The first stuck state a run meets is a `recover` refusal, and its `blocked_paths` record is the worked example. For every entry, judge what the path is and take the least destructive action that clears it. The judgment is the rule; these classes are illustrative:

- **Not project state** — operating-system or editor metadata, ignored build output, a file the project's own tooling produced. Delete it, or exclude it locally with `scripts/add-local-git-exclude.sh <path> '**/<path>'`. Never commit it.
- **This session's own earlier write** that discovery cannot attribute, such as a release bump this run started and abandoned. Revert it or finish it so the bytes are what the pipeline expects.
- **The user's own uncommitted work** that lands on a shared path, such as a hand-edited changelog or REQ file. Keep every byte and leave it exactly where it is: name the path in the run's progress output and carry on. Never author a commit to hold it — the pipeline commits only what the REQ declares, so a commit made here attributes somebody else's bytes to this request, can land a half-written draft on main, and takes away the owner's chance to finish it. If a canonical command still refuses because of that path, it is the shared-state class below: stop and name the verb that resolves it.
- **Shared state whose owner you cannot decide** under plain `run` — an interrupted finalization tail with two possible owners, a foreign hunk in a lifecycle file, a claim another checkout may still hold. Stop, and name the verb that resolves it: `do-work run-with-recovery` when the user knows this checkout is the only writer and releaser, `do-work cleanup` when the archive itself needs repair. Under `run-with-recovery` that ownership question is already answered, so decide it.

The same judgment applies to any later typed refusal whose evidence is dirt or an obstacle rather than a pipeline contradiction: a dirty claim target, a release path the command cannot attribute, a gate hold whose cause is outside the REQ. Judgment covers the obstacle, never the mutation. Do not hand-edit REQ frontmatter, reconstruct staging or provenance, or read a REQ's bytes to guess its owner; after each clearing action re-run the record's exact `verification_argv` and let the command decide. Bound the loop: re-run after every action, and a path that survives the action you expected to clear it stops with the record. Report each cleared path and the action taken in the run's progress output.

## Crash Recovery (Step 1)

Run the canonical `recover` command before selection. Its ordered typed result settles finalization first, then classifies every working claim from structural checkpoint evidence; plain recovery preserves claims and returns exact takeover argv, while explicit `--take-over REQ-NNN` or `--assume-sole-authority` authorizes the canonical reset for every claim except one this run set aside, whose claim recovery preserves so the exclusion holds. Follow the result, never infer ownership from prose or interpolate a writer label into shell source.

Keep the exact `## In Progress (interrupted)` heading as the claim-evidence boundary. Explicit authority removes every same-request entry atomically, including multiple writer labels and unlabelled legacy records, while unrelated entries and project dirt remain byte-identical.

## Repository Gate Deferral and Resumption

This is the full action-owned judgment behind `actions/work.md`'s pre-build and final-gate attribution lanes. The canonical gate is mandatory; deferral changes who owns an unrelated failure and what the selector runs next, never the pass requirement for completion.

### One retry before classification

**Whenever a direct run of the canonical gate argv exits non-zero, run the exact same argv once more — immediately, directly, unpiped, from the project root — before anything classifies that failure.** ***Directly* is a condition on what the gate observes, not on who spawns it:** the gate's own output reaches the console and the gate's own exit status reaches the caller, with no pipe, capture or filter in between. Launching that same argv through `run-timed-command` satisfies it — the wrapper hands the child the console's own handles and exits with the child's status, reporting 128 plus the signal number for a signalled child and 127 for a command that never launched — so every lane below may wrap the gate to attribute its wall time without becoming an indirect run. Only the second run decides: a zero exit is a green gate and the run continues through the ordinary zero-exit path, and a second non-zero exit is the real failure whose status, diagnostic evidence, and fingerprint every branch below reads. Exactly one retry per gate run, in every lane that launches the gate directly — the baseline lane and late attribution alike — so a run that has already been retried is never retried again. Report the retry as one line in the run's progress output naming both exit statuses, and record both in the REQ's `## Testing` section. Step 5 pre-flight applies the same rule inside `<skill-root>/tools/checks/preflight.sh`, which records only the rerun's result.

One transient broken pipe inside a single probe used to cost a deferral, a minted repair REQ, and an archived no-op that changed no code. That is the cost this retry removes; a gate that is genuinely red pays one extra run and defers exactly as before.

### Session state and baseline

At run start hold two session-local sets: **suppressed parents** and **repair closure**. They are scheduling evidence, not REQ fields. A parent enters suppression only from a successful typed `gate_deferral` result and stays there until its repair dependency reaches terminal success. Suppression wins over explicit-REQ provenance, preventing a targeted parent from bypassing the dependency it just gained. Every returned repair id enters the closure even when its `user_request` is outside a targeted UR; this is the only cross-UR widening allowed. Recompute the canonical selector after every deferral and every repair terminal result—never reuse a prior selected record.

Before dispatch or source edits, resolve the project-owned canonical gate once as structured argv and invoke `advance` for the exact working request, passing one `--gate-arg` per token. Consume only the matching request-bound `gate_records`, and take their verdict as given rather than re-deriving it: a satisfied green-gate record authorizes dispatch at its returned baseline revision and the baseline is not rerun, a `needs_input` record requires running its `next_argv` directly and unpiped and returning the exact exit status through the same advance phase so `advance` records the result, and a failed or mismatched record is unverifiable and stops safely. **Fingerprinting is the action's, and the same procedure is used everywhere below:** for a direct run save the current revision, the direct status, bounded diagnostic evidence, and a stable semantic fingerprint that discards volatile timestamps, scratch roots, and ordering noise while retaining the failing command/test identity and normalized diagnostic. When that run was retried under **One retry before classification** above, all four come from the second run. A gate launch or record failure stops safely.

The branch table is exhaustive:

| Claimed REQ and baseline | Action |
|---|---|
| Ordinary REQ, recorded match or direct exit 0 | Save `gate_evidence.baseline_revision` for a match or the freshly recorded revision for a direct run; dispatch normally. |
| Ordinary REQ, still non-zero after its one retry | Defer before source edits through the canonical transaction, consume its typed result, suppress the parent, extend repair closure, and select again. |
| `repository_gate_repair: true`, matching recorded red fingerprint | This is the repair's authorized baseline. Implement without recursive deferral. |
| Repair, exit 0 | Complete through the reviewed no-change path; the defect was repaired elsewhere and parents may resume. |
| Repair, different red fingerprint or launch failure | Stop/fail the repair safely; do not manufacture a second repair from it. Parents remain dependency-gated and unrelated ready work continues. |
| Ready `gate_deferred: true` parent without a saved pair | Require a fresh green baseline, then dispatch normally. |
| Ready deferred parent with a saved pair | Run the resume proof below before deciding whether builder work is reusable. |

### Manifest authoring and collision retry

The action authors exact evidence; `do-work-cli defer-gate --manifest <path>` alone mutates. Copy parent/checkpoint preimages to payload files, carry the exact writer label, structured gate argv, direct non-zero status, fingerprint/evidence, stable root-cause `sweep_key`, and optional paired implementation commits. Scan the same request and `.req-reservations/` evidence as capture, propose read-only max+1, and let `defer-gate` exclusively create the unpadded reservation. Do not call a helper that pre-creates a marker. Retry a collision only for one of two typed results: **(a)** pre-mutation `outcome: refused` with the collision finding, an empty `changes` list, and `rollback.status: not_needed`; or **(b)** post-mutation collision with `outcome: rolled_back` and `rollback.status: succeeded`. Rescan the live repository and propose its new max+1 before retrying. An incomplete/failed rollback, `committed_risk`, any non-collision refusal/finding, stale preimage, or non-empty refused-result changes stops without retry.

On success consume `gate_deferral` fields only: `parent_id`, `parent_path`, `repair_id`, `repair_path`, `repair_outcome`, `repair_dependency`, `diagnostic_fingerprint`, `sweep_key`, command/status, and optional range. Never scrape text rendering or infer the repair from queue order. Validate that the returned parent and fingerprint equal the proposal before updating the session sets.

### Late attribution

The final gate uses the identical argv and fingerprint procedure and never substitutes the baseline record for a direct run. Run the gate directly and unpiped; a non-zero exit gets its one rerun under **One retry before classification** above before any branch here applies, so every branch below reads the second run's status, output, and fingerprint. Return the exact argv and exit status through the current advance phase; require its request-bound green-gate record before continuing. Because recording happens after the gate returns, the recorded revision includes any gate-created `_dev/gate-runs/` log commit. In worktree dispatch mode, a red current tree is attributed by creating an isolated detached diagnostic worktree at the saved `<pre>` and running the gate there directly. Always remove that diagnostic worktree without force after capturing its status and evidence.

- Active `repository_gate_repair: true`: any same-fingerprint red, different-fingerprint red, missing/malformed fingerprint evidence, or launch-failed final gate is a terminal repair failure. Prepare a strict `transition: "fail"` finalization manifest, never invoke `defer-gate`; every parent remains pending behind the failed dependency. Recompute the selector and continue unrelated runnable REQs.
- Base exit 0: current implementation caused the failure; use the bounded remediation loop.
- Base non-zero with the exact current failure fingerprint: unrelated failure; call `defer-gate` with full `<pre>` and `<merge_hash>` commits. After success, the normal archive path will not clean the builder, so immediately run `git worktree remove <path>`, `git branch -d <operative_name>` from the integration branch, then `git worktree prune`. A refusal stops; never add force.
- Base fingerprint mismatch, diagnostic launcher failure, missing/unresolvable range, or inability to isolate the saved base: attribution is unverifiable and stops safely.

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
- **Recorded green revision:** <`recorded_revision` returned by the advance green-gate record>
- **Observed result:** green before implementation; repair already satisfied
- **Verified at:** <now> (current UTC instant — Timestamp rule)
```

Write the mandatory summary exactly as:

```markdown
## Implementation Summary

**Files changed:** None — verified repository-gate repair no-op.

**What was done:** Re-ran the repair's recorded canonical repository gate before source edits and confirmed it is already green; no implementation changes were necessary.
```

**The pre-build run is this branch's only gate lane** — at most two runs of the same argv under the one-retry rule above. After its direct zero exit, return the exact status and argv through the current advance phase, require a satisfied request-bound green-gate record, and write its `recorded_revision` into the evidence block above. Qualification and independent review then each invoke `<skill-root>/tools/do-work-cli.sh --repo-root <project-root> --format json validate-already-green-repair --request-path <exact working REQ path> --writer <exact finalization writer> --at <now>` and consume only their respective `already_green_repair.tdd_allowed` and `already_green_repair.review_allowed` fields; **callers never supply the evidence values, recreate the predicate, relaunch the gate, or run the ordinary diff-requiring qualifier**, and self-review is insufficient. Append `## Qualification\n\nPassed — repository-gate repair no-op; durable gate evidence verified and project diff empty.`, then author the ordinary strict finalization manifest with the already-green evidence and no release payload. Because a no-op repair changes no project path, a later `HEAD` move can never make it the cause of a red gate, and its recorded evidence stays verifiable at its own recorded revision. An ordinary REQ's green record has no such immunity: it can stop matching as `HEAD` moves, and a mismatched record is unverifiable and stops safely (**Session state and baseline** above). Done means that one gate lane plus bookkeeping, reported in the REQ's `## Testing` section, which carries both exit statuses when the retry fired. **No other empty implementation receives any of these exceptions.**

### Continuation and reporting

A failed, cancelled, or still-gated repair never releases its parents: the selector excludes them by `depends_on`, and the action continues unrelated selected REQs instead of ending the run; a successful repair instead triggers a fresh selection, which is what lets its parents resume. Run summaries compose rather than overwrite — report each deferred parent with its typed repair id and create/fold outcome, each no-change repair completion, each resumed parent as reused or rebuilt-after-drift, each repair failure or cancellation, and every unrelated REQ that continued afterward, with dependency-gated parents still under the composed summary's blocked-by-dependencies section. No branch writes `blocked` or `pending-answers`, and no branch asks for user confirmation.

## Worktree Dispatch Mode (Step 1)

**Every run mode, the serial default included.** Each builder runs on its own per-REQ branch, in its own git worktree wherever the harness supports one (*Isolation ladder*, below); the orchestrator merges those branches and remains the only writer of `do-work/` state. Worktree isolation is what makes the ownership boundary hold (**Execution Model — Claim Anywhere, One Releaser**, above): it changes *where* a builder writes — its own tree instead of the main one — so a builder can neither touch queue state nor collide with a sibling, and an interrupted build leaves a branch to merge or discard rather than half-written files in the main tree. **Everything in this section is written per REQ and therefore already holds for any number of concurrent builders** — one `<operative_name>`, one hand-back sequence, one `<pre>..<merge_hash>` range, one cleanup, each per REQ. Fan-Out Dispatch (below) adds only who picks the set and what never parallelises.

**Isolation ladder — three rungs, and only the last one has no isolation.** Probe `git worktree list` (a non-zero exit, or an unrecognized subcommand, means no worktree support) and confirm the harness can run an agent against a working directory you choose.

- **Both present — worktree rung.** Everything in this section applies as written: `git worktree add` under *Naming* below, the builder in its own directory, `git worktree remove` at cleanup.
- **Either missing — branch rung.** The second directory goes; the isolation stays. Create `worktree-agent-REQ-NNN-<suffix>` as a plain branch, so the builder gets its **own per-REQ branch in the shared working tree**: check that branch out, implement and **commit** on it, then check the integration branch back out and run the hand-back sequence below unchanged. The merge, the `<pre>..<merge_hash>` range, qualification, review, post-merge verification and cleanup are identical; cleanup is just `git branch -d <operative_name>` with no worktree to remove. **Reuse that same name — never invent a second prefix:** the sweep in **Crash Recovery (Step 1)** above and `actions/cleanup.md` → **Pass 5: Orphaned Worktrees (consent-gated)** both enumerate `git branch --list 'worktree-agent-*'` independently of worktrees, so a branch-rung leftover is discoverable by the sweeps that already exist. **Builder and orchestrator are roles here, not processes:** one agent usually plays both, and the boundary between them is the branch, not a second process — what the builder writes is committed on this REQ's branch, and every `do-work/` write is made as the orchestrator, on the integration branch, and never committed on this REQ's branch (the queue guard in step 2 of the hand-back sequence is the mechanical check). **This rung isolates committed work only.** One working tree means an uncommitted edit is still sitting in it for the next REQ, so setting a REQ aside on this rung requires the builder's edits committed on its branch first; an uncommitted set-aside is not isolated and must never be reported as one.
- **Not a git repository — no rung.** There is no branch to isolate onto, and `actions/work.md` Step 9 already skips its Commit Phase here. Say so in one progress line — `⚠ REQ-NNN: no git repository, so this implementation is not isolated; setting it aside would leave its edits in the working tree` — and continue.

**The concurrency degrade stays silent; a lost isolation does not.** Dropping from the worktree rung to the branch rung is not reportable, because the REQ is still isolated, and neither is running one REQ at a time because the harness cannot run an agent against a directory you choose: unlike `../../do-work-board/actions/board.md`'s Go check — which reports and stops, because there the toolchain *is* the capability — neither absence is an error and neither may surface as one. The no-rung case is the only genuine loss of isolation, and it is the only one that surfaces.

**A builder tree does not have to be a worktree.** The definition a builder has to satisfy is *own tree, own branch, hands back a branch* — a spawned `git worktree`, a second local workspace, a clone, or a remote/cloud sandbox all satisfy it, and everything downstream of the merge is identical for all four because the orchestrator integrates a branch either way. Three deltas are worth stating, because they are the only places the shape matters:

- **The naming and cleanup mechanics below are worktree-specific.** `git worktree add`, `git worktree remove`, `git worktree prune` and the `worktree-agent-REQ-NNN-<suffix>` convention apply to a spawned worktree. A workspace, clone or sandbox the *user* already opened is not the orchestrator's to name or delete; it hands back a branch and the orchestrator merges it under whatever name it arrived with. **Hold that name as this REQ's operative name exactly the same way** (*The name actually created is this REQ's operative name*, below) — the merge, and any reporting, still needs one string.
- **A remote builder's hand-back travels on the branch.** The absolute-main-tree-path hand-back file (*Sole integrator*, below) is a **local-only** mechanism: it works because the builder shares a filesystem with the main tree. A builder that does not — a cloud sandbox, another machine — commits its manifest **on its own branch** and the orchestrator reads it after the merge. Never hand a remote builder a main-tree path; it resolves to nothing or to something else's.
- **A non-releaser checkout's `do-work/` snapshot is potentially stale, and never authoritative.** Where a consumer commits `do-work/`, every non-releaser checkout carries a copy of the queue as of its last sync: a REQ it shows as `pending` may already be claimed elsewhere, and a `status` it shows may be several transitions behind. Read it as a hint about what exists, never as the current state of anything, and resolve disagreements by syncing rather than by writing. This is the same rule *State stays home* (below) applies to a worktree's snapshot, widened to the checkout that owns one.

**Claim conflicts between checkouts are ordinary git conflicts, and `do-work/CHECKPOINT.md` is where they surface.** Two checkouts claiming the same REQ produce a plain **content** conflict on the REQ file — never a rename conflict, because both sides perform the identical `do-work/queue/` → `do-work/working/` rename and git resolves it silently — made entirely of the two claim writes (`claimed_at`, and whatever sections each side generated). With byte-identical claim writes the REQ file does not conflict at all, and the `writer:` label is then the *only* thing that surfaces the double claim. **Expect a `CHECKPOINT.md` conflict on every concurrent claim, including two that overlap in nothing** — two single-line appends land at the same position, so git reports `CONFLICT (add/add)` (`AA`) or `CONFLICT (content)` (`UU`) depending on whether the file already existed. **Resolve it by keeping every entry from both sides.** That is the merge-time reading of the checkpoint's own contract (**In-Progress Record (Step 1)**, below): every live claim survives. Both one-sided resolutions lose data: taking ours strips another checkout's live claim record by hand; taking theirs discards this checkout's own record and makes its own crash unrecoverable. One id under two labels needs no reconciling: it is the honest record of two checkouts. On the REQ file itself only one claim survives — a human decision, evidenced by which checkout actually has the work. All of that is observed behavior, recorded in `do-work/archive/UR-018/REQ-095-two-clone-acceptance-run.md`'s `## Testing` section, and it holds only where the consumer **commits `do-work/`**; where the directory is untracked nothing syncs, so no conflict surfaces and `duplicate-req-id` is the only detector left.

**Red Flag — a second checkout running the release tail.** Capturing, claiming and building from another checkout are now in contract; the violation to watch for is a second checkout merging integration, bumping the version, prepending to `CHANGELOG.md`, moving files into `archive/`, or closing a UR. Two changelog prepends against one queue is the shape it usually takes, and unique version numbers do not make it safe (**Serial-only**, below).

**Naming.** The worktree directory's basename and the branch name are the **same string**: `worktree-agent-REQ-NNN-<suffix>`. Embedding the REQ id is what lets any later sweep correlate a leftover with its REQ by name alone; sharing one string is what makes a single grep find both the directory and the branch. Derive `<suffix>` from the REQ's filename slug, lowercased and reduced to `[a-z0-9-]`, **as a text operation before you compose any shell command** — never pipe raw REQ text through `tr`/`sed` inside a quoted command line, where an apostrophe breaks the quoting and becomes an injection vector. **On a name collision at creation** — `git worktree add` or branch creation fails because `worktree-agent-REQ-NNN-<suffix>` already exists as a leftover (typically a crash-recovered REQ re-dispatching into the name its own reported-but-not-deleted leftover still holds) — do **not** delete or force. Dispatch under a fresh unique variant: append an **incrementing numeric token** to the suffix — `-2`, then `-3` if that collides too — keeping the `worktree-agent-REQ-NNN-` prefix intact so both names still correlate to the REQ by name (this is what the sweeps grep for). **One scheme, not a choice:** a free pick between a counter and a timestamp token lets two runs shape the same collision differently, so the name stops being readable evidence of what happened. Report the coexistence and leave the original leftover to its owners — the crash sweep above if it turns out merged, `actions/cleanup.md` → **Pass 5: Orphaned Worktrees (consent-gated)** if unmerged.

**The name actually created is this REQ's operative name.** Whatever `git worktree add` succeeded with — the derived name in the common case, the variant after a collision — is the one string every later worktree/branch operation uses: the hand-back merge's `git merge` argument (below), Step 8's `git worktree remove` and `git branch -d` (*Cleanup — happy path*, below), the crash sweep's own-session bookkeeping (**Crash Recovery (Step 1)**, above), and anything reported back to the user. This reference calls it **`<operative_name>`** wherever a command needs it; the worktree path is that same name under the worktrees parent directory. Hold it exactly the way `<pre>`/`<merge_hash>` are held — known from this session's own context and re-typed as a literal into each fresh command, never a shell variable (**Hold both endpoints as re-typed literals**, below). Nothing persists it, because nothing outside this session's run consumes it: the sweeps and `actions/cleanup.md` Pass 5 discover leftover names by enumerating git, never by re-deriving them. **Re-deriving the name from the slug at cleanup time is the failure this closes** — after a variant dispatch the derived string names the *leftover*, so `git worktree remove`/`git branch -d` target unmerged work, refuse, and halt the run on a false "the merge was skipped or lost" alarm while the variant worktree is never cleaned at all. **With no collision the operative name *is* the derived name**, so the common path behaves exactly as before.

**Where worktrees live: outside the repo working tree.** A sibling directory (`../<repo>-worktrees/worktree-agent-REQ-NNN-…`) or a scratch directory — never nested inside the repo. Inside, it is a second checkout of the repo sitting in the repo: `actions/cleanup.md` Pass 3a scans the filesystem for any `do-work/` directory outside the project root, and where the consumer commits `do-work/` it would find the builder's checkout of one and try to relocate it into the canonical queue. The extra tree also reads as stray untracked residue to every status and stray-file check downstream. That is a corruption path, not just untidiness.

**State stays home.** **Every path under `do-work/` exists in the main tree only and is the orchestrator's** — the queue, `working/`, `CHECKPOINT.md` and `runs/` are examples of the rule, not its extent, so a directory added later is covered the moment it exists rather than when someone remembers to list it. Builders receive their REQ brief in the dispatch prompt and never read or write queue state from inside a worktree. Untracked files do not propagate into a new worktree, so on the common install (where `do-work/` is untracked) a builder simply finds nothing there. The trap is the other install: where a consumer **commits** `do-work/`, the worktree carries a *stale snapshot* of the queue as it stood at the branch point. Treat that snapshot as absent — never read a status from it, never write to it. Every claim, status flip, and archive move happens in the main tree, by the orchestrator. **On the branch rung there is only one tree, so this is a role boundary rather than a filesystem one:** the same agent may be builder and orchestrator, and what keeps the rule mechanical is that no `do-work/` path is ever committed on a builder branch — step 2's queue guard below is what proves it.

**Sole integrator.** The builder never writes the main tree or its branch, **with exactly one exception: its own `do-work/runs/work-<YYYY-MM-DD-HHMMSS>/REQ-NNN-handback.md`**, reached by the absolute main-tree path the orchestrator hands it (the repo-relative trap below cuts both directions) and **never staged, committed, or merged** — it is an orchestrator-owned working file, not branch content. That is one path, derived from the builder's own REQ id: a sibling's hand-back, `manifest.md`, anything else under `do-work/runs/`, and every other main-tree path remain violations. The exception exists because the hand-back has to survive a dead transcript (`crew-members/background-agents.md`); without it the mandatory file has nowhere legal to go. It commits on its own branch and hands back its file manifest; the orchestrator merges. A shared file needing one line of wiring — a `<link>` in a shared template, a registry entry, an export barrel — is an **integration seam**: the builder hands back the exact line and where it goes, and the orchestrator applies it in the main tree **inside the merge commit** — step 3 of the hand-back sequence below, which is what keeps the seam inside this REQ's merge range. A builder that edits the seam itself writes the main tree the orchestrator alone owns.

**Merge, never rebase.** Integrate with `git merge --no-ff <branch>` (the full invocation, which adds `--no-commit` so the integration seam has somewhere to go, is step 2 of the hand-back sequence below). Rebasing rewrites the builder's commits, so `git branch -d` no longer recognizes the work as merged and the free merged-ness assertion below is destroyed. `--no-ff` also preserves the merge commit as the "integrated by orchestrator" provenance record.

**When to merge, and the range every evidence step reads.** The orchestrator merges each builder branch **at hand-back — end of Step 6, before Step 6.25 (Implementation Summary)** — so every downstream evidence step (6.25, 6.3, 6.5, 7, 8, 9) observes one integrated main tree, not a clean one. Any position after 6.25 leaves qualify (Step 6.3) and review (Step 7) reading a clean main tree with nothing to check. Run this hand-back sequence on the integration branch — step 0 first, then the four that produce the merge:

0. **Settle owner-written run artifacts and the index first, then capture `<pre>` — in that order.** Queue-mode `advance` already committed the queue-to-working move plus `do-work/CHECKPOINT.md`; never stage either here. First read `git status --short --untracked-files=all -- do-work/` and sort every path into exactly one category: **stage** the exact owner-written run artifacts (`manifest.md` plus each `REQ-NNN-brief.md`); **allow but never stage** each expected `REQ-NNN-handback.md` named by that run; **leave alone and name** every other `do-work/` path, including claim or checkpoint dirt this REQ did not produce. The hand-back is scratch the builder writes to the main tree, so excluding it from the stage set must not reclassify it as an error. A third-category path belongs to another session or to the user: print it once in the run's progress output, keep it out of the stage set, and continue — it is never a reason to stop, never reported as a blocker, and never committed here, because **the pipeline commits only what this REQ declares.** Only a dirty path **this REQ itself owns** still stops: its run artifacts, its working REQ file and checkpoint entry, its `write_set`, its release paths — the bytes the merge and Step 9 are about to read. Check the index too (`git diff --cached --name-only`), because every commit in this sequence takes the whole index and would adopt anything staged in it. A staged path outside the stage set is taken **out of the index** with `git restore --staged -- <path>`, which leaves the working-tree bytes exactly as they are, and named alongside the third category; unstaging is not a stop and not a repair. **One case is not lossless and the progress line must say so:** where the path was staged and then modified again (`MM`), the staged snapshot differs from the working tree, and unstaging drops that snapshot to a dangling blob reachable only through `git fsck --lost-found`. The working tree still holds the owner's newest bytes; an intermediate version they staged deliberately does not survive. Name that case explicitly rather than reporting it as an ordinary unstage. Stage only the exact run-artifact paths with `git add -A -- <exact-run-artifact-paths>`, re-run the cached-name guard, then use plain `git commit`. **Never use `git commit -- do-work/`:** a path-limited commit takes tracked paths from the working tree and ignores the index content the guard just inspected, so an unstaged edit can be swept in while an untracked run artifact is left behind. If the exact stage set has no changes, skip the commit and continue. Where `do-work/` is untracked, there is no run-artifact commit. **Step 0 ends with the index holding nothing but this REQ's own stage set** — that is what makes the plain `git commit` safe — and never requires a clean working tree. Ordering is load-bearing: any run-artifact commit must land *below* `<pre>`, so the next step's capture keeps it outside `<pre>..<merge_hash>`; commit after capturing and the owner's own artifact falls inside the merge range, where qualify and review read it as an undeclared touch of `do-work/`.

1. **Capture `<pre>`** — run `git rev-parse --short HEAD` and read the printed hash. It is the integration tip before this REQ's **first** merge and the lower bound of the merge range. Capture it **once per REQ**; a remediation re-merge does not re-capture it (below). Never recover it afterwards as `HEAD^1` or live `HEAD` — both move as soon as the orchestrator commits the changelog (Step 9) or merges the next sibling.
2. **Guard the queue, then merge without committing** — first run `git diff --name-only <pre>...<operative_name> -- do-work/` (three dots: merge-base to branch tip — which is why this runs **before** the merge; once the branch is merged that diff goes empty). Owner bookkeeping sits below `<pre>` (step 0) and the hand-back file is never committed, so any path printed is queue state committed on the builder's branch — the write *State stays home* forbids. Stop and drop/revert those commits on the branch before integrating (on a remediation re-merge whose fix branch was cut from the integrated tip, owner or sibling commits can surface here too — the cumulative range's safe-direction over-inclusion: judge, don't auto-delete). This guard is the only mechanical check that ever sees them: `tools/checks/qualify.sh` and `tools/checks/scope-drift.sh` both exclude `do-work/` by contract, and `queue-kanban verify`'s committed-queue-state probe reads this same three-dot diff — blind after the merge, and Step 8 deletes the branch before Step 9 runs verify. Then `git merge --no-ff --no-commit <operative_name>` (this REQ's operative name, *Naming* above — the branch the builder was actually dispatched on, which is the collision variant where there was one), then resolve any conflict. `--no-ff` forces a merge commit even where the branch could fast-forward (Merge, never rebase, above); `--no-commit` is what leaves the integration seam somewhere to go. If git answers `Already up to date.` the builder committed nothing: no `MERGE_HEAD` is set, so a `git commit` here would either fail or fabricate a non-merge commit — stop and treat the hand-back as empty instead.
3. **Apply the integration seams, then commit** — stage the handed-back seam lines (Sole integrator, above) and `git commit`. Folding the seam into the merge commit is the only placement that puts it inside the merge range: a seam committed *after* the merge is that merge commit's **child**, hence outside `<pre>..<merge_hash>`, and qualify, review, and Step 9's validation would never see it.
4. **Capture `<merge_hash>`** — `git rev-parse --short HEAD` on the commit just made. It is the upper bound of the merge range and the supplied-provenance hash that finalization records in the REQ's `commit:` field.

**Hold both endpoints as re-typed literals, never as shell variables.** The canonical [State across command blocks](../docs/prescribed-shell-primitives.md#state-across-command-blocks) rule applies because the consumers sit in later blocks with model round-trips in between — a `"$pre..$merge_hash"` composed in a fresh shell expands to `".."`, which git rejects. Hold both hashes known from this session's own context and re-typed into each fresh command, never a shell variable. `tools/checks/qualify.sh` hard-FAILs on a range it cannot resolve rather than reading an empty diff, so a lost endpoint surfaces as a qualification failure naming the range instead of a vacuous pass.

**The merge range is `<pre>..<merge_hash>`**, and in worktree dispatch mode it replaces the working diff wherever an evidence step reads changes: **Step 6.3** (`tools/checks/qualify.sh` via `DO_WORK_DIFF_RANGE="<pre>..<merge_hash>"`), **Step 7** (review's Get-the-Diff), **Step 8** (post-merge verification, below), and **Step 9's** finalization-allowlist validation all consume it. Use the *captured* `<merge_hash>` as the upper bound, never live `HEAD`: HEAD moves when finalization commits its lifecycle/release group, so `<pre>..HEAD` would sweep in that commit and misattribute it to this REQ. `<pre>` is `<merge_hash>`'s first-parent ancestor — its direct first parent after a single merge — so `merge-base(<pre>, <merge_hash>) == <pre>` and this is already the merge-base form.

**Remediation re-merges: the range is cumulative.** A failed review sends the REQ back to Step 6 (`actions/work.md`) and the builder's fix branch is merged again. Repeat steps 2–4 above but **not step 1**: keep the first `<pre>` and re-capture only `<merge_hash>` from the newest merge commit, making the range `<pre₁>..<merge_hash₂>` — first pre-merge tip to latest merge — which covers the original work *and* the fix. Re-capturing `<pre>` would instead give `<pre₂>..<merge_hash₂>`, covering only the fix delta: review would read the fix in isolation and every originally-touched file would WARN as listed-but-not-in-the-diff. **Step 9 records the latest `<merge_hash>`** in `commit:`. The cumulative range's one cost is over-inclusion — any orchestrator commit that landed on the integration branch between the two merges falls inside it and surfaces as an undeclared touch for your judgment. That is the safe direction of error: it shows up and gets judged, where under-inclusion silently hides the REQ's own work.

**Post-merge verification, before finalization.** The builder verified its own branch; nobody has verified the merged result. Re-run the REQ's acceptance checks — the `## Red-Green Proof` GREEN condition, the project's test command, whatever `## Scope` named — against the merged main tree **before** Step 9 finalizes anything (`actions/work.md`). A green builder branch must not compose into a red main that archives as done: **the unit you verify is the unit you roll back**, and a red merged state is not an archive-plus-follow-up — stop, revert to the last verified state, and re-dispatch.

**Cleanup — happy path (Step 9, after typed finalization success).** After finalization reports `cleanup_complete`, remove the builder's worktree and branch **by this REQ's operative name** (*Naming*, above — never re-derived from the slug here): `git worktree remove <path>` (no `--force`, where `<path>` is the worktree whose basename is `<operative_name>`), then `git branch -d <operative_name>`, then `git worktree prune`. Run `branch -d` **from the integration branch you merged into**: `-d` tests merged-ness against the current HEAD (or the branch's configured upstream), so from anywhere else a perfectly-merged branch refuses and an unmerged one can pass — "refusal = unmerged" silently becomes "refusal = wrong branch." Both refusals are signal, not friction: `worktree remove` refuses on a dirty worktree (uncommitted builder work you are about to lose), `branch -d` refuses on an unmerged branch (a merge that was skipped or lost). **Never `-D`, never `--force`.** Report the refusal and stop — forcing destroys the only evidence that the integration didn't happen.

**Cleanup — crash path.** Leftovers from an interrupted run are swept by **Crash Recovery (Step 1)** above: only a clean, merged leftover whose exact REQ is positively settled outside `do-work/working/` is removed mechanically. Dirty, unmerged, still-working, absent, ambiguous, malformed, or unreadable cases are reported and never auto-deleted. Discarding any such builder lane belongs to `actions/cleanup.md` → **Pass 5: Orphaned Worktrees (consent-gated)**, which asks first and only acts when a human can answer.

**Fan-Out Dispatch — several builders, one releaser.** Every guarantee above is already per REQ, so raising the builder count adds no coordination and no durable state.

**Reached two ways, and the default is neither.** `actions/work.md` processes one REQ at a time unless it is asked not to: `do-work run --fan-out [N]` puts it in **auto-wave mode**, where the loop computes the ready set itself and dispatches builders with **no confirmation gate** (*Auto-wave*, below). Without that flag the action runs the serial loop — one REQ at a time, each on its own branch or worktree exactly as this section describes — and a human or advanced harness can still drive the concurrent path by hand — which is what every step here describes, and what auto-wave automates rather than replaces. The floor is why the flag exists instead of the behavior: `actions/work.md` must stay followable by the simplest agent that can read files and run shell commands, so concurrency must be something a reader opts into rather than something sitting in front of them. `--wave N` remains a *scoping* flag — it chooses which dependency depth runs — and composes with `--fan-out`, which chooses how many of the chosen set run at once. The dispatch *mechanism* stays deliberately unspecified (below) in both modes.

What fan-out adds:

- **The set is either picked by a human or computed by auto-wave; `write_set` never gates the first and is not read at all by the second.** In the manual path a human chooses which REQs run together; in auto-wave the loop delegates the whole computation to queue-mode `advance` and dispatches what it claims (*Auto-wave*, below, states the policy around that delegation; the predicate itself is the command's, and no list here is a second copy of it). A REQ's declared `write_set` is **advisory input to that pick, never a gate** in the manual path, and **not read at all** by the computed one — it is display-only, nothing schedules on it, and the board's `overlaps` badge misses glob-vs-glob, `**`, and directory entries, so **absence reads as unknown, not safe** (`../../do-work-board/actions/board.md`). That is why it cannot be a scheduling input: a field whose absence means *unknown* can only ever inform a judgment, never make one.
- **The non-interference proof is the merge, not the pick.** `git merge --no-ff --no-commit` refusing is the only mechanical evidence that two builders' work does not collide. **Its limit is honest: git detects conflicts by line proximity, not meaning.** Two REQs each appending an entry to a shared registry merge cleanly and can still be jointly wrong. The **integration seam** rule (*Sole integrator*, above) is what covers that — and it works only because one integrator applies every seam by hand, inside the merge commit.
- **Integration is serial.** Implementation parallelises; merge → qualify → test → review → changelog → archive runs one REQ at a time. Each merge also invalidates the previous REQ's *Post-merge verification* (above), so those checks re-run per REQ against the tree as it then stands. Expect the wall-clock saving in the build phase only, and say so rather than promising more.
- **A worktree per builder is mandatory, not optional.** Sharing one working tree was considered and ruled out: every test run, qualification check, and review diff would then read a tree carrying the other builder's unfinished edits, so each REQ's evidence steps stop meaning anything and nothing downstream can tell. (The staging race is the lesser problem.) Keep this reason here — without it a later reader re-offers the shared tree as a simplification.

**Auto-wave — what the loop computes, and what it deliberately does not.** `actions/work.md` Step 1 delegates the complete computation and claim transaction to queue-mode `advance`: its typed claimed members are the wave, and its typed exclusions are the composed reasons. **This prose states policy and must never become a second queue scan.** Readiness, the dependency-source-ready set, the live-claim veto, the `assigned_to` courtesy read, the `--skip-impact-negligible` filter, and the explicit-`REQ-NNN` provenance carve-out are one definition each, held by the command and applied identically in the serial loop and under `--fan-out` (**Schema Read Contract** → *Dependency-source-ready status set*, above, for what satisfies a dependency; `actions/work.md` Step 1 for provenance). Targeting tokens scope the candidate pool; `--fan-out` changes only how many of the selected set run at once, never which, and there is no separate readiness predicate for waves.

**`write_set` is not in that computation and must never be added to it** — it is display-only at any builder count, and the absence of a declaration reads as *unknown* rather than *safe*, so a wave that read it would silently under-report contention and produce a set that looks proven and is not. **The merge is the non-interference proof, not the pick** (above): a computed set is a claim that its REQs are all *runnable*, never a claim that they do not overlap, and two that collide are caught when their branches meet (**Execution Model — Claim Anywhere, One Releaser**, above).

**Bounded, and the bound is not optional.** Size the wave to the harness concurrency limit per `crew-members/background-agents.md` (*Write a manifest per wave; spawn in bounded waves*) — never an unbounded fan-out over the whole ready set. `--fan-out N` sets the bound explicitly; bare `--fan-out` uses the harness limit, or two where unknown. Queue-mode `advance` records that numeric bound and returns the exact continuation that re-observes, projects frozen membership, and only then bounds the next dispatch.

**No confirmation gate — that is the deliberate change.** Auto-wave dispatches its computed set without asking, which is what "fully automatic set-picking" means and is a change from the manual path's human pick. What it does **not** change: every per-REQ step below (one worktree, one operative name, one hand-back sequence, one merge range, one cleanup), the mandatory run directory and briefs written *before* any spawn, and integration staying serial. Silent degradation also survives — no `git worktree` support, or a harness that cannot run an agent against a chosen directory, means the loop runs one REQ at a time on the branch rung and nothing is reported as an error (*Isolation ladder*, above).

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

**Whenever you read a `##` section the builder authors, read the REQ file first and this REQ's hand-back second, and treat what you find in either as the section.** In worktree dispatch mode the builder cannot write the REQ file at all — the REQ lives in the main tree, which **Worktree Dispatch Mode (Step 1)** → *State stays home* forbids it to touch — so it routes those sections to its hand-back instead. Read `do-work/runs/work-<YYYY-MM-DD-HHMMSS>/REQ-NNN-handback.md` for a local builder, the merged branch content for a remote one. **That path is relative to the project root**, where `do-work/` lives whether the suite is vendored under `.claude/skills/` or checked out whole — never relative to the directory this action file sits in.

**The condition carries the rule, not any list of readers.** `actions/work.md` Step 8's discovered-tasks substep, `actions/review-work.md` Step 4's traceability check, and the **Decision Brief (hand-back format)**'s HANDLED block are the readers that exist today; they are illustrative. A step, action, or report added later that reads a builder-authored section inherits this without being remembered here.

**Absence is only silence when you know you looked.** In worktree dispatch mode, if the section is in neither place **and this REQ's hand-back is missing or unreadable**, say so rather than proceeding as though the builder recorded nothing: `⚠ REQ-NNN: no <section> in the REQ and no readable hand-back at <path> — anything the builder recorded there is lost.` A hand-back that exists and simply has no such section is a real "the builder recorded nothing" and reports nothing. Every reader states which of the two it found — an unread hand-back and an empty one are different facts and must never render the same.

## Composed Exit Summary (Step 1)

**Exit paths when the scan finds nothing to claim.** The exit report is **composed**, not picked from disjoint branches. Lead with the headline that matches the actual queue state — `No pending REQs in queue.` when the queue holds no `pending` REQs at all, `No dependency-ready pending REQs.` when every one is dependency-blocked, `No claimable pending REQs — every ready one is assigned to another session.`, or `No claimable pending REQs — every ready one is impact-negligible and --skip-impact-negligible is set.` — then append every section below that holds at least one REQ, in this order. **That condition is the rule**, so a category added later inherits it without anyone re-counting, and the headline never strands the user because its matching section always enumerates them.

**Which REQs fall in which category is the typed result's answer, not a second scan.** Each trigger below names evidence the run already holds — `advance`'s queue exclusions and holds, and `recover`'s finalization records — and this table says only what to render from it. Each row's Section cell carries the **exact heading string** to render — write it as written, with `N` replaced by the count — and under it one indented `REQ-NNN — [title] (<fields>)` line per REQ.

| # | Section | Trigger | Fields on each line | Remedy line |
|---|---|---|---|---|
| 1 | Finished-awaiting-archive section — `⚠ N finished REQs awaiting archive (UR-137: 3 REQs, UR-138: 1 REQ, ...):` | a REQ in `do-work/queue/` normalizes to a terminally resolved status | the resolved status, written `(complete → completed)` where an alias was normalized; group the heading by `user_request` (`UR-137: 3 REQs, UR-138: 1 REQ, …`) | `Run do-work cleanup to archive completed work, then do-work recap to see full history.` |
| 2 | Pending-answers section — `⚠ N REQs awaiting clarification:` | status `pending-answers` | status only — frontmatter at this stage; counting `## Open Questions` is `do-work clarify`'s job | `Run do-work clarify to batch-review the open questions; resolved REQs flip to pending and re-enter the queue.` |
| 3 | Blocked-on-external-condition section — `⚠ N REQs blocked on external conditions:` | status `blocked` | `blocked by: <condition>`, the age from `blocked_at`, and `[probe failed this run \| no auto-probe]` — Step 1 already re-ran every `blocked_check` | `When a condition is satisfied, re-run do-work run (a blocked_check is re-probed automatically and unblocks on exit 0), or confirm a human-checkable one via do-work clarify. To give up on one, do-work abandon REQ-NNN.` |
| 3a | …the stakeholder form of row 3 | that REQ also carries `stakeholder:` | `questions for <stakeholder> (N open, K irreversible; since <age>) — report: <latest bundle path from blocked_by>`, led by `⚠ IRREVERSIBLE` when K > 0 | row 3's remedy plus `To ingest a stakeholder's reply, run do-work stakeholder-answers REQ-NNN — share the report path with them first if you haven't.` |
| 4 | Blocked-archive-collision section — `⚠ N REQs held by archive-collision guard:` | status `blocked-archive-collision` | the queue file path, then indented `already archived at <archive-path>` and `recover: rename the queue file (if this is an intentional re-do) or delete it (if it's a stale duplicate), then flip status back to pending` | — |
| 5 | Blocked-by-dependencies section — `⚠ N REQs blocked by unmet dependencies:` | a `pending` REQ has an unmet `depends_on`, or a REQ is `blocked-dependency-cycle` | `(pending; depends on REQ-MMM, status: <status>)`, or `(blocked-dependency-cycle; chain: REQ-PPP → REQ-QQQ → REQ-PPP)`. Pending REQs stay `pending` — the gating is dynamic; only cycle-detected REQs carry a held status | `Resolve the blocking REQs first, then re-run. To force a scoped run that ignores dependency gating, use do-work run REQ-NNN. To break a cycle, edit depends_on and flip the status back to pending. A dependency on a cancelled or failed REQ never self-resolves — re-point it, or abandon the dependent too.` |
| 6 | Assigned-elsewhere section — `⚠ N REQs assigned to another session:` | a `pending` REQ carries a non-empty `assigned_to` | `assigned to <assigned_to, verbatim>` — rendered as written, never normalized | `Skipped as a courtesy, not blocked — nothing confirms that session is running. To take one over here, name it explicitly (do-work run REQ-NNN), which clears the assignment as part of the claim. To drop an assignment without running it, remove the field by hand.` |
| 7 | Skipped-as-negligible section — `⚠ N REQs skipped as impact-negligible:` | `--skip-impact-negligible` is set and an otherwise-claimable `pending` REQ resolves to `impact-negligible` | the resolved token, never the raw one | `Skipped by --skip-impact-negligible, not blocked — nothing was written to these REQs. Re-run without the flag to include them, or name one explicitly, which overrides the skip for that REQ. A REQ with no impact: reads as impact-user-visible and never appears here.` |
| 8 | Held-for-heavy-testing section — `⚠ N REQs held for heavy testing:` | a claimed REQ still carries `## Heavy Verification Plan` without `heavy_verified_revision` after the drain ran (`actions/work.md` Step 7.7) | `held: <finding code> <lane>` — a REQ whose selected lanes all executed was finalized by the drain and never appears here | `A skipped browser lane needs an engine: set QUEUE_KANBAN_BROWSER=<path> and re-run do-work run. Plan drift or a stored historical-revalidation plan is a typed finding for a human; the next drain retries once the cause is fixed.` |
| 9 | Set-aside-by-recovery section — `⚠ N REQs set aside by recovery:` | a `recover` finalization record whose `reason_codes` include `FINALIZATION-SET-ASIDE` | `REQ-NNN (set aside: <reason codes, comma-separated>)`, then an indented `recover: <resolving verb>` | the resolving verb itself, judged as below |

Four rows carry a judgment no typed record makes for you. **Rows 6 and 7 are not held statuses:** those REQs stay `pending`, nothing was written to them, and the same REQ becomes claimable the moment a user names it explicitly — row 7 also renders on the targeted exit path (`actions/work.md` Step 1 → Targeted mode), scoped there to the resolved token set rather than the whole queue. **Row 3a's counts are the one bounded exception to this summary's frontmatter-only stance:** open the body of stakeholder REQs only — at most one per stakeholder, by construction — and fall back to row 3's plain line when a body cannot be read. **Row 9's resolving verb is this summary's judgment, not a field to copy:** pick it exactly as **Stuck Runs Hand Off to Judgment (any step)** picks it, above — `do-work run-with-recovery` when this checkout is the only writer and releaser of the queue, `do-work cleanup` when the archive itself needs repair. That record's `next_argv` is empty on purpose, because the only verb the command could name is the one that just refused, so a missing verb here means the summary skipped the judgment rather than that none exists. The set-aside REQ keeps the claim it already had — recovery does not release it, under `do-work run-with-recovery` either — and its unfinished journal is still on disk.

**After rendering all applicable sections, exit the work loop.** There is no claimed member: Step 1's contract on this path is "render the composed summary, then stop", and the only path that continues is one whose typed queue result contains a committed claim. If **no section applies** (no REQs at all in `do-work/queue/`), report completion and exit. Never silently exit when any section applies — every non-pending or non-ready REQ is something the user needs to see, and a queue hitting several categories renders each section back-to-back instead of one branch's slice.

## In-Progress Record (Step 1)

**Canonical lifecycle transaction boundary.** `do-work-cli recover`, queue-mode `advance`, `advance --checkpoint`, `complete`, `fail`, and `cancel` own the deterministic request, checkpoint, archive, UR, calibration, and provenance changes; `claim`, `unblock`, and `recover-claim` remain compatibility primitives behind those public compositions. The action supplies authority and judgment and then consumes the typed result — it never executes a second move, collision check, UR consolidation, or calibration append, and never reproduces a refused mutation by hand. A refusal is actionable: read it with **Stuck Runs Hand Off to Judgment (any step)**, above.

**One calibration judgment stays here.** The span is read from the just-archived REQ file's own `claimed_at` and `completed_at` at calculation time — never from either stamp held in context earlier in the run.

`do-work/CHECKPOINT.md`'s exact `## In Progress (interrupted)` section is structural claim evidence, and the `writer:` label on an entry is what makes a claim attributable across checkouts. `advance` writes or refreshes the current checkout's entry in the same guarded commit that moves a selected REQ to `working/`; authorized recovery and finalization remove every same-request entry on departure while preserving unrelated labelled or unlabelled records byte-for-byte. Any move of a REQ out of `working/` owes that same removal, and it belongs to the move rather than to a later sweep. The checkpoint grants no lock and no liveness claim.

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

## Pre-Build Evidence Record Template

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

## Testing Section Template

```markdown
## Testing

**Tests run:** [command]
**Result:** ✓ All passing (X tests)

**Repository gate retry:** *(only when a red gate was rerun)* first run exited [X], rerun exited [Y]

**Red-green validation:** *(for bug fixes and new features)*
- [test name/file]: ✗ before implementation → ✓ after
- [test name/file]: ✗ before implementation → ✓ after

**New tests added:**
- [list]

**Existing tests updated (cross-REQ impact):**
- [test file] (from REQ-NNN): [what changed and why — intentional behavior change]

**Heavy verification plan:** *(when lanes were selected)*
- Range: [base revision]..[target revision]
- [lane id]: [exact argv] — [selection reason]

*Verified by work action*
```

## Finding-Closure Ratchet (Step 6.5)

**A review- or triage-finding-origin REQ closes only when its captured GREEN names a regression test/check that fails before the fix and passes after, or when the exact named finding surface is deleted.** Closure evidence must match that named check or deletion surface. A bare patch, unrelated green tests, `tdd: false`, and a high review score are not closure evidence.

## Deferred Prime-Link Path Computation (Step 7.5)

**Path computation rule (for use in Step 8):** the lesson bullet is written to the prime's satellite `lessons-<name>.md`, which sits in the prime's own directory, so the link path is relative to that directory — not the repo root. Count how many directories deep the satellite sits (i.e., the number of path components before the filename). Prepend that many `../` steps to the REQ's repo-root-relative archive path. Examples:
- Satellite at `lessons-auth.md` (0 dirs deep) → `do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`
- Satellite at `src/utils/lessons-auth.md` (2 dirs deep: `src/` and `utils/`) → `../../do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`
- Satellite at `web/src/auth/lessons-auth.md` (3 dirs deep) → `../../../do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`

Step 8 validates the resolved target against the lifecycle plan's exact archive postimage before either the lesson or archive mutation runs. A target that will not exist in the installed package because consumers never receive `do-work/archive/` takes the canonical repository URL instead; any other path absent from the plan is reported, never written.

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

   The builder should classify each discovered task when appending it, using the impact vocabulary (**Request File Schema — Full Frontmatter** above). Every noncritical line ends with the literal suffix `→ report only`:
   ```
   ## Discovered Tasks
   - **impact-critical** SQL injection vulnerability in user search endpoint
   - **impact-user-visible** The retry banner never clears after a successful reconnect → report only
   - **impact-rule-change** Three adapters hand-roll the retry loop the shared client already provides → report only
   - **impact-negligible** Variable naming inconsistency in auth module → report only
   ```

   If the builder did not classify them, classify each with the two questions in `actions/review-work.md` Step 10 — that step is the single home of the rubric, and a second copy here would drift from it.

   **For every noncritical discovery:** Keep it in the current REQ's `## Discovered Tasks` section (or the builder hand-back that Step 8 copies into the archived REQ) and end its line `→ report only`. Do not create, append to, convert, or otherwise mutate a follow-up REQ, sweep, `pending-answers` item, prose backlog, or other deferred-work list. Test-only mechanical hygiene has no carve-out.

   If a maintainer decides one deserves queue work, report the promotion command: invoke `do-work capture` with the complete finding line quoted as the capture source. That explicit capture is new user intent; the builder path creates no placeholder.

   **For `impact-critical` discoveries only:** Run the fold-first scan (`actions/capture-reference.md` → **Fold-First Rule**). Append to an eligible root-cause home or create a follow-up REQ with `status: pending`; critical findings skip confirmation and go straight into the queue. Set `impact: impact-critical` verbatim and mirror `[<impact token>] ` in its title under the REQ Title Convention. `effort_estimate` is the other axis and is judged separately, as size: judge the fix and emit either `effort-mechanical` or `effort-substantive`; when the size is genuinely unclear, put that judgment to the user. Omit the field only when neither judging nor asking was possible, and never copy a default in either direction. Add `- [x] Auto-approved: critical severity (security/data/production risk). → Added to queue immediately.` to Open Questions and report `⚠ CRITICAL discovered: [description] — auto-queued as REQ-NNN` prominently.

## Failure Classification (Step 8)

This classification runs at any generation: `review_generated: true` on the failed REQ does **not** suppress its failure follow-up. The critical-only automatic rule in `actions/review-work.md` Step 10 governs review findings, not failed-work recovery; a failed REQ with no successor would die silently.

Before classifying via the symptom table below, **check for upstream failure**. Cascades from a failed prerequisite often present as plausible-looking `code` or `spec` symptoms in the downstream REQ; misclassifying them sends the builder chasing phantom bugs in the wrong domain.

**Upstream-failure short-circuit:**

Read the frontmatter of every REQ this one depends on — the `addendum_to` parent if set, and every `depends_on` entry, each including the legacy aliases the schema documents. Resolving those ids to exactly one record, and refusing an ambiguous one, is the canonical resolver's job and never a glob spelled out here. If any of them carries `status: failed`, skip the symptom table and short-circuit classification:

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

**Procedure:** run the upstream-failure short-circuit first and take its verdict when it fires, otherwise classify from the symptom table. Either way the chosen `status`, `error`, and `error_type`, and any follow-up REQ, are authored as finalization intent and applied by the canonical lifecycle transaction — never by a hand edit and never by a hand move into `archive/`. Report `[REQ-NNN] failed ([type]): [description]. Follow-up: [REQ-NNN] / None.`, prefixed `(upstream cascade — original failure at REQ-NNN)` when the short-circuit fired.

## Changelog Entry Procedure (Step 9)

This procedure owns release judgment and human-facing payload content. A successful REQ gets one changelog entry unless it is an already-green repair no-op; failed/cancelled work and that no-op carry no release manifest or `release_at`. Match an existing repository convention, or use the house heading `## X.Y.Z — [Short Descriptive Title] (YYYY-MM-DD)` with a unique title that tells a scanning reader what changed.

Choose the version source from affirmative project ownership: prefer the project's own release file, otherwise its release tags, otherwise the changelog counter (starting at `0.1.0`). If project-owned version files disagree, leave them unchanged and use the changelog counter; if the chosen source trails the newest changelog entry, bump from the higher value. Judge the delivered change as breaking, additive, or patch; breaking outranks additive, uncertainty takes the smaller bump, and below `1.0.0` breaking changes bump the minor version.

Include a committed lockfile only when it records this package's own version, and author exact old/new payload bytes that change only that package entry. Classify every release target exactly once as project-owned or a suite-maintainer mirror; do not infer ownership from path spelling or Git tracking, create a missing lockfile, rewrite dependencies, or ask a package manager to produce the payload. The planner discovers the mirror set itself: every tracked `VERSION` file or `**Current version**:` line still carrying the old version, and every tracked changelog byte-identical to the declared changelog preimage, must be a declared target, or the release refuses with `RELEASE-MIRROR-UNDECLARED` naming each path.

The house voice is one or two brief sentences leading with why the delivery matters, followed by specific bullets. Load `crew-members/anti-slop.md` for this human-facing text. The action passes its judged title, prose, ownership, bump, mirrors, and exact payload bytes through the single finalization manifest; deterministic validation and publication stay with the finalizer.

## Commit & Metadata-Commit Procedure (Step 9)

The action judges and authors one strict finalization manifest: the exact request id and path, the expected request and checkpoint digests, the terminal transition and timestamp, the writer, the exact implementation/lifecycle/release allowlist, the commit message, any release payload, and the provenance mode. Use `supplied_commit` with the retained merge hash when an isolated implementation already landed; use `primary_commit`, without an implementation hash, when finalization creates the primary commit from the declared paths. Pass that one manifest to the current `advance` continuation. Finalization remains the sole archive, release, commit, provenance, verification, rollback, and journal authority, and there is no hand-edit, helper, or free-form Git fallback.

Consume ordered `finalizations` **one record at a time** — singular `finalization` is compatibility-only — and judge each record on its own. A record is settled when its terminal phase is `cleanup_complete` with empty `blocked_paths` and `reason_codes`. A record whose `reason_codes` carry `FINALIZATION-SET-ASIDE` is one REQ recovery could not finish, and **that excludes the REQ from this run's selection and nothing else**: its own reason codes say what refused, the remaining records still count as settled, and the run keeps draining the queue (**Composed Exit Summary (Step 1)** → row 9 renders it). A typed refusal is the whole-run stop, and it is what dirt no REQ owns looks like — a dirty index or shared lifecycle, release, or protected paths outside the recovery group refuse without naming a REQ, because no REQ owns that cause: its finding names no REQ, and the verb that resolves it comes from **Stuck Runs Hand Off to Judgment (any step)**, above. On interruption, run public `recover` before any later queue read; it runs finalization discovery first.

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
