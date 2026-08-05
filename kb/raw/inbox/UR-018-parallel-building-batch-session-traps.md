---
source_type: req_lesson
req_id: UR-018
req_path: do-work/archive/UR-018/
date: 2026-08-05
domain: general
module: tools/checks
tags: [tooling, contract-tests, worktree-dispatch, traps]
---

# Lessons from UR-018: Traps the parallel-building batch session already hit

## What the REQ was about

UR-018 re-grained session ownership to "claim anywhere, one releaser" — parallel building across
checkouts (REQ-094 through REQ-104, plus REQ-108/109; ADR-018). Mid-batch, the session wrote itself a
handdown at `do-work/HANDDOWN-UR-018.md` whose *Traps this session already hit* section recorded the
operational potholes it had already fallen into: exact flag orders, exact command forms, the pinned
contract phrases near the batch's files, and the two process failures that produced its follow-up REQs.
That handdown was deleted as a stale session artifact in commit `b1792d0` once the batch shipped; the
traps section outlived the handdown's status table, so it is recovered here. **This is repo-internal
maintainer knowledge about do-work's own tooling — not consumer-facing skill documentation.**

Every trap below was re-verified against the working tree before being written down, and the
verification method is stated per trap so a future reader can re-run it rather than trust this file.
That is not a formality: an adversarial review of this entry re-ran them and found the stale-binary
bullet's evidence wrong on every point, which is why that bullet now records what does **not** prove
staleness. Re-run rather than trust applies to this file too.

## What didn't work

- **Builders paraphrasing a canonical condition at echo sites.** Both of this batch's mid-flight
  follow-ups traced to the same cause: a builder restated a canonical rule in its own words at a site
  that merely *echoes* the rule, and the restatement drifted from the canonical wording. The fix that
  worked is a brief-level instruction — **make builders QUOTE the canonical wording at echo sites**
  rather than summarize it. *Verified still live:* grepped `crew-members/*.md`, `actions/work.md` and
  `actions/work-reference.md` for any echo-site or paraphrase rule — there is none. Nothing in the
  shipped instructions codifies this, so it remains unwritten advice a brief has to carry explicitly.

- **Trusting `queue-kanban verify` findings from a binary nobody rebuilt.** A stale compiled binary
  reported a ghost-REQ false positive that had already been fixed in 0.169.9 (the checkpoint mention
  scan matching the `REQ-0` prefix of a quoted shell glob `REQ-0[0-9][0-9]-*.md`). *Scope corrected on
  recovery — narrower than the original claimed:* the two shipped call sites that run `verify` already
  rebuild in the same command, and both did so when this trap was first written —
  `actions/forensics.md` Check 14 and `actions/work.md` Step 9 each prescribe
  `(cd <skill-root>/tools/queue-kanban && go build -o queue-kanban .) && … verify …`. The residual
  exposure is a **hand-run `verify` outside those prescribed blocks**, which is how the original
  incident happened. Rebuild first:

  ```
  cd tools/queue-kanban && go build -o queue-kanban .
  ```

  Two ways of *testing* for staleness that do not work, recorded because this recovery pass tried both
  and drew a false conclusion from each. Comparing the binary against a fresh `go build` byte-for-byte
  always differs once HEAD has moved, and `go version -m` reports `vcs.time` — the **commit's**
  timestamp, not the binary's mtime. Go stamps `vcs.revision`, `vcs.time` and `vcs.modified` into every
  build, so that delta is the stamp rather than the code: `go build -buildvcs=false` twice is
  byte-identical, and a one-commit-old binary produces identical `verify` output. Compare mtimes against
  the `.go` sources if you must, but rebuilding is cheaper than proving you need to.

## Worth knowing

- **`tools/checks/record-commit-hash.sh` recognizes `--verify` in first position only.** The flag order
  is `--verify <req-file> <hash>`. Putting it last — `<req-file> <hash> --verify` — is not accepted.
  *Verified by running both forms:* the script tests `[ "${1:-}" = "--verify" ]` at line 64 and shifts,
  then hard-requires exactly two remaining arguments, so a trailing flag makes `$#` equal 3 and the
  script exits 2 with `usage: … [--verify] <req-file> <hash>   — exactly 2 arguments; quote a path
  containing spaces`. The good news, worth knowing because it bounds the damage: the arg-count check
  runs before any write path, so a mis-ordered verify invocation cannot be mistaken for a write — you
  get a usage error, not a second rewrite of the REQ.

- **Pass the contract suite to `tools/checks/preflight.sh` as a quoted `bash …` string, never as a bare
  path.** Use `"bash _dev/tests/contract-regressions.sh"`. The bare path hits `Permission denied`
  because the suite is not executable — mode `100644` in the index (`git ls-files -s`), so this is true
  in every checkout, not a local artifact. *Verified by reproducing it* in a scratch directory holding a
  mode-644 copy: preflight prints `WARN: could not run the test command — no baseline recorded`, the
  underlying `Permission denied`, and the hint `(pass the command as separate words, or as one quoted
  string — see the usage header)`. Two annotations on the original trap: preflight now *detects* this
  (exit 126/127 records `"launched": false` and writes no failures file, so Step 6.5 cannot read a
  fictional red baseline as pre-existing), and it deliberately supports the quoted form by handing a
  single whitespace-containing argument to `sh -c` — which is also why `"cd app && npm test"` works.
  The trap is therefore loud rather than silent now, but it still costs you the baseline: preflight
  always exits 0, so a mis-invocation reads as a passing preflight run unless you read the WARN lines.

- **Contract-suite pinned phrases near this batch's files — reword AROUND them, never through them.**
  Verified present in `_dev/tests/contract-regressions.sh` at v0.174.11 by literal grep, with match
  counts: `never grow into one` (3), `absent checkpoint is ambiguous` (1), `foreign claim` (6),
  `Crash Recovery's input` (1), `claim held by` (2), `writer: <hostname>:<absolute-checkout-path>` (2),
  `entry this checkout did not write through verbatim` (2), `no entry this checkout did not write
  remains` (1). The Step-1 line-order pin is also still live but is positional rather than a phrase —
  `actions/work.md` Step 1 must read `do-work/CHECKPOINT.md` *before* the **Crash Recovery:** paragraph,
  enforced by a line-number comparison (suite line ~637, REQ-071). Treat this list as a snapshot, not an
  authority: it is a closed enumeration and will go stale. Re-derive it with a literal grep of the suite
  before a rewrite, and when a re-grain has to change a pinned sentence, update the assertion in the
  same commit rather than weakening other pins.

- **`actions/version.md` and `CHANGELOG.md` are serial-only and integrator-only.** Both are bumped once
  per REQ at Step 9, by the queue owner — a builder on its own worktree branch bumps neither. Changelog
  titles say what shipped, no codenames. After writing, confirm title-uniqueness and version
  monotonicity with `queue-kanban verify`. *Verified still current:* `actions/work-reference.md:360`
  states the serial-only rule for `actions/version.md` plus `CHANGELOG.md` ("one changelog entry per
  REQ, written by the owner at merge time. Unique version numbers do not make a shared prepend safe"),
  `:299` flags a second checkout prepending to `CHANGELOG.md` as a Red Flag, and `verify` is documented
  as the mechanical half of the pre-commit ritual (checking version/changelog agreement and entry-title
  reuse) and wired into `actions/forensics.md` Check 14. Note the dependency between two traps: `verify`
  is what proves the bump, and the stale-binary trap above is what makes that proof unreliable — rebuild
  before believing it.

## Back-reference

Recovered from `do-work/HANDDOWN-UR-018.md`, which was deleted along with five other stale
session-handoff files in commit `b1792d0` (v0.174.11). Retrieve the original with:

```
git show 'b1792d0^:do-work/HANDDOWN-UR-018.md'
```

Quote the `^` as shown — it is a glob operator in zsh and an escape character in cmd.exe. The handdown
carried a SUPERSEDED banner from 2026-08-04T21:38Z; its status table and per-REQ notes were already
wrong when it was deleted and are deliberately **not** reproduced here — only the traps section, which
the banner itself flagged as still accurate. The batch is fully archived at `do-work/archive/UR-018/`
(REQ-094 through REQ-104, REQ-108, REQ-109); the decision it implemented is
`decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md`.
