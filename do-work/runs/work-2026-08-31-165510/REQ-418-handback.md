# REQ-418 Handback — Toolbox Commands and Absorbed Audit Metrics

## Outcome

Implemented the seven canonical shared-CLI surfaces on branch `worktree-agent-REQ-418-toolbox-migration` from base `b713fa8b`:

- `do-work-note`
- `architecture-report-preflight`
- `generate-report-image`
- `generate-report-image-batch`
- `publish-portfolio-summary`
- `install-last30days`
- `audit-metrics inventory|folders|churn|hotspots`

The shared result has exact compatibility text, typed normalized audit JSON, and a narrow cleaned-interruption exit override. Mutators use exact-path Git transactions where their state is committable; last30days remains a private ignored tree and refuses meaningless `--commit`. Report images use one in-process engine, inert argv prompts, adjacent/private staging, Unix process groups with TERM/grace/KILL/reap cleanup, concurrent batch workers, status-backed nonempty publication, all-failed success, and Windows fail-closed ownership.

Retained toolbox scripts and standalone audit-metrics sources are unchanged. REQ-419 actions/recipes and REQ-420 shims/removal are not included. No lifecycle file or REQ-462-owned inventory file changed.

## Scope

- Exact changed set: all and only the frozen 30 paths (6 existing plus 24 new).
- New package: `skills/do-work/tools/do-work-cli/internal/toolboxcommands`.
- REQ-462 overlap: none; both `internal/corehelpers/inventory*` paths remain untouched.
- Scratch artifacts are untracked and were not staged.

## Verification

Green:

- focused toolbox/result/runtime tests;
- `go test -race ./internal/toolboxcommands`;
- full do-work-cli `go test -count=1 ./...` and `go vet ./...`;
- Go 1.25 compatibility;
- Windows toolbox compile/fail-closed process boundary;
- retained audit-metrics `go test -count=1 ./...` and `go vet ./...`;
- exact status/Markdown differential for inventory, folders, churn, and hotspots on the real repository;
- real-CLI registration proof for all seven names;
- retained architecture/image/image-batch/portfolio/last30days focused fixture files;
- full prescribed-shell suite (118 named cases), contract regressions, staged-skills contract;
- installer and updater behavior suites;
- exact 30/30 scope comparison and `git diff --check`;
- unpiped `_dev/tests/maintainer-verify.sh` (browser lane explicitly skipped by the canonical harness because no browser was available; the gate passed).

Focused adversarial coverage includes note dirty/dry-run/exact commit, portfolio inode separation/suffixing, architecture collision allocation, media prompt inertness/all-failed cleanup/process-group escalation, last30days target-Python absence, audit distributions/bands/rename-copy resolution, normalized JSON collections, and interruption override restriction.

## Integration notes

- Registering the toolbox handler map is isolated at `cmd/do-work-cli/main.go`.
- Compatibility callers should continue using retained sources until REQ-419/REQ-420 land; this branch deliberately does not redirect them.
- The standalone audit binary remains the differential oracle and produced byte-identical Markdown for all four modes.
