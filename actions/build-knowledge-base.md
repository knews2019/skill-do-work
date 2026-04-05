# Build Knowledge Base Action

> **Part of the do-work skill.** Implements the LLM Knowledge Base pattern — a persistent, compounding Markdown wiki compiled from raw source documents.

The core idea: instead of re-deriving knowledge from scratch on every query (RAG), the LLM incrementally compiles raw sources into a structured, interlinked Markdown wiki. `raw/` is the source code, the LLM is the compiler, the wiki is the executable.

## Sub-Commands

The `bkb` command accepts a sub-command as its first argument. If no sub-command is given, show the help menu.

| Sub-command | What it does |
|---|---|
| `init [path]` | Initialize a new knowledge base at the given path (default: `./kb`) |
| `triage` | Sort inbox items into capture subdirectories, update the queue |
| `ingest [target]` | Compile source(s) into wiki pages (file, batch, or "today") |
| `query [question]` | Search the wiki and synthesize an answer |
| `lint [scope]` | Health check — contradictions, orphans, broken links, stale claims |
| `resolve` | Walk through open contradictions and resolve them one by one |
| `close` | Finalize the daily log, verify indexes, report summary |
| `rollup` | Monthly rollup — trends, volume stats, recommendations |
| `status` | Show current KB stats — article counts, pending items, recent activity |
| (none) | Show help menu |

---

## Locating the Knowledge Base

Before executing any sub-command (except `init`), find the KB root:

1. Check if `$ARGUMENTS` includes an explicit `--kb <path>` flag — use that path.
2. Look for a `kb/` directory in the current working directory.
3. Look for a `knowledge-base/` directory in the current working directory.
4. Search parent directories (up to 3 levels) for a directory containing both `raw/` and `wiki/` subdirectories.
5. If not found, tell the user: "No knowledge base found. Run `do work bkb init` to create one."

---

## Sub-Command: `init [path]`

Create the full KB directory structure at the specified path (default: `./kb`).

### Step 1: Create the Raw Pipeline

```
<path>/
├── raw/
│   ├── inbox/                      # Zero-friction drop zone
│   ├── capture/                    # Type-sorted staging
│   │   ├── web/
│   │   ├── papers/
│   │   ├── repos/
│   │   ├── images/
│   │   ├── notes/
│   │   ├── audio/
│   │   └── video/
│   ├── processed/                  # Ingested sources, organized by date (YYYY-MM-DD/)
│   └── _inbox_queue.md             # LLM work list
```

### Step 2: Create the Wiki Structure

```
<path>/
├── wiki/
│   ├── _master_index.md            # Top-level nav (~50 lines max)
│   ├── topics/                     # Second-level indexes
│   ├── concepts/                   # Concept articles
│   ├── entities/                   # Person/org/product articles
│   ├── sources/                    # Source summaries
│   ├── comparisons/                # Filed query outputs
│   ├── daily/                      # Daily changelog
│   ├── monthly/                    # Monthly rollup + trends
│   ├── log.md                      # Append-only timeline
│   └── overview.md                 # High-level synthesis
```

### Step 3: Create Seed Files

**`raw/_inbox_queue.md`:**

```markdown
# Inbox Queue

Items pending triage. Updated automatically during triage.

| # | File | Source Type | Status |
|---|---|---|---|
```

**`raw/processed/_manifest.md`:**

```markdown
# Processing Manifest

| File | Date Processed | Processed Path | Wiki Articles Produced | Status |
|---|---|---|---|---|
```

**`wiki/_master_index.md`:**

```markdown
# Master Index

Last updated: {today} | Total articles: 0 | Topic clusters: 0

## Topic Clusters

(none yet — run `do work bkb ingest` to add your first source)

## Recent Activity

- {today}: Knowledge base initialized
```

**`wiki/log.md`:**

```markdown
# Activity Log

## [{today}] init | Knowledge base created

Structure initialized. Ready for first source.
```

**`wiki/overview.md`:**

```markdown
# Knowledge Base Overview

This knowledge base was initialized on {today}. No sources have been ingested yet.

Add sources to `raw/inbox/` and run `do work bkb triage` followed by `do work bkb ingest` to begin building.
```

### Step 4: Create the Schema File

Create `<path>/CLAUDE.md` with the KB schema (conventions, frontmatter format, workflow triggers). Use the schema content from the "Schema File" section below.

### Step 5: Initialize Git

If the KB path is not already inside a git repository, run `git init` in the KB root.

### Step 6: Report

```
Knowledge base initialized at <path>/

  raw/     Source pipeline (inbox → capture → processed)
  wiki/    LLM-maintained wiki (master index → topic indexes → articles)

Next steps:
  Drop files into <path>/raw/inbox/
  do work bkb triage         Sort inbox items
  do work bkb ingest         Compile sources into wiki
```

---

## Sub-Command: `triage`

Sort new items from `raw/inbox/` into `raw/capture/` subdirectories by type.

### Steps

1. **Scan** `raw/inbox/` for all files (non-recursive — files only, not subdirectories).
2. **Classify** each file by extension and content:
   - `.pdf` → `capture/papers/`
   - `.md` from web clippers (has URL in frontmatter or content) → `capture/web/`
   - `.md` that looks like personal notes → `capture/notes/`
   - `.png`, `.jpg`, `.jpeg`, `.gif`, `.svg`, `.webp` → `capture/images/`
   - `.mp3`, `.wav`, `.m4a`, `.ogg` → `capture/audio/`
   - `.mp4`, `.webm`, `.mkv`, transcript files → `capture/video/`
   - Code files, `README.md` from repos → `capture/repos/`
   - Unknown → leave in inbox, flag for user
3. **Move** each classified file to its target directory. If a file with the same name already exists in the target (from a previous triage that hasn't been ingested yet), prefix with the current time: `HHMMSS-filename.ext`.
4. **Update** `raw/_inbox_queue.md` — append only the files moved from inbox in **this** triage pass, marked as "ready". Do NOT re-scan all of `capture/`; the queue is an append-only ledger of triage batches.
5. **Report**: Items triaged, items skipped (with reasons), items ready for ingestion.

If inbox is empty, say so and suggest adding files.

---

## Sub-Command: `ingest [target]`

Compile source documents into wiki pages. This is the core operation.

### Target Resolution

- `ingest` (no target) → process all "ready" items in `raw/_inbox_queue.md`
- `ingest <filename>` → process a specific file from `raw/capture/`
- `ingest <path>` → process a specific file by path (can be outside `capture/` — e.g., a file the user hasn't triaged yet)

### Steps

1. **Read** the target source file(s) from `raw/capture/` (or the specified path).
2. **Handle non-text sources**:
   - **Images** (`.png`, `.jpg`, `.svg`, etc.): Use LLM vision to describe the image. Generate a summary from the visual content. If a companion `.md` file exists alongside (e.g., `diagram.png` + `diagram.md`), use both.
   - **Audio** (`.mp3`, `.wav`, etc.): Check for a companion transcript file (e.g., `podcast.mp3` + `podcast.txt` or `podcast.md`). If found, ingest the transcript. If no transcript exists, skip the file and flag it: "Audio file needs a transcript — add a .txt or .md alongside it."
   - **Video** (`.mp4`, `.webm`, etc.): Same as audio — look for a companion transcript. Skip and flag if none found.
   - **Text files** (`.md`, `.pdf`, `.txt`, code files): Process normally.
3. **For each source**, discuss key takeaways briefly, then:
   a. **Duplicate check** — before creating any wiki page, search for existing pages covering the same topic:
      - **Exact duplicate** (same source re-ingested, or same content from a different URL): update the existing page — add the new file to its `sources:` frontmatter list, refresh any stale claims, note "additional source" in the daily log. Do NOT create a second page.
      - **Near-duplicate** (same topic, different angle or data): create a separate page but add bidirectional `[[wiki-links]]` in both pages' `related:` frontmatter. Note the relationship in the daily log (e.g., "New page X complements existing page Y").
      - **No duplicate**: proceed normally.
   b. **Create or update** a summary page in `wiki/sources/` with YAML frontmatter.
   c. **Create or update** relevant concept pages in `wiki/concepts/`.
   d. **Create or update** relevant entity pages in `wiki/entities/`.
   e. **Check for contradictions** with existing wiki content. Flag any found.
   f. **Update cross-references** — add `[[wiki-links]]` between related pages.
4. **Update indexes**:
   a. Determine which topic cluster(s) the new content belongs to.
   b. Create new topic index (`wiki/topics/_index_[topic].md`) if needed.
   c. Update existing topic index(es) with new article entries.
   d. Update `wiki/_master_index.md` — article counts, topic list, recent activity.
5. **Write daily log**: Create or append to `wiki/daily/{today}.md` listing everything ingested, created, updated, and any contradictions flagged.
6. **Move to processed and mark done** (per-file): After each file completes steps 3–5 successfully, immediately move it to `raw/processed/{today}/` (create the date directory if needed) from wherever it currently lives (`raw/capture/`, `raw/inbox/`, or an external path). If the file was in the queue, mark its row as "done" in `raw/_inbox_queue.md`; if it was ingested directly by path (bypassing triage), add a "done" entry to the queue for traceability. If a file with the same name already exists in the target directory, prefix with the current time: `HHMMSS-filename.ext`. Update `raw/processed/_manifest.md` with the original path, processed path, and wiki articles produced. **This is per-file, not per-batch** — if file 4 of 5 fails, files 1–3 are already safely processed and marked done.
7. **Append to activity log**: Add entry to `wiki/log.md`.
8. **Report**: Sources processed, pages created/updated, contradictions found, skipped files (with reasons), index changes.

### Page Conventions

Every wiki page MUST have YAML frontmatter:

```yaml
---
title: Page Title
type: concept | entity | source-summary | comparison | daily-log | monthly-rollup
topic_cluster: [which topic index this belongs to]
sources: [list of raw/processed/ paths — the file's final stable location]
related: [list of wiki pages linked]
created: YYYY-MM-DD
updated: YYYY-MM-DD
confidence: high | medium | low
---
```

> **`sources:` always uses the `raw/processed/` path** (e.g., `raw/processed/2026-04-05/moe-paper.pdf`). This is the file's stable final location. Never use `capture/` paths — those are transient.

### Index Size Rules

- `_master_index.md` must stay under 80 lines.
- Each topic index must stay under 60 lines.
- When a topic index exceeds 80 articles, split it and update `_master_index.md`.
- Every wiki article must appear in exactly one topic index.
- Every topic index must appear in `_master_index.md`.

---

## Sub-Command: `query [question]`

Search the wiki to answer a question. Uses two-hop navigation.

### Steps

1. **Read** `wiki/_master_index.md` to identify relevant topic cluster(s).
2. **Read** the relevant topic index(es) from `wiki/topics/`.
3. **Read** the specific articles identified as relevant (2–5 articles typically).
4. **Synthesize** an answer using `[[wiki-link]]` citations to wiki pages.
5. **Offer to file**: If the answer is substantive, ask the user if they want to save it as a new wiki page in `wiki/comparisons/`.
6. **If filed**: Create the page with proper frontmatter, update the relevant topic index and `_master_index.md`, append to `wiki/log.md`.

If the wiki has no content on the topic, say so and suggest sources to ingest.

---

## Sub-Command: `lint [scope]`

Health check the wiki for consistency and accuracy.

### Scope

- `lint` (no scope) → quick lint (changed pages since last lint)
- `lint [topic]` → full check of one topic cluster
- `lint full` → cross-cluster structural integrity check
- `lint all` → complete re-verification (quarterly cadence)

### Checks

1. **Contradictions** — pages that make conflicting claims about the same thing.
2. **Orphan pages** — wiki pages with no inbound `[[wiki-links]]`.
3. **Missing pages** — concepts mentioned 3+ times across pages without their own page.
4. **Stale claims** — content superseded by newer sources (check `updated` dates).
5. **Index integrity**:
   - Every wiki article appears in exactly one topic index.
   - Every topic index appears in `_master_index.md`.
   - Article counts in indexes match actual file counts.
   - No topic index exceeds the split threshold (80 articles).
6. **Broken links** — `[[wiki-links]]` pointing to non-existent pages.
7. **Daily log coverage** — daily logs exist for every day that had ingestion activity.
8. **Frontmatter completeness** — all required fields present in every wiki page.

### Report

Write findings to `wiki/daily/{today}.md` and print a summary:

```
Lint results (scope: [scope]):
  Contradictions: N
  Orphan pages: N
  Missing pages suggested: N
  Stale claims: N
  Index issues: N
  Broken links: N
```

Suggest specific fixes for each finding.

---

## Sub-Command: `resolve`

Walk through open contradictions and resolve them interactively.

### Steps

1. **Find open contradictions**: Scan `wiki/daily/` logs and `wiki/log.md` for entries containing "contradiction" or "conflicting". A contradiction is **open** if no subsequent log entry marks it as `[RESOLVED]`. Convention: resolution log entries use the format `[RESOLVED] contradiction: <description>` so they can be matched against the original flag.
2. **For each contradiction**, present it to the user:
   - Show the two (or more) conflicting claims with their source pages and original raw sources.
   - Propose a resolution: which claim is more recent, better sourced, or more authoritative?
   - Ask the user to confirm, adjust, or skip.
3. **Apply resolution**: Update the wiki page(s) — correct the stale/wrong claim, add a note about what changed and why, update the `confidence:` frontmatter if needed.
4. **Log resolution**: Append to `wiki/log.md` and `wiki/daily/{today}.md` using the format `[RESOLVED] contradiction: <description>` — this marks it as closed so future `resolve` runs skip it. Include which pages were updated and how.
5. **Report**: Contradictions resolved, contradictions skipped, contradictions remaining.

If no open contradictions are found, say so.

---

## Sub-Command: `close`

Finalize the day's work.

### Steps

1. **Finalize daily log**: Ensure `wiki/daily/{today}.md` has all changes from today.
2. **Verify indexes**: Check that `_master_index.md` article counts are accurate.
3. **Refresh overview**: Re-read the current state of the wiki (master index, recent daily logs, topic clusters) and regenerate `wiki/overview.md` with an up-to-date high-level synthesis of what the knowledge base covers.
4. **Report**:
   ```
   Day closed ({today}):
     Articles created: N
     Articles updated: N
     Sources ingested: N
     Contradictions pending: N
   ```
5. **Suggest git commit** (do not auto-commit): If there are uncommitted changes in the KB directory, print: "Uncommitted KB changes — run `do work commit` or `git add . && git commit` when ready."

---

## Sub-Command: `rollup`

Generate the monthly summary. Run on the 1st of each month or on demand.

### Steps

1. **Read** all `wiki/daily/` entries for the target month (default: previous month).
2. **Create** `wiki/monthly/{YYYY-MM}.md` with:
   - **Volume stats**: sources ingested, articles created/updated, contradictions found/resolved.
   - **Theme evolution**: which topics grew, which went stale, emerging patterns.
   - **Integrity summary**: lint results, confidence changes.
   - **Recommendations**: topic splits needed, new clusters suggested, gap areas.
3. **Evaluate** whether any topic index needs splitting (threshold: 80+ articles).
4. **Update** `wiki/_master_index.md` with monthly activity line.
5. **Append** to `wiki/log.md`.

---

## Sub-Command: `status`

Quick snapshot of the KB state.

### Steps

1. **Read** `wiki/_master_index.md` for article counts and topic clusters.
2. **Count** files in `raw/inbox/` (pending triage) and items marked "ready" in `raw/_inbox_queue.md` (pending ingestion).
3. **Read** the most recent `wiki/daily/` entry for last activity date.
4. **Report**:
   ```
   Knowledge Base Status:
     Location: <path>
     Total articles: N across M topic clusters
     Inbox: N items pending triage
     Queue: N items ready for ingestion
     Last activity: YYYY-MM-DD
     Last lint: YYYY-MM-DD (or "never")
   ```

---

## Help Menu

When invoked with no sub-command or with `help`:

```
do work bkb — LLM Knowledge Base builder

  Setup:
    do work bkb init              Initialize a new knowledge base in ./kb
    do work bkb init ~/research   Initialize at a custom path

  Daily workflow:
    do work bkb triage            Sort inbox items into capture directories
    do work bkb ingest            Compile all ready sources into wiki
    do work bkb query [question]  Search the wiki and synthesize an answer
    do work bkb close             Finalize today's daily log

  Maintenance:
    do work bkb lint              Quick health check (recent changes)
    do work bkb lint full         Full cross-cluster integrity check
    do work bkb resolve           Walk through and resolve contradictions
    do work bkb rollup            Monthly summary and trend analysis
    do work bkb status            Show KB stats and pending items

  Typical flow:
    1. Drop files into kb/raw/inbox/
    2. do work bkb triage
    3. do work bkb ingest
    4. do work bkb query "what are the tradeoffs of X vs Y?"
    5. do work bkb close
```

---

## Schema File Content

When `init` creates `<path>/CLAUDE.md`, use this content:

```markdown
# LLM Knowledge Base Schema

## Project Structure
- `raw/` — source documents with lifecycle pipeline. NEVER modify originals.
- `raw/inbox/` — zero-friction drop zone. Sort into capture/ before processing.
- `raw/capture/` — type-sorted staging area.
- `raw/processed/YYYY-MM-DD/` — ingested sources, moved here after successful compilation.
- `raw/_inbox_queue.md` — append-only triage ledger. Only updated with files moved in the current triage pass.
- `wiki/` — LLM-generated wiki. You own this entirely.
- `wiki/_master_index.md` — top-level catalog. Read FIRST on every query.
- `wiki/topics/_index_[topic].md` — second-level indexes by topic cluster.
- `wiki/daily/YYYY-MM-DD.md` — daily changelog.
- `wiki/monthly/YYYY-MM.md` — monthly rollup and trends.
- `wiki/log.md` — append-only activity log.

## Page Conventions
Every wiki page MUST have YAML frontmatter:

    ---
    title: Page Title
    type: concept | entity | source-summary | comparison | daily-log | monthly-rollup
    topic_cluster: [which topic index this belongs to]
    sources: [list of raw/processed/ paths — stable final location]
    related: [list of wiki pages linked]
    created: YYYY-MM-DD
    updated: YYYY-MM-DD
    confidence: high | medium | low
    ---

## Index Rules
- _master_index.md: max 80 lines, one line per topic cluster
- Topic indexes: max 60 lines, one line per article in the cluster
- Split threshold: 80 articles per topic index
- Every article in exactly one topic index
- Every topic index listed in _master_index.md

## Workflows
- **triage**: Sort inbox → capture, append only new items to _inbox_queue.md
- **ingest**: Read source → duplicate check → create/update wiki pages → update indexes → write daily log → move source to processed/{today}/ → update manifest → update queue
- **query**: Read master index → topic index → articles → synthesize → optionally file answer
- **lint**: Check contradictions, orphans, missing pages, stale claims, index integrity, broken links
- **resolve**: Walk through open contradictions, propose and apply resolutions with user confirmation
- **close**: Finalize daily log, verify index counts, refresh overview.md, suggest git commit
- **rollup**: Monthly summary with volume, themes, integrity, recommendations
```

---

## Next Steps (shown after each sub-command)

**After init:**
```
Next steps:
  Drop files into <path>/raw/inbox/
  do work bkb triage            Sort inbox items
  do work bkb status            Check KB state
```

**After triage:**
```
Next steps:
  do work bkb ingest            Compile ready sources into wiki
  do work bkb ingest <file>     Ingest a specific source
```

**After ingest:**
```
Next steps:
  do work bkb query [question]  Ask the wiki a question
  do work bkb lint              Health check after ingestion
  do work bkb close             Finalize the day
```

**After query:**
```
Next steps:
  do work bkb ingest            Add more sources
  do work bkb lint              Health check
```

**After lint:**
```
Next steps:
  do work bkb resolve           Resolve flagged contradictions
  do work bkb ingest            Address gaps with new sources
  do work bkb close             Finalize the day
```

**After resolve:**
```
Next steps:
  do work bkb lint              Verify fixes
  do work bkb close             Finalize the day
```

**After close:**
```
Next steps:
  do work bkb rollup            Generate monthly summary (if end of month)
  do work bkb status            Review KB state
```

**After rollup:**
```
Next steps:
  do work bkb lint full         Full integrity check
  do work bkb status            Review KB state
```
