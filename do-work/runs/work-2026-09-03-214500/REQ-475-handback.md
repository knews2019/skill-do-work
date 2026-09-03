# REQ-475 builder handback

## Builder increment

- Branch: `worktree-agent-REQ-475-confine-all-configured-memory-tree-readers`
- Commit: `0f288b7ccbf454c7c73935a8dd6aa3b8f211932b` (`[REQ-475] Confine configured Memory readers`)
- Worktree was clean after the commit.

## Changed-file manifest

- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go`
- `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands_test.go`
- `skills/do-work-knowledge/actions/memory.md`
- `skills/do-work-knowledge/actions/memory-reference.md`

`commands_test.go` remained unchanged: the handler-level test helper marshals the typed result and invokes both real `resultmodel` text and JSON renderers, so a separate runtime-projection edit was not necessary. No lifecycle metadata, run manifest, queue, release/version/changelog, other REQ, BKB implementation, generated site, or screenshot file changed.

## Implementation

- Added one bounded rooted-read function family around an opened `os.Root`: it validates parent components, refuses linked/special final objects, checks pre-open/opened/post-read identity and metadata, distinguishes missing optional objects, and returns opened-file metadata without a later pathname stat.
- Added bounded, identity-checked `logs/` enumeration and a shared deterministic log inventory for broad recall, lexical recall, status, audit, remember/forget, and bootstrap-adjacent reads.
- Migrated configured working-file, log, ledger, and sentinel reads across read-only and mutating Memory surfaces. Ledger parsing is now separate from acquisition; BKB keeps its pre-existing acquisition behavior.
- Enforced repository containment for read-only as well as mutating configured roots, and refused a linked configured root.
- Propagated the precise configured child path on refusal. No partial hits/evidence or refused target bytes reach the typed result or its text/JSON projections; recall refuses before its best-effort ledger append.
- Documented transport ceilings: working file 64 KiB; each log and Memory ledger 8 MiB; bootstrap sentinel 128 bytes; log directory 4,096 entries. The standing-memory semantic cap remains 2,500 characters.

## RED evidence

Command, run after adding the initial adversarial handler matrix and before the production migration:

```text
gofmt -w internal/knowledgecommands/memory_commands_test.go && go test -count=1 -run 'TestConfiguredMemoryReadersRefuseLinkedObjectsWithoutDisclosingTargetBytes|TestConfiguredMemoryReadersRefuseOutsideRepositoryRoot' ./internal/knowledgecommands
```

Result: exit 1. Representative exact failing assertions:

```text
TestConfiguredMemoryReadersRefuseLinkedObjectsWithoutDisclosingTargetBytes/broad_recall/working:
result = ... Outcome:"success" ... Evidence:["updated: configured-memory-outside-canary-475"], want one failure finding

TestConfiguredMemoryReadersRefuseLinkedObjectsWithoutDisclosingTargetBytes/broad_recall/logs_directory:
result = ... Outcome:"success" ... configured-memory-outside-canary-475 ..., want one failure finding

TestConfiguredMemoryReadersRefuseLinkedObjectsWithoutDisclosingTargetBytes/broad_recall/log_file:
result = ... Outcome:"success" ... configured-memory-outside-canary-475 ..., want one failure finding
```

The status and Memory-audit rows for linked working/log/ledger positions also returned success instead of the asserted single failure; linked ledger bytes changed status/audit evidence. All three outside-root rows returned success and projected `configured-memory-outside-root-canary-475` instead of the asserted failure. The shared assertion was `result = ... Outcome:"success"..., want one failure finding`.

## GREEN and validation evidence

- Initial RED command after implementation: PASS, `ok .../internal/knowledgecommands 1.090s`.
- Final focused package: `go test -count=1 ./internal/knowledgecommands` — PASS, `ok ... 7.574s`.
- Race: `go test -race -count=1 ./internal/knowledgecommands` — PASS, `ok ... 7.757s`.
- Vet: `go vet ./...` — PASS, exit 0 with no output.
- Full module: `go test -count=1 ./...` — PASS; all packages passed, including `internal/knowledgecommands 10.885s`.
- Contract suite: `bash _dev/tests/contract-regressions.sh` from repository root — PASS, ending `Contract regression checks passed.` The action shell-block, prescribed-shell, shipped-package-reference, suite-manifest, timestamp, launcher, hook, estimator, selection, and defensive-surface probes all passed.
- Final limit seam after tightening the test to exercise limit dispatch by filename: `go test -count=1 -run 'TestConfiguredMemoryReadLimitsAreInclusiveAndDirectoryEnumerationIsBounded' ./internal/knowledgecommands` — PASS, `ok ... 0.538s`.
- `git diff --check` — PASS.
- Browser validation: not applicable; REQ-475 changes Go command behavior and Markdown contracts, with no browser/UI surface.

Adversarial coverage now includes broad recall, lexical recall, status, Memory audit, remember, and forget across linked working files, log directories, log files, and relevant ledgers; linked configured roots; FIFO/directory/regular-file special-position substitutions; exact and limit+1 reads for every file class; exact and max+1 directory enumeration; typed/text/JSON canary absence; exact affected paths; and Memory-byte preservation. Ordinary status and audit parity/byte preservation remains green alongside the existing scoring, retained lexical-script differential, ledger-redaction, audit-boundary, and BKB tests.

## Guidance read

- `CLAUDE.md`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`
- `_dev/primes/prime-action-files.md`
- `_dev/primes/lessons-action-files.md`
- `_dev/crew/crew-member-general.md`
- `_dev/crew/crew-member-coding-guardrails.md`
- `_dev/crew/crew-member-backend.md`
- `_dev/crew/crew-member-testing.md`
- `_dev/crew/crew-member-communication-style.md`
- The authoritative REQ, brief, full exploration, REQ-417 archive/remediation context, and exact source/tests/action contracts.

## Decisions, risks, and discoveries

- Missing optional logs, ledger, and sentinel remain ordinary absence. Present links, special objects, oversize objects, identity changes, and outside-repository roots are refusals.
- `memory-audit --engine both` retains deterministic BKB inspection but marks the overall result failed when the Memory read is unsafe. A genuinely absent Memory root retains the prior `Absent` classification.
- Read-only external `--memory-root` paths are now rejected. This is the intentional compatibility change required by the repository-rooted contract.
- The finite log/ledger/directory ceilings were previously unspecified; they are now named, tested at the inclusive boundary, and documented. Repositories with a single log/ledger over 8 MiB or more than 4,096 direct log-directory entries will receive an actionable refusal rather than a partial scan.
- During the expanded matrix, lexical recall and initial forget discovery were found to report only the root on a child refusal. Both now report the exact refused configured child path through the shared typed error.
- `os.OpenRoot` is required by the existing module/toolchain baseline. No new dependency was added.
- The archived REQ-417 re-review file referenced by the REQ remains unavailable in the current tree/Git inventory; the durable archive summary and remediation commit were used, as recorded in exploration. No new follow-up task was discovered.

## Merge guidance

Cherry-pick `0f288b7ccbf454c7c73935a8dd6aa3b8f211932b` into the orchestrator integration lane. Resolve no generated or lifecycle files from this branch: it contains only the four builder files above. After integration, the orchestrator should rerun its integration/full gates and retain ownership of request completion/archive, dependency unblocking, queue/run-manifest stamping, version/changelog/release work, cancellation handling, and any `--heavy` verification required for the overall wave.

## Independent-review remediation

- Remediation commit: `1ace19970e242d4c61409a53e81ab78800fb8065` (`[REQ-475] Verify ledger root before append`)
- Additional changed files: `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go` and `memory_commands_test.go` only.
- Fix: the best-effort `appendMemoryLedger` acquisition now calls the identity-checked `openMemoryRoot` rather than `os.OpenRoot` directly. A deterministic pre-open seam swaps the already-scanned configured root to an outside symlink and proves a successful recall does not mutate the outside ledger or disclose its canary.

Remediation RED:

```text
go test -count=1 -run TestMemoryRecallLedgerAppendRefusesConfiguredRootSwap ./internal/knowledgecommands
--- FAIL: TestMemoryRecallLedgerAppendRefusesConfiguredRootSwap
outside ledger changed after configured-root swap:
before={"event":"configured-memory-ledger-root-swap-canary-475"}
after={"event":"configured-memory-ledger-root-swap-canary-475"}
{"engine":"memory","event":"recall","hits":1,"note":"","query":"command platform","source":"do-work-cli","ts":"2026-09-03T22:19:16Z"}
FAIL
```

Remediation GREEN and gates:

- Target seam: `go test -count=1 -run TestMemoryRecallLedgerAppendRefusesConfiguredRootSwap ./internal/knowledgecommands` — PASS, `ok ... 0.330s`.
- Focused package: `go test -count=1 ./internal/knowledgecommands` — PASS, `ok ... 6.058s`.
- Race: `go test -race -count=1 ./internal/knowledgecommands` — PASS, `ok ... 16.134s`.
- Vet: `go vet ./...` — PASS, exit 0 with no output.
- Full module: `go test -count=1 ./...` — PASS; all packages passed, including `internal/knowledgecommands 14.312s`.
- Contract suite: `bash _dev/tests/contract-regressions.sh` — PASS, ending `Contract regression checks passed.`
- `git diff --check` — PASS.

Updated merge guidance: because the original builder commit is already integrated, cherry-pick only remediation commit `1ace19970e242d4c61409a53e81ab78800fb8065`. For a fresh integration, apply `0f288b7ccbf454c7c73935a8dd6aa3b8f211932b` first. The branch remains limited to the original four-file builder manifest; the remediation itself touches only the two Go paths. Orchestrator-owned lifecycle, queue, run-manifest, version, changelog, release, and integration work remains unchanged.
