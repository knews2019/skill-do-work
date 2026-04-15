# ADR Log

> Create or update a project-wide Architecture Decision Record log at `decisions/`, modeled on the BKB wiki pattern. Mines `CHANGELOG.md` for load-bearing, still-in-force decisions. Resumable, idempotent, handles supersession.

**Aliases:** `adr`, `decisions`

**When to use:**
- The repo has a rich `CHANGELOG.md` or release history and decisions are scattered / undocumented
- You want a durable, wiki-linked record of *why* things are the way they are
- Someone is about to join the project and needs to understand the load-bearing choices

**Inputs / flags:**
- `--since <version>` — mine CHANGELOG from this version forward (default: full file)
- `--last <N>` — mine only the last N CHANGELOG entries
- `--dry-run` — print the plan without writing files or committing
- `--no-push` — commit but don't push
- `--batch-size <N>` — ADRs per commit batch (default: 3)

---

## Instructions for the executing agent

You are maintaining an Architecture Decision Record (ADR) log at `decisions/` in the current repo, modeled on the BKB wiki pattern from this skill (see `actions/build-knowledge-base.md` for the wiki-structure conventions this prompt reuses — `_master_index.md`, topic clusters, typed `[[wiki-links]]`, daily `log.md`).

This prompt is **idempotent and resumable**: it detects whether the log already exists and whether a prior run was interrupted, and resumes from there without re-scaffolding or duplicating ADRs.

### Phase 0 — Pre-flight

1. **Confirm the working tree is clean.** Run `git status --porcelain`. If dirty, stop and ask the user to stash or commit first; do not mix unrelated changes into the ADR batches.
2. **Confirm the current branch is not `main`/`master`.** If it is, stop and ask the user for a feature branch name (or offer to create `adr-log` as one). Every batch gets pushed — pushing to main is never acceptable without explicit user authorization.
3. **Confirm `CHANGELOG.md` exists** in the repo root. If not, stop and tell the user this prompt requires a CHANGELOG as source material.
4. **Parse arguments.** Apply `--since`, `--last`, `--dry-run`, `--no-push`, `--batch-size` as stated above.

### Phase 1 — State detection

Check for `decisions/_master_index.md`:

- **Absent → CREATE mode.** Proceed to Phase 2 (scaffolding).
- **Present → UPDATE mode.** Read `decisions/_progress.md`:
  - If `status: complete`, re-run the mining pass (Phase 3) with current CHANGELOG and only create ADRs for decisions not already covered. Skip Phase 2 entirely.
  - If `status: in-progress`, jump to the batch after the last completed one. Do not re-run earlier batches.
  - If `_progress.md` is missing but the index exists, reconstruct state by scanning `decisions/ADR-*.md` for the highest number and assume `status: complete`.

Record which mode you're in at the top of your response to the user.

### Phase 2 — Scaffolding (CREATE mode only, BATCH 1)

Create the following files. All new content goes in `decisions/`.

**`decisions/_master_index.md`** — the wiki entry point. Lists topic clusters and links to every ADR. Initial skeleton (populated in Phase 4 as ADRs are written):

```markdown
# Architecture Decisions — Master Index

Project-wide ADR log. Each entry records a load-bearing decision, its alternatives, and its consequences. Links use `[[ADR-NNNN-kebab-slug]]` syntax.

## Topic clusters

_(populated as ADRs are written)_

## All ADRs

_(populated as ADRs are written — most recent first)_

## Status legend

- **accepted** — currently in force
- **superseded** — replaced by a newer ADR (see `superseded_by`)
- **deprecated** — no longer in force, not yet replaced
- **proposed** — under discussion, not yet adopted
```

**`decisions/log.md`** — the chronological timeline:

```markdown
# ADR Timeline

Most recent decisions at the top. Each entry: date, ADR id + title, one-line summary.

_(populated as ADRs are written)_
```

**`decisions/_progress.md`** — the resume-state ledger:

```markdown
# ADR Log Progress

status: in-progress
mode: create
started_at: <ISO-8601 timestamp>
source: CHANGELOG.md
scope: <"full" | "since X.Y.Z" | "last N">
batch_size: 3

## Batches

- [x] BATCH 1 — scaffolding
- [ ] BATCH 2 — mining + shortlist
- [ ] BATCH 3+ — ADRs in groups of 3
- [ ] FINAL — index reconciliation, log sort, cross-link audit

## Mined decisions (shortlist)

_(populated in Phase 3)_

## Completed ADRs

_(appended as they're written)_
```

**`decisions/_topic-clusters.md`** — the cluster taxonomy (living document):

```markdown
# Topic Clusters

Clusters are populated as ADRs are mined. A cluster is created when 2+ ADRs share a theme. Until then, ADRs live under "Uncategorized" in the master index.

_(populated as ADRs are written)_
```

**Commit this batch** with message `adr: scaffold decisions/ log` and push (unless `--no-push`). Update `_progress.md` to mark BATCH 1 complete before committing.

### Phase 3 — Mining (BATCH 2)

Read `CHANGELOG.md` (scoped by `--since` / `--last` if provided) and extract a shortlist of **load-bearing, still-in-force decisions**.

**Selection test — a CHANGELOG entry qualifies if ALL of these hold:**

1. It records a **decision** (an intentional choice of A over B), not a bugfix, typo, or purely mechanical refactor.
2. It is **referenced by current code, config, or docs** — a reader today still sees its fingerprint. Spot-check by grepping for the relevant files/symbols; if they no longer exist, the decision is obsolete.
3. It would **cause meaningful rework** if reversed (behavior change, migration, retraining, documentation churn).
4. It has **not been explicitly undone** by a later CHANGELOG entry.

**Supersession signals while mining:** if a later CHANGELOG entry explicitly replaces an earlier one, record the pair. The earlier decision may still warrant an ADR (with `status: superseded`) if the replacement is itself load-bearing and the history matters.

**In UPDATE mode, deduplicate:** for each candidate, check every existing ADR's `source:` frontmatter. If the CHANGELOG version/section is already cited there, skip the candidate. If a candidate supersedes an existing ADR, queue a supersession update (Phase 4b below).

Write the shortlist into `_progress.md` under `## Mined decisions (shortlist)` as a table: `#`, `CHANGELOG source`, `proposed slug`, `cluster guess`, `supersedes (if any)`. Commit this batch with message `adr: mine CHANGELOG shortlist` and push. Update `_progress.md` to mark BATCH 2 complete.

### Phase 4 — ADRs in groups of 3 (BATCH 3 … N)

For each group of N ADRs (default 3, controlled by `--batch-size`):

1. Allocate sequential numbers starting from `max(existing ADR number) + 1`, zero-padded to 4 digits. Never reuse a number.
2. Write each ADR to `decisions/ADR-NNNN-kebab-slug.md` using the template below.
3. Update `decisions/_master_index.md` — add each new ADR under its cluster and under "All ADRs".
4. Update `decisions/log.md` — prepend a one-line timeline entry per new ADR (most recent first).
5. Update `decisions/_topic-clusters.md` if a new cluster emerged.
6. If any ADR in this group supersedes an existing one, run Phase 4b for that pair before committing.
7. Append the batch's ADRs to `_progress.md` under `## Completed ADRs`.
8. Tick the batch checkbox in `_progress.md`.
9. Commit with message `adr: ADR-NNNN..ADR-NNNN <short description>` and push.

**ADR file template:**

```markdown
---
id: ADR-NNNN
title: <Short imperative title>
status: accepted            # accepted | superseded | deprecated | proposed
date: YYYY-MM-DD            # date the decision was made, per CHANGELOG
cluster: <cluster-slug>     # or "uncategorized"
source:                     # the CHANGELOG evidence this ADR summarizes
  - CHANGELOG.md: "X.Y.Z — Codename"
supersedes: []              # list of ADR ids this replaces
superseded_by: null         # filled in when this one is later replaced
tags: []                    # optional freeform tags
---

# ADR-NNNN — <Title>

## Context

What forces were in play when this decision was made? What pressures, constraints, or goals made a choice necessary? Cite the relevant CHANGELOG entry and any code/config that embodies the constraint. Use `[[wiki-links]]` for related ADRs.

## Decision

The choice that was made, stated in a single clear sentence, then expanded with the specifics — what was added, changed, removed, or standardized. Name files/modules/patterns when relevant.

## Alternatives considered

- **<Option A>** — <why it was rejected or deferred>
- **<Option B>** — <why it was rejected or deferred>
- <At least 2 alternatives. If the CHANGELOG doesn't discuss alternatives, infer plausible ones from context and mark them `(inferred)`.>

## Consequences

Positive, negative, and neutral effects of the decision. What becomes easier? What becomes harder? What new obligations does this create? What is now load-bearing that wasn't before?

## References

- CHANGELOG.md `X.Y.Z — Codename`
- `path/to/code.ext` — the embodiment of the decision
- `[[ADR-MMMM-related-slug]]` — sibling or predecessor decisions
```

**Rules for each ADR:**

- **Every section must be present.** If there's truly nothing to say under Alternatives, say so explicitly — don't omit the heading.
- **Use `[[ADR-NNNN-kebab-slug]]` links** when pointing at sibling ADRs. The dispatcher's `list` output uses these names, so keep slugs consistent.
- **Typed links** (BKB pattern): prefix with the relationship, e.g. `[[supersedes::ADR-0003-...]]`, `[[refines::ADR-0007-...]]`, `[[conflicts::ADR-0011-...]]`. Use typed links in the References section whenever the relationship is non-obvious.
- **Cluster slugs** stay stable — once you commit to `build-system` don't later rename to `build`. Edit `_topic-clusters.md` deliberately, not drive-by.

### Phase 4b — Supersession (runs inline during Phase 4 when needed)

When a new ADR supersedes an existing one:

1. In the new ADR: set `supersedes: [ADR-OLD]` and add to References: `[[supersedes::ADR-OLD-slug]]`.
2. In the old ADR: set `status: superseded` and `superseded_by: ADR-NEW`. Append to its References: `Superseded by [[ADR-NEW-slug]]`.
3. **Never delete the old ADR.** History stays intact.
4. In `_master_index.md`, mark the old entry with `(superseded by ADR-NEW)` inline.
5. Include the old-ADR edit in the same commit as the new ADR.

### Phase 5 — Final reconciliation (LAST BATCH)

1. Re-read `_master_index.md` and verify every ADR file is listed exactly once under exactly one cluster, with "All ADRs" listing every ADR regardless of cluster.
2. Re-read `log.md` and verify the timeline is date-sorted, newest first.
3. Audit wiki-links: grep all `[[…]]` references and verify each target ADR exists. Report broken links to the user; do not silently "fix" them by dropping the link.
4. Set `_progress.md` → `status: complete` and add a `completed_at: <ISO-8601>` field.
5. Commit with message `adr: finalize index + log` and push.

### Phase 6 — Report

Tell the user:

- Mode (create / update / resume-from-batch-N)
- Counts: new ADRs written, existing ADRs superseded, batches committed, ADRs already covered (skipped)
- Clusters that emerged
- Broken-link findings, if any
- Branch name and last commit hash

### Rules

- **Never delete an existing ADR.** Supersede, don't rewrite history.
- **Never renumber existing ADRs.** Numbers are load-bearing — external references depend on them.
- **Never skip the supersession check.** If you write a new ADR that contradicts an old one without updating the old one's `status` + `superseded_by`, you've created a silent conflict.
- **Never push to `main`/`master`** unless the user has already explicitly authorized it for this specific run.
- **`--dry-run` means read-only.** No file writes, no commits, no pushes. Report what *would* have happened.
- **If the shortlist is empty in UPDATE mode**, report "no new load-bearing decisions found since last run" and stop — do not invent ADRs.

### Common rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "The CHANGELOG entry is small, I'll merge two decisions into one ADR" | One decision per ADR, even if small | ADRs are referenced by id; merged ADRs become unlinkable the moment one half is superseded |
| "This old ADR is clearly wrong now, I'll just edit its Decision section" | Supersede it with a new ADR | Editing the Decision destroys the historical record of what was believed true at the time |
| "The user didn't specify a cluster, I'll leave it uncategorized forever" | Promote it to a named cluster as soon as 2+ ADRs share the theme | Uncategorized bloat defeats the point of topic clusters |
| "I'll skip `_progress.md` this batch, I'll remember where I am" | Always update `_progress.md` before committing the batch | Agents don't persist across sessions; the file IS the memory |
| "The CHANGELOG doesn't list alternatives, I'll omit that section" | Infer 1–2 plausible alternatives and mark them `(inferred)` | An ADR with no Alternatives section reads like a decree, not a decision |

### Red flags (observable symptoms something went wrong)

- Two ADRs share a number, or an ADR id has a gap (`ADR-0003` exists, `ADR-0004` doesn't, `ADR-0005` exists)
- A superseded ADR's `status` is still `accepted`
- `_master_index.md` lists an ADR file that doesn't exist, or vice versa
- `log.md` entries are out of date order
- `_progress.md` says `status: complete` but the mined shortlist still has unticked items
- Duplicate ADRs for the same CHANGELOG source (idempotency failure)
- A commit batch included unrelated files (ADR batches must be surgical)

### Verification checklist

- [ ] Phase 0 pre-flight all passed (clean tree, non-main branch, CHANGELOG present)
- [ ] Mode correctly identified (create vs update vs resume) and reported to user
- [ ] `_progress.md` exists and every completed batch is ticked
- [ ] Every ADR has all five sections (Context, Decision, Alternatives, Consequences, References)
- [ ] Every ADR's `source:` cites at least one CHANGELOG entry
- [ ] Every supersession has both sides updated in the same commit
- [ ] `_master_index.md` and `log.md` reconciled in the final batch
- [ ] All batches pushed (or `--no-push` honored uniformly)
- [ ] No ADR numbers reused; no gaps in the sequence
- [ ] Final report to user includes counts, branch, last commit hash, and any broken-link findings
