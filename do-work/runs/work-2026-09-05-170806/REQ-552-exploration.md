# Pre-explore: REQ-552 — replace two coreutils exec sites with pure Go

Read-only exploration. Repo `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`, HEAD `71eb49f3`
("[work run] checkpoint after REQ-588"). Audit commit in the REQ was `dc8a64e3` (2026-09-03).

---

## 1. Re-verify the claim at HEAD

Command run verbatim from the REQ (§Detailed Requirements last bullet):

```
rg -n 'exec\.Command(Context)?\((ctx, )?"(find|cp|mkdir|grep|sed|ls|rm|mv|cat|touch|head|tail|wc)"' \
  skills/do-work/tools/do-work-cli skills/do-work-board/tools/queue-kanban --glob '!*_test.go'
```

Output at HEAD `71eb49f3`, exactly two lines, exit 0:

```
skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go:133:		if output, copyError := exec.Command("cp", draftPath, stagedPath).CombinedOutput(); copyError != nil {
skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:724:	command := exec.Command("find", archiveRoot, "-name", "REQ-*.md", "-print0")
```

**Both sites are still present. Neither has gone away. The claim holds at HEAD.** The RED
condition in §Red-Green Proof is real right now.

Without the `!*_test.go` glob there is a third hit that the lock-in must NOT count:

```
skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go:25:	if output, err := exec.Command("cp", "-R", sourceRoot, stagingRoot).CombinedOutput(); err != nil {
```

So the lock-in **must** keep `--glob '!*_test.go'` or it is red on day one.

Baseline (before any change), both target packages green:

```
ok  github.com/knews2019/skill-do-work/do-work-cli/internal/corehelpers      14.168s
ok  github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands   4.409s
```

---

## 2. The `find` site

### The call, with context

`skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:715-729`:

```go
715	func archiveWalkFailure(repositoryRoot string) string {
716		archiveRoot := filepath.Join(repositoryRoot, "do-work", "archive")
717		info, err := os.Stat(archiveRoot)
718		if os.IsNotExist(err) {
719			return ""
720		}
721		if err != nil || !info.IsDir() {
722			return fmt.Sprint(err)
723		}
724		command := exec.Command("find", archiveRoot, "-name", "REQ-*.md", "-print0")
725		if output, err := command.CombinedOutput(); err != nil {
726			return strings.TrimSpace(string(output) + " " + err.Error())
727		}
728		return ""
729	}
```

The sole caller, `commands.go:659-668`:

```go
659	func runTimestampCommand(repositoryRoot string, scope doctor.TimestampScope, apply, dryRun bool) resultmodel.CommandResult {
660		commandName := CommandRepairTimestamps
661		if scope == doctor.TimestampScopeArchive {
662			commandName = CommandAuditTimestamps
663			if walkFailure := archiveWalkFailure(repositoryRoot); walkFailure != "" {
664				output := "audit-archive-timestamps: the archive walk failed — nothing was inspected.\n"
665				return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, RepositoryRoot: repositoryRoot, ExactTextOutput: &output,
666					Findings: []resultmodel.CommandFinding{helperFinding("TIMESTAMP-ARCHIVE-WALK-FAILED", resultmodel.SeverityError, []string{"do-work/archive"}, walkFailure, resultmodel.FixabilityManual, "the archive was not inspected", nil, nil)}}
667			}
668		}
669		snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
```

### What it is for, and which command path reaches it

- `archiveWalkFailure` has **exactly one caller** (`rg -n 'archiveWalkFailure'` → the definition
  plus `commands.go:663`).
- Command path: `handleAuditTimestamps` (`commands.go:643`) → `runTimestampCommand(..., doctor.TimestampScopeArchive, ...)`.
  The CLI verb is `audit-archive-timestamps` (`CommandAuditTimestamps`). The sibling verb
  `repair-req-timestamps` uses `TimestampScopeActive` and **never** reaches this probe.
- The `find` **stdout is discarded**. Nothing reads the `-print0` list. Only the error branch
  matters: on a non-zero exit the combined stdout+stderr and `err.Error()` become the string.

### Observable behaviour when the archive is unreadable

`audit-archive-timestamps` returns immediately with:

- `Outcome = OutcomeFindings` (non-zero exit for the caller),
- `ExactTextOutput = "audit-archive-timestamps: the archive walk failed — nothing was inspected.\n"`,
- one finding `TIMESTAMP-ARCHIVE-WALK-FAILED`, severity `Error`, affected path `do-work/archive`,
  fixability `Manual`, refusal reason `"the archive was not inspected"`, and **Evidence = the
  non-empty walk-failure string**.

Nothing else runs — `DiscoverRepository` is never called on that path.

**Why this hard gate exists at all** (this is the reason not to just delete it):
`repositorymodel.DiscoverRepository` walks the same tree at
`internal/repositorymodel/repository_model.go:175-179`, and its walk **swallows** entry errors:

```go
175		walkError := filepath.WalkDir(doWorkRoot, func(absolutePath string, directoryEntry fs.DirEntry, entryError error) error {
176			if entryError != nil {
177				snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot inspect %s: %v", absolutePath, entryError))
178				return nil
179			}
```

So without the probe, an unreadable archive would produce a *clean-looking* audit with a warning.
The replacement must keep failing hard, not defer to `DiscoverRepository`.

### Existing shape of the error string (what "class" means here)

Measured, `find` against a directory with mode 000 inside the archive root:

- real `/usr/bin/find` (macOS/BSD): `find: <abs path>: Permission denied`, exit status 1
  → stored string becomes `find: <abs path>: Permission denied exit status 1`
- GNU findutils would produce a different spelling; on this machine even the interactive `find`
  is a `bfs` shim printing `bfs: error: <path>: Permission denied.`

**The exact string is already platform-dependent and pinned by no test.** `rg -n
'TIMESTAMP-ARCHIVE-WALK-FAILED|archive walk failed'` over everything except `do-work/archive/`,
`do-work/runs/` and `kb/` finds only the two production lines in `commands.go` (404, 664, 666).
No fixture, no Go test, no `_dev/tests` file asserts the evidence text.

**Error string class to preserve:** a *non-empty* string that names the offending path and the
reason, produced only when the archive exists, is a directory, and cannot be fully traversed.
Empty string on: archive absent, or archive fully traversable.

Also note a pre-existing wart at line 721-722: when `archiveRoot` exists but is **not** a
directory, `err` is `nil`, so the function returns `fmt.Sprint(nil)` = `"<nil>"`. That is a
non-empty string, so the gate still fires, but the evidence reads `<nil>`. Out of scope to fix
under the REQ's "do not fix nearby code" constraint; worth a discovered-task note, not a change.

### The walk that already exists in `internal/corehelpers/inventory.go`

Exact helper (post-REQ-550, the delegate was inlined here):

```go
// inventory.go:247
func AssociateProjectPaths(repositoryRoot string, candidates []string) (map[string]string, error)
```

Its walk, `inventory.go:253-270`:

```go
253		roots := []string{filepath.Join(repositoryRoot, "do-work", "working"), filepath.Join(repositoryRoot, "do-work", "archive")}
254		for _, root := range roots {
255			walkError := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
256				if walkErr != nil {
257					if os.IsNotExist(walkErr) {
258						return nil
259					}
260					return walkErr
261				}
262				if entry.Type()&os.ModeSymlink != 0 {
263					if entry.IsDir() {
264						return filepath.SkipDir
265					}
266					return nil
267				}
268				if entry.IsDir() || !strings.HasPrefix(entry.Name(), "REQ-") || !strings.HasSuffix(entry.Name(), ".md") {
269					return nil
270				}
271				contents, err := os.ReadFile(path)
```

**Predicate** (the part that matches `find -name 'REQ-*.md'`):
`!entry.IsDir() && strings.HasPrefix(entry.Name(), "REQ-") && strings.HasSuffix(entry.Name(), ".md")`,
plus a symlink skip that `find` (without `-L`) also effectively has.

**Caveat the REQ under-states.** The REQ says "reusing the inventory walk". `AssociateProjectPaths`
is **not** a reusable readability probe: it also walks `do-work/working`, parses every matched file
with `requestmodel.ParseDocument`, reads Implementation Summary sections and builds an ownership
map. There is no extractable "is this tree walkable" helper in `corehelpers` today. So the
realistic replacement is a small local `filepath.WalkDir` in `commands.go` that *matches* the
inventory walk's error discipline, not a call into `AssociateProjectPaths`.

### Concrete replacement

Inside `corehelpers` (the REQ gives latitude on placement), replace lines 724-728 with a walk that
returns the first traversal error as its string, e.g.:

```go
	var walkFailure string
	if err := filepath.WalkDir(archiveRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			walkFailure = fmt.Sprintf("%s: %v", path, walkErr)
			return walkErr
		}
		return nil
	}); err != nil {
		return walkFailure
	}
	return ""
```

Requirements the replacement must meet:

- Non-empty string **only** when the tree cannot be traversed; empty otherwise (absent archive is
  already handled above at 717-720 and must stay).
- The string names the unreadable path and the reason, so the `TIMESTAMP-ARCHIVE-WALK-FAILED`
  evidence stays actionable. Do not collapse to a bare `"walk failed"`.
- Do **not** filter to `REQ-*.md` before deciding failure: `find` reports a `Permission denied` on
  any subdirectory regardless of whether it holds a matching file, and the current gate fires in
  that case. Filtering by the name predicate would make an unreadable REQ-free subdirectory pass
  where it fails today.
- `os.Stat` at 717 stays as-is; only 724-728 changes.
- Removes the `os/exec` dependency from this function but **not** from the file — `commands.go:874`
  still has `exec.Command("curl", ...)`, so the `os/exec` import stays.

After the change the archive is walked once by `filepath.WalkDir` here and once by
`DiscoverRepository`. The REQ's "walked twice per run" wording describes the *current* cost
(subprocess + Go walk); the change removes the subprocess, not the second traversal.

---

## 3. The `cp` site

`skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go:113-140`:

```go
113	func architecturePublish(ctx commandruntime.ExecutionContext, draft, candidate string, dryRun, commit bool) resultmodel.CommandResult {
114		draftPath := draft
115		if !filepath.IsAbs(draftPath) {
116			draftPath = filepath.Join(ctx.RepositoryRoot, draftPath)
117		}
118		data, err := os.ReadFile(draftPath)
119		if err != nil {
120			return usageResult(CommandArchitecture, "draft is not a regular readable file: "+draft)
121		}
122		if info, statErr := os.Stat(draftPath); statErr != nil || !info.Mode().IsRegular() {
123			return usageResult(CommandArchitecture, "draft is not a regular file: "+draft)
124		}
125		if os.Getenv("DO_WORK_COMPATIBILITY_SHIM") == "1" {
126			staged, stageError := os.CreateTemp("", "architecture-report-copy.*")
127			if stageError != nil {
128				return architectureFailure(stageError.Error())
129			}
130			stagedPath := staged.Name()
131			_ = staged.Close()
132			defer os.Remove(stagedPath)
133			if output, copyError := exec.Command("cp", draftPath, stagedPath).CombinedOutput(); copyError != nil {
134				return architectureFailure("draft copy failed: " + strings.TrimSpace(string(output)))
135			}
136			data, err = os.ReadFile(stagedPath)
137			if err != nil {
138				return architectureFailure(err.Error())
139			}
140		}
```

Reached from `handleArchitecture` (`architecture.go:38-40`) on `--publish <draft> <candidate>`.
CLI verb `architecture-report-preflight`.

The whole block is a re-read: `data` is already the draft bytes from line 118. Under the shim the
bytes are round-tripped through a temp file and re-read into the same variable. The only
observable difference between shim and non-shim is the failure path (a broken `cp` turns a
success into `ARCHITECTURE-PREFLIGHT-FAILED`).

### The copy primitive in `last30days.go`

```go
// last30days.go:251
func copyLast30DaysTree(source, target string) error
```

Body at `last30days.go:251-294`. Per regular file:

```go
274			input, err := os.Open(path)
...
278			output, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode().Perm())
...
283			_, copyErr := io.Copy(output, input)
```

- **Preserves mode?** Yes for a fresh destination — it passes `info.Mode().Perm()` of the source
  entry to `OpenFile`, and for directories `os.MkdirAll(destination, info.Mode().Perm())`
  (line 269). Rejects non-regular entries (line 271-273).
- **Creates parent directories?** Only for directories it encounters *during the walk*
  (line 268-270). It does not create the target root itself; the caller `publishLast30Days`
  does that (`last30days.go:152-160`, `os.MkdirAll(parent, 0o755)` then `os.MkdirTemp`).

### Two traps for the builder

**T1 — `copyLast30DaysTree` cannot be called directly here.** Two reasons:

1. `filepath.WalkDir` on a *regular file* yields exactly one callback whose `relative` is `"."`,
   and the function returns `nil` at line 260-262. Passing a file as `source` copies **nothing**
   and returns `nil` — a silent no-op.
2. `os.O_EXCL` at line 278 would fail anyway, because `os.CreateTemp` at line 126 has **already
   created** `stagedPath`.

So the REQ's "`last30days.go` has `copyLast30DaysTree` with the copy primitive" means the
*primitive* (`os.Open` + `os.OpenFile` + `io.Copy`), not the function. Reuse the shape, not the
call.

**T2 — mode. Measured, not assumed.** `cp src dst` where `dst` already exists does **not** change
`dst`'s mode:

```
src mode: 644
dst mode after cp: 600      # dst pre-created at 0600, unchanged by cp
dst bytes: srcbytes
```

`os.CreateTemp` creates at `0600`. So today "the same file lands at `stagedPath`" means: draft
**bytes**, temp-file **mode 0600** — the draft's mode is never copied.

Therefore:

- `os.WriteFile(stagedPath, data, 0o600)` — correct. `WriteFile` on an existing file uses
  `O_TRUNC` and the perm argument is ignored for an existing file, so the mode stays `0600` either
  way.
- An `os.OpenFile(stagedPath, os.O_WRONLY|os.O_TRUNC, 0o600)` + `io.Copy` from `os.Open(draftPath)`
  — also correct, and closer to the `last30days` primitive.
- **Wrong:** copying the draft's mode onto `stagedPath` (e.g. `0644`), or using
  `O_CREATE|O_EXCL` (fails on the already-created temp).

Simplest correct edit, replacing 133-135:

```go
		if copyError := os.WriteFile(stagedPath, data, 0o600); copyError != nil {
			return architectureFailure("draft copy failed: " + copyError.Error())
		}
```

Keep the `"draft copy failed: "` prefix — it is the evidence string of
`ARCHITECTURE-PREFLIGHT-FAILED` and is what the fixture in §4 keys off if it is retained.

`strings` and `os/exec` stay imported in `architecture.go` regardless (`strings` used at 49, 108,
162; `exec` at 45 for `git rev-parse`).

---

## 4. Existing tests that cover both sites

### `audit-archive-timestamps`

Go tests — **none exercise the walk-failure path**:

- `internal/corehelpers/commands_test.go:16` — registration list, asserts `CommandAuditTimestamps`
  appears among the dispatch names. Would **not** catch a behaviour change.
- `internal/corehelpers/commands_test.go:263` — `{CommandAuditTimestamps, []string{"--fix", "--dry-run"}}`,
  a recovery-argv table entry. Would **not** catch a behaviour change.
- `internal/doctor/doctor_commands_test.go:127` — asserts the *legacy* mechanic filenames
  `repair-req-timestamps.sh` / `audit-archive-timestamps.sh` are gone. Unrelated. Would **not** catch.

Shell fixture — `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh`, 5+ cases covering
`--fix` repair, report-only, clean archive, impossible ordering. All of them run against a
**readable** archive, so none reaches `archiveWalkFailure`'s error branch. They would catch a
regression that made the probe fire on a healthy archive (a false positive), which is the useful
half. They would **not** catch a loss of the failure detection.

**Net: no test covers the unreadable-archive path.** The lock-in in §5 does not cover it either
(it counts exec sites). If the builder wants a real safety net for the `find` half, the honest
option is one Go test in `internal/corehelpers` that chmods a subdirectory to 000 and asserts
`TIMESTAMP-ARCHIVE-WALK-FAILED` with non-empty evidence — but the REQ's Constraints say "no test
files beyond the lock-in", so this needs the maintainer's call, not a unilateral addition.

### Architecture staging under the compatibility shim

Go tests — `internal/toolboxcommands/architecture_test.go` has 5 tests
(`TestArchitectureNameParserOrdersNumericSuffix`,
`TestRemediationArchitectureAbsoluteScanAndReadDirFailure`,
`TestRemediationArchitectureFailedClaimRemainsOccupied`,
`TestArchitecturePublishUsesFirstFreeNumericBundle`,
`TestArchitecturePublishStopsOnNonCollisionClaimFailure`,
`TestArchitecturePublishCommitCommitsPublishedBundle`).
**None of them sets `DO_WORK_COMPATIBILITY_SHIM`**, so none enters the `cp` block at all. They
would not catch a change to it. `internal/toolboxcommands/commands_test.go:7` is registration only.

`rg -n 'DO_WORK_COMPATIBILITY_SHIM'` over Go tests finds it only in
`internal/corehelpers/checks_test.go:203,245` and `internal/corehelpers/inventory_test.go:269` —
different commands.

### >>> The one test that WILL break <<<

`_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh:195-215`, the case
`# architecture-report-preflight: a failed copy cannot expose partial HTML as a prior baseline.`

```bash
202	cat > "$partial_repo/bin/cp" <<'SH'
203	#!/usr/bin/env bash
204	printf '<!doctype html><title>Partial' > "$2"
205	exit 1
206	SH
207	chmod +x "$partial_repo/bin/cp"
208	(cd "$partial_repo" && PATH="$partial_repo/bin:$PATH" "$preflight_script" --publish drafts/report.html "$publish_candidate" >/dev/null 2>&1)
209	[ "$?" -ne 0 ] || fail_case 'architecture-report-preflight reported a failed copy as published'
210	[ ! -e "$partial_repo/$publish_candidate/index.html" ] \
211	  || fail_case 'architecture-report-preflight exposed partial HTML after a failed copy'
```

The fixture wrapper at line 17 of that file sets `DO_WORK_COMPATIBILITY_SHIM=1`, so the `cp`
branch runs, and `exec.Command("cp", ...)` resolves through the process `PATH` — which the test
has prefixed with a fake `cp` that writes partial bytes and exits 1.

**Once the `cp` is pure Go, this PATH shim is inert.** The copy succeeds, the publication
succeeds, `index.html` is created — and **both** assertions at 209 and 210 flip to FAIL, plus the
follow-up scan at 212-215.

This directly collides with the REQ's Constraints: *"Tests unchanged; the existing package tests
are the safety net"* and *"no test files beyond the lock-in"*. **The builder cannot satisfy both.**
Surface it rather than silently editing or silently leaving it red. The three honest resolutions:

- **O1** Delete the case. It pins a failure mode (`cp` failing) that stops existing. The case-count
  is computed at run time from the files (`_dev/tests/prescribed-shell-case-count.sh`,
  `prescribed-shell-scripts-behavior.sh:45`), so no pinned number needs updating.
- **O2** Rewrite the case to force a Go-side copy failure, e.g. `chmod 500 "$TMPDIR"` so
  `os.CreateTemp` fails, or point `TMPDIR` at a read-only directory — this exercises the
  `stageError` branch at 127-129 rather than the copy branch.
- **O3** Keep the case and accept it red. Not acceptable.

**Recommendation: O1.** The block being removed is a re-read of bytes the function already holds;
after the change there is no realistic partial-write window left to pin, and O2 would test
`os.CreateTemp` rather than the copy.

**Tier note:** `_dev/tests/prescribed-shell-harness.sh:11-14` and
`prescribed-shell-scripts-behavior.sh:8-11` both refuse unless `DO_WORK_MAINTAINER_TIER=heavy`.
So this breakage will **not** show up in the fast gate. It appears only in the heavy lane. Do not
read a green fast gate as proof the `cp` change is clean.

---

## 5. The lock-in

### `_dev/tests/audit-lockins.sh` — structure

156 lines, executable, `#!/usr/bin/env bash`, `set -uo pipefail` (line 3 — no `-e`, so a
non-matching `rg` inside a command substitution does not abort). `repo_root` at line 5,
`failure_count=0` at line 6, and at the end:

```bash
151	if [ "$failure_count" -gt 0 ]; then
152	  exit 1
153	fi
154
155	printf 'Audit lock-in regressions passed.\n'
```

Five blocks, each with the identical shape — a comment naming the audit finding and its REQ, a
`$( ... )` capture into a descriptive variable, then an `if [ -n ... ]` that loops the capture with
a here-string and emits `FAIL:` to stderr while bumping `failure_count`:

1. Finding 10 — exported one-line delegates with no production caller (REQ-550), lines 8-41
2. Finding 5 — caller-less toolbox shell shims (REQ-551), lines 43-58
3. Shipped shell delegating check (REQ-551 companion), lines 60-72
4. Finding 8 — dead path pointers in records (REQ-549), lines 74-131
5. Finding 2 — hand-rolled CLI launcher preamble copies (REQ-553), lines 133-149

### Two assertions, verbatim

```bash
if [ -n "$callerless_shims" ]; then
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    printf 'FAIL: caller-less toolbox shell shim found: %s\n' "$f" >&2
    failure_count=$((failure_count + 1))
  done <<< "$callerless_shims"
fi
```

```bash
if [ -n "$hand_rolled_preambles" ]; then
  while IFS= read -r f; do
    [ -z "$f" ] && continue
    printf 'FAIL: hand-rolled do-work-cli launcher preamble outside the preamble pair: %s\n' "$f" >&2
    failure_count=$((failure_count + 1))
  done <<< "$hand_rolled_preambles"
fi
```

### Counting tools used

Mixed. `rg` at lines 28 (`rg -c`), 48 (`rg -l --fixed-strings`), 63 (`rg -q`), 138 (`rg -l --glob`);
`find` at 11 and 62; `awk` at 11 and 108; `grep`/`sed` for filtering. It does **not** have a single
counting convention — `rg` is the idiom for "search a shipped tree for a pattern", which is the
right fit here.

`*_test.go` exclusion: yes, the file already does it twice — `find ... ! -name '*_test.go'`
(line 11) and `--glob '!*_test.go'` (line 28). The new block should use `--glob '!*_test.go'`,
and **must**, per §1 (`update_transaction_test.go:25`).

### Registration in `_dev/tests/contracts/probe-lanes.sh`

Lines 29-30, already present, unchanged by this REQ:

```bash
register_probe audit_lockins_probe "$repo_root/_dev/tests/audit-lockins.sh" \
  'audit lock-in regressions failed (see the attributed FAIL lines above).'
```

`register_probe` (probe-lanes.sh:5-15) refuses if the script is not executable. `probe-lanes.sh` is
sourced by `_dev/tests/contract-regressions.sh:72-73`. Confirmed: **do not create the file and do
not touch the registration**, exactly as the REQ's Constraints say.

### Paste-ready assertion

Insert after the Finding 2 block (line 149) and before the closing `if [ "$failure_count" -gt 0 ]`
at line 151, so the blocks stay in the file's existing descending-finding order… actually the file
orders 10, 5, companion, 8, 2, so append at the end of the blocks:

```bash
# Finding 9: exec-where-pure-go-exists (REQ-552)
# A coreutils subprocess spawned for work the same Go module already does in the standard
# library. Test files are excluded: fixture setup may shell out, shipped code may not.
coreutils_exec_sites="$(
  rg -n 'exec\.Command(Context)?\((ctx, )?"(find|cp|mkdir|grep|sed|ls|rm|mv|cat|touch|head|tail|wc)"' \
    "$repo_root/skills/do-work/tools/do-work-cli" \
    "$repo_root/skills/do-work-board/tools/queue-kanban" \
    --glob '!*_test.go' 2>/dev/null
)"
if [ -n "$coreutils_exec_sites" ]; then
  while IFS= read -r coreutils_site; do
    [ -z "$coreutils_site" ] && continue
    printf 'FAIL: coreutils spawned where the module already has pure Go: %s\n' "${coreutils_site#"$repo_root/"}" >&2
    failure_count=$((failure_count + 1))
  done <<< "$coreutils_exec_sites"
fi
```

**Verified by running it** as a standalone script against HEAD `71eb49f3`:

```
FAIL: coreutils spawned where the module already has pure Go: skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go:133:		if output, copyError := exec.Command("cp", draftPath, stagedPath).CombinedOutput(); copyError != nil {
FAIL: coreutils spawned where the module already has pure Go: skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:724:	command := exec.Command("find", archiveRoot, "-name", "REQ-*.md", "-print0")
failure_count=2
```

RED now, exactly 2, and it goes to 0 when both sites are replaced. Matches §Constraints' "coreutils
exec sites in the two Go modules: 0 after this REQ (today 2)".

Note the pattern is a **closed list of command names**, which CLAUDE.md § "State conditions, not
lists" and the prime's `[family: closed-enumeration-for-a-condition]` trap both warn about. It is
the REQ's own Reproduce pattern, so keeping it byte-identical is the right call here — the lock-in
must fail on exactly what the REQ's RED proof prints. Say so in the comment rather than widening
it; widening would drag in `open`/`rundll32`/`xdg-open` (`skills/do-work-board/tools/queue-kanban/openbrowser.go:20-26`),
`ps`, `tar`, `sh`, `python3`, `curl` and the argv-driven runners, none of which this REQ addresses.

---

## 6. REQ-550 dependency

`do-work/archive/REQ-550-collapse-four-exported-one-line-delegates-into-their-targets.md` —
status `completed`, `completed_at: 2026-09-04T13:45:00Z`, commit `667896dcb39062a9c889ea7533a9665a73ce01e7`.

1. In `corehelpers` it touched **only** `internal/corehelpers/inventory.go`: inlined the private
   `associatePaths` into the exported `AssociateProjectPaths` (that is why the `filepath.WalkDir`
   now sits at `inventory.go:255` inside the exported function).
2. It did **not** touch `internal/corehelpers/commands.go` — the file holding REQ-552's `find` site.
3. Its other Go edits were in `archivefetch`, `doctor` (×2) and `gateevidence`. None overlaps
   REQ-552's write set.
4. It **created** `_dev/tests/audit-lockins.sh` and registered it in
   `_dev/tests/contracts/probe-lanes.sh` — which is exactly why REQ-552 says "the file already
   exists … do not create it and do not change its registration".
5. **The write set has moved: the block is clear.** `_dev/tests/audit-lockins.sh` is now shared
   between the two REQs (550 wrote it, 552 appends to it), but 550 is done and committed, so there
   is no live overlap. REQ-552 can proceed.

---

## 7. Prime and lessons

### Should `prime-do-work-cli.md` be added to `prime_files`?

**Yes.** `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` exists (25 KB, 112 lines) and
covers both packages explicitly:

- line 19 — `internal/corehelpers/` owns "the remaining utility handlers and leaf check, inventory,
  Git, publication, reservation, and survey mechanics"
- line 22 — `internal/toolboxcommands/` owns "notes, architecture preflight/publication, …"

`prime_files: []` is currently empty. Adding it costs the builder one file read and hands over the
Verify commands it will need anyway.

### Traps that bear on replacing an exec with in-process Go

From `prime-do-work-cli.md`:

- **line 63, `[family: opaque-evidence-projection]`** — *"a generic fallback or opaque aggregate
  can discard the exact blocker a caller must act on; derive output, specific typed records, and
  recovery argv from one observation set, and* ***keep implementation dependencies in the Go
  standard library***.*" This is the trap that argues *for* REQ-552, and it also sets the bar for
  the `find` replacement: the evidence string must still name the exact unreadable path, not
  collapse to a generic message.

- **line 81, `[family: smoke-vs-characterization]`** — *"Registration and smoke checks can stay
  green while replacement semantics diverge; compare status, ordered evidence, actions, paths, and
  effects at authority boundaries before retiring an implementation."* Directly applicable: the
  only Go tests touching either site are registration/argv tables (§4), so a green
  `go test ./internal/corehelpers ./internal/toolboxcommands` proves very little here.

- **line 101, Verify** — *"Toolbox migration: `go test -race ./internal/toolboxcommands` … and
  `GOOS=windows GOARCH=amd64 go test -c ./internal/toolboxcommands -o <temporary-path>`."* Worth
  running: removing a `cp` subprocess is a portability improvement, and the Windows cross-compile
  is cheap.

From `lessons-do-work-cli.md`:

- **line 62, `[family: fixture-cost-is-subprocess-spawning]` (REQ-574)** — *"`t.Setenv` cannot be
  traded for a per-child `Cmd.Env` when the code under test shells out itself, because* ***the
  in-process call resolves its own binary from the process PATH***.*" This is the generalisation of
  the §4 breakage: a fixture that controls behaviour by putting a fake binary on `PATH` only works
  while the code actually shells out. Turning the `cp` into in-process Go silently disarms every
  such fixture. **Search for PATH-shim fixtures before deleting an exec** — here that is
  `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh:202-215`.

- **line 23, `[family: final-boundary-identity]` (REQ-416)** and **prime line 60** — the final
  mutation boundary rules. Not triggered: the `stagedPath` write is to a private temp file that is
  read and removed, not a publication boundary. The real publication boundary is further down
  `architecturePublish` (the `rootedMkdirExclusive` claim) and is untouched.

- `rg` over `lessons-do-work-cli.md` for `exec.Command`, `coreutils`, `in-process`, `subprocess`,
  backticked `cp`/`find`, `CombinedOutput` returns nothing else relevant. There is no prior lesson
  about replacing a coreutils call with Go.

---

## Summary of what the builder needs to know

| # | Item |
|---|------|
| F1 | Both exec sites still exist at HEAD `71eb49f3`; RED condition holds; both packages green at baseline. |
| F2 | `archiveWalkFailure` (`commands.go:715`) has one caller, `runTimestampCommand:663`; only the error string is used; it exists because `DiscoverRepository` downgrades walk errors to warnings. |
| F3 | No test anywhere pins the walk-failure evidence text, and the current text is already platform-dependent (`find:` vs `bfs:`). Only the class — non-empty, names the path — must survive. |
| F4 | `AssociateProjectPaths` is not a reusable readability probe; write a small local `filepath.WalkDir` matching its predicate and error discipline. Do not filter by `REQ-*.md` before deciding failure. |
| F5 | `cp` does **not** change `stagedPath`'s mode (measured: stays `0600` from `os.CreateTemp`). Use `os.WriteFile(stagedPath, data, 0o600)` or `OpenFile(...O_WRONLY\|O_TRUNC...)` + `io.Copy`. Never `O_EXCL`. |
| F6 | `copyLast30DaysTree` cannot be called for a single file — `WalkDir` on a regular file returns at `relative == "."` and copies nothing, and its `O_EXCL` would fail on the pre-created temp. Reuse the primitive, not the function. |
| **R1** | **`_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh:195-215` puts a fake failing `cp` on `PATH` and will go RED once the copy is in-process. Heavy tier only, so the fast gate will not show it. Recommend deleting the case (case counts are computed at run time, nothing else to update) and recording the conflict with the REQ's "tests unchanged" constraint in the REQ file.** |
| F7 | Lock-in: append one block to `_dev/tests/audit-lockins.sh` before line 151, in the file's `$(…)` + `if [ -n … ]` + `printf 'FAIL: …'` style, using `rg` with `--glob '!*_test.go'`. Verified 2 → expected 0. Do not touch `probe-lanes.sh:29-30`. |
| F8 | REQ-550 is completed at `667896dc` and touched `corehelpers/inventory.go` only, never `commands.go`. Write set is clear. |
| F9 | Add `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` to `prime_files`. Key traps: `opaque-evidence-projection` (line 63), `smoke-vs-characterization` (line 81), and the PATH-shim lesson at `lessons-do-work-cli.md:62`. |
