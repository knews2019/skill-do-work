# Dream

A manual, four-phase consolidation pass over a plain-text memory directory — find rot, heal contradictions, merge near-duplicates, rebuild the index. Destructive by design, so it never runs automatically and never runs without a preview-then-confirm round.

> **Not to be confused with cleanup, stray-check, or bkb commands.** `do-work cleanup` works on do-work's own archive bookkeeping. `do-work stray-check` works on repo-wide file hygiene. `do-work bkb lint` produces a read-only health report on a bkb wiki. `do-work bkb garden` does conservative relationship hygiene. Dream is the destructive, one-pass rewrite of an existing memory store after lint has been flagging the same issues for a while.

## What it operates on

A plain-text memory directory shaped like this:

| Component | Recognized names | Required? |
|-----------|------------------|-----------|
| **Index** | `MEMORY.md`, `_master_index.md`, `index.md` | Yes — refuses to dream without one |
| **Wiki pages** | `*.md` files in `<dir>/wiki/`, `<dir>/pages/`, or `<dir>` itself | Yes |
| **Log** | `log.md` at `<dir>/log.md` or `<dir>/wiki/log.md` | Read in Phase 1, appended to in Phase 4 |
| **Sources** | `<dir>/sources/` | Read-only ground truth — dream never edits this folder |
| **Lock** | `<dir>/.lock` | Created in Step 1, removed in Phase 4. Two passes never overlap. |

If nothing matches, dream refuses to initialize anything and points at `do-work bkb init`.

## The four phases

| Phase | Step | What happens | Writes? |
|-------|------|--------------|---------|
| **1 — Orient** | Step 2 | Map the index, list pages, read recent log entries. | No |
| **2 — Gather signal** | Step 3 | Run seven mechanical checks (index drift, broken `[[links]]`, orphan pages, stale frontmatter, relative dates, likely duplicates, sources newer than citers). Collect findings into a worklist. | No |
| **2.5 — Preview & Confirm** | Step 3.5 | Print the worklist preview, ask `Apply these N fixes? [all / dry-run / specific clusters / none]`. Ambiguous responses default to `dry-run`. | No |
| **3 — Consolidate** | Step 4 | Mechanical fixes first (index drift, broken links), then semantic ones (contradictions → newest wins, relative → absolute dates, merge near-duplicates with inbound-link repointing, prune the untrue, tighten transcript creep). | **Yes** — but only with the consent from Step 3.5. |
| **4 — Reindex** | Step 5 | Rebuild the index (≤200 lines / ~25 KB), bump `last_updated` on edited pages, append a `- YYYY-MM-DD [dream] — N merged, M pruned, K links fixed, L dates pinned` line to `log.md`, remove `.lock`. | Yes (or just lock release on dry-run). |

## The Phase 2.5 consent gate

The single-bit `do-work dream` invocation token is **not enough** to consent to a Phase 3 worklist you haven't seen. After Phase 2 builds the worklist, dream prints it and asks:

```
Apply these N fixes? [all / dry-run / specific clusters / none]
```

| Answer | Effect |
|--------|--------|
| `all` | Proceed to Phase 3 with the full worklist. |
| `dry-run` | Skip Phase 3 entirely. Phase 4 only releases the lock; no `[dream]` log entry, no `last_updated` bumps. |
| `specific clusters` | Dream asks which ones (e.g., "merges only", "duplicates 1 and 3"), filters, then proceeds. |
| `none` | Release the lock, exit, write nothing. |
| Ambiguous / no response | Default to `dry-run`. Never escalates to `all` on uncertainty. |

The `--dry-run` mode token short-circuits the prompt: it still prints the preview, but auto-answers `dry-run`.

## What gets written in Phase 3

- **Index drift** — missing pages added under their type heading, dangling entries removed.
- **Broken `[[wiki-links]]`** — renamed (if a typo of an existing page), stubbed (with `TODO: stub created by dream` body), or dropped (if no longer meaningful).
- **Orphan pages** — linked from a related page when one exists, or accepted as a top-level entity (and named in the summary).
- **Contradictions** — newest verified fact wins. Tie-break by `last_updated`, then `trust`/`confidence`, then `log.md` recency. Never leave both claims standing.
- **Relative dates** — rewritten to absolute dates established from `log.md` or frontmatter. If unsourceable, the phrase loses the false precision ("recently" → "at some point") and the page is flagged in the summary.
- **Near-duplicates** — richer page wins, loser's unique facts folded in, **every inbound `[[link]]` in other pages repointed to the survivor**, loser deleted, index updated.
- **Transcript creep** — chat-log–shaped pages rewritten as compiled understanding.
- **Sources untouched** — `<dir>/sources/` is immutable ground truth. Even typos stay.

## Output

A concise, honest summary keyed to each phase. The real audit trail is `git diff <dir>` — dream's summary should match what the diff shows, line for line.

## Usage

```
do-work dream                       Resolve default memory dir, run the full four phases
do-work dream ./memory              Run on a specific dir
do-work dream --dry-run             Build the worklist, print preview, exit without writing
do-work dream --dry-run ./kb/wiki   Both — preview the bkb wiki without touching it
do-work consolidate memory          Same as `do-work dream`
do-work clean up wiki               Same — dream wins over the cleanup verb when "wiki" / "memory" / "notes" appear
do-work lint and merge notes        Same
```

## Key rules

- **Manual only.** Never on a timer, hook, or background trigger. Asked to schedule? Refuse and explain.
- **Phases 1–2 are read-only.** No writes until the Phase 2.5 gate has been answered.
- **Phase 2.5 is non-optional.** Even on re-runs. Even when the user "wants the same fixes as last time."
- **Hold `.lock` for the full pass.** Concurrent passes corrupt memory.
- **`<dir>/sources/` is sacred.** Provenance is never edited, even for typos.
- **Newest verified fact wins** on every contradiction. No "both claims survive" outcomes.
- **Every deletion is named in the summary.** Silent removals are forbidden.

## When NOT to use

- No memory directory exists → `do-work bkb init`.
- You want a read-only health report on a bkb wiki → `do-work bkb lint`.
- You want routine, conservative relationship hygiene on a bkb wiki → `do-work bkb garden`.
- You want structural defragmentation → `do-work bkb defrag`.
- You want repo-wide file hygiene (not memory) → `do-work stray-check`.
- You want do-work archive consolidation → `do-work cleanup`.
