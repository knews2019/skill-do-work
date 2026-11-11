#!/usr/bin/env python3
"""Validate this capture against repository inputs without another model call."""
import ast
import hashlib
import json
from pathlib import Path
import re

PROBE_DIRECTORY = Path(__file__).resolve().parent
REPOSITORY_ROOT = PROBE_DIRECTORY.parents[3]


def verify_evidence():
    ast.parse((PROBE_DIRECTORY / 'run-probe.py').read_text())
    rows = [json.loads(line) for line in (PROBE_DIRECTORY / 'live-events.jsonl').read_text().splitlines()]
    reads, results, tool_calls, read_checks = {}, {}, {}, []
    for row_number, row in enumerate(rows, 1):
        event, phase = row['event'], row['phase']
        content = event.get('message', {}).get('content', [])
        if not isinstance(content, list):
            continue
        for block in content:
            if event.get('type') == 'assistant' and block.get('type') == 'tool_use':
                assert block['name'] == 'Read'
                tool_calls[block['id']] = block['input']['file_path']
                reads.setdefault(phase, []).append(block['input']['file_path'])
            if event.get('type') == 'user' and block.get('type') == 'tool_result':
                assert not block.get('is_error')
                path = Path(tool_calls[block['tool_use_id']])
                relative_path = path.relative_to(Path(rows[0]['event']['cwd']))
                # Read numbers the final empty line after the source's trailing newline.
                expected = (REPOSITORY_ROOT / relative_path).read_text().split('\n')
                actual = re.findall(r'(?m)^(\d+)\t(.*)$', block['content'])
                assert [int(number) for number, _ in actual] == list(range(1, len(expected) + 1))
                assert [text for _, text in actual] == expected, relative_path
                read_checks.append({'file': str(relative_path), 'numbered_lines': len(expected),
                                    'tool_result_row': row_number, 'full_content_matches': True})
    for row in rows:
        event = row['event']
        if event.get('type') == 'result':
            assert not event.get('is_error')
            results[row['phase']] = event.get('result')
    assert len(read_checks) == 8
    assert reads.get('continuation', []) == []
    boundaries = [row['event'] for row in rows if row['event'].get('subtype') == 'compact_boundary']
    assert len(boundaries) == 1 and boundaries[0]['compact_metadata']['trigger'] == 'manual'
    assert 'function sumInvoiceAmounts(' in results['baseline']
    assert 'function sumShipmentWeights(' in results['continuation']
    summary = next(row['event']['message']['content'] for row in rows if row['event'].get('isSynthetic'))
    assert 'Naming for Reach' in summary and 'two-word minimum' in summary
    fingerprints = json.loads((PROBE_DIRECTORY / 'input-fingerprints.json').read_text())
    assert all(hashlib.sha256((REPOSITORY_ROOT / key).read_bytes()).hexdigest() == value
               for key, value in fingerprints.items())
    assert not re.search(r'credential_expires_at|messaging_socket_path|rate_limit_event|sk-ant-',
                         (PROBE_DIRECTORY / 'live-events.jsonl').read_text())
    metadata = json.loads((PROBE_DIRECTORY / 'run-metadata.json').read_text())
    assert metadata['process_reaped'] and metadata['process_exit_code'] == 0
    checks = {'python_syntax': 'PASS', 'full_read_results': read_checks,
              'manual_compact_boundary': 'PASS', 'baseline_function': 'sumInvoiceAmounts',
              'continuation_function': 'sumShipmentWeights', 'continuation_read_calls': 0,
              'summary_retains_naming_rule': True, 'source_fingerprints': 'PASS',
              'unrelated_metadata_removed': 'PASS', 'cli_exit_code': 0, 'process_reaped': True}
    (PROBE_DIRECTORY / 'focused-checks.json').write_text(json.dumps(checks, indent=2) + '\n')
    print(json.dumps(checks, indent=2))


if __name__ == '__main__':
    verify_evidence()
