#!/usr/bin/env python3
"""Bounded read-only CLI observation. Replaces only evidence in this directory."""
import datetime
import json
import os
from pathlib import Path
import selectors
import signal
import subprocess
import time

PROBE_DIRECTORY = Path(__file__).resolve().parent
REPOSITORY_ROOT = PROBE_DIRECTORY.parents[3]
PROGRESS_PATH = Path('/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-11-031446/REQ-718-claude-progress.log')
CLI_ARGUMENTS = ['claude', '--safe-mode', '--setting-sources', '', '--strict-mcp-config',
                 '--tools', 'Read', '--allowedTools', 'Read', '--permission-mode', 'dontAsk',
                 '--no-session-persistence', '--input-format', 'stream-json',
                 '--output-format', 'stream-json', '--verbose', '--max-budget-usd', '3', '-p']


def record_progress(message):
    with PROGRESS_PATH.open('a') as output:
        output.write(datetime.datetime.now(datetime.timezone.utc).isoformat() + ' ' + message + '\n')


def sanitize_event(event):
    """Remove unrelated account/runtime metadata, preserving behavioral evidence."""
    if isinstance(event, dict):
        return {key: sanitize_event(value) for key, value in event.items()
                if key not in {'credential_expires_at', 'messaging_socket_path'}}
    if isinstance(event, list):
        return [sanitize_event(value) for value in event]
    return event


def run_probe():
    record_progress('Starting bounded Claude read-only probe; each turn <=150s, total <=360s.')
    metadata = {'argv': CLI_ARGUMENTS, 'cwd': str(REPOSITORY_ROOT),
                'started_at': datetime.datetime.now(datetime.timezone.utc).isoformat(),
                'cli_version': subprocess.check_output(['claude', '--version'], text=True).strip()}
    process = subprocess.Popen(CLI_ARGUMENTS, cwd=REPOSITORY_ROOT, stdin=subprocess.PIPE,
                               stdout=subprocess.PIPE, stderr=subprocess.PIPE, start_new_session=True)
    selector = selectors.DefaultSelector()
    selector.register(process.stdout, selectors.EVENT_READ, 'stdout')
    selector.register(process.stderr, selectors.EVENT_READ, 'stderr')
    buffers = {'stdout': b'', 'stderr': b''}
    started = time.monotonic()
    observed_boundary = False
    completed_turns = []
    try:
        with (PROBE_DIRECTORY / 'live-events.jsonl').open('w') as events:
            for phase, filename in [('baseline', 'baseline-prompt.txt'), ('compact', 'compact-prompt.txt'),
                                    ('continuation', 'continuation-prompt.txt')]:
                if phase == 'continuation' and not observed_boundary:
                    metadata['blocker'] = 'No system compact_boundary event observed; post-compaction check not run.'
                    break
                prompt = (PROBE_DIRECTORY / filename).read_text().strip()
                request = {'type': 'user', 'message': {'role': 'user', 'content': prompt}}
                process.stdin.write((json.dumps(request) + '\n').encode())
                process.stdin.flush()
                record_progress('Sent ' + phase + ' prompt.')
                deadline = min(started + 360, time.monotonic() + 150)
                completed = False
                while time.monotonic() < deadline and not completed:
                    for key, _ in selector.select(timeout=1):
                        channel = key.data
                        chunk = os.read(key.fileobj.fileno(), 65536)
                        if not chunk:
                            selector.unregister(key.fileobj)
                            continue
                        buffers[channel] += chunk
                        while b'\n' in buffers[channel]:
                            line, buffers[channel] = buffers[channel].split(b'\n', 1)
                            if not line:
                                continue
                            if channel == 'stderr':
                                # Keep stderr scoped; CLI can emit operational errors here.
                                event = {'type': 'stderr', 'text': line.decode(errors='replace')}
                            else:
                                try:
                                    event = json.loads(line)
                                except ValueError:
                                    event = {'type': 'non_json_stdout', 'text': line.decode(errors='replace')}
                            if event.get('type') == 'rate_limit_event':
                                continue
                            event = sanitize_event(event)
                            events.write(json.dumps({'phase': phase, 'event': event}) + '\n')
                            events.flush()
                            if event.get('type') == 'system' and event.get('subtype') == 'compact_boundary':
                                observed_boundary = True
                            if event.get('type') == 'result':
                                completed = True
                                completed_turns.append({'phase': phase, 'is_error': event.get('is_error'),
                                                        'subtype': event.get('subtype')})
                    if process.poll() is not None and not selector.get_map():
                        metadata['blocker'] = 'CLI exited before completing ' + phase
                        break
                if not completed:
                    metadata.setdefault('blocker', 'Turn deadline exceeded: ' + phase)
                    break
                if completed_turns[-1]['is_error']:
                    metadata['blocker'] = 'CLI returned error during ' + phase
                    break
    finally:
        if process.stdin and not process.stdin.closed:
            process.stdin.close()
        try:
            process.wait(timeout=3)
        except subprocess.TimeoutExpired:
            os.killpg(process.pid, signal.SIGTERM)
            try:
                process.wait(timeout=3)
            except subprocess.TimeoutExpired:
                os.killpg(process.pid, signal.SIGKILL)
                process.wait()
        selector.close()
        metadata.update({'process_exit_code': process.returncode, 'elapsed_seconds': round(time.monotonic()-started, 3),
                         'observed_compact_boundary': observed_boundary, 'completed_turns': completed_turns,
                         'process_reaped': True})
        (PROBE_DIRECTORY / 'run-metadata.json').write_text(json.dumps(metadata, indent=2) + '\n')
        record_progress('Probe ended: ' + json.dumps(metadata))
    print(json.dumps(metadata, indent=2))


if __name__ == '__main__':
    run_probe()
