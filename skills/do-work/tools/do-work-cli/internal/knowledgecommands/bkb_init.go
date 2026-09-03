package knowledgecommands

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type scaffoldFile struct{ path, content string }

var scaffoldDirectories = []string{
	"raw/inbox", "raw/capture/web", "raw/capture/papers", "raw/capture/repos", "raw/capture/images", "raw/capture/notes", "raw/capture/audio", "raw/capture/video", "raw/processed",
	"wiki/topics", "wiki/concepts", "wiki/entities", "wiki/sources", "wiki/comparisons", "wiki/daily", "wiki/monthly", "agents",
}

var scaffoldFiles = []scaffoldFile{
	{"raw/_inbox_queue.md", "# Inbox Queue\n\nItems pending triage. Updated automatically during triage.\n\n| # | File | Source Type | Status |\n|---|---|---|---|\n"},
	{"raw/processed/_manifest.md", "# Processing Manifest\n\n| File | Date Processed | Processed Path | Wiki Articles Produced | Status |\n|---|---|---|---|---|\n"},
	{"wiki/_master_index.md", "# Master Index\n\nLast updated: {today} | Total articles: 0 | Topic clusters: 0\n\n## Topic Clusters\n\n(none yet — run `do-work-knowledge bkb ingest` to add your first source)\n\n## Recent Activity\n\n- {today}: Knowledge base initialized\n"},
	{"wiki/log.md", "# Activity Log\n\n## [{today}] init | Knowledge base created\n\nStructure initialized. Ready for first source.\n"},
	{"wiki/overview.md", "# Knowledge Base Overview\n\nThis knowledge base was initialized on {today}. No sources have been ingested yet.\n\nAdd sources to `raw/inbox/` and run `do-work-knowledge bkb triage` followed by `do-work-knowledge bkb ingest` to begin building.\n"},
	{"wiki/agent.md", "# Retrieval Agent\n\nLearned patterns from past queries. Read this file FIRST during `bkb query` to prioritize which topic clusters and articles to check.\n\n## Hot Topics\n\n(none yet — patterns emerge after 3+ queries)\n\n## Query Log\n\n| Date | Query | Topics Checked | Articles Used | Useful? |\n|---|---|---|---|---|\n"},
	{"agents/architect.md", architectTemplate}, {"agents/sorter.md", sorterTemplate}, {"agents/compiler.md", compilerTemplate}, {"agents/seeker.md", seekerTemplate},
	{"agents/connector.md", connectorTemplate}, {"agents/librarian.md", librarianTemplate}, {"agents/reviewer.md", reviewerTemplate}, {"agents/editor.md", editorTemplate},
	{"CLAUDE.md", schemaTemplate},
}

const architectTemplate = `# Architect

You are the Architect. You own the KB's structure and schema.

## Focus
- Directory layout and naming conventions
- Schema enforcement (CLAUDE.md rules)
- Index hierarchy (master → topic → article)
- Init, fill-gaps, and structural repair

## When active
- ` + "`bkb init`" + ` — you design and create the full structure
- ` + "`bkb lint`" + ` — you verify index integrity and schema compliance
- ` + "`bkb defrag`" + ` — you evaluate and reshape cluster boundaries
- ` + "`bkb crew`" + ` — you guide custom agent creation and validate definitions

## Standards
- Master index stays under 80 lines
- Topic indexes stay under 60 lines; split when a cluster exceeds 40 articles
- Every article in exactly one topic index
- Every topic index in the master index
- The KB schema file (` + "`<kb>/CLAUDE.md`" + `) is the single source of truth for conventions
`

const sorterTemplate = `# Sorter

You are the Sorter. You classify and route incoming files.

## Focus
- File type detection by extension and content
- Inbox → capture routing
- Queue management (_inbox_queue.md)
- Filename collision handling

## When active
- ` + "`bkb triage`" + ` — you own the entire triage pass

## Standards
- Classify by extension first, content second
- .md files: check for URL in frontmatter (web) vs. personal notes
- Handle collisions with HHMMSS- prefix
- Append-only to _inbox_queue.md — only add files moved in this pass
- Unknown types stay in inbox with a flag, never silently dropped
`

const compilerTemplate = `# Compiler

You are the Compiler. You transform raw sources into wiki knowledge.

## Focus
- Reading and understanding source material
- Creating source summaries, concept pages, entity pages
- Duplicate detection (exact → merge, near → cross-link)
- Per-file processing with independent fault tolerance

## When active
- ` + "`bkb ingest`" + ` — you own the source-to-wiki compilation (including enhanced transcript handling for audio/video)

## Standards
- Every page gets YAML frontmatter with all required fields
- Sources field always uses raw/processed/ paths (final location)
- New pages default to confidence: medium
- Process each file independently — if file 4 fails, files 1-3 are done
- Move source to processed/ immediately after successful compilation
- Non-text sources: images get LLM vision description, audio/video need transcripts
`

const seekerTemplate = `# Seeker

You are the Seeker. You find and synthesize knowledge from the wiki.

## Focus
- Reading the retrieval agent (wiki/agent.md) for query prioritization
- Two-hop navigation: master index → topic index → articles
- Answer synthesis with [[wiki-link]] citations
- Three-tier query routing (Synthesize / Record / Skip)

## When active
- ` + "`bkb query`" + ` — you own search and synthesis

## Standards
- Always read wiki/agent.md first — check hot topics before scanning cold
- Cite sources with [[wiki-links]], never make unsupported claims
- Synthesize tier: answer connects 2+ sources → file as comparison page
- Record tier: substantive single-source answer → log but don't file
- Skip tier: simple lookup → return only, no logging
- Update wiki/agent.md query log after every query
`

const connectorTemplate = `# Connector

You are the Connector. You discover and maintain relationships between pages.

## Focus
- Typed relationships (extends, contradicts, evidence-for, complements, supersedes, depends-on)
- Bidirectional link maintenance
- Contradiction detection and flagging
- Relationship density management (8-per-page cap)

## When active
- ` + "`bkb ingest`" + ` — you add cross-references after the Compiler creates pages
- ` + "`bkb lint`" + ` — you verify relationship validity and density
- ` + "`bkb defrag`" + ` — you assess how relationships span across cluster boundaries
- ` + "`bkb garden`" + ` — you audit relationship types, reciprocity, and density

## Standards
- Every relationship is bidirectional — if A extends B, B gets a link back to A
- Choose the most specific relationship type; default to complements when unsure
- contradicts auto-flags in the daily log and lowers confidence to low
- When a page hits 8 relationships, drop the weakest (lowest-confidence target or oldest complements)
- Every rel: value must be one of the six allowed types
`

const librarianTemplate = `# Librarian

You are the Librarian. You maintain wiki health and track history.

## Focus
- Lint checks (contradictions, orphans, broken links, stale claims, index integrity)
- Contradiction resolution workflow
- Monthly rollups with trend analysis
- Queue archival and manifest maintenance
- Daily/monthly log management

## When active
- ` + "`bkb lint`" + ` — you run all health checks
- ` + "`bkb resolve`" + ` — you walk through contradictions
- ` + "`bkb rollup`" + ` — you produce the monthly summary
- ` + "`bkb close`" + ` — you finalize the day
- ` + "`bkb garden`" + ` — you audit topic cluster balance and identify orphaned indexes

## Standards
- Lint findings go to wiki/daily/{today}.md AND wiki/log.md
- Contradictions use the [RESOLVED] convention for tracking
- Rollups archive queue entries older than 30 days
- Never auto-fix without reporting what changed
`

const reviewerTemplate = `# Reviewer

You are the Reviewer. You are the QA gate — you verify claims, challenge confidence levels, and flag gaps.

## Focus
- Confidence auditing: are pages rated correctly (high/medium/low)?
- Source verification: do sources actually support the claims made?
- Coverage gaps: concepts mentioned 3+ times without their own page
- Stale claims: content superseded by newer sources
- Untested assertions: claims with no source trail

## When active
- ` + "`bkb lint`" + ` — you check confidence accuracy and source backing
- ` + "`bkb ingest`" + ` — you challenge the Compiler's confidence assignments
- ` + "`bkb resolve`" + ` — you evaluate which side of a contradiction has better evidence

## Standards
- A page rated high must have a primary source or 2+ independent sources agree — flag if not
- A page rated medium with 2+ confirming sources should be upgraded to high — flag if not
- Claims that appear in wiki pages but trace to no raw/processed/ source are untested — flag them
- Never silently accept confidence: high without checking the sources list
`

const editorTemplate = `# Editor

You are the Editor. You ensure the wiki is clear, navigable, and well-structured for human readers.

## Focus
- Article readability: clear titles, logical section flow, concise language
- Navigation quality: can a human find what they need in 2 hops?
- Consistency: similar topics use similar page structures
- Frontmatter hygiene: titles match content, topic_cluster assignments make sense
- Overview freshness: wiki/overview.md reflects the current state

## When active
- ` + "`bkb close`" + ` — you review today's new/updated pages for readability
- ` + "`bkb lint`" + ` — you check for thin articles, unclear titles, and navigation dead ends
- ` + "`bkb rollup`" + ` — you refresh the overview and flag readability issues
- ` + "`bkb defrag`" + ` — you ensure restructured clusters have clear, intuitive names

## Standards
- Articles should be scannable — headers, short paragraphs, no walls of text
- Titles should be specific nouns or noun phrases, not sentences
- Every concept page should be understandable without reading its sources
- Topic cluster names should be intuitive — a new reader should guess what's inside
- Flag pages that are stubs (under 3 substantive sentences) for expansion
`

const schemaTemplate = `# LLM Knowledge Base Schema

## Project Structure
- ` + "`raw/`" + ` — source documents with lifecycle pipeline. NEVER modify originals.
- ` + "`raw/inbox/`" + ` — zero-friction drop zone. Sort into capture/ before processing.
- ` + "`raw/capture/`" + ` — type-sorted staging area.
- ` + "`raw/processed/YYYY-MM-DD/`" + ` — ingested sources, moved here after successful compilation.
- ` + "`raw/_inbox_queue.md`" + ` — append-only triage ledger. Only updated with files moved in the current triage pass.
- ` + "`wiki/`" + ` — LLM-generated wiki. You own this entirely.
- ` + "`wiki/_master_index.md`" + ` — top-level catalog. Read FIRST on every query.
- ` + "`wiki/topics/_index_[topic].md`" + ` — second-level indexes by topic cluster.
- ` + "`wiki/daily/YYYY-MM-DD.md`" + ` — daily changelog.
- ` + "`wiki/monthly/YYYY-MM.md`" + ` — monthly rollup and trends.
- ` + "`wiki/log.md`" + ` — append-only activity log.
- ` + "`wiki/agent.md`" + ` — retrieval agent. Learns query patterns to prioritize future lookups.

## Retrieval Agent
- ` + "`wiki/agent.md`" + ` tracks query history and hot topics.
- Read FIRST during ` + "`bkb query`" + ` to prioritize topic clusters.
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
- ` + "`extends`" + ` — builds on target's ideas
- ` + "`contradicts`" + ` — conflicting claims (auto-flags contradiction)
- ` + "`evidence-for`" + ` — supporting data for target's claims
- ` + "`complements`" + ` — related but distinct ground
- ` + "`supersedes`" + ` — replaces/updates the target
- ` + "`depends-on`" + ` — requires target as prerequisite
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
- Flag format in logs: ` + "`contradiction: <description>`" + `
- Resolution format: ` + "`[RESOLVED] contradiction: <description>`" + `
- A contradiction is open if no ` + "`[RESOLVED]`" + ` entry matches the original flag.

## Index Rules
- _master_index.md: max 80 lines, one line per topic cluster
- Topic indexes: max 60 lines, one line per article in the cluster
- Split threshold: 40 articles per topic index
- Every article in exactly one topic index
- Every topic index listed in _master_index.md

## Crew (Agent Dispatch)
- ` + "`agents/`" + ` — 8 built-in role definitions + custom agents, read before each sub-command (skipped if directory absent — see Agent Dispatch guard)
- **init**: Architect | **triage**: Sorter | **ingest**: Compiler → Connector → Reviewer
- **query**: Seeker | **lint**: Librarian + Reviewer + Connector + Editor
- **resolve**: Librarian + Reviewer | **close**: Librarian + Editor | **rollup**: Librarian + Editor
- **defrag**: Architect + Connector + Editor | **garden**: Connector + Librarian
- **crew**: Architect
- Arrow (→) = sequential handoff. Plus (+) = concurrent standards.
- Custom agents (files with ` + "`## Custom Agent`" + ` section) activate based on their ` + "`## When active`" + ` section.

## Custom Agents
- Custom agent files live in ` + "`agents/`" + ` alongside built-ins.
- Custom agents have a ` + "`## Custom Agent`" + ` section with Created/Updated dates.
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
`

func renderScaffold(content, today string) string {
	return strings.ReplaceAll(content, "{today}", today)
}

var knowledgeMutationHook = func(_, _, _ string) error { return nil }

// rootedObject is one filesystem object this invocation created. identity is nil when the
// object escaped identity capture, which rollback reports instead of removing by pathname.
// digestKnown marks a completely written file, whose content digest separates it from a
// replacement that reused the same name and inode.
type rootedObject struct {
	path        string
	identity    os.FileInfo
	digest      [sha256.Size]byte
	digestKnown bool
	recursive   bool
}

type rootedScaffoldWriter struct {
	root        *os.Root
	flow        string
	directories map[string]os.FileInfo
	created     []rootedObject
}

func newRootedScaffoldWriter(rootPath, flow string) (*rootedScaffoldWriter, error) {
	root, err := os.OpenRoot(rootPath)
	if err != nil {
		return nil, err
	}
	identity, err := root.Stat(".")
	if err != nil {
		_ = root.Close()
		return nil, err
	}
	return &rootedScaffoldWriter{root: root, flow: flow, directories: map[string]os.FileInfo{".": identity}}, nil
}

func (writer *rootedScaffoldWriter) close() { _ = writer.root.Close() }

func (writer *rootedScaffoldWriter) snapshotDirectories(paths []string) error {
	for _, path := range paths {
		path = cleanRootPath(path)
		info, err := writer.root.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("scaffold parent %q is not a real directory", path)
		}
		writer.directories[path] = info
	}
	return nil
}

func (writer *rootedScaffoldWriter) validateParents(path string) error {
	parent := cleanRootPath(filepath.Dir(path))
	for knownPath, identity := range writer.directories {
		if !rootPathContains(parent, knownPath) {
			continue
		}
		var current os.FileInfo
		var err error
		if knownPath == "." {
			current, err = writer.root.Stat(knownPath)
		} else {
			current, err = writer.root.Lstat(knownPath)
		}
		if err != nil {
			return fmt.Errorf("revalidate scaffold parent %q: %w", knownPath, err)
		}
		if !os.SameFile(identity, current) {
			return fmt.Errorf("scaffold parent %q changed after validation", knownPath)
		}
	}
	return nil
}

func (writer *rootedScaffoldWriter) createDirectory(path string, recursive bool) error {
	path = cleanRootPath(path)
	if err := writer.validateParents(path); err != nil {
		return err
	}
	if err := knowledgeMutationHook("before-create", writer.flow, path); err != nil {
		return err
	}
	if err := writer.validateParents(path); err != nil {
		return err
	}
	if err := writer.root.Mkdir(path, 0o755); err != nil {
		return err
	}
	// Mkdir is the ownership event. Recording before the fallible Lstat keeps a directory
	// this invocation created from escaping rollback's inventory.
	created := len(writer.created)
	writer.created = append(writer.created, rootedObject{path: path, recursive: recursive})
	info, err := writer.root.Lstat(path)
	if err != nil {
		return err
	}
	writer.directories[path] = info
	writer.created[created].identity = info
	return knowledgeMutationHook("after-create", writer.flow, path)
}

func (writer *rootedScaffoldWriter) createFile(path, content string) error {
	path = cleanRootPath(path)
	if err := writer.validateParents(path); err != nil {
		return err
	}
	if err := knowledgeMutationHook("before-create", writer.flow, path); err != nil {
		return err
	}
	if err := writer.validateParents(path); err != nil {
		return err
	}
	fileHandle, err := writer.root.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	// The exclusive create is the ownership event, and the open handle names the object it
	// produced. Recording here keeps an incomplete write from escaping rollback's inventory.
	created := len(writer.created)
	writer.created = append(writer.created, rootedObject{path: path})
	if openedInfo, statError := fileHandle.Stat(); statError == nil {
		writer.created[created].identity = openedInfo
	}
	_, writeError := fileHandle.WriteString(content)
	closeError := fileHandle.Close()
	if writeError != nil {
		return writeError
	}
	if closeError != nil {
		return closeError
	}
	info, err := writer.root.Lstat(path)
	if err != nil {
		return err
	}
	writer.created[created].identity = info
	writer.created[created].digest = sha256.Sum256([]byte(content))
	writer.created[created].digestKnown = true
	return knowledgeMutationHook("after-create", writer.flow, path)
}

func (writer *rootedScaffoldWriter) rollback() resultmodel.RollbackResult {
	result := resultmodel.RollbackResult{Status: resultmodel.RollbackSucceeded, Actions: []string{}, Errors: []string{}}
	for index := len(writer.created) - 1; index >= 0; index-- {
		created := writer.created[index]
		_ = knowledgeMutationHook("before-rollback", writer.flow, created.path)
		if created.identity == nil {
			result.Errors = append(result.Errors, fmt.Sprintf("preserve %s: identity was not recorded", created.path))
			continue
		}
		current, err := writer.root.Lstat(created.path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("preserve %s: %v", created.path, err))
			continue
		}
		if !os.SameFile(created.identity, current) || !scaffoldContentStillOwned(writer.root, created) {
			result.Errors = append(result.Errors, fmt.Sprintf("preserve replacement at %s: identity changed", created.path))
			continue
		}
		if created.recursive {
			err = writer.root.RemoveAll(created.path)
		} else {
			err = writer.root.Remove(created.path)
		}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("remove owned %s: %v", created.path, err))
			continue
		}
		result.Actions = append(result.Actions, "removed owned "+created.path)
	}
	if len(result.Errors) > 0 {
		result.Status = resultmodel.RollbackIncomplete
	}
	return result
}

// scaffoldContentStillOwned separates a completely written file from a replacement that
// reused its name and inode, which removing and recreating a file in the same directory
// routinely does.
func scaffoldContentStillOwned(root *os.Root, created rootedObject) bool {
	if !created.digestKnown {
		return true
	}
	fileHandle, err := root.Open(created.path)
	if err != nil {
		return false
	}
	defer fileHandle.Close()
	openedInfo, statError := fileHandle.Stat()
	if statError != nil || !os.SameFile(created.identity, openedInfo) {
		return false
	}
	contents, readError := io.ReadAll(fileHandle)
	return readError == nil && sha256.Sum256(contents) == created.digest
}

func cleanRootPath(path string) string {
	clean := filepath.ToSlash(filepath.Clean(path))
	if clean == "" {
		return "."
	}
	return clean
}
func rootPathContains(path, candidate string) bool {
	path, candidate = cleanRootPath(path), cleanRootPath(candidate)
	return candidate == "." || path == candidate || strings.HasPrefix(path, candidate+"/")
}

func runGitInRoot(root *os.Root, rootPath string, arguments ...string) ([]byte, error) {
	command := exec.Command("git", arguments...)
	beforeIdentity, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}
	command.Dir = rootPath
	output, err := command.CombinedOutput()
	afterIdentity, identityError := os.Stat(rootPath)
	if identityError != nil || !os.SameFile(beforeIdentity, afterIdentity) {
		return output, fmt.Errorf("root changed during Git operation")
	}
	if err != nil {
		return output, fmt.Errorf("git %s: %s: %w", strings.Join(arguments, " "), strings.TrimSpace(string(output)), err)
	}
	return output, nil
}

func handleBKBInit(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseBKBOptions(arguments, true)
	if err != nil {
		return usageResult(CommandBKBInit, err)
	}
	targetPath, err := ensureSafeTarget(executionContext.RepositoryRoot, options.target)
	if err != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{knowledgeFinding(CommandBKBInit, "BKB-INIT-UNSAFE-TARGET", resultmodel.SeverityError, []string{options.target}, err.Error(), resultmodel.FixabilityRefused, "a symlink or unsafe target component prevents confined publication")}}
	}
	if err := validateScaffoldTarget(targetPath); err != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{knowledgeFinding(CommandBKBInit, "BKB-INIT-UNSAFE-COLLISION", resultmodel.SeverityError, []string{options.target}, err.Error(), resultmodel.FixabilityRefused, "a scaffold path has an unsafe type and cannot be preserved as a file or directory")}}
	}
	if !options.fillGaps && isKnowledgeBase(targetPath) {
		finding := knowledgeFinding(CommandBKBInit, "BKB-INIT-EXISTS", resultmodel.SeverityWarning, []string{options.target}, "knowledge base already contains raw/ and wiki/", resultmodel.FixabilityAutomatic, "existing content is never overwritten; use --fill-gaps")
		finding.NextArgv = []string{"do-work-cli", CommandBKBInit, "--kb", options.target, "--fill-gaps"}
		finding.NextJustRecipe = CommandBKBInit + " " + quoteRecipeArgument(options.target) + " --fill-gaps"
		finding.VerificationArgv = []string{"do-work-cli", "--format", "json", CommandBKBInit, "--kb", options.target, "--fill-gaps", "--dry-run"}
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{finding}}
	}
	if info, statError := os.Stat(targetPath); statError == nil && info.IsDir() {
		if gitRoot, gitError := enclosingGitRoot(targetPath); gitError == nil {
			return initializeInRepository(gitRoot, targetPath, options, nowUTC())
		}
	}
	if gitRoot, gitError := enclosingGitRoot(executionContext.RepositoryRoot); gitError == nil {
		if targetRelative, relativeError := filepath.Rel(gitRoot, targetPath); relativeError == nil && targetRelative != ".." && !strings.HasPrefix(targetRelative, ".."+string(filepath.Separator)) {
			return initializeInRepository(gitRoot, targetPath, options, nowUTC())
		}
	}
	if options.dryRun {
		return plannedScaffoldResult(targetPath, options.target, nowUTC())
	}
	return initializeStandaloneTarget(targetPath, options, nowUTC())
}

func plannedScaffoldResult(targetPath, targetRelative string, now time.Time) resultmodel.CommandResult {
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded}}
	for _, file := range scaffoldFiles {
		if _, err := os.Lstat(filepath.Join(targetPath, filepath.FromSlash(file.path))); os.IsNotExist(err) {
			result.Changes = append(result.Changes, resultmodel.RecordedChange{Path: filepath.ToSlash(filepath.Join(targetRelative, file.path)), Kind: "planned_create", Detail: "canonical BKB scaffold file"})
		} else {
			result.SkippedWork = append(result.SkippedWork, resultmodel.SkippedWork{Code: "BKB-INIT-PRESERVE-EXISTING", Reason: filepath.ToSlash(filepath.Join(targetRelative, file.path))})
		}
	}
	return result
}

func initializeInRepository(gitRoot, targetPath string, options bkbOptions, now time.Time) resultmodel.CommandResult {
	targetRelative, err := filepath.Rel(gitRoot, targetPath)
	if err != nil || targetRelative == ".." || strings.HasPrefix(targetRelative, ".."+string(filepath.Separator)) {
		return usageResult(CommandBKBInit, fmt.Errorf("target is outside the enclosing Git repository"))
	}
	missingFiles, existingFiles := scaffoldInventory(targetPath, filepath.ToSlash(targetRelative))
	missingDirectories := missingDirectoryInventory(targetPath, filepath.ToSlash(targetRelative))
	if len(missingFiles) == 0 {
		result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess}
		for _, path := range existingFiles {
			result.SkippedWork = append(result.SkippedWork, resultmodel.SkippedWork{Code: "BKB-INIT-PRESERVE-EXISTING", Reason: path})
		}
		return result
	}
	preflight := gittransaction.PreflightTargets(context.Background(), gitRoot, missingFiles, options.commit)
	if preflight.Failure != nil {
		outcome := resultmodel.OutcomeFailure
		if preflight.Failure.Kind == gittransaction.FailureDirtyIndex || preflight.Failure.Kind == gittransaction.FailureDirtyTarget {
			outcome = resultmodel.OutcomeRefused
		}
		return gittransaction.BuildCommandResult(CommandBKBInit, gittransaction.TransactionResult{Outcome: outcome, RepositoryRoot: preflight.RepositoryRoot, Failure: preflight.Failure, Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded}})
	}
	if options.dryRun {
		result := plannedScaffoldResult(targetPath, filepath.ToSlash(targetRelative), now)
		for _, path := range existingFiles {
			result.SkippedWork = append(result.SkippedWork, resultmodel.SkippedWork{Code: "BKB-INIT-PRESERVE-EXISTING", Reason: path})
		}
		return result
	}
	writer, openError := newRootedScaffoldWriter(gitRoot, "git")
	if openError != nil {
		return initializationFailure(options.target, openError, resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded})
	}
	defer writer.close()
	if err := writer.snapshotDirectories(scaffoldRootDirectories(filepath.ToSlash(targetRelative))); err != nil {
		return initializationFailure(options.target, err, resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded})
	}
	if err := knowledgeMutationHook("after-validation", writer.flow, filepath.ToSlash(targetRelative)); err != nil {
		return initializationFailure(options.target, err, writer.rollback())
	}
	if err := applyRootedScaffold(writer, filepath.ToSlash(targetRelative), missingDirectories, missingFiles, now); err != nil {
		return initializationFailure(options.target, err, writer.rollback())
	}
	if options.commit {
		if _, err := runGitInRoot(writer.root, gitRoot, append([]string{"add", "-A", "--"}, missingFiles...)...); err != nil {
			unstageRootedPaths(writer.root, gitRoot, missingFiles)
			return initializationFailure(options.target, err, writer.rollback())
		}
		if _, err := runGitInRoot(writer.root, gitRoot, "commit", "-m", "bkb: initialize knowledge base"); err != nil {
			unstageRootedPaths(writer.root, gitRoot, missingFiles)
			return initializationFailure(options.target, err, writer.rollback())
		}
		if err := verifyRootedCommit(writer.root, gitRoot, missingFiles); err != nil {
			return rootedCommittedRisk(writer.root, gitRoot, options.target, err)
		}
	}
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: gitRoot, Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded}}
	for _, path := range missingFiles {
		result.Changes = append(result.Changes, resultmodel.RecordedChange{Path: path, Kind: "created", Detail: "canonical BKB scaffold file"})
	}
	for _, path := range existingFiles {
		result.SkippedWork = append(result.SkippedWork, resultmodel.SkippedWork{Code: "BKB-INIT-PRESERVE-EXISTING", Reason: path})
	}
	return result
}

func scaffoldRootDirectories(prefix string) []string {
	paths := make([]string, 0, len(scaffoldDirectorySuffixes()))
	for _, suffix := range scaffoldDirectorySuffixes() {
		paths = append(paths, cleanRootPath(filepath.Join(prefix, suffix)))
	}
	return paths
}

func applyRootedScaffold(writer *rootedScaffoldWriter, prefix string, missingDirectories, missingFiles []string, now time.Time) error {
	for _, directory := range missingDirectories {
		if err := writer.createDirectory(directory, false); err != nil {
			return err
		}
	}
	for _, file := range scaffoldFiles {
		path := cleanRootPath(filepath.Join(prefix, file.path))
		if containsString(missingFiles, path) {
			if err := writer.createFile(path, renderScaffold(file.content, now.Format("2006-01-02"))); err != nil {
				return err
			}
		}
	}
	return nil
}

func standaloneRoot(targetPath string) (string, string, error) {
	ancestor := filepath.Dir(targetPath)
	for {
		info, err := os.Stat(ancestor)
		if err == nil {
			if !info.IsDir() {
				return "", "", fmt.Errorf("standalone ancestor %q is not a directory", ancestor)
			}
			prefix, relErr := filepath.Rel(ancestor, targetPath)
			return ancestor, cleanRootPath(prefix), relErr
		}
		if !os.IsNotExist(err) {
			return "", "", err
		}
		parent := filepath.Dir(ancestor)
		if parent == ancestor {
			return "", "", fmt.Errorf("no existing standalone ancestor for %q", targetPath)
		}
		ancestor = parent
	}
}

func initializationFailure(target string, cause error, rollback resultmodel.RollbackResult) resultmodel.CommandResult {
	outcome := resultmodel.OutcomeRolledBack
	if rollback.Status == resultmodel.RollbackNotNeeded {
		outcome = resultmodel.OutcomeFailure
	}
	if rollback.Status == resultmodel.RollbackIncomplete {
		outcome = resultmodel.OutcomeRisk
	}
	return resultmodel.CommandResult{Outcome: outcome, Findings: []resultmodel.CommandFinding{knowledgeFinding(CommandBKBInit, "BKB-INIT-FAILED", resultmodel.SeverityError, []string{target}, cause.Error(), resultmodel.FixabilityManual, "initialization failed; inspect preserved replacements before retrying")}, Rollback: rollback}
}

func verifyRootedCommit(root *os.Root, rootPath string, expected []string) error {
	output, err := runGitInRoot(root, rootPath, "show", "--name-only", "--pretty=format:", "HEAD")
	if err != nil {
		return err
	}
	actual := strings.Fields(string(output))
	expected = append([]string(nil), expected...)
	sort.Strings(actual)
	sort.Strings(expected)
	if strings.Join(actual, "\x00") != strings.Join(expected, "\x00") {
		return fmt.Errorf("committed paths %v differ from exact scaffold paths %v", actual, expected)
	}
	if _, err := runGitInRoot(root, rootPath, "diff", "--cached", "--quiet", "--exit-code"); err != nil {
		return fmt.Errorf("index is not empty after commit: %w", err)
	}
	return nil
}

func unstageRootedPaths(root *os.Root, rootPath string, paths []string) {
	if _, err := runGitInRoot(root, rootPath, append([]string{"reset", "--quiet", "--"}, paths...)...); err != nil {
		_, _ = runGitInRoot(root, rootPath, append([]string{"rm", "--cached", "-r", "--force", "--ignore-unmatch", "--"}, paths...)...)
	}
}

func rootedCommittedRisk(root *os.Root, targetPath, target string, operationError error) resultmodel.CommandResult {
	output, err := runGitInRoot(root, targetPath, "rev-parse", "HEAD")
	sha := strings.TrimSpace(string(output))
	finding := knowledgeFinding(CommandBKBInit, "BKB-INIT-COMMITTED-RISK", resultmodel.SeverityError, []string{target}, operationError.Error(), resultmodel.FixabilityManual, "the commit landed but exact paths could not be verified; revert that commit before retrying")
	if err == nil && sha != "" {
		finding.NextArgv = []string{"git", "-C", targetPath, "revert", sha}
		finding.VerificationArgv = []string{"git", "-C", targetPath, "show", "--name-only", "--pretty=format:", sha}
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRisk, Findings: []resultmodel.CommandFinding{finding}}
}

func displayScaffoldPath(rootPath, prefix, target string) string {
	suffix := strings.TrimPrefix(strings.TrimPrefix(cleanRootPath(rootPath), cleanRootPath(prefix)), "/")
	return filepath.ToSlash(filepath.Join(target, filepath.FromSlash(suffix)))
}

func initializeWithoutRepository(root string, options bkbOptions, now time.Time) resultmodel.CommandResult {
	targetPath, err := ensureSafeTarget(root, options.target)
	if err != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{knowledgeFinding(CommandBKBInit, "BKB-INIT-UNSAFE-TARGET", resultmodel.SeverityError, []string{options.target}, err.Error(), resultmodel.FixabilityRefused, "the target is unsafe")}}
	}
	return initializeStandaloneTarget(targetPath, options, now)
}

func initializeStandaloneTarget(targetPath string, options bkbOptions, now time.Time) resultmodel.CommandResult {
	base, prefix, err := standaloneRoot(targetPath)
	if err != nil {
		return initializationFailure(options.target, err, resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded})
	}
	missingFiles, existingFiles := scaffoldInventory(targetPath, filepath.ToSlash(prefix))
	missingDirectories := missingDirectoryInventory(targetPath, filepath.ToSlash(prefix))
	if len(missingFiles) == 0 {
		result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded}}
		for _, path := range existingFiles {
			result.SkippedWork = append(result.SkippedWork, resultmodel.SkippedWork{Code: "BKB-INIT-PRESERVE-EXISTING", Reason: displayScaffoldPath(path, prefix, options.target)})
		}
		return result
	}
	writer, err := newRootedScaffoldWriter(base, "standalone")
	if err != nil {
		return initializationFailure(options.target, err, resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded})
	}
	defer writer.close()
	if err := writer.snapshotDirectories(scaffoldRootDirectories(prefix)); err != nil {
		return initializationFailure(options.target, err, resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded})
	}
	if err := knowledgeMutationHook("after-validation", writer.flow, prefix); err != nil {
		return initializationFailure(options.target, err, writer.rollback())
	}
	if err := applyRootedScaffold(writer, prefix, missingDirectories, missingFiles, now); err != nil {
		return initializationFailure(options.target, err, writer.rollback())
	}
	if _, err := exec.LookPath("git"); err != nil {
		return initializationFailure(options.target, fmt.Errorf("Git is required to initialize a standalone knowledge base: %w", err), writer.rollback())
	}
	gitPath := cleanRootPath(filepath.Join(prefix, ".git"))
	if err := writer.createDirectory(gitPath, true); err != nil {
		return initializationFailure(options.target, err, writer.rollback())
	}
	targetRoot, err := writer.root.OpenRoot(prefix)
	if err != nil {
		return initializationFailure(options.target, err, writer.rollback())
	}
	defer targetRoot.Close()
	if _, err := runGitInRoot(targetRoot, targetPath, "init"); err != nil {
		return initializationFailure(options.target, err, writer.rollback())
	}
	if options.commit && len(missingFiles) > 0 {
		commitPaths := make([]string, 0, len(missingFiles))
		for _, path := range missingFiles {
			commitPaths = append(commitPaths, strings.TrimPrefix(strings.TrimPrefix(path, cleanRootPath(prefix)), "/"))
		}
		if _, err := runGitInRoot(targetRoot, targetPath, append([]string{"add", "-A", "--"}, commitPaths...)...); err != nil {
			return initializationFailure(options.target, err, writer.rollback())
		}
		if _, err := runGitInRoot(targetRoot, targetPath, "commit", "-m", "bkb: initialize knowledge base"); err != nil {
			return initializationFailure(options.target, err, writer.rollback())
		}
		if err := verifyRootedCommit(targetRoot, targetPath, commitPaths); err != nil {
			return rootedCommittedRisk(targetRoot, targetPath, options.target, err)
		}
	}
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded}}
	for _, relative := range missingFiles {
		result.Changes = append(result.Changes, resultmodel.RecordedChange{Path: displayScaffoldPath(relative, prefix, options.target), Kind: "created", Detail: "canonical BKB scaffold file"})
	}
	for _, relative := range existingFiles {
		result.SkippedWork = append(result.SkippedWork, resultmodel.SkippedWork{Code: "BKB-INIT-PRESERVE-EXISTING", Reason: displayScaffoldPath(relative, prefix, options.target)})
	}
	return result
}

func standaloneCommittedRisk(targetPath, targetRelative string, operationError error) resultmodel.CommandResult {
	commitOutput, commitError := exec.Command("git", "-C", targetPath, "rev-parse", "HEAD").Output()
	commitSHA := strings.TrimSpace(string(commitOutput))
	finding := knowledgeFinding(CommandBKBInit, "BKB-INIT-COMMITTED-RISK", resultmodel.SeverityError, []string{targetRelative}, operationError.Error(), resultmodel.FixabilityManual, "the commit landed but exact paths could not be verified; revert that commit before retrying")
	if commitError == nil && commitSHA != "" {
		finding.NextArgv = []string{"git", "-C", targetPath, "revert", commitSHA}
		finding.NextJustRecipe = ""
		finding.VerificationArgv = []string{"git", "-C", targetPath, "show", "--name-only", "--pretty=format:", commitSHA}
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRisk, Findings: []resultmodel.CommandFinding{finding}}
}

func scaffoldInventory(targetPath, targetRelative string) ([]string, []string) {
	missing, existing := []string{}, []string{}
	for _, file := range scaffoldFiles {
		relative := filepath.ToSlash(filepath.Join(targetRelative, file.path))
		if _, err := os.Lstat(filepath.Join(targetPath, filepath.FromSlash(file.path))); os.IsNotExist(err) {
			missing = append(missing, relative)
		} else {
			existing = append(existing, relative)
		}
	}
	sort.Strings(missing)
	sort.Strings(existing)
	return missing, existing
}

func missingDirectoryInventory(targetPath, targetRelative string) []string {
	allDirectories := []string{}
	for _, directory := range scaffoldDirectorySuffixes() {
		allDirectories = append(allDirectories, filepath.ToSlash(filepath.Join(targetRelative, directory)))
	}
	seen, missing := map[string]bool{}, []string{}
	for _, relative := range allDirectories {
		if seen[relative] {
			continue
		}
		seen[relative] = true
		suffix := strings.TrimPrefix(filepath.FromSlash(relative), filepath.FromSlash(targetRelative))
		suffix = strings.TrimPrefix(suffix, string(filepath.Separator))
		absolutePath := targetPath
		if suffix != "" {
			absolutePath = filepath.Join(targetPath, suffix)
		}
		if _, err := os.Lstat(absolutePath); os.IsNotExist(err) {
			missing = append(missing, relative)
		}
	}
	sort.Slice(missing, func(i, j int) bool {
		ai, aj := strings.Count(missing[i], "/"), strings.Count(missing[j], "/")
		if ai != aj {
			return ai < aj
		}
		return missing[i] < missing[j]
	})
	return missing
}

func enclosingGitRoot(root string) (string, error) {
	output, err := exec.Command("git", "-C", root, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return filepath.Clean(strings.TrimSpace(string(output))), nil
}

func isKnowledgeBase(path string) bool {
	raw, rawErr := os.Stat(filepath.Join(path, "raw"))
	wiki, wikiErr := os.Stat(filepath.Join(path, "wiki"))
	return rawErr == nil && wikiErr == nil && raw.IsDir() && wiki.IsDir()
}

func validateScaffoldTarget(targetPath string) error {
	for _, directory := range scaffoldDirectorySuffixes() {
		path := filepath.Join(targetPath, filepath.FromSlash(directory))
		info, err := os.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("directory scaffold path %q is not a real directory", path)
		}
	}
	for _, file := range scaffoldFiles {
		path := filepath.Join(targetPath, filepath.FromSlash(file.path))
		info, err := os.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("file scaffold path %q is not a regular file", path)
		}
	}
	return nil
}

func scaffoldDirectorySuffixes() []string {
	directories, seen := []string{"."}, map[string]bool{".": true}
	for _, directory := range scaffoldDirectories {
		parts := strings.Split(directory, "/")
		for index := range parts {
			path := filepath.ToSlash(filepath.Join(parts[:index+1]...))
			if !seen[path] {
				seen[path] = true
				directories = append(directories, path)
			}
		}
	}
	sort.Slice(directories, func(i, j int) bool {
		leftDepth, rightDepth := strings.Count(directories[i], "/"), strings.Count(directories[j], "/")
		if leftDepth != rightDepth {
			return leftDepth < rightDepth
		}
		return directories[i] < directories[j]
	})
	return directories
}
func containsString(values []string, wanted string) bool {
	index := sort.SearchStrings(values, wanted)
	return index < len(values) && values[index] == wanted
}
