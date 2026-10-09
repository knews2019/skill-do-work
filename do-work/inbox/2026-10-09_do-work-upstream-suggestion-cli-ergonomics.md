# Upstream suggestion for `knews2019/skill-do-work` — four do-work-cli additions to replace hand-rolled python heredocs: auto-manifest finalize, frontmatter set, worktree lifecycle, capture-files init

**How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
consumer repos it was authored in. Observed against **v0.305.84**; all line numbers are from that tag. Paths
are relative to `skills/do-work/`.

---

## Request

Four steps of the run and capture pipeline are written as prose that a session must turn into shell by
hand. Sessions do that with python heredocs (assert an anchor, then `.replace()`), hand-written manifest
generators and ad hoc `git worktree` chains. The same glue is rebuilt in session after session, and three
of these procedures now survive only as per-machine memory notes ("finalization-manifest",
"single-req-run", "manual-release-manifest"). Please add four CLI commands. Each one stands alone, so each
can be captured and shipped as its own REQ:

1. **`finalize --auto-manifest REQ-N`**: build the finalization manifest's mechanical fields from the REQ and
   the planner, preflight the tree, print a short outcome table.
2. **`frontmatter set <file|REQ-N> <field> (<value> | --at now)`** next to the existing `frontmatter get`, plus
   **`req append-section REQ-N --section <name> --from <file>`**.
3. **`worktree new|status|merge|cleanup REQ-N`**: the builder worktree lifecycle that
   `actions/work-reference.md` already specifies step by step.
4. **`capture-files init <dir>`**: write a filled `capture-files` manifest skeleton and the UR/REQ payload
   templates.

No new frontmatter field, no new status, no change to who judges what. Every command wraps a procedure the
actions already define; none invents a rule.

Dedupe: no entry from 0.305.60 to 0.305.87 adds any of the four, and the upstream queue is empty.
0.305.87 moved the worktree dispatch text (naming, hand-back merge, cleanup) into
`actions/fan-out-reference.md` without changing a rule; the `work-reference.md` lines cited below are its
v0.305.84 home. The archived request that gave frontmatter a CLI read surface (UR-021 / REQ-112, "Give
frontmatter.go a CLI surface so prose can stop reimplementing it") put a `set` verb explicitly out of
scope to avoid a new write surface. Item 2 asks to revisit that: the write happens anyway, by hand, with no
validation.

## What happened, in consumer repos (evidence)

Counts come from 30 days of session transcripts (2026-09-09 to 2026-10-09) across four repos: three
consumer repos and the suite's own maintainer checkout. A second pass re-counted each figure; the numbers
below are the re-counted ones.

| Procedure | Commands | Sessions | Distinct days |
| --- | --- | --- | --- |
| `advance ... --finalization-manifest <hand-built json>` | 171 | ~30 | 15 |
| python heredoc that stamps `review_at:` / `integration_at:` into a REQ file | 112 | ~24 | 13 |
| `git worktree add` for a builder by hand | 85 | ~31 | 14 |
| `date -u +%Y-%m-%dT...` for a stamp, against `do-work-cli now` | 356 vs 40 | 75 vs 10 | |

Dated snippets, anonymised:

- 2026-09-27, one consumer repo: a session wrote `make_manifest.py` into its scratch directory
  (`import json,hashlib,sys,os,re,datetime,shutil`). Hand-written finalize helpers under four names
  (`make_manifest.py`, `finalize.py`, `fin.sh`, `adv.sh`) appear in two sessions; most sessions inline the
  same logic instead.
- 2026-10-02, the suite's own checkout: a python heredoc set `anchor="\n## Decisions\n"`, asserted
  `s.count(anchor)==1`, then ran
  `s.replace("\nintegration_at:","\nreview_at: 2026-10-02T18:47:04Z\nintegration_at:",1)` to stamp one field.
- 2026-10-03, the same consumer repo, finalizing a REQ: `grep '"version"' package.json | head -1 && test -z
  "$(git diff --cached --name-only)" && do-work-cli.sh ... advance REQ-N ... --finalization-manifest
  .../finalization.json > .../finalize-result.json; python3 -c "import json;..."` to read the result back.
- 2026-10-05, the suite's own checkout: the first `capture-files --dry-run` refused with
  `CAPTURE-RAW-INPUT-NOT-CONTAINED` because the UR payload copied the four-backtick `text` fence from the
  reference example. The command derives a bare fence from the raw bytes.
- 2026-10-07, the same consumer repo: `git worktree add -b $name $WT/$name main && ln -s $PWD/node_modules
  $WT/$name/node_modules && ln -s $PWD/env.vars $WT/$name/env.vars`.
- 2026-10-09, a consumer repo outside the four counted above: a session grepped `json:"` in `capture_files.go`
  to recover the manifest schema. The same hunt appears in two other repos. The tags it wanted are in a
  different file.

The per-machine notes say why. One reads: "There is no dry-run and no validator ... so you iterate on typed
refusals." Another: "the manifest schema lives only in Go ... and no example exists on disk, so re-deriving
it cost ~10 reads." A third lists `commit_paths` as "submit a first guess and read the refusal, which names
the exact missing paths."

The finalize, stamp and worktree items recur strongly. The capture-files item has the thinnest count of
the four; it is included because its schema hunt and the fence mismatch above are concrete.

## Where the behaviour lives today

**Item 1, finalization manifest.**

- `actions/work-reference.md:810`: "The action judges and authors one strict finalization manifest: the
  exact request id and path, the expected request and checkpoint digests, the terminal transition and
  timestamp, the writer, the exact implementation/lifecycle/release allowlist, the commit message, any
  release payload, and the provenance mode." Judged and mechanical fields sit in one hand-built file.
- `actions/work.md:475`: "Run the exact `advance` continuation with the selected request path and the single
  action-authored finalization manifest".
- `tools/do-work-cli/internal/lifecycleadvance/advance_commands.go:263-265`: the missing-evidence hint is the
  literal placeholder `<action-authored-finalization-manifest>`.
- `tools/do-work-cli/internal/finalization/finalization_types.go:32-49`: the `Manifest` struct, the only
  schema. No example manifest ships in `actions/` or `docs/`.
- `tools/do-work-cli/internal/finalization/finalization_commands.go:318-334`: `finalize` accepts only
  `--manifest <file>`. No dry-run, no validate mode.
- `tools/do-work-cli/internal/finalization/finalization_prepare.go:158-174`: the finalizer already computes
  the required commit paths (`statePlan.TargetPaths` plus release postimages) and refuses with
  "commit_paths omits planned lifecycle or release targets". The CLI knows the set; it only reports it after
  a refusal.
- `tools/do-work-cli/internal/resultmodel/result_model.go:648-662`: refused already exits 1. Keep that.

**Item 2, frontmatter stamps and REQ sections.**

- `tools/do-work-cli/internal/corehelpers/commands.go:77-79`: `handleFrontmatter` accepts only `get`
  ("usage: frontmatter get <file> <field> [--normalize] [--in-set SET]").
- `actions/work.md:317` ("stamp `integration_at: <now>` ... only if the field is absent") and
  `actions/work.md:383` (`review_at`, `remediation_at`, `re_review_at`, same absent-only rule).
- `actions/work-reference.md:59` (Timestamp rule: UTC `YYYY-MM-DDTHH:MM:SSZ`) and `:73` ("**Stamps are
  append-only.**"). Both are rules a CLI can enforce and a heredoc cannot.
- `tools/replace-text-section.sh:10`: `replace-section` works on marker-delimited managed blocks
  (`--begin-marker` / `--end-marker`), not on REQ `## Testing` or `## Review` sections.
- `actions/sample-archived-req.md:21-111`: the de facto section order (What, Red-Green Proof, ..., Testing,
  Review, Lessons Learned, Orientation). No command checks it, and an out-of-order write is refused later by
  `advance` ("later lifecycle evidence exists before ...",
  `tools/do-work-cli/internal/lifecycleadvance/advance_commands.go:328`).

**Item 3, builder worktree lifecycle.**

- `actions/work-reference.md:416` (Naming: `worktree-agent-REQ-NNN-<suffix>`, numeric collision suffix),
  `:418` (the created name is the operative name), `:420` (worktrees live outside the repo).
- `actions/work-reference.md:426` ("Integrate with `git merge --no-ff <branch>`") and `:430-436` (hand-back
  sequence: settle the index, capture `<pre>`, queue guard
  `git diff --name-only <pre>...<operative_name> -- do-work/`, merge, capture `<merge_hash>`).
- `actions/work-reference.md:445`: cleanup is `git worktree remove` (no `--force`), `git branch -d` from the
  integration branch, `git worktree prune`.
- Nothing in `actions/` or `docs/` covers dependency directories. A JS consumer needs `node_modules` and its
  env file in the builder tree, so every session adds its own `ln -s` lines.

**Item 4, capture-files manifest.**

- `actions/capture.md:228`: prepare "one strict `capture-files` JSON manifest as regular payload files"; the
  shape is described in prose only. `:241` is the only command line shown.
- `actions/capture-reference.md:9` and `:84`: the templates define payload bytes; the manifest schema is a
  link to Go source.
- `tools/do-work-cli/internal/publication/publication_types.go:19-26` (`Manifest`) and `:53-60`
  (`CaptureManifest`): the actual JSON tags.
- `actions/capture-reference.md:191`: the UR example uses `> ````text`, while
  `tools/do-work-cli/internal/publication/publication_manifest.go:107-123` builds a bare fence of
  max(longest backtick run + 1, 3) with no info string, and
  `tools/do-work-cli/internal/publication/capture_files.go:89-90` refuses anything else.

## Proposed direction

**1. `finalize --auto-manifest REQ-N`.** Keep the judged fields with the action, as
`actions/work-reference.md:810` requires. Take them as inputs: `--transition`, `--terminal-status`,
`--message-file`, `--provenance supplied_commit --implementation-hash <merge>` or `--provenance
primary_commit`, optional `--release-manifest`. Fill the rest: request path from the working REQ,
`expected_request_sha256` and `expected_checkpoint_sha256` from the live files, `completed_at` from `now`,
writer label, and `commit_paths` as the planner's required set from `finalization_prepare.go:158-174` plus any
`--extra-path`. Preflight before writing anything: the project version still matches the release manifest's
old version, the index is clean, nothing staged lies outside `commit_paths`. `--emit <path>` writes the
manifest and stops, which is the validate mode the notes ask for. On success or refusal print one table:
outcome, phase, archived path, commit hash, then one line per finding (code and evidence) and per gate. Exit
codes stay as `result_model.go:648-662` defines them.

**2. `frontmatter set` and `req append-section`.** `frontmatter set <file|REQ-N> <field> <value>`, or `--at
now` in place of `<value>` for `*_at` fields. It refuses to overwrite an existing `*_at` field (the
append-only rule), formats the instant per the Timestamp rule, and keeps the Frontmatter Quoting contract.
`req append-section REQ-N --section Testing --from <file>` inserts the section in canonical order, refuses a
second copy of a section that already exists, and is a no-op with exit 0 when the same bytes are already
there. Both resolve `REQ-N` to its one file in `working/` or the queue and refuse on ambiguity.

**3. `worktree new|status|merge|cleanup REQ-N`.** `new` derives the name and collision suffix exactly as
`actions/work-reference.md:416` says, creates the branch and worktree from the integration branch outside the
repo (`:420`), links the paths listed in an optional per-repo config (for example `do-work/worktree-links`,
a new file with one repo-relative path per line such as `node_modules` and `env.vars`), writes the brief
skeleton, and prints the operative name. `status` lists each `worktree-agent-REQ-*` worktree with
ahead/behind against the integration branch, dirty or clean, and last-commit age. `merge` runs the
`:430-436` sequence: clean index, `<pre>` capture, queue guard, `git merge --no-ff --no-commit`, commit with
a `[REQ-N]` message, then prints `<pre>` and `<merge_hash>`. `cleanup` removes the links, runs
`git worktree remove` and `git branch -d` from the integration branch, then `git worktree prune`. Never
`--force`, never `-D`; a refusal is reported and stops.

**4. `capture-files init <dir> [--raw-input <file>]`.** Writes `manifest.json` with every required key
filled with a placeholder or a scanned value (next REQ id and reservation path, from the read-only scan
that `actions/capture.md:78-80` describes in prose: "The scan proposes IDs but writes nothing"), plus UR
and REQ payload templates. With `--raw-input`, the UR template already carries the byte-derived verbatim
block from `containedOutsideBytes`. Fix the `actions/capture-reference.md:191-193` example in the same
change: for that one-line input the command derives a three-backtick fence with no info string, so the
four-backtick `text` opening line and the four-backtick closing line both become three backticks. If a new
subcommand is too much, a `capture-files --example` that prints one valid manifest covers most of the gap.

Once each command ships, change the action line that describes the manual procedure to name the command,
and tell the memory-note owners the note can be deleted.

## Acceptance check

- 1: on a fixture REQ ready for finalization, `finalize --auto-manifest REQ-N --emit m.json` writes a manifest
  that `finalize --manifest m.json` accepts unchanged. A staged unrelated path makes it refuse with exit 1
  before any journal is written. The judged fields are never invented: omitting `--message-file` refuses.
- 2: `frontmatter set REQ-N review_at --at now` adds the field once; a second call exits non-zero and leaves
  the bytes identical. `req append-section` run twice with the same file leaves one section; a different
  file for an existing section refuses. A `## Testing` appended to a REQ that already has `## Review` lands
  before it.
- 3: `worktree new` then `cleanup` on a fixture leaves no worktree, branch or link behind. `merge` refuses
  when the builder branch committed anything under `do-work/`. A dirty worktree makes `cleanup` refuse.
- 4: `capture-files init` output plus filled placeholders passes `capture-files --dry-run` on the first try,
  including a raw input that contains a triple-backtick run.

## Out of scope

- A gate runner with failure-only summaries and flake repeats. It was seen in transcripts but with one clear
  example; it can be its own request later.
- Changing who judges the commit message, release payload or transition. Item 1 fills mechanical fields only.
- Changing worktree naming, merge, or cleanup rules. Item 3 wraps them.
- Removing the hand path. The prose procedures stay valid for harnesses without the new commands.
