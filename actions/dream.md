# Dream Action

> **Part of the do-work skill.** A manual, four-phase consolidation pass over a plain-text memory directory — find rot, heal contradictions, merge near-duplicates, rebuild the index. Destructive by design, so it never runs automatically. User-facing walkthrough: [`docs/dream-guide.md`](../docs/dream-guide.md).

Operates on any plain-markdown memory store with an index file (`MEMORY.md`, `_master_index.md`, or `index.md`), a `wiki/` page directory, and an append-only `log.md` timeline. Reproduces the "Auto Dream" Orient → Gather → Consolidate → Prune pattern, but only on explicit invocation.

The four phases run in order. Phases 1–2 are strictly read-only reconnaissance; Phase 3 is the only phase that writes; Phase 4 rebuilds the index and closes the pass. A `.lock` file in the memory directory bounds the pass so two runs never overlap.

## When to Use

**Use when:**
- Memory has drifted: contradictions, stale relative dates ("yesterday"), broken `[[wiki-links]]`, near-duplicate pages, or an index bloated past budget.
- After a burst of capture/ingest work — best run between bursts, not mid-task.
- User says "dream", "consolidate my memory", "clean up the wiki", "lint and merge my notes", "memory cleanup".
- The target is a bkb wiki (`kb/wiki/`) that has accumulated multiple `bkb lint` warnings over several sessions and the user wants a single aggressive pass.

**Do NOT use when:**
- No memory directory exists — dream consolidates existing memory, never initializes it. Redirect to `do-work bkb init` if a knowledge base is wanted.
- The user only wants a read-only health report on a bkb wiki — that's `do-work bkb lint`. Reach for dream when bkb lint has been flagging the same issues for a while and the user is ready to spend the destructive pass.
- Routine relationship hygiene on a bkb wiki — that's `do-work bkb garden` (conservative) or `do-work bkb defrag` (structural).
- Asked to schedule dream on a timer or hook — refuse and explain that consolidation is destructive and must stay explicit.

## Input

`$ARGUMENTS` accepts (in any order):

- `<path>` — absolute or relative path to the memory directory. If empty, resolve via the default path resolution order in Step 1.
- `--dry-run` (mode token) — run Phases 1–2 to build the worklist, print the Phase 2.5 preview, then **stop without writing**. The lock is released cleanly. Mirrors `stray-check`'s `report` mode. A `--dry-run` invocation skips the confirmation prompt and the answer is implicitly "none."

## Steps

### Step 1: Resolve target directory and acquire lock

1. **Peel out the mode token first.** If `$ARGUMENTS` contains `--dry-run` as a standalone token, set `mode = dry-run` and strip it from `$ARGUMENTS`. The remaining text (if any) is the path. If `--dry-run` is absent, `mode = interactive`.

2. **Resolve the path.** If the cleaned `$ARGUMENTS` is non-empty, treat it as the target. Otherwise, try these in order — first hit wins:
   - `./memory/`
   - `./wiki/`
   - `./kb/wiki/`
   - `./knowledge-base/wiki/`

   If none of those exist, stop and report:
   ```
   No memory directory found.

   I checked: ./memory, ./wiki, ./kb/wiki, ./knowledge-base/wiki.

   Dream consolidates an existing plain-text memory store — it does not create one.
   If you want to start a knowledge base, run: do-work bkb init
   If your memory lives elsewhere, run: do-work dream <path>
   ```
   Do **not** create the directory yourself.

3. **Find the index.** Inside the resolved directory, try `MEMORY.md` → `_master_index.md` → `index.md`. The first that exists is `<index>`. If none exists, stop with: `Target directory has no index file (MEMORY.md / _master_index.md / index.md). Refusing to dream a memory store without an index.`

4. **Find the wiki dir.** Try `<dir>/wiki/` → `<dir>/pages/` → `<dir>` itself — pick the first that contains `*.md` files (excluding the index). Call this `<wiki>`.

5. **Acquire the lock.** Check `<dir>/.lock`:
   - If absent: create it with the current UTC timestamp inside, then continue.
   - If present and its mtime is within the last 5 minutes: stop with `Another dream pass may be running (.lock is recent). If you're sure it isn't, delete <dir>/.lock and re-run.`
   - If present and older than 5 minutes: treat as stale — note it in the eventual summary, remove it, create a fresh one, continue.

### Step 2: Phase 1 — Orient (read-only)

Build a map of current state before changing anything. **Do not load full page bodies yet.**

1. Read `<index>` in full — it's the navigation.
2. List `<wiki>/*.md`. For each page (excluding the index files and `log.md`):
   - Parse YAML frontmatter (the block between the first two `---` markers). Capture `name`, `type`, `last_updated`, `updated`, `created`, plus any `sources:` / `trust:` / `confidence:` fields you encounter.
   - Extract the H1 (`^#\s+(.+)$` on the first matching line). Fall back to the file stem.
   - Record the one-line summary (the first non-empty paragraph after the H1, truncated to ~120 chars). This is enough — body parsing happens lazily in Phase 3.
3. Read the most recent entries of `<dir>/log.md` (or `<wiki>/log.md` if that's where it lives) — enough to know what was touched recently and what dates are anchorable.

You now know what exists, what's typed how, and what was touched recently.

### Step 3: Phase 2 — Gather signal (read-only deterministic checks)

Run all seven checks below and collect findings into a worklist for Phase 3. **Each check is mechanical** — the LLM is doing what a script would do, deterministically.

**Setup.** Build `pages = {stem -> info}` from the Phase-1 map. For each page, also extract:
- `links`: every `[[target]]` match via regex `\[\[([^\]]+?)\]\]`. Normalize each target (strip `.md`, take basename after last `/`, strip surrounding whitespace).
- `relative_dates`: every match of the relative-date regex (see Check 6 below).

Then run each check in turn:

**Check 1 — Pages on disk missing from index.**
Pull `index_stems` from `<index>` via the same wiki-link regex plus markdown links: `\]\(([^)]+\.md)\)`. Normalize identically. Compute `set(pages.keys()) - index_stems`. Worklist payload: `Add to index: <stem>`.

**Check 2 — Index entries with no matching page.**
Compute `index_stems - set(pages.keys())`. Worklist payload: `Remove from index (dangling): <stem>`.

**Check 3 — Broken `[[wiki-links]]`.**
For each `(stem, info)` in `pages`, for each `target` in `info.links`: if `target not in pages`, add `<stem>.md -> [[<target>]]` to the worklist.

**Check 4 — Orphan pages (no inbound links).**
Initialize `inbound = {stem: 0 for stem in pages}`. For each `(stem, info)`, for each `target` in `info.links`: if `target in inbound and target != stem`, increment `inbound[target]`. Worklist: every stem with `inbound == 0`. Note: orphans aren't automatically wrong (some pages are genuine top-level entities) — Phase 3 decides per case.

**Check 5 — Stale pages (frontmatter date older than 90 days).**
For each page, prefer `last_updated`, then `updated`, then `created`. Parse with `%Y-%m-%d` or `%Y/%m/%d`. Compute `(today - parsed_date).days`. If `> 90`, add `<stem>.md (updated <N> days ago)` to the worklist. Use the system's current date as `today`.

**Check 6 — Relative-date occurrences.**
For each page body, grep with this regex (case-insensitive, word-bounded):

```
\b(yesterday|today|tomorrow|tonight|last\s+(?:night|week|month|year|monday|tuesday|wednesday|thursday|friday|saturday|sunday)|next\s+(?:week|month|year|monday|tuesday|wednesday|thursday|friday|saturday|sunday)|this\s+(?:week|month|year|morning|afternoon|evening)|a\s+(?:few\s+)?(?:days?|weeks?|months?)\s+ago|recently|just\s+now|earlier\s+today|the\s+other\s+day)\b
```

Shell equivalent (preferred if `grep` is available):

```
grep -E -i -no '<the regex above>' <wiki>/*.md
```

Worklist payload per page: `<stem>.md — <comma-separated unique lowercased matches>`.

**Check 7 — Likely-duplicate pages by title.**
Compare every unordered pair of page titles (lowercased, whitespace-normalized). Flag a pair as a likely duplicate if **any** of these hold:
- One title is a substring of the other (after stripping trailing punctuation).
- They share ≥80% of word tokens (count tokens that appear in both / count tokens in the shorter title).
- They differ only by trailing plural / version suffix / punctuation (e.g., `redis-decision` vs `redis-decisions`, `migration-v1` vs `migration-v2`).

This approximates SequenceMatcher ratio ≥0.82. Worklist payload: `<a>.md ~ <b>.md (similar title)`.

**Bonus check — Sources newer than citing pages.** If `<dir>/sources/` exists, list each source file's mtime (`ls -la --time-style=long-iso <dir>/sources/` or equivalent). For each page whose body contains a literal reference to a source filename, compare the source's mtime with the page's `last_updated`. Flag any page older than its source.

### Step 3.5: Phase 2.5 — Preview & Confirm (consent gate before Phase 3)

The worklist is now complete and the disk is still untouched. Before any Phase 3 write, present the worklist to the user and require explicit consent.

1. **Print the preview.** Render the Phase 2 findings in the exact format the final summary uses for the "Phase 2 — Gather signal" block (counts + identifying stems for each of the seven checks, plus the bonus check if it fired). The user must be able to see *what* will be touched, not just *how many*.

2. **Ask for consent.** Use your environment's ask-user prompt with the exact wording:

   ```
   Apply these N fixes? [all / dry-run / specific clusters / none]
   ```

   Replace `N` with the total worklist item count.

3. **Resolve the answer:**
   - `all` — proceed to Phase 3 with the full worklist.
   - `dry-run` — skip Phase 3 entirely, jump to Phase 4 to release the lock and emit the summary marked `(dry-run)`. Do not bump `last_updated` on any page. Do not append a `[dream]` line to `log.md` (a dry-run is not a pass).
   - `specific clusters` — ask the user to name which clusters (e.g., "merges only", "links only", "duplicates 1 and 3"), filter the worklist to just those, then proceed.
   - `none` — release the lock and exit without writing.
   - **Ambiguous response or no response** (timeout, unparseable token, anything not in the four above) — default to `dry-run`. Never escalate to `all` on uncertainty.

4. **Mode-token short-circuit.** If the user invoked `do-work dream --dry-run`, skip the ask-user prompt entirely and treat the answer as `dry-run`. Still print the preview — that's the whole point of dry-run mode.

5. **Record the choice** in the eventual summary's Phase 3 section (e.g., `Phase 3 — Consolidate (mode: all)` or `Phase 3 — Skipped (mode: dry-run)`).

### Step 4: Phase 3 — Consolidate and heal (this is where edits happen)

Work through the worklist. **Mechanical fixes first** (low-risk), then **semantic fixes** (judgment-heavy).

**Mechanical (act directly on the linter output):**

- **Index drift (Checks 1 + 2):** add missing pages to `<index>` under the appropriate type heading, format `- [[<stem>]] — <one-line summary>`. Remove every dangling line.
- **Broken `[[wiki-links]]` (Check 3):** for each, choose one of:
  - **Rename**: if the target is an obvious typo of an existing page, fix the link.
  - **Stub**: if the page genuinely should exist, create a minimal page with frontmatter (`name`, `type`, `last_updated: <today>`) and a `TODO: stub created by dream` body.
  - **Drop**: if the link no longer makes sense, remove it from the source page.
- **Orphan pages (Check 4):** link from a related page when one exists, or accept as a top-level entity and note in the summary.

**Semantic (your judgment — the part a script can't do):**

- **Contradictions — newest wins.** When two pages, or two sections of one page, assert conflicting facts, edit to state the current truth. Don't leave both standing. Tie-break by newer `last_updated`, then by `trust`/`confidence` if present, then by what `log.md` shows happened most recently. Never leave both claims in place.
- **Relative → absolute dates (Check 6).** Rewrite every flagged phrase into an absolute date you can establish from `log.md` or frontmatter. Example: "Yesterday we switched to Redis" → "On 2026-05-18 we switched to Redis." If you cannot establish the date, rewrite to lose the false precision ("recently" → "at some point") and flag the page in the eventual summary — never guess.
- **Merge near-duplicates (Check 7).** Read both bodies. If they cover the same concept: keep the richer one, fold in the loser's unique facts, **repoint every inbound `[[link]]`** in other pages to the survivor, delete the loser, update the index. If the bodies are genuinely distinct (similar title only), leave both and note that in the summary.
- **Prune the untrue.** Remove facts no longer accurate, transient state that was wrongly made durable, and anything the user has asked to forget. Deletion must surface in the summary and `log.md` — never silent.
- **Tighten transcript creep.** Any page that reads like a chat log → rewrite into compiled understanding. The synthesis is the value, not the raw exchange.

**Rules during Phase 3:**
- Edit in place.
- Never touch `<dir>/sources/` — it's immutable ground truth.
- Every deletion or merge must be visible in the eventual summary.

### Step 5: Phase 4 — Prune and reindex

**If Step 3.5 resolved to `dry-run` or `none`** (no Phase 3 writes occurred): skip substeps 1–4 (they're all writes), do only substep 5 (release the lock), and Step 6's summary is marked `(dry-run)` or `(declined)` accordingly.

1. **Rebuild the index** (`<index>`):
   - One line per page, grouped by `type` frontmatter field (or by directory if no type is set).
   - Format: `- [[<stem>]] — <one-line summary>`.
   - Keep the file ≤200 lines / ~25 KB (the Claude Code Auto Memory ceiling — preserves compatibility with native tooling).
2. **Demote any verbose prose** that crept into the index back into its topic page.
3. **Bump `last_updated`** to today's absolute date (`YYYY-MM-DD`) on every page you edited in Phase 3.
4. **Append to `log.md`** (top of file, one line):
   ```
   - YYYY-MM-DD [dream] — N merged, M pruned, K links fixed, L dates pinned
   ```
   Use the system's current date and the actual counts from Phase 3.
5. **Remove `<dir>/.lock`.**

### Step 6: Emit summary

Print a concise, honest report. The user should be able to `git diff` the memory directory and see exactly what dream did — the diff is the real audit trail.

## Output Format

```
# Dream pass: <dir>

**Phase 1 — Orient**
Found <N> pages in <wiki>, index = <MEMORY.md / _master_index.md / index.md>.
Most recent log entry: <date> — <activity>.

**Phase 2 — Gather signal**
- Missing from index: <count> (<stems>)
- Index dangling: <count> (<stems>)
- Broken wiki-links: <count> (<from -> target>)
- Orphan pages: <count> (<stems>)
- Stale pages (>90d): <count> (<stems with days>)
- Relative dates: <count> across <N> pages
- Likely duplicates: <count> (<pairs>)

**Phase 3 — Consolidate**
- Merged: <count> (<from -> into>)
- Pruned: <count> (<stems>)
- Links fixed: <count>
- Dates pinned: <count>
- Contradictions resolved: <count> (<one-line description each>)
- Deliberately left alone: <count> (<reason — e.g., genuine top-level orphan, unsourceable date>)

**Phase 4 — Reindex**
- <index> rebuilt: <N> lines
- log.md entry: - <YYYY-MM-DD> [dream] — <counts>
- .lock removed

Audit: run `git diff <dir>` for the full record of changes.
```

## Rules

- **Manual only.** Never schedule dream on a timer, hook, or background trigger. If asked, refuse and explain that consolidation is destructive.
- **Phases 1–2 are read-only.** No writes until Phase 3. If you find yourself editing before the worklist is complete, stop and restart.
- **Phases 1–2 must produce a visible worklist before Phase 3 may begin.** The scan phase makes zero writes. The Phase 2.5 preview + consent gate is non-optional — even on a re-run, even when the user "just wants the same fixes as last time."
- **Hold `.lock` for the full pass.** Create in Step 1, remove in Phase 4. Never skip the lock — concurrent passes corrupt the memory.
- **Never touch `<dir>/sources/`.** It's immutable ground truth. If `sources/` is modified, the pass has violated provenance.
- **Newest verified fact wins** on every contradiction. Never leave both claims standing.
- **Edit in place.** No silent deletions; every removal mentioned in the summary.
- **If no memory dir is found, stop.** Do not initialize one — dream consolidates, it does not create.
- **`_master_index.md` is a first-class index.** Auto-detect it so bkb wikis work without extra arguments.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "I'll skip the seven checks — I can spot the issues by reading the pages" | Run all seven checks first; collect the worklist before any judgment | Mechanical checks are deterministic and free; semantic judgment is expensive. Get the worklist before spending tokens. |
| "I'll merge these two pages without checking inbound links" | Repoint every inbound `[[link]]` before deleting the loser | Silent link breakage is the worst kind of rot — it propagates to every page that referenced the loser. |
| "I can't establish a date for 'recently' — I'll guess from context" | Rewrite to lose the false precision ("recently" → "at some point") and flag in the summary | A confident wrong date is worse than vague phrasing. |
| "The index looks fine — I'll skip the rebuild" | Always rebuild `<index>` in Phase 4 | Edits in Phase 3 may have added or removed pages; the index must reflect disk. |
| "I'll create the memory dir if it doesn't exist — saves the user a step" | Stop with the "no memory directory found" message | Dream is consolidation, not initialization. Creating a dir silently sets up a future Phase-3 nuking of empty defaults. |
| "Just one tiny edit to `sources/` to fix a typo" | Leave `sources/` alone, even for typos | Provenance is sacred. If a source is wrong, that's a sources-management problem, not a dream one. |
| "The user already typed `do-work dream`, that's consent enough — skip the Phase 2.5 ask" | Always present the Phase 2.5 worklist preview and the `Apply these N fixes?` prompt | The invocation token is single-bit consent. The user can't consent to a worklist they haven't seen. Phase 2.5 is the only point where they see what's about to change. |
| "The worklist is short — just one merge, I'll do it without asking" | Run the Phase 2.5 ask anyway; default to `dry-run` if the response is unclear | One destructive write is enough to lose work. Short worklists are exactly when shortcuts feel safe and don't belong. |

## Red Flags

- `.lock` still on disk after the pass — Phase 4 didn't complete; rerun or investigate.
- `log.md` has no new `[dream]` entry but pages were modified — the audit trail is broken.
- Summary reports "0 merged, 0 pruned, 0 links fixed" but `git diff` shows page edits — the action wrote silently.
- A page was deleted but the index still lists it (or vice versa) — Phase 4 reindex was skipped.
- `<dir>/sources/` was modified — provenance was violated; the pass is invalid.
- The rebuilt `<index>` is over 200 lines / 25 KB — Phase 4 budget exceeded; verbose content leaked in.
- The summary lists merges but inbound `[[links]]` to the deleted page still exist in other pages — repoint step was skipped.
- The summary lists Phase 3 work but no Phase 2.5 preview was printed earlier in the run — the consent gate was bypassed.
- `--dry-run` was passed but `log.md` gained a `[dream]` line or pages were edited — the dry-run contract was violated.

## Verification Checklist

- [ ] Target directory resolved (or stopped with "no memory directory found")
- [ ] Index file (`MEMORY.md` / `_master_index.md` / `index.md`) located before Phase 1
- [ ] `.lock` created in Step 1 and removed in Phase 4
- [ ] All seven Phase-2 checks executed; findings collected into a worklist
- [ ] Phase 2.5 confirmation was presented (worklist preview + `Apply these N fixes? [all / dry-run / specific clusters / none]`) before any Phase 3 write
- [ ] No Phase 3 writes occurred if the user declined or chose `dry-run`
- [ ] Mechanical fixes applied before semantic fixes
- [ ] Every edited page has bumped `last_updated` to today's absolute date
- [ ] `<index>` rebuilt and ≤200 lines / ~25 KB
- [ ] `log.md` gained one `- YYYY-MM-DD [dream] — N merged, M pruned, K links fixed, L dates pinned` line at the top
- [ ] `<dir>/sources/` untouched (if it exists)
- [ ] Every merge repointed inbound `[[links]]` before deletion
- [ ] Summary is concise, honest, and auditable against `git diff <dir>`
