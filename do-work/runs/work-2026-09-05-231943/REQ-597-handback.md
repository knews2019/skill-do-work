# REQ-597 hand-back — correct the stale claims in the prescribed-shell guide and its two callers

Branch `worktree-agent-REQ-597-guide-and-callers`, head `7df648828a2e71f26a83967fa4a57f2da9f3a56c`, three commits on top of the pre-flight commit `5a7d70e`:

- `a1e652f` stage 1 — `skills/do-work-toolbox/actions/inspect.md`: the two prescribed blocks and their lead-ins
- `6913dc4` stage 2 — `skills/do-work/docs/prescribed-shell-primitives.md`: sixteen claims in five sections, plus line 24 and two phrases on line 137
- `7df6488` stage 3 — `skills/do-work/actions/commit.md` and `inspect.md`: the association prose

Files changed across the branch: `skills/do-work-toolbox/actions/inspect.md`, `skills/do-work/docs/prescribed-shell-primitives.md`, `skills/do-work/actions/commit.md`. No Go source, no `_dev/tests/`, no changelog entry (see F13). Nothing in the main checkout was edited, staged or committed; this file is the only write there and is not staged.

Every sentence below was derived from the cited code and checked by running the command it describes. Where the verification's `suggested_text` matched what was measured it was reused; where it did not, it was dropped, and the drops are listed under "Drafts not used".

## 1. The broken commands: before and after

`inspect.md` lines 64 and 114 prescribed `protected-inventory.sh start "$(git rev-parse --show-toplevel)" do-work-inspect-secret-quarantine` and the same shape for `associate`. `scripts/protected-inventory.sh:6` is a bare `exec` that forwards `"$@"` after the command token and translates nothing; `internal/corehelpers/inventory.go:339-353` accepts only `--dry-run` and `--quarantine-name` after the mode and returns `usageResult("unknown option …")` for anything else.

Before, run from the fixture root (`scratchpad/req597-inspect/build-fixture.sh`, outputs in `before-start.{out,err}` and `before-associate.{out,err}`):

```
start  "<root>" do-work-inspect-secret-quarantine   -> exit 2, stdout "finding HELPER-USAGE [error]: unknown option <root>", stderr empty
associate "<root>" do-work-inspect-secret-quarantine -> exit 2, identical finding
ls "$(git rev-parse --git-path do-work-inspect-secret-quarantine)" -> No such file
```

`inspect.md:117` reads exit 2 as its skip condition, so as shipped the action never associated a file with a REQ in any project.

After, the blocks as committed (extracted by grep, run against a fresh fixture with `<skill-root>` = `skills/do-work-toolbox`; stage 3 re-ran them on its own fixture, `scratchpad/req597-callers/R-inspect-*.{out,err}`):

```
<skill-root>/../do-work/scripts/protected-inventory.sh start --quarantine-name do-work-inspect-secret-quarantine
  -> exit 0; rows XD config/app.env, M src/alpha.txt … X .env.local, X do-work/.req-reservations/123
  -> .git/do-work-inspect-secret-quarantine created, mode 0600, holding the two X paths
<skill-root>/../do-work/scripts/protected-inventory.sh associate --quarantine-name do-work-inspect-secret-quarantine
  -> exit 0; rows REQ-002 config/app.env, REQ-101 src/alpha.txt, REQ-002 src/beta.txt, REQ-103 src/delta.txt, REQ-105 src/eta.txt, REQ-104 src/zeta.txt
```

`commit.md`'s two blocks (`<skill-root>/scripts/protected-inventory.sh start` and `… associate`, no arguments) were correct as shipped and were left unchanged; run verbatim on the same fixture (`R-commit-*`) they exit 0 with the same rows and write the default `do-work-commit-secret-quarantine` (`inventory.go:337`).

## 2. Stage 1 — inspect.md blocks and lead-ins (commit a1e652f)

| Sentence | Code | Fixture |
|---|---|---|
| Line 64 and 114 blocks: `--quarantine-name do-work-inspect-secret-quarantine`, no root argument | `inventory.go:339-353`; launcher line 6 | Before/after above |
| Line 61: "from the project root … reads the repository root from the current directory and rejects a root argument; run from a subdirectory, `start` prints the same rows but Step 3's `associate` exits 2 as if `do-work/` were missing" | `commandruntime/command_runtime.go:103-107` (root defaults to `os.Getwd()`); `inventory.go:351-352` (positional rejected); `inventory.go:203-207` (NO-DO-WORK-DIR, exit 2) whose text the shim loop at `:445-456` erases | `repo-sub/src`: start exit 0 same four rows; associate exit 2, both streams empty. `repo-nodowork`: associate exit 2, empty |
| Line 61: "Without `--quarantine-name` the wrapper writes `commit`'s file, `do-work-commit-secret-quarantine`" | `inventory.go:337` | `repo-commit`: bare `start` created `.git/do-work-commit-secret-quarantine` |
| Line 111: "`associate` reads the quarantine file `start` wrote, and exits 2 with a `HELPER-USAGE` finding when that file is missing" | `inventory.go:412-415` | `repo-nostart`: exit 2, `finding HELPER-USAGE [error]: protected inventory has not been started with a regular Git-private quarantine file` |
| Retained "worktree-safe" | `inventory.go:358` uses `git rev-parse --git-path` | `repo-wt-linked` (a `git worktree add` checkout): quarantine written under `.git/worktrees/<name>/`, main `.git` untouched; associate exit 0 |

Draft not used: the verification's "the command resolves the repository itself and takes no repository-root argument". The second half is true; the first is false (the root is the current directory, not resolved).

## 3. Stage 2 — prescribed-shell-primitives.md (commit 6913dc4)

Fixtures under `scratchpad/req597-guide/` (`dwcli` built from the worktree with go1.26.1; `pi-repo`, `timing/repo`, `gitfx/repo`, `rename/`, `download/`, `shot/`, `portfolio/repo`; `apply-edits.py` holds the exact old→new strings).

**Shipped executable homes, line 24.** "`scripts/protected-inventory.sh` translates nothing: it is a bare `exec` that forwards its arguments unchanged, so a positional after the mode reaches the command as written and is rejected with `unknown option`, and a caller spells the flag itself (`--quarantine-name <name>`). What that launcher adds is `DO_WORK_COMPATIBILITY_SHIM=1`, which selects the `<tag>\t<path>` output …" — `protected-inventory.sh:6`; `inventory.go:340-352`; `:44-63,445`. `strace -f -e trace=execve` shows argv `[…, "protected-inventory", "start", "<root>", "qn"]` reaching the binary unchanged; the shim-less binary prints `finding INVENTORY-M [info]` blocks instead of tab rows.

**Lifecycle timing.**

- Clock: "Neither accepts a duration — elapsed seconds are always derived by the command … `record-timing-event` stamps the end itself and, given no start, chains the start to the previous event's end (or to its own end on an empty stream); `--started-at <RFC3339>` pins the start … whose POSIX floor is `date -u`." — `timing_commands.go:188-193,211`; `lifecycle_timing.go:135,170-172,280-282,287,322-329`; `work.md:284`; `work-reference.md:67-68`. Measured: `--elapsed 5` → exit 2 `TIMING-USAGE unknown option`; explicit `--started-at` stored verbatim with `elapsed_source wall_clock_difference`; next event with no start chained to the previous `ended_at`; first event on an empty stream has `started_at == ended_at`.
- Wrapper: "hands the child the CLI's own stderr handle for both its stdout and its stderr — a file descriptor, not a pipe — … redirect the CLI's stdout and you capture the result and none of the child's output" — `lifecycle_timing.go:167-168`; `main.go:47` vs `:80`. Measured: `>out.txt 2>err.txt` puts only the rendered result in out.txt and both child lines in err.txt; `/proc/$$/fd/1` and `fd/2` in the child both point at the file; exit 3 / 143 / 127 fidelity measured.
- Evidence: "A command reaches the stream as its executable's base name plus an argv token count … That redaction covers the command and nothing else. `--operation` is free caller text: it reaches the stream as written, and the slowest stage's and the slowest command's operation is printed verbatim into the folded `## Timing` section … `--agent` and `--revision` are free text too and reach the stream, but not the folded section. All three are only stripped of control characters and `|`, collapsed to one line, and cut to 120 characters; that is bounding, not redaction." — `lifecycle_timing.go:702-711,307,298,304-305,242,483-509,382-386,716-732,104`. Measured: an operation containing `/home/alice/.ssh/id_rsa token ghp_ABC123SECRET` was written verbatim into the request file's `Slowest stage:` line; `--agent`/`--revision` appear in the `.jsonl` rows and zero times in the folded file; `printf 'a|b\tc\x01d'` stored as `ab cd`, a 121-character agent cut to 120.
- Fold: "turns one request's stream — a run holds one per REQ — into a single `## Timing` section and deletes that stream, leaving a sibling request's stream in place" — `lifecycle_timing.go:637-649,351-391`. Measured: after folding REQ-006, `runA/` still holds `REQ-007.jsonl`.

**Merge-aware commit diff.** "It resolves the argument with `git rev-parse --verify <commit>^{commit}`, tests for a second parent with `git rev-parse --verify -q <commit>^2` — both handed to `git` as argv, with no shell in between to quote for — and emits `git show --first-parent -m` for a merge or ordinary `git show` otherwise." — `git_helpers.go:99,104,105-106`; `show-commit-diff.sh:7-9`. strace shows exactly three git execs, no shell.

**Commit file listing.**

- "start from:" (was "use:") — no shipped file runs the printed form; the three Go readers each change the flags (`gate_evidence.go:377`, `git_transaction.go:1375`, `finalization_apply.go:545`).
- "This prints no commit header and no message, so message text cannot become a phantom filename, and `-m` is what makes it list a merge's paths at all — without it a merge prints nothing." — every fixture commit message contains `phantom.txt`; clean merge without `-m` printed 0 bytes.
- New: "the repository's first commit prints nothing without `--root`, and a filename carrying a double quote, a backslash, a control character or (under the default `core.quotePath`) a non-ASCII byte comes back C-quoted without `-z`; a plain space does not trigger quoting." — measured on `quo"te.txt`, `back\slash.txt`, `ta<TAB>b.txt`, `café.txt`, `sp ace.txt`.
- New: "On a merge, `-m` prints a path once per parent that changed it, so a path changed on both sides appears twice — de-duplicate before counting." — conflicted merge printed `shared.txt` twice.
- "`git show --name-only --format=` is acceptable for a commit known to have one parent … never point it at a merge, where it prints the combined diff — nothing for a conflict-free merge, only the conflict-resolved paths otherwise" — measured on both merges; the one shipped consumer (`bkb_init.go:733`) reads a commit it just created.
- "No shipped file runs the form above as printed; the three Go readers each carry what their own input needs: gate evidence walks every commit after a green gate with `-m -z`, the Git transaction reads back the commit it just made with `--root -z` and no `-m`, and finalization matches candidate commits against a recorded diff digest with neither." — the three sites above; `grep -rn diff-tree skills/` finds only the guide.

**Verified exact publication.**

- Opening: "the publishing step must be able to say what actually landed there; a rename's or link's exit status alone is not that answer, and shell and syscall disagree about what it means." — only `publication.go:178-183` verifies after the fact; `commands.go:888-891`, `mutation.go:190-204,270` do not.
- Shell sentence kept and scoped: "In hand-written shell, `ln` and `mv` treat a directory standing in the destination's place as a container rather than a collision …" (the phrase is pinned by `prescribed-shell-canonicalization.sh:127`).
- New: "None of the shipped publications below runs `ln` or `mv`. Each publishes through Go's `Rename` or `Link`, which cannot nest anything: `link(2)` refuses every existing destination, and `rename(2)` refuses a directory — but silently replaces an occupying regular file, so a rename is only as safe as the check made on the occupant before it." — the only shell `mv` under `skills/` is the installer bootstrap heredoc into a fresh `mktemp -d`.
- New: "One shipped publication then proves the result: `capture-screenshot` links its stage into place, stats both names, and compares them with `os.SameFile`. The three others below make a misplaced write impossible in advance instead …" — `publication.go:132-133,178,180-183`; `commands.go:869-871,888`; `portfolio.go:198-214` + `mutation.go:196`; `report_image.go:147-148,193,218` → `mutation.go:228-270`. strace of `capture-screenshot`: `linkat` then `newfstatat` on both names; of `publish-portfolio-summary`: `linkat`/`renameat` with no stat or read afterwards.
- "What a helper does *about* an occupied destination is its own policy" (was "a nested write").

**Portfolio summary publication.**

- "The canonical command reads that source once, requires it to be a regular file, and writes every output from those same bytes: each output is staged in a private dot-file beside its own destination … written, synced and closed, then linked or renamed into place. Nothing reads a stage back; what guards a replacement is that the canonical file is still, by inode and by bytes, the regular file the command read just before staging." — `portfolio.go:30,34,123,147`; `mutation.go:140-141,146-157,162-163,171-179,184-190,196`. strace: exactly one `O_RDONLY` open of the source, zero `O_RDONLY` opens of any `.publishing-` name, the canonical re-opened immediately before `renameat`.
- "`--with-snapshot` first publishes the snapshot with an exclusive link onto a candidate it checked free, advancing an occupied candidate to the next numeric suffix (`-2`, `-3`, …), and only then atomically replaces the canonical file from the same bytes." — measured: directory-occupied candidate → `snapB-2.md`; file-occupied → `snapD-2.md`.
- "there is no script to fall back to — neither skill ships one for this publication." — `find skills -iname '*portfolio*'` → `portfolio.go` only.
- "Each publication answers the … check by making a misplaced write impossible rather than by re-reading what it wrote: an occupied snapshot candidate, directory or not, advances to the next suffix; a canonical path occupied by anything but a regular file fails closed; and neither leaves a private file nested inside the occupant." — measured with a directory, a symlink and a fifo in the canonical's place: `PORTFOLIO-CANONICAL-UNSAFE`, exit 1, occupant untouched.

**Report image batch publication, line 137 (beyond the sixteen, same defect class).** "— there is no shipped script for this publication to orchestrate by hand." and "there is no script to fall back to." — `find skills -iname '*report*image*'` → `report_image.go`, `report_image_process.go` only.

### Rename and link measurements (scratchpad/req597-guide/rename)

Fresh case each time: `.stage` holds `PAYLOAD`, `occupied_file` holds `OLD`, `occupied_dir` is an empty directory.

| Operation | `rename(2)` / `link(2)` (python `os.rename`/`os.link`) | Go `os.Rename`/`os.Link` and `Root.*` | shell `mv`/`ln` |
|---|---|---|---|
| rename over a regular file | returned 0, file now `PAYLOAD` | nil, replaced | exit 0, replaced |
| rename over a directory | `EISDIR`, nothing written | `os.Rename` "file exists" (its own pre-check), `Root.Rename` "is a directory"; directory 0 entries | exit 0, `.stage` moved **inside** the directory |
| link over a regular file | `EEXIST` | "file exists" | exit 1 "File exists" |
| link over a directory | `EEXIST`, directory 0 entries | "file exists", 0 entries | exit 0, a hard link created **inside** the directory |

Derived: `mv` and `ln` nest into an occupying directory and exit 0; `link(2)` refuses every existing destination; `rename(2)` refuses a directory and silently replaces a regular file. Because `os.Rename` and `rename(2)` report different errnos for the directory case, the guide names no errno.

### Drafts not used in stage 2

- "Go's `os.Rename` and `os.Link` refuse a destination that already exists" — false for a regular file under rename.
- "refuse that destination with EEXIST and EISDIR" — no single errno pair is true of every shipped call.
- "compares the canonical against what it read a moment earlier / since the command started" — the read is inside `rootedPublishFile`, immediately before staging.
- "`--agent` and `--revision` reach the folded section" — they reach the stream only.
- "single-parent" for the Git transaction's commit — parent count not verified, adjective dropped.
- The quoting-trigger list that omitted control characters — a tab is C-quoted too.

## 4. Stage 3 — commit.md and inspect.md association prose (commit 7df6488)

Fixture `scratchpad/req597-callers/build-fixture.sh` → `repo` and copies (`repo-malformed`, `repo-nodowork`, `repo-nostart`, `repo-allx`, `repo-clean`, `repo-inspect`, `repo-commit`, `repo-sock`, `repo-badfm`, `not-a-repo`); every run's stdout/stderr saved as `<label>.out`/`.err`; `apply-edits.py` holds the exact old→new strings. The fixture holds REQ files at depth 1, 2 and 3 under `archive/` and `working/`, an archive REQ with status `claimed`, a symlinked REQ file and a symlinked directory holding a REQ, an untracked `.env.local` claimed by REQ-101, a deleted tracked `config/app.env` claimed by REQ-002, an unclaimed `src/orphan.txt`, and three hidden untracked paths under `do-work/`.

Sentences changed, with the old text in the commit diff:

| # | Where | New text (abridged) | Code | Fixture run |
|---|---|---|---|---|
| S1 | commit.md:53 | "Start the protected inventory wrapper from the project root; it owns the worktree-safe run quarantine, takes no `--repo-root`, and uses the current directory as the repository root, so from a subdirectory Step 3's `associate` exits 2 as if `do-work/` were missing:" — the clause "delegates low-level classification to the existing checks" deleted | `command_runtime.go:103-107`; `inventory.go:366,439` call `readInventory`/`handleAssociate` in-process, nothing under `tools/checks/` is run; launcher line 6 puts `"$@"` after the command token | `Y`: `start --repo-root <root>` → exit 2 `unknown option --repo-root`. `F`: associate from `repo/src` → exit 2, both streams empty; `F2`: start from `repo/src` → exit 0, same rows |
| S2 | commit.md:59, inspect.md:67 | "except an untracked hidden file under `do-work/`, which it drops as editor or operating-system metadata, and prints one `<tag>\t<path>` row per remaining file" | `inventory.go:102-104,141-143` | `A`: `do-work/.DS_Store` and `do-work/working/.REQ-101-fixture.md.swp` absent, `do-work/.req-reservations/123` present. `U`: a tracked hidden file, modified → `M` row (predicate needs `??`). `Q`: a hidden file as the only change → `INVENTORY-CLEAN`, exit 1 |
| S3 | commit.md:79 | "; run it from the project root, as in Step 1:" | as S1 | `F` |
| S4 | commit.md:85, inspect.md:117 | "The wrapper moves paths through files …" — "re-derives the repository root" deleted | `command_runtime.go:103-107`; `inventory.go:204` stats `<cwd>/do-work` | `F` |
| S5 | same | "walks `do-work/working/` and `do-work/archive/` to any depth, skipping symlinks, reads the `## Implementation Summary` file list of each `REQ-*.md` there that counts (the next paragraph links what counts)" | `inventory.go:254-256` (`WalkDir` over both roots), `:263-268` (symlink file skipped, symlink dir `SkipDir`), `:269` (name shape), `:281-284` (archive REQ needs a terminal-success status) | `B`: rows for REQ-103 (`archive/` depth 1), REQ-104 (`working/nested/`, depth 2), REQ-105 (`archive/UR-002/deeper/`, depth 3); no row for `src/iota.txt` (REQ-107 behind a symlinked directory), `src/kappa.txt` (REQ-108 a symlinked file), `src/theta.txt` (REQ-106, archive status `claimed`) |
| S6 | same | "prints one `<owner>\t<path>` row only for a candidate a REQ claims; it prints no placeholder for the rest" | `inventory.go:445-453` (only `ASSOCIATION-FOUND` becomes a row); the `-` placeholder is `inventory.go:236`, reached only by `tools/checks/associate-files.sh` | `B`: no row for `src/orphan.txt`; `E`: `associate-files.sh` prints `-\tsrc/orphan.txt` and also `REQ-101\t.env.local`, which is why the action must not substitute it |
| S7 | same | "Exit 2 means nothing was associated. A failure before association starts, such as a wrong invocation, `associate` before `start`, or a directory that is not a Git repository, prints a `HELPER-USAGE` finding that names it; a failure inside association, such as no `do-work/` directory, a REQ file it cannot read, or an unmatched backtick in an `## Implementation Summary` line, prints nothing. Either way, skip REQ tracing …" (inspect.md keeps "the skip condition already stated above" and omits the `associate`-before-`start` example, which its line 111 already states) | every `usageResult` before `inventory.go:439` renders as a `HELPER-USAGE` finding; every failure after it (`:205` NO-DO-WORK-DIR, `:212` PARSE-FAILED, `:215` walk error) is erased by the shim loop `:445-456` | `J` unknown option, `N` not started, `O`/`O2` outside a repository → exit 2 with the finding on stdout. `I` no `do-work/`, `G` unmatched backtick (`REQ-109-bad.md`), `V2` unreadable REQ file (a unix socket named `REQ-110-sock.md`), `F` subdirectory → exit 2, stdout and stderr empty. `X`: a REQ with unterminated frontmatter is skipped, not an error (exit 0) |
| S8 | commit.md:89, inspect.md:121 | "A path from Step 1 that appears in no row is unassociated and moves to Step 4. A quarantined `X` path is absent from this output too, even when a REQ claims it, so match rows only against the M/A/D/XD set that survived the quarantine overlay; a missing row is never permission to read, stage, or commit an `X` path." (inspect.md: "to read an `X` path") | `inventory.go:384-390` (X rows are never candidates), `:420-437` (union with the retained quarantine, quarantined candidates dropped), `:445-453` | `B`: `.env.local` is claimed by REQ-101 and gets no row. `D`: `src/alpha.txt` appended to the quarantine as an earlier inventory's X → no `REQ-101\tsrc/alpha.txt` row although the file is `M` now. `T`: a `secrets.pem` created after `start` is appended to the quarantine by `associate` |

Also measured for prose left as it was: `K` all-X tree → associate exit 1, empty stdout; `P` clean tree → both modes exit 1 with an `INVENTORY-CLEAN` findings block; `B` shows an XD path (`config/app.env`) does get a row, so "M/A/D/XD participate in association" holds.

Draft not used in stage 3: "reads the `## Implementation Summary` file list of every `REQ-*.md` it finds under them" — an `archive/` REQ without a terminal-success status is not read for its list (`inventory.go:283`, measured on REQ-106).

## 5. Guards and gate

Stage 1 (worktree at a1e652f): `audit-lockins.sh` exit 0; `prescribed-shell-canonicalization.sh` exit 0; `quiet-grep-pipeline-audit.sh` exit 0; gate `Maintainer verification passed.`, exit 0, 83s.

Stage 2 (6913dc4): the same three guards exit 0 (`quiet-grep pipeline audit passed (94 tracked shell files, 19 must-flag and 7 must-not-flag shapes)`); gate `Maintainer verification passed.`, exit 0, 86s, 796 Go tests.

Stage 3 (7df6488): `Audit lock-in regressions passed.` exit 0; `Prescribed shell primitive canonicalization checks passed.` exit 0; `quiet-grep pipeline audit passed (94 tracked shell files, 19 must-flag and 7 must-not-flag shapes).` exit 0; gate run as `DO_WORK_GATE_ROOT=<worktree> bash scratchpad/gate.sh` → `Maintainer verification passed.`, exit 0, 82s, one `SKIP` line (the heavy-only probes). Logs: `scratchpad/req597-callers/guard-*.log`, `gate.log`.

Acceptance criteria: the two `inspect.md` commands exit 0 and print association rows (section 1); every replacement is derived from code and a fixture (sections 2-4); the shell-`mv`-nests / `rename(2)`-refuses sentences are consistent (section 3, Verified exact publication); no sentence generalizes across differing commands (the drafts that did are listed); the three guards exit 0; the four out-of-scope callers are untouched.

## 6. Found and not fixed

F1. The exit-2 response is unchanged: both actions still say "skip REQ tracing" for every exit 2, including a printed `HELPER-USAGE` finding. That is the reading that let the broken `inspect.md` blocks survive in every project. Keying the skip on the silent shape and reporting a printed finding as a broken invocation is a one-sentence behaviour change in each action; left out because the request allows no behaviour change beyond the two blocks.

F2. `inventory.go:445-456`: the protected-inventory shim loop replaces the result text unconditionally, so the `NO-DO-WORK-DIR: nothing to associate against` and `PARSE-FAILED: …` lines that `handleAssociate` (`:205`, `:212`) prepares for exactly this compatibility mode never reach a caller, and a walk error's `HELPER-USAGE` finding is erased too. Keeping the text when the outcome is a failure would make the three silent exit-2 causes readable. Code change.

F3. `scripts/protected-inventory.sh:6` puts `"$@"` after the command token, so no caller can pass `--repo-root`; `tools/checks/associate-files.sh:10-17` shows the translating shape. A launcher-level fix would let both actions run from anywhere; until then both say "run from the project root". Code change, one place.

F4. `commit.md:67` tells the agent to re-run the inventory "using the retained quarantine — … do not truncate it, append the new X rows". `start` replaces the quarantine with the current X set (`inventory.go:393`); only `associate` unions (`:420-425`). Which command a re-run should use is a design call the prose predates; left alone.

F5. `commit.md:61` "Exit 2 means this is not a git repo" for `start` also covers a `git status` failure (`inventory.go:368`) and a quarantine write failure (`:394`); both are I/O faults, left alone.

F6. `atomic-download`'s occupancy policy is asymmetric and the guide's Atomic download section does not say so: `--dry-run` refuses any existing target (`commands.go:863-865`, exit 2) while the live run refuses only a directory (`:869`) and `os.Rename` then silently replaces an occupying regular file, reported as `created`, exit 0 (measured, `scratchpad/req597-guide/download`). Line 125 stays true but a reader of it alone does not learn the file case. Also `commands.go:891` `info, _ := os.Stat(targetPath)` followed by `info.Size()` would panic if that stat failed.

F7. `finalization_apply.go:545` runs `diff-tree` without `-m`, so a merge commit among the candidates lists no paths and the `exact` loop stays true; only the preceding `git diff --binary` digest match keeps it unreachable. Worth a REQ, not a guide edit.

F8. `internal/lifecycletiming`: `record-timing-event` after a `start --dry-run` — not timing, but the same family: `inventory.go:412-415` requires a regular quarantine file and `start --dry-run` writes none (`:392-396`), so an `associate` after a dry-run start exits 2 with the not-started finding. Not stated in the actions; noted only because a reader might add `--dry-run` to a block.

F9. The same false script claims live in four callers outside the write set, captured as REQ-601 and untouched here: `present-work.md:136` ("the helper verifies each output against the source separately" — nothing is compared to the source) and `:140` ("the compatibility script"); `ai-report-reference.md:31,37,47` (staging "adjacent to `generated/`" is the system temp directory; "a retained script" does not exist); `architecture-report.md:46` and `install.md:261` ("the compatibility script", not checked); `install.md:50,335` (`SKILL.md.download` is never created; the stray is `SKILL.md.download.<random>`); `board.md:87` ("skips silently" — it prints a `GIT-EXCLUDE-NOT-A-REPOSITORY` warning); `install.md:246` (the installer implements the exclude contract inline, it does not call the helper). `work-reference.md:322` says the wrapper "hands the child the console's own handles" — the twin of the corrected guide sentence (both child streams get the CLI's stderr).

F10. `docs/commit-guide.md` and `docs/inspect-guide.md` carry none of the corrected claims (`grep -rn 'for unassociated\|come back `-`\|archive/\*\*\|re-derives the repository root' skills/` matched only the two actions before the edit; `roadmap.md:56` and `abandon.md:34` use the `**` glob for their own by-hand scans, not as a description of this command).

F11. The lifecycle timing category vocabulary is closed (`recovery-selection, planning, exploration-preflight, builder-work, handback-merge, verification-gate, review, remediation, finalization, cleanup`); `verification` alone is rejected. Not documented where a caller would look.

F12. Line 137 of the guide (report image batch) was corrected beyond the sixteen because it had the same "retained script" defect; if the stage must be limited to the named claims, that is two phrases in one paragraph of `6913dc4`.

F13. No changelog entry was added by any stage; `prime-releases.md` makes each of the three commits a release, and `work-reference.md` Step 9 (finalization) owns the changelog transaction — it should write one entry covering `a1e652f`, `6913dc4` and `7df6488`.
