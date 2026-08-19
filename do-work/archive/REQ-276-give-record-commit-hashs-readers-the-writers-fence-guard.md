---
id: REQ-276
title: Give record-commit-hash's readers the same fence guard its writer already has
status: completed
created_at: 2026-08-18T23:38:58Z
claimed_at: 2026-08-19T20:10:08Z
completed_at: 2026-08-19T20:29:22Z
commit:
status_changed_at: 2026-08-19T13:45:20Z
user_request: UR-056
addendum_to: REQ-267
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
route: A
estimate:
  p50_active_minutes: 10
  confidence: medium
  calculated_at: 2026-08-19T20:10:08Z
  basis:
    - Route A
    - 2-file write set
    - 4 acceptance criteria
write_set:
- skills/do-work/tools/checks/record-commit-hash.sh
- _dev/tests/record-commit-hash-guards.sh
---

# Give record-commit-hash's Readers the Same Fence Guard Its Writer Already Has

## What

`skills/do-work/tools/checks/record-commit-hash.sh` guards against an unterminated frontmatter fence **on the write path only**. Its rewrite awk (`:472-487`) buffers and gates on `frontmatter_closed`. Its **read** helpers (`:108-122`, called at `:247` and `:318-338`) still scan to end of file — so on a file whose opening `---` is never closed, `--verify` can read a body `commit:` line as the frontmatter one.

The script's own header at `:103-106` claims both halves use one parser. They do not.

**Why this is worth more than its severity suggests:** `record-commit-hash.sh --verify` is the last check every REQ in this pipeline passes through, and it exists because free-form frontmatter edits once truncated six archived REQs to zero bytes in a consumer repo. A guard that is asymmetric between its reader and its writer is exactly the shape that lets a corrupted file report as verified.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

Scope was narrowed by the user during clarify (see `## Decision Record` below). Refuse the file at
the door rather than teaching each reader a guard.

- Add one `require_closed_frontmatter` helper — an awk that fails when a file opens `---` on line 1
  and never closes it — and call it at exactly **two** sites:
  - on `$request_file` at startup, beside the existing CRLF refusal (`:310`), which covers every
    reader call on that file (`:318`, `:322`, `:328`, `:336-338`, `:566-570`) without editing any of them;
  - on the parent blob fetched by `--verify` (`:247`), the only reader input that is not `$request_file`.
- The header comment at `:103-106` stops claiming the readers and the writer share one guard, and says
  what is actually true: an unterminated fence is refused up front, so no reader ever sees one.
- Two lock-in cases: a file whose fence never closes and whose body carries a column-0 `commit:` line is
  refused on the write path and on the `--verify` path. Each names the failure it pins.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

**Explicitly out of scope** — do not do these under this REQ:

- Refactoring the three readers to share one parser with the writer. A single precondition is smaller
  than threading a guard through each awk, and buys the same property.
- The repo-wide sweep of every shipped awk or shell frontmatter parser. Different scope, different
  write set; the user declined to capture it as a separate REQ during this clarify.

## Context

Discovered by REQ-267's builder while fixing the same defect class in `repair-req-timestamps.sh`, and independently confirmed by REQ-267's reviewer during its Restatement Sweep — which noted the repairer's new buffer-and-gate follows a pattern `record-commit-hash.sh`'s *writer* already established, making the reader's omission an inconsistency inside one file rather than a missing feature.

Classified `[normal]`, so it enters as `pending-answers` under the Discovered Tasks consent flow.

## Open Questions

- [x] REQ-267's builder found that `record-commit-hash.sh` guards against an unterminated frontmatter fence when writing but not when reading, so `--verify` can read a body `commit:` as the frontmatter one — on the last check every REQ in the pipeline passes through. Should I process this as a new task? → Yes, add to queue — with the scope simplified to one up-front guard (see the Decision Record).
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the writer's guard means the pipeline never produces such a file itself.

## Decision Record

- **[2026-08-19] Scope narrowed to one up-front refusal.** The user reviewed this REQ during
  `do-work clarify`, judged the original requirements over-complicated, and chose the simplified
  form now in `## Requirements`. Reasoning: the three readers already run the same four-line awk;
  the defect is that none of them gates on the closing fence, and a single precondition checked once
  at startup fixes every reader at once. Sharing one parser with the writer would be a larger diff
  for the same property.
- **Put out of scope by the same decision:** the reader/writer parser-sharing refactor, and the
  repo-wide sweep of other shipped frontmatter parsers. The user was offered the sweep as its own
  REQ and declined it, so it is not captured anywhere — reopen deliberately if wanted.
- **Decided by:** user, via `do-work clarify`.

---

## Triage

**Route: A** - Simple

**Reasoning:** The user narrowed the scope during clarify and the Requirements now name the exact helper to add, the exact two call sites, the header comment to correct, and the two lock-in cases. Nothing is left to discover.

**Planning:** Not required

*Skipped by work action*

## Decisions

- **D-01**: The two lock-in cases went into `_dev/tests/record-commit-hash-guards.sh`, not
  `_dev/tests/prescribed-shell-scripts-behavior.sh` as `write_set` said. The guards file is
  this script's dedicated behavioral probe suite — it builds a throwaway git repo, runs the
  real script against real REQ fixtures, and already holds nineteen probes including the
  existing unterminated-frontmatter one these two sit beside. The behavior suite greps
  prescribed shell across the whole repo and has no fixture repo for `--verify`'s parent
  blob at all, so the cases could not have been written there without duplicating the
  harness. `write_set` corrected in place. DECIDE & STATE.
- **D-02**: `require_closed_frontmatter` is defined beside the three readers (above the
  `--verify` block), not at the startup guard where the Requirements' first bullet
  describes its call. The parent-blob call site is inside `--verify`, which exits long
  before the startup guards are reached — a definition at the startup site would have been
  missing exactly where the second call needs it. First attempt did that and would have
  turned every `--verify` run into a `command not found` failure. The two call sites are
  the two the Requirements name; only the definition point moved. DECIDE & STATE.
- **D-03**: The parent-blob refusal exits 1 with a `FAIL:` line rather than 2. Exit 2 is
  this script's bad-*input* code and the input here is the caller's file, which was fine;
  what is unreadable is the committed parent, which is a verification failure. It reads
  like the other `--verify` failures beside it. DECIDE & STATE.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/checks/record-commit-hash.sh` (modified)
- `_dev/tests/record-commit-hash-guards.sh` (modified)

**What was done:** Added one `require_closed_frontmatter` precondition — an awk that fails
when a file opens `---` on line 1 and never closes it — defined beside the three readers it
protects, and called at exactly the two sites the Requirements name: on `$request_file` at
startup beside the CRLF refusal, which covers every reader call on that file without
editing any of them, and on the parent blob `--verify` fetches, the only reader input that
is not `$request_file`. The reader block's header comment stops implying the readers carry
the writer's guard themselves and says what is actually true: they are safe because the
file is refused at the door, and a fifth reader inherits that with nothing to remember. Two
lock-in probes pin both paths; the existing unterminated-frontmatter probe was updated
because the refusal now fires earlier and as bad input.

## Qualification

Passed — 2 files verified in the diff, 4 requirements traced, P-A-U confirmed against the
diff. The reviewer independently re-enumerated every reader call site and confirmed all
eight sit after one of the two guards.

## Testing

**Tests run:** `_dev/tests/record-commit-hash-guards.sh` (via the path-rewrite invocation
`contract-regressions.sh` uses), then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing (maintainer-verify exit 0)

**Red-green validation:**
- Probe 20 `unclosed fence` (write path, 5 assertions): ✗ before → ✓ after
- Probe 21 `verify parent fence` (`--verify` parent blob, 3 assertions): ✗ before → ✓ after

Run RED by stashing only the script and keeping the probes. The RED output is the defect
itself rather than a bare failure: on the write path the old script answered
`frontmatter has 2 'commit:' lines — ambiguous`, which is the readers having counted the
body's fenced `commit: deadbee` as a schema field; on the `--verify` path it named
`commit: deadbee   # body prose, MUST NOT be rewritten` as "HEAD^'s own frontmatter line"
and computed the expected removal from it.

**New tests added:**
- Probe 20 and Probe 21 in `_dev/tests/record-commit-hash-guards.sh`

**Existing tests updated (cross-REQ impact):**
- Probe 7 `unterminated` (from the original write-back guard work): its exit expectation
  moved 1 → 2. The property it pins is unchanged — nonzero, and the file byte-identical —
  but the refusal now fires at the door as bad input rather than at the writer's END flush,
  which puts it in the same class as the CRLF refusal in Probe 8 directly below and
  asserted the same way. Comment updated to say so.

*Verified by work action*

## Review

**Reviewer:** independent agent, orchestrated mode against the working-tree diff.

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 88% → higher after the three Minor fixes |
| Test Adequacy | 88% → higher after the Probe 20 comment fix |
| Scope Discipline | 100% |
| Risk | None |
| **Acceptance** | **Pass** |
| **Overall** | **94%** |

**No Important findings.** The reviewer enumerated all eight reader call sites and confirmed
none is reachable unguarded, exercised the awk directly against seven shapes on this
machine's BSD awk and bash 3.2 (closed, unclosed, no leading fence, empty, bare `---` with
and without a trailing newline, CRLF), confirmed `exit` inside a rule reaches END and END's
own `exit 1` wins, confirmed the second `git cat-file blob` reads the same immutable object
with no TOCTOU and cannot misfire on the legitimate no-parent-file path, and confirmed both
exit-code choices against the script's own contract. It also verified scope discipline
against the user's two explicit exclusions: the readers are byte-for-byte unchanged and no
other frontmatter parser in the repo was touched.

**Three Minor findings, all fixed here** (all inside the write set, all one comment or one
line):

- The script's top-of-file **exit-code contract** at `:49-50` listed the exit-2 causes and
  never learned the new one. That is the co-location rule and precisely what the Restatement
  Sweep is for; I updated the reader-block header the Requirements named and missed the one
  they did not.
- The write-awk's `exit 4` **unterminated branch** is now reachable only through a narrow
  race (`$request_file` rewritten between the startup check and the awk read), but its
  comment still read as the primary defense. Kept the branch — it is what makes that race
  non-destructive — and said what it now is.
- **Probe 20's comment overclaimed its own fixture.** Its seed carries a frontmatter
  `commit:` as well as the body bait, so the pre-fix script failed it with
  `frontmatter has 2 'commit:' lines`, not with a silent misattribution. The assertions
  always discriminated correctly (exit 2 + "never closed" versus exit 1 + "ambiguous"), but
  the comment now states what the fixture actually proves — that the duplicate count *is*
  the readers counting a body line as schema — and points at Probe 21 as the one that
  exercises the silent case end to end.

**Nit, addressed in a comment rather than in code:** `require_closed_frontmatter` is
CRLF-blind, because `$0` keeps the trailing `\r` so a CRLF `---` is not the delimiter. For
`$request_file` this is moot — the CRLF refusal runs first. The parent blob has **no** CRLF
check at all, before or after this REQ. That gap is pre-existing, out of this REQ's two
named call sites, and fails safe: a CRLF parent reads as having no frontmatter `commit:`, so
`--verify` expects an insert and reports a mismatch rather than passing something wrong. It
is now stated at the guard instead of left for the next reader to rediscover, and
deliberately not minted as a REQ — a false FAIL on an already-refused file shape is not
worth queue traffic.

**Nit, no action:** the REQ's Requirements cite `:566-570` as reader calls on
`$request_file`; those calls read `$temp_file`, the write-awk's own output, which is
well-formed by construction. An imprecise citation in the REQ, not a defect in either.

## Lessons Learned

**What worked:** Running the new probes RED by stashing *only* the script and keeping the
probes. The failure output named the defect in the script's own words — a duplicate-count
error that was really the reader misreading body prose — which is stronger evidence than a
bare non-zero exit, and it is what let the reviewer catch that one probe's comment described
a different failure from the one its fixture produces.

**What didn't:** Placing the helper at the call site the Requirements described. `--verify`
returns long before the startup guards, so the definition would have been missing exactly
where the second call needed it — every `--verify` run would have died on `command not
found`, on the last check every REQ in the pipeline passes through. Caught by running the
suite, not by reading. **A shell function used by an early-exiting branch has to be defined
above that branch, and "beside the guard it belongs to" is not the same as "before the code
that calls it."**

**Worth knowing:** this file has two independent readers of the same shape — `$request_file`
and `--verify`'s parent blob — and only the first is covered by the startup preconditions.
Every input guard added at the top (CRLF, 0 bytes, symlink, now the fence) covers exactly
one of the two. Anything that must hold for both needs a second explicit call, which is why
the CRLF gap on the parent path exists at all.

## Orientation

`record-commit-hash.sh` — the guard on the last write every REQ in this pipeline passes
through — now refuses a file whose frontmatter fence never closes, before any of its three
readers can scan past the missing delimiter and take a body `commit:` line for the schema
field. One precondition at the door rather than a guard threaded through each reader, so a
fourth reader inherits it with nothing to remember. Lives in the REQ-bookkeeping guard
tooling (`_dev/primes/prime-shell-commands.md`). No map change. Prime spot-check:
`prime-shell-commands.md`'s referenced paths all still exist.
