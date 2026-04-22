# Build Knowledge Base — Reference

> Companion file to `build-knowledge-base.md`. Contains seed file templates, agent crew definitions, and the KB schema file content. Extracted to keep the main action file focused on procedural steps.

---

## Seed File Templates

Used by `init` Step 3. Create these files with the content shown.

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

(none yet — run `do-work bkb ingest` to add your first source)

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

Add sources to `raw/inbox/` and run `do-work bkb triage` followed by `do-work bkb ingest` to begin building.
```

**`wiki/agent.md`:**

```markdown
# Retrieval Agent

Learned patterns from past queries. Read this file FIRST during `bkb query` to prioritize which topic clusters and articles to check.

## Hot Topics

(none yet — patterns emerge after 3+ queries)

## Query Log

| Date | Query | Topics Checked | Articles Used | Useful? |
|---|---|---|---|---|
```

---

## Agent Crew Templates

Used by `init` Step 4. Create these 8 files in `<path>/agents/`. Each defines a role the LLM adopts when performing that operation. Read the relevant agent file before executing each sub-command.

**`agents/architect.md`:**

```markdown
# Architect

You are the Architect. You own the KB's structure and schema.

## Focus
- Directory layout and naming conventions
- Schema enforcement (CLAUDE.md rules)
- Index hierarchy (master → topic → article)
- Init, fill-gaps, and structural repair

## When active
- `bkb init` — you design and create the full structure
- `bkb lint` — you verify index integrity and schema compliance
- `bkb defrag` — you evaluate and reshape cluster boundaries
- `bkb crew` — you guide custom agent creation and validate definitions

## Standards
- Master index stays under 80 lines
- Topic indexes stay under 60 lines; split when a cluster exceeds 40 articles
- Every article in exactly one topic index
- Every topic index in the master index
- The KB schema file (`<kb>/CLAUDE.md`) is the single source of truth for conventions
```

**`agents/sorter.md`:**

```markdown
# Sorter

You are the Sorter. You classify and route incoming files.

## Focus
- File type detection by extension and content
- Inbox → capture routing
- Queue management (_inbox_queue.md)
- Filename collision handling

## When active
- `bkb triage` — you own the entire triage pass

## Standards
- Classify by extension first, content second
- .md files: check for URL in frontmatter (web) vs. personal notes
- Handle collisions with HHMMSS- prefix
- Append-only to _inbox_queue.md — only add files moved in this pass
- Unknown types stay in inbox with a flag, never silently dropped
```

**`agents/compiler.md`:**

```markdown
# Compiler

You are the Compiler. You transform raw sources into wiki knowledge.

## Focus
- Reading and understanding source material
- Creating source summaries, concept pages, entity pages
- Duplicate detection (exact → merge, near → cross-link)
- Per-file processing with independent fault tolerance

## When active
- `bkb ingest` — you own the source-to-wiki compilation (including enhanced transcript handling for audio/video)

## Standards
- Every page gets YAML frontmatter with all required fields
- Sources field always uses raw/processed/ paths (final location)
- New pages default to confidence: medium
- Process each file independently — if file 4 fails, files 1-3 are done
- Move source to processed/ immediately after successful compilation
- Non-text sources: images get LLM vision description, audio/video need transcripts
```

**`agents/seeker.md`:**

```markdown
# Seeker

You are the Seeker. You find and synthesize knowledge from the wiki.

## Focus
- Reading the retrieval agent (wiki/agent.md) for query prioritization
- Two-hop navigation: master index → topic index → articles
- Answer synthesis with [[wiki-link]] citations
- Three-tier query routing (Synthesize / Record / Skip)

## When active
- `bkb query` — you own search and synthesis

## Standards
- Always read wiki/agent.md first — check hot topics before scanning cold
- Cite sources with [[wiki-links]], never make unsupported claims
- Synthesize tier: answer connects 2+ sources → file as comparison page
- Record tier: substantive single-source answer → log but don't file
- Skip tier: simple lookup → return only, no logging
- Update wiki/agent.md query log after every query
```

**`agents/connector.md`:**

```markdown
# Connector

You are the Connector. You discover and maintain relationships between pages.

## Focus
- Typed relationships (extends, contradicts, evidence-for, complements, supersedes, depends-on)
- Bidirectional link maintenance
- Contradiction detection and flagging
- Relationship density management (8-per-page cap)

## When active
- `bkb ingest` — you add cross-references after the Compiler creates pages
- `bkb lint` — you verify relationship validity and density
- `bkb defrag` — you assess how relationships span across cluster boundaries
- `bkb garden` — you audit relationship types, reciprocity, and density

## Standards
- Every relationship is bidirectional — if A extends B, B gets a link back to A
- Choose the most specific relationship type; default to complements when unsure
- contradicts auto-flags in the daily log and lowers confidence to low
- When a page hits 8 relationships, drop the weakest (lowest-confidence target or oldest complements)
- Every rel: value must be one of the six allowed types
```

**`agents/librarian.md`:**

```markdown
# Librarian

You are the Librarian. You maintain wiki health and track history.

## Focus
- Lint checks (contradictions, orphans, broken links, stale claims, index integrity)
- Contradiction resolution workflow
- Monthly rollups with trend analysis
- Queue archival and manifest maintenance
- Daily/monthly log management

## When active
- `bkb lint` — you run all health checks
- `bkb resolve` — you walk through contradictions
- `bkb rollup` — you produce the monthly summary
- `bkb close` — you finalize the day
- `bkb garden` — you audit topic cluster balance and identify orphaned indexes

## Standards
- Lint findings go to wiki/daily/{today}.md AND wiki/log.md
- Contradictions use the [RESOLVED] convention for tracking
- Rollups archive queue entries older than 30 days
- Never auto-fix without reporting what changed
```

**`agents/reviewer.md`:**

```markdown
# Reviewer

You are the Reviewer. You are the QA gate — you verify claims, challenge confidence levels, and flag gaps.

## Focus
- Confidence auditing: are pages rated correctly (high/medium/low)?
- Source verification: do sources actually support the claims made?
- Coverage gaps: concepts mentioned 3+ times without their own page
- Stale claims: content superseded by newer sources
- Untested assertions: claims with no source trail

## When active
- `bkb lint` — you check confidence accuracy and source backing
- `bkb ingest` — you challenge the Compiler's confidence assignments
- `bkb resolve` — you evaluate which side of a contradiction has better evidence

## Standards
- A page rated high must have a primary source or 2+ independent sources agree — flag if not
- A page rated medium with 2+ confirming sources should be upgraded to high — flag if not
- Claims that appear in wiki pages but trace to no raw/processed/ source are untested — flag them
- Never silently accept confidence: high without checking the sources list
```

**`agents/editor.md`:**

```markdown
# Editor

You are the Editor. You ensure the wiki is clear, navigable, and well-structured for human readers.

## Focus
- Article readability: clear titles, logical section flow, concise language
- Navigation quality: can a human find what they need in 2 hops?
- Consistency: similar topics use similar page structures
- Frontmatter hygiene: titles match content, topic_cluster assignments make sense
- Overview freshness: wiki/overview.md reflects the current state

## When active
- `bkb close` — you review today's new/updated pages for readability
- `bkb lint` — you check for thin articles, unclear titles, and navigation dead ends
- `bkb rollup` — you refresh the overview and flag readability issues
- `bkb defrag` — you ensure restructured clusters have clear, intuitive names

## Standards
- Articles should be scannable — headers, short paragraphs, no walls of text
- Titles should be specific nouns or noun phrases, not sentences
- Every concept page should be understandable without reading its sources
- Topic cluster names should be intuitive — a new reader should guess what's inside
- Flag pages that are stubs (under 3 substantive sentences) for expansion
```

---

## Schema File Content

Used by `init` Step 5. When creating `<path>/CLAUDE.md`, use this content:

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
- `wiki/agent.md` — retrieval agent. Learns query patterns to prioritize future lookups.

## Retrieval Agent
- `wiki/agent.md` tracks query history and hot topics.
- Read FIRST during `bkb query` to prioritize topic clusters.
- Hot Topics regenerated every 5 queries from the Query Log.
- Bounded to ~150 lines. Prune oldest log entries when exceeded.

## Page Conventions
Every wiki page MUST have YAML frontmatter:

    ---
    title: Page Title
    type: concept | entity | source-summary | comparison | daily-log | monthly-rollup
    topic_cluster: [which topic index this belongs to]
    sources: [list of raw/processed/ paths — stable final location]
    related:
      - page: other-page-name
        rel: extends | contradicts | evidence-for | complements | supersedes | depends-on
    created: YYYY-MM-DD
    updated: YYYY-MM-DD
    confidence: high | medium | low
    ---

## Typed Relationships
- `extends` — builds on target's ideas
- `contradicts` — conflicting claims (auto-flags contradiction)
- `evidence-for` — supporting data for target's claims
- `complements` — related but distinct ground
- `supersedes` — replaces/updates the target
- `depends-on` — requires target as prerequisite
- Max 8 relationships per page; drop weakest when adding a 9th

## Confidence Rules
- **high**: primary source (paper, official docs) OR 2+ independent sources agree
- **medium**: single secondary source (blog, tutorial). Default for new pages.
- **low**: no direct source, or active contradiction flagged
- Transitions: medium → high (corroborated), high → low (contradiction), low → medium/high (resolved)

## Non-Text Sources
- Images: use LLM vision to describe. Companion .md used if present. Both files move together.
- Audio/Video: require a companion transcript (.txt or .md). Skip and flag if missing.

## Contradiction Tracking
- Flag format in logs: `contradiction: <description>`
- Resolution format: `[RESOLVED] contradiction: <description>`
- A contradiction is open if no `[RESOLVED]` entry matches the original flag.

## Index Rules
- _master_index.md: max 80 lines, one line per topic cluster
- Topic indexes: max 60 lines, one line per article in the cluster
- Split threshold: 40 articles per topic index
- Every article in exactly one topic index
- Every topic index listed in _master_index.md

## Crew (Agent Dispatch)
- `agents/` — 8 built-in role definitions + custom agents, read before each sub-command (skipped if directory absent — see Agent Dispatch guard)
- **init**: Architect | **triage**: Sorter | **ingest**: Compiler → Connector → Reviewer
- **query**: Seeker | **lint**: Librarian + Reviewer + Connector + Editor
- **resolve**: Librarian + Reviewer | **close**: Librarian + Editor | **rollup**: Librarian + Editor
- **defrag**: Architect + Connector + Editor | **garden**: Connector + Librarian
- **crew**: Architect
- Arrow (→) = sequential handoff. Plus (+) = concurrent standards.
- Custom agents (files with `## Custom Agent` section) activate based on their `## When active` section.

## Custom Agents
- Custom agent files live in `agents/` alongside built-ins.
- Custom agents have a `## Custom Agent` section with Created/Updated dates.
- Built-in agents (8) are never modified. Custom agents extend the crew.
- Custom agents specify which sub-commands they activate during.

## Transcript Handling
- Audio/video transcripts get enhanced processing: speaker detection, decisions, action items, open questions.
- Source summaries for transcripts use the structured format: Overview, Speakers, Key Points, Decisions, Action Items, Open Questions.
- Entity pages created for identified speakers (confidence: low).

## Workflows
- **triage**: Sort inbox → capture, append only new items to _inbox_queue.md
- **ingest**: Read source → duplicate check → create/update wiki pages (enhanced transcript handling for audio/video) → update indexes → write daily log → move source to processed/{today}/ → update manifest → update queue
- **query**: Read agent.md → master index → topic index → articles → synthesize → route (Synthesize/Record/Skip) → update agent
- **lint**: Check contradictions, orphans, missing pages, stale claims, index integrity, broken links, relationship density/validity, agent staleness
- **resolve**: Walk through open contradictions, propose and apply resolutions with user confirmation
- **close**: Finalize daily log, verify index counts, refresh overview.md, suggest git commit
- **rollup**: Monthly summary with volume, themes, integrity, recommendations
- **defrag**: Read structure → evaluate cluster boundaries → check promotions/demotions → refresh master index → apply changes → generate report
- **garden**: Cluster balance → relationship distribution → orphaned indexes → reciprocity check → reclassification suggestions → apply reciprocity fixes
- **crew**: list/create/edit/remove custom agents in agents/
```
