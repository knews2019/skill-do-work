#!/usr/bin/env python3
"""Regenerate release payloads from the LIVE tree right before finalization.
usage: prep-payloads.py <repo-root> <out-dir> <entry-file> ; prints OLD NEW versions."""
import sys, os, re, hashlib
root, out, entry_file = sys.argv[1:4]
old = open(os.path.join(root, 'VERSION')).read().strip()
major, minor, patch = old.split('.')
new = f"{major}.{minor}.{int(patch)+1}"
entry = open(entry_file).read()
assert entry.startswith('## ' + new + ' '), f"entry must start with '## {new} ', got {entry[:40]!r}"
if not entry.endswith('\n\n'):
    entry = entry.rstrip('\n') + '\n\n'
anchor = "CHANGELOG-2026-04-07-up-to-v0.49.0.md)\n\n"
def w(name, data):
    p = os.path.join(out, name); open(p, 'w').write(data); return p
for mirror in ('VERSION', 'skills/do-work/VERSION'):
    live = open(os.path.join(root, mirror)).read()
    assert live.strip() == old, mirror
w('VERSION.old', open(os.path.join(root, 'VERSION')).read())
w('VERSION.new', new + '\n')
vm = open(os.path.join(root, 'skills/do-work/actions/version.md')).read()
assert f"**Current version**: {old}" in vm
w('version-md.old', vm)
w('version-md.new', vm.replace(f"**Current version**: {old}", f"**Current version**: {new}", 1))
cl = open(os.path.join(root, 'CHANGELOG.md')).read()
assert cl == open(os.path.join(root, 'skills/do-work/CHANGELOG.md')).read(), 'changelog mirror differs'
assert cl.count(anchor) == 1, 'anchor not unique'
assert ('## ' + new + ' ') not in cl
w('CHANGELOG.old', cl)
w('CHANGELOG.new', cl.replace(anchor, anchor + entry, 1))
print(old, new)
