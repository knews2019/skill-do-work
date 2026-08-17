---
id: REQ-216
title: Teach atomic-download retry and optional credentials
status: pending
created_at: 2026-08-17T17:16:28Z
user_request: UR-049
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-217, REQ-218]
batch: resilient-upstream-fetch
write_set:
  - skills/do-work/scripts/atomic-download.sh
  - skills/do-work/docs/prescribed-shell-primitives.md
  - _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Teach Atomic-Download Retry and Optional Credentials

## What

Add transient-failure retry and opt-in GitHub credentials to the shipped download primitive at
`skills/do-work/scripts/atomic-download.sh`, so every caller inherits the resilience instead of it being
patched into one download site.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read `_dev/primes/prime-shell-commands.md` and the existing script. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

Feedback-brief finding **F1** — the four download sites in this repo are each a single `curl` attempt with
no retry and no credential support. Finding **F2** — `atomic-download.sh` is the canonical primitive for
exactly this and is bypassed by both `tools/` scripts, which is why the fix would otherwise have to be
written three or four times.

## Detailed Requirements

- Add `--retry 3 --retry-delay 2 --retry-max-time 60` to the existing `curl` at
  `skills/do-work/scripts/atomic-download.sh:26`.
- Do **not** use `--retry-all-errors` — it needs curl ≥ 7.71, and plain `--retry` already treats HTTP 429 as
  transient from curl 7.51.0 (2016). The 7.51 floor is consistent with how this repo has reasoned about curl
  versions before: `skills/do-work/CHANGELOG.md:2127` rejected `--remove-on-error` for requiring ≥ 7.83.
- When `GH_TOKEN` or `GITHUB_TOKEN` is non-empty, send `-H "Authorization: Bearer <token>"`. Opt-in only —
  absent or empty means the request goes out exactly as it does today.
- Preserve the publication contract **exactly**: private adjacent `mktemp` sibling, rename only on success,
  non-zero exit preserved, no `.download.*` scratch left behind on any path. The script runs under `set -u`
  (not `-e`) and captures `$?` explicitly into `download_status` — keep that shape.
- Update the **Atomic download publication** section of `skills/do-work/docs/prescribed-shell-primitives.md`
  (heading at line 63) so the canonical guide describes the retry and credential behavior. The heading text
  itself must not change — `_dev/tests/prescribed-shell-canonicalization.sh` greps for it with `grep -Fqx`.

## Constraints

- **Additive flags only.** The brief explicitly does not ask for a change to the publication contract, and
  neither does this REQ.
- The existing behavior case at `_dev/tests/prescribed-shell-scripts-behavior.sh:54-62` must keep passing.
  Its fake `curl` parses `-o` and ignores unrecognized flags, so the new flags do not disturb it — verify
  rather than assume.
- No new runtime dependency.

## Dependencies

None. REQ-217 depends on this.

## Builder Guidance

Certainty: **Firm.** The change is small and the shape is specified. Resist expanding it — the git fallback
belongs to REQ-217 and the ratchet to REQ-218.

## Red-Green Proof

**RED prompt/case:** A fake `curl` placed on `PATH` that returns 429 on its first invocation and succeeds on
its second — the same idiom as the existing partial-download case at
`_dev/tests/prescribed-shell-scripts-behavior.sh:54-62`. Invoke `atomic-download.sh` against it and assert
the target file is published complete.
**Why RED now:** `atomic-download.sh:26` is a single `curl -fsSL -o` attempt. The first 429 fails the script
outright, the target is never published, and the fake `curl`'s second invocation never happens.
**GREEN when:** The fetch succeeds on the retry and the target contains the complete second-attempt payload,
with no `.download.*` sibling left behind.
**Validation:** User confirmed (capture approved the three-REQ split and the C1 scope; route 1 dropped).

## Finding-Closure

Origin: external feedback brief, triaged via `do-work-toolbox validate-feedback` — findings F1, F2, and
proposal C1, all verdict **Accept**. Surface-cost: **earned** — the incident is named and reproducible
(sustained codeload 429 on `do-work update`), the added surface is one flag set plus one conditional header,
and the RED test above replays it. The named regression test is RED test 1, landing in
`_dev/tests/prescribed-shell-scripts-behavior.sh`.

## Full Context

See `do-work/user-requests/UR-049/input.md` for the complete verbatim brief.
