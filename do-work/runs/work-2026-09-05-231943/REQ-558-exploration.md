# REQ-558 Exploration — Keep one nil-root guard in git_transaction.go and delete the other eight

**Route: C (plan-first).** The REQ's central factual claim is inverted. Eight of the nine guards protect a real nil-pointer dereference; only one is genuinely redundant. A builder who follows the Detailed Requirements as written reintroduces panics on the transaction rollback boundary.

**Prior exploration:** none. `do-work/runs/work-2026-09-05-170806/` holds explorations for REQ-486, 552, 554, 556, 583, 587 and 591 only. No REQ-558 artifact exists anywhere under `do-work/runs/`. Everything below is measured fresh at HEAD `fce57fc`.

---

## 1. What still holds, and what has gone stale

The REQ's *measurements* survive. Its *provenance* and its *reasoning* do not.

| REQ claim | Verdict | Evidence at HEAD |
|---|---|---|
| Nine `root [=!]= nil` sites in the file | Holds | 9 lines: `git_transaction.go:297, 998, 1009, 1026, 1114, 1143, 1176, 1193, 1276` |
| No test covers any nil branch | Holds | `rg -l 'OpenRoot\|rollback root is unavailable\|rooted filesystem handle is unavailable' internal/gittransaction/*_test.go` matches nothing |
| Nil is producible at exactly one point | Holds | 8 `os.OpenRoot` sites; 7 return on error; only `:994` falls through |
| That point is in `rollbackTransaction` | **Stale/wrong** | No such function exists. It is `rollbackFailure`, `git_transaction.go:989` |
| The other eight are redundant via `rootedOpenSnapshot` | **Wrong** | Four of the five named callers dereference `root` directly and never reach it |
| Introduced by `b877eb69`, `0a5d4e44`, `a43b2587`, `01d920dd`; audited `dc8a64e3`; report `83594c5e` | **Stale** | All six: `fatal: Not a valid object name`. Blame shows one squashed commit `7bd3464` |
| File has 17 commits behind it | **Stale** | `git log --oneline -- <file> \| wc -l` = **1** |
| Expected net delta −25 | **Wrong** | Only ~2 lines are pure deletion; the rest requires restructuring, giving ≈ −10 to −15 |
| Depends on REQ-557 to keep `internal/` clear | **Wrong reason, and unmet** | REQ-557 is `status: pending`; its write_set never touches `gittransaction`. Real overlap is `_dev/tests/audit-lockins.sh` |

The Reproduce command runs verbatim and reproduces exactly:

```
$ rg -n 'root [=!]= nil' skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go | wc -l
9
$ rg -l 'OpenRoot|rollback root is unavailable|rooted filesystem handle is unavailable' \
    skills/do-work/tools/do-work-cli/internal/gittransaction/*_test.go || echo 'NO TEST covers any nil-root branch'
NO TEST covers any nil-root branch
```

Baseline is green: `go build ./...` exit 0, `go test ./internal/gittransaction/` → `ok ... 1.780s`.

---

## 2. The fact everything turns on

Every method on a nil `*os.Root` panics. Measured on the toolchain in this environment (`go version go1.26.1 linux/amd64`; module pins `go 1.25.0`):

```
Close      PANIC: runtime error: invalid memory address or nil pointer dereference
Lstat      PANIC: ...
Mkdir      PANIC: ...
Remove     PANIC: ...
Rename     PANIC: ...
Open       PANIC: ...
OpenFile   PANIC: ...
MkdirAll   PANIC: ...
```

`os.Root` wraps an unexported `*root` field. There is no nil-safe method and no usable zero value. So a nil-root guard in this file is not style — it is the difference between a degraded rollback and a crash halfway through restoring a repository.

---

## 3. Where nil comes from — the REQ's one correct structural insight

Eight `os.OpenRoot` calls exist in the file. Seven return immediately on error, so `root` is non-nil by construction below them:

| Site | Enclosing function | Error handling |
|---|---|---|
| `:138` | `RecordCreatedDirectory` | `:139-141` return |
| `:221` | `RecordTouched` | `:222-224` return |
| `:324` | `RecordCreated` | `:325-327` return |
| `:338` | `captureCreatedDirectoryIdentities` | `:339-341` return |
| `:370` | `RecordPublished` | `:371-373` return |
| `:467` | transaction preflight | `:468-470` `failTransaction` |
| `:917` | `changedTargets` | `:918-920` return |
| **`:994`** | **`rollbackFailure`** | **`:995-997` appends to `rollback.Errors` and falls through** |

```go
// git_transaction.go:994
root, rootError := os.OpenRoot(repositoryRoot)
if rootError != nil {
    rollback.Errors = append(rollback.Errors, "open rollback root: "+rootError.Error())
}
if root != nil {
    defer root.Close()
}
```

`:994` is the sole producer of a nil `*os.Root` that keeps executing. That part of the REQ is right. Only the function's name is wrong.

---

## 4. The nine guards, each judged by what happens if you delete it

| # | Line | Site | Verdict if deleted | First dereference reached |
|---|---|---|---|---|
| G1 | `:297` | `inspectCreatedObject` | **Redundant — safe** | none; `rootedOpenSnapshot:1276` catches it |
| G2 | `:998` | `defer root.Close()` | **PANIC** | `Close` on nil receiver |
| G3 | `:1009` | `root != nil && privateStateStillOriginal(...)` (dirty tracked) | **PANIC** | `:1250` `root.Lstat` |
| G4 | `:1026` | `root != nil && privateStateStillOriginal(...)` (private untracked) | **PANIC** | `:1250` `root.Lstat` |
| G5 | `:1114` | `if root != nil { root.Lstat }` in created-paths loop | **PANIC** | `:1115` `root.Lstat` |
| G6 | `:1143` | `if !recorded \|\| root == nil` in created-directories loop | **PANIC** | `:1147` `root.Lstat` |
| G7 | `:1176` | `rollbackDirtyTracked` | **PANIC** | `:1180` `root.Lstat`, else `:1213` `root.Mkdir` |
| G8 | `:1193` | `trackedPublicationStillOwned` | **PANIC** | `:1197` `root.Lstat` |
| G9 | `:1276` | `rootedOpenSnapshot` | **PANIC** | `:1280` `root.Lstat` |

**Eight load-bearing, one redundant.** The REQ proposes keeping one and deleting eight; the evidence says keep eight and delete one.

### Why G1 really is redundant

With `:297-299` removed, `inspectCreatedObject` calls `rootedCreatedTargetSnapshot` (`:1265`) → `rootedOpenSnapshot` (`:1276`), which returns `errors.New("rooted filesystem handle is unavailable")`. `isMissingPathError` (`:1351`) tests `os.IsNotExist` / `os.ErrNotExist` / `syscall.ENOENT` and does not match a bare `errors.New`, so the switch at `:301-308` falls to `default` and returns `createdObjectReplaced` — identical to what the guard returned at `:298`. Zero behaviour delta. This is the single guard the REQ could delete as-is.

---

## 5. Why the REQ's redundancy proof fails

The REQ states:

> the two `privateStateStillOriginal` callers, `trackedPublicationStillOwned`, and `rollbackDirtyTracked` all reach `rootedOpenSnapshot`, which tests the same value again.

Each of those three functions is **two-branch**, and the proof only walks one branch. On the other branch each calls `root.Lstat` directly:

```go
// privateStateStillOriginal — git_transaction.go:1248
func privateStateStillOriginal(root *os.Root, state targetState) bool {
	if !state.existed {
		_, err := root.Lstat(filepath.FromSlash(state.path))   // :1250 — direct, no guard
		return os.IsNotExist(err)
	}
	info, digest, err := rootedRegularSnapshot(root, state.path) // :1253 — reaches :1276
	...
}
```

```go
// trackedPublicationStillOwned — git_transaction.go:1192
	if !published.existed {
		_, err := root.Lstat(filepath.FromSlash(path))          // :1197 — direct, no guard
		return os.IsNotExist(err)
	}
	info, digest, err := rootedRegularSnapshot(root, path)       // :1200 — reaches :1276
```

```go
// rollbackDirtyTracked — git_transaction.go:1175
	if !published.existed {
		if _, statError := root.Lstat(filepath.FromSlash(state.path)); ... // :1180 — direct
```

`rootedOpenSnapshot` is the safety net only when the target **existed**. When it did not, control never gets there. The audit read the happy branch of each function and generalised.

---

## 6. A live panic at HEAD the REQ does not know about

The private-untracked branch of `rollbackFailure`:

```go
// git_transaction.go:1023
if state.privateUntracked {
    published, recorded := recorder.publishedPrivate[state.path]
    if !recorded {
        if root != nil && privateStateStillOriginal(root, state) { continue }   // :1026 — G4
        rollback.Errors = append(rollback.Errors, "private target was not identity-recorded: "+state.path)
        continue
    }
    action, privateRollbackError := quarantineAndRollbackPrivate(root, state, published)  // :1032 — NO GUARD
```

`quarantineAndRollbackPrivate` (`:1204`) has no nil test anywhere. Its first root use is `root.Mkdir(directory, 0o700)` at `:1213`.

**Reachable when:** a private untracked target was identity-recorded via `RecordPublished` (`:370`), the operation then failed, and `os.OpenRoot` at `:994` fails at rollback time (repository directory moved, removed, permissions changed, or a transient `EMFILE`). Result: nil dereference panic in the middle of a rollback.

The guard set in this file is not over-defensive. It is **incomplete and inconsistent** — `rollbackDirtyTracked` is guarded, its sibling `quarantineAndRollbackPrivate` is not.

Three root-taking functions have no nil test: `quarantineAndRollbackPrivate` (`:1204`, reachable with nil), `privateStateStillOriginal` (`:1248`, reachable only through G3/G4), `rootedCreateRegular` (`:1310`, masked because `:1213` would panic first). The REQ counts two and reads the gap as evidence the guards are noise; the gap is evidence they are unfinished.

---

## 7. The design space, and the constraint that cannot be met

The REQ's Builder Guidance grants exactly the latitude needed: *"Firm on one producible-nil point; latitude on whether downstream functions take a non-nil root or an explicit no-root branch."*

**O1 — Early exit at the produce point (recommended).** Turn `:995-1000` into an early return: record `"open rollback root: ..."`, run the tail (index check `:1164-1170`, result assembly `:1163-1173`, extracted as `finalizeRollbackResult`), return. Then `defer root.Close()` unconditionally. Root is non-nil by construction below, so all nine guards go, **and the `:1032` panic stops existing** — no targeted bug fix required, because the nil path is gone.

**O2 — Nil-safe wrapper type.** A `rollbackRootHandle` struct holding `openRoot *os.Root`, with methods returning a sentinel error when absent. One nil test in one accessor. Preserves more behaviour, but adds a type and ~30 lines against CLAUDE.md's *delete before you add*, and net delta goes to roughly zero.

**O3 — Split into rooted and root-free rollback functions.** Duplicates the git-only cleanup across two functions. Rejected.

### The constraint problem

The REQ's Constraints say *"Behaviour preserved on every rollback path"*. No option meets that literally:

- **O1** skips, on `OpenRoot` failure, the `git restore --staged` unstaging (`:1004`, `:1081`, `:1102`), the `existingUntrackedAllowed` plain-`os` restore (`:1038-1058`), and several error strings. In practice those are doomed too — `runGit` and `os.WriteFile` both target the same unreachable directory — but not provably in every case (a transient `EMFILE` would let git succeed).
- **O2** changes the error string at `:1144` from `"created directory was not identity-recorded"` to `"created directory changed after publication; preserved replacement"`.

And the REQ's stated safety net cannot arbitrate: *"proved by the existing tests, never by a new mock"* — but **no test reaches a nil root at all**. The suite is green at HEAD and would stay green through a change that introduces a panic there. The REQ forbids the one instrument that would make the change safe.

This is why the route is **C**. The Detailed Requirements and the Constraints both need maintainer amendment before code is written.

---

## 8. Recommended plan

1. **Amend the REQ** to record that the redundancy premise is inverted (8 load-bearing, 1 redundant) and that `rollbackTransaction` is `rollbackFailure`.
2. **Take O1.** Early exit at `:994`; extract the tail as `finalizeRollbackResult`; delete all nine guards.
3. **Renegotiate the behaviour constraint** to "the rollback Outcome, FailureKind and Status are preserved; the error list on the unopenable-root path may shrink." O1 preserves `OutcomeRisk` / `FailureRollback` / `RollbackIncomplete` exactly.
4. **Renegotiate the test constraint** to allow one focused test on the nil-root path. Without it the only path the change touches ships unverified.
5. **Add the lock-in** to `_dev/tests/audit-lockins.sh` (section 9).
6. **Sequence after REQ-557**, but for the real reason: both append to `_dev/tests/audit-lockins.sh`.

---

## 9. Lock-in assertion (paste-ready)

Append after the REQ-553 hand-rolled-preamble block, before the `if [ "$failure_count" -gt 0 ]` tail. Matches the file's house style: `$repo_root`-relative path variable, count, `printf 'FAIL: ...' >&2`, `failure_count=$((failure_count + 1))`.

```bash
# Finding 3: nil-root-guards-git-transaction (REQ-558)
# A nil *os.Root is producible at exactly one site: rollbackFailure's os.OpenRoot
# (git_transaction.go:994), the only one of the eight OpenRoot calls in the file that does
# not return on error. Every other root-taking function receives a value that is non-nil by
# construction, so a second nil test there guards a state the control flow cannot reach.
# Pinned at "at most one": an implementation that keeps `if root == nil` after the open
# scores 1, and one that tests the os.OpenRoot error instead scores 0. Both are correct.
# grep -c exits non-zero on no-match and prints 0 -- its documented interface, so the
# collapsing fallback is safe here. The [ -f ] test is what separates "zero guards" from
# "file moved", which the count alone could not distinguish.
nil_root_guard_file="$repo_root/skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go"
if [ -f "$nil_root_guard_file" ]; then
  nil_root_guard_count="$(grep -c 'root [=!]= nil' "$nil_root_guard_file" || true)"
  if [ "$nil_root_guard_count" -gt 1 ]; then
    printf 'FAIL: %s holds %s nil-root guards; at most 1 is allowed (nil is producible only where os.OpenRoot fails without an early return)\n' \
      "${nil_root_guard_file#"$repo_root/"}" "$nil_root_guard_count" >&2
    failure_count=$((failure_count + 1))
  fi
else
  printf 'FAIL: nil-root guard lock-in cannot run; %s is missing\n' \
    "${nil_root_guard_file#"$repo_root/"}" >&2
  failure_count=$((failure_count + 1))
fi
```

**Why `-gt 1` and not `-eq 1`.** The REQ says "pin at 1". If the builder resolves nil by testing `rootError` at `:995` rather than writing `if root == nil`, the `root [=!]= nil` count is **0**, and `-eq 1` would fail a correct implementation. "At most one" is green at both 0 and 1, red at 2+.

**Trap check against `_dev/primes/prime-shell-commands.md` § Unchecked Exit Status Reads as Content.** The prime names `grep -c` returning non-zero on no-match as a case where "the emptiness **is** the information" and the collapsing fallback is correct. The residual ambiguity — file missing vs. zero guards — is closed by the `[ -f ]` branch, which fails loudly instead of scoring 0.

**Dry-run verified at HEAD:**

```
FAIL: skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go holds 9 nil-root guards; at most 1 is allowed
failure_count=1
```

RED today, green at 0 or 1. Registration at `_dev/tests/contracts/probe-lanes.sh:29` is untouched, per the REQ.

**Naming check.** `nil_root_guard_file`, `nil_root_guard_count` — multi-word, findable by plain-text search, per `skills/do-work/crew-members/coding-guardrails.md` § 5 Naming for Reach. Any Go identifier the builder introduces (`finalizeRollbackResult`) must clear the same bar.

---

## 10. Red-green

**RED (verified at HEAD):** the Reproduce command prints 9 sites and `NO TEST covers any nil-root branch`; `_dev/tests/audit-lockins.sh` with the new block exits 1.

**Second RED the plan should add:** drive a transaction with one `privateUntracked` target through `RecordPublished`, fail the mutation, make `repositoryRoot` unopenable before `rollbackFailure` runs. Today this panics at `:1213`.

**GREEN when all four hold:**

1. `rg -n 'root [=!]= nil' <file>` prints at most one line.
2. `go test ./internal/gittransaction/` green, unchanged.
3. The nil-root test returns `RollbackIncomplete` instead of panicking.
4. `bash _dev/tests/audit-lockins.sh` exits 0.

---

## 11. Files to change

| Path | Change |
|---|---|
| `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` | Early exit at `:994`; extract `finalizeRollbackResult`; delete guards at `:297-299`, `:998-1000`, the conjuncts at `:1009` and `:1026`, `:1114-1119`, the disjunct at `:1143`, `:1176-1178`, `:1193-1195`, `:1276-1279` |
| `/home/user/skill-do-work/_dev/tests/audit-lockins.sh` | Append the section 9 block. Do not create the file; do not touch `probe-lanes.sh:29` |
| `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` | **Conditional on maintainer sign-off** (REQ says "no test files beyond the lock-in"). One test pinning that a failed rollback-root open returns `RollbackIncomplete` rather than panicking |

---

*Generated by Explore agent*