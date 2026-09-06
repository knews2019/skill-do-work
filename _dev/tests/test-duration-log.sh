#!/usr/bin/env bash
# Shared append-only timing record for shell and Go test files.

if [ -z "${repo_root:-}" ]; then
  script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  repo_root="$(cd "$script_dir/../.." && pwd)"
fi

duration_log_path="${DO_WORK_TEST_DURATION_LOG:-${repo_root}/do-work/test-durations.tsv}"
duration_run_id="${DO_WORK_TEST_RUN_ID:-standalone-$(date -u +%Y%m%dT%H%M%SZ)-$$}"
duration_other_gate_processes="${DO_WORK_TEST_OTHER_GATE_PROCESSES:-0}"
duration_log_header="run_id\tfile\tseconds\tother_gate_processes"
export DO_WORK_TEST_DURATION_LOG="$duration_log_path"
export DO_WORK_TEST_RUN_ID="$duration_run_id"
export DO_WORK_TEST_OTHER_GATE_PROCESSES="$duration_other_gate_processes"

if [ ! -e "$duration_log_path" ]; then
  duration_header_candidate="$duration_log_path.header.$$"
  printf '%b\n' "$duration_log_header" > "$duration_header_candidate"
  ln "$duration_header_candidate" "$duration_log_path" 2>/dev/null || true
  rm -f -- "$duration_header_candidate"
fi
if [ "$(sed -n '1p' "$duration_log_path" 2>/dev/null)" != "$(printf '%b' "$duration_log_header")" ]; then
  printf 'FAIL: test duration log has an invalid header: %s\n' "$duration_log_path" >&2
  exit 1
fi

record_test_duration() {
  local test_file="$1"
  local elapsed_seconds="$2"
  printf '%s\t%s\t%s\t%s\n' \
    "$duration_run_id" "$test_file" "$elapsed_seconds" "$duration_other_gate_processes" \
    >> "$duration_log_path"
}

# Efficiency baseline measurement extensions (opt-in via DO_WORK_TEST_EFFICIENCY_LOG).
efficiency_log_header="run_id\tcase_name\tcache_condition\trevision\ttoolchain\tconcurrency\twall_seconds\tuser_cpu_seconds\tsys_cpu_seconds\ttotal_cpu_seconds\tgo_accum_seconds\tsubprocess_counts\texit_status"

init_test_efficiency_log() {
  local target_log="${1:-${DO_WORK_TEST_EFFICIENCY_LOG:-}}"
  if [ -z "$target_log" ]; then
    return 0
  fi
  local target_dir
  target_dir="$(dirname "$target_log")"
  [ -d "$target_dir" ] || mkdir -p "$target_dir"
  if [ ! -e "$target_log" ]; then
    local header_candidate="$target_log.header.$$"
    printf '%b\n' "$efficiency_log_header" > "$header_candidate"
    ln "$header_candidate" "$target_log" 2>/dev/null || cp "$header_candidate" "$target_log" 2>/dev/null || true
    rm -f -- "$header_candidate"
  fi
}

record_test_efficiency() {
  local log_path="${DO_WORK_TEST_EFFICIENCY_LOG:-}"
  if [ -z "$log_path" ]; then
    return 0
  fi
  init_test_efficiency_log "$log_path"
  local case_name="$1"
  local cache_condition="$2"
  local revision="$3"
  local toolchain="$4"
  local concurrency="$5"
  local wall_sec="$6"
  local user_cpu="$7"
  local sys_cpu="$8"
  local total_cpu="$9"
  local go_accum="${10}"
  local subproc_counts="${11}"
  local exit_status="${12}"

  printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
    "$duration_run_id" "$case_name" "$cache_condition" "$revision" "$toolchain" "$concurrency" \
    "$wall_sec" "$user_cpu" "$sys_cpu" "$total_cpu" "$go_accum" "$subproc_counts" "$exit_status" \
    >> "$log_path"
}

measure_command_efficiency() {
  if [ "$#" -lt 2 ]; then
    printf 'usage: measure_command_efficiency <case-name> <cache-condition> <command ...>\n' >&2
    return 2
  fi
  local case_name="$1"
  local cache_condition="$2"
  shift 2

  local measure_py
  measure_py="$(cat <<'PY'
import os, sys, time, resource, subprocess, tempfile, shutil, json, pathlib

case_name = sys.argv[1]
cache_condition = sys.argv[2]
repo_root = sys.argv[3]
run_id = sys.argv[4]
eff_log = sys.argv[5]
concurrency = sys.argv[6]
cmd = sys.argv[7:]

shim_dir = tempfile.mkdtemp(prefix="eff-shim-")
counts_file = os.path.join(shim_dir, "counts.log")
intercepted = ["git", "go", "bash", "python3", "shellcheck", "shasum", "find", "curl"]
orig_path = os.environ.get("PATH", "")
path_dirs = [p for p in orig_path.split(os.path.pathsep) if p != shim_dir]

for b in intercepted:
    real_bin = None
    for d in path_dirs:
        cand = os.path.join(d, b)
        if os.path.isfile(cand) and os.access(cand, os.X_OK):
            real_bin = cand
            break
    if real_bin:
        shim_path = os.path.join(shim_dir, b)
        with open(shim_path, "w") as f:
            f.write(f"#!/bin/sh\necho {b} >> \"{counts_file}\"\nexec \"{real_bin}\" \"$@\"\n")
        os.chmod(shim_path, 0o755)

sub_env = dict(os.environ)
sub_env["PATH"] = shim_dir + os.path.pathsep + orig_path
if concurrency:
    sub_env["GOMAXPROCS"] = concurrency

t0 = time.perf_counter()
ru0 = resource.getrusage(resource.RUSAGE_CHILDREN)
proc_exit = 0
out = ""
err = ""

try:
    proc = subprocess.run(cmd, env=sub_env, capture_output=True, text=True)
    proc_exit = proc.returncode
    out = proc.stdout
    err = proc.stderr
finally:
    ru1 = resource.getrusage(resource.RUSAGE_CHILDREN)
    t1 = time.perf_counter()

wall = t1 - t0
user_cpu = ru1.ru_utime - ru0.ru_utime
sys_cpu = ru1.ru_stime - ru0.ru_stime
total_cpu = user_cpu + sys_cpu

counts = {}
if os.path.exists(counts_file):
    with open(counts_file) as f:
        for line in f:
            n = line.strip()
            if n:
                counts[n] = counts.get(n, 0) + 1
shutil.rmtree(shim_dir, ignore_errors=True)
counts_str = ",".join(f"{k}:{v}" for k, v in sorted(counts.items())) or "none"

go_accum = "n/a"
pkg_wall = "n/a"

if "-json" in cmd or any("Test" in arg for arg in cmd):
    durations = []
    for line in out.splitlines():
        try:
            ev = json.loads(line)
        except Exception:
            continue
        test = ev.get("Test")
        action = ev.get("Action")
        elapsed = float(ev.get("Elapsed", 0.0))
        if test and "/" not in test and action in {"pass", "fail"}:
            durations.append(elapsed)
        elif not test and action in {"pass", "fail"} and "Package" in ev:
            pkg_wall = f"{elapsed:.3f}"

    if durations:
        go_accum = f"{sum(durations):.3f}"

rev = "unknown"
try:
    rev = subprocess.check_output(["git", "-C", repo_root, "rev-parse", "--short", "HEAD"], text=True).strip()
except Exception:
    pass

toolchain = "unknown"
try:
    gver = subprocess.check_output(["go", "version"], text=True).strip().split()[2]
    toolchain = f"{gver},py:{sys.version.split()[0]}"
except Exception:
    pass

if eff_log:
    row = f"{run_id}\t{case_name}\t{cache_condition}\t{rev}\t{toolchain}\t{concurrency}\t{wall:.3f}\t{user_cpu:.3f}\t{sys_cpu:.3f}\t{total_cpu:.3f}\t{go_accum}\t{counts_str}\t{proc_exit}\n"
    header = "run_id\tcase_name\tcache_condition\trevision\ttoolchain\tconcurrency\twall_seconds\tuser_cpu_seconds\tsys_cpu_seconds\ttotal_cpu_seconds\tgo_accum_seconds\tsubprocess_counts\texit_status\n"
    if not os.path.exists(eff_log):
        os.makedirs(os.path.dirname(os.path.abspath(eff_log)), exist_ok=True)
        with open(eff_log, "w") as f:
            f.write(header)
    with open(eff_log, "a") as f:
        f.write(row)

res = {
    "case": case_name,
    "condition": cache_condition,
    "wall_seconds": round(wall, 3),
    "user_cpu_seconds": round(user_cpu, 3),
    "sys_cpu_seconds": round(sys_cpu, 3),
    "total_cpu_seconds": round(total_cpu, 3),
    "go_accum_seconds": go_accum,
    "package_wall_seconds": pkg_wall,
    "subprocess_counts": counts,
    "exit_status": proc_exit
}
print(json.dumps(res))
if proc_exit != 0 and err:
    print(err, file=sys.stderr)
sys.exit(proc_exit)
PY
)"

  python3 -c "$measure_py" \
    "$case_name" \
    "$cache_condition" \
    "$repo_root" \
    "${duration_run_id}" \
    "${DO_WORK_TEST_EFFICIENCY_LOG:-}" \
    "${GOMAXPROCS:-4}" \
    "$@"
}
