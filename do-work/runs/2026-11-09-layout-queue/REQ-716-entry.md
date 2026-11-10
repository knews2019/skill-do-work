## 0.315.2 — Bootstrap Runs the Packaged Installer; Root tools/ Copies Removed (2026-11-10)

The README install block now runs the installer from `skills/do-work/tools/` inside the downloaded archive, and the five duplicate scripts in this repository's root `tools/` folder are gone.

- `install-do-work-suite.sh --print-bootstrap-command` prints the new `source/skills/do-work/tools/` path.
- Updates from older versions are unaffected: an installed updater never runs a script from the downloaded archive.
- A saved copy of the old README block stops with "No such file or directory" before it writes anything. Copy the current README block.
