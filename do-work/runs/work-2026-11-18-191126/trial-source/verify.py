"""Verify locally, or re-extract the frozen rows with --repo REPOSITORY."""
import argparse
import csv
import hashlib
import io
from collections import Counter, defaultdict
from decimal import Decimal as D
from pathlib import Path
from statistics import median
from subprocess import check_output

REV = 'a1b622fb2fb3b6fe4a7b5c0d8ed0d2e300f7e46a'
RUN = 'standalone-20260906T134958Z-1358'
HERE = Path(__file__).resolve().parent
args = argparse.ArgumentParser(description=__doc__)
args.add_argument('--repo', type=Path)
repo = args.parse_args().repo
data = (HERE / 'selected.tsv').read_bytes()
if repo:
    header, *lines = check_output(
        ['git', '-C', str(repo), 'show', f'{REV}:do-work/test-efficiency.tsv']
    ).splitlines(keepends=True)
    selected = [(i, line) for i, line in enumerate(lines, 2)
                if line.split(b'\t', 1)[0] == RUN.encode()]
    assert [i for i, _ in selected] == list(range(11, 38))
    assert header + b''.join(line for _, line in selected) == data
assert hashlib.sha256(data).hexdigest() == (
    'f9ae5ae0c2c3fd8371222ea44104f6096373d8a5e22e6fcc2f61e35ca8fa6fcd'
)

# Literal reference medians (wall, total CPU, Go accumulated) and count vectors.
expected = {
    ('inventory', 'cold'): ('8.973 1.629 8.370', 'bash:2,git:18,go:1'),
    ('inventory', 'warm'): ('9.242 1.756 8.540', 'bash:2,git:18,go:1'),
    ('finalization', 'cold'): ('48.606 37.840 47.780', 'git:3019,go:1'),
    ('finalization', 'warm'): ('50.524 38.891 48.860', 'git:3019,go:1'),
    ('session-start', 'warm'): ('4.534 1.207 n/a', 'bash:15,find:10,git:27,go:15,shasum:24'),
    ('shell-audit', 'warm'): ('5.849 3.253 n/a', 'bash:111,find:2,git:1,shellcheck:108'),
    ('heavy-cli-build', 'cold'): ('7.810 4.375 7.180', 'bash:15,git:126,go:1'),
    ('heavy-cli-build', 'warm'): ('7.803 4.364 7.160', 'bash:15,git:126,go:1'),
    ('reused-stage', 'reused'): ('3.924 1.185 n/a', 'bash:48,git:69,go:46'),
}
metrics = ('wall_seconds', 'user_cpu_seconds', 'sys_cpu_seconds',
           'total_cpu_seconds', 'go_accum_seconds')
groups = defaultdict(list)
totals = Counter()
print('Source lines and identities (line, case, condition, observation within group)')
for line, row in enumerate(csv.DictReader(io.StringIO(data.decode()), delimiter='\t'), 11):
    assert tuple(row[k] for k in ('run_id', 'revision', 'toolchain', 'concurrency', 'exit_status')) == (
        RUN, '12370cdb', 'go1.26.1,py:3.9.6', '4', '0')
    key = row['case_name'], row['cache_condition']
    assert row['subprocess_counts'] == expected[key][1]
    values = [None if row[k] == 'n/a' else D(row[k]) for k in metrics]
    assert abs(values[1] + values[2] - values[3]) <= D('.001')
    groups[key].append(values)
    totals.update({name: int(count) for name, count in
                   (part.split(':') for part in row['subprocess_counts'].split(','))})
    print(line, *key, len(groups[key]))
assert set(groups) == set(expected)
assert all(len(rows) == 3 for rows in groups.values())
assert totals == dict(bash=624, find=36, git=19269, go=201, shasum=72, shellcheck=324)

print('\nSeconds by metric (median [min, max], spread=max-min)')
medians = {}
for key, rows in groups.items():
    medians[key] = [None if column[0] is None else median(column) for column in zip(*rows)]
    assert ' '.join(str(medians[key][i]) if medians[key][i] is not None else 'n/a'
                    for i in (0, 3, 4)) == expected[key][0]
    print(*key)
    for label, mid, column in zip(metrics, medians[key], zip(*rows)):
        if mid is None:
            assert all(value is None for value in column)
            print(f'  {label} n/a')
        else:
            print(f'  {label} {mid} [{min(column)}, {max(column)}], spread {max(column)-min(column)}')

print('\nWarm minus cold medians in seconds (percent relative to cold)')
for case, differences in {
    'inventory': '.269 .041 .085 .127 .170',
    'finalization': '1.918 .534 .517 1.051 1.080',
    'heavy-cli-build': '-.007 -.003 -.004 -.011 -.020',
}.items():
    cold, warm = (medians[case, condition] for condition in ('cold', 'warm'))
    delta = [w - c for c, w in zip(cold, warm)]
    assert delta == list(map(D, differences.split()))
    print(case, ', '.join(f'{name} {diff:+.3f} ({100*diff/base:+.2f}%)'
                         for name, diff, base in zip(metrics, delta, cold)))
print('\nIntercepted totals across 27 rows:', dict(sorted(totals.items())))
print('PASS: 27 rows, nine groups of three, row identities, counts, medians and deltas')
