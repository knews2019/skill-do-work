# Dead Code Audit

Scan of the `skill-do-work` repository looking for files, references, and dispatch paths that nothing else points to. This is a Markdown-based skill — there's no JavaScript, no `package.json`, no env vars in the traditional sense — so the analogues are:

- **Action files** that no router or action loads
- **Crew-member files** that no action loads (statically or via the domain dispatch contract)
- **Companion reference files** that their parent action never reads
- **Specs / prompts / interviews / hooks / docs** with no inbound reference
- **Dispatched modes / sub-commands** declared in a router that have no matching target

Nothing has been deleted. Items are grouped by how confident the audit is that they are unreachable.

## Scope

Scanned: `SKILL.md`, `CLAUDE.md`, `README.md`, `AGENTS.md`, `next-steps.md`, `CHANGELOG.md`, and all `*.md` under `actions/`, `crew-members/`, `specs/`, `prompts/`, `interviews/`, `docs/`, `decisions/`, plus `hooks/` and `.vscode/`.

Reference patterns matched: `actions/X.md`, `crew-members/X.md`, `crew-members/[domain].md` dynamic dispatch (resolved against the canonical `domain` enum), `specs/X.md`, `prompts/X.md`, `interviews/X.md`, `docs/X.md`, and `hooks/X` filenames.

---

## Certain

Files that cannot be reached by any documented dispatch path.

### `crew-members/performance.md`

**Why it looks dead:**

- No action loads it statically (no `crew-members/performance.md` reference exists in any action file).
- Its only plausible loader is the dynamic `crew-members/[domain].md` dispatch in `actions/work.md:512`.
- That dispatch normalizes `domain` against the canonical enum and falls back to `general` for any unknown value. The canonical enum (`actions/work.md:153`, reinforced at `actions/work.md:194`, mirrored at `actions/capture.md:88`) is **`frontend | backend | ui-design | general`** — `performance` is not in it.
- `actions/capture.md:190` (Schema Read Contract) explicitly requires a typo warning + user confirmation before writing a non-canonical `domain` value. So a REQ with `domain: performance` would either get rewritten to a canonical value at capture time, or get warned-and-fallen-back to `general` at work time. Either way, `performance.md` is never loaded.

**Mismatch with documentation:** `CLAUDE.md:144` lists `performance.md` as one of the example domain crew files, alongside `backend.md`, `frontend.md`, `ui-design.md`. The example is stale — the contract no longer accepts `performance` as a domain.

**Resolution options (don't pick one; just listing them):**

1. Add `performance` to the canonical `domain` enum in `actions/work.md` and `actions/capture.md`, and the contract starts loading the file.
2. Remove `crew-members/performance.md` and drop the `performance.md` mention from `CLAUDE.md:144`.
3. Load it directly from a non-domain code path (the way `actions/code-review.md:140` loads `crew-members/security.md`), if there's a specific action where its content belongs.

---

## Probable

Files that have no inbound reference and look like historical archives, but where it's plausible the author intends them to be browsable without being individually linked.

### `decisions/imported-specs/` — both files

- `decisions/imported-specs/2026-04-12_close-gaps-in-interview.md`
- `decisions/imported-specs/2026-04-17_improve-weekly-diff-skill.md`

**Why it looks dead:**

- No `actions/*.md`, `SKILL.md`, `CLAUDE.md`, `README.md`, or `next-steps.md` mentions the `imported-specs/` path.
- The sibling `decisions/records/`, `decisions/topics/`, and `decisions/log.md` are part of an ADR system that the `architecture-decisions-log_create-or-expand` prompt knows how to extend, but `imported-specs/` is not mentioned in the prompt's resolution flow either.

**Plausible non-dead use:** these are imported source documents preserved as evidence behind specific ADRs. Linked from inside individual ADR records, but not from any dispatcher. That's a normal pattern for an audit trail — they're "leaves," not entry points. **Recommendation: keep, but verify each is cited from at least one ADR; if not, they can be moved or pruned.**

### Per-action user guides in `docs/` that are not linked from anywhere except the generic `docs/` directory mention

`README.md:67` says *"Per-action guides live in [`docs/`](./docs/)"* and then individually links only `docs/capture-guide.md` and `docs/work-guide.md`. `actions/interview.md:483` links `docs/interview-guide.md`. Every other guide is reachable only by directory browse:

| Guide | Inbound references besides the generic `docs/` link |
|---|---|
| `docs/bkb-guide.md` | none |
| `docs/cleanup-guide.md` | none |
| `docs/code-review-guide.md` | none |
| `docs/commit-guide.md` | none |
| `docs/forensics-guide.md` | none |
| `docs/inspect-guide.md` | none |
| `docs/present-work-guide.md` | none |
| `docs/prime-guide.md` | none |
| `docs/prompts-guide.md` | none |
| `docs/quick-wins-guide.md` | none |
| `docs/review-work-guide.md` | none |
| `docs/roadmap-guide.md` | none |
| `docs/ui-review-guide.md` | none |
| `docs/verify-requests-guide.md` | none |
| `docs/version-guide.md` | none |

**Why probable rather than certain:** the generic `docs/` link in `README.md` is enough for a reader to discover them. But each guide's action file does not point back to its guide, which means a user reading e.g. `actions/commit.md` has no signal that `docs/commit-guide.md` exists. They are reachable but undiscoverable from inside the skill.

---

## Needs Human Check

Items where an automated audit can't tell intent — they could be deliberate or stale.

### `prompts/weekly-signal-diff-personal.md`

`prompts/README.md` describes it as *"Placeholder template for the personal sidecar. Ships with no real lanes... Not run directly."* So it's intentionally inert from the dispatcher's perspective — it's a copy-paste source for users. But `actions/prompts.md` (the dispatcher) does not know to skip it. A user typing `do-work prompts run weekly-signal-diff-personal` would get the placeholder body adopted as instructions. **Worth confirming:** is the dispatcher supposed to refuse-with-explanation when a prompt's header carries an opt-out marker, or is being a copy-paste template a fine de-facto state?

### `AGENTS.md`

Contains exactly one line: `READ CLAUDE.md`. This is the [AGENTS.md convention](https://agents.md/) — agents that look for `AGENTS.md` by default get redirected to the project's `CLAUDE.md`. Not dead, but worth flagging that anyone who deletes `CLAUDE.md` will leave a stub pointing at nothing.

### `.vscode/tasks.json`

A single VSCode task ("Open current HTML in browser"). Useful only if you (a) use VSCode and (b) preview HTML deliverables produced by `present-work`. Not dead in any reachable sense, but it's the kind of editor-specific file that quietly accumulates. **Worth confirming:** is this still in active use by the maintainer?

### `decisions/_progress.md`, `decisions/_master_index.md`, `decisions/log.md`

These are the running ADR ledger — they're updated by the `architecture-decisions-log_create-or-expand` prompt. No action loads them directly (they're meant to be human-readable artifacts). Not dead — but a casual reader scanning for "where is this file used?" would find nothing. Flagged for completeness, not for action.

### `CHANGELOG.md` history older than a few releases

Not strictly dead code, but the file has grown to ~82KB / 2000+ lines. `actions/version.md` only ever shows the last 5 entries. The full history is preserved in git, so the older entries are duplicative storage. **Worth confirming:** is the full in-tree history intentional (e.g., for offline browsing), or is rolling-window archival OK?

---

## Verified Live (no findings)

Recording these so the audit can be re-run against the same scope:

- **All 28 files in `actions/`** are reachable. 24 are routed from `SKILL.md`'s routing table; `kb-lessons-handoff.md` is invoked from `CLAUDE.md:170` and `actions/work.md` Step 7.5; the four companion references (`bkb-reference.md`, `deep-explore-reference.md`, `interview-reference.md`, `sample-archived-req.md`) are each cited multiple times from their sister action.
- **All 4 specs in `specs/`** (`api-endpoint`, `ui-component`, `refactor`, `bug-fix`) are explicitly listed in `actions/work.md:396` for task-type matching and in `specs/README.md`'s "Available specs" section.
- **All 20 prompt files in `prompts/`** are cataloged in `prompts/README.md`'s "Available prompts" table — the dispatcher's source of truth.
- **All 3 hooks in `hooks/`** are documented in `README.md:135–140` and referenced from `CLAUDE.md`'s file index.
- **All 12 crew-members are loaded except `performance.md`:**
  - Always-loaded: `general.md`, `karpathy.md` (`actions/work.md:510–511`).
  - Action-specific direct loads: `security.md` (code-review), `ui-design.md` (ui-review), `interviewer.md` (interview), `approach-directives.md` (work, parallel-REQ mode), `caveman.md` (work, when `caveman` frontmatter set), `debugging.md` (work, on review-fail retry + after 2+ test failures), `testing.md` (work, when `tdd: true` or `domain: testing`).
  - Domain-dispatched: `backend.md`, `frontend.md`, `ui-design.md` are loaded by `actions/work.md:512` when the REQ's `domain` matches the canonical enum.
- **The one interview template** (`interviews/work-operating-model.md`) is referenced from `actions/interview.md` and `CLAUDE.md`.
- **The full `decisions/records/` directory** (12 ADRs) is the corpus the `architecture-decisions-log_create-or-expand` prompt manages; ADR-001 through ADR-012 are sequenced and linked from the topic indexes (`decisions/topics/_index_*.md`).

---

## Method Notes

- "Reference" means a literal string match for the file's path or basename inside another tracked file. Dynamic dispatch was resolved by reading the actual lookup code: e.g., `crew-members/[domain].md` was checked against every documented value the `domain` frontmatter is allowed to take, not against grep alone.
- Documentation describing a file (e.g., `CLAUDE.md`'s file inventory listing every action) was not counted as a "reference" — that's catalog text, not a dispatch. A reference had to be a path or basename appearing in an executable context (a router table, a load instruction, a "see also" link a reader would follow).
- This is a static analysis. A file flagged as dead could still be loaded by a user pasting it into an LLM directly or by a feature added after this audit. Re-run the audit when the routing table, the Schema Read Contract, or the crew-loading mechanism changes.
