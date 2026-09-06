# REQ-557 Exploration — Deduplicate six Go helper names across do-work-cli

**Request:** `/home/user/skill-do-work/do-work/queue/REQ-557-deduplicate-six-go-helper-names-defined-fourteen-times-across-do-work-cli.md`
**HEAD:** `fce57fc` — `[work run] REQ-592 pre-flight: green gate at 15e2ec3, baseline recorded, builder dispatched`
**Working tree:** clean (`git status --porcelain` empty)
**Prior exploration:** none. `do-work/runs/work-2026-09-05-170806/` holds explorations for REQ-486, REQ-552, REQ-554, REQ-556, REQ-583, REQ-587 and REQ-591 only. `find do-work/runs -name '*557*'` returns nothing. This is the first exploration of REQ-557.

**Route recommendation: C (plan-first).**

---

## 1. Verdict up front

The request cannot go straight to a builder. Four of its factual claims are wrong at HEAD, two files it must edit are missing from its `write_set`, and its central constraint contradicts the repository as it actually stands. None of that makes the request wrong in spirit. The duplication is real and the target of one definition per name is right. But the plan it hands a builder would not compile the intent it describes, and one of its "duplicates" is not a duplicate at all.

Five things changed or were never true:

| Code | Claim in the REQ | State at HEAD |
|---|---|---|
| F1 | Fourteen definitions | **Fifteen.** A 15th copy landed after the audit |
| F2 | Every duplicator already imports the canonical home | **False for four of six packages** |
| F3 | `finalization`'s `uniqueSorted` duplicates `corehelpers` | **Different contract.** Normalizes paths, returns nil on empty or absolute |
| F4 | Three semantic reconciliations | **Four.** `requestIDLess` diverges too |
| F5 | `release.go` "rejects" unparseable semver | **It returns 0**, same as the other copy. The split is strict vs. lenient parsing |

---

## 2. The Reproduce command at HEAD

The REQ's own command, run from the repository root:

```
rg -n --glob '*.go' --glob '!*_test.go' '^func (uniqueSorted|subtractPaths|requestIDLess|firstError|compareSemver|physicalPath)\(' skills/do-work/tools/do-work-cli/internal/ | sort
```

Prints **15 lines**, not 14:

| # | File:line | Name |
|---|---|---|
| 1 | `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go:827` | `subtractPaths` |
| 2 | `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go:846` | `uniqueSorted` |
| 3 | `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go:871` | `firstError` |
| 4 | `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go:307` | `requestIDLess` |
| 5 | `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go:365` | `subtractPaths` |
| 6 | `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go:379` | `uniqueSorted` |
| 7 | `skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go:147` | `physicalPath` |
| 8 | `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go:1337` | `uniqueSorted` |
| 9 | `skills/do-work/tools/do-work-cli/internal/knowledgecommands/interview_commands.go:1535` | `compareSemver` |
| 10 | `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go:83` | `requestIDLess` |
| 11 | `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go:361` | `firstError` |
| 12 | `skills/do-work/tools/do-work-cli/internal/publication/release.go:210` | `compareSemver` |
| 13 | **`skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green.go:442`** | **`uniqueSorted`** |
| 14 | `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go:654` | `requestIDLess` |
| 15 | `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go:209` | `physicalPath` |

Row 13 is the new one. `git log -S 'func uniqueSorted(values []string) []string' -- .../repairvalidation/already_green.go` names commit `7bd3464`, `[work run] work-2026-09-05-124800: REQ-585 qualified, tested, reviewed 94%, held for heavy lanes` — two days after the audit.

Re-running the same pattern without `--glob '!*_test.go'` also returns 15, so no test file defines any of the six names. The exclusion is decorative here, which is worth knowing when writing the lock-in.

**Two consequences.** The REQ's Red-Green Proof line "It prints fourteen definitions" is stale and must be corrected to fifteen before handback. And the audit's proposed lock-in at `do-work/audits/audit-2026-09-03.md:334` — `-le 14`, "currently 14; red at 15" — is **already red at HEAD**, which directly contradicts the REQ's constraint that the assertion be "pinned at today's value so it is green on day one".

### 2.1 Files touched vs. files declared

The REQ's `write_set` (line 18) names nine Go files. The fifteen definitions live in **ten** Go files. Missing:

- `skills/do-work/tools/do-work-cli/internal/publication/release.go` — holds `compareSemver` (row 12)
- `skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green.go` — holds `uniqueSorted` (row 13)

Neither is optional. `compareSemver` has exactly two definitions, one in `interview_commands.go` and one in `release.go`. Deleting only the first leaves one definition, which is correct, but the REQ's Detailed Requirements describe `release.go`'s copy as the one that "rejects" bad input, implying it is the survivor. That works, and `release.go` still needs an edit if the canonical helper moves out of `publication`. `already_green.go` is not optional at all: leave it and the count lands at 7, not 6.

---

## 3. The import claim is the real blocker

The REQ says, and the audit says: *"in every case the duplicating package already imports the package holding an earlier copy, so no import cycle forced the copy"* (`do-work/audits/audit-2026-09-03.md:307`). Its Constraints repeat it as a hard rule: *"No import cycles introduced: the canonical home is always a package the duplicator already imports"* (REQ line 50).

Checked per package with `rg -n 'internal/(corehelpers|repositorymodel)"'`, and confirmed with a full-text `rg -n 'corehelpers'` over each package directory:

| Package | Imports `corehelpers`? | Imports `repositorymodel`? | Duplicates it holds |
|---|---|---|---|
| `finalization` | **yes** — `finalization_apply.go:15`, `finalization_discovery.go:15` | yes — `finalization_prepare.go:16` | `subtractPaths`, `uniqueSorted` |
| `dependencygraph` | no | **yes** — `dependency_graph.go:11` | `requestIDLess` |
| `nextselection` | no | **yes** — `next_types.go:10` | `requestIDLess` |
| `knowledgecommands` | **no (zero hits)** | **no** | `uniqueSorted`, `compareSemver`, `physicalPath` |
| `suiteinstall` | **no (zero hits)** | **no** | `physicalPath` |
| `publication` | **no (zero hits)** | yes — `defer_gate.go:16`, `answer.go:15` | `firstError`, `compareSemver` |
| `repairvalidation` | **no (zero hits)** | yes — `already_green.go:21` | `uniqueSorted` |

The claim holds only for the three `requestIDLess` copies and for `finalization`. For the path and error helpers the REQ wants in `corehelpers`, **four packages would need a brand-new import edge**. The premise that "no import cycle forced the copy" may still be true historically, but the constraint as written is unsatisfiable today.

### 3.1 No cycle, but bad layering

`go list -deps ./internal/corehelpers` (internal packages only):

```
resultmodel, commandruntime, atomicfile, ownedprocess, gittransaction,
schemanormalization, requestmodel, repositorymodel, dependencygraph,
requeststate, cleanup, doctor, nextselection, corehelpers
```

None of `knowledgecommands`, `publication`, `suiteinstall` or `repairvalidation` appears, and those four are imported only by `cmd/do-work-cli/main.go` (plus `publication` by `finalization` and `lifecycleadvance`). **The four new edges compile.** The problem is not cycles, it is direction: `corehelpers` is the command-handler registry (`internal/corehelpers/commands.go:49  func Handlers() map[string]commandruntime.CommandHandler`), and importing it drags thirteen internal packages behind a six-line helper.

One edge is genuinely impossible: `corehelpers` imports `nextselection` (`commands.go:19`), so `nextselection` can never import `corehelpers`. This does not bite, because `nextselection`'s duplicate is `requestIDLess`, whose home is `repositorymodel` — and `repositorymodel` imports only `atomicfile` and `requestmodel`, so it stays a low-level package.

### 3.2 Three options for the canonical home

- **O1 — `corehelpers`, as the REQ says.** Add four import edges. Compiles, no cycle. Cost: `knowledgecommands`, `publication`, `suiteinstall` and `repairvalidation` inherit the whole command-registry dependency graph. Contradicts the REQ's own constraint.
- **O2 — a new stdlib-only leaf package (recommended).** `internal/sharedhelpers/shared_helpers.go` holding `UniqueSorted`, `SubtractPaths`, `FirstError`, `PhysicalPath`, `CompareSemver`. Zero internal imports means zero cycle risk from any direction, forever. Cost: a new package the `write_set` does not authorise, and a junk-drawer risk that a package doc comment stating the admission rule (stdlib only, no do-work domain types, no internal imports) should hold off. `internal/` has no existing generic utility package: of the ten leaf packages (`archivefetch`, `atomicfile`, `managedsection`, `ownedprocess`, `releaseownership`, `resultmodel`, `schemanormalization`, `settingshooks`, `suitemanifest`, plus `commandruntime` and `requestmodel` at two deps) every one is a named domain, so there is nowhere existing that these five belong.
- **O3 — reduce scope.** Deduplicate only where the import already exists: `subtractPaths` into `corehelpers`, `requestIDLess` into `repositorymodel`. Gets 15 down to 11, not 6. Lock-in pinned at 11. Leaves the request half done and its stated target unmet.

`repositorymodel` is the right home for `requestIDLess` under every option. Naming for Reach (`skills/do-work/crew-members/coding-guardrails.md` § 5) is satisfied by all the proposed exported names: `UniqueSorted`, `SubtractPaths`, `RequestIDLess`, `FirstError`, `CompareSemver`, `PhysicalPath` are each two words and each findable by plain-text search. `sharedhelpers` follows the existing `corehelpers` compound-of-two-words convention.

---

## 4. What each duplicate actually is

### 4.1 `subtractPaths` — a true duplicate

`corehelpers/checks.go:827` uses the package's `stringSet` helper (`checks.go:859`); `finalization/finalization_prepare.go:365` builds the same map inline. Same inputs, same outputs, same empty-slice-not-nil return. Clean delete.

Call sites: `checks.go:177`, `checks.go:178`, `finalization_prepare.go:162`.

### 4.2 `uniqueSorted` — four copies, three contracts

| Site | Behaviour |
|---|---|
| `corehelpers/checks.go:846` | Dedupe and sort, **verbatim**. Keeps `""`. No path normalization |
| `knowledgecommands/interview_commands.go:1337` | Dedupe and sort, **drops `""`** |
| `repairvalidation/already_green.go:442` | Dedupe and sort, **drops `""`**. Same contract as the row above |
| `finalization/finalization_prepare.go:379` | **Neither.** `result, _ := normalizeRepositoryPaths(paths); return result` |

That last one is the trap. `normalizeRepositoryPaths` at `finalization_prepare.go:345-363`:

```go
if path == "" || filepath.IsAbs(path) {
    return nil, fmt.Errorf("path must be non-empty and repository-relative: %s", path)
}
cleaned := filepath.Clean(path)
if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
    return nil, fmt.Errorf("path escapes repository: %s", path)
}
set[filepath.ToSlash(cleaned)] = true
```

The error is discarded by `uniqueSorted`, so a single empty or absolute path anywhere in the input makes the whole call return `nil`. And `a//b` and `a/b` collapse to one entry, which no other copy does. This is a **different function that shares a name**, not a duplicate.

`uniqueSorted` has 24 call sites overall. Eighteen are in `finalization/finalization_discovery.go` (lines 67, 70, 132, 162, 188, 194, 195, 205, 573, 627, 640, 649, 844, 895, 1422, 1432, 1470, 1518), three in `finalization/finalization_apply.go` (233, 518, 522), one in `finalization_prepare.go:162`, plus `finalization_recovery_test.go:574`. Several feed refusal codes such as `FINALIZATION-DISCOVERY-PROTECTED-STAGED` and `FINALIZATION-MULTIPLE-TAILS`. In practice the inputs come from `git status` and are relative, so a swap would very likely leave the tests green while silently widening the contract. That is the worst kind of change: invisible until it is not.

### 4.3 `requestIDLess` — three copies, and the divergence the REQ missed

| Site | Parser | Behaviour on `"REQ-12abc"` |
|---|---|---|
| `dependencygraph/dependency_graph.go:307` | `requestNumber` (`dependency_graph.go:316-330`): manual prefix check, digit scan stops at first non-digit | parses **12** |
| `repositorymodel/repository_model.go:654` | `requestNumberFromText` (`:624-631`) on `requestNumberPattern = regexp.MustCompile("(?i)^REQ-0*([0-9]+)")` (`:24`) | parses **12** |
| `nextselection/next_types.go:83` | `numericID` (`:60-73`): **every** character after the prefix must be a digit | **rejects**, falls to string comparison |

The REQ flags the first pair as "two different number parsers". Different implementations, yes. But I traced them and they agree on every input I could construct, including leading zeros, leading whitespace, and integer-overflow digits (both return `false`). The pair is behaviourally equivalent.

The genuine divergence is the third copy, which the REQ does not mention at all. Under the REQ's plan `nextselection` adopts the lenient parser, which changes sort order for malformed identifiers at `next_targets.go:97` and `next_targets.go:151`. That is a fourth reconciliation, and the REQ says plainly at line 44 that a silent pick is a review refusal.

Deleting `dependencygraph`'s copy also orphans `requestNumber`: `dependency_graph.go:308` and `:309` are its only callers. Go will not warn, `vet` will not warn, and the REQ-550 delegate lock-in will not catch it because it is unexported. It must be deleted deliberately. `nextselection`'s `numericID` stays — eleven other call sites in `next_targets.go` still use it.

### 4.4 `firstError` — byte-identical, the easy one

`publication/capture_files.go:361-366` and `corehelpers/checks.go:871-876` are the same six lines with the same parameter names. Call sites: `checks.go:153`, `capture_files.go:151`, `publication/release.go:44`, `publication/answer.go:307`.

### 4.5 `compareSemver` — the REQ describes it wrong

The REQ (line 42): *"returns 0 for unparseable input while `internal/publication/release.go` rejects it."*

`release.go:210-224` does not reject:

```go
oldParts, oldOK := parseSemver(oldVersion)
newParts, newOK := parseSemver(newVersion)
if !oldOK || !newOK {
    return 0
}
```

Both copies return 0. The split is **strict vs. lenient parsing**:

- `release.go` via `parseSemver` (`:227-244`): exactly three parts, no empty part, no leading zero on a multi-digit part, no negative. Anything else scores 0 and the release is treated as no-change.
- `interview_commands.go:1535-1547`: splits on `.`, discards every `Atoi` error so a non-numeric component silently becomes 0, and compares only `min(3, len(a), len(b))` components.

They diverge on partially malformed input. `compareSemver("1.0.0", "1.2.x")` is `-1` in the interview copy and `0` in the release copy.

Good news for the swap: `interview_commands.go:811-815` already calls `semverMajor` on both versions and returns an error unless each is three dot-separated parts with an integer major, *before* reaching the comparison at line 820. So fully unparseable input never gets there. Adopting the strict contract is low-risk and matches the release path.

`parseSemver` stays in `publication`: `release_mirrors.go:109` still calls it. The canonical `CompareSemver` should carry its own copy of the strict parse rather than dragging `parseSemver` out of `publication`, which would widen scope into a file the REQ never names.

### 4.6 `physicalPath` — two real contracts

| Site | Contract |
|---|---|
| `suiteinstall/update_transaction.go:209-215` | `filepath.EvalSymlinks` then `filepath.Abs`. **Errors if the path does not exist** |
| `knowledgecommands/commands.go:147-172` | `Lstat`-walks up to the first existing ancestor, `EvalSymlinks` that, rejoins the missing tail, returns `filepath.Clean`. **No `Abs`** |

The knowledgecommands version is a strict superset for existing paths, except that it never makes a relative result absolute. So the merged contract is: walk missing ancestors, then `Abs`.

The behaviour change lands at `update_transaction.go:179`. Today a nonexistent `InstalledSkillRoot` fails there with `"the installed do-work skill could not be located"`. With ancestor-walking it fails later at line 188 with `"SKILL.md is missing at ..."`. The existence guard is not lost — lines 171-173 already `os.Stat` `ProjectRoot` separately, and lines 188 and 191 check `SKILL.md` and `actions/version.md`. `rg` over every `*_test.go` in the CLI finds no test asserting either message, so nothing turns red. That is exactly why it needs to be a written decision.

One test names the helper directly: `knowledgecommands/memory_commands_test.go:390`.

---

## 5. The four decisions that need recording

The REQ demands three named reconciliations in the Implementation Summary. There are four.

**D1 — `uniqueSorted` empty-string handling, and the `finalization` normalizer.** Two sub-parts. First: does the canonical helper keep `""` (the `corehelpers` behaviour, 6 call sites) or drop it (the `knowledgecommands` and `repairvalidation` behaviour, 12 call sites)? Dropping is the majority and empty path entries are noise, but `corehelpers/inventory.go:419` composes `nonblankLines(existing)` with a `protected` slice that is not pre-filtered, so the change is visible there. No `corehelpers` test names `uniqueSorted` at all, so tests will not decide this. Second, and more important: `finalization`'s copy is a path normalizer, not this helper. Either its 24 call sites move to the plain canonical helper and lose `Clean`/`ToSlash`/nil-on-invalid, or it stays and is renamed (`normalizedUniquePaths` reads correctly and satisfies Naming for Reach) so the name collision disappears and the count still reaches 6. **The rename is the safer default and the one I would plan for.**

**D2 — `requestIDLess` parser strictness.** Canonical is `repositorymodel`'s lenient regex parser. `nextselection` currently rejects any identifier with a non-digit after the prefix and falls back to string comparison. Record that `nextselection` adopts the lenient parse and that malformed identifiers now sort numerically by their leading digits.

**D3 — `physicalPath` missing-ancestor paths.** Canonical is ancestor-walking plus `Abs`. Record that `suiteinstall`'s nonexistent-skill-root failure moves from `update_transaction.go:179` to the `SKILL.md` check at line 188, with the message change spelled out.

**D4 — `compareSemver` parse strictness.** Canonical is the strict three-part parse from `release.go`. Record that `interview_commands.go:820` now scores a malformed minor or patch as equal instead of comparing component by component, and that `semverMajor` at line 811 already blocks fully unparseable input from reaching it.

---

## 6. Lock-in assertion, paste-ready

Append to `/home/user/skill-do-work/_dev/tests/audit-lockins.sh`, immediately before the closing `if [ "$failure_count" -gt 0 ]` block. The file already carries four blocks in this shape (REQ-550 at line 9, REQ-551 at line 48, REQ-549 at line 76, REQ-553 at line 124), and is registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh:29`. Do not change that registration.

```bash
# Finding 4: per-req-duplicate-go-helpers (REQ-557)
# Six helper names carried fifteen definitions across do-work-cli packages. Each now has
# exactly one canonical definition, so this pins the total at 6. Case-insensitive because
# the canonical copies are exported: a re-introduced unexported duplicate is then counted
# alongside them. Equality rather than a ceiling, so a lost canonical definition is caught
# by the same assertion as a regrown duplicate. --count-matches sums per-file totals so a
# file holding two definitions is not undercounted; the awk sum prints 0 when rg finds
# nothing and exits non-zero, which is why the count is summed rather than piped to wc.
helper_definition_count="$(
  rg --glob '*.go' --glob '!*_test.go' --no-filename --count-matches \
    '(?i)^func (uniqueSorted|subtractPaths|requestIDLess|firstError|compareSemver|physicalPath)\(' \
    "$repo_root/skills/do-work/tools/do-work-cli/internal/" 2>/dev/null \
    | awk '{ total += $1 } END { print total + 0 }'
)"
if [ "$helper_definition_count" -ne 6 ]; then
  printf 'FAIL: expected 6 definitions of the six shared do-work-cli helper names (uniqueSorted, subtractPaths, requestIDLess, firstError, compareSemver, physicalPath), found %s\n' \
    "$helper_definition_count" >&2
  failure_count=$((failure_count + 1))
fi
```

I ran this text against HEAD. It prints **15**. I also ran the same shape against a symbol that matches nothing and it prints **0**, so an empty match set degrades to a clean failure rather than an unset variable or an arithmetic error. That matters because the file runs under `set -uo pipefail` with no `-e`, so `rg`'s exit 1 through a pipe would otherwise be the trap here.

Three properties worth keeping when this is edited later:

1. **Case-insensitive.** After the change the canonical copies are exported. A lowercase-only pattern would count 0 and pass a `-le` check while every duplicate quietly regrew.
2. **`-ne 6`, not `-le 6`.** A ceiling passes when a canonical definition is accidentally deleted. Equality catches both directions.
3. **`^func ` anchoring.** The canonical helpers must stay plain package-level functions. Give one a receiver and the line starts `func (`, the count silently drops to 5, and the assertion fails for a reason nobody will guess from the message.

One more constraint the assertion imposes: it reaches 6 only if **both** files missing from the `write_set` are edited. `publication/release.go:210` and `repairvalidation/already_green.go:442` each hold one of the fifteen.

---

## 7. Red-green recipe

**RED, from `/home/user/skill-do-work`:**

```
rg -n --glob '*.go' --glob '!*_test.go' '^func (uniqueSorted|subtractPaths|requestIDLess|firstError|compareSemver|physicalPath)\(' skills/do-work/tools/do-work-cli/internal/ | sort
```

Prints 15 lines. Correct the REQ's Red-Green Proof from fourteen to fifteen before starting, otherwise the handback reports a number that does not match the request.

**Second RED:** paste the assertion block into a scratch copy of `_dev/tests/audit-lockins.sh` and run it. It prints `FAIL: expected 6 definitions ... found 15`.

**GREEN, all four:**

1. The Reproduce command, in its case-insensitive form, prints exactly 6 lines, one per name. The original lowercase-only pattern would print 0 and hide the regression.
2. `bash _dev/tests/audit-lockins.sh` exits 0 and prints `Audit lock-in regressions passed.`
3. From `skills/do-work/tools/do-work-cli`: `go build ./...` and `go vet ./...` exit 0, `go test ./...` green.
4. D1 to D4 are written into the REQ's Implementation Summary with the behaviour each caller keeps.

**Baseline confirmed green at HEAD**, so any failure belongs to this change: `go build ./...` exit 0, `go vet ./...` exit 0, and `go test` over all nine affected packages returns `ok` for `corehelpers` (3.9s), `finalization` (22.8s), `knowledgecommands` (6.7s), `publication` (8.4s), `repositorymodel` (0.03s), `nextselection` (7.6s), `dependencygraph` (0.01s), `suiteinstall` (2.7s), `repairvalidation` (0.9s).

---

## 8. Files to change

New file, needs sign-off first (option O2):

- `skills/do-work/tools/do-work-cli/internal/sharedhelpers/shared_helpers.go` — stdlib-only leaf holding `UniqueSorted`, `SubtractPaths`, `FirstError`, `PhysicalPath`, `CompareSemver`. **Full bodies, never one-line delegates**: the REQ-550 block at `audit-lockins.sh:9-46` already fails any exported function whose body is a single call with no other production caller. Plain funcs, no receivers, so the lock-in regex keeps matching. Under O1 this file disappears and everything below is identical apart from the import path.

Declared in `write_set`:

- `internal/corehelpers/checks.go` — delete `subtractPaths` (827), `uniqueSorted` (846), `firstError` (871). Repoint `uniqueSorted` at 741, 782, 804, 906; `subtractPaths` at 177, 178; `firstError` at 153. Delete `stringSet` (859) if nothing else calls it.
- `internal/finalization/finalization_prepare.go` — delete `subtractPaths` (365), repoint 162. Under D1, rename rather than delete `uniqueSorted` (379).
- `internal/knowledgecommands/interview_commands.go` — delete `uniqueSorted` (1337) and `compareSemver` (1535). Repoint 495, 584, 596, 601, 1323, 1332, 820.
- `internal/knowledgecommands/commands.go` — delete `physicalPath` (147), repoint 144.
- `internal/dependencygraph/dependency_graph.go` — delete `requestIDLess` (307) **and** the now-dead `requestNumber` (316). Repoint 115, 117, 249, 303.
- `internal/repositorymodel/repository_model.go` — export `requestIDLess` (654) as `RequestIDLess`, body unchanged. Repoint 305, 319.
- `internal/nextselection/next_types.go` — delete `requestIDLess` (83), keep `numericID` (60). Repoint `next_targets.go:97`, `next_targets.go:151`.
- `internal/publication/capture_files.go` — delete `firstError` (361), repoint 151.
- `internal/suiteinstall/update_transaction.go` — delete `physicalPath` (209), repoint 175, 179, 199.
- `_dev/tests/audit-lockins.sh` — one assertion block, section 6 above.

**Missing from `write_set` and required:**

- `internal/publication/release.go` — holds `compareSemver` (210). Repoint 26. Leave `parseSemver` (227) alone, `release_mirrors.go:109` still uses it.
- `internal/repairvalidation/already_green.go` — holds `uniqueSorted` (442). Repoint 266, 283, 284, 408, 416. Without this the count lands at 7.
- `internal/corehelpers/inventory.go` — call sites 419, 460.
- `internal/finalization/finalization_discovery.go` — 18 call sites (67, 70, 132, 162, 188, 194, 195, 205, 573, 627, 640, 649, 844, 895, 1422, 1432, 1470, 1518).
- `internal/finalization/finalization_apply.go` — call sites 233, 518, 522.
- `internal/knowledgecommands/memory_commands.go` — `uniqueSorted` at 563, 756; `physicalPath` at 1216, 1227.
- `internal/publication/answer.go` — `firstError` at 307.
- `internal/finalization/finalization_recovery_test.go:574` and `internal/knowledgecommands/memory_commands_test.go:390` — the only two tests naming a deleted private helper directly, which is precisely the carve-out at REQ line 51.

---

## 9. Dependencies, re-checked

**REQ-550** is done: `do-work/archive/REQ-550-collapse-four-exported-one-line-delegates-into-their-targets.md`. Its lock-in block sits at `audit-lockins.sh:9-46` and is the reason the canonical helpers must be full bodies.

**REQ-552** does not gate this on its merits. Its `write_set` (line 18) is `internal/corehelpers/commands.go`, `internal/toolboxcommands/architecture.go` and `_dev/tests/audit-lockins.sh`. It replaces an `exec.Command("find")` probe and an `exec.Command("cp")` with pure Go. It never touches `internal/corehelpers/checks.go` and adds or removes none of the six names. The REQ's stated reason — "so `corehelpers` is settled before helpers move in" — describes a file-level overlap that does not exist. **The only real coupling is `_dev/tests/audit-lockins.sh`**, where both REQs append one assertion. Sequence them or accept a trivial merge on that one file.

**REQ-558** depends on this REQ and also appends to `audit-lockins.sh`. Same coupling.

## 10. Provenance is unverifiable here

`git cat-file -t` returns MISSING for all seven hashes the REQ and audit cite: `dc8a64e3` (audited commit), `83594c5e` (report commit), `761d8e6a`, `01d920dd`, `ac2e3acd`, `625d49aa`, `cf111a50`. The audit was taken against a history this clone does not contain. Every "introduced by" attribution in the REQ should be dropped from the Implementation Summary rather than repeated as fact. The `file:line` evidence at HEAD stands on its own and is what this report relies on.

## 11. Net line delta

The REQ expects −70. That estimate assumed no new package and fourteen definitions. Deleting nine definitions removes roughly 125 lines. Under O2 the new leaf package adds roughly 85 back, plus import lines and call-site qualification across seven extra files, so the realistic net is **−30 to −50**. Under O1 it is closer to the REQ's −70. Neither number is a reason to pick an option, but the handback should not claim −70 without measuring.

---

*Generated by Explore agent*