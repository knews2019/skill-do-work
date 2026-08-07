# Shell-Logic-in-Prose Census — `actions/` and `prompts/`

**Ran:** 2026-08-05 · **Scope:** all 45 files in `actions/`, all 18 in `prompts/` · **Read-only.**

**This document has been trimmed to its durable findings.** It originally carried a 169-row table scoring every
mechanical step in both directories FULL/PARTIAL/NONE against the shipped tooling, with 415 line-number citations. That
table is gone, and the reason is the finding it became: **it started lying within hours.** Merging `origin/main` added
~25 lines to `actions/work-reference.md`, and every citation past the insertion point shifted — the table claimed
L790 / L794–809 / L821 while the cited text had moved to L815 / L822–826 / L846. A coordinate-based index of files
under active development is a liability, not an asset, so the coordinates were removed rather than maintained.

The full table remains in git history at commit `2cf7c0a` for anyone who wants the original evidence. Nothing below
depends on it.

---

## What the census found, and what happened to it

### Finding 1 — `frontmatter.go` had no CLI surface. **Fixed.**

The board tool registered seven subcommands and **none took a file-and-field argument**, so `splitFrontmatter` /
`parseFrontmatterFields` / `lenientFrontmatterFields` were unreachable from any action file. Every frontmatter read in
`actions/` was therefore a hand reimplementation *by construction*, not by oversight — roughly 95 of them, against
three parsers that already shipped (the Go one, plus awk inside `tools/checks/record-commit-hash.sh` and
`tools/checks/blanked-req-scan.sh`).

The consequence was measurable rather than theoretical: the `status` vocabulary alone is read at ~35 prose sites, and
five separate Red Flags across `actions/` document the *same* resulting bug — filtering on the literal `completed` and
silently dropping `completed-with-issues`. Five documented instances of one bug class was the census's strongest
evidence that the contract was fine and only its enforcement was 35 hand copies.

**Delivered by REQ-112** (v0.175.0): `queue-kanban frontmatter get <file> <field> [--normalize] [--in-set …]`.

### Finding 2 — seven of nine Schema Read Contract fields had no normalizer. **Fixed.**

The Schema Read Contract defines nine enum-or-boolean fields with alias maps, canonical enums, and documented
defaults. Only two had a mechanical implementation anywhere in the repo: `normalizeStatus` and
`normalizeTestingStatus`. `domain` was read verbatim, so `domain: back-end` silently mis-selected the crew file a REQ
meant to load, with nothing anywhere to catch it.

**Delivered by REQ-111** (v0.174.15): all seven remaining fields normalize through one table, an unrecognized value
falls back to its documented default *and* reports itself unrecognized, and an absent field stays absent.

### Finding 3 — `prompts/` is almost entirely clean. **No action needed.**

17 of 18 prompt files contain zero shell commands, zero frontmatter reads, and zero output parsing. The single
exception is `architecture-decisions-log_create-or-expand.md`, which gates on `git status --porcelain` and a
branch-name check, and audits `[[wiki-link]]` targets. Worth knowing so nobody re-audits the directory expecting to
find something.

### Findings 4–6 — residual extraction candidates. **Carried by REQ-114.**

The census ranked five extraction candidates by execution frequency × bug risk. The top two became REQ-111 and
REQ-112 above. The remaining three — consolidating the merge-commit-aware diff idiom, the uncommitted-changes
inventory plus REQ association, and writer-label claim classification — are carried by **REQ-114**, restated as
greps rather than line numbers so they cannot decay the way this table did. Candidate B was separately approved and
delivered as REQ-121; Candidates A and C remain unapproved, separate work.

---

## The two structural lessons worth keeping

**A finding about the absence of a mechanism does not rot; a finding about a location does.** Findings 1 and 2 were
still exactly true after a merge that invalidated a quarter of the table's citations, because "there is no CLI
surface" and "seven fields have no normalizer" are claims about what does not exist. That distinction is the one
transferable result here, and it is why REQ-114 describes its candidates by what to grep for.

**An audit's coverage baseline has to be complete or it understates coverage everywhere at once.** The original
inventory treated `tools/checks/` and `queue-kanban` as the whole toolbox and omitted `hooks/*.sh` and
`tools/do-work-update.sh` — both shipped paths. A PR review caught two rows understated as a result (the memory
capture/redaction/ledger mechanics all ship in `hooks/memory-stop-capture.sh`; pipeline state-parsing ships in
`hooks/pipeline-guard.sh`), and applying this repo's grep-the-primitive rule to the exposed class found a third.
An incomplete baseline is the one error that actively misdirects the extraction work an audit exists to inform.

---

## Method, and its limit

Coverage was verified against source — `tools/checks/*.sh` headers, the board tool's subcommand dispatch, its verify
probe constants, and its field readers — never from filenames. Every claim in the original table cited lines that were
actually read.

**But read depth was not uniform, and that was the audit's own soft spot.** 14 files were read end-to-end:
`work.md`, `work-reference.md`, `forensics.md`, `cleanup.md`, `version.md`, `commit.md`, `inspect.md`, `board.md`,
`tidy-repo.md`, `stray-check.md`, `ai-report.md`, `review-work.md`, `memory-value.md`, `validate-feedback.md`. The
other 31 action files and all 18 prompt files were scanned with a keyword pattern (backticked shell commands, `glob`,
`frontmatter`, `scan`/`parse`/`compare`/`filter`/`count`) and only matching lines were read. A mechanic phrased
without any of those tokens would have been missed there. Anyone re-running this should read the remaining 49 in full.
