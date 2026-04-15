# ADR Log

> Create or update a project-wide Architecture Decision Record log at `decisions/`, modeled on the BKB wiki pattern. Layered source mining (implementation-history → lessons-learned → code). Idempotent via REQ/UR keys. Resumable, supersession-aware, pre-flight-checked.

**Aliases:** `adr`, `decisions`

**When to use:**
- The repo has an `implementation-history.md` (REQ/UR ledger) or a rich `CHANGELOG.md` and the architectural *why* is scattered
- You want a durable, wiki-linked record that new contributors can navigate in minutes
- An existing ADR log needs to absorb recent work without duplication

**Inputs / flags:**
- `--from <UR-NNN|REQ-NNN>` — skip every UR/REQ before this one (default: mine everything)
- `--batch-size <N>` — ADRs per commit (default: 3)
- `--dry-run` — print the plan without writing, committing, or pushing
- `--no-push` — commit each batch but skip push

---

## Instructions for the executing agent

You are maintaining an Architecture Decision Record (ADR) log at `decisions/` in the current repo, modeled on the BKB wiki pattern (see `actions/build-knowledge-base.md`). This prompt is **idempotent and resumable**: it detects existing state and picks up where a prior run left off, keyed on REQ/UR IDs so the same decision is never captured twice.

### Sources to mine (priority order)

1. **`implementation-history.md`** — primary spine. Scan every UR/REQ entry. This is the natural key.
2. **`lessons-learned/`** — supplementary context, rationale, and nuance.
3. **Current code** — verify each candidate decision is still visibly in force.
4. **`CHANGELOG.md`** — fallback spine if `implementation-history.md` is absent. When using the fallback, derive synthetic REQ/UR keys from the version string (`changelog:X.Y.Z`) and record them in the `sources:` field.

### Phase 0 — Pre-flight

1. Run `git status --porcelain`. If the working tree is dirty, **stop** and tell the user to stash or commit first. Do not mix unrelated changes into ADR batches.
2. Run `git rev-parse --abbrev-ref HEAD`. If the branch is `main` or `master`, **stop** and ask the user for a feature branch. Every batch is pushed — pushing to main is never acceptable without explicit authorization.
3. Confirm at least one source exists: `implementation-history.md` OR `CHANGELOG.md`. If neither, **stop**.
4. Parse `--from`, `--batch-size`, `--dry-run`, `--no-push`.

### Phase 1 — State detection

Check for `decisions/_master_index.md`:

- **Absent → CREATE mode.** Proceed to Phase 2 (scaffolding).
- **Present → UPDATE mode.** Read `decisions/_progress.md` for the next ADR number and the list of deferred candidates. Skip Phase 2 entirely. Do NOT re-scaffold `_master_index.md`, `_progress.md`, `log.md`, or any `topics/_index_*.md`.
- **Present but `_progress.md` missing** → RESUME mode. Reconstruct state by scanning `decisions/records/adr-*.md` for the highest existing number and every `req:`/`ur:` already referenced; then treat the run as UPDATE mode.

Report the detected mode to the user before doing any writes.

### Phase 2 — Scaffolding (CREATE mode only, BATCH 1)

Create:

```
decisions/
├── _master_index.md           # wiki entry point: clusters, ADR index, legend
├── _progress.md               # resume state: next ADR number, deferred list, completed inventory
├── log.md                     # chronological timeline, newest first
├── records/                   # individual ADR files live here
└── topics/                    # per-cluster wiki pages (created as clusters emerge)
```

**`_master_index.md`** seed:

```markdown
# Architecture Decisions — Master Index

Project-wide ADR log. Wiki-links use `[[name]]` syntax.

## Topic clusters

_(populated as clusters emerge — each links to topics/_index_<cluster>.md)_

## All ADRs

_(populated as ADRs are written — most recent first)_

## Status legend

- **accepted** — currently in force
- **superseded** — replaced by a newer ADR (see `related.superseded-by`)
- **deprecated** — no longer in force, not yet replaced
- **proposed** — under discussion, not yet adopted
```

**`_progress.md`** seed:

```markdown
---
status: in-progress
mode: create
started_at: <ISO-8601>
next_adr_number: 1
batch_size: 3
---

# ADR Log Progress

## Deferred candidates
_(UR/REQ items mined but not yet written — populated in Phase 3)_

## Completed ADRs
_(appended as each batch commits)_
```

**`log.md`** seed:

```markdown
# ADR Timeline

Most recent decisions on top.

_(populated as ADRs are written)_
```

Commit with `docs(adr): scaffold decisions/ log` and push (unless `--no-push`). Tick BATCH 1 in `_progress.md` before committing.

### Phase 3 — Mining (BATCH 2)

Walk the sources in priority order. For each UR/REQ:

**Skip if**: a `req:` or `ur:` field in any existing `decisions/records/adr-*.md` already references this ID (or, for CHANGELOG fallback, any `sources:` entry already cites the version).

**ADR selection criteria** — write an ADR only if ALL of these hold:

1. The decision is **architecturally meaningful** — shapes *how* the system is built, not just *what* was done. Bug fixes, typos, and mechanical refactors are out.
2. The **current code still embodies it** — spot-check by reading the relevant files. If the decision is invisible today, it's obsolete.
3. It would **cause meaningful rework if reversed** (behavior change, migration, retraining, doc churn).
4. It is **not already covered** by another ADR's `req:`/`ur:` fields (idempotency).

**Supersession signal**: if a later UR/REQ explicitly replaces an earlier decision, queue the pair. The older ADR may still be worth writing (with `status: superseded`) if the history matters.

Write the shortlist to `_progress.md` under `## Deferred candidates` as a table: `#`, `UR/REQ`, `proposed slug`, `cluster guess`, `supersedes (if any)`, `confidence`. Commit with `docs(adr): mine candidates from implementation-history` and push. Tick BATCH 2.

### Phase 4 — ADRs in groups of N (BATCH 3 … N)

For each batch of N ADRs (default `--batch-size` = 3):

1. Allocate sequential numbers starting at `_progress.md.next_adr_number`, zero-padded to 3 digits. Never reuse a number.
2. Write each ADR to `decisions/records/adr-NNN-<kebab-slug>.md` using the template below.
3. Run supersession updates (Phase 4b) for any ADR in this batch that replaces an existing one — both sides must be edited in the same commit.
4. Update `decisions/topics/_index_<cluster>.md` for every cluster touched (create if new, see Topic cluster rules below).
5. Update `decisions/_master_index.md`: new ADRs added under their cluster AND under "All ADRs".
6. Prepend entries to `decisions/log.md` (newest first): `YYYY-MM-DD — [[adr-NNN-slug]] — <one-line summary>`.
7. Update `decisions/_progress.md`: bump `next_adr_number`, append completed inventory, remove written candidates from the deferred list.
8. Commit: `docs(adr): add ADR-NNN through ADR-NNN — <short description>` and push.

**ADR file template:**

````markdown
---
title: "ADR-NNN: <Short imperative title>"
type: architecture-decision-record
status: accepted            # accepted | superseded | deprecated | proposed
topic_cluster: <cluster-slug>
decided: <YYYY-MM-DD>       # the date the decision was made (per source)
req:
  - REQ-NNN
ur:
  - UR-NNN                  # use [] if unknown
sources:
  - <file paths or URLs backing the decision>
related:
  - page: adr-MMM-other-slug
    rel: complements        # complements | depends-on | extends | supersedes | superseded-by
created: <today>
updated: <today>
confidence: high            # high | medium | low (see criteria below)
---

# ADR-NNN: <Title>

Topic cluster: [[_index_<cluster>]]

See also: [[adr-MMM-other-slug]], [[adr-PPP-third-slug]]

## Context

What forces were in play? What pressures, constraints, or goals made a choice necessary? Cite the originating UR/REQ and any code/config that embodies the constraint. Link related ADRs with `[[wiki-links]]`.

## Decision

The choice, stated in a single clear sentence, then expanded with the specifics — what was added, changed, removed, or standardized. Name files/modules/patterns.

## Alternatives

- **<Option A>** — <why rejected or deferred>
- **<Option B>** — <why rejected or deferred>

If the source doesn't discuss alternatives, infer 1–2 plausible ones from context and mark them `(inferred)`. An ADR with no Alternatives reads like a decree, not a decision.

## Consequences

Positive, negative, and neutral effects. What becomes easier? What becomes harder? What new obligations does this create? What is now load-bearing that wasn't before?

## References

- `implementation-history.md` — <UR-NNN section>
- `lessons-learned/<file>.md` — <if applicable>
- `path/to/code.ext` — the embodiment of the decision
- `[[adr-MMM-related-slug]]` — sibling or predecessor decisions
````

**Confidence criteria:**

- **high** — decision is explicitly documented in a source with rationale + alternatives
- **medium** — decision is visible in code but rationale is inferred from context
- **low** — rationale is mostly reconstructed; flag for human review

### Phase 4b — Supersession (runs inline during Phase 4 when needed)

When a new ADR supersedes an existing one:

1. **New ADR**: add `related: [{page: adr-OLD-slug, rel: supersedes}]`.
2. **Old ADR**: flip `status: accepted` → `status: superseded`. Append `{page: adr-NEW-slug, rel: superseded-by}` to its `related` list. Bump its `updated:` field to today.
3. **Old ADR body**: append to References: `Superseded by [[adr-NEW-slug]] — <date>`.
4. **Master index**: annotate the old entry with `(superseded by ADR-NEW)`.
5. **Never delete** the old ADR. History stays intact.
6. Both edits ship in the **same commit** as the new ADR.

### Topic cluster rules

- Each cluster has a first-class wiki page at `decisions/topics/_index_<cluster-slug>.md` with a short description and a list of its ADRs.
- Assign new ADRs to the most relevant existing cluster. Create a new `_index_<cluster>.md` only when the theme genuinely doesn't fit any existing cluster or when 2+ pending ADRs share a new theme.
- New clusters must be added to `_master_index.md` under "Topic clusters" in the same commit that introduces them.
- Cluster slugs are **stable** — once you commit to `build-system`, don't later rename to `build`. Migrations are expensive.

### Phase 5 — Final reconciliation (LAST BATCH)

1. Re-read `_master_index.md`: verify every ADR is listed exactly once under exactly one cluster, and every cluster index is linked.
2. Re-read `log.md`: verify entries are date-sorted, newest first.
3. Audit wiki-links: grep all `[[…]]` references in `decisions/` and verify every target exists. Report broken links — do not silently drop them.
4. Flip `_progress.md` → `status: complete`; add `completed_at: <ISO-8601>`.
5. Commit `docs(adr): finalize index + log` and push.

### Phase 6 — Completion report

Print this to the user, verbatim structure:

```markdown
### ADR Extraction Status
| URs/REQs covered this run | ADRs written |
|---|---|
| UR-XXX (REQ-NNN, REQ-NNN) | ADR-NNN, ADR-NNN |
| UR-YYY (REQ-NNN) | ADR-NNN |

### Remaining Candidates Estimate
- **UR-NNN** — <title> — ~<S|M|L> effort — <why it qualifies as an ADR>
- **UR-NNN** — <title> — ~<S|M|L> effort — <why it qualifies as an ADR>

Total remaining: N UR groups → estimated N–N more ADRs

### Run metadata
- Mode: <create | update | resume-from-batch-N>
- Branch: <branch>
- Last commit: <sha>
- ADRs superseded: N
- Broken links found: N
```

Scan `implementation-history.md` (or the fallback source) and count UR sections that are NOT yet referenced in any ADR's `req:`/`ur:` frontmatter. Effort sizing:

- **S** — one REQ, one cluster, rationale obvious from source
- **M** — multi-REQ or rationale partially inferred
- **L** — cross-cluster, supersession involved, or confidence low

### Rules

- **Never delete an existing ADR.** Supersede, don't rewrite history.
- **Never renumber existing ADRs.** Numbers are external references.
- **Never skip the supersession flip.** A new ADR that contradicts an old one without flipping the old one's `status` creates a silent conflict.
- **Never push to `main`/`master`** unless the user has already explicitly authorized it for this specific run.
- **`--dry-run` means read-only.** No file writes, no commits, no pushes. Print the plan and the would-be completion report.
- **If the shortlist is empty in UPDATE mode**, print the completion report with zero ADRs written and a remaining-candidates section of `(none)`. Do not invent ADRs to justify the run.
- **One decision per ADR**, even if the source lumps two decisions into one entry. Split them.

### Common rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "The UR is small, I'll merge two decisions into one ADR" | One decision per ADR | ADRs are referenced by id; merged ADRs break when one half is superseded |
| "This old ADR is clearly wrong now, I'll edit its Decision section" | Supersede it with a new ADR | Editing destroys the historical record of what was believed true at the time |
| "I'll leave it `uncategorized` forever" | Promote to a named cluster when 2+ ADRs share the theme | Uncategorized bloat defeats the point of clusters |
| "I'll skip `_progress.md` this batch, I'll remember" | Always update `_progress.md` before committing | Agents don't persist across sessions; the file IS the memory |
| "The source doesn't list alternatives, I'll omit the section" | Infer 1–2 and mark `(inferred)` | Alternatives-less ADRs read like decrees |
| "The push failed, I'll `--force` it" | Retry with exponential backoff; if still failing, report and stop | Force-pushing ADR commits can rewrite shared history |
| "No `implementation-history.md`, I'll stop" | Fall back to `CHANGELOG.md` with synthetic keys | The prompt is portable by design |

### Red flags

- Two ADRs share a number, or there's a gap (`adr-003` exists, `adr-004` missing, `adr-005` exists)
- A superseded ADR's `status` is still `accepted`
- `_master_index.md` lists an ADR file that doesn't exist, or vice versa
- `log.md` is not date-sorted
- `_progress.md.status: complete` but deferred candidates remain
- Two ADRs with overlapping `req:`/`ur:` IDs — idempotency failure
- A commit batch includes files outside `decisions/` — ADR batches must be surgical
- A new cluster was created without updating `_master_index.md` in the same commit

### Verification checklist

- [ ] Pre-flight passed (clean tree, non-main branch, at least one source present)
- [ ] Mode detected and reported (create / update / resume)
- [ ] `_progress.md` updated before every commit
- [ ] Every ADR has all five body sections (Context, Decision, Alternatives, Consequences, References)
- [ ] Every ADR's `req:` and `ur:` fields are populated (or `[]` where unknown)
- [ ] Every supersession has both sides updated in the same commit
- [ ] Every new cluster has a `topics/_index_<cluster>.md` AND a `_master_index.md` entry, in the same commit
- [ ] `log.md` prepended with one line per new ADR, date-sorted
- [ ] No ADR numbers reused; no gaps in the sequence
- [ ] Commit messages match `docs(adr): …` format
- [ ] All batches pushed (or `--no-push` honored uniformly)
- [ ] Completion report printed with extraction-status table, remaining-candidates estimate, and run metadata
