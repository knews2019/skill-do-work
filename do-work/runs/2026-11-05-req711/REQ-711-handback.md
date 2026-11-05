# REQ-711 runtime fingerprint fixture repair handback

Implemented and committed the fixture repair. Focused verification is green with both inherited host settings and a temporary hostile GOENV file. Parent integration and the canonical maintainer gate remain pending.

## Integration identity

- Checkout: `/Users/t2/.t3/worktrees/skill-do-work2/t3-715b4f5b`
- Branch: `t3/runtime-fingerprint-fixture-repair`
- Commit: `3c6efd59ec172b8481cbca0ccdacc044573bd078`
- Subject: `[REQ-711] isolate runtime evidence fixture Go settings`
- Starting HEAD: `792cf17de031b0784584da3d88701892fd5a8c2f`
- Working tree clean after commit. No PR created or pushed. Parent is sole integrator.

## File manifest

- Modified and committed: `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go` — four added lines in `newRuntimeEvidenceRepository`: one explanatory comment and `t.Setenv("GOENV", "off")`, `t.Setenv("GOFLAGS", "")`, `t.Setenv("GOWORK", "off")` before fixture construction and runtime probes.
- Created outside the implementation checkout, as explicitly authorized: this handback only.

No production code, release metadata, request state, or other main-checkout files were edited. The changed test file has `export-ignore: set`.

## RED/GREEN evidence

Commands ran from `skills/do-work/tools/do-work-cli` using `go1.27.3 darwin/arm64`. Each test invocation used `-count=1`, Go's `-timeout=25s`, and an outer Python `subprocess.run(..., timeout=28)` deadline. No invocation approached 30 seconds.

RED, before editing:

```sh
go test -count=1 -timeout=25s -v ./internal/heavyverification -run '^(TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes|TestShippedGitIsolationPreservesGenericLaneInheritance)$'
```

Exit 1; wall 2.317s; package 2.165s. Runtime subcases `GOOS`, `GOARCH`, `CGO_ENABLED`, and `DO_WORK_TEST_DO_WORK_CLI_BINARY` failed with `default runtime must have a determinable fingerprint`. `GOENV` passed. Git isolation failed with `shipped runtime probe did not isolate host Git configuration`. This independently reproduced the authoritative request's prior baseline failure.

GREEN, after editing, with the same inherited host configuration:

```sh
go test -count=1 -timeout=25s -v ./internal/heavyverification -run '^(TestShipped|TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint)'
```

Exit 0; wall 4.552s; package 4.254s. All six top-level tests in the changed file passed, including all five effective-setting/binary subcases, Git isolation, and `TestShippedRuntimeRefusesOpaqueOverrides`.

Host-isolation replay: created a temporary GOENV file containing exactly `GOFLAGS=-trimpath\n`, supplied only its GOENV path to the subprocess environment, and reran the same GREEN command. `go env` first confirmed effective `GOFLAGS=-trimpath`. Exit 0; wall 4.362s; package 4.227s. The same six tests and five runtime subcases passed. The temporary directory was removed automatically.

Ran `gofmt -w` on the declared file and `git diff --check` successfully. Final diff: one file, four insertions. Verified no diff in `_dev/tests/heavy-runtime-fingerprint.py`.

## Host configuration evidence

Before RED, process-level `GOENV`, `GOFLAGS`, and `GOWORK` were all unset. Effective `go env -json GOENV GOFLAGS GOWORK` reported:

```json
{"GOENV":"/Users/t2/Library/Application Support/go/env","GOFLAGS":"-trimpath","GOWORK":""}
```

The same settings remained after GREEN and the temporary-file replay. The persistent GOENV file SHA-256 before and after the replay was identical: `d1e002f0702167b58f97c02dc83c2e6585ae559640819bc3cb164e99be62372c`. No `go env -w` or persistent configuration mutation was performed. The ordinary RED/GREEN invocations did not externally clear host settings.

## Decisions and discoveries

- Existing failing tests supply the regression proof. No additional test was necessary; the temporary GOENV replay proves host-file isolation without adding redundant coverage.
- Fixture setup precedes all deliberate per-test overrides. The GOENV-content mutation and effective-setting/binary invalidation assertions still run and pass.
- Production refusal of nonempty effective GOFLAGS and external GOWORK is unchanged. Existing opaque-override tests pass; no separate external-GOWORK scenario was added or run.
- The apparent Git-isolation defect was caused by Go configuration contamination, consistent with the supplied diagnosis.
- Full maintainer gate and integration were intentionally not run; the parent owns those steps. Focused green is not a claim that the canonical gate has passed.

## Guidance and lesson paths read

Read `CLAUDE.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, and the authoritative main request at `do-work/working/REQ-711-repair-runtime-fingerprint-baseline.md` (read-only).

Read touch-relevant entries in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: `time-is-not-authorization`, `shipped-module-test-self-containment`, `fixture-cost-is-subprocess-spawning`, the toolchain build-staleness/host `GOFLAGS=-trimpath` lesson, heavyverification package routing, Verify, and Scoped command contracts. This was targeted lesson coverage, not a claim to read the entire satellite.

Read crews under `skills/do-work/crew-members/`: `general.md`, `coding-guardrails.md`, `testing.md`, `shared-principles.md`, `communication-style.md`, and `anti-slop.md`. Also read `_dev/primes/prime-releases.md`; no release changes were made.

## Discovered Tasks

None beyond the scoped repair.
