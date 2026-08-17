---
id: REQ-218
title: Ratchet the tools download and correct the stale gitattributes claim
status: pending
created_at: 2026-08-17T17:16:28Z
user_request: UR-049
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-217]
maintenance: false
related: [REQ-216, REQ-217]
batch: resilient-upstream-fetch
effort_estimate: trivial
write_set:
  - _dev/tests/prescribed-shell-canonicalization.sh
  - .gitattributes
---

# Ratchet the Tools Download and Correct the Stale Gitattributes Claim

## What

Extend `_dev/tests/prescribed-shell-canonicalization.sh` so `skills/*/tools/` cannot hand-roll the download
primitive a fifth time, and correct the `.gitattributes` header comment, which documents a
belt-and-suspenders `tar --exclude` net that does not exist.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read the ratchet, both `tools/` scripts as REQ-217 leaves them, and `.gitattributes`. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

Feedback-brief finding **F2**: the REQ-167/REQ-171 canonicalization campaign swept `actions/` and `scripts/`
and never named `tools/` — both archived REQs' `write_set` fields confirm it. The ratchet at
`_dev/tests/prescribed-shell-canonicalization.sh:9-21` enumerates `scripts/` entries only, so it cannot see
the two `tools/` entry points that restated the download primitive. Without this REQ, the same duplication
returns the next time a tool script needs to fetch something.

Separately, finding **F3.1**: the `.gitattributes` header at lines 9-11 claims the install/update `tar`
command "also passes `--exclude` for do-work/kb/ai-reports/_dev plus .vscode/decisions as a
belt-and-suspenders fallback". No `--exclude` exists at any of the three extraction sites today —
`do-work-update.sh:75`, `install-do-work-suite.sh:28`, `install-do-work-suite.sh:178`. `export-ignore` is the
only thing holding that line, exactly as the comment's own later paragraph warns. A future reader trusting
the stale sentence would under-estimate how load-bearing `export-ignore` is.

## Detailed Requirements

- Add `tools/fetch-upstream-archive.sh` to the ratchet's enumerated prescribed-script list.
- Add a check that no file under `skills/*/tools/` contains a bare `curl -fsSL -o`. This needs **two**
  documented exemptions, not the one the brief proposed:
  1. the `BOOTSTRAP` heredoc in `install-do-work-suite.sh` — nothing is installed yet at that point, so its
     inline `curl` is correct layering; and
  2. `fetch-upstream-archive.sh` itself, which is the primitive's tool-side home.
- Correct `.gitattributes:9-11`. **Preferred: delete the stale sentence** rather than restore the
  `--exclude` net — `export-ignore` is what actually holds the line, the comment's later paragraph already
  says so, and the `CLAUDE.md`/`AGENTS.md` paragraph at lines 35-43 explains why a `tar --exclude` fallback
  is deliberately *not* wanted for those two (basename matching at any depth is the same trap as `diff -x`).
  Restoring the net is the acceptable alternative if the builder finds a reason to prefer it; silently
  leaving the claim stale is not.
- Leave the load-bearing paragraph at `.gitattributes:13-20` and the `/do-work` + `/kb` lines untouched —
  `_dev/tests/contract-regressions.sh` asserts both lines exist.

## Constraints

- Keep the exemptions **stated in the check itself**, not implied — a future reader must be able to see why
  each one is allowed without archaeology.
- Prefer keying the check on the condition rather than a hand-maintained path list where the shell allows it;
  see `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale.

## Dependencies

`depends_on: [REQ-217]` — the new check is RED until REQ-217 removes the bare `curl` from
`do-work-update.sh:72` and `install-do-work-suite.sh:173`, and the ratchet's new list entry requires the
fetcher to exist and be executable.

## Builder Guidance

Certainty: **Firm.** Small and mechanical. The judgment call is the `.gitattributes` wording — delete the
stale sentence unless there is a concrete reason to restore the `--exclude` net instead, and say which you
chose and why.

## Red-Green Proof

**RED prompt/case:** Run `bash _dev/tests/prescribed-shell-canonicalization.sh` against the tree as it stands
**before** REQ-217 lands, with the new "no bare `curl -fsSL -o` under `skills/*/tools/`" check in place.
**Why RED now:** `skills/do-work/tools/do-work-update.sh:72` and
`skills/do-work/tools/install-do-work-suite.sh:173` both contain a bare `curl -fsSL -o` outside the bootstrap
heredoc, so the check reports two violations and the script exits non-zero. The ratchet as it stands today
enumerates `scripts/` only and reports nothing.
**GREEN when:** With REQ-217 landed, the same run exits `0` — both former sites now call the fetcher, and the
only surviving `curl -fsSL -o` occurrences under `skills/*/tools/` are the two documented exemptions.
**Validation:** User confirmed (capture approved the three-REQ split and the two-exemption adjustment).

## Finding-Closure

Origin: external feedback brief, triaged via `do-work-toolbox validate-feedback` — finding F3.1 and proposal
C4 (adjusted from one exemption to two), verdict **Accept**. Surface-cost: **earned** — the ratchet is the
mechanism that already exists for exactly this class, the incident is the four-copy duplication F1
documented, and the added surface is one check in an existing test. The `.gitattributes` half is a
correction/deletion, so **N/A** there. Named regression check: the new `skills/*/tools/` clause in
`_dev/tests/prescribed-shell-canonicalization.sh`.

## Full Context

See `do-work/user-requests/UR-049/input.md` for the complete verbatim brief.
