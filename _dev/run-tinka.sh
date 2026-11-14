#!/usr/bin/env bash
# Start the local Tinka canvas for this repository, stopping any existing managed
# instance or process blocking the configured TCP port first, then open the board
# directly in the default browser.
# Usage: run-tinka.sh <project-root> [board-id|home]
set -u

# Tinka's runtime.mjs calls `ps -p <pid> -o command=` without `-ww`; on macOS BSD ps
# truncates at $COLUMNS (80 in non-TTY / narrow panes) and breaks PID validation.
export COLUMNS=4096

project_root="${1:?usage: run-tinka.sh <project-root> [board-id|home]}"
requested_board="${2-}"
tinka_directory="$project_root/.tinka"
tinka_cli="$tinka_directory/src/cli.mjs"
config_path="$tinka_directory/config.json"

if ! command -v node >/dev/null 2>&1; then
  echo "tinka: node (22 or newer) is required and was not found on PATH." >&2
  exit 1
fi

if [ ! -f "$tinka_cli" ] || [ ! -f "$config_path" ]; then
  echo "tinka: missing $tinka_cli or $config_path — build a canvas first (see skills/do-work-toolbox/docs/ai-report-guide.md#review-on-a-tinka-canvas)." >&2
  exit 1
fi

# Stop this config's own managed server first if a PID record exists.
if [ -f "${config_path}.pid.json" ]; then
  stop_output="$(node "$tinka_cli" stop --config "$config_path" 2>/dev/null || true)"
  if [ -n "$stop_output" ]; then
    printf '%s\n' "$stop_output"
  fi
fi

port="$(node -e '
  const fs = require("node:fs");
  const cfg = JSON.parse(fs.readFileSync(process.argv[1], "utf8"));
  if (!Number.isInteger(cfg.port) || cfg.port < 1 || cfg.port > 65535) process.exit(1);
  process.stdout.write(String(cfg.port));
' "$config_path")" || {
  echo "tinka: could not read a valid integer port from $config_path" >&2
  exit 1
}

# Kill any remaining process blocking TCP $port so Tinka binds its configured port.
if command -v lsof >/dev/null 2>&1; then
  listener_pids="$(lsof -nP -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null || true)"
  if [ -n "$listener_pids" ]; then
    while IFS= read -r listener_pid; do
      [ -z "$listener_pid" ] || [ "$listener_pid" = "$$" ] && continue
      listener_command="$(ps -ww -p "$listener_pid" -o args= 2>/dev/null || true)"
      echo "tinka: killing process blocking port $port (pid $listener_pid): ${listener_command:-unknown}" >&2
      kill "$listener_pid" 2>/dev/null || true
    done <<< "$listener_pids"

    wait_count=0
    while [ "$wait_count" -lt 30 ]; do
      remaining_pids="$(lsof -nP -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null || true)"
      [ -z "$remaining_pids" ] && break
      sleep 0.1
      wait_count=$((wait_count + 1))
    done

    remaining_pids="$(lsof -nP -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null || true)"
    if [ -n "$remaining_pids" ]; then
      while IFS= read -r stubborn_pid; do
        [ -z "$stubborn_pid" ] || [ "$stubborn_pid" = "$$" ] && continue
        echo "tinka: force-killing stubborn process on port $port (pid $stubborn_pid)" >&2
        kill -9 "$stubborn_pid" 2>/dev/null || true
      done <<< "$remaining_pids"

      wait_count=0
      while [ "$wait_count" -lt 20 ]; do
        still_busy="$(lsof -nP -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null || true)"
        [ -z "$still_busy" ] && break
        sleep 0.1
        wait_count=$((wait_count + 1))
      done
    fi

    final_pids="$(lsof -nP -ti tcp:"$port" -sTCP:LISTEN 2>/dev/null || true)"
    if [ -n "$final_pids" ]; then
      echo "tinka: port $port is still occupied (pid $final_pids); refusing to start." >&2
      exit 1
    fi
  fi
fi

if ! start_output="$(node "$tinka_cli" start --config "$config_path")"; then
  exit 1
fi
printf '%s\n' "$start_output"

target_url="$(node -e '
  const fs = require("node:fs");
  const configPath = process.argv[1];
  const requestedBoard = process.argv[2] || "";
  let port;
  try {
    port = JSON.parse(fs.readFileSync(`${configPath}.pid.json`, "utf8")).port;
  } catch {
    port = JSON.parse(fs.readFileSync(configPath, "utf8")).port;
  }
  if (!Number.isInteger(port) || port < 1 || port > 65535) process.exit(1);
  const baseUrl = `http://127.0.0.1:${port}`;
  (async () => {
    if (requestedBoard === "home") {
      process.stdout.write(`${baseUrl}/`);
      return;
    }
    if (requestedBoard) {
      process.stdout.write(`${baseUrl}/boards/${encodeURIComponent(requestedBoard)}/`);
      return;
    }
    try {
      const response = await fetch(`${baseUrl}/api/boards`, { signal: AbortSignal.timeout(2000) });
      if (response.ok) {
        const payload = await response.json();
        if (Array.isArray(payload.boards) && payload.boards.length > 0 && typeof payload.boards[0].id === "string" && payload.boards[0].id) {
          process.stdout.write(`${baseUrl}/boards/${encodeURIComponent(payload.boards[0].id)}/`);
          return;
        }
      }
    } catch {}
    process.stdout.write(`${baseUrl}/`);
  })();
' "$config_path" "$requested_board")" || {
  echo "tinka: started server, failed to resolve board URL." >&2
  exit 1
}

case "$(uname -s)" in
  Darwin)
    open "$target_url" || echo "tinka: browser launch failed; open $target_url manually." >&2
    ;;
  MINGW*|MSYS*|CYGWIN*)
    rundll32 url.dll,FileProtocolHandler "$target_url" >/dev/null 2>&1 &
    ;;
  *)
    if command -v xdg-open >/dev/null 2>&1; then
      xdg-open "$target_url" >/dev/null 2>&1 &
    else
      echo "tinka: no browser opener found; open $target_url manually." >&2
    fi
    ;;
esac

printf 'Opened %s\n' "$target_url"
