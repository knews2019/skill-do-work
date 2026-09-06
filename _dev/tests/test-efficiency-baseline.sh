#!/usr/bin/env bash
# Reproducible test-efficiency baseline runner for representative cases.
# Measures cold/warm execution, descendant CPU, wall time, Go test phase attribution,
# and intercepted subprocess counts.
set -euo pipefail

script_path="${BASH_SOURCE[0]}"
script_directory="${script_path%/*}"
if [ "$script_directory" = "$script_path" ]; then
  script_directory='.'
fi
repo_root="$(cd "$script_directory/../.." && pwd)"

# shellcheck source=_dev/tests/test-duration-log.sh disable=SC1091
source "$repo_root/_dev/tests/test-duration-log.sh"

runs=3
concurrency=4
target_case="all"
tsv_output="${DO_WORK_TEST_EFFICIENCY_LOG:-$repo_root/do-work/test-efficiency.tsv}"

while [ "$#" -gt 0 ]; do
  case "$1" in
    --runs)
      runs="$2"
      shift 2
      ;;
    --concurrency)
      concurrency="$2"
      shift 2
      ;;
    --case)
      target_case="$2"
      shift 2
      ;;
    --output)
      tsv_output="$2"
      shift 2
      ;;
    --quick)
      runs=1
      shift
      ;;
    *)
      printf 'usage: %s [--runs N] [--concurrency N] [--case <name|all>] [--output <file>] [--quick]\n' "$0" >&2
      exit 2
      ;;
  esac
done

export GOMAXPROCS="$concurrency"
export DO_WORK_TEST_EFFICIENCY_LOG="$tsv_output"
init_test_efficiency_log "$tsv_output"

raw_metrics_file="$(mktemp "${TMPDIR:-/tmp}/baseline-metrics.XXXXXX")"
cleanup_metrics() {
  rm -f -- "$raw_metrics_file"
}
trap cleanup_metrics EXIT

record_metric_json() {
  printf '%s\n' "$1" >> "$raw_metrics_file"
}

run_case_samples() {
  local case_id="$1"
  local cond="$2"
  shift 2

  for ((i = 1; i <= runs; i++)); do
    if [ "$cond" = "cold" ]; then
      go clean -testcache 2>/dev/null || true
    fi
    local sample_json
    sample_json="$(measure_command_efficiency "$case_id" "$cond" "$@")"
    record_metric_json "$sample_json"
  done
}

run_all_cases() {
  # 1. inventory (A3)
  if [ "$target_case" = "all" ] || [ "$target_case" = "inventory" ]; then
    printf 'Measuring case: inventory (cold, %d runs)...\n' "$runs" >&2
    run_case_samples "inventory" "cold" \
      go test -C "$repo_root/skills/do-work/tools/do-work-cli" -json -count=1 -run '^TestInventory' ./internal/corehelpers/...
    printf 'Measuring case: inventory (warm, %d runs)...\n' "$runs" >&2
    run_case_samples "inventory" "warm" \
      go test -C "$repo_root/skills/do-work/tools/do-work-cli" -json -count=1 -run '^TestInventory' ./internal/corehelpers/...
  fi

  # 2. finalization (A4)
  if [ "$target_case" = "all" ] || [ "$target_case" = "finalization" ]; then
    printf 'Measuring case: finalization (cold, %d runs)...\n' "$runs" >&2
    run_case_samples "finalization" "cold" \
      go test -C "$repo_root/skills/do-work/tools/do-work-cli" -json -count=1 -run '^TestRecoverFinalization' ./internal/finalization/...
    printf 'Measuring case: finalization (warm, %d runs)...\n' "$runs" >&2
    run_case_samples "finalization" "warm" \
      go test -C "$repo_root/skills/do-work/tools/do-work-cli" -json -count=1 -run '^TestRecoverFinalization' ./internal/finalization/...
  fi

  # 3. session-start (A5)
  if [ "$target_case" = "all" ] || [ "$target_case" = "session-start" ]; then
    printf 'Measuring case: session-start (%d runs)...\n' "$runs" >&2
    run_case_samples "session-start" "warm" \
      bash "$repo_root/_dev/tests/session-start-hook-behavior.sh"
  fi

  # 4. shell-audit (A6)
  if [ "$target_case" = "all" ] || [ "$target_case" = "shell-audit" ]; then
    printf 'Measuring case: shell-audit (%d runs)...\n' "$runs" >&2
    run_case_samples "shell-audit" "warm" \
      bash -c "bash '$repo_root/_dev/tests/action-shell-blocks.sh' && bash '$repo_root/_dev/tests/quiet-grep-pipeline-audit.sh'"
  fi

  # 5. heavy-cli-build (A2)
  if [ "$target_case" = "all" ] || [ "$target_case" = "heavy-cli-build" ]; then
    printf 'Measuring case: heavy-cli-build (cold, %d runs)...\n' "$runs" >&2
    run_case_samples "heavy-cli-build" "cold" \
      env DO_WORK_HEAVY_TESTS=1 go test -C "$repo_root/skills/do-work/tools/do-work-cli" -json -count=1 -run '^TestBuiltInstall' ./internal/suiteinstall/...
    printf 'Measuring case: heavy-cli-build (warm, %d runs)...\n' "$runs" >&2
    run_case_samples "heavy-cli-build" "warm" \
      env DO_WORK_HEAVY_TESTS=1 go test -C "$repo_root/skills/do-work/tools/do-work-cli" -json -count=1 -run '^TestBuiltInstall' ./internal/suiteinstall/...
  fi

  # 6. reused-stage (A8)
  if [ "$target_case" = "all" ] || [ "$target_case" = "reused-stage" ]; then
    printf 'Measuring case: reused-stage (%d runs)...\n' "$runs" >&2
    run_case_samples "reused-stage" "reused" \
      bash "$repo_root/_dev/tests/fast-stage-reuse-behavior.sh"
  fi

  # 7. go-discovery (A7)
  if [ "$target_case" = "all" ] || [ "$target_case" = "go-discovery" ]; then
    printf 'Measuring case: go-discovery (%d runs)...\n' "$runs" >&2
    run_case_samples "go-discovery" "warm" \
      env DO_WORK_TEST_REPO_ROOT="$repo_root" QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior \
      bash "$repo_root/_dev/tests/run-go-tests-with-budget.sh" "$repo_root/skills/do-work-board/tools/queue-kanban" -run '^TestVerify' ./...
  fi
}

run_all_cases

# Aggregate statistics and format evidence table
python3 - "$raw_metrics_file" "$repo_root" "$tsv_output" "$concurrency" "$runs" <<'PY'
import json, sys, subprocess, os

metrics_file = sys.argv[1]
repo_root = sys.argv[2]
tsv_path = sys.argv[3]
concurrency = sys.argv[4]
runs = int(sys.argv[5])

entries = []
with open(metrics_file, "r", encoding="utf-8") as f:
    for line in f:
        line = line.strip()
        if line:
            try:
                entries.append(json.loads(line))
            except Exception:
                pass

groups = {}
for e in entries:
    key = (e["case"], e["condition"])
    groups.setdefault(key, []).append(e)

rev = subprocess.check_output(["git", "-C", repo_root, "rev-parse", "HEAD"], text=True).strip()
dirty = subprocess.check_output(["git", "-C", repo_root, "status", "--porcelain"], text=True).strip()
dirty_summary = f"{len(dirty.splitlines())} modified/untracked paths" if dirty else "clean"

go_ver = "unknown"
try:
    go_ver = subprocess.check_output(["go", "version"], text=True).strip()
except Exception:
    pass

py_ver = f"Python {sys.version.split()[0]}"
bash_ver = "unknown"
try:
    bash_ver = subprocess.check_output(["bash", "--version"], text=True).splitlines()[0]
except Exception:
    pass

descriptions = {
    "inventory": ("A3: Inventory synthetic status matrix", "Replace porcelain subprocess execution with in-process byte parsing for 45 status matrix cases"),
    "finalization": ("A4: Recovery fixture setup", "Copy prepared recovery template states instead of repeating seed commits and Git history setup"),
    "session-start": ("A5: Fixture integrity verification", "Optimize full-tree find/shasum/sort hashing pipelines run after each case"),
    "shell-audit": ("A6: Shell blocks & quiet-grep lints", "Batch compatible ShellCheck files and consolidate redundant AST/grep passes"),
    "heavy-cli-build": ("A2: Integration CLI builds", "Compile test CLI binary once per tested source rather than rebuilding inside 3 test functions"),
    "reused-stage": ("A8: Fast-stage evidence reuse", "Profile fingerprint sealing and avoid redundant reads of identical resolved inputs")
}

priorities = {
    "inventory": "P1 (A3)",
    "heavy-cli-build": "P2 (A2)",
    "finalization": "P3 (A4)",
    "session-start": "P4 (A5)",
    "shell-audit": "P5 (A6)",
    "reused-stage": "P6 (A8)"
}

def calc_median_spread(vals):
    if not vals:
        return 0.0, 0.0
    s = sorted(vals)
    n = len(s)
    med = s[n // 2] if n % 2 == 1 else (s[n // 2 - 1] + s[n // 2]) / 2.0
    spread = s[-1] - s[0]
    return med, spread

print("\n# Test-Efficiency Baseline Evidence Table")
print(f"\n- **Source Revision:** `{rev}` ({dirty_summary})")
print(f"- **Toolchain:** `{go_ver}` | `{py_ver}` | `{bash_ver}`")
print(f"- **Fixed Concurrency:** `GOMAXPROCS={concurrency}` across all runs")
print(f"- **Sample Strategy:** {runs} comparable runs per selection, reporting median ± spread (`max - min`)")
print(f"- **Detailed TSV Log:** `{tsv_path}`\n")

print("| Priority | Case | Target Selection | Condition | Runs | Wall Time (median ± spread) | Total CPU (median ± spread) | Go Accum vs Pkg Wall | Subprocess Counts | Removable Work / Primary Bottleneck |")
print("| --- | --- | --- | --- | ---: | ---: | ---: | ---: | --- | --- |")

for key, items in sorted(groups.items(), key=lambda x: (priorities.get(x[0][0], "P9"), x[0][1])):
    c_id, cond = key
    prio = priorities.get(c_id, "-")
    target_desc, removable = descriptions.get(c_id, (c_id, "-"))
    walls = [item["wall_seconds"] for item in items]
    cpus = [item["total_cpu_seconds"] for item in items]
    med_w, spr_w = calc_median_spread(walls)
    med_c, spr_c = calc_median_spread(cpus)

    # Subprocess counts median
    subproc_keys = set()
    for item in items:
        subproc_keys.update(item.get("subprocess_counts", {}).keys())
    subproc_parts = []
    for sk in sorted(subproc_keys):
        k_vals = [item.get("subprocess_counts", {}).get(sk, 0) for item in items]
        k_med, _ = calc_median_spread(k_vals)
        subproc_parts.append(f"{sk}:{int(k_med)}")
    subproc_str = ", ".join(subproc_parts) if subproc_parts else "none"

    # Go accum vs pkg wall
    go_accums = []
    pkg_walls = []
    for item in items:
        ga = item.get("go_accum_seconds")
        pw = item.get("package_wall_seconds")
        if ga and ga != "n/a":
            try:
                go_accums.append(float(ga))
            except Exception:
                pass
        if pw and pw != "n/a":
            try:
                pkg_walls.append(float(pw))
            except Exception:
                pass

    if go_accums and pkg_walls:
        med_ga, _ = calc_median_spread(go_accums)
        med_pw, _ = calc_median_spread(pkg_walls)
        accum_vs_pkg = f"{med_ga:.2f}s accum / {med_pw:.2f}s pkg"
    elif go_accums:
        med_ga, _ = calc_median_spread(go_accums)
        accum_vs_pkg = f"{med_ga:.2f}s accum"
    else:
        accum_vs_pkg = "n/a (shell)"

    print(f"| {prio} | `{c_id}` | {target_desc} | `{cond}` | {len(items)} | {med_w:.2f}s (±{spr_w:.2f}s) | {med_c:.2f}s (±{spr_c:.2f}s) | {accum_vs_pkg} | `{subproc_str}` | {removable} |")

print("\n### Unsupported Counters")
print("- **Kernel-level untraced execve calls:** Capturing non-PATH absolute binary executions (e.g. `/bin/ps`, `/usr/bin/awk`, `/bin/date`) and internal Go runtime thread forks without exec is unsupported on macOS without elevated root privileges / disabled System Integrity Protection (SIP).")
print("- **Opt-in Intercepted Binaries:** PATH shims record exact counts for: `git`, `go`, `bash`, `python3`, `shellcheck`, `shasum`, `find`, `curl`.")
print("\n### Measurement Observations & Removable Work Ranking")
print("1. **A3 (Inventory):** Significant process spawn churn (`git` and `do-work-cli`) across 45 synthetic status combinations. In-process byte testing directly cuts this overhead.")
print("2. **A2 (Heavy CLI Build):** Calling `buildTestCLIBinary` 3 times inside integration tests invokes redundant Go builds and links. Caching/sharing one compiled binary saves multiple compilation seconds.")
print("3. **A4 (Finalization Recovery):** Repeated Git repository construction (`seedSemanticLegacyTail`) spends seconds creating history from scratch. Pre-baking template states saves repetitive git commits.")
print("4. **A5 (Session-Start):** After each test case, traversal and multi-stage hashing (`find | shasum | sort | shasum`) dominate hook execution. A single-process or manifest digest check removes 20+ subprocesses per run.")
print("5. **A6 (Shell Audit):** Running ShellCheck per file/fence starts the linter dozens of times. Batching files into fewer invocations reduces startup overhead.")

PY
