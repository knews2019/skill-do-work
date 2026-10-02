# Builder brief — REQ-625: Add a low-disk-space probe to the board VERIFY band

**Worktree (your working directory, the only tree you write):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-625-low-disk-space-verify-probe`
**Branch:** `worktree-agent-REQ-625-low-disk-space-verify-probe`, based on `8812a445`. Commit on this branch only, message prefix `[REQ-625]`.
**Route:** B. **TDD: yes** (RED before GREEN, evidence required). **Impact:** user-visible. **Estimate P50:** 35 active minutes.
**Hand-back file (the ONE main-tree path you may write):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-02-140440/REQ-625-handback.md`. Never stage or commit it.

## Read first (absolute paths; the REQ and exploration live in the MAIN tree, read-only)

1. The REQ: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-625-low-disk-space-verify-probe.md` — What, Detailed Requirements, Out of scope, Red-Green Proof, Exploration, Scope. The `## Scope` "Files I will touch" list is your write boundary.
2. Exploration with file:line anchors: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-02-140440/REQ-625-exploration.md`.
3. In your worktree: `_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`, and its lessons satellite `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (binding for anything touching the board).
4. Crew rules in your worktree under `skills/do-work/crew-members/`: `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `backend.md`, `testing.md`. Naming rule: two words minimum for anything with reach.

## Fixed design (decided by triage and exploration; do not re-open)

- **Constants** in `verify.go`: `lowDiskSpaceWarningBytes = 10 << 30` and `lowDiskSpaceCriticalBytes = 3 << 30`, documented. "Below" is strict: free < threshold. At exactly 3 GiB the finding is warning-level; at exactly 10 GiB there is no finding.
- **Category** `verifyCategoryLowDiskSpace = "low-disk-space"` in the existing const block. No registration anywhere else (the board suppression list and renderers need nothing).
- **Measurement type and seam.** A small struct (free bytes, total bytes, device id, all uint64) and a per-platform `measureDiskSpace(directory string) (<struct>, error)`:
  - `disk_space_unix.go` with `//go:build aix || darwin || dragonfly || freebsd || ios || linux`: `syscall.Statfs`, free = `uint64(Bavail) * uint64(Bsize)` (clamp a negative Bavail to 0), total = `uint64(Blocks) * uint64(Bsize)`, device id = `uint64(stat.Dev)` from `os.Stat(directory).Sys().(*syscall.Stat_t)`. Every field conversion stays in this file; types differ per OS.
  - `disk_space_windows.go` with `//go:build windows`: `GetDiskFreeSpaceExW` via `syscall.NewLazyDLL("kernel32.dll")`, following `atomic_replace_windows.go` for UTF-16 conversion and error handling; three `*uint64` out-pointers. Device id: FNV-64 of the upper-cased `filepath.VolumeName` of the absolute path (document that mount points below a drive are not distinguished).
  - `disk_space_unsupported.go` with the negation of both tag sets, returning a typed sentinel `errDiskSpaceUnsupported` (wrap it with `%w` so `errors.Is` works).
  - `appendDiskSpaceFindings(report *VerifyReport, repoRoot string, worktreePathsByName map[string]string, measure func(string) (<struct>, error))` takes the measurement function as a parameter; unit tests call it directly with fakes. `collectVerifyFindings` passes a package-level `var diskSpaceMeasurer = measureDiskSpace` (pick a two-word name that does not collide with the type) after `appendWorktreeFindings`, enumerating worktrees with `listWorktreeAgentWorktrees` only when `gitBinaryAvailable()`; an enumeration error becomes a SkippedProbes line for the worktree half and the repo root is still measured.
- **Finding shape.** One finding per device id, measuring the repo root first, then worktree names in sorted order. `Subject` is the absolute `repoRoot` for the root (reduces to `.`) and the worktree NAME (the map key) for a worktree — never the worktree path, which lies outside the repo. `Fixable: false`. Detail: `"<free> free of <total> (<pct>%) — below the <warning|critical> threshold (<threshold>)"` with a new `formatGibibytes` helper (`"%.1f GiB"`); Detail must not start with the Subject. Remedy exactly: `free space: clear regenerable QA output, finished builder worktrees (do-work cleanup), browser caches; ` + "`du -sh * | sort -h`" + ` at the repo root shows the largest directories`.
- **Skips.** `errDiskSpaceUnsupported` → one SkippedProbes entry `"disk-space probe: unsupported on <GOOS>"` and return, zero findings. Any other measurement error → `"disk-space probe for <directory>: <err>"` for that directory only; keep going.
- **Read-only.** Measure and report; delete nothing.
- **Tests** in a NEW file `disk_space_test.go` (keep `verify_test.go` untouched for the per-file 30 s budget), each named for the failure it pins: above/below thresholds with fakes (2 GiB critical, 9 GiB warning, 20 GiB clean); exact 3 GiB and exact 10 GiB boundaries; two worktree names on one device → one finding; unsupported sentinel → one skip, zero findings; per-directory error → skip for that directory, root still reported. In `generate_test.go`: install a plenty-free fake for the package variable in the existing `TestMain` (so the three real-disk zero-finding tests stay deterministic), and extend `TestGeneratedVerifyPayloadCarriesNoAbsolutePaths` to seed a low-space fake for the repo root plus a worktree name, assert a `low-disk-space` entry is present, and add `Subject` to the texts it scans for absolute paths.
- **Docs.** `prime-do-kanban.md`: one Traps bullet near the verify traps (lines 31-34) naming the probe, its thresholds, the Subject rule and the narrow unix tag list. `lessons-do-kanban.md`: prepend at the top of the entry list `- [family: disk-space-blind-spot] 0.305.60: **<bold lesson>.** <prose>` — a long fan-out run with browser QA can fill a disk in hours while nothing in the suite noticed; the board is where the operator looks, so the probe lives in verify. Do NOT edit `do-work/lessons-index.md` (orchestrator's).

## Never touch

Anything under `do-work/` in either tree except your hand-back file; `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md` (release is the orchestrator's); `go.mod`, `go.sum`; `verify_test.go`; `atomic_replace_*.go`; `web/`; `work.md`, `work-reference.md`. A file outside Scope that you believe you need: stop and report it in the hand-back instead of editing it.

## Testing (record each command's own exit line)

From `skills/do-work-board/tools/queue-kanban` in your worktree:

```
QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 -run 'DiskSpace|GeneratedVerifyPayloadCarriesNoAbsolutePaths|Verify' .
QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 ./...
gofmt -l . ; go vet .
GOOS=windows GOARCH=amd64 go vet . && GOOS=linux go vet . && GOOS=freebsd go vet . && GOOS=openbsd go vet . && GOOS=netbsd go vet . && GOOS=solaris go vet . && GOOS=js GOARCH=wasm go vet .
```

TDD order: write `disk_space_test.go` first against a stub, run it, capture the failing output (the RED evidence), then implement, then rerun (GREEN). Keep the RED output verbatim for the hand-back.

## Hand-back format (write the file, then reply with ONE line)

```
# Hand-back — REQ-625
**Branch:** worktree-agent-REQ-625-low-disk-space-verify-probe  **Base:** 8812a445  **Head:** <hash>
## File manifest
- path (new|modified|deleted) — one line each
## Integration seams
None, or the exact line and the file it belongs in
## P-A-U
- [PLAN] <two or three sentences of the approach you took>
- [APPLY] <confirmation scope stayed inside the declared files>
- [UNIFY] <git diff --stat output, linters run, files checked>
## Red-green evidence
- <test name>: RED output (verbatim, trimmed) → GREEN at <commit>
## Tests run
- <command> → exit <n>, <n> tests  (one per command above, including every cross-vet GOOS)
## Decisions
- D-01: ... (reasoning; mark DECIDE & STATE or ESCALATE with Value/Risk)
## Discovered Tasks
- ... or None
## Lessons read
- each required/touched lessons file and whether whole-satellite
## Blockers
- None, or what stopped you
```
