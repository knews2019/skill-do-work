---
id: REQ-290
title: Surface impact in REQ titles and add a run filter that skips negligible work
status: completed
completed_at: 2026-08-19T16:35:00Z
commit: 225e287
claimed_at: 2026-08-19T15:54:19Z
created_at: 2026-08-19T14:33:51Z
user_request: UR-060
domain: general
route: B
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: [REQ-289]
maintenance: false
related: [REQ-289]
batch: impact-effort-split
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-08-19T15:54:19Z
  basis:
    - Route B
    - 4-file write set
    - 3 acceptance criteria
write_set:
- skills/do-work/actions/work.md
- skills/do-work/actions/work-reference.md
- skills/do-work/actions/capture.md
- skills/do-work/actions/capture-reference.md
- skills/do-work/actions/review-work.md
- skills/do-work/docs/capture-guide.md
- skills/do-work-toolbox/actions/code-review.md
- skills/do-work/docs/work-guide.md
---

# Surface Impact in REQ Titles and Add a Run Filter That Skips Negligible Work

## What

Make the `impact:` field REQ-289 introduces actually usable for the decision the user wants to make:
put the token in the REQ title so it is searchable today, and give `do-work run` a flag that skips
negligible-impact work.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Approach below.
- [x] **[APPLY]:** Applied as planned across the 7 write-set files; nothing outside it written.
- [x] **[UNIFY]:** `git diff --stat` reviewed, then every hunk of the 7-file diff read back (`git diff -U0 skills/`). No debug artifacts, no stray scratch files — the one throwaway Go probe was deleted and its absence confirmed. `_dev/tests/contract-regressions.sh` and `_dev/tests/maintainer-verify.sh` both exit 0, unpiped.

### Plan

**Part 1 — title convention.** New `## REQ Title Convention` section in `capture-reference.md`, placed
immediately after the file's intro and before `## Request File Formats` so every template below it
inherits it. Shape: a bracketed classification tag ahead of any kind prefix —
`title: "[impact-negligible] Review fix: …"` (D-02). Tag emitted only when `impact:` is something
other than the `impact-user-visible` default, matching the board's chip rule (D-03). Tag lives in
`title:` only — never the filename slug, never the body H1 (D-04). Every one of the five emitters
(`capture-reference.md` Simple + Addendum, `review-work.md:362`, `work-reference.md:629`,
`code-review.md:299`) gets the tag placeholder and a pointer; `docs/capture-guide.md` gets the
user-facing mirror. `work-reference.md:137`'s schema line states the source-of-truth rule once.

**Part 2 — the run filter.** `--skip-impact-negligible`, one boolean flag:
1. `work.md` `## Input` — new bullet after `--wave N`, stating the compose/override contract.
2. `work.md:107` strip list + `:110` usage string (both branches).
3. `work.md:157` queue status summary — conditional `(N skipped as impact-negligible)` suffix, which
   is the reporting path for a run that *does* find work.
4. `work.md:193` — new paragraph immediately after the `assigned_to` skip, same shape.
5. `work.md:195` — add the skip to the auto-wave "every filter above still applies" list.
6. `work.md:199` — fourth headline + `skipped-as-negligible` in the section order.
7. `work.md:664` — Step 0 flag list.
8. `work-reference.md` auto-wave ready set — fifth condition, plus the "four conditions" /
   "conditions 2 and 4" counts in the paragraph below it.
9. `work-reference.md` Composed Exit Summary — new section 7 + the fourth headline in its lead.
10. `work-reference.md:237` Schema Read Contract row — name Step 1's filter as a reader, so the
    normalization (absent ⇒ `impact-user-visible` ⇒ never skipped) is attached to the filter.

**Verification:** `_dev/tests/maintainer-verify.sh` and `_dev/tests/contract-regressions.sh`, both
unpiped, exit 0. The bare axis words `user-visible`/`rule-change` must never appear un-prefixed
(`contract-regressions.sh` REQ-289 check).

## Why

A field nobody can filter on does not let the user stop work. Today `do-work run` has no size or
impact filter of any kind — `work.md:103-113` accepts only targeting tokens, `--fan-out`, and
`--wave`, and Step 1's selection scan reads status, dependency readiness, claim state, `assigned_to`,
and wave depth. Nothing else.

## Context

The board's search box already matches on `request.title` (`web/board-filters.js:20-31`), so a token
in the title is filterable with zero board work. A dedicated board filter control would need Go, JS,
and CSS changes; the title gets there for free.

There is no REQ title format rule anywhere in the repo today. `work-reference.md:119` and
`capture-reference.md:14` carry only the placeholder "Brief descriptive title". Prefix conventions
already exist in practice but are unenforced: `capture-reference.md:166` ("Addendum: ...") and
`review-work.md:362` ("Review fix: ...").

## Detailed Requirements

### The token goes in the title, not the filename

- The `impact:` token appears in the REQ's `title:` frontmatter value.
- It does **not** go into the filename slug. `REQ-NNN-slug.md` stays as it is. Impact can be revised
  at clarify, and a filename-borne token would mean renaming files mid-pipeline for a metadata edit.
  The board searches titles, not filenames, so nothing is lost.
- Give the convention a home in `capture-reference.md` alongside the templates, and follow the
  existing prefix precedents rather than inventing a third shape.
- `capture.md` and `review-work.md` emit titles conforming to it.
- `impact:` frontmatter stays the source of truth. The title is a mirror; when they disagree, the
  field wins.

### The run filter

- `do-work run --skip-impact-negligible` omits `impact: impact-negligible` REQs from Step 1's
  selection.
- One boolean flag, not a general `--impact <token>` selector. YAGNI until a second use appears.
- `work.md:107` rejects unrecognized arguments, so the usage string at `:110` changes in the same
  edit.
- The flag changes *which* REQs are selected. State how it composes with `--wave` and with explicit
  `REQ-NNN` targeting: explicit targeting names a REQ deliberately and must not be silently skipped.
- Report what was skipped. A run that silently drops REQs reads as "the queue is empty" when it is
  not.

## Constraints

- Blocked on REQ-289 — there is no field to filter on until it lands. `depends_on` enforces this.
- Write-set overlap with REQ-289 on `work.md`, `capture.md`, `capture-reference.md`, and
  `review-work.md` is expected and safe because the two are strictly sequential. Do not run them in
  parallel.
- No new board code. If a board filter control is wanted later it is a separate REQ.

## Dependencies

REQ-289 must complete first.

## Red-Green Proof

**RED prompt/case:** With a queue holding both negligible-impact and user-visible REQs, try to run
only the ones worth doing, and try to find the negligible ones on the board.

**Why RED now:** `do-work run` has no flag for it — `work.md:110`'s usage string offers targeting
tokens, `--fan-out`, and `--wave`, and unrecognized arguments are rejected. The board's search box
matches titles, and no title carries the impact token, so searching finds nothing.

**GREEN when:** `do-work run --skip-impact-negligible` runs the queue and omits exactly the
`impact-negligible` REQs, reporting how many it skipped. Typing `impact-negligible` into the board's
search box lists exactly those REQs, with no board code changed.

**Validation:** User confirmed — "Field + title prefix + run filter" chosen via the ask tool, with
the stated purpose "if I ever want to stop spawning and processing them I can".

## Full Context

See `do-work/user-requests/UR-060/input.md` for complete verbatim input.

---

## Triage

**Route: B** - Medium

**Reasoning:** Four action files, and the REQ cites the exact argument-parsing sites, so the "what" and most of the "where" are settled. Exploration is still worth one pass: the title convention has to follow the two existing prefix precedents rather than invent a third shape, and the flag's composition with `--wave` and with explicit `REQ-NNN` targeting has to be read off Step 1's actual selection scan rather than assumed.

**Planning:** Not required

## Exploration

**Title emitters — the REQ's premise was short by three.** It names `capture.md` and `review-work.md`
as the emitters. The actual set that templates a `title:` value:

| Site | Template |
|---|---|
| `capture-reference.md:14` | Simple REQ — `title: Brief descriptive title` |
| `capture-reference.md:167` | Addendum REQ — `title: "Addendum: dark mode sidebar support"` |
| `review-work.md:362` | Review follow-up / sweep — `title: "Review fix: [brief description]"` |
| `work-reference.md:629` | **Builder-Decided Follow-up (Step 8)** — `title: "Confirm: [brief description of the choice]"` |
| `do-work-toolbox/actions/code-review.md:299` | `title: "Code review: [brief description]"` |

`work-reference.md:119` carries a **`Short` descriptive title** placeholder, not "Brief" — a grep for
the phrase the REQ quotes misses it. `docs/capture-guide.md:34` is the user-facing mirror.
Discovered-Tasks and Failure classification create follow-ups in prose and emit no `title:` at all.

**The existing prefix shape.** All four precedents agree: leading `<Prefix>: `, sentence-case prefix
(`Review fix`, not `Review Fix`), the whole YAML value **double-quoted**, and the body H1 repeating
the prefix in Title Case. The quoting is load-bearing, not style —
`queue-kanban/frontmatter.go:98-106` documents that an unquoted colon-bearing title drops the whole
frontmatter block into lenient line-based recovery, and `frontmatter_test.go:172,179` pins that.

**The prefix slot is already occupied.** A review-fix REQ that is also negligible needs both, and
`"impact-negligible: Review fix: …"` is a double-colon title. This is the REQ's real design question.

**Where the filter slots into Step 1.** `work.md:193`'s `assigned_to` skip is the exact model — the
only existing flag-independent-value skip, reported in the exit summary and overridden by explicit
naming. The new filter belongs immediately after it, and must ALSO be added at:
- `work.md:195` — the auto-wave restatement of the filter stack.
- `work-reference.md:374-383` — the auto-wave ready-set conditions, which close with **"Nothing else
  enters the computation."** Miss this and `--skip-impact-negligible` silently does nothing under
  `--fan-out`.
- `work.md:664` — Step 0 of the Orchestrator Checklist, which lists the flags.
- The Composed Exit Summary (`work-reference.md:407-479`, mirrored at `work.md:199`), whose rule at
  `:411` is explicitly extensible: "that condition is the rule, and the list below is the set as it
  stands today". Section 6 (assigned-elsewhere, `:467-475`) is the closest render template.

**Composition precedent.** `--wave` is *rejected* alongside targeting tokens (parse-time,
`work.md:105,151`); `depends_on` and `assigned_to` are *overridden* by them (`work.md:155,193`). The
REQ wants override semantics, which matches the `assigned_to` precedent, not the `--wave` one.

**What could break.** `contract-regressions.sh:3311-3316` slices `work.md`'s `## Input` … `## Steps`
block and asserts it contains `UR-NNN` — a `contains`, so adding a flag is safe, but restructuring
the section boundary is not. REQ-289's own check at `contract-regressions.sh:1736-1742` **fails if
the bare axis words `user-visible`/`rule-change` appear without the `impact-` prefix**, so the title
convention must write full tokens, never a shortened form.

**Board search confirmed.** `web/board-filters.js:20-31` matches case-insensitive substring on
`requestId`, `request.title`, and `request.userRequestId`. A title-borne token is filterable with no
board change, exactly as the REQ claims.

*Generated by Explore agent; findings spot-checked by the orchestrator.*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/capture.md`
- `skills/do-work/actions/capture-reference.md`
- `skills/do-work/actions/review-work.md`
- `skills/do-work/docs/capture-guide.md`
- `skills/do-work-toolbox/actions/code-review.md`
- `skills/do-work/docs/work-guide.md`

**Acceptance criteria (restated from the REQ's Red-Green Proof):**
1. `do-work run --skip-impact-negligible` runs the queue and omits exactly the `impact-negligible`
   REQs, reporting how many it skipped.
2. The flag composes with `--wave` and `--fan-out`, and an explicitly-named `REQ-NNN` is never
   silently skipped by it.
3. Typing `impact-negligible` into the board's search box lists exactly those REQs, with no board
   code changed.
4. The impact token appears in `title:`, never in the filename slug; `impact:` frontmatter stays the
   source of truth and wins when the two disagree.
5. Every REQ-title emitter conforms to the convention.
6. `bash _dev/tests/maintainer-verify.sh` exits 0.

## Decisions

- **D-01**: Extend the write set from 4 files to 7. Reasoning: exploration found three more sites
  the REQ's premise missed. `work-reference.md` is **required** on two independent counts — it
  templates a REQ title at `:629`, and its auto-wave ready-set block at `:374-383` closes with
  "Nothing else enters the computation", so a filter not added there silently no-ops under
  `--fan-out`, which is the REQ's own composition requirement. `code-review.md` and
  `capture-guide.md` are the remaining emitter and its user-facing mirror; leaving them ships a
  convention its own emitters do not follow — the same half-applied-convention defect REQ-289's
  review caught as F7. ESCALATE — widening the write set is a scope judgment.
  **Value:** the flag actually works under `--fan-out`, and the convention holds everywhere a title
  is minted.
  **Risk:** low and reversible; REQ-290 is the last REQ in this UR's original batch and nothing runs
  in parallel with it.

- **D-02**: The impact token rides in the title as a **bracketed classification tag** ahead of any kind
  prefix — `title: "[impact-negligible] Review fix: …"` — not as a fifth `<Prefix>: ` kind prefix.
  Reasoning: the `Prefix: ` slot is already occupied by four kind prefixes (`Addendum`, `Review fix`,
  `Confirm`, `Code review`), and reusing it forces either a double-colon title or an ordering rule
  nobody would remember. The REQ's "don't invent a third shape" constraint is about not inventing a
  third *kind-prefix* form; a bracketed tag is a different slot — a classification tag, not a kind —
  so the two compose instead of competing. Verified mechanically: strict YAML parses
  `title: "[impact-negligible] Review fix: guard misses hex shorthand"` as a plain string and the
  board's `parseFrontmatterFields` returns `status` and `depends_on` intact, so the bracket never
  reaches the lenient recovery path. ESCALATE — a title format is user-facing and reasonable people
  would disagree about the shape.
  **Value:** one convention that never collides with the four existing prefixes, and a token the
  board's title-matching search box finds verbatim with no board code changed.
  **Risk:** low and reversible — the tag is a string in one field, not a parsed structure, and nothing
  reads it. Titles already written are unaffected; a later reshape is a find-and-replace. The one
  real cost is the leading bracket eating card width on the board.

- **D-03**: Emit the tag **only when `impact:` is something other than the `impact-user-visible`
  default**, mirroring the board's impact chip rule exactly. Reasoning: tagging every REQ would put
  `[impact-user-visible]` on the large majority of titles for zero filtering benefit, and would make
  title and card disagree about which REQs are worth marking. Absent `impact:` reads as the default,
  so legacy REQs need no retrofit. ESCALATE — same user-facing-format reason as D-02.
  **Value:** the common case stays unadorned, and the title agrees with the chip a reader is already
  looking at.
  **Risk:** a title alone cannot distinguish "judged user-visible" from "never judged" — the same
  ambiguity the chip already has, and `impact:` is the source of truth for anything that acts on it.

- **D-04**: The tag lives in `title:` only — never the filename slug (the REQ mandates this), and
  **never the body `# H1`** (my extension of the same rule). The H1 is Title-Cased prose, and
  `[Impact-Negligible]` would break the exact string the tag exists to be searched for. DECIDE &
  STATE — the H1 clause follows directly from the token-must-stay-verbatim requirement.

- **D-05**: `skills/do-work-toolbox/actions/code-review.md` conforms to the convention by carrying
  **no** tag, because that action writes no `impact:` field at all, so its REQs read as the default.
  Its title template is left as-is and a one-line pointer states why. Adding an `impact:` verdict to
  that emitter is REQ-289-shaped work, not this REQ's — filed under Discovered Tasks. DECIDE & STATE.

- **D-06**: Added the missing `impact:` line to `skills/do-work/docs/capture-guide.md`'s REQ schema
  block. Not a scope grab: the paragraph this REQ requires there describes a tag mirroring a field
  the guide had never shown, so the prose would have referenced a field the reader has never seen.
  DECIDE & STATE.

- **D-07**: Added a **fourth** composed-exit-summary headline — `No claimable pending REQs — every
  ready one is impact-negligible and --skip-impact-negligible is set.` Reasoning: none of the three
  existing headlines fits a queue whose pending REQs exist and are dependency-ready but were all
  dropped by the flag; the closest (`No pending REQs in queue.`) is the exact "reads as the queue is
  empty when it is not" failure the REQ names. DECIDE & STATE — the headline set is already keyed to
  queue state, so this is filling in a state, not adding a mechanism.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/capture-reference.md` (modified) — new `## REQ Title Convention` section as the canonical home; the Simple REQ and Addendum templates carry the quoting and tag comment.
- `skills/do-work/actions/capture.md` (modified) — Step 1's impact assessment mirrors a non-default verdict into the title; a Verification Checklist line covers the tag and its absence from filenames.
- `skills/do-work/actions/review-work.md` (modified) — Step 10's follow-up template title carries the tag ahead of `Review fix: `; the prose above states the emit-only-when-non-default rule.
- `skills/do-work/actions/work-reference.md` (modified) — the `impact:` schema line gains the source-of-truth-vs-mirror rule and names the new filter; the Schema Read Contract row lists the filter as a reader; the auto-wave ready set gains condition 5 with its counts corrected; the Composed Exit Summary gains a fourth headline and a seventh section; the Builder-Decided Follow-up template title is tagged.
- `skills/do-work/actions/work.md` (modified) — `--skip-impact-negligible` added to `## Input`, the argument-strip list, both usage-string branches, the queue status summary suffix, a new Step 1 skip paragraph after the `assigned_to` skip, the auto-wave filter list, the exit-path headline and section order, and Step 0 of the Orchestrator Checklist.
- `skills/do-work/docs/capture-guide.md` (modified) — user-facing mirror: quoted title, `impact:` added to the schema block, one paragraph on the tag and the flag.
- `skills/do-work/docs/work-guide.md` (modified) — the flag's user-facing documentation, placed beside the existing `--wave` / `--fan-out` prose (D-08, orchestrator).
- `skills/do-work-toolbox/actions/code-review.md` (modified) — a one-line pointer to the convention; its titles stay untagged because the action writes no `impact:` (D-05).

**What was done:** Made REQ-289's `impact:` field usable for the decision the user actually wants to make. The token now appears as a bracketed classification tag at the front of a REQ's `title:` — `"[impact-negligible] Review fix: …"` — which composes with the four existing kind-prefixes instead of competing for the same slot, and is found verbatim by the board's existing title search with no board code changed. `do-work run --skip-impact-negligible` omits negligible REQs from selection, reports how many it passed over, stacks with `--wave`, composes with `--fan-out`, and is overridden by explicit `REQ-NNN` naming on the same provenance rule that governs `depends_on` and `assigned_to`. The filter was added to the auto-wave ready-set conditions as well as to Step 1, so it does not silently no-op under `--fan-out`.

## Qualification

**Passed — 8 project files verified, 6 acceptance criteria traced, P-A-U confirmed.**

Mechanical (`tools/checks/qualify.sh`): OK. Scope drift (`tools/checks/scope-drift.sh`): "Implementation Summary matches the Scope declaration" — no undeclared touch, no unused declaration.

Judgment checks, run against git state rather than the builder's report:

- **Substantive (check 2).** All eight files carry real edits. The two that decide whether the REQ works were read line by line: `work.md`'s new Step 1 skip paragraph, and `work-reference.md`'s auto-wave condition 5 — which is the one that stops the flag no-opping under `--fan-out`. Its neighbouring counts were corrected with it ("four conditions" → "five", "carve-outs on conditions 2 and 4" → "2, 4, and 5"), which is the tell that the edit was made rather than pasted.
- **Requirements traced (check 3).** All six acceptance criteria map to a diff site; the composition rule stated in `work.md` and the auto-wave conditions in `work-reference.md` agree with each other.
- **Flowing (check 6).** The conservative-resolution property is the one that matters here — an absent or unrecognized `impact:` must read as `impact-user-visible` so an unjudged REQ is never dropped. Verified live in REQ-289's fixture parse and restated at both the schema line and the Step 1 paragraph.

**One finding, corrected rather than passed downstream (D-09).** The convention's always-double-quote rule carried a justification that does not survive checking: it claimed an unquoted tagged title leaves "the REQ's status, UR pointer, and dependencies" riding on strict parsing. `frontmatter.go:98-105` says the opposite — strict rejection *would* drop those fields, and the line-based recovery exists so it does not. Tested directly: a REQ with an unquoted `title: [impact-negligible] Review fix: …` parses to identical `status`, `impact`, `depends_on`, and `user_request` values as its quoted twin, with no data warning raised. The rule is right and unchanged; the reason was rewritten to the true one.

## Testing

**Test commands run:**

- `bash _dev/tests/maintainer-verify.sh` — **exit 0**, unpiped, after the final edit.
- `bash _dev/tests/contract-regressions.sh` — exit 0.
- `tools/checks/scope-drift.sh` — clean.

**Baseline comparison:** Step 5.75 recorded a green baseline, so every pass above is a true pass rather than a pre-existing failure. No new regressions.

**Red-green validation.** `tdd: false`, and this REQ changes instructions rather than executable behavior, so regression evidence stands in for a red/green pair. Traced to the REQ's `## Red-Green Proof`:

- *RED as captured* — `do-work run` had no impact filter of any kind, and no REQ title carried the token, so the board's title search found nothing.
- *GREEN as captured* — the flag exists in `## Input`, both usage-string branches, Step 1's scan, the auto-wave ready set, and the Orchestrator Checklist; it reports what it dropped by two paths (a `(K skipped as impact-negligible)` suffix on the queue status line for runs that *do* find work, and a seventh composed-exit-summary section for runs that do not). The title convention has a canonical home and every emitter conforms.

**Live parse evidence** (`queue-kanban` against purpose-built fixtures): a bracketed tagged title round-trips verbatim through the board's reader — `'[impact-negligible] Review fix: guard misses hex shorthand'` — with `status`, `impact`, `depends_on`, and `user_request` all intact, which is what makes the board's existing title search find the token with no board code changed (acceptance criterion 3).

**Acceptance criteria, each with its evidence:**

1. Flag omits negligible REQs and reports the count — `work.md` Input, Step 1 paragraph, queue-summary suffix, exit-summary section 7.
2. Composes with `--wave`/`--fan-out`; explicit `REQ-NNN` overrides — stated in `work.md` and matched by auto-wave condition 5.
3. Board search finds the token, no board code changed — proven on a live parse; `git status` shows no board-tool file touched.
4. Token in `title:` only, field is source of truth — convention rules 3 and 6, plus the schema line.
5. Every emitter conforms — all five named in the convention's canonical-home sentence.
6. `maintainer-verify.sh` exits 0 — confirmed.

## Review

**Overall: 70%** | 2026-08-19T16:35:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 80% |
| Code Quality | 78% |
| Test Adequacy | 65% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token):**
- F1 `work-reference.md`'s "What fan-out adds" bullet still said auto-wave computes its set from `depends_on`, claim state and `assigned_to` **alone** — contradicting new condition 5 thirteen lines below it — impact-user-visible → **fixed in place**
- F2 `work.md:35` carried the same stale four-condition gloss — impact-user-visible → **fixed in place**
- F3 The canonical Request File Schema title line (`work-reference.md:119`) was left unquoted, untagged, and with no pointer to the new convention, contradicting it on both counts — impact-rule-change → **fixed in place**
- F4 Discovered Tasks Classification mints REQs and is required to set `impact:` on every one, but was omitted from the convention's emitter list, so a builder-discovered negligible REQ would carry the field and no tag — impact-user-visible → **fixed in place**, and the emitter list re-keyed on the condition rather than the enumeration
- F5 Targeted mode reaches neither reporting path cleanly, and the skipped count's scope is undefined — impact-user-visible → **REQ-297 created**

**Minor findings:** 4 — three **fixed in place** (the `work-guide.md` OPTIONAL overstatement plus its missing never-skips-an-unjudged-REQ reassurance; the convention's "revised at clarify" justification, which named a flow `clarify.md` does not contain; the stale Discovered Task that D-08 had already closed). The fourth — no lock-in check pins the flag's six declaration sites — was **appended as instance F6 to the existing REQ-293 sweep**, whose root cause it shares.

**Nit:** the reviewer confirmed D-09's correction is accurate by reproducing it against the real parser, and found the true reason is *stronger* than what I wrote: recovery parses a value opening with `[` and closing with `]` as a YAML flow list. I verified it — `title: [impact-negligible] Retitle export, again [v2]` reads back as `[impact-negligible] Retitle export again [v2]`, comma silently eaten, no warning. That is real corruption of exactly the titles this convention mints, so the justification now carries it.

**Acceptance:** Partial as reviewed — held there by F1/F2, which told a reader the wave computes from four inputs "alone" and so could have no-opped the flag under `--fan-out`. All four blocking findings were remediated in place and the gate re-run; the board half was already proven end-to-end.

**Post-remediation evidence:** `bash _dev/tests/maintainer-verify.sh` exit 0, unpiped. Live acceptance test against the real repo queue: typing `impact-negligible` into the board's search box matches exactly REQ-296, whose `impact:` field independently reads `impact-negligible` — title and field agree, and no board code was changed.

**Restatement sweep:** performed repo-wide across flag enumerations (`SKILL.md`, `help.md`, `README.md`, `docs/`, `decisions/`, `_dev/tests/`), REQ-title templates, and Step 1 filter lists. `SKILL.md` and `help.md` enumerate no flags, so nothing drifted there. `work.md`'s exit-summary section list carries an explicit "not a second copy to keep in sync" escape hatch and is correctly not a finding. `adr-018` restates the four ready-set inputs and is correctly left alone — an ADR is a point-in-time record.

**Suggested testing:** 3 items (carried into REQ-297 and REQ-293's F6)
**Follow-ups created:** REQ-297; **sweeps appended to:** REQ-293

*Reviewed by review-work action; F1-F4 and three Minors remediated in place by the orchestrator.*

## Lessons Learned

**What worked:**
- Directing the builder to a specific structural site (`work-reference.md`'s auto-wave conditions,
  which close with "Nothing else enters the computation") rather than to the feature in the
  abstract. That block is where the flag would otherwise have silently no-opped, and naming it in
  the dispatch is why condition 5 exists at all.
- Dogfooding the vocabulary immediately. Writing this UR's own follow-ups with `impact:` tokens and
  then searching the live board for `impact-negligible` is what turned acceptance criterion 3 from
  an assertion into an observation.

**What didn't:**
- **Exploration dismissed the Discovered Tasks flow as emitting "no `title:` at all".** True of its
  template and false of its behavior — the flow mints REQs and is explicitly required to stamp
  `impact:` on every one. Reading a template for what it *contains* rather than what the surrounding
  prose *instructs* is how an emitter goes missing from an emitter list. The fix re-keys the list on
  the condition ("any flow that mints a REQ carrying an `impact:` value") so the next one inherits it.
- **A REQ that adds a condition to a list must sweep every gloss of that list, not just the list.**
  Three restatements of the ready-set conditions survived inside the two files this REQ was already
  editing, one of them thirteen lines from the condition it contradicted. The canonical list was
  updated correctly; the prose *about* it was not.

**Worth knowing:**
- A YAML title that opens with `[` and closes with `]` is parsed as a flow list by the board's
  lenient recovery path and comes back **altered** — commas inside it are eaten as separators, with
  no warning. Quoting is what takes the class off the table. This was found only because a reviewer
  checked a justification that was already "close enough", which is the argument for checking
  reasoning rather than verdicts.
- `--skip-impact-negligible` is deliberately conservative: absent and unrecognized both resolve to
  `impact-user-visible`, so an unjudged REQ is never dropped. That property is what makes the flag
  safe to add to a queue whose REQs mostly predate the field, and it is the one thing no test pins
  (tracked as REQ-293's F4).

## Orientation

Now you can stop building work nobody would notice. `do-work run --skip-impact-negligible` drops
`impact-negligible` REQs from selection and reports what it passed over, and the impact token rides
at the front of a REQ's title so the board's existing search box finds them — no board code changed.
Both live in the queue-selection subsystem (`skills/do-work/actions/work.md` Step 1 and its reference)
plus the capture templates that mint titles.

**[MAP CHANGED]** — `impact:` gains its first *reader that acts on it*. Until now it was an authored
display field; it is now a selection input, which means the auto-wave ready set has a fifth condition
and every gloss of that set is a thing that can go stale. The title convention is also new shared
shape: `capture-reference.md` → **REQ Title Convention** is canonical, and the rule is keyed on the
condition — any flow that mints a REQ carrying an `impact:` value follows it — rather than on a list
of emitters that would drift.

**Prime staleness spot-check:** `_dev/primes/prime-action-files.md` re-read; every path it cites
still resolves and nothing in it is falsified by this change. Its own lesson about closing a site
class by its condition rather than by a line list is what the emitter-list fix applied.

## Discovered Tasks

- **impact-user-visible** — ~~`skills/do-work/docs/work-guide.md` does not mention `--skip-impact-negligible`.~~ **Closed during this REQ by D-08** — the orchestrator extended the write set and wrote the paragraph. Left here with its resolution rather than deleted, so the builder's stop-and-report is still legible in the trail.

- **impact-rule-change** — `skills/do-work-toolbox/actions/code-review.md` creates REQs with no
  `impact:` field, so every code-review follow-up reads as `impact-user-visible` and is invisible to
  `--skip-impact-negligible`. It already judges severity per finding; wiring that to an impact
  verdict would make the third REQ-minting action honor the field REQ-289 introduced.
- **impact-negligible** — nothing retrofits the tag onto REQs already in `do-work/queue/` carrying a
  non-default `impact:`. The convention is emit-time only, so today's queue is searchable by field
  (`grep`) but not by the board's search box until each REQ is retitled by hand.
- **D-08** (orchestrator): Extended the write set to 8 files, adding
  `skills/do-work/docs/work-guide.md`, and wrote the flag's user-facing paragraph there. The builder
  reported this as the one thing it could not do — the file was outside its boundary, so it stopped
  and reported rather than writing silently, which is the correct behavior. Reasoning: the guide
  documents `--wave` and `--fan-out` in detail at two places, so a user-facing flag missing from it
  is undiscoverable by the people it exists for, and "make impact actionable" is this REQ's stated
  What. Shipping the flag without its guide entry would be an incomplete deliverable, not a smaller
  one. DECIDE & STATE — one paragraph, additive, in a file nothing else in this UR touches.
- **D-09** (orchestrator): Corrected the stated reason behind the always-double-quote rule in
  `capture-reference.md`'s title convention. The rule is right and stays; its justification was not.
  The builder wrote that an unquoted title "drops the entire frontmatter block into the board
  parser's lenient line-based recovery — the REQ's status, UR pointer, and dependencies all ride on
  that block parsing strictly," which reads the parser's own counterfactual as the present
  consequence. `frontmatter.go:98-105` says strict YAML rejection *would* drop those fields and that
  the line-based extraction exists so it does not. I tested it: a REQ with an unquoted
  `title: [impact-negligible] Review fix: …` parses to byte-identical `status`, `impact`,
  `depends_on`, and `user_request` values as its quoted twin, with no data warning. Nothing is lost
  today. The rewritten reason is the true one — an unquoted tagged title is served by a last-resort
  salvage path rather than the parser proper, which is a real cost and not a data-loss one.
  DECIDE & STATE — one sentence, in a declared file, rule unchanged.
  Recorded because this repo has shipped correct verdicts backed by false arguments before, and a
  reason that does not survive being checked is the defect regardless of whether the rule is right.
