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

(None.) An earlier draft of this report classified `crew-members/performance.md` as certain-dead based on the `domain` Schema Read Contract at `actions/work.md:180` falling back to `general` for unknown values. That was too strong — see the entry in **Probable** below for the corrected analysis. Nothing else in the scan rose to certain-dead.

---

## Probable

Files where the audit can't find an inbound reference that would actually fire under normal use, but where a plausible path of use exists in edge cases or in documentation that may have drifted.

### `crew-members/performance.md`

**Why it looks dead:**

- No action loads it statically — no `crew-members/performance.md` reference exists anywhere.
- Three sites in `actions/work.md` dispatch on `crew-members/[domain].md` dynamically: Step 4 Route C planning (`actions/work.md:407`), Step 6 implementation (`actions/work.md:512`), and the Step 9 review-work spawn (`actions/work.md:658`).
- `performance` is not in the canonical `domain` enum. The enum (`actions/work.md:153`, table at `actions/work.md:194`, mirrored at `actions/capture.md:88`) is **`frontend | backend | ui-design | general`**. Capture (`actions/capture.md:190`) prompts the user to correct a non-canonical `domain` value before writing the REQ.
- The Schema Read Contract narrative at `actions/work.md:180` says *"every read site in this file ... honors a uniform normalize-and-warn contract"*, which would mean Route C and review-spawning both fall back to `general` for `domain: performance` and `performance.md` is never loaded.

**Why it might still be reachable (the reason this isn't certain-dead):**

- The Step 6 entry at `actions/work.md:512` explicitly calls out normalization (*"normalize the REQ's `domain` frontmatter per the Schema Read Contract first"*). The Route C entry at `actions/work.md:407` and the review-spawning entry at `actions/work.md:658` do not — they just check *"if the file exists, load it"*. So whether the contract is enforced at those sites is interpretation-dependent.
- A manually-authored REQ in `do-work/queue/` (skipping capture) with `domain: performance` would, under the literal reading of Route C and review-spawning, find `crew-members/performance.md` and load it.
- The Schema Read Contract's per-field table for `domain` lists *"Step 6 crew load"* as the only read site, which lines up with the literal reading rather than the narrative.

**Mismatch with documentation:** `CLAUDE.md:144` lists `performance.md` as one of the example domain crew files. Either the example is stale (and `performance` was dropped from the enum) or the enum drifted away from intent. Don't delete based on this audit — resolve the contract ambiguity first.

**Resolution options (don't pick one; just listing them):**

1. Add `performance` to the canonical `domain` enum in `actions/work.md` and `actions/capture.md` — then all three dispatch sites load the file consistently.
2. Update `actions/work.md:407` and `actions/work.md:658` to explicitly normalize via the Schema Read Contract (matching the narrative at `actions/work.md:180`), and the file becomes truly unreachable for non-canonical domains. Then it can be removed alongside the `CLAUDE.md:144` mention.
3. Load it directly from a non-domain code path (the way `actions/code-review.md:140` loads `crew-members/security.md`).

### `decisions/imported-specs/2026-04-17_improve-weekly-diff-skill.md`

**Why it looks dead:** No file outside this audit references it. The only repo-wide grep hits on `improve-weekly-diff-skill` are this report itself, and the `2026-04-17` date string elsewhere is just metadata (ADR-012 frontmatter, topic-index `updated:` field) rather than a citation of this file.

**Plausible non-dead use:** imported specs are typically preserved as evidence behind specific ADRs. If a future ADR needs to cite this spec, it can. Until then, it's an orphaned evidence file.

**Sibling file is NOT dead — `decisions/imported-specs/2026-04-12_close-gaps-in-interview.md`:** an earlier draft of this report grouped both imported specs together as probable-dead. That was wrong — `decisions/records/adr-012-interview-v2-gap-closure.md` cites the 2026-04-12 spec three separate times (frontmatter `sources:` at line 8, in-body reference at line 34, References section at line 85). The 2026-04-12 file is a live ADR-evidence document; only the 2026-04-17 file is orphaned.

**Incidental finding while verifying:** `decisions/records/adr-012-interview-v2-gap-closure.md:86` links to `decisions/imported-specs/2026-04-16_expand-skill-do-work-interview.md`, which does not exist in the directory. Broken link — separate issue, worth fixing in a follow-up.

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
- **11 of 12 crew-members have a clear inbound load path:**
  - Always-loaded: `general.md`, `karpathy.md` (`actions/work.md:510–511`).
  - Action-specific direct loads: `security.md` (code-review), `ui-design.md` (ui-review), `interviewer.md` (interview), `approach-directives.md` (work, parallel-REQ mode), `caveman.md` (work, when `caveman` frontmatter set), `debugging.md` (work, on review-fail retry + after 2+ test failures), `testing.md` (work, when `tdd: true` or `domain: testing`).
  - Domain-dispatched: `backend.md`, `frontend.md`, `ui-design.md` are loaded by `actions/work.md:512` when the REQ's `domain` matches the canonical enum, and the same `crew-members/[domain].md` pattern fires again at `actions/work.md:407` (Route C planning) and `actions/work.md:658` (review-work spawn).
  - The one with an ambiguous inbound path — `performance.md` — is detailed in the **Probable** section above.
- **The one interview template** (`interviews/work-operating-model.md`) is referenced from `actions/interview.md` and `CLAUDE.md`.
- **The full `decisions/records/` directory** (12 ADRs) is the corpus the `architecture-decisions-log_create-or-expand` prompt manages; ADR-001 through ADR-012 are sequenced and linked from the topic indexes (`decisions/topics/_index_*.md`).

---

## Method Notes

- "Reference" means a literal string match for the file's path or basename inside another tracked file. Dynamic dispatch was resolved by enumerating every site that could fire the dispatch and reading what each site actually does — not just one site. A first pass missed the secondary `crew-members/[domain].md` load sites at `actions/work.md:407` (Route C) and `actions/work.md:658` (review-work spawn); both treat `crew-members/[domain].md` as a file-existence check rather than going through the Schema Read Contract's normalize-and-warn step. That's why `performance.md` is **Probable** and not **Certain** — the contract narrative claims universal enforcement, but two of the three read sites don't restate it.
- Documentation describing a file (e.g., `CLAUDE.md`'s file inventory listing every action) was not counted as a "reference" — that's catalog text, not a dispatch. A reference had to be a path or basename appearing in an executable context (a router table, a load instruction, a "see also" link a reader would follow).
- Per-ADR evidence files (the `decisions/imported-specs/` directory) require checking the full `decisions/records/` ADR corpus, not just the dispatchers. A first pass missed that ADR-012 cites `2026-04-12_close-gaps-in-interview.md` three times; the audit now scans the ADR corpus for citations before classifying an imported spec as orphaned.
- This is a static analysis. A file flagged as dead could still be loaded by a user pasting it into an LLM directly or by a feature added after this audit. Re-run the audit when the routing table, the Schema Read Contract, or the crew-loading mechanism changes.
