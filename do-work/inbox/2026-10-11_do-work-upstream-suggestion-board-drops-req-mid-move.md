# Upstream suggestion for `knews2019/skill-do-work` — the board must not drop a REQ that moves between its walk and its read

**How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
consumer repo it was authored in. Observed against **v0.305.114**; all line numbers are from that tag. Still
present in v0.305.118.

---

## Request

The live board (`queue-kanban serve`) builds its model in two steps: it walks `do-work/` to list every
`REQ-*.md`, then reads each listed path. A claim moves a REQ from `queue/` to `working/` with one atomic
`os.Rename` in between those two steps often enough to matter. When that happens the read of the old path
fails with "no such file", the loop skips the file as "unreadable", and for that one response the REQ exists
nowhere: not Pending, not Claimed, not Done. Every consumer of the model then reports a real-looking problem
from a false premise. The dependency resolver prints a data warning that a sibling REQ "depends on REQ-NNN,
which is not in the do-work tree — treated as unmet". The verify pass, which reads `CHECKPOINT.md` fresh
and checks its ids against the model, prints a `checkpoint-names-missing-req` finding whose remedy tells the
reader to hand-edit the checkpoint. The page does not poll, so the user keeps looking at that response until
they reload. The user's question was "why is the REQ missing?", and the honest answer was "it never was".

Please make the board treat a file that vanishes between the walk and the read as "the tree moved under
me", not as an unreadable file:

1. **`skills/do-work-board/tools/queue-kanban/model.go:468-474`** (`buildBoard`, the parse loop): when
   `parseRequestTicket` fails with `fs.ErrNotExist`, re-run `enumerateDoWorkTree` once and parse again from
   the fresh listing. One retry is enough: a lifecycle transition renames one file, and the second walk lands
   after it. Keep the existing skip for any other read or parse error.
2. **`skills/do-work-board/tools/queue-kanban/serve.go:344-390`** (`refreshBoardData`): do not store the
   mtime fingerprint for a build that needed the retry, so the next request rebuilds instead of trusting a
   fingerprint taken before the move. Today the fingerprint comes from a separate walk at line 353, before
   `buildBoard` walks again at `model.go:452`, so the cached fingerprint and the cached model can describe
   two different trees.
3. **`skills/do-work-board/tools/queue-kanban/verify.go:731-747`** (`appendCheckpointGhostFindings`): no
   wording change, but this finding and its remedy ("edit that id's own entry out of the checkpoint") are
   what the user sees, so the acceptance check below exercises it.

Nothing here changes the CLI or the claim contract. `ApplyPlan` already does the right thing in the right
order: it replaces the request file in place (`internal/requeststate/state_apply.go:59`), renames it with one
`os.Rename` (`state_apply.go:66`), and writes `CHECKPOINT.md` last (`state_apply.go:111`). The window is on
the reader's side only.

## What happened, in one consumer repo (timeline)

A consumer repo running two work lanes in parallel, suite v0.305.114 throughout, with the live board open in
a browser tab.

| When (UTC) | REQ | Event |
| --- | --- | --- |
| Oct 11, 13:41:33 | A | A run orchestrator claimed REQ A (a confirm-style REQ that an earlier clarify had returned to the queue). The CLI replaced the file in place, renamed it from `queue/` to `working/`, wrote the checkpoint entry, and committed `[REQ A] claim request lifecycle`. |
| Oct 11, 13:41:34 | B | A second lane claimed REQ B the same way, one second later. |
| Oct 11, 13:41 | A, B, C | The user's board tab had requested `board-data.js` during the first claim. That response showed Claimed 0, Pending 9 with B still pending, A in no column, one data warning "REQ C depends on REQ A, which is not in the do-work tree — treated as unmet", and one verify finding "REQ A: checkpoint names missing req: do-work/CHECKPOINT.md names REQ A, which does not exist in queue/, working/, or archive/". |
| Oct 11, 13:42:44 | A, B | A fresh request to the same server returned A and B in Claimed, A's `treeSection` as `working`, no data warning, no finding. Nothing in `do-work/` had changed since 13:41:34. |

The description of the screenshot: the top of the page carries a red "1 data warning" box and an amber
"VERIFY 1 finding" strip naming REQ A, the Claimed column says "Nothing here", and the Pending column lists
B first. Cost: no data was lost and no file was touched, but the user asked a session to explain the
screenshot, the session spent ten minutes proving a one-second race, and until it did the verify remedy was
inviting a hand edit of the checkpoint that would have deleted a live claim.

## Where the behaviour lives today

- `skills/do-work-board/tools/queue-kanban/model.go:452` (`buildBoard`): "discovered, enumerateError :=
  enumerateDoWorkTree(repoRoot)". The walk lists paths and tree sections; nothing is read yet.
- `skills/do-work-board/tools/queue-kanban/model.go:469-472` (`buildBoard` parse loop): "ticket, parseError
  := parseRequestTicket(reference.AbsolutePath, reference.TreeSection); if parseError != nil { continue //
  best-effort: skip an unreadable REQ file }". A missing file and a corrupt file get the same silent skip.
- `skills/do-work-board/tools/queue-kanban/model.go:735` (`parseRequestTicket`): "contentBytes, readError :=
  os.ReadFile(filePath)". This is the read that races the rename.
- `skills/do-work-board/tools/queue-kanban/serve.go:353` and `:373` (`refreshBoardData`): the fingerprint
  walk ("currentFileMtimes := buildTreeMtimeFingerprint(discovered)") and the build walk inside `buildBoard`
  are two separate walks, and line 384 stores the first against the result of the second.
- `skills/do-work-board/tools/queue-kanban/model.go:1867` (dependency resolution): "%s depends on %s, which
  is not in the do-work tree — treated as unmet". Correct for a truly missing id, wrong here.
- `skills/do-work-board/tools/queue-kanban/verify.go:739-746` (`appendCheckpointGhostFindings`): "if _,
  exists := board.RequestsById[mentionedId]; exists { continue }" then the `checkpoint-names-missing-req`
  finding. It reads the checkpoint fresh on every request but checks it against the cached model, so it is
  the surface that turns the dropped file into an instruction.
- `skills/do-work-board/tools/queue-kanban/web/board.js:68`: "setInterval(refreshTickingSurfaces, 1000)" is
  the only timer. It re-renders relative times; it never refetches `board-data.js`, so a transient response
  stays on screen until the user reloads.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go:59-67` and `:111` (`ApplyPlan`):
  replace in place, `os.Rename`, then checkpoint. The sibling that already handles its half correctly.

## Proposed direction

Option A (recommended, board only): retry the walk once on a vanished file and refuse to cache that build.
Option B: make the client poll `board-data.js` every few seconds so a transient response heals on its own.
B hides the symptom and costs a request per tab per interval; A removes the false model, which is what the
dependency warning and the verify finding are built on. Do A. B can be its own suggestion if wanted.

**`skills/do-work-board/tools/queue-kanban/model.go`, `buildBoard`, replace the parse loop at 467-474:**

> Pull the walk-then-read into one helper, `readRequestTickets(repoRoot) (tickets, discovered,
> treeMovedDuringRead bool, err)`. It walks, reads every listed file, and if any read fails with
> `errors.Is(err, fs.ErrNotExist)` it walks and reads once more from scratch. It returns
> `treeMovedDuringRead = true` when the retry happened, and keeps the existing silent skip for every other
> read or parse error. `buildBoard` stores the flag on `Board` as `TreeMovedDuringRead` and, when it is
> set, appends one warning: "do-work/ changed while the board was reading it; this response was rebuilt from
> a second walk". The comment on the skip should say which errors it is for: "a corrupt or unreadable file,
> never a file that moved; a moved file re-walks above".

**`skills/do-work-board/tools/queue-kanban/serve.go`, `refreshBoardData`, after line 373:**

> If `board.TreeMovedDuringRead` is set, serve the result but leave `cachedFileMtimes` empty, so the next
> request rebuilds unconditionally. Better still, take the fingerprint from the walk `buildBoard` actually
> used (return `discovered` alongside the board) so the cache key and the model always describe the same
> tree; that also removes one of the two walks per cache miss.

**Test, `skills/do-work-board/tools/queue-kanban/board_synthetic_test.go` style:**

> Give `readRequestTickets` an injectable `readFile func(string) ([]byte, error)`. The test writes a REQ in
> `queue/`, and the injected reader, on its first call for that path, renames the file into `working/` and
> then returns the real `os.ReadFile` error. Assert: the REQ is on the board with `TreeSection == "working"`,
> `TreeMovedDuringRead` is true, a REQ that lists it in `depends_on` carries no "not in the do-work tree"
> warning, and `appendCheckpointGhostFindings` with a checkpoint naming it reports nothing. A second test
> keeps a genuinely corrupt file and asserts it is still skipped with the flag false.

## Acceptance check

- The parse loop in `model.go` no longer has a bare `continue` on every error: `grep -n 'fs.ErrNotExist'
  skills/do-work-board/tools/queue-kanban/model.go` finds the retry branch.
- The two tests above pass, and `go test ./...` under `skills/do-work-board/tools/queue-kanban` stays green.
- With `queue-kanban serve` running, claim a REQ while `curl -s localhost:8090/board-data.js` runs in a tight
  loop (for example 50 requests during one `do-work` claim). No response lists the REQ in zero columns, and
  no response carries the `checkpoint-names-missing-req` finding for it.
- Unchanged: a REQ file that is corrupt or unreadable for any reason other than "no such file" is still
  skipped silently; a checkpoint id with no file anywhere in `queue/`, `working/` or `archive/` still raises
  `checkpoint-names-missing-req`; a dependency on an id that truly does not exist still warns.

## Out of scope

- Any change to `ApplyPlan` or the claim commit. Its order is already right and atomic per file.
- Client-side polling or a server push for the board page (Option B); separate suggestion if wanted.
- The verify remedy wording for a real ghost id stays as it is.
