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
