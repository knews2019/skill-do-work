# BKB — Build Knowledge Base

An LLM-compiled Markdown wiki built from raw sources. Instead of re-deriving knowledge on every query, the LLM incrementally compiles sources into a structured, interlinked wiki. The raw sources are the input, the LLM is the compiler, the wiki is the output.

Aliases: `bkb`, `kb`, `build knowledge base`, `knowledge base`.

## Folder structure

After `do work bkb init`, the KB looks like this:

```
kb/
├── raw/                          # Source document pipeline
│   ├── inbox/                    # Drop zone — put anything here
│   │   └── clippings/            # Browser clipper target (auto-classified as web)
│   ├── capture/                  # Type-sorted staging (transient)
│   │   ├── web/                  #   Web articles
│   │   ├── papers/               #   PDFs, academic papers
│   │   ├── repos/                #   Code files, READMEs
│   │   ├── images/               #   PNG, JPG, SVG
│   │   ├── notes/                #   Personal notes
│   │   ├── audio/                #   MP3, WAV, M4A
│   │   └── video/                #   MP4, WEBM, MKV + transcripts
│   ├── processed/                # Ingested sources archived by date
│   │   ├── 2026-04-05/
│   │   └── _manifest.md          # Permanent audit trail
│   └── _inbox_queue.md           # Append-only triage ledger
│
├── wiki/                         # LLM-maintained wiki
│   ├── _master_index.md          # Top-level nav (≤80 lines)
│   ├── topics/                   # Second-level cluster indexes
│   │   └── _index_[topic].md    #   One per topic (≤80 articles each)
│   ├── concepts/                 # Concept articles
│   ├── entities/                 # Person/org/product articles
│   ├── sources/                  # Source summaries
│   ├── comparisons/              # Filed query outputs
│   ├── daily/                    # Daily changelogs (YYYY-MM-DD.md)
│   ├── monthly/                  # Monthly rollups (YYYY-MM.md)
│   ├── log.md                    # Append-only activity timeline
│   ├── overview.md               # High-level synthesis
│   └── agent.md                  # Retrieval agent — learns query patterns
│
└── CLAUDE.md                     # Schema — conventions, frontmatter, workflows
```

## Journey of a file

A source file moves through four stages:

```
inbox/ ──triage──> capture/ ──ingest──> processed/
                                  │
                                  └──> wiki pages created/updated
```

### Stage 1: Inbox

Drop files into `raw/inbox/`. PDFs, articles, notes, images, audio, video — anything. No sorting needed.

### Stage 2: Triage

Run `do work bkb triage`.

- Scans `raw/inbox/` for files
- Classifies each by extension (`.pdf` → `papers/`, `.png` → `images/`, etc.)
- Moves files from `inbox/` to the matching `capture/` subdirectory
- Appends entries to `raw/_inbox_queue.md` with status "ready"
- Handles filename collisions by prefixing with `HHMMSS-`

After triage, inbox is empty and capture holds sorted files waiting to be ingested.

### Stage 3: Ingest

Run `do work bkb ingest` (all ready items) or `do work bkb ingest <filename>` (one file).

For each file:

1. **Read** the source from `raw/capture/`
2. **Handle non-text sources** — images get an LLM-generated description; audio/video require a companion transcript
3. **Check for duplicates** — exact duplicates update existing pages; near-duplicates get cross-linked
4. **Create or update wiki pages** — source summaries, concept pages, entity pages as needed
5. **Flag contradictions** with existing wiki content
6. **Update indexes** — topic indexes, master index, cross-references via `[[wiki-links]]`
7. **Move source** from `capture/` to `raw/processed/{today}/` and mark "done" in the queue
8. **Log the work** in `wiki/daily/{today}.md`

Each file is processed independently — if file 4 of 5 fails, files 1-3 are already done.

### Stage 4: Processed

Final resting place. Files are archived by ingestion date (`raw/processed/2026-04-05/`). The `_manifest.md` records every file, its processing date, and which wiki articles it produced. Files here are never re-triaged.

## Wiki page format

Every wiki page has YAML frontmatter:

```yaml
---
title: Page Title
type: concept | entity | source-summary | comparison
topic_cluster: [which topic index this belongs to]
sources: [raw/processed/ paths — always the final location]
related:
  - page: other-page-name
    rel: extends | contradicts | evidence-for | complements | supersedes | depends-on
created: YYYY-MM-DD
updated: YYYY-MM-DD
confidence: high | medium | low
---
```

**Typed relationships** describe *how* pages connect, not just *that* they connect. Six types: `extends`, `contradicts`, `evidence-for`, `complements`, `supersedes`, `depends-on`. Max 8 per page.

**Confidence levels:**

| Level | Meaning |
|-------|---------|
| high | Primary source (peer-reviewed, official docs) or 2+ sources agree |
| medium | Single secondary source — default for new pages |
| low | No direct source (inferred) or active contradiction flagged |

## Navigation

Three-layer hierarchy that scales to thousands of articles:

1. **`_master_index.md`** — lists all topic clusters with article counts
2. **`topics/_index_[topic].md`** — lists all articles in one cluster (splits at 80 articles)
3. **Individual pages** — concepts, entities, sources, comparisons

Any article is reachable in two hops from the master index.

## Commands

| Command | What it does |
|---------|--------------|
| `bkb init [path]` | Create the full directory structure and seed files |
| `bkb triage` | Sort inbox + clippings into capture directories |
| `bkb ingest [file]` | Compile sources into wiki pages |
| `bkb query [question]` | Search the wiki and synthesize an answer |
| `bkb lint` | Quick health check (orphans, broken links, stale claims) |
| `bkb lint full` | Full structural check |
| `bkb resolve` | Walk through open contradictions |
| `bkb close` | Finalize daily log, refresh overview, suggest commit |
| `bkb rollup` | Monthly summary |
| `bkb status` | KB stats and pending items |

## Retrieval agent

The file `wiki/agent.md` learns from your queries over time. It tracks which topic clusters and articles get used most often, so future queries check the most relevant areas first instead of scanning cold.

- Activates after 3+ queries (cold start threshold)
- Hot Topics section regenerated every 5 queries
- Bounded to ~150 lines — oldest log entries pruned automatically

## Query routing

Queries are classified into three tiers to prevent wiki bloat:

| Tier | When | What happens |
|------|------|--------------|
| **Synthesize** | Answer connects 2+ sources | Filed as a wiki page in `comparisons/`, indexes updated |
| **Record** | Substantive answer, single source | Returned to user, brief log entry, no wiki page |
| **Skip** | Simple factual lookup | Returned to user, nothing logged |

## Capture shortcuts

- **General drop zone**: `kb/raw/inbox/` — drop any file
- **Browser clippings**: `kb/raw/inbox/clippings/` — point your browser extension here (Obsidian Web Clipper, MarkDownload, etc.). Files here are auto-classified as web content during triage.

## Typical workflow

```
do work bkb init                  # one-time setup
# drop files into kb/raw/inbox/ (or clip via browser extension)
do work bkb triage                # sort them
do work bkb ingest                # compile into wiki
do work bkb query [question]      # ask the wiki anything
do work bkb lint                  # check for issues
do work bkb close                 # wrap up the day
```
