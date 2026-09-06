# Exploration: REQ-552 — replace two coreutils exec sites with the pure Go the package already has

Read-only exploration. Repo `/home/user/skill-do-work`, HEAD `15e2ec33cb34df4194db262e2aae460c01862858`
("[work run] REQ-592 pipeline: triage, estimate, exploration, scope"). Audit commit named in the REQ
is `dc8a64e3` (2026-09-03). Prior exploration at `do-work/runs/work-2026-09-05-170806/REQ-552-exploration.md`
was written against HEAD `71eb49f3` on a different machine.

Environment note that changes one of the prior exploration's observations: this container runs Linux
with GNU findutils 4.9.0 and GNU coreutils, as root. The prior exploration measured macOS/BSD `find`
and a `bfs` shim. Every measurement below was re-taken here.

---

## 0. What is stale in the prior exploration

| # | Prior statement | Status at HEAD |
|---|-----------------|----------------|
| S1 | "Both exec sites still present, `commands.go:724` and `architecture.go:133`" | **Holds.** Same two lines, same numbers. |
| S2 | "Without the `!*_test.go` glob there is a third hit … the lock-in **must** keep it" | **Holds.** `update_transaction_test.go:25` is still there. |
| S3 | "`_dev/tests/audit-lockins.sh` is 156 lines, five blocks, ends at 151-155" | **Holds.** Byte-for-byte the same structure. |
| S4 | "REQ-550 completed at `667896dc`, corehelpers write set clear" | **Holds.** `status: completed`, `completed_at: 2026-09-04T13:45:00Z`. |
| S5 | "No test covers the unreadable-archive path; the shell fixtures all run against a readable archive" | **WRONG.** There is a second PATH-shim fixture the prior exploration never opened. See §5.2. It was already in the tree when that exploration ran (added in `7bd3464`, 2026-09-05 16:24). |
| S6 | "The one test that WILL break" (singular, the `cp` fixture) | **Understated.** Two case files break, not one. Five `fail_case` lines, not three. |
| S7 | Error-string measurements (`find:` vs `bfs:` spellings) | Machine-specific and no longer reproducible here. The conclusion — the string is platform-dependent and pinned by no test — still holds. |

The REQ's own §What is also stale on one number: it says "Of 90 `exec.Command` sites in do-work-cli".
At HEAD there are **103** non-test sites (197 including tests). The `85 run git` half is exact
(`rg -c 'exec\.Command(Context)?\((ctx, )?"git"' … --glob '!*_test.go'` sums to 85). The two coreutils
sites — the only figure the work acts on — are exact.

---

## 1. The claim, re-verified at HEAD

Command run verbatim from the REQ §Detailed Requirements last bullet:

```
rg -n 'exec\.Command(Context)?\((ctx, )?"(find|cp|mkdir|grep|sed|ls|rm|mv|cat|touch|head|tail|wc)"' \
  skills/do-work/tools/do-work-cli skills/do-work-board/tools/queue-kanban --glob '!*_test.go'
```

Output, exactly two lines, exit 0:

```
skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:724:	command := exec.Command("find", archiveRoot, "-name", "REQ-*.md", "-print0")
skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go:133:		if output, copyError := exec.Command("cp", draftPath, stagedPath).CombinedOutput(); copyError != nil {
```

**Both sites exist. The RED condition in §Red-Green Proof is real right now.**

Dropping the glob adds a third hit that the lock-in must not count:

```
skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go:25:	if output, err := exec.Command("cp", "-R", sourceRoot, stagingRoot).CombinedOutput(); err != nil {
```

That is a fixture spawning `cp -R` to stage a tree, not shipped code, and it is outside this REQ's
scope. **Without `--glob '!*_test.go'` the lock-in would sit at 1 forever after both production sites
are fixed — red on day one and never green.** Keep the glob.

Green simulation: the same `rg` with the two target files also excluded returns **exit 1, no matches**.
So the pattern reaches 0 once both sites are replaced, and nothing else in the two modules matches.
Other non-git exec sites in those modules are `sh` (x2), `xdg-open`, `tar`, `rundll32`, `python3`,
`ps`, `open`, `curl` — none is in the pattern's closed command list.

Baseline, both target packages, run **in tree**:

```
ok  github.com/knews2019/skill-do-work/do-work-cli/internal/corehelpers      3.844s
ok  github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands  2.012s
```

---

## 2. The `find` site

### The call

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

The sole caller is `commands.go:663`, inside `runTimestampCommand` (`commands.go:659`):

```go
663			if walkFailure := archiveWalkFailure(repositoryRoot); walkFailure != "" {
664				output := "audit-archive-timestamps: the archive walk failed — nothing was inspected.\n"
665				return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, RepositoryRoot: repositoryRoot, ExactTextOutput: &output,
666					Findings: []resultmodel.CommandFinding{helperFinding("TIMESTAMP-ARCHIVE-WALK-FAILED", …)}}
667			}
```

`rg -n 'archiveWalkFailure' skills/` returns exactly the definition (715) and that one call (663).

Command path: `handleAuditTimestamps` (`commands.go:642`) → `runTimestampCommand(…, doctor.TimestampScopeArchive, …)`.
CLI verb `audit-archive-timestamps`. The sibling verb `repair-req-timestamps` uses `TimestampScopeActive`
(`commands.go:640`) and never reaches the probe.

The `find` **stdout is discarded** — nothing reads the `-print0` list. Only the error branch matters.

### Why the gate exists (the reason not to just delete it)

`repositorymodel.DiscoverRepository` walks the same tree and **swallows** entry errors —
`internal/repositorymodel/repository_model.go:175-179`:

```go
175	walkError := filepath.WalkDir(doWorkRoot, func(absolutePath string, directoryEntry fs.DirEntry, entryError error) error {
176		if entryError != nil {
177			snapshot.WarningMessages = append(snapshot.WarningMessages, fmt.Sprintf("cannot inspect %s: %v", absolutePath, entryError))
178			return nil
179		}
```

Without the probe an unreadable archive would produce a clean-looking audit with a warning. The
replacement must keep failing hard, not defer to `DiscoverRepository`.

### The error string class to preserve

`rg -n 'TIMESTAMP-ARCHIVE-WALK-FAILED|archive walk failed'` excluding `do-work/archive/`,
`do-work/runs/` and `kb/` finds only three production lines — `commands.go:404` (the severity table),
`commands.go:664` and `commands.go:666`. **No Go test, no fixture and no `_dev/tests` file asserts the
evidence text.** Only the class must survive:

- non-empty **only** when the tree exists, is a directory, and cannot be fully traversed;
- empty when the archive is absent or fully traversable;
- the string names the offending path and the reason, so the finding stays actionable
  (`prime-do-work-cli.md:63`, `[family: opaque-evidence-projection]` — do not collapse it to a
  generic "walk failed").

### The walk the REQ points at, and why it is not directly reusable

`internal/corehelpers/inventory.go:247` `func AssociateProjectPaths(repositoryRoot string, candidates []string) (map[string]string, error)`,
walk at `inventory.go:253-270`:

```go
253		roots := []string{filepath.Join(repositoryRoot, "do-work", "working"), filepath.Join(repositoryRoot, "do-work", "archive")}
255		walkError := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
268			if entry.IsDir() || !strings.HasPrefix(entry.Name(), "REQ-") || !strings.HasSuffix(entry.Name(), ".md") {
```

Same predicate as `find -name 'REQ-*.md'`, plus a symlink skip (262-267) that `find` without `-L`
effectively has too. But `AssociateProjectPaths` **also** walks `do-work/working`, parses every match
with `requestmodel.ParseDocument`, reads Implementation Summary sections and builds an ownership map.
There is no extractable "is this tree walkable" helper in `corehelpers` today. So the REQ's phrase
"reusing the inventory walk" means *matching its error discipline*, not calling it.

### Concrete replacement (compiled and behaviour-checked)

Replace `commands.go:724-728` with:

```go
	var archiveWalkEvidence string
	if walkError := filepath.WalkDir(archiveRoot, func(entryPath string, entry os.DirEntry, entryError error) error {
		if entryError != nil {
			archiveWalkEvidence = fmt.Sprintf("%s: %v", entryPath, entryError)
			return entryError
		}
		return nil
	}); walkError != nil {
		return archiveWalkEvidence
	}
	return ""
```

Requirements it meets, and traps it avoids:

- **Do not filter to `REQ-*.md` before deciding failure.** `find` reports `Permission denied` on any
  unreadable subdirectory whether or not it holds a matching file, and the current gate fires in that
  case. Filtering by the name predicate would let an unreadable REQ-free subdirectory pass where it
  fails today.
- The `os.Stat` guard at 717-722 stays untouched.
- **No import changes.** `fmt`, `os` and `path/filepath` are already imported, and `os/exec` stays for
  the `curl` call at `commands.go:874`.
- `archiveWalkEvidence` is two words and plain-text findable, per
  `skills/do-work/crew-members/coding-guardrails.md` § 5 Naming for Reach (lines 115-132).

One accepted semantic difference: `find` reports every unreadable subdirectory and still exits
non-zero; this walk stops at the first. Both produce a non-empty evidence string, so the gate behaves
identically — the evidence just names one path instead of several. That is cheaper and still specific.

**Pre-existing wart, out of scope.** At 721-722, when `archiveRoot` exists but is not a directory,
`err` is `nil`, so the function returns `fmt.Sprint(nil)` = `"<nil>"`. Non-empty, so the gate fires,
but the evidence reads `<nil>`. The REQ says "do not fix nearby code" — leave it, note it as a
discovered task.

---

## 3. The `cp` site

`skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go:113-140`:

```go
118		data, err := os.ReadFile(draftPath)
125		if os.Getenv("DO_WORK_COMPATIBILITY_SHIM") == "1" {
126			staged, stageError := os.CreateTemp("", "architecture-report-copy.*")
130			stagedPath := staged.Name()
131			_ = staged.Close()
132			defer os.Remove(stagedPath)
133			if output, copyError := exec.Command("cp", draftPath, stagedPath).CombinedOutput(); copyError != nil {
134				return architectureFailure("draft copy failed: " + strings.TrimSpace(string(output)))
135			}
136			data, err = os.ReadFile(stagedPath)
```

Reached from `handleArchitecture` on `--publish <draft> <candidate>`; CLI verb
`architecture-report-preflight`. The whole block is a re-read: `data` already holds the draft bytes
from line 118, and under the shim they are round-tripped through a temp file into the same variable.

### The copy primitive in `last30days.go`

`last30days.go:251` `func copyLast30DaysTree(source, target string) error`, body 251-294. Per regular
file: `os.Open` (274) → `os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode().Perm())`
(278) → `io.Copy` (283).

**It cannot be called here**, for two independent reasons:

1. `filepath.WalkDir` on a *regular file* yields one callback whose `relative` is `"."`, and the
   function returns `nil` at 260-262. Passing a file as `source` copies nothing and returns `nil` — a
   silent no-op.
2. `os.O_EXCL` at 278 would fail anyway, because `os.CreateTemp` at `architecture.go:126` has already
   created `stagedPath`.

So "the `io.Copy` primitive the package already has" means the *shape*, not the function.

### Mode, measured on this machine

```
src mode: 644
dst mode after cp: 600      # dst pre-created at 0600, unchanged by cp
dst bytes: srcbytes
```

`cp src dst` where `dst` exists does not change `dst`'s mode. `os.CreateTemp` creates at `0600`. So
"the same file lands at `stagedPath` with the same mode" means: draft **bytes**, temp-file mode
**0600**. The draft's mode is never copied.

### Concrete replacement (compiled and behaviour-checked)

Replace `architecture.go:133-135` with:

```go
		if copyError := os.WriteFile(stagedPath, data, 0o600); copyError != nil {
			return architectureFailure("draft copy failed: " + copyError.Error())
		}
```

- `os.WriteFile` on an existing file uses `O_TRUNC` and ignores the perm argument, so the mode stays
  `0600` either way. Correct.
- Keep the `"draft copy failed: "` prefix — it is the `ARCHITECTURE-PREFLIGHT-FAILED` evidence string.
- **Wrong alternatives:** copying the draft's mode onto `stagedPath`; any use of `O_CREATE|O_EXCL`
  (the temp already exists); calling `copyLast30DaysTree`.
- **No import changes.** `os` is imported; `os/exec` stays for the git calls at 45 and 98; `strings`
  stays for 49 and 103.

---

## 4. Existing Go test coverage of both sites

### `audit-archive-timestamps`

No Go test exercises the walk-failure path:

- `internal/corehelpers/commands_test.go:16` — registration list only.
- `internal/corehelpers/commands_test.go:263` — a recovery-argv table entry `{CommandAuditTimestamps, ["--fix","--dry-run"]}`.
- `internal/doctor/doctor_commands_test.go:127` — asserts the legacy mechanic filenames are gone.

None would catch a behaviour change.

### Architecture staging under the compatibility shim

`internal/toolboxcommands/architecture_test.go` has six tests (lines 15, 21, 49, 70, 90, 116). **None
sets `DO_WORK_COMPATIBILITY_SHIM`**, so none enters the `cp` block at all.
`rg -n 'DO_WORK_COMPATIBILITY_SHIM' --glob '*_test.go'` finds it only in
`internal/corehelpers/checks_test.go:203,245` and `internal/corehelpers/inventory_test.go:269` —
different commands.

This is `prime-do-work-cli.md:81`, `[family: smoke-vs-characterization]` exactly: registration and
smoke checks stay green while replacement semantics diverge. **A green
`go test ./internal/corehelpers ./internal/toolboxcommands` proves almost nothing about this change.**

Empirically confirmed: with both sites patched, `go test ./internal/toolboxcommands` is `ok` and the
`corehelpers` failure count is identical to the unpatched baseline. Zero Go-test delta.

> Testing note the builder needs: **these packages only pass in tree.** A copy of the module built
> elsewhere fails 59 `corehelpers` subtests, all with
> `inventory_test.go:136: retained inventory script: stat …/checks/uncommitted-inventory.sh: no such file or directory`
> — the test resolves a sibling `checks/` directory relative to the module path. Do not read those as
> real failures.

---

## 5. The two shell fixtures that break — both proven

Method: copied the module into the scratchpad, built an unpatched binary and a binary with both sites
replaced as in §2 and §3, and ran each fixture's scenario against both. No repo file was modified.

### 5.1 `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh:195-215`

Case header at 195: `# architecture-report-preflight: a failed copy cannot expose partial HTML as a prior baseline.`

```bash
202	cat > "$partial_repo/bin/cp" <<'SH'
203	#!/usr/bin/env bash
204	printf '<!doctype html><title>Partial' > "$2"
205	exit 1
206	SH
207	chmod +x "$partial_repo/bin/cp"
208	(cd "$partial_repo" && PATH="$partial_repo/bin:$PATH" "$preflight_script" --publish drafts/report.html "$publish_candidate" >/dev/null 2>&1)
209	[ "$?" -ne 0 ] || fail_case 'architecture-report-preflight reported a failed copy as published'
210	[ ! -e "$partial_repo/$publish_candidate/index.html" ] || fail_case 'architecture-report-preflight exposed partial HTML after a failed copy'
```

The fixture wrapper at line 17 of that file sets `DO_WORK_COMPATIBILITY_SHIM=1`, so the `cp` branch
runs and `exec.Command("cp", …)` resolves through the process `PATH` the test prefixed.

Measured:

```
base    exit=1   index.html absent
patched exit=0   index.html EXISTS
```

Assertions at **209**, **210** and the follow-up scan at **214** all flip to FAIL.

### 5.2 `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh:220-236` — missed by the prior exploration

Case header at 220: `# audit-archive-timestamps: a file walk that could not complete is a failure, never a clean answer with a count of zero.` (REQ-268 instance 2)

```bash
229	audit_walk_bin="$fixture_root/audit-walk-bin"
231	printf '%s\n' '#!/usr/bin/env bash' 'exit 3' > "$audit_walk_bin/find"
232	chmod +x "$audit_walk_bin/find"
233	audit_walk_output="$(PATH="$audit_walk_bin:$PATH" "$core_scripts/audit-archive-timestamps.sh" "$audit_walk_project" 2>&1)" \
234	  && fail_case 'audit-archive-timestamps failed-walk case exited zero after the walk failed'
235	printf '%s' "$audit_walk_output" | grep -q 'audit clean' \
236	  && fail_case 'audit-archive-timestamps failed-walk case reported clean for an archive it never scanned'
```

`skills/do-work/scripts/audit-archive-timestamps.sh:7` execs `do-work-cli.sh --repo-root "$1" --format text audit-archive-timestamps`,
so the fake `find` is what `archiveWalkFailure` resolves. Neither `do-work-cli.sh` nor either preamble
calls `find` itself (grepped, no hits), so line 724 is the only consumer of the shim.

Measured against a fixture repo with one clean archived REQ:

```
base    exit=1   "audit-archive-timestamps: the archive walk failed — nothing was inspected."
patched exit=0   "do-work: archive audit clean (1 file(s) scanned)."
```

Assertions at **234** and **236** both flip to FAIL. (`audit clean` is emitted at `commands.go:793`.)

This case was added in `7bd3464` (2026-09-05 16:24, "REQ-585"), an ancestor of both the prior
exploration's baseline and HEAD. The prior exploration's §4 statement that all the shell fixtures run
against a readable archive is simply wrong.

### 5.3 Consequences

**Tier.** `_dev/tests/prescribed-shell-harness.sh:11-14` refuses unless `DO_WORK_MAINTAINER_TIER=heavy`.
Neither breakage appears in the fast gate. **A green fast gate is not proof this change is clean.**

**The REQ's Constraints cannot be satisfied as written.** §Constraints says "Tests unchanged; the
existing package tests are the safety net" and "no test files beyond the lock-in". Five `fail_case`
lines across two files say otherwise. Options:

- **O1 — delete both case blocks** (`architecture-report-preflight.sh:195-215`,
  `audit-archive-timestamps.sh:220-236`). Nothing else needs updating: case counts are computed at run
  time by `_dev/tests/prescribed-shell-case-count.sh` and accumulated in
  `prescribed-shell-scripts-behavior.sh:19-27,45-46`; no figure is stored anywhere.
- **O2 — rewrite each case to force a Go-side failure.** For the `cp` case: make `os.CreateTemp` fail
  (read-only `TMPDIR`), which tests the `stageError` branch at 127-129, not the copy. For the `find`
  case: `chmod 0000` a subdirectory — **not portable, it is a no-op for root**, so it would pass
  vacuously in any root container.
- **O3 — keep both cases and accept them red.** Not acceptable.

**Recommendation: O1 for the `cp` case, O1 for the `find` case.** The `cp` block being removed is a
re-read of bytes the function already holds, so there is no partial-write window left to pin. The
`find` case's intent (a failed walk is never a clean answer) survives in the code as the hard gate at
`commands.go:663-667`, but after O1 nothing tests it — record that honestly rather than smuggling in a
Go test the Constraints forbid.

**Write set must grow by two files.** Neither is claimed by any other pending REQ (checked every
`write_set` in `do-work/queue/` and `do-work/working/`).

Root cause worth carrying forward — `lessons-do-work-cli.md:62`,
`[family: fixture-cost-is-subprocess-spawning]` (REQ-574): a fixture that controls behaviour by putting
a fake binary on `PATH` only works while the code actually shells out. **Search for PATH-shim fixtures
before deleting any exec.** The search that finds them:

```
rg -n "bin/(cp|find|git|mkdir)|PATH=\"" _dev/tests/prescribed-shell-cases/
```

---

## 6. The lock-in

### `_dev/tests/audit-lockins.sh` at HEAD

156 lines, `-rwxr-xr-x`, `#!/usr/bin/env bash`, `set -uo pipefail` at line 3 (no `-e`, so a
non-matching `rg` inside a command substitution does not abort). `repo_root` at 5, `failure_count=0`
at 6, and the close:

```bash
151	if [ "$failure_count" -gt 0 ]; then
152	  exit 1
153	fi
154
155	printf 'Audit lock-in regressions passed.\n'
```

Five blocks, each the same shape — a comment naming the finding and its REQ, a `$(…)` capture into a
descriptive variable, then `if [ -n … ]` looping the capture with a here-string, emitting `FAIL:` to
stderr and bumping `failure_count`:

1. Finding 10 — exported one-line delegates (REQ-550), lines 8-41
2. Finding 5 — caller-less toolbox shims (REQ-551), lines 43-58
3. Shipped shell delegating check (REQ-551 companion), lines 60-72
4. Finding 8 — dead path pointers (REQ-549), lines 74-131
5. Finding 2 — hand-rolled launcher preambles (REQ-553), lines 133-149

The file already strips `repo_root` from reported paths the same way (line 102,
`"${citing_file#"$repo_root/"}"`), and already excludes `*_test.go` twice — `find … ! -name '*_test.go'`
(line 11) and `--glob '!*_test.go'` (line 28).

### Registration — do not touch

`_dev/tests/contracts/probe-lanes.sh:29-30`:

```bash
register_probe audit_lockins_probe "$repo_root/_dev/tests/audit-lockins.sh" \
  'audit lock-in regressions failed (see the attributed FAIL lines above).'
```

`register_probe` (probe-lanes.sh:5-15) refuses if the script is not executable. `probe-lanes.sh` is
sourced by `_dev/tests/contract-regressions.sh:72-73`, whose tier defaults to `fast`
(`contract-regressions.sh:6`). So this lock-in runs in the fast gate, as the REQ says. **Do not create
the file and do not change the registration.**

### Paste-ready assertion

Insert **after line 149** (the closing `fi` of the Finding 2 block) and **before line 151**, keeping a
blank line either side to match the file's spacing:

```bash
# Finding 9: exec-where-pure-go-exists (REQ-552)
# A coreutils subprocess spawned for work the same Go module already does in the standard
# library. Test files are excluded on purpose: fixture setup may shell out, shipped code may
# not — suiteinstall/update_transaction_test.go:25 spawns `cp -R` and is not this finding.
# The command list is the audit's own Reproduce pattern, kept byte-identical so this fails
# on exactly what the REQ's RED proof prints; widening it would drag in open/xdg-open/
# rundll32/ps/tar/sh/python3/curl, which this finding does not address.
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

**Verified by running it standalone against HEAD:**

```
FAIL: coreutils spawned where the module already has pure Go: skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:724:	command := exec.Command("find", archiveRoot, "-name", "REQ-*.md", "-print0")
FAIL: coreutils spawned where the module already has pure Go: skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go:133:		if output, copyError := exec.Command("cp", draftPath, stagedPath).CombinedOutput(); copyError != nil {
failure_count=2
```

RED now with exactly 2; 0 once both sites are replaced (proven by the exclusion run in §1).

`coreutils_exec_sites` and `coreutils_site` are both two-word, plain-text-findable names, per
coding-guardrails § 5.

### Pinned value: 0, not 2

`do-work/audits/audit-2026-09-03.md:426` says "currently 2; red at 3" (`<= 2`). The REQ's Constraints
line 47 says "0 after this REQ (today 2)". **The REQ is stricter and is the binding one** — write it so
any match fails, which is what the block above does.

### The closed enumeration

The pattern is a closed list of command names, which `CLAUDE.md` § "State conditions, not lists" and
`_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale both warn against. Keeping it
byte-identical to the REQ's Reproduce is right here: the lock-in must fail on exactly what the RED
proof prints. The comment says so, so the next reader does not silently widen it.

---

## 7. Dependencies and write-set contention

`do-work/archive/REQ-550-collapse-four-exported-one-line-delegates-into-their-targets.md` —
`status: completed` (line 4), `completed_at: 2026-09-04T13:45:00Z` (line 27). Its `corehelpers` edit
was `inventory.go` only; `commands.go` was never touched. **Dependency satisfied, write set clear.**

`_dev/tests/audit-lockins.sh` is in the `write_set` of six pending REQs — 552, 554, 555, 556, 557, 558
— each appending one block. Run them serially or expect append conflicts. `REQ-557` depends on this
REQ and touches `corehelpers/checks.go`, not `commands.go`: no overlap.

Neither prescribed-shell case file appears in any pending `write_set`.

---

## 8. Prime files and lessons

`prime_files: []` is empty in the REQ. **Add `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`.**
It exists (25 KB) and covers both packages by name — line 19 (`internal/corehelpers/` owns the leaf
check, inventory, Git, publication, reservation and survey mechanics) and line 22
(`internal/toolboxcommands/` owns notes and architecture preflight/publication).

Traps that bear directly on this work:

- `prime-do-work-cli.md:63`, `[family: opaque-evidence-projection]` — "keep implementation dependencies
  in the Go standard library". This is the argument *for* the REQ, and it sets the bar for the `find`
  replacement: the evidence must still name the exact unreadable path.
- `prime-do-work-cli.md:81`, `[family: smoke-vs-characterization]` — registration and smoke checks stay
  green while replacement semantics diverge. See §4: the only Go tests touching either site are
  registration and argv tables.
- `prime-do-work-cli.md:101` (Verify) — `go test -race ./internal/toolboxcommands` and
  `GOOS=windows GOARCH=amd64 go test -c ./internal/toolboxcommands -o <temp>`. Worth running: dropping
  a `cp` subprocess is a portability win and the cross-compile is cheap.
- `lessons-do-work-cli.md:62`, `[family: fixture-cost-is-subprocess-spawning]` (REQ-574) — the
  generalisation of §5. `rg` over `lessons-do-work-cli.md` for `exec.Command`, `coreutils`,
  `in-process`, `subprocess`, `CombinedOutput` finds nothing else relevant; there is no prior lesson
  about replacing a coreutils call with Go. **This work should produce one.**

---

## 9. Route recommendation

**Route B — explore-then-build.**

The two code edits are mechanical and settled: exact anchors, exact replacements, both compiled and
behaviour-checked, no import changes, no new names beyond one local. That looks like Route A.

What makes it Route B is that the REQ's Constraints are provably unsatisfiable. "Tests unchanged; no
test files beyond the lock-in" collides with five `fail_case` lines across two heavy-tier PATH-shim
fixtures, measured, not guessed. Discovering that collision, proving it, and deciding what happens to
each case is exploration work, and it changes the write set. It is not Route C: there is no
architecture to design, the fix is two small diffs, and `Builder Guidance` already fixes the approach
("Firm: pure Go on both sides").

Per `CLAUDE.md` § Communication, the builder should not stall on this. Write the challenge and the
chosen resolution into the REQ file, take O1 on both cases, and continue.

---

## 10. Verification recipe for the builder

```bash
# RED (must print two lines before the change)
rg -n 'exec\.Command(Context)?\((ctx, )?"(find|cp|mkdir|grep|sed|ls|rm|mv|cat|touch|head|tail|wc)"' \
  skills/do-work/tools/do-work-cli skills/do-work-board/tools/queue-kanban --glob '!*_test.go'

# GREEN: same command prints nothing, then
cd skills/do-work/tools/do-work-cli && go test ./internal/corehelpers ./internal/toolboxcommands
gofmt -l internal/corehelpers internal/toolboxcommands   # must print nothing
go vet ./internal/corehelpers ./internal/toolboxcommands
GOOS=windows GOARCH=amd64 go test -c ./internal/toolboxcommands -o /dev/null

# Fast gate (includes the new lock-in)
bash _dev/tests/audit-lockins.sh          # "Audit lock-in regressions passed."
bash _dev/tests/contract-regressions.sh

# Heavy gate — REQUIRED here, the two fixture edits show up nowhere else. Needs user permission.
DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh
DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/prescribed-shell-cases/architecture-report-preflight.sh
```

Run the Go tests **in tree**. See the note in §4.

---

## Summary table

| # | Item |
|---|------|
| F1 | Both exec sites exist at HEAD `15e2ec33`, unchanged at `commands.go:724` and `architecture.go:133`. RED condition real. Both packages green at baseline in tree. |
| F2 | The lock-in must keep `--glob '!*_test.go'` — `suiteinstall/update_transaction_test.go:25` spawns `cp -R` and is out of scope; without the glob the lock-in never reaches 0. |
| F3 | `archiveWalkFailure` (`commands.go:715`) has one caller (`:663`); only the error string is used; it exists because `DiscoverRepository` (`repository_model.go:175-179`) downgrades walk errors to warnings. |
| F4 | No test pins the walk-failure evidence text. Preserve the class only: non-empty and path-naming when the tree cannot be traversed, empty otherwise. |
| F5 | `AssociateProjectPaths` is not a reusable readability probe. Write a local `filepath.WalkDir` in `commands.go` matching its error discipline. Do not filter by `REQ-*.md` before deciding failure. |
| F6 | Measured: `cp` does not change an existing destination's mode. `stagedPath` is `0600` from `os.CreateTemp`. Use `os.WriteFile(stagedPath, data, 0o600)`. Never `O_EXCL`. |
| F7 | `copyLast30DaysTree` cannot be called for a single file — `WalkDir` returns at `relative == "."` copying nothing, and its `O_EXCL` fails on the pre-created temp. Reuse the primitive, not the function. |
| **R1** | **`_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh:220-236` shims `find` on PATH. Proven RED after the change: base exit=1 "walk failed" → patched exit=0 "audit clean". Missed by the prior exploration.** |
| **R2** | **`_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh:195-215` shims `cp` on PATH. Proven RED: base exit=1 no index.html → patched exit=0 index.html created.** |
| R3 | Both are heavy-tier only (`prescribed-shell-harness.sh:11-14`). A green fast gate proves nothing about them. |
| F8 | Deleting either case costs nothing else — case counts are computed at run time (`prescribed-shell-case-count.sh`), no stored figure. |
| F9 | Lock-in: one block appended after `audit-lockins.sh:149`, verified 2 → 0, pinned at 0 (stricter than the audit's own `<= 2`). Registration at `probe-lanes.sh:29-30` untouched. |
| F10 | REQ-550 completed 2026-09-04, touched `inventory.go` only. Write set clear. `audit-lockins.sh` is shared by six pending REQs — serialize. |
| F11 | Add `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` to `prime_files`. Key traps at lines 63, 81, 101 and `lessons-do-work-cli.md:62`. |
| F12 | REQ §What's "90 exec.Command sites" is stale — 103 at HEAD. "85 run git" is exact. Neither changes the work. |

---

*Generated by Explore agent*