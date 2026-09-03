---
id: REQ-461
title: '[impact-user-visible] Require affirmative project-owned release targets'
status: claimed
created_at: 2026-09-01T00:12:38Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-413]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-413
claimed_at: 2026-09-03T10:32:09Z
route: B
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/release.go
  - skills/do-work/tools/do-work-cli/internal/publication/release_test.go
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-09-03T10:33:01Z
  basis:
    - Route B
    - 2-file write set
    - 1 subsystem involved
    - 5 acceptance criteria
    - cross-route regression gates
---

# Require Affirmative Project-Owned Release Targets

## What

Replace convention-based installed/generated path exclusions with condition-complete evidence that every release target is a project-owned source or declared maintainer mirror.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read both `prime_files`, `prime-do-work-cli.md` and `lessons-do-work-cli.md` in full (including this family's REQ-460 entry and REQ-528's necessary-is-not-sufficient entry), and the crew members. Approach: prove ownership from two declarations the repository already owns — Git's index and the `SKILL.md` package marker — with no directory name in the code.
- [x] **[APPLY]:** Two files, both inside the declared Scope.
- [x] **[UNIFY]:** Orchestrator independently re-ran `go build ./...`, `go vet ./...` (clean), `gofmt -l .` (silent) and `go test -count=1 ./internal/publication/` (`ok … 20.1s`), and read the whole new ownership block. Grepped the file for every directory name the old predicate carried: the only surviving mentions are two words inside a doc comment explaining what each half catches. No name is tested anywhere in the code.

## Finding Provenance

REQ-413's fresh re-review found that release exclusion recognizes named conventions such as `vendor`, `node_modules`, `.codex/skills`, and literal generated directories, but accepts other dependency or generated locations such as `third_party/do-work`, `dist/skills/do-work`, and cache-owned skill trees.

## Detailed Requirements

- Require affirmative, verifiable project-owned classification for every consumer release target instead of inferring safety from the absence of known bad directory names.
- Permit declared maintainer mirrors only through the existing explicit mirror contract and retain byte-identity validation for changelog mirrors.
- Refuse installed skills, dependencies, vendored packages, caches, distribution outputs, and generated trees regardless of their directory spelling.
- Add fixtures for non-example paths including `third_party/do-work`, `dist/skills/do-work`, and an arbitrarily named cache tree.
- Keep refusal actionable by identifying the target and the ownership evidence that is missing or inconsistent.

## Constraints

- Do not replace one directory-name denylist with a larger denylist.
- Preserve caller-selected semantic version and changelog-content judgment.

## Dependencies

No request prerequisite.

## Builder Guidance

Certainty level: Firm on the invariant, flexible on the proof mechanism. Prefer an explicit manifest or repository-owned declaration that the planner can verify mechanically.

## Red-Green Proof

**RED prompt/case:** Supply otherwise valid release targets beneath `third_party`, `dist`, and an unrecognized cache/generated subtree.
**Why RED now:** The current exclusion is a finite convention list, so targets outside those spellings are accepted.
**GREEN when:** Undeclared or non-project-owned targets refuse independent of path spelling, while verified project sources and declared maintainer mirrors apply atomically.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-413-rereview.md` for source context and independent evidence.

---
*Source: REQ-413 fresh re-review finding F7.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The invariant is firm and the single function is obvious, but the REQ deliberately leaves the proof mechanism open ("flexible on the proof mechanism"), and what "project-owned" means *mechanically* is not stated anywhere — that had to be established from what the repository can already verify.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

`installedReleasePath` (`internal/publication/release.go:170-189`) is a denylist of directory names: `vendor`, `vendored`, `node_modules`, a `.codex`/`.claude` followed by `skills`, a `generated`/`.generated` with a `skills` or `do-work` descendant, and the literal prefix `skills/do-work/`. Its two callers (`:34`, `:72`) use it to stop a **consumer** release from mutating installed or vendored suite metadata — a consumer releasing their own project must not bump the version of the do-work copy installed inside their repo.

Everything the REQ names slips through, and I confirmed each is only a spelling away from a listed one: `third_party/do-work`, `dist/skills/do-work`, and any cache tree whose directory happens not to be called `generated`. This is the same shape as REQ-460 — a contract that states a condition, implemented as membership in a list of examples — and the same family, `closed-enumeration-for-a-condition`.

Two mechanisms the repository can already verify, either of which is affirmative rather than exclusionary:

1. **Package markers.** An installed suite package carries its own `SKILL.md` at the top of its directory. The board tool's repo-root walk already relies on exactly this to skip directories "merely *named* `do-work` that are skill installs (SKILL.md at their top level)" — so a target with a package marker at or above it is an installed copy regardless of what the enclosing directory is called.
2. **Git tracking.** A project-owned source is tracked in this repository; a distribution output or cache tree generally is not. This catches `dist/` and cache trees by what they *are* rather than by what they are named.

Neither alone is complete — a repository that commits its `vendor/` tree is tracked but not project-owned, and a package marker does not exist for a generated non-skill artifact — which is why the affirmative test likely needs both, and why `maintainer_release: true` must remain the only door to the suite's own metadata.

The constraint is the sharp edge here: *do not replace one directory-name denylist with a larger denylist*. A fix that adds `third_party`, `dist` and `.cache` to the list satisfies the fixtures and fails the REQ.

*Generated in-session (single-pass discovery)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/publication/release.go` (modify) — replace the convention denylist with an affirmative project-owned test, and make the refusal name the missing evidence
- `skills/do-work/tools/do-work-cli/internal/publication/release_test.go` (modify) — fixtures for the three non-example paths the REQ names plus an arbitrarily-named cache tree, and controls for a genuine project source and a declared maintainer mirror

**Files I will NOT touch:** the semantic-version comparison and changelog-content judgment (the REQ preserves caller-selected judgment), and the changelog byte-identity mirror validation.

**Acceptance criteria (restated from REQ):**
- [x] Every consumer release target requires affirmative, verifiable project-owned classification rather than the absence of known-bad directory names
- [x] Declared maintainer mirrors are permitted only through the existing explicit mirror contract, with changelog byte-identity validation retained
- [x] Installed skills, dependencies, vendored packages, caches, distribution outputs and generated trees refuse regardless of directory spelling — **the created-target half was not delivered by the first implementation** (review F1/F2: an unattested location was attested by the repository root, and the only negative test passed by fixture accident); delivered by the remediation (D-05, D-06)
- [x] Fixtures cover `third_party/do-work`, `dist/skills/do-work`, and an arbitrarily named cache tree
- [x] Refusals name the target and the ownership evidence that is missing or inconsistent
- [x] No directory-name denylist, larger or otherwise, replaces the old one

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/publication/release.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/release_test.go` (modified)

**What was done:** `installedReleasePath` is deleted. `releaseTargetOwnershipGap` replaces it and returns the missing or contradicted evidence rather than a boolean, so the refusal can name it. Ownership rests on two declarations the repository itself owns, both required: Git's index (a target that must already exist has to be tracked; one the release creates has to be un-ignored and sit in the nearest existing ancestor that holds tracked sources) and the absence of a `SKILL.md` between the repository root and the target's parent — the same discriminator `repositorymodel.FindRepositoryRoot` and the board's repo-root walk use to tell an install from a directory merely named `do-work`. The refusal code becomes `RELEASE-TARGET-OWNERSHIP-UNVERIFIED`. `maintainer_release: true` remains the only exemption and is checked before any ownership work runs. No directory name is tested anywhere in the new code.

## Decisions

- **D-01** — Proof mechanism: Git's index **plus** the `SKILL.md` package marker, both required. **ESCALATE.** **Value:** both are declarations the repository owns, so the planner verifies rather than trusting the manifest, and no directory name appears in the code — which is the REQ's constraint, not merely its requirement. **Risk:** the release planner now shells out to git and fails closed when git cannot answer, so a non-git root that plans today refuses tomorrow. Reversible in one function.
- **D-02** — Reused the package-marker *idea*, not code. **DECIDE & STATE.** The CLI prime's **Package direction** rule forbids importing `queue-kanban`, and `repositorymodel.FindRepositoryRoot`'s marker test is inlined rather than exported, so nothing concrete crosses the boundary. The test is five lines and the doc comment cites both precedents. Inside the package the existing `gitPathTracked` helper (`defer_gate.go:521`) was reused rather than duplicated.
- **D-03** — Refusal code renamed `RELEASE-INSTALLED-METADATA-REFUSED` → `RELEASE-TARGET-OWNERSHIP-UNVERIFIED`. **DECIDE & STATE.** The old name stopped describing the condition once it began firing for an untracked cache tree. Grepped repo-wide first: the string existed only in `release.go`, `release_test.go` and a stale compiled binary — no action file, ADR or doc keys on it. **Value:** the code names the real condition, so the next hardening cannot read it narrowly, which is exactly REQ-460's failure mode. **Risk:** a consumer runbook grepping the old code; trivially reversible.
- **D-04** — Residual gaps, named rather than assumed away (REQ-528's lesson applied). **ESCALATE.**
  - **R1** A committed vendored copy with *no* `SKILL.md` is tracked and unmarked, so it is accepted. For the class this guard exists to protect — installed do-work suite metadata — the gap is closed, because `suitemanifest.ValidateSuite` requires a non-empty `SKILL.md` for every module. For a generic third-party library a consumer committed, it is open. No name-free way to close it was found: every "package manifest" test (`package.json`, `Cargo.toml`, `go.mod`) is a denylist wearing a different hat, which the REQ's constraint forbids.
  - **R2** The mirror image: a consumer who **authors** a skill in their own repository and wants to version `skills/my-skill/VERSION` in place is now refused by their own marker, with `maintainer_release` as the only door. The old predicate accepted it, because it denylisted only the literal `skills/do-work/`. The filesystem cannot distinguish a skill you wrote from one you installed — both are tracked and both carry the marker. **This is a behavior regression beyond the REQ's ask** and is routed to its own REQ rather than absorbed here.
  - **R3** A repository with nothing tracked at all (fresh `git init`, no commit) cannot attest a bootstrap changelog location, so its first release refuses. Fail-closed with an actionable message; no carve-out added, because a carve-out here would be speculative defense.

## Discovered Tasks

- **The same defect one package over, and narrower:** `internal/finalization/finalization_discovery.go:1045` `installedReleasePathForDiscovery` carries a prefix denylist (`.claude/skills/`, `.codex/skills/`, `node_modules/`, `vendor/`, `vendored/`, `generated/`, `.generated/`), called at `:725` while enumerating configured release members. It admits the same three paths this REQ fixed, and being prefix-anchored rather than per-segment it also loses `packages/vendor/do-work/VERSION`, which the old release predicate caught. Half the fix is already present — the enumeration only walks tracked paths — so the missing half is exactly the package marker. Its `suite/modules.tsv`-at-HEAD test for "is this the maintainer repo" is itself an affirmative repository-owned declaration, which is supporting precedent for D-01.
- `finalization_discovery.go:624` `releaseMetadataPath` is a closed list of metadata filenames. It is a *positive* enumeration, so a missing name means "not discovered" rather than "wrongly permitted" — stale rather than unsafe. Worth re-keying on a condition, or marking illustrative-and-extensible.
- Checked and cleared as not this shape: `repositorymodel/repository_model.go:186` (queue structure), and `cleanup/cleanup_plan.go:344` plus `repositorymodel/repository_model.go:128`, which both already use the `SKILL.md` marker affirmatively.

## Testing

**Tests run:** `go build ./... && go vet ./... && gofmt -l .`; `go test -count=1 ./internal/publication/`; `go test ./...`; canonical repository gate `bash _dev/tests/maintainer-verify.sh`.
**Result:** ✓ `internal/publication` green (`ok … 20.1s`). Gate exits 1 on two failures, neither this diff's — see attribution below.

**Red-green validation:** traces the REQ's Red-Green Proof, which asked for otherwise-valid targets beneath `third_party`, `dist`, and an unrecognized cache subtree.
- `installed package under an unlisted dependency directory` (`third_party/do-work/VERSION`): ✗ `target without ownership evidence accepted: … Refusal:(*publication.Refusal)(nil)` → ✓
- `distribution output tree` (`dist/skills/do-work/VERSION`): ✗ accepted with a full plan → ✓
- `arbitrarily named cache tree` (`.stash-pantry/packages/do-work/VERSION` — a name appearing in no list anywhere): ✗ accepted → ✓
- `TestBuildReleasePlanRefusesBootstrapChangelogInAnUnattestedLocation` (`dist/skills/do-work/CHANGELOG.md`, created rather than replaced): ✗ accepted → ✓

**Three revert differentials, not one.** The four legacy denylist rows (`vendor/`, `node_modules/`, `.claude/skills/`, `skills/do-work/`) go red under a naive revert on the **refusal-code rename alone**, which proves nothing about the new condition. So each half was isolated separately:

| Experiment | Rows that redden |
|---|---|
| old predicate, new refusal code — isolates the *condition* | third_party, dist, cache tree, bootstrap-create |
| new code, `enclosingPackageMarker` disabled | third_party, vendor, `.claude/skills`, `skills/do-work` |
| new code, tracked check disabled | dist, cache tree, node_modules |

Every refusing row is a lock-in for one specific half, and no row is decorative. Each is named by the class it represents rather than by its path, per the REQ's requirement that examples are fixtures.

**Refusal text, captured from the built code** — target and missing evidence both named:
- `consumer release target third_party/do-work/VERSION is not proven project-owned: third_party/do-work/SKILL.md declares an installed package that owns this subtree; only maintainer_release may mutate a suite package's own metadata`
- `… dist/skills/do-work/VERSION … : the repository does not track it, so Git does not attest it as a project source; …`
- `… dist/skills/do-work/CHANGELOG.md … : the repository tracks no source in dist/skills/do-work, so the new target's location is unattested; …`

**Maintainer-path regression proof.** `TestBuildReleasePlanPlansThisRepositorysOwnMaintainerReleaseShape` builds this repository's exact release shape — `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, with `skills/do-work/SKILL.md` present, all six committed, `RequiredMirrors` set to the three version targets and `maintainer_release: true`. It plans 5 mutations with no refusal. Three of those targets sit under the very marker that refuses a consumer release, so this is the highest-risk regression in the change and it is now pinned rather than assumed.

**Canonical repository gate — attribution.** Exits 1 on two failures, neither attributable to this diff, which touches only `internal/publication`:
- `internal/toolboxcommands` → `TestRemediationCancellationReachesMediaGitCommitAndRollback`: in the recorded red baseline since `008f3d3`, identical text. Tracked as **REQ-524**.
- `internal/suiteinstall` → `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation/update-suite/INT`: **fourth sighting** of the load-dependent signal-timing flake, on a third distinct subtest. Confirmed `-count=3` green in isolation this run. Tracked as **REQ-525**; it has now disrupted roughly half this session's gate runs, which is itself an argument for pulling that REQ forward.

*Verified by work action*

## Remediation

The review scored 73% with Acceptance Partial, which mandates a follow-up REQ per Important finding. Three of the four were fixed here instead, because all three are this REQ's own requirements rather than adjacent work: F1 and F2 *are* requirement 3 on the created-target path, and F3 *is* requirement 5's actionability half. The fourth — the `work-reference.md:1020` prose disagreement — was folded into REQ-532, which already owns the discovery half of that same sentence.

**Remediation commit:** see `commit:` below. Supersedes the doc comment and refusal strings quoted in the first `## Testing` block above.

- **D-05** — Created-target attestation anchors on the target's **immediate parent** (it must exist and the index must claim a source strictly inside it); the ancestor walk is deleted. **ESCALATE.** The alternative — keep walking, but forbid the root attesting a deeper target — was checked and agrees on every reviewer probe, differing only where an intermediate directory exists but the immediate parent does not (`a/` tracked, `a/b/` absent, target `a/b/CHANGELOG.md`), which the walking version accepts on a location nothing has ever declared. **Value:** an attestation is a statement about the directory the file lands in, which is the only directory the index can speak to; it is also one `filepath.Dir` instead of a loop. **Risk:** a consumer release can no longer create a *new* directory at all — `docs/CHANGELOG.md` into a not-yet-existing `docs/` now refuses. The only real bootstrap shape, a root `CHANGELOG.md`, still plans, and the refusal names the fix ("commit a project source in `docs/` first"). The maintainer path is untouched because it short-circuits. Reversible in one function.
- **D-06** — F2 fixed by filtering the index listing to entries strictly *inside* the directory (`ls-files -z`, dropping the entry equal to the directory itself). **DECIDE & STATE.** A gitlink or tracked symlink is one index entry *at* the path, not a source *in* it. `-z` also removes git's path quoting, which would otherwise have failed the prefix test on unusual filenames — fail-closed, but a false refusal.
- **D-07** — `releaseTargetOwnershipGap` returns a typed `releaseOwnershipGap{MissingEvidence, Remedy}` instead of a bare string, each gap carrying its own remedy; the zero value means proven, so it cannot be confused with an empty message. **DECIDE & STATE.** **This makes F3 structurally impossible to regress**: a new gap that forgets its remedy produces a visibly truncated message rather than silently inheriting the wrong one. The free `ownershipRefusalReason` is now `gap.refusalReason(path)`.
- **D-08** — `enclosingPackageMarker` no longer requires the marker to be a non-directory. **DECIDE & STATE.** Dropping `!info.IsDir()` is *less* code and fails closed: a `SKILL.md` directory now still counts as a marker. A legitimate install cannot have one — `suitemanifest.ValidateSuite` and `suiteinstall/update_transaction.go:188` both require a non-empty regular file — so no real install changes behavior.
- **D-09** — `os.Stat`'s case-insensitivity on macOS/Windows is accepted and recorded in the doc comment rather than fixed. **DECIDE & STATE.** The failure mode is a stray ancestor `skill.md` refusing a release that could have been allowed, which is the safe direction; making it case-exact means `os.ReadDir` on every ancestor of every target. No directory-name test was added in any form, per the REQ's constraint.
- **D-10** — Deleted `nearestExistingAncestorDirectory`, orphaned by D-05, rather than leaving it unused.

**The fixture accident, from both sides.** Before the fix, deleting the single `writeFixture(t, root, "dist/skills/do-work/placeholder", …)` line from `TestBuildReleasePlanRefusesBootstrapChangelogInAnUnattestedLocation` turned it into an accepted plan creating `dist`, `dist/skills`, `dist/skills/do-work`. After the fix, the same deletion leaves it **still refusing** — the guard now carries the test rather than the fixture. And the neuter table shows it from the other side: restoring the ancestor walk reddens the three new create rows but **not** that bootstrap test, which is exactly why it never caught this.

**Neuter-and-confirm, five guards, one at a time, source byte-compared after each restore:** the parent anchor reddens only the three create rows plus the unattested-location remedy row; counting the directory's own index entry reddens only the single-object row; giving every gap the maintainer remedy reddens only the three remedy rows; skipping the marker half for created targets reddens only the marker-above-an-attested-location row; and treating a git error as attestation reddens only the two fail-closed rows.

**Refusal text now matches each gap** — the F3 fix, captured from the built code:
- uncommitted existing target → `commit the target so Git attests it as a project source, then release again`
- ignored created target → `drop the ignore rule covering the target, or release a target the repository claims as a source`
- unattested created location → `commit a project source in dist/skills/do-work first, or create the target where the repository already tracks sources`
- package marker → `only maintainer_release may mutate a suite package's own metadata` (the one case where that remedy is correct)

**Verification:** `go build`/`go vet`/`gofmt -l` clean; `go test -count=1 ./internal/publication/ ./internal/finalization/` both green (the latter is the other `BuildReleasePlan` consumer); `go test ./...` red only on REQ-524's known baseline failure. `maintainer_release` still short-circuits ahead of all ownership work, re-confirmed by panicking `releaseTargetOwnershipGap` and watching both maintainer tests pass. Repo-grep for directory-name literals in `release.go`: **zero**.
