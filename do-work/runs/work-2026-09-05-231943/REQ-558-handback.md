# REQ-558 hand-back — nil-root guards in git_transaction.go

Branch `worktree-agent-REQ-558-nil-root-guards`, head `d46f5d4f3437f776dee84f635d63d5d1c51c55ad`.
Files changed: `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go`,
`_dev/tests/audit-lockins.sh`. Nothing else touched. This file is written, not staged and not
committed.

## D1 — The pinned number is 8, not the 1 the request named

The request assumed one producer of a nil `*os.Root` means one needed test. Three independent
traces plus my own reading say the opposite: one nil value fans out to eleven consumers inside
`rollbackFailure` alone, across four independent loops, and each consumer has to decide for
itself what an unusable handle means for the target it is holding. There is no chokepoint that
dominates them, so a count of producers is not a count of guards.

Eight of the nine sites are reached with a nil handle in unmodified end-to-end runs of the
exported `ExecuteTransaction`; seven of those eight turn a reported incomplete rollback into a
nil-pointer panic when deleted individually. One site — `rootedOpenSnapshot` — was reachable by
nothing. It was deleted. The lock-in pins 8.

## Guards deleted (1)

**`rootedOpenSnapshot`, was line 1276** — `if root == nil { return nil, empty, nil,
errors.New("rooted filesystem handle is unavailable") }`

- **Why nil cannot reach it.** `rootedOpenSnapshot` is reached only through three wrappers
  (`rootedRegularSnapshot`, `rootedCreatedTargetSnapshot`, `rootedRegularPreimage`) from thirteen
  call sites, all in this one file, none in a test. I checked each site's handle back to its
  `os.OpenRoot`. Nine of them (lines 227, 252, 375, 496, 514, 683, 712, 935, 947) carry a handle
  from an `os.OpenRoot` at line 221, 324, 370, 467 or 917, and every one of those five sites
  returns immediately when the open fails — so the handle is non-nil at every use. The remaining
  four sit behind a guard that is being kept and that returns before the call: line 300 behind the
  guard in `inspectCreatedObject`, line 1200 behind the guard in `trackedPublicationStillOwned`,
  line 1253 behind the `root != nil &&` short circuits in `rollbackFailure` (both of them), and
  line 1222 behind the guard in `rollbackDirtyTracked`.
- **What the call shape now guarantees instead.** The invariant is the same one the request asked
  for, and it was already expressed in the code: a caller either holds a handle whose open failure
  already returned, or it takes its own explicit no-root branch first. Those branches are the
  eight surviving guards, and each returns the conservative answer its own caller needs —
  `createdObjectReplaced`, `false`, `"rollback root is unavailable"`, or skipping the check. The
  deleted site added no information to any of them; it re-answered a question the call site had
  already answered. Go cannot express a non-nil pointer in a type, so what replaces the deleted
  test is a doc comment on `rootedOpenSnapshot` that states the precondition, names today's
  no-root branches as the current set rather than a closed list, and names the one caller the
  deleted test never covered anyway (below, F1).
- **Behaviour delta: zero.** On the only path where the deleted test could ever have fired
  (`inspectCreatedObject` with the guard at 297 also removed, which is not being done), the error
  it returned is not an `os.ErrNotExist`, so `isMissingPathError` was false and the switch fell to
  `default` and returned `createdObjectReplaced` — the same value the surviving guard at 297
  returns directly.

## Guards kept (8)

Every one of these is marked `reachable-keep` by the traces, and I re-read each call site myself.
Line numbers are at my head.

| Line | Function | Why it stays |
|---|---|---|
| 297 | `inspectCreatedObject` | Reached with nil from `rollbackFailure:1122`. It is the only thing keeping control off the unguarded `root.Remove` at line 1131 in the same loop. |
| 998 | `rollbackFailure` | Four lines below the producer. `(*os.Root).Close` panics on a nil receiver, and it would panic while returning, destroying the rollback report that was just assembled. |
| 1009 | `rollbackFailure` | The `root != nil &&` short circuit. Without it, `privateStateStillOriginal` calls `root.Lstat` directly for a tracked deletion preimage — a documented, tested input shape. |
| 1026 | `rollbackFailure` | Same shape, private untracked branch. An absent private target is the normal case there, and `TestPrivateRollbackIgnoresDeclaredTargetsNotYetMutated` drives exactly this path. |
| 1114 | `rollbackFailure` | Direct `root.Lstat` on the possibly-nil handle, in the loop that removes this invocation's own created files. No downstream guard exists on this path. |
| 1143 | `rollbackFailure` | The `|| root == nil` half is the only thing between a nil handle and `root.Lstat` at 1147 and `root.Remove` at 1152 when the directory identity was recorded. |
| 1176 | `rollbackDirtyTracked` | Its one caller passes the possibly-nil handle unchecked. Both branches below it dereference: `root.Lstat` at 1180, and `root.Mkdir` at 1213 via `quarantineAndRollbackPrivate`. |
| 1193 | `trackedPublicationStillOwned` | Two callers pass the possibly-nil handle. The `!published.existed` branch does `root.Lstat` directly, and the `false` it returns is the conservative answer both callers turn into a preserved-replacement error. |

No guard came back `cannot-establish`, so nothing is kept on a maybe.

## The lock-in

One block in `_dev/tests/audit-lockins.sh`, written in that file's existing per-Finding shape (the
Finding 4 / REQ-557 block is the model). It scans the one file with the audit's own reproduce
pattern `root [=!]= nil`, reads `rg`'s exit status rather than a piped total (status > 1 means the
scan could not run and fails loudly; status 1 means every guard is gone and fails), and compares
the count for equality — floor and ceiling. A ceiling alone catches a tenth guard accreting; a
floor alone catches a load-bearing guard removed as redundant, which is the move this REQ set out
to make nine times over and found it could make once. The file was not created and its
registration in `_dev/tests/contracts/probe-lanes.sh` was not changed.

## Evidence (verbatim)

Build, vet, format, in `skills/do-work/tools/do-work-cli`:

```
=== go build ===
build exit=0
=== go vet ===
vet exit=0
=== gofmt -l . ===
gofmt exit=0 (no paths listed = clean)
=== go test ./internal/gittransaction/ ===
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction	2.353s
```

The same suite before the change: `ok  github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction  2.205s`. Green and unchanged.

Lock-in RED direction 1 — the deleted guard grows back (9 sites):

```
FAIL: 9 nil-root guards in skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go; REQ-558 pinned exactly 8 — one per consumer of the rollback handle, and no more:
  297:	if root == nil {
  998:	if root != nil {
  1009:				if root != nil && privateStateStillOriginal(root, state) {
  1026:				if root != nil && privateStateStillOriginal(root, state) {
  1114:			if root != nil {
  1143:		if !recorded || root == nil {
  1176:	if root == nil {
  1193:	if root == nil {
  1284:	if root == nil {
LOCKIN_EXIT=1
```

Lock-in RED direction 2 — a load-bearing guard is deleted (`trackedPublicationStillOwned`, 7 sites):

```
FAIL: 7 nil-root guards in skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go; REQ-558 pinned exactly 8 — one per consumer of the rollback handle, and no more:
  297:	if root == nil {
  998:	if root != nil {
  1009:				if root != nil && privateStateStillOriginal(root, state) {
  1026:				if root != nil && privateStateStillOriginal(root, state) {
  1114:			if root != nil {
  1143:		if !recorded || root == nil {
  1176:	if root == nil {
LOCKIN_EXIT=1
```

Both probes restored the file afterwards (`TARGET_RESTORED`). Lock-in GREEN at head:

```
Audit lock-in regressions passed.
LOCKIN_EXIT=0
```

The request's reproduce command at my head:

```
$ rg -n 'root [=!]= nil' skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go; rg -l 'OpenRoot|rollback root is unavailable|rooted filesystem handle is unavailable' skills/do-work/tools/do-work-cli/internal/gittransaction/*_test.go || echo 'NO TEST covers any nil-root branch'
297:	if root == nil {
998:	if root != nil {
1009:				if root != nil && privateStateStillOriginal(root, state) {
1026:				if root != nil && privateStateStillOriginal(root, state) {
1114:			if root != nil {
1143:		if !recorded || root == nil {
1176:	if root == nil {
1193:	if root == nil {
NO TEST covers any nil-root branch
```

Eight sites, not the one the Red-Green Proof asked for. The second half of that output is
unchanged and stays true: no test in the package reaches a no-handle branch, which is why the
lock-in has to carry the floor.

Gate — `DO_WORK_GATE_ROOT=<worktree> bash .../gate.sh`:

```
GATE_EXIT=0
Maintainer verification passed.
```

The gate's default invocation skips the heavy tier (`SKIP: staged skills, updater, and installer
probes require maintainer-verify.sh --heavy`), so the two heavy-verification tests that fail for
environmental reasons at the branch point did not run in this pass and are not mine either way.
ShellCheck on the changed lock-in file reports five pre-existing style findings at lines 28, 48,
118, 229 and 230; the new block is clean.

## Found and not fixed

**F1 — a live nil-handle panic on the rollback path, present at HEAD with all nine guards.**
`rollbackFailure:1032` calls `quarantineAndRollbackPrivate(root, state, published)` with the
possibly-nil handle and no check, and that function dereferences it at line 1213,
`root.Mkdir(directory, 0o700)`. Any transaction with an identity-recorded private untracked
target panics mid-rollback if the open at line 994 fails. Two of the three traces reproduced it
against the unmodified file, under two different producer conditions. It is out of this REQ's
scope ("do not fix nearby code"), and it is not made worse by this change: the guard I deleted
sat downstream of line 1213 and could never have run on that path. The doc comment on
`rootedOpenSnapshot` names this caller so the next reader does not mistake the enumeration for a
complete set of no-root branches.

**F2 — two helpers rely entirely on their call sites.** `privateStateStillOriginal`
(`root.Lstat` at 1250) and `rootedCreateRegular` (`root.MkdirAll` / `root.OpenFile` at 1315-1319)
take the same handle and never test it. The request cited that as evidence of inconsistency
arguing for fewer guards; it argues the other way. They are safe only because the guards at 1009,
1026 and 1176 stand in front of them, so any future edit to those call sites breaks them silently.

**F3 — the real remedy the request half-named is still undone.** Deciding once at line 994 —
either abandon the rooted half of rollback with a typed finding, or run it under a proven handle —
would make all eight surviving guards genuinely dead and would subsume F1. That is a change to
`rollbackFailure`'s error handling, not a guard deletion: today it keeps going with a nil handle on
purpose so it can still unstage paths and return a typed incomplete rollback. It needs the first
test in this package that exercises a no-handle rollback, and the practical way to write one is to
chmod the worktree root to 0111 from inside the mutation callback and skip when running as root
(git only needs search permission, so it keeps working while `os.OpenRoot` fails). Worth its own
REQ; it changes reported outcomes for every failure mode that reaches rollback.

**F4 — no changelog entry.** `git_transaction.go` ships, so per `_dev/primes/prime-releases.md`
the commit is a release and needs a `CHANGELOG.md` entry plus the byte-identical mirror at
`skills/do-work/CHANGELOG.md`. Both files are outside this REQ's write set and outside the file
list I was given, so the finalize step still owes that entry.
