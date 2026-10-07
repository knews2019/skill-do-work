---
id: REQ-638
title: '[impact-negligible] Release guard reads suite/modules.tsv once instead of merging two reads of the same declaration'
status: claimed
created_at: 2026-10-06T22:45:18Z
user_request: UR-137
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
related: [REQ-635, REQ-636, REQ-637]
batch: review-0305-69-findings
required_lessons: [_dev/primes/lessons-releases.md]
write_set: [skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go, skills/do-work/CHANGELOG.md]
claimed_at: 2026-10-07T15:01:06Z
---
# Release Guard Reads suite/modules.tsv Once Instead of Merging Two Reads of the Same Declaration

## What
Review finding F8 (cosmetic). `releaseShippedChangeError` (`skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:36-55`) calls `DeclaredMaintainerReleaseRoots`, then `DeclaredModuleSources`, then merges the two through a set. `DeclaredMaintainerReleaseRoots` (`releaseownership/release_ownership.go:127-139`) already calls `DeclaredModuleSources` and filters it, so roots are a subset of sources and the merge is redundant; `headReleaseImage` is also built twice. Simplify to one read with identical behaviour.

## Detailed Requirements
1. Call `DeclaredModuleSources` once. If no declared source has a tracked `<source>/VERSION`, return nil (the repository is a consumer project and is not guarded), exactly as today's `len(roots) == 0` branch does.
2. Otherwise the shipped roots are the declared sources plus `suite` and `tools`, sorted. The output must be byte-identical to today's for every input.
3. Build `headReleaseImage(repositoryRoot)` once.
4. Keep the existing tests green and run the full Go test suite for the do-work CLI module. No new test is required: the change is a deletion with identical behaviour and no runnable RED exists.
5. Release per `_dev/primes/prime-releases.md` with a CHANGELOG entry (shipped Go changes).

## Constraints
- Behaviour-preserving only; no change to `DeclaredMaintainerReleaseRoots` or `DeclaredModuleSources` themselves, which have other callers.
- Batch constraint: each REQ in this batch is its own release and its own commit.

## Dependencies
None. Independent of REQ-635, REQ-636 and REQ-637.

## Builder Guidance
Certainty is high. This is a deletion. Builder latitude: whether the tracked-`VERSION` check stays a small loop inline or becomes a one-line helper.

## Red-Green Proof
**RED prompt/case:** Not applicable as a failing test: the change preserves behaviour. Proof is the existing `finalization` test suite green before and after, plus a diff that removes the second `DeclaredModuleSources` read, the second `headReleaseImage` call and the set merge.
**Why RED now:** The function parses `suite/modules.tsv` twice and merges a subset into its superset.
**GREEN when:** One parse, one image build, existing tests green, `go vet` clean.
**Validation:** Inferred during capture. Confirmed by reading both functions during triage.

## Full Context
See `do-work/user-requests/UR-137/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: review finding F8 from the pasted consumer code review of 0.305.58 to 0.305.69, triaged with `do-work validate-feedback` on 2026-10-06 and accepted.*
