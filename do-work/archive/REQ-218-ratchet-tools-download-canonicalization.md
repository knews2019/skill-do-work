---
id: REQ-218
title: Ratchet the tools download and correct the stale gitattributes claim
status: completed
completed_at: 2026-08-17T19:21:16Z
commit:
claimed_at: 2026-08-17T19:16:51Z
created_at: 2026-08-17T17:16:28Z
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-17T19:17:20Z
  basis:
    - trivial short-circuit
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
- [x] **[PLAN]:** Add the fetcher to the enumerated executable list. Add a `skills/*/tools/` clause that flags any *direct* `curl` invocation rather than the narrower `curl -fsSL -o` string, since REQ-216 already proved the flag list is not fixed. Exempt quoted-heredoc bodies by condition (emitted text, not executed code), which covers the installer's BOOTSTRAP block without naming it. Delete the stale `--exclude` claim from `.gitattributes` rather than restoring the net.
- [x] **[APPLY]:** Two files, both in `write_set`. One departure from the REQ recorded as D-01.
- [x] **[UNIFY]:** `git diff --stat` → `_dev/tests/prescribed-shell-canonicalization.sh` (+38/−1), `.gitattributes` (+9/−9). ShellCheck via `maintainer-verify.sh` clean on the modified test. No debug artifacts. `git archive HEAD | tar t` re-run after the comment edits: zero `do-work/`, `kb/`, `_dev/`, `ai-reports/`, or `CLAUDE.md` entries, confirming the export-ignore rules themselves were untouched.

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

## Triage

**Route:** A — Direct to Builder

**Reasoning:** Two files, both named, both with the change spelled out. The only open question was the `.gitattributes` wording, and the REQ states its own preference.

**Confidence:** high

*Triaged by work action*

## Plan

Planning not required — Route A.

## Decisions

- **D-01: one exemption, not two — and the second is refused on the merits.** The REQ specifies two documented exemptions: the `BOOTSTRAP` heredoc, and `fetch-upstream-archive.sh` itself. The first is real and is implemented. The second is not: as REQ-217 built it, the fetcher **delegates** to `skills/do-work/scripts/atomic-download.sh` and never invokes `curl` at all, so it needs no exemption today — and granting one in advance would pre-authorize exactly the duplication this ratchet exists to prevent. If the fetcher ever needs a direct download, that should re-open this check rather than pass silently. The check says all of this in place, so a reader finds the reasoning where the exemption would have been. **Value:** the primitive's own home stays under the same rule as everything else. **Risk:** a future direct `curl` in the fetcher fails the ratchet and needs a deliberate decision (low, and that is the intent).
- **D-02: flag any direct `curl`, not the literal `curl -fsSL -o`.** The REQ's phrasing names the exact string, but REQ-216 had just inserted `--retry 3 --retry-delay 2 --retry-max-time 60` between `-fsSL` and `-o` — a string match would already have missed the shipped bootstrap line, and would miss any future flag change. The clause matches `curl` used as a command word instead, which is the actual rule: tool scripts delegate their downloads. DECIDE & STATE.
- **D-03: correct the load-bearing paragraph too, minimally.** The REQ says to leave lines 13–20 untouched, but that paragraph repeated the same false claim ("These two lines plus the tar --exclude fallbacks…"), and deleting the stale sentence while leaving its restatement standing would have left the file self-contradictory — the very defect F3.1 reports. The edit removes five words and changes nothing the contract tests assert (`/do-work` and `/kb` export-ignore lines are byte-identical, verified by re-running the suite and by re-listing `git archive HEAD`). Two adjacent sentences that implied *other* paths do have a `--exclude` fallback were corrected for the same reason.

## Implementation Summary

**Files changed:**
- `_dev/tests/prescribed-shell-canonicalization.sh` (modified)
- `.gitattributes` (modified)

**What was done:** Closed the `tools/` gap in the canonicalization ratchet and removed a stale safety claim from the archive-exclusion header.

*The ratchet.* `skills/do-work/tools/fetch-upstream-archive.sh` joins the enumerated must-exist-and-be-executable list. A new clause walks every `skills/*/tools/**.sh` and fails on any direct `curl` invocation, reporting the file and the offending line numbers. Quoted-heredoc bodies are skipped by condition — text inside `<<'DELIMITER'` is emitted for someone else to run, not executed here — which is what makes the installer's `BOOTSTRAP` block legitimate without naming it. Comment lines are skipped. The awk clause recognizes heredoc delimiters generically, so it does not depend on the block continuing to be called `BOOTSTRAP`.

*The header.* The sentence claiming the install/update `tar` passes `--exclude` for `do-work/kb/ai-reports/_dev` plus `.vscode/decisions` is deleted — no such `--exclude` exists at any of the three extraction sites, and `export-ignore` is what actually holds the line. Its restatement in the load-bearing paragraph is corrected (D-03), and the paragraph now says plainly that both fetch routes honor `export-ignore` while a worktree copy would not, which is why REQ-217's fallback repacks with `git archive`. Two neighbouring sentences that implied other paths *do* carry a `--exclude` fallback were adjusted to match.

**Tests touched:** the ratchet itself is the deliverable.

## Qualification

Passed — 2 files verified, all requirements traced, no debug artifacts.

- Both files are in `write_set`; nothing undeclared was touched.
- Substantive: a 38-line check replacing nothing, and a comment correction that removes a false claim in three places.
- Requirements traced: fetcher enumerated; `skills/*/tools/` clause added with its exemption stated in place; `.gitattributes` claim deleted rather than restored; the `/do-work` and `/kb` lines and the paragraph the contract tests assert are intact.
- Flowing: the clause actually walks files and actually fails — proven by the RED run below, not by inspection.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-canonicalization.sh` (RED, GREEN); `bash _dev/tests/maintainer-verify.sh`; `git archive HEAD | tar t` as an export-ignore spot check

**Result:** ✓ ratchet exit 0; ✓ maintainer-verify exit 0, zero FAIL lines; ✓ the archive still contains zero `do-work/`, `kb/`, `_dev/`, `ai-reports/`, or `CLAUDE.md` entries.

**Red-green validation:** ✗ RED — restoring only the two caller bodies from the pre-REQ-217 commit, with the new clause in place, exits 1 and names both sites at the exact lines the REQ predicted: `do-work-update.sh` line 72 and `install-do-work-suite.sh` line 173. The bootstrap heredoc's `curl` on line 25 of the same installer is correctly **not** flagged, which is the exemption working rather than being absent. → ✓ GREEN — with REQ-217 in place the same run exits 0.

**Existing tests updated:** none.

*Verified by work action*

## Review

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Both halves delivered; two documented departures, each argued rather than assumed |
| Code Quality | 93% | The awk heredoc skip is the only non-obvious part and carries its reasoning above it |
| Test Adequacy | 95% | RED replayed against real history, not a synthetic mutation |
| Scope | 100% | Two declared files |
| Risk | None | A test clause and a comment correction |
| Acceptance | Pass | Fails on the pre-REQ-217 tree, passes on this one, exempts exactly the right line |

**Verdict: Approve** — the gap that let one `curl` become four is now closed by a check rather than by vigilance.

### Findings

**Minor:**
- The clause covers `skills/*/tools/`. The root `tools/` tree — which carries the canonical mirrors — is not walked, because those files are byte-identical to their `skills/` copies and are therefore covered transitively by the mirror contract. That reasoning is real but implicit; if the mirror contract ever loosens, this check silently narrows with it.

**Nit:**
- `find skills/*/tools` relies on the glob matching at least one directory. It does today and would fail loudly if it ever did not, so no guard was added.

### Restatement Sweep

**Triggered** — the `.gitattributes` edit removes a claim about how archives are built. Swept `--exclude`, `export-ignore`, and `belt-and-suspenders` across the repo. Results: `.gitattributes` carried the claim in three places, all corrected here; no shipped file, doc, or test repeated it; `_dev/tests/contract-regressions.sh`'s four assertions are on the export-ignore *lines*, which are untouched. The REQ-217 lesson in `_dev/primes/prime-shell-commands.md` already states the `git archive` requirement correctly.

### Requirements Checklist

- [x] `fetch-upstream-archive.sh` added to the enumerated list — delivered
- [x] No direct download under `skills/*/tools/`, with the exemption stated in the check — delivered (one exemption; the second refused on the merits, D-01)
- [x] Stale `.gitattributes` claim corrected by deletion — delivered
- [x] Load-bearing paragraph and the `/do-work` + `/kb` lines preserved — delivered (five words removed from the paragraph, D-03; the asserted lines byte-identical)

### Acceptance Testing

**Result: Pass**
- Ratchet run against the pre-REQ-217 caller bodies: exit 1, both sites named with line numbers.
- Ratchet run against the current tree: exit 0.
- `bash _dev/tests/maintainer-verify.sh`: exit 0, zero FAIL lines.
- `git archive HEAD | tar t`: zero maintainer-only entries, confirming the comment edits did not disturb the rules they describe.
- Finding-Closure Ratchet: the named regression check is the new `skills/*/tools/` clause, measured failing before and passing after.

### Suggested Additional Testing

- A self-test that mutates a tool script to add a `curl` and asserts the ratchet catches it would make this check replayable the way the REQ-203 and REQ-206 detectors are. Worth doing if a third ratchet clause ever lands; one clause with a real RED behind it does not yet justify the harness.

### Follow-up REQs Created

None.

## Lessons Learned

**What worked:** Reverting the two caller bodies out of git history to produce the RED, rather than hand-writing a mutation. The failure came back naming the exact lines the REQ had predicted, which validated the REQ's diagnosis and the new check in the same run.

**What didn't:** Taking the REQ's `curl -fsSL -o` phrasing literally would have shipped a check that its own tree already evaded — REQ-216 had inserted retry flags between `-fsSL` and `-o` hours earlier. A specified string is a snapshot of the thing to catch, not the thing itself.

**Worth knowing:** A stale claim is rarely in one place. The `--exclude` fallback was asserted in the header, restated in the load-bearing paragraph as part of what makes two lines "the ONLY thing", and implied twice more by sentences explaining why *other* paths lack it. Deleting only the sentence the finding named would have left a file that contradicted itself in the direction of false safety.

## Orientation

A script under any skill's `tools/` directory can no longer hand-roll a download: the canonicalization ratchet now walks that tree and fails on a direct `curl`, exempting only text a script emits for someone else to run. The archive-exclusion header no longer claims a `tar --exclude` safety net that does not exist, so a reader can see that `export-ignore` is the only thing keeping the maintainer's queue and knowledge base out of consumer installs. Lives in the prescribed-shell canonicalization contract and the repo's archive-exclusion rules. No runtime behavior changed — this closes the sweep REQ-216 and REQ-217 opened, so the system's shape is unchanged.

Prime staleness spot-check: `_dev/primes/prime-shell-commands.md` — all referenced paths resolve; its **Closed Enumerations Go Stale** section is what D-02 applied, and needs no correction.
