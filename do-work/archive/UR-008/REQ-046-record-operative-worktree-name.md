---
id: REQ-046
title: "Record the collision-variant worktree/branch name as the REQ's operative name for cleanup and sweeps"
status: completed
route: B
created_at: 2026-07-29T09:30:45Z
claimed_at: 2026-07-29T13:05:00Z
completed_at: 2026-07-29T13:14:29Z
commit: 6b4dcc1
user_request: UR-008
addendum_to: REQ-038
domain: general
prime_files: []
tdd: false
depends_on: []
related: [REQ-043]
batch: deep-review-followups
write_set: [actions/work.md, actions/work-reference.md]
maintenance: false
---

# Record the Operative Worktree Name After a Collision Variant (REQ-038 Follow-Up)

## What

REQ-038 made a crash-recovered REQ re-dispatch under a fresh unique name variant (`-2`/`-3` or timestamp token), but nothing records that variant as the name later steps operate on — both cleanup sites still *re-derive* the original slug-based name.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** `prime_files` is empty; loaded `crew-members/general.md` + `crew-members/coding-guardrails.md`. Approach: sweep every worktree/branch-naming site and split them into *re-derives from the slug* (the bug surface) vs *discovers by enumerating git* (already correct) — see `## Exploration`. Then (a) state the operative-name rule once, as a new paragraph immediately after the collision (Naming) paragraph in `actions/work-reference.md`, introducing `<operative_name>` as the held literal in the same held-in-memory pattern as `<pre>`/`<merge_hash>`/`session_id` (no new frontmatter field, no new persisted state); (b) repoint every re-deriving command in the two write_set files to `<operative_name>`; (c) narrow the two-option variant scheme to the incrementing numeric token alone, keeping the `worktree-agent-REQ-NNN-` prefix so the sweeps' greps and the two `contract-regressions.sh` asserts stay green.
- [x] **[APPLY]:** Five prose sites edited plus one rationalization row, all inside the declared write_set. No file outside `actions/work.md` + `actions/work-reference.md` touched; `actions/cleanup.md` deliberately left alone (out of write_set, and it needs no change — see `## Exploration`).
- [x] **[UNIFY]:** `git diff --stat` → `actions/work-reference.md | 10 ++++++----`, `actions/work.md | 7 ++++---` (2 files, 10 insertions, 7 deletions). `git status --porcelain` shows exactly those two files modified and nothing else. Verified each: **`actions/work-reference.md`** — the crash-sweep exemption bullet (`:214`), the Naming paragraph's variant narrowing + the new operative-name paragraph (`:229-231`), the hand-back merge argument (`:244`), the Step 8 cleanup commands (`:256`); all internal cross-references resolve to real section headings and their above/below directions are correct relative to their new positions. **`actions/work.md`** — hand-back merge step 2 (`:416`), the hold-as-literals sentence (`:419`), Step 8 substep 8 (`:600`), one new Common Rationalizations row (`:755`). Native check: `bash _dev/tests/contract-regressions.sh` → `Contract regression checks passed.` (both worktree asserts — `worktree-agent-REQ-` present in work-reference.md, `git branch -d` still present — hold). Prose files, so no linter/build applies. No debug artifacts; the only remaining `worktree-agent-REQ-NNN-<suffix>` occurrences are the creation definition and the one gloss stating that absent a collision the operative name *is* that derived name.

## Prior Implementation

REQ-038 (archived, commit `efb6300`, v0.146.1): on a name collision at worktree creation, dispatch under a fresh unique variant keeping the `worktree-agent-REQ-NNN-` prefix; the crash sweep reports (never deletes) the unmerged leftover, which coexists with the variant until cleanup Pass 5.

## Detailed Requirements

1. `actions/work.md` Step 8 cleanup (~:584) and `actions/work-reference.md` (~:243) name the branch/worktree by the *derived* form (`worktree-agent-REQ-NNN-<suffix>`, suffix from the slug). After a variant dispatch, cleanup targets the **unmerged leftover** instead of the variant: `git worktree remove`/`git branch -d` refuses, tripping the "refusal means a merge was skipped or lost — stop and report" branch and halting a successful run on a false alarm, while the variant worktree is never cleaned. Add one rule at the collision paragraph: the name actually created is the REQ's operative name for every later worktree/branch operation (Step 8 cleanup, the crash sweep's bookkeeping, the hand-back report), and cleanup instructions reference the recorded name, not a re-derivation.
2. While there, pick one variant scheme (recommend the incrementing token) instead of offering two with no selection rule — two sessions offered a free choice can diverge between runs.

## Constraints

- The `worktree-agent-REQ-NNN-` prefix must keep matching all variants so sweeps still correlate names to REQs — that invariant is what REQ-038 preserved; don't trade it away.

## Red-Green Proof
**RED prompt/case:** Simulate REQ-038's own scenario one step further: re-dispatch REQ-NNN under `worktree-agent-REQ-NNN-slug-2`, complete the build, run Step 8 cleanup as written — the prescribed commands name `worktree-agent-REQ-NNN-slug` (the leftover), the removal refuses, and the run halts on the "merge was skipped or lost" false alarm.
**Why RED now:** REQ-038 defined how to *create* the variant but not that later steps must *use* it.
**GREEN when:** The operative-name rule is stated at the collision paragraph and both cleanup sites operate on the recorded name; a variant-dispatched REQ completes Step 8 cleanup without touching the leftover (which stays for Pass 5, as designed); one variant scheme is prescribed.
**Validation:** User confirmed (approved capture of the reviewed finding set)

## Full Context
See `do-work/user-requests/UR-008/input.md`.

## Triage

**Route:** B (Explore then Build)
**Reasoning:** Clear outcome (record the operative worktree name; both cleanup sites reference it, not a re-derivation), but the "every later worktree/branch operation" clause requires finding all sites that name a worktree/branch by the derived form (Step 8 cleanup, the crash sweep, the hand-back report, and the collision paragraph itself) before editing — an exploration/scope pass, not a blind single-spot edit. Two files, one small design choice (item 2). Route B.
**Complexity indicators:** 2 requirements; a completeness requirement ("every later worktree/branch operation" — must enumerate the sites); the `worktree-agent-REQ-NNN-` prefix invariant must be preserved. `maintenance: false` per the REQ marker (a fix + a narrowing of a two-option choice, not a delete-oriented instruction-maintenance pass).
**Rigor:** Standard independent review (main-context) against the diff + a completeness check that every worktree-name-operation site now uses the recorded operative name.

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## Exploration

Swept `actions/work.md`, `actions/work-reference.md`, `actions/cleanup.md` (plus `docs/`, `_dev/tests/`) for every site that names a worktree or branch. Two distinct classes emerged: sites that **re-derive** the name from the REQ slug (the bug surface) and sites that **discover** it by enumerating git (already correct — nothing to change).

**Re-derives the name from the slug (must use the recorded operative name):**

| Site | What it does | Status |
| --- | --- | --- |
| `actions/work-reference.md:229` — *Naming* | Derives `worktree-agent-REQ-NNN-<suffix>` at creation; defines the collision variant. **The home of the new rule.** | RE-DERIVES (by design — this is creation) |
| `actions/work-reference.md:242` — hand-back sequence step 2 | `git merge --no-ff --no-commit worktree-agent-REQ-NNN-<suffix>` | RE-DERIVES → **fixed**. Worse than cleanup: after a variant dispatch this merges the *leftover*, not the builder's work. |
| `actions/work-reference.md:254` — *Cleanup — happy path (Step 8)* | `git worktree remove <path>` + `git branch -d worktree-agent-REQ-NNN-<suffix>` | RE-DERIVES → **fixed** (the REQ's named bug) |
| `actions/work.md:416` — Step 6 hand-back merge step 2 (added by REQ-043) | `git merge --no-ff --no-commit worktree-agent-REQ-NNN-<suffix>` | RE-DERIVES → **fixed** |
| `actions/work.md:600` — Step 8 substep 8 | `git branch -d worktree-agent-REQ-NNN-<suffix>` | RE-DERIVES → **fixed** |

**Discovers the name by enumeration (no re-derivation — correct as written):**

| Site | Why it needs no change |
| --- | --- |
| `actions/work-reference.md:212-217` — Crash Recovery *Worktree sweep* | Enumerates `git worktree list --porcelain` + `git branch --list 'worktree-agent-*'` and operates on the names it finds. **One bookkeeping gap fixed:** its exemption bullet (`:214`) gated by REQ **id**, so after a variant dispatch both the variant and its predecessor leftover carry the same id and *both* got exempted — silently un-reporting the leftover the earlier sweep had reported. |
| `actions/cleanup.md:126-128` — Pass 5 | Same enumeration pattern; never re-derives. Its live-claim gate is also id-based, but Pass 5 runs cross-session where another session's operative name is unknowable, and over-exemption there means *skip and report*, never delete — the safe direction. **No change needed, no Discovered Task.** |
| `actions/work-reference.md:231` — *Where worktrees live* | `../<repo>-worktrees/worktree-agent-REQ-NNN-…` is an illustrative location form, not an operation on a specific worktree. |
| `actions/work.md:754` — rationalization row | About `-D` vs `-d`, not about which name to target. |
| `docs/cleanup-guide.md:25,49`, `_dev/tests/contract-regressions.sh:447,462` | Reference the `worktree-agent-*` **prefix** only — preserved by this change, so both stay green. |

## Scope

**Files I will touch:**
- `actions/work-reference.md` — the Naming/collision paragraph (rule's home), the hand-back merge sequence, the Step 8 cleanup paragraph, the Crash Recovery worktree-sweep exemption bullet
- `actions/work.md` — Step 6 hand-back merge step 2 + the hold-as-literals sentence, Step 8 substep 8, one Common Rationalizations row

Exactly the declared `write_set`. `actions/cleanup.md` is out of scope by declaration and needs no change on the merits (it discovers names by enumeration). If anything else had needed editing I would have stopped and reported rather than widening the set.

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified) — narrowed the collision-variant scheme to the incrementing numeric token; added the operative-name rule as a new paragraph right after it; repointed the hand-back merge argument and both Step 8 cleanup commands to `<operative_name>`; tightened the crash sweep's own-session exemption to the recorded name.
- `actions/work.md` (modified) — repointed Step 6's hand-back merge argument and Step 8 substep 8's cleanup commands to `<operative_name>`; added `<operative_name>` to the "hold as re-typed literals" instruction; added one Common Rationalizations row against re-deriving the name from the slug.

**What was done:**

*Requirement 1 — the operative-name rule.* A new paragraph at the collision site states that whatever `git worktree add` succeeded with is the REQ's **operative name**, the single string used by every later worktree/branch operation: the hand-back `git merge` argument, Step 8's `git worktree remove` / `git branch -d`, the crash sweep's own-session bookkeeping, and anything reported to the user. It is held in orchestrator memory and re-typed as a literal into each fresh command — the identical pattern already used for `<pre>`/`<merge_hash>` and the lock's `session_id` — under the name `<operative_name>`. The paragraph also states why nothing persists it: every consumer sits inside the same session's run, and the cross-session consumers (the crash sweep, cleanup Pass 5) discover leftover names by enumerating git rather than re-deriving them. All five re-deriving command sites now name `<operative_name>`; the failure the change closes is spelled out at the rule (after a variant dispatch, the derived string names the *leftover*, so removal refuses and the run halts on a false lost-merge alarm while the variant is never cleaned).

One site beyond the REQ's named list was fixed because it is the same bug, one step earlier and more damaging: the **hand-back merge** in both files named the branch by the derived form, so after a variant dispatch it would have merged the leftover instead of the builder's actual work.

*Requirement 2 — one variant scheme.* "Append an incrementing token (`-2`, `-3`, …) or a short timestamp token" became "append an **incrementing numeric token** — `-2`, then `-3` if that collides too", with a one-clause reason for there being no choice. The `worktree-agent-REQ-NNN-` prefix is untouched on every variant, so the sweeps' `worktree-agent-*` greps, `actions/cleanup.md` Pass 5's scope, `docs/cleanup-guide.md`, and both `_dev/tests/contract-regressions.sh` worktree asserts are unaffected.

*Behavioral unchanged-ness.* The rule states explicitly that with no collision the operative name **is** the derived name, so the common worktree path is identical to before and serial mode — which has no worktree, no merge, and no cleanup substep — is untouched.

**Testing:** `bash _dev/tests/contract-regressions.sh` → passed. Prose-only change with no executable surface; the RED case is a documentation-logic trace (Step 8 as previously written names the leftover after a variant dispatch), and it is now GREEN by construction — no command in either file derives a cleanup or merge target from the slug.

## Decisions

- **D-01 (DECIDE & STATE): The operative name is held in orchestrator memory, not persisted to REQ frontmatter or the lock file.** Every consumer of the name — the hand-back merge, Step 8 cleanup, the run's own reporting — executes inside the session that created the worktree, and that session already holds `<pre>`, `<merge_hash>`, and `session_id` the same way. The consumers that *do* outlive the session (Crash Recovery's sweep, cleanup Pass 5) find leftovers by enumerating `git worktree list` / `git branch --list`, so they never needed the name recorded in the first place. Persisting it would have added a schema field, a writer, a crash-recovery reset rule, and a board-parser mirror for zero new capability.
- **D-02 (DECIDE & STATE): Named the held literal `<operative_name>` and introduced it in `actions/work-reference.md`.** Angle-bracket literal matches the established `<pre>`/`<merge_hash>` convention so the "re-type it, never a shell variable" discipline transfers without restatement; it is two words and greppable. The reference file is the canonical home, and `actions/work.md` uses the token with a one-clause gloss.
- **D-03 (DECIDE & STATE): Fixed the hand-back merge sites too, though the REQ named only the cleanup sites.** Same defect class, strictly worse consequence: a merge of the leftover branch integrates the wrong commits (or answers `Already up to date.` and gets treated as an empty hand-back) instead of merely failing loudly at cleanup. Leaving it would have left the REQ's own "every later worktree/branch operation" clause unsatisfied.
- **D-04 (DECIDE & STATE): Tightened the crash sweep's exemption to the recorded name for this session's own claims only.** The bullet gated by REQ **id**, so after a variant dispatch the variant and its predecessor leftover share an id and both were exempted — silently un-reporting a leftover an earlier sweep had reported. Own-session claims can now distinguish the two; another session's operative name is not knowable from here, so its claim still exempts by id alone (conservative: skip-and-report, never delete).
- **D-05 (DECIDE & STATE): Added one Common Rationalizations row to `actions/work.md`.** Re-deriving the name at cleanup time is exactly what the file said to do until this change, which makes it the shortcut an agent will reach for; the row names the specific failure (false lost-merge halt, variant never cleaned) and the specific step, per the earned-not-mandatory test.

## Review

**Acceptance: Pass — overall ~96%.** Main-context review against the full diff + a completeness sweep of every worktree/branch-naming site.

**Requirements (both met):**
1. New "operative name" rule stated once — held-in-memory `<operative_name>`, same pattern as `<pre>`/`session_id`; EVERY re-deriving site repointed: the hand-back merge (work-reference.md step 2 + work.md Step 6), Step 8 cleanup (both homes), and the crash-sweep own-session exemption (exempt by recorded operative name so a same-id/different-name leftover isn't silently skipped). **Beyond-scope catch:** the hand-back *merge* also re-derived the name — after a variant dispatch that would merge the leftover branch (wrong work), worse than the cleanup bug the REQ named; the sweep found it.
2. Single scheme prescribed (incrementing numeric token `-2`/`-3`); timestamp option removed with rationale; `worktree-agent-REQ-NNN-` prefix invariant preserved.

**Completeness:** remaining `git branch -d`/`git merge` mentions are the discover-by-enumeration crash-sweep path (names found by listing, not re-derived — correct) or the generic "Merge, never rebase" statement whose operative invocation is step 2. No slug re-derivation left in any merge/cleanup context.
**Common case unchanged:** no collision ⇒ operative name = derived name; serial mode has no worktrees. contract-regressions passes; cleanup.md correctly untouched (Pass 5 discovers names by enumerating git).

No Important/Critical findings. No follow-ups.

## Lessons Learned
**What worked:** Splitting naming sites into "re-derives from slug" (the bug surface) vs "discovers by enumerating git" (already correct) made the fix boundary obvious and kept cleanup.md legitimately out of scope.
**Worth knowing:** The dangerous site wasn't the one the REQ named (Step 8 cleanup) but the hand-back merge added by REQ-043 — after a collision variant it would merge the leftover branch. Any REQ-038-style variant work must repoint the *merge* argument, not just cleanup. The crash-sweep can exempt precisely by operative name only for THIS session's own claims; another session's operative name isn't knowable, so its claim exempts by id alone.

## Orientation
Worktree dispatch now threads a single held `<operative_name>` (the name `git worktree add` actually succeeded with — the collision variant where there was one) through every later worktree/branch operation: hand-back merge, Step 8 cleanup, and the own-session crash-sweep exemption. One variant scheme (incrementing token). No collision ⇒ operative name = derived name, so the common path and serial mode are unchanged. Lives in `actions/work-reference.md` (Worktree Dispatch Mode + Crash Recovery sweep) and `actions/work.md` (Step 6 hand-back, Step 8 substep 8). No map change — hardens the REQ-038 subsystem.
